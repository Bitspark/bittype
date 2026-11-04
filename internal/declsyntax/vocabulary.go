// Package declsyntax shares declaration AST framing, not type semantics.
package declsyntax

const (
	KindProgram       = "program"
	KindModule        = "module"
	KindImport        = "import"
	KindDeclaration   = "declaration"
	KindTypeParameter = "type-parameter"
	KindInterface     = "interface"
	KindParameter     = "parameter"
	KindReference     = "reference"
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
