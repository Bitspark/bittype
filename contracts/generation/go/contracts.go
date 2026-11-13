// Package generation defines declarative mapping and construction contracts.
package generation

import (
	"context"
	"errors"

	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

var (
	ErrInvalid     = errors.New("invalid generation description")
	ErrUnsupported = errors.New("generation capability unavailable")
	ErrScope       = errors.New("incompatible generation scope or association")
	ErrLimit       = errors.New("generation limit exhausted")
	ErrIncomplete  = errors.New("generation obligations incomplete")
	ErrAmbiguous   = errors.New("ambiguous generation mapping")
	// ErrRejected: a source judgment refuted a use from the captured evidence.
	ErrRejected = errors.New("source judgment rejected the use")
)

type Limits struct {
	Operations, Depth, Plans, References, Formals, Reservations, OutputNodes, OutputBytes int
}

func DefaultLimits() Limits { return Limits{100000, 128, 4096, 65536, 4096, 4096, 200000, 32 << 20} }
func (l Limits) Valid() bool {
	return l.Operations > 0 && l.Depth > 0 && l.Plans > 0 && l.References > 0 && l.Formals > 0 && l.Reservations > 0 && l.OutputNodes > 0 && l.OutputBytes > 0
}

// Target validates target intent during planning. Build receives only the closed
// target program; it has no source reader, selectors or mapping callbacks.
// Target implementations are trusted semantic compilers of their target.
// bittype:dimension Syntax
type Target[T any] interface {
	CheckExpression(context.Context, Expression[T]) error
	CheckDeclaration(context.Context, TargetDeclaration[T]) error
	CheckProgram(context.Context, Program[T]) error
	Build(context.Context, BuildInput[T]) (Assembly, error)
}

type BuildInput[T any] struct {
	Program  Program[T]
	Language representation.Language
	Output   ast.TreeID
	Limits   Limits
}
type Execution[T any] struct {
	Target   Target[T]
	Checker  check.Checker
	Language representation.Language
	Output   ast.TreeID
	Limits   Limits
}
