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

	YamlValid bool `json:"omitempty"`
}

type RenderError struct {
	Message string `json:"message,omitempty"`
}
