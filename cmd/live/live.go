package main

import (
	"context"
	"encoding/json"
	"github.com/ldassonville/helm-live/internal/evaluation/helm"
	"github.com/ldassonville/helm-live/internal/server"
	"github.com/ldassonville/helm-live/internal/workspace"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func isValuePathValid(valueFile string) bool {
	_, error := os.Stat(valueFile)
	if os.IsNotExist(error) {
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
			valueFile := getAbsolutePath(viper.GetString(pValueFile))
			if !isValuePathValid(valueFile) {
				os.Exit(1)
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

			renderConfig := &helm.RendererConfig{
				SchemaLocations: []string{
					schemaPath,
					"https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json",
				},
				ValuesFile: valueFile,
			}

			var httpServer = server.NewServer()
			renderer := helm.NewRenderer(wks)

			go server.RunServer(httpServer, func() *helm.Render {
				return renderer.Render(context.Background(), renderConfig)
			}, staticPath, port)

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
	rootCmd.Flags().StringP(pValueFile, "f", "./values-example.yaml", "values files")
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
