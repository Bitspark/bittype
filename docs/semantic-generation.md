# Semantic native type generation

**Owner direction, 4 November 2026.** Generate native
types and interfaces from checked Horos declarations. This work adds explicit
semantic transformation input and pagination mappings for TypeScript, Go, Rust
and C++. It adds no runtime adapters, validators, serialization or behavior.

## Composition

Capture an admitted source snapshot once. Resolve it explicitly and analyze it
with that resolution receipt. Pass the same snapshot and explicit analysis
result to each selected transformer. The transformer requires the matching
operation, profile, configuration, completed obligations and available model.
Equal tree addresses or bytes from another capture are not compatible evidence.

ASTTransformer retains FromID and ToID. Transform receives context and an input
containing the captured source and ordered evidence. Its typed adapter admits
the evidence against the transformer's declared requirements before calling the
mapping. Receipt admission shares the existing semantic dependency rules.
There is no hidden resolver or analyzer invocation. Intermediate route outputs
receive new snapshots; old receipts are never carried across a transformation.

```text
AST[Horos] + explicit resolution/analysis evidence
  -> AST[TypeScript | Go | Rust | C++]
  -> existing target formatter -> CST -> existing printer -> source
```

The semantic graph supplies declaration and parameter identities, Reference,
Application and ProviderExpression terms, positional argument bindings, bounds
and dependencies. The source AST supplies current origins and checked module
pins. Model-local entity IDs are not persistent module/declaration identities.
No expanded generic interface, second canonical algebra or shape-based contract
lookup is introduced.

## First mapping profile

For example, this readable source encodes ordinary Horos revision-2 declarations:

```text
module example.pagination;
type Page<T> = { items: list<T>; next?: text; };
type Person = { name: text; active: bool; };
type PeoplePage = Page<Person>;
type NestedPage = Page<Page<Person>>;
type DataPage<D: Data> = Page<D>;
```

Resolution gives `Page.T` and `DataPage.D` separate binder identities and each
application an exact target and positional arguments. Generation reads these
facts, emits one open generic `Page`, and retains `PeoplePage` as an application
alias. It does not specialize or copy the `Page` body for each consumer.

The first profile covers named/imported declarations, generic binders and
applications, closed records, UTF-8 text, booleans, lists and optional fields.
Canonical module identity and target package/file placement remain separate.
Mappings retain declaration/binder links and current source origins. A target
configuration supplies an output tree ID, a Go package name or a C++ namespace.
This first profile explicitly bundles the complete checked module closure into
one source unit. Root declarations use exported native names; imported
declarations use `M<first-16-hex-of-SHA256(pin-text)>_<Name>`. Collisions fail;
the hash is only a spelling, never semantic identity. Per-module output/import
placement is deferred. Generation has no implicit network or filesystem IO.
The [Syntax and Files proposal](source.md) specifies a multi-unit syntax result
and symmetric Files parsing, formatting and printing. Source exposes established
Semantics alongside concrete file organization. The current single-unit mapping
does not establish native Semantics or implement those proposed Files stages.

| Horos form | TypeScript | Go | Rust | C++ |
| --- | --- | --- | --- | --- |
| UTF-8 text | string | string | String | std::string |
| Boolean | boolean | bool | bool | bool |
| list-of(T) | T[] | []T | Vec<T> | std::vector<T> |
| Named closed record | interface | struct | struct | struct |
| Optional field | field?: T | explicit BittypeOptional[T] | Option<T> | std::optional<T> |
| Type parameter | native parameter | native parameter | native parameter | template parameter |
| Application | F<T> | F[T] | F<T> | F<T> |

TypeScript requires strictNullChecks and exactOptionalPropertyTypes. Go's
presence helper is a generated type declaration with Present and Value fields;
the value is ignored when absent. Its public name is `BittypeOptional`. No helper methods or runtime implementation
are generated. Absence is never silently changed into null or a zero value.

Horos checks Type/Data bounds on the source and all canonical applications.
Native declarations are value representations, not conformance validators.
Native strings and records do not enforce every Horos invariant, and an
unconstrained native parameter does not establish the Data bound. The mapping
profile must explicitly expose that limit rather than advertise native proof.
Unsupported constructors, target forms and required representation policies
fail without usable partial output.

## Subsequent mappings and library ownership

Arbitrary integers require arbitrary-width representations (TypeScript bigint,
Go big integers, Rust BigInt and C++ multiprecision), with dependencies and
ownership rules declared by the target profile. Atoms and tuples follow as
separately qualified representations. Sets/maps need structural identity and
canonical-order rules; unions need qualified sum representations; recursion
needs target layout/indirection policies. Exact values, refinements,
intersections and complements must not silently widen to an ordinary carrier.
These forms are outside the initial pagination implementation.

