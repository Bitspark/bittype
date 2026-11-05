# Research: Compositional generic generation without semantic kernel cases

**ID:** 0008
**Date:** 5 November 2026
**Status:** applied (documentation only; implementation deferred)
**Run-ID:** run_aa6f8ad3-4bbc-4b25-a894-fda48f62c8e5
**Reviewed:** https://github.com/Bitspark/bittype/issues/84#issuecomment-6664097367
**Graded:** https://github.com/Bitspark/bittype/issues/84#issuecomment-6664723752
**Owner:** codex-api-retirement-audit
**Issue:** https://github.com/Bitspark/bittype/issues/84
**Code baseline:** f2179283150b1e7344844fc97503ca26fe40795c
**Companion:** Research 0009, independently owned by codex-hydrated-type-design

## Question and requested decision

What is the smallest coherent composition protocol for translating source generic
declarations and applications when a source parameter may require native payload,
source-identity and membership information, without teaching a general generation
kernel the source kind algebra, target language, or library contract names?

We want a conceptual decision and a proposed Go contract, not implementation.
Please challenge whether parameter expansion or witnesses are the right model at
all. We are not asking you to bless the sketches below. An architecture that
avoids expansion while meeting the explicit guarantees is equally welcome.

This is one of two owner-requested consultations. This brief owns **composition**:
declaration headers, parameter binding, application arguments, plugin ownership,
explicit overrides and external artifacts. The independent companion owns
**sound native constraint encodings** in TypeScript, Go, Rust and C++, including
the precise supported domains, native escape hatches and cross-package evidence
construction. Each expert may challenge assumptions in both; neither response
will be treated as implementing the other expert's missing half.

The requested outcome is a documented design. No runtime adapters, generated
function bodies, value validators, new compiler, or production code is in scope.
The explicitly invoked BB Research workflow authorizes this consultation's
frozen input/recovery record and advice verification; it does not authorize
additional build/test/check workflows or production implementation.

## System, vocabulary and fixed boundaries

**Horos** is a canonical declaration and type-description library. Its semantic
implementation lives in a separate Go module. **Ontos** supplies values and
encoding. Horos descriptions classify data or capability-containing types;
ordinary kind definitions constrain those descriptions and type constructors.
Horos has generic declarations, applications, bounds and explicit named imports.
It does not generate TypeScript/Go/Rust/C++ source.

**Bittype**, implemented in Go, supplies language-neutral AST/CST carriers,
semantic graph storage, stage contracts, an explicitly configured engine,
source-specific providers and target-language tooling. AST means abstract syntax
tree; CST means concrete syntax tree. `AST[L]` uses a Go marker type `L` for its
language. Syntax is the text/CST/AST dimension; Files is the
files/organized-source/semantics dimension. Current native generators produce
type declarations only, as target AST collections, followed by AST-to-CST
formatters and a generic CST printer. There is no text/reparse generation step.

A **provider** implements an explicitly registered pipeline stage. A **snapshot**
captures the source syntax and context. A **receipt** is a published resolution
or analysis result tied to that exact snapshot. The current transformer bridge
admits explicitly supplied receipts; it never invokes analysis implicitly.
Semantic entities represent declarations/binders. Terms represent references,
applications and provider-owned expressions. Numeric IDs are graph-local
addresses, not durable declaration identities.

A **proper type** classifies values. A **constructor** accepts type arguments.
An **open declaration** has unbound parameters. A Horos **kind** is a predicate on
type descriptions or constructors, not a unique native representation tag.
Kinds can overlap. `T: h.Data` means T is constrained by the ordinary Data kind
from an explicitly imported prelude `h`. The prelude names Data, Type, Text,
List and so on are not reserved dispatch keys. A user can define a new kind
without adding a mapping implementation for that kind's name.

A **generation plugin** is ordinary compiled Go integration code selected through
explicit configuration. It is not a runtime registry, dynamically loaded shared
object, source parser extension, or native value adapter. Type/library owners
supply the optional representation mapping for each supported target language.
Builtin representations use that same supplied mechanism. **Mapping selector**
means an exact canonical builtin/declaration identity; this is unrelated to
structural observation selectors inside the Horos kind language.

A **representation profile** states a promise. The currently implemented profile
represents applications already admitted by Horos, and explicitly permits
erasure of generic kind constraints in native output. An **enforcing profile**
additionally promises that native applications within its stated domain cannot
claim unsupported source membership or misassociate it with another native
representation. It does not make arbitrary native runtime values conform to
Horos refinements. What exact native domains can meet that promise is the
companion consultation's question; this brief asks how such a domain composes.

A **source witness** is a proposed compile-time native type identifying a source
expression; **payload** is its native value representation. Whether witness identity
is a declaration, normalized application, or contextual expression remains an
open design question. Neither an exact mapping selector nor invocation-local Ref
settles it. **Evidence** here
means a proposed target encoding of an established source-membership fact or
conditional implication, tied to its representation. This is distinct from the
existing source semantic receipt. It is not an executable proof engine or an
arbitrary public `Proof<S,N,K>` constructor. No witness/evidence APIs are
implemented yet.

The following boundaries are owner decisions:

- Function, Stream, Cell and future contracts are ordinary library declarations.
  The general generator, engine, runtime, Wire and hydration must not recognize
  their names, structural fingerprints or behavior. Runtime behavior is out of
  this task's scope.
- Source-specific adapters may understand the Horos grammar and kind algebra.
  Target backends may understand native grammar, scopes and imports. Type-owned
  integrations choose the association between a source contract and native form.
- The engine registers one whole-language `ASTTransformer` edge per source/target
  pair. Type-owned mappings are collaborators within that transformer, not new engine
  providers or edges. Consumer configuration is explicit; no hidden default bundle.
- Exact identity controls plugin selection. Kind membership is an admission
  obligation, never a most-specific-kind dispatch search. Failure is not a signal
  to retry a weaker representation or a structural fallback.
- Preserve open generics and real binder identities. No silent monomorphization,
  first-application specialization, alias weakening or source-name matching.
- Use the existing generic syntax/semantic infrastructure. Justify any extension
  with a concrete case, rather than inventing another canonical type algebra.

Terminology used in the code appendices: `HeadKind`, `TermKind` and `ast.Kind`
are discriminator tags, distinct from Horos membership predicates. SPI means
extension interface; IR means intermediate representation. `decl.Digest` is a
module-content pin; a declaration identity adds its declared name, and a parameter
identity adds its ordered slot. Endpoint `representation.Language.Profile()`
identifies an explicit syntax/provider compatibility contract. A generation
representation/enforcement profile states the mapping promise; it is not that
endpoint profile. Existing API fields named `Evidence` hold source receipts,
while proposed native membership evidence is generated type-level support.
In sketches, S is source identity, N (or Go T) is native payload, K is a source
kind and P is a proposed native evidence parameter.

## Current source and prior research

The actual native implementation is still a deliberately limited pagination
generator in `internal/horosgen`, with a fixed private `Dialect` interface and
four concrete transformers. The proposed plugin SPI in Appendix B is not code.
Current generation supports proper generic parameters, direct applications,
named closed records, text, bool, lists and optional fields. In the private
generator, generic applications are nonempty and saturated; nongeneric uses are
references rather than zero-argument Application terms. It retains parameter
references and source origins. It rejects higher-order/computed constructors,
unused proper parameters and recursive native declarations. Ignore below is a
requirement for the proposed protocol beyond that implemented backend domain.
The enforcement profile returns `Unsupported`
with no output. Native ASTs/text tooling exist; native semantic providers do not.

