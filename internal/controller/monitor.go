package controller

import (
	"context"
	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/internal/component/monitor/grafana"
	"github.com/falser101/pulsar-operator/internal/component/monitor/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *PulsarClusterReconciler) reconcileMonitor(c *v1alpha1.PulsarCluster) error {
	if !c.Spec.Monitor.Enable {
		return nil
	}

	for _, fun := range []reconcileFunc{
		r.reconcileMonitorPrometheus,
		r.reconcileMonitorGrafana,
	} {
		if err := fun(c); err != nil {
			r.log.Error(err, "Reconciling PulsarCluster Monitor Error", c)
			return err
		}
	}

	//if c.Spec.Monitor.Ingress.Enable &&
	//	(c.Spec.Monitor.Grafana.Host != "" ||
	//		c.Spec.Monitor.Prometheus.Host != "") {
	//	if err := r.reconcileMonitorIngress(c); err != nil {
	//		r.log.Error(err, "Reconciling PulsarCluster Monitor Ingress Error", c)
	//		return err
	//	}
	//}
	return nil
}

func (r *PulsarClusterReconciler) reconcileMonitorPrometheus(c *v1alpha1.PulsarCluster) error {
	for _, fun := range []reconcileFunc{
		r.reconcileMonitorPrometheusRBAC,
		r.reconcileMonitorPrometheusConfigMap,
		r.reconcileMonitorPrometheusStatefulSet,
		r.reconcileMonitorPrometheusService,
	} {
		if err := fun(c); err != nil {
			return err
		}
	}
	return nil
}

