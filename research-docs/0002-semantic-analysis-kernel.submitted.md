# Research: An agnostic semantic-analysis kernel and Horos provider

**ID:** 0002
**Date:** 2 November 2026
**Status:** reviewed
**Owner:** codex-api-retirement-audit
**Issue:** https://github.com/Bitspark/bittype/issues/16
**Reviewed:** https://github.com/Bitspark/bittype/issues/16#issuecomment-6578794456

## Question

What is the smallest sound shared semantic representation and provider boundary
for a Go-implemented, language-agnostic syntax engine, and how should its first
real analyzer adapt existing Horos declaration semantics without moving those
semantics into the engine? The owner explicitly requested expert consultation
followed by implementation, starting from the minimal model reproduced below.
We seek an implementable design, including adverse cases and a bounded first
qualification, not a universal compiler framework.

## Context

### System overview and vocabulary

Bittype is a types/declarations toolchain implemented in Go. An AST is an abstract
syntax tree. A CST is a concrete syntax tree retaining tokens/trivia (whitespace
and comments) and a view of
its AST. A language marker is a distinct Go type parameter separating, for example,
AST[Go] from AST[Horos]. A profile identifies the explicitly admitted syntax or
analysis contract; it is not a guessed compiler version or structural hash.
SPI means service provider interface: the interface implemented by a provider
and called by the engine.

The architectural change under discussion is:

```text
Current Horos: syntax admission -> resolve/check all semantics -> admitted AST
Target: syntax admission -> admitted unresolved AST -> explicit analyzer
        -> shared semantic facts + diagnostics + actual analysis coverage
```

Intended agnostic consumer queries are: resolve an occurrence to its declaration;
inspect the arguments and parameter bindings of an application; traverse recorded
declaration dependencies; and determine whether required analyses completed with
the needed verdicts. Descriptor interpretation remains provider-owned. There are
no production consumers of a shared semantic model or actual AST transformers
yet. The first deliverable must make these queries testable without promising
complete cross-language type comparison or code generation.

Current production source is bittype
45f7c2782d7f77b9574dae1a9893ee81b2bba427. The independent Horos source is
097496351b70336462c659e3deabea833e294b85. All code excerpts below describe those
revisions. There is one Go module in bittype, targeting Go 1.27 with CI toolchain
1.27.3. It has AST/CST carriers, typed syntax for TypeScript/Go/Rust/C++, native
text parsers and formatters, generic CST printing, explicit language registration
and directed AST-transformer routing. No actual language-to-language transformer
is implemented. Native checkers admit a bounded type/declaration syntax slice;
they are not full compiler semantic analyzers.

Ontos is the independent immutable value model: atoms contain arbitrary bytes;
tuples contain ordered values. Existing embeddings interpret these as text,
integers, maps and lists. Horos owns data-type descriptors over Ontos values and
value conformance. It also now owns a declaration layer: modules, pinned imports,
type parameters, and applications. A module pin is SHA-256 of its canonical Ontos
codec-v1 bytes: an exact document address, not semantic equivalence. A descriptor
is an Ontos tuple whose first element is an atom such as `horos/list-of`, followed
by constructor operands. `exact(v)` denotes the singleton type containing
precisely the literal value v; it is not a place to interpret declaration syntax.
Bittype's existing Horos frontend
projects a resolved Horos Program into an AST and can recover the Program.

The generic kernel must never import Horos/Ontos or native language implementations.
Only the Horos frontend and explicitly scoped future transformers may depend on
them. Conversely, Horos must not depend on bittype. Function, Stream and Cell are
ordinary consumer library contracts, not privileged semantic/runtime kinds.
The owner's previous rejected implementation dispatched on such names in general
runtime code; relocating that dispatch into a semantic model is also forbidden.

### Architectural context and current layout

The enforced layout is stage-first, intentionally a matrix across languages:

```text
trees/ast/go                 generic syntax carrier
trees/cst/go                 concrete carrier with AST view
trees/representation/go      owned heterogeneous AST/CST snapshots
contracts/language/go        ID and Name
contracts/text/go            typed Parse/Check/Format/Print providers, diagnostics
contracts/transform/go       ASTTransformer and typed adapter
engine/go                   explicit registry, routes and execution
syntax/<language>/go        syntax vocabulary and marker
parse/<language>/go         text -> CST
check/<language>/go         syntax admission (Horos currently also resolves)
format/<language>/go        AST -> CST
provider/<language>/go      explicit stage composition
resolve/horos/go            canonical values <-> checked Horos AST/model
internal/horos              private Horos projection/recovery implementation
integration/go             independent consumers and architecture tests
```

