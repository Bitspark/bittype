# The contribution contract

**Design, 13 November 2026, from research 0017. Decided by the research owner, `claude-contract-chain-issues`, who also wrote the research question. Not implemented. It becomes final when the gate at the end of this page is met; until then the prototype may change it, and each change is recorded here.**

This page is the contract of a contribution to bittype's generation kernel:
what a contribution is, what its rule may ask and must return, what
identifies it, and what the kernel carries for it without reading it. It is
step 1 of the [order of work](#the-order-of-work) that
[research 0017](../research-docs/0017-contribution-contract.md) decided for
[#196](https://github.com/Bitspark/bittype/issues/196), item 6. Step 2, the
prototype, is done on a branch that never merges.
[What the prototype showed](#what-the-prototype-showed) says what it built,
what its tests show and what they do not show. The mapping contract on `main`
does not change with this page. [The code today](#the-code-today) describes
the present behaviour, and nothing else on this page is implemented on `main`.

The page rests on two pages:
[binding generation](binding-generation.md) (direction of 9 November 2026,
research 0014) and [representation families](representation-families.md)
(design of 9 November 2026, research 0015). Both were accepted by the owner
of research 0014 and 0015, and "accepted page" below means one of them or
another page of this repository that states its own standing. One sentence of
binding generation is amended in substance: a backend has a fourth input (row
23, O15). Every other insertion in the two pages is a dated clarification that
leaves the meaning of its sentence. All of them are listed under
[Amendments to accepted pages](#amendments-to-accepted-pages).

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
- **O1 to O10:** the ten
  [decisions of 13 November 2026](https://github.com/Bitspark/bittype/issues/196#issuecomment-6841190077)
  on the points that the first version of this page left open. The research
  owner made them. O2, O4 and O6 each correct a row of the table (rows 8, 17
  and 31b), and the table carries the correction with its date.
- **O11:** the
  [decision of the same day](https://github.com/Bitspark/bittype/issues/196#issuecomment-6841518102)
  that corrects the wording of O8: a facet is registered by the consumer's
  composition.
- **O12 to O20:** the nine
  [decisions of 13 November 2026](https://github.com/Bitspark/bittype/issues/196#issuecomment-6841866946)
  after an independent reader, who did not write this page, read it at
  `cab8faf` and reported 25 findings. The research owner made them.
- **O21 and O22:** the two
  [decisions of the same day](https://github.com/Bitspark/bittype/issues/196#issuecomment-6842186587)
  that correct O12 and O19 from the accepted contract of a family, after
  both were written into this page.
- **O23 to O26:** the four
  [decisions of the same day](https://github.com/Bitspark/bittype/issues/196#issuecomment-6842416622)
  after the same reader read the changes at `0102bfe` and reported 13
  findings. O23 corrects words of O13, O24 corrects O22, and O26 corrects
  one point of O18.
- **O27 to O39:** the thirteen
  [decisions of 13 November 2026](https://github.com/Bitspark/bittype/issues/196#issuecomment-6844431866)
  after the prototype was done. O27 extends O25, O32 corrects O8 and O16, O36
  ends the deferral of row 26a, O37 extends O9 and one point of O18, and O38
  decides row 11.
- **O40:** the
  [decision of the same day](https://github.com/Bitspark/bittype/issues/196#issuecomment-6844971187)
  on three points that writing O27 to O39 into this page showed. It widens
  one sentence of O17.
- **The prototype's report**, by its section or by its finding, F1 to F27:
  the [report](../research-docs/0017-contribution-contract.prototype.md) of
  the prototype of step 2, written by the agent that built it. The prototype
  is the branch `prototype/contribution-contract-196` at `6080420`, and it is
  thrown away after the migration. The report is evidence and proposals, and
  no decision. The decisions on it are O27 to O39. **Fixed by the prototype**
  on this page means: the form that the prototype built and its report
  states. Every such form is provisional until the gate. The closing
  paragraph of the decisions takes seven of them by name, and this page says
  which. The others are stated as what the prototype built.
- **Verification report 1, 2 or 3**, in the gate V1, V2 or V3: the reports of
  three verification agents, each on a set of the recommendations of research
  0017, checked against the code at `7e9dc34`. They lie beside the research
  document:
  [1](../research-docs/0017-contribution-contract.verify-1.md),
  [2](../research-docs/0017-contribution-contract.verify-2.md),
  [3](../research-docs/0017-contribution-contract.verify-3.md). A report is a
  proposal and no decision. The decisions are the rows. What this page takes
  from a report is marked "derived here".
- **An accepted page**, quoted, with its file and line.
- **AC-01 to AC-10:** the
  [numbered architecture rules](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#numbered-rules)
  of platform, read at `5e988b0`.
- **BC:** the service library's
  [binding contracts](https://github.com/Bitspark/bitspark-service/blob/c36722db901ef7a259242e168c4c35d6f777fbf4/docs/design/binding-contracts.md)
  at `c36722d`, accepted by that repository's owner seat on 11 November 2026.

A sentence marked **derived here** follows from those sources and is stated by
none of them. It is this page's reading and no decision. A point that two
sources leave unsettled, or settle differently, is not resolved on this page.
The first version listed ten such points, and an independent reading found
more, and the prototype more again. O1 to O39 decide them
([Points decided after the first version](#points-decided-after-the-first-version)).
File and line references are `path:line` at commit `04f9a75`, the head of
`main` when this page was last changed. The prototype's report cites accepted
pages at `2a29674`, where some of their lines have other numbers.

**Words used on this page.**

| Word | Meaning here |
| --- | --- |
| The kernel | The generation contract package `contracts/generation/go` with its planner, and the registry that selects and composes |
| The source integration | The package that reads one source language for the kernel. For horos it is `transform/syntax/horos/source/go` |
| A rule | The part of a contribution that the planner runs for one selection |
| A decision record | The frozen value that holds what a rule decided for one selection |
| A projector | The registered function that computes native intent from a frozen decision record |
| The plan | The selected representation plan of binding generation, tracked in [#195](https://github.com/Bitspark/bittype/issues/195) |
| A generation request | What a consumer asks of generation: the roots, the target, the profile and the products |
| A product | What a generation request asks generation to make possible. O5: "A generation request states whether it asks for types alone or for types and bindings." |
| A reader | A program that consumes the plan for a product: the completeness test backend, or a binding backend |
| A rule's request | One thing that a rule names in step 1 of the [protocol](#what-a-rule-does-the-protocol) for the planner to resolve. There are three kinds: a child use, a field position, and a shared helper. Where this page says "request" alone, it means a rule's request |
| A dependency | What the planner resolved for a rule's request: the representation that was selected for it |
| The source view | What a rule reads of the source for its subject, through the reader of the source integration. A projector gets none |
| A member fragment | What an optional contribution answers on the native side for one field position: the member or members of the native record that represent the field, in the member vocabulary of the target's model ([row 13](#optional-fields)). "Fragment" alone means the same |
| A bundle | A contribution together with the facet packages that are registered for it (O16) |
| The contribution's facet list | The facets that the consumer's composition registered for a contribution. The contribution states none (O32, which corrects O16) |
| The realization manifest | What relates a plan, an output and its realization (row 23). It is not the publisher's manifest, which row 15a defers |
| The test tree | Packages that only a test composition imports. The prototype put them under a top-level directory `completeness/`, which is a proposal for the migration's review (O39, [the roles the prototype proposes](#the-roles-the-prototype-proposes)) |
| The prototype | Step 2 of the [order of work](#the-order-of-work): the branch `prototype/contribution-contract-196` at `6080420`. It never merges |

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
(`docs/representation-families.md:36`). Row 1 adopts the definition. The
four kinds of domain below are the advice's (section 1.1):

| Kind of domain | Example | Standing |
| --- | --- | --- |
| An exact declaration | `Cell<T>` of one module pin | Advice, section 1.1 |
| A builtin form | The closed record, the list | Advice, section 1.1 |
| The declarations that need the default native declaration policy | Every named declaration without an exact mapping | Advice, section 1.1 |
| A role | "Optional record field" | Advice, section 1.1. It stands since [row 11](#optional-fields) is decided as D3 (O38) |

The kernel selects by the domain and never by registration order. Two
accepted pages say it: "Registration order is never a preference. More than
one applicable demanded offer is ambiguous."
(`docs/type-generation-plugins.md:182-183`), and "Multiple demanded offers
fail as ambiguous; registration order is not a preference."
(`docs/kernel-composition.md:87-88`). Row 10 and the advice (section 1.5) add
that a refusal must not make the kernel try another contribution. Two
contributions that claim one domain under one profile need an explicit
composition by the consumer (row 15, from the advice, section 2.1).

**One target and one profile.** A contribution states both. The target is "a
stated field of the contribution, not a part of the key" (O4,
[identities](#identities)). The profile is
stated by the contribution and is no argument of its registration:
"A mapping states its profile in its own description, and its constructor
takes none." (`docs/representation-families.md:361`). What a profile is:
"Representation profiles state what the native type retains and what remains a
source-level obligation." (`docs/type-generation-plugins.md:507-508`). A
profile also carries the coherence requirement that bounds
[evidence](#evidence).

**"Library × target" (row 1).** Binding generation says: "The unit of
implementation is library × target." (`docs/binding-generation.md:269`). Row
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

**The kernel carries a binder until the declaration states its native
convention (O28).** The prototype's finding F1: "Under the protocol the
declaration contribution decides after its body is planned. While the body of
`Page<T>` is planned, the native meaning of `T` is not known." (O28) So: "The
kernel represents a bound use as a scoped reference and replaces it when the
declaration's projector has returned its convention. No rule and no projector
sees it." (O28) This belongs to the first line of the table.

O12 names the second line: the declaration contribution "is the part of the
division of generics that a contribution owns". The first line agrees with the
planner page: "The kernel owns capture identity, substitution and scope."
(`docs/generation-planner.md:435`). It changes two
texts that list generics beside the families. For item 1 of
[#207](https://github.com/Bitspark/bittype/issues/207) the
[comment of 13 November 2026](https://github.com/Bitspark/bittype/issues/207#issuecomment-6840545162)
on that issue states the new reading. The delivery row of representation families that names
"generics" among the delivered families
(`docs/representation-families.md:374`) gets a dated clarification.

## What a rule does: the protocol

**One round (row 3).** Row 3: "a rule names its dependencies and the evidence
it needs, the planner resolves them, the rule decides with typed answers, the
planner freezes the record, the projector runs."

| Step | Who | What happens |
| --- | --- | --- |
| 1. Name | The rule | From the checked subject and the source view, the rule names its requests. A request names one of three things: a child use that the kernel minted, a [field position](#optional-fields) (presence is delegated, O38), or a [shared helper](#identities) (row 25), which is refused until a case needs it (O34). For each the rule says what it needs of the dependency. There are three needs: a reference, evidence ([row 9](#evidence)), or conversion (O31, below). O3: "A request may state what the rule already knows from the source without evidence." |
| 2. Resolve | The planner | It selects a contribution for each request, plans it, and hands the rule the facts of the kinds it asked for that were stated for what was selected. The planner carries a fact unread (O13, O23, [evidence](#evidence)) |
| 3. Decide | The rule | With the source view and the typed answers, the rule returns its decision or a [refusal](#refusals). The decision is a typed value (O9, [below](#what-the-kernel-carries-without-reading-it)). Beside it the rule states the facts that others may ask of it (O23) |
| 4. Freeze | The payload codec, and the planner | The codec encodes, validates and freezes the value "at the moment the rule returns, for every request" (O9, [row 28](#what-the-kernel-carries-without-reading-it)). There "request" is the generation request: it is done for types alone as for bindings. The planner adds the rule's requests that it resolved, child uses included, and the answers it issued. O3: the requests "are part of the frozen record, as links that the kernel carries" |
| 5. Project | The registered projector | It computes native intent from the frozen record. O9: "The projector receives the typed view decoded from the frozen value." The entry of the plan is made from the same record |

Row 3 says that "the planner freezes the record". Derived here, from O9 and
O3: the planner has the codec freeze the returned value, and then adds the
links of the requests it resolved. Both together are the frozen record.

The advice (section 1.2): "One selection does not imply one callback
invocation." The kernel keeps one selected contribution and advances it
through the steps.

**The kernel drives (row 3).** A rule never invokes the planner. Row 3 fixes
the order: the rule names, "the planner resolves them", and the rule decides.
Naming a request is a statement to the planner and runs nothing.

**No analysis on demand (row 7).** What a rule may name in step 1 is bounded.
It may request a child use, a field position or a shared helper, and of what
was selected for a request it may ask for evidence that the selection states
about itself. It may not ask for source analysis. Row 7: "no analysis runs on
demand". The accepted sentence behind it: "Neither the shared host nor a
plugin invokes analysis on demand." (`docs/type-generation-plugins.md:535`).

**A third need: conversion (O31).** Row 9 gives a request two needs, a
reference or evidence. The prototype's finding F8: an exact contribution
names its requests itself and does not know which products a generation
request asks for, so it cannot leave out a child that only conversion needs.
O31: "A rule may name a child that only conversion needs. It is resolved when
bindings are asked for, left out for types alone, and recorded on the link
either way." It is O5 in the protocol, and it keeps the accepted sentence:
"Only demanded native actuals need mappings." (`docs/generation-spi.md:170`).
How such a need counts on a cycle is
[not specified here](#what-is-delegated-and-what-is-not-specified-here).

**The shape the prototype gave the protocol.** Fixed by the prototype
(report, section 7). No decision names it, and no signature is decided: "A
contribution states a reference, target, profile, domain, schema reference,
rule, one projector, declared facts and owned declarations. A rule is two
functions without a target type: Name(step) returns requests, Decide(step)
returns a Decision{Payload, Facts}. A projector receives a Projection with no
source type parameter." The migration writes the final types.

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
  (`docs/generation-spi.md:378-381`).

**What belongs in the decision record (row 5).** Row 5: "A per-selection
choice is a field. An invariant of a semantic revision is carried by reference
to that revision".

- **A per-selection choice** is anything that can differ between two
  selections of the same contribution and that affects representation or
  conversion. The advice (section 1.2) lists member names, the presence
  strategy, the selected construction strategy, a layout that depends on an
  operand, and applicable configuration. For member names O29 says where each
  is recorded: "A record's payload has the names of the fields it lays out
  itself, and the name of a delegated field is in its request. Both row 5 and
  O3 hold that way."
- **An invariant** is what every selection of one semantic revision does the
  same way. It needs no field. The plan names the selected contribution
  through its reference ([identities](#identities)), and the reference
  carries the invariant. By O4 it is the contribution's own representation
  revision that "fixes the invariants of one target".
- **The rule for a projector.** The advice (section 1.2): "A projector may
  implement a recorded policy; it must not choose an unrecorded policy."

**An invariant and the plan's list (O1).** Binding generation lists what the
plan carries, among it "the rules for absence, presence, null and native
`undefined`" (`docs/binding-generation.md:58`). Row 5 keeps an invariant out
of the payload. O1: "Both hold: the plan carries a rule either as a recorded
field, where it varies between selections, or through the contribution
reference whose semantic revision fixes it." Its reason: the completeness test
of that page already names the selected contributions as an input of a
backend (`docs/binding-generation.md:74`). The page gets a dated clarification
([Amendments](#amendments-to-accepted-pages)). Which revision fixes the
invariant of one target is decided by O4 ([identities](#identities)): the
contribution's own representation revision.

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
  (`docs/representation-families.md:352-353`).

The sentence gets a dated clarification
([Amendments](#amendments-to-accepted-pages)).

**There is no general loop (row 3).** A rule gets one round of requests. It
cannot ask again after it has seen an answer. Row 3: "No general loop: no case
found needs a second round." The grade gives the second reason: "an untested
loop would be final without evidence."

What would add a second round is a named case that needs one. Derived here,
from verification report 1: such a case is one in which the set of
requests, or the evidence wanted of one, depends on an answer of the first
round. None of the delivered contributions, 21 packages under six names (O12,
O21), and none of the six vertical cases of the plan (`docs/representation-families.md:109-116`) is such a case by its
description. Adding the round then is a change of the shared contract and
needs the justification of AC-08.

**An illustration.** The advice sketches the protocol and says of its sketch:
"This is a protocol sketch, not a proposed Go API." The same holds for the
lines below. They show the order of one round and decide no signature.

```text
name(checked subject, source view)
    -> requests, each for a reference, for evidence or for conversion
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
| Where does an answer come from? | From what was selected for the dependency. No provider computes it, and the planner reads no payload (O13 and O23, below) | O13, O23 |
| What does the kernel check? | That an answer belongs to the dependency it was asked about: the selected representation and its whole contribution reference, the target, profile and configuration, the generic environment, the kind and version of the requirement, and the scope. O19: an answer of evidence binds "to the whole contribution reference, both of its parts" | Advice, section 1.3; O19; [guarantee 4](#what-the-kernel-carries-without-reading-it) |
| What does the kernel not check? | That the fact is true. "It cannot establish a domain assertion merely because a contribution supplied it." | Advice, section 1.3; AC-05 |

Facts about the source stay with the source. Representation families divide
the judgment: "horos semantics and bittype's source integration establish
canonical facts", and "target families establish native discrimination and
compatibility" (`docs/representation-families.md:222-225`). Evidence in the
sense of this section is the second kind only. A rule reads source facts from
the source view.

**Where an answer comes from (O13).** Row 7 bounds it: only facts of the
selected representation, and no analysis on demand. O13 concludes:

> "So an answer is a fact that the selected contribution of the dependency
> stated with its frozen decision, or that its revision fixes. The planner
> reads it there and issues it to the asking rule, bound to that rule, the
> dependency and the dependency's contribution reference. No provider computes
> anything. Where the selection states no such fact, the answer is `unknown`.
> The answers a rule received are links of its frozen record and so part of
> the node's identity." (O13)

O23 corrects the words "reads it there". The frozen decision is a value that
the kernel carries unread: "It cannot: the kernel carries that value unread,
and the contract package imports neither horos nor ontos." So:

> "So a fact is a typed value of an evidence kind. A rule states it beside
> its decision in what it returns, and a contribution declares at
> registration the facts that its revision fixes. The kernel carries a fact
> unread and hands it to a rule that asked for that kind about that
> dependency." (O23)

| Question | Answer (O13, as O23 corrects it) |
| --- | --- |
| What is a fact? | A typed value of an evidence kind |
| Where is a fact stated? | Beside the decision, in what the dependency's rule returns. A fact that the contribution's revision fixes is declared by the contribution at registration |
| Does the planner read the frozen decision for it? | No. It carries the fact unread and hands it to a rule that asked for that kind about that dependency |
| Who computes it? | Nobody at that moment. The dependency stated the fact when it decided, or its contribution declared it |
| What is an answer bound to? | The asking rule, the dependency and the dependency's contribution reference |
| What if the selection states no such fact? | The answer is `unknown` |
| Is an answer part of the frozen record? | Yes, as a link, and so part of the [identity of the node](#identities) |

O23 left two things to the prototype: "the form of an evidence kind, and how
a rule could present an answer that was not issued to it, which is the entry
point of guarantee 4's test." Both are fixed by the prototype, and the
decisions take them as it fixed them (report, section 7).

| What | The form the prototype fixed. Provisional until the gate |
| --- | --- |
| An evidence kind | "EvidenceKind{Identity, Version}, a comparable value owned by a vocabulary of the target or product." A kind is a qualified identity with a version |
| A fact | "A fact is Fact{Kind, Outcome, Content}: the outcome is proved, refuted or exhausted, and the content is the owner's string, carried unread." A stated fact has one of three outcomes. The other two outcomes of an answer are not stated by anyone: `unknown` is what the planner answers where nothing is stated, and deferred is its answer for a binder |
| How a rule presents an answer | "An answer is a value with unexported fields that a rule reads only through the step it was issued to: DecideStep.Present(dependency, answer). It is bound by value to the asking contribution, the subject's place, the dependency, the dependency's whole reference and the generic environment. A rejected answer fails the planning of the subject even where the rule drops the error." This is the entry point of the test of [guarantee 4](#what-the-kernel-carries-without-reading-it) |

**Evidence through a declaration (O35).** The prototype's finding F23: "What
is selected for a use of `Person` is the declaration contribution's
application, not the record in its body." (O35) So a rule that asks about an
operand of a named type asks the declaration contribution. O35 decides the
case of an alias and leaves the nominal case.

| The declaration | What O35 decides | Its reason |
| --- | --- | --- |
| An alias | "an alias is interchangeable with its target, so the facts of an alias's body are the facts of the alias, and the kernel forwards them. The declaration contribution states with its decision that it is transparent." | The accepted pages: "If the source language makes an alias interchangeable with its target, choosing incompatible native nominal types may violate the selected representation profile." (`docs/type-generation-plugins.md:209-211`), and "Administrative aliases and renamed binders are permitted" (`docs/type-generation-plugins.md:485-486`) |
| A nominal declaration | "nothing is decided: it answers for itself, bittype's own states no fact, so the answer is `unknown`, which is the safe answer under coherence." | "Which facts a nominal declaration restates is target-specific and is decided with the first family that needs evidence through one, #209 (maps and sets)." |

The prototype did not build the forwarding. Its fixture policy asked about
its body and stated the answer again (report, F23).

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

**The outcomes (row 8, O2).** An answer has one of five outcomes. Row 8, as
O2 corrects it: "The outcomes of evidence are proved, refuted, unknown and
exhausted, plus deferred to a binder requirement."

| Outcome | Meaning | What a rule may conclude |
| --- | --- | --- |
| Proved | The selected representation establishes the fact | The fact |
| Refuted | It establishes the opposite | The opposite |
| Unknown | Nothing establishes either | Nothing. The advice (section 1.3): "Absence of proof that an operand excludes null is not proof that it does." |
| Exhausted | The check reached its limit | Nothing |
| Deferred to a binder requirement | The operand is bound to a binder. The answer is what the binder's declared requirement establishes, and each application must discharge that requirement | The fact, if the requirement states it. Otherwise nothing |

**The accepted name `overlap` stays (O2).** Representation families names
the outcomes of one analysis, the disjointness of a union: "`proved`,
`overlap`, `unknown` and `exhausted`, and neither of the last two establishes
anything" (`docs/representation-families.md:230-232`). The first form of row
8 called its own four names the accepted ones. O2: "Row 8 was imprecise." And:
"Evidence in general has proved, refuted, unknown and exhausted. For
disjointness the refutation keeps the accepted name `overlap`." So `overlap`
is the refuted outcome of that one judgment, and no accepted name changes.

The fifth outcome is the accepted course of a generic family: "A generic
family proves its requirements for every permitted actual, declares
requirements that actuals must satisfy, or defers qualification to checked
applications." (`docs/representation-families.md:283-285`).

**A binder's requirement: what the accepted pages fix (O13).** O13 asks this
page to quote them and to list what they leave open.

| Question | What an accepted page says, and this page's reading |
| --- | --- |
| Why is a requirement needed? | "A native type parameter does not carry the adapter, equality, discriminator or layout evidence needed to bind values of that type" (`docs/representation-families.md:271-272`) |
| Who declares a binder's requirement? | The accepted sentence above names three courses of a generic family, joined by "or". Declaring is the second: the family "declares requirements that actuals must satisfy" (`docs/representation-families.md:283-284`). Derived here: in the division of [row 2](#what-a-contribution-is) the family is the selected representation contribution |
| Where is it carried? | With each generic use: "Each generic use carries its scoped binder identities, native type arguments, adapter arguments and requirements" (`docs/representation-families.md:106-108`) |
| Who discharges it at an application? | No accepted sentence says it. The third course of that sentence, to defer "qualification to checked applications", is another course and not the discharge of a declared requirement. "Qualifying `U<int>` does not qualify `U<text>`." (`docs/representation-families.md:285`) speaks of qualification evidence and not of a planning check. Derived here, from the advice: the planner makes the "Instantiation and compatibility checks" (section 1.1), and "A concrete application must discharge it." (section 1.5) |
| Under which assumption is an open declaration planned? | "Open eligibility must be established under binder assumptions over the profile's explicitly represented substitutions" (`docs/type-generation-plugins.md:503-505`) |

**How a binder's requirement is declared, recorded, found and discharged.**
Fixed by the prototype within those bounds, and the decisions take it as the
prototype fixed it (report, section 7).

| Step | The form the prototype fixed. Provisional until the gate |
| --- | --- |
| Declared | "Declared by the selected representation contribution in step 1, as part of its demand (Require)." |
| Recorded | In the frozen record of that rule, as an answer that is deferred and required, and "collected on the node of the open declaration as (binder, kind)." |
| Found | "Found by the planner on the open declaration, which is planned before any application." |
| Discharged | "Discharged by the planner at each application from the facts of what was selected for the actual". A proved fact is recorded with the actual's reference. A refuted one is refused as requirement unsatisfied, and an unknown or exhausted one as requirement unresolved. "An actual that is itself a binder moves the requirement to that binder." |

So an application that cannot discharge a requirement is refused with one of
two classes, by the outcome ([refusals](#refusals)). Only the reason is its
owner's.

**A choice made for want of proof (row 8).** A rule that has no proof may
take a representation that needs none, if its contribution offers one, or it
may refuse. Two limits hold.

- **It never replaces what was selected for the operand.** The choice is the
  rule's own, and it is recorded in the rule's decision record. The accepted
  pages: "A validated canonical-value wrapper is a deliberate selected mapping
  for hard forms. It is never a fallback."
  (`docs/representation-families.md:266-267`), and "No retry selects a weaker
  representation." (`docs/kernel-composition.md:105`).
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

**What "fails" means there (O29).** O29 corrects the last words of that test.
They "cannot be a refusal of the kernel: the kernel gives no rule an
application to read, and it does not know what a rule does with an answer."
Inside `Page<T>` the answer about `T` is deferred whatever `Page` is applied
to, and a rule can still decide differently on a deferred answer than on a
proved one (report, F5). So: "The gate's test reads: a rule that decides
differently inside a generic declaration than for its written-out application
fails its family's coherence test." (O29) Coherence is the contribution's to
keep, as the accepted page says: "This is a trusted-plugin obligation, not a
property Go interfaces can prove." (`docs/type-generation-plugins.md:487-488`).

**References against evidence (row 9).** Row 9 tells two needs of a request
apart. O31 adds a third, conversion ([the protocol](#what-a-rule-does-the-protocol)).

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
- **What makes a cycle an evidence cycle (O30).** The sentence above does not
  tell the two cycles apart: the declaration contribution itself demands
  evidence about its body, so every recursive declaration has a demand for
  evidence on its cycle (report, F6). O30: "A cycle is an evidence cycle when
  every open request on it demands evidence (F6). With one reference on it,
  it is a recursive reference. Without this the two refusals of row 9 cannot
  be told apart, because the declaration contribution itself demands evidence
  about its body."
- **A need for conversion on a cycle counts as a reference (O40).** O31's
  third need "asks the planner to represent a child for a binding product and
  demands no answer before the rule decides. So it counts as a reference."
  (O40)
- **Recursive references are deferred (row 9a).** Accepted for the families
  that need them: "Recursive adapters are built through explicit plan
  references (placeholders), separately from native value indirection."
  (`docs/representation-families.md:254-255`). They are not built with this
  contract. A recursive reference is refused. The trigger is under
  [Deferred](#deferred-with-triggers). The prototype shows nothing of
  recursion beyond the two refusals: vertical case 1 is not in it (O6,
  [the order of work](#the-order-of-work)).
- **Two refusals that can be told apart.** The distinction of row 9 can be
  observed only if the refusal of an evidence cycle and the refusal of a
  recursive reference carry different reasons. Fixed by the prototype
  (report, section 7): "Recursive reference and evidence cycle are both
  unsupported, with the kernel as owner and two reasons." Its test runs on
  the independent fixture source, because the horos source refuses a
  recursive declaration before the kernel meets it (report, F7).

## Refusals

**What exists (the generation contract at `04f9a75`).** Listed here because
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
| Unsupported | The selected contribution recognizes the case and cannot represent or check it yet. A later slice can admit the same input. O17: it "is for a case that a selected contribution recognizes and does not cover, as row 21 uses it". The prototype also uses it where the owner is the kernel or a target's member vocabulary (below) | `ErrUnsupported`, narrowed. Today that value also covers a missing mapping, a missing helper and recursion |
| Refused by design | The case lies outside the domain that the representation promises. Another contribution may cover it | New. It is not `ErrRejected`, which is the source's refutation of a use |
| Requirement unsatisfied | A selected dependency fails a declared requirement: the evidence is refuted. Also an application whose actual refutes a binder's requirement | New |
| Requirement unresolved | The evidence a decision requires is unknown or exhausted. Also an application whose actual states nothing about a binder's requirement | The same meaning as `semantic.ErrEvidenceUnavailable` has for source evidence |
| Missing facet | The selected contribution has no facet for a requested product | New. It stays a class of its own |
| Ambiguous or missing contribution | More than one contribution, or none, answers a domain. A family that is not built is the second case (O17). They are two values (O26, below) | `ErrAmbiguous` for the first. The second gets a value of its own. Today it is under `ErrUnsupported` |
| Limit | A resource bound of the request is reached | `ErrLimit` |

The first four meanings follow the advice (section 1.5). Which outcome of
[evidence](#evidence) leads to the third and which to the fourth is derived
here. "Missing facet" stays
recognizable because representation families requires it: "Unsupported use,
missing facet, incompatible plan or runtime, invalid native value, structural
mismatch, domain violation, declared coverage restriction, cycle and
exhaustion stay distinct." (`docs/representation-families.md:334-336`). A
limit is no judgment about the input: "Resource exhaustion lies outside these
sets: a valid value stays valid when a call lacks the resources to convert
it." (`docs/representation-families.md:148-150`).

**What stays as it is. Derived here** from "refined, not replaced":
`ErrInvalid`, `ErrScope`, `ErrIncomplete`, `ErrRejected`, the four verdicts of
a source judgment, and cancellation. A source judgment stays the source's. A
refusal of this section is a contribution's or the planner's.

**What a refusal carries (row 10).** "an owner-qualified reason and the
subject". The advice (section 1.5) lists what to preserve: the source
subject, the selected contribution, the target and profile, the path of
dependencies and the requirement concerned.

- **The class is the kernel's, and the set is closed in the prototype's
  form.** Failure meaning is one of the things composition must preserve
  (AC-05). A closed set of classes in the kernel is therefore no catalogue of
  contributions. O17: "The set of refusal values is closed once the prototype
  has fixed the members that the page lists as undecided, and the page says
  so." The prototype fixed them, at the end of this section: "The set of
  eight values can close as the page lists it." (report, section 7). The
  eight are the seven of row 10 with missing and ambiguous apart.
- **Missing and ambiguous are two values (O26).** Row 10 lists them as one
  item. O26: "They differ in what a consumer does about them", and the
  accepted page tells them apart: "More than one applicable demanded offer is
  ambiguous." (`docs/type-generation-plugins.md:183`).
- **The reason is its owner's, and the set is open.** A reason explains. It
  never adds a class or moves a refusal to another one. Derived here.

**Every refusal of a rule has a value (row 10).** "A rule's refusal that
matches no value today gets one." The test of the row: "Every refusal of the
21 mapping packages matches exactly one value of the contract." The prototype
did not run this test in full: ten of the 21 packages are moved, and five
refusals of moved rules match exactly one value each
([what the prototype showed](#what-the-prototype-showed)).

**No fallback (row 10).** "A selected exact mapping's refusal never falls
back." This keeps three accepted sentences:

- "There is never a fallback to structural conversion, in the same way that
  an exact use mapping never silently becomes a structural publication today."
  (`docs/binding-generation.md:103`)
- "An exact binding takes precedence, and its failure never falls back to
  that policy." (`docs/generation-spi.md:134-135`)
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
  (`docs/naming.md:317-319`,
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

**What row 10 left undecided.** Fixed by the prototype (report, section 7).
The decisions take the last line as the prototype fixed it. No decision names
the first three.

| What | The form the prototype fixed. Provisional until the gate |
| --- | --- |
| The form of a refusal | "One type with a class and one value for each class." "A rule may use four classes: unsupported, refused by design, requirement unsatisfied, requirement unresolved." |
| The two refusals of [row 9](#evidence) | "Recursive reference and evidence cycle are both unsupported, with the kernel as owner and two reasons." |
| An unexpected error of a rule | "An unexpected error of a rule is RuleFailure, which matches no refusal value." |
| A member fragment that a target's member vocabulary does not admit | Unsupported, with the member vocabulary of the target as owner. The report's reason: "a later slice of that vocabulary can admit the same fragment, which is the page's criterion for unsupported." It is not refused by design, "because the contribution's domain covers the position." |

**The owner of an unsupported refusal (O40).** O17 writes of "a case
that a selected contribution recognizes and does not cover". In these forms
of the prototype the owner is not a selected contribution: the kernel for a
recursive reference and for an evidence cycle, and a target's member
vocabulary for a fragment. O40 widens O17's sentence: "The owner of the
refusal is whoever would build that slice: a selected contribution, the
kernel, or a vocabulary." The criterion stays: a later slice can admit the
same input. "An evidence cycle is unsupported and not refused by design,
because the advice names what would lift it" (O40): "a declared evaluation
discipline" (advice, section 1.4). What O17 decided stays: "a family that is
not built has no contribution to select, and that is a missing
contribution." (O40)

## Optional fields

> **Row 11 is decided: presence is delegated (D3, O38).** Until 13 November
> 2026 this whole section was conditional on an answer of the service
> library's owner. O10 gave the decision "to the answer on
> bitspark-service#38 or, if none had come when the prototype was complete,
> to me on its evidence. None has come." (O38) "Me" is the research owner.
> The section states the decided form, and below it what an objection would
> reopen: "An objection with a reason, before the migration merges, reopens
> this." (O38)

**The decision and its evidence (O38).** "The evidence: delegation was built
in Go and in the stand-in; neither side needed the other contribution's
payload, a helper key from the consumer, or a layout across fields; the four
states behave as O14 decides." (O38)

| | Go | The TypeScript stand-in |
| --- | --- | --- |
| The member fragment | One member of its own type, `BittypeOptional[T]` | One optional property |
| The field-state protocol | Two operations. `Construct(present, abstract value)` returns the member, and `State(member)` returns whether it is present and the abstract value | Two operations. `write(container, member key, present, canonical value)`, and `read(container, member key)`, which returns whether the field is present and its value |
| What the record's factory hands over | A member value, and it takes one back | The container object and the member's key, "because an absent property is a fact about the object and has no value" |
| What the optional contribution had to know | The adapter selected for its operand and the declaration it owns. Not the record, the member's name or the canonical name | The adapter selected for its operand |

The table is the report's section 8.

- **The shape of the slot.** O38: "One asymmetry: for Go a value in and a
  value out is enough, and TypeScript needs the containing object and the
  member's key, because an absent property is a fact about the object. A slot
  kind for every target has the TypeScript shape."
- **What it does not establish.** O38: "What this does not establish is the
  service library's own adapter interface." One seat wrote both sides, and
  the stand-in is bittype's test code.
- **Put to that owner.** The shape is the concrete form of the first question
  on bitspark-service#38, in a
  [comment of 13 November 2026](https://github.com/Bitspark/bitspark-service/issues/38#issuecomment-6844432275).
  The slot kind is not agreed.

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
| The accepted text that supports it | The payload column of the plan lists "Record field layout and presence behaviour" as one item (`docs/representation-families.md:87`) | The family list names "optional and nullable fields" (`docs/representation-families.md:10-11`), and the delivery table names "optional fields" among the delivered families (`:374`) |
| What it needs from the service library | Nothing | A field-level conversion slot |

**bittype's position (row 11).** Row 11: "the standard record representation
delegates per-field presence through a published protocol (D3)". Its two
reasons, in the row's words: "the accepted family list names optional fields
as a family, and finding 7 of `docs/package-boundaries.md`". The row also says
what weakens them: "The accepted texts also support the other reading, and the
choice needs a field-level conversion slot in the service library's adapter
interface."

Finding 7 is bittype's own derivation and no owner sentence. With O38 it
stands. It reads "Record
decides optional fields itself in Go, Rust and C++"
(`docs/package-boundaries.md:159`) and names three targets. Derived here:
under D3 the TypeScript record contribution asks as well, because the
position is about the standard record representation of every target. The
prototype built it so, with a new TypeScript package `optional`.

**What row 11 foresaw, and what happened.**

| The case | What row 11 and O10 said | What happened |
| --- | --- | --- |
| The slot is admitted | D3 stands. Each target's record contribution asks, each target has an optional contribution, and the move is step 5 of the [order](#the-order-of-work) | No answer has come |
| The slot is refused | Row 11: "presence becomes the record's own layout (D1), finding 7 is withdrawn, and the Go optional package folds into the record contribution" | No answer has come |
| No answer yet | The prototype builds D3. Row 11: "put to that library's owner on bitspark-service#38 and built in the prototype" | Built, in Go and in the stand-in |
| No answer when the prototype is complete | The research owner decides row 11 on the evidence of the prototype and records the reasons (O10) | Decided as D3 by O38, with the evidence above |

**What an objection would reopen.** If the owner of the service library
objects with a reason before the migration merges, row 11 is open again. If
the slot is then refused, the form is D1, and these statements of this page
change with it:

- Presence becomes the record's own layout, finding 7 is withdrawn, and the
  Go optional package folds into the record contribution (row 11).
- The record's rule changes as the report says (section 12): "the record's
  Name step keeps Reference requests for optional fields and the payload
  gains the presence decision; the field request, the role and the member
  vocabulary are not needed."
- The subject at a field position and the member-fragment answer are not
  built, and the role leaves the [kinds of domain](#what-a-contribution-is)
  (O10).
- The delivered contributions have five names and not six (O12).
- The line of row 33 on a contextual representation at a field position is
  out of the [extension domain](#the-extension-domain) (O10).
- The gate's tests of the delegated form leave it (O18). The vectors of the
  four states stay and run against the record contribution (O26).

**What is in the contract under either answer (O10).** "The round of requests
is part of the contract under either answer, because evidence needs it. The
subject at a field position and the member-fragment answer are built only if
presence is delegated." So a refusal would remove nothing from the
[protocol](#what-a-rule-does-the-protocol). Presence is delegated since O38,
so both are built.

**Why the fallback is no design decided by an existing interface.** Row 11
lets the answer of another repository decide between D1 and D3. AC-09 keeps
an existing interface from deciding a design. The decisions of 13 November
2026 say why this is not that case: "the accepted texts support both
readings, and the slot belongs to an interface that another repository owns,
so the choice is made with that owner. If no answer comes, O10 decides it on
the prototype's evidence, not on what the interface offers today." That is
what O38 did.

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
| Conversion | "a field-state protocol for every target, not only TypeScript. Its slot kind is agreed with the service library (row 11)". It is not agreed yet. Its shape is put to that owner (O38, above) |

- **The native half.** The record contribution and the optional contribution
  both use the member vocabulary of the target's native model. Neither imports
  the other. The record chooses the member's name, as the naming design has
  it: "A member's name and accessibility" belongs to "The declaration, so the
  mapping of that pair" (`docs/naming.md:154`).
- **How the name travels (O3).** "A request may state what the rule already
  knows from the source without evidence." The record contribution knows the
  name from the source view, so it states the name in its request. The
  resolved requests are part of the frozen record, and O3 concludes: "So a
  name stated in a request is recorded once, there, and the decision does not
  restate it." Row 5 makes a member's name a field of the payload. O29 says
  how both hold: "A record's payload has the names of the fields it lays out
  itself, and the name of a delegated field is in its request." The prototype
  built it so: the record's projector pairs the names of its typed view with
  the fragments of its dependency views by the order of the links (report,
  F3).
- **The member fragment.** Fixed by the prototype (report, section 7 and
  stage 6): a fragment may hold more than one member in the Go and the
  TypeScript model of the prototype, a fragment that is no member is refused
  ([refusals](#refusals)), and a target without a member vocabulary refuses a
  field position. No decision names the first.
- **The conversion half.** Only a record can omit a field, in every target.
  So no target converts an optional field from its retrieved value alone. For
  TypeScript representation families already says it: "the adapter separates an
  absent own property from a present `undefined` at run time"
  (`docs/representation-families.md:198-199`). The advice (section 1.6) lists
  what the protocol must be able to do: test presence, obtain a present value,
  delegate that value to the adapter selected for the operand, and construct
  or omit the native field according to the recorded decision.
- **Who invokes it.** The record's factory invokes the protocol for a field.
  It does not read the payload of the optional contribution.
- **The states.** The protocol tells four states of a field apart. They, and
  what a present `undefined` means, are under
  [Inputs that are not bittype's alone](#inputs-that-are-not-bittypes-alone).
- **Its law** is [not specified here](#what-is-delegated-and-what-is-not-specified-here). The prototype
  did not fix it (report, section 7).

**The limits (row 14).**

- **A layout shared across fields is out.** Row 14: "A presence layout shared
  across fields needs the record policy replaced." A request concerns one
  field. A bitmap for all optional fields of a record cannot be answered field
  by field.
- **A contribution without a value adapter is in.** Row 14: "The contract can
  express a contribution with no standalone value adapter". An optional
  contribution has a field-state protocol and nothing that converts a value on
  its own. The advice (section 1.6): "Contextual field conversion is an
  admitted product shape that the contract needs to express now." This
  sentence of row 14 holds under either decision of row 11. O10 made
  conditional only the subject at a field position, the member-fragment
  answer and one line of row 33.

**What the decision changes outside this page.** Three accepted sentences
were written without a field-level slot. Until O38 they waited for the
decision of row 11. They are not edited with this change. O40: bittype's two
"are annotated with the migration, in the change that builds delegation on
`main`, and not before: until then an objection can reopen row 11. The
service library's sentence is its owner's."

| Sentence | Where | Owner |
| --- | --- | --- |
| "A family reaches other families only through its adapter slots." | The [decision of 11 November 2026](https://github.com/Bitspark/bitspark-service/issues/38#issuecomment-6803983510) on bitspark-service#38 | The service library, `codex-platform-status` |
| "Run-time factory reference and typed argument-adapter slots, converting the value inside a message" | `docs/binding-generation.md:98` | Research 0014 |
| "Record field layout and presence behaviour" | `docs/representation-families.md:87` | Research 0015. This chain holds its delivery |

Finding 7 (`docs/package-boundaries.md:159`) carries a dated note since this
page merged. With this change the note says that the finding stands
([Amendments](#amendments-to-accepted-pages)).

## Identities

The [list at the end of this section](#the-identities-in-one-list) names
every identity with its issuer.

**The contribution identity and reference (row 15).** Row 15 adopts the form
of the advice (section 2.1). The lines show the parts and decide no
encoding:

```text
contribution identity  = (publisher, local key)
the family's meaning   = (meaning identity, meaning revision),
                         the same in every target
contribution reference = (contribution identity, the family's meaning,
                          representation revision)
```

The second and third line are as O4 decides and O19 names them,
[below](#identities). Row 15 wrote the reference with one semantic revision.

- **The publisher.** Row 15: "The publisher is the library that publishes the
  mapping, the accepted pages' "library's owner"." The accepted sentence:
  "Contributions belong to their library's owner, outside canonical horos
  bytes." (`docs/binding-generation.md:105`).
- **What the publisher is not** (advice, section 2.1): the owner of the
  source declaration, the module pin of the source that is mapped, or the
  profile. A third party may map another owner's declaration, and several
  contributions state one profile.
- **The local key** tells apart the contributions of one publisher. O4: "The
  publisher's local key is unique for that publisher across targets, and the
  target is a stated field of the contribution, not a part of the key." The
  test of row 15: "Two contributions with one local key under two publishers
  register side by side, and two with one publisher and one key fail."
- **A name does not authenticate.** The advice (section 2.1): "A namespace
  string alone does not authenticate its owner; consumer composition must bind
  it to the selected publisher and manifest." That manifest is the
  publisher's. O16: "Row 15a's manifest is the publisher's, bound by the
  composition lock, and stays deferred." Derived here: until then the
  registration in one process is the binding.
- **It does not find a factory.** Row 15 corrects a premise of the submitted
  document. A backend finds a factory by what the plan entry names: "By
  default a plan entry names the family's run-time factory and gives it the
  family's payload as arguments." (`docs/representation-families.md:63-64`).
  The contribution identity says what was selected. It is what duplicates and
  version checks compare.

**The five axes stay apart (row 16).** Representation families: "**Five version
axes:** the envelope protocol; the family payload schema; the
representation's semantic identity; generated-code/runtime compatibility;
source semantics." (`docs/representation-families.md:102-104`). Row 16: "The
contribution reference carries the semantic revision. The other four axes are
stated where each belongs and are not folded into the identity".

| Axis | What it identifies (advice, section 2.1) | Where it is stated. Derived here, except the third |
| --- | --- | --- |
| Envelope protocol | How the shared graph is represented | The plan (#195) |
| Payload schema | How a family's recorded decisions are typed and encoded | The family's payload schema ([bootstrapping](#payload-bootstrapping)) |
| Representation semantics | What the selected representation means | The contribution reference (row 16), in two parts: the revision of the family's meaning and the contribution's representation revision (O4, O19) |
| Generated-code and runtime compatibility | Which generated code and runtime facilities can cooperate | The facet that names a factory, among its requirements (`docs/binding-generation.md:100`) |
| Source semantics | Which source interpretation the checked references use | The source-semantics input ([row 34](#inputs-that-are-not-bittypes-alone)) |

A digest of an implementation is provenance. It replaces none of the five.
The test of row 16: "Each axis is exercised alone: an unknown required value
of one axis is refused while the other four are unchanged." For the fourth
axis it is a reader's test. O39: "The test of row 16 for the runtime
compatibility is a reader's test, because the kernel carries that facet
unread (F13)." For the third
axis the test runs for each of its two parts. O19: the test of row 16 binds
"to the whole contribution reference, both of its parts". The axes stay apart
from the representation revision too. O24: "A new schema version or a new
runtime compatibility does not by itself make a new representation revision."

**Two identities: the family's meaning and a representation revision (row
17, O4, O19).** The advice asked for a reference to a family's contract that
spans targets. The first form of row 17 folded it into one string of the
contribution reference. That left no place for an invariant that only one
target has, such as the advice's own example of an optional representation
that "always ignores its unused value when absent" (section 1.2). O4: "Row 17
was too short." The decision, which the row now carries:

> "There are two identities. A family contract has an identity and a revision
> that are the same in every target: what values the family covers and what
> they mean. A contribution has its own representation revision, which fixes
> the invariants of one target. A contribution reference names both." (O4)

**The first identity is named "the family's meaning" (O19).** O4's words
"family contract" collide with an accepted use of "contract": "A family
contributes, per target, one contract with seven parts"
(`docs/representation-families.md:32`). One of the seven parts is "Value
meaning": "Promised canonical domain, domain predicates, native equivalence,
any permitted normalization" (`docs/representation-families.md:38`). O19
took that whole part as the one that is the same in every target, and named
O4's identity after it. So this page writes "the family's meaning", with a
meaning identity and a meaning revision, where O4 wrote "family contract".

**Which half of "Value meaning" (O22, corrected by O24 below).** O22 corrects
O19: "The first two
speak of canonical values and are the same in every target: they are the
family's meaning. The native equivalence and a permitted normalization speak
of native values, which differ by target: the contribution's representation
revision fixes them, with the other six parts." (O22) The first two are the
first two things that the part lists.

**The family's meaning is the application's promised domain (O24, which
corrects O22).** O22 did not use the sets that representation families
defines under "Domains". `V_T` is "The application's promised domain,
including identified domain predicates"
(`docs/representation-families.md:142`). `A_T,p` is "The canonical values
this representation promises to cover"
(`docs/representation-families.md:143`), and "`A_T,p ⊊ V_T` is **declared
partial coverage**" (`docs/representation-families.md:147-148`). O24: "The
family's meaning is `V_T` with its predicates. What a representation covers,
`A_T,p`, and a declared partial coverage belong to the contribution's
representation revision, with the native equivalence and a permitted
normalization."

O24 also corrects the last words of O22, that the representation revision
fixes the other six parts: "That is too wide: row 16 keeps the other version
axes apart, and two of them, the payload schema and the runtime
compatibility, are among those parts. The representation revision fixes the
invariants of one target's representation, as row 5 says."

**The family's meaning has one home (O27, which extends O25).** O25 put
identities into a vocabulary package of the family for the target. The
prototype then spelled the meaning identity once in each target's vocabulary
package and compared the spellings in a test (report, F2). O27: "That is two
spellings of one name." bittype's rule puts identities into vocabulary
packages (`AGENTS.md:55-56`), "and a name that two parties agree on has one
owner. So: a vocabulary package of the family, without a target, holds the
meaning identity and its revision, and each target's vocabulary imports it."
The prototype did not build that package. O27 leaves one thing as O19 left it:
"Who owns it where two publishers contribute stays unspecified".

| | The family's meaning | The representation revision |
| --- | --- | --- |
| Where it is spelled | Once, in a vocabulary package of the family, without a target (O27) | In the vocabulary package of the contribution, which is for one target (O25) |
| What it is | `V_T`, the application's promised domain, with its predicates (O24) | What the representation covers, `A_T,p`, and a declared partial coverage; the native equivalence and a permitted normalization (O24). Not the payload schema and not the runtime compatibility, which keep their own axes |
| Whose | Not specified where two publishers contribute to one family. O19: "Who owns a family's meaning when two publishers contribute to it is not specified." | The contribution's publisher |
| The same in every target | Yes | No. It belongs to one contribution, so to one target |
| What it fixes | Which canonical values the application promises, and what they mean | The invariants of one target's representation ([row 5](#what-a-rule-does-the-protocol)) |
| What changes when one target's invariant changes | Nothing | This revision, for that contribution alone |

The use of the family's meaning is recognition across targets. The binding
contracts say that a namespace is "a field the plan selects as the structural
record family" (BC, section 2). A root facet would recognize a declared
meaning, not the name of one target's contribution. That is the advice's
proposal and the service library's to decide. It was not put to that owner,
and the root facet is deferred (row 17a).

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
| "The memo is keyed by the selected representation instance plus input identity, including the descriptor and generic environment." (`docs/representation-families.md:293-294`) | The applied representation, which several uses share. It is not the per-use record |

Read with the first sense, the second sentence would give every use its own
memo key, and one object reached through two fields would convert twice.

**The identity of a node (row 19).** Row 19: "Identity is over the selected
semantic representation without diagnostic provenance".

- **What is in it.** Everything that makes two selections represent
  differently. The advice (section 2.2) lists what a shorter key would miss:
  the recorded decisions, the target configuration, the representations
  selected for the dependencies, the identity of a nominal native
  declaration, the generic conventions and the evidence. O13 confirms the
  last: the answers a rule received "are links of its frozen record and so
  part of the node's identity".
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
(`docs/representation-families.md:106-108`).

- **Two conventions, kept apart.** The native arguments and the adapter
  arguments of an application are two lists, each tied to binders. They are
  never matched by position. A native declaration may drop a source
  parameter that conversion still needs, and it may add a native parameter
  that no source parameter stands for.
- **An omitted argument is kept for conversion.** Row 20: "An argument that
  the native convention omits is still represented for conversion". The test:
  "A native convention that omits a parameter still yields an adapter argument
  for it."
- **One sentence of the overview is read as the native side. Derived
  here.** Row 20 keeps a native convention and puts an adapter convention
  beside it. The overview of generation says: "A single generic convention
  relates source binders, native parameters and complete mapped actuals."
  (`docs/generation.md:45-47`). That page describes type generation, and this
  page reads its sentence as a statement about the native convention. The
  sentence is not annotated
  ([Amendments](#amendments-to-accepted-pages)).
- **A request names its products (O5).** "A generation request states whether
  it asks for types alone or for types and bindings." The decision bounds row
  20: "Row 20 applies when bindings are asked for: an argument that the native
  convention omits is still represented, because its adapter is needed." For
  types alone the accepted sentence holds: "Only demanded native actuals need
  mappings." (`docs/generation-spi.md:170`). So a request for types does not
  fail for want of something that only conversion needs, as binding
  generation requires: "A native mapping without an adapter facet still
  supports types" (`docs/binding-generation.md:103`). In the protocol this is
  the third need of a request, conversion (O31). Fixed by the prototype
  (report, section 7; no decision names it): "A request names a product by
  its identity."

**Specialization (row 21).** An open declaration is emitted once, so its
layout is fixed. A family whose layout depends on its operand has four
courses inside an open generic (advice, section 2.3): a representation that is
valid for every admitted argument, a requirement on the binder, a native
generic convention that can express the variation, or a specialization that
changes the native representation and the conversion together. The accepted
rule for the last: "Explicit specialization is represented as such
([#166](https://github.com/Bitspark/bittype/issues/166)), never as ordinary
substitution into an unchanged generic declaration."
(`docs/representation-families.md:285-288`).

Row 21: "Explicit specialization (#166) is in the extension domain. The
contract can name an application, its selection and the resulting native
representation. The first implementation refuses a specialization as
unsupported". The test: "A specialization request is refused as unsupported
with the application named." The prototype built the refusal and no design of
#166. In it a specialization's domain names a declaration and its actuals,
each by a source selector, and an actual that is itself an application is
matched by its callee alone (report, F9). What that leaves of #166 is
[not specified here](#what-is-delegated-and-what-is-not-specified-here).

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
  realized after planning. O7 says what it does hold: "the plan carries the
  logical symbol identity and the spelled intent, and the place and the
  imports come from the realization."
- **The link from plan to output is published.** Row 23: "The inventory
  record of #281 gains it as a planner-stated field." It is an addition after
  pull request #281, recorded on
  [#272](https://github.com/Bitspark/bittype/issues/272#issuecomment-6840544720).
- **The realization manifest relates the three.** Row 23: "A serializable
  manifest relates plan, output and realization." O16 names it "the
  realization manifest", apart from the publisher's manifest of row 15a and
  from a contribution's facet list. The advice (section 2.5) lists what it
  identifies: the plan and its content digest, the generation run, the target
  configuration, the placement, and the realized symbol of each logical
  symbol. Who writes it and in which form is
  [not specified here](#what-is-delegated-and-what-is-not-specified-here).
  In the prototype the test composition wrote it, which is a shortcut
  (report, F11).
- **A mismatch is a linking error.** The test of row 23: "A plan read with
  the manifest of another output fails as a linking error, not as an
  incomplete plan."
- **A backend has a fourth input (O15).** Row 23: "the placed result is a
  fourth input of a backend." This amends an accepted sentence in substance:
  binding generation says that a backend works from three inputs alone
  (`docs/binding-generation.md:74`). The naming design had left it open:
  "that page stands as written until
  [#195](https://github.com/Bitspark/bittype/issues/195) decides"
  (`docs/naming.md:729-730`). O15: "Row 23 is that decision: I hold #195, and
  the row is linked to it." It is recorded on
  [#195](https://github.com/Bitspark/bittype/issues/195#issuecomment-6841872859),
  the insertion in binding generation is headed "Amended", and the naming
  design gets a dated note
  ([Amendments](#amendments-to-accepted-pages)). Two more sentences list
  native symbols and imports among what the plan carries
  (`docs/binding-generation.md:56`, `docs/representation-families.md:92`).
  O7 gives both the same clarification.

**Read the realization now (row 24).** A reader takes a native name and place
from the realization, even where the realized name equals the spelled intent.
The naming design: "A consumer learns the realized name and the placed site
from a reader on the formatter's result, keyed by the output's handles."
(`docs/naming.md:668-669`). Row 24 limits "now": "For TypeScript, with the
reader on the formatter's result of #273. Go and Rust are row 24a". Neither
language has a Files stage for generated declarations yet. The prototype
could not run the test of this row: `main` has no reader on the formatter's
result (report, F11). O39: "#273 lands before the gate can be met on `main`
(F11)."

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
- **The request is in the contract and refused until a case needs it (O34).**
  "As explicit specialization is (row 21). Under the delegated form no moved
  contribution needs it. Trigger: the first contribution that needs a helper
  another contribution owns." (O34) In the prototype Go's record names no
  helper, the optional contribution owns its helper, and the helper key that
  a consumer creates on `main` is gone (report, F10).

### The identities in one list

"Leaves the process" means that the identity has a place in the serialized
plan or in the realization manifest. The serialized form of each is
[provisional](#what-is-final-and-what-is-provisional).

| Identity | Who issues it | Leaves the process | Not to be confused with | Row |
| --- | --- | --- | --- | --- |
| Contribution identity: publisher and local key | The publisher states it | Yes | A source selector, a module pin, a profile, a registry's handle, a factory reference | 15 |
| Contribution reference: the identity, the family's meaning and the representation revision | The publisher | Yes | The other four version axes, a digest of an implementation | 15, 16, O4, O19, O24 |
| The family's meaning: a meaning identity and a meaning revision | Not specified where two publishers contribute to one family (O19). The same in every target, and spelled once, in a vocabulary package of the family (O27) | Yes, inside the reference | The contract of seven parts that a family contributes per target (`docs/representation-families.md:32`), the name of one target's contribution, a representation revision | 17, O4, O19 |
| Representation revision | The publisher of the contribution, for that contribution | Yes, inside the reference | The revision of the family's meaning, the version of a payload schema, a runtime compatibility | 17, O4, O19, O24 |
| Profile identity | The profile's owner. A contribution states it | Yes | A contribution identity, a version axis | [Row 1](#what-a-contribution-is) |
| Source identity of a declaration or form | The source integration | Yes | A contribution identity | `docs/representation-families.md:87` |
| Role selector of a field position | The source integration | Not said. Row 11 is decided as D3 (O38), and the report does not say whether the serialized plan carries the selector | A use, a type | 12 |
| A use | The kernel, for one run | No | A use entry | 12 |
| Use entry | The planner | Yes | A memo key, a representation node | 18 |
| Representation node | The planner | Yes | A use entry, the memo token, the native type | 18, 19 |
| Application | The planner | Yes | A specialized declaration, a match of two argument lists by position | 20, 21 |
| Applied representation token | The run-time of the service library | It exists at run time only | A use, an entry number, a contribution identity, a descriptor alone | 22 |
| Logical symbol identity | The planner | Yes | A handle, a spelling | 23 |
| Declaration handle | One sealed output | No | The key of a plan entry | 23 |
| Realized symbol: name and place | The language's Files stage | Yes, through the realization manifest | The naming intent, the spelled name of generation | 23, 24 |
| Key of an owned declaration | The contribution that introduces the declaration | Yes | The declaration's spelling | 25 |
| Facet kind | The owner of the kind | Yes, with its contract version, where the kind states a transport (O36) | The issuer of an owner-issued value of #281 | [26](#what-the-kernel-carries-without-reading-it) |
| Evidence kind: an identity and a version | A vocabulary of the target or of the product (fixed by the prototype) | Not said | A facet kind | 7, O23 |
| Factory reference: module and export | The facet of the contribution that names it | Yes | The contribution identity | 15 |

## What the kernel carries without reading it

**A facet is a typed extension contract (row 26).** Row 26: "A facet is a
kind-identified, versioned, validated value on a contribution." The accepted
elements of a paired contribution are facets in this sense: the native facet,
the value adapter facet, the root facet and further product facets
(`docs/binding-generation.md:97-101`), and a family's lowering facet
(`docs/representation-families.md:65-67`).

A facet kind has, by the advice (section 3.1):

- a stable qualified identity and a version of its contract;
- a typed way to construct or issue a value in the process;
- a validator;
- the categories of reference it may hold, and its compatibility
  requirements;
- a transport descriptor where it crosses a process boundary. Row 26a
  deferred that part until a facet crossed one. It did in the prototype, and
  O36 decides the descriptor, below.

The contract of a kind belongs to the product that the facet serves, so to
that product's readers. The kernel knows no kind.

**Two carriers, and the difference (row 26).** Pull request #281 for
[#272](https://github.com/Bitspark/bittype/issues/272) added owner-issued
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

**The transport descriptor of a facet (row 26a, O36).** The trigger of row
26a fired. O36: "A facet crossed the process boundary: the one that names a
run-time factory." The TypeScript stand-in finds a factory by what the plan
entry names, so the reference to the factory has to be in the serialized plan
(report, section 9). O36 takes the prototype's design:

> "a facet kind may state a media type and an encoder of its owner; when the
> planner writes a selected contribution into the plan it stores, for each
> facet a requested product requires, the kind's identity and contract
> version, the media type and the bytes, and reads none of it; a reader finds
> the facet by the kind's identity, checks version and media type, and
> decodes. A kind without a transport never leaves the process. Row 26a is no
> longer deferred." (O36)

| Who | What it does with a facet that has a transport |
| --- | --- |
| The owner of the kind | May state a media type and an encoder |
| The planner | Calls the encoder for each facet that a requested product requires, and stores the kind's identity and contract version, the media type and the bytes with the contribution's entry. It reads none of it |
| A reader | Finds the facet by the kind's identity, checks the version and the media type, and decodes the value |

In the prototype the value that crossed is the factory reference: a module,
an export and a runtime. The Go lowering facets state no transport and never
leave the process (report, section 9).

**Selected together: a contribution, its facet packages, the bundle (O8, O11,
O16, O32).** "A native facet and the facets a binding product needs are selected
together." (`docs/binding-generation.md:103`). O8 said what is selected:

> "A contribution is a bundle. Its manifest lists its facets, and the
> implementation of a facet may live in another package of the same owner. The
> kernel admits a facet that names a registered contribution of the same
> publisher, registered by the owner's composition." (O8)

Three later decisions correct its words.

- **O11: whose composition.** "It is the consumer's composition." Binding
  generation gives the consumer build the "Selection and composition of the
  concrete contributions it uses" (`docs/binding-generation.md:41`), and "A
  consumer supplies its contributions explicitly."
  (`docs/binding-generation.md:115`).
- **O16: a bundle is not a package.** bittype's rule stands: "One
  contribution, meaning one thing a consumer registers or selects on its own,
  is one package." (`AGENTS.md:54-55`). O16: "A consumer's composition
  registers a mapping package, and it registers a facet package apart from
  it; each is one package. They are selected together. The page calls the
  mapping unit the contribution, a separately registered facet unit a facet
  package of that contribution, and the two together the bundle."
- **O32: a contribution does not state its facets.** O8 said "Its manifest
  lists its facets", and O16 read that as a facet list of the contribution's
  own. The prototype could not build it (report, F15). O32: "A production
  mapping cannot name facet kinds that live in the test tree without spelling
  their identities a second time. The facet list of a contribution is the set
  of facets that the composition registered for it, each checked for the same
  publisher, and checked against what each requested product requires." That
  set meets both accepted sentences: composition must "reject duplicates,
  incompatible profiles, missing facets and unsupported applications"
  (`docs/binding-generation.md:113`), and "A native facet and the facets a
  binding product needs are selected together."
  (`docs/binding-generation.md:103`).

| Word | What it is | Per |
| --- | --- | --- |
| The contribution | The mapping unit: one package that a consumer's composition registers | Itself |
| A facet package | A separately registered package that holds a facet of one contribution | Contribution |
| The bundle | A contribution and its facet packages, which are selected together | Contribution |
| The contribution's facet list | What O8 called "Its manifest". Since O32 it is the set of facets that the consumer's composition registered for the contribution. The contribution states none | Contribution |
| The publisher's manifest | What a composition lock would bind a publisher's name to. Deferred (row 15a) | Publisher |
| The realization manifest | What relates plan, output and realization ([row 23](#identities)) | Output |
| The vocabulary package of a family for a target | Holds the contribution's identity and the generated payload types. The mapping package, its codec and its facet packages import it (O25) | Contribution |

- **Where an implementation lives.** In a facet package of the same
  publisher. A reader may also load it in another process, for example a
  factory by its module and export.
- **What is admitted.** A facet of the same publisher as the contribution it
  names. No other party attaches behaviour by repeating a contribution
  identity (advice, section 3.1).
- **The builtin contributions.** O8: "For the builtin families, the facets
  that the completeness test backend reads live in the test tree, as part of
  bittype's own bundle. A production mapping imports nothing of them, and
  bittype's native vocabularies stay types-only." Since O16 a bundle is per
  contribution, so "bittype's own bundle" is read as the bundles of bittype's
  contributions. The Go conversion and lowering that the gate executes are
  facet packages of the builtin contributions, and they are test code.
- **Where the test tree is (O16, O39).** O16 asks for the role of
  `LAYOUT.md` that holds test-only packages. The layout has a role whose
  packages only tests and test compositions import at `04f9a75`:
  "Contribution: verification law" (`LAYOUT.md:162`). It is no role for a
  facet package of a mapping contribution, and the prototype did not use it,
  "because it is bound to one pair by name" (report, section 10). The gate's
  fixture contributions have a place: "Independent consumer fixtures, stage
  composition, language conformance and native compiler qualification remain
  under integration." (`LAYOUT.md:130-131`). A facet package of a builtin
  contribution is no fixture, and it needs a role: "The checks also classify
  every production package" (`LAYOUT.md:255`). The prototype put the test
  tree under a top-level directory `completeness/` and added roles for it to
  the architecture test. They are a proposal
  ([the roles the prototype proposes](#the-roles-the-prototype-proposes)).
  O39: "The migration's review judges them."
- **Who imports what (O25, O27, O37).** The rule that O16 kept goes on:
  "Identities live in vocabulary packages. Only a composition imports a
  contribution." (`AGENTS.md:55-56`). O25: "So a facet package and a codec
  may not import the mapping package. The contribution's identity and the
  generated payload types live in a vocabulary package of that family for
  that target. The mapping package, its codec and its facet packages import
  that vocabulary, and none imports another." Two decisions add an import
  each. Each target's vocabulary imports the vocabulary package of the
  family, which holds the meaning identity (O27, [identities](#identities)).
  And a codec may import the ontos value packages (O37, below).

**What the kernel guarantees (row 27).** Row 27: "Guarantees 1 to 6 are the
kernel's, each with its test." The advice (section 3.2) says what kind they
are: "These are structural, authority, scope, and association guarantees."

| # | The kernel guarantees that | Its test | Red today because |
| --- | --- | --- | --- |
| 1 | Required facets belong to the selected contribution | Remove a facet that a bindings product requires. The request for bindings fails as "missing facet", and a request for types succeeds. What is checked is the set of facets that the composition registered for the contribution (O32) | A request names no product |
| 2 | A facet's kind, target and compatibility match | A facet from another issuer under the same display name, for another target, or of another contract version is rejected | No kind and no version exist |
| 3 | A reference has the right plan, category and scope | A reference of another plan, a member where a type is expected, and a binder outside its scope are each rejected | A reference has no category |
| 4 | Evidence belongs to the dependency it was asked about | An answer issued for another dependency, another generic environment or another contribution reference is rejected. The reference counts whole, both of its parts (O19). A rule presents an answer through the step it was issued to, and a rejected answer fails the planning of the subject (fixed by the prototype, [evidence](#evidence)) | No evidence path exists |
| 5 | The payload snapshot is immutable | Mutate the value a rule returned, after it returned. Neither the native output nor the stored plan changes. O20: the guarantee "has one moment since O9, after the rule returned" | No record exists |
| 6 | The exact payload passed validation | A fixture family whose schema has a constraint that a Go type cannot express, and a rule that returns a value violating it: the value is refused before any backend runs (O20, derived here). A receipt issued for another snapshot or another schema fails the same way. A receipt is the kernel's record that one frozen value passed validation against one schema. The codec that validates is trusted (O37, below) | Nothing validates a payload |

The tests are the advice's, except where O20 is named. Guarantee 1 rests on
O5: a request that names its products "is what guarantee 1 needs to tell a
missing facet from one that nobody asked for". It checks the facets
registered for the contribution (O32). Guarantee 3 needs the reference
categories of the [next section](#what-lowering-must-be-able-to-refer-to)
(derived here, from verification report 1). Which mutants turn each test red
is under [what the prototype showed](#what-the-prototype-showed).

**What guarantee 6 rests on (O37).** The kernel cannot validate. The
prototype's finding F18: the kernel finds the codec by the schema reference
that the contribution states and trusts it, and a codec that does not
validate is not detected by the kernel. O37: "Guarantee 6 says that a codec
is trusted as a target is: the kernel's part is the association of value,
schema reference and receipt." How the kernel finds the codec is fixed by the
prototype, and the decisions take it as the prototype fixed it (report,
section 7): "The consumer's composition registers it. A codec states its
schema reference, a registry holds one codec for each reference, a
contribution states the schema reference of its payload, and the planner
takes the codec registered under exactly that reference and refuses to plan
without one. The receipt is a digest over the schema reference and the frozen
bytes."

**Two that are not the kernel's (row 27).**

| The advice's guarantee | Whose it is | Source |
| --- | --- | --- |
| Unknown required versions fail explicitly | The planner states the axes, and the reader refuses | Row 27: "the planner emits them and the plan's reader refuses". The test is that of [row 16](#identities) |
| Extension needs no shared case | Nobody's at run time. It is evidence for the extension boundary | Row 27: it "is qualification evidence, not a kernel guarantee". It is line 7 of the [gate](#the-gate) |

**What none of them establishes.** The advice (section 3.2): "They do not
prove that a claimed equality relation is correct, that the projector uses the
intended field, or that the TypeScript factory obeys the conversion laws.
Qualification must test those claims."

**What a rule returns, and when it is frozen (row 28, O9).** Representation
families fixes what a payload is: "A payload is an ontos value. Its family
declares the payload type in horos; bittype generates the Go and TypeScript
types, and horos conformance validates a payload."
(`docs/representation-families.md:98-100`). Row 28 adopts the advice's order
(section 3.3), and O9 decides its technique:

> "A rule returns a typed Go value of a type that bittype generates from the
> family's payload schema. Such a type is plain Go and imports neither horos
> nor ontos. The family's payload codec, a package of the new role, encodes
> that value to an ontos value, validates it against the schema and freezes
> it, at the moment the rule returns, for every request. The projector
> receives the typed view decoded from the frozen value." (O9)

| Step | Who | What |
| --- | --- | --- |
| 1 | The rule | Returns a value of a generated plain Go type |
| 2 | The payload codec of the value's schema | Encodes the value as an ontos value, validates it against the schema and freezes it, when the rule returns |
| 3 | The projector | Receives the typed view that is decoded from the frozen value |
| 4 | Serialization | Publishes that same frozen value |

- **For every generation request.** A request for types alone goes through
  the codec as a request for bindings does. O9 gives the reason: "One
  technique for every request keeps the invariant simple: one immutable
  decision value."
- **A mapping package gains one import (O25).** The generated type is plain
  Go, and it does not live in the mapping package. O25: "A mapping package
  therefore gains one import, of its own vocabulary, and still imports
  neither horos nor ontos. The path and the role entry come with the
  prototype."
- **The test of row 28:** "Mutating the value a rule returned, after it
  returned, changes neither the native type nor the serialized plan."

**The payload codec is a new package role (row 28, O9).** Row 28: "Encoding
and validation live in a new package role outside the contract package, which
may import neither horos nor ontos. The horos source reader is not its home.
`LAYOUT.md` names the home".

- **The role has two parts, and O37 divides them.** A shared part, and one
  codec package for each payload schema. O18: "A family's codec is its own
  package, and a new contribution adds one." O9 says that a family's codec
  encodes, validates and freezes, and O18 calls the shared part "the shared
  encoding and validation mechanism of the new role". O37 extends both: "a
  family's codec may import the ontos value packages, and reaches horos only
  through the shared part, which loads the pinned schema, validates by
  conformance, and encodes and decodes canonical bytes."

  | Part | What it does | What it imports |
  | --- | --- | --- |
  | A family's codec | Two hand-written functions between the generated Go type and an ontos value (report, section 7) | Its own vocabulary, the shared part, and the ontos value packages (O37). Never horos |
  | The shared part | "loads a pinned schema module, resolves the ground descriptor of the payload type, validates by horos conformance, encodes canonical ontos bytes and decodes them" (report, section 7) | horos and ontos |

  So O9's "encodes, validates and freezes" is done by the family's codec
  through the shared part.
- **Why O37.** "The prototype kept codec packages away from ontos and had to
  change the shared part for one family's value form." (O37) The prototype is
  built in the older form: its codec packages import neither library, and the
  shared part exports the value forms a codec builds (report, F26).
- **One accepted sentence waits for the migration.** bittype's `AGENTS.md`
  lists who consumes the model dependencies: "Only the Horos frontend
  (including its text stages), its private model helpers, source-specific
  transformation readers and the concrete Horos verification
  projections/composition consume external Horos/Ontos model dependencies."
  (`AGENTS.md:273-275`). O37: "the codec role is added there with the
  migration." This change does not edit it.
- **For the delivered contributions it is written by hand:** 21 packages
  under six names. O9: "the
  codec is written by hand and pinned with its schema, as row 38 decides for
  the bootstrap" ([payload bootstrapping](#payload-bootstrapping)). Which six
  names is decided by O12, below.
- **The path is a postponed part of row 28.** The row says that `LAYOUT.md`
  names the home. `LAYOUT.md` names the role and no path. The prototype
  proposes `transform/syntax/shared/payload/go` for the shared part and
  `transform/syntax/<source>-to-<native>/codec/<name>/go` for a codec
  ([the roles the prototype proposes](#the-roles-the-prototype-proposes)).
  The migration edits `LAYOUT.md`, and its review judges the proposal (O39).
- **Why it is outside the contract package.** "Generic carriers, contracts,
  engine, printer and native syntax stages have no external Horos/Ontos
  imports." (`CHARTER.md:129-131`). The kernel carries an opaque handle to the
  frozen value and the result of its validation, bound to that value and to
  its schema. It imports neither library.
- **Why it is not the horos source reader.** Derived here, from verification
  report 1: a payload's schema is a horos type whatever the source language
  of the generation is. A source that is not horos plans through the same
  kernel and needs the same codecs.

**The six names of delivered contributions that return a decision (O12,
O21).** A contribution is for one target, so the delivered contributions are
21 packages. They have six names. O9 spoke
of "the six delivered families", and row 2 says that generics are no family
with a contribution. O12 says what the six are:

> "A codec belongs to a payload schema, and a payload schema belongs to each
> delivered contribution that returns a decision. There are six: text,
> booleans, lists, closed records, optional fields, and the declaration
> contribution, which is the part of the division of generics that a
> contribution owns (native parameter order, alias or defined declaration). A
> payload may be empty; the record still exists." (O12)

The packages are those of `transform/syntax/horos-to-<native>/mapping/<name>/go`
at `04f9a75`.

| Package `<name>` | TypeScript | Go | Rust | C++ | Returns a decision | Whose contract the schema belongs to, per target (O21) |
| --- | --- | --- | --- | --- | --- | --- |
| `text` | Yes | Yes | Yes | Yes | Yes | Text |
| `boolean` | Yes | Yes | Yes | Yes | Yes | Booleans |
| `list` | Yes | Yes | Yes | Yes | Yes | Lists |
| `record` | Yes | Yes | Yes | Yes | Yes | Closed records |
| `declarations` | Yes | Yes | Yes | Yes | Yes: the native parameter order, and alias or defined declaration | The declaration contribution |
| `optional` | No package | Yes | No package | No package | Yes, since presence is delegated (O38). Before the move of step 5 Go's package returns none: on `main` it is selected for nothing. Fixed by the prototype, and the decisions take it so | Optional fields |

- **Which payloads are empty.** O12 says that a payload may be empty and that
  the record still exists. Fixed by the prototype, and the decisions take it
  as the prototype fixed it (report, section 7). Empty are text, boolean,
  list and optional, in Go and in TypeScript. Not empty are record, with
  `members: list<text>`, and declarations, with `name: text`,
  `nominal: bool` and `parameters: list<text>`.
- **Six, and five only if row 11 is reopened.** O12: "If presence is not
  delegated (row 11 decided as D1), the optional schema folds into the
  record's and there are five." Row 11 is decided as D3 (O38), so there are
  six.
- **Every target gets an `optional` package.** Derived here, in
  [Optional fields](#optional-fields). On `main` only Go has one. The
  prototype added one for TypeScript. For Rust and C++ the report says
  (section 12): they "decide their optional representation inside the record
  today by an external type; under D3 that moves into a new optional package
  for each."
- **A schema per target (O21).** O12's "six" counts names, and O21 corrects
  it. One of the seven parts of a family's contract for one target is the
  "Selected-use payload": "A family-owned, versioned schema for the decisions
  recorded at each use" (`docs/representation-families.md:39`). O21: "The
  contract is per target, so the schema is too." The packages of one name in
  two targets have two schema references. "The family's owner may publish one
  schema module for several targets where the recorded decisions are the
  same, and then the two references are equal. The contract requires
  neither." (O21)

## What lowering must be able to refer to

A lowering facet turns recorded decisions into target code: "typed target
expressions or helpers, such as generated accessors or constructors for Go",
and "source-code strings never stand in for typed target expressions"
(`docs/representation-families.md:66-69`). Row 29 takes the advice's list
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
| A typed access or construction operation | The vocabulary of the product that the facet serves |
| A helper that the lowering itself introduces | Its key is a reference of the contract ([row 25](#identities)). Its body is in the product's vocabulary. Derived here |
| The value of the run-time call context | The vocabulary of the product that the facet serves |

- **A member is referred to, not spelled.** The advice (section 3.4): "A
  free string alone does not provide that association." The name of a member
  may stay data of the payload. A lowering operation refers to the member and
  to the type that contains it.
- **The contract gains no expression of any language.** The kernel checks the
  category, plan and scope of a reference
  ([guarantee 3](#what-the-kernel-carries-without-reading-it)). What is done
  with a reference is the product's.
- **The set of categories** above is the advice's list read by row 29.
  Fixed by the prototype (report, section 7; no decision names it): "The
  reference categories issued are type, binder, adapter, declaration, owned
  helper and member; external symbol and requirement are not issued. A member
  reference is issued by the plan from the node of the native type; the
  projector states the member names beside the type." So two lines of the
  table above have no issued category yet: an external symbol, and a generic
  requirement with its evidence.
- **A helper that the lowering introduces** is keyed in the prototype's Go
  backend by contribution reference and name, and the backend issues its
  spelling. Whether that key is the owned key of [row 25](#identities) was
  seen and not assessed (report, section 13).

**What executed Go lowering shows, and what it does not (O33).** The accepted
sentence is quoted above: source-code strings never stand in for typed target
expressions. The prototype's lowering does not meet it (report, F14). O33:
"The prototype's Go lowering is text with typed holes: every type, member,
helper and adapter argument enters through a hole that holds a checked
reference of the plan, and the text between the holes is the facet's own Go.
So line 5 shows the contract's part, the reference categories and their
checks, and executed code. It does not show a typed expression vocabulary for
Go, which is the reading product's vocabulary and belongs with value-level Go
to the future Go runtime library." And: "The instrument stays a test."

**Where the prototype's value vocabulary lives (row 29).** "The prototype's
Go value vocabulary and its facets stay in the test tree". The reason is the
accepted scope: "bittype's native vocabularies stay types-only."
(`docs/binding-generation.md:87`). O8 keeps both: the facets that the
completeness test backend reads are facet packages of bittype's contributions
and live in the test tree, and a production mapping imports nothing of them.
The prototype put it at `completeness/<native>/values/go`
([the roles the prototype proposes](#the-roles-the-prototype-proposes)).
The prototype's vocabulary is a test instrument. The advice (section 3.4)
says what it must do: "It must compile and run generated construction, member
access, an owned helper, and a generic adapter argument." The production
owner of value-level Go is not bittype: "A Go backend belongs to the future
Go runtime library." (`docs/binding-generation.md:46`).

## The extension domain

**What "changes once" means (row 33).** Representation families says that
what an existing mapping returns changes once: "After them a family adds a
contribution and changes no existing contract."
(`docs/representation-families.md:359`). Row 33 adopts the advice's wording
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
| A further family | In. Derived here | Every family of the delivery table of representation families (`docs/representation-families.md:372-382`) enters as contributions. The accepted promise: "Adding a family through the published identity, dependency, invocation, plan-graph and target-extension contracts changes only that family and the consumers that compose it." (`:22-24`) | A use of a family that is not built is refused as "missing contribution": there is no contribution to select (O17) |
| Further products and facets | In | A further backend for a target, and a further kind of facet, enter through the [facet contract](#what-the-kernel-carries-without-reading-it) with no edit of the kernel | A request for a product whose facet a contribution lacks fails as "missing facet" |
| A contextual representation at a field position | In. Presence is delegated (row 11 decided as D3, O38) | The request of [row 12](#optional-fields) and the member-fragment answer of [row 13](#optional-fields). O10: "Row 33's line on a contextual representation at a field position is conditional on that", which is the decision of row 11. If an objection reopens row 11 and presence becomes the record's layout, the item is out. That the contract can express a contribution without a value adapter of its own ([row 14](#optional-fields)) holds either way | A fragment that the target's member vocabulary does not admit, as unsupported |
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

**The partition (row 30).** The evaluation labels the partition that the
submitted research document proposed "contract final, envelope unstable". Row
30 withdraws it: "The partition of the head of
section 6 is replaced by the advised one: meanings final, encoding and layout
provisional". The reason is that the contract and the envelope share their
identities, their categories of dependency, their generic environments,
their requirements and their references to symbols (advice, section 4.1).

| Must have its intended meaning. Final with the gate | Stated under | May stay provisional or change under a version |
| --- | --- | --- |
| The meaning and the ownership of each identity | [Identities](#identities) | The layout and the encoding of the serialized plan |
| The categories and the scope of a request and of its answer | [The protocol](#what-a-rule-does-the-protocol), [Optional fields](#optional-fields) | Caches and performance strategies |
| The generic conventions and requirements | [Identities](#identities), [Evidence](#evidence) | Coverage of further families |
| The distinction between a reference and a demand for evidence, and the meaning of a refusal | [Evidence](#evidence), [Refusals](#refusals) | Revisions of a family's payload schema |
| The validation of a facet and what it is attached to | [What the kernel carries](#what-the-kernel-carries-without-reading-it) | Further implementations of a product |
| The relation between a native symbol and its realization | [Identities](#identities) | The presentation of diagnostics |
| The inputs that lowering needs | [What lowering must be able to refer to](#what-lowering-must-be-able-to-refer-to) | Optimizations of lowering |

The two outer columns are the advice's table, with one change in the fourth
row of the left column. A row of the left column is not paired with the entry
beside it.

**Recursive references are not final with the gate (O17).** The advice's
table has "Recursive references and refusal meaning" in the left column, and
row 30 took it over. O17 takes recursive references out: "since O6 the
prototype shows nothing of recursion beyond its refusal, and I refused the
general request loop for the same reason, that it would be final without
evidence. It becomes final with #164." What the gate does fix of recursion is
what row 9 puts into the contract: a rule names a reference apart from a
demand for evidence, and both a recursive reference and an evidence cycle are
refused.

- **Unstable does not mean without meaning.** The advice (section 4.1):
  "Calling the transport unstable does not excuse changing the meaning of an
  admitted identity later." The accepted sentence stands and is read so: "The
  envelope keeps an explicitly unstable version until the completeness test
  backend passes" its six cases (`docs/representation-families.md:109-110`).
  The version is what is unstable. It gets a dated clarification
  ([Amendments](#amendments-to-accepted-pages)).
- **Decided is not demonstrated.** Derived here, from AC-07: a row of the
  left column is decided by this page and demonstrated only by its test in
  the gate. The relation between a symbol and its realization is demonstrated
  for TypeScript only, until Go and Rust have a Files stage (row 24a).
- **Until the gate is met** nothing on this page is final. The prototype may
  change any of it, and each change is recorded here with its date. The
  prototype's changes are the decisions O27 to O39 of 13 November 2026 and
  the forms that this page marks "fixed by the prototype".
- **After the gate** the left column is fixed, by AC-10: a first change
  "establishes the final structure: its inputs and identities have their
  intended shape". The right column changes under its version axes
  ([row 16](#identities)).

## Inputs that are not bittype's alone

**The source semantics, for a reader in another process (row 34).** A
backend's second input is "the source semantics"
(`docs/binding-generation.md:26`), and it emits "ground-descriptor constants"
(`docs/binding-generation.md:81`). No serialized form of that input exists.
Serializing the plan alone therefore shows nothing about a reader in another
process. Row 34: "For the prototype: a snapshot of ground descriptors and
module pins exported from the Go side, whose source identities match the
plan's. A general reader protocol is row 34a".

- **What the snapshot holds (O29).** Row 34 says "ground descriptors". The
  nodes of the plan describe the open declaration `Page<T>`, and both readers
  of the prototype needed the canonical name of each field of its body
  (report, F16). O29 words the snapshot of row 34 as "the descriptors of the
  declarations as the source states them, open where the declaration is
  generic". Its reason: "A ground descriptor does not exist for `Page<T>`",
  and open generics remain open. The ground-descriptor constants that a backend
  emits (`docs/binding-generation.md:81`) are for ground uses.
- **How a reader pairs a link with a child (O39).** "A link carries the
  provider's locator of its child, passed through unread, and the snapshot
  gives the children of any descriptor by position, so that a reader has no
  case for a source form (F17, F27)." The accepted page has the kernel carry
  "provider-expression operands, with locators passed through unread"
  (`docs/generation-spi.md:200`). The prototype built the first half. Its two
  readers still have a case for the descriptor kinds record and list, which
  the report calls a shortcut (F27).
- **What a reader may do with it.** The advice (section 5.1): "The backend may
  read source meaning; it must not reconstruct representation choices."
- **A snapshot that does not match is a linking error.** Derived here, from
  verification report 3: a snapshot of a module with the same name and shape
  and another pin is refused, and not read as an incomplete plan. Binding
  generation lists the case among its falsification tests: "a changed module
  digest with the same name and shape" (`docs/binding-generation.md:258`).

**Target configuration (row 35).** The plan records the configuration of the
target that affects selection, what is admissible natively, or whether
generated code is compatible. The envelope has the place: "Version axes;
target and its configuration" (`docs/representation-families.md:86`).

For TypeScript one option matters to optional fields. bittype checks its
generated TypeScript under `exactOptionalPropertyTypes`. That option changes
what the compiler admits. It does not change that JavaScript tells an absent
property from a present one (advice, section 5.2). Row 35: "The optional
family's contract states it, and the four states are tested apart."

| State of an optional field | Tested |
| --- | --- |
| Absent | By its own vector, in Go and in TypeScript. In the prototype it converts there and back as itself on both sides |
| Present with a value | By its own vector, in Go and in TypeScript. In the prototype it converts there and back as itself on both sides |
| Present with `null`, where the type admits it | Derived here: none of the delivered families admits `null`, so this state enters with [#211](https://github.com/Bitspark/bittype/issues/211). In the prototype a present `null` is refused by the operand's adapter on both sides |
| Present with `undefined` | Refused, by its own vector, in TypeScript (O14). The state does not exist in Go, "where a member is a value of its type" (report, section 8) |

- **A present `undefined` is refused (O14).** Two accepted sentences settle
  it. Who decides: "The representation owner owns its factories' semantic
  contract and qualification: bittype's structural profile for the structural
  families" (`docs/representation-families.md:74-76`). And what the adapter
  does: "TypeScript uses `field?: T | null`, and the adapter separates an
  absent own property from a present `undefined` at run time."
  (`docs/representation-families.md:198-199`). So reading a present
  `undefined` as absence was never an option of the profile. O14: "A present
  `undefined` is a present field whose value is no value of the field's type.
  Conversion refuses it, and the vectors pin the refusal."
- **What was asked, and the correction.** Row 35 says of the difference
  between the two repositories: "That difference is put to its owner". The
  first question on bitspark-service#38 offered two behaviours, refused or
  read as absent. O14 calls that framing wrong, and the
  [correction of 13 November 2026](https://github.com/Bitspark/bitspark-service/issues/38#issuecomment-6841872432)
  on bitspark-service#38 states the decision and asks that owner for
  objections. They are weighed by the research owner. The record adapter of
  the service library reads a present `undefined` as absence. O14: "that is a
  difference from the accepted sentence, in code that is not yet the profile's
  package."
- **"Round-trip apart".** The test of [row 13](#optional-fields) says that
  the four states "round-trip apart through the record factory". Derived
  here: a present `undefined` has no canonical value of its own, so it does
  not come back as itself. It is refused. The test is read with row 35: each
  state has its own vector and its own expected outcome. The
  decisions of 13 November 2026 confirm the reading: "the four presence states
  are tested apart and a present `undefined` does not round-trip as itself".

**What goes to the service library (rows 36 and others).** Row 36 is the one
stale row: "Stale for bittype: laws K1 to K4 of the binding contracts state it
and bittype confirmed them on 12 November. The two points they leave unstated
go to the service library". All of the following is in the
[comment of 13 November 2026](https://github.com/Bitspark/bitspark-service/issues/38#issuecomment-6840511750)
on bitspark-service#38, and question 2 in its
[correction](https://github.com/Bitspark/bitspark-service/issues/38#issuecomment-6841872432).
The grade says of the three questions: "None of them blocks the prototype."

| What | Row | Form |
| --- | --- | --- |
| A field-level conversion slot for an optional record field | 11, 13, O38 | Question 1. No answer has come. Its shape is put to that owner in a [comment of 13 November 2026](https://github.com/Bitspark/bitspark-service/issues/38#issuecomment-6844432275), with the decision of row 11 and what reopens it |
| What conversion does with a present `undefined` | 35, O14 | Question 2, then corrected: it is decided as refused, and that owner is asked for objections |
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
  second, and row 37 does not. The decisions of 13 November 2026 confirm that
  one route replaces the two.
- **Its test.** Derived here, from verification report 3: a Go module outside
  this repository authors and registers a contribution. It is red when the
  module needs a symbol that is internal or that lives in the test tree.
- **That contribution states native types only (O18).** "That is intended for
  now: every binding facet kind and the Go value vocabulary live in the test
  tree, and a production home for them belongs to the first backend's owner
  (rows 31a and 37a)."
- **It makes its own bootstrap artifact (O39).** "The out-of-tree
  contribution of line 7 makes its bootstrap artifact with public packages
  only. Today it cannot: the canonical module, the pin and the generated type
  need an internal package (F25)." In the prototype the module outside the
  repository registers a contribution and a codec with public packages only.
  The test wrote its canonical module and pin with an internal package and
  its type by hand, and it shows no generation with that contribution
  (report, stage 9 and F25). The public way from a schema source to the
  pinned artifact is
  [not specified here](#what-is-delegated-and-what-is-not-specified-here).
- **The first contribution of the service library** in its own repository is
  deferred (row 37a).

## Payload bootstrapping

Row 38: "A pinned bootstrap artifact: the schema source, its canonical module
with its pin, and the committed generated types with a regeneration test".
O9 adds the codec to it: "the schema source, its canonical module, the
committed generated types, the codec, and a test that regenerates the types
and round-trips every value".

bittype generates the types of a payload from a horos schema, with the
families whose payloads those types describe. The families that exist cannot
type every field a payload wants. The advice (section 5.5): "The limitation
to six currently supported families is a bootstrap issue, not a reason to
encode unsuitable payload semantics."

- **One schema language.** "No separate schema language is introduced."
  (`docs/representation-families.md:100-101`). The advice repeats it: "Do not
  introduce a second schema language."
- **A field the families cannot type yet.** Derived here, from verification
  report 3: such a field is carried by structure where that says the same,
  for example by the order of a list where an index was wanted. Otherwise it
  is recorded as a bootstrap exception that names the family that ends it.
- **Two kinds of later change.** The advice (section 5.5): "Adding a semantic
  integer field later is normally a family schema/version change. Discovering
  that a payload has been hiding an orchestration reference is a contract
  defect." The second is ruled out by representation families: "Nothing the
  orchestration must resolve hides inside a payload."
  (`docs/representation-families.md:96-97`).
- **What the prototype has of it.** Eleven vocabulary packages, each with the
  schema source, its canonical module and pin, the committed Go type and a
  hand-written codec. No TypeScript payload type is generated, though
  representation families gives bittype both: "bittype generates the Go and
  TypeScript types" (`docs/representation-families.md:99-100`). The test was not
  shown red (report, stage 1).
- **Its tests.** The regeneration test fails when the committed types differ
  from what the schema generates. A payload that does not conform to its
  schema is refused before any backend runs
  ([guarantee 6](#what-the-kernel-carries-without-reading-it)).

## What a test must show

Row 39: "The boundary is enforced by what a projector can reach. Mutants must
disagree on presence, layout, slot identity or environment."

**The mutant that proves nothing.** Row 2 of research 0015 validated the one
origin of a decision by "a seeded mutant that builds the plan entry in a
separate pass from the native expression". A second pass that computes the
same decisions from the same source is an equivalent mutant. Every
behavioural test stays green on it.

**What replaces it.** Three layers. The first two are row 39's: its two
sentences. The examples of mutants in the second are verification report 1's.
The third is derived here, from that report.

| Layer | What it shows | How |
| --- | --- | --- |
| Structure | No second decision point is reachable | The projector cannot reach the source reader, and a reader of the plan has nothing to select with |
| Mutants that disagree | The tests detect a wrong implementation | Each mutant differs in one thing: presence in the entry and not in the native type, the order of fields in one of the two, two crossed links, an adapter argument tied to the wrong binder |
| A decision that the source does not determine | A second origin would be visible | A fixture contribution takes a member's name from its own configuration. A pass that rebuilds the entry from the source then differs from the native type |

The first layer prevents an unauthorized decision point. The other two
detect an incorrect implementation. AC-07 keeps the two kinds of evidence
apart.

**The texts that change.** Row 39: "#207's verification and row 2 of research
0015 are reworded". For #207 the
[comment of 13 November 2026](https://github.com/Bitspark/bittype/issues/207#issuecomment-6840545162)
on that issue states the new verification. The validation cell of row 2 in
the table of [research 0015](../research-docs/0015-representation-families.md)
is reworded with this change, dated, and keeps its former text readable. The
text of #207 itself was reworded when this page merged
([the notes posted with the merge](#what-is-delegated-and-what-is-not-specified-here)).

## Amendments to accepted pages

**One amendment in substance, and clarifications (O15).** Row 23 amends one
accepted sentence in substance: binding generation says that a backend works
from three inputs alone, and row 23 adds a fourth. Its insertion is headed
"Amended". Every other insertion is headed "Clarified" or "Noted" and leaves
the meaning of its sentence. No accepted sentence is rewritten: each stays
readable as it was accepted, and the dated insertion stands beside it. Two
insertions do not stand beside their sentence and name it: the one for
"consumes three inputs" stands after the completeness test, and the one for
the memo sentence stands with the paragraph on the instance. All of them are
made in the change that adds this page.

**Who was told.** The holder of research 0014, `claude-bittype-gaps`, was
told on the board on 13 November 2026 of the amendment and the insertions in
binding generation, and of the reworded cell in the table of research 0015.
No answer had come when this was written. Objections of that seat are
weighed.

| Page and sentence | What the clarification says | Why | Row | Who owns the page |
| --- | --- | --- | --- | --- |
| [Binding generation](binding-generation.md#order-of-work-and-qualification): "The unit of implementation is library × target." (`:269`) | The sentence names ownership and distribution. The unit that a consumer selects is one contribution, and a library may publish several | Read as the unit of selection, the sentence would make every rule of a library one indivisible choice | 1 | Research 0014, held by `claude-bittype-gaps` |
| [Binding generation](binding-generation.md#the-selected-representation-plan): "The backend can implement each selected representation from the plan, the source semantics and the selected contributions alone, without guessing a rule." (`:74`), and "consumes three inputs" (`:23`). **Amended in substance** | For a generated native declaration the three inputs are not sufficient. Its name and place exist only after a Files stage has placed it. A plan entry refers to it by an identity the planner issues, and a backend reads the realization from the placed result, which is a fourth input | The [naming design](naming.md#reconciliation-with-accepted-documents) stated the tension and left it to #195, which decided it with row 23 | 23, O15 | Research 0014, held by `claude-bittype-gaps` |
| [Naming](naming.md#reconciliation-with-accepted-documents): "that page stands as written until [#195](https://github.com/Bitspark/bittype/issues/195) decides" (`:729-730`) | A dated note: #195 decided on 13 November 2026, and binding generation carries the amendment | The sentence waited for that decision | 23, O15 | Research 0016. This chain holds it |
| [Binding generation](binding-generation.md#the-selected-representation-plan): the plan carries "the rules for absence, presence, null and native `undefined`" (`:58`) | The plan carries a rule either as a recorded field, where it varies between selections, or through the contribution reference whose revision fixes it | Row 5 keeps an invariant out of the payload, and both sentences hold | 5, O1 | Research 0014, held by `claude-bittype-gaps` |
| The same list: "native symbols, distinguishing run-time imports from type-only imports" (`:56`) | The plan carries the logical symbol identity and the spelled intent of a generated declaration. Its place and the imports come from the realization | The same reason as for the fourth input | 23, O7 | Research 0014, held by `claude-bittype-gaps` |
| [Representation families](representation-families.md#the-plan): "Native symbols, imports, factory references, runtime requirements" (`:92`) | The same clarification as the row above | The same | 23, O7 | Research 0015. This chain holds its delivery |
| [Representation families](representation-families.md#the-representation-contract): "so a native declaration and its conversion cannot disagree." (`:46-47`) | "Cannot disagree" means that no sanctioned second decision point exists. It does not mean that a projector or a factory cannot be wrong | Read as a claim about every implementation, the sentence contradicts AC-05 | 6 | Research 0015. This chain holds its delivery |
| The same paragraph: "one immutable representation instance per use" (`:44-45`), and "The memo is keyed by the selected representation instance" (`:293`) | "Instance" has two senses. In the first sentence it is the selection at one use, a use entry. In the second it is the applied representation that several uses share | Read with one sense, every use would have its own memo key | 18 | Research 0015 |
| [Representation families](representation-families.md#the-plan): "The envelope keeps an explicitly unstable version" (`:109-110`) | What is unstable is the version: the layout and the encoding. The meanings that the plan shares with the contribution contract are fixed with that contract | The partition "contract final, envelope unstable" is withdrawn | 30 | Research 0015 |
| [Representation families](representation-families.md#delivery): "After them a family adds a contribution and changes no existing contract." (`:359`), and "generics" among the delivered families (`:374`) | "Changes once" is a commitment for the written extension domain. "Generics" names a division of responsibility and no family with a contribution of its own | Rows 33 and 2 | 33, 2 | Research 0015 |
| [Package boundaries](package-boundaries.md#findings-on-main-at-3bf45c3), finding 7 (`:159`) | The finding is bittype's derivation. It was conditional on the answer of the service library until row 11 was decided. Since O38 it stands, and an objection of that owner before the migration merges reopens it | Row 11 and O38 | 11, O38 | bittype, case bittype/ACX-009 |

**The edit of binding generation is three insertions.** Pull request
[#260](https://github.com/Bitspark/bittype/pull/260) of `claude-bittype-gaps`
edits the same page and is open. This change adds one paragraph after the
list of what the plan carries (O1 and O7), one after the completeness test
(row 23, the amendment) and one after the sentence on library × target (row
1). It touches
no line that #260 changes. In representation families it adds four
clarifications, none in a passage that #260 changes.

**Pointers, not amendments.** Seven pages gain a pointer to this page or a
dated note, and change no statement of their own: the [SPI](generation-spi.md) and the
[registration page](kernel-composition.md) say that the mapping contract is
going to change, [semantic generation](semantic-generation.md),
[package boundaries](package-boundaries.md) and
[bittype/ACX-009](architecture-cases/0009-package-boundaries.md) said that
the optional move is conditional, `LAYOUT.md` names the new package role
without a path, and the [documentation map](README.md) lists this page.

**After the prototype (O38, O39).** The change that writes the prototype's
results into this page adds dated notes and edits no accepted sentence. In
[package boundaries](package-boundaries.md), in
[bittype/ACX-009](architecture-cases/0009-package-boundaries.md) and in
[semantic generation](semantic-generation.md) the notes on the optional move
say that row 11 is decided as D3 and what reopens it. The paragraph of
`LAYOUT.md` on this design says that the prototype proposes roles and where
they are listed. Its role table is not edited: that is the migration's, whose
review judges the proposal (O39).

**Sentences that are not amended now.** Each needs its owner or an answer
first. They are listed where they arise:

| Sentence | Why not now | Where |
| --- | --- | --- |
| The three sentences that a field-level slot would change | Row 11 is decided (O38). No decision says when they are annotated, and an objection before the migration merges would reopen row 11 | [Optional fields](#optional-fields) |
| "Only the Horos frontend (including its text stages), its private model helpers, source-specific transformation readers and the concrete Horos verification projections/composition consume external Horos/Ontos model dependencies." (`AGENTS.md:273-275`) | The codec role is added there with the migration (O37) | [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| "A single generic convention relates source binders, native parameters and complete mapped actuals." (`docs/generation.md:45-47`) | The page describes type generation as it is implemented. It changes with the migration, and this page reads the sentence as the native convention | [Identities](#identities) |
| The memo sentence of the binding contracts | It is the service library's | [What is delegated](#what-is-delegated-and-what-is-not-specified-here) |

## Deferred, with triggers

Eight rows are deferred. A ninth, row 26a, was deferred until 13 November
2026, when its trigger fired in the prototype. A trigger is a state of a
repository that anyone can check.

| Row | What is deferred | The trigger that ends the deferral | Tracked |
| --- | --- | --- | --- |
| 9a | Recursive representation references | "#164 is claimed and vertical case 1 is being built." Until then recursion is refused: `grep -n "recursive native declaration" contracts/generation/go/planning.go` finds the refusal | [#164](https://github.com/Bitspark/bittype/issues/164) |
| 15a | A composition lock that binds a publisher's name to the publisher's manifest. A contribution's facet list is not deferred (O16, O32) | "a second publisher registers a contribution (bitspark-service#39 item 2 starts)" | bitspark-service#39 |
| 17a | A root facet of `Service` that recognizes a family's meaning | "bitspark-service#39 starts its root facet. No root facet exists in either repository" | bitspark-service#39 |
| 24a | A reader of the realization for Go and Rust | "slice 1 of #161 or #159 delivers a generated-Files formatter (`ls format/files/golang format/files/rust`)" | [#161](https://github.com/Bitspark/bittype/issues/161), [#159](https://github.com/Bitspark/bittype/issues/159) |
| 26a | A transport descriptor for a facet. **No longer deferred** (O36) | "the first facet crosses a process boundary in the prototype of row 31". It fired: the facet that names a run-time factory crossed. The descriptor is [decided](#what-the-kernel-carries-without-reading-it) | #196 |
| 31a | TypeScript end to end through the service library's own factory | "`packages/structural-adapters` exists on bitspark-service `main`, or its owner records on bitspark-service#38 that a prototype precedes integers" | bitspark-service#38 |
| 31b | Vertical case 1 as a whole, and vertical case 2 in real form (O6) | "#164, #210 and #211 deliver their families" | [#164](https://github.com/Bitspark/bittype/issues/164), [#210](https://github.com/Bitspark/bittype/issues/210), [#211](https://github.com/Bitspark/bittype/issues/211) |
| 34a | A reader protocol for source semantics | "the design of bitspark-service#40 item 3 starts" | bitspark-service#40 |
| 37a | A contribution that is authored in the service library's repository | "bitspark-service#39 item 2 starts" | bitspark-service#39 |

Three remarks, derived here.

- **Row 26a ended inside the prototype.** The stand-in reader of TypeScript
  reads a serialized plan in another process. The facet that names a factory
  crossed the boundary there, the prototype designed the transport
  descriptor, and O36 takes it
  ([what the kernel carries](#what-the-kernel-carries-without-reading-it)).
- **Row 31a is why a stand-in stands in.** What it proves and what it does
  not is under [The order of work](#the-order-of-work).
- **A deferral is no narrowing.** Each deferred part is in the
  [extension domain](#the-extension-domain) or belongs to another owner. None
  is removed from what the accepted pages promise.

**Two things that wait for a trigger by a decision, and not by a row.**

| Decision | What waits | The trigger |
| --- | --- | --- |
| O34 | The request for a shared helper. It is in the contract and refused until then ([identities](#identities)) | "the first contribution that needs a helper another contribution owns" |
| O35 | Which facts a nominal declaration restates. Until then it answers `unknown` ([evidence](#evidence)) | "the first family that needs evidence through one, #209 (maps and sets)" |

## Points decided after the first version

No point is open between two rows, or between a row and an accepted page.
This section lists what was open and how it was decided. What the prototype
left unsettled is under
[Not specified here](#what-is-delegated-and-what-is-not-specified-here), and
where its report and a later decision differ is under
[What the prototype showed](#what-the-prototype-showed).

**O1 to O10.** The first version of this page listed ten points on which two
rows, or a row and an accepted page, did not fit without a reading. It gave
the options and chose none. The research owner decided all ten in the
[decisions of 13 November 2026](https://github.com/Bitspark/bittype/issues/196#issuecomment-6841190077)
on #196. Three of them correct a row of the research table.

**O11.** One further point came from the wording of O8 and was decided
[the same day](https://github.com/Bitspark/bittype/issues/196#issuecomment-6841518102).

**O12 to O20.** An independent reader then read the page at `cab8faf` against
the table, the eleven decisions and the accepted pages, and reported 25
findings. Twelve needed a decision, and the research owner made
[nine](https://github.com/Bitspark/bittype/issues/196#issuecomment-6841866946).
The other findings are corrections of wording, standing and form, which this
page takes as the reader proposed.

**O21 and O22.** Writing O12 and O19 into the page showed that each said
more than `docs/representation-families.md` allows, and the research owner
[corrected both](https://github.com/Bitspark/bittype/issues/196#issuecomment-6842186587).

| | The point | The decision, in short | Stated under |
| --- | --- | --- | --- |
| O1 | An invariant is carried by reference, and the plan "carries" the rules | Both hold. A rule is carried as a recorded field or through the contribution reference | [The protocol](#what-a-rule-does-the-protocol), [Amendments](#amendments-to-accepted-pages) |
| O2 | Row 8 called "refuted" an accepted name, and the accepted page has `overlap` | Row 8 is corrected. Evidence has proved, refuted, unknown and exhausted, and `overlap` stays the name for a union's disjointness | [Evidence](#evidence) |
| O3 | A member's name leaves a rule in a request, before the record is frozen | A request may state what the rule knows from the source. Resolved requests are part of the frozen record | [The protocol](#what-a-rule-does-the-protocol), [Optional fields](#optional-fields) |
| O4 | One semantic string for all targets, and an invariant of one target | Row 17 is corrected. Two identities: one that is the same in every target, and a contribution's own representation revision. O19 names the first | [Identities](#identities) |
| O5 | A dependency that only conversion needs could fail a request for types | A generation request names its products. Row 20 applies when bindings are asked for | [Identities](#identities), [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| O6 | Vertical case 1 at the kernel's level while recursion is refused | Row 31b is corrected. Case 1 is not in the prototype | [The order of work](#the-order-of-work), [Deferred](#deferred-with-triggers) |
| O7 | "Native symbols" in two sentences that row 23 does not name | Both get the clarification of the fourth input | [Identities](#identities), [Amendments](#amendments-to-accepted-pages) |
| O8 | Where the Go facets that the gate executes are declared | The facets of the builtin contributions live in the test tree. O11 and O16 correct two of its words | [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| O9 | How a decision becomes an ontos value | A rule returns a value of a generated plain Go type. A payload codec encodes, validates and freezes it | [What the kernel carries](#what-the-kernel-carries-without-reading-it), [Payload bootstrapping](#payload-bootstrapping) |
| O10 | Whether the migration waits for the answer on optional fields | The round of requests is in the contract either way. The gate requires row 11 decided | [Optional fields](#optional-fields), [The extension domain](#the-extension-domain), [The gate](#the-gate) |
| O11 | Whose composition registers a facet | The consumer's | [What the kernel carries](#what-the-kernel-carries-without-reading-it), below |
| O12 | What "the six delivered families" are, when generics are no family | Six names of delivered contributions (21 packages) that return a decision, the declaration contribution among them. Five if presence is not delegated | [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| O13 | Where an answer of evidence comes from | From what the dependency's selected contribution stated or its revision fixes. No provider computes it | [Evidence](#evidence) |
| O14 | Who decides what a present `undefined` does | bittype's structural profile. It is refused | [Inputs](#inputs-that-are-not-bittypes-alone) |
| O15 | The page said it amends no accepted page in substance | Row 23 amends one sentence in substance, and it is the decision of #195 that the naming design waited for | [Identities](#identities), [Amendments](#amendments-to-accepted-pages) |
| O16 | A bundle across packages against the rule that one contribution is one package, and three things called a manifest | A bundle is a contribution and its facet packages. Three names for the three lists | [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| O17 | The class of an unbuilt family, the closed set of classes, and recursive references final without a test | "Missing contribution". The set closes with the prototype. Recursive references become final with #164 | [Refusals](#refusals), [What is final](#what-is-final-and-what-is-provisional) |
| O18 | What the gate leaves to guess on `main` | Fixtures for the delegated form, the shared part of the codec role, a native-only contribution out of tree, line 7 shown twice | [The gate](#the-gate) |
| O19 | Which revision evidence binds to, and the name "family contract" | The whole contribution reference. The name is "the family's meaning" | [Identities](#identities), [Evidence](#evidence) |
| O20 | Rows without a test, and two tests that could not fail | Tests for rows 1, 5, 17 and 22, and new tests for guarantees 5 and 6 | [The gate](#the-gate), [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| O21 | Whether the packages of one name share one payload schema | O12 is corrected. A schema belongs to a family's contract for one target; the family's owner may publish one module for several | [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| O22 | Whether the native equivalence is part of the family's meaning | O19 is corrected. The native equivalence and a normalization belong to the representation revision. O24 corrects the rest | [Identities](#identities) |
| O23 | The planner cannot read a fact in a value that the kernel carries unread | A fact is a typed value of an evidence kind, stated beside the decision or declared at registration. The kernel carries it unread | [Evidence](#evidence), [The protocol](#what-a-rule-does-the-protocol) |
| O24 | Which domain is the same in every target, and whether the representation revision fixes two version axes | O22 is corrected. The family's meaning is `V_T` with its predicates. Coverage belongs to the representation revision. The axes of row 16 stay apart | [Identities](#identities) |
| O25 | Where the identity and the payload types live, when only a composition imports a contribution | In a vocabulary package of the family for the target. A mapping package gains one import | [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| O26 | The four states under D1, the test of row 1, and one refusal value or two | The four states stay in the gate. The test compares payloads, and frozen records only where nothing depends on the replaced contribution. Missing and ambiguous are two values | [The gate](#the-gate), [Refusals](#refusals) |
| O27 | The prototype spelled a family's meaning identity once for each target | A vocabulary package of the family, without a target, holds it, and each target's vocabulary imports it. It extends O25 | [Identities](#identities), [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| O28 | While the body of a generic declaration is planned, the native meaning of a binder is not known | The kernel carries a bound use as a scoped reference until the declaration's projector has returned its convention | [What a contribution is](#what-a-contribution-is) |
| O29 | Three sentences that could not be built or shown as written | The coherence test is the family's and no refusal of the kernel. The snapshot holds the declarations' descriptors as the source states them. The name of a delegated field is in its request | [Evidence](#evidence), [Inputs](#inputs-that-are-not-bittypes-alone), [Optional fields](#optional-fields), [The gate](#the-gate) |
| O30 | What makes a cycle an evidence cycle | Every open request on it demands evidence. With one reference on it, it is a recursive reference | [Evidence](#evidence) |
| O31 | A child that only conversion needs | A third need of a request: conversion | [The protocol](#what-a-rule-does-the-protocol) |
| O32 | A production mapping cannot name its facet kinds | A contribution states no facets. Its facet list is what the composition registered for it. It corrects O8 and O16 | [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| O33 | The prototype's Go lowering is text with typed holes | Line 5 shows the reference categories, their checks and executed code. It shows no typed expression vocabulary | [What lowering must be able to refer to](#what-lowering-must-be-able-to-refer-to), [The gate](#the-gate) |
| O34 | The request for a shared helper is stated and not built | It is in the contract and refused until a case needs it | [Identities](#identities), [Deferred](#deferred-with-triggers) |
| O35 | Evidence does not pass through a named declaration by itself | The kernel forwards the facts of an alias's body. A nominal declaration answers for itself until #209 decides more | [Evidence](#evidence), [Deferred](#deferred-with-triggers) |
| O36 | A facet crossed the process boundary | The transport descriptor that the prototype designed. Row 26a is no longer deferred | [What the kernel carries](#what-the-kernel-carries-without-reading-it), [Deferred](#deferred-with-triggers) |
| O37 | The shared part of the codec role changed for one family's value form | A codec may import the ontos value packages and reaches horos only through the shared part. A codec is trusted as a target is. It extends O9 and one point of O18 | [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| O38 | Row 11, with no answer of the service library | Presence is delegated (D3). An objection with a reason, before the migration merges, reopens it | [Optional fields](#optional-fields) |
| O39 | Six points for the migration | A link carries its child's locator. One composition for each process. A public route to the bootstrap artifact. The fourth axis is a reader's test. #273 lands first. The prototype's roles are a proposal | [Inputs](#inputs-that-are-not-bittypes-alone), [The gate](#the-gate), [Identities](#identities), [What the prototype showed](#what-the-prototype-showed) |
| O40 | Who owns an unsupported refusal; a need for conversion on a cycle; when three accepted sentences are annotated | Whoever would build the slice: a contribution, the kernel or a vocabulary. It counts as a reference. With the migration | [Refusals](#refusals), [Evidence](#evidence), [Optional fields](#optional-fields) |

**O23 to O26.** The same reader read the changes at `0102bfe` and reported
13 findings, N1 to N13. In four of them an accepted page says what this page
or a decision had left open or said differently, and the research owner made
[four decisions](https://github.com/Bitspark/bittype/issues/196#issuecomment-6842416622).
The other findings are taken as the reader proposed.

**O27 to O39.** The prototype of step 2 was then built. Its report has 24
findings, and "Five are sentences of the page that could not be built or
shown as written." The research owner made
[thirteen decisions](https://github.com/Bitspark/bittype/issues/196#issuecomment-6844431866)
on it. They also take, as the prototype fixed them, seven points that this
page had left to it.

**Whose composition registers a facet (O11).** O8 said "registered by the
owner's composition", and no accepted text names an owner of a composition.
O11: "It is the consumer's composition."

- A facet is registered by the composition of the build that makes the
  request, like every contribution.
- A publisher may ship one registration for its bundles. The consumer's
  composition calls it, and the consumer still selects.
- Nothing is registered apart from a consumer's composition. A publisher that
  registered its facets on its own would be a built-in registration, which
  the contract-owned binding decision of 29 October 2026 rules out.
- The kernel's check stays as O8 said: a facet names a registered
  contribution of the same publisher.

For the builtin contributions bittype is the publisher, and bittype's own
test composition is the consumer that registers their bundles.

## What is delegated, and what is not specified here

**Delegated**, each with the one who decides:

| What | To | Row |
| --- | --- | --- |
| Whether a field-level conversion slot is admitted, and its slot kind. Row 11 is decided without that answer (O38). An objection with a reason, before the migration merges, reopens it | The owner of the service library, on bitspark-service#38 | 11, 13, O38 |
| Objections of the service library's owner to refusing a present `undefined` | Weighed by the research owner. The decision itself is made (O14) | 35, O14 |
| The decision of row 11, if no answer has come when the prototype is complete. Made on 13 November 2026 as O38 | The research owner, on the evidence of the prototype, recorded with its reasons | O10, O38 |
| Whether a prototype of one family may precede integers in the service library | The owner of the service library, on bitspark-service#38 | 31a |
| The form of the memo token, the states of a memo entry, and the cost charged at a memo hit | The service library, on bitspark-service#40 | 22, 36 |
| The wording of the memo sentence of the binding contracts | The service library | 18 |
| The shapes of the entries, the serialized form of the plan, and who writes the realization manifest in which form | [#195](https://github.com/Bitspark/bittype/issues/195) | 18, 23 |
| The planner's declaration identity as a field of the inventory record | [#272](https://github.com/Bitspark/bittype/issues/272), delivered with #195 after pull request #281 | 23 |
| The reader on the result of the TypeScript formatter. It lands before the gate can be met on `main` (O39) | [#273](https://github.com/Bitspark/bittype/issues/273) | 24, O39 |
| Recursive references. Also whether the horos analysis has to change first: the horos source refuses a recursive declaration before the kernel meets it (report, F7, not confirmed; [noted on #164](https://github.com/Bitspark/bittype/issues/164#issuecomment-6844432717)) | [#164](https://github.com/Bitspark/bittype/issues/164) | 9a |
| Explicit specialization | [#166](https://github.com/Bitspark/bittype/issues/166) | 21 |
| The report of binding support | [#197](https://github.com/Bitspark/bittype/issues/197) | The delivery note of the table |
| Value-level Go in production, and with it a typed expression vocabulary for Go (O33) | "the future Go runtime library" (`docs/binding-generation.md:46`) | 29, O33 |
| Which facts a nominal declaration restates | The first family that needs evidence through one, [#209](https://github.com/Bitspark/bittype/issues/209) | O35 |
| The roles and paths of the new packages in `LAYOUT.md` and in the architecture test | The migration, whose review judges the prototype's proposal | 28, O39 |
| The entry for the codec role in the list of `AGENTS.md` of who consumes the horos and ontos model dependencies | The migration | 28, O37 |

**The notes posted with the merge.** The decisions of 13 November 2026 list
four notes that are posted when this page merges. It merged on 13 November
2026 as `2a29674`, and all four are posted.

- **On #207, and in its text.** The text carries a dated note on item 1
  (row 2) and the replaced sentence on the mutant (row 39), with the former
  sentence kept. A
  [comment](https://github.com/Bitspark/bittype/issues/207#issuecomment-6842729411)
  corrects the earlier one: vertical case 1 is not exercised with a mock (O6).
- **The grade comment on #196** said the same of vertical case 1. The
  [comment of the merge](https://github.com/Bitspark/bittype/issues/196#issuecomment-6842728735)
  corrects it.
- **The memo sentence.** Row 18 asks for a clarification of the binding
  contracts' memo sentence "by request". The request, with the two senses of
  "instance", is a
  [comment on bitspark-service#40](https://github.com/Bitspark/bitspark-service/issues/40#issuecomment-6842729795).
- **The seven lines of the gate** are in the text of #196, under item 6. Row
  40 makes them the merge gate of that item.

A note that row 39 names is made with this change: the validation cell of row
2 of research 0015 ([what a test must show](#what-a-test-must-show)).

**Not specified here.** The prototype fixed many of the points that this
list held. Each is stated where it belongs, as fixed by the prototype and
provisional until the gate, and
[What the prototype showed](#what-the-prototype-showed) lists them. What is
left, and what the prototype found newly open:

- **The final Go types and signatures of the contract.** O9 decides that a
  rule returns a typed value of a generated plain Go type. The prototype
  wrote one form of the contract
  ([the protocol](#what-a-rule-does-the-protocol)). No decision takes it, and
  the migration writes the final one.
- **How the planner finds that two uses share a node.** Row 18 fixes the
  outcome for two uses of `Page<Person>`. Row 19 fixes what identity is over.
  In the prototype a node's identity is a digest of the JSON encoding of the
  plan, which its report lists as a shortcut (section 11). The procedure
  belongs to #195.
- **The form of a publisher's name, of a local key, of a meaning identity
  and of a meaning revision.** O4, O19 and O24 fix what each means.
- **Who owns a family's meaning, in part.** For an exact mapping the accepted
  pages say it: the representation owner is "bittype's structural profile for
  the structural families, a library for its exact mappings"
  (`docs/representation-families.md:75-76`), the contract or domain library
  owns "its protocol or domain meaning" (`docs/binding-generation.md:39`),
  and domain predicates live in the domain library, "with a stable identity
  and version referenced by every target's binding"
  (`docs/representation-families.md:154-155`). Not said: how the identity of
  a domain predicate relates to the meaning identity, who owns a family's
  meaning when two publishers contribute to it (O19), and who then owns the
  vocabulary package of the family (O27).
- **Which facts a nominal declaration restates** when a rule asks about an
  operand of a named type. Until it is decided the answer is `unknown`. O35
  gives the trigger: "the first family that needs evidence through one, #209
  (maps and sets)".
- **How a declaration contribution states that it is transparent.** O35
  says that it states it with its decision. The prototype did not build it,
  and its form of a fact has no such statement.
- **What the prototype's refusal of a specialization leaves of #166.** In the
  prototype a specialization's domain names a declaration and its actuals,
  and "An actual that is itself an application is matched by its callee
  alone. Enough for the refusal, no design of #166." (report, F9)
- **The reference categories for an external symbol and for a generic
  requirement.** The prototype issued neither
  ([lowering](#what-lowering-must-be-able-to-refer-to)).
- **The key under which a contribution asks for a shared helper.** Row 25
  says that it makes a request. Whose key that request uses is not said. The
  request is refused until a case needs it (O34).
- **The law of a field-state protocol.** The laws of binding generation are
  written "For a type T and representation profile p"
  (`docs/binding-generation.md:125`). An optional field is no type. The
  protocol is bounded by the two laws for the record that contains the field,
  which is a type: `D(E(n)) ≈ n` for every admissible native value and
  `E(D(v)) = v` for every abstract value
  (`docs/binding-generation.md:133-136`). What the protocol's own law is
  within that bound is not said. The prototype did not fix it.
- **The kinds of evidence.** Each belongs to its vocabulary. This page fixes
  the outcomes and the rule for a binder. The accepted table of required
  obligations bounds the kinds (`docs/representation-families.md:274-281`):
  for `Page<T>`, "`T`'s selected adapter travels with its native type"; for a
  map, "Compatible key semantics, or a uniform canonical-map representation";
  for a union, "Discrimination or overlap coherence for the admitted
  actuals"; for a nullable `T`, "The null representation cannot collide with
  a valid `T`". The prototype has one kind in production code: each target's
  model states the fact that a representation has the record form, which the
  declaration contribution asks about (report, sections 11 and 12).
- **How a contribution after the delivered ones gets its codec, and the
  public way from a schema source to the pinned artifact.** O9 decides the
  hand-written form for the delivered contributions. O39 asks that a
  contribution outside the repository make its artifact with public packages
  only. The report (F25): "the route needs a public way from a schema source
  to the pinned artifact. Not examined: #172."
- **What the prototype's proposal of roles lacks.** A role for the vocabulary
  package of a family (O27), which the prototype did not build, and the
  import of the ontos value packages by a codec (O37)
  ([the roles the prototype proposes](#the-roles-the-prototype-proposes)).
- **When the three accepted sentences that a field-level slot changes are
  annotated** ([optional fields](#optional-fields)).
- **Composition across several profiles.** The advice (section 5.6) names it
  as needing an explicit contract, and no row decides it. Accepted sentences
  decide part: "a request selects the libraries that state the profile it
  names", and a bounded profile "reuses no builtin mapping, because each
  builtin mapping states the structural profile"
  (`docs/representation-families.md:361`); composition must "reject
  duplicates, incompatible profiles, missing facets and unsupported
  applications" (`docs/binding-generation.md:113`). "A request selects by
  its profile." (`docs/semantic-generation.md:59`) So contributions of two
  profiles do not meet in one plan. Not said: whether a later contract lets
  them meet. The prototype did not touch it.
- **Which contribution is the materially different one of the gate.** The
  advice (section 4.2) names three candidates. Step 3 is not done. The
  prototype showed line 7 in a small form.
- **How the authoring route meets an application command line**
  ([#172](https://github.com/Bitspark/bittype/issues/172)). Not examined
  here, and not by the prototype.

## Where each row is stated

"Line" is the line of the [gate](#the-gate) whose tests validate the row.
O1 to O39 are the [decisions](#points-decided-after-the-first-version) that
bear on a row. A dash means that the row has no test of its own before the
gate is met.

| Row | Verdict | Stated under | Line |
| --- | --- | --- | --- |
| 1 | Holds | [What a contribution is](#what-a-contribution-is), [Amendments](#amendments-to-accepted-pages), O20, O26 | 7 |
| 2 | Holds | [What a contribution is](#what-a-contribution-is), [Amendments](#amendments-to-accepted-pages), O12, O28 | 1 |
| 3 | Adapted | [The protocol](#what-a-rule-does-the-protocol), O3, O9, O13, O28, O31 | 2 |
| 4 | Holds | [The protocol](#what-a-rule-does-the-protocol) | 1 |
| 5 | Holds | [The protocol](#what-a-rule-does-the-protocol), O1, O20, O29 | 2 |
| 6 | Holds | [The protocol](#what-a-rule-does-the-protocol), [Amendments](#amendments-to-accepted-pages), O20 | 1 |
| 7 | Adapted | [Evidence](#evidence), O13, O23, O29, O35 | 2 |
| 8 | Adapted | [Evidence](#evidence), O2, O20 | 2 |
| 9 | Adapted | [Evidence](#evidence), O17, O30 | 2 |
| 9a | Deferred | [Deferred](#deferred-with-triggers), O6 | - |
| 10 | Adapted | [Refusals](#refusals), O17, O26 | 2 |
| 11 | Adapted | [Optional fields](#optional-fields), [Amendments](#amendments-to-accepted-pages), O10, O38 | The condition beside the seven lines. Met by O38 |
| 12 | Adapted | [Optional fields](#optional-fields), O38 | 2 |
| 13 | Adapted | [Optional fields](#optional-fields), O3, O29, O38 | 4 |
| 14 | Holds | [Optional fields](#optional-fields), [The extension domain](#the-extension-domain), O10, O38 | 2 |
| 15 | Adapted | [Identities](#identities), O4, O16 | 2 |
| 15a | Deferred | [Deferred](#deferred-with-triggers), O16 | - |
| 16 | Holds | [Identities](#identities), O4, O19, O24, O39 | 3 |
| 17 | Adapted | [Identities](#identities), O4, O19, O20, O22, O24, O27 | 3 |
| 17a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 18 | Holds | [Identities](#identities), [Amendments](#amendments-to-accepted-pages) | 2 |
| 19 | Holds | [Identities](#identities) | 2 |
| 20 | Holds | [Identities](#identities), O5, O31 | 2 |
| 21 | Holds | [Identities](#identities), [The extension domain](#the-extension-domain) | 2 |
| 22 | Adapted | [Identities](#identities), O20 | 4 |
| 23 | Adapted | [Identities](#identities), [Amendments](#amendments-to-accepted-pages), O7, O15, O16 | 6 |
| 24 | Adapted | [Identities](#identities), O39 | 6 |
| 24a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 25 | Holds | [Identities](#identities), O34 | 2, 5 |
| 26 | Adapted | [What the kernel carries](#what-the-kernel-carries-without-reading-it), O8, O11, O16, O20, O32, O36 | 2 |
| 26a | Deferred until 13 November 2026. Decided by O36 | [What the kernel carries](#what-the-kernel-carries-without-reading-it), [Deferred](#deferred-with-triggers), O36 | 4 |
| 27 | Adapted | [What the kernel carries](#what-the-kernel-carries-without-reading-it), O5, O19, O20, O23, O32, O37 | 2, 3 |
| 28 | Adapted | [What the kernel carries](#what-the-kernel-carries-without-reading-it), O9, O12, O21, O25, O27, O37, O39 | 1, 3 |
| 29 | Adapted | [What lowering must be able to refer to](#what-lowering-must-be-able-to-refer-to), O8, O33, O39 | 5 |
| 30 | Holds | [What is final and what is provisional](#what-is-final-and-what-is-provisional), [Amendments](#amendments-to-accepted-pages), O17 | 7 |
| 31 | Adapted | [The order of work](#the-order-of-work), O6, [What the prototype showed](#what-the-prototype-showed) | All |
| 31a | Deferred | [Deferred](#deferred-with-triggers), [The order of work](#the-order-of-work) | - |
| 31b | Deferred | [Deferred](#deferred-with-triggers), [The order of work](#the-order-of-work), O6 | - |
| 32 | Holds | [The order of work](#the-order-of-work) | 2, 6 |
| 33 | Holds | [The extension domain](#the-extension-domain), [Amendments](#amendments-to-accepted-pages), O10, O38 | 7 |
| 34 | Adapted | [Inputs](#inputs-that-are-not-bittypes-alone), O29, O39 | 4 |
| 34a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 35 | Holds | [Inputs](#inputs-that-are-not-bittypes-alone), O14, O26 | 4 |
| 36 | Stale | [Inputs](#inputs-that-are-not-bittypes-alone) | - |
| 37 | Adapted | [Inputs](#inputs-that-are-not-bittypes-alone), [The extension domain](#the-extension-domain), O18, O39 | 7 |
| 37a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 38 | Holds | [Payload bootstrapping](#payload-bootstrapping), O9, O39 | 3 |
| 39 | Holds | [What a test must show](#what-a-test-must-show) | 1 |
| 40 | Adapted | [The gate](#the-gate), O10, O18, O20, O26, O29, O33, O38, O39 | All |

18 rows hold, 21 are adapted, eight are deferred, one was deferred until its
trigger fired (row 26a) and one is stale.

## The code today

At `04f9a75`, on `main`. Nothing in this section is the design, and no fact
here is a reason for a choice above. The prototype is another branch and is
not described here. Paths without a directory are in
`contracts/generation/go`.

| Fact | Where |
| --- | --- |
| A registered unit is a profile and three lists: use rules, declaration recipes and helper recipes | `description.go:81-86` |
| A use rule is a pattern with an optional capture, and a fixed output or an expand callback. A declaration recipe has an expand callback | `description.go:21-31`, `:66-71` |
| The expand callback can take the source reader | `description.go:122-129`; `transform/syntax/horos-to-golang/mapping/record/go/record.go:40-41` |
| A rule's callback runs before its children are planned. `Mapped` records a placeholder, and the planner represents the child when it meets the placeholder afterwards | `planning.go:580`, `:592`, `:612-617`; `description.go:142-180` |
| A rule returns native intent, and nothing else of its decision is kept. The program holds declarations, and a step is five strings for inspection | `planning.go:19-28`, `:40-43`, `:596` |
| The declaration recipe takes the record form when the source body is a record. Each target's model requires that the selected body of that form is its record form | `transform/syntax/horos/mapping/declarations/go/declarations.go:52-55`; `transform/syntax/golang/model/go/model.go:105-109` |
| A rule declines by returning an error, which the planner passes on unclassified | `planning.go:407-409`, `:585-587` |
| The record rule's refusal is a text diagnostic. By reading, it matches none of the seven values | `record.go:50` of the Go record mapping; `transform/syntax/horos/source/go/read.go:126-133` |
| A recursive declaration is refused | `planning.go:350-353` |
| An exact mapping without a recipe does not fall back to the structural policy | `planning.go:361-365` |
| A rule's names: a string that is unique in its description, and a counter of the registry. The planner joins them, and only a step keeps the result | `record.go:40` (`record/golang/1`); `transform/syntax/shared/generation/go/registry.go:19-22`, `:126-127`, `:242`; `planning.go:150` |
| A replacement reaches use rules only | `registry.go:251-262` |
| A request names one profile and no product. The registry plans with every unit that states that profile | `registry.go:31-38`, `:185` |
| One native declaration is planned for each source declaration. Its identity is issued at reservation | `planning.go:350`, `:321-329` |
| A use of a declaration is selected by the callee alone | `planning.go:482-485`, `:497` |
| An application keeps its native arguments only. An actual that the native convention omits is admitted and not represented | `planning.go:547`, `:535-539` |
| One generic convention holds native parameters, bindings indexed by source parameter and arguments indexed by native parameter | `description.go:38-45` |
| The optional representation is chosen inside each record rule: Go through a helper key, Rust and C++ by an external type, TypeScript by the flag of the field | `record.go:70-72` (Go), and `record.go:53-54` of the Rust and C++ mappings, `record.go:55` of the TypeScript mapping |
| Go's optional helper is its own contribution. The consumer creates its key, which is an allocation | `transform/syntax/horos-to-golang/mapping/optional/go/optional.go:21-26`; `expression.go:41-46` |
| The `field` form of the Go model carries a name and the optional flag as strings | `transform/syntax/golang/model/go/model.go:43-45` |
| The source reader returns a field's name, its flag, its use and a site | `transform/syntax/horos/source/go/read.go:25-30` |
| A sealed output holds a snapshot, an inventory of the declarations it owns, the exposed sites and the external uses. An inventory record has a handle, a site, a role, a provenance and owner-issued values. It has no identity of the planner | `output.go:164-170`, `:137-143` |
| The contract package imports six packages of bittype and neither horos nor ontos | `integration/go/architecture_test.go:182` |
| There are 21 mapping packages: five each for TypeScript, Rust and C++, six for Go | The directories `transform/syntax/horos-to-<native>/mapping/<name>/go` |
| TypeScript is checked with `exactOptionalPropertyTypes` | `internal/typescriptsem/compiler.cjs:30` |
| Generated Files formatters exist for horos and TypeScript | The directories under `format/files` |

**Merged while this page was written.** Pull request
[#281](https://github.com/Bitspark/bittype/pull/281) merged as `04632ae` on
13 November 2026. It added the owner-issued values, the handles and the
inventory of the row above. The references of this page were first written
at `ec96b92`, moved to `04632ae`, and moved to `04f9a75` with the change that
writes the prototype's results into this page.

**What does not exist.** No Go file on `main` defines a plan entry, a
payload schema, a facet, a payload codec, evidence about a selected
representation, or a completeness test backend. The prototype branch has a
form of each, and it never merges. The service library has no
package `structural-adapters`. Its owner decided that it is created with
integers
([decision of 11 November 2026](https://github.com/Bitspark/bitspark-service/issues/38#issuecomment-6803983510)).

## The order of work

**The steps (row 31, as the grade lists them).**

| Step | What | Where it lands |
| --- | --- | --- |
| 1 | This page | `main` |
| 2 | A prototype | "on a branch that is thrown away" (row 31). **Done on 13 November 2026**: `prototype/contribution-contract-196` at `6080420`. It does not merge ([what the prototype showed](#what-the-prototype-showed)) |
| 3 | A materially different contribution, with the shared packages fixed | Not done. The prototype showed line 7 of the gate in a small form. It is shown again on `main` with step 4 |
| 4 | The contract in its final form and the migration of the 21 mapping packages, in one change | `main`, when the [gate](#the-gate) is met. It waits for #273 (O39) |
| 5 | The optional move | `main`. Row 11 is decided as D3 (O38), so the move happens unless an objection reopens the row. The production record packages move here and not in step 4 (O18) |
| 6 | The additive work of [#207](https://github.com/Bitspark/bittype/issues/207) | `main` |

**Why not the advised order.** The advice (section 4.2) asks for TypeScript
end to end through the service library's own factory before the migration.
The table's introduction says why that cannot be done: "The package it wants
the proof to run through, `@bitspark/structural-adapters`, does not exist. Its
owner decided on bitspark-service#38 that it is created with integers, which
the plan puts after the migration." There "the plan" is the plan of work of
[#206](https://github.com/Bitspark/bittype/issues/206), not the plan of the
word table. The same facts bound the risk: "no other repository builds on the
contract before integers."

**What remains before step 4.** Derived here, from the report and O39.

- Step 3.
- #273, for line 6 of the gate and the test of row 24. O39: "#273 lands before
  the gate can be met on `main` (F11)."
- The test of row 10 over all 21 packages. The prototype moved ten, and the
  migration moves the other eleven.
- A public route from a schema source to the bootstrap artifact, for the
  contribution of line 7 that is authored out of tree (O39).
- What four decisions ask and the prototype did not build: the vocabulary
  package of a family (O27), the forwarding of an alias's facts (O35), a
  codec that imports the ontos value packages (O37), and a snapshot that
  gives the children of any descriptor by position (O39).

**What the prototype was to contain (rows 31 and 32).** It contains each,
with the limits under [What the prototype showed](#what-the-prototype-showed).

- The rule in its staged form, one round ([the protocol](#what-a-rule-does-the-protocol)).
- The frozen decision record, the plan and the serialized form of the plan.
- Two consumers of the evidence path. Row 3: "the declaration recipe's need to
  know that a selected body is a record, and by a fixture family".
- Go handlers and executed Go lowering, in the test tree.
- A TypeScript reader as a stand-in, in the test tree, with the snapshot of
  [row 34](#inputs-that-are-not-bittypes-alone).
- The optional field in the delegated form ([row 11](#optional-fields)).
- Fixtures for vertical cases 3 to 6, and vertical case 2 at the kernel's
  level, through the independent fixture source, with a mock family. Vertical
  case 1 is not in it (O6).
- A transport descriptor for the first facet that crosses a process boundary,
  and a report of it on this page (row 26a).
- What the prototype fixed of the points that this page had left to it
  ([what the prototype fixed](#what-the-prototype-fixed)).
- The four cases of row 32: "repeated uses share a node, equal native types
  with different adapters stay apart, an erased native argument is kept for
  conversion, one plan under two TypeScript placements".

**The six vertical cases in the prototype.** The cases are those of the plan
(`docs/representation-families.md:111-116`). Row 31 puts fixtures for cases 3
to 6 into the prototype. Row 31b, as O6 corrects it, keeps case 1 out and
case 2 at the kernel's level. O6 gives the reason: "Recursion is refused and
deferred (row 9a), so a mock cannot exercise it at the level of the kernel."
The table says what the prototype built for each case (report, stage 8).

| Case | Its form in the prototype | What that form does not show |
| --- | --- | --- |
| 1. Recursive `ContentValue` | Not built. O6: "Case 1 waits for #164 as a whole." | Anything of recursion beyond its two refusals |
| 2. `Page<Outcome>` with optional nullable fields | At the kernel's level, through the independent fixture source, with a mock family, the contribution `mock-optional-nullable`, which the report names as a mock (O6). The kernel plans a generic applied to a record whose every field is delegated, through a source that is not horos | Null against absence, discriminators, anything about the family the mock stands for, and conversion |
| 3. The exact external generic `Cell<T>` | An exact generic that omits its parameter natively. Its argument is linked for conversion, and its value facet is found | An additional native parameter, capabilities, and what `Cell` does with its slot |
| 4. Go typed construction with a generic adapter argument | `Page[Person]` itself, lowered from the serialized plan and executed | A typed expression vocabulary (O33), and lowering of families that do not exist |
| 5. A generic map whose key requirement is rejected | A fixture generic with a binder requirement. It is discharged by a proved fact, refused as unsatisfied by a refuted one and as unresolved where nothing is stated | The map family and the meaning of its keys |
| 6. An independent family without a shared change | Line 7 of the gate in its small form (stage 9) | Authorship by another owner, TypeScript, and the contribution of step 3 |

The six cases in real form stay what representation families makes them: the
condition for a stable version of the envelope
(`docs/representation-families.md:109-110`). The prototype does not discharge
that condition.

**What a stand-in proves (row 31).** Row 31: "the serialized plan is
sufficient for a second reader in another process, not for the service
library's factory." The table details that sentence. It is derived here,
from verification report 3. The prototype's stand-in is such a reader, and
its report says the same of it (stage 7).

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

## What the prototype showed

Step 2 of the order is done. This section says what the prototype built, what
its tests show, and what they do not show. Its source is the
[prototype's report](../research-docs/0017-contribution-contract.prototype.md),
which the agent that built the prototype wrote. The report is evidence. The
[decisions O27 to O39](https://github.com/Bitspark/bittype/issues/196#issuecomment-6844431866)
say what this page takes from it. By AC-07 this page claims no more than the
tests establish: "do not use one kind of evidence to claim another
guarantee".

**The branch.** `prototype/contribution-contract-196` at `6080420`, with the
branch point `2a29674`. It opens no pull request, never merges, and is thrown
away after the migration. Its log is `PROTOTYPE.md` at its root.

**The result in numbers (report, sections 1 and 2).**

- All nine stages are built. The decisions say of it: "The design could be
  built."
- 41 tests pass. The gate runs were made on the commit `4442047`, one commit
  before the head, which changes the log and one comment.
- A test logs green for the kernel as built and red for a kernel with one
  property switched off. The switch is an instrument of the branch. 25 such
  mutants each turn a test red.
- A test without a mutant is red in one sense only: it does not build against
  the branch point.
- The whole test suite passes on the head, except seven Rust tests that need
  a helper of the Rust compiler. They fail at the branch point too.
- 24 findings, F1 to F27. None stopped a stage. "Five concern a sentence of
  the page that could not be built or shown as written (F5, F8, F11, F15,
  F16)."

**How the old contract stayed (report, section 1).** The mapping contract of
`main` stands beside the new one on the branch. Each moved mapping package
keeps its old library and gains a contribution, so that today's tests run
unchanged and one test compares the two contracts on one module. Moved are
ten of the 21 packages: for Go text, boolean, list, record, declarations and
optional, and for TypeScript text, list, record and declarations. A new
TypeScript package `optional` is added. Rust, C++ and the TypeScript boolean
package stay on the mapping contract.

### What each stage built, shows and does not show

The table is the report's section 3, shortened. The third column is as
important as the second.

| Stage | Its tests show | They do not show |
| --- | --- | --- |
| 1. The staged rule and the frozen record | The moved contributions generate the same Go text as the mapping contract for today's pagination module, and the same TypeScript text for that module without booleans. The projector's capability mentions no use, declaration, binder, selector or reader, and three attempts to read the source do not type-check. A value mutated after the rule returned changes neither the native type nor the serialized plan. With a fixture whose member names come from its configuration, an entry rebuilt from the source is detected. For eleven vocabulary packages the committed module, pin and Go type are what the schema source gives, and every listed value round-trips through its hand-written codec | Rust, C++ and TypeScript booleans. Equality of the origins inside the native tree. That a contribution cannot hand source data to its projector through state of its own package. Immutability against a codec that keeps the value. That no second pass can be written: the mutant is a pass of the test, not of the kernel. TypeScript payload types, which are not generated. A field the existing families cannot type. The bootstrap test was not shown red |
| 2. Evidence | A rule gets a typed answer about what was selected for its operand before it decides. With a declared binder requirement, `Page<text>` and the written-out record have the same native type, the open declaration records the requirement and the application its discharge. Each of the five outcomes reaches a rule as its own answer. The declaration contribution decides nominal or alias from the answer about its body and reads no payload of the record. Two representation revisions give equal payloads and different nodes. A recursive reference and an evidence cycle are refused with two reasons. An answer for another dependency, generic environment, representation revision or meaning revision is rejected | That a fact is true. That the kernel refuses a rule that decides differently inside a generic declaration (F5, O29). A check that reaches a limit: exhausted is stated, nothing computes. Anything of recursion beyond its refusal, and the two refusals run on the independent fixture source only (F7). That a rule cannot act on what it remembers without presenting it |
| 3. Identities and refusals | One key under two publishers registers, and one publisher with one key twice fails, also for a second target. For five families Go and TypeScript carry one meaning and two representation revisions, and a reference that lacks either part is refused. Missing and ambiguous are two values. An unbuilt family is a missing contribution. A specialization is refused as unsupported with the application named while another application plans. Five refusals of moved rules match exactly one value each. Two contributions that spell a helper alike get two symbols | That a publisher's name is its owner's. That the meaning is one identity and not two equal spellings (F2, O27). The refusals of the eleven packages that are not moved. A request for a shared helper (F10, O34) |
| 4. The plan and its serialized form | Two uses of `Page<Person>` are two use entries with two sites on one node. Asked for bindings, an application of a declaration that omits its parameter natively has no native and one adapter argument. Asked for types alone, the omitted argument is not planned. A reader that only parsed the serialized plan derives equal keys for equal nodes and another key with another selected dependency. An unknown envelope protocol, payload schema, meaning revision, representation revision and source semantics are each refused under their own name. The inventory record carries the planner's identity. One plan links with the manifests of two TypeScript placements, and the manifest of another output is a linking error | An adapter. The memo. The fourth axis, which stage 7 shows in a reader. A reader on the formatter's result (F11) |
| 5. Facets, Go handlers and executed Go lowering | A fresh process lowers Go conversion for `Page[Person]` from the serialized plan, the snapshot and the facets. The code compiles, converts five values there and back and refuses three outside the type. It has a typed construction, member access, the owned helper and a generic adapter argument. A lowering by member reference compiles and runs, and a lowering by string does not compile. Guarantees 1, 2, 3 and 6 | A typed expression vocabulary: the lowering is Go text with typed holes (F14, O33). The conversion laws |
| 6. Optional fields by delegation | In lowered Go and in the stand-in, absent and present with a value convert there and back as themselves, and present with `null` is refused. In the stand-in a present `undefined` is refused, and that vector is red when the protocol reads it as absence. A request with the operand of another record, or with a role issued for another pair, fails. A fragment that is no member is refused, a target without a member vocabulary refuses a field position, and a fragment of two members is admitted | A `null` that a type admits. The service library's slot. Independent authorship |
| 7. A second reader in another process | The TypeScript stand-in, run by Node, converts `Page<Person>`, a nested page and a list there and back from the plan and the snapshot alone. It refuses five axes under their own names, the runtime compatibility among them. Without the snapshot both readers fail explicitly, and a snapshot of another pin is a linking error in both. It derives the same memo keys as the Go reader | The service library's factory, its budgets, its call context and its memo. Independent authorship: one seat wrote both sides |
| 8. The fixtures | Vertical cases 2, 3 and 5, in the forms that [the order of work](#the-order-of-work) lists. Case 4 is the lowering of stage 5, and case 6 is stage 9 | Case 1, which is not built. For case 2: anything about the family that the mock stands for |
| 9. The gate, and line 7 in its small form | After a replacement the payloads of the 9 other selections are byte-identical, 5 frozen records that do not depend on the replaced contribution are identical, and 11 that depend on it are other nodes. With both registered the domain is ambiguous until the consumer replaces explicitly. `Page<Person>` regenerates with the new contribution and converts through its field-state protocol, with the backend unchanged. A Go module outside the repository authors and registers a contribution and a codec with public packages only, and a file that needs an internal package does not build. A contribution at a field position cannot request a sibling's operand | Another owner, or TypeScript. A generation with the module outside the repository, and how its author gets the schema artifact (F25). The commit pair of the migration on `main`. One file outside the contribution changed (F24) |

### The gate on the prototype branch

The lines are those of [the gate](#the-gate). The table is the report's
section 4.

| Line | State on the branch | Not run, or its limit |
| --- | --- | --- |
| 1. Native projection | Green. Red under two mutants and against a pass that derives the entry from the source | None reported |
| 2. Dependency composition | Green | The test of row 10 over all 21 packages is not run in full. Ten are moved, and the eleven others return text diagnostics on the mapping contract |
| 3. Serialization | Green | The bootstrap test was not shown red. For the axes the test plans each case and edits the serialized bytes only for the envelope and the source semantics, because a node's identity covers its schema reference and its contribution reference (F12) |
| 4. A second reader in another process | Green. Red under a protocol that reads a present `undefined` as absence | Not run: present with `null` where the type admits it (#211) |
| 5. Executed Go lowering | Green, and the lowering by string does not compile | The instrument is Go text with typed holes (F14, O33) |
| 6. Realization for TypeScript | The test of the realization manifest is green, with a stand-in for the reader | The test of row 24 is not run: `main` has no reader on the formatter's result (#273, F11) |
| 7. The shared packages fixed | Green in the small form | One file outside the contribution changed (F24). Not shown: the commit pair of the migration on `main` |

- **The two checks that stay green** are shown for Go and TypeScript, by the
  first test of stage 1. Rust and C++ are not moved.
- **The condition beside the lines.** The prototype decided nothing about row
  11. Its section 8 is the evidence, and O38 is the decision
  ([optional fields](#optional-fields)).

**The six guarantees (report, section 5).**

| # | The guarantee | Red under the mutants |
| --- | --- | --- |
| 1 | Required facets belong to the selected contribution | No mutant. Without one facet a request for bindings is "missing facet" and a request for types succeeds |
| 2 | Kind, target and compatibility match | `facet-ignores-issuer`, `facet-ignores-target`, `facet-ignores-version` |
| 3 | A reference has the right plan, category and scope | `reference-ignores-plan`, `reference-ignores-category`, `reference-ignores-scope` |
| 4 | Evidence belongs to its dependency | `evidence-ignores-dependency`, `evidence-ignores-environment`, `evidence-ignores-representation-revision`, `evidence-ignores-meaning-revision` |
| 5 | The payload snapshot is immutable | `projector-gets-returned-value`, `freeze-at-serialization` |
| 6 | The exact payload passed validation | `codec-skips-validation`. A value that the Go type admits and the schema does not is refused before the projector runs. A receipt for another frozen value or another schema fails |

### What is not shown

Gathered from the report, so that no line above is read as more than it is.

- **Eleven of the 21 packages are not moved:** every package of Rust and of
  C++, and the TypeScript boolean package. Nothing is shown for them.
- **Two tests of the gate were not run:** the test of row 24, and the test of
  row 10 over all 21 packages.
- **The stand-in is no factory of the service library,** and one author wrote
  both sides. It is evidence about the plan and the protocol, and not about
  that library's adapter interface.
- **The Go lowering is no typed expression vocabulary** (O33).
- **The bootstrap test was not shown red,** and no TypeScript payload type is
  generated.
- **Every facet of the prototype is test code.** No facet package exists in
  production.
- **Not every test has a mutant.** The report says "No mutant" of three: the
  test that compares the native text of the two contracts, the test of the
  replaced contribution, and the test of guarantee 1. A test without a mutant
  is red only in that it does not build against the branch point.
- **The origins inside the native tree were not compared.** The printed text
  is equal.
- **Recursion** beyond its two refusals. The horos source refuses a recursive
  declaration before the kernel meets it, so those refusals are shown on the
  fixture source only (F7).
- **A present `null` where a type admits it** (#211).
- **Not assessed by the prototype** (report, section 13): whether the budgets
  of `main` fit the new contract, concurrency, composition across several
  profiles, and how the authoring route meets a command line (#172).

### Decisions that are later than the prototype

Each of these asks for something that the branch does not have. The page
states the decision, and the prototype shows nothing of it.

| Decision | What it asks | What the branch has |
| --- | --- | --- |
| O27 | One vocabulary package of the family, without a target, for the meaning identity | The identity spelled in each target's vocabulary package, and a test that compares the spellings (F2) |
| O35 | The kernel forwards the facts of an alias's body | A fixture policy that asks about its body and states the answer again. bittype's own declaration contribution states nothing (F23) |
| O37 | A codec may import the ontos value packages | Codec packages that import neither horos nor ontos, and a shared part that exports the value forms a codec builds (F26) |
| O39 | The snapshot gives the children of any descriptor by position | Two readers with a case for each of the descriptor kinds record and list (F27) |
| O39 | A contribution out of tree makes its bootstrap artifact with public packages only | A test that writes the canonical module and the pin with an internal package and the type by hand (F25) |

### What the prototype fixed

Each of these points was left to the prototype. Each is stated where it
belongs, and all of them are provisional until the gate. "Taken" means that
the decisions take the point as the prototype fixed it, by an O-number or in
their closing paragraph.

| Point | Stated under | Taken |
| --- | --- | --- |
| The form of an evidence kind and of a fact | [Evidence](#evidence) | Yes |
| How a rule presents an answer | [Evidence](#evidence) | Yes |
| How a binder's requirement is declared, recorded, found and discharged | [Evidence](#evidence) | Yes |
| The class of the refusal of a member fragment | [Refusals](#refusals) | Yes |
| How a codec is registered and found | [What the kernel carries](#what-the-kernel-carries-without-reading-it) | Yes |
| Which payloads are empty | [What the kernel carries](#what-the-kernel-carries-without-reading-it) | Yes |
| That Go's optional package returns no decision before presence is delegated | [What the kernel carries](#what-the-kernel-carries-without-reading-it) | Yes |
| What the shared part of the codec role does and what a codec does | [What the kernel carries](#what-the-kernel-carries-without-reading-it) | O37, which changes what a codec may import |
| The facet list of a contribution | [What the kernel carries](#what-the-kernel-carries-without-reading-it) | O32 |
| The transport descriptor of a facet | [What the kernel carries](#what-the-kernel-carries-without-reading-it) | O36 |
| The roles and paths of the new packages | [Below](#the-roles-the-prototype-proposes) | O39, as a proposal |
| The shape of a contribution, a rule and a projector | [The protocol](#what-a-rule-does-the-protocol) | No decision names it |
| The form of a refusal, the two refusals of row 9, and an unexpected error | [Refusals](#refusals) | No decision names it. By O17 the set of values is closed once the prototype has fixed them |
| That a request names a product by its identity | [Identities](#identities) | No decision names it |
| The reference categories that are issued, and who issues a member reference | [What lowering must be able to refer to](#what-lowering-must-be-able-to-refer-to) | No decision names it |
| That a member fragment may hold more than one member | [Optional fields](#optional-fields) | No decision names it |

### The roles the prototype proposes

The prototype added roles and edges to the architecture test. O39: they "are
the proposal for the edit of `LAYOUT.md`", and "The migration's review judges
them." This page lists them and edits no role table. The table is the
report's section 10.

| Role | Path | Imports | The report's reason |
| --- | --- | --- | --- |
| Vocabulary of a contribution | `transform/syntax/<source>-to-<native>/vocabulary/<name>/go` | The generation contract and the publisher name. Never a mapping, a codec, a facet package or another vocabulary | O25 |
| Contribution: payload codec | `transform/syntax/<source>-to-<native>/codec/<name>/go` | The generation contract, the shared part of the codec role and its own vocabulary. Never a mapping, a facet package or another codec | Row 28, O18 |
| Contribution: facet package | `completeness/<source>-to-<native>/facets/<name>/go` | The generation contract, the value vocabulary of its target's product and its own contribution's vocabulary. Never a mapping, a codec or another facet package | O8, O16, O25 |
| Vocabulary of a completeness product | `completeness/<native>/values/go` | The generation contract and the syntax of its target | Row 29: the facet kinds belong to the product |
| Contribution: completeness backend | `completeness/<native>/backend/go` | The generation contract and its product's vocabulary | Row 31 |
| Foundation and kernel: the shared part of the codec role | `transform/syntax/shared/payload/go` | The generation contract, and of the external models horos core, horos declarations, ontos core, data and codec | None given |
| Foundation and kernel: the publisher name | `transform/syntax/shared/publisher/go` | Nothing | None given |

- **A changed role.** A mapping may import its own vocabulary, and is
  forbidden every codec, every facet package and every other vocabulary.
- **The test tree** is the top-level directory `completeness/`. "No other
  role permits an import of it, so only a composition and its own packages
  import it." (report, section 10)
- **Not classified by the test:** the JavaScript files of the stand-in.
- **Where two decisions go beyond the proposal.** O27 asks for a vocabulary
  package of a family, which each target's vocabulary imports. The proposal
  has no such role, and it forbids a vocabulary to import another. O37 lets a
  codec import the ontos value packages, which the proposal does not allow.
  The migration's edit has to carry both.

**What else the report has for the migration.** Its section 11 tells apart
the files that belong to the contract from the shortcuts of the branch. The
shortcuts are not for the migration: the mapping contract left beside the new
one, the mutant switch, the JSON encoding of the plan with a node's identity
as a digest of it, the snapshot that a test writes with an internal package,
and everything under `completeness/` in the form it has. Its section 12 says
what each package gets: a vocabulary package, a codec package, the rule in
two steps with its projector, and its refusals as values with a class.

## The gate

Row 40: "Adopted as the merge gate of the migration, each line in the form
bittype can show: a frozen decision survives native projection, dependency
composition, serialization, a second reader in another process, executed Go
lowering, and realization for TypeScript, with every shared package named and
fixed while a materially different contribution is added".

The migration of step 4 does not merge while one line is red. The gate is
shown first on the prototype branch and then on `main` (row 31). On the
prototype branch it is shown with limits, and two of its tests were not run
([what the prototype showed](#what-the-prototype-showed)). Each line
below is a predicate. Its tests come from the Validation cells of the table,
marked by their row, and from the verification reports, marked V1, V2 or V3.
The tests marked O20 are that decision's, and it calls them derived here.
O26 corrects the test of row 1 in line 7: "A selection that depends on the
replaced contribution has another frozen record, because an answer is bound
to the dependency's contribution reference."
Every test in the table is red on `main` today, for the reason given.

| Line | The predicate | Tests | Red today because |
| --- | --- | --- | --- |
| 1. Native projection | The native type is computed from the frozen record, by a projector that cannot read the source | A projector that tries to read the source does not compile or fails the contract test (row 4). Mutating the value a rule returned changes neither the native type nor the serialized plan (row 28). A fixture whose decision the source does not determine is green on the planner and red on a mutant that derives the entry from the source (rows 6 and 39, V1) | The one callback takes the reader. No record exists |
| 2. Dependency composition | A rule decides with typed answers about what was selected, and every reference is checked | A fixture family whose decision depends on evidence about its operand plans and generates (row 3). `Page<text>` and the written-out record generate the same native type, and a rule that decides differently inside a generic declaration than for its written-out application fails its family's coherence test (rows 7 and 8, O29). Two contributions that differ only in their representation revision give equal payloads and different node identities (row 5, O20). An evidence cycle and a recursive reference are refused with different reasons (row 9, V1, O30). Every refusal of the 21 mapping packages matches exactly one value (row 10). A request at a field position with the operand of another record fails (row 12, V2). One local key under two publishers registers, and one publisher with one key twice fails (row 15). Two uses of `Page<Person>` are two use entries on one node (row 18). Equal native types with different adapters stay apart (rows 19 and 32). An omitted native parameter still yields an adapter argument (row 20). A specialization is refused as unsupported with the application named (row 21). Two contributions that spell a helper alike get two symbols (row 25, V2). A member fragment that the target's member vocabulary does not admit is refused as unsupported (row 14, V2). Guarantees 1 to 6 (rows 26 and 27) | A rule runs before its children. No evidence, no owner of an identity and no use entry exist |
| 3. Serialization | The plan, read from its serialized form in a fresh process, is the plan that was written | The completeness test runs from the serialized plan in a process that never held the planner's objects (V3). Each version axis is exercised alone (row 16). For the runtime compatibility that is a test of the facet's reader, because the kernel carries that facet unread (O39). The contributions of one family for two targets carry one meaning identity and two representation revisions, and a reference that lacks either part is refused (row 17, O20). The payload types regenerate from their schema and every value round-trips through its codec (row 38, O9). A payload outside its schema is refused before a backend runs (row 38, V3) | No plan is serialized. No plan carries the axes |
| 4. A second reader in another process | A reader in another language converts from the plan and the snapshot alone | The TypeScript stand-in round-trips `Page<Person>` and refuses a payload schema version it does not know (V3). Given the plan without the snapshot it fails explicitly, and a snapshot with another pin is a linking error (row 34, V3). Each state of an optional field has its own vector, in Go and in the stand-in, and a present `undefined` is refused (rows 13 and 35, O14). These vectors stay under either decision of row 11 (O26). The stand-in derives a memo key from the serialized plan alone: two uses with equal nodes, applications and selected dependencies get equal keys, and a different selected dependency gives a different key (row 22, O20) | No reader and no conversion exist |
| 5. Executed Go lowering | Go code that is lowered from the plan compiles and runs | Construction of `Page[Person]`, access to a member, an owned helper and a generic adapter argument compile and run (advice, section 3.4). A lowering that takes a member by a string the naming policy changed fails to compile, and with a member reference it compiles (row 29, V3). The line shows the reference categories, their checks and executed code. It does not show a typed expression vocabulary for Go (O33) | A member is a string. No lowering exists |
| 6. Realization for TypeScript | A reader takes every native name and place from the realization | A plan read with the realization manifest of another output fails as a linking error, not as an incomplete plan (row 23). One plan under two TypeScript placements resolves each logical symbol to another module path and the same declaration (row 32, V2). No native name reaches emitted glue except through the reader, and a formatter mutant that realizes another spelling is followed (row 24, V2). That test needs the reader of #273, which lands before the gate can be met on `main` (O39) | The inventory record has no planned identity, and the formatter's result has no reader |
| 7. The shared packages fixed | A materially different contribution is added, and no shared package changes | One registered contribution is replaced by another for the same domain. The payloads of every other selection are byte-identical, and so are the frozen records of every selection that does not depend on the replaced one (row 1, O20, O26). The contribution registers, and `Page<Person>` regenerates with it. The diff of that test touches no shared package (row 40, V3). A Go module outside the repository authors and registers a contribution with public packages only, and makes its bootstrap artifact with them (row 37, V3, O39). The composition is one for each process that plans or reads (O39). One fixture for each item that is in the [extension domain](#the-extension-domain), and one recorded refusal for each that is out (row 33, V3) | A replacement reaches use rules only. No request names a product |

**Two checks that are green today and stay green.** The fixtures of today's
generics pass on the new contract with no contribution for generics (row 2,
V1). The native declarations of today's pagination fixtures do not change in
any of the four targets (#207, item 5).

**One condition beside the seven lines: row 11 is decided (O10). It is met
(O38).** "The migration does not wait without limit: its gate requires row 11
to be decided. It is decided by the answer on bitspark-service#38, or, if
none has come when the prototype is complete, by me on the evidence of the
prototype, recorded here with its reasons." No answer came, and O38 decides
row 11 as D3 on that evidence ([optional fields](#optional-fields)). The test
of row 11 in the table names the answer: "The answer on bitspark-service#38
is recorded here with the form that was built. Red until then". What is
recorded is the decision, with the form that was built. An objection with a
reason, before the migration merges, makes the condition open again.

**The gate on `main` (O18).** Four things that the lines above leave to
guess.

- **The tests of the delegated form run against fixtures.** "The gate's tests
  of the delegated form (the request at a field position, the member
  fragment, the four states) run on `main` against fixture contributions in
  the test tree. The production record packages move in step 5. If row 11 is
  decided as D1, those tests leave the gate with the form." O26 corrects that
  sentence in one point: the four states do not leave. "Row 35 and O14 do not
  depend on row 11. Under D1 the record contribution itself refuses a present
  `undefined`, and the vectors of the four states run against it." Row 11 is
  decided as D3 (O38), so the tests of the delegated form stay in the gate.
- **The codec on the list of shared packages** is the shared part of the
  role. "A family's codec is its own package, and a new contribution adds
  one." O37 says what the shared part does, and the prototype had to change
  it for one family's value form before O37 (report, F26).
- **The contribution that is tested out of tree states native types only.**
  "That is intended for now: every binding facet kind and the Go value
  vocabulary live in the test tree, and a production home for them belongs to
  the first backend's owner (rows 31a and 37a)."
- **Line 7 is shown twice on `main`.** "Out of tree for the native half,
  where no file of bittype changes by construction. In the test tree for the
  binding facets, as its own commit of the migration's pull request after the
  migration commit, with its diff recorded in the pull request before the
  squash." So the one change of step 4 has at least two commits before it is
  squashed, and the diff that touches no shared package is that of the second.
- **The composition is one for each process (O39).** O39 says that the
  composition of line 7 "is one for each process that plans or reads (F24)."
  A reader registers the facet and the codec of every selected contribution.
  In the prototype one file outside the contribution of line 7 changed for
  that: a table lets a contribution's file add both (report, F24).

**The shared packages of line 7.** Row 40 asks that they be "named and
fixed". Derived here, from the list of generation packages in `LAYOUT.md`:

- the contract package, `contracts/generation/go`;
- the registry, `transform/syntax/shared/generation/go`;
- the shared native construction, `transform/syntax/shared/target/go`;
- the source reader and the common mapping parts of the source,
  `transform/syntax/horos/source/go` and `transform/syntax/horos/mapping`;
- the native model and target of each language;
- every existing mapping package;
- the shared part of the payload codec role. A codec of one schema is not on
  this list (O18);
- the completeness test backend.

The prototype proposes paths for the last two and adds a package for the
publisher name
([the roles the prototype proposes](#the-roles-the-prototype-proposes)).

**What the gate depends on.** Derived here. Line 6 needs the reader of #273
and the inventory of pull request #281 with the planner's identity. O39: "#273
lands before the gate can be met on `main` (F11)." Line 4 has its last state,
a present `null`, only with #211. Line 7 needs a public route to the
bootstrap artifact (O39). The contribution of line 7 depended on the decision
of row 11, which is made (O38). The Go facets of line 5 are facet packages of
bittype's own contributions in the test tree (O8, O16).

**What the gate does not establish.** It does not establish the service
library's factory (row 31a), realization for Go and Rust (row 24a), recursion
(row 9a), vertical case 1 as a whole, or vertical case 2 in real form (row
31b, O6). Each is recorded as not demonstrated, with its trigger, under
[Deferred](#deferred-with-triggers).

**After the gate.** Derived here, from verification report 3: the
contribution of line 7, pinned at the commit of the migration, still builds
at the head of each later family. It is red when a later family changes a
type of the contract (row 30).

**When the gate is met** this page's first line changes: the design is final,
with the date and the commit of the migration. Until then each change that
the prototype makes to the design is recorded on this page with its date.
