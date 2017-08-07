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
		err := cmd.Label(labelOptions)
		if err != nil {
			log.Println(err)
		}
		return nil
	}

	return labelCmd
}
