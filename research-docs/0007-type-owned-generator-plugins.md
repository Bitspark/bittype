# Research: Type-owned generator plugins and explicit representation overrides

**ID:** 0007
**Date:** 5 November 2026
**Status:** applied (documented proposal only; implementation explicitly deferred)
**Run-ID:** run_35a886b0-509f-4467-8535-9267d74ba42f
**Owner:** codex-api-retirement-audit
**Issue:** https://github.com/Bitspark/bittype/issues/79
**Reviewed:** https://github.com/Bitspark/bittype/issues/79#issuecomment-6658592424
**Graded:** https://github.com/Bitspark/bittype/issues/79#issuecomment-6659575195

## Question

What is the smallest coherent Go design for composing type-owned native-generation
plugins, treating builtin mappings as ordinary supplied plugins and allowing
explicit consumer overrides, while keeping generic binding, semantic identity,
AST ownership and language-specific syntax correct?

The requested deliverable is a documented proposal: Go interfaces, file layout
and composition strategy. No implementation, runtime bindings, new capability
providers, or execution of code validation workflows is authorized in this task.
The user explicitly invoked bb-research for expert consultation. Its submission
records and response verification belong to that requested consultation; they do
not authorize implementing the advice in production code. Accepted advice will
be incorporated into design documents only.

## Context

### System overview and vocabulary

Bitspark has two independent Go modules relevant here:

- **Horos** owns canonical type declarations, generic binding and semantic
  judgments. Its descriptions describe Ontos values and capability-containing
  types. Ontos is the value/encoding library. Ground values are inert data; live
  communication capabilities can occur in the broader type domain.
- **Bittype** owns language-neutral tree/semantic carriers, provider contracts,
  explicit pipeline composition, and language-specific syntax/transformation
  implementations. Horos does not depend on Bittype.

Bittype is implemented in Go and targets TypeScript, Go, Rust and C++. Its generic
`AST[L]` is an abstract syntax tree whose marker `L` identifies the language at
Go type-checking time. `CST[L]` is a concrete syntax tree retaining syntax and an
AST view. A **provider** implements an independently registered pipeline stage.
**Syntax** names the text/CST/AST dimension; **Files** names the whole-source-base
dimension. SPI means service provider interface, the extension contract. The paired
Files dimension uses `Source[L]` (file organization plus concrete trees and
semantic evidence) and `Semantics[L]` (established provider-owned facts).

An **entity** is a semantic declaration or binder. A **term** describes a reference,
a generic application, or an expression encoded by its owning provider. Numeric
entity/term IDs are addresses within one captured graph, not durable identities.
A **snapshot** captures input trees and context. A **receipt** is an immutable
analysis/resolution result tied to that exact snapshot. Current transformers
consume explicit receipts; they do not resolve or analyze implicitly.

An **ordinary declaration** is a named module declaration, such as Page or a
library's Function. A **builtin constructor** is a form defined by the source
language, such as Horos's scalar or list descriptor. These are separate identity
domains. A Horos **kind**, such as Data or List<Scalar>, is a potentially overlapping
constraint on types/constructors, not a unique tag naming a native representation.

Notation used below: an **open declaration** still has unbound generic
parameters. A **proper type** classifies values, whereas a type constructor accepts
type arguments. The **proper-type top** kind admits every proper type. `D: h.Data`
means a parameter constrained by the imported Horos Data kind. `Byte` is an
illustrative user-defined byte type, not a newly proposed builtin or current
supported generator case. `List<Byte>` and `list<T>` are explanatory notation for
applications of Horos's list-of constructor; `list-of(text)` describes the same
constructor directly. `List<Scalar>` in a kind judgment instead constrains a list's
element type; it does not identify the constructor for dispatch.

A **generation plugin** in this question is compiled Go code supplied through
normal package imports and explicit configuration. It produces target type AST
fragments. It is not a dynamically loaded shared object, runtime value adapter,
parser extension, compiler execution hook, or global registration mechanism.

A **representation profile** says what mapping promises: representing uses
already admitted by Horos is weaker than enforcing every Horos constraint against
arbitrary native callers. Current generators expose separate profile names and
reject the unimplemented enforcement profile; compiler acceptance alone would
not establish equivalence.

### Existing ownership and prior decisions

The user requires Function, Stream, Cell and future contracts to remain ordinary
library declarations with library-owned behavior. Generic trees, the generation
kernel, Wire, hydration and routing must not import their names, select behavior
by a shape fingerprint, or acquire central cases/registrations for them.

Research 0003 already advised exact-declaration representation recipes, explicit
consumer selection, ordered parameter mapping, structured target AST or external
types, and declared representation limits. The current proposed recipe document
selects open declarations only. It has no implemented Go plugin API. The new
request extends the proposal to builtin mappings and explicit overrides as the
normal generation mechanism, not an optional bypass around fixed builtin cases.

Research 0005 replaced a closed Type/Data calculus with general structural kinds.
It retained explicit source semantic obligations and native representation limits.
Its overlapping membership facts must not become a most-specific-kind dispatcher.
Research 0004's superseded closed kind model is not a basis for this design.

### Current concrete flow and limitations

There is one root Bittype Go module. The engine registers one ASTTransformer per
source/target language pair. Four concrete Horos-to-native transformers call one
private source-specific generator with a target `Dialect[L]`. That generator
currently has fixed Text, Boolean, List, Record and Alias operations. It is not
a public plugin framework. Generic parameters and applications survive generation.

The implemented pagination subset maps text, boolean, lists, named closed records,
optional fields, proper generic parameters and ordinary applications. Unsupported
shapes, unknown profiles and failed semantic admission return no usable output.
Integers, arbitrary refinements, higher-order/computed constructors and recursive
declarations are not supported by this profile. Imported type declarations are
flattened with disambiguated native names. Each target returns a one-tree ASTSet.

The generic output boundary already supports multiple target trees. Native ASTs
already contain imports/includes/qualified references and native callable/channel
syntax where appropriate. Target semantic providers and native Files
parsing/formatting are not implemented; output can use existing native Syntax
formatters/printers. This proposal must not invent native semantic receipts.

### Relevant code

The excerpts below are from Bittype main at `7e55f7689d6cf8439b13a3ead6681eae0d2ad9e5`.
They include the types needed to understand the extension boundary; the expert
has no repository access and should not need to follow local links.

`contracts/transform/go/transformer.go:13`

