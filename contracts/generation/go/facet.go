package generation

// PROTOTYPE (bittype#196): facets, products and reference categories (rows
// 26, 27 and 29). A facet is a kind-identified, versioned, validated value on
// a contribution. The kernel knows no kind: it compares issuers, identities,
// versions and targets, and hands a value back to its kind's owner.

import (
	"fmt"
)

// FacetKindReference is the stable qualified identity of a facet kind and the
// version of its contract. It leaves the process.
type FacetKindReference struct {
	Identity string `json:"identity"`
	Version  string `json:"version"`
}

// FacetTransport is the transport descriptor of a kind whose values cross a
// process boundary (row 26a): a media type and the kind owner's encoder. The
// kernel calls the encoder when it serializes a plan and does not read the
// result.
type FacetTransport struct {
	MediaType string
	Encode    func(value any) ([]byte, error)
}

type FacetKindDescription struct {
	Identity, Version string
	// Validate is the validator of the kind. It may be nil for a kind whose
	// Go type is its whole contract.
	Validate  func(value any) error
	Transport *FacetTransport
}

// FacetKind is a typed extension contract that belongs to the product the
// facet serves. Its issuer is an identity of the process: a kind made by
// another owner under the same identity and version is another kind.
type FacetKind struct {
	issuer      *identity
	description FacetKindDescription
}

func NewFacetKind(d FacetKindDescription) (FacetKind, error) {
	if d.Identity == "" || d.Version == "" || d.Transport != nil && (d.Transport.MediaType == "" || d.Transport.Encode == nil) {
		return FacetKind{}, ErrInvalid
	}
	return FacetKind{fresh(), d}, nil
}
func (k FacetKind) Valid() bool { return k.issuer != nil }
func (k FacetKind) Reference() FacetKindReference {
	return FacetKindReference{k.description.Identity, k.description.Version}
}

// Facet is a value of one kind for one contribution and one target. Publisher
// is who registered the facet package: the kernel admits a facet that names a
// registered contribution of the same publisher (O8).
type Facet struct {
	kind         FacetKind
	publisher    string
	contribution ContributionIdentity
	target       string
	value        any
}

// Issue is the typed way to make a facet value in the process. It runs the
// kind's validator.
func (k FacetKind) Issue(publisher string, contribution ContributionIdentity, target string, value any) (Facet, error) {
	if !k.Valid() || publisher == "" || contribution.Publisher == "" || contribution.Key == "" || target == "" {
		return Facet{}, ErrInvalid
	}
	if k.description.Validate != nil {
		if err := k.description.Validate(value); err != nil {
			return Facet{}, fmt.Errorf("%w: facet of kind %s: %v", ErrInvalid, k.description.Identity, err)
		}
	}
	return Facet{k, publisher, contribution, target, value}, nil
}

// Read hands a facet's value to the owner of its kind, and to nobody else.
func (k FacetKind) Read(f Facet) (any, bool) {
	if !k.Valid() || f.kind.issuer != k.issuer {
		return nil, false
	}
	return f.value, true
}
func (f Facet) Valid() bool                        { return f.kind.Valid() }
func (f Facet) Kind() FacetKindReference           { return f.kind.Reference() }
func (f Facet) Publisher() string                  { return f.publisher }
func (f Facet) Contribution() ContributionIdentity { return f.contribution }
func (f Facet) Target() string                     { return f.target }

// ErrFacet: a facet under a required kind's identity that does not match the
// kind: another issuer, another target or another contract version.
var ErrFacet = fmt.Errorf("%w: a facet does not match the kind the product requires", ErrScope)

// TypesProduct is the product every request asks for: native types. A
// request names further products to ask for bindings (O5).
const TypesProduct = "types"

// Product is what a generation request asks generation to make possible. It
// states, per kind of subject, the facet kinds that every selected
// contribution of that kind must have. A contribution without a standalone
// value adapter is expressible (row 14): a product requires another kind of a
// field position than of a use.
type Product struct {
	Identity string
	Requires map[SubjectKind][]FacetKind
}

// FacetOf returns the facet of one kind that the plan's composition holds for
// a contribution, checked as the product check does.
func matchFacet(facets []Facet, kind FacetKind, contribution ContributionIdentity, target string) (Facet, error) {
	want := kind.Reference()
	var found *Facet
	for i, f := range facets {
		if f.contribution != contribution || f.kind.description.Identity != want.Identity {
			continue
		}
		if !PrototypeMutants["facet-ignores-issuer"] && f.kind.issuer != kind.issuer {
			return Facet{}, fmt.Errorf("%w: kind %s of contribution %s/%s comes from another issuer", ErrFacet, want.Identity, contribution.Publisher, contribution.Key)
		}
		if !PrototypeMutants["facet-ignores-target"] && f.target != target {
			return Facet{}, fmt.Errorf("%w: kind %s of contribution %s/%s is for target %s", ErrFacet, want.Identity, contribution.Publisher, contribution.Key, f.target)
		}
		if !PrototypeMutants["facet-ignores-version"] && f.kind.description.Version != want.Version {
			return Facet{}, fmt.Errorf("%w: kind %s of contribution %s/%s has contract version %s, the product requires %s", ErrFacet, want.Identity, contribution.Publisher, contribution.Key, f.kind.description.Version, want.Version)
		}
		found = &facets[i]
	}
	if found == nil {
		return Facet{}, ErrMissingFacet
	}
	return *found, nil
}

