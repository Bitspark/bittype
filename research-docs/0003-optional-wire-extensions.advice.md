# Advice for consult
**Research doc:** doc_04fefdb0-fc0a-45f6-a4dd-f6c5ef9e7fe4.md
**Document:** doc_04fefdb0-fc0a-45f6-a4dd-f6c5ef9e7fe4
**Advisor:** ChatGPT, model "6 Pro" as labelled on the answer page
**Session:** nsess_f4bf5e68-6d60-4145-9e35-a45c7fbaf842
**Run:** run_31c1813d-dd34-46a5-9873-105a934581fb
**Chat:** https://chatgpt.com/c/6aea1fc4-047c-83eb-83e6-e7f144e19f72
**Verify:** nightfall consult verify run_31c1813d-dd34-46a5-9873-105a934581fb --output <this file>

The central recommendation is to proceed with the selected `ontos/wire` and `horos/wire` ownership model, but strengthen its semantic boundaries before documenting it as a normative design. The architecture is coherent, provided that value admission, type interpretation, runtime capability recognition, and protocol behavior remain four distinct responsibilities.

The following is a technical review addressing all five questions. It is an independent architectural assessment, not a claim of external human-expert consultation or implementation validation.

# Expert Assessment: Optional Wire Value and Type Extensions

## 1. Semantic boundaries and preservation of the frozen ground models

The selected direction is sound. `ontos/wire` and `horos/wire` should be optional semantic extensions, not modifications to Ontos L0 or Horos's existing ground language.

The essential distinction is:

- **Ontos L0:** portable, finite, immutable ground values with structural identity and canonical encoding.
- **Ontos/Wire:** a mixed value domain containing ground data and opaque live sending capabilities.
- **Horos ground:** descriptions of sets of Ontos L0 values.
- **Horos/Wire:** descriptions of values admitted by the optional mixed profile, including capability-bearing values.
- **Bittype:** language-agnostic AST, semantic orchestration, resolution and transformation infrastructure.

A Wire capability is not an Ontos atom, even if both occupy an own-value slot in Deixis. Likewise, a Wire type is not a Horos ground constructor.

Three qualifications are particularly important.

**First, total embedding does not imply total recovery.** Ground Ontos values can be embedded into the mixed domain, but arbitrary wire-free Deixis trees need not be valid Ontos values. Recovery must remain partial and exact.

**Second, structural lifting is not automatically denotationally sound.** A translated expression must preserve membership, not merely resemble the original syntax. Every lifting rule needs a denotational justification.

**Third, freezing the tree structure does not freeze its payloads.** Byte buffers must be captured immutably. Borrowed mutable buffers cannot satisfy the ground preservation laws merely because their containing Deixis nodes are immutable.

Recommendation: preserve all existing specifications and introduce dated, explicitly scoped amendments to the Ontos and Horos architectural guidance. Do not silently revise historical ADRs or existing ground constructors.

## 2. Mixed Deixis value profile

The proposed `MixedPayload` is a reasonable starting point:

`None | Some(BytesCase(bytes)) | Some(WireCase(handle))`

The profile should specify admission through an explicit recognizer, not through structural resemblance alone.

The normative rules should be:

1. Atom nodes contain immutable, owned byte snapshots and no children.
2. Wire nodes contain an opaque capability handle and no children.
3. Tuple nodes contain `None` and exactly dense, canonical positional children.
4. All descendants must independently satisfy the profile.
5. Unknown payload alternatives, malformed children and noncanonical positions are rejected.
6. Recognition does not grant authority to export, invoke or transport a capability.

Ground recovery succeeds only for the wire-free subset and must satisfy:

- `R(P(v)) = v`.
- Recognition accepts exactly the image of the embedding.
- Recovery preserves Ontos structural identity.
- Ground byte recovery followed by the Ontos codec produces the original canonical Ontos bytes.

These laws must not be confused with equality of Deixis codec bytes.

**Handle aliasing requires particular attention.**

Two Wire nodes containing the same handle should preserve the same underlying capability reference through structural copying and reconstruction. However, the profile must not promise that independently obtained handles are equal, interchangeable or portable.

I recommend explicitly distinguishing:

- Structural identity of the tree.
- Ground identity of recoverable Ontos values.
- Opaque handle identity, defined only by an explicitly supplied runtime recognition/equivalence mechanism.
- Export identity, owned by the runtime's hydration scope.

