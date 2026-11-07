# Planning and running type generation

**Draft proposal · 7 November 2026**

A caller should be able to ask Bittype to generate `pagination.Page` for
TypeScript, inspect which mappings the engine chose, and run that plan. The
project configures its generation libraries and representation preferences once.
The engine assembles the pipeline. During planning, mapping libraries choose
the target types and bind their uses. A run constructs those choices as target
ASTs. A caller can therefore inspect what `list<T>` will become before the
TypeScript builder allocates an array node.

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

## Semantic mapping and AST construction

For a source language A and target language B, the proposed `PlanGeneration`
would perform **semantic mapping**. Its inputs are the captured A declarations,
the facts established about them, the requested roots and the configured
mapping libraries. Its inspectable result describes the intended B declarations,
type expressions, generic binders, reference destinations and declarations to
expose. `Run` would perform **AST construction**. The B builder realizes those
choices as native declaration nodes and retains the destination of each named use.

| Stage | Input | Result |
| --- | --- | --- |
| Read A | The application supplies source files and the selected dependencies. | Syntax parsing yields CSTs with `AST[A]` views. The A resolver binds occurrences in their scopes. The A analyzer uses that resolution to publish `Semantics[A]` for the captured program. |
| Map to B | Captured A sites and facts, requested roots and registered mappings. | The planner selects rules and records the B declarations, expressions and bound uses they prescribe. |
| Construct B | A run of the plan with the selected target builder. | A syntax-checked `AST[B]` set realizes the declarations. Bound uses and the requested public surface must survive this step. |
| Organize B | The intended bindings, requested public surface, selected dependencies and one consumer placement. | The B Files formatter applies that placement and realizes the appropriate imports, exports and local names in each unit. |
| Check and write B | The organized units in their target environment. | The B resolver and analyzer establish `Semantics[B]` for this capture. Files checks the actual bindings and publishes `Source[B]`. The Files printer returns paths and text for the application to write. |

