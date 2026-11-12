// Package horos names the semantic identities of the horos frontend: its
// profiles, keys, prelude pin and contracts. It resolves and analyzes nothing.
package horos

import (
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
)

const ResolutionProfile = "horos-declarations-4/resolution-1"
const AnalysisProfile = "horos-declarations-4/analysis-1"

// Key returns a horos-owned vocabulary key; it performs no registration.
func Key(name string) semantics.Key       { return semantics.Key{Owner: "horos", Name: name} }
func ReferenceResolution() semantics.Key  { return Key("reference-resolution") }
func Validity() semantics.Key             { return Key("program-validity") }
func Indexing() semantics.Key             { return Key("declaration-index") }
func Bindings() semantics.Key             { return Key("occurrence-bindings") }
func Applications() semantics.Key         { return Key("applications") }
func Dependencies() semantics.Key         { return Key("reference-dependencies") }
func TypeFacts() semantics.Key            { return Key("type-facts") }
func TargetRole() semantics.Key           { return Key("target") }
func ExpressionRole() semantics.Key       { return Key("expression") }
func ReferenceDependency() semantics.Key  { return Key("declaration-reference") }
func SortRelation() semantics.Key         { return Key("sort") }
func ResultSortRelation() semantics.Key   { return Key("result-sort") }
func DeclaredKindRelation() semantics.Key { return Key("declared-kind") }
func ResultKindRelation() semantics.Key   { return Key("result-kind") }
func SymbolicBodyRelation() semantics.Key { return Key("symbolic-body") }
func MembershipRelation() semantics.Key   { return Key("membership") }
func IdentityRelation() semantics.Key     { return Key("canonical-identity") }

// Kinds of graph entities; ModuleKind is also the kind of a Files group.
func ModuleKind() semantics.Key         { return Key("module") }
func TypeDefinitionKind() semantics.Key { return Key("type-definition") }

// The kind of the Files relation from an importing file to the imported one.
func ImportRelation() semantics.Key { return Key("import") }

// Operators of the terms that sort and identity facts carry.
func ProperSort() semantics.Key     { return Key("sort/proper") }
func ModuleIdentity() semantics.Key { return Key("identity/module") }

// LimitsSetting contains a JSON array [depth,nodes,bytes,modules]. All four
// positive limits are required; depth must not exceed 512. Its absence selects
// Horos's fixed defaults.
// This bounds authoritative resolution separately from shared graph publication.
func LimitsSetting() semantics.Key { return Key("declaration-limits-json-1") }

// PreludePin identifies the revision's bundled module, not a source location.
func PreludePin() string { return decl.PreludePin().String() }

func ResolutionContract() resolution.Contract {
	keys := []semantics.Key{ReferenceResolution(), Validity(), Indexing(), Bindings(), Applications(), Dependencies(), TypeFacts()}
	contract := semantic.Contract{Profile: ResolutionProfile}
	for _, k := range keys {
		contract.Obligations = append(contract.Obligations, semantic.Obligation{Key: k, Domain: Key("entry-transitive-module-closure")})
		verdict := semantic.Undetermined
		if k == Validity() || k == ReferenceResolution() {
			verdict = semantic.Holds
		}
		contract.SemanticsRequirements = append(contract.SemanticsRequirements, semantic.Requirement{Obligation: k, Verdict: verdict})
	}
	return resolution.Contract{Semantic: contract, Bindings: ReferenceResolution()}
}

func AnalysisContract() analysis.Contract { return AnalysisContractWithResolutionConfiguration(nil) }

func AnalysisContractWithResolutionConfiguration(settings []semantic.Setting) analysis.Contract {
	resolved := ResolutionContract()
	requirement := resolution.Requirement{Profile: resolved.Semantic.Profile, Configuration: settings, Reports: resolved.Semantic.SemanticsRequirements}.Clone()
	contract := semantic.Contract{Profile: AnalysisProfile}
	for _, obligation := range resolved.Semantic.Obligations {
		if obligation.Key != resolved.Bindings {
			contract.Obligations = append(contract.Obligations, obligation)
		}
	}
	for _, required := range resolved.Semantic.SemanticsRequirements {
		if required.Obligation != resolved.Bindings {
			contract.SemanticsRequirements = append(contract.SemanticsRequirements, required)
		}
	}
	return analysis.Contract{Semantic: contract, Resolution: &requirement}
}
