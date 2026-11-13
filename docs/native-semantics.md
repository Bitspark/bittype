# Native semantic providers

The Go, TypeScript and Rust providers use the existing [resolution](resolution.md) and
[analysis](analysis.md) contracts. Native providers own their compiler integration
and vocabulary in the [agnostic semantic graph](semantics.md). They do not depend
on Horos, Logos, generation plugins or runtime bindings. The
[joint verifier](verification.md) consumes these established facts and
explicit native queries alongside its separate Horos evidence.

## Composition and ownership

```text
captured, syntax-admitted ASTs
  -> resolve/files/golang/provider/go -> go/ast -> go/types
  -> immutable declarations, bindings, type terms, facts and reports
  -> analyze/files/golang/provider/go (explicit resolution receipt; reuses the graph)
  -> semantics/golang/go (typed read-only view, contract identities)
  -> analyze/files/golang/queries/go (explicit native query session, when a consumer needs admission checking)
```

Resolution and type checking are coupled inside `go/types`. The resolver declares
both responsibilities; analysis consumes its exact compatible receipt without
running the compiler again. Syntax admission still permits unresolved names.
There is no language switch in the engine and no automatic registration.
Use resolution.NewResolver and analysis.NewAnalyzer with their respective Go contracts and providers.
The contracts are `golang.ResolutionContract()` and `golang.AnalysisContract()`, published by the
vocabulary package `semantics/golang/go`. The providers are the `Provider` types of
`resolve/files/golang/provider/go` and `analyze/files/golang/provider/go`.

The package layout is:

| Package | Responsibility |
| --- | --- |
| `semantics/golang/go` | Go-owned vocabulary, contract identities, environment identity, immutable typed graph readers and the names the first Go Files slice admits |
| `internal/golangsem` | Direct AST bridge, bounded compiler invocation and fact export |
| `resolve/files/golang/provider/go` | Configured resolution, coverage and native validity reports |
| `analyze/files/golang/provider/go` | Receipt-consuming analysis |
| `analyze/files/golang/queries/go` | Explicit compiler-backed admission queries over an opened model |
| `integration/go` | Provider, environment, coverage and graph qualification |
| `verification/files/horos-to-go/go` | Joint representation claims using the native provider results |

The former `verification/files/target/golang/go` projection is removed. No
compatibility entry point remains. No type-generation behavior changes.

## First Go profile

<!-- bittype:overview scope=first-go-profile part=features -->
The supported subset remains the captured type-declaration fragment: one supplied
package across at most 16 files, defined types and aliases, generic parameters and
applications, named/predeclared types, slices, ordinary struct fields and tags,
embedded interface type terms, approximations and unions. Under an explicit
import environment, the captured files form several packages by their file paths,
and named and aliased imports of captured or explicitly supplied packages bind
qualified type references to the imported declarations.
<!-- bittype:overview:end -->

<!-- bittype:overview scope=first-go-profile part=exclusions -->
Imports outside an explicit import environment, dot and blank imports, the
standard library, methods, embedded struct fields, pointer/array/map/function
types and executable bodies remain unsupported. Syntax acceptance does not
promise semantic coverage.
<!-- bittype:overview:end -->

The language version is Go 1.25. Without an import environment, the package path
is `captured` and the only external scope is the predeclared universe. The native compiler is the `go/types`
library linked into the process; its actual Go toolchain version is part of the
required request configuration. Sizes are explicitly `gc/amd64`, matching the
previous default. CI qualifies Go 1.27.3. The provider does not discover a project,
read ambient packages, resolve imports or download dependencies. Changing any of
these environment choices requires a deliberate profile change.

Every result records separate coverage, binding and validity reports. A native
type error gives a rejected validity report with a captured-node diagnostic;
unsupported syntax gives unsupported coverage; exhausted bounds give incomplete
reports. None publishes partial semantics. Invalid configuration, forged supplied
references, foreign receipts and cancellation remain operation errors.

A call whose context ends returns that context's error without waiting for
`go/types`. The whole native work of a resolution runs beside the call as one
unit: conversion, every package check and the export, whose type comparisons
`go/types` cannot interrupt any more than a check. Reconstructing the compiler
state for queries and an admission query run the same way, an admission
including its wait for an earlier admission on the same queries. When the
context ends the unit is abandoned, and nothing it produces is read or
published. Abandoned work stops at the next step it reaches between two
`go/types` calls. A call it is inside keeps its goroutine and its processor
until that call ends by itself, and an abandoned admission keeps later
admissions on its queries waiting, each until its own context ends.

