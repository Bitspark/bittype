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
| `semantics/golang/go` | Go-owned vocabulary, contract identities, environment identity and immutable typed graph readers |
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

The `package-name` fact was added on 13 November 2026 under the unchanged
profile identity: the profile strings and the base configuration are what they
were, and every graph, the base profile's too, carries the fact from then on.
A reader that needs the fact fails with an error on a graph that lacks it. The
import environment's setting has the new key `imports-2`.

**13 November 2026.** The vocabulary publishes the kind `builtin` as
`BuiltinKind()` and the role `name` as `NameRole()`, and the bridge writes both
from there. No graph changes. An entity of the kind `builtin` is a predeclared
basic type, such as `string` and `bool`: its name is the identifier of the
universe scope, and it has no owner and no declaration site. A generated use of
such a type records the predeclared symbol as its
[destination](destinations.md#go), and that entity is what the destination is
compared with. A declaration or a type parameter of the same name is another
entity, of the kind `defined`, `alias` or `parameter`, with a declaration site.
Limit: `builtin` covers the basic types. `error` and `comparable` are named
types of the universe scope and `any` is an alias, so their entities have the
kinds `defined` and `alias`, without an owner and without a site. A reader
tells one of them from a declaration by "no owner and no site", and no marker
says it.

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
  `golang.CapturedImportPath(environment, filePath)` is this derivation. The
  bridge derives each captured file's package once with it and groups the
  files by the result.
- **Supplied packages.** Every other imported package is supplied as Go source:
  an import path and up to 16 files. Supplied text participates in receipt
  identity and is never an output. A supplied package may import only other
  supplied packages, and it may not use a path at or below the output root,
  where the captured packages are. A package of the module outside the output
  root is hand-written code of the consumer's own module and may be supplied:
  `example.org/app/model` beside the output root `generated`. With the output
  root `.` every path under the module is refused. A supplied package of the
  module is below the module root, so its path there follows the rule for
  unsupported directories below: `example.org/app/testdata/x` is refused. A
  supplied path that is the module path, or a path below it, only when case
  is ignored is refused with `ErrInvalidRequest`, whatever its elements are:
  `Example.org/App/model` and `example.org/App`. To Go it is the path of
  another module, and on a file system that ignores case it collides with the
  module's own tree.
- **No discovery.** Import paths need a first element containing a dot, Go's
  module-path rule. That also keeps supplied packages from naming standard
  library paths, which this profile does not provide. Nothing is read from a
  module cache, `GOPATH` or the network.

**The output root must be stated.** An environment without `OutputRoot` is
refused with `ErrInvalidRequest`, and no output root is assumed. The module
root has the one spelling `.`, which a caller writes; an empty field is never
read as it. Every other value is a clean relative path with forward slashes.
An empty, `.` or `..` element is refused with `ErrInvalidRequest`, which
covers an absolute path and a trailing slash. So one directory has one
spelling, and equal environments give equal settings.

**An element of the output root follows Go's rule for an element of an import
path**, as `checkElem` states it in `golang.org/x/mod/module` (`module.go:417-482`
of the pinned go1.27.0 toolchain): ASCII letters, digits and `-._~`, no
trailing dot, and, in the text before the first dot, no Windows device name
in any case and no tilde that only digits follow to its end. `con`, `NUL`,
`com1`, `lpt9.x`, `gen~1` and `~1` are refused with `ErrInvalidRequest`.
`console`, `x.con`, `gen~x` and `gen~` are admitted, as Go admits them. The
rule here is narrower than Go's in two places, a leading dot and a plus sign,
and wider in none. One function of `semantics/golang/go` holds it.
`ValidImportPath`, which judges module, supplied and captured paths, does not
apply the device names and the tilde form yet
([#287](https://github.com/Bitspark/bittype/issues/287)). Two roots are
admitted that no file system or module can hold, because Go's rule for an
element does not refuse them either: an element longer than 255 bytes, and
`go.mod` as the first element, where the module's own file of that name
stands.

**A root's validity depends on the length of the module path.** The bound of
512 bytes on an import path holds the module path, the output root and a
file's directory together. A root that exceeds it with the module path is
refused with `ErrInvalidRequest`. Under a root that fills it, a file in a
subdirectory has no package, and the bridge reports unsupported coverage at
that file.

**One rule for the directories the profile does not cover.**
`golang.UnsupportedDirectory(element)` is its one declaration: `testdata`,
`vendor`, and a name that starts with `_` or `.`. The go tool leaves
`testdata` and a leading underscore or dot out of the `./...` pattern and still
builds such a package when it is imported explicitly. It refuses an import
through a `vendor` directory. The profile covers none of them. The rule
applies to every directory below the module root, and one enters at three
places. At each the answer is unsupported:

- in the path of a captured file, the bridge reports unsupported coverage at
  the file, as it did before the rule had one owner, and `CapturedImportPath`
  gives no identity for the file. The derivation judges the file's
  directories and not its own name;
- as an element of the output root, `ImportConfiguration` returns an error
  that wraps `ErrUnsupported` and writes no setting. The bridge reads the
  Files path, so without this refusal the output root `vendor` would put every
  captured file at a vendor path unnoticed. `CapturedImportPath` gives no
  identity under such a root;
- as an element below the module root in the path of a supplied package of
  the module, `ImportConfiguration` answers the same way:
  `example.org/app/testdata/x`, `example.org/app/_x`, `example.org/app/vendor/x`
  and `example.org/app/.x` under the module `example.org/app`. Only a path
  that is below the module root as it is written is judged here. One that
  differs from it in case is an invalid request before that:
  `Example.org/App/testdata/x`.

A supplied path outside the module is not below the module root, and the
directory rule does not judge it. There the earlier rule holds unchanged: a
`vendor` element is `ErrInvalidRequest`, and `example.org/lib/testdata/x` and
`example.org/lib/_x` are admitted.

An output root is judged in three steps, in this order: a path that is not
clean is invalid, then a directory of the rule is unsupported, then an element
that is no import-path element is invalid. So `./vendor` and `vendor/` are
invalid, and `.generated`, `a b/vendor` and `C:/vendor` are unsupported
although the last two are no path that Go accepts. The path of a supplied
package of the module below the module root is judged in the same order,
before the rule for the output root.

**The setting's bytes.** The key is `imports-2`. The value is the UTF-8 JSON
object with the members `Module`, `OutputRoot` and `Packages` in this order.
A package has the members `ImportPath` and `Files`, and a file `Name` and
`Text`. The setting is part of every receipt under an import environment.
`TestGoImportEnvironmentMustStateTheOutputRoot` pins the key and three literal
values:

- `{"Module":"example.org/app","OutputRoot":"generated","Packages":[]}`
- `{"Module":"example.org/app","OutputRoot":".","Packages":[]}`
- `{"Module":"example.org/app","OutputRoot":"generated","Packages":[{"ImportPath":"example.org/cell","Files":[{"Name":"cell.go","Text":"package cell\n"}]}]}`

The earlier key `imports-1` carried no output root. A setting under that key
is refused, and so is a value under `imports-2` that states none: neither is
read with an assumed output root. The same test builds eight settings by
hand and requires `ErrInvalidRequest` for each from `Imports`,
`AnalysisContractFor` and the resolver: the old key with the old and with the
new bytes, the new key with no root, an empty root, an unclean root, the
members in another order, a root in a `vendor` directory and a root that is a
device name. A setting is what `ImportConfiguration` writes. It writes none
for a root it refuses, so a hand-written one that states such a root is
refused as an altered setting, whatever the reason for the root's refusal.
The same holds for a supplied path that it refuses:
`TestGoSuppliedPackageInsideTheModuleAndOutsideTheOutputRoot` states each of
its refused paths in hand-written bytes and requires `ErrInvalidRequest` from
the same readers.

**The go tool's layout rules apply to captured files:**

- test files (`_test.go`) and files starting with `_` or `.` are not part of a
  package, and are unsupported;
- a file in a directory that `UnsupportedDirectory` names is unsupported, and
  a supplied import path may not contain a `vendor` element, below the
  module root by that rule and elsewhere by its own;
- an import of a path with an `internal` element is allowed only from within the
  tree rooted at its parent;
- a package named `main` cannot be imported;
- import paths, and file names within one package, may not differ only in case.

Supplied files follow the same file rules: a supplied package lists package
source files only, never test or ignored files.

Build constraints, including GOOS and GOARCH file-name suffixes, are not
interpreted: every captured file belongs to its directory's package.

**Classification.** The checks run in a fixed order, so no result or diagnostic
depends on capture order:

1. The environment: the Files layout of every tree, then every supplied package,
   parsed and checked in its own import order.
2. Captured file paths, in file-path order: test and ignored files, files in
   an unsupported directory and invalid package directories. These come before
   any captured resource limit.
3. Captured resource limits, then each tree's syntax, in file-path order:
   unsupported constructs, including dot and blank imports.
4. Imports outside the environment.
5. The go tool's rejections: case collisions of import paths and file names,
   then disagreeing package clauses, then the `internal` and `main` rules and
   import cycles.
6. The type check, in import order through a closed importer.

| Case | Result |
| --- | --- |
| An environment without an output root; an output root that is not `.` and has an empty, `.` or `..` element, an element that is no import-path element, or with the module path more than 512 bytes; a supplied path at or below the output root; a supplied path that is the module path or a path below it only when case is ignored | `ErrInvalidRequest` from `ImportConfiguration`; no setting |
| An output root with a `testdata` or `vendor` element, or one that starts with `_` or `.` | An error that wraps `ErrUnsupported` from `ImportConfiguration`; no setting |
| A supplied package of the module whose path below the module root has a `testdata` or `vendor` element, or one that starts with `_` or `.` | An error that wraps `ErrUnsupported` from `ImportConfiguration`; no setting |
| A setting that `ImportConfiguration` did not write: the key `imports-1`, altered bytes, or a root or a supplied path that it refuses for either reason | `ErrInvalidRequest`; no judgment of the captured code |
| A supplied package that fails to parse or type-check, imports outside the supplied packages, breaks the `internal` rule, imports a `main` package or forms an import cycle; supplied paths or file names that differ only in case, a supplied test or ignored file, or a supplied path with a `vendor` element that is not below the module root; a missing Files layout, or a captured file that is not a `.go` path | `ErrInvalidRequest`; no judgment of the captured code |
| An import outside the environment, a dot or blank import, a test or ignored file, a file in a `testdata` or `vendor` directory or in one that starts with `_` or `.`, a package directory that is not a valid import path or is too long under the output root | Unsupported coverage at the import or file |
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
when the package entity does not carry exactly one name fact in the published
form, and neither reads the path in place of the name. The form is a provider
expression with the operator `go/package-name`, the encoding `go/utf8`, no
operand and text that is not empty and is UTF-8. A fact whose term differs in
one of these is refused as an absent fact is. The vocabulary publishes the
names as `PackageNameRelation()`, `PackageNameOperator()` and `TextEncoding()`,
and the bridge writes them from there. The readers do not apply Go's rule for
a package name to the text. In the base profile the identity is the constant
`captured`, and the name is the one the package clause declares.

**The limit of package identity.** A consumer that relocates declarations can
compare the owner of each declaration, `Type.Package().Path`, with
`CapturedImportPath` of the file it intended. That comparison establishes the
package of that file's directory, and that the checker put the declaration
into that package. Two files of one directory share the identity, so the
comparison does not say which file holds the declaration: that comes from the
declaration's site, `Type.Site()`, whose tree is the file's. It cannot certify
where the consumer writes the files.
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
identities `rust.ResolutionContract()` and `rust.AnalysisContractFor(settings)`,
the `Environment` type, and the graph names that the helper writes and Go reads
by name. The helper is Rust and spells those names itself, and
`TestRustHelperNamesMatchVocabulary` runs it and keeps the two sides equal.
`resolve/files/rust/provider/go.New` takes a
`resolve.Compiler` with absolute `HelperPath` and `SysrootPath`. Constructors
perform no IO, discovery or installation. The caller provisions that helper and
the matching compiler sysroot explicitly. Both operations require the settings
of one stated [environment](#rust-environment),
`rust.ConfigurationFor(environment)`. The helper checks a library crate against
the sysroot only and generates no code, under the `captured-types-2` helper
profile. Compiler version, helper profile and environment are receipt identity;
executable locations are configuration of the collaborator.

### Rust environment

**13 November 2026** ([#278](https://github.com/Bitspark/bittype/issues/278),
rows 28 and 35 of [research 0016](../research-docs/0016-native-files-placement.md)).
`rust.Environment` states what the captured crate is compiled under. A consumer
states every field, and no field has a default. Each field has a vocabulary,
and of that vocabulary this profile covers one value. A refusal comes before
the compiler runs, and no field is accepted and ignored. Later coverage admits
more values of the same fields and adds no constructor.

| Field | Its vocabulary. Another text is invalid input | This profile covers | Another value of the vocabulary |
| --- | --- | --- | --- |
| `Edition` | The four editions that the pinned compiler takes: `2015`, `2018`, `2021`, `2024` | `2024` (`rust.Edition2024`) | Unsupported coverage |
| `Crate.Kind` | The seven crate types that the pinned compiler takes: `lib`, `rlib`, `staticlib`, `dylib`, `cdylib`, `bin`, `proc-macro` | `lib` (`rust.KindLibrary`) | Unsupported coverage |
| `Crate.Name` | What the pinned compiler takes as a crate name: a text that is not empty, of characters that are alphanumeric or an underscore | An ASCII identifier that is no keyword of the edition and not the name of a crate that the sysroot loads | Unsupported coverage |
| `Target` | A well-formed target text | `x86_64-unknown-linux-gnu` (`rust.TargetLinuxX64`) | Unsupported coverage |
| `Dependencies` | Each name is judged as a crate name | None. The sysroot is the only other crate | Unsupported coverage |

Unsupported coverage is `rust.ErrUnsupported`: the value is well formed under
the field's vocabulary, this profile does not cover it, and a later profile may
admit it unchanged. Invalid input is `semantic.ErrInvalidRequest`: a text
outside the field's vocabulary, which includes an empty field, and settings
that `ConfigurationFor` returns for no environment. No refusal is both.
Validity is judged before coverage, over every field and every dependency, so
an environment with an invalid and an unsupported value is invalid.

**The vocabularies of an edition, a crate kind and a crate name are the pinned
compiler's,** read in the source of rustc 1.90.0 and put to that compiler by
`TestRustEnvironmentEqualsCompiler`, which compares its answers with
`ConfigurationFor`. Each is declared once. The editions, the crate kinds and
the supported values are in `semantics/rust/go/environment.go`. The rule for a
crate name and its table of characters are in `internal/rust/cratename.go`,
the lexical helper, because the Rust symbol package judges the crate of a
destination by the same rule and may not import the semantic vocabulary.
The `compiler` setting and the profile strings identify this vocabulary: the
editions, the crate kinds, the characters of a crate name and the keywords are
those of rustc 1.90.0. A new compiler pin is a new vocabulary, under its own
setting and profile strings, and the table of characters is regenerated for
it as the comment in `internal/rust/cratename.go` says.

- **An edition.** `compiler/rustc_span/src/edition.rs:114-126` names five. The
  stable compiler refuses the fifth, `future`
  (`compiler/rustc_session/src/config.rs:2127-2137`), so it is invalid here,
  like `2027` and `banana`.
- **A crate kind.** `compiler/rustc_session/src/config.rs:2953-2980` names
  eight. The stable compiler refuses the eighth, `sdylib`
  (`compiler/rustc_interface/src/passes.rs:271-272`). The compiler's option
  also takes a list such as `lib,rlib`. A list is no crate kind, and the field
  states one. A crate of several kinds needs a new field or key.
- **A crate name.** `compiler/rustc_session/src/output.rs:57-79`: the name is
  not empty, and each character satisfies `c.is_alphanumeric() || c == '_'`,
  with `char::is_alphanumeric` of `library/core/src/char/methods.rs:922-924`
  (`is_alphabetic` at 778-783, `is_numeric` at 1003-1008). So `x²` and `ⅷ` are
  crate names, and `a-b` and a name with a combining mark are not. The
  characters are the compiler's own table, 144435 of them: a program built by
  the pinned compiler printed it, and the test builds that program again and
  compares every code point. Go's `unicode` package is not used, because its
  tables follow the Go release.
- **A crate of the sysroot.** A captured crate may not be named `std`, `core`
  or `alloc`, compared exactly. The reason is two entities of one name in a
  graph: the graph names an entity by its crate and its path, a captured
  source reaches the standard library's entities in these three crates, and
  with the crate name `core` two entities are named `core` and a lookup of the
  crate root is refused as ambiguous. The sysroot loads sixteen further crates
  for the standard library, among them `libc`, `hashbrown`, `memchr` and
  `object`, and a captured crate may have any of their names. The evidence,
  from the independent review of 13 November 2026: under each of the nineteen
  names one source was analyzed, and only these three give two entities of
  one name; over 83 sources no graph held an entity of another crate of the
  sysroot. Its limit: that no path exposes an entity of a private dependency
  of the standard library is surveyed and not proven. The test asserts the
  three as unsupported, admits every other crate of the sysroot's library
  directory, and finds the crate root under each of the sixteen names.
- **A target** is the one field where invalid and unsupported cannot be told
  apart by what a compiler has. The compiler takes any text as the name of a
  target, and a text that ends in `.json` as the path of a file
  (`compiler/rustc_session/src/config.rs:2348-2359`). The vocabulary holds no
  list of targets. A well-formed target text is not empty, consists of ASCII
  letters, digits, `_`, `-` and `.`, and does not end in `.json`. That is the
  alphabet of the 290 targets of the pinned compiler, which the test checks.
  A well-formed text other than the supported target is unsupported, whether
  it names a target of some compiler, like `aarch64-unknown-linux-gnu`, or
  none, like `banana`. A path is invalid: an environment names no file.

**The crate name is the consumer's.** The compiler analyzes under
`Crate.Name`, so every compiler path of the model starts with it, as the
advertised public path of an exposed declaration does. The profile has no
crate name of its own.

**The compiler and the sysroot are not fields.** The compiler is pinned by the
profile and is a setting of its own. The sysroot is a path given to the
resolver. A receipt binds the compiler version that the helper was built with,
and no identity of the sysroot: two sysroots of that version with different
content give receipts with the same identity. This is a stated limit.

**The settings.** `ConfigurationFor` returns three settings of owner `rust`, in
this order. They are part of every receipt:

| Key | Bytes |
| --- | --- |
| `compiler` | `1.90.0` |
| `environment-1` | The environment as JSON, with every member, in the order of the fields. An empty dependency collection is `[]` |
| `helper` | `captured-types-2` |

```json
{"Edition":"2024","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[]}
```

`rust.EnvironmentOf(settings)` reads the environment back. It accepts exactly
the settings that `ConfigurationFor` returns for a supported environment. It
judges the form first: it decodes the value, builds the settings of what it
decoded and compares them, bytes for bytes. An unknown member, a repeated
member, another spelling, another key, one more setting and another order are
invalid input, whatever the values say. It then judges the values as
`ConfigurationFor` does. The resolver and the analyzer both read the
environment of their request this way, and analysis also requires its
configuration to equal its resolution's. `rust.EvidenceFor(settings)` and
`rust.AnalysisContractFor(settings)` exist only for such settings, so no
evidence can be required without an environment.

**The supported values exist in one place,** `semantics/rust/go/environment.go`.
The bridge passes the environment of the request to the helper, and the helper
passes it to the compiler and judges no value.

**An earlier configuration.** Until this change the configuration was fixed:
`compiler`, an `environment` setting with the text
`edition=2024;target=x86_64-unknown-linux-gnu;crate=captured;lib;sysroot-only;no-codegen`,
and `helper` with `captured-types-1`. Both providers refuse it as invalid
input, and no evidence requirement exists for it, so a receipt of the earlier
helper is not admitted as one of this helper and `rust.Open` refuses it.

The profile strings changed with the configuration, see
[the vocabulary](#rust-vocabulary).

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

**13 November 2026.** A path into the sysroot is bound segment by segment, also
when it is absolute. For `::std::string::String` the helper binds `std` to the
root module of the crate `std`, `string` to the module `alloc::string` and
`String` to the struct `alloc::string::String`. It names an entity by the crate
and the path of its definition, which is not the path that was written:
`::std::option::Option` binds to `core::option::Option`. An entity of the
sysroot has no declaration site. A primitive type such as `bool` is no path
segment and has no binding, so a declaration of that name that captures a use
is not visible in the graph.
[A test](../integration/go/rust_semantics_test.go) pins this, and
[Destinations](destinations.md#rust) states what it means for a recorded
destination.

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
Every captured item and field carries its declared visibility, every captured
item whether another crate can name it, and a declaration whose interface names
a more private type carries the compiler's observation.
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

### Rust vocabulary

The provider-owned `rust` vocabulary is:

| Graph element | Meaning |
| --- | --- |
| `module`, `alias`, `struct`, `enum`, `variant`, `field`, `trait`, `impl`, `parameter` entities | Native declaration identities, captured sites and owners |
| `Reference` / `Application` | Native types and effective arguments; shared `DirectParameterBindings` applies |
| `target` bindings | Compiler-resolved path segments at captured nodes |
| `declared-type` facts/occurrences | Checked type of each captured type declaration or field |
| `native-type`, `display` facts | Native type handles and diagnostic spelling; display is never identity |
| `fields`, `variants` facts | Ordered local ADT structure, substituted field types, field names and declaration identities. A field's visibility is the `declared-visibility` fact of its entity |
| `default`, `predicate`, `implements` facts | Generic defaults, declared trait predicates and checked impl headers |
| `declared-visibility` fact | On every captured module below the crate root, alias, struct, enum, trait and field, exactly one: the visibility as declared, in the compiler's resolved form. Its term is a `public` expression, or a `restricted` expression whose one operand refers to the `module` entity it is restricted to. The crate root is a module entity |
| `effective-visibility` fact | On every captured module below the crate root, alias, struct, enum and trait, exactly one: an `effective-visibility` expression whose data state the `Audience`, the `Path` and whether the audience can name the declaration through it, `Nameable` |
| `private-interface` fact | On a captured declaration whose interface names a declaration more private than itself, one for each such observation of the compiler: a reference to the more private declaration |
| `primitive`, `tuple`, `slice`, `array`, `pointer`, `reference` expressions | Rust-owned `json-1` metadata and ordered type operands |
| `public`, `restricted`, `effective-visibility` expressions | The forms of the two visibility facts |

**The profile has a new version since 13 November 2026:**
`rustc-1.90.0/captured-types-2/resolution-1` and
`rustc-1.90.0/captured-types-2/analysis-1`, where the strings said
`captured-types-1`. The change that added the three visibility facts also
removed the `Public` member from the data of the `fields` term. An analysis
profile is "A versioned contract defining obligations, domains, vocabulary and
requirements for a model" ([analysis](analysis.md#concepts)), and "Matching
profiles are explicit compatibility promises"
([resolution](resolution.md#capture-once-invoke-explicitly)). A reader written
for the earlier vocabulary cannot read a graph without that member, so the
promise of the earlier version does not hold and the version changed. A receipt
of the earlier version is refused twice, by its profile and by its
configuration. Analysis meets the profile first: the shared admission of its
resolution receipt compares the profile before the configuration. `rust.Open`
meets the configuration first: it reads the environment from the receipt's
settings before it admits the receipt.

Generic parameters retain their actual owning declaration and slot. Rust's
implicit trait `Self` is a distinct binder even though rustc internally reuses
the trait's `DefId`. Same-spelled `T` parameters in different declarations never
collapse. The native type cache includes the binder environment for parameterized
types. Recursive local types use entity/fact links, not cyclic term containment.
Aliases keep their declaration identities while their native type can normalize
to another declaration. Source path bindings still identify the alias. This slice
does not export a semantic occurrence for every nested source type expression.

### Rust visibility facts

**13 November 2026** ([#278](https://github.com/Bitspark/bittype/issues/278),
rows 23 and 24 of [research 0016](../research-docs/0016-native-files-placement.md)).
A Rust Files stage decides the visibility of every item it places. These facts
are what that decision is checked with. They are two facts and they are kept
apart, because an item's own `pub` does not establish that anyone can name it.

**The declared fact** is what the declaration states, resolved by the compiler:

| Written | Declared fact |
| --- | --- |
| `pub` | Public |
| `pub(crate)` | Restricted to the crate root |
| No modifier | Restricted to the enclosing module, which for a root item is the crate root |
| A field of an enum variant, which takes no modifier | The visibility of its enum |

`Declaration.Visibility()` reads it as `rust.Visibility{Public, Scope}`. `Scope`
is the module declaration the item is restricted to and is invalid for a public
one. A field has the same fact in the same form. `Field.Declaration` is the
field's declaration, and `Field.Public` no longer exists: a yes or no could not
tell `pub(crate)` from private.

**The effective fact is defined by an audience and a path.** It answers one
question: can the audience name the declaration through the path. In this
profile:

- **The audience** is `another-crate` (`rust.AudienceAnotherCrate`): code of a
  crate other than the captured one. This is the audience of
  [exposure in Rust](naming.md#exposure-and-its-audience): "Nameable from
  another crate through its advertised public path".
- **The path** is `definition` (`rust.PathDefinition`): the path of the
  definition, which is the crate name, the enclosing modules and the name.
- **`Nameable`** holds exactly when every segment of that path is visible to
  the audience: when the declaration and every module around it are declared
  public.

The profile admits no `use`, so no re-export exists and the definition path is
the one public path a declaration can have. The fact is still defined by
audience and path. A profile that admits re-exports adds a path, under its own
helper profile, and changes nothing about what this fact says.

The helper takes the answer from the compiler's own effective visibility, at
the level that means direct access. It is not what the compiler calls
reachable. In `mod m { pub struct D {} } pub struct A { pub d: m::D }` the
type `D` reaches other crates through the field and they cannot name it. `D`
is declared public, and its effective fact says that it is not nameable.

`Declaration.Effective()` reads it as `rust.Effective{Audience, Path, Nameable}`
and refuses an audience or a path that it does not know.

**A missing fact never means private.** Both readers fail with
`rust.ErrMissingFact` when the graph has no such fact for the declaration,
more than one, or one in another form. The crate root, variants, impls,
parameters and sysroot declarations carry neither fact, because Rust gives
none of them a visibility modifier of its own, and a reader that is asked for
one fails. Fields carry the declared fact only.

**The forms, exactly.** A declared term is `public` without an operand or
`restricted` with one operand, a module, and carries no data. The data of an
effective term is one object with the members `Audience`, `Path` and
`Nameable`, each once and spelled so, with a text, a text and a truth value,
and nothing after the object. The data of a `fields` term is a list with one
object for each field, each with exactly the members `Name` and `Entity`. A
repeated member, a member in another case, an unknown member and data after
the value are another form, and the reader fails. `Type.Fields()` fails with
`rust.ErrMissingFact` like the two visibility readers, also for a field whose
entity the graph does not have and for a fact that is present twice. It
answers "no fields", without an error, only for a type that has no `fields`
fact. The helper spells these member names itself, and the tests that run it
read every one of them through these readers.

**An interface that names a more private type** is the Rust-owned rule of row
24. `pub struct Public { pub field: Hidden } struct Hidden {}` is legal Rust,
and another crate cannot use the type of that field. The compiler reports it
with the lint `private_interfaces`, as a warning. The helper allows every other
warning and forces this one with `--force-warn`, which holds against
`-A warnings`. `-W` does not: the compiler then reports nothing. Denying the
lint would turn the case into a compiler error and into a rejected validity.
The helper receives its own diagnostics, withholds every warning from its
error stream and publishes each forced one as a `private-interface` fact: on
the declaration that the compiler reports it at, which here is the field, with
a reference to the more private declaration. The compiler reports it for a
field of a struct, of a tuple struct and of a variant, and for an alias. For
the default of a type parameter, `pub struct P<T = Hidden>`, it reports the
item itself, and the fact is on `P`. A test of the bridge package runs the
helper and reads its error stream: it is empty for a source with the lint and
holds the compiler's errors for a source that the compiler rejects.

- **The source stays analyzed.** Coverage, bindings and validity hold, and the
  graph is published. The profile keeps analyzing hand-written Rust that has
  such an interface.
- **`Declaration.PrivateInterface()`** returns the more private declarations,
  and none when the compiler observed none.
- **The refusal belongs to Rust Files.** It refuses an exposed declaration that
  has the fact, as unsupported coverage and never as invalid input, and it
  never exports the hidden type without authorization. That is slice 1 of
  [#159](https://github.com/Bitspark/bittype/issues/159) and not part of this
  provider.
- **A diagnostic the helper cannot place** on two captured declarations is
  unsupported coverage without a graph. Dropping it would publish a
  declaration as usable that the compiler reported.
- **The fact follows the declared visibility of the type, and that is its
  limit.** In `mod m { pub struct D {} } pub struct A { pub d: m::D }` the
  type `D` is declared public, so the compiler reports nothing and the field
  carries no fact. Another crate still cannot name `D`, and the effective fact
  of `D` says so. A Rust Files stage that decides whether a public record is
  usable therefore needs both: no `private-interface` fact on the
  declarations of the record, and the effective fact of every declaration in
  its public representation.
  [The naming design](naming.md#exposure-and-its-audience) states that rule:
  an unexposed source declaration in the body of an exposed one is refused
  "because a caller in another package or crate could not construct the
  value".
- **A limit.** The compiler reports a bound on a more private trait,
  `pub struct Boxed<T: Hidden>`, with another lint, `private_bounds`. The
  helper does not force it and publishes no fact for it.
- **A limit.** The fields reader does not require a field's entity to be of
  the field kind, because the vocabulary publishes no name for that kind.
- **A limit.** The helper does not report which facts it writes. The helper
  profile is a setting that the Go side states alone, so a helper binary of
  an earlier change is caught today only because the request gained a member
  that the earlier helper refuses. It goes with the next change of the helper,
  the binding of `bool` in
  [#276](https://github.com/Bitspark/bittype/issues/276).

`rust.Open(analyzed)` admits the exact completed analysis receipt of a
supported environment and reads the graph without a compiler invocation.
`Lookup` takes a compiler path such as `pages::models::Person` for a crate
named `pages`; paths select declarations but do not define equality.
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
captures, false references, settings of no environment and foreign receipts are
request errors. An environment outside the profile is refused as unsupported,
before the compiler runs. Helper/toolchain failures are operational errors. Cancellation stops the
child process. The bounds are 4,096 AST nodes and 1 MiB source, 16,000 native type
handles, 32,768 semantic records, export depth 128, a 30-second process deadline, 16 MiB response and
256 KiB stderr. Shared capture/publication budgets apply independently.

The [provider tests](../integration/go/rust_semantics_test.go) cover pagination,
nested substitution, defaults, binder ownership, recursion, module paths, marker
traits, invalid/unsupported/incomplete outcomes, source correspondence and receipt
isolation. They cover the declared and the effective fact of every kind of
item, the declared fact of fields, a graph with one fact removed, the
private-interface observation and the environment's refusals. The
[unit tests of the vocabulary](../semantics/rust/go/environment_test.go) pin the
bytes of the settings and read every malformed form of a fact. The [installed consumer](../integration/go/testdata/installed/rust_semantics_test.go)
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
