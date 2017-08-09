package report

import (
	"context"
	"log"
	"math"
	"strings"
	"time"

	"github.com/containous/myrmica-aloba/internal/search"
	"github.com/containous/myrmica-aloba/label"
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
}

func makePRSummaries(client *github.Client, ctx context.Context,
	owner string, repositoryName string,
	members []*github.User,
	transform func(client *github.Client, ctx context.Context, owner string, repositoryName string, members []*github.User, issue github.Issue) prSummary,
	searchFilter ...search.Parameter) ([]prSummary, error) {

	issues, err := search.FindOpenPR(ctx, client, owner, repositoryName, searchFilter...)
	if err != nil {
		log.Fatal(err)
	}

	var summaries []prSummary

	for _, issue := range issues {
		summary := transform(client, ctx, owner, repositoryName, members, issue)
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func newPRSummary(issue github.Issue, approved []string, requestChanges []string) prSummary {

	var areas []string
	var size string
	for _, lbl := range issue.Labels {
		if strings.HasPrefix(lbl.GetName(), "area/") {
			area := strings.TrimPrefix(lbl.GetName(), "area/")
			area = strings.TrimPrefix(area, "provider/")
			area = strings.TrimPrefix(area, "middleware/")
			areas = append(areas, area)
		}
		if strings.HasPrefix(lbl.GetName(), label.SizeLabelPrefix) {
			size = lbl.GetName()
		}
	}

	return prSummary{
		Number:              issue.GetNumber(),
		Title:               issue.GetTitle(),
		DaysBetweenCreation: int(math.Floor(time.Now().Sub(issue.GetCreatedAt()).Hours() / 24)),
		Approved:            approved,
		Areas:               areas,
		Size:                size,
		ChangesRequested:    requestChanges,
	}
}
