# Package boundaries and forbidden dependencies

**Implemented. State of 12 November 2026.** The five steps of the
[order of work](#order-of-work) are on `main`. Go's part of steps 2 and 3
merged with pull request [#246](https://github.com/Bitspark/bittype/pull/246).
The horos names of finding 12 merged with pull requests
[#247](https://github.com/Bitspark/bittype/pull/247) and
[#248](https://github.com/Bitspark/bittype/pull/248). Every confirmed
finding that the order of work covers is corrected. Finding 7 is open and
belongs to [#196](https://github.com/Bitspark/bittype/issues/196). Findings 9,
10 and 11 are candidates and are not confirmed. The
[findings table](#findings-on-main-at-3bf45c3) gives the state of each
finding, and [what the test enforces now](#what-the-test-enforces-now) says
which rules are checked and how. Two
[independent reviews](#independent-review) read moves from the code: one the
moves that merged on 11 November 2026, and one the first commit of #246.
Neither read the second commit of #246, #247 or #248. The reasoning sections
are unchanged since the proposal.

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
column "State" describes `main` at `abfc948`, after pull request #248, on
12 November 2026. It was read in source at `1a5aa2d` on 11 November 2026, and
finding 8 again at `71bfec9`. Only test files and documents changed between
those two commits. Findings 5, 6 and 9 were read again at `c8d0344`, the merge
of #246, and finding 12 at `abfc948`. Finding 12 was found on 11 November
2026, after the scan.

| # | Where | Kind | Assessment | State | Tracked |
| --- | --- | --- | --- | --- | --- |
| 1 | `transform/syntax/horos-to-{typescript,golang,rust,cpp}/mapping/go` each hold Text, Bool, List, Record and Declarations. Go also holds Optional | Grouping | Confirmed | Corrected with [#226](https://github.com/Bitspark/bittype/pull/226), [#227](https://github.com/Bitspark/bittype/pull/227), [#230](https://github.com/Bitspark/bittype/pull/230) and [#232](https://github.com/Bitspark/bittype/pull/232). Each mapping is its own package, `transform/syntax/horos-to-<native>/mapping/<name>/go`, and exports `Library` | [#222](https://github.com/Bitspark/bittype/issues/222) |
| 2 | `Libraries(profile)` in each of those packages lists the package's own mappings | Grouping | Confirmed | Corrected with the same four pull requests. `Libraries` is deleted and a consumer registers each library | [#222](https://github.com/Bitspark/bittype/issues/222) |
| 3 | `verification/files/horos-to-go/mappings/go` holds the Text, Boolean, List and Record laws | Grouping | Confirmed | Corrected with [#228](https://github.com/Bitspark/bittype/pull/228). Each law is its own package, `verification/files/horos-to-go/mappings/{scalar,list,optional,record}/go` | [#222](https://github.com/Bitspark/bittype/issues/222) |
| 4 | `transform/syntax/horos/mapping/go` holds the profile gate, the pattern constructors, the declaration recipe and the naming function | Grouping | Confirmed | Corrected with [#238](https://github.com/Bitspark/bittype/pull/238). The adapter `Library` with the profile gate stays at that path. The pattern constructors, the declaration recipe and the naming function are the packages `pattern`, `declarations` and `naming` under `transform/syntax/horos/mapping/<name>/go`. The division is a move. The profile gate is unchanged | [#222](https://github.com/Bitspark/bittype/issues/222), done-when 4; the gate in [#216](https://github.com/Bitspark/bittype/issues/216), which is untouched |
| 5 | `analyze/files/{golang,typescript}/provider/go` each hold the analysis provider and a native query service that share no declaration | Grouping | Confirmed | Corrected. TypeScript with [#233](https://github.com/Bitspark/bittype/pull/233): the query service is `analyze/files/typescript/queries/go`. Go with [#246](https://github.com/Bitspark/bittype/pull/246): the query service is `analyze/files/golang/queries/go`. `verification/files/horos-to-go/go` imports that package and no longer the analysis provider, under an accepted exception that names `Admission` | [#221](https://github.com/Bitspark/bittype/issues/221), closed |
| 6 | Eight imports of a provider package use nothing but an identity. horos has no vocabulary package for its semantic names | Identity kept with an implementation | Confirmed | Corrected. Rust with [#229](https://github.com/Bitspark/bittype/pull/229), TypeScript with [#233](https://github.com/Bitspark/bittype/pull/233), horos with [#234](https://github.com/Bitspark/bittype/pull/234), which adds `semantics/horos/go`, and Go with [#246](https://github.com/Bitspark/bittype/pull/246). All eight imports are gone, and no stage provider package exports `Contract()`. The horos names that were spelled inline are finding 12 | [#220](https://github.com/Bitspark/bittype/issues/220) |
| 7 | Record decides optional fields itself in Go, Rust and C++ | Coupling | Confirmed; the contract change is designed | Open. [#227](https://github.com/Bitspark/bittype/pull/227) removed the package variable that Go's Record shared with Optional: the consumer now creates a helper key and passes it to both libraries. Record still decides the representation of an optional field in all three targets | [#196](https://github.com/Bitspark/bittype/issues/196), item 6 |
| 8 | The existing check passes unused allowances and cannot see which symbols an import uses | Enforcement gap | Confirmed | Corrected with [#225](https://github.com/Bitspark/bittype/pull/225), which names the symbols of every provider import, fails unused allowances and removed 94 of them. [#244](https://github.com/Bitspark/bittype/pull/244) checks the allowances of a role for use as well, and reads the files that build constraints exclude | [#219](https://github.com/Bitspark/bittype/issues/219), closed; [#218](https://github.com/Bitspark/bittype/issues/218) |
| 9 | Thirteen direct calls into a same-language provider: `Check` from parse, format and the horos resolver; `Parse` from the Rust and TypeScript semantic helpers | Coupling candidate | Read as legitimate: the caller needs its own grammar's answer. `LAYOUT.md` permits them | Not confirmed. The thirteen imports are accepted entries of the exception table, each with its symbol and reason. The table's fourteenth entry is not one of these calls: it is the verification pair's use of `Admission` under finding 5. Whether the admission function should move to a language-private package is not decided | Recorded in [#219](https://github.com/Bitspark/bittype/issues/219); the candidate stays on [#218](https://github.com/Bitspark/bittype/issues/218) |
| 10 | `Exported` appends `_` to `Self` for all four targets. `Self` is a keyword only in Rust | Coupling candidate | Not confirmed | Not confirmed and unchanged. Decided on #222 on 11 November 2026: the rule stays where it is and is a question for #216 | [#216](https://github.com/Bitspark/bittype/issues/216) |
| 11 | `CheckSingleScope` states one placement's rule in the shared target construction package | Grouping candidate | Not confirmed | Not confirmed and unchanged. #222 did not decide it | With the second design point on [#161](https://github.com/Bitspark/bittype/issues/161) |
| 12 | horos-owned keys are spelled inline where they are used: `horos/module` and `horos/import` in `internal/horosfiles`, and `horos/type-definition`, `horos/sort/proper` and `horos/identity/module` in `transform/syntax/horos/source/go` | Identity without a single owner | Identified on 11 November 2026 by the change that added `semantics/horos/go` | Corrected. [#247](https://github.com/Bitspark/bittype/pull/247) published the five keys from `semantics/horos/go` as `ModuleKind()`, `ImportRelation()`, `TypeDefinitionKind()`, `ProperSort()` and `ModuleIdentity()`. That change found three further names that a writer and a separate reader each spelled, on 12 November 2026: the descriptor operator, the term encoding and the two profile names. [#248](https://github.com/Bitspark/bittype/pull/248) gave them one owner in the same package. No production package outside the vocabulary spells any of the eight | [#220](https://github.com/Bitspark/bittype/issues/220), added to its scope |

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
configuration and no other production package uses it. On that day Go still
published `Contract()` from `resolve/files/golang/provider/go` and from
`analyze/files/golang/provider/go`.

Finding 6 on 12 November 2026, with pull request [#246](https://github.com/Bitspark/bittype/pull/246).
`semantics/golang/go` exports `ResolutionContract()` and `AnalysisContract()`
in `contracts.go`. The two Go provider packages no longer export `Contract()`,
and the Go analysis provider no longer imports the Go resolver. No alias
remains at an old path. Go has no `AnalysisContractFor(settings)` on `main`.
Open pull request #204 adds it, because it binds analysis to an import
environment that only #204 has.

The Go move changed its place in the order on 12 November 2026. It was first
prepared as draft pull request #231 on top of #204. #231 is closed. The move
went to `main` first, as #246, because
[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
puts the move of a misplaced responsibility before an addition. #204 rebases
onto it and blocks nothing in this plan. The reasons are in a
[comment on #204](https://github.com/Bitspark/bittype/pull/204#issuecomment-6812747024).

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

Finding 12 on 12 November 2026. Pull request [#247](https://github.com/Bitspark/bittype/pull/247)
published the five keys from `semantics/horos/go`. Their writers in
`internal/horosmodel` and `internal/horosfiles` and their readers in
`internal/horosfiles` and `transform/syntax/horos/source/go` call the published
functions. The change that moved them found three further names that a writer
and a separate reader each spelled. Pull request
[#248](https://github.com/Bitspark/bittype/pull/248) gave each one owner
in the same package:

- **The descriptor operator.** `DescriptorOperator(constructor)` is
  `horos/descriptor/<constructor>`, and `DescriptorConstructor(operator)` is
  its inverse. `internal/horosmodel` writes the operator and the source reader
  reads it.
- **The term encoding.** `TermEncoding()` is `ontos/codec-v1`. ontos owns the
  codec. That horos terms carry it is an agreement between the horos writer
  and the horos reader, so the horos vocabulary states it.
- **The two profile names.** `ResolutionProfile` and `AnalysisProfile` were
  published before. `internal/horosfiles` repeated both as string literals and
  now compares with the constants.

Both pull requests are moves: every key value and every string is unchanged.
Read in source at `abfc948`. The horos syntax has node kinds with the same
spelling as two of these names, `type-definition` and `descriptor`. Those are
other names: kinds of tree nodes, not keys of the semantic graph.

Other horos names stay where `internal/horosmodel` writes them. They are the
entity kinds `type-parameter`, `kind-definition` and `kind-parameter`, the
parameter sorts `proper` and `constructor`, the operators `kind/<form>`,
`sort/arrow`, `membership/assumption`, `membership/proven`,
`identity/declaration` and `identity/parameter`, and the diagnostic codes. Each
has a production writer and no production reader that names it. Tests read
some of them. These are not instances of finding 12. A name needs a single
owner when two parties must agree on it, and here there is no second party.
One of them becomes an instance when a production package other than its
writer reads it by name.

Not examined, and so neither an instance nor cleared:

- **The encoding keys of the other languages.** The Go, TypeScript and Rust
  models write `json-1` and `utf8` in `internal/golangsem/export.go`,
  `internal/typescriptsem/compiler.cjs` and
  `internal/rustsem/helper/src/main.rs`. Their readers were not looked for.
- **The other names of the other languages.** One pair was seen and not
  assessed: the TypeScript compiler helper `internal/typescriptsem/compiler.cjs`
  writes the entity kind `module` and the binding roles `name` and `import`,
  and `format/files/typescript/provider/go` reads them as `n.Key("module")`,
  `n.Key("name")` and `n.Key("import")`.
- **The tag prefixes of the canonical horos value.** `horos/decl/` and `horos/`
  are spelled in three production packages: `internal/horosmodel` in
  `SourceModules`, and `resolve/files/horos/provider/go` and
  `verification/files/source/horos/go`, which each configure the same
  `Frontend` with them. All three write the canonical value, and the external
  horos library reads it. They are data tags of that library's format and not
  semantic keys. Whether they need one owner in bittype is not decided.
- **The node kinds of the horos syntax.** `syntax/horos/go` and
  `internal/declsyntax` each declare them, and neither imports the other.

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
explicit entries for twelve packages of which there is only one. Its reason is
that a pattern for a single package would not be honest.

### What the test enforces now

Since [#244](https://github.com/Bitspark/bittype/pull/244) a package gets its
import rule from the role that its path gives. This section was read in
`integration/go/architecture_test.go`, `architecture_roles_test.go`,
`architecture_check_test.go` and `architecture_exceptions_test.go` at
`c8d0344`, the merge of pull request
[#246](https://github.com/Bitspark/bittype/pull/246). The counts were made
there. #247 and #248 changed none of these files and added no package. The
counts at `71bfec9`, before the Go move, are in
[bittype/ACX-009](architecture-cases/0009-package-boundaries.md#corrections-on-main).

Every production package of the module is classified in one of three ways. A
package with none of them fails as unclassified.

- **By its role.** `packageRole` derives the rule from the path. 84 of the 127
  production packages are classified this way. The
  [layout](../LAYOUT.md#dependency-direction) lists the roles and what each
  may import.
- **By its entry in `foundationRules`.** The entry is the exact list of the
  package's module imports. 31 packages have one.
- **By a hand entry.** The entry is again the exact list. 12 packages have
  one, and all are in `oneOfAKindRules`. `TestProductionArchitecture` holds no
  entry of its own since #246. A hand entry takes precedence over a role.

The tables:

| Table | Entries | Holds |
| --- | --- | --- |
| `foundationRules` | 31 | The carriers, the contracts, the engine, the shared transform, check and print packages, and the five private helpers of no single language: `internal/tree`, `internal/grammar`, `internal/lex`, `internal/parseutil` and `internal/concrete` |
| `oneOfAKindRules` | 12 | Packages of which there is only one: `internal/declsyntax`, `internal/horosmodel`, `internal/horosfiles`, `internal/typescriptimports`, `verification/files/source/horos/go`, `verification/files/horos-to-go/go`, `verification/logos/go`, the four packages `tools/provider-matrix/go`, `report/go`, `capability/go` and `observation/go`, and `tools/horos-module/go` |
| `roleExtensions` | 10 | Per package, the imports that one language needs and its role does not give: six horos stage providers, three TypeScript packages and the horos declaration recipe. An extension must be used, must not repeat the role and lifts no prohibition |
| `languageOwners` | 5 | The language of a package whose path gives no role: `internal/declsyntax`, `internal/horosmodel`, `internal/horosfiles` and `verification/files/source/horos/go` belong to horos, and `internal/typescriptimports` to TypeScript |
| `providerExceptions` | 14 | Every import of a stage provider package or of a query service package outside a composition, with its symbols and reason. All fourteen are accepted and no entry names an issue. Thirteen name a stage provider: the calls of finding 9. One names a query service: `verification/files/horos-to-go/go` takes `Admission` from `analyze/files/golang/queries/go`. #246 deleted the two open entries, Go's `Contract` and `Admission` from the analysis provider |
| `modelDependencies` and `modelExternalImports` | 8 and 10 | The external model packages that the production graph may hold, and the packages that may import each directly |

How the test realises each of the six rules:

| Rule | Realised as |
| --- | --- |
| 1 | A package in `foundationRules` imports exactly what its entry lists. It never imports a package outside that table, also when its entry lists one. So it imports no vocabulary, no contribution and no private part of a language |
| 2 | The four vocabulary shapes have role lists of foundations and their own language's vocabulary. `semantics/<language>/go` and `transform/syntax/<native>/symbol/go` may also import their language's lexical helper. A vocabulary never imports a contribution, and `semantics/<language>/go` never imports a semantic helper. Both prohibitions hold under a hand entry and under an extension as well |
| 3 | Each contribution shape has its own role list. A package that names a language imports no package of another language. A mapping names two. A role grants `internal/<language>` and `internal/<language>sem` only where the imported package is that language's helper by its own role. The role lists name a contribution in three places. A parser and a formatter import their own checker, and a semantic helper its own parser, each under rule 5. A mapping and a common mapping part import the source reader of their source language |
| 4 | Enforced for stage provider packages and for query service packages, the two kinds that `selectedKind` names: an import of either by a package that is not a composition fails without an entry in the exception table. For the other contributions the role lists decide. No role lists a native target, a mapping, a law, a query service or an extractor. One explicit entry lists a query service: `verification/files/horos-to-go/go` names `analyze/files/golang/queries/go`, and its exception names `Admission` only. A role does list the source reader, as named under rule 3, and its importer may use any of its symbols. The check has no general statement that only a composition imports a contribution |
| 5 | The exception table. An entry names importer, provider, symbols and reason. A use of another symbol fails. An unused entry fails, and so does an unused symbol. An entry that is listed twice, names no symbol, gives no reason, names a package that is neither a stage provider nor a query service, or names a composition as importer fails |
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
  role and needs no exception to import a provider or a query service. The
  test does not know a `cmd/**` composition. Adding one is an edit.
- An entry for a package that does not exist is ignored, so a package can be
  deleted without an edit to the test. The entries of the five deleted packages
  `transform/syntax/horos-to-<native>/mapping/go` and
  `verification/files/horos-to-go/mappings/go` are gone.
- A hand entry replaces the allow-list of a role and none of its
  prohibitions. No package has both a role and a hand entry since #246.
- A new language whose parser and formatter call its own checker adds two
  entries to the exception table. The role permits the import, and the
  exception still names the symbols.

Nothing in the test waits for Go any more. Pull request #246 deleted the two
open entries of the exception table and converted the eight hand entries that
`TestProductionArchitecture` held. `semantics/golang/go`, `internal/golangsem`,
`resolve/files/golang/provider/go`, `analyze/files/golang/provider/go` and
`analyze/files/rust/provider/go` take their role. `internal/typescriptimports`,
`verification/files/horos-to-go/go` and `verification/logos/go` are entries of
`oneOfAKindRules`.

The conversion is not a pure move of the rules. The
[review of #246](#the-go-move-pull-request-246) compared each package's hand
entry with its role and found three differences:

- **A gap, corrected.** After the first commit the verification pair could use
  any symbol of the Go query service. Its exception on `main` had named
  `Admission` only. The second commit of #246 extended the exception rule to
  query services and gave the pair an accepted entry for `Admission`.
- **Three Go packages may import more than their hand entry listed.**
  `internal/golangsem` gains `contracts/parse/go` and `files/go`. The Go
  resolver gains `trees/representation/go`. The Go vocabulary gains
  `contracts/resolution/go`, which it now imports, and `internal/golang`. Each
  is what the role gives the same package of TypeScript and Rust. None is a
  package of another language or a contribution. This was accepted on #246:
  one rule per role is the design decided on #218, and the hand entries were
  held only so that two pull requests would merge cleanly.
- **Unused allowances of the five converted packages.** An allowance that one
  of them does not use no longer fails for that package. A role is checked for
  use per shape. This was accepted for the same reason, and it already held
  for every other language.

The last two are limits of the check and not defects that were fixed. They are
listed again under
[what the check cannot establish](#what-the-check-cannot-establish).

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
  Since #246 this holds for the Go packages too: `internal/golangsem` may
  import `contracts/parse/go` and `files/go`, the Go resolver
  `trees/representation/go` and the Go vocabulary `internal/golang`, and none
  of them does.
- **A selected contribution of another kind.** The exception rule names two
  kinds, a stage provider and a query service. An import of a source reader
  is permitted by a role, and its importer may use any of its symbols.
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
[#244](https://github.com/Bitspark/bittype/pull/244) built, with the tables
and the exception rule as #246 left them. All files are under
`integration/go`.

- `architecture_test.go` holds `TestProductionArchitecture`, the tables
  `foundationRules`, `oneOfAKindRules`, `languageOwners` and `roleExtensions`,
  and the external model allowances.
  `go list -deps -json` runs there and gives each package's imports, name,
  directory and files. The loader `architecturePackages` adds the imports and
  symbol uses of the files that the build constraints exclude.
- `architecture_roles_test.go` holds `packageRole`, which derives a package's
  rule from its path. A role is written once as patterns over `<language>`. It
  also holds `selectedKind`, which names the two kinds of package whose import
  outside a composition needs an exception.
- `architecture_check_test.go` holds the check, the function
  `checkArchitecture` over package data.
- `architecture_exceptions_test.go` holds the exception table
  `providerExceptions`.
- Symbol use comes from `go/parser` over each package's non-test files:
  resolve the import names per file, then collect the selectors on them. Use is
  syntactic and nothing is type-checked. The standard library is enough, and
  `go.mod` gained no requirement.
- `architecture_synthetic_test.go` holds 37 tests that feed the check
  hand-built package sets. Each failing set breaks one rule and must produce
  that rule's message and no other. One test adds a third language, a pair
  into it, one more mapping and one more common mapping part to a conforming
  set, and expects no new line in the foundation table, the hand entries or
  the owners. One test writes a module at test time with a forbidden import
  behind a build tag, in a file of another operating system and in a cgo file,
  and runs the real loader with cgo off and on. One test, added with #246,
  imports a query service from a package that is not a composition: without an
  exception, with one, and with a use of a symbol that the exception does not
  name.
- A dot import of a module package fails, because it hides symbol use.

## Order of work

State on `main` at `abfc948`, after pull request #248, 12 November 2026.

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
   ([#234](https://github.com/Bitspark/bittype/pull/234)), and for Go
   ([#246](https://github.com/Bitspark/bittype/pull/246)). All eight
   exceptions are removed. The issue also gained the inline horos keys of
   finding 12. Done with [#247](https://github.com/Bitspark/bittype/pull/247)
   for the five keys and with [#248](https://github.com/Bitspark/bittype/pull/248)
   for three further names that the first change found.
3. **[#221](https://github.com/Bitspark/bittype/issues/221): the query services
   get their own packages.** Done for TypeScript
   ([#233](https://github.com/Bitspark/bittype/pull/233)) and for Go
   ([#246](https://github.com/Bitspark/bittype/pull/246)). The one
   exception on the analysis provider is removed. The verification pair
   imports the query service under an accepted exception that names
   `Admission`. #246 closed the issue.
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
   set keeps explicit entries, and so do twelve packages of which there is
   only one. Eight hand entries stayed in the test for Go, and
   [#246](https://github.com/Bitspark/bittype/pull/246) converted them.

No step is recorded as not done. #218 and #220 were still open on GitHub when
this was written. Finding 7 and the three candidates are not steps of this
order: they stay with #196, #218, #216 and #161 as the findings table says.
Open pull request #204 is not part of this plan. It adds imported packages to
Go's semantics and will add `AnalysisContractFor(settings)` to
`semantics/golang/go` when it rebases.

The steps ran as parallel lanes, as the lane map on #218 records. The path
rules for four package shapes went in with step 1, so that a new package of
one of those shapes was classified without an edit to the test while the
other lanes ran.

Go's part of steps 2 and 3 ran last and in a changed order. It was first
prepared as draft pull request #231 on top of open pull request #204. #231 is
closed, the move merged as #246, and #204 rebases onto it. The paragraphs on
finding 6 under the [findings table](#findings-on-main-at-3bf45c3) give the
reason.

Each pull request that moved code up to #246 reports the same three checks:
the hosted checks passed, the identity values or the generated declarations
were rendered before and after and were byte-identical, and a seeded forbidden
import failed the architecture test. For #246 the hosted `docs`, `contracts`
and `native` checks passed on both of its commits, `7ecdaa1` and `052a3e3`.
#247 and #248 add no import between packages of the module and seeded none.
Their authors compared the published semantic graph before and after, and the
hosted checks passed on both.
[bittype/ACX-009](architecture-cases/0009-package-boundaries.md#correction-and-qualification-status)
records this per pull request.

No other platform repository imports a bittype package, and bittype has no
release tag (checked on 11 November 2026). The moves cost nothing outside this
repository today. That changes when `libraries/service` starts to consume the
generation packages.

## Independent review

Two reviews read moves from the code. The first read the moves that merged on
11 November 2026. The second read the Go move of #246 and is
[below](#the-go-move-pull-request-246). No independent review read #247 or
#248.

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
- The range also ends before #246, which moves production code: Go's contract
  identities and Go's query service. This review did not read that move. The
  second review did.

### The Go move, pull request #246

A reviewer that wrote none of #246 read its first commit, `7ecdaa1`, against
`main` on 12 November 2026. The body of
[#246](https://github.com/Bitspark/bittype/pull/246) records the review
and the decision on each finding.

**Verdict.** The production code is a move, and the test edits change no
assertion. The conversion of the architecture test is not a pure move.

**Findings and decisions.**

1. **A gap, corrected.** The verification pair could use any symbol of the Go
   query service. On `main` its exception had named `Admission` only. The
   review seeded the pair calling `NewQueries`: it passed on `7ecdaa1` and
   failed on `main`. The second commit of #246, `052a3e3`, extended the
   exception rule to query services and gave the pair an accepted entry that
   names `Admission`. Its author reports two seeds on that commit: the pair
   calling `NewQueries` fails, and removing the entry fails. A synthetic test
   covers the rule.
2. **Three Go packages may import more than their hand entry listed.**
   Accepted. `internal/golangsem` gains `contracts/parse/go` and `files/go`,
   the Go resolver gains `trees/representation/go`, and the Go vocabulary
   gains `contracts/resolution/go`, which it now imports, and
   `internal/golang`. Each is what the role gives the same package of
   TypeScript and Rust. The review seeded the four imports that no package
   makes: each passed on `7ecdaa1` and failed on `main`.
3. **An unused allowance is checked per role for the five converted
   packages**, and no longer per package. Accepted for the same reason.
4. **Two comments and the wording of the identity proof.** A comment in the
   synthetic test and one in `contracts.go` had become false, and the second
   commit corrected them. The proof was worded as if the raw output were
   stable, and `%#v` prints a pointer address. The pull request body now says
   that the address is replaced by a fixed word before the comparison.

Findings 2 and 3 are limits of the check. They are not defects that were
fixed.

**Method, as far as its logs show.** The reviewer compared the hand entry of
each converted package with its role and seeded imports into the real tree.
Nine forbidden imports each failed `TestProductionArchitecture` on `7ecdaa1`,
one of them behind a build tag that is never set. The pull request body gives
their number as eight. It also says that the reviewer repeated the comparison
of the contract values with its own program and found them equal.

**Limits.**

- The reviewer did not read the second commit, `052a3e3`. The correction of
  finding 1 is qualified by its author and by the hosted checks only.
- The review's programs and logs are in
  `C:/Development/bitspark/scratch-bittype-review-246/` on the workstation.
  They are not in the repository, so the review cannot be repeated from the
  repository.
- The review is recorded in the pull request body, which the author of #246
  wrote. There is no separate report by the reviewer in the repository or on
  an issue.

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
