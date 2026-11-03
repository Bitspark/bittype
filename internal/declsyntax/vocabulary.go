// Package declsyntax shares declaration AST framing, not type semantics.
package declsyntax

const (
	KindProgram       = "program"
	KindModule        = "module"
	KindImport        = "import"
	KindDeclaration   = "declaration"
	KindTypeParameter = "type-parameter"
	KindParameter     = "parameter"
	KindApplication   = "application"
	KindDescriptor    = "descriptor"
	KindAtom          = "atom"
	KindTuple         = "tuple"
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
