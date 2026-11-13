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
