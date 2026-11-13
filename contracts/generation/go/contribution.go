package generation

// PROTOTYPE of the contribution contract (docs/contribution-contract.md,
// bittype#196). This file states what a contribution is: its identities, its
// domain, its rule in two steps, its projector, the payload codec, evidence
// and refusals. The planner is in contribution_plan.go. The mapping contract
// of main stays beside it on this branch (PROTOTYPE.md says why).

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Bitspark/bittype/trees/ast/go"
)

// ---------------------------------------------------------------------------
// Identities (rows 15 to 17, O4, O19, O24).

// ContributionIdentity is the publisher and the publisher's local key.
type ContributionIdentity struct{ Publisher, Key string }

// Meaning is the family's meaning: the same in every target.
type Meaning struct{ Identity, Revision string }

// ContributionReference names what was selected: the identity, the family's
// meaning and the contribution's own representation revision.
type ContributionReference struct {
	Identity       ContributionIdentity
	Meaning        Meaning
	Representation string
}

// Valid reports whether every part is stated. A reference that lacks either
// part of the third axis is no reference (row 17, O20).
func (r ContributionReference) Valid() bool {
	return r.Identity.Publisher != "" && r.Identity.Key != "" && r.Meaning.Identity != "" && r.Meaning.Revision != "" && r.Representation != ""
}
func (r ContributionReference) String() string {
	return fmt.Sprintf("%s/%s@%s/%s#%s", r.Identity.Publisher, r.Identity.Key, r.Meaning.Identity, r.Meaning.Revision, r.Representation)
}

// SchemaReference is the payload schema axis: which schema types and encodes
// the recorded decisions of a contribution.
type SchemaReference struct{ Identity, Version string }

func (s SchemaReference) Valid() bool { return s.Identity != "" && s.Version != "" }

// ---------------------------------------------------------------------------
// The applicability domain (row 1).

type DomainKind uint8

const (
	// ExactDomain: the uses of one exact declaration.
	ExactDomain DomainKind = iota + 1
	// FormDomain: one builtin form of the source.
	FormDomain
	// DeclarationPolicyDomain: every named declaration without an exact mapping.
	DeclarationPolicyDomain
	// RoleDomain: a role that the source integration issues for a field position.
	RoleDomain
	// SpecializationDomain: one application of one declaration (row 21).
	SpecializationDomain
)

// Domain is what a contribution answers for. The kernel selects by it and
// never by registration order.
type Domain struct {
	kind     DomainKind
	selector Selector
	role     string
	actuals  []Selector
}

func ForExact(declaration Selector) Domain { return Domain{kind: ExactDomain, selector: declaration} }
func ForForm(form Selector) Domain         { return Domain{kind: FormDomain, selector: form} }
func ForDeclarationPolicy() Domain         { return Domain{kind: DeclarationPolicyDomain} }
func ForRole(role string) Domain           { return Domain{kind: RoleDomain, role: role} }
func ForSpecialization(declaration Selector, actuals ...Selector) Domain {
	return Domain{kind: SpecializationDomain, selector: declaration, actuals: append([]Selector(nil), actuals...)}
}
func (d Domain) Kind() DomainKind { return d.kind }

// Key is the domain as selection compares it.
func (d Domain) Key() string {
	switch d.kind {
	case ExactDomain:
		return "exact:" + d.selector.Key()
	case FormDomain:
		return "form:" + d.selector.Key()
	case DeclarationPolicyDomain:
		return "policy"
	case RoleDomain:
		return "role:" + d.role
	case SpecializationDomain:
		keys := []string{d.selector.Key()}
		for _, a := range d.actuals {
			keys = append(keys, a.Key())
		}
		return "specialization:" + strings.Join(keys, "|")
	}
	return ""
}

// ---------------------------------------------------------------------------
// Roles at a field position (row 12).

// RoleIssuer is the source integration's authority to issue role selectors.
// A source hands its issuer to the kernel with what it opens (Source.Roles).
type RoleIssuer[F any] struct{ owner *identity }

// Role is a role selector: the source's statement that one operand of one
// parent use has a role. The kernel reads the key to select and nothing else.
type Role[F any] struct {
	owner           *identity
	key             string
	parent, operand *occurrence[F]
}

func NewRoleIssuer[F any]() RoleIssuer[F] { return RoleIssuer[F]{fresh()} }
func (i RoleIssuer[F]) Issue(key string, parent, operand Use[F]) Role[F] {
	if i.owner == nil || key == "" || parent.o == nil || operand.o == nil {
		return Role[F]{}
	}
	return Role[F]{i.owner, key, parent.o, operand.o}
}
func (r Role[F]) Valid() bool { return r.owner != nil && r.key != "" }
func (r Role[F]) Key() string { return r.key }

