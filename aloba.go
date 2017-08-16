package main

import (
	"log"
	"os"

	"github.com/containous/aloba/cmd"
	"github.com/containous/aloba/options"
	"github.com/containous/flaeg"
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

	// Run command
	flag.Run()
}

func createRootCommand() *flaeg.Command {

	emptyConfig := &options.Empty{}

	rootCmd := &flaeg.Command{
		Name:                  "aloba",
		Description:           "Myrmica Aloba: Manage GitHub labels.",
		Config:                emptyConfig,
		DefaultPointersConfig: &options.Empty{},
		Run: func() error {
			// No op
			return nil
		},
	}

	return rootCmd
}

func createReportCommand() *flaeg.Command {

	reportOptions := &options.Report{
		Slack: &options.Slack{
			IconEmoji: ":captainpr:",
			BotName:   "CaptainPR",
		},
		DryRun: true,
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

		required(reportOptions.GitHub.Token, "github.token")
		required(reportOptions.GitHub.Owner, "github.owner")
		required(reportOptions.GitHub.RepositoryName, "github.repo-name")
		// FIXME
		required(reportOptions.Slack.Token, "slack.token")
		required(reportOptions.Slack.ChannelID, "slack.channel")

		err := cmd.Report(reportOptions)
		if err != nil {
			log.Println(err)
		}
		return nil
	}

	return reportCmd
}

func createLabelCommand() *flaeg.Command {

	labelOptions := &options.Label{
		RulesFilePath: "./rules.toml",
		DryRun:        true,
	}

	labelCmd := &flaeg.Command{
		Name:                  "label",
		Description:           "Add labels to Pull Request",
		Config:                labelOptions,
		DefaultPointersConfig: &options.Label{GitHub: &options.GitHub{}},
	}
	labelCmd.Run = func() error {
		if labelOptions.DryRun {
			log.Print("IMPORTANT: you are using the dry-run mode. Use `--dry-run=false` to disable this mode.")
		}

		if len(labelOptions.GitHub.Token) == 0 {
			labelOptions.GitHub.Token = os.Getenv("GITHUB_TOKEN")
		}

		required(labelOptions.GitHub.Token, "github.token")
		required(labelOptions.GitHub.Owner, "github.owner")
		required(labelOptions.GitHub.RepositoryName, "github.repo-name")
		required(labelOptions.RulesFilePath, "rules-path")

		err := cmd.Label(labelOptions)
		if err != nil {
			log.Println(err)
		}
		return nil
	}

	return labelCmd
}

func required(field string, fieldName string) error {
	if len(field) == 0 {
		log.Fatalf("%s is mandatory.", fieldName)
	}
	return nil
}