Native providers have text capabilities. Horos has only syntax/resolve/check/provider,
with profile `horos-declarations-1/ast-1`; it has no readable parser or formatter.
There is no production language catalogue or automatic registration. Tests
discover stage directories and require the declared capability matrix. Every
production package, including private helpers, has enforced directed imports.
Providers are trusted Go code, not sandboxed plugins. The engine still protects
caller/result ownership against accidental mutation and checks type boundaries.

### Relevant code: syntax and registration boundary

`trees/ast/go/node.go:6-44`, `ast.go:11-35`, `view.go:5-10` and
`origin.go:3-32` define the following carrier.
Kinds, attribute names and edge roles have language-owned meaning. NodeID is an
index in an arena, TreeID is a tooling address, and neither proves semantic identity.
The actual AST also carries its root and node slice, plus origins with source
spans and historical node references.

```go
type Kind[L any] string
type TreeID string
type NodeID uint32
type NodeRef struct { Tree TreeID; Node NodeID }
type Attribute struct { Name string; Text string }
type Child struct { Role string; Node NodeID }
type Reference struct { Role string; Target NodeRef }
type Node[L any] struct {
    Kind Kind[L]
    Attributes []Attribute
    Children []Child
    References []Reference
    Origin Origin
}
type AST[L any] struct { ID TreeID; Root NodeID; Nodes []Node[L] }
type View[L any] interface { ASTView() AST[L] }
type SourceID string
type Span struct { Source SourceID; Start, End uint64 }
type Origin struct { Spans []Span; Nodes []NodeRef }
func (tree AST[L]) ASTView() AST[L] { return tree }
```

`Reference` means a supplied resolved link. Unresolved spellings belong in
attributes, never in a fabricated NodeRef. Generic structure checks enforce
containment, valid references and resource bounds, but not name resolution.
Current shared limits are 200,000 nodes, depth 256, 16 MiB text and 1.6 million
list entries. Cross-tree reference closure is not a generic single-tree check.

`contracts/language/go/language.go:11-24` and `contracts/text/go/provider.go:11-31`:

```go
type ID string
type LanguageProvider interface { ID() ID; Name() string }

// In package text, with ast/cst/language imports:
type CheckerProvider[L any] interface {
    language.LanguageProvider
    Check(ast.View[L]) error
}
type ParserProvider[L any] interface {
    language.LanguageProvider
    Parse(Input) (cst.CST[L], error)
}
type FormatterProvider[L any] interface {
    language.LanguageProvider
    Format(ast.View[L]) (cst.CST[L], error)
}
type PrinterProvider[L any] interface {
    language.LanguageProvider
    Print(cst.CST[L]) (string, error)
}
type Provider[L any] interface {
    ParserProvider[L]; CheckerProvider[L]
    FormatterProvider[L]; PrinterProvider[L]
}
```

`trees/representation/go/language.go:15-45,75-81` captures a concrete Go marker
alongside metadata; `payload` is a private interface with marker() and clone(),
implemented by a typed AST wrapper. Here and in subsequent excerpts, imports and
unrelated functions are omitted, not error branches within the shown functions.

```go
type Language struct {
    id language.ID
    name, profile string
    marker reflect.Type
    check func(payload) error
}
func NewLanguage[L any](provider text.CheckerProvider[L], profile string) (Language, error) {
    if IsNil(provider) { return Language{}, ErrInvalidLanguage }
    id, name := provider.ID(), provider.Name()
    if strings.TrimSpace(string(id)) == "" || strings.TrimSpace(name) == "" || strings.TrimSpace(profile) == "" {
        return Language{}, ErrInvalidLanguage
    }
    return Language{id: id, name: name, profile: profile, marker: reflect.TypeFor[L](),
        check: func(p payload) error { return provider.Check(p.(typed[L]).value) }}, nil
}
func (l Language) Compatible(other Language) bool {
    return l.Validate() == nil && other.Validate() == nil &&
        l.id == other.id && l.profile == other.profile && l.marker == other.marker
}
func (l Language) Check(input AST) error {
    if err := l.match(input); err != nil { return err }
    return l.check(input.payload.clone())
}
```

