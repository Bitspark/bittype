# Template lifting and library-owned native type representations

**Owner direction, 3 November 2026. Status: documentation proposal only.**
This refines the [contract-language proposal](hydrated-types.md) and
[generation pipeline](generation.md). Ontos and Horos remain data-only. A
separate Bittype language provider describes Wire and mixed compositions;
Function, Cell and Stream remain ordinary library declarations.

The existing ASTTransformer interface and syntax engine need no new operation
for a named library contract. The contract provider, lifting transformer, native
transformers and representation configuration below are not implemented.
Examples describe the intended design, not executable configuration or new APIs.
No runtime adapter, validator, callable implementation or transport is added.

## 1. Three separate inputs

| Input | Owner and purpose |
| --- | --- |
| Canonical declaration modules | Libraries define types, imports, parameters and applications. Horos owns data declarations; contract modules use the separate Bittype language. |
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

## 2. Embedding and lifting are explicit, different operations

**Embedding** starts with a closed, resolved Horos data expression, specializes
it to a well-formed ground descriptor G, then produces `Data(G)`. This preserves
its ground domain. A closed `Page<Person>` can take this path. G contains no
generic parameter or unresolved application and is opaque to contract
substitution. The tooling retains the original declarations/applications as
provenance rather than treating specialization as a replacement identity.

**Template lifting** translates a selected Horos declaration into a separately
identified contract declaration, preserving its parameters. A data slot becomes
a contract-type slot only in this explicitly derived declaration. This does not
widen the source declaration or change what Horos accepts.

The consuming build selects embedding or the bounded structural lifting profile
below, and the declarations to transform. These choices configure a concrete
`Horos -> Contract` ASTTransformer. Route discovery does not choose a different
value domain for the caller. If both direct data generation and the lifted route
are installed, the caller explicitly selects the latter when it wants a lift;
a failed direct mapping never triggers that route as a fallback.

### Structural lifting rules

