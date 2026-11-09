# Advice for consult
**Research doc:** doc_e9c3587a-c86e-4f68-9328-436988661a2a.md
**Document:** doc_e9c3587a-c86e-4f68-9328-436988661a2a
**Advisor:** ChatGPT, model not read from the answer page
**Session:** nsess_d0e45633-9a44-4f75-ab80-04821f5fd2fe
**Run:** run_250b7afe-dcac-4972-a8b2-4d3962226188
**Chat:** https://chatgpt.com/c/6af1096c-008c-83ed-82b8-ba6239e4c83c
**Verify:** nightfall consult verify run_250b7afe-dcac-4972-a8b2-4d3962226188 --output <this file>

## Recommendation

**I would choose a binding-owned build-time backend that consumes bittype’s authoritative representation plan.** Libraries should contribute their native representation and adapter together, and the backend should emit static composition plus whatever typed glue each target needs.

That is option **B**, with a stronger interface than the representation manifest sketched in the document. Bittype should remain the authority for selecting native representations; the binding backend should consume those decisions without reconstructing them.

I would also adopt:

- **Horos/ontos-conforming message forms** for structural values.
- **Named JSON projections of abstract ground values**, composed with the native adapters.
- **Documentation metadata outside canonical horos identity**, normally authored beside the declarations.
- **Separate domain outcomes and invocation failures**, owned respectively by the domain and Function libraries.

The central principle is that **one semantic owner does not require one emitter**. It requires an authoritative decision, an explicit interface carrying that decision, and tests that establish that its consumers preserve the promised meaning.

I’m treating the repository descriptions and code excerpts in the attachment as the baseline. Section references below refer to the research document; the external comparisons are supported by primary sources.

## 1. Who should emit, and what should be emitted?

### Choose the boundary by responsibility

The two arguments in §6.3 each identify a real concern, but neither establishes where all emission must live.

The argument for A correctly identifies the need for agreement with native representation. Its overreach is the claim that another emitter must rederive bittype’s choices. It need not do so if those choices are exposed authoritatively. Conversely, putting two independently implemented passes inside bittype would not establish agreement.

The argument for B correctly identifies ownership of adapter composition and its runtime interface. Repository proximity is not itself a semantic justification, however. A binding generator should understand the adapter interface and composition facilities without acquiring Function’s protocol, Stream’s cancellation rules, or ContentHash’s validity predicate.

I would establish these responsibilities:

| Owner | Responsibility |
|---|---|
| **Horos** | Declaration semantics, checked identities and applications, descriptors, ground conformance |
| **Bittype’s planning mechanism** | Exact contribution selection, substitutions, references, and the selected representation plan |
| **Each representation implementation** | Native layout, presence rules, numeric representation, construction and access operations |
| **Binding backend** | Adapter composition, runtime imports, generated glue, and assembly of requested binding products |
| **Contract/domain library** | Its native recipe, adapter implementation, protocol or domain meaning, and associated laws |
| **Bitruntime** | Capability recognition, export/import identity, lifetime, admission, and transport |
| **Consumer build** | Selection and composition of the concrete contributions it uses |

A representation implementation may be a target’s structural profile or a library’s exact custom mapping. Both should use the same public contribution mechanism.

The binding backend can run in the same Go process and planning pipeline as bittype. Whether it is distributed as a separate module, an optional bittype product, or a revived bitlink should follow this interface decision. The dormant repository and existing platform generator do not decide the boundary.

### Publish a selected plan, not merely a table of names

A manifest containing declaration names, field names, optional encodings, and external symbols is useful but insufficient.

The interface needs to describe the representation of each relevant **use and generic instance**, including:

- Its checked source identity and actual generic arguments.
- The selected contribution, target, and representation profile.
- Native types and symbols, including runtime versus type-only imports.
- Field access and construction operations.
- Absence, presence, null, and native `undefined` rules.
- Tuple layout and calling convention.
- Union discrimination and construction.
- Numeric range and exactness requirements.
- Recursive references.
- Selected adapter factories, argument-adapter slots, and resource requirements.

These need not all be fields in a universal serialized schema. They can be typed references and target-owned expressions.

