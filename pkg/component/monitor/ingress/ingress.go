package ingress

//
//import (
//	"fmt"
//	cache"github.com/falser101/pulsar-operator/api/v1alpha1"
//	"pulsar-operator/controllers/monitor/grafana"
//
//	"k8s.io/api/extensions/v1beta1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/apimachinery/pkg/util/intstr"
//)
//
//func MakeIngress(c *cachev1alpha1.PulsarCluster) *v1beta1.Ingress {
//	return &v1beta1.Ingress{
//		TypeMeta: metav1.TypeMeta{
//			Kind:       "Ingress",
//			APIVersion: "v1",
//		},
//		ObjectMeta: metav1.ObjectMeta{
//			Name:        MakeIngressName(c),
//			Namespace:   c.Namespace,
//			Labels:      cachev1alpha1.MakeComponentLabels(c, cachev1alpha1.MonitorComponent),
//			Annotations: c.Spec.Monitor.Ingress.Annotations,
//		},
//		Spec: makeIngressSpec(c),
//	}
//}
//
//func MakeIngressName(c *cachev1alpha1.PulsarCluster) string {
//	return fmt.Sprintf("%s-monitor-ingress", c.GetName())
//}
//
//func makeIngressSpec(c *cachev1alpha1.PulsarCluster) v1beta1.IngressSpec {
//	s := v1beta1.IngressSpec{
//		Rules: make([]v1beta1.IngressRule, 0),
//	}
//
//	if c.Spec.Monitor.Grafana.Host != "" {
//		s.Rules = append(s.Rules, makeGrafanaRule(c))
//	}
//
//	if c.Spec.Monitor.Prometheus.Host != "" {
//		s.Rules = append(s.Rules, makePrometheusRule(c))
//	}
//	return s
//}
//
//func makeGrafanaRule(c *cachev1alpha1.PulsarCluster) v1beta1.IngressRule {
//	r := v1beta1.IngressRule{
//		Host: c.Spec.Monitor.Grafana.Host,
//		IngressRuleValue: v1beta1.IngressRuleValue{
//			HTTP: &v1beta1.HTTPIngressRuleValue{
//				Paths: make([]v1beta1.HTTPIngressPath, 0),
//			},
//		},
//	}
//	path := v1beta1.HTTPIngressPath{
//		Path: "/",
//		Backend: v1beta1.IngressBackend{
//			ServiceName: grafana.MakeServiceName(c),
//			ServicePort: intstr.FromInt(cachev1alpha1.PulsarGrafanaServerPort),
//		},
//	}
//	r.IngressRuleValue.HTTP.Paths = append(r.IngressRuleValue.HTTP.Paths, path)
//	return r
//}
//
//func makePrometheusRule(c *cachev1alpha1.PulsarCluster) v1beta1.IngressRule {
//	r := v1beta1.IngressRule{
//		Host: c.Spec.Monitor.Prometheus.Host,
//		IngressRuleValue: v1beta1.IngressRuleValue{
//			HTTP: &v1beta1.HTTPIngressRuleValue{
//				Paths: make([]v1beta1.HTTPIngressPath, 0),
//			},
//		},
//	}
//	path := v1beta1.HTTPIngressPath{
//		Path: "/",
//		Backend: v1beta1.IngressBackend{
//			ServiceName: prometheus.MakeServiceName(c),
//			ServicePort: intstr.FromInt(cachev1alpha1.PulsarPrometheusServerPort),
//		},
//	}
//	r.IngressRuleValue.HTTP.Paths = append(r.IngressRuleValue.HTTP.Paths, path)
//	return r
//}
