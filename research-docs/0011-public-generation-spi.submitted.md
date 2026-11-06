# Research: Public generation operations and scoped association products

**ID:** 0011
**Date:** 6 November 2026
**Status:** reviewed
**Owner:** codex-hydrated-type-design
**Reviewed:** https://github.com/Bitspark/bittype/issues/102#issuecomment-6678276376
**Issue:** https://github.com/Bitspark/bittype/issues/102
**Code baseline:** Bittype e13a6f456d0bdb27e7e7facfd514dbd0cb0dbf67; Horos Go 0.10.2; public Logos contract fc6132c238eae843720acf096e884550454351c8

## Question

Complete the smallest usable public Go service-provider interface (SPI) for
type-owned generation: independent Use, Publish and Define capabilities, sharing
a paired Introduce/Apply convention for generic declarations. How should public
issuance, immutable access, source assumptions, native formal layouts and result
associations fit together without turning the agnostic host into a semantic
evaluator or allowing it to bless arbitrary payload/support pairs?

The owner authorized implementation of this first SPI slice, tests and docs.
The next milestone replaces all four fixed production generators. This slice
must be executable through independent-package fixtures, not merely a diagram,
but does not implement the full mapping selector/composer, native constraint
enforcement, runtime bindings or adapters.
This slice uses preselected collaborators and public lifecycle operations driven
by independent fixtures. Automatic selection and definition-queue orchestration
remain the later milestone.

## Context

### System overview and vocabulary

Bittype is a Go toolchain for canonical type descriptions and TypeScript, Go,
Rust and C++ type declarations. Horos is a separate library owning canonical
module/import/generic/kind semantics. Its descriptions use immutable Ontos data.
Native syntax has a language-marker parameter, for example AST[Go]. Bittype's
agnostic semantic graph stores provider-owned facts; it does not interpret them.
A kind classifies type descriptions. A proper type describes values; a type
constructor accepts type arguments. An open declaration is the unspecialized
declaration with its real binders. A profile specifies an integration's supported
domain and representation/enforcement promise.

A **snapshot** is one immutable, syntax-admitted input capture. Equal syntax
captured twice has different snapshot identities. A **receipt** is an immutable
result of a particular resolution/analysis operation on that snapshot. A
**source use** is a type occurrence interpreted under its actual source arguments
and assumptions. A **binder** is the actual source declaration parameter, not a
name, index or target parameter position. A **native payload** is the target type
expression selected to represent a source use. **Support** is integration-owned
information connecting that payload with its source use and assumptions; it is
not necessarily a proof term or a generated native witness type.

A **generation product** combines that scoped representation with its support.
A **convention** fixes how one open canonical declaration introduces native
formals and how its uses apply those formals. A **signature/layout** is the
authoritative target parameter/constraint arrangement, which may have more or
fewer native parameters than source binders. A **reservation** supplies a stable
native declaration symbol before its body is defined. A **fragment** is reusable
native syntax with symbolic references, before assembly into a target AST.

### Current architecture and established decisions

The current pipeline is explicit:

```text
Horos source -> parse -> syntax capture -> resolve -> analyze
  -> configured ASTTransformer -> target ASTSet
  -> formatter -> CST -> printer -> native text
```

An ASTSet is a collection of abstract syntax trees; CST is a concrete syntax tree
including tokens/trivia. Transformations return syntax, never hidden target
analysis or a proof result. Source analysis is not rerun when a mapping needs
additional facts. A missing capability produces an explicit unsupported result.

The four current generators share a private fixed Dialect with text, boolean,
list, record and alias operations. It supports pagination with ordinary proper
type parameters. It is not the desired public SPI and will be replaced in the
next milestone. Native enforcement is explicitly Unsupported; representation
generation currently permits erasing kind constraints for source-checked uses.

PR101 implemented an agnostic artifact-association contract and a concrete
Horos/Go verification composition. Horos supplies authoritative membership
queries; a scoped Go projection translates AST directly to go/ast and uses
go/types. An owning mapping constructs the expected Logos statement separately
from a candidate. Logos checks the statement/candidate under explicit policy.
TextBox's admission counterexample and a conditional Page representation theorem
and nested reuse now pass. Logos owns proof syntax/checking, not Bittype's core.

