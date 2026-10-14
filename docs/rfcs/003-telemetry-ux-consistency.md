# RFC: Consistent telemetry command behavior

**Status:** Draft · **Scope:** Metrics, logs, traces, and profiles

## Summary

Define a consistent UX contract across all relevant metrics, logs, traces, and profiles subcommands: discovery and catalogs, primary queries, time-series queries, object retrieval, comparisons and diagnostics, artifact outputs, and adaptive telemetry operations. Normalize equivalent behavior while retaining domain-specific capabilities. The target uses `--query` throughout, defaults window-based reads to the last hour, and wraps structured reads in shared context/status metadata while retaining domain payloads. Explicit names/values discovery is one application of that contract.

The target experience is specified independently of compatibility. Adoption of the common structured contract is opt-in during migration; the mechanism and retirement policy belong in the separate compatibility discussion below. This proposal does not change the query languages, flatten trace types or scopes, or promise identical backend capabilities.

## Motivation

An investigation moves between metrics, logs, traces, and profiles. Users should be able to transfer knowledge about discovery, time bounds, datasource selection, and error recovery. Agents benefit from the same predictability, but the benefit must be established through task success and total token spend rather than inferred from similar command names.

The [earlier cross-signal proposal](../adrs/signal-provider-ux/001-cross-signal-command-consistency.md) established shared datasource flags, primary query commands, and `labels --label NAME`. Several conventions are now implemented. This RFC develops a clearer discovery interface and addresses behavioral gaps that remain. The earlier document's pre-GA migration policy is historical; the current [compatibility contract](../../CONSTITUTION.md#cli-grammar) governs delivery.

## Scope

The audit covers the complete four signal command trees and their datasource counterparts. It includes label and series discovery, metric-name search and metadata, profile-type and data-range discovery, primary and time-series queries, trace retrieval/diff/baseline, profile exemplars/anomalies, artifact rendering, and adaptive telemetry read and mutation commands. A coverage matrix below records where behavior should align and where a domain-specific operation has no equivalent.

The accepted label decisions remain in force: explicit names/values operations, an opt-in common schema, a one-hour window, a 100-item cap, a common single-expression filter, and strict discovery-filter validation. Those exact defaults do not automatically apply to ranked analyses, whole-object retrieval, retention inspection, or mutations. Their target behavior is specified by operation class.

The scope is CLI behavior and its backing-API mapping. A universal query language, backend API redesign, new analytical algorithms, and implementing missing products are outside scope. Compatibility and migration remain a separate team discussion.

## Current behavior and backing capabilities

The audit covers command registration, options, execution, codecs, and backing clients across the four signal trees. These are source-confirmed findings, not live integration results. Public backend source and API contracts inform the mappings; supported deployed versions remain a release gate.

### Complete command-family coverage

“Shared” means the signal and datasource paths reuse constructors. Equivalent operations should share behavior even when a documented naming exception exists.

| Family | Current surface and backing operation | Target assessment |
| --- | --- | --- |
| Primary queries | Metrics `query`: PromQL instant/range through Grafana query API. Logs `query`: log lines through the same transport. Traces `query`/`search`: Tempo search. Profiles `query`: merged stack/span/profile APIs. All shared. | Align expression input, time validation, datasource selection, errors, and output promises. Preserve native languages and distinct result types. Eliminate silently ignored options. |
| Time-series queries | Logs and traces `metrics`; profiles `metrics`; metrics range `query`. Shared. Tempo uses metrics query/range, Pyroscope SelectSeries. | Apply `--step` only where meaningful; distinguish evaluation resolution, ranking, and result caps. Fix unsupported wide output and incomplete links. |
| Label discovery | Four shared `labels` families; traces also `tags`. Native label/tag APIs. | Explicit names/values commands and the settled discovery contract below. |
| Series discovery | Metrics `series`, logs `series`, profiles `series`; only metrics lacks a datasource mount. Native series APIs. | Common window, selector cardinality, cap/coverage, and mount behavior; retain profile label projection. Resolve Pyroscope read-path-dependent multi-selector semantics. |
| Metric catalogs | Metrics `list-names`; shared `metadata`. Label-values API versus descriptive metadata API. | Names inherit discovery bounds; metadata remains a descriptive catalog with strict input validation and bounded output. No invented occurrence window for metadata. |
| Ranked search | Metrics `search {metric-names,label-names,label-values}`; shared constructors exposed as Prometheus `search-*`. Experimental search endpoints. | Retain score/sort/type/unit, feature gates, warnings, and existing incomplete-result metadata. Add explicit discovery defaults; do not turn ranking into alphabetical enumeration. |
| Cardinality | Metrics `cardinality {label-names,label-values}`. Mimir current-ingester APIs. | Preserve `inmemory`/`active`, per-label result units, and supported limit range. Expose datasource counterpart and state that this is current-state analysis, not historical discovery. |
| Billing | Metrics `billing {query,labels,series}`. Ordinary constructors with a billing datasource fallback. | Inherit applicable contracts; retain convenience routing and explicit datasource override. |
| Trace retrieval | Shared `get`; Tempo trace-by-ID API, optional time bounds and client-side filtering/pruning. | Preserve whole-object retrieval, tree views, and ID semantics; disclose transformations and Explore link fidelity. No arbitrary item cap on a trace. |
| Trace comparison | Shared `diff` and `baseline`; trace-diff API, or seed lookup plus generated search. | Preserve experimental capability errors, two-ID diff, seed-centered baseline window, unranked candidates, and partial-seed warnings. Bound candidate collections without presenting them as exhaustive or ranked. |
| Profile catalogs | Shared `list-profile-types` and `data-range`; ProfileTypes and GetProfileStats. | Types use discovery windows and empty-result contract. Data range remains tenant-wide ingestion/queryable-bound diagnostics; no selector or one-hour filter. |
| Profile investigation | Shared `exemplars {profile,span}` and `query-anomalies`; exemplar/heatmap and anomaly APIs. | Preserve profile types, units and drilldown IDs. Disclose series-selection bounds separately from exemplar/anomaly result caps; align empty assertions and links. |
| Adaptive Metrics | `recommendations {list,diff,apply}`; `rules`, `segments`, `exemptions`, each with list/get/create/update/delete. Tenant management APIs. | Consistent inventory bounds, snapshot provenance, strict inputs, preview/apply outcomes; keep domain schemas and tenant authentication. |
| Adaptive Logs | `patterns {list,stats}`; `exemptions`, `segments` list/create/update/delete; `drop-rules` list/get/create/update/delete. Tenant management APIs. | Separate presentation rollups from data limits, disclose snapshot provenance, add supported `get ID` for exemptions/segments, and align mutation outcomes. |
| Adaptive Traces and Profiles | Traces `policies` CRUD and `recommendations` list/apply/dismiss; tenant APIs. Profiles `adaptive` is a stub. | Align actual management operations; report unsupported Profiles capability as failure instead of printing unavailable and returning success. Do not implement a missing product for symmetry. |

Registration evidence: [signal providers](../../internal/providers/), [datasource mounts](../../internal/datasources/providers/). Execution and API mapping evidence: [datasource commands](../../internal/datasources/), [backend query clients](../../internal/query/), and [adaptive providers](../../internal/providers/). This inventory includes the management leaves; it does not equate them with datasource queries.

### Material inconsistencies beyond labels

