package wronginput

import (
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type A struct{}
type B struct{}

func Mix(request semantic.Request[A]) ast.AST[B] { return request.Inputs()[0] }
