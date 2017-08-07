# Myrmica Aloba: Manage GitHub labels.


```
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
```
Add labels to Pull Request

Usage: label [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: label [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)

Flags:
    --debug        Debug mode.                        (default "false")
    --dry-run      Dry run mode.                      (default "true")
    --github-token GitHub token.                      
-o, --owner        Repository owner.                  
-r, --repo-name    Repository name.                   
    --rules-path   xxxxxxxxxxxx                       (default "./rules.toml")
    --web-hook     Run as WebHook.                    (default "false")
-h, --help         Print Help (this message) and exit 

```

Sub-command `report`:
```
Create a report and publish on Slack.

Usage: report [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: report [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)

Flags:
    --bot-icon     Bot icon emoji.                    (default ":captainpr:")
    --bot-name     Bot name.                          (default "CaptainPR")
-c, --channel-id   Slack channel ID.                  
    --debug        Debug mode.                        (default "false")
    --dry-run      Dry run mode.                      (default "true")
    --github-token GitHub token.                      
-o, --owner        Repository owner.                  
-r, --repo-name    Repository name.                   
    --slack-token  Slack token.                       
-h, --help         Print Help (this message) and exit 
```

## Examples:

Sub-command `label`:
```bash
aloba label -o containous -r traefik --web-hook --github-token="xxxxxxxxxx"
```


Sub-command `report`:
```bash
aloba report -o containous -r traefik --github-token="xxxxxxxxxx" --slack-token="xxxxxxxxxx" -c C0CDT22PJ
```