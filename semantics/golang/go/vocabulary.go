// Package golang reads the immutable facts published by the Go semantic providers.
package golang

import (
	"errors"
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
func Bindings() semantics.Key       { return Key("bindings") }
func Coverage() semantics.Key       { return Key("captured-type-coverage") }
func Validity() semantics.Key       { return Key("captured-type-validity") }

// The graph names below are written by the compiler bridge internal/golangsem
// and read by name in model.go.

// Kinds of graph entities.
func PackageKind() semantics.Key   { return Key("package") }
func DefinedKind() semantics.Key   { return Key("defined") }
func AliasKind() semantics.Key     { return Key("alias") }
func ParameterKind() semantics.Key { return Key("parameter") }

// Relations of the facts about a native type.
func IdenticalRelation() semantics.Key   { return Key("identical") }
func UnderlyingRelation() semantics.Key  { return Key("underlying") }
func ConstraintRelation() semantics.Key  { return Key("constraint") }
func AliasTargetRelation() semantics.Key { return Key("alias-target") }
func EmptyRelation() semantics.Key       { return Key("empty") }
func DisplayRelation() semantics.Key     { return Key("display") }

// Operators of the terms of a basic, a slice and a struct type.
func BasicOperator() semantics.Key  { return Key("basic") }
func SliceOperator() semantics.Key  { return Key("slice") }
func StructOperator() semantics.Key { return Key("struct") }

// Configuration fixes the compiler environment of this bounded provider.
func Configuration() []semantic.Setting {
	return []semantic.Setting{
		{Key: Key("language-version"), Value: []byte(LanguageVersion)},
		{Key: Key("toolchain"), Value: []byte(runtime.Version())},
		{Key: Key("package-path"), Value: []byte(PackagePath)},
		{Key: Key("sizes"), Value: []byte("gc/amd64")},
	}
}

// MatchesConfiguration accepts the base Configuration or the base extended by
// exactly one normalized ImportConfiguration.
func MatchesConfiguration(settings []semantic.Setting) bool {
	_, _, err := Imports(settings)
	return err == nil
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
