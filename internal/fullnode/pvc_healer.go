package fullnode

import (
	"context"
	"errors"
	"github.com/samber/lo"
	"math"
	"time"

	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type StatusSyncer interface {
	SyncUpdate(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error
}

type PVCHealer struct {
	client StatusSyncer
	now    func() time.Time
}

func NewPVCHealer(client StatusSyncer) *PVCHealer {
	return &PVCHealer{
		client: client,
		now:    time.Now,
	}
}

// SignalPVCResize patches the CosmosFullNode.status.selfHealing with the new calculated PVC size as a resource quantity.
// Assumes CosmosfullNode.spec.selfHealing.pvcAutoScaling is set or else this method may panic.
// The CosmosFullNode controller is responsible for increasing the PVC disk size.
//
// Returns true if the status was patched.
//
// Returns false and does not patch if:
// 1. The PVCs do not need resizing
// 2. The status already has >= calculated size.
// 3. The maximum size has been reached. It will patch up to the maximum size.
//
// Returns an error if patching unsuccessful.
func (healer PVCHealer) SignalPVCResize(ctx context.Context, crd *cosmosv1.CosmosFullNode, results []PVCDiskUsage) (bool, error) {
	var (
		spec    = crd.Spec.SelfHeal.PVCAutoScale
		trigger = int(spec.UsedSpacePercentage)
	)

	var joinedErr error

	status := crd.Status.SelfHealing.PVCAutoScale

	patches := make(map[string]*cosmosv1.PVCAutoScaleStatus)

	now := metav1.NewTime(healer.now())

	for _, pvc := range results {
		if pvc.PercentUsed < trigger {
			// no need to expand
			continue
		}

		newSize, err := healer.calcNextCapacity(pvc.Capacity, spec.IncreaseQuantity)
		if err != nil {
			joinedErr = errors.Join(joinedErr, err)
			continue
		}

		if status != nil {
			if pvcStatus, ok := status[pvc.Name]; ok && pvcStatus.RequestedSize.Value() == newSize.Value() {
				// already requested
				continue
			}
		}

		if max := spec.MaxSize; !max.IsZero() {
			if pvc.Capacity.Cmp(max) >= 0 {
				// already at max size
				continue
			}

			if newSize.Cmp(max) >= 0 {
				// Cap new size to the max size
				newSize = max
			}
		}

		patches[pvc.Name] = &cosmosv1.PVCAutoScaleStatus{
			RequestedSize: newSize,
			RequestedAt:   now,
		}
	}

	if len(patches) == 0 {
		return false, joinedErr
	}

	return true, errors.Join(joinedErr, healer.client.SyncUpdate(ctx, client.ObjectKeyFromObject(crd), func(status *cosmosv1.FullNodeStatus) {
		if status.SelfHealing.PVCAutoScale == nil {
			status.SelfHealing.PVCAutoScale = patches
			return
		}
		for k, v := range patches {
			status.SelfHealing.PVCAutoScale[k] = v
		}
	}))
}

func (healer PVCHealer) calcNextCapacity(current resource.Quantity, increase string) (resource.Quantity, error) {
	var (
		merr     error
		quantity resource.Quantity
	)

	// Try to calc by percentage first
	v := intstr.FromString(increase)
	percent, err := intstr.GetScaledValueFromIntOrPercent(&v, 100, false)
	if err == nil {
		addtl := math.Round(float64(current.Value()) * (float64(percent) / 100.0))
		quantity = *resource.NewQuantity(current.Value()+int64(addtl), current.Format)
		return quantity, nil
	}

	merr = errors.Join(merr, err)

	// Then try to calc by resource quantity
	addtl, err := resource.ParseQuantity(increase)
	if err != nil {
		return quantity, errors.Join(merr, err)
	}

	return *resource.NewQuantity(current.Value()+addtl.Value(), current.Format), nil
}

func (healer PVCHealer) UpdatePodFailure(ctx context.Context, crd *cosmosv1.CosmosFullNode, podName string) (bool, error) {
	var regenPVCStatus map[string]*cosmosv1.RegenPVCStatus
	if crd.Status.SelfHealing.RegenPVCStatus != nil {
		regenPVCStatus = crd.Status.SelfHealing.RegenPVCStatus
		if regenPVCStatus[podName].Phase != nil && *regenPVCStatus[podName].Phase == cosmosv1.RegenPVCPhaseRegeneratingPVC {
			return false, nil
		}
	}

	currentRegenPVCStatus := regenPVCStatus[podName]

	now := metav1.NewTime(healer.now())

	if currentRegenPVCStatus != nil {
		currentRegenPVCStatus.FailureTimes = lo.FilterMap(currentRegenPVCStatus.FailureTimes, func(item metav1.Time, index int) (metav1.Time, bool) {
			collectionDuration := crd.Spec.SelfHeal.HeightDriftMitigation.RegeneratePVC.FailedCountCollectionDuration
			if (item.Add(collectionDuration.Duration)).Before(now.Time) {
				return metav1.Time{}, false
			}
			return item, true
		})
	} else {
		currentRegenPVCStatus = new(cosmosv1.RegenPVCStatus)
		currentRegenPVCStatus.Phase = ptr(cosmosv1.RegenPVCPhaseNotYet)
	}

	currentRegenPVCStatus.FailureTimes = append(currentRegenPVCStatus.FailureTimes, now)

	currentFailureCount := uint32(len(currentRegenPVCStatus.FailureTimes))

	return currentFailureCount > crd.Spec.SelfHeal.HeightDriftMitigation.RegeneratePVC.ThresholdCount, healer.client.SyncUpdate(ctx, client.ObjectKeyFromObject(crd), func(status *cosmosv1.FullNodeStatus) {
		if status.SelfHealing.RegenPVCStatus == nil {
			status.SelfHealing.RegenPVCStatus = map[string]*cosmosv1.RegenPVCStatus{}
		}
		status.SelfHealing.RegenPVCStatus[podName] = currentRegenPVCStatus
	})
}
