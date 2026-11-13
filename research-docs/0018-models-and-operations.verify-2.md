# verify-order: rows R21 to R31 of research 0018

The lists behind the counts were kept on the workstation that ran the verification and are not part of the repository.

**State read.** bittype `origin/main` is 3998db5. `origin/feat/contribution-contract-196` is af1dcde, 10 commits ahead of main with merge base 3998db5 (the document names 8d10e00). `origin/feat/other-pairs-196` is 447c44c and lacks the migration's last commit. `origin/feat/names-per-file-273` is cbe0291. `docs/contribution-contract.md` is identical on main and on the migration branch. Nothing was built or run; every import fact is textual.

## 1. Contradictions, first

1. **Language first below `model` (R21) contradicts a standing layout rule and the owner's sketch.**
   - `LAYOUT.md:134-135`: "Do not create a language-first packages tree, additional Go modules, compatibility aliases or placeholder source directories."
   - bittype `AGENTS.md:236-237`: "Do not add compatibility aliases or a language-first packages tree."
   - The owner's message 7 sketches `model/{syntax,semantics,...}/go/...`, sort first.
   - The advice gives no reason for language first, and its own A4.3 says the opposite: "organized around the kinds of models being declared, with language identity as a further dimension". R21 and R22 as extracted carry both orders.
2. **The advice shows a planned wait that AC-10 does not allow.** Section 6.5 of the document says the first step "is taken when the migration's branch no longer edits those files". AC-10 (`architectural-correctness.md:345-348`): "Deferring the move itself is permitted only when its destination needs a huge effort. Record such a deferral in the case with its reason and the trigger that ends it". No case records it; the corpus ends at ACX-012.
3. **R24 against AC-11.** Read as one checker per language, "an operation-side checker" would be a per-language copy. `pairs-report.txt` AG-7 already records five near copies of the intended-forms mechanism.
4. **R24 and R28 against R17.** The first move needs declared form descriptions, and none exists on main. R17 puts intent schemas last (syntax, then facts, then intent), while R28 moves the intended-forms packages first.
5. **R25 against three accepted taxonomies.** Its "rejected" means a policy refusal; the accepted `Rejected` means a refutation by evidence. It has no limit category, and every accepted taxonomy has one.
6. **R30 rule 2 cannot be carried by an import test.** The advice lists it under what "the import test should enforce". A check that tried would be the keyword heuristic the owner forbids.
7. **R27 against the foundation role and one `LAYOUT.md` sentence.** The foundation role is cut by the root (carriers to model, contracts to ops). `LAYOUT.md:145-146` says "Directory groups are not Go packages or an assertion that all members share the same dependency level"; a root that means "model imports no ops" is exactly such an assertion.
8. **Three details of the advice are not true of the code.**
   - Its tree gives cpp `syntax/` and `semantics/`. cpp has a syntax vocabulary and an intended-forms package, and no semantic vocabulary or symbol package.
   - It lists `ops/check/files`. There is no Files check stage.
   - It says "build" the two pairs. They are built: draft PR #311 on the migration branch.

## 2. The rows

### R21. Refined P3: two roots, `model/shared`, `model/language/<language>/...`, `model/relations`, `model/contributions`, `ops/...`

**Verdict: NEEDS ADAPTATION.**

**One module.** `git ls-files | grep -E '(^|/)go\.(mod|work)$'` gives `go.mod` only. `integration/go` is not its own module (`LAYOUT.md:140-141`: "no nested Go module is committed"). The layout needs no second module, so the rule "One public root Go module" (`LAYOUT.md:6`) holds. Its cost: the module graph cannot carry rule 1; only the import test and Go's `internal` rule can.

**`internal/` below `ops`.** Go lets only code under `ops/` import `ops/internal/...`.
- Five production imports break if `internal/` moves under `ops` before the splits:
  - `trees/representation/go` imports `internal/tree` (ast.go, cst.go, semantics.go);
  - `semantics/rust/go` and `transform/syntax/rust/symbol/go` import `internal/rust`;
  - `semantics/typescript/go` and `transform/syntax/typescript/symbol/go` import `internal/typescript`.
