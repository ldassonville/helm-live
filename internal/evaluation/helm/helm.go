package helm

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
	"time"
)

import "helm.sh/helm/v3/pkg/action"

type Helm struct {
	actionConfig *action.Configuration
}

func PrintValues(options values.Options) {

	for _, v := range options.ValueFiles {

		dat, err := os.ReadFile(v)
		if err != nil {
			continue
		}

		log.Info().Msgf("ValueFile content : %s", v)
		fmt.Printf(string(dat))
	}
}

func GetMergedValues(options values.Options) (string, error) {
	merge, err := options.MergeValues(getter.Providers{})
	if err != nil {
		log.Log().Err(err).Msg("fail to merge values")
		return "", err
	}

	dat, err := yaml.Marshal(merge)
	if err != nil {
		log.Log().Err(err).Msg("fail to merge values")
		return "", err
	}

	return string(dat), nil
}

func PrintMergedValues(options values.Options) {

	dat, _ := GetMergedValues(options)
	fmt.Printf(string(dat))

}

// ValuesHasFile create a temporary file with the given yaml content
func valuesHasFile(yaml string) (string, error) {

	file, err := os.CreateTemp("/tmp", "kube-assistant-values")
	if err != nil {
		return "", err
	}

	defer file.Close()
	_, err = file.WriteString(yaml)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

// NewLocalHelm create a new Helm client contextualized to the given namespace
func NewLocalHelm(namespace string) (*Helm, error) {
	return NewHelm(&genericclioptions.ConfigFlags{}, namespace)
}

func NewHelm(k8sFlags *genericclioptions.ConfigFlags, namespace string) (*Helm, error) {
	logger := func(msg string, args ...interface{}) {
		log.Debug().Msgf(msg, args...)
	}

	k8sFlags.Namespace = &namespace
	actionConfig := &action.Configuration{}

	actionConfig.Init(k8sFlags, namespace, "memory", logger)

	// create the HELM client
	return &Helm{
		//k8sFlags:     k8sFlags,
		actionConfig: actionConfig,
	}, nil
}

// Template render the helm chart.
// Path is the path to the chart
// Release is the release name for the templating
// valueOpts is the values to apply to the chart
func (h *Helm) Template(path string, release string, valueOpts values.Options) (result *release.Release, err error) {

	cmd := action.NewInstall(h.actionConfig)

	p := getter.Providers{}

	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		log.Err(err).Msg("fail to merge value")
		return nil, err
	}

	chart, err := loader.Load(path)

	cmd.DryRun = true
	cmd.ClientOnly = true
	cmd.InsecureSkipTLSverify = true
	cmd.Timeout = 20 * time.Second
	cmd.Wait = false
	cmd.APIVersions = chartutil.DefaultVersionSet
	cmd.ReleaseName = release
	cmd.Devel = true

	resp, err := cmd.Run(chart, vals)

	if err != nil {
		log.Err(err).
			Str("release", release).
			Msg("helm.Template failed")

	}

	return resp, err
}
