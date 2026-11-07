package verification

import "github.com/Bitspark/logos-contract/contract/go/typed"

// Constructors identify the source and target operators for this captured
// occurrence, in the order of the child correspondences supplied by its owner.
type Constructors struct{Source,Target string; Arity int}
type Evidence struct{Source,Target,Proof typed.Expr}

// Law supplies the type of one explicitly attributed assumption and constructs
// candidate syntax using evidence for this occurrence's children. The expected
// judgment, capture, dependency policy and Logos outcome are held separately.
type Law interface{
 Signature(Constructors)typed.Expr
 Apply(assumption typed.Expr,children []Evidence)typed.Expr
}

// Congruence provides the ordinary n-ary structural law used by mappings
// whose representations follow from one representation premise per child.
type Congruence struct{}
func(Congruence)Signature(c Constructors)typed.Expr{
 typ:=quantify(con("H"),c.Arity,func(h []expression)expression{
  return quantify(con("N"),c.Arity,func(n []expression)expression{
   result:=rep(app(con(c.Source),h...),app(con(c.Target),n...))
   for i:=c.Arity-1;i>=0;i--{result=arrow(rep(h[i],n[i]),result)}
   return result
  })
 })
 return typ(0)
}
func(Congruence)Apply(proof typed.Expr,children []Evidence)typed.Expr{
 for _,c:=range children{proof=typed.App(proof,c.Source)}
 for _,c:=range children{proof=typed.App(proof,c.Target)}
 for _,c:=range children{proof=typed.App(proof,c.Proof)}
 return proof
}
