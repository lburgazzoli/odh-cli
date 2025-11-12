package printer

import (
	"fmt"

	"k8s.io/cli-runtime/pkg/genericiooptions"
)

// OutputFormat specifies the output format for printing results.
type OutputFormat string

const (
	// JSON specifies JSON output format.
	JSON OutputFormat = "json"
	// Table specifies table output format.
	Table OutputFormat = "table"
)

func (f *OutputFormat) String() string {
	return string(*f)
}

// Set sets the output format from a string value.
func (f *OutputFormat) Set(v string) error {
	switch v {
	case string(JSON), string(Table):
		*f = OutputFormat(v)

		return nil
	default:
		return fmt.Errorf("invalid format: %s (must be '%s' or '%s')", v, Table, JSON)
	}
}

// Type returns the type name for the flag value.
func (f *OutputFormat) Type() string {
	return "OutputFormat"
}

// Options contains configuration for creating a printer.
type Options struct {
	IOStreams    genericiooptions.IOStreams
	OutputFormat OutputFormat
}
