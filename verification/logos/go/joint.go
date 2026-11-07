// Package logos checks questions issued by explicitly selected interpretations.
package logos

import (
	"context"
	"errors"
	"slices"

	association "github.com/Bitspark/bittype/contracts/verification/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/logos-contract/contract/go/typed"
)

var ErrInvalid = errors.New("invalid or foreign joint verification request")

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
	policy      typed.Policy
	limits      Limits
	authorities map[semantics.Key]*identity
}
type Composition struct{ state *composition }

// Problem is constructed by the selected interpretation from its typed subjects.
// The receiver admits its captures and locations before checking a candidate.
type Problem struct {
	Input        association.Input
	Requirements []association.Requirement
	Statement    typed.Statement
	Origins      []Origin
}
type identity struct{ key semantics.Key }
type Authority struct{ id *identity }

// Interpretation binds one typed projection to an unforgeable receiver selection.
type Interpretation[S any] struct {
	authority Authority
	describe  func(S) (Problem, error)
}

func Interpret[S any](key semantics.Key, describe func(S) (Problem, error)) Interpretation[S] {
	if !key.Valid() || describe == nil {
		return Interpretation[S]{}
	}
	return Interpretation[S]{Authority{&identity{key}}, describe}
}
func (i Interpretation[S]) Authority() Authority { return i.authority }

// A question retains the typed subject supplied to its selected interpreter.
// Its logical problem is frozen once, before any candidate can be submitted.
type issuance interface{
 issuer() Authority
 description() Problem
}
type issued[S any] struct{
 owner Authority
 subject S
 problem Problem
}
func(q *issued[S])issuer()Authority{return q.owner}
func(q *issued[S])description()Problem{return q.problem}
type Claim struct{state issuance}

func(i Interpretation[S])Ask(subject S)(Claim,error){
 if i.authority.id==nil{return Claim{},ErrInvalid}
 p,err:=i.describe(subject)
 if err!=nil{return Claim{},err}
 return Claim{&issued[S]{i.authority,subject,p.snapshot()}},nil
}

func(p Problem)snapshot()Problem{
 p.Input.Participants=slices.Clone(p.Input.Participants)
 for i:=range p.Input.Participants{
  entry:=&p.Input.Participants[i]
  entry.Evidence=slices.Clone(entry.Evidence)
 }
 p.Requirements=slices.Clone(p.Requirements)
 for i:=range p.Requirements{
  req:=&p.Requirements[i]
  req.Evidence=slices.Clone(req.Evidence)
  for j:=range req.Evidence{req.Evidence[j]=req.Evidence[j].Clone()}
 }
 p.Statement=cloneStatement(p.Statement)
 p.Origins=cloneOrigins(p.Origins)
 return p
}
// Answers identifies this particular issued question, including its subject.
func (r Request) Answers(q Claim) bool { return r.Valid() && r.state.claim == q.state }

type request struct {
	owner       *composition
	claim       issuance
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
func New(policy typed.Policy, limits Limits, selected ...Authority) (Composition, error) {
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
	if len(selected) > limits.Association.Evidence {
		return Composition{}, association.ErrLimit
	}
	authorities := make(map[semantics.Key]*identity, len(selected))
	for _, a := range selected {
		if a.id == nil || authorities[a.id.key] != nil {
			return Composition{}, ErrInvalid
		}
		bytes += len(a.id.key.Owner) + len(a.id.key.Name)
		if bytes > limits.Association.MetadataBytes {
			return Composition{}, association.ErrLimit
		}
		authorities[a.id.key] = a.id
	}
	return Composition{&composition{clonePolicy(policy), limits, authorities}}, nil
}
func (c Composition) Prepare(ctx context.Context, input Claim) (Request, error) {
	if ctx == nil || c.state == nil || input.state == nil {
		return Request{}, ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return Request{}, err
	}
	authority := input.state.issuer().id
	if authority == nil || c.state.authorities[authority.key] != authority {
		return Request{}, ErrInvalid
	}
	b:=input.state.description()
	a, err := association.Capture(ctx, b.Input, b.Requirements, c.state.limits.Association)
	if err != nil {
		return Request{}, err
	}
	count, bytes := 0, 0
	if len(b.Origins) > c.state.limits.Association.Sites {
		return Request{}, association.ErrLimit
	}
	for _, origin := range b.Origins {
		if !validLocation(origin.Location, b.Statement) {
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
	return Request{&request{c.state, input.state, a, cloneStatement(b.Statement), cloneOrigins(b.Origins)}}, nil
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
	return r.state.claim.issuer().id.key
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
