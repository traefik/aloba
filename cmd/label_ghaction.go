package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/containous/aloba/internal/gh"
	"github.com/containous/aloba/options"
	"github.com/google/go-github/v27/github"
	"github.com/ldez/ghwebhook/v2/eventtype"
)

// RunGitHubAction Performs the GitHub Action.
func RunGitHubAction(options *options.GitHubAction, gitHubToken string) error {
	if options.Debug {
		log.Println(options)
	}

	ctx := context.Background()
	client := gh.NewGitHubClient(ctx, gitHubToken)

	owner, repoName := getRepoInfo()

	labeler := NewLabeler(client, owner, repoName)
	labeler.DryRun = options.DryRun

	eventName := os.Getenv("GITHUB_EVENT_NAME")
	eventPath := os.Getenv("GITHUB_EVENT_PATH")

	switch eventName {
	case eventtype.Issues:
		event := &github.IssuesEvent{}
		err := readEvent(eventPath, event)
		if err != nil {
			return fmt.Errorf("unable to read the event file %q: %w", eventPath, err)
		}

		if event.GetAction() == stateOpened {
			return labeler.onIssueOpened(ctx, event)
		}

	case eventtype.PullRequest:
		event := &github.PullRequestEvent{}
		err := readEvent(eventPath, event)
		if err != nil {
			return fmt.Errorf("unable to read the event file %q: %w", eventPath, err)
		}

		rulesPath := filepath.Join(os.Getenv("GITHUB_WORKSPACE"), ".github", "aloba-rules.toml")
		if _, err = os.Stat(rulesPath); os.IsNotExist(err) {
			return fmt.Errorf("unable to read the rules file %q: %w", rulesPath, err)
		}

		rc := &RulesConfiguration{}
		meta, err := toml.DecodeFile(rulesPath, rc)
		if err != nil {
			return err
		}

		if options.DryRun {
			log.Printf("Rules: %+v\n", meta)
		}

		if event.GetAction() == stateOpened {
			return labeler.onPullRequestOpened(ctx, event, rc)
		}

	default:
		return fmt.Errorf("unsupported event type: %s", eventName)
	}

	return nil
}

func readEvent(eventPath string, event interface{}) error {
	content, err := ioutil.ReadFile(eventPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, event)
}

func getRepoInfo() (string, string) {
	githubRepository := os.Getenv("GITHUB_REPOSITORY")

	parts := strings.SplitN(githubRepository, "/", 2)

	return parts[0], parts[1]
}
