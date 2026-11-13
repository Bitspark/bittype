# Prototype of the contribution contract (bittype#196, research 0017 row 31)

**A throwaway branch. Nothing here merges.** Its products are evidence and
findings about the design page `docs/contribution-contract.md` (merged as
`2a29674`). The brief is
`C:/Development/bitspark/scratch-bittype-check/owner-1112/delivery-0017/prototype/BRIEF.md`.
The evidence directory is the directory of the brief. Whoever resumes this
work starts here and at the pushed branch `prototype/contribution-contract-196`.

How to run anything here (the disk rule of the brief):

```sh
E=C:/Development/bitspark/scratch-bittype-check/owner-1112/delivery-0017/prototype
export GOTMPDIR=$E/gotmp GOCACHE=Y:/scratch/bittype-196p/gocache
export BITTYPE_TYPESCRIPT_NODE=$(command -v node)
export BITTYPE_TYPESCRIPT_MODULE=C:/Development/bitspark/scratch-bittype-check/owner-1112/ts/node_modules/typescript/lib/typescript.js
df -BM /c    # no go command below 10000M free; one go command at a time
go test ./integration/go/ -run 'TestPrototype' -count=1 -v
```

Every test of the prototype is named `TestPrototype...` and lies in
`integration/go/prototype_*_test.go`. A test logs `GREEN <property>` for the
kernel as built and `RED <property> under mutant "<name>"` for a kernel with
one property switched off, so one run shows both. The mutants are the switch
`generation.PrototypeMutants` in the contract package. It is an instrument of
this branch and no part of the contract.

## What stands where

| What | Where |
| --- | --- |
| The contract: identities, domains, the rule in two steps, projector, codec interface, evidence, refusals | `contracts/generation/go/contribution.go` |
| The planner of the contract | `contracts/generation/go/contribution_plan.go` |
| The plan, its nodes, use entries, symbols, serialization, receipts | `contracts/generation/go/contribution_record.go` |
| Facets, products, reference categories | `contracts/generation/go/facet.go` |
| Registration and selection | `transform/syntax/shared/generation/go/contributions.go` |
| The shared part of the payload codec role | `transform/syntax/shared/payload/go` |
| The publisher name of bittype's own contributions | `transform/syntax/shared/publisher/go` |
| The vocabulary of a contribution for a target (identity, schema, generated payload type) | `transform/syntax/horos-to-<native>/vocabulary/<name>/go` |
| The payload codec of a contribution's schema | `transform/syntax/horos-to-<native>/codec/<name>/go` |
| The moved rule and projector of a mapping | `transform/syntax/horos-to-<native>/mapping/<name>/go/contribution.go` |
| The horos declaration policy, shared by the targets | `transform/syntax/horos/mapping/declarations/go/contribution.go` |
| The test composition | `integration/go/prototype_composition_test.go` |

## Two decisions about the branch itself

**The mapping contract stays beside the contribution contract.** The brief
allows the unmoved mapping packages "behind an adapter, or out of the build,
whichever is cheapest". The cheapest is neither: `generation.Plan` and
`MappingLibrary` stay as they are, and `generation.PlanContributions` stands
beside them. The two share the walk, admission, the native program, the
target and the sealed output. They do not meet in one plan. Rust, C++ and the
TypeScript boolean package stay on the mapping contract and are untouched.
The migration deletes the mapping contract.

**A moved mapping package keeps its `Library`.** Each moved package has a new
file `contribution.go` that exports `Contribution`, and keeps its old file
with `Library`. So the tests of today run unchanged, and
`TestPrototypeStage1NativeDeclarationsDoNotChange` compares the two contracts
on one module. The migration deletes `Library`.

Moved: for Go `text`, `boolean`, `list`, `record`, `declarations`, `optional`;
for TypeScript `text`, `list`, `record`, `declarations`, and a new package
`optional`. Go's `boolean` is moved beyond what the brief lists, because the
payload schema of the declaration contribution has a boolean field and the
payload types are generated with the Go contributions.

## The form of the contract, in short

Signatures are the prototype's (the page decides none).

- **A contribution** (`Contribution[F, T, R]`) states a reference, a target,
  a profile, a domain, a payload schema, a rule, one projector, the facts its
  revision fixes and the declarations it owns.
- **A rule** (`Rule[F, R]`) has two functions and no target type parameter:
  `Name(NameStep) -> []Request` and `Decide(DecideStep) -> Decision`. It never
  gets the planner. A request is `Reference(child)`, `Evidence(child,
  demands...)`, `Field(operand, role, stated)` or `SharedHelper(key)`.
- **A decision** is `Decision{Payload any, Facts []Fact}`: a value of the
  generated plain Go type, and beside it the facts others may ask of it.
- **The codec** (`PayloadCodec`) is found by the schema reference that the
  contribution states. `Freeze(value) -> bytes` encodes and validates.
  `View(bytes) -> value` decodes a fresh typed value. The kernel keeps the
  bytes and issues a receipt over schema reference and bytes.
- **A projector** (`Projector[T]`, or `DeclarationProjector[T]` for the
  declaration policy) gets a `Projection[T]`: the typed view, the dependency
  views, the stated name of a field position, owned declarations, diagnostic
  sites. It has no type parameter of the source and no reader.
- **The plan** (`RepresentationPlan`) holds nodes, use entries, symbols and
  the selected contributions as values, and serializes to JSON.

## Stage 1: the staged rule and the frozen record (gate line 1)

**Built.** The five steps of the page's protocol, in `cplanner.run` (steps 1
to 4) and `cplanner.represent` / `cplanner.declare` (step 5). The codec role
with its shared part and eleven codec packages. Eleven vocabulary packages
with the pinned bootstrap artifact. The moved rules and projectors.

**What the page said, and what happened.**

- *"the planner resolves them, the rule decides with typed answers"* (row 3).
  Built as written. The planner plans each requested child completely, with
  its own projector, before the asking rule decides.