This verification consumes **completed source and target captures**. Generation
products precede final target assembly. They must not retag a checked result onto
future output. The eventual output is captured and explicitly projected/checked
after assembly; correspondence to a reused theorem remains a separate obligation.

### Prior expert advice reused

Consultation0007 (run_35a886b0-509f-4467-8535-9267d74ba42f) established exact
mapping selection, guarded overrides, separate Use/Publish/Define, immediate
reservation/deferred definition and target-owned symbolic fragments.
Consultation0008 (run_aa6f8ad3-4bbc-4b25-a894-fda48f62c8e5) corrected its
binder-to-native-index scheme to **real binder -> scoped product**, with paired
Introduce/Apply and complete application-result support. Its public construction
API was explicitly deferred. Consultation0009 separated native admission from
payload representation. Consultation0010 separated language judgments, proof
checking, artifact association and model/compiler adequacy. Their accepted
documentation is delivered; the generation SPI is not implemented.

We are not reopening these decisions. The new question is the concrete public
issuance/lifecycle seam now that the first final-output verification exists.

### Responsibility split

| Concern | Owner |
| --- | --- |
| Source applicability, real binders, substitutions and original assumptions | Source integration consuming supplied snapshot-bound semantic receipts and explicit immutable provider-owned query results |
| Native representation choice and associated laws | Selected builtin or library mapping/integration |
| Formal introduction, complete native application and result association | Configured source/target convention |
| Native syntax, symbols, imports, formal constraints and assembly | Target implementation |
| Invocation/scope identity, reservations, operation lifecycle, budgets and whole-output failure | Agnostic generation host |
| Claims over final captured target artifacts | Explicit verification composition and Logos checker |

All builtin mappings use the same supplied mechanism as library mappings. No
central Type/Data/Function/Stream/Cell cases, automatic imports/registrations or
universal type-system/proof enum are permitted. A library-owned Cell mapping
changes that library and consumer configuration only. Go plugins are trusted
in-process code, not a sandbox; nonetheless accidental mixing should be prevented
by construction rather than detected only in convention documentation.

### Five operations and their interactions

These names differ from identically named operations in the code appendix:
semantic.Request.Publish seals a semantic receipt, and Logos typed.Define creates
a logical definition. Generation Publish exposes a native declaration; generation
Define fills a reserved native declaration. Neither issues a semantic verdict.

- **Use** obtains the selected representation product for an admitted use.
- **Publish** exposes a requested root declaration. A missing capability is an
  error for that request; the root must not silently disappear.
- **Define** fills an owned reservation in its own formal environment.
- **Introduce** fixes native formals/constraints and binds every real source
  binder exactly once to an appropriate product. Extra native formals create no
  additional source binder. This happens for the canonical open declaration.
- **Apply** consumes effective source actuals and products in the caller's
  environment and produces the complete native use and result association.

Optional capabilities are separate interfaces. Inline primitive uses need no
Publish or Define. External generic types need Apply without local ownership.
An inline use can still require an owned helper. These are not exclusive modes.
Nongeneric declarations need an explicit route without fake generic methods.

There are two different environments: source binder -> effective source actual,
and source binder -> native representation product. A target formal environment
also owns native symbols and constraint syntax. Immutable handles can still
capture the wrong environment, so copying alone does not establish scope safety.

### Concrete cases that the SPI must express

```text
type Page<T> = { items: list<T>; next?: text; };
type Person = { name: text; active: bool; };
type Outer<U: K> = Page<U>;
type Nested<U: K> = Page<Page<U>>;
type Ignore<T: K> = text;
```

K is an ordinary resolved Horos kind; the integration must supply the actual
membership/constructor guarantees required by its profile. The generic host
does not interpret K. For a suitable selected K, the Page theorem says:

```text
forall T, N.
  Member(T,K) and Rep(T,N)
  -> Member(Page<T>,K) and Rep(Page<T>, NativePage<N>)
```

