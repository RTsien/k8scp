package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/rtsien/k8scp/pkg/scp"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scp",
	Short: "A command line tool for copying files to K8s pods.",
	Long:  `A command line tool for copying files to K8s pods.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if err := cp.Do(); err != nil {
			os.Exit(1)
		}
	},
}
var cp scp.Copy

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scp.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&cp.Src, "src", "s", "", "source file path")
	_ = rootCmd.MarkFlagRequired("src")
	rootCmd.Flags().StringVarP(&cp.ServerURL, "url", "u", "", "server url")
	_ = rootCmd.MarkFlagRequired("url")
	rootCmd.Flags().StringVarP(&cp.Dst, "dst", "d", "", "destination file path")
	_ = rootCmd.MarkFlagRequired("dst")
	rootCmd.Flags().StringVarP(&cp.Namespace, "namespace", "n", "", "k8s namespace")
	_ = rootCmd.MarkFlagRequired("namespace")
	rootCmd.Flags().StringVarP(&cp.Pod, "pod", "p", "", "pod name")
	_ = rootCmd.MarkFlagRequired("pod")
	rootCmd.Flags().StringVarP(&cp.Container, "container", "c", "", "container name")
	_ = rootCmd.MarkFlagRequired("container")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func main() {
	Execute()
}
