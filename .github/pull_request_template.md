<!--
Title: use Conventional Commits, for example `fix(slo): handle an empty objective list`.
PRs are squash-merged, so the title becomes the commit message and feeds the changelog.
Mark breaking changes with `!`, for example `feat(resources)!: ...`.
-->

## What

<!-- What changes: behavior added, changed or removed. One short paragraph per change. -->

## Why

<!-- The problem this solves and why this approach. -->

## How it was tested

<!-- Checks actually run, each with the command or behavior and what you saw. Leave a box unticked if a check was skipped, and say why. -->

- [ ] `mise run gate` passes (lint, tests, build)
- [ ] `GCX_AGENT_MODE=false mise run reference` run, if commands, flags, config or env vars changed
- [ ]

## Related

<!--
`Fixes #123` for the issue this PR resolves, or `Related to #123`. Delete this section if nothing applies.
New commands, flags, providers and datasource kinds need an agreed feature proposal first:
https://github.com/grafana/gcx/blob/main/CONTRIBUTING.md#new-commands-need-an-issue-first
-->

---

**Before review**

- [ ] All commits are signed ([required in every Grafana repository](https://community.grafana.com/t/action-required-signed-commits-mandatory-for-all-grafana-repositories/163404))
- Authorship (tick one):
  - [ ] Written by a person
  - [ ] Written with an AI agent, directed and reviewed by the author
  - [ ] Written mostly by an AI agent, reviewed by the author
