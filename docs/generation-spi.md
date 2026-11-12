# Declarative generation SPI

**Implemented, 8 November 2026.** Mapping libraries describe intended native
types. Planning resolves those descriptions. Construction receives the closed
native program and builds ASTs. The procedural User, Publisher, Definer,
Convention, Plan, Session and MappingProvider interfaces have been removed,
including their tokens and operation tickets. There is no compatibility path.

This design preserves the responsibilities and scope guarantees without
prescribing their former interfaces. The architectural correction is recorded in
[bittype/ACX-005](architecture-cases/0005-generation-interface-design.md).
[Kernel registration](kernel-composition.md) supplies the complete calling example.

## Mapping descriptions

These are implemented signatures in package generation. F and T identify
languages; R is the source owner's reading API. The application registers a
MappingLibrary. The registry calls Describe once, outside its lock, and copies
the returned containers. Describe receives no source instance or AST builder.
[Definitions](../contracts/generation/go/description.go)

```go
type MappingLibrary[F, T, R any] interface {
    Describe(from, to representation.Language) (MappingDescription[F,T,R], error)
}

type MappingDescription[F, T, R any] struct {
    Profile      string
    Bindings     []BindingDescription[F,T,R]
    Declarations []DeclarationRecipe[F,T,R]
    Helpers      []HelperRecipe[T]
}

type BindingDescription[F, T, R any] struct {
    ID     string
    Input  Pattern[F,R]
    Output Expression[T]
    Expand func(context.Context, Expansion[F,T,R]) (Expression[T], error)
}

type Pattern[F, R any] struct {
    Selector Selector
    Slots    []Slot
    Capture  func(context.Context, R, UseView[F]) ([]Capture[F], error)
}
type Capture[F any] struct {
    Slot Slot
    Use  Use[F]
}
```

An exact selector belongs to the source language. Mapping selectors are frozen
in the capture's source language, and every selector the source returns must
have the same language and profile. Otherwise planning fails with `ErrScope`
before any exact or default rule is chosen. Within that language the kernel
compares keys exactly, never the spelling Cell, Page or Data. A fixed rule supplies Output. A
source-dependent rule supplies Expand instead, for example to enumerate record
fields. Capture and Expand run only during planning. They do not recursively
invoke the planner or construct AST nodes.

This is the TypeScript list rule, with error handling abbreviated. The Horos
pattern owns the element operand; the mapping chooses an array; the TypeScript
model owns that native form. `patterns` is the package
`transform/syntax/horos/mapping/pattern/go`. [Implemented mapping](../transform/syntax/horos-to-typescript/mapping/list/go/list.go)

```go
element := generation.TypeSlot("E")
pattern, err := patterns.Element(from, element)

binding := generation.BindingDescription[horos.Language, ts.Language, source.Reader]{
    ID:     "list/typescript/1",
    Input:  pattern,
    Output: model.Collection(generation.Mapped[ts.Language](element)),
}
```

Slot identity is distinct from its label. The registry rejects undeclared slots.
The planner accepts a captured occurrence only if the kernel minted it as a
direct operand or actual of the selected occurrence in this run. A same-named
foreign slot, the occurrence itself, or a handle kept from another run cannot
stand in for it.
Dynamic expansion receives these implemented methods. The kernel creates and
revokes the request; a retained request is unusable. A rejected capture poisons
that expansion even if its callback swallows the error.

```go
func (r Expansion[F,T,R]) Reader() (R, error)
func (r Expansion[F,T,R]) Use() (UseView[F], error)
func (r Expansion[F,T,R]) Mapped(use Use[F]) (Expression[T], error)
```

Mapped records a dependency. It does not execute another mapping. Native
expressions contain target-owned forms, immutable scalar metadata and child
expressions. The shared interpreter understands substitution and references,
not the meanings of these forms. There is no AST kind, source string or
arbitrary payload field in an expression.
[Expression contract](../contracts/generation/go/expression.go)

## Declarations, publication and helpers

