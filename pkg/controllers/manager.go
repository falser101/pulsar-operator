package controllers

import (
	"context"
	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/pkg/component/manager"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *PulsarClusterReconciler) reconcileManager(c *v1alpha1.Pulsar) (err error) {
	for _, fun := range []reconcileFunc{
		r.reconcileManagerConfigMap,
		r.reconcileManagerStatefulSet,
		r.reconcileManagerService,
		r.reconcileManagerJob,
	} {
		if err = fun(c); err != nil {
			r.log.Error(err, "Reconciling Pulsar Manager Error", c)
			return
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileManagerConfigMap(c *v1alpha1.Pulsar) (err error) {
	cmCreate := manager.MakeConfigMap(c)
	configMap := &v1.ConfigMap{}
	if err = r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: cmCreate.Namespace,
		Name:      cmCreate.Name,
	}, configMap); err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, cmCreate, r.scheme); err != nil {
			return err
		}

		if err = r.client.Create(context.TODO(), cmCreate); err == nil {
			r.log.Info("Create pulsar manager config map success",
				"ConfigMap.Namespace", c.Namespace,
				"ConfigMap.Name", cmCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileManagerStatefulSet(c *v1alpha1.Pulsar) (err error) {
	sCreate := manager.MakeStatefulSet(c)
	sCur := &appsv1.StatefulSet{}
	if err = r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: sCreate.Namespace,
		Name:      sCreate.Name,
	}, sCur); err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, sCreate, r.scheme); err != nil {
			return err
		}

		if err = r.client.Create(context.TODO(), sCreate); err == nil {
			r.log.Info("Create pulsar manager statefulSet success",
				"StatefulSet.Namespace", c.Namespace,
				"StatefulSet.Name", sCreate.GetName())
		}
	}
	r.log.Info("Manager node num info",
		"Replicas", sCur.Status.Replicas,
		"ReadyNum", sCur.Status.ReadyReplicas,
		"CurrentNum", sCur.Status.CurrentReplicas,
	)
	return
}

func (r *PulsarClusterReconciler) reconcileManagerJob(c *v1alpha1.Pulsar) (err error) {
	jobCreate := manager.MakeJob(c)
	jobCur := &batchv1.Job{}
	if err = r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: jobCreate.Namespace,
		Name:      jobCreate.Name,
	}, jobCur); err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, jobCreate, r.scheme); err != nil {
			return err
		}

		if err = r.client.Create(context.TODO(), jobCreate); err == nil {
			r.log.Info("Create Pulsar Manager Job Success",
				"Job.Namespace", c.Namespace,
				"Job.Name", jobCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileManagerService(c *v1alpha1.Pulsar) (err error) {
	svcCreate := manager.MakeService(c)
	svcCur := &v1.Service{}
	if err = r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: svcCreate.Namespace,
		Name:      svcCreate.Name,
	}, svcCur); err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, svcCreate, r.scheme); err != nil {
			return err
		}

		if err = r.client.Create(context.TODO(), svcCreate); err == nil {
			r.log.Info("Create Pulsar Manager Service Success",
				"Service.Namespace", c.Namespace,
				"Service.Name", svcCreate.GetName())
		}
	}
	return
}
