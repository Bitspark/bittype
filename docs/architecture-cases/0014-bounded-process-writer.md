# bittype/ACX-014: The bounded writer for a child process had one copy per caller

**Status:** Corrected, with limits; state of 13 November 2026. One package,
`internal/processoutput`, holds the writer, the three callers use it and no
copy is left. The fixture of the tests, which this correction had first
written a third and a fourth time, has one owner too. The [limits](#limits)
say what the correction does not establish. An independent review read the
first form at `c8c3509`, and its findings are corrected in the
[second form](#independent-review).
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
  draws on it.
- **Two differences in behaviour, and no caller sees either.** The loader's
  copy kept the part of a refused write that still fitted. The one writer
  keeps nothing of a refused write. And `Bytes` returns nothing once its
  bound is passed: from the first refused write on, no writer of the bound
  takes another byte, so what each holds may be a stream with its end cut
  off, and a cut stream must not read as a whole one. Each caller asks
  `Exceeded` before it reads `Bytes` and returns no output after a passed
  bound.
- **The cancel function is the caller's.** It is called once, after the
  bound reports itself passed and outside the bound's lock, so it may ask
  the bound.
- **Each caller keeps its meaning.** The bridges read `Exceeded` and return
  `semantic.ErrLimit` and no output. The loader returns its error "Go
  command output exceeds N bytes". The package knows neither.
- **The architecture test** gains the entry of the package in
  `foundationRules`, with no import, and the import for each of the three
  importers: a role extension for `internal/rustsem` and for
  `internal/typescriptsem`, and the package in the entry of
  `tools/provider-matrix/report/go`. Each of the three is checked for use by
  its own package.

### The fixture of the tests

The same rule applies to test code. Pull request #298 gave each bridge a
test with a real child process: the test binary, started as the child,
writes stated numbers of bytes to its streams, and in some cases connects to
the test first and never exits by itself, so that only a cancellation ends
it. The first form of this correction wrote that child and its harness a
third time for the loader and a fourth time for the writer, and recorded
"whether the fixture is one mechanism was not assessed". The review measured
it: 60 code lines were identical in all four test files, and 116 in the two
bridge files. The protocol was one. What differed is what each package
owns: how its test binary knows that it is the child, what the output stream
starts with, whether the child reads its input stream, and the call under
test.

[`internal/processoutput/processoutputtest`](../../internal/processoutput/processoutputtest/processoutputtest.go)
now holds the child, its description file and the harness, and the four test
files use it. It imports no package of the module and has its entry in
`foundationRules`. Only test files import it, and the architecture test
reads no test file, so no importer has an allowance. The two bridge test
files of #298 are edited for it: their cases, bounds and assertions are
unchanged, and their own copy of the child and of the harness is gone.

The copies also shared a defect, which is corrected once. The harness
waited for the first of two events, the child's connection and the return
of the call, and failed when the return came first. Nothing orders the two
in the test's process: a child that made its parent return had connected
before it wrote, and the test may still learn of the return first. The
review saw that false failure in 1 to 10 of 100 runs on a loaded host, and
the same lines were in the bridge tests on `main`. The harness now keeps a
return that comes first and awaits the connection after it, with a bound.

### Independent review

A reviewer that wrote none of the change read its first form at `c8c3509`
on 13 November 2026. The record is a
[comment on the pull request](https://github.com/Bitspark/bittype/pull/313#issuecomment-6858433447).
It found no defect in the writer or in the three callers, by reading and
under 16 mutants of its own, and confirmed the red run on `main`, the
place and the three import allowances. Its seven findings and what the
second form does:

1. **The false failure of the "process stays" tests.** Corrected once, in
   the shared harness.
2. **The fixture in four copies.** One package, as above.
3. **No race detector had run two writers on one bound.** The hosted job
   has a step that runs the tests of `internal/processoutput` and of the
   loader's package under the detector.
4. **`Bytes` returned a cut stream with no sign that it was cut.** It
   returns nothing once its bound is passed.
5. **`docs/package-boundaries.md` disagreed with itself.** Its counts were
   made again.
6. **"Cancel is called outside the lock" had no test.**
   `TestCancelMayAskTheBound` gives the bound a cancel function that asks
   the bound and writes to it. Called inside the lock it never returns, and
   the test fails when its wait passes.
7. **One `WaitDelay` failure in 300 runs under load**, in a case of the
   writer's own test where the process exits by itself. `os/exec` gives the
   copies of the streams the command's `WaitDelay` to end after the process
   has exited, and returns an error when it passes. That test had taken the
   bridges' one second for a call of its own. Its wait is now the harness's
   bound, because the call is the test's own and the wait is not what the
   test is about. The tests of the bridges and of the loader go through the
   caller, so they run the caller's real wait, one second and two seconds.
   A failure there would report the caller's own limit, which this case
   leaves out.

The reviewer did not read the second form.

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
  with one shared bound, let its cancel function ask the bound, and write to
  it from eight goroutines while a ninth reads the bound and the writers.
- **The bridges.** `TestHelperOutputBound` passes in both packages with its
  cases and assertions unchanged. `TestBoundedBufferOnlyWrites` named the
  fields of a bridge's own type and went with that type; the tests of the
  owner assert what it asserted.
- **A copy fails.** A bridge or the loader as it is on `main`, with its own
  writer and without the import, fails `TestProductionArchitecture`: its
  extension or its entry is unused.
- **Mutants.** 65 single changes: 21 of the writer, 2 of the fixture, 12 of
  each bridge's use of the writer, 8 of the loader's use of it and 10 of the
  entries in the architecture test. 64 fail a test. The one that does not
  removes the lock from `Exceeded` and `Bytes`: without the race detector
  no test can fail for it, because a test that orders its read after the
  writes reads the right value without the lock, and a test that does not
  order them accepts either value.
- **The false failure.** The process tests of the four packages ran 300
  times each with the default number of processors and 300 times each on
  one: 2400 runs and no failure. The test files of the first form, run after
  them under the same conditions, failed 5 of 300 runs of the
  writer's test and 1 of 100 of the loader's.
- **The race detector.** The hosted step runs the tests of the writer and of
  the loader's package under it, `TestConcurrentWritersKeepTheBound` among
  them.
- **Output.** 159 values that the integration tests generate, format or
  place have the same SHA-256 on `main` and with this change.

The pull request holds the runs.

### Limits

- **A copy that keeps the import.** The architecture test sees imports. A
  package that uses the shared writer for one stream and a correct writer of
  its own for another passes every test. Only a reader sees it. The review
  confirmed this with two mutants.
- **A fourth caller.** Nothing fails when a new package starts a process and
  writes its own bounded writer, or a test its own child process. The layout
  says that it must not, and the import check cannot say it.
- **The lock of `Exceeded` and `Bytes`** is shown only by the race
  detector, in the hosted step. That the detector reports the mutant was
  reasoned and not run: the workstation has no detector.
- **The Rust helper** does not run on the workstation. The hosted checks run
  its tests against the one writer.
- **Other limits of the three callers** are not examined: the deadlines, the
  wait after a cancellation, and a process that a helper starts itself.
- **The harness's check that the process ended** has no mutant: with these
  callers, a call that returns has ended its process, so no single change
  makes a process outlive its call.

## Assessment lesson

The first correction asked whether each copy held its bound, answered it
with good evidence, and then weighed the cost of one shared type against the
size of the helper. Placement is not a cost question. The question that
finds this fault is the third of AC-11: would another caller copy the
mechanism? One already had, with the defect that the correction had just
removed from the other two.

The correction then repeated the fault it corrected, one level down. It
copied the child process of the tests twice more and wrote that the question
was not assessed. A rule about mechanisms does not stop at the line between
code and its tests, and "not assessed" is not an answer to a question that
the change itself raises by adding the copy.

[Package boundaries](../package-boundaries.md) · [Layout](../../LAYOUT.md#dependency-direction) · [Case index](README.md)
