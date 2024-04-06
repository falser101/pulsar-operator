/*
Copyright 2024.

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
	"github.com/falser101/pulsar-operator/internal/component/metadata"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type reconcileFunc func(cluster *v1alpha1.PulsarCluster) error

// PulsarClusterReconciler reconciles a PulsarCluster object
type PulsarClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	log    logr.Logger
}

//+kubebuilder:rbac:groups=message.apache.com,resources=pulsarclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=message.apache.com,resources=pulsarclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=message.apache.com,resources=pulsarclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PulsarCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *PulsarClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.log.Info("[Start] Reconciling PulsarCluster")
	defer func() {
		r.log.Info("[End] Reconciling PulsarCluster")
	}()

	pulsar := &v1alpha1.PulsarCluster{}
	if err := r.Get(ctx, req.NamespacedName, pulsar); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	changed := pulsar.SpecSetDefault()
	if changed {
		r.log.Info("Setting spec default settings for pulsar-cluster")
		if err := r.Update(ctx, pulsar); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}
	changed = pulsar.StatusSetDefault()
	if changed {
		r.log.Info("Setting status default settings for pulsar-cluster")
		if err := r.Status().Update(context.TODO(), pulsar); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}
	for _, fun := range []reconcileFunc{
		r.reconcileZookeeper,
		r.reconcileBookie,
		r.reconcileAutoRecovery,
		r.reconcileBroker,
		r.reconcileProxy,
		r.reconcileManager,
		r.reconcileMonitor,
		r.reconcilePulsarCluster,
	} {
		if err := fun(pulsar); err != nil {
			return reconcile.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PulsarClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PulsarCluster{}).
		Complete(r)
}

// Reconcile pulsarCluster resource
func (r *PulsarClusterReconciler) reconcilePulsarCluster(c *v1alpha1.PulsarCluster) error {
	if err := r.reconcileInitPulsarClusterMetaData(c); err != nil {
		return err
	}

	if err := r.reconcileAuthentication(c); err != nil {
		return err
	}

	if err := r.reconcilePulsarClusterPhase(c); err != nil {
		return err
	}
	return nil
}

// Init pulsar metaData
func (r *PulsarClusterReconciler) reconcileInitPulsarClusterMetaData(c *v1alpha1.PulsarCluster) (err error) {
	job := metadata.MakeInitClusterMetaDataJob(c)

	jobCur := &batchv1.Job{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      job.Name,
		Namespace: job.Namespace,
	}, jobCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, job, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), job); err == nil {
			r.log.Info("Start init pulsar cluster metaData job",
				"Job.Namespace", job.Namespace,
				"Job.Name", job.Name)
		}

	} else if err == nil && jobCur.Status.Succeeded == 1 {
		// Init pulsar cluster success
		c.Status.Phase = v1alpha1.PulsarClusterLaunchingPhase
		if err = r.Status().Update(context.TODO(), c); err == nil {
			r.log.Info("Init pulsar cluster metaData success",
				"PulsarCluster.Namespace", c.Namespace,
				"PulsarCluster.Name", c.Name)
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcilePulsarClusterPhase(c *v1alpha1.PulsarCluster) (err error) {
	if c.Status.Phase == v1alpha1.PulsarClusterLaunchingPhase &&
		r.isZookeeperRunning(c) &&
		r.isBookieRunning(c) &&
		r.isBrokerRunning(c) {
		c.Status.Phase = v1alpha1.PulsarClusterRunningPhase
		if err = r.Status().Update(context.TODO(), c); err == nil {
			r.log.Info("start pulsar cluster success",
				"PulsarCluster.Namespace", c.Namespace,
				"PulsarCluster.Name", c.Name)
		}
	}
	return
}
