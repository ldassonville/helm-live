package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/ldassonville/helm-live/internal/evaluation/helm"
	"github.com/ldassonville/helm-live/internal/kubernetes"
	"github.com/ldassonville/helm-live/internal/server"
	"github.com/ldassonville/helm-live/internal/validation"
	"github.com/ldassonville/helm-live/internal/workspace"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	extensioncs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"os"
	"path/filepath"
)

var (
	Commit  = ""
	Version = "x.x.x"
	Date    = ""
)

func main() {

	Execute()
}

const (
	pChartPath  = "chart-path"
	pSchemaPath = "schema-path"
	pValueFile  = "value-file"
	pStaticPath = "static-path"
	pPort       = "port"
)

func init() {
	cobra.OnInitialize(initConfig)
}
func initConfig() {
	viper.SetEnvPrefix("helm-live")
	viper.AutomaticEnv()
}

func getAbsolutePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	absPath, _ := filepath.Abs(path)
	return absPath
}

func isSchemaPathValid(path string) bool {
	return true
}

// isValuePathValid checks if the value file provided is valid
func isValuePathValid(valueFile string) bool {

	// Check the value existence
	_, err := os.Stat(valueFile)
	if os.IsNotExist(err) {
		log.Error().Msgf("Values file %s does not exist", valueFile)
		return false
	}
	return true
}

func isChartPathValid(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Error().Msgf("Chart directory %s does not exist", path)
		return false
	}

	if info.IsDir() {
		// check if the directory contains a Chart.yaml file
		_, err = os.Stat(filepath.Join(path, "Chart.yaml"))
		if os.IsNotExist(err) {
			log.Error().Err(err).Msgf("Invalid chart directory %s. Missing Chart.yaml", path)
			return false
		}

		// check if the templates directory exists
		infoManifest, err := os.Stat(filepath.Join(path, "templates"))
		if os.IsNotExist(err) {
			log.Error().Err(err).Msgf("Invalid chart directory %s. Missing templates directory", path)
			return false
		}
		// check if the templates is a directory
		if !infoManifest.IsDir() {
			log.Error().Err(err).Msgf("Invalid chart directory %s. Manifests is not a directory", path)
			return false
		}
	}
	return !os.IsNotExist(err)
}

func getRootCmd() *cobra.Command {

	var rootCmd = &cobra.Command{
		Use:   "helm-live",
		Short: "helm-live - Launch helm live viewer",
		Long:  `helm-live`,
		Run: func(cmd *cobra.Command, args []string) {

			_ = viper.BindPFlag(pChartPath, cmd.Flags().Lookup(pChartPath))
			chartPath := getAbsolutePath(viper.GetString(pChartPath))
			if !isChartPathValid(chartPath) {
				os.Exit(1)
			}

			_ = viper.BindPFlag(pSchemaPath, cmd.Flags().Lookup(pSchemaPath))
			schemaPath := getAbsolutePath(viper.GetString(pSchemaPath))
			if !isSchemaPathValid(schemaPath) {
				os.Exit(1)
			}

			_ = viper.BindPFlag(pValueFile, cmd.Flags().Lookup(pValueFile))
			valueFile := viper.GetString(pValueFile)
			if valueFile != "" {
				// If value file is defined, check if it exists
				valueFile = getAbsolutePath(valueFile)
				if !isValuePathValid(valueFile) {
					os.Exit(1)
				}
			}

			_ = viper.BindPFlag(pStaticPath, cmd.Flags().Lookup(pStaticPath))
			staticPath := viper.GetString(pStaticPath)
			if staticPath != "" {
				staticPath = getAbsolutePath(staticPath)
			}

			_ = viper.BindPFlag(pPort, cmd.Flags().Lookup(pPort))
			port := viper.GetInt(pPort)

			wks := &workspace.Workspace{
				Path: chartPath,
			}

			schemaLocations := []string{
				schemaPath,
				"https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json",
			}

			renderConfig := &helm.RendererConfig{
				SchemaLocations: schemaLocations,
				ValuesFile:      valueFile,
			}

			var httpServer = server.NewServer()
			renderer := helm.NewRenderer(wks)

			renderRegister := func(router *gin.Engine) {
				// Add one shoot render endpoint
				router.GET("/_render", func(c *gin.Context) {

					render := renderer.Render(context.Background(), renderConfig)
					c.JSON(200, render)
				})
			}

			// Initialise the kubernetes client
			restConfig, _ := kubernetes.GetkubeConfig()
			k8sClientSet, _ := extensioncs.NewForConfig(restConfig)
			crdClient := k8sClientSet.ApiextensionsV1().CustomResourceDefinitions()

			validationHandler := validation.New(schemaPath, schemaLocations, crdClient)

			var handlerRegisters = []func(engine *gin.Engine){
				renderRegister, validationHandler.Register,
			}
			go server.RunServer(httpServer, handlerRegisters, staticPath, port)

			// Web socket to notify the client when the workspace changes
			wks.OnChangeFnc = func() {
				render := renderer.Render(context.Background(), renderConfig)
				strJson, _ := json.Marshal(render)
				httpServer.Message <- string(strJson)
			}

			wks.Live()
		},
	}

	rootCmd.Flags().StringP(pSchemaPath, "s", "./catalog-crds-json-schema/{{.Group}}/{{.ResourceKind}}_{{.ResourceAPIVersion}}.json", "schemas path ")
	rootCmd.Flags().StringP(pChartPath, "c", ".", "Chart path")
	rootCmd.Flags().StringP(pValueFile, "f", "", "values files")
	rootCmd.Flags().String(pStaticPath, "", "statics path files")
	rootCmd.Flags().StringP(pPort, "p", "8085", "HTTP Port used")
	return rootCmd
}

func Execute() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := getRootCmd().Execute(); err != nil {
		log.Error().Err(err).Msg("There was an error while executing your CLI")
		os.Exit(1)
	}
}
