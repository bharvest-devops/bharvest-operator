package fullnode

import (
	"context"
	"errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand"
	"strconv"
	"testing"
	"time"

	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type mockStatusSyncer func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error

func (fn mockStatusSyncer) SyncUpdate(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
	if ctx == nil {
		panic("nil context")
	}
	return fn(ctx, key, update)
}

func TestPVCHealer_SignalPVCResize(t *testing.T) {
	t.Parallel()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	ctx := context.Background()

	panicSyncer := mockStatusSyncer(func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
		panic("should not be called")
	})

	t.Run("happy path", func(t *testing.T) {
		var (
			capacity  = resource.MustParse("100Gi")
			stubNow   = time.Now()
			zeroQuant resource.Quantity
		)
		const (
			usedSpacePercentage = 80
			name                = "auto-scale-test"
			namespace           = "strangelove"
		)

		for _, tt := range []struct {
			Increase string
			Max      resource.Quantity
			Want     resource.Quantity
		}{
			{"20Gi", resource.MustParse("500Gi"), resource.MustParse("120Gi")},
			{"10%", zeroQuant, resource.MustParse("110Gi")},
			{"0.5Gi", zeroQuant, resource.MustParse("100.5Gi")},
			{"200%", zeroQuant, resource.MustParse("300Gi")},
			// Weird user input cases
			{"1", zeroQuant, *resource.NewQuantity(capacity.Value()+1, resource.BinarySI)},
		} {
			var crd cosmosv1.CosmosFullNode
			crd.APIVersion = "v1"
			crd.Name = name
			crd.Namespace = namespace
			crd.Spec.SelfHeal = &cosmosv1.SelfHealSpec{
				PVCAutoScale: &cosmosv1.PVCAutoScaleSpec{
					UsedSpacePercentage: usedSpacePercentage,
					IncreaseQuantity:    tt.Increase,
					MaxSize:             tt.Max,
				},
			}

			var patchCalled bool
			syncer := mockStatusSyncer(func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
				require.Equal(t, name, key.Name)
				require.Equal(t, namespace, key.Namespace)

				var got cosmosv1.FullNodeStatus
				update(&got)
				gotStatus := got.SelfHealing.PVCAutoScale
				require.Equal(t, stubNow, gotStatus["pvc-"+name+"-0"].RequestedAt.Time, tt)
				require.Truef(t, tt.Want.Equal(gotStatus["pvc-"+name+"-0"].RequestedSize), "%s:\nwant %+v\ngot  %+v", tt, tt.Want, gotStatus["pvc-"+name+"-0"].RequestedSize)

				patchCalled = true
				return nil
			})

			scaler := NewPVCHealer(syncer)
			scaler.now = func() time.Time {
				return stubNow
			}

			trigger := 80 + r.Intn(20)
			usage := []PVCDiskUsage{
				{Name: "pvc-" + name + "-0", PercentUsed: trigger, Capacity: capacity},
				{Name: "pvc-" + name + "-1", PercentUsed: 10},
				{Name: "pvc-" + name + "-2", PercentUsed: 79},
			}
			got, err := scaler.SignalPVCResize(ctx, &crd, lo.Shuffle(usage))

			require.NoError(t, err, tt)
			require.True(t, got, tt)
			require.True(t, patchCalled, tt)
		}
	})

	t.Run("does not exceed max", func(t *testing.T) {
		var (
			capacity = resource.MustParse("100Ti")
			maxSize  = resource.MustParse("200Ti")
		)
		const usedSpacePercentage = 80

		var crd cosmosv1.CosmosFullNode
		name := "name"
		crd.Name = name
		crd.Spec.SelfHeal = &cosmosv1.SelfHealSpec{
			PVCAutoScale: &cosmosv1.PVCAutoScaleSpec{
				UsedSpacePercentage: usedSpacePercentage,
				IncreaseQuantity:    "300%",
				MaxSize:             maxSize,
			},
		}

		var patchCalled bool
		syncer := mockStatusSyncer(func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
			var got cosmosv1.FullNodeStatus
			update(&got)
			gotStatus := got.SelfHealing.PVCAutoScale
			require.Equal(t, maxSize.Value(), gotStatus["pvc-"+name+"-0"].RequestedSize.Value())
			require.Equal(t, maxSize.Format, gotStatus["pvc-"+name+"-0"].RequestedSize.Format)

			patchCalled = true
			return nil
		})
		scaler := NewPVCHealer(syncer)

		usage := []PVCDiskUsage{
			{Name: "pvc-" + name + "-0", PercentUsed: 80, Capacity: capacity},
		}
		got, err := scaler.SignalPVCResize(ctx, &crd, lo.Shuffle(usage))

		require.NoError(t, err)
		require.True(t, got)
		require.True(t, patchCalled)
	})

	t.Run("capacity at or above max", func(t *testing.T) {
		for _, tt := range []struct {
			Max, Capacity resource.Quantity
		}{
			{resource.MustParse("5Ti"), resource.MustParse("5Ti")}, // the same
			{resource.MustParse("1G"), resource.MustParse("2G")},   // greater
		} {
			const usedSpacePercentage = 60

			var crd cosmosv1.CosmosFullNode
			name := "name"
			crd.Name = name
			crd.Spec.SelfHeal = &cosmosv1.SelfHealSpec{
				PVCAutoScale: &cosmosv1.PVCAutoScaleSpec{
					UsedSpacePercentage: usedSpacePercentage,
					IncreaseQuantity:    "10Gi",
					MaxSize:             tt.Max,
				},
			}

			scaler := NewPVCHealer(panicSyncer)
			usage := []PVCDiskUsage{
				{Name: "pvc-" + name + "-0", PercentUsed: 80, Capacity: tt.Capacity},
			}
			got, err := scaler.SignalPVCResize(ctx, &crd, usage)

			require.NoError(t, err, tt)
			require.False(t, got, tt)
		}
	})

	t.Run("no patch needed", func(t *testing.T) {
		for _, tt := range []struct {
			DiskUsage []PVCDiskUsage
		}{
			{nil}, // tests zero state
			{[]PVCDiskUsage{
				{PercentUsed: 79},
				{PercentUsed: 1},
				{PercentUsed: 10},
			}},
		} {
			var crd cosmosv1.CosmosFullNode
			crd.Spec.SelfHeal = &cosmosv1.SelfHealSpec{
				PVCAutoScale: &cosmosv1.PVCAutoScaleSpec{
					UsedSpacePercentage: 80,
					IncreaseQuantity:    "10Gi",
				},
			}

			scaler := NewPVCHealer(panicSyncer)
			got, err := scaler.SignalPVCResize(ctx, &crd, lo.Shuffle(tt.DiskUsage))

			require.NoError(t, err)
			require.False(t, got)
		}
	})

	t.Run("patch already signaled", func(t *testing.T) {
		const usedSpacePercentage = 90

		var crd cosmosv1.CosmosFullNode
		name := "name"
		crd.Name = name
		crd.Spec.SelfHeal = &cosmosv1.SelfHealSpec{
			PVCAutoScale: &cosmosv1.PVCAutoScaleSpec{
				UsedSpacePercentage: usedSpacePercentage,
				IncreaseQuantity:    "10Gi",
			},
		}
		crd.Status.SelfHealing.PVCAutoScale = map[string]*cosmosv1.PVCAutoScaleStatus{
			"pvc-" + name + "-0": {
				RequestedSize: resource.MustParse("100Gi"),
			},
		}

		scaler := NewPVCHealer(panicSyncer)
		usage := []PVCDiskUsage{
			{Name: "pvc-" + name + "-0", PercentUsed: usedSpacePercentage, Capacity: resource.MustParse("90Gi")},
		}
		got, err := scaler.SignalPVCResize(ctx, &crd, usage)

		require.NoError(t, err)
		require.False(t, got)
	})

	t.Run("invalid increase quantity", func(t *testing.T) {
		const usedSpacePercentage = 80

		for _, tt := range []struct {
			Increase string
		}{
			{""}, // CRD validation should prevent this
			{"wut"},
		} {
			var crd cosmosv1.CosmosFullNode
			name := "name"
			crd.Name = name
			crd.Spec.SelfHeal = &cosmosv1.SelfHealSpec{
				PVCAutoScale: &cosmosv1.PVCAutoScaleSpec{
					UsedSpacePercentage: usedSpacePercentage,
					IncreaseQuantity:    tt.Increase,
				},
			}

			scaler := NewPVCHealer(panicSyncer)
			usage := []PVCDiskUsage{
				{Name: "pvc-" + name + "-0", PercentUsed: usedSpacePercentage},
			}
			_, err := scaler.SignalPVCResize(ctx, &crd, lo.Shuffle(usage))

			require.Error(t, err)
			require.Contains(t, err.Error(), "invalid value for IntOrString: invalid type: string is not a percentage")
		}
	})

	t.Run("patch error", func(t *testing.T) {
		const usedSpacePercentage = 50

		var crd cosmosv1.CosmosFullNode
		crd.Spec.SelfHeal = &cosmosv1.SelfHealSpec{
			PVCAutoScale: &cosmosv1.PVCAutoScaleSpec{
				UsedSpacePercentage: usedSpacePercentage,
				IncreaseQuantity:    "10%",
			},
		}

		scaler := NewPVCHealer(mockStatusSyncer(func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
			return errors.New("boom")
		}))
		usage := []PVCDiskUsage{
			{Name: "pvc-0", PercentUsed: usedSpacePercentage},
		}
		_, err := scaler.SignalPVCResize(ctx, &crd, lo.Shuffle(usage))

		require.Error(t, err)
		require.EqualError(t, err, "boom")
	})
}

