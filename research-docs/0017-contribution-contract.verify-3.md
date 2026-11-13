# Research 0017, verification 3: evidence, order and boundaries (R29 to R38, R40)

> **What this file is.** The report of one verification agent of
> [research 0017](0017-contribution-contract.md) on the recommendations
> R29 to R38 and R40 of the advice, checked against the code at `7e9dc34` on
> 13 November 2026. The agent did not write the research document. The
> report is a proposal and no decision. The decisions are the rows of the
> table "Applied" of the research document, and where a row and this
> report differ, the row governs. The design page
> [the contribution contract](../docs/contribution-contract.md) cites this
> report as verification report 3, and in its gate as V3.
>
> The text below is the report as it was written. Paths that begin with
> `C:/` are paths on the research owner's workstation and are kept as
> text. A path such as `platform/...` or `v2-inputs/` names a file that
> is not in this repository. 3 code expressions of the form `Name[T] (argument)` have a space before the parenthesis that the report did not have, because the documentation check reads the unspaced form as a link.

Read on 13 November 2026. Reading only: nothing was built or run, no board tool was used, nothing was posted.

**States read.**

| What | State |
| --- | --- |
| bittype `main` | `7e9dc34` (local equals `origin/main`). It differs from the expert's code baseline `2daea52` only by #280, a documents change |
| bittype open pull requests on the same contract package | #281 (inventory and handles, #272), #282 (#271), #260 (documents). None merged |
| bitspark-service `main` | `c456bc7`, 12 November 04:13Z (local checkout equals `origin/main`). Binding contracts on `main` are byte-identical to `c36722d` (blob `b849f97`) |
| horos, ontos | As pinned in `go.mod`: horos `core/go` v0.10.2, ontos v0.16.0, read from the Go module cache. horos `main` read through the GitHub API |
| Research 0017 | `Status: advised`, branch `research/contribution-contract-0017` at `5b23405` |

Short names: BG = `docs/binding-generation.md`, RF = `docs/representation-families.md`, BC = bitspark-service `docs/design/binding-contracts.md` at `c36722d`, AC = platform `docs/concepts/architectural-correctness.md`, SUB = the submitted research document.

## Facts the recommendations are checked against

