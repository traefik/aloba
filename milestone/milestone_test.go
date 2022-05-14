package milestone

import (
	"reflect"
	"testing"

	"github.com/google/go-github/v44/github"
)

func Test_find(t *testing.T) {
	testCases := []struct {
		name           string
		metas          []*MetaMilestone
		baseRef        string
		expectNotFound bool
		expectedMeta   *MetaMilestone
	}{
		{
			name: "branch with a 'v'",
			metas: []*MetaMilestone{
				{name: "1.2", ID: 1, weight: 20},
				{name: "1.3", ID: 2, weight: 30},
				{name: "1.1", ID: 3, weight: 20},
			},
			baseRef:      "v1.1",
			expectedMeta: &MetaMilestone{name: "1.1", ID: 3, weight: 20},
		},
		{
			name: "not existing milestone",
			metas: []*MetaMilestone{
				{name: "1.2", ID: 1, weight: 20},
				{name: "1.3", ID: 2, weight: 30},
				{name: "1.1", ID: 3, weight: 20},
			},
			baseRef:        "v1.0",
			expectNotFound: true,
		},
		{
			name: "branch without a 'v'",
			metas: []*MetaMilestone{
				{name: "1.2", ID: 1, weight: 20},
				{name: "1.3", ID: 2, weight: 30},
				{name: "1.1", ID: 3, weight: 20},
			},
			baseRef:      "1.1",
			expectedMeta: &MetaMilestone{name: "1.1", ID: 3, weight: 20},
		},
		{
			name: "branch master",
			metas: []*MetaMilestone{
				{name: "1.2", ID: 1, weight: 20},
				{name: "1.3", ID: 2, weight: 30},
				{name: "1.1", ID: 3, weight: 20},
			},
			baseRef:        "master",
			expectNotFound: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			meta := find(test.metas, test.baseRef)

			if !test.expectNotFound && meta == nil {
				t.Fatalf("Got %v, want no meta.", meta)
			}
			if test.expectNotFound && meta != nil {
				t.Fatalf("Got no meta, want %v.", test.expectedMeta)
			}

			if !reflect.DeepEqual(meta, test.expectedMeta) {
				t.Errorf("Got %v, want %v.", meta, test.expectedMeta)
			}
		})
	}
}

func Test_weightCalculator(t *testing.T) {
	testCases := []struct {
		name           string
		stone          *github.Milestone
		exceptedWeight int64
		exceptError    bool
	}{
		{
			name: "with 2 digits",
			stone: &github.Milestone{
				Title: github.String("1.4"),
			},
			exceptedWeight: 100400000,
		},
		{
			name: "with 3 digits",
			stone: &github.Milestone{
				Title: github.String("1.3.2"),
			},
			exceptedWeight: 100300200,
		},
		{
			name: "with 1 digits",
			stone: &github.Milestone{
				Title: github.String("1"),
			},
			exceptError: true,
		},
		{
			name: "with letters",
			stone: &github.Milestone{
				Title: github.String("foo"),
			},
			exceptError: true,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			weight, err := weightCalculator(test.stone)
			if test.exceptError && err == nil {
				t.Fatalf("Expect an error.")
			}
			if !test.exceptError && err != nil {
				t.Fatalf("Got %v, want no error.", err)
			}
			if weight != test.exceptedWeight {
				t.Errorf("Got %d, want %d.", weight, test.exceptedWeight)
			}
		})
	}
}
