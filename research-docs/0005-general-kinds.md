# Research: general structural kinds for Horos and Bittype

**ID:** 0005
**Date:** 4 November 2026
**Status:** reviewed
**Owner:** codex-hydrated-type-design
**Issue:** https://github.com/Bitspark/bittype/issues/59
**Reviewed:** https://github.com/Bitspark/bittype/issues/59#issuecomment-6636131079

## Question

What general, implementable kind system should replace the current closed
Type/Data interface and dedicated Data proof system? The owner explicitly asks:

> Then why don't we make scalar, record, list, etc. proper kinds and define Data
> as their union? I want a general kind system, not special handling for Data and Type.

The owner then requested this expert consultation. This is a substantive
correction to research 0004, not a request to rename its two enum members or add
more built-in cases. We need kinds that users can define, reference, parameterize
and combine, with Data expressible through those ordinary mechanisms. Structural
families such as Scalar, Record and List must carry useful constraints, including
constraints on nested types. Constructor parameters must continue to work.

Please answer the five questions below in at most 6,000 words. Specify the
meaning, exact core grammar, checking rules and counterexamples before proposing
Go APIs. Recommend a coherent implementation and qualification boundary. Explain
any genuinely necessary primitive universe or syntax distinction rather than
silently reintroducing privileged Type/Data names. The examples below are tests
of competing interpretations, not a frozen solution.

## Context

### Vocabulary and ownership

A **value** is a runtime datum or, in the optional mixed domain, a datum containing
an endpoint capability. A **proper type** describes values. A **type constructor**
accepts type-level arguments and produces a type. A **kind** classifies types or
type constructors; the requested general system must go beyond constructor arity.
Whether a type constructor may itself return another constructor is reopened
here; “produces a type” is not intended to preserve the previous restriction to
proper results. A **kind constructor** instead accepts kind arguments and
constructs a kind, as in a List kind constrained by its element kind. These are
different levels, even when their names are similar.
An **AST node kind** in Bittype's carrier is merely a language-owned syntax tag,
not the semantic kind system under discussion.

| Component | Existing responsibility |
| --- | --- |
| Ontos | Finite immutable ground values: atoms and tuples, structured identity, canonical encoding, and data embeddings such as text/list/map |
| Horos ground API | Data descriptors and bounded well-formedness/value conformance, implemented in Go, Rust and TypeScript |
| Horos declarations | Canonical modules, immutable imports, names, generic binders, applications, static Wire expressions and specialization; currently Go only |
| Optional Ontos/wire | Mixed runtime value carrier, above unchanged ground Ontos values |
| Bittype | Go implementation of generic AST/CST carriers, provider contracts, semantic graph, language engine and concrete language tooling |
| Bitwire/Bitruntime | Live capabilities, transport and hydration; outside this static change |
| Contract libraries | Ordinary Function, Stream, Cell and future declarations and behavior; no special recognition in general infrastructure |

Wire is a static type constructor for endpoint capabilities with a message type.
No live endpoint occurs inside a canonical module: the *description* of Wire is
ground data, while its instances are capabilities. A valid type expression does
not establish protocol behavior or runtime authority. A structural resemblance
to Cell or Function never authorizes selecting that library's behavior.

The owner requires canonical definitions as input, not TypeScript extraction:

```text
canonical Horos module values
  -> Horos resolution/checking (authoritative semantics)
  -> Bittype Horos AST + provider-owned semantic facts
  -> AST[TypeScript | Go | Rust | C++] -> CST[target] -> source
```

Here CST means concrete syntax tree including tokens/trivia. Parsing, formatting
and printing are separate registered stages. Bittype's generic engine checks
carrier/snapshot/receipt integrity and finds registered transformation paths;
it must never implement Horos kinds or know names of behavioral libraries.
Existing source-to-native generation is types/interfaces only, with one open
Page generic and ordinary applications in four target languages. Unsupported
representations return a located error and no usable output. No runtime adapters,
serializers, generators of function implementations or protocol code are in scope.

### The new cases, before implementation details

The immediate target is ordinary definitions such as the following. This is
illustrative kind notation, not canonical syntax:

```text
kind Flat = Scalar | List<Scalar>
kind Rows = Record<Flat>
kind Tree = recursive X. Scalar | List<X> | Record<X>
kind Data = recursive D. <all appropriate data families, recursively using D>
```

