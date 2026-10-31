// Package golang declares the Go type-syntax vocabulary over the generic AST.
// It contains no parser, builder, grammar checker, generator or formatter.
// The node shapes and this vocabulary's limits are documented in docs/syntax.md.
package golang

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/language/go"
)

// Language identifies this syntax family and is its distinct Go AST marker.
type Language struct{}

// ID is the stable registration key for this language family.
func (Language) ID() language.ID { return "go" }

// Name is the human-readable language name, not a registration key.
func (Language) Name() string { return "Go" }

var _ language.LanguageProvider = Language{}

// AST is the existing generic carrier specialized for this language.
type AST = ast.AST[Language]

// Node is a node in this language's unchecked syntax carrier.
type Node = ast.Node[Language]

// Kind is this language's syntax discriminator.
type Kind = ast.Kind[Language]

// Node kinds for file structure, type declarations and native type expressions.
const (
	KindFile              Kind = "file"
	KindImport            Kind = "import"
	KindTypeDeclaration   Kind = "type-declaration"
	KindTypeAlias         Kind = "type-alias"
	KindTypeParameter     Kind = "type-parameter"
	KindNamedType         Kind = "named-type"
	KindQualifiedType     Kind = "qualified-type"
	KindTypeApplication   Kind = "type-application"
	KindStructType        Kind = "struct-type"
	KindField             Kind = "field"
	KindEmbeddedField     Kind = "embedded-field"
	KindInterfaceType     Kind = "interface-type"
	KindMethod            Kind = "method"
	KindEmbeddedType      Kind = "embedded-type"
	KindFunctionType      Kind = "function-type" // Native signature syntax; no library Function semantics.
	KindParameter         Kind = "parameter"
	KindPointerType       Kind = "pointer-type"
	KindSliceType         Kind = "slice-type"
	KindArrayType         Kind = "array-type"
	KindMapType           Kind = "map-type"
	KindChannelType       Kind = "channel-type"
	KindTypeUnion         Kind = "type-union" // A constraint union, not a value union.
	KindTypeApproximation Kind = "type-approximation"
)

// Containment roles. Parameters and results are separate ordered lists.
const (
	RoleImport        = "import"
	RoleDeclaration   = "declaration"
	RoleType          = "type"
	RoleTypeParameter = "type-parameter"
	RoleConstraint    = "constraint"
	RoleConstructor   = "constructor"
	RoleArgument      = "argument"
	RoleField         = "field"
	RoleMember        = "member"
	RoleSignature     = "signature"
	RoleParameter     = "parameter"
	RoleResult        = "result"
	RoleElement       = "element"
	RoleKey           = "key"
	RoleTerm          = "term"
)

// Reference roles. Target links name declarations, imports or lexical binders.
const (
	RoleTarget = "target"
)

// Attribute names. Package qualifiers are identifiers, not dot-joined type source.
const (
	AttributePackage   = "package"
	AttributeName      = "name"
	AttributePath      = "path"
	AttributeAlias     = "alias"
	AttributeQualifier = "qualifier"
	AttributeTag       = "tag"
	AttributeLength    = "length"
	AttributeDirection = "direction"
	AttributeVariadic  = "variadic"
)
