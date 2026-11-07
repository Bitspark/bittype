# Planning and running type generation

**Draft proposal · 7 November 2026**

A caller should be able to ask Bittype to generate `pagination.Page` for
TypeScript, inspect which mappings the engine chose, and run that plan. The
project configures its generation libraries and representation preferences once.
The engine assembles the pipeline. Mapping libraries declare how source types
become target types, including how a generic element is mapped.

## What exists today

| On main | Proposed change |
| --- | --- |
| `Engine.Plan(from, to)` finds a language route through registered AST transformers. | Add `PlanGeneration` for decisions about particular source declarations. |
| `Registry.Transformer(request)` captures an explicit provider selection and ties it to the exact source snapshot. Its execution opens the providers and collects their bindings. | The planner resolves the reachable mapping choices before returning a runnable plan. |
| A mapping provider's `Open` can supply exact bindings, a structural declaration policy, and helper plans. The TypeScript list mapping calls `Child` on its element and builds an array node in Go. | Advertise bindings at registration and express the list rule as an executable input pattern and target expression. |
| `Engine.Generate` takes the analyzed input, roots, target, output identity, provider handles and limits. It returns a syntax-checked AST set. | Use configured providers and default limits. Keep `Generate` as the convenience call that plans and then runs. |
| The TypeScript target already emits package imports and marks declarations exported in its logical AST. Its Files formatter rebuilds imports after placement. | Plan bound target uses and publication intent. Files formatting is the first stage that constructs concrete import/export syntax for generated declarations. |