Here Scalar, List and Record are *kind families*. Lowercase `text`, `list<T>`
and `wire<T>` below denote *type expressions*; Person and Page are named type
declarations. WireFamily denotes the unconstrained outer capability-type family
when discussing kind expressions. Capitalization is explanatory notation, not a
proposed namespace or parsing rule. A bare Record/List family in a comparison
means its unconstrained outer shape, not an implicit partial kind application.

Record<K> means that each field type has classification K, so a data-only
definition needs Record<Data>, not unrestricted Record. The `recursive` spelling
leaves least/greatest fixed-point semantics for the expert to choose explicitly,
including how ground `rec`/`ref` types are classified. It does not imply that
recursive kind binding is the same operation as a recursive ground descriptor.

Flat, Rows, Tree and Data must all use the same definition/checking mechanisms.
An imported user-defined category must work without adding a checker branch,
renaming a flag or extending a central named-kind enumeration. The full set of
required cases and alternative denotations appears below.

### Inspected revisions and prior consultation

- Horos main: `619a10351d33b6aa4d150781eaa4b845b719cd3a` (revision-3 declaration
  model; source changes merged, pending 0.10.0 release held for this correction).
- Bittype main: `1672ed4c07ed19c1b1a4fd60ec06bd8788d753e5` (revision-2 frontend
  and implemented native pagination).
- Bittype PR 57 code candidate: `7082f9ce0a1f2ad0ac81fa1394454f16147cf367`
  (revision-3 syntax, generic application callees and subject-associated facts,
  authoritative Horos bridge and preserved native pagination; not merged).
  Later commit `8ecc0b6` only marks consultation 0004 superseded.

Research 0004 recommended saturated higher-order constructor inputs with a
separate compositional Data analysis. Its input was SHA-256
`e7941dee4413a1becd65bd32e9a256bcaead4b04619949fd7f3565877fad773d`;
run `run_03dd46c7-4356-4db4-92cf-2e2fd3f98645` completed, and verified advice
SHA-256 is `b56af93bfc33b9d46e98c96289b2bea224a27af8a8fbf8631a610828b3f95469`.
The core of that advice was:

```text
erased shape: K = Type | (K1, ..., Kn) => Type
interface:    S = Type | Data | (S1, ..., Sn) => R
result:       R = Type | Data
```

Its implementation is useful evidence and a migration starting point, but the
fixed proper-result vocabulary and Data-specific analysis do not meet the new
requirement. Do not treat that prior choice or its deferrals as architectural
constraints on this consultation. The frozen advice and historical validation
remain preserved, not a second implementation to keep indefinitely.

Reuse checks covered research files/branches/worktrees, open PRs, the platform
research reservation and Nightfall's collection/consultation list. Research
0002 established provider-owned semantics and exact receipts; 0003 established
data/identity/runtime-authority boundaries. Neither answers this general-kind
question. The independent source-composition documentation proposal concerns
Files/Source/Semantics and Syntax/CST/AST; its new APIs are not yet implemented
and do not determine the kind calculus.

### Ground descriptors and the current declaration format

An Ontos atom is an arbitrary finite byte string; a tuple is a finite ordered
sequence of Ontos values. Data embeddings are conventional labelled tuples.
Horos descriptor syntax is `Tuple(Atom("horos/" + form), operands...)`.
The existing form meanings are fixed and must not change as a side effect:

| Form | Meaning and type-bearing positions |
| --- | --- |
| any, atom and atom refinements | All ground values, or atom subsets; bounds/literals are metadata |
| scalar(int / utf8-text / bool), text refinements | Data embeddings; scalar selector and bounds are metadata |
| exact(value) | Singleton ground value; its operand is opaque literal data, not recursively interpreted type syntax |
| tuple(T...), compound(label,T...) | Ordered child types; compound label is metadata |
| seq(T), list-of(T), set-of(T) | Homogeneous bare tuples, embedded lists or sets |
| map-of(K,V) | Embedded map with key and value types |
| record(open/closed, fields...) | Embedded text-keyed map; each field has a name, req/opt marker and type; open mode permits additional unconstrained ground values |
| union(T...) | Union of value sets; empty union has no values |
| and(T...) | Intersection of value sets; empty intersection admits all ground values |
| not(T) | Complement relative to the ground value domain |
| rec(name,T), ref(name) | Lexically bound recursive descriptor; every ref must occur after strict value descent |
| wire(T) | Declaration-layer capability type; message position is a proper type; not accepted by ground WF/Conforms |

