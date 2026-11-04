// Package horosformat constructs concrete Horos syntax directly from an AST.
package horosformat

import (
	"encoding/hex"
	"reflect"
	"strconv"
	"unicode/utf8"

	check "github.com/Bitspark/bittype/check/horos/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/internal/concrete"
	lit "github.com/Bitspark/bittype/internal/horos"
	"github.com/Bitspark/bittype/internal/horosmodel"
	"github.com/Bitspark/bittype/internal/lex"
	"github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

type doc = concrete.Doc[s.Language]
type formatter struct {
	tree.View[s.Language]
	concrete.Builder[s.Language]
	parameters []string
	err        error
}

// Format preserves canonical values, including unresolved applications. A
// multi-module Program is emitted as sections with its selected entry first.
// Formatting does not resolve imports, specialize generics or repair assertions.
func Format(view ast.View[s.Language]) (cst.CST[s.Language], error) {
	if view == nil {
		return cst.CST[s.Language]{}, check.Check(view)
	}
	v := reflect.ValueOf(view)
	switch v.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface, reflect.Func, reflect.Chan:
		if v.IsNil() {
			return cst.CST[s.Language]{}, check.Check(view)
		}
	}
	t := view.ASTView()
	if err := check.Check(t); err != nil {
		return cst.CST[s.Language]{}, err
	}
	pins, err := horosmodel.SourcePins(t)
	if err != nil {
		return cst.CST[s.Language]{}, err
	}
	f := formatter{View: tree.View[s.Language]{Tree: t}, Builder: concrete.Builder[s.Language]{Vocabulary: concrete.Vocabulary[s.Language]{Layout: s.ConcreteLayout, Token: s.ConcreteToken, Whitespace: s.ConcreteWhitespace, Comment: s.ConcreteComment}}}
	for id, pin := range pins {
		if f.Attr(id, s.AttributePin) != pin {
			return cst.CST[s.Language]{}, text.AtNode(t, id, "module pin does not match its declaration bytes")
		}
	}
	root := t.Nodes[t.Root-1].References[0].Target.Node
	modules := []ast.NodeID{root}
	for _, child := range t.Nodes[t.Root-1].Children {
		if child.Node != root {
			modules = append(modules, child.Node)
		}
	}
	var parts []doc
	for i, id := range modules {
		if i > 0 {
			parts = append(parts, f.Line("    ", 0))
		}
		parts = append(parts, f.module(id))
	}
	if f.err != nil {
		return cst.CST[s.Language]{}, f.err
	}
	result, err := concrete.Build(t, f.Group(t.Root, parts...))
	if err != nil {
		return cst.CST[s.Language]{}, err
	}
	bytes, leaves := 0, 0
	for _, n := range result.Nodes {
		if n.Text != nil {
			bytes += len(*n.Text)
			leaves++
		}
	}
	if bytes > lex.MaxBytes || leaves > lex.MaxTokens {
		return cst.CST[s.Language]{}, &text.Diagnostic{Code: text.LimitExceeded, Message: "formatted source exceeds parser limits"}
	}
	return result, nil
}

func (f *formatter) bad(id ast.NodeID, message string) doc {
	if f.err == nil {
		f.err = text.AtNode(f.Tree, id, message)
	}
	return f.Seq()
}

func (f *formatter) module(id ast.NodeID) doc {
	parts := []doc{f.T("module"), f.Space(), f.T(lit.ModuleName(f.Attr(id, s.AttributeName))), f.T(";"), f.Line("    ", 0)}
	for _, imp := range f.Children(id, s.RoleImport) {
		parts = append(parts, f.Group(imp, f.T("import"), f.Space(), f.T(lit.Name(f.Attr(imp, s.AttributeAlias))), f.Space(), f.T("from"), f.Space(), f.T(lit.ModuleName(f.Attr(imp, s.AttributeName))), f.Space(), f.T("pin"), f.Space(), f.T(lit.Quote("sha256:"+f.Attr(imp, s.AttributePin))), f.T(";")), f.Line("    ", 0))
	}
	for _, declaration := range f.Children(id, s.RoleDeclaration) {
		parts = append(parts, f.Line("    ", 0), f.declaration(declaration), f.Line("    ", 0))
	}
	return f.Group(id, parts...)
}