- *A consequence the page does not name: the binder placeholder.* On `main`
  the declaration recipe runs before its body, so the body is resolved with
  the native meaning of each binder already known. Under the protocol the
  declaration contribution decides after its body is planned, and its
  projector states the native parameters. So while the body of `Page<T>` is
  planned, the native meaning of `T` is not known yet. The planner represents
  a bound use by a placeholder (`BinderExpression`) and replaces it when the
  declaration's projector has returned the convention. No rule or projector
  sees the difference. This is finding F1 (a consequence, no contradiction).
- *"The codec encodes, validates and freezes the value at the moment the rule
  returns"* (O9). Built as written. `Freeze` runs directly after `Decide`
  returns, for every request.
- *"The projector receives the typed view decoded from the frozen value"*
  (O9). Built as written: `View(frozen bytes)`.
- *"A projector that tries to read the source does not compile or fails the
  contract test"* (row 4). Both hold. `Projector[T]` has no type parameter
  for the source language or the reader, and `Projection[T]` has no method
  that mentions a use, a declaration, a binder, a selector or the reader.
- *Diagnostic sites* (row 4: "through the record or the kernel's views").
  Chosen: the kernel's views. A projection gives the site of its subject, of
  each dependency and of each binder. The sites are the kernel's own reading
  of the graph, where the old recipes took them from the source reader. The
  printed native text is unchanged; the origins inside the native AST were
  not compared.

**What the page left open and what was chosen.**

- *Any Go type or signature.* As above.
- *What the shared part of the codec role does and what a family's codec
  does.* The shared part (`transform/syntax/shared/payload/go`) loads a pinned
  schema module, resolves the ground descriptor of the payload type, validates
  by `horos.Conforms`, encodes canonical ontos bytes and decodes them. A
  family's codec is two hand-written functions between its generated Go type
  and an ontos value. So "encodes, validates and freezes" is done by the
  family's codec *through* the shared mechanism: the family maps, the shared
  part validates and encodes.
- *How a codec is registered and how the kernel knows that the codec that
  validated is the one of the payload's schema.* The consumer's composition
  registers a codec (`Registry.RegisterCodec`). A codec states its schema
  reference, and a registry holds one codec per reference. A contribution
  states the schema reference of its payload. The planner takes the codec
  registered under exactly that reference, and refuses to plan when none is
  registered (`ErrIncomplete`). The receipt is a digest over the schema
  reference and the frozen bytes, issued only where the planner freezes.
- *Which payloads are empty.* Empty: text, boolean, list and optional, for
  both targets. Not empty: record (`members: list<text>`, the names of the
  members the record lays out itself) and declarations (`name: text`,
  `nominal: bool`, `parameters: list<text>`).
- *Whether Go's optional package returns a decision before presence is
  delegated.* Not decided by the prototype: the prototype builds the
  delegated form only. In that form the package is selected for a field
  position and returns an empty decision, and the record exists. The helper
  `BittypeOptional` is a declaration the contribution owns. As it stands on
  `main` the package is selected for nothing: a helper recipe is found by a
  key and no rule runs, so there is no selection that could have a record.
- *The path of a codec package and of a vocabulary package.* See the table
  above. The roles are in `integration/go/architecture_roles_test.go`
  (`bundleRole`).

**Findings of this stage.**

- **F1.** The binder placeholder, above.
- **F2. The home of a family's meaning.** O4 and O19: the family's meaning is
  "the same in every target". O25 puts a contribution's identity into "a
  vocabulary package of that family for that target". No page names a package
  that the vocabularies of two targets share. The prototype spells the meaning
  identity in each target's vocabulary package (`bittype.structural/text` in
  two packages), and a test of stage 3 compares them. That is one identity
  spelled twice, which `AGENTS.md` rules out ("Identities live in vocabulary
  packages").
- **F3. The record's member names are recorded in two places.** Row 5 makes a
  member name a field of the payload. O3 says a name stated in a request "is
  recorded once, there, and the decision does not restate it". Both hold only
  if the record's payload has the names of the fields it lays out itself, and
  the name of a delegated field is in its request. Built so. The record's
  projector reads names from the typed view and fragments from the dependency
  views, and zips them by the order of the links.
- **F4. The catalogue test of mapping packages.** A test on `main`
  (`TestBuiltinMappingsStateTheStructuralProfile`) names every mapping package
  of a pair by directory and requires each to export `Library`. The new
  TypeScript `optional` package fails it. It is a shared file that a family
  touches. One line excludes the package on this branch.

**Shared files changed in this stage, and why.**

| File | Change | Contract or shortcut |
| --- | --- | --- |
| `contracts/generation/go/expression.go` | `BinderExpression` | Contract (F1) |
| `contracts/generation/go/source.go` | `Source.Roles`, `Source.Semantics` | Contract (rows 12, 16) |
| `contracts/generation/go/output.go`, `planning.go`, `output_test.go` | `OwnedDeclaration.Symbol`, the planner-stated field of the inventory record | Contract (row 23) |
| `transform/syntax/shared/generation/go/registry.go` | one field for the registration state | Contract |
| `transform/syntax/horos/source/go` | the role issuer, `Field.Role`, the two constants | Contract (row 12) |
| `transform/syntax/{golang,typescript}/model/go` | `EvidenceRecordForm`, `CheckMember` | Contract (rows 3, 13) |
| `transform/syntax/{golang,typescript}/target/go` | `Provider.CheckMember` | Contract (row 13) |
| `integration/go/architecture_roles_test.go`, `architecture_test.go` | the roles and entries of the new packages | Contract (the layout edit) |
| `integration/go/horos_generation_profile_test.go` | one exclusion | F4 |

**Tests** (`integration/go/prototype_stage1_test.go`,
`prototype_bootstrap_test.go`). Evidence: `stage1-run1.txt` (verbose, red and
green), `stage1-full.txt` (the contract, transform and integration packages;
the only failures are the Rust tests that need the rustc helper, and F4 before
its exclusion), `architecture-red-stage1.txt` (the architecture test before
the roles were added).

| Test | Shows | Does not show |
| --- | --- | --- |
| `TestPrototypeStage1NativeDeclarationsDoNotChange` | The moved contributions generate, through the five steps, the same Go text as the mapping contract for today's pagination module, and the same TypeScript text for that module without booleans | Rust, C++; TypeScript booleans; equality of AST origins |
| `TestPrototypeStage1ProjectorCannotReadTheSource` | The projector's capability mentions no source handle (reflection), three attempts to read the source do not type-check, and the same contract test is red on the callback of the mapping contract | That a contribution cannot hand source data to its projector through its own package state (the page's "A limit") |
| `TestPrototypeStage1FrozenValueSurvivesMutation` | A value mutated after the rule returned changes neither the native type nor the serialized plan. Red under `projector-gets-returned-value` and under `freeze-at-serialization` | Immutability against a codec that keeps the value |
| `TestPrototypeStage1DecisionTheSourceDoesNotDetermine` | With a fixture whose member names come from its configuration, the plan entry and the native type agree, and an entry rebuilt from the source is detected | That no second pass can be written. The mutant is a pass of the test, not of the kernel |
| `TestPrototypePayloadBootstrap` | For eleven vocabulary packages: the committed module and pin are what the schema source lowers to, the committed Go type is what bittype generates with the prototype's own Go contributions, every listed value round-trips | TypeScript types of a payload (none are generated); a payload field the existing families cannot type |
| `TestProductionArchitecture` | Green with the new roles | That the roles are right |