- All five are code that R7d and R7h place in ops. After those splits the compiler itself enforces "model imports no private part of ops".
- Compositions import helpers too: `integration/go` imports `internal/horosfiles` and `internal/tree`; `tools/horos-module/go` imports `internal/horosmodel` and `internal/lex`; `tools/provider-matrix/extract/grammar/go` imports `internal/grammar`. With helpers at `ops/internal`, integration and tools must be under `ops/`. That part of the expert's tree follows from Go's rule.
- The model root has no place for a private helper of its own. `model/internal` would be invisible to ops.

**Fixed points touched.**

| Fixed point | Quote | Effect |
| --- | --- | --- |
| 20 | "A package that names a language imports no package of another language. A mapping names two." | Kept. `model/contributions` adds model-side packages that name two languages; the sentence needs that case |
| 21 | Go's `internal` rule | As above: five imports, and tools and integration forced under ops |
| 22 | "One public root Go module" | Kept. "placeholder source directories" forbids creating `model/relations` empty |
| 23 | "Do not create a language-first packages tree" | Contradicted by `model/language/<language>/...` |
| 24 | "organize by processing stage and enforce a uniform structure across languages"; `<stage>/<dimension>/<scope>/<component>/go` | Kept below ops. The expert's tree drops the trailing `go` and gives no reason |
| 25 | "A checker gets no concrete tree, a formatter no lexer, and no Syntax stage a semantic contract." | Untouched |
| `LAYOUT.md:131-133` | "Independent consumer fixtures, stage composition, language conformance and native compiler qualification remain under integration." | Becomes `ops/integration` |
| `LAYOUT.md:145-147` | "Directory groups are not ... an assertion that all members share the same dependency level." | No longer true of the two roots |

**The placement table.** Command for the 129: `git ls-tree -r --name-only origin/main | grep '\.go$' | grep -v '_test\.go$' | grep -Ev '(^|/)testdata/' | xargs -n1 dirname | sort -u`. The right column is my derivation from the tree in A4.1; the expert gave directories, not a mapping. Where a package splits outside R24, the line is another verifier's (R7a to R7i); I state the parts from the code study and spot checks.

Carriers (6)

| Today | Under the expert's layout |
| --- | --- |
| `trees/ast/go`, `trees/cst/go` | `model/shared/tree`, whole |
| `trees/representation/go` | Splits. `Language` and the envelope types to `model/shared/representation`. `Pack`, `Unpack`, `PackCST`, `UnpackCST`, `PackSet`, `PackSemantics`, `UnpackSemantics` to ops, no directory named. The envelopes have private state, so an ops package cannot construct them |
| `files/go` | `model/shared/files` (`Capture` validates: R7a) |
| `semantics/go` | `model/shared/graph` (`New` validates: R7a) |
| `source/go` | **No place named.** `model/shared` lists tree, graph, files, representation. Types would be a fifth entry; `Publisher.Publish` (source.go:254-448) to ops |

Vocabularies (16)

| Today | Under the expert's layout |
| --- | --- |
| `syntax/{cpp,golang,horos,rust,typescript}/go` | `model/language/<language>/syntax`, whole |
| `semantics/horos/go` | `model/language/horos/semantics`, whole |
| `semantics/{golang,rust,typescript}/go` | Split. Names, profiles, contract identities, input types to `model/language/<language>/semantics`. Typed readers and the settings encoding to ops, no directory named |
| `transform/syntax/{cpp,golang,rust,typescript}/model/go` | Split (R24). Declarations to `model/language/<language>/intent` |
| `transform/syntax/{golang,rust,typescript}/symbol/go` | Split. Symbol types to `model/language/<language>/symbols`. Admission and issuing (`Predeclared`, `ExternalUse`, `ReadExternal`) to ops, no directory named |

Contracts (13)

| Today | Under the expert's layout |
| --- | --- |
| `contracts/language/go` | Model by R8 (identity only). **No place named** |
| `contracts/{text,pipeline,check,print,parse,format,transform,verification}/go` | `ops/contracts/<name>`, whole |
| `contracts/{semantic,resolution,analysis}/go` | `ops/contracts`. Rule 1 forces a split: the semantic vocabulary packages declare identities as values of `Contract`, `Requirement`, `Obligation`, `Setting`, `EvidenceRequirement` |
| `contracts/generation/go` | Splits. Intent data types (`Form`, `FormIssuer`, `Expression`, `NativeDeclaration`) are the "shared carrier" of R15, **no place named**. `Target`, planner, sealing, source walk to `ops/contracts/generation` |