Source graph membership facts and assumptions already exist. `Parameter.Sort`
only identifies proper/constructor sorts, not a rich kind contract.
`DirectParameterBindings` pairs a direct callee's real binders with supplied
arguments but does not establish applicability. A provider-specific fact encodes
its own semantics; a generic graph integrity check is not a semantic proof.

Research 0007 already answered the basic plugin question. Its completed response
was evaluated and incorporated into the current proposal, with production
implementation deferred. Relevant accepted advice:

1. Keep the existing transformer bridge; use exact bindings and explicit guarded
   overrides. `BindingID` identifies a configured mapping, not just its Go code.
2. Separate `Use` (proper type expression), optional `Publish` (requested open
   declaration) and optional `Define` (fill a reserved declaration). No exclusive
   inline/external/owned mode enum: those effects can coexist.
3. Reserve symbols immediately, queue immutable definition requests and define
   each owned resource once. Do not capture mutable definition closures or require
   a general public planning phase merely to reserve symbols.
4. Preserve full source substitution environments separately from target formal
   environments; provide factories for opaque references and symbols.
5. Reusable native fragments are immutable templates. Materializing each use
   allocates distinct containment nodes. Symbols are references, not fake AST IDs.
6. Require uniform generic families: mapping after source substitution must agree
   with substituting the mapped arguments, under the representation promise.
7. Select a declaration mapping before expanding its body; an external mapping
   may avoid source-body recursion. No eager flattening into primitives.

We are not asking for this review again. The newly exposed gap is that a single
`BinderLink{Source, TargetIndex}` cannot express native projections and additional
arguments for witnesses/evidence. More header slots do not define how the same
declaration is used. Exact mapping selection also says nothing by itself about
whether a native generic bound still matches an overridden child representation.

Reuse inspection found research 0002 on the semantic kernel, 0003 on optional
Wire/type-owned mappings, 0005 on general kinds, 0006 on the source-authored
prelude, and 0007 on plugins. 0001/0004 are superseded historical models. No
existing response settles this combined signature/application problem. Both new
consultations share the same code baseline and own distinct research records.

## Shared examples and counterexamples

The following is Horos declaration syntax, assuming an explicit import of its
ordinary prelude as `h`. Pins/import boilerplate are omitted, not inferred:

```text
type Page<T> = { items: list<T>; next?: text; };
type DataPage<D: h.Data> = Page<D>;
type Person = { name: text; active: bool; };
type People = DataPage<Person>;

kind Tree = nu {
  X = h.Text | h.Bool | h.List<X> | h.ClosedRecord<X>;
} in X;
type TreePage<T: Tree>: Tree = Page<T>;
type TreePeople = TreePage<Person>;
type Outer<T: Tree>: Tree = TreePage<T>;
```

`nu` introduces a guarded greatest-fixed-point kind. Tree admits text, bool,
lists whose element satisfies Tree and closed records whose fields do likewise.
It is a user-authored example, not a builtin. The source establishes TreePage's
result membership under the assumption that its argument satisfies Tree. In
`Outer<T>`, the actual T is Outer.T, not TreePage's formal binder. A wire type is
outside Tree. A descriptor with a union at its root is outside Tree even if its
alternatives are text/bool; Tree does not include a union-node pattern.

For comparison, `kind Leaf = h.Text | h.Bool` admits the text descriptor or the
bool descriptor. A native `T extends string | boolean` also admits the native
union; that is not enough to recover which source descriptor was meant. Two
source descriptions can map to the same payload type but differ in membership.

The candidate enforcing outputs in the docs are only public-shape sketches:

```typescript
type TreePage<S extends TreeMember> = Page<Native<S>>;
type People = TreePage<PersonType>;
```

```go
type TreePage[S, T any, P TreeEvidence[S, T]] = Page[T]
type People = TreePage[PersonType, Person, PersonInTree]
```

```rust
pub struct TreePage<S: TreeMember> { pub page: Page<S::Repr> }
pub type People = TreePage<PersonType>;
```

```cpp
template<class S> requires tree_member<S>
using TreePage = Page<typename S::representation>;
using People = TreePage<PersonType>;
```

`Native`, `TreeMember`, `TreeEvidence`, `PersonType`, `PersonInTree` and the
membership concept are proposed support definitions, not given magical
implementations. The companion expert must determine whether/how these can be
sound. Rust's wrapper is observably distinct from Page; a representation promise
requiring alias interchangeability might reject it. TypeScript witness unions
can lose correlation; returning `never` for a bad type is not formation rejection.
Go's illustrative private evidence constructor does not automatically compose
with a different package's private evidence type. C++ specializations raise
questions about who may assert associations. These are reasons to leave the
encoding abstract until its guarantees are stated.

For library composition, assume the independently authored, exactly identified
declaration Cell<T>, with a selected library representation:

```text
type EditablePeoplePage = Page<Cell<Person>>;
```

A structural Page mapping composes its list/field mappings, leaving Page generic.
The Cell plugin selects `cell.Cell<Person>` or an alternative representation.
Page must not need a Cell-specific rule. If Cell includes capabilities, this Page
application can be admitted but DataPage<Cell<Person>> cannot acquire Data
membership merely because Cell happens to map to an ordinary native struct.

Further cases a protocol should explain:

| Case | Required distinction |
| --- | --- |
| `Ignore<T: Tree> = text` | An ignored argument still has to satisfy its source bound; payload generation cannot drop admission obligations. |
| `Outer<T: Tree> = TreePage<T>` | Open assumptions and dependent native arguments flow through actual binder identities. |
| One source T represented by S, N and P | Additional native parameters are generated support, not invented source binders. |
| A field lowers T to `S::Repr` or `Native<S>` | Binder substitution may be a structured native expression, not one parameter-symbol lookup. |
| Text override from string to application Text | Source membership is unchanged; a native string constraint/evidence association might no longer fit. |
| List override from slice to Collection<T> | Open Page must use a coherent generic family; concrete-byte-specific storage cannot silently change only a later application. |
| Imported Page compiled with another mapping | Local overrides do not rewrite that external artifact. Its export contract must match or composition fails. |
| Same source expression in two effective configurations | Its source identity can agree while native association and support resources differ. |
| Recursive helper dependencies | Symbol reservation prevents emission-order dependence but does not establish legal native recursion or membership. |
| Inline library representation plus evidence helpers | One plugin can use inline syntax, external imports and owned declarations together. |

## Current proposal: flow and the missing connection

The proposed kernel captures one effective binding table and no semantic
catalogue. The source adapter admits the exact source evidence and exposes roots,
declarations, effective applications and actual type-child projections. It retains
lexical/substitution/assumption contexts in opaque source references. Exact
selectors distinguish canonical builtins from module-pinned named declarations.
The ordinary structural declaration policy is another explicitly supplied mapping.

For each requested named root, the kernel selects its mapping before expansion
and invokes Publish. The plugin calls Own with the open source declaration,
header and binder links; this reserves output symbols and queues Define.
Use constructs a complete proper native type, lowering only the children it needs
through the same configuration. Define is later called once with fresh formal
bindings. Helpers use mapping-local keys and do not capture source binders.
The target backend allocates names/imports/units and assembles structured ASTs.
Each native grammar remains in its own syntax/checker/formatter stages.

Ownership is `(invocation, configured mapping, open source declaration or helper
key)`. A source definition is never keyed by its first application. Definitions
may request other helpers. Errors discard the entire invocation; no partial AST
or weaker retry. A configured `BindingID` helps guard replacement and resources
but is not proof of semantic or external-artifact compatibility. The output is
an AST set, not filesystem layout or a native semantic-provider receipt.

With one native parameter per source parameter, the kernel can lower a parameter
through a simple formal environment. With the sketches above, the payload of T
might instead be a native projection, and using the named declaration might need
several linked arguments. The current header/binder protocol does not specify:

