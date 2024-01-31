package controllers

import (
	"context"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/pkg/component/authentication"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *PulsarClusterReconciler) reconcileAuthentication(c *v1alpha1.Pulsar) (err error) {
	if !c.Spec.Authentication.Enabled {
		return
	}
	for _, fun := range []reconcileFunc{
		r.reconcileServiceAccount,
		r.reconcileRole,
		r.reconcileRoleBinding,
		r.reconcileJob,
	} {
		if err = fun(c); err != nil {
			r.log.Error(err, "Reconciling Pulsar Authentication Error", c)
			return
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileServiceAccount(c *v1alpha1.Pulsar) (err error) {
	sa := authentication.MakeServiceAccount(c)
	saCur := &corev1.ServiceAccount{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      sa.Name,
		Namespace: sa.Namespace,
	}, saCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, sa, r.scheme); err != nil {
			return
		}

		if err = r.client.Create(context.TODO(), sa); err == nil {
			r.log.Info("Create pulsar secret sa success",
				"ServiceAccount.Namespace", sa.Namespace,
				"ServiceAccount.Name", sa.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileRole(c *v1alpha1.Pulsar) (err error) {
	role := authentication.MakeRole(c)
	roleCur := &rbacv1.Role{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      role.Name,
		Namespace: role.Namespace,
	}, roleCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, role, r.scheme); err != nil {
			return err
		}

		if err = r.client.Create(context.TODO(), role); err == nil {
			r.log.Info("Create pulsar secret role success",
				"Role.Namespace", role.Namespace,
				"Role.Name", role.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileRoleBinding(c *v1alpha1.Pulsar) (err error) {
	roleBinding := authentication.MakeRoleBinding(c)
	roleBindingCur := &rbacv1.RoleBinding{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: roleBinding.Namespace,
		Name:      roleBinding.Name,
	}, roleBindingCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, roleBinding, r.scheme); err != nil {
			return err
		}

		if err = r.client.Create(context.TODO(), roleBinding); err == nil {
			r.log.Info("Create pulsar secret RoleBinding success",
				"RoleBinding.Namespace", roleBinding.Namespace,
				"RoleBinding.Name", roleBinding.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileJob(c *v1alpha1.Pulsar) (err error) {
	job := authentication.MakeJob(c)
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
			r.log.Info("Create pulsar secret job success",
				"job.Namespace", job.Namespace,
				"job.Name", job.GetName())
		}
	}
	return
}
