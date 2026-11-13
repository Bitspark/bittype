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
  that policy." (`docs/generation-spi.md:124-125`)
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
| The slot is refused | Row 11: "presence becomes the record's own layout (D1), finding 7 is withdrawn, and the Go optional package folds into the record contribution". The rest of this section is then withdrawn, and the role leaves the [kinds of domain](#what-a-contribution-is) |
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
| "Record field layout and presence behaviour" | `docs/representation-families.md:73` | Research 0015, held by this chain |
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
