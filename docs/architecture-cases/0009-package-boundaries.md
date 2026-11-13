# bittype/ACX-009: Package boundaries do not follow contributions

**Status:** Corrected, with limits; state of 12 November 2026. Every
confirmed instance is corrected on `main`. The last were Go's contract
identities and query service
([#246](https://github.com/Bitspark/bittype/pull/246)) and the horos
names that a writer and a reader each spelled
([#247](https://github.com/Bitspark/bittype/pull/247) and
[#248](https://github.com/Bitspark/bittype/pull/248)), and the tag
prefixes of canonical horos values, which three packages spelled
([#250](https://github.com/Bitspark/bittype/pull/250)). Two further instances
were confirmed after this status was first written, by issue
[#251](https://github.com/Bitspark/bittype/issues/251), and are corrected:
the horos syntax node kinds, attributes and roles, which two packages declared
([#254](https://github.com/Bitspark/bittype/pull/254)), and sixteen
TypeScript graph names that Go read by spelling
([#255](https://github.com/Bitspark/bittype/pull/255)). The status
stands. Two independent reviews read the moves of 11 November 2026 and the
first commit of #246. No independent review read the second commit of #246,
#247, #248, #250, #254 or #255. The [limits](#limits) say what the corrections
and the check do not establish and which names were not examined. Two
candidates are not confirmed and are not part of this case. A third, the
`Self` rule, was answered in principle on 12 November 2026 and is not
corrected ([Candidates](#candidates)).
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
[#196](https://github.com/Bitspark/bittype/issues/196), item 6. Noted on
13 November 2026: whether the second is a fault is conditional since research
0017 (row 11), as the [limit](#limits) on the helper key says.

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

**Corrected, with the limits below. State of 12 November 2026.** No confirmed
instance remains. The state below was read in source on `main` at `1a5aa2d` on
11 November 2026, and the architecture test again at `71bfec9` on 12 November
2026. Only test files and documents changed between those two commits. The Go
rows and the tables of the architecture test were read at `c8d0344`, the merge
of pull request #246, on 12 November 2026. The horos rows were read at
`abfc948`, the merge of pull request #248, on the same day, and the row of the
tag prefixes on the head of pull request #250, `cadf495`. The two rows of
issue #251 were read at `f821eb3`, the merge of pull request #255. The order is
the one in
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
| Go's Record reached the Optional helper through the package variable `optionalKey` | The consumer creates a helper key and passes it to `record.Library(key)` and to `optional.Library(key)`. At the time of #227 both also took a profile, which #265 removed. Neither package creates the key and neither imports the other | [#227](https://github.com/Bitspark/bittype/pull/227) |
| The common horos mapping package held the adapter with the profile gate, the pattern constructors, the declaration recipe and the naming function | The adapter `Library` with the profile gate stays at `transform/syntax/horos/mapping/go`. `pattern`, `declarations` and `naming` are packages under `transform/syntax/horos/mapping/<name>/go`. `Pattern` is now `pattern.Head` and `Declarations` became `declarations.Library`. #238 left the profile gate unchanged. #265 removed the gate, and `declarations.Library` is now `declarations.Recipe`, which states no profile ([bittype/ACX-008](0008-structural-profile-gate.md)) | [#238](https://github.com/Bitspark/bittype/pull/238) |
| The laws package held the `Text`, `Boolean`, `List` and `Record` laws | Each law is its own package, `verification/files/horos-to-go/mappings/{scalar,list,optional,record}/go`. `scalar` exports the rules `Text()` and `Boolean()` of one law. The other three export `Rule` | [#228](https://github.com/Bitspark/bittype/pull/228) |
| The TypeScript analysis provider package held a query service | `Queries`, `NewQueries` and `Assignability` are in `analyze/files/typescript/queries/go` | [#233](https://github.com/Bitspark/bittype/pull/233) |
| Rust, TypeScript and horos identities were published by provider packages. horos had no vocabulary package | `semantics/rust/go`, `semantics/typescript/go` and the new `semantics/horos/go` export `ResolutionContract()` and `AnalysisContract()`. `semantics/horos/go` owns the horos roles, relations and profiles, the limits setting and the prelude pin. The TypeScript `Compiler` type is in `semantics/typescript/go`. Seven of the eight imports in the table above are gone | [#229](https://github.com/Bitspark/bittype/pull/229) rust, [#233](https://github.com/Bitspark/bittype/pull/233) typescript, [#234](https://github.com/Bitspark/bittype/pull/234) horos |
| The check accepted any use of an allowed import and passed unused allowances | Every import of a stage provider package outside a composition needs an entry that names its symbols and its reason. A use outside the named symbols fails. An unused entry, an unused symbol and an unused allowance fail. A dot import of a module package fails. 94 unused allowances were removed | [#225](https://github.com/Bitspark/bittype/pull/225), closing [#219](https://github.com/Bitspark/bittype/issues/219) |
| Every new package edited the test's list | A package gets its import rule from the role that its path gives. At `71bfec9`, 78 of the 126 production packages are classified this way, and a new package of an existing role needs no edit to the test. A package that names a language imports no package of another language. 48 packages keep an explicit entry: the 31 of the closed foundation and kernel set, 9 packages of which there is only one, and 8 entries held for Go. The explicit entries for the foundation and kernel set are a narrowing of the plan, decided on #218 on 11 November 2026: their layering among each other is a fact per package, which one rule for the role would lose | [#225](https://github.com/Bitspark/bittype/pull/225) for four package shapes, [#244](https://github.com/Bitspark/bittype/pull/244) for the others |
| Go's analysis provider package held a query service, and the horos-to-go verification pair imported that package for `Admission` | `Queries`, `NewQueries` and `Admission` are the package `queries` at `analyze/files/golang/queries/go`. `verification/files/horos-to-go/go` imports that package and no longer the analysis provider. It does so under an accepted exception that names `Admission`. No alias remains | [#246](https://github.com/Bitspark/bittype/pull/246) |
| Go's contract identities were published by its resolver and analysis provider packages, and the analysis provider imported the resolver for one | `semantics/golang/go` exports `ResolutionContract()` and `AnalysisContract()` in `contracts.go`. The two provider packages no longer export `Contract()`, and the Go analysis provider no longer imports the Go resolver. The last of the eight imports in the table above is gone. No alias remains | [#246](https://github.com/Bitspark/bittype/pull/246) |
| The exception table held two open entries for Go, and `TestProductionArchitecture` held eight hand entries until Go moved | The two open entries are deleted. Fourteen entries remain, all accepted: the thirteen same-language calls and the verification pair's use of `Admission`. The eight hand entries are converted. `semantics/golang/go`, `internal/golangsem`, `resolve/files/golang/provider/go`, `analyze/files/golang/provider/go` and `analyze/files/rust/provider/go` take their role. `internal/typescriptimports`, `verification/files/horos-to-go/go` and `verification/logos/go` are entries of `oneOfAKindRules`. At `c8d0344`, 84 of the 127 production packages are classified by their role. 43 keep an explicit entry: the 31 of the closed foundation and kernel set and 12 packages of which there is only one | [#246](https://github.com/Bitspark/bittype/pull/246) |
| After the first commit of #246 the verification pair could use any symbol of the Go query service. The exception rule covered stage provider packages only. The independent review of #246 found it | The exception rule covers a query service package as well as a stage provider package. An import of either outside a composition needs an entry that names its symbols. The pair's entry names `Admission` | [#246](https://github.com/Bitspark/bittype/pull/246), second commit `052a3e3` |
| Five horos-owned keys were spelled inline by a writer and by a separate reader: `horos/module`, `horos/import`, `horos/type-definition`, `horos/sort/proper` and `horos/identity/module` | `semantics/horos/go` publishes them as `ModuleKind()`, `ImportRelation()`, `TypeDefinitionKind()`, `ProperSort()` and `ModuleIdentity()`. The writers in `internal/horosmodel` and `internal/horosfiles` and the readers in `internal/horosfiles` and `transform/syntax/horos/source/go` call these functions | [#247](https://github.com/Bitspark/bittype/pull/247) |
| Three further horos names were spelled by a writer and by a separate reader. The change that moved the five keys found them on 12 November 2026: the descriptor operator `horos/descriptor/<constructor>`, the term encoding `ontos/codec-v1`, and the two profile names, which `internal/horosfiles` repeated as string literals | `semantics/horos/go` publishes `DescriptorOperator(constructor)` with its inverse `DescriptorConstructor(operator)`, and `TermEncoding()`. `internal/horosfiles` compares with the constants `ResolutionProfile` and `AnalysisProfile`, which were published before. No production package outside the vocabulary spells any of the three | [#248](https://github.com/Bitspark/bittype/pull/248) |
| The ontos tag prefixes of canonical horos values, `horos/decl/` and `horos/`, were spelled in three production packages. `resolve/files/horos/provider/go` and `verification/files/source/horos/go` each configure the same `Frontend` with them and must agree, and `internal/horosmodel` spelled them again in `SourceModules`. The tuple tagged `horos/decl/arrow` was built by hand in `verification/files/source/horos/go` and again in `verification/files/horos-to-go/go`. The change that wrote this record found both on 12 November 2026 | `semantics/horos/go` publishes the constants `DeclarationTagPrefix` and `DescriptorTagPrefix`, and the three packages use them. `verification/files/source/horos/go` exports `ArrowValue(kind)`, and the pair calls it. The pair imports no ontos package directly, and its two external allowances for `ontos/core/go` and `ontos/data/go` are removed from `modelExternalImports`. The horos resolver now imports its vocabulary, which its role permits. Since 13 November 2026 the two constants are in `syntax/horos/go`: see the next row | [#250](https://github.com/Bitspark/bittype/pull/250) |
| Names of canonical horos data were spelled by more than one production package, or were declared where two of their parties could not import them. The parser wrote the field tag `horos/field` and the record words `open`, `closed`, `req` and `opt`, and the formatter compared with all five. The source reader read `opt` and built the record tag from the prefix and `record`. The verification model read `closed` and `opt`, and the four record mappings named `closed` as the discriminator of their pattern. The label `arrow` was written by `internal/horosmodel` and by `ArrowValue`. The two prefixes were in `semantics/horos/go`, which the role of a Syntax stage does not permit. [#256](https://github.com/Bitspark/bittype/issues/256) confirmed these as its items 6, 21 and 22 on 12 and 13 November 2026 | `syntax/horos/go` declares them in `tags.go`: `DeclarationTagPrefix`, `DescriptorTagPrefix`, the labels `ArrowLabel` and `RecordLabel`, `FieldTag`, and the words `RecordOpen`, `RecordClosed`, `FieldRequired` and `FieldOptional`. The parser, the formatter, `internal/horosmodel`, the horos resolver, the source reader, the verification model and the four record mappings take them from there. Each already imported that package. `semantics/horos/go` keeps neither prefix and no alias, and the horos resolver no longer imports it. No role, import or allowance is added, and every value is unchanged. The declaration is a stand-in for an export of the horos library ([Bitspark/horos#53](https://github.com/Bitspark/horos/issues/53)). **Partial, and pending for three names.** Two spellings by literal remain. The constructor name `record` is the same name as `RecordLabel`: the parser writes it into the tree, and the formatter, the admission grammar, the source reader, the declaration recipe, the verification model and the four record mappings compare with it. The selector values `open` and `closed` in `internal/declsyntax` are the same names as `RecordOpen` and `RecordClosed`. Item 11 of #256, which covers the constructor names as a whole, closes this | [#285](https://github.com/Bitspark/bittype/pull/285) |
| `syntax/horos/go` and `internal/declsyntax` each declared the 33 node kinds, the 9 attribute names and the 18 role names of the horos syntax, with equal names and values. Neither imported the other. The parser and the formatter used one declaration, and the checker and the private model the other. Issue #251 confirmed it on 12 November 2026 | One declaration, in `syntax/horos/go`. `internal/horosmodel` takes the names from there and no longer imports `internal/declsyntax`. `internal/declsyntax` defines its 60 constants from the vocabulary's by constant conversion, so no name is spelled there. Two entries of `oneOfAKindRules` follow: `internal/declsyntax` gains `syntax/horos/go`, and `internal/horosmodel` gains `syntax/horos/go` and loses `internal/declsyntax` | [#254](https://github.com/Bitspark/bittype/pull/254) |
| `format/files/typescript/provider/go` read seven TypeScript graph names by spelling, and the model readers of `semantics/typescript/go` spelled nine more, two of them twice. The compiler helper `internal/typescriptsem/compiler.cjs` writes all sixteen. The scan on #251 found the outside reader on 12 November 2026 | `semantics/typescript/go` publishes the sixteen names as functions: seven entity kinds, six fact relations, one term operator and two binding roles. The Files formatter and the model readers call them. The helper is JavaScript, cannot import a Go package and is not edited. Tests that run the real helper keep the two sides equal, and `TestTypeScriptHelperNamesMatchVocabulary` is new | [#255](https://github.com/Bitspark/bittype/pull/255) |
| The allow-list of a role was not checked for use | Every allowance is checked for use. An entry and an extension are checked per package. A role is checked per shape: each import it permits must be imported by at least one package of that role. External model allowances are checked too | [#244](https://github.com/Bitspark/bittype/pull/244) |
| The check read only the files of the build that ran it, and the vocabulary rule accepted any name as a language. The independent review found both | The check reads every non-test file of a package, also one that a build tag, another operating system or disabled cgo excludes. A role grants `internal/<language>` and `internal/<language>sem` only where the imported package is that language's helper by its own role | [#244](https://github.com/Bitspark/bittype/pull/244) |

[#224](https://github.com/Bitspark/bittype/pull/224) prepared the mapping and
law splits. It put each pair's list of libraries and the laws behind one test
seam file under `integration/go`. It changed no production code.

The Go move changed its place in the order on 12 November 2026. It was first
prepared as draft pull request #231 on top of open pull request #204, which
adds imported packages to Go's semantics. #231 is closed. The move went to
`main` first, as #246, because
[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
puts the move of a misplaced responsibility before an addition. #204 rebases
onto it and blocks nothing in this case. The reasons are in a
[comment on #204](https://github.com/Bitspark/bittype/pull/204#issuecomment-6812747024).
`AnalysisContractFor(settings)` for Go is not part of #246. It arrives with
#204.

### Qualification

The statements below come from the pull requests named, unless a bullet names
another source. Their authors made the before and after comparisons. #227,
#228, #229, #233, #234, #246, #247, #248, #250, #254 and #255 used programs
that are not committed.

- **Hosted checks.** The hosted `docs`, `contracts` and `native` checks passed
  on each of the eleven pull requests, and on #244. They passed on both
  commits of #246, `7ecdaa1` and `052a3e3`, on #247 at `838d5e2`, on #248 at
  `3db9bc9`, on #250 at `cadf495`, on #254 at `da487a9` and on #255 at
  `82b6f49`.
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
- **The check itself.** `integration/go/architecture_synthetic_test.go` holds 37
  tests that feed the check hand-built package sets. Each failing set breaks
  one rule and must produce that rule's message and no other. One test adds a
  third language, a pair into it, one more mapping and one more common mapping
  part without a new line in the foundation table, the hand entries or the
  owners. One test writes a module at test time and runs the real loader on
  files that a build tag, another operating system and cgo exclude. One test,
  added with #246, imports a query service from a package that is not a
  composition: without an exception, with one, and with a use of a symbol that
  the exception does not name.
- **Go's Record and Optional.** `TestGolangRecordRequiresItsOptionalHelper`
  covers an invalid key, a key that no library supplies and the working case.
- **Go's identities and query service, by their author.** #246 reports the
  following, run on its head and on `main` at `3c3b748`. `go build ./...` and
  `go vet ./...` pass. `go test ./... -count=1` has the same outcome per
  package and per test before and after. The one difference is a new line for
  the new package, which has no test files. A program rendered
  `resolve.Contract()`, `analyze.Contract()` and the analysis contract's
  resolution requirement on `main`, and `ResolutionContract()`,
  `AnalysisContract()` and the same requirement on the branch. The rendering
  prints a pointer address. With that address replaced by a fixed word, both
  outputs are 2522 bytes with the same SHA-256. Six forbidden imports were
  seeded, each as a file that was removed again, and each failed
  `TestProductionArchitecture`: the analysis provider importing its resolver,
  the verification pair importing the analysis provider, the query service
  importing its resolver, the query service importing the analysis provider,
  the resolver importing the query service, and the Go vocabulary importing
  the Rust vocabulary. For the second commit, `052a3e3`, the author reports
  two more seeds: the verification pair calling `NewQueries` fails, and
  removing the pair's exception fails.
- **Independent review of the Go move.** A reviewer that wrote none of #246
  read its first commit, `7ecdaa1`, against `main` on 12 November 2026. Its
  verdict is that the production code is a move and that the test edits
  change no assertion. It found that the conversion of the architecture test
  is not a pure move, in three ways. The verification pair could use any
  symbol of the Go query service, which the second commit corrected. Three Go
  packages may import more than their hand entry listed. An unused allowance
  of the five converted packages is checked per role and no longer per
  package. The last two were accepted on #246 and are under
  [limits](#limits). The review seeded the real tree: the pair calling
  `NewQueries` passed on `7ecdaa1` and failed on `main`, the four newly
  permitted imports that no package makes passed on `7ecdaa1` and failed on
  `main`, and nine forbidden imports failed on `7ecdaa1`, one of them behind a
  build tag that is never set. The review is recorded in the body of
  [#246](https://github.com/Bitspark/bittype/pull/246), and
  [package boundaries](../package-boundaries.md#the-go-move-pull-request-246)
  gives each finding and its decision.
- **The five horos keys, by their author.** #247 reports the following.
  A test that is not committed compared each published key with the literal
  it replaced, and all five are equal. Another parsed the two-file pagination
  fixture through the horos Files parser and wrote the whole semantic graph,
  the groups and the relations. Before and after, the output is 431382 bytes
  with the same SHA-256. `go test ./... -count=1` has the same failing
  packages and the same failing tests before and after, with and without the
  pinned TypeScript. `TestProductionArchitecture` passes, and no file of the
  architecture test is edited. The integrator of the plan read the diff.
- **The three further horos names, by their author.** #248 reports the
  following. A test that is kept in a scratch directory and not committed
  wrote the semantic graph of three fixtures as JSON before and after: 263878,
  239082 and 216059 bytes, each with the same SHA-256. It wrote the formatted
  TypeScript, Go, Rust and C++ output of two fixtures, and all eight files are
  identical. It forged published graphs to show that the reader refuses the
  same terms as before: a descriptor operator with another owner or another
  prefix, and a payload or a module identity with another encoding, each end
  in `semantic.ErrInvalidRequest` before and after. Four seeded mutants of the
  changed code each made that test fail and were reverted. The per-package and
  per-test outcomes of `go test ./... -count=1` are identical before and
  after.
- **The tag prefixes, by their author.** #250 reports the following. Its
  author is the integrator of the plan. `go build ./...` and `go vet ./...`
  pass. `go test ./... -count=1` with the pinned TypeScript fails only the
  tests that need the Rust helper, which is the failure list of `main` on that
  workstation. A program that is not committed compared both constants with
  the old literals, and `ArrowValue(v)` with the old hand-built value, also in
  encoded form. They are identical. Neither literal is left in production Go
  code outside `semantics/horos/go`. The architecture tests pass, and the only
  edit to them is the removal of the two unused allowances. The change that
  wrote this record repeated the search for the two literals on `cadf495` and
  found them in `semantics/horos/go/vocabulary.go` only.
- **The horos syntax names, by their author.** #254 reports the following,
  run on `main` at `9867398` and on its head `da487a9`. A program that is not
  committed parsed both declarations before the change and found 33 kinds, 9
  attributes and 18 roles in both, with no name in one only and no different
  value. Another printed every constant of both packages with its type and
  value, 121 lines, identical before and after. A test that is not committed
  dumped 52 files of behaviour: four horos sources through parse, check,
  print, format, the private reader, resolution and analysis, two canonical
  fixtures through the private projection, and four generated outputs. All 52
  have the same size and SHA-256 before and after. The sorted outcomes of
  `go test ./... -count=1` per package and per test are identical, with the
  pinned TypeScript and without the Rust helper. One value was seeded,
  `KindAtom` in `syntax/horos/go`, in plain copies of both trees. Before the
  change the two declarations then disagree and 26 tests fail. After it the
  framing follows the vocabulary and the test outcomes equal the unseeded
  ones. The integrator of the plan read the diff.
- **The TypeScript graph names, by their author.** #255 reports the following.
  A test that is not committed compared each of the sixteen functions with
  the key it replaced, and all are equal. Another wrote the semantic graph,
  the Files groups and relations, the printed text and what the model readers
  return for three fixtures in seven dumps. Each dump has the same size and
  SHA-256 before and after. Each published name was then changed in the
  vocabulary in turn, with the helper untouched, and the two test binaries
  that link the vocabulary were run in full. Tests that already existed fail
  for twelve names. `TestTypeScriptHelperNamesMatchVocabulary`, which is
  committed, fails for `BuiltinKind`, `ParameterKind` and
  `NativeTypeRelation`. No test fails for `ImportRole`. The mutation runs were
  made before `main` was merged into the branch and were not repeated on the
  merged head.
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

- The comparisons of #227, #228, #229, #233, #234, #246, #247, #248 and #250
  cannot be repeated from the repository, because the programs that made them
  were not committed. The tests that compare fragments of the generated
  declarations remain. #248 commits no test for `DescriptorOperator` and
  `DescriptorConstructor`.
- The independent review covered the moves in `d23542a..e83cf52`. It did not
  verify anything that needs the Rust helper, the native compiler job, or the
  race and fuzz steps. The hosted checks ran those on each pull request. Its
  range ends before #244, which changes test code only.
- That range also ends before #246, which moves production code. A second
  independent review read the first commit of #246, `7ecdaa1`. It did not read
  the second commit, `052a3e3`, which extends the exception rule to query
  services. That commit is qualified by its author's two seeds, by the
  synthetic test and by the hosted checks.
- The second review's programs and logs are in
  `C:/Development/bitspark/scratch-bittype-review-246/` on the workstation and
  not in the repository. The review is recorded in the body of #246, which the
  author of #246 wrote.
- No independent review read #247, #248, #250, #254 or #255. Their
  qualification is their authors', and the integrator of the plan read the
  diffs of #247 and #254 and wrote #250.
- **`ImportRole` is published and untested.** `ImportRole()` is one of the
  sixteen TypeScript names. No test fails when its value changes. The Files
  formatter admits bindings of that role and stores them, and nothing looks
  them up, so the read has no effect that a test could observe. Whether that
  read should exist is not decided.
- **Only the Go side of the TypeScript names was mutated.** The tests were
  shown to fail when a name changes in the vocabulary. The helper's spelling
  was never changed, so the tests were not shown to fail for a changed helper.
  The helper is JavaScript and spells every name itself. Nothing but those
  tests keeps the two sides equal.
- **The language marker of the horos framing is vestigial.** `internal/horosmodel`
  and `internal/declsyntax` are still generic over a language marker `L`. Pull
  request #39 introduced it for a second language, which pull request #49
  removed the same day. Every production use is over the horos `Language`. One
  unit test builds a tree over `struct{}`. The parameter costs the conversions
  that #254 added, and removing it is a change of its own. The constants of
  `internal/declsyntax` stay exported, and no package outside it uses them.
- **The `typescript/name` role of the symbol package is not an instance.**
  Decided on 12 November 2026 by the integrator of the plan.
  `transform/syntax/typescript/symbol/go` declares the role as its private
  `nameRole`, writes it in `ExternalUse` and reads it in `ReadExternal`. Its
  consumers call those two functions and never name the key, and no production
  code compares it with the binding role that the helper writes. It has one
  party. The role that the architecture test gives the package does not permit
  `semantics/typescript/go`, and that stays.
- **The tag prefixes have their owner outside bittype.** The horos library
  reads the canonical values, so it owns `horos/decl/` and `horos/`. horos
  0.10.2, the version that bittype pins, spells them in
  `declarations/grammar.go` and exports neither. bittype's one spelling is
  therefore a stand-in and not the owner's. `semantics/horos/go` held it from
  #250, and `syntax/horos/go` holds it since pull request
  [#285](https://github.com/Bitspark/bittype/pull/285) of 13 November 2026,
  because the parser and the formatter agree on names of the same kind and
  may not import the semantic vocabulary. It can go when horos exports the
  names ([Bitspark/horos#53](https://github.com/Bitspark/horos/issues/53)).
  This limit first said that the tag names after the prefix have one writer in
  bittype. That was wrong for the field tag `horos/field`, which the parser
  wrote and the formatter read, and for the label `arrow`, which two packages
  wrote. Both are declared in `syntax/horos/go` now, with the label `record`
  and the four words of a record payload. The other 28 labels after the
  declaration prefix, such as `node`, `nu` and `module`, have one writer in
  bittype, `internal/horosmodel`, and the horos library reads them. The
  correction of #285 is partial. The four record mappings read `RecordClosed`.
  The constructor name `record`, which the parser writes and which is the same
  name as `RecordLabel`, and the selector values `open` and `closed` in
  `internal/declsyntax`, which are the same names as `RecordOpen` and
  `RecordClosed`, are still spelled by literal. They stay pending under item
  11 of [#256](https://github.com/Bitspark/bittype/issues/256), which covers
  the constructor names as a whole.
- **Converted packages may import what their role gives.** Since #246 three
  Go packages may import more than their hand entry listed.
  `internal/golangsem` gains `contracts/parse/go` and `files/go`. The Go
  resolver gains `trees/representation/go`. The Go vocabulary gains
  `contracts/resolution/go`, which it now imports, and `internal/golang`. Each
  is what the role gives the same package of TypeScript and Rust, and none is
  a package of another language or a contribution. This was accepted on #246.
  It is a limit of the check and not a defect that was fixed.
- **Unused allowances of the converted packages.** For the five converted
  packages an unused allowance no longer fails per package. It is checked per
  role, as for every other language. This was accepted on #246 for the same
  reason and is the limit that the bullet on shapes below states.
- The author of #246 ran the tests on a workstation that lacks the pinned
  TypeScript and the Rust helper. The tests that need them fail there before
  and after the move, so that run shows nothing about them. The authors of
  #247, #248 and #250 had the pinned TypeScript and no Rust helper. The hosted
  checks run those tests. The seeded imports of #246 are not in the
  repository either.
- **Names with a writer and no reader are not instances.** A name needs a
  single owner when two parties must agree on it. Where a name has a
  production writer and no production reader that names it, there is no second
  party. Such a name becomes an instance when a production package other than
  its writer reads it by name. Two groups were checked and are of this kind.
  - The horos names that `internal/horosmodel` still spells where it writes
    them: the entity kinds `type-parameter`, `kind-definition` and
    `kind-parameter`, the parameter sorts `proper` and `constructor`, the
    operators `kind/<form>`, `sort/arrow`, `membership/assumption`,
    `membership/proven`, `identity/declaration` and `identity/parameter`, and
    the diagnostic codes. Tests read some of them. Read in source at
    `abfc948`.
  - The encoding keys `json-1` and `utf8` that the Go, TypeScript and Rust
    models write in `internal/golangsem/export.go`,
    `internal/typescriptsem/compiler.cjs` and
    `internal/rustsem/helper/src/main.rs`. No production Go code compares a
    term's encoding with either. A search of the production Go files for
    `.Encoding` finds the writers in `internal/horosmodel`, the generic
    validity checks in `semantics/go/model.go` and the horos reader, which
    compares with `TermEncoding()`. A search for the two names finds them as
    keys only in `internal/golangsem/export.go`, where they are written. The
    integrator of the plan checked this on `main` at `abfc948`, and the change
    that wrote this record repeated both searches on `cadf495`.
- **What issue #251 examined.** This list first named two groups as not
  examined: the names that the native compiler helpers write, and the horos
  syntax node kinds. [#251](https://github.com/Bitspark/bittype/issues/251)
  took both on 12 November 2026.
  - The horos syntax node kinds, attributes and roles were a confirmed
    instance, corrected with #254.
  - The [scan on #251](https://github.com/Bitspark/bittype/issues/251#issuecomment-6814719220)
    searched `main` at `9867398` for graph names spelled outside the
    vocabulary packages, the compiler bridges and the horos private parts.
    Only `format/files/typescript/provider/go` read native graph names by
    spelling. That is a confirmed instance, corrected with #255. Rust and Go
    have no such reader.
- **Not examined.** The following were not assessed. They are neither
  instances nor cleared.
  [#256](https://github.com/Bitspark/bittype/issues/256) tracks them.
  Each says only what was seen in source at `f821eb3`.
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
  part of this case. Noted on 13 November 2026: research 0017 (row 11) makes
  that move conditional on an answer of the service library, and it follows
  the gate of [the contribution contract](../contribution-contract.md#optional-fields).
- A green run still shows only that the listed edges are absent
  ([AC-07](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-07)).

### Instances that remain

No confirmed instance remains on `main` at `f821eb3`, after pull request #255.
The last six left this list on 12 November 2026:

1. **Go's query service**, corrected with pull request #246.
2. **Go's contract identities**, corrected with pull request #246. With
   instance 1 they were the two open entries of the exception table. No hand
   entry is left in `TestProductionArchitecture`.
3. **Inline horos keys**, corrected with pull request #247 for the five keys
   that [#220](https://github.com/Bitspark/bittype/issues/220) tracked, and
   with pull request #248 for three further names of the same kind: the
   descriptor operator, the term encoding and the two profile names. The
   change that moved the five keys found the three.
4. **The tag prefixes of canonical horos values**, corrected with pull request
   #250. The change that wrote this record found them while it checked what
   remains, and listed them first as not examined. They were then confirmed:
   two packages configure the same frontend and must agree, and a third
   spelled the prefixes again. Pull request #285 moved them from
   `semantics/horos/go` to `syntax/horos/go` on 13 November 2026, with the
   other names of canonical horos data that
   [#256](https://github.com/Bitspark/bittype/issues/256) confirmed. That
   correction is partial: the constructor name `record` and the selector
   values of a record's mode are still spelled by literal and stay pending
   under item 11 of #256.
5. **The horos syntax node kinds, attributes and roles**, corrected with pull
   request #254. Issue #251 confirmed them after the status "Corrected, with
   limits" was written. Two packages declared the same 60 names, and the
   parser and the checker each used a different declaration.
6. **Sixteen TypeScript graph names**, corrected with pull request #255. Issue
   #251 confirmed them the same way. The JavaScript helper writes them, and Go
   packages read them by spelling.

Instances 5 and 6 did not change the status. They were confirmed and corrected
on the same day, and the status is dated for the state after both.

[#219](https://github.com/Bitspark/bittype/issues/219),
[#220](https://github.com/Bitspark/bittype/issues/220),
[#221](https://github.com/Bitspark/bittype/issues/221) and
[#222](https://github.com/Bitspark/bittype/issues/222) are closed. The plan in
[#218](https://github.com/Bitspark/bittype/issues/218) is complete with the
change that wrote this record.

A new confirmed instance reopens this list. The names under "Not examined" in
the [limits](#limits) are where one would be looked for first. Issue #251 is
closed with the change that wrote this paragraph, and
[#256](https://github.com/Bitspark/bittype/issues/256) tracks those names.

Pull request #204 is not part of this case. It adds imported packages to
Go's semantics and, merged after the move, `AnalysisContractFor(settings)` to
`semantics/golang/go`.

The profile gate in the adapter of the common horos mapping package is not an
instance of this case. It belongs to
[#216](https://github.com/Bitspark/bittype/issues/216) and
[bittype/ACX-008](0008-structural-profile-gate.md), and #238 did not change it.
The change for #216 removed it on 12 November 2026.

### Candidates

Two candidates are not confirmed and are not part of this case until they
are. The `Self` rule was the third. It is answered in principle and not
corrected.

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
  profile, so the question was passed to #216. #216 did not decide it.
  [#268](https://github.com/Bitspark/bittype/issues/268) asked who owns the rule.
  **Answered in principle on 12 November 2026**
  ([comment on #268](https://github.com/Bitspark/bittype/issues/268#issuecomment-6837825043)): a reserved word is the target's, because it is
  a fact about one language, and no accepted document promises equal native
  names in every target. With that answer the rule is in the wrong place
  ([AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02)):
  the shared function holds a word that only Rust reserves. Nothing is
  corrected. The function still appends `_` to `Self` for all four targets.
  The answer is one part of the [naming design](../naming.md#reserved-words),
  and the correction is delivered by [#274](https://github.com/Bitspark/bittype/issues/274) with #268: the shared
  function holds no reserved word, and each target that needs one states its
  own. The pending record is the Go naming instance of
  [bittype/ACX-004](0004-generation-file-syntax.md), which names the same
  shared function. #268 stays open until then and records the result here.
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
