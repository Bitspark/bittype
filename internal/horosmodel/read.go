package horosmodel

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"

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
func (r reader[L]) index(id ast.NodeID, name string) (core.Value, error) {
	spelling := r.attr(id, name)
	if len(spelling) > 20 {
		return nil, r.bad(id, "invalid lexical index")
	}
	index, err := strconv.Atoi(spelling)
	if err != nil || index < 0 {
		return nil, r.bad(id, "invalid lexical index")
	}
	return data.EncodeInt(big.NewInt(int64(index))), nil
}

func (r reader[L]) value(id ast.NodeID) (core.Value, error) {
	n := r.node(id)
	if s.Kind(n.Kind) == s.KindAtom {
		bytes, err := hex.DecodeString(r.attr(id, s.AttributeHex))
		if err != nil {
			return nil, r.bad(id, "invalid atom hex")
		}
		return core.NewAtom(bytes), nil
	}
	if s.Kind(n.Kind) == s.KindParameter || s.Kind(n.Kind) == s.KindKindParameterReference {
		index, err := r.index(id, s.AttributeIndex)
		if err != nil {
			return nil, err
		}
		tag := "param"
		if s.Kind(n.Kind) == s.KindKindParameterReference {
			tag = "kind-param"
		}
		return tagged(r.declarationPrefix, tag, index), nil
	}
	if s.Kind(n.Kind) == s.KindRecursiveVariable {
		depth, err := r.index(id, s.AttributeDepth)
		if err != nil {
			return nil, err
		}
		index, err := r.index(id, s.AttributeIndex)
		if err != nil {
			return nil, err
		}
		return tagged(r.declarationPrefix, "recursive-var", depth, index), nil
	}
	if s.Kind(n.Kind) == s.KindReference || s.Kind(n.Kind) == s.KindKindReference {
		prefix := ""
		if s.Kind(n.Kind) == s.KindKindReference {
			prefix = "kind-"
		}
		switch r.attr(id, s.AttributeScope) {
		case "local":
			return tagged(r.declarationPrefix, prefix+"local", data.MustEncodeText(r.attr(id, s.AttributeName))), nil
		case "imported":
			return tagged(r.declarationPrefix, prefix+"imported", data.MustEncodeText(r.attr(id, s.AttributeAlias)), data.MustEncodeText(r.attr(id, s.AttributeName))), nil
		default:
			return nil, r.bad(id, "invalid reference scope")
		}
	}
	// Selector and equation nodes are AST framing, not extra canonical forms.
	if s.Kind(n.Kind) == s.KindNodePattern {
		var selectors [][2]core.Value
		var pattern core.Value
		for _, child := range n.Children {
			if child.Role == s.RoleSelector {
				selectors = append(selectors, [2]core.Value{data.MustEncodeText(r.attr(child.Node, s.AttributeName)), data.MustEncodeText(r.attr(child.Node, s.AttributeValue))})
			} else {
				value, err := r.value(child.Node)
				if err != nil {
					return nil, err
				}
				pattern = value
			}
		}
		values, err := data.EncodeMap(selectors)
		if err != nil {
			return nil, r.bad(id, "duplicate or malformed selectors")
		}
		if pattern == nil {
			return nil, r.bad(id, "missing child sequence")
		}
		return tagged(r.declarationPrefix, "node", data.MustEncodeText(r.attr(id, s.AttributeConstructor)), values, pattern), nil
	}
	if s.Kind(n.Kind) == s.KindRecursiveKind {
		var equations []core.Value
		for _, child := range n.Children {
			wrapper := r.node(child.Node)
			if len(wrapper.Children) != 1 {
				return nil, r.bad(child.Node, "missing recursive equation")
			}
			value, err := r.value(wrapper.Children[0].Node)
			if err != nil {
				return nil, err
			}
			equations = append(equations, value)
		}
		entry, err := r.index(id, s.AttributeIndex)
		if err != nil {
			return nil, err
		}
		return tagged(r.declarationPrefix, "nu", data.EncodeList(equations), entry), nil
	}
	values := make([]core.Value, len(n.Children))
	for i, child := range n.Children {
		value, err := r.value(child.Node)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}
	switch s.Kind(n.Kind) {
	case s.KindTuple:
		return core.NewTuple(values...), nil
	case s.KindDescriptor:
		return tagged(r.descriptorPrefix, r.attr(id, s.AttributeConstructor), values...), nil
	case s.KindProperSort:
		return tagged(r.declarationPrefix, "sort-proper"), nil
	case s.KindArrowSort, s.KindArrow:
		if len(values) < 2 {
			return nil, r.bad(id, "arrow requires inputs and result")
		}
		tag := "arrow"
		if s.Kind(n.Kind) == s.KindArrowSort {
			tag = "sort-arrow"
		}
		return tagged(r.declarationPrefix, tag, data.EncodeList(values[:len(values)-1]), values[len(values)-1]), nil
	case s.KindTop, s.KindBottom:
		tag := "top"
		if s.Kind(n.Kind) == s.KindBottom {
			tag = "bottom"
		}
		return tagged(r.declarationPrefix, tag, values...), nil
	case s.KindApplication, s.KindKindApplication:
		if len(values) < 2 {
			return nil, r.bad(id, "application requires callee and arguments")
		}
		tag := "apply"
		if s.Kind(n.Kind) == s.KindKindApplication {
			tag = "kind-apply"
		}
		return tagged(r.declarationPrefix, tag, values[0], data.EncodeList(values[1:])), nil
	case s.KindUnion, s.KindIntersection:
		if len(values) < 3 {
			return nil, r.bad(id, "kind Boolean form requires sort and operands")
		}
		tag := "or"
		if s.Kind(n.Kind) == s.KindIntersection {
			tag = "and"
		}
		return tagged(r.declarationPrefix, tag, values[0], data.EncodeList(values[1:])), nil
	case s.KindSequenceEpsilon:
		return tagged(r.declarationPrefix, "epsilon"), nil
	case s.KindSequenceItem:
		return tagged(r.declarationPrefix, "item", values...), nil
	case s.KindSequenceConcat, s.KindSequenceAlternative:
		tag := "concat"
		if s.Kind(n.Kind) == s.KindSequenceAlternative {
			tag = "alternative"
		}
		return tagged(r.declarationPrefix, tag, data.EncodeList(values)), nil
	case s.KindSequenceStar:
		return tagged(r.declarationPrefix, "star", values...), nil
	default:
		return nil, r.bad(id, "expected a revision-4 expression or value fragment")
	}
}