That distinction matters. Bittype’s existing design already lets targets own the meanings of native forms while shared machinery handles substitution and references (§3.6). Extend that principle: **keep target-specific operational representation information target-owned, and retain horos as the source of type semantics**. Do not introduce another canonical type algebra inside the manifest.

I would begin with a typed build-time interface. Serialize it only when a real process boundary requires that. Bittype’s charter would need to permit this selected-plan interface or equivalent product hooks; it would not need to acquire runtime adapters or protocol implementations.

A useful completeness test is:

> Can the binding backend implement the selected representation using only this plan, referenced source semantics, and the selected contributions—without guessing a representation rule?

If it must independently decide how Go optional values work or how a custom union is represented, the boundary is incomplete.

### Emit static composition and typed glue

The code-versus-descriptions choice should be made per responsibility.

I recommend generating:

- Ordinary static imports of selected factories.
- Factory calls composing child adapters.
- Ground descriptor constants.
- Operation and presentation metadata.
- Target-specific accessors, constructors, or conversion glue where necessary.

Shared structural factories may interpret a closed description. An interpreter becomes a second policy owner only when it independently chooses representation rules. Executing explicit decisions from the selected plan is ordinary implementation.

TypeScript may need mostly factory composition. Go may benefit from generated typed field access and construction. The design should permit both without forcing Go into a JavaScript-style object interpreter.

The descriptions option also needs a small executable emitter: static imports and factory calls are executable syntax. The real distinction is how much conversion logic is generated. The types-only AST limitation in §3.6 calls for a separately scoped binding emitter, with appropriate syntax and compiler checks; it does not require extending bittype into a general-purpose programming-language compiler.

### The directly transferable protobuf lesson

Protobuf’s plugin interface passes generators file descriptors, imports, generation parameters, and resolved names. That demonstrates a useful separation between compiler analysis and independent emission. Its request does **not** convey another generator’s chosen native layout, so it does not establish that source descriptors alone are enough for your two outputs to agree. Your boundary needs the additional selected representation information described above. [GitHub](https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/compiler/plugin.proto?utm_source=chatgpt.com)

## 2. What should the library contribution contract contain?

### Pair the native and adapter contributions

A library should publish a contribution for a particular declaration and supported target representation:

| Contribution element | Required meaning |
|---|---|
| **Source selector** | Exact module digest and declaration identity |
| **Representation context** | Target and supported profile |
| **Native facet** | Bittype recipe or native intent |
| **Adapter facet** | Runtime factory reference and typed argument-adapter slots |
| **Requirements** | Adapter-interface compatibility, resource facilities, supported scope |
| **Optional product facets** | Library-owned operation or presentation interpretation |

The native and adapter facets must be selected as a compatible pair. A custom native mapping without an adapter facet may still support type generation, but a request for bindings must fail explicitly. It must not silently fall back to structural conversion.

The contribution belongs under the **library owner**, outside canonical horos bytes. A library-controlled Go mapping package is a reasonable implementation given the build-time environment. It can reference the TypeScript or Go runtime factory whose import will be emitted. A declarative sidecar is also reasonable for simple contributions.

The decisive property is joint ownership and selection, not whether the contribution is written in Go or metadata. Native package paths and adapter export names should not become part of abstract declaration identity.

### Preserve identities through composition

The composition process should:

1. Resolve declarations and applications through the existing checked graph.
2. Preserve the original declaration associations through aliases.
3. Select the contribution by exact identity.
4. Bind generic slots using checked parameter identities.
5. Compose child adapters and emit the selected factory call.
6. Reject duplicate contributions, incompatible profiles, missing required facets, and unsupported applications.

Selection must happen before structural expansion destroys the association with a library declaration. Two declarations with identical bodies can legitimately have different adapters and behavioral meaning.

The mechanism may index contributions explicitly supplied by a consumer. That is different from a framework catalogue that must be edited for every new library. Consumers necessarily select dependencies; shared machinery must have no preferred Function, Stream, or Cell entries.

Adapter caching also needs care. A cache key may need the selected contribution, generic actuals, and representation profile. Runtime export reuse additionally needs the appropriate ownership scope. “Same native object” alone is not sufficient justification to reuse an export under a different adapter or lifetime owner.

