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
   the word "instance" for more than one of them. A fourth joined them on
   12 November 2026: the identity of a generated native declaration, whose
   name and import path are now decided only after generation.

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

The code in section 3 is evidence of the present state. It is never a reason
for or against an option.

**How this document is organised.** Section 1 introduces the system and its
terms. Section 2 gives the architecture rules, the accepted designs, the work
items and the contracts of the library that consumes the plan. Section 3
shows the contract as it is today, with the code inlined, and section 4
traces one request. Section 5 lists the fixed points with their standing.
Section 6 holds what we have considered, each point with options, our leaning
and what speaks against it. Section 7 holds the questions.

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
| Contribution | One thing a consumer registers on its own: everything one owner publishes for one horos declaration or builtin form, one target and one profile. A family is delivered as one contribution per target. What exists today is a "mapping library" that holds one rule; quoted documents call it a "mapping" (the "Record mapping"), and section 3 calls it a **registered unit** to keep it apart from the contribution of the design |
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
| Files stage, placement | The stage after generation that puts generated declarations into files, packages or modules according to a **placement** the consumer supplies, and that realizes each declaration's final name (section 2.6) |

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
Files stage: files, imports,           factory calls and glue
realized names (section 2.6)
```

The left branch exists and is what section 3 shows. The right branch is the
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
  types. These authors are other repositories, and the design promises that
  they need no change to bittype.
- **Backends.** The TypeScript backend belongs to the service library, runs in
  Node and reads a serialized plan. The owner of a Go backend, a Go runtime
  library, does not exist yet. bittype owns only the completeness test
  backend.
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
3.6):

> In a Record mapping, `r.Mapped(f.Use)` asks the kernel for the native type of a field. Whichever mapping the consumer selected answers. This is the correct form.
>
> Two lines below, the same Record mapping decides how an optional field is represented. Rust builds `std::option::Option`, C++ builds `std::optional`, and Go wraps the type in a helper through a package variable. A consumer who wants another representation of optional fields must replace Record.

Its findings table (line 155) records the case as "Record decides optional
fields itself in Go, Rust and C++", kind "Coupling", status "Open", tracked by
item 6 of issue #196. Of its dependency check it says (lines 551-553): "A
decision made inline. Record building `std::optional` imports nothing
forbidden. That needs the contract change of #196 and a reviewer who asks the
question above." So no import rule can find or prevent this fault. Only the
contract can.

### 2.2 Binding generation: who owns what

Source: bittype `docs/binding-generation.md`. Its first line gives its
standing: "Accepted direction, 9 November 2026, from research 0014. Not
implemented." Research 0014 was an expert consultation like this one. Standing
of everything quoted in this subsection: **accepted design** (binding
generation, 9 November 2026).

**What it is** (lines 7-15). From a horos contract a build needs more than
native types: "conversion between native values and hydrated message values"
(message values in which references to remote endpoints are live objects),
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

**The plan** (lines 52-72). "A manifest of names is not enough. For each use
and generic instance, the plan carries" a list that includes "the selected
contribution, target and representation profile", "field access and
construction", "the rules for absence, presence, null and native `undefined`",
"recursive references" and "the selected adapter factories, their argument
slots and their resource requirements". Its form is "typed references and
target-owned expressions", and it "starts as a typed build-time interface and
is serialized only where a real process boundary needs it. A TypeScript
backend in Node is such a boundary." Then:

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

**The advice behind these sentences.** The expert's answer to research 0014
(`research-docs/0014-binding-generation.advice.md`, line 123) names the
property that fixes the unit, and our issues quote it as the source:

> The decisive property is joint ownership and selection, not whether the contribution is written in Go or metadata. Native package paths and adapter export names should not become part of abstract declaration identity.

**The laws every contribution must satisfy** (lines 121-145). With `V_T` the
admissible abstract values of a type T, `N_T,p` the admissible native values
under profile p, E for encoding and D for decoding:

```text
D(E(n)) ≈ n    for every n in N_T,p   (the equivalence states whether copying, class identity, branch tags and sharing matter)
E(D(v)) = v    for every v in V_T     (abstract-value equality; any normalization is specified)
```

plus eight further laws, among them capture, generic composition and
cross-target agreement. They matter here for one reason: a family's recorded
decisions have to contain whatever these laws depend on. Two examples from the
first slice: in Go, `Value` is ignored when `Present` is false, so `≈`
identifies two absent values with different `Value`; in TypeScript, an absent
property and a present `undefined` are different native states. Neither fact
is visible in the native type.

**Additivity is the acceptance criterion** (lines 257-265). With shared
machinery frozen, Go bindings may change "only the Go backend and the Go
implementations each semantic owner supplies". "The unit of implementation is
library × target."

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

The advice it comes from (same file, line 44) says why a factory with plain
data cannot be the only path:

> A TypeScript runtime can access a field from its name. A Go family may instead need generated accessors, constructors, typed closures, or helper declarations. An external native type may expose construction only through a particular public API.

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

**Generic requirements** (lines 233-241). "A native type parameter does not
carry the adapter, equality, discriminator or layout evidence needed to bind
values of that type". For `Page<T>`, "`T`'s selected adapter travels with its
native type". For "Nullable `T`", "The null representation cannot collide with
a valid `T`": a requirement that one family places on whatever another family
selected for its operand.

**Absence and null** (lines 164-168), the family behind question 1:

> An optional nullable field has three states. Go uses `Optional[Nullable[T]]` with defined treatment of unused payload fields; a pointer form needs proof that its `nil` cannot collide with a valid value (`**big.Int`). TypeScript uses `field?: T | null`, and the adapter separates an absent own property from a present `undefined` at run time.

**Conversion state** (lines 254-257), the sentence behind question 3:

> Active traversal detects cycles, and completed conversions give sharing. The memo is keyed by the selected representation instance plus input identity, including the descriptor and generic environment. An object-only cache is unsound when one object appears under two mappings.

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
for option C of section 6.1. Section 3.7 shows what is kept and what is lost.

Its verification: "After this slice, a fixture family is added by registering
one contribution. The diff of that test touches no file under
`contracts/generation/go` and no existing mapping."

**Issue #195** is the plan. Its items restate section 2.3 and add two
sentences: "A payload that fails its type is refused before any backend runs"
and "A test backend in this repository, in Go, implements each supported
representation from the plan, the source semantics and the selected
contributions alone. A plan that omits a decision fails the test." It assigns
vertical case 4, the Go case with typed construction, to itself and to no
slice.

**Issue #207** is the first slice. For text, booleans, lists, named closed
records, optional fields and generic parameters with their applications, on
TypeScript and Go (items 1, 3, 5, 9 and 11 of eleven; links reduced to issue
numbers; a number in brackets is a row of research 0015's result table):

> 1. **Contract.** Each family is one representation contract with all seven parts, contributed through the public interface of #196. The kernel, the planner and the completeness test backend gain no case for a family.
> 3. **Conversion.** Go conversion runs in the completeness test backend of #195. The TypeScript factories live in the structural run-time package (bitspark-service#38) and pass this repository's conformance vectors there.
> 5. **Existing output unchanged.** The native declarations generated for today's pagination fixtures stay byte-identical in all four targets. Existing binding identities keep their meaning. The profile's own identity is settled first, in #216.
> 9. **Optional fields.** The TypeScript adapter separates an absent own property from a present `undefined` at run time. Go keeps `BittypeOptional[T]` and ignores `Value` when absent [21].
> 11. **Rust and C++** keep their type support. Their binding support is reported as absent, not implied (#197).

Its verification: "`Page<Person>` goes from horos source to native types, a
plan and a round trip in the Go completeness test backend, with no rule
guessed by the backend." and "A seeded mutant that builds a plan entry in a
separate pass from the native expression turns the completeness test red".
("Binding identities" in item 5 are the rule IDs such as `list/golang/1` of
sections 3.4 and 3.5.)

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

**A caution about our own reading.** On 12 November 2026 we first recorded
what was left of item 6 as "details inside a settled design" (#196, 19:38Z).
A reading of every accepted document against the first slice then found three
shapes that no document gives (#196, 22:40Z). Earlier that day two design
questions on this chain had been answered from the code, wrongly. The reading
behind this document may make the opposite error and read more into the
documents than they say. Please treat every line marked as our derivation
accordingly.

### 2.5 The other side of the interface: the service library's binding contracts

The service library wrote three contracts that the plan and the facets must
agree with: bitspark-service `docs/design/binding-contracts.md`, read at commit
`c36722db901ef7a259242e168c4c35d6f777fbf4`. Its status line: "accepted by
libraries/service on 11 November 2026. Nothing here is implemented yet". Its
lifetime laws (L1 to L14) are partly under revision; bittype has confirmed its
laws on capture (C) and on the call context (K), quoted below, and not the
group L. Standing: **accepted design** (binding contracts of the service
library, 11 November 2026), with L1 marked unconfirmed.

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
an identity, and proved equal in meaning by shared test vectors. No compiler
connects the two.

**What every family declares and every factory receives.** Law C2 (line
316): "Each family declares, for each native value it reads, one of: it
**copies** (the message owns the data), it **retains** a reference, or it
needs the value **stable while reading**." The call context (lines 327-329):
"There is one call context per message per side [...] It is owned by this
repository's adapter runtime and passed to every factory."

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

### 2.6 Names and places are decided after generation

One more accepted decision bears on what a plan can record. bittype has a
stage after generation, called **Files**. By an accepted decision of
7 November 2026, generation decides nothing about files, packages, modules,
imports or module-level visibility. The Files stage of each target takes the
generated declarations and one **placement** supplied by the consumer (which
declaration goes into which file and native scope), writes imports and
visibility, and checks the placed program with the target's own compiler. The
same generated declarations serve every placement. TypeScript has such a
stage today; Go's and Rust's are being built.

Research 0016, a consultation answered on 12 November 2026, extended this to
names. Its result table (`research-docs/0016-native-files-placement.md`, rows
11 and 13) records, as recommendation and accepted action:

- Row 13, "Go's final spelling belongs to Files; identity, naming intent,
  realized spelling": "Accepted. It follows from the accepted decision: Files
  realizes visibility under the language's rules, and in Go that is the case
  of the name."
- Row 11, "The module-pin prefix leaves the default naming": "The
  qualification becomes an explicit, selectable naming policy first, with
  today's as one option and output unchanged."

The decisions recorded on the work item for it (issue #274, comments of
12 November 2026, 22:34Z and 23:05Z) say what generation records instead of a
name: "An intent contains the mapping's chosen name and its qualification,
kept apart", "A qualification is recorded as the module's identity. The
spelled prefix is one policy's realization of it", and "A consumer learns the
realized name and the placed site from a reader on the formatter's result,
keyed by the output's handles." A **handle** is an identity that the sealed
generation output is to issue for each declaration it contains (issue #272,
not implemented; today a consumer finds a declaration by walking the syntax
tree and reading names). A member's name, such as a record field's, stays
with the declaration and is fixed by generation.

So three things that are one string today become three: the **identity** of a
native declaration, its **naming intent** (a chosen name plus a
qualification), and its **realized name and place** (the spelling, and the
file, package or module, which exist only after Files has run for one
placement). Today the declaration recipe of section 3.5 spells the final name
itself: it capitalises the horos name and, for a declaration outside the
entry module, prefixes `M`, sixteen hexadecimal digits of the module pin and
`_`.

Standing: **accepted design** (the decision of 7 November 2026, and rows 11,
13 and 14 of research 0016's result table, 12 November 2026). The design page
for names is written and not yet merged, and nothing of it is implemented.

This does not fit two accepted sentences of sections 2.2 and 2.3 as they are
written. The plan carries "native symbols, distinguishing run-time imports
from type-only imports", the envelope has "Native symbols, imports, factory
references, runtime requirements", and a backend works "from the plan, the
source semantics and the selected contributions alone". A plan issued by the
planner, before any placement, cannot hold a realized name or an import path.
Section 6.5 gives the alternatives.

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

`contracts/generation/go/mapping.go:5-11`, `expression.go:9-14, 30-35, 41-44,
50-74` (accessors omitted). In the first comment, "profile" is a property of a
language identity; it is not the representation profile of section 1.2:

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

// HelperKey has the same two fields.

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
```

