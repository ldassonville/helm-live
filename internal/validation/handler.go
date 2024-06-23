package validation

import (
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
)

func New(schemaPath string, crdi v1.CustomResourceDefinitionInterface) *Handler {
	return &Handler{
		registry: NewRegistry(schemaPath, crdi),
	}
}

type Handler struct {
	registry *Registry
}

func (h *Handler) Register(router *gin.Engine) {
	router.GET("/import-schema", h.ImportSchema)
}

func (h *Handler) ImportSchema(ctx *gin.Context) {

	group := ctx.Query("group")
	version := ctx.Query("version")
	kind := ctx.Query("kind")

	h.registry.ImportSchema(ctx, group, version, kind)
}
