# Research: The value-type core of a wire-independent contract language

**ID:** 0001
**Date:** 28 September 2026
**Status:** draft
**Owner:** the coding agent working bittype #1 for the project lead
**Issue:** https://github.com/Bitspark/bittype/issues/1 (coordinated under https://github.com/Bitspark/bittype/issues/5)

## Question

We are designing a new **contract language**, from its **value-type core** outward:

- A contract language declares what data looks like and which operations exist.
- The value-type core is the part that declares data: records, choices, enumerations, aliases, optional and nullable fields, generics, and imports.

The language will be named **bittype**. It must be **wire-independent**: a declaration has to mean the same thing whatever network protocol, storage format or code generator later consumes it. Every declaration also needs a **canonical identity**: a digest computed from a precisely specified encoding of its content, which the same implementation must reproduce in several programming languages.

The language replaces an older one, **nightseam**'s declaration language, which several real services depend on. We treat the old language as evidence to learn from, not a specification to copy.

We need advice on:

1. The **minimal set of kinds** the core should admit, and which to exclude or defer.
2. How to model **absence versus null** as two separate facts.
3. How **generic parameters and generic application** should work, so that compile-time specialization and run-time application reach the same identity.
4. **What a canonical declaration identity contains**, and how its encoding is specified and versioned.
5. How **abstract operations and callable signatures** can later extend the core without naming any wire protocol.

