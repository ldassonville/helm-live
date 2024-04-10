package helm

import (
	"sort"
	"time"
)

func NewRenderBuilder() *renderBuilder {
	return &renderBuilder{
		render: &Render{},
	}
}

type renderBuilder struct {
	render *Render
}

func (rb *renderBuilder) Build() *Render {

	sort.SliceStable(rb.render.Sources, func(i, j int) bool {
		m1 := rb.render.Sources[i]
		m2 := rb.render.Sources[j]

		if m1.Source == m2.Source {
			return m1.Source < m2.Source
		}
		return m1.Source < m2.Source
	})

	for _, source := range rb.render.Sources {
		sort.SliceStable(source.Manifests, func(i, j int) bool {
			m1 := source.Manifests[i]
			m2 := source.Manifests[j]

			if m1.Name == m2.Name {
				return m1.Kind < m2.Kind
			}
			return m1.Name < m2.Name
		})
	}

	return rb.render
}

func (rb *renderBuilder) AddError(err error) *renderBuilder {

	rb.render.Errors = append(rb.render.Errors, &RenderError{Message: err.Error()})
	return rb
}

func (rb *renderBuilder) getInfo() *Info {
	if rb.render.Info == nil {
		rb.render.Info = &Info{}
	}
	return rb.render.Info
}

func (rb *renderBuilder) SetExecutionTime(executionTime time.Time) *renderBuilder {

	rb.getInfo().ExecutionTime = executionTime
	return rb
}

func (rb *renderBuilder) SetStatus(status string) *renderBuilder {

	rb.getInfo().Status = status
	return rb
}

func (rb *renderBuilder) AddRenderError(err *RenderError) *renderBuilder {
	rb.render.Errors = append(rb.render.Errors, err)
	return rb
}

func (rb *renderBuilder) AddSources(sources []*SourceFile) *renderBuilder {

	rb.render.Sources = append(rb.render.Sources, sources...)
	return rb
}

func (rb *renderBuilder) SetMergedValues(values string) {

	rb.render.MergedValues = &Values{
		Data: values,
	}
}

func (rb *renderBuilder) SetRawManifest(manifest string) {

	rb.render.RawManifest = manifest
}

func (rb *renderBuilder) SetChartPath(path string) {
	rb.render.ChartPath = path
}

func (rb *renderBuilder) SetChartName(name string) {
	rb.render.ChartName = name
}
