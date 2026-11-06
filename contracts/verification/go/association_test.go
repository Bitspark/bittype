package verification_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/verification/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type fixtureLanguage struct{}
func (fixtureLanguage) ID() language.ID { return "foreign-association-fixture" }
func (fixtureLanguage) Name() string { return "Foreign association fixture" }
func (fixtureLanguage) Check(ast.View[fixtureLanguage]) error { return nil }
func key(name string) semantics.Key { return semantics.Key{Owner: "fixture", Name: name} }

func participant(t *testing.T, name string, status semantic.Completion, verdict semantic.Verdict) (verification.Participant, verification.Requirement) {
	t.Helper()
	ctx := context.Background()
	definition, err := representation.NewLanguage[fixtureLanguage](fixtureLanguage{}, "syntax/v1")
	if err != nil { t.Fatal(err) }
	checker, err := check.NewChecker[fixtureLanguage](definition, fixtureLanguage{})
	if err != nil { t.Fatal(err) }
	tree, err := representation.Pack[fixtureLanguage](definition, ast.AST[fixtureLanguage]{ID: "same-address", Root: 1, Nodes: []ast.Node[fixtureLanguage]{{Kind: "root"}}})
	if err != nil { t.Fatal(err) }
	snapshot, err := semantic.Capture(ctx, checker, semantic.Input{Entry: ast.NodeRef{Tree: "same-address", Node: 1}, Trees: []representation.AST{tree}})
	if err != nil { t.Fatal(err) }
	contract := semantic.Contract{Profile: "facts/v1", Obligations: []semantic.Obligation{{Key: key("admission"), Domain: key("root")}}}
	op, err := semantic.NewOperation[fixtureLanguage](definition, contract, key("facts"), fixtureLanguage{}, func(_ context.Context, request semantic.Request[fixtureLanguage]) (semantic.Result, error) {
		return request.Publish(semantic.Publication{Reports: []semantic.Report{{Obligation: key("admission"), Status: status, Verdict: verdict}}})
	})
	if err != nil { t.Fatal(err) }
	settings := []semantic.Setting{{Key: key("profile"), Value: []byte("a")}}
	result, err := op.Run(ctx, snapshot, semantic.Options{Configuration: settings}, nil)
	if err != nil { t.Fatal(err) }
	p := verification.Participant{Key: key(name), Snapshot: snapshot, Evidence: []semantic.Result{result}}
	r := verification.Requirement{Key: p.Key, Language: definition, Evidence: []semantic.EvidenceRequirement{{Operation: key("facts"), Profile: "facts/v1", Configuration: settings, Reports: []semantic.Requirement{{Obligation: key("admission"), Verdict: semantic.Holds}}}}}
	return p, r
}

func fixture(t *testing.T) (verification.Input, []verification.Requirement) {
	t.Helper()
	a, ar := participant(t, "source", semantic.Complete, semantic.Holds)
	b, br := participant(t, "target", semantic.Complete, semantic.Holds)
	return verification.Input{Source: a.Key, Target: b.Key, Participants: []verification.Participant{a, b}}, []verification.Requirement{ar, br}
}

func TestAssociationOwnsConstructionAndViews(t *testing.T) {
	input, requirements := fixture(t)
	source := input.Participants[0].Snapshot
	a, err := verification.Capture(context.Background(), input, requirements, verification.Limits{})
	if err != nil { t.Fatal(err) }
	input.Participants[0].Evidence[0] = semantic.Result{}
	input.Participants[0].Snapshot = semantic.Snapshot{}
	requirements[0].Evidence[0].Configuration[0].Value[0] = 'z'
	requirements[0].Evidence[0].Reports[0].Verdict = semantic.Rejected
	view := a.Participants()
	view[0].Evidence[0] = semantic.Result{}
	view[0].Snapshot = semantic.Snapshot{}
	policy := a.Requirements()
	policy[0].Evidence[0].Configuration[0].Value[0] = 'z'
	p, ok := a.Participant(key("source"))
	if !ok || !p.Snapshot.Same(source) || !p.Evidence[0].Valid() { t.Fatal("capture or receipt lost") }
	if string(a.Requirements()[0].Evidence[0].Configuration[0].Value) != "a" || a.Requirements()[0].Evidence[0].Reports[0].Verdict != semantic.Holds { t.Fatal("mutable requirements escaped") }
	b := a
	a = verification.Association{}
	if a.Valid() || !b.Valid() || !b.Same(b) { t.Fatal("mutable outer handle") }
}

