// Package golangsem owns the private compiler bridge for the captured Go fragment.
package golangsem

import (
	"context"
	"fmt"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	goast "go/ast"
	"go/token"
	"go/types"
	"slices"
	"sort"
	"strconv"
)

// Failure distinguishes unsupported/limited coverage from an actual type error.
type Failure struct {
	Err     error
	At      ast.NodeRef
	Invalid bool
}

func (e *Failure) Error() string                 { return e.Err.Error() }
func (e *Failure) Unwrap() error                 { return e.Err }
func (c *converter) pos(id ast.NodeID) token.Pos { return c.positions.Pos(int(id-1) * 1024) }
func setPosition(e goast.Expr, p token.Pos) {
	switch x := e.(type) {
	case *goast.ArrayType:
		x.Lbrack = p
	case *goast.StructType:
		x.Struct = p
	case *goast.InterfaceType:
		x.Interface = p
	case *goast.UnaryExpr:
		x.OpPos = p
	case *goast.BinaryExpr:
		x.OpPos = p
	case *goast.IndexListExpr:
		x.Lbrack = p
	}
}

type converter struct {
	tree        s.AST
	origins     map[*goast.Ident]ast.NodeRef
	refs        map[*goast.Ident][]ast.Reference
	positions   *token.File
	expressions map[goast.Expr]ast.NodeRef
	// imports is true only under an ImportConfiguration; the base profile
	// admits neither import declarations nor qualified types.
	imports bool
	sites   []importSite
	// rhs records each declaration's right-hand side, so the exporter can tell
	// when a captured declaration's structure comes from a supplied package.
	rhs map[*goast.Ident]goast.Expr
}

// importSite is one captured import declaration and the path it names.
type importSite struct {
	At   ast.NodeRef
	Path string
}

func (c *converter) node(id ast.NodeID) s.Node { return c.tree.Nodes[id-1] }
func attribute(n s.Node, name string) string {
	for _, a := range n.Attributes {
		if a.Name == name {
			return a.Text
		}
	}
	return ""
}
func children(n s.Node, role string) []ast.NodeID {
	var r []ast.NodeID
	for _, ch := range n.Children {
		if ch.Role == role {
			r = append(r, ch.Node)
		}
	}
	return r
}
func (c *converter) ident(id ast.NodeID) *goast.Ident {
	node := c.node(id)
	x := goast.NewIdent(attribute(node, s.AttributeName))
	c.origins[x] = ast.NodeRef{Tree: c.tree.ID, Node: id}
	c.refs[x] = node.References
	x.NamePos = c.pos(id)
	return x
}
func (c *converter) one(node s.Node, role string, depth int) (goast.Expr, error) {
	x := children(node, role)
	if len(x) != 1 {
		return nil, n.ErrUnsupported
	}
	return c.expression(x[0], depth+1)
}

