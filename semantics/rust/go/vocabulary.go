// Package rust reads immutable facts established by the pinned Rust compiler.
package rust

import (
	"errors"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
)

const CompilerVersion = "1.90.0"

// The two profiles are in their second version. The first one's fields term
// carried a member that this one's does not, so a reader of the first cannot
// read a graph of this one.
const ResolutionProfile = "rustc-1.90.0/captured-types-2/resolution-1"
const AnalysisProfile = "rustc-1.90.0/captured-types-2/analysis-1"

// HelperProfile names the facts that the compiler helper writes. It is a
// setting of every configuration and changes when facts are added, so a
// receipt of another helper is never read as one of this helper. The two
// profile strings above are not derived from it: they change on their own
// occasion. An added entity kind and binding changed this name and kept
// those, because nothing was removed. Lookup changed with the addition: it
// never selects a primitive entity. A reader of before is kept from a graph
// with these bindings by this setting, whose new value it refuses.
//
// The helper states the same name in every response, and the bridge refuses a
// response that states another one. So the name is declared here and spelled
// once more in the helper, and a helper binary of another change of the facts
// is never taken for this one.
const HelperProfile = "captured-types-3"

var ErrForeign = errors.New("foreign or missing Rust semantic handle")
var ErrCompiler = errors.New("Rust compiler helper failed")

// ErrUnsupported is unsupported coverage: the input is well formed under the
// vocabulary of this profile, and the profile does not cover it. A later
// profile may admit the same input unchanged. It never means invalid input.
var ErrUnsupported = errors.New("Rust semantics outside supported fragment")

// ErrMissingFact reports that a graph lacks a fact its reader requires, or
// carries it in another form. An absent visibility fact never reads as
// private.
var ErrMissingFact = errors.New("Rust semantic fact is absent or malformed")

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

// Relations of the visibility facts about a captured declaration.
//
// The declared fact is the visibility as written, in the compiler's resolved
// form. Its term has PublicOperator, or RestrictedOperator with one operand
// that refers to the module entity the declaration is restricted to. Every
// captured module below the crate root, alias, struct, enum, trait and field
// has exactly one.
//
// The effective fact says whether an audience can name the declaration through
// a path. Its term has EffectiveVisibilityOperator and names both. Every
// captured module below the crate root, alias, struct, enum and trait has
// exactly one.
//
// A private-interface fact is the compiler's observation that the interface of
// the declaration names a declaration more private than itself. Its term
// refers to that more private declaration. A declaration has one for each
// observation and none without one.
func DeclaredVisibilityRelation() semantics.Key  { return Key("declared-visibility") }
func EffectiveVisibilityRelation() semantics.Key { return Key("effective-visibility") }
func PrivateInterfaceRelation() semantics.Key    { return Key("private-interface") }

// The operators of the terms that the two visibility facts carry.
func PublicOperator() semantics.Key              { return Key("public") }
func RestrictedOperator() semantics.Key          { return Key("restricted") }
func EffectiveVisibilityOperator() semantics.Key { return Key("effective-visibility") }

// The kind of a module entity. The scope of a restricted visibility is one,
// and the crate root is one.
func ModuleKind() semantics.Key { return Key("module") }

// PrimitiveKind is the kind of the entity of a primitive type, such as bool
// and u8. The entity's name is the name that the language gives the type. A
// primitive type is defined in no crate, so the entity has no owner, no
// declaration site and no compiler path, and it carries no visibility fact.
// A graph holds one such entity for each primitive type that its source uses.
//
// The helper binds every primitive type node, under the role of a path
// segment, to what the compiler resolved the name to: this entity, or the
// declaration or the generic parameter of that name that captures the use.
func PrimitiveKind() semantics.Key { return Key("primitive") }

// The audience and the path of the effective fact of this profile: another
// crate, through the path of the definition. The profile has no re-export, so
// this is the one public path a declaration can have.
const (
	AudienceAnotherCrate = "another-crate"
	PathDefinition       = "definition"
)

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

// AnalysisContractFor is the contract of the captured Rust analysis profile
// under exactly the supplied configuration. It carries the resolution
// contract's obligations and requires a receipt of the resolution profile
// under that configuration.
func AnalysisContractFor(settings []semantic.Setting) (analysis.Contract, error) {
	evidence, err := EvidenceFor(settings)
	if err != nil {
		return analysis.Contract{}, err
	}
	c := ResolutionContract().Semantic
	c.Profile = AnalysisProfile
	return analysis.Contract{Semantic: c, Resolution: &resolution.Requirement{Profile: ResolutionProfile, Configuration: evidence.Configuration}}, nil
}