```go
var ErrInvalidTransformer = errors.New("invalid AST transformer")

// Input supplies captured syntax and ordered receipts explicitly. Callers must
// not mutate Evidence during invocation. No provider runs implicitly.
type Input struct {
	Source   semantic.Snapshot
	Evidence []semantic.Result
}

// Contract declares exactly the evidence required by a typed mapping.
type Contract struct {
	Evidence []semantic.EvidenceRequirement
}

// ASTTransformer is a directed, potentially partial AST mapping. Endpoint IDs
// are stable. Implementations return an error for unsupported input, preserve
// necessary provenance, perform no implicit IO, and never return usable partial
// output on failure. Transforming syntax is not proof of semantic equivalence.
type ASTTransformer interface {
	FromID() language.ID
	ToID() language.ID
	Transform(context.Context, Input) (representation.ASTSet, error)
}

// DescribedTransformer additionally advertises typed/profile endpoints so a
// registry can reject incompatibility before execution. Minimal implementations
// are still checked against registered endpoints at every execution boundary.
type DescribedTransformer interface {
	ASTTransformer
	FromLanguage() representation.Language
	ToLanguage() representation.Language
}
```

`contracts/transform/go/adapter.go:14`

```go
type Request[L any] struct {
	source   ast.AST[L]
	entry    ast.NodeRef
	evidence []semantic.Result
	inputs   []representation.AST
}

func (r Request[L]) Source() ast.AST[L] { return tree.Clone(r.source) }
func (r Request[L]) Entry() ast.NodeRef { return r.entry }

// Inputs retains the complete captured source collection for multi-tree mappings.
func (r Request[L]) Inputs() []representation.AST { return slices.Clone(r.inputs) }
func (r Request[L]) Evidence(index int) (semantic.Result, bool) {
	if index < 0 || index >= len(r.evidence) {
		return semantic.Result{}, false
	}
	return r.evidence[index], true
}

// NewTransformer captures requirements and shares semantic receipt admission.
// Endpoint syntax admission remains the engine's checker responsibility.
func NewTransformer[F, T any](from, to representation.Language, contract Contract, stage func(context.Context, Request[F]) (ast.Set[T], error)) (ASTTransformer, error) {
	if stage == nil {
		return nil, ErrInvalidTransformer
	}
	if err := representation.CheckMarker[F](from); err != nil {
		return nil, err
	}
	if err := representation.CheckMarker[T](to); err != nil {
		return nil, err
	}
	if err := semantic.ValidateEvidenceRequirements(contract.Evidence); err != nil {
		return nil, err
	}
	requirements := make([]semantic.EvidenceRequirement, len(contract.Evidence))
	for i, requirement := range contract.Evidence {
		requirements[i] = requirement.Clone()
	}
	return bound[F, T]{from: from, to: to, requirements: requirements, stage: stage}, nil
}

type bound[F, T any] struct {
	from, to     representation.Language
	requirements []semantic.EvidenceRequirement
	stage        func(context.Context, Request[F]) (ast.Set[T], error)
}

func (b bound[F, T]) FromID() language.ID                   { return b.from.ID() }
func (b bound[F, T]) ToID() language.ID                     { return b.to.ID() }
func (b bound[F, T]) FromLanguage() representation.Language { return b.from }
func (b bound[F, T]) ToLanguage() representation.Language   { return b.to }
```

`semantics/go/model.go:19`

```go
// IDs are one-based addresses local to one Semantics[L], never persistent identities.
type EntityID uint32
type TermID uint32

// Key names vocabulary owned by a provider. It causes no registration or lookup.
type Key struct{ Owner, Name string }

func (k Key) Valid() bool {
	return strings.TrimSpace(k.Owner) != "" && strings.TrimSpace(k.Name) != "" && utf8.ValidString(k.Owner) && utf8.ValidString(k.Name)
}

type Entity struct {
	Kind         Key
	Name         string
	Owner        EntityID
	Declarations []ast.NodeRef
	Parameters   []Parameter
	Body         TermID
}

// Parameter retains a real binder and its coarse provider-owned sort. Sort is
// not a structured kind; providers may associate richer facts with the entity.
type Parameter struct {
	Entity EntityID
	Sort   Key
}

type TermKind uint8

const (
	Reference TermKind = iota + 1
	Application
	ProviderExpression
)

// Term has exactly one variant. Application arguments are effective positional
// arguments, not a mixture of written, inferred, defaulted or packed arguments.
// Reference owns Target; Application owns Callee and Arguments. Applicability,
// arity and saturation belong to the provider, not graph admission.
// A provider expression's encoding defines its bytes and operand locators.
type Term struct {
	Kind      TermKind
	Target    EntityID
	Callee    TermID
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

// Occurrence identifies a semantic position within a current syntax node.
// Role and Ordinal are provider-owned; an occurrence is not a candidate set.
type Occurrence struct {
	Site    ast.NodeRef
	Role    Key
	Ordinal uint32
}

type Binding struct {
	At     Occurrence
	Target EntityID
}
```

`trees/ast/go/ast.go:1`

```go
// Package ast defines language-parameterized syntax carriers.
//
// The language parameter separates ASTs and node kinds in Go's type system.
// A language package supplies its own named marker type, kinds and grammar;
// this package contains no language catalogue or grammar implementation.
//
// These are data types, not checked trees. Structural and grammatical checks
// live in separate packages; the carrier performs no validation or allocation.
package ast

// TreeID identifies a particular AST within a set of inputs and outputs.
// It is an opaque provenance/addressing key, not a declaration or content hash.
type TreeID string

// NodeID is a one-based index into an AST's Nodes slice. Zero names no node.
// IDs are local to an AST; use NodeRef when crossing tree boundaries.
type NodeID uint32

// AST holds one syntax tree in a flat node arena.
//
// L is a language-owned named marker type. It is not a runtime registration.
// A populated tree has a nonempty ID and a Root in the range 1..len(Nodes).
// Node i+1 is stored at Nodes[i]. The zero value is an empty carrier.
//
// Children describe syntax containment; References describe resolved links.
// Recursive declarations use references rather than containment cycles.
// Fields remain mutable: callers must treat stage inputs as read-only and
// create separate output storage. No immutability or validity is enforced here.
//
// bittype:dimension Syntax
type AST[L any] struct {
	ID    TreeID
	Root  NodeID
	Nodes []Node[L]
}
```

`trees/ast/go/node.go:1`

```go
package ast

// Kind is a language-owned syntax discriminator. The same spelling in two
// different named language markers does not give the kinds the same Go type.
// There are no reserved kinds or built-in language/contract names.
type Kind[L any] string

// Attribute is an uninterpreted textual syntax atom.
//
// A grammar defines the name, allowed text and meaning. Text can preserve an
// identifier or an exact literal spelling without rounding or normalization.
// It is not an escape hatch for executable source or a runtime value codec.
type Attribute struct {
	Name string
	Text string
}

// Child is an ordered containment edge to another node in the same AST.
// The grammar owns Role names and cardinality; repeated roles can form a list.
type Child struct {
	Role string
	Node NodeID
}

// Reference is a resolved, non-containment link.
// A grammar owns Role and the allowed target. Unresolved spellings belong in
// syntax attributes, not in a fabricated target with a missing node.
type Reference struct {
	Role   string
	Target NodeRef
}

// Node contains language-owned syntax and uninterpreted metadata.
//
// Attributes and Children preserve their supplied order. A grammar decides
// which order matters and whether duplicate attributes or roles are allowed.
// No grammar, name resolution or normalization is performed by these fields.
type Node[L any] struct {
	Kind       Kind[L]
	Attributes []Attribute
	Children   []Child
	References []Reference
	Origin     Origin
}
```

