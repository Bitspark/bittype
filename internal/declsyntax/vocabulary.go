// Package declsyntax shares declaration AST framing, not type semantics.
package declsyntax

import h "github.com/Bitspark/bittype/syntax/horos/go"

// The horos syntax vocabulary owns every name below. The admission grammar is
// typed by language and takes the node kinds from the vocabulary directly.
// Attribute and role names are strings. No name is spelled here.
const (
	AttributeName        = h.AttributeName
	AttributePin         = h.AttributePin
	AttributeAlias       = h.AttributeAlias
	AttributeIndex       = h.AttributeIndex
	AttributeDepth       = h.AttributeDepth
	AttributeValue       = h.AttributeValue
	AttributeScope       = h.AttributeScope
	AttributeConstructor = h.AttributeConstructor
	AttributeHex         = h.AttributeHex
)

const (
	RoleModule      = h.RoleModule
	RoleImport      = h.RoleImport
	RoleDeclaration = h.RoleDeclaration
	RoleParameter   = h.RoleParameter
	RoleContract    = h.RoleContract
	RoleSort        = h.RoleSort
	RoleInput       = h.RoleInput
	RoleResult      = h.RoleResult
	RoleBody        = h.RoleBody
	RoleCallee      = h.RoleCallee
	RoleArgument    = h.RoleArgument
	RoleSelector    = h.RoleSelector
	RolePattern     = h.RolePattern
	RoleEquation    = h.RoleEquation
	RoleOperand     = h.RoleOperand
	RoleElement     = h.RoleElement
	RoleRoot        = h.RoleRoot
	RoleTarget      = h.RoleTarget
)
