package label

import (
	"strings"

	"github.com/google/go-github/v47/github"
)

// status labels.
const (
	WIP                     = "WIP"
	StatusPrefix            = "status/"
	StatusNeedsTriage       = StatusPrefix + "0-needs-triage"
	StatusNeedsDesignReview = StatusPrefix + "1-needs-design-review"
	StatusNeedsReview       = StatusPrefix + "2-needs-review"
	StatusNeedsMerge        = StatusPrefix + "3-needs-merge"
)

// contributor labels.
const (
	ContributorPrefix                  = "contributor/"
	ContributorWaitingForCorrections   = ContributorPrefix + "waiting-for-corrections"
	ContributorNeedMoreInformation     = ContributorPrefix + "need-more-information"
	ContributorWaitingForDocumentation = ContributorPrefix + "waiting-for-documentation"
	ContributorWaitingForFeedback      = ContributorPrefix + "waiting-for-feedback"
)

// size labels.
const (
	SizeLabelPrefix = "size/"
	Small           = "S"
	Medium          = "M"
	Large           = "L"
)

// HasStatus check if issue labels contains a status label.
func HasStatus(issueLabels []*github.Label) bool {
	for _, lbl := range issueLabels {
		if strings.HasPrefix(lbl.GetName(), StatusPrefix) || strings.HasSuffix(lbl.GetName(), WIP) {
			return true
		}
	}
	return false
}

// HasLabel check if issue labels  contains a specific label.
func HasLabel(issueLabels []*github.Label, label string) bool {
	for _, lbl := range issueLabels {
		if lbl.GetName() == label {
			return true
		}
	}
	return false
}
