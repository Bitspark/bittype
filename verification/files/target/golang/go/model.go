// Package golang explicitly checks the supported captured Go type fragment.
// It translates AST to go/ast and uses go/types; it never prints and reparses.
package golang

import (
	"context"
	"errors"
	"fmt"
	goast "go/ast"
	"go/token"
	"go/types"
	"strconv"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

const Profile = "go1.25/captured-types-1"

var ErrUnsupported = errors.New("Go type projection outside supported fragment")
var ErrForeign = errors.New("foreign or missing Go projection handle")

func operation() semantics.Key { return semantics.Key{Owner: "go", Name: "captured-type-check"} }
func validity() semantics.Key  { return semantics.Key{Owner: "go", Name: "captured-type-validity"} }
func configuration() []semantic.Setting {
	return []semantic.Setting{{Key: semantics.Key{Owner: "go", Name: "language-version"}, Value: []byte("go1.25")}}
}
func Evidence() semantic.EvidenceRequirement {
	return semantic.EvidenceRequirement{Operation: operation(), Profile: Profile, Configuration: configuration(), Reports: []semantic.Requirement{{Obligation: validity(), Verdict: semantic.Holds}}}
}

type state struct {
	snapshot semantic.Snapshot
	receipt  semantic.Result
	pkg      *types.Package
	sites    map[types.Object]ast.NodeRef
}
type Model struct{ state *state }
type Type struct {
	owner *state
	value types.Type
}
type Field struct {
	Name     string
	Type     Type
	Site     ast.NodeRef
	Embedded bool
	Tag      string
}
type Admission struct {
	owner       *state
	constructor Type
	arguments   []Type
	admitted    bool
	detail      string
}

func (m Model) Valid() bool       { return m.state != nil }
func (m Model) Same(n Model) bool { return m.Valid() && m.state == n.state }
func (m Model) Snapshot() semantic.Snapshot {
	if !m.Valid() {
		return semantic.Snapshot{}
	}
	return m.state.snapshot
}
func (m Model) Receipt() semantic.Result {
	if !m.Valid() {
		return semantic.Result{}
	}
	return m.state.receipt
}
func (m Model) Lookup(name string) (Type, error) {
	if !m.Valid() {
		return Type{}, ErrForeign
	}
	o, ok := m.state.pkg.Scope().Lookup(name).(*types.TypeName)
	if !ok {
		return Type{}, ErrForeign
	}
	return Type{m.state, o.Type()}, nil
}
func (m Model) Basic(kind types.BasicKind) (Type, error) {
	if !m.Valid() || kind != types.String && kind != types.Bool {
		return Type{}, ErrUnsupported
	}
	return Type{m.state, types.Typ[kind]}, nil
}
func (t Type) Valid() bool  { return t.owner != nil && t.value != nil }
func (t Type) Model() Model { return Model{t.owner} }
func (t Type) Same(u Type) bool {
	return t.Valid() && t.owner == u.owner && types.Identical(t.value, u.value)
}
func (t Type) IsBasic(kind types.BasicKind) bool {
	b, ok := t.value.(*types.Basic)
	return t.Valid() && ok && b.Kind() == kind
}
func (t Type) IsNamed() bool { _, ok := t.value.(*types.Named); return t.Valid() && ok }
func (t Type) IsAlias() bool { _, ok := t.value.(*types.Alias); return t.Valid() && ok }
func (t Type) Site() ast.NodeRef {
	if n, ok := t.value.(*types.Named); ok {
		return t.owner.sites[n.Obj()]
	}
	if p, ok := t.value.(*types.TypeParam); ok {
		return t.owner.sites[p.Obj()]
	}
	return ast.NodeRef{}
}
func (t Type) Parameters() []Type {
	n, ok := t.value.(*types.Named)
	if !ok || t.owner == nil {
		return nil
	}
	p := n.TypeParams()
	result := make([]Type, p.Len())
	for i := range result {
		result[i] = Type{t.owner, p.At(i)}
	}
	return result
}
func (t Type) Arguments() []Type {
	n, ok := t.value.(*types.Named)
	if !ok || t.owner == nil {
		return nil
	}
	p := n.TypeArgs()
	result := make([]Type, p.Len())
	for i := range result {
		result[i] = Type{t.owner, p.At(i)}
	}
	return result
}
func (t Type) Origin() Type {
	if n, ok := t.value.(*types.Named); ok {
		return Type{t.owner, n.Origin()}
	}
	return Type{}
}
func (t Type) SliceElement() (Type, bool) {
	n, ok := t.value.(*types.Slice)
	if !ok || t.owner == nil {
		return Type{}, false
	}
	return Type{t.owner, n.Elem()}, true
}
func (t Type) Fields() ([]Field, bool) {
	if !t.Valid() {
		return nil, false
	}
	n, ok := t.value.Underlying().(*types.Struct)
	if !ok {
		return nil, false
	}
	r := make([]Field, n.NumFields())
	for i := range r {
		f := n.Field(i)
		r[i] = Field{f.Name(), Type{t.owner, f.Type()}, t.owner.sites[f], f.Embedded(), n.Tag(i)}
	}
	return r, true
}
func (t Type) Unconstrained() bool {
	p, ok := t.value.(*types.TypeParam)
	if !ok {
		return false
	}
	i, ok := p.Constraint().Underlying().(*types.Interface)
	return ok && i.Complete().Empty()
}
func (t Type) String() string {
	if !t.Valid() {
		return ""
	}
	return types.TypeString(t.value, nil)
}

func (m Model) Admit(ctx context.Context, constructor Type, arguments ...Type) (Admission, error) {
	if ctx == nil || !m.Valid() || constructor.owner != m.state {
		return Admission{}, ErrForeign
	}
	if err := ctx.Err(); err != nil {
		return Admission{}, err
	}
	n, ok := constructor.value.(*types.Named)
	if !ok || n.TypeParams().Len() != len(arguments) || len(arguments) > 16 {
		return Admission{}, ErrUnsupported
	}
	args := make([]types.Type, len(arguments))
	for i, t := range arguments {
		if t.owner != m.state {
			return Admission{}, ErrForeign
		}
		args[i] = t.value
	}
	_, err := types.Instantiate(nil, n, args, true)
	if e := ctx.Err(); e != nil {
		return Admission{}, e
	}
	r := Admission{owner: m.state, constructor: constructor, arguments: append([]Type(nil), arguments...), admitted: err == nil}
	if err != nil {
		r.detail = err.Error()
	}
	return r, nil
}
func (a Admission) Model() Model      { return Model{a.owner} }
func (a Admission) Constructor() Type { return a.constructor }
func (a Admission) Arguments() []Type { return append([]Type(nil), a.arguments...) }
func (a Admission) Admitted() bool    { return a.owner != nil && a.admitted }
func (a Admission) Detail() string    { return a.detail }

// Inspect is an explicit native semantic step. It neither generates syntax nor
// runs during joint preparation. No imports, methods or executable bodies are
// admitted by this first profile; the environment is the supplied package plus
// the fixed Go predeclared universe.
func Inspect(ctx context.Context, snapshot semantic.Snapshot) (Model, error) {
	if ctx == nil || !snapshot.Valid() {
		return Model{}, semantic.ErrInvalidRequest
	}
	l := semantic.DefaultLimits()
	l.MaxInputs = 16
	l.MaxInputNodes = 4096
	l.MaxInputEntries = 32768
	l.MaxInputBytes = 1 << 20
	if err := snapshot.CheckLimits(l); err != nil {
		return Model{}, err
	}
	contract := semantic.Contract{Profile: Profile, Obligations: []semantic.Obligation{{Key: validity(), Domain: semantics.Key{Owner: "go", Name: "captured-package"}}}}
	var captured *state
	op, err := semantic.NewOperation[s.Language](snapshot.Language(), contract, operation(), s.Language{}, func(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
		model, e := inspect(ctx, snapshot)
		if e != nil {
			return semantic.Result{}, e
		}
		captured = model
		return request.Publish(semantic.Publication{Reports: []semantic.Report{{Obligation: validity(), Status: semantic.Complete, Verdict: semantic.Holds}}})
	})
	if err != nil {
		return Model{}, err
	}
	receipt, err := op.Run(ctx, snapshot, semantic.Options{Configuration: configuration()}, nil)
	if err != nil {
		return Model{}, err
	}
	captured.receipt = receipt
	return Model{captured}, nil
}

type converter struct {
	tree    s.AST
	origins map[*goast.Ident]ast.NodeRef
	refs    map[*goast.Ident][]ast.Reference
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
	n := c.node(id)
	x := goast.NewIdent(attribute(n, s.AttributeName))
	c.origins[x] = ast.NodeRef{Tree: c.tree.ID, Node: id}
	c.refs[x] = n.References
	return x
}
func (c *converter) one(n s.Node, role string, depth int) (goast.Expr, error) {
	x := children(n, role)
	if len(x) != 1 {
		return nil, ErrUnsupported
	}
	return c.expression(x[0], depth+1)
}
func (c *converter) expression(id ast.NodeID, depth int) (goast.Expr, error) {
	if depth > 64 {
		return nil, semantic.ErrLimit
	}
	n := c.node(id)
	switch n.Kind {
	case s.KindNamedType:
		return c.ident(id), nil
	case s.KindSliceType:
		e, err := c.one(n, s.RoleElement, depth)
		return &goast.ArrayType{Elt: e}, err
	case s.KindTypeApplication:
		f, err := c.one(n, s.RoleConstructor, depth)
		if err != nil {
			return nil, err
		}
		var args []goast.Expr
		for _, child := range children(n, s.RoleArgument) {
			x, e := c.expression(child, depth+1)
			if e != nil {
				return nil, e
			}
			args = append(args, x)
		}
		return &goast.IndexListExpr{X: f, Indices: args}, nil
	case s.KindStructType:
		fields := &goast.FieldList{}
		for _, child := range children(n, s.RoleField) {
			f := c.node(child)
			if f.Kind != s.KindField {
				return nil, ErrUnsupported
			}
			ty, e := c.one(f, s.RoleType, depth)
			if e != nil {
				return nil, e
			}
			field := &goast.Field{Names: []*goast.Ident{c.ident(child)}, Type: ty}
			if tag := attribute(f, s.AttributeTag); tag != "" {
				field.Tag = &goast.BasicLit{Kind: token.STRING, Value: strconv.Quote(tag)}
			}
			fields.List = append(fields.List, field)
		}
		return &goast.StructType{Fields: fields}, nil
	case s.KindInterfaceType:
		fields := &goast.FieldList{}
		for _, child := range children(n, s.RoleMember) {
			f := c.node(child)
			if f.Kind != s.KindEmbeddedType {
				return nil, ErrUnsupported
			}
			ty, e := c.one(f, s.RoleType, depth)
			if e != nil {
				return nil, e
			}
			fields.List = append(fields.List, &goast.Field{Type: ty})
		}
		return &goast.InterfaceType{Methods: fields}, nil
	case s.KindTypeApproximation:
		e, err := c.one(n, s.RoleType, depth)
		return &goast.UnaryExpr{Op: token.TILDE, X: e}, err
	case s.KindTypeUnion:
		var out goast.Expr
		for _, child := range children(n, s.RoleTerm) {
			e, err := c.expression(child, depth+1)
			if err != nil {
				return nil, err
			}
			if out == nil {
				out = e
			} else {
				out = &goast.BinaryExpr{X: out, Op: token.OR, Y: e}
			}
		}
		if out == nil {
			return nil, ErrUnsupported
		}
		return out, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupported, n.Kind)
	}
}
func (c *converter) file() (*goast.File, error) {
	n := c.node(c.tree.Root)
	if n.Kind != s.KindFile || len(children(n, s.RoleImport)) != 0 {
		return nil, ErrUnsupported
	}
	f := &goast.File{Name: goast.NewIdent(attribute(n, s.AttributePackage))}
	for _, id := range children(n, s.RoleDeclaration) {
		d := c.node(id)
		if d.Kind != s.KindTypeDeclaration && d.Kind != s.KindTypeAlias {
			return nil, ErrUnsupported
		}
		ty, err := c.one(d, s.RoleType, 0)
		if err != nil {
			return nil, err
		}
		spec := &goast.TypeSpec{Name: c.ident(id), Type: ty}
		if d.Kind == s.KindTypeAlias {
			spec.Assign = token.Pos(1)
		}
		params := children(d, s.RoleTypeParameter)
		if len(params) > 16 {
			return nil, ErrUnsupported
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
func inspect(ctx context.Context, snapshot semantic.Snapshot) (*state, error) {
	var files []*goast.File
	origins := map[*goast.Ident]ast.NodeRef{}
	refs := map[*goast.Ident][]ast.Reference{}
	for _, input := range snapshot.Inputs() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		tree, err := representation.Unpack[s.Language](snapshot.Language(), input)
		if err != nil {
			return nil, err
		}
		c := converter{tree, origins, refs}
		f, err := c.file()
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	info := &types.Info{Defs: map[*goast.Ident]types.Object{}, Uses: map[*goast.Ident]types.Object{}}
	conf := types.Config{GoVersion: "go1.25"}
	pkg, err := conf.Check("captured", token.NewFileSet(), files, info)
	if err != nil {
		return nil, err
	}
	sites := map[types.Object]ast.NodeRef{}
	for id, obj := range info.Defs {
		if obj != nil {
			sites[obj] = origins[id]
		}
	}
	for id, obj := range info.Uses {
		for _, ref := range refs[id] {
			if ref.Role != s.RoleTarget || ref.Target != sites[obj] {
				return nil, semantic.ErrInvalidRequest
			}
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return &state{snapshot: snapshot, pkg: pkg, sites: sites}, nil
}
