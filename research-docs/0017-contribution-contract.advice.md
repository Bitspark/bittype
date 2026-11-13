# Advice for consult
**Research doc:** doc_d8e4f94a-2737-4a39-a3cb-bfd505f8f0c4.md
**Document:** doc_d8e4f94a-2737-4a39-a3cb-bfd505f8f0c4
**Advisor:** ChatGPT, model not read from the answer page
**Session:** nsess_2a0716f5-8aee-4b96-8938-a4169a2bdb8c
**Run:** run_49f8e87c-d37a-4578-ae05-65eaf57a60f5
**Chat:** https://chatgpt.com/c/6af656d1-3bbc-83ed-93c0-b737e261317c
**Verify:** nightfall consult verify run_49f8e87c-d37a-4578-ae05-65eaf57a60f5 --output <this file>

## Recommendation

Keep the central ownership model and **R2’s single decision record**, but **do not finalize the contract sketched in §6.5 yet**. It needs three substantive additions:

1. A rule must be able to obtain evidence about selected dependencies **before committing its own representation decisions**.
2. An optional field needs a **contextual field representation and conversion contract**. A native member plus an ordinary value adapter is insufficient.
3. The contract must distinguish **a selection at a source use, a shared representation description, an instantiated runtime representation, and a realized native symbol**.

I would also change the implementation order: prove the contract through the TypeScript process boundary and compiled Go lowering before treating the migration of all four targets as final.

Section references below are to research(20261113-001620).md[your research document](sandbox:/workspace/scratch/9c11a7251d63/upload/research(20261113-001620).md). I treat its quoted owner rules and accepted decisions as the premises of this review.

### Choices against your option tables

| Question | Recommendation |
|---|---|
| Contribution unit | **U1, generalized** to one independently selectable representation policy over a declared applicability domain |
| Rule result | **R2**, with dependency resolution before finalization and a frozen decision snapshot |
| Optional fields | **D3**, returning a target-owned member fragment together with a contextual field-conversion contract |
| Contribution identity | **I2, corrected** to use the mapping publisher’s namespace |
| Use entries and representation nodes | **N3** |
| Generics | **O1 with explicit application nodes and requirements**; specialization must coordinate native representation and conversion |
| Generated native symbols | **S2 with a checked realization manifest**; packaging the pair is useful, but Files should not rewrite payloads |
| Facets | **F2**, with typed, reader-owned facet contracts and explicit validation |
| Evidence before finalization | **P1 + P2 + P3**, including **executable V2**; V3 is preparation, not sufficient evidence |

## 1. What a rule returns and what it may ask

### 1.1 The contribution unit: independently selectable policy

I favor U1, but “one horos declaration or builtin form” is too narrow. It immediately needs exceptions for the default declaration recipe and optional fields.

Define the unit as:

> One independently selectable representation policy, for one target and one profile, over a declared applicability domain, including its paired facets.

That domain can be an exact declaration, a builtin form, the role “optional record field,” or declarations requiring the default native declaration policy.

This preserves the owner’s definition of a contribution as something the consumer selects independently. A library can distribute several contributions together and provide a convenience registration function. That packaging need not remove independent replacement.

“Library × target” then describes implementation ownership and distribution; it does not force every rule from that library into one indivisible selection.

For generics, separate mechanism from policy:

| Responsibility | Owner |
|---|---|
| Binder identity, scope, substitution, applications, environment links | Kernel and source integration, through their respective contracts |
| Native parameter order, omission, additional parameters, alias versus defined declaration | Declaration or exact-mapping contribution |
| Adapter parameters and requirements | Selected representation contribution |
| Instantiation and compatibility checks | Planner invoking the appropriate contracts |

The work item’s “generics family” should document this division. It should not force substitution into a contribution merely to make the family list uniform.

### 1.2 R2 needs a staged planning protocol

The current sketch returns the final decision record and only then resolves its slots. That is inadequate whenever a decision depends on a selected operand.

Use this lifecycle:

1. Inspect the checked subject and declare dependencies.
2. Resolve selected dependency references and requested evidence.
3. Finalize the representation decisions.
4. Freeze and validate the decision snapshot.
5. Project native intent and publish the plan from that snapshot.

