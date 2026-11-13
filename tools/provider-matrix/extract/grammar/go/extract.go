// Package grammar interprets the declaration format of the agnostic syntax
// model: a grammar table written as a composite literal of its table type.
// It reads initializers without importing or executing inspected packages.
package grammar

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"strconv"

	c "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
)

// owner is the package of the table type and of the shared predicate names.
const owner = "github.com/Bitspark/bittype/model/syntax/go"

// Format identifies this interpretation, independently of each selected root.
const Format = owner + "#Rules"

type Extractor struct{}

func (Extractor) Selectors() []c.Selector {
	symbol := c.Symbol{Package: owner, Name: "Rules"}
	return []c.Selector{{NamedType: &symbol}}
}

func (Extractor) Extract(ctx context.Context, request c.DefinitionRequest) (c.Inventory, error) {
	limits, err := c.EffectiveLimits(request.Limits)
	if err != nil {
		return c.Inventory{}, err
	}
	if ctx == nil || request.Source == nil {
		return c.Inventory{}, fmt.Errorf("missing extraction context or source")
	}
	if err := ctx.Err(); err != nil {
		return c.Inventory{}, err
	}
	root, err := request.Source.Declaration(request.Root)
	if err != nil {
		return c.Inventory{}, err
	}
	if root.Symbol != request.Root {
		return c.Inventory{}, fmt.Errorf("source returned another definition")
	}
	e := evaluator{ctx: ctx, request: request, limits: limits, files: map[c.FileRef]*parsed{}, active: map[visit]bool{}}
	if root.Value.File == "" {
		result := c.Inventory{Extent: c.PartlyDecoded, Limitations: []c.Limitation{{At: root.Span, Message: "Selected definition has no initializer to interpret"}}}
		if err := c.Validate(request, result); err != nil {
			return c.Inventory{}, err
		}
		return result, nil
	}
	x := e.expression(root.Value)
	if e.err != nil {
		return c.Inventory{}, e.err
	}
	facts := e.facts(x)
	if facts.Type.Origin == nil || *facts.Type.Origin != (c.Symbol{Package: owner, Name: "Rules"}) {
		return c.Inventory{}, fmt.Errorf("selected initializer is not the grammar Rules type")
	}
	e.result.Extent = c.Decoded
	e.rules(x, 0)
	if e.err != nil {
		return c.Inventory{}, e.err
	}
	// A later map assignment may have replaced an opaque earlier value. Extent
	// describes the retained result, while unresolved keys remain localized on
	// their collection and cannot be erased by a known-name overwrite.
	if e.uncertain && hasUncertainty(e.result) {
		e.result.Extent = c.PartlyDecoded
	}
	if err := c.Validate(request, e.result); err != nil {
		return c.Inventory{}, err
	}
	return e.result, nil
}

func hasUncertainty(inventory c.Inventory) bool {
	if len(inventory.Limitations) > 0 {
		return true
	}
	for _, feature := range inventory.Features {
		if feature.Roles.State != c.Known || feature.Attributes.State != c.Known || feature.References.State != c.Known {
			return true
		}
		for _, role := range feature.Roles.Value {
			if role.Min.State != c.Known || role.Max.State != c.Known || role.Kinds.State != c.Known {
				return true
			}
		}
		for _, attribute := range feature.Attributes.Value {
			if attribute.Required.State != c.Known || attribute.Values.State != c.Known {
				return true
			}
		}
	}
	return false
}

type parsed struct {
	file  c.File
	token *token.File
	nodes map[[2]int]ast.Node
}

type expression struct {
	node  ast.Expr
	file  *parsed
	env   *bindings
	items []*expression // non-nil only for a captured variadic parameter
	at    c.Span
}

type bindings struct {
	values map[string]*expression
}

type visit struct {
	node ast.Node
	op   string
	env  *bindings
}

type evaluator struct {
	ctx        context.Context
	request    c.DefinitionRequest
	limits     c.Limits
	files      map[c.FileRef]*parsed
	active     map[visit]bool
	expansions int
	err        error
	uncertain  bool
	result     c.Inventory
}