## Stage 2: evidence (gate line 2, rows 3, 5, 7, 8, 9; guarantee 4)

**Built.** Evidence kinds, facts, demands, answers with their five outcomes,
binder requirements with their discharge at an application, the two refusals
of row 9, and the binding of an answer. Two consumers of the evidence path:
the declaration contribution (its body's record form) and the fixture family
`evidenceList` (`integration/go/prototype_fixtures_test.go`). A second set of
fixture contributions for the independent fixture source
(`integration/go/prototype_independent_test.go`, with
`integration/go/testdata/generationfixture/contributions.go`).

**What the page left open and what was chosen.**

- *The form of an evidence kind.* `EvidenceKind{Identity, Version}`: a
  qualified identity and a version, a plain comparable value. Its owner is a
  vocabulary of the target or product (`model.EvidenceRecordForm` in the Go
  and the TypeScript model). A fact is `Fact{Kind, Outcome, Content}` with the
  outcome proved, refuted or exhausted. `Content` is a string that the kind's
  owner reads and the kernel carries unread. A rule states facts in
  `Decision.Facts`; a contribution declares them in `Contribution.Facts`.
- *How a rule could present an answer that was not issued to it.* An answer
  is a value (`Answer`) with unexported fields. A rule reads it only through
  the step it was issued to: `DecideStep.Present(dependency, answer)`.
  `DecideStep.Evidence(dependency, kind)` is the ordinary way and calls
  `Present` with the step's own answer. An answer is bound, by value, to the
  asking contribution's reference, the subject's place in its declaration,
  the index of the dependency, the dependency's whole contribution reference
  and the generic environment (the declaration being planned). A rejected
  answer fails the planning of the subject even if the rule drops the error.
  The binding is by value so that the test can keep everything equal but one
  part. An answer kept from an identical run would be read; it says the same.
- *A binder's requirement: in which record, in what form, how found.*
  - **Declared** by the selected representation contribution, in step 1, with
    the demand itself: `Demand{Kind, Require: true}`. Where the operand turns
    out to be bound to a binder, the planner records "this binder's actuals
    must prove this kind".
  - **Stated** in the frozen record of the rule that declared it, as the
    answer `deferred, required` on the link, and collected on the node of the
    open declaration that binds the binder (`PlanNode.Requirements`:
    binder index and kind).
  - **Found** by the planner on the open declaration, which it has planned
    before it plans any application of it.
  - **Discharged** by the planner at each application
    (`cplanner.application`): it reads the facts of what was selected for the
    actual. Proved: recorded in `PlanApplication.Discharges` with the actual's
    contribution reference. Refuted: "requirement unsatisfied". Unknown or
    exhausted: "requirement unresolved". An actual that is itself bound to a
    binder moves the requirement to that binder (outcome `deferred`).
  - Without `Require` the answer about a binder is `deferred` and establishes
    nothing.
- *The value of the two refusals of row 9.* Both are class "unsupported",
  owner `bittype/kernel`, with the reasons `recursive reference` and
  `evidence cycle`. Unsupported, because a later slice can admit both inputs:
  #164 for the first, "a declared evaluation discipline" for the second.

**Findings of this stage.**

- **F5. "a rule that reads one application to decide fails" (rows 7, 8).**
  The kernel gives no rule an application to read: no rule runs for an
  application, and inside `Page<T>` the answer about `T` is `deferred`
  whatever `Page` is applied to. What a rule can still do is decide
  differently on `deferred` than on `proved`. The kernel cannot refuse that:
  it does not know what a rule does with an answer. So "fails" can only mean
  "fails the coherence test of its family". Built so: the test compares
  `Page<text>` with the written-out record, is green for the rule that
  declares the requirement and red for the rule that makes a choice for want
  of proof inside the generic only.
- **F6. What makes a cycle an evidence cycle.** The page: "A rule that
  demands evidence from a dependency whose own decision demands evidence from
  that rule cannot be planned." The declaration contribution itself demands
  evidence about its body (row 3). Read as "a cycle with an evidence demand
  on it", every recursive declaration would be an evidence cycle, and the two
  refusals of row 9 could not be told apart. Built: a cycle is an evidence
  cycle when *every* open request on it demands evidence, because then no
  decision on it can finish. With one reference on it, it is a recursive
  reference, which placeholders would admit.
- **F7. The horos source cannot state the recursion cases.** The horos
  analysis publishes no sort for a recursive declaration, so the horos source
  refuses `type Node = { next?: Node; }` as "required semantic evidence is
  not established" before the kernel meets the recursion. The test of row 9
  runs on the independent fixture source. This is seen and not assessed
  beyond that.
- **F8. A request for conversion only.** O5 lets a request for types alone
  succeed without what only conversion needs, and row 20 keeps an omitted
  argument "for conversion". For a declaration that the policy represents,
  the planner knows from the native convention which actuals the native type
  demands. For an exact contribution the rule names its requests itself, and
  it does not know the products. So a request needs a third need beside
  "reference" and "evidence": "for conversion only". Built as
  `generation.Conversion(child)`: resolved when bindings are asked for, left
  out for types alone, and recorded on the link either way. The page's step 1
  names two needs.

**Tests** (`integration/go/prototype_stage2_test.go`; evidence
`stage2-run2.txt`, and `stage2-run1.txt` for the first run, in which the
recursion test failed for the reason of F7).

| Test | Shows | Does not show |
| --- | --- | --- |
| `...Stage2FixtureFamilyDecidesOnEvidence` | A rule gets a typed answer about what was selected for its operand before it decides, and its native type follows it: `[]string` where the fact is declared, `[][]string` where nothing is stated | That the fact is true (AC-05) |
| `...Stage2GenericAndWrittenOutAgree` | With a declared requirement, `Page<text>` and the written-out record have the same native type, the open declaration records the requirement and the application its discharge. Red for a rule that decides differently inside the generic | That the kernel refuses the incoherent rule (F5) |
| `...Stage2FiveOutcomes` | Each of the five outcomes reaches a rule as its own answer | A check that reaches a limit: "exhausted" is stated, nothing computes |
| `...Stage2DeclarationPolicyAsksAboutItsBody` | The declaration contribution decides nominal or alias from the answer about its body and reads no payload of the record | A record under another contribution |
| `...Stage2RevisionIsIdentityAndNoPayload` | Two representation revisions: equal payloads, different nodes, also for the node that depends on it | What a revision means |
| `...Stage2CycleAndRecursionAreToldApart` | `Rec = Page<Rec>` and `Loop = probe(Loop)` are refused as unsupported with two reasons, on the independent source | Anything of recursion beyond its refusal |
| `...Stage2AnswersAreBoundToTheirDependency` | An answer for another dependency, another environment, another representation revision or another meaning revision of the dependency is rejected. Each is red under a kernel without that one comparison | That a rule cannot act on what it remembers without presenting it |

## Stage 3: identities and refusals (rows 10, 15, 17, 21, 25)

**Built.** The contribution identity and reference as the page writes them,
checked at registration. The refusal as one type with a class
(`generation.Refusal`) and a value for each class. Selection by domain with
"missing" and "ambiguous" as two values. The specialization domain and its
refusal. Declarations that a contribution owns, keyed by reference and local
name. No kernel code was added in this stage beyond one line; the stage is
its tests.

**What the page left open and what was chosen.**

- *The form of a refusal in Go.* One type with a class, and one value per
  class that `errors.Is` matches. The existing values are kept:
  `ErrUnsupported`, `ErrAmbiguous`, `ErrLimit`. New values:
  `ErrRefusedByDesign`, `ErrRequirementUnsatisfied`,
  `ErrRequirementUnresolved`, `ErrMissingFacet`, `ErrMissingContribution`. A
  refusal carries its owner and reason, the subject (source identity and
  site), the selected contribution with target and profile, the path of
  dependencies, the requirement and, for a specialization, the application.
- *What an unexpected error of a rule is.* `generation.RuleFailure`: an
  implementation failure that names the contribution. It does not unwrap, so
  the error of a rule cannot add a class. A rule that refuses with a class of
  the planner (missing facet, missing or ambiguous contribution, limit) is an
  implementation failure too. A rule may use four classes: unsupported,
  refused by design, requirement unsatisfied, requirement unresolved.
- *The class of the member-fragment refusal.* Fixed in stage 6: unsupported.
- *The form of a publisher's name, of a local key, of a meaning identity and
  of a meaning revision.* Strings. The prototype uses the module path as
  bittype's publisher name, `horos-to-<native>/<name>` as the local key,
  `bittype.structural/<family>` as the meaning identity, `1` as its revision
  and `<native>/1` as the representation revision. Nothing depends on the
  forms.
- *The classes of the moved packages' refusals.* "nested records require a
  named declaration in this profile": refused by design (the profile's
  representation does not promise it, and another contribution may). A name
  that the naming rule does not support: unsupported (the naming design
  admits the same input later).

