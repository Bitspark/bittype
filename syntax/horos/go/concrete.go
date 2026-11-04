package horos

import "github.com/Bitspark/bittype/trees/cst/go"

// Concrete kinds are spelling/layout owned by Horos, not by the printer.
const (
	ConcreteLayout     cst.Kind[Language] = "layout"
	ConcreteToken      cst.Kind[Language] = "token"
	ConcreteWhitespace cst.Kind[Language] = "whitespace"
	ConcreteComment    cst.Kind[Language] = "comment"
)
