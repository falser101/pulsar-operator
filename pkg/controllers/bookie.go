package controllers

import (
	"context"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/pkg/component/bookie"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *PulsarClusterReconciler) reconcileBookie(c *v1alpha1.Pulsar) error {
	// if c.Status.Phase == v1alpha1.PulsarClusterInitializingPhase {
	// 	return nil
	// }

	for _, fun := range []reconcileFunc{
		r.reconcileBookieConfigMap,
		r.reconcileBookieStatefulSet,
		r.reconcileBookieService,
	} {
		if err := fun(c); err != nil {
			r.log.Error(err, "Reconciling Pulsar Broker Error", c)
			return err
		}
	}
	return nil
}

func (r *PulsarClusterReconciler) reconcileBookieConfigMap(c *v1alpha1.Pulsar) (err error) {
	cmCreate := bookie.MakeConfigMap(c)

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
			r.log.Info("Create pulsar bookie config map success",
				"ConfigMap.Namespace", c.Namespace,
				"ConfigMap.Name", cmCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileBookieStatefulSet(c *v1alpha1.Pulsar) (err error) {

	ssCreate := bookie.MakeStatefulSet(c)

	ssCur := &appsv1.StatefulSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      ssCreate.Name,
		Namespace: ssCreate.Namespace,
	}, ssCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, ssCreate, r.scheme); err != nil {
			return err
		}

		if err = r.client.Create(context.TODO(), ssCreate); err == nil {
			r.log.Info("Create pulsar bookie statefulSet success",
				"StatefulSet.Namespace", c.Namespace,
				"StatefulSet.Name", ssCreate.GetName())
		}
	} else if err != nil {
		return err
	} else {
		if c.Spec.Bookie.Size != *ssCur.Spec.Replicas {
			old := *ssCur.Spec.Replicas
			ssCur.Spec.Replicas = &c.Spec.Bookie.Size
			if err = r.client.Update(context.TODO(), ssCur); err == nil {
				r.log.Info("Scale pulsar bookie statefulSet success",
					"OldSize", old,
					"NewSize", c.Spec.Bookie.Size)
			}
		}
	}

	r.log.Info("Bookie node num info",
		"Replicas", ssCur.Status.Replicas,
		"ReadyNum", ssCur.Status.ReadyReplicas,
		"CurrentNum", ssCur.Status.CurrentReplicas,
	)
	return
}

func (r *PulsarClusterReconciler) reconcileBookieService(c *v1alpha1.Pulsar) (err error) {
	sCreate := bookie.MakeService(c)

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
			r.log.Info("Create Pulsar Bookie Service Success",
				"Service.Namespace", c.Namespace,
				"Service.Name", sCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) isBookieRunning(c *v1alpha1.Pulsar) bool {
	ss := &appsv1.StatefulSet{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      bookie.MakeStatefulSetName(c),
		Namespace: c.Namespace,
	}, ss)
	if err == nil {
		return ss.Status.ReadyReplicas == c.Spec.Bookie.Size
	}
	return false
}