T ranges over canonical types under retained source assumptions; N ranges over
native representations in the selected mapping's supported domain. Member(T,K)
is Horos kind membership. Rep(T,N) is the selected static type/schema relation,
not runtime conversion or a claim about every constructible native value. The
independently established constructor guarantee Page: K -> K supplies the source
membership part.

Outer passes Outer.U's product, never Page.T's product. Nested's outer Apply
consumes support for the inner **Page<U> result**, not support for U beside a
Page-shaped payload. Ignore still requires whole-use source admission although
its payload never traverses T. Calling Page<Person> first must not set Page's
open signature or capture Person. A helper queued while visiting Outer needs its
own explicit formals, not an implicit link to Outer's frame. Equal-looking
captures, signatures or native type text do not authorize reuse across owners.

### Candidate public surface to critique

We are considering contracts/generation/go for public carriers/operations and a
small generation/syntax/go invocation driver for the demonstrated lifecycle.
The eventual selector/queue composer, language backends and independent mapping
modules remain separate. This first step must not expose the fixed Dialect.

One possible ownership arrangement lets the integration own an opaque immutable
product token P. The host never creates P by pairing arbitrary fields; a configured
association protocol exposes/admits its payload and source/scope metadata.
Alternative: core-owned opaque products with request-bound publication that calls
the selected semantic integration to capture and validate support before sealing.
We need a concrete choice. A generic parameter such as C any or a struct copy
does not ensure immutability; mutable slices/pointers/closures remain mutable.
Reflection cloning or a universal string-keyed evidence bag is unacceptable.

The following is conceptual Go notation, not a promised signature. F and T are
source/target language markers; P is the integration's product token. Request,
layout, publication and definition carriers would be immutable typed handles
issued through the public operation lifecycle, not freely writable records:

```go
type User[F, T, P any] interface {
    Use(context.Context, UseRequest[F,T,P]) (P, error)
}
type Publisher[F, T, P any] interface {
    Publish(context.Context, PublishRequest[F,T,P]) (Publication[T], error)
}
type Definer[F, T, P any] interface {
    Define(context.Context, DefineRequest[F,T,P]) (Definition[T], error)
}
type Convention[F, T, P any] interface {
    Introduce(context.Context, IntroduceRequest[F,T,P]) (Introduction[F,T,P], error)
    Apply(context.Context, ApplyRequest[F,T,P]) (P, error)
}
```

The host should mediate Apply, select the exact saved convention/signature,
admit the whole source use before child traversal and validate result association.
The integration owns support vocabulary, required semantic judgments and safe
issuance/access. Target helpers own native AST construction and do not decide
source applicability. No arbitrary source-to-product retagging operation exists.

For a locally owned declaration, reservation may first allocate an owner/scope,
then run Introduce once on the
open declaration so the target can allocate the complete formal layout, then
freeze it before Define/Apply. Reserving source-arity-many target slots would be
wrong for witness-dependent or extra-parameter conventions. Define should not be
able to replace/drop that layout: either the target assembles a saved header with
the returned body, or explicitly checks full declarations against the saved
layout. The agnostic host cannot inspect native constraint syntax itself.
External/inline conventions need an equivalent immutable signature route without
mandatory local ownership.

For reusable syntax, the existing AST disallows shared containment nodes and
cycles. Each materialized use must have fresh containment nodes. Shared symbols
remain references; an external artifact/export must not be given a fabricated
resolved AST NodeRef. We want a small typed fragment boundary, not a second AST
or grammar and not raw source strings.

### Limits, failures, validation and scope

Invocations are bounded by explicit work/operation/depth/output limits. No
cross-invocation cache, automatic monomorphization or hidden fallback is planned.
Error/cancellation invalidates the invocation's output, rather than returning
partially usable generated declarations. A provider may be unusable for one
requested capability without needing placeholder implementations for others.