`Validate` rejects missing metadata/marker/checker. `match` rejects an invalid
envelope or mismatched ID, profile or actual marker. IsNil includes typed nils.
`trees/representation/go/ast.go:53-80` is the admission and ownership boundary:

```go
func Pack[L any](definition Language, view ast.View[L]) (AST, error) {
    if err := definition.acceptMarker(reflect.TypeFor[L]()); err != nil { return AST{}, err }
    if IsNil(view) { return AST{}, ErrInvalidAST }
    input := view.ASTView()
    if err := check.Structure(input); err != nil { return AST{}, fmt.Errorf("%w: %w", ErrInvalidAST, err) }
    result := AST{id: definition.id, profile: definition.profile, payload: typed[L]{tree.Clone(input)}}
    if err := definition.Check(result); err != nil { return AST{}, err }
    return result, nil
}
func Unpack[L any](definition Language, input AST) (ast.AST[L], error) {
    if err := definition.acceptMarker(reflect.TypeFor[L]()); err != nil { return ast.AST[L]{}, err }
    if err := definition.match(input); err != nil { return ast.AST[L]{}, err }
    return tree.Clone(input.payload.(typed[L]).value), nil
}
```

This `representation.AST` differs from the generic `ast.AST[L]` shown above.
Its actual storage (`trees/representation/go/ast.go:25-38`) is:

```go
type AST struct {
    id language.ID
    profile string
    payload payload
}
type payload interface { marker() reflect.Type; clone() payload }
type typed[L any] struct { value ast.AST[L] }
func (p typed[L]) marker() reflect.Type { return reflect.TypeFor[L]() }
func (p typed[L]) clone() payload { return typed[L]{tree.Clone(p.value)} }
```

The clone copies the node arena and its mutable child, attribute, reference and
origin slices. AST envelope fields are private;
its zero value is invalid. It exposes language/profile and Snapshot, but no
marker-agnostic node catalogue. A future semantic graph's NodeRefs therefore
need validation against an input catalogue established by the typed adapter, or
a deliberately narrow new envelope operation. Historical Origin.Nodes must not
be mistaken for direct semantic declaration sites in the current input closure.

`contracts/transform/go/adapter.go:14-48` provides the adapter precedent:

```go
func NewTransformer[F, T any](from, to representation.Language,
    stage pipeline.Stage[ast.AST[F], ast.AST[T]]) (ASTTransformer, error) {
    if stage == nil { return nil, ErrInvalidTransformer }
    if err := representation.CheckMarker[F](from); err != nil { return nil, err }
    if err := representation.CheckMarker[T](to); err != nil { return nil, err }
    return bound[F, T]{from: from, to: to, stage: stage}, nil
}
func (b bound[F, T]) Transform(input representation.AST) (representation.AST, error) {
    if err := b.from.Check(input); err != nil { return representation.AST{}, err }
    source, err := representation.Unpack[F](b.from, input)
    if err != nil { return representation.AST{}, err }
    result, err := b.stage(source)
    if err != nil { return representation.AST{}, err }
    return representation.Pack(b.to, result)
}
```

Stage is a typed function `func(I) (O,error)`. The private bound struct holds
from/to descriptors and that stage. ASTTransformer has FromID, ToID and Transform;
the typed adapter also advertises its descriptors for early compatibility checks.
The engine registers a single language per ID and one transformer per directed
edge. It snapshots registry dependencies under a mutex and releases the lock
before callbacks. Duplicate/ambiguous registration is rejected; routes never
silently fall back after a semantic failure. These same ownership and explicit
composition principles should apply to analysis, not necessarily routing.

The current registry entry and duplicate boundary
(`engine/go/engine.go:37-40`, `registry.go:26-37`) are:

```go
type registration struct {
    language representation.Language
    text *Provider
}
func (e *Engine) register(r registration) error {
    e.mu.Lock()
    defer e.mu.Unlock()
    id := r.language.ID()
    if _, ok := e.languages[id]; ok {
        return fmt.Errorf("%w: %s", ErrDuplicateLanguage, id)
    }
    if e.languages == nil { e.languages = make(map[language.ID]registration) }
    e.languages[id] = r
    return nil
}
```

