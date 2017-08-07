package label

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/go-github/github"
)

type Rule struct {
	Label string
	Regex string
}

func DetectAreas(client *github.Client, ctx context.Context, owner string, repositoryName string, prNumber int, rules []Rule) ([]string, error) {

	areasMap := make(map[string]struct{})

	opt := &github.ListOptions{
		PerPage: 150,
	}

	for {
		cfs, resp, err := client.PullRequests.ListFiles(ctx, owner, repositoryName, prNumber, opt)
		if err != nil {
			return nil, err
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

	areas := []string{}

	for area := range areasMap {
		areas = append(areas, area)
	}

	return areas, nil
}

func isRelatedTo(filename string, exp string) bool {
	return regexp.MustCompile(exp).MatchString(filename)
}
