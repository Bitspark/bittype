// Package target constructs a horos AST from a closed target program: horos
// as the TARGET of a generation.
//
// It exists for the pair typescript-to-horos, the second of the two other
// pairs of bittype#196 (decision O43). It builds one module of type
// definitions whose bodies are closed records of required fields and the
// scalars text and bool: the forms of the horos model beside it. The horos
// syntax checker admits the tree and the horos formatter prints it.
//
// Two things make a horos tree unlike the tree of a native target, and both
// are findings of the two other pairs (finding 2):
//
//   - A module node carries the pin of its own canonical bytes, and a field's
//     name is the literal layout of an ontos text value. So this package
//     imports the private module lowering of the horos frontend
//     (internal/horosmodel) and the ontos value packages, which the role of a
//     target does not give: the architecture test names both as an extension.
//   - The program node refers to its root module. The shared construction
//     lets a lowering refer to a planned declaration or parameter and to
//     nothing else, so this package completes the tree after the shared
//     construction has returned it.
//
// It refuses, as unsupported: a generic definition and a reference from one
// definition to another. Neither is needed by the pair.
//
// COPY (AC-11, finding AG-9 of the two other pairs): CheckExpression and
// CheckProgram forward to the model and to the shared single-scope check, and
// Build starts as a store and the shared compile, as in the target of every
// language. That wiring interprets nothing of horos and has no agnostic
// implementation: each target repeats it. What is horos's here is the
// lowering of each form, the spelling check of a name and the completion of
// the tree.
package target

