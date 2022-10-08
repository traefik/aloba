package gh

import (
	"context"
	"fmt"

	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"
)

// Review status.
const (
	Approved         = "APPROVED"
	ChangesRequested = "CHANGES_REQUESTED"
	Commented        = "COMMENTED"
)

// NewGitHubClient create a new GitHub client.
func NewGitHubClient(ctx context.Context, token string) *github.Client {
	if token == "" {
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return github.NewClient(oauth2.NewClient(ctx, ts))
}

// GetReviewStatus get reviews status of a Pull Request.
func GetReviewStatus(ctx context.Context, client *github.Client, owner, repositoryName string, members []*github.User, prNumber int) (map[string]string, map[string]string, error) {
	opts := &github.ListOptions{
		PerPage: 100,
	}

	uniqueReviews := make(map[string]string)

	for {
		reviews, resp, err := client.PullRequests.ListReviews(ctx, owner, repositoryName, prNumber, opts)
		if err != nil {
			return nil, nil, err
		}

		for _, review := range reviews {
			if review.GetState() != Commented && isTeamMember(members, review.User.GetLogin()) {
				uniqueReviews[review.User.GetLogin()] = review.GetState()
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	approvedReviews := make(map[string]string)
	changesRequestedReviews := make(map[string]string)
	for login, state := range uniqueReviews {
		if state == Approved {
			approvedReviews[login] = state
		} else if state == ChangesRequested {
			changesRequestedReviews[login] = state
		}
	}

	return approvedReviews, changesRequestedReviews, nil
}

// GetTeamMembers get members of a team.
func GetTeamMembers(ctx context.Context, client *github.Client, owner, teamName string) ([]*github.User, error) {
	team, err := getTeamByName(ctx, client, owner, teamName)
	if err != nil {
		return nil, err
	}

	orgTeamMemberOpts := &github.TeamListTeamMembersOptions{}
	orgTeamMemberOpts.PerPage = 100

	members, _, err := client.Teams.ListTeamMembersBySlug(ctx, owner, team.GetSlug(), orgTeamMemberOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to get team %s on %s: %w", teamName, owner, err)
	}

	return members, nil
}

func getTeamByName(ctx context.Context, client *github.Client, owner, teamName string) (*github.Team, error) {
	teams, _, err := client.Teams.ListTeams(ctx, owner, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams on %s: %w", owner, err)
	}

	for _, team := range teams {
		if team.GetName() == teamName {
			return team, nil
		}
	}
	return nil, fmt.Errorf("team %q not found", teamName)
}

func isTeamMember(members []*github.User, login string) bool {
	for _, member := range members {
		if member.GetLogin() == login {
			return true
		}
	}
	return false
}