Ground `rec`/`ref` are not generic parameters or general type-level recursion.
Structural child descent through tuple/collection/record guards a recursive
value check; union/intersection/complement do not. All values themselves remain
finite acyclic trees. Ground equality is exact structural identity, not semantic
type equivalence. Ground conformance provides explicit exhaustion diagnostics;
the general-kind model must not turn exhausted analysis into nonmembership.

Current declaration revision 3 uses this exact notation and grammar. `A` is a
UTF-8 atom, `T` is a text embedding, `I` an integer embedding, `L` a list
embedding, `M` a canonical text-keyed map, and `D(tag,...)` prefixes a tuple with
`A("horos/decl/" + tag)`:

```text
Claim       = A("type") | A("data")
Interface   = D("proper", Claim)
            | D("constructor", L(Interface, ...), Claim)
Parameter   = D("parameter", T(name), Interface)
Declaration = D("type", L(Parameter, ...), Claim, E)
Module      = D("module", I(3), T(moduleName), importsMap, declarationsMap)
Import      = D("import", T(expectedModuleName), Atom(32-byte module digest))

E = Horos descriptor with E only in schema-designated child-type positions
  | Tuple(A("horos/wire"), E)
  | D("param", I(zeroBasedSlot))
  | D("local", T(name))
  | D("imported", T(alias), T(name))
  | D("apply", E, L(E, ...))
```

All imports pin exact canonical module bytes with SHA-256; declaration identity
is `(module digest, declaration name)`, and parameter identity adds the ordered
slot. Constructor arguments retain their definition scope; other arguments
retain caller scope. Empty application is rejected; zero-parameter declarations
are referenced directly. All current declarations are exported. Descriptor
metadata never becomes a type parameter just because it resembles tagged syntax.

Existing data/capability admission is intentionally nonuniform. Tuple, compound,
list-of, closed record and union may contain mixed types. Seq, set-of, map-of,
open record, and, not and guarded recursion currently require Data operands.
This protects operations whose mixed interpretation is not implemented. The new
model must express these form contracts without a checker branch for a named
Data kind, or explicitly justify a changed contract as separately scoped work.

### Actual specialized implementation

Horos `core/go/internal/declmodel/interfaces.go:10-53` contains these declarations
and the shape-erasure body (accessor cloning details omitted):

```go
type Claim string
const (
    TypeClaim Claim = "type"
    DataClaim Claim = "data"
)
type Interface struct {
    inputs []Interface
    result Claim
}
func (i Interface) IsConstructor() bool { return len(i.inputs) != 0 }
func (i Interface) Kind() Kind {
    result := Kind{inputs: make([]Kind, len(i.inputs))}
    for n, child := range i.inputs {
        result.inputs[n] = child.Kind()
    }
    return result
}
type Kind struct{ inputs []Kind }
```

Thus Scalar, Record, List and Wire proper expressions all erase to the same zero
shape. `ProperInterface`, `validateInterface` and `readClaim` reject any result
other than TypeClaim/DataClaim. `compatible` at lines 151-169 compares input
interfaces contravariantly and rejects a non-Data result when Data is expected.
It has no general kind expressions, named kind declarations, union or inclusion.

`core/go/internal/declmodel/bounds.go:135-152` dispatches membership as follows:

```go
func (c *proofChecker) satisfies(v *semanticValue, expected Interface, depth int) (bool, error) {
    shape, err := sameShape(v.interface_, expected, c.work, depth)
    if err != nil || !shape {
        return false, err
    }
    if !expected.IsConstructor() {
        if expected.result == TypeClaim {
            return true, nil
        }
        status, err := c.data(v, depth+1)
        return status == DataProven, err
    }
    if v.rigid {
        return compatible(v.interface_, expected, c.work, depth+1)
    }
    return c.declarationSatisfies(v.declaration, expected, depth+1)
}
```