The [profile query implementation](../../internal/datasources/pyroscope/query.go) accepts profile-ID and stacktrace filters but its pprof branch omits them, although the backing merge-profile API supports them. Selecting an output format can therefore change the requested data. Trace search and profile merge accept `--step` and discard it; profile metrics ranking can override an explicitly supplied step. The [wide codec](../../internal/datasources/query/codecs.go) is registered for profile query and trace metrics without matching encoders, producing an error after fetching data.

Limits also mean different things. Logs and trace search describe zero as unlimited while omitting the native limit leaves backend defaults in force. Adaptive Metrics recommendations trim every format with `--top`; [Adaptive Logs patterns](../../internal/providers/logs/adaptive/commands.go) use `--top` only for a table rollup and return the full filtered set in structured formats. Several [adaptive lists](../../internal/providers/metrics/adaptive/commands.go) and [profile exemplar](../../internal/query/pyroscope/exemplars.go)/[anomaly results](../../internal/query/pyroscope/anomalies.go) trim data without machine-readable truncation information, unlike metric names and ranked search.

Profile series advertises selector union, but public Pyroscope source combines multiple selectors differently across read paths. Logs series lacks time flags supported by its API. Discovery and query leaves differ in explicit-empty validation, type checking, empty assertions, and Explore support. These are behavioral issues with effects on task correctness, not merely inconsistent spelling.

PR #1325 improves Loki label discovery by adding selector filtering, explicit-empty validation, positional-argument rejection, and shared codec routing. Its single `--query` matches the endpoint's single-selector capability. It does not introduce the existing names/values operation switching or resolve the broader gaps. Its overall direction improves UX; the common contract should guide subsequent work without treating useful incremental filtering as a regression.

## Proposed UX

### Common contract, specialized operations

Consistency applies to equivalent actions: selecting a datasource, carrying a window, filtering, understanding bounds, interpreting empty results, and recovering from errors. A trace lookup, cardinality diagnostic, ranked anomaly result, and configuration update remain distinct operations. Help and agent command metadata must identify operation class and supported capabilities without requiring trial requests.

| Operation class | Target time behavior | Target result controls |
| --- | --- | --- |
| Historical inventory: labels, names, series, profile types; historical ranked search | Explicit absolute bounds; default last hour | Inventory default 100 items; common coverage. Ranked search retains its domain ranking and current 50-item default, with explicit bounds and warnings. |
| Expression and time-series queries | Default window-based reads to last hour; instant evaluation is explicit; report effective mode/window | Retain domain controls and units. Step only for evaluation. Zero removes only a gcx cap where applicable. |
| Whole-object retrieval and comparisons | Optional lookup bounds; baseline defaults to seed ±30 minutes; diff retains its shared lookup window | Whole objects remain whole unless an explicit transform is requested. Candidate searches use a named candidate cap and coverage. |
| Descriptive catalogs, current-state diagnostics, adaptive observations/configuration | Metadata is descriptive; cardinality is current-state; adaptive observations are latest snapshots; configurations are current state | Bound collection reads and disclose coverage; preserve diagnostic units, snapshot freshness, and ranking. No fictional historical filter. |
| Mutations and artifacts | Preserve domain scope and exact source request | Explicit mutation outcomes; artifact formats preserve filters and represent generation failures. No inventory cap applied to artifact bytes or update batches. |

The shared query vocabulary, read envelope, and one-hour window-based default are target decisions. Numeric caps remain operation-specific, as specified below; the agreed 100-item label default is unchanged.

### Query grammar and time modes

Use `--query EXPR` as the canonical expression flag on execution and discovery commands. An execution command that takes an expression permits one positional expression as shorthand; supplying both forms, repeating the flag, or explicitly providing a blank expression is a usage error. Commands whose positional argument is an identifier keep that identifier positional. Their optional narrowing expression is `--query` where the backing operation supports one. A command without an expression capability does not advertise the flag. Native PromQL, LogQL, TraceQL and profile selectors remain unchanged.

Ranked metadata search has separate search terms and selector inputs: positional `TERM...` retains its native OR-combined term semantics, while `--query` narrows the candidate series with one native selector. This is not the positional-expression shorthand, so terms plus `--query` are valid together. Repeated native `--match` union remains available where supported and conflicts with `--query`, not with search terms. Preserve terms, selector and ranking independently in result context.

```sh
# Proposed: OR-combined search terms plus selector narrowing.
gcx metrics search metric-names request latency --query '{job="api"}' --since 1h -d metrics-uid
```

Window-based reads default to the last hour, including log lines, trace search, merged profiles, profile exemplars/anomalies, and range evaluation. `--since DURATION` ends at `--to` or now; `--from` requires `--to`; `--since` and `--from` conflict. Resolve relative times against one clock value and return absolute effective bounds. `--to` alone is invalid. The chosen codec and human/agent mode never change the window.

Time-series operations default to range mode and accept `--instant` only where the backend provides an instant operation. Prometheus/LogQL instant evaluation uses `--time` (default now), mutually exclusive with range bounds. Tempo instant metrics aggregate over an explicit or default one-hour window; that distinction is named `window_aggregate` in result context. It is not presented as a point evaluation. Profile leaderboard aggregation is separately selected with `--rank`; it retains a window and does not advertise instant evaluation. Trace search, log retrieval and profile merge have no instant or step capability.

`--step` controls range evaluation resolution only. A supplied step must be positive and accepted by the backend; reject it in instant, window-aggregate, search, merge, and rank modes before I/O. When omitted, preserve the backend/datasource resolution policy, record the effective value if known, and mark a server-chosen value as unknown rather than inventing one. Log-line sampling remains a separate native capability. No format may silently change query mode or resolution.

```sh
# Proposed target: the same expression input and window across signals.
gcx metrics query --query 'rate(http_requests_total[5m])' --since 1h --step 1m -d metrics-uid
gcx logs metrics --query 'sum(rate({job="api"}[5m]))' --since 1h --step 1m -d logs-uid
gcx traces query --query '{ resource.service.name = "api" }' --since 1h -d traces-uid
gcx profiles metrics --query '{service_name="api"}' --profile-type process_cpu:cpu:nanoseconds:cpu:nanoseconds --since 1h --step 1m -d profiles-uid

# Point evaluation is explicit; --time without --instant is rejected.
gcx metrics query --query 'up' --instant --time 2026-10-12T10:00:00Z -d metrics-uid
```

### Operation vocabulary and result units

Keep familiar domain operations, but use explicit enumeration verbs and a single meaning for each common control. The following is the target command surface; compatibility spellings are a separate concern.