func (r reader[L]) declaration(id ast.NodeID) (core.Value, error) {
	n := r.node(id)
	tag, parameterTag := "type", "parameter"
	if s.Kind(n.Kind) == s.KindKindDefinition {
		tag, parameterTag = "kind", "kind-parameter"
	} else if s.Kind(n.Kind) != s.KindTypeDefinition {
		return nil, r.bad(id, "expected declaration")
	}
	var parameters []core.Value
	var result, body core.Value
	for _, child := range n.Children {
		switch child.Role {
		case s.RoleParameter:
			parameter := r.node(child.Node)
			if len(parameter.Children) != 1 {
				return nil, r.bad(child.Node, "expected an explicit parameter contract")
			}
			contract, err := r.value(parameter.Children[0].Node)
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, tagged(r.declarationPrefix, parameterTag, data.MustEncodeText(r.attr(child.Node, s.AttributeName)), contract))
		case s.RoleResult:
			value, err := r.value(child.Node)
			if err != nil {
				return nil, err
			}
			result = value
		case s.RoleBody:
			value, err := r.value(child.Node)
			if err != nil {
				return nil, err
			}
			body = value
		default:
			return nil, r.bad(child.Node, "unknown declaration child role")
		}
	}
	if result == nil || body == nil {
		return nil, r.bad(id, "missing declaration result or body")
	}
	return tagged(r.declarationPrefix, tag, data.EncodeList(parameters), result, body), nil
}