Engine holds a sync.RWMutex, languages map and directed transformer edges. Execute
captures its route, descriptors and transformer handles under RLock, explicitly
unlocks, then checks the input and invokes each transformer (`transform.go:30-63`).
It checks each returned endpoint, wraps step errors, and snapshots each output.
RegisterAnalyzer must be an explicit capability attachment, not a second
RegisterLanguage call that silently replaces an existing entry.

### Relevant code: current Horos frontend and declaration authority

`bittype/check/horos/go/check.go:13-16` simply calls `bridge.Read(view)` and returns
its error. `provider/horos/go/provider.go:13-17` forwards Check and implements
CheckerProvider[horos.Language]. `resolve/horos/go/resolve.go:16-31` exposes
Resolve(rootValue, dependencyValues, treeID), Project(program, treeID), and Model(view).
The private bridge implements this current call chain:

```text
representation.Pack -> Language.Check -> Horos Check -> bridge.Read
  -> generic structural admission and AST decoding
  -> Horos declarations.Resolve(rootValue, dependencyValues)
  -> projectChecked(program, input.TreeID)
  -> equivalent(expectedAST, suppliedAST)
  -> checked Program or error
```

`internal/horos/read.go:158-222` includes these crucial lines after decoding:

```go
program, err := decl.Resolve(rootValue, dependencies)
if err != nil { return nil, fmt.Errorf("%w: %w", ErrInvalid, err) }
expected, err := projectChecked(program, tree.ID)
if err != nil { return nil, err }
if err := equivalent(expected, tree); err != nil { return nil, err }
return program, nil
```

`equivalent` compares kinds, attributes, containment and every reference under
an arena-address mapping. It requires all expected references, rejects wrong
targets, and never repairs the supplied tree. Thus adding an analyzer over the
existing envelope alone would be circular: unresolved Horos ASTs cannot be packed.
The current profile explicitly promises a complete resolved closure in one tree.

Horos `core/go/declarations/model.go:16-203` supplies immutable checked facts:

```go
type Digest [32]byte
type DeclarationID struct { Module Digest; Name string }
type ParameterID struct { Declaration DeclarationID; Index int }
type Kind string
const (
    Descriptor Kind = "descriptor"
    Parameter Kind = "parameter"
    Application Kind = "application"
)
type Operand struct { path []int; expression *Expression }
func (o Operand) Path() []int { return slices.Clone(o.path) }
func (o Operand) Expression() *Expression { return copyExpression(o.expression) }
type Expression struct {
    kind Kind
    value core.Value
    constructor string
    operands []Operand
    parameter ParameterID
    target DeclarationID
    arguments []*Expression
}
func (e *Expression) Kind() Kind { return e.kind }
func (e *Expression) Value() core.Value { return e.value }
func (e *Expression) Constructor() string { return e.constructor }
func (e *Expression) Operands() []Operand { return slices.Clone(e.operands) }
func (e *Expression) Parameter() (ParameterID, bool) { return e.parameter, e.kind == Parameter }
func (e *Expression) Target() (DeclarationID, bool) { return e.target, e.kind == Application }
// Arguments() returns separately owned child handles.
// Program.Modules(), Module.Declarations(), Declaration.Expression() likewise
// return independent handles; Ontos Value itself is immutable.
```

Program contains the root module digest, canonical module closure, declarations
and limits. Modules expose Pin(), Name(), Value(), Imports(), Declarations(). An
Import exposes Alias(), Name(), Pin(). Declarations expose ID(), Parameters()
(ordered display names), and Expression(). Program.Declaration(ID) performs
identity-preserving lookup. All mutable accessors copy; the model is safe for
concurrent reading/specialization. No generic semantic model exists in bittype yet.

`core/go/declarations/resolve.go:17-111` performs all-or-nothing ResolveWithLimits.
It validates budgets, canonical bytes/pins, transitive imports and expected module
names, declaration references, arity, generic scope, descriptor well-formedness,
recursion boundaries and declaration/import cycle restrictions. It returns no
partial Program. Default limits: depth256, 100,000 visited nodes, 16MiB values,
256 supplied modules. ParameterID includes the defining declaration; nested
application arguments retain caller bindings. The declaration format has only
type parameters, no inference/defaults/value parameters/higher-kinded slots.
All revision-1 declarations are exported. Import aliases and declaration names
occupy separate namespaces, selected by the explicit local/imported target form.

