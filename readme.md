# Myrmica Aloba: Manage GitHub labels.

[![Build Status](https://travis-ci.org/containous/aloba.svg?branch=master)](https://travis-ci.org/containous/aloba)
[![Docker Build Status](https://img.shields.io/docker/build/containous/aloba.svg)](https://hub.docker.com/r/containous/aloba/builds/)


```shell
Myrmica Aloba: Manage GitHub labels.

Usage: aloba [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: aloba [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)

Available Commands:
	label                                              Add labels to Pull Request
	report                                             Create a report and publish on Slack.
Use "aloba [command] --help" for more information about a command.

Flags:
-h, --help Print Help (this message) and exit 
```

Sub-command `label`:
```shell
Add labels to Pull Request

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

Sub-command `report`:
```shell
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

## Examples:

Sub-command `label`:
```shell
aloba label -o containous -r traefik --web-hook --github.token="xxxxxxxxxx"
```

Sub-command `report`:
```shell
aloba report -o containous -r traefik --github.token="xxxxxxxxxx" --slack.token="xxxxxxxxxx" --slack.channel=C0CDT22PJ
```