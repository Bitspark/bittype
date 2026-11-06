// Package verification is the explicit Horos/Go static representation composition.
// It has no runtime bindings and installs no language or checker registrations.
package verification

import (
	"context"
	"errors"
	"slices"
	"strings"

	association "github.com/Bitspark/bittype/contracts/verification/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/logos-contract/contract/go/typed"
)

var ErrInvalid = errors.New("invalid or foreign joint verification request")
var ErrCoverage = errors.New("static representation is outside the selected mapping domain")

type Origin struct {
	Location typed.Location
	Sites    []association.Site
}
type Diagnostic struct {
	Location typed.Location
	Code     typed.Code
	Detail   string
	Sites    []association.Site
}
type Limits struct {
	Association association.Limits
	Checker     typed.Limits
}

func DefaultLimits() Limits { return Limits{association.DefaultLimits(), typed.DefaultLimits()} }

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

func cloneStatement(s typed.Statement) typed.Statement {
	s.Declarations = slices.Clone(s.Declarations)
	s.Context = slices.Clone(s.Context)
	return s
}
func clonePolicy(p typed.Policy) typed.Policy {
	p.Dependencies = slices.Clone(p.Dependencies)
	return p
}
func cloneOrigins(in []Origin) []Origin {
	out := slices.Clone(in)
	for i := range out {
		out[i].Sites = slices.Clone(in[i].Sites)
	}
	return out
}
func logicalOrigins(declarations []typed.Declaration, sites []association.Site) []Origin {
	out := []Origin{{Location: typed.Location{Phase: typed.ExpectedPhase}, Sites: slices.Clone(sites)}}
	for i, d := range declarations {
		name := d.Name()
		var selected []association.Site
		for _, s := range sites {
			source := (strings.HasPrefix(name, "source.") || strings.HasPrefix(name, "h.") || strings.Contains(name, ".h.") || name == "argument.member") && s.Participant == sourceKey()
			target := (strings.HasPrefix(name, "native.") || strings.HasPrefix(name, "n.") || strings.Contains(name, ".n.")) && s.Participant == targetKey()
			mapping := strings.HasPrefix(name, "law.") || strings.Contains(name, ".law.") || name == "family.theorem"
			if source || target || mapping {
				selected = append(selected, s)
			}
		}
		if len(selected) > 0 {
			out = append(out, Origin{Location: typed.Location{Phase: typed.DeclarationPhase, Index: uint32(i), HasIndex: true}, Sites: selected})
		}
	}
	return out
}

func validLocation(location typed.Location, statement typed.Statement) bool {
	switch location.Phase {
	case typed.DeclarationPhase:
		return location.HasIndex && uint64(location.Index) < uint64(len(statement.Declarations))
	case typed.ContextPhase:
		return location.HasIndex && uint64(location.Index) < uint64(len(statement.Context))
	case typed.ExpectedPhase:
		return !location.HasIndex && location.Index == 0
	default:
		return false
	}
}

// New captures receiver policy. Unknown foundations/dependencies and checker
// resource exhaustion remain Logos outcomes; this constructor does not precheck
// or silently expand the receiver's policy/budget.
func New(policy typed.Policy, limits Limits) (Composition, error) {
	if !limits.Association.Valid() {
		return Composition{}, ErrInvalid
	}
	if len(policy.Dependencies) > limits.Association.Evidence {
		return Composition{}, association.ErrLimit
	}
	bytes := len(policy.Foundation) + len(policy.Checker)
	for _, dep := range policy.Dependencies {
		if len(dep) > limits.Association.MetadataBytes-bytes {
			return Composition{}, association.ErrLimit
		}
		bytes += len(dep)
	}
	if bytes > limits.Association.MetadataBytes {
		return Composition{}, association.ErrLimit
	}
	return Composition{&composition{clonePolicy(policy), limits}}, nil
}
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
func (r Request) Valid() bool             { return r.state != nil }
func (r Request) Same(other Request) bool { return r.Valid() && r.state == other.state }
func (r Request) Association() association.Association {
	if !r.Valid() {
		return association.Association{}
	}
	return r.state.association
}
func (r Request) Question() semantics.Key {
	if !r.Valid() {
		return semantics.Key{}
	}
	return r.state.claim.question
}
func (r Request) Statement() typed.Statement {
	if !r.Valid() {
		return typed.Statement{}
	}
	return cloneStatement(r.state.statement)
}
func (r Request) Policy() typed.Policy {
	if !r.Valid() {
		return typed.Policy{}
	}
	return clonePolicy(r.state.owner.policy)
}
func (r Request) Limits() Limits {
	if !r.Valid() {
		return Limits{}
	}
	return r.state.owner.limits
}
func (r Request) Origins() []Origin {
	if !r.Valid() {
		return nil
	}
	return cloneOrigins(r.state.origins)
}
func (r Result) Valid() bool { return r.state != nil }
func (r Result) Request() Request {
	if !r.Valid() {
		return Request{}
	}
	return r.state.request
}
func (r Result) Candidate() typed.Expr {
	if !r.Valid() {
		return typed.Expr{}
	}
	return r.state.candidate
}
func (r Result) Outcome() typed.Outcome {
	if !r.Valid() {
		return typed.Outcome{}
	}
	return r.state.outcome
}
func (r Result) Diagnostics() []Diagnostic {
	if !r.Valid() || r.state.outcome.IsAccepted() {
		return nil
	}
	o := r.state.outcome
	d := Diagnostic{Location: o.Location(), Code: o.Code(), Detail: o.Detail()}
	for _, origin := range r.state.request.state.origins {
		if origin.Location == d.Location {
			d.Sites = append(d.Sites, origin.Sites...)
		}
	}
	return []Diagnostic{d}
}
