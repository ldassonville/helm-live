package helm

type GroupVersionKind struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

type Manifest struct {
	GroupVersionKind `json:",inline"`
	Name             string `json:"name,omitempty"`
	Namespace        string `json:"manifest,omitempty"`
	Content          string `json:"content,omitempty"`

	// KubeConformValidation result
	KubeConformValidation *KubeConformValidation `json:"kubeConformValidation,omitempty"`

	IsYamlValid bool   `json:"isYamlValid"`
	YamlError   string `json:"yamlError,omitempty"`
}

type KubeConformValidation struct {
	Status           string            `json:"status,omitempty"`
	ErrMsg           string            `json:"errMsg,omitempty"`
	ValidationErrors []ValidationError `json:"validationErrors,omitempty"`
}

type ValidationError struct {
	Path string `json:"path"`
	Msg  string `json:"msg"`
}

type RenderError struct {
	Message string `json:"message,omitempty"`
}