**One selection does not imply one callback invocation.** The kernel can retain one selected contribution while advancing its planning operation through these stages.

One possible abstract interface is:

```text
Advance(checked request, source view, resolved answers)
    -> Need(additional typed demands)
     | Accept(decision draft)
     | Refuse(structured refusal)

Project(frozen decision, typed dependency views, symbol views)
    -> typed native result
```

This is a protocol sketch, not a proposed Go API. Its important property is that `Accept` happens after the evidence needed to decide is available. A bounded work queue can implement it without exposing an arbitrary solver.

Register the projector on the contribution. **Do not return a fresh projector closure from each decision call.** Such a closure could capture the source reader or an unstored choice, defeating the claimed authority of the decision record.

The projector should receive no source reader and perform no new selection.

#### What belongs in the decision record?

Record every **per-selection choice** affecting representation or conversion: member names, presence strategy, selected construction strategy, operand-dependent layout, and applicable configuration.

You do not need a payload field for every invariant of a representation’s semantic version. If a particular optional representation always ignores its unused value when absent, that can be fixed meaning of that semantic version. Both implementations must obey it. A payload field is needed when the behavior varies between selections.

This distinction answers one of the objections to R2. A projector may implement a recorded policy; it must not choose an unrecorded policy.

R2 nevertheless does not prove agreement between arbitrary Go and TypeScript implementations. Its advantage is an authoritative origin, restricted inputs, and inspectable decisions. The accepted phrase “cannot disagree” should mean that there is no sanctioned independent selection or reconstruction point. Taken as a claim that incorrect projector or factory code is impossible, it conflicts with AC-05’s warning that selection does not establish the truth of claims.

### 1.3 Operand requirements need typed evidence

A parent should request a typed capability or fact about **the selected dependency**, rather than inspect another family’s payload.

Examples include:

- The operand’s admitted native domain excludes the sentinel proposed for null.
- The selected key representation supplies equality and hashing consistent with the required source-key semantics.
- A branch supplies the discriminator or construction interface required by a union representation.

The evidence contract belongs to the relevant target or product vocabulary. The kernel should not acquire cases for nullability, comparable keys, or every future domain property.

The kernel can enforce that evidence is associated with the correct:

- Selected representation and semantic revision.
- Target, profile, and relevant configuration.
- Generic environment.
- Requirement kind and version.
- Source and dependency scope.

It cannot establish a domain assertion merely because a contribution supplied it. The evidence provider and its qualification must establish that assertion.

Also distinguish **false**, **unknown**, and **deferred under a binder requirement**. Absence of proof that an operand excludes null is not proof that it does. A selected nullable policy may choose a recorded wrapper strategy when proof is unavailable, or refuse the application. It must not silently replace the operand’s selected representation.