The profile's limits bound sizes, not this cost, for captured code as for
supplied source. Comparing two structurally equal chains of struct aliases
doubles with every step. As captured declarations in the base profile such
chains took 0.6 seconds at depth 20 and 8.8 seconds at depth 24, and a supplied
package of 1,618 bytes took 11 seconds to check. A consumer that must bound
that work, for source it does not trust, runs the provider in a process it can
end.

Those incomplete reports cover the provider's own input, conversion and export
bounds. The shared capture/publication boundary can also return `ErrLimit` when
the caller's input or result budget is too small; that is not a native type
rejection. No budget is raised implicitly.

The graph contains package/type/alias/parameter/field entities, resolved name
bindings, expression terms and generic applications. Go-owned facts retain
underlying types, alias targets, constraints, instantiated fields and native type
identity. Aliases remain distinguishable from definitions even when Go considers
their types identical. IDs belong to that exact model, never to a spelling.

The vocabulary owner is `go`. The published representation is:

| Graph element | Meaning |
| --- | --- |
| `package`, `defined`, `alias`, `parameter`, `field`, `builtin` entities | Package and declaration identities; parameters carry their declaration owner and ordinary `type` sort |
| `Reference` / `Application` terms | Bound named types/parameters and effective generic arguments; the shared `DirectParameterBindings` reader applies |
| `name` bindings / `type-expression` occurrences | Resolved identifier targets and native type terms at current AST nodes |
| `underlying`, `alias-target`, `constraint` facts | Go's checked relationships, including substituted underlying types of applications |
| `identical` fact | A model-local representative established by `types.Identical`, not identity inferred from spelling or serialized shape |
| `package-name` fact | On every `package` entity, captured and supplied: the name that its package clauses declare, `types.Package.Name()`. The entity's own name is the package's path, which does not determine the declared name |
| `basic`, `slice`, `struct`, `interface`, `union` expressions | Native structure with ordered operands and Go-owned `json-1` metadata |

Struct metadata records field names, tags, embedding flags and model-local field
entities; corresponding operands are their types. Interface metadata records
embedded-term count, explicit method names, comparability and whether it is a
method-set interface. Union metadata records each term's approximation flag.
The builtin `error` interface needs signature/tuple facts internally, even though
authoring method/function syntax remains unsupported. `display` facts are UTF-8
diagnostic text, never semantic identity. Typed readers return copies or immutable
handles; callers cannot mutate a published graph through them.

### Import environment

`golang.ImportConfiguration(golang.ImportEnvironment{Module, OutputRoot, Packages})`
extends the base configuration with one normalized `imports-2` setting. It
follows Go's own module layout:

- **The output root.** `Module` is the module path and nothing else. The output
  root is the directory that the paths of captured files are relative to.
  `OutputRoot` states where it sits in the module: its directory relative to
  the module root, such as `generated` or `internal/gen`. When the output root
  is the module root, that is stated as `golang.ModuleRoot`, which is `.`.
- **Captured packages.** Each captured tree must be linked to a `.go` file in the
  request's Files context. Its package import path is the module path, the
  output root and the file's directory, joined. `pages/types.go` under
  `example.org/app` belongs to `example.org/app/generated/pages` with the
  output root `generated`, and to `example.org/app/pages` with the output root
  `.`. Files in one directory form one package.
  `golang.CapturedImportPath(environment, filePath)` is this derivation, and
  the bridge groups captured files by the same function.
- **Supplied packages.** Every other imported package is supplied as Go source:
  an import path and up to 16 files. Supplied text participates in receipt
  identity and is never an output. A supplied package may import only other
  supplied packages, and it may not use a path at or below the output root,
  where the captured packages are. A package of the module outside the output
  root is hand-written code of the consumer's own module and may be supplied:
  `example.org/app/model` beside the output root `generated`. With the output
  root `.` every path under the module is refused.
- **No discovery.** Import paths need a first element containing a dot, Go's
  module-path rule. That also keeps supplied packages from naming standard
  library paths, which this profile does not provide. Nothing is read from a
  module cache, `GOPATH` or the network.

