package milestone

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"

	"github.com/google/go-github/v27/github"
)

// MetaMilestone represent a milestone.
type MetaMilestone struct {
	weight int64
	name   string
	ID     int
}

type byWeight []*MetaMilestone

func (m byWeight) Len() int {
	return len(m)
}

func (m byWeight) Less(i, j int) bool {
	return m[i].weight > m[j].weight
}

func (m byWeight) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

var expMilestone = regexp.MustCompile(`(\d+)\.(\d+)(?:\.(\d+))?`)

// Detect the possible milestone of a PR.
func Detect(ctx context.Context, client *github.Client, owner, repoName string, pr *github.PullRequest) (*MetaMilestone, error) {
	opt := &github.MilestoneListOptions{}
	opt.State = "all"
	opt.PerPage = 10

	stones, _, err := client.Issues.ListMilestones(ctx, owner, repoName, opt)
	if err != nil {
		return nil, err
	}

	metas := makeMetas(stones)

	// if the base branch is a possible version
	if meta := find(metas, pr.Base.GetRef()); meta != nil {
		return meta, nil
	}

	return nil, nil
}

func makeMetas(stones []*github.Milestone) []*MetaMilestone {
	var metas []*MetaMilestone

	for _, ml := range stones {
		weight, errWeight := weightCalculator(ml)
		if errWeight != nil {
			log.Println(errWeight)
			continue
		}
		meta := &MetaMilestone{
			name:   ml.GetTitle(),
			ID:     ml.GetNumber(),
			weight: weight,
		}
		metas = append(metas, meta)
	}

	sort.Sort(byWeight(metas))

	return metas
}

func find(metas []*MetaMilestone, baseRef string) *MetaMilestone {
	for _, meta := range metas {
		if baseRef == meta.name || baseRef == "v"+meta.name {
			return meta
		}
	}
	return nil
}

func weightCalculator(ml *github.Milestone) (int64, error) {
	parts := expMilestone.FindStringSubmatch(ml.GetTitle())

	if len(parts) != 4 {
		return 0, fmt.Errorf("invalid milestone title %s", ml.GetTitle())
	}

	part1, err := parseNumber(parts[3], 100)
	if err != nil {
		return 0, err
	}
	part2, err := parseNumber(parts[2], 100000)
	if err != nil {
		return 0, err
	}
	part3, err := parseNumber(parts[1], 100000000)
	if err != nil {
		return 0, err
	}

	return part1 + part2 + part3, nil
}

func parseNumber(raw string, multiplier int64) (int64, error) {
	if raw == "" {
		return 0, nil
	}

	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return value * multiplier, nil
}
