## gcx irm oncall teams get-working-hours

Get the business hours of a team.

### Synopsis

Get the business hours of a team.

A team's business hours are a list of named intervals. Each interval holds up
to three time ranges per weekday. A team is in business hours when the current
time falls inside any of its intervals; the "pause escalation until team's
working hours" escalation step waits for that moment.

By default the command prints a weekly grid with one row per interval. Use
-o yaml or -o json to emit the document {"irm_working_hours": [...]} (the list
is null when the team has no business hours). Edit that document, then pass it
back through update-working-hours.

```
gcx irm oncall teams get-working-hours <id> [flags]
```

### Options

```
  -h, --help            help for get-working-hours
      --jq string       jq expression to apply to JSON output. Mutually exclusive with --json.
      --json string     Comma-separated list of fields to include in JSON output, or 'list' (or '?') to discover available fields
  -o, --output string   Output format. One of: agents, json, table, yaml (default "table")
```

### Options inherited from parent commands

```
      --agent                       Enable agent mode (JSON output, no color). Auto-detected from known agent identity variables. Set GCX_AGENT_NAME to identify a supported harness, or GCX_AGENT_MODE to control the mode.
      --config string               Path to the configuration file to use
      --context string              Name of the context to use (overrides current-context in config)
      --insecure-log-http-payload   Log full HTTP request/response bodies including raw credentials, authorization tokens, cookies, and OAuth refresh tokens. Requires -vvv. Do not ship these logs.
      --no-color                    Disable color output
      --no-truncate                 Disable table column truncation (auto-enabled when stdout is piped)
  -v, --verbose count               Verbose mode. Multiple -v options increase the verbosity (maximum: 3).
```

### SEE ALSO

* [gcx irm oncall teams](gcx_irm_oncall_teams.md)	 - Manage OnCall teams.

