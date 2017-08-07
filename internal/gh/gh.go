package gh

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	approved         = "APPROVED"
	changesRequested = "CHANGES_REQUESTED"
	commented        = "COMMENTED"
)

func NewGitHubClient(ctx context.Context, token string) *github.Client {
	var client *github.Client
	if len(token) == 0 {
		client = github.NewClient(nil)
	} else {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}
	return client
}

func GetReviewStatus(client *github.Client, ctx context.Context, owner string, repositoryName string, prNumber int) (map[string]string, map[string]string, error) {
	opts := &github.ListOptions{
		PerPage: 50,
	}

	reviews, _, err := client.PullRequests.ListReviews(ctx, owner, repositoryName, prNumber, opts)
	if err != nil {
		return nil, nil, err
	}

	uniqueReviews := make(map[string]string)
	for _, review := range reviews {
		if review.GetState() != commented {
			uniqueReviews[review.User.GetLogin()] = review.GetState()
		}
	}

	approvedReviews := make(map[string]string)
	changesRequestedReviews := make(map[string]string)
	for login, state := range uniqueReviews {
		if state == approved {
			approvedReviews[login] = state
		} else if state == changesRequested {
			changesRequestedReviews[login] = state
		}
	}

	return approvedReviews, changesRequestedReviews, nil
}
