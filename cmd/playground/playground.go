package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ldassonville/helm-playground/internal/evaluation/helm"
	"github.com/ldassonville/helm-playground/internal/server"
	"github.com/ldassonville/helm-playground/internal/workspace"
	"github.com/spf13/cobra"
	"os"
)

var (
	Commit  = ""
	Version = "x.x.x"
	Date    = ""
)

func main() {

	Execute()

}

var rootCmd = &cobra.Command{
	Use:   "helm-playground",
	Short: "helm-playground - Launch helm playground",
	Long:  `helm-playground`,
	Run: func(cmd *cobra.Command, args []string) {
		wks := &workspace.Workspace{
			Path: "/home/ldassonville/git/leroymerlin/packaging/aap-core-argocd/helm-expo",
		}

		var httpServer = server.NewServer()
		go server.RunServer(httpServer)

		renderer := helm.NewRenderer(wks)

		wks.OnChangeFnc = func() {
			render := renderer.Render(context.Background(), "/home/ldassonville/git/leroymerlin/packaging/aap-core-argocd/helm-expo/values-example.yaml")

			strJson, _ := json.Marshal(render)
			httpServer.Message <- string(strJson)
		}

		wks.Live()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