func (r *PulsarClusterReconciler) reconcileMonitorPrometheusRBAC(c *v1alpha1.PulsarCluster) (err error) {
	// cluster role
	crCreate := prometheus.MakeClusterRole(c)
	crCur := &rbacv1.ClusterRole{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name: crCreate.Name,
	}, crCur, &client.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		if err = r.Create(context.TODO(), crCreate); err == nil {
			r.log.Info("Create pulsar monitor prometheus cluster role success",
				"ClusterRole.Name", crCreate.GetName())
		}
	}

	// service account
	saCreate := prometheus.MakeServiceAccount(c)
	saCur := &v1.ServiceAccount{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      saCreate.Name,
		Namespace: c.GetNamespace(),
	}, saCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, saCreate, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), saCreate); err == nil {
			r.log.Info("Create pulsar monitor prometheus service account success",
				"ServiceAccount.Name", saCreate.GetName())
		}
	}

	// cluster role and service account binding
	rbCreate := prometheus.MakeClusterRoleBinding(c)
	rbCur := &rbacv1.ClusterRoleBinding{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name: rbCreate.Name,
	}, rbCur)
	if err != nil && errors.IsNotFound(err) {
		if err = r.Create(context.TODO(), rbCreate); err == nil {
			r.log.Info("Create pulsar monitor prometheus cluster role binding success",
				"ClusterRoleBinding.Name", rbCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileMonitorPrometheusConfigMap(c *v1alpha1.PulsarCluster) (err error) {
	cmCreate := prometheus.MakeConfigMap(c)

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
			r.log.Info("Create pulsar monitor prometheus config map success",
				"ConfigMap.Namespace", c.Namespace,
				"ConfigMap.Name", cmCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileMonitorPrometheusStatefulSet(c *v1alpha1.PulsarCluster) (err error) {
	sSetCreate := prometheus.MakeStatefulSet(c)

	sSetCur := &appsv1.StatefulSet{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      sSetCreate.Name,
		Namespace: sSetCreate.Namespace,
	}, sSetCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, sSetCreate, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), sSetCreate); err == nil {
			r.log.Info("Create pulsar monitor prometheus statefulset success",
				"StatefulSet.Namespace", c.Namespace,
				"StatefulSet.Name", sSetCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileMonitorPrometheusService(c *v1alpha1.PulsarCluster) (err error) {
	sCreate := prometheus.MakeService(c)

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
			r.log.Info("Create pulsar monitor prometheus service success",
				"Service.Namespace", c.Namespace,
				"Service.Name", sCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileMonitorGrafana(c *v1alpha1.PulsarCluster) error {
	for _, fun := range []reconcileFunc{
		r.reconcileMonitorGrafanaConfigMap,
		r.reconcileManagerGrafanaSecret,
		r.reconcileManagerGrafanaPVC,
		r.reconcileMonitorGrafanaDeployment,
		r.reconcileMonitorGrafanaService,
	} {
		if err := fun(c); err != nil {
			return err
		}
	}
	return nil
}

func (r *PulsarClusterReconciler) reconcileMonitorGrafanaConfigMap(c *v1alpha1.PulsarCluster) (err error) {
	cmCreate := grafana.MakeConfigMap(c)
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
			r.log.Info("Create pulsar monitor grafana configmap success",
				"ConfigMap.Namespace", c.Namespace,
				"ConfigMap.Name", cmCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileManagerGrafanaSecret(c *v1alpha1.PulsarCluster) (err error) {
	secCreate := grafana.MakeSecret(c)
	secCur := &v1.Secret{}
	if err = r.Get(context.TODO(), types.NamespacedName{
		Namespace: secCreate.Namespace,
		Name:      secCreate.Name,
	}, secCur); err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, secCreate, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), secCreate); err == nil {
			r.log.Info("Create pulsar monitor grafana secret success",
				"Secret.Namespace", c.Namespace,
				"Secret.Name", secCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileManagerGrafanaPVC(c *v1alpha1.PulsarCluster) (err error) {
	if c.Spec.Monitor.Grafana.StorageClassName == "" {
		return
	}
	pvcCreate := grafana.MakePVC(c)
	pvcCur := &v1.PersistentVolumeClaim{}
	if err = r.Get(context.TODO(), types.NamespacedName{
		Namespace: pvcCreate.Namespace,
		Name:      pvcCreate.Name,
	}, pvcCur); err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, pvcCreate, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), pvcCreate); err == nil {
			r.log.Info("Create pulsar monitor grafana pvc success",
				"PVC.Namespace", c.Namespace,
				"PVC.Name", pvcCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileMonitorGrafanaDeployment(c *v1alpha1.PulsarCluster) (err error) {
	dmCreate := grafana.MakeDeployment(c)

	dmCur := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      dmCreate.Name,
		Namespace: dmCreate.Namespace,
	}, dmCur)
	if err != nil && errors.IsNotFound(err) {
		if err = controllerutil.SetControllerReference(c, dmCreate, r.Scheme); err != nil {
			return err
		}

		if err = r.Create(context.TODO(), dmCreate); err == nil {
			r.log.Info("Create pulsar monitor grafana deployment success",
				"Deployment.Namespace", c.Namespace,
				"Deployment.Name", dmCreate.GetName())
		}
	}
	return
}

func (r *PulsarClusterReconciler) reconcileMonitorGrafanaService(c *v1alpha1.PulsarCluster) (err error) {
	sCreate := grafana.MakeService(c)

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
			r.log.Info("Create pulsar monitor grafana service success",
				"Service.Namespace", c.Namespace,
				"Service.Name", sCreate.GetName())
		}
	}
	return
}

//func (r *PulsarClusterReconciler) reconcileMonitorIngress(c *cachev1alpha1.PulsarCluster) (err error) {
//	inCreate := ingress.MakeIngress(c)
//
//	inCur := &v1beta1.Ingress{}
//	err = r.Get(context.TODO(), types.NamespacedName{
//		Name:      inCreate.Name,
//		Namespace: inCreate.Namespace,
//	}, inCur)
//	if err != nil && errors.IsNotFound(err) {
//		if err = controllerutil.SetControllerReference(c, inCreate, r.Scheme); err != nil {
//			return err
//		}
//
//		if err = r.Create(context.TODO(), inCreate); err == nil {
//			r.log.Info("Create pulsar monitor ingress success",
//				"Ingress.Namespace", c.Namespace,
//				"Ingress.Name", inCreate.GetName())
//		}
//	}
//	return
//}
