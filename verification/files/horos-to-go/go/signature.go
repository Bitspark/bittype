package verification

import (
	"slices"

	association "github.com/Bitspark/bittype/contracts/verification/go"
	joint "github.com/Bitspark/bittype/verification/logos/go"
	"github.com/Bitspark/logos-contract/contract/go/typed"
)

// signature records the exact sites supplied at declaration construction.
type signature struct {
	declarations []typed.Declaration
	origins      []joint.Origin
}

func (s *signature) add(d typed.Declaration, sites ...association.Site) {
	if len(sites) > 0 {
		s.origins = append(s.origins, joint.Origin{Location: typed.Location{Phase: typed.DeclarationPhase, Index: uint32(len(s.declarations)), HasIndex: true}, Sites: slices.Clone(sites)})
	}
	s.declarations = append(s.declarations, d)
}
func (s *signature) include(part signature) {
	offset := uint32(len(s.declarations))
	s.declarations = append(s.declarations, part.declarations...)
	for _, o := range part.origins {
		o.Location.Index += offset
		s.origins = append(s.origins, o)
	}
}
func (s signature) problem(input association.Input, req []association.Requirement, expected typed.Expr, sites []association.Site) joint.Problem {
	origins := append([]joint.Origin{{Location: typed.Location{Phase: typed.ExpectedPhase}, Sites: sites}}, s.origins...)
	return joint.Problem{Input: input, Requirements: req, Statement: typed.Statement{Foundation: typed.Foundation, Declarations: s.declarations, Expected: expected}, Origins: origins}
}
