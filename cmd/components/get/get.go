package get

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/lburgazzoli/odh-cli/pkg/components"
	utilclient "github.com/lburgazzoli/odh-cli/pkg/util/client"
)

const (
	cmdName  = "get"
	cmdShort = "Get a specific component by type"
	cmdLong  = `Get a specific ODH/RHOAI component by type name.

The command intelligently matches the component type (case-insensitive) and
returns the singleton instance of that component type.

Components are cluster-scoped resources and follow a singleton pattern - each
type typically has one instance (e.g., "default-kserve", "default-dashboard").

Examples:
  kubectl odh components get kserve
  kubectl odh components get dashboard
  kubectl odh components get DataSciencePipelines`
)

// AddCommand adds the get subcommand to the components command.
func AddCommand(parent *cobra.Command, flags *genericclioptions.ConfigFlags) {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   cmdName + " <component-type>",
		Short: cmdShort,
		Long:  cmdLong,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			componentType := args[0]
			ctx := context.Background()

			client, err := utilclient.NewClient(flags)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			component, err := components.GetComponentByType(ctx, client, componentType)
			if err != nil {
				return fmt.Errorf("failed to get component: %w", err)
			}

			switch outputFormat {
			case "json":
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				if err := encoder.Encode(component); err != nil {
					return fmt.Errorf("failed to encode as JSON: %w", err)
				}
			case "yaml":
				yamlData, err := yaml.Marshal(component)
				if err != nil {
					return fmt.Errorf("failed to marshal as YAML: %w", err)
				}
				fmt.Fprint(cmd.OutOrStdout(), string(yamlData))
			default:
				return fmt.Errorf("unsupported output format: %s (supported: json, yaml)", outputFormat)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output format (json|yaml)")

	parent.AddCommand(cmd)
}
