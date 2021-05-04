# Myrmica Aloba: Add labels and milestone on pull requests and issues.

[![Release](https://img.shields.io/github/release/traefik/aloba.svg?style=flat)](https://github.com/traefik/aloba/releases/latest)
[![Build Status](https://github.com/traefik/aloba/workflows/Main/badge.svg?branch=master)](https://github.com/traefik/aloba/actions)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/traefik/aloba?label=Docker%20Image)](https://hub.docker.com/r/traefik/aloba/builds/)

## Overview

Available Commands:
- `action`: GitHub Action (Add labels and milestone on pull requests and issues.)
- `label`: Add labels and milestone on pull requests and issues.
- `report`: Create a report and publish on Slack.
- `version`: Display the version.

## Manage GitHub labels.

- on new issue: adds the label `status/0-needs-triage`
- on new pull request:
    - adds the label `status/0-needs-triage`
    - add labels based on [rules](#rules).
    - adds a milestone (if a milestone matches the based branch of the PR).
    - adds a label related to the size of the pull request.

### Command `action`

```yaml
GitHub Action

Usage: action [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: action [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)

Flags:
    --debug   Debug mode.                        (default "false")
    --dry-run Dry run mode.                      (default "true")
-h, --help    Print Help (this message) and exit 
```

- `GITHUB_TOKEN`: Github Token.
- `.github/aloba-rules.toml`: the rules to apply.

#### Examples:

```hcl
workflow "Aloba: Issues" {
  on = "issues"
  resolves = ["issue-labels"]
}

action "issue-labels" {
  uses = "docker://traefik/aloba"
  secrets = ["GITHUB_TOKEN"]
  args = "action --dry-run=false"
}

workflow "Aloba: Pull Requests" {
  on = "pull_request"
  resolves = ["pull-request-labels"]
}

action "pull-request-labels" {
  uses = "docker://traefik/aloba"
  secrets = ["GITHUB_TOKEN"]
  args = "action --dry-run=false"
}
```

### Command `label`

```yaml
Add labels and milestone on pull requests and issues.

Usage: label [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: label [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)

Flags:
    --debug            Debug mode.                        (default "false")
    --dry-run          Dry run mode.                      (default "true")
    --github           GitHub options.                    (default "true")
-o, --github.owner     Repository owner.
-r, --github.repo-name Repository name.
    --github.token     GitHub token.
    --rules-path       Path to the rule file.             (default "./rules.toml")
    --web-hook         Run as WebHook.                    (default "true")
    --web-hook.port    WebHook port.                      (default "80")
    --web-hook.secret  WebHook secret.
-h, --help             Print Help (this message) and exit
```

- `GITHUB_TOKEN`: Github Token.
- `WEBHOOK_SECRET`: Github WebHook Secret.

#### Examples:

```shell
aloba label -o traefik -r traefik --web-hook --github.token="xxxxxxxxxx"
```

### Rules

```toml
[[Rules]]
  Label = "area/vegetable"
  Regex = "(?i).*(tomate|carotte).*"

[[Rules]]
  Label = "area/cheese"
  Regex = "cheese/.*"

[[Rules]]
  Label = "area/infrastructure"
  Regex = "(?i)(\\.github|script/).*"

[Limits]
  [Limits.Small]
    SumLimit = 150
    DiffLimit = 70
    FilesLimit = 20
  [Limits.Medium]
    SumLimit = 400
    DiffLimit = 200
    FilesLimit = 50
```

## Command `report`

```yaml
Create a report and publish on Slack.

Usage: report [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: report [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)

Flags:
    --debug            Debug mode.                        (default "false")
    --dry-run          Dry run mode.                      (default "true")
    --github           GitHub options.                    (default "true")
-o, --github.owner     Repository owner.
-r, --github.repo-name Repository name.
    --github.token     GitHub token.
    --slack            Slack options.                     (default "true")
    --slack.bot-icon   Bot icon emoji.                    (default ":captainpr:")
    --slack.bot-name   Bot name.                          (default "CaptainPR")
    --slack.channel    Slack channel ID.
    --slack.token      Slack token.
-h, --help             Print Help (this message) and exit
```

- `GITHUB_TOKEN`: Github Token.
- `SLACK_TOKEN`: Slack Token.

### Examples:

```shell
aloba report -o traefik -r traefik --github.token="xxxxxxxxxx" --slack.token="xxxxxxxxxx" --slack.channel=C0CDT22PJ
```

## The Mymirca colony

- [Myrmica Lobicornis](https://github.com/traefik/lobicornis) 🐜: Update and merge pull requests.
- [Myrmica Aloba](https://github.com/traefik/aloba) 🐜: Add labels and milestone on pull requests and issues.
- [Messor Structor](https://github.com/traefik/structor) 🐜: Manage multiple documentation versions with Mkdocs.
- [Lasius Mixtus](https://github.com/traefik/mixtus) 🐜: Publish documentation to a GitHub repository from another.
- [Myrmica Bibikoffi](https://github.com/traefik/bibikoffi) 🐜: Closes stale issues
- [Chalepoxenus Kutteri](https://github.com/traefik/kutteri) 🐜: Track a GitHub repository and publish on Slack.
- [Myrmica Gallienii](https://github.com/traefik/gallienii) 🐜: Keep Forks Synchronized

## What does Myrmica Aloba mean?

![Myrmica Aloba](http://www.antwiki.org/wiki/images/8/8c/Myrmica_aloba_H_casent0907652.jpg)
