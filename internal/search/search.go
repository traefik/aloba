package search

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/go-github/github"
)

type byCreated []github.Issue

func (a byCreated) Len() int      { return len(a) }
func (a byCreated) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byCreated) Less(i, j int) bool {
	return a[i].GetCreatedAt().Before(a[j].GetCreatedAt())
}

// FindOpenPR search and find open Pull Requests.
func FindOpenPR(ctx context.Context, client *github.Client, owner string, repositoryName string, parameters ...Parameter) ([]github.Issue, error) {

	var filter string
	for _, param := range parameters {
		if param != nil {
			filter += param()
		}
	}

	query := fmt.Sprintf("repo:%s/%s type:pr state:open %s", owner, repositoryName, filter)

	options := &github.SearchOptions{
		Sort:        "created",
		Order:       "desc",
		ListOptions: github.ListOptions{PerPage: 25},
	}

	issues, err := findIssues(ctx, client, query, options)
	if err != nil {
		return nil, err
	}
	sort.Sort(byCreated(issues))

	return issues, nil
}

func findIssues(ctx context.Context, client *github.Client, query string, searchOptions *github.SearchOptions) ([]github.Issue, error) {
	var allIssues []github.Issue
	for {
		issuesSearchResult, resp, err := client.Search.Issues(ctx, query, searchOptions)
		if err != nil {
			return nil, err
		}
		allIssues = append(allIssues, issuesSearchResult.Issues...)
		if resp.NextPage == 0 {
			break
		}
		searchOptions.Page = resp.NextPage
	}
	return allIssues, nil
}
