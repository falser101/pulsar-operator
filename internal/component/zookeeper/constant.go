package zookeeper

import (
	"fmt"
	"github.com/falser101/pulsar-operator/api/v1alpha1"
)

// ServiceAnnotations Annotations
var ServiceAnnotations map[string]string
var StatefulSetAnnotations map[string]string

func init() {
	// Init Service Annotations
	ServiceAnnotations = make(map[string]string)
	ServiceAnnotations["service.alpha.kubernetes.io/tolerate-unready-endpoints"] = "true"

	// Init StatefulSet Annotations
	StatefulSetAnnotations = make(map[string]string)
	StatefulSetAnnotations["pod.alpha.kubernetes.io/initialized"] = "true"
	StatefulSetAnnotations["prometheus.io/scrape"] = "true"
	StatefulSetAnnotations["prometheus.io/port"] = "8000"
}

func MakeName(c *v1alpha1.PulsarCluster) string {
	return fmt.Sprintf("%s-%s", c.Name, v1alpha1.ZookeeperComponent)
}
