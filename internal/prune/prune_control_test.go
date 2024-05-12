package prune

import (
	"context"
	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	"github.com/bharvest-devops/cosmos-operator/internal/fullnode"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

type MockCandidateCollector func(ctx context.Context, controller client.ObjectKey) []*corev1.Pod

func (fn MockCandidateCollector) SyncedPods(ctx context.Context, controller client.ObjectKey) []*corev1.Pod {
	if ctx == nil {
		panic("nil context")
	}
	return fn(ctx, controller)
}

var nopCollector = MockCandidateCollector(func(ctx context.Context, controller client.ObjectKey) []*corev1.Pod {
	return nil
})

func TestPruneControl_FindCandidate(t *testing.T) {
	ctx := context.Background()

	var crd = cosmosv1.CosmosFullNode{
		Spec: cosmosv1.FullNodeSpec{
			SelfHeal: ptr(cosmosv1.SelfHealSpec{
				PruningSpec: ptr(cosmosv1.PruningSpec{
					UsedSpacePercentage: 50,
				}),
			}),
		},
	}

	var defaultResults = []fullnode.PVCDiskUsage{
		{
			Name:        "pvc-cosmoshub-0",
			PercentUsed: 10,
			Capacity:    resource.MustParse("100Gi"),
		},
		{
			Name:        "pvc-cosmoshub-1",
			PercentUsed: 90,
			Capacity:    resource.MustParse("100Gi"),
		},
	}

	t.Run("happy path", func(t *testing.T) {
		cacheController := MockCandidateCollector(
			func(ctx context.Context, controller client.ObjectKey) []*corev1.Pod {
				return []*corev1.Pod{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "pod",
						},
						Spec: corev1.PodSpec{
							Volumes: []corev1.Volume{
								{
									Name: "vol-chain-home",
								},
							},
						},
					},
				}
			},
		)

		pruner := NewPruner(cacheController)

		pod, err := pruner.FindCandidate(ctx, ptr(crd), defaultResults)

		require.NoError(t, err)
		require.Equal(t, pod.Name, "cosmoshub")
	})

	t.Run("find failed", func(t *testing.T) {
		cacheController := MockCandidateCollector(
			func(ctx context.Context, controller client.ObjectKey) []*corev1.Pod {
				return []*corev1.Pod{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "pod",
						},
						Spec: corev1.PodSpec{
							Volumes: []corev1.Volume{
								{
									Name: "vol-chain-home",
								},
							},
						},
					},
				}
			},
		)

		noExceededDiskUsages := []fullnode.PVCDiskUsage{
			{
				Name:        "pvc-cosmoshub-0",
				PercentUsed: 40,
				Capacity:    resource.MustParse("100Gi"),
			},
			{
				Name:        "pvc-cosmoshub-1",
				PercentUsed: 30,
				Capacity:    resource.MustParse("100Gi"),
			},
		}

		pruner := NewPruner(cacheController)

		pod, err := pruner.FindCandidate(ctx, ptr(crd), noExceededDiskUsages)

		require.NoError(t, err)
		require.Equal(t, pod, nil)
	})

}
