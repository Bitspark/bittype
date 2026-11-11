# Package boundaries and forbidden dependencies

**Implemented except for Go. State of 12 November 2026.** Steps 1, 4 and 5 of
the [order of work](#order-of-work) are on `main` at `71bfec9`. Steps 2 and 3
are on `main` except for Go, whose part waits for draft pull request #231 on
open pull request #204. The inline horos keys of finding 12 are open as well. The
[findings table](#findings-on-main-at-3bf45c3) gives the state of each
finding, and [what the test enforces now](#what-the-test-enforces-now) says
which rules are checked and how. An
[independent review](#independent-review) read the merged moves. The reasoning
sections are unchanged since the proposal.

The owner asked on 11 November 2026 to check for packages "that contain multiple
implementations / functionalities that either already are separate and neither
conceptually nor practically depend on each other, or which still depend on
each other but shouldn't, because it violates our kernel principles (e.g., a
function calling another function directly instead of invoking the kernel for
invoking the function's semantics, the kernel then deciding which registered
implementation to call)", and for a way to enforce the result "through
forbidden dependencies that we could check for and enforce".

This document records the difference between the two faults, what to look for,
the findings of a scan of `main` at `3bf45c3`, the rules to enforce and the
order of work.

[Plan, #218](https://github.com/Bitspark/bittype/issues/218) ·
[Open findings, bittype/ACX-009](architecture-cases/0009-package-boundaries.md) ·
[Layout](../LAYOUT.md) ·
[Architecture test](../integration/go/architecture_test.go) ·
[The check](../integration/go/architecture_check_test.go) ·
[The roles](../integration/go/architecture_roles_test.go) ·
[Exception table](../integration/go/architecture_exceptions_test.go) ·
[Kernel-owned composition](kernel-composition.md) ·
[Shared guidance](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#package-boundaries-and-forbidden-dependencies)

## Terms

- **Kernel.** The shared mechanisms that select and compose: the engine, the
  generation registry and planner, and the contracts they are written against.
- **Contribution.** One thing a consumer registers or selects on its own: a
  stage provider, a source, a target, one mapping library, one verification
  law.
- **Vocabulary.** The names two parties must agree on and nothing else: syntax
  kinds, semantic keys, profile and contract identities, intended native forms.
  A vocabulary package imports no contribution.
- **Composition.** Code that chooses contributions and registers them: the
  integration tests, the tools and a future command.
- **Private part.** Code that serves one owner and is never selected by a
  consumer, such as a lexical helper of one language.

## Two different faults

| | Wrong grouping | Wrong coupling |
| --- | --- | --- |
| What it is | Parts that do not need each other share one package | A contribution names another contribution that the consumer is meant to choose |
| What is wrong | The package boundary. No dependency between the parts is wrong | The dependency. The caller makes the kernel's choice |
| Consequence | An importer of one part receives all of them, and the import check cannot tell the parts apart | The kernel's selection, scope and failure guarantees do not apply to that call. Replacing the callee means editing the caller |
| Correction | A move. Nothing is rewired | A demand that the kernel resolves. This can need a kernel capability |

The first fault hides the second. Inside one package Go lets every declaration
reach every other, and the import check works per package. Coupling between
two contributions in one package is therefore invisible to it. Splitting comes
first, which is the order
[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
requires: move what is misplaced, then add what is missing.

## The question that separates them

Not every direct call is wrong. Ask of each call from one part to another:

> Whose answer does the caller need: its own, or the composition's?

- **Its own.** The caller relies on a private invariant, a shared mechanism or
  a vocabulary. A direct call is right. If the callee sits in a contribution
  package, the question is only whether it is in the right package.
- **The composition's.** The caller needs whatever the consumer registered for
  that meaning. It must state a demand and let the kernel resolve it. A direct
  call fixes the consumer's choice inside the caller.

Two examples from `main` at `3bf45c3`:

- In a Record mapping, `r.Mapped(f.Use)` asks the kernel for the native type
  of a field. Whichever mapping the consumer selected answers. This is the
  correct form.
- Two lines below, the same Record mapping decides how an optional field is
  represented. Rust builds `std::option::Option`, C++ builds `std::optional`,
  and Go wraps the type in a helper through a package variable. A consumer who
  wants another representation of optional fields must replace Record.

Parts that share a private invariant belong together. The planner and the plan
types in `contracts/generation/go` are one unit, because only the planner may
issue a plan.

## What to look for

Signs of wrong grouping:

- **Several contributions in one package.** A consumer can register each one
  separately, and they share a package only because they were written together.
- **A function that lists the package's contributions.** Adding a contribution
  edits the list.
- **An importer that uses a slice of the package.** It imports a provider
  package and uses one identity, or one service that shares no declaration
  with the provider.
- **A name kept with an implementation.** A profile, a contract, a relation
  key or a pin is published by a provider package and not by a vocabulary.
- **A policy next to a shared mechanism.** A shared package also holds one
  case's rule.

Signs of wrong coupling:

- **A contribution imports a contribution.** Each such import needs a reason
  that survives the question above.
- **A callback that builds where it could ask.** It constructs the native form
  of something another contribution owns.
- **A package variable shared by two contributions.** One reaches the other
  without declaring it.
- **A shared function that carries one case's meaning.** Every caller inherits
  a rule that belongs to one of them.

## Findings on `main` at `3bf45c3`

The scan covered bittype only. Each finding was read in source after the scan
reported it. The columns "Where" and "Assessment" describe `3bf45c3`. The
column "State" describes `main` at `71bfec9` on 12 November 2026. It was read
in source at `1a5aa2d` on 11 November 2026, and finding 8 again at `71bfec9`.
Only test files and documents changed between those two commits. Finding 12
was found on 11 November 2026, after the scan.

| # | Where | Kind | Assessment | State | Tracked |
| --- | --- | --- | --- | --- | --- |
| 1 | `transform/syntax/horos-to-{typescript,golang,rust,cpp}/mapping/go` each hold Text, Bool, List, Record and Declarations. Go also holds Optional | Grouping | Confirmed | Corrected with [#226](https://github.com/Bitspark/bittype/pull/226), [#227](https://github.com/Bitspark/bittype/pull/227), [#230](https://github.com/Bitspark/bittype/pull/230) and [#232](https://github.com/Bitspark/bittype/pull/232). Each mapping is its own package, `transform/syntax/horos-to-<native>/mapping/<name>/go`, and exports `Library` | [#222](https://github.com/Bitspark/bittype/issues/222) |
| 2 | `Libraries(profile)` in each of those packages lists the package's own mappings | Grouping | Confirmed | Corrected with the same four pull requests. `Libraries` is deleted and a consumer registers each library | [#222](https://github.com/Bitspark/bittype/issues/222) |
| 3 | `verification/files/horos-to-go/mappings/go` holds the Text, Boolean, List and Record laws | Grouping | Confirmed | Corrected with [#228](https://github.com/Bitspark/bittype/pull/228). Each law is its own package, `verification/files/horos-to-go/mappings/{scalar,list,optional,record}/go` | [#222](https://github.com/Bitspark/bittype/issues/222) |
| 4 | `transform/syntax/horos/mapping/go` holds the profile gate, the pattern constructors, the declaration recipe and the naming function | Grouping | Confirmed | Corrected with [#238](https://github.com/Bitspark/bittype/pull/238). The adapter `Library` with the profile gate stays at that path. The pattern constructors, the declaration recipe and the naming function are the packages `pattern`, `declarations` and `naming` under `transform/syntax/horos/mapping/<name>/go`. The division is a move. The profile gate is unchanged | [#222](https://github.com/Bitspark/bittype/issues/222), done-when 4; the gate in [#216](https://github.com/Bitspark/bittype/issues/216), which is untouched |
| 5 | `analyze/files/{golang,typescript}/provider/go` each hold the analysis provider and a native query service that share no declaration | Grouping | Confirmed | Partly corrected. TypeScript is corrected with [#233](https://github.com/Bitspark/bittype/pull/233): the query service is `analyze/files/typescript/queries/go`. Go is open: its query service is still in `analyze/files/golang/provider/go` | [#221](https://github.com/Bitspark/bittype/issues/221); Go in draft pull request #231 |
| 6 | Eight imports of a provider package use nothing but an identity. horos has no vocabulary package for its semantic names | Identity kept with an implementation | Confirmed | Partly corrected. Rust with [#229](https://github.com/Bitspark/bittype/pull/229), TypeScript with [#233](https://github.com/Bitspark/bittype/pull/233), horos with [#234](https://github.com/Bitspark/bittype/pull/234), which adds `semantics/horos/go`. Seven of the eight imports are gone. Go is open: `analyze/files/golang/provider/go` still reads `Contract` from `resolve/files/golang/provider/go` | [#220](https://github.com/Bitspark/bittype/issues/220); Go in draft pull request #231 |
| 7 | Record decides optional fields itself in Go, Rust and C++ | Coupling | Confirmed; the contract change is designed | Open. [#227](https://github.com/Bitspark/bittype/pull/227) removed the package variable that Go's Record shared with Optional: the consumer now creates a helper key and passes it to both libraries. Record still decides the representation of an optional field in all three targets | [#196](https://github.com/Bitspark/bittype/issues/196), item 6 |
| 8 | The existing check passes unused allowances and cannot see which symbols an import uses | Enforcement gap | Confirmed | Corrected with [#225](https://github.com/Bitspark/bittype/pull/225), which names the symbols of every provider import, fails unused allowances and removed 94 of them. [#244](https://github.com/Bitspark/bittype/pull/244) checks the allowances of a role for use as well, and reads the files that build constraints exclude | [#219](https://github.com/Bitspark/bittype/issues/219), closed; [#218](https://github.com/Bitspark/bittype/issues/218) |
| 9 | Thirteen direct calls into a same-language provider: `Check` from parse, format and the horos resolver; `Parse` from the Rust and TypeScript semantic helpers | Coupling candidate | Read as legitimate: the caller needs its own grammar's answer. `LAYOUT.md` permits them | Not confirmed. The thirteen imports are accepted entries of the exception table, each with its symbol and reason. Whether the admission function should move to a language-private package is not decided | Recorded in [#219](https://github.com/Bitspark/bittype/issues/219); the candidate stays on [#218](https://github.com/Bitspark/bittype/issues/218) |
| 10 | `Exported` appends `_` to `Self` for all four targets. `Self` is a keyword only in Rust | Coupling candidate | Not confirmed | Not confirmed and unchanged. Decided on #222 on 11 November 2026: the rule stays where it is and is a question for #216 | [#216](https://github.com/Bitspark/bittype/issues/216) |
| 11 | `CheckSingleScope` states one placement's rule in the shared target construction package | Grouping candidate | Not confirmed | Not confirmed and unchanged. #222 did not decide it | With the second design point on [#161](https://github.com/Bitspark/bittype/issues/161) |
| 12 | horos-owned keys are spelled inline where they are used: `horos/module` and `horos/import` in `internal/horosfiles`, and `horos/type-definition`, `horos/sort/proper` and `horos/identity/module` in `transform/syntax/horos/source/go` | Identity without a single owner | Identified on 11 November 2026 by the change that added `semantics/horos/go` | Open | [#220](https://github.com/Bitspark/bittype/issues/220), added to its scope |

Finding 6 in detail, at `3bf45c3`. 22 imports go from a production package into
a stage provider package. Eight use only an identity (`Contract`, `Profile`, a
relation key, `PreludePin` or the `Compiler` configuration type), one uses the
co-located query service of finding 5, and thirteen are the calls of finding 9.
Go, TypeScript and Rust have `semantics/<language>/go` for their vocabulary and
still publish `Contract()` from the resolver package. horos has no such package:
most of its names are defined in `internal/horosmodel` and re-exported from the
resolver.

Finding 6 on 11 November 2026. `semantics/rust/go`, `semantics/typescript/go`
and the new `semantics/horos/go` export `ResolutionContract()` and
`AnalysisContract()`. TypeScript also exports `AnalysisContractFor(settings)`,
and horos exports `AnalysisContractWithResolutionConfiguration(settings)`.
`semantics/horos/go` owns the horos roles, relations and profiles, the limits
setting and the prelude pin. The TypeScript `Compiler` type and its `Validate`
method are in `semantics/typescript/go`. `Validate` only requires two absolute
paths. It reads no file and starts no process, so a vocabulary can hold it. Rust's
`Compiler` stays in `resolve/files/rust/provider/go`: it is the resolver's own
configuration and no other production package uses it. Go still publishes
`Contract()` from `resolve/files/golang/provider/go` and from
`analyze/files/golang/provider/go`. Their move is prepared in draft pull request
#231, which is stacked on open pull request #204 and is not on `main`.

Finding 7 in detail. The helper key is an interim form. In Go a consumer writes
`key := g.NewHelperKey("horos/optional")` and passes `key` to
`record.Library(profile, key)` and to `optional.Library(profile, key)`. Neither
package creates the key and neither imports the other. What a mapping's `Expand`
returns is unchanged, so Record still wraps the optional field itself.

Finding 12 in detail. The keys are written as
`semantics.Key{Owner: "horos", Name: "module"}` and the like at the place of
use. `internal/horosmodel` writes `type-definition`, `sort/proper` and
`identity/module` when it exports facts, and the source reader repeats them to
read those facts. `internal/horosfiles` compares an entity kind with
`horos/module`, which `internal/horosmodel` also writes. It is the same kind of
fault as finding 6: a name two parties must agree on has no single owner. #220
is done for these keys when `semantics/horos/go` publishes them and both sides
use it. Read in source at `1a5aa2d`.

## Rules to enforce

A package's role follows from its path, which the [layout](../LAYOUT.md) fixes.

| Role | Paths |
| --- | --- |
| Foundation | `trees/*/go`, `files/go`, `source/go`, `semantics/go`, `contracts/*/go` |
| Kernel | `engine/go`, `transform/syntax/shared/*/go`, `check/syntax/shared/*/go`, `print/*/shared/provider/go` |
| Vocabulary | `syntax/<language>/go`, `semantics/<language>/go`, `transform/syntax/<native>/model/go`, `transform/syntax/<native>/symbol/go` |
| Contribution | `<stage>/<dimension>/<language>/provider/go`, `analyze/files/<language>/queries/go`, `transform/syntax/<language>/{source,target}/go`, each mapping, each verification law and each capability extractor |
| Common mapping part | `transform/syntax/<source>/mapping/go` and `transform/syntax/<source>/mapping/<part>/go` |
| Private part | `internal/<language>` and `internal/<language>sem`, bound to their language by name. The other packages under `internal/` have explicit entries: five helpers of no single language in the foundation table, and four that the table `languageOwners` binds to one language |
| Composition | `integration/go`, `verification/integration/go`, `tools/**`. A future `cmd/**` is intended and is not yet known to the test |

1. A foundation or kernel package imports no vocabulary and no contribution.
2. A vocabulary package imports foundations and its own language's vocabulary,
   never a contribution.
3. A contribution imports foundations, kernel contracts, the vocabularies of
   its own scope and its own private parts.
4. Only a composition imports a contribution.
5. Every exception to rule 3 or 4 names the importer, the imported package, the
   exact exported symbols and the reason. A use outside the named symbols
   fails. An unused exception fails, and so does an unused symbol in one.
6. A production package whose path gives no role fails.

Rule 4 is the one that stops direct calls. If nothing but a composition can
import a contribution, no contribution can call another one.

The plan in #218 asked for a test that "holds no allow-list per package". That
was narrowed on #218 on 11 November 2026. The closed set of foundation and
kernel packages keeps one explicit entry per package. The reason is that their
layering among each other is a fact per package: `trees/ast/go` imports
nothing, `trees/cst/go` imports only the AST, and `contracts/text/go` holds
diagnostics only. One rule for the role would lose those facts. A new
foundation is rare and is meant to be an edit. Every role that can grow gets
its rule from its path: vocabularies, stage providers, sources, targets,
mappings, laws and the private helpers of a language. The test also keeps
explicit entries for nine packages of which there is only one. Its reason is
that a pattern for a single package would not be honest.

### What the test enforces now

Since [#244](https://github.com/Bitspark/bittype/pull/244) a package gets its
import rule from the role that its path gives. This section was read in
`integration/go/architecture_test.go`, `architecture_roles_test.go`,
`architecture_check_test.go` and `architecture_exceptions_test.go` at
`71bfec9`. The counts were made there.

Every production package of the module is classified in one of three ways. A
package with none of them fails as unclassified.

- **By its role.** `packageRole` derives the rule from the path. 78 of the 126
  production packages are classified this way. The
  [layout](../LAYOUT.md#dependency-direction) lists the roles and what each
  may import.
- **By its entry in `foundationRules`.** The entry is the exact list of the
  package's module imports. 31 packages have one.
- **By a hand entry.** The entry is again the exact list. 9 are in
  `oneOfAKindRules`, and 8 are written in `TestProductionArchitecture` and wait
  for Go. A hand entry takes precedence over a role.

The tables:

| Table | Entries | Holds |
| --- | --- | --- |
| `foundationRules` | 31 | The carriers, the contracts, the engine, the shared transform, check and print packages, and the five private helpers of no single language: `internal/tree`, `internal/grammar`, `internal/lex`, `internal/parseutil` and `internal/concrete` |
| `oneOfAKindRules` | 9 | Packages of which there is only one: `internal/declsyntax`, `internal/horosmodel`, `internal/horosfiles`, `verification/files/source/horos/go`, the four packages `tools/provider-matrix/go`, `report/go`, `capability/go` and `observation/go`, and `tools/horos-module/go` |
| Hand entries in `TestProductionArchitecture` | 8 | `semantics/golang/go`, `internal/golangsem`, `resolve/files/golang/provider/go`, `analyze/files/golang/provider/go`, `internal/typescriptimports`, `analyze/files/rust/provider/go`, `verification/files/horos-to-go/go` and `verification/logos/go`. They convert when pull requests #204 and #231 have merged |
| `roleExtensions` | 10 | Per package, the imports that one language needs and its role does not give: six horos stage providers, three TypeScript packages and the horos declaration recipe. An extension must be used, must not repeat the role and lifts no prohibition |
| `languageOwners` | 5 | The language of a package whose path gives no role: `internal/declsyntax`, `internal/horosmodel`, `internal/horosfiles` and `verification/files/source/horos/go` belong to horos, and `internal/typescriptimports` to TypeScript |
| `providerExceptions` | 15 | Every import of a stage provider package outside a composition, with its symbols and reason. Thirteen are accepted: the calls of finding 9. Two are open and both are Go's: `Contract` until #220 and `Admission` until #221 |
| `modelDependencies` and `modelExternalImports` | 8 and 10 | The external model packages that the production graph may hold, and the packages that may import each directly |

How the test realises each of the six rules:

| Rule | Realised as |
| --- | --- |
| 1 | A package in `foundationRules` imports exactly what its entry lists. It never imports a package outside that table, also when its entry lists one. So it imports no vocabulary, no contribution and no private part of a language |
| 2 | The four vocabulary shapes have role lists of foundations and their own language's vocabulary. `semantics/<language>/go` and `transform/syntax/<native>/symbol/go` may also import their language's lexical helper. A vocabulary never imports a contribution, and `semantics/<language>/go` never imports a semantic helper. Both prohibitions hold under a hand entry and under an extension as well |
| 3 | Each contribution shape has its own role list. A package that names a language imports no package of another language. A mapping names two. A role grants `internal/<language>` and `internal/<language>sem` only where the imported package is that language's helper by its own role. The role lists name a contribution in three places. A parser and a formatter import their own checker, and a semantic helper its own parser, each under rule 5. A mapping and a common mapping part import the source reader of their source language |
| 4 | Enforced for stage provider packages: an import by a package that is not a composition fails without an entry in the exception table. For the other contributions the role lists decide. No role lists a native target, a mapping, a law, a query service or an extractor. A role does list the source reader, as named under rule 3, and its importer may use any of its symbols. The check has no general statement that only a composition imports a contribution |
| 5 | The exception table. An entry names importer, provider, symbols and reason. A use of another symbol fails. An unused entry fails, and so does an unused symbol. An entry that is listed twice, names no symbol, gives no reason, names a package that is not a stage provider or names a composition as importer fails |
| 6 | A production package with neither a role nor an entry fails. A package named `internal/<name>` has a helper's role only when `syntax/<name>/go` exists |

Further failures:

- **Every allowance is used.** An entry and an extension are checked per
  package: an import that they allow and the package does not make fails. A
  role is checked per shape: each import it permits must be imported by at
  least one package that the role classifies. External allowances are checked
  too, per package in `modelExternalImports` and per shape for the two models
  that the law role carries.
- **Files that build constraints exclude.** The check reads every non-test
  file of a package, also one behind a build tag, one of another operating
  system or architecture, and a cgo file when cgo is off. `go list` leaves
  their imports out.
- **A package that is not listed.** An import of a module package that the
  package set does not hold fails, because nothing checks that package.
- **Dot import.** A dot import of a package of this module fails.

Four details of the implementation:

- The compositions `integration/go` and `verification/integration/go` are not
  classified at all. A package under `tools/` has a hand entry or the extractor
  role and needs no exception to import a provider. The test does not know a
  `cmd/**` composition. Adding one is an edit.
- An entry for a package that does not exist is ignored, so a package can be
  deleted without an edit to the test. The entries of the five deleted packages
  `transform/syntax/horos-to-<native>/mapping/go` and
  `verification/files/horos-to-go/mappings/go` are gone.
- A hand entry replaces the allow-list of a role and none of its
  prohibitions. `semantics/golang/go`, `internal/golangsem`, the two Go
  provider packages and `analyze/files/rust/provider/go` have both.
- A new language whose parser and formatter call its own checker adds two
  entries to the exception table. The role permits the import, and the
  exception still names the symbols.

Still open:

- **Go.** The two open entries of the exception table and the eight hand
  entries wait for draft pull request #231, which is stacked on open pull
  request #204. Neither is on `main`.

### What the check cannot establish

- **A decision made inline.** Record building `std::optional` imports nothing
  forbidden. That needs the contract change of #196 and a reviewer who asks the
  question above.
- **Whether a package holds one contribution.** That is a judgment. The scan
  lists candidates, and it over-reports: a setting key that a provider reads
  through a helper shows as a separate group.
- **Correct ownership.** A green run shows that the forbidden edges are absent
  ([AC-07](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-07)).
- **A re-exported symbol.** A helper that holds an exception can re-export the
  excepted symbol to its own importers. The
  [independent review](#independent-review) confirmed by seed that this
  passes.
- **An aliased type.** A package can alias an excepted type for another
  package. The review confirmed by seed that this passes too.
- **Use per package within a role.** Use is checked per shape. One package may
  import what a sibling of the same shape imports, whether it needs it or not.
- **An unused entry in `modelDependencies`.** The list of external packages
  that the graph may hold is not checked for use.
- **A package whose every file is excluded and that nothing imports.**
  `go list` does not list it, so it is neither classified nor checked on that
  build.

The [agent instructions](../AGENTS.md) say not to replace reasoning about
responsibilities with keyword heuristics or new lint gates. These rules do not
replace it. They enforce boundaries that the reasoning has already decided, in
the test that already enforces directed imports.

## How the check is built

This section described the plan. It now describes what
[#225](https://github.com/Bitspark/bittype/pull/225) and
[#244](https://github.com/Bitspark/bittype/pull/244) built. All files are
under `integration/go`.

- `architecture_test.go` holds `TestProductionArchitecture`, the tables
  `foundationRules`, `oneOfAKindRules`, `languageOwners` and `roleExtensions`,
  the eight hand entries and the external model allowances.
  `go list -deps -json` runs there and gives each package's imports, name,
  directory and files. The loader `architecturePackages` adds the imports and
  symbol uses of the files that the build constraints exclude.
- `architecture_roles_test.go` holds `packageRole`, which derives a package's
  rule from its path. A role is written once as patterns over `<language>`.
- `architecture_check_test.go` holds the check, the function
  `checkArchitecture` over package data.
- `architecture_exceptions_test.go` holds the exception table
  `providerExceptions`.
- Symbol use comes from `go/parser` over each package's non-test files:
  resolve the import names per file, then collect the selectors on them. Use is
  syntactic and nothing is type-checked. The standard library is enough, and
  `go.mod` gained no requirement.
- `architecture_synthetic_test.go` holds 36 tests that feed the check
  hand-built package sets. Each failing set breaks one rule and must produce
  that rule's message and no other. One test adds a third language, a pair
  into it, one more mapping and one more common mapping part to a conforming
  set, and expects no new line in the foundation table, the hand entries or
  the owners. One test writes a module at test time with a forbidden import
  behind a build tag, in a file of another operating system and in a cgo file,
  and runs the real loader with cgo off and on.
- A dot import of a module package fails, because it hides symbol use.

## Order of work

State on `main` at `71bfec9`, 12 November 2026.

1. **[#219](https://github.com/Bitspark/bittype/issues/219): name the symbols
   and fail unused allowances.** Done with
   [#225](https://github.com/Bitspark/bittype/pull/225); the issue is closed.
   It moved nothing. It made the 22 provider imports explicit: thirteen
   accepted, eight open until #220 and one open until #221. It also added the
   path rules for the four package shapes that the later steps create.
2. **[#220](https://github.com/Bitspark/bittype/issues/220): identities move to
   the vocabularies.** Done for Rust
   ([#229](https://github.com/Bitspark/bittype/pull/229)), TypeScript
   ([#233](https://github.com/Bitspark/bittype/pull/233)) and horos
   ([#234](https://github.com/Bitspark/bittype/pull/234)). Seven of the eight
   exceptions are removed. Not done for Go: draft pull request #231 is stacked
   on open pull request #204. The issue also gained the inline horos keys of
   finding 12, which are not done.
3. **[#221](https://github.com/Bitspark/bittype/issues/221): the query services
   get their own packages.** Done for TypeScript
   ([#233](https://github.com/Bitspark/bittype/pull/233)). Not done for Go: the
   same draft pull request #231. The one exception remains until it merges.
4. **[#222](https://github.com/Bitspark/bittype/issues/222): one mapping
   contribution per package.** Done for the four pair packages
   ([#226](https://github.com/Bitspark/bittype/pull/226),
   [#227](https://github.com/Bitspark/bittype/pull/227),
   [#230](https://github.com/Bitspark/bittype/pull/230),
   [#232](https://github.com/Bitspark/bittype/pull/232)) and for the
   verification laws ([#228](https://github.com/Bitspark/bittype/pull/228)).
   [#224](https://github.com/Bitspark/bittype/pull/224) prepared it with one
   test seam file per pair and one for the laws under `integration/go`.
   [#238](https://github.com/Bitspark/bittype/pull/238) divided the common
   package `transform/syntax/horos/mapping/go` into the adapter, which stays at
   that path, and the packages `pattern`, `declarations` and `naming`. It is a
   move and does not remove the profile gate.
   [#216](https://github.com/Bitspark/bittype/issues/216) is untouched and
   stays whole. [#207](https://github.com/Bitspark/bittype/issues/207), which
   writes the first families into mapping packages, is open.
5. **[#218](https://github.com/Bitspark/bittype/issues/218): rules by role.**
   Done with [#244](https://github.com/Bitspark/bittype/pull/244). A package
   gets its import rule from the role that its path gives, and a new package
   of an existing role needs no edit to the test. The step was narrowed as
   [rules to enforce](#rules-to-enforce) records: the foundation and kernel
   set keeps explicit entries, and so do nine packages of which there is only
   one. Not done: eight hand entries stay in the test until pull requests #204
   and #231 have merged. #218 stays open until then.

The steps ran as parallel lanes, as the lane map on #218 records. The path
rules for four package shapes went in with step 1, so that a new package of
one of those shapes was classified without an edit to the test while the
other lanes ran.

Each pull request that moved code reports the same three checks: the hosted
checks passed, the identity values or the generated declarations were rendered
before and after and were byte-identical, and a seeded forbidden import failed
the architecture test.
[bittype/ACX-009](architecture-cases/0009-package-boundaries.md#correction-and-qualification-status)
records this per pull request.

No other platform repository imports a bittype package, and bittype has no
release tag (checked on 11 November 2026). The moves cost nothing outside this
repository today. That changes when `libraries/service` starts to consume the
generation packages.

## Independent review

A reviewer that wrote none of the changes read the merged moves from the code,
in the range `d23542a..e83cf52`. The review is recorded in a
[comment on #218](https://github.com/Bitspark/bittype/issues/218#issuecomment-6809378108)
of 11 November 2026.

**Verdict.** They are moves. The review found no change of behaviour, of an
identity value or of generated output apart from three differences that the
pull requests had declared.

**Method.**

- **Declarations.** 237 production declarations were compared with the old
  ones. 133 are token-identical to an old declaration. The other 104 were
  read. Each difference is a rename, a qualifier or a declared difference.
- **Identities.** Two programs, one per revision, print 3352 lines: the horos,
  Rust and TypeScript profiles, keys, roles, relations, settings, pin and
  contracts, the five verification rules, and the description of every
  mapping library under three profiles. Both outputs have the same SHA-256.
- **Generated output.** The output of `TestSemanticPaginationMappings` for the
  four targets is identical at both revisions. The plan snapshots are
  unchanged.
- **Consumers.** The installed-consumer fixture, built with the module fetched
  from the module proxy at `e83cf52`, passes seven of its eight tests. The
  eighth needs the Rust helper.

**Not verified.** Anything that needs the Rust helper, the native compiler
job, and the race and fuzz steps. The hosted checks ran those on each pull
request.

**Defects.** The review found four, all in the check or in comments and none
in moved code. #244 corrected them.

1. The check did not read files that build constraints exclude. A forbidden
   import passed behind a build tag, in a `_linux.go` file when the test runs
   on Windows, and in a cgo file with cgo off.
2. The vocabulary rule accepted any name as a language, so a new
   `semantics/horosmodel/go` could import the private `internal/horosmodel`.
3. The entries of the five deleted packages were still in the list.
4. Four comments stated something that was no longer true.

**Limits.**

- The review confirmed by seed two things that no import check sees. They are
  under [what the check cannot establish](#what-the-check-cannot-establish): a
  re-exported symbol and an aliased type.
- The review's programs and outputs are in a scratch directory on the
  workstation, which the comment names. They are not in the repository, so
  the comparison cannot be repeated from the repository.
- The range ends before #244. That pull request changes test code only. Its
  qualification is in
  [bittype/ACX-009](architecture-cases/0009-package-boundaries.md#qualification).

## Repeating the scan

Two questions, answered from source without type checking:

1. Per package, which groups of top-level declarations never reference one
   another? Two or more groups that each export something are grouping
   candidates.
2. Per import of a module package, which exported symbols does the importer
   use? An importer that uses one identity from a provider package is a
   vocabulary candidate.

The tool's source is in [#218](https://github.com/Bitspark/bittype/issues/218).
It is a discovery aid and not a check. Read every candidate in source before
recording it.
