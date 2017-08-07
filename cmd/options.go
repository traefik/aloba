package cmd

// NoOption empty struct.
type NoOption struct{}

// ReportOptions report options
type ReportOptions struct {
	SlackToken     string `long:"slack-token" description:"Slack token."`
	ChannelID      string `long:"channel-id" short:"c" description:"Slack channel ID."`
	GitHubToken    string `long:"github-token" description:"GitHub token."`
	Owner          string `short:"o" description:"Repository owner."`
	RepositoryName string `long:"repo-name" short:"r" description:"Repository name."`
	IconEmoji      string `long:"bot-icon" description:"Bot icon emoji."`
	BotName        string `long:"bot-name" description:"Bot name."`
	Debug          bool   `description:"Debug mode."`
	DryRun         bool   `long:"dry-run" description:"Dry run mode."`
}

// LabelOptions label options
type LabelOptions struct {
	GitHubToken    string `long:"github-token" description:"GitHub token."`
	Owner          string `short:"o" description:"Repository owner."`
	RepositoryName string `long:"repo-name" short:"r" description:"Repository name."`
	RulesFilePath  string `long:"rules-path" description:"xxxxxxxxxxxx"`
	Debug          bool   `description:"Debug mode."`
	DryRun         bool   `long:"dry-run" description:"Dry run mode."`
	WebHook        bool   `long:"web-hook" description:"Run as WebHook."`
}
