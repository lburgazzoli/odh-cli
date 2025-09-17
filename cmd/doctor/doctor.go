package doctor

import (
	"fmt"

	"github.com/lburgazzoli/odh-cli/pkg/doctor"
	"github.com/lburgazzoli/odh-cli/pkg/doctor/checks"
	"github.com/lburgazzoli/odh-cli/pkg/printer"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	cmdName  = "doctor"
	cmdShort = "Run diagnostic checks on ODH installation"
)

func AddCommand(root *cobra.Command, flags *genericclioptions.ConfigFlags) {
	printerOpts := printer.Options{
		OutputFormat: printer.Table,
		IOStreams: genericiooptions.IOStreams{
			In:     root.InOrStdin(),
			Out:    root.OutOrStdout(),
			ErrOut: root.ErrOrStderr(),
		},
	}

	cmd := &cobra.Command{
		Use:   cmdName,
		Short: cmdShort,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if printerOpts.OutputFormat != "table" && printerOpts.OutputFormat != "json" {
				return fmt.Errorf("invalid output format: %s (must be 'table' or 'json')", printerOpts.OutputFormat)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := flags.ToRESTConfig()
			if err != nil {
				return fmt.Errorf("failed to create REST config: %w", err)
			}

			c, err := client.New(config, client.Options{})
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			runner := doctor.NewRunner(c, []doctor.DiagnosticCheck{
				checks.NewComponentsCheck(),
			})

			results, err := runner.RunAllChecks()
			if err != nil {
				return fmt.Errorf("failed to run checks: %w", err)
			}

			return printer.NewPrinter(printerOpts).PrintResults(results)
		},
	}

	cmd.Flags().VarP(&printerOpts.OutputFormat, "output", "o", "Output format (table|json)")
	flags.AddFlags(cmd.Flags())

	root.AddCommand(cmd)
}
