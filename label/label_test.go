package label

import (
	"testing"

	"github.com/google/go-github/github"
)

func TestHasLabel(t *testing.T) {
	basicLabels := []github.Label{
		{
			Name: github.String("foo"),
		},
		{
			Name: github.String("bar"),
		},
		{
			Name: github.String("fii"),
		},
		{
			Name: github.String("bir"),
		},
	}

	testCases := []struct {
		name           string
		labels         []github.Label
		label          string
		expectedResult bool
	}{
		{
			name:           "label exists",
			labels:         basicLabels,
			label:          "foo",
			expectedResult: true,
		},
		{
			name:           "label not exists",
			labels:         basicLabels,
			label:          "fuu",
			expectedResult: false,
		},
		{
			name:           "empty labels",
			labels:         []github.Label{},
			label:          "fuu",
			expectedResult: false,
		},
		{
			name:           "empty label",
			labels:         basicLabels,
			label:          "",
			expectedResult: false,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := HasLabel(test.labels, test.label)

			if result != test.expectedResult {
				t.Errorf("Got %t, want %t", result, test.expectedResult)
			}
		})
	}
}

func TestHasStatus(t *testing.T) {
	testsCases := []struct {
		name           string
		labels         []github.Label
		expectedResult bool
	}{
		{
			name: "with status label",
			labels: []github.Label{
				{
					Name: github.String("foo"),
				},
				{
					Name: github.String("bar"),
				},
				{
					Name: github.String("status/test"),
				},
			},
			expectedResult: true,
		},
		{
			name: "without status label",
			labels: []github.Label{
				{
					Name: github.String("foo"),
				},
				{
					Name: github.String("bar"),
				},
				{
					Name: github.String("fii"),
				},
			},
			expectedResult: false,
		},
		{
			name: "with WIP label",
			labels: []github.Label{
				{
					Name: github.String("foo"),
				},
				{
					Name: github.String("bar"),
				},
				{
					Name: github.String("WIP"),
				},
			},
			expectedResult: true,
		},
	}

	for _, test := range testsCases {
		test := test
		t.Run(test.name, func(t *testing.T) {

			result := HasStatus(test.labels)

			if result != test.expectedResult {
				t.Errorf("Got %t, want %t", result, test.expectedResult)
			}
		})
	}
}
