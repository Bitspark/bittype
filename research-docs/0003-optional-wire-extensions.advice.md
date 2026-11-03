# Advice for consult
**Research doc:** doc_1ad678c3-1fa0-4220-9d39-efc374e4a5ae.md
**Document:** doc_1ad678c3-1fa0-4220-9d39-efc374e4a5ae
**Advisor:** ChatGPT, model "6 Pro" as labelled on the answer page
**Session:** nsess_f4bf5e68-6d60-4145-9e35-a45c7fbaf842
**Run:** run_38589d6b-89c8-4c4d-8809-1ddd2738aa3c
**Chat:** https://chatgpt.com/c/6aea1d46-53d0-83ed-8412-f1b4141d96b8
**Verify:** nightfall consult verify run_38589d6b-89c8-4c4d-8809-1ddd2738aa3c --output <this file>

I have reviewed the research document and its five expert questions.

My overall assessment is that the selected architecture is sound and preferable to putting the canonical mixed-value/type contract in Bittype, but the proposal needs several important semantic refinements before it should be treated as a stable specification.

The most important distinction is that the three concerns must remain separate:

1. Deixis owns the structural carrier.
2. Ontos and Horos optional extensions own the mixed value and type semantics.
3. Bitruntime owns the recognition, authority, hydration and lifecycle of live capabilities.

Bittype should remain an agnostic language-tooling and transformation system, not become another semantic authority.

Below is my technical review, written as a response to the research document.

# Expert Review: Optional Wire Extensions for Ontos and Horos

**Research reviewed:** 0003 — Optional Wire value and type extensions without losing the existing design\
**Review date:** 3 November 2026\
**Review status:** Architectural assessment and recommendations\
**Scope:** Semantic design, type ownership, value representation, generic lifting, identity, tooling and future qualification.

## Executive assessment

I recommend proceeding with the selected architecture:

- `ontos/wire` owns the optional mixed data/capability value profile.
- `horos/wire` owns the corresponding optional mixed-domain type and declaration semantics.
- Deixis remains the generic structural carrier.
- Existing Ontos L0/L1/L2 and existing Horos ground semantics remain frozen.
- Bittype consumes the new profiles through explicit semantic providers and transformations.
- Bitwire and Bitruntime retain their separate transport and runtime responsibilities.

This is a coherent separation of concerns.

However, I would make five qualifications before treating the design as specification-ready.

1. A structurally admitted mixed value must not automatically become a runtime-recognized capability.
2. Ground embedding and mixed-domain operations must have explicitly different denotations.
3. Template lifting must be defined by a denotational preservation theorem, not merely by structural translation rules.
4. Declaration identity and profile compatibility must include the semantic profile, not just a module digest and declaration name.
5. Capability aliasing and byte immutability must be addressed explicitly.

These refinements do not require abandoning the proposed design. They make its existing boundaries enforceable.

---

## 1. Semantic boundaries and frozen ground models

### Recommendation

The optional extensions should be understood as new semantic domains, not extensions to the existing ground algebra.

In particular:

- Ontos L0 remains `Atom(Bytes) | Tuple(Value*)`.
- Horos ground descriptors continue describing sets of Ontos L0 values.
- Existing ground conformance remains unchanged.
- Live capabilities are admitted only through the optional mixed profile.
- Ground descriptors do not acquire an implicit ability to describe live resources.

The crucial distinction is between embedding a ground value into a larger domain and widening the meaning of the ground value itself.

Let:

- (O) be the existing Ontos ground domain.
- (M) be the admitted mixed domain.
- (e: O \rightarrow M) be the ground embedding.

For an existing ground descriptor (G), define:

{e(v)\mid v\in O,\ Conforms(v,G)}\
]

This definition is important because it prevents `Data(any)` from becoming the universal mixed-domain type.

`Data(any)` accepts all embedded ground values, not all values in the mixed domain.

Likewise:

[\
Data(not(G))\
]

means the ground complement embedded into the mixed domain. It does not mean the complement of G over the entire mixed domain.

This should be normative language in `horos/wire`.

### A necessary distinction

There are at least three different notions of validity:

1. Structural admission: the value belongs to the mixed structural profile.
2. Runtime recognition: a capability occurrence is recognized by the runtime as an exportable sending capability.
3. Operational validity: the recognized capability is usable under the relevant runtime scope, endpoint and lifecycle conditions.

These must never be conflated.

A structurally valid `WireCase(handle)` does not establish that the handle is authorized for export.

A runtime-recognized handle does not establish that its remote peer implements the declared message protocol.

A successfully admitted send does not establish remote execution or successful delivery.

These are separate contracts.

### Counterexample

Suppose an arbitrary application defines:

`MyObject.Send(HydratedValue) error`

It implements the Go interface used by a hydrated Wire.

That alone must not make it an exportable capability.

The existing Bitruntime recognition boundary correctly prevents this. The optional Ontos profile must preserve that distinction.

**Decision:** Accept the proposed ownership architecture, with explicit structural, runtime and operational validity boundaries.

---

## 2. Mixed Deixis value profile

