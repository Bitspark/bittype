package analyze

import (
	"context"
	"reflect"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/internal/typescriptsem"
	resolve "github.com/Bitspark/bittype/resolve/files/typescript/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Queries struct { model n.Model; compiler resolve.Compiler }
type Assignability struct { model n.Model; source n.Type; target n.Type; holds bool }
// NewQueries configures explicit compiler queries; it starts no hidden process.
func NewQueries(model n.Model,compiler resolve.Compiler) (Queries,error) {
	if !model.Valid() { return Queries{},n.ErrForeign }; if err:=compiler.Validate(); err!=nil { return Queries{},err }
	if _,err:=n.Open(model.Receipt()); err!=nil { return Queries{},err }; return Queries{model,compiler},nil
}
func (q Queries) Assignable(ctx context.Context,source,target n.Type) (Assignability,error) {
	if ctx==nil || !q.model.Valid() || !source.Valid() || !target.Valid() || !q.model.Same(source.Model()) || !q.model.Same(target.Model()) { return Assignability{},n.ErrForeign }
	if err:=ctx.Err(); err!=nil { return Assignability{},err }
	var trees []s.AST
	for _,input:=range q.model.Snapshot().Inputs() { tree,err:=representation.Unpack[s.Language](q.model.Snapshot().Language(),input); if err!=nil { return Assignability{},err }; trees=append(trees,tree) }
	out,err:=typescriptsem.Check(ctx,q.compiler.NodePath,q.compiler.TypeScriptPath,trees,q.model.Snapshot().Files(),&typescriptsem.Query{Source:uint32(source.ID()),Target:uint32(target.ID())})
	if err!=nil { return Assignability{},err }
	if out.Status!="ok" || out.Facts==nil || !reflect.DeepEqual(*out.Facts,q.model.Graph().Data()) || out.Assignable==nil { return Assignability{},semantic.ErrEvidenceUnavailable }
	return Assignability{q.model,source,target,*out.Assignable},nil
}
func (a Assignability) Model() n.Model { return a.model }
func (a Assignability) Source() n.Type { return a.source }
func (a Assignability) Target() n.Type { return a.target }
func (a Assignability) Holds() bool { return a.model.Valid() && a.holds }
