package validation

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	extensioncs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"testing"
)

func GetkubeConfig() (*rest.Config, error) {
	var config *rest.Config
	var err error

	var home = homedir.HomeDir()

	// Initialize Kubernetes client
	config, err = rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	configPath := filepath.Join(home, ".kube/config")
	log.Info().Msgf("loading kube config from : %v", configPath)

	logrus.Infof("loading kube config: %s", configPath)

	config, err = clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		logrus.Fatalf("Error creating dynamic client: %v", err)
	}

	config.TLSClientConfig = rest.TLSClientConfig{
		Insecure: true,
	}

	return config, nil
}

func TestResolveCRD(t *testing.T) {
	restConfig, err := GetkubeConfig()
	if err != nil {
		logrus.Fatal("Error getting config", err)
	}
	k8sClientSet, err := extensioncs.NewForConfig(restConfig)

	resolver := NewResolver(k8sClientSet.ApiextensionsV1().CustomResourceDefinitions())

	crd, err := resolver.Resolve(context.TODO(), "argoproj.io", "Application")

	if crd == nil {
		t.Failed()
	}

}