Here `semanticValue` is a checked proper-expression closure, a known declaration
constructor, or a fresh rigid symbolic parameter with an Interface. `sameShape`
recursively compares input counts; `c.work` is the operation's aggregate budget.
`declarationSatisfies` checks original argument admission, then evaluates the
known body under fresh symbolic expected inputs only if the expected result is
Data. Abstract constructor parameters use only their declared guarantees.

The dedicated summary representation (`interfaces.go:170-213`) is:

```go
type SummaryKind string
const (
    SummaryTrue  SummaryKind = "true"
    SummaryFalse SummaryKind = "false"
    SummaryNeed  SummaryKind = "need"
    SummaryAnd   SummaryKind = "and"
    SummaryProbe SummaryKind = "probe"
)
type DataStatus string
const (
    DataProven      DataStatus = "proven"
    DataFailed      DataStatus = "failed"
    DataConditional DataStatus = "conditional"
)
type DataSummary struct {
    kind      SummaryKind
    parameter ParameterID
    children  []*DataSummary
    probe     *Expression
}
```

Need references a proper parameter; Probe retains checked Reference/Application
syntax. Its binding environment is supplied separately by operation-local
semantic values and the summary evaluator, not stored in the exposed Probe.
Together these preserve conditional Data information through
Wrapped and Compose, rather than permanently classifying every unknown call as
false. A failed compositional Data certificate does not assert that a non-data
inhabitant exists. No emptiness/equivalence solver is implemented.

Public `Program.Satisfies(id, expected Interface)` asks about a named declaration;
`Program.Specialize(id, []Argument)` substitutes typed arguments. `Argument`
handles are sealed and belong to one checked Program, even if another Program
has identical source bytes. `Program.Proper(value)` admits closed proper types;
`Program.Constructor(id)` admits named constructors. `DataDescriptor` specializes
and additionally establishes Data before returning a ground descriptor for
WF/Conforms. There is no mixed value validator.

The replacement must remove this dedicated proof route rather than keep it as
a fallback next to a general kind checker. A ground API boundary may still need
to establish that a result is acceptable to ground Horos; please specify how it
does so through the ordinary kind model or another justified boundary check.

The primitive formation rules are specialized too. Horos
`core/go/declarations/grammar.go:101-116` selects a private
`DomainRule{RequiresData bool; Capability bool}`:

```go
func descriptorRule(constructor string, args []core.Value) model.DomainRule {
    switch constructor {
    case "wire":
        return model.DomainRule{Capability: true}
    case "tuple", "compound", "list-of", "union":
        return model.DomainRule{}
    case "record":
        if len(args) > 0 {
            if mode, ok := args[0].(core.Atom); ok &&
                string(mode.Bytes()) == "closed" {
                return model.DomainRule{}
            }
        }
    }
    return model.DomainRule{RequiresData: true}
}
```

This rule drives both open-body checks and closed classification. Closed proper
arguments also invoke the grammar's separate Data classifier before assigning
their current status (`core/go/internal/declmodel/specialize.go:77-89`). Public
`DataDescriptor` specializes and then calls that classifier again. Generalizing
only source annotations leaves these authorities in place. Primitive descriptor
semantics remain Horos-owned; the open question is how their operand/result
formation contracts are expressed and consumed by the general kind algebra.

The concrete call flow makes the duplication visible:

```text
Resolve
  -> read/check all declaration bodies and dependency cycles
  -> checkBounds: rigid inputs, operand admission, Data summaries/result claim
Specialize
  -> re-admit Program-owned arguments
  -> satisfies -> expand -> validate final result
DataDescriptor
  -> Specialize -> descriptorData -> return unchanged ground descriptor
```

These calls are in `core/go/declarations/model.go:76-84,145-159`,
`core/go/internal/declmodel/resolve.go:115-137`, `bounds.go:3-20,74-100`, and
`specialize.go:54-124` in that same internal directory.

### Bittype: useful generic infrastructure and specialized consumers

Bittype PR 57's `semantics/go/model.go:20-121` stores provider-owned graph data.
The following are exact field shapes; IDs are local one-based uint32 addresses,
and `ast.NodeRef` is a tree-qualified AST address:

```go
type EntityID uint32
type TermID uint32
type Key struct{ Owner, Name string }
type Parameter struct {
    Entity EntityID
    Sort   Key // coarse provider category, not a structural semantic kind
}
type Term struct {
    Kind      TermKind // Reference, Application, or ProviderExpression
    Target    EntityID // Reference only
    Callee    TermID   // Application only
    Arguments []TermID
    Operator  Key     // ProviderExpression only
    Encoding  Key
    Data      []byte
    Operands  []Operand
}
type Operand struct {
    Locator []uint32
    Term    TermID
}
type Subject struct { Entity EntityID; Term TermID } // exactly one nonzero
type Fact struct {
    Subject  Subject
    Relation Key
    Term     TermID
}
```

`Entity` holds a provider-owned kind key, name, owner, ordered real parameters,
body term and declaring AST sites. The generic model allows multiple different
fact values for the same subject/relation, rejects exact duplicate triples,
copies all mutable input, validates links and applies resource budgets. Facts
are associations, not containment. Callee/argument/operand links are containment
and must form an acyclic bounded graph. A provider owns the interpretation of
keys, payloads, kinds, application applicability and arity. The generic engine
does not infer those rules from spelling or graph shape.

Resolution adapts the authoritative checked Horos Program into these terms.
Analysis consumes an explicit compatible resolution receipt and publishes the
same exact-snapshot evidence; there is no hidden resolve fallback. Completion,
validity and model availability are separate. A missing required obligation
prevents model publication. Any input/configuration/dependency change requires
a fresh snapshot and receipt. The consumer then explicitly supplies semantic
input to the transformer.
Concretely, Decode reconstructs canonical values in
`internal/horosmodel/read.go:254-292`; Resolve calls the configured authoritative
resolver and verifies projection correspondence at lines 295-309. The binding
to `decl.ResolveWithLimits` is in `resolve/horos/go/resolver.go:17-23`.
`internal/horosmodel/resolver.go:188-206` orchestrates those calls and exports
the resulting facts. `analyze/horos/go/analyze.go:46-68` reuses the explicit
resolution result.

The current Horos parser, canonical bridge and semantic fact producer encode
only Type/Data proper interfaces and constructor arrows. The native reader
likewise understands those two claim payloads and currently expects one
interface fact. These provider-owned readers require replacement; changing
`Parameter.Sort` or the generic carrier is not sufficient. A future principal
classification and a set of established memberships must not be conflated.

Concrete examples from the inspected candidate:

```go
// parse/horos/go/parse.go:188-198
func (p *parser) properInterface() ast.NodeID {
    start := p.Token().Start
    name := "type"
    switch p.Take().Text {
    case "Type":
    case "Data":
        name = "data"
    default:
        p.Fail(text.SyntaxError, "expected Type or Data interface")
    }
    return p.Add(s.KindInterface, start, tr.Attrs(s.AttributeName, name))
}

// internal/horosgen/generate.go:181-186
if fact.Relation == resolve.InterfaceRelation() {
    if g.interfaces[fact.Subject] != 0 {
        return fail(semantic.ErrInvalidRequest)
    }
    g.interfaces[fact.Subject] = fact.Term
}

// internal/horosgen/generate.go:235-263 (claim check excerpt)
if claim := atom(value); claim != "type" && claim != "data" {
    return false, semantic.ErrInvalidRequest
}
```

The parser's builder creates an AST node at a source offset; syntax admission
allows unresolved names generally, but interface positions currently cannot hold
them. The native reader's map stores one interface term per semantic Subject.
`internal/horosmodel/facts.go:41-161` exports `kind/proper`, `kind/constructor`,
`interface/proper`, `interface/constructor`, `data-summary/*` and `data-status/*`
terms as ordinary provider expressions, and `internal/horosmodel/project.go:200-222` projects
`decl.Claim`/`decl.Interface` into syntax. These are concrete replacement points,
not rules to move into the generic model.

### Current examples and required generalized cases

The current readable source is:

```text
type Person = { name: text; };
type Page<T> = { items: list<T>; next?: text; };
type DataPage<D: Data> = Page<D>;
type Ignore<T> = text;
type Send<T> = wire<T>;
type Wrapped<F: Type -> Type, T> = { value: F<T>; };
type Compose<F: Type -> Type, G: Type -> Type, T> = F<G<T>>;
```

One Page identity must work for data and capability elements; no authored second
Page or generated template-lifting copy. Known Page establishes a data result
for data inputs, Ignore produces text regardless of its admitted input, and
Send produces a capability. Under abstract F the checker may use only F's
declared contract, and specialization cannot repair an invalid open body.