Tests must exercise independent packages importing only public interfaces:
Page/Outer/Nested, ignored constrained actuals, renamed/same-spelling binders,
foreign captures/configurations/scopes, layout disagreement, source-support vs
result-support mismatch, stale publication, caller-frame escape, mutation,
missing optional capability, cancellation and limits. A small Go fixture should
assemble actual type ASTs and compile remotely; it must not claim native source
constraint enforcement. Existing production transformers remain unchanged until
their explicit replacement milestone. Native/generic grammar expansion is not
required merely to demonstrate these public contracts.

The owner requires remote compilation/testing and no runtime bindings/adapters.
The requested consultation retains its exact submitted input, advice and recovery
identifiers. As in consultation0010, the owner's standing rule excludes a
separate checksum or administrative/provenance validation workflow.

## Questions for the expert

The intended control flow is: admit supplied source evidence; invoke the selected
capability; introduce the open layout through the applicable local/external/inline
route; apply the retained convention in the caller's environment; complete/assemble
native definitions; separately capture and verify final output when requested.

1. What concrete public issuance/access design best keeps source admission,
   immutable representation support and host scope/lifecycle separate? Please
   choose between or replace the candidate product arrangements, including who
   may mint each handle and which obligations cannot be enforced by Go typing.
2. How should generic Use/Publish/Define share one Introduce/Apply convention,
   callee acquisition and authoritative formal layout? How should inline,
   external and owned cases compose, while nongeneric mappings can omit generic
   capabilities?
3. How should source substitutions, native product environments and application
   result support be carried so Outer, Nested and Ignore compose correctly?
   Please show a complete minimal Go boundary and data/control flow, not only
   methods whose critical public construction or safe-access types are omitted.
4. What should this first implemented SPI guarantee before target assembly, and
   how should later explicit final-output verification use it without treating
   symbolic products as already-checked target artifacts? Which tests and scope
   limits are essential before exposing the SPI?
5. What important concern or simpler decomposition have we missed? Please
   challenge any unnecessary coupling or abstraction, and distinguish work
   required for this first slice from full generator/plugin migration.

## Relevant current code

The extracts below are actual code at the stated baseline. They supply the
supporting types and call paths; the expert need not access the repository.

### contracts/transform/go/transformer.go:13

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

### contracts/transform/go/adapter.go:1

```go
package transform

import (
	"context"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/internal/tree"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"slices"
)

// Request supplies an owned typed source and admitted immutable evidence.
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
func (b bound[F, T]) Transform(ctx context.Context, input Input) (representation.ASTSet, error) {
	if ctx == nil {
		return representation.ASTSet{}, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return representation.ASTSet{}, err
	}
	evidence, err := semantic.AdmitEvidence(input.Source, b.requirements, input.Evidence)
	if err != nil {
		return representation.ASTSet{}, err
	}
	envelope, err := input.Source.EntryTree()
	if err != nil {
		return representation.ASTSet{}, err
	}
	source, err := representation.Unpack[F](b.from, envelope)
	if err != nil {
		return representation.ASTSet{}, err
	}
	result, err := b.stage(ctx, Request[F]{source: source, entry: input.Source.Entry(), evidence: evidence, inputs: input.Source.Inputs()})
	if err != nil {
		return representation.ASTSet{}, err
	}
	if err := ctx.Err(); err != nil {
		return representation.ASTSet{}, err
	}
	return representation.PackSet(b.to, result)
}
```

### contracts/semantic/go/operation.go:18

```go
type snapshotState struct {
	syntax    representation.Language
	files     *files.Context
	entry     ast.NodeRef
	trees     []representation.AST
	catalogue map[ast.TreeID]ast.NodeID
}

// Snapshot identifies one admitted input collection. Identical addresses or
// bytes captured again are a different snapshot. Storage is private and immutable.
type Snapshot struct{ state *snapshotState }

func (s Snapshot) Valid() bool { return s.state != nil }
func (s Snapshot) LanguageID() language.ID {
	if !s.Valid() {
		return ""
	}
	return s.state.syntax.ID()
}
func (s Snapshot) Entry() ast.NodeRef {
	if !s.Valid() {
		return ast.NodeRef{}
	}
	return s.state.entry
}
```

### contracts/semantic/go/operation.go:381

