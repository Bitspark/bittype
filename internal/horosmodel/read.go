package horosmodel

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/Bitspark/bittype/contracts/text/go"
	s "github.com/Bitspark/bittype/internal/declsyntax"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

func tagged(prefix, name string, args ...core.Value) core.Value {
	return core.NewTuple(append([]core.Value{core.NewAtom([]byte(prefix + name))}, args...)...)
}

type reader[L any] struct {
	tree                                ast.AST[L]
	declarationPrefix, descriptorPrefix string
}

func (r reader[L]) node(id ast.NodeID) ast.Node[L] { return r.tree.Nodes[id-1] }
func (r reader[L]) attr(id ast.NodeID, name string) string {
	for _, attribute := range r.node(id).Attributes {
		if attribute.Name == name {
			return attribute.Text
		}
	}
	return ""
}
func (r reader[L]) bad(id ast.NodeID, message string) error {
	return text.AtNode(r.tree, id, message)
}
func (r reader[L]) value(id ast.NodeID) (core.Value, error) {
	n := r.node(id)
	switch n.Kind {
	case s.KindAtom:
		bytes, err := hex.DecodeString(r.attr(id, s.AttributeHex))
		if err != nil {
			return nil, r.bad(id, "invalid atom hex")
		}
		return core.NewAtom(bytes), nil
	case s.KindParameter:
		spelling := r.attr(id, s.AttributeIndex)
		// Horos's declaration slot is a Go int. Bound the decimal work before
		// conversion; even a 64-bit int needs fewer than twenty decimal digits.
		if len(spelling) > 20 {
			return nil, r.bad(id, "invalid parameter index")
		}
		index, err := strconv.Atoi(spelling)
		if err != nil || index < 0 {
			return nil, r.bad(id, "invalid parameter index")
		}
		return tagged(r.declarationPrefix, "param", data.EncodeInt(big.NewInt(int64(index)))), nil
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
			return tagged(r.descriptorPrefix, r.attr(id, s.AttributeConstructor), values...), nil
		}
		var target core.Value
		switch r.attr(id, s.AttributeScope) {
		case "local":
			target = tagged(r.declarationPrefix, "local", data.MustEncodeText(r.attr(id, s.AttributeName)))
		case "imported":
			target = tagged(r.declarationPrefix, "imported", data.MustEncodeText(r.attr(id, s.AttributeAlias)), data.MustEncodeText(r.attr(id, s.AttributeName)))
		default:
			return nil, r.bad(id, "invalid application scope")
		}
		return tagged(r.declarationPrefix, "apply", target, data.EncodeList(values)), nil
	default:
		return nil, r.bad(id, "expected a descriptor expression or value fragment")
	}
}

func (r reader[L]) declaration(id ast.NodeID) (core.Value, error) {
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
			parameter := r.node(child.Node)
			if len(parameter.Children) != 1 || parameter.Children[0].Role != s.RoleBound || r.node(parameter.Children[0].Node).Kind != s.KindBound {
				return nil, r.bad(child.Node, "expected an explicit parameter bound")
			}
			bound := r.attr(parameter.Children[0].Node, s.AttributeName)
			parameters = append(parameters, tagged(r.declarationPrefix, "parameter",
				data.MustEncodeText(r.attr(child.Node, s.AttributeName)), core.NewAtom([]byte(bound))))
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
	return tagged(r.declarationPrefix, "type", data.EncodeList(parameters), body), nil
}

