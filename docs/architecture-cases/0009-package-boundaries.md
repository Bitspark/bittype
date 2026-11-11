# bittype/ACX-009: Package boundaries do not follow contributions

**Status:** Partly corrected; state of 12 November 2026. The import rules by
role are corrected and the merged moves are independently reviewed. Correction
is pending for Go's contract identities and query service, with the eight hand
entries in the test that wait for them, and for the inline horos keys. Three
further candidates are not confirmed.
**Scope:** The native and common horos mapping packages, the horos-to-Go
verification laws, the Go and TypeScript analysis provider packages, the
publication of contract identities and horos semantic names, and the import
check in
[`integration/go/architecture_test.go`](https://github.com/Bitspark/bittype/blob/3bf45c3/integration/go/architecture_test.go),
at `3bf45c3`. Identified on 11 November 2026 after the owner asked for this
check. Plan in [#218](https://github.com/Bitspark/bittype/issues/218); design in
[package boundaries](../package-boundaries.md).

## Rules involved

[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
states that a misplaced responsibility is incorrect and moves first. The
instances below are misplacements at the level of the package: a part sits in a
package whose other parts it does not need, or a name is published by an
implementation.

[AC-07](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-07)
limits what evidence supports. `LAYOUT.md` says the architecture checks
"classify every production package and enforce its directed dependencies". That
is true per package. It does not establish that two contributions inside one
package are independent, that an allowed import uses only what it should, or
that an allowance is needed.

No instance is recorded here as a violation of
[AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02)
or [AC-03](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-03).
An independent mapping library can still enter the kernel without editing any
of these packages. Two neighbouring faults are recorded elsewhere: the closed
profile gate is [bittype/ACX-008](0008-structural-profile-gate.md), and the
optional-field decision inside Record is the contract change of
[#196](https://github.com/Bitspark/bittype/issues/196), item 6.

## Intended boundary

One contribution, meaning one thing a consumer registers or selects on its own,
is one package. A name that two parties must agree on lives in a vocabulary
package that imports no implementation. Only a composition imports a
contribution. A contribution that needs another's result states a demand, and
the kernel resolves it with whatever the consumer registered.

## Actual implementation and consequence

This section records what was found at `3bf45c3` and is kept as found. The
state after the corrections of 11 November 2026 is under
[Correction and qualification status](#correction-and-qualification-status).

**Contributions grouped in one package.**

- `transform/syntax/horos-to-{typescript,golang,rust,cpp}/mapping/go` each hold
  `Text`, `Bool`, `List`, `Record` and `Declarations`. The Go package also
  holds `Optional`.
- Each of the four exports `Libraries(profile)`, a list of its own mappings.
  [`docs/generation.md`](../generation.md) states "There are no per-pair
  assembly factories or presets."
- `verification/files/horos-to-go/mappings/go` holds the `Text`, `Boolean`,
  `List` and `Record` laws.
- `transform/syntax/horos/mapping/go` holds the profile gate with its `Library`
  adapter, the pattern constructors, the declaration recipe and the naming
  function.
- `analyze/files/golang/provider/go` and `analyze/files/typescript/provider/go`
  each hold the analysis provider and a native query service. The two share no
  declaration.

**Identities published by an implementation.** 22 imports go from a production
package into a stage provider package. Eight use only an identity:

| Importer | Provider package | Symbols |
| --- | --- | --- |
| `analyze/files/golang/provider/go` | `resolve/files/golang/provider/go` | `Contract` |
| `analyze/files/rust/provider/go` | `resolve/files/rust/provider/go` | `Contract` |
| `analyze/files/typescript/provider/go` | `resolve/files/typescript/provider/go` | `Contract`, `Compiler` |
| `analyze/files/horos/provider/go` | `resolve/files/horos/provider/go` | `Contract` |
| `internal/horosfiles` | `resolve/files/horos/provider/go` | `PreludePin` |
| `transform/syntax/horos/source/go` | `analyze/files/horos/provider/go` | `Contract`, `Profile` |
| `transform/syntax/horos/source/go` | `resolve/files/horos/provider/go` | `ExpressionRole`, `IdentityRelation`, `ResultSortRelation`, `SortRelation` |
| `verification/files/source/horos/go` | `resolve/files/horos/provider/go` | `Contract`, `Profile` |

A ninth, `verification/files/horos-to-go/go`, imports the Go analysis provider
only for the query service's `Admission`.

**The check.** The four native targets may import their checker, and none does.
The test passes. An allowed import may use any exported symbol.

Consequences:

- Coupling between two mappings of one pair is invisible to the import check.
  One exists: Go's Record reaches the Optional helper through the package
  variable `optionalKey`
  ([source](https://github.com/Bitspark/bittype/blob/3bf45c3/transform/syntax/horos-to-golang/mapping/go/mapping.go#L72-L88)).
- A representation family written into one of these packages has to move when
  the package is divided. [#207](https://github.com/Bitspark/bittype/issues/207)
  writes the first ones.
- A reader of horos semantic facts imports the resolver implementation to name
  a relation. A second resolver would have to be imported next to the first, or
  repeat the names.
- Every new mapping edits `Libraries`, and every new package edits the test's
  list.

## What worked and what it did not establish

The architecture test passes and its rules hold: the kernel, the source and the
targets import no concrete mapping package. Independent contributions enter the
kernel through the public contract, as the composition tests show. Those
results establish that the kernel is open and that the listed edges are absent.
They do not establish that packages follow contributions, because the test
reads imports per package and accepts any use of an allowed import.

## Required correction

1. Each builtin mapping and each verification law is its own package, and no
   package lists them ([#222](https://github.com/Bitspark/bittype/issues/222)).
2. The query services leave the analysis provider packages
   ([#221](https://github.com/Bitspark/bittype/issues/221)).
3. Contract identities and horos semantic names are published by vocabulary
   packages, and no package imports a provider for a name
   ([#220](https://github.com/Bitspark/bittype/issues/220)).
4. Every remaining import of a provider names its symbols and its reason, and
   an unused allowance fails
   ([#219](https://github.com/Bitspark/bittype/issues/219)).
5. The import rules are stated per role, so a new package of an existing role
   needs no edit to the test
   ([#218](https://github.com/Bitspark/bittype/issues/218)).
6. Every identity value and every generated declaration stays byte-identical.
   These are moves.

## Correction and qualification status

**Partly corrected. The case stays open.** The state below was read in source
on `main` at `1a5aa2d` on 11 November 2026, and the architecture test again at
`71bfec9` on 12 November 2026. Only test files and documents changed between
those two commits. The order is the one in
[package boundaries](../package-boundaries.md#order-of-work). The pair packages
were split while [#207](https://github.com/Bitspark/bittype/issues/207), which
writes the first families into mapping packages, is still open. That is the
order
[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
requires.

### Corrections on `main`

| Instance | Correction | Pull requests |
| --- | --- | --- |
| The four native mapping packages each held every mapping of their pair | Each mapping is its own package, `transform/syntax/horos-to-<native>/mapping/<name>/go`, with `text`, `boolean`, `list`, `record` and `declarations`, and `optional` for Go. Each exports `Library`. The old packages are deleted and no alias remains | [#226](https://github.com/Bitspark/bittype/pull/226) rust, [#227](https://github.com/Bitspark/bittype/pull/227) golang, [#230](https://github.com/Bitspark/bittype/pull/230) cpp, [#232](https://github.com/Bitspark/bittype/pull/232) typescript |
| Each of the four exported `Libraries(profile)` | `Libraries` is deleted. A consumer registers each library | The same four |
| Go's Record reached the Optional helper through the package variable `optionalKey` | The consumer creates a helper key and passes it to `record.Library(profile, key)` and to `optional.Library(profile, key)`. Neither package creates the key and neither imports the other | [#227](https://github.com/Bitspark/bittype/pull/227) |
| The common horos mapping package held the adapter with the profile gate, the pattern constructors, the declaration recipe and the naming function | The adapter `Library` with the profile gate stays at `transform/syntax/horos/mapping/go`. `pattern`, `declarations` and `naming` are packages under `transform/syntax/horos/mapping/<name>/go`. `Pattern` is now `pattern.Head` and `Declarations` is now `declarations.Library`. The profile gate is unchanged | [#238](https://github.com/Bitspark/bittype/pull/238) |
| The laws package held the `Text`, `Boolean`, `List` and `Record` laws | Each law is its own package, `verification/files/horos-to-go/mappings/{scalar,list,optional,record}/go`. `scalar` exports the rules `Text()` and `Boolean()` of one law. The other three export `Rule` | [#228](https://github.com/Bitspark/bittype/pull/228) |
| The TypeScript analysis provider package held a query service | `Queries`, `NewQueries` and `Assignability` are in `analyze/files/typescript/queries/go` | [#233](https://github.com/Bitspark/bittype/pull/233) |
| Rust, TypeScript and horos identities were published by provider packages. horos had no vocabulary package | `semantics/rust/go`, `semantics/typescript/go` and the new `semantics/horos/go` export `ResolutionContract()` and `AnalysisContract()`. `semantics/horos/go` owns the horos roles, relations and profiles, the limits setting and the prelude pin. The TypeScript `Compiler` type is in `semantics/typescript/go`. Seven of the eight imports in the table above are gone | [#229](https://github.com/Bitspark/bittype/pull/229) rust, [#233](https://github.com/Bitspark/bittype/pull/233) typescript, [#234](https://github.com/Bitspark/bittype/pull/234) horos |
| The check accepted any use of an allowed import and passed unused allowances | Every import of a stage provider package outside a composition needs an entry that names its symbols and its reason. A use outside the named symbols fails. An unused entry, an unused symbol and an unused allowance fail. A dot import of a module package fails. 94 unused allowances were removed | [#225](https://github.com/Bitspark/bittype/pull/225), closing [#219](https://github.com/Bitspark/bittype/issues/219) |
| Every new package edited the test's list | A package gets its import rule from the role that its path gives. At `71bfec9`, 78 of the 126 production packages are classified this way, and a new package of an existing role needs no edit to the test. A package that names a language imports no package of another language. 48 packages keep an explicit entry: the 31 of the closed foundation and kernel set, 9 packages of which there is only one, and 8 entries held for Go. The explicit entries for the foundation and kernel set are a narrowing of the plan, decided on #218 on 11 November 2026: their layering among each other is a fact per package, which one rule for the role would lose | [#225](https://github.com/Bitspark/bittype/pull/225) for four package shapes, [#244](https://github.com/Bitspark/bittype/pull/244) for the others |
| The allow-list of a role was not checked for use | Every allowance is checked for use. An entry and an extension are checked per package. A role is checked per shape: each import it permits must be imported by at least one package of that role. External model allowances are checked too | [#244](https://github.com/Bitspark/bittype/pull/244) |
| The check read only the files of the build that ran it, and the vocabulary rule accepted any name as a language. The independent review found both | The check reads every non-test file of a package, also one that a build tag, another operating system or disabled cgo excludes. A role grants `internal/<language>` and `internal/<language>sem` only where the imported package is that language's helper by its own role | [#244](https://github.com/Bitspark/bittype/pull/244) |

[#224](https://github.com/Bitspark/bittype/pull/224) prepared the mapping and
law splits. It put each pair's list of libraries and the laws behind one test
seam file under `integration/go`. It changed no production code.

### Qualification

The statements below come from the pull requests named, unless a bullet names
another source. Their authors made the before and after comparisons. #227,
#228, #229, #233 and #234 used programs that are not committed.

- **Hosted checks.** The hosted `docs`, `contracts` and `native` checks passed
  on each of the eleven pull requests, and on #244.
- **Identity values.** #229 rendered the two Rust contract values before and
  after. #233 rendered the TypeScript contracts, the configured analysis
  contract and the `Compiler` type. #234 rendered the horos keys, profiles,
  prelude pin and contracts. Each reports byte-identical output.
- **Generated output.** #226, #230 and #232 report that the Rust, C++ and
  TypeScript source logged by `TestSemanticPaginationMappings` is
  byte-identical before and after. #227 reports the same for the formatted Go
  of the pagination module, every plan step and the description of the six
  libraries. #238 reports that the source the same test logs for all four
  targets is byte-identical before and after.
- **Law names and keys.** #228 printed each rule's key, dependency name and
  law type before and after and reports the same output.
- **Seeded forbidden imports.** Each of these was added, failed
  `TestProductionArchitecture` and was removed again: a mapping importing
  another mapping of its pair (#226, #227, #230, #232), a law importing another
  law (#228), the Rust analyzer importing its resolver again (#229), the
  TypeScript query service importing its resolver (#233), `semantics/horos/go`
  importing a stage provider (#234), the `pattern` package importing `naming`
  (#238), and a formatter importing its parser and using an unnamed symbol of
  its checker (#225).
- **The check itself.** `integration/go/architecture_synthetic_test.go` holds 36
  tests that feed the check hand-built package sets. Each failing set breaks
  one rule and must produce that rule's message and no other. One test adds a
  third language, a pair into it, one more mapping and one more common mapping
  part without a new line in the foundation table, the hand entries or the
  owners. One test writes a module at test time and runs the real loader on
  files that a build tag, another operating system and cgo exclude.
- **Go's Record and Optional.** `TestGolangRecordRequiresItsOptionalHelper`
  covers an invalid key, a key that no library supplies and the working case.
- **The rules by role, by their author.** #244 reports the following on the
  real tree. Five seeded forbidden imports failed `TestProductionArchitecture`:
  a Rust parser importing TypeScript's syntax, a vocabulary importing a stage
  provider, a mapping importing a native target, the engine importing a
  language's syntax and a lexical helper importing a stage contract. Three
  seeds in files that the build excludes failed: behind a build tag, in a
  `_linux.go` file on Windows and in a cgo file with cgo off. Fourteen scratch
  packages passed `TestProductionArchitecture` without an edit to the test
  files: a new language, two mappings, a law and an extractor.
- **The rules by role, by the integrator.** The integrator reports in a
  [status comment on #218](https://github.com/Bitspark/bittype/issues/218#issuecomment-6809884454)
  that ten further forbidden imports were seeded on the head of #244 before it
  merged, and that each failed. The comment gives the count. The list is the
  integrator's, as given to the change that wrote this record: a mapping
  importing another pair's mapping, a formatter importing its own parser, a syntax vocabulary importing
  the text contract, the Rust analyzer importing its resolver, a native target
  importing a source reader, the horos resolver importing TypeScript's syntax,
  a law importing the Go analysis provider, the generation contract importing
  the horos vocabulary, `naming` importing `pattern`, and an unknown package
  `internal/foo`. The same ten were repeated for this record at `71bfec9` in a
  scratch copy. Each failed `TestProductionArchitecture`. Four failed with the
  reason of a prohibition. Five failed as an import that the package's rule
  does not name, and three of those also as a provider import without an
  exception. `internal/foo` failed as unclassified.
- **Independent review of the moves.** A reviewer that wrote none of the
  changes read the range `d23542a..e83cf52` from the code. Its verdict is that
  they are moves. It compared 237 production declarations with the old ones.
  It printed every identity at both revisions, 3352 lines with the same
  SHA-256. It compared the generated output of the four targets, and it built
  the installed-consumer fixture, which passed seven of its eight tests. It
  found four defects, all in the check or in comments and none in moved code,
  and #244 corrected them. The review is recorded in a
  [comment on #218](https://github.com/Bitspark/bittype/issues/218#issuecomment-6809378108),
  and [package boundaries](../package-boundaries.md#independent-review) gives
  its method.

### Limits

- The comparisons of #227, #228, #229, #233 and #234 cannot be repeated from
  the repository, because the programs that made them were not committed. The
  tests that compare fragments of the generated declarations remain.
- The independent review covered the moves in `d23542a..e83cf52`. It did not
  verify anything that needs the Rust helper, the native compiler job, or the
  race and fuzz steps. The hosted checks ran those on each pull request. Its
  range ends before #244, which changes test code only.
- The review's programs and outputs are in a scratch directory on the
  workstation and not in the repository. The same holds for the seeded
  imports of #244 and of the integrator, and for their repeat at `71bfec9`.
  None of them can be repeated from the repository alone.
- Use of a role's allowances is checked per shape, not per package. One
  package may import what a sibling of the same shape imports.
- An import check does not see a re-export or an alias. The review confirmed
  by seed that a helper holding an exception can re-export the excepted
  symbol, and that a package can alias an excepted type for another package.
  Both pass. [Package boundaries](../package-boundaries.md#what-the-check-cannot-establish)
  lists what else the check cannot establish.
- The helper key that Go's Record receives is an interim form. Record still
  decides how an optional field is represented in Go, Rust and C++. That is
  [#196](https://github.com/Bitspark/bittype/issues/196), item 6, and is not
  part of this case.
- A green run still shows only that the listed edges are absent
  ([AC-07](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-07)).

### Instances that remain

1. **Go's query service.** `analyze/files/golang/provider/go` still holds
   `Queries`, `NewQueries` and `Admission` next to the analysis provider.
   `verification/files/horos-to-go/go` imports that package for `Admission`.
2. **Go's contract identities.** `resolve/files/golang/provider/go` and
   `analyze/files/golang/provider/go` still publish `Contract()`, and the
   analyzer imports the resolver package for it.
3. **Inline horos keys.** `horos/module` and `horos/import` are spelled inline
   in `internal/horosfiles`. `horos/type-definition`, `horos/sort/proper` and
   `horos/identity/module` are spelled inline in
   `transform/syntax/horos/source/go`. The change that added
   `semantics/horos/go` found them. They are tracked in
   [#220](https://github.com/Bitspark/bittype/issues/220).

Instances 1 and 2 are the two open entries of the exception table. Their
correction is prepared in draft pull request #231. That pull request is
stacked on open pull request #204 and is not on `main`.

Eight hand entries in `TestProductionArchitecture` wait for the same two pull
requests. Five are entries that the pull requests edit: the four Go packages
`semantics/golang/go`, `internal/golangsem`,
`resolve/files/golang/provider/go` and `analyze/files/golang/provider/go`, and
`verification/files/horos-to-go/go`. Three are neighbouring lines kept in
place so that the pull requests merge cleanly. A comment in the test says how
the eight convert afterwards. The case stays open while these and the three
instances remain.

The profile gate in the adapter of the common horos mapping package is not an
instance of this case. It stays with
[#216](https://github.com/Bitspark/bittype/issues/216) and
[bittype/ACX-008](0008-structural-profile-gate.md), and #238 did not change it.

### Candidates

Three candidates are not confirmed and are not part of this case until they
are.

- **The thirteen direct calls into a same-language checker or parser.** They
  are accepted entries of the exception table, each with its symbol and its
  reason. The closing comment of #219 records what a move would cost. horos
  could import `internal/declsyntax` directly, which would remove three of the
  eleven `Check` entries. The four native languages keep `Check`, its table and
  its helpers in the provider package. `Parse` is the whole parser. The
  question stays a candidate on #218.
- **The `Self` rule in the shared naming function.** Decided on #222 on
  11 November 2026: the rule stays where it is. The naming function gives one
  declaration the same name in every target. That is a property of the
  profile, so the question belongs to #216.
- **`CheckSingleScope` in the shared target package.** #222 did not decide it.
  It stays with the second design point of
  [#161](https://github.com/Bitspark/bittype/issues/161).

## Assessment lesson

The import check was read as proof that contributions are independent. It
proves that listed edges between packages are absent. Two questions were never
asked of it: what can reach what inside one package, and how much of an allowed
import is used. When a dependency check is green, ask what its unit is, and
whether the thing that must be independent is smaller than that unit.

The independent review found that the unit was also smaller than "the
package". The check read the files that the build running the test selects,
so a file behind a build tag or for another operating system was never read.
Ask as well which files, builds and settings the check's unit leaves out.
