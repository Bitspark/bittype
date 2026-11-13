# verify-line: research 0018, rows R1 to R20 against bittype `main` and the rules

Read at bittype `main` 3998db5 (the document's baseline a28febc plus two documentation commits). Nothing was changed, built or run.

Verdicts: 15 HOLDS and 13 NEEDS ADAPTATION across the 28 rows. None is stale, rejected or deferred as a whole. One part of R13 is deferred with a trigger.

---

## 1. Where the advice contradicts a rule or an accepted design, or shows two of them contradicting

**C1. Validation in a vocabulary package: accepted pages against rule 2 under the advice's test.** Two accepted pages give the intended-forms package validity:
- `LAYOUT.md:67`: "native semantic forms and validation".
- `docs/generation-composition.md:17`: "| Native model | Native form identity and validity, without AST syntax |".

A third admits a method with a criterion, `docs/package-boundaries.md:226`: "It reads no file and starts no process, so a vocabulary can hold it."

R2, R7a and R7i take every evaluating function out of a model package. The owner's rule 2 outranks the pages only if evaluation counts as an operation. The owner asked that (messages 9 and 10) and did not decide it. The same page's definition ("The names two parties must agree on and nothing else", `docs/package-boundaries.md:98`) already contradicts these sentences, and the advice sides with the definition.

**C2. The vocabulary role cannot be wholly on the model side.** `docs/package-boundaries.md:98-99` lists "profile and contract identities" as vocabulary. In code those identities are values of contract-package types:
- `semantics/golang/go/contracts.go:15-21`: `ResolutionContract() resolution.Contract`.
- `semantics/golang/go/vocabulary.go:96-99`: `Evidence() semantic.EvidenceRequirement`.

R8 puts contract packages under ops, and rule 1 forbids a model to import them. So the owner's direction, R8 and the accepted role definition cannot all hold for one package. The vocabulary role has to exist on both sides. The advice does not say this.

**C3. R15 part 3 against "one unit".** `docs/package-boundaries.md:144-146`: "Parts that share a private invariant belong together. The planner and the plan types in `contracts/generation/go` are one unit, because only the planner may issue a plan." `Expression[T]` holds a kind that only the planner creates (`contracts/generation/go/planning.go:547`, `:750`; comment at `expression.go:62-63`). A model-side "generic intent expression carrier" would separate that kind from its only issuer.

**C4. R7d against the accepted Source design.** `docs/source.md:379` declares `Source[L]` as "private sealed storage", and its publication names two receipts (`source/go/source.go:38-44`: `Snapshot semantic.Snapshot`, `Analysis semantic.Result`). Under R8 both are operation-owned. "Generic carriers: model" cannot hold for `source/go` as a whole.

**C5. R19's content digest against three sentences.**
- `trees/representation/go/language.go:12-14`: "Profiles are caller-owned compatibility promises, not hashes or automatically negotiated versions."
- `docs/engine.md:45-47`: "It is not a generated fingerprint".
- `docs/contribution-contract.md:1034`: "A digest of an implementation is provenance. It replaces none of the five."

Also bittype `AGENTS.md:5-7` (owner, 4 November): no "checksum/validation/check workflows, or provenance workflows" as implicit prerequisites. A digest can be provenance. It cannot be the identity.

**C6. R17's order against AC-10 and AC-11.** The advice orders syntax schemas, then fact schemas, then intent schemas. The intent side is where a forbidden per-language copy exists today (AG-7, five copies with the horos model on `origin/feat/other-pairs-196`). AC-11: "Per-language copies of that mechanism are forbidden." AC-10: the misplaced part moves first, the missing part is added afterwards. Fact schemas are missing, not misplaced. The advice gives no reason for putting intent last.

**C7. The advice against itself, with AC-05.** The closing principle says "A model must be complete as a declaration without depending on any operation that consumes it." R18 keeps "complex relationships between nodes" and "evaluating lexical constraints" as code. For Go that leaves three things stated only by code:
- the identifier rule (`check/syntax/golang/provider/go/check.go:13`, through `go/token`);
- two rules between nodes (`check.go:51-81`);
- the universe scope (`transform/syntax/golang/symbol/go/symbol.go:35`).

The declared model is then weaker than what the checker admits. AC-05 (owner clarification, 13 November) requires that such a relationship be stated, "such as a conservative upper bound rather than equality". The advice does not state it.

**C8. R16's reason against D2 is wider than the charter sentence.** `CHARTER.md:129-131` restricts "Generic carriers, contracts, engine, printer and native syntax stages". The accepted `docs/representation-families.md:98-101` already says: "Payload schemas are horos types. ... No separate schema language is introduced." The two do not contradict each other. They do show that horos types (not horos kinds) are an existing declaration mechanism for values, which the advice does not consider.

---

## 2. The five things asked beyond a verdict

### 2.1 The test made concrete

**`transform/syntax/golang/model/go/model.go` (112 lines).**

| Lines | Declaration | Side under the advice's test |
| --- | --- | --- |
| 12 | `var forms = g.NewFormIssuer[s.Language]()` | **No answer.** It creates an authority (a process-local pointer). It neither declares nor evaluates |
| 18-26, 30-33 | 7 form names, 2 primitive names | Model |
| 35 | `Form(name)` issues a form for any non-empty name | Ops (construction) |
| 36 | `Name(f)` reads back, refuses a foreign form | Ops (evaluation) |
| 37-51 | Five constructors. They refuse nothing | Ops ("Building intent expressions: Operation code") |
| 52-93 | `CheckExpression` | Ops. Holds shape conditions and three coverage answers (`:63` the primitive set, `:87-88` external, `:89-90` unknown name) |
| 94-112 | `CheckDeclaration` | Ops. `:105-109` is the record-body condition (R14) |

Nine string constants are left on the model side. No `g.Form[s.Language]` can be produced from them without the ops side, so the model side alone states no intent.

**`semantics/golang/go` (855 lines).**
- Model: `vocabulary.go:13-16` (profile strings), `:22-71` (21 key functions), `model.go:32-39` (`FieldInfo`, the payload shape of the struct term).
- Ops: `model.go:15-31` and `:41-317` (`Model`, `Type`, `Field`, `Open`, 33 reader methods; `:164-167` is a coverage answer), all of `environment.go`, `vocabulary.go:18-19` (errors), `:85-88` (`MatchesConfiguration`).
- No answer:
  - `vocabulary.go:74-81`, `Configuration()`: a compiler run's settings, including `runtime.Version()`, which reads the running toolchain.
  - `vocabulary.go:89-99` and `contracts.go:15-40`: identities by R2, typed by contract packages that are ops by R8.

**`syntax/golang/go` (107 lines).** Wholly model: a marker type, two constant-returning methods, 23 node kinds, 15 roles, 1 reference role, 9 attribute names, 4 concrete kinds.

**Go's grammar table, `check/syntax/golang/provider/go/check.go`.**
- Model: `:15-16` (kind lists) and `:17-41` (the table, 23 entries).
- Ops: `:13` (`identifier`), `:45-83` (`Check`).
- The table binds 17 predicate function values: 10 `identifier`, 2 `g.String`, one each of `g.Nonempty`, `g.Bool`, `g.Decimal`, `g.Set(...)` and an anonymous function. So it cannot move to a model package until predicates are named by identity.
- No answer: `:19`, `func(v string) bool { return v == "." || identifier(v) }`, a constraint that exists only as its evaluator; `:38`, `g.Set("both", "send", "receive")`, data wrapped in a closure.
- `:51-81` are two rules between nodes, declared nowhere.

**`trees/ast/go` (130 lines), `trees/cst/go` (30 lines).** Plain exported data, wholly model. `ast.go:8`: "the carrier performs no validation or allocation".

**`semantics/go/model.go` (560 lines).**
- Model: `:19-125`, the plain records and `Data`.
- Ops: `:127-140` (limits), `:146` (`Semantics[L]`, sealed), `:230-405` (`New`), the accessors.
- No answer: `:26-28` (`Key.Valid`), `:507` (`Facts`), `:524` (`DirectParameterBindings`). They are queries on the sealed graph, neither declaration nor reconstruction of domain meaning.

**`trees/representation/go` (490 lines).** The envelopes exist only through `Pack`, which runs the structure check (`ast.go:56-68`). Under R7d's own second sentence the whole package is ops. The descriptor `Language` mixes identity data with a `reflect.Type` and a refusing constructor (`language.go:15-30`).

**`contracts/transform/go` (151 lines).** Wholly ops: request, contract, interfaces, admission. Nothing model-side.

**`contracts/generation/go` (2,135 lines).** `Form`, `Slot`, `HelperKey`, `Expression` and the issuers have unexported fields and identity pointers (`expression.go:11-14`, `:103-106`). The model candidates of R15 cannot be separated from their constructors (see below).

**`internal/rust` (310 lines).**
- Model: the reserved-word list (`names.go:11`) and the 840 character ranges (`cratename.go`).
- Ops: `Identifier`, `Declaration`, `PathSegment`, `CrateName` and the literal decoders.
- Rules such as "a raw identifier is not `_`, `self`, `Self`, `super`, `crate`" (`names.go:16-18`) exist only as code.

**`internal/golang` does not exist.** Go's lexical rules are the building toolchain's `go/token` and `go/types.Universe` (`check.go:13`, `target/go/provider.go:31`, `symbol.go:35`).

**The unexported-field question.** Go does not allow a type with unexported fields and its only constructor to live in different packages. "Split: declarations model, executable validation ops" therefore has two realisations, and the advice names neither:
1. Plain exported data on the model side and a sealed type on the ops side. This is what trees have (`ast.AST[L]` and `representation.AST`) and what the graph already has inside one file (`Data` at `:118`, `Semantics[L]` at `:146`).
2. The sealed type moves whole to the ops side and only names stay. This is the result for `Form[T]`.

In both, the model side is usable without the ops side as data and carries no guarantee. A model-side type also cannot name an admitted type in its own fields.

**Fit with the rules.**
- **Owner rule of 6 November ("Shared mechanisms own invariant guarantees"): no collision.** The mechanisms that own guarantees are operations under R7f. The rule and the test agree that a model promises nothing. A collision would arise only if a validating constructor stayed in a model-side carrier (`semantics.New`, `Files.Capture`, the issuers).
- **AC-05:** guarantees travel with admitted values, so nothing is weakened as long as no representation is opened to make the split. The 13 November clarification adds the obligation of C7.
- **AC-11:** the test fits it (one checker over declared data).
- **AC-06:** fits.

### 2.2 Evidence per carrier (R5, R6)

| Carrier | Unchecked | Check | Admitted or sealed | What the evidence records | What R6 would change |
| --- | --- | --- | --- | --- | --- |
| Syntax tree | `ast.AST[L]`, `cst.CST[L]` | Level 1: `check.Structure` (`check/syntax/shared/ast/go/structure.go:66`). Level 2: the language checker | Level 1: `representation.AST` (`ast.go:23-29`). Level 2: `semantic.Snapshot` (`contracts/semantic/go/operation.go:18-28`), for a set of trees | Envelope: language ID, profile string, Go marker. Snapshot: language descriptor, entry, envelopes, Files context. Not recorded: which checker ran | See the two gaps below the table |
| Semantic graph | `semantics.Data` | `semantics.New`: integrity and bounds only | `Semantics[L]`, `representation.Semantics` | The receipt `semantic.Result` (`request.go:17-25`, `:98-104`): snapshot, operation key, dependencies, contract (profile, obligations, requirements), requested keys, configuration bytes, limits, reports (obligation, completion, verdict), diagnostics. Admission is `AdmitEvidence` (`evidence.go:44-62`) | Almost nothing: the receipt already is scoped evidence. Missing is a check of fact shapes |
| Files | `files.Files`, `files.Context` | `Files.Capture`, `Context.Capture`, `Cover` | No type of its own: `Capture` returns `Files` again (`files/go/files.go:55-85`). Sealed states are the snapshot's private context and `source.Source[L]` | For `Source`: the receipt (`Source.Evidence()`), matched against `Policy.Evidence` | A distinct admitted type for captured `Files` |

The two syntax gaps:
- A parser checks grammar (`parse/syntax/golang/provider/go/parse.go:89`) and returns an envelope that records only level 1 (`contracts/parse/go/adapter.go:52`). `semantic.Capture` then runs the checker again.
- No native language declares its syntax profile. Compositions pass the literal `native-types-v1` (11 occurrences in 9 files, none in production code). horos declares one (`syntax/horos/go/nodes.go:26`).

The fact-shape gap: nothing checks that a `go/struct` term has one operand per payload entry. The reader refuses on read (`semantics/golang/go/model.go:289`).

**Accepted designs already decide most of it.**
- `docs/analysis.md:28`: "An owned, structurally checked AST snapshot; Capture separately requires syntax admission".
- `docs/analysis.md:45-47`: "A syntax profile specifies admitted AST forms. An analysis profile specifies what the analyzer establishes about those forms."
- `docs/abstraction-levels.md` (table of correspondence), on `semantic.Snapshot`: "It is not itself a successful semantic judgment".
- `docs/source.md:321-325`: "Matching addresses or bytes cannot substitute for exact receipt association."
- `docs/semantics.md:38-41`: "The kernel validates graph integrity; it cannot prove a provider's claims".
- `docs/semantics.md:119-121`: "the kernel does not decode them". A fact-shape check therefore cannot be part of graph admission. It would be a separate mechanism with a provider-supplied schema.
- `contracts/transform/go/transformer.go:22-25`: "Contract declares exactly the evidence required by a typed mapping."

**Profile keys.** `go1.25/captured-types-1/analysis-1`, `rustc-1.90.0/captured-types-2/analysis-1` and the Rust helper setting `captured-types-3` are the model identity and version that evidence carries today. The revision practice is decided on #278 and #276 (`docs/native-semantics.md:833-846`, `:874-881`): the helper profile changes when facts are added; the analysis profile changes when a member is removed or changed.

### 2.3 R13 and R14 against the code

**The check, `transform/syntax/golang/model/go/model.go:105-110`:**

```go
if n == FormDefined {
    body, ok := Name(d.Body.Form())
    if !ok || body != FormRecord {
        return g.ErrInvalid
    }
}
```

Four facts:
1. It returns `ErrInvalid`, so the code treats it as a rule of the form, not as coverage. `source/go/source.go:80-86` defines the difference: well formed under the vocabulary and outside what is covered is unsupported.
2. Go's own declared syntax admits any type under a defined type: `check.go:20`, `g.One(s.RoleType, types...)`.
3. The Go target does not need the restriction: `transform/syntax/golang/target/go/provider.go:96-106` lowers a `defined` declaration with any body to `type-declaration`.
4. The condition is the declaration recipe's policy (`transform/syntax/horos/mapping/declarations/go/declarations.go:53-55`: record form when the source body is a record, otherwise alias).

**The page that speaks of it.** `docs/contribution-contract.md:2502`: "Each target's model requires that the selected body of that form is its record form". That section says of itself "Nothing in this section is the design". The design part of the same page gives the choice to the contribution, `:203`: "alias or defined declaration | The contribution that supplies the declaration: the default declaration policy or an exact mapping". No accepted page states the restriction as a rule of Go.

**Result.** It is not a rule of Go and not a rule the target relies on. It is a mapping policy enforced in the model package under the wrong category. For TypeScript's `interface` and Rust's `struct` the same condition is a rule of the form. Only Go's is not.

**The Go rows of the expert's table against section 3.3 and the code.**

| Row | Finding |
| --- | --- |
| Slice | Matches: `check.go:35`, `vocabulary.go:70`, `model.go:20` |
| Defined type | Matches, with R14 as above. The analysed `defined` entity's body is a reference to itself today (`internal/golangsem/export.go:207`) |
| Alias | Matches in name. The analysed kinds are wider than the declaration forms: `vocabulary.go:40-42`, "error and comparable are named types of the universe scope and any is an alias, so their entities have DefinedKind and AliasKind" |
| Field | Matches. `field` is an unpublished literal of the bridge (`export.go:135`) |
| Primitive, builtin, basic | Matches |
| Record, struct | Matches. The Go target reads only the field name and ignores the optional flag (`provider.go:122-123`) |
| External | Matches (`model.go:87-88`) |
| Type parameter | **Does not match.** Go's intent has no notion of its own to share. The binder is the kernel's (`ParameterExpression`, `expression.go:55`; `NativeParameter`, `native.go:21`), by `docs/generation-planner.md:434-436`: "The kernel owns capture identity, substitution and scope." |
| Generic application, package | Match |
| Interface, union, signature, tuple, boolean | Written by the bridge as unpublished literals (`export.go:253`, `:259`, `:267`, `:271`, `:277`). No reader names them |

### 2.4 R15 to R17 against AG-7

**Does it give AG-7 one agnostic owner?** It permits one and does not state it. R15 says "the operation that constructs and validates intent". A checker per language would satisfy those words and still break AC-11. The requirement has to be added: one checker in the agnostic layer, with each language supplying a declared table.

The pairs lane says the same: "(4) Yes, as a table of forms for one checker."

**The declared part for two Go forms.** This is a sketch; no such table type exists.

- **`field`.** Data: `name`, a non-empty string; `optional`, `true` or `false`. Arguments: exactly one, of any form. Today the constructor at `model.go:43-45` and the check at `:70-73` state this twice.
- **`record`.** Data: none. Arguments: any number, each of the form `field`, with the first datum unique among them. Today `model.go:46-48` and `:74-86`.

The table type needs: data arity with a per-datum constraint (non-empty, one of a set); argument count; a constraint on the form of arguments; a sibling-unique datum; declaration forms with an optional body-form constraint. The record rule is identical in all five copies, which is the concrete need that R18's "unless" clause asks for.

**What could not be declared.**
- That a field name is a Go identifier. The model does not check it today; the Go checker refuses it at sealing (`check.go:27`).
- What `optional` means for Go.
- What each form becomes as nodes. That is the target's lowering.
- Coverage: Go's refusal of `external`, and which primitives exist.
- The body rule of Go's `defined` (R14).
- Where a form may occur. A `slice` of a `field` passes the model's check and fails only at sealing.

**One real fact written and read.**

Writer, `internal/golangsem/export.go:235-243`: for a struct, one `n.FieldInfo{Name, Tag, Embedded, Entity}` and one operand per field, then `e.expression(n.StructOperator(), fields, operands...)`. The helper at `:109-116` JSON-encodes the data under the literal `n.Key("json-1")` and gives operand `i` the locator `[i]`. `:286` adds the fact `underlying`.

Reader, `semantics/golang/go/model.go:283-302`:

```go
v := t.Underlying().term()
if v.Operator != StructOperator() { return nil, false }
... json.Unmarshal(v.Data, &metadata) != nil || len(metadata) != len(v.Operands)
```

- **Declared today:** the Go struct `FieldInfo` and the operator name, shared by import.
- **Not declared:** the encoding key (the reader never compares `v.Encoding`); the locator convention (the reader goes by position); that the entity has the kind `field`; that `underlying` has exactly one value (`model.go:191`).
- **A fact schema would hold:** the relation `go/underlying` (subject a term, one value); the operator `go/struct` (encoding `go/json-1`; data a list of name, tag, embedded and an entity of kind `go/field`; one operand per entry at locator `[i]`).
- **It could not hold:** that the entity is the right field.
- **Second example:** the payload of `basic` is the number of `go/types.BasicKind` (`export.go:198`, `model.go:210-214`). That is a Go standard library value which a reader in another language cannot interpret.

**A further copy the pairs lane did not look at.** The helper "JSON data under `json-1`, operands by position" is written three times: `export.go:109-116`, `internal/typescriptsem/compiler.cjs:151`, `internal/rustsem/helper/src/main.rs:361`. The Go bridge writes 11 keys as literals that the vocabulary does not publish (command at the end).

### 2.5 R19: identities today

| What | Identity | Depends on a Go path or marker? |
| --- | --- | --- |
| Semantic vocabulary | `semantics.Key{Owner: "go", Name: name}` (`vocabulary.go:22`) | No |
| Syntax vocabulary | `type Kind = ast.Kind[Language]` (`syntax/golang/go/nodes.go:31`); the ID `"go"` (`:16`); the descriptor's `Compatible` compares `l.id == other.id && l.profile == other.profile && l.marker == other.marker` (`language.go:41-43`) | The marker, in process only. It is a `reflect.Type` and is never serialised. `docs/engine.md:48`: "Matching requires all three of language ID, profile and Go marker." |
| Profile | Strings (`vocabulary.go:13-14`; `semantics/rust/go/vocabulary.go:18-19`, `:34`) | No. For native syntax no profile is declared at all |
| Form | `Form[T]{owner *identity; name string}` (`expression.go:11-14`); "identity is an allocation-unique owner" (`:103`) | The marker `T` and a process-local pointer. No declared identity or revision exists. The plan on the migration branch does not serialise forms |
| Contribution, on `main` | The mapping's `Profile` string, `"structural-representation/1"` (`transform/syntax/horos/mapping/go/mapping.go:28`) | No |
| Contribution, on `origin/feat/contribution-contract-196` | `ContributionReference` (publisher, key, meaning with revision, representation), `SchemaReference`, `EvidenceKind` (`contracts/generation/go/contribution.go:22-46`, `:138`), for example `"bittype/golang/record-form"`, version `"1"` | No. The strings are spelled after the directory name `golang`, while the language ID is `go` (pairs report F3) |

---

## 3. Rows

**R1. The test: declares, or determines, interprets, produces.** NEEDS ADAPTATION.
- Evidence: 2.1. The test decides most declarations and gives no answer for five kinds of code: issuers; constructors that refuse nothing (`expression.go:76-78`); constraints that exist only as their evaluator; values typed by a contract package that a vocabulary publishes; plain queries on a sealed value.
- Adjustment, two sentences added:
  1. A type whose non-zero values only a refusing function creates is the admitted state and lives with that function on the ops side. Model-side types are plain data (`docs/package-boundaries.md:144-146` supports this; it is R6).
  2. A constraint without a separate declaration is named by the model and stated by its evaluator. The declared model is then an upper bound of what is admitted, and the page says so (AC-05, 13 November).
- No owner rule is contradicted.

**R2. What a model package may and must not contain.** NEEDS ADAPTATION.
- The lists hold for the syntax vocabularies, the tree packages and the plain records of the graph.
- Adjustments:
  1. "Identity and version information" excludes values typed by a contract package (C2). They go to an ops-side vocabulary: `semantics/golang/go/contracts.go`, `vocabulary.go:74-99`.
  2. "Builders" must be read with R1's first sentence, or `Expression[T]` cannot exist as a model type.

**R3. A common value type is owned independently of either operation.** HOLDS.
- `language.ID` and `language.Provider` (`contracts/language/go`, 25 lines) are imported by every syntax vocabulary and become model-side shared values.
- The semantic handles already have an owner that is neither resolution nor analysis (`contracts/semantic/go`).
- Limit, from R8: a receipt of one operation that another consumes stays with its producer. `generation.Output` is used by `contracts/format/go`; an ops-to-ops import is allowed.

**R4. An operation consumes the model as explicit input and does not own the declaration it interprets.** HOLDS.
- The generic checker already takes its table as an argument (`internal/grammar/rules.go:80`).
- Operations own declarations today in three places: the five grammar tables (unexported variables of checkers, or `internal/declsyntax`); the 11 literal keys of the Go bridge; the names spelled in the TypeScript and Rust bridges.
- The semantic owner stays the language (`docs/semantics.md:38-40`, `docs/generation-planner.md:529-531`). Only the package changes.

**R5. Four levels; level 3 is external correctness; no guarantee inferred from another.** HOLDS.
- The document's own definition of correctness names "a compiler's judgment", so its placing of level 3 under consistency contradicted itself.
- In code, level 3 is decided by compiler bridges and carried as receipts.
- The non-inference is accepted in several places: `trees/representation/go/ast.go:71`; `contracts/semantic/go/request.go:107`; `docs/semantics.md:206`; `docs/parsing-formatting.md:121`; `docs/generation.md:94`; AC-07.
- Not established: a page that states the levels. None exists.

**R6. Checked state separate for every sort of model; generalise the envelope.** NEEDS ADAPTATION.
- The pattern exists for all three carriers, and accepted pages decide it (2.2).
- Adjustment: introduce no new generic evidence or admitted-artifact concept. The envelope, the snapshot, the receipt and `Source` are that concept (derivation rule).
- Four gaps to close:
  1. Native syntax vocabularies declare their profile.
  2. A parsed or formatted envelope does not record that grammar admission ran.
  3. No fact-shape check exists. It must be separate from graph admission (`docs/semantics.md:119-121`).
  4. `Files.Capture` returns the unchecked type.

**R7a. Constructors and checks: split.** NEEDS ADAPTATION.
- The premise "constructors that enforce invariants are executing validation" covers three different things in the code:
  - constructors that refuse nothing (`model.go:37-51`);
  - refusing constructors of sealed types (`expression.go:18-23`, `language.go:21-30`, `symbol.go:34-43`);
  - admission of plain data (`semantics.New`, `Files.Capture`).
- Adjustment: R1's first sentence. Do not split a sealed type from its constructor; move both.

**R7b. Typed readers: ops; fact schemas: model.** HOLDS.
- Paths in 2.1. It resolves the known contradiction of `LAYOUT.md` ("immutable native type readers" in a vocabulary) in favour of the definition.
- Not answered: which role an ops-side reader package has. It is imported by a query service and by verification, so it is neither a contribution nor a private part.

**R7c. Structure checks and grammar tables: split.** HOLDS.
- `docs/abstraction-levels.md:22-24` has this shape: "A grammar, for example, can supply declarative rules."
- Here: the table types of `internal/grammar/rules.go:15-40` and the five tables go to the model side; `Check` (`:80-146`), `check.Structure` and the rules in code stay ops.
- Precondition: predicates named by identity (17 function values in Go's table).

**R7d. Generic carriers: model.** NEEDS ADAPTATION.
- True for `trees/ast/go`, `trees/cst/go` and the plain records of `semantics/go` and `files/go`.
- Not true for `trees/representation/go` (admitted state and packing) and for `source/go` (C4): `Policy` holds two operation callbacks (`source.go:124-130`).
- Adjustment: of six carrier packages, two are wholly model, two split, two are ops except `Unit`, `Group` and `Relation`.

**R7e. Contract packages: split by responsibility.** HOLDS.
- By reading the declarations, a model-side part exists in two of the thirteen packages: `language`, and `generation` (the intent data, see R15).
- `text.Diagnostic` is an operation's result and stays ops.
- The document's test found ten packages holding both. The difference is the advice's test for a value (R8).

**R7f. Engine, registry, planner: ops.** HOLDS. It fits the owner rule of 6 November and the term "kernel" in `docs/package-boundaries.md:93-94`.

**R7g. Mapping identities and decision schemas: model declarations.** HOLDS.
- On `main` this is two profile strings (`transform/syntax/horos/mapping/go/mapping.go:28`, `:34`).
- The migration branch has 24 vocabulary packages and the pairs branch 33, each with a reference, a schema reference, a horos schema and generated types. The 33 codec packages are ops.
- Precondition: the identity types (`contribution.go:22-46`, `:138`, plain string structs) leave the generation contract package for a model-side owner.

**R7h. Lexical rules: split.** HOLDS.
- All six call sites of lexical evaluators in vocabulary packages are inside validating code, which is ops under R7a and R7b: `semantics/rust/go/environment.go:156`, `:160`; `semantics/typescript/go/environment.go:37`; the Rust symbol package `:56`, `:58`, `:73`; the TypeScript symbol package `:24`.
- So the "needed on both sides" problem of case (h) disappears.
- For Go the declaration is the toolchain's standard library, which no profile records today.

**R7i. Intended forms: split.** NEEDS ADAPTATION.
- Adjustments:
  1. The page wording of C1 changes.
  2. The ops-side checker is one agnostic mechanism (AC-11).
  3. The accepted SPI already gives the checks to the target (`docs/generation-spi.md:354-360`), so the target passes its language's table to the shared checker and keeps coverage. Rust's does so already (`transform/syntax/rust/model/go/model.go:57-60`).
  4. The constructors are used by six Go mapping files. They become ops-side helpers over one shared constructor.

**R8. No third root; a receipt belongs to its operation's contract.** HOLDS.
- The owner's message 8 names two roots. `docs/package-boundaries.md:144-146` agrees on receipts.
- Consequence: of the 24 imports from vocabulary packages into contract packages, the 5 into `language` become model to model. The other 19 must end, through R7b, R15 and C2.

**R9. Several related models per language.** HOLDS.
- `docs/generation.md:77-79` and `docs/kernel-composition.md:22-23` separate intent from analysed facts.
- Note: `docs/abstraction-levels.md` has five levels and intent is none of them. The row adds a sort of model without a level.
- C++ has syntax and intent and no semantics. `semantics/` holds `golang`, `horos`, `rust`, `typescript`.

**R10. U1 first, with declared and testable relations.** NEEDS ADAPTATION.
- U1 follows the platform table: "Independent identities or lifetimes | Keep separate concepts and define their relationship."
- Adjustment: "declared" cannot mean a table for the relation from a form to nodes. R18 keeps it code, and the target is the "trusted interpreter" of its forms. It means a stated relation on the language's page, plus a test per form that builds, checks, analyses and compares.
- The relation from a form to analysed meaning is exercised today only by verification laws (`verification/files/horos-to-go/mappings/{list,record,optional}/go/law.go`).
- A relation as data waits for a reader (derivation rule).

**R11. U2 selectively, by criteria.** HOLDS.
- The criteria are the platform's (`docs/concepts/semantic-ownership-and-generalization.md:250-253`).
- The advice's list omits "allowed operations" and the substitution condition. The platform's wording is the binding one.

**R12. Reject U3 as the starting point.** HOLDS.
- Syntax has notions without a type-level counterpart (`file`, `import`, `parameter`).
- One Rust form becomes a path of several nodes or two different shapes (`transform/syntax/rust/target/go/provider.go:161-167`).
- C++ has no analysis.
- Fixed points 15 to 17 and the derivation rule agree.

**R13. The Go rows.** NEEDS ADAPTATION.
- Nine of eleven rows match the code (2.3).
- Adjustments:
  1. Type parameter: nothing Go-owned on the intent side to share.
  2. The analysed `alias` and `defined` kinds are wider than the declaration forms.
  3. The four "shared concept identity" rows have no reader. Only `integration/go/destinations_test.go` imports both vocabularies, on `main` and on the pairs branch. That part is deferred.
- Trigger: a non-test file imports both `transform/syntax/golang/model/go` and `semantics/golang/go`. Command: `git grep -l 'transform/syntax/golang/model/go"' -- '*.go' | xargs grep -l 'bittype/semantics/golang/go"'`.

**R14. The record-body restriction is not Go's; the target owns it.** NEEDS ADAPTATION.
- The first half holds (2.3).
- Adjustment to the second half: the target does not need the restriction. Its owner is the declaration contribution (`docs/contribution-contract.md:203`). The condition leaves the model package. It does not become a target rule.
- Not established: why the check was written (it came with c054f80, #155).

**R15. The intended-forms package becomes five things.** NEEDS ADAPTATION.
- Adjustments:
  1. Part 4 is one agnostic checker and one agnostic constructor (AC-11, AG-7).
  2. Part 3 must say which part of `Expression[T]` is the carrier. Only the native kind (form, data, arguments) is intent. The five other kinds are the planner's protocol, and one is planner-issued (C3).
  3. The form's identity today is an issuer pointer, so parts 1 and 2 need a declared identity (R19).
- The horos model on the pairs branch confirms "horos as a target declares its own intent vocabulary". It is the fifth copy.

**R16. D1 first, not D2, not D3.** HOLDS.
- D1 extends a mechanism that exists.
- D3 is ruled out by `docs/representation-families.md:100-101` and the derivation rule.
- Note C8: for value-shaped schemas, horos types with a generator and a regeneration test already exist on the migration branch. The advice does not weigh them.

**R17. Three schema families; syntax, then facts, then intent.** NEEDS ADAPTATION.
- The syntax family already exists almost entirely in `internal/grammar`: required attributes, roles, cardinalities, allowed kinds, order (`rules.go:117`), one reference role.
- Adjustments:
  1. The order follows AC-10 (C6): the intent table and the placement of the grammar tables are corrections; fact schemas are additions.
  2. "Parameter and reference structures" belong to the shared carrier, not to a language's intent schema.
  3. "Independently versioned" uses profiles.
- The concrete cases the derivation rule asks for exist: AG-7 for forms; hand-synchronised writers and readers for facts (AC-05, 13 November).

**R18. What stays code.** HOLDS.
- It matches the code item by item. For Go: two rules between nodes, the `identifier` predicate, a bridge of 1,241 lines, 33 reader methods, the target's switch.
- Its "unless" clause is already met for the record rule (five identical copies).
- With C7: the page must say the declared grammar is an upper bound.

**R19. Explicit identity and revision, independent of path and marker.** NEEDS ADAPTATION.
- Table in 2.5.
- Adjustments:
  1. Use the existing concept, the profile, and the decided practice of #278. The five version axes are accepted for contributions (`docs/representation-families.md:102-104`).
  2. A digest is provenance only (C5).
  3. What is missing: a declared profile for each native syntax vocabulary, and an identity and revision for each form vocabulary.
  4. Of the four revision events, two are distinguished today (removal or change; addition). The other two need a case. One exists for interpretation: Rust's `Lookup` changed with an added fact (`docs/native-semantics.md:848-860`).
- Note: lexical rules that follow the Go release (`unicode.IsLetter` in `rules.go:66`, `go/token`) are in no profile.

**R20. A conformance suite.** NEEDS ADAPTATION.
- Items 1 to 3 exist. The checker interprets the table, and the ten parser and formatter providers call their checker.
- Missing: fact-shape validation, and reader tests against schemas.
- Item 6 ("targets produce trees consistent with intent schemas") has no checkable meaning until R10's relation is stated. Today the output is checked against the syntax table at sealing.
- For the bridges, nothing is generated on `main`. Rust has a run-time handshake on one name (`semantics/rust/go/vocabulary.go:29-33`, `main.rs:93`).
- AC-05 (13 November) requires the mechanism to be general, "without requiring a new case-specific check for each implementation change". The generator from one description meets that. The payload generator on the branch is the precedent.

---

## 4. What the advice did not answer, and what two rows give together

**Not answered.**
- Whether the code remainder (rules between nodes, predicates) is still model. Question 3 asked it.
- Which realisation of "split" applies to a sealed type (2.1).
- The role of ops-side packages that are neither contribution nor private part: readers, contract identities, intent constructors.
- Whether the form checker is one mechanism. Section 6.3 left it open and the advice passes over it.
- What a shared concept identity is concretely, and who reads it.
- Whether fact payloads and form data are declared as Go tables or as horos types.
- How the declared model relates to a toolchain-owned rule (Go's keywords and universe scope).

**Rows together.**
- **R1 with R6:** read together they are consistent only as "plain data on the model side, admitted state on the ops side". Read literally, R7a alone would open sealed types.
- **R2 with R8:** the vocabulary role splits across the roots (C2).
- **R3 with R8:** R3 does not reach receipts.
- **R7d with R8:** `source/go` and the envelope package are ops.
- **R7a and R7b with R7h:** the lexical helpers are needed on one side only.
- **R15 with R16 and R17, plus AC-11:** one table type and one checker for forms. The syntax side already has that. AG-7 gets its owner only if the single checker is made a requirement.
- **R14 with R15:** Go's intent schema declares no body rule for `defined`, while TypeScript's, Rust's and C++'s do. So "a declaration form's body" is language-owned data, as the pairs lane wrote.
- **R17 with R20, on `docs/package-boundaries.md` finding 9:** moving the table alone does not remove the parser and formatter exception, because full admission also needs the rules in code.
- **Outside my rows, noted for their owners:**
  - The layout sketch in advice 4.1 gives C++ `syntax/` and `semantics/`. C++ has syntax and intent and no semantics.
  - The four refusal categories of R25 already exist for Files (`source/go/source.go:56-87`) and as completion against verdict in `contracts/semantic/go/contract.go:47-67`.

---

## 5. Table

| Row | Verdict | One line |
| --- | --- | --- |
| R1 | NEEDS ADAPTATION | The test leaves five kinds of code undecided; add the sealed-type rule and the upper-bound rule |
| R2 | NEEDS ADAPTATION | Contract identities and evidence requirements cannot be model-side; "builders" needs R1's rule |
| R3 | HOLDS | `language.ID` becomes shared; receipts stay with their operation |
| R4 | HOLDS | Operations own declarations today in the grammar tables and the bridges' literal keys |
| R5 | HOLDS | Level 3 as external correctness matches the receipts; non-inference is already accepted |
| R6 | NEEDS ADAPTATION | The pattern exists and is decided; no new concept, four gaps |
| R7a | NEEDS ADAPTATION | Move a sealed type with its constructor; do not split them |
| R7b | HOLDS | Readers to ops; `FieldInfo` and the keys stay |
| R7c | HOLDS | Tables and table types to model, once predicates are named by identity |
| R7d | NEEDS ADAPTATION | Two carrier packages are model, two split, two are ops |
| R7e | HOLDS | A model-side part exists in two of thirteen contract packages |
| R7f | HOLDS | Fits the owner rule of 6 November |
| R7g | HOLDS | Pair vocabularies are model, once their identity types leave the contract package |
| R7h | HOLDS | All callers of lexical evaluators end up on the ops side |
| R7i | NEEDS ADAPTATION | One agnostic checker; coverage to the target; page wording changes |
| R8 | HOLDS | 19 of 24 vocabulary-to-contract imports must end |
| R9 | HOLDS | Intent is a sort of model without one of the five levels |
| R10 | NEEDS ADAPTATION | A stated relation and a test, not a table |
| R11 | HOLDS | Use the platform's six criteria and its substitution condition |
| R12 | HOLDS | Syntax-only notions and one-to-many lowering are in the code |
| R13 | NEEDS ADAPTATION | Type parameter has nothing to share; shared identities deferred until a reader exists |
| R14 | NEEDS ADAPTATION | Not Go's rule, but the declaration contribution's, not the target's |
| R15 | NEEDS ADAPTATION | One checker and one constructor; split `Expression` kinds; forms need a declared identity |
| R16 | HOLDS | D1 extends what exists; D3 is ruled out by an accepted sentence |
| R17 | NEEDS ADAPTATION | Order by AC-10: intent and table placement before fact schemas |
| R18 | HOLDS | Matches the code; the record rule already meets its exception |
| R19 | NEEDS ADAPTATION | Profiles are the identity; a digest is provenance; native syntax and forms lack one |
| R20 | NEEDS ADAPTATION | Three items exist; item 6 is undefined; the mechanism must be general |

**Commands behind the counts** (run at the repository root):
- Literal keys of the Go bridge: `grep -oh 'n\.Key("[a-z0-9-]*")' internal/golangsem/*.go | sort | uniq -c` (11).
- Published key functions: `grep -c 'semantics.Key *{ *return Key(' semantics/golang/go/vocabulary.go` (21).
- Predicates in Go's table: `grep -o 'g\.Attr([a-z]*, [^)]*' check/syntax/golang/provider/go/check.go` (17).
- Reader methods: `grep -c '^func (m Model) \|^func (t Type) ' semantics/golang/go/model.go` (33).
- Vocabulary imports: a loop over `syntax/*/go semantics/*/go transform/syntax/*/{model,symbol}/go` listing module imports of non-test files (24 into contract packages).
- Vocabulary and codec packages on the branches: `git ls-tree -r --name-only <ref> -- transform | grep "/vocabulary/.*/go/vocabulary.go" | wc -l` (24 and 33), and the same for `codec.go` (33).
- The profile literal: `git grep -c 'native-types-v1' -- '*.go'` (11 in 9 files).
- Providers that import their checker: `git grep -l 'check/syntax/[a-z]*/provider/go"' -- 'parse/**/*.go' 'format/**/*.go'` (10).
- Line counts: `wc -l` over the non-test files of each package.
- Not counted by a command: the 23 table entries (read from `check.go:18-40`) and the classification of contract packages in R7e (by reading).
