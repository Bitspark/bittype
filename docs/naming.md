# Naming: generation records an intent, each language realizes the spelling

**Accepted design, 13 November 2026, from research 0016 (rows 11, 13 and 14). Not implemented.**
Decided by [research 0016](../research-docs/0016-native-files-placement.md)
and by the comments of its owner on
[#274](https://github.com/Bitspark/bittype/issues/274). It answers
[#268](https://github.com/Bitspark/bittype/issues/268) and is delivered by
#274, per target, in the [order](#delivery-order-per-target) near the end. No
target has changed a name. [The code today](#the-code-today) describes the
present behaviour, and nothing else on this page is implemented.

**How to read the sources.** Each statement names its source.

- **An accepted document**, quoted, with its file and line.
- **The advice** of research 0016, section 2, "Names"
  ([advice](../research-docs/0016-native-files-placement.advice.md)), with the
  row of the table "Applied" of the research document that accepted it.
- **D1 to D6:** the six
  [decisions of 12 November 2026](https://github.com/Bitspark/bittype/issues/274#issuecomment-6837793448)
  on #274. The integrator who wrote them is the owner of research 0016.
- **E1 to E7:** the seven
  [decisions of 13 November 2026](https://github.com/Bitspark/bittype/issues/274#issuecomment-6838436844)
  on the points that the first version of this page left open. E4 defers one
  thing that D4 defined, Go's unexported spelling.
- **R8 to R15:** the eight
  [decisions after an independent reading](https://github.com/Bitspark/bittype/issues/274#issuecomment-6839006959)
  of this page, also of 13 November 2026. They are numbered 8 to 15 in that
  comment.
- **R14c:** the [correction of decision 14, point 3](https://github.com/Bitspark/bittype/issues/274#issuecomment-6839289646), of 13 November
  2026, after a second reading. It withdraws that point and accepts four small
  findings.

A sentence marked **derived here** follows from those sources and is stated by
none of them. Code references are `path:line` at commit `9961adb`.

## Three things that were one

A native declaration has one name today, a string that generation fixes
(`contracts/generation/go/native.go:11`). The design separates three things
that this string stood for.

| | What it is | Who states it | Source |
| --- | --- | --- | --- |
| **Declaration identity** | Which declaration this is, within one sealed output. It alone is authoritative for references. It is never a spelling | The generation kernel | Advice: "authoritative for references" (row 13). [`docs/source.md:58-59`](source.md#from-generated-declarations-to-files): "A spelling alone is insufficient: two distinct types named `Person` must not become interchangeable when placement changes." [`docs/binding-generation.md:101`](binding-generation.md#paired-contributions): "Native package paths and export names are not part of declaration identity." |
| **Naming intent** | What generation records about a name: the chosen name, the module's identity, the policy's identity and the spelling the policy produced | The mapping, with the policy that the consumer's composition gave it | Advice: "the mapping’s chosen name or naming prescription" (row 13). D2, E2, R10 |
| **Realized spelling** | The identifier that a language's Files stage publishes. In the first slices it is the spelled intent, once that language's rule has checked it | The language | Advice: "the identifier published by Files" (row 13). [`docs/source.md:50`](source.md#from-generated-declarations-to-files): "Realize the requested visibility under the language's rules and establish that uses still resolve to their intended targets." R8, R13 |

**Identity.** A use already identifies its destination by a link to the
declaration's node, not by its name:
"References to generated declarations and lexical parameters already identify
their destination nodes."
([`docs/generation-planner.md:99-100`](generation-planner.md#the-result-of-run)).
The name that a target also writes on each use
(`transform/syntax/golang/target/go/provider.go:42`) is a copy. The handle that
a sealed output issues for each declaration is the subject of
[#272](https://github.com/Bitspark/bittype/issues/272). The identity that a
plan entry records is another question, and this page does not define it
([What is fixed, and when](#what-is-fixed-and-when)).

**Intent (D2, E2, R10).** D2: "An intent contains the mapping's chosen name and
its qualification, kept apart. Today's spelling fuses the two." The fused form
is `"M" + decl.ModulePin[:16] + "_" + name`
(`transform/syntax/horos/mapping/declarations/go/declarations.go:46`). R10 says
what is recorded: "An intent records the chosen name, the module's identity
(its full pin), the policy's identity, and the spelling the policy produced."

- **The chosen name.** The mapping derives it from the source declaration, or
  a helper recipe states it for its helper.
- **The module's identity.** E2: "A qualification is recorded as the module's
  identity." It is the full pin of the module that the declaration comes from,
  not a truncation. A helper has no source module
  ([#272, decision 3](https://github.com/Bitspark/bittype/issues/272#issuecomment-6837792602):
  "A helper has no source declaration.").
- **The policy's identity.** R10: "The policy is a named value that the
  consumer's composition gives to the declarations mapping." And: "A second
  qualified policy is then told from the first by its identity."
- **The spelling the policy produced.** E2: "The spelled prefix is one
  policy's realization of it." Under the qualified policy the spelling has
  today's prefix, and under the unqualified policy it is the chosen name
  ([The qualification policy](#the-qualification-policy)).

**Who spells, and where the spelling is (R10).** "The mapping applies it at
generation." And: "The unplaced tree's name is that spelling, and the
builder's comparator and every Files stage read it there." So the name of a
native declaration and the `name` attribute of the unplaced tree stay one
string, the spelled intent. No target and no Files stage turns a module
identity into a prefix. R10: "Names stay with the mapping, as `docs/source.md`
has it, and each language checks or refuses the spelling." The words about
`docs/source.md` are a reading, and R14c says whose: "my reading of that
page's sentence about what generation builds". That sentence gives the
selected mapping and the generation target the task "Build declaration bodies
and type expressions" (`docs/source.md:48`), and this page uses the same
reading under [Who owns what](#who-owns-what). The decision of 7 November has
no sentence that says who spells a name.

An intent states no visibility, no reserved-word escape and no placement.

**The chosen name of a source declaration (E3).** It is the horos name in a
capitalised form: the first letter in upper case, and `Type` before a leading
underscore. E3: "The `Type` prefix before a leading underscore stays with the
chosen name: it is the mapping's way to give every name a capitalised form,
the same for all targets." R8: "That the mapping capitalises a chosen name is
the mapping's style for every target and says nothing about access." The
other two things that the shared naming function does today leave that
function: the `Self` rule and the refusal of `__`, both under
[Reserved words](#reserved-words).

**The alternative that was not taken.** An intent that is only the spelling
written today, from which a language derives its other spellings (#274,
point 2, option b). The issue's reason was: "With b the shared capitalization
rule stays the source of every target's names." That reason does not separate
the two options, because under E3 the capitalised form stays the mapping's
for every target. Derived here, what separates them: with b the fused string
is the only record, and a reader has to take it apart to learn the name, the
module and the rule that joined them. With the chosen design the parts are
recorded beside the spelling, and no language derives one spelling from
another. Go's rule checks the spelling and "renames nothing" (R8).

**Where the intent travels.** In the inventory record of #272, "as an
owner-issued value" (row 10 of the table). The generation kernel carries it
and does not interpret it. How the four parts reach that record from the
mapping is not specified here
([What is delegated](#what-is-delegated-and-what-is-not-specified-here)).

**Realized spelling.** Advice, accepted as row 13: "Files selects the
appropriate realization, rewrites linked uses, checks collisions, and reports
the realized name." In the first slices the selection is a check: a Files
stage publishes the spelled intent or refuses (R8, R13).

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
Files must use an explicitly authorized naming rule or fail." Of the two, the
first slices take the second for every target. R13: "in the first slices a
target refuses a spelling it cannot publish. An escape rule is a later,
explicit rule that turns one refusal into a success."

## Who owns what

| Decision | Owner | Source |
| --- | --- | --- |
| The chosen name of a declaration | The mapping: the declaration recipe for a source declaration, the helper recipe for a helper | [`docs/source.md:48`](source.md#from-generated-declarations-to-files) gives the "Selected mapping and generation target" the task "Build declaration bodies and type expressions". Advice: "Generation should choose naming intent; Files should realize it under exposure." |
| The spelling of an intent | The mapping, which applies the selected policy at generation | R10 |
| The qualification policy | The consumer's composition. It gives a named value to the declarations mapping of each pair, which has a default | D3, row 11, R10 |
| The names inside one declaration: type parameters, members, duplicates among them, malformed identifiers | Generation | Advice: "Retain checks whose meaning is wholly inside a declaration: duplicate type parameters, malformed identifiers, and similar local defects." (row 8) |
| A member's name and accessibility | The declaration, so the mapping of that pair | [`docs/source.md:77-78`](source.md#from-generated-declarations-to-files): "The distinction is about module visibility. A member's accessibility, when it is part of the chosen type representation, still belongs to the declaration." |
| The access that a declaration requires | The request for its roots. A contribution for every declaration it introduces | D5, R11. See [Exposure](#exposure-and-its-audience) |
| Go's rule: the check that the spelled intent of a public declaration is exported | Go | Row 13, R8. See [Go](#go) |
| Reserved words, and what is refused | Each target | D1, E1, E3, R13. See [Reserved words](#reserved-words) |
| Collisions between declarations in a placed scope | The language's Files stage | Advice: "Collisions belong to the language’s placed scopes" (row 8). Not delivered: [#273](https://github.com/Bitspark/bittype/issues/273) is to deliver it for TypeScript, and slice 1 of [#161](https://github.com/Bitspark/bittype/issues/161) and of [#159](https://github.com/Bitspark/bittype/issues/159) for Go and Rust |
| A type parameter that captures a use of a declaration in its placed scope | The language's Files stage, which refuses the placement | Advice, row 8: "Files must either perform an explicitly permitted alpha-renaming or refuse that placement." R13 takes the refusal for the first slices |
| How a use reaches its declaration, and reporting the realized name | The language's Files stage | `docs/source.md:50`. Advice, row 13, quoted above |
| What "exposed" requires, and which outputs are admitted | Each language | Row 14, E4, E5. See [Exposure](#exposure-and-its-audience) |

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
- **Why it is not a profile.** D3 decides it. The argument is derived here.
  [`docs/semantic-generation.md:51`](semantic-generation.md#representation-and-enforcement-profiles)
  defines a profile: "A profile is a promise about the native type." The rule
  that a contribution states its own profile is about that promise: "A
  composition chooses which libraries it registers. It cannot choose what a
  builtin library promises." (`docs/semantic-generation.md:56-58`). A policy
  changes how a declaration is spelled. The native type, its members and what
  it represents are the same under either policy. So a composition that
  selects a policy chooses its own names and nothing that a library answers
  for.

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

**Required access has two owners (D5, R11).** D5: "The access a declaration
requires has two owners: the request, for its roots, and a contribution, for
the helpers it introduces." R11 widens the second: "The second is wider than
helpers. A contribution states the access of every declaration it introduces."

- **The request** states its roots. That is the accepted decision: "The plan
  records which declaration the caller requested to expose."
  (`docs/source.md:48`).
- **A contribution** states the access of a declaration that it introduces.
  R14c says what that is: it "means one whose access no request can state". A
  source declaration is not one, because a request can name it as a root. A
  helper is one case.
  [#272, decision 4](https://github.com/Bitspark/bittype/issues/272#issuecomment-6837792602):
  "a helper recipe states the access its declaration requires. The record
  carries it as an owner-issued value." The reason is that no request can
  state it. A root is a horos declaration site
  (`contracts/generation/go/planning.go:173-191`), and a helper has no source
  declaration.
- **A lifted inner binder is another case, and it is not decided here.** R11:
  "A lifted inner binder (#164) is one, and its contribution states its access
  when that family is built. No rule for it is invented here."

**The Go optional helper (D5).** "An exposed record whose field has the
helper's type is usable from another package only if the helper is nameable
there." The helper is the type of every optional field of a Go record
(`transform/syntax/horos-to-golang/mapping/record/go/record.go:70-72`). "The
optional contribution states that its helper is part of the public
representation, and Files realizes it exported." The same holds for any helper
in a public representation. This replaces one item of row 13's action, which
listed "the optional helper across packages" among the refusals. R15: "the
helper is realized exported and is no refusal."

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
rule is each language's. Row 24 calls Rust's "A Rust-owned rule". R14: "Go's
exact rule for what is in a public representation is written in slice 1 of
#161. The set of admitted outputs only grows afterwards."

**Which outputs Go and Rust admit (E4).** "Go and Rust admit only outputs in
which every source declaration in the public representation of an exposed one
is itself a requested root." E4 accepts this as the consequence of D5, and
says how far it reaches: "In the structural profile every member of a Go or
Rust record is accessible, so this covers every dependency." The Go record
mapping capitalizes each field (`record.go:62`), and the Rust target writes
every field public (`transform/syntax/rust/target/go/provider.go:115`).

- **The refusal says what to do.** E4: "The refusal names the declaration to
  request." R14: "A refusal names every offending declaration, each by its
  handle and its spelled intent."
- **Why a refusal and not an export.** For a source declaration, E4: "The
  authority to expose stays with the caller". That is `docs/source.md:48`:
  "The plan records which declaration the caller requested to expose." For a
  declaration that a contribution introduces, the contribution has that
  authority (R11).
- **What follows for Go's unexported spelling.** E4: "So Go's unexported
  spelling is realized by nothing in the first slices." It is deferred. See
  [Go](#go).

**The rule is each language's (E5).**

| Language | An unexposed source declaration in the body of an exposed one | Reason, in E5's words |
| --- | --- | --- |
| TypeScript | Admitted within one module. A placement that separates the two is refused as invalid input | "because structural typing keeps the exposed type usable there" |
| Go | Refused as unsupported coverage | "because a caller in another package or crate could not construct the value" |
| Rust | Refused as unsupported coverage | The same |

`TestGeneratedFilesPreservePublicationIntent`
(`integration/go/typescript_placement_test.go:323-376`) pins the TypeScript
behaviour, and it stays (E5).

**Which category a refusal has (R14c).** The rule, stated once: "a refusal is
unsupported coverage when a later slice can admit the same input. It is
invalid input when no later slice can."

- **Go and Rust: unsupported coverage.** "Go and Rust refuse the output
  whatever its placement, and a later rule of the language can admit that
  output unchanged, so their refusal is unsupported coverage."
- **TypeScript: invalid input.** "TypeScript admits the output and refuses one
  placement of it. A caller who wants that placement has to change the
  request, which makes it a different input." The refusal "stays invalid
  input, as #271 decided on 12 November", for the reason #271 gives: "it can
  never succeed for the same output, because Files never exports a
  declaration without authorization." It is invalid input today
  (`format/files/typescript/provider/go/relocate.go:163-165`), and that does
  not change.
- **What is withdrawn.** Point 3 of R14 said that this TypeScript refusal is
  unsupported coverage, as in Go and Rust. It reversed the decision on #271
  without saying so, and R14c withdraws it.

**How this relates to the decision of 7 November.** That decision names one
stater of visibility, the caller. A contribution's stated access is a second
one, and it is D5 and R11 that add it, not a reading of the decision. It keeps
what the decision forbids: no placement changes what is accessible, and Files
exposes nothing that an owner did not state before placement. `Exposed()` does
not list the helper. R14: "A consumer learns realized names and access from
the formatter's result. Nothing offers them before formatting."

**Deferred: row 14a.** "A record of access beyond the requested roots". Its
trigger: "a generation request can state access other than its roots". An
example is access within one Rust crate. R15 says what the row does not
defer: "Row 14a defers access that a request states beyond its roots, while a
contribution's stated access is delivered with #272."

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

**Two policies.** Each is a named value with its own identity, and the mapping
applies the selected one at generation (R10).

| Policy | The spelling it produces | Result |
| --- | --- | --- |
| Qualified | Today's rule: the chosen name for a declaration of the entry module, and `M`, 16 digits of the pin, `_` and the chosen name for every other | Today's output, byte for byte, in all four targets |
| Unqualified | The chosen name | One horos declaration has the same native name whether its module is the entry of the request or an import |

Under either policy an intent records the module's full pin and the policy's
identity beside the spelling (R10).

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
be refused. When each change must have landed is under
[What is fixed, and when](#what-is-fixed-and-when).

**Never a hash added only when a collision appears** (advice, row 11): "That
makes adding an unrelated declaration capable of renaming an existing API." A
name never depends on which other declarations are in the output. Under the
unqualified policy it depends on its own declaration only. Under the qualified
policy it also depends on which module the request enters through, which is
the fourth defect above.

## Reserved words

**D1.** "A reserved word is the target's. A reserved word is a fact about one
language. The shared naming function knows none." This is the answer that #268
asks for. The alternative was a profile that promises one name in every
target. D1: "No accepted document promises equal native names across targets:
`docs/binding-generation.md` promises cross-target agreement on values and
protocol observations." The law reads: "Different native representations agree
on abstract values and protocol observations"
(`docs/binding-generation.md:145`).

**In the first slices a target refuses, and does not rename (R13).** "in the
first slices a target refuses a spelling it cannot publish. An escape rule is
a later, explicit rule that turns one refusal into a success." R13 states for
every target what three decisions had said of one case each: D4 of Go's
unexported form, E1 of Rust and `Self`, E3 of C++ and `__`. The answer on #268
and the advice both leave an escape possible ("an explicitly authorized
naming rule or fail"), and R13 does not exclude one later. The reason is E1's:
"A refusal can be lifted later, a spelling cannot be changed."
[`docs/representation-families.md:26-28`](representation-families.md#delivered-complete-the-first-time)
has the same rule for native types: "An explicit refusal that can later be
lifted additively is complete; replacing a delivered native type (for example
`number` by `bigint`) is not."

**Where each language's words are.** Each language's syntax checker already
holds its own list. The one word that the shared function knows is Rust's.

| Language | Its reserved words | What the chosen name of a source declaration can meet | Rule |
| --- | --- | --- | --- |
| Go | Keywords, in the checker (`check/syntax/golang/provider/go/check.go:13`). Predeclared identifiers are in no list | Nothing: every keyword and every predeclared identifier of Go begins with a lower-case letter, and the mapping writes the chosen name of a source declaration with a capital. A helper's name is its recipe's and can be any | Go's [check](#go) refuses a public declaration whose spelling is not exported. That covers every keyword and every predeclared identifier |
| Rust | `internal/rust/names.go:11` | `Self`, the only entry that begins with a capital | Rust refuses `Self` as unsupported coverage (E1) |
| TypeScript | `internal/typescript/names.go:11`, and the predefined type words at `:24` | Nothing: every entry begins with a lower-case letter | TypeScript refuses a reserved word, if a helper's name or a later naming rule produces one (R13) |
| C++ | `internal/cpp/names.go:11`. A name that contains `__` is a reserved identifier of C++ and is in no list | No keyword: every entry begins with a lower-case letter. A chosen name that contains `__`, which the shared function refuses today | C++ refuses a name that contains `__` (E3) |

- **`Self` (E1).** "Rust refuses it as unsupported coverage. The other targets
  realize `Self` unchanged." Today the shared function appends `_` to it for
  all four targets
  (`transform/syntax/horos/mapping/naming/go/naming.go:28-30`). After the
  delivery the function holds no reserved word. A horos type named `Self` is
  `Self` in Go, TypeScript and C++, and Rust refuses it, where Rust output has
  `Self_` today. Two issue texts still say "a legal name in Rust": #274
  (done-when 7) and #268 (done-when 2). E1 replaces both.
- **A name that contains `__` (E3).** The shared function refuses it for every
  target today (`naming.go:15`). E3: "The refusal of `__` is a
  reserved-identifier rule of C++ and moves to the C++ target. The other
  targets admit such a name."
- **`Type` before a leading underscore (E3)** is no reserved-word rule. It
  stays with the chosen name, as [above](#three-things-that-were-one).
- **`type`, `func`, `class` and `interface`** (#268, done-when 3). A
  declaration or parameter of that name has the chosen name `Type`, `Func`,
  `Class` or `Interface`, and none of the four lists holds those. So no target
  meets a reserved word through the chosen name of a source declaration.
- **Field names.** No target lowercases a field name. TypeScript, Rust and
  C++ keep the horos spelling
  (`transform/syntax/horos-to-typescript/mapping/record/go/record.go:55`,
  `transform/syntax/horos-to-rust/mapping/record/go/record.go:56`,
  `transform/syntax/horos-to-cpp/mapping/record/go/record.go:56`), and Go
  capitalizes it (`transform/syntax/horos-to-golang/mapping/record/go/record.go:62`).
  A member's name is the pair's mapping's, and a reserved word in it is that
  language's to refuse (R13).
- **Capturing a predeclared or global name** is not a reserved-word question.
  A declaration named `string` in Go is legal and rebinds every use of the
  predeclared type. That is checked through recorded destinations,
  [#276](https://github.com/Bitspark/bittype/issues/276).

## Go

In Go the case of the first letter of a package-level name is its visibility.
Row 13's action: "It follows from the accepted decision: Files realizes
visibility under the language's rules, and in Go that is the case of the
name."

**Go's rule is Go's own, and it is a check (R8).** "a spelling is exported when
it begins with an upper-case letter, by Go's definition. For a declaration
that requires public access, Go's rule checks the spelled intent and refuses
it as unsupported coverage when it is not exported. It renames nothing."

- **What requires public access.** A requested root, and a declaration whose
  contribution states that it is part of a public representation (D5, R11).
  The optional helper is one.
- **A source declaration.** Its spelling comes from the mapping. That the
  mapping writes a capital is its style for every target (R8). Go's rule does
  not rest on it. The rule reads the spelling it is given.
- **A helper.** Its name is the one its recipe states. A public helper whose
  recipe states a name that is not exported is refused, for example a helper
  named `optional` (R8 and R14c, in the validation below).
- **Nothing is renamed.** In the first slices the realized name of an admitted
  Go declaration is its spelled intent.

**What Go Files refuses in the first slices.**

- A public declaration whose spelled intent is not exported (R8).
- An output with an unexposed source declaration in the public representation
  of an exposed one (D5, E4). In the structural profile that is every
  unexposed source declaration (E4).
- **Derived here:** a declaration that a contribution introduces and for
  which it states no public access. Files does not widen access (D5), and no
  unexported spelling is built (E4). No such declaration exists today.

Each is unsupported coverage, and each names what it refuses (R14).

**Deferred: the unexported spelling.** This page specifies no unexported form,
and nothing in the first slices builds one.

- **D4, of 12 November 2026, defined it:** "Go's unexported spelling is the
  lower-case form, with refusals where that form is a predeclared identifier, a
  keyword or a collision. No escape rule: a refusal can be lifted later, a
  spelling cannot be changed."
- **E4, of 13 November 2026, defers it:** "So Go's unexported spelling is
  realized by nothing in the first slices. Decision 4 above defined it. It is
  now deferred, with a trigger: a representation exists in which an unexposed
  source declaration is not part of an exposed one's public representation, or
  a request can state access beyond its roots (row 14a). Until then it is not
  built and not specified further."
- **Where the deferral is recorded.** Row 13a of the table "Applied" of
  [research 0016](../research-docs/0016-native-files-placement.md).
- **How to see whether the trigger has occurred.** The first half is about a
  Go representation.
  `grep -n "naming.Exported(f.Name)" transform/syntax/horos-to-golang/mapping/record/go/record.go`
  shows the one Go representation with members, which makes every member
  accessible. The second half has row 14a's command,
  `grep -n "Roots" contracts/generation/go/planning.go`.
- **Texts that still assume an unexported spelling.** They wait for the
  trigger. That is derived here from E4.
  - #274: done-when 9 has the Go rule return "the unexported spelling or a
    refusal", and two tests of its Verify table are about an unexported
    spelling.
  - #161, slice 1, item 6: "exported for an exposed declaration, unexported
    for every other". Its test list expects that the compiler's entity for
    `Hidden` "has an unexported name". It also says of the helper that "it
    becomes unexported", which D5 decides otherwise.
  - The Go naming instance of bittype/ACX-004 points to slice 1 of #161 for
    the realization, and says so.

**Fields.** A Go field's case is a member's accessibility and stays with the Go
record mapping (`docs/source.md:77-78`).

**Type parameters.** A type parameter is a binder inside one declaration. It
has no exposure and keeps its chosen name.

**The misplaced decision.** Go generation decides the case of every name
today, before exposure is known, and no Go-owned rule exists:

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
or not, and nothing of Go's checks whether that is right.
[bittype/ACX-004](architecture-cases/0004-generation-file-syntax.md) records
it as a pending instance.

**How the correction is validated (R8 and R14c, the validation of row 13).**
Go Files refuses a public helper named `optional`, refuses an unexposed source
declaration in the public representation of an exposed one, and realizes the
optional helper exported. Red today: no Go-owned check exists. R8 wrote the
first of the three with a helper named `string`. R14c replaces the example by
"a helper name that only the exportedness check refuses, such as `optional`",
because "A helper named `string` could be refused by the capture check of row
25 alone." After the
correction the mapping still writes a capital, as its style. What decides
whether a Go declaration is published exported is Go's check and the access
that its owner stated.

**Where the Go rule lives.** In a Go-owned package that the Go target imports,
and later Go Files (#274, done-when 9). It holds the check. The model is
`transform/syntax/typescript/symbol/go`, which "owns the TypeScript public
type symbols shared by generation and Files". Go Files applies it in slice 1
of #161.

## What is fixed, and when

**E6.** "the record of intent and each language's rule are fixed before the
plan envelope becomes stable and before the first release (#3). The envelope
is unstable until the six vertical cases pass, and a plan entry written before
that may still change its recorded qualification."

The text of #274, and the first version of this page, drew the line earlier:
before [#207](https://github.com/Bitspark/bittype/issues/207) and
[#195](https://github.com/Bitspark/bittype/issues/195) write their first plan
entry. E6 calls that wrong and gives the reason: "The default policy for Go
and Rust can only change with their Files stages (#161, #159), which come
later." The first release is
[#3](https://github.com/Bitspark/bittype/issues/3). The six vertical cases are
those of the plan envelope
([`docs/representation-families.md:95-102`](representation-families.md#the-plan)).
The fourth is a Go case.

Fixed before the envelope is declared stable and before the first release
([AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10):
a first change establishes the final structure):

- **The record of intent:** the chosen name, the module's full pin, the
  policy's identity and the spelling the policy produced (R10).
- **The rule of TypeScript, Go and Rust:** what "exposed" requires in that
  language, Go's check, each target's reserved words with what it refuses,
  and which outputs it admits. Of the last, R14 says: "The set of admitted
  outputs only grows afterwards."
- **Each target's default policy (R12).** "Each change of a target's default
  policy lands before the plan envelope is declared stable and before the
  first release. #195 gains that condition: the envelope is not declared
  stable before slice 1 of #161 and of #159 have changed their defaults."
  That note is not on #195 yet (see
  [the notes not yet posted](#what-is-delegated-and-what-is-not-specified-here)).
  TypeScript's default changes with #273, Go's with slice 1 of #161, Rust's
  with slice 1 of #159 (E6).

**C++ is excepted.** R12: "C++ is excepted by row 36 of research 0016." It
keeps the qualified policy until #160, the guarantee of research 0016 excludes
legacy C++ output (row 36), and its exposure is deferred (row 32).

**A plan entry (R9).** That a plan entry records a declaration's identity and
naming intent, and no realized native name, is not fixed by this design. R9:
"It is what this design asks of #195, pending research 0017." Done-when 2 of
#274 asks only that #207 and #195 each carry a comment that says so.

- **Which identity.** R9: "Which identity a plan entry records is not defined
  here: an output-issued handle cannot be the key of a serialized plan."
- **Why a handle cannot be that key.** Derived here from three sentences. A
  plan entry exists at planning, before any output:
  "The planner selects one immutable representation instance per use. The
  native type and the plan entry are both projections of that instance."
  (`docs/representation-families.md:44-46`). A handle belongs to one sealed
  output: "A handle issued by one `Output` is refused when it is resolved
  against another `Output`. That includes a second run of the same plan."
  (#272, done-when 5). And the plan crosses a process boundary: "The plan
  starts as a typed build-time interface and is serialized only where a real
  process boundary needs it." (`docs/binding-generation.md:70`).

What may change afterwards, without changing a name that exists:

- **A refusal that is lifted**, by an explicit rule (R13): Rust's of `Self`,
  C++'s of `__`, Go's of a spelling that is not exported, Go's and Rust's of an
  output with an unexposed dependency.
- **What is deferred with a trigger:** Go's unexported spelling (row 13a) and
  access that a request states beyond its roots (row 14a).
- **A field added to the inventory record.** #272: "Adding a field to the
  record is additive."

## How a consumer learns the realized name and the placed site

**D6.** "A consumer learns the realized name and the placed site from a reader
on the formatter's result, keyed by the output's handles." The handles are
those of #272. R14 adds access, and that there is no earlier source: "A
consumer learns realized names and access from the formatter's result. Nothing
offers them before formatting."

- **Why a reader is needed.** `Source` reports neither today. The TypeScript
  formatter's map from an original site to its placed site is private
  (`format/files/typescript/provider/go/relocate.go:28`).
- **Why on the formatter's result.** A name is realized only once a language's
  Files stage has admitted it, and the placed site exists only after
  placement. The alternative was to read both from the published semantic
  graph (#274, point 6, option b).
- **What a binding backend reads** is with #195 (R9, E7), under
  [Reconciliation](#reconciliation-with-accepted-documents).

## Reconciliation with accepted documents

**[Representation families](representation-families.md#families), lines 207
to 209. Derived here.** "A lifted inner binder gets a deterministic name, with
collision and reserved-word handling, and keeps its free generic parameters."
The sentence does not say who handles a collision or how. The following
reading is this page's, from rows 8 and 9, D1 and R13. It is not a decision of
the owner of that family. R13: "The page marks the reconciliation with the
lifted binder of `docs/representation-families.md` as derived, and #164
carries the note." That note is not on #164 yet (see
[the notes not yet posted](#what-is-delegated-and-what-is-not-specified-here)).

| Part of that sentence | Where this reading puts it |
| --- | --- |
| "a deterministic name" | Generation. The recursion family states a deterministic **intent** for the lifted binder. It depends on the declaration that the binder is lifted from and on the binder, never on which other declarations are in the output |
| Collisions among the names of one source declaration: its type parameters, its members, and the binders lifted from it | Generation. Their meaning is wholly inside one declaration |
| Collisions between top-level declarations, which a lifted binder becomes | The language's Files stage, in the placed scope. In the first slices it refuses (row 8, R13, and "never a hash", row 11) |
| "reserved-word handling" | The target (D1). In the first slices it refuses (R13) |

- **What the reading gives up.** The advice of research 0015 behind that
  sentence derives the name "from its identity and native scope". Generation
  does not know the placed scope under research 0016. So generation cannot
  give a lifted binder another name when its name meets that of another
  declaration of the output.
- **The access of a lifted binder** is not part of this reading. Its
  contribution states it when that family is built (R11).

**[Binding generation](binding-generation.md#the-selected-representation-plan):
a tension that this page states and does not resolve (E7, R9).** The plan
carries "native symbols, distinguishing run-time imports from type-only
imports" (`docs/binding-generation.md:56`), and the shared envelope of
[representation families](representation-families.md#the-plan) lists "Native
symbols, imports, factory references, runtime requirements". The completeness
test says: "The backend can implement each selected representation from the
plan, the source semantics and the selected contributions alone, without
guessing a rule." (`docs/binding-generation.md:72`).

- **The tension.** E7: "With names realized by Files, a plan entry cannot hold
  a realized name or an import path."
- **What this design asks of #195.** E7: "A plan entry refers to a native
  declaration by its identity, and the backend reads the realized name and
  place from the formatter's result. That is a fourth input to the backend and
  an amendment to an accepted sentence."
- **Who decides.** E7: "It belongs to #195 and goes into research 0017 as part
  of its question on the identities of the plan." This page amends no sentence
  of binding generation, and that page stands as written until
  [#195](https://github.com/Bitspark/bittype/issues/195) decides.

**[Semantic generation](semantic-generation.md#representation-and-enforcement-profiles),
lines 47 to 49. Derived here.** "The version stays `/1`: nothing is released
and the generated output did not change with the names." That sentence is
about the renaming of the two profile identities, and it gives two reasons
together. A change of a target's default policy keeps the first and not the
second: it lands before the first release (R12), and it changes the names of
declarations of imported modules for a consumer that passes no policy. The
profile's version does not record that change, because a policy is not a
profile (D3) and a profile is "a promise about the native type"
(`docs/semantic-generation.md:51`). So the sentence is not a rule that
unchanged output is a condition of `/1`.

**[Source and Files](source.md#from-generated-declarations-to-files).** No
sentence of the decision of 7 November says who spells a declaration's name.
The statement that Files never renames a declaration was a derivation of the
research document (its statement 18) and is withdrawn (row 13). In the first
slices no Files stage renames one all the same (R8, R13). The one place where
this design adds to the decision is the second owner of required access,
described under [Exposure](#exposure-and-its-audience).

## What is delegated, and what is not specified here

The first version of this page listed seven open points, and E1 to E7 decide
them. A second version said that nothing is open. R14 replaces that sentence
with this list.

**Decided by R14.** Points 1, 2 and 4 in its words. Point 3 is withdrawn by
R14c:

| | Decision | Stated under |
| --- | --- | --- |
| 1 | "Go's exact rule for what is in a public representation is written in slice 1 of #161. The set of admitted outputs only grows afterwards." | [Exposure](#exposure-and-its-audience) |
| 2 | "A refusal names every offending declaration, each by its handle and its spelled intent." | [Exposure](#exposure-and-its-audience) |
| 3 | Withdrawn. It made TypeScript's refusal of a separating placement unsupported coverage. R14c keeps it invalid input, as #271 decided, and states the rule for the category of a refusal | [Exposure](#exposure-and-its-audience) |
| 4 | "A consumer learns realized names and access from the formatter's result. Nothing offers them before formatting." | [How a consumer learns](#how-a-consumer-learns-the-realized-name-and-the-placed-site) |

**Delegated**, each with the one who decides:

| What | To | Source |
| --- | --- | --- |
| Whether a plan entry may omit native symbols, and which identity it records | #195 and research 0017 | E7, R9 |
| That the envelope is not declared stable before the defaults of Go and Rust have changed | #195 | R12 |
| The access of a lifted inner binder, and the reading of its name | [#164](https://github.com/Bitspark/bittype/issues/164) | R11, R13 |
| Go's exact rule for what is in a public representation | Slice 1 of #161 | R14 |
| Rust's rule for a usable public record | [#278](https://github.com/Bitspark/bittype/issues/278) and slice 1 of #159 | Row 24 |
| C++ exposure, and the C++ default | #160 | Rows 32 and 36 |

**The notes not yet posted.** R12 says that #195 gains a condition, and R13
that #164 carries a note. Neither issue has the note yet. R14c: "The notes for
#164 and #195 are not posted yet. I post them when this page merges."

**Not specified here:**

- **Go's unexported spelling.** Deferred with a trigger (row 13a).
- **Access that a request states beyond its roots.** Deferred with a trigger
  (row 14a).
- **How the four parts of an intent reach the inventory record.** R10 fixes
  what is recorded and that the mapping spells. The spelling travels as the
  declaration's name. The carrier of the other three parts, from the mapping
  to the record that the target issues, is left to the delivery of #272 and
  of done-when 8 of #274.
- **The identities of the two policies.** R10 makes each a named value. This
  page gives neither a name.
- **What a helper's intent records in place of a module and a policy.** A
  helper has no source module, and the policy is given to the declarations
  mapping, not to the mapping that states a helper. Its spelling is the name
  its recipe states.
- **An escape rule.** None exists. One is "a later, explicit rule" (R13).

**Where each decision of 13 November is stated:**

| | Decision | Stated under |
| --- | --- | --- |
| E1 | Rust refuses `Self`. The other targets keep it | [Reserved words](#reserved-words) |
| E2 | A qualification is recorded as the module's identity | [Three things that were one](#three-things-that-were-one) |
| E3 | `Type` before a leading underscore stays with the chosen name. The refusal of `__` moves to the C++ target | [Three things that were one](#three-things-that-were-one), [Reserved words](#reserved-words) |
| E4 | Which outputs Go and Rust admit. Go's unexported spelling is deferred | [Exposure](#exposure-and-its-audience), [Go](#go) |
| E5 | The rule is each language's. TypeScript admits within one module | [Exposure](#exposure-and-its-audience) |
| E6 | Fixed before the envelope is stable and before the first release | [What is fixed, and when](#what-is-fixed-and-when) |
| E7 | "Native symbols" in the plan: stated and not resolved | [Reconciliation](#reconciliation-with-accepted-documents) |
| R8 | Go's rule is a check and renames nothing | [Go](#go) |
| R9 | What this design asks of #195 about a plan entry | [What is fixed, and when](#what-is-fixed-and-when) |
| R10 | The mapping spells. What an intent records | [Three things that were one](#three-things-that-were-one) |
| R11 | A contribution states the access of every declaration it introduces | [Exposure](#exposure-and-its-audience) |
| R12 | Each default change lands before the envelope is declared stable | [What is fixed, and when](#what-is-fixed-and-when) |
| R13 | In the first slices a target refuses and does not rename | [Reserved words](#reserved-words) |
| R14 | Points 1, 2 and 4 above. Point 3 is withdrawn | This section |
| R15 | Annotations of the research table | Rows 13, 14, 14a, 23 and 29 of the table "Applied" |
| R14c | TypeScript's refusal of a separating placement stays invalid input. The rule for the category of a refusal. Four small corrections | [Exposure](#exposure-and-its-audience), [Go](#go), [Three things that were one](#three-things-that-were-one), this section |

## Delivery order per target

From #274, "Order", as E6 and R12 correct it, and rows 8, 10, 11 and 12 of the
table "Applied".

| Step | Needs | Item |
| --- | --- | --- |
| This design | Nothing | #274, done-when 1 to 4 |
| The inventory with handles. Each record can carry an intent, and a contribution's stated access | Nothing. It can start with [#271](https://github.com/Bitspark/bittype/issues/271) | #272 (row 10, and its decision 4) |
| The builder's comparator orders by name, then origin, then planned identity. The name it reads is the spelled intent (R10) | In the same change as the first removal of the flat-name rule | #272 (row 12) |
| The policy is a named value that the mapping applies. The qualified policy reproduces today's output in all four targets | Nothing | #274, done-when 5 and 6 |
| The shared naming function holds no reserved word and no refusal of `__`. Rust refuses `Self`, and C++ refuses a name that contains `__` | Nothing | #274, done-when 7 as E1 and E3 change it, and #268 |
| Every inventory record carries an intent for TypeScript, Go and Rust | #272 | #274, done-when 8 |
| A Go-owned package holds Go's check | #274 names no prerequisite | #274, done-when 9 as R8 changes it, without the unexported spelling (row 13a) |

Then per target (row 8: "TypeScript first, Go and Rust with their slice 1, C++
with #160"):

| Target | The flat-name rule leaves, and Files checks collisions per placed scope | What Files does with a spelled intent | The default becomes unqualified |
| --- | --- | --- | --- |
| TypeScript | #273, second change, after #272 | TypeScript has no case rule. Its Files stage already writes the export from exposure (`relocate.go:89-91`). Its refusal of a separating placement stays invalid input (R14c) | After #273 (#274, done-when 10) |
| Go | Slice 1 of #161 | Slice 1 of #161: the check that a public declaration's spelling is exported, and the refusal of an output with an unexposed dependency. Nothing is renamed. The unexported spelling is deferred (row 13a) | Slice 1 of #161 |
| Rust | Slice 1 of #159 | Slice 1 of #159 realizes visibility and refuses an output with an unexposed dependency. Rust has no case rule | Slice 1 of #159 |
| C++ | With #160. Until then C++ selects the flat rule in its own target | Deferred until [#167](https://github.com/Bitspark/bittype/issues/167) supplies a C++ semantic provider | With #160. Until then the qualified policy |

The record of intent and the rules of TypeScript, Go and Rust are fixed, and
each of their defaults has changed, before the plan envelope is declared
stable and before the first release (E6, R12). #274's "Order" said that all of
it lands before #207 and #195 record a native name, and E6 corrects that.

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
| A type parameter named like any declaration of the output is refused | `transform/syntax/shared/target/go/check.go:25` |
| TypeScript Files reads a requested root as exported and writes the export | `format/files/typescript/provider/go/generated.go:92-97`, `format/files/typescript/provider/go/relocate.go:89-91` |
| TypeScript Files refuses a use, from another file, of a declaration that is not exposed, as invalid input | `relocate.go:163-165` |
| TypeScript Files compares the realized exposure with the requested one in both directions | `format/files/typescript/provider/go/format.go:201-205` |
| TypeScript Files refuses two declarations of one name in one file | `relocate.go:120-126` |
| TypeScript Files spells a use in the same file from its destination declaration | `relocate.go:155-160` |
