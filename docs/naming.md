# Naming: generation records an intent, each language realizes the spelling

**Accepted design, 12 November 2026, from research 0016 (rows 11, 13 and 14). Not implemented.**
Decided by [research 0016](../research-docs/0016-native-files-placement.md)
and by the [decisions on #274](https://github.com/Bitspark/bittype/issues/274#issuecomment-6837793448)
of its owner. It answers [#268](https://github.com/Bitspark/bittype/issues/268)
and is delivered by [#274](https://github.com/Bitspark/bittype/issues/274), per
target, in the [order](#delivery-order-per-target) at the end. No target has
changed a name. [The code today](#the-code-today) describes the present
behaviour, and nothing else on this page is implemented.

**How to read the sources.** Each statement names one of three sources.

- **An accepted document**, quoted, with its file and line.
- **The advice** of research 0016, section 2, "Names"
  ([advice](../research-docs/0016-native-files-placement.advice.md)), with the
  row of the table "Applied" of the research document that accepted it.
- **A decision on #274**, written D1 to D6 after the six points of the
  decisions comment, with the reason the comment gives. The integrator who
  wrote it is the owner of research 0016.

A sentence marked **derived here** follows from those sources and is stated by
none of them. Code references are `path:line` at commit `9961adb`.

## Three things that were one

A native declaration has one name today, a string that generation fixes
(`contracts/generation/go/native.go:11`). The design separates three things
that this string stood for.

| | What it is | Who states it | Source |
| --- | --- | --- | --- |
| **Declaration identity** | Which declaration this is, within one sealed output. It alone is authoritative for references. It is never a spelling | The generation kernel | Advice: "authoritative for references" (row 13). [`docs/source.md:58-59`](source.md#from-generated-declarations-to-files): "A spelling alone is insufficient: two distinct types named `Person` must not become interchangeable when placement changes." [`docs/binding-generation.md:101`](binding-generation.md#paired-contributions): "Native package paths and export names are not part of declaration identity." |
| **Naming intent** | What generation records about a name: the name the mapping chose and, apart from it, its qualification | The mapping, and for the qualification the consumer's composition | Advice: "the mapping’s chosen name or naming prescription" (row 13). D2 |
| **Realized spelling** | The identifier that a language's Files stage publishes | The language | Advice: "the identifier published by Files" (row 13). [`docs/source.md:50`](source.md#from-generated-declarations-to-files): "Realize the requested visibility under the language's rules and establish that uses still resolve to their intended targets." |

**Identity.** A use already identifies its destination by a link to the
declaration's node, not by its name:
"References to generated declarations and lexical parameters already identify
their destination nodes."
([`docs/generation-planner.md:99-100`](generation-planner.md#the-result-of-run)).
The name that a target also writes on each use
(`transform/syntax/golang/target/go/provider.go:42`) is a copy. The handle that
a sealed output issues for each declaration is the subject of
[#272](https://github.com/Bitspark/bittype/issues/272).

**Intent (D2).** "An intent contains the mapping's chosen name and its
qualification, kept apart. Today's spelling fuses the two." The fused form is
`"M" + decl.ModulePin[:16] + "_" + name`
(`transform/syntax/horos/mapping/declarations/go/declarations.go:46`). An
intent states:

- **the chosen name**, which the mapping derives from the source declaration,
  or which a helper recipe states for its helper;
- **the qualification**, which is either none or the module that the
  declaration comes from. The [qualification policy](#the-qualification-policy)
  decides which. A helper has no source module and no qualification
  ([#272, decision 3](https://github.com/Bitspark/bittype/issues/272#issuecomment-6837792602):
  "A helper has no source declaration.").

An intent states no visibility, no reserved-word escape and no placement. The
chosen name is not a spelling of any language. The horos side writes it with a
capital first letter, as #274 does in its tests ("the unexported spelling of
the intent `Person`"), and that capital says nothing about access.

The alternative was an intent that is the spelling written today, from which a
language derives its other spellings. It was rejected on #274 (point 2): "With
b the shared capitalization rule stays the source of every target's names."

**Where the intent travels.** In the inventory record of #272, "as an
owner-issued value" (row 10 of the table). The target issues and reads it. The
generation kernel carries it and does not interpret it.

**Realized spelling.** Advice, accepted as row 13: "Files selects the
appropriate realization, rewrites linked uses, checks collisions, and reports
the realized name."

**A name does not change between placements.** Exposure is fixed when an
output is sealed (`contracts/generation/go/output.go:52-56`), and a realization
depends on the intent, the exposure and the language's rule only. Advice: "A
fixed output has fixed exposure intent, so an exposed `Person` can remain
`Person` in every admitted placement." This keeps
[`docs/source.md:66-67`](source.md#from-generated-declarations-to-files): "Both
placements use the same generated declarations, bindings and publication
intent."

**Realization repairs nothing.** Advice: "If realization would produce an
invalid identifier, capture a builtin, or conflict with another declaration,
Files must use an explicitly authorized naming rule or fail."

## Who owns what

| Decision | Owner | Source |
| --- | --- | --- |
| The chosen name of a declaration | The mapping: the declaration recipe for a source declaration, the helper recipe for a helper | [`docs/source.md:48`](source.md#from-generated-declarations-to-files) gives the "Selected mapping and generation target" the task "Build declaration bodies and type expressions". Advice: "Generation should choose naming intent; Files should realize it under exposure." |
| The names inside one declaration: type parameters, members, duplicates among them, malformed identifiers | Generation | Advice: "Retain checks whose meaning is wholly inside a declaration: duplicate type parameters, malformed identifiers, and similar local defects." (row 8) |
| A member's name and accessibility | The declaration, so the mapping of that pair | [`docs/source.md:77-78`](source.md#from-generated-declarations-to-files): "The distinction is about module visibility. A member's accessibility, when it is part of the chosen type representation, still belongs to the declaration." |
| The qualification policy | The consumer's composition, as an argument of the declarations mapping of each pair, with a default | D3, row 11 |
| The case rule and the exported and unexported forms | Go | Row 13, D4. See [Go](#go) |
| Reserved words, and what is refused or escaped | Each target | D1. See [Reserved words](#reserved-words) |
| Collisions between declarations in a placed scope | The language's Files stage | Advice: "Collisions belong to the language’s placed scopes" (row 8). Delivered per target by [#273](https://github.com/Bitspark/bittype/issues/273) for TypeScript and by slice 1 of [#161](https://github.com/Bitspark/bittype/issues/161) and [#159](https://github.com/Bitspark/bittype/issues/159) |
| Rewriting the linked uses and reporting the realized name | The language's Files stage | Advice, row 13, quoted above |
| What "exposed" requires | Each language | Row 14. See [Exposure](#exposure-and-its-audience) |

**The mapping chooses the name.** The declaration recipe of the horos side does
so today (`transform/syntax/horos/mapping/declarations/go/declarations.go:41`),
and the optional helper's recipe states `BittypeOptional`
(`transform/syntax/horos-to-golang/mapping/optional/go/optional.go:24`). That
ownership stays.

**The consumer's composition selects the qualification policy (D3).** "The
qualification policy is selected in the consumer's composition, as an argument
of the declarations mapping of each pair, with a default. It is a choice about
the consumer's names. It is not a profile, so it does not contradict #265."

- **Why an argument of each pair's declarations mapping.** Row 11 asks for a
  default per target, and the pair's own package is where a target's default
  lives (#274, point 3). The consumer already passes a value of its own to a
  mapping library this way: the helper key in `record.Library(key)` and
  `optional.Library(key)`
  ([semantic generation](semantic-generation.md#pagination-subset)).
- **Why it is not a profile.**
  [`docs/semantic-generation.md:51`](semantic-generation.md#representation-and-enforcement-profiles)
  defines a profile: "A profile is a promise about the native type." The rule
  that a contribution states its own profile is about that promise: "A
  composition chooses which libraries it registers. It cannot choose what a
  builtin library promises." (`docs/semantic-generation.md:56-58`). A
  qualification changes no promise. The native type, its members and what it
  represents are the same under either policy, and a name is not part of what
  is promised: "Native package paths and export names are not part of
  declaration identity." (`docs/binding-generation.md:101`). So a composition
  that selects a policy chooses its own names and nothing that the library
  answers for.

**Each language owns its realization.** A reserved word, the case of a Go name
and the scope in which two names collide are facts of one language. A shared
mechanism holds none of them
([AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02)).
The shared contracts gain no access level and no keyword of a target language
(#274, "Exposure is the language's to interpret", from row 32).

**Two meanings of "qualification".** On this page a qualification is part of a
declaration's own name. The accepted decision uses the word for a use:
"Choose the local spelling, qualification or alias, and construct the
necessary import, export or other module syntax." (`docs/source.md:50`). There
it means how a use reaches a declaration in another scope, such as
`people.Person`. That stays with Files and is not changed here.

## Exposure and its audience

**`Exposed()` stays the list of requested roots** (row 14).
[`docs/generation-planner.md:100-101`](generation-planner.md#the-result-of-run):
"Exposed lists the realized sites requested by the consumer, not every
dependency or helper." The only access input of a request is
`PlanningInput.Roots` (`contracts/generation/go/planning.go:72`).

**The audience per language.** The accepted decision says "the requested
visibility under the language's rules" (`docs/source.md:50`) and does not say
to whom a declaration is visible. Row 14 states it, in the advice's words:

| Language | What an exposed declaration must be |
| --- | --- |
| TypeScript | "Nameable by an importer through its advertised module export" |
| Go | "Nameable from another package through its advertised package/name" |
| Rust | "Nameable from another crate through its advertised public path" |

C++ exposure is deferred until slice 1 of
[#160](https://github.com/Bitspark/bittype/issues/160) starts (row 32). It is
held here as one constraint: no shared operation means "private" in every
language.

**The two-sided rule** (advice, row 14): "Required access must be established,
and realized access must remain within the authorized access." A declaration
that nothing requires to be accessible gets only the language's local
visibility. Advice: "For nonexposed declarations, the initial policy should
authorize only the language’s local visibility."

**Required access has two owners (D5).** "The access a declaration requires has
two owners: the request, for its roots, and a contribution, for the helpers it
introduces."

- **The request** states its roots. That is the accepted decision: "The plan
  records which declaration the caller requested to expose."
  (`docs/source.md:48`).
- **A contribution** states the access of a helper that it introduces.
  [#272, decision 4](https://github.com/Bitspark/bittype/issues/272#issuecomment-6837792602):
  "a helper recipe states the access its declaration requires. The record
  carries it as an owner-issued value." The reason is that no request can
  state it. A root is a horos declaration site
  (`contracts/generation/go/planning.go:173-191`), and a helper has no source
  declaration.

**The Go optional helper (D5).** "An exposed record whose field has the
helper's type is usable from another package only if the helper is nameable
there." The helper is the type of every optional field of a Go record
(`transform/syntax/horos-to-golang/mapping/record/go/record.go:70-72`). "The
optional contribution states that its helper is part of the public
representation, and Files realizes it exported." The same holds for any helper
in a public representation. This replaces one item of row 13's action, which
listed "the optional helper across packages" among the refusals.

**Files never widens access on its own (D5).** "an unexposed source declaration
in the public representation of an exposed one is refused as unsupported
coverage, in Go as in Rust." Two accepted sentences forbid the other way out:
"moving a declaration cannot authorize exposing a private type or substituting
another type of the same name." (`docs/source.md:68-69`), and "An inaccessible
dependency cannot silently become public to make a placement work."
(`docs/generation-planner.md:236`). The advice says the same of Rust: "Files
must never respond by exporting `Hidden` without authorization." The refusal is
unsupported coverage and never invalid input (row 24, and
[#271](https://github.com/Bitspark/bittype/issues/271) for the category).

**What the public representation is.** It is what the audience meets when it
uses an exposed declaration: the declarations that its body names where the
audience can reach them. Which members the audience can reach is the
declaration's own statement (`docs/source.md:77-78`, quoted above). The exact
rule is each language's. Row 24 calls Rust's "A Rust-owned rule", and Go's is
written with slice 1 of #161.

**A consequence, derived here.** In the structural profile every member is
accessible in every target: the Go record mapping capitalizes each field
(`record.go:62`) and the Rust target writes every field public
(`transform/syntax/rust/target/go/provider.go:115`). An output holds the
requested roots and what they demand. So an unexposed source declaration is
always named in the body of an exposed one, or of a declaration that is itself
refused for that reason. Under D5 Go and Rust therefore admit, in this
profile, only outputs in which every source declaration is a requested root.
The unexported Go spelling is first realized for a source declaration when a
representation has members that its audience cannot reach, or when row 14a is
delivered. [Open point 4](#open-inside-this-design) asks whether that is
intended, and open point 5 asks what TypeScript does.

**How this relates to the decision of 7 November.** That decision names one
stater of visibility, the caller. A helper's stated access is a second one, and
it is D5 that adds it, not a reading of the decision. It keeps what the
decision forbids: no placement changes what is accessible, and Files exposes
nothing that an owner did not state before placement. `Exposed()` does not
list the helper.

**Deferred: row 14a.** "A record of access beyond the requested roots". It
stays deferred for access that a request states beyond its roots (D5), for
example access within one Rust crate. Its trigger: "a generation request can
state access other than its roots".

## The qualification policy

**Today.** A declaration of the module that the request enters through keeps
its chosen name. Every other declaration gets `M`, the first 16 hexadecimal
digits of its module's pin and `_` before the name
(`transform/syntax/horos/mapping/declarations/go/declarations.go:45-47`). The
prefix is mandatory: the recipe takes the two native forms and nothing else
(`declarations.go:26`). No accepted document gives it a purpose.
[Semantic generation](semantic-generation.md#generic-and-naming-limits)
describes it as present behaviour and adds: "This product spelling is not
semantic identity."

**What is wrong with it as a mandatory default.** The advice (row 11):

- "A content change changes the module pin and therefore potentially every
  prefixed public name."
- "A truncated pin is not a uniqueness guarantee."
- "Different versions can coexist in different native scopes without different
  local names."

The verification of the advice added a fourth (section "Applied"): the name
depends on which module the request enters through. `Root` is true for a
declaration of the entry module
(`transform/syntax/horos/source/go/read.go:62`), so the same horos declaration
is `Names` in one request and `M<pin>_Names` in another.

**Two policies.**

| Policy | What an intent carries | Result |
| --- | --- | --- |
| Qualified | The module qualification, for a declaration outside the entry module. This is today's rule | Today's output, byte for byte, in all four targets |
| Unqualified | No qualification | One horos declaration has the same native name whether its module is the entry of the request or an import |

Under the unqualified policy two declarations of the same name can be in one
output. Advice: "If two declarations named `Person` are deliberately placed in
one conflicting scope, refuse the placement. The consumer can choose separate
scopes or an explicit naming policy."

**The order in which the default changes** (row 11): "The qualification becomes
an explicit, selectable naming policy first, with today's as one option and
output unchanged. Unqualified names become the default per target once that
target's Files stage owns collisions."

1. The policy becomes selectable in all four pairs. The default is the
   qualified policy, and no output changes.
2. A target's default becomes the unqualified policy in the change that gives
   its Files stage the collision check for a placed scope: TypeScript after
   #273, Go in slice 1 of #161, Rust in slice 1 of #159. C++ keeps the
   qualified policy until #160.

The second step waits because the flat-name rule still refuses two equal names
in one output
(`transform/syntax/shared/target/go/check.go:17-19`,
`transform/syntax/shared/target/go/compile.go:22-23`). Without the prefix, and
before a Files stage can separate scopes, an output that generates today would
be refused.

**Never a hash added only when a collision appears** (advice, row 11): "That
makes adding an unrelated declaration capable of renaming an existing API." A
name depends on its own declaration and the selected policy, never on which
other declarations are in the output.

## Reserved words

**D1.** "A reserved word is the target's. A reserved word is a fact about one
language. The shared naming function knows none." This is the answer that #268
asks for. The alternative was a profile that promises one name in every
target. D1: "No accepted document promises equal native names across targets:
`docs/binding-generation.md` promises cross-target agreement on values and
protocol observations." The law reads: "Different native representations agree
on abstract values and protocol observations"
(`docs/binding-generation.md:145`).

**Where each language's words are.** Each language's syntax checker already
holds its own list. The one word that the shared function knows is Rust's.

| Language | Its reserved words | What arises from a chosen name | Rule |
| --- | --- | --- | --- |
| Go | Keywords, in the checker (`check/syntax/golang/provider/go/check.go:13`). Predeclared identifiers are in no list | Nothing in the exported spelling: every keyword and every predeclared identifier of Go begins with a lower-case letter. The unexported spelling can be either | Refusal, no escape. See [Go](#go) |
| Rust | `internal/rust/names.go:11` | `Self`, the only entry that begins with a capital | Rust renames `Self` to a legal name. See [open point 1](#open-inside-this-design) |
| TypeScript | `internal/typescript/names.go:55`, and the predefined type words at `:68` | Nothing: every entry begins with a lower-case letter | Refusal by TypeScript, if a later naming rule can produce one |
| C++ | `internal/cpp/names.go:40` | Nothing: every entry begins with a lower-case letter | Refusal by C++, if a later naming rule can produce one |

- **`Self`.** The shared function appends `_` to it for all four targets
  (`transform/syntax/horos/mapping/naming/go/naming.go:28-30`). After the
  delivery it holds no reserved word. A horos type named `Self` is `Self` in
  Go, TypeScript and C++ and a legal name in Rust (#274, done-when 7).
- **`type`, `func`, `class` and `interface`** (#268, done-when 3). A
  declaration or parameter of that name has the chosen name `Type`, `Func`,
  `Class` or `Interface`, and none of the four lists holds those. They do not
  arise in an exported or unqualified spelling. `Type` and `Func` are refused
  as unexported Go spellings.
- **Field names.** No target writes a field name in lower case. TypeScript,
  Rust and C++ keep the horos spelling
  (`transform/syntax/horos-to-typescript/mapping/record/go/record.go:55`,
  `transform/syntax/horos-to-rust/mapping/record/go/record.go:56`,
  `transform/syntax/horos-to-cpp/mapping/record/go/record.go:56`), and Go
  capitalizes it (`transform/syntax/horos-to-golang/mapping/record/go/record.go:62`).
  A member's name is the pair's mapping's, and a reserved word in it is that
  language's to refuse.
- **Capturing a predeclared or global name** is not a reserved-word question.
  A declaration named `string` in Go is legal and rebinds every use of the
  predeclared type. That is checked through recorded destinations,
  [#276](https://github.com/Bitspark/bittype/issues/276).

## Go

In Go the case of the first letter of a package-level name is its visibility.
Row 13's action: "It follows from the accepted decision: Files realizes
visibility under the language's rules, and in Go that is the case of the
name."

| | Spelling | When |
| --- | --- | --- |
| Exported | The intent as written: the chosen name, after its qualification if it has one. The chosen name of a source declaration begins with an upper-case ASCII letter (`naming.go:23-27`), and so does a qualified one (`declarations.go:46`) | The declaration is a requested root, or a helper whose contribution states that it is part of a public representation |
| Unexported | The same with its first letter in lower case. `Person` gives `person` | Every other declaration |

**Refusals (D4).** "Go's unexported spelling is the lower-case form, with
refusals where that form is a predeclared identifier, a keyword or a
collision." Examples: `String` would give `string`, `Any` would give `any`,
`Type` would give `type`. A collision is with another name of the package in
which the declaration is placed. Each refusal is unsupported coverage.

**No escape rule (D4).** "No escape rule: a refusal can be lifted later, a
spelling cannot be changed." A refusal that is lifted admits an output that
was refused and changes no name that exists. An escape fixes a spelling that
consumers then depend on.
[`docs/representation-families.md:26-28`](representation-families.md#delivered-complete-the-first-time)
has the same rule for native types: "An explicit refusal that can later be
lifted additively is complete; replacing a delivered native type (for example
`number` by `bigint`) is not."

**Fields.** A Go field's case is a member's accessibility and stays with the Go
record mapping (`docs/source.md:77-78`).

**Type parameters.** A type parameter is a binder inside one declaration. It
has no exposure and keeps its chosen name.

**The misplaced decision.** Go generation decides the case of every name
today, before exposure is known:

- The shared naming function capitalizes the first letter (`naming.go:25-26`).
  Its package comment says: "It is one rule for every target." (`naming.go:4`).
- The declaration recipe applies it to every declaration
  (`declarations.go:41`) and "is the same for every target"
  (`declarations.go:5`).
- The recipe runs before exposure exists. The planner declares a root and
  only then marks it exposed (`contracts/generation/go/planning.go:182-190`).
  A dependency that is not a root keeps `Exposed == false` and its capitalized
  name.
- The Go target writes that name on the declaration
  (`transform/syntax/golang/target/go/provider.go:76`).

So every generated Go declaration is exported, whether the caller exposed it
or not. A decision that the accepted design gives to Files is made by
generation, through a recipe that all four targets share.
[bittype/ACX-004](architecture-cases/0004-generation-file-syntax.md) records
it as a pending instance.

**Where the Go rule lives.** In a Go-owned package that the Go target imports,
and later Go Files (#274, done-when 9). For an intent it returns the exported
spelling, and the unexported spelling or a refusal. The model is
`transform/syntax/typescript/symbol/go`, which "owns the TypeScript public type
symbols shared by generation and Files". Go Files realizes a spelling in
slice 1 of #161.

## What must not change after the first slice records a name

[#207](https://github.com/Bitspark/bittype/issues/207) and
[#195](https://github.com/Bitspark/bittype/issues/195) record declarations in
plan entries. A plan entry records a declaration's identity and naming intent,
and no realized native name (#274, done-when 2). From then on a later change
may not alter
([AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10):
a first change establishes the final structure):

- **the record of intent:** that it states the chosen name and the
  qualification, kept apart;
- **the default policy of each target**, as the
  [order](#delivery-order-per-target) gives it;
- **each language's realization rule:** what "exposed" requires in that
  language, the Go spellings and their refusals, and each target's reserved
  words with what it refuses or renames.

What may change:

- **A refusal that is lifted** (D4).
- **A field added to the inventory record.** #272: "Adding a field to the
  record is additive."

The default of a target changes once, in the step that the order names. For
Go and Rust that step is slice 1 of #161 and #159, which is later than the
other parts of #274. [Open point 6](#open-inside-this-design) asks what a plan
entry records for those targets before it. C++ is the stated exception: its
names keep today's rule until #160, and the guarantee of research 0016
excludes legacy C++ output (row 36).

## How a consumer learns the realized name and the placed site

**D6.** "A consumer learns the realized name and the placed site from a reader
on the formatter's result, keyed by the output's handles." The handles are
those of #272.

- **Why a reader is needed.** `Source` reports neither today. The TypeScript
  formatter's map from an original site to its placed site is private
  (`format/files/typescript/provider/go/relocate.go:28`).
- **Why on the formatter's result.** The realized name exists only after a
  language's Files stage has realized it, and the placed site only after
  placement. The alternative was to read both from the published semantic
  graph (#274, point 6, option b). With the reader, #195 pins identity and
  intent in a plan entry and reads the spelling after placement.

## Reconciliation with accepted documents

**[Representation families](representation-families.md#families), lines 207
to 209.** "A lifted inner binder gets a deterministic name, with collision and
reserved-word handling, and keeps its free generic parameters." The sentence
does not say who handles a collision or how. Read with rows 8 and 9 and with
D1:

| Part of that sentence | Where it is |
| --- | --- |
| "a deterministic name" | Generation. The recursion family states a deterministic **intent** for the lifted binder. It depends on the declaration that the binder is lifted from and on the binder, never on which other declarations are in the output |
| Collisions among the names of one source declaration: its type parameters, its members, and the binders lifted from it | Generation. Their meaning is wholly inside one declaration |
| Collisions between top-level declarations, which a lifted binder becomes | The language's Files stage, in the placed scope. It refuses and does not rename (row 8, and "never a hash", row 11) |
| "reserved-word handling" | The target (D1) |

One reading of that sentence does not survive: that generation gives a lifted
binder another name when its name meets that of another declaration of the
output. Generation does not know the placed scope, and a name must not depend
on which other declarations are in the output.

**[Binding generation](binding-generation.md#the-selected-representation-plan).**
The plan carries "native symbols, distinguishing run-time imports from
type-only imports" (`docs/binding-generation.md:56`), and the shared envelope
of [representation families](representation-families.md#the-plan) lists
"Native symbols, imports, factory references, runtime requirements". For a
declaration that generation owns, this design reads "native symbol" as the
declaration's identity and naming intent. Its realized name and placed site
come from the reader of D6. A symbol of an external package keeps the name
that its owner published. The completeness test of that page says: "The
backend can implement each selected representation from the plan, the source
semantics and the selected contributions alone, without guessing a rule."
(`docs/binding-generation.md:72`). Reading a name from the formatter's result
guesses no rule. It is still an input beyond the three that the sentence
lists, and the owner of #195 has to confirm this reading. It is
[open point 7](#open-inside-this-design).

**[Source and Files](source.md#from-generated-declarations-to-files).** No
sentence of the decision of 7 November says who spells a declaration's name.
The statement that Files never renames a declaration was a derivation of the
research document (its statement 18) and is withdrawn (row 13). The one place
where this design adds to the decision is the second owner of required access,
described under [Exposure](#exposure-and-its-audience).

## Open inside this design

1. **Rust's legal name for `Self`.** #274 requires "a legal name in Rust" and
   does not say which.
   - Option a: `Self_`, which Rust output has today. #268, done-when 4, keeps
     today's output "except for a name the answer deliberately changes".
   - Option b: a refusal, by the reason D4 gives for Go: "a refusal can be
     lifted later, a spelling cannot be changed." This changes done-when 7 of
     #274.
2. **What a qualification is, as a recorded value.** D2 keeps the chosen name
   and the qualification apart and does not say what the qualification holds.
   Today it is spelled as `M`, 16 digits of the pin and `_`
   (`declarations.go:46`).
   - Option a: the module's identity. Each language's realization spells it,
     and under the qualified policy it spells today's form.
   - Option b: the spelled prefix. The truncation is then part of what a plan
     entry records.
3. **The other two rules of the shared naming function.** It puts `Type`
   before a leading underscore (`naming.go:23-24`) and refuses a name that
   contains `__` (`naming.go:15`). No document gives either a purpose. A name
   that contains `__`, or that begins with an underscore and a capital, is
   reserved in C++. #268 excludes both rules "only if answer 1 is taken
   further", and #274 takes it further.
   - Option a: they are restrictions of the chosen name, and stay with the
     horos side.
   - Option b: they are C++'s reserved identifiers, and move to C++ with
     done-when 7 of #274.
4. **Which unexposed source declarations Go and Rust admit.** The consequence
   derived under [Exposure](#exposure-and-its-audience): in the structural
   profile, none. Row 13's validation expects that "an unexposed Go declaration
   is realized unexported and its uses are rewritten".
   - Option a: the consequence stands. The unexported spelling is tested on
     the Go rule itself, slice 1 of #161 tests the refusal, and the row's
     validation is met by the first representation with members that its
     audience cannot reach.
   - Option b: the public representation is narrower than every declaration
     that a body names. Then Go and Rust each need a stated rule for it before
     their slice 1.
5. **TypeScript and an unexposed declaration in the public representation of an
   exposed one.** D5 names Go and Rust. TypeScript admits the case today when
   both are in one file, and refuses a use across files
   (`format/files/typescript/provider/go/relocate.go:163-165`).
   `TestGeneratedFilesPreservePublicationIntent`
   (`integration/go/typescript_placement_test.go:323-376`) pins both.
   - Option a: TypeScript keeps admitting it. A TypeScript type is structural,
     so an importer can read and build the value without naming the type.
   - Option b: TypeScript refuses it as Go and Rust do. The test changes.
6. **What a plan entry records for Go and Rust before their slice 1.** The
   default policy of Go and Rust changes in slice 1 of #161 and #159. A plan
   entry that #207 writes before that, for a consumer that passes no policy,
   records a qualified intent that the slice then changes.
   - Option a: #207 records no intent for Go and Rust before their slice 1.
   - Option b: the default of Go and Rust changes with the other parts of
     #274. That needs their collision check earlier than the table says.
7. **The plan's native symbols.** The reading under
   [Reconciliation](#reconciliation-with-accepted-documents) is this page's.
   The owner of #195 confirms it or states another.

## Delivery order per target

From #274, "Order", and rows 8, 10, 11 and 12 of the table "Applied".

| Step | Needs | Item |
| --- | --- | --- |
| This design | Nothing | #274, done-when 1 to 4 |
| The inventory with handles. Each record can carry an intent | Nothing | #272 (row 10) |
| The builder's comparator orders by name, then origin, then planned identity | In the same change as the first removal of the flat-name rule | #272 (row 12) |
| The policy is selectable. The qualified policy reproduces today's output in all four targets | Nothing | #274, done-when 5 and 6 |
| The shared naming function holds no reserved word. Rust states its own | Nothing | #274, done-when 7, and #268 |
| Every inventory record carries an intent for TypeScript, Go and Rust. A Go-owned package holds the Go rule | #272 | #274, done-when 8 and 9 |

Then per target (row 8: "TypeScript first, Go and Rust with their slice 1, C++
with #160"):

| Target | The flat-name rule leaves, and Files checks collisions per placed scope | Files realizes the spelling | The default becomes unqualified |
| --- | --- | --- | --- |
| TypeScript | #273, second change, after #272 | TypeScript has no case rule. Its Files stage already writes the export from exposure (`relocate.go:89-91`) | After #273 (#274, done-when 10) |
| Go | Slice 1 of #161 | Slice 1 of #161: exported or unexported, uses rewritten | Slice 1 of #161 |
| Rust | Slice 1 of #159 | Slice 1 of #159 realizes visibility. The name changes only for `Self` | Slice 1 of #159 |
| C++ | With #160. Until then C++ selects the flat rule in its own target | Deferred until [#167](https://github.com/Bitspark/bittype/issues/167) supplies a C++ semantic provider | With #160. Until then the qualified policy |

All of it lands before #207 and #195 record a native name.

## The code today

At `9961adb`. Nothing in this section is the design.

| Fact | Where |
| --- | --- |
| One function spells every target's names: it capitalizes the first letter, puts `Type` before a leading underscore, turns `Self` into `Self_` and refuses `__` | `transform/syntax/horos/mapping/naming/go/naming.go:14-32` |
| The declaration recipe applies it to every declaration and every type parameter | `transform/syntax/horos/mapping/declarations/go/declarations.go:41`, `:62` |
| The module-pin prefix, for a declaration outside the entry module | `declarations.go:45-47`, `transform/syntax/horos/source/go/read.go:62` |
| Nothing selects a naming policy. Each pair calls the recipe with its two forms | `declarations.go:26`, and line 23 of `transform/syntax/horos-to-typescript/mapping/declarations/go/declarations.go` and of the Go pair, line 22 of the Rust and C++ pairs |
| A native declaration's name is one string | `contracts/generation/go/native.go:11` |
| Exposure is set after the recipe ran, and only for the requested roots | `contracts/generation/go/planning.go:182-190` |
| `Exposed()` holds the sites of the requested roots | `contracts/generation/go/output.go:42-43`, `:59` |
| The Go helper's name is a literal, and the record mapping wraps every optional field in the helper | `transform/syntax/horos-to-golang/mapping/optional/go/optional.go:24`, `transform/syntax/horos-to-golang/mapping/record/go/record.go:70-72` |
| The Go target writes the name on the declaration and on each use | `transform/syntax/golang/target/go/provider.go:76`, `:42` |
| The Rust target writes every struct, alias and field public | `transform/syntax/rust/target/go/provider.go:92`, `:95`, `:115` |
| Two equal names in one output are refused, in each target's program rule and again in the shared builder | `transform/syntax/shared/target/go/check.go:17-19`, `transform/syntax/shared/target/go/compile.go:22-23` |
| TypeScript Files reads a requested root as exported and writes the export | `format/files/typescript/provider/go/generated.go:92-97`, `format/files/typescript/provider/go/relocate.go:89-91` |
| TypeScript Files refuses a use, from another file, of a declaration that is not exposed | `relocate.go:163-165` |
| TypeScript Files compares the realized exposure with the requested one in both directions | `format/files/typescript/provider/go/format.go:201-205` |
| TypeScript Files refuses two declarations of one name in one file | `relocate.go:120-126` |
| TypeScript Files spells a use in the same file from its destination declaration | `relocate.go:155-160` |
