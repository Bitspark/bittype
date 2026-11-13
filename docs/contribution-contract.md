# The contribution contract

**Design, 13 November 2026, from research 0017. Decided by the research owner, `claude-contract-chain-issues`, who also wrote the research question. Not implemented. It becomes final when the gate at the end of this page is met; until then the prototype may change it, and each change is recorded here.**

This page is the contract of a contribution to bittype's generation kernel:
what a contribution is, what its rule may ask and must return, what
identifies it, and what the kernel carries for it without reading it. It is
step 1 of the [order of work](#the-order-of-work) that
[research 0017](../research-docs/0017-contribution-contract.md) decided for
[#196](https://github.com/Bitspark/bittype/issues/196), item 6. The mapping
contract of today does not change with this page.
[The code today](#the-code-today) describes the present behaviour, and nothing
else on this page is implemented.

The page rests on two accepted pages and amends neither in substance:
[binding generation](binding-generation.md) (accepted direction of 9 November
2026, research 0014) and [representation families](representation-families.md)
(accepted design of 9 November 2026, research 0015). Where the table of
research 0017 clarifies a sentence of theirs, the clarification is dated and
listed under [Amendments to accepted pages](#amendments-to-accepted-pages).

**How to read the sources.** Each statement names its source.

- **Row N:** a row of the table "Applied" of
  [research 0017](../research-docs/0017-contribution-contract.md#applied). The
  table has 49 rows. Each has a verdict and an Action, and the Action is the
  decision. Where the Action differs from the advice, the Action governs.
- **The grade:** the
  [comment of 13 November 2026](https://github.com/Bitspark/bittype/issues/196#issuecomment-6840503518)
  on #196 in which the research owner graded that table. bittype has no
  architecture seat, so the owner graded and said so there.
- **The advice:** the
  [expert's reply](../research-docs/0017-contribution-contract.advice.md), by
  its section number. It gives the reasoning of a row and is quoted only where
  the row accepted it.
- **An accepted page**, quoted, with its file and line.
- **AC-01 to AC-10:** the
  [numbered architecture rules](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#numbered-rules)
  of platform, read at `5e988b0`.
- **BC:** the service library's
  [binding contracts](https://github.com/Bitspark/bitspark-service/blob/c36722db901ef7a259242e168c4c35d6f777fbf4/docs/design/binding-contracts.md)
  at `c36722d`, accepted by that repository's owner seat on 11 November 2026.

A sentence marked **derived here** follows from those sources and is stated by
none of them. It is this page's reading and no decision. A point that two
sources leave unsettled, or settle differently, is not resolved on this page:
it stands under [Open inside this design](#open-inside-this-design) with its
options. File and line references are `path:line` at commit `ec96b92`, the
head of `main` when this page was written. The line numbers of the two
accepted pages that this change annotates are those of that commit.

**Words used on this page.**

| Word | Meaning here |
| --- | --- |
| The kernel | The generation contract package `contracts/generation/go` with its planner, and the registry that selects and composes |
| The source integration | The package that reads one source language for the kernel. For horos it is `transform/syntax/horos/source/go` |
| A rule | The part of a contribution that the planner runs for one selection |
| A decision record | The frozen value that holds what a rule decided for one selection |
| A projector | The registered function that computes native intent from a frozen decision record |
| The plan | The selected representation plan of binding generation, tracked in [#195](https://github.com/Bitspark/bittype/issues/195) |
| A product | Something a request asks generation to make possible: native types, or one kind of binding |
| A reader | A program that consumes the plan: the completeness test backend, or a binding backend |

## What a contribution is

**The unit (row 1).** Row 1 adopts the advice's definition (section 1.1):
"One independently selectable representation policy, for one target and one
profile, over a declared applicability domain, including its paired facets."

It is the owner's definition of a contribution applied to representation:
"A contribution, meaning one thing a consumer registers or selects on its own,
is one unit." (platform `docs/concepts/architectural-correctness.md:463-464`).
A consumer can therefore replace one contribution and keep every other.

**The applicability domain.** A contribution declares the domain it answers
for. Representation families already require that as the first part of a
family's contract: "The source identities and forms it accepts, its generic
assumptions, its targets, its explicit refusals"
(`docs/representation-families.md:36`). The advice names four kinds of
domain, and row 1 adopts them with the definition:

| Kind of domain | Example | Standing |
| --- | --- | --- |
| An exact declaration | `Cell<T>` of one module pin | Row 1 |
| A builtin form | The closed record, the list | Row 1 |
| The declarations that need the default native declaration policy | Every named declaration without an exact mapping | Row 1 |
| A role | "Optional record field" | Row 1, and conditional on [row 11](#optional-fields) |

The kernel selects by the domain and never by registration order (row 10,
from the advice, section 1.5). Two contributions that claim one domain under
one profile need an explicit composition by the consumer (row 15, from the
advice, section 2.1).

**One target and one profile.** A contribution states both. The profile is
stated by the contribution and is no argument of its registration:
"A mapping states its profile in its own description, and its constructor
takes none." (`docs/representation-families.md:330`). What a profile is:
"Representation profiles state what the native type retains and what remains a
source-level obligation." (`docs/type-generation-plugins.md:507-508`). A
profile also carries the coherence requirement that bounds
[evidence](#evidence).

**"Library × target" (row 1).** Binding generation says: "The unit of
implementation is library × target." (`docs/binding-generation.md:265`). Row
1: that sentence "names ownership and distribution". A library owns and
distributes the contributions it publishes for one target, and it may offer
one function that registers all of them. That packaging does not make them one
selection. The unit of selection is the single contribution. The sentence gets
a dated clarification, listed under
[Amendments](#amendments-to-accepted-pages).

**Generics are a division of responsibility (row 2).** No contribution owns
substitution, and "generics" is no family with a contribution of its own. Row
2 documents the division "as advised" (advice, section 1.1):

| Responsibility | Owner |
| --- | --- |
| Binder identity, scope, substitution, applications, the links of an environment | The kernel and the source integration, each through its contract |
| Native parameter order, omitted and additional native parameters, alias or defined declaration | The contribution that supplies the declaration: the default declaration policy or an exact mapping |
| Adapter parameters and requirements | The selected representation contribution |
| Instantiation and compatibility checks | The planner, which invokes those contracts |

This agrees with the planner page: "The kernel owns capture identity,
substitution and scope." (`docs/generation-planner.md:284`). It changes two
texts that list generics beside the families. Item 1 of
[#207](https://github.com/Bitspark/bittype/issues/207) is reworded by the
[comment of 13 November 2026](https://github.com/Bitspark/bittype/issues/207#issuecomment-6840545162)
on that issue. The delivery row of representation families that names
"generics" among the delivered families
(`docs/representation-families.md:341`) gets a dated clarification.

## What a rule does: the protocol

**One round (row 3).** Row 3: "a rule names its dependencies and the evidence
it needs, the planner resolves them, the rule decides with typed answers, the
planner freezes the record, the projector runs."

| Step | Who | What happens |
| --- | --- | --- |
| 1. Name | The rule | From the checked subject and the source view, the rule names its dependencies, and for each the evidence it needs. A dependency is a child use that the kernel minted, or a [request](#optional-fields) for whatever the consumer registered for a meaning. Each is named as a reference or as a demand for evidence ([row 9](#evidence)) |
| 2. Resolve | The planner | It selects a contribution for each named dependency, plans it, and collects the evidence asked for |
| 3. Decide | The rule | With the source view and the typed answers, the rule returns a decision record or a [refusal](#refusals) |
| 4. Freeze | The planner | It freezes and validates the record ([row 28](#what-the-kernel-carries-without-reading-it)) before anything reads it |
| 5. Project | The registered projector | It computes native intent from the frozen record. The entry of the plan is made from the same record |

The advice (section 1.2): "One selection does not imply one callback
invocation." The kernel keeps one selected contribution and advances it
through the steps.

**The kernel drives.** A rule never invokes the planner. That is so today, and
it stays: "Capture and Expand run only during planning. They do not
recursively invoke the planner or construct AST nodes."
(`docs/generation-spi.md:58-59`). Naming a dependency is a statement to the
planner, as it is today: "Mapped records a dependency. It does not execute
another mapping." (`docs/generation-spi.md:92`).

**No analysis on demand (row 7).** What a rule may name in step 1 is bounded.
It may ask the planner to represent a child use or a source-issued subject,
and it may ask for evidence that the selected dependency states about itself.
It may not ask for source analysis. Row 7: "no analysis runs on demand". The
accepted sentence behind it: "Neither the shared host nor a plugin invokes
analysis on demand." (`docs/type-generation-plugins.md:535`).

**The projector (row 4).** Row 4: "The projector is registered, reads no
source and selects nothing".

- **Registered.** It belongs to the contribution. It is not a function that a
  decision returns. The advice (section 1.2) gives the reason: a function made
  per decision "could capture the source reader or an unstored choice".
- **What it may read.** The frozen decision record, typed views of the
  dependencies that the planner resolved, and views of the symbols the plan
  issued. It gets no source reader.
- **Diagnostic sites.** A native expression carries the source site it came
  from. Row 4: "Source sites for diagnostics reach it through the record or
  the kernel's views".
- **A limit, derived here.** The projector's inputs keep the source reader
  away from it. They cannot keep a contribution from passing what it read to
  its projector through state of its own package. A contribution stays
  trusted in that respect, as a target is today
  (`docs/generation-spi.md:369-372`).

**What belongs in the decision record (row 5).** Row 5: "A per-selection
choice is a field. An invariant of a semantic revision is carried by reference
to that revision".

- **A per-selection choice** is anything that can differ between two
  selections of the same contribution and that affects representation or
  conversion. The advice (section 1.2) lists member names, the presence
  strategy, the selected construction strategy, a layout that depends on an
  operand, and applicable configuration.
- **An invariant** is what every selection of one semantic revision does the
  same way. It needs no field. The plan names the selected contribution and
  its semantic revision ([identities](#identities)), and that reference
  carries the invariant.
- **The rule for a projector.** The advice (section 1.2): "A projector may
  implement a recorded policy; it must not choose an unrecorded policy."

Binding generation lists what the plan carries, among it "the rules for
absence, presence, null and native `undefined`"
(`docs/binding-generation.md:58`). Row 5 lets an invariant rule be carried by
reference. Whether that sentence needs a clarification is
[open](#open-inside-this-design), point 1.

**"Cannot disagree" (row 6).** Representation families: "No later pass
reconstructs a decision from the source, so a native declaration and its
conversion cannot disagree." (`docs/representation-families.md:46-47`). Row 6
reads the last two words as "no sanctioned second decision point".

- **What the sentence claims.** There is one origin of every decision, the
  frozen record. No interface lets a later pass select again or rebuild a
  decision from the source.
- **What it does not claim.** That a projector or a factory cannot be wrong.
  AC-05: "Selecting a contribution does not establish the truth of its
  claims." The same page already says of its tests: "Deliberately broken
  rules show that the tests detect a failure. They do not prove the laws"
  (`docs/representation-families.md:321-322`).

The sentence gets a dated clarification
([Amendments](#amendments-to-accepted-pages)).

**There is no general loop (row 3).** A rule gets one round of requests. It
cannot ask again after it has seen an answer. Row 3: "No general loop: no case
found needs a second round." The grade gives the second reason: "an untested
loop would be final without evidence."

What would add a second round is a named case that needs one. Derived here,
from verification report 1: such a case is one in which the set of
dependencies, or the evidence wanted of one, depends on an answer of the first
round. None of the six delivered families and none of the six vertical cases
of the plan (`docs/representation-families.md:95-102`) is such a case by its
description. Adding the round then is a change of the shared contract and
needs the justification of AC-08.

**An illustration.** The advice sketches the protocol and says of its sketch:
"This is a protocol sketch, not a proposed Go API." The same holds for the
lines below. They show the order of one round and decide no signature.

```text
name(checked subject, source view)
    -> dependencies, each a reference or a demand for evidence
decide(checked subject, source view, typed answers)
    -> decision record | refusal
project(frozen record, dependency views, symbol views)
    -> native intent
```

**Why the shared contract changes. Derived here.** AC-08 asks three things
of a change to a shared contract. The pull request that changes the contract
states them in full. In short:

- **The missing capability.** A rule cannot learn what was selected for a
  dependency, and what it decides leaves no record that a second reader can
  use.
- **Why it belongs to the kernel.** Selection, scope and the order of
  planning are the planner's: "Exact contribution selection, substitutions,
  references, and the selected representation plan"
  (`docs/binding-generation.md:36`).
- **How it serves independently owned contributions.** A structural
  contribution and an exact mapping name dependencies and receive answers in
  the same way. The kinds of evidence belong to vocabularies, and the kernel
  gains no case for any of them.

## Evidence

**Typed evidence about the selected dependency (row 7).** A rule that needs
to know something about a dependency asks for a typed fact about what was
selected for it. It does not read the payload of another family.

| Question | Answer | Source |
| --- | --- | --- |
| What can be evidence? | A fact that the selected representation states about its own native side. Row 7: "Only facts of the selected representation are evidence" | Row 7 |
| Who owns a kind of evidence? | The vocabulary of the target or of the product that the fact concerns. "The kernel should not acquire cases for nullability, comparable keys, or every future domain property." | Advice, section 1.3 |
| What does the kernel check? | That an answer belongs to the dependency it was asked about: the selected representation and its semantic revision, the target, profile and configuration, the generic environment, the kind and version of the requirement, and the scope | Advice, section 1.3; [guarantee 4](#what-the-kernel-carries-without-reading-it) |
| What does the kernel not check? | That the fact is true. "It cannot establish a domain assertion merely because a contribution supplied it." | Advice, section 1.3; AC-05 |

Facts about the source stay with the source. Representation families divide
the judgment: "horos semantics and bittype's source integration establish
canonical facts", and "target families establish native discrimination and
compatibility" (`docs/representation-families.md:191-194`). Evidence in the
sense of this section is the second kind only. A rule reads source facts from
the source view, as it does today.

**Bounded by coherence (row 7).** A profile requires that a generic
declaration and its applications agree:

> "Exact selection alone is insufficient: a plugin could inspect an argument
> and secretly change representation even though the selected key remained
> the same. The profile therefore requires one coherent representation family
> over its supported substitutions: mapping after substitution must agree with
> substituting mapped arguments into the mapped generic declaration."
> (`docs/type-generation-plugins.md:480-485`)

And: "Open eligibility must be established under binder assumptions over the
profile's explicitly represented substitutions; success on a few concrete
applications does not establish it." (`docs/type-generation-plugins.md:503-505`).

The submitted research document did not show this requirement to the expert.
The verification found it, and row 7 bounds evidence by it:

> **The rule for an operand bound to a binder.** Row 7: "evidence about an
> operand bound to a binder comes from the binder's requirements and never
> from one application."

Inside `Page<T>`, a rule that asks about the operand `T` is answered from
what `T`'s declared requirements establish. What `Person` is at
`Page<Person>` does not reach the rule that planned the body of `Page<T>`.
The open declaration is emitted once, so nothing learned at one application
can change it.

**The outcomes (row 8).** An answer has one of five outcomes. Row 8: "The
accepted outcome names are used (proved, refuted, unknown, exhausted) plus
deferred to a binder requirement."

| Outcome | Meaning | What a rule may conclude |
| --- | --- | --- |
| Proved | The selected representation establishes the fact | The fact |
| Refuted | It establishes the opposite | The opposite |
| Unknown | Nothing establishes either | Nothing. The advice (section 1.3): "Absence of proof that an operand excludes null is not proof that it does." |
| Exhausted | The check reached its limit | Nothing |
| Deferred to a binder requirement | The operand is bound to a binder. The answer is what the binder's declared requirement establishes, and each application must discharge that requirement | The fact, if the requirement states it. Otherwise nothing |

The names of the first four come from representation families, which gives
them for one analysis as "`proved`, `overlap`, `unknown` and `exhausted`, and
neither of the last two establishes anything"
(`docs/representation-families.md:199-201`). That page says `overlap` where
row 8 says refuted. The difference is
[open](#open-inside-this-design), point 2.

The fifth outcome is the accepted course of a generic family: "A generic
family proves its requirements for every permitted actual, declares
requirements that actuals must satisfy, or defers qualification to checked
applications." (`docs/representation-families.md:252-254`).

**A choice made for want of proof (row 8).** A rule that has no proof may
take a representation that needs none, if its contribution offers one, or it
may refuse. Two limits hold.

- **It never replaces what was selected for the operand.** The choice is the
  rule's own, and it is recorded in the rule's decision record. The accepted
  pages: "A validated canonical-value wrapper is a deliberate selected mapping
  for hard forms. It is never a fallback."
  (`docs/representation-families.md:235-236`), and "No retry selects a weaker
  representation." (`docs/kernel-composition.md:98`).
- **It is the same inside a generic declaration and outside it.** Row 8: "A
  choice made for want of proof must be the same for the generic declaration
  and for its written-out application".

An example, derived here from verification report 1. A rule takes a direct
form for the operand `text`, because the fact is proved there, and a wrapper
for the operand `T`, because nothing is proved about a binder. `Page<text>`
then has a wrapped field, and the same record written out with `text` has a
direct one. That is the disagreement the coherence requirement forbids. Such
a rule has to treat the binder through a declared requirement, or use one
representation for every operand it admits. The test of row 7 states the
outcome: "`Page<text>` and the written-out record generate the same native
type when a rule uses evidence, and a rule that reads one application to
decide fails."

**References against evidence (row 9).** Step 1 of the protocol names each
dependency as one of two things.

| | A reference | A demand for evidence |
| --- | --- | --- |
| What the rule needs | To refer to the other representation | An answer before its own decision can finish |
| May it be recursive? | In the design, yes. Today it is refused (row 9a) | No. A cycle is refused |

Row 9: "The distinction and the refusal of an evidence cycle enter the
contract now."

- **Why the distinction enters now.** Derived here, from AC-10: a first
  change "explicitly refuses the cases it cannot yet check". If the kind of a
  dependency were not part of how a rule names it, lifting the refusal of
  recursion later would change how every rule names a dependency.
- **An evidence cycle is refused.** A rule that demands evidence from a
  dependency whose own decision demands evidence from that rule cannot be
  planned. The advice (section 1.4) says what would lift the refusal: "a
  declared evaluation discipline", and not recursion through callbacks.
- **Recursive references are deferred (row 9a).** Accepted for the families
  that need them: "Recursive adapters are built through explicit plan
  references (placeholders), separately from native value indirection."
  (`docs/representation-families.md:223-224`). They are not built with this
  contract. A recursive reference is refused. The trigger is under
  [Deferred](#deferred-with-triggers). What the prototype shows of recursion
  is [open](#open-inside-this-design), point 6.
- **Two refusals that can be told apart. Derived here.** The distinction of
  row 9 can be observed only if the refusal of an evidence cycle and the
  refusal of a recursive reference carry different reasons
  ([refusals](#refusals)).

## Refusals

**What exists (the generation contract at `ec96b92`).** Listed here because
row 10 refines these and replaces none.

| Value or verdict | Its text or meaning | Where |
| --- | --- | --- |
| `ErrInvalid` | "invalid generation description" | `contracts/generation/go/contracts.go:14` |
| `ErrUnsupported` | "generation capability unavailable" | `contracts.go:15` |
| `ErrScope` | "incompatible generation scope or association" | `contracts.go:16` |
| `ErrLimit` | "generation limit exhausted" | `contracts.go:17` |
| `ErrIncomplete` | "generation obligations incomplete" | `contracts.go:18` |
| `ErrAmbiguous` | "ambiguous generation mapping" | `contracts.go:19` |
| `ErrRejected` | "source judgment rejected the use" | `contracts.go:21` |
| The four verdicts of a source judgment | `Admitted`, `Unsupported`, `Rejected`, `Undetermined` | `contracts/generation/go/source.go:83-92` |
| The value of the verdict `Undetermined` | `semantic.ErrEvidenceUnavailable`, from the semantic contract | `source.go:187-188` |

**The refinement (row 10).** Row 10: "The generation contract's existing
values and verdicts are refined, not replaced: unsupported, refused by design,
requirement unsatisfied, requirement unresolved, missing facet, ambiguous or
missing contribution, limit, with an owner-qualified reason and the subject."

| Refusal | Meaning | Relation to what exists. Derived here, from verification report 1 |
| --- | --- | --- |
| Unsupported | The selected contribution recognizes the case and cannot represent or check it yet. A later slice can admit the same input | `ErrUnsupported`, narrowed. Today that value also covers a missing mapping, a missing helper and recursion |
| Refused by design | The case lies outside the domain that the representation promises. Another contribution may cover it | New. It is not `ErrRejected`, which is the source's refutation of a use |
| Requirement unsatisfied | A selected dependency fails a declared requirement: the evidence is refuted | New |
| Requirement unresolved | The evidence a decision requires is unknown or exhausted | The same meaning as `semantic.ErrEvidenceUnavailable` has for source evidence |
| Missing facet | The selected contribution has no facet for a requested product | New. It stays a class of its own |
| Ambiguous or missing contribution | More than one contribution, or none, answers a domain | `ErrAmbiguous` for the first. The second is under `ErrUnsupported` today |
| Limit | A resource bound of the request is reached | `ErrLimit` |

The first four meanings follow the advice (section 1.5). Which outcome of
[evidence](#evidence) leads to the third and which to the fourth is derived
here. "Missing facet" stays
recognizable because the accepted design requires it: "Unsupported use,
missing facet, incompatible plan or runtime, invalid native value, structural
mismatch, domain violation, declared coverage restriction, cycle and
exhaustion stay distinct." (`docs/representation-families.md:303-305`). A
limit is no judgment about the input: "Resource exhaustion lies outside these
sets: a valid value stays valid when a call lacks the resources to convert
it." (`docs/representation-families.md:117-119`).

**What stays as it is. Derived here** from "refined, not replaced":
`ErrInvalid`, `ErrScope`, `ErrIncomplete`, `ErrRejected`, the four verdicts of
a source judgment, and cancellation. A source judgment stays the source's. A
refusal of this section is a contribution's or the planner's.

**What a refusal carries (row 10).** "an owner-qualified reason and the
subject". The advice (section 1.5) lists what to preserve: the source
subject, the selected contribution, the target and profile, the path of
dependencies and the requirement concerned.

- **The class is the kernel's, and the set is closed.** Failure meaning is
  one of the things composition must preserve (AC-05). A closed set of classes
  in the kernel is therefore no catalogue of contributions.
- **The reason is its owner's, and the set is open.** A reason explains. It
  never adds a class or moves a refusal to another one. Derived here.

**Every refusal of a rule has a value (row 10).** "A rule's refusal that
matches no value today gets one." The test of the row: "Every refusal of the
21 mapping packages matches exactly one value of the contract."

**No fallback (row 10).** "A selected exact mapping's refusal never falls
back." This keeps three accepted sentences:

- "There is never a fallback to structural conversion, in the same way that
  an exact use mapping never silently becomes a structural publication today."
  (`docs/binding-generation.md:99`)
- "An exact binding takes precedence, and its failure never falls back to
  that policy." (`docs/generation-spi.md:126-127`)
- The advice (section 1.5): "Applicability and refusal are not a mechanism
  for trying contributions in registration order."

So a refusal ends the planning of its subject. It is never an instruction to
try the next contribution.

**Why this is not the classification of the Files stage (row 10).** "The
Files categories of #271 stay their own classification: Files has one
formatter per language and no selection".

- **Files** tells unsupported coverage from invalid input by one question:
  "a refusal is unsupported coverage when a later slice can admit the same
  input. It is invalid input when no later slice can."
  (`docs/naming.md:316-318`,
  [#271](https://github.com/Bitspark/bittype/issues/271)). One formatter
  answers for one language.
- **Generation** has three parties in a refusal: the source, the selected
  contribution and a dependency. "Refused by design" has no counterpart in
  Files, because it means that the input is valid and that this contribution
  will never cover it. That can only be said where something is selected.
- **What the two share. Derived here**, from verification report 1: the
  criterion for unsupported, the meaning of missing evidence, and the limit.
  They share no value: the contract package imports no package of the Files
  stage (`integration/go/architecture_test.go:182`).

**Not decided by row 10**, and listed under
[Not specified here](#what-is-delegated-and-what-is-not-specified-here): the
form of a refusal in Go, the value of the two refusals of
[row 9](#evidence), whether a missing contribution and an ambiguous one are
one value or two, and what an unexpected error of a rule is.

## Optional fields

> **This whole section is conditional.** It states bittype's position (row
> 11). The position needs a slot in the service library's adapter interface
> that its owner has not agreed to. The question is
> [with that owner](https://github.com/Bitspark/bitspark-service/issues/38#issuecomment-6840511750)
> on bitspark-service#38 since 13 November 2026. Nothing in this section is
> fixed until the answer is recorded here and in the table.

**Two domains, both legitimate (row 11).** An optional field is no type in
horos. It is a flag on a field of a record. Who decides its native
representation can be read two ways, and the advice (section 1.6) says so:
"The general owner rule alone does not prove that presence representation
must be independently replaceable. D1 could be legitimate if presence belongs
wholly to the record policy."

| | D1: presence is the record's own layout | D3: the record delegates presence |
| --- | --- | --- |
| Who decides how an optional field is represented | The record contribution | A separately selected contribution, whose domain is the role "optional record field" |
| What a consumer replaces to get another representation | The record contribution | The optional contribution |
| The accepted text that supports it | The payload column of the plan lists "Record field layout and presence behaviour" as one item (`docs/representation-families.md:73`) | The family list names "optional and nullable fields" (`docs/representation-families.md:10-11`), and the delivery table names "optional fields" among the delivered families (`:341`) |
| What it needs from the service library | Nothing | A field-level conversion slot |

**bittype's position (row 11).** Row 11: "the standard record representation
delegates per-field presence through a published protocol (D3)". Its two
reasons, in the row's words: "the accepted family list names optional fields
as a family, and finding 7 of `docs/package-boundaries.md`". The row also says
what weakens them: "The accepted texts also support the other reading, and the
choice needs a field-level conversion slot in the service library's adapter
interface."

Finding 7 is bittype's own derivation and no owner sentence. It reads "Record
decides optional fields itself in Go, Rust and C++"
(`docs/package-boundaries.md:156`) and names three targets. Derived here:
under D3 the TypeScript record contribution asks as well, because the
position is about the standard record representation of every target.

**What happens under each answer (row 11).**

| The answer on bitspark-service#38 | What follows |
| --- | --- |
| The slot is admitted | D3 stands. Each target's record contribution asks, each target has an optional contribution, and the move is step 5 of the [order](#the-order-of-work) |
| The slot is refused | Row 11: "presence becomes the record's own layout (D1), finding 7 is withdrawn, and the Go optional package folds into the record contribution". Whether the request of row 12 then stays in the contract is [open](#open-inside-this-design), point 10 |
| No answer yet | The prototype builds D3. Row 11: "put to that library's owner on bitspark-service#38 and built in the prototype" |

**The subject of the request (row 12).** A record contribution that delegates
presence makes a request in step 1 of the [protocol](#what-a-rule-does-the-protocol).
Row 12: "The subject is the parent record's use, an operand use that the
kernel mints, and a role selector that the source integration issues. The
kernel keeps minting every use".

- **It is a field position, not a type.** The advice (section 1.6): "this
  subject represents a **field position**, not a horos type."
- **The kernel keeps its scope check.** It checks that the operand use is a
  child of the parent use, and that the role selector was issued for that
  pair. The advice (section 1.6): "Accepting an arbitrary selector string and
  arbitrary operand would weaken the present checked-child-use boundary."
- **The source integration issues the role.** Whether a field is optional is
  a fact of the source. The kernel learns a selector and reads neither the
  flag nor the field's name. Derived here, from verification report 2.
- **No use is invented.** A source integration does not present a use for a
  field position. Only the kernel mints uses, and row 12 keeps that.

**The answer has two halves (row 13).**

| Half | What row 13 decides |
| --- | --- |
| Native | "the target's existing `field` form is the shared member vocabulary, and the member's name comes with the request" |
| Conversion | "a field-state protocol for every target, not only TypeScript. Its slot kind is agreed with the service library (row 11)" |

- **The native half.** The record contribution and the optional contribution
  both use the member vocabulary of the target's native model. Neither imports
  the other. The record chooses the member's name, as the naming design has
  it: "A member's name and accessibility" belongs to "The declaration, so the
  mapping of that pair" (`docs/naming.md:154`). How a name travels in a
  request is [open](#open-inside-this-design), point 3.
- **The conversion half.** Only a record can omit a field, in every target.
  So no target converts an optional field from its retrieved value alone. For
  TypeScript the accepted design already says it: "the adapter separates an
  absent own property from a present `undefined` at run time"
  (`docs/representation-families.md:167-168`). The advice (section 1.6) lists
  what the protocol must be able to do: test presence, obtain a present value,
  delegate that value to the adapter selected for the operand, and construct
  or omit the native field according to the recorded decision.
- **Who invokes it.** The record's factory invokes the protocol for a field.
  It does not read the payload of the optional contribution.
- **The states.** The protocol tells four states of a field apart. They, and
  what a present `undefined` means, are under
  [Inputs that are not bittype's alone](#inputs-that-are-not-bittypes-alone).

**The limits (row 14).**

- **A layout shared across fields is out.** Row 14: "A presence layout shared
  across fields needs the record policy replaced." A request concerns one
  field. A bitmap for all optional fields of a record cannot be answered field
  by field.
- **A contribution without a value adapter is in.** Row 14: "The contract can
  express a contribution with no standalone value adapter." An optional
  contribution has a field-state protocol and nothing that converts a value on
  its own. The advice (section 1.6): "Contextual field conversion is an
  admitted product shape that the contract needs to express now."

**What this position changes outside this page, if it stands.** None of the
following is edited now. Each waits for the answer.

| Sentence | Where | Owner |
| --- | --- | --- |
| "A family reaches other families only through its adapter slots." | The [decision of 11 November 2026](https://github.com/Bitspark/bitspark-service/issues/38#issuecomment-6803983510) on bitspark-service#38 | The service library, `codex-platform-status` |
| "Run-time factory reference and typed argument-adapter slots, converting the value inside a message" | `docs/binding-generation.md:94` | Research 0014 |
| "Record field layout and presence behaviour" | `docs/representation-families.md:73` | Research 0015. This chain holds its delivery |
| Finding 7 | `docs/package-boundaries.md:156` | bittype. It gets a dated note now, which says that it is conditional ([Amendments](#amendments-to-accepted-pages)) |

## Identities

The [list at the end of this section](#the-identities-in-one-list) names
every identity with its issuer.

**The contribution identity and reference (row 15).** Row 15 adopts the form
of the advice (section 2.1). The two lines show the parts and decide no
encoding:

```text
contribution identity  = (publisher, local key)
contribution reference = (contribution identity, semantic revision)
```

- **The publisher.** Row 15: "The publisher is the library that publishes the
  mapping, the accepted pages' "library's owner"." The accepted sentence:
  "Contributions belong to their library's owner, outside canonical horos
  bytes." (`docs/binding-generation.md:101`).
- **What the publisher is not** (advice, section 2.1): the owner of the
  source declaration, the module pin of the source that is mapped, or the
  profile. A third party may map another owner's declaration, and several
  contributions state one profile.
- **The local key** tells apart the contributions of one publisher. The test
  of row 15: "Two contributions with one local key under two publishers
  register side by side, and two with one publisher and one key fail."
- **A name does not authenticate.** The advice (section 2.1): "A namespace
  string alone does not authenticate its owner; consumer composition must bind
  it to the selected publisher and manifest." The lock that would bind them is
  deferred (row 15a). Derived here: until then the registration in one
  process is the binding.
- **It does not find a factory.** Row 15 corrects a premise of the submitted
  document. A backend finds a factory by what the plan entry names: "By
  default a plan entry names the family's run-time factory and gives it the
  family's payload as arguments." (`docs/representation-families.md:49-50`).
  The contribution identity says what was selected. It is what duplicates and
  version checks compare.

**The five axes stay apart (row 16).** The accepted design: "**Five version
axes:** the envelope protocol; the family payload schema; the
representation's semantic identity; generated-code/runtime compatibility;
source semantics." (`docs/representation-families.md:88-90`). Row 16: "The
contribution reference carries the semantic revision. The other four axes are
stated where each belongs and are not folded into the identity."

| Axis | What it identifies (advice, section 2.1) | Where it is stated. Derived here, except the third |
| --- | --- | --- |
| Envelope protocol | How the shared graph is represented | The plan (#195) |
| Payload schema | How a family's recorded decisions are typed and encoded | The family's payload schema ([bootstrapping](#payload-bootstrapping)) |
| Representation semantics | What the selected representation means | The contribution reference (row 16) |
| Generated-code and runtime compatibility | Which generated code and runtime facilities can cooperate | The facet that names a factory, among its requirements (`docs/binding-generation.md:96`) |
| Source semantics | Which source interpretation the checked references use | The source-semantics input ([row 34](#inputs-that-are-not-bittypes-alone)) |

A digest of an implementation is provenance. It replaces none of the five.
The test of row 16: "Each axis is exercised alone: an unknown required value
of one axis is refused while the other four are unchanged."

**A family's semantic identity spans targets (row 17).** The advice asked for
a separate reference to a family's contract. Row 17 folds it into the
semantic axis: "a family's semantic identity is the same string in every
target's contribution reference". The use is recognition across targets. The
service root must know that a field was selected as the structural record
family (BC, section 2), and it should recognize a declared contract, not the
name of one target's contribution. The root facet that does so is deferred
(row 17a). What this means for an invariant that only one target has is
[open](#open-inside-this-design), point 4.

**Use entries and representation nodes (row 18).** The plan holds both.

| | A use entry | A representation node |
| --- | --- | --- |
| What it records | One selection at one source use or derived demand, with its identity and its provenance | The decisions and the dependencies of a selected representation |
| How many | One for each use | One for all the uses that selected the same representation |
| What it refers to | Its node | Other nodes, through typed links |

The test of row 18: "Two uses of `Page<Person>` are two use entries on one
representation node."

The accepted pages use the word "instance" for both. Row 18 gives them "a
dated clarification, in bittype for `docs/representation-families.md` and by
request for the binding contracts' memo sentence".

| Accepted sentence | Sense |
| --- | --- |
| "The planner selects one immutable representation instance per use." (`docs/representation-families.md:44-45`) | The use entry and the frozen decision it refers to |
| "The memo is keyed by the selected representation instance plus input identity, including the descriptor and generic environment." (`docs/representation-families.md:262-263`) | The applied representation, which several uses share. It is not the per-use record |

Read with the first sense, the second sentence would give every use its own
memo key, and one object reached through two fields would convert twice.

**The identity of a node (row 19).** Row 19: "Identity is over the selected
semantic representation without diagnostic provenance".

- **What is in it.** Everything that makes two selections represent
  differently. The advice (section 2.2) lists what a shorter key would miss:
  the recorded decisions, the target configuration, the representations
  selected for the dependencies, the identity of a nominal native
  declaration, the generic conventions and the evidence.
- **What is not in it.** Where in the source a selection was made. Two uses
  at two source sites can share a node. Derived here: a diagnostic site that
  reaches a projector through the record ([row 4](#what-a-rule-does-the-protocol))
  is therefore no part of what identifies the node.
- **No proof of equivalence.** The advice (section 2.2): "Do not require the
  kernel to prove arbitrary semantic equivalence between different
  implementations." Two different contributions never share a node, whatever
  they produce.
- **How the planner finds that two uses share a node** is
  [not specified here](#what-is-delegated-and-what-is-not-specified-here).

**Open generics are primary (row 20).** `Page<T>` has one complete open
description. `Page<Person>` is an application: a node that links the binders
of `Page` to the source arguments, the native arguments, the adapter arguments
and the evidence of this use (advice, section 2.3). The accepted pages agree:
"Open generics remain open; ground specialization is not a prerequisite."
(`docs/generation.md:45`), and "Each generic use carries its scoped binder
identities, native type arguments, adapter arguments and requirements"
(`docs/representation-families.md:92-94`).

- **Two conventions, kept apart.** The native arguments and the adapter
  arguments of an application are two lists, each tied to binders. They are
  never matched by position. A native declaration may drop a source
  parameter that conversion still needs, and it may add a native parameter
  that no source parameter stands for.
- **An omitted argument is kept for conversion.** Row 20: "An argument that
  the native convention omits is still represented for conversion". The test:
  "A native convention that omits a parameter still yields an adapter argument
  for it."
- **The native convention is the one of today. Derived here.** "A single
  generic convention relates source binders, native parameters and complete
  mapped actuals." (`docs/generation.md:45-47`) describes the native side and
  stays true of it. The adapter convention stands beside it.
- **A request for types only** must not fail for want of something that only
  conversion needs. How that fits one round is
  [open](#open-inside-this-design), point 5.

**Specialization (row 21).** An open declaration is emitted once, so its
layout is fixed. A family whose layout depends on its operand has four
courses inside an open generic (advice, section 2.3): a representation that is
valid for every admitted argument, a requirement on the binder, a native
generic convention that can express the variation, or a specialization that
changes the native representation and the conversion together. The accepted
rule for the last: "Explicit specialization is represented as such
([#166](https://github.com/Bitspark/bittype/issues/166)), never as ordinary
substitution into an unchanged generic declaration."
(`docs/representation-families.md:254-257`).

Row 21: "Explicit specialization (#166) is in the extension domain. The
contract can name an application, its selection and the resulting native
representation. The first implementation refuses a specialization as
unsupported". The test: "A specialization request is refused as unsupported
with the application named."

**The memo token (row 22).** Row 22: "The memo is the service library's.
bittype's plan carries what the token is derived from: representation nodes,
applications and selected dependencies. The token's form goes to
bitspark-service#40".

- **What the token is not** (advice, section 2.4): the identity of a use, the
  number of an entry, a contribution identity, a source descriptor alone, or
  a freshly made adapter object.
- **Not every node has one.** A fragment for a field position converts no
  value on its own. The advice (section 2.4): "do not require a value-memo
  token for every plan node."

**A generated declaration in the plan (row 23).** Three things are related,
and none is derived from another by name.

| | What it is | Who issues it |
| --- | --- | --- |
| The logical symbol identity | Which planned declaration an entry of the plan means | The planner. Row 23: "The planner's declaration identity is the logical symbol identity." |
| The handle | One declaration of one sealed output | The sealed output ([#272](https://github.com/Bitspark/bittype/issues/272)) |
| The realization | The name and the place that a language's Files stage gave the declaration | The Files stage of the language |

- **A plan entry refers by the first.** It holds no handle: "an
  output-issued handle cannot be the key of a serialized plan"
  (`docs/naming.md:644`). It holds no realized name, because a name is
  realized after planning.
- **The link from plan to output is published.** Row 23: "The inventory
  record of #281 gains it as a planner-stated field." It is an addition after
  pull request #281, recorded on
  [#272](https://github.com/Bitspark/bittype/issues/272#issuecomment-6840544720).
- **A manifest relates the three.** Row 23: "A serializable manifest relates
  plan, output and realization." The advice (section 2.5) lists what it
  identifies: the plan and its content digest, the generation run, the target
  configuration, the placement, and the realized symbol of each logical
  symbol. Who writes it and in which form is
  [not specified here](#what-is-delegated-and-what-is-not-specified-here).
- **A mismatch is a linking error.** The test of row 23: "A plan read with
  the manifest of another output fails as a linking error, not as an
  incomplete plan."
- **A backend has a fourth input.** Row 23: "the placed result is a fourth
  input of a backend." The sentence of binding generation is annotated
  ([Amendments](#amendments-to-accepted-pages)). What "native symbols" in the
  plan then means is [open](#open-inside-this-design), point 7.

**Read the realization now (row 24).** A reader takes a native name and place
from the realization, even where the realized name equals the spelled intent.
The naming design: "A consumer learns the realized name and the placed site
from a reader on the formatter's result, keyed by the output's handles."
(`docs/naming.md:668-669`). Row 24 limits "now": "For TypeScript, with the
reader on the formatter's result of #273. Go and Rust are row 24a". Neither
language has a Files stage for generated declarations yet.

**Declarations that a contribution owns (row 25).** A contribution may
introduce a declaration that represents no source declaration, such as Go's
optional helper.

- **The key is owned and logical.** The contribution states a key made of its
  contribution reference and a local name (advice, section 2.5). The planner
  issues the symbol and shares the declaration by that key.
- **The spelling is not the key.** Row 25: "The spelling is not the key." The
  advice (section 2.5): "Coincidentally identical names do not establish
  shared ownership."
- **A shared helper is asked for.** Row 25: "A contribution that needs
  whichever shared helper the composition selects makes a request". A helper
  that is private stays its owner's.

### The identities in one list

"Leaves the process" means that the identity has a place in the serialized
plan or in the manifest. The serialized form of each is
[provisional](#what-is-final-and-what-is-provisional).

| Identity | Who issues it | Leaves the process | Not to be confused with | Row |
| --- | --- | --- | --- | --- |
| Contribution identity: publisher and local key | The publisher states it | Yes | A source selector, a module pin, a profile, a registry's handle, a factory reference | 15 |
| Contribution reference: the identity and its semantic revision | The publisher | Yes | The other four version axes, a digest of an implementation | 15, 16 |
| A family's semantic identity | The family's owner, the same string for every target | Yes, inside the reference | The name of one target's contribution | 17 |
| Profile identity | The profile's owner. A contribution states it | Yes | A contribution identity, a version axis | [Row 1](#what-a-contribution-is) |
| Source identity of a declaration or form | The source integration | Yes | A contribution identity | `docs/representation-families.md:73` |
| Role selector of a field position | The source integration | Conditional on [row 11](#optional-fields) | A use, a type | 12 |
| A use | The kernel, for one run | No | A use entry | 12 |
| Use entry | The planner | Yes | A memo key, a representation node | 18 |
| Representation node | The planner | Yes | A use entry, the memo token, the native type | 18, 19 |
| Application | The planner | Yes | A specialized declaration, a match of two argument lists by position | 20, 21 |
| Applied representation token | The run-time of the service library | It exists at run time only | A use, an entry number, a contribution identity, a descriptor alone | 22 |
| Logical symbol identity | The planner | Yes | A handle, a spelling | 23 |
| Declaration handle | One sealed output | No | The key of a plan entry | 23 |
| Realized symbol: name and place | The language's Files stage | Yes, through the manifest | The naming intent, the spelled name of generation | 23, 24 |
| Key of an owned declaration | The contribution that introduces the declaration | Yes | The declaration's spelling | 25 |
| Facet kind | The owner of the kind | Yes | The issuer of an owner-issued value of #281 | [26](#what-the-kernel-carries-without-reading-it) |
| Factory reference: module and export | The facet of the contribution that names it | Yes | The contribution identity | 15 |

## What the kernel carries without reading it

**A facet is a typed extension contract (row 26).** Row 26: "A facet is a
kind-identified, versioned, validated value on a contribution." The accepted
elements of a paired contribution are facets in this sense: the native facet,
the value adapter facet, the root facet and further product facets
(`docs/binding-generation.md:93-97`), and a family's lowering facet
(`docs/representation-families.md:51-53`).

A facet kind has, by the advice (section 3.1):

- a stable qualified identity and a version of its contract;
- a typed way to construct or issue a value in the process;
- a validator;
- the categories of reference it may hold, and its compatibility
  requirements;
- a transport descriptor where it crosses a process boundary. That part is
  deferred (row 26a).

The contract of a kind belongs to the product that reads it. The kernel knows
no kind.

**Two carriers, and the difference (row 26).** Pull request #281 for
[#272](https://github.com/Bitspark/bittype/issues/272) adds owner-issued
values. Row 26: a facet "is a second carrier beside the owner-issued values of
#281, and the difference is stated: those attach to a declaration of one
output and never leave the process, a facet attaches to a contribution and may
cross a process boundary."

| | An owner-issued value (#281) | A facet |
| --- | --- | --- |
| Attaches to | One declaration of one sealed output | A contribution |
| Its identity | Lives in the process | A kind with a qualified identity and a version |
| Validated | No. The kernel stores it and hands it back | Yes, by the validator of its kind |
| Leaves the process | Never. A Files stage of the same process reads it | It may |
| Used for | The naming intent, the access a helper states | A value adapter, a root facet, a lowering facet |

The [comment of 13 November 2026](https://github.com/Bitspark/bittype/issues/272#issuecomment-6840544720)
on #272 records it and adds: "Neither generalizes the other."

**Selected together.** "A native facet and the facets a binding product needs
are selected together." (`docs/binding-generation.md:99`). A reader may load
the implementation that a facet names from elsewhere, for example a factory
by its module and export in another process. It stays a member of the one
selected contribution. Only what the selected contribution declares is
admitted, so no other package attaches behaviour by repeating a contribution
identity (advice, section 3.1). Where the Go facets that the gate executes
are declared is [open](#open-inside-this-design), point 8.

**What the kernel guarantees (row 27).** Row 27: "Guarantees 1 to 6 are the
kernel's, each with its test." The advice (section 3.2) says what kind they
are: "These are structural, authority, scope, and association guarantees."

| # | The kernel guarantees that | Its test | Red today because |
| --- | --- | --- | --- |
| 1 | Required facets belong to the selected contribution | Remove a facet that a bindings product requires. The request for bindings fails as "missing facet", and a request for types succeeds | A request names no product |
| 2 | A facet's kind, target and compatibility match | A facet from another issuer under the same display name, for another target, or of another contract version is rejected | No kind and no version exist |
| 3 | A reference has the right plan, category and scope | A reference of another plan, a member where a type is expected, and a binder outside its scope are each rejected | A reference has no category |
| 4 | Evidence belongs to the dependency it was asked about | An answer issued for another dependency, another generic environment or another semantic revision is rejected | No evidence path exists |
| 5 | The payload snapshot is immutable | Mutate the value a rule returned, once after it returned and once between freeze and projection. Neither the native output nor the stored plan changes | No record exists |
| 6 | The exact payload passed validation | Invalid data, a receipt issued for another snapshot, and a receipt from a validator that the composition did not register each fail before any backend runs | Nothing validates a payload |

The tests are the advice's, with the additions of verification report 1 in
rows 5 and 6. Three preconditions, derived here from that report: guarantee 1
needs a request that names the products it wants; guarantee 3 needs the
reference categories of the [next section](#what-lowering-must-be-able-to-refer-to);
guarantee 6 needs a validator that the composition registers.

**Two that are not the kernel's (row 27).**

| The advice's guarantee | Whose it is | Source |
| --- | --- | --- |
| Unknown required versions fail explicitly | The planner states the axes, and the reader refuses | Row 27: "the planner emits them and the plan's reader refuses". The test is that of [row 16](#identities) |
| Extension needs no shared case | Nobody's at run time. It is evidence for the extension boundary | Row 27: it "is qualification evidence, not a kernel guarantee". It is line 7 of the [gate](#the-gate) |

**What none of them establishes.** The advice (section 3.2): "They do not
prove that a claimed equality relation is correct, that the projector uses the
intended field, or that the TypeScript factory obeys the conversion laws.
Qualification must test those claims."

**Freeze before either projection (row 28).** The accepted design fixes what
a payload is: "A payload is an ontos value. Its family declares the payload
type in horos; bittype generates the Go and TypeScript types, and horos
conformance validates a payload." (`docs/representation-families.md:84-86`).
Row 28 adopts the advice's order for it (section 3.3):

1. The rule returns its decision.
2. The decision is encoded, validated against its schema and frozen as one
   canonical snapshot.
3. The projector receives a typed view that is derived from that snapshot.
4. Serialization publishes that same snapshot.

The invariant is one immutable decision value. The storage technique is not
fixed by it. The test of row 28: "Mutating the value a rule returned, after it
returned, changes neither the native type nor the serialized plan."

**A new package role (row 28).** Row 28: "Encoding and validation live in a
new package role outside the contract package, which may import neither horos
nor ontos. The horos source reader is not its home. `LAYOUT.md` names the
home".

- **The role.** This page calls it the **payload codec**: the package that
  encodes a decision as an ontos value, validates it with horos conformance
  against the family's schema, and freezes the snapshot. The name is this
  page's.
- **Why it is outside the contract package.** "Generic carriers, contracts,
  engine, printer and native syntax stages have no external Horos/Ontos
  imports." (`CHARTER.md:129-131`). The kernel carries an opaque handle to the
  snapshot and the result of its validation, bound to that snapshot and to its
  schema. It imports neither library.
- **Why it is not the horos source reader.** Derived here, from verification
  report 1: a payload's schema is a horos type whatever the source language
  of the generation is. A source that is not horos plans through the same
  kernel and needs the same codec.
- **Who registers it.** The consumer's composition, like a source and a
  target. Derived here.
- **Its path.** `LAYOUT.md` does not name it yet. It gets the path and the
  import allowances with the prototype, and both are
  [not specified here](#what-is-delegated-and-what-is-not-specified-here).
  How a rule's decision becomes an ontos value is
  [open](#open-inside-this-design), point 9.

## What lowering must be able to refer to

A lowering facet turns recorded decisions into target code: "typed target
expressions or helpers, such as generated accessors or constructors for Go",
and "source-code strings never stand in for typed target expressions"
(`docs/representation-families.md:52-55`). Row 29 takes the advice's list
(section 3.4) as "the checklist of the design", and says whose each part is:
"Reference categories are the contract's. Access and construction operations
belong to the reading product's vocabulary."

| What a lowering facet must be able to express | Whose |
| --- | --- |
| The native type of a selected dependency | A reference category of the contract |
| A scoped native binder, and an adapter argument | Reference categories of the contract |
| A generic requirement and its evidence | A reference category of the contract |
| A generated declaration, an external symbol, an owned helper | Reference categories of the contract |
| A generated member, tied to the native type that contains it | A reference category of the contract |
| A typed access or construction operation | The vocabulary of the product that reads the facet |
| A helper that the lowering itself introduces | Its key is a reference of the contract ([row 25](#identities)). Its body is in the product's vocabulary. Derived here |
| The value of the run-time call context | The vocabulary of the product that reads the facet |

- **A member is referred to, not spelled.** The advice (section 3.4): "A
  free string alone does not provide that association." The name of a member
  may stay data of the payload. A lowering operation refers to the member and
  to the type that contains it.
- **The contract gains no expression of any language.** The kernel checks the
  category, plan and scope of a reference
  ([guarantee 3](#what-the-kernel-carries-without-reading-it)). What is done
  with a reference is the product's.
- **The set of categories** above is the advice's list read by row 29. The
  prototype fixes the set, and how a member reference is issued is
  [not specified here](#what-is-delegated-and-what-is-not-specified-here).

**Where the prototype's value vocabulary lives (row 29).** "The prototype's
Go value vocabulary and its facets stay in the test tree". The reason is the
accepted scope: "bittype's native vocabularies stay types-only."
(`docs/binding-generation.md:83`). The prototype's vocabulary is a test
instrument. The advice (section 3.4) says what it must do: "It must compile
and run generated construction, member access, an owned helper, and a generic
adapter argument." The production owner of value-level Go is not bittype: "A
Go backend belongs to the future Go runtime library."
(`docs/binding-generation.md:46`).

## The extension domain

**What "changes once" means (row 33).** Representation families says that
what an existing mapping returns changes once: "After them a family adds a
contribution and changes no existing contract."
(`docs/representation-families.md:328`). Row 33 adopts the advice's wording
for that commitment (section 4.4): "One migration to the demonstrated
extension contract for the agreed extension domain, before independently
owned consumers rely on it."

It is a bounded commitment, as the rules have it: "It cannot prove that all
imaginable future requirements fit the interface. The promise is about the
explicitly declared extension domain and the responsibilities that support
it." (platform `docs/concepts/architectural-correctness.md:501-503`). A
requirement outside the domain is met by a justified change under AC-08. It
does not break the commitment.

**The domain (row 33).** Row 33: "further products and facets in; a
contextual representation at a field position in; cross-field layout
coordination out; application-specific specialization in and refused until
#166; one authoring route, a Go contribution package".

| Item | In or out | What that means | What is refused until built |
| --- | --- | --- | --- |
| A further family | In. Derived here | Every family of the delivery table of representation families (`docs/representation-families.md:339-349`) enters as contributions. The accepted promise: "Adding a family through the published identity, dependency, invocation, plan-graph and target-extension contracts changes only that family and the consumers that compose it." (`:22-24`) | A family that is not built is unsupported |
| Further products and facets | In | A further backend for a target, and a further kind of facet, enter through the [facet contract](#what-the-kernel-carries-without-reading-it) with no edit of the kernel | A request for a product whose facet a contribution lacks fails as "missing facet" |
| A contextual representation at a field position | In | The request of [row 12](#optional-fields), and a contribution without a value adapter of its own ([row 14](#optional-fields)). Row 33 states the item without a condition, and row 11 is conditional. What that means under a refusal is [open](#open-inside-this-design), point 10 | A fragment that the target's member vocabulary does not admit |
| Cross-field layout coordination | Out | A layout that several fields share needs the record policy replaced ([row 14](#optional-fields)) | A contribution that would need it replaces the record contribution |
| Application-specific specialization | In, and refused until [#166](https://github.com/Bitspark/bittype/issues/166) | The contract can name an application, its selection and the resulting native representation ([row 21](#identities)) | Every specialization, as unsupported, with the application named |
| Authoring routes | One | A Go contribution package ([row 37](#inputs-that-are-not-bittypes-alone)). A declarative sidecar is not promised | Nothing is refused. No other route exists |

The list narrows nothing that an accepted page promises. The rule: "The
extension domain comes from the agreed responsibility, not retrospectively
from the cases already implemented. Do not narrow the stated contract to make
a closed implementation appear compliant." (platform
`docs/concepts/architectural-correctness.md:364-366`).

**How the domain is tested.** Derived here, from verification report 3: one
fixture for each item that is in, added with the shared packages unchanged,
and one recorded refusal for each item that is out or not yet built. They
belong to line 7 of the [gate](#the-gate).

## What is final and what is provisional

**The partition (row 30).** The submitted research document proposed "contract
final, envelope unstable". Row 30 withdraws it: "The partition of the head of
section 6 is replaced by the advised one: meanings final, encoding and layout
provisional". The reason is that the contract and the envelope share their
identities, their categories of dependency, their generic environments,
their requirements and their references to symbols (advice, section 4.1).

| Must have its intended meaning. Final with the gate | Stated under | May stay provisional or change under a version |
| --- | --- | --- |
| The meaning and the ownership of each identity | [Identities](#identities) | The layout and the encoding of the serialized plan |
| The categories and the scope of a request and of its answer | [The protocol](#what-a-rule-does-the-protocol), [Optional fields](#optional-fields) | Caches and performance strategies |
| The generic conventions and requirements | [Identities](#identities), [Evidence](#evidence) | Coverage of further families |
| Recursive references, and the meaning of a refusal | [Evidence](#evidence), [Refusals](#refusals) | Revisions of a family's payload schema |
| The validation of a facet and what it is attached to | [What the kernel carries](#what-the-kernel-carries-without-reading-it) | Further implementations of a product |
| The relation between a native symbol and its realization | [Identities](#identities) | The presentation of diagnostics |
| The inputs that lowering needs | [What lowering must be able to refer to](#what-lowering-must-be-able-to-refer-to) | Optimizations of lowering |

The two outer columns are the advice's table. A row of the left column is
not paired with the entry beside it.

- **Unstable does not mean without meaning.** The advice (section 4.1):
  "Calling the transport unstable does not excuse changing the meaning of an
  admitted identity later." The accepted sentence stands and is read so: "The
  envelope keeps an explicitly unstable version until the completeness test
  backend passes" its six cases (`docs/representation-families.md:95-96`).
  The version is what is unstable. It gets a dated clarification
  ([Amendments](#amendments-to-accepted-pages)).
- **Decided is not demonstrated.** Derived here, from AC-07: a row of the
  left column is decided by this page and demonstrated only by its test in
  the gate. The relation between a symbol and its realization is demonstrated
  for TypeScript only, until Go and Rust have a Files stage (row 24a).
- **Until the gate is met** nothing on this page is final. The prototype may
  change any of it, and each change is recorded here with its date.
- **After the gate** the left column is fixed, by AC-10: a first change
  "establishes the final structure: its inputs and identities have their
  intended shape". The right column changes under its version axes
  ([row 16](#identities)).

## Inputs that are not bittype's alone

**The source semantics, for a reader in another process (row 34).** A
backend's second input is "the source semantics"
(`docs/binding-generation.md:26`), and it emits "ground-descriptor constants"
(`docs/binding-generation.md:77`). No serialized form of that input exists.
Serializing the plan alone therefore shows nothing about a reader in another
process. Row 34: "For the prototype: a snapshot of ground descriptors and
module pins exported from the Go side, whose source identities match the
plan's. A general reader protocol is row 34a".

- **What a reader may do with it.** The advice (section 5.1): "The backend may
  read source meaning; it must not reconstruct representation choices."
- **A snapshot that does not match is a linking error.** Derived here, from
  verification report 3: a snapshot of a module with the same name and shape
  and another pin is refused, and not read as an incomplete plan. Binding
  generation lists the case among its falsification tests: "a changed module
  digest with the same name and shape" (`docs/binding-generation.md:254`).

**Target configuration (row 35).** The plan records the configuration of the
target that affects selection, what is admissible natively, or whether
generated code is compatible. The envelope has the place: "Version axes;
target and its configuration" (`docs/representation-families.md:72`).

For TypeScript one option matters to optional fields. bittype checks its
generated TypeScript under `exactOptionalPropertyTypes`. That option changes
what the compiler admits. It does not change that JavaScript tells an absent
property from a present one (advice, section 5.2). Row 35: "The optional
family's contract states it, and the four states are tested apart."

| State of an optional field | Tested |
| --- | --- |
| Absent | By its own vector, in Go and in TypeScript |
| Present with a value | By its own vector, in Go and in TypeScript |
| Present with `null`, where the type admits it | Derived here: none of the delivered families admits `null`, so this state enters with [#211](https://github.com/Bitspark/bittype/issues/211) |
| Present with `undefined` | By its own vector, in TypeScript. What conversion does with it is not decided |

- **A present `undefined` is not decided.** Row 35: "bittype checks under
  `exactOptionalPropertyTypes`, and the service library compiles without it
  and reads a present `undefined` as absence. That difference is put to its
  owner". The question asks which of two behaviours the vectors pin: refused,
  or read as absent. The decision is the representation owner's, which for
  the structural families is bittype's structural profile
  (`docs/representation-families.md:60-62`). It is recorded here when it is
  made.
- **"Round-trip apart".** The test of [row 13](#optional-fields) says that
  the four states "round-trip apart through the record factory". Derived
  here: a present `undefined` has no canonical value of its own, so under
  either answer it does not come back as itself. The test is read with row
  35: each state has its own vector and its own expected outcome.

**What goes to the service library (rows 36 and others).** Row 36 is the one
stale row: "Stale for bittype: laws K1 to K4 of the binding contracts state it
and bittype confirmed them on 12 November. The two points they leave unstated
go to the service library". All of the following is in the
[comment of 13 November 2026](https://github.com/Bitspark/bitspark-service/issues/38#issuecomment-6840511750)
on bitspark-service#38. The grade says of the three questions: "None of them
blocks the prototype."

| What | Row | Form |
| --- | --- | --- |
| A field-level conversion slot for an optional record field | 11, 13 | Question 1 |
| What conversion does with a present `undefined` | 35 | Question 2 |
| Whether a prototype of one delivered family may precede integers | 31a | Question 3 |
| A backend has a fourth input. Item 9 of bitspark-service#40 repeats the sentence about three | 23 | Told |
| The form of the memo token is theirs, on bitspark-service#40 | 22 | Told |
| Two points that K1 to K4 leave unstated: the states of a memo entry, and that a memo hit still charges the cost of an expanded occurrence | 36 | Told |

**One authoring route (row 37).** Row 37: "One route: a Go contribution
package, tested out of tree. A declarative sidecar is not promised."

- **What the route is.** A Go package that states a contribution through the
  public packages of bittype only. A library whose repository is written in
  TypeScript publishes such a package beside its TypeScript factory (advice,
  section 5.4).
- **What it replaces.** Row 12 of research 0014 named two forms: "as a Go
  mapping package or a declarative sidecar". No design page promised the
  second, and row 37 does not.
- **Its test.** Derived here, from verification report 3: a Go module outside
  this repository authors and registers a contribution. It is red when the
  module needs a symbol that is internal or that lives in the test tree.
- **The first contribution of the service library** in its own repository is
  deferred (row 37a).

## Payload bootstrapping

Row 38: "A pinned bootstrap artifact: the schema source, its canonical module
with its pin, and the committed generated types with a regeneration test".

bittype generates the types of a payload from a horos schema, with the
families whose payloads those types describe. The families that exist cannot
type every field a payload wants. The advice (section 5.5): "The limitation
to six currently supported families is a bootstrap issue, not a reason to
encode unsuitable payload semantics."

- **One schema language.** "No separate schema language is introduced."
  (`docs/representation-families.md:86-87`). The advice repeats it: "Do not
  introduce a second schema language."
- **A field the families cannot type yet.** Derived here, from verification
  report 3: such a field is carried by structure where that says the same,
  for example by the order of a list where an index was wanted. Otherwise it
  is recorded as a bootstrap exception that names the family that ends it.
- **Two kinds of later change.** The advice (section 5.5): "Adding a semantic
  integer field later is normally a family schema/version change. Discovering
  that a payload has been hiding an orchestration reference is a contract
  defect." The second is ruled out by the accepted design: "Nothing the
  orchestration must resolve hides inside a payload."
  (`docs/representation-families.md:82-83`).
- **Its tests.** The regeneration test fails when the committed types differ
  from what the schema generates. A payload that does not conform to its
  schema is refused before any backend runs
  ([guarantee 6](#what-the-kernel-carries-without-reading-it)).

## What a test must show

Row 39: "The boundary is enforced by what a projector can reach. Mutants must
disagree on presence, layout, slot identity or environment."

**The mutant that proves nothing.** Row 2 of research 0015 validates the one
origin of a decision by "a seeded mutant that builds the plan entry in a
separate pass from the native expression". A second pass that computes the
same decisions from the same source is an equivalent mutant. Every
behavioural test stays green on it.

**What replaces it.** Three layers. The first is row 39's. The other two are
derived here, from verification report 1.

| Layer | What it shows | How |
| --- | --- | --- |
| Structure | No second decision point is reachable | The projector cannot reach the source reader, and a reader of the plan has nothing to select with |
| Mutants that disagree | The tests detect a wrong implementation | Each mutant differs in one thing: presence in the entry and not in the native type, the order of fields in one of the two, two crossed links, an adapter argument tied to the wrong binder |
| A decision that the source does not determine | A second origin would be visible | A fixture contribution takes a member's name from its own configuration. A pass that rebuilds the entry from the source then differs from the native type |

The first layer prevents an unauthorized decision point. The other two
detect an incorrect implementation. AC-07 keeps the two kinds of evidence
apart.

**The texts that change.** Row 39: "#207's verification and row 2 of research
0015 are reworded". The verification of #207 is reworded by the
[comment of 13 November 2026](https://github.com/Bitspark/bittype/issues/207#issuecomment-6840545162)
on that issue. The validation cell of row 2 in the table of research 0015
still has the old text. This change does not edit that table
([the notes not yet made](#what-is-delegated-and-what-is-not-specified-here)).

## Amendments to accepted pages

Each amendment is a dated clarification beside the sentence it concerns. The
sentence itself is not changed, so it stays readable as it was accepted. All
of them are made in the change that adds this page.

| Page and sentence | What the clarification says | Why | Row | Who owns the page |
| --- | --- | --- | --- | --- |
| [Binding generation](binding-generation.md#order-of-work-and-qualification): "The unit of implementation is library × target." (`:265`) | The sentence names ownership and distribution. The unit that a consumer selects is one contribution, and a library may publish several | Read as the unit of selection, the sentence would make every rule of a library one indivisible choice | 1 | Research 0014, held by `claude-bittype-gaps` |
| [Binding generation](binding-generation.md#the-selected-representation-plan): "The backend can implement each selected representation from the plan, the source semantics and the selected contributions alone, without guessing a rule." (`:72`), and "consumes three inputs" (`:23`) | For a generated native declaration the three inputs are not sufficient. Its name and place exist only after a Files stage has placed it. A plan entry refers to it by an identity the planner issues, and a backend reads the realization from the placed result, which is a fourth input | The [naming design](naming.md#reconciliation-with-accepted-documents) stated the tension and left it to research 0017 | 23 | Research 0014, held by `claude-bittype-gaps` |
| [Representation families](representation-families.md#the-representation-contract): "so a native declaration and its conversion cannot disagree." (`:46-47`) | "Cannot disagree" means that no sanctioned second decision point exists. It does not mean that a projector or a factory cannot be wrong | Read as a claim about every implementation, the sentence contradicts AC-05 | 6 | Research 0015. This chain holds its delivery |
| The same paragraph: "one immutable representation instance per use" (`:44-45`), and "The memo is keyed by the selected representation instance" (`:262`) | "Instance" has two senses. In the first sentence it is the selection at one use, a use entry. In the second it is the applied representation that several uses share | Read with one sense, every use would have its own memo key | 18 | Research 0015 |
| [Representation families](representation-families.md#the-plan): "The envelope keeps an explicitly unstable version" (`:95-96`) | What is unstable is the version: the layout and the encoding. The meanings that the plan shares with the contribution contract are fixed with that contract | The partition "contract final, envelope unstable" is withdrawn | 30 | Research 0015 |
| [Representation families](representation-families.md#delivery): "After them a family adds a contribution and changes no existing contract." (`:328`), and "generics" among the delivered families (`:341`) | "Changes once" is a commitment for the written extension domain. "Generics" names a division of responsibility and no family with a contribution of its own | Rows 33 and 2 | 33, 2 | Research 0015 |
| [Package boundaries](package-boundaries.md#findings-on-main-at-3bf45c3), finding 7 (`:156`) | The finding is bittype's derivation and is conditional on the answer of the service library. If the slot is refused, the finding is withdrawn | Row 11 | 11 | bittype, case bittype/ACX-009 |

**The edit of binding generation is two insertions.** Pull request
[#260](https://github.com/Bitspark/bittype/pull/260) of `claude-bittype-gaps`
edits the same page and is open. This change adds one paragraph after each of
the two sentences and touches no line that #260 changes. In representation
families it adds three clarifications, none in a passage that #260 changes.

**Pointers, not amendments.** Five pages gain a pointer to this page and
change no statement of their own: the [SPI](generation-spi.md) and the
[registration page](kernel-composition.md) say that the mapping contract is
going to change, [semantic generation](semantic-generation.md) and the limit
of [bittype/ACX-009](architecture-cases/0009-package-boundaries.md) say that
the optional move is conditional, and `LAYOUT.md` names the new package role
without a path.

**Sentences that are not amended now.** Each needs its owner or an answer
first. They are listed where they arise:

| Sentence | Why not now | Where |
| --- | --- | --- |
| "the rules for absence, presence, null and native `undefined`" in the plan's list (`docs/binding-generation.md:58`) | A reading, not a decision | [Open](#open-inside-this-design), point 1 |
| "native symbols, distinguishing run-time imports from type-only imports" (`docs/binding-generation.md:56`), and "Native symbols, imports, factory references, runtime requirements" (`docs/representation-families.md:78`) | Row 23 names one sentence only | [Open](#open-inside-this-design), point 7 |
| The three sentences that a field-level slot would change | They wait for the answer on bitspark-service#38 | [Optional fields](#optional-fields) |
| The memo sentence of the binding contracts | It is the service library's | [What is delegated](#what-is-delegated-and-what-is-not-specified-here) |

## Deferred, with triggers

Nine rows are deferred. A trigger is a state of a repository that anyone can
check.

| Row | What is deferred | The trigger that ends the deferral | Tracked |
| --- | --- | --- | --- |
| 9a | Recursive representation references | "#164 is claimed and vertical case 1 is being built." Until then recursion is refused: `grep -n "recursive native declaration" contracts/generation/go/planning.go` finds the refusal | [#164](https://github.com/Bitspark/bittype/issues/164) |
| 15a | A composition lock that binds a publisher's name to its manifest | "a second publisher registers a contribution (bitspark-service#39 item 2 starts)" | bitspark-service#39 |
| 17a | A root facet of `Service` that recognizes a family's contract | "bitspark-service#39 starts its root facet. No root facet exists in either repository" | bitspark-service#39 |
| 24a | A reader of the realization for Go and Rust | "slice 1 of #161 or #159 delivers a generated-Files formatter (`ls format/files/golang format/files/rust`)" | [#161](https://github.com/Bitspark/bittype/issues/161), [#159](https://github.com/Bitspark/bittype/issues/159) |
| 26a | A transport descriptor for a facet | "the first facet crosses a process boundary in the prototype of row 31" | #196 |
| 31a | TypeScript end to end through the service library's own factory | "`packages/structural-adapters` exists on bitspark-service `main`, or its owner records on bitspark-service#38 that a prototype precedes integers" | bitspark-service#38 |
| 31b | Vertical cases 1 and 2 in real form | "#164, #210 and #211 deliver their families" | [#164](https://github.com/Bitspark/bittype/issues/164), [#210](https://github.com/Bitspark/bittype/issues/210), [#211](https://github.com/Bitspark/bittype/issues/211) |
| 34a | A reader protocol for source semantics | "the design of bitspark-service#40 item 3 starts" | bitspark-service#40 |
| 37a | A contribution that is authored in the service library's repository | "bitspark-service#39 item 2 starts" | bitspark-service#39 |

Three remarks, derived here.

- **Row 26a ends inside the prototype.** The stand-in reader of TypeScript
  reads a serialized plan in another process. The facet that names a factory
  crosses the boundary there, so the transport descriptor is designed in the
  prototype and recorded on this page.
- **Row 31a is why a stand-in stands in.** What it proves and what it does
  not is under [The order of work](#the-order-of-work).
- **A deferral is no narrowing.** Each deferred part is in the
  [extension domain](#the-extension-domain) or belongs to another owner. None
  is removed from what the accepted pages promise.

## Open inside this design

Ten points on which two rows, or a row and an accepted page, do not fit
without a reading. This page does not choose the reading. Each point gives
the two texts, the options, and what the rest of the page assumes meanwhile.
The research owner decides each, and the decision is recorded here.

**1. An invariant carried by reference, and the plan's list (row 5 and binding
generation).** "For each use and generic instance, the plan carries" a list
that includes "the rules for absence, presence, null and native `undefined`"
(`docs/binding-generation.md:52`, `:58`). Row 5: "An invariant of a semantic
revision is carried by reference to that revision".

- (a) The plan carries a rule also when it names it through the selected
  contribution and its semantic revision. The page of research 0014 gets a
  dated clarification from its owner.
- (b) Every such rule is a field of a payload, also where it never varies.

Meanwhile the page follows row 5 and edits nothing.

**2. The name of the second outcome (row 8 and representation families).** Row
8: "The accepted outcome names are used (proved, refuted, unknown,
exhausted)". The accepted page: "`proved`, `overlap`, `unknown` and
`exhausted`" (`docs/representation-families.md:200`).

- (a) Refuted is the general name. `overlap` is what one analysis, the
  disjointness of union branches, calls its refuted outcome.
- (b) The four accepted names are used as they are, `overlap` included.
- (c) The parenthesis of row 8 is corrected.

Meanwhile the page uses the names of row 8.

**3. A choice that leaves a rule before it decides (rows 3 and 5, and row
13).** Row 13: "the member's name comes with the request". A request is made
in step 1, and the record is frozen after step 3 (row 3). A member's name is
a per-selection choice of the record contribution (row 5).

- (a) What a rule states in step 1 is part of its decision and is frozen with
  it. The kernel carries the arguments of a request and does not read them.
- (b) The answer leaves the name out, and the record contribution adds it in
  its own decision.
- (c) A member's name is fixed from the source view and the naming policy
  before any rule runs, and is no choice of a rule.

The point is void if [row 11](#optional-fields) falls to D1.

**4. One semantic string for all targets, and an invariant of one target
(rows 5 and 17).** Row 5 lets an invariant be carried by the semantic
revision. Row 17: "a family's semantic identity is the same string in every
target's contribution reference". The advice's own example of an invariant
belongs to one target: an optional representation that "always ignores its
unused value when absent" (section 1.2). If that changes for Go alone, no
text says which string changes.

- (a) The semantic revision has two parts: the family's contract, shared by
  all targets, and the revision of one target's representation.
- (b) It is one string. A change that concerns one target is a new
  contribution with a new local key.
- (c) An invariant of one target is a field of the payload after all.

Connected and also undecided: whether the local key of a contribution
contains its target.

**5. A dependency that only one product needs (rows 3 and 20, and two
accepted sentences).** Row 3: the planner resolves the named dependencies
before the rule decides. Row 20: "An argument that the native convention
omits is still represented for conversion". Accepted: "Only demanded native
actuals need mappings." (`docs/generation-spi.md:162`), and "A native mapping
without an adapter facet still supports types"
(`docs/binding-generation.md:99`). If every named dependency is planned, a
request for types can fail for something that only conversion needs.

- (a) A rule names each dependency with the product that needs it, a request
  names its products, and the planner resolves what the requested products
  need.
- (b) Every named dependency is planned. This contradicts the second
  accepted sentence.
- (c) A dependency that cannot be planned yields a refused answer, and only
  the product that uses the answer fails.

[Guarantee 1](#what-the-kernel-carries-without-reading-it) needs a request
that names its products under every option.

**6. What the prototype shows of recursion (rows 31b and 9a).** Row 31b:
vertical cases 1 and 2 "are exercised at the kernel's level through the
independent fixture source". Row 9a defers recursive references and says that
recursion "is refused today".

- (a) The prototype shows the two refusals of [row 9](#evidence) with their
  different reasons, and nothing else of case 1.
- (b) The prototype lifts the refusal for the fixture source, to show that a
  reference is reserved before its body is filled. The record says that no
  family stands behind it.

Under (a), "recursive references" in the left column of
[what is final](#what-is-final-and-what-is-provisional) is decided and not
demonstrated.

**7. "Native symbols" in the plan (row 23 and two accepted sentences).** The
plan carries "native symbols, distinguishing run-time imports from type-only
imports" (`docs/binding-generation.md:56`), and the envelope lists "Native
symbols, imports, factory references, runtime requirements"
(`docs/representation-families.md:78`). Row 23 annotates one other sentence.
The naming design: "With names realized by Files, a plan entry cannot hold a
realized name or an import path." (`docs/naming.md:721-722`).

- (a) A native symbol in the plan is a reference: the logical symbol identity
  of a generated declaration, or the external symbol of an external one, each
  with whether its use is needed at run time. An import is derived by the
  backend from the realization. Both sentences get a dated clarification from
  their owners.
- (b) The plan and the manifest travel as one input, and the two sentences
  are read of that input. The advice (section 2.5) allows the packaging.

Meanwhile the page states only what row 23 states.

**8. Where the Go facets that the gate executes are declared (rows 26, 29 and
40).** Row 26 puts a facet "on a contribution". Row 29: "The prototype's Go
value vocabulary and its facets stay in the test tree". Row 40 asks for
executed Go lowering on the prototype branch and then on `main`. A production
mapping imports "Never another mapping and never a native target"
(`LAYOUT.md:160`), and nothing of the test tree.

- (a) The gate runs Go conversion and lowering on contributions of the test
  tree, which state identities of their own. The production Go contributions
  carry no binding facet until a Go runtime library exists. The record says
  that the six production families were not what was lowered.
- (b) A production contribution declares its facet as a reference, and the
  test composition supplies the implementation that the reference names.
- (c) A contribution of the test tree wraps a production rule and adds the
  facets.

**9. How a decision becomes an ontos value (row 28, the charter and the
mapping role).** Row 28 puts encoding into a package role of its own. A
payload "is an ontos value", and "bittype generates the Go and TypeScript
types" (`docs/representation-families.md:84-86`). A mapping package may not
import ontos (`LAYOUT.md:160`). To turn a value of a generated Go type into
an ontos value is the Go conversion of the payload's own families, and the
charter excludes "application-value validators, serialization" from bittype
(`CHARTER.md:119`).

- (a) Each contribution encodes its payload by hand with the constructors of
  ontos. Every mapping package then gets an allowance to import ontos.
- (b) An encoder is generated and committed with the
  [bootstrap artifact](#payload-bootstrapping).
- (c) A rule builds its decision through a builder that the payload codec
  supplies, and never holds a value of the generated type. The generated
  types serve the readers.

This decides what a rule returns. It is settled in the prototype, before the
gate, and not with `LAYOUT.md` alone.

**10. The final contract and the answer on optional fields (rows 11, 12, 31
and 33).** Row 31 puts the migration before the optional move, and the
delivery note of the table puts that move "by the answer of row 11". The test
of row 11 is "Red until then". Row 33 puts "a contextual representation at a
field position in" the extension domain without a condition. The request of
row 12 is what the final contract holds for that item.

- (a) The migration waits for the answer. Under a refusal the request of row
  12 is not built, and the item leaves the extension domain.
- (b) The request is part of the contract under either answer and is shown by
  a fixture. The answer decides only whether the standard record
  contributions use it. Under a refusal no TypeScript conversion exists for
  such a contribution, so the item is in for types only.
- (c) The migration goes ahead with the delegated form. A refusal later
  removes it by a second change, justified under AC-08.

Meanwhile the page marks [Optional fields](#optional-fields) as conditional
and states the item of the domain as row 33 does.

## What is delegated, and what is not specified here

**Delegated**, each with the one who decides:

| What | To | Row |
| --- | --- | --- |
| Whether a field-level conversion slot is admitted, and its slot kind | The owner of the service library, on bitspark-service#38 | 11, 13 |
| Which behaviour the vectors pin for a present `undefined` | Asked of the service library. Decided by bittype's structural profile after the answer | 35 |
| Whether a prototype of one family may precede integers in the service library | The owner of the service library, on bitspark-service#38 | 31a |
| The form of the memo token, the states of a memo entry, and the cost charged at a memo hit | The service library, on bitspark-service#40 | 22, 36 |
| The wording of the memo sentence of the binding contracts | The service library | 18 |
| The shapes of the entries, the serialized form of the plan, and who writes the manifest in which form | [#195](https://github.com/Bitspark/bittype/issues/195) | 18, 23 |
| The planner's declaration identity as a field of the inventory record | [#272](https://github.com/Bitspark/bittype/issues/272), delivered with #195 after pull request #281 | 23 |
| The reader on the result of the TypeScript formatter | [#273](https://github.com/Bitspark/bittype/issues/273) | 24 |
| Recursive references | [#164](https://github.com/Bitspark/bittype/issues/164) | 9a |
| Explicit specialization | [#166](https://github.com/Bitspark/bittype/issues/166) | 21 |
| The report of binding support | [#197](https://github.com/Bitspark/bittype/issues/197) | The delivery note of the table |
| Value-level Go in production | "the future Go runtime library" (`docs/binding-generation.md:46`) | 29 |

**The notes not yet made.** Two things that a row says are done or asked, and
that this change found undone, and one text that does not carry the gate yet:

- **Row 18** asks for a clarification of the binding contracts' memo sentence
  "by request". The comment of 13 November 2026 on bitspark-service#38 tells
  that owner about the memo token (row 22). Neither it nor a comment on
  bitspark-service#40 names the two senses of "instance".
- **Row 39** says that row 2 of research 0015 is reworded. The validation
  cell of that row has the old text. This change does not edit the table of
  another research record.
- **The gate on #196.** Row 40 makes the seven lines the merge gate of item 6
  of #196. The grade on that issue names the gate. The text of the issue
  itself does not carry the seven lines.

**Not specified here:**

- **Any Go type or signature of the contract.** The sketches on this page
  decide none. The prototype writes them.
- **The form of a refusal in Go:** values, or one type with a class. The
  value of the refusal of an evidence cycle and of a recursive reference.
  Whether a missing contribution and an ambiguous one are one value or two.
  What an unexpected error of a rule is: the advice calls it an
  implementation failure, and row 10 lists no such class.
- **How a request names its products**, and whether a rule names a dependency
  with the product that needs it ([open](#open-inside-this-design), point 5).
- **How the planner finds that two uses share a node.** Row 18 fixes the
  outcome for two uses of `Page<Person>`. Row 19 fixes what identity is over.
  The procedure belongs to #195.
- **The form of a publisher's name and of a local key**, and whether the key
  contains the target ([open](#open-inside-this-design), point 4).
- **The path and the import allowances of the payload codec.** `LAYOUT.md`
  gets both with the prototype.
- **The set of reference categories**, and who issues a reference to a
  member.
- **The key under which a contribution asks for a shared helper.** Row 25
  says that it makes a request. Whose key that request uses is not said.
- **Whether a member fragment may hold more than one member.** The target's
  member vocabulary decides, and what it does not admit is refused.
- **The law of a field-state protocol.** The laws of binding generation are
  written "For a type T and representation profile p"
  (`docs/binding-generation.md:121`). An optional field is no type.
- **The kinds of evidence.** Each belongs to its vocabulary. This page fixes
  the outcomes and the rule for a binder.
- **Composition across several profiles.** The advice (section 5.6) names it
  as needing an explicit contract. No row decides it.
- **Which contribution is the materially different one of the gate.** The
  advice (section 4.2) names three candidates. One of them, an optional
  fragment of another shape, exists only under D3.
- **How the authoring route meets an application command line**
  ([#172](https://github.com/Bitspark/bittype/issues/172)). Not examined
  here.

## Where each row is stated

"Line" is the line of the [gate](#the-gate) whose tests validate the row.
A dash means that the row has no test of its own before the gate is met.

| Row | Verdict | Stated under | Line |
| --- | --- | --- | --- |
| 1 | Holds | [What a contribution is](#what-a-contribution-is), [Amendments](#amendments-to-accepted-pages) | 7 |
| 2 | Holds | [What a contribution is](#what-a-contribution-is), [Amendments](#amendments-to-accepted-pages) | 1 |
| 3 | Adapted | [The protocol](#what-a-rule-does-the-protocol) | 2 |
| 4 | Holds | [The protocol](#what-a-rule-does-the-protocol) | 1 |
| 5 | Holds | [The protocol](#what-a-rule-does-the-protocol), [Open](#open-inside-this-design) 1 | 1 |
| 6 | Holds | [The protocol](#what-a-rule-does-the-protocol), [Amendments](#amendments-to-accepted-pages) | 1 |
| 7 | Adapted | [Evidence](#evidence) | 2 |
| 8 | Adapted | [Evidence](#evidence), [Open](#open-inside-this-design) 2 | 2 |
| 9 | Adapted | [Evidence](#evidence) | 2 |
| 9a | Deferred | [Deferred](#deferred-with-triggers), [Open](#open-inside-this-design) 6 | - |
| 10 | Adapted | [Refusals](#refusals) | 2 |
| 11 | Adapted | [Optional fields](#optional-fields), [Amendments](#amendments-to-accepted-pages), [Open](#open-inside-this-design) 10 | Its own condition, beside the seven lines |
| 12 | Adapted | [Optional fields](#optional-fields) | 2 |
| 13 | Adapted | [Optional fields](#optional-fields), [Open](#open-inside-this-design) 3 | 4 |
| 14 | Holds | [Optional fields](#optional-fields), [The extension domain](#the-extension-domain) | 2 |
| 15 | Adapted | [Identities](#identities) | 2 |
| 15a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 16 | Holds | [Identities](#identities) | 3 |
| 17 | Adapted | [Identities](#identities), [Open](#open-inside-this-design) 4 | 3 |
| 17a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 18 | Holds | [Identities](#identities), [Amendments](#amendments-to-accepted-pages) | 2 |
| 19 | Holds | [Identities](#identities) | 2 |
| 20 | Holds | [Identities](#identities), [Open](#open-inside-this-design) 5 | 2 |
| 21 | Holds | [Identities](#identities), [The extension domain](#the-extension-domain) | 2 |
| 22 | Adapted | [Identities](#identities) | 3 |
| 23 | Adapted | [Identities](#identities), [Amendments](#amendments-to-accepted-pages), [Open](#open-inside-this-design) 7 | 6 |
| 24 | Adapted | [Identities](#identities) | 6 |
| 24a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 25 | Holds | [Identities](#identities) | 2, 5 |
| 26 | Adapted | [What the kernel carries](#what-the-kernel-carries-without-reading-it), [Open](#open-inside-this-design) 8 | 2 |
| 26a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 27 | Adapted | [What the kernel carries](#what-the-kernel-carries-without-reading-it) | 2, 3 |
| 28 | Adapted | [What the kernel carries](#what-the-kernel-carries-without-reading-it), [Open](#open-inside-this-design) 9 | 1, 3 |
| 29 | Adapted | [What lowering must be able to refer to](#what-lowering-must-be-able-to-refer-to) | 5 |
| 30 | Holds | [What is final and what is provisional](#what-is-final-and-what-is-provisional), [Amendments](#amendments-to-accepted-pages) | 7 |
| 31 | Adapted | [The order of work](#the-order-of-work) | All |
| 31a | Deferred | [Deferred](#deferred-with-triggers), [The order of work](#the-order-of-work) | - |
| 31b | Deferred | [Deferred](#deferred-with-triggers), [The order of work](#the-order-of-work), [Open](#open-inside-this-design) 6 | - |
| 32 | Holds | [The order of work](#the-order-of-work) | 2, 6 |
| 33 | Holds | [The extension domain](#the-extension-domain), [Amendments](#amendments-to-accepted-pages), [Open](#open-inside-this-design) 10 | 7 |
| 34 | Adapted | [Inputs](#inputs-that-are-not-bittypes-alone) | 4 |
| 34a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 35 | Holds | [Inputs](#inputs-that-are-not-bittypes-alone) | 4 |
| 36 | Stale | [Inputs](#inputs-that-are-not-bittypes-alone) | - |
| 37 | Adapted | [Inputs](#inputs-that-are-not-bittypes-alone), [The extension domain](#the-extension-domain) | 7 |
| 37a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 38 | Holds | [Payload bootstrapping](#payload-bootstrapping) | 3 |
| 39 | Holds | [What a test must show](#what-a-test-must-show) | 1 |
| 40 | Adapted | [The gate](#the-gate) | All |

18 rows hold, 21 are adapted, nine are deferred and one is stale.

## The code today

At `ec96b92`. Nothing in this section is the design, and no fact here is a
reason for a choice above. Paths without a directory are in
`contracts/generation/go`.

| Fact | Where |
| --- | --- |
| A registered unit is a profile and three lists: use rules, declaration recipes and helper recipes | `description.go:72-77` |
| A use rule is a pattern with an optional capture, and a fixed output or an expand callback. A declaration recipe has an expand callback | `description.go:21-31`, `:61-66` |
| The expand callback can take the source reader | `description.go:113-120`; `transform/syntax/horos-to-golang/mapping/record/go/record.go:40-41` |
| A rule's callback runs before its children are planned. `Mapped` records a placeholder, and the planner represents the child when it meets the placeholder afterwards | `planning.go:548`, `:560`, `:580-585`; `description.go:133-171` |
| A rule returns native intent, and nothing else of its decision is kept. The program holds declarations, and a step is five strings for inspection | `planning.go:14-22`, `:33-36`, `:564` |
| The declaration recipe takes the record form when the source body is a record. Each target's model requires that the selected body of that form is its record form | `transform/syntax/horos/mapping/declarations/go/declarations.go:52-55`; `transform/syntax/golang/model/go/model.go:105-109` |
| A rule declines by returning an error, which the planner passes on unclassified | `planning.go:378-380`, `:553-555` |
| The record rule's refusal is a text diagnostic. By reading, it matches none of the seven values | `record.go:50` of the Go record mapping; `transform/syntax/horos/source/go/read.go:126-133` |
| A recursive declaration is refused | `planning.go:322-325` |
| An exact mapping without a recipe does not fall back to the structural policy | `planning.go:333-337` |
| A rule's names: a string that is unique in its description, and a counter of the registry. The planner joins them, and only a step keeps the result | `record.go:40` (`record/golang/1`); `transform/syntax/shared/generation/go/registry.go:19-22`, `:126-127`, `:242`; `planning.go:143` |
| A replacement reaches use rules only | `registry.go:251-262` |
| A request names one profile and no product. The registry plans with every unit that states that profile | `registry.go:31-38`, `:185` |
| One native declaration is planned for each source declaration. Its identity is issued at reservation | `planning.go:322`, `:309-317` |
| A use of a declaration is selected by the callee alone | `planning.go:450-453`, `:465` |
| An application keeps its native arguments only. An actual that the native convention omits is admitted and not represented | `planning.go:515`, `:503-507` |
| One generic convention holds native parameters, bindings indexed by source parameter and arguments indexed by native parameter | `description.go:38-45` |
| The optional representation is chosen inside each record rule: Go through a helper key, Rust and C++ by an external type, TypeScript by the flag of the field | `record.go:70-72` (Go), and `record.go:53-54` of the Rust and C++ mappings, `record.go:55` of the TypeScript mapping |
| Go's optional helper is its own contribution. The consumer creates its key, which is an allocation | `transform/syntax/horos-to-golang/mapping/optional/go/optional.go:21-26`; `expression.go:41-46` |
| The `field` form of the Go model carries a name and the optional flag as strings | `transform/syntax/golang/model/go/model.go:43-45` |
| The source reader returns a field's name, its flag, its use and a site | `transform/syntax/horos/source/go/read.go:25-30` |
| A sealed output holds a snapshot, the exposed sites and the external uses. It lists no declaration | `output.go:52-56` |
| The contract package imports six packages of bittype and neither horos nor ontos | `integration/go/architecture_test.go:182` |
| There are 21 mapping packages: five each for TypeScript, Rust and C++, six for Go | The directories `transform/syntax/horos-to-<native>/mapping/<name>/go` |
| TypeScript is checked with `exactOptionalPropertyTypes` | `internal/typescriptsem/compiler.cjs:30` |
| Generated Files formatters exist for horos and TypeScript | The directories under `format/files` |

**Not on `main`.** Pull request
[#281](https://github.com/Bitspark/bittype/pull/281) is open. Read on its
branch at `ca77f65`, it adds owner-issued values, a handle for each
declaration of a sealed output, and an inventory record with a handle, a
site, a role, a provenance and values. The record has no identity of the
planner.

**What does not exist.** No Go file of the repository defines a plan entry, a
payload schema, a facet, a payload codec, evidence about a selected
representation, or a completeness test backend. The service library has no
package `structural-adapters`. Its owner decided that it is created with
integers
([decision of 11 November 2026](https://github.com/Bitspark/bitspark-service/issues/38#issuecomment-6803983510)).

## The order of work

**The steps (row 31, as the grade lists them).**

| Step | What | Where it lands |
| --- | --- | --- |
| 1 | This page | `main` |
| 2 | A prototype | "on a branch that is thrown away" (row 31). It does not merge |
| 3 | A materially different contribution, with the shared packages fixed | The prototype branch |
| 4 | The contract in its final form and the migration of the 21 mapping packages, in one change | `main`, when the [gate](#the-gate) is met |
| 5 | The optional move | `main`, by the answer to [row 11](#optional-fields) |
| 6 | The additive work of [#207](https://github.com/Bitspark/bittype/issues/207) | `main` |

**Why not the advised order.** The advice (section 4.2) asks for TypeScript
end to end through the service library's own factory before the migration.
The table's introduction says why that cannot be done: "The package it wants
the proof to run through, `@bitspark/structural-adapters`, does not exist. Its
owner decided on bitspark-service#38 that it is created with integers, which
the plan puts after the migration." The same facts bound the risk: "no other
repository builds on the contract before integers."

**What the prototype contains (rows 31 and 32).**

- The rule in its staged form, one round ([the protocol](#what-a-rule-does-the-protocol)).
- The frozen decision record, the plan and the serialized form of the plan.
- Two consumers of the evidence path. Row 3: "the declaration recipe's need to
  know that a selected body is a record, and by a fixture family".
- Go handlers and executed Go lowering, in the test tree.
- A TypeScript reader as a stand-in, in the test tree, with the snapshot of
  [row 34](#inputs-that-are-not-bittypes-alone).
- The optional field in the delegated form ([row 11](#optional-fields)).
- Fixtures for vertical cases 3 to 6.
- The four cases of row 32: "repeated uses share a node, equal native types
  with different adapters stay apart, an erased native argument is kept for
  conversion, one plan under two TypeScript placements".

**The six vertical cases in the prototype.** The cases are those of the plan
(`docs/representation-families.md:97-102`). Row 31 puts fixtures for cases 3
to 6 into the prototype, and row 31b the first two "at the kernel's level".
The form of each fixture below is derived here, from verification report 3.

| Case | Its form in the prototype | What that form does not show |
| --- | --- | --- |
| 1. Recursive `ContentValue` | At the kernel's level, with a mock family that the record names as a mock ([open](#open-inside-this-design), point 6) | The recursion family, the memo, the budgets |
| 2. `Page<Outcome>` with optional nullable fields | At the kernel's level, with a mock family. The optional field of `Page<Person>` runs through the field-state protocol | Null against absence, and discriminators |
| 3. The exact external generic `Cell<T>` | An exact generic with an omitted or an additional native parameter and a value facet | Capabilities, and what `Cell` does with its slot |
| 4. Go typed construction with a generic adapter argument | `Page[Person]` itself, lowered and executed | Lowering of families that do not exist |
| 5. A generic map whose key requirement is rejected | A fixture generic that requires evidence about its argument, accepted for one operand and refused for another | The map family and the meaning of its keys |
| 6. An independent family without a shared change | The contribution of step 3 | Authorship by another owner |

The six cases in real form stay what the accepted design makes them: the
condition for a stable version of the envelope
(`docs/representation-families.md:95-96`). The prototype does not discharge
that condition.

**What a stand-in proves (row 31).** Row 31: "the serialized plan is
sufficient for a second reader in another process, not for the service
library's factory."

| A TypeScript reader written by bittype, run from a Go test, shows | It does not show |
| --- | --- |
| No identity in the plan is a pointer of the planning process | That the service library's factory can be written against the plan: its adapter interface with budgets, its call context and its memo |
| Every link, version axis and refusal survives serialization | Independent authorship. One seat wrote both sides |
| The plan, the payloads, the realization and the descriptors suffice for a reader in another language | Anything about their package rules or their CI |

The proof through the service library's own factory is deferred (row 31a).

**A materially different contribution (step 3).** The rules say what counts:
"implement a materially different contribution using its public contracts in
a separately owned package or fixture. Keep the shared mechanism fixed."
(platform `docs/concepts/architectural-correctness.md:487-489`). And what does
not: "A renamed copy of the original example, a private test-only constructor
or a fixture that activates a preinstalled special case does not establish
that." (`:496-497`). The advice (section 4.2): "A second unconditional scalar
is weak P2 evidence."

## The gate

Row 40: "Adopted as the merge gate of the migration, each line in the form
bittype can show: a frozen decision survives native projection, dependency
composition, serialization, a second reader in another process, executed Go
lowering, and realization for TypeScript, with every shared package named and
fixed while a materially different contribution is added".

The migration of step 4 does not merge while one line is red. The gate is
shown first on the prototype branch and then on `main` (row 31). Each line
below is a predicate. Its tests come from the Validation cells of the table,
marked by their row, and from the verification reports, marked V1, V2 or V3.
Every test in the table is red today, for the reason given.

| Line | The predicate | Tests | Red today because |
| --- | --- | --- | --- |
| 1. Native projection | The native type is computed from the frozen record, by a projector that cannot read the source | A projector that tries to read the source does not compile or fails the contract test (row 4). Mutating the value a rule returned changes neither the native type nor the serialized plan (row 28). A fixture whose decision the source does not determine is green on the planner and red on a mutant that derives the entry from the source (row 39, V1) | The one callback takes the reader. No record exists |
| 2. Dependency composition | A rule decides with typed answers about what was selected, and every reference is checked | A fixture family whose decision depends on evidence about its operand plans and generates (row 3). `Page<text>` and the written-out record generate the same native type, and a rule that reads one application fails (row 7). An evidence cycle and a recursive reference are refused with different reasons (row 9, V1). Every refusal of the 21 mapping packages matches exactly one value (row 10). A request at a field position with the operand of another record fails (row 12, V2). One local key under two publishers registers, and one publisher with one key twice fails (row 15). Two uses of `Page<Person>` are two use entries on one node (row 18). Equal native types with different adapters stay apart (rows 19 and 32). An omitted native parameter still yields an adapter argument (row 20). A specialization is refused as unsupported with the application named (row 21). Two contributions that spell a helper alike get two symbols (row 25, V2). A fragment that the target's member vocabulary does not admit is refused with a class (row 14, V2). Guarantees 1 to 6 (row 27) | A rule runs before its children. No evidence, no owner of an identity and no use entry exist |
| 3. Serialization | The plan, read from its serialized form in a fresh process, is the plan that was written | The completeness test runs from the serialized plan in a process that never held the planner's objects (V3). Each version axis is exercised alone (row 16). The regeneration test of the payload types, and a payload outside its schema is refused before a backend runs (row 38, V3) | No plan is serialized. No plan carries the axes |
| 4. A second reader in another process | A reader in another language converts from the plan and the snapshot alone | The TypeScript stand-in round-trips `Page<Person>` and refuses a payload schema version it does not know (V3). Given the plan without the snapshot it fails explicitly, and a snapshot with another pin is a linking error (row 34, V3). Each state of an optional field has its own vector, in Go and in the stand-in (rows 13 and 35) | No reader and no conversion exist |
| 5. Executed Go lowering | Go code that is lowered from the plan compiles and runs | Construction of `Page[Person]`, access to a member, an owned helper and a generic adapter argument compile and run (advice, section 3.4). A lowering that takes a member by a string the naming policy changed fails to compile, and with a member reference it compiles (row 29, V3) | A member is a string. No lowering exists |
| 6. Realization for TypeScript | A reader takes every native name and place from the realization | A plan read with the manifest of another output fails as a linking error, not as an incomplete plan (row 23). One plan under two TypeScript placements resolves each logical symbol to another module path and the same declaration (row 32, V2). No native name reaches emitted glue except through the reader, and a formatter mutant that realizes another spelling is followed (row 24, V2) | The inventory record has no planned identity, and the formatter's result has no reader |
| 7. The shared packages fixed | A materially different contribution is added, and no shared package changes | The contribution registers, and `Page<Person>` regenerates with it. The diff of that test touches no shared package (row 40, V3). A Go module outside the repository authors and registers a contribution with public packages only (row 37, V3). One fixture for each item that is in the [extension domain](#the-extension-domain), and one recorded refusal for each that is out (row 33, V3) | A replacement reaches use rules only. No request can ask at a field position |

**Two checks that are green today and stay green.** The fixtures of today's
generics pass on the new contract with no contribution for generics (row 2,
V1). The native declarations of today's pagination fixtures do not change in
any of the four targets (#207, item 5).

**One condition beside the seven lines.** The test of
[row 11](#optional-fields): "The answer on bitspark-service#38 is recorded
here with the form that was built. Red until then". Whether the migration
waits for it is [open](#open-inside-this-design), point 10.

**The shared packages of line 7.** Row 40 asks that they be "named and
fixed". Derived here, from the list of generation packages in `LAYOUT.md`:

- the contract package, `contracts/generation/go`;
- the registry, `transform/syntax/shared/generation/go`;
- the shared native construction, `transform/syntax/shared/target/go`;
- the source reader and the common mapping parts of the source,
  `transform/syntax/horos/source/go` and `transform/syntax/horos/mapping`;
- the native model and target of each language;
- every existing mapping package;
- the payload codec;
- the completeness test backend.

**What the gate depends on.** Derived here. Line 6 needs the reader of #273
and the inventory of pull request #281 with the planner's identity. Line 4
has its last state, a present `null`, only with #211. The contribution of
line 7 depends on the answer to row 11. Where the Go facets of line 5 are
declared is [open](#open-inside-this-design), point 8.

**What the gate does not establish.** It does not establish the service
library's factory (row 31a), realization for Go and Rust (row 24a), recursion
(row 9a), or vertical cases 1 and 2 in real form (row 31b). Each is recorded
as not demonstrated, with its trigger, under
[Deferred](#deferred-with-triggers).

**After the gate.** Derived here, from verification report 3: the
contribution of line 7, pinned at the commit of the migration, still builds
at the head of each later family. It is red when a later family changes a
type of the contract (row 30).

**When the gate is met** this page's first line changes: the design is final,
with the date and the commit of the migration. Until then each change that
the prototype makes to the design is recorded on this page with its date.
