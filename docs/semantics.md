# Agnostic semantic facts

**2 November 2026. Implemented in Go source.** The
[semantic model](../semantics/go/model.go) records facts established by language
providers. The separate [analysis contract](analysis.md) captures a request,
reports coverage and publishes an immutable result. Horos supplies the first
production analyzer; the four native languages still supply syntax tooling only.

The later [Syntax and Files proposal](source.md) names the abstract representation
Semantics and the concrete representation Source, matching AST and CST. It proposes
Semantics[L] as a future replacement for the current Model type, not a second fact
graph. This reference still describes the implemented API; no rename or alias is
implemented by the proposal.

The owner requested this boundary after the initial minimal Go sketch. The
[reviewed expert consultation](../research-docs/0002-semantic-analysis-kernel.md)
records how the sketch was narrowed for implementation. This model has no scope
graph, substitution map, universal type algebra or conditional proof system.
None is required for the current four consumer queries: resolve an occurrence,
inspect an expression, inspect an application's parameter bindings, and enumerate
a declaration's dependencies.

## Ownership

Syntax records what was supplied; analysis records established meanings and
bindings. Bittype owns the representation and publication invariants. A provider
owns lookup, semantic rules, vocabulary, coverage and the interpretation of its
payloads. The kernel validates graph integrity; it cannot prove a provider's
claims or make native meanings equivalent to Horos.

Horos remains the [canonical data-type authority](model.md). Adapting its checked
Program to this graph does not reproduce its descriptor rules. A future native
analyzer may retain different concepts and rules under the same generic interface.
Canonicalization remains an explicit, potentially partial AST transformation.

## Representation

The public construction records are plain Go data. `semantics.New` validates,
bounds and copies them into a `Model` with private storage. Analysis providers
normally publish them through `Request.Publish`, which also checks request
provenance and model requirements.

```go
type EntityID uint32
type TermID uint32

type Key struct { Owner, Name string }

type Entity struct {
    Kind         Key
    Name         string
    Owner        EntityID
    Declarations []ast.NodeRef
    Parameters   []Parameter
    Body         TermID
}
type Parameter struct {
    Entity EntityID
    Sort   Key
}

type Term struct {
    Kind      TermKind // Reference, Application or ProviderExpression
    Target    EntityID
    Arguments []TermID
    Operator  Key
    Encoding  Key
    Data      []byte
    Operands  []Operand
}
type Operand struct {
    Locator []uint32
    Term    TermID
}
```

An entity identifies a module, declaration, parameter or provider-owned intrinsic.
Its name is display metadata. Same-spelled declarations can have different IDs;
one entity can have several declaration sites; an intrinsic can have none.
`Kind` distinguishes provider-owned categories such as aliases and defined types.
`Owner` and `Body` are optional links, not a universal namespace or type-of rule.

IDs are one-based arena positions local to one model. Zero is valid only for
optional links; equal numeric IDs in different models establish nothing. Current
AST sites use their actual tree-qualified, one-based node addresses. Historical
origins are not current declaration sites, module pins or semantic identities.

The three term variants are disjoint:

| Variant | Stored meaning |
| --- | --- |
| Reference | A target entity, including an open parameter or intrinsic |
| Application | A target entity and its complete ordered effective arguments |
| ProviderExpression | An operator key, encoding key, exact copied bytes and located semantic operands |

An unapplied reference and a zero-argument application are distinct. Application
arity must equal the target entity's parameter count. `ParameterBindings` pairs
that parameter list with the argument list; no second substitution map can drift
from the application. A provider must make any defaults explicit before export.
Parameter sorts are keys, so the carrier does not assume that every slot is a type.

Provider expressions retain information that references/applications cannot
express: exact literals, refinements, or native operators. Encoding and operator
keys belong to the provider. Operand locators are distinct provider-defined paths
within that payload; the kernel does not decode them or infer semantic equality
from equal bytes. It checks their term links and storage bounds.

Horos exports canonical encoding of `Expression.Value()` and operands from
`Expression.Operands()`. Its authoritative walker determines type positions.
Literal tuples resembling applications remain bytes, not discovered bindings.
No descriptor-constructor switch is copied into the adapter or kernel.

## Occurrences and dependencies

```go
type Occurrence struct {
    Site    ast.NodeRef
    Role    Key
    Ordinal uint32
}
type Binding struct { At Occurrence; Target EntityID }
type Expression struct { At Occurrence; Term TermID }
type Dependency struct {
    From, To EntityID
    Relation Key
    At       *Occurrence
}
```

Role and ordinal distinguish several semantic occurrences at one syntax node.
`Resolve` returns an established entity binding, never a speculative candidate
set. `TermAt` returns an established term association. Absence is no recorded
answer; it does not prove a failed lookup. A dependency is an explicit directed,
labeled declaration edge. Consumers do not infer dependencies by guessing from
payloads or syntax shape.

Provider keys are owned vocabulary, not executable discovery or registration.
Consumers can inspect the graph generically and interpret particular keys only
under the provider's documented analysis profile.

## Pagination

The [canonical example](horos.md) has this explanatory meaning:

```text
pagination.Page<T> = record { items: list-of(T), next?: text }
people.Person = record { name: text }
people.PeoplePage = pagination.Page<people.Person>
```

Keep distinct entities for both modules, Page, Page's T, Person and PeoplePage.
PeoplePage's body is an application targeting Page. Its argument is itself the
checked zero-argument application of Person. The T in Page's list operand is a
Reference to Page's parameter entity. For `Outer<T> = Page<T>`, the argument
references Outer's T while the application binds Page's T; equal spelling does
not merge the binders.

The application remains inspectable without expansion. Horos's separate
`Specialize` operation can produce a ground descriptor on explicit request.
Analysis does not automatically specialize, choose native representations, or
generate runtime behavior.

## Publication invariants

`Model` accessors return copied mutable records, including nested byte arrays,
operand locators and optional occurrence pointers. Producer mutation after
publication and caller mutation of an accessor result cannot change the model.
Producers must not mutate construction storage concurrently with capture.

Validation checks IDs, term variants, unique occurrence associations, parameter
ownership, effective arity, sites and configured record/edge/byte/depth budgets.
Term containment through arguments and operands must be a DAG. Entity references
and declaration dependencies may cycle; the core imposes no universal semantic
recursion prohibition. A provider may impose its own language's stricter rules.

`Model.Valid()` establishes a constructed graph, not a language-valid program.
Consumers must inspect the enclosing analysis reports and model availability.
The [analysis contract](analysis.md) makes these separate questions explicit.

## Qualification and remaining scope

[Model tests](../semantics/go/model_test.go) qualify identity, ownership, malformed
graphs, current provenance, exact bounds and DAG depth. An independently authored
[foreign-provider fixture](../integration/go/semantic_extension_test.go) exercises
the same API with unfamiliar keys, non-type parameters, aliases, source-less
intrinsics and recursive entities. This proves structural extensibility, not
production semantics for Go or another native language.

Scopes and namespace lookup await a consumer needing them. Assignability,
conversion, inference, conditional judgments and partial checked models require
separate contracts before implementation. No cache key, persistent semantic
identity, automatic normalization or cross-language equivalence is specified.

Function, Stream, Cell and future contracts remain ordinary library declarations.
Forbidden are core contract kinds, semantic name switches, builtin imports or
registrations, behavioral shape fingerprints and generated central dispatch
tables. Adding a contract does not modify this model, the engine, Wire or runtime.
