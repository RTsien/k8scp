package main

import (
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/rtsien/k8scp/pkg/common"
	"github.com/rtsien/k8scp/pkg/svr"
)

var kubeconfig string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "svr",
	Short: "A server for copying files to K8s pod.",
	Long:  `A server for copying files to K8s pod.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		kByte, err := os.ReadFile(kubeconfig)
		common.AssertErr(err, "kubeconfig file not found: %s", kubeconfig)

		http.HandleFunc("/upload", svr.UploadHandler(string(kByte)))

		if err = http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "kubeconfig file path")
	_ = rootCmd.MarkFlagRequired("kubeconfig")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func main() {
	Execute()
}
