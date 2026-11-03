// Package horosmodel shares tooling mechanics between explicit Horos profiles.
// It is not part of the agnostic kernel. It has no profile registry, constructor
// dispatch, default frontend or behavioral contract knowledge.
package horosmodel

import (
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
)

type Program interface {
	Root() decl.Digest
	Modules() []*decl.Module
	Module(decl.Digest) (*decl.Module, bool)
	Declaration(decl.DeclarationID) (*decl.Declaration, bool)
}
type Frontend[L any] struct {
	DeclarationPrefix, DescriptorPrefix string
	Owner, ResolutionProfile            string
	ResolveModel                        func(core.Value, []core.Value, decl.Limits) (Program, error)
	PinModel                            func(core.Value, decl.Limits) (decl.Digest, error)
}
