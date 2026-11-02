// Package horos describes canonical declaration syntax over the generic AST.
// It imports no semantic model, parser, engine, generator or behavioral library.
package horos

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type Language struct{}

func (Language) ID() language.ID { return "horos" }
func (Language) Name() string    { return "Horos" }

var _ language.LanguageProvider = Language{}

type AST = ast.AST[Language]
type Node = ast.Node[Language]
type Kind = ast.Kind[Language]

// Profile identifies this AST-only projection of declaration revision 1.
const Profile = "horos-declarations-1/ast-1"

const (
	KindProgram       Kind = "program"
	KindModule        Kind = "module"
	KindImport        Kind = "import"
	KindDeclaration   Kind = "declaration"
	KindTypeParameter Kind = "type-parameter"
	KindParameter     Kind = "parameter"
	KindApplication   Kind = "application"
	KindDescriptor    Kind = "descriptor"
	KindAtom          Kind = "atom"
	KindTuple         Kind = "tuple"
)

const (
	AttributeName        = "name"
	AttributePin         = "pin"
	AttributeAlias       = "alias"
	AttributeIndex       = "index"
	AttributeScope       = "scope"
	AttributeConstructor = "constructor"
	AttributeHex         = "hex"
)

const (
	RoleModule      = "module"
	RoleImport      = "import"
	RoleDeclaration = "declaration"
	RoleParameter   = "parameter"
	RoleBody        = "body"
	RoleArgument    = "argument"
	RoleOperand     = "operand"
	RoleElement     = "element"
	RoleRoot        = "root"
	RoleTarget      = "target"
)