### Recommendation

The proposed `MixedPayload` model is a good starting point:

- `None` represents tuple nodes.
- `Some(BytesCase(bytes))` represents ground atoms.
- `Some(WireCase(handle))` represents live sending capabilities.

The admitted value domain should be exactly the subset described by the proposal.

I would specify its recognizer as a partial structural recognizer with the following rules.

**Atom case**

A node is an atom if:

- Its own value is `Some(BytesCase(bytes))`.
- It has no children.
- Its bytes are captured immutably.

**Wire case**

A node is a Wire capability if:

- Its own value is `Some(WireCase(handle))`.
- It has no children.
- Its handle is retained as an opaque reference.
- No equality, hashing or serialization operation is implied.

**Tuple case**

A node is a tuple if:

- Its own value is `None`.
- Its child keys are exactly the canonical dense positional keys.
- Every child recursively satisfies the mixed-value recognizer.

All other structures are rejected.

This gives the profile a simple, exact structural grammar.

### 2.1 Byte immutability is a normative requirement

There is a subtle issue in the proposal that deserves particular attention.

Deixis freezes the tree structure. It does not freeze the contents of an arbitrary `T`.

For example, in Go:

`[]byte` is mutable.

If a BytesCase retains a caller-owned slice, the caller can modify the bytes after admission. That would violate the intended ground-value preservation law.

I therefore recommend:

- Copy bytes when constructing or admitting an atom.
- Do not expose mutable internal byte storage.
- Copy keys when admitting child structures.
- Define the immutable observation semantics independently of the native representation.

The implementation may use slices internally, but the public contract must provide immutable behavior.

### 2.2 Handle aliasing must be preserved

Consider a mixed tuple containing the same capability twice:

`Tuple(w, w)`

The structural reconstruction must retain the same underlying capability reference in both positions.

It must not:

- Clone the capability.
- Reconstruct two independent handles.
- Invoke the capability.
- Bind it to a new runtime scope.

This is reference preservation, not value equality.

The profile should explicitly distinguish:

- Reference identity preservation.
- Structural equality.
- Capability equality.
- Export identity.

Only the first is required at the generic structural layer.

Export identity remains a runtime concern.

### 2.3 Ground recovery

The proposed partial ground recovery is correct.

For wire-free admitted values:

[\
R(e(v))=v\
]

For values containing Wire cases, ground recovery fails.

This gives a clean separation between the ground and mixed domains.

I would also explicitly require:

[\
R(m)\text{ succeeds}\
\iff\
m\text{ is an admitted wire-free mixed value}\
]

The recovered Ontos value must then be encoded using the existing Ontos codec.

The Deixis codec must not be substituted for the Ontos codec.

### 2.4 Identity and encoding

I recommend the following rules:

| Operation                    | Mixed-domain guarantee             |
| ---------------------------- | ---------------------------------- |
| Structural admission         | Yes, for recognized profile values |
| Ground recovery              | Yes, for wire-free values          |
| Ground equality              | Yes, after recovery                |
| Capability equality          | No generic guarantee               |
| Mixed-value hashing          | No generic guarantee               |
| Ground canonical encoding    | Yes, after recovery                |
| Capability portability       | No                                 |
| Runtime export               | Runtime responsibility             |
| Runtime proxy reconstruction | Runtime responsibility             |

The last two operations must not leak into Ontos.

### 2.5 Dependency direction

My preferred dependency graph is:

```
Deixis
  ↑
Ontos
  ↑
Ontos/Wire
  ↑
Horos/Wire
  ↑
Bittype semantic provider
```

The actual library relationships should retain the existing independent responsibilities, but this illustrates the intended direction of semantic ownership.

Bitruntime should compose the structural profile with its own capability recognition and lifecycle machinery.

Neither Ontos/Wire nor Horos/Wire should depend on Bitruntime.

Ontos/Wire should not depend on Bitwire merely to obtain a sending interface. An opaque carrier parameter or supplied recognition boundary is preferable.

**Decision:** Accept the mixed Deixis profile, subject to explicit immutability, alias preservation and partial identity guarantees.

---

## 3. Horos/Wire, generic declarations and partial lifting

This is the most important part of the proposal.

I agree that Horos/Wire should own the optional type semantics, reusing existing Horos declaration machinery where possible.

It should not modify the existing ground declaration revision.

### 3.1 Wire\<E>

The intended interpretation should be explicit:

`Wire<E>` describes a live sending capability whose admissible message domain is E.

It does not describe:

- An arbitrary object implementing a similarly named interface.
- An endpoint with receive or close authority.
- A complete remote protocol implementation.
- A guarantee of successful delivery.
- A guarantee that the receiving implementation conforms to E.

The type describes the sending contract, not the entire runtime behavior.

This is particularly important because Bitwire distinguishes a Wire from an Endpoint.

The optional type profile must preserve that distinction.

### 3.2 Data boundaries

I strongly recommend retaining explicit ground boundaries.

For example:

```
Data(H("tuple"))
Wire<E>
```

These are fundamentally different.

The