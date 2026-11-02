# Hydrated Horos types and library composition

**2 November 2026. Proposal only; no implementation.** The owner requested a
native Wire atomic type in the canonical model, with Function, Stream and Cell
declared by ordinary libraries. Horos owns the proposed exact format and meaning;
bittype owns its future AST/semantic projections and native transformations.
The current [Horos frontend](horos.md), ground declaration revision 1, native
text tooling and agnostic engine remain unchanged.

## The missing type boundary

Bitwire already defines a value domain containing opaque sending capabilities:

```text
GroundValue   = ByteAtom | Tuple(GroundValue*)
HydratedValue = ByteAtom | HydratedTuple(HydratedValue*) | Wire
```

Horos currently classifies the first domain. The [Horos hydrated type proposal](https://github.com/Bitspark/horos/blob/main/docs/spec/hydrated-types.md)
specifies declaration revision 2 for the second. A Wire is atomic in the sense
of a live leaf; it is not an ontos byte atom. No new reference encoding or
hydration mechanism is required. Native Wire recognition and lifetime remain
with the [existing bitwire contract](https://github.com/Bitspark/bitwire/blob/e2d2d641744a49f1d8098b6ea5f2724454d72b96/docs/wire/hydrated-native.md)
and its runtime realization.

Type descriptions, modules and their pins remain ground ontos data. They describe
capabilities without containing any. Loading/analyzing a schema must not import
a live runtime, create/export endpoints or perform sends.

## Exact canonical input

The proposed primitive has exactly one child type expression, its message type:

```text
H("wire", E) = Tuple(Atom(UTF8("horos/wire")), E)
```

It conveys sending authority only. There are no core operation names, reply
policies or Function/Stream/Cell discriminators. Using the existing Horos value
abbreviations, its ordinary generic library declaration is:

```text
WireLibrary = D("module", I(2), T("example.wire"), M{}, M{
  "Wire": D("type", L(T("M")),
    H("wire", D("param", I(0))))
})
```

`D` prefixes declaration tags with `horos/decl/`; `H` prefixes type tags with
`horos/`; `I`, `T`, `L` and `M` are canonical integer, text, list and text-keyed
map embeddings. `A` in the full specification is a raw UTF-8 atom. These are
exact value equations, not implemented authoring functions or a readable parser.
The [owning specification](https://github.com/Bitspark/horos/blob/main/docs/spec/hydrated-types.md#2-exact-wire-type-format)
defines arity, operands, domain selection and the full library module values.

Revision 2 is explicit. The proposal initially admits only revision-2 imports;
there is no implicit widening of a revision-1 dependency or fallback for an
unknown constructor. A specialized expression retains its hydrated profile even
if it is simply `any()` or a complement with no Wire tag. Bittype must not feed
such an expression into the current ground checker and report the narrower
answer as equivalent. Current source accepts neither this module revision nor
this new constructor.

## How the libraries compose

The following is explanatory shorthand for the exact pinned module equations in
Horos's proposal. It is not an additional canonical format or an implemented
bittype source grammar:

```text
// example.wire
Wire<M> = primitive wire(M)

// example.contracts.function imports the exact example.wire module
Unit = tuple()
Function<A,B> = Wire<tuple(A, Wire<B>)>

// example.contracts.cell imports the exact Function module
Cell<T> = record closed {
  get: Function<Unit,T>
  set: Function<T,Unit>
}

// example.contracts.stream imports the exact Function module
Step<T> = union(compound("item", T), compound("end"))
Stream<T> = record closed {
  next: Function<Unit,Step<T>>
  cancel: Function<Unit,Unit>
}
```

Function is a success-only illustrative request/reply representation: each
request contains A and a fresh reply Wire for B. Cell composes two Functions;
this Stream is a pull composition with next/cancel Functions. The libraries
retain ordinary module/declaration identities and indexed binders. These examples
do not replace current service contracts or establish deployed compatibility.

The declaration describes representation. Function's implementation owns reply
matching, failure/completion rules and resource management; Cell owns its state
and mutation guarantees; Stream owns demand, ordering and terminal/cancellation
behavior. Those laws must be specified and tested in the libraries. A generic
record composition does not derive them. Unit denotes a bare empty tuple; a
native `void`, callback or iterator surface requires an explicit library mapping.

A future generator follows imports and generic applications, preserving identities
through aliases and nested uses. If a library selects a custom native surface or
adapter, its binding association must identify the resolved declaration, not a
bare name or an expanded shape. Consumers explicitly select these library-owned
exports. The exact adapter metadata/API and executable bindings remain outside
this documentation task and the current types-only implementation milestone.

## Projection through the agnostic kernel

The generic AST and semantic model already have carriers for this proposal:

| Layer | Proposed treatment |
| --- | --- |
| Horos authoritative model | Recognize revision 2, the Wire message-type position and the hydrated profile's semantics |
| Horos AST vocabulary | Reuse a Descriptor node with constructor attribute `wire` and one operand; keep ordinary declarations, parameters and applications |
| Horos projection/resolution | Obtain the message operand at tuple path `[1]` from Horos's authoritative expression model; preserve input profile and binding identities |
| Shared semantic facts | Reuse ProviderExpression with exact encoded schema bytes and the located message operand; use existing Reference/Application terms for generic arguments |
| Analysis provider | Report validity, coverage and model availability for the selected declaration profile; supply facts without inspecting a live capability |
| Engine and registration | Explicitly register the matching Horos syntax descriptor/analyzer; no automatic provider, named contract lookup or new generic core kind |
| Future target transformer | Produce the selected native type AST using preserved declarations and explicit library representation choices |
| Existing formatter/printer | Consume the target language AST/CST through their existing contracts |

The syntax and analysis profile identifiers must distinguish the new domain from
the current `horos-declarations-1/ast-2` and
`horos-declarations-1/analysis-2`. The exact identifiers are assigned with a
qualified implementation, not silently reused. An engine instance currently
fixes one syntax profile per language ID; registering the new profile must obey
that existing rule rather than introduce negotiation or implicit conversions.

No new universal semantic term kind or generic core Wire interpretation follows
from this proposal. The Horos provider retains the primitive's meaning in its
owned vocabulary. The agnostic kernel validates graph integrity and publication;
it does not decode the schema payload or decide what a Wire message means.
Providers and transformers consume authoritative operand information instead of
copying a second Horos constructor switch.

For `Cell<Person>`, the application binds Cell's T to the caller's Person. The
get member applies Function to Unit and that bound Person; set reverses those
argument positions. Function's body applies the imported Wire declaration to
its request tuple, whose reply position applies Wire to its B. The same graph
queries can inspect each binding without erasing the named declarations or
automatically expanding them. A field or alias called `Cell` alone establishes
none of these identities.

## Structure, behavior and ground identity

Recognizing a live Wire proves only that it is a sending capability admitted by
the runtime. `Wire<Text>` additionally declares a message contract; observing
the leaf cannot prove that the remote receiver follows it. Typed bindings can
check actual messages at their boundaries. Neither a semantic analysis Holds
report nor a local shape match establishes remote execution, one reply,
cancellation or protocol equivalence. No variance or assignability is inferred.

Horos's proposed profile specifies which structural constructors lift. Tuples,
lists, record fields and map values can contain Wires. Map keys and set elements
remain ordered ground data; exact literals remain ground; open record extras
range over the selected hydrated domain. Do not use existing ground container
objects to smuggle live leaves into ontos. Descriptions remain serializable,
while values containing Wires have no newly invented structural equality or
content hash. Recursion cannot capture an outside rec/ref binder through a Wire
message boundary; message expressions may use generic parameters or their own
closed, guarded recursion. The owning specification supplies the exact rules.

These distinctions matter for native generation too: a target type that can
store a capability does not prove that its code implements a library's protocol.
The native representation profile must state its limits. Current native ASTs,
text round trips and compiler tests establish none of these future mappings.

## Repository layout and dependency direction

This proposal preserves the [stage-first matrix](../LAYOUT.md):

| Existing package | Future responsibility, when implemented |
| --- | --- |
| `syntax/horos/go` | Describe the qualified new AST profile using existing carrier vocabulary |
| `check/horos/go` | Admit that syntax profile, preserving the syntax/semantic distinction |
| `resolve/horos/go` and `internal/horos` | Project/adapt authoritative revision-2 declarations and operands |
| `analyze/horos/go` | Publish the new profile's checked facts and reports |
| `provider/horos/go` | Explicit composition of the qualified capabilities |
| `trees/*`, `semantics/go`, `contracts/*`, `engine/go` | Remain independent of concrete Wire and contract-library semantics |

Future `transforms/horos-to-<native>/go` mappings retain their planned location.
Library-specific native types and adapters belong in their owning libraries;
consuming composition imports them. No placeholder package, production
dependency or generator implementation is introduced now. The generic kernel
must not depend on bitwire, bitruntime or the Function/Stream/Cell libraries.

The dependency direction is: application composition -> chosen transformers and
contract libraries; Horos frontend -> Horos's schema model; behavioral library
implementation -> existing Wire/hydration. Loading a schema model need not load
the implementation of the live capability it describes.

Forbidden are core Function/Stream/Cell kinds, name switches, built-in imports,
default registrations, behavior chosen from structural fingerprints, generated
central dispatch tables, or a new reference/export engine. Teaching the
Horos-specific frontend its one new primitive is not permission to teach the
general kernel behavioral contracts. Adding another contract after this boundary
is qualified changes the library and its consumers, not the generic core.

## Qualification before implementation claims

Future qualification must independently test exact canonical Wire/module forms,
revision mismatch, generic/alias identity, located message operands, literal
boundaries, profile-preserving specialization and separate syntax/semantic
outcomes. It must include wire-free expressions whose meaning depends on domain,
ground-only map/set positions and closed recursion boundaries. Input capture and
analysis must never send, export or construct live endpoints.

Cross-package tests must keep generic infrastructure independent of concrete
semantic packages. Native mapping/compilation, runtime capability recognition,
message validation and Function/Cell/Stream behavior each require their own
evidence. This documentation change supplies no new executable implementation,
test vectors, parser, generator, runtime adapter or package release.
