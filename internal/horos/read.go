package horosbridge

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"unicode/utf8"

	structure "github.com/Bitspark/bittype/check/ast/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

func tagged(prefix, name string, args ...core.Value) core.Value {
	return core.NewTuple(append([]core.Value{core.NewAtom([]byte(prefix + name))}, args...)...)
}

type reader struct{ tree s.AST }

func (r reader) node(id ast.NodeID) s.Node { return r.tree.Nodes[id-1] }
func (r reader) attr(id ast.NodeID, name string) string {
	for _, attribute := range r.node(id).Attributes {
		if attribute.Name == name {
			return attribute.Text
		}
	}
	return ""
}
func (r reader) bad(id ast.NodeID, message string) error {
	return text.AtNode(r.tree, id, message)
}
func (r reader) value(id ast.NodeID) (core.Value, error) {
	n := r.node(id)
	switch n.Kind {
	case s.KindAtom:
		bytes, err := hex.DecodeString(r.attr(id, s.AttributeHex))
		if err != nil {
			return nil, r.bad(id, "invalid atom hex")
		}
		return core.NewAtom(bytes), nil
	case s.KindParameter:
		index, err := strconv.Atoi(r.attr(id, s.AttributeIndex))
		if err != nil || index < 0 {
			return nil, r.bad(id, "invalid parameter index")
		}
		return tagged("horos/decl/", "param", data.EncodeInt(big.NewInt(int64(index)))), nil
	case s.KindTuple, s.KindDescriptor, s.KindApplication:
		values := make([]core.Value, len(n.Children))
		for i, child := range n.Children {
			value, err := r.value(child.Node)
			if err != nil {
				return nil, err
			}
			values[i] = value
		}
		if n.Kind == s.KindTuple {
			return core.NewTuple(values...), nil
		}
		if n.Kind == s.KindDescriptor {
			return tagged("horos/", r.attr(id, s.AttributeConstructor), values...), nil
		}
		var target core.Value
		switch r.attr(id, s.AttributeScope) {
		case "local":
			target = tagged("horos/decl/", "local", data.MustEncodeText(r.attr(id, s.AttributeName)))
		case "imported":
			target = tagged("horos/decl/", "imported", data.MustEncodeText(r.attr(id, s.AttributeAlias)), data.MustEncodeText(r.attr(id, s.AttributeName)))
		default:
			return nil, r.bad(id, "invalid application scope")
		}
		return tagged("horos/decl/", "apply", target, data.EncodeList(values)), nil
	default:
		return nil, r.bad(id, "expected a descriptor expression or value fragment")
	}
}

func (r reader) declaration(id ast.NodeID) (core.Value, error) {
	if r.node(id).Kind != s.KindDeclaration {
		return nil, r.bad(id, "expected declaration")
	}
	var parameters []core.Value
	var body core.Value
	for _, child := range r.node(id).Children {
		switch child.Role {
		case s.RoleParameter:
			if r.node(child.Node).Kind != s.KindTypeParameter {
				return nil, r.bad(child.Node, "expected parameter declaration")
			}
			parameters = append(parameters, data.MustEncodeText(r.attr(child.Node, s.AttributeName)))
		case s.RoleBody:
			if body != nil {
				return nil, r.bad(id, "duplicate declaration body")
			}
			value, err := r.value(child.Node)
			if err != nil {
				return nil, err
			}
			body = value
		default:
			return nil, r.bad(child.Node, "unknown declaration child role")
		}
	}
	if body == nil {
		return nil, r.bad(id, "missing declaration body")
	}
	return tagged("horos/decl/", "type", data.EncodeList(parameters), body), nil
}

func (r reader) module(id ast.NodeID) (core.Value, error) {
	if r.node(id).Kind != s.KindModule {
		return nil, r.bad(id, "expected module")
	}
	var imports, declarations [][2]core.Value
	for _, child := range r.node(id).Children {
		switch child.Role {
		case s.RoleImport:
			if r.node(child.Node).Kind != s.KindImport {
				return nil, r.bad(child.Node, "expected import")
			}
			bytes, err := hex.DecodeString(r.attr(child.Node, s.AttributePin))
			if err != nil {
				return nil, r.bad(child.Node, "invalid import pin")
			}
			imported := tagged("horos/decl/", "import",
				data.MustEncodeText(r.attr(child.Node, s.AttributeName)), core.NewAtom(bytes))
			imports = append(imports, [2]core.Value{data.MustEncodeText(r.attr(child.Node, s.AttributeAlias)), imported})
		case s.RoleDeclaration:
			value, err := r.declaration(child.Node)
			if err != nil {
				return nil, err
			}
			declarations = append(declarations, [2]core.Value{data.MustEncodeText(r.attr(child.Node, s.AttributeName)), value})
		default:
			return nil, r.bad(child.Node, "unknown module child role")
		}
	}
	importsValue, err := data.EncodeMap(imports)
	if err != nil {
		return nil, r.bad(id, "duplicate or malformed imports")
	}
	declarationsValue, err := data.EncodeMap(declarations)
	if err != nil {
		return nil, r.bad(id, "duplicate or malformed declarations")
	}
	return tagged("horos/decl/", "module", data.EncodeInt(big.NewInt(1)),
		data.MustEncodeText(r.attr(id, s.AttributeName)), importsValue, declarationsValue), nil
}