Let D be a resolved declaration in a checked Horos module bundle. Let L(E)
translate an expression in D's lexical environment. `C(...)` and `Data(G)` use
the [exact proposed contract vocabulary](hydrated-types.md#2-exact-proposed-wire-expression-and-ground-embedding).
The following rules apply in order, so a closed application or structural node
is not opportunistically collapsed into Data:

| Horos input | Contract output and obligation |
| --- | --- |
| Parameter at slot i | `C("param", I(i))`, bound to the corresponding parameter of the derived declaration; retain the source binder as provenance. |
| Application of resolved declaration D to arguments A | Application of the derived declaration L(D) to L(A), preserving argument order, arity, lexical scope and the original application. Follow the exact pinned target, never a same-spelled name. |
| Tuple | `C("tuple", L(E)...)`; retain the ordered elements, including an empty tuple. |
| Labelled compound | `C("compound", A(label), L(E)...)`; preserve the label as literal data. |
| Typed list | `C("list-of", L(element))`. |
| Closed record | `C("record", A("closed"), fields...)`; preserve each literal field name, required/optional presence and translated field type. Keys remain ground canonical text. |
| Union with at least two alternatives | `C("union", L(E)...)`; preserve the alternatives and their order. Do not select a behavioral binding by inspecting the branches. |
| Union with one alternative | L(the alternative), retaining the source union as provenance; the contract union constructor requires at least two alternatives. |
| Empty union | `Data(H("union"))`, preserving the empty ground type without introducing an upper bottom constructor. |
| Any other subtree that can independently specialize to a closed, well-formed Horos descriptor G | `Data(G)`. Retain source references in provenance; do not inspect literal contents or widen G's domain. |
| Any other subtree | Reject with a located unsupported-lift diagnostic. Do not erase constraints, substitute a top type or invent an upper constructor. |

The structural rules apply even when their children happen to be closed. This
makes the chosen template translation stable before arguments are substituted.
Memoize declaration translation by resolved source identity and the selected
lifting configuration, not by expanded shape. The first profile follows the
existing acyclic declaration/application restriction; it does not introduce
recursive named contract applications.

The closed-subtree rule is deliberately narrower than merely seeing no generic
parameter token. The subtree must have no free generic or recursive variable,
and all applications it contains must resolve and specialize within Horos.
For example:

- `any`, scalar types, byte/text refinements and exact literals remain ground
  inside Data. `Data(any)` never accepts a Wire.
- A closed complement, set, map, open record or guarded recursive descriptor
  may stay inside Data with its existing meaning. No live keys, set elements,
  open-record extras or recursion are introduced.
- `not(T)`, `set-of(T)` and `rec(...T...)` have no lift in this profile.
  A closed data specialization may instead be explicitly embedded, but a failed
  template lift cannot silently change to that operation.
- An exact literal that resembles `param`, `apply` or `wire` is still an opaque
  literal. Substitution and lifting visit declared type positions only.

### Derived modules and identity

The selected roots determine a transitive declaration closure. Translate each
reachable declaration once, grouped by its source module; imports in the output
refer to the corresponding derived contract modules. Preserve parameter order,
names and source declaration/application origins. Reject missing imports, cycles
or output naming conflicts before publishing any usable result.

Validate all original module imports under Horos's rules before selecting that
output closure, including imports unused by the selected roots. Pruning generated
output is not permission to skip the source language's resolution obligations.

The build supplies output module names and records its selected roots, source
module pins, lifting profile/version and target naming choices. Derived modules
use the contract module envelope and its canonical encoding/import pins from
[the contract proposal](hydrated-types.md#3-structural-composition-and-generic-declarations).
Dependency modules are encoded and pinned before their importers. The output
module pin plus declaration name identifies the derived declaration; provenance
records the distinct original Horos identity. A pin is not a behavioral
fingerprint or a claim of protocol equivalence.

No second Page is authored and no `LiftedPage` spelling is required. The compiler
may name the derived declaration Page in its selected output module. Its semantic
identity remains distinct from the original data declaration even when native
representations can be shared deliberately.

For admitted templates and data arguments, the lift must preserve the original
ground layout and membership after binding. Replacing an argument with a
capability type is meaningful only in the derived contract declaration. This is
a qualification obligation of this specific lifting profile, not a theorem
about every Horos constructor or every native representation.

## 3. Pagination through the type pipeline

These equations are explanatory notation, not a readable source grammar:

```text
Horos source:
  Page<T> = record closed {
    items: list-of(T)
    next?: scalar(utf8-text)
  }

Derived contract declaration:
  Page<T> = record closed {
    items: list-of(T)
    next?: Data(scalar(utf8-text))
  }
```

Horos's `Page<Person>` stays a data application. The derived Page can be applied
to an explicitly embedded Person data type or to `Cell<Data(PersonDescriptor)>`.
Writing `Page<Cell<Person>>` below is shorthand for those resolved upper-layer
references, not permission to pass a Cell into a Horos data slot.

```text
AST[Horos] ---------------------------------> AST[Native] -> CST -> source
    |
    +-- explicit selected lift -> AST[Contract] --+
                                                +-> AST[Native] -> CST -> source
```

Both native paths retain generic declarations and applications. Generate Page
once per selected module/representation, not a new record per application.
Specialization for a ground descriptor is a separate operation; it is not a
prerequisite for native generic generation. Representative target shapes are:

| Target ID | Page's type declaration | Application containing a Cell |
| --- | --- | --- |
| `typescript` | `interface Page<T>` with `items: T[]` and `next?: string` | `Page<Cell<Person>>` |
| `go` | `type Page[T any] struct` with `Items []T` and an explicit presence representation for Next | `Page[cell.Cell[Person]]` |
| `rust` | `struct Page<T>` with `items: Vec<T>` and `next: Option<String>` | `Page<cell::Cell<Person>>` |
| `cpp` | `template<typename T> struct Page` with `std::vector<T>` and `std::optional<std::string>` | `Page<cell::Cell<Person>>` |

These are candidate representations and illustrative symbols, not delivered
packages. Each native profile specifies container, presence, null, numeric and
ownership choices. A Go zero value alone cannot distinguish absence from a
present zero value. Optional text in this example does not admit null.
Mappings must report refinements/invariants that native types cannot express;
a strict unsupported case is an error, not silent semantic equivalence.

A library may explicitly nominate the same native generic Page for both the
Horos and derived contract views when their selected representations agree.
That association names both source identities and the qualified representation;
it is not inferred from equal-looking ASTs. Preserve separate semantic metadata
and later validation/binding choices. Native source alone does not recover them.

The generator of Page only traverses a parameterized record/list. The application
site supplies the native representation of T. For example, a selected
TypeScript representation can produce:

```typescript
export interface Page<T> {
  items: T[];
  next?: string;
}

export type PeoplePage = Page<Person>;
export type EditablePeoplePage = Page<Cell<Person>>;
```

The imports and Person/Cell declarations are deliberately omitted from this
small illustration; it is not a standalone compiler fixture. There is no Cell
case in the pagination transformer.

## 4. Ordinary libraries and optional native representations

The [complete proposed module equations](hydrated-types.md#4-library-owned-function-cell-and-stream)
declare Function over Wire, Cell as get/set functions and a pull Stream as
next/cancel functions with an ordinary Step union. Their readable forms are:

```text
Unit          = Data(empty tuple)
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
| Arguments | Explicit mapping from resolved source parameter slots to target type operands. Native constraints, extra parameters or indirections must be declared by the profile. |
| Imports and dependencies | Target import/header/crate identity and exported symbol path; exact package artifacts supplied by the consuming dependency lock. |
| Qualification | Supported applications/constraints and native representation limits. An unsupported argument fails explicitly. |

An external-type recipe emits a reference and required imports; it does not
redefine the external type. A declaration recipe emits owned type/interface AST
nodes when building the library's type artifact. Consumers reference that
artifact. A structural mapping recursively uses the selected representations of
referenced declarations and arguments. All three preserve origins and generic
bindings, and return no usable partial output on error.

The initial selection unit is an open declaration, not a runtime application
value or shape. For an application `Function<A,B>`, look up the resolved
Function declaration once, translate A and B, then instantiate the selected
native recipe in its defined parameter order. Aliases retain their own origins
and the resolved target; the resolver's explicit alias relation authorizes
following a target recipe. An independently declared lookalike does not.

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

1. Supply pinned declarations and any explicitly selected lift results to the
   appropriate resolver/provider. Resolve imports, slots and applications.
2. Construct one concrete native transformer per requested language pair, with
   its target profile, dependency lock and selected library recipes. Capture
   immutable configuration; there are no mutable process-global registrations.
3. Register those transformers with the existing engine and select the intended
   route. Format their target ASTs and print the resulting CSTs normally.

For example, in deliberately schematic Go notation:

```go
tsMapping := contractToTypeScript.New(tsProfile, modules, libraries.TypeScriptTypes())
goMapping := contractToGo.New(goProfile, modules, libraries.GoTypes())
rustMapping := contractToRust.New(rustProfile, modules, libraries.RustTypes())
cppMapping := contractToCPP.New(cppProfile, modules, libraries.CPPTypes())

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
| `transforms/horos-to-contract/go` | Explicit closed embedding and the separately selected partial lifting profile; resolved origins and derived modules. |
| `transforms/horos-to-<native>/go` | Direct data type generation, independently useful without capabilities. |
| `transforms/contract-to-<native>/go` | Structural type lowering and uniformly configured native recipes; no named-contract imports. |
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

- Data-only and capability-bearing pagination use one authored template, with
  distinct original/derived identities and stable open generics in all targets.
- Lifting preserves supported ground membership/layout, lexical slot scope,
  nested/imported applications, source origins, presence and literal opacity.
  Unsupported parameterized forms fail; closed forms stay explicitly ground.
- Root/dependency selection and canonical encoding give reproducible derived
  modules. Changing the source pin/profile cannot silently reuse old output.
- An unfamiliar library contract uses ordinary declarations and optional
  recipes without a core source, dependency, registry or bootstrap change.
- Recipe selection survives aliases/applications, distinguishes unrelated
  lookalikes, and rejects missing/conflicting profiles and target dependencies.
- One generic Page composes with data, Function, Cell and Stream target types.
  A library-authorized shared native representation retains separate source
  identities; native erasure is not evidence of semantic equivalence.
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
