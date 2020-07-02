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

func (l *Labeler) runWebHook(ctx context.Context, rc *RulesConfiguration, opts *options.WebHook) error {
	handlers := ghw.NewEventHandlers().
		OnPullRequest(l.onPullRequest(ctx, rc)).
		OnPullRequestReview(l.onPullRequestReview(ctx)).
		OnIssues(l.onIssue(ctx))

	hook := ghw.NewWebHook(handlers,
		ghw.WithPort(opts.Port),
		ghw.WithSecret(opts.Secret),
		ghw.WithEventTypes(eventtype.PullRequest, eventtype.PullRequestReview, eventtype.Issues))

	return hook.ListenAndServe()
}

func (l *Labeler) onIssue(ctx context.Context) func(*url.URL, *github.WebHookPayload, *github.IssuesEvent) {
	return func(_ *url.URL, _ *github.WebHookPayload, event *github.IssuesEvent) {
		if event.GetAction() == stateOpened {
			go func(event *github.IssuesEvent) {
				err := l.onIssueOpened(ctx, event)
				if err != nil {
					log.Println(err)
				}
			}(event)
		}
	}
}

func (l *Labeler) onPullRequest(ctx context.Context, rc *RulesConfiguration) func(*url.URL, *github.WebHookPayload, *github.PullRequestEvent) {
	return func(_ *url.URL, _ *github.WebHookPayload, event *github.PullRequestEvent) {
		if event.GetAction() == stateOpened {
			go func(event *github.PullRequestEvent) {
				err := l.onPullRequestOpened(ctx, event, rc)
				if err != nil {
					log.Println(err)
				}
			}(event)
		}
	}
}

func (l *Labeler) onPullRequestReview(ctx context.Context) func(*url.URL, *github.WebHookPayload, *github.PullRequestReviewEvent) {
	return func(_ *url.URL, _ *github.WebHookPayload, event *github.PullRequestReviewEvent) {
		if event.GetAction() == "submitted" {
			if strings.EqualFold(event.Review.GetState(), gh.ChangesRequested) {
				go func(event *github.PullRequestReviewEvent) {
					issue, _, err := l.client.Issues.Get(ctx, l.owner, l.repoName, event.PullRequest.GetNumber())
					if err != nil {
						log.Println(err)
						return
					}

					if label.HasLabel(issue.Labels, label.ContributorWaitingForCorrections) {
						return
					}

					if l.DryRun {
						log.Printf("#%d: Add %s\n", issue.GetNumber(), label.ContributorWaitingForCorrections)
					} else {
						_, _, err = l.client.Issues.AddLabelsToIssue(ctx, l.owner, l.repoName, issue.GetNumber(), []string{label.ContributorWaitingForCorrections})
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
					issue, _, err := l.client.Issues.Get(ctx, l.owner, l.repoName, event.PullRequest.GetNumber())
					if err != nil {
						log.Println(err)
						return
					}

					err = l.removeLabel(ctx, issue, label.ContributorWaitingForCorrections)
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

func (l *Labeler) removeLabel(ctx context.Context, issue *github.Issue, labelName string) error {
	if label.HasLabel(issue.Labels, labelName) {
		if l.DryRun {
			log.Printf("#%d: Remove %s\n", issue.GetNumber(), labelName)
		} else {
			_, err := l.client.Issues.RemoveLabelForIssue(ctx, l.owner, l.repoName, issue.GetNumber(), labelName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