func (e *evaluator) guard(x *expression, operation string, depth int) func() {
	if e.err != nil {
		return nil
	}
	if err := e.ctx.Err(); err != nil {
		e.err = err
		return nil
	}
	e.expansions++
	if depth > e.limits.MaxDepth || e.expansions > e.limits.MaxExpansions {
		e.err = fmt.Errorf("grammar extraction traversal limit exceeded")
		return nil
	}
	if x == nil {
		e.err = fmt.Errorf("missing source expression")
		return nil
	}
	key := visit{x.node, operation, x.env}
	if x.node != nil && e.active[key] {
		e.err = fmt.Errorf("cycle in grammar declaration references")
		return nil
	}
	if x.node != nil {
		e.active[key] = true
	}
	return func() { delete(e.active, key) }
}

func (e *evaluator) load(ref c.FileRef) *parsed {
	if e.err != nil {
		return nil
	}
	if file := e.files[ref]; file != nil {
		return file
	}
	file, err := e.request.Source.File(ref)
	if err != nil {
		e.err = err
		return nil
	}
	set := token.NewFileSet()
	node, err := parser.ParseFile(set, file.Path, file.Text, parser.SkipObjectResolution)
	if err != nil {
		e.err = fmt.Errorf("parse captured source: %w", err)
		return nil
	}
	p := &parsed{file: file, token: set.File(node.Pos()), nodes: map[[2]int]ast.Node{}}
	ast.Inspect(node, func(n ast.Node) bool {
		if n != nil {
			p.nodes[[2]int{p.token.Offset(n.Pos()), p.token.Offset(n.End())}] = n
		}
		return true
	})
	e.files[ref] = p
	return p
}

func (e *evaluator) expression(span c.Span) *expression {
	if e.err != nil {
		return nil
	}
	if err := c.ValidateSpan(e.request.Source, span); err != nil {
		e.err = err
		return nil
	}
	p := e.load(span.File)
	if p == nil {
		return nil
	}
	node, ok := p.nodes[[2]int{span.Start, span.End}].(ast.Expr)
	if !ok {
		e.err = fmt.Errorf("source span is not an expression")
		return nil
	}
	return &expression{node: node, file: p, at: span}
}

func child(x *expression, node ast.Expr) *expression {
	return &expression{node: node, file: x.file, env: x.env, at: c.Span{File: x.file.file.Ref, Start: x.file.token.Offset(node.Pos()), End: x.file.token.Offset(node.End())}}
}

func (e *evaluator) facts(x *expression) c.Facts {
	if e.err != nil || x == nil || x.node == nil {
		return c.Facts{}
	}
	facts, err := e.request.Source.Facts(x.at)
	if err != nil {
		e.err = err
	}
	return facts
}

// follow resolves only a source alias/initializer, never an evaluated value.
func (e *evaluator) follow(x *expression) *expression {
	if p, ok := x.node.(*ast.ParenExpr); ok {
		return child(x, p.X)
	}
	if id, ok := x.node.(*ast.Ident); ok {
		if x.env != nil {
			if bound := x.env.values[id.Name]; bound != nil {
				return bound
			}
		}
	}
	switch x.node.(type) {
	case *ast.Ident, *ast.SelectorExpr:
		facts := e.facts(x)
		if facts.Symbol != nil {
			decl, err := e.request.Source.Declaration(*facts.Symbol)
			if err == nil && decl.Symbol == *facts.Symbol && decl.Value.File != "" {
				return e.expression(decl.Value)
			}
		}
	}
	return nil
}

func (e *evaluator) callee(x *expression) c.Symbol {
	facts := e.facts(x)
	if facts.Callee != nil {
		return *facts.Callee
	}
	if call, ok := x.node.(*ast.CallExpr); ok {
		facts = e.facts(child(x, call.Fun))
		if facts.Symbol != nil {
			return *facts.Symbol
		}
	}
	return c.Symbol{}
}

func (e *evaluator) appendCall(x *expression) bool {
	return e.callee(x) == (c.Symbol{Package: "builtin", Name: "append"})
}

