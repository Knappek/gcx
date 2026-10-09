# Requests for Comments

RFCs propose engineering changes for review: the problem, user experience, technical design, tradeoffs, and open questions. A proposal is not evidence that a feature has been implemented or deployed. Existing architecture decisions remain documented in [ADRs](../adrs/).

## Index

| RFC | Status |
| --- | --- |
| [001: Unify alerting configuration and resource workflows](001-alerting-provider-refactor.md) | Proposed |
| [002: Manage Assistant Watchers through gcx](002-assistant-watchers.md) | Proposed |

## Authoring and maintenance

An RFC is needed when design tradeoffs or scope need agreement before implementation; the maintainer who agrees with the feature proposal says so (see [CONTRIBUTING.md](../../CONTRIBUTING.md#new-commands-need-an-issue-first)). An RFC records design and scope decisions for implementation PRs; those PRs may use [OpenSpec](../../openspec/) when they need a tracked implementation plan. Start agent-assisted contribution work with [`contribute`](../../.claude/skills/contribute/SKILL.md).

Use `NNN-title.md`, choosing the next unused number starting at `001`, title the document `RFC NNN: <Title>`, and add it to this index. Record status in the index rather than in the RFC. Keep later revisions in the same file. Define proposed terms in the RFC. When a feature ships, move its terms to the [glossary](../glossary/README.md) and link to them from the RFC instead of repeating them.

The repository provides these contributor skills:

- [create-rfc](../../.claude/skills/create-rfc/SKILL.md): research and develop a proposal through review, using the [Rust-based template](../../.claude/skills/create-rfc/references/rfc-template.md).
- [review-rfc](../../.claude/skills/review-rfc/SKILL.md): review an RFC before it is shared or accepted, checking its claims, scope, rules, consistency and structure.
- [update-rfc](../../.claude/skills/update-rfc/SKILL.md): reconcile an existing RFC with implementation evidence when explicitly requested, distinguishing implementation, verification, and deployment.

These workflows do not authorize implementation or publication. RFCs and their examples must be suitable for this public repository; restricted research belongs outside the checkout.
