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
languages; S is the source owner's reading API. The application registers a
MappingLibrary. The registry calls Describe once, outside its lock, and copies
the returned containers. Describe receives no source instance or AST builder.
[Definitions](../contracts/generation/go/description.go)

```go
type MappingLibrary[F, T any, S SourceSession[F]] interface {
    Describe(from, to representation.Language) (MappingDescription[F,T,S], error)
}

type MappingDescription[F, T any, S SourceSession[F]] struct {
    Profile      string
    Bindings     []BindingDescription[F,T,S]
    Declarations []DeclarationRecipe[F,T,S]
    Helpers      []HelperRecipe[T]
}

type BindingDescription[F, T any, S SourceSession[F]] struct {
    ID     string
    Input  Pattern[F,S]
    Output Expression[T]
    Expand func(context.Context, Expansion[F,T,S]) (Expression[T], error)
}

type Pattern[F any, S SourceSession[F]] struct {
    Selector Selector
    Slots    []Slot
    Capture  func(context.Context, S, UseData[F]) ([]Capture[F], error)
}
type Capture[F any] struct {
    Slot Slot
    Use  Use[F]
}
```

An exact selector belongs to the source language. The kernel compares identity,
never the spelling Cell, Page or Data. A fixed rule supplies Output. A
source-dependent rule supplies Expand instead, for example to enumerate record
fields. Capture and Expand run only during planning. They do not recursively
invoke the planner or construct AST nodes.

This is the TypeScript list rule, with error handling abbreviated. The Horos
pattern owns the element operand; the mapping chooses an array; the TypeScript
model owns that native form. [Implemented mapping](../transform/syntax/horos-to-typescript/mapping/go/mapping.go)

```go
element := generation.TypeSlot("E")
pattern, err := rules.Element(from, element)

binding := generation.BindingDescription[horos.Language, ts.Language, source.Reader]{
    ID:     "list/typescript/1",
    Input:  pattern,
    Output: model.Collection(generation.Mapped[ts.Language](element)),
}
```

Slot identity is distinct from its label. The registry rejects undeclared slots.
The planner authenticates every captured occurrence as a direct operand or actual
of the selected parent. A same-named foreign slot cannot stand in for it.
Dynamic expansion receives these implemented methods. The kernel creates and
revokes the request; a retained request is unusable. A rejected capture poisons
that expansion even if its callback swallows the error.

```go
func (r Expansion[F,T,S]) Source() (S, error)
func (r Expansion[F,T,S]) Use() (UseData[F], error)
func (r Expansion[F,T,S]) Mapped(use Use[F]) (Expression[T], error)
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
type DeclarationRecipe[F, T any, S SourceSession[F]] struct {
    ID       string
    Selector Selector
    Expand   func(context.Context, DeclarationRequest[F,T,S]) (DeclarationDescription[T], error)
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
policy. The declaration request exposes Source, Open and Mapped. Its Mapped
operation accepts that declaration's body. Helpers have owner-issued keys and
their own native parameters. They have no synthetic source declaration or
inherited source environment.

## Source and representation environments

One declarative convention supplies both sides of a generic. Bindings is
indexed by the source declaration's authenticated parameters. Arguments is
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

## Source admission

The source integration consumes an already admitted capture and its explicit
evidence. Its existing Source and Catalog contracts remain; narrowing their
traversal responsibilities is a separate proposed change. The new child check
authenticates the relation needed by declarative captures.
[Complete source SPI](../contracts/generation/go/source.go) ·
[Catalog SPI](../contracts/generation/go/mapping.go)

```go
type SourceSession[F any] interface {
    Source[F]
    Catalog[F]
    ResolveRoots(context.Context, []ast.NodeRef) ([]Declaration[F], error)
}
type SourceProvider[F any, S SourceSession[F]] interface {
    Contract() transform.Contract
    Open(context.Context, transform.Request[F]) (S, error)
}
// Source method excerpts
CheckInput(context.Context, semantic.Snapshot, []semantic.Result) error
Admit(context.Context, Environment[F], Use[F]) (UseData[F], error)
CheckChild(context.Context, Environment[F], Use[F], Use[F]) error
```

The kernel invokes these methods. Horos supplies them from its captured
declaration semantics; the independent fixture supplies another source.
Generation never invokes an implicit resolver or analyzer.

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
outside this slice. The separate ASTTransformer routing API remains available
for independent syntax transformations.
