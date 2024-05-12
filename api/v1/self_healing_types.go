package v1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SelfHealingController, PruningController are the canonical controller name.
const SelfHealingController = "SelfHealing"
const PruningController = "Pruning"

// SelfHealSpec is part of a CosmosFullNode but is managed by a separate controller, SelfHealingReconciler.
// This is an effort to reduce complexity in the CosmosFullNodeReconciler.
// The controller only modifies the CosmosFullNode's status subresource relying on the CosmosFullNodeReconciler
// to reconcile appropriately.
type SelfHealSpec struct {
	// Automatically increases PVC storage as they approach capacity.
	//
	// Your cluster must support and use the ExpandInUsePersistentVolumes feature gate. This allows volumes to
	// expand while a pod is attached to it, thus eliminating the need to restart pods.
	// If you cluster does not support ExpandInUsePersistentVolumes, you will need to manually restart pods after
	// resizing is complete.
	// +optional
	PVCAutoScale *PVCAutoScaleSpec `json:"pvcAutoScale"`

	// Take action when a pod's height falls behind the max height of all pods AND still reports itself as in-sync.
	//
	// +optional
	HeightDriftMitigation *HeightDriftMitigationSpec `json:"heightDriftMitigation"`

	// PruningSpec configures strategy of pruning.
	//
	// In node operating, the most important is reliable service.
	// but to achieve this, you should resize disks when the node's disk size almost fulled.
	// or you can prune node every interval.
	//
	// This configuration supports you to prune nodes without manual tasks, through job will be run automatically at the same time every day.
	//
	// If you configure this, it'll be run before autoScaling pvc.
	//
	// +optional
	PruningSpec *PruningSpec `json:"pruningSpec"`
}

type PVCAutoScaleSpec struct {
	// The percentage of used disk space required to trigger scaling.
	// Example, if set to 80, autoscaling will not trigger until used space reaches >=80% of capacity.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:MaxSize=100
	UsedSpacePercentage int32 `json:"usedSpacePercentage"`

	// How much to increase the PVC's capacity.
	// Either a percentage (e.g. 20%) or a resource storage quantity (e.g. 100Gi).
	//
	// If a percentage, the existing capacity increases by the percentage.
	// E.g. PVC of 100Gi capacity + IncreaseQuantity of 20% increases disk to 120Gi.
	//
	// If a storage quantity (e.g. 100Gi), increases by that amount.
	IncreaseQuantity string `json:"increaseQuantity"`

	// A resource storage quantity (e.g. 2000Gi).
	// When increasing PVC capacity reaches >= MaxSize, autoscaling ceases.
	// Safeguards against storage quotas and costs.
	// +optional
	MaxSize resource.Quantity `json:"maxSize"`
}