- how the header and every application share one representation of the generic
  signature, including bounds and assumptions;
- who derives the source argument's native representation and supporting
  associations, and how several plugins consume that result coherently;
- whether payload-only lowering is enough, or whether reusable products must
  carry opaque source/witness/constraint information;
- how a pending helper definition receives scoped support without captured
  closures or invented source binders;
- how an imported declaration communicates the exact argument transformation
  it expects, including changed output after an explicit override.

The existing proposal's generic-family law must therefore cover full applications,
not just payload AST text: mapping a substituted source declaration and applying
the translated generic should agree on admissibility, payload, witness/evidence
arguments and result association within the selected profile. The law is an
obligation on trusted Go mappings; the kernel cannot infer it from identical
names, equal-looking syntax or `BindingID`.

For `Outer<T: Tree>: Tree = TreePage<T>`, the current provider exports parameter
contracts, each declaration's result/constructor contract, symbolic body, and
established memberships with their original rigid assumptions. Outer.T and
TreePage.T remain different semantic entities. DirectParameterBindings supplies
the correspondence from TreePage.T to the reference to Outer.T; it does not prove
new membership. The proposed source adapter may carry that substitution through
actual child references while preserving the captured assumption environment.
Any requested judgment absent from the admitted model requires an explicit
extension of the source-provider contract and corresponding input receipts;
generation does not invent it. The current Horos analyzer reuses its explicit
resolver's semantic graph/reports, rather than exposing an on-demand proof oracle.

The protocol must transport capabilities and encoding domains whose guarantees
are established by their owning integrations. It must not promise a uniform
enforcing encoding before the companion study determines the target-specific
limits. Membership soundness, completeness within a declared domain, rejection
at type formation, and native alias interchangeability are separate promises.

## Alternatives to assess, not predetermined decisions

**Keep the core payload-only and let plugins manage everything.** A library
mapping could create its header, extra arguments and helpers privately. This
preserves a small kernel but might duplicate source-generic reasoning and make
independent mappings disagree on parameter conventions. Does a source/target
integration collaborator suffice without adding a common generation product?

**One reusable signature/application elaboration object.** A source/target-owned
object could derive native formals, a source-binder-to-payload projection, and
the argument transformation for every application. The kernel would carry opaque
handles and enforce lifecycle/ownership only. We have not defined its fields or
whether this is simply a restricted artifact contract rather than a new phase.
Avoid an all-purpose planning IR merely because a reservation queue exists.

**One witness-facing parameter with native representation projection.** This can
reduce public parameter count in some targets but is not equally expressible in
all four languages. A universal source/witness/payload/evidence tuple might bake
one encoding into the supposedly agnostic core. Native-specific shape must remain
possible without turning every support item into an untyped bag.

**Explicit elaboration before AST generation.** A source/target integration might
resolve representation obligations and generic conventions before native syntax
construction. This could separate semantic obligations from graph ownership,
but it must not become a second Horos evaluator, duplicate canonical kinds, or
re-run analysis without admitted receipts. Is the existing semantic model enough
to describe the selected associations, or should this be invocation-local
generation metadata? A new public provider is not presumed necessary.

**Limit the first enforcing profile.** A closed generated domain or unsupported
open constraint may be an honest starting point. We should still know which
extension seam will support ordinary library composition, and must not disguise
finite evidence enumeration as arbitrary recursive/open native admission.

## Questions for the expert

1. What is the right conceptual decomposition between source membership,
   representation choice, generic signature elaboration and native syntax
   construction? Identify the minimum shared object or law that makes open
   declaration publication and every application coherent. Contrast it with
   keeping all extra parameter logic private to plugins.
2. What concrete Go contract would express that decomposition while retaining
   the existing transformer/receipt boundary and type-owned mapping selection?
   Walk Page<T>, TreePage<T> and Outer<T> through formal binding, nested payload
   projection and concrete application. Specify ownership/lifetime of any new
   handles and how Use, Publish and Define share the same convention. If your
   recommendation uses expanded arguments or
   projections, illustrate those forms; otherwise explain how it satisfies the
   same source cases without them.
3. How should independently supplied mappings, explicit builtin/library overrides
   and imported compiled generic artifacts compose? Define the information that
   must travel with an export/application, its ownership, and the situations that
   must fail. Identify required semantic guarantees as explicit prerequisites;
   leave per-target soundness arguments to 0009. Do not invent a central kind
   catalogue, dynamic proof registry or hidden weakening.
4. What smallest implementation sequence and supported generic domain would make
   this design reviewable? Identify the concrete examples that separate a coherent
   minimal protocol from needless machinery, and the limitations that should stay
   explicit rather than being encoded into an overly broad public SPI.
5. What have we misunderstood or omitted? Challenge the premise, the proposed
   witness/evidence approach, and the boundary between the two consultations.
   Propose a simpler or different conceptual approach if it better satisfies the
   intended guarantees.

## Evidence appendices

The excerpts below are included so the answer does not depend on repository
access. Appendix A is implemented source at the stated baseline. Appendix B is
the proposed contract from the main documentation, not compiled code. Native
encoding sketches above deliberately have unresolved helper construction.

### Appendix A: implemented contracts and call sites

`contracts/transform/go/transformer.go:1-44` (implemented excerpt):

```go
// Package transform defines directed AST transformations and their typed bridge.
// Concrete mappings belong to separate transforms/<source>-to-<target> packages.
package transform

import (
	"context"
	"errors"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

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

`contracts/transform/go/adapter.go:1-92` (implemented excerpt):

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

`trees/ast/go/set.go:1-8` (implemented excerpt):

```go
package ast

// Set is construction storage for one language's complete multi-unit output.
// Entry is explicit, even for a one-tree set. It carries no file organization.
type Set[L any] struct {
	Entry NodeRef
	Trees []AST[L]
}
```

`semantics/go/model.go:1-132` (implemented excerpt):

```go
// Package semantics stores established, provider-owned semantic facts.
// It checks graph integrity, never a language's lookup or type rules.
package semantics

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Bitspark/bittype/trees/ast/go"
)

var ErrInvalid = errors.New("invalid semantic model")
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

type Limits struct {
	MaxRecords int
	MaxEdges   int
	MaxBytes   int
	MaxDepth   int
}
```

`Semantics[L]` is an immutable handle over a copied Data graph. Its implemented definition at `semantics/go/model.go:146` is:

```go
type Semantics[L any] struct{ state *Data }
```

`semantics/go/model.go:500-517` (implemented excerpt):

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
```

`semantics/go/model.go:519-553` (implemented excerpt):

```go
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
```

The exporter below is Horos-owned interpretation, not kernel logic. `e.parameters` maps Horos ParameterID to graph EntityID. `e.kind` serializes the resolved kind as provider-owned terms. The generated membership fact retains each original rigid assumption. Current providers export established facts, not an arbitrary on-demand proof oracle. `e.provider` serializes a provider expression and its explicitly located term operands; `e.fact` associates it with a subject/relation.

`internal/horosmodel/facts.go:50-72` (implemented excerpt):

