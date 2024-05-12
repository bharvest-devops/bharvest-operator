package prune

import (
	"context"
	"errors"
	"fmt"
	"strings"

	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	cosmosalpha "github.com/bharvest-devops/cosmos-operator/api/v1alpha1"
	"github.com/bharvest-devops/cosmos-operator/internal/kube"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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

func (control FullNodeControl) CheckPruningComplete(ctx context.Context, crd *cosmosv1.CosmosFullNode) (bool, error) {
	var pods corev1.PodList
	if err := control.client.List(ctx, &pods,
		client.InNamespace(crd.Namespace),
		client.MatchingFields{kube.ControllerOwnerField: crd.Name},
	); err != nil {
		return false, fmt.Errorf("list pods: %w", err)
	}

	for _, pruningCandidate := range crd.Status.SelfHealing.CosmosPruningStatus.Candidates {
		for _, p := range pods.Items {
			if p.Name == GetPrunerPodName(pruningCandidate.PodName) {
				for _, containerStatus := range p.Status.ContainerStatuses {
					if containerStatus.State.Terminated == nil {
						return false, nil
					}
				}
			}
		}
	}
	return true, nil
}

func (control FullNodeControl) SignalPodReplace(ctx context.Context, crd *cosmosv1.CosmosFullNode, pods []*corev1.Pod) error {
	var joinedErr error
	for _, candidate := range pods {
		key := control.sourceKey(candidate.Name, candidate.Namespace)
		objKey := client.ObjectKey{Name: crd.Name, Namespace: crd.Namespace}
		if err := control.statusClient.SyncUpdate(ctx, objKey, func(status *cosmosv1.FullNodeStatus) {
			if status.SelfHealing.CosmosPruningStatus.Candidates == nil {
				status.SelfHealing.CosmosPruningStatus.Candidates = make(map[string]cosmosv1.PruningCandidate)
			}
			status.SelfHealing.CosmosPruningStatus.Candidates[key] = cosmosv1.PruningCandidate{PodName: candidate.Name, Namespace: candidate.Namespace}
		}); err != nil {
			joinedErr = errors.Join(joinedErr, err)
		}
	}
	return joinedErr
}

// ConfirmPodReplaced ConfirmPodDeletion returns a nil error if the pod is replaced.
// Any non-nil error is transient, including if the pod has not been replaced yet.
// If CosmosPruning.Status.Candidates are no, reconciler will be misunderstand it's working good.
func (control FullNodeControl) ConfirmPodReplaced(ctx context.Context, crd *cosmosv1.CosmosFullNode) error {
	var pods corev1.PodList
	if err := control.client.List(ctx, &pods,
		client.InNamespace(crd.Namespace),
		client.MatchingFields{kube.ControllerOwnerField: crd.Name},
	); err != nil {
		return fmt.Errorf("list pods: %w", err)
	}
	for _, pod := range pods.Items {
		for _, c := range crd.Status.SelfHealing.CosmosPruningStatus.Candidates {
			if pod.Name == c.PodName {
				return fmt.Errorf("pod %s not replaced yet", pod.Name)
			}
		}
	}
	return nil
}

// SignalPodRestoration updates the LocalFullNodeRef's status to indicate it should recreate the pod candidate.
// Any error returned can be treated as transient and retried.
func (control FullNodeControl) SignalPodRestoration(ctx context.Context, crd *cosmosv1.CosmosFullNode) error {
	var joinedErr error
	for _, candidate := range crd.Status.SelfHealing.CosmosPruningStatus.Candidates {
		key := control.sourceKey(candidate.PodName, candidate.Namespace)
		objKey := client.ObjectKey{Name: crd.Name, Namespace: crd.Namespace}
		if err := control.statusClient.SyncUpdate(ctx, objKey, func(status *cosmosv1.FullNodeStatus) {
			delete(status.SelfHealing.CosmosPruningStatus.PodPruningStatus, key)
		}); err != nil {
			joinedErr = errors.Join(joinedErr, err)
		}
	}
	return joinedErr
}

// ConfirmPodRestoration verifies the pod has been restored.
func (control FullNodeControl) ConfirmPodRestoration(ctx context.Context, crd *cosmosv1.CosmosFullNode) error {
	var (
		fullnode cosmosv1.CosmosFullNode
		getKey   = client.ObjectKey{Name: crd.Name, Namespace: crd.Namespace}
	)

	if err := control.client.Get(ctx, getKey, &fullnode); err != nil {
		return fmt.Errorf("get CosmosFullNode: %w", err)
	}

	for _, candidate := range crd.Status.SelfHealing.CosmosPruningStatus.Candidates {
		if _, exists := fullnode.Status.ScheduledSnapshotStatus[control.sourceKey(candidate.PodName, candidate.Namespace)]; exists {
			return fmt.Errorf("pod %s not restored yet", candidate.PodName)
		}
	}

	return nil
}

func (control FullNodeControl) sourceKey(candidateName, namespace string) string {
	key := strings.Join([]string{namespace, candidateName, cosmosalpha.GroupVersion.Version, cosmosalpha.GroupVersion.Group}, ".")
	// Remove all slashes because key is used in JSONPatch where slash "/" is a reserved character.
	return strings.ReplaceAll(key, "/", "")
}

func GetPrunerPodName(podName string) string {
	return fmt.Sprintf("%s-%s", podName, "pruning")
}
