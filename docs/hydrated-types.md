# Unified Horos types and value domains

**Owner decisions, 4 November 2026. General kinds implemented in Go source.**
One Horos declaration language describes data and capability types. Revision 4
uses ordinary structural kinds, including an imported Data definition, instead
of closed Type/Data bounds or a dedicated Data evaluator. Ontos/wire remains an
optional *value* model. This distinction requires no second declaration language
or copied Page. The current migration does not change the runtime model below.

## 1. Ownership and value domains

| Component | Responsibility |
| --- | --- |
| Deixis | Generic structure over opaque contents |
| Ontos core | Ground data values, structural identity and encoding |
| Ontos/wire | Mixed value structure and exact ground embedding/recovery |
| Horos | Ground conformance; unified definitions, ordinary kinds, resolution and type-level application |
| Bittype concrete frontend | Project authoritative Horos syntax and semantic facts |
| Bittype agnostic kernel | Generic trees, facts, stage contracts and registration |
| Bitwire / bitruntime | Sending capability contracts and hydration |
| Contract libraries | Function, Stream, Cell types, representations and behavior |

The [Horos frontend guide](horos.md) documents the revision-4 migration and APIs. The
[preservation map](wire-extension-preservation.md) records retained concepts and
superseded proposals. Native pagination has a separate representation profile;
runtime adapters remain future work.

```text
GroundValue = ByteAtom | Tuple(GroundValue*)
MixedValue  = ByteAtom | MixedTuple(MixedValue*) | Wire
```

Deixis can carry opaque values, but its arbitrary nodes are broader than this
profile. Admission must distinguish childless byte and Wire leaves from tuples
with dense canonical positional keys; bytes and child structure must be captured
immutably. A Wire remains an opaque live sending capability, not a ground atom,
receiver, closure authority or arbitrary object with a Send method. Wire-free
values project to ground only after mixed-profile admission. This introduces
neither a total pure live-to-ground codec nor universal Wire equality or hashes.

Bitwire 0.7.0 specializes the shared native `ontos/wire` model for its sending
contract. Bitruntime 0.8.0 uses its immutable tuple factory with runtime-owned
handle recognition; sends do not copy every value into a second heap of Deixis
nodes. Runtime admission of a capability remains distinct from structural
profile admission. See [runtime values](runtime-values.md) for the actual current
Go types, exports, proxies and lifetime ownership. The static frontend is moving
to Horos declaration revision 4; Bittype does not import the live carrier or runtime.

Every type description is itself ground data. Describing a Wire does not create
one, and loading, resolving or analyzing a description performs no sends,
endpoint allocation, capability binding or export. Static tooling needs no live runtime.

## 2. Canonical syntax and ordinary kinds

[Horos revision 4](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md)
retains the module envelope and ordinary imports, with kind and type definitions
in one declarations map. Kind definitions have sort-annotated parameters; type
definitions have kind contracts and an explicit result kind. Indexed binders and
expression-callee applications retain their scopes. There is no schema-digest
field or independently managed formation artifact.

`H("wire", E)` remains a capability description. A ground descriptor appears
directly, without a schema-level Data wrapper. Every type description is ground
data regardless of what its instances would contain. `Data` is an ordinary
recursive kind in the explicitly imported prelude; the general machinery also
handles unrelated categories such as Scalar, List, Record, Flat, Rows and Tree.

Tuple, compound, list-of, closed record and union retain mixed proper operands.
Wire is outside the ordinary ground kind even when its message is data. Other
primitive operand restrictions are expressed through ordinary kind contracts,
not a name-sensitive Data flag. Ground any and complement keep their data
universe. Literal payloads remain opaque. No mixed equality, key semantics or
recursive capability law follows from this shared syntax.

## 3. Pagination in one language

```text
type Page<T> = { items: list<T>; next?: text; };
type DataPage<D: h.Data>: h.Data = Page<D>;
type Person = { name: text; };
```