The new model must also express constraints unrelated to Data, for example
“any list whose elements are scalar types,” “any record whose fields are scalar
or list-of-scalar types,” and an imported user-defined recursive category of
trees of records and lists. Changing that category's name, defining another
equivalent expression, or adding a differently constrained category must not
require checker or Bittype-kernel code changes.

The Flat/Rows/Tree examples above classify *types* with constrained children.
Clarify the denotation rather than assuming that notation provides an algorithm.
Scalar selectors, field labels/presence/open
mode, literal values and recursion binders cannot be silently erased if they
affect validity or an advertised kind guarantee.

The following distinctions must be resolved with concrete judgments:

1. `list<wire<Person>>` is a list-shaped type but is not Data. Likewise a closed
   record with a Wire field is a record but not Data. A shallow
   `Scalar | Record | List` definition incorrectly admits both.
2. A single type `union(text, {name:text})` has values of two shapes. A kind union
   `Scalar | Record` is a union of *sets of types*: it need not contain that
   union-shaped type. Compare classifying syntax families with classifying
   all types whose value sets lie inside a given value category.
3. `Record<Scalar | WireFamily>` can permit one scalar field and one Wire field;
   `Record<Scalar> | Record<WireFamily>` need not. Do not assume distributivity of
   kind parameters without defining it.
4. `exact(D("param", I(0)))` is a ground singleton; its literal remains opaque.
   Decide whether exact-text is Scalar by semantic shape, an Exact family,
   both, or an explicitly established relation. No arbitrary equivalence solver
   is currently available.
5. `union()` is empty but data-admissible. `wire<union()>` remains a capability
   type even if its message type is empty. Intersections/complements can create
   semantic equivalences or empty types that syntactic classifications miss.
6. Ground guarded recursive records can be Data, but recursive kind definitions
   need their own positivity/guard/termination rules. Merely finding a memoized
   cycle must not prove an unguarded kind or ground recursion valid.
7. Page's outer record shape is fixed, but its nested-kind guarantees depend on
   its argument classification. Explain whether ordinary
   kind parameters/quantification, a symbolic structural result, multiple
   established arrows, or another mechanism expresses this *generally*.
   For a concrete non-Data example, an abstract F admitting Flat inputs and
   guaranteeing Record<Flat> results must let `F<A>` establish Rows for A in
   Flat. A consumer requiring Rows must accept that result. An abstract F
   guaranteeing only Flat cannot establish that
   its result is Record<Flat>, because it may return text; requiring that result
   must fail. For the first F, which guarantees Rows, a wrapper `{value: F<A>}`
   has Record<Rows> structure, but does not automatically establish Rows:
   Rows membership alone does not establish Flat membership. Rebuilding the dedicated Data conditional-proof
   subsystem is insufficient. Quantification, joint arrows and symbolic result
   structure remain candidate mechanisms, not assumed syntax.
8. A union of two constructor interfaces need not guarantee both contracts.
   Distinguish an unknown alternative from a constructor satisfying both input/
   result guarantees, and define arrow variance and saturation separately.
9. An empty closed record may satisfy `Record<K>` vacuously for many K. An
   optional Wire field still permits a capability when present. A record with
   both a Wire field and an impossible required field is uninhabited; structural
   Data and an extensional all-inhabitants-are-data test can disagree.
10. Subtracting the outer Wire family from a universal kind still leaves
    `Record<WireFamily>` and `List<WireFamily>`. `Type minus WireFamily` is not a recursive Data
    definition. If extensional Data admits an empty but syntactically mixed
    type, its unchanged descriptor still contains Wire and fails ground WF;
    the ground boundary must not pretend the subset proof rewrites syntax.

### Bounded operation and identity constraints

Canonical input is untrusted and resolution performs no IO. Default Horos
limits are depth 256, 100,000 nodes/work units, 16 MiB and 256 modules, with an
absolute depth maximum 512. Raw admission, aggregate semantic work and final
expanded output have separate budgets. Generated cache keys, environments,
interface comparisons and growth consume aggregate work/bytes; output charges
repeated occurrences even if internal structure is shared. Caches are local to
one operation; checked Programs and exposed facts are immutable and can be read
concurrently. Failures publish no usable partial Program/model/output.

