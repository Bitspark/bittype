package logos

import (
	"context"
	"errors"
	"testing"

	checkgolang "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	association "github.com/Bitspark/bittype/contracts/verification/go"
	parsegolang "github.com/Bitspark/bittype/parse/syntax/golang/provider/go"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"github.com/Bitspark/logos-contract/contract/go/typed"
)

// Exercise checker failures using an independently supplied interpretation.
var fixture = Interpret(semantics.Key{Owner: "test", Name: "fault"}, func(p *Problem) (Problem, error) { return *p, nil })
var src = semantics.Key{Owner: "test", Name: "source"}

func faultyInterpretation(t *testing.T, statement typed.Statement) *Problem {
	t.Helper()
	language, err := representation.NewLanguage[s.Language](s.Language{}, "syntax/v1")
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[s.Language](language, checkgolang.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	tree, err := parsegolang.Parse(context.Background(), parse.Input{SourceID: "input", TreeID: "input", Text: "package example; type Text string"})
	if err != nil {
		t.Fatal(err)
	}
	packed, err := representation.Pack[s.Language](language, tree.ASTView())
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := semantic.Capture(context.Background(), checker, semantic.Input{Entry: ast.NodeRef{Tree: "input", Node: tree.Abstract.Root}, Trees: []representation.AST{packed}})
	if err != nil {
		t.Fatal(err)
	}
	input := association.Input{Source: src, Target: src, Participants: []association.Participant{{Key: src, Snapshot: snapshot}}}
	location := typed.Location{Phase: typed.ExpectedPhase}
	if len(statement.Declarations) > 0 {
		location = typed.Location{Phase: typed.DeclarationPhase, Index: 0, HasIndex: true}
	}
	origin := Origin{Location: location, Sites: []association.Site{{Participant: src, Node: snapshot.Entry()}}}
	return &Problem{Input: input, Requirements: []association.Requirement{{Key: src, Language: language}}, Statement: statement, Origins: []Origin{origin}}
}
func TestInvalidSignatureIsAnUnchangedLogosOutcome(t *testing.T) {
	statement := typed.Statement{Foundation: typed.Foundation, Declarations: []typed.Declaration{typed.Define("bad", typed.Bool(), typed.Zero())}, Expected: typed.Bool()}
	claim := faultyInterpretation(t, statement)
	c, err := New(typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker}, DefaultLimits(), fixture.Authority())
	if err != nil {
		t.Fatal(err)
	}
	r, err := c.Prepare(context.Background(), ask(t, claim))
	if err != nil {
		t.Fatal("preparation must not fake a logical precheck", err)
	}
	out, err := c.Check(context.Background(), r, typed.True())
	if err != nil {
		t.Fatal(err)
	}
	expected := typed.Check(statement, typed.True(), r.Policy(), r.Limits().Checker)
	if !out.Valid() || out.Outcome().Status() != typed.InvalidStatus || out.Outcome().Code() != expected.Code() || out.Outcome().Location() != expected.Location() {
		t.Fatal("invalid signature was rewritten")
	}
	d := out.Diagnostics()
	if len(d) != 1 || len(d[0].Sites) != 1 || d[0].Location.Phase != typed.DeclarationPhase {
		t.Fatal(d)
	}
	d[0].Sites[0].Node = ast.NodeRef{}
	if out.Diagnostics()[0].Sites[0].Node.Node == 0 {
		t.Fatal("diagnostic storage escaped")
	}
}

type cancelBeforePublication struct {
	context.Context
	calls int
}

