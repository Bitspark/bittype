package independent
import (
 "github.com/Bitspark/bittype/ast/go"
 "github.com/Bitspark/bittype/transform/go"
)
type First struct{}
type Second struct{}
type Third struct{}
func Wrong(from,to transform.Language){
 _,_=transform.NewTransformer[First,Second](from,to,func(ast.AST[First])(ast.AST[Third],error){return ast.AST[Third]{},nil})
}
