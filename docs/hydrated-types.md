# Optional Wire types and Bittype tooling

**Owner decision, 3 November 2026. Status: Go profiles, frontend and explicit lifting implemented.**
Optional `ontos/wire` owns the mixed value profile; optional `horos/wire` owns
its type/declaration model and lifting semantics. Bittype supplies tooling for
those definitions. This replaces the unimplemented Bittype-owned contract
language, including its `bittype/contract/*` syntax and proposed `contract`
stage directories. There is no compatibility alias or fallback decoder.

Ground Ontos, ground Horos declaration revision 1, the existing Horos frontend
retain their public semantics; the agnostic engine remains unchanged. This does not revive the earlier
proposal to add Wire to ground Horos revision 2. The
[preservation map](wire-extension-preservation.md) records where every prior
concept now belongs, including ideas retained for later implementation.

## 1. Ownership and value domains

| Component | Responsibility and status |
| --- | --- |
| Deixis | Existing generic structure over opaque contents; not an AST or type model |
| Ground Ontos | Existing immutable data values, structural identity and encoding |
| Optional ontos/wire | Implemented Go mixed data/capability profile over Deixis |
| Ground Horos | Existing data descriptors, declarations, binding, specialization and conformance |
| Optional horos/wire | Implemented static Go mixed types, explicit ground embedding and partial template lifting |
| Bittype providers/transformers | Tooling projections of the selected Horos profile and native type representations |
| Bittype agnostic kernel | Generic trees, semantic facts and explicit registration; no concrete Wire semantics |
| Bitwire / bitruntime | Existing sending capability contracts and context-dependent hydration |
| Contract libraries | Function, Stream, Cell declarations, native representations and behavioral laws |

