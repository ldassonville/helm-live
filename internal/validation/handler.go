package validation

import (
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	"net/http"
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
	router.GET("/schemas/_import", h.ImportSchema)
}

func (h *Handler) ImportSchema(ctx *gin.Context) {

	group := ctx.Query("group")
	kind := ctx.Query("kind")

	err := h.registry.ImportSchema(ctx, group, kind)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)

}
