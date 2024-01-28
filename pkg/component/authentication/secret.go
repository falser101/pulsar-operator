package authentication

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeSecret(c *v1alpha1.Pulsar) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeSecretName(c),
			Namespace: c.Namespace,
		},
	}
}

func MakeBrokerSecret(c *v1alpha1.Pulsar) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeBrokerSecretName(c),
			Namespace: c.Namespace,
		},
	}
}

func MakeBrokerSecretName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-token-broker-admin", c.Name)
}

func MakeSecretName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-token-asymmetric-key", c.Name)
}

func GenerateAsymmetricKey(c *v1alpha1.Pulsar) (private []byte, public []byte, err error) {
	cmd := exec.Command("bash", "-c", "bin/pulsarctl-amd64-linux/pulsarctl token create-key-pair -a RS256 --output-private-key private.key --output-public-key public.key")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	fmt.Println("exec result:", string(output))
	private, err = os.ReadFile("private.key")
	if err != nil {
		return
	}
	if err = os.Remove("private.key"); err != nil {
		return
	}
	public, err = os.ReadFile("public.key")
	if err != nil {
		return
	}
	if err = os.RemoveAll("public.key"); err != nil {
		return
	}
	return
}

func GenerateTokenKey(c *v1alpha1.Pulsar) (token []byte, err error) {
	cmd := exec.Command("bash", "-c", "bin/pulsarctl-amd64-linux/pulsarctl token create -a RS256 --private-key-file private.key --subject broker-admin 2&> broker-admin")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	fmt.Println("exec result:", string(output))
	token, err = os.ReadFile("broker-admin")
	if err != nil {
		return
	}
	if err = os.Remove("broker-admin"); err != nil {
		return
	}
	return
}
