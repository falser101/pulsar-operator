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

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"pulsar-operator/pkg/api/v1alpha1"
	"pulsar-operator/pkg/component/bookie"
	"pulsar-operator/pkg/component/metadata"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
)

type reconcileFunc func(cluster *v1alpha1.PulsarCluster) error

func NewPulsarClusterReconciler(
	client client.Client,
	scheme *runtime.Scheme,
	record record.EventRecorder,
	log logr.Logger,
) *PulsarClusterReconciler {
	return &PulsarClusterReconciler{
		client:   client,
		scheme:   scheme,
		recorder: record,
		log:      log,
	}
}

// PulsarClusterReconciler reconciles a PulsarCluster object
type PulsarClusterReconciler struct {
	client   client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
	log      logr.Logger
}

//+kubebuilder:rbac:groups=cache.example.com,resources=pulsarclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.example.com,resources=pulsarclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.example.com,resources=pulsarclusters/finalizers,verbs=update
// +kubebuilder:resource:path=foos,scope=Cluster
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the PulsarCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *PulsarClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.log.Info("[Start] Reconciling PulsarCluster")
	defer func() {
		r.log.Info("[End] Reconciling PulsarCluster")
	}()

	pulsarCluster := &v1alpha1.PulsarCluster{}
	if err := r.client.Get(ctx, req.NamespacedName, pulsarCluster); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	changed := pulsarCluster.SpecSetDefault()
	if changed {
		r.log.Info("Setting spec default settings for pulsar-cluster")
		if err := r.client.Update(ctx, pulsarCluster); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}
	changed = pulsarCluster.StatusSetDefault()
	if changed {
		r.log.Info("Setting status default settings for pulsar-cluster")
		if err := r.client.Status().Update(context.TODO(), pulsarCluster); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}
	for _, fun := range []reconcileFunc{
		r.reconcileZookeeper,
		r.reconcileBookie,
		r.reconcileAutoRecovery,
		r.reconcileBroker,
		r.reconcileManager,
		r.reconcileMonitor,
		r.reconcilePulsarCluster,
	} {
		if err := fun(pulsarCluster); err != nil {
			return reconcile.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PulsarClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PulsarCluster{}).
		Owns(&appsv1.Deployment{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 2}).
		Complete(r)
}

// Reconcile pulsarCluster resource
func (r *PulsarClusterReconciler) reconcilePulsarCluster(c *v1alpha1.PulsarCluster) error {
	if err := r.reconcileInitPulsarClusterMetaData(c); err != nil {
		return err
	}

	if err := r.reconcilePulsarClusterPhase(c); err != nil {
		return err
	}
	return nil
}

// Init pulsar metaData
func (r *PulsarClusterReconciler) reconcileInitPulsarClusterMetaData(c *v1alpha1.PulsarCluster) (err error) {
	if c.Status.Phase == v1alpha1.PulsarClusterInitializingPhase && r.isZookeeperRunning(c) {
		job := metadata.MakeInitClusterMetaDataJob(c)

		jobCur := &batchv1.Job{}
		err = r.client.Get(context.TODO(), types.NamespacedName{
			Name:      job.Name,
			Namespace: job.Namespace,
		}, jobCur)
		if err != nil && errors.IsNotFound(err) {
			if err = controllerutil.SetControllerReference(c, job, r.scheme); err != nil {
				return err
			}

			if err = r.client.Create(context.TODO(), job); err == nil {
				r.log.Info("Start init pulsar cluster metaData job",
					"Job.Namespace", job.Namespace,
					"Job.Name", job.Name)
			}

		} else if err == nil && jobCur.Status.Succeeded == 1 {
			// Init pulsar cluster success
			c.Status.Phase = v1alpha1.PulsarClusterLaunchingPhase
			if err = r.client.Status().Update(context.TODO(), c); err == nil {
				r.log.Info("Init pulsar cluster metaData success",
					"PulsarCluster.Namespace", c.Namespace,
					"PulsarCluster.Name", c.Name)
			}
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
		if err = r.client.Status().Update(context.TODO(), c); err == nil {
			r.log.Info("start pulsar cluster success",
				"PulsarCluster.Namespace", c.Namespace,
				"PulsarCluster.Name", c.Name)
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileBookieJob(c *v1alpha1.PulsarCluster) (err error) {
	jobCreate := bookie.MakeJob(c)
	jobCur := &batchv1.Job{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      jobCreate.Name,
		Namespace: jobCreate.Namespace,
	}, jobCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, jobCreate, r.scheme); err != nil {
			return err
		}

		if err = r.client.Create(context.TODO(), jobCreate); err == nil {
			r.log.Info("Start init pulsar bookie job",
				"Job.Namespace", jobCreate.Namespace,
				"Job.Name", jobCreate.Name)
		}

	} else if err == nil && jobCur.Status.Succeeded == 1 {
		// Init pulsar cluster success
		if err = r.client.Status().Update(context.TODO(), c); err == nil {
			r.log.Info("Init pulsar bookie success",
				"Job.Namespace", c.Namespace,
				"Job.Name", c.Name)
		}
	}
	return
}
