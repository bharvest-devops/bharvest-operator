package prune

import (
	"context"
	"errors"
	"fmt"
	"github.com/bharvest-devops/cosmos-operator/internal/fullnode"
	"testing"

	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type mockStatusSyncer func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error

func (fn mockStatusSyncer) SyncUpdate(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
	if ctx == nil {
		panic("nil context")
	}
	return fn(ctx, key, update)
}

var nopStatusSyncer = mockStatusSyncer(func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
	return nil
})

type mockReader struct {
	Lister func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error
	Getter func(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error
}

func (m mockReader) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	if ctx == nil {
		panic("nil context")
	}
	if len(opts) > 0 {
		panic("unexpected opts")
	}
	if m.Getter == nil {
		panic("get called with no implementation")
	}
	return m.Getter(ctx, key, obj, opts...)
}

func (m mockReader) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	if ctx == nil {
		panic("nil context")
	}
	if m.Lister == nil {
		panic("list called with no implementation")
	}
	return m.Lister(ctx, list, opts...)
}

var nopReader = mockReader{
	Lister: func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error { return nil },
	Getter: func(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
		return nil
	},
}

func TestFullNodeControl_SignalPodReplace(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var crd cosmosv1.CosmosFullNode
	crd.Namespace = "default" // Tests for slash stripping.
	crd.Name = "cosmoshub"

	t.Run("happy path", func(t *testing.T) {
		syncer := mockStatusSyncer(func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
			require.Equal(t, "default/cosmoshub", key.String())

			var got cosmosv1.FullNodeStatus

			update(&got)

			candidateKey := "default.cosmoshub-0.v1.cosmos.bharvest"

			require.Equal(t,
				cosmosv1.SelfHealingCandidate{
					PodName:   "cosmoshub-0",
					Namespace: "default",
				},
				got.SelfHealing.CosmosPruningStatus.Candidates[candidateKey])

			return nil
		})

		reader := mockReader{
			Getter: func(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return nil
			},
			Lister: func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
				return nil
			},
		}

		control := NewFullNodeControl(syncer, reader)
		candidates := []*corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cosmoshub-0",
					Namespace: "default",
				},
			},
		}

		err := control.SignalPodReplace(ctx, &crd, candidates)

		require.NoError(t, err)
	})

	t.Run("signal failed", func(t *testing.T) {
		syncer := mockStatusSyncer(func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
			return errors.New("boom")
		})

		candidates := []*corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cosmoshub-0",
					Namespace: "default",
				},
			},
		}

		control := NewFullNodeControl(syncer, nopReader)
		err := control.SignalPodReplace(ctx, &crd, candidates)

		require.Error(t, err)
		require.EqualError(t, err, "boom")
	})
}

