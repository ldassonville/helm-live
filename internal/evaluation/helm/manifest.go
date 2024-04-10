package helm

import (
	"github.com/ldassonville/helm-live/internal/evaluation/utils"
	"github.com/ldassonville/helm-live/internal/evaluation/utils/pseudoyaml"
	"github.com/rs/zerolog/log"
	"strings"
)

const yamlSeparator = "---\n"

const helmSourceSeparator = "---\n# Source: "

func LoadHelmYAMLManifests(data string) []*SourceFile {

	var res []*SourceFile
	sourceIndex := make(map[string]*SourceFile)
	strSources := strings.Split(data, helmSourceSeparator)

	for _, strSource := range strSources {

		parts := strings.SplitN(strSource, "\n", 2)
		if len(parts) >= 2 {

			sourcePath := parts[0]
			yamlData := parts[1]
			var src = sourceIndex[sourcePath]

			// If source not found, create a new one
			if src == nil {
				src = &SourceFile{
					Source: sourcePath,
				}
				sourceIndex[sourcePath] = src
				res = append(res, src)
			}

			src.Manifests = append(src.Manifests, loadYAMLManifests(yamlData)...)
		}
	}
	return res
}

func loadYAMLManifests(data string) []*Manifest {

	var res []*Manifest
	strManifests := strings.Split(data, yamlSeparator)

	for _, strManifest := range strManifests {
		m := &Manifest{}

		// Skip empty content
		if strings.TrimSpace(strManifest) == "" {
			continue
		}

		if err := loadYAMLManifest(strManifest, m); err != nil {
			log.Err(err).Msgf("fail to parse : %s", strManifest)
			m.YamlError = err.Error()
		}
		res = append(res, m)
	}
	return res
}

func loadYAMLManifest(data string, m *Manifest) error {

	manifestMeta := &utils.ManifestMeta{}

	m.IsYamlValid = true
	err := utils.UnmarshalYamlManifest([]byte(data), manifestMeta)
	if err != nil {
		m.Name = pseudoyaml.LookupFirstMatch([]string{"metadata", "name"}, data)
		if m.Name == "" {
			m.Name = "name not found"
		}
		m.IsYamlValid = false
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
