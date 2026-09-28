# Advice for consult
**Research doc:** doc_a89555cd-e8c1-4fb1-981b-68a17d0bfe8d.md
**Document:** doc_a89555cd-e8c1-4fb1-981b-68a17d0bfe8d
**Advisor:** ChatGPT GPT-6 Pro
**Session:** nsess_575afb85-6234-4f2d-a1aa-1a9c97156ad7
**Run:** run_727b1126-6e4d-44a9-9996-16f72909a29a
**Chat:** https://chatgpt.com/c/6ab9c357-6230-83eb-a17e-539d8fb80b03
**Verify:** nightfall consult verify run_727b1126-6e4d-44a9-9996-16f72909a29a --output <this file>
**Citations:** 30 source marker(s) linked from the answer page by the controller

# Recommendation for bittype’s value-type core

I recommend **an abstract-value language with structural type identity, closed records, labelled unions, compositional nullability and constraints, productive recursion, and first-order generics**. JSON should be a specified representation of those values—not their universal definition.

The most important decision is what “the same declaration” means after normalization. That decision must precede the encoding design. In particular, **hashing a generic constructor and its argument digests is not automatically equivalent to hashing the resulting type**. Both are defensible identity schemes, but they identify different things.

Your document correctly separates declaration identity, wire profiles, generated-code compatibility and behavior. I would make the same separation between **type meaning, JSON validation descriptions, and native memory layout**. Several of the recorded contradictions arise because those three are currently doing each other’s work. [file: research.md]

## 1. Kinds and primitives

### Recommendation and deciding principle

Admit a constructor when it expresses a useful distinction in values, has portable semantics, and cannot be expressed adequately through the existing basis. Do not admit a constructor merely because a particular wire format or target language has one.

The supplied consumer evidence makes labelled unions, recursion and bytes more urgent than entities, inheritance or family parameters: unions are being simulated with optional fields, and recursive values are escaping into `json`. [file: research.md]

I would choose this first-revision basis:

| Disposition | Constructs | Decision |
|---|---|---|
| Admit | Closed records | Unordered named fields, with presence recorded on each field. |
| Admit | Labelled unions | Exactly one named variant, carrying any admitted value type. |
| Admit | Lists and string-keyed maps | Finite collections; list order matters, map order does not. |
| Admit | Nullability and typed refinements | Type-level operations, discussed below. |
| Admit | Named recursive definitions | Subject to explicit guardedness, productivity and regularity rules. |
| Derive | Unit, enums and transparent aliases | Unit is the empty record; enums are unions of unit payloads; aliases disappear during normalization. |
| Derive | Literal types | Singleton refinements of appropriate scalar types, with exact literal semantics. |
| Defer | Entities, nominal newtypes, family parameters, open records/unions, intersections and untagged unions | None should enter accidentally through another constructor. |

Keep `enum` as convenient surface syntax, but give it no independent semantic identity. An enum and the corresponding explicitly written unit union must have the same identity and the same representation under a given representation specification.

I would omit `extends` initially. When admitted, it should mean **checked syntactic composition**, not subtyping: expand it before identity, reject overlapping member names, and give inheritance order no residual meaning.

Entities do not need a new value constructor. A future entity declaration can associate a record with a typed key projection. Composite keys can then be ordinary records rather than a special concatenated-string convention. Any contractual meaning of that key designation belongs in the entity declaration’s identity; it must not silently become a property of the underlying record.

### Keep variant labels; exclude discriminator-member names

A union needs semantic alternatives:

```text
Part =
    text(TextPart)
  | image(ImagePart)
  | count(int53)
```

The labels `text`, `image` and `count` belong to the language. A JSON property named `"type"`, `"kind"` or `"tag"` does not.

A representation may encode `image(p)` as an adjacent tag and payload, a single-entry object, or a binary discriminant. Those representations must preserve the selected alternative and its complete payload.

The deciding law is:

\[
\operatorname{decode}_{P,T}(\operatorname{encode}_{P,T}(v))=v
\]

for every supported value \(v\). This rules out flattening that identifies a scalar payload with a record containing a similarly named member.

Smithy makes this distinction explicitly: union member names identify alternatives, exactly one member is selected, and serialization is determined by a protocol. It also supports unit-valued alternatives. That is a close precedent for the boundary bittype needs. [Smithy: Aggregate types — Smithy 2.0](https://smithy.io/2.0/spec/aggregate-types.html)

### Define values abstractly, with a separate JSON projection

The bitschema contradiction requires an explicit specification change, not a reinterpretation of its existing discriminator field.

I recommend distinguishing:

```text
S       closed semantic type
P       named JSON value-projection specification
D_P(S)  closed JSON validation description for S under P
```

A discriminator-member name may appear in `D_P(S)` because that description validates a concrete JSON representation. It must not appear in `S` or its declaration digest.

Publish the JSON projection independently of any network protocol. Bitschema could own that specification’s text; bittype’s exporter implements it from the specification and shared fixtures, without importing bitschema’s code. Bitlink selects that representation or supplies an explicit translation when binding operations.

The byte-exact boundary then becomes:

\[
\operatorname{descriptionBytes}(S,P)
\]

rather than an underspecified “description of `S`.” Fixtures must identify both the description revision and projection revision.

This preserves the fixed byte-exact boundary while acknowledging that the present bitschema design mixes type structure with representation choices. [file: research.md]

### Primitive domains

On the evidence supplied, I would start with these domains:

| Type | Proposed meaning |
|---|---|
| `boolean` | Exactly two values. |
| `string` | Finite sequences of Unicode scalar values; exact sequence equality, no implicit normalization. |
| `bytes` | Finite sequences of octets. Base64 is a representation choice. |
| `int53` | Exact integers from −9,007,199,254,740,991 through +9,007,199,254,740,991. |
| `float64` | Finite binary64 values, with positive and negative zero identified; no NaN or infinities. |
| `null` | One explicit null value, distinct from unit and absence. |
| `timestamp` | A precisely specified instant domain, not “whatever the platform accepts as RFC 3339.” |
| `json` | A finite JSON-shaped value tree with explicitly fixed numeric and string domains—not arbitrary host-language values or preserved JSON text. |

For the initial `json` domain, I would use finite binary64 numbers and the same string domain as above. This deliberately does **not** preserve arbitrary-precision JSON number spellings. Duplicate object names must be rejected before a parser loses them.

The name `int53` makes the limit honest. A source spelling `integer` could transparently alias it. Narrower ranges derive from refinements. Add `int64`, `uint64` or arbitrary integers only with an explicit exact native rendering and JSON projection; never silently widen `integer`. The I-JSON specification documents precisely this interoperability boundary and recommends string representations when greater numeric precision is required. [RFC Editor: www.rfc-editor.org](https://www.rfc-editor.org/rfc/rfc7493.html)

For `timestamp`, my concrete starting proposal is a proleptic-Gregorian, POSIX-scale instant with nanosecond precision, restricted to years 0001–9999, represented internally by normalized seconds and nanoseconds. Exclude leap-second spellings and unknown-offset timestamps. Preserve an original offset only through a separate structured type. RFC 3339 alone is insufficient to settle the domain: its grammar includes offset and leap-second cases that require additional policy. [RFC Editor: www.rfc-editor.org](https://www.rfc-editor.org/rfc/rfc3339.html)

Defer `date` and `decimal` until their distinct semantics are needed and specified. Neither should be approximated by `timestamp` or `float64`.

**Essential vectors:** unit versus null; enum versus its expanded union; scalar versus record union payloads; empty bytes versus absence; integer boundaries; floating zero and overflow; Unicode supplementary characters; duplicate map keys; timestamp offset equivalence.

## 2. Presence, nullness, constraints and recursion

### Presence belongs to fields; nullability belongs to types

I would retain the predecessor’s stated principle, but repair its implementation: **a field may be absent; a value cannot be absent from itself**. The document shows that nightseam’s supposedly equivalent nullable spellings currently diverge in both checking and identity. [file: research.md] [file: research.md]

Use one semantic representation:

```text
Field = { name, required, type }
Type  includes Nullable(T)
```

No second nullable flag survives parsing.

| Field declaration | Absent | Present `null` | Present value of `T` |
|---|---:|---:|---:|
| Required `T` | Reject | Reject* | Accept |
| Optional `T` | Accept | Reject* | Accept |
| Required `Nullable<T>` | Reject | Accept | Accept |
| Optional `Nullable<T>` | Accept | Accept | Accept |

\* Unless `T` itself admits null.

Nullability should be idempotent:

\[
N(N(T))=N(T)
\]

Also specify cases such as `N(null) = null` and `N(json) = json` for the proposed JSON domain. Normalization must occur after alias resolution and generic substitution, not just during parsing.

For a first revision, I would expose only the type-level spelling. Any later field-level sugar must desugar **before** checking, description generation and identity.

JSON Schema independently separates required properties from the types of their values and explicitly distinguishes missing properties from properties containing null. That supports the semantic distinction, without requiring bittype to adopt JSON Schema’s entire composition system. [JSON Schema: JSON Schema - object](https://json-schema.org/understanding-json-schema/reference/object)

### Constraints refine types, not storage locations

A constraint should be available wherever its base type is available:

```text
Email = Refine(string, Pattern(...))

List<Email>
Map<string, Email>
Nullable<Email>
Page<Email>
```

The meaning of a nullable email is:

```text
Nullable<Refine(string, Pattern(...))>
```

The pattern constrains non-null strings. It is not a vaguely defined predicate applied to null.

Every constraint needs a typed domain, an exact comparison or matching rule, and a normalization rule. I recommend:

- Integer bounds use exact integers.
- Floating bounds are converted once, by specified round-to-nearest/ties-to-even rules, to finite binary64 values; identity records those exact values, for example as fixed-width hexadecimal bit patterns.
- Timestamp bounds compare normalized instants, not source strings.
- String lengths count Unicode scalar values; byte lengths count octets; collection lengths count elements or entries.

This is a deliberate alternative to recording exact decimal bounds while validating approximately. For a `float64` constraint, the declared bound is a floating value. Exact decimal constraints belong to a future decimal domain.

Capture declaration numeric tokens losslessly until their target domain is known. Otherwise a general JSON parser may destroy information before the checker can make the required distinction.

For `pattern`, “ECMAScript minus lookaround and backreferences” is not a complete portable specification. Specify a small grammar and matching semantics, including full-string versus substring matching, character classes, anchors and Unicode treatment. Pin it by revision. Do not delegate semantics to whatever regular-expression engine happens to ship with Go or JavaScript.

Normalize nested refinements through explicit rules, such as combining interval bounds. **Do not promise that every two predicates accepting the same values receive one identity.** That would be a much stronger equivalence relation than ordinary declaration normalization.

### Do not give global `unique` a local-validator meaning

Separate “elements of this collection are distinct” from “this field is unique across all entities held by a service.”

The former is local, but requires a specified value-equality relation. The latter is a property of external state. A validator inspecting one record cannot establish it.

For revision one, I would defer both unless a concrete scenario requires local collection uniqueness. Never silently reinterpret the predecessor’s global uniqueness annotation as local validation.

### “Absent equals empty” is not automatic validation

Preserve absence and emptiness in the core. Field-specific equivalence belongs in an explicit application normalization or interpretation rule.

Otherwise the same validator would sometimes preserve a distinction and sometimes erase it based on prose. That would recreate the current difficulty in which optional fields carry several different application meanings. The empty-bytes example is particularly important: an explicitly present empty byte string must remain present. [file: research.md]

If normalization is eventually declared in the language, its identifier and parameters become contract content. Until then, keep it outside validation.

### Recursion: separate finite values from finite Go layouts

I would reject the collections-only rule. It confuses a semantic question with a rendering limitation.

These definitions should be admissible:

```text
Node = { next?: Node }
Node = { next: Nullable<Node> }
Tree = leaf(string) | branch(List<Tree>)
```

Each has finite values. By contrast:

```text
Loop = { next: Loop }
Loop = more(Loop)
Loop = { children: List<Loop> with minimum length 1 }
```

has no finite construction.

Use explicit guardedness and structural-productivity checks. Optional fields, nullable alternatives, terminating union variants and possibly empty collections can supply terminating cases. Positive collection minimum lengths must be considered. This is not a blanket promise to solve satisfiability for every combination of refinements.

Separately, generate indirection where Go requires it—for example `Optional[*Node]`, with constructors and codecs enforcing valid states. Go’s specification permits recursive structures through pointers but forbids recursion composed solely through inline structs and arrays. That is a native-layout rule, not a reason to prohibit the abstract values. [Go: The Go Programming Language Specification - The Go Programming Language](https://go.dev/ref/spec)

Schema recursion does not make cyclic host-language object graphs valid values. The value model can remain finite trees.

For TypeScript, generate `x?: T`, `x: T | null` and `x?: T | null` under `strictNullChecks` and `exactOptionalPropertyTypes`. Explicit `undefined` should not become a hidden fourth contract state. The latter compiler option exists specifically to distinguish an absent property from one explicitly assigned `undefined`. [TypeScript: TypeScript: TSConfig Option: exactOptionalPropertyTypes](https://www.typescriptlang.org/tsconfig/exactOptionalPropertyTypes.html)

**Essential laws:** all four presence/null combinations; nullable idempotence; equivalent constrained aliases; preservation of present-empty values; accepted terminating recursive schemas; rejected nonterminating schemas; identical semantic results through alternative native layouts.

Protocol Buffers’ documentation records the consequence of failing to track presence: a present empty value may disappear during a serialization round trip. Bittype should make that loss impossible rather than merely conventional to avoid. [Protocol Buffers: Application Note: Field Presence | Protocol Buffers Documentation](https://protobuf.dev/programming-guides/field_presence/)

## 3. Generics and generic identity

### First-order, explicitly scoped type parameters

Start with type parameters only. Each parameter ranges over a value type. A generic declaration has a type-theoretic kind such as:

\[
\mathrm{Value}^n \rightarrow \mathrm{Value}
\]

but generic constructors themselves are not arguments in revision one.

Require explicit, complete application. Defer inference, bounds, higher-kinded parameters and nested generic binders. Reject unused parameters after normalization; do not accidentally introduce phantom nominal identities into an otherwise structural system.

Resolve named arguments to declaration-order slots before canonicalization. Parameter names are local binders, not identity-bearing names. `Page<T>` and an otherwise identical `Page<U>` therefore have equal constructor identities.

The theory must be extended here. Its existing holes are not already a theory of binders; the document explicitly identifies that gap. Specify scope, alpha-renaming and capture-avoiding substitution, using positional bound-variable references internally. [file: research.md]

CDDL provides a relevant modest precedent: its generic arguments are bound within the scope of the generic rule. Bittype needs similarly explicit scoping, plus stronger identity laws. [RFC Editor: www.rfc-editor.org](https://www.rfc-editor.org/rfc/rfc8610.html)

### Family parameters are a module-system feature

A family parameter buys a coherent collection of related declarations and projections such as `S.Envelope`. It also introduces module signatures, member resolution, sharing relationships, and questions about whether identity depends on the whole supplied family or only selected members.

That is materially more than a bounded type parameter. Defer it until a consumer demonstrates why several ordinary type arguments are insufficient.

Do not emulate deferred family parameters by filling generic slots with unchecked JSON. That loses the very type and identity information application is supposed to establish.

### Choose the identity of the normalized application

For the structural design recommended here:

\[
\operatorname{AppliedType}(F,A_1,\ldots,A_n)
=
N\!\left(\operatorname{body}(F)[A_1,\ldots,A_n]\right)
\]

and:

\[
\operatorname{TypeId}(F<A_1,\ldots,A_n>)
=
H\!\left(E(\operatorname{AppliedType}(F,A_1,\ldots,A_n))\right)
\]

Here `N` is the specified semantic normalization, and `E` is the canonical encoding.

A specialized named alias and the corresponding runtime application must both reach this expression. A manually written equivalent structural type should reach it too.

By contrast:

```text
H("apply", constructorDigest, argumentDigests)
```

identifies an application recipe unless additional laws make it otherwise.

For example:

```text
Maybe<T> = Nullable<T>
```

Then `Maybe<string>` and `Maybe<Nullable<string>>` normalize to the same type. Their argument digests differ, so naïve recipe hashing gives different identities.

An applicative, constructor-sensitive identity scheme is a legitimate alternative, but then expanded structural equality is not your identity relation. **Do not mix the two schemes.**

### Runtime application needs more than opaque digests

For this structural choice, a runtime constructor must retain a compact normalized template, and each argument must supply a closed type witness containing the semantic description needed for substitution and normalization.

The runtime library therefore needs graph substitution, normalization and canonical encoding. It does not need the source parser, an import resolver, a transport or a peer.

A bare digest is generally insufficient to reconstruct the structure needed by those operations. Generated packages should carry the witnesses explicitly rather than infer contract identity from native type names.

Descriptions can remain closed:

```text
open template + closed witnesses
    → closed applied semantic type
    → canonical identity
    → closed JSON description under P
```

The template is not a bitschema description. Bitschema should not own generic substitution. This resolves the apparent conflict between closed descriptions and runtime application without weakening either requirement. [file: research.md]

### Restrict recursion so applications remain finite descriptions

Require regular recursive generics in revision one: recursive occurrences preserve the recursive group’s parameter vector. Accept `Tree<T>` referring to `Tree<T>`; reject a definition that generates `Nest<List<T>>`, then `Nest<List<List<T>>>`, indefinitely.

Finite runtime **values** do not by themselves guarantee a finite closed **type description**. These are separate checks.

### Laws to state and test

The semantic suite should include:

\[
S[\mathrm{id}]=S
\qquad
(S[\sigma])[\tau]=S[\sigma\star\tau]
\]

\[
N(N(S))=N(S)
\qquad
N(S[\sigma])=N(N(S)[N(\sigma)])
\]

together with alpha-renaming invariance and simultaneous substitution. In particular, replacing a hole once with `Box<hole>` must not recursively reapply that same substitution forever.

The two project-specific laws are:

\[
\operatorname{Id}(\operatorname{specialize}(F,A))
=
\operatorname{Id}(\operatorname{runtimeApply}(F,A))
\]

\[
\operatorname{gen}(F[A])
\cong
\operatorname{instantiate}(\operatorname{gen}(F),\operatorname{gen}(A))
\]

The second equivalence should be checked through value conversions and codec behavior, not generated-source equality.

Dhall is a useful identity precedent: its semantic integrity checks hash a standardized binary representation of a normal form. The transferable lesson is **normalize before hashing**, not to adopt its evaluation language or entire normalizer. [GitHub](https://raw.githubusercontent.com/dhall-lang/dhall-lang/master/standard/binary.md)

## 4. Canonical declaration identity

### First specify the equivalence relation

I recommend **structural equality of normalized declarations**, including explicit field and variant labels.

Exclude declaration names, import aliases, filenames and inline-record positions. Include field names, variant names, presence, primitive domains and constraints. Transparent aliases disappear.

Consequently, two differently named records can have one type identity. A renamed field cannot. Moving an unchanged inline record does not change that record’s type identity, though changing its enclosing field can change the enclosing record’s identity.

A source name alone therefore does not distinguish `UserId = string` from `OrderId = string`. Distinguish them through an actual wrapper today, or introduce an explicit nominal constructor in a later revision. Do not hide nominality in directory paths.

Unison demonstrates both the benefit and the limit of name-independent identity: it separates names from content hashes, but also supplies unique types when structural identity is insufficient. Unlike Unison’s structural examples, bittype should retain field and variant labels because those labels are declared contract content. [Unison: 💡 The big idea · Unison programming language (+1 more source)](https://www.unison-lang.org/docs/the-big-idea/)

Also narrow the glossary’s “equal meanings” promise. Specify **equality under the published normalization rules**, not every possible mathematical equivalence between validators. Two differently expressed regular expressions may accept the same strings without being identical declarations.

### Identity content and ordering

For a value-type root, include its normalized reachable type graph. For future operation and contract roots, additionally include operation names, signatures, typed errors, operation mode and declared roles.

Exclude documentation, source locations, native rendering choices, wire property names, paths, carrier/profile revisions and unrelated exports.

That satisfies the ruling that imported declarations and operations must contribute by content, without repeating the predecessor’s whole-family invalidation of otherwise unrelated callables. [file: research.md] [file: research.md]

Field, variant and operation maps are unordered. Parameter slots are ordered. List values remain ordered. Source ordering can remain presentation metadata.

Reject duplicate members before constructing those maps. Do not turn malformed declarations into valid ones by silently deduplicating them.

### Recursive graphs need a canonical graph rule

Sorting JSON keys is insufficient. It does not settle recursive references, duplicate subgraphs or arbitrary node numbering.

For the proposed structural equality, I recommend:

1. Resolve references and perform the specified normalization.
2. Construct the finite rooted, labelled type graph.
3. Merge nodes equivalent under labelled graph bisimulation.
4. Number the resulting nodes by first visit from the root, using a specified order for outgoing edges.
5. Encode references by those canonical node numbers.

“Bisimulation” here means matching node content and corresponding labelled edges recursively. Merging equivalent nodes makes identity independent of whether equivalent structure was written twice, shared through an alias, or represented by equivalent recursive definitions.

The specification must define this process, not merely name it. In particular, do not attempt to compute cyclic content identities by repeatedly hashing references until the hashes “stabilize.”

### Imports pin content; families remain packaging units

Retain a family-like module as a publishing and import convenience, but not as the indivisible identity unit for every exported type.

A release manifest can map exported names to declaration identities and pin dependency manifests by digest. Names and version labels help locate content; verified digests establish the content actually used.

Updating a module’s unrelated export changes its manifest identity, but should not change the identity of an unchanged exported type.

For revision one, keep module imports acyclic. Mutually recursive declarations belong in one publishing unit. This avoids circular manifest hashes while still admitting recursive value types.

An unavailable dependency is an unresolved-input failure, not permission to hash a name-only placeholder. Likewise, accepted identity vectors must not contain import arrangements the language rejects—the predecessor currently violates that boundary. [file: research.md]

### A concrete encoding proposal

I would choose **`bittype-c14n/1`**, using RFC 8785 JCS over a fully specified normalized graph representation.

The name must commit to normalization and graph encoding, not just the JSON printer.

Define:

```text
preimage = UTF8("bittype-c14n/1") || 0x00 || JCS(normalizedGraph)
digest   = SHA-256(preimage)
identity = "bittype-c14n/1:sha256:" || lowercaseHex(digest)
```

Specify exact schemas for every node kind, materialize semantic defaults, and give numeric metadata canonical encodings that do not depend on an intermediate JSON number. Integer metadata can use canonical decimal strings; floating bounds can use normalized binary64 bit strings.

JCS supplies UTF-16-code-unit property ordering, UTF-8 output and precise escaping rules. It does not normalize Unicode strings or reorder arrays. Its specification also recommends string representations for numbers needing precision beyond binary64. Those rules must be followed directly, not approximated by either language’s default JSON serializer. [RFC Editor: RFC 8785: JSON Canonicalization Scheme (JCS) (+1 more source)](https://www.rfc-editor.org/rfc/rfc8785.html)

Deterministic CBOR is also viable, but merely saying “canonical CBOR” would leave important choices unresolved. RFC 8949 distinguishes ordering variants and requires applications to settle numeric representation choices. Given your JSON-heavy fixtures and modest scale, JCS is the simpler starting choice. [RFC Editor: RFC 8949: Concise Binary Object Representation (CBOR)](https://www.rfc-editor.org/rfc/rfc8949.html)

### Seed vectors

For illustration, suppose the proposed graph grammar encodes the primitive string type as exactly:

```json
{"nodes":[{"kind":"string"}],"root":"0"}
```

With the prefix above and no trailing newline, its proposed digest is:

```text
7ae40a3fd1297b7e39b2eb5cd12a7092542b7de660b49feef31c36a32c0012e1
```

A closed record with one required string field named `value` would be:

```json
{"nodes":[{"fields":{"value":{"required":true,"type":"1"}},"kind":"record"},{"kind":"string"}],"root":"0"}
```

Its proposed digest is:

```text
f0102f7c83c68978842034af1428d475385c4a9dd705b02a0571d53f28b7f542
```

These are **proposed seed vectors, not existing bittype identities**. I constructed the bytes explicitly and checked the hashes with two hashing tools. A transparent alias, a manually expanded record and `Box<string>` where `Box<T> = { value: T }` should all reach the second vector.

The full suite needs positive equivalence classes, negative distinctions, recursive graph cases and refusals. Suggested refusal messages include:

| Code | Proposed exact message |
|---|---|
| `duplicate_member` | `Duplicate member name "x".` |
| `unbound_parameter` | `Unbound type parameter "T".` |
| `open_description` | `Description contains an unfilled type parameter.` |
| `nonregular_recursion` | `Recursive application changes its type arguments.` |
| `import_digest_mismatch` | `Imported content does not match its pinned digest.` |

Specify error precedence and ordering, structured paths, and message rendering. Keep machine-specific filenames outside the normative diagnostic.

Finally, **test vectors do not prove universal agreement**. They provide independent evidence. The specification should establish that equivalent normal forms have equal bytes, distinct normal forms have distinct bytes, and decoding canonical bytes recovers the encoded normal form. Hash equality then additionally relies on the hash’s collision resistance.

Avro’s Parsing Canonical Form illustrates why scope matters: it deliberately retains only parsing-relevant attributes. That is appropriate for its purpose, but too narrow for a contract identity that must include constraints and operations. Protocol Buffers separately warns that deterministic serialization is not canonical serialization. [Apache Avro: Specification | Apache Avro (+1 more source)](https://avro.apache.org/docs/1.12.0/specification/)

## 5. Future operations, and what is still missing

### Reserve the right distinction: signature versus reference

Keep these separate:

```text
Signature                         not a data value
CapabilityReference<Signature>    a value referring to a callable thing
EntityKey<Entity>                 ordinary identifying data, not authority
```

This leaves room for callable values without pretending a function implementation is serializable data. A capability reference’s abstract type can be wire-independent even though its concrete token representation is not.

A future signature language should distinguish at least:

```text
RequestReply<Input, Output, ErrorUnion>
OneWay<Input>
```

A request-reply operation returning unit is **not** a one-way operation. The former has a reply outcome; the latter has no application reply through that invocation. Transport and local delivery failures remain separate from declared remote errors.

Typed errors can be ordinary labelled unions, with stable alternatives and structured payloads. Smithy provides a useful precedent for operations explicitly referring to input, output and structured error shapes. [Smithy: Service types — Smithy 2.0](https://smithy.io/2.0/spec/service-types.html)

Reverse calls then need no special transport concept in the core. An operation can accept a capability reference supplied by its caller, or a contract can declare abstract provided and required interfaces. Bitlink decides how those roles map to connections and messages.

Cap’n Proto demonstrates the value of this separation: structures can contain interface references, and interfaces are passed by reference rather than as ordinary structure data. Its concrete RPC mechanics need not be copied. [Cap'n Proto: Cap'n Proto: Schema Language](https://capnproto.org/language.html)

### Derive liveness through generic application

Compute a conservative “may contain capabilities” property through records, unions and collections. For an open generic, this property may depend on its arguments; after application it must be determined from the closed witnesses.

Do not let a declaration assert that it is plain while permitting an argument containing a capability. Conversely, do not ban all such arguments merely to avoid defining the effect calculation.

Lifetime, revocation, authorization and remote-object identity are not established by a matching type digest. A token claiming a signature is not proof that its holder has authority or that its remote implementation is faithful.

The first consumer requires errors and reverse calls, but this round explicitly does not design operations. Therefore record **`CALLABLE-REFERENCE-SEMANTICS`** as a named blocker for that slice until these distinctions and their ownership are accepted—not as a reason to encode callbacks as `json`. [file: research.md] [file: research.md]

### Additional decisions I would make explicit

**Compatibility is directional and policy-dependent.** Equal declaration identities establish an exact declared-content match under a known encoding, not behavioral equivalence or faithful generation. Different identities do not by themselves prove that communication is impossible. Preserve `unverified` for missing or unsupported evidence. Under an equality-only policy, a mismatch can be rejected as incompatible *under that policy*, without claiming a general compatibility theorem.

Closed records and unions make mixed-revision operation restrictive. Adding an optional field still introduces values an old closed reader rejects. If independently rolling upgrades are essential, decide an explicit unknown-field and unknown-variant policy before publishing revision one; do not let codecs invent it. The predecessor already supplies evidence of this failure mode. [file: research.md]

**Native mappings need conversion laws.** Mapping a declaration onto an existing Go or TypeScript type should use an explicit generation-side mapping and conversion functions, held to round-trip and validation laws. Keep that mapping outside declaration identity. The source document’s native-model mismatch is a real consumer gap, not merely an ergonomic preference. [file: research.md]

**Description hashes must have an explicit purpose.** A hash of `D_P(S)` is naturally scoped to a description revision and projection. Do not present it as a second interchangeable declaration identity. If bitschema retains a structural hash, state exactly what it erases and what comparisons it supports.

**Value hashing is a separate specification.** Decide whether an application hashes bytes, a logical value, or a selected projection of that value. A declaration canonicalizer does not answer that question. Do not derive durable application identities from incidental generated-struct JSON.

**Resource exhaustion is not semantic invalidity.** Distinguish language refusals from operational limits encountered while resolving, instantiating or validating. An implementation running out of resources must not quietly establish a different language or produce an identity from partial content.

**Unknown semantic features must fail closed.** An older implementation must reject an unknown constructor or constraint rather than omit it from validation or hashing. Keep open-ended tool annotations outside semantic content unless their meaning has been explicitly standardized.

## What I would require before closing bittype #1

The three deliverables should agree on one sequence:

```text
resolve and check
    → normalize under a published equality relation
    → form the canonical semantic graph
    → encode and identify
    → render native types or JSON descriptions
```

The written core specification should settle primitive domains, binders, recursion and the distinction between fields and values. `bittype-c14n/1` should settle normalization, graph equivalence, numbering, bytes and hash framing. The independent vectors should cover both routes through generic application, the byte-exact description boundary, and every recorded predecessor defect.

**My first priority would be to settle normalized type equality and the JSON-description boundary before implementing either generator.** Those two decisions determine whether the rest of the system can genuinely share one meaning and one identity.