func TestFullNodeControl_ConfirmPodReplaced(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var crd cosmosv1.CosmosFullNode
	crd.Name = "cosmoshub"
	crd.Namespace = "default"

	t.Run("happy path", func(t *testing.T) {

		reader := mockReader{Lister: func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
			list.(*corev1.PodList).Items = []corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: fullnode.GetPrunerPodName("cosmoshub-0"),
					},
				},
			}
			return nil
		}}

		fullNodeControl := NewFullNodeControl(nopStatusSyncer, reader)

		podName := "cosmoshub-0"
		crd.Status.SelfHealing.CosmosPruningStatus = ptr(cosmosv1.CosmosPruningStatus{
			Candidates: map[string]cosmosv1.SelfHealingCandidate{
				fullNodeControl.sourceKey(podName, crd.Namespace): {PodName: podName, Namespace: crd.Namespace},
			},
		})

		err := fullNodeControl.ConfirmPodReplaced(ctx, ptr(crd))

		require.NoError(t, err)
	})

	t.Run("failed path", func(t *testing.T) {
		podName := "cosmoshub-0"

		reader := mockReader{Lister: func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
			list.(*corev1.PodList).Items = []corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: podName,
					},
				},
			}
			return nil
		}}

		fullNodeControl := NewFullNodeControl(nopStatusSyncer, reader)

		crd.Status.SelfHealing.CosmosPruningStatus = ptr(cosmosv1.CosmosPruningStatus{
			Candidates: map[string]cosmosv1.SelfHealingCandidate{
				fullNodeControl.sourceKey(podName, crd.Namespace): {PodName: podName, Namespace: crd.Namespace},
			},
		})

		err := fullNodeControl.ConfirmPodReplaced(ctx, ptr(crd))

		require.Error(t, err, fmt.Sprintf("pod %s not replaced yet", podName))
	})

	t.Run("empty candidates path", func(t *testing.T) {
		podName := "cosmoshub-0"

		reader := mockReader{Lister: func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
			list.(*corev1.PodList).Items = []corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: podName,
					},
				},
			}
			return nil
		}}

		fullNodeControl := NewFullNodeControl(nopStatusSyncer, reader)

		crd.Status.SelfHealing.CosmosPruningStatus = ptr(cosmosv1.CosmosPruningStatus{
			Candidates: map[string]cosmosv1.SelfHealingCandidate{},
		})

		err := fullNodeControl.ConfirmPodReplaced(ctx, ptr(crd))

		require.Error(t, err, errors.New(NO_CANDIDATES_ERR))
	})

	t.Run("get error", func(t *testing.T) {
		var reader mockReader
		reader.Lister = func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
			return errors.New("boom")
		}

		control := NewFullNodeControl(nopStatusSyncer, reader)
		err := control.ConfirmPodReplaced(ctx, &crd)

		require.Error(t, err)
		require.EqualError(t, err, "list pods: boom")
	})
}

func TestFullNodeControl_CheckPruningComplete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var crd cosmosv1.CosmosFullNode
	crd.Namespace = "default"
	crd.Name = "cosmoshub"
	candidatePodName := "cosmoshub-0"

	t.Run("happy path", func(t *testing.T) {

		read := mockReader{Lister: func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
			list.(*corev1.PodList).Items = []corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: fullnode.GetPrunerPodName(candidatePodName),
					},
					Status: corev1.PodStatus{
						ContainerStatuses: []corev1.ContainerStatus{
							{
								State: corev1.ContainerState{
									Terminated: ptr(corev1.ContainerStateTerminated{}),
								},
							},
						},
					},
				},
			}
			return nil
		}}

		fullNodeControl := NewFullNodeControl(nopStatusSyncer, read)

		crd.Status.SelfHealing.CosmosPruningStatus = ptr(cosmosv1.CosmosPruningStatus{
			Candidates: map[string]cosmosv1.SelfHealingCandidate{
				fullNodeControl.sourceKey(candidatePodName, crd.Namespace): {
					PodName:   candidatePodName,
					Namespace: crd.Namespace,
				},
			},
		})

		ok, err := fullNodeControl.CheckPruningComplete(ctx, ptr(crd))
		require.Equal(t, true, ok)

		require.NoError(t, err)
	})

	t.Run("failed path", func(t *testing.T) {

		read := mockReader{Lister: func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
			list.(*corev1.PodList).Items = []corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      fullnode.GetPrunerPodName(candidatePodName),
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Volumes: []corev1.Volume{
							{
								Name: "vol-chain-home",
								VolumeSource: corev1.VolumeSource{
									PersistentVolumeClaim: ptr(corev1.PersistentVolumeClaimVolumeSource{
										ClaimName: "pvc-cosmoshub",
									}),
								},
							},
						},
					},
					Status: corev1.PodStatus{
						ContainerStatuses: []corev1.ContainerStatus{
							{
								State: corev1.ContainerState{
									Waiting:    nil,
									Running:    nil,
									Terminated: nil,
								},
							},
						},
					},
				},
			}

			return nil
		}}

		fullNodeControl := NewFullNodeControl(nopStatusSyncer, read)

		crd.Status.SelfHealing.CosmosPruningStatus = ptr(cosmosv1.CosmosPruningStatus{
			Candidates: map[string]cosmosv1.SelfHealingCandidate{
				fullNodeControl.sourceKey(candidatePodName, crd.Namespace): {
					PodName:   candidatePodName,
					Namespace: crd.Namespace,
				},
			},
		})

		ok, err := fullNodeControl.CheckPruningComplete(ctx, &crd)

		require.NoError(t, err)
		require.Equal(t, false, ok)

	})

}

