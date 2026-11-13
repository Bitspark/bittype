# Research: The contribution contract of a type generator that must also carry conversion decisions: what a rule returns, requests between contributions, facets and the identities of the plan

**ID:** 0017
**Date:** 12 November 2026
**Status:** draft
**Owner:** claude-contract-chain-issues
**Issue:** https://github.com/Bitspark/bittype/issues/196 (item 6), with #195, #197 and #207
**Code baseline:** bittype, code baseline 2daea52 (2daea52e18e20369401500ddeeb19efae581b306); documents at 7e9dc34

## Question

**In one sentence:** what should the contract between a type generator's
shared kernel and its independently owned contributions be, so that one
selection per use of a type yields both the native type and everything a
separately owned conversion backend needs, and so that the contract need not
change again once other repositories build on it?

**horos** is a declarative type language. **bittype** is a Go library that reads
horos declarations and generates native type declarations for four languages:
TypeScript, Go, Rust and C++. It is extensible: every rule that turns one kind
of horos type into a native type is supplied from outside the shared mechanism,
by a package that a consumer registers. We call the shared mechanism the
**kernel** and one registered unit a **contribution**.

Today the unit that generation registers describes native types and nothing
else. Two designs accepted on 9 November 2026 give it a second duty. Services
declared in horos need conversion code between native values and message
values (the values our services exchange, section 1.1). That code is emitted
by a **binding backend**, one per language, which is owned outside bittype
and, for TypeScript, runs in another process. bittype stays the only authority
that selects native representations. It has to publish every selection, with
every decision the selection contains, as a typed **plan** that a backend can
implement "without guessing a rule".

**"Once".** The owner's direction of 9 November 2026 is:

> I want this properly done, so that we don't accumulate technical debt and have to refactor later.

From it and from the owner's rule AC-10 (section 2.1) we derived that the
existing contract between kernel and contribution "changes once": every
existing rule of the four targets moves to the new contract before any new
kind of type is added, and afterwards a new kind of type is a new contribution
that edits nothing shared. That sentence is ours. We wrote it into a work item
and, as a delivery note, onto an accepted design page; it is no part of what
was accepted, and question 4 asks about it.

What "once" binds, as we have decided it: the thing to avoid is a change of
the contract after independently owned contributions and a backend in another
repository are built on it. Nothing is released, and no consumer outside our
own repositories depends on any interface shown here. A prototype on a
branch, including one target moved end to end with one TypeScript factory
round trip, is allowed before the contract is fixed on the main branch.

**The shape we would build** (standing: our derivation, open). It bundles
four claims of different standing:

- The contribution states its own profile (a named promise about the native
  types). Accepted design, and implemented.
- The unit: one contribution is everything one owner publishes for one horos
  declaration or builtin form (a type constructor built into horos, such as
  list or record), one target and one profile. Our derivation; the accepted
  texts say "library × target" and "per target" (section 6.1).
- For one use of a type, a rule returns **one decision record** together with
  the function that computes the native type from it, so that each decision
  is written once. Our derivation from the accepted property "one instance,
  two projections" ("instance" is the immutable result of selecting a rule
  for one use).
- The kernel issues the result (creates it and gives it an identity that no
  contribution can forge), links it to the results of other uses and reads no
  meaning of the decision record. Our derivation.

**What is open besides that shape.**

1. **How one contribution asks the kernel for another's answer.** The case at
   hand: by our own finding, the rule for horos record types should not decide
   how an optional field of the record is represented.
2. **Whether a contribution carries an open set of parts ("facets")** whose
   kinds the kernel does not enumerate, among them a lowering facet (an
   optional part through which a family emits typed target expressions, such
   as a Go constructor, from its recorded decisions).
3. **Which identities the plan records**: that of a contribution, which a
   backend in another process must recognise; that of a plan entry; the one a
   conversion's memo is keyed by at run time (a conversion keeps a
   per-message table of objects already converted, to preserve sharing and
   detect cycles); and that of a generated native declaration, whose file,
   package or module is decided only after generation.

**How to read the fixed points.** Every statement this document treats as
given has exactly one of four standings (section 5):

