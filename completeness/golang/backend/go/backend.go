// Package backend is the completeness test backend for Go (prototype,
// bittype#196, rows 29 and 31). It lowers a conversion to Go code from the
// serialized plan, the source snapshot and the facets of the selected
// contributions, and from nothing else: it never holds an object of the
// planner. It is a test instrument and no Go backend: "A Go backend belongs
// to the future Go runtime library."
//
// The kernel concepts are lowered here, once, and no family has a case here:
// a planned declaration becomes a function that takes one adapter for each
// binder, an application becomes a call of it with the adapter arguments, and
// a binder becomes the function's parameter. Everything a family decides is
// lowered by that family's facet.
package backend

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/Bitspark/bittype/completeness/golang/values/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
)

// Descriptor is the source's description of one type, as the snapshot of the
// source semantics carries it (row 34). Position is the source's own label of
// a child, the same that a plan's link carries.
type Descriptor struct {
	Kind        string       `json:"kind"` // record | list | text | bool | reference | parameter
	Fields      []Field      `json:"fields,omitempty"`
	Element     *Child       `json:"element,omitempty"`
	Declaration string       `json:"declaration,omitempty"`
	Arguments   []Descriptor `json:"arguments,omitempty"`
	Index       int          `json:"index,omitempty"`
}
type Field struct {
	Name     string     `json:"name"`
	Optional bool       `json:"optional"`
	Position []uint32   `json:"position"`
	Type     Descriptor `json:"type"`
}
type Child struct {
	Position []uint32   `json:"position"`
	Type     Descriptor `json:"type"`
}
type Module struct {
	Name string `json:"name"`
	Pin  string `json:"pin"`
}

// Snapshot is the source-semantics input of a reader in serialized form: the
// module pins and the descriptor of each declaration, by source identity.
type Snapshot struct {
	Semantics    string                `json:"semantics"`
	Modules      []Module              `json:"modules"`
	Declarations map[string]Descriptor `json:"declarations"`
}

// Input is everything the backend reads.
type Input struct {
	Plan     []byte
	Snapshot []byte
	Facets   []g.Facet
	Codecs   []g.PayloadCodec
	// Semantics are the source interpretations this reader knows.
	Semantics []string
}

type scope struct {
	node    g.PlanNode
	binders []string
}

type lowerer struct {
	plan      g.RepresentationPlan
	in        Input
	snapshot  Snapshot
	functions map[string]string // declaration node -> function name
	helpers   map[string]string // helper key -> spelled name
	body      strings.Builder
}

// Link checks that a snapshot belongs to a plan: every source declaration of
// the plan is in it. A module of the same name under another pin is a linking
// error, and no snapshot at all is an explicit failure.
func Link(plan g.RepresentationPlan, serialized []byte) (Snapshot, error) {
	var snapshot Snapshot
	if len(serialized) == 0 {
		return snapshot, fmt.Errorf("%w: a reader needs the source snapshot, and none was given", g.ErrIncomplete)
	}
	if err := json.Unmarshal(serialized, &snapshot); err != nil {
		return snapshot, fmt.Errorf("%w: source snapshot: %v", g.ErrInvalid, err)
	}
	if snapshot.Semantics != plan.Source().Semantics {
		return snapshot, fmt.Errorf("%w: the snapshot is of source semantics %q and the plan of %q", g.ErrLinking, snapshot.Semantics, plan.Source().Semantics)
	}
	for _, s := range plan.Symbols() {
		identity, ok := strings.CutPrefix(s.ID, "source:")
		if !ok {
			continue
		}
		if _, ok := snapshot.Declarations[identity]; !ok {
			return snapshot, fmt.Errorf("%w: the snapshot has no declaration %s: it is of another module pin", g.ErrLinking, identity)
		}
	}
	return snapshot, nil
}

