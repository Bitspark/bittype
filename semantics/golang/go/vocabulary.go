// Package golang reads the immutable facts published by the Go semantic providers.
package golang

import (
	"errors"
	"reflect"
	"runtime"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
)

const ResolutionProfile = "go1.25/captured-types-1/resolution-1"
const AnalysisProfile = "go1.25/captured-types-1/analysis-1"
const LanguageVersion = "go1.25"
const PackagePath = "captured"

var ErrForeign = errors.New("foreign or missing Go semantic handle")
var ErrUnsupported = errors.New("Go semantics outside supported fragment")

// Key returns a Go-owned vocabulary key; it performs no registration.
func Key(name string) semantics.Key { return semantics.Key{Owner: "go", Name: name} }
func Bindings() semantics.Key { return Key("bindings") }
func Coverage() semantics.Key { return Key("captured-type-coverage") }
func Validity() semantics.Key { return Key("captured-type-validity") }

// Configuration fixes the compiler environment of this bounded provider.
func Configuration() []semantic.Setting {
	return []semantic.Setting{
		{Key: Key("language-version"), Value: []byte(LanguageVersion)},
		{Key: Key("toolchain"), Value: []byte(runtime.Version())},
		{Key: Key("package-path"), Value: []byte(PackagePath)},
		{Key: Key("sizes"), Value: []byte("gc/amd64")},
	}
}
func MatchesConfiguration(settings []semantic.Setting) bool {
	return reflect.DeepEqual(settings, Configuration())
}
func Requirements() []semantic.Requirement {
	return []semantic.Requirement{
		{Obligation: Coverage(), Verdict: semantic.Holds},
		{Obligation: Bindings(), Verdict: semantic.Holds},
		{Obligation: Validity(), Verdict: semantic.Holds},
	}
}
func Evidence() semantic.EvidenceRequirement {
	return semantic.EvidenceRequirement{Operation: analysis.Kind(), Profile: AnalysisProfile,
		Configuration: Configuration(), Reports: Requirements(), Semantics: true}
}