`trees/ast/go/set.go:1`

```go
package ast

// Set is construction storage for one language's complete multi-unit output.
// Entry is explicit, even for a one-tree set. It carries no file organization.
type Set[L any] struct {
	Entry NodeRef
	Trees []AST[L]
}
```

`internal/horosgen/generate.go:104`

```go
type Field struct {
	Name     string
	Optional bool
	Type     ast.NodeID
	Origin   ast.NodeRef
}
type Dialect[L any] interface {
	ReservedNames() []string
	Parameter(*Arena[L], string, ast.NodeRef) ast.NodeID
	Reference(*Arena[L], string, ast.NodeID, []ast.NodeID, ast.NodeRef) ast.NodeID
	Text(*Arena[L], ast.NodeRef) ast.NodeID
	Boolean(*Arena[L], ast.NodeRef) ast.NodeID
	List(*Arena[L], ast.NodeID, ast.NodeRef) ast.NodeID
	Record(*Arena[L], ast.NodeID, string, []ast.NodeID, []Field)
	Alias(*Arena[L], ast.NodeID, string, []ast.NodeID, ast.NodeID)
	Finish(*Arena[L], []ast.NodeID) ast.NodeID
	FieldName(string) (string, error)
}

```

`internal/horosgen/generate.go:550`

```go
func (g *generator[L]) expression(id semantics.TermID, depth int) (ast.NodeID, error) {
	if err := g.tick(depth); err != nil {
		return 0, err
	}
	term := g.facts.Terms[id-1]
	origin := g.sites[id]
	if term.Kind == semantics.Reference || term.Kind == semantics.Application {
		entity, err := g.referenceTarget(id)
		if err != nil {
			return 0, err
		}
		g.used[entity] = true
		var arguments []ast.NodeID
		for _, argument := range term.Arguments {
			mapped, err := g.expression(argument, depth+1)
			if err != nil {
				return 0, err
			}
			arguments = append(arguments, mapped)
		}
		return g.dialect.Reference(&g.arena, g.names[entity], g.targets[entity], arguments, origin), g.arena.Err
	}
	if term.Operator.Owner != "horos" {
		return 0, g.unsupported(id, "unrecognized semantic operator")
	}
	switch term.Operator.Name {
	case "descriptor/scalar":
		tuple, err := g.payload(id)
		if err != nil {
			return 0, err
		}
		switch atom(tuple.At(1)) {
		case "utf8-text":
			return g.dialect.Text(&g.arena, origin), g.arena.Err
		case "bool":
			return g.dialect.Boolean(&g.arena, origin), g.arena.Err
		}
	case "descriptor/list-of":
		element, err := g.operand(id, 1)
		if err != nil {
			return 0, err
		}
		mapped, err := g.expression(element, depth+1)
		if err != nil {
			return 0, err
		}
		return g.dialect.List(&g.arena, mapped, origin), g.arena.Err
	}
	return 0, g.unsupported(id, "pagination mapping does not represent "+term.Operator.Name+" here; nested records must be named declarations")
}
```

`internal/horosmodel/facts.go:73`

```go
func (e *exporter[L]) identity(subject semantics.Subject, name string, value core.Value) error {
	term, err := e.provider("identity/"+name, value, nil, 1)
	if err != nil { return err }
	return e.fact(subject, e.frontend.IdentityRelation(), term)
}
func (e *exporter[L]) moduleIdentity(id semantics.EntityID, pin decl.Digest) error {
	return e.identity(semantics.Subject{Entity: id}, "module", core.NewAtom(pin[:]))
}
func (e *exporter[L]) declarationIdentity(id semantics.EntityID, value decl.DeclarationID) error {
	return e.identity(semantics.Subject{Entity: id}, "declaration", core.NewTuple(core.NewAtom(value.Module[:]), data.MustEncodeText(value.Name)))
}
func (e *exporter[L]) parameterIdentity(id semantics.EntityID, value decl.ParameterID) error {
	return e.identity(semantics.Subject{Entity: id}, "parameter", core.NewTuple(core.NewAtom(value.Declaration.Module[:]), data.MustEncodeText(value.Declaration.Name), data.EncodeInt(big.NewInt(int64(value.Index)))))
}
```

`semantics/go/model.go:500`

```go
type ArgumentBinding struct {
	Parameter EntityID
	Argument  TermID
}

// Facts returns the matching fact terms in provider order. Absence is not a
// negative judgment; the provider's reports determine established coverage.
func (m Semantics[L]) Facts(subject Subject, relation Key) []TermID {
	var result []TermID
	if m.Valid() {
		for _, fact := range m.state.Facts {
			if fact.Subject == subject && fact.Relation == relation {
				result = append(result, fact.Term)
			}
		}
	}
	return result
}

// DirectParameterBindings pairs the real parameters of a direct Reference
// callee with an equally sized argument list. Parameter-headed applications,
// other callees and arity mismatches have no such view. This structural query
// does not establish applicability.
// A direct zero-parameter application returns an empty successful view.
func (m Semantics[L]) DirectParameterBindings(id TermID) ([]ArgumentBinding, bool) {
	if !m.Valid() || id == 0 || uint64(id) > uint64(len(m.state.Terms)) {
		return nil, false
	}
	term := m.state.Terms[id-1]
	if term.Kind != Application {
		return nil, false
	}
	callee := m.state.Terms[term.Callee-1]
	if callee.Kind != Reference {
		return nil, false
	}
	entity := m.state.Entities[callee.Target-1]
	if entity.Owner != 0 {
		for _, parameter := range m.state.Entities[entity.Owner-1].Parameters {
			if parameter.Entity == callee.Target {
				return nil, false
			}
		}
	}
	parameters := entity.Parameters
	if len(parameters) != len(term.Arguments) {
		return nil, false
	}
	result := make([]ArgumentBinding, len(parameters))
	for i, parameter := range parameters {
		result[i] = ArgumentBinding{Parameter: parameter.Entity, Argument: term.Arguments[i]}
	}
	return result, true
}

// View exposes immutable semantics in the same language.
//
// bittype:dimension Files
type View[L any] interface{ SemanticsView() Semantics[L] }
```

`transforms/horos-to-go/go/transform.go:15`

