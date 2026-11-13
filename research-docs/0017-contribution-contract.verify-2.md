# Research 0017, verification 2: optional fields, identities, generics, native declarations (R11 to R25)

> **What this file is.** The report of one verification agent of
> [research 0017](0017-contribution-contract.md) on the recommendations
> R11 to R25 of the advice, checked against the code at `7e9dc34` on
> 13 November 2026. The agent did not write the research document. The
> report is a proposal and no decision. The decisions are the rows of the
> table "Applied" of the research document, and where a row and this
> report differ, the row governs. The design page
> [the contribution contract](../docs/contribution-contract.md) cites this
> report as verification report 2, and in its gate as V2.
>
> The text below is the report as it was written. Paths that begin with
> `C:/` are paths on the research owner's workstation and are kept as
> text. A path such as `platform/...` or `v2-inputs/` names a file that
> is not in this repository.

Read-only verification. Nothing was built, run, edited or posted. Two things touched state outside this file: `git fetch` of the two pull-request branches in the bittype checkout and of `origin` in the service library checkout (remote-tracking refs only, no working tree change), and the fetched input copies named below.

- **Code read:** bittype at `C:/Development/bitspark/platform/libraries/bittype`, `main` = `7e9dc34`. `git diff --stat 2daea52 HEAD` shows 12 files, all documents (`docs/naming.md` new, 874 lines; `docs/package-boundaries.md`, `docs/representation-families.md`, `docs/semantic-generation.md`, `docs/generation-planner.md`, `research-docs/0016-*`, two architecture cases, `CHANGELOG.md`, `LAYOUT.md`). **No Go file changed since `2daea52`**, so every code premise the expert was told still stands on `main`.
- **Open pull requests read:** #281 (`origin/feat/output-inventory-272`, head `c84ef50`), #282 (diff stat only), #260 (comments, second reading).
- **Other side:** bitspark-service `binding-contracts.md` at `c36722d`, issues #38 and #40, and `packages/api/src/wire.d.ts` in the local checkout at `c456bc7`.
- **Paths** below are relative to the bittype root unless they start with `platform/` or name another repository. "RF" is `docs/representation-families.md`, "BG" `docs/binding-generation.md`, "BC" the service library's binding contracts.
- Copies of the fetched inputs are in `C:/Development/bitspark/scratch-bittype-check/owner-1112/research-0017/v2-inputs/` (`binding-contracts.md`, `issue38.txt`, `issue40.txt`, `pr260.txt`, `issue166.txt`, `issue195.txt`, `issue196.txt`, `issue206.txt`, `issue207.txt`).

Terms: the **six delivered families** are text, booleans, lists, closed records, optional fields, generics (RF:341). The **six vertical cases** are RF:97-102.

---

## A. Optional fields (R11 to R14)

### Facts common to R11 to R14

**The finding, and who recorded it.** `docs/package-boundaries.md:156`, findings table, row 7: "Record decides optional fields itself in Go, Rust and C++ | Coupling | Confirmed; the contract change is designed | Open. […] Record still decides the representation of an optional field in all three targets". The reasoning at `:100-103`: "the same Record mapping decides how an optional field is represented. Rust builds `std::option::Option`, C++ builds `std::optional`, and Go wraps the type in a helper through a package variable. A consumer who wants another representation of optional fields must replace Record."