**Findings of this stage.**

- **F9. A specialization's domain needs the source identities of its
  arguments.** Row 21: "The contract can name an application". Built as
  `ForSpecialization(declaration, actuals...)`, each a selector of the
  source. The planner matches an application by the selector of its callee
  and of each actual. An actual that is itself an application is matched by
  its callee alone. That is enough for the refusal the row asks for, and it
  is not a design of #166.
- **F10. The shared helper request is stated and not built.** Row 25: "A
  contribution that needs whichever shared helper the composition selects
  makes a request". The page leaves open whose key the request uses. The
  request kind exists (`SharedHelper(key)`) and the planner refuses it as
  unsupported. Under the delegated form no moved contribution needs one: Go's
  record no longer names a helper, and the optional contribution owns its
  helper. What on `main` needs the consumer to create a key
  (`optional.Library(key)` and `record.Library(key)`) needs none here.

**Tests** (`integration/go/prototype_stage3_test.go`; evidence
`stage3-run2.txt`, and `stage3-run1.txt` for the first run with one wrong
assertion of the test). The tests of row 5 and of recursion are in stage 2.

| Test | Shows | Does not show |
| --- | --- | --- |
| `...Stage3PublisherAndLocalKey` | One key under two publishers registers; one publisher with one key twice fails, also for a second target; the stated target is checked | That a publisher's name is its owner's (row 15a) |
| `...Stage3OneMeaningTwoRepresentationRevisions` | For text, list, record, declarations and optional: Go and TypeScript carry one meaning and two representation revisions; a reference that lacks either part is refused | Who owns a meaning; that the meaning is one identity and not two equal spellings (F2) |
| `...Stage3RefusalValues` | Missing and ambiguous are two values; an unbuilt family is "missing contribution"; a specialization is refused as unsupported with the application named while another application plans; five refusals of moved rules match exactly one value each; an unexpected error matches none | The refusals of the ten unmoved packages. The gate's test of row 10 over all 21 is not run |
| `...Stage3OwnedDeclarationsAreKeyedByOwner` | Two contributions that spell a helper alike get two symbols. Red under a kernel that shares by spelling. A shared helper request is refused | What a Files stage does with two declarations of one name; a shared helper |