The plan describes what the selected mappings intend. Only analysis of the
realized B program establishes `Semantics[B]`. Today generation consumes an
analyzed source and returns a syntax-checked target AST set. The separation of
mapping decisions into an inspectable plan is proposed here.
[Current Files input and proposed handoff](#where-the-plan-stops) ·
[Generation design](generation.md#semantic-mapping-and-target-ast-construction) ·
[Current generation contract](semantic-generation.md#explicit-semantic-input) ·
[Syntax and Files roles](source.md#symmetric-role-contracts).

### Pagination on read and write

Suppose the caller requests `Page`, `Person` and `Cells` as public TypeScript
declarations. The application supplies three Horos files with `people.horos`
as the entry. Each pin placeholder stands for the 64 hexadecimal digits of its
module's pin. The `Cell` declaration follows the independent library fixture.
[Horos Files reading](source.md#horos-composition-and-pagination) ·
[Cell mapping fixture](../integration/go/horos_generation_composition_test.go).

```horos
// pagination.horos
module example.pagination;
type Page<T> = { items: list<T>; next?: text; };

// cells.horos
module example.cells;
type Cell<V> = wire<V>;

// people.horos
module example.people;
import pages from example.pagination pin "sha256:<pagination-pin>";
import cells from example.cells pin "sha256:<cells-pin>";
type Person = { name: text; };
type Cells = pages.Page<cells.Cell<Person>>;
```

**Read Horos.** Syntax parsing retains the written names and nested applications.
The resolver follows the two import pins and resolves `Person` in
`example.people`. Analysis establishes that the inner application binds
`Cell.V` to that `Person`, and the outer one binds `Page.T` to the whole
`Cell<Person>` application. The `T` in `items` still belongs to the open
`Page` declaration. These facts concern the Horos input and its scopes.
[Pagination analysis](analysis.md#pagination-from-input-to-queries) ·
[Horos Files fixture](../integration/go/files_test.go).

**Map to TypeScript.** The structural mappings select an interface for each of
`Page` and `Person` and a type alias for `Cells`. The list mapping chooses an
array whose element is the target parameter of `Page`. The Cell library's
mapping matches the pinned Horos `Cell` declaration, selects the public type
export `Cell` from `@bitspark/cell`, and uses the mapped `Person` as its
argument. The outer application uses that entire result as the argument to
the generated `Page`. **Semantic mapping ends here.** The representation,
each reference destination and the three declarations to expose are fixed.

**Construct the TypeScript AST.** The builder now has enough information to
construct the interfaces, type parameter, properties and alias. The outline
shows containment and the intended destination of each reference. The
external package identity is an association on the use.

```text
interface Page
  type-parameter T
  property items → array-type → type-reference to Page.T
  optional property next → keyword-type string
interface Person
  property name → keyword-type string
type-alias Cells
  type-application
    constructor → type-reference to generated Page
    argument    → type-application
      constructor → type-reference to Cell exported by @bitspark/cell
      argument    → type-reference to generated Person
```

A syntax check admits these node kinds, roles and containment. The builder
retains the associations through construction so Files can later compare
each reference with the target the compiler finds. [TypeScript AST vocabulary](syntax.md#typescript).

**Organize the output.** Run once. Put `Page` and `Cells` in `pages/page.ts`
and the public `Person` in `model/person.ts`. Given the selected Cell package,
TypeScript Files produces the following module.

```typescript
import type { Person } from "../model/person.js";
import type { Cell } from "@bitspark/cell";

export interface Page<T> { items: T[]; next?: string; }
export type Cells = Page<Cell<Person>>;
```

Change only the placement to put `Person` beside `Page`. The pages module now
also contains `export interface Person { name: string; }`. Its relative import
disappears. The package import and the alias body stay the same.

**Check the bindings before writing.** Before publishing either organized `Source`,
resolve and analyze its units in the selected TypeScript package environment.
Compare the `Page` and `Person` uses with their declarations at their new
sites, the `T` use with the binder on that `Page`, and the `Cell` use with the
package export. Each capture has its own entity IDs. The comparison follows
the declaration into its new unit or identifies the selected package and export.
Reading TypeScript that spells the import as
`import type { Cell as C } from "@bitspark/cell"` must recover that same export
from the `C` in `Page<C<Person>>`. If a use binds to a different same-named
type, or the package export is unavailable, Files must reject the result.
The printer receives the checked `Source` and the application chooses where
to write it. [TypeScript placement and binding checks](typescript-type-imports.md#from-generation-to-files) ·
[Placement and alias fixtures](../integration/go/typescript_placement_test.go).

### Where the plan stops

`Publish` records which declarations to expose. Each named use retains the
identity of its owned declaration, lexical parameter or external public symbol.
A change of placement uses the same plan. Selecting a different external type
for Cell changes a mapping decision and requires a new plan.

The proposal still needs a way to carry the bound uses and publication intent
from a run into target semantic admission. The proposed `ASTSet` return type
alone leaves that association unspecified. Today's TypeScript target emits
package imports and export modifiers into a logical unit, which is resolved
and analyzed to supply the existing Files formatter's input graph, snapshot
and receipt. After placement the formatter obtains a fresh analysis and checks
the relocated bindings before it publishes `Source`. Under the proposed split,
the intent from the plan must reach target admission directly. Analysis of a
logical unit, where available, establishes facts only for that unit. The
organized units still need their own binding checks.
[Generation and Files boundary](source.md#from-generated-declarations-to-files) ·
[Current TypeScript behavior](typescript-type-imports.md) ·
[Open architectural case](architecture-cases/0004-generation-file-syntax.md).

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
of that exact occurrence. During planning, `ArrayOf` selects an array of that
mapped element. During construction, the TypeScript builder creates the array
node with the realized child in its element position. The same expression drives
both stages. For `list<bool>` this yields
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
   readonly-list override. Compare the planned element of `items` with the
   reference to the target parameter in the constructed AST.
3. Migrate the remaining providers and callers, including helper, external and
   declaration operations for the four targets. Replace the per-request mapping
   list and superseded procedural rules on the normal generation path. Keep
   `Generate` as the one-call entry. This changes the public provider contract
   and the requests made by installed consumers.
4. Carry bound external uses and publication intent across the generation/Files
   handoff. Remove early import/export emission from TypeScript generation.
   Format one run of `Cells` with both pagination placements. Check each
   organized use against its intended declaration and compile both layouts
   against the selected package. An inaccessible binding must fail even if a
   different declaration shares its name. Apply the same boundary when migrating
   other targets that already emit includes or module visibility.

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
