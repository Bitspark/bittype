package golang

import (
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
)

// The contract identities are what a consumer registers the Go providers under
// and what a dependent names to require them. They are published here, with the
// profiles and keys they are built from, so that naming one imports no provider.

// ResolutionContract is the contract of the first Go profile's resolution.
func ResolutionContract() resolution.Contract {
	c := semantic.Contract{Profile: ResolutionProfile, SemanticsRequirements: Requirements()}
	for _, key := range []semantics.Key{Coverage(), Bindings(), Validity()} {
		c.Obligations = append(c.Obligations, semantic.Obligation{Key: key, Domain: Key("captured-package")})
	}
	return resolution.Contract{Semantic: c, Bindings: Bindings()}
}

// AnalysisContract is the contract of the first Go profile's analysis. It
// requires resolution under Configuration.
func AnalysisContract() analysis.Contract {
	c := ResolutionContract().Semantic
	c.Profile = AnalysisProfile
	return analysis.Contract{Semantic: c, Resolution: &resolution.Requirement{Profile: ResolutionProfile, Configuration: Configuration()}}
}