Kernel (7): `engine/go` to `ops/engine`; `transform/syntax/shared/{generation,target}/go`, `check/syntax/shared/{ast,cst}/go`, `print/{syntax,files}/shared/provider/go` to `ops/<same path>`.

Stage providers and query services (28), all to `ops/<same path>`:
- whole: `parse/syntax/{cpp,golang,horos,rust,typescript}/provider/go`, `parse/files/horos/provider/go`, `format/syntax/{cpp,golang,horos,rust,typescript}/provider/go`, `format/files/{horos,typescript}/provider/go`, `resolve/files/{golang,horos,rust,typescript}/provider/go`, `analyze/files/{golang,horos,rust,typescript}/provider/go`, `analyze/files/{golang,typescript}/queries/go`, `check/syntax/horos/provider/go`;
- split by R7c: `check/syntax/{cpp,golang,rust,typescript}/provider/go`. The grammar table (`rules`, for Go at provider line 17) goes to the language's syntax model.

Generation by language and pair (30)

| Today | Under the expert's layout |
| --- | --- |
| `transform/syntax/{cpp,golang,rust,typescript}/target/go` | `ops/<same path>`; receives the coverage decisions |
| `transform/syntax/horos/source/go`, `transform/syntax/horos/mapping/{declarations,naming,pattern}/go` | `ops/<same path>` |
| `transform/syntax/horos/mapping/go` | Splits. Two profile identities to model, area not said. `Library[T]` to ops (ACX-010 already moves it) |
| 21 mappings: `horos-to-{cpp,rust,typescript}/mapping/{boolean,declarations,list,record,text}/go`, `horos-to-golang/mapping/{boolean,declarations,list,optional,record,text}/go` | `ops/<same path>` |

Private helpers (16), to `ops/internal/<name>`:
- whole: `internal/{tree,lex,parseutil,concrete,golangsem,rustsem,typescriptsem,horosmodel,horosfiles,typescriptimports}`;
- `internal/grammar` splits: table types and entry constructors to model, **no place named**; `Check` (grammar.go:80) to ops;
- `internal/{cpp,horos,rust,typescript}` split by R7h: word lists and tables to model, **no place named**; predicates and literal decoding to ops;
- `internal/declsyntax` splits: the horos table to model, `Check` to ops.

Verification (7): `verification/files/horos-to-go/go`, `verification/files/horos-to-go/mappings/{list,optional,record,scalar}/go`, `verification/files/source/horos/go`, `verification/logos/go` to `ops/verification/...`.

Tools (6): `tools/horos-module/go`, `tools/provider-matrix/{go,report/go,capability/go,observation/go,extract/grammar/go}` to `ops/tools/...`.

Test-only (3): `integration/go` to `ops/integration`; `verification/integration/go` and `tools/provider-matrix/extract/structfixture/go` follow their parents.

Sum: 11 packages wholly model, 27 split, 91 wholly ops.

**The adjustment.**
- Below `model`, sort first then language, as the owner's sketch and the standing rule have it. Sort first also keeps a sort's agnostic carrier beside its languages, as `semantics/go` sits beside `semantics/golang/go` today.
- Name the six missing places before any move: the Source carrier, language identity, the intent carrier, the grammar-table type, lexical constraints, and the ops halves of split vocabulary packages.
- Keep `internal/` at the module root until the five imports are split away.
- Keep the trailing `go` unless a reason is given.

**Not established:** whether `model/relations` has any content today. Nothing on main or on the branches is a declared relation in the sense of R10.

### R22. The stage matrix stays below `ops`; the model root is organized by sort with language as a further dimension

**Verdict: HOLDS.**

- The first half keeps the decision of 31 October (`LAYOUT.md:3-5`) below ops and matches the owner's `ops/{transform,...}/go/...`.
- The second half is sort first, which is what R21 should say.
- Concretely: 91 packages take the prefix `ops/` and nothing else. `model` and `symbol` leave the stage directory `transform/`.
- `TestLanguageStageMatrix` (`architecture_test.go:15`) lists `syntax` as the first column of `nativeStages`. That column moves to the model root, so the matrix check reads both roots.
- STALE detail: the `ops/check/files` of the tree.

### R23. Contribution identities and decision schemas: identities independent of the implementation, under a model area; codec and mapping are ops

