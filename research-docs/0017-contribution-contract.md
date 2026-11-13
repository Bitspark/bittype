# Research: The contribution contract of a type generator that must also carry conversion decisions: what a rule returns, requests between contributions, facets and the identities of the plan

**ID:** 0017
**Date:** 12 November 2026
**Status:** delivering
**Run-ID:** run_49f8e87c-d37a-4578-ae05-65eaf57a60f5
**Reviewed:** https://github.com/Bitspark/bittype/issues/196#issuecomment-6839823652
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
**kernel** and a unit that a consumer registers a **contribution**.

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
  families live. Three of its generic types appear in this document:
  `Function<A, B, E>` (a callable with an argument, a result and a declared
  failure), `Stream` (another type of the same kind) and `Service<Ops>` (a
  record `Ops` of operations served at one address). A fourth, `Cell<T>` (a
  holder of a current value), belongs to an independent fixture library and
  is our test case for an exact mapping. A value of one of them inside a
  message is a live reference to something remote, a **capability**. Their
  native types are shipped by the library, not generated.
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
| Source integration, use | The **source integration** is the package that reads horos for the kernel: it classifies declarations and uses and gives rules a typed reader. A **use** is one occurrence of a type expression, reached under one declaration along one path: in `Page`, `list<T>` is one use and `T` inside it another |
| Source identities | Three are distinct. The **selector**: what kind of thing a use is, equal for every `list<...>`; `builtin/7:list-of/0:` for a builtin form (each part length-prefixed: the constructor, an empty discriminator), `declaration/<module pin>/<name>` for a named declaration. The **source-use identity**: which occurrence. The **source identity of a type**: which type; for a type without free parameters its descriptor. For an open type no accepted document says what it is, and section 6.4 depends on that |
| Builtin form, native form | A **builtin form** is a type constructor built into horos (text, list, record). A **native form** is a constructor of native intent, issued by one target's native model, such as Go's `slice`, `record` and `field` |
| Native intent | bittype's description of a native type before any syntax exists: a tree of native forms with string data and arguments. In the code: `Expression[T]`. It describes types only. A **native model** is the package of one target that issues its native forms and checks native intent |
| Kernel, planner | The **kernel** is bittype's shared generation mechanism: the contract package and the registry. Quoted documents and code comments also say "shared mechanism", "orchestration" and "composer". Its **planner** walks the uses, selects a rule for each and assembles the result |
| Rule | A function or template that maps the uses of one selector to native intent. The code and quoted documents call a rule for uses a "binding". In our own words "binding" always means value conversion, so we say **use rule**, and **declaration recipe** for the rule that turns a named horos declaration into a native declaration |
| Helper | A native declaration that is not generated from a horos declaration and that a rule refers to through a key. Go's `BittypeOptional` is the only one. The accepted design also says "helpers" for generated value-level functions; we call those lowering output |
| Family | A representation family: one kind of horos type together with its native representation (text, lists, closed records, optional fields, integers, unions). A **structural family** is one contributed under the structural profile |
| Profile | A named promise about what the native types retain. Its identity is its name string. The only implemented one is `structural-representation/1`, the **structural profile**: native types that mirror the horos structure and represent every value of the horos type. A second, `structural-enforcement/1`, is declared and stated by no rule |
| Exact mapping | A library's own representation of one named horos declaration, for example of its `Cell<T>` as a class it ships. Selected in place of the structural families for that declaration |
| Library | In our words and in quoted designs: a separately owned code base that owns horos declarations and their native types. In the code (`MappingLibrary`, `Library()`, `library-3`) it is one registered unit |
| Contribution, registered unit | A **contribution** is one thing a consumer registers on its own (the owner's definition, section 2.1). What one contribution comprises is open (6.1). A **registered unit** is what exists today: a "mapping library" holding one rule; quoted documents call it a "mapping" |
| Contract | Without qualification: the contribution contract, between kernel and contribution. Other uses are qualified: the "representation contract" (what a family must define, 2.3), the service library's "binding contracts" (2.5), the Go "contract package"; quoted "contracts" can also mean service types such as `Function` |
| Facet | A part of a contribution that one product reads: a native facet (the type), a value adapter facet (the conversion; quotes also say "adapter facet"), a root facet, an operation facet, a lowering facet |
| Product | One kind of output requested from the same selections: native type declarations, or a **binding product** such as value adapters, JSON presentation or operation data. The accepted designs use the word without defining it; this is our reading |
| Selection, instance | A **selection** is the planner's choice of one rule for one use. An **instance** is the accepted design's "representation instance": the immutable result of one selection for one use. Quoted "generic instance" means an application such as `Page<Person>` |
| Plan, plan entry | The **plan** is the "selected representation plan": the typed record of every selection of one request, published for backends. It does not exist yet; the Go type `GenerationPlan` is something else. A **plan entry** is one item of it; whether there is one per use or per distinct selection is open (6.4) |
| Envelope, payload | The **envelope** is the part of the plan whose shape all families share: identities, links, versions. A **payload** is the family-owned part of an entry: an ontos value, typed by a horos type that the family declares |
| Decision record | Our word for what a rule returns for one use under the shape we would build. It is not a horos record type: "the record rule", "the record family" and "Record" in quotes always mean the horos record type |
| Slot | A declared place in a rule's result that stands for the representation of another use. The planner fills it |
| Demand | The owner's word (2.1): a request from a contribution that the kernel resolves to whatever the consumer registered. Today there are two kinds: "represent this child use" and "the helper under this key" |
| Generic environment | The mapping from a generic declaration's parameters (its **binders**) to the actual arguments of one application: `{T -> Person}` for `Page<Person>`. An **adapter argument** is the adapter selected for an actual argument |
| Native symbol | The accepted design's word for a reference to a named native entity that generated code must spell and import: a generated declaration, a helper, an external type, a factory. The design does not say whether it is an identity or a spelled name with an import path (6.5) |
| Backend, factory | A **backend** is a program, one per target and owned outside bittype, that reads the plan and emits the code that composes conversions. A **factory** is a run-time function that builds an adapter from a payload and the adapters of the slots |
| Completeness test backend, conformance vectors | The first is a backend written in Go inside bittype's tests: the acceptance instrument of the plan, never shipped. The second are test vectors that bittype publishes per family; every implementation of the family's conversion, in any repository, must pass them |
| Source semantics | The checked semantic graph of the horos modules. A backend reads it beside the plan. No document we know gives it a serialized form |
| Files stage, placement | The stage after generation that puts generated declarations into files, packages or modules according to a **placement** the consumer supplies (2.6) |

### 1.3 The stages

A consumer's composition registers one source integration, one target builder
per target, and the registered units. For each use reachable from the
requested declarations the kernel selects one rule by the use's selector,
runs it, and resolves the rule's slots by planning the uses they name. The
result is native intent, which the target builder turns into a syntax tree
and, for TypeScript today, a Files stage puts into files (2.6). That exists
(section 3). The accepted design adds a second result of the same walk, the
plan, read by a backend outside bittype. The contract in question is what a
registered unit is, what its rule returns and what it may ask of the kernel.

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

**Backends and process boundaries.** The TypeScript backend runs in Node and
reads a serialized plan. The owner of a Go backend, a Go runtime library,
does not exist yet; bittype owns only the completeness test backend. Kernel,
registered units and that test backend share one Go process, where an
identity can be a pointer. The TypeScript backend and factories do not.

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

**What violates AC-05 and AC-10:**

> Openness is achieved by accepting arbitrary payloads, unchecked claims or callbacks that discard the guarantees the mechanism promises.

> Delivering a first change whose input or identity later has to change.

The first sentence matters for the question on facets: an open set must not
become "arbitrary payloads". The second is the reason for this consultation.

**Wrong coupling, and the demand that corrects it** (owner request of
11 November 2026, section "Package boundaries and forbidden dependencies").
A table distinguishes two faults. The second, with its cells joined into one
paragraph:

> Wrong coupling. What it is: A contribution names another contribution that the consumer is meant to choose. What is wrong: The dependency. The caller makes the shared mechanism's choice. Consequence: The mechanism's selection, scope and failure guarantees do not apply to that call. Replacing the callee means editing the caller. Correction: A demand that the shared mechanism resolves.

And below it:

> Ask of each call from one part to another whose answer the caller needs: its own, or the composition's. A caller that relies on its own invariant, a shared mechanism or a vocabulary calls directly. A caller that needs whatever the consumer registered for a meaning states a demand and lets the shared mechanism resolve it.

The same section defines the unit (two list items joined):

> A contribution, meaning one thing a consumer registers or selects on its own, is one unit. It depends on shared contracts, the vocabularies of its own scope and its own private parts. Only a composition depends on a contribution.

**Evidence for an extension boundary:** "implement a materially different
contribution using its public contracts in a separately owned package or
fixture. Keep the shared mechanism fixed."

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
change is designed". We no longer hold the second half: this document exists
because several shapes of that change are given nowhere. No import rule can
see this kind of decision: "Record building `std::optional` imports nothing
forbidden."

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
and further products such as JSON presentation and operation data. "That set
is **binding generation**."

**The decision** (list items joined, links reduced to issue numbers):

> One semantic owner does not require one emitter. It requires an authoritative decision, an explicit interface carrying that decision, and tests that establish that its consumers preserve the promised meaning.
>
> **bittype** remains the authority for selecting native representations. It publishes each selection as a typed **selected representation plan** (#195).
>
> **A binding backend per target** consumes three inputs: the plan; the source semantics; each library's **paired contribution** (#196). From these it emits static composition and typed glue.

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
requirements". ("Generic instance" is an application such as
`Page<Person>`.) Its form is "typed references and target-owned expressions", and
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

**The laws every contribution must satisfy.** The page states two round-trip
laws, decoding after encoding up to a declared equivalence of native values
and encoding after decoding exactly, and eight further laws. They matter here
for one reason: what a family records must contain whatever these laws
depend on. Two examples from the first slice: in Go, `Value` is ignored when
`Present` is false, so the equivalence identifies two absent values with
different `Value`; in TypeScript, an absent property and a present
`undefined` are different native states. Neither fact is visible in the
native type.

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

> A family is delivered with its native type, its plan entries, its value conversion and its qualification together. It is never a types-only mapping that a later binding step reopens.
>
> **The promise.** Adding a family through the published identity, dependency, invocation, plan-graph and target-extension contracts changes only that family and the consumers that compose it. Existing representation identities keep their meaning.

The page defines none of these five contracts further. In our reading they
are the identities of question 2, the demand of question 1, the way a backend
calls a factory, the envelope, and the facets of question 3.

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

> “Defined once” should mean that representation decisions have one authoritative origin. It need not mean that a Go planner and a TypeScript runtime somehow share one executable implementation. Their agreement still needs a shared contract and verification.

**Factories and the lowering facet:**

> By default a plan entry names the family's run-time factory and gives it the family's payload as arguments. A family may also supply a target-owned **lowering facet** that turns its recorded decisions into typed target expressions or helpers, such as generated accessors or constructors for Go. The binding backend owns imports, names, placement and assembly. It never switches over families, and source-code strings never stand in for typed target expressions. A shared interpreter is a valid factory implementation only if its handlers are registered contributions executing selected decisions. Generated source is never modified as an extension mechanism.

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

**Issue #207** is the six families (items 1 and 5 of eleven):

> 1. **Contract.** Each family is one representation contract with all seven parts, contributed through the public interface of #196. The kernel, the planner and the completeness test backend gain no case for a family.
> 5. **Existing output unchanged.** The native declarations generated for today's pagination fixtures stay byte-identical in all four targets. Existing binding identities keep their meaning. The profile's own identity is settled first, in #216.

("Binding identities" are the rule IDs such as `list/golang/1` of section
3.4.) Its verification: "`Page<Person>` goes from horos source to native
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

**A caution about our own reading.** On 12 November 2026 two design
questions in this line of work were answered from the code, wrongly. This
document may make the opposite error and read more into the accepted pages
than they say.

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
run-time requirement, each shown by its own test". Its structural run-time
package
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
"is either a namespace, which is a field the plan selects as the structural
record family, or a leaf whose selected contribution supplies an operation
facet". So one contribution's facet must recognise what the plan selected for
another use: a second reader of a contribution's identity besides the
backend.

**"Instance" in two senses.** Our pull request #260 brings bittype's two
design pages into agreement with these contracts; it is open. Its second
reading (12 November 2026, 18:50Z; line numbers are those of its branch)
found:

> `docs/representation-families.md:44-45` says "one immutable representation instance per use", and `:259` keys the memo by "the selected representation instance". The page now uses "instance" in two senses.

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
selected mapping and the target builder build declaration bodies and bind
each named use to its intended declaration. The consumer must "Supply one
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
2026, 22:34Z and 23:05Z). An intent keeps the chosen name and the
qualification apart, and "A consumer learns the realized name and the
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
| 3.4 | Four names exist for selection. Only the source's selector means anything outside the process | 2 |
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
| `Desc` | Abbreviates `g.MappingDescription[h.Language, s.Language, source.Reader]` |
| `Rule` | Abbreviates `g.BindingDescription[h.Language, s.Language, source.Reader]` |

