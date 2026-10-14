## gcx irm oncall teams update-working-hours

Replace the business hours of a team.

### Synopsis

Replace the business hours of a team.

The file holds the document that get-working-hours emits. The intervals in it
become the team's entire set: an interval that the file omits is deleted. Pass
{"irm_working_hours": null} to clear every business hour of the team.

Times are wall-clock times (HH:MM:SS) in the team's timezone. A weekday with no
ranges is a day off, but every interval needs a time range on at least one
weekday.

Example file:

  irm_working_hours:
    - name: Office hours
      irm_working_hours:
        monday:    [{start: "09:00:00", end: "17:00:00"}]
        tuesday:   [{start: "09:00:00", end: "17:00:00"}]
        wednesday: [{start: "09:00:00", end: "17:00:00"}]
        thursday:  [{start: "09:00:00", end: "17:00:00"}]
        friday:    [{start: "09:00:00", end: "17:00:00"}]
        saturday:  []
        sunday:    []

The command prints the stored hours as a weekly grid. Use -o yaml or -o json
to emit the stored document instead.

```
gcx irm oncall teams update-working-hours <id> [flags]
```

### Options

```
  -f, --filename string   File containing the resource definition (JSON/YAML, use - for stdin)
  -h, --help              help for update-working-hours
      --jq string         jq expression to apply to JSON output. Mutually exclusive with --json.
      --json string       Comma-separated list of fields to include in JSON output, or 'list' (or '?') to discover available fields
  -o, --output string     Output format. One of: agents, json, table, yaml (default "table")
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

