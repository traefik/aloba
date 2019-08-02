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

// RunGitHubAction Performs the GitHub Action
func RunGitHubAction(options *options.GitHubAction, gitHubToken string) error {
	if options.Debug {
		log.Println(options)
	}

	ctx := context.Background()
	client := gh.NewGitHubClient(ctx, gitHubToken)

	eventName := os.Getenv("GITHUB_EVENT_NAME")
	eventPath := os.Getenv("GITHUB_EVENT_PATH")

	switch eventName {
	case eventtype.Issues:
		event := &github.IssuesEvent{}
		err := readEvent(eventPath, event)
		if err != nil {
			return fmt.Errorf("unable to read the event file %q: %v", eventPath, err)
		}

		owner, repoName := getRepoInfo()
		if event.GetAction() == stateOpened {
			return onIssueOpened(ctx, client, event, owner, repoName, options.DryRun)
		}

	case eventtype.PullRequest:
		event := &github.PullRequestEvent{}
		err := readEvent(eventPath, event)
		if err != nil {
			return fmt.Errorf("unable to read the event file %q: %v", eventPath, err)
		}

		rulesPath := filepath.Join(os.Getenv("GITHUB_WORKSPACE"), ".github", "aloba-rules.toml")
		if _, err = os.Stat(rulesPath); os.IsNotExist(err) {
			return fmt.Errorf("unable to read the rules file %q: %v", rulesPath, err)
		}

		rc := &RulesConfiguration{}
		meta, err := toml.DecodeFile(rulesPath, rc)
		if err != nil {
			return err
		}

		if options.DryRun {
			log.Printf("Rules: %+v\n", meta)
		}

		owner, repoName := getRepoInfo()
		if event.GetAction() == stateOpened {
			return onPullRequestOpened(ctx, client, event, owner, repoName, rc, options.DryRun)
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
