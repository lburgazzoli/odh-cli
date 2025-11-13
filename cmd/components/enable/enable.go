package enable

import (
	"fmt"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	cmdName  = "enable"
	cmdShort = "Enable a component"
	cmdLong  = `Enable an ODH/RHOAI component.

This command is not yet implemented.`
)

// AddCommand adds the enable subcommand to the components command.
func AddCommand(parent *cobra.Command, _ *genericclioptions.ConfigFlags) {
	cmd := &cobra.Command{
		Use:   cmdName + " <component-name>",
		Short: cmdShort,
		Long:  cmdLong,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement enable functionality
			return fmt.Errorf("enable command not yet implemented")
		},
	}

	parent.AddCommand(cmd)
}

