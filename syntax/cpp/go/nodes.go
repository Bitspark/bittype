// Package cpp declares the C++ type-syntax vocabulary over the generic AST.
// It contains no parser, builder, grammar checker, generator or formatter.
// The node shapes and this vocabulary's limits are documented in docs/syntax.md.
package cpp

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/language/go"
)

// Language identifies this syntax family and is its distinct Go AST marker.
type Language struct{}

// ID is the stable registration key for this language family.
func (Language) ID() language.ID { return "cpp" }

// Name is the human-readable language name, not a registration key.
func (Language) Name() string { return "C++" }

var _ language.Language = Language{}

// AST is the existing generic carrier specialized for this language.
type AST = ast.AST[Language]

// Node is a node in this language's unchecked syntax carrier.
type Node = ast.Node[Language]

// Kind is this language's syntax discriminator.
type Kind = ast.Kind[Language]

// Node kinds for declarations, templates, named types and declarator type structure.
const (
	KindTranslationUnit Kind = "translation-unit"
	KindInclude Kind = "include"
	KindNamespace Kind = "namespace"
	KindRecord Kind = "record"
	KindTypeAlias Kind = "type-alias"
	KindEnum Kind = "enum"
	KindEnumerator Kind = "enumerator"
	KindBase Kind = "base"
	KindField Kind = "field"
	KindMethod Kind = "method"
	KindDestructor Kind = "destructor"
	KindParameter Kind = "parameter"
	KindTypeParameter Kind = "type-parameter"
	KindNamedType Kind = "named-type"
	KindNameSegment Kind = "name-segment"
	KindBuiltinType Kind = "builtin-type"
	KindPointerType Kind = "pointer-type"
	KindLValueReferenceType Kind = "lvalue-reference-type"
	KindRValueReferenceType Kind = "rvalue-reference-type"
	KindCVQualifiedType Kind = "cv-qualified-type"
	KindArrayType Kind = "array-type"
	KindFunctionType Kind = "function-type" // Native signature syntax; no library Function semantics.
	KindMemberPointerType Kind = "member-pointer-type"
)

// Containment roles. Each qualified-name segment owns its own template arguments.
const (
	RoleInclude = "include"
	RoleDeclaration = "declaration"
	RoleTemplateParameter = "template-parameter"
	RoleType = "type"
	RoleDefault = "default"
	RoleBase = "base"
	RoleMember = "member"
	RoleUnderlying = "underlying"
	RoleEnumerator = "enumerator"
	RoleParameter = "parameter"
	RoleResult = "result"
	RoleSegment = "segment"
	RoleArgument = "argument"
	RoleElement = "element"
	RoleOwner = "owner"
)

// Reference roles. Target links retain declarations, imports or lexical template binders.
const (
	RoleTarget = "target"
)

// Attribute names. Qualified types and template applications are structured children, never source snippets.
const (
	AttributeName = "name"
	AttributePath = "path"
	AttributeSystem = "system"
	AttributeRecordKey = "record-key"
	AttributeForward = "forward"
	AttributeAccess = "access"
	AttributeScoped = "scoped"
	AttributeValue = "value"
	AttributeParameterKey = "parameter-key"
	AttributeGlobal = "global"
	AttributeTypename = "typename"
	AttributeTemplate = "template"
	AttributeConst = "const"
	AttributeVolatile = "volatile"
	AttributeVirtual = "virtual"
	AttributePure = "pure"
	AttributeNoexcept = "noexcept"
	AttributeRefQualifier = "ref-qualifier"
	AttributeVariadic = "variadic"
	AttributeLength = "length"
)
