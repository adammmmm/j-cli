# j-cli

```
A CLI tool for Junos devices, capable of configuration changes, configuration backups, operational commands and more.

Usage:
  j-cli [command]

Available Commands:
  cfg         Configuration Commands
  completion  Generate the autocompletion script for the specified shell
  diff        Show difference between current and previous commit
  help        Help about any command
  op          Operational Commands
  running     Show running configuration

Flags:
  -a, --authentication string   Authentication file in json format (default "auth.json")
      --backoff duration        Initial backoff duration between retries (default 2s)
  -d, --devices string          Device File in json format (default "devices.json")
  -h, --help                    help for j-cli
      --json                    Output results as JSON
      --retries int             Number of retries per device on failure (default 2)
  -t, --timeout duration        Timeout per device operation (default 30s)
  -w, --workers int             Number of concurrent device workers (default 8)
```
