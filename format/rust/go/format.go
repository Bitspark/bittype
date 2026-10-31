package rustformat

import (
	check "github.com/Bitspark/bittype/check/rust/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/internal/concrete"
	"github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
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
	var parts []doc
	for i, id := range f.Children(t.Root, "item") {
		if i > 0 {
			parts = append(parts, f.Line("    ", 0))
		}
		parts = append(parts, f.item(id, 0), f.Line("    ", 0))
	}
	return concrete.Build(t, f.Group(t.Root, parts...))
}
func (f formatter) visibility(id ast.NodeID) doc {
	switch f.Attr(id, "visibility") {
	case "public":
		return f.Seq(f.T("pub"), f.Space())
	case "crate":
		return f.Seq(f.T("pub"), f.T("("), f.T("crate"), f.T(")"), f.Space())
	}
	return f.Seq()
}
func (f formatter) bounds(id ast.NodeID) doc {
	var out []doc
	for _, b := range f.Children(id, "bound") {
		out = append(out, f.typ(b))
	}
	return f.Join(out, f.Space(), f.T("+"), f.Space())
}
func (f formatter) generics(id ast.NodeID) doc {
	var docs []doc
	for _, p := range f.Children(id, "generic-parameter") {
		name := f.Attr(p, "name")
		if f.Kind(p) == "lifetime-parameter" {
			name = "'" + name
		}
		parts := []doc{f.T(name)}
		if len(f.Children(p, "bound")) > 0 {
			parts = append(parts, f.T(":"), f.Space(), f.bounds(p))
		}
		if d := f.Child(p, "default"); d != 0 {
			parts = append(parts, f.Space(), f.T("="), f.Space(), f.typ(d))
		}
		docs = append(docs, f.Group(p, parts...))
	}
	if len(docs) == 0 {
		return f.Seq()
	}
	return f.Seq(f.T("<"), f.Join(docs, f.T(","), f.Space()), f.T(">"))
}
func (f formatter) where(id ast.NodeID) doc {
	var docs []doc
	for _, p := range f.Children(id, "where") {
		docs = append(docs, f.Group(p, f.typ(f.Child(p, "subject")), f.T(":"), f.Space(), f.bounds(p)))
	}
	if len(docs) == 0 {
		return f.Seq()
	}
	return f.Seq(f.Space(), f.T("where"), f.Space(), f.Join(docs, f.T(","), f.Space()))
}
func (f formatter) fields(id ast.NodeID, named bool, level int) doc {
	fields := f.Children(id, "field")
	var parts []doc
	if named {
		parts = append(parts, f.T("{"))
	} else {
		parts = append(parts, f.T("("))
	}
	for i, field := range fields {
		if named {
			parts = append(parts, f.Line("    ", level+1))
		} else if i > 0 {
			parts = append(parts, f.Space())
		}
		p := []doc{f.visibility(field)}
		if named {
			p = append(p, f.T(f.Attr(field, "name")), f.T(":"), f.Space())
		}
		p = append(p, f.typ(f.Child(field, "type")))
		parts = append(parts, f.Group(field, p...), f.T(","))
	}
	if named {
		if len(fields) > 0 {
			parts = append(parts, f.Line("    ", level))
		}
		parts = append(parts, f.T("}"))
	} else {
		parts = append(parts, f.T(")"))
	}
	return f.Seq(parts...)
}
func (f formatter) item(id ast.NodeID, level int) doc {
	kind := f.Kind(id)
	parts := []doc{f.visibility(id)}
	switch kind {
	case "use":
		parts = append(parts, f.T("use"), f.Space(), f.typ(f.Child(id, "path")))
		if a := f.Attr(id, "alias"); a != "" {
			parts = append(parts, f.Space(), f.T("as"), f.Space(), f.T(a))
		}
		parts = append(parts, f.T(";"))
	case "module":
		parts = append(parts, f.T("mod"), f.Space(), f.T(f.Attr(id, "name")), f.Space(), f.T("{"))
		for _, item := range f.Children(id, "item") {
			parts = append(parts, f.Line("    ", level+1), f.item(item, level+1))
		}
		parts = append(parts, f.Line("    ", level), f.T("}"))
	case "type-alias":
		parts = append(parts, f.T("type"), f.Space(), f.T(f.Attr(id, "name")), f.generics(id), f.where(id), f.Space(), f.T("="), f.Space(), f.typ(f.Child(id, "type")), f.T(";"))
	case "struct", "tuple-struct":
		parts = append(parts, f.T("struct"), f.Space(), f.T(f.Attr(id, "name")), f.generics(id))
		if kind == "tuple-struct" {
			parts = append(parts, f.fields(id, false, level), f.where(id), f.T(";"))
		} else {
			parts = append(parts, f.where(id), f.Space(), f.fields(id, true, level))
		}
	case "enum":
		parts = append(parts, f.T("enum"), f.Space(), f.T(f.Attr(id, "name")), f.generics(id), f.where(id), f.Space(), f.T("{"))
		for _, v := range f.Children(id, "variant") {
			parts = append(parts, f.Line("    ", level+1))
			p := []doc{f.T(f.Attr(v, "name"))}
			if f.Kind(v) == "tuple-variant" {
				p = append(p, f.fields(v, false, level+1))
			} else if f.Kind(v) == "record-variant" {
				p = append(p, f.Space(), f.fields(v, true, level+1))
			}
			parts = append(parts, f.Group(v, p...), f.T(","))
		}
		parts = append(parts, f.Line("    ", level), f.T("}"))
	case "trait":
		parts = append(parts, f.T("trait"), f.Space(), f.T(f.Attr(id, "name")), f.generics(id))
		if len(f.Children(id, "bound")) > 0 {
			parts = append(parts, f.T(":"), f.Space(), f.bounds(id))
		}
		parts = append(parts, f.where(id), f.Space(), f.T("{"))
		for _, m := range f.Children(id, "member") {
			parts = append(parts, f.Line("    ", level+1))
			var p []doc
			if f.Kind(m) == "associated-type" {
				p = []doc{f.T("type"), f.Space(), f.T(f.Attr(m, "name"))}
				if len(f.Children(m, "bound")) > 0 {
					p = append(p, f.T(":"), f.Space(), f.bounds(m))
				}
			} else {
				p = append(p, f.T("fn"), f.Space(), f.T(f.Attr(m, "name")), f.generics(m), f.T("("))
				var params []doc
				if r := f.Child(m, "receiver"); r != 0 {
					params = append(params, f.receiver(r))
				}
				params = append(params, f.parameters(m)...)
				p = append(p, f.Join(params, f.T(","), f.Space()), f.T(")"))
				if result := f.Child(m, "result"); result != 0 {
					p = append(p, f.Space(), f.T("->"), f.Space(), f.typ(result))
				}
				p = append(p, f.where(m))
			}
			p = append(p, f.T(";"))
			parts = append(parts, f.Group(m, p...))
		}
		parts = append(parts, f.Line("    ", level), f.T("}"))
	}
	return f.Group(id, parts...)
}
func (f formatter) receiver(id ast.NodeID) doc {
	var parts []doc
	form := f.Attr(id, "receiver-form")
	if form != "value" {
		parts = append(parts, f.T("&"))
		if l := f.Child(id, "lifetime"); l != 0 {
			parts = append(parts, f.typ(l), f.Space())
		}
		if form == "mutable" {
			parts = append(parts, f.T("mut"), f.Space())
		}
	}
	parts = append(parts, f.T("self"))
	return f.Group(id, parts...)
}
func (f formatter) parameters(id ast.NodeID) []doc {
	var docs []doc
	for _, p := range f.Children(id, "parameter") {
		var parts []doc
		if name := f.Attr(p, "name"); name != "" {
			parts = append(parts, f.T(name), f.T(":"), f.Space())
		}
		parts = append(parts, f.typ(f.Child(p, "type")))
		docs = append(docs, f.Group(p, parts...))
	}
	return docs
}
func (f formatter) element(id ast.NodeID) doc {
	e := f.Child(id, "element")
	d := f.typ(e)
	if f.Kind(e) == "trait-object-type" && len(f.Children(e, "bound")) > 1 {
		return f.Seq(f.T("("), d, f.T(")"))
	}
	return d
}
func (f formatter) typ(id ast.NodeID) doc {
	var parts []doc
	switch f.Kind(id) {
	case "primitive-type":
		parts = []doc{f.T(f.Attr(id, "name"))}
	case "lifetime-reference":
		parts = []doc{f.T("'" + f.Attr(id, "name"))}
	case "path-type":
		if f.Flag(id, "absolute") {
			parts = append(parts, f.T("::"))
		}
		var segments []doc
		for _, seg := range f.Children(id, "segment") {
			p := []doc{f.T(f.Attr(seg, "name"))}
			args := f.Children(seg, "argument")
			if len(args) > 0 {
				var docs []doc
				for _, a := range args {
					docs = append(docs, f.typ(a))
				}
				p = append(p, f.T("<"), f.Join(docs, f.T(","), f.Space()), f.T(">"))
			}
			segments = append(segments, f.Group(seg, p...))
		}
		parts = append(parts, f.Join(segments, f.T("::")))
	case "reference-type":
		parts = append(parts, f.T("&"))
		if l := f.Child(id, "lifetime"); l != 0 {
			parts = append(parts, f.typ(l), f.Space())
		}
		if f.Flag(id, "mutable") {
			parts = append(parts, f.T("mut"), f.Space())
		}
		parts = append(parts, f.element(id))
	case "raw-pointer-type":
		qualifier := "const"
		if f.Flag(id, "mutable") {
			qualifier = "mut"
		}
		parts = []doc{f.T("*"), f.T(qualifier), f.Space(), f.element(id)}
	case "slice-type":
		parts = []doc{f.T("["), f.typ(f.Child(id, "element")), f.T("]")}
	case "array-type":
		parts = []doc{f.T("["), f.typ(f.Child(id, "element")), f.T(";"), f.Space(), f.T(f.Attr(id, "length")), f.T("]")}
	case "tuple-type":
		var elems []doc
		for _, e := range f.Children(id, "element") {
			elems = append(elems, f.typ(e))
		}
		parts = []doc{f.T("("), f.Join(elems, f.T(","), f.Space())}
		if len(elems) == 1 {
			parts = append(parts, f.T(","))
		}
		parts = append(parts, f.T(")"))
	case "never-type":
		parts = []doc{f.T("!")}
	case "trait-object-type":
		parts = []doc{f.T("dyn"), f.Space(), f.bounds(id)}
	case "function-type":
		if f.Flag(id, "unsafe") {
			parts = append(parts, f.T("unsafe"), f.Space())
		}
		if abi := f.Attr(id, "abi"); abi != "" {
			parts = append(parts, f.T("extern"), f.Space(), f.T(strconv.Quote(abi)), f.Space())
		}
		parts = append(parts, f.T("fn"), f.T("("), f.Join(f.parameters(id), f.T(","), f.Space()), f.T(")"))
		if result := f.Child(id, "result"); result != 0 {
			parts = append(parts, f.Space(), f.T("->"), f.Space(), f.typ(result))
		}
	}
	return f.Group(id, parts...)
}

type Provider struct{ s.Language }

func (Provider) Format(view ast.View[s.Language]) (cst.CST[s.Language], error) { return Format(view) }

var _ text.FormatterProvider[s.Language] = Provider{}
