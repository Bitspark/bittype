# Research 0017: the report of the prototype

> **What this file is.** The report of the prototype of the contribution
> contract: step 2 of the order of work of
> [research 0017](0017-contribution-contract.md), row 31. The prototype is
> the branch `prototype/contribution-contract-196` of this repository at
> `6080420`, with the branch point `2a29674`. The agent that built the
> prototype wrote this report on 13 November 2026.
>
> **Its standing.** Evidence and proposals. It is no decision. The decisions
> on it are O27 to O39 in the
> [comment of 13 November 2026](https://github.com/Bitspark/bittype/issues/196#issuecomment-6844431866)
> on #196, and the design page
> [the contribution contract](../docs/contribution-contract.md), which states
> them and cites this report by its section and by its finding, F1 to F27.
> Where a decision and this report differ, the decision governs.
>
> **The branch is thrown away after the migration.** It opens no pull request
> and never merges. Its log is `PROTOTYPE.md` at its root.
>
> The text below is the report as it was written, character for character. It
> was written as plain text and stands in a fenced block, so that its angle
> brackets and its lines stay as they are. Paths that begin with `C:/` are
> paths on the research owner's workstation and are kept as text. The files
> it names under "EVIDENCE FILES" are on that workstation and not in this
> repository. Line numbers of accepted pages in it are those at `2a29674`.

```text
PROTOTYPE OF THE CONTRIBUTION CONTRACT: REPORT
bittype#196, research 0017 row 31, step 2 of the order of work. 13 November 2026.

1. THE BRANCH

Repository Bitspark/bittype, branch prototype/contribution-contract-196, head 6080420, pushed and read back from the remote. Branch point 2a29674. Twelve commits: one for each stage, one for the reader's composition hook and two for the log. The gate runs were made on 4442047. The one commit after it changes the log and one comment. No pull request was opened. The worktree is C:/Development/bitspark/platform-wt/bittype-196-prototype-20261113. The log is PROTOTYPE.md at the root of the branch. The evidence is in C:/Development/bitspark/scratch-bittype-check/owner-1112/delivery-0017/prototype/.

All nine stages are built. The design could be built. Twenty-four findings came out of it (F1 to F27; F20 to F22 of the log are folded into others). None stopped a stage. Five concern a sentence of the page that could not be built or shown as written (F5, F8, F11, F15, F16). Two things the gate asks for were not run: the test of row 24 (no reader on the formatter's result exists on main) and the test of row 10 over all 21 packages (ten of the 21 are moved).

Two decisions about the branch itself. First, the mapping contract of main stays beside the new contract: generation.Plan and MappingLibrary are untouched, and generation.PlanContributions stands beside them. Rust, C++ and the TypeScript boolean package stay on the mapping contract. This was cheaper than an adapter and cheaper than taking them out of the build. Second, each moved mapping package keeps its old Library and gains a file contribution.go, so that today's tests run unchanged and one test compares the two contracts on one module. Moved: for Go text, boolean, list, record, declarations and optional. For TypeScript text, list, record, declarations and a new package optional. Go's boolean is moved beyond the brief's list, because the payload schema of the declaration contribution has a boolean field and the payload types are generated with the Go contributions.

2. EVIDENCE FILES

gate-run-prototype.txt: one verbose run of all prototype tests and TestProductionArchitecture. 41 tests pass, with 119 GREEN lines and 25 RED lines. A test logs GREEN for the kernel as built and RED for a kernel with one property switched off (the switch generation.PrototypeMutants, an instrument of the branch).
gate-run-all.txt: go test ./... on the head. Every package passes except seven Rust tests of integration/go that need the rustc helper. They fail at the branch point too.
gate-red-at-branch-point.txt: the prototype's tests do not build against 2a29674. A test that has no mutant is red in that sense only.
stage1-run1.txt to stage9-run2.txt: the runs of each stage, including the first runs that failed.
architecture-red-stage1.txt: TestProductionArchitecture before the new roles.
stage5-lowered.go.txt, stage5-plan.json, stage5-snapshot.json, stage7-plan.json, stage7-snapshot.json: the lowered Go, two serialized plans and two snapshots.
stage9-diff-stat.txt: the commit of the contribution of line 7 adds one file.

3. PER STAGE

All test names begin with TestPrototype. They are in integration/go/prototype_*_test.go.

Stage 1, the staged rule and the frozen record. Built.
- Stage1NativeDeclarationsDoNotChange. Shows: the moved contributions generate, through name, resolve, decide, freeze and project, the same Go text as the mapping contract for today's pagination module, and the same TypeScript text for that module without booleans. Does not show: Rust, C++, TypeScript booleans, or equality of the origins inside the native AST. No mutant.
- Stage1ProjectorCannotReadTheSource. Shows: the projector's capability mentions no use, declaration, binder, selector or reader (by reflection); three attempts to read the source do not type-check; the same contract test is red on the callback of the mapping contract. Does not show: that a contribution cannot hand source data to its projector through state of its own package. The page says so itself.
- Stage1FrozenValueSurvivesMutation. Shows: a value mutated after the rule returned changes neither the native type nor the serialized plan. Red under projector-gets-returned-value and freeze-at-serialization. Does not show: immutability against a codec that keeps the value.
- Stage1DecisionTheSourceDoesNotDetermine. Shows: with a fixture whose member names come from its configuration, the plan entry and the native type agree, and an entry rebuilt from the source is detected. Does not show: that no second pass can be written. The mutant is a pass of the test, not of the kernel.
- PayloadBootstrap. Shows: for eleven vocabulary packages the committed module and pin are what the schema source lowers to, the committed Go type is what bittype generates with the prototype's own Go contributions, and every listed value round-trips through the hand-written codec. Does not show: TypeScript payload types (none are generated), or a field the existing families cannot type. It was not shown red.

Stage 2, evidence. Built.
- Stage2FixtureFamilyDecidesOnEvidence. Shows: a rule gets a typed answer about what was selected for its operand before it decides: []string where the dependency declares the fact, [][]string where nothing is stated. Does not show: that the fact is true.
- Stage2GenericAndWrittenOutAgree. Shows: with a declared binder requirement, Page<text> and the written-out record have the same native type; the open declaration records the requirement and the application its discharge. Red for a rule that decides differently inside the generic. Does not show: that the kernel refuses such a rule (F5).
- Stage2FiveOutcomes. Shows: proved, refuted, unknown, exhausted and deferred each reach a rule as their own answer. Does not show: a check that reaches a limit. Exhausted is stated, nothing computes.
- Stage2DeclarationPolicyAsksAboutItsBody. Shows: the declaration contribution decides nominal or alias from the answer about its body and reads no payload of the record.
- Stage2RevisionIsIdentityAndNoPayload. Shows: two representation revisions give equal payloads and different nodes, also for the node that depends on them.
- Stage2CycleAndRecursionAreToldApart. Shows: Rec = Page<Rec> and Loop = probe(Loop) are refused as unsupported with two reasons. Runs on the independent fixture source (F7). Does not show: anything of recursion beyond its refusal.
- Stage2AnswersAreBoundToTheirDependency. Shows: an answer for another dependency, another generic environment, another representation revision or another meaning revision of the dependency is rejected. Red under four mutants. Does not show: that a rule cannot act on what it remembers without presenting it.

Stage 3, identities and refusals. Built.
- Stage3PublisherAndLocalKey. Shows: one key under two publishers registers; one publisher with one key twice fails, also for a second target. Does not show: that a publisher's name is its owner's.
- Stage3OneMeaningTwoRepresentationRevisions. Shows: for five families Go and TypeScript carry one meaning and two representation revisions; a reference that lacks either part is refused. Does not show: that the meaning is one identity and not two equal spellings (F2).
- Stage3RefusalValues. Shows: missing and ambiguous are two values; an unbuilt family is "missing contribution"; a specialization is refused as unsupported with the application named while another application plans; five refusals of moved rules match exactly one value each; an unexpected error of a rule matches none. Does not show: the refusals of the eleven unmoved packages.
- Stage3OwnedDeclarationsAreKeyedByOwner. Shows: two contributions that spell a helper alike get two symbols. Red under owned-shared-by-spelling. Does not show: a shared helper request (F10).

Stage 4, the plan and its serialized form. Built.
- Stage4TwoUsesOneNode. Shows: two uses of Page<Person> are two use entries with two sites on one node; the open declaration is one node; the application links its binder to the node of Person.
- Stage4OmittedArgumentIsKeptForConversion. Shows: asked for bindings, an application of a declaration that omits its parameter natively has no native and one adapter argument, and two such applications with equal native types are two nodes; asked for types alone the omitted argument is not planned. Red under two mutants. Does not show: an adapter.
- Stage4MemoKeyFromTheSerializedPlan. Shows: a reader that only parsed the serialized plan derives equal keys for equal nodes and another key with another selected dependency. Does not show: the memo.
- Stage4VersionAxesAlone. Shows: an unknown envelope protocol, payload schema, meaning revision, representation revision and source semantics are each refused under their own name. Does not show: the fourth axis (stage 7).
- Stage4RealizationManifest. Shows: the inventory record carries the planner's identity; one plan links with the manifests of two TypeScript placements and resolves each symbol to another module path and the same name; the manifest of another output is a linking error. Red under link-ignores-plan-digest, which then reports an incomplete plan. Does not show: a reader on the formatter's result (F11).

Stage 5, facets, Go handlers and executed Go lowering. Built, with the limit of F14.
- Stage5LoweredGoCompilesAndRuns. Shows: a fresh process lowers Go conversion for Page[Person] from the serialized plan, the snapshot and the facets; the code compiles, converts five values there and back and refuses three outside the type; it has a typed construction, member access, the owned helper and a generic adapter argument. Does not show: a typed expression vocabulary, or the conversion laws.
- Stage5MemberByReferenceNotByString. Shows: with member names from a fixture's configuration, the lowering by member reference compiles and runs and the lowering by string does not compile.
- Stage5FacetsBelongAndMatch, Stage5ReferencesAreChecked, Stage5ExactPayloadPassedValidation: guarantees 1, 2, 3 and 6, in section 5.

Stage 6, optional fields by delegation. Built in both targets.
- Stage6FourStatesOfAnOptionalField. Shows: in lowered Go and in the stand-in, absent and present with a value convert there and back as themselves; present with null is refused; in the stand-in a present undefined is refused, and that vector is red when the protocol reads it as absence. Does not show: a null that a type admits, the service library's slot, independent authorship.
- Stage6FieldPositionKeepsTheScopeCheck. Shows: a request with the operand of another record, or with a role issued for another pair, fails. Red under field-ignores-parent.
- Stage6MemberFragmentOutsideTheVocabulary. Shows: a fragment that is no member is refused as unsupported with the member vocabulary as owner; a target without a member vocabulary refuses a field position; a fragment of two members is admitted.

Stage 7, a second reader in another process. Built.
- Stage7StandinRoundTripsPageOfPerson. Shows: the TypeScript stand-in, run with node, converts Page<Person>, a nested page and a list there and back from the plan and the snapshot alone. Does not show: the service library's factory, its budgets, call context and memo, or independent authorship.
- Stage7StandinRefusesUnknownVersions. Shows: five axes refused under their own names, the runtime compatibility among them.
- Stage7ReadersNeedTheMatchingSnapshot. Shows: without the snapshot the stand-in and the Go backend fail explicitly; a snapshot of another pin is a linking error in both.
- Stage7StandinDerivesMemoKeys. Shows: equal keys for two uses, another key with another dependency, and the same key as the Go reader derives.

Stage 8, the fixtures. Built. Case 1 is not built.
- Stage8Case2PageOfOutcomeWithAMockFamily. Case 2 at the kernel's level, through the independent fixture source, with a MOCK family (the contribution mock-optional-nullable). Shows: the kernel plans a generic applied to a record whose every field is delegated, through a source that is not horos. Does not show: null against absence, discriminators, anything about the family the mock stands for, or conversion.
- Stage8Case3ExactGenericWithAnOmittedParameter. Shows: an exact generic omits its parameter natively, its argument is linked for conversion, its value facet is found. Does not show: an additional native parameter.
- Stage8Case5GenericRequirementAcceptedAndRefused. Shows: a binder requirement is discharged by a proved fact, refused as unsatisfied by a refuted one and as unresolved where nothing is stated.
- Case 4 is Stage5LoweredGoCompilesAndRuns. Case 6 is stage 9. The identity cases of row 32 are the four tests of stage 4.

Stage 9, the gate and line 7 in its small form. Built.
- Stage9ReplacedContributionLeavesTheOthers. Shows: after the replacement the payloads of the 9 other selections are byte-identical; 5 frozen records that do not depend on the replaced contribution are identical; 11 that depend on it are other nodes (O26); with both registered the domain is ambiguous until the consumer replaces explicitly. No mutant.
- Stage9PageOfPersonRegeneratesAndConverts. Shows: Page<Person> regenerates with the new contribution and the conversion lowered by the unchanged backend converts through its field-state protocol. Does not show: another owner, or TypeScript.
- Stage9OutOfTreeModule. Shows: a Go module outside the repository authors and registers a contribution and a codec with public packages only; a file that needs an internal package does not build. Does not show: a generation with it, or how its author gets the schema artifact (F25).
- Stage9ExtensionDomainRefusals. Shows: a contribution at a field position cannot request a sibling's operand.

4. THE GATE TABLE

Line 1, native projection. Stage1ProjectorCannotReadTheSource (row 4), Stage1FrozenValueSurvivesMutation (row 28), Stage1DecisionTheSourceDoesNotDetermine (rows 6 and 39). Green. Red under two mutants and against a pass that derives the entry from the source.

Line 2, dependency composition. Stage2FixtureFamilyDecidesOnEvidence (row 3), Stage2GenericAndWrittenOutAgree (rows 7 and 8), Stage2RevisionIsIdentityAndNoPayload (row 5), Stage2CycleAndRecursionAreToldApart (row 9), Stage3RefusalValues (rows 10 and 21), Stage6FieldPositionKeepsTheScopeCheck (row 12), Stage3PublisherAndLocalKey (row 15), Stage4TwoUsesOneNode (row 18), Stage4OmittedArgumentIsKeptForConversion (rows 19, 20 and 32), Stage3OwnedDeclarationsAreKeyedByOwner (row 25), Stage6MemberFragmentOutsideTheVocabulary (row 14), guarantees 1 to 6. Green. Not run in full: row 10 over all 21 packages. Reason: ten of the 21 are moved (the eleventh moved package, the TypeScript optional package, is new), and the eleven unmoved packages return text diagnostics on the mapping contract.

Line 3, serialization. Stage5LoweredGoCompilesAndRuns (a fresh process), Stage4VersionAxesAlone and Stage7StandinRefusesUnknownVersions (row 16), Stage3OneMeaningTwoRepresentationRevisions (row 17), PayloadBootstrap (row 38), Stage5ExactPayloadPassedValidation (row 38). Green. The bootstrap test was not shown red.

Line 4, a second reader in another process. Stage7StandinRoundTripsPageOfPerson, Stage7StandinRefusesUnknownVersions, Stage7ReadersNeedTheMatchingSnapshot (row 34), Stage6FourStatesOfAnOptionalField (rows 13 and 35), Stage7StandinDerivesMemoKeys (row 22). Green. Red under a protocol that reads a present undefined as absence. Not run: present with null where the type admits it (#211).

Line 5, executed Go lowering. Stage5LoweredGoCompilesAndRuns, Stage5MemberByReferenceNotByString (row 29). Green, and the lowering by string does not compile. Limit: F14.

Line 6, realization for TypeScript. Stage4RealizationManifest (rows 23 and 32): green, with a stand-in for the reader. The test of row 24 is NOT RUN. Reason: main has no reader on the formatter's result (#273 is open), F11.

Line 7, the shared packages fixed. Stage9ReplacedContributionLeavesTheOthers (row 1), Stage9PageOfPersonRegeneratesAndConverts (row 40), Stage9OutOfTreeModule (row 37), Stage9ExtensionDomainRefusals with the fixtures of earlier stages (row 33). Green in the small form. One file outside the contribution changed (F24). Not shown: the commit pair of the migration on main.

The two checks that stay green: shown for Go and TypeScript by Stage1NativeDeclarationsDoNotChange. Rust and C++ are not moved.
The condition beside the lines: row 11 is not decided here. Section 8 is the evidence.

5. THE GUARANTEE TABLE

1. Required facets belong to the selected contribution. Stage5FacetsBelongAndMatch: without one facet a request for bindings is "missing facet" and a request for types succeeds. No mutant.
2. Kind, target and compatibility match. Stage5FacetsBelongAndMatch. Red under facet-ignores-issuer, facet-ignores-target, facet-ignores-version.
3. A reference has the right plan, category and scope. Stage5ReferencesAreChecked. Red under reference-ignores-plan, reference-ignores-category, reference-ignores-scope.
4. Evidence belongs to its dependency. Stage2AnswersAreBoundToTheirDependency. Red under evidence-ignores-dependency, evidence-ignores-environment, evidence-ignores-representation-revision, evidence-ignores-meaning-revision.
5. The payload snapshot is immutable. Stage1FrozenValueSurvivesMutation. Red under projector-gets-returned-value and freeze-at-serialization.
6. The exact payload passed validation. Stage5ExactPayloadPassedValidation: a value that the Go type admits and the schema (one of two exact atoms) does not is refused before the projector runs. Red under codec-skips-validation. A receipt for another frozen value or another schema fails.

6. THE FINDINGS

Line numbers of accepted pages are those at 2a29674.

F1. The binder placeholder.
Page: "The protocol", row 3: "the planner resolves them, the rule decides with typed answers".
What happened: on main the declaration recipe runs before its body, so the body is resolved with the native meaning of each binder known. Under the protocol the declaration contribution decides after its body is planned and its projector states the native parameters. While the body of Page<T> is planned, the native meaning of T is not known (cplanner.declare, Stage1NativeDeclarationsDoNotChange).
Deviation: the planner represents a bound use by a placeholder (BinderExpression) and replaces it when the declaration's projector has returned its convention. No rule or projector sees it.
Proposal: say under "Generics are a division of responsibility" that the kernel carries a binder as a scoped reference until the declaration contribution has stated the native convention. Accepted page: "The kernel owns capture identity, substitution and scope." (docs/generation-planner.md:435).

F2. The home of a family's meaning.
Page: O4 and O19, the family's meaning is "the same in every target". O25: identity and payload types live "in a vocabulary package of that family for that target".
What happened: no page names a package that two targets' vocabularies share. The meaning identity is spelled in each target's vocabulary package and a test compares the spellings (Stage3OneMeaningTwoRepresentationRevisions).
Deviation: two spellings.
Proposal: name a vocabulary package of the family, without a target, that holds the meaning identity and its revision, and let each target's vocabulary import it. Accepted pages: "Identities live in vocabulary packages." (AGENTS.md:55-56), and of domain predicates "with a stable identity and version referenced by every target's binding" (docs/representation-families.md:154-155). Who owns that package where two publishers contribute stays open, as O19 says.

F3. A record's member names are in two places.
Page: row 5, a member name is a field of the payload. O3: a name stated in a request "is recorded once, there, and the decision does not restate it".
What happened: both hold only if the record's payload has the names of the fields it lays out itself and the name of a delegated field is in its request.
Deviation: none, built so. The record's projector zips names from the typed view with fragments from the dependency views by link order.
Proposal: say it under "How the name travels".

F4. The catalogue test of mapping packages.
What happened: TestBuiltinMappingsStateTheStructuralProfile on main names every mapping package of a pair by directory and requires each to export Library. The new TypeScript optional package failed it.
Deviation: one exclusion line in integration/go/horos_generation_profile_test.go.
Proposal: no change of the page. The migration replaces that test. Under D3 it changes again for Rust and C++.

F5. "a rule that reads one application to decide fails".
Page: "Evidence", rows 7 and 8, and the gate, line 2.
What happened: the kernel gives no rule an application to read. Inside Page<T> the answer about T is "deferred" whatever Page is applied to. A rule can still decide differently on "deferred" than on "proved", and the kernel cannot refuse that, because it does not know what a rule does with an answer.
Deviation: "fails" is built as "fails the coherence test of its family" (Stage2GenericAndWrittenOutAgree).
Proposal: word the gate's test so: "a rule that decides differently inside a generic declaration than for its written-out application fails the family's coherence test". The accepted page already says whose obligation it is: "This is a trusted-plugin obligation" (docs/type-generation-plugins.md:487), and of such tests "Deliberately broken rules show that the tests detect a failure. They do not prove the laws" (docs/representation-families.md, quoted on the page).

F6. What makes a cycle an evidence cycle.
Page: "A rule that demands evidence from a dependency whose own decision demands evidence from that rule cannot be planned."
What happened: the declaration contribution itself demands evidence about its body (row 3). Read as "a cycle with an evidence demand on it", every recursive declaration is an evidence cycle and the two refusals of row 9 cannot be told apart.
Deviation: a cycle is an evidence cycle when every open request on it demands evidence. With one reference on it, it is a recursive reference.
Proposal: state that criterion. Accepted page: "Recursive adapters are built through explicit plan references (placeholders)", so a cycle with a reference on it is the case placeholders admit.

F7. The horos source cannot state the recursion cases.
What happened: the horos analysis establishes no sort for a recursive declaration, so the horos source refuses it as "required semantic evidence is not established" before the kernel meets it (stage2-run1.txt). The test of row 9 runs on the independent fixture source.
Proposal: none for the page. Seen and not assessed: whether #164 needs a change of the horos analysis first.

F8. A request for conversion only.
Page: step 1 names two needs, "a reference or evidence". O5: a request for types "does not fail for want of something that only conversion needs".
What happened: for a declaration that the policy represents, the planner knows from the native convention which actuals are demanded. An exact contribution names its requests itself and does not know the products.
Deviation: a third need, generation.Conversion(child): resolved when bindings are asked for, left out for types alone, recorded on the link either way (Stage4OmittedArgumentIsKeptForConversion, Stage8Case3).
Proposal: add the need to step 1. Accepted page: "Only demanded native actuals need mappings." (docs/generation-spi.md:170).

F9. A specialization's domain.
Page: row 21, "The contract can name an application".
Deviation: ForSpecialization(declaration, actuals...), each a source selector. An actual that is itself an application is matched by its callee alone. Enough for the refusal, no design of #166.
Proposal: none.

F10. The shared helper request is stated and not built.
Page: row 25, and "The key under which a contribution asks for a shared helper" is not specified.
Deviation: the request kind exists and is refused as unsupported. Under the delegated form no moved contribution needs it: Go's record names no helper, and the optional contribution owns its helper. The consumer-created helper key of main is gone.
Proposal: defer the shared helper request with a trigger, or say that it is in the contract and refused until a case needs it.

F11. Gate line 6 cannot be shown in full.
Page: "Line 6 needs the reader of #273", and the test of row 24.
What happened: #273 is open. The TypeScript formatter's map from a generated site to its placed site is private (format/files/typescript/provider/go/relocate.go).
Deviation: the test composition writes the manifest from the placement the consumer asked for and the generation output, runs the real Files formatter for that placement and checks its text. That is no reader. The test of row 24 is not run.
Proposal: none for the design. For the order of work: #273 has to land before the gate can be met on main.

F12. Tampering cannot exercise an axis alone.
What happened: a node's identity covers its schema reference and its contribution reference. Changing either in the bytes breaks the node's identity first.
Deviation: the test of row 16 plans each case, and edits the bytes only for the envelope and the source semantics.
Proposal: none. A note for the migration's test.

F13. The fourth axis has no place the kernel reads.
Page: row 16 puts it in "The facet that names a factory, among its requirements".
What happened: the kernel carries a facet unread, so only the reader of that facet can refuse an unknown value. Built so in the stand-in.
Proposal: say that the test of row 16 for this axis is a reader's test.

F14. The Go value vocabulary of the prototype is Go text with typed holes.
Page: row 29, the instrument "must compile and run generated construction, member access, an owned helper, and a generic adapter argument".
What happened: it does. Every type, member, owned helper and adapter argument enters the lowered text through a hole that holds a reference of the plan, which the backend checks and resolves. The text between the holes is the facet's own Go.
Deviation: no typed expression vocabulary for Go was built.
Proposal: none for the contract. The report must be read so: line 5 shows the reference categories and their checks, and executed code. It does not meet "source-code strings never stand in for typed target expressions" (docs/representation-families.md:68-69).

F15. A production mapping cannot name its facet kinds.
Page: the contribution's facet list is "What a contribution states about its facets" (O16). O18: "every binding facet kind and the Go value vocabulary live in the test tree". O8: "A production mapping imports nothing of them".
What happened: a mapping package of bittype can state its facet kinds only by spelling their identities a second time, which AGENTS.md rules out.
Deviation: the contribution states nothing. Its facet list is what the consumer's composition registered for it, checked against what each requested product requires.
Proposal: define the facet list as the facets registered for a contribution by its publisher. Accepted page: composition must "reject duplicates, incompatible profiles, missing facets and unsupported applications" (docs/binding-generation.md:113), and "A native facet and the facets a binding product needs are selected together." (:103). Both are met by the registered set.

F16. The snapshot of ground descriptors does not serve an open declaration.
Page: row 34, "a snapshot of ground descriptors and module pins".
What happened: the plan's nodes describe the open declaration Page<T>, which has no ground descriptor. Both readers need, for the body of Page<T>, the canonical name of each field.
Deviation: the snapshot carries for each declaration the source's description of its type as written, with parameters and references, and each child with its position.
Proposal: word row 34 as "the descriptors of the declarations as the source states them, open where the declaration is generic". Accepted pages: "Open generics remain open; ground specialization is not a prerequisite." (docs/generation.md, quoted on the page). The list of what a backend emits has "ground-descriptor constants" (docs/binding-generation.md:81); those are for ground uses.

F17. A link needs the source's label of its child.
What happened: a reader has to pair a link of a record node with a field of the source's description. The plan carried no pairing.
Deviation: a link carries position, the provider's own locator of the operand, passed through unread, and the snapshot gives each child the same position.
Proposal: add it to what a link carries. Accepted page: the kernel carries "provider-expression operands, with locators passed through unread" (docs/generation-spi.md:200).

F18. Guarantee 6 rests on the codec.
What happened: the kernel finds the codec by the schema reference the contribution states and trusts it. A codec that does not validate is not detected by the kernel. The receipt adds that a frozen value cannot be exchanged after validation.
Proposal: say under guarantee 6 that the codec is trusted as a target is, and that the kernel's part is the association of value, schema and receipt.

F19. The optional flag of the Go field form.
What happened: the Go model's field form carries an "optional" flag as data. To keep the native text unchanged the Go optional contribution writes it into its fragment. Nothing reads it in the Go target.
Proposal: none. Seen, not assessed.

F20. Finding numbers F20 to F22 of the log are folded into F4, F13 and F16 here.

F23. Evidence does not pass through a named declaration by itself.
Page: row 7, "Only facts of the selected representation are evidence".
What happened: what is selected for a use of Person is the declaration contribution's application, not the record in Person's body. In case 5 the fixture policy had to ask about its body and state the answer again. It has to name each kind, and kinds belong to vocabularies it may not know.
Deviation: built in the fixture policy only. bittype's own declaration contribution states nothing, so a rule that asks about an operand of a named record type is answered "unknown".
Proposal: decide whether a declaration is transparent for evidence. The accepted pages give a line: "If the source language makes an alias interchangeable with its target, choosing incompatible native nominal types may violate the selected representation profile." (docs/type-generation-plugins.md:208-210) and "Administrative aliases and renamed binders are permitted" (:485-486). From them: the kernel forwards the facts of an alias's body, and a nominal declaration answers for itself.

F24. A reader has its own composition.
What happened: the lowering process registers the facet and the codec of every selected contribution. For the contribution of line 7 the file integration/go/prototype_stage5_test.go changed (commit cb9c851): a table lets a contribution's file add both.
Proposal: say under line 7 that "the composition" is one for each process that plans or reads.

F25. The authoring route cannot make its own bootstrap artifact with public packages.
Page: row 37 and row 38.
What happened: the module outside the repository registers with public packages only. It cannot make the canonical module of its schema, the pin or the generated Go type: the test writes the first two with internal/horosmodel and the type by hand.
Proposal: the route needs a public way from a schema source to the pinned artifact. Not examined: #172.

F26. The shared part of the codec role was touched for one family.
What happened: the fixture family of guarantee 6 has a field that is one of two atoms. The shared package exports the value forms a codec builds, and a codec package may import neither horos nor ontos, so Atom and ReadAtom were added to transform/syntax/shared/payload/go in stage 5.
Proposal: let a codec package import the ontos value packages. Then the shared part validates and encodes, and a family needs no shared change for a new value form. The reason against it today: "Only the Horos frontend ... consume external Horos/Ontos model dependencies" (AGENTS.md:273). That sentence needs an entry for the codec role either way, because the shared part imports both.

F27. The two readers have a case for each source form.
What happened: to pair a link with the source's description of that child, the Go backend and the stand-in switch over the snapshot's descriptor kinds record and list. A family with another source form would add a case to both.
Proposal: none for the page. A shortcut of the prototype's snapshot. In the migration the snapshot gives the children of any descriptor by position, so that a reader has no case. Accepted page: the backend "never switches over families" (docs/representation-families.md:68).

7. WHAT THE PAGE LEFT TO THE PROTOTYPE

Go types and signatures. A contribution states a reference, target, profile, domain, schema reference, rule, one projector, declared facts and owned declarations. A rule is two functions without a target type: Name(step) returns requests, Decide(step) returns a Decision{Payload, Facts}. A projector receives a Projection with no source type parameter.

The form of an evidence kind. EvidenceKind{Identity, Version}, a comparable value owned by a vocabulary of the target or product. A fact is Fact{Kind, Outcome, Content}: the outcome is proved, refuted or exhausted, and the content is the owner's string, carried unread. Page: say that a kind is a qualified identity with a version and that a fact carries one of three stated outcomes.

How a rule could present an answer that was not issued to it. An answer is a value with unexported fields that a rule reads only through the step it was issued to: DecideStep.Present(dependency, answer). It is bound by value to the asking contribution, the subject's place, the dependency, the dependency's whole reference and the generic environment. A rejected answer fails the planning of the subject even where the rule drops the error. Page: name this as the entry point.

A binder's requirement. Declared by the selected representation contribution in step 1, as part of its demand (Require). Stated in the frozen record of that rule as the answer "deferred, required", and collected on the node of the open declaration as (binder, kind). Found by the planner on the open declaration, which is planned before any application. Discharged by the planner at each application from the facts of what was selected for the actual: proved is recorded with the actual's reference, refuted is "requirement unsatisfied", unknown or exhausted is "requirement unresolved". An actual that is itself a binder moves the requirement to that binder. This is within the bounds the page quotes.

The class of the member-fragment refusal. Unsupported, with the owner "member vocabulary of <target>". Reason: a later slice of that vocabulary can admit the same fragment, which is the page's criterion for unsupported. It is not "refused by design", because the contribution's domain covers the position.

The form of a refusal, the two refusals of row 9, an unexpected error. One type with a class and one value for each class. Recursive reference and evidence cycle are both unsupported, with the kernel as owner and two reasons. An unexpected error of a rule is RuleFailure, which matches no refusal value. A rule may use four classes: unsupported, refused by design, requirement unsatisfied, requirement unresolved. The set of eight values can close as the page lists it.

Where the packages live and what they import. Section 10.

What the shared part of the codec role does and what a family's codec does. The shared part loads a pinned schema module, resolves the ground descriptor of the payload type, validates by horos conformance, encodes canonical ontos bytes and decodes them. A family's codec is two hand-written functions between its generated Go type and an ontos value. Page: "encodes, validates and freezes" is done by the family's codec through the shared mechanism. See F26 for the import rule.

How a codec is registered and how the kernel knows it is the one of the schema. The consumer's composition registers it. A codec states its schema reference, a registry holds one codec for each reference, a contribution states the schema reference of its payload, and the planner takes the codec registered under exactly that reference and refuses to plan without one. The receipt is a digest over the schema reference and the frozen bytes.

Which payloads are empty. Empty: text, boolean, list and optional, in both targets. Not empty: record (members: list<text>) and declarations (name: text, nominal: bool, parameters: list<text>).

Whether Go's optional package returns a decision before presence is delegated. On main it is selected for nothing: a helper recipe is found by a key and no rule runs, so there is no selection that could have a record. In the delegated form it is selected for a field position and returns an empty decision. Page: say "No" for the state before step 5.

The form of a contribution's facet list. F15.

Further points of "Not specified here". A request names a product by its identity. The reference categories issued are type, binder, adapter, declaration, owned helper and member; external symbol and requirement are not issued. A member reference is issued by the plan from the node of the native type; the projector states the member names beside the type. A fragment may hold more than one member in the Go and TypeScript models of the prototype. The realization manifest is written by the test composition. The law of a field-state protocol is not fixed.

8. ROW 11: THE FIELD-STATE PROTOCOL

Go. The fragment is one member of its own type, BittypeOptional[T]. The protocol has two operations: Construct(present, abstract value) returns the member, and State(member) returns (present, abstract value). The record's factory had to know which abstract field belongs to the position (from the snapshot) and the member (a member reference). It hands over a member value and takes one back. The optional contribution had to know the adapter selected for its operand and the declaration it owns. It did not need the record, the member's name or the canonical name.

The TypeScript stand-in. The fragment is one optional property. The protocol has two operations: write(container, member key, present, canonical value), and read(container, member key) returns {present, value}. The record's factory hands over the container object and the member's key, because an absent property is a fact about the object and has no value. The optional contribution had to know the adapter selected for its operand.

Neither side needed the other contribution's payload, a helper key from the consumer, or a layout across fields.

The four states. Absent and present with a value convert there and back as themselves on both sides. Present with null is refused by the operand's adapter on both sides. Present with undefined is refused by the protocol in the stand-in. It does not exist in Go, where a member is a value of its type.

Delegation of presence was workable in both. The one asymmetry is the shape of the slot: value in and value out is enough for Go, and TypeScript needs the container and the key. A slot kind for every target has the TypeScript shape. One seat wrote both sides, and the stand-in is not the service library's adapter interface, so this is evidence about the plan and the protocol and not about that interface.

9. ROW 26A: A FACET CROSSED THE PROCESS BOUNDARY

Yes: the facet that names a run-time factory. The stand-in finds a factory by what the plan entry names, so the factory reference has to be in the serialized plan. The transport descriptor designed here: a facet kind may state a media type and an encoder of its owner. When the planner writes a selected contribution into the plan, it calls the encoder for each facet that a requested product requires and stores the kind's identity and contract version, the media type and the bytes with the contribution's entry. It reads none of it. A reader finds the facet by the kind's identity, checks the version and the media type and decodes the value. A kind without a transport, such as the Go lowering facets, never leaves the process. The value that crossed is {module, export, runtime}.

What the reader needed that the plan did not carry: the pairing of a link with a source child (now carried as position, F17), and the derivation of the receipt, which the stand-in has written into it and which belongs to what the envelope version means.

10. ROLES AND EDGES ADDED TO THE ARCHITECTURE TEST, AS AN EDIT OF LAYOUT.md

New rows of the role table:

Vocabulary of a contribution. Path transform/syntax/<source>-to-<native>/vocabulary/<name>/go. Imports the generation contract and the publisher name. Never a mapping, a codec, a facet package or another vocabulary. Reason: O25.

Contribution: payload codec. Path transform/syntax/<source>-to-<native>/codec/<name>/go. Imports the generation contract, the shared part of the codec role and its own vocabulary. Never a mapping, a facet package or another codec. Reason: row 28, O18.

Contribution: facet package. Path completeness/<source>-to-<native>/facets/<name>/go. Imports the generation contract, the value vocabulary of its target's product and its own contribution's vocabulary. Never a mapping, a codec or another facet package. Reason: O8, O16, O25.

Vocabulary of a completeness product. Path completeness/<native>/values/go. Imports the generation contract and the syntax of its target. Reason: row 29, the facet kinds belong to the product.

Contribution: completeness backend. Path completeness/<native>/backend/go. Imports the generation contract and its product's vocabulary. Reason: row 31.

New entries of the foundation and kernel table: transform/syntax/shared/payload/go (imports the generation contract; external: horos core, horos declarations, ontos core, data and codec) and transform/syntax/shared/publisher/go (imports nothing).

Changed role: a mapping may import its own vocabulary, and is forbidden every codec, facet package and other vocabulary.

The test tree is the top-level directory completeness/. No other role permits an import of it, so only a composition and its own packages import it. The path is the prototype's choice. The page's nearest role, the verification law, was not used, because it is bound to one pair by name.

Not classified by the test: the JavaScript files of the stand-in.

11. SHARED PACKAGES AND FILES THE PROTOTYPE CHANGED

Belongs to the contract (part of the migration):
contracts/generation/go/contribution.go, contribution_plan.go, contribution_record.go, facet.go: the contract, its planner, the plan and the facets. New.
contracts/generation/go/expression.go: BinderExpression (F1).
contracts/generation/go/source.go: Source.Roles and Source.Semantics (rows 12 and 16).
contracts/generation/go/output.go, planning.go, output_test.go: the planner's identity on the inventory record (row 23).
transform/syntax/shared/generation/go/contributions.go and one field in registry.go: registration and selection.
transform/syntax/shared/payload/go: the shared part of the codec role. Atom and ReadAtom in it are F26.
transform/syntax/shared/publisher/go: bittype's publisher name.
transform/syntax/horos/source/go: the role issuer, Field.Role, the role key and the semantics version (row 12).
transform/syntax/horos/mapping/declarations/go/contribution.go: the declaration policy in two steps.
transform/syntax/golang/model/go and typescript/model/go: the evidence kind of the record form and the member vocabulary. The same for the two targets in target/go.
integration/go/architecture_roles_test.go and architecture_test.go: the roles above.
integration/go/testdata/generationfixture: the fixture reader issues roles, and the prototype's program.

Shortcuts of the prototype (not for the migration):
The mapping contract left in place beside the new one, and Library kept in each moved package.
generation.PrototypeMutants and every reading of it, in the contract, the registry's callers and the shared codec part.
integration/go/horos_generation_profile_test.go: the exclusion of F4.
The JSON encoding of the plan, and node identity as a digest of that JSON.
The snapshot written by a test with internal/horosmodel, and its cases in both readers (F27).
All of completeness/ in the form it has: the vocabulary of text with holes, one JavaScript module for all stand-in factories.
The lenient target in one test, and the switch in the shared codec part.

12. WHAT THE MIGRATION WOULD HAVE TO DO PER PACKAGE

As far as the eleven packages of the prototype show it (ten moved and one new), each package gets four things.

1. A vocabulary package: the reference (publisher, local key, meaning, representation revision), the schema reference, the schema source, its canonical module and pin, and the generated payload type. For seven of eleven the payload is empty.
2. A codec package: two hand-written functions and a constructor. About ten lines where the payload is empty.
3. The rule in two steps and the projector, in the mapping package. For text, boolean and list this is mechanical: a fixed output becomes a projector, a capture becomes one Reference request. For record it is a real rewrite: read the fields in Name, state one request for each field, state member names in Decide, and build the native record in the projector from the typed view and the dependency views. For declarations the shared horos recipe splits into requests, a decision and a description; the test "is the body a record" becomes a demand for evidence, and the target's model gains the evidence kind.
4. Its refusals as values with a class: decide for each text diagnostic of today whether it is unsupported or refused by design.

Beyond the packages: each target's model states its member vocabulary and the record-form fact; each record contribution declares that fact; each composition registers codecs beside contributions; the tests that list mapping packages change. Under D3 every target gains an optional package, and Go's record and optional packages lose the consumer-created helper key. Rust and C++ decide their optional representation inside the record today by an external type; under D3 that moves into a new optional package for each. If row 11 is decided as D1, the record's Name step keeps Reference requests for optional fields and the payload gains the presence decision; the field request, the role and the member vocabulary are not needed.

Not shown by the prototype: Rust and C++ at all, the TypeScript boolean package, and a facet package in production. Every facet of the prototype is test code.

13. SEEN AND NOT ASSESSED

- The origins inside the native AST now come from the kernel's reading of the graph where the old recipes took them from the source reader. The printed text is equal. The origins were not compared.
- Limits. The frozen payload bytes and the lowered expressions are charged to the output budget. Whether the budgets of main are the right ones for the new contract was not examined.
- Concurrency. PrototypeMutants is a package variable. The prototype's tests do not run in parallel.
- The registry's ReplaceContribution is new and has one test. The Replace of the mapping contract is untouched.
- A helper that a lowering introduces is keyed by contribution reference and name in the Go backend, and its spelling is issued by the backend. Whether that key is the "reference of the contract" of row 29 was not examined against the owned key of row 25.
- The Go backend lowers a record only as the body of a nominal declaration and with one member for each link.
- The stand-in validates no payload against its schema. It reads the canonical bytes as data and recomputes the receipt.
- Composition across several profiles, and how the authoring route meets a command line (#172), were not touched.
- Whether the horos analysis has to change before #164 (F7).
- Disk: C: stayed between 35800M and 37100M free. No go command ran in the background.
- The shared temporary directory. Several of my small patch scripts were written to /tmp, which on this workstation is the shared directory C:/Users/julia/AppData/Local/Temp. Their names were p.py, fix.py and patch5.py. If another lane kept a file of one of those names there, I overwrote it. Once, near the end, a command of mine ran a file fix2.py in that directory that was not mine: my own write of that name had been blocked together with a cleanup command. That script stopped at its second line, on a directory that does not exist, and changed nothing. Later scripts are in the evidence directory.
- The directory gotmp/ in the evidence directory still holds what the go commands left. A cleanup of it was blocked by a safety check and was not done.
```
