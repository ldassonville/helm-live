package helm

import (
	"github.com/ldassonville/helm-playground/internal/evaluation/utils"
	"github.com/ldassonville/helm-playground/internal/evaluation/utils/pseudoyaml"
	"github.com/rs/zerolog/log"
	"strings"
)

const yamlSeparator = "---\n"

func LoadYAMLManifests(data string) ([]*Manifest, error) {

	var res []*Manifest
	strManifests := strings.Split(data, yamlSeparator)

	for _, strManifest := range strManifests {
		m := &Manifest{}

		if strings.TrimSpace(strManifest) == "" {
			continue
		}

		if err := LoadYAMLManifest(strManifest, m); err == nil {
			log.Err(err).Msgf("fail to parse : %s", strManifest)
		}
		res = append(res, m)
	}
	return res, nil
}

func LoadYAMLManifest(data string, m *Manifest) error {

	manifestMeta := &utils.ManifestMeta{}

	err := utils.UnmarshalYamlManifest([]byte(data), manifestMeta)
	if err != nil {
		m.Name = pseudoyaml.LookupFirstMatch([]string{"metadata", "name"}, data)
		if m.Name == "" {
			m.Name = "name not found"
		}
		m.YamlValid = true
		m.Content = data
		m.GroupVersionKind = GroupVersionKind{
			Kind:    pseudoyaml.LookupFirstMatch([]string{"kind"}, data),
			Group:   pseudoyaml.LookupFirstMatch([]string{"group"}, data),
			Version: pseudoyaml.LookupFirstMatch([]string{"version"}, data),
		}
		return err
	}

	m.Content = data
	m.Namespace = manifestMeta.Namespace
	m.Name = manifestMeta.Name

	//TODO : Yaml error

	m.GroupVersionKind = GroupVersionKind{
		Kind:    manifestMeta.Kind,
		Group:   "",
		Version: manifestMeta.APIVersion,
	}

	// TODO: handle default group
	// split : apps/v1
	if strings.Contains(manifestMeta.APIVersion, "/") {
		parts := strings.Split(manifestMeta.APIVersion, "/")

		m.Version = parts[len(parts)-1]
		m.Group = strings.Join(parts[:len(parts)-1], "/")
	}

	return nil
}
