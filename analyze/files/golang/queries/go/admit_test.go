package queries

import (
	"bytes"
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	analyze "github.com/Bitspark/bittype/analyze/files/golang/provider/go"
	checkgolang "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	check "github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	resolution "github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	parsegolang "github.com/Bitspark/bittype/parse/syntax/golang/provider/go"
	resolvegolang "github.com/Bitspark/bittype/resolve/files/golang/provider/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func admittedModel(t *testing.T, text string) n.Model {
	t.Helper()
	ctx := context.Background()
	parsed, err := parsegolang.Parse(ctx, parse.Input{SourceID: "0", TreeID: "0", Text: text})
	if err != nil {
		t.Fatal(err)
	}
	tree := parsed.ASTView()
	definition, err := representation.NewLanguage[s.Language](s.Language{}, "go-syntax-1")
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[s.Language](definition, checkgolang.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := representation.Pack[s.Language](definition, tree)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := semantic.Capture(ctx, checker, semantic.Input{Entry: ast.NodeRef{Tree: tree.ID, Node: tree.Root}, Trees: []representation.AST{envelope}})
	if err != nil {
		t.Fatal(err)
	}
	resolver, err := resolution.NewResolver[s.Language](snapshot.Language(), n.ResolutionContract(), resolvegolang.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: n.Configuration()})
	if err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[s.Language](snapshot.Language(), n.AnalysisContract(), analyze.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	result, err := analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: n.Configuration()}})
	if err != nil {
		t.Fatal(err)
	}
	model, err := n.Open(result)
	if err != nil {
		t.Fatal(err)
	}
	return model
}

// An admission that cannot go on returns when its context ends. The test holds
// the state's lock, as an earlier admission abandoned inside go/types would,
// and ends the context only once the admission is seen waiting for that lock.
func TestAdmitReturnsWhenContextEndsWhileItCannotGoOn(t *testing.T) {
	model := admittedModel(t, "package example; type Text string; type Box[T any] struct { Value T }")
	q, err := NewQueries(context.Background(), model)
	if err != nil {
		t.Fatal(err)
	}
	box, err := model.Lookup("Box")
	if err != nil {
		t.Fatal(err)
	}
	text, err := model.Lookup("Text")
	if err != nil {
		t.Fatal(err)
	}
	if admission, err := q.Admit(context.Background(), box, text); err != nil || !admission.Admitted() {
		t.Fatal(admission, err)
	}

	q.state.mu.Lock()
	defer q.state.mu.Unlock() // lets the abandoned admission end
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	returned := make(chan error, 1)
	go func() {
		_, err := q.Admit(ctx, box, text)
		returned <- err
	}()
	// The bounds only turn a hang into a failure; nothing is inferred from time.
	patience := time.After(time.Minute)
	stacks := make([]byte, 1<<20)
	for !bytes.Contains(stacks[:runtime.Stack(stacks, true)], []byte("Queries.admit")) {
		select {
		case err := <-returned:
			t.Fatal("the admission returned although the lock is held:", err)
		case <-patience:
			t.Fatal("the admission never reached the lock")
		default:
			runtime.Gosched()
		}
	}
	cancel()
	select {
	case err := <-returned:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Minute):
		t.Fatal("the admission waited for a lock that is not released")
	}
}