// inline interprets only one-return structural helper syntax, substituting its
// formal expressions. Blocks, mutation and arbitrary execution remain opaque.
func (e *evaluator) inline(x *expression) *expression {
	call, ok := x.node.(*ast.CallExpr)
	if !ok {
		return nil
	}
	symbol := e.callee(x)
	if symbol.Package == "" || symbol.Package == owner || symbol.Package == "builtin" {
		return nil
	}
	decl, err := e.request.Source.Declaration(symbol)
	if err != nil || decl.Symbol != symbol {
		return nil
	}
	p := e.load(decl.Span.File)
	if p == nil {
		return nil
	}
	fn, ok := p.nodes[[2]int{decl.Span.Start, decl.Span.End}].(*ast.FuncDecl)
	if !ok || fn.Recv != nil || fn.Body == nil || len(fn.Body.List) != 1 {
		return nil
	}
	ret, ok := fn.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return nil
	}
	parameters := []string{}
	variadic := false
	if fn.Type.Params != nil {
		for _, parameter := range fn.Type.Params.List {
			if _, ok := parameter.Type.(*ast.Ellipsis); ok {
				variadic = true
			}
			for _, name := range parameter.Names {
				parameters = append(parameters, name.Name)
			}
		}
	}
	if !variadic && len(parameters) != len(call.Args) || variadic && len(call.Args) < len(parameters)-1 {
		return nil
	}
	var env *bindings
	if len(parameters) > 0 {
		env = &bindings{values: map[string]*expression{}}
	}
	for i, name := range parameters {
		if variadic && i == len(parameters)-1 {
			if call.Ellipsis.IsValid() {
				env.values[name] = child(x, call.Args[i])
			} else {
				values := make([]*expression, 0, len(call.Args)-i)
				for _, arg := range call.Args[i:] {
					values = append(values, child(x, arg))
				}
				env.values[name] = &expression{items: values, at: x.at}
			}
		} else {
			env.values[name] = child(x, call.Args[i])
		}
	}
	result := ret.Results[0]
	return &expression{node: result, file: p, env: env, at: c.Span{File: p.file.Ref, Start: p.token.Offset(result.Pos()), End: p.token.Offset(result.End())}}
}

func known[T any](at c.Span, value T) c.Field[T] {
	return c.Field[T]{State: c.Known, Value: value, At: at}
}

func unknown[T any](e *evaluator, at c.Span, reason string) c.Field[T] {
	e.uncertain = true
	return c.Field[T]{State: c.Unknown, At: at, Reason: reason}
}

func (e *evaluator) limitation(at c.Span, message string) {
	e.uncertain = true
	e.result.Limitations = append(e.result.Limitations, c.Limitation{At: at, Message: message})
}

func (e *evaluator) constant(x *expression, kind string, depth int) (string, bool) {
	done := e.guard(x, "constant", depth)
	if done == nil {
		return "", false
	}
	defer done()
	if value := e.facts(x).Constant; value != nil && value.Kind == kind {
		return value.Value, true
	}
	if next := e.follow(x); next != nil {
		return e.constant(next, kind, depth+1)
	}
	return "", false
}

func (e *evaluator) boolean(x *expression, depth int) c.Field[bool] {
	if value, ok := e.constant(x, "bool", depth); ok && (value == "true" || value == "false") {
		return known(x.at, value == "true")
	}
	return unknown[bool](e, x.at, "Boolean expression could not be decoded")
}

func (e *evaluator) integer(x *expression, depth int) c.Field[int] {
	if value, ok := e.constant(x, "int", depth); ok {
		if number, err := strconv.Atoi(value); err == nil {
			return known(x.at, number)
		}
	}
	return unknown[int](e, x.at, "Integer expression could not be decoded")
}

func (e *evaluator) isNil(x *expression) bool {
	symbol := e.facts(x).Symbol
	return symbol != nil && *symbol == (c.Symbol{Package: "builtin", Name: "nil"})
}

func (e *evaluator) rules(x *expression, depth int) {
	done := e.guard(x, "rules", depth)
	if done == nil {
		return
	}
	defer done()
	if next := e.follow(x); next != nil {
		e.rules(next, depth+1)
		return
	}
	if e.isNil(x) {
		return
	}
	literal, ok := x.node.(*ast.CompositeLit)
	if !ok {
		e.limitation(x.at, "Rules initializer is outside the supported declaration format")
		return
	}
	for _, element := range literal.Elts {
		pair, ok := element.(*ast.KeyValueExpr)
		if !ok {
			e.limitation(x.at, "Rules source entry has no decodable key/value pair")
			continue
		}
		key := child(x, pair.Key)
		name, ok := e.constant(key, "string", depth+1)
		if !ok || name == "" {
			e.limitation(key.at, "Unresolved feature key; reported entries do not establish the effective map")
			continue
		}
		e.result.Features = append(e.result.Features, e.feature(name, child(x, pair.Value), depth+1))
	}
}

