package typescript

import (
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
)

// ResolutionContract is the identity of the captured TypeScript resolution
// profile. A resolver is registered under it and a consumer names it; neither
// needs the other's package.
func ResolutionContract() resolution.Contract {
	c := semantic.Contract{Profile: ResolutionProfile, SemanticsRequirements: Requirements()}
	for _, key := range []semantics.Key{Coverage(), Bindings(), Validity()} {
		c.Obligations = append(c.Obligations, semantic.Obligation{Key: key, Domain: Key("captured-program")})
	}
	return resolution.Contract{Semantic: c, Bindings: Bindings()}
}

// AnalysisContract is the identity of the captured TypeScript analysis profile.
// It requires resolution under the original captured program configuration.
func AnalysisContract() analysis.Contract {
	c := ResolutionContract().Semantic
	c.Profile = AnalysisProfile
	return analysis.Contract{Semantic: c, Resolution: &resolution.Requirement{Profile: ResolutionProfile, Configuration: Configuration()}}
}

// AnalysisContractFor binds analysis to resolution with exactly the supplied package environment.
func AnalysisContractFor(settings []semantic.Setting) (analysis.Contract, error) {
	evidence, err := EvidenceFor(settings)
	if err != nil {
		return analysis.Contract{}, err
	}
	c := AnalysisContract()
	c.Resolution.Configuration = evidence.Configuration
	return c, nil
}
