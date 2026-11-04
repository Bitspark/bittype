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

var _ language.Provider = Language{}

type AST = ast.AST[Language]
type Node = ast.Node[Language]
type Kind = ast.Kind[Language]

// Profile admits declaration syntax before resolution. Binding annotations are
// optional assertions; the separate resolver establishes or rejects them.
const Profile = "horos-declarations-3/ast-1"

// Interface names are type, data or constructor. Proper interfaces are leaves;
// constructors have nonempty input interfaces followed by a proper result.
// Declarations always contain a proper result, including an implicit Type claim.
// References carry names/binding assertions; applications contain a callee and
// nonempty arguments without duplicating that reference's attributes.
const (
	KindProgram       Kind = "program"
	KindModule        Kind = "module"
	KindImport        Kind = "import"
	KindDeclaration   Kind = "declaration"
	KindTypeParameter Kind = "type-parameter"
	KindInterface     Kind = "interface"
	KindParameter     Kind = "parameter"
	KindReference     Kind = "reference"
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
	RoleInterface   = "interface"
	RoleInput       = "input"
	RoleResult      = "result"
	RoleBody        = "body"
	RoleCallee      = "callee"
	RoleArgument    = "argument"
	RoleOperand     = "operand"
	RoleElement     = "element"
	RoleRoot        = "root"
	RoleTarget      = "target"
)
