package printer

import (
	"fmt"

	"k8s.io/cli-runtime/pkg/genericiooptions"
)

type OutputFormat string

const (
	JSON  OutputFormat = "json"
	Table OutputFormat = "table"
)

func (f *OutputFormat) String() string {
	return string(*f)
}

func (f *OutputFormat) Set(v string) error {
	switch v {
	case string(JSON), string(Table):
		*f = OutputFormat(v)
		return nil
	default:
		return fmt.Errorf("invalid format: %s (must be '%s' or '%s')", v, Table, JSON)
	}
}

func (f *OutputFormat) Type() string {
	return "OutputFormat"
}

type Options struct {
	IOStreams    genericiooptions.IOStreams
	OutputFormat OutputFormat
}
