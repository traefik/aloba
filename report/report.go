package report

import (
	"context"
	"fmt"
	"log"

	"github.com/containous/myrmica-aloba/internal/gh"
	"github.com/containous/myrmica-aloba/internal/search"
	"github.com/containous/myrmica-aloba/label"
	"github.com/google/go-github/github"
)

type model struct {
	withReviews  []prSummary
	noReviews    []prSummary
	contrib      []prSummary
	designReview []prSummary
}

var loginMap = map[string]string{
	"timoreimann": "ttr",
	"emilevauge":  "emile",
}

func MakeReport(client *github.Client, ctx context.Context, owner string, repositoryName string) (*model, error) {

	rp := &model{}
	var err error

	// reviews + status-2 + no contrib/
	rp.withReviews, err = makePRSummaries(client, ctx, owner, repositoryName,
		makeWithLGTM,
		search.WithReview,
		search.WithLabels(label.StatusNeedsReview),
		search.WithExcludedLabels(
			label.WIP,
			label.ContributorWaitingForCorrections,
			label.ContributorWaitingForFeedback,
			label.ContributorWaitingForDocumentation,
			label.ContributorNeedMoreInformation,
			label.StatusNeedsDesignReview,
			label.StatusNeedsMerge))
	if err != nil {
		return nil, err
	}

	// no review + status-2 + no contrib/
	rp.noReviews, err = makePRSummaries(client, ctx, owner, repositoryName,
		makeWithoutLGTM,
		search.WithReviewNone,
		search.WithLabels(label.StatusNeedsReview),
		search.WithExcludedLabels(
			label.WIP,
			label.ContributorWaitingForCorrections,
			label.ContributorWaitingForFeedback,
			label.ContributorWaitingForDocumentation,
			label.ContributorNeedMoreInformation,
			label.StatusNeedsDesignReview,
			label.StatusNeedsMerge))
	if err != nil {
		return nil, err
	}

	// contrib/
	rp.contrib, err = makePRSummaries(client, ctx, owner, repositoryName,
		makeWithLGTM,
		search.WithReview,
		search.WithLabels(
			label.StatusNeedsReview,
			label.ContributorWaitingForCorrections),
		search.WithExcludedLabels(label.WIP))
	if err != nil {
		return nil, err
	}

	// design review
	rp.designReview, err = makePRSummaries(client, ctx, owner, repositoryName,
		makeWithoutLGTM,
		search.WithLabels(label.StatusNeedsDesignReview),
		search.WithExcludedLabels(label.WIP))
	if err != nil {
		return nil, err
	}

	return rp, nil
}

func DisplayReport(rp *model) {
	if len(rp.withReviews) != 0 {
		fmt.Println("With reviews:")
		fmt.Println(makeMessage(rp.withReviews))
	}
	if len(rp.noReviews) != 0 {
		fmt.Println("No reviews:")
		fmt.Println(makeMessage(rp.noReviews))
	}
	if len(rp.contrib) != 0 {
		fmt.Println("waiting-for-corrections:")
		fmt.Println(makeMessage(rp.contrib))
	}
	if len(rp.designReview) != 0 {
		fmt.Println("Need design review:")
		fmt.Println(makeMessage(rp.designReview))
	}
}

func makeMessage(summaries []prSummary) string {
	var message string
	for _, summary := range summaries {
		message += makeLine(summary)
	}
	return message
}

func makeLine(summary prSummary) string {
	line := fmt.Sprintf("- <https://github.com/containous/traefik/pull/%[1]d|#%[1]d>:", summary.Number)
	line += fmt.Sprintf(" %3d days,", summary.DaysBetweenCreation)
	line += fmt.Sprintf(" %d LGTM", summary.LGTM)
	if len(summary.ChangesRequested) != 0 {
		line += fmt.Sprintf(", changes requested by %v", summary.ChangesRequested)
	}
	line += fmt.Sprintf(" -")
	if len(summary.Areas) != 0 {
		line += fmt.Sprintf(" %v", summary.Areas)
	}
	if summary.Size != "" {
		line += fmt.Sprintf("%7s", summary.Size)
	}
	line += fmt.Sprintf(" - _%s_", summary.Title)
	line += fmt.Sprintln()

	return line
}

func makeWithLGTM(client *github.Client, ctx context.Context, owner string, repositoryName string, issue github.Issue) prSummary {

	approvedReviews, changesRequestedReviews, err := gh.GetReviewStatus(client, ctx, owner, repositoryName, issue.GetNumber())
	if err != nil {
		log.Fatal(err)
	}

	crb := []string{}
	for gitHublogin := range changesRequestedReviews {
		slackLogin, ok := loginMap[gitHublogin]
		if !ok {
			slackLogin = gitHublogin
		}

		crb = append(crb, fmt.Sprintf("<@%s>", slackLogin))
	}

	return newPRSummary(issue, len(approvedReviews), crb)
}

func makeWithoutLGTM(_ *github.Client, _ context.Context, _ string, _ string, issue github.Issue) prSummary {
	return newPRSummary(issue, 0, nil)
}
