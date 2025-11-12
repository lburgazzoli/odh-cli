package table

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/olekukonko/tablewriter"

	"github.com/lburgazzoli/odh-cli/pkg/util"
)

// Option is a functional option for configuring a Renderer.
type Option = util.Option[Renderer]

// WithWriter sets the output writer for the table renderer.
func WithWriter(w io.Writer) Option {
	return util.FunctionalOption[Renderer](func(r *Renderer) {
		r.writer = w
	})
}

// WithHeaders sets the column headers for the table.
func WithHeaders(headers ...string) Option {
	return util.FunctionalOption[Renderer](func(r *Renderer) {
		r.headers = headers
	})
}

// WithFormatter adds a column-specific formatter function.
func WithFormatter(columnName string, formatter ColumnFormatter) Option {
	return util.FunctionalOption[Renderer](func(r *Renderer) {
		if r.formatters == nil {
			r.formatters = make(map[string]ColumnFormatter)
		}

		r.formatters[strings.ToUpper(columnName)] = formatter
	})
}

// WithTableOptions sets the underlying tablewriter options.
func WithTableOptions(values ...tablewriter.Option) Option {
	return util.FunctionalOption[Renderer](func(r *Renderer) {
		r.tableOptions = append(r.tableOptions, values...)
	})
}

// JQFormatter creates a ColumnFormatter that executes a jq query on the input value.
// The query is compiled once at creation time.
// Panics if the query compilation fails (fail fast at setup time).
func JQFormatter(query string) ColumnFormatter {
	compiledQuery, err := gojq.Parse(query)
	if err != nil {
		panic("failed to compile jq query: " + err.Error())
	}

	return func(value any) any {
		// Convert Go types to JSON-compatible types for gojq
		// This handles slices, maps, structs, etc.
		var normalizedValue any
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return err.Error()
		}

		if err := json.Unmarshal(jsonBytes, &normalizedValue); err != nil {
			return err.Error()
		}

		// Run the query against the normalized value
		iter := compiledQuery.Run(normalizedValue)

		// Get the first result
		result, ok := iter.Next()
		if !ok {
			return nil
		}

		// Check for errors
		if err, isErr := result.(error); isErr {
			// Return error as string for display
			return err.Error()
		}

		return result
	}
}

// ChainFormatters composes multiple formatters into a single formatter pipeline.
// The output of each formatter is passed as input to the next formatter.
// This enables building transformation pipelines like: JQ extraction → colorization → truncation.
func ChainFormatters(formatters ...ColumnFormatter) ColumnFormatter {
	if len(formatters) == 0 {
		return func(value any) any {
			return value
		}
	}

	if len(formatters) == 1 {
		return formatters[0]
	}

	return func(value any) any {
		result := value
		for _, formatter := range formatters {
			result = formatter(result)
		}

		return result
	}
}