| # | Fact | Source |
| --- | --- | --- |
| F1 | 21 mapping packages (5 each for TypeScript, Rust, C++; 6 for Go), 25 to 78 lines each. `contracts/generation/go` has 2,009 lines with tests, `registry.go` 273 | `transform/syntax/horos-to-*/mapping/*/go`, counted |
| F2 | bittype has no plan type, no payload schema, no completeness test backend, no conformance vectors, no serialized plan. "completeness test backend" occurs only in documents | search of `*.go` and `*.md` |
| F3 | `@bitspark/structural-adapters` does not exist. bitspark-service has the packages `api`, `cli`, `contract-generator`, `contracts`, `tools`, the branches `main`, `feat/tools-gate-4a-inputs`, `fix/web-retention-20261016`, no Go file and no horos module | GitHub tree and branches API |
| F4 | The package is decided to start with integers: "The package is created with its first real family, integers (bittype#163), once the unstable payload for integers exists. Our contribution rules forbid empty packages." | bitspark-service#38, decision comment, 11 November 18:27Z; `CONTRIBUTING.md:28` "New packages are created only with real implementation." |
| F5 | Integers come after the first slice: "2. Integers, TypeScript half \| #163 \| step 1" | bittype#206, Track B |
| F6 | The mismatch between F4 and #207 item 3 was never raised on bitspark-service#38. Its only comment is the decision of 11 November. bittype#196 (12 November 22:40Z) says "That question goes to that issue." | issue comments API |
| F7 | The TypeScript backend has no package and no code: "Whether the backend is contract-generator's successor or a new package is decided and recorded here, once bittype#195 and #196 fix the interface." Its own order: "Design of items 1 to 3 can start once bittype publishes the first unstable plan fixture." | bitspark-service#40, items 1 and Order |
| F8 | BC: "accepted by libraries/service on 11 November 2026. Nothing here is implemented yet" | BC:3-4 |
| F9 | service-api has structural adapters of the older route (`dataAdapter`, `recordAdapter`, `arrayAdapter`, `tupleAdapter`, `unionAdapter`, `genericAdapter`). They take a schema or child adapters, not a plan payload, and the record path reads a present `undefined` as absence: `if (source[name] === undefined && !acceptsUnit(selected)) {`, which fails for a required field and encodes the empty tuple for an optional one | `packages/api/src/wire.d.ts:54-65`, `wire.mjs:121-123`; bitspark-service#40 Gap: "the shared `dataAdapter` ... chooses representation rules itself" |
| F10 | The owner seat of bitspark-service is `codex-platform-status` (BC:6, #38, #40). Its last signed acts, all on 12 November: a comment on bitspark-service#47 at 07:51Z, comments on bittype#196 at 07:51Z and 14:48Z. The repository's last merge is #48 at 04:13Z. Its stated queue is lifetime: "#47's text goes first ... #42's implementation follows" (bittype#196, 07:51Z) | comment and commit dates through the API. Every session posts as `jm9e`, so the seat is read from the `Acted-by` line |
| F11 | horos v0.10.2 has `func Diagnose(v, t Value) Diagnosis` (verdicts `conforms`, `mismatch`, `not_well_formed`, `exhausted`) and `func (p *Program) GroundDescriptor(value TypeValue) (core.Value, error)`, with `Apply` beside it. bittype tests already call `GroundDescriptor`. horos TypeScript has `core.ts`, `diagnose.ts`, `cli.ts`, `index.ts` only: ground checking, no declarations model | module cache `core/go@v0.10.2/diagnose.go:191`, `declarations/model.go:152,180`; `format/syntax/horos/provider/go/format_test.go:148`; horos tree API |
| F12 | Go's vocabularies are types only. Syntax kinds: file, import, type declarations, type forms. Model forms: `Primitive`, `Collection`, `Field`, `Record`, `External`. A field is `g.Native(Form(FormField), []string{name, strconv.FormatBool(optional)}, ty)`: its name is a string | `syntax/golang/go/nodes.go:35-57`, `transform/syntax/golang/model/go/model.go:37-51` |
| F13 | Tests already compile and run handwritten Go against generated types in a temporary module (`go test ./...`, under `BITTYPE_NATIVE=1`), and type-check generated Go in process with `go/types`. CI also has a job with a temporary module that requires external modules ("Check an independently installed public consumer") | `integration/go/horos_generation_pagination_test.go:478-493`, `generation_independent_test.go:40-48`, `.github/workflows/go.yml:144` |
| F14 | An independent fixture source exists: "SourceLanguage is the independent fixture language. It is not horos." It has a generic `Page` with an additional native parameter, a `Cell` of a separate mapping library and a helper | `integration/go/testdata/generationfixture/fixture.go` |

## R31. The order (the weightiest)

### 1. Premises

| Premise of the advice | Verdict | Evidence |
| --- | --- | --- |
| P3 can run "through the actual service-library factory" | **FALSE today.** No such factory, package, plan reader or backend exists (F3, F7, F8, F9), and the owner decided the package starts with integers (F4), which the plan of work puts after the migration (F5) | The expert was told only "It will also host the structural run-time package" (SUB:140) and "one TypeScript factory round trip in the service library" (SUB 6.6, P3). It was not told F4 |
| V2 can be a "small Go value-expression vocabulary and emitter in the completeness test backend" that is "no production backend" | **PARTLY.** No such vocabulary exists (F12) and no test backend exists (F2). The toolchain to compile and run generated Go in a test exists (F13). See R29 for the charter | |
| The six vertical cases can be exercised in "thin executable versions" before the migration | **PARTLY.** Cases 3, 4, 5 and 6 can be made thin on the six delivered families with fixtures. Cases 1 and 2 need families bittype does not have (table below) | `docs/semantic-generation.md:164-167`: "Open records, Wire, tuples, sets/maps, unions/intersections, complements, exact/refinement forms and ground recursion remain outside this profile." |
| "Nothing is released", so a prototype before the rewrite is allowed | **TRUE**, and stronger than told. The first dependent outside bittype starts with integers (F4, F5), not at step 3 as SUB:292-294 says. bittype#196 (13 November 00:15Z): "A prototype on a branch before the contract is fixed on `main` is allowed" | |
| Moving the 21 packages twice is cheap; the cost lies in other repositories | **TRUE** (F1; SUB:252-255) | |
| Rust and C++ "migrate with binding support explicitly absent" | **TRUE**, already required: #207 item 11; BG:99 "a request for bindings that needs a missing facet fails explicitly" | |

### 2. Conflicts

- **The plan of work, #206.** Track C puts the structural run-time package after "the first payloads from Track B" and the backend after "the serialized plan of #195". Step 2 of the expert's order inverts that. It is possible only from a prototype branch, and only if `codex-platform-status` agrees to precede integers.
- **AC-10 and the open violation.** Record deciding the optional representation is "Open" (bittype/ACX-009, finding 7; `docs/package-boundaries.md:156`). The move needs the new contract (derivation Q7). Every step placed before the migration keeps that violation open on `main`. AC:292-295: "Deferring the move itself is permitted only when its destination needs a huge effort. Record such a deferral in the case with its reason and the trigger that ends it". The expert's order does not mention this.
- **"Delivered complete" (RF:17-19).** Not violated: "The prototype families need not become production deliveries" (advice 4.2). A thin recursion or union family on a branch must not land on `main` as a types-first family.
- **Research 0015 row 39 (HOLDS).** The six cases are the freeze gate of the envelope. Thin versions before the migration do not discharge that row. No conflict if this is recorded.
- **Derivation Q4** recommended "option 1 in the first slice, and record that case 4 is not claimed by it", and SUB leaned to V3. Both are this seat's derivations, not accepted designs. Accepted rows point the expert's way: 0015 row 3 validation, "a Go vertical case (row 39, case 4) lowers through a family facet with no backend change", and row 39, "The Go case runs in bittype's own completeness test backend (#195)."
- **Other recommendations.** Step 2 presupposes R37 (an authoring route) and R34 (a source-semantics transport). Step 3 presupposes R26 (who issues the lowering facet kind) and R29. "Including ... realization" presupposes R23; only TypeScript has a Files stage (`format/files/typescript`; Go and Rust have none).
- **The lane.** #281, #273 and #274 edit `contracts/generation/go` now. #207: "one lane, or an agreed sequence." A long-lived prototype branch must be stacked on them.

### 3. What each of the seven steps needs that is not there

| Step | Needs from bitspark-service | Needs from horos or ontos | Toolchain | bittype alone |
| --- | --- | --- | --- | --- |
| 1. P1, specify | Nothing | Nothing | None | Yes. Not done: SUB 6.6 "We have done none of the following" |
| 2. P3, TypeScript through the real factory | The package (absent, F3; decided otherwise, F4); factories for text, booleans, lists, records, optional fields and generic applications (absent, F9); a reader of the serialized plan (absent, F7); a call context with the four measures (BC, not implemented; today's `AdapterContext` has `mark`, `undo`, `find`, `add`, `rollback`, `commit`). Owner `codex-platform-status`, active on 12 November (F10), with a queue of #47, #42, #39 | horos TypeScript `diagnose` exists. Ground descriptors must come from the Go side, because horos TypeScript has no declarations model (F11) | Node and `tsc` are in bittype's hosted job. bittype pins TypeScript 5.9.3 (`semantics/typescript/go/vocabulary.go:12`); bitspark-service uses 7.0.2 (`package.json`) | Only the bittype half: prototype contract, plan, payload schemas, serialized form, a snapshot of descriptors |
| 3. V2, executed Go lowering | Nothing | horos and ontos Go, as pinned | Go, present (F13). New: generated glue that imports ontos needs a temporary module with requirements; the native job's temporary modules have none today | Yes, with a new value vocabulary (R29) |
| 4. Thin six cases | Real `Cell<T>` needs Function's contribution and the lifetime work (bitspark-service#39, #42, #47). The TypeScript side of all six is bitspark-service#40 item 9 | Forms exist in horos. A recursion case should not be built on v0.10.2: horos 0.11.0 notes "Recursive descriptors resolve `ref` lexically in every core, which changes some verdicts" (derivation Q8). Expanded size needs ontos v0.17.0 | Go | Cases 3 (ground stand-in), 4, 5 (fixture requirement), 6. Cases 1 and 2 only as kernel-level fixtures or as parts of #164, #210, #211, #212 |
| 5. P2, a materially different contribution | For "independently authored" in the strict sense: another owner, which needs the route of R37 | Nothing | Go | Yes, as a fixture in a separately owned package (AC:487-491). Independence of authorship then rests on a second seat |
| 6. Migrate 21 packages | Nothing | Nothing | Go | Yes |
| 7. Publish schemas, vectors, declarations | TypeScript qualification runs in their CI (#206: "Until that package qualifies a family, the family reports type support only for TypeScript") | Nothing | | The Go half |

**What a stand-in is worth (AC-07).**

| Stand-in | Establishes | Does not establish |
| --- | --- | --- |
| The Go test backend reads the **serialized** plan in a process that never held the planner's objects | No identity is a pointer; every link, version axis and refusal survives serialization | Anything about TypeScript |
| TypeScript factories written by bittype's seat on the prototype branch, run by Node from a Go test | The serialized plan, the payloads, the realization and the descriptors suffice for **a** reader in another language; the record factory can convert an optional field from the containing object; a contribution is found by its identity without Go | That **their** factory can be written: their adapter interface with budgets, their call context and memo, their package rules. Independent authorship: one seat wrote both sides |
| The same factories offered as a pull request or branch in bitspark-service, reviewed by its owner | Additionally their interface and their review | Their CI against published vectors, until vectors exist |
| The owner's written review of the first unstable plan fixture (their own trigger, F7) | A design-level statement by the second owner | Executed behaviour |

### 4. The six vertical cases, thin or not

| Case (owner) | Needs beyond the six delivered families | Owning issues | Thinnest form that still tests the contract | What that form does not show |
| --- | --- | --- | --- | --- |
| 1. Recursive `ContentValue`, shared DAG, cycle, expanded size (#164) | Recursion; `ContentValue` is "an untagged native union, `Uint8Array \| ContentMap \| readonly ContentValue[]`" (bitspark-content#54), so also bytes, maps, unions; a call context with memo, active set and expanded size; ontos v0.17.0 | #164, #208, #209, #210; Track B steps 4 to 9 | A recursive use graph in the fixture source (F14) with two fixture contributions: the planner reserves a reference before it fills the body, and a cyclic evidence requirement is refused (R9). Plan level, no values | The recursion family, naming of lifted binders, native indirection, memo and budgets. Go values of the six families have no pointer identity (derivation Q4: "Cycle cannot be built from these Go types") |
| 2. `Page<Outcome>`, optional nullable fields, exact discriminator (#210, #211) | Nullable values, unions, exact forms; "a recorded answer on #193" (#206) | #211, #210, #212, #193 | The optional field of `Page<Person>` converted through the contextual protocol in Go and in TypeScript, plus a second optional representation that needs operand evidence: Go's pointer form, which RF:165-166 names ("a pointer form needs proof that its `nil` cannot collide with a valid value") | Null against absence, and discriminator decisions. Those stay on paper (P1) |
| 3. Exact external `Cell<T>` (#196) | A paired contribution with capability conversion and lifetime; under bitspark-service#42 "Case 3 then expects a paired contribution" (bittype#196, 12 November) | #196; bitspark-service#39, #42, #47 | A ground exact generic with an erased or additional native parameter and a value facet. Fixtures for the type half exist: `TestGenerationPlanOmittedNativeArgument`, and F14's `Page[T0, Storage]` | Capabilities, holds, Cell's slot behaviour |
| 4. Go typed construction with a generic adapter argument (#195) | No family. A lowering facet kind, member and helper references, a value vocabulary (R29) | #195 | `Page[Person]` itself. This is the real case, not a thin one | Lowering of families that do not exist |
| 5. Generic map whose key requirement is rejected (#209) | Maps (after bytes and integers in Track B); a map form in the native models | #209, #208, #163 | A fixture generic whose contribution requires evidence about its argument's selected representation, accepted for text and refused for a list, through the real demand path | The map family and its key semantics |
| 6. Independent family, no shared change (#196) | Nothing | #196 | The pointer-form optional fixture of case 2, in a separately owned package | Authorship by another owner |

"Thin" stops testing the contract when the plan entries are written by hand instead of produced by the planner from a source, or when the fixture activates a path only it can reach (AC:496-497: "A renamed copy of the original example, a private test-only constructor or a fixture that activates a preinstalled special case does not establish that"). A fixture source with fixture contributions tests the kernel's protocol and mocks the family. That must be said in the record.

### 5. Three orders side by side

**(A) The document's order** (SUB 1.4, derivation Q7). **(B) The expert's order** (advice 4.2). **(C) The smallest order I can build from what bittype can do alone plus what bitspark-service has delivered** (a design document and the `Adapter` interface of service-api, nothing else).

Order C in steps:

1. P1 on #196: the two-column table of R30, the eight items of R29 answered per target, the refusal classes, the extension domain of R33.
2. One branch, stacked on #281, carrying the contract change for all 21 packages with byte-identical output. The four targets share the contract package, so the tree builds only if all move (or a second contract stands beside the first on the branch); the Rust and C++ packages are 25 to 61 lines each. It is not merged yet.
3. On that branch, for `Page<Person>` on Go and TypeScript: the plan at the unstable version, the payload schemas of the six families, the serialized form, and a snapshot of ground descriptors and module pins (R34).
4. Go: handlers for the six families reading the **serialized** plan; lowering for record, optional and generic application only, type-checked in process and compiled and run in the hosted job (construction of `Page[Person]`, member access, the optional helper, a generic adapter argument).
5. TypeScript: a Node reader with factories for text, lists, records with the optional member, and the generic application, written to the shape of service-api's `Adapter`, on the branch only. It reads the realization of the TypeScript Files stage.
6. Fixtures with the shared packages fixed: the pointer-form optional with operand evidence (cases 2 and 6, and R32's memo case, since a pointer has identity); the ground exact generic with an erased argument (case 3); the requirement-refusing generic (case 5); the recursive reference in the fixture source (case 1).
7. The first unstable plan fixture goes to `codex-platform-status` with two questions: whether it is sufficient for their adapter interface and call context, and the first-family mismatch of F6. Their answer is recorded. It does not gate the merge.
8. Merge the contract change once. Then the optional move as its own change. Then the additive work on `main`.

| | (A) Document | (B) Expert | (C) Smallest |
| --- | --- | --- | --- |
| Before the 21 packages change on `main` | Paper families; a fixture with the change | P1; TypeScript through bitspark-service; executed Go lowering; six thin cases; a different contribution | P1; the migration on a branch; steps 3 to 7 above |
| Packages prototyped | 0 | 11 (5 TypeScript, 6 Go) | 21 on the branch, which is the migration itself |
| Targets with executed conversion | 0 | 2 | 2 (TypeScript by a stand-in) |
| Families with executed conversion | 0 | 6 delivered, plus thin forms of recursion, nullable, union with exact discriminator and map | 6 delivered; 0 new builtin forms; 4 fixtures |
| Repositories that must act | 1 | 2, and bitspark-service must reverse F4 | 1; one request to the second |
| Establishes | Native output is unchanged; one fixture registers without a shared edit | Everything in R40 | Native projection; dependency composition with operand evidence; serialization; typed Go lowering compiled and run; a contextual optional field in both languages; a reader in another process and language; a different contribution with the mechanism fixed |
| Does not establish | That a decision record, a request or a facet has a reader at all. SUB 2.4 says so itself: "A contract can pass step 1 and be wrong in every part that the first slice exists to add" | Nothing it names, if it can be done | The service library's own factory and call context; recursion, unions, nullable and maps with values; real `Cell<T>`; realization for Go and Rust |
| Cost and risk | Lowest now. A second contract change is likely when the first reader arrives; inside bittype that is cheap (F1), and no other repository builds on it before integers (F4, F5) | Step 2 is circular under the recorded decisions (F4, F5) and waits on another owner. Step 4 in real form is Track B steps 2 to 9. The optional violation stays open on `main` for all of it | A long-lived branch on a contested package. The stand-in is test code in TypeScript that bittype's accepted texts do not provide for (#195 item 8 says "in Go"), so it stays on the branch or moves to bitspark-service |
| Against AC-10 and AC-07 | AC-07: byte-identical output supports no claim about the new half. AC-10 is at risk through the optional move: "Delivering a first change whose input or identity later has to change" | Meets both, if the deferral of the move is recorded. Otherwise the deferral is unrecorded, and its size is "huge" by another repository's schedule | My reading: AC-10 binds the contract through the first move. The optional move fixes the demand, its answer, the helper declaration and the replacement of one contribution. Order C's evidence covers exactly those parts in both languages. What it leaves is missing coverage or a later, justified AC-08 change |

Is anything "huge" in AC-10's sense? Step 2 of (B) in its literal form (another owner's package, against that owner's decision) and step 4 of (B) in real form (issues #164, #208, #209, #210, #211, #212). Steps 3 and 5 are not. In (C) nothing is; the deferral of the optional move still has to be written into bittype/ACX-009 with its reason and its trigger ("the evidence list on #196 is green and the contract change is merged").

### 6. Proposed verdict

**NEEDS ADAPTATION.** Adopt the principle (no final contract on `main` before a reader in another process and a compiled Go lowering have consumed it) in the form of order C. Per step:

| Step | Verdict | Trigger or adjustment |
| --- | --- | --- |
| 1 | HOLDS | |
| 2 | Real form: candidate DEFERRED. Stand-in: in | Trigger: `packages/structural-adapters` exists on bitspark-service `main` and its CI runs bittype's vectors for one of the six delivered families; or its owner records on #38 that a prototype branch precedes integers. Until then the record says "not demonstrated: the service library's factory" |
| 3 | HOLDS, with the home of R29 | |
| 4 | NEEDS ADAPTATION | Cases 3, 4, 5, 6 thin as above; cases 1 and 2 at kernel level and on paper. The real six stay the freeze gate (0015 row 39) |
| 5 | HOLDS | A second seat writes the fixture from the contract's documents only |
| 6, 7 | HOLDS | Step 6 is the merge of the branch |

### 7. Validation

A gate recorded on #196 item 6, each line a named test on the branch. Red before, green after:

- The completeness test runs from the serialized plan in a fresh process. Red while any link is a Go pointer.
- A seeded mutant in which the Go lowering takes a member by a string that the naming policy changed: the generated Go fails to compile. With a member reference it compiles.
- The pointer-form optional is registered and `Page<Person>` regenerates with the other optional representation. The diff of that test touches no file under `contracts/generation/go`, `transform/syntax/shared/generation/go`, the record mappings or the test backend.
- The Node reader round-trips `Page<Person>` and refuses a plan whose payload schema version it does not know.
- The fixture package, pinned at the commit of the contract change, still builds at the head of each later family. Red when a later family changes a contract type.

## R29. What the lowering environment must express

1. **Premises.**
   - "A member name as a free string is not enough": TRUE of today's code (F12), and of SUB's sketch, which keeps member names in payloads.
   - "Reflection can recover information that a typed emitter cannot obtain": TRUE; SUB 6.3 concedes it ("Reflection in V1 hides exactly these needs").
   - "It need not ... enter bittype's types-only native vocabulary": PARTLY, see 2.
   - Of the eight items, against `main` and #281:

   | Item of R29 | Today | With #281 |
   | --- | --- | --- |
   | Native type of a selected dependency | As a type operand only: `Mapped[T] (slot)`, opaque to the rule | Same |
   | Scoped native binders | `Parameter[T] (index)`, `Argument[T] (index)`; "indices only have meaning in their declaration's single generic convention" (`expression.go:61-62`) | Same |
   | Adapter arguments | Not expressible. `GenericConvention` has two native lists, indexed by source and by target parameter (`description.go:38-45`) | Same |
   | Generic requirements and evidence | Not expressible. Only `ParameterDescription.RequiredUse bool` | Same |
   | Generated declarations, external symbols, owned helpers | Type level: `DeclarationExpression` (issued by the planner), `External(module, name)`, `Helper(key, ...)` | A handle per top-level declaration, with roles requested, dependency, helper. It is local to one output: "an output-issued handle cannot be the key of a serialized plan" (`docs/naming.md:644`) |
   | Generated members with their containing type | Not expressible; a string in `data` | Not expressible. "One record per declaration of the snapshot"; a record states "nothing about compiler entities" |
   | Typed access and construction | Not expressible; no value-level kind or form exists (F12) | Same |
   | Lowering-owned helper declarations | Not expressible for a function; `HelperRecipe` holds a type declaration | `Values` can carry an owner-issued value; still no function |
   | The runtime call-context value | Not expressible. No Go call context exists anywhere | Same |

2. **Conflicts.**
   - BG:83: "It is scoped separately, with its own syntax and compiler checks. bittype's native vocabularies stay types-only." `docs/runtime-values.md:318-321`: "Generating such bodies requires separately scoped target syntax, checking, formatting and compiler qualification." RF:54-55: "source-code strings never stand in for typed target expressions."
   - A vocabulary in the test tree does not add to bittype's native vocabularies, so the letter of BG:83 holds. Three conditions make it hold in substance: the facet returns typed nodes and only the test's printer writes text; the Go compiler is the "compiler check"; and nothing in production imports it.
   - The third condition bites. A production Go mapping may import "the generation contract, the representation, the syntax of both languages, the source reader and the common mapping parts of its source language, and the model of its native language" (`LAYOUT.md`, dependency table, "Contribution: mapping"). It cannot import a vocabulary from the test tree. So in the prototype the lowering facets are fixture contributions too, or they are registered beside the family, which is the lookup by identity that R26 distinguishes from "selected together".
   - The production owner of value-level Go is unnamed: BG:46, "A Go backend belongs to the future Go runtime library." Accepted rows already require case 4 in bittype's test backend through a family facet (0015 rows 3 and 39), so this question exists whatever is decided here. The expert moves it forward. It is the derivation's stop condition 5.

3. **Requires that does not exist.** A Go value vocabulary (no owner; bittype can decide for its own test instrument and must record that the production home is open); a lowering facet kind and its issuer (R26); member and owned-helper references in the contract; a minimal call context for the Go test backend (bittype's own test code; derivation Q4).

4. **Size.** The list is eight rows to answer in P1 for two targets. Executed, the least is one function per nominal type of `Page[Person]`: a generic function declaration, parameters, a keyed composite literal, a selector, a call, a type instantiation, a symbol reference and a return, about eight node kinds, with lists converted by a handwritten generic helper in the fixture module. Three lowering facets: record, optional, declaration with its generic application. Nothing here is huge.

5. **Verdict.** HOLDS as the checklist of P1. NEEDS ADAPTATION for where each item lives: references that the orchestration must resolve (member, owned helper, binder, adapter argument, requirement) are categories of the kernel's contract; access and construction operations and the call-context value belong to the vocabulary of the product that reads the facet. The contract gains reference categories and no Go expression kind.

6. **Validation.** Red today: a fixture lowering for the record cannot name the member `Items` of the generated declaration except by a string. Green: it obtains a reference, and the generated function compiles (`go/types`, in process) and runs in the hosted job. Mutant as under R31.

## R30. What must be final and what may stay provisional

1. **Premises.** TRUE. SUB, head of section 6, states the partition and its own doubt: "Either the identities belong to the contract and are final, or AC-10 does not bind an artifact that is declared unstable." What the accepted page keeps unstable is the envelope's **version** (RF:95-96). The partition is this seat's derivation, and the derivation's own table already read "unstable" as "Not a provisional shape".
2. **Conflicts.** None with an accepted text. One limit: "the native symbol and realization relationship" can be demonstrated for TypeScript only; Go and Rust have no Files stage, and bittype#195 (13 November) already holds the stable version until "those two slices have changed their defaults". For three targets the left column is a design statement.
3. **Requires.** Nothing that does not exist. Seven rows on the left, each decided by another recommendation (identity R15, requests R3 and R7, generics R20, references and refusals R9 and R10, facets R26 and R27, symbols R23, lowering inputs R29).
4. **Size.** A table on #196 and #195. Under AC-07 a row is "decided" by text and "demonstrated" only by a test of order C.
5. **Verdict.** HOLDS. Replace the partition sentence; give each left-hand row its decision and its evidence, or the words "not demonstrated".
6. **Validation.** The pinned-fixture build of R31, validation 5.

## R32. Identity and linking cases

1. **Premises.** TRUE that the cases probe the contribution contract: they read what contributions supply. Which can be tests today:

   | Case | Today | When |
   | --- | --- | --- |
   | Repeated `Page<Person>` uses share a runtime token | No. No plan and no token exist (F2) | Plan level with the first plan. At value level only with a native that has identity: TypeScript objects, or a Go pointer. A Go struct of the six families has none |
   | Equal native types with different adapters stay distinct | No. No adapter exists | With a fixture exact mapping of a named text declaration to `string` beside the structural text family |
   | An erased native argument is retained for conversion | The type half is a test today: `TestGenerationPlanOmittedNativeArgument`, "A native projection may omit a source argument, but cannot grant admission" (`integration/go/generation_planning_test.go:379-427`) | The conversion half with adapter arguments in the plan |
   | One semantic plan under two placements | The type half can be written today for TypeScript: one output, two `tsfiles.Organization` values (`typescript_placement_test.go`; `docs/source.md:66`). I did not find a test that does exactly this | The plan half with R23. Go and Rust after slice 1 of #161 and #159 |

2. **Conflicts.** None. bittype#195's criterion of 12 November already states the first case.
3. **Requires.** The plan; R22 and R23.
4. **Size.** Four tests, three of them at plan level. Not huge.
5. **Verdict.** HOLDS. Add them to #195's Verify, with the note that value-level sharing in Go needs an identity-bearing native.
6. **Validation.** Mutants: the token taken from the use entry (case 1 red); the token taken from the native type (case 2 red); the adapter list zipped with the native argument list (case 3 red).

## R33. "Changes once" as a bounded commitment

1. **Premises.** TRUE. "Changes once" is this seat's delivery note (RF:328) and derivation; the rules bound it the same way: AC:501-503, "It cannot prove that all imaginable future requirements fit the interface. The promise is about the explicitly declared extension domain"; AC-08 "does not require an API to remain frozen forever".
   Is the extension domain written down? PARTLY, in three places, none complete: RF:22-24 ("Adding a family through the published identity, dependency, invocation, plan-graph and target-extension contracts changes only that family and the consumers that compose it"), the family list and delivery table of RF, and BG:259-263 (a second contract library, a second service, Go bindings).
2. **Conflicts.** AC:364-366: "The extension domain comes from the agreed responsibility, not retrospectively from the cases already implemented. Do not narrow the stated contract". So the written domain may not drop anything the accepted pages promise. The owner's "properly done ... no refactor later" (9 November) is not contradicted: the bound is the declared domain, not the delivered cases.
3. **Draft from the accepted documents only.**

   | Item the expert lists | In the domain? | Source |
   | --- | --- | --- |
   | New products | IN for a further backend per target and for a contribution's own facets: BG:97 "Optional product facets \| For example an operation facet"; BG:263 "Go bindings \| Only the Go backend and the Go implementations each semantic owner supplies"; BC S3. Whether a new **kind** of facet needs no kernel edit is stated nowhere; it is derivation Q4 | BG, BC |
   | Contextual representations | PARTLY. "optional and nullable fields" are families (RF:10-11). That the optional representation is replaceable apart from Record is bittype's finding (ACX-009 finding 7), while RF:73 lists "Record field layout and presence behaviour" as the record's payload. R11 decides this | RF, package boundaries |
   | Cross-field layout coordination | OUT by silence. No accepted text promises it; RF:73 leaves field layout with the record family | RF |
   | Application-specific specialization | IN. RF:254-257: "Explicit specialization is represented as such (#166), never as ordinary substitution into an unchanged generic declaration." bittype#195 (12 November), criterion 5: "A native recipe may depend on an actual ... or record it as explicit specialization (#166)" | RF, #206 step 11 |
   | Further contribution-authoring formats | UNDECIDED. Named once, in an accepted row: "as a Go mapping package or a declarative sidecar" (0014 row 12). No design page promises a format | 0014 |

4. **Requires.** A decision by the research owner for the two silent items, and the service library's view on products (operation and root facets are theirs). **Size.** One paragraph in RF and on #196. **Is it huge?** No.
5. **Verdict.** HOLDS, with the adjustment that the statement is assembled from the accepted pages and narrows none of them. Specialization being IN activates R21: the contract must be able to name it now.
6. **Validation.** One fixture per IN item with an empty shared diff; one recorded refusal per OUT item.

## R34. The source-semantics input

1. **Premises.** TRUE. BG:23-27 makes "the source semantics" one of three inputs; BG:77 has the backend emit "ground-descriptor constants"; bitspark-service#40 item 6 checks "against descriptors emitted at build time". No document or code gives a backend in Node a way to obtain them. What exists: horos modules are canonical ontos values with pins, so the **source** has a serialized form; the **checked** program (`Apply`, `GroundDescriptor`) exists in Go only (F11). bittype's source reader is a Go interface (`transform/syntax/horos/source/go`). bitspark-service has no horos dependency and no horos module (F3).
2. **Conflicts.** None. The envelope already lists "Checked source references and source semantic identity" (RF:73).
3. **Requires.** A serialized snapshot or a reader protocol. A TypeScript declarations model would be horos's work and is not needed for the six families. bittype alone can export, from its Go side, the module pins and the ground descriptor of every ground use and root.
4. **Size.** For `Page<Person>`: one table of descriptors keyed by the plan's source identities, and the pins. Not huge. A general reader protocol (operation and presentation data) is larger and has no consumer yet.
5. **Verdict.** HOLDS for the premise and for order C's step 3. NEEDS ADAPTATION: the smallest form is the descriptor and pin snapshot, packaged with the plan as the backend's input. The general protocol is a candidate for DEFERRED; trigger: design of bitspark-service#40 item 3 ("operation and presentation metadata") starts.
6. **Validation.** Red: the Node reader, given the plan only, cannot check `Page<Person>` and must fail explicitly. Green with the snapshot. Mutant: a snapshot of a module with the same name and shape and another pin is refused as a linking error (BG:254 lists that falsification test).

## R35. Target configuration

1. **Premises.** TRUE. Where it is recorded today:
   - `internal/typescriptsem/compiler.cjs:30`: `strict: true, exactOptionalPropertyTypes: true, noUncheckedIndexedAccess: true`.
   - `semantics/typescript/go/vocabulary.go:54`: a setting `environment` = `ES2020;ESNext;strict;exactOptionalPropertyTypes;noUncheckedIndexedAccess;lib.es2020;no-ambient-packages`, bound to the semantic receipt.
   - Test flags and one assertion: `// @ts-expect-error exact presence does not admit present undefined` (`horos_generation_pagination_test.go:472`).
   - `docs/semantic-generation.md:136`: "TypeScript uses strictNullChecks and exactOptionalPropertyTypes."
   - Not in the generation request or output. The TypeScript target has no configuration; Go's has `Package`.
   What the accepted designs say about absent and `undefined`: that the adapter "separates an absent own property from a present `undefined` at run time" (RF:167-168; #207 item 9; 0015 row 21), and that the plan carries "the rules for absence, presence, null and native `undefined`" (BG:58). None says what conversion **does** with a present `undefined`.
2. **Conflicts.** None in the documents. A real difference between the repositories: bitspark-service compiles with `strict` only (`test/types/tsconfig.json`, the fixture `tsconfig.json` files), on TypeScript 7.0.2, and its record adapter reads a present `undefined` as absence (F9). bittype qualifies its types under the option and TypeScript 5.9.3.
3. **Requires.** One decision by the structural profile's owner, bittype (RF:60-62): the admissible native values and the equivalence of the optional family on TypeScript. The envelope field "target and its configuration" gets its content.
4. **Size.** One decision and four vectors for one family on one target (absent, present value, present `undefined`; present null with #211). Not huge.
5. **Verdict.** HOLDS.
6. **Validation.** The vectors. Red against service-api's present `recordAdapter` for the present-`undefined` vector, which shows that it is not the family's factory. A compile fixture built with and without the option shows that type admissibility changes and the run-time result does not.

## R36. Budgets and memo

1. **Premises.** PARTLY: already specified. BC: "The memo key is the selected representation instance (its descriptor and its generic environment) plus the input's identity"; K1 "One budget per message per side, shared by all nested conversions"; K2 "A DAG whose expanded size exceeds the limit fails as exhaustion, even when its distinct objects fit"; K3 "One object under two representations converts once per representation"; K4 "A cycle is detected through active traversal and fails as `cycle` unless the family supports it". bittype confirmed K1 to K4 on 12 November (bittype#195, #196). RF:288-291: "ontos owns the encoded size, families account for their work, the service library owns the call context, and Content owns its limits." None of it is implemented (bitspark-service#40 item 5).
2. **Conflicts.** None. Not stated in BC: the state of a memo entry after a speculative rollback (L13), and that a hit charges the expanded size again; K2 implies the second.
3. **Requires / who owns.** The shared runtime contract is bitspark-service's. bittype owns each family's operational declaration (contract part 6, "Budget use") and its test backend's own minimal call context. For the six delivered families the question is mostly empty in Go: no identity, no recursion.
4. **Size.** Nothing before the migration.
5. **Verdict.** STALE as new work for bittype: the accepted BC already states it. Record a pointer and pass the two unstated points to `codex-platform-status` on #40.
6. **Validation.** BC's own test `K2 a small DAG with a large expansion is exhausted`; in bittype, vertical case 1 (#164). Red when a memo hit skips the expanded-size charge.

## R37. A contribution-authoring route

1. **Premises.** TRUE that none is resolved. SUB:265-269 says so. Today a contribution is a Go value registered in bittype's process; every library fixture lives in bittype's test tree. bitspark-service#39 item 2 asks for "Paired contributions for TypeScript through the public contribution interface of bittype#196" in `packages/contracts` and does not say in which language; that repository has no Go file (F3). Its present route is of the older kind: ordinary npm packages whose factories a generator imports statically (`test/fixtures/binding-libraries`). No sidecar format exists or is designed.
2. **Conflicts.** None with a design page. Related and open: bittype#172, "An application CLI lets platform's build run generation without its own Go composition" (#206, Freeze and release). A CLI cannot load a Go contribution package it was not built with, so the route and the CLI are one question.
3. **Requires.** A decision that reaches bitspark-service and platform's build, to be made with `codex-platform-status`. What bittype can do alone: show that an out-of-tree Go module can author and register a contribution using public packages only. CI already checks "an independently installed public consumer" in a temporary module (F13).
4. **Size.** Place the fixture of R31 step 5 in that out-of-tree module: one fixture, no new mechanism. A sidecar format would be a design of its own.
5. **Verdict.** NEEDS ADAPTATION. Record on #196 that the one route of the first contract is a Go contribution package, that "declarative sidecar" is named in 0014 row 12 and not promised, and test the Go route out of tree. The service library's own contribution is a candidate for DEFERRED; trigger: bitspark-service#39 item 2 starts against bittype#196's unstable interface.
6. **Validation.** Red: the out-of-tree fixture needs a symbol that is internal or lives in the test tree. Green: it builds and registers against the installed module.

## R38. Payload bootstrapping

1. **Premises.** TRUE. The limit is this seat's derivation (SUB, head of section 6: "A payload cannot hold an integer, such as a field index, until the integer family exists"). It comes from RF:84-87 ("bittype generates the Go and TypeScript types") together with the code: the native models admit two primitives (`PrimitiveString`, `PrimitiveBool`, `model.go:31-32`) and "Canonical integers remain Unsupported" (`docs/semantic-generation.md:164`). bittype#163 is integers; its Go half waits for #161 "because `math/big` is an import that only Go Files may place" (#206).
2. **Conflicts.** None: "No separate schema language is introduced" (RF:86-87; 0015 row 9). A point the expert did not see: the six families cannot type an **enumeration** either. A "presence strategy" or "construction strategy" (R5) is a union of exact values, which are #210 and #212. Typed as `text`, the schema validates less than it means, and narrowing it later is a change of that family's schema version.
3. **What a "pinned bootstrap artifact" is here.** bittype generates the payload types from the families whose payloads they are: a compiler bootstrap. The artifact is the family's schema as horos source, its canonical module value with its pin, and the generated Go and TypeScript types committed beside it, with a test that regenerates them and compares. horos does the same for its prelude (`core/go/declarations/standard/horos.kinds.horos` and `horos.kinds.ontos`; its README: "Bittype's ordinary single-module encoder lowers that source into the packaged canonical Ontos asset").
4. **Requires / size.** Five or six schemas for two targets. For each field that wants an integer or an enumeration: say so, and either carry it by structure (list order instead of an index, a boolean for a two-way choice) or record a bootstrap exception with its end (#163, #212). Not huge.
5. **Verdict.** HOLDS.
6. **Validation.** The regeneration test (red on drift). A payload with a strategy outside its domain is refused before any backend runs; red if the schema typed the field as plain text and no exception is recorded.

## R40. The acceptance condition of #196 item 6

1. **Premises.** TRUE that the present condition is weaker. Item by item:

   | The expert's condition | #196 item 6 and its Verify today | Where the item sits today |
   | --- | --- | --- |
   | Native projection | "The native declarations generated for today's pagination fixtures stay byte-identical in all four targets" (#207 item 5) | In the slice. Met by order A's step 1 |
   | Dependency composition | `Page<Person>` round trip in the Go test backend (#207 Verify). No operand evidence | Step 3 of the slice, after the contract is final |
   | Serialization | #195 item 9 and Verify 4 | Step 3, after |
   | Independent factory resolution | #207 item 3: "The TypeScript factories live in the structural run-time package ... and pass this repository's conformance vectors there"; Order 2: "waits for bitspark-service#38" | Outside the slice, in another repository |
   | Typed lowering | Vertical case 4 (#195 item 10) | Freeze gate of the envelope; the derivation would not claim it for the slice |
   | Native-symbol realization | Not in #196. Open on #195 (13 November, point 2) | Research 0017 |
   | Mechanism fixed while a materially different contribution is added | "a fixture family is added by registering one contribution. The diff of that test touches no file under `contracts/generation/go` and no existing mapping" | In the slice, weaker in two ways: it does not say "materially different" (AC:496), and it names neither the registry and planner package nor the test backend (#207 item 1 does) |

2. **Conflicts.** With #206's order, as under R31. "Independent factory resolution" in its real form cannot be a condition of a change in bittype without making bittype's `main` wait on F4.
3. **Requires.** As R31.
4. **Size.** As order C.
5. **Verdict.** NEEDS ADAPTATION. Make the seven lines the merge gate of item 6, each with the form it can honestly take: five real (native projection, composition with evidence, serialization, typed Go lowering, the different contribution), one for TypeScript only (realization), one by stand-in with the real form deferred under R31's trigger. Reword the Verify line to "materially different" and to the full set of shared packages. The priority list of R40 (operand evidence, generic conventions, identity, owned symbols, target configuration, authoring and process boundary) matches SUB 6.7's suspects 1, 6, 5, 4, 3 and 9 and adds nothing to verify here.
6. **Validation.** The gate of R31, validation.

## Proposed verdicts

| R | Short | Proposed verdict |
| --- | --- | --- |
| R29 | What lowering must express | HOLDS as P1's checklist; NEEDS ADAPTATION for the home: reference categories in the contract, operations in the reading product's vocabulary, prototype facets and vocabulary outside production packages. Five of eight items are not expressible today |
| R30 | Final against provisional | HOLDS |
| R31 | The order | NEEDS ADAPTATION: order C. Step 2 in real form DEFERRED (trigger: the package on bitspark-service `main` runs bittype's vectors for one delivered family, or its owner records that a prototype precedes integers). Step 4: four cases thin, two at kernel level and on paper |
| R32 | Identity and linking cases | HOLDS. Two type halves are testable today; the rest come with the plan |
| R33 | Bounded "once", written domain | HOLDS, assembled from accepted pages without narrowing; two items need a decision |
| R34 | Source-semantics input | HOLDS; smallest form a descriptor and pin snapshot; the general reader protocol DEFERRED (trigger: bitspark-service#40 item 3 design) |
| R35 | Target configuration | HOLDS |
| R36 | Memo and budgets | STALE for bittype: BC K1 to K4 already state it. Pointer and two questions to the service library |
| R37 | Authoring route | NEEDS ADAPTATION: one route now (a Go contribution package, tested out of tree); sidecar recorded as not promised; the service library's own contribution DEFERRED (trigger: bitspark-service#39 item 2 starts) |
| R38 | Payload bootstrapping | HOLDS; enumerations are the tighter limit than integers |
| R40 | Acceptance of #196 item 6 | NEEDS ADAPTATION: seven lines as the merge gate, each in its honest form |

## The three orders in one table

| | (A) Document | (B) Expert | (C) Smallest that meets AC-10 and AC-07 as I read them |
| --- | --- | --- | --- |
| Sequence | Contract change for 21 packages on `main`; optional move; additive work | Specify; TypeScript through the real factory; executed Go lowering; six thin cases; a different contribution; migrate; publish | Specify; migration on a branch; plan, serialized form, Go handlers and lowering, a TypeScript stand-in, four fixtures on that branch; ask the service library; merge once; optional move; additive work |
| Evidence before `main` changes | Paper and one fixture | All of R40 | Five of R40's seven lines real, one for TypeScript only, one by stand-in |
| Needs another repository | No | Yes, against its recorded decision (F4) | No. One request, not a gate |
| New builtin forms before the migration | 0 | About 4 (recursion, nullable, union with exact values, map) | 0 |
| Open optional violation on `main` | Until step 2 | Until after step 5, unbounded by step 2 | Until the branch merges; recorded with its trigger |
| Main risk | A second contract change | Stalls on another owner; grows into Track B | A long branch on `contracts/generation/go`; TypeScript test code that the accepted texts do not provide for |

## Three findings the research owner must not miss

1. **Step 2 of the expert's order is circular under the recorded decisions, and the expert was not told.** The package that P3 must run through does not exist. Its owner decided on 11 November that it "is created with its first real family, integers (bittype#163), once the unstable payload for integers exists". bittype#206 puts integers after step 1, which contains the migration. The question that would break the cycle was announced on bittype#196 on 12 November and never posted on bitspark-service#38. The same facts shorten the exposure the document feared: no other repository builds on the contract before integers. A stand-in proves the process boundary and a second reader. It does not prove their factory, and the record has to say so.

2. **Executed Go lowering fits the charter only if the prototype's vocabulary and its lowering facets both stay out of production packages.** BG:83 keeps bittype's vocabularies types-only, and a production mapping may not import a test vocabulary. Accepted rows (0015 rows 3 and 39) already put case 4 in bittype's test backend through a family facet, so the missing owner of value-level Go is a question of the plan, not of the advice. Today five of the eight things R29 lists cannot be expressed: a member is a string, and adapter arguments, requirements, value-level helpers and a call context have no form. #281's inventory lists top-level declarations only, with handles local to one output.

3. **"Thin versions of all six cases" is two different amounts of work.** Cases 3, 4, 5 and 6 fit on the six delivered families with four fixtures. Cases 1 and 2 need recursion, unions, nullable values, exact forms, maps and bytes (#164, #210, #211, #212, #209, #208). In real form that is Track B steps 2 to 9 placed before the migration, with the optional violation of bittype/ACX-009 open on `main` throughout. AC-10 allows that only as a recorded deferral for huge effort. The form that stays small tests the kernel's protocol with a fixture source and says plainly that it mocks the family.

## What I could not establish

- The board state of `codex-platform-status` (heartbeat, claims, inbox). I read issue and commit dates only; the last is 12 November 14:48Z.
- Whether bitspark-service's owner would accept a prototype branch or a pull request with factories for the six families ahead of integers. Nobody has asked (F6).
- Whether the service library's future factories will implement today's `Adapter` interface unchanged. #40 lists "`packages/api` (the adapter interface, the call context)" among what it touches.
- Whether `GroundDescriptor` yields the descriptor that a conversion test needs for `Page<Person>`, and whether horos TypeScript `diagnose` agrees with Go v0.10.2 on it. Nothing was run. horos 0.11.0 changes verdicts for recursive descriptors.
- Whether the hosted native job may fetch ontos and horos for a temporary module of generated glue. The installed-consumer job does fetch modules; the native job's temporary modules have no requirements today.
- Whether an existing test already places one output under two TypeScript organizations.
- The size of the work in days. I give counts of packages, targets, families, fixtures and repositories only.
- Not read: the bodies of research 0014 and 0015 beyond their "Applied" rows; `docs/naming.md` except one line; `LAYOUT.md` beyond the dependency table and the test compositions; the bodies of #163, #166, #209 to #214; pull request #260; bitspark-service research 0001; the four review reports and any other verification of this research.