**Verdict: NEEDS ADAPTATION.** The first half is already true on the migration branch.

- No identity is a Go path. The strings are `horos-to-golang/record`, `bittype.payload/horos-to-golang/record`, `bittype.structural/closed-record`, `bittype/golang/record-form` and the publisher `github.com/Bitspark/bittype`. Nothing derives one by reflection: `git grep -E 'PkgPath|reflect\.TypeOf'` over the branch's production files is empty. A package move changes no identity, pin or serialized plan.
- Identities already sit apart from the implementation: 24 vocabulary packages and 6 family packages, with 0 functions (`git grep -E '^func '`). The role forbids them the mapping and the codec.
- They import `contracts/generation/go` for the types `ContributionReference`, `ContributionIdentity`, `Meaning`, `SchemaReference` (30 imports). Under rule 1 those types and `EvidenceKind` need a model-side owner first. That is the split of R7e, in the package the branch changes most (13 files).
- "No assumption that horos is the source" has residue recorded in `pairs-report.txt`:
  - the shared codec part sits at `transform/syntax/horos/payload/go` though it serves every pair;
  - the publisher has three spellings (AG-6);
  - a family's meaning is undefined for a source that is not horos (F8).

**The adjustment.** Decide the owner of the identity types before placing the 31 packages. Say where a family of one language lives: the advice does not tell `model/relations` from `model/contributions`.

### R24. The first move: coverage to the target, well-formedness to an operation-side checker, form descriptions to the model, constructors by declare or construct

**Verdict: NEEDS ADAPTATION.** The premise is true: mappings use constructors and names and never a check; targets use names and checks and no constructor (`git grep -o 'model\.[A-Z][A-Za-z]+'` over mappings and targets).

`transform/syntax/golang/model/go/model.go` on main, piece by piece:

| Lines | Piece | Under R24 |
| --- | --- | --- |
| 12 | `var forms = g.NewFormIssuer[s.Language]()` | **The test gives no clean answer.** It is the vocabulary's identity and an issuing mechanism of the contract package at once. R19 says identity must not rest on a Go marker, so it leaves the model |
| 18-26 | `FormPrimitive` ... `FormDefined` | Declaration, stays |
| 30-33 | `PrimitiveString`, `PrimitiveBool` | Declaration of names, stays |
| 35-36 | `Form`, `Name` | Expression machinery. Bound to `forms` by Go visibility; goes with it |
| 37-51 | `Primitive`, `Collection`, `Field`, `Record`, `External` | Operation: the advice's table says "Building intent expressions: Operation code". Each signature is also the only typed statement of a form's shape |
| 53-59 | Context check; foreign form gives `ErrScope` | Checker, agnostic |
| 62-65 | Primitive: arity | Checker, from a table |
| 63 | `d[0] != PrimitiveString && d[0] != PrimitiveBool` | Coverage, to the target. Returned as `ErrInvalid` today |
| 66-73 | Slice and field arity, non-empty name, boolean flag | Checker, from a table |
| 74-86 | Record: members are fields, no name twice | Checker. A relation among siblings, so code unless the table can state it. Identical in all five models, so agnostic |
| 87-88 | External gives `ErrUnsupported` | Coverage, to the target |
| 89-90 | Unknown form gives `ErrUnsupported` | Against a declared table this is "not declared". R25 would call it invalid; the Files rule calls it unsupported |
| 102 | Empty name; form is neither alias nor defined | Checker |
| 105-110 | Defined needs a record body | Coverage, to the target (R14). `ErrInvalid` today |

Added by the migration (`model.go:114-127` on the branch):
- `EvidenceRecordForm`: an identity, a declaration; stays on the model side.
- `CheckMember`: a check, operation side. It is a third statement in code of "a record's members are fields", after lines 79-85 and the target's switch.

The target (`transform/syntax/golang/target/go/provider.go`):
- `CheckExpression` and `CheckDeclaration` (24-29) forward to the model.
- `Backend.Expression` (111-128) checks again and reads data by position.
- `Declaration` (98-99) reads the form name.

It would call the checker with Go's table and add its three coverage decisions.

