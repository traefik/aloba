package label

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
)

// Limits a set of Pull request limits by size.
type Limits struct {
	Small  Limit
	Medium Limit
}

// Limit a set of Pull request limits.
type Limit struct {
	SumLimit   int
	DiffLimit  int
	FilesLimit int
}

// Changes represent the changes of a Pull Request.
type Changes struct {
	Number       int
	AdditionSum  int
	DeletionSum  int
	ChangedFiles int
}

// GetCurrentSize get the size of a Pull Request.
func GetCurrentSize(issueLabels []github.Label) string {
	for _, lbl := range issueLabels {
		if strings.HasPrefix(lbl.GetName(), SizeLabelPrefix) {
			return lbl.GetName()
		}
	}
	return ""
}

// GetSizeLabel evaluate PR size (exclude vendor files)
func GetSizeLabel(ctx context.Context, client *github.Client, owner string, repositoryName string, prNumber int, limits Limits) (string, error) {
	changes, err := calculateChanges(ctx, client, owner, repositoryName, prNumber)
	if err != nil {
		return "", err
	}

	return SizeLabelPrefix + getSizeLabel(changes, limits), nil
}

// calculateChanges count changes (exclude vendor files)
func calculateChanges(ctx context.Context, client *github.Client, owner string, repositoryName string, prNumber int) (*Changes, error) {

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
				changes.ChangedFiles++
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