func TestAssociationRejectsForeignReceiptAtEqualAddresses(t *testing.T) {
	input, requirements := fixture(t)
	input.Participants[0].Evidence[0] = input.Participants[1].Evidence[0]
	a, err := verification.Capture(context.Background(), input, requirements, verification.Limits{})
	if !errors.Is(err, semantic.ErrInvalidRequest) || a.Valid() { t.Fatal(a.Valid(), err) }
}

func TestAssociationHasFreshIdentityAndQualifiedSites(t *testing.T) {
	input, requirements := fixture(t)
	a, err := verification.Capture(context.Background(), input, requirements, verification.Limits{})
	if err != nil { t.Fatal(err) }
	b, err := verification.Capture(context.Background(), input, requirements, verification.Limits{})
	if err != nil { t.Fatal(err) }
	if a.Same(b) { t.Fatal("reconstruction reused association identity") }
	for _, part := range []string{"source", "target"} {
		if !a.Contains(verification.Site{Participant: key(part), Node: ast.NodeRef{Tree: "same-address", Node: 1}}) { t.Fatal(part) }
	}
	for _, site := range []verification.Site{
		{Participant: key("foreign"), Node: ast.NodeRef{Tree: "same-address", Node: 1}},
		{Participant: key("source"), Node: ast.NodeRef{Tree: "foreign", Node: 1}},
		{Participant: key("source"), Node: ast.NodeRef{Tree: "same-address", Node: 2}},
		{},
	} {
		if a.Contains(site) || !errors.Is(a.CheckSites([]verification.Site{site}), verification.ErrInvalid) { t.Fatal(site) }
	}
}

func TestAssociationPreservesUnavailableEvidence(t *testing.T) {
	for _, status := range []semantic.Completion{semantic.Complete, semantic.Incomplete, semantic.Unsupported} {
		p, requirement := participant(t, "source", status, semantic.Undetermined)
		input := verification.Input{Source: p.Key, Target: p.Key, Participants: []verification.Participant{p}}
		a, err := verification.Capture(context.Background(), input, []verification.Requirement{requirement}, verification.Limits{})
		if !errors.Is(err, semantic.ErrEvidenceUnavailable) || a.Valid() { t.Fatal(status, err) }
		if p.Evidence[0].Reports()[0].Status != status { t.Fatal("receipt was rewritten") }
	}
}

func TestAssociationBoundsAggregateInputsAndRejectsMalformedScope(t *testing.T) {
	input, requirements := fixture(t)
	limits := verification.DefaultLimits()
	limits.InputNodes = 1 // Either snapshot fits alone; together they exceed it.
	if _, err := verification.Capture(context.Background(), input, requirements, limits); !errors.Is(err, verification.ErrLimit) { t.Fatal(err) }
	limits = verification.DefaultLimits()
	limits.Evidence = 1
	if _, err := verification.Capture(context.Background(), input, requirements, limits); !errors.Is(err, verification.ErrLimit) { t.Fatal(err) }
	limits = verification.DefaultLimits()
	limits.MetadataBytes = 1
	if _, err := verification.Capture(context.Background(), input, requirements, limits); !errors.Is(err, verification.ErrLimit) { t.Fatal(err) }
	input.Target = key("absent")
	if _, err := verification.Capture(context.Background(), input, requirements, verification.Limits{}); !errors.Is(err, verification.ErrInvalid) { t.Fatal(err) }
	input.Target = input.Source
	input.Participants[1] = input.Participants[0]
	if _, err := verification.Capture(context.Background(), input, requirements, verification.Limits{}); !errors.Is(err, verification.ErrInvalid) { t.Fatal(err) }
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := verification.Capture(ctx, input, requirements, verification.Limits{}); !errors.Is(err, context.Canceled) { t.Fatal(err) }
	if _, err := verification.Capture(nil, input, requirements, verification.Limits{}); !errors.Is(err, verification.ErrInvalid) { t.Fatal(err) }
}
