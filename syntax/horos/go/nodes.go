// Package horos describes canonical declaration syntax over the generic AST.
// It imports no semantic model, parser, engine, generator or behavioral library.
package horos

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type Language struct{}

// bittype:scope docs/syntax.md#identity-and-registration
func (Language) ID() language.ID { return "horos" }
// bittype:scope docs/syntax.md#identity-and-registration
func (Language) Name() string    { return "Horos" }

var _ language.Provider = Language{}

type AST = ast.AST[Language]
type Node = ast.Node[Language]
type Kind = ast.Kind[Language]

// Profile admits declaration syntax before resolution. Binding annotations are
// optional assertions; the separate resolver establishes or rejects them.
const Profile = "horos-declarations-4/ast-1"

// Kind/type definitions have separate binders and bodies, and share a module's
// declaration namespace. Sorts, kinds and sequence patterns are ordinary owned
// AST children. Boolean kinds store their explicit sort. Recursive references
// use lexical group depth and equation index; optional equation names are source
// spelling metadata. Type applications own an arbitrary expression callee, while
// kind applications own a named kind reference. None of these nodes resolves a
// name or establishes a semantic judgment.
const (
	KindProgram                Kind = "program"
	KindModule                 Kind = "module"
	KindImport                 Kind = "import"
	KindTypeDefinition         Kind = "type-definition"
	KindKindDefinition         Kind = "kind-definition"
	KindTypeParameter          Kind = "type-parameter"
	KindKindParameter          Kind = "kind-parameter"
	KindParameter              Kind = "parameter"
	KindReference              Kind = "reference"
	KindApplication            Kind = "application"
	KindDescriptor             Kind = "descriptor"
	KindAtom                   Kind = "atom"
	KindTuple                  Kind = "tuple"
	KindProperSort             Kind = "proper-sort"
	KindArrowSort              Kind = "arrow-sort"
	KindTop                    Kind = "kind-top"
	KindBottom                 Kind = "kind-bottom"
	KindKindParameterReference Kind = "kind-parameter-reference"
	KindKindReference          Kind = "kind-reference"
	KindKindApplication        Kind = "kind-application"
	KindUnion                  Kind = "kind-union"
	KindIntersection           Kind = "kind-intersection"
	KindArrow                  Kind = "kind-arrow"
	KindNodePattern            Kind = "kind-node"
	KindSelector               Kind = "selector"
	KindRecursiveKind          Kind = "kind-recursive"
	KindEquation               Kind = "kind-equation"
	KindRecursiveVariable      Kind = "kind-recursive-variable"
	KindSequenceEpsilon        Kind = "sequence-epsilon"
	KindSequenceItem           Kind = "sequence-item"
	KindSequenceConcat         Kind = "sequence-concat"
	KindSequenceAlternative    Kind = "sequence-alternative"
	KindSequenceStar           Kind = "sequence-star"
)

const (
	AttributeName        = "name"
	AttributePin         = "pin"
	AttributeAlias       = "alias"
	AttributeIndex       = "index"
	AttributeDepth       = "depth"
	AttributeValue       = "value"
	AttributeScope       = "scope"
	AttributeConstructor = "constructor"
	AttributeHex         = "hex"
)

const (
	RoleModule      = "module"
	RoleImport      = "import"
	RoleDeclaration = "declaration"
	RoleParameter   = "parameter"
	RoleContract    = "contract"
	RoleSort        = "sort"
	RoleInput       = "input"
	RoleResult      = "result"
	RoleBody        = "body"
	RoleCallee      = "callee"
	RoleArgument    = "argument"
	RoleSelector    = "selector"
	RolePattern     = "pattern"
	RoleEquation    = "equation"
	RoleOperand     = "operand"
	RoleElement     = "element"
	RoleRoot        = "root"
	RoleTarget      = "target"
)
