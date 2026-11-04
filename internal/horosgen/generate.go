// Package horosgen reads established Horos terms for the native pagination profile.
// It is source-specific transformation machinery, not a second canonical model.
package horosgen

import (
	"context"
	"crypto/sha256"
	"fmt"
	"slices"
	"strings"

	analyze "github.com/Bitspark/bittype/analyze/horos/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	"github.com/Bitspark/bittype/semantics/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/bitspark/ontos/codec/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

const MaxNodes = 100000
const MaxDepth = 128

func Contract() transform.Contract {
	return transform.Contract{Evidence: []semantic.EvidenceRequirement{{Operation: analysis.Kind(), Profile: analyze.Profile, Reports: analyze.Contract().Semantic.ModelRequirements, Model: true}}}
}

// Arena allocates target syntax only. Its sticky error prevents partial output.
type Arena[L any] struct {
	Tree ast.AST[L]
	Ctx  context.Context
	Err  error
}

func (a *Arena[L]) Add(kind ast.Kind[L], attributes []ast.Attribute, origin ast.NodeRef, children ...ast.Child) ast.NodeID {
	if a.Err != nil {
		return 0
	}
	if err := a.Ctx.Err(); err != nil {
		a.Err = err
		return 0
	}
	if len(a.Tree.Nodes) >= MaxNodes {
		a.Err = &text.Diagnostic{Code: text.LimitExceeded, Message: "native generation node limit", Node: &origin}
		return 0
	}
	n := ast.Node[L]{Kind: kind, Attributes: attributes, Children: children}
	if origin.Tree != "" {
		n.Origin.Nodes = []ast.NodeRef{origin}
	}
	a.Tree.Nodes = append(a.Tree.Nodes, n)
	return ast.NodeID(len(a.Tree.Nodes))
}
func (a *Arena[L]) Set(id ast.NodeID, kind ast.Kind[L], attributes []ast.Attribute, children ...ast.Child) {
	if a.Err != nil {
		return
	}
	n := &a.Tree.Nodes[id-1]
	n.Kind = kind
	n.Attributes = attributes
	n.Children = children
}
func (a *Arena[L]) Target(id, target ast.NodeID) {
	if a.Err == nil {
		a.Tree.Nodes[id-1].References = []ast.Reference{{Role: "target", Target: ast.NodeRef{Tree: a.Tree.ID, Node: target}}}
	}
}
func (a *Arena[L]) Origin(id ast.NodeID) ast.NodeRef {
	if id > 0 && int(id) <= len(a.Tree.Nodes) && len(a.Tree.Nodes[id-1].Origin.Nodes) > 0 {
		return a.Tree.Nodes[id-1].Origin.Nodes[0]
	}
	return ast.NodeRef{}
}
func Attr(pairs ...string) []ast.Attribute {
	var out []ast.Attribute
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, ast.Attribute{Name: pairs[i], Text: pairs[i+1]})
	}
	return out
}
func Child(role string, id ast.NodeID) ast.Child { return ast.Child{Role: role, Node: id} }
func Children(role string, ids []ast.NodeID) []ast.Child {
	out := make([]ast.Child, len(ids))
	for i, id := range ids {
		out[i] = Child(role, id)
	}
	return out
}

// Field is an already mapped record member, never a canonical type expression.
type Field struct {
	Name     string
	Optional bool
	Type     ast.NodeID
	Origin   ast.NodeRef
}
type Dialect[L any] interface {
	ReservedNames() []string
	Parameter(*Arena[L], string, ast.NodeRef) ast.NodeID
	Reference(*Arena[L], string, ast.NodeID, []ast.NodeID, ast.NodeRef) ast.NodeID
	Text(*Arena[L], ast.NodeRef) ast.NodeID
	Boolean(*Arena[L], ast.NodeRef) ast.NodeID
	List(*Arena[L], ast.NodeID, ast.NodeRef) ast.NodeID
	Record(*Arena[L], ast.NodeID, string, []ast.NodeID, []Field)
	Alias(*Arena[L], ast.NodeID, string, []ast.NodeID, ast.NodeID)
	Finish(*Arena[L], []ast.NodeID) ast.NodeID
	FieldName(string) (string, error)
}

