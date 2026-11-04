# Unified Horos types and value domains

**Owner decision, 4 November 2026.** One Horos declaration language describes
data and capability types, with explicit [proper/constructor interfaces](kinds.md).
Data is a refined kind of Type, certified compositionally. The former
Horos/Wire frontend, package and profiles are removed without aliases. Ontos/wire
remains an optional *value* model. This distinction does not require a second
declaration language or copied Page definition.

## 1. Ownership and value domains

| Component | Responsibility |
| --- | --- |
| Deixis | Generic structure over opaque contents |
| Ontos core | Ground data values, structural identity and encoding |
| Ontos/wire | Mixed value structure and exact ground embedding/recovery |
| Horos | Ground conformance; unified declarations, kinds/interfaces, Data certificates, resolution and specialization |
| Bittype concrete frontend | Project authoritative Horos syntax and semantic facts |
| Bittype agnostic kernel | Generic trees, facts, stage contracts and registration |
| Bitwire / bitruntime | Sending capability contracts and hydration |
| Contract libraries | Function, Stream, Cell types, representations and behavior |

The [Horos frontend guide](horos.md) documents implemented APIs. The
[preservation map](wire-extension-preservation.md) records retained concepts and
superseded proposals. [Native pagination generators](semantic-generation.md) are
implemented; higher-kind representations and runtime bindings remain future work.

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
Go types, exports, proxies and lifetime ownership. The static frontend uses Horos
revision 3 and ground Ontos 0.16.0; Bittype does not import the live carrier or runtime.

Every type description is itself ground data. Describing a Wire does not create
one, and loading, resolving or analyzing a description performs no sends,
endpoint allocation, capability binding or export. Static tooling needs no live runtime.

## 2. Canonical syntax and interfaces

[Horos revision 3](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md)
owns the only declaration grammar: D("module", I(3), ...), parameters with
proper/constructor interfaces, explicit proper result claims, indexed slots,
declaration references and expression-callee applications. The
[kind guide](kinds.md#exact-canonical-revision-3) gives the exact shape.
`H("wire", E)` describes a sending capability. A ground descriptor G appears
directly; there is no schema-level Data wrapper. Type descriptions remain data
regardless of what their instances would contain.

Tuple, compound, list-of, closed record and union admit either domain. Their
operands determine whether the result is data. Wire never becomes Data merely
because its messages are data. Other constructors retain the Data restrictions
specified by Horos; any and complement keep their data universe. Literal payloads
remain opaque. No mixed equality, key semantics or recursive capability law is
inferred from the shared syntax.

## 3. Pagination in one language

```text
Page<T> = record closed { items: list-of(T), next?: scalar(utf8-text) }
DataPage<D: Data> = Page<D>
Person = record closed { name: scalar(utf8-text) }
```

Page<Person>, DataPage<Person> and DataPage<Page<Person>> establish data results.
Page<Cell<Person>> uses the same Page identity but can contain capabilities;
DataPage<Cell<Person>> fails if that ordinary library Cell contains a Wire.
Forward<T> = DataPage<T> fails symbolically; Forward<D: Data> succeeds. Phantom
parameters do not contaminate a data result. Actual parameter use and target
interfaces matter, not names or shape fingerprints. Known constructor arguments
retain this precision through Wrapped/Compose; abstract constructors expose only
their declared admission and result guarantees.

The [complete canonical pagination input](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md#5-complete-pagination-example)
retains pins, caller bindings and ordered slots. Specialize produces closed type
expressions without mutating generics. DataDescriptor additionally requires data
before ground WF/Conforms. Template lifting, closed schema embedding and
cross-profile transforms are removed; Ontos value embedding is unchanged.

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
family. `internal/horosmodel` adapts authoritative modules, bindings and interfaces;
`internal/declsyntax` shares private syntax mechanics. The generic semantic graph
stores expression callees and subject-associated facts. Check validates interface
syntax; resolution establishes its semantics;
analysis consumes explicit compatible evidence and the same immutable graph.

Kinds, interfaces, Data summaries/statuses and canonical identities are
provider-owned terms queried through Model.Facts. The required type-facts report
establishes their enumeration coverage. The kernel stores their links and bytes
without interpreting Type, Data or Wire. Native pagination mappings use
`transforms/horos-to-<native>/go` and reject residual higher kinds. No concrete type introduces
an engine case, built-in registration or runtime dependency.

Finite inputs still need bounded admission, resolution, symbolic checking and
specialization. Exhaustion is an incomplete request, not semantic rejection or
usable partial output. Source/context attribution and exact operand locations
remain separate from historical provenance.

## 6. Qualification and forbidden shortcuts

The AST/semantic interface fixtures, existing
[binding tests](../resolve/horos/go/resolve_test.go) and independently installed
[pagination consumer](../integration/go/testdata/installed/consumer_test.go)
qualify one language, explicit interfaces, unchanged binder identity and data/mixed
specialization. Horos additionally tests library-owned Function/Cell composition,
phantom arguments, bounded proofs and retired-input rejection. Ground conformance,
mixed structural admission, runtime recognition, authority and protocol behavior
remain different judgments with separate evidence.

Forbidden: old-profile aliases/readers, widening ground any or values, guessing
behavior from library names/structure, core Function/Stream/Cell switches or
imports, a new Bittype application-value algebra, pure serialization of live
handles, or a replacement hydration/export engine. Static native pagination
generation creates no runtime bindings.
