package cmd

import (
	"context"
	"log"
	"net/url"
	"strings"

	"github.com/containous/aloba/internal/gh"
	"github.com/containous/aloba/label"
	"github.com/containous/aloba/options"
	"github.com/google/go-github/v27/github"
	ghw "github.com/ldez/ghwebhook/v2"
	"github.com/ldez/ghwebhook/v2/eventtype"
)

func (l *Labeler) runWebHook(ctx context.Context, owner string, repositoryName string, rc *RulesConfiguration, opts *options.WebHook, dryRun bool) error {
	handlers := ghw.NewEventHandlers().
		OnPullRequest(l.onPullRequest(ctx, owner, repositoryName, rc, dryRun)).
		OnPullRequestReview(l.onPullRequestReview(ctx, owner, repositoryName, dryRun)).
		OnIssues(l.onIssue(ctx, owner, repositoryName, dryRun))

	hook := ghw.NewWebHook(handlers,
		ghw.WithPort(opts.Port),
		ghw.WithSecret(opts.Secret),
		ghw.WithEventTypes(eventtype.PullRequest, eventtype.PullRequestReview, eventtype.Issues))
	return hook.ListenAndServe()
}

func (l *Labeler) onIssue(ctx context.Context, owner string, repositoryName string, dryRun bool) func(*url.URL, *github.WebHookPayload, *github.IssuesEvent) {
	return func(_ *url.URL, _ *github.WebHookPayload, event *github.IssuesEvent) {
		if event.GetAction() == stateOpened {
			go func(event *github.IssuesEvent) {
				err := l.onIssueOpened(ctx, event, owner, repositoryName, dryRun)
				if err != nil {
					log.Println(err)
				}
			}(event)
		}
	}
}

func (l *Labeler) onPullRequest(ctx context.Context, owner string, repositoryName string, rc *RulesConfiguration, dryRun bool) func(*url.URL, *github.WebHookPayload, *github.PullRequestEvent) {
	return func(_ *url.URL, _ *github.WebHookPayload, event *github.PullRequestEvent) {
		if event.GetAction() == stateOpened {
			go func(event *github.PullRequestEvent) {
				err := l.onPullRequestOpened(ctx, event, owner, repositoryName, rc, dryRun)
				if err != nil {
					log.Println(err)
				}
			}(event)
		}
	}
}

func (l *Labeler) onPullRequestReview(ctx context.Context, owner string, repositoryName string, dryRun bool) func(*url.URL, *github.WebHookPayload, *github.PullRequestReviewEvent) {
	return func(_ *url.URL, _ *github.WebHookPayload, event *github.PullRequestReviewEvent) {
		if event.GetAction() == "submitted" {
			if strings.EqualFold(event.Review.GetState(), gh.ChangesRequested) {
				go func(event *github.PullRequestReviewEvent) {
					issue, _, err := l.client.Issues.Get(ctx, owner, repositoryName, event.PullRequest.GetNumber())
					if err != nil {
						log.Println(err)
						return
					}

					if label.HasLabel(issue.Labels, label.ContributorWaitingForCorrections) {
						return
					}

					if dryRun {
						log.Printf("#%d: Add %s\n", issue.GetNumber(), label.ContributorWaitingForCorrections)
					} else {
						_, _, err = l.client.Issues.AddLabelsToIssue(ctx, owner, repositoryName, issue.GetNumber(), []string{label.ContributorWaitingForCorrections})
						if err != nil {
							log.Println(err)
							return
						}
					}
				}(event)
				return
			}

			if strings.EqualFold(event.Review.GetState(), gh.Approved) {
				go func(event *github.PullRequestReviewEvent) {
					issue, _, err := l.client.Issues.Get(ctx, owner, repositoryName, event.PullRequest.GetNumber())
					if err != nil {
						log.Println(err)
						return
					}

					err = l.removeLabel(ctx, owner, repositoryName, issue, label.ContributorWaitingForCorrections, dryRun)
					if err != nil {
						log.Println(err)
						return
					}
				}(event)
				return
			}
		}
	}
}

func (l *Labeler) removeLabel(ctx context.Context, owner string, repositoryName string, issue *github.Issue, labelName string, dryRun bool) error {
	if label.HasLabel(issue.Labels, labelName) {
		if dryRun {
			log.Printf("#%d: Remove %s\n", issue.GetNumber(), labelName)
		} else {
			_, err := l.client.Issues.RemoveLabelForIssue(ctx, owner, repositoryName, issue.GetNumber(), labelName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
