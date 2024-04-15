package fullnode

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	"github.com/bharvest-devops/cosmos-operator/internal/cosmos"
	"github.com/bharvest-devops/cosmos-operator/internal/kube"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DriftDetection detects pods that are lagging behind the latest block height.
type DriftDetection struct {
	available      func(pods []*corev1.Pod, minReady time.Duration, now time.Time) []*corev1.Pod
	collector      StatusCollector
	computeRollout func(maxUnavail *intstr.IntOrString, desired, ready int) int
}

func NewDriftDetection(collector StatusCollector) DriftDetection {
	return DriftDetection{
		available:      kube.AvailablePods,
		collector:      collector,
		computeRollout: kube.ComputeRollout,
	}
}

// LaggingPods returns pods that are lagging behind the latest block height.
func (d DriftDetection) LaggingPods(ctx context.Context, crd *cosmosv1.CosmosFullNode) []*corev1.Pod {
	var lagging []*corev1.Pod
	synced := d.collector.Collect(ctx, client.ObjectKeyFromObject(crd)).Synced()

	lagging = lo.FilterMap(synced, func(item cosmos.StatusItem, _ int) (*corev1.Pod, bool) {
		var initDuration metav1.Duration
		if item.HeightRetainTime == initDuration {
			return item.GetPod(), false
		}
		isLagging := item.HeightRetainTime.Duration >= crd.Spec.SelfHeal.HeightDriftMitigation.ThresholdTime.Duration
		return item.GetPod(), isLagging
	})
	if len(lagging) > 0 {
		return lagging
	}

	maxHeight := lo.MaxBy(synced, func(a cosmos.StatusItem, b cosmos.StatusItem) bool {
		return a.Status.LatestBlockHeight() > b.Status.LatestBlockHeight()
	}).Status.LatestBlockHeight()

	thresh := uint64(crd.Spec.SelfHeal.HeightDriftMitigation.ThresholdHeight)
	lagging = lo.FilterMap(synced, func(item cosmos.StatusItem, _ int) (*corev1.Pod, bool) {
		isLagging := maxHeight-item.Status.LatestBlockHeight() >= thresh
		return item.GetPod(), isLagging
	})

	avail := d.available(synced.Pods(), 5*time.Second, time.Now())
	rollout := d.computeRollout(crd.Spec.RolloutStrategy.MaxUnavailable, int(crd.Spec.Replicas), len(avail))
	return lo.Slice(lagging, 0, rollout)
}
