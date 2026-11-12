# Kernel-owned generation composition

**Implemented, 8 November 2026.** The application registers source readers,
native targets and declarative mapping libraries. The kernel plans the requested
declarations and composes the target program. Run constructs ASTs; Files realizes
their placement. No library supplies a complete pipeline or learns its consumers.

[Declarative SPI](generation-spi.md) · [Planner walkthrough](generation-planner.md) ·
[Generation assembly correction](architecture-cases/0002-generation-assembly.md) ·
[Runnable Page/Cell walkthrough](page-cell-example.md)

## Uniform layout

Implementations follow stage/dimension/scope/component/go. Scope is a language,
an ordered pair or shared. The native model describes intended types; the target
constructs their syntax. Analyzed native facts remain under semantics.

```text
contracts/generation/go/                     descriptions, scopes, planning and output
engine/go/                                  endpoint checks and public entry points
transform/syntax/shared/generation/go/       registration, selection and composition
transform/syntax/shared/target/go/           scoped AST construction and relocation
transform/syntax/horos/source/go/            horos interpretation and reader
transform/syntax/horos/mapping/go/           profile identities and the mapping library function type
transform/syntax/horos/mapping/{pattern,declarations,naming}/go/ horos patterns, declaration recipe, naming
transform/syntax/<native>/model/go/          intended native forms, no AST construction
transform/syntax/<native>/target/go/         native grammar and construction
transform/syntax/horos-to-<native>/mapping/<name>/go/ one declarative representation library each
format/files/typescript/provider/go/        placement, imports, exports and binding checks
```

Parsing, checking, formatting, resolution and analysis keep their existing stage
paths. The [full layout](../LAYOUT.md) documents those capabilities. Mappings
depend on their source's reader and target models, never AST construction
packages. Targets depend on their model and shared constructor, never source
readers or mapping libraries. The shared kernel imports neither; it reads the
shared semantic graph itself.

The guarded planner in contracts/generation owns issuance of the closed Program,
scoped substitutions and GenerationPlan. The stage registry supplies selected
descriptions to it. Keeping those private fields with their invariant interpreter
prevents an unchecked public constructor for a supposedly completed plan.

## Registration and request

These are implemented generic functions in the composition package. The caller
supplies concrete implementations. The registry admits language/profile/Go
marker and source API identities, freezes descriptions and invokes no provider
under its lock. [Registry](../transform/syntax/shared/generation/go/registry.go)

```go
func RegisterSource[F, R any](
    r *Registry, language representation.Language,
    p generation.SourceProvider[F,R],
) error

func RegisterTarget[T any](
    r *Registry, language representation.Language,
    p generation.Target[T],
) error

func RegisterMapping[F, T, R any](
    r *Registry, from, to representation.Language,
    p generation.MappingLibrary[F,T,R],
) (MappingRef, error)

type Request struct {
    Input   transform.Input
    To      representation.Language
    Roots   []ast.NodeRef
    Output  ast.TreeID
    Profile string
    Limits  generation.Limits
}
```

One source and target registration identify each configured language. Compatible
mapping libraries participate in the requested generation profile. The request
does not repeat a list of mappings. Planning chooses exact rules at demanded
occurrences. Multiple demanded offers fail as ambiguous; registration order is
not a preference.

An explicit replacement names both registrations and their binding IDs. The
registry owns these handles, so another registry's handles cannot select them.

```go
type BindingRef struct { Provider MappingRef; ID string }
type Replacement struct {
    Selector generation.Selector
    Expected BindingRef
    With     BindingRef
}
func (r *Registry) Replace(Replacement) error
```

A replacement must retain the exact selector and identify both offers in the
applicable configuration. Source failures, missing mappings and expansion errors
are terminal. No retry selects a weaker representation. A structural declaration
recipe applies only when no exact use binding displaces it.

## Kernel lifecycle

1. Copy the request and snapshot compatible registrations, profile and explicit
   replacements. Check roots, output identity, limits and language/API identities.
2. Admit supplied evidence against the exact source snapshot and unpack the one
   graph its designated evidence publishes. Open the source over that capture.
   Resolve the requested sites to their declaring entities; the source
   classifies them.
3. Walk the graph: open declarations, follow references and applications, and
   judge every occurrence and actual through the source. Select descriptions and
   expand their source-dependent parts. Check that captures are structural
   children. Substitute lexical binders and complete actuals.
4. Close every demanded declaration/helper dependency and check the native
   program. Return an inspectable plan. No target AST has been constructed.
5. On Run, construct fresh ASTs from that fixed native program, check their
   syntax and seal the generation output. No mapping library is invoked.
6. Pass the complete output to Files with placement and the target dependency
   environment. Files constructs module syntax and checks fresh native bindings.

The [abstraction-level context matrix](abstraction-levels.md#agnostic-kernel-and-language-specific-rules)
distinguishes source facts, target intent and realized target syntax. Only native
analysis establishes facts about the constructed program.

## Horos to TypeScript

The following consumer excerpt assumes the engine already has both language
descriptors and their syntax checkers, and that snapshot/evidence/roots belong
to one analyzed Horos capture. Each provider's implemented contract is shown
in the [SPI guide](generation-spi.md).

```go
registry := eng.Generation()
err := compose.RegisterSource[horos.Language, source.Reader](
    registry, from, source.Provider{},
)
if err != nil { return err }

err = compose.RegisterTarget(registry, to, target.Provider{})
if err != nil { return err }

libraries := []generation.MappingLibrary[horos.Language, typescript.Language, source.Reader]{
    text.Library(),
    boolean.Library(),
    list.Library(),
    record.Library(),
    declarations.Library(),
}
for _, library := range libraries {
    if _, err := compose.RegisterMapping(registry, from, to, library); err != nil {
        return err
    }
}

plan, err := eng.PlanGeneration(ctx, compose.Request{
    Input: transform.Input{Source: snapshot, Evidence: []semantic.Result{analysis}},
    To: to, Roots: roots, Output: "pagination-types",
    Profile: rules.RepresentationProfile, Limits: generation.DefaultLimits(),
})
if err != nil { return err }
out, err := plan.Run(ctx)
```

`text`, `boolean`, `list`, `record` and `declarations` are the five packages
under `transform/syntax/horos-to-typescript/mapping/<name>/go`. Each exports
`Library()`, which takes no profile. A library states the structural
representation profile in its own description, and the request names the
profile it wants. No package lists them. The consumer names every library it
registers, as the [installed consumer fixture](../integration/go/testdata/installed/horos_generation_test.go)
does.

Horos supplies interpretation and admission. The mapping library selects native
representations. TypeScript supplies intended native forms and their AST
construction. For Page<Cell<Person>>, the Cell library registers one more exact
description selecting its package export. The kernel and target remain unchanged.
[Cell and reversed Pair qualification](../integration/go/generation_planning_test.go)

Go's target configuration supplies Package, and C++'s supplies Namespace.
Both are recorded as pending placement instances in
[ACX-004](architecture-cases/0004-generation-file-syntax.md): the package or
namespace of a placed declaration belongs to the consumer's placement.
Configuration is stable for a registered target. There is no provider opening
that hands a target builder to a mapping, no procedural operation interface and
no compatibility alias.