No default hashing, canonical encoding or equality should be defined for the mixed domain as a whole.

For dependency direction, `ontos/wire` should define the profile and its generic capability carrier without importing Bitruntime or requiring Bitwire's operational interfaces. Runtime composition should provide capability recognition and export/import behavior through explicit providers.

The existing runtime carrier can remain an optimized realization. A Deixis allocation rewrite is not a prerequisite.

## 3. Horos/Wire types, declaration reuse and partial lifting

`horos/wire` should reuse Horos declaration resolution, parameter identities, application semantics, immutable module bundles and specialization machinery.

It should introduce a separate profile and envelope identity. Existing `horos/<constructor>` ground descriptors must retain their current meaning and reject Wire-specific syntax.

`Wire<E>` should mean a recognized sending capability whose message values conform to E under the selected mixed profile.

This does not establish that a particular runtime handle is recognized, that it is exportable, or that its remote peer conforms. Those are separate runtime and protocol obligations.

### Partial lifting

The ordered lifting proposal is fundamentally correct, but needs stronger invariants.

I recommend retaining its rules with the following qualifications:

| Input                                             | Required behavior                                                                    |
| ------------------------------------------------- | ------------------------------------------------------------------------------------ |
| Parameter                                         | Preserve exact binder identity and provenance.                                       |
| Resolved application                              | Preserve exact target identity, argument order, arity and caller-scope substitution. |
| Tuple, compound, typed list, closed record        | Translate recursively, preserving all structural distinctions.                       |
| Multi-branch union                                | Translate every branch in original order.                                            |
| Singleton union                                   | Translate its branch while retaining union provenance.                               |
| Empty union                                       | Embed the ground empty type.                                                         |
| Other closed, independently specializable subtree | Embed as `Data(G)`.                                                                  |
| Unsupported expression                            | Return a located failure without partial usable output.                              |

Two additional requirements are essential.

**A. Lifting decisions must be position-sensitive.**

An expression appearing as a type operand is not necessarily a type expression everywhere else. Literal contents of `exact(value)`, labels, field names and scalar selectors must never be recursively interpreted as parameters, applications or Wire constructors.

For example, a literal Ontos tuple that happens to resemble `param(0)` remains literal data.

**B. Closed embedding must remain explicit.**

Consider:

`Other<T> = not(T)`

`Root = not(Other<Person>)`

The complete expression for Root can be specialized and embedded as ground data. It need not lift Other.

However:

`Alias = Other<Person>`

preserves the top-level application and therefore requires a supported lift of Other. It must fail under the initial profile rather than silently switching to embedding.

This distinction is important: the choice is determined by the requested transformation and expression position, not by which engine route happens to be available.

In particular:

- `Data(any)` remains ground.
- `Data(not(G))` complements ground Ontos values only.
- `not(T)` is not independently liftable.
- `set-of(T)` remains unsupported.
- Recursion containing free generic parameters is not liftable.
- Closed specializations of these expressions may be explicitly embedded.

The empty-union rule also deserves a denotational test: its embedding must preserve the empty set of admissible values, not introduce a representation-specific bottom value.

### Required lifting qualification

For every admitted template D and closed ground arguments A, the derived template, with arguments explicitly embedded as Data, must describe precisely the embedded instances of the specialization of D.

This must be tested for:

- Records with optional fields.
- Nested applications and transitive dependencies.
- Singleton and empty unions.
- Closed complements and unsupported open expressions.
- Parameters occurring at different descriptor operand positions.
- Exact literals resembling declaration syntax.
- Recursive descriptors with and without free generic parameters.

Source AST equality, identical module pins and identical codec bytes are not required.

## 4. Identity, profiles, modules and native generation

The most important identity decision is to distinguish semantic compatibility from declaration identity.

A declaration identity should continue to identify an exact module digest and declaration name. It must not be replaced by structural shape equivalence.

I recommend defining a separate immutable profile identity containing at least:

- Profile identifier and version.
- Owning extension and semantic revision.
- Admitted value domain.
- Interpretation and recognition rules.
- Ground embedding and recovery laws.
- Relevant positional and encoding conventions.

A matching profile name alone is insufficient. Compatibility requires matching the relevant semantic revision and admitted-domain contract.

Derived declarations must receive distinct identities and retain:

- Original source module pins.
- Selected roots.
- Lifting version.
- Source declaration identities.
- Ordered parameter mappings.
- Derived output names.
- Provenance and dependency relationships.