```go
// RepresentationProfile permits kind-constraint erasure for Horos-checked uses.
const RepresentationProfile = g.RepresentationProfile

// EnforcementProfile requires native kind enforcement and is not implemented.
const EnforcementProfile = g.EnforcementProfile

type Config struct {
	TreeID  ast.TreeID
	Package string
	Profile string
}

func New(from, to representation.Language, config Config) (transform.ASTTransformer, error) {
	if from.ID() != (h.Language{}).ID() || to.ID() != (s.Language{}).ID() {
		return nil, transform.ErrInvalidTransformer
	}
	if from.Profile() != h.Profile || config.TreeID == "" || !g.KnownProfile(config.Profile) || !g.Identifier(config.Package) || token.Lookup(config.Package).IsKeyword() {
		return nil, transform.ErrInvalidTransformer
	}
	return transform.NewTransformer[h.Language, s.Language](from, to, g.Contract(), func(ctx context.Context, request transform.Request[h.Language]) (ast.Set[s.Language], error) {
		output, err := g.Generate(ctx, request, config.TreeID, config.Profile, &dialect{pkg: config.Package})
		if err != nil {
			return ast.Set[s.Language]{}, err
		}
		return ast.Set[s.Language]{Entry: ast.NodeRef{Tree: output.ID, Node: output.Root}, Trees: []ast.AST[s.Language]{output}}, nil
	})
}

type dialect struct {
	pkg      string
	optional ast.NodeID
}

func (*dialect) ReservedNames() []string { return []string{"BittypeOptional"} }

func (*dialect) FieldName(name string) (string, error) { return g.Exported(name) }
func (*dialect) Parameter(a *g.Arena[s.Language], name string, o ast.NodeRef) ast.NodeID {
	constraint := a.Add(s.KindInterfaceType, nil, o)
	return a.Add(s.KindTypeParameter, g.Attr("name", name), o, g.Child(s.RoleConstraint, constraint))
}
func (*dialect) Reference(a *g.Arena[s.Language], name string, target ast.NodeID, args []ast.NodeID, o ast.NodeRef) ast.NodeID {
	ref := a.Add(s.KindNamedType, g.Attr("name", name), o)
	a.Target(ref, target)
	if len(args) == 0 {
		return ref
	}
	return a.Add(s.KindTypeApplication, nil, o, append([]ast.Child{g.Child(s.RoleConstructor, ref)}, g.Children(s.RoleArgument, args)...)...)
}
func (*dialect) Text(a *g.Arena[s.Language], o ast.NodeRef) ast.NodeID {
	return a.Add(s.KindNamedType, g.Attr("name", "string"), o)
}
func (*dialect) Boolean(a *g.Arena[s.Language], o ast.NodeRef) ast.NodeID {
	return a.Add(s.KindNamedType, g.Attr("name", "bool"), o)
}
func (*dialect) List(a *g.Arena[s.Language], element ast.NodeID, o ast.NodeRef) ast.NodeID {
	return a.Add(s.KindSliceType, nil, o, g.Child(s.RoleElement, element))
}
func (d *dialect) presence(a *g.Arena[s.Language]) ast.NodeID {
	if d.optional != 0 {
		return d.optional
	}
	d.optional = a.Add(s.KindTypeDeclaration, nil, ast.NodeRef{})
	p := d.Parameter(a, "T", ast.NodeRef{})
	value := d.Reference(a, "T", p, nil, ast.NodeRef{})
	boolean := d.Boolean(a, ast.NodeRef{})
	d.Record(a, d.optional, "BittypeOptional", []ast.NodeID{p}, []g.Field{{Name: "Present", Type: boolean}, {Name: "Value", Type: value}})
	return d.optional
}
func (d *dialect) Record(a *g.Arena[s.Language], id ast.NodeID, name string, params []ast.NodeID, fields []g.Field) {
	var members []ast.Child
	for _, field := range fields {
		ty := field.Type
		if field.Optional {
			ty = d.Reference(a, "BittypeOptional", d.presence(a), []ast.NodeID{ty}, field.Origin)
		}
		member := a.Add(s.KindField, g.Attr("name", field.Name), field.Origin, g.Child(s.RoleType, ty))
		members = append(members, g.Child(s.RoleField, member))
	}
	body := a.Add(s.KindStructType, nil, a.Origin(id), members...)
	a.Set(id, s.KindTypeDeclaration, g.Attr("name", name), append(g.Children(s.RoleTypeParameter, params), g.Child(s.RoleType, body))...)
}
func (*dialect) Alias(a *g.Arena[s.Language], id ast.NodeID, name string, params []ast.NodeID, body ast.NodeID) {
	a.Set(id, s.KindTypeAlias, g.Attr("name", name), append(g.Children(s.RoleTypeParameter, params), g.Child(s.RoleType, body))...)
}
func (d *dialect) Finish(a *g.Arena[s.Language], decls []ast.NodeID) ast.NodeID {
	if d.optional != 0 {
		decls = append([]ast.NodeID{d.optional}, decls...)
	}
	return a.Add(s.KindFile, g.Attr("package", d.pkg), ast.NodeRef{}, g.Children(s.RoleDeclaration, decls)...)
}
```

The following small existing carriers make graph and tree ownership explicit.
`Subject` selects exactly one entity or term. Neither address is a durable source
identity.

```go
// semantics/go/model.go
type Subject struct {
    Entity EntityID
    Term   TermID
}
// trees/ast/go/origin.go
type NodeRef struct {
    Tree TreeID
    Node NodeID
}
```

The current private target arena owns one tree. `Add` appends a fresh containment
node; `Set` fills a reserved node; `Target` records a reference within that tree.
It is shown as current implementation context, not as the proposed public API.

```go
type Arena[L any] struct {
	Tree ast.AST[L]
	Ctx  context.Context
	Err  error
}

func (a *Arena[L]) Add(kind ast.Kind[L], attributes []ast.Attribute, origin ast.NodeRef, children ...ast.Child) ast.NodeID {
	if a.Err != nil {
		return 0
	}
	if err := a.Ctx.Err(); err != nil {
		a.Err = err
		return 0
	}
	if len(a.Tree.Nodes) >= MaxNodes {
		a.Err = &text.Diagnostic{Code: text.LimitExceeded, Message: "native generation node limit", Node: &origin}
		return 0
	}
	n := ast.Node[L]{Kind: kind, Attributes: attributes, Children: children}
	if origin.Tree != "" {
		n.Origin.Nodes = []ast.NodeRef{origin}
	}
	a.Tree.Nodes = append(a.Tree.Nodes, n)
	return ast.NodeID(len(a.Tree.Nodes))
}
func (a *Arena[L]) Set(id ast.NodeID, kind ast.Kind[L], attributes []ast.Attribute, children ...ast.Child) {
	if a.Err != nil {
		return
	}
	n := &a.Tree.Nodes[id-1]
	n.Kind = kind
	n.Attributes = attributes
	n.Children = children
}
func (a *Arena[L]) Target(id, target ast.NodeID) {
	if a.Err == nil {
		a.Tree.Nodes[id-1].References = []ast.Reference{{Role: "target", Target: ast.NodeRef{Tree: a.Tree.ID, Node: target}}}
	}
}
```

### Data and control flow

Current generation is:

1. An application supplies source modules and explicitly registers languages,
   checkers, resolver/analyzer and one configured transformer per language pair.
2. Capture, resolution and analysis establish a snapshot and compatible receipts.
3. NewTransformer admits the receipts and gives the concrete mapping a typed
   request containing the full source collection.
4. Horos's reader interprets source-owned operator encodings, identity facts,
   parameter scopes and published type/kind facts. Generic core code interprets none.
5. `prepare` reserves target declaration placeholders and generic parameter nodes.
   The current traversal rejects cycles and maps expressions through the fixed
   Dialect method set and descriptor switch.
6. The target dialect writes one AST arena. The outer adapter packs the complete
   ASTSet. The engine admits/captures outputs. Formatting and printing are separate.

Current AST containment forbids multiple parents, orphan nodes and containment
cycles. Reference edges can point to declarations in other supplied target trees.
A reusable lowered type cannot therefore be one reusable raw `ast.NodeID` attached
at several sites. A plugin must produce owned fragments or request fresh
containment occurrences; declaration references can share symbolic targets.

Current limits include a 100,000-node/128-depth private generation budget and a
1-to-256 tree output-set boundary. Plugins must participate in bounded work;
there is no measured need for parallel plugin execution in the first increment.

## Proposed starting design (submitted candidate, superseded by the evaluation)

The candidate below records what was submitted. The evaluated design in
[the public proposal](../docs/type-generation-plugins.md) supersedes its virtual
structural-default replacement rule, Mode enum and incomplete planning surface.
It is retained as consultation context, not a competing current specification.

This is a candidate to critique, not an implementation commitment. The existing
ASTTransformer and engine edge model stay unchanged. One explicitly configured
transformer composes a source adapter, a target assembler and supplied mappings.

### Identity and deterministic selection

A source adapter owns durable selector construction and current-use resolution.
It can expose Horos-specific helpers outside the generic contract:

- `BuiltinScalar("utf8-text")`, distinct from bool and int;
- `BuiltinConstructor("list-of")`;
- `Declaration(modulePin, declarationName)`.

Here `modulePin` is Horos's `decl.Digest` identifying the exact canonical module
value. `Declaration(modulePin, name)` corresponds to existing `decl.DeclarationID`
with fields Module and Name. A parameter has `decl.ParameterID` = (declaration
identity, ordered zero-based slot). These differ from graph-local EntityID/TermID
and syntax NodeRef. A renamed declaration or changed pinned module has a new
selector; configuration does not follow it by name resemblance.

The generic package compares opaque keys bound to the source language/profile;
it does not decode Horos tags, module pins or builtin names. Source-owned helper
names above are illustrative. The actual Horos descriptor/scalar operator alone
is insufficient to distinguish text from boolean; literal discriminator bytes
must be interpreted by its source adapter. A broad scalar-family rule could be
expanded into exact keys by that adapter; it is not a priority predicate in the
core. Literal payloads of exact descriptors remain opaque data.

Configuration provides one defaults set and one explicit overrides set. Defaults
are ordinary packages imported by application composition. There are no default
plugins hidden inside the engine or generic generation package.

- Duplicate selectors in either layer are errors, including two conflicting
  dependency-supplied defaults. Dependency order never chooses a winner.
- Each override replaces exactly one effective default at the same key. For an
  ordinary declaration that default can be the explicitly configured Structural
  policy; this is still an exact-declaration override, not a category predicate.
  A proposed optional
  `Replaces` plugin ID guards against replacing a different default after upgrades.
- Explicit library bindings are additions to the assembled defaults. An override
  of a builtin with no supplied default is an error. The Source adapter can
  distinguish an exact declaration eligible for Structural from a missing builtin.
  Whether this effective-default distinction needs a different API is a question.
- Select once before emission; the effective table is immutable per transformer.
- A selected mapping's error or Unsupported is terminal. No retry of the replaced
  default, structural mapping, or lower-priority candidate is allowed.
- A mapping can recursively request child mapping under the same effective table.
  That recursion is deliberate composition, not fallback.
- An ordinary declaration with no custom mapping may use the explicitly selected
  structural-declaration policy. Its body then composes builtin mappings. An
  explicit declaration mapping is selected before expanding its body.
- Aliases retain their own identities and constraints. A configured alias mapping
  applies at that boundary; otherwise its body follows normal declared references.
  A lookalike declaration never inherits a library's mapping by structural shape.

Actual kind membership can establish a mapping's eligibility under captured
assumptions. It is not a unique dispatch key: list-of(text) can satisfy Data,
List<Scalar> and proper-type top simultaneously.

### Builtins, ordinary types and overrides

A consumer or Horos-owned integration module supplies default composition. It
imports a Horos-owned integration package
providing scalar/list/record mappings. It explicitly supplies those mappings to
the same binding mechanism used by a Function or Cell library. Horos's semantic
module remains independent of Bittype; the integration module can depend on both.

For example, explicit configuration could select:

```text
Default: Horos utf8-text -> Go string plugin
Default: Horos bool      -> Go bool plugin
Default: Horos list-of   -> Go slice plugin
Default: Horos closed-record descriptor -> Go record-body plugin
Structural declaration policy: reserve a named declaration and map its body
Default: Function pin+name -> Function library's Go callable plugin
Override: Function pin+name, replacing callable plugin -> invocation-handle plugin
```

This is illustrative configuration, not a claim those Function plugins exist.
The same source can select a TypeScript callable, Go func, Rust library-owned
callable handle or C++ library-owned type. Builtin/native AST syntax still belongs
to the corresponding language syntax provider. The mapping plugin chooses and
constructs that syntax; it neither modifies the formatter nor implements calls.

A builtin mapping may itself generate a wrapper although there is no named source
declaration to put in DefineRequest.Input. The proposed protocol must identify
and reserve such mapping-owned declarations separately (for example a helper key
scoped to plugin configuration), including their generic parameters. This is a
concrete reason to examine Plan instead of assuming Definer is sufficient.

A text-to-custom-wrapper override must provide a valid target declaration or a
locked external symbol and preserve the intended text domain. An unbounded Horos
integer cannot become a bounded native integer with a claim of full fidelity.
The representation profile must state accepted inputs and retained obligations.
Optional presence cannot silently become a zero-value sentinel. No protocol
behavior follows from native type syntax.

The simple path is selection -> plugin -> target fragment -> assembly. For
`list-of(text)`, the source adapter selects the list builtin key; that plugin asks
Context.Lower for its element under the same effective configuration. The selected
text plugin supplies a Go string fragment (or the explicit text override supplies
a wrapper reference); the list plugin returns a slice fragment. The target
assembler gives each occurrence its own contained nodes and resolves symbols.

### Candidate Go interfaces

All declarations in this section are proposed. Private fields mean opaque,
invocation-bound handles with copying accessors, not missing implementation code.
`F` and `T` are existing source and target language marker types. The contract
would live at `contracts/generation/go`, with package name `generation`.