The application requests source roots. A declaration recipe describes each
owned declaration. The planner records the corresponding target identity and
the constructor emits its AST site into Output.Exposed. An external use rule
alone cannot satisfy a publication request.

```go
type DeclarationRecipe[F, T, R any] struct {
    ID       string
    Selector Selector
    Expand   func(context.Context, DeclarationRequest[F,T,R]) (DeclarationDescription[T], error)
}
type DeclarationDescription[T any] struct {
    Name    string
    Form    Form[T]
    Origin  ast.NodeRef
    Generic GenericConvention[T]
    Body    Expression[T]
}
type HelperRecipe[T any] struct {
    Key         HelperKey
    Declaration NativeDeclaration[T]
}
```

A zero recipe selector explicitly supplies a structural declaration policy.
An exact binding takes precedence, and its failure never falls back to that
policy. The declaration request exposes Reader, Declaration (the kernel's
DeclarationView) and Mapped. Its Mapped operation accepts that declaration's
body occurrence only. Helpers have owner-issued keys and
their own native parameters. They have no synthetic source declaration or
inherited source environment.

## Source and representation environments

One declarative convention supplies both sides of a generic. Bindings is
indexed by the source declaration's real binders, in parameter order. Arguments is
indexed by the native parameters. The names do not establish either relation.

```go
type GenericConvention[T any] struct {
    Parameters []ParameterDescription
    Bindings   []Expression[T]
    Arguments  []Expression[T]
}
type ParameterDescription struct {
    Name        string
    Origin      ast.NodeRef
    RequiredUse bool
}
func Parameter[T any](index int) Expression[T]
func Argument[T any](index int) Expression[T]
```

For Page<T>, Bindings contains Parameter(0), and Arguments contains Argument(0).
Reordering native parameters changes both arrays. An extra native parameter
can have an argument expression such as Collection(Argument(0)). The independent
Go fixture uses this to generate Page<T, Storage> with Storage supplied as []T.

The planner substitutes complete mapped actuals into Arguments. Thus the
additional argument for Page<Page<Person>> is []Page<Person, []Person>, not
[]Person. Every source actual still requires admission, including an actual
omitted by the native convention. Only demanded native actuals need mappings.
RequiredUse expresses the selected recipe's requirement that a target parameter
occur in its body; the default Horos structural recipe requires this.

Parameter references carry the owning planned declaration. Native construction
checks that they occur in that scope. A helper cannot capture a caller's
parameter. These guarantees replace mutable operation tickets and payload
retagging. Generation is not a proof of source/native type equivalence.
[Independent fixture](../integration/go/generation_independent_test.go) ·
[Verification](verification.md)

## Semantic traversal and source interpretation

The kernel admits the request and walks the semantic graph itself. Its source
contract designates exactly one evidence requirement that publishes a graph.
`Plan` admits the input's snapshot and evidence against the provider's
contract, unpacks that result's graph as `Semantics[F]`, and opens the provider
over exactly that capture. No caller supplies an opened source or a graph.
Handles resolve only against that capture. Interpreting that capture is the provider's responsibility: the
kernel cannot check where a provider's answers come from.
`Plan` checks the mapping descriptions and limits before it admits the
evidence. The registry checks the request's evidence count and root sites
first, then the descriptions, then calls `Plan`. `RegisterSource` rejects an
unusable contract early; each run admits against the provider's `Contract` at
that time and checks it again. The kernel reads only structure:

- a root site's declaring entity;
- a declaration's body and real binders;
- term variants, reference targets, application callees and arguments, paired
  with real binders by `DirectParameterBindings`;
- provider-expression operands, with locators passed through unread.

It never reads entity kinds, names, operators, payloads, locator meanings or
facts. [Contract](../contracts/generation/go/source.go) ·
[Walk and handles](../contracts/generation/go/occurrence.go)

