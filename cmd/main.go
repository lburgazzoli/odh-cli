package main

import (
	"os"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/lburgazzoli/odh-cli/cmd/components"
	"github.com/lburgazzoli/odh-cli/cmd/version"
)

func main() {
	flags := genericclioptions.NewConfigFlags(true)

	cmd := &cobra.Command{
		Use:   "kubectl-odh",
		Short: "kubectl plugin for ODH diagnostic and inspection",
	}

	version.AddCommand(cmd, flags)
	components.AddCommand(cmd, flags)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