// ---------------------------------------------------------------------------
// Evidence (rows 7 to 9, O2, O13, O23).

// EvidenceKind is the form the prototype fixes: a qualified identity and a
// version, owned by the vocabulary of the target or product that the fact
// concerns. The kernel compares kinds and reads no content.
type EvidenceKind struct{ Identity, Version string }

type Outcome uint8

const (
	Proved Outcome = iota + 1
	Refuted
	Unknown
	Exhausted
	Deferred
)

func (o Outcome) String() string {
	return [...]string{"", "proved", "refuted", "unknown", "exhausted", "deferred"}[o]
}

// Fact is a typed value of an evidence kind. A rule states it beside its
// decision, or a contribution declares it at registration for what its
// revision fixes. Content is the kind owner's; the kernel carries it unread.
type Fact struct {
	Kind    EvidenceKind
	Outcome Outcome // Proved, Refuted or Exhausted
	Content string
}

func (k EvidenceKind) Proved(content string) Fact  { return Fact{k, Proved, content} }
func (k EvidenceKind) Refuted(content string) Fact { return Fact{k, Refuted, content} }
func (k EvidenceKind) Exhausted() Fact             { return Fact{k, Exhausted, ""} }

// Demand asks for the facts of one kind about a dependency. Require declares,
// where the dependency is bound to a binder, that the binder's actuals must
// prove the fact: the requirement of the binder.
type Demand struct {
	Kind    EvidenceKind
	Require bool
}

// Answer is what the planner issued to one rule about one dependency. It is
// bound to the asking contribution, the subject, the dependency, the
// dependency's whole contribution reference and the generic environment. A
// rule reads it through the step it was issued to (DecideStep.Present).
type Answer struct {
	kind        EvidenceKind
	outcome     Outcome
	content     string
	established bool // Deferred: the declared requirement establishes the fact
	bound       answerBinding
}
type answerBinding struct {
	asker       ContributionReference
	subject     string
	dependency  int
	reference   ContributionReference
	environment string
}

func (a Answer) Kind() EvidenceKind { return a.kind }

// ---------------------------------------------------------------------------
// The rule (row 3): one round.

type SubjectKind uint8

const (
	UseSubject SubjectKind = iota + 1
	DeclarationSubject
	FieldSubject
)

// FieldPosition is the subject of a contribution whose domain is a role: the
// parent record's use, the operand use the kernel minted and the role selector
// the source issued. Stated is what the asking rule already knew from the
// source (O3), for a member fragment the member's name.
type FieldPosition[F any] struct {
	Parent, Operand Use[F]
	Role            Role[F]
	Stated          string
}

// Subject is the checked subject of one selection.
type Subject[F any] struct {
	Kind        SubjectKind
	Use         UseView[F]
	Declaration DeclarationView[F]
	Field       FieldPosition[F]
}

type requestKind uint8

const (
	childRequest requestKind = iota + 1
	fieldRequest
	helperRequest
)

// Request is one thing a rule names in step 1.
type Request[F any] struct {
	kind       requestKind
	use        Use[F]
	role       Role[F]
	stated     string
	demands    []Demand
	helper     string
	conversion bool
}

// Conversion names a child use that the rule needs for conversion only: the
// native type does not refer to it. The planner resolves it when the request
// asks for bindings and leaves it out for types alone (O5, row 20): "Only
// demanded native actuals need mappings."
func Conversion[F any](child Use[F]) Request[F] {
	return Request[F]{kind: childRequest, use: child, conversion: true}
}

// Reference names a child use the rule needs to refer to.
func Reference[F any](child Use[F]) Request[F] { return Request[F]{kind: childRequest, use: child} }

// Evidence names a child use about whose selection the rule needs answers.
func Evidence[F any](child Use[F], demands ...Demand) Request[F] {
	return Request[F]{kind: childRequest, use: child, demands: append([]Demand(nil), demands...)}
}

// Field names a field position: an operand of the subject with the role the
// source issued for that pair. stated travels with the request (O3).
func Field[F any](operand Use[F], role Role[F], stated string) Request[F] {
	return Request[F]{kind: fieldRequest, use: operand, role: role, stated: stated}
}