| Use case | Target command(s) | Result and controls |
| --- | --- | --- |
| Evaluate metrics | `metrics query` | Range matrix by default; point evaluation with `--instant`. Preserve complete samples within each returned series. |
| Retrieve logs / evaluate LogQL | `logs query`, `logs metrics` | Log entries versus metric matrix. Line order is explicit through `--direction`; time-series semantics follow query modes above. |
| Find traces / derive metrics | `traces query`, `traces metrics` | Trace summaries versus metric matrix/window aggregate. Search does not imply ranked similarity. |
| Merge profiles / derive metrics | `profiles query`, `profiles metrics` | Stack/profile payload versus time series. `--rank` selects whole-window leaderboard; `--group-by` and `--aggregation` retain domain meaning. |
| Discover labels and series | `<signal> list-labels`, `list-label-values NAME`; `metrics/logs/profiles list-series` | Names, typed values, or series label sets; no trace-series facade. |
| Browse metric catalogs | `metrics list-names`, `metrics metadata` | Names versus descriptive name/type/help/unit records. Metadata retains `--metric` narrowing. |
| Search metric metadata | `metrics search {metric-names,label-names,label-values}` | Ranked matches; native term/score/sort controls remain distinct from selector narrowing. |
| Diagnose cardinality | `metrics cardinality {label-names,label-values}` | Current-state counts; limit units are explicitly per requested label where required. |
| Retrieve and compare traces | `traces get ID`, `traces diff ID_A ID_B`, `traces baseline ID` | Whole trace, structured diff, or candidate trace summaries. |
| Profile discovery and drilldown | `profiles list-profile-types`, `data-range`, `exemplars {profile,span}`, `query-anomalies` | Type catalog, diagnostic bounds, exemplar references, or scored anomaly records. |
| Inspect adaptive state | `<signal> adaptive <collection> list/get`; domain recommendation diff and pattern stats | Current resource objects, generated observations or domain comparisons; only supported collections are mounted. |
| Change adaptive state | Existing domain create/update/delete/apply/dismiss operations | Existing resource schemas and explicit per-item mutation outcomes. No fabricated symmetry for missing products. |

`--limit` always caps returned items independent of codec, with the unit declared by the command and result. `--rank` chooses profile aggregation/ranking, `--sort` chooses a supported ordering, and `--table-top` controls only a presentation rollup. Remove overloaded `--top`/`--top-n` from the target vocabulary through the deferred migration. Existing search and cardinality ranking retain their domain algorithms; identical names do not imply interchangeable scoring.

| Result family | Target default and unit | Zero / other bounds |
| --- | --- | --- |
| Historical inventories and metric metadata | 100 names, values, label sets, types, or metadata records | Zero removes gcx cap. |
| Ranked metric search | 50 matches | Zero removes gcx cap; backend stream/scan limits remain visible. |
| Logs / trace search / baseline | 100 log entries / 20 traces / 20 candidates | Zero removes gcx cap, not backend defaults; adapters disclose inability to establish exhaustion. |
| Profile metrics / exemplars / anomalies | 10 series or ranked groups / 100 exemplars / 100 anomalies | Zero removes gcx item cap; series-selection budgets are separately identified. |
| Adaptive inventories and observations | 50 resource or observation items | Zero removes gcx cap; stats aggregates refer to the stated full snapshot or explicitly bounded source. |
| Cardinality | Existing default 20, native allowed range 0–500, per-label where specified | Native zero semantics documented for this diagnostic, not generalized to gcx inventory caps. |
| Whole traces, diffs, merged profiles and metric matrices | No newly imposed collection cap on whole objects or arbitrary samples | Explicit pruning, `--max-nodes`, native query limits and transport bounds reported separately. |

For log retrieval, target default order is newest first; `--direction forward` selects oldest first. Preserve backend ordering for unranked trace search and mark it as backend-defined. Inventory order is deterministic; ranked output identifies ranking criteria and uses stable identity tie-breaks where gcx controls ordering. Client-side caps cannot claim to reduce backend work. Never truncate sample points to simulate a series cap or truncate an artifact byte stream to simulate a node bound.

### Retrieve, compare, and drill down

Trace `get` retains full-trace lookup without a default window, so an older ID remains retrievable. Optional lookup bounds constrain the backend search. Span filtering uses the common `--query` expression flag; hierarchy, depth and pruning remain explicit transformation controls. Reject hierarchy/depth controls without their required filter and pruning controls without `--prune`. Results state transformations applied and distinguish an intentionally reduced representation from incomplete backend retrieval.

`diff` retains two-ID semantics and structured differences; it does not need a table merely for symmetry. `baseline` uses the seed's start/end plus 30 minutes of padding by default. `--window` changes padding; explicit `--from/--to` overrides it only when `--window` was not explicitly supplied. Conflicting explicit controls are rejected. Candidate coverage and seed completeness are separate; candidates are not ranked similarity results. IDs from query results can be passed unchanged into all three operations.

```sh
# Proposed target: use returned identifiers without rewriting them.
gcx traces get 0123456789abcdef0123456789abcdef -d traces-uid
gcx traces baseline 0123456789abcdef0123456789abcdef --window 30m --limit 20 -d traces-uid
gcx traces diff 0123456789abcdef0123456789abcdef fedcba9876543210fedcba9876543210 -d traces-uid

gcx profiles exemplars profile --query '{service_name="api"}' --since 1h -d profiles-uid
gcx profiles query --query '{service_name="api"}' --profile-type process_cpu:cpu:nanoseconds:cpu:nanoseconds --profile-id 00000000-0000-4000-8000-000000000001 --since 1h -d profiles-uid
```

Exemplar/anomaly records retain profile type, unit, observation time/window and the exact drilldown identifiers available from the backend. Missing correlation fields are absent, never fabricated. Use the returned window and profile type for a reproducible follow-up rather than recomputing “last hour.” IDs do not authorize the CLI to choose a different datasource silently. Profile/span/trace ID combinations retain native validation restrictions; error guidance names a supported alternative.

A backend limit on selected series is not a cap on exemplars. Report it as an execution bound separately from output coverage. Anomalies retain scores and algorithm-specific meaning, not a common cross-signal severity scale. A valid empty result succeeds; the opt-in `--error-on-empty` applies to query collections that can determine emptiness. Missing capability/configuration remains a failure.

### Artifacts and human presentation

Every advertised codec must encode the command's actual result. Provide meaningful `wide` views for profile query and trace metrics; validate unsupported format/mode combinations before I/O. Graphs, trace trees, raw logs, DOT and pprof remain useful domain views. JSON/YAML and agent output carry the structured contract below, while tables show concise context, warnings and coverage rather than dumping metadata columns into each row.

Output selection must preserve the effective request. Apply profile-ID and stacktrace filters to pprof export. Keep `--max-nodes` explicit and independent of the selected format; target omission sends the existing normal-profile default of 50,000 nodes for every profile representation. Explicit zero requests the native unbounded-node behavior only where supported, with transport bounds still in force. Do not silently change 50,000 nodes to unlimited or 100 because the codec changed. Format-specific backend execution is allowed only with equivalent filtering and disclosed representation semantics.

For `-o pprof`, retain `--pprof-path` and `--pprof-overwrite`; create the complete artifact before reporting success, preserve an existing file on failure, and never overwrite without the explicit option. Agent receipts identify path, format, source request and known size, not binary data in a JSON stream. Retain the existing rejection of `--error-on-empty` for pprof until artifact emptiness can be established correctly. For explicit DOT requests, an unsupported backend must produce an unsupported-format error rather than returning a table as successful DOT. A suggested table retry is guidance, not automatic format substitution.

Explore links have `exact`, `broader`, or `unavailable` fidelity in structured context. Exact links reproduce all material query inputs; broader pivots explicitly list omitted transforms, such as trace pruning. `--share-link` may return a labeled broader pivot; `--open` must not silently open it as an equivalent result. If only a broader pivot exists, return the successful query and the link explanation without auto-opening. Link failure does not invalidate successful data retrieval.

### Catalogs, series, and diagnostics

Apply discovery time controls to metric names, series, profile types, and ranked search where supported. Preserve metric-name local prefix/suffix/contains filters as explicit client-side narrowing, search scores and backend warnings, and profile series label projection. Metadata remains descriptive HELP/type information; cardinality remains a current-ingester measurement with its existing 0–500 limit range and per-label value-limit semantics. Profile `data-range` diagnoses ingestion and queryable bounds; unknown bounds and an ingestion-history indicator must not imply data is currently queryable.

