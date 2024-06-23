package kubernetes

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
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
	config.TLSClientConfig = rest.TLSClientConfig{
		Insecure: true,
	}
	return config, nil
}