// SharedHelper names whichever shared helper the composition selects for key
// (row 25). The prototype states the request kind and refuses it as
// unsupported: whose key the request uses is not specified.
func SharedHelper[F any](key string) Request[F] { return Request[F]{kind: helperRequest, helper: key} }

type NameStep[F, R any] struct {
	Subject Subject[F]
	Reader  R
}

type answerState struct {
	issued   [][]Answer
	expected []answerBinding
	failure  error
}

type DecideStep[F, R any] struct {
	Subject Subject[F]
	Reader  R
	state   *answerState
}

// Answers returns the answers issued about dependency i, in demand order.
func (s DecideStep[F, R]) Answers(i int) []Answer {
	if s.state == nil || i < 0 || i >= len(s.state.issued) {
		return nil
	}
	return append([]Answer(nil), s.state.issued[i]...)
}

// Present reads an answer as an answer about dependency i. This is the entry
// point of guarantee 4: an answer that was not issued to this step for this
// dependency, in this environment and for the dependency's contribution
// reference, is rejected with ErrScope, and the planning of the subject fails
// even if the rule drops the error.
func (s DecideStep[F, R]) Present(i int, a Answer) (Outcome, string, error) {
	if s.state == nil || i < 0 || i >= len(s.state.expected) {
		return 0, "", ErrScope
	}
	want, got := s.state.expected[i], a.bound
	if PrototypeMutants["evidence-ignores-dependency"] {
		got.dependency, got.subject = want.dependency, want.subject
	}
	if PrototypeMutants["evidence-ignores-environment"] {
		got.environment = want.environment
	}
	if PrototypeMutants["evidence-ignores-reference"] {
		got.reference = want.reference
	}
	if PrototypeMutants["evidence-ignores-representation-revision"] {
		got.reference.Representation = want.reference.Representation
	}
	if PrototypeMutants["evidence-ignores-meaning-revision"] {
		got.reference.Meaning = want.reference.Meaning
	}
	if got != want {
		err := fmt.Errorf("%w: an evidence answer was not issued for this dependency (issued for %s in %q about %s, presented for %s in %q about %s)", ErrScope, a.bound.subject, a.bound.environment, a.bound.reference, want.subject, want.environment, want.reference)
		if s.state.failure == nil {
			s.state.failure = err
		}
		return 0, "", err
	}
	return a.outcome, a.content, nil
}

// Evidence is the answer of one kind about dependency i. Unknown where the
// rule did not ask for the kind. For Deferred, established says whether the
// binder's declared requirement establishes the fact.
func (s DecideStep[F, R]) Evidence(i int, kind EvidenceKind) (outcome Outcome, content string, established bool) {
	for _, a := range s.Answers(i) {
		if a.kind == kind {
			o, c, err := s.Present(i, a)
			if err != nil {
				return Unknown, "", false
			}
			return o, c, a.established
		}
	}
	return Unknown, "", false
}

// Decision is what a rule returns: a value of the generated plain Go type of
// its payload schema, and beside it the facts that others may ask of it.
type Decision struct {
	Payload any
	Facts   []Fact
}

// Rule is the part of a contribution that the planner runs for one selection.
// Name may be nil for a rule without requests. Neither function is given the
// planner, a target type or native intent.
type Rule[F, R any] struct {
	Name   func(context.Context, NameStep[F, R]) ([]Request[F], error)
	Decide func(context.Context, DecideStep[F, R]) (Decision, error)
}

// ---------------------------------------------------------------------------
// The payload codec (row 28, O9).

// PayloadCodec is the codec of one payload schema. Freeze encodes the value a
// rule returned, validates the encoding against the schema and returns the
// canonical bytes of the frozen value. View decodes frozen bytes into a fresh
// typed value. The contract package imports neither horos nor ontos: the
// kernel holds the bytes as an opaque handle and never reads them.
type PayloadCodec interface {
	Schema() SchemaReference
	Freeze(value any) ([]byte, error)
	View(frozen []byte) (any, error)
}

// ErrPayload: a rule returned a value that its payload schema does not admit.
var ErrPayload = fmt.Errorf("%w: payload outside its schema", ErrInvalid)

// ---------------------------------------------------------------------------
// The projector (row 4).

// Member is one member of a projected native type, for member references. A
// member is named, or it stands for the members of a dependency's fragment.
type Member struct {
	Name     string
	Fragment int // 1+index of the dependency whose fragment supplies the members, when Name is empty
}

func FragmentMembers(dependency int) Member { return Member{Fragment: dependency + 1} }