Key actual expression branches (`resolve.go:228-282`):

```go
if args, ok := form(value, "param", 1); ok {
    slot, ok := index(args[0])
    if !ok || slot >= len(owner.parameters) {
        return nil, nil, failure(Invalid, "unbound or out-of-range type parameter")
    }
    return &Expression{kind: Parameter, value: value,
        parameter: ParameterID{Declaration: owner.id, Index: slot}}, descriptorAny(), nil
}
if args, ok := form(value, "apply", 2); ok {
    target, err := r.target(args[0], owner.id)
    if err != nil { return nil, nil, err }
    arguments, err := data.ReadList(args[1])
    if err != nil || len(arguments) != len(r.program.declarations[target].parameters) {
        return nil, nil, failure(Invalid, "application argument list or arity mismatch")
    }
    result := &Expression{kind: Application, value: value, target: target}
    r.edges[owner.id] = append(r.edges[owner.id], target)
    for i, argument := range arguments {
        childPath := append(slices.Clone(path), 2, i+1)
        child, skeleton, err := r.read(argument, owner, childPath, depth+1)
        if err == nil { err = descriptorWF(skeleton) }
        if err != nil { return nil, nil, locate(err, owner.id, childPath) }
        result.arguments = append(result.arguments, child)
    }
    return result, descriptorAny(), nil
}
```

`form` recognizes the exact tagged declaration tuple and arity; `index` decodes
a nonnegative integer; `target` resolves local/imported references in the owner
module. `descriptorAny` substitutes only already-checked parameter/application
positions for a final descriptor-grammar check; it is not the resulting semantic
meaning. `descriptorWF` delegates to existing Horos descriptor rules.

For Descriptor expressions, Horos itself computes which nested tuple positions
contain types. `Operands()` exposes each path and checked child expression; Value()
retains every original literal operand. Bittype can traverse those positions
without another descriptor-constructor switch. A field's name, requiredness,
record open/closed mode, scalar kind and refinement bounds are literal operands.
For example, exact(tuple tagged `horos/decl/param`) is literal data; it is not an
open slot and traversal must not capture it. Existing frontend projection already
uses these operand paths to distinguish type positions from literal fragments.

Horos errors (`model.go:30-58`) are currently coarse:

```go
type Code string
const (
    Invalid Code = "invalid"
    Resolution Code = "resolution"
    Cycle Code = "cycle"
    Exhausted Code = "exhausted"
)
type Error struct {
    Code Code
    Module Digest
    Declaration string
    Path []int // tuple indices within declaration body; empty for module error
    Message string
}
```

Resolution covers both missing supplied pinned module (resolve.go61-64) and
definite name/pin mismatch or unknown alias/declaration (75-77,205-216). Invalid
includes invalid checker limits (21-24) as well as invalid declarations. Module
error locations are not always populated. Exhaustion is not proof of invalidity.
An adapter cannot reliably produce fine verdicts by switching on these broad
codes or parsing English messages. A narrowly scoped structured-error extension
in Horos is possible; it must preserve semantic ownership, not reimplement lookup
in bittype to infer why Horos failed.

### The documented minimal semantic model

`bittype/docs/semantics.md:49-145` is an illustrative proposal, not a frozen API:

```go
type EntityID uint32
type ScopeID uint32
type TermID uint32
type Key struct { Owner string; Name string }
type Model struct {
    Language language.ID
    Profile string
    Snapshot string
    Entities map[EntityID]Entity
    Scopes map[ScopeID]Scope
    Terms map[TermID]Term
    Bindings map[ast.NodeRef][]EntityID
    Judgments []Judgment
    Completed []Key
}
type Entity struct {
    Name string
    Declarations []ast.NodeRef
    DeclaredIn ScopeID
    Meaning TermID
    Parameters []Parameter
    Members []EntityID
}
type Scope struct { Declarations []EntityID; Links []ScopeLink }
type ScopeLink struct { Relation Key; Target ScopeID }
type Parameter struct { Entity EntityID; Sort Key }
type Term struct {
    Head EntityID
    Arguments []TermID
    Substitution map[EntityID]TermID
}
type Judgment struct {
    Relation Key
    Operands []TermID
    At ast.NodeRef
    Verdict Verdict
    Message string
}
type Verdict uint8
const (
    Undetermined Verdict = iota
    Holds
    Rejected
    Conditional
    Unsupported
    Incomplete
)
```

