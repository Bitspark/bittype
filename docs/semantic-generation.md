# Semantic native type generation

**Owner direction, 4 November 2026. Implementation in progress.** Generate native
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

The first profile covers named/imported declarations, generic binders and
applications, closed records, UTF-8 text, booleans, lists and optional fields.
Canonical module identity and target package/file placement remain separate.
Mappings retain declaration/binder links and current source origins. A target
configuration explicitly supplies native module locations. Generation has no
implicit network or filesystem IO.

| Horos form | TypeScript | Go | Rust | C++ |
| --- | --- | --- | --- | --- |
| UTF-8 text | string | string | String | std::string |
| Boolean | boolean | bool | bool | bool |
| list-of(T) | T[] | []T | Vec<T> | std::vector<T> |
| Closed record | interface/object | struct | struct | struct |
| Optional field | field?: T | explicit Optional[T] | Option<T> | std::optional<T> |
| Type parameter | native parameter | native parameter | native parameter | template parameter |
| Application | F<T> | F[T] | F<T> | F<T> |

TypeScript requires strictNullChecks and exactOptionalPropertyTypes. Go's
presence helper is a generated type declaration with Present and Value fields;
the value is ignored when absent. No helper methods or runtime implementation
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
on failure, imported/nested pagination, equal-spelled binders, aliases, names,
origin links, semantic rejection and unsupported types. Independent consumers
compile all four generated outputs through maintained hosted gates. AST and
semantic assertions supplement compilation; compiler success alone does not
prove Horos value conformance. No native generator runs on the workstation.