```go
func (r Result) Snapshot() Snapshot {
	if !r.Valid() {
		return Snapshot{}
	}
	return Snapshot{state: r.state.owner.source}
}
```

### contracts/semantic/go/operation.go:396

```go
func (s Snapshot) Language() representation.Language {
	if !s.Valid() {
		return representation.Language{}
	}
	return s.state.syntax
}
func (s Snapshot) Same(other Snapshot) bool { return s.Valid() && s.state == other.state }
```

### contracts/semantic/go/request.go:1

```go
package semantic

import (
	"errors"
	"fmt"
	"slices"
	"unicode/utf8"

	"github.com/Bitspark/bittype/internal/tree"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// capture is also the invocation identity. It is private, nonempty storage;
// equal entry addresses in different captures do not establish compatibility.
type capture struct {
	source        *snapshotState
	operation     semantics.Key
	dependencies  []Dependency
	contract      Contract
	requested     []semantics.Key
	configuration []Setting
	limits        Limits
}

func (c *capture) contains(site ast.NodeRef) bool {
	return site.Tree != "" && site.Node > 0 && site.Node <= c.source.catalogue[site.Tree]
}

type requestState[L any] struct {
	capture *capture
	inputs  []ast.AST[L]
}

// Request owns a captured syntax/context snapshot. Its zero value is invalid.
// Accessors return owned copies. No accessor fetches dependencies or resolves names.
//
// bittype:dimension Files
type Request[L any] struct{ state *requestState[L] }

func (r Request[L]) Entry() ast.NodeRef {
	if r.state == nil {
		return ast.NodeRef{}
	}
	return r.state.capture.source.entry
}

func (r Request[L]) Inputs() []ast.AST[L] {
	if r.state == nil {
		return nil
	}
	result := make([]ast.AST[L], len(r.state.inputs))
	for i, input := range r.state.inputs {
		result[i] = tree.Clone(input)
	}
	return result
}

func (r Request[L]) Requested() []semantics.Key {
	if r.state == nil {
		return nil
	}
	return slices.Clone(r.state.capture.requested)
}

func (r Request[L]) Configuration() []Setting {
	if r.state == nil {
		return nil
	}
	return cloneSettings(r.state.capture.configuration)
}

func (r Request[L]) Contract() Contract {
	if r.state == nil {
		return Contract{}
	}
	return r.state.capture.contract.Clone()
}

func (r Request[L]) Limits() Limits {
	if r.state == nil {
		return Limits{}
	}
	return r.state.capture.limits
}

// Publish seals output while the provider owns its mutable construction storage.
// It does not attach this request's identity to a previously published result.
// Re-publication is allowed, but returning an earlier invocation's Result is not.
func (r Request[L]) Publish(publication Publication) (Result, error) {
	if r.state == nil {
		return Result{}, ErrInvalidRequest
	}
	return publish[L](r.state.capture, publication)
}

type resultState struct {
	owner        *capture
	graph        representation.Semantics
	hasSemantics bool
	reports      []Report
	diagnostics  []Diagnostic
}

// Result is an immutable request-bound publication. A valid result can contain
// diagnostics and reports without a model. Valid does not mean semantic validity.
//
// bittype:dimension Files
type Result struct{ state *resultState }

func (r Result) Valid() bool        { return r.state != nil }
func (r Result) HasSemantics() bool { return r.state != nil && r.state.hasSemantics }
```

### contracts/semantic/go/evidence.go:13