// Exported is deliberately conservative and collision-checked after conversion.
func Exported(name string) (string, error) {
	if !Identifier(name) {
		return "", fmt.Errorf("identifier %q is outside the pagination profile", name)
	}
	if name[0] == '_' {
		name = "Type" + name
	} else if name[0] >= 'a' && name[0] <= 'z' {
		name = string(name[0]-32) + name[1:]
	}
	if name == "Self" {
		name += "_"
	}
	return name, nil
}
func Identifier(name string) bool {
	if name == "" || strings.Contains(name, "__") {
		return false
	}
	for i, b := range []byte(name) {
		if !(b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b == '_' || i > 0 && b >= '0' && b <= '9') {
			return false
		}
	}
	return name != "_"
}

type generator[L any] struct {
	arena      Arena[L]
	dialect    Dialect[L]
	source     h.AST
	facts      semantics.Data
	interfaces map[semantics.Subject]semantics.TermID
	sites      map[semantics.TermID]ast.NodeRef
	names      map[semantics.EntityID]string
	targets    map[semantics.EntityID]ast.NodeID
	parameters map[semantics.EntityID][]ast.NodeID
	payloads   map[semantics.TermID]core.Tuple
	order      []semantics.EntityID
	state      map[semantics.EntityID]byte
	used       map[semantics.EntityID]bool
	work       int
}

// Generate bundles the checked program closure into one source unit. Imported
// declarations are qualified by their module pin; no native import resolution runs.
func Generate[L any](ctx context.Context, request transform.Request[h.Language], id ast.TreeID, dialect Dialect[L]) (ast.AST[L], error) {
	fail := func(err error) (ast.AST[L], error) { return ast.AST[L]{}, err }
	source := request.Source()
	if id == "" || id == source.ID || request.Entry() != (ast.NodeRef{Tree: source.ID, Node: source.Root}) {
		return fail(semantic.ErrInvalidRequest)
	}
	evidence, ok := request.Evidence(0)
	if !ok {
		return fail(semantic.ErrEvidenceUnavailable)
	}
	model, ok := evidence.Model()
	if !ok {
		return fail(semantic.ErrEvidenceUnavailable)
	}
	g := generator[L]{arena: Arena[L]{Tree: ast.AST[L]{ID: id}, Ctx: ctx}, dialect: dialect, source: source, facts: model.Data(), sites: map[semantics.TermID]ast.NodeRef{}, interfaces: map[semantics.Subject]semantics.TermID{}, names: map[semantics.EntityID]string{}, targets: map[semantics.EntityID]ast.NodeID{}, parameters: map[semantics.EntityID][]ast.NodeID{}, payloads: map[semantics.TermID]core.Tuple{}, state: map[semantics.EntityID]byte{}, used: map[semantics.EntityID]bool{}}
	for _, e := range g.facts.Expressions {
		if e.At.Role == resolve.ExpressionRole() {
			g.sites[e.Term] = e.At.Site
		}
	}
	for _, fact := range g.facts.Facts {
		if fact.Relation == resolve.InterfaceRelation() {
			if g.interfaces[fact.Subject] != 0 {
				return fail(semantic.ErrInvalidRequest)
			}
			g.interfaces[fact.Subject] = fact.Term
		}
	}
	if err := g.prepare(); err != nil {
		return fail(err)
	}
	var declarations []ast.NodeID
	for _, entityID := range g.order {
		entity := g.facts.Entities[entityID-1]
		target := g.targets[entityID]
		params := g.parameters[entityID]
		term := g.facts.Terms[entity.Body-1]
		if term.Operator == (semantics.Key{Owner: "horos", Name: "descriptor/record"}) {
			fields, err := g.fields(entity.Body)
			if err != nil {
				return fail(err)
			}
			dialect.Record(&g.arena, target, g.names[entityID], params, fields)
		} else {
			body, err := g.expression(entity.Body, 0)
			if err != nil {
				return fail(err)
			}
			dialect.Alias(&g.arena, target, g.names[entityID], params, body)
		}
		for _, p := range entity.Parameters {
			if !g.used[p.Entity] {
				return fail(g.unsupported(entity.Body, "unused generic parameters require a separate representation profile"))
			}
		}
		declarations = append(declarations, target)
	}
	g.arena.Tree.Root = dialect.Finish(&g.arena, declarations)
	if g.arena.Err != nil {
		return fail(g.arena.Err)
	}
	if err := ctx.Err(); err != nil {
		return fail(err)
	}
	return g.arena.Tree, nil
}

