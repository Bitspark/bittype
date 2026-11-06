# Invocation observations

The `go` package imports Go test events and associates them with independently
discovered top-level test declarations. It is development tooling, not a runtime
capability provider or an assessment of language completeness.

Three things remain separate:

- A declaration has an owning source identity, tested package, internal/external
  package variant, filename and exact top-level test name.
- The reporting host associates authored capability cases with that declaration.
- An invocation supplies observations for exact emitted names, including separate
  occurrences when a name runs repeatedly.

One passing function referenced by several case annotations remains one observed
function execution. Passing negative tests are passes; the importer does not
interpret their operation-level expectations.

## Capture and publication

`../capture.py` uses Python's standard library. `run` receives the existing
`go test -json` command explicitly, preserves its selection/options and returns
its process failure. It records only requested environment settings. For example:

```sh
python3 tools/provider-matrix/capture.py run \
  --root "$PWD" --output "$RUNNER_TEMP/observations/contracts" \
  --command-id contracts -- go test -json -count=1 ./...
```

Each invocation directory contains exactly these capture inputs:

```text
invocation.json   versioned invocation, source/build context and exit status
stdout.jsonl      unmodified Go test/build JSON output
stderr.log        separate command stderr, including non-JSON early errors
```

The initial invocation record is incomplete. The helper records completed capture
only after the child process exits and both output files close. Interrupted or
malformed captures cannot supply applicable successes. Distinct command IDs use
distinct directories; existing captures are never overwritten.

The reporting host captures its current source context independently:

```sh
python3 tools/provider-matrix/capture.py context \
  --root "$PWD" --output "$RUNNER_TEMP/current-context.json"
```

For a race invocation, also pass `--build-flag=-race` and load declarations under
that build flag. Do not adopt the invocation's context as the declaration
snapshot. Installed consumer modules have a separate harness/module identity;
their same-spelled test names do not become repository test declarations.

Context includes the actual checkout revision, module selections/replacements,
Go version and build selection. Git's existing source-state information discloses
formatter/dependency changes and untracked inputs. Modified or incomplete context
does not establish the recorded clean revision. Local replacements and overlays
remain unqualified; there is no checksum or source-equivalence mechanism.

The existing CI jobs upload captures and reports together, including after test
failure. Results are artifacts of that invocation, not committed documentation.
Static reports contain declarations, authored scopes/cases and the CI route.
They contain no execution results, generation timestamp or current HEAD field.

## Import and association

`Load` reads the invocation file and its fixed sibling outputs. `Decode` accepts
the same invocation plus readers. The schema cannot choose additional paths.
`Associate` receives a separately supplied `Snapshot` and produces explicit
matched, mismatched, insufficient, ambiguous or unassociated relationships.

Both internal and external tests can have the same name in the tested package.
Such observations remain ambiguous. Event order, `_test` package spelling and
generated name suffixes never resolve source ownership. Child events remain
subordinate execution details; they do not introduce declared capability cases.

The decoder preserves named pass/fail/skip outcomes, package outcomes, build
events and incomplete executions separately. It never reads outcomes from
`Output` text. A package pass does not supply a missing named observation. A
parent pass does not turn skipped children into successful assertions. A named
pass before another test fails remains visible beside the enclosing failure.
Repeated occurrences never replace earlier failures.

Profiles and assertion labels remain authored associations. Ordinary Go events
do not prove a test selected its annotated profile, invoked a prepared compiler,
or exercised every branch. Exact invocation arguments remain visible; the tool
does not infer fresh execution from timestamps or cached-output text. CI preserves
`-count=1` for its ordinary and race commands.

The default import bounds are 64 MiB stdout, 16 MiB stderr and 200,000 events.
Explicit limits cannot exceed 1 GiB per stream or 10,000,000 events. Invocation
metadata is limited to 1 MiB. Cancellation is checked between reads; readers and
in-process work do not provide hard process isolation. Raw artifacts remain
available when parsing, source discovery or report generation fails.
