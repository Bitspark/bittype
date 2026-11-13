# Package boundaries and forbidden dependencies

**Implemented. State of 12 November 2026.** The five steps of the
[order of work](#order-of-work) are on `main`. Go's part of steps 2 and 3
merged with pull request [#246](https://github.com/Bitspark/bittype/pull/246).
The horos names of finding 12 merged with pull requests
[#247](https://github.com/Bitspark/bittype/pull/247),
[#248](https://github.com/Bitspark/bittype/pull/248) and
[#250](https://github.com/Bitspark/bittype/pull/250). Issue
[#251](https://github.com/Bitspark/bittype/issues/251) then examined two
groups of names that this document had listed as not examined. It confirmed an
instance in each, and pull requests
[#254](https://github.com/Bitspark/bittype/pull/254) and
[#255](https://github.com/Bitspark/bittype/pull/255) corrected them. Every
confirmed finding that the order of work covers is corrected. Finding 7 is
open and
belongs to [#196](https://github.com/Bitspark/bittype/issues/196). Findings 9
and 11 are candidates and are not confirmed. Finding 10 was answered in
principle on 12 November 2026 and is not corrected. The
[findings table](#findings-on-main-at-3bf45c3) gives the state of each
finding, and [what the test enforces now](#what-the-test-enforces-now) says
which rules are checked and how. Two
[independent reviews](#independent-review) read moves from the code: one the
moves that merged on 11 November 2026, and one the first commit of #246.
Neither read the second commit of #246, #247, #248, #250, #254 or #255. The
reasoning sections are unchanged since the proposal.

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
column "State" describes `main` at `f821eb3`, after pull request #255, on
12 November 2026.
It was read in source at `1a5aa2d` on 11 November 2026, and finding 8 again at
`71bfec9`. Only test files and documents changed between those two commits.
Findings 5, 6 and 9 were read again at `c8d0344`, the merge of #246. Finding
12 was read at `abfc948`, the merge of #248, and on the head of #250,
`cadf495`. Finding 12 was found on 11 November 2026, after the scan.

| # | Where | Kind | Assessment | State | Tracked |
| --- | --- | --- | --- | --- | --- |
| 1 | `transform/syntax/horos-to-{typescript,golang,rust,cpp}/mapping/go` each hold Text, Bool, List, Record and Declarations. Go also holds Optional | Grouping | Confirmed | Corrected with [#226](https://github.com/Bitspark/bittype/pull/226), [#227](https://github.com/Bitspark/bittype/pull/227), [#230](https://github.com/Bitspark/bittype/pull/230) and [#232](https://github.com/Bitspark/bittype/pull/232). Each mapping is its own package, `transform/syntax/horos-to-<native>/mapping/<name>/go`, and exports `Library` | [#222](https://github.com/Bitspark/bittype/issues/222) |
| 2 | `Libraries(profile)` in each of those packages lists the package's own mappings | Grouping | Confirmed | Corrected with the same four pull requests. `Libraries` is deleted and a consumer registers each library | [#222](https://github.com/Bitspark/bittype/issues/222) |
| 3 | `verification/files/horos-to-go/mappings/go` holds the Text, Boolean, List and Record laws | Grouping | Confirmed | Corrected with [#228](https://github.com/Bitspark/bittype/pull/228). Each law is its own package, `verification/files/horos-to-go/mappings/{scalar,list,optional,record}/go` | [#222](https://github.com/Bitspark/bittype/issues/222) |
| 4 | `transform/syntax/horos/mapping/go` holds the profile gate, the pattern constructors, the declaration recipe and the naming function | Grouping | Confirmed | Corrected with [#238](https://github.com/Bitspark/bittype/pull/238). The adapter `Library` with the profile gate stays at that path. The pattern constructors, the declaration recipe and the naming function are the packages `pattern`, `declarations` and `naming` under `transform/syntax/horos/mapping/<name>/go`. The division is a move and left the profile gate unchanged. The change for #216 removed the gate on 12 November 2026 | [#222](https://github.com/Bitspark/bittype/issues/222), done-when 4; the gate in [#216](https://github.com/Bitspark/bittype/issues/216) and [bittype/ACX-008](architecture-cases/0008-structural-profile-gate.md) |
| 5 | `analyze/files/{golang,typescript}/provider/go` each hold the analysis provider and a native query service that share no declaration | Grouping | Confirmed | Corrected. TypeScript with [#233](https://github.com/Bitspark/bittype/pull/233): the query service is `analyze/files/typescript/queries/go`. Go with [#246](https://github.com/Bitspark/bittype/pull/246): the query service is `analyze/files/golang/queries/go`. `verification/files/horos-to-go/go` imports that package and no longer the analysis provider, under an accepted exception that names `Admission` | [#221](https://github.com/Bitspark/bittype/issues/221), closed |
| 6 | Eight imports of a provider package use nothing but an identity. horos has no vocabulary package for its semantic names | Identity kept with an implementation | Confirmed | Corrected. Rust with [#229](https://github.com/Bitspark/bittype/pull/229), TypeScript with [#233](https://github.com/Bitspark/bittype/pull/233), horos with [#234](https://github.com/Bitspark/bittype/pull/234), which adds `semantics/horos/go`, and Go with [#246](https://github.com/Bitspark/bittype/pull/246). All eight imports are gone, and no stage provider package exports `Contract()`. The horos names that were spelled inline are finding 12 | [#220](https://github.com/Bitspark/bittype/issues/220) |
| 7 | Record decides optional fields itself in Go, Rust and C++ | Coupling | Confirmed; the contract change is designed | Open. [#227](https://github.com/Bitspark/bittype/pull/227) removed the package variable that Go's Record shared with Optional: the consumer now creates a helper key and passes it to both libraries. Record still decides the representation of an optional field in all three targets | [#196](https://github.com/Bitspark/bittype/issues/196), item 6 |
| 8 | The existing check passes unused allowances and cannot see which symbols an import uses | Enforcement gap | Confirmed | Corrected with [#225](https://github.com/Bitspark/bittype/pull/225), which names the symbols of every provider import, fails unused allowances and removed 94 of them. [#244](https://github.com/Bitspark/bittype/pull/244) checks the allowances of a role for use as well, and reads the files that build constraints exclude | [#219](https://github.com/Bitspark/bittype/issues/219), closed; [#218](https://github.com/Bitspark/bittype/issues/218) |
| 9 | Thirteen direct calls into a same-language provider: `Check` from parse, format and the horos resolver; `Parse` from the Rust and TypeScript semantic helpers | Coupling candidate | Read as legitimate: the caller needs its own grammar's answer. `LAYOUT.md` permits them | Not confirmed. The thirteen imports are accepted entries of the exception table, each with its symbol and reason. The table's fourteenth entry is not one of these calls: it is the verification pair's use of `Admission` under finding 5. Whether the admission function should move to a language-private package is not decided | Recorded in [#219](https://github.com/Bitspark/bittype/issues/219); the candidate stays on [#218](https://github.com/Bitspark/bittype/issues/218) |
| 10 | `Exported` appends `_` to `Self` for all four targets. `Self` is a keyword only in Rust | Coupling | Answered in principle on 12 November 2026: a reserved word is the target's, so the rule is in the wrong place | Not corrected, and the code is unchanged. Decided on #222 on 11 November 2026: the rule stays where it is and is a question for #216. #216 did not decide it. The [answer on #268](https://github.com/Bitspark/bittype/issues/268#issuecomment-6837825043) gives a reserved word to the target, because it is a fact about one language. The [naming design](naming.md#reserved-words) states it with the rest of the design: generation records a naming intent, and each language realizes the spelling. [#274](https://github.com/Bitspark/bittype/issues/274) delivers it with #268. Until then the function appends `_` to `Self` in every target | [#268](https://github.com/Bitspark/bittype/issues/268), [#274](https://github.com/Bitspark/bittype/issues/274); pending in [bittype/ACX-004](architecture-cases/0004-generation-file-syntax.md) |
| 11 | `CheckSingleScope` states one placement's rule in the shared target construction package | Grouping candidate | Not confirmed | Not confirmed and unchanged. #222 did not decide it | With the second design point on [#161](https://github.com/Bitspark/bittype/issues/161) |
| 12 | horos-owned keys are spelled inline where they are used: `horos/module` and `horos/import` in `internal/horosfiles`, and `horos/type-definition`, `horos/sort/proper` and `horos/identity/module` in `transform/syntax/horos/source/go` | Identity without a single owner | Identified on 11 November 2026 by the change that added `semantics/horos/go` | Corrected. [#247](https://github.com/Bitspark/bittype/pull/247) published the five keys from `semantics/horos/go` as `ModuleKind()`, `ImportRelation()`, `TypeDefinitionKind()`, `ProperSort()` and `ModuleIdentity()`. That change found three further names that a writer and a separate reader each spelled, on 12 November 2026: the descriptor operator, the term encoding and the two profile names. [#248](https://github.com/Bitspark/bittype/pull/248) gave them one owner in the same package. The documents change of #218 found one more on the same day: the two tag prefixes of canonical horos values, spelled in three packages. [#250](https://github.com/Bitspark/bittype/pull/250) moved them into that package, and [#285](https://github.com/Bitspark/bittype/pull/285) moved them on to `syntax/horos/go` on 13 November 2026, where the parser and the formatter can reach them. No production package outside the two vocabularies spells any of these names. Two further instances of the same kind lie outside the horos semantic names. [#251](https://github.com/Bitspark/bittype/issues/251) confirmed them, and they are corrected: the horos syntax node kinds, attributes and roles were declared twice ([#254](https://github.com/Bitspark/bittype/pull/254)), and sixteen TypeScript graph names were read in Go by spelling ([#255](https://github.com/Bitspark/bittype/pull/255)) | [#220](https://github.com/Bitspark/bittype/issues/220), closed; the tag prefixes under [#218](https://github.com/Bitspark/bittype/issues/218); the two further instances under [#251](https://github.com/Bitspark/bittype/issues/251) |

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
remains at an old path. Go's `AnalysisContractFor(settings)` arrived with pull
request #204, because it binds analysis to an import environment that only #204
has.

The Go move changed its place in the order on 12 November 2026. It was first
prepared as draft pull request #231 on top of #204. #231 is closed. The move
went to `main` first, as #246, because
[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
puts the move of a misplaced responsibility before an addition. #204 rebases
onto it and blocks nothing in this plan. The reasons are in a
[comment on #204](https://github.com/Bitspark/bittype/pull/204#issuecomment-6812747024).

Finding 7 in detail. The helper key is an interim form. In Go a consumer writes
`key := g.NewHelperKey("horos/optional")` and passes `key` to
`record.Library(key)` and to `optional.Library(key)`. Neither
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

The change that wrote this record found one more instance while it checked
what remains, and pull request [#250](https://github.com/Bitspark/bittype/pull/250)
corrected it:

- **The tag prefixes of canonical horos values.** `horos/decl/` and `horos/`
  are the ontos tag namespaces of declaration forms and of descriptor
  constructors. Three production packages spelled them:
  `resolve/files/horos/provider/go` and `verification/files/source/horos/go`,
  which each configure the same `Frontend` and must agree, and
  `internal/horosmodel` in `SourceModules`. With #250, `semantics/horos/go`
  published the constants `DeclarationTagPrefix` and `DescriptorTagPrefix`, and
  all three used them. Both constants are in `syntax/horos/go` since 13
  November 2026, as the second paragraph below this list says.
- **The arrow value.** `verification/files/source/horos/go` and
  `verification/files/horos-to-go/go` each built the tuple tagged
  `horos/decl/arrow` by hand. The first now exports `ArrowValue(kind)` and the
  pair calls it. The pair therefore imports no ontos package directly, and its
  two external allowances for `ontos/core/go` and `ontos/data/go` are removed
  from `modelExternalImports`.

It is a move as well: the constants and the arrow value are unchanged. Read in
source on the head of #250, `cadf495`. The owner of the two prefixes is the
horos library, which reads the values. horos 0.10.2, the version that bittype
pins, spells them in `declarations/grammar.go` and exports neither. bittype has
one spelling of its own until horos exports them.

The names of canonical horos data on 13 November 2026. This paragraph first
said that the tag names after the prefix have one writer in bittype,
`internal/horosmodel`. That was wrong for two of them, as
[#256](https://github.com/Bitspark/bittype/issues/256) found. The parser wrote
the field tag `horos/field` and the formatter compared with it. The label
`arrow` was written by `internal/horosmodel` and again by `ArrowValue` in
`verification/files/source/horos/go`. The words `open`, `closed`, `req` and
`opt` of a record payload were written by the parser and read by the
formatter, the source reader and the verification model, and the four record
mappings named `closed` as the discriminator of their pattern. The parser and the
formatter could not take any of these from `semantics/horos/go`, because the
role of a Syntax stage does not permit the semantic vocabulary. Pull request
[#285](https://github.com/Bitspark/bittype/pull/285) therefore moved the two
prefixes to `syntax/horos/go`, which every one of these packages may import
and already imported. It declares with them the labels `arrow` and `record`,
the field tag and the four words, in
[`tags.go`](../syntax/horos/go/tags.go). Every party writes or reads the
horos tree, or builds the canonical value that the tree mirrors. The tree
itself carries the field tag, the four words and the constructor `record`: the
[horos document](horos.md#syntax-vocabulary) says that descriptor operands
"preserve the exact canonical value layout", and that labels and field
wrappers "retain their literal roles". The tree carries neither `horos/decl/`
nor `arrow`. `internal/horosmodel` adds them when it lowers the tree to a
canonical value, and `ArrowValue` builds a value outside any tree.
`semantics/horos/go` keeps neither prefix and no alias. The move added no
role, no import and no allowance, and every value is unchanged. The owner of
all of these names is still the horos library
([Bitspark/horos#53](https://github.com/Bitspark/horos/issues/53) asks it to
export them), and `tags.go` is a stand-in until it does. The other 28 labels
that `internal/horosmodel/read.go` writes after the declaration prefix, such
as `node`, `nu` and `module`, have one writer in bittype.

The correction is partial for three of the nine names. The four record
mappings read `RecordClosed` since the same pull request. Two spellings by
literal remain, and both are the same names as constants of `tags.go`. The
constructor name `record` is the same name as `RecordLabel`: the parser writes
it into the tree, `internal/horosmodel` appends it to the descriptor prefix,
and the formatter, the admission grammar, the source reader, the declaration
recipe, the verification model and the four record mappings compare with it.
The selector values `open` and `closed` in `internal/declsyntax` are the same
names as `RecordOpen` and `RecordClosed`, because the horos library derives
the selector value from the mode atom of the record. Both belong to item 11 of
#256, which covers the constructor names as a whole and closes this.

Other names stay where they are written. These are not instances of finding
12. A name needs a single owner when two parties must agree on it, and here
there is no second party. One of them becomes an instance when a production
package other than its writer reads it by name.

- **horos names that `internal/horosmodel` writes.** The entity kinds
  `type-parameter`, `kind-definition` and `kind-parameter`, the parameter sorts
  `proper` and `constructor`, the operators `kind/<form>`, `sort/arrow`,
  `membership/assumption`, `membership/proven`, `identity/declaration` and
  `identity/parameter`, and the diagnostic codes. Each has a production writer
  and no production reader that names it. Tests read some of them.
- **The encoding keys `json-1` and `utf8` of the Go, TypeScript and Rust
  models.** `internal/golangsem/export.go`,
  `internal/typescriptsem/compiler.cjs` and
  `internal/rustsem/helper/src/main.rs` write them. No production Go code
  compares a term's encoding with either. A search of the production Go files
  for `.Encoding` finds the writers in `internal/horosmodel`, the generic
  validity checks in `semantics/go/model.go` and the horos reader, which
  compares with `TermEncoding()`. A search for the two names finds them as
  keys only in `internal/golangsem/export.go`, where they are written. The
  integrator of the plan checked this on `main` at `abfc948`, and the change
  that wrote this record repeated both searches on the head of #250.

What [#251](https://github.com/Bitspark/bittype/issues/251) examined. This
section first listed two groups of names as not examined: the horos syntax
node kinds, and the names that the native compiler helpers write. The issue
took both on 12 November 2026. Each held a confirmed instance, and both are
corrected. Read in source at `f821eb3`.

- **The horos syntax node kinds, attributes and roles.** `syntax/horos/go` and
  `internal/declsyntax` each declared 33 node kinds, 9 attribute names and 18
  role names. The two lists were equal in every name and value, and neither
  package imported the other. The parser and the formatter used one
  declaration, and the checker and the private model the other. Pull request
  [#254](https://github.com/Bitspark/bittype/pull/254) left one
  declaration, in `syntax/horos/go`. `internal/horosmodel` takes the names
  from there and no longer imports `internal/declsyntax`. `internal/declsyntax`
  defines its 60 constants from the vocabulary's by constant conversion, so no
  name is spelled there. The second declaration came from pull request #39 of
  4 November 2026 (`d59c139`). It added a second language and made the private
  model generic over a language marker `L`. Pull request #49 (`2810050`)
  removed that language the same day, and the second declaration stayed.
- **The TypeScript graph names.** The
  [scan on #251](https://github.com/Bitspark/bittype/issues/251#issuecomment-6814719220)
  searched `main` at `9867398` for graph names spelled outside the vocabulary
  packages, the compiler bridges and the horos private parts. Only
  `format/files/typescript/provider/go` read native graph names by spelling.
  Rust and Go have no such reader. Pull request
  [#255](https://github.com/Bitspark/bittype/pull/255) published
  sixteen names from `semantics/typescript/go`: seven that the Files formatter
  reads, and nine that the vocabulary package's own model readers read. Each
  of the sixteen is spelled once in Go production code, apart from the symbol
  package below.
- **How the helper and the vocabulary stay equal.** The helper
  `internal/typescriptsem/compiler.cjs` writes all sixteen. It is JavaScript,
  cannot import a Go package and is not edited. Nothing makes the two
  spellings one. Tests keep them equal: they run the real helper and read its
  graph through the Go readers. The author of #255 changed each published name
  in the vocabulary in turn. Existing tests fail for twelve. The new test
  `TestTypeScriptHelperNamesMatchVocabulary` fails for three more:
  `BuiltinKind`, `ParameterKind` and `NativeTypeRelation`. No test fails for
  `ImportRole`, because the formatter stores bindings of that role and nothing
  looks them up. Only the vocabulary was changed in these runs. The helper's
  spelling was never changed.
- **The symbol package keeps its own `typescript/name` role.** Decided on
  12 November 2026 by the integrator of the plan.
  `transform/syntax/typescript/symbol/go` declares the role as its private
  `nameRole`, writes it in `ExternalUse` and reads it in `ReadExternal`. Its
  consumers call those two functions and never name the key. No production
  code compares it with the binding role that the helper writes. So the role
  has one party and needs no second owner. The role that the architecture
  test gives the package does not permit `semantics/typescript/go`, and that
  stays.

Limits of these two corrections:

- No independent review read #254 or #255. Their qualification is their
  authors'.
- `ImportRole` is published and no test fails when its value changes.
- The tests were shown to fail for a changed vocabulary. They were not shown
  to fail for a changed helper.
- The language marker `L` of `internal/horosmodel` and `internal/declsyntax`
  is a remainder of the removed second language. Only horos uses the framing.
  The constants of `internal/declsyntax` stay exported, and no package outside
  it uses them.

Still not examined, and so neither an instance nor cleared.
[#256](https://github.com/Bitspark/bittype/issues/256) tracks these. Each
bullet says only what was seen in source at `f821eb3`:

- **What the TypeScript helper reads.** `internal/typescriptsem/compiler.cjs`
  reads the syntax node kinds, child roles and attribute names of its request
  by spelling: a set of supported kinds from line 23, the attributes `operator`
  and `name`, the roles `default` and `type-parameter`, and a switch over node
  kinds from line 95. `syntax/typescript/go` declares these names in Go.
- **The status of the helper's response.** `ok`, `invalid`, `unsupported`,
  `incomplete`, `mismatch` and `environment` are spelled in the helper and in
  `internal/typescriptsem/compiler.go`.
  `resolve/files/typescript/provider/go/resolver.go` spells `ok`, `invalid` and
  `unsupported`, and `analyze/files/typescript/queries/go/queries.go` spells
  `ok`. It is a string of the bridge's output and not a graph name.
- **The generation form names of TypeScript.** `alias`, `interface` and
  `record` are spelled in `transform/syntax/typescript/model/go/model.go`.
  `transform/syntax/typescript/target/go/provider.go` spells `alias` and
  `record`, and the declarations mapping
  `transform/syntax/horos-to-typescript/mapping/declarations/go` spells
  `interface` and `alias`. They are not graph names either.
- **The Go and Rust vocabulary packages' own readers.** The scan did not
  cover spellings inside a vocabulary package. Two were seen:
  `semantics/golang/go/model.go` spells `Key("alias")` and `Key("defined")`
  twice each, and the bridge `internal/golangsem/export.go` writes both kinds.
  `semantics/rust/go/model.go` spells `Key("native-type")` and
  `Key("declared-type")`, which the Rust helper
  `internal/rustsem/helper/src/main.rs` writes. The two sides were not
  compared name by name, and what test fails when they differ was not asked.
- **The admission grammars of the native checkers.** `rules` in
  `check/syntax/<language>/provider/go/check.go` is a table keyed by string
  literals for Go, TypeScript, Rust and C++. It spells node kinds, child roles
  and attribute names that `syntax/<language>/go` also declares. Each of the
  four checkers also compares a node kind with an inline string in a few
  places, for example `"path-type"` and `"lifetime-reference"` in the Rust
  checker. As in `internal/declsyntax` before #254, a checker spells names that
  its vocabulary declares.

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
there. #247 and #248 changed none of these files and added no package. #250
changed one thing in them: it removed two external allowances from one entry
of `modelExternalImports`. #254 changed two entries of `oneOfAKindRules`:
`internal/declsyntax` gains `syntax/horos/go`, and `internal/horosmodel` gains
`syntax/horos/go` and loses `internal/declsyntax`. #255 changed none of these
files. Every count was made again at `f821eb3`, and none changed. The counts
at `71bfec9`, before the Go move, are in
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
| `modelDependencies` and `modelExternalImports` | 8 and 10 | The external model packages that the production graph may hold, and the packages that may import each directly. Since #250 the entry of `verification/files/horos-to-go/go` names the horos declarations and the logos contract, and no ontos package |

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

State on `main` after pull request #250, 12 November 2026.

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
   for three further names that the first change found. The issue is closed.
   One more instance of finding 12, the tag prefixes of canonical horos
   values, was found after that and is done with
   [#250](https://github.com/Bitspark/bittype/pull/250) under #218.
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
   move and does not remove the profile gate. The change for
   [#216](https://github.com/Bitspark/bittype/issues/216) removed the gate on
   12 November 2026
   ([bittype/ACX-008](architecture-cases/0008-structural-profile-gate.md)).
   [#207](https://github.com/Bitspark/bittype/issues/207), which
   writes the first families into mapping packages, is open.
5. **[#218](https://github.com/Bitspark/bittype/issues/218): rules by role.**
   Done with [#244](https://github.com/Bitspark/bittype/pull/244). A package
   gets its import rule from the role that its path gives, and a new package
   of an existing role needs no edit to the test. The step was narrowed as
   [rules to enforce](#rules-to-enforce) records: the foundation and kernel
   set keeps explicit entries, and so do twelve packages of which there is
   only one. Eight hand entries stayed in the test for Go, and
   [#246](https://github.com/Bitspark/bittype/pull/246) converted them.

No step is recorded as not done, and #219, #220, #221 and #222 are closed. The
plan in #218 is complete with the change that wrote this record. Finding 7,
the two candidates and finding 10 are not steps of this order: they stay with
#196, #268 and #161 as the findings table says, and the candidate of finding 9 stays
recorded here and in
[bittype/ACX-009](architecture-cases/0009-package-boundaries.md#candidates).
[#251](https://github.com/Bitspark/bittype/issues/251) examined the names
that this document first listed as not examined, and is closed with the change
that wrote this paragraph. What it left unexamined is listed under
[finding 12](#findings-on-main-at-3bf45c3) and tracked in
[#256](https://github.com/Bitspark/bittype/issues/256). Pull
request #204 is not part of this plan. It adds imported packages to
Go's semantics and, merged after the move, `AnalysisContractFor(settings)` to
`semantics/golang/go`.

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
#247, #248 and #250 seeded no forbidden import. Their authors compared values
before and after: the published semantic graph for #247 and #248, and the two
constants and the arrow value for #250. The hosted `docs`, `contracts` and
`native` checks passed on #247, #248 and #250.
[bittype/ACX-009](architecture-cases/0009-package-boundaries.md#correction-and-qualification-status)
records this per pull request.

No other platform repository imports a bittype package, and bittype has no
release tag (checked on 11 November 2026). The moves cost nothing outside this
repository today. That changes when `libraries/service` starts to consume the
generation packages.

## Independent review

Two reviews read moves from the code. The first read the moves that merged on
11 November 2026. The second read the Go move of #246 and is
[below](#the-go-move-pull-request-246). No independent review read #247, #248,
#250, #254 or #255.

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
one of them behind a build tag that is never set. The pull request body also
says that the reviewer repeated the comparison of the contract values with
its own program and found them equal.

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
