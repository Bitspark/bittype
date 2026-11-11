// Package parse parses readable Horos declarations without resolving them.
package parse

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"
	"strings"

	check "github.com/Bitspark/bittype/check/syntax/horos/provider/go"
	structure "github.com/Bitspark/bittype/check/syntax/shared/ast/go"
	parsecontract "github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/internal/concrete"
	lit "github.com/Bitspark/bittype/internal/horos"
	"github.com/Bitspark/bittype/internal/horosmodel"
	"github.com/Bitspark/bittype/internal/lex"
	"github.com/Bitspark/bittype/internal/parseutil"
	tr "github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

var vocabulary = concrete.Vocabulary[s.Language]{Layout: s.ConcreteLayout, Token: s.ConcreteToken, Whitespace: s.ConcreteWhitespace, Comment: s.ConcreteComment}

type parser struct {
	*parseutil.Parser[s.Language]
	parameters     map[string]int
	kindParameters map[string]int
	groups         []map[string]int
	ascribed       map[ast.NodeID]bool
}

func (p *parser) Add(kind s.Kind, start int, attrs []ast.Attribute, children ...ast.Child) ast.NodeID {
	end := start
	if p.Pos > 0 {
		end = p.Tokens[p.Pos-1].End
	}
	return p.AddSpan(kind, start, end, attrs, children...)
}

func (p *parser) AddSpan(kind s.Kind, start, end int, attrs []ast.Attribute, children ...ast.Child) ast.NodeID {
	if p.Err != nil {
		return 0
	}
	if len(p.Tree.Nodes) >= structure.MaxNodes {
		p.Fail(text.LimitExceeded, "too many abstract nodes")
		return 0
	}
	return p.Parser.AddSpan(kind, start, end, attrs, children...)
}

// Parse admits one or more module sections. The first selects the program entry;
// subsequent sections explicitly supply dependency documents. No IO is done.
// bittype:scope docs/horos-syntax.md#pipeline-and-limits
func Parse(ctx context.Context, input parsecontract.Input) (cst.CST[s.Language], error) {
	if ctx == nil {
		return cst.CST[s.Language]{}, errors.New("nil context")
	}
	if err := ctx.Err(); err != nil {
		return cst.CST[s.Language]{}, err
	}
	base, err := parseutil.New[s.Language](input, lex.Options{BacktickStrings: true, Punctuation: []string{"->"}})
	if err != nil {
		return cst.CST[s.Language]{}, err
	}
	p := parser{Parser: base}
	var modules []ast.Child
	for !p.EOF() && p.Err == nil {
		modules = append(modules, tr.Child(s.RoleModule, p.module()))
	}
	if len(modules) == 0 {
		p.Fail(text.SyntaxError, "expected module declaration")
	}
	if err := p.End(); err != nil {
		return cst.CST[s.Language]{}, err
	}
	root := p.AddSpan(s.KindProgram, 0, len(input.Text), nil, modules...)
	tree := p.Compact(root)
	if p.Err != nil {
		return cst.CST[s.Language]{}, p.Err
	}
	tree.Nodes[tree.Root-1].References = []ast.Reference{{Role: s.RoleRoot, Target: ast.NodeRef{Tree: tree.ID, Node: tree.Nodes[tree.Root-1].Children[0].Node}}}
	if err := check.Check(tree); err != nil {
		return cst.CST[s.Language]{}, err
	}
	pins, err := horosmodel.SourcePins(tree)
	if err != nil {
		return cst.CST[s.Language]{}, err
	}
	for id, pin := range pins {
		for i := range tree.Nodes[id-1].Attributes {
			if tree.Nodes[id-1].Attributes[i].Name == s.AttributePin {
				tree.Nodes[id-1].Attributes[i].Text = pin
			}
		}
	}
	return concrete.Parsed(tree, input, p.All, vocabulary)
}

func (p *parser) name() string {
	t := p.Take()
	name, valid := sourceName(t.Text)
	if !valid {
		p.Fail(text.SyntaxError, "expected identifier (escape reserved names with backticks)")
	}
	return name
}

func sourceName(name string) (string, bool) {
	escaped := strings.HasPrefix(name, "`") && strings.HasSuffix(name, "`") && len(name) >= 2
	if escaped {
		name = name[1 : len(name)-1]
	}
	return name, lit.Identifier(name) && (escaped || !lit.Reserved(name))
}

func (p *parser) moduleName() string {
	name := p.name()
	for p.Eat(".") && p.Err == nil {
		name += "." + p.name()
	}
	return name
}

func (p *parser) quoted() string {
	value, err := lit.Unquote(p.Take().Text)
	if err != nil {
		p.Fail(text.SyntaxError, err.Error())
	}
	return value
}

