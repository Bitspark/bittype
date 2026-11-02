package horosbridge

import (
	"encoding/hex"
	"fmt"
	"math/big"

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
		index, ok := new(big.Int).SetString(r.attr(id, s.AttributeIndex), 10)
		if !ok || index.Sign() < 0 {
			return nil, r.bad(id, "invalid parameter index")
		}
		return tagged("horos/decl/", "param", data.EncodeInt(index)), nil
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

// ValidationError locates a frontend assertion that disagrees with authoritative
// resolution. Context distinguishes document-pin assertions from source syntax.
type ValidationError struct {
	Site ast.NodeRef
	Context bool
	Message string
}

func (e *ValidationError) Error() string { return fmt.Sprintf("%s#%d: %s", e.Site.Tree, e.Site.Node, e.Message) }
func (e *ValidationError) Unwrap() error { return ErrInvalid }

// Decoded retains the source correspondence, not a partially checked Program.
// It is private frontend state consumed synchronously by explicit resolution.
type Decoded struct {
	tree s.AST
	root core.Value
	dependencies []core.Value
	modules map[decl.Digest]ast.NodeID
	declarations map[decl.DeclarationID]ast.NodeID
	parameters map[decl.ParameterID]ast.NodeID
	imports map[importID]ast.NodeID
}

func (d *Decoded) site(id ast.NodeID) ast.NodeRef { return ast.NodeRef{Tree: d.tree.ID, Node: id} }
func (d *Decoded) assertion(id ast.NodeID, context bool, message string) error {
	return &ValidationError{Site: d.site(id), Context: context, Message: message}
}

// Decode reconstructs values from syntax already admitted by the public stage.
// Horos computes every module pin under the selected limits. No descriptor
// positions are guessed; this private bridge never calls a stage provider.
func Decode(view ast.View[s.Language], limits decl.Limits) (*Decoded, error) {
	tree := view.ASTView()
	r := reader{tree: tree}
	d := &Decoded{tree: tree, modules: make(map[decl.Digest]ast.NodeID), declarations: make(map[decl.DeclarationID]ast.NodeID), parameters: make(map[decl.ParameterID]ast.NodeID), imports: make(map[importID]ast.NodeID)}
	root := tree.Nodes[tree.Root-1].References[0].Target.Node
	for _, child := range tree.Nodes[tree.Root-1].Children {
		value, err := r.module(child.Node)
		if err != nil { return nil, err }
		pin, err := decl.PinWithLimits(value, limits)
		if err != nil { return nil, err }
		if r.attr(child.Node, s.AttributePin) != pin.String() { return nil, d.assertion(child.Node, true, "module pin does not match its declaration bytes") }
		if d.modules[pin] != 0 { return nil, d.assertion(child.Node, true, "duplicate module document in input") }
		d.modules[pin] = child.Node
		for _, member := range r.node(child.Node).Children {
			if member.Role == s.RoleImport {
				d.imports[importID{module: pin, alias: r.attr(member.Node, s.AttributeAlias)}] = member.Node
				continue
			}
			id := decl.DeclarationID{Module: pin, Name: r.attr(member.Node, s.AttributeName)}
			d.declarations[id] = member.Node
			index := 0
			for _, parameter := range r.node(member.Node).Children {
				if parameter.Role == s.RoleParameter {
					d.parameters[decl.ParameterID{Declaration: id, Index: index}] = parameter.Node
					index++
				}
			}
		}
		if child.Node == root { d.root = value } else { d.dependencies = append(d.dependencies, value) }
	}
	return d, nil
}

// Resolve returns a complete checked model only after every present annotation
// and every type/literal position agrees with Horos's authoritative projection.
func (d *Decoded) Resolve(limits decl.Limits) (*decl.Program, error) {
	program, err := decl.ResolveWithLimits(d.root, d.dependencies, limits)
	if err != nil { return nil, err }
	expected, err := projectChecked(program, d.tree.ID)
	if err != nil { return nil, err }
	_ = d.correspond(expected)
	return program, nil
}

// Read is an explicitly semantic convenience API. Pack/Check no longer calls it.
func Read(view ast.View[s.Language]) (*decl.Program, error) {
	d, err := Decode(view, decl.DefaultLimits())
	if err != nil { return nil, err }
	return d.Resolve(decl.DefaultLimits())
}

func (d *Decoded) DeclarationSite(id decl.DeclarationID) ast.NodeRef { return d.site(d.declarations[id]) }
func (d *Decoded) ParameterSite(id decl.ParameterID) ast.NodeRef { return d.site(d.parameters[id]) }
func (d *Decoded) ModuleSite(id decl.Digest) ast.NodeRef { return d.site(d.modules[id]) }
func (d *Decoded) BodySite(id decl.DeclarationID) ast.NodeRef {
	node := d.declarations[id]
	if node != 0 { for _, child := range d.tree.Nodes[node-1].Children { if child.Role == s.RoleBody { return d.site(child.Node) } } }
	return ast.NodeRef{}
}

// OperandSite follows a canonical tuple path through the existing AST encoding.
// It interprets only declaration syntax framing, never descriptor constructors.
func (d *Decoded) OperandSite(site ast.NodeRef, path []int) ast.NodeRef {
	if site.Tree != d.tree.ID || site.Node == 0 || int(site.Node) > len(d.tree.Nodes) { return ast.NodeRef{} }
	for len(path) > 0 {
		node := d.tree.Nodes[site.Node-1]
		index := path[0]
		path = path[1:]
		switch node.Kind {
		case s.KindDescriptor:
			index-- // descriptor tag has no separate AST node
		case s.KindApplication:
			if index != 2 || len(path) == 0 { return ast.NodeRef{} }
			index, path = path[0]-1, path[1:] // encoded argument list tag
		case s.KindTuple:
		default:
			return ast.NodeRef{}
		}
		if index < 0 || index >= len(node.Children) { return ast.NodeRef{} }
		site.Node = node.Children[index].Node
	}
	return site
}

func (d *Decoded) ArgumentSite(site ast.NodeRef, index int) ast.NodeRef {
	return d.OperandSite(site, []int{2, index+1})
}

func (d *Decoded) Locate(detail *decl.Error) *ast.NodeRef {
	if detail == nil { return nil }
	var site ast.NodeRef
	if detail.Declaration != "" {
		id := decl.DeclarationID{Module: detail.Module, Name: detail.Declaration}
		site = d.BodySite(id)
		if len(detail.Path) > 0 { if exact := d.OperandSite(site, detail.Path); exact.Node != 0 { site = exact } }
	} else if detail.ImportAlias != "" {
		site = d.site(d.imports[importID{module: detail.Module, alias: detail.ImportAlias}])
	} else {
		site = d.ModuleSite(detail.Module)
	}
	if site.Node == 0 { return nil }
	return &site
}

func (d *Decoded) correspond(expected s.AST) error {
	mapping := map[ast.NodeID]ast.NodeID{expected.Root: d.tree.Root}
	wantReader := reader{tree: expected}
	gotReader := reader{tree: d.tree}
	var walk func(ast.NodeID, ast.NodeID) error
	walk = func(wantID, gotID ast.NodeID) error {
		want, got := expected.Nodes[wantID-1], d.tree.Nodes[gotID-1]
		bad := func() error { return d.assertion(gotID, false, "Horos syntax does not match checked declaration positions") }
		if want.Kind != got.Kind || len(want.Attributes) != len(got.Attributes) || len(want.Children) != len(got.Children) { return bad() }
		for _, attribute := range want.Attributes { if gotReader.attr(gotID, attribute.Name) != attribute.Text { return bad() } }
		mapping[wantID] = gotID
		for i, child := range want.Children {
			actual := got.Children[i]
			if want.Kind == s.KindModule {
				name := s.AttributeName
				if child.Role == s.RoleImport { name = s.AttributeAlias }
				for _, candidate := range got.Children {
					if candidate.Role == child.Role && gotReader.attr(candidate.Node, name) == wantReader.attr(child.Node, name) { actual = candidate; break }
				}
			}
			if child.Role != actual.Role { return bad() }
			if err := walk(child.Node, actual.Node); err != nil { return err }
		}
		return nil
	}
	for _, child := range expected.Nodes[expected.Root-1].Children {
		bytes, _ := hex.DecodeString(wantReader.attr(child.Node, s.AttributePin))
		var pin decl.Digest
		copy(pin[:], bytes)
		actual := d.modules[pin]
		if actual == 0 { return fmt.Errorf("%w: resolved module has no source", ErrInvalid) }
		if err := walk(child.Node, actual); err != nil { return err }
	}
	for wantID, gotID := range mapping {
		want := expected.Nodes[wantID-1]
		for _, actual := range d.tree.Nodes[gotID-1].References {
			found := false
			for _, reference := range want.References {
				if reference.Role == actual.Role && actual.Target.Tree == d.tree.ID && actual.Target.Node == mapping[reference.Target.Node] { found = true; break }
			}
			if !found { return d.assertion(gotID, false, "incorrect declaration, import or parameter annotation") }
		}
	}
	return nil
}