```go
type PlanningInput[F, T, R any] struct {
    From      representation.Language
    Input     transform.Input
    Source    SourceProvider[F, R] // Plan admits Input against it and opens it
    Roots     []ast.NodeRef
    Mappings  []SelectedDescription[F, T, R]
    Execution Execution[T]
}
type SourceProvider[F, R any] interface {
    Contract() transform.Contract // exactly one requirement publishes the graph
    Open(context.Context, Captured[F]) (Source[F, R], error)
}
type Source[F, R any] struct {
    Interpretation Interpretation[F]
    Reader         R // the reading API that this source's mapping libraries use
}
type Interpretation[F any] interface {
    Declaration(context.Context, Declaration[F]) (Selector, bool, error)
    Expression(context.Context, Use[F]) (Selector, error)
    JudgeDeclaration(context.Context, Declaration[F]) (Judgment, error)
    JudgeUse(context.Context, Use[F]) (Judgment, error)
}
type Judgment struct {
    Verdict    Verdict // Admitted, Unsupported, Rejected or Undetermined
    Diagnostic error
}
```

The source owns three things:

- **Classification:** whether an entity is a declaration it generates, and that
  declaration's stable selector.
- **Descriptor selectors:** the selector for each provider expression.
- **Judgments:** whether a declaration, or a whole occurrence, is admissible under
  its declaration's assumptions. A judgment uses only the captured evidence.

The kernel maps each verdict to one failure meaning:

| Verdict | Kernel failure |
| --- | --- |
| Unsupported | `ErrUnsupported` |
| Rejected | `ErrRejected` |
| Undetermined | `semantic.ErrEvidenceUnavailable` |
| zero value | `ErrInvalid` |

An uncovered fact is therefore never reported as a rejection.

**Handles.** The kernel issues `Declaration[F]`, `Binder[F]` and `Use[F]` for one
run. A `Use[F]` is an *occurrence*: a graph term reached under one declaration
along one structural path. A term shared by two bodies is two occurrences. A
holder resolves a handle with `In(captured)`, which succeeds only for the
capture that the handle's run walks. Two captures that share a graph but differ
in evidence are different.

**Scope.** A reference to a parameter of the declaration being planned is a bound
use. A reference to another declaration's parameter is `ErrScope`. Two kinds of
form are unsupported at their site, because the structure cannot form a complete
native application:

- an application without a direct declaration callee and complete arguments;
- a generic declaration used without arguments.

These are not source judgments.

**Whole-use admission.** Before any binding or convention represents an
occurrence, the kernel admits the occurrence and its whole subtree. Every
occurrence in the subtree needs a structural reading in its declaration's scope
and an Admitted judgment. That covers actuals the native convention omits,
content nested inside them, and operands no mapping captures. A generated
declaration's body is admitted before its recipe runs, even if the recipe does
not represent it. Representation stays demand-driven; admission does not.
Admissions are memoized per occurrence.

The judgment of a use covers that use. A callee declaration is classified only
when a use of it is represented, and judged only when it is generated. A source
that does not admit some uses of an entity must say so in `JudgeUse`.
Declarations are memoized per entity, and recursion is detected there.

A source diagnostic only explains a verdict. It cannot add an error class:
`errors.Is` sees only the verdict's class, while `errors.As` still reaches the
diagnostic and any site it carries.

A kernel failure of an occurrence is reported at its term's syntax site when the
graph records exactly one. A term recorded at several sites is reported at the
declaration the occurrence is in. A use whose callee the source does not
generate, or whose callee selector is outside the source language, is a failure
of that use, not of the callee.

A capture that is not a child the kernel minted, or a declaration recipe mapping
anything but its body, fails as `ErrScope` with its own message.

Mapping callbacks receive the kernel's views and the source reader, never the
graph or the interpretation:

```go
type UseView[F any] struct {
    Use      Use[F]
    Kind     UseKind        // BoundUse, DeclarationUse or ExpressionUse
    Bound    Binder[F]
    Callee   Declaration[F]
    Actuals  []Actual[F]    // {Binder, Use}, in the callee's parameter order
    Operands []Operand[F]   // {Locator, Use}
    Body     bool           // the occurrence is its declaration's whole body
}
func (v UseView[F]) Actual(Binder[F]) (Use[F], bool)
func (v UseView[F]) Operand(locator ...uint32) (Use[F], bool)
type DeclarationView[F any] struct {
    Declaration Declaration[F]
    Formals     []Binder[F]
    Body        Use[F]
}
```

The horos reader interprets descriptors, record fields, declaration identity and
binder names over these handles. Its field and element readers select operand
occurrences from the kernel's view by horos locator. The independent fixture
publishes its own graph through its own analysis and supplies only its own
interpretation, reader and library. Both drive the same kernel walk. Generation
never invokes an implicit resolver or analyzer.
[horos source](../transform/syntax/horos/source/go/provider.go) ·
[Independent source](../integration/go/testdata/generationfixture/fixture.go) ·
[Traversal qualification](../integration/go/generation_traversal_test.go)

## Target construction

A target integration implements this interface. The planner invokes its checks
against intended native forms and the completed program. Run invokes Build.
Build receives no source reader, mapping callback or selector. It also receives no
source-binder or application templates. Those remain in the planner's convention
state. Its declarations are target-only:

```go
type NativeDeclaration[T any] struct {
    Name       string
    Form       Form[T]
    Origin     ast.NodeRef
    Parameters []NativeParameter
    Body       Expression[T]
}
type NativeParameter struct {
    Name   string
    Origin ast.NodeRef
}
```

[Target contract](../contracts/generation/go/contracts.go)

```go
type Target[T any] interface {
    CheckExpression(context.Context, Expression[T]) error
    CheckDeclaration(context.Context, NativeDeclaration[T]) error
    CheckProgram(context.Context, Program[T]) error
    Build(context.Context, BuildInput[T]) (Assembly, error)
}
type BuildInput[T any] struct {
    Program  Program[T]
    Language representation.Language
    Output   ast.TreeID
    Limits   Limits
}
```

The four integrations delegate AST scheduling, fresh containment, scoped symbols,
relocation and publication association to the shared constructor. Their native
lowering methods supply grammar. Their model packages describe native intent
without constructing syntax. Native symbols and fragments belong to one run;
a foreign symbol cannot bind merely because its name or numeric position matches.
[Constructor](../transform/syntax/shared/target/go/compile.go) ·
[TypeScript target](../transform/syntax/typescript/target/go/provider.go)

Targets are trusted native semantic compilers. The kernel checks captures,
scopes, limits and resulting syntax; it does not prove an arbitrary target
implementation's semantics. Native analysis and the verification SPI remain
separate evidence.

## Plan and output

The engine exposes PlanGeneration and Generate. The latter plans and runs once.
Only successful planning returns a usable GenerationPlan. It retains a closed,
immutable native program and target configuration. Steps returns copies.

```go
func (e *Engine) PlanGeneration(context.Context, compose.Request) (generation.GenerationPlan, error)
func (e *Engine) Generate(context.Context, compose.Request) (generation.Output, error)
func (p GenerationPlan) Steps() []Step
func (p GenerationPlan) Run(context.Context) (Output, error)
```

Each Run creates fresh construction state and a fresh checked snapshot. There is
no mapping callback during Run, so a mutable array selected during planning
cannot become a readonly array through a second mapping operation. Changing
registered libraries affects later plans. Failure, cancellation or a limit
returns no usable output. Recursion is explicitly unsupported by the current
native construction profiles.

Output carries checked syntax, requested public sites and external destinations
into Files. Files chooses placement and module syntax, then checks the actual
native bindings. [Output handoff](typescript-type-imports.md) ·
[Planning tests](../integration/go/generation_planning_test.go)

Native kind enforcement, executable code, adapters and runtime behavior remain
outside this slice. [GEN-02](generation-gaps.md#gen-02) records the absent native
bound contract and the work required before an enforcing profile can be claimed.
The separate ASTTransformer routing API remains available
for independent syntax transformations.
