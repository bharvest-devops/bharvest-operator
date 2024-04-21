/*
Copyright 2024 B-Harvest Corporation.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"time"

	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	"github.com/bharvest-devops/cosmos-operator/internal/cosmos"
	"github.com/bharvest-devops/cosmos-operator/internal/fullnode"
	"github.com/bharvest-devops/cosmos-operator/internal/healthcheck"
	"github.com/bharvest-devops/cosmos-operator/internal/kube"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// SelfHealingReconciler reconciles the self healing portion of a CosmosFullNode object
type SelfHealingReconciler struct {
	client.Client
	cacheController *cosmos.CacheController
	diskClient      *fullnode.DiskUsageCollector
	driftDetector   fullnode.DriftDetection
	pvcHealer       *fullnode.PVCHealer
	recorder        record.EventRecorder
}

func NewSelfHealing(
	client client.Client,
	recorder record.EventRecorder,
	statusClient *fullnode.StatusClient,
	httpClient *http.Client,
	cacheController *cosmos.CacheController,
) *SelfHealingReconciler {
	return &SelfHealingReconciler{
		Client:          client,
		cacheController: cacheController,
		diskClient:      fullnode.NewDiskUsageCollector(healthcheck.NewClient(httpClient), client),
		driftDetector:   fullnode.NewDriftDetection(cacheController),
		pvcHealer:       fullnode.NewPVCHealer(statusClient),
		recorder:        recorder,
	}
}

// Reconcile reconciles only the self-healing spec in CosmosFullNode. If changes needed, this controller
// updates a CosmosFullNode status subresource thus triggering another reconcile loop. The CosmosFullNode
// uses the status object to reconcile its state.
func (r *SelfHealingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithName(cosmosv1.SelfHealingController)
	logger.V(1).Info("Entering reconcile loop", "request", req.NamespacedName)

	crd := new(cosmosv1.CosmosFullNode)
	if err := r.Get(ctx, req.NamespacedName, crd); err != nil {
		// Ignore not found errors because can't be fixed by an immediate requeue. We'll have to wait for next notification.
		// Also, will get "not found" error if crd is deleted.
		// No need to explicitly delete resources. Kube GC does so automatically because we set the controller reference
		// for each resource.
		return stopResult, client.IgnoreNotFound(err)
	}

	if crd.Spec.SelfHeal == nil {
		return stopResult, nil
	}

	reporter := kube.NewEventReporter(logger, r.recorder, crd)

	r.pvcAutoScale(ctx, reporter, crd)
	r.mitigateHeightDrift(ctx, reporter, crd)

	return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
}

func (r *SelfHealingReconciler) regeneratePVC(ctx context.Context, reporter kube.Reporter, crd *cosmosv1.CosmosFullNode, pod *corev1.Pod) {
	if crd.Spec.SelfHeal.HeightDriftMitigation.RegeneratePVC == nil {
		return
	}

	regenPVC, err := r.pvcHealer.UpdatePodFailure(ctx, crd, pod.Name)
	if err != nil {
		reporter.Error(err, "Failed to update podFailureStatus")
		reporter.RecordError("PVCRegenerating", err)
	}

	if regenPVC {
		pvc := new(corev1.PersistentVolumeClaim)

		// Find matching PVC to capture its actual capacity
		name := fullnode.PVCName(pod)
		key := client.ObjectKey{Namespace: pod.Namespace, Name: name}
		if err := r.Client.Get(ctx, key, pvc); err != nil {
			reporter.Error(err, "Failed to get pvc ", pod.Name)
			reporter.RecordError("PVCRegenerating", err)
			return
		}

		if err := r.Delete(ctx, pvc); err != nil {
			reporter.Error(err, "Failed to delete pvc", pvc.Name)
			reporter.RecordError("PVCRegenerating", err)
			return
		}

		msg := fmt.Sprintf("Pod %s has overed thresholdCount[%s]. Re-generating PVC...", pod.Name, crd.Spec.SelfHeal.HeightDriftMitigation.RegeneratePVC.FailedCountCollectionDuration.String())
		reporter.RecordInfo("PVCRegenerating", msg)
	}

	return
}

func (r *SelfHealingReconciler) pvcAutoScale(ctx context.Context, reporter kube.Reporter, crd *cosmosv1.CosmosFullNode) {
	if crd.Spec.SelfHeal.PVCAutoScale == nil {
		return
	}
	usage, err := r.diskClient.CollectDiskUsage(ctx, crd)
	if err != nil {
		reporter.Error(err, "Failed to collect pvc disk usage")
		// This error can be noisy so we record a generic error. Check logs for error details.
		reporter.RecordError("PVCAutoScaleCollectUsage", errors.New("failed to collect pvc disk usage"))
		return
	}
	didSignal, err := r.pvcHealer.SignalPVCResize(ctx, crd, usage)
	if err != nil {
		reporter.Error(err, "Failed to signal pvc resize")
		reporter.RecordError("PVCAutoScaleSignalResize", err)
		return
	}
	if !didSignal {
		return
	}
	const msg = "PVC auto scaling requested disk expansion"
	reporter.Info(msg)
	reporter.RecordInfo("PVCAutoScale", msg)
}

func (r *SelfHealingReconciler) mitigateHeightDrift(ctx context.Context, reporter kube.Reporter, crd *cosmosv1.CosmosFullNode) {
	if crd.Spec.SelfHeal.HeightDriftMitigation == nil {
		return
	}

	pods := r.driftDetector.LaggingPods(ctx, crd)
	var deleted int
	for _, pod := range pods {
		// CosmosFullNodeController will detect missing pod and re-create it.
		if err := r.Delete(ctx, pod); kube.IgnoreNotFound(err) != nil {
			reporter.Error(err, "Failed to delete pod", "pod", pod.Name)
			reporter.RecordError("HeightDriftMitigationDeletePod", err)
			continue
		}
		reporter.Info("Deleted pod for meeting height drift or heightRetainTime threshold", "pod", pod.Name)
		r.regeneratePVC(ctx, reporter, crd, pod)
		deleted++
	}
	if deleted > 0 {
		msg := fmt.Sprintf("Height lagged behind by more than %d blocks or overed heightRetainTime than (%d); deleted %d pod(s)",
			crd.Spec.SelfHeal.HeightDriftMitigation.ThresholdHeight, crd.Spec.SelfHeal.HeightDriftMitigation.ThresholdTime, deleted)
		reporter.RecordInfo("HeightDriftMitigation", msg)
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *SelfHealingReconciler) SetupWithManager(_ context.Context, mgr ctrl.Manager) error {
	// We do not have to index Pods because the CosmosFullNodeReconciler already does so.
	// If we repeat it here, the manager returns an error.
	return ctrl.NewControllerManagedBy(mgr).
		For(&cosmosv1.CosmosFullNode{}).
		Complete(r)
}
