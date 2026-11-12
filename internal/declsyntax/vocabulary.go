// Package declsyntax shares declaration AST framing, not type semantics.
package declsyntax

import h "github.com/Bitspark/bittype/syntax/horos/go"

// The horos syntax vocabulary owns every name below. The admission grammar is
// a table keyed by strings, so each kind is the vocabulary's constant converted
// to a string. No name is spelled here.
const (
	KindProgram                = string(h.KindProgram)
	KindModule                 = string(h.KindModule)
	KindImport                 = string(h.KindImport)
	KindTypeDefinition         = string(h.KindTypeDefinition)
	KindKindDefinition         = string(h.KindKindDefinition)
	KindTypeParameter          = string(h.KindTypeParameter)
	KindKindParameter          = string(h.KindKindParameter)
	KindParameter              = string(h.KindParameter)
	KindReference              = string(h.KindReference)
	KindApplication            = string(h.KindApplication)
	KindDescriptor             = string(h.KindDescriptor)
	KindAtom                   = string(h.KindAtom)
	KindTuple                  = string(h.KindTuple)
	KindProperSort             = string(h.KindProperSort)
	KindArrowSort              = string(h.KindArrowSort)
	KindTop                    = string(h.KindTop)
	KindBottom                 = string(h.KindBottom)
	KindKindParameterReference = string(h.KindKindParameterReference)
	KindKindReference          = string(h.KindKindReference)
	KindKindApplication        = string(h.KindKindApplication)
	KindUnion                  = string(h.KindUnion)
	KindIntersection           = string(h.KindIntersection)
	KindArrow                  = string(h.KindArrow)
	KindNodePattern            = string(h.KindNodePattern)
	KindSelector               = string(h.KindSelector)
	KindRecursiveKind          = string(h.KindRecursiveKind)
	KindEquation               = string(h.KindEquation)
	KindRecursiveVariable      = string(h.KindRecursiveVariable)
	KindSequenceEpsilon        = string(h.KindSequenceEpsilon)
	KindSequenceItem           = string(h.KindSequenceItem)
	KindSequenceConcat         = string(h.KindSequenceConcat)
	KindSequenceAlternative    = string(h.KindSequenceAlternative)
	KindSequenceStar           = string(h.KindSequenceStar)
)

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
