package kubeconform

import (
	"github.com/yannh/kubeconform/pkg/config"
	"github.com/yannh/kubeconform/pkg/output"
	"github.com/yannh/kubeconform/pkg/validator"
	"os"
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