type HeightDriftMitigationSpec struct {
	// If pod's height falls behind the max height of all pods by this value or more AND the pod's RPC /status endpoint
	// reports itself as in-sync, the pod is deleted. The CosmosFullNodeController creates a new pod to replace it
	// Pod deletion respects the CosmosFullNode.Spec.RolloutStrategy and will not delete more pods than set
	// by the strategy to prevent downtime.
	// This workaround is necessary to mitigate a bug in the Cosmos SDK and/or CometBFT where pods report themselves as
	// in-sync even though they can lag thousands of blocks behind the chain tip and cannot catch up.
	// A "rebooted" pod /status reports itself correctly and allows it to catch up to chain tip.
	// +kubebuilder:validation:Minimum:=1
	ThresholdHeight uint32 `json:"thresholdHeight"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Schemaless
	// +optional
	MaxHeightRetentionTime metav1.Duration `json:"maxHeightRetentionTime"`

	// RegeneratePVC specifies if delete pvc according to pods' starting failure count.
	// In most cases, unhealthy status involves invalid volume contents not only pod(computing resource); like containing appHash record.
	//
	// Through this field, you could regenerate pvc using spec.volumeClaimTemplate when the starting failure count of pod has over threshold.
	// If not set, deletion of pvc will be no.
	// +optional
	RegeneratePVC *RegeneratePVCSpec `json:"regeneratePVC"`
}

type RegeneratePVCSpec struct {

	// FailedCountCollectionDuration specifies how long self-healing controller will sum failure counts.
	// You could set this like "1m", "10m".
	// +kubebuilder:validator:Schemaless
	FailedCountCollectionDuration metav1.Duration `json:"failedCountCollectionDuration"`

	// ThresholdCount determines when regeneration logic will be run.
	ThresholdCount uint32 `json:"thresholdCount"`
}

// PruningSpec specifies whether you are going to prune data when node exceed threshold.
// It's similar with PVCAutoScaling, but more efficient way to save disks.
// Meanwhile, it could cause some non-reliable service providing.
type PruningSpec struct {

	// The percentage of used disk space required to trigger pruning.
	// Example, if set to 80, autoscaling will not trigger until used space reaches >=80% of capacity.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:MaxSize=100
	UsedSpacePercentage int32 `json:"usedSpacePercentage"`

	// The image url of you'll use for pruning.
	// If not set, defaults to "ghcr.io/bharvest-devops/cosmos-pruner:latest"
	// command
	// +optional
	Image string `json:"image"`

	// The description of shell used in pruning.
	// pruning will be process with "/bin/sh" cmd for flexibility.
	// And also it'll be used as args at 2nd like below;
	// args:
	//    - "-c"
	//    - "cosmos-pruner prune /home/operator/cosmos/data/ -b=0 ....
	//
	// If not set, it'll be set as default options.
	//   -b=0 # every blocks in tendermint data (also tx_index if you set true --tx_index) will be pruned.
	//   -v=0 # every versions in cosmos app will be pruned.
	//   --tx_index=true # prune also tx_index
	//   --compact=true # compact golevelDB after pruning
	//   --cosmos-sdk=true # pruning also cosmos app data
	//
	// +optional
	PruningCommand string `json:"pruningCommand"`
}

type DBBackend string

// CosmosPruningStatus shows status of process for pruning.
type CosmosPruningStatus struct {

	// Candidates describes what pod is currently pruning.
	// +optional
	Candidates map[string]PruningCandidate `json:"candidate"`

	// +optional
	PodPruningStatus map[string]PodPruningStatus `json:"podPruningStatus"`

	// The phase of the pruning.
	CosmosPruningPhase CosmosPruningPhase `json:"cosmosPruningPhase"`
}

type PruningCandidate struct {
	PodName   string `json:"podName"`
	Namespace string `json:"namespace"`
}

type PodPruningStatus struct {
	// LastPruned shows when does pod pruned.
	// +optional
	LastPruneTime *metav1.Time `json:"lastPruned"`
}

type CosmosPruningPhase string

const (
	// CosmosPruningPhaseFindingCandidate indicates it is waiting for exceeding of threshold, and also finding candidates.
	CosmosPruningPhaseFindingCandidate CosmosPruningPhase = "FindingCandidate"

	// CosmosPruningPhaseWaitingForPodReplaced indicates controller is waiting for the fullNodeRef to delete the candidate pod.
	CosmosPruningPhaseWaitingForPodReplaced CosmosPruningPhase = "WaitingForPodDeletion"

	// CosmosPruningPhaseWaitingForComplete indicates controller is waiting for complete pruning.
	CosmosPruningPhaseWaitingForComplete CosmosPruningPhase = "WaitingForComplete"

	// CosmosPruningPhaseRestorePod signals the fullNodeRef it can recreate the temporarily deleted pod.
	CosmosPruningPhaseRestorePod CosmosPruningPhase = "RestoringPod"

	// CosmosPruningPhaseConfirmPodRestoration confirms if pod is restored on cosmosFullNode
	CosmosPruningPhaseConfirmPodRestoration CosmosPruningPhase = "ConfirmPodRestoration"
)

type SelfHealingStatus struct {
	// PVC auto-scaling status.
	// +mapType:=granular
	// +optional
	PVCAutoScale map[string]*PVCAutoScaleStatus `json:"pvcAutoScaler"`

	// Re-generating PVC status.
	// +mapType:=granular
	// +optional
	RegenPVCStatus map[string]*RegenPVCStatus `json:"regenPVCStatus"`

	// CosmosPruning status.
	// +optional
	CosmosPruningStatus *CosmosPruningStatus `json:"cosmosPruningStatus"`
}

type RegenPVCStatus struct {
	FailureTimes []string `json:"podStartingFailureTimes"`

	// The phase of the controller.
	Phase *RegenPVCPhase `json:"phase"`
}

type RegenPVCPhase string

const (
	RegenPVCPhaseRegeneratingPVC RegenPVCPhase = "RegeneratingPVC"
	RegenPVCPhaseNotYet          RegenPVCPhase = "NotYet"
)

type PVCAutoScaleStatus struct {
	// The PVC size requested by the SelfHealing controller.
	RequestedSize resource.Quantity `json:"requestedSize"`
	// The timestamp the SelfHealing controller requested a PVC increase.
	RequestedAt metav1.Time `json:"requestedAt"`
}
