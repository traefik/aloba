package search

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/go-github/v27/github"
)

type byCreated []github.Issue

func (a byCreated) Len() int      { return len(a) }
func (a byCreated) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byCreated) Less(i, j int) bool {
	return a[i].GetCreatedAt().Before(a[j].GetCreatedAt())
}

// FindOpenPR search and find open Pull Requests.
func FindOpenPR(ctx context.Context, client *github.Client, owner, repositoryName string, parameters ...Parameter) ([]github.Issue, error) {
	query := createQuery(owner, repositoryName, parameters)

	searchOptions := &github.SearchOptions{
		Sort:        "created",
		Order:       "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var issues []github.Issue
	for {
		issuesSearchResult, resp, err := client.Search.Issues(ctx, query, searchOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to search PR on %s/%s: %w", owner, repositoryName, err)
		}

		issues = append(issues, issuesSearchResult.Issues...)
		if resp.NextPage == 0 {
			break
		}

		searchOptions.Page = resp.NextPage
	}

	sort.Sort(byCreated(issues))

	return issues, nil
}

func createQuery(owner, repositoryName string, parameters []Parameter) string {
	var filter string
	for _, param := range parameters {
		if param != nil {
			filter += param()
		}
	}

	return fmt.Sprintf("repo:%s/%s type:pr state:open %s", owner, repositoryName, filter)
}