Required intent: IDs are local to one snapshot; zero is absent only where allowed.
Names are display metadata. Open parameters are resolved entities, not unknown
names. Bindings record established targets, never ambiguous candidate sets.
Scope links do not impose universal lexical precedence or acyclicity. A named
alias and its target remain distinguishable. Reference syntax and semantic
identity are distinct. Keys carry provider-owned vocabulary, not auto-discovery
or semantic truth. Completed means analysis completed, not all judgments hold.
Assignability, identity, conversion and constraint satisfaction are different
relations. Unknown/unsupported/incomplete never become true or false implicitly.
Conditional results are not proof without represented assumptions. Input trees
are not mutated. Published snapshots must not expose mutable storage.

Known inadequacies to settle rather than hide:

- Term currently cannot distinguish reference from a zero-argument application.
- Constructor plus type-argument list loses literal operands: scalar int vs text,
  exact("a") vs exact("b"), required vs optional fields, labels and record mode.
- A caller/provider supplied Snapshot string does not bind a result to exact
  input contents. Two distinct ASTs can have the same caller-chosen TreeID.
- A useful shared model must not just wrap an opaque Horos Program, but it also
  must not become a universal Horos type algebra or promise native equivalence.
- Structured premises, constant terms and rich resolution outcomes are not yet
  designed. Native Rust includes lifetimes and associated types; Go has named
  versus alias types, embedded constraints and array lengths; the current
  TypeScript AST includes keyof/indexed access, while mapped types are a future
  counterexample outside its admitted syntax; C++ includes pointer/reference distinctions.
  We need a bounded initial contract and explicit coverage instead of pretending
  the first Horos mapping gives complete semantics for these languages.

### Pagination and required adversarial examples

Explanatory notation (not implemented readable Horos grammar):

```text
pagination.Page<T> = record closed { items: list-of(T), next?: text }
people.Person = record closed { name: text }
people.PeoplePage = pagination.Page<people.Person>
```

The shared graph must distinguish Page, its T, Person and PeoplePage. The
application head is Page, argument is caller-side Person, and substitution maps
Page's parameter identity to Person's term. An imported module's same-spelled
Person and another declaration's T cannot capture those uses. Aliases and open
declarations remain inspectable. No expansion happens simply because a binding
exists. Horos's explicit Program.Specialize(declarationID, groundArguments)
already returns a ground descriptor, preserving the input model.
A ground descriptor contains no declaration parameters or applications and is
closed with respect to Horos recursion references; guarded recursive descriptors
can still be ground.

Adversarial acceptance cases: missing/wrong binding annotation; missing pinned
dependency versus a definitely absent name in a supplied complete module; changed
dependency/options with same TreeID; wrong syntax marker/profile; stale result
from a different request; analyzer mutation/retained slices; output plus error;
budget exhaustion; literals resembling declaration forms; exact-value and record
metadata preservation; zero-argument application; concurrent/reentrant callbacks;
unfamiliar third-party language registration without kernel edits.

### Candidate composition and constraints

Proposed additions are `semantics/go`, `contracts/analysis/go`,
`analyze/horos/go`, an optional engine analyzer slot, and Horos provider composition.
Check remains syntax admission; Analyze establishes semantic facts and validates
any supplied semantic annotations. Horos's profile must change explicitly, not
quietly retain the old resolved-AST guarantee. No compatibility aliases or dual
fallback implementation: migrate existing tests/callers to the clean contract.

Candidate typed interface:

```go
type Provider[L any] interface {
    language.LanguageProvider
    Analyze(context.Context, Request[L]) (Result, error)
}
```