// Read validates AST structure and reuses Horos for all declaration semantics.
// Reprojection checks the supplied binding edges and literal/type boundaries;
// it never repairs the input tree or invents missing references.
func Read(view ast.View[s.Language]) (*decl.Program, error) {
	if view == nil {
		return nil, ErrInvalid
	}
	reflected := reflect.ValueOf(view)
	switch reflected.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface, reflect.Func, reflect.Chan:
		if reflected.IsNil() {
			return nil, ErrInvalid
		}
	}
	tree := view.ASTView()
	if err := structure.Structure(tree); err != nil {
		return nil, err
	}
	// Canonical names are Ontos text values. Reject invalid Go strings before
	// calling the text embedding; a malformed AST must never panic.
	for i, node := range tree.Nodes {
		for _, attribute := range node.Attributes {
			if !utf8.ValidString(attribute.Text) {
				return nil, text.AtNode(tree, ast.NodeID(i+1), "attribute text is not UTF-8")
			}
		}
	}
	r := reader{tree: tree}
	if r.node(tree.Root).Kind != s.KindProgram {
		return nil, r.bad(tree.Root, "expected program root")
	}
	root := ast.NodeID(0)
	for _, reference := range r.node(tree.Root).References {
		if reference.Role == s.RoleRoot && reference.Target.Tree == tree.ID {
			root = reference.Target.Node
		}
	}
	var rootValue core.Value
	var dependencies []core.Value
	for _, child := range r.node(tree.Root).Children {
		if child.Role != s.RoleModule {
			return nil, r.bad(child.Node, "expected module child")
		}
		value, err := r.module(child.Node)
		if err != nil {
			return nil, err
		}
		if child.Node == root {
			rootValue = value
		} else {
			dependencies = append(dependencies, value)
		}
	}
	if rootValue == nil {
		return nil, r.bad(tree.Root, "root reference must name a contained module")
	}
	program, err := decl.Resolve(rootValue, dependencies)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	expected, err := Project(program, tree.ID)
	if err != nil {
		return nil, err
	}
	if err := equivalent(expected, tree); err != nil {
		return nil, err
	}
	return program, nil
}

// Arena addresses and provenance may differ; the normalized containment and
// semantic edges must match. No structural hash selects a behavioral contract.
func equivalent(expected, actual s.AST) error {
	mapping := make(map[ast.NodeID]ast.NodeID)
	var walk func(ast.NodeID, ast.NodeID) error
	walk = func(wantID, gotID ast.NodeID) error {
		want, got := expected.Nodes[wantID-1], actual.Nodes[gotID-1]
		bad := func() error { return text.AtNode(actual, gotID, "Horos syntax or binding does not match checked declarations") }
		if want.Kind != got.Kind || len(want.Attributes) != len(got.Attributes) || len(want.Children) != len(got.Children) || len(want.References) != len(got.References) {
			return bad()
		}
		attrs := make(map[string]string)
		for _, attribute := range got.Attributes {
			attrs[attribute.Name] = attribute.Text
		}
		for _, attribute := range want.Attributes {
			if value, ok := attrs[attribute.Name]; !ok || value != attribute.Text {
				return bad()
			}
		}
		mapping[wantID] = gotID
		for i, child := range want.Children {
			if child.Role != got.Children[i].Role {
				return bad()
			}
			if err := walk(child.Node, got.Children[i].Node); err != nil {
				return err
			}
		}
		return nil
	}
	if err := walk(expected.Root, actual.Root); err != nil {
		return err
	}
	for wantID, gotID := range mapping {
		got := actual.Nodes[gotID-1]
		refs := make(map[string]ast.NodeRef)
		for _, reference := range got.References {
			refs[reference.Role] = reference.Target
		}
		for _, reference := range expected.Nodes[wantID-1].References {
			target, ok := refs[reference.Role]
			if !ok || target.Tree != actual.ID || target.Node != mapping[reference.Target.Node] {
				return text.AtNode(actual, gotID, "incorrect declaration, import or parameter reference")
			}
		}
	}
	return nil
}
