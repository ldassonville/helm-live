package kubeconform

import (
	"bytes"
	"github.com/yannh/kubeconform/pkg/config"
	"github.com/yannh/kubeconform/pkg/output"
	"github.com/yannh/kubeconform/pkg/validator"
	"os"
	"strings"
	"text/template"
)

func Validate(cfg config.Config) (validator.Validator, error) {
	var o output.Output
	var useStdin bool = false
	var err error

	if o, err = output.New(os.Stdout, cfg.OutputFormat, cfg.Summary, useStdin, cfg.Verbose); err != nil {
		return nil, err
	}

	println(o)

	var v validator.Validator
	v, err = validator.New(cfg.SchemaLocations, validator.Opts{
		Cache:                cfg.Cache,
		Debug:                cfg.Debug,
		SkipTLS:              cfg.SkipTLS,
		SkipKinds:            cfg.SkipKinds,
		RejectKinds:          cfg.RejectKinds,
		KubernetesVersion:    cfg.KubernetesVersion,
		Strict:               cfg.Strict,
		IgnoreMissingSchemas: cfg.IgnoreMissingSchemas,
	})
	if err != nil {
		return nil, err
	}

	return v, nil
}

// GetSchemaPath returns the path to the schema file for a given resource kind and API version
func GetSchemaPath(tpl, resourceKind, resourceAPIVersion, k8sVersion string, strict bool) (string, error) {
	normalisedVersion := k8sVersion
	if normalisedVersion != "master" {
		normalisedVersion = "v" + normalisedVersion
	}

	strictSuffix := ""
	if strict {
		strictSuffix = "-strict"
	}

	groupParts := strings.Split(resourceAPIVersion, "/")
	versionParts := strings.Split(groupParts[0], ".")

	kindSuffix := "-" + strings.ToLower(versionParts[0])
	if len(groupParts) > 1 {
		kindSuffix += "-" + strings.ToLower(groupParts[1])
	}

	tmpl, err := template.New("tpl").Parse(tpl)
	if err != nil {
		return "", err
	}

	tplData := struct {
		NormalizedKubernetesVersion string
		StrictSuffix                string
		ResourceKind                string
		ResourceAPIVersion          string
		Group                       string
		KindSuffix                  string
	}{
		normalisedVersion,
		strictSuffix,
		strings.ToLower(resourceKind),
		groupParts[len(groupParts)-1],
		groupParts[0],
		kindSuffix,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tplData)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
