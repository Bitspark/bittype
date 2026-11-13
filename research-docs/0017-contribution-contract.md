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

### 2.3 Representation families: the contract, the instance and the plan

Source: bittype `docs/representation-families.md`. Its first line: "Accepted
design, 9 November 2026. Not implemented." It was decided by research 0015, a
second consultation by the same research owner, and "details the selected
representation plan and the paired contributions of binding generation for
concrete value families". Where the two pages differ, representation
families is the later amendment. Standing of everything quoted from the page
in this subsection: **accepted design** (representation families, 9 November
2026). Quotations of the expert's advice are background and have no standing.

**Delivered complete:**

> A family is delivered with its native type, its plan entries, its value conversion and its qualification together. It is never a types-only mapping that a later binding step reopens. The conversion constrains which native types are admissible at all: a TypeScript `number` cannot hold every horos integer.
>
> **The promise.** Adding a family through the published identity, dependency, invocation, plan-graph and target-extension contracts changes only that family and the consumers that compose it. Existing representation identities keep their meaning.

The page names these five contracts and defines none of them further. In our
reading, identity is the subject of question 2, dependency is the demand of
question 1, plan-graph is the envelope, target-extension is the facets of
question 3, and invocation is how a backend calls a factory.

**The representation contract.** "A family contributes, per target, one
contract with seven parts:"

| Part | What it defines |
| --- | --- |
| Applicability | The source identities and forms it accepts, its generic assumptions, its targets, its explicit refusals |
| Native representation | Native type, dependencies, construction, access, discrimination, admissible native states |
| Value meaning | Promised canonical domain, domain predicates, native equivalence, any permitted normalization |
| Selected-use payload | A family-owned, versioned schema for the decisions recorded at each use |
| Conversion implementation | Run-time factory reference, child adapter slots, runtime compatibility, conversion laws; optionally a lowering facet (below) |
| Operational requirements | Budget use, mutation and ownership, recursion support, sharing behaviour, failures |
| Qualification | Separate evidence for type generation and binding support, tied to the representation and implementation versions |

This is what a contribution must supply. The seven elements of section 2.2
list what a library publishes; the two lists describe the same thing from two
sides. Note "its explicit refusals" under applicability, and "per target" as
this page's word on granularity.

**One instance, two projections.** The sentence the shape we would build
rests on:

> The planner selects one immutable representation instance per use. The native type and the plan entry are both projections of that instance. No later pass reconstructs a decision from the source, so a native declaration and its conversion cannot disagree.

The advice behind it:

> That prevents a subtle failure mode: generating a native declaration, then independently reconstructing conversion decisions from the source. Two passes can read the same declaration and still disagree about a nullable pointer, a union wrapper, or a generic specialization.
>
> “Defined once” should mean that representation decisions have one authoritative origin. It need not mean that a Go planner and a TypeScript runtime somehow share one executable implementation. Their agreement still needs a shared contract and verification.

**Factories and the lowering facet:**

> By default a plan entry names the family's run-time factory and gives it the family's payload as arguments. A family may also supply a target-owned **lowering facet** that turns its recorded decisions into typed target expressions or helpers, such as generated accessors or constructors for Go. The binding backend owns imports, names, placement and assembly. It never switches over families, and source-code strings never stand in for typed target expressions. A shared interpreter is a valid factory implementation only if its handlers are registered contributions executing selected decisions. Generated source is never modified as an extension mechanism.

The advice behind it says why a factory with plain data cannot be the only
path: "A TypeScript runtime can access a field from its name. A Go family may
instead need generated accessors, constructors, typed closures, or helper
declarations."

**The plan**, verbatim:

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

Followed by five properties:

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

(`ContentValue` is a recursive type of the Content service and `Outcome` a
union with a constant tag field.) Four remarks for a reader outside the
project.

- "Record field layout and presence behaviour" stands in the payload column
  as one item. It is the second text behind option D1 of section 6.2.
- "Selected contribution and profile identities" requires an identity of a
  contribution that survives serialization. The page does not say what it is.
- "Typed dependency slots" requires every link between two entries to be
  visible to the kernel. A family that asks for another's answer must do it
  in a way that produces a link.
- Vertical case 4 is the lowering facet's test. Representation families
  assigns it to no slice. The six cases were set for the envelope; the
  contribution contract has no such list.

**Generic requirements.** "A native type parameter does not carry the
adapter, equality, discriminator or layout evidence needed to bind values of
that type". For `Page<T>`: "`T`'s selected adapter travels with its native
type". For a nullable `T`: "The null representation cannot collide with a
valid `T`", which is a requirement that one family places on whatever another
family selected for its operand.

**Absence and null**, the family behind question 1:

> An optional nullable field has three states. Go uses `Optional[Nullable[T]]` with defined treatment of unused payload fields; a pointer form needs proof that its `nil` cannot collide with a valid value (`**big.Int`). TypeScript uses `field?: T | null`, and the adapter separates an absent own property from a present `undefined` at run time.

The page also keys a conversion's memo "by the selected representation
instance plus input identity"; section 2.5 quotes the service library's
version of that sentence, on which question 2 rests.

### 2.4 The work items, the order and the first slice

Four open issues of the bittype repository carry the work. Their texts are
ours. Where they restate an accepted design they add nothing. Where they go
beyond it, standing: **our derivation**.

