// Package target constructs native ASTs from a closed native semantic plan.
package target

import (
	"context"
	"slices"
	"strings"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	native "github.com/Bitspark/bittype/internal/cpp"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	model "github.com/Bitspark/bittype/transform/syntax/cpp/model/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type Provider struct{ Namespace string }
type Backend struct {
	*m.Store[s.Language]
	config  Provider
	headers []string
}

func (p Provider) CheckExpression(ctx context.Context, e g.Expression[s.Language]) error {
	return model.CheckExpression(ctx, e)
}
func (p Provider) CheckDeclaration(ctx context.Context, d g.TargetDeclaration[s.Language]) error {
	return model.CheckDeclaration(ctx, d)
}
func (p Provider) Build(ctx context.Context, in g.BuildInput[s.Language]) (g.Assembly, error) {
	if p.Namespace == "" || !native.Identifier(p.Namespace) || strings.HasPrefix(p.Namespace, "_") || p.Namespace == "std" {
		return g.Assembly{}, g.ErrInvalid
	}
	b := &Backend{Store: m.New[s.Language](ctx, in.Limits), config: p}
	return m.Compile(ctx, in, b.Store, b)
}
func (b *Backend) Parameter(name string, origin ast.NodeRef) (m.Fragment[s.Language], error) {
	return b.Node(m.Node(s.KindTypeParameter, m.Attributes("name", name), origin))
}
func (b *Backend) Reference(symbol m.Symbol[s.Language], args []m.Fragment[s.Language], origin ast.NodeRef) (m.Fragment[s.Language], error) {
	var segments []m.Child[s.Language]
	var attrs []ast.Attribute
	if !symbol.IsFormal() {
		ns, err := b.Node(m.Node(s.KindNameSegment, m.Attributes("name", b.config.Namespace), origin))
		if err != nil {
			return ns, err
		}
		segments = append(segments, m.Child[s.Language]{Role: s.RoleSegment, Fragment: ns})
		attrs = m.Attributes("global", "true")
	}
	segment, err := b.Node(m.Node(s.KindNameSegment, m.Attributes("name", symbol.Name()), origin), edges(s.RoleArgument, args)...)
	if err != nil {
		return segment, err
	}
	segment, err = b.Link(segment, symbol, s.RoleTarget)
	if err != nil {
		return segment, err
	}
	segments = append(segments, m.Child[s.Language]{Role: s.RoleSegment, Fragment: segment})
	return b.Node(m.Node(s.KindNamedType, attrs, origin), segments...)
}
func edges(role string, fragments []m.Fragment[s.Language]) []m.Child[s.Language] {
	out := make([]m.Child[s.Language], len(fragments))
	for i, f := range fragments {
		out[i] = m.Child[s.Language]{Role: role, Fragment: f}
	}
	return out
}
func attribute(n s.Node, key string) string {
	for _, a := range n.Attributes {
		if a.Name == key {
			return a.Text
		}
	}
	return ""
}
func (b *Backend) Include(header string) error {
	if header == "" || strings.ContainsAny(header, "\r\n<>") {
		return g.ErrInvalid
	}
	for _, old := range b.headers {
		if old == header {
			return nil
		}
	}
	b.headers = append(b.headers, header)
	return nil
}

// Path is native syntax; selecting a library artifact and its meaning is the caller's responsibility.
func (b *Backend) Path(names []string, args []m.Fragment[s.Language], origin ast.NodeRef, absolute bool) (m.Fragment[s.Language], error) {
	if len(names) == 0 {
		return m.Fragment[s.Language]{}, g.ErrInvalid
	}
	var segments []m.Child[s.Language]
	for i, name := range names {
		var children []m.Child[s.Language]
		if i == len(names)-1 {
			children = edges(s.RoleArgument, args)
		}
		p, err := b.Node(m.Node(s.KindNameSegment, m.Attributes("name", name), origin), children...)
		if err != nil {
			return p, err
		}
		segments = append(segments, m.Child[s.Language]{Role: s.RoleSegment, Fragment: p})
	}
	var attrs []ast.Attribute
	if absolute {
		attrs = m.Attributes("global", "true")
	}
	return b.Node(m.Node(s.KindNamedType, attrs, origin), segments...)
}
func (b *Backend) Declaration(d g.TargetDeclaration[s.Language], params []m.Fragment[s.Language], body m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	name, origin := d.Name, d.Origin
	form, _ := model.Name(d.Form)
	record := form != model.FormAlias

	children := edges(s.RoleTemplateParameter, params)
	if record {
		n, fields, err := b.Inspect(body)
		if err != nil {
			return body, err
		}
		if n.Kind != s.KindRecord {
			return body, g.ErrInvalid
		}
		for _, f := range fields {
			fn, _, err := b.Inspect(f.Fragment)
			if err != nil {
				return body, err
			}
			for _, p := range params {
				pn, _, err := b.Inspect(p)
				if err != nil {
					return body, err
				}
				if attribute(fn, "name") == attribute(pn, "name") {
					return body, g.ErrUnsupported
				}
			}
		}
		children = append(children, fields...)
		return b.Node(m.Node(s.KindRecord, m.Attributes("name", name, "record-key", "struct"), origin), children...)
	}
	children = append(children, m.Child[s.Language]{Role: s.RoleType, Fragment: body})
	return b.Node(m.Node(s.KindTypeAlias, m.Attributes("name", name), origin), children...)
}
func (b *Backend) Finish(declarations []m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	var children []m.Child[s.Language]
	headers := slices.Clone(b.headers)
	slices.Sort(headers)
	for _, header := range headers {
		n, err := b.Node(m.Node(s.KindInclude, m.Attributes("path", header, "system", "true"), ast.NodeRef{}))
		if err != nil {
			return n, err
		}
		children = append(children, m.Child[s.Language]{Role: s.RoleInclude, Fragment: n})
	}
	ns, err := b.Node(m.Node(s.KindNamespace, m.Attributes("name", b.config.Namespace), ast.NodeRef{}), edges(s.RoleDeclaration, declarations)...)
	if err != nil {
		return ns, err
	}
	children = append(children, m.Child[s.Language]{Role: s.RoleDeclaration, Fragment: ns})
	return b.Node(m.Node(s.KindTranslationUnit, nil, ast.NodeRef{}), children...)
}
func (b *Backend) Expression(e g.Expression[s.Language], args []m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	if err := model.CheckExpression(context.Background(), e); err != nil {
		return m.Fragment[s.Language]{}, err
	}
	name, _ := model.Name(e.Form())
	data, origin := e.Data(), e.Origin()
	switch name {
	case model.FormPrimitive:
		if data[0] == model.PrimitiveString {
			if err := b.Include("string"); err != nil {
				return m.Fragment[s.Language]{}, err
			}
			return b.Path([]string{"std", "string"}, nil, origin, true)
		}
		return b.Node(m.Node(s.KindBuiltinType, m.Attributes("name", "bool"), origin))
	case model.FormSequence:
		if err := b.Include("vector"); err != nil {
			return m.Fragment[s.Language]{}, err
		}
		return b.Path([]string{"std", "vector"}, args, origin, true)
	case model.FormField:
		return b.Node(m.Node(s.KindField, m.Attributes("name", data[0]), origin), m.Child[s.Language]{Role: s.RoleType, Fragment: args[0]})
	case model.FormRecord:
		return b.Node(m.Node(s.KindRecord, nil, origin), edges(s.RoleMember, args)...)
	case model.FormExternal:
		if data[0] == "std" && data[1] == "optional" {
			if err := b.Include("optional"); err != nil {
				return m.Fragment[s.Language]{}, err
			}
		}
		return b.Path(append(strings.Split(data[0], "::"), data[1]), args, origin, true)
	}
	return m.Fragment[s.Language]{}, g.ErrUnsupported
}

func (p Provider) CheckProgram(ctx context.Context, program g.Program[s.Language]) error {
	return m.CheckSingleScope(ctx, program)
}

// CheckMember states the member vocabulary of the target (bittype#196, row
// 13): generation.MemberVocabulary. A member of a record is one field form
// that this target admits. The rule is the shared one, selected here with
// this target's member form and its own expression rule.
func (p Provider) CheckMember(ctx context.Context, e g.Expression[s.Language]) error {
	return m.CheckMemberForm(ctx, e, model.Form(model.FormField), p.CheckExpression)
}
