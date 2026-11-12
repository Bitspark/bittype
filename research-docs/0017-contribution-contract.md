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