Rust’s compiler documentation provides a useful structural comparison: it separates selection, nested obligations, and fulfillment, and represents obligations supplied by a generic parameter environment. Borrow that distinction; bittype does not need to reproduce Rust’s full trait solver. [Rust Compiler Development Guide](https://rustc-dev-guide.rust-lang.org/traits/resolution.html?utm_source=chatgpt.com)

### 1.4 Recursive references and evidence cycles are different

The request protocol must distinguish:

- A reference to another representation, which may be recursive.
- Evidence that must be available before a representation decision can finish.

Reserve representation references before filling recursive bodies. Otherwise an ordinary recursive type becomes recursive callback evaluation.

An evidence cycle requires a separate policy. The first implementation can explicitly refuse cyclic evidence requirements while supporting recursive representation graphs. Supporting such evidence later would require a declared evaluation discipline—for example, a bounded monotone computation—rather than accidental recursion through callbacks.

Put this distinction in the contract now. It affects scheduling, termination, and failure meaning.

### 1.5 Refusals should be structured outcomes

At minimum, distinguish these cases:

| Outcome | Meaning |
|---|---|
| Unsupported | The selected contribution recognizes the case but cannot yet implement or check it |
| Refused by design | The case is outside that representation’s promised domain |
| Requirement unsatisfied | A selected dependency fails a declared requirement |
| Requirement unresolved | Required evidence is unavailable |
| Mechanism failure | Missing or ambiguous contribution, missing facet, invalid reference, incompatible scope |
| Implementation failure | A bug or unexpected execution error |

Use a small stable classification plus an owner-qualified reason code and structured details. Preserve the source subject, selected contribution, target/profile, dependency path, and relevant requirement.

An unresolved generic requirement may become an explicit binder obligation where the declaration contract supports that. A concrete application must discharge it.

**A selected exact mapping’s refusal must never trigger structural fallback.** Applicability and refusal are not a mechanism for trying contributions in registration order.

### 1.6 Optional fields: D3, with a stronger answer than “Member”

The document correctly identifies two possible extension domains. The general owner rule alone does not prove that presence representation must be independently replaceable. D1 could be legitimate if presence belongs wholly to the record policy.

I recommend explicitly choosing the other domain:

> The standard record representation delegates per-field optional presence through a published field-layout and conversion protocol.

Given that decision, D3 is appropriate.

The source integration should issue or validate a subject containing the parent record occurrence, field identity, operand use, and provenance. The kernel must preserve its scope guarantees. Accepting an arbitrary selector string and arbitrary operand would weaken the present checked-child-use boundary.

D2 is not inherently dishonest: a source integration may expose derived semantic subjects. The useful distinction is that this subject represents a **field position**, not a horos type. D3 makes that distinction explicit.

#### Native syntax is only half the answer

A TypeScript optional field cannot be converted correctly from its retrieved value alone. Missing and present-with-`undefined` require observing the containing object.

Therefore the selected field representation needs a contextual protocol that can:

- Test presence.
- Obtain a present value.
- Delegate that value to the selected operand adapter.
- Construct or omit the native field according to the recorded decision.

The record factory invokes this protocol. It does not interpret the optional family’s private payload.

The native part should use a **target-owned member-fragment vocabulary**, shared by the record and optional contributions. This avoids a dependency on the record contribution’s private implementation.

A fragment can accommodate a field pair as well as a single wrapped field. But state the limit: a presence bitmap shared across many fields requires record-wide coordination. Either that coordination belongs to the published host protocol, or such a representation requires replacing the record policy. A promise to replace arbitrary presence layouts cannot be implemented by a single-member interface.

This also means **not every representation contribution supplies an ordinary standalone value adapter**. Contextual field conversion is an admitted product shape that the contract needs to express now.

## 2. Identities, generics, and native declarations

### 2.1 Contribution identity: I2 with publisher ownership

Use:

```text
ContributionID  = (PublisherID, PublisherLocalKey)
ContributionRef = (ContributionID, RepresentationSemanticRevision)
```

The publisher is the owner of the mapping contribution.

It is not automatically:

- The owner of the source declaration.
- The module pin of the source being mapped.
- The profile whose promise the contribution claims to satisfy.

Those substitutions fail for third-party mappings and multiple implementations of the same profile.

A controlled package or domain namespace is practical. A namespace string alone does not authenticate its owner; consumer composition must bind it to the selected publisher and manifest. A trusted composition lock is sufficient for the stated problem. There is no need to invent a global registry or signature infrastructure merely to prevent accidental collisions.

The publisher-local key distinguishes variants from the same owner. Conflicting definitions of the same semantic reference must fail, and overlapping selection candidates still require explicit composition.

I would reject I1 in its proposed form. It lacks ownership and combines several meanings in one string. I3 remains useful for integrity and caching, but a schema digest cannot establish representation semantics.

Keep the five accepted axes explicit:

| Axis | What it identifies |
|---|---|
| Envelope protocol | How the shared graph is represented |
| Payload schema | How a family’s recorded decisions are typed and encoded |
| Representation semantics | What the selected representation means |
| Generated-code/runtime compatibility | Which generated code and runtime facilities can cooperate |
| Source semantics | Which source interpretation the checked references use |

Implementation artifact digests are useful provenance and qualification inputs; they do not replace these axes.

For `Service`’s recognition of structural records, provide a distinct **family contract reference** that spans targets. Its root facet should recognize that declared contract and supported revision, rather than parse a target-specific contribution name.

### 2.2 Preserve both use selections and representation nodes

Choose N3.

A use entry records an individually inspectable selection with its source-use or derived-demand identity and provenance. It refers to a representation node containing the shared decisions and dependencies.

That preserves “one immutable representation instance per use” while allowing sharing of its semantic description. Clarify the terminology in the accepted pages: the runtime memo must not use the unique identity of that per-use record.

The N2 tuple proposed in §6.4 is insufficient as a representation-equality key. It omits potentially meaningful differences in:

- Recorded decisions.
- Target configuration.
- Selected dependency representations.
- Native constructors or nominal declaration identity.
- Generic conventions and relevant evidence.

Two representations can have the same horos descriptor and produce different native classes. Two adapters can have the same native type and different conversion meaning.

Define identity over the selected semantic representation, excluding diagnostic provenance. Do not require the kernel to prove arbitrary semantic equivalence between different implementations.

Also, **N2 does not logically require O2**. Open representation nodes can be shared. The difficulty arises only if a plan node’s identity is required to be the complete runtime memo identity.

### 2.3 Open generics should be primary

Choose O1 with explicit application nodes.

`Page<T>` should have a complete open representation description. `Page<Person>` supplies an environment connecting its binders to source arguments, native arguments, adapter arguments, and evidence.

A free `T` needs a checked scoped binder reference. An optional field needs a checked contextual subject. Neither needs an invented ground descriptor.

Native and adapter argument conventions must be separate. Do not zip two positional lists and assume correspondence.

For example, an exact native mapping might erase `T` while conversion still depends on `T`’s adapter. Conversely, an additional native parameter may come from a representation requirement and have no direct source-parameter counterpart.

#### O2 does not repair argument-dependent native layout

This is a decisive constraint.

If `Page<T>` is emitted as one ordinary generic native declaration, its layout is fixed. Learning something about `Person` at `Page<Person>` cannot silently change that declaration.

An operand-sensitive family inside an open generic must use one of these approaches:

1. A representation valid uniformly for every admitted argument.
2. A binder requirement guaranteeing that the chosen representation is valid.
3. An explicit native generic convention capable of expressing the variation.
4. A specialization mechanism that changes native representation and conversion together.

Monomorphizing only the plan does not solve this problem. A runtime adapter environment does not solve it either.

If application-specific specialization belongs to the promised extension domain, the contract must already be able to name its application, selection, and resulting native representation. A first implementation can refuse unsupported cases, but that refusal must not conceal an input model incapable of expressing the promised specialization.

### 2.4 Runtime memo identity is an applied representation identity

Conceptually:

```text
MemoKey = (ConversionDirection, AppliedRepresentationToken, InputIdentity)
```

Separate encode and decode memo tables can make direction implicit.

The applied token identifies the selected semantic representation after applying its relevant argument representations and evidence. It is not a source-use ID, plan entry number, contribution ID, source descriptor alone, or freshly constructed adapter object.

For two fields containing `Page<Person>`:

- Their use entries differ.
- Their applied representation token should agree when their selected semantics agree.
- The same input object should then convert once for that representation.

If the selected adapter for `Person` changes, the token must change even when its native type does not.

Likewise, the token for `list<Person>` should depend on the selected list representation and selected element representation, rather than on which enclosing declaration’s binder happened to supply `Person`.

The runtime can issue opaque interned tokens from a checked linked graph. A universal global hashing scheme is unnecessary. Recursive graphs need reserved references and a seal step, or an equivalent treatment of recursive components; naïvely hashing every child before its parent cannot handle cycles.

Finally, do not require a value-memo token for every plan node. A contextual optional-field fragment may have graph identity and provenance without denoting a standalone value conversion.

### 2.5 Native symbols: S2 with a checked realization manifest

The placement contradiction is real. A plan produced before placement cannot contain the final location of a declaration whose placement is deliberately deferred. The “three sufficient inputs” statement needs clarification or amendment.

Use S2, with an explicit relationship between:

1. A serializable logical symbol identity issued for the plan.
2. The handle issued by a particular sealed native output.
3. The target-specific realization produced by Files.

A realization manifest should identify:

```text
plan artifact
plan content digest
generation run / sealed output
target configuration
placement
logical symbol -> realized native symbol
```

The logical-symbol-to-output-handle correspondence is established during generation. It must not be reconstructed from names.

Keep content identity and output identity distinct. Identical runs may have equal semantic digests while still rejecting each other’s output-local handles, as the naming design requires.

A mismatch is a **linking/input error**, distinct from a representation plan that omitted a decision.

The realized symbol should provide structured location, export, and accessibility facts. An import statement or local alias depends on where the backend places its glue; the target’s symbol/import resolver should derive those from the realization.

Do not turn the current equality of naming-intent spelling and realized spelling into a permanent API assumption. Read the realization now, even when it returns the same spelling. That preserves the intended place for a later escaping rule.

#### Contribution-owned declarations

For `BittypeOptional`, the optional contribution declares an owned logical key, such as its contribution reference plus `"optional-helper"` and relevant declaration parameters. The planner issues the actual symbol and shares that declaration according to the ownership/key contract.

Its chosen spelling is not the deduplication key.

If another contribution needs whichever shared helper implementation composition selects, it makes a demand. A genuinely private helper remains private to its owner. Coincidentally identical names do not establish shared ownership.

You may package the immutable plan and checked realization into one backend input artifact. That achieves the useful part of S4 without making Files understand or rewrite family payloads. The additional realization information remains a real dependency and should be documented honestly.

## 3. What the kernel carries and what it guarantees

### 3.1 Choose F2 as a typed extension contract

F2 is the right logical model. An open set of facet kinds should have reader-owned contracts, not merely string keys and arbitrary values.

A facet kind needs:

- A stable qualified identity and contract version.
- An in-process typed construction or issuance API.
- An appropriate validator.
- Declared reference categories and compatibility requirements.
- A transport descriptor where the facet crosses a process boundary.

The selected contribution’s manifest declares its facets and their implementations.

Resolving a Node factory by identity does not turn this into F3. The implementation can be loaded remotely while remaining a member of the one selected, owner-approved bundle. An unrelated package should not acquire permission to attach behavior merely by repeating someone else’s contribution ID.

MLIR’s interfaces are a useful comparison: generic transformations operate through interfaces instead of dialect-specific cases; external models have explicit registration; and promised interfaces can fail when their implementation was not registered. Its documentation also warns about ownership arrangements that permit duplicate or divergent implementations. [MLIR](https://mlir.llvm.org/docs/Interfaces/?utm_source=chatgpt.com)

### 3.2 Opaque does not mean unchecked

The kernel can preserve substantial guarantees without understanding family semantics:

| Guarantee | Evidence that should exercise it |
|---|---|
| Required facets belong to the selected bundle | Remove a required facet; bindings fail while an admitted types-only request still succeeds |
| Facet kind, target, and compatibility match | Substitute a wrong issuer, target, or facet contract version; reject |
| References have the correct graph, category, and scope | Inject another plan’s reference, a member in a type slot, or an escaped binder; reject |
| Evidence belongs to the selected dependency | Substitute evidence from another environment or semantic revision; reject |
| The payload snapshot is immutable | Mutate the submitter’s original maps/slices; the stored plan and native output remain unchanged |
| The exact payload passed schema validation | Submit invalid data or a receipt for another snapshot; fail before backend execution |
| Unknown required versions fail explicitly | Exercise the version axes independently |
| Extension requires no shared case | Add a materially different family and product facet with shared packages held fixed |

These are structural, authority, scope, and association guarantees.

They do not prove that a claimed equality relation is correct, that the projector uses the intended field, or that the TypeScript factory obeys the conversion laws. Qualification must test those claims.

Similarly, an opaque carrier cannot inspect arbitrary implementation code and prove that it hides no dependency. The enforceable contract is that orchestration resolves dependencies only through declared typed references; product validators and qualification check that implementations obey that boundary.

### 3.3 Freeze the record before either projection consumes it

Your question about typed Go values versus encoded payloads matters.

A sound arrangement is:

1. The family produces its generated typed decision value.
2. The Horos/Ontos integration encodes, validates, and freezes a canonical snapshot.
3. Native projection receives a typed view derived from that snapshot.
4. Serialization publishes that same snapshot.

The kernel can carry an opaque payload handle and validation result without importing Horos or Ontos. Validation must be bound to the exact schema and snapshot. “This Go type has a validator” is insufficient.

A carefully frozen typed representation is another possible implementation, but it must establish that the encoded value is the same decision the native projection consumed.

The invariant to preserve is **one immutable decision value**, not one particular storage technique.

### 3.4 Lowering needs references that reflection conceals

Before finalization, the lowering environment must be able to express:

- The native type of a selected dependency.
- Scoped native binders and adapter arguments.
- Generic requirements and their evidence.
- Generated declarations, external symbols, and owned helpers.
- Generated members, associated with their containing native type.
- Typed access and construction operations.
- Lowering-owned helper declarations and dependencies.
- The runtime call-context value required by emitted conversions.

Member names can remain payload data. A typed lowering operation also needs to know which containing type and member it is accessing. A free string alone does not provide that association.

This is why V2 is necessary. Reflection can recover information that a typed emitter cannot obtain from the proposed contract.

A small Go value-expression vocabulary and emitter in the completeness test backend is sufficient for this proof. It need not become a production backend or enter bittype’s types-only native vocabulary. It must compile and run generated construction, member access, an owned helper, and a generic adapter argument.

## 4. Evidence before calling the contract final

### 4.1 The proposed stabilization partition is too broad

“Contract final, every envelope field and link unstable” is not a viable boundary. The contribution contract and envelope share identities, dependency categories, generic environments, requirements, and symbol references.

AC-10 and an unstable envelope can coexist if the distinction is:

| Must have its intended semantic shape | May remain provisional or evolve under versioning |
|---|---|
| Identity meanings and ownership | Wire table layout and encoding details |
| Request/result categories and scope | Caches and performance strategies |
| Generic conventions and requirements | Additional admitted family coverage |
| Recursive references and refusal meaning | Family payload schema revisions |
| Facet validation and association | Additional product implementations |
| Native symbol and realization relationship | Diagnostic presentation |
| Inputs required by lowering | Lowering optimizations |

Calling the transport unstable does not excuse changing the meaning of an admitted identity later.

MLIR’s bytecode documentation makes a related distinction: overall encoding versions and dialect versions are separate, and transport compatibility does not by itself preserve changed dialect meaning. That supports keeping your version axes separate; it does not prove the proposed semantic contract is complete. [MLIR](https://mlir.llvm.org/docs/BytecodeFormat/?utm_source=chatgpt.com)

### 4.2 Require P3 and executable V2 before the final migration

I recommend this order:

1. **P1:** specify the semantic distinctions, ownership, typed demands, refusal classes, and reference scope.
2. **P3:** build TypeScript end to end through the actual service-library factory, including serialization and realization.
3. **V2:** compile and execute Go lowering with typed construction and a generic adapter argument.
4. Exercise thin executable versions of all six accepted vertical cases.
5. **P2:** with shared code held fixed, add a materially different independently authored contribution.
6. Finalize the semantic contract and migrate all 21 rule packages.
7. Publish the qualified family schemas, vectors, compatibility declarations, and protocol support.

The prototype families need not become production deliveries before the migration. Their job is to expose missing capabilities while the interface is still legitimately experimental.

A second unconditional scalar is weak P2 evidence. An optional-field fragment of a different shape, an exact generic with a nontrivial argument convention, or a family requiring operand evidence would be materially stronger.

Existing Rust and C++ native support can migrate with binding support explicitly absent. Do not manufacture binding qualification to make the target matrix look uniform.

### 4.3 The six cases test the contribution contract too

| Accepted vertical case | What it probes in this contract |
|---|---|
| Recursive `ContentValue`, shared DAG, cycle, expanded-size failure | Recursive references, adapter construction, memo identity, call-context propagation, resource accounting |
| `Page<Outcome>` with optional nullable fields and exact discriminator | Contextual field conversion, operand evidence, presence/null separation, discriminator decisions |
| Exact external `Cell<T>` | Independent ownership, factory lookup, generic conventions, native/adapter argument separation |
| Typed Go construction with a generic adapter argument | Member/type/helper references and lowering |
| Rejected generic map key | Obligations, evidence association, deferred checks, structured refusal |
| Independent family without shared changes | Actual extensibility |

These cannot be confined to “envelope tests” while the producer contract is frozen earlier. They exercise information the contributions must supply.

Also add focused identity and linking cases: repeated `Page<Person>` uses sharing a runtime token, equal native types with different adapters remaining distinct, an erased native argument retained for conversion, and one semantic plan working under two placements.

### 4.4 Interpret “changes once” as a bounded compatibility commitment

A defensible commitment is:

> One migration to the demonstrated extension contract for the agreed extension domain, before independently owned consumers rely on it.

That is consistent with AC-10.

A promise that no future requirement can ever require a shared change is not something this consultation or any finite test suite can establish. AC-08 already defines the process for a genuinely necessary change.

Write down the extension domain. In particular, decide whether it includes new products, contextual representations, cross-field layout coordination, application-specific specialization, and additional contribution-authoring formats. Otherwise the phrase “a new family changes nothing shared” has no testable boundary.

## 5. Additional issues to address

### 5.1 The source-semantics input is an unproven process boundary

The document explicitly says that the backend needs source semantics and that no serialized form is defined (§1.2). Serializing the plan alone therefore does not demonstrate the accepted backend architecture.

P3 must include either a checked reachable semantic snapshot or a pinned reader protocol. Its source identities must match those recorded in the plan. The backend may read source meaning; it must not reconstruct representation choices.

Protocol Buffers offers a useful concrete precedent: its plugin request carries generation roots, the transitive descriptor set, parameters, and compiler version, with fully qualified type references. The relevant lesson is to transport the semantic inputs explicitly. [GitHub](https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/compiler/plugin.proto?utm_source=chatgpt.com)

### 5.2 Target configuration needs precise meaning

Record configuration that affects selection, native admissibility, or generated-code compatibility.

For TypeScript, `exactOptionalPropertyTypes` changes whether an optional property implicitly admits assignment of `undefined`; it does not change JavaScript’s distinction between absence and a present property. [typescriptlang.org](https://www.typescriptlang.org/tsconfig/exactOptionalPropertyTypes.html?utm_source=chatgpt.com)

Consequently, the representation contract must state what conversion does with present `undefined`, regardless of compiler checking. Exercise absence, present null where admitted, present value, and present `undefined` separately.

### 5.3 Runtime memoization must preserve resource accounting

Planning recursion and detecting cycles in input objects are separate concerns.

For the shared-DAG and expanded-size case, a memo hit must not make the resource cost of another occurrence disappear. If the limit concerns expanded encoded size, repeated occurrences still consume that budget even when conversion work is reused.

The shared runtime contract needs to specify memo states, cycle behavior, and resource accounting. A successful recursive plan does not establish any of those properties.

### 5.4 Resolve an actual contribution-authoring route

A library owned in a TypeScript repository can initially publish a companion Go contribution package plus its TypeScript factory and manifest. Repository ownership does not require every implementation to use the repository’s primary language.

If declarative sidecars are the promised route, one must participate in the prototype. Plan serialization alone does not establish that an independent library can author a contribution.

Avoid promising both routes without testing either.

### 5.5 Payload bootstrapping should not distort the contract

The limitation to six currently supported families is a bootstrap issue, not a reason to encode unsuitable payload semantics.

Choose the appropriate payload schema, then make its generation reproducible through a pinned bootstrap artifact or the required contribution support. Do not introduce a second schema language.

Adding a semantic integer field later is normally a family schema/version change. Discovering that a payload has been hiding an orchestration reference is a contract defect. Keep those cases distinct.

### 5.6 Qualification needs behavioral and structural evidence

The proposed mutant in §2.4—building the plan entry in another pass—will not necessarily fail a behavioral test if that pass happens to compute identical decisions.

Enforce the intended boundary through APIs and dependencies: projectors cannot access the source reader, and conversion consumers cannot perform selection. Then use mutants that actually disagree on presence, layout, slot identity, or generic environments to test semantic detection.

This is the distinction between preventing an unauthorized decision point and detecting an incorrect implementation.

For §6.7, I would prioritize **operand evidence, generic conventions, identity, owned native symbols, target configuration, and the real authoring/process boundary** before finalization. Refusal classes and multi-profile composition also need explicit contracts. Payload evolution is more readily isolated once those contracts are sound.

**The acceptance condition for #196 should therefore be stronger than “all mappings implement the new result type.”** It should require a checked, frozen decision to survive native projection, dependency composition, serialization, independent factory resolution, typed lowering, and native-symbol realization—with the shared mechanism held fixed while a materially different contribution is added. That is the evidence that would justify making the migration final.