// Lower returns the Go source of the conversion functions of the plan's
// roots, for the package main of the completeness program. The program's
// runtime is Runtime.
func Lower(in Input) (string, error) {
	plan, err := g.ParsePlan(in.Plan)
	if err != nil {
		return "", err
	}
	knowledge := g.ReaderKnowledge{SourceSemantics: in.Semantics}
	for _, codec := range in.Codecs {
		knowledge.Schemas = append(knowledge.Schemas, codec.Schema())
	}
	for _, c := range plan.Contributions() {
		// The selected contributions are an input of a backend: it knows
		// those it has a facet for, by their whole reference.
		for _, f := range in.Facets {
			if f.Contribution() == c.Reference.Identity {
				knowledge.Contributions = append(knowledge.Contributions, c.Reference)
				break
			}
		}
	}
	if err := plan.CheckKnown(knowledge); err != nil {
		return "", err
	}
	snapshot, err := Link(plan, in.Snapshot)
	if err != nil {
		return "", err
	}
	l := &lowerer{plan: plan, in: in, snapshot: snapshot, functions: map[string]string{}, helpers: map[string]string{}}
	var roots []string
	for _, root := range plan.Roots() {
		symbol, _ := plan.Symbol(root)
		node, _ := plan.Node(symbol.Node)
		if node.Binders != 0 {
			continue // an open declaration is lowered when an application reaches it
		}
		function, err := l.declaration(node)
		if err != nil {
			return "", err
		}
		roots = append(roots, fmt.Sprintf("\t%q: root(%s()),\n", root[strings.LastIndex(root, "/")+1:], function))
	}
	sort.Strings(roots)
	return "// Code lowered from the serialized plan by the completeness test backend. DO NOT EDIT.\n\npackage main\n\nimport \"fmt\"\n\nvar _ = fmt.Sprint\n\nvar roots = map[string]func(Value) (Value, error){\n" + strings.Join(roots, "") + "}\n" + l.body.String(), nil
}

func (l *lowerer) descriptor(node g.PlanNode) *Descriptor {
	if d, ok := l.snapshot.Declarations[node.Source]; ok {
		return &d
	}
	return nil
}

// declaration emits the adapter function of a planned declaration, once.
func (l *lowerer) declaration(node g.PlanNode) (string, error) {
	if name, ok := l.functions[node.ID]; ok {
		return name, nil
	}
	symbol, ok := l.plan.Symbol(node.Symbol)
	if !ok || node.Kind != "declaration" {
		return "", fmt.Errorf("%w: node %s is no planned declaration", g.ErrScope, node.ID)
	}
	name := "adapt" + symbol.Intent
	l.functions[node.ID] = name
	s := scope{node: node}
	var formals, parameters []string
	for i := 0; i < node.Binders; i++ {
		s.binders = append(s.binders, fmt.Sprintf("B%d", i))
		formals = append(formals, fmt.Sprintf("B%d any", i))
		parameters = append(parameters, fmt.Sprintf("b%d Adapter[B%d]", i, i))
	}
	generic := ""
	if len(formals) > 0 {
		generic = "[" + strings.Join(formals, ", ") + "]"
	}
	code, err := l.facet(node, s, false, l.descriptor(node))
	if err != nil {
		return "", err
	}
	fmt.Fprintf(&l.body, "\nfunc %s%s(%s) Adapter[%s] {\n\treturn %s\n}\n", name, generic, strings.Join(parameters, ", "), l.self(symbol, s), code)
	return name, nil
}

// self is the Go type of a declaration inside its own adapter function. Go
// has no Files stage for generated declarations (row 24a): the name is the
// spelled intent of the plan.
func (l *lowerer) self(symbol g.PlanSymbol, s scope) string {
	if len(s.binders) == 0 {
		return symbol.Intent
	}
	return symbol.Intent + "[" + strings.Join(s.binders, ", ") + "]"
}

// expression lowers the node of one dependency in the scope of a declaration.
func (l *lowerer) expression(id string, s scope, declared bool, descriptor *Descriptor) (string, error) {
	node, ok := l.plan.Node(id)
	if !ok {
		return "", fmt.Errorf("%w: no node %s", g.ErrIncomplete, id)
	}
	switch node.Kind {
	case "binder":
		// A generic adapter argument: the adapter that the binder brings into
		// the scope of its declaration. The plan checks the scope.
		reference, err := l.plan.AdapterOf(s.node.ID, node.Binder.Index)
		if err != nil {
			return "", err
		}
		if err := l.plan.Check(reference, g.AdapterReference, node.Binder.Scope); err != nil {
			return "", err
		}
		return fmt.Sprintf("b%d", reference.Index()), nil
	case "application":
		callee, _ := l.plan.Node(node.Application.Declaration)
		function, err := l.declaration(callee)
		if err != nil {
			return "", err
		}
		arguments := make([]string, callee.Binders)
		for _, a := range node.Application.Adapters {
			var argument *Descriptor
			if descriptor != nil && a.Binder < len(descriptor.Arguments) {
				argument = &descriptor.Arguments[a.Binder]
			}
			if arguments[a.Binder], err = l.expression(a.Node, s, false, argument); err != nil {
				return "", err
			}
		}
		for i, a := range arguments {
			if a == "" {
				return "", fmt.Errorf("%w: the application %s has no adapter argument for binder %d: the plan was not asked for bindings", g.ErrIncomplete, node.Source, i)
			}
		}
		return function + "(" + strings.Join(arguments, ", ") + ")", nil
	case "use", "field":
		return l.facet(node, s, declared, descriptor)
	}
	return "", fmt.Errorf("%w: a %s node has no conversion", g.ErrUnsupported, node.Kind)
}

