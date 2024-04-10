package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ldassonville/helm-live/internal/evaluation/helm"
	"github.com/ldassonville/helm-live/internal/server"
	"github.com/ldassonville/helm-live/internal/workspace"
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
		fmt.Println(fmt.Sprintf("Values file %s does not exist", valueFile))
		return false
	}
	return true
}

func isChartPathValid(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("Chart directory %s does not exist", path))
		return false
	}

	if info.IsDir() {
		// check if the directory contains a Chart.yaml file
		_, err = os.Stat(filepath.Join(path, "Chart.yaml"))
		if os.IsNotExist(err) {
			return false
		}

		// check if the templates directory exists
		infoManifest, err := os.Stat(filepath.Join(path, "templates"))
		if os.IsNotExist(err) {
			return false
		}
		// check if the templates is a directory
		if !infoManifest.IsDir() {
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
			staticPath := getAbsolutePath(viper.GetString(pStaticPath))

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
			}, staticPath)

			wks.OnChangeFnc = func() {
				render := renderer.Render(context.Background(), renderConfig)
				strJson, _ := json.Marshal(render)
				httpServer.Message <- string(strJson)
			}

			wks.Live()
		},
	}

	rootCmd.Flags().String(pSchemaPath, "./catalog-crds-json-schema/{{.Group}}/{{.ResourceKind}}_{{.ResourceAPIVersion}}.json", "schemas path ")
	rootCmd.Flags().String(pChartPath, ".", "Chart path")
	rootCmd.Flags().String(pValueFile, "./values-example.yaml", "values files")
	rootCmd.Flags().String(pStaticPath, "./statics", "statics path files")

	return rootCmd
}

func Execute() {
	if err := getRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