```go
// Each membership is one established source guarantee with its original rigid
// assumptions. Query results never mutate this graph or imply a principal kind.
func (e *exporter[L]) memberships(subject semantics.Subject, memberships []decl.Membership, depth int) error {
	for _, membership := range memberships {
		kind, err := e.kind(membership.Kind(), ast.NodeRef{}, 0, false, depth+1)
		if err != nil { return err }
		operands := []semantics.Operand{{Locator: []uint32{0}, Term: kind}}
		for i, assumption := range membership.Assumptions() {
			reference, err := e.allocateTerm(depth+2)
			if err != nil { return err }
			if err := e.fillTerm(reference, semantics.Term{Kind: semantics.Reference, Target: e.parameters[assumption.Parameter()]}); err != nil { return err }
			contract, err := e.kind(assumption.Kind(), ast.NodeRef{}, 0, false, depth+2)
			if err != nil { return err }
			value, err := e.provider("membership/assumption", core.NewAtom(nil), []semantics.Operand{{Locator: []uint32{0}, Term: reference}, {Locator: []uint32{1}, Term: contract}}, depth+1)
			if err != nil { return err }
			operands = append(operands, semantics.Operand{Locator: []uint32{1, uint32(i)}, Term: value})
		}
		value, err := e.provider("membership/proven", core.NewAtom(nil), operands, depth)
		if err != nil { return err }
		if err := e.fact(subject, e.frontend.MembershipRelation(), value); err != nil { return err }
	}
	return nil
}
```

The existing exporter separately retains bounds and result contracts. In this block, `subject` is the declaration entity, `f` supplies Horos-owned relation keys, `pid` is the canonical parameter identity and `e.parameters` maps it to the local entity. `e.kind` exports the actual kind graph rather than a coarse sort. `internal/horosmodel/export.go:224-236`:

```go
				for index, parameter := range declaration.Parameters() {
					pid := decl.ParameterID{Declaration: declaration.ID(), Index: index}
					kind, err := e.kind(parameter.Kind, input.ParameterContractSite(pid), id, true, 1)
					if err != nil { return semantics.Data{}, err }
					if err := e.fact(semantics.Subject{Entity: e.parameters[pid]}, f.DeclaredKindRelation(), kind); err != nil { return semantics.Data{}, err }
				}
				result, err := e.kind(declaration.Result(), input.ResultSite(declaration.ID()), id, true, 1)
				if err != nil { return semantics.Data{}, err }
				if err := e.fact(subject, f.ResultKindRelation(), result); err != nil { return semantics.Data{}, err }
				contract, err := e.kind(declaration.Contract(), ast.NodeRef{}, 0, false, 1)
				if err != nil { return semantics.Data{}, err }
				if err := e.fact(subject, f.DeclaredKindRelation(), contract); err != nil { return semantics.Data{}, err }
				if err := e.memberships(subject, declaration.Memberships(), 1); err != nil { return semantics.Data{}, err }
```

The current private implementation maps parameter uses and arguments one-for-one. `Arena` is per-invocation target AST construction storage with a sticky error, bounded node allocation and fresh contained nodes. `Add` allocates, `Set` fills a reserved declaration, and `Target` adds a resolved syntax reference. The fixed Dialect is the implementation the plugin proposal will replace.

`internal/horosgen/generate.go:103-121` (implemented excerpt):

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

`internal/horosgen/generate.go:170-269` (implemented excerpt):

```go
// Generate bundles every exported type definition into one source unit. Kind
// definitions remain compile-time inputs. The representation profile erases
// kind contracts for uses admitted by Horos; it promises no native enforcement.
// Imported declarations retain module-qualified names; no import resolution runs.
func Generate[L any](ctx context.Context, request transform.Request[h.Language], id ast.TreeID, profile string, dialect Dialect[L]) (ast.AST[L], error) {
	fail := func(err error) (ast.AST[L], error) { return ast.AST[L]{}, err }
	source := request.Source()
	if id == "" || id == source.ID || request.Entry() != (ast.NodeRef{Tree: source.ID, Node: source.Root}) {
		return fail(semantic.ErrInvalidRequest)
	}
	switch profile {
	case RepresentationProfile:
	case EnforcementProfile:
		site := request.Entry()
		return fail(&text.Diagnostic{Code: text.Unsupported, Message: "native kind enforcement requires an explicit target representation", Node: &site})
	default:
		return fail(transform.ErrInvalidTransformer)
	}
	evidence, ok := request.Evidence(0)
	if !ok {
		return fail(semantic.ErrEvidenceUnavailable)
	}
	model, ok := evidence.Semantics()
	if !ok {
		return fail(semantic.ErrEvidenceUnavailable)
	}
	g := generator[L]{arena: Arena[L]{Tree: ast.AST[L]{ID: id}, Ctx: ctx}, dialect: dialect, source: source, facts: model.Data(), sites: map[semantics.TermID]ast.NodeRef{}, sorts: map[semantics.Subject]semantics.TermID{}, results: map[semantics.EntityID]semantics.TermID{}, identities: map[semantics.EntityID]semantics.TermID{}, names: map[semantics.EntityID]string{}, targets: map[semantics.EntityID]ast.NodeID{}, parameters: map[semantics.EntityID][]ast.NodeID{}, payloads: map[semantics.TermID]core.Tuple{}, state: map[semantics.EntityID]byte{}, used: map[semantics.EntityID]bool{}}
	g.sources = map[ast.TreeID]ast.AST[h.Language]{}
	for _, envelope := range request.Inputs() {
		input, err := representation.Unpack[h.Language](evidence.Snapshot().Language(), envelope)
		if err != nil {
			return fail(err)
		}
		if input.ID == id {
			return fail(semantic.ErrInvalidRequest)
		}
		g.sources[input.ID] = input
	}
	for _, e := range g.facts.Expressions {
		if e.At.Role == resolve.ExpressionRole() {
			g.sites[e.Term] = e.At.Site
		}
	}
	for _, fact := range g.facts.Facts {
		switch fact.Relation {
		case resolve.SortRelation():
			if g.sorts[fact.Subject] != 0 {
				return fail(semantic.ErrInvalidRequest)
			}
			g.sorts[fact.Subject] = fact.Term
		case resolve.ResultSortRelation():
			if fact.Subject.Entity == 0 || g.results[fact.Subject.Entity] != 0 {
				return fail(semantic.ErrInvalidRequest)
			}
			g.results[fact.Subject.Entity] = fact.Term
		case resolve.IdentityRelation():
			if fact.Subject.Entity == 0 || g.identities[fact.Subject.Entity] != 0 {
				return fail(semantic.ErrInvalidRequest)
			}
			g.identities[fact.Subject.Entity] = fact.Term
		}
	}
	if err := g.prepare(); err != nil {
		return fail(err)
	}
	var declarations []ast.NodeID
	for _, entityID := range g.order {
		entity := g.facts.Entities[entityID-1]
		target := g.targets[entityID]
		params := g.parameters[entityID]
		term := g.facts.Terms[entity.Body-1]
		if term.Operator == (semantics.Key{Owner: "horos", Name: "descriptor/record"}) {
			fields, err := g.fields(entity.Body)
			if err != nil {
				return fail(err)
			}
			dialect.Record(&g.arena, target, g.names[entityID], params, fields)
		} else {
			body, err := g.expression(entity.Body, 0)
			if err != nil {
				return fail(err)
			}
			dialect.Alias(&g.arena, target, g.names[entityID], params, body)
		}
		for _, p := range entity.Parameters {
			if !g.used[p.Entity] {
				return fail(g.unsupported(entity.Body, "unused generic parameters require a separate representation profile"))
			}
		}
		declarations = append(declarations, target)
	}
	g.arena.Tree.Root = dialect.Finish(&g.arena, declarations)
	if g.arena.Err != nil {
		return fail(g.arena.Err)
	}
	if err := ctx.Err(); err != nil {
		return fail(err)
	}
	return g.arena.Tree, nil
}
```

The formal-parameter loop inside `generator.prepare` is implemented at `internal/horosgen/generate.go:394-423`. Here `declarations` is the already-reserved declaration entity list, `occupied` contains reserved native names, `properSort` reads the provider-owned sort fact, and Exported chooses an admissible native spelling:

```go
	// Parameter names cannot shadow any generated declaration or support type.
	for _, id := range declarations {
		e := g.facts.Entities[id-1]
		local := map[string]bool{}
		for _, p := range e.Parameters {
			pe := g.facts.Entities[p.Entity-1]
			if pe.Kind != (semantics.Key{Owner: "horos", Name: "type-parameter"}) || len(pe.Declarations) != 1 {
				return semantic.ErrInvalidRequest
			}
			proper, err := g.properSort(g.sorts[semantics.Subject{Entity: p.Entity}])
			if err != nil {
				return err
			}
			if !proper {
				site := pe.Declarations[0]
				return &text.Diagnostic{Code: text.Unsupported, Message: "constructor parameters require a separate native representation profile", Node: &site}
			}
			name, err := Exported(pe.Name)
			if err != nil {
				return g.unsupported(e.Body, err.Error())
			}
			if occupied[name] || local[name] {
				return g.unsupported(e.Body, "native generic name collision: "+name)
			}
			local[name] = true
			g.names[p.Entity] = name
			g.targets[p.Entity] = g.dialect.Parameter(&g.arena, name, pe.Declarations[0])
			g.parameters[id] = append(g.parameters[id], g.targets[p.Entity])
		}
	}
```

The actual application lowering below recursively maps every argument to exactly one target type, then delegates the complete native reference. `referenceTarget` rejects computed and constructor-parameter callees. `tick` accounts work/depth/cancellation; `sites` stores source origins; `targets` maps semantic entities to reserved native AST nodes.

`internal/horosgen/generate.go:493-514` (implemented excerpt):

```go
func (g *generator[L]) referenceTarget(id semantics.TermID) (semantics.EntityID, error) {
	term := g.facts.Terms[id-1]
	target := term.Target
	if term.Kind == semantics.Application {
		callee := g.facts.Terms[term.Callee-1]
		if callee.Kind != semantics.Reference {
			return 0, g.unsupported(id, "computed constructor callees require a separate native representation profile")
		}
		target = callee.Target
		entity := g.facts.Entities[target-1]
		if entity.Kind != (semantics.Key{Owner: "horos", Name: "type-definition"}) {
			return 0, g.unsupported(id, "constructor parameter applications require a separate native representation profile")
		}
		if len(entity.Parameters) == 0 || len(entity.Parameters) != len(term.Arguments) {
			return 0, semantic.ErrInvalidRequest
		}
	}
	if target == 0 || g.targets[target] == 0 {
		return 0, semantic.ErrInvalidRequest
	}
	return target, nil
}
```

`internal/horosgen/generate.go:550-599` (implemented excerpt):

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

`transforms/horos-to-go/go/transform.go:1-104` (implemented excerpt):

```go
// Package golang maps analyzed Horos pagination types to Go syntax.
package golang

import (
	"context"
	"github.com/Bitspark/bittype/contracts/transform/go"
	g "github.com/Bitspark/bittype/internal/horosgen"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"go/token"
)

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

### Appendix B: proposed generation contracts from main

The following is the complete proposed Go-contract section from `docs/type-generation-plugins.md:373-896`, copied for this consultation. Its lifecycle and missing witness expansion are deliberately visible. The definitions describe interfaces, not a compiled package. Target convenience helpers are illustrative, with semantics described immediately around them.

#### Proposed Go contracts

Everything in this section is proposed. All shared declarations belong to package
`generation` at `contracts/generation/go`. `F` and `T` are the existing source and
target language marker types. Existing `transform`, `representation`, `semantics`
and `ast` types retain their current meanings. Opaque handles have construction
capabilities below, so independent source and target implementations can create
them without importing private Bittype packages.

### Configuration

```go
type PluginID string
type BindingID string

// ID identifies this configured binding, including its revision/configuration.
// Plugin.ID identifies the implementation for diagnostics.
type Mapping[F, T any] struct {
    ID     BindingID
    Plugin Plugin[F, T]
}

type Binding[F, T any] struct {
    Source  Selector[F]
    Mapping Mapping[F, T]
}

type Override[F, T any] struct {
    Source   Selector[F]
    Replaces BindingID
    With     Mapping[F, T]
}

type Config[F, T any] struct {
    Source     Source[F]
    Target     Target[T]
    Bindings   []Binding[F, T]
    Overrides  []Override[F, T]
    Structural *Mapping[F, T]
    Limits     Limits
}
```

Bindings contain builtin and library mappings uniformly. Structural is a selected
source/target integration plugin, not a hidden implementation in the kernel.
The factory belongs to the proposed `generation/syntax/go` package:

```go
func NewTransformer[F, T any](
    from, to representation.Language,
    config Config[F, T],
) (transform.ASTTransformer, error)
```

It captures the effective configuration and adapts through the existing
`transform.NewTransformer` bridge. It registers nothing. Composition provides
all required source evidence through Source.Contract and the invocation's
receipts. Plugins never invoke analysis implicitly. A mapping that needs more
facts requires an explicitly enlarged source contract and corresponding receipts;
there is no new per-plugin analyzer registry in this initial proposal.

### Use, publication and definition

```go
type Plugin[F, T any] interface {
    ID() PluginID
    Use(context.Context, UseRequest[F, T]) (Type[T], error)
}

type Publisher[F, T any] interface {
    Publish(
        context.Context,
        PublishRequest[F, T],
    ) (Symbol[T], error)
}

type Definer[F, T any] interface {
    Define(
        context.Context,
        DefineRequest[F, T],
    ) (Decl[T], error)
}

type UseRequest[F, T any] struct {
    Input   Ref[F]
    Context Context[F, T]
}

type PublishRequest[F, T any] struct {
    Input   SourceDecl[F] // Always open.
    Context Context[F, T]
}

type DefineRequest[F, T any] struct {
    Spec    DefinitionSpec[F, T]
    Output  Owned[T]
    Context Context[F, T]
}
```

Use creates a proper type expression such as `[]Person` or `Page<Person>`.
Publish exposes a requested open declaration such as Page<T> and returns a local
owned symbol. An open generic declaration is not itself a proper type expression;
calling Use on it for side effects would conflate these operations.

Publisher and Definer are optional capabilities of the selected Plugin. Publish
must explicitly arrange an owned alias/wrapper when the normal representation is
inline or external, or report unsupported publication. It cannot silently omit a
requested API root. Every requested definition requires Definer. There is no
Mode enum: an inline type can need a generated helper, and an owned declaration
can reference an external type.

### Owned declarations and helpers

```go
type DefinitionSubject[F any] struct { /* private tagged alternative */ }

func SourceDefinition[F any](open Ref[F]) DefinitionSubject[F]
func HelperDefinition[F any](localKey string) DefinitionSubject[F]

func (s DefinitionSubject[F]) Source() (Ref[F], bool)
func (s DefinitionSubject[F]) Helper() (string, bool)

// Explicit correspondence; TargetIndex is zero-based.
type BinderLink[F any] struct {
    Source      Ref[F]
    TargetIndex int
}

type DefinitionSpec[F, T any] struct {
    Subject DefinitionSubject[F]
    Header  DeclarationHeader[T]
    Binders []BinderLink[F]
}

type Owned[T any] struct {
    Symbol     Symbol[T]
    Parameters []Symbol[T]
    Scope      Scope[T]
}