func (p *parser) module() ast.NodeID {
	start := p.Token().Start
	p.Want("module")
	name := p.moduleName()
	p.Want(";")
	var children []ast.Child
	imports, declarations := map[string]bool{}, map[string]bool{}
	for p.Peek() == "import" && p.Err == nil {
		begin := p.Take().Start
		alias := p.name()
		if imports[alias] {
			p.Fail(text.SyntaxError, "duplicate import alias")
		}
		imports[alias] = true
		p.Want("from")
		moduleName := p.moduleName()
		p.Want("pin")
		pin := p.quoted()
		if !strings.HasPrefix(pin, "sha256:") || len(pin) != 71 {
			p.Fail(text.SyntaxError, "expected sha256: followed by 64 hexadecimal digits")
		}
		pin = strings.TrimPrefix(pin, "sha256:")
		if _, err := hex.DecodeString(pin); err != nil {
			p.Fail(text.SyntaxError, "invalid SHA256 pin")
		}
		p.Want(";")
		children = append(children, tr.Child(s.RoleImport, p.Add(s.KindImport, begin, tr.Attrs(s.AttributeAlias, alias, s.AttributeName, moduleName, s.AttributePin, strings.ToLower(pin)))))
	}
	for !p.EOF() && p.Peek() != "module" && p.Err == nil {
		id := p.declaration()
		if p.Err != nil {
			return 0
		}
		declName := (tr.View[s.Language]{Tree: p.Tree}).Attr(id, s.AttributeName)
		if declarations[declName] {
			p.Fail(text.SyntaxError, "duplicate declaration name")
		}
		declarations[declName] = true
		children = append(children, tr.Child(s.RoleDeclaration, id))
	}
	return p.Add(s.KindModule, start, tr.Attrs(s.AttributeName, name, s.AttributePin, strings.Repeat("0", 64)), children...)
}

func (p *parser) declaration() ast.NodeID {
	start := p.Token().Start
	kindDefinition := p.Eat("kind")
	if !kindDefinition {
		p.Want("type")
	}
	name := p.name()
	p.parameters = map[string]int{}
	p.kindParameters = map[string]int{}
	p.groups = nil
	var children []ast.Child
	if p.Eat("<") {
		for p.Err == nil && !p.EOF() && p.Peek() != ">" {
			begin := p.Token().Start
			parameter := p.name()
			scope := p.parameters
			if kindDefinition {
				scope = p.kindParameters
			}
			if _, exists := scope[parameter]; exists {
				p.Fail(text.SyntaxError, "duplicate declaration parameter")
			}
			index := len(scope)
			scope[parameter] = index
			var contract ast.NodeID
			parameterKind, role := s.KindTypeParameter, s.RoleContract
			if kindDefinition {
				parameterKind, role = s.KindKindParameter, s.RoleSort
				if p.Eat(":") {
					contract = p.sortExpression()
				} else {
					contract = p.properSort(p.Token().Start)
				}
			} else if p.Eat(":") {
				contract = p.kindExpression()
			} else {
				contract = p.properTop(p.Token().Start)
			}
			children = append(children, tr.Child(s.RoleParameter, p.Add(parameterKind, begin, tr.Attrs(s.AttributeName, parameter, s.AttributeIndex, strconv.Itoa(index)), tr.Child(role, contract))))
			if !p.Eat(",") {
				break
			}
		}
		if len(children) == 0 {
			p.Fail(text.SyntaxError, "empty parameter list; omit angle brackets")
		}
		p.Want(">")
	}
	var result ast.NodeID
	if kindDefinition {
		if p.Eat(":") {
			result = p.sortExpression()
		} else {
			result = p.properSort(p.Token().Start)
		}
	} else if p.Eat(":") {
		result = p.kindExpression()
	} else {
		result = p.properTop(p.Token().Start)
	}
	children = append(children, tr.Child(s.RoleResult, result))
	p.Want("=")
	declarationKind := s.KindTypeDefinition
	var body ast.NodeID
	if kindDefinition {
		declarationKind = s.KindKindDefinition
		body = p.kindExpression()
	} else {
		body = p.expr()
	}
	children = append(children, tr.Child(s.RoleBody, body))
	p.Want(";")
	return p.Add(declarationKind, start, tr.Attrs(s.AttributeName, name), children...)
}

func (p *parser) descriptor(name string, start int, operands ...ast.NodeID) ast.NodeID {
	return p.Add(s.KindDescriptor, start, tr.Attrs(s.AttributeConstructor, name), tr.Children(s.RoleOperand, operands)...)
}

