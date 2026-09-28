# Research: The value-type core of a wire-independent contract language

**ID:** 0001
**Date:** 28 September 2026
**Status:** advised
**Run-ID:** run_727b1126-6e4d-44a9-9996-16f72909a29a
**Document-ID:** doc_a89555cd-e8c1-4fb1-981b-68a17d0bfe8d
**Reviewed:** https://github.com/Bitspark/bittype/issues/1#issuecomment-5861714323
**Owner:** the coding agent working bittype #1 for the project lead
**Issue:** https://github.com/Bitspark/bittype/issues/1 (coordinated under https://github.com/Bitspark/bittype/issues/5)

## Question

We are designing a new **contract language**, from its **value-type core** outward:

- A contract language declares what data looks like and which operations exist.
- The value-type core is the part that declares data: records, choices, enumerations, aliases, optional and nullable fields, constraints, generic types, and imports.

The language will be named **bittype**. Two requirements shape it:

- **Wire independence.** A declaration has to mean the same thing whatever network protocol, storage format or code generator later consumes it.
- **Canonical identity.** Every declaration needs a digest computed from a precisely specified encoding of its content. Separate libraries in several programming languages must reproduce that digest exactly.

It replaces an older language, the one inside **nightseam**, which five production projects depend on. We treat the old language as evidence to learn from, not a specification to copy.

We would like your advice on five things, which map one to one onto the questions at the end:

1. **The kinds and primitives** the core should admit, derive or exclude.
2. **Presence, nullness and constraints**: which facts belong to a field and which to a type.
3. **Generic parameters and their application**, so that generating a specialized type ahead of time and applying a generic at run time reach one identity.
4. **Canonical declaration identity**: what it contains, what imports pin, and how the encoding is specified, named and versioned.
5. **Anything we are not seeing**, and in particular how to shape the core now so that **abstract operations**, typed errors, one-way operations and reverse calls can extend it later without naming any protocol. Operations are not designed in this round; the core just must not rule them out.

**What the answer is for.** The design issue, bittype #1, closes on three deliverables, reviewed against a written contract theory:
- a written specification of the value-type core;
- a named, versioned canonical encoding;
- shared test vectors, written independently of any implementation.

