package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/viveksinghggits/kreate/pkg/utils"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/kubectl/pkg/cmd/util"
)

var rootCmd = &cobra.Command{
	Use:   utils.Name,
	Short: "Create Kubernetes resources that are not supported by `kubectl create`",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	ioStreams := genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	rootCmd.AddCommand(NewCmdCreatePVC(
		util.NewFactory(util.NewMatchVersionFlags(defaultConfigFlags().WithWarningPrinter(ioStreams))),
		ioStreams,
	))

}

func defaultConfigFlags() *genericclioptions.ConfigFlags {
	return genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
}
