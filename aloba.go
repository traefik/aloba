package main

import (
	"log"
	"os"

	"github.com/containous/flaeg"
	"github.com/containous/myrmica-aloba/cmd"
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

	emptyConfig := &cmd.NoOption{}

	rootCmd := &flaeg.Command{
		Name:                  "aloba",
		Description:           "Myrmica Aloba: Manage GitHub labels.",
		Config:                emptyConfig,
		DefaultPointersConfig: &cmd.NoOption{},
		Run: func() error {
			// No op
			return nil
		},
	}

	return rootCmd
}

func createReportCommand() *flaeg.Command {

	reportOptions := &cmd.ReportOptions{
		IconEmoji: ":captainpr:",
		BotName:   "CaptainPR",
		DryRun:    true,
	}

	reportCmd := &flaeg.Command{
		Name:                  "report",
		Description:           "Create a report and publish on Slack.",
		Config:                reportOptions,
		DefaultPointersConfig: &cmd.ReportOptions{},
	}
	reportCmd.Run = func() error {
		if reportOptions.DryRun {
			log.Print("IMPORTANT: you are using the dry-run mode. Use `--dry-run=false` to disable this mode.")
		}
		required(reportOptions.GitHubToken, "github-token")
		required(reportOptions.Owner, "owner")
		required(reportOptions.RepositoryName, "repo-name")
		required(reportOptions.SlackToken, "slack-token")
		required(reportOptions.ChannelID, "channel-id")

		err := cmd.Report(reportOptions)
		if err != nil {
			log.Println(err)
		}
		return nil
	}

	return reportCmd
}

func createLabelCommand() *flaeg.Command {

	labelOptions := &cmd.LabelOptions{
		RulesFilePath: "./rules.toml",
		DryRun:        true,
	}

	labelCmd := &flaeg.Command{
		Name:                  "label",
		Description:           "Add labels to Pull Request",
		Config:                labelOptions,
		DefaultPointersConfig: &cmd.LabelOptions{},
	}
	labelCmd.Run = func() error {
		if labelOptions.DryRun {
			log.Print("IMPORTANT: you are using the dry-run mode. Use `--dry-run=false` to disable this mode.")
		}
		required(labelOptions.GitHubToken, "github-token")
		required(labelOptions.Owner, "owner")
		required(labelOptions.RepositoryName, "repo-name")
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
