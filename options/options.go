package options

// Empty options.
type Empty struct{}

// Report options.
type Report struct {
	Debug      bool    `description:"Debug mode."`
	DryRun     bool    `long:"dry-run" description:"Dry run mode."`
	ServerMode bool    `long:"server" description:"Server mode."`
	ServerPort int     `long:"port" description:"Server port."`
	GitHub     *GitHub `description:"GitHub options."`
	Slack      *Slack  `description:"Slack options."`
}

// Slack options.
type Slack struct {
	Token     string `description:"Slack token."`
	ChannelID string `long:"channel" description:"Slack channel ID."`
	IconEmoji string `long:"bot-icon" description:"Bot icon emoji."`
	BotName   string `long:"bot-name" description:"Bot name."`
}

// GitHub options.
type GitHub struct {
	Token          string `description:"GitHub token."`
	Owner          string `short:"o" description:"Repository owner."`
	RepositoryName string `long:"repo-name" short:"r" description:"Repository name."`
}

// Label options.
type Label struct {
	Debug         bool     `description:"Debug mode."`
	DryRun        bool     `long:"dry-run" description:"Dry run mode."`
	WebHook       *WebHook `long:"web-hook" description:"Run as WebHook."`
	RulesFilePath string   `long:"rules-path" description:"Path to the rule file."`
	GitHub        *GitHub  `description:"GitHub options."`
}

// GitHubAction options.
type GitHubAction struct {
	Debug  bool `description:"Debug mode."`
	DryRun bool `long:"dry-run" description:"Dry run mode."`
}

// WebHook options.
type WebHook struct {
	Port   int    `description:"WebHook port."`
	Secret string `description:"WebHook secret."`
}
