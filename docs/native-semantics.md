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
  -> semantics/golang/go (typed read-only view)
  -> explicit native query session, when a consumer needs admission checking
```

Resolution and type checking are coupled inside `go/types`. The resolver declares
both responsibilities; analysis consumes its exact compatible receipt without
running the compiler again. Syntax admission still permits unresolved names.
There is no language switch in the engine and no automatic registration.
Use resolution.NewResolver and analysis.NewAnalyzer with their respective Go contracts and providers.

The package layout is:

| Package | Responsibility |
| --- | --- |
| `semantics/golang/go` | Go-owned vocabulary, environment identity and immutable typed graph readers |
| `internal/golangsem` | Direct AST bridge, bounded compiler invocation and fact export |
| `resolve/files/golang/provider/go` | Configured resolution, coverage and native validity reports |
| `analyze/files/golang/provider/go` | Receipt-consuming analysis and explicit compiler-backed admission queries |
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
references, foreign receipts and cancellation remain operation errors. Native
checking is synchronous: cancellation is checked around compiler calls and does
not claim to preempt `go/types` while it runs.

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

`golang.ImportConfiguration(golang.ImportEnvironment{Module, Packages})` extends
the base configuration with one normalized `imports-1` setting. It follows Go's
own module layout:

- **Captured packages.** Each captured tree must be linked to a `.go` file in the
  request's Files context. Its package import path is the module path joined with
  the file's directory: `pages/page.go` under `example.org/app` belongs to
  `example.org/app/pages`. Files in one directory form one package.
- **Supplied packages.** Every other imported package is supplied as Go source:
  an import path and up to 16 files. Supplied text participates in receipt
  identity and is never an output. A supplied package may import only other
  supplied packages, and it may not use a path under the module.
- **No discovery.** Import paths need a first element containing a dot, Go's
  module-path rule. That also keeps supplied packages from naming standard
  library paths, which this profile does not provide. Nothing is read from a
  module cache, `GOPATH` or the network.

Packages are checked in import order through a closed importer. The results
are classified as follows:

| Case | Result |
| --- | --- |
| A captured type error, an unused import, an unexported cross-package use, an import cycle through a captured package, inconsistent package clauses | Rejected validity, with a diagnostic at the captured node |
| An import outside the environment, a dot or blank import, a package directory that is not a valid import path | Unsupported coverage at the import or file |
| A supplied package that fails to parse or type-check, imports outside the environment; a missing Files layout or a captured file that is not a `.go` path | `ErrInvalidRequest`; no judgment of the captured code |

The graph gains:

- one `package` entity per captured package, in import-path order, owning that
  package's declarations, followed by one per supplied package;
- an `import` binding from each captured import declaration to its package's
  entity;
- `name` bindings from qualified type references to the imported declarations.

Supplied declarations are exported by identity only: their reference,
application and display terms, never `underlying`, `alias-target` or
`constraint` facts. Their internals may lie outside this profile.
`Model.LookupIn(importPath, name)` reads a package-qualified declaration.
`Model.Lookup` refuses a name that more than one package declares. Use
`analyze.ContractFor(settings)` and `golang.EvidenceFor(settings)` with the
same configuration for resolution and analysis.

A consumer explicitly constructs both adapters, resolves with
`semantic.Options{Configuration: golang.Configuration()}`, and calls the analyzer
with that same configuration and `Resolution: &resolved`. `golang.Open(analyzed)`
then opens the typed view. A resolution receipt alone, a receipt from a different
capture or a different compiler configuration cannot substitute for that input.

`semantics/golang.Open` admits the analysis receipt and reads this graph without
compilation. Native `NewQueries` is a separate, explicit operation: it rebuilds
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
  -> semantics/typescript/go (typed read-only view)
```

The package roles mirror Go. `internal/typescriptsem` owns the compiler bridge,
including its embedded CommonJS helper. Each TypeScript stage supplies its own provider. No TypeScript types, rules or compiler handles enter the
agnostic kernel; no native stage imports Horos, Logos or generation packages.

The resolver constructor requires a `resolve.Compiler` containing absolute
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
and their declaration text participate in receipt identity. Use `analyze.ContractFor(settings)`
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
immutable graph without invoking the compiler. `analyze.NewQueries(model,
compiler)` separately configures explicit native queries. `Assignable(ctx,
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
Its Cargo package is unpublished. `resolve/files/rust/provider/go.New` takes a
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
