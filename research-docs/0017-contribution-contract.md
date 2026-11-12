# Research: The contribution contract of a type generator that must also carry conversion decisions: demands between contributions, open facets and the identities of the plan

**ID:** 0017
**Date:** 12 November 2026
**Status:** draft
**Owner:** claude-contract-chain-issues
**Issue:** https://github.com/Bitspark/bittype/issues/196 (item 6), with #195, #197 and #207
**Code baseline:** bittype 2daea52e18e20369401500ddeeb19efae581b306

## Question

**horos** is a declarative type language. **bittype** is a Go library that reads
horos declarations and generates native type declarations for four languages:
TypeScript, Go, Rust and C++. It is extensible: every rule that turns one kind
of horos type into a native type is supplied from outside the shared mechanism,
by a package that a consumer registers. We call the shared mechanism the
**kernel** and one registered unit a **contribution**.

Today a contribution describes native types and nothing else. Two accepted
designs of 9 November 2026 give it a second duty. Services declared in horos
need conversion code between native values and message values. That code is
emitted by a **binding backend**, one per language, which is owned outside
bittype and, for TypeScript, runs in another process. bittype stays the only
authority that selects native representations. It has to publish every
selection, with every decision the selection contains, as a typed **plan** that
a backend can implement "without guessing a rule".

The accepted design says the existing contract between kernel and contribution
"changes once", for every existing rule of the four languages, before any new
kind of type is added. After that change, a new kind of type must be a new
contribution that edits nothing shared. The owner's direction of 9 November
2026 is:

> I want this properly done, so that we don't accumulate technical debt and have to refactor later.

Nothing is released. No consumer outside our own repositories depends on any
interface shown here, so everything implemented can change now. What we must
avoid is a second change of this contract after contributions and a backend in
another repository are built on it.

