package doctor

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Runner struct {
	client client.Client
	checks []DiagnosticCheck
}

func NewRunner(c client.Client, checks []DiagnosticCheck) *Runner {
	return &Runner{
		client: c,
		checks: checks,
	}
}

func (r *Runner) RunAllChecks() (*CheckResults, error) {
	ctx := context.Background()

	categories := make([]Category, 0)
	summary := Summary{}

	for _, diagnosticCheck := range r.checks {
		checkCategories := diagnosticCheck.Execute(ctx, r.client)
		categories = append(categories, checkCategories...)
	}

	// Aggregate summaries from all categories
	for _, category := range categories {
		categorySummary := ComputeSummary(category)
		summary.OK += categorySummary.OK
		summary.Warning += categorySummary.Warning
		summary.Error += categorySummary.Error
	}

	return &CheckResults{
		Categories: categories,
		Summary:    summary,
	}, nil
}
