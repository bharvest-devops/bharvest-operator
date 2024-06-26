package prune

import (
	"context"
	"errors"
	"fmt"
	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	"github.com/bharvest-devops/cosmos-operator/internal/fullnode"
	"github.com/bharvest-devops/cosmos-operator/internal/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"
)

const NO_CANDIDATES_ERR = "there are no candidates"

type StatusSyncer interface {
	SyncUpdate(ctx context.Context, key client.ObjectKey, update func(status *cosmosv1.FullNodeStatus)) error
}

// FullNodeControl manages a ScheduledVolumeSnapshot's spec.fullNodeRef.
type FullNodeControl struct {
	client       client.Reader
	statusClient StatusSyncer
}

func NewFullNodeControl(statusClient StatusSyncer, client client.Reader) *FullNodeControl {
	return &FullNodeControl{client: client, statusClient: statusClient}
}

func (control FullNodeControl) SignalPodReplace(ctx context.Context, crd *cosmosv1.CosmosFullNode, pods []*corev1.Pod) error {
	for _, candidate := range pods {
		key := control.sourceKey(candidate.Name, candidate.Namespace)
		objKey := client.ObjectKey{Name: crd.Name, Namespace: crd.Namespace}
		if err := control.statusClient.SyncUpdate(ctx, objKey, func(status *cosmosv1.FullNodeStatus) {
			if status.SelfHealing.CosmosPruningStatus == nil {
				status.SelfHealing.CosmosPruningStatus = new(cosmosv1.CosmosPruningStatus)
			}
			if status.SelfHealing.CosmosPruningStatus.Candidates == nil {
				status.SelfHealing.CosmosPruningStatus.Candidates = make(map[string]cosmosv1.SelfHealingCandidate)
			}
			status.SelfHealing.CosmosPruningStatus.Candidates[key] = cosmosv1.SelfHealingCandidate{PodName: candidate.Name, Namespace: candidate.Namespace}
			status.SelfHealing.CosmosPruningStatus.CosmosPruningPhase = cosmosv1.CosmosPruningPhaseWaitingForPodReplaced

			if status.SelfHealing.CosmosPruningStatus.PodPruningStatus == nil {
				status.SelfHealing.CosmosPruningStatus.PodPruningStatus = make(map[string]cosmosv1.PodPruningStatus)
			}
			status.SelfHealing.CosmosPruningStatus.PodPruningStatus[key] = cosmosv1.PodPruningStatus{LastPruneTime: ptr(metav1.NewTime(time.Now()))}
		}); err != nil {
			return err
		}
	}
	return nil
}

// ConfirmPodReplaced ConfirmPodDeletion returns a nil error if the pod is replaced.
// Any non-nil error is transient, including if the pod has not been replaced yet.
// If CosmosPruning.Status.Candidates are no, reconciler will be misunderstood it's working good.
func (control FullNodeControl) ConfirmPodReplaced(ctx context.Context, crd *cosmosv1.CosmosFullNode) error {
	var (
		pods      corev1.PodList
		joinedErr error
	)
	if err := control.client.List(ctx, &pods,
		client.InNamespace(crd.Namespace),
		client.MatchingFields{kube.ControllerOwnerField: crd.Name},
	); err != nil {
		return fmt.Errorf("list pods: %w", err)
	}

	var existsCandidate bool
	for _, pod := range pods.Items {
		for _, c := range crd.Status.SelfHealing.CosmosPruningStatus.Candidates {
			if pod.Name == c.PodName {
				return fmt.Errorf("pod %s not replaced yet", pod.Name)
			} else {
				existsCandidate = true
			}
		}
	}
	pruningPhase := cosmosv1.CosmosPruningPhaseWaitingForComplete

	if !existsCandidate {
		pruningPhase = cosmosv1.CosmosPruningPhaseFindingCandidate
		joinedErr = errors.New(NO_CANDIDATES_ERR)
	}
	return errors.Join(joinedErr, control.statusClient.SyncUpdate(ctx, client.ObjectKeyFromObject(crd), func(status *cosmosv1.FullNodeStatus) {
		status.SelfHealing.CosmosPruningStatus.CosmosPruningPhase = pruningPhase
	}))
}

