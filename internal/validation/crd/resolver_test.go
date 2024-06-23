package crd

import (
	"context"
	"github.com/ldassonville/helm-live/internal/kubernetes"
	"github.com/sirupsen/logrus"
	extensioncs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"testing"
)

func TestResolveCRD(t *testing.T) {
	restConfig, err := kubernetes.GetkubeConfig()
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