func (g *generator[L]) unsupported(id semantics.TermID, message string) error {
	site := g.sites[id]
	return &text.Diagnostic{Code: text.Unsupported, Message: message, Node: &site}
}

// The checked Horos interface facts distinguish proper types from constructors.
// This reader selects a native profile; it does not infer or check those facts.
func (g *generator[L]) constructor(subject semantics.Subject) (bool, error) {
	id := g.interfaces[subject]
	if id == 0 {
		return false, semantic.ErrInvalidRequest
	}
	term := g.facts.Terms[id-1]
	if term.Kind != semantics.ProviderExpression || term.Encoding != (semantics.Key{Owner: "ontos", Name: "codec-v1"}) {
		return false, semantic.ErrInvalidRequest
	}
	value, err := codec.Decode(term.Data)
	if err != nil {
		return false, err
	}
	if claim := atom(value); claim != "type" && claim != "data" {
		return false, semantic.ErrInvalidRequest
	}
	switch term.Operator {
	case semantics.Key{Owner: "horos", Name: "interface/proper"}:
		if len(term.Operands) != 0 {
			return false, semantic.ErrInvalidRequest
		}
		return false, nil
	case semantics.Key{Owner: "horos", Name: "interface/constructor"}:
		if len(term.Operands) == 0 {
			return false, semantic.ErrInvalidRequest
		}
		return true, nil
	default:
		return false, semantic.ErrInvalidRequest
	}
}