func (f *formatter) declaration(id ast.NodeID) doc {
	f.parameters = nil
	var parameters []doc
	seen := map[string]bool{}
	for i, parameter := range f.Children(id, s.RoleParameter) {
		name := f.Attr(parameter, s.AttributeName)
		if f.Attr(parameter, s.AttributeIndex) != strconv.Itoa(i) || seen[name] {
			return f.bad(parameter, "parameters require distinct names and ordered contiguous indices")
		}
		seen[name] = true
		f.parameters = append(f.parameters, name)
		bound := f.Child(parameter, s.RoleBound)
		parts := []doc{f.T(lit.Name(name))}
		if f.Attr(bound, s.AttributeName) == "data" {
			parts = append(parts, f.T(":"), f.Space(), f.Group(bound, f.T("Data")))
		}
		parameters = append(parameters, f.Group(parameter, parts...))
	}
	parts := []doc{f.T("type"), f.Space(), f.T(lit.Name(f.Attr(id, s.AttributeName)))}
	if len(parameters) > 0 {
		parts = append(parts, f.T("<"), f.Join(parameters, f.T(","), f.Space()), f.T(">"))
	}
	parts = append(parts, f.Space(), f.T("="), f.Space(), f.typ(f.Child(id, s.RoleBody), 0, 0), f.T(";"))
	return f.Group(id, parts...)
}

func (f *formatter) literal(id ast.NodeID) core.Value {
	value, err := horosmodel.SourceLiteral(f.Tree, id)
	if err != nil {
		f.bad(id, err.Error())
	}
	return value
}

func (f *formatter) atom(id ast.NodeID) (string, bool) {
	value, ok := f.literal(id).(core.Atom)
	if !ok {
		f.bad(id, "expected literal atom")
		return "", false
	}
	return string(value.Bytes()), true
}

func (f *formatter) args(name string, values []doc) doc {
	return f.Seq(f.T(name), f.T("<"), f.Join(values, f.T(","), f.Space()), f.T(">"))
}

func (f *formatter) call(name string, values []doc) doc {
	return f.Seq(f.T(name), f.T("("), f.Join(values, f.T(","), f.Space()), f.T(")"))
}

func (f *formatter) typ(id ast.NodeID, parent, level int) doc {
	if f.err != nil {
		return f.Seq()
	}
	var result doc
	prec := 3
	switch f.Node(id).Kind {
	case s.KindParameter:
		index, err := strconv.Atoi(f.Attr(id, s.AttributeIndex))
		if err != nil || index < 0 || index >= len(f.parameters) {
			return f.bad(id, "parameter index has no source binder")
		}
		result = f.T(lit.Name(f.parameters[index]))
	case s.KindApplication:
		name := f.Attr(id, s.AttributeName)
		parts := []doc{}
		if f.Attr(id, s.AttributeScope) == "imported" {
			parts = append(parts, f.T(lit.Name(f.Attr(id, s.AttributeAlias))), f.T("."))
		} else {
			shadowed := false
			for _, parameter := range f.parameters {
				if parameter == name {
					shadowed = true
				}
			}
			if shadowed {
				parts = append(parts, f.T("self"), f.T("."))
			}
		}
		parts = append(parts, f.T(lit.Name(name)))
		var arguments []doc
		for _, argument := range f.Children(id, s.RoleArgument) {
			arguments = append(arguments, f.typ(argument, 0, level))
		}
		if len(arguments) > 0 {
			parts = append(parts, f.T("<"), f.Join(arguments, f.T(","), f.Space()), f.T(">"))
		}
		result = f.Seq(parts...)
	case s.KindDescriptor:
		name := f.Attr(id, s.AttributeConstructor)
		operands := f.Children(id, s.RoleOperand)
		arity := func(n int) bool {
			if len(operands) != n {
				f.bad(id, "wrong source constructor arity")
				return false
			}
			return true
		}
		switch name {
		case "any", "atom":
			if !arity(0) {
				return f.Seq()
			}
			result = f.T(name)
		case "scalar":
			if !arity(1) {
				return f.Seq()
			}
			kind, ok := f.atom(operands[0])
			if !ok {
				return f.Seq()
			}
			switch kind {
			case "utf8-text":
				kind = "text"
			case "int", "bool":
			default:
				return f.bad(id, "unknown scalar spelling")
			}
			result = f.T(kind)
		case "tuple", "list-of", "seq", "set-of", "map-of", "not", "wire", "union", "and":
			switch name {
			case "list-of", "seq", "set-of", "not", "wire":
				if !arity(1) {
					return f.Seq()
				}
			case "map-of":
				if !arity(2) {
					return f.Seq()
				}
			}
			operator := ""
			if len(operands) >= 2 {
				if name == "union" {
					operator, prec = "|", 1
				}
				if name == "and" {
					operator, prec = "&", 2
				}
			}
			var values []doc
			for _, operand := range operands {
				childParent := 0
				if operator != "" {
					childParent = prec
				}
				values = append(values, f.typ(operand, childParent, level))
			}
			if operator != "" {
				result = f.Join(values, f.Space(), f.T(operator), f.Space())
			} else {
				switch name {
				case "list-of":
					name = "list"
				case "set-of":
					name = "set"
				case "map-of":
					name = "map"
				}
				result = f.args(name, values)
			}
		case "record":
			result = f.record(id, operands, level)
		case "exact", "atom-exact", "atom-in", "atom-len", "text-len", "key-grammar", "text-non-blank", "lowercase-uuid", "compound", "rec", "ref":
			count := 1
			switch name {
			case "key-grammar", "text-non-blank", "lowercase-uuid":
				count = 0
			case "atom-len", "text-len", "rec":
				count = 2
			case "compound":
				count = -1
			}
			if count >= 0 && !arity(count) {
				return f.Seq()
			}
			if name == "compound" && len(operands) == 0 {
				return f.bad(id, "compound requires a label")
			}
			var values []doc
			for i, operand := range operands {
				if (name == "compound" || name == "rec") && i > 0 {
					values = append(values, f.typ(operand, 0, level))
				} else {
					values = append(values, f.Group(operand, f.value(f.literal(operand))))
				}
			}
			result = f.call(name, values)
		default:
			return f.bad(id, "constructor has no admitted source spelling")
		}
	default:
		return f.bad(id, "expected type expression")
	}
	if prec <= parent {
		result = f.Seq(f.T("("), result, f.T(")"))
	}
	return f.Group(id, result)
}