For `list-series`, one common `--query` means one native selector. Expose repeated `--match` only when the adapter can guarantee union semantics. Metrics and logs can map native unions. For profiles, initially support the single-selector contract and reject multiple selectors under the target interface until the supported route proves stable union behavior. Do not silently map a union request onto conjunction. Explicit fan-out/dedup is an alternative requiring a separate cost and completeness design; it is not introduced here. This series gate is independent of the settled profile-label single-selector rule.

Mount applicable metric names, series, and cardinality constructors in the Prometheus datasource tree. Keep documented search naming differences where both paths reuse behavior. Billing stays a signal convenience; its reused constructors inherit the chosen migration. Adaptive operations remain tenant-scoped rather than acquiring misleading datasource flags.

### Adaptive telemetry reads and mutations

Adaptive inventories, recommendations, and configuration reads use format-independent `--limit` for returned items, with structured truncation/coverage. Use common read coverage, translating existing list metadata where full enumeration or known truncation is established and representing unknown explicitly otherwise. Preserve recommendation scores, patterns, volumes, policy schemas, resource identifiers, and observed totals. Ordinary configuration lists target the existing common 50-item management default; recommendation ranking remains domain-specific while observation collections use the 50-item cap above. Negative bounds and unexpected positional arguments are rejected before I/O.

Separate table rollups from result selection. For Adaptive Logs patterns, the presentation control `--table-top` can retain the useful “everything else” aggregate without changing fetched data. A format-independent item cap is `--limit`. A rollup must describe whether its tail covers the entire filtered snapshot or only a capped subset. Metrics recommendation ranking remains explicit and stable across codecs. Snapshot time and observation window are reported when the API supplies them; unknown freshness is marked unknown. These services return generated snapshots, not arbitrary last-hour queries.

Add `get ID` to Logs exemptions and segments using existing backing-client capability. Across all adaptive mutations, validate identity/scope before I/O, retain preview/dry-run and confirmation rules, and produce machine-readable per-item outcomes for batches. Partial failure must not return an indistinguishable overall success; use existing failure taxonomy and retain successful item results. Retry guidance must distinguish idempotent reads from potentially applied writes. Read/merge updates must preserve omitted fields; explicitly empty values that intentionally clear a field remain valid. These rules do not equate applying recommendations, dismissing them, and CRUD updates or introduce cross-service transactions.

For multi-item mutations, preflight the complete supplied request before any writes: parse and validate every item, resolve identity/scope and dependencies, and perform available read-only prerequisite checks. An invalid item fails the batch without writes. Preflight is not a transaction or a guarantee against concurrent changes. A dry run returns `planned` outcomes and proposed changes, explicitly identifying client-side validation; it must not send a mutation to test whether a backend supports dry-run.

After successful preflight and existing confirmation, execute independent items with bounded concurrency. Continue independent items after a failure; mark dependent items `not_attempted` when their prerequisite fails. Return each selected item's stable identity, action and outcome: `succeeded`, `unchanged`, `failed`, `unknown`, or `not_attempted`. A timeout after sending a write is `unknown` unless read-back or an idempotency guarantee establishes its effect. Do not blindly retry it. The result gives concrete read-back/reconciliation guidance and preserves confirmed successes; summary counts cover every selected item exactly once.

Mixed success and failure/uncertainty uses existing partial-failure exit 4; a wholly failed or uncertain batch uses the applicable general/auth failure exit; cancellation keeps exit 5 and reports known outcomes where possible. All-success/unchanged is exit 0. A failure-free preview returns planned items with `dry_run: true`, never claims changes succeeded, and does not reuse a stale preview as execution authorization. Listing a capped recommendation set must not silently define an “apply all” scope: application resolves and states its target set independently, requiring explicit all-selection where the command supports it.

Profiles `adaptive` must report an unsupported capability with a nonzero structured error until implemented. Its current success exit after an unavailable message misleads agents.

### Discover names, then values

The proposed target commands make the requested operation explicit:

```sh
# Proposed: names and values are separate operations.
gcx metrics list-labels -d metrics-uid
gcx logs list-label-values job -d logs-uid
gcx traces list-label-values resource.service.name -d traces-uid
gcx profiles list-label-values service_name -d profiles-uid
```

`list-labels` accepts no positional arguments. `list-label-values` requires exactly one nonempty label identifier. These names follow the [enumeration naming convention](../design/command-naming.md). The same constructors are mounted under `gcx datasources prometheus`, `loki`, `tempo`, and `pyroscope`; those paths return the same data and accept the same flags as their signal counterparts.

Human output uses a names or values table with trace scope/type columns where meaningful. Machine formats use the common schema below. Agent mode changes presentation defaults, not the selected window, limit, filters, or available capabilities.

### Carry the investigation window across signals

Discovery uses the shared one-hour window and time grammar defined above. Validate that the start precedes the end and pass explicit absolute bounds to the backend.

```sh
# Proposed: the same explicit window for discovery across signals.
gcx metrics list-label-values job --since 6h -d metrics-uid
gcx logs list-label-values job --since 6h -d logs-uid
gcx traces list-label-values resource.service.name --since 6h -d traces-uid
gcx profiles list-label-values service_name --since 6h -d profiles-uid
```

Use absolute `--from`/`--to` values when several calls must share exactly the same interval. The output reports the resolved bounds and whether they came from the default or explicit flags. Bounds constrain the backend's discovery operation; they do not prove a returned label occurred in a sample strictly inside that interval. Human help and machine metadata describe approximate metadata semantics where applicable.

The same one-hour default applies to window-based execution queries. Instant evaluation, whole-object lookup, seed-centered baseline and timeless/snapshot operations retain their explicit exceptions above. Verify mode mapping through Grafana transport rather than inferring it from a legacy request field.

### Scope discovery consistently

Every new label discovery command accepts one `--query` expression. Its purpose is the same across signals: scope discovery server-side. The accepted language remains native and narrower than a general expression where the endpoint requires that.

```sh
# Proposed: one discovery expression, interpreted by each backend.
gcx metrics list-label-values instance --query '{job="api"}' --since 1h -d metrics-uid
gcx logs list-label-values instance --query '{job="api"}' --since 1h -d logs-uid
gcx traces list-label-values resource.service.name --query '{ span.http.status_code >= 500 }' --since 1h -d traces-uid
gcx profiles list-label-values service_name --query '{namespace="payments"}' --since 1h -d profiles-uid
```

Metrics label discovery additionally accepts repeatable `--match` with native selector-union semantics. Profile label endpoints combine their matcher-array entries conjunctively, so profiles exposes only the common single `--query` form; it does not claim a union capability. Loki series supports unions, but Loki label discovery does not. `--query` and `--match` are mutually exclusive, including when a supplied value is empty. Repeated `--query` is rejected instead of silently keeping the last occurrence. Multiple constraints inside one selector retain their native conjunction semantics; repeated selectors retain their documented union semantics. Label discovery does not emulate missing unions by issuing several requests.

Metrics retains `--metric` as a narrowing modifier applied to every selector. Traces retains `--scope`, including supported intrinsic discovery, and returns identifiers suitable for passing directly to `list-label-values`. Intrinsic identifiers require their own mapping; they are not formed by blindly prepending `intrinsic.`. A qualified identifier and a contradictory `--scope` are rejected.

