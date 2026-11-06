package verification

import (
	"context"
	"errors"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	association "github.com/Bitspark/bittype/contracts/verification/go"
	p "github.com/Bitspark/bittype/provider/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"github.com/Bitspark/logos-contract/contract/go/typed"
)

// Inject defects at the private trusted interpretation boundary. No production
// constructor accepts these arbitrary statements or permits this mutation.
func faultyInterpretation(t *testing.T, statement typed.Statement) Claim {
	t.Helper()
	language, err := representation.NewLanguage[s.Language](p.Provider{}, "syntax/v1")
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[s.Language](language, p.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	tree, err := (p.Provider{}).Parse(context.Background(), parse.Input{SourceID: "input", TreeID: "input", Text: "package example; type Text string"})
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
	input := association.Input{Source: sourceKey(), Target: sourceKey(), Participants: []association.Participant{{Key: sourceKey(), Snapshot: snapshot}}}
	location := typed.Location{Phase: typed.ExpectedPhase}
	if len(statement.Declarations) > 0 {
		location = typed.Location{Phase: typed.DeclarationPhase, Index: 0, HasIndex: true}
	}
	origin := Origin{Location: location, Sites: []association.Site{site(sourceKey(), snapshot.Entry())}}
	return Claim{&claim{input: input, requirements: []association.Requirement{{Key: sourceKey(), Language: language}}, question: question("fault-injection"), statement: statement, origins: []Origin{origin}}}
}
func TestInvalidSignatureIsAnUnchangedLogosOutcome(t *testing.T) {
	statement := typed.Statement{Foundation: typed.Foundation, Declarations: []typed.Declaration{typed.Define("bad", typed.Bool(), typed.Zero())}, Expected: typed.Bool()}
	claim := faultyInterpretation(t, statement)
	c, err := New(typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker}, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	r, err := c.Prepare(context.Background(), claim)
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
	c, err := New(typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker}, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	r, err := c.Prepare(context.Background(), claim)
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
	c, _ := New(typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker}, DefaultLimits())
	claim.state.origins[0].Sites[0].Participant = targetKey()
	if r, err := c.Prepare(context.Background(), claim); !errors.Is(err, association.ErrInvalid) || r.Valid() {
		t.Fatal(err)
	}
	claim.state.origins[0].Sites[0].Participant = sourceKey()
	claim.state.origins[0].Location = typed.Location{Phase: typed.DeclarationPhase, Index: 99, HasIndex: true}
	if r, err := c.Prepare(context.Background(), claim); !errors.Is(err, ErrInvalid) || r.Valid() {
		t.Fatal("accepted foreign logical location", err)
	}
	claim.state.origins[0].Location = typed.Location{Phase: typed.ExpectedPhase}
	limits := DefaultLimits()
	limits.Association.Sites = 1
	c, _ = New(typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker}, limits)
	claim.state.origins = append(claim.state.origins, claim.state.origins[0])
	if r, err := c.Prepare(context.Background(), claim); !errors.Is(err, association.ErrLimit) || r.Valid() {
		t.Fatal(err)
	}
}