type Context[F, T any] interface {
    Source() SourceSession[F]
    Syntax() Builder[T]

    Lower(context.Context, Ref[F]) (Type[T], error)

    Own(
        context.Context,
        DefinitionSpec[F, T],
    ) (Owned[T], error)

    External(
        context.Context,
        ExternalRef[T],
    ) (Symbol[T], error)
}
```

Own reserves the declaration and parameter symbols immediately, records a copied
specification, and queues the selected mapping's Define operation. It returns
without recursively executing Define. Repeated reservations must agree on the
header and binder correspondence. There are no captured definition closures.

A resource is owned by `(invocation, configured mapping, source declaration or
helper key)`. SourceDefinition always names the canonical open declaration, not
the first concrete application that encountered it. A builtin text wrapper can
use HelperDefinition("text-wrapper") without inventing a source declaration.
Helpers with the same local name in different mappings remain distinct.

The host enters Define with a fresh formal environment. BinderLink explicitly
pairs real source binders with reserved target parameter slots. A helper begins
without captured source binders; HelperDefinition requires an empty Binders list,
and its native parameters are explicit in the header and Output.Parameters. Nested
Lower preserves the lexical target environment while changing resource ownership
to the newly selected mapping.

Header.Parameters can reserve native slots, but BinderLink only associates a
source binder with a target slot. It does not specify a source-witness/native-type/
evidence tuple, a representation projection from a witness, or synthesis of extra
application arguments. Reserving more slots alone therefore does not implement
the companion's witness sketches. Such an encoding needs an explicit binder and
application contract shared by Use, Publish and Define; it remains deferred.
Do not invent source binders for generated native-only evidence parameters.

### Source references, binding and projection

```go
type Selector[F any] struct { /* private, comparable identity */ }
type Ref[F any] struct { /* private invocation and source token */ }

func NewSelector[F any](
    source representation.Language,
    identitySchema string,
    canonical []byte,
) (Selector[F], error)

func (s Selector[F]) SourceLanguage() representation.Language
func (s Selector[F]) IdentitySchema() string
func (s Selector[F]) Canonical() []byte // Copy.

// Created by the host for one invocation and supplied to Source.Open.
// The source adapter keeps it private.
type RefFactory[F any] interface {
    Intern(localToken string) (Ref[F], error)
    Token(Ref[F]) (string, error)
}

type SelectorClass uint8

const (
    BuiltinSelector SelectorClass = iota + 1
    DeclarationSelector
)

type HeadKind uint8

const (
    ParameterHead HeadKind = iota + 1
    BuiltinHead
    DeclarationHead
)

type Parameter[F any] struct {
    Binder   Ref[F]
    NameHint string
    Sort     semantics.Key
}

type Argument[F any] struct {
    Parameter Ref[F]
    Value     Ref[F]
}

type View[F any] struct {
    Input   Ref[F] // Effective reference after captured substitutions.
    Subject semantics.Subject
    Origins []ast.NodeRef

    Kind HeadKind

    Selector Selector[F] // BuiltinHead or DeclarationHead.
    Binder   Ref[F]      // ParameterHead: canonical real binder.
    Open     Ref[F]      // DeclarationHead: original open declaration.

    Arguments []Argument[F] // Effective declaration argument order.
}

type SourceDecl[F any] struct {
    Open       Ref[F]
    Selector   Selector[F]
    NameHint   string
    Parameters []Parameter[F]
    Body       Ref[F]
    Origins    []ast.NodeRef
}

type Source[F any] interface {
    Language() representation.Language
    Contract() transform.Contract

    ClassifySelector(Selector[F]) (SelectorClass, error)

    Open(
        context.Context,
        transform.Request[F],
        RefFactory[F],
    ) (SourceSession[F], error)
}

type SourceSession[F any] interface {
    // Initial profile: requested named declarations, in stable order.
    // Each reference is the canonical open declaration.
    Roots() []Ref[F]

    Inspect(Ref[F]) (View[F], error)

    // Accepts an open declaration or an admitted direct application.
    // Body carries the environment of the supplied reference.
    Declaration(Ref[F]) (SourceDecl[F], error)

    // Projects an actual child while preserving the parent's effective
    // substitution, lexical, and captured-assumption environment.
    Project(
        parent Ref[F],
        directChild semantics.Subject,
    ) (Ref[F], error)

    Model() semantics.Semantics[F]
}
```

Source here means the input-language collaborator, not the Files carrier
`source.Source[F]`. View is an operational view over the existing semantic graph,
not another canonical type algebra. Its three heads distinguish parameter uses,
builtin dispatch and declaration dispatch; they do not enumerate type libraries.

A selector is source/profile-qualified and uses a source-owned identity schema.
Horos's declaration key corresponds to `(decl.Digest, declaration name)`;
parameter identity adds the ordered slot. ClassifySelector interprets that
encoding without requiring the declaration to be in the current invocation.
The generic selection code compares identities without decoding Horos values.

Inspect applies captured substitutions while retaining declaration boundaries.
Declaration(view.Open) returns an open body; Declaration(use.Input) returns a
body in that application's environment. Project admits a real direct child and
preserves the parent's lexical/substitution/assumption context. It is not a
constructor for arbitrary graph addresses. Crossing a declaration boundary uses
Declaration. Horos-specific readers decode designated descriptor operands and
call Project; literal data in exact descriptors is never traversed as types.

Two environments remain separate: source binder -> effective source argument
belongs to SourceSession; source binder -> target parameter type belongs to the
lowering context. Matching parameter names or slice positions is insufficient.
RefFactory interns the complete source context and rejects references from other
invocations. Source adapters keep their factory private.

Roots means requested named API declarations, not every resolved dependency.
Source.Open admits source evidence and provides views; it does not eagerly reject
all recursive or unsupported bodies before representation selection. Required
negative or absent source admission is still a failure.

### Native fragments and symbols

```go
type Scope[T any] struct { /* private session identity */ }
type Symbol[T any] struct { /* private session identity */ }

type Fragment[T any] struct { /* immutable native syntax template */ }
type Type[T any] struct { /* fragment classified as a type use */ }
type Decl[T any] struct { /* fragment classified as a declaration */ }

func (t Type[T]) Fragment() Fragment[T]
func (d Decl[T]) Fragment() Fragment[T]
```

```go
type NameSlot[T any] struct {
    Attribute string
    Symbol    Symbol[T]
}

type Child[T any] struct {
    Role string
    Node Fragment[T]
}

type Link[T any] struct {
    Role   string
    Target Symbol[T] // Local resolved-link target only.
}

type NodeSpec[T any] struct {
    Kind       ast.Kind[T]
    Attributes []ast.Attribute
    Names      []NameSlot[T]
    Children   []Child[T]
    References []Link[T]

    Defines Symbol[T] // Zero, or the symbol anchored by this node.
    Origin  ast.Origin
}

type UseForm[T any] string // Target-owned construction vocabulary.

type SymbolUse[T any] struct {
    Form   UseForm[T]
    Symbol Symbol[T]
    Origin ast.Origin
}

type FragmentView[T any] struct {
    Node *NodeSpec[T]
    Use  *SymbolUse[T]
}

func (f Fragment[T]) View() FragmentView[T] // Copies; exactly one alternative.

type Builder[T any] interface {
    Node(NodeSpec[T]) (Fragment[T], error)
    Use(SymbolUse[T]) (Fragment[T], error)
    Type(Fragment[T]) (Type[T], error)
    Declaration(Fragment[T]) (Decl[T], error)
}
```

These fragments are immutable target syntax templates. Builder copies input
slices and enforces session/structural ownership. Type and Declaration mark the
intended native syntactic category; they do not add a language-independent
catalogue of valid native kinds. Native category/grammar admission remains with
the target adapter and existing checker.

A target package supplies typed convenience functions over this surface:

```go
gotarget.String(builder, origin)
gotarget.Slice(builder, element, origin)
gotarget.Named(builder, symbol, arguments, origin)
gotarget.DefinedType(builder, owned, body, origin)
```

SymbolUse is a target-owned construction placeholder, expanded into structured
native syntax during assembly. It is not a canonical type kind or runtime
binding. The generic kernel never appends generic arguments to arbitrary output;
the selected plugin and target helper construct the complete native type use.

Names as well as reference links stay symbolic until allocation. Every materialized
occurrence gets fresh containment nodes; reusable fragments do not share contained
AST nodes. Definitions and parameter symbols have one primary anchor. Free symbols
must remain in their native scope. For example, a Page<T> field fragment cannot be
moved into an unrelated declaration while retaining Page's T.

### Target session and construction capabilities

```go
type UnitKey string