Five constructors build expressions: `Native(form, data, args...)`,
`Mapped(slot)`, `Parameter(index)`, `Argument(index)` and `Helper(key,
args...)`.

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

A declaration recipe's callback receives a `DeclarationRequest`, which is an
`Expansion` with one more accessor for the kernel's view of the declaration.
It returns a `DeclarationDescription`: the native name, a declaration form,
the correspondence between the horos generic parameters and the native ones
(the "generic convention"), and a body expression. A `NativeDeclaration[T]`,
as a helper supplies it, is a name, a declaration form, native parameter names
and a body expression.

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

`transform/syntax/shared/generation/go/registry.go:19-30`, and the signature
of the registration function (`:111`):

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

func RegisterMapping[F, T, R any](r *Registry, from, to representation.Language, p g.MappingLibrary[F, T, R]) (MappingRef, error)
```

`RegisterMapping` calls `Describe` once, passes the result to
`FreezeDescription`, increments a counter and stores the description with its
profile under the handle `MappingRef{r, counter}`. `FreezeDescription` copies
the lists and checks them: a non-empty profile, rule IDs unique within the one
description, valid selectors of the source language, slots declared before
use, and size limits. It checks nothing across two registered units.

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
name and files it under its selector's key
(`contracts/generation/go/planning.go:142-144`; recipes are filed alike, and a
second helper under one key fails as ambiguous):

```go
		parts = append(parts, g.SelectedDescription[F, T, R]{ID: fmt.Sprintf("library-%d", e.ref.id), Description: d})