import (
	"context"
	"encoding/hex"
	"strings"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	native "github.com/Bitspark/bittype/internal/horos"
	"github.com/Bitspark/bittype/internal/horosmodel"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/horos/model/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// Provider is the horos target. Module is the name of the module that the
// generated definitions form: a horos tree is a program of modules, and a
// generation request names none.
type Provider struct{ Module string }

type backend struct {
	*m.Store[s.Language]
	module string
}

// moduleName reports whether a name is a horos module name: identifiers
// joined by dots.
func moduleName(name string) bool {
	for _, segment := range strings.Split(name, ".") {
		if !native.Identifier(segment) {
			return false
		}
	}
	return true
}

func (p Provider) CheckExpression(ctx context.Context, e g.Expression[s.Language]) error {
	return model.CheckExpression(ctx, e)
}

// CheckDeclaration adds to the model's check the spelling rule of horos: the
// name of a definition is a horos identifier.
func (p Provider) CheckDeclaration(ctx context.Context, d g.TargetDeclaration[s.Language]) error {
	if err := model.CheckDeclaration(ctx, d); err != nil {
		return err
	}
	if !native.Identifier(d.Name) {
		return g.ErrInvalid
	}
	return nil
}

// CheckProgram requires distinct definition names: the definitions form one
// module, which is one declaration namespace.
func (p Provider) CheckProgram(ctx context.Context, program g.Program[s.Language]) error {
	return m.CheckSingleScope(ctx, program)
}

func (p Provider) Build(ctx context.Context, in g.BuildInput[s.Language]) (g.Assembly, error) {
	if !moduleName(p.Module) {
		return g.Assembly{}, g.ErrInvalid
	}
	b := &backend{Store: m.New[s.Language](ctx, in.Limits), module: p.Module}
	assembly, err := m.Compile(ctx, in, b.Store, b)
	if err != nil {
		return g.Assembly{}, err
	}
	// The shared construction returned one tree whose root is the program
	// node. Complete it: the root reference and the pin of the module.
	tree, err := representation.Unpack[s.Language](in.Language, assembly.Syntax.EntryAST())
	if err != nil {
		return g.Assembly{}, err
	}
	program := &tree.Nodes[tree.Root-1]
	if program.Kind != s.KindProgram || len(program.Children) != 1 {
		return g.Assembly{}, g.ErrInvalid
	}
	module := program.Children[0].Node
	program.References = []ast.Reference{{Role: s.RoleRoot, Target: ast.NodeRef{Tree: tree.ID, Node: module}}}
	pins, err := horosmodel.SourcePins(tree)
	if err != nil {
		return g.Assembly{}, err
	}
	pin, ok := pins[module]
	if !ok {
		return g.Assembly{}, g.ErrInvalid
	}
	attributes := tree.Nodes[module-1].Attributes
	for i := range attributes {
		if attributes[i].Name == s.AttributePin {
			attributes[i].Text = pin
		}
	}
	assembly.Syntax, err = representation.PackSet(in.Language, ast.Set[s.Language]{Entry: ast.NodeRef{Tree: tree.ID, Node: tree.Root}, Trees: []ast.AST[s.Language]{tree}})
	if err != nil {
		return g.Assembly{}, err
	}
	return assembly, nil
}

func (b *backend) atom(value string, origin ast.NodeRef) (m.Fragment[s.Language], error) {
	return b.Node(m.Node(s.KindAtom, m.Attributes(s.AttributeHex, hex.EncodeToString([]byte(value))), origin))
}

// literal writes an ontos value as the literal atoms and tuples of a horos
// tree: the tree carries canonical values literally (syntax/horos/go).
func (b *backend) literal(value core.Value, origin ast.NodeRef) (m.Fragment[s.Language], error) {
	switch v := value.(type) {
	case core.Atom:
		return b.atom(string(v.Bytes()), origin)
	case core.Tuple:
		children := make([]m.Child[s.Language], v.Len())
		for i := range children {
			element, err := b.literal(v.At(i), origin)
			if err != nil {
				return element, err
			}
			children[i] = m.Child[s.Language]{Role: s.RoleElement, Fragment: element}
		}
		return b.Node(m.Node(s.KindTuple, nil, origin), children...)
	}
	return m.Fragment[s.Language]{}, g.ErrInvalid
}

// Parameter and Reference are the two lowerings that this target does not
// have: a generic definition and a reference from one definition to another.
func (b *backend) Parameter(string, ast.NodeRef) (m.Fragment[s.Language], error) {
	return m.Fragment[s.Language]{}, g.ErrUnsupported
}
func (b *backend) Reference(m.Symbol[s.Language], []m.Fragment[s.Language], ast.NodeRef) (m.Fragment[s.Language], error) {
	return m.Fragment[s.Language]{}, g.ErrUnsupported
}

func (b *backend) Expression(e g.Expression[s.Language], args []m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	if err := model.CheckExpression(context.Background(), e); err != nil {
		return m.Fragment[s.Language]{}, err
	}
	name, _ := model.Name(e.Form())
	values, origin := e.Data(), e.Origin()
	operands := func(fragments ...m.Fragment[s.Language]) []m.Child[s.Language] {
		out := make([]m.Child[s.Language], len(fragments))
		for i, f := range fragments {
			out[i] = m.Child[s.Language]{Role: s.RoleOperand, Fragment: f}
		}
		return out
	}
	switch name {
	case model.FormScalar:
		selector, err := b.atom(values[0], origin)
		if err != nil {
			return selector, err
		}
		return b.Node(m.Node(s.KindDescriptor, m.Attributes(s.AttributeConstructor, "scalar"), origin), operands(selector)...)
	case model.FormField:
		// The wrapper of one record field: (FieldTag, name, required, type).
		tag, err := b.atom(s.FieldTag, origin)
		if err != nil {
			return tag, err
		}
		text, err := data.EncodeText(values[0])
		if err != nil {
			return m.Fragment[s.Language]{}, g.ErrInvalid
		}
		label, err := b.literal(text, origin)
		if err != nil {
			return label, err
		}
		required, err := b.atom(s.FieldRequired, origin)
		if err != nil {
			return required, err
		}
		children := make([]m.Child[s.Language], 0, 4)
		for _, element := range []m.Fragment[s.Language]{tag, label, required, args[0]} {
			children = append(children, m.Child[s.Language]{Role: s.RoleElement, Fragment: element})
		}
		return b.Node(m.Node(s.KindTuple, nil, origin), children...)
	case model.FormRecord:
		mode, err := b.atom(s.RecordClosed, origin)
		if err != nil {
			return mode, err
		}
		return b.Node(m.Node(s.KindDescriptor, m.Attributes(s.AttributeConstructor, s.RecordLabel), origin), operands(append([]m.Fragment[s.Language]{mode}, args...)...)...)
	}
	return m.Fragment[s.Language]{}, g.ErrUnsupported
}

// Declaration builds a type definition with the result kind that an omitted
// annotation elaborates to: top of the proper sort.
func (b *backend) Declaration(d g.TargetDeclaration[s.Language], params []m.Fragment[s.Language], body m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	if len(params) != 0 {
		return body, g.ErrUnsupported
	}
	sort, err := b.Node(m.Node(s.KindProperSort, nil, ast.NodeRef{}))
	if err != nil {
		return sort, err
	}
	result, err := b.Node(m.Node(s.KindTop, nil, ast.NodeRef{}), m.Child[s.Language]{Role: s.RoleSort, Fragment: sort})
	if err != nil {
		return result, err
	}
	return b.Node(m.Node(s.KindTypeDefinition, m.Attributes(s.AttributeName, d.Name), d.Origin),
		m.Child[s.Language]{Role: s.RoleResult, Fragment: result}, m.Child[s.Language]{Role: s.RoleBody, Fragment: body})
}

// Finish builds the module of the definitions and the program that holds it.
// The pin is a placeholder of the right form until Build has computed it.
func (b *backend) Finish(declarations []m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	children := make([]m.Child[s.Language], len(declarations))
	for i, f := range declarations {
		children[i] = m.Child[s.Language]{Role: s.RoleDeclaration, Fragment: f}
	}
	module, err := b.Node(m.Node(s.KindModule, m.Attributes(s.AttributeName, b.module, s.AttributePin, strings.Repeat("0", 64)), ast.NodeRef{}), children...)
	if err != nil {
		return module, err
	}
	return b.Node(m.Node(s.KindProgram, nil, ast.NodeRef{}), m.Child[s.Language]{Role: s.RoleModule, Fragment: module})
}

var _ g.Target[s.Language] = Provider{}