func (p *parser) expr() ast.NodeID {
	if !p.Enter() {
		return 0
	}
	defer p.Leave()
	start := p.Token().Start
	parts := []ast.NodeID{p.intersection()}
	for p.Err == nil && p.Eat("|") {
		parts = append(parts, p.intersection())
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return p.descriptor("union", start, parts...)
}

func (p *parser) intersection() ast.NodeID {
	start := p.Token().Start
	parts := []ast.NodeID{p.primary()}
	for p.Err == nil && p.Eat("&") {
		parts = append(parts, p.primary())
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return p.descriptor("and", start, parts...)
}

func (p *parser) arguments() []ast.NodeID {
	p.Want("<")
	var out []ast.NodeID
	for p.Err == nil && !p.EOF() && p.Peek() != ">" {
		out = append(out, p.expr())
		if !p.Eat(",") {
			break
		}
	}
	p.Want(">")
	return out
}

func (p *parser) primary() ast.NodeID {
	start := p.Token().Start
	callee := p.primaryAtom()
	for p.Err == nil && p.Peek() == "<" {
		arguments := p.arguments()
		if len(arguments) == 0 {
			p.Fail(text.SyntaxError, "empty application argument list; use a reference")
			return 0
		}
		children := []ast.Child{tr.Child(s.RoleCallee, callee)}
		children = append(children, tr.Children(s.RoleArgument, arguments)...)
		callee = p.Add(s.KindApplication, start, nil, children...)
	}
	return callee
}

func (p *parser) primaryAtom() ast.NodeID {
	if p.Err != nil {
		return 0
	}
	start := p.Token().Start
	if p.Eat("(") {
		id := p.expr()
		p.Want(")")
		if id != 0 {
			p.Tree.Nodes[id-1].Origin.Spans = []ast.Span{{Source: p.Input.SourceID, Start: uint64(start), End: uint64(p.Tokens[p.Pos-1].End)}}
		}
		return id
	}
	if p.Peek() == "{" || p.Peek() == "open" {
		return p.record()
	}
	// Hyphens are allowed only in the fixed constructor spellings, never names.
	word := p.Peek()
	if word == "atom" || word == "text" || word == "int" || word == "key" || word == "lowercase" {
		for i := 1; p.At(i) == "-" && p.Pos+i+1 < len(p.Tokens); i += 2 {
			if p.Tokens[p.Pos+i-1].End != p.Tokens[p.Pos+i].Start || p.Tokens[p.Pos+i].End != p.Tokens[p.Pos+i+1].Start {
				break
			}
			word += "-" + p.At(i+1)
		}
	}
	switch word {
	case "any", "atom":
		p.Take()
		return p.descriptor(word, start)
	case "text", "int", "bool":
		p.Take()
		kind := word
		if kind == "text" {
			kind = "utf8-text"
		}
		return p.descriptor("scalar", start, p.fragment(core.NewAtom([]byte(kind)), start))
	case "tuple", "list", "seq", "set", "map", "union", "and", "not", "wire":
		p.Take()
		args := p.arguments()
		arity := -1
		switch word {
		case "list", "seq", "set", "not", "wire":
			arity = 1
		case "map":
			arity = 2
		}
		if arity >= 0 && len(args) != arity {
			p.Fail(text.SyntaxError, "wrong constructor arity")
		}
		switch word {
		case "list":
			word = "list-of"
		case "set":
			word = "set-of"
		case "map":
			word = "map-of"
		}
		return p.descriptor(word, start, args...)
	case "exact", "atom-exact", "atom-in", "atom-len", "text-len", "int-min", "int-max", "key-grammar", "text-non-blank", "lowercase-uuid", "compound", "rec", "ref":
		for i, part := range strings.Split(word, "-") {
			if i > 0 {
				p.Want("-")
			}
			p.Want(part)
		}
		p.Want("(")
		var args []ast.NodeID
		switch word {
		case "key-grammar", "text-non-blank", "lowercase-uuid":
		case "atom-len", "text-len":
			args = append(args, p.literalNode())
			p.Want(",")
			args = append(args, p.literalNode())
		case "compound":
			args = append(args, p.literalNode())
			for p.Err == nil && p.Eat(",") {
				args = append(args, p.expr())
			}
		case "rec":
			args = append(args, p.literalNode())
			p.Want(",")
			args = append(args, p.expr())
		default:
			args = append(args, p.literalNode())
		}
		p.Want(")")
		return p.descriptor(word, start, args...)
	}
	local := p.Eat("self")
	if local {
		p.Want(".")
	}
	name := p.name()
	attrs := tr.Attrs(s.AttributeScope, "local", s.AttributeName, name)
	if !local && p.Eat(".") {
		attrs = tr.Attrs(s.AttributeScope, "imported", s.AttributeAlias, name, s.AttributeName, p.name())
	} else if !local {
		if index, ok := p.parameters[name]; ok {
			return p.Add(s.KindParameter, start, tr.Attrs(s.AttributeIndex, strconv.Itoa(index)))
		}
	}
	return p.Add(s.KindReference, start, attrs)
}

func (p *parser) record() ast.NodeID {
	start := p.Token().Start
	mode := "closed"
	if p.Eat("open") {
		mode = "open"
	}
	modeNode := p.fragment(core.NewAtom([]byte(mode)), start)
	p.Want("{")
	operands := []ast.NodeID{modeNode}
	for p.Err == nil && !p.EOF() && p.Peek() != "}" {
		begin := p.Token().Start
		name := ""
		if strings.HasPrefix(p.Peek(), "\"") {
			name = p.quoted()
		} else {
			name = p.name()
		}
		nameNode := p.fragment(data.MustEncodeText(name), begin)
		presence := "req"
		if p.Eat("?") {
			presence = "opt"
		}
		presenceNode := p.fragment(core.NewAtom([]byte(presence)), p.Token().Start)
		p.Want(":")
		body := p.expr()
		p.Want(";")
		label := p.AddSpan(s.KindAtom, begin, begin, tr.Attrs(s.AttributeHex, hex.EncodeToString([]byte("horos/field"))))
		operands = append(operands, p.Add(s.KindTuple, begin, nil, tr.Children(s.RoleElement, []ast.NodeID{label, nameNode, presenceNode, body})...))
	}
	p.Want("}")
	return p.descriptor("record", start, operands...)
}

func (p *parser) literalNode() ast.NodeID {
	start := p.Token().Start
	value := p.value()
	if p.Err != nil {
		return 0
	}
	return p.fragment(value, start)
}

func (p *parser) value() core.Value {
	if !p.Enter() {
		return nil
	}
	defer p.Leave()
	if p.Eat("#") {
		hexadecimal := p.Eat("hex")
		value := p.quoted()
		if hexadecimal {
			bytes, err := hex.DecodeString(value)
			if err != nil {
				p.Fail(text.SyntaxError, "invalid byte atom hex")
			}
			return core.NewAtom(bytes)
		}
		return core.NewAtom([]byte(value))
	}
	if strings.HasPrefix(p.Peek(), "\"") {
		return data.MustEncodeText(p.quoted())
	}
	switch p.Peek() {
	case "true", "false":
		return data.EncodeBool(p.Take().Text == "true")
	case "null":
		p.Take()
		return data.EncodeNull()
	case "tuple", "set":
		kind := p.Take().Text
		p.Want("(")
		var values []core.Value
		for p.Err == nil && !p.EOF() && p.Peek() != ")" {
			values = append(values, p.value())
			if !p.Eat(",") {
				break
			}
		}
		p.Want(")")
		if p.Err != nil {
			return nil
		}
		if kind == "tuple" {
			return core.NewTuple(values...)
		}
		value, err := data.EncodeSet(values)
		if err != nil {
			p.Fail(text.SyntaxError, err.Error())
		}
		return value
	}
	negative := p.Eat("-")
	spelling := p.Take().Text
	if negative {
		spelling = "-" + spelling
	}
	if !lit.Integer(spelling) || spelling == "-0" {
		p.Fail(text.SyntaxError, "expected literal value (decimal integers allow at most 4096 digits)")
		return nil
	}
	integer, _ := new(big.Int).SetString(spelling, 10)
	return data.EncodeInt(integer)
}

func (p *parser) fragment(value core.Value, start int) ast.NodeID {
	var visit func(core.Value) ast.NodeID
	visit = func(value core.Value) ast.NodeID {
		if !p.Enter() {
			return 0
		}
		defer p.Leave()
		if len(p.Tree.Nodes) >= structure.MaxNodes {
			p.Fail(text.LimitExceeded, "too many abstract nodes")
			return 0
		}
		switch v := value.(type) {
		case core.Atom:
			return p.AddSpan(s.KindAtom, start, start, tr.Attrs(s.AttributeHex, hex.EncodeToString(v.Bytes())))
		case core.Tuple:
			var children []ast.Child
			for i := 0; i < v.Len() && p.Err == nil; i++ {
				children = append(children, tr.Child(s.RoleElement, visit(v.At(i))))
			}
			return p.AddSpan(s.KindTuple, start, start, nil, children...)
		default:
			p.Fail(text.SyntaxError, "invalid literal value")
			return 0
		}
	}
	id := visit(value)
	if id != 0 {
		end := start
		if p.Pos > 0 && p.Tokens[p.Pos-1].End > end {
			end = p.Tokens[p.Pos-1].End
		}
		p.Tree.Nodes[id-1].Origin.Spans[0].End = uint64(end)
	}
	return id
}

// Provider supplies parsing independently of the composed Horos provider.
type Provider struct{ s.Language }

func (Provider) Parse(ctx context.Context, input parsecontract.Input) (cst.CST[s.Language], error) {
	return Parse(ctx, input)
}

var _ parsecontract.Provider[parsecontract.Input, cst.CST[s.Language]] = Provider{}
