package options

// Empty options.
type Empty struct{}

// Report options
type Report struct {
	Debug  bool    `description:"Debug mode."`
	DryRun bool    `long:"dry-run" description:"Dry run mode."`
	GitHub *GitHub `description:"GitHub options."`
	Slack  *Slack  `description:"Slack options."`
}

// Slack options
type Slack struct {
	Token     string `description:"Slack token."`
	ChannelID string `long:"channel" description:"Slack channel ID."`
	IconEmoji string `long:"bot-icon" description:"Bot icon emoji."`
	BotName   string `long:"bot-name" description:"Bot name."`
}

// GitHub options
type GitHub struct {
	Token          string `description:"GitHub token."`
	Owner          string `short:"o" description:"Repository owner."`
	RepositoryName string `long:"repo-name" short:"r" description:"Repository name."`
}

// Label options
type Label struct {
	Debug         bool    `description:"Debug mode."`
	DryRun        bool    `long:"dry-run" description:"Dry run mode."`
	WebHook       bool    `long:"web-hook" description:"Run as WebHook."`
	RulesFilePath string  `long:"rules-path" description:"Path to the rule file."`
	GitHub        *GitHub `description:"GitHub options."`
}
