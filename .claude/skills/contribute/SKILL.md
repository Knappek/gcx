---
name: contribute
description: >-
  Prepare a contribution to the gcx repository: report or fix a bug, make a
  maintenance change, or propose, place and implement a capability (provider,
  datasource kind, resource adapter, cloud command or bundled skill), then
  self-review the diff. Use inside a gcx checkout. Trigger on phrases like
  "contribute to gcx", "integrate with gcx", "add my product to gcx", "add a
  provider", "add a datasource", "new gcx command", "expose this API through
  gcx", "fix this gcx bug" or "get my gcx PR ready for review". NOT for
  operating a Grafana instance with gcx (use the gcx skill or a product skill
  such as slo-manage, synth-manage-checks or create-dashboard). NOT for
  installing or configuring gcx (use setup-gcx).
---

# Contribute to gcx

Read [AGENTS.md](../../../AGENTS.md) and [CONTRIBUTING.md](../../../CONTRIBUTING.md).
CONTRIBUTING decides which changes need an issue first; this skill applies it.
The contributor does not need to choose a provider, datasource, RFC or OpenSpec
workflow upfront.

## 1. Establish the route and current step

Identify the requested outcome and any existing issue, maintainer decision, RFC
or implementation branch. Inspect those records and the relevant code before
asking questions. Then state the route and the current step:

- **Direct PR:** bug fixes, documentation, tests, and performance or reliability
  work. Link the bug issue when one exists.
- **Issue first:** new commands, subcommand groups, providers, datasource kinds
  and resource types, and changes to an existing command's path, flags,
  positional arguments or output shape. These need a feature proposal and a
  maintainer comment on the issue agreeing with the placement before
  implementation. Grafana product teams don't need a proposal inside their own
  product area; shared commands and new top-level areas still do.

Reuse settled scope and design at each handoff. Work within the request and the
authorization already given: a draft-only or review-only request ends with that
artifact. Prepare the concrete result before asking for permission to publish,
and reuse permission already given for the same destination and scope.

## 2. Prepare the issue when the route needs one

Search existing issues and PRs first, and continue on a matching one.

- **Bug report:** use the [Bug report form](../../../.github/ISSUE_TEMPLATE/1-bug-report.yml).
  Establish what happened against what was expected, the steps to reproduce,
  and the gcx version. Record verified findings separately from assumptions.
- **Feature proposal:** use [propose-feature](../propose-feature/SKILL.md) to
  research and draft it.

Follow the GitHub Issues privacy rules in AGENTS.md. When filing is authorized,
create the issue with `gh issue create --body-file`, apply the form's
`action/needs-triage` label, and set the issue type (`Bug` or `Feature`) the way
propose-feature describes. Read the issue back and report its URL.

## 3. Load only the planning that is needed

Most changes go from the agreed issue straight to implementation. The
maintainer who agrees with the proposal says when more planning is needed:

- An RFC when design tradeoffs or scope need agreement: use
  [create-rfc](../create-rfc/SKILL.md), [review-rfc](../review-rfc/SKILL.md) for
  review, and [update-rfc](../update-rfc/SKILL.md) only when explicitly asked to
  reconcile an RFC with implemented changes.
- A tracked implementation plan: use [openspec-propose](../openspec-propose/SKILL.md)
  for **one PR**, linking the issue and any RFC. Use the OpenSpec
  apply/update/sync/archive workflows as needed, and get artifact IDs, paths and
  requirements from the `openspec` CLI's active schema.

## 4. Implement with the appropriate references

For a bug, inspect the relevant subsystem and add regression coverage. For
documentation or test-only work, inspect the affected material.

For a capability or command change, read
[capability implementation](references/capability-implementation.md) for
placement, backend readiness and the command contract. Record those decisions
in the issue or the existing planning record rather than a new document. Then
load:

- [Provider implementation](references/provider-implementation.md) for a Cloud
  product provider, including conditional resource adapters and smoke tests.
- [Datasource implementation](references/datasource-implementation.md) for a new
  datasource kind, including shared queries, Explore links, codecs and routing.
- [Distribution and gates](references/distribution-and-gates.md) when resource or
  datasource wiring, command annotations, or bundled skill coverage changes.

Use the existing subsystem implementation for other changes. Resolve routine
choices from repository evidence and preserve CONSTITUTION, DESIGN and
architecture rules.

**Ask** when the repo, the governing docs and a verified probe can't settle
something that changes the work: missing product or API facts, unclear
ownership or API stability, conflicting governing docs, a new verb or placement
that the naming guide doesn't cover, or a scope choice that changes what gets
built. Group related questions, each with the evidence and a recommended
answer. Don't ask just because the change is large.

**Stop** and report instead of implementing when:

- the route needs an issue first and no maintainer has agreed with the placement;
- the change violates CONSTITUTION with no compliant alternative, or needs a waiver;
- a missing backend, auth or ownership prerequisite prevents an honest
  implementation (report the prerequisite and its owner).

## 5. Verify and prepare review

Read [self-review](references/self-review.md) and apply the checks the actual
diff triggers, including incremental fixes. Run the AGENTS.md quality checks
and relevant domain verification with the built binary. Report failed or
unavailable checks with their reason: a check that did not run is not a pass.

When OpenSpec is used, keep task state consistent with implementation and
verification evidence, and finish and archive the PR-scoped change through its
workflow. Implementation PRs reference the issue; the PR that completes it
closes it. Fill in the [pull request template](../../../.github/pull_request_template.md)
and end with this summary:

```text
CONTRIBUTION SUMMARY — <change>

Route: <direct PR | issue first: #issue>
Placement: <necessity + path + wiring, or n/a> — <one-line rationale>
Changed: <what this diff does>
Non-goals: <out of scope; product-owned items with owners>

Unresolved decisions or risks: <what a reviewer must decide, or none>
Unverified assumptions: <claims not probed, and how to probe them, or none>
Checks failed or skipped: <gate: reason, or none>
Architecture deviations: <deviation + rationale, or none>
```

An empty field is a claim: write `none` only when it is true.
