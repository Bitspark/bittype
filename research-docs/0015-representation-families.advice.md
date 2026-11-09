# Advice for consult
**Research doc:** doc_3dc21dae-8122-425a-9f63-0499b0f16d98.md
**Document:** doc_3dc21dae-8122-425a-9f63-0499b0f16d98
**Advisor:** ChatGPT, model not read from the answer page
**Session:** nsess_0150a4c5-cd83-4970-b354-bb23d7909bb0
**Run:** run_4164744c-2dd0-427b-89b3-fd2325e0a7f3
**Chat:** https://chatgpt.com/c/6af21959-cf78-83ed-b93d-362811f5cefc
**Verify:** nightfall consult verify run_4164744c-2dd0-427b-89b3-fd2325e0a7f3 --output <this file>

**I recommend keeping the architecture established in research 0014 and adopting factory-based contributions, with a few changes before freezing the design.** Each family should own a versioned representation contract, and bittype should produce one selected instance of that contract from which both native declarations and bindings are derived.

The most important correction concerns unions: **overlap does not, by itself, make a horos union impossible to bind.** The real obstruction is losing information that the round-trip laws promise to preserve. A related issue appears in sharing: preserving a compact DAG during conversion cannot establish either preservation of sender aliasing across the network or a bound on the expanded canonical encoding.

Below I answer the five questions, treating the implementation details in the attachment as the supplied baseline. Section references refer to that document.

## 1. The unit of a family and its owners

### Make the unit a representation contract, instantiated once per use

The proposed native facet, value facet, and plan entries are the right starting point. I would make their common contract explicit:

| Part of the family contract | What it must define |
|---|---|
| Applicability | Source identities and forms accepted, generic assumptions, supported targets, and explicit refusals |
| Native representation | Native type, dependencies, construction, access, discrimination, and admissible native states |
| Value meaning | Promised canonical domain, domain predicates, native equivalence, and any permitted normalization |
| Selected-use schema | A family-owned, versioned schema for the decisions recorded at each use |
| Conversion implementation | Factory or implementation reference, child adapter slots, runtime compatibility, and conversion laws |
| Operational requirements | Budget use, mutation and ownership requirements, recursion support, sharing behavior, and failures |
| Qualification | Separate evidence for type generation and binding support, tied to the actual representation and implementation versions |

The planner should select **one immutable representation instance per use**. Its native type and binding-plan entry should be two projections of that instance.

That prevents a subtle failure mode: generating a native declaration, then independently reconstructing conversion decisions from the source. Two passes can read the same declaration and still disagree about a nullable pointer, a union wrapper, or a generic specialization.

“Defined once” should mean that representation decisions have one authoritative origin. It need not mean that a Go planner and a TypeScript runtime somehow share one executable implementation. Their agreement still needs a shared contract and verification.

### Prefer A, with a contribution-owned code-generation extension

Factory arguments are the best default. A record family can supply field metadata and child adapters; a map family can supply key and value adapters; a recursive family can supply references to adapters constructed later.

However, **factory names plus plain data should not be the only permitted implementation path**.

A TypeScript runtime can access a field from its name. A Go family may instead need generated accessors, constructors, typed closures, or helper declarations. An external native type may expose construction only through a particular public API.

Provide an extension point through which the selected family can translate its recorded decisions into target-owned expressions or helpers. The binding backend should manage imports, names, references, placement, and assembly. The family’s lowering implementation should supply the family-specific fragments.

This preserves the intended ownership:

- The family chooses the representation.
- Its lowering implementation translates that choice into target syntax.
- The backend composes the resulting pieces.

A generic backend must not acquire a growing switch over records, maps, unions, and future families. Likewise, arbitrary source-code strings should not become the substitute for typed target expressions and symbol references.

A shared runtime interpreter remains a valid implementation of these factories if its handlers are registered contributions. It must execute selected decisions rather than reconstructing them.

### Put runtime authority with the representation owner

**The structural profile’s owner should own the factories’ semantic contract and qualification.** I recommend a small runtime package per target, maintained as part of that profile’s implementation.

