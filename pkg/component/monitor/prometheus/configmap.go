package prometheus

import (
	"fmt"
	"github.com/falser101/pulsar-operator/api/v1alpha1"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeConfigMap(c *v1alpha1.Pulsar) *v1.ConfigMap {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeConfigMapName(c),
			Namespace: c.Namespace,
		},
		Data: map[string]string{
			"prometheus.yml": prometheusConfig,
			"rules.yml":      "    groups: null",
		},
	}
}

func MakeConfigMapName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-monitor-prometheus-configmap", c.GetName())
}

var prometheusConfig = `
    global:
      scrape_interval: 15s
    rule_files:
      - 'rules.yml'
    alerting:
      alertmanagers:
      - static_configs:
        - targets: ['test001-tonglinkq-alert-manager:9093']
        path_prefix: /
    scrape_configs:
    - job_name: 'prometheus'
      static_configs:
      - targets:
        - '127.0.0.1:9090'
      metrics_path: /metrics
    - job_name: 'kubernetes-pods'
      bearer_token_file: /pulsar/tokens/client/token
      kubernetes_sd_configs:
      - role: pod
        namespaces:
          names:
            - test
      relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
        target_label: __address__
      - action: labelmap
        regex: __meta_kubernetes_pod_label_(.+)
      - source_labels: [__meta_kubernetes_namespace]
        action: replace
        target_label: kubernetes_namespace
      - source_labels: [__meta_kubernetes_pod_label_component]
        action: replace
        target_label: job
      - source_labels: [__meta_kubernetes_pod_name]
        action: replace
        target_label: kubernetes_pod_name
      metric_relabel_configs:
`
