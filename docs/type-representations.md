# Library-owned native type representations

**Owner direction updated 5 November 2026. Status: direct declaration frontends and pagination transformers implemented; custom native representations proposed.**
This refines the [unified Horos tooling design](hydrated-types.md) and
[generation pipeline](generation.md). Ground values/checks remain data-only;
Ontos/wire owns mixed values and the one Horos language describes their types.
Bittype consumes that model; Function, Cell and Stream remain ordinary libraries.

The existing ASTTransformer interface and syntax engine need no new operation
for a named library contract. The unified Horos frontend is [implemented](horos.md).
The [pagination transformers](semantic-generation.md) are implemented; custom
representation configuration below remains proposed. The
[type-owned generator plugin design](type-generation-plugins.md) is the current
proposal for Go contracts, builtin mappings, explicit overrides and composition.
It supersedes this page's earlier open-declaration-only selection scope. Recipe
examples below explain representations, not executable APIs.
No runtime adapter, validator, callable implementation or transport is added.

The separate [runtime value explanation](runtime-values.md) connects this static
design to existing Wire/hydration. Native type recipes describe declarations;
future value adapters convert actual hydrated instances to and from native T.
That page also separates compiled typed adapters from dynamic descriptions and
marks the runtime-binding API as unimplemented.

## 1. Three separate inputs

| Input | Owner and purpose |
| --- | --- |
| Canonical declaration modules | Libraries define kinds, types, imports, parameters and applications. Horos owns the general kind contracts and one declaration language. |
| Native representation choices | A library supplies optional target recipes for its declarations; the consuming build selects exact recipes and dependencies. Ordinary structural declarations need no custom recipe. |
| Runtime bindings | Later library work supplies behavior and conversion where needed. Type generation neither implements nor loads it. |

Loading a module does not install a runtime adapter. Selecting a native type
does not prove its behavior. A native function-signature AST node describes
language syntax; it does not select a Function protocol.

The consuming application supplies an immutable, pinned module bundle and
target configuration. Dependency retrieval and package installation happen
outside transformation; no resolver or transformer performs implicit network IO.
Providers establish the semantic prerequisites using the existing
[capture, resolution and analysis boundaries](resolution.md). A syntactically
admitted AST is insufficient evidence for a semantic transformation.

## 2. Direct declarations and generic bounds

Author one declaration with general kind contracts. The ordinary Type prelude
definition denotes proper-sort top; ordinary Data recursively describes ground
type forms. `Page<T>` and `DataPage<D: h.Data> = Page<D>`, using an explicit
prelude import, express the distinction without separate
languages, copied templates or a schema-level Data wrapper. Ground descriptors
occur directly in expressions. Ontos value embedding/recovery is unchanged.

Template lifting, the second frontend and closed schema embedding are removed
without aliases. Horos Apply binds checked Program-owned values; Descriptor
materializes a proper result, while GroundDescriptor also requires the ordinary
ground-kind obligation and existing WF. Neither erases the authored
open declaration. Canonical expressions, bounds, pins and lexical binding come
from [Horos](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md).
AST addresses and origins record correspondence, not semantic equivalence or
behavioral authority. Native names and recipe configuration remain separate.

## 3. Pagination through the type pipeline

Author this generic once in Horos. These equations are explanatory notation,
not a readable source grammar:

```text
Page<T> = record closed {
  items: list-of(T)
  next?: scalar(utf8-text)
}
DataPage<D: h.Data> = Page<D>
Person = record closed { name: scalar(utf8-text) }
PeoplePage = Page<Person>
EditablePeoplePage = Page<Cell<Person>>
```

Both applications select the same Page. Cell is an ordinary pinned library
declaration. No expanded Page record or second Page module is needed. DataPage imposes an explicit data-only constraint on the same Page. The future route is:

```text
AST[Horos] -> AST[Native] -> CST[Native] -> source
```

The route consumes checked source semantics, including generic bounds. Native generation
retains generic declarations and applications; specialization is not a required
prepass. Generate Page once per selected module/representation, not a record per
application. Representative target shapes are:

