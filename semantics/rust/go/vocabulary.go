// Package rust reads immutable facts established by the pinned Rust compiler.
package rust

import (
	"errors"
	"reflect"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
)

const CompilerVersion = "1.90.0"
const ResolutionProfile = "rustc-1.90.0/captured-types-1/resolution-1"
const AnalysisProfile = "rustc-1.90.0/captured-types-1/analysis-1"

var ErrForeign = errors.New("foreign or missing Rust semantic handle")
var ErrCompiler = errors.New("Rust compiler helper failed")

func Key(name string) semantics.Key { return semantics.Key{Owner: "rust", Name: name} }
func Bindings() semantics.Key       { return Key("bindings") }
func Coverage() semantics.Key       { return Key("captured-type-coverage") }
func Validity() semantics.Key       { return Key("captured-type-validity") }

// The graph names below are written by the compiler helper, which is Rust
// and spells them itself, and read by name in model.go.

// Relations of facts: about a declaration and a native type.
func DeclaredTypeRelation() semantics.Key { return Key("declared-type") }
func PredicateRelation() semantics.Key    { return Key("predicate") }
func DefaultRelation() semantics.Key      { return Key("default") }
func NativeTypeRelation() semantics.Key   { return Key("native-type") }
func DisplayRelation() semantics.Key      { return Key("display") }
func FieldsRelation() semantics.Key       { return Key("fields") }

// The operator of the term that a fields fact carries.
func FieldsOperator() semantics.Key { return Key("fields") }

func Configuration() []semantic.Setting {
	return []semantic.Setting{
		{Key: Key("compiler"), Value: []byte(CompilerVersion)},
		{Key: Key("environment"), Value: []byte("edition=2024;target=x86_64-unknown-linux-gnu;crate=captured;lib;sysroot-only;no-codegen")},
		{Key: Key("helper"), Value: []byte("captured-types-1")},
	}
}
func MatchesConfiguration(settings []semantic.Setting) bool {
	return reflect.DeepEqual(settings, Configuration())
}
func Requirements() []semantic.Requirement {
	return []semantic.Requirement{{Obligation: Coverage(), Verdict: semantic.Holds}, {Obligation: Bindings(), Verdict: semantic.Holds}, {Obligation: Validity(), Verdict: semantic.Holds}}
}

// ResolutionContract is the contract of the captured Rust resolution profile.
// Resolver and consumers agree on this one value.
func ResolutionContract() resolution.Contract {
	c := semantic.Contract{Profile: ResolutionProfile, SemanticsRequirements: Requirements()}
	for _, key := range []semantics.Key{Coverage(), Bindings(), Validity()} {
		c.Obligations = append(c.Obligations, semantic.Obligation{Key: key, Domain: Key("captured-program")})
	}
	return resolution.Contract{Semantic: c, Bindings: Bindings()}
}

// AnalysisContract is the contract of the captured Rust analysis profile. It
// carries the resolution contract's obligations and requires a receipt of the
// resolution profile under the pinned configuration.
func AnalysisContract() analysis.Contract {
	c := ResolutionContract().Semantic
	c.Profile = AnalysisProfile
	return analysis.Contract{Semantic: c, Resolution: &resolution.Requirement{Profile: ResolutionProfile, Configuration: Configuration()}}
}

func Evidence() semantic.EvidenceRequirement {
	return semantic.EvidenceRequirement{Operation: analysis.Kind(), Profile: AnalysisProfile, Configuration: Configuration(), Reports: Requirements(), Semantics: true}
}
