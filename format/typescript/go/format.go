package typescriptformat

import (
	"github.com/Bitspark/bittype/trees/ast/go"
	check "github.com/Bitspark/bittype/check/typescript/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/internal/concrete"
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/internal/tree"
	lit "github.com/Bitspark/bittype/internal/typescript"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/contracts/text/go"
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
	return concrete.Build(t, f.file(t.Root))
}
func (f formatter) name(name string) doc {
	if !g.Identifier(name) {
		name = lit.Quote(name)
	}
	return f.T(name)
}
func (f formatter) file(id ast.NodeID) doc {
	var out []doc
	for _, imp := range f.Children(id, "import") {
		parts := []doc{f.T("import"), f.Space(), f.T("type"), f.Space()}
		bindings := f.Children(imp, "binding")
		if len(bindings) == 1 && f.Kind(bindings[0]) == "namespace-import" {
			parts = append(parts, f.T("*"), f.Space(), f.T("as"), f.Space(), f.T(f.Attr(bindings[0], "name")))
		} else {
			var names []doc
			for _, b := range bindings {
				p := []doc{f.T(f.Attr(b, "name"))}
				if alias := f.Attr(b, "alias"); alias != "" {
					p = append(p, f.Space(), f.T("as"), f.Space(), f.T(alias))
				}
				names = append(names, f.Group(b, p...))
			}
			parts = append(parts, f.T("{"), f.Space(), f.Join(names, f.T(","), f.Space()), f.Space(), f.T("}"))
		}
		parts = append(parts, f.Space(), f.T("from"), f.Space(), f.T(lit.Quote(f.Attr(imp, "module"))), f.T(";"), f.Line("  ", 0))
		out = append(out, f.Group(imp, parts...))
	}
	for i, decl := range f.Children(id, "declaration") {
		if i > 0 || len(out) > 0 {
			out = append(out, f.Line("  ", 0))
		}
		parts := []doc{}
		if f.Flag(decl, "exported") {
			parts = append(parts, f.T("export"), f.Space())
		}
		if f.Kind(decl) == "type-alias" {
			parts = append(parts, f.T("type"), f.Space(), f.T(f.Attr(decl, "name")), f.generics(decl, 0), f.Space(), f.T("="), f.Space(), f.typ(f.Child(decl, "type"), 0, 0), f.T(";"))
		} else {
			parts = append(parts, f.T("interface"), f.Space(), f.T(f.Attr(decl, "name")), f.generics(decl, 0))
			bases := f.Children(decl, "base")
			if len(bases) > 0 {
				var ds []doc
				for _, b := range bases {
					ds = append(ds, f.typ(b, 0, 0))
				}
				parts = append(parts, f.Space(), f.T("extends"), f.Space(), f.Join(ds, f.T(","), f.Space()))
			}
			parts = append(parts, f.Space(), f.members(decl, 0))
		}
		out = append(out, f.Group(decl, parts...), f.Line("  ", 0))
	}
	return f.Group(id, out...)
}
func (f formatter) generics(id ast.NodeID, level int) doc {
	var docs []doc
	for _, p := range f.Children(id, "type-parameter") {
		parts := []doc{f.T(f.Attr(p, "name"))}
		if c := f.Child(p, "constraint"); c != 0 {
			parts = append(parts, f.Space(), f.T("extends"), f.Space(), f.typ(c, 0, level))
		}
		if d := f.Child(p, "default"); d != 0 {
			parts = append(parts, f.Space(), f.T("="), f.Space(), f.typ(d, 0, level))
		}
		docs = append(docs, f.Group(p, parts...))
	}
	if len(docs) == 0 {
		return f.Seq()
	}
	return f.Seq(f.T("<"), f.Join(docs, f.T(","), f.Space()), f.T(">"))
}
func (f formatter) signature(id ast.NodeID, arrow bool, level int) doc {
	var params []doc
	for _, p := range f.Children(id, "parameter") {
		parts := []doc{}
		if f.Flag(p, "rest") {
			parts = append(parts, f.T("..."))
		}
		parts = append(parts, f.T(f.Attr(p, "name")))
		if f.Flag(p, "optional") {
			parts = append(parts, f.T("?"))
		}
		parts = append(parts, f.T(":"), f.Space(), f.typ(f.Child(p, "type"), 0, level))
		params = append(params, f.Group(p, parts...))
	}
	parts := []doc{f.generics(id, level), f.T("("), f.Join(params, f.T(","), f.Space()), f.T(")")}
	if arrow {
		parts = append(parts, f.Space(), f.T("=>"))
	} else {
		parts = append(parts, f.T(":"))
	}
	parts = append(parts, f.Space(), f.typ(f.Child(id, "result"), 0, level))
	return f.Group(id, parts...)
}
func (f formatter) members(id ast.NodeID, level int) doc {
	parts := []doc{f.T("{")}
	members := f.Children(id, "member")
	for _, m := range members {
		parts = append(parts, f.Line("  ", level+1))
		var p []doc
		switch f.Kind(m) {
		case "property":
			if f.Flag(m, "readonly") {
				p = append(p, f.T("readonly"), f.Space())
			}
			p = append(p, f.name(f.Attr(m, "name")))
			if f.Flag(m, "optional") {
				p = append(p, f.T("?"))
			}
			p = append(p, f.T(":"), f.Space(), f.typ(f.Child(m, "type"), 0, level+1))
		case "method":
			p = append(p, f.name(f.Attr(m, "name")))
			if f.Flag(m, "optional") {
				p = append(p, f.T("?"))
			}
			p = append(p, f.signature(f.Child(m, "signature"), false, level+1))
		case "call-signature":
			p = append(p, f.signature(f.Child(m, "signature"), false, level+1))
		case "index-signature":
			if f.Flag(m, "readonly") {
				p = append(p, f.T("readonly"), f.Space())
			}
			parameter := f.Child(m, "parameter")
			p = append(p, f.T("["), f.T(f.Attr(parameter, "name")), f.T(":"), f.Space(), f.typ(f.Child(parameter, "type"), 0, level+1), f.T("]"), f.T(":"), f.Space(), f.typ(f.Child(m, "type"), 0, level+1))
		}
		p = append(p, f.T(";"))
		parts = append(parts, f.Group(m, p...))
	}
	if len(members) > 0 {
		parts = append(parts, f.Line("  ", level))
	}
	parts = append(parts, f.T("}"))
	return f.Seq(parts...)
}
func (f formatter) typ(id ast.NodeID, parent, level int) doc {
	prec := 5
	var parts []doc
	t := func(role string, p int) doc { return f.typ(f.Child(id, role), p, level) }
	switch f.Kind(id) {
	case "type-reference", "keyword-type":
		parts = []doc{f.T(f.Attr(id, "name"))}
	case "qualified-type":
		parts = []doc{t("qualifier", 5), f.T("."), f.T(f.Attr(id, "name"))}
	case "type-application":
		var args []doc
		for _, a := range f.Children(id, "argument") {
			args = append(args, f.typ(a, 0, level))
		}
		parts = []doc{t("constructor", 5), f.T("<"), f.Join(args, f.T(","), f.Space()), f.T(">")}
	case "literal-type":
		parts = []doc{f.T(f.Attr(id, "text"))}
	case "array-type":
		prec = 4
		parts = []doc{t("element", 4), f.T("["), f.T("]")}
	case "indexed-access-type":
		prec = 4
		parts = []doc{t("object", 4), f.T("["), t("index", 0), f.T("]")}
	case "type-operator":
		prec = 3
		parts = []doc{f.T(f.Attr(id, "operator")), f.Space(), t("type", 3)}
	case "union-type", "intersection-type":
		prec = 1
		symbol := "|"
		if f.Kind(id) == "intersection-type" {
			prec = 2
			symbol = "&"
		}
		var types []doc
		for _, c := range f.Children(id, "type") {
			types = append(types, f.typ(c, prec+1, level))
		}
		parts = []doc{f.Join(types, f.Space(), f.T(symbol), f.Space())}
	case "function-type":
		prec = 0
		parts = []doc{f.signature(id, true, level)}
	case "object-type":
		parts = []doc{f.members(id, level)}
	case "tuple-type":
		var elems []doc
		for _, e := range f.Children(id, "element") {
			var p []doc
			if f.Flag(e, "rest") {
				p = append(p, f.T("..."))
			}
			name := f.Attr(e, "name")
			if name != "" {
				p = append(p, f.T(name))
				if f.Flag(e, "optional") {
					p = append(p, f.T("?"))
				}
				p = append(p, f.T(":"), f.Space())
			}
			p = append(p, f.typ(f.Child(e, "type"), 0, level))
			if name == "" && f.Flag(e, "optional") {
				p = append(p, f.T("?"))
			}
			elems = append(elems, f.Group(e, p...))
		}
		parts = []doc{f.T("["), f.Join(elems, f.T(","), f.Space()), f.T("]")}
	}
	d := f.Group(id, parts...)
	if prec < parent {
		return f.Seq(f.T("("), d, f.T(")"))
	}
	return d
}

type Provider struct{ s.Language }

func (Provider) Format(view ast.View[s.Language]) (cst.CST[s.Language], error) { return Format(view) }

var _ text.FormatterProvider[s.Language] = Provider{}