**What the checker is as a package.** By AC-11 it is one agnostic mechanism over a per-language declared table, not one per language. The kernel package `transform/syntax/shared/target/go` already holds such a check: every target's `CheckProgram` calls `m.CheckSingleScope` (provider.go:130-132). The target role already permits that import (`architecture_roles_test.go:356-358`), so no new edge is needed. The planner still reaches the check only through the `Target` interface.

**The constructors need a home that mappings may import.** A mapping may not import a target (bittype `AGENTS.md:190`: "Mapping packages depend on native models, never target AST constructors"). So they need an ops-side package per language, or they stay as typed declarations of shape.

**How many packages.**
- On main: 4 model, 4 target, `contracts/generation/go`, `transform/syntax/shared/target/go`. That is 10 existing packages, plus one new model-side home for the intent carrier and table type.
- On the pairs branch: 2 more, the horos model and target.
- The migration touches all 10.

**The adjustment.**
1. One checker, in the kernel, over declared form tables.
2. The form table is designed in this move, so intent schemas come first, not third.
3. The issuer and the constructors get a stated home.
4. The change of category at lines 63 and 105-110 from invalid to unsupported is made on purpose and recorded, since AC-05 names failure meaning.

### R25. Every operation distinguishes invalid, unsupported, incomplete, rejected

**Verdict: NEEDS ADAPTATION.** The distinction between invalid and unsupported is accepted and implemented. The uniform set of four is not, and it adds a concept without a case the existing contracts cannot express.

What exists:
- **Files formatting** (`docs/source.md:498-507`, 13 November): "Invalid input", "Unsupported coverage", "Missing evidence", "Exceeded limit", with the rule "a refusal is unsupported coverage when a later slice can admit the same input, and it is invalid input when no later slice can". Values: `source.ErrInvalid`, `source.ErrUnsupported`, `semantic.ErrEvidenceUnavailable`, `source.ErrLimit`.
- **Resolution and analysis** (`contracts/semantic/go/contract.go:46-61`, `docs/analysis.md:300-303`): completion (`Complete`, `Incomplete`, `Unsupported`) apart from the verdict (`Undetermined`, `Holds`, `Rejected`). "Report separates completion from truth."
- **Generation** (`contracts/generation/go/contracts.go:14-21`, `source.go:83-92`): `ErrInvalid`, `ErrUnsupported`, `ErrScope`, `ErrLimit`, `ErrIncomplete`, `ErrAmbiguous`, `ErrRejected`, and the verdicts `Admitted`, `Unsupported`, `Rejected`, `Undetermined`. The contribution contract refines these to eight refusal classes.
- **Parse and check** (`contracts/text/go`): the codes `syntax`, `unsupported`, `invalid-ast`, `invalid-cst`, `limit`.

Against the four:
- "Invalid" and "unsupported" match.
- "Incomplete" matches missing evidence, `Incomplete` and `Undetermined`.
- "Rejected: a policy explicitly refuses" does not match the accepted `Rejected`, "the evidence refutes the subject's admissibility". The contribution contract page (table at 645-653) calls the policy case "Refused by design" and says "It is not `ErrRejected`".
- The four have no limit. `docs/analysis.md:367`: "Exhaustion cannot masquerade as an empty valid model or a rejected type."

**The adjustment.** Keep each operation's accepted classes. Take the principle only. Apply it where it bites: the two coverage answers of R24 that return `ErrInvalid`.

### R26. A model is not narrowed to what an operation supports

**Verdict: HOLDS.**

- `syntax/golang/go/nodes.go:51-55` declares pointer, array, map and channel types.
- The Go target writes ten kinds (`git grep -o 's\.Kind[A-Za-z]+'` on the target).
- AC-06: "Keep supported feature coverage separate from the intended extension boundary."
- Concretely, `model.go:63` and `:105-110` narrow, and they leave with R24.
- The horos model on the pairs branch says of itself "Coverage of this model, and no rule of horos" (model.go:79-81).
- A2.4 of the advice calls the intent model "the target language's supported type-intent vocabulary", which reads against this row. A missing form is missing coverage; a refusal inside the model is narrowing.

### R27. Before any move: the vocabulary, the ownership rules, the roles stay below the roots, U1

**Verdict: NEEDS ADAPTATION.** Deciding first is right: AC-10 asks that the first change "establishes the final structure", and a dependency check "enforces a boundary that reasoning has already decided".

"The existing roles remain below the two roots" is not true of all of them:

| Role today | Under the two roots |
| --- | --- |
| Foundation | **Cut in two.** Carriers to model, the contract packages to ops. `foundationRules` (31 entries) holds foundation, kernel and shared helpers in one closed table. Five of its edges would run from model to ops: `trees/representation/go` to the two checks, to `contracts/language/go` and to `internal/tree`; `source/go` to `contracts/semantic/go` |
| Kernel | ops, whole |
| Shared private helper | `ops/internal`; `internal/grammar` splits |
| Vocabulary | model. Three of its four sorts split |
| Contribution, common mapping part | ops, whole |
| Private part | `ops/internal`; the four lexical helpers split |
| Composition | `ops/integration`, `ops/tools` |

The two halves of a vocabulary package that holds checks:
- The declaring half stays a vocabulary. For the first time the definition "The names two parties must agree on and nothing else" would be true of it.
- For the intended-forms packages the other half is kernel (the one checker) and contribution (coverage in the target). No new role.
- For the typed readers and for symbol admission no role fits. They are specific to a language, not selected by a consumer, and imported by several roles, so "private part" is wrong. A role must be named, or they join the query service of their language.

The eight terms: most map onto words the repository has (vocabulary, carrier, contract, snapshot, receipt, envelope, admission, coverage, evidence). The derivation rule asks to map before adding.

**The adjustment.** Decide the roles per root, with the foundation table divided and the role of the reading half named. Amend `LAYOUT.md:145-147`.

### R28. Then move misplaced responsibilities as independent changes, intended-forms first, then protocols from carriers

**Verdict: HOLDS.** AC-10 says "as its own change". The platform page says "Correct the grouping first, as AC-10 orders" (`:557-558`).

- The first move is the one that collides with the migration in every package (R24, R29).
- The later splits do not: the migration touches 37 of the 129 packages on main and leaves 92 alone, among them every carrier, every `syntax` and `semantics` vocabulary, the three symbol packages, 12 of the 13 contract packages and all of `internal/`.
- No rule orders two independent corrections. "First, then" may run in parallel.

### R29. Do not merge the work in flight into the present structure; retarget the migration

**Verdict: NEEDS ADAPTATION.** The order holds. Three premises of the advice do not.

**Facts: what "retarget" would mean.** The migration adds 62 packages and removes one; its diff is 282 files (200 added, 52 modified, 28 deleted, 2 renamed). The 98 packages it adds or touches, by `mig-pkgs-all.txt`:

| Count | What | Under R24, R7, R23 |
| --- | --- | --- |
| 4 | Intended-forms packages | Split |
| 1 | `contracts/generation/go` (2,135 lines on main, 4,025 on the branch) | Splits |
| 4 | Targets | Receive coverage and the checker call |
| 31 | 24 contribution vocabularies, 6 family vocabularies, the mapping identities | Declarations only. Leave the stage directory for `model` |
| 58 | 24 codecs, 24 mappings, 3 shared policies, the shared codec part, 2 kernel packages, the source reader, 3 tools | Placed correctly apart from a path prefix |

**Facts: is the root move mechanical?**
- Given a table from old path to new, yes for: 2,121 import strings in 365 Go files on main (959, 871 and 291 in production, test and fixture files); 465 tracked files under the 18 code directories; 36 path mentions in three workflow files.
- Not mechanical:
  - the table itself (the classification);
  - the 27 packages that split;
  - `packageRole`, a switch over segment count and position, with 1,166 quoted paths in 17 test files (638 in the synthetic sets);
  - the five tables (31, 12, 5, 12 and 14 entries);
  - `internal/`, which cannot go under `ops` before five imports are gone.
- The branch has 2,474 import strings in 437 files; the pairs branch 2,701 in 480.

**Facts: does anything move twice?**
- If the migration merges first:
  - its 62 new packages are delivered at paths the direction already supersedes and are renamed once afterwards;
  - `CheckMember` is delivered into four model packages and moved out;
  - 31 identity packages are re-placed.
- No identity changes in either order.
- The page says of itself (`contribution-contract.md:2990-2992`): "the contribution of line 7, pinned at the commit of the migration, still builds at the head of each later family". A root move after the gate breaks that pin by import path alone.
- If a split package is moved whole and split later, the split-off half moves twice. So splits precede the root move for those packages.

