# AST transformation and pipeline design

**Current architecture, 8 November 2026.** Generation maps admitted source
semantics to intended native types, then constructs their ASTs. File organization
is a later operation. Horos-to-TypeScript, Go, Rust and C++ representation mappings
are implemented through the [declarative SPI](generation-spi.md) and
[kernel-owned composition](kernel-composition.md). They do not register
ASTTransformer edges. [Documentation map](README.md) identifies the owning guides.

The output scope is types and interfaces, including open generics and applications.
Executable bodies, value adapters, serialization and protocol implementations
remain outside this work. Function, Stream and Cell are ordinary libraries;
neither the generation kernel nor Wire/runtime learns their names.

## Target pipeline

```text
supplied Horos text or canonical module values + supplied dependencies
  -> Horos ASTs and explicit resolution/analysis
  -> captured source Semantics + receipts
  -> PlanGeneration: selected mappings -> closed native Program
  -> GenerationPlan.Run: target construction -> checked generation.Output
  -> TypeScript Files: placement + dependency environment -> checked Source
  -> shared Files printer -> paths and text
```

The application owns IO, dependency acquisition and target configuration. Horos
owns declaration meaning, kinds, binding, application and ground extraction.
Bittype's Horos stages project those meanings into the shared fact model.
Generation does not invoke an implicit resolver or analyzer.

All four targets construct ASTs and have Syntax formatting/printing. The complete
generated-output handoff to Files currently exists for TypeScript. Go/Rust/C++
placement remains [proposed](native-file-placement.md). Rust visibility and C++
include emission still violate the intended Files boundary; their correction is
tracked in [ACX-004](architecture-cases/0004-generation-file-syntax.md).

The [pagination profile](semantic-generation.md) maps text, bool, lists, named
closed records, optional fields and ordinary proper generic parameters and
applications. The application selects public roots. Native construction follows
the representation dependencies demanded by those roots; source admission covers
the source profile's reachable module closure. Unsupported demanded mappings fail
without usable output. An exact external mapping need not expand its source body.

Open generics remain open; ground specialization is not a prerequisite. A single
generic convention relates source binders, native parameters and complete mapped
actuals. It can reorder, omit or add native arguments while preserving their
scope. Omitted source arguments still require source admission. See
[generic elaboration](generic-elaboration.md).

### Semantic mapping and target AST construction

| Operation | Owner and product |
| --- | --- |
| Source resolution/analysis | Source-language stages establish references, applications and facts for an exact capture and environment. |
| Semantic mapping | Supplied mapping libraries choose native expressions, declarations and generic conventions. The kernel selects and composes them into a closed native Program. |
| AST construction | The native target realizes that Program with fresh AST sites and bound references. Syntax checking admits the result. |
| File organization | The native Files stage uses one consumer placement to choose module syntax, aliases and visibility, then checks realized native bindings. |
| Text emission | Syntax formatting constructs CSTs; the shared printer emits their leaves. |

These implemented interfaces delimit planning and construction. The consumer
registers libraries and a target; registration invokes Describe, the planner
invokes target checks, and Run invokes Build. Excerpts use the generation package.

```go
type MappingLibrary[F, T, R any] interface {
    Describe(from, to representation.Language) (MappingDescription[F,T,R], error)
}
type Target[T any] interface {
    CheckExpression(context.Context, Expression[T]) error
    CheckDeclaration(context.Context, NativeDeclaration[T]) error
    CheckProgram(context.Context, Program[T]) error
    Build(context.Context, BuildInput[T]) (Assembly, error)
}
```

Build receives no source reader or mapping callback. Native intent is not an
analyzed `Semantics[T]`: target analysis establishes facts about the actual
constructed/organized program. An external Cell use identifies its public package
symbol before any import exists. TypeScript Files later verifies that the emitted
occurrence resolves to that symbol even if an alias is necessary.

The engine entry points are implemented as follows (`compose` denotes
`transform/syntax/shared/generation/go`):

```go
func (e *Engine) PlanGeneration(ctx context.Context, r compose.Request) (generation.GenerationPlan, error)
func (e *Engine) Generate(ctx context.Context, r compose.Request) (generation.Output, error)
func (p GenerationPlan) Run(context.Context) (Output, error)
```

Generate plans and runs once. Output retains its checked snapshot, requested
public declaration sites and external occurrences. It establishes neither
native kind enforcement nor a proof that the chosen representation is correct.
The [Page/Cell walkthrough](page-cell-example.md) follows the complete consumer
composition; the [context matrix](abstraction-levels.md) distinguishes the products.

## Go layout

One root Go module contains a stage-first matrix. [LAYOUT.md](../LAYOUT.md) owns
the complete hierarchy and import rules. The generation slice is:

```text
contracts/generation/go/                      descriptions, invariant planner, output
engine/go/                                   public operations and endpoint admission
transform/syntax/shared/generation/go/        registry and selected composition
transform/syntax/shared/target/go/            shared AST construction mechanics
transform/syntax/horos/source/go/             Horos evidence reader
transform/syntax/horos/mapping/go/            source patterns and declaration policy
transform/syntax/horos-to-<native>/mapping/go/ representation choices
transform/syntax/<native>/model/go/           native intent vocabulary
transform/syntax/<native>/target/go/          native intent to AST
format/files/typescript/provider/go/          implemented generated Files handoff
```

Mappings depend on source interpretation and native models. Targets depend on
native models and construction mechanics. Neither the kernel nor a target imports
concrete mapping libraries. There are no per-pair assembly factories or presets.
The consumer supplies implementations; the kernel performs their composition.

## Language identity and registration

Language markers implement the existing language package contract:

```go
type Provider interface {
    ID() ID
    Name() string
}
```

A registered representation captures language ID, profile and actual Go marker
identity. Syntax/Files stage adapters also capture their explicitly supplied
collaborators. Matching strings alone do not authorize a foreign typed payload.
MappingLibrary and Target are different contracts: they do not pretend to be
single-language parse/format providers. Their endpoint identities are supplied
at registration. See [registration signatures](kernel-composition.md#registration-and-request)
and [independent text stages](engine.md).

## Typed stage contract

The small pipeline package defines the implemented function type:

```go
type Stage[I, O any] func(context.Context, I) (O, error)
```

It imports context from the standard library, not ASTs or languages. A stage must
respect input ownership and error semantics; this type alone does not enforce
those guarantees. Shared generic carriers admit structure. Language checkers own
syntax, and semantic stages own their language judgments. A shape-compatible
function is not automatically a registered generation or semantic capability.

## ASTTransformer and syntax engine

ASTTransformer remains an independent syntax-transformation contract:

```go
type ASTTransformer interface {
    FromID() language.ID
    ToID() language.ID
    Transform(context.Context, transform.Input) (representation.ASTSet, error)
}
func (e *Engine) TransformTo(ctx context.Context, input transform.Input, targetID language.ID) (representation.ASTSet, error)
```

The application supplies edges and the engine invokes them. Typed adapters admit
explicit semantic evidence and return owned AST sets, not an AST[any] payload.
The engine rejects missing or ambiguous shortest routes, duplicate edges and
endpoint mismatches. Each leg admits its output and requires fresh evidence for
a later semantic operation. Failure never retries a weaker route. A same-language
request returns a checked snapshot, not implicit normalization.

Generation registrations do not add these edges. There are no production native
canonicalizers or automatic TypeScript -> Horos -> Rust generation routes. Future
canonicalization needs explicit semantic choices: native numeric and nullable
forms do not uniquely identify source domains or absence semantics. Reviewing an
adopted canonical definition can make it authoritative; parsing native syntax
alone cannot. The [engine reference](engine.md) owns routing details.

## Composition and dependencies

| Component | Allowed responsibility/dependencies |
| --- | --- |
| Generic AST/CST and representation | Standard-library storage/identity plus their shared carrier contracts; no language rules |
| Generic pipeline | context; no AST or language dependency |
| Language identity | Stable ID/name contract; no language catalogue |
| Semantic graph and capture | Generic addresses, evidence and ownership; no source/native theory |
| Generation planner and registry | Generic captures, scoped substitution, selection and native intent; no concrete language mappings |
| Horos integration | Horos-owned meanings and Ontos module values under explicit package allowances |
| Native generation model/target | Its native forms and syntax; no Horos/library semantics |
| Native Files integration | Placement, dependency environment and native binding checks |
| Optional verification | Agnostic artifact association, library-owned interpretations/laws and explicit Logos checking |

The [dependency guide](dependencies.md) records module versions. Package isolation
inside the root module does not erase module-level dependencies. TypeScript and
Rust semantic integrations use explicitly provisioned compiler helpers; their
Syntax-only parsing/formatting remains independent of those helpers.

Generating several targets means separate requests with explicit target
configuration. Applications choose scheduling; input captures remain read-only
and runs own their output storage. No central target or contract catalogue is
needed.

## Evidence and next increments

[GEN-01](generation-gaps.md#gen-01) records the correction to common semantic
navigation, its scope and its independent verification. [GEN-02](generation-gaps.md#gen-02) records absent
native constraint enforcement. The representation profile explicitly permits
constraint erasure; requesting the enforcement profile fails Unsupported.
[Native placement](native-file-placement.md) records the remaining Files work.

Generation, source membership, native compiler acceptance and cross-language
verification answer different questions. [Verification](verification.md) uses
captured artifacts and separately selected mapping laws; it is not automatically
performed by Generate. The [type-system model](type-system-model.md),
[constraint proposal](native-kind-constraints.md), [Logos preservation proposal](logos-preservation.md)
and [Go admission proposal](go-native-admission.md) retain their explicit limits.
Marker-method emission is a proposed scope extension, not implemented type output.

Use the existing Go/native and documentation workflows for relevant qualification.
No extra manifest or provenance workflow is required. Public builders, a general
application CLI, runtime bindings and arbitrary native-language adoption remain
outside the implemented scope.