**The output root must be stated.** An environment without `OutputRoot` is
refused with `ErrInvalidRequest`, and no output root is assumed. The module
root has the one spelling `.`, which a caller writes; an empty field is never
read as it. Every other value is a clean relative path with forward slashes
whose elements are import-path elements. Refused are an absolute path, a
backslash, a trailing slash, an empty element and a `.` or `..` element, so one
directory has one spelling and equal environments give equal settings. Refused
too is an output root with a `testdata` or `vendor` element or an element that
starts with `_` or `.`: the rules below keep such a directory out of a package
tree for a captured file, and an output root inside one would put every
captured file there.

**The setting's bytes.** The key is `imports-2`. The value is the UTF-8 JSON
object with the members `Module`, `OutputRoot` and `Packages` in this order.
For an environment without supplied packages it is
`{"Module":"example.org/app","OutputRoot":"generated","Packages":[]}`. The
setting is part of every receipt under an import environment. The earlier key
`imports-1` carried no output root. A setting under that key is refused, and
so is a value under `imports-2` that states none: neither is read with an
assumed output root. `TestGoImportEnvironmentMustStateTheOutputRoot` pins the
key, the bytes and both refusals.

**The go tool's layout rules apply to captured files:**

- test files (`_test.go`), files and directories starting with `_` or `.`, and
  `testdata` directories are not part of a package, and are unsupported;
- `vendor` directories, whose import paths Go resolves differently, are
  unsupported, and a supplied import path may not contain a `vendor` element;
- an import of a path with an `internal` element is allowed only from within the
  tree rooted at its parent;
- a package named `main` cannot be imported;
- import paths, and file names within one package, may not differ only in case.

Supplied files follow the same file rules: a supplied package lists package
source files only, never test or ignored files.

Build constraints, including GOOS and GOARCH file-name suffixes, are not
interpreted: every captured file belongs to its directory's package.

**What the first Go Files slice admits by name.** A build selects a Go file by
its name, and the profile above does not look: it analyzes `types_linux.go` as
part of its package, and a build for another system drops that file. A Go
Files stage that chooses file names must therefore not write such a name. Two
functions of `semantics/golang/go` state what the first slice admits:

- `golang.AdmitFileName(filePath)` refuses a base name with an underscore, or
  one that starts with a dot. That refuses every GOOS and GOARCH suffix, every
  test file and every file the go tool ignores, and more: `types_gen.go` is
  refused although Go builds it. The rule needs no list of systems that
  changes between toolchains. The directory is not judged.
- `golang.AdmitPackageName(name)` admits ordinary packages. It refuses a
  command, which is package `main`, and a test package, whose name ends in
  `_test`.

Both return `nil` for an admitted name. A name that Go accepts and the slice
does not cover gives an error that wraps `golang.ErrUnsupported`, which a Files
stage reports as unsupported coverage. What is no Go file path or no package
name gives one that wraps `semantic.ErrInvalidRequest`. A later slice may admit
a refused name and may not give an admitted name another meaning. The
functions do not change the semantic profile: the bridge refuses what it
refused before, and a file with a system suffix stays analyzable.

**Classification.** The checks run in a fixed order, so no result or diagnostic
depends on capture order:

1. The environment: the Files layout of every tree, then every supplied package,
   parsed and checked in its own import order.
2. Captured file paths, in file-path order: test, ignored and vendored files and
   invalid package directories. These come before any captured resource limit.
3. Captured resource limits, then each tree's syntax, in file-path order:
   unsupported constructs, including dot and blank imports.
4. Imports outside the environment.
5. The go tool's rejections: case collisions of import paths and file names,
   then disagreeing package clauses, then the `internal` and `main` rules and
   import cycles.
6. The type check, in import order through a closed importer.

| Case | Result |
| --- | --- |
| A supplied package that fails to parse or type-check, imports outside the supplied packages, breaks the `internal` rule, imports a `main` package or forms an import cycle; supplied paths or file names that differ only in case, a supplied test or ignored file, or a supplied path with a `vendor` element; a missing Files layout, or a captured file that is not a `.go` path | `ErrInvalidRequest`; no judgment of the captured code |
| An import outside the environment, a dot or blank import, a test, ignored or vendored file, a package directory that is not a valid import path | Unsupported coverage at the import or file |
| A captured type error, an unused import, an unexported cross-package use, an import cycle through a captured package (including a self-import), an `internal` violation, an import of a `main` package, inconsistent package clauses, captured import paths or file names that differ only in case | Rejected validity, with a diagnostic at the captured node |