| Target ID | Page's type declaration | Application containing a Cell |
| --- | --- | --- |
| `typescript` | `interface Page<T>` with `items: T[]` and `next?: string` | `Page<Cell<Person>>` |
| `go` | `type Page[T any] struct` with `Items []T` and an explicit presence representation for Next | `Page[cell.Cell[Person]]` |
| `rust` | `struct Page<T>` with `items: Vec<T>` and `next: Option<String>` | `Page<cell::Cell<Person>>` |
| `cpp` | `template<typename T> struct Page` with `std::vector<T>` and `std::optional<std::string>` | `Page<cell::Cell<Person>>` |

These are candidate representations and illustrative symbols, not delivered
packages. Each native profile specifies container, presence, null, numeric and
ownership choices. A Go zero value alone cannot distinguish absence from a
present zero value. Optional text here does not admit null. Mappings must report
refinements or invariants native types cannot express; unsupported strict mappings
fail instead of claiming semantic equivalence.

Page's mapping traverses a parameterized record/list. Its application supplies
the representation of T. A selected TypeScript representation could produce:

```typescript
export interface Page<T> {
  items: T[];
  next?: string;
}

export type PeoplePage = Page<Person>;
export type EditablePeoplePage = Page<Cell<Person>>;
```

Imports and Person/Cell declarations are omitted from this illustration. There
is no Cell case in the pagination transformer. Native source does not by itself
recover canonical identities, constraints or later validation/binding choices.

## 4. Ordinary libraries and optional native representations