func (r reader[L]) module(id ast.NodeID) (core.Value, error) {
	if s.Kind(r.node(id).Kind) != s.KindModule {
		return nil, r.bad(id, "expected module")
	}
	var imports, declarations [][2]core.Value
	for _, child := range r.node(id).Children {
		switch child.Role {
		case s.RoleImport:
			if s.Kind(r.node(child.Node).Kind) != s.KindImport {
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
	return tagged(r.declarationPrefix, "module", data.EncodeInt(big.NewInt(4)),
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
	parts        []*Decoded[L]
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

func (d *Decoded[L]) partForModule(pin decl.Digest) *Decoded[L] {
	if len(d.parts) == 0 {
		return d
	}
	for _, part := range d.parts {
		if part.modules[pin] != 0 {
			return part
		}
	}
	return nil
}
func (d *Decoded[L]) partForTree(id ast.TreeID) *Decoded[L] {
	if len(d.parts) == 0 {
		if d.tree.ID == id {
			return d
		}
		return nil
	}
	for _, part := range d.parts {
		if part.tree.ID == id {
			return part
		}
	}
	return nil
}
func (d *Decoded[L]) DeclarationSite(id decl.DeclarationID) ast.NodeRef {
	part := d.partForModule(id.Module)
	if part == nil {
		return ast.NodeRef{}
	}
	return part.site(part.declarations[id])
}
func (d *Decoded[L]) ParameterSite(id decl.ParameterID) ast.NodeRef {
	part := d.partForModule(id.Declaration.Module)
	if part == nil {
		return ast.NodeRef{}
	}
	return part.site(part.parameters[id])
}
func (d *Decoded[L]) ModuleSite(id decl.Digest) ast.NodeRef {
	part := d.partForModule(id)
	if part == nil {
		return ast.NodeRef{}
	}
	return part.site(part.modules[id])
}
func (d *Decoded[L]) childSite(site ast.NodeRef, role string) ast.NodeRef {
	part := d.partForTree(site.Tree)
	if part == nil || site.Node == 0 || int(site.Node) > len(part.tree.Nodes) {
		return ast.NodeRef{}
	}
	for _, child := range part.tree.Nodes[site.Node-1].Children {
		if child.Role == role {
			return part.site(child.Node)
		}
	}
	return ast.NodeRef{}
}
func (d *Decoded[L]) BodySite(id decl.DeclarationID) ast.NodeRef {
	return d.childSite(d.DeclarationSite(id), s.RoleBody)
}
func (d *Decoded[L]) ResultSite(id decl.DeclarationID) ast.NodeRef {
	return d.childSite(d.DeclarationSite(id), s.RoleResult)
}
func (d *Decoded[L]) ParameterContractSite(id decl.ParameterID) ast.NodeRef {
	site := d.ParameterSite(id)
	if contract := d.childSite(site, s.RoleContract); contract.Node != 0 {
		return contract
	}
	return d.childSite(site, s.RoleSort)
}

// OperandSite follows canonical framing through an existing tree-qualified AST.
// It never interprets a descriptor constructor or traverses opaque literals as types.
func (d *Decoded[L]) OperandSite(site ast.NodeRef, path []int) ast.NodeRef {
	part := d.partForTree(site.Tree)
	if part == nil || site.Node == 0 || int(site.Node) > len(part.tree.Nodes) {
		return ast.NodeRef{}
	}
	for len(path) > 0 {
		node := part.tree.Nodes[site.Node-1]
		index := path[0]
		path = path[1:]
		listIndex := func() int {
			if len(path) == 0 || path[0] <= 0 {
				return -1
			}
			item := path[0] - 1
			path = path[1:]
			return item
		}
		equation := false
		switch s.Kind(node.Kind) {
		case s.KindTuple:
		case s.KindDescriptor:
			index--
		case s.KindApplication, s.KindKindApplication:
			if index == 1 {
				index = 0
			} else if index == 2 {
				item := listIndex()
				if item < 0 {
					return ast.NodeRef{}
				}
				index = item + 1
			} else {
				return ast.NodeRef{}
			}
		case s.KindArrowSort, s.KindArrow:
			if index == 1 {
				index = listIndex()
				if index >= len(node.Children)-1 {
					return ast.NodeRef{}
				}
			} else if index == 2 {
				index = len(node.Children) - 1
			} else {
				return ast.NodeRef{}
			}
		case s.KindTop, s.KindBottom, s.KindSequenceItem, s.KindSequenceStar:
			if index != 1 {
				return ast.NodeRef{}
			}
			index = 0
		case s.KindUnion, s.KindIntersection:
			if index == 1 {
				index = 0
			} else if index == 2 {
				item := listIndex()
				if item < 0 {
					return ast.NodeRef{}
				}
				index = item + 1
			} else {
				return ast.NodeRef{}
			}
		case s.KindNodePattern:
			if index != 3 {
				return ast.NodeRef{}
			}
			index = len(node.Children) - 1
		case s.KindRecursiveKind:
			if index != 1 {
				return ast.NodeRef{}
			}
			index = listIndex()
			equation = true
		case s.KindSequenceConcat, s.KindSequenceAlternative:
			if index != 1 {
				return ast.NodeRef{}
			}
			index = listIndex()
		default:
			return ast.NodeRef{}
		}
		if index < 0 || index >= len(node.Children) {
			return ast.NodeRef{}
		}
		site.Node = node.Children[index].Node
		if equation {
			site = part.childSite(site, s.RoleBody)
			if site.Node == 0 {
				return ast.NodeRef{}
			}
		}
	}
	return site
}
func (d *Decoded[L]) CalleeSite(site ast.NodeRef) ast.NodeRef { return d.OperandSite(site, []int{1}) }
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
		switch detail.Region {
		case "body", "kind-body":
			site = d.BodySite(id)
		case "result", "result-sort":
			site = d.ResultSite(id)
		case "parameter-kind", "kind-parameter-sort":
			site = d.ParameterContractSite(decl.ParameterID{Declaration: id, Index: detail.ParameterIndex})
		default:
			site = d.DeclarationSite(id)
		}
		if len(detail.Path) > 0 {
			if exact := d.OperandSite(site, detail.Path); exact.Node != 0 {
				site = exact
			}
		}
	} else if detail.ImportAlias != "" {
		part := d.partForModule(detail.Module)
		if part != nil {
			site = part.site(part.imports[importID{module: detail.Module, alias: detail.ImportAlias}])
		}
	} else {
		site = d.ModuleSite(detail.Module)
	}
	if site.Node == 0 {
		return nil
	}
	return &site
}

func (d *Decoded[L]) CalleeDeclarationSite(detail *decl.Error) *ast.NodeRef {
	if detail == nil {
		return nil
	}
	var site ast.NodeRef
	if detail.CalleeParameter.Declaration.Name != "" {
		site = d.ParameterSite(detail.CalleeParameter)
	} else if detail.Target.Name != "" {
		site = d.DeclarationSite(detail.Target)
	}
	if site.Node == 0 {
		return nil
	}
	return &site
}

// correspond maps canonical positions to tree-qualified sites. Files stay separate.
func (d *Decoded[L]) correspond(expected ast.AST[L]) error {
	mapping := map[ast.NodeID]ast.NodeRef{expected.Root: {Tree: d.tree.ID, Node: d.tree.Root}}
	contexts := map[ast.NodeID]bool{}
	rootSite := d.site(d.tree.Nodes[d.tree.Root-1].References[0].Target.Node)
	wantReader := reader[L]{tree: expected}
	var walk func(ast.NodeID, ast.NodeRef, bool) error
	walk = func(wantID ast.NodeID, gotSite ast.NodeRef, context bool) error {
		part := d.partForTree(gotSite.Tree)
		if part == nil {
			return ErrInvalid
		}
		gotID := gotSite.Node
		want, got := expected.Nodes[wantID-1], part.tree.Nodes[gotID-1]
		gotReader := reader[L]{tree: part.tree}
		bad := func() error {
			return part.annotation(gotID, context, "Horos syntax does not match checked declaration positions")
		}
		gotAttributes := len(got.Attributes)
		if s.Kind(got.Kind) == s.KindEquation {
			for _, attribute := range got.Attributes {
				if attribute.Name == s.AttributeName {
					gotAttributes--
				}
			}
		}
		if want.Kind != got.Kind || len(want.Attributes) != gotAttributes || len(want.Children) != len(got.Children) {
			return bad()
		}
		for _, attribute := range want.Attributes {
			if gotReader.attr(gotID, attribute.Name) != attribute.Text {
				return bad()
			}
		}
		mapping[wantID] = gotSite
		contexts[wantID] = context
		type memberKey struct{ role, name string }
		members := map[memberKey]ast.Child{}
		if s.Kind(want.Kind) == s.KindModule || s.Kind(want.Kind) == s.KindNodePattern {
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
			if s.Kind(want.Kind) == s.KindModule || s.Kind(want.Kind) == s.KindNodePattern {
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
			if err := walk(child.Node, part.site(actual.Node), context); err != nil {
				return err
			}
		}
		return nil
	}
	for _, child := range expected.Nodes[expected.Root-1].Children {
		bytes, _ := hex.DecodeString(wantReader.attr(child.Node, s.AttributePin))
		var pin decl.Digest
		copy(pin[:], bytes)
		actual := d.ModuleSite(pin)
		if actual.Node == 0 {
			// The revision's bundled module is authoritative context, not fabricated source.
			if pin == decl.PreludePin() {
				continue
			}
			return fmt.Errorf("%w: resolved module has no source", ErrInvalid)
		}
		if err := walk(child.Node, actual, actual != rootSite); err != nil {
			return err
		}
	}
	for wantID, gotSite := range mapping {
		part := d.partForTree(gotSite.Tree)
		for _, actual := range part.tree.Nodes[gotSite.Node-1].References {
			found := false
			for _, reference := range expected.Nodes[wantID-1].References {
				if reference.Role == actual.Role && actual.Target == mapping[reference.Target.Node] {
					found = true
					break
				}
			}
			if !found {
				return part.annotation(gotSite.Node, contexts[wantID], "incorrect declaration, import or parameter annotation")
			}
		}
	}
	return nil
}

// DecodeCollection pools canonical module documents, never concatenated syntax.
func (f Frontend[L]) DecodeCollection(inputs []ast.AST[L], entry ast.NodeRef, limits decl.Limits) (*Decoded[L], error) {
	if len(inputs) == 0 {
		return nil, ErrInvalid
	}
	combined := &Decoded[L]{frontend: f}
	pins := map[decl.Digest]bool{}
	for _, input := range inputs {
		part, err := f.Decode(input, limits)
		if err != nil {
			return nil, err
		}
		if entry.Tree == input.ID {
			if entry.Node != input.Root {
				return nil, ErrInvalid
			}
			combined.tree = input
			combined.root = part.root
		}
		for pin, node := range part.modules {
			if pins[pin] {
				return nil, part.assertion(node, true, "duplicate module document across input files")
			}
			pins[pin] = true
		}
		combined.parts = append(combined.parts, part)
	}
	if combined.root == nil {
		return nil, ErrInvalid
	}
	for _, part := range combined.parts {
		if part.tree.ID != entry.Tree {
			combined.dependencies = append(combined.dependencies, part.root)
		}
		combined.dependencies = append(combined.dependencies, part.dependencies...)
	}
	return combined, nil
}