// Check consumes syntax-admitted trees. Token positions are synthetic addresses
// used only to return compiler diagnostics to captured nodes, never source spans.
// captured is the request's Files context. It is required only under an import
// configuration, where it assigns each tree to a package by its file path.
func Check(ctx context.Context, inputs []s.AST, captured *files.Context, settings []semantic.Setting) (*Checked, error) {
	if ctx == nil || len(inputs) == 0 {
		return nil, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	environment, imports, err := n.Imports(settings)
	if err != nil {
		return nil, semantic.ErrInvalidRequest
	}
	fset := token.NewFileSet()
	// The environment is validated completely before any captured code is
	// judged, and captured trees are then processed in file-path order, so no
	// result depends on capture order.
	var env *environmentState
	if imports {
		if env, err = prepareEnvironment(ctx, fset, environment, inputs, captured); err != nil {
			return nil, err
		}
		inputs = slices.Clone(inputs)
		sort.SliceStable(inputs, func(i, j int) bool { return env.paths[inputs[i].ID] < env.paths[inputs[j].ID] })
		if err := checkLayout(env, inputs); err != nil {
			return nil, err
		}
	}
	if err := bound(inputs); err != nil {
		return nil, err
	}
	var goFiles []*goast.File
	var sites []importSite
	origins := map[*goast.Ident]ast.NodeRef{}
	refs := map[*goast.Ident][]ast.Reference{}
	expressions := map[goast.Expr]ast.NodeRef{}
	positions := map[token.Pos]ast.NodeRef{}
	rhs := map[*goast.Ident]goast.Expr{}
	for _, tree := range inputs {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		file := fset.AddFile(string(tree.ID), -1, len(tree.Nodes)*1024+1024)
		for i := range tree.Nodes {
			positions[file.Pos(i*1024)] = ast.NodeRef{Tree: tree.ID, Node: ast.NodeID(i + 1)}
		}
		c := converter{tree: tree, origins: origins, refs: refs, positions: file, expressions: expressions, imports: imports, rhs: rhs}
		f, err := c.file()
		if err != nil {
			if _, ok := err.(*Failure); !ok {
				err = &Failure{Err: err, At: ast.NodeRef{Tree: tree.ID, Node: tree.Root}}
			}
			return nil, err
		}
		goFiles = append(goFiles, f)
		sites = append(sites, c.sites...)
	}
	info := &types.Info{Defs: map[*goast.Ident]types.Object{}, Uses: map[*goast.Ident]types.Object{}, Types: map[goast.Expr]types.TypeAndValue{}}
	var program *checkedProgram
	if imports {
		program, err = checkCaptured(ctx, fset, env, inputs, goFiles, sites, info, positions)
	} else {
		program, err = checkSingle(fset, inputs, goFiles, info, positions)
	}
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	if err != nil {
		return nil, err
	}
	defined := map[types.Object]ast.NodeRef{}
	for id, obj := range info.Defs {
		if obj != nil {
			defined[obj] = origins[id]
		}
	}
	for id, obj := range info.Uses {
		for _, ref := range refs[id] {
			if ref.Role != s.RoleTarget || ref.Target != defined[obj] {
				return nil, semantic.ErrInvalidRequest
			}
		}
	}
	return export(ctx, program, info, origins, expressions, positions, rhs)
}

func typeFailure(err error, positions map[token.Pos]ast.NodeRef) error {
	var at ast.NodeRef
	if e, ok := err.(types.Error); ok {
		at = positions[e.Pos]
	}
	return &Failure{Err: err, At: at, Invalid: true}
}

// checkSingle is the base profile: one captured package at n.PackagePath.
func checkSingle(fset *token.FileSet, inputs []s.AST, goFiles []*goast.File, info *types.Info, positions map[token.Pos]ast.NodeRef) (*checkedProgram, error) {
	conf := types.Config{GoVersion: n.LanguageVersion, Sizes: types.SizesFor("gc", "amd64")}
	pkg, err := conf.Check(n.PackagePath, fset, goFiles, info)
	if err != nil {
		return nil, typeFailure(err, positions)
	}
	captured := capturedPackage{pkg: pkg}
	for _, tree := range inputs {
		captured.roots = append(captured.roots, ast.NodeRef{Tree: tree.ID, Node: tree.Root})
	}
	return &checkedProgram{captured: []capturedPackage{captured}}, nil
}

func bound(inputs []s.AST) error {
	if len(inputs) > 16 {
		return semantic.ErrLimit
	}
	nodes, entries, size := 0, 0, 0
	for _, t := range inputs {
		nodes += len(t.Nodes)
		size += len(t.ID)
		for _, v := range t.Nodes {
			size += len(v.Kind)
			entries += len(v.Attributes) + len(v.Children) + len(v.References) + len(v.Origin.Spans) + len(v.Origin.Nodes)
			for _, a := range v.Attributes {
				size += len(a.Name) + len(a.Text)
			}
			for _, c := range v.Children {
				size += len(c.Role)
			}
			for _, r := range v.References {
				size += len(r.Role) + len(r.Target.Tree)
			}
			for _, s := range v.Origin.Spans {
				size += len(s.Source)
			}
			for _, o := range v.Origin.Nodes {
				size += len(o.Tree)
			}
		}
		if nodes > 4096 || entries > 32768 || size > 1<<20 {
			return semantic.ErrLimit
		}
	}
	return nil
}
func (c *converter) expression(id ast.NodeID, depth int) (out goast.Expr, err error) {
	defer func() {
		if err != nil {
			if _, ok := err.(*Failure); !ok {
				err = &Failure{Err: err, At: ast.NodeRef{Tree: c.tree.ID, Node: id}}
			}
			return
		}
		if out != nil {
			c.expressions[out] = ast.NodeRef{Tree: c.tree.ID, Node: id}
			setPosition(out, c.pos(id))
		}
	}()
	if depth > 64 {
		return nil, semantic.ErrLimit
	}
	node := c.node(id)
	switch node.Kind {
	case s.KindNamedType:
		return c.ident(id), nil
	case s.KindQualifiedType:
		if !c.imports {
			return nil, fmt.Errorf("%w: %s", n.ErrUnsupported, node.Kind)
		}
		// The qualifier names an import of the same file; the selected identifier
		// carries the node's origin, so the binding is the imported declaration.
		qualifier := goast.NewIdent(attribute(node, s.AttributeQualifier))
		qualifier.NamePos = c.pos(id)
		return &goast.SelectorExpr{X: qualifier, Sel: c.ident(id)}, nil
	case s.KindSliceType:
		e, err := c.one(node, s.RoleElement, depth)
		return &goast.ArrayType{Elt: e}, err
	case s.KindTypeApplication:
		f, err := c.one(node, s.RoleConstructor, depth)
		if err != nil {
			return nil, err
		}
		var args []goast.Expr
		for _, child := range children(node, s.RoleArgument) {
			x, e := c.expression(child, depth+1)
			if e != nil {
				return nil, e
			}
			args = append(args, x)
		}
		return &goast.IndexListExpr{X: f, Indices: args}, nil
	case s.KindStructType:
		fields := &goast.FieldList{}
		for _, child := range children(node, s.RoleField) {
			f := c.node(child)
			if f.Kind != s.KindField {
				return nil, n.ErrUnsupported
			}
			ty, e := c.one(f, s.RoleType, depth)
			if e != nil {
				return nil, e
			}
			field := &goast.Field{Names: []*goast.Ident{c.ident(child)}, Type: ty}
			if tag := attribute(f, s.AttributeTag); tag != "" {
				field.Tag = &goast.BasicLit{ValuePos: c.pos(child), Kind: token.STRING, Value: strconv.Quote(tag)}
			}
			fields.List = append(fields.List, field)
		}
		return &goast.StructType{Fields: fields}, nil
	case s.KindInterfaceType:
		fields := &goast.FieldList{}
		for _, child := range children(node, s.RoleMember) {
			f := c.node(child)
			if f.Kind != s.KindEmbeddedType {
				return nil, n.ErrUnsupported
			}
			ty, e := c.one(f, s.RoleType, depth)
			if e != nil {
				return nil, e
			}
			fields.List = append(fields.List, &goast.Field{Type: ty})
		}
		return &goast.InterfaceType{Methods: fields}, nil
	case s.KindTypeApproximation:
		e, err := c.one(node, s.RoleType, depth)
		return &goast.UnaryExpr{Op: token.TILDE, X: e}, err
	case s.KindTypeUnion:
		var joined goast.Expr
		for _, child := range children(node, s.RoleTerm) {
			e, err := c.expression(child, depth+1)
			if err != nil {
				return nil, err
			}
			if joined == nil {
				joined = e
			} else {
				joined = &goast.BinaryExpr{X: joined, Op: token.OR, Y: e}
			}
		}
		if joined == nil {
			return nil, n.ErrUnsupported
		}
		return joined, nil
	default:
		return nil, fmt.Errorf("%w: %s", n.ErrUnsupported, node.Kind)
	}
}
func (c *converter) file() (*goast.File, error) {
	node := c.node(c.tree.Root)
	if node.Kind != s.KindFile || !c.imports && len(children(node, s.RoleImport)) != 0 {
		return nil, n.ErrUnsupported
	}
	f := &goast.File{Name: goast.NewIdent(attribute(node, s.AttributePackage)), Package: c.pos(c.tree.Root)}
	f.Name.NamePos = c.pos(c.tree.Root)
	for _, id := range children(node, s.RoleImport) {
		spec, err := c.importSpec(id)
		if err != nil {
			return nil, &Failure{Err: err, At: ast.NodeRef{Tree: c.tree.ID, Node: id}}
		}
		f.Imports = append(f.Imports, spec)
		f.Decls = append(f.Decls, &goast.GenDecl{Tok: token.IMPORT, TokPos: c.pos(id), Specs: []goast.Spec{spec}})
	}
	for _, id := range children(node, s.RoleDeclaration) {
		d := c.node(id)
		if d.Kind != s.KindTypeDeclaration && d.Kind != s.KindTypeAlias {
			return nil, n.ErrUnsupported
		}
		ty, err := c.one(d, s.RoleType, 0)
		if err != nil {
			return nil, err
		}
		spec := &goast.TypeSpec{Name: c.ident(id), Type: ty}
		c.rhs[spec.Name] = ty
		if d.Kind == s.KindTypeAlias {
			spec.Assign = c.pos(id)
		}
		params := children(d, s.RoleTypeParameter)
		if len(params) > 16 {
			return nil, n.ErrUnsupported
		}
		if len(params) > 0 {
			spec.TypeParams = &goast.FieldList{}
		}
		for _, p := range params {
			constraint, e := c.one(c.node(p), s.RoleConstraint, 0)
			if e != nil {
				return nil, e
			}
			spec.TypeParams.List = append(spec.TypeParams.List, &goast.Field{Names: []*goast.Ident{c.ident(p)}, Type: constraint})
		}
		f.Decls = append(f.Decls, &goast.GenDecl{Tok: token.TYPE, Specs: []goast.Spec{spec}})
	}
	return f, nil
}
