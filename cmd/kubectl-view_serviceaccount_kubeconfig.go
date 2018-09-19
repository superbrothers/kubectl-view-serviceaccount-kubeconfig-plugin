package main

import (
	"flag"
	"os"

	"github.com/spf13/pflag"
	"github.com/superbrothers/kubectl-view-serviceaccount-kubeconfig-plugin/pkg/cmd"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func init() {
	// Initialize glog flags
	flag.CommandLine.Set("logtostderr", "true")
}

func main() {
	flags := pflag.NewFlagSet("kubectl-view-serviceaccount-kubeconfig", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewCmdViewServiceaccountKubeconfig(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
