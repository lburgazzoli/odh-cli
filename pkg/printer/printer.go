package printer

import (
	"encoding/json"
	"io"

	"github.com/fatih/color"
	"github.com/lburgazzoli/odh-cli/pkg/doctor"
	"github.com/lburgazzoli/odh-cli/pkg/printer/table"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgHiYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
)

type Printer interface {
	PrintResults(results *doctor.CheckResults) error
}

func NewPrinter(opts Options) Printer {
	switch opts.OutputFormat {
	case JSON:
		return &JSONPrinter{out: opts.IOStreams.Out}
	case Table:
		return &TablePrinter{out: opts.IOStreams.Out}
	default:
		return &TablePrinter{out: opts.IOStreams.Out}
	}
}

type TablePrinter struct {
	out io.Writer
}

func (p *TablePrinter) PrintResults(results *doctor.CheckResults) error {
	renderer := table.NewRenderer(
		table.WithWriter(p.out),
		table.WithHeaders("CHECK", "STATUS", "MESSAGE"),
		table.WithFormatter("STATUS", func(value interface{}) any {
			v := value.(string)

			switch v {
			case string(doctor.StatusOK):
				v = green(v)
			case string(doctor.StatusWarning):
				v = yellow(v)
			case string(doctor.StatusError):
				v = red(v)
			}

			return v
		}),
	)

	for _, category := range results.Categories {
		if err := p.appendCategoryWithTree(renderer, category); err != nil {
			return err
		}
	}

	return renderer.Render()
}

func (p *TablePrinter) appendCategoryWithTree(renderer *table.Renderer, category doctor.Category) error {
	if err := renderer.Append([]any{category.Name, string(category.Status), category.Message}); err != nil {
		return err
	}

	// Add individual checks with tree formatting
	for i, check := range category.Checks {
		var prefix string
		if i == len(category.Checks)-1 {
			// Last check uses └──
			prefix = "└── "
		} else {
			// Other checks use ├──
			prefix = "├── "
		}

		if err := renderer.Append([]any{prefix + check.Name, string(check.Status), check.Message}); err != nil {
			return err
		}
	}

	return nil
}

type JSONPrinter struct {
	out io.Writer
}

func (p *JSONPrinter) PrintResults(results *doctor.CheckResults) error {
	encoder := json.NewEncoder(p.out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}