type Placement struct {
    Unit  UnitKey
    Scope string // Target-configured declaration scope within/for the unit.
}

type ParameterHeader[T any] struct {
    Form     ast.Kind[T]
    NameHint string
}

type DeclarationHeader[T any] struct {
    Placement Placement
    Form      ast.Kind[T]
    NameHint  string
    Parameters []ParameterHeader[T]
}

type ExternalRef[T any] struct {
    Artifact string // Key into the caller-supplied locked artifact catalogue.
    Export   string // Target-owned export/member identity.
}

type ResourceKey struct { /* host-created ownership identity */ }

type Target[T any] interface {
    Language() representation.Language

    Open(
        context.Context,
        Storage[T],
    ) (TargetSession[T], error)
}

type TargetSession[T any] interface {
    EntryScope() Scope[T]

    Reserve(ResourceKey, DeclarationHeader[T]) (Owned[T], error)
    External(ExternalRef[T]) (Symbol[T], error)

    Builder(Scope[T]) Builder[T]
    Parameter(Scope[T], Symbol[T]) (Type[T], error)

    Define(Owned[T], Decl[T]) error

    Finish(
        context.Context,
        []Symbol[T], // Published roots.
    ) (ast.Set[T], error)
}
```

```go
type Limits struct {
    MaxNodes       int
    MaxDepth       int
    MaxUnits       int
    MaxWork        int
    MaxDefinitions int
    MaxImports     int
    MaxBytes       int
}