// Projected is the native intent of a use or of a field position. A use
// projects a type. A field position projects a member fragment: the member or
// members of the native record that represent the field.
type Projected[T any] struct {
	Type     Expression[T]
	Fragment []Expression[T]
	Members  []Member
}

// DependencyView is the typed view of a resolved request.
type DependencyView[T any] struct {
	Type     Expression[T]   // a child use: the native type selected for it
	Fragment []Expression[T] // a field position: the member fragment
	Stated   string
	Origin   ast.NodeRef
}

type projectionState[T any] struct {
	payload      any
	dependencies []DependencyView[T]
	stated       string
	origin       ast.NodeRef
	binders      []ast.NodeRef
	owned        func(name string, args []Expression[T]) (Expression[T], error)
}

// Projection is everything a projector may read: the typed view decoded from
// the frozen decision, typed views of the dependencies the planner resolved,
// views of the symbols the plan issued, and diagnostic sites. It has no source
// reader, no use, no declaration handle and no selector.
type Projection[T any] struct{ state *projectionState[T] }

func (p Projection[T]) Payload() any      { return p.state.payload }
func (p Projection[T]) Dependencies() int { return len(p.state.dependencies) }
func (p Projection[T]) Dependency(i int) DependencyView[T] {
	if i < 0 || i >= len(p.state.dependencies) {
		return DependencyView[T]{}
	}
	return p.state.dependencies[i]
}

// Stated is what the asking rule stated in the request that selected this
// contribution, for a field position the member's name.
func (p Projection[T]) Stated() string      { return p.state.stated }
func (p Projection[T]) Origin() ast.NodeRef { return p.state.origin }

// Binders is the number of binders of the declaration being projected, and
// BinderOrigin the declaring site of one of them.
func (p Projection[T]) Binders() int { return len(p.state.binders) }
func (p Projection[T]) BinderOrigin(i int) ast.NodeRef {
	if i < 0 || i >= len(p.state.binders) {
		return ast.NodeRef{}
	}
	return p.state.binders[i]
}

// Owned refers to a declaration that the contribution owns (row 25), applied
// to args. The planner issues its symbol by the owned key.
func (p Projection[T]) Owned(name string, args ...Expression[T]) (Expression[T], error) {
	if p.state.owned == nil {
		return Expression[T]{}, ErrScope
	}
	return p.state.owned(name, args)
}

type Projector[T any] func(context.Context, Projection[T]) (Projected[T], error)
type DeclarationProjector[T any] func(context.Context, Projection[T]) (DeclarationDescription[T], error)

// OwnedDescription is a declaration that a contribution introduces and that
// represents no source declaration. Its key is the contribution reference and
// Name. The spelling in Declaration.Name is not the key.
type OwnedDescription[T any] struct {
	Name        string
	Declaration NativeDeclaration[T]
}

// ---------------------------------------------------------------------------
// The contribution (row 1).

// Contribution is one independently selectable representation policy for one
// target and one profile over a declared domain.
type Contribution[F, T, R any] struct {
	Reference ContributionReference
	// Target is a stated field (O4): the identity of the target language.
	Target  string
	Profile string
	Domain  Domain
	Schema  SchemaReference
	Rule    Rule[F, R]
	// Project computes the native intent of a use or a field position. Declare
	// computes it for the declaration policy domain. Exactly one is set.
	Project Projector[T]
	Declare DeclarationProjector[T]
	// Facts are the facts that the contribution's revision fixes for every
	// selection, declared at registration (O23).
	Facts []Fact
	Owned []OwnedDescription[T]
}

func (c Contribution[F, T, R]) validate() error {
	if !c.Reference.Valid() || c.Target == "" || c.Profile == "" || !c.Schema.Valid() || c.Rule.Decide == nil || c.Domain.Key() == "" {
		return fmt.Errorf("%w: incomplete contribution %s", ErrInvalid, c.Reference)
	}
	policy := c.Domain.kind == DeclarationPolicyDomain
	if (c.Declare != nil) != policy || (c.Project != nil) == policy {
		return fmt.Errorf("%w: contribution %s has the wrong projector for its domain", ErrInvalid, c.Reference)
	}
	names := map[string]bool{}
	for _, o := range c.Owned {
		if o.Name == "" || names[o.Name] {
			return fmt.Errorf("%w: contribution %s owns an unnamed or repeated declaration", ErrInvalid, c.Reference)
		}
		names[o.Name] = true
	}
	return nil
}

// ---------------------------------------------------------------------------
// Refusals (row 10, O17, O26). The form the prototype fixes: one type with a
// class. Each class has a value that errors.Is matches. The existing values
// are refined and none is replaced.