```go
// Existing types imported from contracts/semantic, contracts/transform,
// semantics, trees/ast and trees/representation are reused.
type PluginID string

// Source-owned canonical selector bytes plus a source/profile domain.
// This is not a syntax name, local entity ID or semantic equivalence claim.
type Selector[F any] struct { /* private domain and exact identity */ }

// Captured semantic subject plus its lexical/substitution environment.
// Model()/Subject() expose the existing immutable graph and semantic.Subject.
// Direct argument access retains real binder-to-argument pairs.
type Ref[F any] struct { /* private invocation, subject and bindings */ }

type Mode uint8
const (
    Inline Mode = iota + 1 // use-site target type expression
    External              // reference a supplied native type artifact
    Owned                 // generate native declarations once
)

type Binding[F, T any] struct {
    Source Selector[F]
    Plugin Plugin[F, T]
}
type Override[F, T any] struct {
    Source   Selector[F]
    Replaces PluginID
    With     Plugin[F, T]
}
type Config[F, T any] struct {
    Source     Source[F]
    Target     Target[T]
    Defaults   []Binding[F, T]
    Overrides  []Override[F, T]
    Structural Plugin[F, T] // explicit ordinary-declaration policy
}

// A target-typed immutable expression/declaration fragment and symbolic target.
// These own storage; they are not raw arena indices or source-code strings.
type Type[T any] struct { /* private fragment */ }
type Declarations[T any] struct { /* private declaration fragments */ }
type Symbol[T any] struct { /* private session-owned declaration/import symbol */ }

type UseRequest[F, T any] struct {
    Input   Ref[F]
    Context Context[F, T]
}
type DefineRequest[F, T any] struct {
    Input      Ref[F]       // selected source declaration
    Symbol     Symbol[T]    // reserved before body emission
    Parameters []Symbol[T]  // ordered target binder symbols
    Context    Context[F, T]
}

type Plugin[F, T any] interface {
    ID() PluginID
    Mode() Mode
    Use(context.Context, UseRequest[F, T]) (Type[T], error)
}
type Definer[F, T any] interface {
    Plugin[F, T]
    Define(context.Context, DefineRequest[F, T]) (Declarations[T], error)
}

type Context[F, T any] interface {
    Lower(context.Context, Ref[F]) (Type[T], error)
    Target() Target[T]
    // Source argument/binder access goes through Ref and the source adapter;
    // the context does not expose a mutable source model or registry.
}

// Source/Target are collaborators, not new registered engine providers.
// Their complete small method sets are a principal consultation question.
type Source[F any] interface {
    Contract() transform.Contract
    Open(transform.Request[F]) (SourceSession[F], error)
}
type SourceSession[F any] interface {
    Declarations() []Ref[F]
    Selector(Ref[F]) (Selector[F], error)
    Model() semantics.Semantics[F]
}
type Target[T any] interface {
    // Supplies typed fragment construction, symbol/import allocation and
    // final ast.Set[T] assembly; no source language or canonical type cases.
}

// Proposed concrete composition factory in generation/go, not a new engine SPI.
func NewTransformer[F, T any](
    from, to representation.Language,
    config Config[F, T],
) (transform.ASTTransformer, error)
```

`Mode` is a storage/ownership category, not a list of source/native types. An
Owned mapping must implement Definer. External mappings produce references and
explicit imports; they must not fabricate native declarations for external types.
Inline mappings produce complete native type expressions. Native generic syntax
cannot be appended uniformly to every output expression: the plugin must decide
whether to reference a generic symbol, expand a callable type or emit other
valid target syntax.

The candidate needs a small, explicit target fragment and planning surface:

- build target-owned kinds/attributes/children using scoped handles;
- bind references to reserved local/external symbols without raw NodeRef creation;
- request dependency imports from caller-supplied artifacts (no retrieval);
- allocate names, lexical scopes and output units deterministically;
- allow declared reusable helpers, such as Go's presence wrapper, without a
  builtin-name switch or inline recursive allocation during body emission;
- finish all units and remap symbolic refs to ordinary existing AST addresses.

Alternative: a single `Generate(Request) -> ast.NodeID` callback with a public
mutable target Arena is smaller, but risks cross-session indices, shared
containment and mixing definition planning with expression emission. Another
alternative is a mandatory `Plan` method before Use/Define; we ask whether that
is justified for helpers, imported artifacts and recursive ownership, rather
than adding speculative public interfaces.

The plugin closure/configuration is treated as immutable for one transformer.
The host owns cancellation, common limits, symbol tables and final output.
Arbitrary Go plugin code is trusted code, not sandboxed by these interfaces.
No cross-invocation cache is proposed; model-local addresses cannot be cache keys.

### Generics, application overrides and recursion

Default and declaration overrides apply coherently to open declarations and
ordered applications. `Page<T>` is still emitted once; `DataPage<D: h.Data>`
remains an ordinary constrained declaration, not a copied Page.

Concrete-application overrides, such as `List<Byte> -> ByteBuffer`, are useful but
are not covered by open-declaration keys. They require source-owned durable
application identity and a representation-family/specialization policy. The
counterexample is:

```text
type Page<T> = { items: list<T>; };
type Bytes = Page<Byte>;
```

Generating native Page<T> with a generic collection does not make its field become
ByteBuffer when instantiated as Page<Byte>. Applying the override only at direct
List<Byte> uses would violate this requirement: lowering after substitution must
agree with substituting into the lowered generic representation.
An explicit specialization profile could handle that; silent monomorphization or
silently ignoring the override would violate current promises. Candidate initial
scope therefore admits exact builtin/declaration selectors and defers application
selectors until this coherence policy is specified. Higher-order/computed
constructors remain a separately supported source case, not automatically enabled.

Planning must select representations before rejecting cycles: a recursive source
can map to a supplied external handle. Conversely, reserving symbols does not make
an infinitely sized Rust struct or invalid native alias cycle legal. A selected
representation must support the recursion, introduce its declared indirection,
or return Unsupported. No generic kernel guesses pointer or ownership semantics.

### Target file layout and dependency direction

Proposed paths only; no source directories are created by this documentation task:

```text
bittype/
  contracts/generation/go/            shared typed plugin/configuration contracts
    selector.go request.go plugin.go fragment.go
  generation/go/                      selection, bounded composition and assembly
    bindings.go planning.go transform.go
  generation/source/horos/go/          Horos semantic reading and exact selector helpers
    source.go identity.go
  generation/target/<language>/go/     target syntax, symbols/imports and units
    fragments.go symbols.go assemble.go
  transforms/horos-to-<target>/go/     configuration-only ASTTransformer factories
  docs/type-generation-plugins.md      proposed public design

horos-owned integration module/        independent of Horos's core module
  generation/syntax/<language>/go/     builtin mappings owned alongside Horos integration
    scalar.go list.go record.go
  transforms/horos-to-<target>/go/     opinionated defaults composition

function library/                     illustrative, no new repo required by this task
  generation/syntax/<language>/go/     Function-specific native type plugins
cell library/
  generation/syntax/<language>/go/     optional Cell-specific native type plugins
consuming application/
  generation/                         selected defaults, libraries and overrides
```

