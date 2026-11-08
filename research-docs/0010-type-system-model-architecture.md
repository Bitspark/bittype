# Research: Type-system model architecture and modular reasoning

> Historical consultation and dated evaluation. See the [research map](README.md)
> for current contracts and subsequent decisions; original advice/submissions
> retain their inspected baseline.

**ID:** 0010
**Date:** 5 November 2026
**Status:** applied
**Run-ID:** run_e274981d-8f47-46db-b52c-1f0b931b567d
**Owner:** codex-hydrated-type-design
**Issue:** https://github.com/Bitspark/bittype/issues/92
**Reviewed:** https://github.com/Bitspark/bittype/issues/92#issuecomment-6669329915
**Graded:** https://github.com/Bitspark/bittype/issues/92#issuecomment-6669978474

## Applied

Evaluation scope is **design documentation only**. The complete frozen
[submitted brief](0010-type-system-model-architecture.submitted.md) and
[raw advice](0010-type-system-model-architecture.advice.md) were read. Attempt 1,
document `doc_d799c49f-8580-4342-a6f1-cf7005119237`, run
`run_e274981d-8f47-46db-b52c-1f0b931b567d`, answered by the service's labelled
6 Pro advisor, describes Bittype `93713a0`, Horos `ad7c6f4` and Logos `d0df8bbb`.
Those code revisions were used for the independent evaluations.

The original capture failed to observe completion controls. The supported capture
recovery published the existing canonical response as
`result_73d7f50e-29c2-4235-b65b-7c9a9a7e385f`; no second consultation was submitted.
Attach exited 0, and the service's run/document/process association and downloaded
headers agree. Content and question coverage were assessed manually; the recovered
status supplied no automated advice verdict. Under the owner's no-extra-checksum
or administrative-workflow rule, no separate digest or verifier workflow was
run, and none is claimed. The recovery record retains that distinction.

Independent read-only evaluations: `model10_bittype` checked carriers, receipts,
transform outputs and layout; `model10_horos` checked constructor guarantees,
inspection, original constraints and recursion; `model10_logos` checked the
production SPI and proposed typed foundation. No compiler, software test or proof
execution was performed. The Go, TypeScript, Rust and C++ primary references in
the architecture document were checked for the specific language distinctions.

The main adaptation is that Horos already supports `Satisfies(Page, K => K)` and
public declaration inspection. An extra DataPage declaration or new inspection
API is not a prerequisite. The missing work is the precisely scoped evidence
bridge. The other adaptations preserve the existing root Go module, avoid
premature transport decisions, and use functional input/assumption association
without adding administrative manifests, checksums or provenance workflows.

Every HOLDS/NEEDS ADAPTATION row below is a documentation action, owned by
`codex-hydrated-type-design` (abbreviated **owner**). Links identify the concrete
design sections to deliver. Validation means checking that the documented
contract covers the stated case and agrees with the cited code; it does not mean
that a proposed implementation test has run.

All accepted documentation actions are now implemented in the linked design and
the reconciled model, preservation, generation, semantics and layout documents.
Final consistency review found no actionable contradictions. Local file/heading
links and Git whitespace checks passed. No production source, dependency or
provider changed; no compiler, software test or proof execution is claimed.
The three deferred rows retain explicit implementation/certificate triggers.

