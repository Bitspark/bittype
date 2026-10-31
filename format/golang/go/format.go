package golangformat

import (
	"github.com/Bitspark/bittype/ast/go"
	check "github.com/Bitspark/bittype/check/golang/go"
	"github.com/Bitspark/bittype/cst/go"
	"github.com/Bitspark/bittype/internal/concrete"
	"github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/text/go"
	"strconv"
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
func (f formatter) file(id ast.NodeID) doc {
	parts := []doc{f.T("package"), f.Space(), f.T(f.Attr(id, "package"))}
	for i, imp := range f.Children(id, "import") {
		if i == 0 {
			parts = append(parts, f.Line("\t", 0))
		}
		parts = append(parts, f.Line("\t", 0), f.T("import"), f.Space())
		if alias := f.Attr(imp, "alias"); alias != "" {
			parts = append(parts, f.T(alias), f.Space())
		}
		parts = append(parts, f.Group(imp, f.T(strconv.Quote(f.Attr(imp, "path")))))
	}
	for _, decl := range f.Children(id, "declaration") {
		parts = append(parts, f.Line("\t", 0), f.Line("\t", 0), f.decl(decl))
	}
	parts = append(parts, f.Line("\t", 0))
	return f.Group(id, parts...)
}
func (f formatter) decl(id ast.NodeID) doc {
	parts := []doc{f.T("type"), f.Space(), f.T(f.Attr(id, "name"))}
	params := f.Children(id, "type-parameter")
	if len(params) > 0 {
		var docs []doc
		for _, p := range params {
			docs = append(docs, f.Group(p, f.T(f.Attr(p, "name")), f.Space(), f.typ(f.Child(p, "constraint"), 0)))
		}
		parts = append(parts, f.T("["), f.Join(docs, f.T(","), f.Space()), f.T("]"))
	}
	parts = append(parts, f.Space())
	if f.Kind(id) == "type-alias" {
		parts = append(parts, f.T("="), f.Space())
	}
	parts = append(parts, f.typ(f.Child(id, "type"), 0))
	return f.Group(id, parts...)
}
func (f formatter) parameters(ids []ast.NodeID, level int) doc {
	var docs []doc
	for _, id := range ids {
		var parts []doc
		if name := f.Attr(id, "name"); name != "" {
			parts = append(parts, f.T(name), f.Space())
		}
		if f.Flag(id, "variadic") {
			parts = append(parts, f.T("..."))
		}
		parts = append(parts, f.typ(f.Child(id, "type"), level))
		docs = append(docs, f.Group(id, parts...))
	}
	return f.Join(docs, f.T(","), f.Space())
}
func (f formatter) signature(id ast.NodeID, level int) doc {
	parts := []doc{f.T("("), f.parameters(f.Children(id, "parameter"), level), f.T(")")}
	results := f.Children(id, "result")
	if len(results) == 1 && !f.Has(results[0], "name") {
		parts = append(parts, f.Space(), f.typ(f.Child(results[0], "type"), level))
	} else if len(results) > 0 {
		parts = append(parts, f.Space(), f.T("("), f.parameters(results, level), f.T(")"))
	}
	return f.Group(id, parts...)
}
func (f formatter) typ(id ast.NodeID, level int) doc {
	t := func(role string) doc { return f.typ(f.Child(id, role), level) }
	var parts []doc
	switch f.Kind(id) {
	case "named-type":
		parts = []doc{f.T(f.Attr(id, "name"))}
	case "qualified-type":
		parts = []doc{f.T(f.Attr(id, "qualifier")), f.T("."), f.T(f.Attr(id, "name"))}
	case "type-application":
		var args []doc
		for _, a := range f.Children(id, "argument") {
			args = append(args, f.typ(a, level))
		}
		parts = []doc{t("constructor"), f.T("["), f.Join(args, f.T(","), f.Space()), f.T("]")}
	case "pointer-type":
		parts = []doc{f.T("*"), t("element")}
	case "slice-type":
		parts = []doc{f.T("["), f.T("]"), t("element")}
	case "array-type":
		parts = []doc{f.T("["), f.T(f.Attr(id, "length")), f.T("]"), t("element")}
	case "map-type":
		parts = []doc{f.T("map"), f.T("["), t("key"), f.T("]"), t("element")}
	case "channel-type":
		d := f.Attr(id, "direction")
		if d == "receive" {
			parts = append(parts, f.T("<-"))
		}
		parts = append(parts, f.T("chan"))
		if d == "send" {
			parts = append(parts, f.T("<-"))
		}
		parts = append(parts, f.Space())
		if f.Kind(f.Child(id, "element")) == "channel-type" {
			parts = append(parts, f.T("("), t("element"), f.T(")"))
		} else {
			parts = append(parts, t("element"))
		}
	case "function-type":
		parts = []doc{f.T("func"), f.signature(id, level)}
	case "struct-type", "interface-type":
		keyword, role := "struct", "field"
		if f.Kind(id) == "interface-type" {
			keyword, role = "interface", "member"
		}
		parts = []doc{f.T(keyword), f.Space(), f.T("{")}
		members := f.Children(id, role)
		for _, m := range members {
			parts = append(parts, f.Line("\t", level+1))
			var member []doc
			switch f.Kind(m) {
			case "field":
				member = append(member, f.T(f.Attr(m, "name")), f.Space(), f.typ(f.Child(m, "type"), level+1))
			case "embedded-field", "embedded-type":
				member = append(member, f.typ(f.Child(m, "type"), level+1))
			case "method":
				member = append(member, f.T(f.Attr(m, "name")), f.signature(f.Child(m, "signature"), level+1))
			}
			if f.Has(m, "tag") {
				member = append(member, f.Space(), f.T(strconv.Quote(f.Attr(m, "tag"))))
			}
			parts = append(parts, f.Group(m, member...))
		}
		if len(members) > 0 {
			parts = append(parts, f.Line("\t", level))
		}
		parts = append(parts, f.T("}"))
	case "type-union":
		var terms []doc
		for _, term := range f.Children(id, "term") {
			terms = append(terms, f.typ(term, level))
		}
		parts = []doc{f.Join(terms, f.Space(), f.T("|"), f.Space())}
	case "type-approximation":
		parts = []doc{f.T("~"), t("type")}
	}
	return f.Group(id, parts...)
}

type Provider struct{ s.Language }

func (Provider) Format(view ast.View[s.Language]) (cst.CST[s.Language], error) { return Format(view) }

var _ text.FormatterProvider[s.Language] = Provider{}
