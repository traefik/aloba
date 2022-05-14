package label

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/v44/github"
)

// Rule Area rule.
type Rule struct {
	Label string
	Regex string
}

// DetectAreas for a Pull Request.
func DetectAreas(ctx context.Context, client *github.Client, owner, repositoryName string, prNumber int, rules []Rule) ([]string, error) {
	areasMap := make(map[string]struct{})

	opt := &github.ListOptions{
		PerPage: 150,
	}

	for {
		cfs, resp, err := client.PullRequests.ListFiles(ctx, owner, repositoryName, prNumber, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list files from PR %d: %w", prNumber, err)
		}

		for _, cf := range cfs {
			for _, rule := range rules {
				if isRelatedTo(cf.GetFilename(), rule.Regex) && !strings.HasPrefix(cf.GetFilename(), "vendor/") {
					areasMap[rule.Label] = struct{}{}
				}
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var areas []string
	for area := range areasMap {
		areas = append(areas, area)
	}
	return areas, nil
}

func isRelatedTo(filename, exp string) bool {
	return regexp.MustCompile(exp).MatchString(filename)
}