func (control FullNodeControl) CheckPruningComplete(ctx context.Context, crd *cosmosv1.CosmosFullNode) (bool, error) {
	var (
		pods            corev1.PodList
		existsCandidate bool
		pruningStatus   = crd.Status.SelfHealing.CosmosPruningStatus
	)
	if pruningStatus == nil {
		return false, errors.New(NO_CANDIDATES_ERR)
	} else if pruningStatus.Candidates == nil {
		return false, errors.New(NO_CANDIDATES_ERR)
	}

	if err := control.client.List(ctx, &pods,
		client.InNamespace(crd.Namespace),
		client.MatchingFields{kube.ControllerOwnerField: crd.Name},
	); err != nil {
		return false, fmt.Errorf("list pods: %w", err)
	}

	var existsPod bool
	for _, pruningCandidate := range pruningStatus.Candidates {
		for _, p := range pods.Items {
			if p.Name == fullnode.GetPrunerPodName(pruningCandidate.PodName) {
				for _, status := range p.Status.ContainerStatuses {
					if status.State.Terminated != nil {
						return true, control.statusClient.SyncUpdate(ctx, client.ObjectKeyFromObject(crd), func(status *cosmosv1.FullNodeStatus) {
							status.SelfHealing.CosmosPruningStatus.CosmosPruningPhase = cosmosv1.CosmosPruningPhaseRestorePod
						})
					}
				}
				existsCandidate = true
			}
			if strings.Contains(p.Name, pruningCandidate.PodName) {
				existsPod = true
			}
		}
	}

	if !existsPod {
		return true, control.statusClient.SyncUpdate(ctx, client.ObjectKeyFromObject(crd), func(status *cosmosv1.FullNodeStatus) {
			status.SelfHealing.CosmosPruningStatus.CosmosPruningPhase = cosmosv1.CosmosPruningPhaseRestorePod
		})
	}

	if !existsCandidate {
		return false, errors.New(NO_CANDIDATES_ERR)
	}
	return false, nil
}

// SignalPodRestoration updates the LocalFullNodeRef's status to indicate it should recreate the pod candidate.
// Any error returned can be treated as transient and retried.
func (control FullNodeControl) SignalPodRestoration(ctx context.Context, crd *cosmosv1.CosmosFullNode) error {
	var pruningStatus = crd.Status.SelfHealing.CosmosPruningStatus

	if pruningStatus == nil {
		return errors.New(NO_CANDIDATES_ERR)
	}

	for _, candidate := range pruningStatus.Candidates {
		key := control.sourceKey(candidate.PodName, candidate.Namespace)
		objKey := client.ObjectKey{Name: crd.Name, Namespace: crd.Namespace}
		return control.statusClient.SyncUpdate(ctx, objKey, func(status *cosmosv1.FullNodeStatus) {
			status.SelfHealing.CosmosPruningStatus.CosmosPruningPhase = cosmosv1.CosmosPruningPhaseConfirmPodRestoration
			delete(status.SelfHealing.CosmosPruningStatus.Candidates, key)
		})
	}
	return errors.New(NO_CANDIDATES_ERR)
}

// ConfirmPodRestoration verifies the pod has been restored.
func (control FullNodeControl) ConfirmPodRestoration(ctx context.Context, crd *cosmosv1.CosmosFullNode) error {
	var (
		existsCandidate bool
		pods            corev1.PodList
	)

	if err := control.client.List(ctx, &pods,
		client.InNamespace(crd.Namespace),
		client.MatchingFields{kube.ControllerOwnerField: crd.Name}); err != nil {
		return fmt.Errorf("list pods: %w", err)
	}

	for _, pod := range pods.Items {
		if _, exists := crd.Status.SelfHealing.CosmosPruningStatus.Candidates[control.sourceKey(pod.Name, pod.Namespace)]; exists {
			return fmt.Errorf("pod %s not restored yet", pod.Name)
		}
		existsCandidate = true
	}
	if !existsCandidate {
		return errors.New(NO_CANDIDATES_ERR)
	}

	return control.statusClient.SyncUpdate(ctx, client.ObjectKeyFromObject(crd), func(status *cosmosv1.FullNodeStatus) {
		status.SelfHealing.CosmosPruningStatus.CosmosPruningPhase = cosmosv1.CosmosPruningPhaseFindingCandidate
	})
}

func (control FullNodeControl) sourceKey(candidatePodName, namespace string) string {
	key := strings.Join([]string{namespace, candidatePodName, cosmosv1.GroupVersion.Version, cosmosv1.GroupVersion.Group}, ".")
	// Remove all slashes because key is used in JSONPatch where slash "/" is a reserved character.
	return strings.ReplaceAll(key, "/", "")
}