Wire requires an explicitly configured native type representation. Function,
Stream, Cell and future contracts remain ordinary library declarations. Optional
library recipes select structured target types by exact source declaration and
target profile, with explicit argument mapping and locked dependencies. They
never add a named contract to the engine, carrier, formatter or runtime. Recipe
SPIs and capability representations remain a later types-only increment; this
implementation does not pretend missing representations are available.

## Layout and qualification

Shared request admission lives in contracts/semantic/go and
contracts/transform/go. The engine passes explicit input and captures checked
intermediate outputs. Source-specific reading helpers live in
internal/horosgen. Concrete mappings live under
transforms/horos-to-typescript/go, transforms/horos-to-go/go,
transforms/horos-to-rust/go and transforms/horos-to-cpp/go. Each owns only its
target syntax construction; the generic kernel has no Horos/Ontos imports.

Qualification covers stale/mismatched/incomplete receipts, multiple targets
sharing one snapshot, route isolation, cancellation, bounded work, empty output
on failure, imported and nested pagination applications, equal-spelled binders, aliases, names,
origin links, semantic rejection and unsupported types. Independent consumers
compile all four generated outputs through maintained hosted gates. AST and
semantic assertions supplement compilation; compiler success alone does not
prove Horos value conformance. No native generator runs on the workstation.

## Public composition

Each concrete mapping exposes `New(from, to, Config)`. The caller explicitly
registers the endpoint languages, checkers and transformer. The source descriptor
uses the Horos AST profile; the destination profile is the caller's promise about
its native checker. The target marker and language ID must match the mapping.
The engine has no language or constructor catalogue.

```go
resolved, err := engine.Resolve(ctx, snapshot, semantic.Options{})
// Handle err before continuing.
analyzed, err := engine.Analyze(ctx, analysis.Input{
    Snapshot: snapshot, Resolution: &resolved,
})
// Handle err before continuing.
input := transform.Input{
    Source: snapshot,
    Evidence: []semantic.Result{analyzed},
}
typescriptAST, err := engine.TransformTo(ctx, input, "typescript")
goAST, err := engine.TransformTo(ctx, input, "go")
```

The mapping requires the Horos analysis operation/profile, its complete validity
and enumeration reports, an available model and empty analysis configuration.
Resolution configuration is part of the analyzer
provider contract; the mapping does not resolve or reinterpret it. The adapter checks capture identity and evidence before the mapper runs.
Syntax-only mappings declare an empty evidence list. Multi-hop routes discard the
previous evidence and capture the checked intermediate tree. If a later edge
needs semantics, the caller must execute that leg separately after explicitly
resolving/analyzing its new snapshot; the engine never runs hidden analysis.

## Profile limits

Names must fit the documented target syntax. Declarations and generic parameters
capitalize their initial ASCII letter (`page` becomes `Page`); a leading underscore
becomes `Type_`, and `Self` becomes `Self_`. Conversion collisions, declaration
shadowing by generic parameters and unsupported field spellings are errors. Go
fields also capitalize their initial letter; TypeScript retains property names.
Rust and C++ reject keyword field names. No member name is guessed from its type.

Records must be named declarations. Nested applications such as
`Page<Page<Person>>` are supported; anonymous records inside fields/arguments are
rejected. Unused generic parameters, recursive declarations, open records and
constructors outside the table are rejected with source-node diagnostics and no
output. This is an explicit first subset, not a fallback that drops constraints.
Name anonymous records before mapping. Phantom/recursive representations and
per-library recipes need separately specified profiles.

Go alone reserves the generated support name `BittypeOptional`; other mappings
do not know that support declaration. C++ declaration references use the fully
qualified configured namespace, so a member named `Person` can still have the
type `Person`. A field cannot shadow a C++ template parameter.
The configured C++ namespace must be an unreserved identifier and cannot be `std`.

The native graph retains separate target nodes for equal-spelled binders and
references to their actual declarations. Origins point back to current canonical
AST nodes. Standard-library types are qualified in Rust/C++. Dependencies are
ordered before uses, and each imported declaration is generated once. The helper
indexes the admitted semantic graph and decodes descriptor payloads on demand;
it never calls Horos resolution, specialization or template lifting. Generation
bounds target nodes to 100,000, recursive traversal depth to 128 and total visits
to 800,000, with context cancellation. Formatter/carrier budgets still apply.