```

```go
		for _, b := range part.Description.Bindings {
			p.bindings[b.Input.Selector.Key()] = append(p.bindings[b.Input.Selector.Key()], selectedBinding[F, T, R]{part.ID + "/" + b.ID, b})
		}
```

So a rule is known to the planner as, for example, `library-3/list/golang/1`.
The first part depends on the order in which the consumer registered its
units. The second part is a string its author chose. The page that documents
the implemented generator says of it
(`docs/type-generation-plugins.md:203-205`): "A binding's ID is local to its
mapping registration. It is not a semantic proof, a globally unique
implementation revision or a configuration digest. There is no PluginID API."

Per use, the planner takes the use rules filed under the use's selector key
(`planning.go:465-475`). More than one is the error "ambiguous generation
mapping", never a preference. Exactly one is applied. With none, a use of a
builtin form fails with "no mapping for" its key, and a use of a named
declaration is planned as a reference to the native declaration that a
declaration recipe produces. So a use rule for a named declaration is an
exact mapping and wins, and a recipe with a zero selector is the default for
every declaration.

A consumer resolves an ambiguity with `Registry.Replace(Replacement)`. A
replacement names the selector and two use rules, each by a registry handle
and its local ID. The handle is a pointer to the registry plus a counter, so
it exists only in the process and registry that issued it. The registry
removes the `Expected` rule before planning. Only use rules can be replaced.
A declaration recipe cannot, and two units that supply a helper under the same
key fail as ambiguous with no way to choose.

### 3.5 One Go family of each kind

The first two are the whole rule package of one family for Go, minus package
clause, imports and comments. The third is the recipe that Go's declarations
unit wraps. The record family follows in section 3.6.

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
`transform/syntax/horos-to-golang/mapping/list/go/list.go:18-27`:

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

`patterns.Element(from, slot)` returns the pattern for the builtin form
`list-of`, declares the slot and supplies the `Capture` function, which asks
the source integration's reader for the list's element use and binds it to
the slot. The template says "a Go slice of whatever is selected for `E`".
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
		// decl the source integration's reading of it; name is the native
		// name computed from decl.Name (error handling omitted)
		body, descriptor, err := src.Descriptor(open.Body)
		// ...
		form := alias
		if descriptor && body.Constructor == "record" {
			form = record
		}
		desc := g.DeclarationDescription[T]{Name: name, Form: form, Origin: decl.Origin}
		for i, f := range open.Formals {
			// ... parameter is the reader's view of f, pname its native name
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
  around the field's type. A representation of another shape (a pointer, a
  field pair, a modifier) needs an edit of the record rule.
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
target builder, which adds the `#include` when it meets them
(`transform/syntax/cpp/target/go/provider.go:191-192`):

```go
		if data[0] == "std" && data[1] == "optional" {
			if err := b.Include("optional"); err != nil {
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
33-43`):

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
   its selector, so the default declaration recipe runs: name `People`, form
   `alias`, a slot for the body. The body is a use of `Page` applied to
   `Person`, so the planner plans the declaration `Page` with the same recipe:
   name `Page`, form `defined`, one parameter `T`, a slot for the body.
3. `Page`'s body is a use of the builtin form "closed record". The record
   rule's callback reads two fields, asks for a slot per field type, wraps the
   second in the helper reference and returns
   `record(field(Items, false, slot), field(Next, true, helper(key, slot)))`.
4. The planner fills the first slot by planning the use `list<T>`: the list
   rule's template `slice(slot E)`, whose slot is the use `T`, replaced from
   the generic environment. It fills the second slot by planning the use
   `text`: the constant `primitive(string)`. It resolves the helper reference
   by planning the declaration `BittypeOptional` once.
5. Back at `People`: the actual argument `Person` is planned like `Page`. The
   body of `People` becomes a reference to the declaration `Page` with one
   argument, the reference to `Person`.
6. The result is a `Program` of four native declarations (`People`, `Page`,
   `BittypeOptional`, `Person`) and a flat list of step strings.

Eight selections were made (the default recipe three times, the record rule
twice, the list rule once, the text rule twice) and one helper was resolved.
None of them exists as a value afterwards.

**Under the accepted design** the same walk has to leave a plan. From
sections 2.2 and 2.3, for this request it holds at least: for each use the
selected contribution and profile, by an identity that a backend in another
process recognises; for the record of `Page` the field layout and, for `next`,
the presence behaviour, as a payload typed in horos; links, visible to the
kernel, from the record's node to the nodes selected for `list<T>` and `text`,
and from the list's node to whatever is selected for `T`; for the application
`Page<Person>` the scoped binder `T`, the native type argument and the adapter
argument; the five version axes, the native symbols and the factory
references. The Go native type of `Page`, byte for byte as today, must be a
projection of the same selections. The completeness test backend then converts
a Go value of type `Page[Person]` to an ontos value and back, using the plan,
the source semantics and the registered contributions, "without guessing a
rule".

**What has no place to come from today.** The trace shows the three open
points concretely.

- At step 3 the record rule decides the optional representation itself. Under
  the design it must ask. There is no use to ask about, only a flag
  (question 1).
- After step 6 the Go test backend must construct a value of the named struct
  type `Page[Person]`. A factory that receives field names as data can do that
  only through reflection. The design's alternative is a lowering facet that
  yields typed Go code, for which bittype has no vocabulary (question 2).
- `library-3/list/golang/1` is the only name the selected list rule has. Its
  first part is a registration counter (question 3).

## 5. Fixed points and their standing

The four standings are those of the Question section. "BG" is binding
generation and "RF" is representation families, both accepted on 9 November
2026 (sections 2.2 and 2.3). "BC" is the service library's binding contracts,
accepted on 11 November 2026 (section 2.5).