The [complete proposed module equations](hydrated-types.md#4-library-owned-function-cell-and-stream)
declare Function over Wire, Cell as get/set functions and a pull Stream as
next/cancel functions with an ordinary Step union. Their readable forms are:

```text
Unit          = Data(H("tuple"))
Function<A,B> = Wire<tuple(A, Wire<B>)>
Cell<T>       = record { get: Function<Unit,T>, set: Function<T,Unit> }
Step<T>       = compound("item", T) | compound("end")
Stream<T>     = record {
  next: Function<Unit,Step<T>>
  cancel: Function<Unit,Unit>
}
```

The Function shape is illustrative and success-only. Reply discipline, failures,
timeouts, resource lifetime, stream demand/cancellation and cell mutation are
library behavior. None follows from these declarations or native compiler
acceptance. Unit is an empty ground tuple; any mapping to native void must be
explicitly selected and qualified, never inferred from a name.

When the selected target profile supports the underlying forms, ordinary
structural generation suffices for these declarations. For example, generated
Cell fields refer to the selected Function type. This needs no Cell-specific
generator extension. Step needs a qualified union representation; a target with
no exact native sum representation must document its retained invariants or
reject that requested profile.

A library may instead supply a custom native representation, such as its own
callable handle, iterator-like type or opaque cell handle. This is an optional
recipe attached to the library's exact declaration. The library owns the
association and its target support; the consuming build selects it explicitly.
The generator does not infer one from a name, record shape or native signature.

### Recipe contract

This is the proposed information a recipe carries, not a serialized schema or
an implemented Go interface:

| Field | Required meaning |
| --- | --- |
| Source | Exact source language/profile, module pin and declaration identity; ordered canonical type parameters. |
| Target | Target language ID, syntax/profile compatibility and the selected native representation profile/version. These are distinct from source identity. |
| Native type | Either an external named generic type or a type-only target AST recipe. No raw source string or executable body. |
| Arguments | Explicit mapping from resolved source parameter slots to target type operands. Source kind contracts and their native enforcement limits must be explicit. Native constraints, extra parameters or indirections must be declared by the profile. |
| Imports and dependencies | Target import/header/crate identity and exported symbol path; exact package artifacts supplied by the consuming dependency lock. |
| Qualification | Supported applications/constraints and native representation limits. An unsupported argument fails explicitly. |

An external-type recipe emits a reference and required imports; it does not
redefine the external type. A declaration recipe emits owned type/interface AST
nodes when building the library's type artifact. Consumers reference that
artifact. A structural mapping recursively uses the selected representations of
referenced declarations and arguments. All three preserve origins and generic
bindings, and return no usable partial output on error.

An ordinary library representation selects an exact declaration, including an
open declaration, rather than a runtime application value or shape. The
[plugin proposal](type-generation-plugins.md#exact-selection-and-explicit-overrides)
also supports exact source-owned builtin selectors through the same mechanism.
For an application `Function<A,B>`, look up the resolved
Function declaration once, translate A and B, then instantiate the selected
native recipe in its defined parameter order. Aliases retain their own origins
and resolved targets. Following a target recipe uses provider-established
reference/application facts; for Horos, an alias-like declaration has a resolved
application body, not a new general alias-relation fact. An independently
declared lookalike does not acquire that target's recipe.

Missing required custom recipes, unsupported targets/arguments, incompatible
profiles and conflicting recipes for one source/target selection are errors.
Do not use registration order, last-writer-wins, shape matching or fallback to a
structural Wire representation. Structural generation is a deliberate choice
in the selected profile, including for declarations without a custom recipe.

## 5. Configuration for TypeScript, Go, Rust and C++

The following table describes one illustrative external Function type recipe
for each target. `FunctionDecl` means the already-resolved exact declaration,
including its module pin and source profile. It is not a lookup by the text
"Function". Package/header names below are examples, not existing dependencies.

| Target | Library-owned recipe configuration | Emitted type use |
| --- | --- | --- |
| TypeScript (`typescript`) | Type-only module import `@example/function-types`, exported symbol Function, arguments slots 0 then 1; selected TypeScript syntax and module-resolution profiles. | `Function<A, B>` through a generated type-only import/alias. |
| Go (`go`) | Import path `example.org/contracts/function/types`, exported identifier Function, arguments slots 0 then 1; selected Go profile and generic constraints. The import alias is allocated without collisions. | `fn.Function[A, B]` when fn is the allocated alias. |
| Rust (`rust`) | Locked crate package and local crate name `example_function_types`, exported path Function, arguments slots 0 then 1; edition, ownership/lifetime and trait-bound choices explicitly supplied by the profile. | `example_function_types::Function<A, B>`. |
| C++ (`cpp`) | Locked library artifact, header `example/function_types.hpp`, qualified template `example::function::Function`, arguments slots 0 then 1; selected C++ standard and ownership/completeness constraints. | `example::function::Function<A, B>`. |

One recipe does not automatically establish support in another language. A
library can supply metadata for external symbols without implementing its
runtime in Go. Any executable generation helper is implemented in Go, like the
rest of the toolchain, and constructs the selected target AST. Native execution
and behavior, when later required, need a separately available target library.

For the default structural Cell and Stream, there is no additional native
recipe to install: their generated fields compose the configured Function
representation and ordinary type mappings. If Cell or Stream selects a custom
surface, its own library supplies the same fields in a recipe for each supported
target. Neither contract receives a privileged registration route.

The recipe's logical symbol and dependency are stable inputs; local import
aliases and qualification are allocated by the native transformer. Conflicting
symbols/dependencies must be diagnosed or resolved by explicit configuration.
The formatter handles spelling, escaping and layout after these choices; the
printer only emits CST leaves. Missing AST vocabulary is an unsupported-target
error, not permission to print source and parse it back.

### Composition with the existing engine

The consuming build performs three operations:

1. Supply directly authored, pinned declarations to the
   appropriate resolver/provider. Resolve imports, slots and applications.
2. Construct one concrete native transformer per requested language pair, with
   its target profile, dependency lock and selected library recipes. Capture
   immutable configuration; there are no mutable process-global registrations.
3. Register those transformers with the existing engine and select the intended
   route. Format their target ASTs and print the resulting CSTs normally.

For example, in deliberately schematic Go notation:

```go
tsMapping := horosToTypeScript.New(tsProfile, modules, libraries.TypeScriptTypes())
goMapping := horosToGo.New(goProfile, modules, libraries.GoTypes())
rustMapping := horosToRust.New(rustProfile, modules, libraries.RustTypes())
cppMapping := horosToCPP.New(cppProfile, modules, libraries.CPPTypes())

engine.RegisterTransformer(tsMapping)
engine.RegisterTransformer(goMapping)
engine.RegisterTransformer(rustMapping)
engine.RegisterTransformer(cppMapping)
```

These constructors/methods are proposed composition notation, not current APIs;
real composition must handle construction and registration errors. The selected
library set is assembled by the application, not a built-in library catalogue.
Concrete transformers consume uniform recipe contracts and must not import
specific Function/Cell/Stream libraries. Application composition depends on the
selected library recipes; recipes depend on target AST/configuration APIs.

The implemented [ASTTransformer](../contracts/transform/go/transformer.go)
continues to expose FromID, ToID and Transform. Its configuration is captured
by the concrete implementation. No contract-specific method or callback is
added to the engine. There is one edge per language pair, not one edge per
library; alternative profiles use separate configured engine instances under
the current [registry policy](engine.md). Language-provider registration and
module resolution are also distinct from native recipe selection.

Native profile/recipe metadata is generation configuration, not a new Horos
constructor or a runtime registry. No init-time discovery, built-in default
contract set, dynamic import by contract name or live endpoint acquisition
occurs during static tooling.

## 6. Ownership, layout and later runtime work

The existing [stage-first layout](../LAYOUT.md) remains in force:

| Proposed location/owner | Responsibility |
| --- | --- |
| `transform/syntax/shared/generation/go` | Kernel registration and composition of supplied providers for explicit roots. |
| Bittype Horos integration packages | Builtin target mappings in `transform/syntax/horos-to-<target>/mapping/go`, separate from the shared composer. |
| Each independent contract library | Canonical module, optional target recipes/type artifacts and future runtime bindings. No per-library stage is added to bittype. |
| Consuming application composition | Pin modules/dependencies, select profiles and recipes, register transformers, arrange output artifacts. |

Here `<native>` is `typescript`, `golang`, `rust` or `cpp`; Go's language ID
is `go`. No placeholder packages or registration SPI are created by this
document. Future provider/transformer dependency allowances require review and
tests; current generic and native-stage import boundaries stay unchanged.

Later runtime adapters can compose through generic parameters as Nightseam's
[generated Page adapter](https://github.com/Bitspark/nightseam/blob/85a238979129a8f361a9c9c05ea38ecd0343a569/cmd/nightseam/testdata/golden-families/api/ts/boxes-client/src/types.ts#L180)
did: a container receives the element's interpretation. That inspected source
also had a built-in callable kind. It is evidence for parameter composition,
not a specification for our kernel, type language or runtime.

Types-only generation emits no adapter factories, sends, exports, proxy
creation or behavior. Later generated consuming bindings can import selected
library adapters statically. Function owns invocation; Stream owns its stream
laws; Cell owns mutation. Wire/hydration retain their existing responsibilities
and never discover these libraries by name.

## 7. Qualification required before implementation is claimed

Future fixtures must establish:

- Data and capability applications use the same authored Page, preserving
  its declaration identity, lexical parameter scope and open generics in all targets.
- Native mappings preserve general kind obligations, imported applications, source
  origins, presence and literal opacity. Native erasure never proves a Data
  bound or ground equivalence.
- The complete pinned module context is checked. Missing imports, inconsistent
  profiles and unsupported target representations produce explicit diagnostics.

- An unfamiliar library contract uses ordinary declarations and optional
  recipes without a core source, dependency, registry or bootstrap change.
- Recipe selection survives aliases/applications, distinguishes unrelated
  lookalikes, and rejects missing/conflicting profiles and target dependencies.
- One generic Page composes with data, Function, Cell and Stream target types.
  Profile-qualified source identity remains explicit even when native types erase it.
- Target AST checks, formatter/printer composition and native compiler tests
  qualify TypeScript, Go, Rust and C++ against the selected external type
  artifacts. Native precision limits are explicit. No runtime behavior is
  claimed from those tests.
- Architecture checks prove no named-contract knowledge in generic trees,
  semantics infrastructure, engine, Wire or hydration, and no runtime loading
  or implicit IO during generation.

This documentation task adds no such implementation or conformance tests. It
forbids silent domain widening, raw-source recipes, shape-selected behavior,
central Function/Cell/Stream switches or defaults, and any compatibility route
that retains those forbidden special cases.
