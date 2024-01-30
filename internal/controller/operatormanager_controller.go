/*
Copyright 2023.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komk8s "github.com/kkb0318/kom/internal/kubernetes"
	"github.com/kkb0318/kom/internal/tool/factory"
)

const komFinalizer = "kom.kkb.jp/finalizers"

// OperatorManagerReconciler reconciles a OperatorManager object
type OperatorManagerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=kom.kkb.jp,resources=operatormanagers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kom.kkb.jp,resources=operatormanagers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kom.kkb.jp,resources=operatormanagers/finalizers,verbs=get;create;update;patch;delete
// +kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=helmrepositories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=helm.toolkit.fluxcd.io,resources=helmreleases,verbs=get;list;watch;create;update;patch;delete

// SetupWithManager sets up the controller with the Manager.
func (r *OperatorManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&komv1alpha1.OperatorManager{}).
		Complete(r)
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the OperatorManager object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *OperatorManagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, retErr error) {
	log := ctrllog.FromContext(ctx)
	log.Info("start reconciling")
	obj := &komv1alpha1.OperatorManager{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Apply Resources to pull helm, oci, git
	if !controllerutil.ContainsFinalizer(obj, komFinalizer) {
		controllerutil.AddFinalizer(obj, komFinalizer)
		if err := r.Update(ctx, obj); err != nil {
			log.Error(err, "Failed to update custom resource to add finalizer")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}
	// Examine if the object is under deletion.
	if !obj.DeletionTimestamp.IsZero() {
		retErr = r.reconcileDelete(ctx, obj)
		return
	}
	if err := r.reconcile(ctx, obj); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("successfully reconciled")
	return ctrl.Result{}, nil
}

func (r *OperatorManagerReconciler) reconcile(ctx context.Context, obj *komv1alpha1.OperatorManager) error {
	log := ctrllog.FromContext(ctx)
	rm := factory.NewResourceManager(*obj)
	handler, err := komk8s.NewHandler(obj, r.Client, komk8s.Owner{Field: "kom"})
	if err != nil {
		return err
	}
	err = handler.ApplyAll(ctx, rm)
	if err != nil {
		log.Error(err, "server-side apply failed")
		return err
	}
	log.Info("server-side apply completed")
	return nil
}

func (r *OperatorManagerReconciler) reconcileDelete(ctx context.Context, obj *komv1alpha1.OperatorManager) error {
	// Remove our finalizer from the list
	log := ctrllog.FromContext(ctx)
	rm := factory.NewResourceManager(*obj)
	handler, err := komk8s.NewHandler(obj, r.Client, komk8s.Owner{Field: "kom"})
	if err != nil {
		return err
	}
	err = handler.DeleteAll(ctx, rm)
	if err != nil {
		log.Error(err, "deletion failed")
		return err
	}
	controllerutil.RemoveFinalizer(obj, komFinalizer)
	return r.Update(ctx, obj)
}