Loki label discovery accepts a stream selector, not log pipelines. Tempo discovery supports a restricted spanset filter rather than arbitrary TraceQL. Its backend can tolerate invalid or incomplete filters and return broader results. The new commands therefore require strict validation of the supported discovery grammar: reject malformed expressions and valid language features unsupported by discovery before I/O. A successful parse as general TraceQL is insufficient. The validator must verify the entire expression, not accept a valid prefix or silently remove an unsupported predicate. Tempo support cannot ship until this is demonstrated against the actual discovery route. Enumeration coverage describes result bounds and must never excuse uncertain filter enforcement.

### Understand result bounds

For label discovery, the default is 100 returned items. `--limit N` is a nonnegative cap on normalized items across the whole response, including all trace scopes. `--limit 0` removes the gcx item cap; it never promises that a backend default, scan budget, or transport bound disappears.

Push limits to the backend only where the endpoint supports them. Where possible, fetch one extra item to establish that more results exist. Normalize and deduplicate names by actionable identifier and values by logical type plus lossless representation before the final item cap; sort that returned inventory deterministically. Sorting does not promise that a capped backend returns the same subset on every call. A per-scope Tempo bound must be reconciled with the global item cap. Where an endpoint has no limit, fetch its bounded HTTP response before applying the item cap; disclose that the limit reduces output rather than backend work. Do not replace label discovery with expensive series enumeration merely to simulate a missing feature.

A result identifies whether it is complete under the endpoint's enumeration semantics, known partial, or of unknown completeness. Fewer than N items alone is not proof of completeness. Unknown totals remain unknown. Exceeding the HTTP response safety bound is an error, not a successful truncated JSON document. No automatic all-results scan or fan-out is introduced.

### Common structured read contract

All structured telemetry reads use one small metadata contract with operation-specific payloads. A discriminated `type` and string `schema_version` identify the shape. Common fields are `signal`, the effective `datasource_uid` for datasource reads or non-secret tenant context for adaptive reads, `status`, and `warnings`. Historical reads include `window`; point evaluations use `evaluation_time`. The effective `filter`, mode, resolution, ordering/ranking, and execution bounds appear only when applicable. Values distinguish explicit inputs from defaults; unknown backend decisions are marked unknown, not guessed.

Collections use `items`; whole objects, matrices, merged profiles, comparisons and aggregate diagnostics use `data` with their native typed structure. Shared metadata does not flatten spans, samples, stack trees, scores or units. A matrix preserves series and sample arrays; a trace preserves its trace schema. A missing object produces not-found failure, not an empty collection. Do not attach an enumeration coverage claim to a whole object: backend incompleteness belongs in warnings and status, while deliberate transforms are recorded separately.

`status` is `ok` or `partial` for successful reads; known backend data loss or a deliberate item cap makes it partial. Unknown enumeration coverage is expressed by `coverage.status: unknown` without asserting known data loss. `coverage` describes the returned collection only; `execution` describes upstream scan/series/node budgets, ranking and transformations that affect the analytical meaning. These fields cannot substitute for each other. Warnings use stable codes plus readable messages and optional recovery hints; no warning contains credentials or a raw authenticated URL.

Human codecs show a compact context line and actionable warnings/coverage. Structured codecs serialize the same result once; stderr never contains the only copy of a material warning for machine consumers. Agent spill refers to the complete structured document. Implement existing `ListEnvelope` selection for collections; preserve the shared siblings during item selection. For whole-object reads, `--json` field paths are relative to `data`; selection replaces only `data` and preserves context/status/warnings. Field discovery lists domain payload fields. Unknown paths fail with usage guidance rather than silently returning an empty object. Implement and test this alongside existing collection selection, without pretending the object is a collection. `--jq` remains an explicitly requested arbitrary transformation.

```json
{
  "type": "telemetry_trace_search",
  "schema_version": "1",
  "signal": "traces",
  "datasource_uid": "traces-uid",
  "status": "partial",
  "window": {"from": "2026-10-12T09:00:00Z", "to": "2026-10-12T10:00:00Z", "source": "explicit", "semantics": "backend_search"},
  "filter": {"query": "{ resource.service.name = \"api\" }", "language": "traceql"},
  "mode": "search",
  "items": [{"traceID": "0123456789abcdef0123456789abcdef"}],
  "coverage": {"status": "partial", "returned": 1, "limit": 1, "unit": "traces", "reason": "gcx_item_cap"},
  "warnings": [{"code": "RESULT_CAPPED", "message": "Additional matching traces were observed; increase --limit to return more."}]
}
```

This illustrative search item omits optional native fields for brevity; implementations preserve them. The cap warning requires observed evidence, such as an extra returned trace, not merely a response whose length equals the cap. One envelope serves human, JSON, YAML and agent rendering; raw/artifact output has its explicit contract above.

### Label discovery payload

The following is the proposed version-one envelope. Exact field names are proposal details, not a shipped interface. JSON, YAML, and agent rendering derive from the same result; human codecs display the same coverage information in a concise form.

```json
{
  "type": "telemetry_label_values",
  "schema_version": "1",
  "signal": "traces",
  "datasource_uid": "traces-uid",
  "status": "ok",
  "warnings": [],
  "label": "span.http.status_code",
  "window": {
    "from": "2026-10-12T09:00:00Z",
    "to": "2026-10-12T10:00:00Z",
    "source": "explicit",
    "semantics": "backend_metadata"
  },
  "filter": {
    "query": "{ resource.service.name = \"checkout\" }",
    "language": "traceql_discovery"
  },
  "items": [
    { "value": "200", "value_type": "int" },
    { "value": "500", "value_type": "int" }
  ],
  "coverage": {
    "status": "unknown",
    "returned": 2,
    "limit": 100,
    "unit": "label_values",
    "reason": "backend_did_not_establish_exhaustion"
  }
}
```

For names, `type` is `telemetry_label_names`, `label` is absent, and each item contains `name` with optional `scope`. A trace name is an actionable identifier, such as `resource.service.name`, not an ambiguous flattened name. Other signals use their original label names. Values preserve the backend's JSON representation and carry its logical type; numeric strings must not be silently coerced or lose precision. String-label backends use `value_type: "string"`. Backend type vocabulary and unknown future types remain representable. Empty results use `items: []`. Backend decoding and subsequent output normalization must preserve JSON numbers losslessly: the current generic map conversion can turn numeric values into `float64`, so passing a typed envelope through it is not sufficient. Use lossless number decoding through the shared selection/encoding path and test its existing consumers before changing that helper.

`filter` is absent only when no narrowing input was requested. Record `query` or the union `match` array, plus `metric` and `scope` when supplied; these modifiers remain visible even when the result is empty. The adapter must test that this recorded context matches the request it sends. Union filtering does not invent one expression. `coverage.status` is `complete`, `partial`, or `unknown`. A `complete` result means the adapter has established exhaustion for this operation under the backend contract, not that all historical telemetry or every matching sample was examined. `partial` requires observed evidence of a bound; its reason distinguishes a gcx cap from a reported backend cap. Include a total only when observed, and a continuation only when it is actionable. Unknown backend truncation is never relabeled as known completeness.