func fields(x *expression, names ...string) (map[string]*expression, bool) {
	literal, ok := x.node.(*ast.CompositeLit)
	if !ok {
		return nil, false
	}
	result := map[string]*expression{}
	for i, element := range literal.Elts {
		if pair, ok := element.(*ast.KeyValueExpr); ok {
			name, ok := pair.Key.(*ast.Ident)
			if !ok {
				return nil, false
			}
			result[name.Name] = child(x, pair.Value)
		} else if i < len(names) {
			result[names[i]] = child(x, element)
		} else {
			return nil, false
		}
	}
	return result, true
}

func (e *evaluator) feature(key string, x *expression, depth int) c.Feature {
	f := c.Feature{Key: key, At: x.at}
	done := e.guard(x, "feature", depth)
	if done == nil {
		return f
	}
	defer done()
	if next := e.follow(x); next != nil {
		return e.feature(key, next, depth+1)
	}
	fs, ok := fields(x, "Attributes", "Roles", "Reference")
	if !ok {
		if next := e.inline(x); next != nil {
			return e.feature(key, next, depth+1)
		}
		f.Roles = unknown[[]c.Role](e, x.at, "Rule expression could not be decoded")
		f.Attributes = unknown[[]c.Attribute](e, x.at, "Rule expression could not be decoded")
		f.References = unknown[[]string](e, x.at, "Rule expression could not be decoded")
		return f
	}
	f.Roles = known[[]c.Role](x.at, nil)
	f.Attributes = known[[]c.Attribute](x.at, nil)
	f.References = known[[]string](x.at, nil)
	if value := fs["Roles"]; value != nil {
		f.Roles = e.roles(value, depth+1)
	}
	if value := fs["Attributes"]; value != nil {
		f.Attributes = e.attributes(value, depth+1)
	}
	if value := fs["Reference"]; value != nil {
		allowed := e.boolean(value, depth+1)
		if allowed.State != c.Known {
			f.References = unknown[[]string](e, value.at, allowed.Reason)
		} else if allowed.Value {
			f.References = known(value.at, []string{"target"})
		}
	}
	return f
}

func unique(values []string) []string {
	var result []string
	seen := map[string]bool{}
	for _, value := range values {
		if !seen[value] {
			result = append(result, value)
			seen[value] = true
		}
	}
	return result
}

func merge[T any](e *evaluator, at c.Span, values ...c.Field[[]T]) c.Field[[]T] {
	result := known[[]T](at, nil)
	for _, value := range values {
		result.Value = append(result.Value, value.Value...)
		if value.State != c.Known {
			e.uncertain = true
			result.State = c.Partial
			result.Reason = "Collection contains an undecoded portion"
		}
	}
	return result
}

func (e *evaluator) stringArgs(x *expression, args []ast.Expr, ellipsis bool, depth int) c.Field[[]string] {
	parts := []c.Field[[]string]{}
	for i, arg := range args {
		value := child(x, arg)
		if ellipsis && i == len(args)-1 {
			parts = append(parts, e.strings(value, depth+1))
		} else if text, ok := e.constant(value, "string", depth+1); ok {
			parts = append(parts, known(value.at, []string{text}))
		} else {
			parts = append(parts, unknown[[]string](e, value.at, "String expression could not be decoded"))
		}
	}
	result := merge(e, x.at, parts...)
	result.Value = unique(result.Value)
	return result
}

func (e *evaluator) strings(x *expression, depth int) c.Field[[]string] {
	done := e.guard(x, "strings", depth)
	if done == nil {
		return c.Field[[]string]{}
	}
	defer done()
	if x.items != nil {
		parts := []c.Field[[]string]{}
		for _, item := range x.items {
			if value, ok := e.constant(item, "string", depth+1); ok {
				parts = append(parts, known(item.at, []string{value}))
			} else {
				parts = append(parts, unknown[[]string](e, item.at, "Variadic string could not be decoded"))
			}
		}
		return merge(e, x.at, parts...)
	}
	if next := e.follow(x); next != nil {
		return e.strings(next, depth+1)
	}
	if e.isNil(x) {
		return known[[]string](x.at, nil)
	}
	if literal, ok := x.node.(*ast.CompositeLit); ok {
		return e.stringArgs(x, literal.Elts, false, depth+1)
	}
	if call, ok := x.node.(*ast.CallExpr); ok && e.appendCall(x) && len(call.Args) > 0 {
		result := merge(e, x.at, e.strings(child(x, call.Args[0]), depth+1), e.stringArgs(x, call.Args[1:], call.Ellipsis.IsValid(), depth+1))
		result.Value = unique(result.Value)
		return result
	}
	if next := e.inline(x); next != nil {
		return e.strings(next, depth+1)
	}
	return unknown[[]string](e, x.at, "Allowed-kind expression could not be decoded")
}