Bittype factories take supplied bindings and do not import the integration module.
Opinionated assembly lives in that integration module or the consuming application,
so there is no Bittype -> integration -> Bittype module dependency cycle.

The builtin integration module is source-owner-owned code, not a mandatory new
repository or a backdependency from Horos core to Bittype. Each target language
uses the same package layout, with Go as the implementation language. Type-owned
plugins can be separate Go integration modules within their library repositories.
Composition changes when a consumer selects a new library; generic generator,
ASTTransformer, engine, parser, formatter, Wire and runtime do not change.

Dependencies: consumer -> selected plugins + generic generation + source/target
adapters; plugin -> shared contracts + its source semantics + target syntax;
generic generation -> shared contracts/carriers only; source adapter -> Horos
frontend; target adapter -> own native syntax. The engine never imports these
implementations. There is no dynamic discovery, init registration, global
capability table or default library catalogue.

### Prior art and limits of analogy

- MLIR Dialect Conversion's TypeConverter separates type mapping hooks from
  value materialization. Its reverse-registration-order conversions and value
  operations are not proposed here. Context affects cacheability.
  https://mlir.llvm.org/docs/DialectConversion/#type-converter
- OpenAPI Generator's schema/import mappings let consumers reference existing
  model types. Its name-based configuration and generator-specific support do not
  establish canonical identity or uniform semantic plugins.
  https://openapi-generator.tech/docs/customization/#bringing-your-own-models
- Rust coherence illustrates why overlapping implementations need a clear policy.
  Bittype requires consumer-scoped explicit overrides rather than Rust's global
  orphan rules, and therefore should not copy those restrictions.
  https://doc.rust-lang.org/reference/items/implementations.html#trait-implementation-coherence

## Questions for the Expert

1. What selector and override model best supports builtin mappings and exact
   library declarations while preserving source ownership, alias identity and
   deterministic consumer policy? Is defaults-plus-explicit-replacement sufficient?
2. What is the smallest useful Go plugin contract and target fragment/symbol
   surface for composing type uses and owned/external declarations? Please give
   concrete signatures, including complete Source/SourceSession contracts, their
   supporting types and package/file boundaries.
3. How should open generic coherence constrain application-specific overrides,
   source kinds and representation eligibility? Which boundaries should the first
   proposal explicitly support or defer?
4. What planning and ownership protocol should handle helpers, imports, repeated
   type occurrences, recursive declarations and multiple AST output units while
   preserving target ownership and explicit native semantic choices? Include a
   builtin text mapping that emits its own wrapper without a source declaration.
5. What important simplification, failure case or alternative architecture have
   we missed, given the existing ASTTransformer, semantic evidence and types-only
   scope? Please challenge assumptions rather than optimizing the candidate blindly.


## Applied

Evaluation scope is documentation only, as explicitly requested by the owner.
The consulted baseline is Bittype `7e55f7689d6cf8439b13a3ead6681eae0d2ad9e5`.
Attempt 1 submitted document `doc_4c3de827-b0ed-45cc-8870-6c6b892772f0` to run
`run_35a886b0-509f-4467-8535-9267d74ba42f`. The frozen input digest is
`ed5c1cd6cd937fa2c2f35aed51d2476f4adef72559d5a360d28f09a7b3e8515d`.
The completed raw advice digest is
`637f56db9c8b4153ce7e7a821f74dc18ef8849878bc5d6da98581774b2e52934`.
The requested Nightfall response verification completed with exit 0. These records
belong to the explicitly invoked consultation, not a code-validation workflow.
The submitted input and raw advice remain intact. Living-document drift consists
of lifecycle metadata, the supersession notice and this evaluation.

All accepted actions below mean documenting the proposed design, not implementing
its Go types. The no-implementation instruction overrides the skill's default
production delivery. No build, test, linter, documentation checker or generated
provider matrix was run. The source excerpts in this document and direct source
inspection supply the evaluation basis; native/compiler behavior is not claimed
from this task. The existing pagination implementation is unchanged.

During documentation delivery, main advanced to `2c1000e4147830e8074f4d0fbf029e739b46f38d`
with the separate native kind-constraint proposal. That documentation change was
merged into this branch and cross-linked. It changed no Go source from the consulted
baseline. The plugin proposal preserves those enforcement obligations without
claiming that its initial profile settles witness/evidence parameter expansion.

