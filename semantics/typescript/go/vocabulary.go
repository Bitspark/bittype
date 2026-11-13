// Package typescript reads immutable TypeScript compiler facts.
package typescript

import (
	"errors"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
)

const CompilerVersion = "5.9.3"

// The profile has its second version since 13 November 2026 (bittype#310). A
// published member changed: the body of an entity is the term of its
// declaration's written definition. Until then it was the declared type of
// the name, which the declared-type fact carries then and now. A receipt of
// captured-types-1 is no receipt of this profile.
const ResolutionProfile = Profile + "/resolution-1"
const AnalysisProfile = Profile + "/analysis-1"

// Profile is the name of the profile without its operation: the compiler
// version and the helper profile. The two profile strings are built from it.
const Profile = "typescript-" + CompilerVersion + "/" + HelperProfile

// HelperProfile names what the compiler helper writes. It is the value of the
// helper setting, so it is part of every receipt, and it changed with the
// profile: settings with captured-types-1 are no configuration of this one.
// The helper is JavaScript and spells the name once itself, for its messages.
const HelperProfile = "captured-types-2"

var ErrForeign = errors.New("foreign or missing TypeScript semantic handle")
var ErrCompiler = errors.New("TypeScript compiler helper failed")

func Key(name string) semantics.Key { return semantics.Key{Owner: "typescript", Name: name} }
func Bindings() semantics.Key       { return Key("bindings") }
func Coverage() semantics.Key       { return Key("captured-type-coverage") }
func Validity() semantics.Key       { return Key("captured-type-validity") }

// The graph names below are written by the compiler helper, which is
// JavaScript and spells them itself, and read by name in Go.

// Kinds of graph entities; ModuleKind is also the kind of a Files group.
func ModuleKind() semantics.Key      { return Key("module") }
func PackageKind() semantics.Key     { return Key("package") }
func PackageTypeKind() semantics.Key { return Key("package-type") }
func BuiltinKind() semantics.Key     { return Key("builtin") }
func AliasKind() semantics.Key       { return Key("alias") }
func InterfaceKind() semantics.Key   { return Key("interface") }
func ParameterKind() semantics.Key   { return Key("parameter") }

// Relations of facts: about a module, a declaration and a native type.
func NamedExportRelation() semantics.Key { return Key("named-export") }

// DeclaredTypeRelation relates the entity of a captured alias, interface,
// type parameter and property to the term of its type. A supplied package
// type and its type parameters carry no such fact. The term is not the
// entity's body: the body of an alias and of an interface is the term of its
// written definition, and a type parameter and a property have none.
// Model.DeclaredType reads the fact.
func DeclaredTypeRelation() semantics.Key { return Key("declared-type") }
func ConstraintRelation() semantics.Key   { return Key("constraint") }
func NativeTypeRelation() semantics.Key   { return Key("native-type") }
func DisplayRelation() semantics.Key      { return Key("display") }
func PropertiesRelation() semantics.Key   { return Key("properties") }

// The operator of the term that a properties fact carries. It is also the
// operator of the body of an interface without an extends clause and of the
// last operand of an extension.
func ObjectOperator() semantics.Key { return Key("object") }

// ExtensionOperator is the operator of the body of an interface that has an
// extends clause. The term's encoding is JSONEncoding and its data an
// ExtensionInfo. Its operands are the base type terms as written, in written
// order, and then one object term of the members that the interface declares
// itself. Model.Extension reads it.
func ExtensionOperator() semantics.Key { return Key("extension") }

// JSONEncoding is the encoding of a term whose data is JSON.
func JSONEncoding() semantics.Key { return Key("json-1") }

// Roles of binding occurrences: a type reference and an import binding.
func NameRole() semantics.Key   { return Key("name") }
func ImportRole() semantics.Key { return Key("import") }

func Configuration() []semantic.Setting {
	return []semantic.Setting{
		{Key: Key("compiler"), Value: []byte(CompilerVersion)},
		{Key: Key("environment"), Value: []byte("ES2020;ESNext;strict;exactOptionalPropertyTypes;noUncheckedIndexedAccess;lib.es2020;no-ambient-packages")},
		{Key: Key("helper"), Value: []byte(HelperProfile)},
	}
}
func MatchesConfiguration(settings []semantic.Setting) bool {
	_, _, err := Imports(settings)
	return err == nil
}
func Requirements() []semantic.Requirement {
	return []semantic.Requirement{{Obligation: Coverage(), Verdict: semantic.Holds}, {Obligation: Bindings(), Verdict: semantic.Holds}, {Obligation: Validity(), Verdict: semantic.Holds}}
}
func Evidence() semantic.EvidenceRequirement {
	return semantic.EvidenceRequirement{Operation: analysis.Kind(), Profile: AnalysisProfile, Configuration: Configuration(), Reports: Requirements(), Semantics: true}
}