This envelope needs tri-state coverage. The existing [list metadata helper](../../internal/output/listmeta.go) treats absent `list_meta` as completeness and cannot express unknown coverage unchanged. Use its proven truncation calculations where appropriate, but do not attach a contradictory `list_meta` contract to this result. Implement the existing [ListEnvelope interface](../../internal/output/field_select.go) with `ListItemsKey()` returning `items`; its item-level `--json` selection preserves sibling metadata, including window and coverage. Verify this with the new result types rather than introducing another output-selection mechanism. `--jq` remains an explicit user transformation that can discard metadata. Agent spill receipts must lead to the complete envelope, not an independently truncated payload.

### Other inventory result shapes

Series, metric-name inventories, and profile-type inventories use the same versioned context/coverage envelope with distinct operation types (`telemetry_series`, `telemetry_metric_names`, `telemetry_profile_types`) and domain-specific `items`. Preserve series label sets, metric-name records (`{"name": "..."}`), and profile-type identifiers/units; do not coerce them into label-value objects. Their 100-item default counts one returned series label set, metric name, or profile type respectively. Deduplicate by the operation's identity before applying the cap, retaining uncertainty about backend exhaustion. All inventory items are objects, including names, because existing item-field selection expects object records. Implement `ListEnvelope` for each so field selection and spill preserve context and coverage.

Descriptive metric metadata targets a 100-record cap, where one record is a metric name plus its metadata entry; retain multiple distinct type/help/unit entries for the same name. It uses a descriptive operation type and coverage, with no invented window. Ranked metric search retains its domain items, scores, warnings and incomplete-stream evidence inside the common read contract. Management inventories likewise preserve domain item schemas inside it. Translate existing list metadata into coverage without inventing completeness; migration selects compatibility wrappers for consumers of legacy envelopes.

### Datasource selection and capability failures

Use common selection precedence: explicit `--datasource/-d`, configured signal datasource, then existing discovery. An explicit selection failure never triggers a different datasource. Reuse already-known type information; an explicit UID must not require a new metadata-read permission when its query route is authorized. A confirmed type mismatch is a usage failure; an unavailable metadata check is not proof of a mismatch. Validate through the operation route and retain its authentication/authorization errors.

Capabilities come from command/adapter support and verified backend responses. Help and the agent command catalog identify supported modes, formats, filter grammar and limit units. A failed request must not be retried with removed filters, bounds, identifiers, or a different format. Equivalent transport fallback is permitted only when the full request contract survives; authentication/authorization failures do not trigger capability fallback. Capability absence and transient service failure remain distinguishable structured errors.

### Errors and empty results

Validate explicit empty or whitespace-only identifiers, selectors, durations, and datasource flags before configuration loading or network requests. Omission may select a default; an explicitly empty scope-affecting value must not broaden the request. Structural validation is shared. Backend-language validation remains the responsibility of the appropriate adapter, with strict validation of the supported discovery grammar required for tolerant endpoints.

Use the existing structured error protocol and exit taxonomy: invalid CLI input returns usage failure, authentication and authorization retain their meaning, and backend unavailability does not become an empty successful inventory. Unsupported capabilities must not be retried by dropping filters or bounds. A valid empty inventory succeeds with an empty array and the requested window/coverage. Existing query `--error-on-empty` behavior remains available where the result type supports it; this RFC does not make absence of telemetry a default error.

### Partial results, completeness assertions, and deadlines

Usable reads with known partial results or unknown coverage succeed by default and carry the common status, warnings and coverage. `--require-complete` changes that policy: success requires established completeness for the requested operation; partial or unknown coverage fails while preserving the usable payload in the existing structured error/partial-result protocol. This is a general assertion failure (exit 1), not exit 4 merely because a collection is capped. The latter remains reserved for mixed operation outcomes. Emit one valid document rather than a success document followed by a separate error document.

Expose the assertion on reads with an enumerable completeness contract. For analytical queries, the adapter must account for known upstream bounds as well as the output list; complete enumeration of returned top groups does not prove an unbounded computation. Unknown upstream completeness cannot satisfy the assertion. Whole-object and diagnostic operations without a meaningful completeness contract reject the flag before I/O and explain the supported alternative. Deliberate transforms describe the requested representation; they cannot hide backend partial retrieval.

`--error-on-empty` asserts that the result contains data; `--require-complete` asserts that its completeness is established. They are independent: a complete empty inventory passes only the latter, and a nonempty partial inventory passes only the former. Combining them requires both conditions. Neither triggers automatic pagination, fan-out, retries, wider windows or removal of server bounds.

The invocation context/deadline covers datasource resolution and all constituent requests, including baseline lookup/search and artifact production; do not reset a fresh budget for each internal call. Preserve existing timeout controls rather than inventing per-signal flags. A transport/deadline error is an error unless the backend delivered a valid usable partial response with explicit semantics; incomplete JSON is never a successful partial result. Mutation timeouts use the uncertain-write contract above.

## Proposed architecture

### Components and ownership

The command constructors remain in `internal/datasources/{kind}` and are mounted by the existing signal providers. The core CLI owners maintain shared input/time/datasource validation, output and error contracts, common read context/envelopes, discovery options, integration with existing output selection, and conformance tests by operation class. Signal owners retain query modes, domain result types, artifact generation, and adaptive mutation semantics. Signal owners maintain expression validation, backend request translation, scope/type preservation, and capability tests against their supported versions.

`internal/query/{kind}` continues owning HTTP paths and backend wire types. It translates normalized operation requests into native requests and returns result data plus evidence about bounds. Each adapter returns the effective operation context and domain result; shared code assembles the read envelope from that evidence; it must not infer completeness from a response length alone. `internal/output` retains ownership of codecs and agent protocols. Adaptive providers continue to own tenant management clients and mutation execution; common list/validation/output helpers provide consistency without routing management through datasource clients. A dedicated generic provider framework, API-generated command tree, or persistent capability registry is unnecessary for these four concrete adapters.

### Shared execution boundary

Build an immutable effective request after local validation: operation, datasource or tenant identity, query/identifier inputs, absolute time semantics, mode, bounds, and output requirements. Resolve supported capabilities without speculative writes. Backend clients execute the request; response adapters preserve native payloads and return warnings, observed limits and completeness evidence. Codecs receive that result and cannot select a different window or discard filters.

Keep capability declarations beside command constructors and use them to register only valid flags/formats and populate the existing agent command catalog. This is a small shared contract, not a new remotely synchronized capability service. Common preflight checks cover expression source conflicts, prerequisite flags, format compatibility and numeric/time validation. Signal adapters own grammar and native combination rules. Targeted route fixtures verify that every accepted input reaches its intended request field.

Reuse [mutation and artifact result families](../../internal/output/mutation.go) instead of applying a read envelope to writes or files. The existing batch family counts successes/skips and enumerates failures; it does not currently represent every item or uncertain writes. Extend it through a versioned contract where required, and update shared consumers before telemetry providers adopt it. The artifact receipt also needs a compatible/versioned extension: per-file known byte size and an effective source-request context, preserving existing path/format/failure fields. Shared receipt serializers, agent consumers and spill handling must be tested with both existing and extended types. Preserve the existing structured error protocol and single-document output on failures.

### Discovery request flow

1. Parse the operation and validate flags, arguments, filter form, and time range before I/O. Resolve the clock once and retain the normalized request.
2. Resolve the datasource using the common precedence. Validate any already-known type and preserve permission errors without switching identities.
3. Execute the signal-specific discovery operation with the explicit window, filter, and supported backend bound. Do not remove constraints on failure.
4. Normalize identifiers and typed values, apply any remaining global cap, and attach evidence-based coverage.
5. Encode through the selected codec. Field selection and spill retain the context needed to interpret the result.