## Stage 4: the plan and its serialized form (rows 16, 18, 19, 20, 22, 23, 32)

**Built.** The plan as values: representation nodes identified by a digest
of their content, use entries, applications with two lists tied to binders,
logical symbols, the selected contributions, the products, the target with
its configuration and the source with its semantics version
(`contracts/generation/go/contribution_record.go`). `Serialize` and
`ParsePlan` (JSON; the encoding is provisional). `ReaderKnowledge` and
`CheckKnown` for the version axes. The realization manifest and
`LinkRealization`. The planner's identity as a field of the inventory record
(`OwnedDeclaration.Symbol`).

**What the page said, and what happened.**

- *"One for all the uses that selected the same representation"* (row 18).
  Built by content: a node's identity is a digest over its kind, source
  identity, contribution reference, schema reference, frozen payload, receipt,
  links with their answers, stated facts, symbol, binder, application,
  requirements and member names. It holds no source site. "How the planner
  finds that two uses share a node" is #195's; the digest is the prototype's.
- *Open generics are primary* (row 20). `Page<T>` is one declaration node
  with its binder count. A use of `T` inside it is a binder node with the
  declaration's logical symbol as its scope. `Page<Person>` is an application
  node that names the declaration node and has two lists, `native` and
  `adapters`, each of pairs (binder index, node). They are never matched by
  position.
- *"A generation request states whether it asks for types alone or for types
  and bindings"* (O5). A request names products by identity; `types` is
  always among them. A product is registered by the composition and states,
  per kind of subject (use, declaration, field position), the facet kinds it
  requires. With no further product the planner plans only the arguments that
  the native convention demands. With one it plans an adapter argument for
  every binder.
- *"The planner's declaration identity is the logical symbol identity"* (row
  23). Built as a string: `source:<source identity of the declaration>` for a
  source declaration and `owned:<contribution reference>/<local name>` for an
  owned one. It is in the plan, on the declaration's node and on the
  inventory record of the sealed output.

**What the page left open and what was chosen.**

- *How a generation request names one binding product among several.* By the
  product's identity, a string that the composition registered a product
  under (`Registry.RegisterProduct`, `ContributionRequest.Products`).
- *Who writes the realization manifest and in which form.* On this branch
  the test composition writes it (`realizationTS`), as the struct
  `RealizationManifest`: the plan's content digest, the run, the target with
  its configuration, the placement and one realized symbol (logical symbol,
  name, module path) per planned declaration. #195 decides the real writer.

**Findings of this stage.**

