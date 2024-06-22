package validation

import (
	"github.com/gin-gonic/gin"
	"github.com/ldassonville/helm-live/internal/validation/schema"
	v1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	"net/http"
)

func New(schemaPath string, schemaLocations []string, crdi v1.CustomResourceDefinitionInterface) *Handler {
	return &Handler{
		crdSchemaResolver: schema.New(schemaLocations),
		registry:          NewRegistry(schemaPath, crdi),
	}
}

type Handler struct {
	crdSchemaResolver schema.Resolver
	registry          *Registry
}

func (h *Handler) Register(router *gin.Engine) {
	router.GET("/schemas/_import", h.ImportSchema)
	router.GET("/schemas/crd/:group/:version/:kind", h.ResolveCRDSchema)
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

func (h *Handler) ResolveCRDSchema(ctx *gin.Context) {

	group := ctx.Param("group")
	version := ctx.Param("version")
	kind := ctx.Param("kind")

	data, err := h.crdSchemaResolver.ResolveSchema(group, kind, version)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "application/json", data)

}