```go
// request to run a provider. Semantics additionally requires a published graph.
type EvidenceRequirement struct {
	Operation     semantics.Key
	Profile       string
	Configuration []Setting
	Reports       []Requirement
	Semantics     bool
}

func (r EvidenceRequirement) Clone() EvidenceRequirement {
	r.Configuration = cloneSettings(r.Configuration)
	r.Reports = slices.Clone(r.Reports)
	return r
}

// ValidateEvidenceRequirements applies one shared metadata budget to the set.
func ValidateEvidenceRequirements(requirements []EvidenceRequirement) error {
	limits := DefaultLimits()
	if len(requirements) > limits.MaxInputs {
		return ErrLimit
	}
	budget := newMetadataBudget(limits)
	for _, r := range requirements {
		if err := validateDependencySpecification(r.Operation, r.Profile, r.Configuration, r.Reports, &budget); err != nil {
			return err
		}
	}
	return nil
}

// AdmitEvidence shares semantic dependency admission, then requires complete
// reports and any requested model. It never invokes a resolver or analyzer.
func AdmitEvidence(source Snapshot, requirements []EvidenceRequirement, results []Result) ([]Result, error) {
	if !source.Valid() || len(requirements) != len(results) {
		return nil, ErrInvalidRequest
	}
	if err := ValidateEvidenceRequirements(requirements); err != nil {
		return nil, err
	}
	for i, r := range requirements {
		dependency, err := admitDependency(source.state, Dependency{Result: results[i], Operation: r.Operation, Profile: r.Profile, Configuration: r.Configuration, Reports: r.Reports})
		if err != nil {
			return nil, err
		}
		if !dependency.reportsEstablished() || r.Semantics && !results[i].HasSemantics() {
			return nil, ErrEvidenceUnavailable
		}
	}
	return slices.Clone(results), nil
}
```

### semantics/go/model.go:17

```go
var ErrLimit = errors.New("semantic model resource limit")

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

type Expression struct {
	At   Occurrence
	Term TermID
}

type Dependency struct {
	From     EntityID
	To       EntityID
	Relation Key
	At       *Occurrence
}

// Subject selects exactly one model-local entity or term. A subject needs no
// fabricated syntax site; current occurrence associations remain separate.
type Subject struct {
	Entity EntityID
	Term   TermID
}

// Fact associates a subject with a provider-owned fact term. A relation may have
// several values; its meaning and cardinality belong to the provider. Fact links
// are associations, not term containment or evidence of semantic equality.
type Fact struct {
	Subject  Subject
	Relation Key
	Term     TermID
}

// Data is construction storage. New bounds, validates and copies it before
// publication. The producer must not mutate it concurrently with New.
type Data struct {
	Entities     []Entity
	Terms        []Term
	Bindings     []Binding
	Expressions  []Expression
	Dependencies []Dependency
	Facts        []Fact
}

```

### trees/ast/go/ast.go:1

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

### trees/ast/go/node.go:1

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

### trees/representation/go/ast.go:16

```go
var (
	ErrInvalidLanguage  = errors.New("invalid language descriptor")
	ErrLanguageMismatch = errors.New("language or profile mismatch")
	ErrTypeMismatch     = errors.New("AST marker type mismatch")
	ErrInvalidAST       = errors.New("invalid AST envelope")
)

// AST is an owned, structurally checked, opaque snapshot. Copying the envelope is safe;
// Unpack returns fresh mutable storage. Its zero value is invalid.
type AST struct {
	id      language.ID
	profile string
	payload payload
}

type payload interface {
	marker() reflect.Type
	clone() payload
	inspect() Info
	referencesValid(map[ast.TreeID]ast.NodeID) bool
}
type typed[L any] struct{ value ast.AST[L] }

func (p typed[L]) marker() reflect.Type { return reflect.TypeFor[L]() }
func (p typed[L]) clone() payload       { return typed[L]{tree.Clone(p.value)} }

func (a AST) LanguageID() language.ID { return a.id }
func (a AST) Profile() string         { return a.profile }

// Snapshot owns a deep copy, including attributes, references and origin slices.
func (a AST) Snapshot() (AST, error) {
	if a.payload == nil || a.id == "" || a.profile == "" {
		return AST{}, ErrInvalidAST
	}
	return AST{id: a.id, profile: a.profile, payload: a.payload.clone()}, nil
}

// Pack accepts either a typed AST or its existing CST view. It validates and
// snapshots storage; language syntax admission is a separate check.Provider.
func Pack[L any](definition Language, view ast.View[L]) (AST, error) {
	if err := definition.acceptMarker(reflect.TypeFor[L]()); err != nil {
		return AST{}, err
	}
	if IsNil(view) {
		return AST{}, ErrInvalidAST
	}
	input := view.ASTView()
	if err := check.Structure(input); err != nil {
		return AST{}, fmt.Errorf("%w: %w", ErrInvalidAST, err)
	}
	result := AST{id: definition.id, profile: definition.profile, payload: typed[L]{tree.Clone(input)}}
	return result, nil
}

// Unpack checks the identity/profile/actual Go type and returns a deep copy.
// It does not imply language syntax admission or semantic validity.
func Unpack[L any](definition Language, input AST) (ast.AST[L], error) {
	if err := definition.acceptMarker(reflect.TypeFor[L]()); err != nil {
		return ast.AST[L]{}, err
	}
	if err := definition.match(input); err != nil {
		return ast.AST[L]{}, err
	}
	return tree.Clone(input.payload.(typed[L]).value), nil
}
```