### Operation discovery must use the same extension boundary

The document defines an operation as a field whose type is a library Function (§1.2). That convention must not become a Function-specific switch in the general generator.

Function’s contribution should expose its argument, result, and failure positions through a build-time product or operation facet. A consumer-selected service or presentation contribution can combine that information with structural field paths.

General traversal may know that it is visiting a record field. It should not secretly know that the first argument of a particular declaration means “request” and the second means “response.”

Likewise, a mapped Cell should not be opened and flattened into `get` and `set` operations unless Cell or an explicitly selected service projection supplies that interpretation.

### State contribution laws over explicit domains

For a type \(T\), define:

- \(V_T\): its admissible abstract values, including any library-owned semantic predicates.
- \(N_{T,p}\): admissible native values under representation profile \(p\).
- \(E\) and \(D\): encoding and decoding.

The basic data laws are:

\[
D(E(n)) \approx_N n
\quad\text{for } n \in N_{T,p}
\]

\[
E(D(v)) = v
\quad\text{for } v \in V_T
\]

Here, equality on the second line is abstract-value equality. Any normalization must be explicitly specified. The first equivalence must say whether copying, native class identity, branch tags, and sharing matter.

Every contribution should also establish:

| Law | What it establishes |
|---|---|
| **Conformance** | Produced values inhabit the declared descriptor where ground conformance applies |
| **Semantic validity** | Domain predicates survive both directions |
| **Coverage** | Every promised abstract value is representable; restrictions are explicit |
| **Capture** | Later mutation of caller-owned storage cannot change the captured message |
| **Generic composition** | Substitution, aliases, nested applications, and recursion preserve the laws |
| **Resource atomicity** | Failed preparation releases newly acquired resources without damaging existing ones |
| **Capability preservation** | Authority, identity/reuse rules, scope, and lifetime remain valid |
| **Cross-target agreement** | Different native representations agree on abstract values and protocol observations |

Do not weaken coverage to “the values this implementation happens to accept.” If the declaration admits values a native representation cannot express, either choose a faithful representation, strengthen the declared domain, or expose and reject the unsupported scope.

For capabilities, round trips require contextual behavioral equivalence, not ordinary equality. The proposed `stub(bind(f)) ≈ f` law (§3.10) needs explicit assumptions about connection state, lifetime, cancellation, and observable failures. It cannot promise local-function timing, reference identity, or transport-free execution.

Rollback also has a precise limit: it can undo newly prepared exports. It cannot undo an admitted message or arbitrary external effects. Conversion should not invoke a protocol merely to determine how a value should be encoded.

### Worked case: Function failures

I support `Function<A,B,E>`, subject to the Function library owning the complete design.

Its declared application result is:

\[
\operatorname{Outcome}(B,E)
=
\operatorname{Ok}(B)
\mid
\operatorname{DomainError}(E)
\]

That describes a completed application answer. Invocation failures occupy a separate channel.

A native API could return `Promise<Result<B,E>>` while using a distinct Function-owned failure mechanism for invocation failures. Alternatively, it could return an explicitly tagged encompassing result. The representation contribution owns that native choice.

The distinction must preserve these meanings:

| Observation | What it establishes |
|---|---|
| A received domain error | The operation produced the declared failure; it may already have changed state |
| Local non-admission | This send was not admitted, according to the runtime’s guarantee |
| Disconnect or missing reply after admission | The result is unknown |
| Malformed or undecodable reply | The caller lacks a valid outcome; execution may already have completed |
| Cancellation | Whatever cancellation guarantee the Function protocol actually supplies |

Remote protocol errors that travel in replies must appear in the Function library’s **full declared reply protocol**. An `ok | error(E)` declaration is incomplete if the peer may also send `protocol` or `limit` replies. Local refusal need not be a remote reply variant.

If native implementations throw exceptions, domain-error conversion must be explicit. Recognizing arbitrary exceptions by a familiar `code` property is insufficient. Error-class reconstruction belongs to the domain’s native adapter.

Declared failures can remove repetitive transport wiring, but they cannot generate business decisions about which failures an operation promises or what side effects precede them.

## 3. How should JSON and documentation work?