An `ast.NodeRef` is a site in a syntax tree, used for diagnostics, and
`g.At(e, site)` attaches one. `from` and `to` are language identities: a
language with a syntax revision that its owner calls a "profile", unrelated
to the representation profile. Code comments use words that our prose has
replaced: "provider" (a registered unit, or in a file name the target
builder), "Run" (executing the plan closure), "composer" (the kernel),
"occurrence" (a use), "binding" (a use rule; in `Generic.Bindings` a
parameter binding) and "capture" (binding a slot to a use).

### 3.2 What a rule can say: native intent

`contracts/generation/go/mapping.go:5-11`, `expression.go:9-14, 30-35,
61-74` (`HelperKey` has the same two fields as `Slot`):

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

An expression has one of six kinds: a native form with string data and
argument expressions (`g.Native`); a slot (`g.Mapped`, "captured"), which the
planner replaces with the native intent it plans for the use bound to it; a
generic parameter of the native declaration being described; an actual
argument at an application; a reference, through a key, to a helper
declaration that some registered unit must supply (`g.Helper`); and a
reference to a planned native declaration, which only the planner creates.

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

`contracts/generation/go/description.go:18-52, 61-84` (a `Capture[F]` pairs a
slot with a use):

```go
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
`Registry.Replace`, naming the selector and two use rules, the expected one
and its replacement, each by handle and rule ID. The registry removes the
expected rule before planning. Only use rules can be replaced. A declaration
recipe cannot, and two units that supply a helper under one key fail as
ambiguous.

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
name, the flag `Optional` and one child use, the field's type
(`transform/syntax/horos/source/go/read.go:25-30`). Only the field's type is
a use that the kernel walks and selects for. The flag is content of the
record form. No selector exists for an optional field, so nothing can be
registered for it today.

**Go.** The record rule is a use rule with the ID `record/golang/1` for the
builtin form `record` with discriminator `closed`. Its constructor takes a
helper key, `Library(optional g.HelperKey)`. Its callback refuses a record
that is not the whole body of a declaration, by returning the located
diagnostic `src.Unsupported(use.Use, "nested records require a named
declaration in this profile")`. It then reads the fields and runs this loop
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
(`transform/syntax/horos-to-golang/mapping/optional/go/optional.go`). Its
description holds one helper recipe under the key it was given: the native
declaration `BittypeOptional`, of the form `defined`, with one parameter `T`
and a record body of two fields, `Present` of type `bool` and `Value` of type
`T`.

Neither package imports the other. The consumer's composition creates the key
and gives it to both (`integration/go/libraries_golang_test.go:26, 31, 33`):

```go
	optionalKey := g.NewHelperKey("horos/optional")
	// ...
		record.Library(optionalKey),
	// ...
		optional.Library(optionalKey),
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

Later families raise questions that we cannot yet classify as the same or as
different, and we have not checked each against horos. As far as we can tell,
a nullable value is a type in horos (a union with null), so a use with a
selector, selected like a list. A map's requirement on its key and "The null
representation cannot collide with a valid `T`" are requirements on what was
selected for an operand. A union's discriminator is a fact about the union's
branches taken together, neither a use nor a flag.

### 3.7 What the planner keeps per use, and what it drops

Two of the planner's output types (`contracts/generation/go/planning.go:33-43`):

```go
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

A third, `Program[T]`, is the closed list of native declarations; only the
target builder receives it. `GenerationPlan` is a closure over it plus the
`Step` list, which exists for inspection in tests. `Output` is the sealed
result of running the closure: the checked syntax trees, the sites of the
declarations the request exposed, and the external uses of 3.8.

After applying a use rule the planner resolves the result: it replaces every
slot with the result of planning the use bound to it, and every helper
reference with a reference to the helper's declaration. It then appends one
step (`planning.go:564`):

```go
	p.steps = append(p.steps, Step{Binding: b.name, Source: key, Form: out.form.Name(), Expression: describe(out)})
```

`b.name` is the prefixed rule name, `key` the selector key, and `describe`
prints an expression as a string. A finished declaration is stored
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
"Present", "Value")`). This has the shape of option F3 in
section 6.3 (a per-family part registered separately) and of I2 in 6.4 (an
identity with an owner part), and it is a later pass that reads source and
native type again.

**External symbols** (`contracts/generation/go/output.go:14-41`). When a rule
maps a horos declaration to a type that a library ships, such as `Cell`
exported by the package `@bitspark/cell`, generation records the use as an
`ExternalUse`: a site and an `ExternalSymbol`, which is an issuer-owned key
(for TypeScript the module and the export name). The sealed output carries
these, and TypeScript's Files stage builds the import from them for each
placement. This is the one
place where generation already records a native symbol by identity for a
later stage to realize, which is the shape of option S2 in section 6.5.

## 4. One request traced through today's code

The request: Go types for `People` from the module of section 1.1, with the
six Go units of 3.6 registered. Read from the code and followed by hand; not
run.

1. The planner starts at the declaration `People`. No recipe and no use rule
   is filed under its selector, so the default declaration recipe runs: name
   `People`, form `alias`, a slot for the body. The body uses `Page`, which
   is planned with the same recipe: form `defined`, one parameter `T`.
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
implemented today, which are changeable. "BG" is binding generation and "RF"
representation families (9 November 2026, accepted by the owner of research
0014 and 0015); "BC" is the service library's binding contracts (11 November
2026, accepted by that repository's owner seat). Our derivations open
section 6.

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
| 22 | Everything in section 3: a unit is a profile and three lists; a rule returns native intent and declines by an error; the planner keeps a native program and debug strings; a rule's only names are a string local to its unit and a registration counter; a replacement names use rules only; one helper exists, wired by a key | Implemented today | 3 |

Out of scope: the design of the Files stage itself; capability lifetime and
the laws L1 to L14; the report of binding support; the choice of native
representation for any particular family.

## 6. What we have considered

Everything in this section is ours and open. "We lean to" marks what we would
do today. No option is ruled out because the present code would have to
change. Option labels are unique across the section, so that an answer can
cite them.

**Our derivations, all open.**

- The contract "changes once", in the order of section 1.4 (from AC-10; 2.4
  gives what speaks against the order).
- The record rules deciding the optional representation is a case of wrong
  coupling (a finding of this seat; 2.1, 6.2).
- The shape we would build: the unit, the rule's result, and a kernel that
  reads no meaning of what it carries (6.1).
- Fixed points 12 and 15 together: the plan's types live with the planner,
  which may import neither horos nor ontos, so the planner holds a payload as
  something opaque, and encoding and validation happen outside the kernel.
- The payload types of the first slice can use only the six existing
  families, because bittype generates their Go and TypeScript types. A
  payload cannot hold an integer, such as a field index, until the integer
  family exists.
- The Go conversions of the first slice are test code (6.3).
- What a naming intent records, and how a plan refers to a native
  declaration (2.6, 6.5).

**What would be final, what stays unstable, what is versioned.** Final after
the contract change, as we intend it: the unit, what a rule returns and how
it declines a use, the request capability, how facets are carried, and the
contribution identity. Unstable until the six vertical cases pass (fixed
point 13): the envelope, that is the shape of an entry, the links and every
envelope field. Versioned per family: payload schemas.

**Two contradictions between fixed points that we see and have not
resolved.** First, AC-10 requires of a first change that "its inputs and
identities have their intended shape", while fixed point 13 keeps the
envelope explicitly unstable, and the plan's identities sit in the envelope.
Either the identities belong to the contract and are final, or AC-10 does not
bind an artifact that is declared unstable. Second, fixed points 4 and 11
against 19: a plan that a backend needs nothing beside cannot hold an import
of a declaration whose place is decided later (6.5).

### 6.1 The shape we would build: the unit and what a rule returns

**The unit.** The owner's definition is "one thing a consumer registers or
selects on its own". Binding generation says "library × target";
representation families says a family contributes "per target".

| Option | One contribution is | Where the default declaration recipe sits | Where "generics" sits |
| --- | --- | --- | --- |
| U1 | One horos declaration or builtin form, one target, one profile, with all its facets | It is neither: one contribution selected for every named declaration without an exact mapping, a special case of the unit | Split. "Record or alias" is that contribution's decision; binders, arguments and substitution are the planner's. The family has no contribution of its own, against issue #207's "one representation contract with all seven parts" |
| U2 | Everything one owner publishes for one target: the structural profile's Go rules are one contribution with parts | A part of the structural library | A part of it. Families stay the unit of delivery and of the plan's identities, not of registration; whether a consumer can still replace one part is then a second question |
| U3 | One family for all its targets: one payload schema and one set of vectors with a native rule per target | Per target inside it | As in U1. Cross-target agreement has one owner, but adding a target edits every family |

We would build U1 and hold it weakly against U2. Under U2 the request of
section 6.2 is a matter inside one library, and option D1 there is the
natural reading.

**What a rule returns.** Fixed point 9 gives the property; issue #196 item 6
"prescribes no signature".

| Option | What a rule returns for one use | What it guarantees | Weakness |
| --- | --- | --- | --- |
| R1. Two results | Native intent and a payload, from one call | No second pass: both come from one call | Each decision is written twice. Agreement of the two rests on each author, or on a check that compares them |
| **R2. One decision record and its projection** | The family's decisions for this use as one value, with the slots it needs. The contribution also supplies a function from (decision record, what the kernel filled into the slots) to native intent. The plan entry is the record plus the planner's links | The native type is computed from what the plan carries, so each decision is written once | The native rule must be a function of the record. Every decision that is control flow becomes a field. The projection gets no source reader |
| R3. Native intent is the record | A richer native intent only: typed, serializable, with annotations for what no type shows. The plan entry is derived from it | One value and one vocabulary per target | Fixed point 12 asks for a family-owned schema in horos, and this is a target-owned one. Decisions that show in no type ("`Value` is ignored when absent") need annotations, a second vocabulary inside the first |

An illustration for Go's record family and the body of `Page<T>`. It is our
derivation, and no signature is proposed:

```text
R1  rule(use) returns two values:
      native  = record(field "Items" of s1, field "Next" of s2)
      payload = {fields: [{source: "items", native: "Items"}, {source: "next", native: "Next", optional: true}]}
      slots   = s1: the use list<T>;  s2: the use text, already wrapped as the rule saw fit
    The names and the fact "next is optional" are written in both values.

