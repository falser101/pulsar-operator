package controllers

import (
	"context"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/pkg/component/broker"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *PulsarClusterReconciler) reconcileBroker(c *v1alpha1.Pulsar) error {
	for _, fun := range []reconcileFunc{
		r.reconcileBrokerConfigMap,
		r.reconcileBrokerDeployment,
		r.reconcileBrokerService,
		r.reconcileAuthentication,
	} {
		if err := fun(c); err != nil {
			r.log.Error(err, "Reconciling Pulsar Broker Error", c)
			return err
		}
	}
	return nil
}

func (r *PulsarClusterReconciler) reconcileBrokerConfigMap(c *v1alpha1.Pulsar) (err error) {
	cmCreate := broker.MakeConfigMap(c)

	cmCur := &v1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      cmCreate.Name,
		Namespace: cmCreate.Namespace,
	}, cmCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, cmCreate, r.scheme); err != nil {
			return err
		}

		if err = r.client.Create(context.TODO(), cmCreate); err == nil {
			r.log.Info("Create pulsar broker config map success",
				"ConfigMap.Namespace", c.Namespace,
				"ConfigMap.Name", cmCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileBrokerDeployment(c *v1alpha1.Pulsar) (err error) {
	dmCreate := broker.MakeDeployment(c)

	dmCur := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      dmCreate.Name,
		Namespace: dmCreate.Namespace,
	}, dmCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, dmCreate, r.scheme); err != nil {
			return err
		}

		if err = r.client.Create(context.TODO(), dmCreate); err == nil {
			r.log.Info("Create pulsar broker deployment success",
				"Deployment.Namespace", c.Namespace,
				"Deployment.Name", dmCreate.GetName())
		}
	} else if err != nil {
		return err
	} else {
		if c.Spec.Broker.Replicas != *dmCur.Spec.Replicas {
			old := *dmCur.Spec.Replicas
			dmCur.Spec.Replicas = &c.Spec.Broker.Replicas
			if err = r.client.Update(context.TODO(), dmCur); err == nil {
				r.log.Info("Scale pulsar broker deployment success",
					"OldSize", old,
					"NewSize", c.Spec.Broker.Replicas)
			}
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileBrokerService(c *v1alpha1.Pulsar) (err error) {
	sCreate := broker.MakeService(c)

	sCur := &v1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      sCreate.Name,
		Namespace: sCreate.Namespace,
	}, sCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, sCreate, r.scheme); err != nil {
			return err
		}

		if err = r.client.Create(context.TODO(), sCreate); err == nil {
			r.log.Info("Create pulsar broker service success",
				"Service.Namespace", c.Namespace,
				"Service.Name", sCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) isBrokerRunning(c *v1alpha1.Pulsar) bool {
	dm := broker.MakeDeployment(c)

	dmCur := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      dm.Name,
		Namespace: dm.Namespace,
	}, dmCur)
	return err == nil && dmCur.Status.ReadyReplicas == c.Spec.Broker.Replicas
}
