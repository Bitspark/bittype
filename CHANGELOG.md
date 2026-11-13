# Changelog

## The Go bridge publishes the package name, and the Go environment states where the output root sits — 13 November 2026

A Go Files stage will decide a package's name and, through the directory, its identity. The Go semantic profile now gives it something to check each with. This is [#277](https://github.com/Bitspark/bittype/issues/277), which delivers rows 3, 21 and 22 of [research 0016](research-docs/0016-native-files-placement.md). No Files stage is added, and nothing that generation produces changes. The rule for the file names and package names that a slice admits (row 31) is the Go Files formatter's and comes with slice 1 of #161.

- **The package name is a fact.** Every `package` entity, captured and supplied, carries one `package-name` fact with `types.Package.Name()`. The entity's own name stays the path, and `LookupIn` keeps reading it. Two captures at one path that differ only in `package pages` and `package people` had equal graphs. They now differ in the one term of that fact and in nothing else. The base profile's graphs gain the fact too. Its profile strings and base configuration are unchanged: no accepted document asks for a new identity when a graph gains a fact.
- **The form of the fact, and its readers.** The term is a provider expression with the operator `go/package-name`, the encoding `go/utf8`, no operand and the name as UTF-8 text. The Go vocabulary publishes `PackageNameRelation()`, `PackageNameOperator()` and `TextEncoding()`, and the bridge writes from them. `Model.PackageName(importPath)` reads the fact. `Type.Package()` reads the package of a package-level declaration from the graph's ownership: identity and name. Both fail when the package does not carry exactly one fact in that form, and neither reads the path in place of the name.
- **The output root is stated.** `ImportEnvironment` gains `OutputRoot`, the directory of the output root relative to the module root. `Module` now means the module path and nothing else. `pages/types.go` under `example.org/app` with the output root `generated` belongs to `example.org/app/generated/pages`. `CapturedImportPath(environment, filePath)` derives that identity. The bridge derives each captured file's package once with it and groups by the result. It replaces `CapturedImportPath(module, filePath)`.
- **Nothing is defaulted.** An environment without an output root is refused with `ErrInvalidRequest`. The item decided that. It did not say how the module root is stated: it is `golang.ModuleRoot`, which is `.`, the relative spelling of that directory. A caller writes it, so an empty field is never read as it, and it has no second spelling.
- **Which roots are refused, beyond the item.** A root with an empty, `.` or `..` element is not a clean path and is an invalid request: the empty string, an absolute path and a trailing slash among them. An element follows Go's rule for an element of an import path, as `golang.org/x/mod/module` states it in the pinned go1.27.0 toolchain (`module.go:417-482`): `con`, `NUL`, `com1`, `lpt9.x`, `gen~1` and `~1` are refused, and `console`, `x.con`, `gen~x` are admitted. The rule is one function. It is narrower than Go's for a leading dot and a plus sign. `ValidImportPath` does not apply the device names and the tilde form to module, supplied and captured paths yet, which is [#287](https://github.com/Bitspark/bittype/issues/287). A root whose path with the module path exceeds 512 bytes is refused, so a root's validity depends on the length of the module path.
- **One rule for the directories the profile does not cover.** `UnsupportedDirectory` in `semantics/golang/go` is the one declaration: `testdata`, `vendor`, and a name that starts with `_` or `.`. The go tool leaves `testdata` and a leading underscore or dot out of `./...` and still builds them on an explicit import, and it refuses an import through `vendor`. The bridge spelled the rule for a captured file's path and now calls the one declaration. What it refuses, and with which error, is unchanged. The same directory as part of the output root is refused by `ImportConfiguration` with `ErrUnsupported`, the category the bridge gives, and no setting is written. Without that refusal the output root `vendor` would put every captured file at a vendor path while the bridge, which reads the Files path, saw none. `CapturedImportPath` gives no identity for a file in such a directory.
- **The setting's bytes change.** The key is `imports-2`, and the value is `{"Module":…,"OutputRoot":…,"Packages":[…]}`. This setting is part of every receipt under an import environment. A setting under the old key `imports-1` is refused, and so is a value under the new key that states no output root. `TestGoImportEnvironmentMustStateTheOutputRoot` pins the key, three literal values, one of them with a supplied package, and the refusal of eight hand-written settings. A consumer that built an environment adds `OutputRoot: golang.ModuleRoot` to keep what it had.
- **A supplied package inside the module.** A supplied path at or below the output root stays refused. A package of the module outside the output root is hand-written code and may be supplied, such as `example.org/app/model` beside the output root `generated`. Before, every path under the module was refused, so generated code could not import a hand-written package of its own module.
- **The limit, in the document.** [Native semantic providers](docs/native-semantics.md#import-environment) states the environment, the rules for the output root, the bytes of the setting, the name fact, its form and the readers, and extends the classification table. It also states what the comparison of a declaration's owner with the derived identity of a file establishes: the package of that file's directory. Which file holds the declaration comes from its site. Nothing certifies where the consumer writes the files, because the type checker is told each package's path.
- **A record brought up to date.** [Package boundaries](docs/package-boundaries.md) said that no production Go code compares a term's encoding with `utf8`. The name readers now do, so the Go vocabulary publishes the key and the bridge writes it from there at both places.
- **Tests.** Seven tests in `integration/go/go_imports_semantics_test.go` and four in `semantics/golang/go` are new. Thirteen existing lines that construct an environment gain `OutputRoot: n.ModuleRoot`, and no existing assertion changes. 85 mutants, each removing one check or restoring one old behaviour, each fail a test. A mutant of each part of the directory rule fails at both entrances, the captured file's path and the output root. The rows of the element rule were also put to the go tool itself: it accepts every admitted row.
- **Unchanged output.** No generation or formatting package imports `semantics/golang/go` or `internal/golangsem`. The lists of `go test ./...` before and after differ only in the eleven new tests.

Not done:

- The rule for admitted file names and package names, with the Go formatter's admission and its guard. It is slice 1 of #161, and the two functions that an earlier state of this change held are saved for it.
- `ValidImportPath` for module, supplied and captured paths: #287.
- The derivation judges a file's directories and not its name. The bridge refuses test files and files that start with `_` or `.` itself, as before.
- The generated reports `docs/language-support.md` and `docs/provider-capabilities.md` are not edited here. The workflow on `main` regenerates them with the seven new cases.

## The contribution contract and the naming design say that their notes are posted — 13 November 2026

Both design pages listed notes to other issues as "not yet posted", to be posted when the page merged. Both have merged and every note is posted, so the pages say so and link each note. Documents only.

- **[The contribution contract](docs/contribution-contract.md)** merged as `2a29674`. Posted: the reworded text of #207 with a correction about vertical case 1, the same correction of the grade on #196, the seven lines of the gate in the text of #196, and the request about the memo sentence on bitspark-service#40.
- **[Naming](docs/naming.md)** merged as `7e9dc34`. Posted on 12 November 2026: the condition on #195 and the note on #164.

## The contribution contract: a design with its gate, from research 0017 — 13 November 2026

[The contribution contract](docs/contribution-contract.md) is the design of a contribution to the generation kernel: what it is, what its rule may ask and must return, what identifies it, and what the kernel carries for it unread. It is step 1 of the order that [research 0017](research-docs/0017-contribution-contract.md) decided for [#196](https://github.com/Bitspark/bittype/issues/196), item 6. The research owner decided it, and the page says so. It is not implemented, and it becomes final only when its gate is met. This entry adds and changes documents only. The mapping contract is unchanged.

- **Every row has its place.** The page states each of the 49 rows of the research table with its source, quotes the accepted sentence a row rests on with file and line, and marks what it derives itself. A table names the section and the gate line of each row.
- **The protocol.** A rule gets one round: it names its dependencies and the evidence it needs, the planner resolves them, the rule decides with typed answers, the decision is frozen when the rule returns, and a registered projector that reads no source computes the native intent. There is no general loop.
- **Evidence** is bounded by the coherence that a profile requires of a generic declaration and its applications. An operand bound to a binder is answered from the binder's requirements, never from one application. An answer is read from what the dependency's selected contribution stated, and no provider computes it.
- **Refusals** refine the seven values and four verdicts of the generation contract and replace none. A refusal of a selected exact mapping never falls back.
- **Identities.** One list names each identity of the design with its issuer, whether it leaves the process, and what it must not be confused with. A contribution reference names the family's meaning, which is the same in every target, and the contribution's own representation revision.
- **Optional fields are conditional.** The section states the position of bittype and what happens under each answer of the service library on bitspark-service#38. The round of requests is in the contract under either answer. A present `undefined` in a TypeScript optional field is refused: two accepted sentences settle it, and the question first put to the service library is corrected there.
- **The gate.** Seven lines, each a predicate with the tests that are red today, and one condition beside them: the question of optional fields is decided. The migration of the 21 mapping packages does not merge while one is red.
- **Twenty-six decisions after the first version.** The first version of the page stopped at ten points on which two rows, or a row and an accepted page, did not fit without a reading. The research owner decided them as [O1 to O10](https://github.com/Bitspark/bittype/issues/196#issuecomment-6841190077) on #196, with [O11](https://github.com/Bitspark/bittype/issues/196#issuecomment-6841518102) for one word of O8. An independent reader then read the page and reported 25 findings, and [O12 to O20](https://github.com/Bitspark/bittype/issues/196#issuecomment-6841866946) decide the twelve that needed a decision. [O21 and O22](https://github.com/Bitspark/bittype/issues/196#issuecomment-6842186587) correct O12 and O19 from the accepted contract of a family. The same reader read the changes and reported 13 more findings, and [O23 to O26](https://github.com/Bitspark/bittype/issues/196#issuecomment-6842416622) decide the four in which an accepted page says what the page or a decision had left open or said differently: a fact of evidence is a typed value that the kernel carries unread, the family's meaning is the application's promised domain, the contribution's identity and its generated payload types live in a vocabulary package, and missing and ambiguous are two refusal values. The page states each decision with its reason, and nothing is open inside it. A generation request names its products: types alone, or types and bindings. A rule returns a value of a generated plain Go type, and a payload codec encodes, validates and freezes it when the rule returns. A bundle is a contribution and the facet packages registered for it by the consumer's composition, and the facet packages of the builtin contributions are test code. The delivered contributions, 21 packages under six names, return a decision, the declaration contribution among them. Vertical case 1 is not in the prototype, and recursive references become final with #164, not with the gate.
- **The research tables.** In the table "Applied" of [research 0017](research-docs/0017-contribution-contract.md), the Action of rows 8, 17 and 31b is corrected as O2, O4 and O6 decide, each dated and with its former text kept, and sixteen rows carry a dated note for the decisions that bear on them. In the table of [research 0015](research-docs/0015-representation-families.md), the validation of row 2 is reworded as row 39 of research 0017 decides, and keeps its former text. The submitted documents and the advice are unchanged.
- **The verification reports.** The reports of the three verification agents of research 0017 are committed beside the research document as `.verify-1.md`, `.verify-2.md` and `.verify-3.md`, each with a header that says it is a proposal and no decision. The page cites them about twenty times, and they were on a workstation only.
- **One amendment and dated clarifications of accepted pages.** No accepted sentence is rewritten. One is amended in substance, in an insertion headed "Amended": [binding generation](docs/binding-generation.md) said that a backend works from three inputs alone, and it has a fourth, the placed result. That is the decision of #195 that the [naming design](docs/naming.md) waited for, and that page gets a dated note. The other insertions clarify. Binding generation: "library × target" names ownership and distribution, and two items of the plan's list are read with the amendment: a rule is carried as a field or through the contribution reference, and "native symbols" means an identity and a spelled intent. [Representation families](docs/representation-families.md): what "cannot disagree" claims, the two senses of "instance", what is unstable in the envelope, what "native symbols" in the envelope means, what "change once" commits to, and that generics are no family with a contribution of its own. Finding 7 of [package boundaries](docs/package-boundaries.md) and [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) say, in two places each, that the optional move is conditional.
- **Pointers.** The [SPI](docs/generation-spi.md) and the [registration page](docs/kernel-composition.md) say that the mapping contract is going to change and where the design is. [Semantic generation](docs/semantic-generation.md), `LAYOUT.md` and the [documentation map](docs/README.md) link the page.

Not done: no code and no prototype. `LAYOUT.md` names no path for the new payload codec role and has no role for a facet package of a mapping contribution or for a vocabulary package of a family for a target. It gets them with the prototype's report. Four notes that the page lists as not yet posted, to be posted when it merges: a correction on #207 and in its text, a correction of the grade comment on #196 about vertical case 1, the wording of the binding contracts' memo sentence, to bitspark-service#40, and the seven lines of the gate in the text of #196.

## A sealed output lists every declaration it owns and issues the handles a placement is built from — 13 November 2026

`generation.Output` has a complete inventory of its declarations, and each record carries a handle that belongs to that one output. This is [#272](https://github.com/Bitspark/bittype/issues/272) and delivers rows 10 and 12, the inventory part of row 19 and the order basis of row 5 of [research 0016](research-docs/0016-native-files-placement.md#applied). The generated output of every existing test is unchanged: a temporary probe recorded 144 records in a full run before and after, 88 formatted texts, 38 output trees, 8 placed TypeScript trees and 10 Files records, and each has the same SHA-256. The placement types, the shared admission and the removal of the flat-name rule are #273.

- **The inventory.** `Output.Inventory()` returns one record for each declaration of the snapshot and no other record, in the order of the sites. A record carries a handle, the site, the role, the provenance and the values that the owner of the description issued. The roles are requested, dependency and helper. They are not a closed set: a further role may be added, and a reader refuses a role it does not know. The provenance is the list of the declaring sites of the source declarations that the declaration represents, and it is empty for a helper. There is no grouping field.
- **Where the facts come from.** The planner knew the source declaration of each planned declaration and which of them are helpers, and `PlannedDeclaration` had no field for either. It now has `Role` and `Provenance`. `Role` replaces `Exposed`, which stated one of the three roles. The shared constructor held the node of every declaration in a local map and returned only the requested ones. It now returns the site of every planned declaration in `Assembly.Declarations`, once for all four targets. `Assembly.Exposed` is removed: the kernel derives the exposed sites from the plan and the listed sites, so the requested records are exactly the sites of `Exposed()` by construction. The planner hands a target a copy of each declaration it asks the target to check or to build, so a target that overwrites what it was handed changes no record.
- **The handle.** A `DeclarationHandle` carries its output, as an `ExternalSymbol` carries its issuer. `Output.Resolve` returns the record of a handle that this output issued and fails with `ErrScope` for any other. A second run of the same plan has the same tree identity and the same node numbers, and it refuses the handles of the first run.
- **Issued values.** A naming intent and the access a helper requires are not meanings of the kernel. Their owner creates a `ValueIssuer`, a declaration recipe or a helper recipe states the issued value in `Values`, and the kernel carries it to the record. Only the issuer reads the content through this API. The planner refuses a value that no issuer issued with `ErrInvalid` and charges each value to the output budget, the bytes of its content and one byte for the value. Limit: the content is hidden by the API only, as an external symbol's key is. Reflection reads it, a value compares equal wherever it is carried, and whoever holds a value can state it again in another description. Formatting a value, or a record or a description that holds one in an exported field, prints whether it was issued and never its content; formatting an `Output`, a `Program` or another struct that holds a value under an unexported field prints the content. No builtin mapping states a value yet. The naming intent is #274.
- **What sealing checks.** `seal` fails with `ErrIncomplete` when the assembly lists another number of declarations than the plan has, and with `ErrScope` when a site lies outside the snapshot, is listed twice or lies inside another listed site. The last check needs no knowledge of a language: it refuses the root node beside any other declaration, and a node inside another listed declaration. The shared constructor guarantees that a site is the node the target's lowering returned for that declaration, emitted once. Where that node sits is the lowering's decision: C++ declarations are the children of a namespace. The constructor returns no assembly when a planned declaration has no node in the unit.
- **The comparator.** The shared constructor seeds its walk in the order of `compareDeclarations`. The walk emits each declaration after the declarations it uses, so the comparator fixes the node numbers and the sites. Its key is the name, then the provenance, then the planned identity. The provenance is the planner's own fact. The origin on a description is the mapping's optional statement, so it is not part of the key. Before, the comparator returned 1 for two equal names in both directions. That could not be reached, because the constructor refuses equal names first, and it still does. Limit: two declarations of one name and one provenance, or two helpers of one name, are ordered by the planned identity, which follows the order of the requested roots.
- **The shared contract changed**, and platform's rule AC-08 asks why. The capability is to address every owned declaration without reading native syntax, and to bind a placement to the output it was written for. The generation kernel is the layer, because it issues declaration identity: only the planner knows the roles and the source declarations, and only the kernel seals the capture that a site belongs to. Every language's Files stage uses the same records and handles, and a mapping or a target language adds a meaning by issuing a value without a change to the kernel. The boundary change is not only an addition: a target that fills its own `Assembly` now lists every declaration. The four targets build through the shared constructor and are not edited.
- **Tests.** Twenty tests are new, with 36 subtests. The independent fixture generates three roots that need one dependency and one helper, and the test selects the five declarations by role, provenance and issued value and reads no `name` attribute of the output. One plan is run twice and each run refuses the handles of the other. A fixture target that omits, adds or repeats a site, or lists one outside the snapshot or inside another listed site, is not sealed. A target that overwrites every value and every parameter name of the declarations it is asked to check changes neither the records nor the generated unit, for a source declaration and for a helper. The pagination module gives a complete inventory in TypeScript, Go, Rust and C++, and four independent records requested in descending order of their names get ascending sites in all four. An assembly of two trees gives its records in the order of the trees. The comparator is checked for every pair of ten declarations, whose origins and planned identities run against the order wherever an earlier part of the key decides. 55 seeded faults each fail a test: four restore the behaviour before this change, 48 remove or change one check, one copy or one part of the key, and three are the faults that an independent review found surviving the first version of these tests. The two existing tests of `seal` build the assembly and the plan in the new form and assert what they asserted. The independent fixture's library states a naming intent for each declaration it describes and an access for its helper, which changes no generated text.
- **Documents.** [Planning and running type generation](docs/generation-planner.md#the-declaration-inventory-and-its-handles) describes the inventory, the handle, the issued values and the order of sites, and records the justification and what a later change does not alter. [The generation SPI](docs/generation-spi.md) and [generation](docs/generation.md) follow.

What a later change does not alter: the accessor `Inventory` and that it is complete, the fields of a record and their meaning, the meaning of each of the three roles, the handle type and that it identifies one declaration of one sealed output, and the comparator's key. Adding a field to the record is additive, and so is a further role.

Not done:

- The kernel cannot tell which declaration a listed site realizes. A target that fills its own assembly and lists the sites of two declarations in each other's place is not detected. A test compares the sites with the declaration nodes of each of the four targets.
- A single listed site that is the root of its tree is not refused, because no other site is inside it.
- The TypeScript Files formatter and the placement tests still find declarations by walking the tree. They move to the inventory with #273.
- No target reads an issued value, and a target cannot add one when it builds. A value enters through a description.
- One source declaration still yields one native declaration. The record allows several, and the planner does not produce them.

## Research 0017: the contribution contract, evaluated and in delivery — 13 November 2026

[Research 0017](research-docs/0017-contribution-contract.md) asks what the mapping contract must become when a rule also has to yield what a binding backend needs: what a rule returns, how one contribution asks for the answer of another, what the kernel carries without reading it, and what identifies a contribution and a plan. It came from a reading of the accepted designs for the first slice (#207, #195, #196, #197) that left three points open. Four readers reviewed it, and it was answered through the consult service. This entry adds documents only.

- **Evaluated.** The advice was split into 40 recommendations, each checked against the code at `7e9dc34` and the state of the service library. 18 rows hold as written, 21 are accepted with a named adjustment, nine are deferred with a trigger and one is stale. None is rejected. The table is the section "Applied" of the document.
- **What it decides.** A rule gets one round of requests before it decides, and a registered projector that reads no source computes the native type from the frozen decision. Evidence about an operand is bounded by the coherence of generic declarations. The identity of a contribution is its publisher and a local key, with a semantic revision. A plan has use entries and representation nodes, and open generics are primary. A facet is a typed, validated value on a contribution. A backend gets the placed result as a fourth input.
- **What it withdraws.** That the contract could be final while nothing reads its new half. The node key and the identity form the document leaned to. "Member" as a sufficient answer for an optional field.
- **Optional fields.** The position of bittype is that the standard record representation delegates presence through a published protocol. That needs a field-level conversion slot in the service library, so it is asked on bitspark-service#38 and built in the prototype first.
- **The order.** The design page, a prototype on a branch that is thrown away, a materially different contribution, then the contract in its final form and the migration of the 21 mapping packages on `main` in one change. "Changes once" means one migration to the demonstrated extension contract for the agreed extension domain, before independently owned consumers rely on it.

## The naming design: generation records an intent, each language realizes the spelling — 13 November 2026

[Naming](docs/naming.md) is the accepted design for who decides a native name and its visibility. It comes from rows 11, 13 and 14 of [research 0016](research-docs/0016-native-files-placement.md) and from the comments of its owner on [#274](https://github.com/Bitspark/bittype/issues/274): the six decisions of 12 November 2026, the [seven of 13 November 2026](https://github.com/Bitspark/bittype/issues/274#issuecomment-6838436844) on the points that the first version of the page left open, the [eight after an independent reading](https://github.com/Bitspark/bittype/issues/274#issuecomment-6839006959) of the page, and [one correction](https://github.com/Bitspark/bittype/issues/274#issuecomment-6839289646) after its second reading. It answers [#268](https://github.com/Bitspark/bittype/issues/268). This entry adds and changes documents only. No target has changed a name, and the page says so at its head.

- **Three things that were one string.** A declaration's identity, which alone is authoritative for references. The naming intent, which generation records: the mapping's chosen name, the module's identity as its full pin, the identity of the naming policy, and the spelling that the policy produced. The realized spelling, which a language's Files stage publishes.
- **Owners.** The mapping chooses the name, in a capitalised form that is the same for all targets, and spells it under the policy. The consumer's composition selects the policy, a named value that it gives to the declarations mapping of each pair. The unplaced tree's name is the spelled intent. Each language checks that spelling under its own rules: its reserved words, in Go the case of a name, and in its Files stage the collisions of a placed scope.
- **Refuse, not rename.** In the first slices a target refuses a spelling it cannot publish. An escape rule is a later, explicit rule. Rust refuses `Self` as unsupported coverage, and the other targets keep `Self`. The refusal of a name that contains `__` moves from the shared naming function to the C++ target.
- **Exposure.** `Exposed()` stays the list of requested roots. The page states what an exposed declaration must be in TypeScript, Go and Rust, and the rule that required access is established and realized access stays within authorized access. The request states the access of its roots, and a contribution states the access of every declaration it introduces. So the Go optional helper is realized exported. Files never widens access on its own.
- **Which outputs a language admits.** Go and Rust admit only outputs in which every source declaration in the public representation of an exposed one is itself a requested root. A refusal names every offending declaration. In the structural profile that covers every dependency. TypeScript admits an unexposed declaration in the body of an exposed one within one module, as `TestGeneratedFilesPreservePublicationIntent` pins. Its refusal of a placement that separates the two stays invalid input, as #271 decided. The rule for both: a refusal is unsupported coverage when a later slice can admit the same input, and invalid input when no later slice can.
- **Go.** Go's rule is its own, and it is a check: a spelling is exported when it begins with an upper-case letter, and Go Files refuses a declaration that requires public access and whose spelled intent is not exported. It renames nothing. Go's unexported spelling was defined by the decisions of 12 November 2026 and is deferred by those of 13 November 2026, because nothing in the first slices realizes it. The page specifies no unexported form.
- **The module-pin prefix** is today's default and becomes one of two policies. It becomes selectable first, with unchanged output. Unqualified names become a target's default once that target's Files stage owns collisions: TypeScript with #273, Go and Rust with slice 1 of #161 and #159.
- **What is fixed, and when.** The record of intent and the rules of TypeScript, Go and Rust, before the plan envelope is declared stable and before the first release. Each change of a default lands before that too. C++ is excepted. That a plan entry records no realized native name is what the design asks of #195, and is not fixed by it.
- **Research 0016.** In the table "Applied", row 13 has the validation that the decisions of 13 November 2026 give, the new row 13a defers Go's unexported spelling with its trigger, and rows 13, 14, 14a, 23 and 29 carry dated annotations. The submitted document and the advice are unchanged.
- **The case record.** [bittype/ACX-004](docs/architecture-cases/0004-generation-file-syntax.md) gains the Go naming instance as pending: generation decides the case of every Go name, through a recipe that all four targets share, before exposure is known, and no Go-owned rule checks it. The same shared function holds Rust's reserved word `Self`.
- **The `Self` rule.** Finding 10 of [package boundaries](docs/package-boundaries.md) and the candidate in [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md#candidates) record the answer: a reserved word is the target's. The rule is in the wrong place and not corrected. Two candidates of that case stay unconfirmed, where the documents said three.
- **Other documents.** [Semantic generation](docs/semantic-generation.md#generic-and-naming-limits) says which of its naming statements describe today's default and why they change. [Representation families](docs/representation-families.md#families) points to the page's reading of its sentence on a lifted binder's name, which the page marks as derived. [The planner page](docs/generation-planner.md) lists the Go naming case among the open ones. `LAYOUT.md`, the [documentation map](docs/README.md) and the [case index](docs/architecture-cases/README.md) link the page.

Not done:

- The delivery of #274, per target: the selectable policy, the reserved word and the `__` refusal out of the shared function, the intent in the inventory of [#272](https://github.com/Bitspark/bittype/issues/272), the Go-owned package for Go's check, and each target's default. Its use in Go Files is slice 1 of [#161](https://github.com/Bitspark/bittype/issues/161).
- Delegated, as the page lists: to [#195](https://github.com/Bitspark/bittype/issues/195) and research 0017, whether a plan entry may omit the "native symbols" that [binding generation](docs/binding-generation.md#the-selected-representation-plan) puts into the plan, and which identity it records. To [#164](https://github.com/Bitspark/bittype/issues/164), the access of a lifted inner binder and the reading of its name. To slice 1 of #161, Go's exact rule for what is in a public representation.
- Not specified by the page: how the parts of an intent reach the inventory record, and the identities of the two policies.
- The texts of #274, #161 and #268 still have an unexported Go spelling, "a legal name in Rust" and #274's "Order". The decisions of 13 November 2026 change them, and this change edits no issue.
- #268, #207, #195 and #164 are not commented or edited by this change.
- No code, test or generated report is changed.

## Research 0016: native Files placement, evaluated and in delivery — 12 November 2026

[Research 0016](research-docs/0016-native-files-placement.md) asks what the first Go and Rust Files slices (#161, #159) must fix for good: the placement input, the unplaced generated output and the placement check. It was written after a reading of the accepted documents left five points open, reviewed by four readers, and answered through the consult service. This entry adds documents only.

- **Evaluated.** The advice was split into 36 recommendations, and each was checked against the code at `7e6e02e`. 15 rows hold as written, 17 are accepted with a named adjustment, one is stale and five are deferred with a trigger. None is rejected. The table is the section "Applied" of the document.
- **What it decides.** Each language owns its placement type, and one shared admission runs the common rules. A placement is a typed value built from handles that the sealed output issues, so it cannot be used with another output. Generation records a naming intent, and each language's Files stage realizes the spelling, which in Go includes the case of the name. Unplaced TypeScript, Go and Rust declarations get their own root. Rust external destinations are recorded now. The Files stage tells unsupported coverage from invalid input.
- **What it withdraws.** The statement that Files never renames a declaration. The leaning to one shared placement type. The deferral of Rust external paths in [bittype/ACX-004](docs/architecture-cases/0004-generation-file-syntax.md), which the case record follows with #276.
- **C++.** Its move stays deferred until C++ has a semantic provider (#167), and the guarantee that generated output does not change after the first slices excludes legacy C++ output. #160 records it.
- **Delivery.** #271 to #278, then slice 1 of #161 and of #159. The document reads `delivering` until every accepted row is implemented and validated.

## horos scope values and the record tag take their names from the vocabularies, and the published horos names are pinned — 12 November 2026

The two values of the horos `scope` attribute have one declaration in `syntax/horos/go`, the horos source reader builds the record tag from the vocabulary's prefix, and a test pins every published name of `semantics/horos/go`. These are items 8, 6 and 9 of [#256](https://github.com/Bitspark/bittype/issues/256). It is a move: every value is unchanged, no import is new and no architecture test is edited. One of the two tags of item 6 is confirmed and not corrected, because its owner is not importable by the two packages that spell it.

- **The scope values, item 8.** `syntax/horos/go` declares `ScopeLocal`, which is `local`, and `ScopeImported`, which is `imported`, beside `AttributeScope`. Four packages spelled them sixteen times. The parser and the model projection in `internal/horosmodel/project.go` write them. The canonical reader in `internal/horosmodel/read.go`, the admission grammar and two comparisons in `internal/declsyntax`, and the formatter read them. All use the two constants now, and the grammar's value set of `reference` and `kind-reference` is `g.Set(h.ScopeLocal, h.ScopeImported)`. The capability report still lists `local, imported` for both.
- **Two names with one spelling.** The same lines of `read.go` build the tags `horos/decl/local`, `horos/decl/imported`, `horos/decl/kind-local` and `horos/decl/kind-imported`. Those are forms of the canonical declaration format, which the horos library owns and reads. The attribute value is a name of bittype's tree. One can change without the other, so the tag names stay spelled where they are written and do not use the two constants. With a scope value changed in the vocabulary, every canonical value that the tests produce is unchanged.
- **The record tag, item 6.** `Fields` in `transform/syntax/horos/source/go` compared a payload's tag with `horos/record`. `internal/horosmodel` writes that tag from `DescriptorTagPrefix` and the constructor. The reader now compares with `DescriptorTagPrefix + "record"`. It already imported the vocabulary.
- **The field tag, item 6: confirmed, not corrected.** The parser writes the label `horos/field` of a record field wrapper and the formatter compares with it. The tag is in the same namespace as the record tag, and [the horos document](docs/horos.md) says that descriptor operands "preserve the exact canonical value layout" and that field wrappers "retain their literal roles". The admission grammar does not know the label, so it is not a syntax name. bittype's one spelling of that namespace is `DescriptorTagPrefix` in `semantics/horos/go`, and the role of a Syntax stage does not permit that package. Seeded, the import fails `TestProductionArchitecture` for both stages. No allowance is added. The decision is open on #256.
- **Pins, item 9.** `semantics/horos/go/vocabulary_test.go` is new. It compares the four published strings and 26 keys with literals, the two contracts and the configured analysis contract with literal obligations, domains and verdicts, and the prelude pin with the horos library's. It tests `DescriptorOperator` and `DescriptorConstructor`: the round trip for the 24 constructors that the parser writes, the bare prefix, a nested prefix, another owner and operators that are no descriptor. One test reads the package's source and fails for an exported name that the file does not cover. `syntax/horos/go/nodes_test.go` is new and pins the two scope values.
- **What fails when a value changes.** 42 changes were made to `vocabulary.go`, one per run: each published value, seven rules of the contracts, two of `DescriptorConstructor` and one new exported key. Each fails the new file. The eleven test packages that link a horos vocabulary were run for eight of them and for the two scope values. No test that existed fails for a changed `ScopeLocal`, `ScopeImported`, `ModuleKind()`, `TypeDefinitionKind()`, `ImportRelation()`, `ProperSort()`, `ModuleIdentity()` or descriptor operator prefix. That confirms by mutation what the review of #247 and #248 had found by search. One existing test fails for the term encoding, and 117 for `DescriptorTagPrefix`.
- **Unchanged output.** A temporary probe recorded what the existing tests produce in a full run: every tree the parser and the projection return, every admission result, every canonical module and literal value as codec bytes, every formatted text and every list of record fields. Before and after, the 2136 distinct records have the same SHA-256. The before and after lists of `go test ./...` differ by the seven new tests and nothing else.

Not done:

- The field tag, as above. Two ways were seen. The canonical value tags get a vocabulary that a Syntax stage may import, which is a decision about where #250 put the two prefixes. Or the private model projection, which both text stages already import for pins and literals and which writes every other tag, also names the field label. Seeded, the second passes the architecture test without a new allowance, which shows only that no edge is forbidden.
- The other atoms of the record payload have the same writers and readers and are not assessed: `req` and `opt`, which the parser writes and the formatter, the source reader and `verification/files/source/horos/go` read, and `open` and `closed`. The constructor names, `record` among them, are item 11 of #256.
- The tag names after the two prefixes stay spelled in `internal/horosmodel/read.go`. Item 7 of #256 corrects the sentence in the documents that gives them one writer.
- `horos/files/` begins a tree identifier in `parse/files/horos/provider/go` and in `internal/horosfiles`. Both write it, and no reader of that spelling was found. Not assessed.
- No test file other than the two new ones is edited, and no document besides this entry.

## The structural profile has its own identity, and each mapping states its profile — 12 November 2026

The two generation profiles of the horos mapping side are named for what they promise, and a mapping states its profile itself. No shared code lists, writes or refuses a profile. This is [#216](https://github.com/Bitspark/bittype/issues/216), pull request [#265](https://github.com/Bitspark/bittype/pull/265), and corrects [bittype/ACX-008](docs/architecture-cases/0008-structural-profile-gate.md). The generated output is unchanged between the merge base `d35f1fa` and this change: 108 records of existing integration tests, 79 formatted texts and 29 output trees, have the same SHA-256, and so has the text of three further fixtures in the four targets, twelve files. 62 of the 79 texts come from formatter-only tests that do not generate.

- **The identities.** `RepresentationProfile` is `structural-representation/1` and `EnforcementProfile` is `structural-enforcement/1`. Before, both began with the word `pagination`, the name of the first fixture. The owner of research 0015 chose the names and recorded the reason on #216. The Go constant names and the version `/1` are unchanged. There is no alias: nothing is released, and a request under an old name finds no mapping and fails.
- **A mapping states its profile.** A profile is a promise about the native type, so the mapping that makes it states it. Every builtin mapping writes `RepresentationProfile` into its own description, and its constructor takes no profile: `text.Library()` where it was `text.Library(profile)`, and for Go `record.Library(key)` and `optional.Library(key)`. That includes the declarations mapping of each target. A composition chooses which libraries it registers and which profile it requests.
- **The common declaration recipe states no profile.** `declarations.Library(profile, record, alias)` in `transform/syntax/horos/mapping/declarations/go` is now `declarations.Recipe(record, alias)`. It returns the recipe, and a mapping puts it into its own description. A contribution of another profile uses the same recipe.
- **No closed list, and no refusal by name.** `Library.Describe` in `transform/syntax/horos/mapping/go` refused every profile but `RepresentationProfile`. `Library` is now a function type with the method `Describe`, and it returns what the function returns. The registry already selects the mappings that state the requested profile, and that selection is the only rule that admits one. A second profile is a contribution that states its own identity in its own package, and it changes nothing in the horos mapping package.
- **The enforcement profile.** No mapping in bittype states it, so a request for it fails with `ErrUnsupported` in the registry, with no plan and no output. A profile that no registered mapping states fails the same way. That behaviour existed and is unchanged. The registry and `contracts/**` are not edited.
- **Two earlier forms were replaced before merging.** At `fd75abe` this pull request kept `Library(profile)` and a refusal of the one name `EnforcementProfile` in the adapter. An independent review showed that a five-line wrapper went around the refusal, and the holder of #216 decided that the profile moves to the mapping that promises it. At `d61bc5d` the common recipe still wrote the structural profile for the four declarations mappings, and no test failed when the shared function type refused a stated profile again. A second independent review found both. bittype/ACX-008 records the forms and the reasons.
- **One name, one place.** The integration fixture `horosgeneration.Request` spelled the identity as a string and now takes it from the horos mapping package. One test beside the declaration pins both values on purpose. No other test or fixture spells either identity.
- **Tests.** In `integration/go/horos_generation_profile_test.go`, one test builds every builtin library through its exported constructor, describes it and compares the stated profile. It reads the pairs and the mapping packages from the directories by name. Another registers the builtin TypeScript family and a fixture contribution, `integration/go/testdata/readonlyprofile`, in one registry: the fixture is a `mapping.Library` function value with its own profile, maps a list to `ReadonlyArray<E>` and uses the common recipe, and each request receives its own family. Ported to the merge base, where the gate exists, that test fails. The file covers the two refusals. `mapping_test.go` pins the identities and covers the function type for functions that succeed and that fail. Seeded faults fail these tests: the function type refusing a stated profile, a builtin constructor that takes a caller's profile, a builtin mapping that states another profile, a mapping package that the test does not name, the fixture stating the builtin profile, and either identity changed.
- **Edges.** A nil description function fails with `ErrInvalid`. A description that states no profile is refused at registration by the contract's `FreezeDescription` with `ErrInvalid`. On `main` the adapter answered a library without rules with `ErrInvalid` under the structural profile and with `ErrUnsupported` under any other, and an empty profile with `ErrUnsupported`.
- **Imports.** The four target declarations packages now build their own description and import what the other mappings of their target import: the generation contract, the horos syntax, the horos source reader and the representation package. Their role allows it. The common recipe package no longer imports the horos mapping package, and its allowance in the architecture test is removed. No exception is added. The fixture `horosgeneration` imports the horos mapping package for the identity. The new fixture package and the tests import what they use.
- **Documents.** [Semantic generation](docs/semantic-generation.md) describes the profiles and the three rules. [Generation gaps](docs/generation-gaps.md), [native kind constraints](docs/native-kind-constraints.md), the [Page and Cell example](docs/page-cell-example.md) and [representation families](docs/representation-families.md) use the new identities. The examples in [kernel composition](docs/kernel-composition.md), [generation](docs/generation.md) and [package boundaries](docs/package-boundaries.md) call the constructors without a profile. `AGENTS.md`, `LAYOUT.md`, [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) and the two index tables follow. Six documents that still named the profile after the fixture in prose say "structural".

Not done:

- Two independent reviews read this change, at `fd75abe` and at `d61bc5d`. The fixes after the second are checked by its reviewer again, and that check is not in this entry.
- The registry does not judge a stated profile. A contribution of another author that states the enforcement profile, for example a wrapper around a builtin library's description, is selected by a request for it and yields representation output. That author states the profile and answers for it. `TestRegistryDoesNotJudgeAStatedProfile` reproduces this.
- A further profile can use the parts that state no profile: the function type, the patterns, the naming function and the declaration recipe. It reuses no builtin mapping, because each builtin mapping states the structural profile.
- The test of the builtin mappings names one exported constructor per package. A second exported constructor in an existing mapping package, a package variable or a setter would pass it.
- The fixture contribution exists for TypeScript only. No bounded profile exists yet.
- The provider matrix shows the builtin mappings of a pair as one entry, the shared function type. That was so before and is not examined here. The source link of that entry moves in three generated documents, which the matrices workflow regenerates.
- The `Self` rule in the shared naming function, which bittype/ACX-009 passed to #216, is not decided here. It is [#268](https://github.com/Bitspark/bittype/issues/268).
- The documents under `research-docs/` are verbatim records and keep the old names. Earlier entries of this file do not spell them.

## Bridge statuses declared once, and the helpers' request names pinned to the syntax vocabulary — 12 November 2026

The status of a compiler bridge's `Output` has one declaration per name in the bridge, and tests pin the names that the TypeScript helper reads from its request. These are items 3 and 4 of [#256](https://github.com/Bitspark/bittype/issues/256). The issue named TypeScript. The Rust bridge had the same shape for the status and is corrected with it. It is a move: every value is unchanged, no import is new and no architecture test is edited. The two helpers are not edited.

- **The status, item 4.** `internal/typescriptsem` and `internal/rustsem` each declare `StatusOK`, which is `ok`, `StatusInvalid`, which is `invalid`, and `StatusUnsupported`, which is `unsupported`. `resolve/files/typescript/provider/go`, `resolve/files/rust/provider/go` and `analyze/files/typescript/queries/go` compared `Output.Status` with those spellings and now use the constants. `incomplete` was spelled three times in the TypeScript bridge and four times in the Rust bridge and by no other production package. Each bridge declares it once as the unexported `statusIncomplete`.
- **Statuses that stay where they are.** `mismatch` and `environment` are read once in each bridge and become errors there. `compiler-errors` is read once in the Rust bridge, which answers `invalid` in its place. `missing-context` is written once in each bridge, and no package names it: the resolvers report every status they do not name as incomplete. The Go bridge `internal/golangsem` has no status. It returns an error, and its resolver reads the error's type.
- **No package shared between the languages.** Four spellings are equal in the two bridges. Each bridge owns the output of its own helper, the Rust helper never writes `invalid`, and no package compares a status of one language with another's.
- **What the TypeScript helper reads, item 3.** `compiler.cjs` reads 24 syntax names of its request by spelling: eighteen node kinds, the child roles `default`, `type-parameter` and `argument`, the reference role `target` and the attribute names `operator` and `name`. The issue listed two roles and no reference role. A script compared each with `syntax/typescript/go`, and all 24 have a constant of the same class and value. The helper is JavaScript and cannot import the package, so nothing but tests keeps the two sides equal.
- **What fails when a request name changes.** Each of the 24 names was renamed in the request on its way to the helper, one per run, with the helper and every Go package untouched. Tests that already existed fail for 22. None failed for the role `type-parameter`, which the helper reads to refuse more than sixteen generic parameters, or for the attribute `name`, which it reads to refuse the `this` type. `TestTypeScriptHelperReadsVocabularyNames` is new and fails for all 24. It runs the helper on a capture that holds every kind the helper lists and on four captures that each turn on one name.
- **What fails when a status changes.** Each status was changed on the Go side of the TypeScript bridge in turn. Tests that already existed fail for `ok`, `invalid`, `unsupported` and `mismatch`. None failed for `incomplete` as the helper writes it or for `environment`, so `TestHelperStatusesIncompleteAndEnvironment` is new: it runs the helper past its depth bound and with a foreign compiler module. No test fails for `missing-context`, which has a writer and no reader.
- **The Rust helper.** `internal/rustsem/helper/src/main.rs` reads 26 node kinds, the suffix `-variant` of three more, the child roles `member` and `lifetime`, the reference role `target` and the attribute name `name`. The script found a constant of the same class and value in `syntax/rust/go` for each, and the three kinds that end in `-variant`. The helper is not built on the workstation that made this change, so no name and no status of the Rust side was changed there. `TestRustSemantic*` and `TestSemanticProviders/rust` run the helper in the hosted checks.
- **The Go bridge.** `internal/golangsem` is Go and takes every kind, role and attribute name it reads from `syntax/golang/go`. It spells none.

Not done:

- The tests were shown to fail for a changed request and for a changed Go spelling. Neither helper was edited, so they were not shown to fail for a changed helper.
- Seen and not assessed. The TypeScript and Rust parsers `parse/syntax/<language>/provider/go` write node kinds through the constants of `syntax/<language>/go` and spell every child role and every attribute name by literal. For those names the tree that reaches a helper is written by the parser's spelling and not by the vocabulary's constant. The helpers also compare three attribute values, `readonly`, `this` and `static`, that the parsers, the TypeScript checker and the Rust lexical helper spell and no vocabulary declares.
- [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) and [package boundaries](docs/package-boundaries.md) still list both items under "Not examined". This change edits no document besides this entry.

## Rust graph names published by the vocabulary package — 12 November 2026

The seven Rust graph names that the compiler helper writes and the model readers of `semantics/rust/go` read by name are published by that package. This is the Rust half of item 2 of [#256](https://github.com/Bitspark/bittype/issues/256). It is a move: each key value is unchanged, and no Go code writes any of the seven.

- **Assessment.** The helper `internal/rustsem/helper/src/main.rs` is a Rust program and the reader is a Go package, so they are two parties. The bridge `internal/rustsem` decodes the helper's graph as a whole and names none of its keys. `model.go` of the vocabulary package is the only production Go code that reads a Rust graph name: no package under `resolve`, `analyze`, `format`, `verification` or `transform` does.
- **Published.** Fact relations: `DeclaredTypeRelation()` is `rust/declared-type`, `PredicateRelation()` is `rust/predicate`, `DefaultRelation()` is `rust/default`, `NativeTypeRelation()` is `rust/native-type`, `DisplayRelation()` is `rust/display` and `FieldsRelation()` is `rust/fields`. Term operator: `FieldsOperator()` is `rust/fields`.
- **Writers and readers.** The helper writes all seven. It cannot import a Go package and is unchanged. `model.go` reads the seven through these functions. Its private `fact` method takes the relation key and no longer a name.
- **How the helper and the vocabulary are kept equal.** By tests that run the helper and read its graph through the Go readers. `TestRustHelperNamesMatchVocabulary` is new. It reads the helper's graph through each of the seven names and reports one failure per name. Before it, no test opened a native type by id, which is the only read of `NativeTypeRelation()`.

Not done:

- No name was changed in turn to list the failing tests, as #255 did for TypeScript. The Rust helper is not built on the workstation where this change was written, so every test that runs it fails there before and after, and a changed name shows nothing. Which existing test reads which name was established by reading the tests: `TestSemanticProviders/rust` reads the declared type, the fields relation and operator and the display, `TestRustSemanticNativeTypes` the declared type, the fields and the default, `TestRustSemanticNamespaceAndCoordinates` the declared type and the fields, and `TestRustSemanticBoundsAndOutcomes` the predicate. The hosted `contracts` check builds the helper and runs these tests and the new one.
- Names with a writer and no Go reader that names them stay where they are written. From the helper: the entity kinds `module`, `struct`, `enum`, `variant`, `field`, `alias`, `trait`, `impl` and `parameter`, the parameter sort `type`, the operators `primitive`, `tuple`, `slice`, `array`, `pointer`, `reference`, `variants`, `native-type`, `trait-predicate` and `display`, the encodings `json-1` and `utf8`, the relations `variants` and `implements`, the expression role `declared-type`, the binding role `target`, and the diagnostic codes `unsupported`, `incomplete`, `mismatch`, `units`, `syntax` and `driver`. From Go: the diagnostic codes `limit`, `missing-context` and `compiler-error` of `internal/rustsem`. The helper also spells the subject of its diagnostics, `captured-type-coverage`, which the vocabulary publishes as `Coverage()`. No production code compares a diagnostic's subject or code with a name.
- Two neighbouring groups were seen and not examined. The status of the helper's response is spelled in the helper, in `internal/rustsem/compiler.go` and in `resolve/files/rust/provider/go/resolver.go`. The helper reads the syntax node kinds, child roles and attribute names of its request by spelling, and `syntax/rust/go` declares them in Go.
- The Go half of item 2 is its own change, pull request #264. It follows pull request #204, which edited the same files.
- [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) and [package boundaries](docs/package-boundaries.md) still list the Rust names under "Not examined". This change edits no document besides this entry.

## Go graph names published by the vocabulary package — 12 November 2026

The thirteen Go graph names that the compiler bridge `internal/golangsem` writes and the model readers of `semantics/golang/go` read by name are published by that vocabulary package. This is the Go half of item 2 of [#256](https://github.com/Bitspark/bittype/issues/256). It follows the change that added imported packages, which edited the same files. It is a move: each key value is unchanged, and the semantic graph and everything the model readers return for three fixtures are byte-identical before and after, in seven dumps.

- **Assessment.** The bridge and the vocabulary are two Go packages, so they are two parties: `export.go` of the bridge wrote each name as a string and `model.go` of the vocabulary package spelled it again to read it. `model.go` is the only production code that reads a Go graph name. The resolver, the analysis provider, the query service, the verification pair and its laws use the model's methods and name no graph key.
- **Published.** Entity kinds: `PackageKind()` is `go/package`, `DefinedKind()` is `go/defined`, `AliasKind()` is `go/alias` and `ParameterKind()` is `go/parameter`. Fact relations: `IdenticalRelation()` is `go/identical`, `UnderlyingRelation()` is `go/underlying`, `ConstraintRelation()` is `go/constraint`, `AliasTargetRelation()` is `go/alias-target`, `EmptyRelation()` is `go/empty` and `DisplayRelation()` is `go/display`. Term operators: `BasicOperator()` is `go/basic`, `SliceOperator()` is `go/slice` and `StructOperator()` is `go/struct`.
- **Writers and readers.** `export.go` writes the thirteen through these functions. Its private `expression` and `fact` take the key and no longer a name, and `object` holds the entity kind as a key. `model.go` reads the thirteen through them, and its private `fact` takes the relation key. Both sides are Go, so the vocabulary now makes the two spellings one.
- **Imports.** None is new. The bridge already imported the vocabulary, and no architecture test is edited.
- **What fails when a name changes.** Each published value was changed in the vocabulary in turn, and the four test packages that link it were run: `integration/go`, `verification/integration/go`, `internal/golangsem` and `analyze/files/golang/queries/go`. The writer and the reader follow the vocabulary together, so only a test that spells a name can fail. `PackageKind()` makes three tests and one subtest fail, because they spell `go/package`. The other twelve make none fail. Each name was then changed in `model.go` alone, so that the reader differs from the writer as it could before this change. Tests fail for twelve names, between one and twenty-three tests each. For `DisplayRelation()` none fails: no test of these packages depends on what `Type.String()` returns.

Not done:

- Names with a writer and no reader that names them stay where they are written, as `n.Key(name)` in `export.go`: the entity kinds `field` and `builtin`, the parameter sort `type`, the operators `interface`, `boolean`, `union`, `signature`, `tuple` and `display`, the encodings `json-1` and `utf8`, the binding roles `name` and `import`, and the expression role `type-expression`. The `import` role came with imported packages. The same holds for the diagnostic codes `incomplete`, `type-error` and `unsupported` of `resolve/files/golang/provider/go`. One of them becomes an instance when a production package other than its writer reads it by name.
- The setting key `imports-1` and the obligation domain `captured-package` are written and read inside the vocabulary package, which is one party.
- Test files still spell `go/package`, `go/underlying`, `go/display`, `go/name` and `go/import`. No test file is edited.
- [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) and [package boundaries](docs/package-boundaries.md) still list the Go names under "Not examined". This change edits no document besides this entry.

## Checker grammars typed by language, with their names taken from the syntax vocabulary — 12 November 2026

The admission tables of the four native checkers take every node kind, child role and attribute name from `syntax/<language>/go`, and the shared table type is typed by language. This is item 1 of [#256](https://github.com/Bitspark/bittype/issues/256). It is a move: every value is unchanged. The five `rules` tables, dumped as sorted text, and the five grammar inventories that `tools/provider-matrix` extracts from source are byte-identical before and after.

- **Who spelled them.** `syntax/<language>/go` declares the names. `rules` in `check/syntax/{golang,typescript,rust,cpp}/provider/go/check.go` was keyed by string literals and spelled each name again, in the table, in the kind lists beside it and in the code below it. The checker is a reader of trees that other packages build, so these are two parties.
- **The table type.** `internal/grammar` declares `Rules[L]`, `Rule[L]` and `Role[L]`. A table is keyed by `ast.Kind[L]` and a role lists `ast.Kind[L]`. `One`, `Opt`, `Many` and `Check` take the language's kinds. The compiler refuses a kind of another language as a key, among the kinds of a role and in a kind list, and it refuses a table applied to a tree of another language. Role and attribute names stay strings, as the vocabularies declare them. Checking and every diagnostic message are unchanged.
- **The four checkers.** 619 literals became vocabulary constants: 110 in Go, 162 in TypeScript, 200 in Rust and 147 in C++. Of them 301 are node kinds, 185 are child roles and 133 are attribute names. A comparison of a node's kind reads the typed kind of the node. No vocabulary lacked a constant, and `syntax/<language>/go` is not edited.
- **Left as literals.** Attribute values and diagnostic messages. The values are the channel directions of Go, the keyword names, literal kinds, modifiers and operators of TypeScript, the visibilities, primitive names and receiver forms of Rust, and the access levels, record keys, parameter keys and reference qualifiers of C++. No vocabulary declares them, and the checker is the party that states which values a tree may carry.
- **horos.** `internal/declsyntax` uses the kinds of `syntax/horos/go` directly, and its 33 conversion constants are deleted. `Check` and its two helpers lose the language parameter that [#254](https://github.com/Bitspark/bittype/pull/254) left without a second language, because the typed table admits a horos tree only. The caller `check/syntax/horos/provider/go` is unchanged, and `internal/horosmodel` is not edited.
- **The extractor.** `tools/provider-matrix/extract/grammar/go` is unchanged. It selects a table by the origin of its type, which is `internal/grammar#Rules` for every instantiation, and it reads a typed kind constant as a string constant. Its test calls the typed helpers, and `TestTableTypedByLanguageReadsVocabularyConstants` is new: it reads a table typed by language whose names come from a vocabulary package.
- **Imports.** None is new, and no architecture test is edited.
- **What fails when a name changes.** Thirteen constants were changed in their vocabulary one at a time: a kind, a role and an attribute name in each of the four native languages, and one horos kind. A changed kind no longer sets the checker against the parser and the target, which take kinds from the vocabulary. On `main` the checker refused the changed Go kind in 44 failing tests and subtests, and now in one, whose test tree spells the kind. The other failures of a changed kind are printed text and expected values that differ, because the formatter compares kinds by spelling and tests spell them. A changed role still makes the checker refuse trees, now those of the parser, which writes roles by spelling. A changed attribute name makes it refuse the trees of the parser and of the target, which both write attribute names by spelling. The changed horos kind makes no test fail: every horos party takes it from the vocabulary and no test spells it.

Not done:

- **The parsers.** `parse/syntax/{golang,typescript,rust,cpp}/provider/go` take node kinds from the vocabulary and write child roles and attribute names by spelling, 42, 64, 60 and 52 times. They are writers that a separate reader checks, so this is the same finding. It is outside the checker and is not corrected here.
- **The formatters.** `format/syntax/{golang,typescript,rust,cpp}/provider/go` use no vocabulary constant. They read roles and attribute names by spelling, 32, 58, 61 and 59 times, and compare node kinds by spelling, 15, 19, 21 and 11 times. With the role `type` changed in the vocabulary, the Go formatter and the TypeScript formatter each panicked on a generated tree that the checker had admitted.
- **The targets.** `transform/syntax/{golang,typescript,rust,cpp}/target/go` take node kinds and child roles from the vocabulary and write attribute names by spelling. The entry below recorded this as seen. The attribute mutations confirm it: the checker refuses generated trees when an attribute name changes.
- **The reference role.** `internal/grammar` compares a reference's role with the literal `target`, which each vocabulary declares as `RoleTarget` and which the targets and two bridges write through the constant. The shared package may not import a language, so a correction puts the name into the table, where the extractor reads `Reference` as a boolean. That changes the declaration format and is not a move.
- **A literal still converts.** An untyped string literal is accepted as a key and beside a typed kind. Only a call whose kinds are all literals is refused. The type separates languages and does not forbid spelling.
- The generated documents are not regenerated here. Regenerated before and after, they differ in 21 line anchors and in nothing else, and the provider-matrix workflow writes them after the merge.
- [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) and [package boundaries](docs/package-boundaries.md) still list the checker grammars under "Not examined". This change edits no document besides this entry.

## Go imported packages — 12 November 2026

- The first Go profile accepts an explicit import environment (#168). Captured files form packages by their paths under a module, and supplied dependency packages are given as source in receipt identity. Packages are checked in import order through a closed importer, without module or standard-library discovery. Qualified type references bind to imported declarations; supplied declarations are exported by identity only. The base configuration is unchanged.
- Supplied packages are kept out of the graph's structural projection for containment, not for secrecy: properties the compiler establishes for captured types, and human-readable text, may name supplied declarations. A supplied import path with a `vendor` element is refused.
- A Go resolution whose context ends returns the context's error without waiting for `go/types`, in the base profile too. Its whole native work runs beside the call as one unit: conversion, the package checks and the export, whose type comparisons cannot be interrupted either. Reconstructing the compiler state for queries and an admission query run the same way. `go/types` cannot be interrupted, so abandoned work keeps its goroutine until the call it is inside ends; [native semantics](docs/native-semantics.md) states the limit and the measured costs. Before, a 1,618-byte supplied package held a call with a 20 ms deadline for 11 seconds, and captured alias chains of depth 24 held a base-profile resolution for 8.8 seconds.
- An admission query's constructor must be a generic defined type. An instance used as constructor was accepted before, and a non-generic type with no arguments made `go/types` panic; both are now refused as unsupported.
- Analysis compares its whole configuration with the resolution's, and a test now fails when only lengths or keys are compared.

## Generation form names declared once by each language's model package — 12 November 2026

The names of the intended native forms have one declaration per language, in `transform/syntax/<language>/model/go`. This is item 5 of [#256](https://github.com/Bitspark/bittype/issues/256). The issue named TypeScript. Go, Rust and C++ had the same shape and are corrected with it. It is a move: every value is unchanged, and the plan steps, the generated tree and the printed text of the four targets are byte-identical before and after, in fifteen dumps.

- **Who spelled them.** A model package issues its language's forms through `Form(name)` and reads a name back through `Name(form)`. It had no named values, so every user passed or compared a string. Three packages per language did: the model itself, the target `transform/syntax/<language>/target/go`, which read the expression form names and `alias`, and the declarations mapping `transform/syntax/horos-to-<language>/mapping/declarations/go`, which wrote the two declaration form names. Go's optional helper wrote `defined` as well.
- **Published.** Exported constants beside `Form` and `Name`. In all four models `FormPrimitive` is `primitive`, `FormField` is `field`, `FormRecord` is `record`, `FormExternal` is `external` and `FormAlias` is `alias`. TypeScript adds `FormArray`, which is `array`, and `FormInterface`, which is `interface`. Go adds `FormSlice`, which is `slice`, and `FormDefined`, which is `defined`. Rust and C++ add `FormSequence`, which is `sequence`, and `FormStruct`, which is `struct`. The model's constructors and checks, the target's switch and the mappings use them. Go's `external` has one party, because only the model writes and reads it. It is declared with the other six so that the block is the whole list.
- **The two primitive names.** The primitive form carries a name, `string` or `bool`. The issue did not list them. The text and boolean mappings wrote the name by spelling, the model compared it by spelling, and the TypeScript, Rust and C++ targets compared `string` by spelling. Each model declares `PrimitiveString` and `PrimitiveBool`. The Go target compares neither. It writes the name into the tree as the Go type name, as before.
- **No package shared between the languages.** `primitive`, `field`, `record`, `external` and `alias` are spelled alike in the four models. Each model owns its own forms, a form can be read only by its issuer, and no package compares a form name of one language with another's or with a literal. The planner copies a form's name into the steps of the inspectable plan and compares nothing. The common declaration recipe receives the two declaration forms as values.
- **Imports.** None is new. The target and the mapping packages already imported the model of their language, and no architecture test is edited.
- **What fails when a name changes.** The 36 published values were changed in their owners, and `integration/go`, the one test binary that links the model packages, was run for each change. Six were changed one at a time and make tests fail: `FormAlias` and `PrimitiveString` of TypeScript, and `FormAlias`, `FormDefined`, `PrimitiveString` and `PrimitiveBool` of Go. A test composition spells five of them, and the two Go primitive names reach the generated Go code, which changes with them. The other thirty were changed together in one run, each to a value of its own, and no test failed. A further changed name can add a failure and cannot remove one, so none of the thirty makes a test fail alone. Every production party now takes these names from the owner, so a changed value cannot make two of them disagree, and the tests reach these forms through the model's constructors and the mapping libraries.

Not done:

- The test compositions `integration/go/testdata/generationfixture` and `integration/go/generation_planning_test.go` still pass `alias`, `defined` and `string` by spelling. No test file is edited.
- The C++ record mapping writes the external type `std`, `optional`, and the C++ target compares that pair to add the header. This is the place where Record decides how an optional field is represented. It belongs to [#196](https://github.com/Bitspark/bittype/issues/196), item 6, and is not changed here.
- Seen and not assessed. The four targets spell the syntax attribute names and values (`name`, `optional`, `package`, `visibility`, `absolute`, `global`, `record-key`, `path`, `system`) by literal beside the `Attribute*` constants of `syntax/<language>/go`. The horos descriptor constructors `record`, `scalar` and `list-of` and the discriminators `closed`, `utf8-text` and `bool` are spelled in the common mapping parts, in the horos source reader, in the mapping libraries and in `verification/files/source/horos/go`.
- [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) and [package boundaries](docs/package-boundaries.md) still list the form names under "Not examined". This change edits no document besides this entry.

## Documents record what #251 examined and what stays unexamined — 12 November 2026

[Package boundaries](docs/package-boundaries.md) and [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) listed two groups of names as not examined. [#251](https://github.com/Bitspark/bittype/issues/251) examined both, and the two entries below record the corrections. This entry changes documents only.

- **The scan.** A search of `main` at `9867398` for graph names spelled outside the vocabulary packages, the compiler bridges and the horos private parts found one outside reader, `format/files/typescript/provider/go`. Rust and Go have none. The scan did not cover spellings inside the Go and Rust vocabulary packages.
- **The case.** ACX-009 gains the two corrections of #254 and #255 with their authors' qualification and their limits. Its status stays "Corrected, with limits". The two instances were confirmed and corrected on the same day.
- **The symbol package.** `transform/syntax/typescript/symbol/go` keeps its own `typescript/name` role. It writes and reads the role itself and its consumers never name the key, so the role has one party. Decided on 12 November 2026 by the integrator of the plan.
- **Limits recorded.** No independent review read #254 or #255. `ImportRole()` is published and no test fails when it changes. The tests were shown to fail for a changed vocabulary and not for a changed helper. The language marker `L` of the horos framing is a remainder of a removed language.
- **Other documents.** `AGENTS.md`, `LAYOUT.md`, [native semantics](docs/native-semantics.md) and the two index tables follow. The counts of the architecture test were made again at `f821eb3`, and none changed.

Not examined, and tracked in [#256](https://github.com/Bitspark/bittype/issues/256):

- What the TypeScript helper reads by spelling: the syntax node kinds, child roles and attribute names of its request.
- The status of the helper's response, which the helper, the bridge, the resolver and the query service each spell.
- The TypeScript generation form names `alias`, `interface` and `record`, which the model, the target and the declarations mapping spell.
- The Go and Rust vocabulary packages' own readers against the names that their bridge and helper write.
- The admission grammars of the four native checkers, which are keyed by string literals beside the constants of `syntax/<language>/go`.

## TypeScript graph names published by the vocabulary package — 12 November 2026

The sixteen TypeScript graph names that the compiler helper writes and a Go package reads by name are published by `semantics/typescript/go`. This is the Go side of observation 1 of [#251](https://github.com/Bitspark/bittype/issues/251) for TypeScript. It is a move: each key value is unchanged, and the semantic graph, the groups, the relations and the printed text in seven dumps of three fixtures are byte-identical before and after.

- **Published.** Entity kinds: `ModuleKind()` is `typescript/module`, `PackageKind()` is `typescript/package`, `PackageTypeKind()` is `typescript/package-type`, `BuiltinKind()` is `typescript/builtin`, `AliasKind()` is `typescript/alias`, `InterfaceKind()` is `typescript/interface` and `ParameterKind()` is `typescript/parameter`. Fact relations: `NamedExportRelation()` is `typescript/named-export`, `DeclaredTypeRelation()` is `typescript/declared-type`, `ConstraintRelation()` is `typescript/constraint`, `NativeTypeRelation()` is `typescript/native-type`, `DisplayRelation()` is `typescript/display` and `PropertiesRelation()` is `typescript/properties`. Term operator: `ObjectOperator()` is `typescript/object`. Binding roles: `NameRole()` is `typescript/name` and `ImportRole()` is `typescript/import`.
- **Writers and readers.** `internal/typescriptsem/compiler.cjs` writes all sixteen. It is JavaScript, cannot import a Go package and is unchanged. `format/files/typescript/provider/go` reads the module, package, package-type and builtin kinds, the named-export relation and the two roles through these functions, and writes the kind of its Files group through `ModuleKind()`. `model.go` of the vocabulary package reads the alias, interface and parameter kinds, the other five relations and the object operator through them. Its two private `fact` methods take the relation key and no longer a name.
- **How the helper and the vocabulary are kept equal.** By tests that run the helper and read its graph through the Go readers. Each published name was changed in the vocabulary in turn, with the helper untouched, and both test binaries that link the vocabulary were run. Twelve names made tests that already existed fail. `BuiltinKind()`, `ParameterKind()` and `NativeTypeRelation()` made none fail, so `TestTypeScriptHelperNamesMatchVocabulary` is new: it formats a module that refers to `Array` and opens a parameter and a native type by id, and it fails for each of the three. `ImportRole()` makes no test fail. The formatter stores the binding of an import binding and looks up only the bindings of type references, so this read has no effect that a test could observe.

Not done:

- `transform/syntax/typescript/symbol/go` still spells `typescript/name` as its own `nameRole`, which it writes in `ExternalUse` and reads in `ReadExternal`. The role that the architecture test gives that package does not permit `semantics/typescript/go`: with the import, `TestProductionArchitecture` fails. This change edits no architecture test, so that second spelling of `typescript/name` remains.
- Names with a writer and no separate Go reader stay where they are written. From the helper: the entity kinds `program`, `script` and `property`, the parameter sort `type`, the operators `primitive`, `literal`, `tuple`, `array`, `union`, `intersection`, `boolean` and `display`, the encodings `json-1` and `utf8`, the relation `structure`, the expression roles `type-expression` and `native-type`, and the diagnostic codes. From Go: the diagnostic codes `limit` and `missing-context` of `internal/typescriptsem`, and the Files relation `type-import` and the setting `declaration-placement-1` of the Files formatter.
- Two neighbouring groups were seen and not examined. The helper reads the syntax node kinds, child roles and attribute names of its request by spelling, and `syntax/typescript/go` owns them in Go. The status of the helper's response (`ok`, `invalid`, `unsupported` and the rest) is spelled in the helper, in `internal/typescriptsem`, in the resolver and in the query service.
- The Rust and Go helpers of #251 are not part of this change. Observation 2 of #251 has its own entry below.
- [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) and [package boundaries](docs/package-boundaries.md) still list the TypeScript names under "Not examined". This change edits no document besides this entry.

## horos syntax node kinds, attributes and roles declared once, in the vocabulary package — 12 November 2026

This answers observation 2 of [#251](https://github.com/Bitspark/bittype/issues/251). `syntax/horos/go` and `internal/declsyntax` each declared the 33 node kinds, the 9 attribute names and the 18 role names of the horos syntax, and neither imported the other. The two lists were equal in every name and every value. `syntax/horos/go` owns the names, and no other production package spells one. It is a move: every value, the parsed and formatted fixtures, the syntax projected from canonical values and the published semantic graph are byte-identical before and after.

- **Why there were two.** Pull request #39 of 4 November 2026 made the private horos model generic over the language marker, so that horos and a second language, horos/wire, could share it. It added `internal/declsyntax` with the names as untyped strings, which fit a tree of either language. Pull request #49 of the same day removed the second language. Since then every production use of the framing is over `syntax/horos/go`'s `Language`. One unit test of `internal/horosmodel` builds a tree over `struct{}`.
- **`internal/horosmodel`** takes the names from `syntax/horos/go` and no longer imports `internal/declsyntax`. Its code is generic over the language marker, so it reads a tree's kind as a horos kind, `s.Kind(n.Kind)`, where it compares one, and the projector's `add` takes a horos kind. `resolver.go` spelled the attribute name `pin` inline and now uses the constant.
- **`internal/declsyntax`** keeps its constant names, because its admission grammar is a table keyed by strings. Each constant is defined from the vocabulary: a kind as `string(h.KindProgram)`, an attribute or a role as `h.AttributeName`. The grammar table `rules` is unchanged. `Check` spelled the root kind inline as `"program"` and now uses the constant.
- **The architecture test.** Two entries of `oneOfAKindRules` change. `internal/declsyntax` gains `syntax/horos/go`. `internal/horosmodel` gains `syntax/horos/go` and loses `internal/declsyntax`. Both packages belong to horos by `languageOwners`, and the import goes from a private part to the vocabulary of its own language. `syntax/horos/go` imports nothing new.

Not done:

- The framing stays generic over the language marker, although only horos uses it. Removing the parameter would remove the conversions. It changes the signatures of `internal/horosmodel` and `internal/declsyntax` and is a change of its own.
- The constants of `internal/declsyntax` stay exported. No package outside it uses them now.
- This change edits no document besides this entry. [Package boundaries](docs/package-boundaries.md) and [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) still list the node kinds as not examined.
- Observation 1 of #251, the names that the native compiler helpers write, is open.
- The checkers of the native languages were not examined. The Rust checker compares node kinds with inline strings beside the constants of its vocabulary.

## Go's contract identities and query service move, and the documents follow — 12 November 2026

Pull request [#246](https://github.com/Bitspark/bittype/pull/246) made the Go part of [#220](https://github.com/Bitspark/bittype/issues/220) and [#221](https://github.com/Bitspark/bittype/issues/221). It is a move. This entry also brings the documents up to date with [#247](https://github.com/Bitspark/bittype/pull/247) and [#248](https://github.com/Bitspark/bittype/pull/248), which have their own entries below, and it records [#250](https://github.com/Bitspark/bittype/pull/250), which has none. [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) records the corrections, their qualification and their limits.

- **Contract identities.** `semantics/golang/go` exports `ResolutionContract()` and `AnalysisContract()` in `contracts.go`. `resolve/files/golang/provider/go` and `analyze/files/golang/provider/go` no longer export `Contract()`. The Go analysis provider no longer imports the Go resolver.
- **Query service.** `Queries`, `NewQueries` and `Admission` are the package `queries` at `analyze/files/golang/queries/go`. `verification/files/horos-to-go/go` imports that package and no longer the analysis provider.
- **No alias.** Nothing remains at an old path. A consumer that called `Contract()` on a Go provider package, or `NewQueries` on the Go analysis provider package, names the new places.
- **The exception rule covers a query service.** An import of a query service package (`analyze/files/<language>/queries/go`) by a package that is not a composition needs an entry in the exception table that names its symbols, as an import of a stage provider package does. The second commit of #246 added this after the independent review. The verification pair has an accepted entry for `Admission`.
- **The architecture test.** The two open entries of the exception table are deleted. Fourteen remain, and all are accepted: the thirteen same-language `Check` and `Parse` calls and the pair's `Admission`. The eight hand entries in `TestProductionArchitecture` are converted: `semantics/golang/go`, `internal/golangsem`, the two Go provider packages and `analyze/files/rust/provider/go` take their role, and `internal/typescriptimports`, `verification/files/horos-to-go/go` and `verification/logos/go` are entries of `oneOfAKindRules`. 84 of the 127 production packages are classified by their role. 43 keep an explicit entry: 31 in `foundationRules` and 12 in `oneOfAKindRules`. The synthetic tests are 37.
- **Tag prefixes of canonical horos values.** Pull request #250 is one more move of the kind of #247 and #248. `semantics/horos/go` publishes the constants `DeclarationTagPrefix` (`horos/decl/`) and `DescriptorTagPrefix` (`horos/`). `resolve/files/horos/provider/go`, `verification/files/source/horos/go` and `internal/horosmodel` spelled both and now use the constants. `verification/files/source/horos/go` exports `ArrowValue(kind)`, and `verification/files/horos-to-go/go` calls it where it built the tuple tagged `horos/decl/arrow` by hand. The pair imports no ontos package directly, and its two external allowances for `ontos/core/go` and `ontos/data/go` are removed from the architecture test. The horos library owns the two prefixes. horos 0.10.2 spells them in `declarations/grammar.go` and exports neither, so bittype has one spelling of its own until it does. Its author reports that build and vet pass, that the tests with the pinned TypeScript fail only where the Rust helper is missing, and that a program that is not committed found the constants and the arrow value identical to the old literals and the hand-built value, also in encoded form.
- **Order.** The move was first prepared as draft pull request #231 on top of open pull request #204. #231 is closed. The move went to `main` first, because a misplaced responsibility moves before an addition ([AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)), and #204 rebases onto it. The reasons are in a [comment on #204](https://github.com/Bitspark/bittype/pull/204#issuecomment-6812747024).
- **Qualification, by the author of #246.** `go build ./...` and `go vet ./...` pass. `go test ./... -count=1` has the same outcome per package and per test on the branch and on `main` at `3c3b748`, apart from one new line for the new package, which has no test files. The tests that need the pinned TypeScript and the Rust helper fail on the author's workstation in both runs. A program that is not committed rendered the two contract values and the analysis contract's resolution requirement before and after. With the printed pointer address replaced by a fixed word, each output is 2522 bytes with the same SHA-256. Six seeded forbidden imports each failed `TestProductionArchitecture`. On the second commit, the pair calling `NewQueries` fails and removing the pair's exception fails.
- **Hosted checks.** `docs`, `contracts` and `native` passed on both commits of #246, `7ecdaa1` and `052a3e3`, and on #250 at `cadf495`.
- **Independent review.** A reviewer that wrote none of #246 read its first commit, `7ecdaa1`, on 12 November 2026. Its verdict is that the production code is a move and that no test assertion changed. It found one gap, which the second commit corrected: the verification pair could use any symbol of the query service. It found two loosenings, which were accepted and are recorded as limits of the check: three Go packages may import what their role gives the same package of TypeScript and Rust, and an unused allowance of the five converted packages is checked per role. The reviewer did not read the second commit. Its evidence is in a scratch directory on the workstation and not in the repository.
- **Documents.** `AGENTS.md`, `LAYOUT.md`, [package boundaries](docs/package-boundaries.md), [native semantics](docs/native-semantics.md), the [type-system architecture](docs/type-system-architecture.md), [analysis](docs/analysis.md) and the two index tables follow #246, #247, #248 and #250. Finding 12 of package boundaries and the inline horos keys of ACX-009 are recorded as corrected. ACX-009 has no confirmed instance left. Its status is corrected, with limits. Issues #219 to #222 are closed, and the plan in [#218](https://github.com/Bitspark/bittype/issues/218) is complete with this change. The generated provider tables change through their workflow.

Not done:

- Go has no `AnalysisContractFor(settings)`. It arrives with open pull request #204, which is not part of the plan in #218.
- No independent review read the second commit of #246, #247, #248 or #250.
- Two groups of names were not examined and are listed in the limits of ACX-009: the names that the native compiler helpers write and Go packages read, and the horos syntax node kinds that `syntax/horos/go` and `internal/declsyntax` each declare. [#251](https://github.com/Bitspark/bittype/issues/251) tracks them. Two of the three native writers are not Go, so one owner there needs a design.
- The profile gate and the profile names ([#216](https://github.com/Bitspark/bittype/issues/216), [bittype/ACX-008](docs/architecture-cases/0008-structural-profile-gate.md)) and what a mapping's `Expand` returns ([#196](https://github.com/Bitspark/bittype/issues/196), item 6) are unchanged.
- Three candidates are not confirmed: the thirteen same-language `Check` and `Parse` calls ([#218](https://github.com/Bitspark/bittype/issues/218)), the `Self` rule (a question for #216) and `CheckSingleScope` (with [#161](https://github.com/Bitspark/bittype/issues/161)).

## horos descriptor operators, term encoding and profile names owned by the vocabulary package — 12 November 2026

Three more names that a writer and a separate reader have to agree on have one owner, `semantics/horos/go`. They were found with the five keys of [#220](https://github.com/Bitspark/bittype/issues/220), under the plan in [#218](https://github.com/Bitspark/bittype/issues/218). It is a move: each key value, each string, the published semantic graph and the generated output are byte-identical before and after.

- **Published.** `DescriptorOperator(constructor)` is `horos/descriptor/<constructor>`. `DescriptorConstructor(operator)` is its inverse: it returns the constructor of such a key and reports false for every other key. `TermEncoding()` is `ontos/codec-v1`, the encoding that the data of horos terms carries. ontos owns the codec; that horos terms carry it is a horos agreement. `ResolutionProfile` and `AnalysisProfile` were published before.
- **Writers and readers.** `internal/horosmodel` writes the descriptor operators and the encoding of every provider term through these functions. `transform/syntax/horos/source/go` reads the constructor through `DescriptorConstructor` and compares the encoding of a payload and of a module identity with `TermEncoding()`. It refuses the same terms as before with `semantic.ErrInvalidRequest`, and an operator that is the prefix alone still reads as an empty constructor. `internal/horosfiles` compares the profiles of its stages with the two constants. No production package outside `semantics/horos/go` spells any of the three.

Not done:

- No document describes these three names as inline, and this change edits no document besides this entry.
- ontos 0.16.0 exports no identity for its codec from a Go package, so `ontos/codec-v1` is spelled once, in `semantics/horos/go`.
- A test in `analyze/files/horos/provider/go` still spells the encoding itself. It pins the value independently of the vocabulary.
- The remaining entity kinds, parameter sorts, term operators and diagnostic codes have a production writer and no production reader, and stay where they are written.

## horos keys published by the vocabulary package — 12 November 2026

The five horos-owned keys of [#220](https://github.com/Bitspark/bittype/issues/220) that a writer and a reader each spelled inline are published by `semantics/horos/go`, under the plan in [#218](https://github.com/Bitspark/bittype/issues/218). It is a move: each key value and the published semantic graph are byte-identical before and after.

- **Published.** `ModuleKind()` is `horos/module`, `TypeDefinitionKind()` is `horos/type-definition`, `ImportRelation()` is `horos/import`, `ProperSort()` is `horos/sort/proper` and `ModuleIdentity()` is `horos/identity/module`.
- **Writers and readers.** `internal/horosmodel` writes the two entity kinds and the two operators through these functions. `internal/horosfiles` writes the group kind and the import relation and reads the module kind through them. `transform/syntax/horos/source/go` reads the type-definition kind and the two operators through them and no longer declares its own copies. No production package outside `semantics/horos/go` spells any of the five.

Not done:

- [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) and row 12 of [package boundaries](docs/package-boundaries.md) still describe these five keys as inline. This change edits no document besides this entry.
- Other horos-owned names are still spelled where they are used. The `descriptor/` operator prefix has a writer in `internal/horosmodel` and a reader in `transform/syntax/horos/source/go`. The profile names that `semantics/horos/go` already publishes are repeated as literals in `internal/horosfiles`. The remaining entity kinds, parameter sorts, term operators and diagnostic codes have a production writer and no production reader.

## Import rules by role and the independent review — 12 November 2026

The last step of the plan in [#218](https://github.com/Bitspark/bittype/issues/218) merged with [#244](https://github.com/Bitspark/bittype/pull/244), at `71bfec9`. It changes test code only. [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) records the correction, its qualification and what remains.

- **Rules by role.** The architecture test derives a package's import rule from the role that its path gives (`integration/go/architecture_roles_test.go`). 78 of the 126 production packages are classified this way, and a new package of an existing role needs no edit to the test. 48 keep an explicit entry: 31 foundation and kernel packages in `foundationRules`, 9 packages of which there is only one in `oneOfAKindRules`, and 8 hand entries held for Go. The explicit entries for the foundation and kernel set were decided on #218 on 11 November 2026, because their layering among each other is a fact per package.
- **New in the check.** A package that names a language imports no package of another language. Every allowance is checked for use: an entry and an extension per package, a role per shape. The check reads the files that build constraints exclude. A role grants `internal/<language>` and `internal/<language>sem` only where the imported package is that language's helper. The entries of deleted packages are gone.
- **Independent review.** A reviewer that wrote none of the changes read the merged moves `d23542a..e83cf52` and found that they are moves. It found four defects in the check and in comments and none in moved code, and #244 corrected them. The review is recorded in a [comment on #218](https://github.com/Bitspark/bittype/issues/218#issuecomment-6809378108). Its programs are in a scratch directory on the workstation and not in the repository.
- **Documents.** `LAYOUT.md` replaces its per-component dependency table with a role table. [Package boundaries](docs/package-boundaries.md) says how the test realises each of the six rules, what the check cannot establish and what the review covered. `AGENTS.md` says what adding a package needs.

Not done:

- Go's contract identities and Go's query service are still in its provider packages. Their move is draft pull request #231, stacked on open pull request #204. Two open entries of the exception table and the eight hand entries in the test wait for it.
- Some horos-owned keys are still spelled inline ([#220](https://github.com/Bitspark/bittype/issues/220)).
- The profile gate and the profile names ([#216](https://github.com/Bitspark/bittype/issues/216), [bittype/ACX-008](docs/architecture-cases/0008-structural-profile-gate.md)) and what a mapping's `Expand` returns ([#196](https://github.com/Bitspark/bittype/issues/196), item 6) are unchanged.
- Three candidates are not confirmed: the thirteen same-language `Check` and `Parse` calls ([#218](https://github.com/Bitspark/bittype/issues/218)), the `Self` rule (a question for #216) and `CheckSingleScope` (with [#161](https://github.com/Bitspark/bittype/issues/161)).

## Documents follow the divided common mapping package — 11 November 2026

- `LAYOUT.md`, `AGENTS.md`, [generation](docs/generation.md), [kernel composition](docs/kernel-composition.md) and [type-owned mappings](docs/type-generation-plugins.md) list the adapter package and its siblings `pattern`, `declarations` and `naming` ([#238](https://github.com/Bitspark/bittype/pull/238)). The examples in [generation SPI](docs/generation-spi.md) and [generation planner](docs/generation-planner.md) call `patterns.Element`.

## Package boundaries, first corrections — 11 November 2026

Code changes that merged on 11 November 2026 under the plan in [#218](https://github.com/Bitspark/bittype/issues/218). The pull requests that moved code report byte-identical identity values and generated declarations before and after. [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) records the corrections, their qualification and what remains.

- [#224](https://github.com/Bitspark/bittype/pull/224): tests reach each pair's mapping libraries and the verification laws through one seam file each under `integration/go` (`libraries_*_test.go`, `laws_golang_test.go`).
- [#225](https://github.com/Bitspark/bittype/pull/225), closing [#219](https://github.com/Bitspark/bittype/issues/219): the architecture check names the exact symbols of every import of a stage provider package. An unused allowance, an unused exception and a dot import of a module package fail. Four package shapes are classified by their path. 94 unused allowances are removed from the per-package list.
- [#226](https://github.com/Bitspark/bittype/pull/226), [#227](https://github.com/Bitspark/bittype/pull/227), [#230](https://github.com/Bitspark/bittype/pull/230), [#232](https://github.com/Bitspark/bittype/pull/232): each native mapping library is its own package, `transform/syntax/horos-to-<native>/mapping/<name>/go`, and exports `Library(profile)`. `Libraries(profile)` is deleted. For Go the consumer creates the helper key and passes it to `record.Library(profile, key)` and `optional.Library(profile, key)`.
- [#228](https://github.com/Bitspark/bittype/pull/228): each horos-to-Go verification law is its own package, `verification/files/horos-to-go/mappings/{scalar,list,optional,record}/go`.
- [#229](https://github.com/Bitspark/bittype/pull/229), [#233](https://github.com/Bitspark/bittype/pull/233), [#234](https://github.com/Bitspark/bittype/pull/234): the Rust, TypeScript and horos contract identities are published by `semantics/{rust,typescript,horos}/go` as `ResolutionContract()` and `AnalysisContract()`. `semantics/horos/go` is new and owns the horos roles, relations and profiles, the limits setting and the prelude pin. The TypeScript `Compiler` type moved to `semantics/typescript/go`.
- [#233](https://github.com/Bitspark/bittype/pull/233): the TypeScript query service is `analyze/files/typescript/queries/go`.
- [#238](https://github.com/Bitspark/bittype/pull/238): the common horos mapping package is divided. The adapter `Library` with the profile gate stays at `transform/syntax/horos/mapping/go`. `pattern`, `declarations` and `naming` are packages under `transform/syntax/horos/mapping/<name>/go`. `Pattern` is now `pattern.Head` and `Declarations` is now `declarations.Library`. The profile gate is unchanged ([#216](https://github.com/Bitspark/bittype/issues/216)).
- The documents follow these changes, except #238: `LAYOUT.md` records the path form `<component>/<name>/go` for a component with several independently selectable contributions, and the examples in [kernel composition](docs/kernel-composition.md) and the [Page and Cell example](docs/page-cell-example.md) register each library.

Not done:

- Go's contract identities and Go's query service are still in its provider packages. Their move is draft pull request #231, stacked on open pull request #204.
- The per-package list in the architecture test still classifies every package outside the four path shapes.
- Some horos-owned keys are still spelled inline in `internal/horosfiles` and `transform/syntax/horos/source/go` ([#220](https://github.com/Bitspark/bittype/issues/220)).
- What a mapping's `Expand` returns is unchanged ([#196](https://github.com/Bitspark/bittype/issues/196), item 6).

## Package boundaries — 11 November 2026

- Add [package boundaries](docs/package-boundaries.md): the difference between parts that are wrongly grouped in one package and contributions that are wrongly coupled, the question that separates them, what to look for, the findings of a scan at `3bf45c3`, the proposed import rules by role and the order of work ([#218](https://github.com/Bitspark/bittype/issues/218)). Nothing is implemented.
- Add [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md): the mapping packages, the verification laws and two analysis provider packages group separate contributions, contract identities and horos names are published by provider packages, and the import check cannot see either ([#219](https://github.com/Bitspark/bittype/issues/219), [#220](https://github.com/Bitspark/bittype/issues/220), [#221](https://github.com/Bitspark/bittype/issues/221), [#222](https://github.com/Bitspark/bittype/issues/222)).
- Record the owner's request of 11 November in `AGENTS.md`, and note in `LAYOUT.md`, [generation](docs/generation.md) and [representation families](docs/representation-families.md) what the current check does not establish and that the mapping split precedes the families.

## Order of correction and the profile gate — 11 November 2026

- Record the owner's 11 November direction in [ACX-004](docs/architecture-cases/0004-generation-file-syntax.md): each misplaced placement decision moves to its Files formatter first and coverage follows additively. The Go package clause ([#161](https://github.com/Bitspark/bittype/issues/161)) and Rust visibility ([#159](https://github.com/Bitspark/bittype/issues/159)) can start; C++ ([#160](https://github.com/Bitspark/bittype/issues/160)) is deferred until its semantic provider, with the reason. The Rust external path is recorded as a further pending instance. No correction is implemented.
- Add [bittype/ACX-008](docs/architecture-cases/0008-structural-profile-gate.md): the shared horos mapping policy admits exactly one profile name, and that name is taken from the first fixture ([#216](https://github.com/Bitspark/bittype/issues/216)). Pending.
- Note in [representation families](docs/representation-families.md) the two changes that precede the additive work, and that ontos delivered the expanded-size measure in v0.17.0.

## Chain issues made startable — 10 November 2026

- Record two further pending instances in ACX-004: the Go package clause and the C++ namespace are fixed by target configuration during generation, while the placement proposal takes each declaration's native scope from the consumer's placement ([#161](https://github.com/Bitspark/bittype/issues/161), [#160](https://github.com/Bitspark/bittype/issues/160)). No correction is implemented.
- Point the delivery tables of [representation families](docs/representation-families.md) and [binding generation](docs/binding-generation.md) at the per-family issues split out of #165, the libraries/service and Content tracking issues, and the ordered plan in [#206](https://github.com/Bitspark/bittype/issues/206), which replaces #5.

## Representation families — 9 November 2026

- Record research 0015's accepted design ([representation families](docs/representation-families.md)). Each family is a representation contract, instantiated once per use, whose native type and plan entry are projections of one selected instance. The plan is a graph with a shared envelope and family-owned payloads declared in horos, five version axes and explicit generic requirements; it stays unstable until six vertical cases pass. Correct 0014's union rule: the obstruction is information loss, not overlap. Nothing is implemented yet.

## GEN-01 verified — 9 November 2026

- Record the independent verification of the GEN-01 correction: #192 at
  `6132170`, board record 000005. ACX-006 and GEN-01 read verified. The
  provider-integrity and single-site diagnostic limits remain; #193 tracks
  certifying that a provider's judgments rest on the admitted evidence.

## Binding generation ownership — 9 November 2026

- Record the accepted ownership of binding generation for horos-declared contracts
  (research 0014, [binding generation](docs/binding-generation.md)). bittype stays
  the representation authority and will publish a typed selected representation
  plan. Per-target binding backends compose paired native and adapter contributions
  from it. The charter gains the plan interface only; nothing is implemented yet.

## Generation, semantics and documentation — 8 November 2026

Changes since the previously recorded module-encoding entry:

- Add bounded Go, TypeScript and Rust semantic integrations with explicit compiler
  configuration and capture-bound facts; C++ semantics remain deferred.
- Implement public joint verification association, direct Logos checking and
  library-owned Horos/Go projections/laws, including independent reversed Pair
  composition and checked nested theorem reuse.
- Replace per-pair generation assembly and procedural mapping operations with
  independently registered descriptions, kernel-owned planning, one generic
  convention and separate native AST construction for all four targets.
- Carry generation.Output into TypeScript Files; decide imports/exports after
  placement and check external symbols, public roots and conflicting local names.
- Record open source traversal ownership (GEN-01), native enforcement (GEN-02)
  and Rust/C++ Files ownership (ACX-004).
- Correct GEN-01 (ACX-006, research 0013): the kernel admits the evidence graph and
  walks declarations, references, applications and children itself, with
  run-scoped occurrence handles and binder-keyed actuals. Sources supply
  classification, descriptor selectors and judgments (admitted, unsupported,
  rejected, undetermined) plus a reader for their libraries. Every actual is
  judged, including omitted ones. The Open/Admit/CheckChild/Catalog reader is
  removed without aliases; horos and the independent source, which now publishes
  its own graph, drive the same walk. GEN-02 and ACX-004 remain unfinished.
- Correct the GEN-01 delivery after independent verification rejected it:
  `PlanningInput` takes the source language, input and provider, and `Plan`
  admits the evidence and opens the source itself (`AdmitInput` is no longer
  exported). Handles resolve with `In(Captured)` against the exact capture. A
  source selector outside the capture's source language fails as `ErrScope`. A
  term recorded at several syntax sites is reported at the occurrence's
  declaration. A use whose callee is not generated or whose callee selector is
  foreign is reported as that use. `Plan` checks mapping
  descriptions before it admits evidence, as the registry does. ACX-006 and
  GEN-01 stay open until the correction is verified.
- Reconcile active architecture/API guides, Horos graph examples and dependency
  versions; add a documentation map, runnable Page/Cell fixture walkthrough,
  native placement proposal, exact verification contribution reference and
  current-carrier context table. Repair section links and clarify historical
  consultation authority. Move the maintained matrix identity case into Bittype.


## Horos module encoding — 5 November 2026

- Add the general `horos-module` command: one Horos source module to canonical
  Ontos bytes, sharing the existing parser and module lowering with source pins.
- Add a manual generation-only workflow that builds the encoder and lowers a
  selected repository source on a GitHub runner. No semantic resolution or
  prelude-specific logic is added to the command or agnostic kernel.
- Pin Horos Go 0.10.2, whose complete standard kinds module is authored in Horos
  source and loaded from the generated canonical asset. Keep Ontos at 0.16.0.

## Unified Horos declarations — 4 November 2026

- Replace separate Horos and Horos/Wire frontends with one revision-2 Horos AST
  and provider family. Remove old packages/profiles without compatibility aliases.
- Preserve explicit Type/Data parameter bounds in syntax and provider-owned
  semantic terms. Check spelling structurally; resolve/analysis delegate meaning
  to Horos and retain shared immutable evidence.
- Pin verified public Horos 0.9.0 with Ontos 0.16.0, removing the development
  pseudo-version and retaining the no-replacement installed-consumer gate.
- Migrate canonical fixtures and the independently installed pagination consumer;
  qualify one Page, DataPage aliases, data/mixed specialization and open forwarding.
- Update the provider matrix and ownership/model guides. No generators or runtime
  bindings are introduced; the agnostic kernel and native language stages are unchanged.


## Unreleased

### 4 November 2026

- Remove the Horos-to-HorosWire template-lifting transformer, its tests, exclusive
  indexing helpers and architecture allowances. Author optional pagination directly
  and use the same generic with data and Wire arguments. Preserve closed ground
  embedding and specialization; remove all compatibility paths. Update active
  design documents and align the public Horos dependency to 0.8.0.

- Remove `text.Provider` and the engine's combined provider API without aliases.
  Register parser, checker, formatter and printer adapters independently against
  language identity/profile metadata, in any order. Pack/PackCST validate generic
  carrier structure; explicit checkers preserve syntax admission for Check,
  semantic Capture and every transformation endpoint. Add partial-registration,
  ownership, preflight, concurrency, removed-API and installed-consumer coverage;
  update current documentation and the provider matrix.

- Move parsing, formatting and printing contracts to `parse.Provider`,
  `format.Provider` and `print.Provider`; move source input to `parse.Input`.
  `text.Provider` composes these independent contracts with `check.Provider`.
  Update consumers, independent positive/negative fixtures, dependency checks,
  documentation and the six-language provider matrix. Remove the old names
  without compatibility aliases.

### 3 November 2026

- Move AST syntax admission to `check.Provider` in `contracts/check/go`,
  independent of text and CST contracts. `representation.NewLanguage` accepts
  this interface; `text.Provider` composes it with parsing, formatting and
  printing. Update Horos and native providers, independent consumers, dependency
  checks and the provider matrix. Remove `text.Checker` without an alias.

### 2 November 2026

- Standardize the public provider interfaces as `language.Provider`,
  `resolution.Provider`, `analysis.Provider` and `text.Provider`. Individual
  text capabilities are `text.Parser`, `text.Formatter`, `text.Checker` and
  `text.Printer`. Update all current consumers, fixtures and documentation;
  remove the previous interface names without compatibility aliases.

### 31 October 2026

- Adopt the stage-first hierarchy: trees/{ast,cst,representation},
  contracts/{language,pipeline,text,transform} and uniform native stage paths.
  Move generic AST checks to check/ast. Remove superseded public paths and the
  engine's moved CST exports rather than retaining compatibility aliases.
- Consolidate checked AST/CST snapshots and descriptors in representation;
  preserve sealed payloads and construction-time marker checks across the split.
- Discover the complete native stage matrix in tests, enforce per-package
  production dependencies, and share provider conformance/round-trip/fuzz cases.

- Implement ASTTransformer, typed adapters and owned AST/CST registration
  envelopes with explicit language profiles and marker checks.
- Implement an empty-by-default syntax engine with explicit provider/edge
  registration, unique shortest paths, explicit routes, checked execution and
  existing native text composition. Add independent graph-oracle, ownership,
  negative public API and concurrent-registry tests, including a remote race gate.
- Keep canonical AST and actual language mappings deferred; record the delivered
  frame in [docs/engine.md](docs/engine.md).

- Add the generic CST and borrowed AST-view interface. Native parsers retain
  source tokens/trivia, formatters construct CST directly, and the shared printer
  emits supplied leaves without knowing languages or formatting conventions.
- Implement native TypeScript, Go, Rust and C++ parsing, formatting and syntax
  checks, with independent AST goldens, source round trips, malformed-input
  regressions, fuzzing and remote native compiler qualification.
- Replace the identity interface with LanguageProvider (ID/Name), embedded by
  every specialized text provider; compose providers explicitly. That native-text
  slice added no alias, registry, canonical AST or generator.
- Document the pipeline and ownership in
  [docs/parsing-formatting.md](docs/parsing-formatting.md).

- Add TypeScript, Go, Rust and C++ AST vocabularies over the unchanged generic
  carrier, with independent target consumers and language-separation checks.
- Add an agnostic language identity interface and provider ID/Name methods
  for the owner's intended central registry. That earlier slice added no
  registry, canonical AST or native text stages.
- Document node shapes, native syntax limits and the updated dependency/layout
  boundary in [docs/syntax.md](docs/syntax.md). Generators remain unimplemented;
  native formatting is delivered by the later text slice above.
- Document the Go AST-to-AST generation pipeline, target layout, composition
  strategy and dependency boundaries in [docs/generation.md](docs/generation.md).
- Add generic AST carriers, language-owned node kinds, containment/reference
  links and provenance, plus the independent typed `Stage[I, O]` contract.
- Add remote Go formatting, type-separation and dependency tests, and vet.
  That initial foundation included no language vocabularies or stage executors.
- Add the proposed canonical value-type model and a types-only implementation
  plan for TypeScript, Go and Rust in [docs/model.md](docs/model.md).
- Replace stale bootstrap guidance in the README, charter and contributor
  instructions with the current scope. Function, Stream and Cell stay
  library-owned; the model proposal is not a released language.
- Add remote documentation checks for local links, JSON examples and whitespace.

## Historical bootstrap — 26 September 2026

- Established the repository charter and component-first source layout.
- The original broader scope is retained in Git history. It is superseded for
  current work by the 31 October direction.
- Unmerged research at
  [`91363a5`](https://github.com/Bitspark/bittype/tree/91363a569c62899132ec365103abde7e41c760d9/research-docs)
  remains historical design evidence, not the current model specification.
