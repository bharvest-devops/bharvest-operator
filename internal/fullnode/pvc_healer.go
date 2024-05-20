package fullnode

import (
	"context"
	"errors"
	"fmt"
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

// UpdatePodFailure updates how many times failed on this pod, and also check if this pod should be entered Re-gen pvc step.
func (healer PVCHealer) UpdatePodFailure(ctx context.Context, crd *cosmosv1.CosmosFullNode, podName string) (bool, error) {
	var joinedErr error
	if crd == nil {
		return false, errors.Join(joinedErr, fmt.Errorf("provided CosmosFullNode is nil"))
	}

	regenPVCStatus := crd.Status.SelfHealing.RegenPVCStatus
	if regenPVCStatus == nil {
		regenPVCStatus = ptr(cosmosv1.RegenPVCStatus{
			RegenPVCPhase: cosmosv1.RegenPVCPhaseNotYet,
		})
		crd.Status.SelfHealing.RegenPVCStatus = regenPVCStatus
	}

	current, _ := regenPVCStatus.FailureTimes[podName]

	now := metav1.NewTime(healer.now())

	if current != nil {
		current = lo.FilterMap(current, func(failureTime string, index int) (string, bool) {
			collectionDuration := crd.Spec.SelfHeal.HeightDriftMitigation.RegeneratePVC.FailedCountCollectionDuration
			t, err := time.Parse(failureTime, "2006-01-02 15:04:05")
			if err != nil {
				return "", false
			}
			if (t.Add(collectionDuration.Duration)).Before(now.Time) {
				return "", false
			}
			return failureTime, true
		})
	} else {
		current = []string{}
		regenPVCStatus.RegenPVCPhase = cosmosv1.RegenPVCPhaseNotYet
	}

	current = append(current, now.Format("2006-01-02 15:04:05"))

	isOveredRegenThreshold := uint32(len(current)) > crd.Spec.SelfHeal.HeightDriftMitigation.RegeneratePVC.ThresholdCount
	if isOveredRegenThreshold {
		regenPVCStatus.RegenPVCPhase = cosmosv1.RegenPVCPhaseRegeneratingPVC
		current = []string{}
	}
	regenPVCStatus.FailureTimes[podName] = current

	return isOveredRegenThreshold, errors.Join(joinedErr, healer.client.SyncUpdate(ctx, client.ObjectKeyFromObject(crd), func(status *cosmosv1.FullNodeStatus) {
		status.SelfHealing.RegenPVCStatus = regenPVCStatus
	}))
}