func (c *cancelBeforePublication) Err() error {
	c.calls++
	if c.calls >= 2 {
		return context.Canceled
	}
	return nil
}
func TestCancellationBeforePublicationDoesNotMintResult(t *testing.T) {
	claim := faultyInterpretation(t, typed.Statement{Foundation: typed.Foundation, Expected: typed.Bool()})
	c, err := New(typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker}, DefaultLimits(), fixture.Authority())
	if err != nil {
		t.Fatal(err)
	}
	r, err := c.Prepare(context.Background(), ask(t, claim))
	if err != nil {
		t.Fatal(err)
	}
	ctx := &cancelBeforePublication{Context: context.Background()}
	out, err := c.Check(ctx, r, typed.True())
	if !errors.Is(err, context.Canceled) || out.Valid() {
		t.Fatal("published after cancellation", err)
	}
	if out, err := c.Check(context.Background(), r, typed.True()); err != nil || !out.Outcome().IsAccepted() {
		t.Fatal("cancellation revoked existing request", err)
	}
}
func TestPreparationRejectsForeignSitesAndBoundsWholeMap(t *testing.T) {
	claim := faultyInterpretation(t, typed.Statement{Foundation: typed.Foundation, Expected: typed.Bool()})
	c, _ := New(typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker}, DefaultLimits(), fixture.Authority())
	claim.Origins[0].Sites[0].Participant = semantics.Key{Owner: "test", Name: "other"}
	if r, err := c.Prepare(context.Background(), ask(t, claim)); !errors.Is(err, association.ErrInvalid) || r.Valid() {
		t.Fatal(err)
	}
	claim.Origins[0].Sites[0].Participant = src
	claim.Origins[0].Location = typed.Location{Phase: typed.DeclarationPhase, Index: 99, HasIndex: true}
	if r, err := c.Prepare(context.Background(), ask(t, claim)); !errors.Is(err, ErrInvalid) || r.Valid() {
		t.Fatal("accepted foreign logical location", err)
	}
	claim.Origins[0].Location = typed.Location{Phase: typed.ExpectedPhase}
	limits := DefaultLimits()
	limits.Association.Sites = 1
	c, _ = New(typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker}, limits, fixture.Authority())
	claim.Origins = append(claim.Origins, claim.Origins[0])
	if r, err := c.Prepare(context.Background(), ask(t, claim)); !errors.Is(err, association.ErrLimit) || r.Valid() {
		t.Fatal(err)
	}
}

func ask(t *testing.T, p *Problem) Claim {
	t.Helper()
	claim, err := fixture.Ask(p)
	if err != nil {
		t.Fatal(err)
	}
	return claim
}
func TestIssuerFreezesSubjectProjectionAndRequiresExactAuthority(t *testing.T) {
	p := faultyInterpretation(t, typed.Statement{Foundation: typed.Foundation, Expected: typed.Bool()})
	claim := ask(t, p)
	p.Origins[0].Sites[0].Node = ast.NodeRef{}
	p.Statement.Expected = typed.Nat()
	p.Input.Participants = nil
	p.Requirements[0].Key = semantics.Key{}
	other := Interpret(src, func(*Problem) (Problem, error) { return Problem{}, ErrInvalid })
	// Matching the question key also cannot substitute one authority for another.
	foreign := Interpret(fixture.Authority().id.key, func(*Problem) (Problem, error) { return Problem{}, ErrInvalid })
	policy := typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker}
	c, err := New(policy, DefaultLimits(), fixture.Authority(), other.Authority())
	if err != nil {
		t.Fatal(err)
	}
	req, err := c.Prepare(context.Background(), claim)
	if err != nil || !req.Answers(claim) {
		t.Fatal(err)
	}
	out, err := c.Check(context.Background(), req, typed.True())
	if err != nil || !out.Outcome().IsAccepted() {
		t.Fatal(err)
	}
	c, err = New(policy, DefaultLimits(), foreign.Authority())
	if err != nil {
		t.Fatal(err)
	}
	if r, err := c.Prepare(context.Background(), claim); !errors.Is(err, ErrInvalid) || r.Valid() {
		t.Fatal("foreign authority", err)
	}
	if _, err := New(policy, DefaultLimits(), fixture.Authority(), foreign.Authority()); !errors.Is(err, ErrInvalid) {
		t.Fatal("duplicate interpretation", err)
	}
}