| # | Recommendation | Verdict | Action | Owner | Link | Validation |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | Consume Logos semantic packages; no Bittype-owned type-system foundation | HOLDS | Replace the separate shared-modeling-library direction with explicit consumer integration. | owner | [Ownership](../docs/type-system-architecture.md#ownership-and-dependencies) | Matches Logos ADR0038/design0014; generic Bittype graph retains no rules. |
| 2 | Separate object semantics, foundation soundness, evidence validity and adequacy | HOLDS | Define the four assurances and their limits. | owner | [Assurance](../docs/type-system-architecture.md#four-assurance-boundaries) | A valid receipt/proof is not presented as compiler or extraction correctness. |
| 3 | Relate source, native and proof binders | HOLDS | Require real three-way correspondence, not same spelling or a universal Binder record. | owner | [Binding](../docs/type-system-architecture.md#binding-and-the-expected-statement) | Existing ParameterID/EntityID and generated binder references are retained. |
| 4 | Distinguish proof-context extension from native environment extension | HOLDS | Bind native environment and extension policy as semantic inputs. | owner | [Binding](../docs/type-system-architecture.md#binding-and-the-expected-statement) | C++ specialization and TypeScript merging cannot be treated as harmless logical weakening. |
| 5 | Add only the missing generation/claim association | HOLDS | Keep the first seam private and specific to actual captured source, output and selected mappings. | owner | [Composition](../docs/type-system-architecture.md#explicit-composition) | Does not replace AST, Semantics or Horos authority. |
| 6 | Put optional integration in a separate module | NEEDS ADAPTATION | Preserve the one-root-module rule; describe opt-in stage packages first. A future external module requires a concrete dependency decision. | owner | [Layout](../docs/type-system-architecture.md#proposed-package-roles) | Optional Go packages are not claimed to isolate module dependencies. |
| 7 | Keep ASTTransformer and wrap it explicitly | HOLDS | Capture actual output, explicitly obtain needed native evidence and check the separately constructed statement. | owner | [Composition](../docs/type-system-architecture.md#explicit-composition) | Source receipts do not certify target output or survive another transform leg implicitly. |
| 8 | Use a small Go-facing checking seam | NEEDS ADAPTATION | Document a private shape using bound Logos statement/evidence references; no copied foundation schema or public SPI yet. | owner | [Composition](../docs/type-system-architecture.md#explicit-composition) | Expected statement is not chosen by the candidate proof. |
| 9 | Bind actual helpers and emitted artifacts | HOLDS | Scope AST claims to captured ASTs; file claims require explicit emission correspondence. | owner | [Binding](../docs/type-system-architecture.md#binding-and-the-expected-statement) | Changed helper/body or printer output cannot reuse an unrelated result. |
| 10 | Keep invalid evidence, refutation and incomplete execution separate | HOLDS | Reuse the existing completion/verdict distinction with separate named obligations. | owner | [Outcomes](../docs/type-system-architecture.md#outcomes-and-trust) | Invalid proof does not refute its proposition; Unproved/exhaustion is not false. |
| 11 | Mapping ownership does not establish law validity | HOLDS | Every law has a precise proposition and checked support or an explicitly accepted assumption. | owner | [Representation](../docs/type-system-architecture.md#representation-and-selected-laws) | No self-asserted final preservation theorem or builtin privilege. |
| 12 | A narrow Logos typed checker precedes symbolic certification | NEEDS ADAPTATION | State the recommended prerequisite and required operations; Logos owns its admission. | owner | [Checker](../docs/type-system-architecture.md#the-logos-checking-prerequisite) | Test-private Rust/proposed contracts are not claimed as production APIs. |
| 13 | Use a Go client to a checker process on shared workers | NEEDS ADAPTATION | Keep a process adapter as a candidate, pending the admitted interface. Separate workload-location rules from production deployment design. | owner | [Checker](../docs/type-system-architecture.md#the-logos-checking-prerequisite) | No mandatory service, new protocol implementation or deployment workflow inferred. |
| 14 | Require contexts, abstraction, substitution, application and exact statement checking | HOLDS | Make these explicit capabilities and refusal conditions. | owner | [Checker](../docs/type-system-architecture.md#the-logos-checking-prerequisite) | Ground DAG and RuleAdmissibility alone cannot establish generalization. |
| 15 | State the exact quantified Page theorem | HOLDS | Fix source/native domains, resolved kind, family identities and static Rep relation. | owner | [Page](../docs/type-system-architecture.md#statement-two-symbolic-page) | Quantifies over arbitrary admitted T/N, not sampled types or all kinds K. |
| 16 | Obtain a source schema from DataPage or add inspection/export | NEEDS ADAPTATION | Prefer the existing constructor-arrow query; existing admitted alias inspection is another route. Only the evidence bridge is missing. | owner | [Horos source guarantee](../docs/type-system-architecture.md#the-authoritative-source-guarantee) | `knownConstructor` retains original bounds; no mandatory helper, duplicate Page or new inspection API. |
| 17 | Prove actual named-family/field correspondence | HOLDS | Connect real definitions, binders, fields and helpers without equating a native named type to its anonymous body. | owner | [Representation](../docs/type-system-architecture.md#representation-and-selected-laws) | Same names or unfolded shapes alone are insufficient. |
| 18 | Use checked abstraction instead of fresh-atom generalization | HOLDS | State the proof-context rules and uniform imported-law requirements. | owner | [Page proof](../docs/type-system-architecture.md#checked-abstraction-and-nested-use) | No global parametricity assertion or producer freshness Boolean substitutes for checking. |
| 19 | Instantiate the theorem for nested Page and preserve original contracts | HOLDS | Feed inner membership and representation results into outer application. | owner | [Nested use](../docs/type-system-architecture.md#checked-abstraction-and-nested-use) | Includes ignored-parameter and actual DataPage alias obligations. |
| 20 | Native admission is independent of conditional preservation | HOLDS | Retain the independent implication from accepted native arguments to source/Rep premises. | owner | [Admission](../docs/type-system-architecture.md#native-admission-remains-independent) | Page[N any] is not advertised as enforcing Data. |
| 21 | Keep recursive-kind reasoning authoritative initially | HOLDS | Import a scoped trusted source guarantee; the composition proof adds no new coinduction engine. | owner | [Recursion](../docs/type-system-architecture.md#recursion-and-later-certificates) | Greatest-fixed-point meaning and Unproved inclusion remain intact. |
| 22 | Later independently replay recursive membership/inclusion certificates | DEFERRED | Trigger: a Horos-owned certificate contract specifies actual fixed-point/domain/guard obligations and a qualified checker; inclusion requires separate coverage. | Horos integration | [Deferral](../docs/type-system-architecture.md#recursion-and-later-certificates) | No proof-by-revisitation or implemented certificate is claimed. |
| 23 | Consider a consumer certificate calculus on today's leaf SPI | DEFERRED | Trigger: an explicit delivery requirement cannot use the Logos milestone and a separate sound certificate contract is justified. Never an automatic fallback. | integration owner | [Alternative](../docs/type-system-architecture.md#the-logos-checking-prerequisite) | No parallel binder/checker framework is selected now. |
| 24 | Preserve completion outside binary leaf hooks | NEEDS ADAPTATION | Distinguish binary Leaf/Rule hooks from completion-aware GroundDecider; any accepted receipt comes from the actual checker and exact immutable input. | owner | [Outcomes](../docs/type-system-architecture.md#outcomes-and-trust) | Producer-supplied verified flags do not confer authority. |
| 25 | Start from two exact statements, then derive the minimum API | HOLDS | Define TextBox and Page before public contracts or a general solver. | owner | [Sequence](../docs/type-system-architecture.md#implementation-and-qualification-sequence) | Both exercises precede contract promotion. |
| 26 | Establish a checked TextBox bool counterexample | HOLDS | Require actual native acceptance, fixed interpretation and authoritative source refutation. | owner | [TextBox](../docs/type-system-architecture.md#statement-one-textbox) | Unproved or an unmapped native type is not a counterexample. |
| 27 | Characterize exact and underlying Go domains | HOLDS | Specify the grammar/environment and whole-domain laws; ~string needs broader association, not automatic rejection as unsound. | owner | [TextBox](../docs/type-system-architecture.md#statement-one-textbox) | Aliases, definitions and four test examples are distinguished from a quantified characterization. |
| 28 | Bind the actual Go baseline rather than the moving live spec | HOLDS | Keep Go 1.25 as the example's explicit qualification baseline; require matching future compiler/spec inputs. | owner | [Qualification](../docs/type-system-architecture.md#implementation-and-qualification-sequence) | Current Go 1.27 documentation is not reported as an executed 1.25 result. |
| 29 | Use a scripted proof producer first; no general search prerequisite | HOLDS | Separate proof production from checking in the minimal sequence. | owner | [Sequence](../docs/type-system-architecture.md#implementation-and-qualification-sequence) | No solver performance or completeness claim. |
| 30 | Check the native-model adequacy direction that preservation needs | HOLDS | Require admitted actual cases to be covered by the model; clarify specification well-formedness versus compiler acceptance. | owner | [Adequacy](../docs/type-system-architecture.md#four-assurance-boundaries) | Positive compilations alone cannot exclude omitted counterexamples or C++ no-diagnostic cases. |
| 31 | Add decisive negative cases | HOLDS | Record binding, helper, assumption, dependency, budget and native-profile mutations as future acceptance criteria. | owner | [Qualification](../docs/type-system-architecture.md#implementation-and-qualification-sequence) | Documentation-only criteria; no test execution claimed. |
| 32 | Stress the model with TypeScript, Rust and C++ distinctions | HOLDS | Preserve language-owned judgment and checking-point differences. | owner | [Languages](../docs/type-system-architecture.md#other-native-languages) | Primary sources support bottom/alias/constraint-context distinctions. |
| 33 | Define the strength of Rep; allow justified many-to-one relations | HOLDS | Use a static schema association first; value relations/adapters are separate scope. | owner | [Representation](../docs/type-system-architecture.md#representation-and-selected-laws) | Optional inactive payloads or nil/empty values acquire no runtime claim. |
| 34 | Expose the three elaboration/adequacy boundaries | HOLDS | Separate authoritative source projection, actual target projection and native-model adequacy. | owner | [Assurance](../docs/type-system-architecture.md#four-assurance-boundaries) | Proof validity cannot hide an incorrect artifact projection. |
| 35 | Retain portable identity and assumption/dependency association | NEEDS ADAPTATION | Reuse existing canonical identities/snapshots; reconstruct local handles; add no universal hash, persistent manifest or administrative provenance workflow. | owner | [Binding](../docs/type-system-architecture.md#binding-and-the-expected-statement) | Program ownership is preserved even for equal canonical inputs. |
| 36 | Implement/promote the foundation, bridges and proofs | DEFERRED | Trigger: owner authorizes implementation, Logos admits the required checker contract/executable, and the scoped Horos/native bridges are implemented and qualified. | affected component owners | [Sequence](../docs/type-system-architecture.md#implementation-and-qualification-sequence) | All production capabilities remain explicitly unavailable in this documentation task. |
| 37 | Make missing capabilities and unsatisfied obligations block stronger claims | HOLDS | List precise unsupported/incomplete conditions; preserve current representation/enforcement profiles. | owner | [Limits](../docs/type-system-architecture.md#unsupported-conditions-and-forbidden-shortcuts) | No weaker fallback, erased constraint or undisclosed assumption becomes success. |

## Question

How should we implement a reusable model of type systems and their judgments,
then compose those models with type-generation mappings and Logos reasoning,
without duplicating the existing language semantics or building a second logical
framework? We need a conceptually sound architecture, ownership/dependency
boundaries and a small implementation sequence, before committing to public APIs.

The current task is consultation and design documentation only. Accepted advice
will refine the proposal; production implementation is a later authorization.
Please challenge the candidate model and its placement, rather than assuming a
new Bittype-owned framework is the desired answer.

Reuse finding: Logos research 0043 already recommends typed semantic packages and
a shared checking foundation. The new decision is how Horos/native type theories
become consumers, how Bittype binds their inputs/results, and which sound route
can establish symbolic proofs while the richer Logos production API is still
proposed. We are not reselecting a logical foundation from scratch.

## Context

### System overview and vocabulary

Bitspark has independent repositories with these responsibilities:

| Project or term | Meaning here |
| --- | --- |
| Ontos | Ground values, structural identity and encoding. Values include atoms and compound/list/map structures. It is not the logical term language with variables. |
| Horos | Ground data descriptors and conformance; also a canonical static declaration language with modules, ordinary kind definitions, generic parameters and applications. The declaration layer can describe Wire capabilities without making ground conformance capability-aware. |
| Bittype | A Go toolchain with generic abstract/concrete syntax trees (AST/CST), language-specific syntax, explicit resolution/analysis and AST-to-AST transformation. It consumes Horos as canonical input and generates TypeScript, Go, Rust and C++ type declarations. |
| Logos | Domain-independent logical representation, reasoning services and proof/account checking. Its current production account is a finite ground derivation; newer typed semantic-package work is described below. |
| Language integration | The code and declarative theory that connect one language's actual syntax/semantics to supported judgments. It owns those meanings; the general host does not. |
| Mapping integration | A type owner's selected native representation and laws connecting canonical/native meanings. Builtin mappings use the same mechanism as other libraries. |
| Kind | In Horos, a predicate/contract over static proper types or type constructors. Scalar, List, Record and Data are ordinary source-authored definitions, not kernel tags. |
| Proper type / constructor | A proper type describes values; a type constructor accepts type-level arguments and yields a proper type or another constructor. These are Horos's coarse sorts, not a universal native type classification. |
| Ground / rigid / logical variable | A ground logical term has no logic variables. A rigid symbol encodes an arbitrary but fixed object-language parameter. A logic variable is instantiated during reasoning. Encoding a rigid symbol as a ground atom does not by itself justify universal generalization. |
| Judgment / evidence / interpretation | A judgment is a question with a specific meaning; evidence supports it under declared assumptions; an interpretation connects objects in different semantic domains. None is supplied merely by a shared spelling. |

Static type generation is the entire product scope here. Function, Stream and
Cell are ordinary library contracts. Their runtime behavior, transport,
hydration, invocation, cancellation and value adapters are outside this work.
Neither their names nor the name Data may become generator, runtime or Logos
kernel cases.

Inspected sources: Bittype `93713a007cb28ccbf86bbaf3661cbc7bbaec8ec4`,
Horos `ad7c6f493d726a82e445875d42af319adce0e3ce`, and Logos
`d0df8bbb6914e8002a441cb79241da8c9a3db7b0`. Bittype's Go module targets Go
1.27, depends on Horos Go 0.10.2 and Ontos 0.16.0, and has no direct Logos
dependency. Code inspection is not a benchmark or an executed proof. There is no
measured throughput target or evidence that proof search is a current bottleneck.

### Existing Bittype architecture

The repository uses one Go module and a stage-first hierarchy:

```text
trees/{ast,cst,representation}/go   generic syntax carriers/envelopes
semantics/go                      immutable provider-owned semantic graph
contracts/{language,pipeline,parse,check,format,print}/go
contracts/{semantic,resolution,analysis,transform}/go
engine/go                         explicit registration and orchestration
syntax/<language>/go              language AST vocabulary
parse/syntax/<language>/go         parser implementations
check/syntax/<language>/go         syntax admission
format/syntax/<language>/go        AST -> CST implementations
print/syntax/go                   generic CST printer
resolve/files/horos/go            explicit Horos resolution
analyze/files/horos/go            analysis from an explicit receipt
provider/<language>/go            optional composition, no registration effects
transforms/horos-to-<language>/go  native type transformations
internal/horosmodel               private authoritative-model projection
```

The generic engine has no language list, canonical-language special case,
automatic provider discovery or hidden resolver. Adding a language changes its
packages and caller composition, not the kernel. The generic graph is not a
universal type algebra. Native syntax tooling exists for all four languages;
only Horos has production semantic providers. Syntax admission does not imply
full native typechecking.

Actual carrier excerpts, `semantics/go/model.go:18-72,106-145`:

```go
type EntityID uint32
type TermID uint32
type Key struct{ Owner, Name string }

type Parameter struct {
    Entity EntityID
    Sort   Key
}

type Term struct {
    Kind      TermKind // Reference, Application, ProviderExpression
    Target    EntityID
    Callee    TermID
    Arguments []TermID
    Operator  Key
    Encoding  Key
    Data      []byte
    Operands  []Operand
}

type Subject struct { Entity EntityID; Term TermID }
type Fact struct {
    Subject  Subject
    Relation Key
    Term     TermID
}
```

IDs are local addresses in one admitted immutable `Semantics[L]`, where L is a
Go language marker. Keys identify provider-owned vocabulary; they are not
registrations. A term references an entity, applies another term, or carries a
provider expression with explicit child operands. Operand locators identify
positions in the provider encoding. Coarse parameter Sort distinguishes proper
type/constructor roles, not a full kind contract. The constructor checks
structural integrity and clones storage; it never proves supplied facts true.
Containment is acyclic; recursive entity references can still describe recursive
language objects.

Actual report boundary, `contracts/semantic/go/contract.go:27-85`:

```go
type Obligation struct {
    Key    semantics.Key
    Domain semantics.Key
}
type Report struct {
    Obligation semantics.Key
    Status     Completion // Complete, Incomplete, Unsupported
    Verdict    Verdict    // Undetermined, Holds, Rejected
}
```

Providers interpret the obligation/domain. An incomplete/unsupported report
cannot carry Holds or Rejected. A result's `Valid()` means its publication is
well-associated, not that the language rules are independently proved. Snapshots
and receipts bind actual inputs, operation, provider profile/configuration and
requirements; no hidden analysis repairs missing evidence.

Actual transformation boundary, `contracts/transform/go/transformer.go:15-35`:

```go
type Input struct {
    Source   semantic.Snapshot
    Evidence []semantic.Result
}
type ASTTransformer interface {
    FromID() language.ID
    ToID() language.ID
    Transform(context.Context, Input) (representation.ASTSet, error)
}
```

The typed adapter admits supplied evidence before invoking the mapper
(`contracts/transform/go/adapter.go:65-91`). It never runs a provider itself.
`engine/go/transform.go:73-87` captures each resulting AST and discards the
previous leg's evidence before a later leg. Output ASTSet is syntax only.
Preservation support therefore needs an explicit association with the source,
actual output, mapping choices and requested claim; putting it inside AST nodes
or reusing source receipts as target evidence would violate the boundary.

Current native pagination generation consumes source-checked facts, retains real
binders and outputs ordinary open generics. Its representation profile explicitly
permits native constraint erasure. The `pagination-enforcement/1` profile returns
Unsupported (`internal/horosgen/generate.go:170-184`). Higher-order constructor
parameters and several other type forms are unsupported. No native admission
proof or generic theorem is implemented by today's transformer.

### Authoritative Horos operations and their limits

Public Horos no longer contains its old Logos dependency or proof/witness API.
That code was left in private history during the public cutover. `core/go/go.mod`
depends only on Ontos; `README.md:88-92` and `docs/releases/v0.6.0.md:14` describe
the exclusion. A new optional bridge must not be described as an existing public
capability.

Relevant signatures, `core/go/declarations/model.go:84` onward (reformatted):

```go
type Program struct{ inner *model.Program }
func Resolve(root core.Value, dependencies []core.Value) (*Program, error)
func ResolveWithLimits(root core.Value, dependencies []core.Value,
    limits Limits) (*Program, error)

func (p *Program) CompileKind(module Digest, expression core.Value) (KindHandle, error)
func (p *Program) Value(id DeclarationID) (TypeValue, error)
func (p *Program) Proper(value core.Value) (TypeValue, error)
func (p *Program) Apply(callee TypeValue, arguments []TypeValue) (TypeValue, error)
func (p *Program) Satisfies(value TypeValue, kind KindHandle) (Verdict, error)
func (p *Program) Includes(actual KindHandle, expected KindHandle) (Verdict, error)
func (p *Program) GroundDescriptor(value TypeValue) (core.Value, error)
```

Here `core.Value` is an Ontos value containing canonical declaration syntax;
Digest identifies an imported canonical module; DeclarationID identifies a
declaration within it. Limits bound resolution/query work. TypeValue and
KindHandle are sealed Program-owned handles. Current query verdicts are Proven,
Refuted and Unproved; malformed/foreign input and exhaustion are separate errors.
Inclusion currently returns Proven or Unproved, not Refuted: failed comparison
can lack a realizable counterexample. Required Unproved obligations prevent
successful admission but are not established logical negations.

Relevant internal definitions (`internal/declmodel/model.go:14`,
`specialize.go:9`, within `core/go`):

```go
type DeclarationID struct { Module Digest; Name string }
type ParameterID struct { Declaration DeclarationID; Index int }
type Assumption struct { parameter ParameterID; kind *Kind }
type Membership struct { kind *Kind; assumptions []Assumption }
type TypeValue struct { owner *Program; value *semanticValue }

func (p *Program) hasValue(value TypeValue) bool {
    return p != nil && value.owner == p && value.value != nil
}
```

`Kind` exposes an immutable public view through accessors with private storage;
`semanticValue` is private. Neither is a portable proof object. An integration can
inspect kind structure through that existing view. Even identical modules
admitted in two Programs do not permit handle
exchange. Published memberships retain original parameter assumptions. A new
query does not mutate the Program or infer a principal kind.

Control flow:

```text
canonical modules + explicit dependencies
  -> Horos ResolveWithLimits
  -> resolve imports and real parameter slots
  -> check sorts, primitive formation, kinds and every declaration body
  -> checked Program
  -> Bittype projection of entities, contracts, bodies and original assumptions
  -> explicit resolution receipt
  -> Horos analyzer republishes that graph/reports from the receipt
```

The Bittype resolver calls `decl.ResolveWithLimits` directly
(`resolve/files/horos/go/resolver.go:18-23`). Its analyzer currently reuses the
receipt rather than performing another analysis
(`analyze/files/horos/go/analyze.go:36-68`). Merely requesting another obligation
key does not add reasoning capability; unknown keys return Unsupported.

For generic application, Horos checks nonempty full saturation and finite sorts,
then admits **every** argument against the original declaration's kind, binds
`ParameterID{declaration,slot}`, and evaluates under that environment. An unused
argument is still admitted. A narrower requested result guarantee never repairs
an invalid declaration. An application may return a constructor; no partial
application or source type lambda exists. Known constructors expose checked
bodies/environments; abstract constructors expose declared guarantees.

Recursive kinds use explicit guarded simultaneous greatest fixed points (`nu`).
Named references are acyclic; each recursive cycle crosses a type-child descent.
The checker builds Boolean obligations, initializes candidate pairs true, then
propagates failures to a stable solution (`internal/kindmodel/membership.go:3`).
A revisited pair alone is not proof. Inclusion uses finite cuts over positive
recursive hedge automata, whose free observation universe is larger than admitted
Horos types. For example, an infinite `nu X. Wire<X>` observation need not be a
realizable admitted type. Ordinary least-model Horn recursion or an unchecked
proof cycle is not an adequate replacement.

GroundDescriptor separately requires a closed proper value, membership in the
ordinary bundled ground kind and unchanged ground well-formedness. Static kind
membership neither validates arbitrary native values nor grants Wire authority.

### Logos: current production boundary and overlapping design

The production account is a finite directed acyclic graph of ground conclusions.
It has Rule, Builtin and Leaf nodes. Rule nodes apply a clause substitution and
refer to previously established premises; Builtin nodes ask admitted theories;
Leaf nodes use externally supplied evidence under an explicit verifier.
An unsupported decision must not become false. The general checker has no
Horos, Go, Data or Page semantics.

Actual production excerpts, `contract/go/term.go:99` and
`contract/go/derivation.go:150` (nonessential members omitted):

```go
type Subst map[VarId]Term
type Term interface {
    Equal(Term) bool
    IsGround() bool
    EraseGround() (core.Value, bool)
    Substitute(Subst) Term
    String() string
    isTerm()
}
type RuleNode struct {
    Conc GroundLiteral
    Clause ClauseId
    Subst Subst
    Premises []Premise
    Meta NodeMeta
}
type BuiltinNode struct { Conc GroundLiteral }
type LeafNode struct { Conc GroundLiteral; Evidence Evidence }
type Derivation struct {
    Nodes []Node
    Root uint64
    ClosedWorld []ClosedModelCert
}
```

Terms are atoms, tuples or first-order variables. A Clause is a head with body
premises; ClauseId identifies it in the supplied program. GroundLiteral is a
ground logical claim. Premise identifies earlier positive support or checked
closed-world side conditions, including negation and counting;
NodeMeta carries policy metadata. ClosedModelCert supports explicitly checked
closed-world reasoning; search failure alone is not negation. The checker admits
the inputs, checks any certificates, traverses nodes in order, checks substituted
clauses/premises, then checks that the root matches the requested claim.

Actual extension signatures, `contract/go/spi.go:88` onward:

```go
type Theory interface {
    Owns(Symbol) bool
    Decide(GroundLiteral) TheoryOutcome
}
type GroundDecider interface {
    DecideGround(GroundLiteral) completion.TheoryGround
}
type LeafVerifier[C any] interface {
    Extensional() *PredicateSet
    Verify(GroundLiteral, Evidence, C) LeafOutcome
}
type RuleAdmissibility[C any] interface {
    Admissible(Clause, NodeMeta, C) RuleOutcome
}
```

Symbol/PredicateSet identify consumer vocabulary. C is an explicitly supplied
immutable application context; Evidence is an opaque value interpreted by its
leaf verifier. Current LeafOutcome has LeafAdmit/LeafReject, and RuleOutcome has
RuleAdmit/RuleReject: these hooks do not provide an Unproved/Exhausted arm.
Rejecting a supplied certificate is not proof of the negated source judgment.
A Horos Unproved/exhausted call must remain incomplete at the consumer outcome
boundary, not be silently relabeled logical nonmembership.

The Theory path is different: current readers use completion-aware DecideGround,
so stopped computation remains noncompletion, not logical Fails. Its legacy
Decide form has Holds/Fails/NotMine only. Theories are construction-configured,
deterministic ground decisions. Distinguish this Go `Theory` extension from a
language theory (object-language judgments/rules) and the logical foundation
(the calculus used to check evidence); they have different responsibilities.

RuleAdmissibility receives neither the instantiated substitution, cited premises
nor the enclosing derivation. It cannot alone check a generic scope-discharge
step merely because its name mentions admissibility. A consumer certificate
route would need explicit bound proof data and a real checker. No production
universal-introduction, induction or coinduction operation was found here.

Search and checking are separate. `kernel/go/seam/spi.go:52` exposes:

```go
type EngineProvider[C any] interface {
    Descriptor() ProviderDescriptor
    Probe(RequiredProfile, ProgramAnalysis, ExecutionEnvironment[C]) Eligibility
    Prepare(context.Context, PlanEnvelope, ExecutionEnvironment[C]) (
        PrepareOutcome[EngineInstance[C]], error)
}
type EngineInstance[C any] interface {
    Solve(context.Context, ExecutionEnvironment[C], RunLimits) (CandidateStream, error)
}
```

These named inputs describe requested capability, analyzed program, explicitly
bound environment, selected plan and resource limits. CandidateStream returns
candidate answers/substitutions with optional derivations, not accepted proofs.
The runtime's certification step checks them independently. The convenient
`kernel/go/embed.Query` constructs guest-backed providers; it exposes custom
leaves but no direct arbitrary Theory field. The lower search-side HostSolver
returns bindings for a goal/substitution. Those bindings still require a paired
checker Theory. Adding a ground checker does not supply a search procedure.

**Substantial overlapping work already exists in Logos.**
[Research 0043](https://github.com/Bitspark/logos/blob/d0df8bbb6914e8002a441cb79241da8c9a3db7b0/research-docs/0043-realizing-the-logos-vision.md)
and its expert advice recommend a semantic graph, a common checking foundation
and specialized calculi with justified connections. We reuse that advice, not
re-request the general choice of a logical foundation. Its package components
are syntax/contexts, semantic meaning, judgments, evidence calculi, adequacy or
explicit assumptions, composition laws and actual operational services.

[Proposed ADR 0038](https://github.com/Bitspark/logos/blob/d0df8bbb6914e8002a441cb79241da8c9a3db7b0/docs/architecture/decisions/0038-typed-judgments-and-semantic-packages.md)
uses `T; Gamma |- e : J`: T is the admitted foundation/dependencies, Gamma the
explicit assumptions, e evidence and J the exact typed judgment. Representation,
execution and certification have independent support states. The ADR explicitly
overrides no accepted production contract while proposed.

`docs/design/0014-typed-package-production.md:74` proposes these shapes, **not
exported Go interfaces**:

```text
Statement<F> {
  foundation: FoundationInstance,
  packages: Set<PackageInstance<F>>,
  context: Ordered<Binding<F>>,
  objects: Map<ObjectName, BoundObject<F>>,
  judgment: JudgmentApplication<F>
}
JudgmentApplication<F> {
  declaration: QualifiedDeclaration,
  arguments: Ordered<Expr<F>>
}
SemanticPackage<F> {
  identity, foundation, parameters, signature,
  semantic_definitions, judgments, evidence_rules,
  adequacy, operations, services
}
Justification<F> =
  Checked { proposition, evidence_ref }
  | Assumed { proposition, dependency_identity }
```

F identifies a chosen logical foundation. Package instances bind parameters and
meanings by identity. The draft canonical contract in design 0015 distinguishes
Type, Term, Proposition and Proof under F, with ordered context entries
`object(type)` or `assume(proposition)`. Qualified references replace bare names
as authority. Design 0016 instantiates a proposed Horn/stable/rational portfolio,
not a Horos/native type theory. Production admission remains outstanding.

There is real **test-private Rust** foundation code at
`conformance/rs/tests/heterogeneous/foundation.rs:21,263-293`:

```rust
pub enum Expr {
    Sort(u32), Var(u32), Const(String),
    Pi { domain: Box<Expr>, body: Box<Expr> },
    Lam { domain: Box<Expr>, body: Box<Expr> },
    App { fun: Box<Expr>, arg: Box<Expr> },
    // also dependent pairs, Nat/Bool eliminators and identity/transport
}
pub fn infer(signature: &Signature, context: &[Expr],
    term: &Expr, budget: &mut Budget) -> Outcome<Expr>;
pub fn check(signature: &Signature, context: &[Expr],
    term: &Expr, expected: &Expr, budget: &mut Budget) -> Outcome<Expr>;
```

Pi is a dependent function type, Lam its introduction and App its elimination.
Var uses de Bruijn indices (binder positions, not names); Signature binds checked
constant declarations. Conversion/substitution consume bounded work. The selected
mathematical fragment has stratified universes, dependent functions/pairs,
natural/Boolean eliminators and identity/reflexivity/transport. It has no general
recursion, arbitrary user rewrite system or general inductive declarations.
Its scoped selection explicitly does not activate production semantics, admit a
checker executable or claim a metatheoretic soundness proof. Bittype cannot import
this as a stable Go API today.

The existing implementation direction is Go for Bittype tooling. A proposed
cross-language process/library boundary, a Go realization of a Logos foundation,
or waiting for Logos production promotion must be an explicit dependency choice,
not an assumption that experimental Rust code is already a usable Go library.

The remaining decision is how type-system theories and Bittype's inputs should
consume this direction, which current route can safely deliver the first slice,
and what shared integration concepts are genuinely missing. A new independent
foundation, proof grammar or rule evaluator in Bittype would duplicate this work.

### Proposed reusable model (not an implemented API)

The current Bittype proposal identifies seven conceptual parts:

1. Configuration: language revision/options, module environment, consumer operations.
2. Language-owned term vocabulary: types, kinds, constraints, applications and
   additional static terms such as lifetimes or constant expressions.
3. Binding/context: real identity, lexical scope, capture-avoiding substitution,
   source and target binders distinguished from logical variables.
4. Judgment declarations: named questions with argument domains and meanings.
5. Rule modules: inference, equations, side conditions and admitted foundations.
6. Coverage: supported constructs, operations, argument domains and checking points.
7. Reasoning support: derivations, procedures or certificates for each claim.

Judgment roles include formation, kind membership, equality/normalization,
subtyping, assignability, constraint satisfaction and application. These are not
a closed enum or mandatory methods every language implements. Rust outlives and
trait obligations must be ordinary extensions. Coarse sort, kind membership,
subtyping, assignability and runtime conformance remain distinct relations.

The model would express both within-language questions and cross-language
preservation. The latter needs an intended interpretation, fixed before the
answer, and a selected representation relation. For a source declaration H,
native declaration N, configuration C and declared native domain D:

```text
for all n in D:
  NativeAccepts(N,n)
    implies SourceAccepts(H,Interpret_C(n))
       and Represents_C(SourceResult(H,Interpret_C(n)),NativeResult(N,n))
```

NativeAccepts includes the promised checking point (for example, generic
application formation, alias use or demanded associated-type normalization) and permitted consumer
operations. The interpretation must cover accepted cases in the claimed domain;
it cannot pick any convenient source type afterward. Relative completeness is a
different statement: every source application in the independently specified
supported source domain has a native representation that is admitted. Refusing
everything is not useful preservation. Search incompleteness is not refutation.

Language theories should be declarative owned content. Native adapters connect
that content to actual syntax and dependencies; mapping libraries supply their
laws; the host composes them. The question is how much shared machinery is
actually missing, and whether it belongs in Logos, an integration package or
Bittype. No package split or dependency change has been selected.

### Concrete cases the architecture must explain

#### A. Source membership and Go admission

Horos (with an explicit ordinary prelude import as h):

```text
type TextBox<T: h.Text> = { value: T; };
type Label = TextBox<text>;
```

Selected associations are canonical text -> Go string and bool -> Go bool.
An unconstrained candidate `type TextBox[T any] struct { Value T }` admits bool,
which is a source counterexample. The desired model should derive that failure,
not rely on a plugin's self-reported preservation flag.

A second candidate is:

```go
type TextBox[T interface{ string }] struct { Value T }
type StringAlias = string
type NamedString string
```

For an explicitly bounded fragment of closed concrete arguments, exact type
terms distinguish string/its aliases from bool and the defined NamedString.
The model must characterize the complete admitted domain, not just list four
tests. Changing the term to `~string` changes that domain. Relevant source:
[Go specification](https://go.dev/ref/spec#General_interfaces); the existing
walkthrough fixes Go 1.25 as its baseline. Other operations/versions need their
own modeled coverage. No claim about all Go string values being valid UTF-8 is
part of this static example.

Even correct admission leaves the field/value representation law, required
presence, binder correspondence and capture-avoiding substitution to establish.

#### B. Symbolic pagination, not enumeration

```text
type Page<T> = { items: list<T>; next?: text; };
type DataPage<D: h.Data>: h.Data = Page<D>;
type Person = { name: text; active: bool; };
type People = DataPage<Person>;
type Nested = DataPage<DataPage<Person>>;
```

The selected Go representation family is:

```go
type BittypeOptional[V any] struct { Present bool; Value V }
type Page[N any] struct {
    Items []N
    Next BittypeOptional[string]
}
```

The owner's central requirement is a symbolic theorem for arbitrary T and N:

```text
T : Data and Rep_C(T,N)
  implies Page<T> : Data and Rep_C(Page<T>,GoPage<N>)
```

Rep_C is an established selected representation association including native
well-formedness. It is not a preference for a type name. Introduce fresh rigid
source/target parameters, assume those premises, then compose list, optional
field and record laws. Do not enumerate T, make Data a privileged kind, or infer
that all Data types already have a native mapping.

Current documentation gives illustrative consumer rules such as:

```prolog
kind_in(C,P,T,K) :- proof_scope(C,P), assumed(C,P,member(T,K)).
kind_in(C,P,list(T),K) :- kind_closure(C,K,list), kind_in(C,P,T,K).
rep_in(C,P,list(T),slice(N)) :- mapping_law(C,list_to_slice), rep_in(C,P,T,N).
```

Uppercase tokens denote logic variables, lower-case symbols relation/constructor
names. C binds the actual source/target/mapping configuration and P the local
proof scope. `kind_closure` must be an established consequence for the resolved
ordinary kind; a mapping law is an admitted equation with a stated meaning.
Neither is established merely by declaring the fact. Complete record rules walk
actual paired fields; opening bodies must preserve their original binders.

The finite ground derivation under scope P is **not yet universal evidence**.
The design owes lawful assumption discharge/generalization, freshness, all
reachable dependencies, substitution stability and actual named-family
correspondence. The next model must make those obligations concrete and identify
the checker that enforces them. It must also preserve the original application
contracts, including ignored parameters.

Native admission is independent: `Page[N any]` supplies a representation family
but does not enforce Data. An enforcing candidate must show that every admitted
native argument supplies the theorem's premises under the intended interpretation.
Nested pagination reuses the inner application's result proofs. This is not a
finite-selector-only project.

#### C. Other languages are stress cases, not interchangeable Go encodings

Earlier consultations established these design constraints; this request does
not ask for another round of encoding recipes:

| Language | Distinctions the shared model must be able to express |
| --- | --- |
| TypeScript | Structural assignability vs identity; never/any; variance; conditional/mapped-type normalization; checking point and declaration merging. Independent upper-bound witnesses can admit bottom; finite selector normalization differs from generic formation rejection. |
| Go | Exact vs underlying type sets, nominal identity vs aliases, inference and forwarding. A private singleton representation can still be inferred into another generic parameter; private names alone do not establish issuance or admission. |
| Rust | Trait obligations, associated-type normalization, coherence, sealing, lifetimes and sizedness. An ordinary alias bound and a demanded projection gate are different; effective Sized restrictions must exclude unsupported trait-object witnesses. |
| C++ | Template argument binding, requires/constraint satisfaction, identities, aliases, specialization, visibility and static evaluation. Supporting-template customization and changing satisfaction across contexts must be included or explicitly excluded. |

No realistic complete decision procedure for all these languages is assumed.
Unsupported constructs should be explicit. A theorem proved in an incomplete
model can still misdescribe a real compiler; adequacy and compiler regression
testing are separate from internal proof checking.

### Prior advice, reuse and candidate approaches

[Bittype research 0007](0007-type-owned-generator-plugins.md) separates Use
(produce a native type expression), Publish (expose a declaration) and Define
(fill a reserved declaration/helper). Type owners supply mappings, including
builtins. Mapping selection is exact and explicit; an override replaces a binding,
not a fallback. No runtime protocol is implemented by selecting native syntax.

[Research 0008](0008-generic-generation-protocol.md) adds paired formal-parameter
Introduce and actual-argument Apply operations, real binder-to-scoped-product
association, arbitrary convention-owned immutable support, and result support for
nested applications. The host enforces scope and invocation association, not a
universal witness tuple. Imported fixed mapping choices cannot change through an
element override. Source facts and configuration-specific representation evidence
are separate. This is a proposed composition protocol, not implemented software.

[Research 0009](0009-native-kind-constraint-soundness.md) supplies the native
counterexamples above and requires intended interpretation, explicit domain,
checking point, consumer operations and extension policy. Representation collisions
must respect the queried property; full injectivity is not always necessary.
It did not design a reusable language-theory foundation or establish a complete
native enforcing encoding.

Two independent design axes to assess:

- **Model ownership: direct use of Logos.** Integrations publish language-owned signatures,
  judgments and laws using Logos's existing/proposed semantic-package model;
  Bittype adds only explicit input/result association and orchestration.
- **Model ownership: small shared modeling boundary.** Reusable binding/context and typed judgment
  descriptions above existing carriers, explicitly lowered to Logos. Justify each
  addition by a case that direct Logos use cannot handle; no second logical engine.
- **Evidence route: authoritative procedure or independent certificate.** A Horos adapter can supply declared trusted
  procedure results; native fragments and representation laws are explicit.
  Independent certificate replay can follow. This is incremental but changes the
  admitted foundation and must be visible to the requester. This route can
  coexist with either ownership choice; the alternatives are not mutually exclusive.

Rejected shortcuts: a universal source/native Type enum; named Data or contract
cases; modeling every relation as compatibility; duplicating the Horos checker
inside generation; passing foreign Program handles; treating absence of proof as
false; claiming a symbolic theorem from a few instances; asserting the final
preservation claim as an axiom; assuming source receipts certify new target ASTs;
or labeling an unexecuted proposal as production capability.

We have no preference for the smallest amount of code at the expense of semantic
adequacy. We do want a small coherent first slice rather than premature APIs for
every language feature. The static declaration model, proof foundation and runtime
behavior must remain separate. Repository resource rules require portable tests
on shared/CI workers and forbid all Logos workloads on the workstation.

## Questions for the expert

1. Given Logos's existing semantic-package direction, how should language-owned
   type systems, judgments, binding and assumptions become consumers of it?
   Which parts can use it directly, and what reusable boundary, if any, is missing? Please
   distinguish object-language semantics, proof foundations, evidence and model
   adequacy rather than selecting a universal type algebra.
2. How would you modularize the implementation and its declarative assets across
   the existing Bittype, Horos and Logos boundaries? Give a concrete dependency
   graph, package-role layout and small Go-facing interfaces or data shapes where
   they clarify the boundary. Address current APIs versus proposed Logos APIs,
   explicit input/output association, and type-owned mapping laws without
   premature kernel changes or an accidental dependency cycle.
3. What exact reasoning/checking route should establish the symbolic pagination
   theorem and its reusable nested applications, while keeping native admission
   separate? Explain lawful generalization/substitution, original constraints,
   trusted procedure vs independently checked evidence, and the limits of the
   existing ground DAG. Include how recursive Horos kinds are safely scoped or
   integrated without substituting least-fixed-point semantics. Recommend a
   concrete route between a consumer-owned certificate checker on current APIs
   and the proposed typed-package route, identifying its owning checker, required
   API/soundness prerequisites and what can proceed before production promotion.
4. What first vertical slice and qualification sequence would demonstrate this
   architecture? Use Horos/Go TextBox and symbolic Page, then show how TypeScript,
   Rust and C++ stress the abstraction. Distinguish proof search, independent
   checking, native-model adequacy, logical refutation and incomplete execution;
   identify decisive negative tests and explicit deferral triggers.
5. What are we missing or assuming incorrectly, and what alternative design would
   you recommend? You may reject the proposed shared layer or change the sequence
   if a more coherent architecture follows from the actual evidence.