**Facts: the reverse cost.** If the move lands on main first:
- 77 existing Go files that the branch modifies, deletes or renames are also touched by the move; 71 of them hold a module import string;
- of the 47 modified in place, 41 are rewritten by the move, and in 14 the branch itself changed an import line;
- 102 of its 152 added Go files need the same rewrite;
- its three architecture test files carry 246 added lines that the role rewrite replaces.

**Facts: does the migration add misplacement?**
- Operation code inside packages that will be models: `CheckMember`, 6 lines in each of 4 packages, 24 lines. Main already has 250 lines of checks there (61, 63, 63, 63).
- The pairs branch adds a fifth intended-forms package with 52 lines of checks.
- Model-side packages importing a contract package: 30 new, against 33 on main.
- The contract package nearly doubles while holding both sorts of thing.
- The branch's last commit states the choice: "the models' member vocabulary and evidence kind are each model's own answer".

The strongest argument for the expert's order is therefore small in lines and real in kind.

**Facts: what can be split on main now without touching a file the branch changes.**
- No branch file imports these: `semantics/rust/go`, `internal/grammar`, `internal/rust`, `internal/horos`, `internal/declsyntax`, `check/syntax/shared/*`, `contracts/verification/go`.
- A few branch files import these: `source/go` (2), `files/go` (1), `semantics/golang/go` (2), `semantics/typescript/go` (4).
- `trees/representation/go` is imported by 56 branch files, for `Language`. A split that leaves `Language` in place does not collide.
- Not independent: the four intended-forms packages (every importer is a branch file: 11 of 11, 10 of 10, 8 of 8, 7 of 7), the four targets, `contracts/generation/go`, `transform/syntax/shared/target/go`.

**Facts: the state of the gate.**
- `contribution-contract.md:2886`: "The migration of step 4 does not merge while one line is red."
- The branch has four `gate_*` test functions for parts of line 1; two subtests skip as "IN WORK".
- PR #308 is a draft.
- Line 6 needs the reader of #273, which is open as PR #304 and changes `transform/syntax/shared/target/go` and the TypeScript target.
- The pairs wait for #310.
- The Go and Rust Files formatters exist on no remote branch.

Merging first is therefore not available today, whatever is ruled.

**The three premises that do not hold.**
1. "Merely to move it again": 58 of 98 packages change by a prefix only, and no identity changes.
2. The pairs are built, not to be built.
3. "The agreed architecture" does not exist yet: no accepted page holds the test, the layout or the roles.

**The order I would take, from the rules.**
1. Decide and record the test, the layout below the roots and the roles (R27). Documents only; the branch collides in `LAYOUT.md` alone.
2. At once, on the migration and pairs branches: `CheckMember` leaves the model packages and the target states it. The contract's own wording calls it "the target's member vocabulary", and the interface is implemented by the target. This needs no decision of research 0018.
3. On main, each as its own change, the splits that touch no branch file. Each split-off part is created at its final path, so it never moves again.
4. On main, as its own change, the first move of R24. The branch then merges main.
   - AC-10 asks its own question: "After the first change, does any planned step have to move, rename or remove what that change delivered?" Only main first answers no.
   - Waiting for the gate is a deferral that AC-10 permits for a huge effort only, and this move is 10 packages.
5. On main, the scripted move of the whole packages, after steps 3 and 4 and before the migration merges. The same script runs on both branches.
6. The gate is shown on the new boundaries, the migration merges, then the pairs.
7. The Files formatters for Go and Rust are written at final paths from their first commit.

**What line 7 means for the order.** The gate names its shared packages by path (`:2958-2972`), and "the native model and target of each language" is on the list. R24 changes exactly those. Shown before R24, line 7 would establish "no shared package changes" for a boundary about to be redrawn; AC-07 lets evidence support only the claim it establishes. The page also makes the migration's commit the point of finality (`:1861-1863`, `:2995`). So the boundaries the gate names must be the lasting ones when it is met. Line 7 is no freeze against corrections: it measures the diff of adding a contribution.

**Not established:** whether git's rename detection carries the branch's edits through the move (no merge was tried), and how long the gate needs.

### R30. The import test enforces five rules

**Verdict: NEEDS ADAPTATION.** The mechanism reads the import block and the selected symbols of every production file, also behind build tags. A rule comes from a path shape (`packageRole`) or one of five tables with 74 hand entries on main.

