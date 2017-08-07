package label

import (
	"strings"

	"github.com/google/go-github/github"
)

const (
	ContributorWaitingForCorrections = "contributor/waiting-for-corrections"
	WIP                              = "WIP"
	StatusPrefix                     = "status/"
	StatusNeedsTriage                = "status/0-needs-triage"
	StatusNeedsDesignReview          = "status/1-needs-design-review"
	StatusNeedsReview                = "status/2-needs-review"
	StatusNeedsMerge                 = "status/3-needs-merge"
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
