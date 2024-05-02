package fullnode

import (
	"context"
	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	"github.com/bharvest-devops/cosmos-operator/internal/cosmos"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ResetStatus is used at the beginning of the reconcile loop.
// It resets the crd's status to a fresh state.
func ResetStatus(crd *cosmosv1.CosmosFullNode) {
	crd.Status.ObservedGeneration = crd.Generation
	crd.Status.Phase = cosmosv1.FullNodePhaseProgressing
	crd.Status.StatusMessage = nil
}

type StatusCollector interface {
	Collect(ctx context.Context, controller client.ObjectKey) cosmos.StatusCollection
}

// SyncInfoStatus returns the status of the full node's sync info.
func SyncInfoStatus(
	ctx context.Context,
	crd *cosmosv1.CosmosFullNode,
	collector StatusCollector,
) map[string]*cosmosv1.SyncInfoPodStatus {
	status := make(map[string]*cosmosv1.SyncInfoPodStatus, crd.Spec.Replicas)

	coll := collector.Collect(ctx, client.ObjectKeyFromObject(crd))

	for _, item := range coll {
		podName := item.GetPod().Name
		var (
			stat       cosmosv1.SyncInfoPodStatus
			beforeStat = crd.Status.SyncInfo[podName]
		)

		stat.Timestamp = metav1.NewTime(item.Timestamp())

		comet, err := item.GetStatus()
		if err != nil {
			stat.Error = ptr(err.Error())
			if beforeStat != nil {
				stat.LastBlockTimestamp = beforeStat.LastBlockTimestamp
			} else {
				stat.LastBlockTimestamp = stat.Timestamp
			}
		} else {

			stat.Height = ptr(comet.LatestBlockHeight())
			stat.InSync = ptr(!comet.Result.SyncInfo.CatchingUp)

			if beforeStat != nil &&
				beforeStat.Height != nil &&
				stat.Height != nil &&
				*beforeStat.Height == *stat.Height &&
				beforeStat.LastBlockTimestamp != *new(metav1.Time) {
				stat.LastBlockTimestamp = beforeStat.LastBlockTimestamp
			} else {
				stat.LastBlockTimestamp = stat.Timestamp
			}
		}

		retainDuration := metav1.Duration{
			Duration: stat.Timestamp.Sub(stat.LastBlockTimestamp.Time),
		}

		stat.HeightRetainTime = &retainDuration
		status[podName] = &stat
	}

	return status
}
