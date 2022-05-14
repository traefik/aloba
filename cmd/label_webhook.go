package cmd

import (
	"context"
	"net/url"
	"strings"

	"github.com/google/go-github/v44/github"
	ghw "github.com/ldez/ghwebhook/v3"
	"github.com/ldez/ghwebhook/v3/eventtype"
	"github.com/rs/zerolog/log"
	"github.com/traefik/aloba/internal/gh"
	"github.com/traefik/aloba/label"
	"github.com/traefik/aloba/options"
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

func (l *Labeler) onIssue(ctx context.Context) func(*url.URL, *github.IssuesEvent) {
	return func(_ *url.URL, event *github.IssuesEvent) {
		if event.GetAction() == stateOpened {
			go func(event *github.IssuesEvent) {
				err := l.onIssueOpened(ctx, event)
				if err != nil {
					log.Error().Int("issue", event.Issue.GetNumber()).Err(err).Msg("error")
				}
			}(event)
		}
	}
}

func (l *Labeler) onPullRequest(ctx context.Context, rc *RulesConfiguration) func(*url.URL, *github.PullRequestEvent) {
	return func(_ *url.URL, event *github.PullRequestEvent) {
		if event.GetAction() == stateOpened {
			go func(event *github.PullRequestEvent) {
				err := l.onPullRequestOpened(ctx, event, rc)
				if err != nil {
					log.Error().Int("pr", event.PullRequest.GetNumber()).Err(err).Msg("error")
				}
			}(event)
		}
	}
}

func (l *Labeler) onPullRequestReview(ctx context.Context) func(*url.URL, *github.PullRequestReviewEvent) {
	return func(_ *url.URL, event *github.PullRequestReviewEvent) {
		if event.GetAction() == "submitted" {
			if strings.EqualFold(event.Review.GetState(), gh.ChangesRequested) {
				go func(event *github.PullRequestReviewEvent) {
					issue, _, err := l.client.Issues.Get(ctx, l.owner, l.repoName, event.PullRequest.GetNumber())
					if err != nil {
						log.Error().Int("pr", event.PullRequest.GetNumber()).Err(err).Msg("Failed to get PR.")
						return
					}

					if label.HasLabel(issue.Labels, label.ContributorWaitingForCorrections) {
						return
					}

					if l.DryRun {
						log.Debug().Int("pr", issue.GetNumber()).Msgf("Add %s", label.ContributorWaitingForCorrections)
					} else {
						_, _, err = l.client.Issues.AddLabelsToIssue(ctx, l.owner, l.repoName, issue.GetNumber(), []string{label.ContributorWaitingForCorrections})
						if err != nil {
							log.Error().Err(err).Int("pr", event.PullRequest.GetNumber()).Msg("Failed to add label on PR.")
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
						log.Error().Int("pr", event.PullRequest.GetNumber()).Err(err).Msg("Failed to get PR.")
						return
					}

					err = l.removeLabel(ctx, issue, label.ContributorWaitingForCorrections)
					if err != nil {
						log.Error().Int("pr", event.PullRequest.GetNumber()).Err(err).Msg("Failed to remove label on PR.")
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
			log.Debug().Int("pr", issue.GetNumber()).Msgf("Remove %s", labelName)
		} else {
			_, err := l.client.Issues.RemoveLabelForIssue(ctx, l.owner, l.repoName, issue.GetNumber(), labelName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
