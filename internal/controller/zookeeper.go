package controller

import (
	"context"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/internal/component/zookeeper"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	policyV1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *PulsarClusterReconciler) reconcileZookeeper(c *v1alpha1.PulsarCluster) error {
	for _, fun := range []reconcileFunc{
		r.reconcileZookeeperConfigMap,
		r.reconcileZookeeperRBAC,
		r.reconcileZookeeperStatefulSet,
		r.reconcileZookeeperService,
		r.reconcileZookeeperPodDisruptionBudget,
	} {
		if err := fun(c); err != nil {
			r.log.Error(err, "Reconciling Pulsar Zookeeper Error", c)
			return err
		}
	}
	return nil
}

func (r *PulsarClusterReconciler) reconcileZookeeperConfigMap(c *v1alpha1.PulsarCluster) (err error) {
	cmCreate := zookeeper.MakeConfigMap(c)

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
			r.log.Info("Create pulsar zookeeper config map success",
				"ConfigMap.Namespace", c.Namespace,
				"ConfigMap.Name", cmCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileZookeeperRBAC(c *v1alpha1.PulsarCluster) (err error) {
	if err = r.reconcileZookeeperServiceAccount(c); err != nil {
		return
	}
	if err = r.reconcileZookeeperRole(c); err != nil {
		return
	}
	if err = r.reconcileZookeeperRoleBinding(c); err != nil {
		return
	}
	return
}

func (r *PulsarClusterReconciler) reconcileZookeeperServiceAccount(c *v1alpha1.PulsarCluster) (err error) {
	saCreate := zookeeper.MakeServiceAccount(c)
	saCur := &v1.ServiceAccount{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      saCreate.Name,
		Namespace: saCreate.Namespace,
	}, saCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, saCreate, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), saCreate); err == nil {
			r.log.Info("Create pulsar zookeeper service account success",
				"ServiceAccount.Namespace", c.Namespace,
				"ServiceAccount.Name", saCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileZookeeperRole(c *v1alpha1.PulsarCluster) (err error) {
	roleCreate := zookeeper.MakeRole(c)
	roleCur := &rbacv1.Role{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      roleCreate.Name,
		Namespace: roleCreate.Namespace,
	}, roleCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, roleCreate, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), roleCreate); err == nil {
			r.log.Info("Create pulsar zookeeper role success",
				"Role.Namespace", c.Namespace,
				"Role.Name", roleCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileZookeeperRoleBinding(c *v1alpha1.PulsarCluster) (err error) {
	rbCreate := zookeeper.MakeRoleBinding(c)
	rbCur := &rbacv1.RoleBinding{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      rbCreate.Name,
		Namespace: rbCreate.Namespace,
	}, rbCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, rbCreate, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), rbCreate); err == nil {
			r.log.Info("Create pulsar zookeeper roleBinding success",
				"RoleBinding.Namespace", c.Namespace,
				"RoleBinding.Name", rbCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileZookeeperStatefulSet(c *v1alpha1.PulsarCluster) (err error) {
	ssCreate := zookeeper.MakeStatefulSet(c)

	ssCur := &appsv1.StatefulSet{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      ssCreate.Name,
		Namespace: ssCreate.Namespace,
	}, ssCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, ssCreate, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), ssCreate); err == nil {
			r.log.Info("Create pulsar zookeeper statefulSet success",
				"StatefulSet.Namespace", c.Namespace,
				"StatefulSet.Name", ssCreate.GetName())
		}
	} else if err != nil {
		return err
	} else {
		if c.Spec.Zookeeper.Replicas != *ssCur.Spec.Replicas {
			old := *ssCur.Spec.Replicas
			ssCur.Spec.Replicas = &c.Spec.Zookeeper.Replicas
			if err = r.Update(context.TODO(), ssCur); err == nil {
				r.log.Info("Scale pulsar zookeeper statefulSet success",
					"OldSize", old,
					"NewSize", c.Spec.Zookeeper.Replicas)
			}
		}
	}

	r.log.Info("Zookeeper node num info",
		"Replicas", ssCur.Status.Replicas,
		"ReadyNum", ssCur.Status.ReadyReplicas,
		"CurrentNum", ssCur.Status.CurrentReplicas,
	)
	return
}

func (r *PulsarClusterReconciler) reconcileZookeeperService(c *v1alpha1.PulsarCluster) (err error) {
	sCreate := zookeeper.MakeService(c)

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
			r.log.Info("Create pulsar zookeeper service success",
				"Service.Namespace", c.Namespace,
				"Service.Name", sCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileZookeeperPodDisruptionBudget(c *v1alpha1.PulsarCluster) (err error) {
	pdb := zookeeper.MakePodDisruptionBudget(c)

	pdbCur := &policyV1.PodDisruptionBudget{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      pdb.Name,
		Namespace: pdb.Namespace,
	}, pdbCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, pdb, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), pdb); err == nil {
			r.log.Info("Create pulsar zookeeper podDisruptionBudget success",
				"PodDisruptionBudget.Namespace", c.Namespace,
				"PodDisruptionBudget.Name", pdb.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) isZookeeperRunning(c *v1alpha1.PulsarCluster) bool {
	ss := &appsv1.StatefulSet{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      zookeeper.MakeStatefulSetName(c),
		Namespace: c.Namespace,
	}, ss)
	if err == nil {
		return ss.Status.ReadyReplicas == c.Spec.Zookeeper.Replicas
	}
	return false
}
