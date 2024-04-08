package utils

import (
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type ManifestMeta struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

func UnmarshalYamlManifest(yamlData []byte, destinationObj interface{}) error {

	jsonData, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, destinationObj)
	if err != nil {
		return err
	}
	return nil
}
