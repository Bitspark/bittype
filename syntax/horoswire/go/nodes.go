// Package horoswire describes canonical declaration syntax over the generic AST.
// It imports no semantic model, parser, engine, generator or behavioral library.
package horoswire

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type Language struct{}

func (Language) ID() language.ID { return "horoswire" }
func (Language) Name() string    { return "Horos/Wire" }

var _ language.Provider = Language{}

type AST = ast.AST[Language]
type Node = ast.Node[Language]
type Kind = ast.Kind[Language]

// Profile admits declaration syntax before resolution. Binding annotations are
// optional assertions; the separate analyzer establishes or rejects them.
const Profile = "horos-wire-1/ast-1"

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
