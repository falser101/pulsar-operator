package controller

import (
	"context"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/internal/component/proxy"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *PulsarClusterReconciler) reconcileProxy(c *v1alpha1.PulsarCluster) error {
	for _, fun := range []reconcileFunc{
		r.reconcileProxyConfigMap,
		r.reconcileProxyDeployment,
		r.reconcileProxyService,
	} {
		if err := fun(c); err != nil {
			r.log.Error(err, "Reconciling PulsarCluster Proxy Error", c)
			return err
		}
	}
	return nil
}

func (r *PulsarClusterReconciler) reconcileProxyDeployment(c *v1alpha1.PulsarCluster) (err error) {
	ssCreate := proxy.MakeDeployment(c)
	ssCur := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      ssCreate.Name,
		Namespace: ssCreate.Namespace,
	}, ssCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, ssCreate, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), ssCreate); err == nil {
			r.log.Info("Create pulsar proxy deployment success",
				"Deployment.Namespace", c.Namespace,
				"Deployment.Name", ssCreate.GetName())
		}
	} else if err != nil {
		return err
	} else {
		if c.Spec.Proxy.Replicas != *ssCur.Spec.Replicas {
			old := *ssCur.Spec.Replicas
			ssCur.Spec.Replicas = &c.Spec.Proxy.Replicas
			if err = r.Update(context.TODO(), ssCur); err == nil {
				r.log.Info("Scale pulsar Proxy deployment success",
					"OldSize", old,
					"NewSize", c.Spec.Proxy.Replicas)
			}
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileProxyService(c *v1alpha1.PulsarCluster) (err error) {
	sCreate := proxy.MakeService(c)
	sCur := &v1.Service{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      sCreate.Name,
		Namespace: sCreate.Namespace,
	}, sCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, sCreate, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), sCreate); err == nil {
			r.log.Info("Create Pulsar Proxy Service Success",
				"Service.Namespace", c.Namespace,
				"Service.Name", sCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileProxyConfigMap(c *v1alpha1.PulsarCluster) (err error) {
	cmCreate := proxy.MakeConfigMap(c)

	cmCur := &v1.ConfigMap{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      cmCreate.Name,
		Namespace: cmCreate.Namespace,
	}, cmCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, cmCreate, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), cmCreate); err == nil {
			r.log.Info("Create pulsar Proxy config map success",
				"ConfigMap.Namespace", c.Namespace,
				"ConfigMap.Name", cmCreate.GetName())
		}
	}
	return
}