**What the answer is for.** The design issue (bittype #1) closes on three deliverables, reviewed against a written contract theory:

- a written specification of the value-type core;
- a named, versioned canonical encoding;
- a set of shared, independently written test vectors.

After that come native Go and TypeScript type generation and identity libraries (bittype #3), a separate validator repository whose descriptions must match bittype's output byte for byte (bitschema #2), a generator that maps operations onto a protocol (bitlink #1), and a first real consumer (bittree #53). The project lead decides; this consultation informs that decision.

**Why it matters.**

- **The language is a long-lived public contract.** A published language revision is immutable, so a wrong choice in the core means a new revision and a migration for every consumer.
- **A wrong identity scheme makes two equal declarations look different, or two different ones look equal.** Identity decides whether a generated client and server agree they speak the same contract.
- **A wrong generics design breaks identity equality.** Generating a specialized type ahead of time and applying a generic at run time would then produce different identities for the same type.
- **Five services** wait on this language to leave nightseam.

## Context

### System Overview

**The family.** A group of small repositories, each owning one decision, together replaces a monolithic predecessor. The ones that matter here:

| Repository | Owns | State |
| --- | --- | --- |
| **bittype** | The contract language: syntax and semantics of value types, abstract operations, callable signatures, generics and imports. Also canonical declaration identity, native type rendering per programming language, and small runtime libraries for presence and identity. | Chartered, documents only. **This consultation designs its core.** |
| **bitschema** | **Descriptions** (the runtime, machine-checkable form of a type) and a validator per language that checks a JSON value against one. | Design document only. |
| **bitlink** | The mapping from abstract operations onto a concrete wire protocol (paths, message representation), and the generator of client and server adapters. | Documents only. |
| **bitwire / bitruntime** | The wire contract and its Go/TypeScript implementation. | Released. |
| **bittheory** | A written contract theory with laws that the others are checked against. | Prose, no stable law identifiers yet. |
| **nightseam** | The predecessor. It holds a declaration language, a code generator and a runtime in one repository. | Frozen at v0.6.0 as a compatibility baseline. |

**Implementation languages.** The first targets are **Go and TypeScript**. Every rule must be implementable, with identical results, in both, and later in others.

**The dependency rule** (bitwire decision 0010, quoted):

> "Value types, abstract operations and callable signatures belong to the language. Paths, the profile and message representation belong to Bitlink's binding."

> "bittype, bitschema → nothing in the project / Bitlink generator → bittype, bitschema"

So bittype depends on nothing else in the family. bitschema depends on nothing either, and in particular not on bittype. The generator depends on both.

**Glossary.** Terms are defined here before they are used.

| Term | Meaning |
| --- | --- |
| **Declaration** | Text a person writes, today JSON files, that declares named types and, later, operations. |
| **Family** | nightseam's unit of declarations: one directory, versioned and imported as a whole. Whether the new language keeps this unit is open. |
| **Value type** | A type whose values are plain data (no live capabilities): primitives, records, unions, lists, maps. |
| **Record** | A named set of fields, each with a name and a type. |
| **Union** | Exactly one of several named variants, each possibly carrying a payload; a *tagged* union says in the data which variant it is. |
| **Enum** | A closed set of names with no payload. |
| **Alias** | A second name for a type expression. It is *transparent* if interchangeable with its target, *nominal* if distinct. |
| **Entity** | A record with a declared key field that identifies application data (a durable ID, not a network capability). |
| **Absence / presence** | Whether a field appears in an encoded object at all. |
| **Nullness** | Whether a present field may hold the explicit value `null`. |
| **Generic parameter / application** | A placeholder in a declaration (`Page<T>`), and filling it (`Page<Part>`). |
| **Specialization vs run-time application** | Generating a separate, concrete type for `Page<Part>` ahead of time, versus keeping `Page<T>` generic in generated code and filling `T` in the running program. |
| **Canonical encoding** | A specified, deterministic byte serialization of a declaration's meaning, so equal meanings give equal bytes. |
| **Declaration identity / digest** | A hash (e.g. SHA-256) of the canonical encoding, naming a declaration's content. |
| **Vectors** | Test cases shared by every implementation: declarations with their exact canonical bytes and digests, and declarations that must be refused, with the exact refusal. |
| **Native rendering** | The Go or TypeScript types generated from a declaration. |
| **Description** | bitschema's runtime JSON form of a type, which a validator checks values against. |
| **Wire, profile, adapter** | The network side: how operations become messages (the profile), and the generated code that carries calls (adapters). All of it is outside the language. |
| **Callable / reverse call** | An operation signature treated as a type; a call from server to client, the reverse of the usual direction. |
| **Revision** | A published, immutable version of the language. Any change to what is accepted, rejected or meant is a new revision. |

### Architectural Context

**The rules the language must satisfy** are quoted from bittype's charter, its agent instructions, bitwire decision 0010 and issue #1:

- **Wire independence.** *"Test each declaration with: 'would it make sense, unchanged, with a different adapter and no Bitwire?'"* The language is also meant for storage, UI and command-line generators that have no protocol at all.
- **Four things kept distinct:** *"declaration identity; profile revision; generated-code compatibility; behavior the declaration does not represent. A digest names declared content, not behavior."*
- **One identity implementation per language:** *"Hold every identity library to the shared canonicalization vectors. Two implementations of one encoding in one language are a defect."* Decision 0010 adds that identity must support *"generic application at run time without the parser, the transports or the peer"*: a small runtime library, not the compiler.
- **Revisions:** *"A published language revision is immutable. A change to what is accepted, rejected or meant is a new revision. The canonical encoding is versioned by name. Any change to identity content is an explicit transition, never a side effect of packaging."* Before 1.0 there is no compatibility promise.
- **Independent evidence:** *"Compile-time specialization and run-time generic application must reach the same identity. Tests of the generic-application laws at the semantic level, not only on generated text."* Issue #1 adds that vectors must be written *"rather than using one implementation as its own oracle."*
- **The boundary with bitschema:** *"the descriptions bittype emits must equal the ones bitschema validates against, with neither repository importing the other."* The boundary is kept as byte-exact shared fixtures.
- **Design, not migration:** *"This is a new language, designed rather than migrated. Nightseam's language and the family design are evidence to examine, not specifications to copy."*
- **The first consumer's slice** must cover *"absence and null, errors, reverse calls and generic identity."* Issue #1 is explicit that these needs are *"requirements to consider, not permission to smuggle paths, carriers or profiles into this language."* If the core does not settle how callables extend it, a named design blocker is recorded instead.

**The contract theory** (bittheory, which the design is reviewed against) states:

- *A contract is `C = (S, B)`*: a shape `S` together with a behavioral predicate `B` over behaviors of that shape. *"A shape specifies available sorts, fields, constructors, operations, arguments, results, and roles."* A native declaration *"alone need not express `B`"*.
- **Shapes** are defined as `Shapes(O) = μX. Option(O) × FiniteMap(Key, X)`, where `μX.` means the least fixed point (an inductive type), `O` stands for opaque values, `Option(O)` is an optional value, and `FiniteMap(Key, X)` is a finite map from keys to sub-shapes. *"Child order has no shape meaning. A node with neither a value nor children is a valid, present shape; it is different from a missing child."*
- **On null:** *"`null` is an explicit empty slot; an omitted key differs from a key present with `null`."* And: *"Optionality is an explicit instantiation `T = Option(O)`, not an absent component."*
- **Generics are whole-shape holes:** `Shapes(O, G) = μX. Generic(G) + (Option(O) × FiniteMap(Key, X))`. *"The current language has no binders; introducing binders later would require additional scoping and capture-avoidance rules."*
- **Substitution laws.** The theory writes `S[σ]` for a shape `S` with substitution `σ` applied (holes replaced), and `σ ⋆ τ` for the composition of two substitutions.
  - S1: `Generic(g)[σ] = σ(g)`.
  - S2: substitution distributes over children.
  - S3: the identity substitution changes nothing.
  - S4: `(S[σ])[τ] = S[σ ⋆ τ]`.

  Substitution is simultaneous: *"substituting `g0 := Box[g0]` once into `g0` produces `Box[g0]`, rather than expanding forever."*
- **The generation square**, law S8: *"fill, then transform; or transform the template and every argument, then fill."* Here "transform" means code generation, and the two routes must agree up to an equivalence. *"Exact generated-source equality is a stronger requirement. A specialized generator and a generic template may produce different but equivalent code."*
- **On digests:** *"a declaration digest names content, not behavior"*. *"Equal digests alone do not prove artifact faithfulness, behavioral equivalence, or interoperability."* An injective key encoding *"must refuse keys outside that image rather than normalize them."*
- **What the theory leaves out.** It has **no value-type vocabulary**: no records, unions, enums, aliases, constraints, imports or subtyping. Parameters have no sorts or bounds. The theory assigns *"Canonical declaration and generic-application laws"* to bittype to state and test, and does not state them itself.

### The predecessor: nightseam's declaration language (v0.6.0), as evidence

**The tiers.** nightseam split a family's declarations into three files:
- `model.json`, the types;
- `protocol.json`, the operations, errors and parameters;
- `live.json`, callable types.

*"a declaration refers to its own tier or a lower one, never a higher one."* The protocol file must declare `"profile": "nightseam.duplex/1"`, which **names the wire protocol inside the language**. The new language must not do that.

**A real declaration file** (a test family, descriptions omitted):

```json
{"nightseam": 2, "types": {
  "Option": {"kind":"union","parameters":[{"name":"T"}],"tag":"kind",
    "variants":{"some":"T","none":{"kind":"record","fields":[]}}},
  "Part": {"kind":"union","tag":"type","variants":{"text":"TextPart",
    "image":{"kind":"record","fields":[
      {"name":"url","type":"string","pattern":"^https://[A-Za-z0-9.-]+/[^ ]*$"},
      {"name":"alt","type":{"nullable":"string"},"required":false}]},
    "count":"integer"}},
  "RichPart": {"kind":"union","extends":["Part"],"tag":"type",
    "variants":{"table":{"kind":"record","fields":[{"name":"rows","type":{"array":{"nullable":"string"}}}]}}},
  "Page": {"kind":"record","parameters":[{"name":"T"}],"fields":[
    {"name":"items","type":{"array":"T"}},
    {"name":"next","type":{"nullable":"string"},"required":false}]},
  "Parts": {"kind":"alias","type":{"apply":"Page","with":{"T":"Part"}}}
}}
```

**Kinds:** `entity`, `record` (may `extends` other records, and may be `open` to unknown fields), `enum` (string values only), `alias`, `union` (a `tag` plus `variants`, each variant a type or an inline record), and `callable`.

**Type expressions:**
- primitives `string`, `boolean`, `integer`, `number`, `timestamp`, `json`;
- `{"array":T}` and `{"map":T}` (string keys only);
- `{"nullable":T}`;
- `{"literal":"text"}`;
- `{"ref":"User"}`, an entity key;
- `{"apply":X,"with":{…}}`;
- inline record shapes.

There is **no bytes type, no numeric width and no decimal type**. The meaning of each primitive lives in the validators, not the language:
- `integer` is limited to JavaScript's safe range (at most 2^53 − 1), rendered as Go `int64` and TypeScript `number`.
- `number` is a finite 64-bit float.
- `timestamp` is an RFC 3339 string.
- `json` accepts any JSON value.

**Presence and null**, from nightseam's own decision record:

> "`{"nullable": T}` is a type expression, and a field's `nullable: true` is its sugar … `required` stays where it is, on the member, because presence is a fact of a member and not of a type: a member may be absent from an object, and a value cannot be absent from itself."

The generated code wraps each fact in a small struct:

```go
// Go (generated): a field not required is Optional[T]; nullable is Nullable[T]; both nest.
type Page[T any] struct {
	Items []T                                        `json:"items"`
	Next  runtime.Optional[runtime.Nullable[string]] `json:"next,omitzero"`
}
// runtime: type Optional[T any] struct { Value T; Present bool }
//          type Nullable[T any] struct { Value T; Null bool }   // zero value is NON-null
```

```ts
// TypeScript (generated)
export interface Page<T = unknown> { "items": Array<T>; "next"?: string | null; }
export type Part = { "type": "count"; "value": number } | { "type": "image"; "value": PartImage }
                 | { "type": "text"; "value": TextPart };
```

**Unions are adjacently tagged** on the wire, as `{"type": "image", "value": {...}}`, after a flat encoding lost information. A JSON payload `3` and a record `{"value": 3}` became indistinguishable (nightseam#146). In Go a union renders as a struct with one pointer per variant: *"construction can select zero or several alternatives. The generated codec enforces exactly one; the Go type alone cannot."* **Inline shapes are named by their path**: moving an unchanged inline record renamed its generated type in both languages.

**Generics** have two sorts of parameter:

> "a parameter is of one of two **sorts** … **A family parameter** — it has `of`. It is filled by a **family** … A type is drawn *through* it, `S.Envelope` … **A type parameter** — it has no `of`. … A family parameter is **not** a type parameter with a bound … a family is not a type — nothing is an instance of it."

- Type parameters have no bounds, there are no higher-kinded arguments, and there is no inference.
- Application is always explicit, `{"apply": X, "with": {...}}`.
- Unused parameters are refused.
- The identity of an application is `{"apply":"worker/Handler","arguments":[E1,E2]}`, with arguments in declared order.
- nightseam held one law in Go, TypeScript and on the wire: `gen(bind(C, F)) ≅ gen(C)[F]`. Generating the applied declaration equals applying the generated generic. It was checkable only *"because the generic declaration survives into the package rather than being compiled away."*

**Imports** name another family directory, with no version or digest. Import cycles are refused. A lower-case qualifier names a family (`identity.User`); an upper-camel one names a parameter (`S.Payload`).

**Constraints** are `min`, `max`, `length`, `pattern` and `unique`. They are allowed only on record fields, and only when the field's type is a bare primitive. So neither `pattern` on `{"nullable":"string"}`, nor on an alias `Email = string`, nor on array elements is possible.

**Canonical declaration identity.** The digest is the lowercase hex SHA-256 of a canonical JSON graph:

- **Structure:** the graph is `{"definitions": {...}, "root": ..., "version": 1}`. Definitions are keyed by paths `family/Type`, `family/$server` and `family/$client`.
- **Included:** *"Every field has `name`, `type`, `required`, `nullable`, and `unique`, including the boolean defaults."*
- **Excluded:** *"Error descriptions are documentation and are excluded."*
- **Ordering:** *"All object keys are sorted by their UTF-8 byte sequence … error codes, enum values … are deduplicated and sorted … Field, parameter, inheritance, and application-argument order is meaningful."*
- **Aliases:** *"Pure aliases disappear into their normalized target expression."*
- **Numbers:** numeric constraints become exact decimal strings (`1.00` → `"1"`, `100` → `"1e2"`).
- **Escaping:** the bytes are Go's `encoding/json` output, including HTML escaping of `<`, `>` and `&` as `<` etc. That is **not** RFC 8785, the JSON Canonicalization Scheme.

Real canonical bytes of a one-field record (abridged):

```json
{"definitions":{"same":{"client":{"ref":"same/$client"},"errors":[],"kind":"family","parameters":[],
 "server":{"ref":"same/$server"},"types":{"Payload":{"ref":"same/Payload"}}}, …,
 "same/Payload":{"fields":[{"name":"value","nullable":false,"required":true,
 "type":{"primitive":"string"},"unique":false}],"kind":"record","open":false,"parameters":[]}},
 "root":{"ref":"same"},"version":1}
```

nightseam implemented this canonicalizer three times: in the generator, in the Go runtime (705 lines) and in the TypeScript runtime (481 lines). All three were held to 76 shared cases. A peer could ask for the server's digest with an `identity.check` request.

**Recorded defects and regrets:**

- **The "sugar" is not equivalent.** `{"type":"string","nullable":true}` and `{"type":{"nullable":"string"}}` generate identical Go, yet have **different canonical bytes and different digests** (`48a35247…` versus `83bddc4a…`). The loader never desugars. This was reproduced with the released tool.
- **Recursion is refused.** A cycle through array, map or nullable is refused (`"Value schema cycle reaches Node. [cyclic_type]"`), while inline shapes are exempt. Consumers declare recursive data as `json` instead.
- **The identity was once too narrow** (nightseam#343). The digest hashed only types, so changing a method's result from `string` to `integer` left it unchanged. The operator's ruling: *"hash the declaration, canonically rendered … never a backend-specific wire rendering … operations and imported declarations included by content (not by name)."*
- **Unverified identity counted as a match** (nightseam#720). When a peer answered "method not found" to `identity.check`, the client proceeded as if the identity had matched. The rule now: *"Distinguish compatible, incompatible and unverified outcomes."*
- **One-way operations could not be declared** (nightseam#724). The method schema required a `result`, so a fire-and-forget operation was declared as a method, contradicting the runtime.
- **The only compatibility test is digest equality** (nightseam#670): *"which revisions are compatible is a consumer's policy."* Adding an optional field is a new digest, and closed records make even that change breaking for strict readers.
- **Values have no canonical encoding** (nightseam#669). A lab hashed the JSON of a Go struct, and adding a field changed every digest.
- **The identity covers forms the checker refuses.** Its vectors include mutual imports, which the language refuses. Callables' digests include their whole declaring family, so any unrelated edit changes every callable's identity.
- **Constraint semantics differ from their identity.** min/max are compared as float64, although identity keeps exact decimals. A timestamp `min` is a JSON number compared *lexicographically* against the RFC 3339 string.

### The runtime form: bitschema's description basis (design only)

bitschema owns "the runtime form of a type". Its admission rule: *"A constructor is admitted only with a scenario that no smaller basis passes, and that scenario is a row in `tables/validator.json`."* Anything else is *derived*, i.e. expressible in what is admitted, or excluded. Its current proposal:

| Constructor | Disposition |
| --- | --- |
| Primitives `string`, `boolean`, `integer`, `number`, `timestamp`, `json` | Admitted, a closed set. `bytes`, `date` and `decimal` are open. Numeric width and range semantics are not stated anywhere. |
| Record | Admitted. Members carry two flags: **`required` (default true) and `nullable` (default false).** *"JSON has two distinct absences … Collapsing them into one union member makes a validator unable to tell a consumer which of the two it received."* |
| Discriminated union | Admitted. *"a union is a discriminator name and a set of variants."* |
| Enum | Derived, as a union whose variants carry nothing. Whether it keeps its own spelling is open. |
| Array, map | Admitted. |
| Named references | Admitted. Recursion is derived through names, and *"a cycle through non-optional, non-collection members is refused at load."* |
| Foreign | Admitted: a type from another separately versioned family. |
| Opaque | Admitted: *"JSON carried through unexamined — … what fills a generic slot whose family is chosen elsewhere."* |
| Reference | Admitted: a token claiming a declared interface. *"an arrow is **not a value** here."* |
| Interface | Admitted: named members with request, result and errors. |
| Constraints | Admitted per kind: `min`, `max`, `length`, `pattern`, `unique`, each *"only for the kinds that give it meaning"*. |
| Entity and key reference | Derived. An entity is metadata; a key reference validates as the key's primitive. |
| `extends` | Derived: *"Flattened by bittype before a description exists."* |
| Intersection, untagged union, kinded parameters | **Excluded.** |

- **No parameters:** *"a description is the *runtime* form and is always closed — a parameter is filled where a rendering is instantiated, and what may fill it is bittype's check."*
- **Identity is split:** *"the structural hash is bitschema's; anything beyond structure — identity, scope, provenance — belongs to bittype's own hash … Two descriptions of the same shape share a hash even where the types mean different things, and that is intended."*

### Earlier family design and its study of nightseam

A family design document (September 2026, now design input only) proposed a first tier with these kinds:
- *entity*: *"a record with identity. `key` names the field"*;
- *record*: may `extends` and be `open`;
- *enum*: *"a closed set of string values"*;
- *alias*.

It kept presence and nullness as two field facts: *"A field that is not required is absent, not null."* Generics were **family-level**: *"a parameter is bound where a rendering is instantiated."* It left open where `unique` belongs, composite keys, the callable-reference syntax, *"the exact primitive basis"*, and `date`, `decimal` and `bytes`.

A later study of nightseam recommended:
- **Recursion through collections:** *"no value contains itself except through a collection, which the empty collection always satisfies."*
- **Avoiding the nullable trap** in Go. A self-reference through `{"nullable":"T"}` renders as `type Nullable[T any] struct { Value T; Null bool }`, *"an infinite-size Go type that does not compile … The same defect disqualifies bitschema's variant, which would make any non-required field a guard."* Presence wrappers are value structs, so recursion through an optional field produces an infinitely sized Go type.
- **Continuing to leave out:**
  - a `mu`/`var` recursion constructor;
  - type-parameter bounds (*"so that a declaration cannot be live and say otherwise"*);
  - nominal or branded aliases;
  - open unions;
  - structural subtyping;
  - `derived`/`readOnly`/`writeOnly`;
  - `default`;
  - an open set of kinds.
- **Keeping nightseam's strength:** *"One shipped declaration, interpreted by one validator per language … held to a single table case by case — including the exact words of every refusal."*

It also recorded that nightseam aliases are transparent Go aliases (`type UserId = string`), so `UserId` and `OrderId` are interchangeable, and that *"Nightseam computes no digest of a family at all"*. That last finding predates the canonical graph described above.

### What the consumers actually declare

The five production consumers are bittree, repo-tool, nightforge's concur and ulab-runner, nighthall, and the original bitsystem. The counts come from a script over their declaration files.

| | bittree | repo-tool | concur | ulab-runner | nighthall | bitsystem |
| --- | --- | --- | --- | --- | --- | --- |
| Records / enums / unions / aliases | 44 / 6 / 4 / 2 | 23 / 2 / 0 / 0 | 39 / 11 / 1 / 0 | 22 / 6 / 0 / 0 | 114 / 14 / 0 / 2 | 48 / 10 / 1 / 2 |
| Fields / optional / nullable | 151 / 57 / **0** | 87 / 0 / 0 | 203 / 73 / 0 | 64 / 11 / **3** | 447 / 114 / 0 | 270 / 135 / 0 |
| Constraints | 12 (`min`/`max` only) | 0 | 38 `min` | 17 | 14 `min` | 0 |
| Generics, entities, `extends`, imports | none | none | none | none | 4 imports | none |
| Reverse calls (client methods) | 0 | 0 | 0 | 0, but 11 callables passed as values | 8 | 3 |

- **Null is rare, and meaningful when used.** ulab-runner's `holder` field is always present and `null` means *"nobody holds control"*.
- **Optional fields are common**, and their meaning varies by field. In bittree, *"omitted and empty mean the same bytes"* for one payload. For edits, *"empty `text` means absence … explicitly empty `bytes` remains present"*, and both bindings must round-trip `{"op":"set","bytes":""}` exactly.
- **The most common workaround** is a record of flat optional fields whose presence depends on another field's value, stated only in prose:
  - bittree's `Status` (*"`result` is present exactly when finished"*);
  - bittree's `Segment` (*"Exactly one member"* of four optional fields);
  - bittree's `Difference`;
  - nighthall's `Frame` (*"Present only for cut retention"*), which is validated by hand: *"ValidateFrame adds cross-field invariants that the declaration language does not express."*

  These are really tagged unions of records.
- **`json` is the escape hatch** for recursive data (bittree's tree `Node`) and for untagged unions (bittree's `Edit.to`).
- **Errors** are code lists with no typed payload. bittree carries its error details as untyped bytes. ulab-runner's callables declare their own error codes.
- **Generics appear only in an experimental probe:** a record with four type parameters, and a generic callable applied inside a record. Its measured cost: *"generic parameters flowing through generated signatures plus one codec per native representation."*
- **bittree measured these gaps:**
  - recursive declarations refused;
  - no `bytes` primitive;
  - no way to map a declared type onto an existing native type;
  - Go field representations (`int64`, presence wrappers) that do not convert to its native model;
  - generated TypeScript that fails under the strict compiler option `exactOptionalPropertyTypes`.

### Data flow of a declaration

```text
author writes declarations (bittype's surface syntax)
  → bittype parses, resolves imports, checks            (refusals: exact codes and words, shared vectors)
  → canonical encoding → declaration digest              (one library per language, shared vectors)
  → native types (Go, TypeScript) + generation manifest  (bittype)
  → descriptions (closed; parameters filled)             (must equal bitschema's fixtures byte for byte)
  → operations mapped onto the wire, adapters generated  (bitlink; outside the language)
  → at run time: identity compared between peers         (compatible / incompatible / unverified)
```

A generic type can reach the running program in two ways, and both must reach one identity:
- **specialized:** `Page<Part>` generated ahead of time as its own type;
- **applied at run time:** a generic `Page<T>` in the generated package, filled with `Part` in the running program.

### Contradictions in our own documents

These are unresolved, and the expert's view on them is welcome:

1. **Nullness:** bitschema's design calls member flags *"Nightseam's choice, kept"*. nightseam at v0.6.0 made nullness a type constructor `{"nullable":T}` with the flag as sugar, and the two produce different identities.
2. **Recursion:** bitschema lets an optional member break a cycle. The nightseam study allows only collections, because of the Go value-struct problem.
3. **Parameters in descriptions:** bitschema says descriptions are always closed, with no parameters. The validator work item asks to resolve *"generic substitution"*, and run-time generic application must still reach one identity.
4. **The owner of the concepts beyond value types:** the earlier design placed interfaces, families and imports in a system-declaration layer that decision 0010 did not create. It gives bittype *"abstract operations, callable signatures, generics and imports"*. Whether "family" survives as a unit of versioning and import is open.
5. **What there is to audit:** issue #1 asks to *"audit what Nightseam's canonical declaration graph contained"*. The earlier study found no family digest; the graph described above came later.
6. **`json` versus opaque** in bitschema, and **`unique`**: is it a property of one value (checkable by a validator) or of a collection of entities?
7. **Aliases** are transparent in nightseam, and the study excludes nominal ones. The specification must say whether an alias is visible to identity.
8. **Still undefined:** numeric domains, the exact form of `timestamp`, `bytes`/`date`/`decimal`, composite keys, whether `enum` keeps its own spelling, and nightseam's `literal` and inline forms.

### Constraints

**Fixed:**
- Wire independence.
- bittype depends on nothing in the family.
- **Neither bittype nor bitschema imports the other**; their boundary is byte-exact fixtures.
- One canonical-identity library per language, held to shared vectors that are written independently.
- Identity names declared content, not behavior.
- Language revisions are immutable, and the encoding is versioned by name.
- Specialization and run-time application reach one identity.
- Go and TypeScript first, and the rules must port to other languages.
- The language is designed, not migrated.
- Unresolved semantics are named blockers, not invented defaults.

**Flexible:**
- the surface syntax, which need not be JSON;
- the set of kinds and primitives;
- how presence and null are expressed;
- the generic system, including whether family-level parameters exist at all;
- the canonical encoding format and what it includes;
- nominal or structural identity for named types;
- whether imports pin content;
- how much of the first consumers' needs the first revision admits.

**Scale.** Families range from tens to a few hundred types (nighthall: 114 records, 447 fields). Correctness, portability and stable identity matter far more than performance.

### What We've Considered

These are directions, not decisions:

- **Kinds: a sum-of-products core.** Records, plus tagged unions whose variants may be records, would replace the flat-optional-fields workaround.
  - Enum as a union without payloads.
  - Transparent aliases that are invisible to identity.
  - An entity as a record plus key metadata, not a kind of its own.
  - Extension (`extends`) flattened before identity is computed.
  - Unresolved: how the union *tag* stays out of the language, since the data carrier is bitlink's.
- **Presence and null as two orthogonal member facts**, not a nullable type constructor. That would avoid nightseam's two spellings with two digests. But nullable collection elements (`array of nullable string`) may still need the type form, and the Go rendering must not produce infinitely sized types under recursion.
- **Recursion only through collections**, or through an indirection that the rendering represents as a pointer. Which one is open.
- **Generics.** Type parameters only in the first revision, with family-level parameters deferred. Explicit application. The identity of an application is the constructor's identity plus its ordered arguments' identities. Generic declarations survive into generated packages so that run-time application is possible.
- **Identity.** Hash a canonical encoding of the resolved declaration graph:
  - aliases resolved;
  - defaults written explicitly;
  - documentation excluded;
  - imports included by content;
  - an established canonicalization (RFC 8785 JSON, or deterministic CBOR) rather than one language library's output;
  - the encoding named and versioned;
  - one digest per named declaration rather than per family, so an unrelated edit does not change everything.

  Unresolved: nominal versus structural identity.
- **Callables and operations** as a later extension. Signatures with typed errors, one-way operations as a declared form, and no paths or profile in the language. A reference to a callable is distinct from an entity key.

## Questions for the Expert

1. **The core's kinds and primitives.** How would you choose the smallest set of kinds and primitives for a wire-independent value language that must render natively into Go and TypeScript? Consider sum-of-products versus flags, tagged unions with record variants, enums, aliases and entities, numeric domains, and `bytes` and `timestamp`. Which constructs would you admit in the first revision, derive, or deliberately exclude, and what principle would you use to decide?
2. **Absence and nullness.** How would you model *absent* and *null* as two facts so that they compose cleanly with unions, collections, generics and recursion, render idiomatically in Go and TypeScript, and give one canonical form, with no two spellings producing two identities? How should field-specific meanings such as "absent equals empty" be expressed, if at all?
3. **Generics and generic identity.** How would you design type parameters and their application so that generating `Page<Part>` ahead of time and applying `Page<T>` at run time reach the same identity? Consider kinds or sorts, bounds, whether family-level parameters are worth their complexity, and how application relates to the closed runtime descriptions. What laws would you state and test at the semantic level?
4. **Canonical declaration identity.** What should a declaration's canonical identity contain and exclude? Consider documentation, defaults, alias resolution, imported content, field order, nominal names versus structure, and the granularity of per-declaration versus per-family digests. How would you specify, name and version the encoding, and write the refusal vectors, so that independent implementations in several languages provably agree?
5. **Open.** Looking at the whole picture — the predecessor's defects, the consumers' workarounds, the theory, the split between language, descriptions and wire binding, and the later extension to callables, typed errors, one-way operations and reverse calls without naming a protocol — what would you do differently from the directions above, and what are we not seeing?
