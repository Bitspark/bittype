// Package golang reads the immutable facts published by the Go semantic providers.
package golang

import (
	"errors"
	"runtime"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
)

// The profile has its second version since 13 November 2026 (bittype#310). A
// published member changed: the body of an entity is the term of its
// declaration's written definition, and the type of a type name is the
// declared-type fact of its entity. Until then the body was that type. A
// receipt of captured-types-1 is no receipt of this profile.
const ResolutionProfile = "go1.25/captured-types-2/resolution-1"
const AnalysisProfile = "go1.25/captured-types-2/analysis-1"
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

// BuiltinKind is the kind of the entity of a predeclared basic type, such as
// string and bool. The entity's name is the identifier of the universe scope.
// It has no owner and no declaration site.
//
// Limit: the kind covers the basic types. error and comparable are named types
// of the universe scope and any is an alias, so their entities have
// DefinedKind and AliasKind, without an owner and without a site.
func BuiltinKind() semantics.Key { return Key("builtin") }

// NameRole is the role of the binding of an identifier that is used as a
// name: the binding's site is the node of the use, and its target is the
// entity that the compiler resolved the identifier to.
func NameRole() semantics.Key { return Key("name") }

// DeclaredTypeRelation is the relation of the fact from the entity of a type
// name to the term of its own type: a defined type, an alias, a type parameter
// and a predeclared type, captured, supplied or of the universe scope. Every
// such entity carries exactly one. The term is the one that the facts about
// the type hang on. It is not the entity's body: the body of a captured
// declaration is the term of its written definition, and an entity without a
// captured definition has no body.
func DeclaredTypeRelation() semantics.Key { return Key("declared-type") }

// Relations of the facts about a native type.
func IdenticalRelation() semantics.Key   { return Key("identical") }
func UnderlyingRelation() semantics.Key  { return Key("underlying") }
func ConstraintRelation() semantics.Key  { return Key("constraint") }
func AliasTargetRelation() semantics.Key { return Key("alias-target") }
func EmptyRelation() semantics.Key       { return Key("empty") }
func DisplayRelation() semantics.Key     { return Key("display") }

// Relation of the fact about a package entity: the name that its package
// clauses declare. The entity's own name is the import path. The fact's term
// is a provider expression with PackageNameOperator and TextEncoding and no
// operand, whose data is the name. The name readers refuse another form.
func PackageNameRelation() semantics.Key { return Key("package-name") }
func PackageNameOperator() semantics.Key { return Key("package-name") }

// TextEncoding is the encoding of a term whose data is UTF-8 text.
func TextEncoding() semantics.Key { return Key("utf8") }

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
