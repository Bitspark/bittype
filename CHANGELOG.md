# Changelog

## Rust graph names published by the vocabulary package — 12 November 2026

The seven Rust graph names that the compiler helper writes and the model readers of `semantics/rust/go` read by name are published by that package. This is the Rust half of item 2 of [#256](https://github.com/Bitspark/bittype/issues/256). It is a move: each key value is unchanged, and no Go code writes any of the seven.

- **Assessment.** The helper `internal/rustsem/helper/src/main.rs` is a Rust program and the reader is a Go package, so they are two parties. The bridge `internal/rustsem` decodes the helper's graph as a whole and names none of its keys. `model.go` of the vocabulary package is the only production Go code that reads a Rust graph name: no package under `resolve`, `analyze`, `format`, `verification` or `transform` does.
- **Published.** Fact relations: `DeclaredTypeRelation()` is `rust/declared-type`, `PredicateRelation()` is `rust/predicate`, `DefaultRelation()` is `rust/default`, `NativeTypeRelation()` is `rust/native-type`, `DisplayRelation()` is `rust/display` and `FieldsRelation()` is `rust/fields`. Term operator: `FieldsOperator()` is `rust/fields`.
- **Writers and readers.** The helper writes all seven. It cannot import a Go package and is unchanged. `model.go` reads the seven through these functions. Its private `fact` method takes the relation key and no longer a name.
- **How the helper and the vocabulary are kept equal.** By tests that run the helper and read its graph through the Go readers. `TestRustHelperNamesMatchVocabulary` is new. It reads the helper's graph through each of the seven names and reports one failure per name. Before it, no test opened a native type by id, which is the only read of `NativeTypeRelation()`.

Not done:

- No name was changed in turn to list the failing tests, as #255 did for TypeScript. The Rust helper is not built on the workstation where this change was written, so every test that runs it fails there before and after, and a changed name shows nothing. Which existing test reads which name was established by reading the tests: `TestSemanticProviders/rust` reads the declared type, the fields relation and operator and the display, `TestRustSemanticNativeTypes` the declared type, the fields and the default, `TestRustSemanticNamespaceAndCoordinates` the declared type and the fields, and `TestRustSemanticBoundsAndOutcomes` the predicate. The hosted `contracts` check builds the helper and runs these tests and the new one.
- Names with a writer and no Go reader that names them stay where they are written. From the helper: the entity kinds `module`, `struct`, `enum`, `variant`, `field`, `alias`, `trait`, `impl` and `parameter`, the parameter sort `type`, the operators `primitive`, `tuple`, `slice`, `array`, `pointer`, `reference`, `variants`, `native-type`, `trait-predicate` and `display`, the encodings `json-1` and `utf8`, the relations `variants` and `implements`, the expression role `declared-type`, the binding role `target`, and the diagnostic codes `unsupported`, `incomplete`, `mismatch`, `units`, `syntax` and `driver`. From Go: the diagnostic codes `limit`, `missing-context` and `compiler-error` of `internal/rustsem`. The helper also spells the subject of its diagnostics, `captured-type-coverage`, which the vocabulary publishes as `Coverage()`. No production code compares a diagnostic's subject or code with a name.
- Two neighbouring groups were seen and not examined. The status of the helper's response is spelled in the helper, in `internal/rustsem/compiler.go` and in `resolve/files/rust/provider/go/resolver.go`. The helper reads the syntax node kinds, child roles and attribute names of its request by spelling, and `syntax/rust/go` declares them in Go.
- The Go half of item 2 is its own change, pull request #264. It follows pull request #204, which edited the same files.
- [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) and [package boundaries](docs/package-boundaries.md) still list the Rust names under "Not examined". This change edits no document besides this entry.

## Go imported packages — 12 November 2026

- The first Go profile accepts an explicit import environment (#168). Captured files form packages by their paths under a module, and supplied dependency packages are given as source in receipt identity. Packages are checked in import order through a closed importer, without module or standard-library discovery. Qualified type references bind to imported declarations; supplied declarations are exported by identity only. The base configuration is unchanged.
- Supplied packages are kept out of the graph's structural projection for containment, not for secrecy: properties the compiler establishes for captured types, and human-readable text, may name supplied declarations. A supplied import path with a `vendor` element is refused.
- A Go resolution whose context ends returns the context's error without waiting for `go/types`, in the base profile too. Its whole native work runs beside the call as one unit: conversion, the package checks and the export, whose type comparisons cannot be interrupted either. Reconstructing the compiler state for queries and an admission query run the same way. `go/types` cannot be interrupted, so abandoned work keeps its goroutine until the call it is inside ends; [native semantics](docs/native-semantics.md) states the limit and the measured costs. Before, a 1,618-byte supplied package held a call with a 20 ms deadline for 11 seconds, and captured alias chains of depth 24 held a base-profile resolution for 8.8 seconds.
- An admission query's constructor must be a generic defined type. An instance used as constructor was accepted before, and a non-generic type with no arguments made `go/types` panic; both are now refused as unsupported.
- Analysis compares its whole configuration with the resolution's, and a test now fails when only lengths or keys are compared.

## Generation form names declared once by each language's model package — 12 November 2026

The names of the intended native forms have one declaration per language, in `transform/syntax/<language>/model/go`. This is item 5 of [#256](https://github.com/Bitspark/bittype/issues/256). The issue named TypeScript. Go, Rust and C++ had the same shape and are corrected with it. It is a move: every value is unchanged, and the plan steps, the generated tree and the printed text of the four targets are byte-identical before and after, in fifteen dumps.

- **Who spelled them.** A model package issues its language's forms through `Form(name)` and reads a name back through `Name(form)`. It had no named values, so every user passed or compared a string. Three packages per language did: the model itself, the target `transform/syntax/<language>/target/go`, which read the expression form names and `alias`, and the declarations mapping `transform/syntax/horos-to-<language>/mapping/declarations/go`, which wrote the two declaration form names. Go's optional helper wrote `defined` as well.
- **Published.** Exported constants beside `Form` and `Name`. In all four models `FormPrimitive` is `primitive`, `FormField` is `field`, `FormRecord` is `record`, `FormExternal` is `external` and `FormAlias` is `alias`. TypeScript adds `FormArray`, which is `array`, and `FormInterface`, which is `interface`. Go adds `FormSlice`, which is `slice`, and `FormDefined`, which is `defined`. Rust and C++ add `FormSequence`, which is `sequence`, and `FormStruct`, which is `struct`. The model's constructors and checks, the target's switch and the mappings use them. Go's `external` has one party, because only the model writes and reads it. It is declared with the other six so that the block is the whole list.
- **The two primitive names.** The primitive form carries a name, `string` or `bool`. The issue did not list them. The text and boolean mappings wrote the name by spelling, the model compared it by spelling, and the TypeScript, Rust and C++ targets compared `string` by spelling. Each model declares `PrimitiveString` and `PrimitiveBool`. The Go target compares neither. It writes the name into the tree as the Go type name, as before.
- **No package shared between the languages.** `primitive`, `field`, `record`, `external` and `alias` are spelled alike in the four models. Each model owns its own forms, a form can be read only by its issuer, and no package compares a form name of one language with another's or with a literal. The planner copies a form's name into the steps of the inspectable plan and compares nothing. The common declaration recipe receives the two declaration forms as values.
- **Imports.** None is new. The target and the mapping packages already imported the model of their language, and no architecture test is edited.
- **What fails when a name changes.** The 36 published values were changed in their owners, and `integration/go`, the one test binary that links the model packages, was run for each change. Six were changed one at a time and make tests fail: `FormAlias` and `PrimitiveString` of TypeScript, and `FormAlias`, `FormDefined`, `PrimitiveString` and `PrimitiveBool` of Go. A test composition spells five of them, and the two Go primitive names reach the generated Go code, which changes with them. The other thirty were changed together in one run, each to a value of its own, and no test failed. A further changed name can add a failure and cannot remove one, so none of the thirty makes a test fail alone. Every production party now takes these names from the owner, so a changed value cannot make two of them disagree, and the tests reach these forms through the model's constructors and the mapping libraries.

Not done:

- The test compositions `integration/go/testdata/generationfixture` and `integration/go/generation_planning_test.go` still pass `alias`, `defined` and `string` by spelling. No test file is edited.
- The C++ record mapping writes the external type `std`, `optional`, and the C++ target compares that pair to add the header. This is the place where Record decides how an optional field is represented. It belongs to [#196](https://github.com/Bitspark/bittype/issues/196), item 6, and is not changed here.
- Seen and not assessed. The four targets spell the syntax attribute names and values (`name`, `optional`, `package`, `visibility`, `absolute`, `global`, `record-key`, `path`, `system`) by literal beside the `Attribute*` constants of `syntax/<language>/go`. The horos descriptor constructors `record`, `scalar` and `list-of` and the discriminators `closed`, `utf8-text` and `bool` are spelled in the common mapping parts, in the horos source reader, in the mapping libraries and in `verification/files/source/horos/go`.
- [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) and [package boundaries](docs/package-boundaries.md) still list the form names under "Not examined". This change edits no document besides this entry.

## Documents record what #251 examined and what stays unexamined — 12 November 2026

[Package boundaries](docs/package-boundaries.md) and [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) listed two groups of names as not examined. [#251](https://github.com/Bitspark/bittype/issues/251) examined both, and the two entries below record the corrections. This entry changes documents only.

- **The scan.** A search of `main` at `9867398` for graph names spelled outside the vocabulary packages, the compiler bridges and the horos private parts found one outside reader, `format/files/typescript/provider/go`. Rust and Go have none. The scan did not cover spellings inside the Go and Rust vocabulary packages.
- **The case.** ACX-009 gains the two corrections of #254 and #255 with their authors' qualification and their limits. Its status stays "Corrected, with limits". The two instances were confirmed and corrected on the same day.
- **The symbol package.** `transform/syntax/typescript/symbol/go` keeps its own `typescript/name` role. It writes and reads the role itself and its consumers never name the key, so the role has one party. Decided on 12 November 2026 by the integrator of the plan.
- **Limits recorded.** No independent review read #254 or #255. `ImportRole()` is published and no test fails when it changes. The tests were shown to fail for a changed vocabulary and not for a changed helper. The language marker `L` of the horos framing is a remainder of a removed language.
- **Other documents.** `AGENTS.md`, `LAYOUT.md`, [native semantics](docs/native-semantics.md) and the two index tables follow. The counts of the architecture test were made again at `f821eb3`, and none changed.

Not examined, and tracked in [#256](https://github.com/Bitspark/bittype/issues/256):

- What the TypeScript helper reads by spelling: the syntax node kinds, child roles and attribute names of its request.
- The status of the helper's response, which the helper, the bridge, the resolver and the query service each spell.
- The TypeScript generation form names `alias`, `interface` and `record`, which the model, the target and the declarations mapping spell.
- The Go and Rust vocabulary packages' own readers against the names that their bridge and helper write.
- The admission grammars of the four native checkers, which are keyed by string literals beside the constants of `syntax/<language>/go`.

## TypeScript graph names published by the vocabulary package — 12 November 2026

The sixteen TypeScript graph names that the compiler helper writes and a Go package reads by name are published by `semantics/typescript/go`. This is the Go side of observation 1 of [#251](https://github.com/Bitspark/bittype/issues/251) for TypeScript. It is a move: each key value is unchanged, and the semantic graph, the groups, the relations and the printed text in seven dumps of three fixtures are byte-identical before and after.

- **Published.** Entity kinds: `ModuleKind()` is `typescript/module`, `PackageKind()` is `typescript/package`, `PackageTypeKind()` is `typescript/package-type`, `BuiltinKind()` is `typescript/builtin`, `AliasKind()` is `typescript/alias`, `InterfaceKind()` is `typescript/interface` and `ParameterKind()` is `typescript/parameter`. Fact relations: `NamedExportRelation()` is `typescript/named-export`, `DeclaredTypeRelation()` is `typescript/declared-type`, `ConstraintRelation()` is `typescript/constraint`, `NativeTypeRelation()` is `typescript/native-type`, `DisplayRelation()` is `typescript/display` and `PropertiesRelation()` is `typescript/properties`. Term operator: `ObjectOperator()` is `typescript/object`. Binding roles: `NameRole()` is `typescript/name` and `ImportRole()` is `typescript/import`.
- **Writers and readers.** `internal/typescriptsem/compiler.cjs` writes all sixteen. It is JavaScript, cannot import a Go package and is unchanged. `format/files/typescript/provider/go` reads the module, package, package-type and builtin kinds, the named-export relation and the two roles through these functions, and writes the kind of its Files group through `ModuleKind()`. `model.go` of the vocabulary package reads the alias, interface and parameter kinds, the other five relations and the object operator through them. Its two private `fact` methods take the relation key and no longer a name.
- **How the helper and the vocabulary are kept equal.** By tests that run the helper and read its graph through the Go readers. Each published name was changed in the vocabulary in turn, with the helper untouched, and both test binaries that link the vocabulary were run. Twelve names made tests that already existed fail. `BuiltinKind()`, `ParameterKind()` and `NativeTypeRelation()` made none fail, so `TestTypeScriptHelperNamesMatchVocabulary` is new: it formats a module that refers to `Array` and opens a parameter and a native type by id, and it fails for each of the three. `ImportRole()` makes no test fail. The formatter stores the binding of an import binding and looks up only the bindings of type references, so this read has no effect that a test could observe.

Not done:

- `transform/syntax/typescript/symbol/go` still spells `typescript/name` as its own `nameRole`, which it writes in `ExternalUse` and reads in `ReadExternal`. The role that the architecture test gives that package does not permit `semantics/typescript/go`: with the import, `TestProductionArchitecture` fails. This change edits no architecture test, so that second spelling of `typescript/name` remains.
- Names with a writer and no separate Go reader stay where they are written. From the helper: the entity kinds `program`, `script` and `property`, the parameter sort `type`, the operators `primitive`, `literal`, `tuple`, `array`, `union`, `intersection`, `boolean` and `display`, the encodings `json-1` and `utf8`, the relation `structure`, the expression roles `type-expression` and `native-type`, and the diagnostic codes. From Go: the diagnostic codes `limit` and `missing-context` of `internal/typescriptsem`, and the Files relation `type-import` and the setting `declaration-placement-1` of the Files formatter.
- Two neighbouring groups were seen and not examined. The helper reads the syntax node kinds, child roles and attribute names of its request by spelling, and `syntax/typescript/go` owns them in Go. The status of the helper's response (`ok`, `invalid`, `unsupported` and the rest) is spelled in the helper, in `internal/typescriptsem`, in the resolver and in the query service.
- The Rust and Go helpers of #251 are not part of this change. Observation 2 of #251 has its own entry below.
- [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) and [package boundaries](docs/package-boundaries.md) still list the TypeScript names under "Not examined". This change edits no document besides this entry.

## horos syntax node kinds, attributes and roles declared once, in the vocabulary package — 12 November 2026

This answers observation 2 of [#251](https://github.com/Bitspark/bittype/issues/251). `syntax/horos/go` and `internal/declsyntax` each declared the 33 node kinds, the 9 attribute names and the 18 role names of the horos syntax, and neither imported the other. The two lists were equal in every name and every value. `syntax/horos/go` owns the names, and no other production package spells one. It is a move: every value, the parsed and formatted fixtures, the syntax projected from canonical values and the published semantic graph are byte-identical before and after.

- **Why there were two.** Pull request #39 of 4 November 2026 made the private horos model generic over the language marker, so that horos and a second language, horos/wire, could share it. It added `internal/declsyntax` with the names as untyped strings, which fit a tree of either language. Pull request #49 of the same day removed the second language. Since then every production use of the framing is over `syntax/horos/go`'s `Language`. One unit test of `internal/horosmodel` builds a tree over `struct{}`.
- **`internal/horosmodel`** takes the names from `syntax/horos/go` and no longer imports `internal/declsyntax`. Its code is generic over the language marker, so it reads a tree's kind as a horos kind, `s.Kind(n.Kind)`, where it compares one, and the projector's `add` takes a horos kind. `resolver.go` spelled the attribute name `pin` inline and now uses the constant.
- **`internal/declsyntax`** keeps its constant names, because its admission grammar is a table keyed by strings. Each constant is defined from the vocabulary: a kind as `string(h.KindProgram)`, an attribute or a role as `h.AttributeName`. The grammar table `rules` is unchanged. `Check` spelled the root kind inline as `"program"` and now uses the constant.
- **The architecture test.** Two entries of `oneOfAKindRules` change. `internal/declsyntax` gains `syntax/horos/go`. `internal/horosmodel` gains `syntax/horos/go` and loses `internal/declsyntax`. Both packages belong to horos by `languageOwners`, and the import goes from a private part to the vocabulary of its own language. `syntax/horos/go` imports nothing new.

Not done:

- The framing stays generic over the language marker, although only horos uses it. Removing the parameter would remove the conversions. It changes the signatures of `internal/horosmodel` and `internal/declsyntax` and is a change of its own.
- The constants of `internal/declsyntax` stay exported. No package outside it uses them now.
- This change edits no document besides this entry. [Package boundaries](docs/package-boundaries.md) and [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) still list the node kinds as not examined.
- Observation 1 of #251, the names that the native compiler helpers write, is open.
- The checkers of the native languages were not examined. The Rust checker compares node kinds with inline strings beside the constants of its vocabulary.

## Go's contract identities and query service move, and the documents follow — 12 November 2026

Pull request [#246](https://github.com/Bitspark/bittype/pull/246) made the Go part of [#220](https://github.com/Bitspark/bittype/issues/220) and [#221](https://github.com/Bitspark/bittype/issues/221). It is a move. This entry also brings the documents up to date with [#247](https://github.com/Bitspark/bittype/pull/247) and [#248](https://github.com/Bitspark/bittype/pull/248), which have their own entries below, and it records [#250](https://github.com/Bitspark/bittype/pull/250), which has none. [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) records the corrections, their qualification and their limits.

- **Contract identities.** `semantics/golang/go` exports `ResolutionContract()` and `AnalysisContract()` in `contracts.go`. `resolve/files/golang/provider/go` and `analyze/files/golang/provider/go` no longer export `Contract()`. The Go analysis provider no longer imports the Go resolver.
- **Query service.** `Queries`, `NewQueries` and `Admission` are the package `queries` at `analyze/files/golang/queries/go`. `verification/files/horos-to-go/go` imports that package and no longer the analysis provider.
- **No alias.** Nothing remains at an old path. A consumer that called `Contract()` on a Go provider package, or `NewQueries` on the Go analysis provider package, names the new places.
- **The exception rule covers a query service.** An import of a query service package (`analyze/files/<language>/queries/go`) by a package that is not a composition needs an entry in the exception table that names its symbols, as an import of a stage provider package does. The second commit of #246 added this after the independent review. The verification pair has an accepted entry for `Admission`.
- **The architecture test.** The two open entries of the exception table are deleted. Fourteen remain, and all are accepted: the thirteen same-language `Check` and `Parse` calls and the pair's `Admission`. The eight hand entries in `TestProductionArchitecture` are converted: `semantics/golang/go`, `internal/golangsem`, the two Go provider packages and `analyze/files/rust/provider/go` take their role, and `internal/typescriptimports`, `verification/files/horos-to-go/go` and `verification/logos/go` are entries of `oneOfAKindRules`. 84 of the 127 production packages are classified by their role. 43 keep an explicit entry: 31 in `foundationRules` and 12 in `oneOfAKindRules`. The synthetic tests are 37.
- **Tag prefixes of canonical horos values.** Pull request #250 is one more move of the kind of #247 and #248. `semantics/horos/go` publishes the constants `DeclarationTagPrefix` (`horos/decl/`) and `DescriptorTagPrefix` (`horos/`). `resolve/files/horos/provider/go`, `verification/files/source/horos/go` and `internal/horosmodel` spelled both and now use the constants. `verification/files/source/horos/go` exports `ArrowValue(kind)`, and `verification/files/horos-to-go/go` calls it where it built the tuple tagged `horos/decl/arrow` by hand. The pair imports no ontos package directly, and its two external allowances for `ontos/core/go` and `ontos/data/go` are removed from the architecture test. The horos library owns the two prefixes. horos 0.10.2 spells them in `declarations/grammar.go` and exports neither, so bittype has one spelling of its own until it does. Its author reports that build and vet pass, that the tests with the pinned TypeScript fail only where the Rust helper is missing, and that a program that is not committed found the constants and the arrow value identical to the old literals and the hand-built value, also in encoded form.
- **Order.** The move was first prepared as draft pull request #231 on top of open pull request #204. #231 is closed. The move went to `main` first, because a misplaced responsibility moves before an addition ([AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)), and #204 rebases onto it. The reasons are in a [comment on #204](https://github.com/Bitspark/bittype/pull/204#issuecomment-6812747024).
- **Qualification, by the author of #246.** `go build ./...` and `go vet ./...` pass. `go test ./... -count=1` has the same outcome per package and per test on the branch and on `main` at `3c3b748`, apart from one new line for the new package, which has no test files. The tests that need the pinned TypeScript and the Rust helper fail on the author's workstation in both runs. A program that is not committed rendered the two contract values and the analysis contract's resolution requirement before and after. With the printed pointer address replaced by a fixed word, each output is 2522 bytes with the same SHA-256. Six seeded forbidden imports each failed `TestProductionArchitecture`. On the second commit, the pair calling `NewQueries` fails and removing the pair's exception fails.
- **Hosted checks.** `docs`, `contracts` and `native` passed on both commits of #246, `7ecdaa1` and `052a3e3`, and on #250 at `cadf495`.
- **Independent review.** A reviewer that wrote none of #246 read its first commit, `7ecdaa1`, on 12 November 2026. Its verdict is that the production code is a move and that no test assertion changed. It found one gap, which the second commit corrected: the verification pair could use any symbol of the query service. It found two loosenings, which were accepted and are recorded as limits of the check: three Go packages may import what their role gives the same package of TypeScript and Rust, and an unused allowance of the five converted packages is checked per role. The reviewer did not read the second commit. Its evidence is in a scratch directory on the workstation and not in the repository.
- **Documents.** `AGENTS.md`, `LAYOUT.md`, [package boundaries](docs/package-boundaries.md), [native semantics](docs/native-semantics.md), the [type-system architecture](docs/type-system-architecture.md), [analysis](docs/analysis.md) and the two index tables follow #246, #247, #248 and #250. Finding 12 of package boundaries and the inline horos keys of ACX-009 are recorded as corrected. ACX-009 has no confirmed instance left. Its status is corrected, with limits. Issues #219 to #222 are closed, and the plan in [#218](https://github.com/Bitspark/bittype/issues/218) is complete with this change. The generated provider tables change through their workflow.

Not done:

- Go has no `AnalysisContractFor(settings)`. It arrives with open pull request #204, which is not part of the plan in #218.
- No independent review read the second commit of #246, #247, #248 or #250.
- Two groups of names were not examined and are listed in the limits of ACX-009: the names that the native compiler helpers write and Go packages read, and the horos syntax node kinds that `syntax/horos/go` and `internal/declsyntax` each declare. [#251](https://github.com/Bitspark/bittype/issues/251) tracks them. Two of the three native writers are not Go, so one owner there needs a design.
- The profile gate and the profile names ([#216](https://github.com/Bitspark/bittype/issues/216), [bittype/ACX-008](docs/architecture-cases/0008-structural-profile-gate.md)) and what a mapping's `Expand` returns ([#196](https://github.com/Bitspark/bittype/issues/196), item 6) are unchanged.
- Three candidates are not confirmed: the thirteen same-language `Check` and `Parse` calls ([#218](https://github.com/Bitspark/bittype/issues/218)), the `Self` rule (a question for #216) and `CheckSingleScope` (with [#161](https://github.com/Bitspark/bittype/issues/161)).

## horos descriptor operators, term encoding and profile names owned by the vocabulary package — 12 November 2026

Three more names that a writer and a separate reader have to agree on have one owner, `semantics/horos/go`. They were found with the five keys of [#220](https://github.com/Bitspark/bittype/issues/220), under the plan in [#218](https://github.com/Bitspark/bittype/issues/218). It is a move: each key value, each string, the published semantic graph and the generated output are byte-identical before and after.

- **Published.** `DescriptorOperator(constructor)` is `horos/descriptor/<constructor>`. `DescriptorConstructor(operator)` is its inverse: it returns the constructor of such a key and reports false for every other key. `TermEncoding()` is `ontos/codec-v1`, the encoding that the data of horos terms carries. ontos owns the codec; that horos terms carry it is a horos agreement. `ResolutionProfile` and `AnalysisProfile` were published before.
- **Writers and readers.** `internal/horosmodel` writes the descriptor operators and the encoding of every provider term through these functions. `transform/syntax/horos/source/go` reads the constructor through `DescriptorConstructor` and compares the encoding of a payload and of a module identity with `TermEncoding()`. It refuses the same terms as before with `semantic.ErrInvalidRequest`, and an operator that is the prefix alone still reads as an empty constructor. `internal/horosfiles` compares the profiles of its stages with the two constants. No production package outside `semantics/horos/go` spells any of the three.

Not done:

- No document describes these three names as inline, and this change edits no document besides this entry.
- ontos 0.16.0 exports no identity for its codec from a Go package, so `ontos/codec-v1` is spelled once, in `semantics/horos/go`.
- A test in `analyze/files/horos/provider/go` still spells the encoding itself. It pins the value independently of the vocabulary.
- The remaining entity kinds, parameter sorts, term operators and diagnostic codes have a production writer and no production reader, and stay where they are written.

## horos keys published by the vocabulary package — 12 November 2026

The five horos-owned keys of [#220](https://github.com/Bitspark/bittype/issues/220) that a writer and a reader each spelled inline are published by `semantics/horos/go`, under the plan in [#218](https://github.com/Bitspark/bittype/issues/218). It is a move: each key value and the published semantic graph are byte-identical before and after.

- **Published.** `ModuleKind()` is `horos/module`, `TypeDefinitionKind()` is `horos/type-definition`, `ImportRelation()` is `horos/import`, `ProperSort()` is `horos/sort/proper` and `ModuleIdentity()` is `horos/identity/module`.
- **Writers and readers.** `internal/horosmodel` writes the two entity kinds and the two operators through these functions. `internal/horosfiles` writes the group kind and the import relation and reads the module kind through them. `transform/syntax/horos/source/go` reads the type-definition kind and the two operators through them and no longer declares its own copies. No production package outside `semantics/horos/go` spells any of the five.

Not done:

- [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) and row 12 of [package boundaries](docs/package-boundaries.md) still describe these five keys as inline. This change edits no document besides this entry.
- Other horos-owned names are still spelled where they are used. The `descriptor/` operator prefix has a writer in `internal/horosmodel` and a reader in `transform/syntax/horos/source/go`. The profile names that `semantics/horos/go` already publishes are repeated as literals in `internal/horosfiles`. The remaining entity kinds, parameter sorts, term operators and diagnostic codes have a production writer and no production reader.

## Import rules by role and the independent review — 12 November 2026

The last step of the plan in [#218](https://github.com/Bitspark/bittype/issues/218) merged with [#244](https://github.com/Bitspark/bittype/pull/244), at `71bfec9`. It changes test code only. [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) records the correction, its qualification and what remains.

- **Rules by role.** The architecture test derives a package's import rule from the role that its path gives (`integration/go/architecture_roles_test.go`). 78 of the 126 production packages are classified this way, and a new package of an existing role needs no edit to the test. 48 keep an explicit entry: 31 foundation and kernel packages in `foundationRules`, 9 packages of which there is only one in `oneOfAKindRules`, and 8 hand entries held for Go. The explicit entries for the foundation and kernel set were decided on #218 on 11 November 2026, because their layering among each other is a fact per package.
- **New in the check.** A package that names a language imports no package of another language. Every allowance is checked for use: an entry and an extension per package, a role per shape. The check reads the files that build constraints exclude. A role grants `internal/<language>` and `internal/<language>sem` only where the imported package is that language's helper. The entries of deleted packages are gone.
- **Independent review.** A reviewer that wrote none of the changes read the merged moves `d23542a..e83cf52` and found that they are moves. It found four defects in the check and in comments and none in moved code, and #244 corrected them. The review is recorded in a [comment on #218](https://github.com/Bitspark/bittype/issues/218#issuecomment-6809378108). Its programs are in a scratch directory on the workstation and not in the repository.
- **Documents.** `LAYOUT.md` replaces its per-component dependency table with a role table. [Package boundaries](docs/package-boundaries.md) says how the test realises each of the six rules, what the check cannot establish and what the review covered. `AGENTS.md` says what adding a package needs.

Not done:

- Go's contract identities and Go's query service are still in its provider packages. Their move is draft pull request #231, stacked on open pull request #204. Two open entries of the exception table and the eight hand entries in the test wait for it.
- Some horos-owned keys are still spelled inline ([#220](https://github.com/Bitspark/bittype/issues/220)).
- The profile gate and the profile names ([#216](https://github.com/Bitspark/bittype/issues/216), [bittype/ACX-008](docs/architecture-cases/0008-structural-profile-gate.md)) and what a mapping's `Expand` returns ([#196](https://github.com/Bitspark/bittype/issues/196), item 6) are unchanged.
- Three candidates are not confirmed: the thirteen same-language `Check` and `Parse` calls ([#218](https://github.com/Bitspark/bittype/issues/218)), the `Self` rule (a question for #216) and `CheckSingleScope` (with [#161](https://github.com/Bitspark/bittype/issues/161)).

## Documents follow the divided common mapping package — 11 November 2026

- `LAYOUT.md`, `AGENTS.md`, [generation](docs/generation.md), [kernel composition](docs/kernel-composition.md) and [type-owned mappings](docs/type-generation-plugins.md) list the adapter package and its siblings `pattern`, `declarations` and `naming` ([#238](https://github.com/Bitspark/bittype/pull/238)). The examples in [generation SPI](docs/generation-spi.md) and [generation planner](docs/generation-planner.md) call `patterns.Element`.

## Package boundaries, first corrections — 11 November 2026

Code changes that merged on 11 November 2026 under the plan in [#218](https://github.com/Bitspark/bittype/issues/218). The pull requests that moved code report byte-identical identity values and generated declarations before and after. [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) records the corrections, their qualification and what remains.

- [#224](https://github.com/Bitspark/bittype/pull/224): tests reach each pair's mapping libraries and the verification laws through one seam file each under `integration/go` (`libraries_*_test.go`, `laws_golang_test.go`).
- [#225](https://github.com/Bitspark/bittype/pull/225), closing [#219](https://github.com/Bitspark/bittype/issues/219): the architecture check names the exact symbols of every import of a stage provider package. An unused allowance, an unused exception and a dot import of a module package fail. Four package shapes are classified by their path. 94 unused allowances are removed from the per-package list.
- [#226](https://github.com/Bitspark/bittype/pull/226), [#227](https://github.com/Bitspark/bittype/pull/227), [#230](https://github.com/Bitspark/bittype/pull/230), [#232](https://github.com/Bitspark/bittype/pull/232): each native mapping library is its own package, `transform/syntax/horos-to-<native>/mapping/<name>/go`, and exports `Library(profile)`. `Libraries(profile)` is deleted. For Go the consumer creates the helper key and passes it to `record.Library(profile, key)` and `optional.Library(profile, key)`.
- [#228](https://github.com/Bitspark/bittype/pull/228): each horos-to-Go verification law is its own package, `verification/files/horos-to-go/mappings/{scalar,list,optional,record}/go`.
- [#229](https://github.com/Bitspark/bittype/pull/229), [#233](https://github.com/Bitspark/bittype/pull/233), [#234](https://github.com/Bitspark/bittype/pull/234): the Rust, TypeScript and horos contract identities are published by `semantics/{rust,typescript,horos}/go` as `ResolutionContract()` and `AnalysisContract()`. `semantics/horos/go` is new and owns the horos roles, relations and profiles, the limits setting and the prelude pin. The TypeScript `Compiler` type moved to `semantics/typescript/go`.
- [#233](https://github.com/Bitspark/bittype/pull/233): the TypeScript query service is `analyze/files/typescript/queries/go`.
- [#238](https://github.com/Bitspark/bittype/pull/238): the common horos mapping package is divided. The adapter `Library` with the profile gate stays at `transform/syntax/horos/mapping/go`. `pattern`, `declarations` and `naming` are packages under `transform/syntax/horos/mapping/<name>/go`. `Pattern` is now `pattern.Head` and `Declarations` is now `declarations.Library`. The profile gate is unchanged ([#216](https://github.com/Bitspark/bittype/issues/216)).
- The documents follow these changes, except #238: `LAYOUT.md` records the path form `<component>/<name>/go` for a component with several independently selectable contributions, and the examples in [kernel composition](docs/kernel-composition.md) and the [Page and Cell example](docs/page-cell-example.md) register each library.

Not done:

- Go's contract identities and Go's query service are still in its provider packages. Their move is draft pull request #231, stacked on open pull request #204.
- The per-package list in the architecture test still classifies every package outside the four path shapes.
- Some horos-owned keys are still spelled inline in `internal/horosfiles` and `transform/syntax/horos/source/go` ([#220](https://github.com/Bitspark/bittype/issues/220)).
- What a mapping's `Expand` returns is unchanged ([#196](https://github.com/Bitspark/bittype/issues/196), item 6).

## Package boundaries — 11 November 2026

- Add [package boundaries](docs/package-boundaries.md): the difference between parts that are wrongly grouped in one package and contributions that are wrongly coupled, the question that separates them, what to look for, the findings of a scan at `3bf45c3`, the proposed import rules by role and the order of work ([#218](https://github.com/Bitspark/bittype/issues/218)). Nothing is implemented.
- Add [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md): the mapping packages, the verification laws and two analysis provider packages group separate contributions, contract identities and horos names are published by provider packages, and the import check cannot see either ([#219](https://github.com/Bitspark/bittype/issues/219), [#220](https://github.com/Bitspark/bittype/issues/220), [#221](https://github.com/Bitspark/bittype/issues/221), [#222](https://github.com/Bitspark/bittype/issues/222)).
- Record the owner's request of 11 November in `AGENTS.md`, and note in `LAYOUT.md`, [generation](docs/generation.md) and [representation families](docs/representation-families.md) what the current check does not establish and that the mapping split precedes the families.

## Order of correction and the profile gate — 11 November 2026

- Record the owner's 11 November direction in [ACX-004](docs/architecture-cases/0004-generation-file-syntax.md): each misplaced placement decision moves to its Files formatter first and coverage follows additively. The Go package clause ([#161](https://github.com/Bitspark/bittype/issues/161)) and Rust visibility ([#159](https://github.com/Bitspark/bittype/issues/159)) can start; C++ ([#160](https://github.com/Bitspark/bittype/issues/160)) is deferred until its semantic provider, with the reason. The Rust external path is recorded as a further pending instance. No correction is implemented.
- Add [bittype/ACX-008](docs/architecture-cases/0008-structural-profile-gate.md): the shared horos mapping policy admits exactly one profile name, and that name is taken from the first fixture ([#216](https://github.com/Bitspark/bittype/issues/216)). Pending.
- Note in [representation families](docs/representation-families.md) the two changes that precede the additive work, and that ontos delivered the expanded-size measure in v0.17.0.

## Chain issues made startable — 10 November 2026

- Record two further pending instances in ACX-004: the Go package clause and the C++ namespace are fixed by target configuration during generation, while the placement proposal takes each declaration's native scope from the consumer's placement ([#161](https://github.com/Bitspark/bittype/issues/161), [#160](https://github.com/Bitspark/bittype/issues/160)). No correction is implemented.
- Point the delivery tables of [representation families](docs/representation-families.md) and [binding generation](docs/binding-generation.md) at the per-family issues split out of #165, the libraries/service and Content tracking issues, and the ordered plan in [#206](https://github.com/Bitspark/bittype/issues/206), which replaces #5.

## Representation families — 9 November 2026

- Record research 0015's accepted design ([representation families](docs/representation-families.md)). Each family is a representation contract, instantiated once per use, whose native type and plan entry are projections of one selected instance. The plan is a graph with a shared envelope and family-owned payloads declared in horos, five version axes and explicit generic requirements; it stays unstable until six vertical cases pass. Correct 0014's union rule: the obstruction is information loss, not overlap. Nothing is implemented yet.

## GEN-01 verified — 9 November 2026

- Record the independent verification of the GEN-01 correction: #192 at
  `6132170`, board record 000005. ACX-006 and GEN-01 read verified. The
  provider-integrity and single-site diagnostic limits remain; #193 tracks
  certifying that a provider's judgments rest on the admitted evidence.

## Binding generation ownership — 9 November 2026

- Record the accepted ownership of binding generation for horos-declared contracts
  (research 0014, [binding generation](docs/binding-generation.md)). bittype stays
  the representation authority and will publish a typed selected representation
  plan. Per-target binding backends compose paired native and adapter contributions
  from it. The charter gains the plan interface only; nothing is implemented yet.

## Generation, semantics and documentation — 8 November 2026

Changes since the previously recorded module-encoding entry:

- Add bounded Go, TypeScript and Rust semantic integrations with explicit compiler
  configuration and capture-bound facts; C++ semantics remain deferred.
- Implement public joint verification association, direct Logos checking and
  library-owned Horos/Go projections/laws, including independent reversed Pair
  composition and checked nested theorem reuse.
- Replace per-pair generation assembly and procedural mapping operations with
  independently registered descriptions, kernel-owned planning, one generic
  convention and separate native AST construction for all four targets.
- Carry generation.Output into TypeScript Files; decide imports/exports after
  placement and check external symbols, public roots and conflicting local names.
- Record open source traversal ownership (GEN-01), native enforcement (GEN-02)
  and Rust/C++ Files ownership (ACX-004).
- Correct GEN-01 (ACX-006, research 0013): the kernel admits the evidence graph and
  walks declarations, references, applications and children itself, with
  run-scoped occurrence handles and binder-keyed actuals. Sources supply
  classification, descriptor selectors and judgments (admitted, unsupported,
  rejected, undetermined) plus a reader for their libraries. Every actual is
  judged, including omitted ones. The Open/Admit/CheckChild/Catalog reader is
  removed without aliases; horos and the independent source, which now publishes
  its own graph, drive the same walk. GEN-02 and ACX-004 remain unfinished.
- Correct the GEN-01 delivery after independent verification rejected it:
  `PlanningInput` takes the source language, input and provider, and `Plan`
  admits the evidence and opens the source itself (`AdmitInput` is no longer
  exported). Handles resolve with `In(Captured)` against the exact capture. A
  source selector outside the capture's source language fails as `ErrScope`. A
  term recorded at several syntax sites is reported at the occurrence's
  declaration. A use whose callee is not generated or whose callee selector is
  foreign is reported as that use. `Plan` checks mapping
  descriptions before it admits evidence, as the registry does. ACX-006 and
  GEN-01 stay open until the correction is verified.
- Reconcile active architecture/API guides, Horos graph examples and dependency
  versions; add a documentation map, runnable Page/Cell fixture walkthrough,
  native placement proposal, exact verification contribution reference and
  current-carrier context table. Repair section links and clarify historical
  consultation authority. Move the maintained matrix identity case into Bittype.


## Horos module encoding — 5 November 2026

- Add the general `horos-module` command: one Horos source module to canonical
  Ontos bytes, sharing the existing parser and module lowering with source pins.
- Add a manual generation-only workflow that builds the encoder and lowers a
  selected repository source on a GitHub runner. No semantic resolution or
  prelude-specific logic is added to the command or agnostic kernel.
- Pin Horos Go 0.10.2, whose complete standard kinds module is authored in Horos
  source and loaded from the generated canonical asset. Keep Ontos at 0.16.0.

## Unified Horos declarations — 4 November 2026

- Replace separate Horos and Horos/Wire frontends with one revision-2 Horos AST
  and provider family. Remove old packages/profiles without compatibility aliases.
- Preserve explicit Type/Data parameter bounds in syntax and provider-owned
  semantic terms. Check spelling structurally; resolve/analysis delegate meaning
  to Horos and retain shared immutable evidence.
- Pin verified public Horos 0.9.0 with Ontos 0.16.0, removing the development
  pseudo-version and retaining the no-replacement installed-consumer gate.
- Migrate canonical fixtures and the independently installed pagination consumer;
  qualify one Page, DataPage aliases, data/mixed specialization and open forwarding.
- Update the provider matrix and ownership/model guides. No generators or runtime
  bindings are introduced; the agnostic kernel and native language stages are unchanged.


## Unreleased

### 4 November 2026

- Remove the Horos-to-HorosWire template-lifting transformer, its tests, exclusive
  indexing helpers and architecture allowances. Author optional pagination directly
  and use the same generic with data and Wire arguments. Preserve closed ground
  embedding and specialization; remove all compatibility paths. Update active
  design documents and align the public Horos dependency to 0.8.0.

- Remove `text.Provider` and the engine's combined provider API without aliases.
  Register parser, checker, formatter and printer adapters independently against
  language identity/profile metadata, in any order. Pack/PackCST validate generic
  carrier structure; explicit checkers preserve syntax admission for Check,
  semantic Capture and every transformation endpoint. Add partial-registration,
  ownership, preflight, concurrency, removed-API and installed-consumer coverage;
  update current documentation and the provider matrix.

- Move parsing, formatting and printing contracts to `parse.Provider`,
  `format.Provider` and `print.Provider`; move source input to `parse.Input`.
  `text.Provider` composes these independent contracts with `check.Provider`.
  Update consumers, independent positive/negative fixtures, dependency checks,
  documentation and the six-language provider matrix. Remove the old names
  without compatibility aliases.

### 3 November 2026

- Move AST syntax admission to `check.Provider` in `contracts/check/go`,
  independent of text and CST contracts. `representation.NewLanguage` accepts
  this interface; `text.Provider` composes it with parsing, formatting and
  printing. Update Horos and native providers, independent consumers, dependency
  checks and the provider matrix. Remove `text.Checker` without an alias.

### 2 November 2026

- Standardize the public provider interfaces as `language.Provider`,
  `resolution.Provider`, `analysis.Provider` and `text.Provider`. Individual
  text capabilities are `text.Parser`, `text.Formatter`, `text.Checker` and
  `text.Printer`. Update all current consumers, fixtures and documentation;
  remove the previous interface names without compatibility aliases.

### 31 October 2026

- Adopt the stage-first hierarchy: trees/{ast,cst,representation},
  contracts/{language,pipeline,text,transform} and uniform native stage paths.
  Move generic AST checks to check/ast. Remove superseded public paths and the
  engine's moved CST exports rather than retaining compatibility aliases.
- Consolidate checked AST/CST snapshots and descriptors in representation;
  preserve sealed payloads and construction-time marker checks across the split.
- Discover the complete native stage matrix in tests, enforce per-package
  production dependencies, and share provider conformance/round-trip/fuzz cases.

- Implement ASTTransformer, typed adapters and owned AST/CST registration
  envelopes with explicit language profiles and marker checks.
- Implement an empty-by-default syntax engine with explicit provider/edge
  registration, unique shortest paths, explicit routes, checked execution and
  existing native text composition. Add independent graph-oracle, ownership,
  negative public API and concurrent-registry tests, including a remote race gate.
- Keep canonical AST and actual language mappings deferred; record the delivered
  frame in [docs/engine.md](docs/engine.md).

- Add the generic CST and borrowed AST-view interface. Native parsers retain
  source tokens/trivia, formatters construct CST directly, and the shared printer
  emits supplied leaves without knowing languages or formatting conventions.
- Implement native TypeScript, Go, Rust and C++ parsing, formatting and syntax
  checks, with independent AST goldens, source round trips, malformed-input
  regressions, fuzzing and remote native compiler qualification.
- Replace the identity interface with LanguageProvider (ID/Name), embedded by
  every specialized text provider; compose providers explicitly. That native-text
  slice added no alias, registry, canonical AST or generator.
- Document the pipeline and ownership in
  [docs/parsing-formatting.md](docs/parsing-formatting.md).

- Add TypeScript, Go, Rust and C++ AST vocabularies over the unchanged generic
  carrier, with independent target consumers and language-separation checks.
- Add an agnostic language identity interface and provider ID/Name methods
  for the owner's intended central registry. That earlier slice added no
  registry, canonical AST or native text stages.
- Document node shapes, native syntax limits and the updated dependency/layout
  boundary in [docs/syntax.md](docs/syntax.md). Generators remain unimplemented;
  native formatting is delivered by the later text slice above.
- Document the Go AST-to-AST generation pipeline, target layout, composition
  strategy and dependency boundaries in [docs/generation.md](docs/generation.md).
- Add generic AST carriers, language-owned node kinds, containment/reference
  links and provenance, plus the independent typed `Stage[I, O]` contract.
- Add remote Go formatting, type-separation and dependency tests, and vet.
  That initial foundation included no language vocabularies or stage executors.
- Add the proposed canonical value-type model and a types-only implementation
  plan for TypeScript, Go and Rust in [docs/model.md](docs/model.md).
- Replace stale bootstrap guidance in the README, charter and contributor
  instructions with the current scope. Function, Stream and Cell stay
  library-owned; the model proposal is not a released language.
- Add remote documentation checks for local links, JSON examples and whitespace.

## Historical bootstrap — 26 September 2026

- Established the repository charter and component-first source layout.
- The original broader scope is retained in Git history. It is superseded for
  current work by the 31 October direction.
- Unmerged research at
  [`91363a5`](https://github.com/Bitspark/bittype/tree/91363a569c62899132ec365103abde7e41c760d9/research-docs)
  remains historical design evidence, not the current model specification.
