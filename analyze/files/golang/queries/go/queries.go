// Package queries answers explicit native Go questions about an admitted model.
package queries

import (
	"context"
	"go/types"
	"reflect"
	"sync"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/internal/golangsem"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type queryState struct {
	model   n.Model
	checked *golangsem.Checked
	mu      sync.Mutex
}
type Queries struct{ state *queryState }
type Admission struct {
	model       n.Model
	constructor n.Type
	arguments   []n.Type
	admitted    bool
	detail      string
}

// NewQueries explicitly reconstructs the compiler state for this admitted model.
// Ordinary analysis and graph readers do not invoke this operation implicitly.
func NewQueries(ctx context.Context, model n.Model) (Queries, error) {
	if ctx == nil || !model.Valid() {
		return Queries{}, n.ErrForeign
	}
	if err := ctx.Err(); err != nil {
		return Queries{}, err
	}
	if _, err := n.Open(model.Receipt()); err != nil {
		return Queries{}, err
	}
	var inputs []s.AST
	for _, input := range model.Snapshot().Inputs() {
		tree, err := representation.Unpack[s.Language](model.Snapshot().Language(), input)
		if err != nil {
			return Queries{}, err
		}
		inputs = append(inputs, tree)
	}
	checked, err := golangsem.Check(ctx, inputs)
	if err != nil {
		return Queries{}, err
	}
	if !reflect.DeepEqual(checked.Facts, model.Graph().Data()) {
		return Queries{}, semantic.ErrInvalidResult
	}
	return Queries{&queryState{model: model, checked: checked}}, nil
}
func (q Queries) Admit(ctx context.Context, constructor n.Type, arguments ...n.Type) (Admission, error) {
	if ctx == nil || q.state == nil || !constructor.Valid() || !q.state.model.Same(constructor.Model()) {
		return Admission{}, n.ErrForeign
	}
	if err := ctx.Err(); err != nil {
		return Admission{}, err
	}
	q.state.mu.Lock()
	defer q.state.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return Admission{}, err
	}
	t, ok := q.state.checked.Types[constructor.ID()].(*types.Named)
	if !ok || t.TypeParams().Len() != len(arguments) || len(arguments) > 16 {
		return Admission{}, n.ErrUnsupported
	}
	args := make([]types.Type, len(arguments))
	for i, a := range arguments {
		if !a.Valid() || !q.state.model.Same(a.Model()) {
			return Admission{}, n.ErrForeign
		}
		args[i] = q.state.checked.Types[a.ID()]
		if args[i] == nil {
			return Admission{}, n.ErrForeign
		}
	}
	_, err := types.Instantiate(nil, t, args, true)
	if e := ctx.Err(); e != nil {
		return Admission{}, e
	}
	result := Admission{model: q.state.model, constructor: constructor, arguments: append([]n.Type(nil), arguments...), admitted: err == nil}
	if err != nil {
		result.detail = err.Error()
	}
	return result, nil
}
func (a Admission) Model() n.Model      { return a.model }
func (a Admission) Constructor() n.Type { return a.constructor }
func (a Admission) Arguments() []n.Type { return append([]n.Type(nil), a.arguments...) }
func (a Admission) Admitted() bool      { return a.model.Valid() && a.admitted }
func (a Admission) Detail() string      { return a.detail }