All open definitions and all arguments are checked, including ignored ones.
Every named reference contributes cycle dependencies, including an unused
constructor passed to Ignore. Declaration cycles are currently rejected;
guarded ground recursion is handled separately. Kind-definition recursion will
require an explicit graph rule rather than weakening all cycles. Unknown
membership due to open parameters, completed failure, malformed requests and
exhaustion must remain distinguishable.

For a possible kind denotation based on values, write
`C(S) = {T | values(T) is a subset of S}`. In general
`C(ScalarValues) union C(RecordValues)` differs from
`C(ScalarValues union RecordValues)`: a mixed type union witnesses it.
This distinguishes the two classification levels without assuming either one
is the recommended design. Similarly `not(text)` is a data type because its
complement universe is ground data; it is not the kind complement of Scalar.

Names and import pins identify declarations, not extensional equality.
Subkinding/equivalence may establish relations without changing identity or
selecting a library binding. Kind-definition/schema dependencies and any prelude
must participate in exact receipts and bounded snapshots. A kind name alone
cannot confer trust. The owner requires clean replacements: no revision-3
fallback reader, old Data route, compatibility alias, second language ID or
executable legacy archive merely to keep callers working.

### Native generation and independent qualification

The implemented pagination mapping supports text/bool data scalars,
records, lists, optionals, declarations and ordinary generic applications under
explicit profiles. It preserves open Page and application aliases, rather than
monomorphizing everything. Current native readers reject residual constructor
parameters across the entire exported closure with no output, even if a selected
closed use could be reduced. The general-kind system must not falsely imply
that TypeScript, Go, Rust or C++ can enforce every structural kind natively.
Canonical integer scalar types are admitted by Horos, but their native mapping
is not implemented by this profile and returns Unsupported.
The current native profile also rejects unused proper parameters across that
closure, so the canonical Ignore example above is not itself a claim of native
generation support (`internal/horosgen/generate.go:211-214`).

Advise what semantic information the generator should consume and what should
be erased only under an explicit representation contract. Do not solve target
limitations by silently weakening Horos constraints, dropping unsupported
declarations or adding runtime adapters. Broader native encodings can be a
separate qualified milestone with a concrete trigger.

