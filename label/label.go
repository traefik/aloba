package label

import (
	"strings"

	"github.com/google/go-github/github"
)

const (
	WIP                     = "WIP"
	StatusPrefix            = "status/"
	StatusNeedsTriage       = StatusPrefix + "0-needs-triage"
	StatusNeedsDesignReview = StatusPrefix + "1-needs-design-review"
	StatusNeedsReview       = StatusPrefix + "2-needs-review"
	StatusNeedsMerge        = StatusPrefix + "3-needs-merge"
)

const (
	ContributorPrefix                  = "contributor/"
	ContributorWaitingForCorrections   = ContributorPrefix + "waiting-for-corrections"
	ContributorNeedMoreInformation     = ContributorPrefix + "need-more-information"
	ContributorWaitingForDocumentation = ContributorPrefix + "waiting-for-documentation"
	ContributorWaitingForFeedback      = ContributorPrefix + "waiting-for-feedback"
)

const (
	SizeLabelPrefix = "size/"
	Small           = "S"
	Medium          = "M"
	Large           = "L"
)

func HasStatus(issueLabels []github.Label) bool {
	for _, lbl := range issueLabels {
		if strings.HasPrefix(lbl.GetName(), StatusPrefix) || strings.HasSuffix(lbl.GetName(), WIP) {
			return true
		}
	}
	return false
}

func ExistsLabel(issueLabels []github.Label, label string) bool {
	for _, lbl := range issueLabels {
		if lbl.GetName() == label {
			return true
		}
	}
	return false
}