| # | Recommendation | Verdict | Action | Owner | Link | Validation |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | Keep one existing ASTTransformer edge and explicit composition. | HOLDS | Retain FromID/ToID/Transform and the existing typed bridge; propose no engine method or provider cell. | codex-api-retirement-audit | [Composition](../docs/type-generation-plugins.md#composition-through-the-existing-pipeline) | Existing transform interface/bridge inspected; docs only. |
| 2 | Exact bindings, mandatory replacement guard and separate structural policy. | NEEDS ADAPTATION | Replace the candidate's virtual structural-default overrides with exact additions plus checked replacements. Use BindingID for configured instances and PluginID for implementations. Changed configuration gets a new BindingID. | codex-api-retirement-audit | [Selection](../docs/type-generation-plugins.md#exact-selection-and-explicit-overrides) | Existing registry rejects duplicate edges and semantic source identities distinguish declarations; new selection remains proposed. |
| 3 | Immutable selection, isolated invocations and terminal errors. | HOLDS | Document no order precedence, retries, global state or cross-invocation cache; no partial output. | codex-api-retirement-audit | [Configuration](../docs/type-generation-plugins.md#configuration) | Existing typed bridge snapshots source and returns errors atomically. |
| 4 | Exact alias identity does not prove native equality/assignability. | HOLDS | Keep identity selection separate from source-owned eligibility and profile promises. | codex-api-retirement-audit | [Selection](../docs/type-generation-plugins.md#exact-selection-and-explicit-overrides) | Horos identity facts and alias/application handling inspected; no inferred equality promise. |
| 5 | Complete Source/SourceSession with roots, inspect, declaration and projection. | HOLDS | Specify environment-bearing Ref views over existing Semantics, preserving actual binder identity and designated type operands. | codex-api-retirement-audit | [Source](../docs/type-generation-plugins.md#source-references-binding-and-projection) | Existing Subject, ArgumentBinding, identity and descriptor-operand facts inspected. |
| 6 | Provide real factories for opaque source handles. | HOLDS | Specify host-created RefFactory for independent adapters; retain session ownership and private source tokens. | codex-api-retirement-audit | [Source](../docs/type-generation-plugins.md#source-references-binding-and-projection) | Existing semantic/AST IDs are local addresses, so public opaque fields alone would not expose an extension path. |
| 7 | Remove Mode and separate Use, Publish and Define. | HOLDS | Document optional Publisher/Definer alongside Plugin.Use; requested roots require explicit local publication. | codex-api-retirement-audit | [Operations](../docs/type-generation-plugins.md#use-publication-and-definition) | Current prepare reserves declarations before field/type emission; generic declarations are distinct from proper uses. |
| 8 | Own source definitions or helpers under the selected configured binding. | HOLDS | Specify copied DefinitionSpec, open source identity/helper keys, no definition closures, and explicit binder links. Helpers do not capture source binders. | codex-api-retirement-audit | [Ownership](../docs/type-generation-plugins.md#owned-declarations-and-helpers) | Go's current presence helper demonstrates declaration creation without a source declaration. |
| 9 | Immutable native fragments with symbolic names and references. | NEEDS ADAPTATION | Adopt templates and per-occurrence cloning; clarify that generic wrappers mark intent while target grammar admits native categories. | codex-api-retirement-audit | [Fragments](../docs/type-generation-plugins.md#native-fragments-and-symbols) | Existing AST separates containment, textual attributes and references; native checkers own grammar. |
| 10 | Target lifecycle, storage capabilities and deterministic symbol allocation. | NEEDS ADAPTATION | Adopt target/session/storage interfaces; add opaque ResourceKey.OrderKey for independent deterministic allocation without decoding source identity. | codex-api-retirement-audit | [Target](../docs/type-generation-plugins.md#target-session-and-construction-capabilities) | Existing TreeID/NodeID cannot serve as reusable cross-session handles; target aliases/scopes are native concerns. |
| 11 | External imports are syntax, not fabricated resolved AST links. | HOLDS | State supplied artifact/export requirements and emit native imports/qualification when no declaration tree is supplied. Keep Files/lock/environment separate. | codex-api-retirement-audit | [Target](../docs/type-generation-plugins.md#target-session-and-construction-capabilities) | AST Reference requires an actual target; ast.Set is not a Files Source or lockfile. |
| 12 | Deferred reservation queue rather than mandatory Plan. | HOLDS | Document immediate reservation, one deferred Define per resource, queue closure before assembly and atomic failure. | codex-api-retirement-audit | [Composition](../docs/type-generation-plugins.md#composition-through-the-existing-pipeline) | Current prepare/body/Finish split supplies the concrete need; a whole-program Plan method is not established by this task. |
| 13 | Generic coherence applies to plugin behavior, not just selector shapes. | HOLDS | Require uniform families under supported substitutions; prohibit hidden argument-sensitive storage/helper/ownership specialization. Treat this as a trusted-plugin obligation. | codex-api-retirement-audit | [Generics](../docs/type-generation-plugins.md#generics-and-representation-promises) | Current open Page generation and binder references expose the List<Byte> counterexample. |
| 14 | Keep the first profile narrow. | HOLDS | Propose named roots, proper parameters, direct applications and explicit helpers. Defer application selectors, higher-order/computed constructors, monomorphization and anonymous publication until their semantics are specified. | codex-api-retirement-audit | [Generics](../docs/type-generation-plugins.md#generics-and-representation-promises) | Current pagination already rejects unsupported higher-order and recursive cases. No new capability is claimed. |
| 15 | Builtin wrapper/helper internals use native primitives directly. | HOLDS | Document Text and Presence examples; canonical child types still compose through Lower. | codex-api-retirement-audit | [Builtin helper](../docs/type-generation-plugins.md#a-builtin-override-that-generates-a-helper) | Current Go presence implementation uses native bool and parameter references; text reentry would select the same override. |
| 16 | Select before representation-specific recursion rejection. | HOLDS | Keep required Horos admission; examine demanded representation dependencies after selection. Native recursive legality remains a separate obligation. | codex-api-retirement-audit | [Composition](../docs/type-generation-plugins.md#composition-through-the-existing-pipeline) | Current visitDeclaration rejects source cycles; a selected external mapping need not demand that body. |
| 17 | Exact selection does not prove child/external compatibility. | HOLDS | Require stated target-position and artifact assumptions; no new generic capability calculus. Explain that overrides do not rewrite prebuilt dependencies. | codex-api-retirement-audit | [Selection](../docs/type-generation-plugins.md#exact-selection-and-explicit-overrides) | Current profiles already distinguish representation from enforcement; native semantic providers are unavailable. |
| 18 | Evidence requirements cover plugin eligibility explicitly. | HOLDS | Source.Contract states supported mapping requirements; additional facts require explicit composition and supplied receipts. Missing facts stay unknown. | codex-api-retirement-audit | [Configuration](../docs/type-generation-plugins.md#configuration) | Existing transform.Contract and captured evidence are explicit; no hidden analyzer call is present. |
| 19 | Bound expanded work and isolate invocation state. | HOLDS | Specify materialized-node/byte/import/definition accounting and cooperative cancellation without claiming a Go sandbox. | codex-api-retirement-audit | [Target](../docs/type-generation-plugins.md#target-session-and-construction-capabilities) | Existing arena limits/cancellation and bounded AST sets inspected. No limits changed. |
| 20 | Keep mappings and structural policy in type-owned integrations. | NEEDS ADAPTATION | Adopt the dependency direction; add the existing explicit Syntax dimension to proposed paths, and retain Bittype's root module. Default composition remains outside it. | codex-api-retirement-audit | [Layout](../docs/type-generation-plugins.md#proposed-file-layout-and-dependencies) | Existing LAYOUT and package ownership require stage/dimension symmetry. |
| 21 | Keep the SPI experimental and avoid a richer generic semantic model. | HOLDS | Use existing semantic carriers plus operational views and source-owned adapters. No public API stability or implementation is claimed. | codex-api-retirement-audit | [Decision](../docs/type-generation-plugins.md#expert-recommendation-and-chosen-proposal) | Only Horos supplies production semantic analysis today. |
| 22 | Use prior art selectively. | HOLDS | Retain the type/value separation and external-model precedent; reject MLIR retry precedence and Rust's global orphan policy as unsuitable here. | codex-api-retirement-audit | [Prior art](#prior-art-and-limits-of-analogy) | Primary references and earlier research 0003/0005 inspected; no copied selection semantics. |

### Explicit implementation deferrals

The trigger for implementing the proposed Go contracts, coordinator, adapters,
plugins and migration is a separate owner request to implement this design.
All accepted recommendations in this task are delivered as documentation. No
compatibility layer, runtime adapter, new language capability or production
mapping is introduced by recording the proposal.

Application specialization and higher-order generation need a separately agreed
coherence/profile design even after basic plugin implementation is authorized.
Fragment-local type/lifetime binders and anonymous root publication need concrete
supported target policies. A future mandatory whole-program Plan or per-mapping
evidence API should be justified by a case the documented queue/source contract
cannot handle. Runtime bindings remain separate work with separate authorization.
