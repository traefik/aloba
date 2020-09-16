package report

import (
	"context"
	"log"
	"math"
	"strings"
	"time"

	"github.com/google/go-github/v27/github"
	"github.com/traefik/aloba/internal/search"
	"github.com/traefik/aloba/label"
)

type prSummary struct {
	Number              int
	Title               string
	DaysBetweenCreation int
	Approved            []string
	ChangesRequested    []string
	Size                string
	Areas               []string
	Milestone           string
}

type transformer func(ctx context.Context, members []*github.User, issue github.Issue) prSummary

func (r *Reporter) makePRSummaries(ctx context.Context, members []*github.User, transform transformer, searchFilter ...search.Parameter) []prSummary {
	issues, err := search.FindOpenPR(ctx, r.client, r.owner, r.repoName, searchFilter...)
	if err != nil {
		log.Fatal(err)
	}

	var summaries []prSummary

	for _, issue := range issues {
		summary := transform(ctx, members, issue)
		summaries = append(summaries, summary)
	}

	return summaries
}

func newPRSummary(issue github.Issue, approved, requestChanges []string) prSummary {
	var areas []string
	var size string
	for _, lbl := range issue.Labels {
		if strings.HasPrefix(lbl.GetName(), "area/") && !strings.HasPrefix(lbl.GetName(), "area/infrastructure") || strings.HasPrefix(lbl.GetName(), "kind/bug/") {
			area := strings.TrimPrefix(lbl.GetName(), "area/")
			area = strings.TrimPrefix(area, "provider/")
			area = strings.TrimPrefix(area, "middleware/")
			area = strings.TrimPrefix(area, "kind/")
			areas = append(areas, area)
		}
		if strings.HasPrefix(lbl.GetName(), label.SizeLabelPrefix) {
			size = strings.TrimPrefix(lbl.GetName(), "size/")
		}
	}

	return prSummary{
		Number:              issue.GetNumber(),
		Title:               issue.GetTitle(),
		DaysBetweenCreation: int(math.Floor(time.Since(issue.GetCreatedAt()).Hours() / 24)),
		Approved:            approved,
		Areas:               areas,
		Size:                size,
		ChangesRequested:    requestChanges,
		Milestone:           issue.Milestone.GetTitle(),
	}
}
