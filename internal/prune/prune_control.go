package prune

import (
	"context"
	"errors"
	"fmt"
	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	"github.com/bharvest-devops/cosmos-operator/internal/fullnode"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const cosmosSourceLabel = "cosmos.bharvest/source"

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

func (p *Pruner) FindCandidate(ctx context.Context, crd *cosmosv1.CosmosFullNode, results []fullnode.PVCDiskUsage) (*corev1.Pod, error) {
	var (
		spec    = crd.Spec.SelfHeal.PruningSpec
		trigger = int(spec.UsedSpacePercentage)
	)

	var joinedErr error

	status := crd.Status.SelfHealing.CosmosPruningStatus

	cctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var (
		synced     = p.candidateCollector.SyncedPods(cctx, client.ObjectKey{Namespace: crd.Namespace, Name: crd.Name})
		availCount = int32(len(synced))
	)

	if availCount <= 0 {
		return nil, errors.Join(joinedErr, fmt.Errorf("there are no available pods to prune. pruning must be preceed for synced pods"))
	}

	for _, pvc := range results {
		if pvc.PercentUsed < trigger {
			// no need to prune
			continue
		}

		if status != nil {
			// Finding candidate
			for _, pod := range synced {
				if fullnode.PVCName(pod) != pvc.Name {
					continue
				}
				return pod, nil
			}
		}
	}
	return nil, nil
}