type Storage[T any] interface {
    Scope(parent Scope[T], stableKey string) (Scope[T], error)
    Symbol(scope Scope[T], stableKey string) (Symbol[T], error)
    Builder(scope Scope[T]) Builder[T]

    Limits() Limits
    ChargeOutput(nodes, bytes int) error
}
```

ResourceKey is comparable and supplied by the host, scoped to the active invocation.
The target also needs a stable ordering key for deterministic allocation:

```go
func (k ResourceKey) OrderKey() string
```

OrderKey is opaque to the target and excludes the invocation nonce. It distinguishes
the configured mapping and definition subject; it is not a structural fingerprint
or semantic equivalence claim. Source identity decoding stays in the source adapter.

Only the generation host receives TargetSession and Storage. Plugins receive the
restricted Context surface. The contract package supplies copying constructors
and handle issuers; Storage scopes and symbols are session-owned. Native adapters
own declaration/import scopes, spelling, placement, and artifact/export lookup.
They can capture those configurations in their constructors. No canonical mapping
is hidden in the target backend.

Placement identifies logical AST units and configured native declaration scopes.
The current ASTSet return boundary does not assign filesystem paths or publish a
Files Source. Parameter constraints are built in Define after reservation and
must agree with the native header. Arbitrary fragment-local type/lifetime binders
and arbitrary nested declaration scopes are deferred in the first profile.

ExternalRef identifies a supplied artifact/export, not an AST node. Without its
declaration tree, emit native import/include and qualified syntax rather than a
fabricated NodeRef. Import aliases are allocated per materialization unit under
native scope rules. The consumer retains the dependency environment; ASTSet is
neither a dependency lock nor a semantic receipt.

Limits cover actual materialized occurrences, expanded uses, definitions, imports,
units, payloads and work. They retain existing ceilings rather than counting only
reusable templates. Cancellation is cooperative; arbitrary Go plugins remain
trusted code. No cross-invocation cache or global mutable plugin state is proposed.


### Appendix C: organization and dependency direction

```text
contracts/generation/go/                  shared protocol
generation/syntax/go/                     agnostic selection and ownership
generation/syntax/source/horos/go/        Horos source adapter
generation/syntax/target/<language>/go/   native fragments, symbols, assembly
transforms/horos-to-<target>/go/           configured transformer factories
```

Type-owned integrations reside outside the Bittype Go module, in their owning libraries. Horos-owned integrations supply builtin mappings and structural declaration policy. Dependencies point integration -> Bittype, consumer -> integration + Bittype. Bittype factories never import an opinionated default bundle. Source and target implementations are collaborators within generation, not new registered stages. Stage-first Syntax/Files organization is an established repository convention. No directory in this proposed tree exists merely because it appears in this brief.

### Appendix D: primary prior art and limits of the analogy

[MLIR dialect conversion](https://mlir.llvm.org/docs/DialectConversion/#type-converter)
separates type conversion from value materialization and provides mappings from
one input type to several output types. Its signature conversion records how
original argument positions correspond to new positions, additional inputs, or
replacement values. That is relevant to coherent declaration/use elaboration;
it does not establish Horos membership or native generic constraints. Our
inference is that a shared correspondence can be useful without copying MLIR's
pattern selection, reverse registration precedence, materialization casts or
retry behavior. None of those mechanisms is part of the proposed type-only
Bittype boundary.

[GHC's typechecker evidence API](https://downloads.haskell.org/~ghc/10.0.0.20260917/docs/libraries/ghc-10.0.0.20260917/GHC-Tc-Types-Evidence.html)
represents distinct type/evidence abstractions and applications, with explicit
composition and scope. This is an example of elaboration preserving both the
introduced formals and their uses. The inference for our question is limited:
a header-only expansion misses the corresponding use transformation. We do not
propose executable dictionaries or Haskell instance resolution.

[Practical aspects of evidence-based compilation in System FC](https://www.microsoft.com/en-us/research/wp-content/uploads/2016/07/fc-new-tyco.pdf)
separates source inference from an explicitly typed intermediate language carrying
equality evidence, and discusses evidence-size costs. It suggests separating
established facts from their explicit elaborated representation, with expansion
budgets. Its equality coercions are not Horos membership, and it does not justify
adding a universal proof calculus or inference engine to Bittype. These are
architectural analogies for the expert to assess, not adopted implementation plans.

## Applied

This evaluation covers attempt 1, run
`run_aa6f8ad3-4bbc-4b25-a894-fda48f62c8e5`, against baseline
`f2179283150b1e7344844fc97503ca26fe40795c`. The frozen input is unchanged
(`6771ea6ed7f5c9b7698f1a842821c6b2dda5862eb4c63bb5b19edf8e35fb70d5`);
the genuine completed response is
`cde74f820e72037f88ef9e37cb4d51e7f5b3be504fab64c7062b00c47e8e5bc8`.
The living document differs by lifecycle and evaluation additions only.

Accepted actions below are documentation decisions in this conceptual task.
They do not assert that proposed APIs or native enforcement exist. The owner
prohibits implicit validation workflows; no compiler/test/build run was added.
The Validation column gives the source basis or explicitly unqualified boundary,
not an unperformed passing check. The paired synthesis will integrate 0009's
independently evaluated soundness limits before the proposal is delivered.

| # | Recommendation | Verdict | Action | Owner | Link | Validation |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | Separate source admission, representation selection, generic elaboration, native syntax and coordination | HOLDS | Document the elaboration responsibility in source/target integration; keep engine and kernel agnostic | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Current transformer boundary and proposed ownership support the separation; source inspection only |
| 2 | Use one shared formal/application convention, rather than independent private bookkeeping | HOLDS | Specify Introduce/Apply as paired conceptual operations for one open declaration | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Current one-to-one formal/application loops and BinderLink limitation shown in Appendix A/B |
| 3 | Do not mandate a universal source/payload/proof tuple | HOLDS | Keep support encoding and parameter shape integration-owned; additional discrimination is required only by the chosen promise/domain | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Leaf collision explains information loss, without proving one native encoding |
| 4 | Adopt formal/application and result-association laws | HOLDS | Include scoped assumptions and result composition, including nested TreePage<TreePage<T>> | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Existing memberships retain assumptions in internal/horosmodel/facts.go:50-71; new elaboration remains proposed |
| 5 | Replace binder-to-index with binder-to-scoped-product | HOLDS | Record the replacement seam; identify BinderLink as insufficient for the expanded proposal | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | docs/type-generation-plugins.md:497-556 explicitly lacks projections/extra arguments |
| 6 | Use a typed integration support carrier C | NEEDS ADAPTATION | Treat third Go type parameter and immutable support token as a candidate implementation, not a finalized universal public SPI | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Concrete support domain depends on0009; construction/immutability/lifetime API remains to be specified |
| 7 | Admit the whole source use before payload-driven traversal | HOLDS | Preserve ignored-argument obligations; lazy mapping selection may avoid payload expansion but never source admission | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Horos resolver admits module closure before native generation; current backend rejects unused parameters, so Ignore is a future case |
| 8 | Add an AdmittedUse coverage handle | NEEDS ADAPTATION | Define it as lookup/transport of established coverage, not an analysis call or new theorem | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | AdmitEvidence checks receipt association; DirectParameterBindings is structural only |
| 9 | Missing judgments require an expanded source contract/receipt | NEEDS ADAPTATION | State that an implemented provider/profile capability is required; adding configuration requirements alone produces no new facts | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | analyze/files/horos/go/analyze.go:37-72 reuses resolver facts and reports unknown obligations Unsupported |
| 10 | Do not bless arbitrary payload/support pairs by wrapping them | HOLDS | Require integration-owned association construction/admission and host scope/lifetime discipline; no universal semantic verifier | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Current AST receipts do not establish proposed native association; no new runtime implementation claimed |
| 11 | Describe each signature from the open declaration and fixed configuration | HOLDS | Prohibit first-use specialization and captured caller products/symbols | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Existing ownership proposal already keys definitions by canonical open declaration |
| 12 | Apply returns complete result representation, not only extra argument slices | HOLDS | Include full native use and result association in the conceptual convention | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Nested constrained applications expose missing result support; target-specific construction awaits0009 |
| 13 | Preserve frozen native formal constraints through Define | HOLDS | Target owns grammar/anchors and must not drop an unused argument's admission gate; proposed capability only | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Current Parameter takes only name/origin; current enforcement profile is Unsupported |
| 14 | Keep Lower as payload projection of the complete representation path | HOLDS | State that convenience payload access cannot bypass admission or association obligations | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Existing proposed Lower centralizes child dispatch; richer result not yet implemented |
| 15 | Use explicit native parameters for deferred helpers and scoped products | HOLDS | No captured caller environment, even through immutable references; cache keys include target formal environment | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Source and target environments are distinct; existing AST containment/reference rules require scoped anchors |
| 16 | Keep representation associations outside source Semantics | HOLDS | Store them in invocation-local generation state, with source receipts unchanged | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | semantics/go/model.go owns copied provider facts; engine discards source receipts across transform legs |
| 17 | Separate selector, source identity, BindingID, association and convention identities | HOLDS | Document distinct purposes and conservative source identity; different identities do not establish semantic inequality | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Existing module/declaration/parameter identities at internal/horosmodel/facts.go:73-90 |
| 18 | Qualify Native<S> by representation association or restrict coexistence | HOLDS | Require association-qualified witness, explicit family index, or a domain forbidding conflicting coexistence; do not choose a universal encoding | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Same source under explicit override can map to different native types; a single-valued projection needs this distinction |
| 19 | An override is replacement authority, not compatibility evidence | HOLDS | Preserve source membership separately from native family/constraint validity; reject incompatible changed associations | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Existing exact BindingID guard cannot prove semantic compatibility |
| 20 | Imported generic exports need explicit contracts | HOLDS | Document source/signature identity, convention/domain, support authority, result behavior and artifact/export association | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | ExternalRef currently names only Artifact/Export in the proposal |
| 21 | Distinguish baked-in representation dependencies from variable arguments | HOLDS | Compare relevant Text/List/optional/support requirements, not the entire mapping table; unrelated overrides remain independent | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Page.next uses Text independently of Page's generic argument |
| 22 | Treat export descriptors as producer claims and preserve ASTSet output | HOLDS | Supply descriptors with integrations initially; no mutable side channel/fake AST node/source fact for manifests | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | trees/ast/go/set.go and transform adapter return syntax only; identities do not prove claims |
| 23 | Preserve explicit failure and no weaker retry | HOLDS | Failure includes missing admission, association mismatch, dropped bounds, incompatible imported dependencies and scope escape | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Existing transformer bridge returns no usable partial output; no new fallback authorized |
| 24 | Start with a representation-only convention, then one justified enforcement domain | HOLDS | Document ordered implementation proposal and decisive examples; no implementation undertaken now | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Current representation/enforcement split makes the staged boundary explicit |
| 25 | Finalize public factories, convention construction and support guarantees together with0009 | DEFERRED | Trigger: a specific target/domain and support authority have been selected, and implementation is explicitly requested; then complete public issuer/accessor paths, immutable carrier rules, authoritative formal layout, callee acquisition and result-support demand before production code | codex-api-retirement-audit | [issue84](https://github.com/Bitspark/bittype/issues/84) | Advisor snippets omit constructors and cannot serve as a complete ready-to-compile SPI |
| 26 | Implement generic products, encoding, import contracts and proposed cases | DEFERRED | Trigger: owner authorizes implementation after paired conceptual discussion; no current generator or native guarantee changed | codex-api-retirement-audit | [issue84](https://github.com/Bitspark/bittype/issues/84) | Current Dialect and Unsupported enforcement remain intact |
| 27 | Defer higher-order/computed generics, mixed native binders, automatic adaptation and unrestricted recursion | DEFERRED | Trigger: a concrete separately scoped case exceeds the qualified first-order direct-application domain | codex-api-retirement-audit | [issue84](https://github.com/Bitspark/bittype/issues/84) | Current generator explicitly rejects those constructs; declaration reservation is not legality evidence |
| 28 | Use prior-art elaboration analogies without importing executable dictionaries, universal proofs or fallback | HOLDS | Document limited relevance of paired abstraction/application and generated-support budgets | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Sources in Appendix D support architectural analogies, not Horos/native soundness |
| 29 | Keep optional plugin effects optional | NEEDS ADAPTATION | Generic convention, publication and definition are separate demanded capabilities; specify callee acquisition for owned/inline/external results without forcing Own for every use | codex-api-retirement-audit | [PR86](https://github.com/Bitspark/bittype/pull/86) | Existing proposal permits mixed inline/external/owned effects; raw advisor snippets leave acquisition/issuance incomplete |

No recommendation is interpreted as permission to implement an enforcing target
or change the engine's registered capabilities. New generation-association
contracts remain proposed and must not be advertised in the provider matrix.

The documentation outcome is delivered in [PR86](https://github.com/Bitspark/bittype/pull/86):
the [paired synthesis](../docs/generic-elaboration.md), corrected
[plugin proposal](../docs/type-generation-plugins.md), corrected
[native constraint proposal](../docs/native-kind-constraints.md), and the current
[generation overview](../docs/generation.md). The peer reviewed the shared changes
against consultation0009 before integration. Applied here means the accepted
conceptual documentation actions are complete. No generator implementation,
native encoding, runtime binding or compiler qualification is claimed.
