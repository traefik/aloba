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

		err := required(reportOptions.GitHub.Token, "github.token")
		if err != nil {
			log.Fatal(err)
		}
		err = required(reportOptions.GitHub.Owner, "github.owner")
		if err != nil {
			log.Fatal(err)
		}
		err = required(reportOptions.GitHub.RepositoryName, "github.repo-name")
		if err != nil {
			log.Fatal(err)
		}
		err = required(reportOptions.Slack.Token, "slack.token")
		if err != nil {
			log.Fatal(err)
		}
		err = required(reportOptions.Slack.ChannelID, "slack.channel")
		if err != nil {
			log.Fatal(err)
		}

		err = cmd.Report(reportOptions)
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

	defaultPointerOptions := &options.Label{
		GitHub: &options.GitHub{},
		WebHook: &options.WebHook{
			Port: 80,
		},
	}

	labelCmd := &flaeg.Command{
		Name:                  "label",
		Description:           "Add labels to Pull Request",
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

		err := required(labelOptions.GitHub.Token, "github.token")
		if err != nil {
			log.Fatal(err)
		}
		err = required(labelOptions.GitHub.Owner, "github.owner")
		if err != nil {
			log.Fatal(err)
		}
		err = required(labelOptions.GitHub.RepositoryName, "github.repo-name")
		if err != nil {
			log.Fatal(err)
		}
		err = required(labelOptions.RulesFilePath, "rules-path")
		if err != nil {
			log.Fatal(err)
		}

		err = cmd.Label(labelOptions)
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