type context struct {
	l          *lowerer
	node       g.PlanNode
	payload    any
	descriptor *Descriptor
}

func (c context) Payload() any        { return c.payload }
func (c context) Links() int          { return len(c.node.Links) }
func (c context) Stated(i int) string { return c.node.Links[i].Stated }
func (c context) Fragment(i int) bool { return c.node.Links[i].Kind == "field" }
func (c context) Member(i int) (g.PlanReference, error) {
	if len(c.node.Members) != len(c.node.Links) {
		return g.PlanReference{}, fmt.Errorf("%w: the backend lowers a record with one member for each link", g.ErrUnsupported)
	}
	return c.l.plan.MemberOf(c.node.ID, i)
}
func (c context) Owned(name string) (g.PlanReference, error) {
	return c.l.plan.DeclarationOf("owned:" + c.node.Contribution.String() + "/" + name)
}
func (c context) Canonical(i int) (string, error) {
	if c.descriptor == nil || c.descriptor.Kind != "record" {
		return "", fmt.Errorf("%w: the source snapshot does not describe node %s as a record", g.ErrLinking, c.node.ID)
	}
	for _, f := range c.descriptor.Fields {
		if fmt.Sprint(f.Position) == fmt.Sprint(c.node.Links[i].Position) {
			return f.Name, nil
		}
	}
	return "", fmt.Errorf("%w: the source snapshot has no child at position %v", g.ErrLinking, c.node.Links[i].Position)
}

// child is the source's description of the child at link i.
func (c context) child(i int) *Descriptor {
	if c.descriptor == nil {
		return nil
	}
	position := fmt.Sprint(c.node.Links[i].Position)
	switch c.descriptor.Kind {
	case "record":
		for _, f := range c.descriptor.Fields {
			if fmt.Sprint(f.Position) == position {
				return &f.Type
			}
		}
	case "list":
		if c.descriptor.Element != nil {
			return &c.descriptor.Element.Type
		}
	}
	return c.descriptor
}

// facet lowers one node with the lowering facet of its contribution.
func (l *lowerer) facet(node g.PlanNode, s scope, declared bool, descriptor *Descriptor) (string, error) {
	kind := values.ValueLowering
	if node.Kind == "field" {
		kind = values.FieldStateLowering
	}
	var lowering *values.Lowering
	for _, f := range l.in.Facets {
		if f.Contribution() != node.Contribution.Identity {
			continue
		}
		if value, ok := kind.Read(f); ok {
			v := value.(values.Lowering)
			lowering = &v
		}
	}
	if lowering == nil {
		return "", fmt.Errorf("%w: no %s facet for contribution %s", g.ErrMissingFacet, kind.Reference().Identity, node.Contribution)
	}
	c := context{l: l, node: node, descriptor: descriptor}
	// The backend reads a payload only through the codec of its schema, after
	// it has checked the receipt (guarantee 6).
	frozen, err := hex.DecodeString(node.Payload)
	if err != nil {
		return "", err
	}
	if err := g.CheckReceipt(node.Receipt, *node.Schema, frozen); err != nil {
		return "", err
	}
	for _, codec := range l.in.Codecs {
		if codec.Schema() == *node.Schema {
			if c.payload, err = codec.View(frozen); err != nil {
				return "", err
			}
		}
	}
	if lowering.Helpers != nil {
		helpers, err := lowering.Helpers(c)
		if err != nil {
			return "", err
		}
		for _, h := range helpers {
			key := node.Contribution.String() + "/" + h.Name
			if _, done := l.helpers[key]; done {
				continue
			}
			l.helpers[key] = fmt.Sprintf("helper%d_%s", len(l.helpers), h.Name)
			text, err := l.render(h.Code, c, s, declared)
			if err != nil {
				return "", err
			}
			l.body.WriteString("\n" + text + "\n")
		}
	}
	code, err := lowering.Lower(c)
	if err != nil {
		return "", err
	}
	return l.render(code, c, s, declared)
}

