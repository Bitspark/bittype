package consumer_test

import (
 "github.com/Bitspark/bittype/contracts/analysis/go"
 nativequery "github.com/Bitspark/bittype/analyze/files/golang/go"
	"context"
	"go/types"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	gp "github.com/Bitspark/bittype/provider/golang/go"
	hp "github.com/Bitspark/bittype/provider/horos/go"
	resolve "github.com/Bitspark/bittype/resolve/files/horos/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	hs "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	v "github.com/Bitspark/bittype/verification/files/horos-to-go/go"
	h "github.com/Bitspark/bittype/verification/files/source/horos/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	"github.com/Bitspark/logos-contract/contract/go/typed"
	core "github.com/bitspark/ontos/core/go"
)

type installedSyntax[L any] interface {
	language.Provider
	Check(ast.View[L]) error
	Parse(context.Context, parse.Input) (cst.CST[L], error)
}

func installedCapture[L any](t *testing.T, p installedSyntax[L], profile, source string) semantic.Snapshot {
	t.Helper()
	definition, err := representation.NewLanguage[L](p, profile)
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[L](definition, p)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := p.Parse(context.Background(), parse.Input{SourceID: "installed", TreeID: "installed", Text: source})
	if err != nil {
		t.Fatal(err)
	}
	tree := parsed.ASTView()
	packed, err := representation.Pack[L](definition, tree)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := semantic.Capture(context.Background(), checker, semantic.Input{Entry: ast.NodeRef{Tree: tree.ID, Node: tree.Root}, Trees: []representation.AST{packed}})
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

// A real installed consumer uses public projections and the public checker with
// no workspace, private module, replacement, engine registration or hidden stage.
func TestInstalledJointTextBox(t *testing.T) {
	ctx := context.Background()
	snapshot := installedCapture[hs.Language](t, hp.Provider{}, hs.Profile, "module installed; import h from horos.kinds pin \"sha256:"+decl.PreludePin().String()+"\"; type TextBox<T:h.Text> = {value:T;};")
	resolver, err := resolution.NewResolver[hs.Language](snapshot.Language(), resolve.Contract(), hp.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	receipt, err := resolver.Resolve(ctx, snapshot, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	source, err := h.Inspect(ctx, snapshot, receipt)
	if err != nil {
		t.Fatal(err)
	}
	box, err := source.Declaration(decl.DeclarationID{Module: source.Root(), Name: "TextBox"})
	if err != nil {
		t.Fatal(err)
	}
	bound, err := box.ParameterKind(0)
	if err != nil {
		t.Fatal(err)
	}
	boolean, err := source.Proper(core.NewTuple(core.NewAtom([]byte("horos/scalar")), core.NewAtom([]byte("bool"))))
	if err != nil {
		t.Fatal(err)
	}
	nonmember, err := source.Satisfies(ctx, boolean, bound)
	if err != nil || nonmember.Verdict() != decl.Refuted {
		t.Fatalf("authoritative nonmembership: %v %v", nonmember.Verdict(), err)
	}
	target, err := checkedNative(ctx, installedCapture[gs.Language](t, gp.Provider{}, "go-syntax-1", "package installed; type TextBox[T any] struct { Value T }"))
	if err != nil {
		t.Fatal(err)
	}
	nativeBox, err := target.Lookup("TextBox")
	if err != nil {
		t.Fatal(err)
	}
	nativeBool, err := target.Basic(types.Bool)
	if err != nil {
		t.Fatal(err)
	}
	admission, err := nativeAdmission(t,ctx,target,nativeBox,nativeBool)
	if err != nil {
		t.Fatal(err)
	}
	claim, err := v.Counterexample(box, nativeBox, nonmember, admission)
	if err != nil {
		t.Fatal(err)
	}
	checker, err := v.New(typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker, Dependencies: v.MappingDependencies()}, v.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	request, err := checker.Prepare(ctx, claim)
	if err != nil {
		t.Fatal(err)
	}
	result, err := checker.Check(ctx, request, v.CounterexampleCandidate())
	if err != nil || !result.Outcome().IsAccepted() {
		t.Fatalf("joint check: %v %s %s", err, result.Outcome().Code(), result.Outcome().Detail())
	}
}

// Each semantic role is explicit; opening the typed view compiles nothing.
func checkedNative(ctx context.Context,snapshot semantic.Snapshot) (n.Model,error) {
 resolver,err:=gp.NewFilesResolver(snapshot.Language()); if err!=nil { return n.Model{},err }
 resolved,err:=resolver.Resolve(ctx,snapshot,semantic.Options{Configuration:n.Configuration()}); if err!=nil { return n.Model{},err }
 analyzer,err:=gp.NewFilesAnalyzer(snapshot.Language()); if err!=nil { return n.Model{},err }
 analyzed,err:=analyzer.Analyze(ctx,analysis.Input{Snapshot:snapshot,Resolution:&resolved,Options:semantic.Options{Configuration:n.Configuration()}}); if err!=nil { return n.Model{},err }
 return n.Open(analyzed)
}
func nativeAdmission(t *testing.T,ctx context.Context,model n.Model,constructor n.Type,args ...n.Type) (nativequery.Admission,error) {
 t.Helper(); queries,err:=nativequery.NewQueries(ctx,model); if err!=nil { return nativequery.Admission{},err }; return queries.Admit(ctx,constructor,args...)
}