func TestPVCHealder_UpdatePodFailure(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("happy path", func(t *testing.T) {

		fullnodeStatus := cosmosv1.FullNodeStatus{
			SelfHealing: cosmosv1.SelfHealingStatus{},
		}

		mockSyncer := mockStatusSyncer(func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
			update(&fullnodeStatus)

			require.NotNil(t, fullnodeStatus.SelfHealing.RegenPVCStatus)
			require.Equal(t, 1, len(fullnodeStatus.SelfHealing.RegenPVCStatus.FailureTimes[sourceKey("initia-1", "default")]))
			require.NotNil(t, fullnodeStatus.SelfHealing.RegenPVCStatus.Candidates)

			return nil
		})

		crd := ptr(cosmosv1.CosmosFullNode{
			ObjectMeta: v1.ObjectMeta{
				Namespace: "default",
			},
			Spec: cosmosv1.FullNodeSpec{
				SelfHeal: ptr(cosmosv1.SelfHealSpec{
					HeightDriftMitigation: ptr(cosmosv1.HeightDriftMitigationSpec{
						RegeneratePVC: ptr(cosmosv1.RegeneratePVCSpec{
							ThresholdCount: 3,
						}),
					}),
				}),
			},
			Status: fullnodeStatus,
		})

		healer := NewPVCHealer(mockSyncer)

		isOver, err := healer.UpdatePodFailure(ctx, crd, "initia-1")

		require.NoError(t, err)
		require.Equal(t, false, isOver)
	})

	t.Run("elect candidate path", func(t *testing.T) {

		fullnodeStatus := cosmosv1.FullNodeStatus{
			SelfHealing: cosmosv1.SelfHealingStatus{},
		}

		podName := "initia-1"

		crd := ptr(cosmosv1.CosmosFullNode{
			ObjectMeta: v1.ObjectMeta{
				Namespace: "default",
			},
			Spec: cosmosv1.FullNodeSpec{
				SelfHeal: ptr(cosmosv1.SelfHealSpec{
					HeightDriftMitigation: ptr(cosmosv1.HeightDriftMitigationSpec{
						RegeneratePVC: ptr(cosmosv1.RegeneratePVCSpec{
							ThresholdCount:                3,
							FailedCountCollectionDuration: v1.Duration{Duration: 5 * time.Minute},
						}),
					}),
				}),
			},
			Status: fullnodeStatus,
		})

		mockSyncer := mockStatusSyncer(func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
			update(&crd.Status)
			return nil
		})

		healer := NewPVCHealer(mockSyncer)

		for i := 0; i < 2; i++ {
			isOver, err := healer.UpdatePodFailure(ctx, crd, podName)

			require.NoError(t, err)
			require.Equal(t, false, isOver)

			time.Sleep(1 * time.Second)
			t.Log(crd.Status.SelfHealing.RegenPVCStatus.FailureTimes[sourceKey(podName, "default")])
		}

		mockCheckSyncer := mockStatusSyncer(func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
			update(&crd.Status)

			require.NotNil(t, crd.Status.SelfHealing.RegenPVCStatus)
			require.Equal(t, podName, crd.Status.SelfHealing.RegenPVCStatus.Candidates[sourceKey(podName, "default")].PodName)
			require.Equal(t, 0, len(crd.Status.SelfHealing.RegenPVCStatus.FailureTimes[sourceKey(podName, "default")]))
			require.NotNil(t, crd.Status.SelfHealing.RegenPVCStatus.Candidates)

			return nil
		})

		candidateHealer := NewPVCHealer(mockCheckSyncer)
		isOver, err := candidateHealer.UpdatePodFailure(ctx, crd, podName)

		require.NoError(t, err)
		require.Equal(t, true, isOver)

	})

	t.Run("failure occurs for different pods path", func(t *testing.T) {

		fullnodeStatus := cosmosv1.FullNodeStatus{
			SelfHealing: cosmosv1.SelfHealingStatus{},
		}

		mockCheckSyncer := mockStatusSyncer(func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
			update(&fullnodeStatus)

			require.NotNil(t, fullnodeStatus.SelfHealing.RegenPVCStatus)
			require.Equal(t, 3, len(fullnodeStatus.SelfHealing.RegenPVCStatus.FailureTimes))
			require.NotNil(t, fullnodeStatus.SelfHealing.RegenPVCStatus.Candidates)

			require.Equal(t, 0, len(fullnodeStatus.SelfHealing.RegenPVCStatus.Candidates))

			return nil
		})

		crd := ptr(cosmosv1.CosmosFullNode{
			ObjectMeta: v1.ObjectMeta{
				Namespace: "default",
			},
			Spec: cosmosv1.FullNodeSpec{
				SelfHeal: ptr(cosmosv1.SelfHealSpec{
					HeightDriftMitigation: ptr(cosmosv1.HeightDriftMitigationSpec{
						RegeneratePVC: ptr(cosmosv1.RegeneratePVCSpec{
							ThresholdCount: 3,
						}),
					}),
				}),
			},
			Status: fullnodeStatus,
		})

		mockSyncer := mockStatusSyncer(func(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error {
			update(ptr(crd.Status))
			return nil
		})

		healer := NewPVCHealer(mockSyncer)

		i := 0
		for ; i < 2; i++ {
			isOver, err := healer.UpdatePodFailure(ctx, crd, "initia-"+strconv.Itoa(i))

			require.NoError(t, err)
			require.Equal(t, false, isOver)
		}
		checkHealer := NewPVCHealer(mockCheckSyncer)
		isOver, err := checkHealer.UpdatePodFailure(ctx, crd, "initia-"+strconv.Itoa(i))

		require.NoError(t, err)
		require.Equal(t, false, isOver)

	})

}