### Derive JSON from abstract ground values

For ground data, use a named projection \(P\) and compose:

\[
\operatorname{toJSON}_{T,p}
=
P \circ E_{T,p}
\]

\[
\operatorname{fromJSON}_{T,p}
=
D_{T,p} \circ P^{-1}
\]

Each library then implements its native-to-abstract-value conversion once per target representation. Each presentation format implements its projection once. Libraries avoid a separate native codec for every format.

This is the directly useful Serde lesson: a type supplies its mapping to a shared data model, while a format supplies the model-to-format mapping. Ontos already provides your common value model, so there is no need to introduce another one for this purpose. [Serde](https://serde.rs/data-model.html?utm_source=chatgpt.com)

The existing `ontos-data-json/1` is not enough: the document describes it as one-directional and unable to produce bytes or sets (§3.5). A new lossless profile needs explicit treatment of:

- Bytes.
- Unbounded integers.
- Arbitrary map keys.
- Sets.
- Tuples and compounds.
- Arbitrary ground ontos values.
- Collisions between ordinary data and reserved JSON encodings.

A total fallback for finite ground values is straightforward in principle: represent atoms as tagged byte strings and tuples as tagged arrays of recursively projected values. A friendlier profile can shorten canonical text, records, and collections while defining escapes or fallback forms.

“Two-way” should mean reversible on a specified JSON image:

\[
P^{-1}(P(v)) = v
\]

Accepted alternative spellings may normalize on the return trip, but that normalization must be specified. Integer exactness must survive parsing; passing an unbounded integer through an ordinary finite numeric representation is not a lossless implementation.

Use a **pure ground conversion** for this composition. JSON rendering should not export endpoints as an incidental consequence of calling an arbitrary effectful adapter.

### Separate presentation convenience from serialization meaning

A lossless projection does not automatically produce an excellent form editor.

A library may additionally contribute a hash input widget, a ContentValue editor, or entry hints. Those should operate on the abstract-value or projection interface, so they are reusable across native targets.

Live capabilities need a separate interaction or reference protocol. They cannot become portable JSON data merely by adding a tag. Presentation should report unsupported operations or positions explicitly; one unsupported capability-bearing operation need not prevent every ground-data operation from being presented.

Argument names, arity, flattening, and unit results also need an authoritative source. The CLI should consume the operation projection selected by the build, rather than infer those rules again from native types.

### Use one documentation sidecar with multiple authoring routes

The documentation alternatives in §6.2 are partly complementary:

- Readable-syntax comments can produce documentation metadata.
- Separately authored or localized documentation can produce the same metadata.

I recommend one sidecar format, outside canonical horos identity, keyed by:

- Exact declaration identity.
- Source member path.
- Parameter position or another defined parameter identity.
- A use-site location when service-specific wording overrides a declaration default.

A description of Function itself cannot supply the summary for `content.value.put`. Target-native field names are also unsuitable as stable documentation keys.

Define precedence between declaration defaults and use-site overlays, and diagnose stale references. Prose edits should change documentation output without changing the canonical module digest.

Protobuf provides a useful precedent: its descriptor schema carries optional source information and comments for development tools, and explicitly permits removing that information without harming runtime descriptor functionality. This supports separation of tool metadata; it does not determine horos’s identity rules. [GitHub](https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/descriptor.proto?utm_source=chatgpt.com)

Semantic constraints must retain a semantic owner. A hash pattern, numeric bound, or protocol guarantee cannot be established merely by an “entry hint.”

## 4. What should determine message forms, validation, and the pilot?

### Use descriptor-conforming structural values

I agree with the document’s leaning toward canonical horos/ontos data.

The source declaration should determine the abstract value being communicated. The native representation contribution determines how that value appears in the target language. The structural adapter connects the two.

A custom native representation does not grant permission to emit an unrelated abstract value. If ContentLink is declared as a record, its adapter should produce that record. If the compact tuple representation remains desirable, it should be an explicitly declared protocol representation or named projection.

This gives ground conformance a direct role and keeps target-language choices from changing structural message meaning.

### Resolve the service-root discrepancy explicitly

There is an architectural decision between §§3.1 and 3.8:

- The declared service is a **record of Function values**.
- Today’s `serviceBinding` exports **one capability receiving path-dispatched requests**.

Those are different abstract contracts.

For the canonical default, I would export the service as its declared record of library Function values. If a single dispatcher capability is a requirement, give that root protocol an explicit library-owned declaration and adaptation to the service API.

The general generator must not infer that adaptation because a record happens to contain callable fields. Making arguments descriptor-conforming does not by itself make the root service descriptor-conforming.

This decision should precede the Content pilot because it determines what the generated root binding means.

### Divide validation by judgment

The split should be:

| Judgment | Owner |
|---|---|
| Whether the value is ground or mixed | Ontos/wire value facilities |
| Ground structural conformance | Horos |
| Native representation validity and construction | Structural or custom adapter |
| Additional domain predicates | Domain library |
| Capability recognition, authority, and lifetime | Runtime |
| Protocol messages, state, replies, and cancellation | Contract library |

For ground inputs, establish groundness and use emitted descriptors. Prefer `diagnose` at boundaries where failure classification matters: exhaustion and mismatch are different outcomes (§3.4).

For mixed values, checking only ground leaves is insufficient. The enclosing record or tuple still needs its required fields, labels, arity, and capability positions checked. Recognizing a capability does not establish that its peer implements `E` or obeys a protocol.

Any reusable mixed structural checker needs a specified responsibility consistent with horos and ontos/wire. It should not become an independent type semantics hidden inside the binding runtime. Unsupported mixed constructs should be rejected explicitly until qualified.

Include an explicit `any()` test: a hydrated value must not be cast to a ground value and accidentally admitted through a ground checker.

Validation and conversion also need coordinated resource accounting. Giving every child a fresh whole-operation budget defeats the intended bound. Rechecking every recursive subtree can multiply traversal work. Check maximal ground regions or otherwise account for repeated work deliberately.

### Correct two proposed simplifications

**An integer declaration does not by itself remove Content’s result validator.** The declared integer is signed and unbounded, while `check.objects` currently means a non-negative integer (§§3.2, 3.9). Non-negativity still needs a refinement or domain-owned predicate. The native representation also needs range and exactness rules.

**Synchronous conversion does not by itself establish capture.** A synchronous adapter can return an atom backed by the caller’s mutable byte buffer. Later mutation can still change the message.

The required law is:

> After successful capture, later mutation of caller-owned data cannot change the captured message.

That can be implemented by copying, transferring ownership, or relying on an explicit immutability contract. The synchronous boundary determines *when* the guarantee must hold; it does not establish the guarantee.

Apply the same reasoning to provider results and decoded values. Define the supported mutation model, including whether shared memory or reentrant property access is admitted. DAG memoization may preserve efficiency, but native reference identity is a separate promise.

### What Content should qualify

A completed Content pilot should exercise all six operations and produce:

- Native types from horos.
- Exact selection of library-owned adapters.
- Descriptor-conforming messages.
- Ground validation and domain validation.
- JSON presentation.
- Operation metadata and documentation.
- Declared domain outcomes with separate invocation failures.

It should preserve Content’s existing meaningful guarantees: hash validation, non-negative counts, snapshotting, recursive limits, cycle and duplicate-key rejection, and the intended treatment of sharing (§3.9).

A first `put`/`get` slice can reduce uncertainty, but it should not be reported as qualification of all Content or the full extension boundary.

The most informative tests are:

| Test | Claim it challenges |
|---|---|
| Mutate input bytes immediately after calling; mutate provider results after completion | Capture and ownership |
| Exercise absence, null, unit, bytes, large integers, and recursive values | Representation fidelity |
| Use overlapping unions with distinguishable native alternatives | Automatic union support |
| Feed Go-produced descriptors to the TypeScript checker with valid, invalid, and exhausted cases | Cross-implementation conformance |
| Fail after one new export while retaining an existing export | Rollback and reuse scope |
| Disconnect after admission or corrupt a successful reply | Preservation of execution uncertainty |
| Deliver duplicate or late replies | Function protocol discipline |
| Change module digest while retaining declaration name and shape | Exact identity selection |
| Present ground operations beside an unsupported capability operation | Additive presentation support |

The Go-descriptor/TypeScript-checker test is particularly valuable because the note identifies that exact combination as untested (§3.4).

### Falsify additivity by freezing shared machinery

After establishing the supported structural scope, run these extension experiments with shared orchestration unchanged:

| Extension | Allowed additions | Failure evidence |
|---|---|---|
| **Second contract library** | Its declarations, contributions, adapters, tests, and consumer imports | Generator switches, default registration, bootstrap edits, or structural rediscovery of its behavior |
| **Second service** | Its contract, implementation, documentation, and dependency composition | New service paths or domain cases in shared libraries |
| **Go bindings** | Go backend and Go-specific implementations of selected libraries | Changes to TypeScript contributions or general orchestration to recognize Go-specific cases |

The second library should have distinct behavior and work through aliases and nested generic applications. A renamed Function is too weak a test.

The second service should exercise meaningful differences, including cross-service imports and union or nullable structures. If it reveals a genuinely unsupported structural form, that requires a general profile extension under AC-08; it does not justify a service-specific exception.

“Adding Go changes only Go” needs one clarification: it includes **Go-specific implementations supplied by each semantic owner**. Arbitrary Function, ContentHash, or ContentValue behavior cannot be synthesized from a native type recipe alone. The natural implementation unit is library × target.

## 5. What is missing, and which risk is underestimated?

### The largest data risk is loss of information

The document emphasizes agreement between native types and adapters. I would add a stronger question:

> Does the chosen native representation preserve enough information for the promised conversion laws to be possible?

Consider a horos union of integer values and the exact integer zero. Suppose a native tagged sum distinguishes `GeneralInteger(0)` from `ExactZero`. Both encode to the same canonical integer zero. Decoding cannot recover which native branch was used.

A complete manifest cannot solve that loss. Neither can a single emitter.

The representation must instead:

- Intentionally identify those native alternatives under the stated equivalence.
- Use distinguishable source values, such as explicitly tagged compounds.
- Supply a library-owned representation with the required information.
- Reject the unsupported mapping.

The reverse problem also exists: distinct abstract alternatives may collapse into the same native representation, making encoding ambiguous.

Serde explicitly chooses declaration-order first-success matching for untagged enums. That is a specified policy, not a consequence of union membership. Horos’s set-theoretic union semantics do not by themselves authorize importing that policy. [Serde](https://serde.rs/enum-representations.html?utm_source=chatgpt.com)

This is why **type-generation support and binding support should be qualified separately**. Intersections, complements, recursion, and custom native representations make “every valid type has an automatic adapter” too strong a promise.

### The largest operational risk is effects hidden inside conversion

Capability adaptation is not ordinary data serialization. It may allocate endpoints, reuse exports, transfer lifetime obligations, and prepare resources that outlive the call.

The ownership contract must answer:

- Who owns a callback after the request finishes?
- What does commit transfer, and at which boundary?
- Which newly prepared resources are rolled back on failure?
- Which exports may be reused, and within what scope?
- What remains live after cancellation or an unknown outcome?

A generic transaction object is useful only if these meanings are precise. Successful compilation cannot establish them.

The particularly dangerous combination is speculative union conversion with effectful adapters: trying a branch must not invoke remote behavior, and rollback must not close a previously shared export.

### “Declared once” still needs an explicit semantic boundary

`ContentHash = text` does not declare the hash pattern. The model library’s validator supplies a stronger admissible domain than the descriptor (§4.4).

That can be a sound design, but the claim should be precise: **structure is declared in horos; additional domain laws remain with the model library**. If those laws must be portable and machine-readable from the contract alone, horos needs suitable refinements or another explicitly owned semantic mechanism.

Similarly, if failures must be ground data, express that restriction in the relevant kind or contract. Do not infer it from the word “error.”

Before committing the interfaces, I would therefore write down four short contracts: the supported native representation laws, the service-root protocol, the capture/lifetime model, and the distinction between completed domain outcomes and invocation uncertainty. Then use Content, an independent second library, and a small Go interoperability slice to challenge those contracts. Those experiments will provide much stronger architectural evidence than the location of the emitter or the amount of generated code.