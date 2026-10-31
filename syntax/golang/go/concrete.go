package golang

import "github.com/Bitspark/bittype/cst/go"

// Concrete kinds are owned by this language. The shared printer never
// interprets their spelling. Grammar containers use the native rule names.
const (
	ConcreteLayout cst.Kind[Language] = "layout"
	ConcreteToken cst.Kind[Language] = "token"
	ConcreteWhitespace cst.Kind[Language] = "whitespace"
	ConcreteComment cst.Kind[Language] = "comment"
)

