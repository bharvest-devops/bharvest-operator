package prune

import (
	"context"
	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	"github.com/bharvest-devops/cosmos-operator/internal/fullnode"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type CandidateCollector interface {
	SyncedPods(ctx context.Context, controller client.ObjectKey) []*corev1.Pod
}

type Pruner struct {
	candidateCollector CandidateCollector
}

func NewPruner(candidateCollector CandidateCollector) *Pruner {
	return &Pruner{
		candidateCollector: candidateCollector,
	}
}

func (p *Pruner) FindCandidate(ctx context.Context, crd *cosmosv1.CosmosFullNode, results []fullnode.PVCDiskUsage) *corev1.Pod {
	var spec = crd.Spec.SelfHeal.PruningSpec
	if spec == nil {
		// Pruning not work
		return nil
	}

	var trigger = int(spec.UsedSpacePercentage)

	status := crd.Status.SelfHealing.CosmosPruningStatus

	if status == nil {
		status = new(cosmosv1.CosmosPruningStatus)
	}

	cctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var (
		synced     = p.candidateCollector.SyncedPods(cctx, client.ObjectKey{Namespace: crd.Namespace, Name: crd.Name})
		availCount = int32(len(synced))
		minAvail   = spec.MinAvailable
	)

	if minAvail <= 0 {
		minAvail = 2
	}

	if availCount <= minAvail {
		return nil
	}

	for _, pvc := range results {
		if pvc.PercentUsed < trigger {
			// no need to prune
			continue
		}

		// Finding candidate
		for _, pod := range synced {
			if fullnode.PVCName(pod) != pvc.Name {
				continue
			}
			return pod
		}
	}
	return nil
}
