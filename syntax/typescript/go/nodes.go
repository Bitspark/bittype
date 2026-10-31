// Package typescript declares the TypeScript type-syntax vocabulary over the generic AST.
// It contains no parser, builder, grammar checker, generator or formatter.
// The node shapes and this vocabulary's limits are documented in docs/syntax.md.
package typescript

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/language/go"
)

// Language identifies this syntax family and is its distinct Go AST marker.
type Language struct{}

// ID is the stable registration key for this language family.
func (Language) ID() language.ID { return "typescript" }

// Name is the human-readable language name, not a registration key.
func (Language) Name() string { return "TypeScript" }

var _ language.Language = Language{}

// AST is the existing generic carrier specialized for this language.
type AST = ast.AST[Language]

// Node is a node in this language's unchecked syntax carrier.
type Node = ast.Node[Language]

// Kind is this language's syntax discriminator.
type Kind = ast.Kind[Language]

// Node kinds for declarations, members and type expressions.
const (
	KindSourceFile Kind = "source-file"
	KindImport Kind = "import"
	KindImportBinding Kind = "import-binding"
	KindNamespaceImport Kind = "namespace-import"
	KindTypeAlias Kind = "type-alias"
	KindInterface Kind = "interface"
	KindTypeParameter Kind = "type-parameter"
	KindTypeReference Kind = "type-reference"
	KindQualifiedType Kind = "qualified-type"
	KindTypeApplication Kind = "type-application"
	KindKeywordType Kind = "keyword-type"
	KindLiteralType Kind = "literal-type"
	KindArrayType Kind = "array-type"
	KindTupleType Kind = "tuple-type"
	KindTupleElement Kind = "tuple-element"
	KindObjectType Kind = "object-type"
	KindUnionType Kind = "union-type"
	KindIntersectionType Kind = "intersection-type"
	KindFunctionType Kind = "function-type" // Native signature syntax; no library Function semantics.
	KindParameter Kind = "parameter"
	KindProperty Kind = "property"
	KindMethod Kind = "method"
	KindCallSignature Kind = "call-signature"
	KindIndexSignature Kind = "index-signature"
	KindTypeOperator Kind = "type-operator"
	KindIndexedAccessType Kind = "indexed-access-type"
)

// Containment roles. Their allowed kinds and cardinality are language-owned.
const (
	RoleImport = "import"
	RoleDeclaration = "declaration"
	RoleBinding = "binding"
	RoleType = "type"
	RoleTypeParameter = "type-parameter"
	RoleConstraint = "constraint"
	RoleDefault = "default"
	RoleBase = "base"
	RoleMember = "member"
	RoleQualifier = "qualifier"
	RoleConstructor = "constructor"
	RoleArgument = "argument"
	RoleElement = "element"
	RoleParameter = "parameter"
	RoleResult = "result"
	RoleSignature = "signature"
	RoleObject = "object"
	RoleIndex = "index"
)

// Reference roles. Target links name declarations, imports or lexical binders.
const (
	RoleTarget = "target"
)

// Attribute names. Text is an identifier, literal token or documented atom, never raw source.
const (
	AttributeName = "name"
	AttributeModule = "module"
	AttributeAlias = "alias"
	AttributeExported = "exported"
	AttributeOptional = "optional"
	AttributeReadonly = "readonly"
	AttributeRest = "rest"
	AttributeLiteralKind = "literal-kind"
	AttributeText = "text"
	AttributeOperator = "operator"
)
