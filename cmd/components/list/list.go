package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/lburgazzoli/odh-cli/pkg/components"
	"github.com/lburgazzoli/odh-cli/pkg/printer/table"
	utilclient "github.com/lburgazzoli/odh-cli/pkg/util/client"
)

const (
	cmdName  = "list"
	cmdShort = "List all components"
	cmdLong  = `List all ODH/RHOAI components from the components.platform.opendatahub.io API group.

Components are cluster-scoped resources.`
)

// AddCommand adds the list subcommand to the components command.
func AddCommand(parent *cobra.Command, flags *genericclioptions.ConfigFlags) {
	var outputFormat string

	cmd := &cobra.Command{
		Use:     cmdName,
		Aliases: []string{"ls"},
		Short:   cmdShort,
		Long:    cmdLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			client, err := utilclient.NewClient(flags)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			componentList, err := components.ListComponents(ctx, client)
			if err != nil {
				return fmt.Errorf("failed to list components: %w", err)
			}

			switch outputFormat {
			case "json":
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")

				if err := encoder.Encode(componentList); err != nil {
					return fmt.Errorf("failed to encode components as JSON: %w", err)
				}

				return nil
			case "yaml":
				yamlData, err := yaml.Marshal(componentList)
				if err != nil {
					return fmt.Errorf("failed to marshal as YAML: %w", err)
				}
				fmt.Fprint(cmd.OutOrStdout(), string(yamlData))
				return nil
			case "table":
				renderer := table.NewWithColumns[unstructured.Unstructured](
					cmd.OutOrStdout(),
					table.NewColumn("TYPE").
						JQ(`.kind`),
					table.NewColumn("READY").
						JQ(`.status.conditions[]? | select(.type=="Ready") | .status // "Unknown"`),
					table.NewColumn("MESSAGE").
						JQ(`.status.conditions[]? | select(.type=="Ready") | .message // ""`),
				)

				if err := renderer.AppendAll(componentList.Items); err != nil {
					return fmt.Errorf("failed to append rows: %w", err)
				}

				if err := renderer.Render(); err != nil {
					return fmt.Errorf("failed to render table: %w", err)
				}

				return nil
			default:
				return fmt.Errorf("unsupported output format: %s (supported: table, json, yaml)", outputFormat)
			}
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table|json|yaml)")

	parent.AddCommand(cmd)
}
