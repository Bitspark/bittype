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
more. O1 to O22 decide them
([Points decided after the first version](#points-decided-after-the-first-version)).
File and line references are `path:line` at commit `04632ae`, the head of
`main` when this page was last changed. The line numbers of the pages that
this change annotates are those of that commit.

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
| The contribution's facet list | What a contribution states about its facets. It is per contribution (O16) |
| The realization manifest | What relates a plan, an output and its realization (row 23). It is not the publisher's manifest, which row 15a defers |
| The test tree | Packages that only a test composition imports. Where they live is [not specified here](#what-is-delegated-and-what-is-not-specified-here) (O16) |

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
| A role | "Optional record field" | Advice, section 1.1, and conditional on [row 11](#optional-fields) |

Derived here: the kernel selects by the domain and never by registration
order. Row 10 and the advice (section 1.5) say it of a refusal, which must not
make the kernel try another contribution. Two contributions that claim one
domain under one profile need an explicit composition by the consumer (row
15, from the advice, section 2.1).

**One target and one profile.** A contribution states both. The target is "a
stated field of the contribution, not a part of the key" (O4,
[identities](#identities)). The profile is
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

O12 names the second line: the declaration contribution "is the part of the
division of generics that a contribution owns". The first line agrees with the
planner page: "The kernel owns capture identity, substitution and scope."
(`docs/generation-planner.md:435`). It changes two
texts that list generics beside the families. For item 1 of
[#207](https://github.com/Bitspark/bittype/issues/207) the
[comment of 13 November 2026](https://github.com/Bitspark/bittype/issues/207#issuecomment-6840545162)
on that issue states the new reading. The delivery row of representation families that names
"generics" among the delivered families
(`docs/representation-families.md:341`) gets a dated clarification.

## What a rule does: the protocol

**One round (row 3).** Row 3: "a rule names its dependencies and the evidence
it needs, the planner resolves them, the rule decides with typed answers, the
planner freezes the record, the projector runs."

| Step | Who | What happens |
| --- | --- | --- |
| 1. Name | The rule | From the checked subject and the source view, the rule names its requests. A request names one of three things: a child use that the kernel minted, a [field position](#optional-fields) if presence is delegated, or a [shared helper](#identities) (row 25). For each the rule says what it needs of the dependency: a reference or evidence ([row 9](#evidence)). O3: "A request may state what the rule already knows from the source without evidence." |
| 2. Resolve | The planner | It selects a contribution for each request, plans it, and reads the evidence asked for from what was selected (O13, [evidence](#evidence)) |
| 3. Decide | The rule | With the source view and the typed answers, the rule returns its decision or a [refusal](#refusals). The decision is a typed value (O9, [below](#what-the-kernel-carries-without-reading-it)) |
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
  (`docs/generation-spi.md:371-374`).

**What belongs in the decision record (row 5).** Row 5: "A per-selection
choice is a field. An invariant of a semantic revision is carried by reference
to that revision".

- **A per-selection choice** is anything that can differ between two
  selections of the same contribution and that affects representation or
  conversion. The advice (section 1.2) lists member names, the presence
  strategy, the selected construction strategy, a layout that depends on an
  operand, and applicable configuration.
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
backend (`docs/binding-generation.md:72`). The page gets a dated clarification
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
round. None of the six delivered contributions (O12) and none of the six
vertical cases of the plan (`docs/representation-families.md:95-102`) is such a case by its
description. Adding the round then is a change of the shared contract and
needs the justification of AC-08.

**An illustration.** The advice sketches the protocol and says of its sketch:
"This is a protocol sketch, not a proposed Go API." The same holds for the
lines below. They show the order of one round and decide no signature.

```text
name(checked subject, source view)
    -> requests, each for a reference or for evidence
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
| Where does an answer come from? | From what was selected for the dependency. No provider computes it (O13, below) | O13 |
| What does the kernel check? | That an answer belongs to the dependency it was asked about: the selected representation and its whole contribution reference, the target, profile and configuration, the generic environment, the kind and version of the requirement, and the scope. O19: an answer of evidence binds "to the whole contribution reference, both of its parts" | Advice, section 1.3; O19; [guarantee 4](#what-the-kernel-carries-without-reading-it) |
| What does the kernel not check? | That the fact is true. "It cannot establish a domain assertion merely because a contribution supplied it." | Advice, section 1.3; AC-05 |

Facts about the source stay with the source. Representation families divide
the judgment: "horos semantics and bittype's source integration establish
canonical facts", and "target families establish native discrimination and
compatibility" (`docs/representation-families.md:191-194`). Evidence in the
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

| Question | Answer (O13) |
| --- | --- |
| What does the planner read an answer from? | The frozen decision of the dependency's selected contribution, or what that contribution's revision fixes |
| Who computes it? | Nobody at that moment. The dependency stated the fact when it decided |
| What is an answer bound to? | The asking rule, the dependency and the dependency's contribution reference |
| What if the selection states no such fact? | The answer is `unknown` |
| Is an answer part of the frozen record? | Yes, as a link, and so part of the [identity of the node](#identities) |

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
anything" (`docs/representation-families.md:199-201`). The first form of row
8 called its own four names the accepted ones. O2: "Row 8 was imprecise." And:
"Evidence in general has proved, refuted, unknown and exhausted. For
disjointness the refutation keeps the accepted name `overlap`." So `overlap`
is the refuted outcome of that one judgment, and no accepted name changes.

The fifth outcome is the accepted course of a generic family: "A generic
family proves its requirements for every permitted actual, declares
requirements that actuals must satisfy, or defers qualification to checked
applications." (`docs/representation-families.md:252-254`).

**A binder's requirement: what the accepted pages fix (O13).** O13 asks this
page to quote them and to list what they leave open.

| Question | What an accepted page says |
| --- | --- |
| Why is a requirement needed? | "A native type parameter does not carry the adapter, equality, discriminator or layout evidence needed to bind values of that type" (`docs/representation-families.md:240-241`) |
| Who declares a binder's requirement? | The generic family: it "declares requirements that actuals must satisfy" (`docs/representation-families.md:252-253`). In the division of [row 2](#what-a-contribution-is) that is the selected representation contribution |
| Where is it carried? | With each generic use: "Each generic use carries its scoped binder identities, native type arguments, adapter arguments and requirements" (`docs/representation-families.md:92-94`) |
| Who discharges it at an application? | A checked application does: the family "defers qualification to checked applications", and "Qualifying `U<int>` does not qualify `U<text>`." (`docs/representation-families.md:253-254`). In the division of row 2 the planner makes the "Instantiation and compatibility checks" (advice, section 1.1) |
| Under which assumption is an open declaration planned? | "Open eligibility must be established under binder assumptions over the profile's explicitly represented substitutions" (`docs/type-generation-plugins.md:503-505`) |

What those sentences leave open is
[not specified here](#what-is-delegated-and-what-is-not-specified-here): in
which frozen record a requirement is stated and in what form, how the planner
finds the requirements of a binder when it plans an application, and with
which reason an application that cannot discharge one is refused. The
prototype fixes them within these bounds and reports them on this page (O13).

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
  [Deferred](#deferred-with-triggers). The prototype shows nothing of
  recursion beyond that refusal: vertical case 1 is not in it (O6,
  [the order of work](#the-order-of-work)).
- **Two refusals that can be told apart. Derived here.** The distinction of
  row 9 can be observed only if the refusal of an evidence cycle and the
  refusal of a recursive reference carry different reasons
  ([refusals](#refusals)).

## Refusals

**What exists (the generation contract at `04632ae`).** Listed here because
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
| Unsupported | The selected contribution recognizes the case and cannot represent or check it yet. A later slice can admit the same input. O17: it "is for a case that a selected contribution recognizes and does not cover, as row 21 uses it" | `ErrUnsupported`, narrowed. Today that value also covers a missing mapping, a missing helper and recursion |
| Refused by design | The case lies outside the domain that the representation promises. Another contribution may cover it | New. It is not `ErrRejected`, which is the source's refutation of a use |
| Requirement unsatisfied | A selected dependency fails a declared requirement: the evidence is refuted | New |
| Requirement unresolved | The evidence a decision requires is unknown or exhausted | The same meaning as `semantic.ErrEvidenceUnavailable` has for source evidence |
| Missing facet | The selected contribution has no facet for a requested product | New. It stays a class of its own |
| Ambiguous or missing contribution | More than one contribution, or none, answers a domain. A family that is not built is the second case (O17) | `ErrAmbiguous` for the first. The second is under `ErrUnsupported` today |
| Limit | A resource bound of the request is reached | `ErrLimit` |

The first four meanings follow the advice (section 1.5). Which outcome of
[evidence](#evidence) leads to the third and which to the fourth is derived
here. "Missing facet" stays
recognizable because representation families requires it: "Unsupported use,
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

- **The class is the kernel's, and the set closes with the prototype.**
  Failure meaning is one of the things composition must preserve (AC-05). A
  closed set of classes in the kernel is therefore no catalogue of
  contributions. O17: "The set of refusal values is closed once the prototype
  has fixed the members that the page lists as undecided, and the page says
  so." The undecided members are at the end of this section.
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
  that policy." (`docs/generation-spi.md:127-128`)
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

**Not decided by row 10**, and listed under
[Not specified here](#what-is-delegated-and-what-is-not-specified-here): the
form of a refusal in Go, the value of the two refusals of
[row 9](#evidence), whether a missing contribution and an ambiguous one are
one value or two, what an unexpected error of a rule is, and the class of the
refusal of a member fragment that a target's member vocabulary does not
admit. These are the undecided members that O17 means.

## Optional fields

> **This whole section is conditional.** It states bittype's position (row
> 11). The position needs a slot in the service library's adapter interface
> that its owner has not agreed to. The question is
> [with that owner](https://github.com/Bitspark/bitspark-service/issues/38#issuecomment-6840511750)
> on bitspark-service#38 since 13 November 2026. Nothing in this section is
> fixed until row 11 is decided and the decision is recorded here and in the
> table. O10 says who decides: "It is decided by the answer on
> bitspark-service#38, or, if none has come when the prototype is complete, by
> me on the evidence of the prototype, recorded here with its reasons." "Me"
> is the research owner.

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
| The slot is refused | Row 11: "presence becomes the record's own layout (D1), finding 7 is withdrawn, and the Go optional package folds into the record contribution". The round of requests stays, and the request at a field position is not built (O10, below) |
| No answer yet | The prototype builds D3. Row 11: "put to that library's owner on bitspark-service#38 and built in the prototype" |
| No answer when the prototype is complete | The research owner decides row 11 on the evidence of the prototype and records the reasons (O10) |

**What is in the contract under either answer (O10).** "The round of requests
is part of the contract under either answer, because evidence needs it. The
subject at a field position and the member-fragment answer are built only if
presence is delegated." So a refusal removes nothing from the
[protocol](#what-a-rule-does-the-protocol). It leaves the subject of row 12
and the answer of row 13 unbuilt, and the role leaves the
[kinds of domain](#what-a-contribution-is).

**Why the fallback is no design decided by an existing interface.** Row 11
lets the answer of another repository decide between D1 and D3. AC-09 keeps
an existing interface from deciding a design. The decisions of 13 November
2026 say why this is not that case: "the accepted texts support both
readings, and the slot belongs to an interface that another repository owns,
so the choice is made with that owner. If no answer comes, O10 decides it on
the prototype's evidence, not on what the interface offers today."

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
  mapping of that pair" (`docs/naming.md:154`).
- **How the name travels (O3).** "A request may state what the rule already
  knows from the source without evidence." The record contribution knows the
  name from the source view, so it states the name in its request. The
  resolved requests are part of the frozen record, and O3 concludes: "So a
  name stated in a request is recorded once, there, and the decision does not
  restate it."
- **The conversion half.** Only a record can omit a field, in every target.
  So no target converts an optional field from its retrieved value alone. For
  TypeScript representation families already says it: "the adapter separates an
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
  express a contribution with no standalone value adapter". An optional
  contribution has a field-state protocol and nothing that converts a value on
  its own. The advice (section 1.6): "Contextual field conversion is an
  admitted product shape that the contract needs to express now." This
  sentence of row 14 holds under either decision of row 11. O10 makes
  conditional only the subject at a field position, the member-fragment
  answer and one line of row 33.

**What this position changes outside this page, if it stands.** None of the
following is edited now. Each waits for the decision of row 11.

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
  bytes." (`docs/binding-generation.md:101`).
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
  family's payload as arguments." (`docs/representation-families.md:49-50`).
  The contribution identity says what was selected. It is what duplicates and
  version checks compare.

**The five axes stay apart (row 16).** Representation families: "**Five version
axes:** the envelope protocol; the family payload schema; the
representation's semantic identity; generated-code/runtime compatibility;
source semantics." (`docs/representation-families.md:88-90`). Row 16: "The
contribution reference carries the semantic revision. The other four axes are
stated where each belongs and are not folded into the identity".

| Axis | What it identifies (advice, section 2.1) | Where it is stated. Derived here, except the third |
| --- | --- | --- |
| Envelope protocol | How the shared graph is represented | The plan (#195) |
| Payload schema | How a family's recorded decisions are typed and encoded | The family's payload schema ([bootstrapping](#payload-bootstrapping)) |
| Representation semantics | What the selected representation means | The contribution reference (row 16), in two parts: the revision of the family's meaning and the contribution's representation revision (O4, O19) |
| Generated-code and runtime compatibility | Which generated code and runtime facilities can cooperate | The facet that names a factory, among its requirements (`docs/binding-generation.md:96`) |
| Source semantics | Which source interpretation the checked references use | The source-semantics input ([row 34](#inputs-that-are-not-bittypes-alone)) |

A digest of an implementation is provenance. It replaces none of the five.
The test of row 16: "Each axis is exercised alone: an unknown required value
of one axis is refused while the other four are unchanged." For the third
axis the test runs for each of its two parts. O19: the test of row 16 binds
"to the whole contribution reference, both of its parts".

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

**Which half of "Value meaning" (O22).** O22 corrects O19: "The first two
speak of canonical values and are the same in every target: they are the
family's meaning. The native equivalence and a permitted normalization speak
of native values, which differ by target: the contribution's representation
revision fixes them, with the other six parts." (O22) The first two are the
promised canonical domain and the domain predicates.

| | The family's meaning | The representation revision |
| --- | --- | --- |
| Which part of the family's contract of seven parts | Of "Value meaning": the promised canonical domain and the domain predicates (O22) | Of "Value meaning": the native equivalence and a permitted normalization. And the other six parts (O22) |
| Whose | Not specified where two publishers contribute to one family. O19: "Who owns a family's meaning when two publishers contribute to it is not specified." | The contribution's publisher |
| The same in every target | Yes | No. It belongs to one contribution, so to one target |
| What it fixes | What values the family covers and what they mean | The invariants of one target's representation ([row 5](#what-a-rule-does-the-protocol)) |
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
| "The memo is keyed by the selected representation instance plus input identity, including the descriptor and generic environment." (`docs/representation-families.md:262-263`) | The applied representation, which several uses share. It is not the per-use record |

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
  mappings." (`docs/generation-spi.md:163`). So a request for types does not
  fail for want of something that only conversion needs, as binding
  generation requires: "A native mapping without an adapter facet still
  supports types" (`docs/binding-generation.md:99`).

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
- **A mismatch is a linking error.** The test of row 23: "A plan read with
  the manifest of another output fails as a linking error, not as an
  incomplete plan."
- **A backend has a fourth input (O15).** Row 23: "the placed result is a
  fourth input of a backend." This amends an accepted sentence in substance:
  binding generation says that a backend works from three inputs alone
  (`docs/binding-generation.md:72`). The naming design had left it open:
  "that page stands as written until
  [#195](https://github.com/Bitspark/bittype/issues/195) decides"
  (`docs/naming.md:729-730`). O15: "Row 23 is that decision: I hold #195, and
  the row is linked to it." It is recorded on
  [#195](https://github.com/Bitspark/bittype/issues/195#issuecomment-6841872859),
  the insertion in binding generation is headed "Amended", and the naming
  design gets a dated note
  ([Amendments](#amendments-to-accepted-pages)). Two more sentences list
  native symbols and imports among what the plan carries
  (`docs/binding-generation.md:56`, `docs/representation-families.md:78`).
  O7 gives both the same clarification.

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
plan or in the realization manifest. The serialized form of each is
[provisional](#what-is-final-and-what-is-provisional).

| Identity | Who issues it | Leaves the process | Not to be confused with | Row |
| --- | --- | --- | --- | --- |
| Contribution identity: publisher and local key | The publisher states it | Yes | A source selector, a module pin, a profile, a registry's handle, a factory reference | 15 |
| Contribution reference: the identity, the family's meaning and the representation revision | The publisher | Yes | The other four version axes, a digest of an implementation | 15, 16, O4, O19 |
| The family's meaning: a meaning identity and a meaning revision | Not specified where two publishers contribute to one family (O19). The same in every target | Yes, inside the reference | The contract of seven parts that a family contributes per target (`docs/representation-families.md:32`), the name of one target's contribution, a representation revision | 17, O4, O19 |
| Representation revision | The publisher of the contribution, for that contribution | Yes, inside the reference | The revision of the family's meaning | 17, O4, O19 |
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
| Realized symbol: name and place | The language's Files stage | Yes, through the realization manifest | The naming intent, the spelled name of generation | 23, 24 |
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

**Selected together: a contribution, its facet packages, the bundle (O8, O11,
O16).** "A native facet and the facets a binding product needs are selected
together." (`docs/binding-generation.md:99`). O8 said what is selected:

> "A contribution is a bundle. Its manifest lists its facets, and the
> implementation of a facet may live in another package of the same owner. The
> kernel admits a facet that names a registered contribution of the same
> publisher, registered by the owner's composition." (O8)

Two later decisions correct its words.

- **O11: whose composition.** "It is the consumer's composition." Binding
  generation gives the consumer build the "Selection and composition of the
  concrete contributions it uses" (`docs/binding-generation.md:41`), and "A
  consumer supplies its contributions explicitly."
  (`docs/binding-generation.md:111`).
- **O16: a bundle is not a package.** bittype's rule stands: "One
  contribution, meaning one thing a consumer registers or selects on its own,
  is one package." (`AGENTS.md:54-55`). O16: "A consumer's composition
  registers a mapping package, and it registers a facet package apart from
  it; each is one package. They are selected together. The page calls the
  mapping unit the contribution, a separately registered facet unit a facet
  package of that contribution, and the two together the bundle."

| Word | What it is | Per |
| --- | --- | --- |
| The contribution | The mapping unit: one package that a consumer's composition registers | Itself |
| A facet package | A separately registered package that holds a facet of one contribution | Contribution |
| The bundle | A contribution and its facet packages, which are selected together | Contribution |
| The contribution's facet list | What O8 called "Its manifest": the facets that the contribution has | Contribution |
| The publisher's manifest | What a composition lock would bind a publisher's name to. Deferred (row 15a) | Publisher |
| The realization manifest | What relates plan, output and realization ([row 23](#identities)) | Output |

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
- **Where the test tree is.** O16 asks for the role of `LAYOUT.md` that holds
  test-only packages. There is none. The layout has the test compositions:
  `integration/go` holds "independent consumers and cross-package
  qualification" (`LAYOUT.md:71`), and of the role Composition it says "The
  two test compositions are not classified." (`LAYOUT.md:165`). A package
  that only such a composition imports has no role there. Where the facet
  packages of the builtin contributions live is therefore
  [not specified here](#what-is-delegated-and-what-is-not-specified-here),
  and the prototype reports its choice (O16).

**What the kernel guarantees (row 27).** Row 27: "Guarantees 1 to 6 are the
kernel's, each with its test." The advice (section 3.2) says what kind they
are: "These are structural, authority, scope, and association guarantees."

| # | The kernel guarantees that | Its test | Red today because |
| --- | --- | --- | --- |
| 1 | Required facets belong to the selected contribution | Remove a facet that a bindings product requires. The request for bindings fails as "missing facet", and a request for types succeeds | A request names no product |
| 2 | A facet's kind, target and compatibility match | A facet from another issuer under the same display name, for another target, or of another contract version is rejected | No kind and no version exist |
| 3 | A reference has the right plan, category and scope | A reference of another plan, a member where a type is expected, and a binder outside its scope are each rejected | A reference has no category |
| 4 | Evidence belongs to the dependency it was asked about | An answer issued for another dependency, another generic environment or another contribution reference is rejected. The reference counts whole, both of its parts (O19) | No evidence path exists |
| 5 | The payload snapshot is immutable | Mutate the value a rule returned, after it returned. Neither the native output nor the stored plan changes. O20: the guarantee "has one moment since O9, after the rule returned" | No record exists |
| 6 | The exact payload passed validation | A fixture family whose schema has a constraint that a Go type cannot express, and a rule that returns a value violating it: the value is refused before any backend runs (O20, derived here). A receipt issued for another snapshot fails the same way. A receipt is the kernel's record that one frozen value passed validation against one schema | Nothing validates a payload |

The tests are the advice's, except where O20 is named. Guarantee 1 rests on
O5: a request that names its products "is what guarantee 1 needs to tell a
missing facet from one that nobody asked for". It checks the contribution's
facet list (O16). Two preconditions are derived here, from verification
report 1: guarantee 3 needs the reference categories of the
[next section](#what-lowering-must-be-able-to-refer-to), and guarantee 6 needs
the kernel to know that the codec that validated is the one of the payload's
schema, which is
[not specified here](#what-is-delegated-and-what-is-not-specified-here).

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
(`docs/representation-families.md:84-86`). Row 28 adopts the advice's order
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
- **A mapping package needs no new import.** The generated type is plain Go,
  so a mapping that returns it imports neither horos nor ontos.
- **The test of row 28:** "Mutating the value a rule returned, after it
  returned, changes neither the native type nor the serialized plan."

**The payload codec is a new package role (row 28, O9).** Row 28: "Encoding
and validation live in a new package role outside the contract package, which
may import neither horos nor ontos. The horos source reader is not its home.
`LAYOUT.md` names the home".

- **The role has two parts.** A shared part that encodes, validates and
  freezes, and one codec package for each payload schema. O18: "A family's
  codec is its own package, and a new contribution adds one." A codec is the
  only place where a decision meets ontos and horos conformance.
- **For the six delivered contributions it is written by hand.** O9: "the
  codec is written by hand and pinned with its schema, as row 38 decides for
  the bootstrap" ([payload bootstrapping](#payload-bootstrapping)). Which six
  is decided by O12, below.
- **The path is a postponed part of row 28.** The row says that `LAYOUT.md`
  names the home. This change names the role there and no path. The path and
  the import allowances come with the prototype, and both are
  [not specified here](#what-is-delegated-and-what-is-not-specified-here).
- **Why it is outside the contract package.** "Generic carriers, contracts,
  engine, printer and native syntax stages have no external Horos/Ontos
  imports." (`CHARTER.md:129-131`). The kernel carries an opaque handle to the
  frozen value and the result of its validation, bound to that value and to
  its schema. It imports neither library.
- **Why it is not the horos source reader.** Derived here, from verification
  report 1: a payload's schema is a horos type whatever the source language
  of the generation is. A source that is not horos plans through the same
  kernel and needs the same codecs.

**The six delivered contributions that return a decision (O12).** O9 spoke
of "the six delivered families", and row 2 says that generics are no family
with a contribution. O12 says what the six are:

> "A codec belongs to a payload schema, and a payload schema belongs to each
> delivered contribution that returns a decision. There are six: text,
> booleans, lists, closed records, optional fields, and the declaration
> contribution, which is the part of the division of generics that a
> contribution owns (native parameter order, alias or defined declaration). A
> payload may be empty; the record still exists." (O12)

The packages are those of `transform/syntax/horos-to-<native>/mapping/<name>/go`
at `04632ae`.

| Package `<name>` | TypeScript | Go | Rust | C++ | Returns a decision | Whose schema and codec |
| --- | --- | --- | --- | --- | --- | --- |
| `text` | Yes | Yes | Yes | Yes | Yes | Text |
| `boolean` | Yes | Yes | Yes | Yes | Yes | Booleans |
| `list` | Yes | Yes | Yes | Yes | Yes | Lists |
| `record` | Yes | Yes | Yes | Yes | Yes | Closed records |
| `declarations` | Yes | Yes | Yes | Yes | Yes: the native parameter order, and alias or defined declaration | The declaration contribution |
| `optional` | No package | Yes | No package | No package | Yes, if presence is delegated | Optional fields |

- **A payload may be empty.** O12 says so, and that the record still exists.
  Which payloads are empty is for the prototype to find.
- **Six, or five.** O12: "If presence is not delegated (row 11 decided as
  D1), the optional schema folds into the record's and there are five."
- **Under D3 every target gets an `optional` package.** Derived here, in
  [Optional fields](#optional-fields). Today only Go has one.
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
- **The set of categories** above is the advice's list read by row 29. The
  prototype fixes the set, and how a member reference is issued is
  [not specified here](#what-is-delegated-and-what-is-not-specified-here).

**Where the prototype's value vocabulary lives (row 29).** "The prototype's
Go value vocabulary and its facets stay in the test tree". The reason is the
accepted scope: "bittype's native vocabularies stay types-only."
(`docs/binding-generation.md:83`). O8 keeps both: the facets that the
completeness test backend reads are facet packages of bittype's contributions
and live in the test tree, and a production mapping imports nothing of them.
The prototype's vocabulary is a test instrument. The advice (section 3.4)
says what it must do: "It must compile and run generated construction, member
access, an owned helper, and a generic adapter argument." The production
owner of value-level Go is not bittype: "A Go backend belongs to the future
Go runtime library." (`docs/binding-generation.md:46`).

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
| A further family | In. Derived here | Every family of the delivery table of representation families (`docs/representation-families.md:339-349`) enters as contributions. The accepted promise: "Adding a family through the published identity, dependency, invocation, plan-graph and target-extension contracts changes only that family and the consumers that compose it." (`:22-24`) | A use of a family that is not built is refused as "missing contribution": there is no contribution to select (O17) |
| Further products and facets | In | A further backend for a target, and a further kind of facet, enter through the [facet contract](#what-the-kernel-carries-without-reading-it) with no edit of the kernel | A request for a product whose facet a contribution lacks fails as "missing facet" |
| A contextual representation at a field position | In, if presence is delegated | The request of [row 12](#optional-fields) and the member-fragment answer of [row 13](#optional-fields). O10: "Row 33's line on a contextual representation at a field position is conditional on that", which is the decision of row 11. If presence stays the record's layout, the item is out. That the contract can express a contribution without a value adapter of its own ([row 14](#optional-fields)) holds either way | A fragment that the target's member vocabulary does not admit |
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
| Present with `undefined` | Refused, by its own vector, in TypeScript (O14) |

- **A present `undefined` is refused (O14).** Two accepted sentences settle
  it. Who decides: "The representation owner owns its factories' semantic
  contract and qualification: bittype's structural profile for the structural
  families" (`docs/representation-families.md:60-62`). And what the adapter
  does: "TypeScript uses `field?: T | null`, and the adapter separates an
  absent own property from a present `undefined` at run time."
  (`docs/representation-families.md:167-168`). So reading a present
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
| A field-level conversion slot for an optional record field | 11, 13 | Question 1 |
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
  (`docs/representation-families.md:86-87`). The advice repeats it: "Do not
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
  (`docs/representation-families.md:82-83`).
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
text of #207 itself is not reworded yet
([the notes not yet posted](#what-is-delegated-and-what-is-not-specified-here)).

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
| [Binding generation](binding-generation.md#order-of-work-and-qualification): "The unit of implementation is library × target." (`:265`) | The sentence names ownership and distribution. The unit that a consumer selects is one contribution, and a library may publish several | Read as the unit of selection, the sentence would make every rule of a library one indivisible choice | 1 | Research 0014, held by `claude-bittype-gaps` |
| [Binding generation](binding-generation.md#the-selected-representation-plan): "The backend can implement each selected representation from the plan, the source semantics and the selected contributions alone, without guessing a rule." (`:72`), and "consumes three inputs" (`:23`). **Amended in substance** | For a generated native declaration the three inputs are not sufficient. Its name and place exist only after a Files stage has placed it. A plan entry refers to it by an identity the planner issues, and a backend reads the realization from the placed result, which is a fourth input | The [naming design](naming.md#reconciliation-with-accepted-documents) stated the tension and left it to #195, which decided it with row 23 | 23, O15 | Research 0014, held by `claude-bittype-gaps` |
| [Naming](naming.md#reconciliation-with-accepted-documents): "that page stands as written until [#195](https://github.com/Bitspark/bittype/issues/195) decides" (`:729-730`) | A dated note: #195 decided on 13 November 2026, and binding generation carries the amendment | The sentence waited for that decision | 23, O15 | Research 0016. This chain holds it |
| [Binding generation](binding-generation.md#the-selected-representation-plan): the plan carries "the rules for absence, presence, null and native `undefined`" (`:58`) | The plan carries a rule either as a recorded field, where it varies between selections, or through the contribution reference whose revision fixes it | Row 5 keeps an invariant out of the payload, and both sentences hold | 5, O1 | Research 0014, held by `claude-bittype-gaps` |
| The same list: "native symbols, distinguishing run-time imports from type-only imports" (`:56`) | The plan carries the logical symbol identity and the spelled intent of a generated declaration. Its place and the imports come from the realization | The same reason as for the fourth input | 23, O7 | Research 0014, held by `claude-bittype-gaps` |
| [Representation families](representation-families.md#the-plan): "Native symbols, imports, factory references, runtime requirements" (`:78`) | The same clarification as the row above | The same | 23, O7 | Research 0015. This chain holds its delivery |
| [Representation families](representation-families.md#the-representation-contract): "so a native declaration and its conversion cannot disagree." (`:46-47`) | "Cannot disagree" means that no sanctioned second decision point exists. It does not mean that a projector or a factory cannot be wrong | Read as a claim about every implementation, the sentence contradicts AC-05 | 6 | Research 0015. This chain holds its delivery |
| The same paragraph: "one immutable representation instance per use" (`:44-45`), and "The memo is keyed by the selected representation instance" (`:262`) | "Instance" has two senses. In the first sentence it is the selection at one use, a use entry. In the second it is the applied representation that several uses share | Read with one sense, every use would have its own memo key | 18 | Research 0015 |
| [Representation families](representation-families.md#the-plan): "The envelope keeps an explicitly unstable version" (`:95-96`) | What is unstable is the version: the layout and the encoding. The meanings that the plan shares with the contribution contract are fixed with that contract | The partition "contract final, envelope unstable" is withdrawn | 30 | Research 0015 |
| [Representation families](representation-families.md#delivery): "After them a family adds a contribution and changes no existing contract." (`:328`), and "generics" among the delivered families (`:341`) | "Changes once" is a commitment for the written extension domain. "Generics" names a division of responsibility and no family with a contribution of its own | Rows 33 and 2 | 33, 2 | Research 0015 |
| [Package boundaries](package-boundaries.md#findings-on-main-at-3bf45c3), finding 7 (`:156`) | The finding is bittype's derivation and is conditional on the answer of the service library. If the slot is refused, the finding is withdrawn | Row 11 | 11 | bittype, case bittype/ACX-009 |

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
[bittype/ACX-009](architecture-cases/0009-package-boundaries.md) say that
the optional move is conditional, `LAYOUT.md` names the new package role
without a path, and the [documentation map](README.md) lists this page.

**Sentences that are not amended now.** Each needs its owner or an answer
first. They are listed where they arise:

| Sentence | Why not now | Where |
| --- | --- | --- |
| The three sentences that a field-level slot would change | They wait for the decision of row 11 | [Optional fields](#optional-fields) |
| "A single generic convention relates source binders, native parameters and complete mapped actuals." (`docs/generation.md:45-47`) | The page describes type generation as it is implemented. It changes with the migration, and this page reads the sentence as the native convention | [Identities](#identities) |
| The memo sentence of the binding contracts | It is the service library's | [What is delegated](#what-is-delegated-and-what-is-not-specified-here) |

## Deferred, with triggers

Nine rows are deferred. A trigger is a state of a repository that anyone can
check.

| Row | What is deferred | The trigger that ends the deferral | Tracked |
| --- | --- | --- | --- |
| 9a | Recursive representation references | "#164 is claimed and vertical case 1 is being built." Until then recursion is refused: `grep -n "recursive native declaration" contracts/generation/go/planning.go` finds the refusal | [#164](https://github.com/Bitspark/bittype/issues/164) |
| 15a | A composition lock that binds a publisher's name to the publisher's manifest. A contribution's facet list is not deferred (O16) | "a second publisher registers a contribution (bitspark-service#39 item 2 starts)" | bitspark-service#39 |
| 17a | A root facet of `Service` that recognizes a family's meaning | "bitspark-service#39 starts its root facet. No root facet exists in either repository" | bitspark-service#39 |
| 24a | A reader of the realization for Go and Rust | "slice 1 of #161 or #159 delivers a generated-Files formatter (`ls format/files/golang format/files/rust`)" | [#161](https://github.com/Bitspark/bittype/issues/161), [#159](https://github.com/Bitspark/bittype/issues/159) |
| 26a | A transport descriptor for a facet | "the first facet crosses a process boundary in the prototype of row 31" | #196 |
| 31a | TypeScript end to end through the service library's own factory | "`packages/structural-adapters` exists on bitspark-service `main`, or its owner records on bitspark-service#38 that a prototype precedes integers" | bitspark-service#38 |
| 31b | Vertical case 1 as a whole, and vertical case 2 in real form (O6) | "#164, #210 and #211 deliver their families" | [#164](https://github.com/Bitspark/bittype/issues/164), [#210](https://github.com/Bitspark/bittype/issues/210), [#211](https://github.com/Bitspark/bittype/issues/211) |
| 34a | A reader protocol for source semantics | "the design of bitspark-service#40 item 3 starts" | bitspark-service#40 |
| 37a | A contribution that is authored in the service library's repository | "bitspark-service#39 item 2 starts" | bitspark-service#39 |

Three remarks, derived here.

- **Row 26a ends inside the prototype.** The stand-in reader of TypeScript
  reads a serialized plan in another process. The facet that names a factory
  crosses the boundary there, so the transport descriptor is designed in the
  prototype and reported on this page. It is on the list of
  [what the prototype contains](#the-order-of-work).
- **Row 31a is why a stand-in stands in.** What it proves and what it does
  not is under [The order of work](#the-order-of-work).
- **A deferral is no narrowing.** Each deferred part is in the
  [extension domain](#the-extension-domain) or belongs to another owner. None
  is removed from what the accepted pages promise.

## Points decided after the first version

Nothing is open inside this design. This section lists what was open and how
it was decided.

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
| O12 | What "the six delivered families" are, when generics are no family | Six delivered contributions that return a decision, the declaration contribution among them. Five if presence is not delegated | [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| O13 | Where an answer of evidence comes from | From what the dependency's selected contribution stated or its revision fixes. No provider computes it | [Evidence](#evidence) |
| O14 | Who decides what a present `undefined` does | bittype's structural profile. It is refused | [Inputs](#inputs-that-are-not-bittypes-alone) |
| O15 | The page said it amends no accepted page in substance | Row 23 amends one sentence in substance, and it is the decision of #195 that the naming design waited for | [Identities](#identities), [Amendments](#amendments-to-accepted-pages) |
| O16 | A bundle across packages against the rule that one contribution is one package, and three things called a manifest | A bundle is a contribution and its facet packages. Three names for the three lists | [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| O17 | The class of an unbuilt family, the closed set of classes, and recursive references final without a test | "Missing contribution". The set closes with the prototype. Recursive references become final with #164 | [Refusals](#refusals), [What is final](#what-is-final-and-what-is-provisional) |
| O18 | What the gate leaves to guess on `main` | Fixtures for the delegated form, the shared part of the codec role, a native-only contribution out of tree, line 7 shown twice | [The gate](#the-gate) |
| O19 | Which revision evidence binds to, and the name "family contract" | The whole contribution reference. The name is "the family's meaning" | [Identities](#identities), [Evidence](#evidence) |
| O20 | Rows without a test, and two tests that could not fail | Tests for rows 1, 5, 17 and 22, and new tests for guarantees 5 and 6 | [The gate](#the-gate), [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| O21 | Whether the packages of one name share one payload schema | O12 is corrected. A schema belongs to a family's contract for one target; the family's owner may publish one module for several | [What the kernel carries](#what-the-kernel-carries-without-reading-it) |
| O22 | Whether the native equivalence is part of the family's meaning | O19 is corrected. The canonical domain and its predicates are the family's meaning; the native equivalence and a normalization belong to the representation revision | [Identities](#identities) |

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
| Whether a field-level conversion slot is admitted, and its slot kind | The owner of the service library, on bitspark-service#38 | 11, 13 |
| Objections of the service library's owner to refusing a present `undefined` | Weighed by the research owner. The decision itself is made (O14) | 35, O14 |
| The decision of row 11, if no answer has come when the prototype is complete | The research owner, on the evidence of the prototype, recorded with its reasons | O10 |
| Whether a prototype of one family may precede integers in the service library | The owner of the service library, on bitspark-service#38 | 31a |
| The form of the memo token, the states of a memo entry, and the cost charged at a memo hit | The service library, on bitspark-service#40 | 22, 36 |
| The wording of the memo sentence of the binding contracts | The service library | 18 |
| The shapes of the entries, the serialized form of the plan, and who writes the realization manifest in which form | [#195](https://github.com/Bitspark/bittype/issues/195) | 18, 23 |
| The planner's declaration identity as a field of the inventory record | [#272](https://github.com/Bitspark/bittype/issues/272), delivered with #195 after pull request #281 | 23 |
| The reader on the result of the TypeScript formatter | [#273](https://github.com/Bitspark/bittype/issues/273) | 24 |
| Recursive references | [#164](https://github.com/Bitspark/bittype/issues/164) | 9a |
| Explicit specialization | [#166](https://github.com/Bitspark/bittype/issues/166) | 21 |
| The report of binding support | [#197](https://github.com/Bitspark/bittype/issues/197) | The delivery note of the table |
| Value-level Go in production | "the future Go runtime library" (`docs/binding-generation.md:46`) | 29 |

**The notes not yet posted.** The decisions of 13 November 2026 list four
notes that are posted when this page merges. None is posted yet.

- **On #207, and in its text.** Item 1 and the sentence on the mutant are
  still the old ones (rows 2 and 39), and the comment of 13 November 2026 on
  #207 still says that vertical case 1 is exercised with a mock, which O6
  corrects.
- **The grade comment on #196** says the same of vertical case 1 and is
  corrected by O6.
- **The memo sentence.** Row 18 asks for a clarification of the binding
  contracts' memo sentence "by request". The comment of 13 November 2026 on
  bitspark-service#38 tells that owner about the memo token (row 22). Neither
  it nor a comment on bitspark-service#40 names the two senses of "instance".
  The note goes to bitspark-service#40.
- **The seven lines of the gate** in the text of #196. Row 40 makes them the
  merge gate of item 6 of that issue.

A note that row 39 names is made with this change: the validation cell of row
2 of research 0015 ([what a test must show](#what-a-test-must-show)).

**Not specified here:**

- **Any Go type or signature of the contract.** O9 decides that a rule
  returns a typed value of a generated plain Go type. It decides no
  signature, and the sketches on this page decide none. The prototype writes
  them.
- **The form of a refusal in Go:** values, or one type with a class. The
  value of the refusal of an evidence cycle and of a recursive reference.
  Whether a missing contribution and an ambiguous one are one value or two.
  What an unexpected error of a rule is: the advice calls it an
  implementation failure, and row 10 lists no such class.
- **How a generation request names one binding product among several.** O5
  tells types alone from types and bindings. The [extension domain](#the-extension-domain)
  admits further products.
- **How the planner finds that two uses share a node.** Row 18 fixes the
  outcome for two uses of `Page<Person>`. Row 19 fixes what identity is over.
  The procedure belongs to #195.
- **The form of a publisher's name, of a local key, of a meaning identity
  and of a meaning revision.** O4 and O19 fix what each means. Who owns a
  family's meaning when two publishers contribute to it is not specified
  (O19). What the family's meaning of an exact mapping is, is not said.
- **A binder's requirement.** In which frozen record a requirement is stated
  and in what form, how the planner finds the requirements of a binder when
  it plans an application, and with which reason an application that cannot
  discharge one is refused. The prototype fixes them within the bounds that
  [Evidence](#evidence) quotes, and reports them here (O13).
- **The path and the import allowances of the payload codec packages.**
  `LAYOUT.md` gets both with the prototype. Also not said: how a codec is
  registered, how the kernel knows that the codec that validated is the one
  of the payload's schema, and how a contribution after the six delivered
  ones gets its codec. O9 decides the hand-written form for those six.
- **Where the test tree is.** `LAYOUT.md` has no role for a package that only
  a test composition imports. The prototype reports where it puts the facet
  packages of the builtin contributions (O16).
- **The form of a contribution's facet list** (O8, O16). The publisher's
  manifest and the lock that binds a publisher's name to it are deferred (row
  15a).
- **The class of the refusal of a member fragment that a target's member
  vocabulary does not admit.** It is one of the undecided members of the set
  of refusals (O17).
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
  fragment of another shape, exists only if presence is delegated.
- **How the authoring route meets an application command line**
  ([#172](https://github.com/Bitspark/bittype/issues/172)). Not examined
  here.

## Where each row is stated

"Line" is the line of the [gate](#the-gate) whose tests validate the row.
O1 to O22 are the [decisions](#points-decided-after-the-first-version) that
bear on a row. A dash means that the row has no test of its own before the
gate is met.

| Row | Verdict | Stated under | Line |
| --- | --- | --- | --- |
| 1 | Holds | [What a contribution is](#what-a-contribution-is), [Amendments](#amendments-to-accepted-pages), O20 | 7 |
| 2 | Holds | [What a contribution is](#what-a-contribution-is), [Amendments](#amendments-to-accepted-pages), O12 | 1 |
| 3 | Adapted | [The protocol](#what-a-rule-does-the-protocol), O3, O9, O13 | 2 |
| 4 | Holds | [The protocol](#what-a-rule-does-the-protocol) | 1 |
| 5 | Holds | [The protocol](#what-a-rule-does-the-protocol), O1, O20 | 2 |
| 6 | Holds | [The protocol](#what-a-rule-does-the-protocol), [Amendments](#amendments-to-accepted-pages), O20 | 1 |
| 7 | Adapted | [Evidence](#evidence), O13 | 2 |
| 8 | Adapted | [Evidence](#evidence), O2, O20 | 2 |
| 9 | Adapted | [Evidence](#evidence), O17 | 2 |
| 9a | Deferred | [Deferred](#deferred-with-triggers), O6 | - |
| 10 | Adapted | [Refusals](#refusals), O17 | 2 |
| 11 | Adapted | [Optional fields](#optional-fields), [Amendments](#amendments-to-accepted-pages), O10 | The condition beside the seven lines |
| 12 | Adapted | [Optional fields](#optional-fields) | 2 |
| 13 | Adapted | [Optional fields](#optional-fields), O3 | 4 |
| 14 | Holds | [Optional fields](#optional-fields), [The extension domain](#the-extension-domain), O10 | 2 |
| 15 | Adapted | [Identities](#identities), O4, O16 | 2 |
| 15a | Deferred | [Deferred](#deferred-with-triggers), O16 | - |
| 16 | Holds | [Identities](#identities), O4, O19 | 3 |
| 17 | Adapted | [Identities](#identities), O4, O19, O20, O22 | 3 |
| 17a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 18 | Holds | [Identities](#identities), [Amendments](#amendments-to-accepted-pages) | 2 |
| 19 | Holds | [Identities](#identities) | 2 |
| 20 | Holds | [Identities](#identities), O5 | 2 |
| 21 | Holds | [Identities](#identities), [The extension domain](#the-extension-domain) | 2 |
| 22 | Adapted | [Identities](#identities), O20 | 4 |
| 23 | Adapted | [Identities](#identities), [Amendments](#amendments-to-accepted-pages), O7, O15, O16 | 6 |
| 24 | Adapted | [Identities](#identities) | 6 |
| 24a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 25 | Holds | [Identities](#identities) | 2, 5 |
| 26 | Adapted | [What the kernel carries](#what-the-kernel-carries-without-reading-it), O8, O11, O16, O20 | 2 |
| 26a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 27 | Adapted | [What the kernel carries](#what-the-kernel-carries-without-reading-it), O5, O19, O20 | 2, 3 |
| 28 | Adapted | [What the kernel carries](#what-the-kernel-carries-without-reading-it), O9, O12, O21 | 1, 3 |
| 29 | Adapted | [What lowering must be able to refer to](#what-lowering-must-be-able-to-refer-to), O8 | 5 |
| 30 | Holds | [What is final and what is provisional](#what-is-final-and-what-is-provisional), [Amendments](#amendments-to-accepted-pages), O17 | 7 |
| 31 | Adapted | [The order of work](#the-order-of-work), O6 | All |
| 31a | Deferred | [Deferred](#deferred-with-triggers), [The order of work](#the-order-of-work) | - |
| 31b | Deferred | [Deferred](#deferred-with-triggers), [The order of work](#the-order-of-work), O6 | - |
| 32 | Holds | [The order of work](#the-order-of-work) | 2, 6 |
| 33 | Holds | [The extension domain](#the-extension-domain), [Amendments](#amendments-to-accepted-pages), O10 | 7 |
| 34 | Adapted | [Inputs](#inputs-that-are-not-bittypes-alone) | 4 |
| 34a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 35 | Holds | [Inputs](#inputs-that-are-not-bittypes-alone), O14 | 4 |
| 36 | Stale | [Inputs](#inputs-that-are-not-bittypes-alone) | - |
| 37 | Adapted | [Inputs](#inputs-that-are-not-bittypes-alone), [The extension domain](#the-extension-domain), O18 | 7 |
| 37a | Deferred | [Deferred](#deferred-with-triggers) | - |
| 38 | Holds | [Payload bootstrapping](#payload-bootstrapping), O9 | 3 |
| 39 | Holds | [What a test must show](#what-a-test-must-show) | 1 |
| 40 | Adapted | [The gate](#the-gate), O10, O18, O20 | All |

18 rows hold, 21 are adapted, nine are deferred and one is stale.

## The code today

At `04632ae`. Nothing in this section is the design, and no fact here is a
reason for a choice above. Paths without a directory are in
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
at `ec96b92` and are moved to `04632ae`.

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
| 3 | A materially different contribution, with the shared packages fixed | Shown on the prototype branch, and again on `main` with step 4 |
| 4 | The contract in its final form and the migration of the 21 mapping packages, in one change | `main`, when the [gate](#the-gate) is met |
| 5 | The optional move | `main`, by the decision of [row 11](#optional-fields). The production record packages move here and not in step 4 (O18) |
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

**What the prototype contains (rows 31 and 32).**

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
- What the prototype fixes and reports of the points under
  [Not specified here](#what-is-delegated-and-what-is-not-specified-here).
- The four cases of row 32: "repeated uses share a node, equal native types
  with different adapters stay apart, an erased native argument is kept for
  conversion, one plan under two TypeScript placements".

**The six vertical cases in the prototype.** The cases are those of the plan
(`docs/representation-families.md:97-102`). Row 31 puts fixtures for cases 3
to 6 into the prototype. Row 31b, as O6 corrects it, keeps case 1 out and
case 2 at the kernel's level. O6 gives the reason: "Recursion is refused and
deferred (row 9a), so a mock cannot exercise it at the level of the kernel."
The form of each fixture below is derived here, from verification report 3.

| Case | Its form in the prototype | What that form does not show |
| --- | --- | --- |
| 1. Recursive `ContentValue` | Not in the prototype. O6: "Case 1 waits for #164 as a whole." | Anything of recursion beyond its refusal |
| 2. `Page<Outcome>` with optional nullable fields | At the kernel's level, through the independent fixture source, with a mock family that the prototype's written record names as a mock (O6). The optional field of `Page<Person>` runs through the field-state protocol | Null against absence, and discriminators |
| 3. The exact external generic `Cell<T>` | An exact generic with an omitted or an additional native parameter and a value facet | Capabilities, and what `Cell` does with its slot |
| 4. Go typed construction with a generic adapter argument | `Page[Person]` itself, lowered and executed | Lowering of families that do not exist |
| 5. A generic map whose key requirement is rejected | A fixture generic that requires evidence about its argument, accepted for one operand and refused for another | The map family and the meaning of its keys |
| 6. An independent family without a shared change | The contribution of step 3 | Authorship by another owner |

The six cases in real form stay what representation families makes them: the
condition for a stable version of the envelope
(`docs/representation-families.md:95-96`). The prototype does not discharge
that condition.

**What a stand-in proves (row 31).** Row 31: "the serialized plan is
sufficient for a second reader in another process, not for the service
library's factory." The table details that sentence. It is derived here,
from verification report 3.

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
The tests marked O20 are that decision's, and it calls them derived here.
Every test in the table is red today, for the reason given.

| Line | The predicate | Tests | Red today because |
| --- | --- | --- | --- |
| 1. Native projection | The native type is computed from the frozen record, by a projector that cannot read the source | A projector that tries to read the source does not compile or fails the contract test (row 4). Mutating the value a rule returned changes neither the native type nor the serialized plan (row 28). A fixture whose decision the source does not determine is green on the planner and red on a mutant that derives the entry from the source (rows 6 and 39, V1) | The one callback takes the reader. No record exists |
| 2. Dependency composition | A rule decides with typed answers about what was selected, and every reference is checked | A fixture family whose decision depends on evidence about its operand plans and generates (row 3). `Page<text>` and the written-out record generate the same native type, and a rule that reads one application fails (rows 7 and 8). Two contributions that differ only in their representation revision give equal payloads and different node identities (row 5, O20). An evidence cycle and a recursive reference are refused with different reasons (row 9, V1). Every refusal of the 21 mapping packages matches exactly one value (row 10). A request at a field position with the operand of another record fails (row 12, V2). One local key under two publishers registers, and one publisher with one key twice fails (row 15). Two uses of `Page<Person>` are two use entries on one node (row 18). Equal native types with different adapters stay apart (rows 19 and 32). An omitted native parameter still yields an adapter argument (row 20). A specialization is refused as unsupported with the application named (row 21). Two contributions that spell a helper alike get two symbols (row 25, V2). A member fragment that the target's member vocabulary does not admit is refused, with the class that the prototype fixes (row 14, V2). Guarantees 1 to 6 (rows 26 and 27) | A rule runs before its children. No evidence, no owner of an identity and no use entry exist |
| 3. Serialization | The plan, read from its serialized form in a fresh process, is the plan that was written | The completeness test runs from the serialized plan in a process that never held the planner's objects (V3). Each version axis is exercised alone (row 16). The contributions of one family for two targets carry one meaning identity and two representation revisions, and a reference that lacks either part is refused (row 17, O20). The payload types regenerate from their schema and every value round-trips through its codec (row 38, O9). A payload outside its schema is refused before a backend runs (row 38, V3) | No plan is serialized. No plan carries the axes |
| 4. A second reader in another process | A reader in another language converts from the plan and the snapshot alone | The TypeScript stand-in round-trips `Page<Person>` and refuses a payload schema version it does not know (V3). Given the plan without the snapshot it fails explicitly, and a snapshot with another pin is a linking error (row 34, V3). Each state of an optional field has its own vector, in Go and in the stand-in, and a present `undefined` is refused (rows 13 and 35, O14). The stand-in derives a memo key from the serialized plan alone: two uses with equal nodes, applications and selected dependencies get equal keys, and a different selected dependency gives a different key (row 22, O20) | No reader and no conversion exist |
| 5. Executed Go lowering | Go code that is lowered from the plan compiles and runs | Construction of `Page[Person]`, access to a member, an owned helper and a generic adapter argument compile and run (advice, section 3.4). A lowering that takes a member by a string the naming policy changed fails to compile, and with a member reference it compiles (row 29, V3) | A member is a string. No lowering exists |
| 6. Realization for TypeScript | A reader takes every native name and place from the realization | A plan read with the realization manifest of another output fails as a linking error, not as an incomplete plan (row 23). One plan under two TypeScript placements resolves each logical symbol to another module path and the same declaration (row 32, V2). No native name reaches emitted glue except through the reader, and a formatter mutant that realizes another spelling is followed (row 24, V2) | The inventory record has no planned identity, and the formatter's result has no reader |
| 7. The shared packages fixed | A materially different contribution is added, and no shared package changes | One registered contribution is replaced by another for the same domain, and the frozen records of every other selection are byte-identical (row 1, O20). The contribution registers, and `Page<Person>` regenerates with it. The diff of that test touches no shared package (row 40, V3). A Go module outside the repository authors and registers a contribution with public packages only (row 37, V3). One fixture for each item that is in the [extension domain](#the-extension-domain), and one recorded refusal for each that is out (row 33, V3) | A replacement reaches use rules only. No request names a product |

**Two checks that are green today and stay green.** The fixtures of today's
generics pass on the new contract with no contribution for generics (row 2,
V1). The native declarations of today's pagination fixtures do not change in
any of the four targets (#207, item 5).

**One condition beside the seven lines: row 11 is decided (O10).** "The
migration does not wait without limit: its gate requires row 11 to be
decided. It is decided by the answer on bitspark-service#38, or, if none has
come when the prototype is complete, by me on the evidence of the prototype,
recorded here with its reasons." The test of row 11 stays: "The answer on
bitspark-service#38 is recorded here with the form that was built. Red until
then".

**The gate on `main` (O18).** Four things that the lines above leave to
guess.

- **The tests of the delegated form run against fixtures.** "The gate's tests
  of the delegated form (the request at a field position, the member
  fragment, the four states) run on `main` against fixture contributions in
  the test tree. The production record packages move in step 5. If row 11 is
  decided as D1, those tests leave the gate with the form."
- **The codec on the list of shared packages** is the shared part of the
  role. "A family's codec is its own package, and a new contribution adds
  one."
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

**The shared packages of line 7.** Row 40 asks that they be "named and
fixed". Derived here, from the list of generation packages in `LAYOUT.md`:

- the contract package, `contracts/generation/go`;
- the registry, `transform/syntax/shared/generation/go`;
- the shared native construction, `transform/syntax/shared/target/go`;
- the source reader and the common mapping parts of the source,
  `transform/syntax/horos/source/go` and `transform/syntax/horos/mapping`;
- the native model and target of each language;
- every existing mapping package;
- the shared part of the payload codec role, which encodes, validates and
  freezes. A codec of one schema is not on this list (O18);
- the completeness test backend.

**What the gate depends on.** Derived here. Line 6 needs the reader of #273
and the inventory of pull request #281 with the planner's identity. Line 4
has its last state, a present `null`, only with #211. The contribution of
line 7 depends on the decision of row 11. The Go facets of line 5 are facet
packages of bittype's own contributions in the test tree (O8, O16).

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