**What we treat as settled, and ask you to attack if you see a fault.** One
contribution is everything one owner publishes for one horos declaration or
builtin form, one language and one profile (a named promise about the native
types). The contribution states its profile itself. For one use of a type, a
rule returns **one record of its decisions** together with the function that
computes the native type from that record, so that each decision is written
once. The kernel issues the result, links it to the results of other uses and
reads no meaning of the record. The property behind this ("one instance, two
projections") is stated by the accepted design. The shape is our derivation
from it (section 6.1).

**What is open.** Three parts of the contract have no shape in any accepted
document, and each is public:

1. **How one contribution asks the kernel for another's answer** when the thing
   asked for is not a type occurrence in the source. The case at hand: a record
   must not decide how an optional field is represented. It has to ask, and the
   kernel must resolve the request to whatever the consumer registered.
2. **Whether a contribution carries an open set of parts ("facets")** whose
   kinds the kernel does not enumerate, and who owns the vocabulary of
   value-level code that one such part, the lowering facet, would produce for
   Go.
3. **Which identities the plan records**: the identity of a contribution that a
   backend in another process recognises, the node of the plan, and the
   identity that a conversion's memo is keyed by at run time. Our documents use
   the word "instance" for more than one of them.

Two further questions ask how to validate the contract before every existing
rule is rewritten onto it, and what such a change typically gets wrong.

**How to read the fixed points.** Every statement this document treats as given
has exactly one of four standings, listed in section 5:

- **Owner rule.** Fixed.
- **Accepted design** (the document and its date are named). Fixed, unless you
  see that it contradicts an owner rule or itself. If so, please say so.
- **Implemented today.** A fact of the present code. It can change now, and
  rule AC-09 (section 2.1) gives it no weight in the design.
- **Our assumption or derivation.** Open.

The code in section 4 is evidence of the present state. It is never a reason
for or against an option.

**How this document is organised.**

| Section | Content |
| --- | --- |
| 1 | The system for an outsider, and the terms used throughout |
| 2 | The architecture rules, the two accepted designs, the work items, and the contracts of the library that consumes the plan |
| 3 | One request traced through today's code and through the accepted design |
| 4 | The contract as it is today, with the code inlined |
| 5 | Fixed points and their standing |
| 6 | What we have considered: the settled part and the three open points, each with options, our leaning and what speaks against it |
| 7 | Questions for the expert |

Quotations give file and line at the code baseline. Quotations from GitHub
issues give the issue number and the time of the comment in UTC. All code is
Go unless marked.

## 1. System overview

### 1.1 The components

- **ontos** is a data model with a canonical binary encoding. A value is an
  atom (a byte string) or a tuple of values. Every message that our services
  exchange is an ontos value. ontos has a Go and a TypeScript implementation.
- **horos** is a type language over ontos values. A horos type denotes a set of
  ontos values. It has builtin forms (UTF-8 text, booleans, lists, records with
  required and optional fields, and others such as integers, maps and unions)
  and named declarations with generic parameters. Declarations live in modules;
  a module is identified by a hash of its content, called its **pin**. For a
  type without free parameters horos produces a **descriptor**, which is itself
  an ontos value, and its operation `diagnose` checks whether a value conforms
  to a descriptor.
- **bittype** reads horos modules, checks them, and generates native type
  declarations. It is written in Go, whatever language it generates.
- **The service library** (repository `bitspark-service`, in our workspace
  `libraries/service`) is the TypeScript run-time library of our services. It
  owns the adapter interface (an **adapter** converts between a native value
  and a message value), the per-message **call context**, the contracts
  `Function`, `Stream` and `Service`, and, by a decision of 9 November 2026,
  the TypeScript binding backend. It will also host the run-time package in
  which the TypeScript conversion functions for bittype's builtin
  representations live.
- **Consumers** are builds of services. The pilot is a service called Content.

An example that the rest of the document reuses:

```horos
module example.pagination;
type Page<T> = { items: list<T>; next?: text; };
type Person = { name: text; };
type People = Page<Person>;
```

`next?` is an optional field: the record value may lack it. bittype generates
today, for TypeScript and for Go (whitespace ours):

```typescript
export interface Page<T> { items: T[]; next?: string; }
```

```go
type Page[T interface {}] struct {
	Items []T
	Next  BittypeOptional[string]
}
type BittypeOptional[T interface {}] struct {
	Present bool
	Value   T
}
```

For Rust the field type is `::std::option::Option<::std::string::String>` and
for C++ `::std::optional<::std::string>`. So the same source decision, "this
field is optional", becomes a property of the field in TypeScript, a generated
helper type in Go and an external library type in Rust and C++.

### 1.2 Terms used throughout

Each term has one sense in this document. Where the code or a quoted document
uses a word differently, the row says so.

| Term | Meaning here |
| --- | --- |
| Target | One of the four native languages. The Go interface named `Target` in the code is the language's syntax builder; we call it the **target builder** |
| Source integration | The package that reads horos for the kernel: it classifies declarations and uses, and gives contributions a typed reader for horos-specific content |
| Use | One occurrence of a type expression in the source, reached under one declaration along one path. In `Page`, `list<T>` in the field `items` is one use, and `T` inside it is another. The code also says "occurrence" |
| Selector | The source integration's exact identity for what a use is: a builtin form (for example `builtin/7:list-of/0:`) or a named declaration (`declaration/<module pin>/<name>`) |
| Native intent | bittype's description of a native type before any syntax exists: a tree of forms with string data and arguments. In the code: `Expression[T]`. It describes types only |
| Form | A constructor of native intent, issued by one target's native model, for example Go's `slice`, `record` and `field`. Only its issuer can interpret it |
| Native model | The package of one target that issues its forms and checks native intent built from them |
| Kernel | bittype's shared generation mechanism: the contract package with the planner, and the registry. It contains no rule about any horos form or any target |
| Planner | The part of the kernel that walks the uses reachable from the requested declarations, selects a rule for each and assembles the result |
| Rule | A function or template that maps the uses of one selector to native intent. The code calls a rule for uses a "binding" (`BindingDescription`). In this document "binding" always means value conversion, as in "binding generation", so we say **use rule** for the code's binding, and **declaration recipe** for the rule that turns a named horos declaration into a native declaration |
| Helper | Today only: a native declaration that is not generated from a horos declaration and that a rule refers to through a key. Go's `BittypeOptional` is the one helper |
| Representation family, short **family** | One kind of horos type together with its native representation: text, lists, closed records, optional fields, integers, unions and so on. The term is the accepted design's |
| Profile | A named promise about what the native types retain. The only implemented one is `structural-representation/1` (the **structural profile**): native types that mirror the horos structure and represent every value of the horos type |
| Exact mapping | A library's own representation of one named horos declaration, for example a library that maps its `Cell<T>` to a class it ships. It takes precedence over the structural profile for that declaration |
| Contribution | One thing a consumer registers on its own: everything one owner publishes for one horos declaration or builtin form, one target and one profile. A family is delivered as one contribution per target. Today's closest thing is a "mapping library" that holds one rule |
| Facet | A part of a contribution that one product reads. The accepted design names a native facet (the type), a value adapter facet (the conversion), a root facet, optional product facets such as an operation facet, and a lowering facet |
| Selection | The planner's choice of one contribution for one use |
| Instance | Only the accepted design's "representation instance": the immutable result of one selection for one use. Our documents also use the word for a coarser run-time identity; we never do. Section 6.4 is about this |
| Plan | The "selected representation plan": the typed record of every selection of one generation request, published for binding backends. It does not exist yet. The Go type `GenerationPlan` in today's code is something else (a closure that builds the output, plus debug strings); we always call it by its code name |
| Plan node | One entry of the plan. Whether there is one per use or one per distinct selection is open (section 6.4) |
| Envelope | The part of a plan node, and of the plan, whose shape is shared by all families: identities, links, versions |
| Payload | The family-owned part of a plan node: an ontos value, typed by a horos type that the family declares, holding the family's decisions for the selected use |
| Decision record | Our word for what a rule returns for one use under the settled shape (section 6.1). The payload is the same content as it appears in the plan |
| Slot | A declared place in a rule's result that stands for the representation of another use. The planner fills it. In the plan a filled slot is a link between two nodes |
| Demand | A request from a contribution to the kernel for a representation that the kernel resolves to whatever the consumer registered. Today there is one kind: "represent this child use" |
| Binding backend, short **backend** | A program, one per target and owned outside bittype, that reads the plan and emits the code that composes conversions |
| Factory | A run-time function, exported by a run-time package, that builds an adapter from a payload and the adapters of the slots |
| Completeness test backend | A backend written in Go inside bittype's tests. It is the acceptance instrument of the plan and no product |
| Consumer composition | The code of a consumer that registers contributions and requests generation |

### 1.3 The stages, in one picture

```text
horos modules
   │  parse, resolve, analyse                      (horos frontend inside bittype)
   ▼
checked semantic graph of the source
   │
   │      registered by the consumer composition:
   │        one source integration   (horos)
   │        contributions            (per target and profile)
   │        one target builder       (per target)
   ▼
kernel: for each use reachable from the requested declarations,
        select one contribution by the use's selector, run its rule,
        resolve the rule's slots by planning the uses they name
   │                                   │
   │ native declarations               │ selected representation plan
   │ (native intent)                   │ (accepted design, not implemented)
   ▼                                   ▼
target builder: syntax tree            binding backend, per target, outside bittype:
   │                                   reads the plan, the source semantics and the
   ▼                                   contributions' factories; emits imports,
Files stage: files, imports            factory calls and glue
(out of scope here)
```

The left branch exists and is what section 4 shows. The right branch is the
accepted design of section 2. The contract in question is the arrow from
"contributions" into the kernel: what a contribution is, what its rule returns
and what it may ask of the kernel.

### 1.4 Scale and direction

- **Targets.** Four. TypeScript and Go are to get conversions. Rust and C++
  keep types only for now, and their binding support is to be reported as
  absent.
- **Families.** Six exist as types only: text, booleans, lists, named closed
  records, optional fields, and generic parameters with their applications.
  Each target has five or six rule packages for them. About ten more families
  are planned, in this order of need: integers, bytes, maps and sets, exact
  values and text refinements, tuples and sequences, unions, nullable values,
  recursion, open records, explicit specialization.
- **Exact mappings.** The service library will contribute `Function`,
  `Stream`, `Service` and `Cell`; Content will contribute two or three domain
  types. These authors are other repositories. The design promise is that
  they need no change to bittype.
- **Backends.** The TypeScript backend belongs to the service library, runs in
  Node and reads a serialized plan. A Go backend "belongs to the future Go
  runtime library", which does not exist. bittype owns only the completeness
  test backend.
- **The first slice** is the next piece of work and the reason for this
  consultation. It completes the six existing families: a recorded decision
  per use, a typed plan, conversion in the completeness test backend for Go,
  and a report of binding support. Its acceptance example is `Page<Person>`
  from horos source to native types, a plan and a round trip of values in Go.
  Before it, the contract between kernel and contribution changes once, and
  every existing rule of the four targets is rewritten onto it with
  byte-identical generated output.
- **Process boundaries.** Kernel, contributions and the Go test backend share
  one Go process, where an identity can be a pointer. The TypeScript backend
  and the TypeScript factories do not: whatever they need must be a value in
  the serialized plan.

## 2. Architectural context

### 2.1 Standing architecture rules (AC = architectural correctness)

The owner's rules for all platform work are numbered AC-01 to AC-10. They are
quoted from the platform repository's `docs/concepts/architectural-correctness.md`.
The questions rest on these seven. Standing: **owner rule**.

| Rule | Requirement, verbatim | Lines |
| --- | --- | --- |
| AC-02 | "Let independently owned implementations supply their domain interpretations, policies and laws. Shared mechanisms own their invariant guarantees and depend on extension contracts; consumer composition selects the concrete contributions." | 102-105 |
| AC-03 | "Within the agreed extension domain, adding a contribution must change only its owning implementation and consumer composition, not shared orchestration, closed case catalogues or unrelated implementations." | 124-126 |
| AC-04 | "Builtins and the first implementation must use the same public extension contract available to independently supplied implementations." | 144-145 |
| AC-05 | "Preserve identity, scope, ownership, authority, assumptions, evidence association, resource limits and failure meaning through composition, as applicable to the contract. Selecting a contribution does not establish the truth of its claims." | 163-166 |
| AC-08 | "When a shared contract really must change, explain the missing invariant or capability, why its responsibility belongs to that layer and how the revised contract serves independently owned contributions." | 225-227 |
| AC-09 | "Whenever we design, redesign or improve an interface, choose the best contract for the domain requirements, responsibilities, ownership, composition and guarantees. The currently implemented surface must have **zero influence** on that choice. Implementation or migration effort has **zero weight**, including as a tiebreaker." | 245-249 |
| AC-10 | "Move a misplaced or wrongly named responsibility to its owner first, as its own change, without waiting for the complete implementation of its destination. That first change establishes the final structure: its inputs and identities have their intended shape, it keeps every guarantee the destination's shared mechanism requires for the cases it admits, and it explicitly refuses the cases it cannot yet check. Missing coverage is then added without moving, renaming or removing what exists." | 286-292 |

Three further passages of the same document shape the questions.

**What violates AC-03, AC-05 and AC-10** (lines 132-133, 172-173, 305-306):

> The addition needs a shared switch/type/registry/bootstrap edit, a weaker fallback or a parallel copy of the orchestration to admit its meaning.

> Openness is achieved by accepting arbitrary payloads, unchecked claims or callbacks that discard the guarantees the mechanism promises.

> Delivering a first change whose input or identity later has to change.

The second sentence matters for question 2: an open set of facets must not
become "arbitrary payloads". The third is the reason for this consultation.

**Wrong coupling, and the demand that corrects it** (owner request of
11 November 2026, section "Package boundaries and forbidden dependencies",
lines 441-454). A table distinguishes two faults. The second:

> Wrong coupling. What it is: A contribution names another contribution that the consumer is meant to choose. What is wrong: The dependency. The caller makes the shared mechanism's choice. Consequence: The mechanism's selection, scope and failure guarantees do not apply to that call. Replacing the callee means editing the caller. Correction: A demand that the shared mechanism resolves.

(The table's cells are joined here into one paragraph.) And below it:

> Not every direct call is wrong. Ask of each call from one part to another whose answer the caller needs: its own, or the composition's. A caller that relies on its own invariant, a shared mechanism or a vocabulary calls directly. A caller that needs whatever the consumer registered for a meaning states a demand and lets the shared mechanism resolve it.

The same section defines the unit that question 1 and the settled part use
(lines 463-466, two list items joined):

> A contribution, meaning one thing a consumer registers or selects on its own, is one unit. It depends on shared contracts, the vocabularies of its own scope and its own private parts. Only a composition depends on a contribution.

**Evidence for an extension boundary** (lines 487-497):

> When qualifying a claimed extension mechanism, implement a materially different contribution using its public contracts in a separately owned package or fixture. Keep the shared mechanism fixed. [...] A renamed copy of the original example, a private test-only constructor or a fixture that activates a preinstalled special case does not establish that.

bittype's own record of applying the package-boundary request,
`docs/package-boundaries.md`, names the one confirmed case of wrong coupling in
the generator. It is the case behind question 1 (lines 96-102, written against
an earlier commit; the Go detail has changed since and is shown in section
4.5):

> In a Record mapping, `r.Mapped(f.Use)` asks the kernel for the native type of a field. Whichever mapping the consumer selected answers. This is the correct form.
>
> Two lines below, the same Record mapping decides how an optional field is represented. Rust builds `std::option::Option`, C++ builds `std::optional`, and Go wraps the type in a helper through a package variable. A consumer who wants another representation of optional fields must replace Record.

Its findings table (line 155) records the case as "Record decides optional
fields itself in Go, Rust and C++", kind "Coupling", "Confirmed; the contract
change is designed", status "Open", tracked by item 6 of issue #196. The same
document says of its dependency check (lines 551-553): "A decision made
inline. Record building `std::optional` imports nothing forbidden. That needs
the contract change of #196 and a reviewer who asks the question above." So no
import rule can find or prevent this fault. Only the contract can.

### 2.2 Binding generation: who owns what

Source: bittype `docs/binding-generation.md`. Its first line gives its
standing: "Accepted direction, 9 November 2026, from research 0014. Not
implemented." Research 0014 was an expert consultation like this one. Standing
of everything quoted in this subsection: **accepted design** (binding
generation, 9 November 2026).

**What it is** (lines 7-15). From a horos contract a build needs more than
native types: "conversion between native values and hydrated message values",
conversion to and from JSON for presentation, operation data, shape checks of
received values and declared operation failures. "That set is **binding
generation**."

**The decision** (lines 19-31; list items joined, links reduced to issue numbers):

> One semantic owner does not require one emitter. It requires an authoritative decision, an explicit interface carrying that decision, and tests that establish that its consumers preserve the promised meaning.
>
> **bittype** remains the authority for selecting native representations. It publishes each selection as a typed **selected representation plan** (#195).
>
> **A binding backend per target** consumes three inputs: the plan; the source semantics; each library's **paired contribution** (#196). From these it emits static composition and typed glue.
>
> **bittype does not emit runtime adapters or own their behaviour.** Its charter gains only the plan interface and the hooks that produce it.

**The owners** (lines 33-41), verbatim:

| Owner | Responsibility |
| --- | --- |
| horos | Declaration semantics, checked identities and applications, descriptors, ground conformance |
| bittype planning | Exact contribution selection, substitutions, references, and the selected representation plan |
| Each representation implementation (a target's structural profile, or a library's exact mapping) | Native layout, presence rules, numeric representation, construction and access |
| Binding backend, one per target | Adapter composition, runtime imports, generated glue, assembly of the requested binding products |
| Contract or domain library | Its native recipe, its adapter implementation, its protocol or domain meaning, and the laws below |
| bitruntime | Capability recognition, export and import identity, lifetime, admission, transport |
| Consumer build | Selection and composition of the concrete contributions it uses |

(bitruntime is the transport runtime below the service library; it plays no
part in this consultation.) Two rows carry the settled part of this document.
The planner owns selection, substitution, references and the plan. The
representation implementation owns "native layout, presence rules [...]
construction and access". Nothing in this table gives the planner a decision
about layout, and nothing gives a representation the right to issue a plan.

**Placement of the backends** (lines 45-46): "TypeScript. The backend belongs
to platform's `libraries/service` (Bitspark/bitspark-service) and consumes a
serialized plan." "Go. A Go backend belongs to the future Go runtime library."

**The plan** (lines 52-72). "A manifest of names is not enough. For each use
and generic instance, the plan carries" a list that includes "the selected
contribution, target and representation profile", "field access and
construction", "the rules for absence, presence, null and native `undefined`",
"recursive references" and "the selected adapter factories, their argument
slots and their resource requirements". Then:

> **Form.** These are typed references and target-owned expressions, the same principle as bittype's native intent (the generation SPI). They are not a universal schema and not another type algebra.
>
> **Serialization.** The plan starts as a typed build-time interface and is serialized only where a real process boundary needs it. A TypeScript backend in Node is such a boundary.
>
> **Completeness test.** The backend can implement each selected representation from the plan, the source semantics and the selected contributions alone, without guessing a rule. If it must decide how Go optional values work, or how a custom union is represented, the plan is incomplete.

And on value-level code (lines 81-83):

> **Interpreters.** A shared structural interpreter is allowed only if it executes the plan's decisions rather than choosing representation rules itself.
>
> **The emitter's scope.** It is scoped separately, with its own syntax and compiler checks. bittype's native vocabularies stay types-only.

"Native vocabularies" are the forms of section 1.2 and the syntax node kinds
of each target. They can express a type and a type declaration. They cannot
express a function, a call, a field access or a constructor. This sentence is
one horn of question 2.

**Paired contributions** (lines 87-99). "A library publishes, for one
declaration and one supported target representation:"

| Element | Meaning |
| --- | --- |
| Source selector | Exact module digest and declaration identity |
| Representation context | Target and supported profile |
| Native facet | Its bittype recipe or native intent |
| Value adapter facet | Run-time factory reference and typed argument-adapter slots, converting the value inside a message |
| Root (attachment) facet | Serving the value as one routed endpoint; Function's is `functionBinding`, its value facet `functionAdapter` |
| Requirements | Adapter-interface compatibility, resource facilities, supported scope |
| Optional product facets | For example an operation facet, through which a callable contract exposes its argument, result and failure positions |

> **Selection.** A native facet and the facets a binding product needs are selected together. A native mapping without an adapter facet still supports types, but a request for bindings that needs a missing facet fails explicitly. There is never a fallback to structural conversion, in the same way that an exact use mapping never silently becomes a structural publication today.

The table names four kinds of facet and ends with "Optional product facets:
for example". It neither closes the list nor says who may add a kind. That is
the other horn of question 2.

**Composition** (lines 105-109): "resolve through the checked graph; preserve
associations through aliases; select by exact identity before any structural
expansion; bind slots by checked parameter identity; reject duplicates,
incompatible profiles, missing facets and unsupported applications." (Five
list items joined.) The design does not say what makes two contributions
duplicates of each other, nor who compares them. That is part of question 3.

**The advice behind these sentences.** Two sentences of the expert's answer to
research 0014 (`research-docs/0014-binding-generation.advice.md`, lines 49 and
123) fix the unit and are quoted in our issues as its source:

> A representation implementation may be a target’s structural profile or a library’s exact custom mapping. Both should use the same public contribution mechanism.

> The decisive property is joint ownership and selection, not whether the contribution is written in Go or metadata. Native package paths and adapter export names should not become part of abstract declaration identity.

**The laws every contribution must satisfy** (lines 121-145). With `V_T` the
admissible abstract values of a type T, `N_T,p` the admissible native values
under profile p, E for encoding and D for decoding:

```text
D(E(n)) ≈ n    for every n in N_T,p   (the equivalence states whether copying, class identity, branch tags and sharing matter)
E(D(v)) = v    for every v in V_T     (abstract-value equality; any normalization is specified)
```

plus conformance, semantic validity, coverage, capture, generic composition,
resource atomicity, capability preservation and cross-target agreement. They
matter here for one reason: a family's decision record has to contain whatever
these laws depend on. Two examples from the first slice: in Go, `Value` is
ignored when `Present` is false, so the equivalence `≈` identifies two absent
values with different `Value`; in TypeScript, "an absent own property" and "a
present `undefined`" are different native states. Neither fact is visible in
the native type.

**Additivity as the acceptance criterion** (lines 257-265). With shared
machinery frozen, a second contract library may change "only its declarations,
contributions, adapters, tests and consumer imports", and Go bindings "only
the Go backend and the Go implementations each semantic owner supplies". "The
unit of implementation is library × target."

### 2.3 Representation families: the contract, the instance and the plan

Source: bittype `docs/representation-families.md`. Its first line: "Accepted
design, 9 November 2026. Not implemented." It was decided by research 0015, a
second consultation, and "details the selected representation plan and the
paired contributions of binding generation for concrete value families".
Where the two documents differ, this one is the later amendment. Standing of
everything quoted in this subsection: **accepted design** (representation
families, 9 November 2026).

**Delivered complete** (lines 17-25):

> A family is delivered with its native type, its plan entries, its value conversion and its qualification together. It is never a types-only mapping that a later binding step reopens. The conversion constrains which native types are admissible at all: a TypeScript `number` cannot hold every horos integer.
>
> **The promise.** Adding a family through the published identity, dependency, invocation, plan-graph and target-extension contracts changes only that family and the consumers that compose it. Existing representation identities keep their meaning.

**The representation contract** (lines 32-42). "A family contributes, per
target, one contract with seven parts:"

| Part | What it defines |
| --- | --- |
| Applicability | The source identities and forms it accepts, its generic assumptions, its targets, its explicit refusals |
| Native representation | Native type, dependencies, construction, access, discrimination, admissible native states |
| Value meaning | Promised canonical domain, domain predicates, native equivalence, any permitted normalization |
| Selected-use payload | A family-owned, versioned schema for the decisions recorded at each use |
| Conversion implementation | Run-time factory reference, child adapter slots, runtime compatibility, conversion laws; optionally a lowering facet (below) |
| Operational requirements | Budget use, mutation and ownership, recursion support, sharing behaviour, failures |
| Qualification | Separate evidence for type generation and binding support, tied to the representation and implementation versions |

The seven parts and the seven elements of section 2.2 describe the same
contribution from two sides. Binding generation lists what a library
publishes. This table lists what the contract must define.

**One instance, two projections** (lines 44-47). This is the sentence the
settled part of this document rests on:

> The planner selects one immutable representation instance per use. The native type and the plan entry are both projections of that instance. No later pass reconstructs a decision from the source, so a native declaration and its conversion cannot disagree.

The expert's answer to research 0015 gave the reason
(`research-docs/0015-representation-families.advice.md`, lines 34 and 36):

> That prevents a subtle failure mode: generating a native declaration, then independently reconstructing conversion decisions from the source. Two passes can read the same declaration and still disagree about a nullable pointer, a union wrapper, or a generic specialization.
>
> “Defined once” should mean that representation decisions have one authoritative origin.

**Factories and the lowering facet** (lines 49-58):

> By default a plan entry names the family's run-time factory and gives it the family's payload as arguments. A family may also supply a target-owned **lowering facet** that turns its recorded decisions into typed target expressions or helpers, such as generated accessors or constructors for Go. The binding backend owns imports, names, placement and assembly. It never switches over families, and source-code strings never stand in for typed target expressions. A shared interpreter is a valid factory implementation only if its handlers are registered contributions executing selected decisions. Generated source is never modified as an extension mechanism.

The advice it comes from (same file, lines 44 and 46) says why a factory with
plain data cannot be the only path:

> A TypeScript runtime can access a field from its name. A Go family may instead need generated accessors, constructors, typed closures, or helper declarations. An external native type may expose construction only through a particular public API.
>
> Provide an extension point through which the selected family can translate its recorded decisions into target-owned expressions or helpers. The binding backend should manage imports, names, references, placement, and assembly. The family’s lowering implementation should supply the family-specific fragments.

**Run-time authority** (lines 60-66): "The representation owner owns its
factories' semantic contract and qualification: bittype's structural profile
for the structural families, a library for its exact mappings. [...]
repository placement does not decide semantic authority. The backend never
decides what a structural integer, nullable field or map means."

**The plan** (lines 70-79), verbatim:

| Shared envelope | Family-owned payload |
| --- | --- |
| Version axes; target and its configuration | Integer representation and admitted range |
| Checked source references and source semantic identity | Record field layout and presence behaviour |
| Root products and source-use identities | Union dispatch and construction |
| Scoped binders and actual-argument environments | Map representation and key requirements |
| Selected contribution and profile identities | Refinements and domain-predicate references |
| Typed dependency slots and recursive references | Family-specific construction and access |
| Native symbols, imports, factory references, runtime requirements | Other representation decisions |
| Diagnostic provenance | Family-specific diagnostic detail |

Followed by five properties (lines 81-102):

> - **A graph.** Generic parameters and recursive adapter references are explicit links in the envelope. Nothing the orchestration must resolve hides inside a payload.
> - **Payload schemas are horos types.** A payload is an ontos value. Its family declares the payload type in horos; bittype generates the Go and TypeScript types, and horos conformance validates a payload. No separate schema language is introduced.
> - **Five version axes:** the envelope protocol; the family payload schema; the representation's semantic identity; generated-code/runtime compatibility; source semantics. An unknown required schema, semantic version or runtime requirement fails explicitly. Only metadata designated nonsemantic may be ignored.
> - **Generic requirements from the first version.** Each generic use carries its scoped binder identities, native type arguments, adapter arguments and requirements (below).
> - **Unstable until six vertical cases pass.** The envelope keeps an explicitly unstable version until the completeness test backend passes:
>   1. recursive `ContentValue` with a shared DAG, a cycle and an expanded-size failure;
>   2. `Page<Outcome>` with optional nullable fields and an exact discriminator;
>   3. the exact external generic `Cell<T>`;
>   4. a Go case with typed construction and a generic adapter argument;
>   5. a generic map whose key requirement is rejected;
>   6. an independently registered family that needs no kernel or backend change.

Three remarks for a reader outside the project.

- The envelope row "Selected contribution and profile identities" requires an
  identity of a contribution that survives serialization. The design does not
  say what it is (question 3).
- The envelope row "Typed dependency slots" requires every link between two
  plan nodes to be visible to the kernel. A family that asks for another
  family's answer must therefore do it in a way that produces a link
  (question 1).
- Vertical case 4 is the lowering facet's test. It is not assigned to a slice
  by this document (question 2).

**Generic requirements** (lines 233-250). "A native type parameter does not
carry the adapter, equality, discriminator or layout evidence needed to bind
values of that type". The table that follows requires, for `Page<T>`, that
"`T`'s selected adapter travels with its native type", and for "Nullable `T`"
that "The null representation cannot collide with a valid `T`". The second is
an example of a requirement that one family places on whatever another family
selected for its operand.

**Absence and null** (lines 164-168), the family behind question 1:

> An optional nullable field has three states. Go uses `Optional[Nullable[T]]` with defined treatment of unused payload fields; a pointer form needs proof that its `nil` cannot collide with a valid value (`**big.Int`). TypeScript uses `field?: T | null`, and the adapter separates an absent own property from a present `undefined` at run time.

**Conversion state** (lines 254-257), the sentence behind question 3:

> Active traversal detects cycles, and completed conversions give sharing. The memo is keyed by the selected representation instance plus input identity, including the descriptor and generic environment. An object-only cache is unsound when one object appears under two mappings.

**Qualification** (lines 303-315). Each family and target is qualified for
bindings separately from types, by seven kinds of evidence, among them "both
laws over the declared domains", "substitution, aliases and recursion" and
"additivity: an unrelated contribution leaves existing selections unchanged".
"Deliberately broken rules show that the tests detect a failure. They do not
prove the laws, so each family keeps its law-level argument beside its tests."

**The one-time change** (line 321, dated 11 November 2026):

> Two existing parts change once, ahead of the families: the structural profile's identity and the closed profile gate (#216, bittype/ACX-008), and what an existing mapping returns (item 6 of #196). After them a family adds a contribution and changes no existing contract.

The first of the two was done on 12 November 2026. The second is the subject
of this consultation.

### 2.4 The work items and the first slice

Four open issues of the bittype repository carry the work. Their texts are
ours. Where they restate an accepted design they add nothing. Where they go
beyond it, standing: **our assumption or derivation**.

**Issue #196, item 6** (added 11 November 2026) is the one-time contract
change. Verbatim, with one stale file path corrected in brackets:

> **The existing mapping contract changes once, here.** [...] A binding's `Expand` returns only a native `Expression[T]` (`contracts/generation/go/description.go:26-31`), and a declaration recipe returns only a `DeclarationDescription[T]` (`:61-66`). A decision such as "this field is optional" is made inside the callback and recorded nowhere ([today `transform/syntax/horos-to-golang/mapping/record/go/record.go:70-72`]). Item 1 of #195 requires the native expression and the plan entry to be projections of one selected instance. So what a contribution returns has to change, and the mappings of all four targets implement that contract. The change is made in this slice for every existing mapping, including Rust and C++, so that later families add contributions and change no existing contract. Derive the new result from the domain (AC-09). This item prescribes no signature.

"Recorded nowhere" is imprecise against the code, and the difference matters
for option C of section 6.1. Section 4.6 shows what is kept and what is lost.

Its verification: "After this slice, a fixture family is added by registering
one contribution. The diff of that test touches no file under
`contracts/generation/go` and no existing mapping."

**Issue #195** is the plan. Its items restate section 2.3: one instance and
two projections; envelope and payloads; a graph; payload schemas in horos ("A
payload that fails its type is refused before any backend runs"); five version
axes; generic requirements from the first version; the completeness test
backend ("A test backend in this repository, in Go, implements each supported
representation from the plan, the source semantics and the selected
contributions alone. A plan that omits a decision fails the test."); a
serialized form; unstable until the six vertical cases pass. It assigns
vertical case 4, the Go case with typed construction, to itself and to no
slice.

**Issue #207** is the first slice. For text, booleans, lists, named closed
records, optional fields and generic parameters with their applications, on
TypeScript and Go (items 1 to 3, 5 and 9 to 11 of eleven; links reduced to
issue numbers; a number in brackets is a row of research 0015's result
table):

> 1. **Contract.** Each family is one representation contract with all seven parts, contributed through the public interface of #196. The kernel, the planner and the completeness test backend gain no case for a family.
> 2. **One instance, two projections.** The native type and the plan entry both come from the one selected instance (#195). Each family declares its payload schema in horos.
> 3. **Conversion.** Go conversion runs in the completeness test backend of #195. The TypeScript factories live in the structural run-time package (bitspark-service#38) and pass this repository's conformance vectors there.
> 5. **Existing output unchanged.** The native declarations generated for today's pagination fixtures stay byte-identical in all four targets. Existing binding identities keep their meaning. The profile's own identity is settled first, in #216.
> 9. **Optional fields.** The TypeScript adapter separates an absent own property from a present `undefined` at run time. Go keeps `BittypeOptional[T]` and ignores `Value` when absent [21].
> 10. **Generics.** For `Page<T>`, the adapter selected for `T` travels with its native type argument. The plan carries scoped binders, native type arguments, adapter arguments and requirements [35].
> 11. **Rust and C++** keep their type support. Their binding support is reported as absent, not implied (#197).

Its verification: "`Page<Person>` goes from horos source to native types, a
plan and a round trip in the Go completeness test backend, with no rule
guessed by the backend." and "A seeded mutant that builds a plan entry in a
separate pass from the native expression turns the completeness test red".
("Binding identities" in item 5 are the rule IDs such as `list/golang/1` of
section 4.2.)

**Issue #197** is the report: binding support per family and target,
separately from type support, shown only with the evidence of section 2.3.

**The order we derived** (standing: our derivation, from AC-10 and the
sentence "change once"):

1. The contribution contract changes, in its final shape, for every existing
   rule of the four targets. Generated output stays byte-identical. Types-only
   behaviour.
2. The optional representation moves out of the record rules, as its own
   change. Until then the wrong coupling of section 2.1 stays an open, recorded
   violation. It cannot move first, because the demand it needs returns an
   answer in the new contract's shape.
3. Everything after is additive: the typed plan at an unstable version,
   payload schemas, the Go conversions in the test backend, conformance
   vectors, the serialized form, the report.

Step 1 is what must not happen twice. Its content is this consultation.

**How the three open points were found.** On 12 November 2026 we first
recorded what was left of item 6 as "details inside a settled design" (#196,
19:38Z). A reading of every accepted document against the first slice then
found three shapes that no document gives, and the last comment on #196
(22:40Z) sends them here. The same day, two design questions on this chain had
been answered from the code and answered wrongly. The reading that produced
this document tried the opposite error, and may read more into the documents
than they say. Please treat every line marked as our derivation accordingly.

### 2.5 The other side of the interface: the service library's binding contracts

The service library wrote three contracts that the plan and the facets must
agree with: bitspark-service `docs/design/binding-contracts.md`, read at commit
`c36722db901ef7a259242e168c4c35d6f777fbf4`. Its status line: "accepted by
libraries/service on 11 November 2026. Nothing here is implemented yet". Its
lifetime laws (L1 to L14) are partly under revision; bittype has confirmed the
groups C, K and F below and not the group L. Standing: **accepted design**
(binding contracts of the service library, 11 November 2026), with L1 marked
unconfirmed.

**What a backend pins.** The service library's issue for its TypeScript
backend (bitspark-service#40, item 2): "It reads one pinned serialized plan
version. It refuses an unknown envelope version, payload schema, semantic
version and run-time requirement, each shown by its own test". Item 3: "It
never selects or re-derives a representation, never switches over families
[...] and contains no Function, Stream, Service or Cell name." Its run-time
package for bittype's structural families (bitspark-service#38, decision of
11 November 2026): "The package makes no representation decision: every choice
arrives as the family's plan payload. Its CI runs bittype's conformance
vectors for each family it implements". So the factory for bittype's list
family in TypeScript is code in another repository, found from a plan node by
an identity, and proved equal in meaning by shared test vectors. No
compiler connects the two.

**Capture** (lines 316 and 318):

> C2. Each family declares, for each native value it reads, one of: it **copies** (the message owns the data), it **retains** a reference, or it needs the value **stable while reading**.

> C4. A decoded native value belongs to its receiver. [...] Decoding may therefore share one native value between two positions only if that value is transitively immutable.

**The call context** (lines 327-329): "There is one call context per message
per side [...] It is owned by this repository's adapter runtime and passed to
every factory." Every conversion charges its work to it (law K1: "One budget
per message per side, shared by all nested conversions.").

**Conversion state and the memo key** (lines 354-356 and 367):

> **The memo key** is the selected representation instance (its descriptor and its generic environment) plus the input's identity. An object-only key is unsound when one object appears under two representations.

> K3. One object under two representations converts once per representation.

**Exports** (line 264, not confirmed by bittype): "L1. Converting a native
capability the sender implements makes one export per message, per selected
representation instance, per native value."

**The service root reads the plan** (lines 187-189). For the contract
`Service<Ops>`, a field of the record `Ops` "is either a namespace, which is a
field the plan selects as the structural record family, or a leaf whose
selected contribution supplies an operation facet". So one contribution's
facet is given the identity of what the plan selected for another use, and has
to recognise "the structural record family" by it. This is a second consumer
of a contribution's identity besides the backend.

**"Instance" in two senses.** Our pull request #260 brings bittype's two
design pages into agreement with these contracts. It is open and not merged.
Its second reading (comment of 12 November 2026, 18:50Z) found, as finding 1:

> The contracts say: "The memo key is the selected representation instance (its descriptor and its generic environment) plus the input's identity". The key includes the input's identity, so two uses share the instance part of the key, not one key.
>
> `docs/representation-families.md:44-45` says "one immutable representation instance per use", and `:259` keys the memo by "the selected representation instance". The page now uses "instance" in two senses.
>
> Suggested: "Two uses that select the same contribution and profile for the same descriptor and generic environment have the same run-time identity. The memo key is that identity plus the input's identity (K3)." Say that this identity is coarser than the per-use instance.

Issue #195 gained the matching criterion on 12 November 2026 (08:11Z): "A
run-time identity for the representation instance. [...] The plan gives each
factory that identity. Two uses with the same selected representation,
descriptor and environment share one memo key."

So three things are called, or keyed by, "the selected representation
instance": the per-use result of a selection (section 2.3), the part of a memo
key that is not the input (K3), and the unit of an export (L1). The first is
one per use. The second and third are one per distinct combination of
contribution, descriptor and generic environment. `People = Page<Person>` used
in two fields is two uses and, at run time, one memo identity.

## 3. Relevant code: the contract as it is today

Everything in this section has the standing **implemented today**. It is shown
so that you can see what a rule can and cannot say now. It is not a constraint
on the answer.

### 3.1 Reading guide

| Subsection | The one thing to take from it | Feeds question |
| --- | --- | --- |
| 3.2 | A rule's whole result is native intent: a tree of target-owned forms, string data, and three kinds of placeholder (a slot, a generic parameter, a helper key) | 1, 4 |
| 3.3 | A registered unit is a profile plus three lists of rules. Its only names are a rule ID that is local to the registration and a registry handle that is a pointer | 3 |
| 3.4 | The registry selects by profile, then the planner selects per use by the use's selector. A replacement can name a use rule only | 1, 3 |
| 3.5 | One Go family of each kind. A rule is a constant template, a template with slots, or a callback that reads the source | 4 |
| 3.6 | The same source fact, "this field is optional", is decided inside each target's record rule, in four different ways. Go alone goes through the kernel, by a key that the consumer wires | 1 |
| 3.7 | Per use the planner keeps the resolved native intent and five debug strings. It drops which rule was selected for which use, the slot-to-use association and the generic convention | 3, 4 |

In the excerpts, `F`, `T` and `R` are type parameters: `F` and `T` are marker
types of the source and target language, `R` is the type of the source
integration's reader. Import aliases follow one pattern: `g` is the kernel's
contract package (`contracts/generation/go`), `h` the horos syntax vocabulary,
`s` the syntax vocabulary of the target at hand, `model` that target's native
model, `source` the horos source integration, `rules` the small common package
of the horos rule side (the profile names and one function type), `patterns`
the common source patterns, `naming` the common naming function, `ast` the
syntax tree package (an `ast.NodeRef` is a site in a syntax tree, used for
diagnostics), `representation` the package of language identities. Elisions
are marked `// ...`.

### 3.2 What a rule can say: native intent

`contracts/generation/go/mapping.go:5-12`, `expression.go:9-14, 30-35, 41-44,
50-90` (accessors omitted):

```go
// Selector is an exact source-owned identity, qualified by language/profile/marker.
// Its key is never interpreted by the composer. Display names and graph-local IDs
// are not substitutes for a source integration's stable declaration identity.
type Selector struct {
	language representation.Language
	key      string
}

// Form is a target-owned native semantic constructor. Only its issuer can
// interpret its key. This is not an AST kind or a globally registered type kind.
type Form[T any] struct {
	owner *identity
	name  string
}

// Slot identity is independent of its display name. One description declares
// its slots, and captures supply the exact source occurrences for them.
type Slot struct {
	owner *identity
	name  string
}

type HelperKey struct {
	owner *identity
	name  string
}

type ExpressionKind uint8

const (
	NativeExpression ExpressionKind = iota + 1
	CapturedExpression
	ParameterExpression
	ArgumentExpression
	HelperExpression
	DeclarationExpression
)

// Expression is an immutable template. Parameter and Argument indices only have
// meaning in their declaration's single generic convention. DeclarationExpression
// is issued by the planner after it authenticates the source callee.
type Expression[T any] struct {
	origin      ast.NodeRef
	kind        ExpressionKind
	form        Form[T]
	data        []string
	args        []Expression[T]
	slot        Slot
	index       int
	helper      HelperKey
	declaration uint32
}

func Native[T any](form Form[T], data []string, args ...Expression[T]) Expression[T] { /* ... */ }
func Mapped[T any](slot Slot) Expression[T]                                          { /* ... */ }
func Parameter[T any](index int) Expression[T]                                       { /* ... */ }
func Argument[T any](index int) Expression[T]                                        { /* ... */ }
func Helper[T any](key HelperKey, args ...Expression[T]) Expression[T]               { /* ... */ }
```

`identity` is a private struct whose address is the identity: `Form`, `Slot`
and `HelperKey` values are equal only if they come from the same call of their
constructor in the same process. That makes them unforgeable, and it makes
them meaningless outside the process.

The six kinds, in the words of this document:

- **Native**: a form of the target's native model with string data and
  argument expressions. `model.Primitive("string")` is the Go form `primitive`
  with data `["string"]`.
- **Captured** (`Mapped`): a slot. The planner replaces it with the native
  intent it plans for the use bound to that slot. This is today's one demand.
- **Parameter** and **Argument**: a generic parameter of the native
  declaration being described, and an actual argument at an application.
- **Helper**: a reference, through a key, to a helper declaration that some
  registered unit must supply. The planner replaces it with a reference to
  that declaration.
- **Declaration**: a reference to a planned native declaration. Only the
  planner creates it.

A native model is a handful of constructors over forms, plus a check. Go's,
complete except for the check (`transform/syntax/golang/model/go/model.go:12,
35-51`):

```go
var forms = g.NewFormIssuer[s.Language]()

func Form(name string) g.Form[s.Language]      { return forms.Form(name) }
func Name(f g.Form[s.Language]) (string, bool) { return forms.Read(f) }
func Primitive(name string) g.Expression[s.Language] {
	return g.Native(Form(FormPrimitive), []string{name})
}
func Collection(element g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormSlice), nil, element)
}
func Field(name string, optional bool, ty g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormField), []string{name, strconv.FormatBool(optional)}, ty)
}
func Record(fields ...g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormRecord), nil, fields...)
}
func External(module, name string, args ...g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormExternal), []string{module, name}, args...)
}
```

The Go model has seven forms: `primitive`, `slice`, `field`, `record` and
`external` for expressions, `alias` and `defined` for declarations. It admits
two primitives, `string` and `bool`, and its check refuses every `external`
expression. The other three models are parallel. All four describe types and
type declarations only.

### 3.3 What a registered unit is

`contracts/generation/go/description.go:13-31, 61-84` (one tool annotation
comment omitted):

```go
type Capture[F any] struct {
	Slot Slot
	Use  Use[F]
}

// Pattern fixes selection and declares its captures. Capture reads the source
// meaning of the kernel's view through the source reader R and names children
// of that view; the kernel checks that each captured use is one of them.
type Pattern[F, R any] struct {
	Selector Selector
	Slots    []Slot
	Capture  func(context.Context, R, UseView[F]) ([]Capture[F], error)
}
type BindingDescription[F, T, R any] struct {
	ID     string
	Input  Pattern[F, R]
	Output Expression[T]
	Expand func(context.Context, Expansion[F, T, R]) (Expression[T], error)
}

type DeclarationRecipe[F, T, R any] struct {
	ID string
	// A zero selector explicitly supplies the structural declaration policy.
	Selector Selector
	Expand   func(context.Context, DeclarationRequest[F, T, R]) (DeclarationDescription[T], error)
}
type HelperRecipe[T any] struct {
	Key         HelperKey
	Declaration NativeDeclaration[T]
}

type MappingDescription[F, T, R any] struct {
	Profile      string
	Bindings     []BindingDescription[F, T, R]
	Declarations []DeclarationRecipe[F, T, R]
	Helpers      []HelperRecipe[T]
}

// Describe is called at registration. The registry freezes all returned data.
// No provider receives a target AST builder or is called during Run.
type MappingLibrary[F, T, R any] interface {
	Describe(from, to representation.Language) (MappingDescription[F, T, R], error)
}
```

`Use[F]` is the kernel's handle for one use, and `UseView[F]` its structural
view of that use: its kind (a reference to a generic parameter, a reference to
or application of a declaration, or a builtin form), its children, and whether
it is the whole body of its declaration. Only the kernel creates them.

A use rule (`BindingDescription`) has exactly one of two results: `Output`, a
constant template, or `Expand`, a callback. The callback receives an
`Expansion`, which offers three operations and nothing else (signatures,
`description.go:113-136`):

```go
// Reader is the source's interpretation reader, available during the callback.
func (r Expansion[F, T, R]) Reader() (R, error)
// Use is the kernel's view of the occurrence a binding rule maps.
func (r Expansion[F, T, R]) Use() (UseView[F], error)
// Mapped captures a use for the planner to represent after the callback. A
// binding rule may capture only children of its occurrence; a declaration
// recipe only its body.
func (r Expansion[F, T, R]) Mapped(use Use[F]) (Expression[T], error)
```

`Mapped` returns a placeholder at once. The planner plans the child use after
the callback has returned. So a callback never sees what was selected for a
child, and the native intent of a child is never an input of the parent's
rule. The callback's whole output is one `Expression[T]`.

A declaration recipe returns a `DeclarationDescription`: the native name, a
declaration form, the correspondence between the horos generic parameters and
the native ones (the "generic convention"), and a body expression.

Since 11 November 2026 every registered unit holds exactly one entry in one of
the three lists; the lists are a remainder of an earlier shape in which one
unit held all rules of a target. Since 12 November 2026 each unit writes its
own profile into `Profile`. The common horos package gives every unit the same
Go shape, a function (`transform/syntax/horos/mapping/go/mapping.go:28, 40`):

```go
const RepresentationProfile = "structural-representation/1"

type Library[T any] func(from, to representation.Language) (g.MappingDescription[h.Language, T, source.Reader], error)
```

### 3.4 Registration, selection by profile, replacement

`transform/syntax/shared/generation/go/registry.go:19-30, 111-130`:

```go
type MappingRef struct {
	owner *Registry
	id    uint64
}
type BindingRef struct {
	Provider MappingRef
	ID       string
}
type Replacement struct {
	Selector       g.Selector
	Expected, With BindingRef
}

func RegisterMapping[F, T, R any](r *Registry, from, to representation.Language, p g.MappingLibrary[F, T, R]) (MappingRef, error) {
	// ... argument checks
	// Provider code runs outside the registry lock and only at registration.
	desc, err := p.Describe(from, to)
	if err != nil {
		return MappingRef{}, err
	}
	desc, err = g.FreezeDescription(from, desc, g.DefaultLimits())
	if err != nil {
		return MappingRef{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	ref := MappingRef{r, r.next}
	r.mappings = append(r.mappings, mappingEntry{from, to, reflect.TypeFor[R](), desc.Profile, desc, plan[F, T, R], ref})
	return ref, nil
}
```

`FreezeDescription` copies the lists and checks them: a non-empty profile,
rule IDs unique within the one description, valid selectors of the source
language, slots declared before use, and size limits. It checks nothing across
two registered units.

A request names a profile, and the registry plans with every unit that states
it (`registry.go:184-188, 194-196`):

```go
	for _, m := range r.mappings {
		if m.profile == q.Profile && m.from.Compatible(q.Input.Source.Language()) && m.to.Compatible(q.To) {
			mappings = append(mappings, m)
		}
	}
	// ...
	if len(mappings) == 0 {
		return fail(g.ErrUnsupported)
	}
```

It hands them to the planner under names it makes from its registration
counter (`registry.go:242`), and the planner prefixes each rule's ID with that
name (`contracts/generation/go/planning.go:142-153`):

```go
		parts = append(parts, g.SelectedDescription[F, T, R]{ID: fmt.Sprintf("library-%d", e.ref.id), Description: d})
```

```go
		for _, b := range part.Description.Bindings {
			p.bindings[b.Input.Selector.Key()] = append(p.bindings[b.Input.Selector.Key()], selectedBinding[F, T, R]{part.ID + "/" + b.ID, b})
		}
		for _, r := range part.Description.Declarations {
			p.recipes[r.Selector.Key()] = append(p.recipes[r.Selector.Key()], selectedRecipe[F, T, R]{part.ID + "/" + r.ID, r})
		}
		for _, h := range part.Description.Helpers {
			if _, ok := p.helpers[h.Key]; ok {
				return fail(ErrAmbiguous)
			}
			p.helpers[h.Key] = h
		}
```

So a rule is known to the planner as, for example, `library-3/list/golang/1`.
The first part depends on the order in which the consumer registered its
units. The second part is a string its author chose. The page that documents
the implemented generator says of it (`docs/type-generation-plugins.md:203-205`):
"A binding's ID is local to its
mapping registration. It is not a semantic proof, a globally unique
implementation revision or a configuration digest. There is no PluginID API."

Per use, the planner takes the rules filed under the use's selector key. More
than one is an error, never a preference (`planning.go:465-475`, inside the
function that plans one use):

```go
	choices := p.bindings[key]
	if len(choices) > 1 {
		return Expression[T]{}, fmt.Errorf("%w: %s", ErrAmbiguous, key)
	}
	if len(choices) == 1 {
		return p.bind(choices[0], key, view, f, depth)
	}
	if view.Kind != DeclarationUse {
		return Expression[T]{}, p.walk.unsupported(u.o.term, u.o.declaration, "no mapping for "+key)
	}
	id, err := p.declare(view.Callee, depth+1)
```

A use rule for a named declaration is an exact mapping and wins. Without one,
a use of a named declaration is planned as a reference to the native
declaration that a declaration recipe produces. A recipe with a zero selector
is the default for every declaration.

A consumer resolves an ambiguity with `Registry.Replace(Replacement)`. A
replacement names the selector and two use rules, each by a registry handle
and its local ID. The handle is a pointer to the registry plus a counter, so
it exists only in the process and registry that issued it. The registry
removes the `Expected` rule before planning. Only use rules can be replaced.
A declaration recipe cannot, and two units that supply a helper under the same
key fail as ambiguous with no way to choose.

### 3.5 One Go family of each kind

Each of the following is the whole rule package of one family for Go, minus
its package clause, imports and comments.

**Text: a constant template.**
`transform/syntax/horos-to-golang/mapping/text/go/text.go:16-30`:

```go
func primitive(sourceName, targetName string) rules.Library[s.Language] {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		pattern, err := patterns.Head(from, "scalar", sourceName)
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Profile: rules.RepresentationProfile, Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: sourceName + "/golang/1", Input: pattern, Output: model.Primitive(targetName)}}}, nil
	}
}

func Library() rules.Library[s.Language] {
	return primitive("utf8-text", model.PrimitiveString)
}
```

`patterns.Head(from, "scalar", "utf8-text")` returns a `Pattern` whose
selector is the horos source integration's identity for the builtin form
`scalar` with discriminator `utf8-text`. The rule's ID is `utf8-text/golang/1`.
Nothing about Go strings beyond the type name is stated anywhere: not that a
Go string may hold invalid UTF-8, and not what a conversion does with one.

**List: a template with one slot.**
`transform/syntax/horos-to-golang/mapping/list/go/list.go:18-27`, and the
pattern it uses, `transform/syntax/horos/mapping/pattern/go/pattern.go:27-41`:

```go
func Library() rules.Library[s.Language] {
	slot := g.TypeSlot("E")
	return func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		pattern, err := patterns.Element(from, slot)
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Profile: rules.RepresentationProfile, Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: "list/golang/1", Input: pattern, Output: model.Collection(g.Mapped[s.Language](slot))}}}, nil
	}
}
```

```go
func Element(from representation.Language, slot g.Slot) (g.Pattern[h.Language, source.Reader], error) {
	p, err := Head(from, "list-of", "")
	if err != nil {
		return p, err
	}
	p.Slots = []g.Slot{slot}
	p.Capture = func(ctx context.Context, src source.Reader, u g.UseView[h.Language]) ([]g.Capture[h.Language], error) {
		operand, err := src.Element(u)
		if err != nil {
			return nil, err
		}
		return []g.Capture[h.Language]{{Slot: slot, Use: operand}}, nil
	}
	return p, nil
}
```

The slot `E` is bound to the list's element use by the source integration's
reader, and the template says "a Go slice of whatever is selected for `E`".
This is the correct form of a demand: the list rule does not know what
represents its element. Whether a nil slice and an empty slice are the same
value is stated nowhere.

**Declarations: a recipe that all four targets share.** The horos side owns
it; a target passes in its two declaration forms.
`transform/syntax/horos/mapping/declarations/go/declarations.go:26-73`,
shortened where marked:

```go
func Recipe[T any](record, alias g.Form[T]) g.DeclarationRecipe[h.Language, T, source.Reader] {
	return g.DeclarationRecipe[h.Language, T, source.Reader]{ID: "declarations", Expand: func(ctx context.Context, r g.DeclarationRequest[h.Language, T, source.Reader]) (g.DeclarationDescription[T], error) {
		// ... src is the reader, open the kernel's view of the declaration,
		// decl the source integration's reading of it (error handling omitted)
		name, err := naming.Exported(decl.Name)
		// ...
		if !decl.Root {
			name = "M" + decl.ModulePin[:16] + "_" + name
		}
		body, descriptor, err := src.Descriptor(open.Body)
		// ...
		form := alias
		if descriptor && body.Constructor == "record" {
			form = record
		}
		desc := g.DeclarationDescription[T]{Name: name, Form: form, Origin: decl.Origin}
		for i, f := range open.Formals {
			parameter, err := src.Parameter(f)
			// ...
			pname, err := naming.Exported(parameter.Name)
			// ...
			desc.Generic.Parameters = append(desc.Generic.Parameters, g.ParameterDescription{Name: pname, Origin: parameter.Origin, RequiredUse: true})
			desc.Generic.Bindings = append(desc.Generic.Bindings, g.Parameter[T](i))
			desc.Generic.Arguments = append(desc.Generic.Arguments, g.Argument[T](i))
		}
		desc.Body, err = r.Mapped(open.Body)
		return desc, err
	}}
}
```

The recipe states no profile and has a zero selector. Go's unit wraps it:
`Declarations: []g.DeclarationRecipe[...]{recipe.Recipe(model.Form(model.FormDefined), model.Form(model.FormAlias))}`
under `Profile: rules.RepresentationProfile`
(`transform/syntax/horos-to-golang/mapping/declarations/go/declarations.go:21-24`).
"Generics" as a family is this recipe's `Generic` part together with the
planner's handling of parameters and applications. It makes two decisions
inline that a conversion needs: whether the declaration becomes a named record
type or an alias, and how each horos parameter corresponds to a native one.

### 3.6 Where "this field is optional" is decided, in four targets

An optional field is not a horos type. The horos source integration reads a
record as a list of fields, each with a name, a flag and one child use, the
field's type (`transform/syntax/horos/source/go/read.go:25-30, 77-79`):

```go
type Field struct {
	Name     string
	Optional bool
	Use      g.Use[h.Language]
	Origin   ast.NodeRef
}

// Fields reads a closed or open record payload (horos/record, kind, field...).
// Each field is (horos/field, name, req|opt, type); field k's type is the
// operand at locator [k+2, 3].
```

Only the field's type is a use that the kernel walks and selects for. The
flag `opt` is content of the record form. No selector exists for "an optional
field", so nothing can be registered for it today.

**Go.** The record rule is a callback. It asks the kernel for each field's
type and then wraps an optional one in a helper
(`transform/syntax/horos-to-golang/mapping/record/go/record.go:31-78`, error
handling shortened):

```go
func Library(optional g.HelperKey) rules.Library[s.Language] {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		if !optional.Valid() {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, g.ErrInvalid
		}
		pattern, err := patterns.Head(from, "record", "closed")
		// ...
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Profile: rules.RepresentationProfile, Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: "record/golang/1", Input: pattern, Expand: func(ctx context.Context, r g.Expansion[h.Language, s.Language, source.Reader]) (g.Expression[s.Language], error) {
			src, err := r.Reader()
			// ...
			use, err := r.Use()
			// ...
			if !use.Body {
				return g.Expression[s.Language]{}, src.Unsupported(use.Use, "nested records require a named declaration in this profile")
			}
			view, _, err := src.Descriptor(use.Use)
			// ...
			fields, err := src.Fields(use)
			// ...
			var members []g.Expression[s.Language]
			for _, f := range fields {
				f.Name, err = naming.Exported(f.Name)
				// ...
				ty, err := r.Mapped(f.Use)
				// ...
				if f.Optional {
					ty = g.Helper(optional, ty)
				}
				members = append(members, g.At(model.Field(f.Name, f.Optional, ty), f.Origin))
			}
			return g.At(model.Record(members...), view.Origin), nil
		}}}}, nil
	}
}
```

The helper is supplied by a second unit, which is the whole "optional fields"
family of Go today
(`transform/syntax/horos-to-golang/mapping/optional/go/optional.go:21-28`):

```go
func Library(key g.HelperKey) rules.Library[s.Language] {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Profile: rules.RepresentationProfile, Helpers: []g.HelperRecipe[s.Language]{{Key: key, Declaration: g.NativeDeclaration[s.Language]{
			Name: "BittypeOptional", Form: model.Form(model.FormDefined), Parameters: []g.NativeParameter{{Name: "T"}},
			Body: model.Record(model.Field("Present", false, model.Primitive(model.PrimitiveBool)), model.Field("Value", false, g.Parameter[s.Language](0))),
		}}}}, nil
	}
}
```

Neither package imports the other. The consumer composition creates the key
and gives it to both (`integration/go/libraries_golang_test.go:25-35`):

```go
func golangLibraries() []g.MappingLibrary[h.Language, s.Language, source.Reader] {
	optionalKey := g.NewHelperKey("horos/optional")
	return []g.MappingLibrary[h.Language, s.Language, source.Reader]{
		text.Library(),
		boolean.Library(),
		list.Library(),
		record.Library(optionalKey),
		declarations.Library(),
		optional.Library(optionalKey),
	}
}
```

What this arrangement gives and does not give:

- The kernel resolves the helper. A key that no registered unit supplies makes
  planning fail with "generation capability unavailable" at the first optional
  field, and an optional field is never emitted unwrapped.
- A consumer who wants another helper declaration registers another unit under
  the same key. Selection is by wiring: the key is a pointer identity created
  by the composition, it is no selector, `Replace` does not apply to it, and it
  cannot be written into a plan.
- The record rule still decides that an optional field is a type wrapped
  around the field's type, with one type argument. A representation that is
  not of that shape (a pointer, a field pair, a modifier) is out of reach
  without editing the record rule.
- The page that documents this form calls it interim
  (`docs/semantic-generation.md:160-162`): "This form is interim: #196, item 6,
  changes how a record asks for the optional representation."

**C++ and Rust.** The record rule builds the wrapping type itself. The loop
body of C++'s rule
(`transform/syntax/horos-to-cpp/mapping/record/go/record.go:49-56`):

```go
				ty, err := r.Mapped(f.Use)
				if err != nil {
					return g.Expression[s.Language]{}, err
				}
				if f.Optional {
					ty = model.External("std", "optional", ty)
				}
				members = append(members, g.At(model.Field(f.Name, f.Optional, ty), f.Origin))
```

Rust's is the same with `model.External("std::option", "Option", ty)`. The
pair of strings `"std"`, `"optional"` is written here and compared in the C++
target builder, which adds the `#include`
(`transform/syntax/cpp/target/go/provider.go:190-196`):

```go
	case model.FormExternal:
		if data[0] == "std" && data[1] == "optional" {
			if err := b.Include("optional"); err != nil {
				return m.Fragment[s.Language]{}, err
			}
		}
		return b.Path(append(strings.Split(data[0], "::"), data[1]), args, origin, true)
```

So in C++ the knowledge "an optional field is `std::optional`" is in two
packages that agree by spelling, and neither is a contribution a consumer can
replace.

**TypeScript.** The record rule wraps nothing. The loop body is
`ty, err := r.Mapped(f.Use)` followed by
`members = append(members, g.At(model.Field(f.Name, f.Optional, ty), f.Origin))`.
The flag travels as the second string of the `field` form, and the TypeScript
target builder turns it into the property's `?`
(`transform/syntax/typescript/target/go/provider.go:133-134`):

```go
	case model.FormField:
		return b.Node(m.Node(s.KindProperty, m.Attributes("name", data[0], "optional", data[1]), origin), m.Child[s.Language]{Role: s.RoleType, Fragment: args[0]})
```

All four models have the same `Field(name, optional, ty)` constructor. Go's,
Rust's and C++'s target builders ignore its flag and render the wrapped type;
TypeScript's renders the flag and has no wrapped type.

**The four answers side by side.**

| Target | Native form of `next?: text` | What the representation is | Who decides today |
| --- | --- | --- | --- |
| TypeScript | `next?: string` | A property of the field member. No type stands for it | The record rule, through a flag that the target builder renders |
| Go | `Next BittypeOptional[string]` | A generated generic type around the field's type | The record rule decides "a wrapper type"; a helper unit, wired by key, supplies the declaration |
| Rust | A field `next` of type `::std::option::Option<::std::string::String>` | An external type around the field's type | The record rule |
| C++ | A field `next` of type `::std::optional<::std::string>` | An external type around the field's type, plus an `#include` | The record rule and the target builder, agreeing by spelling |

Later families raise the same shape of question: a nullable value (a wrapper,
a pointer, or `T | null`), a map's requirement on its key, a union's
discriminator, an optional nullable field with three states.

### 3.7 What the planner keeps per use, and what it drops

The planner's output types (`contracts/generation/go/planning.go:14-22,
33-43`, and `native.go:8-16`):

```go
type PlannedDeclaration[T any] struct {
	ID          uint32
	Description NativeDeclaration[T]
	Exposed     bool
}
type Program[T any] struct {
	declarations []PlannedDeclaration[T]
	exposed      []uint32
}

type Step struct {
	Binding, Source, Declaration, Form string
	Expression                         string
}

// GenerationPlan is issued only after all reachable rules and dependencies have
// been resolved. Run does not retain a source, reader or mapping expansion.
type GenerationPlan struct {
	steps []Step
	run   func(context.Context) (Output, error)
}

// NativeDeclaration is target intent after source mapping and substitution.
// It has no source formal bindings, argument templates or mapping callbacks.
type NativeDeclaration[T any] struct {
	Name       string
	Form       Form[T]
	Origin     ast.NodeRef
	Parameters []NativeParameter
	Body       Expression[T]
}
```

`Program[T]` is the closed list of native declarations. It reaches the target
builder and nothing else. `GenerationPlan` is a closure over it plus the
`Step` list, which exists for inspection in tests.

Applying a use rule ends like this (`planning.go:560-565`; `captures` maps the
rule's slots to uses, `f` is the generic environment of the declaration being
planned):

```go
	out, err := p.resolve(expression, captures, f, depth+1)
	if err != nil {
		return Expression[T]{}, err
	}
	p.steps = append(p.steps, Step{Binding: b.name, Source: key, Form: out.form.Name(), Expression: describe(out)})
	return out, nil
```

`resolve` walks the rule's result and replaces every slot with the result of
planning the use bound to it, and every helper reference with a reference to
the helper's declaration. The returned expression is embedded in its parent's
expression. `describe` prints an expression as a string. When a declaration
is finished, the planner stores it through a function whose comment states
what is dropped (`planning.go:756-764`):

```go
// nativeDeclaration deliberately drops source binding/application templates.
// Construction can only observe the native declaration selected by planning.
func nativeDeclaration[T any](d DeclarationDescription[T]) NativeDeclaration[T] {
	out := NativeDeclaration[T]{Name: d.Name, Form: d.Form, Origin: d.Origin, Body: d.Body}
	for _, p := range d.Generic.Parameters {
		out.Parameters = append(out.Parameters, NativeParameter{Name: p.Name, Origin: p.Origin})
	}
	return out
}
```

So, per use:

| Fact | After planning |
| --- | --- |
| The native intent selected for the use | Kept, embedded in the parent's expression inside `Program[T]`. Readable only by the target's native model, which issued its forms |
| Which rule was selected | Kept only as the string `Step.Binding`, for example `library-3/list/golang/1`. The step list is flat: a step is linked neither to its use nor to the step of the parent use |
| Which use it was | Dropped. `Step.Source` holds the selector key, which is the same for every list in the program |
| Which child use filled which slot | Dropped when `resolve` returns |
| A use of a generic parameter | No step at all. The planner substitutes from the environment |
| The generic convention (which horos parameter is which native parameter, and how an application's arguments are ordered) | Dropped by `nativeDeclaration`. The application keeps only its resolved native arguments |
| "This field is optional" | Kept twice, in target-owned form: as the string `"true"` in the `field` form's data, and (Go, Rust, C++) as the wrapping type. Not kept as a decision: nothing says that a Go `BittypeOptional` with `Present` false ignores `Value`, or that a TypeScript reader must tell an absent property from `undefined` |
| A refusal ("nested records require a named declaration in this profile") | An error from inside the callback. Not declared anywhere before the use is met |

The issue text of section 2.4 says such a decision is "recorded nowhere". The
table gives the precise state: the flag survives as an untyped string inside
native intent that only one target's model may interpret, and its meaning for
values survives nowhere.

## 4. One request, today and under the accepted design

The request: Go types for `People` from the module of section 1.1, with the
six Go units of section 3.6 registered. This trace is read from the code of
section 3; we did not run it for this document.

**Today.**

1. The registry collects the six units, since all state
   `structural-representation/1`, and calls the planner.
2. The planner starts at the declaration `People`. No use rule is filed under
   its selector, so the default declaration recipe runs. It returns the name
   `People`, the form `alias` and a slot for the body.
3. The body is a use of `Page` applied to `Person`. No use rule is filed under
   `Page`'s selector, so the planner plans the declaration `Page` with the
   same recipe: name `Page`, form `defined`, one parameter `T`, a slot for the
   body.
4. `Page`'s body is a use of the builtin form "closed record". The record rule
   is filed under that selector. Its callback reads two fields, asks for a slot
   per field type, wraps the second in the helper reference and returns
   `record(field(Items, false, slot), field(Next, true, helper(key, slot)))`.
5. The planner fills the first slot by planning the use `list<T>`: the list
   rule's template `slice(slot E)`, whose slot is the use `T`, a generic
   parameter, replaced from the environment. It fills the second slot by
   planning the use `text`: the constant `primitive(string)`. It resolves the
   helper reference by planning the helper declaration `BittypeOptional` once
   and referring to it with the argument `primitive(string)`.
6. Back at `People`: the actual argument `Person` is a use of a declaration,
   planned like `Page`. The body of `People` becomes a reference to the
   declaration `Page` with one argument, the reference to `Person`.
7. The result is a `Program` of four native declarations (`People`, `Page`,
   `BittypeOptional`, `Person`) and a flat list of step strings. The target
   builder turns the program into a Go syntax tree.

Eight selections were made (the default recipe three times, the record rule
twice, the list rule once, the text rule twice) and one helper was resolved. None of them
exists as a value afterwards.

**Under the accepted design** the same walk has to leave a plan. From the
sentences of sections 2.2 and 2.3, for this request it holds at least:

- for each use, which contribution and profile were selected, by an identity
  a backend in another process recognises;
- for the record of `Page`: the field layout (source names, native names,
  order) and, for `next`, the presence behaviour, as a payload typed by a
  horos type that the record family (or the optional family) declares;
- links, visible to the kernel, from the record's node to the nodes selected
  for `list<T>` and `text`, and from the list's node to "whatever is selected
  for `T`";
- for the application `Page<Person>`: the scoped binder `T`, the native type
  argument `Person` and the adapter argument, so that "`T`'s selected adapter
  travels with its native type";
- the five version axes, the native symbols, the factory references.

And the Go native type of `Page`, byte for byte as today, must be a projection
of the same selections. The completeness test backend then converts a Go
value of type `Page[Person]` to an ontos value and back, using the plan, the
source semantics and the registered contributions, "without guessing a rule".

**What has no place to come from today.** The trace shows the three open
points concretely.

- At step 4 the record rule decides the optional representation itself. Under
  the design it must ask. There is no use to ask about, only a flag
  (question 1).
- At step 7 the Go test backend must construct a value of the named struct
  type `Page[Person]`. A factory that receives field names as data can do that
  only through reflection. The design's alternative is a lowering facet that
  yields typed Go code, for which bittype has no vocabulary (question 2).
- `library-3/list/golang/1` is the only name the selected list rule has. Its
  first part is a registration counter (question 3).
