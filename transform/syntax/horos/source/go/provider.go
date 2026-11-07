package source

import (
	"context"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// Reader is the Horos interpretation API consumed by representation libraries.
// Neither the kernel nor native target implementations interpret these values.
type Reader interface {
	g.SourceSession[h.Language]
	Snapshot() semantic.Snapshot
	Declaration(g.Declaration[h.Language]) (Declaration, error)
	Parameter(g.Binder[h.Language]) (Parameter, error)
	IsBody(g.Environment[h.Language], g.Use[h.Language]) bool
	Read(g.Use[h.Language]) (Expression, error)
	Operand(g.Use[h.Language], ...uint32) (g.Use[h.Language], error)
	Fields(g.Use[h.Language]) ([]Field, error)
	Unsupported(g.Use[h.Language], string) error
}
type Provider struct{}

func (Provider) Contract() transform.Contract { return Contract() }
func (Provider) Open(ctx context.Context, r transform.Request[h.Language]) (Reader, error) {
	return New(ctx, r)
}

func (s *Source) ResolveRoots(ctx context.Context, sites []ast.NodeRef) ([]g.Declaration[h.Language], error) {
	if ctx == nil || len(sites) == 0 {
		return nil, g.ErrInvalid
	}
	bySite := map[ast.NodeRef]g.Declaration[h.Language]{}
	for _, d := range s.order {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		info, err := s.Declaration(d)
		if err != nil {
			return nil, err
		}
		bySite[info.Origin] = d
	}
	out := make([]g.Declaration[h.Language], 0, len(sites))
	seen := map[ast.NodeRef]bool{}
	for _, site := range sites {
		d, ok := bySite[site]
		if !ok || seen[site] {
			return nil, g.ErrScope
		}
		seen[site] = true
		out = append(out, d)
	}
	return out, nil
}

var _ g.SourceProvider[h.Language, Reader] = Provider{}