### Label adapter mapping and integration boundaries

| Signal | Target label discovery mapping | Important boundary |
| --- | --- | --- |
| Metrics | A single `--query` becomes one `match[]`; repeated `--match` stays a union; send bounds and supported limit | Grafana's Prometheus resource handler must forward every parameter; do not assume a latest-doc capability exists on all versions. |
| Logs | One selector becomes `query`; send bounds; apply item cap locally when no server limit exists | Grafana's Loki resource route translates to native label endpoints. Series has different selector cardinality. |
| Traces | Send the supported filter as `q`, bounds in seconds, and supported discovery bounds; preserve scopes/types | Proxy forwarding and Tempo's tolerant discovery grammar require explicit validation coverage. |
| Profiles | Send one selector as a matcher with millisecond bounds; preserve labels and string values | LabelNames/LabelValues matcher-array entries are conjunctive, not a native union. Verify the selected single-selector request through the proxy/Connect route. |

The named mappings are supported by public API documentation and client source. Source inspection of a Grafana release also supports resource query-string forwarding, but does not establish the minimum compatible Grafana/plugin/backend combination. No automatic fallback may broaden an operation. A capability unavailable on a supported target needs an explicit unsupported result or a proven equivalent route, with auth failures kept distinct from capability failures.

### Execution and management API mappings

| Operation | Existing integration boundary | Adapter obligation for target UX |
| --- | --- | --- |
| PromQL and LogQL evaluation / log lines | Shared Grafana query transport with legacy fallback | Carry explicit mode, absolute time and resolution; distinguish log line max from matrix semantics; equivalent fallback must preserve inputs. |
| Tempo search and metrics | Datasource proxy to search and metrics query/range endpoints | Preserve TraceQL, mode and bounds; distinguish point/range/window-aggregate semantics and backend default caps. |
| Trace get, diff and baseline | Proxy to trace-by-ID/diff; baseline orchestrates seed lookup and search | Preserve ID/filter/pruning controls; separate seed completeness from candidate coverage; one normalized clock/window policy across constituent reads. |
| Profile merge and artifacts | Proxy/Connect merge-stack, span-profile and merge-profile requests | Map type/IDs/stack selector/time/node bound for every valid output combination; validate produced artifact before receipt. |
| Profile metrics, exemplars and anomalies | SelectSeries, heatmap/exemplar and anomaly requests | Keep series-selection budget, ranking, resolution and returned-item cap distinct; retain correlation IDs and units. |
| Metric names/search/metadata/cardinality | Prometheus/Mimir datasource routes | Send only supported historical bounds; retain metadata and cardinality semantics, search stream warnings and feature-gate errors. |
| Profile type and data-range catalogs | ProfileTypes and GetProfileStats | Apply historical window to types only; preserve unknown diagnostic bounds and distinguish ingestion history from current queryability. |
| Adaptive reads and mutations | Existing tenant-authenticated management clients | Preserve resource scope, snapshot semantics, credentials and mutation outcome evidence; no datasource transport or fabricated historical filter. |

Native API support is necessary but insufficient: the Grafana resource/proxy route and installed plugin must preserve each relevant field. Test backend-version differences explicitly. No adapter may silently broaden a request to compensate for a route that cannot implement the target contract.

## Backwards compatibility and migration

This section is deliberately separate from the target design. The chosen direction is to preserve compatibility invocations and offer opt-in adoption of the common read schema. Exact aliases, shims, `--expr`/`--query`, `series`/`list-series`, `--top` replacements, default query modes/windows, read envelopes, deprecation dates, default switching, and major-version timing require a separate discussion with signal teams. A plain Cobra alias is insufficient when the legacy operation combines names/values or has different output/default semantics; a compatibility wrapper may be required.

The review must cover existing `labels`/`tags`, `--label`, selectors, JSON/YAML consumers, discovery time defaults, series defaults, accepted no-op flags, and query limit semantics. Inventory supported invocations before selecting the migration. The rollout must preserve released commands, default result shapes, and flags under the current compatibility invariant unless an explicit exception is approved. The target design above does not select such an exception.

## Drawbacks and alternatives

Explicit names and values operations enlarge the command catalog, while the common envelope adds metadata tokens. Those costs must be measured against fewer help lookups, mistaken calls, and ambiguous results. Typed value objects are more verbose than bare strings but preserve trace meaning. Strict filter validation adds parser maintenance and a version-compatibility obligation.

Keeping `labels --label` would avoid a naming migration but retain implicit operation switching. Flattening backend results into string arrays would be compact but lose trace type and scope. Exposing each backend's API vocabulary unchanged would simplify adapters while transferring backend-specific learning and error recovery to every user. A universal query language or union emulation would introduce translation and execution costs beyond the consistency benefit sought here.

The proposal adopts a common operation contract with explicit capabilities, rather than reducing every signal to the least capable API. Common windows apply to historical discovery and window-based queries; explicit point evaluation and domain operations retain their distinct semantics.

## Validation and delivery

### Correctness acceptance

Use table-driven command tests through the real constructors and targeted HTTP fixtures. Cover each distinct mapping once and verify both mounts share it. The contract suite checks names/values dispatch, explicit-empty rejection before I/O, time normalization, filter cardinality and combination, known datasource mismatches, all advertised formats, typed identifiers/values, empty arrays, and coverage preservation across JSON, YAML, field selection, and agent spill.

Bound tests cover zero, one, exactly N, N+1, backend partial responses, unknown exhaustion, and oversized HTTP bodies. Include Tempo's per-scope limit and a backend without server limits. Include numeric values beyond the exact integer range of `float64` (for example 9007199254740993), numeric strings, booleans, and unknown trace types across JSON, `--json value`, YAML, and agent spill. Test new schema stability separately from the compatibility fixtures selected during migration.

Test completeness assertions against complete-empty, nonempty-partial, unknown-coverage and upstream-budget cases, including combined empty/completeness assertions and single-document failure output. Test mutation preflight makes no writes on any invalid item, independent execution continues after failures, dependencies are not attempted, ambiguous timeouts remain unknown, and every selected item is counted once.

Beyond discovery, conformance fixtures must prove that pprof retains every accepted filter, advertised wide codecs work, irrelevant/conflicting step flags fail before I/O, and exemplar/series/result bounds remain distinct. Exercise adaptive list truncation in every format, table rollups with capped input, intentional field clearing, preview/apply consistency, mixed batch outcomes, and unsupported Profiles adaptive failure. Verify retrieval transformations, baseline windows, diff capability errors, and profile data-range unknowns. Both signal and datasource mounts must satisfy the same applicable tests.

Live integration must exercise Grafana-facing routes, not only direct backends. Record Grafana/plugin/backend versions and capability results. Verify explicit windows, selector forwarding, intrinsic tag values, malformed Tempo filters, and profile matcher combinations. Query checks must establish actual Loki mode and limit translation, Tempo zero-limit behavior, Explore fidelity, and permission-sensitive datasource validation. Passing a fixture does not establish these integrations.

### Human and agent evaluation

Use matched baseline/proposal tasks spanning the complete surface:

