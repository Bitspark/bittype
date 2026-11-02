package wronginput

import (
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type A struct{}
type B struct{}

func Mix(request analysis.Request[A]) ast.AST[B] { return request.Inputs()[0] }