- It entered with commit `d23542a` (#223, 11 November 2026), written by this seat. The page says its "reasoning sections are unchanged since the proposal" (`:25-26`) and that the two independent reviews "read moves from the code" (`:23-24`): they reviewed package moves, not this judgment.
- The owner's request (`platform/docs/concepts/architectural-correctness.md:422-430`) names no optional field. The rule it applies is general: "A caller that needs whatever the consumer registered for a meaning states a demand" (`:455-456`). Whether presence is "a meaning the consumer registers" is the premise the finding assumes.
- `docs/architecture-cases/0009-package-boundaries.md:49-51`: "No instance is recorded here as a violation of AC-02 or AC-03."
- **The finding names three targets, not four.** TypeScript's `field?: T` is not in it.

So the expert is right: the finding is this seat's derivation, and no owner sentence decides the domain.

**What the accepted designs say about presence.** They support both readings.

| Text | Reads as |
| --- | --- |
| BG:37, owners table (accepted direction, 9 November 2026): "Each representation implementation (a target's structural profile, or a library's exact mapping) \| Native layout, presence rules, numeric representation, construction and access" | Neutral: an implementation owns presence, it does not say which one |
| BG:58, the plan's contents: "the rules for absence, presence, null and native `undefined`" | Neutral |
| RF:73, payload column (accepted design, 9 November 2026): "Record field layout and presence behaviour" as **one** item | Presence is the record's layout (D1) |
| RF:9-11, family list: "integers, text, bytes, lists, maps and sets, unions, optional and nullable fields, recursive types" and RF:341 "The delivered families (text, booleans, lists, closed records, optional fields, generics)" | Optional fields are a family of their own |
| RF:164-168: "TypeScript uses `field?: T \| null`, and the adapter separates an absent own property from a present `undefined` at run time." | Neutral on owner; fixes the run-time behaviour |
| bittype#207 item 1: "Each family is one representation contract with all seven parts"; item 9: "**Optional fields.** The TypeScript adapter separates an absent own property from a present `undefined` at run time. Go keeps `BittypeOptional[T]`" | One contract per family, so D1 contradicts #207 item 1 as written (an issue, not an owner rule) |

**The four targets as they are.**

| Target | Where presence is decided | Evidence |
| --- | --- | --- |
| TypeScript | In the record rule, as the member's modifier. No optional package exists | `transform/syntax/horos-to-typescript/mapping/record/go/record.go:55` `model.Field(f.Name, f.Optional, ty)`; target `transform/syntax/typescript/target/go/provider.go:134` `m.Attributes("name", data[0], "optional", data[1])` |
| Go | In the record rule, through a key the composition creates | `transform/syntax/horos-to-golang/mapping/record/go/record.go:70-72` `if f.Optional { ty = g.Helper(optional, ty) }`; `…/optional/go/optional.go:21-26` supplies `BittypeOptional` under that key |
| Rust | Inline in the record rule | `transform/syntax/horos-to-rust/mapping/record/go/record.go:53-55` `ty = model.External("std::option", "Option", ty)` |
| C++ | Inline in the record rule, and matched by name in the target | `transform/syntax/horos-to-cpp/mapping/record/go/record.go:53-55` `model.External("std", "optional", ty)`; `transform/syntax/cpp/target/go/provider.go:191` `if data[0] == "std" && data[1] == "optional"` adds the include |

All four record rules pass the flag into the `field` form (`model.Field(f.Name, f.Optional, ty)`), but only TypeScript's target reads it. Go (`golang/target/go/provider.go:93`), Rust (`:115`) and C++ (`:187`) write the name only. The flag is validated in every model (`model.go:70-73`) and is dead data in three of them.

**What each choice means.**

- **D1** (presence is the record's layout). TypeScript: nothing moves. Rust, C++: nothing moves; finding 7 is withdrawn for them. Go: `BittypeOptional` becomes a private part of the record contribution, the composition key disappears, and the `optional` package split of #227 is reversed. `docs/semantic-generation.md:159-162`, finding 7, ACX-009's limit bullet (`0009-package-boundaries.md:543-546`) and #207 item 1 are reworded. A consumer replaces the record contribution to get another presence layout (`Registry.Replace` by the selector `builtin/6:record/6:closed`).
- **D3** (the record delegates). TypeScript, Rust and C++ each gain an optional contribution package; all four record rules stop passing or acting on `f.Optional`; the kernel gains a selection for something that is no use (AC-08 justification needed); the conversion side needs a field-level protocol (R13) that reaches the service library's adapter interface.

---

### R11. Choose the domain explicitly; given delegation, D3

1. **Premises.** TRUE that D1 is legitimate and that the finding does not prove the other reading (quotes above). TRUE that D4 is what Go has today (`record.Library(optional g.HelperKey)`, `record.go:31`) and that it is interim: `docs/semantic-generation.md:159-162` "This form is interim: #196, item 6, changes how a record asks for the optional representation."
   D2, which the expert calls "not inherently dishonest", is **not expressible** in the kernel as it is: "Only the kernel mints occurrences" (`contracts/generation/go/occurrence.go:60-63`), from terms of the one admitted graph ("The kernel walks this graph and no other", `contracts/generation/go/source.go:20`). A source integration cannot present a use for a field position unless the horos graph gains a term for it.
2. **Drift.** None in code. `docs/package-boundaries.md` changed only in its status lines.
3. **Conflicts.** None with an owner rule. D3 contradicts RF:73 read literally ("Record field layout and presence behaviour" as one payload item); D1 contradicts #207 item 1 and the family list read literally. The expert did not show either text to be self-contradictory; the expert showed that neither decides. Whichever is chosen, one accepted sentence needs a clarifying edit by its owner (RF is this chain's page).
4. **Derivation.** No new abstraction in R11 itself. The case is delivered family 5 (optional fields) and vertical case 2 (`Page<Outcome>` with optional nullable fields). Under D1 the existing contract expresses it with no kernel capability.
5. **Change.** A recorded decision before the contract change. It cannot follow additively: under AC-10 a decision that sits in the wrong place moves first.
6. **Verdict: HOLDS**, with one adjustment to the wording the expert proposes: the decision record must (a) name TypeScript, which finding 7 omits, (b) say which accepted sentence it clarifies, and (c) be made with the service library, because the D3 branch changes the adapter interface (R13).
7. **Validation.** A dependency-and-behaviour test: register a second optional contribution for one target with the record contribution unchanged, and require the generated member to change. Red today (in Rust and C++ the only way is to replace `record/rust/1`); green under D3. Under D1 the mirror test: `grep -rn "f.Optional" transform/syntax/horos-to-*/mapping/record` is the only place presence is decided, and the Go `optional` package no longer exists.

### R12. The source integration issues or validates the subject

1. **Premises.** TRUE that an optional field is no use: horos carries it as a flag of the record's field, `transform/syntax/horos/source/go/read.go:77-79` "Each field is (horos/field, name, req|opt, type); field k's type is the operand at locator [k+2, 3]", and `:110` `Optional: atom(f.At(2)) == "opt"`.
   TRUE that a checked-child-use boundary exists and would be weakened by a free selector plus a free operand: `contracts/generation/go/description.go:158` `allowed = childOf(use, s.use.Use)`, error `errNotChild` "a capture must be a child the kernel minted for the mapped occurrence" (`:188`).
   **What the subject would be today.** `Reader.Fields` already returns three of the four parts: `type Field struct { Name string; Optional bool; Use g.Use[h.Language]; Origin ast.NodeRef }` (`read.go:25-30`). The fourth, the parent record occurrence, is the `UseView` the rule passed in. Two gaps:
   - The reader is the **contribution-facing** API. The kernel-facing `Interpretation` has only `Declaration`, `Expression`, `JudgeDeclaration`, `JudgeUse` (`contracts/generation/go/source.go:106-125`). Nothing lets the source issue a selector for a position. The kernel never sees the flag.
   - `Field.Origin` is the site of the field's **type term** (`read.go:110` `Origin: r.s.sites[term]`), not of the field. Provenance of the field itself is not available.
2. **Drift.** None.
3. **Conflicts.** None. Consistent with `transform/syntax/horos/source/go/provider.go:11-14`: the reader "supplies no navigation: declarations, binders, actuals and children come from the kernel's views."
4. **Derivation.** Mostly existing concepts: the kernel already mints the child use at locator `[k+2, 3]` and already checks `childOf`. The new part is one thing: a source-owned selector for "this child, in the role it plays in its parent". Needed by delivered family 5 only under D3. The existing contract cannot express it (selection is always by the use's own selector, `planning.go:450-465`).
5. **Change.** A kernel-facing source capability plus a request form. In the one contract change if D3 is chosen; void under D1.
6. **Verdict: NEEDS ADAPTATION.** Adjustment: define the subject from what exists, (parent use, the kernel-minted operand use, a selector the source interpretation issues for that pair), and keep the field's name and the flag out of the kernel. Record that D2 is ruled out by the kernel's minting rule, not by taste.
7. **Validation.** A fixture contribution requests the optional-field representation with an operand that is a child of another record, or with a required field's operand. Red today (no such request exists; the nearest, `Mapped`, accepts any child regardless of role). Green: both fail with `ErrScope`, and the same request for a real optional field of the mapped record succeeds.

### R13. A member-fragment vocabulary plus a contextual field-conversion protocol

1. **Premises.**
   **(a) "A TypeScript optional field cannot be converted from its retrieved value alone": TRUE**, and already an accepted statement: RF:167-168 "the adapter separates an absent own property from a present `undefined` at run time", repeated as #207 item 9. Reading `obj.next` gives `undefined` in both states, so the test needs the containing object.
   Checked against the first slice's Go side: `BittypeOptional[T]` has `Present` and `Value` (`optional.go:25`), so the **native** state is in the retrieved value. But the **canonical** side is contextual in every target: horos has no optional type, only `req|opt` on a field, so the conversion of one field maps a native field state to "absent, or one canonical value", and only the record can omit or include the field. The expert's claim is therefore stronger than stated: no target has a standalone value adapter for an optional field. This also means BG's laws, written "For a type T and representation profile p" (BG:121), cannot qualify an optional-field family by itself; its law is a pair of maps on field states.
   Checked against BC: C2 "Each family declares, for each native value it reads, one of: it **copies** […], it **retains** a reference, or it needs the value **stable while reading**." A field protocol that tests presence reads the record's object, so it must declare that read. Nothing in C1 to C5 or K1 to K4 forbids a contextual protocol.
   **(b) "A target-owned member-fragment vocabulary shared by the record and optional contributions": PARTLY.** A member vocabulary exists in every target's native model and is target-owned: `FormField = "field"` with `Field(name string, optional bool, ty)` (`transform/syntax/{golang,typescript,rust,cpp}/model/go/model.go:21, 43-45`), issued by the package-private `forms` issuer (`:12`). Go's optional contribution already imports that package and calls `model.Field` (`optional.go:10, 25`), so an optional contribution can use it without importing the record contribution. What does not exist is a **fragment**: the record form requires every argument to be exactly one `field` (`model.go:74-86`), so a field pair or an omitted member cannot be answered. And the member's **name** is the record's decision (`golang/…/record.go:62` `naming.Exported(f.Name)`; `docs/naming.md` owners table: "A member's name and accessibility | The declaration, so the mapping of that pair"), so the request must carry the name.
2. **Drift.** `docs/naming.md` (13 November) fixes member names with the declaration. #281 does not touch members.
3. **Conflicts.** With accepted wording on the other side, which the expert did not address:
   - bitspark-service#38, decision of 11 November 2026: "Each family is its own module. A family reaches other families only through its adapter slots." A field protocol is not an adapter slot as that interface stands: `Adapter<T>` is `toWire(value: T, context?)` and `fromWire(value, context?)` (`platform/libraries/service/packages/api/src/wire.d.ts:33-40`), and presence today is the record adapter's own option (`recordAdapter(properties: … { adapter, optional?: boolean })`, `:60`).
   - BG:94 "Value adapter facet | Run-time factory reference and typed argument-adapter slots, converting the value inside a message" and RF:40 "child adapter slots".
   These need a clarifying amendment by their owners (bittype for BG and RF; `codex-platform-status` for #38 and the adapter interface). The expert showed no contradiction in them; the need follows from RF:167-168 once optional fields are their own contribution.
4. **Derivation.** The native half maps onto the existing `field` form plus, where needed, one new form in a target's own model. The conversion half is a new product shape. Needed by delivered family 5 on TypeScript in the first slice (#207 item 9) and by vertical case 2. The existing contract cannot express it under D3; under D1 the record factory does it alone.
5. **Change.** Native half: in the one contract change (what a request may return). Conversion half: its shape must be settled before the contract is final, and it is delivered in the service library's run-time package. It cannot follow additively without reopening the record factory.
6. **Verdict: NEEDS ADAPTATION.** Adjustments: (i) state the protocol on field states for all targets, not as a TypeScript peculiarity; (ii) the request carries the member name the record chose; (iii) decide the slot kind with the service library before fixing it; (iv) remove the unread `optional` flag from the Go, Rust and C++ `field` forms when presence moves, or record why it stays.
7. **Validation.** Four-state vector for one optional text field on TypeScript: absent, present value, present `undefined`, and (after #211) present null, each asserted separately in both directions. Red today (no conversion exists; with `recordAdapter` the states collapse). Structural half: a test that the record factory's module imports no optional module and contains no presence test.

### R14. State the limit; contextual conversion is a product shape

1. **Premises.** TRUE as a statement of limit: one request per field cannot coordinate a layout shared by several fields. No delivered or listed family uses a shared presence layout: RF:164-168 names only `Optional[Nullable[T]]`, a pointer form and `field?: T | null`. TRUE that not every contribution has a standalone value adapter (R13, premise a).
2. **Drift.** None.
3. **Conflicts.** Same wording points as R13 (BG:94, RF:40, #38).
4. **Derivation.** Adds no abstraction; it bounds one. No delivered family or vertical case needs record-wide coordination.
5. **Change.** A sentence in the contract's applicability and an explicit refusal (AC-10: "explicitly refuses the cases it cannot yet check"). In the one change.
6. **Verdict: HOLDS.** It feeds the extension-domain list of R33.
7. **Validation.** A fixture optional contribution that tries to return two members for one request where the target's fragment admits one. Red today (not expressible). Green: a classified refusal naming the field, never a silently dropped member.

---

## B. Contribution identity (R15 to R17)

### Facts common to R15 to R17

- **What identifies a rule today.** A string unique inside one description (`contracts/generation/go/freeze.go:64` `ids[b.ID]`), such as `"record/golang/1"`, `"list/typescript/1"`, `"declarations"`. A registry counter: `type MappingRef struct { owner *Registry; id uint64 }` (`transform/syntax/shared/generation/go/registry.go:19-22`, issued at `:126-128`). The planner's name is their concatenation, `fmt.Sprintf("library-%d", e.ref.id)` (`:242`) plus `"/" + b.ID` (`planning.go:143`), kept only in the debug `Step.Binding` (`planning.go:33-36`). A replacement names `BindingRef{Provider MappingRef; ID string}` (`registry.go:23-30`) and walks `Description.Bindings` only (`:252`), so it reaches use rules and neither recipes nor helpers.
- **Selector keys are source identities, not contribution identities:** `builtin/<n>:<constructor>/<n>:<discriminator>` and `"declaration/"+modulePin+"/"+name` (`transform/syntax/horos/source/go/identity.go:15-27`), with the comment "These identities describe Horos grammar, never a selected native mapping."
- **No publisher concept exists.** `grep -ri publisher` finds only `source.Publisher`, the Files publishing type. The nearest accepted concept is the library: BG:87 "A library publishes, for one declaration and one supported target representation"; BG:101 "Contributions belong to their library's owner"; RF:60-62 "The representation owner owns its factories' semantic contract and qualification". None has an identity form.
- **What a backend keys on.** Not the contribution identity. BG:94: "Run-time factory reference and typed argument-adapter slots"; RF:49-50: "a plan entry names the family's run-time factory"; bitspark-service#38: "The backend only emits static imports of it into generated consumer code"; #196 item 1: "the runtime factory reference (module and export)". The submitted document's sentence "A backend in another process finds the factory by it" (section 6.4) is weaker than the accepted texts: the factory is found by its module and export. The contribution identity is what the plan records as selected (RF:76 "Selected contribution and profile identities"), what duplicates and version checks compare, and what a root facet recognises.

### R15. `ContributionID = (PublisherID, PublisherLocalKey)`, `ContributionRef = (ContributionID, revision)`

1. **Premises.** TRUE that the publisher is not the source's module pin: `Registry.Replace` exists exactly because two contributions answer one selector (`registry.go:131-145`, "Explicit replacements refer to both original and replacement offers", `:244`), and the selector already contains the pin, so the pin cannot tell them apart. TRUE that it is not the profile: every builtin states the same `structural-representation/1` (`transform/syntax/horos/mapping/go/mapping.go:28`). TRUE that I1's names are unique per registry only.
   PARTLY for the stated risk. The collision the submitted document feared ("the first wrong factory selected by a colliding name") does not arise if entries carry factory references, as the accepted texts say.
2. **Drift.** `docs/naming.md` adds a neighbouring identity, "the policy's identity", stated by the consumer's composition. It does not define a contribution identity.
3. **Conflicts.** None with an accepted design: RF:76 gives "Selected contribution and profile identities" with no form. With "Derive before extending": a "trusted composition lock" and a "manifest" are new artifacts, and no delivered family or vertical case needs them. Today the binding of an identity to an implementation is the registration itself (`RegisterMapping`), in one process.
4. **Derivation.** "Publisher" maps onto the existing concept "the library that publishes" (BG:87). The gap is an identity form, not a concept. Vertical case 6 (an independently registered family) and case 3 (the external `Cell<T>`) need two owners kept apart; the existing contract cannot express it (`library-<counter>` is neither stable nor owned).
5. **Change.** The identity fields a contribution states; `MappingRef` stays a registry-local handle. In the one change (AC-10: "its inputs and identities have their intended shape").
6. **Verdict: NEEDS ADAPTATION.** Adjustments: name the first part after the existing concept (the publishing library) and give it one form; drop the lock and manifest from the contract change and record them as deferred with the trigger "a plan is consumed by a process that did not perform the registration and must check which implementation it loaded"; correct the premise about factory lookup.
7. **Validation.** Two fixture libraries from different owners both state the local key `list` for one target. Red today in two ways: registered under one profile they fail as an ambiguous selector with no owner in the message, and their `Step.Binding` values differ only by `library-<n>`. Green: both register, the plan records two distinct identities that survive re-registration in another order, and a duplicate of one full reference fails.

### R16. Keep the five version axes separate from the identity

1. **Premises.** TRUE. RF:88-91: "**Five version axes:** the envelope protocol; the family payload schema; the representation's semantic identity; generated-code/runtime compatibility; source semantics." Today one of them is folded into the name: `"list/golang/1"`. No document says what that `1` versions (`grep "golang/1" docs/` finds only examples).
2. **Drift.** None.
3. **Conflicts.** None. The profile identity `structural-representation/1` also carries a version and is pinned by `TestProfileIdentitiesArePinned` (`transform/syntax/horos/mapping/go/mapping_test.go:18-25`); it is the profile's own and is not one of the five axes.
4. **Derivation.** No new abstraction; it applies an accepted row: research 0015, Applied, row 10 (HOLDS), "**Separate version axes:** envelope protocol, family payload schema, representation semantics, generated-code/runtime compatibility, source semantics".
5. **Change.** The trailing `/1` leaves the local key. In the one change. One test reads the string: `integration/go/horos_generation_profile_test.go:324`.
6. **Verdict: HOLDS.**
7. **Validation.** A plan fixture with one axis raised at a time is refused by the reader with that axis named, five tests (bitspark-service#40 item 2 asks for the same on the TypeScript side). Red today: no axis exists outside the name.

### R17. A family contract reference that spans targets, for `Service`'s recognition of records

1. **Premises.** TRUE that `Service` must recognise a family: BC, section 2: "A field is either a namespace, which is a field the plan selects as the structural record family, or a leaf whose selected contribution supplies an operation facet". TRUE that today's names are per target (`record/typescript/1`, `record/golang/1`).
   **Not needed by the first slice.** BC status line: "Nothing here is implemented yet". `Service<Ops>` is bitspark-service#39; the first slice is "`Page<Person>` carried end to end" in the Go completeness test backend (bittype#206). No root facet exists in either repository.
2. **Drift.** PR #260 (open) would add the sentence that the root facet "reads from the plan whether a field was selected as the structural record family"; its second reading lists that under "Sentences that decide more than they confirm", item 2.
3. **Conflicts.** None. It decides something for the service library, so it is theirs to agree.
4. **Derivation.** A new identity kind. The accepted axis "the representation's semantic identity" (RF:89) may already be it: if that identity names a contract (what the representation means) and not only a number, the root facet recognises it and no further kind is needed. BC's own wording ("the structural record family") is a single named case; a second profile's record contribution is the concrete case that would show whether recognition is by contract or by name, and none exists.
5. **Change.** Nothing in the first slice. But R15's "RepresentationSemanticRevision" fixes the shape of that axis in the one change, so whether it is `(contract, revision)` or a bare revision must be chosen then.
6. **Verdict: NEEDS ADAPTATION, with the root-facet part a candidate for DEFERRED.** Adjustment: do not add a separate "family contract reference"; decide in R15 that the semantic axis names a contract that may be shared by the contributions of several targets. Trigger for the rest: bitspark-service#39 writes the first root facet (`gh issue view 39 --repo Bitspark/bitspark-service` closed, or `packages/service-contracts` declares `Service`).
7. **Validation.** Later, in the service library: `Service<{ ns: { f: Function<…> } }>` composes with the namespace record selected from a second, independently registered record contribution that states the same semantic contract, and fails S2 with one that states another. Red today: nothing to run.

---

## C. Entries, nodes and the memo (R18, R19, R22)

### Facts common to R18, R19, R22

- RF:44-47: "The planner selects one immutable representation instance per use. The native type and the plan entry are both projections of that instance. No later pass reconstructs a decision from the source, so a native declaration and its conversion cannot disagree."
- RF:261-264: "The memo is keyed by the selected representation instance plus input identity, including the descriptor and generic environment. An object-only cache is unsound when one object appears under two mappings."
- BC, "Call context and conversion state": "**The memo key** is the selected representation instance (its descriptor and its generic environment) plus the input's identity." K3: "One object under two representations converts once per representation." K1: "One budget per message per side, shared by all nested conversions." K2: "A DAG whose expanded size exceeds the limit fails as exhaustion, even when its distinct objects fit." K4: "A cycle is detected through active traversal and fails as `cycle` unless the family supports it."
- Research 0014, Applied, row 15 (HOLDS): "Adapter cache keys include the selected contribution, generic actuals and profile."
- **The second reading of #260** (12 November 2026, finding 1, medium): "`docs/representation-families.md:44-45` says 'one immutable representation instance per use', and `:259` keys the memo by 'the selected representation instance'. The page now uses 'instance' in two senses." It suggests: "Two uses that select the same contribution and profile for the same descriptor and generic environment have the same run-time identity. […] Say that this identity is coarser than the per-use instance."
- **Code today.** No per-use record is kept. `Step` is five debug strings (`planning.go:33-36`), and `Program[T]` holds declarations only. A `Use` is "a graph term reached under one declaration along one structural path in one run. Two paths to a shared term are different occurrences" (`occurrence.go:60-63`): a run-local pointer, not a serializable identity. Sharing exists by construction for named declarations only (`p.declarations[d.entity]`, `planning.go:322`); two uses of `list<Person>` are bound twice.

Read literally, RF:44-47 and RF:261-264 together give every use its own memo key, so one object reached through two fields would convert twice and "Completed conversions preserve sharing" (BC) would fail across uses. The expert's N3 and the token of R22 remove that reading. They are consistent with K1 to K4, with row 15 of research 0014 and with the second reading. They need **a clarification of two accepted sentences**, not a reversal: RF:261-264 (owner: bittype, the holder of #195 and #196) and BC's memo-key bullet (owner: `codex-platform-status`). PR #260 is where the bittype wording is already under review and is not merged.

### R18. N3: use entries and representation nodes

1. **Premises.** TRUE (facts above). The envelope already lists both halves without naming the second: RF:74 "Root products and source-use identities" and RF:75 "Scoped binders and actual-argument environments".
2. **Drift.** None in code. #260's second reading postdates the submitted document's fixed point 16 only in its finding.
3. **Conflicts.** None once RF:44-47 is read as the use entry and RF:261-264 as the applied representation. The expert showed the two-senses problem (advice 2.2: "the runtime memo must not use the unique identity of that per-use record").
4. **Derivation.** "Representation node" is a new name for something the planner already has for declarations (one planned declaration per entity) and lacks for builtin forms. Vertical case 1 (shared DAG) needs a common identity across uses.
5. **Change.** Entry shapes are envelope (#195, explicitly unstable). The contract change must only fix that a selection yields one frozen decision that both a use entry and a node can refer to.
6. **Verdict: HOLDS.** One point the advice leaves open and the owner must settle: whether two structurally equal selections at different uses (two `list<Person>`) are one node. R19 says the kernel proves no equivalence, so sharing in the plan can only be by construction; equality across uses is then the run-time's interning (R22).
7. **Validation.** The second half of row 31's pending test, extended: one object placed in two fields of the same type converts once; the same object under two different mappings converts twice. Red today (no conversion; and under a per-use key the first half would fail).

### R19. The N2 tuple is insufficient as an equality key

1. **Premises.** PARTLY. Within one plan as the code is, (contribution, source identity, generic environment) nearly determines the decisions, because one registry cannot hold two rules for one selector (`planning.go:465-467`, `ErrAmbiguous`). It fails in the cases the expert lists once they exist, and two are already accepted:
   - **Configuration given to a contribution by the composition.** `docs/naming.md` D3: "The qualification policy is selected in the consumer's composition, as an argument of the declarations mapping of each pair". The same contribution with another policy yields other names. Go's `record.Library(key)` is the same pattern today.
   - **Position.** The record rule already decides by position (`record.go:49-51`, `!use.Body`), and an optional-field subject has no "source identity of the type" at all.
   And across plans: one process may load code generated from two compositions; a key that omits the composition's choices would merge them.
2. **Drift.** `docs/naming.md` D3 strengthens the premise.
3. **Conflicts.** None. It refines BC's parenthesis "(its descriptor and its generic environment)", which #260's first review already corrected: "The memo key includes the selected representation."
4. **Derivation.** No new abstraction.
5. **Change.** None in the contribution contract beyond R5 (decisions are recorded). It constrains R22.
6. **Verdict: HOLDS.**
7. **Validation.** R32's case "equal native types with different adapters stay distinct": two plans for one source, differing only in one replaced contribution, loaded into one run-time; the same input object converts under both without a memo hit. Red today: nothing to run.

### R22. `MemoKey = (direction, applied representation token, input identity)`

1. **Premises.** TRUE that the key is none of: a use identity, an entry number, a contribution identity, a descriptor alone. TRUE that not every node denotes a value conversion (R13). The direction component is already implicit: BC has "one call context per message per side", and each converts one way.
   Today's run-time already interns by applied arguments: research 0014 row 15 notes "platform today caches generics by argument-adapter objects (`genericAdapter`)" (`wire.d.ts:65`).
2. **Drift.** None.
3. **Conflicts.** Wording only, as above. Ownership: the memo is the service library's ("owned by this repository's adapter runtime", BC), tracked in bitspark-service#40 item 4.
4. **Derivation.** The token is the existing "selected representation instance" of BC, made precise. Vertical case 1 needs it.
5. **Change.** Nothing in bittype's contribution contract except that a contribution's facets may include one with no value token. The plan must carry the links from which the run-time derives the token; those are envelope.
6. **Verdict: NEEDS ADAPTATION.** Adjustments: route the definition to bitspark-service#40 and agree the wording of both sentences with its owner; drop the direction component or mark it implicit; state that the token is interned by the run-time from (factory reference, frozen payload, child tokens) so that no kernel-side equivalence is implied.
7. **Validation.** BC's own tests `K3 one object under two mappings` and `K2 a small DAG with a large expansion is exhausted`, plus R32's "repeated `Page<Person>` uses share a runtime token". All red today (unimplemented).

---

## D. Generics (R20, R21)

### Facts common to R20, R21

- **One open declaration.** `declare` plans a declaration once per entity (`planning.go:322-327`), with one `GenericConvention` kept per planned declaration (`:409`).
- **Applications are references with arguments.** `planning.go:515` `return Expression[T]{kind: DeclarationExpression, declaration: id, args: args}`. No record of the application is kept; not even a `Step` is appended for it.
- **A bound use is the binder's native expression.** `planning.go:444-449`, from `Generic.Bindings`.
- **What the convention expresses.** `description.go:38-45`: "Convention describes both sides of a generic, including reordered, omitted and additional native parameters. Bindings is indexed by source parameter; Arguments is indexed by target parameter." Checked at `planning.go:385`: `len(desc.Generic.Bindings) != len(view.Formals) || len(desc.Generic.Arguments) != len(desc.Generic.Parameters)`.
- **Accepted statement** (`docs/generation.md:45-48`): "Open generics remain open; ground specialization is not a prerequisite. A single generic convention relates source binders, native parameters and complete mapped actuals. It can reorder, omit or add native arguments while preserving their scope." And `docs/generation-composition.md:51`: "It does not specialize Page from its first use."
- **An omitted actual is never represented.** `planning.go:484-508`: only actuals an `Argument(i)` expression refers to are planned (`if needed[i] { actuals[i], err = p.use(a, f, depth+1) }`). An omitted one is admitted (`admit`, `:275-279`: "including actuals a native convention omits") and no representation is selected for it.
- **Selection ignores the arguments.** A declaration use selects by the callee's selector alone (`planning.go:450-453`, `:465`).

### R20. O1, explicit application nodes, separate native and adapter conventions

1. **Premises.** TRUE; O1 is what the code and `docs/generation.md` already do on the native side, and RF:92-94 asks for the rest: "Each generic use carries its scoped binder identities, native type arguments, adapter arguments and requirements".
   **Could native and adapter arguments be zipped by accident? Yes.** After `template` the application keeps only `args`, indexed by native parameter; the pairing with source binders is gone. A reader that took `args[i]` as "the argument for binder i" would be wrong for every convention that reorders, omits or adds, and one exists: `docs/generation-composition.md:51-52` "The independent Go fixture adds a native parameter computed from the complete mapped actual". And for an erased parameter there is nothing to zip at all, since no representation was selected (fact above). The expert's example "an exact native mapping might erase `T` while conversion still depends on `T`'s adapter" hits a real gap.
2. **Drift.** None.
3. **Conflicts.** None.
4. **Derivation.** "Application node" records something the planner computes and discards. The adapter convention is new and needed by vertical case 3 (`Cell<T>`) and case 4 (a generic adapter argument), and by delivered family 6 (#207 item 10: "the adapter selected for `T` travels with its native type argument").
5. **Change.** In the one change: the convention a declaration or exact-mapping contribution returns gains an adapter side indexed by source binder; the planner represents an omitted actual when a requested product needs its adapter and not for a types-only request (BG:99: "A native mapping without an adapter facet still supports types").
6. **Verdict: HOLDS.**
7. **Validation.** R32's "an erased native argument is retained for conversion": a fixture exact mapping of `Box<T>` whose native type omits `T`. Types-only generation succeeds with no rule for the actual (as today). A bindings request records an adapter argument for `T`, or fails naming the missing rule. Red today: the actual is never planned and nothing can ask for it.

### R21. O2 does not repair argument-dependent layout; decide whether specialization is in the domain

1. **Premises.** TRUE. One native declaration per entity (`planning.go:322`), so a fact about `Person` cannot change `Page`. RF:254-257: "Explicit specialization is represented as such (#166), never as ordinary substitution into an unchanged generic declaration."
   **Is any delivered family operand-sensitive inside an open generic? No.** `optional` of `T` is uniform in all four targets (`g.Helper(optional, ty)`, `Option<T>`, `std::optional<T>`, `field?: T`), and so is list. The first operand-sensitive families are later ones: RF:248 "Nullable `T` | The null representation cannot collide with a valid `T`", RF:246 maps, RF:247 unions.
   **Is application-specific specialization in the promised domain? By the accepted list, yes.** RF:349 delivery row "Explicit specialization | bittype | #166", and #206 Track B step 11 "Explicit specialization | #166 | step 1". RF:328 promises: "After them a family adds a contribution and changes no existing contract."
   **And #166 itself expects a contract change.** Its "Touches" lists `contracts/generation/go` and `transform/syntax/shared/generation/go`; its gap text says "A generic declaration is emitted once, so a mapping cannot change how one particular application is represented"; `docs/type-generation-plugins.md:500` says "Such application selectors and automatic monomorphization are deferred."
2. **Drift.** None.
3. **Conflicts.** RF:328 against #166's "Touches": as the contract stands and as the sketch of section 6.5 stands, #166 is a second change of the mapping contract. The expert did not see #166 and still named the issue.
4. **Derivation.** A selection subject for an application (callee plus the selections of its actuals) is a new input shape. Step 11 needs it; no delivered family and no vertical case does.
5. **Change.** Either the one change makes an application nameable as a subject of selection and refuses it everywhere for now, or specialization is written out of the "changes once" domain with #166 as the recorded second change.
6. **Verdict: HOLDS.** It must be decided before the contract change, and R33's domain list is where the answer is written.
7. **Validation.** #166's own verify fixture ("`Page<Person>` keeps its generic representation and `Page<Byte>` uses the library's specialized representation, in the same output") written now as an expected refusal: a contribution states a selection for the application `Page<Byte>` and planning refuses it with a classified "not yet supported" at that application. Red today: the request cannot be stated, and an exact rule for `Page` would capture every application.

---

## E. Native symbols (R23 to R25)

### Facts common to R23 to R25

- **What a handle binds** (#281, `contracts/generation/go/output.go:94-97` on the branch): `type DeclarationHandle struct { output *identity; index uint32 }`, "identifies one declaration of one sealed output. Only seal issues handles." `Resolve` "refuses a handle that this output did not issue with ErrScope. That includes the handle of another run of the same plan" (`:156-164`).
- **Why it cannot key a serialized plan** (`docs/naming.md:644-656`): "an output-issued handle cannot be the key of a serialized plan", derived from three sentences: a plan entry exists at planning, before any output; a handle belongs to one sealed output; the plan crosses a process boundary. In code the output's identity is an allocation (`owner := fresh()`, branch `output.go:241`), so neither a handle nor an output has a serializable form.
- **A plan-side identity already exists and is dropped.** `PlannedDeclaration.ID` is issued by the planner at reservation (`planning.go:309-317`; branch `:321-325`), before `Run`, and is what every `DeclarationExpression` and `ParameterExpression` refers to. In #281, `seal` receives the planned declarations in planned order beside the assembly's sites (`seal(…, declarations []planned, exposed []uint32, …)`, branch `:177`), then sorts the records by site (`:238`) and numbers the handles by sorted position (`:243`). The record keeps `Handle, Site, Role, Provenance, Values` (`:108-114`) and **no planned identity**. So the correspondence the expert asks for is available at exactly one point, inside `seal`, and is not published.
- **Today a consumer could recover it only partly:** by `Provenance` for a declaration that represents a source declaration, and not at all for a helper ("Empty for a helper, which represents none").
- #281's planner document: "Adding a field to the record is additive", and "A grouping that a later case needs enters as an issued value". `Values` are "what the owner of the declaration's description states"; a planned identity is the planner's own fact, like `Role` and `Provenance`.
- **What the Files formatter's result offers today.** A `source.Source[L]` (`contracts/format/go/generated.go:35-43`). The TypeScript formatter's map from generated to placed site is private (`format/files/typescript/provider/go/relocate.go:25-31`, field `moved`). `docs/naming.md:668-671` D6: "A consumer learns the realized name and the placed site from a reader on the formatter's result, keyed by the output's handles", and R14 there: "Nothing offers them before formatting." That reader is not delivered. Only TypeScript and horos have a generated-Files formatter (`format/files/` holds `horos` and `typescript`); Go and Rust have none.
- **Helpers today.** `type HelperKey struct { owner *identity; name string }`, `NewHelperKey(name)` returns `HelperKey{fresh(), name}` (`contracts/generation/go/expression.go:41-48`). Two keys of one name differ (`integration/go/horos_generation_golang_optional_test.go:55` expects `ErrUnsupported`). The planner shares one declaration per key (`helperIDs map[HelperKey]uint32`, `planning.go:685-695`); two suppliers of one key fail (`planning.go:149-151`, `ErrAmbiguous`); a missing supplier fails (`:681-684`, `ErrUnsupported`). Equal spellings are refused separately by the shared constructor (`docs/naming.md`, "The code today": `transform/syntax/shared/target/go/check.go:17-19`).

### R23. S2 with a checked realization manifest

1. **Premises.** TRUE that the placement contradiction is real and already recorded by this chain: `docs/naming.md:721-727`, E7: "With names realized by Files, a plan entry cannot hold a realized name or an import path. […] That is a fourth input to the backend and an amendment to an accepted sentence. […] It belongs to #195 and goes into research 0017". TRUE that the correspondence must come from generation and not from names (`docs/source.md:58-59`: "A spelling alone is insufficient").
   PARTLY for the manifest: its line "generation run / sealed output" names something with no serializable identity (fact above). What can cross a process boundary is the plan's digest and a table from logical symbol to realized symbol; the join through handles has to happen in the Go process that holds the output and the formatter's result.
2. **Drift.** Large and favourable. `docs/naming.md` now states the tension and defers it to this research; research 0016 row 10 is applied in #281 ("identity and site, role, provenance, and the naming intent of row 13 as an owner-issued value"); row 13 separates "identity, naming intent, realized spelling".
3. **Conflicts.** It amends accepted sentences, and the expert did show the contradiction (advice 2.5; the submitted document's second contradiction):
   - BG:23-27 "**A binding backend per target** consumes three inputs" and BG:72, the completeness test "from the plan, the source semantics and the selected contributions alone". Owner: bittype, #195.
   - RF:78 "Native symbols, imports, factory references, runtime requirements" in the envelope. Owner: bittype.
   - bitspark-service#40 item 9: "The backend builds every selected representation from the plan, the source semantics and the selected contributions alone." Owner: `codex-platform-status`. The advice does not mention this one.
   No conflict with research 0016 rows 10, 13, 14, 19, 34.
4. **Derivation.** The "logical symbol identity" maps onto the existing planned declaration identity; the "realization" onto D6's reader. New: a serializable form of the first, a published link from it to the handle, and the manifest. No delivered family needs an import (#206: "they need no import, so they do not wait for Go Files"); vertical case 3 (`Cell<T>` from a package) and bittype#163's Go half (`math/big`) do.
5. **Change.** In the one contract change: only that an entry refers to a native declaration by a plan-issued identity and never by a name. Additive afterwards: a planner-stated field on the inventory record (#281's own rule), the reader on the formatter's result (#273), the manifest and its serialized form (#195).
6. **Verdict: NEEDS ADAPTATION.** Adjustments: (i) use the planner's existing declaration identity and say what makes it stable (it follows root and traversal order, so it is stable for one plan and not across requests); (ii) publish the link as a planner-stated field of `OwnedDeclaration`, not as an issued value; (iii) replace "generation run / sealed output" in the manifest by what is serializable, and say that the handle join happens in process; (iv) take the amendment to bitspark-service#40 as well.
7. **Validation.** R32's "one semantic plan works under two placements": one plan, two TypeScript placements; the backend input for each resolves every logical symbol to a different module path and the same declaration. Then mix plan A with the realization of plan B and require a linking error that is distinct from "the plan is incomplete". Red today: the record has no planned identity, the formatter has no reader, and a helper cannot be matched at all.

### R24. Read the realization now; structured location, export and accessibility facts

1. **Premises.** TRUE and already this chain's design: `docs/naming.md`, "Three things that were one": "**Realized spelling** | The identifier that a language's Files stage publishes. In the first slices it is the spelled intent, once that language's rule has checked it"; D6 and R14 quoted above. TRUE that the import statement is derived where the glue is placed (`docs/source.md`: the Files formatter must "Choose the local spelling, qualification or alias, and construct the necessary import").
   FALSE for "now" on two of the first slice's targets: there is nothing to read for Go or Rust. No Go generated-Files formatter exists; the Go target still takes the package from generation (`transform/syntax/golang/target/go/provider.go:15, 79`, the open violation bittype/ACX-004), and the first slice's conversions run in the Go completeness test backend.
2. **Drift.** `docs/naming.md` is the drift; it agrees.
3. **Conflicts.** None. Research 0016 row 13a defers Go's unexported spelling with a trigger; that is compatible.
4. **Derivation.** No new abstraction beyond D6's reader.
5. **Change.** Nothing in the contribution contract. The reader is #273 for TypeScript, slice 1 of #161 for Go, slice 1 of #159 for Rust.
6. **Verdict: NEEDS ADAPTATION.** Adjustment: per target. TypeScript reads the realization from its first plan-consuming test. Go and Rust are a candidate for DEFERRED with the trigger "a generated-Files formatter for that language exists" (`ls format/files/` shows `golang`, respectively `rust`); until then the Go test backend states in its own record that it reads the spelled intent, as an open limitation and not as the interface.
7. **Validation.** A source-level check on the TypeScript backend fixture: no native name reaches emitted glue except through the realization reader. Then a mutant formatter that realizes one declaration under another spelling than its intent; the glue must follow it. Red today (no reader; a backend would have to copy the intent).

### R25. A contribution-owned declaration gets an owned logical key

1. **Premises.** TRUE that the spelling is not the deduplication key; that already holds (`HelperKey` is the key, and the name collision rule is separate). TRUE that a contribution needing "whichever shared helper composition selects" makes a demand; `g.Helper(key, args...)` is that demand today, with missing and ambiguous suppliers failing in the planner. What changes: the key stops being an allocation the composition creates and passes to two libraries, and becomes a value its owner states. Only a value can be recorded in a plan or matched in another process.
2. **Drift.** #281 gives a helper its own role and its owner's statements: `HelperDeclaration` "is demanded by a mapping through a helper key. It represents no source declaration", and `NativeDeclaration.Values`. Research 0016 row 13 (annotated 13 November): "the optional helper is realized exported because its contribution states it". `docs/naming.md`, owners table: "The chosen name of a declaration | The mapping: the declaration recipe for a source declaration, the helper recipe for a helper".
3. **Conflicts.** None. With R11: under D3 the Go helper becomes private to the optional contribution, so the shared-key case has no delivered instance left; under D1 it is private to the record contribution. Either way the composition-made key disappears.
4. **Derivation.** An existing concept (the helper key) re-keyed from identity to value. Needed by delivered family 5 on Go (`BittypeOptional`); the independent fixture has a second helper (`integration/go/testdata/generationfixture/fixture.go:259`, `g.NewHelperKey("independent/cursor")`).
5. **Change.** `HelperRecipe.Key` and `Helper(key, …)` in the mapping contract: in the one change. The key on the inventory record is additive, and it is what lets a helper be matched under R23.
6. **Verdict: HOLDS**, with the note that "relevant parameters" needs one sentence: the delivered helper is generic and shared across every operand, so no parameter enters its key today.
7. **Validation.** Two contributions each own a helper they both spell `BittypeOptional`. Green: two distinct symbols in the plan, and the name collision reported by the language's placed-scope check, not merged and not an ambiguity. And one helper demanded from two rules through its owned key appears once in the inventory with that key on its record. Red today: a key cannot be stated without the composition, and the record of a helper carries nothing that identifies it.

---

## Proposed verdicts

| # | Recommendation | Proposed verdict | In the one contract change? |
| --- | --- | --- | --- |
| R11 | Choose the presence domain explicitly; then D3 | HOLDS (record the choice with TypeScript named, the clarified sentence, and the service library's agreement) | Decision precedes it |
| R12 | The source issues or validates the subject | NEEDS ADAPTATION (subject = parent use, kernel-minted operand use, source-issued role selector; D2 is ruled out by the minting rule) | Yes, if D3 |
| R13 | Member fragment plus contextual field conversion | NEEDS ADAPTATION (field-state protocol for all targets; name in the request; slot kind agreed with the service library; drop the unread flag) | Native half yes; conversion shape before final |
| R14 | State the limit; contextual conversion is a shape | HOLDS | Yes, as a stated refusal |
| R15 | Publisher-owned contribution identity | NEEDS ADAPTATION (publisher = the publishing library; lock and manifest deferred; correct the factory-lookup premise) | Yes |
| R16 | Five axes apart from the identity | HOLDS | Yes |
| R17 | Family contract reference across targets | NEEDS ADAPTATION (fold into the semantic axis of R15); root-facet part candidate for DEFERRED, trigger bitspark-service#39 | Shape of the axis only |
| R18 | N3: use entries and representation nodes | HOLDS (clarify RF:44-47 against RF:261-264) | Notion yes; shapes are envelope |
| R19 | N2 tuple insufficient | HOLDS | No |
| R20 | O1, application nodes, separate conventions | HOLDS | Yes |
| R21 | O2 does not repair layout; decide specialization | HOLDS (decide against #166 before the change) | Decision precedes it |
| R22 | Memo key with an applied token | NEEDS ADAPTATION (route to bitspark-service#40; direction implicit; run-time interning) | No |
| R23 | S2 with a checked realization manifest | NEEDS ADAPTATION (existing planned identity; planner-stated inventory field; serializable manifest; amend bitspark-service#40 too) | Only "by plan-issued identity, never a name" |
| R24 | Read the realization now | NEEDS ADAPTATION (TypeScript now; Go and Rust candidate for DEFERRED, trigger: their generated-Files formatter exists) | No |
| R25 | Owned logical key for owned declarations | HOLDS | Yes |

No recommendation in this set is stale or invalid, and none is a candidate for rejection.

## The three findings the research owner must not miss

1. **The optional-field conversion is contextual in every target, and the decision reaches the service library.** horos has no optional type, so even Go's `BittypeOptional` converts to "absent, or a value" and only the record can omit the field. The accepted pages already require the TypeScript half (RF:167-168, #207 item 9). Under D3 this needs a field-level slot that today's adapter interface does not have (`Adapter.toWire(value)`, `wire.d.ts:33-40`) and that bitspark-service#38 excludes as worded ("only through its adapter slots"); BG's laws are per type and cannot qualify the family alone. Under D1 none of this is needed, #227's Go split is reversed and #207 item 1 is reworded. Finding 7 is this seat's derivation, names three targets, and the accepted texts support both readings: this is a real decision, to be made with `codex-platform-status`, before the contract is drawn.
2. **"Changes once" already fails for explicit specialization, and adapter arguments have no source today.** #166 is on the accepted family list and its own "Touches" names `contracts/generation/go`; selection keys on the callee alone (`planning.go:450-465`) and one native declaration exists per entity (`:322`). Separately, an actual that the native convention omits is never represented (`planning.go:484-508`), so "`T`'s selected adapter travels with its native type" has nothing to travel with for an erased parameter, and the retained application is indexed by native parameter only. Both belong in the contract change or in a written exclusion.
3. **The plan-to-handle correspondence exists at one point and is thrown away, and nothing on the output side is serializable.** #281's `seal` holds planned order and sites together, then sorts and renumbers; the inventory record carries no planned identity; a helper has no provenance to match by; handles and outputs are allocations; the formatter's result has no reader; Go and Rust have no Files stage. The amendment of BG:23-27 and BG:72 that `docs/naming.md` E7 anticipated is confirmed by the expert, and it must also reach bitspark-service#40 item 9, which repeats the three-input sentence.

Also relevant, already known to the chain: RF uses "instance" in two senses (found by the second reading of #260, finding 1). R18 and R22 depend on clarifying RF:261-264 and BC's memo-key bullet; PR #260 is open and is the place for the bittype half.

## What I could not establish by reading

- Whether the service library would accept a field-conversion slot, or prefers presence in its record factory. No statement exists; bitspark-service#39 and #47 were not read.
- Whether the planner's declaration identity is deterministic across two plans of the same request. By reading, reservation order follows root order and expression traversal over slices; nothing was run.
- What the reader on the formatter's result will offer. bittype#272, #273 and #274 were read only through `docs/naming.md`, research 0016's Applied rows and #281's diff, not their issue texts.
- Whether today's output stays byte-identical if the `optional` flag leaves the Go, Rust and C++ `field` forms. Those targets do not read it, so it should; not run.
- Research 0014 and 0015: the Applied tables were searched by keyword (memo, instance, identity, generic, optional, specialization, duplicate); their advice texts were not re-read in full.
- #282 was read by diff stat only. It separates unsupported coverage from invalid input in the Files stage, which is where R23's "linking or input error" would be classified.
- Whether any consumer outside the repository reads `Step.Binding`. Inside bittype one test does (`integration/go/horos_generation_profile_test.go:324`).
- The 21 rule packages were not counted; the record and optional mappings, the list and text mappings of Go, the declaration recipe and all four models and targets were read.
