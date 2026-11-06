// Package rust reads immutable facts established by the pinned Rust compiler.
package rust

import (
	"errors"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	"reflect"
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
func Evidence() semantic.EvidenceRequirement {
	return semantic.EvidenceRequirement{Operation: analysis.Kind(), Profile: AnalysisProfile, Configuration: Configuration(), Reports: Requirements(), Semantics: true}
}