func (f *formatter) record(id ast.NodeID, operands []ast.NodeID, level int) doc {
	if len(operands) == 0 {
		return f.bad(id, "record requires mode")
	}
	mode, ok := f.atom(operands[0])
	if !ok {
		return f.Seq()
	}
	parts := []doc{}
	if mode == "open" {
		parts = append(parts, f.T("open"), f.Space())
	} else if mode != "closed" {
		return f.bad(id, "record mode must be open or closed")
	}
	parts = append(parts, f.T("{"))
	for _, field := range operands[1:] {
		args := f.Children(field, s.RoleElement)
		if f.Node(field).Kind != s.KindTuple || len(args) != 4 {
			return f.bad(field, "expected a record field tuple")
		}
		label, ok := f.atom(args[0])
		if !ok || label != "horos/field" {
			return f.bad(field, "expected horos/field label")
		}
		args = args[1:]
		name, err := data.ReadText(f.literal(args[0]))
		if err != nil {
			return f.bad(field, "field name must be text")
		}
		presence, ok := f.atom(args[1])
		if !ok {
			return f.Seq()
		}
		if presence != "req" && presence != "opt" {
			return f.bad(field, "field presence must be req or opt")
		}
		spelling := lit.Quote(name)
		if lit.Identifier(name) {
			spelling = lit.Name(name)
		}
		p := []doc{f.Group(args[0], f.T(spelling))}
		if presence == "opt" {
			p = append(p, f.Group(args[1], f.T("?")))
		}
		p = append(p, f.T(":"), f.Space(), f.typ(args[2], 0, level+1), f.T(";"))
		parts = append(parts, f.Line("    ", level+1), f.Group(field, p...))
	}
	if len(operands) > 1 {
		parts = append(parts, f.Line("    ", level))
	}
	return f.Seq(append(parts, f.T("}"))...)
}

// Prefer readable canonical embeddings, falling back to raw tuples for every
// other value. Recognition never changes the value or descends as type syntax.
func (f *formatter) value(value core.Value) doc {
	if f.err != nil {
		return f.Seq()
	}
	if atom, ok := value.(core.Atom); ok {
		bytes := atom.Bytes()
		if utf8.Valid(bytes) {
			return f.Seq(f.T("#"), f.T(lit.Quote(string(bytes))))
		}
		return f.Seq(f.T("#"), f.T("hex"), f.T(lit.Quote(hex.EncodeToString(bytes))))
	}
	if value == nil {
		return f.Seq()
	}
	if textValue, err := data.ReadText(value); err == nil {
		return f.T(lit.Quote(textValue))
	}
	if integer, err := data.ReadInt(value); err == nil && integer.BitLen() <= lit.MaxIntegerDigits*4 {
		spelling := integer.String()
		if lit.Integer(spelling) {
			return f.T(spelling)
		}
	}
	if boolean, err := data.ReadBool(value); err == nil {
		return f.T(strconv.FormatBool(boolean))
	}
	if data.ReadNull(value) == nil {
		return f.T("null")
	}
	if values, err := data.ReadSet(value); err == nil {
		var parts []doc
		for _, v := range values {
			parts = append(parts, f.value(v))
		}
		return f.call("set", parts)
	}
	tuple := value.(core.Tuple)
	var parts []doc
	for i := 0; i < tuple.Len(); i++ {
		parts = append(parts, f.value(tuple.At(i)))
	}
	return f.call("tuple", parts)
}