### internal/horosgen/generate.go:103

```go
// Field is an already mapped record member, never a canonical type expression.
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

### transforms/horos-to-go/go/transform.go:43

```go
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
```

### verification/files/horos-to-go/go/joint.go:36

```go
type composition struct {
	policy typed.Policy
	limits Limits
}
type Composition struct{ state *composition }
type claim struct {
	input          association.Input
	requirements   []association.Requirement
	question       semantics.Key
	statement      typed.Statement
	origins        []Origin
	family         *family
	argument       *argument
	theorem        Result
	depth          int
	counterexample *counterexample
}

// Claim is sealed by this owning interpretation, never by a caller-supplied
// formula, verdict flag or arbitrary checker outcome.
type Claim struct{ state *claim }
type request struct {
	owner       *composition
	claim       *claim
	association association.Association
	statement   typed.Statement
	origins     []Origin
}
type Request struct{ state *request }
type result struct {
	request   Request
	candidate typed.Expr
	outcome   typed.Outcome
}
type Result struct{ state *result }
```

### verification/files/horos-to-go/go/joint.go:143

```go
func (c Composition) Prepare(ctx context.Context, input Claim) (Request, error) {
	if ctx == nil || c.state == nil || input.state == nil {
		return Request{}, ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return Request{}, err
	}
	b := input.state
	a, err := association.Capture(ctx, b.input, b.requirements, c.state.limits.Association)
	if err != nil {
		return Request{}, err
	}
	count, bytes := 0, 0
	if len(b.origins) > c.state.limits.Association.Sites {
		return Request{}, association.ErrLimit
	}
	for _, origin := range b.origins {
		if !validLocation(origin.Location, b.statement) {
			return Request{}, ErrInvalid
		}
		count += len(origin.Sites)
		if count > c.state.limits.Association.Sites {
			return Request{}, association.ErrLimit
		}
		for _, site := range origin.Sites {
			bytes += len(site.Participant.Owner) + len(site.Participant.Name) + len(site.Node.Tree)
			if bytes > c.state.limits.Association.MetadataBytes {
				return Request{}, association.ErrLimit
			}
		}
		if err := a.CheckSites(origin.Sites); err != nil {
			return Request{}, err
		}
	}
	if err := ctx.Err(); err != nil {
		return Request{}, err
	}
	return Request{&request{c.state, b, a, cloneStatement(b.statement), cloneOrigins(b.origins)}}, nil
}
func (c Composition) Check(ctx context.Context, r Request, candidate typed.Expr) (Result, error) {
	if ctx == nil || c.state == nil || !r.Valid() || r.state.owner != c.state {
		return Result{}, ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	outcome := typed.Check(r.state.statement, candidate, c.state.policy, c.state.limits.Checker)
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	return Result{&result{r, candidate, outcome}}, nil
}
```

### verification/files/horos-to-go/go/mapping.go:285

```go
func (f Family) BindArgument(source h.Declaration, target n.Type, fields []FieldBinding, membership h.Judgment) (Argument, error) {
	if f.state == nil || !source.Valid() || !source.Model().Same(f.state.record.source.Model()) || !target.Model().Same(f.state.record.target.Model()) || !source.Model().Same(membership.Model()) {
		return Argument{}, ErrInvalid
	}
	if membership.Verdict() != decl.Proven {
		return Argument{}, semantic.ErrEvidenceUnavailable
	}
	if !membership.Kind().Same(f.state.kind) || len(source.Definition().Parameters()) != 0 {
		return Argument{}, ErrCoverage
	}
	// Value() of a nullary declaration is its body, so compare the authoritative
	// descriptor with this declaration's exact descriptor after its explicit query.
	descriptor, err := membership.Value().Descriptor()
	if err != nil {
		return Argument{}, err
	}
	if !descriptor.Equal(source.Definition().Expression().Value()) {
		return Argument{}, ErrCoverage
	}
	r, err := bindRecord(source, target, fields, f.state.optional)
	if err != nil {
		return Argument{}, err
	}
	d, proof := argumentLogic(r)
	return Argument{&argument{f.state, r, membership, d, proof}}, nil
}

// Application retains the checked closed theorem as a definition and rechecks
// it with each application. It never turns an acceptance into an assumption.
func (f Family) Application(theorem Result, arg Argument, depth int) (Claim, error) {
	if f.state == nil || arg.state == nil || arg.state.family != f.state || !theorem.Valid() || theorem.Request().state.claim.family != f.state || theorem.Request().Question() != question("conditional-representation") || !theorem.Outcome().IsAccepted() {
		return Claim{}, ErrInvalid
	}
	if depth < 1 || depth > 4 {
		return Claim{}, association.ErrLimit
	}
	d := slices.Clone(f.state.declarations)
	d = append(d, typed.Define("family.theorem", theorem.Request().Statement().Expected, theorem.Candidate()))
	d = append(d, arg.state.declarations...)
	h, n := con("h.argument"), con("n.argument")
	for range depth {
		h = app(con("h.family"), h)
		n = app(con("n.family"), n)
	}
	input, req := associationInput(f.state.record.source.Model(), f.state.record.target.Model())
	sites := append(slices.Clone(f.state.record.sites), arg.state.record.sites...)
	return Claim{&claim{input: input, requirements: req, question: question("application-representation"), statement: typed.Statement{Foundation: typed.Foundation, Declarations: d, Expected: conjunction(member(h), rep(h, n))(0)}, origins: logicalOrigins(d, sites), family: f.state, argument: arg.state, theorem: theorem, depth: depth}}, nil
}
func (a Argument) Candidate(depth int) typed.Expr {
	if a.state == nil || depth < 1 || depth > 4 {
		return typed.Expr{}
	}
	h, n := con("h.argument"), con("n.argument")
	p := pair(conjunction(member(h), rep(h, n)), con("argument.member"), fixed(a.state.candidate))
	for range depth {
		p = app(con("family.theorem"), h, n, fst(p), snd(p))
		h = app(con("h.family"), h)
		n = app(con("n.family"), n)
	}
	return p(0)
}
```

### Supporting definitions and ownership notes

A Setting is a provider-qualified key plus byte value; a Requirement names an
obligation and expected verdict. Evidence admission compares the selected
operation/profile/configuration and complete reports against the exact snapshot.
Semantic Publication is mutable construction storage for a graph, reports and
diagnostics; Request.Publish bounds and copies it, preserving provider ownership
of the assertions. It does not prove their language-level truth.
Result.Snapshot returns the original capture. Recapturing equal Inputs creates a
new identity and must not substitute for it.

Arena.Add in the existing private generator creates a fresh node; Arena.Target
writes a resolved intra-output reference. That mutable arena is not an immutable
symbolic fragment API. ast.NodeRef is {Tree TreeID; Node NodeID}; ast.Origin stores
optional current/historical source locations without affecting native grammar.

```go
// trees/ast/go/set.go:5
type Set[L any] struct {
    Entry NodeRef
    Trees []AST[L]
}
// trees/ast/go/view.go:7
type View[L any] interface { ASTView() AST[L] }
// contracts/semantic/go/contract.go:151
type Publication struct {
    Facts *semantics.Data
    // Must be an admitted dependency; mutually exclusive with Facts.
    ReuseSemanticsFrom *Result
    Reports []Report
    Diagnostics []Diagnostic
}
```