The existing `generation.Plan[F,T]` describes the operations for one mapping or
declaration. A `GenerationPlan` in this proposal covers a complete request and
the dependencies of its roots. The engine already owns assembly, following the
correction recorded in ACX-002. [Current generation lifecycle](kernel-composition.md#kernel-lifecycle) ·
[Language route planner](../engine/go/path.go) ·
[Generation registry](../transform/syntax/shared/generation/go/registry.go) ·
[Generation assembly case](architecture-cases/0002-generation-assembly.md).

## Caller interface

The first public operation would plan against an analyzed source and explicit
roots. For example, the module below has one declaration to publish.

```horos
module pagination;
type Page<T> = { items: list<T>; };
```

The proposed public surface has three operations. Only the engine can issue a
runnable plan. Inspection exposes decisions that callers can read.

```go
func (e *Engine) PlanGeneration(ctx context.Context, r generation.Request) (GenerationPlan, error)
func (p GenerationPlan) Steps() []GenerationStep
func (p GenerationPlan) Run(ctx context.Context) (representation.ASTSet, error)
```

The call uses the declared target language and its profile. The source
language comes from the snapshot. A file-based frontend would parse, resolve
and analyze the supplied module and imports, then resolve `pagination.Page` to
the exact declaration site in that snapshot.

```go
plan, err := eng.PlanGeneration(ctx, generation.Request{
    Input: transform.Input{
        Source:   snapshot,
        Evidence: []semantic.Result{horosAnalysis},
    },
    Roots:  []ast.NodeRef{pageSite},
    To:     typescriptLanguage,
    Output: "pagination-types",
})
if err != nil {
    return representation.ASTSet{}, err
}
return plan.Run(ctx)
```

`Output` identifies the generated tree independently of any file. `plan.Steps()`
would return a read-only view of the requested publications, selected bindings,
scoped operands and dependency edges.
`Run` would return a `representation.ASTSet`. The same plan could be run again
against the captured input with a fresh execution context.

Supplying just the two language descriptors would allow the engine to assemble
a reusable pipeline of compatible providers and configured preferences. Selecting
the rule for a particular field or imported declaration requires the source
snapshot and roots. The first API exposes the source-bound plan. An endpoint-only
pipeline could later bind separate requests through that same planning operation.

### Where the plan stops

`PlanGeneration` chooses representations and binds uses; `Run` constructs
target declarations. Neither operation assigns output paths or imports. The
information that must survive the run is the **intended target of each use**
and **which declarations to expose**. `Publish` expresses the latter; it does
not decide whether a particular file contains an `export` modifier or statement.
Generic parameters keep their lexical identity. A reference to an external
public symbol keeps its dependency and export identity without a fictitious
local declaration or a file inside that dependency.

Planning `Page<T>` alone reaches its list and binder. Planning the separate
alias `Cells = Page<Cell<Person>>` also selects the owned `Person` and the
library mapping for the external `@bitspark/cell` export `Cell`. Inspection
can explain all three references without predicting any import statements.
Run the plan once. Place `Person` in `model/person.ts` and `Cells` beside
`Page` in `pages/page.ts`: TypeScript Files formatting needs both a relative
type import and a package type import. Move only `Person` into `pages/page.ts`
and format the same generation result again: only the package import remains.
Changing which external type represents Cell, by contrast, changes the mapping
selection and requires a new plan.

The consumer supplies the placement. The Files formatter uses it with the
target bindings and admitted semantic environment to choose per-unit spelling
and aliases and express the requested visibility. The printer emits the completed
`Source`. A syntax-only render can show a declaration but cannot establish a
module's public surface or imports. The proposal still needs a way to carry
external bindings and publication intent from a run into target semantic
admission without manufacturing concrete import nodes. Returning an `ASTSet`
alone does not specify that association. This is a handoff to the existing
[Syntax and Files stages](source.md#from-generated-declarations-to-files), not
a second organization model. [Current TypeScript behavior](typescript-type-imports.md)
· [Open architectural case](architecture-cases/0004-generation-file-syntax.md).

## What a mapper declares

Registration would call `Describe` once and capture the declared endpoints,
generation profile and offers. An exact offer contains a provider-local binding
ID, a source-owned pattern and a target-owned expression. The pattern supplies
an exact selector for lookup. The expression specifies the output and refers to
the recursively mapped children. Go signatures below illustrate the proposal.

```go
type MappingProvider[F, T any, S SourceSession[F], N Native[F, T]] interface {
    Describe() (MappingDescription[F, T], error)
    Open(context.Context, MappingRequest[S, N]) (Contribution[F, T], error)
}

type MappingDescription[F, T any] struct {
    From, To          representation.Language
    GenerationProfile string
    Bindings          []BindingDescription[F, T]
    StructuralID      string   // Optional structural declaration policy.
    HelperKeys        []string // Helpers this provider can supply.
}

type BindingDescription[F, T any] struct {
    ID     string
    Input  SourcePattern[F]
    Output TargetExpression[T]
}
```

For the list mapper the entire type rule is **`Horos.list<E> → TypeScript.Array(mapped(E))`**.
The same slot occurs on both sides of the library's declaration.

```go
element := rule.TypeSlot("E")
listArray := gen.BindingDescription[horos.Language, ts.Language]{
    ID:     "list/typescript/1",
    Input:  horosPattern.ListOf(element.Capture()),
    Output: tsType.ArrayOf(element.Mapped()),
}
```

`E` is a slot owned by this rule. `ListOf` identifies the Horos constructor and
locates its element occurrence. `Mapped` requests the selected target representation
of that exact occurrence. `ArrayOf` describes a TypeScript array node with the
result in its element position. Running this expression drives the existing
child-generation and target-builder operations. For `list<bool>` this yields
`boolean[]`. For `list<list<bool>>` it yields `boolean[][]`. Reusing a slot from
another rule or referencing an uncaptured slot is an invalid declaration.
[Current list mapper](../transform/syntax/horos-to-typescript/mapping/go/mapping.go) ·
[Target fragment builder](../transform/syntax/shared/target/go/builder.go).

The source library owns selector identity and operand interpretation. The target
library owns native forms and syntax. The mapping owner chooses the representation.
The kernel matches opaque selectors and manages scoped captures and substitution.
A `Cell<T>` library can use the same interface with an exact module pin and
declaration identity and choose an external target type. A project can replace
the list binding with `ReadonlyArray<mapped(E)>` in its profile. The preference
identifies both the expected and replacement bindings by provider and binding ID.
An unresolved exact conflict at a demanded occurrence is a planning error.
[Source selectors](../transform/syntax/horos/source/go/identity.go) ·
[Existing readonly and Cell examples](../integration/go/horos_generation_composition_test.go).

## What the plan contains

Planning snapshots the configured provider registrations, their compatible
language and Go API identities, generation profile, replacements, output
identity and limits. It admits the supplied semantic receipts against the
exact source snapshot, authenticates the requested declaration sites, and
walks dependencies starting at those sites. An eligible structural policy can
handle a declaration with no exact binding. An exact binding that fails ends
the request. A missing demanded representation or an ambiguous choice is
reported at the requesting source occurrence.

An inspection of the example would expose these decisions.

```text
Publish pagination.Page       TypeScript interface Page<T>
  Introduce Page.T            target parameter bound to this declaration
  Define field items          Horos list whose element refers to Page.T
    Select list/typescript/1  Array(mapped(E))
    Bind E                    this occurrence of Page.T
    Map E                     use Page's target parameter
```

The plan requests publication of `Page`. If the consumer places it in a
TypeScript module as a named export, Files formatting and printing produce:

```typescript
export interface Page<T> {
  items: T[];
}
```

A plan for an alias such as `Nested = Page<Page<Person>>` would include the
required named declarations and both applications. Each use refers to the same
open `Page` declaration. The inner application's result supplies the argument to
the outer one. Rule slots, source generic binders and target parameters each
keep their own identities. The source still admits every actual argument against
its original contract, including an argument the target representation omits.
[Generation operations and binder scopes](generation-spi.md#source-and-representation-environments).

A complete bound plan fixes each reachable representation choice and records
the required child, declaration and helper dependencies. It may store shared
dependencies once and refer to them from several steps. Planning admits only
owner-authenticated references from the captured input, respects the existing
resource bounds, and checks cyclic dependencies against the selected target's
supported domain. Inspection returns copies or immutable views. Edits to the
registry or a different source snapshot require a new plan to take effect.

### Making dependencies visible

Today callbacks such as `Child` and `RequireHelper` can discover work while a
mapping runs. For fixed-shape rules, `mapped(E)` makes the child visible in the
declaration. A record with a variable number of fields needs a source-owned
projection to enumerate the fields of the actual declaration. A Go optional-field
mapping can determine at that point whether it needs the registered presence/value
helper. A binding to an external type can declare the external target reference
and its required children while keeping the source body opaque.

For a source-dependent part of a rule the owner would supply an `Expand`
operation on that part of the registered expression. For example, a record
mapper can declare an object type with repeated members. Its expansion would
read the fields of the captured record and fill those member positions with
references to the mapped field occurrences. The target object form is declared
once. A fixed expression such as `ArrayOf(element.Mapped())` expands directly
from its declaration.

Expansion would use a kernel-issued request scoped to the selected rule and
source occurrence. It would return a target expression with references to its
child uses, owned or external declarations and helper keys. The kernel would
check that each child came from the requested parent, authenticate declaration
and binder identities, and resolve helper keys to their configured owners.
Structural policies and helper-only providers follow the same rules. An offer
that cannot expose its reachable dependencies would report unsupported planning
when selected. Execution may materialize a planned helper on demand.
[Current definition scheduler](../transform/syntax/shared/generation/go/definitions.go) ·
[Go helper use](../transform/syntax/horos-to-golang/mapping/go/mapping.go).

## Execution guarantees

The plan retains immutable descriptions, source occurrences, provider-qualified
bindings and stable provider configuration. Each provider must capture its
configuration at registration and register a new instance for a change. The
kernel snapshots the registry under its lock and invokes owner code after
releasing it. Only successful planning issues a runnable handle.

Each `Run` opens fresh sessions against the captured inputs. The executor binds
the plan's symbolic references to that run's source and target handles and
checks the opened contributions against the selected offers. Construction
tickets, association tokens and native fragments belong to that run.

The executor uses the existing `Publish`, `Introduce`, `Use`, `Apply` and `Define`
operations. It preserves lexical environments, authenticates returned products,
allocates native symbols and nodes, and admits the finished target tree. An
unplanned dependency, foreign handle, provider failure, cancellation, exhausted
limit or invalid target tree fails that run as a whole. The same plan remains
available for another run. A successful run guarantees no particular file
organization. Files formatting must preserve the target bindings and requested
public surface or reject the arrangement; it cannot ask the generator to select
a different type. Cross-language proofs use the verification contract.
[Guarded operations](generation-spi.md#five-operations) ·
[Engine generation entry](../engine/go/generation.go) ·
[TypeScript file placement](typescript-type-imports.md) · [Verification](verification.md).

## Migration and qualification

1. Introduce typed source captures and target expressions over the existing
   generation SPI. Migrate the TypeScript list rule so its declared expression
   produces the output. Exercise open, concrete and nested elements.
2. Extend registration with immutable offers and move selection into
   `PlanGeneration`. Implement `Steps` and `Run` for a named `Page` root,
   including structural field traversal, binder identity, and an explicit
   readonly-list override.
3. Migrate the remaining providers and callers, including helper, external and
   declaration operations for the four targets. Replace the per-request mapping
   list and superseded procedural rules on the normal generation path. Keep
   `Generate` as the one-call entry. This changes the public provider contract
   and the requests made by installed consumers.
4. Carry bound external uses and publication intent across the generation/Files
   handoff. Remove early import/export emission from TypeScript generation.
   Qualify one run of `Cells` against both placements above: inspect the same
   declaration bindings before formatting, then compile both emitted layouts
   against the selected package. An inaccessible binding must fail, never
   resolve by coincidence to another declaration of the same name. Apply the
   same boundary when migrating other targets that already emit includes or
   module visibility.

Qualification should compare planned and one-call output, exercise fresh runs
of one plan after a registry change, and reject ambiguous, stale or foreign
bindings, cross-rule slots and undeclared helper demands. An independent library
should compose a materially different rule, such as a two-operand mapping with
reversed target arguments, inside the existing list and record mappings through
its own declaration and consumer configuration. The installed-consumer fixture
and native compiler tests can exercise the migrated public path. These are
proposed tests for the implementation.
[Kernel fixtures](../integration/go/kernel_generation_test.go) ·
[Pagination fixtures](../integration/go/horos_generation_pagination_test.go) ·
[Installed consumer](../integration/go/testdata/installed/horos_generation_test.go).