The service library can own and maintain the backend and the common call context. It can also host the runtime package operationally. Repository placement does not need to determine semantic authority.

The important boundary is that the service backend must not independently decide what a structural integer, nullable field, or map means. Exact domain mappings should own their factories on the same terms as the structural profile.

There is a relevant failure documented by Protobuf’s maintainers: insertion points allowed plugins to depend on generated implementation details, making ordinary generator changes capable of breaking other teams’ code or violating invariants while still compiling. Their recommended replacement is separate plugins and files composed through public APIs. That supports keeping your family integration explicit and avoiding modifications to generated source as an extension mechanism. [Protocol Buffers Documentation](https://protobuf.dev/design-decisions/insertion-points-and-generic-services/?utm_source=chatgpt.com)

### Keep the shared plan about identity and composition

The plan needs a shared composition protocol, but it does not need a shared catalogue of representation semantics.

I would divide it this way:

| Shared envelope | Family-owned payload |
|---|---|
| Format version, target identity and relevant target configuration | Integer representation and admitted range |
| Checked source references and semantic-version identity | Record field layout and presence behavior |
| Root products and source-use identities | Union dispatch and construction decisions |
| Scoped binders and actual-argument environments | Map representation and key requirements |
| Selected contribution/profile identities | Refinements and domain predicate references |
| Typed dependency slots and recursive references | Family-specific construction and access details |
| Native symbols, imports, factory references, runtime requirements | Other representation-specific decisions |
| Diagnostic provenance and namespaced requirements | Family-specific diagnostic details |

The envelope should represent a **graph**, not just recursively nested payloads. Generic parameters and recursive adapter references must be explicit links. Dependencies that shared orchestration must resolve cannot remain hidden inside opaque family data.

“Opaque to the planner” must also remain compatible with “validated by the family.” Each payload needs a registered reader that verifies its schema and compatibility. Ideally, one authoring schema produces the Go payload definitions, TypeScript reader types, and validation machinery.

### Separate the version axes

A single plan-format version is insufficient. Distinguish:

1. The shared envelope and composition protocol.
2. The family payload schema.
3. The representation’s semantic identity.
4. Generated-code/runtime compatibility.
5. Source semantic identity.

A factory can change meaning without changing its argument shape. A profile can change native equivalence without changing a field name. Generated code can call a runtime export absent from an older package.

MLIR provides a close precedent: its bytecode format and dialect versions are separate, and its format-compatibility promise explicitly depends on stable dialect definitions. Protobuf likewise treats generated-code/runtime compatibility as a separate contract and prohibits pairing newly generated code with a runtime that predates its requirements. [MLIR](https://mlir.llvm.org/docs/BytecodeFormat/?utm_source=chatgpt.com)

Unknown required schemas, semantic versions, or runtime requirements should fail explicitly. Only metadata explicitly designated as nonsemantic should be ignorable.

## 2. Where native and canonical domains differ

### Distinguish four sets

The document’s laws become clearer if each selected representation distinguishes:

- \(H_T\): all canonical values conforming to the horos descriptor.
- \(V_T \subseteq H_T\): the application’s promised domain, including identified domain predicates.
- \(A_{T,p} \subseteq V_T\): the canonical values that this representation promises to cover.
- \(N_{T,p}\): admissible native values, usually a subset of the host-language type.

For a fully supported, lossless binding:

\[
A_{T,p}=V_T.
\]

For an ordinary structural representation without additional domain predicates:

\[
A_{T,p}=V_T=H_T.
\]

A bounded representation of an unbounded integer has \(A_{T,p}\subsetneq V_T\). That is **declared partial coverage**, and should be reported as such. It cannot satisfy the original full-coverage law over all of \(V_T\).

This separation prevents a family from rescuing a failed coverage claim by redefining its domain as “whatever the implementation accepts.”

Resource exhaustion belongs outside these semantic sets. A value remains valid even when the current call lacks the resources to convert it.

### The structural default should be lossless

I agree with the proposed default:

- TypeScript integers use `bigint`.
- Go integers use `*big.Int`.
- Other targets select an arbitrary-precision dependency explicitly before claiming equivalent coverage.

This is a profile policy, rather than something inherently guaranteed by being “structural.”

For Go, ordinary integer admissibility should exclude `nil`. Encoding must not modify the input, and decoding must establish the promised ownership of the returned integer. If a copy is needed, it must follow `big.Int`’s copying contract: the official documentation requires `Int.Set`; shallow struct copies are unsupported. [Go Packages](https://pkg.go.dev/math/big?utm_source=chatgpt.com)

For Content’s count, **non-negative is not sufficient to justify `number`**. A non-negative integer remains unbounded. The domain would need an additional bound such as:

\[
0\le n\le 2^{53}-1.
\]

Otherwise, keep `bigint`.

If a bounded numeric representation is selected, specify its complete admissible domain, including fractional values, infinities, and negative zero. The safe-integer interval is a useful contiguous policy; it is not the set of every individually representable integer.

### Put domain restrictions with the domain contract

Given horos’s deliberate exclusion of integer ranges, I would put a count restriction in the domain library now, with a stable identity and version that every target’s binding references.

Logos can later provide checked evidence for that predicate. It should not become a prerequisite for expressing the operational contract unless the project wants proof-carrying domain claims immediately.

There is an important distinction:

- “This application’s `Count` means an integer in this interval” defines application meaning.
- “This target’s representation handles only this interval” describes a representation restriction.

Those should produce different diagnostics. A valid application value that exceeds a target representation is not a domain violation.

### Wider native types require explicit admissibility

Using `Uint8Array` for a bounded atom or `string` for `ContentHash` is reasonable. The generated native type is then a convenient carrier, not a complete validator.

On encoding, validate native representation requirements, preserve the domain predicate, and establish ground conformance. On decoding, validate the canonical input and construct an admissible native result. Generated types alone do not discharge either direction.

The text issue also extends beyond TypeScript. Go strings can contain arbitrary bytes, so their text adapter must validate UTF-8. This is explicitly documented by the Go project. [The Go Programming Language](https://go.dev/blog/strings?utm_source=chatgpt.com)

For refinements, use horos’s interpretation of the constraint. A text-length bound must not silently become JavaScript UTF-16 length or Go byte length.

### Treat map equality as an API contract

The recommendation to use ordinary native maps only when native key equality agrees with canonical key equality is sound.

The mathematical claim needs a qualification, however: `Map<Uint8Array, V>` is not inherently incapable of a value-preserving conversion. It can represent every finite canonical map if admissible native maps exclude canonically duplicate keys and native equivalence ignores key-object identity.

Its problem is that ordinary lookup does not implement canonical key equality. That makes it a poor default map abstraction.

I would choose:

- Ordinary native maps for keys with established equality compatibility.
- A canonical-map abstraction for general keys.
- An entry sequence with explicit canonical-uniqueness rules as a simpler alternative.

For Content, `ReadonlyMap<string, ContentValue>` is a good TypeScript default. Its contract should say that insertion order is not preserved as part of the message value. An object representation can also work, but then own-property handling, special property names, construction, and extra-property behavior become part of the family contract.

A canonical-map abstraction must capture a stable key representation. Caching encoded bytes while retaining a mutable key as the authoritative value would permit the key and its index to disagree. Also, storing encoded `Uint8Array` objects in an ordinary JavaScript `Map` merely moves the reference-equality problem.

For Go, `comparable` is not sufficient evidence. A `*big.Int` is comparable as a pointer, but pointer identity is not integer equality. Native map legality and canonical key compatibility are separate requirements. [The Go Programming Language](https://go.dev/ref/spec?utm_source=chatgpt.com)

Every map encoder should detect duplicate canonical keys **after key conversion** and sort using the exact codec order specified in §3.1. Incoming noncanonical map order should be rejected where the descriptor requires a canonical map.

### Preserve absence and null through generic composition

For Go, a safe generic default is:

```go
Optional[Nullable[T]]
```

Presence and null then have separate tags. The family must define whether unused payload fields are ignored by native equivalence or rejected as invalid states.

Pointer representations remain possible, but require proof that their sentinels cannot collide. Mapping nullable `T` to `*Mapped(T)` yields `**big.Int` when `Mapped(T)` is `*big.Int`; reusing the existing pointer instead needs evidence that its `nil` state is otherwise inadmissible.

For TypeScript, keep `pending?: Pending | null`, together with runtime checks distinguishing an absent own property from a present `undefined`. `exactOptionalPropertyTypes` strengthens static checking, but the distinction remains observable at runtime and must be preserved by the adapter. [typescriptlang.org](https://www.typescriptlang.org/tsconfig/exactOptionalPropertyTypes.html?utm_source=chatgpt.com)

## 3. Unions: decide by information preservation

### Correct the empty-list example

The conclusion in §4.2 about

```text
union(list-of(int), list-of(text))
```

does not follow from horos’s stated semantics.

The canonical empty list is one value that belongs to both alternatives. It has no branch provenance. A native empty array in `bigint[] | string[]` likewise need not carry branch provenance.

A codec can define:

```text
decode(canonical empty list) = []
encode([]) = canonical empty list
```

Both round-trip laws can hold. Nonempty lists can be validated and converted according to their elements.

Horos union is a set union. JSON Schema’s distinction between `anyOf`, which permits multiple matching alternatives, and `oneOf`, which requires exactly one, is a useful comparison. [Boolean JSON Schema combination](https://json-schema.org/understanding-json-schema/reference/combining?utm_source=chatgpt.com)

### State the actual obstruction

The existing native round-trip law implies:

\[
E(n_1)=E(n_2)\ \Longrightarrow\ n_1\approx n_2.
\]

If two native values encode to the same canonical value, decoding that common value cannot recover two meanings that the contract promises to distinguish.

For example:

```text
Ints([])  → canonical empty list
Texts([]) → canonical empty list
```

cannot preserve the native variant tag if `Ints([])` and `Texts([])` are observably different under the promised equivalence.

This requires a correction to the accepted statement from 0014. **A library-owned representation cannot recover arbitrary information that the canonical contract erases.** It can instead:

- Use a common native representation without the extra distinction.
- Specify a normal form and admit only that form.
- Define the native equivalence so that the extra distinction is irrelevant.
- Change the source contract so that the distinction is actual canonical data.

The last option is what a true tagged sum does. Avro, for example, encodes a union branch index before the branch value. Adding such an index to these horos values would change their canonical contract, so it cannot be introduced merely as binding glue. [Apache Avro](https://avro.apache.org/docs/1.12.0/specification/?utm_source=chatgpt.com)

Native equivalence must be chosen honestly and in advance. It should not be weakened after discovering a failed round trip.

### Source disjointness and native distinguishability are independent

There is also an obstruction in the other direction.

Consider:

```text
union(list-of(int), seq(int))
```

A canonical labelled list and a bare tuple sequence are distinguishable here. Nevertheless, if both standalone families use native `bigint[]`, the native empty array does not say which canonical value to produce.

Each family can validly use arrays by itself. Their composition into this union needs a reconstructible native wrapper or a refusal.

Thus, source disjointness alone does not qualify a binding. The selected child representations must also compose without losing canonical distinctions.

### Use a conservative first family

I would still ship a first structural union family restricted to cases with straightforward evidence:

- Distinct canonical forms.
- Distinct exact compound labels.
- A required discriminator field with disjoint exact values.
- Null plus a type proved to exclude null.
- Compatible native tests and constructors for the selected child representations.

That is a sensible supported subset, rather than a theorem that all overlap is impossible.

Later, an overlap-supporting family can establish coherence. Useful sufficient obligations are:

\[
E_i(n)=E_j(n)
\]

where the native representations overlap, and

\[
D_i(v)\approx D_j(v)
\]

where the canonical alternatives overlap.

A fixed, documented partition can also work. The planner should not invent such a partition opportunistically.

### Divide ownership of the judgment

**Horos semantics and bittype’s source integration** should establish canonical facts: membership, disjointness, overlap witnesses, and valid canonical discriminators. Domain libraries may supply separately identified facts about their additional predicates.

**Target representation families** should establish native discrimination, construction, and compatibility between child representations.

**The planner** should link those facts and requirements without learning their family-specific meaning.

**The backend** should emit the selected strategy without discovering a discriminator or choosing a preferred branch.

A discriminator is not necessarily a membership proof. Seeing a `list` label can identify a candidate branch without proving that every element conforms.

### Keep analysis sound and deliberately incomplete

Start with a bounded analysis that handles the cases above, fixed arities, finite exact-value sets, and simple propagation through intersections.

Its outcomes should distinguish:

```text
proved(strategy/evidence)
overlap(witness)
unknown(reason)
exhausted(budget)
```

Neither `unknown` nor exhaustion establishes overlap or disjointness.

For Drafts’ intersection over a union, logical normalization may be useful as an analysis view, but it must preserve the original declaration identity used for exact mappings. Also, intersecting closed records is not generally equivalent to merging their fields: each original closed-record constraint still applies.

Serde provides a useful example of the consequences of a different policy. It explicitly tries untagged variants in order and returns the first successful deserialization. That defines an ordered decoder; it does not establish preservation of arbitrary native variant tags. Two variants that both serialize an empty collection illustrate the resulting information loss. This is a consequence of the policy, rather than a claim of a Serde implementation defect. [Serde](https://serde.rs/enum-representations.html?utm_source=chatgpt.com)

## 4. Recursion, native identity, sharing, and limits

### Resolve source binders before choosing native names

The source integration should own lexical binder identity, scope, shadowing, guardedness, and substitution.

Native naming should operate on those checked identities. A whole-body `rec` can take its declaration’s public name. A lifted inner binder needs a deterministic name derived from its identity and native scope, including collision and reserved-word handling.

Crucially, lifting an inner binder must preserve its free generic parameters. A helper declaration cannot accidentally detach a recursive body from the generic environment in which it was defined.

### Choose indirection from the selected layout graph

Indirection should be decided from native layout dependencies.

Some edges already cross a representation with finite-size indirection; other edges embed their target directly. The target implementation should identify the remaining layout cycles and choose where to break them.

The shared planner may provide graph algorithms and strongly connected components. The target owner must define the meaning of its layout edges and legal declarations.

Ownership is part of this decision. Rust’s `Box` supports recursive size requirements, while `Rc` supports shared ownership; the Rust documentation illustrates that two lists cannot share the same boxed tail under exclusive ownership. Changing between these later changes the native representation. [The Rust Programming Language](https://doc.rust-lang.org/book/ch15-01-box.html?utm_source=chatgpt.com)

Accordingly, do not freeze a Rust or C++ recursive family’s ownership representation while leaving its promised native sharing behavior undecided.

### Construct recursive adapters through explicit references

The backend should allocate placeholders or lazy references, construct the mutually dependent adapters, then publish the completed adapter graph.

That avoids requiring a child adapter to be fully constructed before its parent can exist. It also keeps adapter recursion separate from native value indirection: they are different graphs with different obligations.

For value conversion, maintain separate state for:

- Active traversal states, used to detect cycles.
- Completed conversions, used to reuse shared results.

A repeated completed value is sharing; recurrence through an active traversal state is a cycle.

Memoization must be keyed by the **selected representation instance and input identity**, including any relevant descriptor or generic environment. An object-only cache is unsound when the same object appears under two different mappings or refinements.

### Distinguish three sharing promises

There are three materially different contracts:

1. Abstract value equality ignores sharing.
2. Conversion preserves available physical DAG sharing.
3. Network transport preserves the sender’s original aliasing.

The current architecture supports the first and can support the second. It cannot generally support the third.

Let `x` and `y` be distinct objects with equal contents:

```text
n₁ = [x, x]
n₂ = [x, y]
```

They have the same canonical tree encoding. A decoder receiving those bytes cannot know which native graph the sender used.

Preserving sharing during an in-memory conversion is still useful. It can reduce allocation and repeated conversion work. It should be stated as an operational guarantee conditioned on the sharing exposed by that input representation.

Hash-consing equal decoded subtrees would also be an operational policy. With mutable native values, it can introduce aliases the sender did not have. It therefore cannot be treated as invisible.

### Content needs both graph and expanded-size accounting

The tension in §6.8 is already concrete because §3.1 states that canonical encoding and equality pay the expanded tree cost.

Define:

```text
x₀ = a small byte value
xₖ = [xₖ₋₁, xₖ₋₁]
```

There are only \(O(k)\) distinct nodes, but \(\Theta(2^k)\) encoded occurrences.

A distinct-object limit and a depth limit do not bound that expansion. Preserving the DAG during conversion merely postpones it.

Use one shared call context containing several explicit measures:

| Measure | Purpose |
|---|---|
| Actual conversion and validation work | Bounds computation performed |
| Distinct objects and allocated bytes | Bounds materialization |
| Maximum path depth | Bounds nesting and recursive processing |
| Expanded canonical occurrences or encoded bytes | Bounds downstream processing and output |

These measures need defined units and owners. Ontos or the message layer should provide the canonical encoded-size calculation; families should account for their actual work and allocations; Content should own its domain-specific limits.

Expanded size can be computed without expanding the DAG. Memoize each subtree’s size, add that size for every incoming occurrence, and saturate arithmetic at the configured limit. Similarly, record subtree height so that a memo hit on a deeper path cannot bypass the depth limit.

A single call-wide budget should mean one shared accounting context. It should not mean that all resource dimensions collapse into one ambiguous node counter. Actual validation and conversion work may both be charged, while logical message size should not double merely because two passes inspect it.

Cap’n Proto documents a closely related amplification problem and accounts for traversed pointed-to data rather than relying only on distinct stored objects. Its handling of zero-size lists also illustrates why compact storage alone cannot establish bounded processing. [capnproto.org](https://capnproto.org/encoding.html?utm_source=chatgpt.com)

### Leave the pending capture boundary with its owner

These conclusions do not require choosing the service library’s pending capture/lifetime contract.

Families can already declare what is mutable, what they copy, what they retain, and what stability they require while reading. The service owner can then define when capture completes and which read-stability guarantees apply.

Copying atom bytes establishes protection against later mutation of those original bytes. It does not, by itself, establish a coherent snapshot of an arbitrary mutable composite graph during conversion.

## 5. What is most likely to force a refactor

### Generic representations need explicit assumptions

This is the largest missing composition requirement.

A native type parameter does not provide the adapter, equality behavior, discriminator, or layout evidence needed to bind values of that type.

| Generic form | Required obligation |
|---|---|
| `Page<T>` | Supply `T`’s selected adapter along with its native type |
| `map-of(K,V)` | Establish compatible key semantics, or use a uniform canonical-map representation |
| `union(T,text)` | Establish discrimination or overlap coherence for the admitted actuals |
| Nullable `T` | Establish that the null representation cannot collide with a valid `T` |
| An inner recursive binder using `T` | Preserve the free parameter and its adapter when lifting the binder |
| An exact external generic mapping | Supply its library-owned factory arguments and requirements |

For `U<T> = union(T,text)`, qualifying `U<int>` does not qualify `U<text>`.

Each generic family must either prove its requirements for every permitted actual, declare requirements that actuals must satisfy, or defer qualification to checked applications. Explicit specialization is another option, but it must be represented as such.

Also, instance-specific binding decisions must remain consistent with the native generic declaration. A plan cannot quietly give one instantiation a different native layout while claiming it is ordinary substitution into an unchanged generic type.

This argues for carrying scoped binder identities, native type arguments, adapter arguments, and requirements together from the first plan version.

### The reported exhaustion behavior is a qualification blocker

Section 3.2 reports that the boolean conformance API can turn exhaustion into `false`, which `not` may then turn into success.

That behavior cannot underpin binding validity or disjointness evidence.

Using `diagnose` is the correct direction, but verify its internal composition. A diagnostic wrapper around an evaluator that already made an unsound logical decision would not repair it.

At minimum:

```text
not(exhausted) = exhausted
not(not_well_formed) = not_well_formed
```

Likewise, a union with no successful alternative and an exhausted alternative has not established mismatch.

Preserve distinct failures for unsupported uses, missing facets, incompatible plans/runtimes, invalid native values, structural mismatch, domain violations, declared coverage restrictions, cycles, and exhaustion. Family-specific detail can remain namespaced. Shared orchestration must never interpret an arbitrary conversion failure as a failed membership test.

### Preserve information in the less prominent families

Several smaller decisions deserve explicit treatment before their families are delivered:

- **Open records:** retain additional fields or refuse the representation. Dropping them breaks the canonical round trip.
- **`any()`:** preserve arbitrary ontos atoms and tuples. Do not reinterpret every familiar-looking label as a normalized data form.
- **Closed records:** define how wider native objects with extra properties are handled.
- **Go collections:** decide whether nil and empty collections are equivalent, and how that interacts with explicit nullable forms.
- **Bytes:** specify ownership of decoded storage, including what repeated occurrences may share.
- **Sequences and lists:** they may share a standalone native representation; their composition must preserve distinctions where both occur.

An explicitly selected validated canonical-value wrapper is also a useful representation for difficult forms. It can preserve full values while richer native representations remain unsupported. That must be a deliberate selected mapping, consistent with the existing prohibition on fallback.

### Freeze the envelope after testing difficult compositions

The proposed delivery order is reasonable for Content. It is insufficient as the order in which to validate and freeze shared interfaces.

Before calling the envelope stable, build representative vertical cases:

1. Recursive `ContentValue`, including a shared DAG, a cycle, and an expanded-size failure.
2. `Page<Outcome>` with optional nullable fields and an exact discriminator.
3. An exact external generic mapping such as `Cell<T>`.
4. A Go example requiring typed construction and a generic adapter argument.
5. A generic map whose key requirement is rejected correctly.
6. An independently registered family that requires no kernel or backend-dispatch change.

Then qualify each family through evidence appropriate to each claim:

- Compile generated types in the target configuration.
- Check both conversion laws over declared domains.
- Test boundary values, malformed input, and classified failures.
- Verify substitution, aliases, and recursive applications.
- Check source conformance independently of successful native construction.
- Verify cross-target canonical agreement.
- Demonstrate that adding an unrelated contribution preserves existing selections.

Tests that deliberately break a rule can establish that the test suite detects that failure. They do not establish complete coverage or prove the conversion laws. Keep the law-level argument alongside the executable evidence.

### Make the no-refactor promise precise

For the extension class described here, a strong promise is practical:

> Adding a family through the published identity, dependency, invocation, graph, and target-extension contracts requires changes only to that family and the consumers composing it. Existing representation identities retain their meaning.

A guarantee covering arbitrary future semantics is not something this review or a finite test suite can establish.

Separate additive availability from consumer migration. A new family can be added without changing existing ones; moving a consumer to a different native representation still changes that consumer’s contract.

The reported Rust integer limitation and absent C++ runtime should therefore remain explicit qualification boundaries. Type generation can be supported separately, but it must not imply that full binding coverage has already been established.

## Assessment of the four hypotheses

| Hypothesis | Assessment |
|---|---|
| **1. Structural profiles never narrow.** | Adopt this as the default profile’s explicit coverage promise. Domain predicates, representation restrictions, and resource exhaustion must remain separate. |
| **2. Union distinguishability belongs on the source side.** | Canonical membership facts belong with source and domain semantic owners. Native discrimination and compatibility belong with representation owners. Binding support needs both. |
| **3. Selected contributions and their arguments suffice beyond source semantics.** | Yes, provided the protocol also carries explicit graph links, scoped generic arguments, requirements, schema identities, runtime compatibility, and target-owned lowering references. |
| **4. Sharing is never part of a family’s semantic contract.** | Too broad. Sharing can be a native or operational promise. Exact sender aliasing cannot generally be preserved across the current canonical tree encoding. |

**My recommended first release is a lossless structural profile using factory-based contributions, explicit generic requirements, conservative union support, graph-based planning, and separately qualified target bindings.** Before freezing it, correct the overlap argument, establish compositional diagnostic outcomes, and exercise recursion and Go construction through the same public extension mechanism that the scalar families use.