// Package rust declares the Rust type-syntax vocabulary over the generic AST.
// It contains no parser, builder, grammar checker, generator or formatter.
// The node shapes and this vocabulary's limits are documented in docs/syntax.md.
package rust

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/language/go"
)

// Language identifies this syntax family and is its distinct Go AST marker.
type Language struct{}

// ID is the stable registration key for this language family.
func (Language) ID() language.ID { return "rust" }

// Name is the human-readable language name, not a registration key.
func (Language) Name() string { return "Rust" }

var _ language.Language = Language{}

// AST is the existing generic carrier specialized for this language.
type AST = ast.AST[Language]

// Node is a node in this language's unchecked syntax carrier.
type Node = ast.Node[Language]

// Kind is this language's syntax discriminator.
type Kind = ast.Kind[Language]

// Node kinds for items, generic binders, paths and type expressions.
const (
	KindFile Kind = "file"
	KindModule Kind = "module"
	KindUse Kind = "use"
	KindTypeAlias Kind = "type-alias"
	KindStruct Kind = "struct"
	KindTupleStruct Kind = "tuple-struct"
	KindEnum Kind = "enum"
	KindUnitVariant Kind = "unit-variant"
	KindTupleVariant Kind = "tuple-variant"
	KindRecordVariant Kind = "record-variant"
	KindField Kind = "field"
	KindTupleField Kind = "tuple-field"
	KindTrait Kind = "trait"
	KindAssociatedType Kind = "associated-type"
	KindMethodSignature Kind = "method-signature"
	KindReceiver Kind = "receiver"
	KindParameter Kind = "parameter"
	KindTypeParameter Kind = "type-parameter"
	KindLifetimeParameter Kind = "lifetime-parameter"
	KindLifetimeReference Kind = "lifetime-reference"
	KindWherePredicate Kind = "where-predicate"
	KindPathType Kind = "path-type"
	KindPathSegment Kind = "path-segment"
	KindPrimitiveType Kind = "primitive-type"
	KindReferenceType Kind = "reference-type"
	KindRawPointerType Kind = "raw-pointer-type"
	KindSliceType Kind = "slice-type"
	KindArrayType Kind = "array-type"
	KindTupleType Kind = "tuple-type"
	KindNeverType Kind = "never-type"
	KindFunctionType Kind = "function-type" // Native bare fn signature; no library Function semantics.
	KindTraitObjectType Kind = "trait-object-type"
)

// Containment roles. Generic arguments may be type nodes or lifetime references.
const (
	RoleItem = "item"
	RolePath = "path"
	RoleGenericParameter = "generic-parameter"
	RoleType = "type"
	RoleDefault = "default"
	RoleBound = "bound"
	RoleWhere = "where"
	RoleSubject = "subject"
	RoleField = "field"
	RoleVariant = "variant"
	RoleMember = "member"
	RoleReceiver = "receiver"
	RoleParameter = "parameter"
	RoleResult = "result"
	RoleSegment = "segment"
	RoleArgument = "argument"
	RoleElement = "element"
	RoleLifetime = "lifetime"
)

// Reference roles. A path segment or lifetime use can retain its resolved binder.
const (
	RoleTarget = "target"
)

// Attribute names. Lifetimes store their identifier without the apostrophe.
const (
	AttributeName = "name"
	AttributeAlias = "alias"
	AttributeVisibility = "visibility"
	AttributeAbsolute = "absolute"
	AttributeMutable = "mutable"
	AttributeLength = "length"
	AttributeReceiverForm = "receiver-form"
	AttributeUnsafe = "unsafe"
	AttributeABI = "abi"
)
