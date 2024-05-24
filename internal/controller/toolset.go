package controller

import (
	"context"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/internal/component/toolset"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *PulsarClusterReconciler) reconcileToolset(c *v1alpha1.PulsarCluster) (err error) {
	if c.Status.Phase != v1alpha1.PulsarClusterRunningPhase {
		return
	}
	for _, fun := range []reconcileFunc{
		r.reconcileToolsetConfigMap,
		r.reconcileToolsetDeployment,
		r.reconcileToolsetService,
	} {
		if err = fun(c); err != nil {
			r.log.Error(err, "Reconciling Pulsar Toolset Error", c)
			return
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileToolsetConfigMap(c *v1alpha1.PulsarCluster) (err error) {
	cmCreate := toolset.MakeConfigMap(c)
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
			r.log.Info("Create pulsar toolset config map success",
				"ConfigMap.Namespace", c.Namespace,
				"ConfigMap.Name", cmCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileToolsetDeployment(c *v1alpha1.PulsarCluster) (err error) {
	depCreate := toolset.MakeDeployment(c)
	depCur := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      depCreate.Name,
		Namespace: depCreate.Namespace,
	}, depCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, depCreate, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), depCreate); err == nil {
			r.log.Info("Create pulsar toolset deployment success",
				"Deployment.Namespace", c.Namespace,
				"Deployment.Name", depCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileToolsetService(c *v1alpha1.PulsarCluster) (err error) {
	svcCreate := toolset.MakeService(c)
	svcCur := &v1.Service{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      svcCreate.Name,
		Namespace: svcCreate.Namespace,
	}, svcCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, svcCreate, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), svcCreate); err == nil {
			r.log.Info("Create pulsar toolset service success",
				"Service.Namespace", c.Namespace,
				"Service.Name", svcCreate.GetName())
		}
	}
	return
}