The **project lead decides**; this consultation informs that decision. The rest of the chain depends on the result:
- Go and TypeScript type generation and identity libraries (bittype #3);
- a validator whose data descriptions must equal bittype's output byte for byte (bitschema #2);
- a generator that maps operations onto a protocol (bitlink #1);
- a first real consumer (bittree #53).

**The form of answer that helps us most.** For each question, please give:
- a recommendation;
- the principle that decides it;
- the laws or test vectors it implies;
- what it rules out;
- where comparable systems made the same choice and how it held up. Examples include Protocol Buffers, Cap'n Proto, Smithy, TypeSpec, Avro, CDDL, JSON Schema, Dhall, Unison's content-addressed definitions, and algebraic data types in ML or Haskell.

**Why it matters.**
- **Revisions are permanent.** A published language revision is immutable, so a wrong choice in the core costs a new revision and a migration for every consumer.
- **Identity must be exact both ways.** A wrong identity scheme makes equal declarations look different, or different ones look equal. Identity is how a generated client and server establish that they speak the same contract.
- **Generics must not split identity.** A wrong generics design would give specialized and run-time-applied types different identities for the same type.

## Context

### System Overview

**The project.** A group of small repositories, each owning one decision, replaces a monolithic predecessor. In this document "the project" means that group. "Family" is reserved for a unit of declarations, defined below.

| Repository | What it is | State |
| --- | --- | --- |
| **bittype** | The contract language: syntax and semantics of value types, abstract operations, callable signatures, generics and imports. It also owns canonical declaration identity, native type rendering per programming language, and small runtime libraries for presence and identity. | Chartered, documents only. **This consultation designs its core.** |
| **bitschema** | **Descriptions**, the runtime form of a type, and one validator per programming language that checks a JSON value against a description. | Design document only |
| **bitlink** | Maps abstract operations onto a concrete wire protocol, and generates client and server code for it. | Documents only |
| **bitwire** | The project's wire contract: interfaces for sending messages, and the JSON-over-WebSocket protocol `bitwire/1`. | Released |
| **bitruntime** | The Go and TypeScript implementation of bitwire. | Released |
| **bittheory** | A written contract theory with laws the other repositories are checked against. | Prose, no stable law identifiers yet |
| **nightseam** | The predecessor, holding a declaration language, a code generator and a runtime together. | Frozen at v0.6.0 |

**The production consumers of nightseam's language:**

| Project | What it is |
| --- | --- |
| **bittree** | "One tree for a whole repository": files and their syntax trees behind one interface. It is the first generated consumer of the new language. |
| **repo-tool** | A persistent repository service with generated API, command-line and tool clients. |
| **nightforge** | Turns closed coding-agent command-line tools into instrumentable runtimes. It contributes two declaration sets, *concur* and *ulab-runner*; the latter runs agent sessions in containers. |
| **nighthall** | Hosts coding-agent sessions on machines or in the cloud, steered from a browser. |
| **bitsystem** | The original prototype of the space model, not to be confused with its successor. |

**Conventions.**
- "nightseam#343" means issue 343 in the nightseam repository. Such issues are summarized in the text, so no link needs to be followed.
- Quotations keep their original capitalization. "Bitlink", "Bitwire" and "Nightseam" are the same projects as bitlink, bitwire and nightseam.
- "Decision 0010" is an architecture decision record in bitwire that set the division of work between all these repositories. "The operator" in quotations is the project lead.

**Implementation languages.** Go and TypeScript come first. Every rule must be implementable with identical results in both, and later in others.

**Glossary.** Terms used across this document:

| Term | Meaning |
| --- | --- |
| **Declaration** | Text a person writes (today JSON files) that declares named types and, later, operations. |
| **Family** | nightseam's unit of declarations: one directory of files, imported and versioned as a whole. Whether the new language keeps this unit is open. |
| **Declaration kind** | A form of declaration: record, union, enum, alias, entity, and so on. The word "kind" in this document means this, *not* the type-theoretic "type of types"; where that meaning is needed, the text says **type-theoretic kind**. |
| **Record** | A named set of fields, each with a name and a type. |
| **Union** | Exactly one of several named variants, each possibly carrying a payload. A *tagged* union names the chosen variant in the data. |
| **Enum** | A closed set of names without payloads. |
| **Alias** | A second name for a type expression. It is *transparent* if interchangeable with its target, and *nominal* if distinct from it. |
| **Entity** | A record with a declared key field that identifies application data. The key is a durable ID. |
| **Presence** | Whether a field appears in an encoded object at all. An *optional* field may be absent. |
| **Nullness** | Whether a field that is present may hold the explicit value `null`. |
| **Constraint** | A restriction on values of a field: `min`, `max`, `length`, `pattern` (a regular expression), `unique`. |
| **Generic parameter / application** | A placeholder in a declaration (`Page<T>`), and filling it (`Page<Part>`). |
| **Specialization vs run-time application** | Generating a concrete type for `Page<Part>` ahead of time, versus keeping `Page<T>` generic in generated code and filling `T` in the running program. |
| **Canonical encoding** | A specified, deterministic byte serialization of a declaration's meaning: equal meanings give equal bytes. |
| **Declaration identity, digest** | A hash, for example SHA-256, of the canonical encoding. It names the declaration's content. |
| **Vectors** | Test cases shared by every implementation: declarations with their exact canonical bytes and digests, plus declarations that must be **refused**. A refusal has a stable code and exact message, e.g. `[cyclic_type]`. |
| **Native rendering** | The Go or TypeScript types generated from a declaration. |
| **Generation manifest** | A file written beside generated code recording which language, encoding and generator revisions produced it. |
| **Description** | bitschema's runtime JSON form of a type, which a validator checks values against. It is always *closed*: it has no unfilled generic parameters. |
| **Wire, profile, carrier, adapter** | The network side, all outside the language:<br>• **wire**: messages between processes;<br>• **profile**: the rules for how operations become messages;<br>• **carrier**: how a value is represented inside a message, e.g. how a union's tag is spelled;<br>• **adapter**: generated code that sends and receives calls. |
| **Abstract operation** | An operation's name, argument, result and errors, with no mapping onto a protocol. |
| **Callable** | An operation signature treated as a type, so values can carry references to callable things. |
| **Live value, capability** | A handle to a running remote object, as opposed to plain data. A value containing one is *live*. |
| **Reverse call** | A call from server to client, the reverse of the usual direction. |
| **Revision** | A published, immutable version of the language. |
| **Slice** | The subset of language features the first consumer needs. |

**How a declaration flows through the project:**

```text
author writes declarations (bittype's surface syntax)
  → bittype parses, resolves imports and checks them     (refusals with exact codes and words, shared vectors)
  → canonical encoding → declaration digest               (one library per programming language, shared vectors)
  → native types (Go, TypeScript) + generation manifest   (bittype)
  → descriptions, closed, with parameters filled          (must equal bitschema's fixtures byte for byte)
  → operations mapped onto the wire, adapters generated   (bitlink, outside the language)
  → at run time, identities compared between two sides    (outcomes: compatible, incompatible, unverified)
```

A generic type can reach a running program either specialized, or generic and applied at run time. Both routes must reach one identity.

**Is JSON the value model?** No decision fixes it:
- In practice, every consumer's values are JSON.
- bitschema validates JSON values, and its descriptions are JSON.
- The current wire, `bitwire/1`, carries JSON.

The language must nonetheless be wire-independent, and is meant to serve storage, UI and command-line generators too. So it is **open** whether the language defines its values abstractly, with JSON as one representation, or defines them as JSON values. This bears on the union tag question in contradiction 3 below.

### Architectural Context

**The rules the language must satisfy.** These are quoted from bittype's charter, its instructions to contributors, decision 0010, and bittype #1:

- **Ownership.** *"Value types, abstract operations and callable signatures belong to the language. Paths, the profile and message representation belong to Bitlink's binding."* Here a path is the address of an operation on the wire, and the binding is bitlink's mapping of operations onto it.
- **Dependencies.** bittype and bitschema depend on nothing else in the project, and in particular not on each other. The generator depends on both.
- **Wire independence.** *"Test each declaration with: 'would it make sense, unchanged, with a different adapter and no Bitwire?'"*
- **Four things kept distinct:** *"declaration identity; profile revision; generated-code compatibility; behavior the declaration does not represent. A digest names declared content, not behavior."*
- **One identity implementation per language.** *"Hold every identity library to the shared canonicalization vectors. Two implementations of one encoding in one language are a defect."* Identity must support *"generic application at run time without the parser, the transports or the peer"*. That means a small runtime library, not the compiler.
- **Revisions.** *"A published language revision is immutable. A change to what is accepted, rejected or meant is a new revision. The canonical encoding is versioned by name. Any change to identity content is an explicit transition, never a side effect of packaging."* There is no compatibility promise before 1.0.
- **Independent evidence.**
  - *"Compile-time specialization and run-time generic application must reach the same identity. Tests of the generic-application laws at the semantic level, not only on generated text."*
  - The construction square: *"instantiating a generic and then adapting it agrees with adapting it and then instantiating."* Vectors must be written *"rather than using one implementation as its own oracle."*
- **The boundary with bitschema.** *"the descriptions bittype emits must equal the ones bitschema validates against, with neither repository importing the other."* The boundary is kept as byte-exact fixtures shared by both sides.
- **Design, not migration.** *"This is a new language, designed rather than migrated. Nightseam's language and the family design are evidence to examine, not specifications to copy."*
- **The first consumer's slice** must cover *"absence and null, errors, reverse calls and generic identity"*. These are *"requirements to consider, not permission to smuggle paths, carriers or profiles into this language."* If the core does not settle how callables extend it, a named design blocker is recorded instead.

**The contract theory (bittheory), which the design is reviewed against.** Notation:
- `S` is a *shape*, the structure of a contract; `B` is a behavioral predicate.
- `O` stands for opaque values.
- `×` is a product (pair) and `+` a sum (either-or).
- `Option(X)` is an optional `X`.
- `FiniteMap(Key, X)` is a finite map from child names to sub-shapes.
- `μX.` is the least fixed point, that is, an inductive definition.

The theory states:

- **Contracts.** *"A contract is `C = (S, B)`"*. *"A shape specifies available sorts, fields, constructors, operations, arguments, results, and roles."* Here "sorts" means the categories of things a shape provides. A native declaration *"alone need not express `B`"*.
- **Shapes.** `Shapes(O) = μX. Option(O) × FiniteMap(Key, X)`: a node has an optional value and a finite map of named children. *"Child order has no shape meaning. A node with neither a value nor children is a valid, present shape; it is different from a missing child."*
- **Generics are holes.** `Shapes(O, G) = μX. Generic(G) + (Option(O) × FiniteMap(Key, X))`, where `Generic(g)` is a hole named `g`. *"The current language has no binders; introducing binders later would require additional scoping and capture-avoidance rules."*
  - **This conflicts with declared type parameters.** `Page<T>` binds `T`, so a design with parameters goes beyond the theory as written (contradiction 7).
- **Substitution laws.** `S[σ]` applies a substitution `σ` that fills holes, and `σ ⋆ τ` composes two substitutions.
  - S1–S3 say that filling a hole yields its replacement, that substitution distributes over children, and that the identity substitution changes nothing.
  - S4: `(S[σ])[τ] = S[σ ⋆ τ]`.
  - Substitution is simultaneous: *"substituting `g0 := Box[g0]` once into `g0` produces `Box[g0]`, rather than expanding forever."*
  - S5–S7 concern navigation and are not relevant here.
- **The generation square (S8):** *"fill, then transform; or transform the template and every argument, then fill."* "Transform" means code generation, and the two routes must agree up to an equivalence between generated representations. *"Exact generated-source equality is a stronger requirement. A specialized generator and a generic template may produce different but equivalent code."*
- **Digests:** *"a declaration digest names content, not behavior"*, and *"equal digests alone do not prove artifact faithfulness, behavioral equivalence, or interoperability."*
- **Nearby statements, which are analogies and not rulings on fields:**
  - About coordinate axes: *"an omitted key differs from a key present with `null`"*.
  - About tree node values: *"Optionality is an explicit instantiation `T = Option(O)`"*.
  - About a tree's key encoding: it must *"refuse keys outside that image rather than normalize them"*. The image is the set of byte strings the encoding can produce.
- **What the theory leaves out.** It has no vocabulary for value types: no records, unions, enums, aliases, constraints, imports or subtyping, and parameters have no sorts or bounds. The theory assigns *"Canonical declaration and generic-application laws"* to bittype to state and test.

### The predecessor: nightseam's declaration language (v0.6.0), as evidence

**Tiers.** A nightseam family has three declaration files:
- `model.json` holds the types;
- `protocol.json` holds the operations, errors and parameters;
- `live.json` holds callable types.

*"A declaration refers to its own tier or a lower one, never a higher one."* The key `"nightseam": 2` in `model.json` is the version of the language. `protocol.json` must declare `"profile": "nightseam.duplex/1"`, **naming the wire protocol inside the language**. The new language must not do that.

**A real declaration** (a test family; `TextPart` is another record in the same family, and descriptions are omitted):

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

**Declaration kinds.**
- `entity`.
- `record`: may `extends` other records, whose fields come first, and may be `open` to unknown fields.
- `enum`: string values only.
- `alias`.
- `union`: a `tag` and `variants`. Each variant is a type or an inline record, and `{"empty": true}` is a variant with no payload. A union may also `extends` another union.
- `callable`: in `live.json` only.

**Type expressions.**
- **Primitives:** `string`, `boolean`, `integer`, `number`, `timestamp`, `json`.
- **Collections:** `{"array":T}`, and `{"map":T}` with string keys.
- **Other forms:** `{"nullable":T}` (`{"nullable":{"nullable":T}}` is refused), `{"literal":"text"}`, `{"ref":"User"}` (an entity's key), `{"apply":X,"with":{…}}`, and inline records.

There is **no bytes type, no numeric width and no decimal type**. The meaning of each primitive lives in the validators, not in the language:
- `integer` is limited to JavaScript's safe range, at most 2^53 − 1, rendered as Go `int64` and TypeScript `number`.
- `number` is a finite 64-bit float.
- `timestamp` is an RFC 3339 string.
- `json` accepts any JSON.

**Presence and null.** nightseam's own decision record:

> "`{"nullable": T}` is a type expression, and a field's `nullable: true` is its sugar … `required` stays where it is, on the member, because presence is a fact of a member and not of a type: a member may be absent from an object, and a value cannot be absent from itself."

Generated code wraps each fact in a small struct, a *presence wrapper*:

```go
// Go (generated): a field not required is Optional[T]; a nullable one is Nullable[T]; both nest.
type Page[T any] struct {
	Items []T                                        `json:"items"`
	Next  runtime.Optional[runtime.Nullable[string]] `json:"next,omitzero"` // omitzero: omit when absent
}
// runtime: type Optional[T any] struct { Value T; Present bool }
//          type Nullable[T any] struct { Value T; Null bool }   // an unset Nullable is NOT null
```

```ts
// TypeScript (generated). PartImage is the name derived from the inline record's position.
export interface Page<T = unknown> { "items": Array<T>; "next"?: string | null; }
export type Part = { "type": "count"; "value": number } | { "type": "image"; "value": PartImage }
                 | { "type": "text"; "value": TextPart };
```

These wrappers are held inline, not through a pointer. A type that contains itself through them would therefore have infinite size in Go.

**Unions are adjacently tagged**: the tag and the payload are two members, `{"type": "image", "value": {...}}`.
- This replaced a *flat* form that merged the tag into the payload, which lost information (nightseam#146): a JSON payload `3` and a record `{"value": 3}` became indistinguishable.
- In Go a union is a struct with one pointer per variant: *"construction can select zero or several alternatives. The generated codec enforces exactly one; the Go type alone cannot."*
- **An inline record is named after its position**, so moving an unchanged inline record renamed its generated type in both languages.

**Generics** have two sorts of parameter:

> "A **family parameter** — it has `of`. It is filled by a **family** … A type is drawn *through* it, `S.Envelope` … A **type parameter** — it has no `of`. … A family parameter is **not** a type parameter with a bound … a family is not a type — nothing is an instance of it."

- **Family parameters.** `S` is filled by a whole family, and `S.Envelope` names that family's member type `Envelope`.
- **Type parameters** have no bounds, no higher-kinded arguments (parameters that are themselves generic) and no inference.
- **Application** is always explicit. Unused parameters are refused.
- In generated Go, a run-time application is bound explicitly, e.g. `schema.Bind(map[string]any{"T": runtime.TypeArgument[T]()})`.
- nightseam held one law, in Go, TypeScript and on the wire, for family parameters: `gen(bind(C, F)) ≅ gen(C)[F]`. Here:
  - `C` is a family with a family parameter;
  - `F` is the family that fills it;
  - `bind` fills the parameter;
  - `gen` is code generation;
  - `≅` means equivalent generated code.

  In words: generating the filled family equals filling the generated one. A later study observed that this was checkable only *"because the generic declaration survives into the package rather than being compiled away."*

**Imports** name another family directory, with no version or digest.
- Import cycles are refused.
- A lower-case qualifier names a family (`identity.User`), and an upper-camel-case one names a parameter (`S.Payload`).

**Constraints.**
- `min`/`max` apply to integer, number and timestamp fields.
- `length` applies to strings (counted in code points) and arrays.
- `pattern` applies to strings, in a restricted ECMAScript regular-expression dialect with no lookaround and no backreferences.
- `unique` applies to entity fields only: *"the holder's to enforce and not the wire's — a validator sees one value."*

Constraints attach only to record fields, never to aliases, array elements or map values. So an alias `Email = string` cannot carry a `pattern`.

**Operations**, for completeness. A family has a server side and a client side. The client side lists reverse calls.
- **Methods** have an optional request object and a required `result`. **Events** carry a single payload and expect no reply.
- **Errors** are codes with a description and no typed payload.
- A **callable** is a type with an optional request, result and error codes. It is never inline, it is identified by its declared path, and liveness is derived: *"contains a callable"*.
- The built-in `duplex.Envelope` and `Handle` types put wire vocabulary into the type table.

**Canonical declaration identity.** The digest is the lowercase hexadecimal SHA-256 of a canonical JSON graph.

*Structure.*
- The graph is `{"definitions": {...}, "root": ..., "version": 1}`.
- Definitions are keyed by a family name (`same`), a type path (`same/Payload`), or a side path (`same/$server`, `same/$client`). A side is a set of operations: methods and events.
- A digest can be taken of a single type, covering only the declarations it reaches.
- Non-generic callables keep their whole declaring family's digest, so any unrelated edit to that family changes every callable's identity.

*Content.*
- *"Every field has `name`, `type`, `required`, `nullable`, and `unique`, including the boolean defaults."* Explicit and omitted defaults give the same digest; this is a shared vector.
- *"Error descriptions are documentation and are excluded."*
- *"Pure aliases disappear into their normalized target expression."*
- References are by qualified name, with the referenced content included, so identity is both nominal and structural. Callables stay nominal even when signatures match.
- A union's identity includes its `tag` and `value` member names, which are **carrier facts** about the wire representation.

*Ordering.*
- *"All object keys are sorted by their UTF-8 byte sequence … error codes, enum values … are deduplicated and sorted … Field, parameter, inheritance, and application-argument order is meaningful."* Note that the theory says child order has no shape meaning.
- An application's identity lists captured family parameters first, then the constructor's own arguments in declared order. It has a printable form, e.g. `worker/Function<integer,string>`.

*Encoding.*
- Numeric constraints become exact decimal strings: `1.00` becomes `"1"`, and `100` becomes `"1e2"`.
- The bytes are Go's `encoding/json` output, including HTML escaping of `<`, `>` and `&` as `\u003c`, `\u003e` and `\u0026`, and escaping of U+2028 and U+2029 as `\u2028` and `\u2029`. This is **not** RFC 8785, the JSON Canonicalization Scheme, which also sorts keys by UTF-16 code units, not UTF-8 bytes.

A real canonical encoding of a one-field record. The `…` omits two empty sides.

```json
{"definitions":{"same":{"client":{"ref":"same/$client"},"errors":[],"kind":"family","parameters":[],
 "server":{"ref":"same/$server"},"types":{"Payload":{"ref":"same/Payload"}}}, …,
 "same/Payload":{"fields":[{"name":"value","nullable":false,"required":true,
 "type":{"primitive":"string"},"unique":false}],"kind":"record","open":false,"parameters":[]}},
 "root":{"ref":"same"},"version":1}
```

nightseam implemented this canonicalizer three times, in the generator, the Go runtime and the TypeScript runtime, all held to 76 shared cases.

**Recorded defects and regrets:**

- **The "sugar" is not equivalent.** A field `{"type":"string","nullable":true}` and a field `{"type":{"nullable":"string"}}` generate identical Go, yet their canonical fragments differ, and so do their digests:
  - `"nullable":true,…"type":{"primitive":"string"}`
  - `"nullable":false,…"type":{"nullable":{"primitive":"string"}}`

  The checker also treats them differently: the first accepts a `pattern`, while the second is refused with `pattern holds a string, not {"nullable":"string"}. [invalid_constraint]`. The loader never desugars. This was reproduced with the released tool.
- **Recursion is refused.** A cycle through an array, map or nullable is rejected with `"Value schema cycle reaches Node. [cyclic_type]"`; inline records are exempt. Consumers declare recursive data as `json` instead.
- **The identity was once too narrow** (nightseam#343). The digest hashed only types, so changing a method's result from `string` to `integer` left it unchanged. The operator ruled: *"hash the declaration, canonically rendered … never a backend-specific wire rendering … operations and imported declarations included by content (not by name)."*
- **A built-in one-way operation was declared as a method** (nightseam#724). Methods require a `result`, so this declaration contradicted the runtime, which treats it as an event.
- **The only compatibility test is equality.** *"Compatibility between distinct revisions remains consumer policy."* Adding an optional field changes the digest. Because records are closed and validation is strict, that change also breaks readers (nightseam#670).
- **Values have no canonical encoding** (nightseam#669). One service hashed the JSON of a Go struct, and adding a field changed every digest.
- **The identity covers forms the checker refuses.** Its vectors include mutually importing families, which the language rejects.
- **Checking and identity disagree on constraint meaning.**
  - `min`/`max` are compared as 64-bit floats, while identity keeps exact decimals.
  - A timestamp `min` is a JSON number, compared *as text* against the RFC 3339 string.

### What the consumers actually declare

Five projects use nightseam's language, and nightforge contributes two declaration sets. These counts come from a script over their declaration files; the bracketed question numbers show where each finding bears.

| | bittree | repo-tool | concur | ulab-runner | nighthall | bitsystem |
| --- | --- | --- | --- | --- | --- | --- |
| Records / enums / unions / aliases | 44 / 6 / 4 / 2 | 23 / 2 / 0 / 0 | 39 / 11 / 1 / 0 | 22 / 6 / 0 / 0 | 114 / 14 / 0 / 2 | 48 / 10 / 1 / 2 |
| Fields / optional / nullable | 151 / 57 / **0** | 87 / 0 / 0 | 203 / 73 / 0 | 64 / 11 / **3** | 447 / 114 / 0 | 270 / 135 / 0 |
| Fields with constraints | 12 (each `min` and `max`) | 0 | 38 (`min`) | 17 | 14 (`min`) | 0 |
| Generics, entities, `extends`, imports | none | none | none | none | 4 imports | none |
| Reverse calls (client methods) | 0 | 0 | 0 | 0, but 11 callables passed as values | 8 | 3 |

- **[Q2] Null is rare, and meaningful when used.** ulab-runner's `holder` field is always present, and `null` means *"nobody holds control"*.
- **[Q2] Optional fields are common, and their meanings vary by field.**
  - For one bittree payload, *"omitted and empty mean the same bytes"*.
  - For edits, *"empty `text` means absence … explicitly empty `bytes` remains present"*, and both implementations must round-trip `{"op":"set","bytes":""}` exactly.
- **[Q1] The most common workaround is a union written as flat optional fields.** Presence depends on another field's value, stated only in prose:
  - bittree's `Status`: *"`result` is present exactly when finished"*.
  - bittree's `Segment`: *"Exactly one member"* of four optional fields.
  - nighthall's `Frame`: one field is *"present only for"* one value of another. It is checked by hand: *"ValidateFrame adds cross-field invariants that the declaration language does not express."*
- **[Q1, Q2] `json` is the escape hatch** for recursive data, such as bittree's syntax-tree `Node`, and for untagged unions, such as bittree's `Edit.to`.
- **[Q3] Generics appear only in an experimental probe:** a record with four type parameters, and a generic callable applied inside a record. The probe measured the cost as *"generic parameters flowing through generated signatures plus one codec per native representation."*
- **[Q5] Errors, callables and reverse calls.**
  - Errors are code lists with no typed payload, and bittree carries its error details as untyped bytes.
  - ulab-runner passes a record of callables to the runner, which calls them back to deliver output: a reverse call through values. Its callables declare their own error codes.
  - nighthall and bitsystem declare client-side methods.
- **bittree's measured gaps:**
  - recursive declarations are refused;
  - there is no `bytes` primitive;
  - a declared type cannot be mapped onto an existing native type;
  - the Go field representations (`int64`, presence wrappers) don't convert to its native model;
  - the generated TypeScript fails under the strict compiler option `exactOptionalPropertyTypes`, which makes `x?: T` refuse an explicit `undefined`.

### The runtime form: bitschema's description basis (design only)

bitschema owns "the runtime form of a type". Its admission rule: *"A constructor is admitted only with a scenario that no smaller basis passes, and that scenario is a row in `tables/validator.json`."* A *constructor* is a building block of descriptions, and the *basis* is the admitted set. Anything else is *derived*, meaning expressible in what is admitted, or excluded.

| Constructor | Disposition |
| --- | --- |
| Primitives `string`, `boolean`, `integer`, `number`, `timestamp`, `json` | Admitted, a closed set. `bytes`, `date` and `decimal` are open. Numeric width and range are not stated. |
| Record | Admitted. Each member carries two flags, `required` (default true) and `nullable` (default false). *"JSON has two distinct absences … Collapsing them into one union member makes a validator unable to tell a consumer which of the two it received."* |
| Discriminated union | Admitted: *"a discriminator name and a set of variants."* |
| Enum | Derived, as a union whose variants carry nothing. Whether it keeps a spelling of its own is open. |
| Array, map | Admitted |
| Named references | Admitted. Recursion is derived through names: *"a cycle through non-optional, non-collection members is refused at load"*, so an optional member may break a cycle. |
| Foreign | Admitted: a type from another separately versioned set of declarations. |
| Opaque | Admitted: *"JSON carried through unexamined — … what fills a generic slot whose family is chosen elsewhere."* |
| Reference | Admitted: a token claiming a declared interface, checked for form only. *"an arrow is not a value here"*, meaning a function type is not a value type. |
| Interface | Admitted: named members, each with a request, a result and errors. |
| Constraints | Admitted per kind, *"only for the kinds that give it meaning"*. |
| Entity and key reference | Derived. An entity is metadata; a key reference validates as the key's primitive. |
| `extends` | Derived: *"Flattened by bittype before a description exists."* |
| Intersection, untagged union, parameters whose arguments are themselves generic | **Excluded** |

Two further rules:
- **Descriptions have no parameters.** *"a description is the runtime form and is always closed — a parameter is filled where a rendering is instantiated, and what may fill it is bittype's check."*
- **Identity is split.** *"the structural hash is bitschema's; anything beyond structure — identity, scope, provenance — belongs to bittype's own hash … Two descriptions of the same shape share a hash even where the types mean different things, and that is intended."*

### Earlier design work in the project, and its study of nightseam

An earlier design document proposed a first tier of entity, record, enum and alias. It kept presence and nullness as two field facts: *"A field that is not required is absent, not null."* Its generics were family-level. It left open where `unique` belongs, composite keys, the callable-reference syntax, the exact set of primitives, and `date`, `decimal` and `bytes`.

A later study of nightseam recommended the following. These are advice, not rules.
- **Recursion through collections only:** *"no value contains itself except through a collection, which the empty collection always satisfies."*
- **Avoiding the Go presence-wrapper trap.** A type that reaches itself through `{"nullable":T}` renders as `Nullable[T]` held inline, *"an infinite-size Go type that does not compile."* The same objection applies to bitschema's rule that an optional member breaks a cycle.
- **Continuing to leave out:**
  - a dedicated recursion constructor;
  - bounds on type parameters (*"so that a declaration cannot be live and say otherwise"*: liveness is derived, never declared);
  - nominal aliases;
  - open unions;
  - structural subtyping;
  - read-only and write-only markers;
  - default values;
  - an open-ended set of kinds.
- **Keeping nightseam's strength:** *"One shipped declaration, interpreted by one validator per language … held to a single table case by case — including the exact words of every refusal."*

It also recorded that nightseam's aliases are transparent Go aliases (`type UserId = string`), so `UserId` and `OrderId` are interchangeable.

### Contradictions in our own documents

These are unresolved, and your view on them is welcome:

1. **Nullness: flag or type?** bitschema's design keeps two member flags and calls them *"Nightseam's choice, kept"*. That describes nightseam before its version 0.4. Since then nightseam has made nullness a type constructor `{"nullable":T}`, with the flag as sugar, and the two now differ in both identity and checking.
2. **Recursion: which edges may break a cycle?** bitschema lets an optional member break one. The study allows only collections, because of the Go wrapper problem.
3. **The union tag.** bitschema's union is *"a discriminator name and a set of variants"*, and bittype's descriptions must equal bitschema's byte for byte. Yet the discriminator's representation is a carrier fact that belongs to bitlink, and nightseam's identity includes the tag names. Where the tag lives depends partly on whether the value model is JSON.
4. **Parameters in descriptions.** bitschema says descriptions are always closed. Its own work item asks to resolve *"generic substitution"*, and run-time generic application must still reach one identity.
5. **Who owns what goes beyond value types.** The earlier design put interfaces, families and imports in a separate system-declaration layer, which decision 0010 did not create. It gives bittype *"abstract operations, callable signatures, generics and imports"*. Whether "family" survives as a unit of versioning and import is open.
6. **`json` versus opaque** in bitschema, and **`unique`**. Is `unique` a property of one value, checkable by a validator, or of a collection of entities, the holder's to enforce?
7. **Binders.** The theory has no binders and warns that adding them needs scoping rules. A design with declared type parameters adds them.
8. **Still undefined:**
   - numeric domains;
   - the exact form of `timestamp`;
   - `bytes`, `date` and `decimal`;
   - composite keys;
   - whether `enum` keeps its own spelling;
   - nightseam's `literal` and inline-record forms, and the position-derived names of inline records.

### Constraints

**Fixed.** These are all quoted under Architectural Context:
- wire independence;
- no dependency on the rest of the project;
- the byte-exact bitschema boundary;
- one identity library per language, with independent vectors;
- identity names content, not behavior;
- immutable revisions and an encoding versioned by name;
- specialization and run-time application reach one identity;
- Go and TypeScript first;
- the language is designed, not migrated;
- unresolved semantics become named blockers, never invented rules.

**Earlier rulings that bind:**
- nightseam#343: operations and imported declarations are included in identity by content, not only by name.
- Identity checks distinguish compatible, incompatible and unverified outcomes, and never count "not checked" as a match.

**Advisory only:** the study's "continue to leave out" list, bitschema's excluded constructors, and everything nightseam did.

**Flexible:**
- the surface syntax, which need not be JSON;
- whether the value model is JSON;
- the set of kinds and primitives;
- how presence, null and constraints are expressed;
- the generic system, including whether family parameters exist;
- the canonical encoding and its content;
- nominal or structural identity;
- what imports pin;
- how much of the consumers' needs the first revision admits.

**Scale.** Families range from tens to a few hundred types; nighthall has 114 records with 447 fields. Correctness, portability and stable identity matter far more than performance.

### What We've Considered

These are directions, not decisions. The numbers in brackets are the contradictions each one touches.

- **Kinds: a sum-of-products core.**
  - Records, plus tagged unions whose variants may be records. That would replace the flat-optional-fields workaround.
  - An enum as a union without payloads.
  - Transparent aliases, invisible to identity [8].
  - An entity as a record plus key metadata.
  - `extends` flattened before identity is computed.

  We have not resolved how a tagged union's discriminator stays out of the language [3].
- **Presence and null as two separate member facts**, with no nullable type constructor. That would take bitschema's side of [1] and avoid nightseam's two spellings with two identities. But a nullable collection element (an array of nullable strings) may still need a type-level form.
- **Recursion only through collections**, or through an indirection that the rendering represents as a pointer [2]. We have not decided which.
- **Constraints** attached to types (aliases, elements) rather than only to fields, with one exact numeric meaning shared by checking and identity [6].
- **Generics** [4, 7]:
  - type parameters only in the first revision, with family parameters deferred;
  - explicit application;
  - the identity of an application is the constructor's identity plus the identities of its ordered arguments;
  - generic declarations survive into generated packages, so run-time application is possible.
- **Identity.** A canonical encoding of the resolved declaration graph:
  - aliases resolved;
  - default flag values written out;
  - documentation excluded;
  - imports included by content;
  - an established canonicalization (RFC 8785 JSON, or deterministic CBOR) rather than one library's output;
  - the encoding named and versioned;
  - one digest per named declaration, so an unrelated edit changes nothing else.

  We have not resolved nominal versus structural identity, or whether field order is meaningful [5].
- **Operations, later.**
  - Callable signatures with typed errors.
  - One-way operations as a declared form.
  - No paths or profiles in the language.
  - A reference to a callable is distinct from an entity key.

## Questions for the Expert

1. **Kinds and primitives.** Which declaration kinds and primitives would you admit, derive or exclude in the first revision of a wire-independent value language that renders natively into Go and TypeScript, and by what principle?
   - How would a tagged union keep its discriminator out of the language, when the discriminator's representation belongs to the wire binding?
   - Should the language's values be defined abstractly, or as JSON?
   - Where comparable languages made these choices, which patterns held up?
2. **Presence, nullness and constraints.** Which facts belong to a field and which to a type? How would you place presence, nullness and constraints so that each:
   - has exactly one canonical spelling;
   - composes with unions, collections, aliases, generics and recursion;
   - renders idiomatically in Go and TypeScript;
   - is checked with the same exact meaning its identity records?

   How should field-specific meanings such as "absent equals empty" be expressed, if at all?
3. **Generics and generic identity.** How would you design type parameters and their application so that generating `Page<Part>` ahead of time and applying `Page<T>` at run time reach one identity?
   - Consider parameter sorts (nightseam's type and family parameters), bounds, and what family-level parameters would buy and cost.
   - How should application relate to descriptions, which must be closed?
   - Which laws would you state and test at the semantic level, given that the theory's generics have no binders?
4. **Canonical identity.** What should a declaration's canonical identity contain? Consider:
   - whether aliases are visible to it;
   - nominal names versus structure, including inline records;
   - field order, given that the theory's children are unordered;
   - what an import pins (a name, a version or content);
   - whether "family" remains the unit of versioning.

   How would you specify, name and version the encoding, and write acceptance and refusal vectors, so that independent implementations in several languages provably agree? Which patterns from content-addressed systems apply?
5. **Open.** How would you shape this core now so that abstract operations, typed errors, one-way operations and reverse calls can extend it later without naming a protocol? And beyond that, looking at the predecessor's defects, the consumers' workarounds, the theory and the split between language, descriptions and wire binding: what would you do differently, and what are we not seeing?