- **F11. Gate line 6 cannot be shown in full: main has no reader on the
  formatter's result.** The page: "Line 6 needs the reader of #273". #273 is
  open, and the TypeScript formatter's map from a generated site to its
  placed site is private (`format/files/typescript/provider/go/relocate.go`).
  The prototype's manifest takes the module path from the placement that the
  consumer asked for and the name from the generation output, runs the real
  Files formatter for that placement and checks that the text it wrote has
  each declaration in the file the manifest names. That is no reader of the
  realization. So the tests of rows 23 and 32 are run with a stand-in for the
  reader, and the test of row 24 ("a formatter mutant that realizes another
  spelling is followed") is not run.
- **F12. Tampering cannot exercise an axis alone.** A node's identity covers
  its schema reference and its contribution reference. Changing either in the
  serialized bytes breaks the node's identity before any axis is compared. So
  the test of row 16 plans each case: with a schema of another version, with
  another meaning revision, with another representation revision. Only the
  envelope protocol and the source semantics, which are no part of a node,
  are changed in the bytes. No change to the page; a note for the migration's
  test.
- **F13. The fourth axis has no place the kernel can read.** Row 16 puts
  generated-code and runtime compatibility "in the facet that names a
  factory, among its requirements". The kernel carries a facet unread, so
  only the reader of that facet can refuse an unknown value. Built so in
  stage 7; `CheckKnown` covers the other four.

**Tests** (`integration/go/prototype_stage4_test.go`; evidence
`stage4-run1.txt`).

| Test | Shows | Does not show |
| --- | --- | --- |
| `...Stage4TwoUsesOneNode` | Two uses of `Page<Person>` are two use entries with two sites on one node; the open declaration is one node; the application links its binder to the node of `Person` | #195's procedure for sharing |
| `...Stage4OmittedArgumentIsKeptForConversion` | Asked for bindings, `Ignore<Person>` has no native and one adapter argument, and two applications with equal native types are two nodes. Asked for types alone the omitted argument is not planned. Red under `application-drops-omitted-argument` and `node-identity-ignores-links`. The same for an exact contribution with a conversion-only request | An adapter |
| `...Stage4MemoKeyFromTheSerializedPlan` | A reader that only parsed the serialized plan derives equal keys for two uses with equal nodes and another key with another selected dependency | The memo and its token (bitspark-service#40) |
| `...Stage4VersionAxesAlone` | An unknown value of the envelope protocol, the payload schema, the meaning revision, the representation revision and the source semantics is each refused under its own name, and the baseline is read | The fourth axis (F13, stage 7) |
| `...Stage4RealizationManifest` | The inventory carries the planner's identity; one plan links with the manifests of two placements and resolves each symbol to another module path and the same name; the manifest of another output is a linking error. Red under `link-ignores-plan-digest`, which then reports an incomplete plan | A reader on the formatter's result (F11); row 24 |

## Stage 5: facets, Go handlers and executed Go lowering (guarantees 1, 2, 3, 6; gate line 5; gate line 3, first test)

**Built.** Facet kinds, facets and products in the contract
(`contracts/generation/go/facet.go`). The test tree of the completeness
backend, a new top-level directory `completeness/`:

| Package | What it is |
| --- | --- |
| `completeness/golang/values/go` | The Go value vocabulary, the two facet kinds of the Go conversion product and the product |
| `completeness/golang/backend/go` | The completeness test backend for Go: it lowers from the serialized plan, the source snapshot and the facets, and holds the runtime of the lowered program |
| `completeness/horos-to-golang/facets/<name>/go` | One facet package for each of text, boolean, list, record, declarations, optional |

The facets are Go handlers: a facet's value is a Go function that the backend
calls in process (`values.Lowering`). No second, interpreting handler that
converts values in process was built: the lowered code is what is executed.

The source snapshot of row 34 is exported by the test composition
(`integration/go/prototype_snapshot_test.go`).

**What the page left open and what was chosen.**

- *The form of a contribution's facet list.* **The list is what the
  consumer's composition registered for the contribution, and the
  contribution itself states nothing** (finding F15). `Registry.RegisterFacet`
  admits a facet that names a registered contribution of the same publisher;
  the planner checks, for each selected contribution and each product the
  request names, that a facet of each required kind is among them.
- *Where a facet package lives and what it imports.*
  `completeness/<source>-to-<native>/facets/<name>/go`. It imports the
  contract, the value vocabulary of its target's product and its own
  contribution's vocabulary. It imports no mapping, no codec and no other
  facet package. The roles are `completenessRole` in
  `integration/go/architecture_roles_test.go`.
- *The set of reference categories, and who issues a member reference.*
  Issued by the plan: type, binder, adapter, declaration, owned helper,
  member. Not issued: external symbol, and requirement with its evidence. A
  member reference is issued by the plan from the node of the native type
  that contains the member (`RepresentationPlan.MemberOf(node, ordinal)`).
  The member's spelled name is stated by the projector beside the native type
  (`Projected.Members`), and the planner carries it on the node.
- *How a reader gets a typed payload.* A facet package may not import the
  codec (O25). The backend's composition gives it the codecs, and the backend
  decodes a node's payload with the codec of the node's schema, after it
  checked the receipt, and hands the typed view to the facet.

**Findings of this stage.**

- **F14. The Go value vocabulary of the prototype is Go text with typed
  holes.** The page (row 29) asks of the prototype's vocabulary that it
  "compile and run generated construction, member access, an owned helper,
  and a generic adapter argument". It does. Every type, member, owned helper
  and adapter argument enters the lowered text through a hole that holds a
  reference of the plan, which the backend checks with the plan and resolves.
  The text between the holes is the facet's own Go. So the accepted sentence
  "source-code strings never stand in for typed target expressions"
  (`docs/representation-families.md:52-55`) is not met by this instrument. A
  typed expression vocabulary for Go was not built. What the prototype shows
  for row 29 is the reference categories and their checks, not a vocabulary.
- **F15. A production mapping cannot name its facet kinds.** O8/O16: the
  contribution's facet list is "What a contribution states about its
  facets". O18: "every binding facet kind and the Go value vocabulary live in
  the test tree". O8: "A production mapping imports nothing of them". And
  `AGENTS.md`: "Identities live in vocabulary packages." So a mapping package
  of bittype can state its facet kinds only by spelling their identities a
  second time. Deviation: the contribution states nothing, and its facet list
  is what the composition registered for it. Guarantee 1 then reads: "the
  facets a product requires are registered for the selected contribution, by
  its publisher".
- **F16. The snapshot of ground descriptors does not serve an open
  declaration.** Row 34: "a snapshot of ground descriptors and module pins".
  The plan's nodes describe the open declaration `Page<T>`, which has no
  ground descriptor. A reader needs, for the body of `Page<T>`, the canonical
  name of each field. The prototype's snapshot therefore carries, for each
  declaration, the source's description of its type as written (with
  parameters and references), and each child with its position.
- **F17. A link needs the source's label of its child.** A reader has to
  pair a link of a record node with a field of the source's description. The
  plan carried no such pairing. Built: a link carries `position`, the
  provider's own locator of the operand (for horos `[k+2, 3]` for field k),
  passed through unread by the kernel, and the snapshot gives each child the
  same position. Without it a reader would have to pair by order, which is a
  decision of the record contribution and not of the source.
- **F18. Guarantee 6 rests on the codec.** The kernel finds the codec by the
  schema reference that the contribution states and trusts it to validate. A
  codec that does not validate is not detected by the kernel (the mutant
  `codec-skips-validation` plans). What the kernel adds is the receipt: a
  frozen value cannot be exchanged after validation without the receipt
  failing.

**Tests** (`integration/go/prototype_stage5_test.go`; evidence
`stage5-run2.txt`, `stage5-lowered.go.txt`, `stage5-plan.json`,
`stage5-snapshot.json`; `stage5-run1.txt` is the first run, in which the
mutant of guarantee 6 was not detected because decoding validated a second
time).

| Test | Shows | Does not show |
| --- | --- | --- |
| `...Stage5LoweredGoCompilesAndRuns` | A fresh process lowers Go conversion for `Page[Person]` from the serialized plan, the snapshot and the facets; the code compiles, converts five values there and back and refuses three that are outside the type; it has a typed construction, member access, the owned helper and a generic adapter argument | A typed expression vocabulary (F14); the conversion laws |
| `...Stage5MemberByReferenceNotByString` | With member names from a fixture's configuration, the lowering by member reference compiles and runs, and the lowering that spells members from the source does not compile | A naming policy of the naming design |
| `...Stage5FacetsBelongAndMatch` | Guarantee 1: without one facet a request for bindings is "missing facet" and a request for types succeeds. Guarantee 2: another issuer under the same name, another target and another contract version are rejected, each red under its mutant. A facet of another publisher is not admitted | What a facet does |
| `...Stage5ReferencesAreChecked` | Guarantee 3: a reference of another plan, a member where a type is expected, a binder outside its scope are rejected, each red under its mutant | External symbols; requirement references |
| `...Stage5ExactPayloadPassedValidation` | Guarantee 6: a value that the Go type admits and the schema (`atom-exact("direct") \| atom-exact("boxed")`) does not is refused before the projector runs; red under a codec that does not validate. A receipt for another frozen value or another schema fails | That the kernel detects a codec that does not validate (F18) |

## Stage 6: optional fields by delegation (row 11 as D3; rows 12, 13, 14, 35)

**Built.** The delegated form, in both targets. The record contribution asks
for a field position; the optional contribution answers with a member
fragment and, for conversion, with a field-state protocol; the record's
factory invokes the protocol and reads no payload of the optional
contribution. A TypeScript `optional` package (mapping, vocabulary, codec,
facet package). The TypeScript stand-in, which stage 7 uses as a reader:
`completeness/typescript/standin/{reader,structural,ontos}.mjs`, with its
vocabulary `completeness/typescript/values/go` and the facet packages
`completeness/horos-to-typescript/facets/<name>/go`.

**What the page said, and what happened.**

- *The subject of the request* (row 12). Built as written. The source reader
  issues the role selector (`source.Field.Role`, key
  `horos/optional-record-field`) for the record's use and the field's operand
  use. The record rule requests `Field(operand, role, memberName)`. The kernel
  checks that the operand is a child of the asking use and that the role was
  issued by the source's issuer for exactly that pair. It reads the role's
  key to select and nothing else: not the flag, not the field's name.
- *The native half* (row 13). Built as written: the fragment is in the
  target's `field` form. The record states the member's name in the request
  and does not record it again. Go's fragment is one field of the type
  `BittypeOptional[T]`; TypeScript's is one optional property.
- *The conversion half* (row 13). Built as a field-state protocol on both
  sides. See "Row 11: what the protocol needed" below.
- *"A contribution with no standalone value adapter"* (row 14). Expressed by
  the product: it requires another facet kind of a contribution at a field
  position than of a use. The optional contribution has a field-state facet
  and no value facet, in both targets.

**What the page left open and what was chosen.**

- *The class of the member-fragment refusal.* **Unsupported**, with the
  owner "member vocabulary of <target>" and the kernel's reason. Why: the
  party that does not admit the fragment is the target's member vocabulary,
  and a later slice of that vocabulary can admit the same fragment. The
  criterion is the page's own for unsupported ("A later slice can admit the
  same input"). It is not "refused by design": the optional contribution's
  domain does cover the field position. A target that states no member
  vocabulary refuses every field position the same way.
- *Whether a member fragment may hold more than one member.* The page: "The
  target's member vocabulary decides". The Go and the TypeScript model of the
  prototype check each member alone (`model.CheckMember`), so a fragment of
  two fields is admitted. The Go lowering of the prototype handles one member
  for each link only.
- *How a target states its member vocabulary.* An optional interface of the
  contract, `generation.MemberVocabulary[T]` with `CheckMember`. The Go and
  the TypeScript target implement it by their model's field form.
- *The law of a field-state protocol.* Not fixed. The prototype's vectors
  bound it by the two laws of the record that contains the field: each state
  that the record's canonical value can have converts there and back as
  itself, and a native state that is no value of the field's type is refused.

**Row 11: what the protocol needed.**

| | Go | The TypeScript stand-in |
| --- | --- | --- |
| The native fragment | One member of its own type, `BittypeOptional[T]` | One optional property `name?: T` |
| Operations | `Construct(present, abstract value) -> member`, `State(member) -> (present, abstract value)` | `write(container, member key, present, canonical value)`, `read(container, member key) -> {present, value}` |
| What the record's factory had to know | Which abstract field belongs to the position (from the source snapshot), and the member (a member reference). It hands over a member *value* and takes one back | The same, and it hands over the *container object and the member's key*, because an absent property is a fact about the object and has no value |
| What the optional contribution had to know | The adapter selected for its operand, and the declaration it owns. Not the record, not the member's name, not the canonical name | The adapter selected for its operand. It is given the container and the key at each call |
| The four states | Absent, present with a value, present with null (refused by the operand's adapter). A present undefined does not exist: a Go member is a value of its type | Absent, present with a value, present with null (refused by the operand's adapter), present with undefined (refused by the protocol) |
| What neither side needed | The payload of the other contribution; a shared helper key; a cross-field layout | |

Delegation of presence was workable in both. The one asymmetry is the slot's
shape: in Go the protocol can be value-in, value-out, and in TypeScript it
must be given the container and the key. A slot kind that is to serve every
target therefore has the TypeScript shape, of which the Go shape is the
special case where the container's member is a value. The decision of row 11
is the research owner's, with the service library's owner; this is evidence
for it and no decision.

**Findings of this stage.**

- **F19. The role of the optional flag in the native field form.** The Go
  model's `field` form carries an "optional" flag as data, and today's Go
  record passes `true` together with the wrapped type. To keep the native
  text unchanged, the Go optional contribution writes the same flag into its
  fragment. The flag is then the optional contribution's own statement in the
  member vocabulary, and the record no longer knows it. Nothing reads the
  flag in the Go target. Seen, not assessed further.
- **F20. Under D3 a new package appears for each target, and tests that list
  mapping packages change.** See F4. The same will hold for Rust and C++.

**Tests** (`integration/go/prototype_stage6_test.go`; evidence
`stage6-run2.txt`; `stage6-run1.txt` is the first run with one wrong
expectation of the test: the target's identity is `go`).

| Test | Shows | Does not show |
| --- | --- | --- |
| `...Stage6FourStatesOfAnOptionalField` | In lowered Go and in the stand-in: absent and present with a value convert there and back as themselves; present with null is refused; in the stand-in a present undefined is refused, and the vector is red when the protocol reads it as absence | A null that a type admits (#211); the service library's slot (row 11); independent authorship |
| `...Stage6FieldPositionKeepsTheScopeCheck` | A request with the operand of another record, or with a role issued for another pair, fails. Red under a kernel without the check | That the source issues the role for optional fields only |
| `...Stage6MemberFragmentOutsideTheVocabulary` | A fragment that is no member is refused as unsupported with the member vocabulary as owner; a target without a member vocabulary refuses a field position; a fragment of two members is admitted | Cross-field layout (out of the domain) |

## Stage 7: a second reader in another process (gate line 4; row 26a)

**Built.** The TypeScript stand-in as a reader
(`completeness/typescript/standin/reader.mjs`), run with node from a Go test
with three files: the serialized plan, the source snapshot and the vectors.
Its factories (`structural.mjs`) and a decoder for the canonical bytes of a
payload (`ontos.mjs`). The Go backend of stage 5 is a second reader of the
same inputs in a fresh Go process.

**Row 26a: a facet value crossed the process boundary.** Yes: the facet that
names a run-time factory. The stand-in finds a factory by what the plan entry
names (`docs/representation-families.md:49-50`), so the factory reference has
to be in the serialized plan. The transport descriptor the prototype designed
(`generation.FacetTransport`):

- A facet kind may state a transport: a media type and an encoder that the
  kind's owner supplies.
- When the planner writes a selected contribution into the plan, it takes
  each facet that a requested product requires of that contribution. If the
  facet's kind has a transport, the planner calls the encoder and stores the
  kind's reference (identity and contract version), the media type and the
  encoded bytes with the contribution's entry. It reads none of it.
- A reader finds the facet by the kind's identity, checks the kind's contract
  version and the media type, and decodes the value itself.
- A facet kind without a transport never leaves the process: the Go lowering
  facets, which are Go functions, do not.

The value that crossed is `{module, export, runtime}`: the factory reference
and the runtime compatibility it requires (the fourth version axis).

**What the reader needed, and where it came from.**

| The reader needed | It came from |
| --- | --- |
| Which contribution was selected for a node, with both revisions | The plan: the node's contribution reference |
| The factory of a contribution | The plan: the transported facet of the contribution |
| The recorded decisions | The plan: the node's payload, canonical ontos bytes, which the stand-in decodes itself |
| That the payload is the validated one | The plan: the receipt, which the stand-in recomputes |
| The member's name of a delegated field | The plan: `stated` on the link (O3) |
| Which abstract field a link of a record is | **Not the plan alone**: the link's `position` and the snapshot's description of the declaration (F17) |
| The canonical field names | The snapshot (source meaning) |
| The adapter of a binder | The plan: the application's adapter arguments |
| How a receipt is computed, and how node identities are computed | **Not carried**: the stand-in has the derivation of a receipt written into it. It is part of what the envelope version means. The stand-in does not recompute node identities |
| That the plan asked for bindings | The plan: `products`, and the adapter arguments being there |

So the plan did not carry everything a reader needed in two places: the
pairing of a link with a source child (now carried as `position`), and the
derivation of the receipt (part of the envelope protocol, to be written down
with #195).

**Findings of this stage.**

- **F21. A reader checks three things per contribution, and the plan has one
  place for each only since this stage.** Payload schema version, both
  revisions of the reference, and the runtime compatibility. The first two
  are in the plan's contribution entry. The third is in the facet. A reader
  that has no facet for a contribution cannot check the fourth axis for it;
  it fails as "missing facet" first.
- **F22. "Ground descriptors" for the stand-in.** As F16: the stand-in walks
  the open declaration `Page<T>` with the adapter argument of the
  application, and needs the description of `Page` as written, not the ground
  descriptor of `Page<Person>`.

**Tests** (`integration/go/prototype_stage7_test.go`; evidence
`stage7-run1.txt`, `stage7-plan.json`, `stage7-snapshot.json`). The vectors
of the four states are in stage 6.

| Test | Shows | Does not show |
| --- | --- | --- |
| `...Stage7StandinRoundTripsPageOfPerson` | The stand-in converts `Page<Person>`, a nested page and a list there and back from the plan and the snapshot alone, and refuses two values outside the type. Every selected contribution crossed with one transported facet | The service library's factory, its budgets, call context and memo (row 31a); independent authorship |
| `...Stage7StandinRefusesUnknownVersions` | An unknown envelope protocol, payload schema version, representation revision, runtime compatibility and source semantics are each refused under their own name | A reader that knows several values of an axis |
| `...Stage7ReadersNeedTheMatchingSnapshot` | Without the snapshot the stand-in and the Go backend fail explicitly. A snapshot of another pin is a linking error in both | A reader protocol for source semantics (row 34a) |
| `...Stage7StandinDerivesMemoKeys` | Equal keys for two uses of `Page<Person>`, another key with another selected dependency, and the same key as the Go reader derives from the same bytes | A memo (bitspark-service#40) |