func (r reader[L]) module(id ast.NodeID) (core.Value, error) {
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
			imported := tagged(r.declarationPrefix, "import",
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
	return tagged(r.declarationPrefix, "module", data.EncodeInt(big.NewInt(2)),
		data.MustEncodeText(r.attr(id, s.AttributeName)), importsValue, declarationsValue), nil
}

// ValidationError locates a frontend assertion that disagrees with authoritative
// resolution. Context distinguishes document-pin assertions from source syntax.
type ValidationError struct {
	Site       ast.NodeRef
	Context    bool
	Annotation bool
	Message    string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s#%d: %s", e.Site.Tree, e.Site.Node, e.Message)
}
func (e *ValidationError) Unwrap() error { return ErrInvalid }

// Decoded retains the source correspondence, not a partially checked Program.
// It is private frontend state consumed synchronously by explicit resolution.
type Decoded[L any] struct {
	frontend     Frontend[L]
	tree         ast.AST[L]
	root         core.Value
	dependencies []core.Value
	modules      map[decl.Digest]ast.NodeID
	declarations map[decl.DeclarationID]ast.NodeID
	parameters   map[decl.ParameterID]ast.NodeID
	imports      map[importID]ast.NodeID
}

func (d *Decoded[L]) site(id ast.NodeID) ast.NodeRef { return ast.NodeRef{Tree: d.tree.ID, Node: id} }
func (d *Decoded[L]) assertion(id ast.NodeID, context bool, message string) error {
	return &ValidationError{Site: d.site(id), Context: context, Message: message}
}

func (d *Decoded[L]) annotation(id ast.NodeID, context bool, message string) error {
	return &ValidationError{Site: d.site(id), Context: context, Annotation: true, Message: message}
}

// Decode reconstructs values from syntax already admitted by the public stage.
// Horos computes every module pin under the selected limits. No descriptor
// positions are guessed; this private bridge never calls a stage provider.
func (f Frontend[L]) Decode(view ast.View[L], limits decl.Limits) (*Decoded[L], error) {
	tree := view.ASTView()
	r := reader[L]{tree: tree, declarationPrefix: f.DeclarationPrefix, descriptorPrefix: f.DescriptorPrefix}
	d := &Decoded[L]{frontend: f, tree: tree, modules: make(map[decl.Digest]ast.NodeID), declarations: make(map[decl.DeclarationID]ast.NodeID), parameters: make(map[decl.ParameterID]ast.NodeID), imports: make(map[importID]ast.NodeID)}
	root := tree.Nodes[tree.Root-1].References[0].Target.Node
	for _, child := range tree.Nodes[tree.Root-1].Children {
		value, err := r.module(child.Node)
		if err != nil {
			var syntax *text.Diagnostic
			if child.Node != root && errors.As(err, &syntax) && syntax.Code != text.LimitExceeded {
				id := child.Node
				if syntax.Node != nil {
					id = syntax.Node.Node
				}
				return nil, d.assertion(id, true, err.Error())
			}
			return nil, err
		}
		pin, err := f.PinModel(value, limits)
		if err != nil {
			return nil, err
		}
		if r.attr(child.Node, s.AttributePin) != pin.String() {
			return nil, d.assertion(child.Node, true, "module pin does not match its declaration bytes")
		}
		if d.modules[pin] != 0 {
			return nil, d.assertion(child.Node, true, "duplicate module document in input")
		}
		d.indexModule(child.Node, pin, r)
		if child.Node == root {
			d.root = value
		} else {
			d.dependencies = append(d.dependencies, value)
		}
	}
	return d, nil
}

// Resolve returns a complete checked model only after every present annotation
// and every type/literal position agrees with Horos's authoritative projection.
func (d *Decoded[L]) Resolve(limits decl.Limits) (Program, error) {
	program, err := d.frontend.ResolveModel(d.root, d.dependencies, limits)
	if err != nil {
		return nil, err
	}
	expected, err := projectChecked[L](program, d.tree.ID)
	if err != nil {
		return nil, err
	}
	if err := d.correspond(expected); err != nil {
		return nil, err
	}
	return program, nil
}

// Read is an explicitly semantic convenience API. Pack/Check no longer calls it.
func (f Frontend[L]) Read(view ast.View[L]) (Program, error) {
	d, err := f.Decode(view, decl.DefaultLimits())
	if err != nil {
		return nil, err
	}
	return d.Resolve(decl.DefaultLimits())
}

func (d *Decoded[L]) DeclarationSite(id decl.DeclarationID) ast.NodeRef {
	return d.site(d.declarations[id])
}
func (d *Decoded[L]) ParameterSite(id decl.ParameterID) ast.NodeRef { return d.site(d.parameters[id]) }
func (d *Decoded[L]) ModuleSite(id decl.Digest) ast.NodeRef         { return d.site(d.modules[id]) }
func (d *Decoded[L]) BodySite(id decl.DeclarationID) ast.NodeRef {
	node := d.declarations[id]
	if node != 0 {
		for _, child := range d.tree.Nodes[node-1].Children {
			if child.Role == s.RoleBody {
				return d.site(child.Node)
			}
		}
	}
	return ast.NodeRef{}
}

// OperandSite follows a canonical tuple path through the existing AST encoding.
// It interprets only declaration syntax framing, never descriptor constructors.
func (d *Decoded[L]) OperandSite(site ast.NodeRef, path []int) ast.NodeRef {
	if site.Tree != d.tree.ID || site.Node == 0 || int(site.Node) > len(d.tree.Nodes) {
		return ast.NodeRef{}
	}
	for len(path) > 0 {
		node := d.tree.Nodes[site.Node-1]
		index := path[0]
		path = path[1:]
		switch node.Kind {
		case s.KindDescriptor:
			index-- // descriptor tag has no separate AST node
		case s.KindApplication:
			if index != 2 || len(path) == 0 {
				return ast.NodeRef{}
			}
			index, path = path[0]-1, path[1:] // encoded argument list tag
		case s.KindTuple:
		default:
			return ast.NodeRef{}
		}
		if index < 0 || index >= len(node.Children) {
			return ast.NodeRef{}
		}
		site.Node = node.Children[index].Node
	}
	return site
}

func (d *Decoded[L]) ArgumentSite(site ast.NodeRef, index int) ast.NodeRef {
	return d.OperandSite(site, []int{2, index + 1})
}

func (d *Decoded[L]) Locate(detail *decl.Error) *ast.NodeRef {
	if detail == nil {
		return nil
	}
	var site ast.NodeRef
	if detail.Declaration != "" {
		id := decl.DeclarationID{Module: detail.Module, Name: detail.Declaration}
		site = d.BodySite(id)
		if len(detail.Path) > 0 {
			if exact := d.OperandSite(site, detail.Path); exact.Node != 0 {
				site = exact
			}
		}
	} else if detail.ImportAlias != "" {
		site = d.site(d.imports[importID{module: detail.Module, alias: detail.ImportAlias}])
	} else {
		site = d.ModuleSite(detail.Module)
	}
	if site.Node == 0 {
		return nil
	}
	return &site
}

func (d *Decoded[L]) correspond(expected ast.AST[L]) error {
	mapping := map[ast.NodeID]ast.NodeID{expected.Root: d.tree.Root}
	contexts := map[ast.NodeID]bool{}
	rootModule := d.tree.Nodes[d.tree.Root-1].References[0].Target.Node
	wantReader := reader[L]{tree: expected}
	gotReader := reader[L]{tree: d.tree}
	var walk func(ast.NodeID, ast.NodeID, bool) error
	walk = func(wantID, gotID ast.NodeID, context bool) error {
		want, got := expected.Nodes[wantID-1], d.tree.Nodes[gotID-1]
		bad := func() error {
			return d.annotation(gotID, context, "Horos syntax does not match checked declaration positions")
		}
		if want.Kind != got.Kind || len(want.Attributes) != len(got.Attributes) || len(want.Children) != len(got.Children) {
			return bad()
		}
		for _, attribute := range want.Attributes {
			if gotReader.attr(gotID, attribute.Name) != attribute.Text {
				return bad()
			}
		}
		mapping[wantID] = gotID
		contexts[wantID] = context
		type memberKey struct{ role, name string }
		var members map[memberKey]ast.Child
		if want.Kind == s.KindModule {
			members = make(map[memberKey]ast.Child, len(got.Children))
			for _, candidate := range got.Children {
				name := s.AttributeName
				if candidate.Role == s.RoleImport {
					name = s.AttributeAlias
				}
				members[memberKey{candidate.Role, gotReader.attr(candidate.Node, name)}] = candidate
			}
		}
		for i, child := range want.Children {
			actual := got.Children[i]
			if want.Kind == s.KindModule {
				name := s.AttributeName
				if child.Role == s.RoleImport {
					name = s.AttributeAlias
				}
				var found bool
				actual, found = members[memberKey{child.Role, wantReader.attr(child.Node, name)}]
				if !found {
					return bad()
				}
			}
			if child.Role != actual.Role {
				return bad()
			}
			if err := walk(child.Node, actual.Node, context); err != nil {
				return err
			}
		}
		return nil
	}
	for _, child := range expected.Nodes[expected.Root-1].Children {
		bytes, _ := hex.DecodeString(wantReader.attr(child.Node, s.AttributePin))
		var pin decl.Digest
		copy(pin[:], bytes)
		actual := d.modules[pin]
		if actual == 0 {
			return fmt.Errorf("%w: resolved module has no source", ErrInvalid)
		}
		if err := walk(child.Node, actual, actual != rootModule); err != nil {
			return err
		}
	}
	for wantID, gotID := range mapping {
		want := expected.Nodes[wantID-1]
		for _, actual := range d.tree.Nodes[gotID-1].References {
			found := false
			for _, reference := range want.References {
				if reference.Role == actual.Role && actual.Target.Tree == d.tree.ID && actual.Target.Node == mapping[reference.Target.Node] {
					found = true
					break
				}
			}
			if !found {
				return d.annotation(gotID, contexts[wantID], "incorrect declaration, import or parameter annotation")
			}
		}
	}
	return nil
}
