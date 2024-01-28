package manager

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	"github.com/falser101/pulsar-operator/pkg/component/bookie"
	"github.com/falser101/pulsar-operator/pkg/component/broker"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeConfigMap(c *v1alpha1.Pulsar) *v1.ConfigMap {
	data := make(map[string]string)
	if c.Spec.Manager.ConfigMap == nil {
		data[BackendEntrypointKey] = MakeBackendEntrypoint(c)
		data[EntrypointKey] = EntrypointValue
	}
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeConfigMapName(c),
			Namespace: c.Namespace,
		},
		Data: data,
	}
}

func MakeConfigMapName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-manager-configmap", c.GetName())
}

func MakeBackendEntrypoint(c *v1alpha1.Pulsar) string {
	if c.Spec.Authentication.Enabled {
		return fmt.Sprintf(BackendEntrypointValue, bookie.MakeServiceName(c), c.GetName(), broker.MakeServiceName(c)) + `--backend.jwt.token=$(cat /pulsar/tokens/pulsar_manager/token) \
		--jwt.broker.token.mode=PRIVATE \
		--jwt.broker.public.key=file:///pulsar/keys/token/public.key \
		--jwt.broker.private.key=file:///pulsar/keys/token/private.key \`
	}
	return fmt.Sprintf(BackendEntrypointValue, bookie.MakeServiceName(c), c.GetName(), broker.MakeServiceName(c))
}