The authoritative boundaries are the
[Ontos optional value profile](https://github.com/Bitspark/ontos/blob/main/docs/design/0007-optional-wire-extension.md)
and [Horos optional type profile](https://github.com/Bitspark/horos/blob/main/docs/spec/wire-types.md).
The [frontend guide](horoswire.md) identifies implemented Go packages and the
remaining runtime/native-generation work.

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

Bitwire's inspected hydrated domain already has the same conceptual atom/tuple/
capability shape. Alignment with `ontos/wire` is future qualification, not an
implemented type alias or a requirement to copy every runtime value into a
second heap of Deixis nodes. Runtime admission of a capability remains distinct
from structural profile admission. See [runtime values](runtime-values.md) for
the actual current Go types, exports, proxies and lifetime ownership.

Every type description is itself ground data. Describing a Wire does not create
one, and loading, resolving or analyzing a description performs no sends,
endpoint allocation, capability binding or export. Static tooling needs no live runtime.

## 2. Canonical format belongs to Horos

[Horos's Wire type proposal](https://github.com/Bitspark/horos/blob/main/docs/spec/wire-types.md#2-canonical-expression-format)
owns the exact value grammar. Its separate `horos/wire/*` envelope and profile
identify mixed declarations; the ground `horos/decl/*` revision remains separate.
The central distinction is:

```text
Data(G)  -- closed, well-formed ground Horos descriptor G
Wire(E)  -- sending capability with declared message type E
```

`Data(any)` excludes values containing a Wire. Complements inside Data remain
ground complements. Exact literals are opaque even when they resemble parameter,
application or Wire syntax. No open upper parameter may be inserted into G.
Wire's schema message operand is traversed as syntax; a live Wire remains a
value leaf, so traversal never follows the receiver.

The optional profile also owns bounded tuple/compound/list/closed-record/union
composition, parameters and applications. Ground layout and literal record keys
are preserved; new live map keys, set elements, mixed top/complement or recursive
named applications are not implicitly introduced. Closed unsupported ground
forms remain available through Data. These are Horos rules consumed by a
frontend, not a second semantic specification in Bittype.

Modules preserve exact import pins, ordered lexical slots, explicit arity,
caller-side argument binding and source declaration identity. A profile remains
explicit even if its expressions contain no Wire. Specialization is separate
from generic native generation; it must not erase the original open declaration
or use equal expanded shapes to substitute a different library.

## 3. Embedding, lifting and pagination

Ground `Page<Person>` can be specialized and embedded as `Data(G)`. Template
lifting instead derives a separately identified mixed declaration, preserving
the authored Page's parameter and applications. It follows
[Horos's partial lifting rules](https://github.com/Bitspark/horos/blob/main/docs/spec/wire-types.md#4-explicit-embedding-and-partial-template-lifting);
Bittype implements the corresponding explicitly selected AST transformation.

```text
Ground Page<T>                 Derived mixed Page<T>
  items: list-of(T)              items: list-of(T)
  next?: scalar(utf8-text)       next?: Data(scalar(utf8-text))
```

The derived Page can accept `Cell<Data(PersonDescriptor)>`; the original ground
Page cannot. No second Page is authored and no `LiftedPage` spelling is required.
Derived module/declaration identities and provenance remain distinct. The
[tooling design](type-representations.md) preserves the full ordered rules,
dependency-closure counterexamples and native representation configuration.
Route discovery never chooses a different domain or turns a failed lift into
closed embedding. Direct ground-to-native generation remains independently useful.

## 4. Library-owned Function, Cell and Stream

The complete canonical module equations now live with
[Horos's proposed format](https://github.com/Bitspark/horos/blob/main/docs/spec/wire-types.md#5-library-owned-function-cell-and-stream).
In explanatory notation, with H constructing a ground Horos descriptor:

```text
Unit = Data(H("tuple"))
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

Bittype will project the optional Horos profile through the existing AST and
semantic SPIs. The current engine fixes one syntax profile per language ID and
rejects same-ID transformation edges. Consequently, an explicit ground-to-mixed
edge needs distinct registered endpoint IDs. Exact IDs, profile strings and Go
marker choices must be qualified with implementation; a different Go marker is
an isolation choice, not an additional existing engine requirement.

The [stage-first layout](../LAYOUT.md) remains in force. Directory
name `horoswire` denotes the optional Horos frontend, not a new Bittype model:

| Component | Responsibility |
| --- | --- |
| `syntax/horoswire/go` | AST projection vocabulary and language identity |
| `check/horoswire/go` | Syntax-only admission; unresolved names may reach resolution |
| `resolve/horoswire/go`, `internal/horosmodel` | Adapt the authoritative module/binding model with immutable capture and explicit dependencies |
| `analyze/horoswire/go` | Adapt established semantics to shared facts/reports |
| `provider/horoswire/go` | Explicit AST-only capability composition |
| `transforms/horos-to-horoswire/go` | Explicitly selected Horos-owned partial lift (implemented) |
| `transforms/horoswire-to-<native>/go` | Native type AST mappings and uniform library recipes (future) |

All listed stages except native mappings now exist. The existing `horos` stages retain
their ground profiles. Shared `contracts/` contains agnostic SPIs, not canonical
behavioral declarations. A new frontend requires narrow reviewed model dependency
allowances and capability tests; these are enforced by the architecture suite.

Provider-owned terms and located operands describe Wire and Data. Generic facts
retain references/applications, graph integrity and request-bound evidence without
decoding their meaning. Keep syntax admission, semantic validity, completion and
model availability separate. Ground conformance, mixed structure checks, runtime
capability recognition and remote behavioral compliance are different judgments.

Finite input is not necessarily bounded enough to process. Admission, resolution,
lifting and generation each need explicit node/depth/byte/module budgets where
applicable, with exhaustion separate from invalid syntax, semantic non-membership
and unsupported mappings.
An incomplete stage publishes no usable partial model or target. These use existing result dimensions, not new kernel cases.

Static dependencies point from the concrete Bittype frontend to Horos's model
and ground Ontos document representation. They do not point from Horos to
Bittype, or from type tools to a live runtime. The eventual runtime realization
of the optional value profile needs its own bounded interface/dependency design;
no import cycle or default runtime is authorized here.

## 6. Qualification and forbidden shortcuts

Ground Horos/Bittype APIs continue to reject the optional format. The separate
[optional frontend tests](../integration/go/horoswire_test.go) and
[lifting tests](../transforms/horos-to-horoswire/go/lift_test.go) qualify profiles and grammar,
operand/literal boundaries, generic scope, identity/provenance, ordered lifting,
and ground-membership/layout preservation. Native representation limits, actual
runtime alignment and library protocol/lifetime tests are separate evidence.

Forbidden: widening ground `WF`/`Conforms`, Ontos containers or generic slots;
guessing a profile from the absence of Wire; compatibility aliases for superseded
proposal tags; core Function/Stream/Cell names, imports, switches, registries or
shape fingerprints; a new Bittype application-value algebra; pure serialization
of live handles; and a replacement hydration/export engine. Relocating a model
does not itself eliminate duplicate implementations: future work must demonstrate
the selected realization and its preservation laws.
