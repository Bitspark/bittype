# bittype/ACX-014: The bounded writer for a child process had one copy per caller

**Status:** Corrected, with limits; state of 13 November 2026. One package,
`internal/processoutput`, holds the writer, the three callers use it and no
copy is left. The [limits](#limits) say what the correction does not
establish. No independent review has read it yet.
**Scope:** The collection of a child process's output stream and error stream
against a bound, on `main` at `3998db5fe4f76942ba1c99870a6be8cc160e316f`:
`boundedBuffer` in
[`internal/rustsem/compiler.go`](https://github.com/Bitspark/bittype/blob/3998db5fe4f76942ba1c99870a6be8cc160e316f/internal/rustsem/compiler.go#L46-L68)
and in
[`internal/typescriptsem/compiler.go`](https://github.com/Bitspark/bittype/blob/3998db5fe4f76942ba1c99870a6be8cc160e316f/internal/typescriptsem/compiler.go#L43-L65),
and `outputBudget` with `boundedOutput` in
[`tools/provider-matrix/report/go/load.go`](https://github.com/Bitspark/bittype/blob/3998db5fe4f76942ba1c99870a6be8cc160e316f/tools/provider-matrix/report/go/load.go#L124-L149).
Identified on 13 November 2026. Tracked in
[#299](https://github.com/Bitspark/bittype/issues/299).

## Rules involved

[AC-11](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-11)
places a mechanism by the knowledge that its contract, algorithm and
guarantees require. This mechanism holds a number of bytes, refuses the write
that would pass it, cancels a process through a function it is given and
refuses every later write. None of that knows Rust, TypeScript or the
capability report. The rule says: "Per-language copies of that mechanism are
forbidden." It also names the reason that was given for the copies: "the size
of a helper cannot justify language-specific ownership of an agnostic
mechanism". The four assessment questions of the rule have these answers:

1. **What is language-specific?** The bounds (16 MiB and 256 KiB for Rust,
   16 MiB and 16 KiB for TypeScript, `MaxMetadataBytes` for the loader), how
   the streams share them, and what a passed bound means to the caller's own
   callers. Each caller keeps these.
2. **What stays the same when those inputs are substituted?** The refusal of
   a write that would pass the bound, the cancellation, the refusal of every
   later write, and that no method takes bytes around the bound.
3. **Did that mechanism have an agnostic implementation?** No. A third
   caller copied it, and that copy was wrong.
4. **Can each caller use the same mechanism and keep its own meaning?** Yes,
   as the correction shows.

[AC-05](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-05)
names resource limits among the guarantees that composition must preserve.
It is the consequence here and not a second finding: a guarantee that every
caller implements again holds only where each copy is right, and one of the
three was not.

[AC-07](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-07)
explains why the first correction left the copies. Its evidence was real:
each bridge had a test with a real process, and 16 mutants per bridge failed
a test. That evidence shows that each copy held its bound. It says nothing
about who owns the mechanism.

## Intended boundary

One implementation of the bounded writer, in a package that names no
language and no tool. A caller supplies the number of bytes, the function
that cancels its process, and its own answer to a passed bound. The
[layout](../../LAYOUT.md#dependency-direction) has the place: the shared
private helpers, "A helper of no single language cannot conceal a language",
which are entries of the closed table `foundationRules` in the architecture
test.

## Actual implementation and consequence

Three packages each declared the writer.

- **The two bridges.** `boundedBuffer` was the same type twice, with the same
  comment. Both embedded `bytes.Buffer` until pull request
  [#298](https://github.com/Bitspark/bittype/pull/298) corrected them on
  13 November 2026, and so neither enforced its bound before that day
  ([#296](https://github.com/Bitspark/bittype/issues/296)).
- **The loader of the reporting tool.** `boundedOutput` embedded
  `bytes.Buffer` beside a budget that the output streams and the error
  streams of two Go commands share. `os/exec` copies a child's stream with
  `io.Copy`, which calls the writer's `ReadFrom` when it has one. So the
  bounded `Write` was not called, and the budget `MaxMetadataBytes` was not
  enforced. The [capability report design](../capability-reports.md) says
  "Loading bounds metadata, file count and captured source bytes". For the
  metadata that was not true.

The consequence is the one the rule predicts. The defect was found in one
copy and corrected in two. The third copy kept it, because nothing joined
the three. The decision on #296 accepted "two stay" for this reason: one
shared type would have needed an entry in the architecture test. An entry is
an edit to make and not a reason to keep a copy.

## What worked and what it did not establish

After #298 each bridge held its bound. `TestHelperOutputBound` in each
bridge package starts a real child process through the bridge's exported
`Check` and requires `semantic.ErrLimit` and no output one byte over each
bound. Those tests are unchanged by this correction and pass against the one
writer.

They did not establish that the mechanism had one owner, and they could not
see the third copy.

## Required correction

1. Implement the bounded writer once, in a package that names no language
   and no tool, and place it from the layout.
2. Move the two bridges and the loader to it and delete their copies.
3. Keep with each caller its bounds, how its streams share them, and what a
   passed bound means.
4. Enforce the loader's budget, and show it with a real process.
5. Give the architecture test the entry of the package and the import of
   each importer.

## Correction and qualification status

Corrected by the change for
[#299](https://github.com/Bitspark/bittype/issues/299), pull request
[#313](https://github.com/Bitspark/bittype/pull/313).

### The change

- **One package.**
  [`internal/processoutput`](../../internal/processoutput/processoutput.go)
  exports `Bound`, `Writer` and `ErrBound`. A `Bound` is a number of bytes
  and the function that cancels the process. `Bound.Writer` returns a writer
  that draws on it, and several writers of one bound share it. `Write` is
  the only method that takes bytes. A write that would pass the bound
  appends nothing, cancels once and returns `ErrBound`, and every later
  write of every writer of that bound is refused. `NewBound` refuses a
  missing cancel function, and a writer that no bound returned refuses every
  write. The package imports no package of the module.
- **The interface is derived from the three uses and not from one copy**
  ([AC-09](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-09)).
  The bridges bound each stream by itself. The loader bounds four streams of
  two processes together. One writer type with its own count serves the
  first and not the second, so the count is its own value, and a writer
  draws on it. The loader's copy kept the part of a refused write that still
  fitted. No caller reads a stream whose bound was passed, so the one writer
  keeps nothing of a refused write.
- **Each caller keeps its meaning.** The bridges read `Exceeded` and return
  `semantic.ErrLimit` and no output. The loader returns its error "Go
  command output exceeds N bytes". The package knows neither.
- **The architecture test** gains the entry of the package in
  `foundationRules`, with no import, and the import for each of the three
  importers: a role extension for `internal/rustsem` and for
  `internal/typescriptsem`, and the package in the entry of
  `tools/provider-matrix/report/go`. Each of the three is checked for use by
  its own package.

### Qualification

- **The loader's budget is enforced.** `TestLoadBoundsGoCommandOutput` puts
  a copy of the test binary in the search path under the name `go`. The four
  streams of the two commands fill the budget exactly, and the loader goes
  on to its next step. One byte more on any of the four, or a first command
  that passes the budget alone, fails the load with the budget's message,
  for a process that exits and for one that never exits by itself.
  `TestLoadFailsWhenTheGoCommandPassesTheBudget` does the same with the Go
  command itself and a budget of sixteen bytes. Both fail against the
  loader of `main`: it read every oversized stream whole and cancelled
  nothing.
- **The one writer.** The tests of `internal/processoutput` read its method
  set, run it against a real child process with one bound per stream and
  with one shared bound, and write to it from eight goroutines.
- **The bridges.** `TestHelperOutputBound` passes unchanged in both
  packages. `TestBoundedBufferOnlyWrites` named the fields of a bridge's
  own type and went with that type; the tests of the owner assert what it
  asserted.
- **A copy fails.** A bridge or the loader as it is on `main`, with its own
  writer and without the import, fails `TestProductionArchitecture`: its
  extension or its entry is unused.
- **Mutants.** 58 single changes each fail a test: 18 of the writer, 12 of
  each bridge's use of it, 8 of the loader's use of it and 8 of the entries
  in the architecture test.
- **Output.** 153 values that the integration tests generate, format or
  place have the same SHA-256 before and after.

The pull request holds the runs.

### Limits

- **A copy that keeps the import.** The architecture test sees imports. A
  package that uses the shared writer for one stream and a correct writer of
  its own for another passes every test. Only a reader sees it.
- **A fourth caller.** Nothing fails when a new package starts a process and
  writes its own bounded writer. The layout says that it must not, and the
  import check cannot say it.
- **The race detector** does not run on the workstation. The test with eight
  goroutines failed in 20 of 20 runs without the lock. The hosted step
  `syntax-race` runs `TestRustSemanticConcurrent` and
  `TestTypeScriptSemanticConcurrentQueries` under the detector, and both
  reach a bridge and so the one writer. The tests of
  `internal/processoutput` themselves do not run under it.
- **The Rust helper** does not run on the workstation. The hosted checks run
  its tests against the one writer.
- **Other limits of the three callers** are not examined: the deadlines, the
  wait after a cancellation, and a process that a helper starts itself.
- **The child-process fixture of the tests** exists once per package: in
  both bridges, in the loader and in the owner. Each starts its process the
  way its package does. Whether the fixture is one mechanism was not
  assessed.

## Assessment lesson

The first correction asked whether each copy held its bound, answered it
with good evidence, and then weighed the cost of one shared type against the
size of the helper. Placement is not a cost question. The question that
finds this fault is the third of AC-11: would another caller copy the
mechanism? One already had, with the defect that the correction had just
removed from the other two.

[Package boundaries](../package-boundaries.md) · [Layout](../../LAYOUT.md#dependency-direction) · [Case index](README.md)