The graph gains:

- one `package` entity per captured package, in import-path order, owning that
  package's declarations, followed by one per supplied package;
- an `import` binding from each captured import declaration to its package's
  entity;
- `name` bindings from qualified type references to the imported declarations.

**Supplied structure stays outside the graph's projection.** The graph projects
captured code; supplied packages take part by identity. A supplied declaration
contributes its entity, the entities of its type parameters with their names,
its reference and application terms and its `display` and `identical` facts,
never its structure: no fields, no `underlying`, `alias-target` or `constraint`
facts. Native identity is still go/types identity: an `identical` fact relates a
supplied alias to its target when that target is in the graph, which `string`
and `bool` always are.
This also holds when a captured declaration is built directly on supplied
structure: a captured alias of a supplied type, a captured declaration whose
right-hand side is one (through chains of any length), or a type parameter
constrained by one. So supplied internals may lie outside this profile without
failing a captured result.

The purpose is that containment, not secrecy: the consumer supplies these
packages itself. Two kinds of output may therefore depend on supplied
declarations:

- **Compiler-established properties of captured types**, such as the comparability
  or emptiness of a captured constraint that embeds a supplied interface, and
  go/types identity in `identical` facts.
- **Human-readable text**: diagnostics, `display` facts and admission detail may
  name supplied types and spell their constraints. This text is never semantic
  identity.

`Type.Unconstrained` reports only what the graph establishes, and a parameter it
reports is one the compiler found unconstrained. It never reports a supplied
parameter, nor a captured parameter whose constraint is a supplied interface or
an alias of one, because the graph holds no structure for that constraint. A
captured constraint that embeds a supplied interface is captured structure and
is reported as the compiler finds it. An admission query's constructor must be a
generic defined type: an instance, a non-generic type and a generic alias are
refused as unsupported.

**Lookups.** `Model.Lookup(name)` considers captured packages only, and refuses a
name that more than one of them declares. `Model.LookupIn(importPath, name)`
reads a package-qualified declaration. A supplied declaration is present only
where captured code uses it.

**Package name and identity.** `Model.PackageName(importPath)` reads the
`package-name` fact of a captured or supplied package. `Type.Package()` reads
the package of a package-level declaration from the graph's ownership: its
identity, which is its import path, and its declared name. An instance of a
generic type has the package of its generic declaration. Both readers fail
when the package entity does not carry exactly one name fact, and neither
reads the path in place of the name. In the base profile the identity is the
constant `captured`, and the name is the one the package clause declares.

**The limit of package identity.** A consumer that relocates declarations can
compare the owner of each declaration, `Type.Package().Path`, with
`CapturedImportPath` of the file it intended. That comparison establishes that
the declaration sits in the tree of the intended file and that the checker put
it into that package. It cannot certify where the consumer writes the files.
The type checker is told each package's path and does not discover a module
layout, and the expected path and the bridge's grouping come from the one
function. Whether the files on disk sit at `OutputRoot` inside the module is
the consumer's statement, which nothing here checks.

**One environment end to end.** Analysis refuses a configuration that differs
from its resolution's. Use `golang.AnalysisContractFor(settings)` and
`golang.EvidenceFor(settings)` with the same configuration for resolution and
analysis. `ImportConfiguration` refuses an environment that one analysis
invocation could not carry; it charges the settings twice, once for the analysis
and once for its resolution. An export-time limit or unsupported construct is
reported at the captured declaration or use being exported.

A consumer explicitly constructs both adapters, resolves with
`semantic.Options{Configuration: golang.Configuration()}`, and calls the analyzer
with that same configuration and `Resolution: &resolved`. `golang.Open(analyzed)`
then opens the typed view. A resolution receipt alone, a receipt from a different
capture or a different compiler configuration cannot substitute for that input.

`semantics/golang.Open` admits the analysis receipt and reads this graph without
compilation. Native `queries.NewQueries(ctx, model)` from
`analyze/files/golang/queries/go` is a separate, explicit operation: it rebuilds
compiler state from the same capture and environment and checks correspondence
with the published graph. A query session retains that state privately, with no
compiler pointers in the generic graph, hidden cache or callback attachment.
`Admit` checks generic instantiation with `types.Instantiate(..., true)` and binds
its answer to the model, constructor and arguments. It proves native argument
admission only, not Horos membership or representation preservation. Query sessions
serialize compiler queries; immutable semantic views may be read concurrently.