The existing Horos specialization identities and module pins should not change.

### Library-owned native recipes

This is a particularly important architectural requirement.

Function, Cell and Stream are ordinary library-owned contracts, not universal runtime primitives.

Their native representations must therefore be owned by the exact library declarations and explicitly selected by the consuming build.

A native recipe should identify:

1. Source language/profile and pinned module.
2. Exact library declaration identity.
3. Target language and representation profile.
4. External generic type or type-only native AST.
5. Ordered generic parameter mapping.
6. Locked native imports and dependencies.
7. Supported arguments and representation limitations.

No name-based matching, shape-based inference, arbitrary raw source, silent fallback or last-writer-wins registration should be permitted.

Structural generation remains an explicit alternative.

For example, a generated `Page<T>` may use a selected library's native `Cell<T>` in Go, Rust, TypeScript or C++, without embedding Cell semantics in Bittype.

Crucially, a native recipe does not implement a runtime adapter or prove protocol behavior.

The proposed Adapter API should remain separate future work. Dynamic descriptions cannot manufacture statically usable native types, and a structural declaration does not prove that a runtime object implements its protocol.

Bittype's existing AST/CST pipeline, provider SPI, resolution receipts and immutable transformer configuration should remain unchanged.

## 5. Missing considerations, simpler composition and falsification criteria

There are five additional concerns worth making explicit.

**1. Capability attenuation and authority.**

A sending capability may expose less authority than its underlying endpoint. The mixed profile must not accidentally preserve or expose receiving and closing operations through a Wire face.

**2. Resource lifetime and stale references.**

A structurally valid Wire node may contain a handle whose scope has expired. Structural admission must not imply operational usability. Runtime operations must independently enforce scope, export identity and cleanup.

**3. Resource limits and adversarial nesting.**

Finite trees can still be excessively large or deeply nested. Recognition, lifting, specialization and transport need explicit node, depth and byte limits. These limits must produce distinct diagnostics rather than silently changing semantics.

**4. Protocol laws need independent specifications.**

Function, Cell and Stream require their own behavioral specifications and qualification suites. In particular, Stream needs explicit rules for concurrent demand, cancellation, terminal races and cleanup. Structural type conformance cannot establish those laws.

**5. Avoid premature generalization.**

A general capability-extension framework may eventually be useful, but the current requirements justify a Wire-specific extension. Introduce common abstractions only when another concrete extension demonstrates genuine shared semantics.

### Future qualification tests that could falsify this design

The recommended architecture should be considered unqualified if any of these tests fail:

- Mixed-profile round trips mutate or lose ground bytes.
- Structural copying changes handle aliasing.
- A malformed or sparse tuple is accepted as canonical.
- A runtime accepts an arbitrary interface implementation as an authorized capability.
- An old Horos ground checker accepts Wire syntax.
- Partial lifting changes membership or silently embeds an unsupported application.
- Literal data is mistaken for declaration syntax.
- Generic specialization substitutes parameters in labels or literal operands.
- Native generation selects an unpinned or ambiguous library representation.
- A generated type compiles but its declared adapter assumptions fail runtime conformance tests.
- Closing a hydration scope fails to invalidate its exports without incorrectly closing borrowed endpoints.

## Final recommendation

Proceed with the selected optional ownership model, subject to these qualifications:

- Ontos owns the mixed value profile, not capability authority.
- Horos owns the optional Wire type and declaration semantics, not runtime recognition.
- Bitruntime owns hydration, export identity, scope and operational admission.
- Contract libraries own protocol definitions, behavioral laws, adapters and native recipes.
- Bittype owns language-agnostic semantic tooling and explicit transformation orchestration, not another canonical type algebra.

Preserve the existing ground models, declaration identities, generic semantics and AST/CST architecture.

Document the ownership transition and supersession boundaries now. Treat profile compatibility, partial-lift soundness, native recipe selection and runtime adapter conformance as future qualification obligations.

**The decisive architectural principle is that a shared value representation does not imply shared semantic authority.** Each layer should own only the laws it can actually enforce and independently qualify.

Overall assessment: The selected direction is viable. The most consequential safeguards are the strict separation of structural admission from runtime authority, the explicit distinction between partial lifting and closed embedding, and exact library-owned native representation selection. These should be treated as design requirements rather than optional implementation refinements.