| Journey | Correct completion requires |
| --- | --- |
| Discover and transfer | Names → values → series/query across signals, widening an empty window, interpreting typed values and uncertain coverage. |
| Investigate a spike | Carry an absolute window across metric/log/trace/profile queries, choose evaluation step only where valid, and interpret bounded results. |
| Retrieve and compare | Retrieve an older trace by ID, obtain seed-centered baseline candidates, interpret diff and partial-seed evidence without claiming ranked or exhaustive matches. |
| Profile drilldown and export | Find profile type, select exemplar/anomaly, preserve IDs/stack filters in the next query and pprof export, and distinguish an unavailable capability from no data. |
| Catalog and management | Interpret current-state cardinality and adaptive snapshots, detect capped lists, preview a recommendation change, and correctly report a mixed batch outcome in a controlled fixture. |

Include malformed input, missing metadata permission with valid query access, unsupported formats, and cross-signal transfer variants. Use fixtures for mutations; this RFC does not authorize live changes.

For humans, measure correct completion first, then meaningful actions to completion, help visits, corrections/retries, and elapsed time. Counterbalance interface and signal order to separate transfer benefits from learning effects. For agents, hold model, instructions, tool setup, fixtures, and budget fixed; measure correct task success and total input/output tokens across the whole attempt, including tool responses and failed calls. Report failures separately and cost per successful task without discarding failed-attempt cost. Repeat runs and report uncertainty rather than choosing a winning anecdote.

The target is non-regressing task success with fewer human actions and lower agent spend, assessed per task and signal as well as overall. Establish baseline variance before agreeing to numerical release thresholds. A lower token count with more incorrect answers fails the objective. Contract bug fixes do not require proving a cosmetic naming preference through an eval.

### Keeping the contract consistent

Signal teams own backend semantics and implementation; core CLI owners own shared behavioral contracts. New or changed commands must identify their operation class, reuse the applicable options/output/error helpers, and pass the corresponding conformance cases. A justified exception states the backend or task constraint, observable UX difference, and recovery guidance in help and agent metadata. It must not silently overload an existing flag. Review exceptions against task success and action/token cost rather than requiring visual symmetry.

Add cross-signal transfer tasks to the eval suite when a new capability creates a comparable workflow. Review per-signal regressions as well as aggregate improvements; a team cannot hide a lower task success rate behind another signal's token savings. Backend-specific features remain welcome when their constraints are explicit. This gives teams freedom over domain capabilities while preventing every integration from redefining common CLI behavior.

### Delivery slices

Deliver by behavior risk rather than label-first scope: (1) repair ignored filters, unusable codecs, false-success capability errors, and misleading bounds; (2) add shared validation and list/coverage contracts across catalogs, query collections, and adaptive reads; (3) implement the label redesign and historical discovery defaults across adapters; (4) complete retrieval/drilldown/link fidelity and management outcome conformance. Each slice includes source mapping, fixtures, applicable live route checks, help, and command metadata. Evaluate complete journeys before switching preferred documentation and skills under the separately agreed migration plan.

The design covers the whole contract; delivery slices do not defer designing difficult semantics. Update the architecture's claim of identical signal flags to describe equivalent operations with explicit capabilities. Keep generated help, command metadata, examples, and portable skills consistent with the chosen migration stage.

## Unresolved questions

### Deferred migration decisions

Compatibility aliases or wrappers, opt-in mechanics, release sequencing, and retirement timing remain for a separate discussion with signal teams. The target design fixes the discovery window at one hour, the label-result cap at 100, and strict discovery-filter validation as a release requirement.

### Gap-to-acceptance checklist

Every audited inconsistency has an observable target and a release check. Passing this checklist requires implementation evidence; the source audit alone does not satisfy it.

| Audited gap | Target remedy | Acceptance evidence |
| --- | --- | --- |
| Vocabulary, missing mounts/read verbs | Canonical query input, explicit discovery verbs, shared datasource constructors, Logs configuration get | Real command-tree tests and both-mount request equality. |
| Missing/default time controls and misleading mode translation | Common absolute window normalization, last-hour default, explicit instant | Frozen-clock fixtures and Grafana-facing mode/time round trips. |
| Ignored/conflicting flags and explicit-empty broadening | Preflight option/capability validation, strict tolerant-endpoint grammar | Zero-I/O negative tests for every relevant operation. |
| Datasource validation divergence | Common precedence without extra mandatory metadata permission | Query-authorized/metadata-denied role test and confirmed mismatch tests. |
| Series union differs by route | Single selector unless union is proven; native union only where guaranteed | Versioned series-route fixtures and live supported-route tests. |
| Overloaded top/limit and hidden truncation | Separate ranking, presentation, item caps and execution bounds | Same selected data across codecs, boundary tests, truthful coverage/units. |
| Unsupported wide and format-dependent profile requests | Working codecs, equivalent filters/node policy, no silent DOT substitution | Format matrix, wire-field assertions, parsed artifact validation. |
| Missing empty assertions / unavailable reported success | Consistent collection assertions and structured capability failure | Empty, not-found, auth, unsupported and service-error cases remain distinct. |
| Incomplete or absent Explore links | Exact/broader/unavailable fidelity and no silent broad auto-open | URL round-trip checks for query/type/IDs/window/step/transforms. |
| Ambiguous catalog/snapshot semantics | Descriptive, current-state and generated-snapshot context | No fake windows; unknown freshness/bounds preserved in output. |
| Loss of scope/types/precision/context | Domain payloads with shared read context and lossless selection | Trace identifiers, numeric precision, JSON/YAML/selection/spill fixtures. |
| Adaptive input/outcome inconsistency | Strict preflight, shared mutation protocol and preserved clearing semantics | No-write validation tests and controlled multi-item outcome tests. |

### Technical release gates

- Establish minimum supported route/backend versions through source and live tests. Record unsupported combinations; never infer that every Grafana 12+ installation exposes a latest API feature.
- Verify the single-selector profile mapping and Tempo actionable identifiers, including intrinsic and quoted names, through golden request/response fixtures and live round trips. Profile label matcher-array conjunction is source-confirmed; it must not become an accidental union claim in help or tests. Separately verify the Series read-path difference and reject unsupported multi-selector unions before I/O.
- Choose a maintainable validator for the exact discovery grammar and demonstrate that malformed and unsupported expressions cannot silently broaden results. Regex checks alone are insufficient evidence.
- Verify datasource metadata permissions versus query permissions. A type check must not prevent an explicit-UID operation that previously worked solely because metadata read permission is absent; test the selected behavior with representative roles.
- Validate the new envelope's ListEnvelope field selection and tri-state coverage against the existing output contract. Reuse metadata-preserving selection; do not interpret the binary ListMeta truncation contract as support for unknown coverage.
- Prove faithful Explore links for each supported query combination, and mark unsupported combinations explicitly. No UI fidelity or live integration result is claimed by this RFC draft.

## Terminology

**Discovery** enumerates label names, label values, metric names, profile types, or series used to construct a query. **Snapshot** is the latest backend-generated observation set, with freshness distinct from an arbitrary query window. **Result cap** bounds returned items; **ranking** chooses or orders items by a domain criterion; **presentation rollup** summarizes items without changing data acquisition. **Signal** means metrics, logs, traces, or profiles. **Selector union** includes results matching any supplied selector; predicates inside an individual selector retain backend semantics. **Coverage** describes whether the returned inventory is complete under the backend operation's contract, known partial, or uncertain. **Task success** means correctly completing the investigation, not merely receiving exit code zero.