R2  decide(use) returns one value and its requests:
      record   = {fields: [{source: "items", native: "Items", slot: s1}, {source: "next", native: "Next", slot: s2}]}
      requests = s1: represent the use list<T>
                 s2: represent "an optional field" with the operand use text       (the request of 6.2)
    project(record, filled) returns, without reading the source:
      record(field record.fields[0].native of filled[s1], filled[s2])              (filled[s2] is a member here)

Entries for Page<Person>, with open nodes (O1 in 6.4):
  e1  record     payload = the record above       links: s1 -> e2, s2 -> e3
  e2  list       payload = {}                     link:  E -> binder T of Page
  e3  optional   payload = {valueIgnoredWhenAbsent: true}   link: operand -> e4; owns the declaration BittypeOptional
  e4  text       payload = {refusesInvalidUTF8: true}
  e5  Page       payload = {form: "defined"}      binder T; body -> e1
  e6  Page<Person>   refers to e5 with T := e7 (Person): native argument Person, adapter argument e7
```

What the illustration leaves open, and we have not decided:

- **Which value the projection reads.** A typed Go value of the family's
  generated payload type, which is encoded as an ontos value afterwards and
  outside the kernel, or the encoded payload. Under the first, "one record"
  depends on a step that the kernel does not own.
- **How a rule declines a use.** Fixed point 8 lists "explicit refusals"
  under applicability, and AC-10 requires that a first change "explicitly
  refuses the cases it cannot yet check". The classes we can tell apart: not
  yet supported, refused by design, a failed requirement on an operand, and
  the kernel's own (no contribution, ambiguous, missing facet).
- **Where a native declaration comes from** that a contribution owns without
  a horos declaration behind it, such as `BittypeOptional`.

**Why the kernel contract may change at all** (AC-08). The planner cannot
keep one origin per representation decision while a rule returns native
intent only and the selection is discarded. Only the planner sees each
selection exactly once. The result shape is the same for a structural family
and for a library's exact mapping.

**What speaks against R2.**

- It buys agreement for the native type only. The conversion is a factory in
  another repository, or a handler in a test, that reads the record and must
  agree with a native type derived from it. The record makes disagreement
  detectable by vectors and not impossible, so R1 with a check may buy the
  same.
- The projection is contribution code that the kernel cannot inspect. What it
  decides in its body ("optional means wrap in `BittypeOptional`") is decided
  off the record. The line between a decision and a computation is drawn by
  each author, which is the weakness charged to R1.
- A family whose decision depends on what was selected for its operand (a
  nullable value that is `T | null` when `T` excludes null and a wrapper
  otherwise) needs facts about the operand's selection as inputs. R2 does not
  say where they come from.
- Diagnostic sites, and the unchanged output we use as a regression oracle,
  push syntax-level data into the record, toward R3.

### 6.2 How a record's optional field gets its representation

**Two readings of whether the record rule must ask at all.**

- *Presence is the record's own layout.* The accepted owners table gives the
  representation implementation "Native layout, presence rules", and the
  payload column lists "Record field layout and presence behaviour" as one
  item (2.2, 2.3). On this reading a consumer who wants another presence
  representation legitimately replaces the record contribution.
- *Presence is another contribution's meaning.* Representation families
  counts optional fields among the families, each of which is "delivered
  with its native type, its plan entries, its value conversion and its
  qualification together", and the owner's rule asks for a demand wherever a
  caller "needs whatever the consumer registered for a meaning". On this
  reading the record rule's
  inline decision is the coupling of section 2.1. This is our finding, and we
  may be wrong about it.

**What no document says**, if the record rule asks: what the request is
keyed by, since the kernel selects by the selector of a use and an optional
field is no use (3.6); what the answer is; and how it is replaced.

| Option | The request | What it guarantees, and to whom | Weakness |
| --- | --- | --- | --- |
| D1. Presence stays in the record family | None. The record's payload holds the presence representation; Go's helper becomes a private part of Go's record contribution | One contribution and the simplest plan. No new kernel capability | A consumer who wants another optional representation replaces the record contribution. "Optional fields" is then a part of the record family and no contribution |
| D2. The source integration presents the optional field as a use | The horos integration gives the kernel a use for "optional field of type X" with its own selector and the field's type as its child. The record rule asks for it as it asks for any field's type | Selection, ambiguity, replacement and the missing-rule failure apply unchanged. The kernel gains no capability | The source integration presents as a use something that horos does not treat as a type. Every later form content of this kind (a union's discriminator) needs the same decision in the integration |
| **D3. A source-owned selector for a form that is not a use** | The source integration publishes a selector for "an optional field of a record". The record contribution asks the kernel for the representation under that selector and passes the field's type use as the operand | As D2 for the consumer. The record's entry links to an optional entry by a slot. The record family never learns what was selected | The kernel gains a capability: selection at the request of a contribution for something that is no occurrence in the source graph. Its general justification (AC-08) must be written |
| D4. A key that the composition wires | The consumer's composition creates a key and gives it to the asker and to the supplier | The kernel resolves it, and a missing supplier fails | Selection is by wiring and not by an identity that the source owns. A key that must also be recorded in a plan has to be a value, and is then a selector by another name |
| D5. A role that the asking family owns | The record family publishes a role ("presence of a field"), and contributions register as its providers for a target | The asker defines exactly what it needs from an answer, including its shape | A second selection mechanism beside selectors. Two families that need presence (records, and tuples with optional elements) would each own a role |

**The form of the answer** matters under D2 to D5:

- **Type.** Native intent with one slot for the operand: `BittypeOptional[·]`,
  `std::optional<·>`. TypeScript then has no optional contribution, its `?`
  stays the record family's layout decision, and the six families mean
  different sets of contributions per target.
- **Member.** Given the field's name and the operand's representation, the
  answer is the member: in TypeScript a property with a modifier, in Go a
  field with a wrapped type and the helper declaration it needs. All four
  targets ask alike. A slot then holds "native intent for a position that the
  asking family defines", not "for a type position".
- **Record-only.** The answer is the optional family's decision record, and
  the record family's projection reads a small published vocabulary of it and
  builds the member itself. One family then interprets another's payload.

**We lean to D3 with Member.** "Optional fields" is then one contribution per
target. Go's `BittypeOptional` becomes a part of Go's optional contribution,
and the consumer-made key disappears. An optional nullable field composes
without touching the record family.

**What speaks against the leaning.**

- The first reading above: D1 may simply be right, and then the finding, step
  2 of the order and this whole capability fall away.
- D2 reaches the same guarantees with no kernel capability, at the price of a
  decision in the source integration that D3 also needs (a selector).
- Member makes the answer's shape depend on the asker: a record wants a
  member, a tuple with an optional element something else. And under Member
  the optional contribution must know the record family's member vocabulary,
  the coupling charged to D5.
- A request that returns a link and no facts cannot serve a family whose
  decision depends on its operand's selection (6.1). If requests must return
  declared facts about what was selected, that is a larger capability.

### 6.3 Facets, and the lowering facet

**What is decided** (fixed points 6, 10, 14). A contribution pairs a native
facet with a value adapter facet; a root facet and "optional product facets"
exist; a lowering facet may turn recorded decisions into "typed target
expressions or helpers"; the backend never switches over families; bittype's
vocabularies stay types-only.

**What no document says.** Whether the set of facet kinds is closed in the
kernel's contract, and in which vocabulary a lowering facet for Go writes an
accessor or a constructor.

| Option | The contribution contract | What it guarantees, and to whom | Weakness |
| --- | --- | --- | --- |
| F1. A closed list | One field per facet kind: native, value adapter, root, operation, lowering | Every reader sees every facet's type, and the kernel can check each | A new binding product needs a kernel edit. That violates AC-03 if binding products are inside "the agreed extension domain", which no document says. The kernel would name `operation`, a concept of the service library |
| **F2. An open set, keyed by the reading product** | The kernel knows the native facet and carries the decision record. Every other facet is carried under a kind that the product reading it issues, as a native model issues its forms. The kernel checks presence and never reads a facet's content | A new product adds a kind and edits nothing shared. The failure "missing facet" stays the kernel's | The kernel carries values it cannot check. Our candidate guarantees, which keep this from being "arbitrary payloads": each kind has an issuer that validates its facets, and a facet refers to slots only through kernel-issued references |
| F3. No facets in the contribution | The contribution is the native facet and the decision record. Everything else is registered separately, with the product that reads it, under the contribution's identity | The smallest kernel contract, and one mechanism on both sides of the process boundary | "Selected together" (fixed point 6) becomes a lookup by identity, and a second owner could attach a facet to someone else's family |

The Go conversion of the first slice is a second choice. Under the order of
section 1.4 the contract is final at step 1, before any conversion exists; so
none of these tests the contract before it is fixed unless it is done on
paper or in a prototype (6.6).

| Option | How the Go test backend converts `Page[Person]` | What it shows | What it leaves open |
| --- | --- | --- | --- |
| V1. Registered handlers that execute the payload | One Go conversion per family, registered with the family's contribution, takes the payload and the conversions of its slots. A named struct is read and built by the field names in the payload, through reflection | Completeness: every name, order and presence rule comes from the payload. The design allows it: "its handlers are registered contributions executing selected decisions" | The lowering facet is not exercised. Nothing shows that typed construction can be expressed from the decision record |
| V2. Lowering to generated Go | Each Go family supplies a lowering facet; the test backend assembles Go source, compiles and runs it | Also the lowering facet and vertical case 4 | Needs typed expressions for value-level Go, which bittype may not own (fixed point 14) and whose owner does not exist |
| V3. A paper lowering | V1 in code, plus the lowering facets of Go's record and list families written out on paper against the contract | That the decision record and the slots carry what a lowering needs, as far as paper shows | No compiler has checked it |

**We lean to F2 and to V3**: handlers in code, the Go lowering on paper
before the contract is fixed, and vertical case 4 built later.

**What speaks against the leaning.**

- A lowering facet is a function of the decision record and of typed
  references to native symbols and slots. If it needs something they do not
  carry (a typed reference to a generated member instead of its name, a
  slot's native type as a value-level type argument, the constructor of a
  declaration the contribution owns), the contribution contract changes a
  second time. Reflection in V1 hides exactly these needs.
- Under F2 "selected together" guarantees that the facets were registered
  together, not that they agree.
- If binding products are outside the extension domain, F1 is admissible and
  keeps the checks that F2 gives up.
- The process boundary forces lookup by identity for TypeScript anyway, so
  F3, one mechanism for all facets, may be more honest than "together" in Go
  and "by identity" elsewhere.

### 6.4 The identities of the plan

Five identities are in play. We use these labels from here on.

| Label | What it identifies | What the fixed points say (numbers of section 5) |
| --- | --- | --- |
| **Contribution identity** | Which contribution was selected. A backend in another process finds the factory by it, and the root facet of `Service` recognises a family by it | The envelope has "Selected contribution and profile identities"; duplicates are rejected (6, 11). No form is given |
| **Use entry** | One use and what was selected for it | "one immutable representation instance per use" (9); the envelope has "source-use identities" |
| **Representation node** | One distinct selection, shared by every use that makes it | Nothing directly |
| **Run-time identity** | The part of a memo key that is not the input | "the selected representation instance (its descriptor and its generic environment)" (16) |
| **Native declaration identity** | Which generated native declaration | "It is never a spelling" (20). Section 6.5 |

**The contribution identity.**

| Option | Form | What it guarantees, and to whom | Weakness |
| --- | --- | --- | --- |
| **I1. Owner-stated name and version** | A string the contribution states, such as `list/golang/1`, unique within a registry | To a backend: a stable key in every process | Uniqueness holds per registry only, which is the one place that does not need it: two owners can publish the same name. The version in the name folds one of the five axes into an identity. The name is per target, while the root facet of `Service` must recognise a family |
| I2. Owner namespace, name and version as separate parts | The namespace is an identity the owner demonstrably controls: for an exact mapping the pin of the module it maps, for a profile's families the profile's identity | Two owners cannot collide by accident, and the version stays an axis | Needs a rule for a contribution that maps someone else's declaration (a Go mapping of another library's type), and for two contributions of one owner for one declaration |
| I3. Derived | A digest over the contribution's payload schema, target, profile and selector | Nobody states it, so nobody states it wrongly | A factory can change meaning without changing any hashed part, which is why the design has a semantic-identity axis |

**Use entries and representation nodes.**

| Option | The plan has | Weakness |
| --- | --- | --- |
| N1. Use entries only | One entry per use. Two fields of type `Page<Person>` are two entries | Two uses of one representation have no common identity, which our own criterion on issue #195 (2.5) asks for, unless each entry also carries a run-time identity |
| **N2. Representation nodes; uses refer to them** | One node per distinct (contribution identity, profile, source identity of the type, generic environment). A use holds a reference to a node | It rewords accepted fixed point 9, "one immutable representation instance per use", as "each use refers to one instance". Per-use facts (provenance, the source-use identity) need a second record. "Source identity of the type" is undefined for an open type, and for an optional field, which is no type |
| N3. Both, under two names | Use entries with provenance and a node reference, and representation nodes as in N2 | Two concepts where the accepted text has one word |

**Open or instantiated nodes for a generic declaration.**

| Option | For `Page<T>` and the application `Page<Person>` | What speaks for it | Weakness |
| --- | --- | --- | --- |
| O1. Open nodes | The body of `Page` is planned once, with a link to the binder `T`. Each application supplies a native argument and an adapter argument | Three accepted sentences read this way: the envelope row "Scoped binders and actual-argument environments"; "Each generic use carries its scoped binder identities, native type arguments, adapter arguments and requirements"; and "`T`'s selected adapter travels with its native type". A generic declaration has a complete plan on its own | The run-time identity of "the list inside `Page<Person>`" is then a pair (node, environment) that exists only at run time, so no entry of the plan is the run-time identity. A rule runs once per open use and cannot depend on the actual argument |
| O2. Instantiated nodes | One set of nodes per application: the plan is monomorphised | A node is the run-time identity, and a decision may depend on the actual argument, as a nullable `T` needs | The plan's size depends on the applications in the program, and the native declaration stays generic while its plan does not. A generic declaration alone has no complete plan |

**We lean to I1 and N2, and have no leaning between O1 and O2.**

**What speaks against the leaning.** On I1: see its weakness; the first wrong
factory selected by a colliding name fails at run time, in another
repository's tests at best. On N2: besides rewording an accepted sentence, it
presupposes O2, since under O1 a node is not the run-time identity either.
We noticed this only while writing this document.

### 6.5 How a plan refers to a generated native declaration

**What is fixed.** The plan carries "native symbols" and imports, and a
backend needs nothing beside the plan, the source semantics and the
contributions (fixed points 4, 11). Generation decides no place (19). In the
first slices the spelled name exists at generation and does not change
between placements (20). So the open part is the place, hence the import
path, and the identity by which an entry refers to a declaration: the naming
page rules out the sealed output's handle as "the key of a serialized plan".

| Option | What a plan entry holds | What it guarantees, and to whom | Weakness |
| --- | --- | --- | --- |
| S1. The plan is produced after placement | Spelled names and import paths | The accepted sentences hold as written: one artifact, three inputs | The plan depends on a placement, and one generated set serves several. The planner cannot issue it alone |
| **S2. The plan carries an identity and the spelled name; the backend joins it with the placement result** | A native declaration identity that the planner issues, and the spelled name. The backend reads the place from the Files formatter's result | The plan is independent of placement, and the planner stays its only issuer | The placed result becomes a fourth input of a backend, which amends an accepted sentence. A plan of one run can be joined with the placement result of another. The planner's identity and the output's handle need a correspondence that nothing defines |
| S3. The plan carries naming intents; the backend derives the rest | Chosen name, module identity and policy | Three inputs, and no dependence on Files | The backend applies a published policy a second time, and cannot derive a place at all without the placement |
| S4. Files seals plan and realization into one artifact per placement | The Files stage takes the plan and emits, per placement, a copy with names and places filled in | A backend keeps three inputs and cannot mismatch plan and placement | The Files stage of each target must read and rewrite the plan, so the plan's envelope becomes an interface of four more components |

**We lean to S2.** Member names, which generation fixes, stay in payloads.

**What speaks against the leaning.** The join: for the TypeScript backend the
placement result must cross the same process boundary as the plan, so it
needs a serialized form and a checkable association with the plan, and
neither exists. "The plan is incomplete" gets a second meaning: the plan is
complete and the join failed. A backend also places its own glue (it "owns
imports, names, placement and assembly"), so it may be a participant in
placement more than a reader of its result, which S4 models and S2 does not.
One sequencing fact, which is no design objection: Go and Rust have no Files
stage yet.

**The contract we lean to, as one sketch.** It collects the leanings of 6.1
to 6.5 (U1, R2, D3 with Member, F2, I1, N2, S2) and proposes no signature.

```text
A contribution states    its contribution identity; the selector it accepts; its target; its profile;
                         the horos type of its decision record, with a version; the uses it declines, by class;
                         a decide function and a projection; further facets, each under a kind its reader issues.
