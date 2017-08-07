package label

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
)

type Limits struct {
	Small  Limit
	Medium Limit
}

type Limit struct {
	SumLimit   int
	DiffLimit  int
	FilesLimit int
}

type Changes struct {
	Number       int
	AdditionSum  int
	DeletionSum  int
	ChangedFiles int
}

// GetCurrentSize
func GetCurrentSize(issueLabels []github.Label) string {
	for _, lbl := range issueLabels {
		if strings.HasPrefix(lbl.GetName(), SizeLabelPrefix) {
			return lbl.GetName()
		}
	}
	return ""
}

// GetSizeLabel evaluate PR size (exclude vendor files)
func GetSizeLabel(client *github.Client, ctx context.Context, owner string, repositoryName string, prNumber int, limits Limits) (string, error) {
	changes, err := calculateChanges(client, ctx, owner, repositoryName, prNumber)
	if err != nil {
		return "", err
	}

	return SizeLabelPrefix + getSizeLabel(changes, limits), nil
}

// calculateChanges count changes (exclude vendor files)
func calculateChanges(client *github.Client, ctx context.Context, owner string, repositoryName string, prNumber int) (*Changes, error) {

	changes := &Changes{
		Number: prNumber,
	}

	opt := &github.ListOptions{
		PerPage: 150,
	}

	for {
		cfs, resp, err := client.PullRequests.ListFiles(ctx, owner, repositoryName, prNumber, opt)
		if err != nil {
			return nil, err
		}

		for _, cf := range cfs {

			if !strings.HasPrefix(cf.GetFilename(), "vendor/") && cf.GetFilename() != "glide.lock" && cf.GetFilename() != "glide.yml" {
				changes.ChangedFiles += 1
				changes.AdditionSum += cf.GetAdditions()
				changes.DeletionSum += cf.GetDeletions()
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return changes, nil
}

func getSizeLabel(changes *Changes, limits Limits) string {
	sum := changes.AdditionSum + changes.DeletionSum

	diff := changes.AdditionSum - changes.DeletionSum
	if diff < 0 {
		diff = -diff
	}

	if sum < limits.Small.SumLimit && diff < limits.Small.DiffLimit && changes.ChangedFiles < limits.Small.FilesLimit {
		return Small
	} else if sum < limits.Medium.SumLimit && diff < limits.Medium.DiffLimit && changes.ChangedFiles < limits.Medium.FilesLimit {
		return Medium
	} else {
		return Large
	}
}
