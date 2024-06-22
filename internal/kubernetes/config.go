package kubernetes

import (
	"github.com/sirupsen/logrus"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func GetkubeConfig() (*rest.Config, error) {

	k8sFlags := genericclioptions.NewConfigFlags(false)
	var home = homedir.HomeDir()
	configPath := filepath.Join(home, ".kube/config")
	k8sFlags.KubeConfig = &configPath

	config, err := k8sFlags.ToRawKubeConfigLoader().ClientConfig()
	if err != nil {
		return nil, err
	}
	/*	config.TLSClientConfig = rest.TLSClientConfig{
		Insecure: true,
	}*/
	return config, nil
}

func KubernetesDynamicClient() (*dynamic.DynamicClient, error) {

	config, err := GetkubeConfig()
	if err != nil {
		logrus.Fatalf("Error getting config", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("Error creating dynamic client: %v", err)
	}

	return dynamicClient, nil

}