// ---------------------------------------------------------------------------
// Reference categories (row 29, guarantee 3).

type ReferenceCategory uint8

const (
	// TypeReference: the native type of a selected dependency.
	TypeReference ReferenceCategory = iota + 1
	// BinderReference: a scoped native binder.
	BinderReference
	// AdapterReference: an adapter argument of an application, or the adapter
	// parameter that a binder brings into its declaration's scope.
	AdapterReference
	// RequirementReference: a generic requirement and its evidence.
	RequirementReference
	// DeclarationReference: a generated declaration, by its logical symbol.
	DeclarationReference
	// OwnedReference: a declaration that a contribution owns.
	OwnedReference
	// MemberReference: a generated member, tied to the native type that contains it.
	MemberReference
)

func (c ReferenceCategory) String() string {
	return [...]string{"", "type", "binder", "adapter", "requirement", "declaration", "owned", "member"}[c]
}

// PlanReference is a reference a reader holds into one plan. Only the plan
// issues one. The contract gains no expression of any language: what is done
// with a reference is the product's.
type PlanReference struct {
	plan     *identity
	category ReferenceCategory
	node     string
	index    int
	scope    string
}

func (r PlanReference) Category() ReferenceCategory { return r.category }
func (r PlanReference) Node() string                { return r.node }
func (r PlanReference) Index() int                  { return r.index }

func (p RepresentationPlan) reference(category ReferenceCategory, node string, index int, scope string) (PlanReference, error) {
	if _, ok := p.index[node]; !ok {
		return PlanReference{}, fmt.Errorf("%w: no such node", ErrScope)
	}
	return PlanReference{p.owner, category, node, index, scope}, nil
}

// TypeOf refers to the native type of a node.
func (p RepresentationPlan) TypeOf(node string) (PlanReference, error) {
	return p.reference(TypeReference, node, 0, "")
}

// MemberOf refers to the member with the given ordinal of the native type of
// a node. The reference is tied to that node.
func (p RepresentationPlan) MemberOf(node string, ordinal int) (PlanReference, error) {
	n, ok := p.Node(node)
	if !ok || ordinal < 0 || ordinal >= len(n.Members) {
		return PlanReference{}, fmt.Errorf("%w: no such member", ErrScope)
	}
	return p.reference(MemberReference, node, ordinal, "")
}

// BinderOf and AdapterOf refer to a binder of a declaration and to the
// adapter that the binder brings into the declaration's scope. Their scope is
// the declaration's logical symbol.
func (p RepresentationPlan) BinderOf(declaration string, index int) (PlanReference, error) {
	return p.scoped(BinderReference, declaration, index)
}
func (p RepresentationPlan) AdapterOf(declaration string, index int) (PlanReference, error) {
	return p.scoped(AdapterReference, declaration, index)
}
func (p RepresentationPlan) scoped(category ReferenceCategory, declaration string, index int) (PlanReference, error) {
	n, ok := p.Node(declaration)
	if !ok || n.Kind != "declaration" || index < 0 || index >= n.Binders {
		return PlanReference{}, fmt.Errorf("%w: no such binder", ErrScope)
	}
	return p.reference(category, declaration, index, n.Symbol)
}

// DeclarationOf refers to a planned declaration by its logical symbol.
func (p RepresentationPlan) DeclarationOf(symbol string) (PlanReference, error) {
	s, ok := p.Symbol(symbol)
	if !ok {
		return PlanReference{}, fmt.Errorf("%w: no such symbol", ErrScope)
	}
	category := DeclarationReference
	if s.Role == "owned" {
		category = OwnedReference
	}
	return PlanReference{p.owner, category, s.Node, 0, ""}, nil
}

// Check is guarantee 3: a reference has the right plan, category and scope.
// scope is the logical symbol of the declaration in whose scope the reader
// uses the reference; it is compared for the scoped categories only.
func (p RepresentationPlan) Check(r PlanReference, want ReferenceCategory, scope string) error {
	if !PrototypeMutants["reference-ignores-plan"] && (r.plan == nil || r.plan != p.owner) {
		return fmt.Errorf("%w: a reference of another plan", ErrScope)
	}
	if !PrototypeMutants["reference-ignores-category"] && r.category != want {
		return fmt.Errorf("%w: a %s reference where a %s reference is expected", ErrScope, r.category, want)
	}
	if !PrototypeMutants["reference-ignores-scope"] && (r.category == BinderReference || r.category == AdapterReference) && r.scope != scope {
		return fmt.Errorf("%w: a binder of %q outside its scope, in %q", ErrScope, r.scope, scope)
	}
	return nil
}

// Member resolves a checked member reference to the spelled name of the
// member in the generated native type. Go and Rust have no Files stage for
// generated declarations (row 24a), so the spelled name is what a Go reader
// has; a TypeScript reader takes the realized name from the realization.
func (p RepresentationPlan) Member(r PlanReference) (string, error) {
	if err := p.Check(r, MemberReference, ""); err != nil {
		return "", err
	}
	n, _ := p.Node(r.node)
	if r.index < 0 || r.index >= len(n.Members) {
		return "", ErrScope
	}
	return n.Members[r.index], nil
}
