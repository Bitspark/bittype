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
