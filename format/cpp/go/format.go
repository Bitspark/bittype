package cppformat

import (
	"github.com/Bitspark/bittype/ast/go"
	check "github.com/Bitspark/bittype/check/cpp/go"
	"github.com/Bitspark/bittype/cst/go"
	"github.com/Bitspark/bittype/internal/concrete"
	"github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	"github.com/Bitspark/bittype/text/go"
	"strings"
)

type doc = concrete.Doc[s.Language]
type formatter struct {
	tree.View[s.Language]
	concrete.Builder[s.Language]
}

func Format(view ast.View[s.Language]) (cst.CST[s.Language], error) {
	t := view.ASTView()
	if err := check.Check(t); err != nil {
		return cst.CST[s.Language]{}, err
	}
	f := formatter{tree.View[s.Language]{Tree: t}, concrete.Builder[s.Language]{Vocabulary: concrete.Vocabulary[s.Language]{Layout: s.ConcreteLayout, Token: s.ConcreteToken, Whitespace: s.ConcreteWhitespace, Comment: s.ConcreteComment}}}
	var parts []doc
	for _, id := range f.Children(t.Root, "include") {
		open, close := "\"", "\""
		if f.Flag(id, "system") {
			open, close = "<", ">"
		}
		parts = append(parts, f.Group(id, f.T("#"), f.T("include"), f.Space(), f.T(open+f.Attr(id, "path")+close)), f.Line("  ", 0))
	}
	for i, id := range f.Children(t.Root, "declaration") {
		if i > 0 || len(parts) > 0 {
			parts = append(parts, f.Line("  ", 0))
		}
		parts = append(parts, f.decl(id, 0), f.Line("  ", 0))
	}
	return concrete.Build(t, f.Group(t.Root, parts...))
}
func (f formatter) templates(id ast.NodeID, level int) doc {
	var params []doc
	for _, p := range f.Children(id, "template-parameter") {
		key := f.Attr(p, "parameter-key")
		if key == "" {
			key = "typename"
		}
		parts := []doc{f.T(key), f.Space(), f.T(f.Attr(p, "name"))}
		if d := f.Child(p, "default"); d != 0 {
			parts = append(parts, f.Space(), f.T("="), f.Space(), f.typ(d))
		}
		params = append(params, f.Group(p, parts...))
	}
	if len(params) == 0 {
		return f.Seq()
	}
	return f.Seq(f.T("template"), f.T("<"), f.Join(params, f.T(","), f.Space()), f.T(">"), f.Line("  ", level))
}
func (f formatter) decl(id ast.NodeID, level int) doc {
	parts := []doc{f.templates(id, level)}
	switch f.Kind(id) {
	case "namespace":
		parts = append(parts, f.T("namespace"), f.Space(), f.T(f.Attr(id, "name")), f.Space(), f.T("{"))
		for _, d := range f.Children(id, "declaration") {
			parts = append(parts, f.Line("  ", level+1), f.decl(d, level+1))
		}
		parts = append(parts, f.Line("  ", level), f.T("}"))
	case "type-alias":
		parts = append(parts, f.T("using"), f.Space(), f.T(f.Attr(id, "name")), f.Space(), f.T("="), f.Space(), f.typ(f.Child(id, "type")), f.T(";"))
	case "record":
		parts = append(parts, f.T(f.Attr(id, "record-key")), f.Space(), f.T(f.Attr(id, "name")))
		if f.Flag(id, "forward") {
			parts = append(parts, f.T(";"))
			break
		}
		var bases []doc
		for _, b := range f.Children(id, "base") {
			p := []doc{}
			if a := f.Attr(b, "access"); a != "" {
				p = append(p, f.T(a), f.Space())
			}
			if f.Flag(b, "virtual") {
				p = append(p, f.T("virtual"), f.Space())
			}
			p = append(p, f.typ(f.Child(b, "type")))
			bases = append(bases, f.Group(b, p...))
		}
		if len(bases) > 0 {
			parts = append(parts, f.Space(), f.T(":"), f.Space(), f.Join(bases, f.T(","), f.Space()))
		}
		parts = append(parts, f.Space(), f.T("{"))
		defaultAccess := "private"
		if f.Attr(id, "record-key") == "struct" { defaultAccess = "public" }
		currentAccess := defaultAccess
		for _, m := range f.Children(id, "member") {
			parts = append(parts, f.Line("  ", level+1))
			a := f.Attr(m, "access")
			explicitAccess := a != ""
			if a == "" { a = defaultAccess }
			if explicitAccess || a != currentAccess {
				parts = append(parts, f.T(a), f.T(":"), f.Line("  ", level+1))
			}
			currentAccess = a
			switch f.Kind(m) {
			case "field":
				parts = append(parts, f.Group(m, f.render(f.Child(m, "type"), f.T(f.Attr(m, "name")), true), f.T(";")))
			case "method":
				parts = append(parts, f.method(m))
			case "destructor":
				parts = append(parts, f.destructor(m, f.Attr(id, "name")))
			default:
				parts = append(parts, f.decl(m, level+1))
			}
		}
		parts = append(parts, f.Line("  ", level), f.T("}"), f.T(";"))
	case "enum":
		parts = append(parts, f.T("enum"), f.Space())
		if f.Flag(id, "scoped") {
			parts = append(parts, f.T("class"), f.Space())
		}
		parts = append(parts, f.T(f.Attr(id, "name")))
		if u := f.Child(id, "underlying"); u != 0 {
			parts = append(parts, f.Space(), f.T(":"), f.Space(), f.typ(u))
		}
		if f.Flag(id, "forward") {
			parts = append(parts, f.T(";"))
			break
		}
		parts = append(parts, f.Space(), f.T("{"))
		for _, e := range f.Children(id, "enumerator") {
			p := []doc{f.T(f.Attr(e, "name"))}
			if value := f.Attr(e, "value"); value != "" {
				p = append(p, f.Space(), f.T("="), f.Space(), f.number(value))
			}
			parts = append(parts, f.Line("  ", level+1), f.Group(e, p...), f.T(","))
		}
		parts = append(parts, f.Line("  ", level), f.T("}"), f.T(";"))
	}
	return f.Group(id, parts...)
}
func (f formatter) number(value string) doc {
	if strings.HasPrefix(value, "-") {
		return f.Seq(f.T("-"), f.T(value[1:]))
	}
	return f.T(value)
}
func (f formatter) params(id ast.NodeID) doc {
	var docs []doc
	for _, p := range f.Children(id, "parameter") {
		name := f.Attr(p, "name")
		d := f.Seq()
		if name != "" {
			d = f.T(name)
		}
		docs = append(docs, f.Group(p, f.render(f.Child(p, "type"), d, name != "")))
	}
	if f.Flag(id, "variadic") {
		docs = append(docs, f.T("..."))
	}
	return f.Join(docs, f.T(","), f.Space())
}
func (f formatter) method(id ast.NodeID) doc {
	var parts []doc
	if f.Flag(id, "virtual") {
		parts = append(parts, f.T("virtual"), f.Space())
	}
	suffix := []doc{f.T(f.Attr(id, "name")), f.T("("), f.params(id), f.T(")")}
	for _, q := range []string{"const", "volatile"} {
		if f.Flag(id, q) {
			suffix = append(suffix, f.Space(), f.T(q))
		}
	}
	if ref := f.Attr(id, "ref-qualifier"); ref == "lvalue" {
		suffix = append(suffix, f.Space(), f.T("&"))
	} else if ref == "rvalue" {
		suffix = append(suffix, f.Space(), f.T("&&"))
	}
	if f.Flag(id, "noexcept") {
		suffix = append(suffix, f.Space(), f.T("noexcept"))
	}
	parts = append(parts, f.render(f.Child(id, "result"), f.Seq(suffix...), true))
	if f.Flag(id, "pure") {
		parts = append(parts, f.Space(), f.T("="), f.Space(), f.T("0"))
	}
	parts = append(parts, f.T(";"))
	return f.Group(id, parts...)
}
func (f formatter) destructor(id ast.NodeID, name string) doc {
	var parts []doc
	if f.Flag(id, "virtual") {
		parts = append(parts, f.T("virtual"), f.Space())
	}
	parts = append(parts, f.T("~"), f.T(name), f.T("("), f.T(")"))
	if f.Flag(id, "noexcept") {
		parts = append(parts, f.Space(), f.T("noexcept"))
	}
	if f.Flag(id, "pure") {
		parts = append(parts, f.Space(), f.T("="), f.Space(), f.T("0"))
	}
	parts = append(parts, f.T(";"))
	return f.Group(id, parts...)
}
func (f formatter) typ(id ast.NodeID) doc { return f.render(id, f.Seq(), false) }
func (f formatter) named(id ast.NodeID) doc {
	var parts []doc
	if f.Flag(id, "typename") {
		parts = append(parts, f.T("typename"), f.Space())
	}
	if f.Flag(id, "global") {
		parts = append(parts, f.T("::"))
	}
	var segments []doc
	for _, seg := range f.Children(id, "segment") {
		p := []doc{}
		if f.Flag(seg, "template") {
			p = append(p, f.T("template"), f.Space())
		}
		p = append(p, f.T(f.Attr(seg, "name")))
		if args := f.Children(seg, "argument"); len(args) > 0 {
			var ds []doc
			for _, a := range args {
				ds = append(ds, f.typ(a))
			}
			p = append(p, f.T("<"), f.Join(ds, f.T(","), f.Space()), f.T(">"))
		}
		segments = append(segments, f.Group(seg, p...))
	}
	parts = append(parts, f.Join(segments, f.T("::")))
	return f.Group(id, parts...)
}
func (f formatter) pointer(id ast.NodeID, qualifiers []doc, name doc) doc {
	element := f.Child(id, "element")
	var parts []doc
	switch f.Kind(id) {
	case "pointer-type":
		parts = append(parts, f.T("*"))
	case "lvalue-reference-type":
		parts = append(parts, f.T("&"))
	case "rvalue-reference-type":
		parts = append(parts, f.T("&&"))
	case "member-pointer-type":
		parts = append(parts, f.named(f.Child(id, "owner")), f.T("::"), f.T("*"))
	}
	parts = append(parts, qualifiers...)
	parts = append(parts, name)
	d := f.Seq(parts...)
	kind := f.Kind(element)
	if kind == "array-type" || kind == "function-type" {
		d = f.Seq(f.T("("), d, f.T(")"))
	}
	return f.render(element, d, true)
}
func (f formatter) render(id ast.NodeID, name doc, hasName bool) doc {
	var out doc
	switch f.Kind(id) {
	case "pointer-type", "lvalue-reference-type", "rvalue-reference-type", "member-pointer-type":
		out = f.pointer(id, nil, name)
	case "array-type":
		out = f.render(f.Child(id, "element"), f.Seq(name, f.T("["), f.T(f.Attr(id, "length")), f.T("]")), true)
	case "function-type":
		parts := []doc{name, f.T("("), f.params(id), f.T(")")}
		if f.Flag(id, "noexcept") {
			parts = append(parts, f.Space(), f.T("noexcept"))
		}
		out = f.render(f.Child(id, "result"), f.Seq(parts...), true)
	case "cv-qualified-type":
		var qualifiers []doc
		for _, q := range []string{"const", "volatile"} {
			if f.Flag(id, q) {
				qualifiers = append(qualifiers, f.T(q), f.Space())
			}
		}
		element := f.Child(id, "element")
		kind := f.Kind(element)
		if kind == "pointer-type" || kind == "member-pointer-type" {
			out = f.pointer(element, append([]doc{f.Space()}, qualifiers...), name)
		} else {
			out = f.Seq(append(qualifiers, f.render(element, name, hasName))...)
		}
	case "named-type", "builtin-type":
		var base doc
		if f.Kind(id) == "named-type" {
			base = f.named(id)
		} else {
			var words []doc
			for _, w := range strings.Fields(f.Attr(id, "name")) {
				words = append(words, f.T(w))
			}
			base = f.Join(words, f.Space())
		}
		if hasName {
			out = f.Seq(base, f.Space(), name)
		} else {
			out = base
		}
	}
	return f.Group(id, out)
}

type Provider struct{ s.Language }

func (Provider) Format(view ast.View[s.Language]) (cst.CST[s.Language], error) { return Format(view) }

var _ text.FormatterProvider[s.Language] = Provider{}
