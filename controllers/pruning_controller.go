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
	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	"github.com/bharvest-devops/cosmos-operator/internal/healthcheck"
	"github.com/bharvest-devops/cosmos-operator/internal/prune"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/retry"
	"net/http"
	"time"

	"github.com/bharvest-devops/cosmos-operator/internal/cosmos"
	"github.com/bharvest-devops/cosmos-operator/internal/fullnode"
	"github.com/bharvest-devops/cosmos-operator/internal/kube"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PruningReconciler reconciles a pruning partion of CosmosFullNode object
type PruningReconciler struct {
	client.Client
	diskClient      *fullnode.DiskUsageCollector
	recorder        record.EventRecorder
	pruner          *prune.Pruner
	fullNodeControl *prune.FullNodeControl
}

func NewPruningReconciler(
	client client.Client,
	recorder record.EventRecorder,
	statusClient *fullnode.StatusClient,
	cacheController *cosmos.CacheController,
	httpClient *http.Client,
) *PruningReconciler {
	return &PruningReconciler{
		Client:          client,
		diskClient:      fullnode.NewDiskUsageCollector(healthcheck.NewClient(httpClient), client),
		recorder:        recorder,
		pruner:          prune.NewPruner(cacheController),
		fullNodeControl: prune.NewFullNodeControl(statusClient, client),
	}
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PruningReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithName(cosmosv1.PruningController)
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

	defer func(Client client.Client, ctx context.Context, obj client.Object, opts ...client.UpdateOption) {
		_ = retry.RetryOnConflict(retry.DefaultRetry, func() error {
			return Client.Update(ctx, obj)
		})
	}(r.Client, ctx, crd)

	reporter := kube.NewEventReporter(logger, r.recorder, crd)

	retryResult := ctrl.Result{RequeueAfter: 10 * time.Second}

	if crd.Status.SelfHealing.CosmosPruningStatus == nil {
		pruningStatus := cosmosv1.CosmosPruningStatus{
			CosmosPruningPhase: cosmosv1.CosmosPruningPhaseFindingCandidate,
		}
		crd.Status.SelfHealing.CosmosPruningStatus = ptr(pruningStatus)
	}
	status := crd.Status.SelfHealing.CosmosPruningStatus

	switch status.CosmosPruningPhase {
	case cosmosv1.CosmosPruningPhaseFindingCandidate:
		usage, err := r.diskClient.CollectDiskUsage(ctx, crd)
		if err != nil {
			reporter.Error(err, "Failed to collect pvc disk usage")
			// This error can be noisy so we record a generic error. Check logs for error details.
			reporter.RecordError("PVCPruning", errors.New("failed to collect pvc disk usage"))
			return retryResult, err
		}

		var candidatePod *corev1.Pod
		candidatePod = r.pruner.FindCandidate(ctx, crd, usage)
		if candidatePod == nil {
			return retryResult, nil
		}

		msg := fmt.Sprintf("Pruning candidate found: %s", candidatePod.Name)
		reporter.Info(msg)
		reporter.RecordInfo("PVCPruning", msg)

		err = r.fullNodeControl.SignalPodReplace(ctx, crd, []*corev1.Pod{candidatePod})
		if err != nil {
			return ctrl.Result{}, err
		}

		crd.Status.SelfHealing.CosmosPruningStatus.CosmosPruningPhase = cosmosv1.CosmosPruningPhaseWaitingForPodReplaced
	case cosmosv1.CosmosPruningPhaseWaitingForPodReplaced:

		if err := r.fullNodeControl.ConfirmPodReplaced(ctx, crd); err != nil {
			reporter.Error(err, "Failed to confirm pod deletion")
			reporter.RecordError("PVCPruning", errors.Join(errors.New("WaitingForPodDeletionError"), err))
			if err.Error() == prune.NO_CANDIDATES_ERR {
				crd.Status.SelfHealing.CosmosPruningStatus.CosmosPruningPhase = cosmosv1.CosmosPruningPhaseFindingCandidate
			}
			return retryResult, nil
		}
		crd.Status.SelfHealing.CosmosPruningStatus.CosmosPruningPhase = cosmosv1.CosmosPruningPhasePruning

	case cosmosv1.CosmosPruningPhaseWaitingForComplete:
		ready, err := r.fullNodeControl.CheckPruningComplete(ctx, crd)
		if err != nil {
			reporter.Error(err, "Failed to check pruning status")
			reporter.RecordError("PruningCheckingError", err)
			return retryResult, nil
		}
		if !ready {
			reporter.Info("Pruning is not complete")
			return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
		}

		reporter.Info("Pruning complete")
		crd.Status.SelfHealing.CosmosPruningStatus.CosmosPruningPhase = cosmosv1.CosmosPruningPhaseRestorePod

	case cosmosv1.CosmosPruningPhaseRestorePod:
		err := r.fullNodeControl.SignalPodRestoration(ctx, crd)
		if err != nil {
			reporter.Error(err, "Failed to restore pod")
			reporter.RecordError("PodRestorationError", err)
			return retryResult, err
		}

		crd.Status.SelfHealing.CosmosPruningStatus.CosmosPruningPhase = cosmosv1.CosmosPruningPhaseConfirmPodRestoration
	case cosmosv1.CosmosPruningPhaseConfirmPodRestoration:
		err := r.fullNodeControl.ConfirmPodRestoration(ctx, crd)
		if err != nil {
			reporter.Error(err, "Failed to confirm pod Restoration")
			reporter.RecordError("ConfirmPodRestorationErr", err)
			return retryResult, err
		}
		// Complete pruning
		crd.Status.SelfHealing.CosmosPruningStatus.CosmosPruningPhase = cosmosv1.CosmosPruningPhaseFindingCandidate
	}

	// Updating status in the defer above triggers a new reconcile loop.
	return stopResult, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PruningReconciler) SetupWithManager(_ context.Context, mgr ctrl.Manager) error {
	// We do not have to index Pods by CosmosFullNode because the CosmosFullNodeReconciler already does so.
	// If we repeat it here, the manager returns an error.
	return ctrl.NewControllerManagedBy(mgr).
		For(&cosmosv1.CosmosFullNode{}).
		Complete(r)
}