Request would contain owned AST inputs, explicit entry/dependency context, profile
and limits. A typed adapter would enforce marker/profile and own snapshots, bind
results to a request token/catalogue, and expose a heterogeneous analyzer for
RegisterAnalyzer. Existing text.Provider must not acquire a required Analyze
method. One explicitly configured analyzer per registered language is enough
initially; no automatic semantic canonicalization through Horos or route search.

Result could distinguish diagnostics, actual coverage, and an optional complete
model. Infrastructure errors return no usable output. Expected semantic rejection
is an analysis result. Horos is all-or-nothing today; preserving that behavior
is preferable to manufacturing checked partial facts. No result caching or
persistent semantic hashes are required. A per-request internal token is a
possible identity boundary, but cannot itself prove provider claims true.

Implementation is authorized through remote CI and main integration. Native
generators, runtime bindings, function/stream implementations, text grammar for
Horos, compiler-complete native analyzers, publication and deployment are outside
this task. A narrow native semantic qualification is reasonable only if needed
to prevent an accidentally Horos-shaped SPI. The owner forbids local portable
compilation/tests/formatting: maintained GitHub gates run exact committed source,
including native fixtures, race and fuzz tests. No performance claim is established;
do not introduce caches or optimize copying without measurements.

### Prior evidence, alternatives and rejected approaches

Historical research 0001-value-type-core (Nightfall document
doc_a89555cd-e8c1-4fb1-981b-68a17d0bfe8d, run
run_727b1126-6e4d-44a9-9996-16f72909a29a) asked about an independent canonical
value algebra and structural identity. Its evaluated document is retained in Git
91363a5 and advice in cff51c6. The owner subsequently selected Horos and cleaned
the current tree. Reusable lessons include correct binder scope, explicit import
context, opaque literal boundaries and exhaustion distinct from mismatch.
Old proposals to erase aliases, hash expanded applications for semantic identity,
introduce a different numeric/nullability algebra or reject well-formed uninhabited
types on productivity grounds are not current authority. This consultation asks
a different SPI question.

Official reference APIs motivating the common concepts are Go types.Info
(Defs/Uses/Types/Scopes/Instances), TypeScript's checker separation between
syntax/symbol/type, and Rust's distinct type/lifetime/constant parameters. These
are conceptual references, not dependencies or a promise to reproduce them.

Considered options:

1. Kernel delegates Analyze but returns an opaque compiler object: smallest
   wiring, insufficient agnostic inspection or validation.
2. Universal built-in type algebra with language/contract name dispatch: makes
   Horos-to-native equivalence assumptions and violates semantic ownership.
3. Generic fact graph plus provider-owned terms/literals and explicit coverage:
   current direction, but must define losslessness, provenance and partial results.
4. Leave Horos Check as semantic admission and analyze only resolved ASTs: easy
   adapter initially, circular for unresolved authored/migration syntax.
5. Add unchecked AST[any] or cast by matching IDs: breaks existing marker safety.
6. Store arbitrary mutable `any` payloads: breaks owned snapshots and meaningful
   generic integrity checks. A bounded opaque atom versus typed syntax-fragment
   reference is a real design choice, not authorization for arbitrary handles.
7. Automatically infer unsupported native relations by projecting through Horos:
   conflicts with partial, explicitly configured language mappings.

## Questions for the expert

1. What minimal shared semantic graph and vocabulary boundary best supports
   useful agnostic inspection while keeping each language's rules authoritative?
   Please recommend which distinctions must be represented in the first
   implementation and which can remain explicitly unsupported, including
   references/applications and literal/field metadata, using concrete Horos and
   native-language counterexamples.
2. How should input context, typed registration, immutable ownership and result
   provenance compose with the existing checked-envelope kernel? Recommend an
   implementable request/result contract, including which integrity checks the
   kernel can establish and which claims remain provider promises.
3. How should syntax admission, semantic rejection, incomplete context/resource
   exhaustion, coverage and diagnostics be separated during the Horos profile
   migration? What is the smallest sound adaptation of its all-or-nothing resolver
   and coarse error API without duplicating its semantic rules?
4. What bounded implementation sequence and discriminating acceptance tests
   would establish the first production Horos provider and avoid freezing a
   misleading language-agnostic interface? Identify what should be deferred with
   concrete triggers, including whether a small native analyzer is needed now.
5. What important assumption or architectural opportunity have we missed, and
   how would you revise this proposal if starting from these constraints?
