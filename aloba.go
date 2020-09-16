package main

import (
	"fmt"
	"log"
	"os"

	"github.com/containous/flaeg"
	"github.com/ogier/pflag"
	"github.com/traefik/aloba/cmd"
	"github.com/traefik/aloba/meta"
	"github.com/traefik/aloba/options"
)

func main() {
	rootCmd := createRootCommand()
	flag := flaeg.New(rootCmd, os.Args[1:])

	// Report
	reportCmd := createReportCommand()
	flag.AddCommand(reportCmd)

	// Label
	labelCmd := createLabelCommand()
	flag.AddCommand(labelCmd)

	// GitHubAction
	ghaCmd := createGitHubActionCommand()
	flag.AddCommand(ghaCmd)

	// Run command

	// version
	versionCmd := createVersionCommand()
	flag.AddCommand(versionCmd)

	// Print help when the command is running without any parameters.
	rootCmd.Run = func() error {
		return flaeg.LoadWithCommand(rootCmd, []string{"-h"}, nil, []*flaeg.Command{rootCmd, reportCmd, labelCmd, ghaCmd, versionCmd})
	}

	// Run command
	err := flag.Run()
	if err != nil && err != pflag.ErrHelp {
		log.Fatalf("Error: %v\n", err)
	}
}

func createVersionCommand() *flaeg.Command {
	return &flaeg.Command{
		Name:                  "version",
		Description:           "Display the version.",
		Config:                &options.Empty{},
		DefaultPointersConfig: &options.Empty{},
		Run: func() error {
			meta.DisplayVersion()
			return nil
		},
	}
}

func createRootCommand() *flaeg.Command {
	emptyConfig := &options.Empty{}

	rootCmd := &flaeg.Command{
		Name:                  "aloba",
		Description:           "Myrmica Aloba: Manage GitHub labels.",
		Config:                emptyConfig,
		DefaultPointersConfig: &options.Empty{},
	}

	return rootCmd
}

func createReportCommand() *flaeg.Command {
	reportOptions := &options.Report{
		Slack: &options.Slack{
			IconEmoji: ":captainpr:",
			BotName:   "CaptainPR",
		},
		DryRun:     true,
		ServerPort: 80,
	}

	reportCmd := &flaeg.Command{
		Name:                  "report",
		Description:           "Create a report and publish on Slack.",
		Config:                reportOptions,
		DefaultPointersConfig: &options.Report{Slack: &options.Slack{}, GitHub: &options.GitHub{}},
	}
	reportCmd.Run = func() error {
		if reportOptions.DryRun {
			log.Print("IMPORTANT: you are using the dry-run mode. Use `--dry-run=false` to disable this mode.")
		}

		if len(reportOptions.GitHub.Token) == 0 {
			reportOptions.GitHub.Token = os.Getenv("GITHUB_TOKEN")
		}

		if len(reportOptions.Slack.Token) == 0 {
			reportOptions.Slack.Token = os.Getenv("SLACK_TOKEN")
		}

		err := validateReportOptions(reportOptions)
		if err != nil {
			return err
		}

		return cmd.Report(reportOptions)
	}

	return reportCmd
}

func createLabelCommand() *flaeg.Command {
	labelOptions := &options.Label{
		RulesFilePath: "./rules.toml",
		DryRun:        true,
	}

	defaultPointerOptions := &options.Label{
		GitHub: &options.GitHub{},
		WebHook: &options.WebHook{
			Port: 80,
		},
	}

	labelCmd := &flaeg.Command{
		Name:                  "label",
		Description:           "Add labels and milestone on pull requests and issues.",
		Config:                labelOptions,
		DefaultPointersConfig: defaultPointerOptions,
	}

	labelCmd.Run = func() error {
		if labelOptions.DryRun {
			log.Print("IMPORTANT: you are using the dry-run mode. Use `--dry-run=false` to disable this mode.")
		}

		if len(labelOptions.GitHub.Token) == 0 {
			labelOptions.GitHub.Token = os.Getenv("GITHUB_TOKEN")
		}

		if labelOptions.WebHook != nil && len(labelOptions.WebHook.Secret) == 0 {
			labelOptions.WebHook.Secret = os.Getenv("WEBHOOK_SECRET")
		}

		err := validateLabelOptions(labelOptions)
		if err != nil {
			return err
		}

		return cmd.Label(labelOptions)
	}

	return labelCmd
}

func createGitHubActionCommand() *flaeg.Command {
	labelOptions := &options.GitHubAction{
		DryRun: true,
	}

	labelCmd := &flaeg.Command{
		Name:                  "action",
		Description:           "GitHub Action",
		Config:                labelOptions,
		DefaultPointersConfig: &options.GitHubAction{},
	}

	labelCmd.Run = func() error {
		if labelOptions.DryRun {
			log.Print("IMPORTANT: you are using the dry-run mode. Use `--dry-run=false` to disable this mode.")
		}

		ghToken := os.Getenv("GITHUB_TOKEN")
		err := required(ghToken, "GITHUB_TOKEN")
		if err != nil {
			return err
		}

		return cmd.RunGitHubAction(labelOptions, ghToken)
	}

	return labelCmd
}

func required(field, fieldName string) error {
	if len(field) == 0 {
		return fmt.Errorf("%s is mandatory", fieldName)
	}
	return nil
}

func validateReportOptions(reportOptions *options.Report) error {
	err := required(reportOptions.GitHub.Token, "github.token")
	if err != nil {
		return err
	}
	err = required(reportOptions.GitHub.Owner, "github.owner")
	if err != nil {
		return err
	}
	err = required(reportOptions.GitHub.RepositoryName, "github.repo-name")
	if err != nil {
		return err
	}
	err = required(reportOptions.Slack.Token, "slack.token")
	if err != nil {
		return err
	}

	if !reportOptions.ServerMode {
		errSlackChan := required(reportOptions.Slack.ChannelID, "slack.channel")
		if errSlackChan != nil {
			return errSlackChan
		}
	}
	return nil
}

func validateLabelOptions(labelOptions *options.Label) error {
	err := required(labelOptions.GitHub.Token, "github.token")
	if err != nil {
		return err
	}
	err = required(labelOptions.GitHub.Owner, "github.owner")
	if err != nil {
		return err
	}
	err = required(labelOptions.GitHub.RepositoryName, "github.repo-name")
	if err != nil {
		return err
	}
	return required(labelOptions.RulesFilePath, "rules-path")
}
