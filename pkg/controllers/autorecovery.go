package controllers

import (
	"context"
	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/pkg/component/autorecovery"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *PulsarClusterReconciler) reconcileAutoRecovery(c *v1alpha1.Pulsar) error {
	if c.Status.Phase == v1alpha1.PulsarClusterInitializingPhase {
		return nil
	}

	for _, fun := range []reconcileFunc{
		r.reconcileBookieAutoRecoveryDeployment,
	} {
		if err := fun(c); err != nil {
			r.log.Error(err, "Reconciling Pulsar AutoRecovery Error", c)
			return err
		}
	}
	return nil
}

func (r *PulsarClusterReconciler) reconcileBookieAutoRecoveryDeployment(c *v1alpha1.Pulsar) (err error) {
	dmCreate := autorecovery.MakeDeployment(c)

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
			r.log.Info("Create pulsar autoRecovery deployment success",
				"Deployment.Namespace", c.Namespace,
				"Deployment.Name", dmCreate.GetName())
		}
	}

	r.log.Info("Pulsar autoRecovery node num info",
		"Replicas", dmCur.Status.Replicas,
		"ReadyNum", dmCur.Status.ReadyReplicas,
	)
	return
}
