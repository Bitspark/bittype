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