func (e *evaluator) roles(x *expression, depth int) c.Field[[]c.Role] {
	done := e.guard(x, "roles", depth)
	if done == nil {
		return c.Field[[]c.Role]{}
	}
	defer done()
	if next := e.follow(x); next != nil {
		return e.roles(next, depth+1)
	}
	if e.isNil(x) {
		return known[[]c.Role](x.at, nil)
	}
	values := x.items
	if literal, ok := x.node.(*ast.CompositeLit); ok {
		values = make([]*expression, 0, len(literal.Elts))
		for _, element := range literal.Elts {
			values = append(values, child(x, element))
		}
	}
	if values != nil {
		result := known[[]c.Role](x.at, nil)
		for _, value := range values {
			role, ok := e.role(value, depth+1)
			if ok {
				result.Value = append(result.Value, role)
			} else {
				result.State, result.Reason = c.Partial, "Role collection contains an undecoded entry"
				e.uncertain = true
			}
		}
		return result
	}
	if call, ok := x.node.(*ast.CallExpr); ok && e.appendCall(x) && len(call.Args) > 0 {
		parts := []c.Field[[]c.Role]{e.roles(child(x, call.Args[0]), depth+1)}
		for i, argument := range call.Args[1:] {
			value := child(x, argument)
			if call.Ellipsis.IsValid() && i == len(call.Args)-2 {
				parts = append(parts, e.roles(value, depth+1))
			} else if role, ok := e.role(value, depth+1); ok {
				parts = append(parts, known(value.at, []c.Role{role}))
			} else {
				parts = append(parts, unknown[[]c.Role](e, value.at, "Appended role could not be decoded"))
			}
		}
		return merge(e, x.at, parts...)
	}
	if next := e.inline(x); next != nil {
		return e.roles(next, depth+1)
	}
	return unknown[[]c.Role](e, x.at, "Role collection could not be decoded")
}

func (e *evaluator) role(x *expression, depth int) (c.Role, bool) {
	done := e.guard(x, "role", depth)
	if done == nil {
		return c.Role{}, false
	}
	defer done()
	if next := e.follow(x); next != nil {
		return e.role(next, depth+1)
	}
	role := c.Role{At: x.at, Min: known(x.at, 0), Max: known(x.at, c.UpperBound{}), Kinds: known[[]string](x.at, nil)}
	if fs, ok := fields(x, "Name", "Min", "Max", "Kinds"); ok {
		name := fs["Name"]
		if name == nil {
			return c.Role{}, false
		}
		role.Name, ok = e.constant(name, "string", depth+1)
		if !ok || role.Name == "" {
			return c.Role{}, false
		}
		if value := fs["Min"]; value != nil {
			role.Min = e.integer(value, depth+1)
		}
		if value := fs["Max"]; value != nil {
			maximum := e.integer(value, depth+1)
			if maximum.State != c.Known {
				role.Max = unknown[c.UpperBound](e, value.at, maximum.Reason)
			} else if maximum.Value < 0 {
				role.Max = known(value.at, c.UpperBound{Unbounded: true})
			} else {
				role.Max = known(value.at, c.UpperBound{Value: maximum.Value})
			}
		}
		if value := fs["Kinds"]; value != nil {
			role.Kinds = e.strings(value, depth+1)
		}
		return role, true
	}
	if next := e.inline(x); next != nil {
		return e.role(next, depth+1)
	}
	return c.Role{}, false
}

