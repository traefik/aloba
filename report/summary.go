package report

import (
	"context"
	"log"
	"math"
	"strings"
	"time"

	"github.com/containous/aloba/internal/search"
	"github.com/containous/aloba/label"
	"github.com/google/go-github/github"
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

func makePRSummaries(ctx context.Context, client *github.Client,
	owner string, repositoryName string,
	members []*github.User,
	transform func(ctx context.Context, client *github.Client, owner string, repositoryName string, members []*github.User, issue github.Issue) prSummary,
	searchFilter ...search.Parameter) []prSummary {

	issues, err := search.FindOpenPR(ctx, client, owner, repositoryName, searchFilter...)
	if err != nil {
		log.Fatal(err)
	}

	var summaries []prSummary

	for _, issue := range issues {
		summary := transform(ctx, client, owner, repositoryName, members, issue)
		summaries = append(summaries, summary)
	}

	return summaries
}

func newPRSummary(issue github.Issue, approved []string, requestChanges []string) prSummary {
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