| # | Statement | Standing | Shown in |
| --- | --- | --- | --- |
| 1 | The rules AC-02, 03, 04, 05, 08, 09 and 10. In particular: adding a contribution changes only its owner and the consumer composition; openness must not come from arbitrary payloads or unchecked claims; a first change is final in its inputs and identities | Owner rule | 2.1 |
| 2 | "Properly done", with no refactor later | Owner rule | Question |
| 3 | A contribution that needs whatever the consumer registered for a meaning states a demand that the shared mechanism resolves. The record rules deciding the optional representation are a confirmed case of the fault | Owner rule, applied in bittype's package-boundaries record | 2.1, 3.6 |
| 4 | bittype selects native representations and publishes each selection in a typed plan. One backend per target, owned outside bittype, consumes the plan, the source semantics and the contributions. bittype emits no run-time adapters | Accepted design (BG) | 2.2 |
| 5 | The planner owns selection, substitutions, references and the plan. Each representation implementation owns native layout, presence rules, construction and access | Accepted design (BG) | 2.2 |
| 6 | A library publishes a contribution "for one declaration and one supported target representation": selector, target and profile, native facet, value adapter facet, requirements, optional product facets. The facets a binding product needs are selected together with the native facet. A missing facet fails explicitly, with no structural fallback. Composition rejects duplicates | Accepted design (BG) | 2.2 |
| 7 | A target's structural profile and a library's exact mapping use the same public contribution mechanism | Owner rule (AC-04) and accepted design (BG) | 2.1, 2.2 |
| 8 | A family contributes, per target, one contract with seven parts, and is delivered complete | Accepted design (RF) | 2.3 |
| 9 | One immutable representation instance per use. The native type and the plan entry are projections of it. No later pass reconstructs a decision from the source | Accepted design (RF) | 2.3 |
| 10 | A plan entry names a run-time factory and gives it the payload. A family may supply a target-owned lowering facet that yields typed target expressions or helpers. The backend never switches over families. Source-code strings never stand in for typed target expressions | Accepted design (RF) | 2.3 |
| 11 | The plan is a shared envelope plus family-owned payloads. It is a graph: links are explicit in the envelope, and nothing the orchestration must resolve hides in a payload | Accepted design (RF) | 2.3 |
| 12 | A payload is an ontos value whose type its family declares in horos. bittype generates the Go and TypeScript types of that type, and horos conformance validates a payload | Accepted design (RF) | 2.3 |
| 13 | Five version axes. Generic requirements from the first version. The envelope keeps an explicitly unstable version until the six vertical cases pass | Accepted design (RF) | 2.3 |
| 14 | bittype's native vocabularies stay types-only. A backend's emitter is scoped separately, with its own syntax and compiler checks | Accepted design (BG) | 2.2 |
| 15 | The TypeScript backend belongs to the service library and reads a serialized plan. A Go backend belongs to a future Go runtime library | Accepted design (BG) | 1.4 |
| 16 | The existing mapping contract changes once, for every existing rule of the four targets, before the families. After it a family adds a contribution and changes no existing contract | Accepted design (RF, sentence of 11 November 2026) | 2.3 |
| 17 | The kernel's contract package and the registry import neither horos nor ontos: "Generic carriers, contracts, engine, printer and native syntax stages have no external Horos/Ontos imports" (`CHARTER.md:129-131`). A dependency test enforces it | Accepted design (bittype's charter and repository layout) | here |
| 18 | The call context is passed to every factory. Each family declares what it copies, retains or needs stable. The memo key is "the selected representation instance (its descriptor and its generic environment) plus the input's identity" | Accepted design (BC). Its law L1, which also uses "instance", is not confirmed | 2.5 |
| 19 | The completeness test backend is in bittype, in Go, and is a test. bittype publishes no conversion code | Accepted design (BG, bittype's charter, and the result table of research 0015: "The Go case runs in bittype's own completeness test backend") | 2.2, 2.4 |
| 20 | Nothing is released. No interface shown here has a consumer outside our repositories | Implemented today | Question |
| 21 | A registered unit is a profile and three lists; a rule returns native intent; the planner keeps a native program and debug strings; rule IDs are local strings and registry handles are pointers; a replacement names use rules only; one kind of helper exists, wired by a key | Implemented today | 3 |
| 22 | The contribution states its own profile, and a request selects the contributions that state the profile it names | Accepted design (BG's element "Representation context", meaning "Target and supported profile") and implemented today | 2.2, 3.3 |
| 23 | The unit of registration is one contribution, and one contribution is one declaration or builtin form, one target and one profile. BG says "one declaration and one supported target representation", RF says "per target"; neither names builtin forms or says "one registration" | Our derivation | 6.1 |
| 24 | A rule returns one decision record and the projection that computes the native type from it. The planner issues the instance, links its slots and reads no meaning of the record | Our derivation from 5, 9 and 11 | 6.1 |
| 25 | The order: the contract change first with byte-identical output, then the optional representation moves, then additive work | Our derivation from AC-10 and 16 | 2.4 |
| 26 | The Go conversions of the first slice are test code, written as registered handlers that execute the payload | Our derivation | 6.3 |
| 27 | The payload types of the first slice can use only the six existing families, because bittype generates their Go and TypeScript types and represents nothing else yet. A payload cannot hold an integer until the integer family exists | Our derivation from 12 | 6.6 |
| 28 | Generation records a naming intent per declaration: a chosen name and, apart from it, a qualification recorded as the module's identity. Each target's Files stage realizes the spelling (in Go including the case of the name) and decides file, package or module, per placement. At generation time neither the final name nor its import path exists | Accepted design (decision of 7 November 2026; research 0016's result table, rows 11, 13 and 14, 12 November 2026) | 2.6 |
| 29 | Statements 4 and 11 ("native symbols", imports, and a backend that needs only the plan, the source semantics and the contributions) do not fit statement 28 as written: a plan entry cannot hold a realized name or an import path | Our reading of two accepted designs against each other | 2.6, 6.5 |
| 30 | A plan entry refers to a native declaration by its identity, and a backend reads the realized name and place from the Files formatter's result, which becomes its fourth input | Our derivation | 6.5 |

Statements 12 and 17 together have a consequence that is easy to miss. The
plan's node type lives with the planner, in the contract package, because
only the planner may issue a plan. A payload is an ontos value validated
against a horos type. The contract package may import neither. So the planner
holds a payload as something opaque to it, and encoding and validation happen
in a package outside the kernel that may import both.

Out of scope: the design of the Files stage itself (research 0016), beyond
its consequence for the plan in statements 28 to 30; capability lifetime and
the laws L1 to L14; the report of binding support; the choice of native
representation for any particular family.
## 6. What we have considered

Everything in this section is ours and open, including 6.1, which we treat as
settled only in the sense that we would build it tomorrow. "Our leaning" marks
what we would do today. No option is ruled out because the present code would
have to change.

### 6.1 The part we treat as settled: what a rule returns

Fixed point 9 gives the property: one instance per use, two projections, no
later pass. Fixed point 5 names two owners: the representation owns the
decisions, the planner owns selection, links and the plan. Issue #196 item 6
"prescribes no signature". Three shapes satisfy the words to different
degrees.

| Option | What a rule returns for one use | What it guarantees, and to whom | Weakness |
| --- | --- | --- | --- |
| A. Two results side by side | Native intent and a payload, from one call | To a test that looks for a second pass: both come from one call | Each decision is written twice inside the rule, once per result. "Cannot disagree" then rests on each author. A rule that sets the flag in one and forgets it in the other compiles and plans |
| **B. One decision record and its projections** | The family's decisions for this use as one value, with the slots it needs declared. The contribution also supplies a function from (decision record, native intent of the filled slots) to native intent. The plan node is the record plus the planner's links | To a backend: nothing was decided off the record, because the native type is computed from it. To the family: each decision is written once. To the planner: links are explicit and it reads no meaning | The native rule must be a function of the record. Every decision that is control flow today becomes a field. The projection gets no source reader, so anything it needs from the source must have been recorded |
| C. Native intent is the record | Native intent only, as today. The planner keeps it per use, and a plan node is derived from it | One value | A node would consist of target forms and strings that only a native model may read (section 3.7). A decision that no type shows has no place: "`Value` is ignored when absent", "an absent property is not `undefined`", text validation. Fixed point 12 asks for a family-owned horos schema |

**Our leaning: B.** Three consequences follow, all ours:

- A declaration recipe returns a decision record too (record type or alias;
  the name). Its generic part, the correspondence of parameters, is link
  material and belongs to the envelope, which fixed point 13 already requires
  to carry "scoped binder identities, native type arguments, adapter
  arguments".
- The record is opaque to the planner and typed by its family. Because of
  fixed point 17 its encoding as an ontos value and its validation happen
  outside the kernel.
- The planner issues the instance: it gives it identity, fills and links its
  slots and makes it immutable. A contribution cannot construct an instance or
  a plan node. This continues an existing rule, that only the planner issues a
  plan.

The unit follows from fixed points 6 and 8: one contribution is one
declaration or builtin form, one target, one profile, with all its facets, and
it is what a consumer registers. The three lists of section 3.3 go away.

**Why the kernel contract may change at all** (AC-08). The missing capability:
the planner cannot keep "one origin per representation decision" while a rule
returns native intent only and the per-use result is discarded. Why this
layer: only the planner sees each selection exactly once and resolves slots,
substitution and scope. How it serves independent contributions: the carrier
is the same for a structural family and for a library's exact mapping, and
the planner learns no payload meaning.

**What speaks against our leaning.** The strongest argument we see is that B
buys "cannot disagree" for the native type only. The conversion is not a
projection that bittype computes from the record. It is a factory in another
repository (TypeScript) or a handler in a test (Go) that reads the record and
must agree with a native type that a function in bittype derived from the
same record. The record makes disagreement detectable, by conformance vectors,
and not impossible. If that is all it buys, option A with a check might buy
the same at less cost to authors. A second argument: B requires the native
type to be a function of recorded decisions and slot results. A family whose
decision depends on a property of what was selected for its operand (a
nullable value that is `T | null` when `T` excludes null and a wrapper
otherwise) needs facts about the operand's selection as inputs, and today a
rule never sees its child's selection (section 3.3). B does not say where
those facts come from. That leads into the first open point.

### 6.2 Open point 1: the kernel's demand capability

**What is decided.** The record rule must not decide the optional
representation. It states a demand, and the kernel resolves it to whatever
the consumer registered (fixed point 3). "Optional fields" is one of the six
families of the first slice, so it has all seven parts on TypeScript and Go.

**What no document says.** (a) What the demand is keyed by: the kernel selects
by the selector of a use, and an optional field is no use (section 3.6).
(b) What the answer is: a type around the field's type in Go, Rust and C++, a
property of the field member in TypeScript. (c) How the demand appears in the
plan. (d) How a consumer replaces the answer.

| Option | The demand | What it guarantees, and to whom | Weakness |
| --- | --- | --- | --- |
| A. Presence stays in the record family | None. The record's payload holds the presence representation; Go's helper becomes a private part of Go's record contribution | One contribution; the simplest plan | Excluded by fixed point 3. A consumer who wants another optional representation still replaces the record contribution |
| **B. A source-owned selector for a form that is not a use** | The source integration publishes a selector for "an optional field of a record", as it does for `list-of`. The record contribution asks the kernel for the representation under that selector and passes the field's type use as the operand. Any contribution registers for that selector like for any other | To the consumer: selection by exact identity, the ambiguity error, replacement and the missing-contribution failure all apply unchanged. To the backend: the record's node links to an optional node by a typed slot. To the record family: it never learns what was selected | The kernel gains a capability: selection for something that is not an occurrence in the source graph, at the request of a contribution. Its general justification (AC-08) must be written. The set of such "modalities" is open-ended, and each needs a selector that the source integration owns |
| C. A key that the composition wires | Today's helper key, made complete: the unit that supplies the key becomes a full contribution with payload and conversion | The kernel resolves it and a missing supplier fails | Selection is by wiring, not by identity. Replacement does not apply. The key is a pointer of one process and cannot be recorded in a plan. Two suppliers are an error with no way to choose |
| D. A demand keyed by a role that the asking family owns | The record family publishes a role identity ("presence of a field"), and contributions register as providers of that role for a target | The asking family defines exactly what it needs from an answer, including its shape | A second selection mechanism beside selectors. A provider depends on the vocabulary of the family that asks, so "optional" becomes the record family's concept and not horos's. Two families that need the same modality (records and tuples with optional elements) would each own a role |

**What the answer is** is a second, independent choice:

- **T. A type.** The answer is native intent with one slot for the operand:
  `BittypeOptional[·]`, `std::optional<·>`. Then TypeScript has no optional
  contribution, the `?` stays the record family's own layout decision, and
  the absence rule of TypeScript lives in the record family. The six families
  mean different sets of contributions per target.
- **M. A field member.** The answer is at the level of the member: given the
  field's name and the operand's representation, it yields the member (in
  TypeScript a property with a modifier, in Go a field with a wrapped type and
  the helper declaration it needs). All four targets ask alike, and each
  target's presence rule has one owner. Then the demand's result is not "a
  type", and the kernel's notion of what a slot can hold widens from "native
  intent for a type position" to "native intent for a position that the
  asking family defines".
- **R. A decision record only.** The answer is the optional family's decision
  record, with no native intent of its own. The record family's projection
  reads a small, published vocabulary of that record ("wrap in this type", or
  "mark the member") and builds the member itself. Then the record family
  interprets another family's payload, which fixed point 11 seems to exclude
  for the orchestration but does not clearly exclude between families.

**Our leaning: B with M.** "Optional fields" is then one contribution per
target, selected for the horos optional-field form. The record's payload keeps
the field layout and, per optional field, a link to the optional node. Go's
`BittypeOptional` declaration becomes a part of Go's optional contribution, so
the consumer-made key and the separate helper registration disappear. An
optional nullable field composes without touching the record family: the
nullable family represents the field's type, the optional family the field.

**What speaks against our leaning.** First, M makes the answer's shape depend
on the asker: a record wants a member, a tuple with an optional trailing
element wants something else, a function parameter something else again. One
selector for "optional" then has to serve several position kinds, or there is
one selector per position kind, and the optional family is no longer one
thing. Under T this problem does not arise, at the price of TypeScript's
asymmetry. Second, B invents source identities for things horos does not
treat as types. The next modality (a map's key requirement, a union's
discriminator) may not be a form of the source at all but a relation between
two selected representations, and then a source-owned selector is the wrong
kind of key. Third, a demand that returns a link and no facts cannot serve a
family whose decision depends on its operand's selection (end of 6.1). If
demands must also return declared facts about the selected contribution, that
is a larger capability than the one we lean to.

### 6.3 Open point 2: facets as an open set, and the lowering facet

**What is decided.** A contribution pairs a native facet with a value adapter
facet; a root facet and "optional product facets" exist; a lowering facet may
turn recorded decisions into "typed target expressions or helpers"; the
backend never switches over families; bittype's vocabularies stay types-only
(fixed points 6, 10 and 14).

**What no document says.** Whether the set of facet kinds is closed in the
kernel's contract. Who owns the vocabulary in which a lowering facet for Go
writes an accessor or a constructor, given that the Go backend's owner does
not exist yet. Whether the first slice must already pass vertical case 4, "a
Go case with typed construction and a generic adapter argument".

Two choices are entangled. The first is the shape of the contract:

| Option | The contribution contract | What it guarantees, and to whom | Weakness |
| --- | --- | --- | --- |
| F1. A closed list | The kernel's contribution type has one field per facet kind: native, value, root, operation, lowering | Every reader sees every facet's type | A new binding product (presentation, documentation, a second backend) needs a kernel edit: the violation AC-03 describes. The kernel would name `operation`, a concept of the service library |
| **F2. An open set, keyed by the reading product** | The kernel knows the native facet, because it plans with it, and the decision record. Every other facet is carried under a kind that the product reading it issues, as a native model issues its forms today. The kernel checks presence ("a bindings request needs kind K; this contribution has none") and never reads a facet's content | A new product adds a kind and edits nothing shared. The failure "missing facet" stays the kernel's | The kernel carries values it cannot check. What keeps this from being the "arbitrary payloads" of AC-05 is only that each kind has an issuer that validates its facets, and that a facet may refer to slots only through kernel-issued references |
| F3. No facets in the contribution | The contribution is the native facet and the decision record. Everything else is registered separately, with the product that reads it, under the contribution's identity | The kernel contract is smallest. A second owner can add a facet to someone else's family | "Selected together" (fixed point 6) becomes a lookup by identity. We considered and dropped this on 12 November 2026 because joint ownership excludes a second owner adding a facet; it is, however, exactly how the TypeScript factory already relates to bittype's family (section 2.5) |

The second is the first slice's Go conversion, which decides whether the
contract is tested against lowering before it is fixed:

| Option | How the Go test backend converts `Page[Person]` | What it proves | What it leaves open |
| --- | --- | --- | --- |
| **V1. Registered handlers that execute the payload** | One Go conversion per family, registered with the family's contribution, takes the payload and the conversions of its slots. For a named struct it reads and builds fields by the names in the payload, through reflection | Completeness: every name, order and presence rule comes from the payload. The design allows it: "its handlers are registered contributions executing selected decisions" | The lowering facet is not exercised. Case 4 stays open. Nothing shows that typed construction can be expressed from the decision record |
| V2. Lowering to generated Go | Each Go family supplies a lowering facet; the test backend assembles Go source, compiles it and runs it | Also the lowering facet, and case 4 | Needs typed expressions for value-level Go: calls, field access, composite literals, closures. bittype may not own that vocabulary (fixed point 14), and its future owner is absent. The first slice would define a second public interface as a side effect |
| V3. A paper lowering | V1 in code, plus the lowering facet of the Go record and list families written out on paper against F2, with a hypothetical expression vocabulary, before the contract is fixed | That F2 and the decision record carry what a lowering needs, to the extent paper shows it | No compiler has checked it |

**Our leaning: F2 with V1, and a record on the issue that case 4 is not
claimed by the first slice.** The lowering facet is "target-owned", and its
vocabulary belongs to a backend's emitter, which is "scoped separately". The
first slice's claim is a complete plan and the conversion laws for six
families, and V1 establishes both. Adding a lowering facet later stays
additive only if the contract does not enumerate facet kinds, which is why F2
cannot be postponed even though lowering can.

**What speaks against our leaning.** The six vertical cases exist because the
expert of research 0015 expected the envelope to be wrong until difficult
compositions had been built. We would be fixing the *contribution contract*,
which "changes once", on the evidence of the easy path only, and reserving
the "unstable" label for the envelope. But a lowering facet is a function of
the decision record and of typed references to native symbols and slots. If
it turns out to need something the record or the slots do not carry (a typed
reference to a generated field instead of its name; the native type of a slot
as a value-level type argument; the helper declaration's constructor), the
contribution contract changes a second time, not the envelope. Reflection in
V1 hides exactly these needs, because it turns every reference into a string
at run time. A second argument: under F2 the kernel checks the presence of a
facet kind and nothing else, so "selected together" guarantees that the
facets were registered together, not that they agree. A third: the process
boundary forces lookup by identity for TypeScript anyway, so F3, one
mechanism for all facets, may be more honest than "together" in Go and "by
identity" elsewhere.

### 6.4 Open point 3: the identities of the plan

**What is decided.** The envelope records the identities of the selected
contribution and profile, and source-use identities. Duplicate contributions
are rejected. There is one instance per use. The memo key is "the selected
representation instance (its descriptor and its generic environment) plus the
input's identity", and one object under two representations converts once per
representation (fixed points 6, 9, 11 and 18).

**What no document says.** Three identities are in play, and no document gives
any of them a form (a fourth, the identity of a native declaration, is the
subject of 6.5):

1. **The contribution identity.** What a plan node records as "selected
   contribution", what a backend in another process uses to find the factory,
   and what the service root's facet compares to recognise "the structural
   record family". Today's only candidates are a registration counter and a
   local string (section 3.4).
2. **The plan node.** One per use, or one per distinct selection.
3. **The run-time identity.** The part of a memo key that is not the input.

For the contribution identity:

| Option | Form | What it guarantees, and to whom | Weakness |
| --- | --- | --- | --- |
| **I1. Owner-stated name and version** | A string the contribution states, such as `list/golang/1`, unique within a registry; the registry rejects a second registration of the same string | To a backend: a stable key in every process. Today's rule IDs already have this form | Uniqueness holds per registry only. Two owners can publish the same name, and a backend that pins it cannot tell them apart. The version in the name is one of the five axes (the semantic identity) folded into an identity |
| I2. Owner namespace, name and version as separate parts | The owner's namespace is an identity the owner demonstrably controls: for an exact mapping the horos module pin of the declaration it maps; for a profile's builtin families the profile's identity | Collisions across owners cannot happen by accident, and the version axis stays an axis | Needs a rule for what a namespace is for a contribution that maps someone else's declaration (a Go mapping of a library's `Channel<T>`), and for two contributions of one owner for one declaration (two profiles, two targets) |
| I3. Derived | A digest over the contribution's payload schema, its target, its profile and its selector | Nobody states it, so nobody can state it wrongly | A factory can change meaning without changing any of the hashed parts, which is why the design has a semantic-identity axis at all |

For the node and the word "instance":

| Option | The plan node | Weakness |
| --- | --- | --- |
| N1. One node per use | Two fields of type `Page<Person>` are two nodes | Two memo identities for one representation, which K3 and issue #195's criterion exclude, unless the run-time identity is a separate field that both nodes share |
| **N2. One node per distinct selection; uses refer to nodes** | A node is identified by (contribution identity, profile identity, source identity of the type, generic environment). Each use holds one reference to a node. "Instance" means the node | Per-use facts (provenance, the source-use identity the envelope lists) need a second, smaller record per use. What "source identity of the type" is for something that is not a type, such as the optional-field form of 6.2, is undefined |
| N3. Two levels, both named | Use entries (one per use, with provenance and a node reference) and representation nodes (as N2). The design's "instance" is the use entry; the run-time identity is the node | Two concepts where the accepted text has one word |

**Our leaning: I1 for the contribution identity, N2 for the node.** bittype's
own families keep the form `list/golang/1`. The registry handle of section
3.4 stays what authenticates a replacement inside one process; the stated
identity is what the plan records.

**What speaks against our leaning.** On I1: the whole point of the identity is
the process boundary and the second repository, and "unique within a
registry" is a guarantee about the one place that does not need it. The
moment the service library and Content both register contributions, nothing
but convention keeps their names apart, and a wrong factory selected by a
colliding name fails at run time, in someone else's CI at best. On N2: we
noticed while writing this document that "generic environment" hides a
choice our derivation did not make. A plan for the open declaration `Page<T>`
naturally has one record node for its body, with a link to the binder `T`,
and each application supplies an adapter argument. The run-time identity of
"the list inside `Page<Person>`" is then (that node, the environment in which
`T` is `Person`'s node), a pair that exists only at run time. Making it a plan
node means one node per instantiation: the plan is monomorphised, its size
depends on the applications in the program, and a generic declaration alone
has no complete plan. If nodes are open, then N2's node is not the run-time
identity either, and there are four identities, not three. We do not know
which of the two the design means by "its descriptor and its generic
environment".

### 6.5 Native symbols in a plan whose names are realized later

**What is decided.** The plan carries "native symbols" and imports, and a
backend needs nothing beyond the plan, the source semantics and the
contributions (fixed points 4 and 11). Generation records a naming intent;
each target's Files stage realizes the spelling and decides file, package or
module, per placement (fixed point 28).

**What no document says.** How a plan entry refers to a generated native
declaration, and from where a backend learns the name and import path it must
write into its glue. This is a fourth identity beside the three of 6.4: the
identity of a native declaration.

| Option | What a plan entry holds | What it guarantees, and to whom | Weakness |
| --- | --- | --- | --- |
| S1. The plan is produced after placement and carries realized symbols | Spelled names and import paths | The accepted sentences hold as written: one artifact, three inputs | The plan depends on a placement, and one generated set serves several. The planner cannot issue it alone, so either Files completes a plan it did not make, or the kernel learns placement. A consumer that only wants conversions for already placed types still runs Files first |
| **S2. The plan carries identities; the backend joins it with the placement result** | For each native declaration its identity (the handle the sealed output issues) and its naming intent. The backend reads the realized name and place from the Files formatter's result, by handle | The plan is independent of placement. Each language's spelling rule has one implementation. The planner stays the plan's only issuer | The placed result becomes a fourth input of a backend, which amends an accepted sentence. Two artifacts can be mismatched: a plan of one run with the placement result of another. The handle must survive serialization, and today an output's identity is a pointer and its sites repeat on every run |
| S3. The plan carries intents; the backend derives spellings | Chosen name and qualification | Three inputs, no dependence on Files | The backend repeats each language's rule (Go's case, the qualification policy, collision handling). Two implementations must agree, which is the guessing that the completeness test excludes |

**Our leaning: S2.** A plan entry refers to a native declaration by its
identity, and the backend reads the realized name and place from the Files
formatter's result. Member names, which generation fixes, stay in payloads.

**What speaks against our leaning.** S2 moves the difficulty into the join.
For the TypeScript backend the placement result must cross the same process
boundary as the plan, so it needs a serialized form and a binding to the plan
that a reader can check, and neither exists. The completeness test then has
four inputs, and "the plan is incomplete" gets a second meaning: the plan is
complete and the join failed. A backend also places its own glue (it "owns
imports, names, placement and assembly"), so it may be less a reader of a
placement result than a participant in placement, which S2 does not model.
And for Go and Rust no Files stage exists yet, so the first slice's Go test
backend would read names that generation still spells, through an interface
that is meant to change.

### 6.6 What the first version must contain, and how to test it before the rewrite

The contract change rewrites 21 rule packages (five for each of TypeScript,
Rust and C++, six for Go) and the kernel's planner. We see two ways to find a
fault in the contract before that work, and have done neither.

- **Paper families.** Write the contribution of one family of each kind
  against the proposed contract, for TypeScript and Go, without code: a
  constant (text), one slot (list), a callback with a demand (record with an
  optional field), the declaration recipe with its generic part, an exact
  external generic (`Cell<T>`), and one later family that places a
  requirement on its operand (a nullable value, or a map with its key). For
  each: the decision record as a horos type, the projection, the slots, the
  facets, the plan nodes for one example, and what a factory receives.
- **A materially different fixture.** The architecture rules ask for one
  (section 2.1), and issue #196 requires it after the change. It could be
  written against the contract's types before the 21 packages move.

What we believe the first version of the envelope needs, from fixed points 11
to 13: one node per selection with the contribution and profile identity; the
source reference and source identity; links by identity and never by nesting,
so that a recursive reference is later one more link; for each generic use
the scoped binders, native type arguments, adapter arguments and
requirements; the five version axes as five values; root products for
requested declarations; places for factory references and run-time
requirements. One limit applies (fixed point 27): a first payload schema
cannot hold an integer, such as a field index or a count, nor a union or an
exact value.

We lean to doing the paper families before the rewrite and the fixture with
it. What speaks against: paper finds what its author already suspects, and
the cases most likely to break the contract (recursion, unions with
discriminators, explicit specialization) are the ones we understand least
and would write worst.

### 6.7 Surprises we already suspect

These are known suspects. Question 5 asks for what lies outside this list.

1. A family whose decision depends on what was selected for its operand
   (6.1, 6.2): today a rule never sees its child's selection.
2. Refusals. A rule refuses a use by returning an error from its callback
   (section 3.6: "nested records require a named declaration in this
   profile"). The representation contract lists "explicit refusals" under
   applicability, as declared facts.
3. Target configuration that changes a decision. TypeScript's meaning of
   `next?: string` depends on the compiler option
   `exactOptionalPropertyTypes`. The envelope has "target and its
   configuration"; a decision record that silently assumes an option is
   wrong under another.
4. A native declaration that a contribution owns without a horos declaration
   behind it (Go's `BittypeOptional`): who states its naming intent, how two
   contributions that need the same one share it, and how a plan refers to it.
5. The semantic version inside an identity (`list/golang/1`) against the
   design's separate version axis (6.4).
6. Generic conventions of exact external generics whose native parameters
   are reordered, omitted or added, and what "the adapter argument travels
   with the native type argument" means for an omitted one.
7. One contribution wanted under two profiles, when a profile is a field
   that the contribution states once.

## 7. Questions for the expert

In section 5, owner rules are fixed, and accepted designs are fixed unless
you see a contradiction. Everything else is open, including what the code
does today and the shape of section 6.1, which we invite you to attack.
Please name the option first and then the reason. Patterns from comparable
systems are welcome where they change the answer, including what failed:
compilers with separately registered lowering passes, schema and interface
compilers with plugin protocols across a process boundary, serialization
frameworks that derive a type and its codec from one description, and
type-class or trait resolution where one instance demands another.

1. **A demand for something that is not a source occurrence.** A record
   contribution must ask the kernel how an optional field is represented,
   and the kernel must resolve the request to whatever the consumer
   registered; later families ask alike for a nullable value, a map's key
   requirement or a union's discriminator. (a) What would you key such a
   demand by, and who should own that key: the source integration, the
   asking family, or the composition? (b) What should the answer be: a type
   around the operand, a member of the asking family's layout, or a decision
   record that the asker reads, given that the same modality is a property
   of a field in TypeScript, a generated type in Go and an external type in
   Rust and C++? (c) Should a demand be able to return declared facts about
   what was selected, for a family whose own decision depends on its
   operand's representation? (d) How should a resolved demand appear in the
   plan, and by what act does a consumer replace the answer? We lean to a
   source-owned selector for the optional-field form with the member as the
   answer (6.2, option B with M); 6.2 also lists what speaks against it. Your
   answer decides the kernel capability we add in the one contract change,
   and whether "optional fields" is one contribution on every target.

2. **Facets and lowering.** A contribution pairs its native facet with a
   conversion facet and may carry more: a root facet, an operation facet, a
   lowering facet that yields typed target code. bittype's own vocabularies
   describe types only, and the owner of Go's value-level code does not
   exist yet. (a) Should the kernel's contribution contract enumerate facet
   kinds, carry an open set whose kinds the reading products issue, or carry
   none and leave every facet to a lookup by the contribution's identity?
   (b) If the set is open, what must the kernel still guarantee about a
   facet it cannot read, so that openness is not "arbitrary payloads"?
   (c) Who should own the vocabulary of typed value-level target expressions,
   and what is the least a contribution contract must fix about it today so
   that a lowering facet can be added later without a second change?
   (d) Should the contract be fixed before or after one Go case with typed
   construction and a generic adapter argument has been built? We lean to an
   open set, to handlers that execute the payload in the first slice, and to
   building that Go case later (6.3, F2 with V1); 6.3 also lists what speaks
   against it. Your answer decides the shape of the contribution type and
   whether the first slice includes that case.

3. **The identities of a plan.** A plan must name the selected contribution
   so that a backend in another process, and a factory package in another
   repository, recognise it; it has nodes; a conversion's memo is keyed at
   run time by "the selected representation instance" plus the input; and it
   must refer to generated native declarations.
   (a) Who should issue the identity of a contribution, what should it
   consist of, and how is a collision between two owners detected, and
   where? (b) Should a plan have one node per use, one per distinct
   selection, or both as two named things? (c) For a generic declaration,
   should nodes be open, with the environment supplied at an application,
   or one per instantiation, and what then is the run-time identity that a
   memo is keyed by? (d) Which of these should the word "instance" name?
   (e) A generated native declaration gets its final name and its file,
   package or module only after generation, per placement (2.6). Should the
   plan be produced after placement and carry realized symbols, carry
   identities that a backend joins with the placement result, or carry
   naming intents, and what must a backend in another process be given so
   that it never guesses a name? We lean to an owner-stated name with a
   version, unique within a registry, to one node per distinct selection,
   and to identities joined with the placement result (6.4, I1 with N2; 6.5,
   S2); both sections also list what speaks against the leaning, including
   a choice between open and instantiated nodes that we have not made. Your
   answer decides the identity fields of the contribution contract and of
   the envelope's first version, and whether a placement result becomes an
   input of a backend.

4. **Testing a contract that may change only once.** One change rewrites
   every existing rule of four targets onto the new contract; the plan's
   envelope may stay unstable until six vertical cases pass, the contract
   may not. (a) How would you test the contract of sections 6.1 to 6.5
   before that rewrite, and which families, written on paper or as fixtures,
   are most likely to break it? (b) What must the first version of the
   contribution contract and of the envelope already contain so that a later
   family adds one contribution and one payload schema and nothing else?
   (c) Is the division sound, in which the contract is final in the first
   change while the envelope is declared unstable, or do the two have to
   stabilise together? We lean to paper families for six kinds before the
   rewrite and a materially different fixture with it (6.6); 6.6 also says
   why paper may not be enough. Your answer decides what we build and write
   before the rewrite starts, and what the first version carries.

5. **Open.** What should we know that these questions do not ask? When a
   type generator's contribution contract is extended to carry conversion
   decisions for backends it does not own, what tends to be discovered late
   and then forces a second change of that contract? Section 6.7 lists the
   surprises we already suspect, so what lies outside that list is the most
   valuable to us. We will use your answer to add fields, refusals or checks
   to the contract before the rewrite.
