package helm

func NewRenderBuilder() *renderBuilder {
	return &renderBuilder{
		render: &Render{},
	}
}

type renderBuilder struct {
	render *Render
}

func (rb *renderBuilder) Build() *Render {
	return rb.render
}

func (rb *renderBuilder) AddError(err error) *renderBuilder {
	rb.render.Errors = append(rb.render.Errors, &RenderError{Message: err.Error()})
	return rb
}

func (rb *renderBuilder) SetStatus(status string) *renderBuilder {
	rb.render.Status = status
	return rb
}

func (rb *renderBuilder) AddRenderError(err *RenderError) *renderBuilder {
	rb.render.Errors = append(rb.render.Errors, err)
	return rb
}

func (rb *renderBuilder) AddManifests(manifests []*Manifest) *renderBuilder {
	rb.render.Manifests = append(rb.render.Manifests, manifests...)
	return rb
}

func (rb *renderBuilder) SetMergedValues(values string) {

	rb.render.MergedValues = &Values{
		Data: values,
	}
}