- **Owner rule.** Stated by the project's owner. Fixed.
- **Accepted design.** We name the document, its date and who accepted it.
  For binding generation and representation families (sections 2.2, 2.3)
  "accepted" means: an expert's advice was evaluated row by row by the owner
  of that research, another seat of our team, and recorded as the project's
  design. For the naming design (section 2.6) it means accepted by the seat
  that writes this document. Fixed unless you see that it contradicts an
  owner rule, another accepted design or itself. If so, please say so.
- **Implemented today.** A fact of the present code. It can change now, and
  rule AC-09 gives it no weight in the design.
- **Our derivation.** Open. Section 6 opens with the list.

The code in section 3 is evidence of the present state and never a reason for
or against an option. One consequence we state once: our work items require
the generated output of today's fixtures to stay byte-identical through the
contract change. That is a regression oracle for the migration. It is no
constraint on the answer.

Quotations of code give file and line at the code baseline. Quotations of
documents give the file. Quotations from GitHub issues give the issue number
and the time of the comment in UTC. All code is Go unless marked.

## 1. System overview

### 1.1 The components

- **ontos** is a data model with a canonical binary encoding. A value is an
  atom (a byte string) or a tuple of values. Every message that our services
  exchange is an ontos value. ontos has a Go and a TypeScript implementation.
- **horos** is a type language over ontos values. A horos type denotes a set of
  ontos values. It has builtin forms (UTF-8 text, booleans, lists, records with
  required and optional fields, and others such as integers, maps and unions)
  and named declarations with generic parameters. Declarations live in modules;
  a module is identified by a hash of its content, its **pin** (quoted
  documents also say "module digest"). For a type without free parameters
  horos produces a **descriptor**, which is itself an ontos value, and its
  operation `diagnose` checks whether a value conforms to a descriptor.
- **bittype** reads horos modules, checks them, and generates native type
  declarations. It is written in Go, whatever language it generates. Its
  charter, the repository's scope document, keeps the shared parts free of the
  source language: "Generic carriers, contracts, engine, printer and native
  syntax stages have no external Horos/Ontos imports" (`CHARTER.md`). In our
  terms: the kernel imports neither horos nor ontos, and a dependency test
  enforces it.
- **The service library** (repository `bitspark-service`) is the TypeScript
  run-time library of our services. It owns the adapter interface (an
  **adapter** converts between a native value and a message value), the
  per-message **call context**, and, by a decision of 9 November 2026, the
  TypeScript binding backend. It will also host the **structural run-time
  package**, in which the TypeScript conversion functions for bittype's own
  families live. Three generic types of the service library appear in this
  document: `Function<A, B, E>` (a callable with an argument, a result and a
  declared failure), `Stream` (another type of the same kind) and
  `Service<Ops>` (a record `Ops` of operations served at one address). A
  fourth, `Cell<T>` (a holder of a current value), belongs to an independent
  fixture library and is our test case for an exact mapping. A value of one
  of them inside a message is a live reference to something remote, a
  **capability**. Their native types are types the library ships, not
  generated ones.
- **Consumers** are builds of our own services. The pilot is a service called
  Content.

An example that the rest of the document reuses:

```horos
module example.pagination;
type Page<T> = { items: list<T>; next?: text; };
type Person = { name: text; };
type People = Page<Person>;
```

`next?` marks an **optional field**. That term has one definition here: in
horos it is a flag on a field of a record form, and no type; in the accepted
design it is one of the representation families; whether it is a contribution
of its own is question 1. bittype generates today, for TypeScript and for Go
(whitespace ours):

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
for C++ `::std::optional<::std::string>`. So the same source fact becomes a
property of the field in TypeScript, a generated helper type in Go and an
external library type in Rust and C++.

### 1.2 Terms used throughout

Each term has one sense in this document. Where the code or a quoted document
uses a word differently, the row says so.

