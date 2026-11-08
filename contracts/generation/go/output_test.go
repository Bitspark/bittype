package generation_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	gen "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type assemblyNative struct {
	*nativeStore
	assemble func(gen.Assembly) gen.Assembly
}

func (n *assemblyNative) Assemble(ctx context.Context, pubs []gen.Publication[sourceLanguage, nativeLanguage], defs []gen.Definition[sourceLanguage, nativeLanguage], limits gen.Limits) (gen.Assembly, error) {
	out, err := n.nativeStore.Assemble(ctx, pubs, defs, limits)
	if err != nil {
		return gen.Assembly{}, err
	}
	return n.assemble(out), nil
}

func TestOutputRejectsUnassociatedAssembly(t *testing.T) {
	for _, tc := range []struct {
		name  string
		want  error
		alter func(gen.Assembly, gen.ExternalUse) gen.Assembly
	}{
		{"unrequested-exposure", gen.ErrIncomplete, func(out gen.Assembly, u gen.ExternalUse) gen.Assembly {
			out.Exposed = []ast.NodeRef{u.At.Site}
			return out
		}},
		{"foreign-use-site", gen.ErrScope, func(out gen.Assembly, u gen.ExternalUse) gen.Assembly {
			u.At.Site.Tree = "foreign"
			out.ExternalUses = []gen.ExternalUse{u}
			return out
		}},
		{"duplicate-destination", gen.ErrScope, func(out gen.Assembly, u gen.ExternalUse) gen.Assembly {
			out.ExternalUses = []gen.ExternalUse{u, u}
			return out
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			external, err := gen.NewSymbolIssuer().Issue(f.target, "owner/type")
			if err != nil {
				t.Fatal(err)
			}
			use := gen.ExternalUse{At: semantics.Occurrence{Site: ast.NodeRef{Tree: "native", Node: 2}, Role: semantics.Key{Owner: "test", Name: "type-use"}}, Symbol: external}
			c := f.configuration(nil)
			c.Native = &assemblyNative{f.native, func(out gen.Assembly) gen.Assembly { return tc.alter(out, use) }}
			run := begin(t, context.Background(), f, c, gen.DefaultLimits())
			if out, err := run.Finish(); !errors.Is(err, tc.want) || out.Snapshot().Valid() {
				t.Fatal("invalid assembly sealed", err)
			}
		})
	}
}

type checkingNative struct {
	fixtureIdentity
	inspect func()
}

func (c checkingNative) Check(ast.View[nativeLanguage]) error { c.inspect(); return nil }

func TestTargetSyntaxCheckerCannotReenterFinish(t *testing.T) {
	f := newFixture(t)
	var session *session
	var nested error
	calls := 0
	checker, err := check.NewChecker[nativeLanguage](f.target, checkingNative{fixtureIdentity("generation-target"), func() {
		calls++
		if calls == 1 {
			_, nested = session.Finish()
		}
	}})
	if err != nil {
		t.Fatal(err)
	}
	config := f.configuration(nil)
	config.Checker = checker
	session = begin(t, context.Background(), f, config, gen.DefaultLimits())
	failedFinish(t, session, gen.ErrInvalid)
	if calls != 1 || !errors.Is(nested, gen.ErrInvalid) {
		t.Fatal("checker escaped commit boundary", calls, nested)
	}
}
