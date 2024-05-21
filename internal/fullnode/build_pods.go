package fullnode

import (
	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	"github.com/bharvest-devops/cosmos-operator/internal/diff"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	configChecksumAnnotation = "cosmos.bharvest/config-checksum"
)

// BuildPods creates the final state of pods given the crd.
func BuildPods(crd *cosmosv1.CosmosFullNode, cksums ConfigChecksums) ([]diff.Resource[*corev1.Pod], error) {
	var (
		builder   = NewPodBuilder(crd)
		overrides = crd.Spec.InstanceOverrides
		pods      []diff.Resource[*corev1.Pod]
	)
	candidates := podSnapshotCandidates(crd)
	for i := int32(0); i < crd.Spec.Replicas; i++ {
		pod, err := builder.WithOrdinal(i).Build()
		if err != nil {
			return nil, err
		}
		if _, shouldSnapshot := candidates[pod.Name]; shouldSnapshot {
			continue
		}

		// If current pod's pvc should be pruned, it'll automatically change current pod into pruningPod.
		if prunerPod := podPruner(crd, pod).BuildPruningContainer(crd); prunerPod != nil {
			pods = append(pods, diff.Adapt(prunerPod, i))
			continue
		}

		// Check if current pod's pvc should be re-generate, and it's true, current pod will be temporary deleted.
		if crd.Status.SelfHealing.RegenPVCStatus != nil {
			var isRegenerated bool
			for _, c := range crd.Status.SelfHealing.RegenPVCStatus.Candidates {
				if c.PodName == pod.Name && c.Namespace == pod.Namespace {
					isRegenerated = true
					break
				}
			}
			if isRegenerated {
				continue
			}
		}

		if len(crd.Spec.ChainSpec.Versions) > 0 {
			instanceHeight := uint64(0)
			if height, ok := crd.Status.Height[pod.Name]; ok {
				instanceHeight = height
			}
			var image string
			for _, version := range crd.Spec.ChainSpec.Versions {
				if instanceHeight < version.UpgradeHeight {
					break
				}
				image = version.Image
			}
			if image != "" {
				setChainContainerImage(pod, image)
			}
		}
		if o, ok := overrides[pod.Name]; ok {
			if o.DisableStrategy != nil {
				continue
			}
			if o.Image != "" {
				setChainContainerImage(pod, o.Image)
			}
		}

		if crd.Spec.PodTemplate.TerminationPolicy == cosmosv1.RemainTerminationPolicy {
			for j, c := range pod.Spec.Containers {
				c.Args = append(c.Args, ";trap : TERM INT; sleep infinity & wait")
				pod.Spec.Containers[j].Args = c.Args
			}
			for j, c := range pod.Spec.InitContainers {
				c.Args = append(c.Args, ";trap : TERM INT; sleep infinity & wait")
				pod.Spec.InitContainers[j].Args = c.Args
			}
		}

		pod.Annotations[configChecksumAnnotation] = cksums[client.ObjectKeyFromObject(pod)]
		pods = append(pods, diff.Adapt(pod, i))
	}
	return pods, nil
}

func setChainContainerImage(pod *corev1.Pod, image string) {
	for i := range pod.Spec.Containers {
		if pod.Spec.Containers[i].Name == mainContainer {
			pod.Spec.Containers[i].Image = image
			break
		}
	}
	for i := range pod.Spec.InitContainers {
		if pod.Spec.InitContainers[i].Name == chainInitContainer {
			pod.Spec.InitContainers[i].Image = image
			break
		}
	}
}

func podSnapshotCandidates(crd *cosmosv1.CosmosFullNode) map[string]struct{} {
	candidates := make(map[string]struct{})
	for _, v := range crd.Status.ScheduledSnapshotStatus {
		candidates[v.PodCandidate] = struct{}{}
	}
	return candidates
}

func podPruner(crd *cosmosv1.CosmosFullNode, pod *corev1.Pod) *PrunerPod {
	pruningStatus := crd.Status.SelfHealing.CosmosPruningStatus
	if pruningStatus == nil {
		return nil
	}
	for _, p := range pruningStatus.Candidates {
		if pod.Name == p.PodName && pod.Namespace == p.Namespace {
			prunerPod := PrunerPod(*pod)
			return ptr(prunerPod)
		}
	}
	return nil
}

type PrunerPod corev1.Pod