## First TypeScript profile

**6 November 2026:** TypeScript uses the same resolution and analysis interfaces
as Go. The generation work's new Source and transformation capture accessors do
not alter those contracts. Native semantic receipts remain independent of
generation support and representation-preservation claims.

```text
syntax-admitted TypeScript ASTs + captured Files and links
  -> resolve/files/typescript/provider/go
  -> private source/AST correspondence check + pinned TypeScript compiler
  -> immutable symbols, native type terms, bindings and reports
  -> analyze/files/typescript/provider/go (explicit receipt; reuses the exact graph)
  -> semantics/typescript/go (typed read-only view, graph names, contract identities, compiler configuration)
  -> analyze/files/typescript/queries/go (explicit native queries over an opened model)
```

The package roles mirror Go. `semantics/typescript/go` publishes the contract identities
`ResolutionContract()`, `AnalysisContract()` and `AnalysisContractFor(settings)`
and the `Compiler` configuration type. It also publishes the sixteen graph names
that the compiler helper writes and Go reads by name, such as `ModuleKind()`,
`NamedExportRelation()` and `NameRole()`. The helper is JavaScript and spells them
itself, and tests that run it keep the two sides equal. The query service is its own package,
`analyze/files/typescript/queries/go`. Go has the same two placements since pull request
[#246](https://github.com/Bitspark/bittype/pull/246) of 12 November 2026. Go publishes
`AnalysisContractFor(settings)` too, since pull request [#204](https://github.com/Bitspark/bittype/pull/204), and has no `Compiler` type. `internal/typescriptsem` owns the compiler bridge,
including its embedded CommonJS helper. Each TypeScript stage supplies its own provider. No TypeScript types, rules or compiler handles enter the
agnostic kernel; no native stage imports Horos, Logos or generation packages.

The resolver constructor requires a `typescript.Compiler` containing absolute
`NodePath` and `TypeScriptPath` paths. The latter names the explicitly provisioned
TypeScript package's `lib/typescript.js`. Construction does no discovery,
installation or IO. Resolution invokes that compiler and requires version
**5.9.3**. The caller supplies `typescript.Configuration()` in both operations;
the configuration fixes ES2020, ESNext modules, strict checking,
`exactOptionalPropertyTypes`, `noUncheckedIndexedAccess`, `lib.es2020.d.ts`, no
ambient packages, and the `captured-types-1` helper profile. The baseline host
reads only supplied source units and that compiler package's standard libraries.
Compiler files are trusted, explicitly provisioned dependencies.
The child process receives no `NODE_OPTIONS`, `NODE_PATH` or credential environment.

The named-import setting from `typescript.ImportConfiguration(packages)` additionally
admits local `.js` specifiers resolved against captured `.ts` files and named
exports of explicitly supplied package declaration modules. Both the package specifiers
and their declaration text participate in receipt identity. Use `typescript.AnalysisContractFor(settings)`
with the same configuration for resolution and analysis. The graph binds a named use
to its compiler symbol and records the imported package type's module, export and
generic arguments. Bittype checks the syntax link to the local import binding.
Package-private member structure is outside the graph projection. The Files formatter
establishes these facts again after placement. [TypeScript type imports](typescript-type-imports.md)

Every supplied AST must correspond to its captured original `.ts` text. The
private bridge reparses that text using Bittype's existing TypeScript parser,
compares kinds, attributes and containment, and recomputes offsets before
associating native compiler syntax. UTF-8 source spans become UTF-16 compiler
offsets. Caller origins are not proof of correspondence. Supplied resolved
references must agree with compiler symbols. No formatter or printer is involved;
checking generated text would require an explicit new capture.

<!-- bittype:overview scope=first-typescript-profile part=features -->
The first profile covers aliases, interfaces and bases, generic parameters and
constraints, explicit applications, references, primitive/literal types, arrays,
tuples, object properties, unions, intersections and readonly arrays/tuples.
Script files share global scope; exported declarations give each file module
scope. Binders and declarations with the same spelling remain distinct symbols.
Named recursive types retain entity/fact links without cyclic term containment.
<!-- bittype:overview:end -->

<!-- bittype:overview scope=first-typescript-profile part=exclusions -->
Without the named-import configuration, imports and external package types are
unsupported. With it, named local/package imports and their supplied declarations
are admitted as described above. Qualified names, callable/method/index signatures, indexed access,
`keyof`, `unique`, `this`, generic defaults, interface declaration merging and
external types outside that explicit import environment or compiler-recognized array constructors remain
outside this first profile. Standard arrays and readonly arrays are supported;
mapped utility types such as `Readonly<T>` are not. Their synthesized modifiers
need a separate extension of native fact export. Authored conditional and mapped
types are also unsupported even though the syntax provider admits them.
<!-- bittype:overview:end -->

The existing syntax parser rejects diagnostic-control comments (`@ts-nocheck`,
`@ts-ignore`, `@ts-expect-error`, `@ts-check`) and source dependency directives.
Supplying an AST separately does not bypass that rule: correspondence reparses
the original source and rejects the incompatible capture with `ErrInvalidRequest`,
before invoking the compiler. No duplicate compiler guard is needed.
Directive-looking text inside a string literal has no such effect.
This is bounded native checking, not support for
all TypeScript syntax or arbitrary `.d.ts` dependency closures. Generic
applications must supply their complete arguments; there is no public dynamic
generic-instantiation query in this slice.

The `typescript` vocabulary is:

| Graph element | Meaning |
| --- | --- |
| `program`, `script`, `module`, `alias`, `interface`, `parameter`, `property`, `builtin` entities | Compiler symbol identities; declaration sites and real generic binder owners |
| `name` bindings | Native lookup results at captured reference nodes |
| `type-expression` occurrences | Source reference/application structure, preserving constructor identity even when aliases disappear from native types |
| `named-export` fact on a module | A reference term for each named type the compiler exports from the captured module. Files checks requested public sites against these facts after placement. |
| `native-type` occurrences and marker facts | Native compiler types; model-local handles, not serialized shape identity |
| `declared-type`, `constraint` facts | Declaration types and explicit generic bounds |
| `structure`, `properties` facts | Native structural descriptions and substituted declared property value types |
| `primitive`, `literal`, `array`, `tuple`, `union`, `intersection`, `object` expressions | TypeScript-owned `json-1` data and ordered operands |
| `display` facts | UTF-8 diagnostic spelling, never type identity |

Property metadata retains optional and readonly flags and declaration entities.
Its operand is the compiler's declared non-missing value type: `age?: number`
exposes `number` plus `Optional: true`, while `age: number | undefined` exposes
the explicit union plus `Optional: false`. This is not a flow-sensitive read type
or a writable acceptance judgment; native assignability uses the complete
compiler type. This follows the pinned checker's
[symbol-at-location behavior](https://github.com/microsoft/TypeScript/blob/v5.9.3/src/compiler/checker.ts).
Aliases retain
declaration identities but can denote the same compiler type. `Type.Same` means
the same native type handle in the same admitted receipt, not mutual assignability.
The shared `DirectParameterBindings` reader works on source application terms;
the corresponding native type can retain an alias or normalize to another form.

`typescript.Open(analyzed)` admits a complete analysis receipt and reads its
immutable graph without invoking the compiler. `queries.NewQueries(model,
compiler)` from `analyze/files/typescript/queries/go` separately configures
explicit native queries. `Assignable(ctx,
source, target)` calls the official checker's directional `isTypeAssignableTo`.
It reconstructs compiler state from the original capture and exact environment,
requires correspondence with the entire published graph, and returns a result
bound to that model and both type handles. It does not test Horos membership or
representation preservation. Each query owns a fresh process; concurrent queries
share only immutable facts. There is no persistent compiler service or cache.

Missing Files context reports incomplete coverage; unsupported constructs report
unsupported coverage; compiler type errors report complete coverage and rejected
validity. None publishes partial semantics. AST/source contradictions, forged
references and wrong configuration are request errors. Missing/wrong compiler
executables or modules are operational compiler errors. Resource limits yield
incomplete coverage or a boundary `ErrLimit`, never a false type judgment.

The profile bounds captures to 16 units, 4,096 nodes and 1 MiB of source, 16
parameters per declaration, 32,768 semantic records and depth 128. Each compiler
process has a 30-second deadline, a 16 MiB response limit and a 16 KiB stderr
limit. Cancellation terminates the child. The shared capture/publication limits
still apply independently. Tests require explicitly provisioned compiler paths;
missing compilers fail rather than silently skipping semantic qualification.

The [provider tests](../integration/go/typescript_semantics_test.go) cover nested
pagination, constraints, scope, recursion, immutable facts, source correspondence,
receipt isolation, assignability and concurrency. The
[installed consumer](../integration/go/testdata/installed/typescript_semantics_test.go)
also qualifies the embedded helper from a fresh public-module installation.
The remote gate runs these with TypeScript 5.9.3 and typechecks the helper itself.

## First Rust profile

**6 November 2026:** Rust implements the same independent resolution and analysis
contracts. Its syntax and generation providers already existed; this adds native
semantic checking without changing generation or runtime behavior.

```text
syntax-admitted Rust AST + captured Files and links
  -> resolve/files/rust/provider/go
  -> original-source/AST correspondence + pinned rustc_driver helper
  -> immutable native declarations, types, bindings, predicates and reports
  -> analyze/files/rust/provider/go (explicit receipt; reuses the exact graph)
  -> semantics/rust/go (typed read-only view)
```

The public adapters and orchestration are Go. The private helper at
`internal/rustsem/helper` uses Rust's compiler APIs, pinned to **1.90.0**.
Its Cargo package is unpublished. `semantics/rust/go` publishes the contract
identities `rust.ResolutionContract()` and `rust.AnalysisContract()`.
`resolve/files/rust/provider/go.New` takes a
`resolve.Compiler` with absolute `HelperPath` and `SysrootPath`. Constructors
perform no IO, discovery or installation. The caller provisions that helper and
the matching compiler sysroot explicitly. Both operations require
`rust.Configuration()`: edition 2024, `x86_64-unknown-linux-gnu`, crate `captured`,
library checking, sysroot-only dependencies, no code generation, and the
`captured-types-1` helper profile. Compiler version, profile and environment are
receipt identity; executable locations are configuration of the collaborator.

The remote gate builds the helper with Rust 1.90.0, `rustc-dev` and
`llvm-tools-preview`; `RUSTC_BOOTSTRAP=1` enables the private compiler API only
while building that helper. The checked input gets no feature-enabling environment.
The helper links its explicit compiler library location and receives an empty
environment from the Go bridge. It supplies original source through `Input::Str`,
denies filesystem source lookup and stops after compiler analysis, before codegen.
There is no Cargo project discovery, external package resolution, dependency
download, build script, proc macro or source emission during semantic checking.
Compiler/helper installation is a trusted prerequisite. This first executable
profile is qualified on Linux x86-64; other hosts require deliberate qualification.

The private Go bridge reparses the original `.rs` capture with the existing parser
and compares node kinds, attributes and containment. It recomputes spans, maps
CRLF/BOM normalization to compiler byte positions, and ignores caller origin claims.
Supplied reference targets must agree with native path resolution. The existing
syntax parser excludes attributes, macros and executable bodies. No source is
printed from ASTs to obtain a compiler input; generated text needs a new explicit
capture. The helper distinguishes compiler-injected standard prelude declarations
from authored declarations.

<!-- bittype:overview scope=first-rust-profile part=features -->
The supported semantic fragment is one captured source file with inline modules,
type aliases, structs/tuple structs, enums with unit/tuple/record variants, fields,
ordinary type parameters and defaults, effective generic applications, paths,
primitives, tuples, numeric-length arrays, slices, raw pointers and explicit
`'static` references. Marker traits, supertraits, ordinary trait bounds/where
predicates and empty type-only trait impls are checked by the compiler. Sysroot
types use their compiler declaration identities: `Vec`, `Option`, `String`, `Box`
and similar types have no name-specific mapping in Bittype. Their internals stay
opaque; effective arguments include native defaults, such as an allocator argument.
<!-- bittype:overview:end -->

<!-- bittype:overview scope=first-rust-profile part=exclusions -->
Imports and multi-file crates, dependency closures, methods/functions/callable
types, trait objects, associated types/projections/GATs, lifetime parameters,
non-static references, const generics and higher-ranked predicates are outside
this first semantic profile. Some have syntax tooling; that does not imply
semantic coverage. An unsupported source construct or exported compiler form
reports unsupported coverage without a partial graph. The standard compiler
prelude and sysroot are part of the explicit environment, not a library-name
switch in the kernel.
<!-- bittype:overview:end -->

The provider-owned `rust` vocabulary is:

| Graph element | Meaning |
| --- | --- |
| `module`, `alias`, `struct`, `enum`, `variant`, `field`, `trait`, `impl`, `parameter` entities | Native declaration identities, captured sites and owners |
| `Reference` / `Application` | Native types and effective arguments; shared `DirectParameterBindings` applies |
| `target` bindings | Compiler-resolved path segments at captured nodes |
| `declared-type` facts/occurrences | Checked type of each captured type declaration or field |
| `native-type`, `display` facts | Native type handles and diagnostic spelling; display is never identity |
| `fields`, `variants` facts | Ordered local ADT structure, substituted field types, field names/visibility and declaration identities |
| `default`, `predicate`, `implements` facts | Generic defaults, declared trait predicates and checked impl headers |
| `primitive`, `tuple`, `slice`, `array`, `pointer`, `reference` expressions | Rust-owned `json-1` metadata and ordered type operands |

Generic parameters retain their actual owning declaration and slot. Rust's
implicit trait `Self` is a distinct binder even though rustc internally reuses
the trait's `DefId`. Same-spelled `T` parameters in different declarations never
collapse. The native type cache includes the binder environment for parameterized
types. Recursive local types use entity/fact links, not cyclic term containment.
Aliases keep their declaration identities while their native type can normalize
to another declaration. Source path bindings still identify the alias. This slice
does not export a semantic occurrence for every nested source type expression.

`rust.Open(analyzed)` admits the exact completed analysis receipt and reads the
graph without a compiler invocation. `Lookup` takes a compiler path such as
`captured::models::Person`; paths select declarations but do not define equality.
Ambiguous paths across native namespaces require selecting the actual entity ID.
`Type.Same` means compiler type identity with the same generic binder environment
in that receipt. It does not imply coercibility, trait satisfaction or Horos
representation preservation. Declared predicates are facts about declarations,
not proofs of every possible specialization. Rust does not enforce ordinary
type-alias bounds at every use site; this provider preserves rustc's behavior.
There is no dynamic admission/coercion query API in this first Rust slice.

Missing Files context reports incomplete coverage. Native source errors report
rejected validity with a captured-node diagnostic. Unsupported forms report
unsupported coverage. None publishes partial semantics. Contradictory source/AST
captures, false references, wrong configuration and foreign receipts are request
errors; helper/toolchain failures are operational errors. Cancellation stops the
child process. The bounds are 4,096 AST nodes and 1 MiB source, 16,000 native type
handles, 32,768 semantic records, export depth 128, a 30-second process deadline, 16 MiB response and
256 KiB stderr. Shared capture/publication budgets apply independently.

The [provider tests](../integration/go/rust_semantics_test.go) cover pagination,
nested substitution, defaults, binder ownership, recursion, module paths, marker
traits, invalid/unsupported/incomplete outcomes, source correspondence and receipt
isolation. The [installed consumer](../integration/go/testdata/installed/rust_semantics_test.go)
uses the published Go module with an explicitly provisioned helper. Immutable
graph reuse is exercised under the race detector. Missing compiler provisioning
fails qualification rather than silently skipping it.

## Extension to C++ and broader native coverage

The same roles and receipt contracts apply, with language-owned facts and query
semantics. There is no universal `Compatible` predicate and no Bittype
reimplementation of all native type systems.

| Language | Planned authority and initial semantic surface |
| --- | --- |
| Rust extension | Broader dependency/lifetime/associated-type coverage and explicit native queries on top of the pinned helper. Compiler internals remain unstable. |
| C++ | Pinned Clang LibTooling helper; declarations, template parameters/applications and constraint checking. |

Public orchestration stays in Go; compiler helpers may use their native language.
Compiler version, language edition, flags, target and dependency environment must
be explicit. Helpers must not discover or download an ambient dependency closure.
Compilers that require source consume captured files with verified AST/file
correspondence. Checking emitted source is an explicit new capture; it cannot be
silently printed/reparsed and presented as analysis of the original AST.

Future providers distinguish missing context, unsupported coverage, resource
exhaustion and proven invalidity. Their first slices should publish no partial
models. Native semantic checking remains separate from
[type generation](semantic-generation.md), [native admission](go-native-admission.md)
and the [joint verification contract](verification.md).

Compiler integration references: [Go types](https://pkg.go.dev/go/types),
[TypeScript Compiler API](https://github.com/microsoft/TypeScript/wiki/Using-the-Compiler-API),
[Rust driver](https://rustc-dev-guide.rust-lang.org/rustc-driver/intro.html), and
[Clang LibTooling](https://clang.llvm.org/docs/LibTooling.html).