Existing qualification supplies a useful base. Horos main `619a103` passed all
six lanes in [CI 41860921097](https://github.com/Bitspark/horos/actions/runs/41860921097),
including independent eager expansion/Data oracle, open-body, identity, literal,
recursion, foreign-Program and budget tests. Bittype `f5256ed` passed
[contracts/native 41860252057](https://github.com/Bitspark/bittype/actions/runs/41860252057)
and [docs 41860252132](https://github.com/Bitspark/bittype/actions/runs/41860252132):
parser/formatter round trips, provider graph/receipt tests, four pagination
outputs and native compilers, and an independently installed public Go consumer.
Candidate `7082f9c` has identical Go/workflow sources to `f5256ed`, with additional
documentation; its exact-head hosted conclusions were still pending when this
brief was prepared. None of these runs qualifies an unimplemented general-kind
replacement. Tests run only on maintained remote/shared gates, not this
workstation. New general-kind tests need an
independent oracle or explicit countermodels, including renamed/user-defined
categories that would expose a hidden Data-specific branch.

## Options and unresolved design choices

1. **Closed Type/Data plus more enum entries:** rejected by owner. It still
   centralizes named categories and cannot express nested constraints generally.
2. **Structural recursive kind expressions:** kind terms denote families of
   normalized type shapes, with unions and child constraints. This appears to
   fit the request, but aliases/applications, unranked record/tuple children,
   recursion, exact literals and subset checking need a precise treatment.
3. **General refinement predicates over types:** expressive structural tests,
   rows, labels and logical combinations can classify more cases. Solver
   dependence, proof status, termination and canonical reproducibility must be
   controlled; “use SMT” alone is not an implementation contract.
4. **Value-denotation classifications:** kinds classify types by the values they
   admit. This makes semantic Scalar membership intuitive but risks importing
   general subtyping, emptiness and equivalence into every kind check.

A union of structural kinds and a union type descriptor operate at different
levels. A standard library may define Data and a broad Type category, but the
checker needs a general way to interpret those definitions. Whether a primitive
proper-type universe is required, how its public spelling works, and how it
interacts with new structural families must be explicit. A central switch moved
into generated code or arbitrary named callbacks is not a general declarative
kind algebra.

### Primary-source precedents and their limits

- Caires and Toninho's [Refinement Kinds](https://arxiv.org/html/1908.00441v1)
  classifies structural properties of types, including record and collection
  structure, through predicates over type trees. It is a direct precedent for
  richer kinds. Its implementation discussion says the underlying SMT theories
  are generally undecidable and completeness is solver-relative. It does not
  supply a bounded complete checker for our requested combination.
- Hosoya, Vouillon and Pierce's [Regular Expression Types for XML](https://www.cis.upenn.edu/~bcpierce/papers/regsub-toplas.pdf)
  models recursive tree languages and checks inclusion, a decidable problem
  with EXPTIME-complete worst-case complexity. Sharing and lazy exploration
  improve practical checking. Applying that machinery to normalized Horos type
  shapes is our inference, not a result proved by that paper. Variadic record/
  tuple operands require sequence structure rather than an assumed fixed arity;
  arbitrary equality of repeated subtrees goes beyond regular shape membership.
- The [CDuce manual](https://www.cduce.org/manual_types_patterns.html) combines
  Boolean types and recursive definitions, but recursive cycles must cross
  structural constructors; Boolean connectives do not guard a cycle. Its
  parameterized recursive definitions also restrict recursive instantiations.
  This motivates an explicit recursive-kind rule, not wholesale import of
  CDuce's value-denotation semantics or its function types.
- Gaster and Jones's [extensible records and variants](https://web.cecs.pdx.edu/~mpj/pubs/96-3.pdf)
  separates row-polymorphic structure from ordinary type parameters. A general
  Record family can constrain all field types without already supporting
  arbitrary field-preserving row transformations. Existing Horos open-record
  mode must not silently become an unconstrained row variable.
- [GHC's kind-polymorphism guide](https://downloads.haskell.org/ghc/latest/docs/users_guide/exts/poly_kinds.html)
  provides explicit/inferred kind schemes and higher-order constructor
  examples, while separately documenting polymorphic recursion and features
  that can permit nontermination. Kind polymorphism and higher-order arguments
  are separable; adopting schemes requires a bounded combined calculus here.

No cited precedent establishes that regular structural kinds, aliases,
recursive descriptors and quantified higher-order interfaces remain decidable
when combined without restriction. That combined obligation is part of the
requested expert assessment. Positive guarded structural schemas with explicit
arrow contracts and limited schemes, full Boolean schemas over a fixed
constructor universe, and restricted refinement/row constraints are distinct
implementation options with different costs.

## Questions for the expert

1. **Meaning and algebra:** What denotation and core kind grammar best meet the
   owner's request? Define structural Scalar/Record/List kinds, their parameters,
   union/intersection if needed, recursive definitions, and a complete ordinary
   definition of Data covering the existing Horos forms. Explain Type/universes
   and distinguish kind union from type union with the counterexamples above.
2. **Constructors and general result information:** How should higher-order
   arguments, arrow contracts and Page-like preservation be represented and
   checked without a Data-specific proof system? State the role of kind
   variables/quantification, known versus abstract constructors, result
   inference, variance, application and saturation, and what is separable.
   Does the algebra admit principal inferred kinds, or should consumers request
   particular membership proofs? Define completeness of the published evidence.
3. **Representation, extension and bounded algorithms:** Give exact canonical
   declaration/kind forms, binding/import/identity rules and an implementable
   membership/subkinding algorithm with explicit decidability and resource
   limits. Distinguish adding an ordinary kind definition from a genuinely new
   primitive type form; explain ownership and trust without central named cases.
4. **Migration, semantic facts and native consumers:** What should Horos expose
   and Bittype project, retaining its agnostic graph and exact receipts? Specify
   replacement of Claim/Interface/DataSummary, preservation of ground APIs and
   pagination, native rejection/representation policy, and independent tests
   that would falsify the design before delivery.
5. **Open assessment:** Which assumption or requirement above is mistaken or
   missing? Is there a simpler general construction that meets the owner's
   intent more directly? Identify any tradeoff that our framing obscures and
   give the strongest counterexample to your recommended design.