decide receives          the kernel's view of one use; the source reader; a way to state requests.
decide returns           a decision record, or a refusal with its class;
                         requests, one per slot: "represent this child use", or
                         "represent this source-owned form with this operand use".
project receives         the decision record and what the kernel filled into each slot. No source reader.
project returns          native intent for the position (a type, or a member), and the native
                         declarations that the contribution owns, each with a naming intent.
The planner issues       the instance for each use: identity, filled slots, immutability;
                         the plan's entries and links; binders, native arguments and adapter arguments;
                         the identity of each native declaration.
The planner checks       selection and ambiguity; that a required facet kind is present;
                         that slots are filled only with its own references.
The planner never reads  a decision record or the content of a facet.
Outside the kernel       encoding a decision record as an ontos value; validating it against its horos type.
```

### 6.6 How to test the contract before the rewrite

We have done none of the following.

| Option | What is done before every rule moves | What it can show | What it cannot show |
| --- | --- | --- | --- |
| P1. Paper families | One family of each kind is written against the sketch of 6.5, for TypeScript and Go, without code: text, list, a record with an optional field, the declaration recipe with its generic part, an exact mapping of a generic (`Cell<T>`), one family with a requirement on its operand (a nullable value or a map), and Go's lowering for record and list (V3). For each: the decision record as a horos type, the projection, the requests, the facets, the entries for one example, what a factory receives | Gaps in what a record, a request or a facet can express | Paper finds what its author already suspects, and the cases most likely to break the contract (recursion, unions, explicit specialization) are the ones we would write worst |
| P2. A materially different fixture | A contribution in a separately owned package, of another profile or shape than the builtins, written against the contract's types with the kernel fixed (the evidence that section 2.1 asks for) | That a second author can use the contract without a shared edit | Nothing about a reader in another process |
| P3. A prototype on a branch | One target moved end to end for `Page<Person>`: contract, plan, serialized form, and one TypeScript factory round trip in the service library | The whole chain once, including the process boundary and the second repository, where sections 2.5 and 6.4 place the risk | It reverses the order of section 1.4 for one target and is thrown away or merged; three targets are still unproven |

**We lean to P1 before the rewrite and P2 with it, and have not decided on
P3.** What speaks against that: neither P1 nor P2 involves the reader in the
other process or repository, and that reader is the reason "once" matters.

### 6.7 What we already suspect

This list is open for answers: please confirm, rank or go beyond it.

1. A family whose decision depends on what was selected for its operand
   (6.1, 6.2).
2. Refusals as declared facts with classes, against errors from a callback
   (6.1).
3. Target configuration that changes a decision. The meaning of TypeScript's
   `next?: string` depends on the compiler option
   `exactOptionalPropertyTypes`; a decision record that silently assumes an
   option is wrong under another. The envelope has "target and its
   configuration".
4. A native declaration that a contribution owns without a horos declaration
   behind it (Go's `BittypeOptional`): who states its naming intent, how two
   contributions that need the same one share it, and how a plan refers to
   it.
5. The semantic version inside a contribution identity (`list/golang/1`)
   against the design's separate version axis, and which of the five axes a
   contribution states at all (6.4).
6. Generic conventions of exact mappings whose native parameters are
   "reordered, omitted and additional" (3.3), and what "`T`'s selected
   adapter travels with its native type" means for an omitted one.
7. One family wanted under two profiles, when a contribution states one
   profile and a request selects by exactly one profile string (3.4).
8. A payload schema that must change when the integer family arrives (a
   field index as an integer): whether that is a versioned change of one
   family or a second change of the contract.
9. How a library in a TypeScript repository authors a contribution (1.4).

## 7. Questions for the expert

In section 5, owner rules are fixed, and accepted designs are fixed unless
you see a contradiction; two that we see are named at the head of section 6.
Everything else is open, including what the code does today and every
leaning. Option labels refer to the tables of section 6. Please name an
option, or state another, and then the reason. Patterns from comparable
systems are welcome where they change the answer, including what failed:
compilers with separately registered lowering passes, interface compilers
with plugin protocols across a process boundary, serialization frameworks
that derive a type and its codec from one description, and type-class or
trait resolution, where resolving one implementation requires resolving
another.

1. **What a rule returns and what it may ask.** One selection per use must
   yield both the native type and what a backend needs to convert values.
   How would you shape (a) the unit of a contribution and the result of its
   rule, including how a rule declines a use (6.1: U1, U2, U3; R1, R2, R3);
   (b) the representation of a record's optional field, which is content of
   the record form and no use that the kernel walks: the record's own
   layout, a use that the source integration presents, or a request under a
   key, and in that case the form of the answer (6.2: D1 to D5; Type,
   Member, Record-only); (c) the way a rule learns facts about what was
   selected for its operand, as a nullable value or a map's key needs (3.3,
   6.1)? We lean to U1 and R2, and to D3 with Member; 6.1 and 6.2 each end
   with what speaks against that. Your answer decides the signature of a
   rule and whether the kernel gains a capability. If you choose another
   result than R2, please say which of the later questions change with it.

2. **Identities.** A backend in another process and a factory package in
   another repository must recognise a selected contribution; a conversion's
   memo is keyed by a representation identity plus the input; a generated
   native declaration has its name at generation and its place only
   afterwards (2.6). How would you form (a) a contribution identity that
   keeps owners apart and stays apart from the version axes (6.4: I1, I2,
   I3); (b) the entries of a generic declaration, open or one set per
   application, and from that the run-time identity and whether a plan needs
   use entries, representation nodes or both (6.4: O1, O2; N1, N2, N3);
   (c) a plan's reference to a native declaration, including one that only a
   contribution owns, so that a backend in another process never guesses a
   name or an import (6.5: S1 to S4)? We lean to I1, N2 and S2 and have no
   leaning between O1 and O2; 6.4 and 6.5 end with what speaks against that.
   Your answer decides the identity fields of the contract and whether a
   placement result becomes an input of a backend.

3. **What the kernel carries and cannot read.** Besides its native facet a
   contribution supplies a value adapter facet and may supply root,
   operation and lowering facets (typed target expressions). The kernel also
   carries each decision record unread. How would you (a) let the contract
   hold these parts: one field per kind, an open set whose kinds the reading
   products issue, or none (6.3: F1, F2, F3); (b) state what the kernel
   guarantees about a value it cannot read, a facet or a decision record,
   with a test for each guarantee (our two candidates stand in the row F2);
   (c) decide what a contribution must already be able to refer to (a
   generated member, a slot's native type, a declaration it owns) so that
   lowering can be added later without changing the contract? We lean to F2;
   6.3 ends with what speaks against it. Your answer decides the shape of
   the contribution type.

4. **Evidence before a contract is final.** One change moves 21 rule
   packages of four targets onto the new contract. Nothing is released, and
   a prototype on a branch is allowed before that. An earlier consultation
   set six vertical cases that must pass before the plan's envelope is
   called stable (2.3); the contribution contract has no such list. How
   would you divide what must be final in that change from what may stay
   provisional, given our partition and the tension with AC-10 at the head
   of section 6, and the suspects of 6.7? And which evidence would you
   require before the rewrite, for the six cases and for the sketch of 6.5:
   paper, a fixture, or one target end to end, with the order of section 1.4
   changed if need be (6.6: P1, P2, P3; 6.3: V1, V2, V3)? We lean to P1 and
   P2 and to V3, and have not decided on P3; 2.4, 6.3 and 6.6 say what
   speaks against that. Your answer decides what we build before the rewrite
   and whether "changes once" holds as we derived it.

5. **Open.** What should we know that these questions do not ask? An answer
   at any level is welcome: the division into kernel, contribution, plan and
   backend, the aim of changing the contract once, or a single field. Where
   a type generator came to carry conversion decisions for backends it does
   not own, what was learned late, and what would have shown it earlier?
   Section 6.7 holds what we already suspect; please confirm, rank or go
   beyond it. From your answer we decide what enters the contract before the
   rewrite and what we reconsider.

## Applied

**Evaluated against:** `0017-contribution-contract.submitted.md` (sha256 `66d4a2ef753841fb887bbe7cc7138183c944fddd6c40af776cabd99ea151d121`), attempt 1, run `run_49f8e87c-d37a-4578-ae05-65eaf57a60f5`; code at bittype `7e9dc34` (`main`), which differs from the described baseline `2daea52` by documents only. The living document differs from the submitted bytes only in its Status and Run-ID lines and in this section. bitspark-service was read at `c456bc7`.
**Graded:** https://github.com/Bitspark/bittype/issues/196#issuecomment-6840503518 (by the research owner, claude-contract-chain-issues; bittype has no architecture seat)

**How it was evaluated.** The advice was split into 40 recommendations. Three agents that did not write the document checked each one's premises against the code, the accepted designs and the state of the service library, by reading, and proposed a verdict with evidence. The research owner decided each row. No premise of the advice was false about the code.

**Delivery.** In this order: the design page `docs/contribution-contract.md`; a prototype on a branch that is thrown away; a materially different contribution with the shared packages fixed; the contract in its final form and the migration of the 21 mapping packages on `main` in one change, gated by row 40; the optional move by the answer of row 11. The work items are #196 (the contract), #195 (the plan), #207 (the slice) and #197 (the report). Three questions are with the owner of the service library on bitspark-service#38. This document reads `applied` only when every accepted row is implemented and validated.

**What the verification added that the advice could not know.**

- **The planner runs a rule before its children.** A rule's callback runs first, and it never sees what was selected for a child (`contracts/generation/go/planning.go`). The advice is right that a decision cannot depend on an operand today. One round of requests between today's two callbacks is enough for every case found.
- **A coherence requirement the submitted document never mentioned.** `docs/type-generation-plugins.md`: a mapping must not "inspect an argument and secretly change representation", and "mapping after substitution must agree with substituting mapped arguments into the mapped generic declaration". It bounds rows 7 and 8.
- **The order the advice asks for is circular as written.** The package it wants the proof to run through, `@bitspark/structural-adapters`, does not exist. Its owner decided on bitspark-service#38 that it is created with integers, which the plan puts after the migration. The same facts mean that no other repository builds on the contract before integers.
- **Optional-field conversion is contextual in every target,** not only in TypeScript: only the record can omit a field. A separate optional contribution therefore needs a field-level conversion slot that the service library's adapter interface does not have.
- **"Changes once" already fails for one accepted family.** Explicit specialization (#166) is on the accepted list, and its work item names the generation contract among what it touches.
- **The link from the plan to the output exists at one point and is discarded.** In #281 the seal holds planned order and sites together and then renumbers. The planner's declaration identity already is the logical symbol identity the advice asks for.
- **A backend finds a factory by module and export name,** not by the contribution identity. The collision argument of section 6.4 was weaker than written.

| # | Recommendation | Verdict | Action | Owner | Link | Validation |
|---|---|---|---|---|---|---|
| 1 | The unit is one independently selectable representation policy over a declared applicability domain | HOLDS | Adopted as the definition. "Library × target" in `docs/binding-generation.md` gets a dated clarification: it names ownership and distribution | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 2 | Generics are a division of responsibility, not a family that owns substitution | HOLDS | The division is documented as advised. #207 item 1 and the family list are reworded | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 3 | One decision record, with a staged planning protocol | NEEDS ADAPTATION | Adopted with one round: a rule names its dependencies and the evidence it needs, the planner resolves them, the rule decides with typed answers, the planner freezes the record, the projector runs. No general loop: no case found needs a second round. The evidence path is exercised in the prototype by the declaration recipe's need to know that a selected body is a record, and by a fixture family | claude-contract-chain-issues | #196 | A fixture family whose decision depends on evidence about what was selected for its operand plans and generates. Red today: a rule runs before its children and cannot see a selection (`planning.go`, the callback runs before `resolve`) |
| 4 | The projector is registered, reads no source and selects nothing | HOLDS | Source sites for diagnostics reach it through the record or the kernel's views | claude-contract-chain-issues | #196 | A projector that tries to read the source does not compile or fails the contract test. Red today: the one callback takes the reader |
| 5 | What belongs in the decision record | HOLDS | A per-selection choice is a field. An invariant of a semantic revision is carried by reference to that revision | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 6 | "Cannot disagree" means no sanctioned second decision point | HOLDS | `docs/representation-families.md` gets a dated clarification | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 7 | Typed evidence about the selected dependency | NEEDS ADAPTATION | Adopted, bounded by coherence: evidence about an operand bound to a binder comes from the binder's requirements and never from one application. Only facts of the selected representation are evidence, and no analysis runs on demand | claude-contract-chain-issues | #196 | `Page<text>` and the written-out record generate the same native type when a rule uses evidence, and a rule that reads one application to decide fails. Red today: no evidence path exists |
| 8 | Evidence states: false, unknown, deferred | NEEDS ADAPTATION | The outcomes of evidence are proved, refuted, unknown and exhausted, plus deferred to a binder requirement. For the disjointness of a union the refutation keeps the accepted name `overlap` of `docs/representation-families.md`. A choice made for want of proof must be the same for the generic declaration and for its written-out application (corrected on 13 November 2026, decision O2 on #196) | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 9 | References may be recursive; evidence cycles are refused | NEEDS ADAPTATION | The distinction and the refusal of an evidence cycle enter the contract now. Recursive references are row 9a | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 9a | Recursive representation references | DEFERRED | trigger: #164 is claimed and vertical case 1 is being built. Recursion is refused today (`grep -n "recursive native declaration" contracts/generation/go/planning.go`) | claude-contract-chain-issues | #164 | — |
| 10 | Refusals are structured outcomes | NEEDS ADAPTATION | The generation contract's existing values and verdicts are refined, not replaced: unsupported, refused by design, requirement unsatisfied, requirement unresolved, missing facet, ambiguous or missing contribution, limit, with an owner-qualified reason and the subject. A rule's refusal that matches no value today gets one. A selected exact mapping's refusal never falls back. The Files categories of #271 stay their own classification: Files has one formatter per language and no selection | claude-contract-chain-issues | #196 | Every refusal of the 21 mapping packages matches exactly one value of the contract. Red today: the record rule returns a diagnostic that matches none |
| 11 | Choose the presence domain explicitly | NEEDS ADAPTATION | bittype's position: the standard record representation delegates per-field presence through a published protocol (D3). The reasons: the accepted family list names optional fields as a family, and finding 7 of `docs/package-boundaries.md`. The accepted texts also support the other reading, and the choice needs a field-level conversion slot in the service library's adapter interface. So it is put to that library's owner on bitspark-service#38 and built in the prototype. If that owner refuses the slot, presence becomes the record's own layout (D1), finding 7 is withdrawn, and the Go optional package folds into the record contribution | claude-contract-chain-issues | #196, bitspark-service#38 | The answer on bitspark-service#38 is recorded here with the form that was built. Red until then |
| 12 | The source issues or validates the subject of the request | NEEDS ADAPTATION | The subject is the parent record's use, an operand use that the kernel mints, and a role selector that the source integration issues. The kernel keeps minting every use | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 13 | Member fragment plus a contextual field-conversion protocol | NEEDS ADAPTATION | Native half: the target's existing `field` form is the shared member vocabulary, and the member's name comes with the request. Conversion half: a field-state protocol for every target, not only TypeScript. Its slot kind is agreed with the service library (row 11) | claude-contract-chain-issues | #196, bitspark-service#38 | Absent, present with a value, present with null where admitted, and present with `undefined` round-trip apart through the record factory, in Go and in the TypeScript stand-in. Red today: no conversion exists |
| 14 | State the limit; contextual conversion is a product shape | HOLDS | A presence layout shared across fields needs the record policy replaced. The contract can express a contribution with no standalone value adapter | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 15 | Contribution identity: publisher and local key, plus a semantic revision | NEEDS ADAPTATION | Adopted. The publisher is the library that publishes the mapping, the accepted pages' "library's owner". A composition lock and a manifest are row 15a. The premise of section 6.4 about factory lookup is corrected | claude-contract-chain-issues | #196 | Two contributions with one local key under two publishers register side by side, and two with one publisher and one key fail. Red today: rule IDs have no owner |
| 15a | A composition lock that binds a publisher name to its manifest | DEFERRED | trigger: a second publisher registers a contribution (bitspark-service#39 item 2 starts) | claude-contract-chain-issues | bitspark-service#39 | — |
| 16 | The five version axes stay apart from the identity | HOLDS | The contribution reference carries the semantic revision. The other four axes are stated where each belongs and are not folded into the identity | claude-contract-chain-issues | #196, #195 | Each axis is exercised alone: an unknown required value of one axis is refused while the other four are unchanged. Red today: no plan carries the axes |
| 17 | A family contract reference that spans targets | NEEDS ADAPTATION | Two identities, not one string. A family contract has an identity and a revision that are the same in every target: what values the family covers and what they mean. A contribution has its own representation revision, which fixes the invariants of one target. A contribution reference names both. The publisher's local key is unique for that publisher across targets, and the target is a stated field of the contribution. A root facet that recognizes a family contract is row 17a (corrected on 13 November 2026, decision O4 on #196) | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 17a | `Service`'s root facet recognizes a family contract | DEFERRED | trigger: bitspark-service#39 starts its root facet. No root facet exists in either repository | claude-contract-chain-issues | bitspark-service#39 | — |
| 18 | N3: use entries and representation nodes | HOLDS | The two senses of "instance" on the accepted pages get a dated clarification, in bittype for `docs/representation-families.md` and by request for the binding contracts' memo sentence | claude-contract-chain-issues | #195 | Two uses of `Page<Person>` are two use entries on one representation node. Red today: a use of a named declaration leaves no record |
| 19 | The N2 tuple is not an equality key | HOLDS | Identity is over the selected semantic representation without diagnostic provenance | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 20 | O1: open generics primary, application nodes, separate argument conventions | HOLDS | An argument that the native convention omits is still represented for conversion. Noted on 13 November 2026 (decision O5 on #196): this applies when a request asks for bindings. A request names its products, and for types alone only demanded native actuals need mappings | claude-contract-chain-issues | #196 | A native convention that omits a parameter still yields an adapter argument for it. Red today: an omitted actual is never represented |
| 21 | O2 does not repair argument-dependent layout; decide specialization | HOLDS | Explicit specialization (#166) is in the extension domain. The contract can name an application, its selection and the resulting native representation. The first implementation refuses a specialization as unsupported | claude-contract-chain-issues | #196 | A specialization request is refused as unsupported with the application named. Red today: selection keys on the callee alone |
| 22 | The memo key uses an applied representation token | NEEDS ADAPTATION | The memo is the service library's. bittype's plan carries what the token is derived from: representation nodes, applications and selected dependencies. The token's form goes to bitspark-service#40 | claude-contract-chain-issues | #195, bitspark-service#40 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 23 | S2 with a checked realization manifest | NEEDS ADAPTATION | The planner's declaration identity is the logical symbol identity. The inventory record of #281 gains it as a planner-stated field. A serializable manifest relates plan, output and realization. `docs/binding-generation.md` is amended: the placed result is a fourth input of a backend. bitspark-service#40 item 9 repeats the sentence and is told | claude-contract-chain-issues | #195, #272 | A plan read with the manifest of another output fails as a linking error, not as an incomplete plan. Red today: neither exists |
| 24 | Read the realization now | NEEDS ADAPTATION | For TypeScript, with the reader on the formatter's result of #273. Go and Rust are row 24a | claude-contract-chain-issues | #273 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 24a | A realization reader for Go and Rust | DEFERRED | trigger: slice 1 of #161 or #159 delivers a generated-Files formatter (`ls format/files/golang format/files/rust`) | claude-contract-chain-issues | #161, #159 | — |
| 25 | An owned logical key for a contribution's declarations | HOLDS | The spelling is not the key. A contribution that needs whichever shared helper the composition selects makes a request | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 26 | F2 as a typed extension contract | NEEDS ADAPTATION | A facet is a kind-identified, versioned, validated value on a contribution. It is a second carrier beside the owner-issued values of #281, and the difference is stated: those attach to a declaration of one output and never leave the process, a facet attaches to a contribution and may cross a process boundary. A transport descriptor is row 26a | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 26a | A transport descriptor for a facet | DEFERRED | trigger: the first facet crosses a process boundary in the prototype of row 31 | claude-contract-chain-issues | #196 | — |
| 27 | What the kernel guarantees about values it cannot read | NEEDS ADAPTATION | Guarantees 1 to 6 are the kernel's, each with its test. Unknown required versions: the planner emits them and the plan's reader refuses. "Extension needs no shared case" is qualification evidence, not a kernel guarantee | claude-contract-chain-issues | #196 | One test per guarantee 1 to 6, each listed on the design page with what makes it red |
| 28 | Freeze the record before either projection consumes it | NEEDS ADAPTATION | Adopted. Encoding and validation live in a new package role outside the contract package, which may import neither horos nor ontos. The horos source reader is not its home. `LAYOUT.md` names the home | claude-contract-chain-issues | #196 | Mutating the value a rule returned, after it returned, changes neither the native type nor the serialized plan. Red today: no record exists |
| 29 | What the lowering environment must be able to express | NEEDS ADAPTATION | It is the checklist of the design. Reference categories are the contract's. Access and construction operations belong to the reading product's vocabulary. The prototype's Go value vocabulary and its facets stay in the test tree: `docs/binding-generation.md` keeps bittype's native vocabularies types-only | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 30 | The stabilization partition | HOLDS | The partition of the head of section 6 is replaced by the advised one: meanings final, encoding and layout provisional | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 31 | Evidence before the migration | NEEDS ADAPTATION | Adopted in the form bittype can do alone: the design page; then a prototype on a branch that is thrown away, with the staged rule, the frozen record, the plan and its serialized form, Go handlers and executed Go lowering in the test tree, a TypeScript reader as a stand-in in the test tree, and fixtures for vertical cases 3 to 6; then a materially different contribution with the shared packages fixed; then the contract in its final form and the migration of the 21 packages on `main` in one change. What a stand-in proves is stated: the serialized plan is sufficient for a second reader in another process, not for the service library's factory. Rows 31a and 31b hold the rest | claude-contract-chain-issues | #196, #207 | The gate of row 40 on the prototype branch, then on `main` |
| 31a | TypeScript end to end through the service library's own factory | DEFERRED | trigger: `packages/structural-adapters` exists on bitspark-service `main`, or its owner records on bitspark-service#38 that a prototype precedes integers | claude-contract-chain-issues | bitspark-service#38 | — |
| 31b | Vertical cases 1 and 2 in real form | DEFERRED | Vertical case 1 is not in the prototype: recursion is refused and deferred (row 9a), so a mock cannot exercise it at the kernel's level, and it waits for #164 as a whole. Vertical case 2 in real form needs unions, nullable values and exact forms. In the prototype it is exercised at the kernel's level through the independent fixture source, and the record says that the family is a mock. trigger: #164, #210 and #211 deliver their families (corrected on 13 November 2026, decision O6 on #196) | claude-contract-chain-issues | #164, #210, #211 | — |
| 32 | The six cases test the contract too; identity and linking cases | HOLDS | Added to the prototype: repeated uses share a node, equal native types with different adapters stay apart, an erased native argument is kept for conversion, one plan under two TypeScript placements | claude-contract-chain-issues | #207 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 33 | "Changes once" as a bounded commitment; write down the extension domain | HOLDS | Adopted in the advised words. The domain is written into the design page: further products and facets in; a contextual representation at a field position in; cross-field layout coordination out; application-specific specialization in and refused until #166; one authoring route, a Go contribution package. Noted on 13 November 2026 (decision O10 on #196): the line on a contextual representation at a field position is conditional on the decision of row 11 | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 34 | The source-semantics input needs a form | NEEDS ADAPTATION | For the prototype: a snapshot of ground descriptors and module pins exported from the Go side, whose source identities match the plan's. A general reader protocol is row 34a | claude-contract-chain-issues | #195 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 34a | A reader protocol for source semantics | DEFERRED | trigger: the design of bitspark-service#40 item 3 starts | claude-contract-chain-issues | bitspark-service#40 | — |
| 35 | Record target configuration; state what a present `undefined` means | HOLDS | The optional family's contract states it, and the four states are tested apart. bittype checks under `exactOptionalPropertyTypes`, and the service library compiles without it and reads a present `undefined` as absence. That difference is put to its owner | claude-contract-chain-issues | #207, bitspark-service#38 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 36 | A memo hit must not hide resource cost | STALE / INVALID | Stale for bittype: laws K1 to K4 of the binding contracts state it and bittype confirmed them on 12 November. The two points they leave unstated go to the service library | — | bitspark-service#38 | — |
| 37 | An actual contribution-authoring route | NEEDS ADAPTATION | One route: a Go contribution package, tested out of tree. A declarative sidecar is not promised. The service library's own first contribution is row 37a | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 37a | A contribution authored in the service library's repository | DEFERRED | trigger: bitspark-service#39 item 2 starts | claude-contract-chain-issues | bitspark-service#39 | — |
| 38 | Payload bootstrapping must not distort the contract | HOLDS | A pinned bootstrap artifact: the schema source, its canonical module with its pin, and the committed generated types with a regeneration test | claude-contract-chain-issues | #196 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 39 | The "another pass" mutant proves nothing | HOLDS | The boundary is enforced by what a projector can reach. Mutants must disagree on presence, layout, slot identity or environment. #207's verification and row 2 of research 0015 are reworded | claude-contract-chain-issues | #207 | its line of the design page `docs/contribution-contract.md`, with the test that the gate of row 40 names for it |
| 40 | A stronger acceptance condition for item 6 | NEEDS ADAPTATION | Adopted as the merge gate of the migration, each line in the form bittype can show: a frozen decision survives native projection, dependency composition, serialization, a second reader in another process, executed Go lowering, and realization for TypeScript, with every shared package named and fixed while a materially different contribution is added | claude-contract-chain-issues | #196 | The seven lines of the gate, each with its test named on the design page. The migration does not merge while one is red |

**Corrected on 13 November 2026.** Writing the table out as the design page `docs/contribution-contract.md` showed ten points that did not fit without a reading. The research owner decided them as [O1 to O10](https://github.com/Bitspark/bittype/issues/196#issuecomment-6841190077) on #196. Three correct an Action above: row 8 (O2), row 17 (O4) and row 31b (O6). Rows 20 and 33 carry a note (O5, O10). The other five change no row.

**Statements of the submitted document that this evaluation withdraws.** The partition "contract final, envelope unstable" at the head of section 6. The leaning N2 and its key. The leaning I1. The answer shape "Member" as sufficient. The order of section 1.4, in which the contract would be final before anything reads its new half. The sentence that a rule returns its decision before its dependencies are resolved. The collision argument of section 6.4. The leaning V3 as sufficient evidence.
