package table_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/lburgazzoli/odh-cli/pkg/printer/table"

	. "github.com/onsi/gomega"
)

type testPerson struct {
	Name   string
	Age    int
	Status string
}

type testPersonWithTags struct {
	Name     string
	Tags     []string
	Metadata map[string]any
}

func TestRendererWithSliceInput(t *testing.T) {
	g := NewWithT(t)

	var buf bytes.Buffer
	renderer := table.NewRenderer(
		table.WithWriter(&buf),
		table.WithHeaders("Name", "Age"),
	)

	err := renderer.Append([]any{"Alice", 30})
	g.Expect(err).ShouldNot(HaveOccurred())

	err = renderer.Render()
	g.Expect(err).ShouldNot(HaveOccurred())

	output := buf.String()
	g.Expect(output).Should(ContainSubstring("Alice"))
	g.Expect(output).Should(ContainSubstring("30"))
}

func TestRendererWithStructInput(t *testing.T) {
	g := NewWithT(t)

	var buf bytes.Buffer
	renderer := table.NewRenderer(
		table.WithWriter(&buf),
		table.WithHeaders("Name", "Age", "Status"),
	)

	person := testPerson{
		Name:   "Alice",
		Age:    30,
		Status: "active",
	}

	err := renderer.Append(person)
	g.Expect(err).ShouldNot(HaveOccurred())

	err = renderer.Render()
	g.Expect(err).ShouldNot(HaveOccurred())

	output := buf.String()
	g.Expect(output).Should(ContainSubstring("Alice"))
	g.Expect(output).Should(ContainSubstring("30"))
	g.Expect(output).Should(ContainSubstring("active"))
}

func TestRendererWithCustomFormatter(t *testing.T) {
	g := NewWithT(t)

	var buf bytes.Buffer
	renderer := table.NewRenderer(
		table.WithWriter(&buf),
		table.WithHeaders("Name", "Status"),
		table.WithFormatter("Name", func(v any) any {
			return strings.ToUpper(v.(string))
		}),
	)

	person := testPerson{
		Name:   "Alice",
		Status: "active",
	}

	err := renderer.Append(person)
	g.Expect(err).ShouldNot(HaveOccurred())

	err = renderer.Render()
	g.Expect(err).ShouldNot(HaveOccurred())

	output := buf.String()
	g.Expect(output).Should(ContainSubstring("ALICE"))
}

func TestRendererWithJQFormatter(t *testing.T) {
	g := NewWithT(t)

	var buf bytes.Buffer
	renderer := table.NewRenderer(
		table.WithWriter(&buf),
		table.WithHeaders("Name", "Tags"),
		table.WithFormatter("Tags", table.JQFormatter(`. | join(", ")`)),
	)

	person := testPersonWithTags{
		Name: "Alice",
		Tags: []string{"admin", "user"},
	}

	err := renderer.Append(person)
	g.Expect(err).ShouldNot(HaveOccurred())

	err = renderer.Render()
	g.Expect(err).ShouldNot(HaveOccurred())

	output := buf.String()
	g.Expect(output).Should(ContainSubstring("Alice"))
	g.Expect(output).Should(ContainSubstring("admin, user"))
}

func TestRendererWithChainedFormatters(t *testing.T) {
	g := NewWithT(t)

	var buf bytes.Buffer
	renderer := table.NewRenderer(
		table.WithWriter(&buf),
		table.WithHeaders("Name", "Status"),
		table.WithFormatter("Name",
			table.ChainFormatters(
				table.JQFormatter("."),
				func(v any) any {
					return strings.ToUpper(v.(string))
				},
				func(v any) any {
					return "[" + v.(string) + "]"
				},
			),
		),
	)

	person := testPerson{
		Name:   "Alice",
		Status: "active",
	}

	err := renderer.Append(person)
	g.Expect(err).ShouldNot(HaveOccurred())

	err = renderer.Render()
	g.Expect(err).ShouldNot(HaveOccurred())

	output := buf.String()
	g.Expect(output).Should(ContainSubstring("[ALICE]"))
}

func TestRendererWithJQExtraction(t *testing.T) {
	g := NewWithT(t)

	var buf bytes.Buffer
	renderer := table.NewRenderer(
		table.WithWriter(&buf),
		table.WithHeaders("Name", "Metadata"),
		table.WithFormatter("Metadata",
			table.JQFormatter(`.location // "Unknown"`),
		),
	)

	person := testPersonWithTags{
		Name:     "Alice",
		Metadata: map[string]any{"location": "NYC"},
	}

	err := renderer.Append(person)
	g.Expect(err).ShouldNot(HaveOccurred())

	err = renderer.Render()
	g.Expect(err).ShouldNot(HaveOccurred())

	output := buf.String()
	g.Expect(output).Should(ContainSubstring("Alice"))
	g.Expect(output).Should(ContainSubstring("NYC"))
}

func TestRendererAppendAll(t *testing.T) {
	g := NewWithT(t)

	var buf bytes.Buffer
	renderer := table.NewRenderer(
		table.WithWriter(&buf),
		table.WithHeaders("Name", "Age"),
	)

	people := []any{
		testPerson{Name: "Alice", Age: 30},
		testPerson{Name: "Bob", Age: 25},
		testPerson{Name: "Charlie", Age: 35},
	}

	err := renderer.AppendAll(people)
	g.Expect(err).ShouldNot(HaveOccurred())

	err = renderer.Render()
	g.Expect(err).ShouldNot(HaveOccurred())

	output := buf.String()
	g.Expect(output).Should(ContainSubstring("Alice"))
	g.Expect(output).Should(ContainSubstring("Bob"))
	g.Expect(output).Should(ContainSubstring("Charlie"))
}

func TestRendererCaseInsensitiveMatching(t *testing.T) {
	g := NewWithT(t)

	var buf bytes.Buffer
	renderer := table.NewRenderer(
		table.WithWriter(&buf),
		table.WithHeaders("name", "AGE"),
	)

	person := testPerson{
		Name: "Alice",
		Age:  30,
	}

	err := renderer.Append(person)
	g.Expect(err).ShouldNot(HaveOccurred())

	err = renderer.Render()
	g.Expect(err).ShouldNot(HaveOccurred())

	output := buf.String()
	g.Expect(output).Should(ContainSubstring("Alice"))
	g.Expect(output).Should(ContainSubstring("30"))
}
