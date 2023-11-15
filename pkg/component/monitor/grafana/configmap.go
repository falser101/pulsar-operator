package grafana

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
			"grafana.ini": grafanaConfig,
		},
	}
}

func MakeConfigMapName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-monitor-grafana-configmap", c.GetName())
}

var grafanaConfig = `[analytics]\ncheck_for_updates = true\n[auth.anonymous]\nenabled =
true\norg_name = Main Org.\norg_role = Admin\n[auth.azuread]\nallow_sign_up =
true\nallowed_domains = \nallowed_groups = \nauth_url = \nclient_id = \nclient_secret
= \nenabled = false\nname = Azure AD\nrole_attribute_strict = true\nscopes = openid
email profile\ntoken_url = \n[grafana_com]\nurl = https://grafana.com\n[log]\nmode
= console\n[log.file]\nformat = text\nlevel = info\n[paths]\ndata = /var/lib/grafana/pulsar/data\nplugins
= /var/lib/grafana/pulsar/plugins\nprovisioning = /var/lib/grafana/pulsar_provisioning\n[security]\nadmin_password
= {{ GRAFANA_ADMIN_PASSWORD }}\nadmin_user = {{ GRAFANA_ADMIN_USER }}\nallow_embedding
= true\n[server]\ndomain = {{ GRAFANA_DOMAIN }}\nroot_url = {{ GRAFANA_ROOT_URL
}}\nserve_from_sub_path = {{ GRAFANA_SERVE_FROM_SUB_PATH }}\n`