Here `h` is an ordinary explicit import of the `horos.kinds` prelude. Omitted
type annotations elaborate to `top(P)`, not a reserved name Type.
Page<Person>, DataPage<Person> and DataPage<Page<Person>> establish data results.
Page<Cell<Person>> uses the same Page identity but can contain capabilities;
DataPage<Cell<Person>> fails if that ordinary library Cell contains a Wire.
Forward<T> = DataPage<T> fails symbolically; Forward<D: h.Data> succeeds. Known
bodies retain actual parameter relationships for requested kind queries. Abstract
constructors expose declared guarantees only. Neither names nor shape fingerprints
establish a library's behavior.

The [canonical model](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md)
retains module identities, caller bindings and ordered slots. Sealed Program-owned
values support application and ordinary membership queries without mutating
generic definitions. GroundDescriptor requires a closed proper value, ordinary
ground-kind membership, materialization and unchanged ground WF. The old
Specialize/DataDescriptor wrappers and dedicated Data proof path are replaced.
Template lifting and schema wrappers remain removed; Ontos value embedding is
unchanged.

## 4. Library-owned Function, Cell and Stream

The complete canonical module equations now live with
[Horos's canonical examples](https://github.com/Bitspark/horos/blob/main/docs/spec/wire-types.md#2-library-owned-function-cell-and-stream).
In explanatory notation, with H constructing a ground Horos descriptor:

```text
Unit = H("tuple")
Function<A,B> = Wire<tuple(A, Wire<B>)>
Cell<T> = record closed { get: Function<Unit,T>, set: Function<T,Unit> }
Step<T> = compound("item", T) | compound("end")
Stream<T> = record closed {
  next: Function<Unit,Step<T>>
  cancel: Function<Unit,Unit>
}
```

This preserves the success-only Function example, get/set Cell and pull Stream.
It does not claim a complete production protocol. Unit describes the bare empty
tuple; native void is a separate explicit mapping. Reply matching, errors,
timeouts, demand, cancellation, mutation and resource cleanup belong to the
libraries. A leaf cannot prove a remote peer follows its declared message type.
No variance, subtyping or protocol-equivalence inference follows from these
shapes or a successful syntax/analysis report.

Native type recipes remain library-owned, selected by exact declaration/profile
identity in consuming build configuration. Ordinary structural composition needs
no custom recipe. Neither names nor expanded shapes select behavior or recipes.
Adding another library requires no generator kernel, runtime, Wire or hydration
case. [Type recipes](type-representations.md) and future
[runtime adapters](runtime-values.md) remain separate work.

## 5. Provider, pipeline and dependencies

The stage-first layout contains one syntax/check/resolve/analyze/provider Horos
family. `internal/horosmodel` adapts authoritative modules, bindings and kinds;
`internal/declsyntax` shares private syntax mechanics. Generic contracts are
unchanged. Check admits kind/type syntax; resolution establishes meaning;
analysis consumes explicit compatible evidence and the same typed Semantics.
Multiple Program trees and their tree-qualified sites remain supported.

Provider-owned facts distinguish sorts, declared/result kinds, checked symbolic
bodies and individual memberships with their actual assumptions. Absence is not
nonmembership; no principal kind is published. The kernel stores these terms
without interpreting Type, Data, Scalar, Record, List or Wire. Native mappings
use `transform/syntax/horos-to-<native>/mapping/go` with an explicit representation profile;
unimplemented enforcement returns Unsupported. No concrete library introduces
an engine case, built-in registration or runtime dependency.

Finite inputs still need bounded admission, resolution, symbolic checking and
specialization. Exhaustion is an incomplete request, not semantic rejection or
usable partial output. Source/context attribution and exact operand locations
remain separate from historical provenance.

## 6. Distinct judgments and forbidden shortcuts

Ground conformance, general kind membership, mixed structural admission, live
runtime recognition, authority and protocol behavior remain different judgments.
The revision-4 source migration uses public Horos Go 0.10.1.

Forbidden: old-profile aliases/readers, widening ground any or values, guessing
behavior from library names/structure, core Function/Stream/Cell switches or
imports, a new Bittype application-value algebra, pure serialization of live
handles, or a replacement hydration/export engine. Native type generation does
not provide runtime bindings.
