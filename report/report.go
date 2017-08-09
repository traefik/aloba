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

const teamName = "traefik"

var loginMap = map[string]string{
	"timoreimann": "ttr",
	"emilevauge":  "emile",
}

func MakeReport(client *github.Client, ctx context.Context, owner string, repositoryName string) (*model, error) {

	var err error

	members, err := gh.GetTeamMembers(client, ctx, owner, teamName)
	if err != nil {
		return nil, err
	}

	rp := &model{}

	// reviews + status-2 + no contrib/
	rp.withReviews, err = makePRSummaries(client, ctx, owner, repositoryName, members,
		makeWithReview,
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
	rp.noReviews, err = makePRSummaries(client, ctx, owner, repositoryName, nil,
		makeWithoutReview,
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
	rp.contrib, err = makePRSummaries(client, ctx, owner, repositoryName, members,
		makeWithReview,
		search.WithReview,
		search.WithLabels(
			label.StatusNeedsReview,
			label.ContributorWaitingForCorrections),
		search.WithExcludedLabels(label.WIP))
	if err != nil {
		return nil, err
	}

	// design review
	rp.designReview, err = makePRSummaries(client, ctx, owner, repositoryName, nil,
		makeWithoutReview,
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
		fmt.Println(makeMessage(rp.withReviews, true))
	}
	if len(rp.noReviews) != 0 {
		fmt.Println("No reviews:")
		fmt.Println(makeMessage(rp.noReviews, true))
	}
	if len(rp.contrib) != 0 {
		fmt.Println("waiting-for-corrections:")
		fmt.Println(makeMessage(rp.contrib, true))
	}
	if len(rp.designReview) != 0 {
		fmt.Println("Need design review:")
		fmt.Println(makeMessage(rp.designReview, true))
	}
}

func makeMessage(summaries []prSummary, details bool) string {
	var message string
	for _, summary := range summaries {
		message += makeLine(summary, details)
	}
	return message
}

func makeLine(summary prSummary, details bool) string {
	line := fmt.Sprintf("- <https://github.com/containous/traefik/pull/%[1]d|#%[1]d>:", summary.Number)
	line += fmt.Sprintf(" %3d days,", summary.DaysBetweenCreation)
	line += fmt.Sprintf(" %d LGTM", len(summary.Approved))
	if details {
		line += fmt.Sprintf(" %v", summary.Approved)
	}
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

func makeWithReview(client *github.Client, ctx context.Context, owner string, repositoryName string, members []*github.User, issue github.Issue) prSummary {

	approvedReviews, changesRequestedReviews, err := gh.GetReviewStatus(client, ctx, owner, repositoryName, members, issue.GetNumber())
	if err != nil {
		log.Fatal(err)
	}

	crb := []string{}
	for gitHubLogin := range changesRequestedReviews {
		slackLogin, ok := loginMap[gitHubLogin]
		if !ok {
			slackLogin = gitHubLogin
		}

		crb = append(crb, fmt.Sprintf("<@%s>", slackLogin))
	}

	ar := []string{}
	for gitHubLogin := range approvedReviews {
		ar = append(ar, gitHubLogin)
	}

	return newPRSummary(issue, ar, crb)
}

func makeWithoutReview(_ *github.Client, _ context.Context, _ string, _ string, _ []*github.User, issue github.Issue) prSummary {
	return newPRSummary(issue, nil, nil)
}