func TestFullNodeControl_SignalPodRestoration(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var crd cosmosv1.CosmosFullNode
	crd.Namespace = "default"
	crd.Name = "cosmoshub"

	t.Run("happy path", func(t *testing.T) {
		syncer := mockStatusSyncer(func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
			require.Equal(t, "default/cosmoshub", key.String())

			var got cosmosv1.FullNodeStatus

			update(&got)

			candidateKey := "default.cosmoshub-0.v1.cosmos.bharvest"

			require.NotNil(t, got.SelfHealing.CosmosPruningStatus.Candidates[candidateKey])

			return nil
		})

		control := NewFullNodeControl(syncer, nopReader)
		candidates := []*corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cosmoshub-0",
					Namespace: "default",
				},
			},
		}

		err := control.SignalPodReplace(ctx, &crd, candidates)

		require.NoError(t, err)
	})

	t.Run("candidates failed", func(t *testing.T) {

		control := NewFullNodeControl(nopStatusSyncer, nopReader)
		err := control.SignalPodRestoration(ctx, &crd)

		require.Error(t, err)
		require.EqualError(t, err, NO_CANDIDATES_ERR)
	})

	t.Run("signal failed", func(t *testing.T) {
		syncer := mockStatusSyncer(func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
			return errors.New("boom")
		})

		podName := "cosmoshub-0"
		candidateName := "default." + crd.Name + "-0.v1.cosmos.bharvest"

		crd.Status.SelfHealing.CosmosPruningStatus = ptr(cosmosv1.CosmosPruningStatus{
			Candidates: map[string]cosmosv1.SelfHealingCandidate{
				candidateName: {
					PodName:   podName,
					Namespace: "default",
				},
			},
		})

		control := NewFullNodeControl(syncer, nopReader)
		err := control.SignalPodRestoration(ctx, &crd)

		require.Error(t, err)
		require.EqualError(t, err, "boom")
	})
}

func TestFullNodeControl_ConfirmPodRestoration(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var crd cosmosv1.CosmosFullNode
	crd.Namespace = "default"
	crd.Name = "cosmoshub"
	crd.Status.SelfHealing.CosmosPruningStatus = new(cosmosv1.CosmosPruningStatus)
	crd.Status.SelfHealing.CosmosPruningStatus.Candidates = make(map[string]cosmosv1.SelfHealingCandidate)

	t.Run("happy path", func(t *testing.T) {
		var reader mockReader

		candidateKey := "default." + crd.Name + "-0.v1.cosmos.bharvest"
		crd.Status.SelfHealing.CosmosPruningStatus.Candidates[candidateKey] = cosmosv1.SelfHealingCandidate{PodName: fullnode.GetPrunerPodName(crd.Name + "-0"), Namespace: crd.Namespace}

		reader.Lister = func(_ context.Context, list client.ObjectList, opts ...client.ListOption) error {
			list.(*corev1.PodList).Items = []corev1.Pod{
				{ObjectMeta: metav1.ObjectMeta{Name: "pod-1"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "pod-2"}},
			}

			require.Len(t, opts, 2)
			var listOpt client.ListOptions
			for _, opt := range opts {
				opt.ApplyToList(&listOpt)
			}
			require.Equal(t, "default", listOpt.Namespace)
			require.Zero(t, listOpt.Limit)
			require.Equal(t, ".metadata.controller=cosmoshub", listOpt.FieldSelector.String())
			return nil
		}

		control := NewFullNodeControl(nopStatusSyncer, reader)

		err := control.ConfirmPodRestoration(ctx, &crd)
		require.NoError(t, err)
	})

	t.Run("happy path - no items", func(t *testing.T) {
		control := NewFullNodeControl(nopStatusSyncer, nopReader)
		err := control.ConfirmPodRestoration(ctx, &crd)

		require.Error(t, err)
		require.EqualError(t, err, NO_CANDIDATES_ERR)
	})

}