func (e *evaluator) attributes(x *expression, depth int) c.Field[[]c.Attribute] {
	done := e.guard(x, "attributes", depth)
	if done == nil {
		return c.Field[[]c.Attribute]{}
	}
	defer done()
	if next := e.follow(x); next != nil {
		return e.attributes(next, depth+1)
	}
	if e.isNil(x) {
		return known[[]c.Attribute](x.at, nil)
	}
	result := known[[]c.Attribute](x.at, nil)
	add := func(key, value *expression) {
		name, ok := e.constant(key, "string", depth+1)
		if !ok || name == "" {
			result.State, result.Reason = c.Partial, "Unresolved attribute key; reported source entries do not establish the effective map"
			e.uncertain = true
			return
		}
		// A map literal cannot repeat a constant key, so a known name stands
		// once. An unresolved key still leaves the collection partial because
		// it can collide with a reported source entry.
		result.Value = append(result.Value, e.attribute(name, value, depth+1))
	}
	if literal, ok := x.node.(*ast.CompositeLit); ok {
		for _, element := range literal.Elts {
			pair, ok := element.(*ast.KeyValueExpr)
			if !ok {
				result.State, result.Reason = c.Partial, "Attribute source entry has no key/value pair"
				e.uncertain = true
				continue
			}
			add(child(x, pair.Key), child(x, pair.Value))
		}
		return result
	}
	if next := e.inline(x); next != nil {
		return e.attributes(next, depth+1)
	}
	return unknown[[]c.Attribute](e, x.at, "Attribute collection could not be decoded")
}

func (e *evaluator) attribute(name string, x *expression, depth int) c.Attribute {
	attribute := c.Attribute{Name: name, At: x.at}
	done := e.guard(x, "attribute", depth)
	if done == nil {
		return attribute
	}
	defer done()
	if next := e.follow(x); next != nil {
		return e.attribute(name, next, depth+1)
	}
	if fs, ok := fields(x, "Required", "Predicate", "OneOf"); ok {
		attribute.Required = known(x.at, false)
		if value := fs["Required"]; value != nil {
			attribute.Required = e.boolean(value, depth+1)
		}
		attribute.Values = e.values(x, fs["Predicate"], fs["OneOf"], depth+1)
		return attribute
	}
	if next := e.inline(x); next != nil {
		return e.attribute(name, next, depth+1)
	}
	attribute.Required = unknown[bool](e, x.at, "Attribute expression could not be decoded")
	attribute.Values = unknown[c.ValueSet](e, x.at, "Attribute expression could not be decoded")
	return attribute
}

// values interprets the two constraints of an attribute entry. A set gives
// its texts, and an empty set restricts nothing, as in the table type. Of the
// predicates only the shared name Bool is interpreted, by its resolved symbol:
// it admits "true" and "false". Any other name is reported as present and its
// behavior is not interpreted. With a set and a predicate both hold.
func (e *evaluator) values(x *expression, predicate, oneOf *expression, depth int) c.Field[c.ValueSet] {
	set := known(x.at, c.ValueSet{Unrestricted: true})
	if oneOf != nil {
		texts := e.strings(oneOf, depth)
		if texts.State != c.Known || len(texts.Value) > 0 {
			set = c.Field[c.ValueSet]{State: texts.State, Value: c.ValueSet{Values: unique(texts.Value)}, At: oneOf.at, Reason: texts.Reason}
		}
	}
	if predicate == nil {
		return set
	}
	switch e.predicate(predicate, depth) {
	case noPredicate:
		return set
	case boolPredicate:
		if set.State != c.Known {
			return set
		}
		values := []string{}
		for _, value := range []string{"true", "false"} {
			if set.Value.Unrestricted || slices.Contains(set.Value.Values, value) {
				values = append(values, value)
			}
		}
		return known(predicate.at, c.ValueSet{Values: values})
	}
	e.uncertain = true
	return c.Field[c.ValueSet]{State: c.Opaque, At: predicate.at, Reason: "Predicate present; behavior not interpreted"}
}

const (
	otherPredicate = iota
	noPredicate
	boolPredicate
)

// predicate classifies the expression of an entry's predicate name: the zero
// value of the name type, the shared name Bool, or any other name.
func (e *evaluator) predicate(x *expression, depth int) int {
	done := e.guard(x, "predicate", depth)
	if done == nil {
		return otherPredicate
	}
	defer done()
	if next := e.follow(x); next != nil {
		return e.predicate(next, depth+1)
	}
	if _, ok := x.node.(*ast.CallExpr); ok && e.callee(x) == (c.Symbol{Package: owner, Name: "Bool"}) {
		return boolPredicate
	}
	if literal, ok := x.node.(*ast.CompositeLit); ok && len(literal.Elts) == 0 {
		if origin := e.facts(x).Type.Origin; origin != nil && *origin == (c.Symbol{Package: owner, Name: "Predicate"}) {
			return noPredicate
		}
	}
	return otherPredicate
}