// render fills the holes of lowered code. Every reference is checked by the
// plan for its plan, its category and its scope before it is resolved.
func (l *lowerer) render(code values.Code, c context, s scope, declared bool) (string, error) {
	text := code.Text
	for i, hole := range code.Holes {
		var filled string
		var err error
		switch {
		case hole.Dependency > 0:
			link := hole.Dependency - 1
			if link >= len(c.node.Links) {
				return "", fmt.Errorf("%w: no dependency %d", g.ErrScope, link)
			}
			child := c.child(link)
			if c.node.Kind == "declaration" {
				child = c.descriptor
			}
			filled, err = l.expression(c.node.Links[link].Node, s, hole.Declared, child)
		case hole.Helper != "":
			name, ok := l.helpers[c.node.Contribution.String()+"/"+hole.Helper]
			if !ok {
				return "", fmt.Errorf("%w: no helper %s", g.ErrScope, hole.Helper)
			}
			filled = name
		case hole.SelfType:
			if !declared {
				return "", fmt.Errorf("%w: the backend lowers a record only as the body of a nominal declaration", g.ErrUnsupported)
			}
			symbol, _ := l.plan.Symbol(s.node.Symbol)
			filled = l.self(symbol, s)
		default:
			if err = l.plan.Check(hole.Reference, hole.Want, s.node.Symbol); err != nil {
				return "", err
			}
			switch hole.Want {
			case g.MemberReference:
				filled, err = l.plan.Member(hole.Reference)
			case g.OwnedReference, g.DeclarationReference:
				for _, symbol := range l.plan.Symbols() {
					if symbol.Node == hole.Reference.Node() {
						filled = symbol.Intent
					}
				}
			default:
				err = fmt.Errorf("%w: the Go vocabulary has no hole for a %s reference", g.ErrUnsupported, hole.Want)
			}
		}
		if err != nil {
			return "", err
		}
		text = strings.ReplaceAll(text, fmt.Sprintf("«%d»", i), filled)
	}
	return text, nil
}

// Runtime is the fixed part of the completeness program: the abstract value,
// the adapter, the field-state protocol and the driver. It knows no family.
const Runtime = `package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Value is an abstract value: text is a string, a boolean a bool, a list a
// []any and a record a map[string]any in which an absent field has no key.
type Value = any

// Adapter converts between abstract values and native values of N.
type Adapter[N any] struct {
	Decode func(Value) (N, error)
	Encode func(N) (Value, error)
}

// FieldState is the field-state protocol on the Go side. A record's adapter
// invokes it for one field whose presence it delegates. Construct builds the
// native member from the state of the field: whether it is present, and its
// abstract value if it is. State reads the native member back into that state.
type FieldState[M any] struct {
	Construct func(present bool, value Value) (M, error)
	State     func(member M) (present bool, value Value, err error)
}

// Fields reads an abstract record and refuses a field outside names.
func Fields(v Value, names ...string) (map[string]Value, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("structural mismatch: a record is expected, got %T", v)
	}
	known := map[string]bool{}
	for _, name := range names {
		known[name] = true
	}
	for name := range m {
		if !known[name] {
			return nil, fmt.Errorf("structural mismatch: the record has no field %q", name)
		}
	}
	return m, nil
}

func root[N any](a Adapter[N]) func(Value) (Value, error) {
	return func(v Value) (Value, error) {
		n, err := a.Decode(v)
		if err != nil {
			return nil, err
		}
		return a.Encode(n)
	}
}

type vector struct {
	Root  string ` + "`json:\"root\"`" + `
	Value Value  ` + "`json:\"value\"`" + `
}
type result struct {
	Root  string ` + "`json:\"root\"`" + `
	Value Value  ` + "`json:\"value\"`" + `
	Error string ` + "`json:\"error,omitempty\"`" + `
}

func main() {
	var vectors []vector
	if err := json.NewDecoder(os.Stdin).Decode(&vectors); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	results := []result{}
	for _, v := range vectors {
		convert, ok := roots[v.Root]
		if !ok {
			results = append(results, result{Root: v.Root, Error: "no such root"})
			continue
		}
		out, err := convert(v.Value)
		if err != nil {
			results = append(results, result{Root: v.Root, Error: err.Error()})
			continue
		}
		results = append(results, result{Root: v.Root, Value: out})
	}
	json.NewEncoder(os.Stdout).Encode(results)
}
`