| Rule | Can an import check carry it? |
| --- | --- |
| 1. No import from ops into model | Yes, as one prefix rule beside "a foundation or kernel package imports only foundation and kernel packages" (`architecture_check_test.go:230-233`). Test files stay out by design |
| 2. No operation implementation under model | **No.** `CheckExpression` in the Go model imports nothing forbidden. The platform page: the check "cannot see a decision made inline". A check on function names or on the presence of functions is the keyword heuristic the owner excludes, and it would be wrong: `semantics/horos/go` states names as 26 one-line functions |
| 3. No dependency across languages except declared ones | Yes; it exists (`roleRule.scope`, `languageOwners`). The language segment moves in the path |
| 4. No private operation implementation imported by model | Yes: a subset of rule 1, and the compiler enforces it once helpers are at `ops/internal` |
| 5. No cycles between model declarations | Vacuous: Go refuses import cycles. As a layering among model areas, the per-package lists of `foundationRules` can carry it |

**The adjustment.**
- Rule 2 stays in review. The advice says so in the same section: the test "cannot decide whether a package has been conceptually classified correctly".
- Each added rule gets a test that breaks it, as the platform page asks.
- "A small number of explicit exceptions" starts from 74 entries.

### R31. No approval for a migration that merely moves packages under two directory names

**Verdict: HOLDS.**

- A prefix-only move cannot satisfy the direction: wherever the 27 mixed packages land, rule 1 or rule 2 is broken.
- `architectural-correctness.md:525`: "Moving, renaming or generating closed case handling does not correct ownership."
- The scripted move is still needed as one step. For 91 packages it is all that happens.

## 3. What the advice did not answer, and what two rows give together

**Not answered, and my rows needed it:**
- When the root move itself happens. The advice's eight steps contain no such step.
- What one module costs for enforcing rule 1. Fixed point 22 asked.
- How an ops package creates a model-side opaque value (an envelope, a form), given Go visibility.
- The six places R21 does not name.
- What `model/relations` is, against `model/contributions`.
- The role of the ops half of a split vocabulary package.
- Whether adding an intended form is an edit of a shared package in the sense of line 7 (contradiction 5 of section 5.3 of the document).
- Any reason for language first.

**Collisions and consequences:**
- **R21 with R7d and R7h:** `ops/internal` cannot exist before five imports are split away, so the splits precede the root move.
- **R24 with AC-11 and R17:** one checker and declared tables, so the intent schema is first.
- **R24 with R25:** two refusals change category; the unknown form gets opposite answers from the four outcomes and from the Files rule.
- **R23 with rule 1:** the 30 identity packages can only be under `model` after their types have a model-side owner.
- **R28 with R29:** the first move touches exactly the 10 packages the migration edits; PR #304 touches two of them as well.
- **R27 with R21:** the closed foundation table spans both roots.
- **R29 with AC-10:** whichever order is taken, a case record is due. Either the correction is recorded, or a deferral with reason and trigger.

## 4. Table

| Row | Verdict | One line |
| --- | --- | --- |
| R21 | NEEDS ADAPTATION | Two roots hold. Sort first below `model`, six places unnamed, `internal` stays at the root until five imports are split away |
| R22 | HOLDS | The matrix stays below `ops`; `ops/check/files` does not exist |
| R23 | NEEDS ADAPTATION | Identities are already independent of paths; their types need a model-side owner first |
| R24 | NEEDS ADAPTATION | One agnostic checker over declared form tables; issuer and constructors need a home; two categories change |
| R25 | NEEDS ADAPTATION | Keep the accepted classes; "rejected" and the missing limit contradict them |
| R26 | HOLDS | `nodes.go:51-55` against ten kinds written; two narrowing conditions leave with R24 |
| R27 | NEEDS ADAPTATION | Decide first, yes; the foundation role is cut by the root and the reading half has no role |
| R28 | HOLDS | Own changes, grouping first; 92 of 129 packages are free of the branch |
| R29 | NEEDS ADAPTATION | Main first, by AC-10's own question and by what line 7 can establish; 58 of 98 packages only take a prefix |
| R30 | NEEDS ADAPTATION | Rules 1, 3 and 4 by imports; 5 is vacuous in Go; 2 is review only |
| R31 | HOLDS | A prefix-only move breaks rule 1 or rule 2 for 27 packages |