var (
	ErrRefusedByDesign        = errors.New("refused by design")
	ErrRequirementUnsatisfied = errors.New("requirement unsatisfied")
	ErrRequirementUnresolved  = errors.New("requirement unresolved")
	ErrMissingFacet           = errors.New("missing facet")
	ErrMissingContribution    = errors.New("missing contribution")
)

type RefusalClass uint8

const (
	RefusalUnsupported RefusalClass = iota + 1
	RefusalByDesign
	RefusalRequirementUnsatisfied
	RefusalRequirementUnresolved
	RefusalMissingFacet
	RefusalAmbiguousContribution
	RefusalMissingContribution
	RefusalLimit
)

// Value is the contract value of a class.
func (c RefusalClass) Value() error {
	switch c {
	case RefusalUnsupported:
		return ErrUnsupported
	case RefusalByDesign:
		return ErrRefusedByDesign
	case RefusalRequirementUnsatisfied:
		return ErrRequirementUnsatisfied
	case RefusalRequirementUnresolved:
		return ErrRequirementUnresolved
	case RefusalMissingFacet:
		return ErrMissingFacet
	case RefusalAmbiguousContribution:
		return ErrAmbiguous
	case RefusalMissingContribution:
		return ErrMissingContribution
	case RefusalLimit:
		return ErrLimit
	}
	return ErrInvalid
}

// The kernel's own reasons. A reason explains; it never adds a class.
const (
	KernelOwner                = "bittype/kernel"
	ReasonRecursiveReference   = "recursive reference"
	ReasonEvidenceCycle        = "evidence cycle"
	ReasonSpecialization       = "explicit specialization is not built (#166)"
	ReasonMemberFragment       = "the target's member vocabulary does not admit the member fragment"
	ReasonNoMemberVocabulary   = "the target states no member vocabulary"
	ReasonSharedHelper         = "shared helper requests are not built in the prototype"
	ReasonExactWithoutPolicy   = "an exact mapping publishes no declaration"
	ReasonRequirement          = "an actual does not discharge a binder's requirement"
	ReasonNoContribution       = "no contribution answers the domain"
	ReasonSeveralContributions = "more than one contribution answers the domain"
	ReasonFacet                = "the selected contribution has no facet of a kind the product requires"
)

// Refusal is a contribution's or the planner's refusal of one subject.
type Refusal struct {
	Class RefusalClass
	// Owner qualifies the reason: a publisher, a vocabulary or the kernel.
	Owner, Reason string
	// The subject: its source identity and site.
	Subject string
	Site    ast.NodeRef
	// The selected contribution, where one was selected, with target and profile.
	Contribution    ContributionReference
	Target, Profile string
	// Path is the path of dependencies from the root, as source identities.
	Path []string
	// Requirement is the requirement concerned, where there is one.
	Requirement EvidenceKind
	// Application names the application of a refused specialization.
	Application string
}

// Refuse is how a rule refuses. A rule may use the first four classes.
func Refuse(class RefusalClass, owner, reason string) *Refusal {
	return &Refusal{Class: class, Owner: owner, Reason: reason}
}
func (r *Refusal) Error() string {
	text := fmt.Sprintf("%v: %s: %s", r.Class.Value(), r.Owner, r.Reason)
	if r.Subject != "" {
		text += " (subject " + r.Subject + ")"
	}
	if r.Application != "" {
		text += " (application " + r.Application + ")"
	}
	if r.Contribution.Valid() {
		text += " (contribution " + r.Contribution.String() + ")"
	}
	if len(r.Path) > 0 {
		text += " (path " + strings.Join(r.Path, " > ") + ")"
	}
	return text
}
func (r *Refusal) Is(target error) bool { return target == r.Class.Value() }

// RuleFailure is an unexpected error of a rule or a projector: an
// implementation failure and no refusal. It matches no refusal value.
type RuleFailure struct {
	Contribution ContributionReference
	Err          error
}

func (f *RuleFailure) Error() string {
	return "contribution " + f.Contribution.String() + " failed: " + f.Err.Error()
}

// As reaches the failure's cause. The failure does not unwrap, so that the
// error of a rule cannot add a refusal class.
func (f *RuleFailure) As(target any) bool { return errors.As(f.Err, target) }

// PrototypeMutants switches single properties of the kernel off, so that a
// test of the prototype can show that it detects their absence. It is an
// instrument of the throwaway branch and no part of the contract.
var PrototypeMutants = map[string]bool{}