**Issue #196, item 6** (added 11 November 2026) is the contract change.
Verbatim, with one stale file path corrected in brackets. In it "a binding"
is a use rule and "the mappings" are the registered units:

> **The existing mapping contract changes once, here.** [...] A binding's `Expand` returns only a native `Expression[T]` (`contracts/generation/go/description.go:26-31`), and a declaration recipe returns only a `DeclarationDescription[T]` (`:61-66`). A decision such as "this field is optional" is made inside the callback and recorded nowhere ([today `transform/syntax/horos-to-golang/mapping/record/go/record.go:70-72`]). Item 1 of #195 requires the native expression and the plan entry to be projections of one selected instance. So what a contribution returns has to change, and the mappings of all four targets implement that contract. The change is made in this slice for every existing mapping, including Rust and C++, so that later families add contributions and change no existing contract. Derive the new result from the domain (AC-09). This item prescribes no signature.

"Recorded nowhere" is imprecise; section 3.7 gives the exact state. Its
verification: "After this slice, a fixture family is added by registering one
contribution. The diff of that test touches no file under
`contracts/generation/go` and no existing mapping."

The same day this seat added a delivery note to the page of section 2.3, in a
pull request that states it changes "No code, no contract and no accepted
design text": "Two existing parts change once, ahead of the families: [...]
and what an existing mapping returns (item 6 of #196). After them a family
adds a contribution and changes no existing contract." So "changes once" has
the standing of our derivation from AC-10, wherever it stands.

**Issue #195** is the plan. Its items restate section 2.3 and add: "A payload
that fails its type is refused before any backend runs" and "A test backend
in this repository, in Go, implements each supported representation from the
plan, the source semantics and the selected contributions alone. A plan that
omits a decision fails the test." It keeps vertical case 4 for itself.

**Issue #207** is the six families (items 1, 5 and 9 of eleven):

> 1. **Contract.** Each family is one representation contract with all seven parts, contributed through the public interface of #196. The kernel, the planner and the completeness test backend gain no case for a family.
> 5. **Existing output unchanged.** The native declarations generated for today's pagination fixtures stay byte-identical in all four targets. Existing binding identities keep their meaning. The profile's own identity is settled first, in #216.
> 9. **Optional fields.** The TypeScript adapter separates an absent own property from a present `undefined` at run time. Go keeps `BittypeOptional[T]` and ignores `Value` when absent [21].

("Binding identities" are the rule IDs such as `list/golang/1` of section
3.4. "[21]" is a row of research 0015's result table, which asserts of
today's Go type that it "already ignores `Value` when absent"; no code states
it.) Its verification: "`Page<Person>` goes from horos source to native
types, a plan and a round trip in the Go completeness test backend, with no
rule guessed by the backend." and "A seeded mutant that builds a plan entry
in a separate pass from the native expression turns the completeness test
red".

**Why this order** (the steps of section 1.4; standing: our derivation). The
contract change comes first because the optional move needs a request whose
answer has the new contract's shape; a request built on today's contract
would return native intent and change again. Everything in step 3 is then
additive. Until step 2, the coupling of section 2.1 stays an open, recorded
violation.

**What speaks against this order.** Step 1 makes the contract final while
nothing reads its new half. With types-only behaviour and unchanged output,
the decision record, the request and the facets have no consumer until step
3, and no reader in another process or repository until the service library
writes its first factory. A contract can pass step 1 and be wrong in every
part that the first slice exists to add. The prototype that our reading of
"once" allows (Question section; 6.6) is the obvious remedy, and it changes
the order: one target end to end first, the other three after.

**A caution about our own reading.** On 12 November 2026 we first recorded
what was left of item 6 as "details inside a settled design" (#196, 19:38Z),
then found by reading the accepted pages against the first slice that several
shapes are given nowhere (#196, 22:40Z). Earlier that day two design
questions in this line of work had been answered from the code, wrongly. This
document may make the opposite error and read more into the pages than they
say.

### 2.5 The other side of the interface: the service library's binding contracts

Source: bitspark-service `docs/design/binding-contracts.md` at commit
`c36722db901ef7a259242e168c4c35d6f777fbf4`. Its status line: "accepted by
libraries/service on 11 November 2026. Nothing here is implemented yet". It
was written and accepted by the owner seat of that repository. Its laws on
capability lifetime (group L) are under revision in bitspark-service#47;
bittype has confirmed the group on the call context (K) and not group L.
Standing: **accepted design** (binding contracts, 11 November 2026), with L1
marked unconfirmed.

**What a backend pins.** The service library's issue for its TypeScript
backend (bitspark-service#40): "It reads one pinned serialized plan version.
It refuses an unknown envelope version, payload schema, semantic version and
run-time requirement, each shown by its own test", and "It never selects or
re-derives a representation, never switches over families [...] and contains
no Function, Stream, Service or Cell name." Its structural run-time package
(bitspark-service#38, decision of 11 November 2026): "The package makes no
representation decision: every choice arrives as the family's plan payload.
Its CI runs bittype's conformance vectors for each family it implements". So
the TypeScript factory for bittype's list family is code in another
repository, found from a plan entry by an identity and held to the same
meaning by shared test vectors. No compiler connects the two.

**Conversion state and the memo key.** Every factory receives the call
context, which holds the memo:

> **The memo key** is the selected representation instance (its descriptor and its generic environment) plus the input's identity. An object-only key is unsound when one object appears under two representations.

Law K3, a table row in the source: "One object under two representations
converts once per representation." Law L1, unconfirmed, uses the same unit
for exporting a capability (making it callable by the receiver): "one export
per message, per selected representation instance, per native value".

**The root facet of `Service` reads the plan.** A field of the record `Ops`
is a nested record of operations (a namespace) or one operation (a leaf): it
"is either a namespace, which is a field the plan selects as the structural
record family, or a leaf whose selected contribution supplies an operation
facet". So one contribution's facet is given what the plan selected for
another use and must recognise "the structural record family" by it: a second
reader of a contribution's identity besides the backend.

**"Instance" in two senses.** Our pull request #260 brings bittype's two
design pages into agreement with these contracts; it is open. Its second
reading (12 November 2026, 18:50Z; line numbers are those of its branch)
found:

> `docs/representation-families.md:44-45` says "one immutable representation instance per use", and `:259` keys the memo by "the selected representation instance". The page now uses "instance" in two senses.
>
> Suggested: "Two uses that select the same contribution and profile for the same descriptor and generic environment have the same run-time identity. The memo key is that identity plus the input's identity (K3)." Say that this identity is coarser than the per-use instance.

Issue #195 gained the matching criterion the same day (08:11Z): "Two uses
with the same selected representation, descriptor and environment share one
memo key." So `People = Page<Person>` used in two fields is two uses and, at
run time, one memo identity. Section 6.4 gives each of these identities a
label.

### 2.6 Names are spelled at generation; places are decided after it

Three statements of different standing bear on what a plan can record about
a generated native declaration.

**(a) Generation decides no place.** Standing: **accepted design** (bittype
`docs/source.md`, "Design decision, 7 November 2026"). bittype has a stage
after generation, called **Files**. The decision divides the work: the
selected mapping and the target builder "Build declaration bodies and type
expressions. Bind each named use to its intended target declaration, lexical
parameter, or selected external public symbol." The consumer must "Supply one
placement for the declarations being emitted: output files and the native
scopes they belong to." The target's Files formatter (the Files stage's
implementation for one target) must "Choose the local spelling, qualification
or alias, and construct the necessary import, export or other module syntax."
One generated set serves every placement: "Both placements use the same
generated declarations, bindings and publication intent." TypeScript has a
Files stage today; Go's and Rust's are being built.

**(b) Identity, naming intent and realized spelling are three things, and
each language's Files stage realizes the spelling.** Standing: a
recommendation of research 0016 (a consultation on the Files stage) that this
seat accepted as research owner on 12 November 2026; its design page
`docs/naming.md` is on the main branch since 13 November 2026, written and
decided by this same seat and read twice by an independent reader. Nothing of
it is implemented. From the page:

| | What it is |
| --- | --- |
| **Declaration identity** | "Which declaration this is, within one sealed output. It alone is authoritative for references. It is never a spelling" |
| **Naming intent** | "What generation records about a name: the chosen name, the module's identity, the policy's identity and the spelling the policy produced" |
| **Realized spelling** | "The identifier that a language's Files stage publishes. In the first slices it is the spelled intent, once that language's rule has checked it" |

A "sealed output" is the immutable result of one generation run. The
"policy" is a naming policy that the consumer's composition selects: with or
without a qualifier that says which module a name comes from. (That
"qualification" of a name is unrelated to the "Qualification" of section 2.3,
which is evidence.) Who spells, in the page's words: "The mapping applies it
at generation", and "The unplaced tree's name is that spelling, and the
builder's comparator and every Files stage read it there." Each language then
checks the spelling or refuses it and "renames nothing": "in the first slices
a target refuses a spelling it cannot publish. An escape rule is a later,
explicit rule that turns one refusal into a success." And "A name does not
change between placements."

**(c) What an intent records, and how a consumer learns a place.** Standing:
**our derivation**, as decisions of this seat on issue #274 (12 November
2026, 22:34Z and 23:05Z). "An intent contains the mapping's chosen name and
its qualification, kept apart." "A consumer learns the realized name and the
placed site from a reader on the formatter's result, keyed by the output's
handles." A declaration **handle** is an identity that a sealed output is to
issue for each declaration it contains (issue #272, not implemented): "A
handle issued by one `Output` is refused when it is resolved against another
`Output`. That includes a second run of the same plan."

**What follows for the plan.** In the first slices the spelled name of a
generated declaration exists at generation and does not depend on placement,
so a plan entry can hold it. A later rule under which a language realizes a
different spelling would be an addition. What a plan entry cannot hold is the
place: the file, package or module, and so the import path that a backend's
glue must write. A member's name, such as a record field's, stays with the
declaration and is fixed at generation.

Section 2.2 lists among what the plan carries "native symbols, distinguishing
run-time imports from type-only imports", the envelope of section 2.3 lists
"Native symbols, imports, factory references, runtime requirements", and a
backend works "from the plan, the source semantics and the selected
contributions alone". For an import of a generated declaration these
sentences and (a) do not fit as written. The naming page states the tension
and leaves it to this consultation, with one constraint of its own: "Which
identity a plan entry records is not defined here: an output-issued handle
cannot be the key of a serialized plan." Section 6.5 gives the alternatives.

## 3. Relevant code: the contract as it is today

Everything in this section has the standing **implemented today**. It shows
what a rule can and cannot say now. It is not a constraint on the answer.

### 3.1 Reading guide

| Subsection | The one thing to take from it | Feeds question |
| --- | --- | --- |
| 3.2 | A rule's whole result is native intent: a tree of native forms, string data and three kinds of placeholder | 1 |
| 3.3 | A registered unit is a profile plus three lists of rules. A callback never sees what was selected for a child use | 1, 2 |
| 3.4 | Four names exist for selection and none of them survives the process | 2 |
| 3.5 | A rule is a constant template, a template with slots, or a callback that reads the source. The declaration recipe spells the native name | 1, 2 |
| 3.6 | "This field is optional" is decided inside each target's record rule, in four different ways | 1 |
| 3.7 | Per use the planner keeps resolved native intent and debug strings, and drops the selection | 2 |
| 3.8 | Two neighbouring mechanisms already register per-family parts apart from the rule, and record a native symbol by identity | 2, 3 |

In the excerpts, `F`, `T` and `R` are type parameters: marker types of the
source and target language, and the type of the source integration's reader.
Elisions are marked `// ...`. Aliases and abbreviations:

| Name | Meaning |
| --- | --- |
| `g` | The kernel's contract package, `contracts/generation/go` |
| `h`, `s` | The syntax vocabulary of horos, and of the target at hand |
| `model`, `source` | The target's native model; the horos source integration |
| `rules`, `patterns` | Common packages of the horos rule side: the profile names and one function type; the source patterns |
| `m`, `b` | The shared tree-building package; the syntax tree under construction in a target builder |
| `Desc` | Abbreviates `g.MappingDescription[h.Language, s.Language, source.Reader]` |
| `Rule` | Abbreviates `g.BindingDescription[h.Language, s.Language, source.Reader]` |
| `Exp`, `Expr` | Abbreviate `g.Expansion[h.Language, s.Language, source.Reader]` and `g.Expression[s.Language]` |

An `ast.NodeRef` is a site in a syntax tree, used for diagnostics, and
`g.At(e, site)` attaches one. `from` and `to` are language identities: a
language with a syntax revision that its owner calls a "profile", unrelated
to the representation profile. Code comments use words that our prose has
replaced: "provider" (a registered unit, or in a file name the target
builder), "Run" (executing the plan closure), "composer" (the kernel),
"occurrence" (a use), "binding" (a use rule; in `Generic.Bindings` a
parameter binding) and "capture" (binding a slot to a use).

### 3.2 What a rule can say: native intent

`contracts/generation/go/mapping.go:5-11`, `expression.go:9-14, 30-35, 41-44,
50-74`:

```go
// Selector is an exact source-owned identity, qualified by language/profile/marker.
// Its key is never interpreted by the composer. Display names and graph-local IDs
// are not substitutes for a source integration's stable declaration identity.
type Selector struct {
	language representation.Language
	key      string
}
// ...
// Form is a target-owned native semantic constructor. Only its issuer can
// interpret its key. This is not an AST kind or a globally registered type kind.
type Form[T any] struct {
	owner *identity
	name  string
}
// ...
// Slot identity is independent of its display name. One description declares
// its slots, and captures supply the exact source occurrences for them.
type Slot struct {
	owner *identity
	name  string
}
// ...
type HelperKey struct {
	owner *identity
	name  string
}
// ...
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

`identity` is a private struct whose address is the identity. A `Slot` or a
`HelperKey` equals only copies of the value that one call of
`g.TypeSlot(name)` or `g.NewHelperKey(name)` returned. A `Form` carries the
identity of its issuer, one per native model, plus a name: two forms of one
model with one name are equal, and no other package can make one. All three
are unforgeable and meaningless outside the process.

The six kinds, in our words:

- **Native**: a native form with string data and argument expressions.
- **Captured**: a slot. The planner replaces it with the native intent it
  plans for the use bound to that slot.
- **Parameter** and **Argument**: a generic parameter of the native
  declaration being described, and an actual argument at an application.
- **Helper**: a reference, through a key, to a helper declaration that some
  registered unit must supply.
- **Declaration**: a reference to a planned native declaration. Only the
  planner creates it.

A native model is a form issuer, a few constructors and a check. Three lines
of Go's (`transform/syntax/golang/model/go/model.go:12, 35, 43-45`):

```go
var forms = g.NewFormIssuer[s.Language]()
// ...
func Form(name string) g.Form[s.Language]      { return forms.Form(name) }
// ...
func Field(name string, optional bool, ty g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormField), []string{name, strconv.FormatBool(optional)}, ty)
}
```

The Go model has seven forms: `primitive`, `slice`, `field`, `record` and
`external` for expressions, `alias` and `defined` for declarations
(`FormField` and its like are constants holding these names). It admits two
primitives, `string` and `bool`, and its check refuses every `external`
expression. The other three models are parallel. So everything a rule can say
today is a type or a type declaration, in strings and forms whose meaning
only the issuing model defines.

### 3.3 What a registered unit is

`contracts/generation/go/description.go:13-52, 61-84`:

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
type ParameterDescription struct {
	Name        string
	Origin      ast.NodeRef
	RequiredUse bool
}

// Convention describes both sides of a generic, including reordered, omitted
// and additional native parameters. Bindings is indexed by source parameter;
// Arguments is indexed by target parameter.
type GenericConvention[T any] struct {
	Parameters []ParameterDescription
	Bindings   []Expression[T]
	Arguments  []Expression[T]
}
type DeclarationDescription[T any] struct {
	Name    string
	Form    Form[T]
	Origin  ast.NodeRef
	Generic GenericConvention[T]
	Body    Expression[T]
}
// ...
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
// ...
type MappingLibrary[F, T, R any] interface {
	Describe(from, to representation.Language) (MappingDescription[F, T, R], error)
}
```

`Use[F]` is the kernel's handle for one use; only the kernel creates one.
`UseView[F]` is its structural view of that use: its kind (a reference to a
generic parameter, a reference to or application of a declaration, or a
builtin form), its children, and whether it is the whole body of its
declaration. A `NativeDeclaration[T]` is a name, a declaration form, native
parameter names and a body expression. `RequiredUse` makes planning fail if
the native parameter is unused in the body.

A use rule (`BindingDescription`) has exactly one of two results: `Output`, a
constant template, or `Expand`, a callback. The callback receives an
`Expansion`, which offers three operations (signatures,
`description.go:113-136`); a declaration recipe's `DeclarationRequest` adds a
fourth, the kernel's view of the declaration:

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

`Mapped` returns a placeholder at once, and the planner plans the child use
after the callback has returned. So a callback never sees what was selected
for a child, and a child's native intent is never an input of its parent's
rule. A rule has no declared way to decline a use: it returns an error from
its callback (an example is in 3.6).

Every registered unit today holds exactly one entry in one of the three
lists, and writes its own profile into `Profile`. All are built as one Go
function type, `rules.Library[T]`, whose value returns the description.

### 3.4 Registration, selection and the names that exist today

| Name | Who issues it | Outside the process | What it selects |
| --- | --- | --- | --- |
| Selector key, such as `builtin/7:list-of/0:` | The source integration | A string, meaningful anywhere | The rules filed for a use |
| Rule ID, such as `list/golang/1` | The unit's author, as a string | A string, checked for uniqueness only within its own unit | With a registry handle, one use rule in a replacement |
| Registry handle, `MappingRef` | The registry: its own address and a counter | Nothing | One registered unit |
| `Slot`, `HelperKey` | Whoever calls the constructor: a unit, or the consumer's composition | Nothing | The use bound to a slot; the helper under a key |

`transform/syntax/shared/generation/go/registry.go:19-30`:

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
```

**Registration.** `RegisterMapping` calls a unit's `Describe` once, checks
the result (a non-empty profile, rule IDs unique within the one description,
valid selectors, declared slots, size limits) and stores it under a new
handle. Nothing is checked across two units.

**Selection by profile.** A request names exactly one profile, and the
registry plans with every unit that states that string and has compatible
languages. A unit that states another profile is not in the plan. So an exact
mapping joins the structural families today by stating their profile string
itself. With no unit the request fails as unsupported.

**The planner's name for a rule.** The registry hands each unit to the planner
under a name made from its counter (`registry.go:242`), and the planner
prefixes each rule's ID with it:

```go
		parts = append(parts, g.SelectedDescription[F, T, R]{ID: fmt.Sprintf("library-%d", e.ref.id), Description: d})
```

So a rule is known as, for example, `library-3/list/golang/1`. The first part
depends on the order of registration. The page that documents the generator
says of the second (`docs/type-generation-plugins.md`): "A binding's ID is
local to its mapping registration. It is not a semantic proof, a globally
unique implementation revision or a configuration digest. There is no
PluginID API."

**Selection per use** (`contracts/generation/go/planning.go:465-475`). The
planner takes the use rules filed under the use's selector key:

- More than one: the error "ambiguous generation mapping", never a
  preference.
- Exactly one: it is applied. A use rule for a named declaration is thereby
  an exact mapping.
- None, for a builtin form: the error "no mapping for" the key.
- None, for a named declaration: the use becomes a reference to the native
  declaration that a declaration recipe produces. A recipe with a zero
  selector is the default for every declaration.

**Replacement.** A consumer resolves an ambiguity with
`Registry.Replace(Replacement)`, naming the selector and two use rules, each
by handle and rule ID. The registry removes the `Expected` rule before
planning. Only use rules can be replaced. A declaration recipe cannot, and
two units that supply a helper under one key fail as ambiguous.

### 3.5 One Go family of each kind

**Text: a constant template.** The whole unit is one use rule with the ID
`utf8-text/golang/1`, a pattern for the builtin form `scalar` with
discriminator `utf8-text`, and the constant output `model.Primitive("string")`
(`transform/syntax/horos-to-golang/mapping/text/go/text.go`). Nothing about Go
strings beyond the type name is stated: not that a Go string may hold invalid
UTF-8, and not what a conversion does with one.

**List: a template with one slot.**
`transform/syntax/horos-to-golang/mapping/list/go/list.go:18-27`:

```go
func Library() rules.Library[s.Language] {
	slot := g.TypeSlot("E")
	return func(from, to representation.Language) (Desc, error) {
		pattern, err := patterns.Element(from, slot)
		if err != nil {
			return Desc{}, err
		}
		return Desc{Profile: rules.RepresentationProfile, Bindings: []Rule{{ID: "list/golang/1", Input: pattern, Output: model.Collection(g.Mapped[s.Language](slot))}}}, nil
	}
}
```

`patterns.Element` returns the pattern for the builtin form `list-of`,
declares the slot and supplies the `Capture` function, which asks the source
integration's reader for the list's element use. The template says "a Go
slice of whatever is selected for `E`": the list rule does not know what
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
		// decl the source integration's reading of it
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
		// ... one native parameter per horos parameter, in the same order
		desc.Body, err = r.Mapped(open.Body)
		return desc, err
	}}
}
```

`src.Descriptor(u)` is the source integration's reading of a use as a builtin
form (its constructor, discriminator and site) with a flag that says whether
the use is one; it is not the horos descriptor of section 1.1. The recipe has
a zero selector and states no profile; Go's `declarations` unit wraps it with
the forms `defined` and `alias` under the structural profile. It makes three
decisions inline that a conversion or a backend needs: the spelled native
name (capitalised by `naming.Exported`, and prefixed for a declaration
outside the request's entry module), whether the declaration becomes a named
record type or an alias, and how each horos parameter corresponds to a native
one. "Generics" as a family is this last part together with the planner's
handling of parameters and applications.

### 3.6 Where "this field is optional" is decided, in four targets

Look for who decides in each row:

| Target | Native form of `next?: text` | What the representation is | Who decides today |
| --- | --- | --- | --- |
| TypeScript | `next?: string` | A property of the field member. No type stands for it | The record rule, through a flag that the target builder renders |
| Go | `Next BittypeOptional[string]` | A generated generic type around the field's type | The record rule decides "a wrapper type"; a helper unit, wired by a key, supplies the declaration |
| Rust | A field of type `::std::option::Option<::std::string::String>` | An external type around the field's type | The record rule |
| C++ | A field of type `::std::optional<::std::string>` | An external type around the field's type, plus an `#include` | The record rule and the target builder, agreeing by spelling |

The horos source integration reads a record as a list of fields, each with a
name, a flag and one child use, the field's type
(`transform/syntax/horos/source/go/read.go:25-30`):

```go
type Field struct {
	Name     string
	Optional bool
	Use      g.Use[h.Language]
	Origin   ast.NodeRef
}
```

Only the field's type is a use that the kernel walks and selects for. The
flag is content of the record form. No selector exists for an optional field,
so nothing can be registered for it today.

**Go.** The record rule is a use rule with the ID `record/golang/1` for the
builtin form `record` with discriminator `closed`. Its constructor takes a
helper key, `Library(optional g.HelperKey)`. Its callback refuses a record
that is not the whole body of a declaration (`src.Unsupported(use.Use, "nested
records require a named declaration in this profile")`), reads the fields and
runs this loop
(`transform/syntax/horos-to-golang/mapping/record/go/record.go:61-75`):

```go
			for _, f := range fields {
				f.Name, err = naming.Exported(f.Name)
				if err != nil {
					return g.Expression[s.Language]{}, err
				}
				ty, err := r.Mapped(f.Use)
				if err != nil {
					return g.Expression[s.Language]{}, err
				}
				if f.Optional {
					ty = g.Helper(optional, ty)
				}
				members = append(members, g.At(model.Field(f.Name, f.Optional, ty), f.Origin))
			}
			return g.At(model.Record(members...), view.Origin), nil
```

The helper is supplied by a second unit, which is all there is of the
"optional fields" family in Go
(`transform/syntax/horos-to-golang/mapping/optional/go/optional.go:21-28`):

```go
func Library(key g.HelperKey) rules.Library[s.Language] {
	return func(from, to representation.Language) (Desc, error) {
		return Desc{Profile: rules.RepresentationProfile, Helpers: []g.HelperRecipe[s.Language]{{Key: key, Declaration: g.NativeDeclaration[s.Language]{
			Name: "BittypeOptional", Form: model.Form(model.FormDefined), Parameters: []g.NativeParameter{{Name: "T"}},
			Body: model.Record(model.Field("Present", false, model.Primitive(model.PrimitiveBool)), model.Field("Value", false, g.Parameter[s.Language](0))),
		}}}}, nil
	}
}
```

Neither package imports the other. The consumer's composition creates the key
and gives it to both (`integration/go/libraries_golang_test.go:26-34`):

```go
	optionalKey := g.NewHelperKey("horos/optional")
	return []g.MappingLibrary[h.Language, s.Language, source.Reader]{
		text.Library(),
		boolean.Library(),
		list.Library(),
		record.Library(optionalKey),
		declarations.Library(),
		optional.Library(optionalKey),
	}
```

**C++ and Rust.** The record rule builds the wrapping type itself: where Go's
loop has `ty = g.Helper(optional, ty)`, C++'s has
`ty = model.External("std", "optional", ty)` and Rust's
`ty = model.External("std::option", "Option", ty)`. The C++ target builder
compares the same two strings, `data[0] == "std" && data[1] == "optional"`,
to add the `#include`. So in C++ the knowledge "an optional field is
`std::optional`" is in two packages that agree by spelling. No unit exists
for the optional representation alone: a consumer who wants another one
replaces the whole record rule, and the target builder too if the new type
needs an include.

**TypeScript.** The record rule wraps nothing. The flag travels as the second
string of the `field` form, and the TypeScript target builder turns it into
the property's `?`. Go's, Rust's and C++'s target builders ignore that
string.

What the Go arrangement gives and does not give:

- The kernel resolves the helper. A key that no registered unit supplies
  makes planning fail at the first optional field, and an optional field is
  never emitted unwrapped.
- A consumer who wants another helper declaration registers another unit
  under the same key in place of this one. Selection is by wiring.
- The record rule still decides that an optional field is a type wrapped
  around the field's type. A representation of another shape (a pointer, a
  field pair, a modifier) needs an edit of the record rule.
- The page that documents this form calls it interim
  (`docs/semantic-generation.md`): "This form is interim: #196, item 6,
  changes how a record asks for the optional representation."

Later families raise questions that we cannot yet classify as the same or as
different. A nullable value is a type of its own in horos, so it is a use
with a selector, and a nullable family is selected like a list. A map's
requirement on its key and "the null representation cannot collide with a
valid `T`" are requirements on what was selected for an operand. A union's
discriminator is content of the union form, like the optional flag.

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
// ...
type Step struct {
	Binding, Source, Declaration, Form string
	Expression                         string
}
// ...
// GenerationPlan is issued only after all reachable rules and dependencies have
// been resolved. Run does not retain a source, reader or mapping expansion.
type GenerationPlan struct {
	steps []Step
	run   func(context.Context) (Output, error)
}
```

`Program[T]` is the closed list of native declarations; only the target
builder receives it. `GenerationPlan` is a closure over it plus the `Step`
list, which exists for inspection in tests. `Output` is the sealed result of
running the closure: the checked syntax trees, the sites of the declarations
the request exposed, and the external uses of 3.8.

Applying a use rule ends like this (`planning.go:560-565`; `expression` is
the rule's template or its callback's result, `captures` maps the rule's
slots to uses, `f` is the generic environment, `b.name` the prefixed rule
name, `key` the selector key):

```go
	out, err := p.resolve(expression, captures, f, depth+1)
	if err != nil {
		return Expression[T]{}, err
	}
	p.steps = append(p.steps, Step{Binding: b.name, Source: key, Form: out.form.Name(), Expression: describe(out)})
	return out, nil
```

`resolve` replaces every slot with the result of planning the use bound to
it, and every helper reference with a reference to the helper's declaration.
`describe` prints an expression as a string. A finished declaration is stored
through `nativeDeclaration`, whose comment reads "nativeDeclaration
deliberately drops source binding/application templates" (`planning.go:756`).

So, per use:

| Fact | After planning |
| --- | --- |
| The native intent selected for the use | Kept, embedded in the parent's expression inside `Program[T]`. Form names and data are plain strings behind exported accessors, and the planner prints them into `Step.Expression`. Only the issuing model can authenticate a form and say what it means |
| Which rule was selected | Kept only as the string `Step.Binding`. The step list is flat: a step is linked neither to its use nor to the step of the parent use |
| Which use it was | Dropped. `Step.Source` holds the selector key, which is the same for every list in the program |
| Which child use filled which slot | Dropped when `resolve` returns |
| A use of a generic parameter | No step. The planner substitutes from the environment |
| A use of a named declaration that no use rule covers | No step for the use. One step per declaration, when it is first planned. The use becomes a reference with resolved arguments |
| The generic convention | Dropped by `nativeDeclaration`. An application keeps only its resolved native arguments |
| "This field is optional" | Kept twice in native intent: as the string `"true"` in the `field` form's data, and (Go, Rust, C++) as the wrapping type. Not kept as a decision: no code says that a Go `BittypeOptional` with `Present` false ignores `Value`, or that a TypeScript reader must tell an absent property from `undefined` |
| A refusal | An error from inside the callback, declared nowhere before the use is met |

### 3.8 Two neighbouring mechanisms

Two implemented mechanisms outside the generation unit already do what some
options of section 6 propose.

**Verification laws.** For Go, each family has a static check, run after
generation over the compiler's view of the generated types, that a native
type represents its horos type. Each is an independently selectable rule in
its own package, registered apart from the family's generation unit, with a
key that has a separate owner part
(`semantics.Key{Owner: "horos-go/static", Name: "optional-presence"}`) and a
versioned name (`horos-go/optional-presence-1`). Nothing links a law to the
rule it checks: the consumer's test composition names the helper type and its
fields a second time (`optionalLaw(native(t, g, "BittypeOptional"),
"Present", "Value")`) and the record's field correspondence a second time
(`{Source: "next", Target: "Next"}`). This has the shape of option F3 in
section 6.3 (a per-family part registered separately) and of I2 in 6.4 (an
identity with an owner part), and it is a later pass that reads source and
native type again.

**External symbols.** `contracts/generation/go/output.go:17-21, 37-40`:

```go
type ExternalSymbol struct {
	owner    *identity
	language representation.Language
	key      string
}
// ...
type ExternalUse struct {
	At     semantics.Occurrence
	Symbol ExternalSymbol
}
```

When a rule maps a horos declaration to a type that a library ships, such as
`Cell` exported by the package `@bitspark/cell`, generation records the use as
an external symbol: an issuer-owned key (for TypeScript the module and the
export name) at a site. The sealed output carries these, and TypeScript's
Files stage builds the import from them for each placement. This is the one
place where generation already records a native symbol by identity for a
later stage to realize, which is the shape of option S2 in section 6.5.

## 4. One request traced through today's code

The request: Go types for `People` from the module of section 1.1, with the
six Go units of 3.6 registered. Read from the code and followed by hand; not
run.

1. The planner starts at the declaration `People`. No recipe and no use rule
   is filed under its selector, so the default declaration recipe runs: name
   `People`, form `alias`, a slot for the body. The body is a use of `Page`
   applied to `Person`, so the planner plans the declaration `Page` with the
   same recipe: name `Page`, form `defined`, one parameter `T`, a slot for
   the body.
2. `Page`'s body is a use of the builtin form "closed record". The record
   rule's callback reads two fields, asks for a slot per field type, wraps the
   second in the helper reference and returns
   `record(field(Items, false, slot), field(Next, true, helper(key, slot)))`.
3. The planner fills the first slot by planning the use `list<T>`: the list
   rule's template `slice(slot E)`, whose slot is the use `T`, replaced from
   the generic environment. It fills the second slot by planning the use
   `text`: the constant `primitive(string)`. It resolves the helper reference
   by planning the declaration `BittypeOptional` once.
4. Back at `People`: `Person` is planned like `Page`, and the body of `People`
   becomes a reference to the declaration `Page` with one argument.

The result is a `Program` of four native declarations (`People`, `Page`,
`BittypeOptional`, `Person`) and nine step strings. Eight rule applications
were made (the default recipe three times, the record rule twice, the list
rule once, the text rule twice) and one helper was resolved. No selection
exists as a value afterwards, and the use `Page<Person>` leaves nothing of
its own.

## 5. Fixed points and their standing

Only fixed points are listed: owner rules, accepted designs, and facts
implemented today, which are changeable. Sections 2 and 3 hold the sentences.
"BG" is binding generation and "RF" representation families (9 November
2026, accepted by the owner of research 0014 and 0015); "BC" is the service
library's binding contracts (11 November 2026, accepted by that repository's
owner seat). Our derivations open section 6.

| # | Statement | Standing | Shown in |
| --- | --- | --- | --- |
| 1 | AC-02, 03, 04, 05, 08, 09, 10 | Owner rule | 2.1 |
| 2 | "Properly done", with no refactor later | Owner rule | Question |
| 3 | A contribution is one thing a consumer registers on its own. A caller that needs whatever the consumer registered for a meaning states a demand that the shared mechanism resolves | Owner rule | 2.1 |
| 4 | bittype selects native representations and publishes each selection in a typed plan. One backend per target, owned outside bittype, works from the plan, the source semantics and the contributions. bittype emits no run-time adapters | Accepted design (BG) | 2.2 |
| 5 | The planner owns selection, substitutions, references and the plan. A representation implementation owns native layout, presence rules, construction and access | Accepted design (BG) | 2.2 |
| 6 | The elements of a paired contribution. Facets a binding product needs are selected together with the native facet; a missing facet fails; no structural fallback; duplicates are rejected. "The unit of implementation is library × target" | Accepted design (BG) | 2.2 |
| 7 | A structural profile and an exact mapping use the same public mechanism | Owner rule (AC-04) and accepted design (BG) | 2.1, 2.2 |
| 8 | A family contributes, per target, one contract with seven parts, explicit refusals among them, and is delivered complete | Accepted design (RF) | 2.3 |
| 9 | One immutable representation instance per use; the native type and the plan entry are projections of it; no later pass | Accepted design (RF) | 2.3 |
| 10 | A factory by default; optionally a target-owned lowering facet with typed target expressions; the backend never switches over families; no source-code strings | Accepted design (RF) | 2.3 |
| 11 | The plan is a shared envelope plus family-owned payloads, and a graph: nothing the orchestration must resolve hides in a payload. The envelope lists native symbols and imports | Accepted design (RF) | 2.3 |
| 12 | A payload is an ontos value whose type its family declares in horos; bittype generates that type's Go and TypeScript types | Accepted design (RF) | 2.3 |
| 13 | Five version axes. Generic requirements from the first version. The envelope keeps an unstable version until six vertical cases pass | Accepted design (RF) | 2.3 |
| 14 | bittype's native vocabularies stay types-only; a backend's emitter is scoped separately. The TypeScript backend belongs to the service library, a Go backend to a future Go runtime library | Accepted design (BG) | 2.2 |
| 15 | The kernel imports neither horos nor ontos | Accepted design (bittype's charter and repository layout, which record owner decisions of 31 October and 6 November 2026) | 1.1 |
| 16 | Every factory receives the call context. The memo key is "the selected representation instance (its descriptor and its generic environment) plus the input's identity" (with law K3) | Accepted design (BC). Its law L1 is unconfirmed | 2.5 |
| 17 | The completeness test backend is in bittype, in Go, and is a test. bittype publishes no conversion code | Accepted design (BG and bittype's charter) | 2.2, 2.4 |
| 18 | A contribution states its target and supported profile | Accepted design (BG), and implemented. That a request selects by exactly one profile string is implemented only | 2.2, 3.4 |
| 19 | Generation decides no file, package, module or import. One generated set serves every placement | Accepted design (`docs/source.md`, design decision of 7 November 2026) | 2.6 |
| 20 | Identity, naming intent and realized spelling are three things. In the first slices the mapping spells the name at generation, and each language checks or refuses it | A recommendation of research 0016 that this seat accepted as research owner; design page on the main branch since 13 November 2026, written and decided by this same seat, read twice by an independent reader | 2.6 |
| 21 | Nothing is released. No interface shown here has a consumer outside our repositories | Implemented today | Question |
| 22 | Everything in section 3: a unit is a profile and three lists; a rule returns native intent and declines by an error; the planner keeps a native program and debug strings; no name of a rule survives the process; a replacement names use rules only; one helper exists, wired by a key | Implemented today | 3 |

Out of scope: the design of the Files stage itself; capability lifetime and
the laws L1 to L14; the report of binding support; the choice of native
representation for any particular family.