func (g *generator[L]) prepare() error {
	root := g.source.Nodes[g.source.Root-1]
	var rootModule ast.NodeRef
	for _, r := range root.References {
		if r.Role == h.RoleRoot {
			rootModule = r.Target
		}
	}
	if rootModule.Tree == "" {
		return semantic.ErrInvalidRequest
	}
	occupied := map[string]bool{}
	for _, name := range g.dialect.ReservedNames() {
		occupied[name] = true
	}
	var declarations []semantics.EntityID
	for i, e := range g.facts.Entities {
		if e.Kind != (semantics.Key{Owner: "horos", Name: "declaration"}) {
			continue
		}
		id := semantics.EntityID(i + 1)
		name, err := Exported(e.Name)
		if err != nil {
			return g.unsupported(e.Body, err.Error())
		}
		owner := g.facts.Entities[e.Owner-1]
		if len(owner.Declarations) != 1 || len(e.Declarations) != 1 {
			return semantic.ErrInvalidRequest
		}
		site := owner.Declarations[0]
		if site.Tree != g.source.ID || site.Node == 0 || int(site.Node) > len(g.source.Nodes) {
			return semantic.ErrInvalidRequest
		}
		if site != rootModule {
			var pin string
			for _, a := range g.source.Nodes[site.Node-1].Attributes {
				if a.Name == h.AttributePin {
					pin = a.Text
				}
			}
			if pin == "" {
				return semantic.ErrInvalidRequest
			}
			sum := sha256.Sum256([]byte(pin))
			name = fmt.Sprintf("M%x_%s", sum[:8], name)
		}
		if occupied[name] {
			return g.unsupported(e.Body, "native declaration name collision: "+name)
		}
		occupied[name] = true
		g.names[id] = name
		g.targets[id] = g.arena.Add("", nil, e.Declarations[0])
		declarations = append(declarations, id)
	}
	// Parameter names cannot shadow any generated declaration or support type.
	for _, id := range declarations {
		e := g.facts.Entities[id-1]
		local := map[string]bool{}
		for _, p := range e.Parameters {
			pe := g.facts.Entities[p.Entity-1]
			if len(pe.Declarations) != 1 {
				return semantic.ErrInvalidRequest
			}
			constructor, err := g.constructor(semantics.Subject{Entity: p.Entity})
			if err != nil {
				return err
			}
			if constructor {
				site := pe.Declarations[0]
				return &text.Diagnostic{Code: text.Unsupported, Message: "constructor parameters require a separate native representation profile", Node: &site}
			}
			name, err := Exported(pe.Name)
			if err != nil {
				return g.unsupported(e.Body, err.Error())
			}
			if occupied[name] || local[name] {
				return g.unsupported(e.Body, "native generic name collision: "+name)
			}
			local[name] = true
			g.names[p.Entity] = name
			g.targets[p.Entity] = g.dialect.Parameter(&g.arena, name, pe.Declarations[0])
			g.parameters[id] = append(g.parameters[id], g.targets[p.Entity])
		}
	}
	for _, id := range declarations {
		if err := g.visitDeclaration(id, 0); err != nil {
			return err
		}
	}
	return g.arena.Err
}
func (g *generator[L]) tick(depth int) error {
	if err := g.arena.Ctx.Err(); err != nil {
		return err
	}
	g.work++
	if depth > MaxDepth || g.work > MaxNodes*8 {
		return &text.Diagnostic{Code: text.LimitExceeded, Message: "native generation traversal limit"}
	}
	return nil
}
func (g *generator[L]) visitDeclaration(id semantics.EntityID, depth int) error {
	if err := g.tick(depth); err != nil {
		return err
	}
	if g.state[id] == 2 {
		return nil
	}
	e := g.facts.Entities[id-1]
	if g.state[id] == 1 {
		return g.unsupported(e.Body, "recursive declarations are outside the pagination profile")
	}
	g.state[id] = 1
	var visit func(semantics.TermID, int) error
	visit = func(id semantics.TermID, level int) error {
		if err := g.tick(level); err != nil {
			return err
		}
		term := g.facts.Terms[id-1]
		if term.Kind == semantics.Reference {
			if g.targets[term.Target] == 0 {
				return semantic.ErrInvalidRequest
			}
			switch g.facts.Entities[term.Target-1].Kind {
			case semantics.Key{Owner: "horos", Name: "declaration"}:
				if err := g.visitDeclaration(term.Target, level+1); err != nil {
					return err
				}
			case semantics.Key{Owner: "horos", Name: "parameter"}:
			default:
				return semantic.ErrInvalidRequest
			}
		}
		if term.Kind == semantics.Application {
			if err := visit(term.Callee, level+1); err != nil {
				return err
			}
		}
		for _, arg := range term.Arguments {
			if err := visit(arg, level+1); err != nil {
				return err
			}
		}
		for _, op := range term.Operands {
			if err := visit(op.Term, level+1); err != nil {
				return err
			}
		}
		return nil
	}
	if err := visit(e.Body, depth+1); err != nil {
		return err
	}
	g.state[id] = 2
	g.order = append(g.order, id)
	return nil
}
func (g *generator[L]) payload(id semantics.TermID) (core.Tuple, error) {
	if value, ok := g.payloads[id]; ok {
		return value, nil
	}
	term := g.facts.Terms[id-1]
	if term.Encoding != (semantics.Key{Owner: "ontos", Name: "codec-v1"}) {
		return core.Tuple{}, semantic.ErrInvalidRequest
	}
	value, err := codec.Decode(term.Data)
	if err != nil {
		return core.Tuple{}, err
	}
	tuple, ok := value.(core.Tuple)
	if !ok {
		return core.Tuple{}, semantic.ErrInvalidRequest
	}
	g.payloads[id] = tuple
	return tuple, nil
}
func atom(value core.Value) string {
	a, ok := value.(core.Atom)
	if !ok {
		return ""
	}
	return string(a.Bytes())
}
func (g *generator[L]) operand(id semantics.TermID, path ...uint32) (semantics.TermID, error) {
	for _, op := range g.facts.Terms[id-1].Operands {
		if slices.Equal(op.Locator, path) {
			return op.Term, nil
		}
	}
	return 0, semantic.ErrInvalidRequest
}
func (g *generator[L]) expression(id semantics.TermID, depth int) (ast.NodeID, error) {
	if err := g.tick(depth); err != nil {
		return 0, err
	}
	term := g.facts.Terms[id-1]
	origin := g.sites[id]
	constructor, err := g.constructor(semantics.Subject{Term: id})
	if err != nil {
		return 0, err
	}
	if constructor {
		return 0, g.unsupported(id, "constructor references require a separate native representation profile")
	}
	if term.Kind == semantics.Reference || term.Kind == semantics.Application {
		entity := term.Target
		if term.Kind == semantics.Application {
			callee := g.facts.Terms[term.Callee-1]
			if callee.Kind != semantics.Reference {
				return 0, g.unsupported(term.Callee, "only named declaration callees are supported by the pagination profile")
			}
			if g.facts.Entities[callee.Target-1].Kind != (semantics.Key{Owner: "horos", Name: "declaration"}) {
				return 0, g.unsupported(term.Callee, "constructor parameter callees require a separate native representation profile")
			}
			entity = callee.Target
		}
		target := g.targets[entity]
		if target == 0 {
			return 0, semantic.ErrInvalidRequest
		}
		g.used[entity] = true
		var arguments []ast.NodeID
		for _, argument := range term.Arguments {
			mapped, err := g.expression(argument, depth+1)
			if err != nil {
				return 0, err
			}
			arguments = append(arguments, mapped)
		}
		return g.dialect.Reference(&g.arena, g.names[entity], target, arguments, origin), g.arena.Err
	}
	if term.Operator.Owner != "horos" {
		return 0, g.unsupported(id, "unrecognized semantic operator")
	}
	switch term.Operator.Name {
	case "descriptor/scalar":
		tuple, err := g.payload(id)
		if err != nil {
			return 0, err
		}
		switch atom(tuple.At(1)) {
		case "utf8-text":
			return g.dialect.Text(&g.arena, origin), g.arena.Err
		case "bool":
			return g.dialect.Boolean(&g.arena, origin), g.arena.Err
		}
	case "descriptor/list-of":
		element, err := g.operand(id, 1)
		if err != nil {
			return 0, err
		}
		mapped, err := g.expression(element, depth+1)
		if err != nil {
			return 0, err
		}
		return g.dialect.List(&g.arena, mapped, origin), g.arena.Err
	}
	return 0, g.unsupported(id, "pagination mapping does not represent "+term.Operator.Name+" here; nested records must be named declarations")
}
func (g *generator[L]) fields(id semantics.TermID) ([]Field, error) {
	tuple, err := g.payload(id)
	if err != nil {
		return nil, err
	}
	if atom(tuple.At(1)) != "closed" {
		return nil, g.unsupported(id, "open records require an explicit native representation")
	}
	var fields []Field
	names := map[string]bool{}
	operands := map[uint32]semantics.TermID{}
	for _, op := range g.facts.Terms[id-1].Operands {
		if len(op.Locator) == 2 && op.Locator[1] == 3 {
			operands[op.Locator[0]] = op.Term
		}
	}
	for i := 2; i < tuple.Len(); i++ {
		if err := g.tick(0); err != nil {
			return nil, err
		}
		field, ok := tuple.At(i).(core.Tuple)
		if !ok {
			return nil, semantic.ErrInvalidRequest
		}
		name, err := data.ReadText(field.At(1))
		if err != nil {
			return nil, err
		}
		name, err = g.dialect.FieldName(name)
		if err != nil {
			return nil, g.unsupported(id, err.Error())
		}
		if names[name] {
			return nil, g.unsupported(id, "native field name collision: "+name)
		}
		names[name] = true
		child := operands[uint32(i)]
		if child == 0 {
			return nil, semantic.ErrInvalidRequest
		}
		mapped, err := g.expression(child, 0)
		if err != nil {
			return nil, err
		}
		fields = append(fields, Field{Name: name, Optional: atom(field.At(2)) == "opt", Type: mapped, Origin: g.sites[child]})
	}
	return fields, nil
}