| Term | Meaning here |
| --- | --- |
| Target | One of the four native languages. The Go interface named `Target` in the code is the language's syntax builder; we call it the **target builder** |
| Source integration | The package that reads horos for the kernel: it classifies declarations and uses, and gives rules a typed reader for horos-specific content |
| Use | One occurrence of a type expression in the source, reached under one declaration along one path. In `Page`, `list<T>` is one use and `T` inside it another. The code says "occurrence" |
| Source identities | Three are distinct. The **selector**: what kind of thing a use is, equal for every `list<...>`; for a builtin form `builtin/7:list-of/0:` (each part length-prefixed: the constructor `list-of`, an empty discriminator), for a named declaration `declaration/<module pin>/<name>`. The **source-use identity**: which occurrence. The **source identity of a type**: which type; for a type without free parameters its descriptor. For an open type no accepted document says what it is, and section 6.4 depends on that |
| Builtin form, native form | A **builtin form** is a type constructor built into horos (text, list, record). A **native form** is a constructor of native intent, issued by one target's native model, such as Go's `slice`, `record` and `field` |
| Native intent | bittype's description of a native type before any syntax exists: a tree of native forms with string data and arguments. In the code: `Expression[T]`. It describes types only. A **native model** is the package of one target that issues its native forms and checks native intent |
| Kernel | bittype's shared generation mechanism: the contract package with the planner, and the registry. Quoted documents and code comments also say "shared mechanism", "orchestration" and "composer" |
| Planner | The part of the kernel that walks the uses reachable from the requested declarations, selects a rule for each and assembles the result |
| Rule | A function or template that maps the uses of one selector to native intent. The code and quoted documents call a rule for uses a "binding" (`BindingDescription`, "a binding's ID"). In our own words "binding" always means value conversion, so we say **use rule**, and **declaration recipe** for the rule that turns a named horos declaration into a native declaration |
| Helper | A native declaration that is not generated from a horos declaration and that a rule refers to through a key. Go's `BittypeOptional` is the only one. The accepted design also says "helpers" for generated value-level functions; we call those lowering output |
| Family | A representation family: one kind of horos type together with its native representation (text, lists, closed records, optional fields, integers, unions). A **structural family** is one contributed under the structural profile |
| Profile | A named promise about what the native types retain. Its identity is its name string. The only implemented one is `structural-representation/1`, the **structural profile**: native types that mirror the horos structure and represent every value of the horos type. A second, `structural-enforcement/1`, is declared and stated by no rule |
| Exact mapping | A library's own representation of one named horos declaration, for example of its `Cell<T>` as a class it ships. Selected in place of the structural families for that declaration |
| Library | In our words and in quoted designs: a separately owned code base that owns horos declarations and their native types (the service library, Content). In the code (`MappingLibrary`, `Library()`, the name `library-3`) it is one registered unit |
| Contribution, registered unit | A **contribution** is one thing a consumer registers on its own (the owner's definition, section 2.1). What one contribution comprises is open (6.1). A **registered unit** is what exists today: a "mapping library" holding one rule; quoted documents call it a "mapping" |
| Contract | Without qualification: the contribution contract, between kernel and contribution. Other uses are always qualified: the "representation contract" (what a family must define, section 2.3), the service library's "binding contracts" (2.5), the Go "contract package", and quoted "contracts" that mean service types such as `Function` |
| Facet | A part of a contribution that one product reads: a native facet (the type), a value adapter facet (the conversion; quotes also say "adapter facet"), a root facet, an operation facet, a lowering facet |
| Product | One kind of output requested from the same selections: native type declarations, or a **binding product** such as value adapters, JSON presentation or operation data. The accepted designs use the word without defining it; this is our reading |
| Selection, instance | A **selection** is the planner's choice of one rule for one use. An **instance** is the accepted design's "representation instance": the immutable result of one selection for one use. Quoted "generic instance" means an application such as `Page<Person>` |
| Plan, plan entry | The **plan** is the "selected representation plan": the typed record of every selection of one request, published for backends. It does not exist yet. The Go type `GenerationPlan` is something else. A **plan entry** is one item of it; whether there is one per use or per distinct selection is open, and section 6.4 names the candidates |
| Envelope, payload | The **envelope** is the part of the plan whose shape all families share: identities, links, versions. A **payload** is the family-owned part of an entry: an ontos value, typed by a horos type that the family declares |
| Decision record | Our word for what a rule returns for one use under the shape we would build. It is not a horos record type: "the record rule", "the record family" and "Record" in quotes always mean the horos record type |
| Slot | A declared place in a rule's result that stands for the representation of another use. The planner fills it |
| Demand | The owner's word (2.1): a request from a contribution that the kernel resolves to whatever the consumer registered. Today there are two kinds: "represent this child use" and "the helper under this key" |
| Generic environment | The mapping from a generic declaration's parameters (its **binders**) to the actual arguments of one application: `{T -> Person}` for `Page<Person>`. An **adapter argument** is the adapter selected for an actual argument |
| Native symbol | The accepted design's word for a reference to a named native entity that generated code must spell and import: a generated declaration, a helper, an external type, a factory. The design does not say whether it is an identity or a spelled name with an import path (6.5) |
| Backend, factory | A **backend** is a program, one per target and owned outside bittype, that reads the plan and emits the code that composes conversions. A **factory** is a run-time function that builds an adapter from a payload and the adapters of the slots |
| Completeness test backend | A backend written in Go inside bittype's tests: the acceptance instrument of the plan, never shipped |
| Conformance vectors | Test vectors that bittype publishes per family, stated against a horos descriptor. Every implementation of the family's conversion, in any repository, must pass them |
| Source semantics | The checked semantic graph of the horos modules. A backend reads it beside the plan. No document we know gives it a serialized form |
| Files stage, placement | The stage after generation that puts generated declarations into files, packages or modules according to a **placement** the consumer supplies (2.6) |

### 1.3 The stages, in one picture

```text
horos modules
   │  parse, resolve, analyse                      (horos frontend inside bittype)
   ▼
checked semantic graph of the source
   │
   │      registered by the consumer's composition:
   │        one source integration   (horos)
   │        registered units         (per target and profile)
   │        one target builder       (per target)
   ▼
kernel: for each use reachable from the requested declarations,
        select one rule by the use's selector, run it,
        resolve the rule's slots by planning the uses they name
   │                                   │
   │ native declarations               │ selected representation plan
   │ (native intent)                   │ (accepted design, not implemented)
   ▼                                   ▼
target builder: syntax tree            binding backend, per target, outside bittype:
   │                                   reads the plan, the source semantics and the
   ▼                                   contributions' factories; emits imports,
Files stage: files, imports            factory calls and glue
(TypeScript today; section 2.6)
```

The left branch exists, with a Files stage for TypeScript only, and is what
section 3 shows. The right branch is the accepted design of section 2. The
contract in question is the arrow from the registered units into the kernel:
what a contribution is, what its rule returns and what it may ask of the
kernel.

### 1.4 Scale and direction

**Targets.** Four. TypeScript and Go are to get conversions. Rust and C++ keep
types only for now, and their binding support is to be reported as absent.

**The six families that exist, as types only.**

| Family | What it is in horos | Packages today | Contributions under our leanings (6.1, 6.2) |
| --- | --- | --- | --- |
| Text | The builtin form `scalar` with discriminator `utf8-text` | `text`, in each target | One per target |
| Booleans | The builtin form `scalar` with discriminator `bool` | `boolean`, in each target | One per target |
| Lists | The builtin form `list-of`, with one operand | `list`, in each target | One per target |
| Closed records | The builtin form `record` with discriminator `closed`; each field has a name, a flag and a type | `record`, in each target | One per target |
| Optional fields | No type: the flag `opt` on a field of a record form | Go: `optional`, which supplies only the helper declaration. Elsewhere none: inside `record` | One per target, if question 1 is answered as we lean. None under option D1 |
| Generics | The parameters (binders) of a named declaration, uses of a binder, and applications | `declarations`, in each target, wrapping one shared recipe; plus planner code for parameters and applications | The `declarations` contribution for "record or alias"; the links between binders and arguments are the envelope's and the planner's |

That is 21 rule packages, five per target and six for Go, of 25 to 78 lines
each. The kernel's contract package has about 1,900 lines of Go and the
registry about 270. So moving every rule twice would be cheap inside bittype.
The cost that "once" guards against lies in other repositories.

**Later families**, in this order of need: integers, bytes, maps and sets,
exact values (a type with one value, used as a discriminator) and text
refinements, tuples and sequences, unions, nullable values, recursion, open
records, explicit specialization (a separate representation for one
application of a generic declaration).

**Other authors.** The service library is to contribute exact mappings for
`Function`, `Stream`, `Service` and `Cell`; Content two or three domain types.
The design promises that they need no change to bittype. How a library in a
TypeScript repository authors a contribution is not decided: a registered
unit today is a Go value in bittype's process, and the result table of
research 0014 says a contribution is "a Go mapping package or a declarative
sidecar". No sidecar exists.

**Backends and process boundaries.** The TypeScript backend belongs to the
service library, runs in Node and reads a serialized plan. The owner of a Go
backend, a Go runtime library, does not exist yet; bittype owns only the
completeness test backend. Kernel, registered units and the Go test backend
share one Go process, where an identity can be a pointer. The TypeScript
backend and factories do not: whatever they need must be a value in a
serialized artifact.

**The steps, and which of them are "the first slice".** Standing: our
derivation (the order is argued, with what speaks against it, in 2.4).

1. The contract changes. Every existing rule of the four targets moves to it.
   Behaviour stays types-only.
2. If our finding about optional fields holds (6.2), the optional
   representation moves out of the record rules, as its own change.
3. Additive work for the six families: the typed plan at an explicitly
   unstable version, a payload schema per family, conversions for Go inside
   the completeness test backend, conformance vectors, the serialized form of
   the plan, and a report of binding support.

"The first slice" is all three steps. Its targets are TypeScript and Go for
plan entries, payloads and conversions, and Rust and C++ for step 1 only. Its
products are the unchanged native types, the plan, Go conversions as test
code, the vectors and the report. The TypeScript factories are written in the
service library against the payloads and vectors of step 3, so step 3 is when
the first dependent outside bittype starts to build on the contract. The
acceptance example is `Page<Person>`: from horos source to native types, a
plan, and a round trip of values in Go.

## 2. Architectural context

### 2.1 Standing architecture rules (AC = architectural correctness)

The owner's rules for all platform work are numbered AC-01 to AC-10. They are
quoted from the platform repository's `docs/concepts/architectural-correctness.md`
at commit `5e988b0`. The questions rest on these seven. Standing: **owner
rule**.

| Rule | Requirement, verbatim |
| --- | --- |
| AC-02 | "Let independently owned implementations supply their domain interpretations, policies and laws. Shared mechanisms own their invariant guarantees and depend on extension contracts; consumer composition selects the concrete contributions." |
| AC-03 | "Within the agreed extension domain, adding a contribution must change only its owning implementation and consumer composition, not shared orchestration, closed case catalogues or unrelated implementations." |
| AC-04 | "Builtins and the first implementation must use the same public extension contract available to independently supplied implementations." |
| AC-05 | "Preserve identity, scope, ownership, authority, assumptions, evidence association, resource limits and failure meaning through composition, as applicable to the contract. Selecting a contribution does not establish the truth of its claims." |
| AC-08 | "When a shared contract really must change, explain the missing invariant or capability, why its responsibility belongs to that layer and how the revised contract serves independently owned contributions." |
| AC-09 | "Whenever we design, redesign or improve an interface, choose the best contract for the domain requirements, responsibilities, ownership, composition and guarantees. The currently implemented surface must have **zero influence** on that choice. Implementation or migration effort has **zero weight**, including as a tiebreaker." |
| AC-10 | "Move a misplaced or wrongly named responsibility to its owner first, as its own change, without waiting for the complete implementation of its destination. That first change establishes the final structure: its inputs and identities have their intended shape, it keeps every guarantee the destination's shared mechanism requires for the cases it admits, and it explicitly refuses the cases it cannot yet check. Missing coverage is then added without moving, renaming or removing what exists." |

Three further passages of the same document shape the questions. Standing:
owner rule.

**What violates AC-03, AC-05 and AC-10:**

> The addition needs a shared switch/type/registry/bootstrap edit, a weaker fallback or a parallel copy of the orchestration to admit its meaning.

> Openness is achieved by accepting arbitrary payloads, unchecked claims or callbacks that discard the guarantees the mechanism promises.

> Delivering a first change whose input or identity later has to change.

The second sentence matters for the question on facets: an open set must not
become "arbitrary payloads". The third is the reason for this consultation.

**Wrong coupling, and the demand that corrects it** (owner request of
11 November 2026, section "Package boundaries and forbidden dependencies").
A table distinguishes two faults. The second, with its cells joined into one
paragraph:

> Wrong coupling. What it is: A contribution names another contribution that the consumer is meant to choose. What is wrong: The dependency. The caller makes the shared mechanism's choice. Consequence: The mechanism's selection, scope and failure guarantees do not apply to that call. Replacing the callee means editing the caller. Correction: A demand that the shared mechanism resolves.

And below it:

> Not every direct call is wrong. Ask of each call from one part to another whose answer the caller needs: its own, or the composition's. A caller that relies on its own invariant, a shared mechanism or a vocabulary calls directly. A caller that needs whatever the consumer registered for a meaning states a demand and lets the shared mechanism resolve it.

The same section defines the unit (two list items joined):

> A contribution, meaning one thing a consumer registers or selects on its own, is one unit. It depends on shared contracts, the vocabularies of its own scope and its own private parts. Only a composition depends on a contribution.

**Evidence for an extension boundary:**

> When qualifying a claimed extension mechanism, implement a materially different contribution using its public contracts in a separately owned package or fixture. Keep the shared mechanism fixed. [...] A renamed copy of the original example, a private test-only constructor or a fixture that activates a preinstalled special case does not establish that.

**bittype's own finding under that rule.** Standing: **our derivation**. It
was written and confirmed by the seat that writes this document, in bittype's
record of applying the package-boundary request, `docs/package-boundaries.md`.
The owner's rule is the general one above; that the following is a case of it
is ours, and section 6.2 gives the reading under which it is not. The record
says (written against an earlier commit; the Go detail has changed since and
is shown in section 3.6, as is the call `r.Mapped`, which asks the kernel to
represent a child use):

> In a Record mapping, `r.Mapped(f.Use)` asks the kernel for the native type of a field. Whichever mapping the consumer selected answers. This is the correct form.
>
> Two lines below, the same Record mapping decides how an optional field is represented. Rust builds `std::option::Option`, C++ builds `std::optional`, and Go wraps the type in a helper through a package variable. A consumer who wants another representation of optional fields must replace Record.

Its findings table records the case as "Record decides optional fields itself
in Go, Rust and C++", kind "Coupling", state "Open", tracked by item 6 of
issue #196 (section 2.4). Its assessment cell says "Confirmed; the contract
change is designed"; section 2.4 says why we no longer hold the second half.
No import rule can see this kind of decision: "Record building
`std::optional` imports nothing forbidden."

### 2.2 Binding generation: who owns what

Source: bittype `docs/binding-generation.md`. Its first line: "Accepted
direction, 9 November 2026, from research 0014. Not implemented." Research
0014 was an expert consultation like this one; its owner, another seat of our
team, evaluated the advice and recorded the result as this page. Standing of
everything quoted from the page in this subsection: **accepted design**
(binding generation, 9 November 2026).

**What it is.** From a horos type of a service a build needs more than native
types: "conversion between native values and hydrated message values"
(message values in which references to remote endpoints are live objects),
conversion to and from JSON for presentation, operation data, shape checks
and declared operation failures. "That set is **binding generation**."

**The decision** (list items joined, links reduced to issue numbers):

> One semantic owner does not require one emitter. It requires an authoritative decision, an explicit interface carrying that decision, and tests that establish that its consumers preserve the promised meaning.
>
> **bittype** remains the authority for selecting native representations. It publishes each selection as a typed **selected representation plan** (#195).
>
> **A binding backend per target** consumes three inputs: the plan; the source semantics; each library's **paired contribution** (#196). From these it emits static composition and typed glue.
>
> **bittype does not emit runtime adapters or own their behaviour.** Its charter gains only the plan interface and the hooks that produce it.

**The owners**, verbatim, without the row for the transport runtime:

| Owner | Responsibility |
| --- | --- |
| horos | Declaration semantics, checked identities and applications, descriptors, ground conformance |
| bittype planning | Exact contribution selection, substitutions, references, and the selected representation plan |
| Each representation implementation (a target's structural profile, or a library's exact mapping) | Native layout, presence rules, numeric representation, construction and access |
| Binding backend, one per target | Adapter composition, runtime imports, generated glue, assembly of the requested binding products |
| Contract or domain library | Its native recipe, its adapter implementation, its protocol or domain meaning, and the laws below |
| Consumer build | Selection and composition of the concrete contributions it uses |

Two rows carry the shape we would build. The planner owns selection,
substitution, references and the plan. The representation implementation owns
"Native layout, presence rules, [...] construction and access". Nothing in
the table gives the planner a decision about layout, and nothing gives a
representation the right to issue a plan. Note "presence rules" in the third
row: it is one of the texts behind option D1 of section 6.2.

**The plan.** "A manifest of names is not enough. For each use and generic
instance, the plan carries" a list that includes "the selected contribution,
target and representation profile", "native symbols, distinguishing run-time
imports from type-only imports", "field access and construction", "the rules
for absence, presence, null and native `undefined`", "recursive references"
and "the selected adapter factories, their argument slots and their resource
requirements". ("Generic instance" is an application such as `Page<Person>`.
Resource requirements are run-time limits and play no part in the
questions.) Its form is "typed references and target-owned expressions", and
it "starts as a typed build-time interface and is serialized only where a
real process boundary needs it. A TypeScript backend in Node is such a
boundary." Then:

> **Completeness test.** The backend can implement each selected representation from the plan, the source semantics and the selected contributions alone, without guessing a rule. If it must decide how Go optional values work, or how a custom union is represented, the plan is incomplete.

And on value-level code:

> **Interpreters.** A shared structural interpreter is allowed only if it executes the plan's decisions rather than choosing representation rules itself.
>
> **The emitter's scope.** It is scoped separately, with its own syntax and compiler checks. bittype's native vocabularies stay types-only.

"Native vocabularies" are the native forms and the syntax node kinds of each
target. They can express a type and a type declaration. They cannot express
a function, a call, a field access or a constructor. Backends are placed with
the owner of each target's adapter interface: the TypeScript backend with the
service library, and "A Go backend belongs to the future Go runtime library."

**Paired contributions.** "A library publishes, for one declaration and one
supported target representation:"

| Element | Meaning |
| --- | --- |
| Source selector | Exact module digest and declaration identity |
| Representation context | Target and supported profile |
| Native facet | Its bittype recipe or native intent |
| Value adapter facet | Run-time factory reference and typed argument-adapter slots, converting the value inside a message |
| Root (attachment) facet | Serving the value as one routed endpoint; Function's is `functionBinding`, its value facet `functionAdapter` |
| Requirements | Adapter-interface compatibility, resource facilities, supported scope |
| Optional product facets | For example an operation facet, through which a callable contract exposes its argument, result and failure positions |

(A root facet serves a value as an endpoint of its own and converts nothing;
`functionBinding` is the service library's name for Function's.)

> **Selection.** A native facet and the facets a binding product needs are selected together. A native mapping without an adapter facet still supports types, but a request for bindings that needs a missing facet fails explicitly. There is never a fallback to structural conversion, in the same way that an exact use mapping never silently becomes a structural publication today.

(The last clause describes today's planner: a use rule registered for a named
declaration is applied or fails, and the planner never falls back to the
default declaration recipe; section 3.4.) The table ends with "Optional
product facets" and an example. It neither closes the list of facet kinds nor
says who may add one.

**Ownership and composition.** "Contributions belong to their library's
owner, outside canonical horos bytes. Native package paths and export names
are not part of declaration identity." Composition must "select by exact
identity before any structural expansion" and "reject duplicates,
incompatible profiles, missing facets and unsupported applications". The
design does not say what makes two contributions duplicates, nor who compares
them.

**The unit.** "The unit of implementation is library × target." This is the
page's only sentence on granularity. Background, with no standing of its own:
the expert's answer to research 0014 named the property behind pairing as
"joint ownership and selection, not whether the contribution is written in Go
or metadata".

**The laws every contribution must satisfy.** With `V_T` the admissible
abstract values of a type T, `N_T,p` the admissible native values under
profile p, E for encoding and D for decoding:

```text
D(E(n)) ≈ n    for every n in N_T,p   (the equivalence states whether copying, class identity, branch tags and sharing matter)
E(D(v)) = v    for every v in V_T     (abstract-value equality; any normalization is specified)
```

plus eight further laws. They matter here for one reason: what a family
records must contain whatever these laws depend on. Two examples from the
first slice: in Go, `Value` is ignored when `Present` is false, so `≈`
identifies two absent values with different `Value`; in TypeScript, an absent
property and a present `undefined` are different native states. Neither fact
is visible in the native type.
