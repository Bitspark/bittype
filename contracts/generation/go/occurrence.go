package generation

import (
	"fmt"
	"slices"

	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// walk is the kernel's index over one captured graph for one planning run.
// It records only structure: declaration sites, binder ownership and the syntax
// sites of terms. Entity kinds, names, operators, payloads, locator meanings
// and facts are never read.
type walk[F any] struct {
	graph     semantics.Semantics[F]
	sites     map[ast.NodeRef][]semantics.EntityID
	binders   map[semantics.EntityID]semantics.EntityID // parameter -> declaring entity
	termSites map[semantics.TermID]ast.NodeRef
	roots     map[semantics.EntityID]*occurrence[F]
}

func newWalk[F any](graph semantics.Semantics[F]) *walk[F] {
	data := graph.Data()
	w := &walk[F]{graph: graph, sites: map[ast.NodeRef][]semantics.EntityID{}, binders: map[semantics.EntityID]semantics.EntityID{}, termSites: map[semantics.TermID]ast.NodeRef{}, roots: map[semantics.EntityID]*occurrence[F]{}}
	for i, e := range data.Entities {
		id := semantics.EntityID(i + 1)
		for _, site := range e.Declarations {
			w.sites[site] = append(w.sites[site], id)
		}
		for _, p := range e.Parameters {
			w.binders[p.Entity] = id
		}
	}
	for _, e := range data.Expressions {
		if _, ok := w.termSites[e.Term]; !ok {
			w.termSites[e.Term] = e.At.Site
		}
	}
	return w
}

// Declaration is a kernel handle for one declaration entity in one run.
type Declaration[F any] struct {
	walk   *walk[F]
	entity semantics.EntityID
}

// Binder is a kernel handle for one real binder (a declaration parameter).
type Binder[F any] struct {
	walk   *walk[F]
	entity semantics.EntityID
}

// Use is a kernel handle for one occurrence: a graph term reached under one
// declaration along one structural path in one run. Two paths to a shared term
// are different occurrences. Only the kernel mints occurrences.
type Use[F any] struct{ o *occurrence[F] }

type occurrence[F any] struct {
	walk        *walk[F]
	term        semantics.TermID
	declaration semantics.EntityID
	parent      *occurrence[F]
	children    map[int]*occurrence[F]
	view        *UseView[F]
	judged      bool
	judgment    error
}

func (d Declaration[F]) Valid() bool { return d.walk != nil && d.entity != 0 }
func (b Binder[F]) Valid() bool      { return b.walk != nil && b.entity != 0 }
func (u Use[F]) Valid() bool         { return u.o != nil }

// In resolves a handle against the graph its holder reads. It succeeds only if
// the handle was issued for exactly that graph, so a source can reject handles
// from a run over different evidence.
func (d Declaration[F]) In(graph semantics.Semantics[F]) (semantics.EntityID, bool) {
	if !d.Valid() || d.walk.graph != graph {
		return 0, false
	}
	return d.entity, true
}
func (b Binder[F]) In(graph semantics.Semantics[F]) (semantics.EntityID, bool) {
	if !b.Valid() || b.walk.graph != graph {
		return 0, false
	}
	return b.entity, true
}
func (u Use[F]) In(graph semantics.Semantics[F]) (semantics.TermID, bool) {
	if !u.Valid() || u.o.walk.graph != graph {
		return 0, false
	}
	return u.o.term, true
}

// Declaration is the declaration whose assumptions apply to this occurrence.
func (u Use[F]) Declaration() Declaration[F] {
	if !u.Valid() {
		return Declaration[F]{}
	}
	return Declaration[F]{u.o.walk, u.o.declaration}
}

type UseKind uint8

const (
	// BoundUse refers to a binder of the declaration being planned.
	BoundUse UseKind = iota + 1
	// DeclarationUse refers to or directly applies a declaration.
	DeclarationUse
	// ExpressionUse is a source-owned provider expression.
	ExpressionUse
)

// Actual pairs a callee's real binder with the occurrence that fills it.
type Actual[F any] struct {
	Binder Binder[F]
	Use    Use[F]
}

// Operand is a provider-expression child. The locator is the provider's own
// label for it, passed through unread; only the source interprets it.
type Operand[F any] struct {
	Locator []uint32
	Use     Use[F]
}

// UseView is the kernel's structural view of one occurrence.
type UseView[F any] struct {
	Use      Use[F]
	Kind     UseKind
	Bound    Binder[F]      // BoundUse
	Callee   Declaration[F] // DeclarationUse
	Actuals  []Actual[F]    // DeclarationUse, in the callee's parameter order
	Operands []Operand[F]   // ExpressionUse
	Body     bool           // the occurrence is its declaration's whole body
}

func (v UseView[F]) Clone() UseView[F] {
	v.Actuals = slices.Clone(v.Actuals)
	v.Operands = cloneOperands(v.Operands)
	return v
}

// Actual returns the occurrence that fills a callee binder.
func (v UseView[F]) Actual(b Binder[F]) (Use[F], bool) {
	for _, a := range v.Actuals {
		if a.Binder == b {
			return a.Use, true
		}
	}
	return Use[F]{}, false
}

// Operand returns the operand with exactly this provider locator.
func (v UseView[F]) Operand(locator ...uint32) (Use[F], bool) {
	for _, o := range v.Operands {
		if slices.Equal(o.Locator, locator) {
			return o.Use, true
		}
	}
	return Use[F]{}, false
}

// DeclarationView is the kernel's view of a declaration being planned.
type DeclarationView[F any] struct {
	Declaration Declaration[F]
	Formals     []Binder[F]
	Body        Use[F]
}

func (v DeclarationView[F]) Clone() DeclarationView[F] {
	v.Formals = slices.Clone(v.Formals)
	return v
}

// declaration returns the kernel's view of a declaration entity.
func (w *walk[F]) declaration(d Declaration[F]) (DeclarationView[F], error) {
	e, ok := w.graph.Entity(d.entity)
	if !ok {
		return DeclarationView[F]{}, ErrScope
	}
	if e.Body == 0 {
		return DeclarationView[F]{}, w.unsupported(0, d.entity, "declaration has no body")
	}
	body := w.roots[d.entity]
	if body == nil {
		body = &occurrence[F]{walk: w, term: e.Body, declaration: d.entity}
		w.roots[d.entity] = body
	}
	view := DeclarationView[F]{Declaration: d, Body: Use[F]{body}}
	for _, p := range e.Parameters {
		view.Formals = append(view.Formals, Binder[F]{w, p.Entity})
	}
	return view, nil
}

func (w *walk[F]) child(parent *occurrence[F], position int, term semantics.TermID) Use[F] {
	if parent.children == nil {
		parent.children = map[int]*occurrence[F]{}
	}
	c := parent.children[position]
	if c == nil {
		c = &occurrence[F]{walk: w, term: term, declaration: parent.declaration, parent: parent}
		parent.children[position] = c
	}
	return Use[F]{c}
}

// structure builds an occurrence's structural view: the term variant, its
// referenced entity and the children the kernel follows. A form with no direct
// structural reading is unsupported here; that is not a source judgment.
func (w *walk[F]) structure(o *occurrence[F]) (UseView[F], error) {
	if o.view != nil {
		return *o.view, nil
	}
	t, ok := w.graph.Term(o.term)
	if !ok {
		return UseView[F]{}, ErrScope
	}
	view := UseView[F]{Use: Use[F]{o}, Body: o.parent == nil}
	switch t.Kind {
	case semantics.Reference:
		if owner, binder := w.binders[t.Target]; binder {
			if owner != o.declaration {
				return UseView[F]{}, &SiteError{Err: ErrScope, Site: w.termSites[o.term], Message: "reference to a binder of another declaration"}
			}
			view.Kind, view.Bound = BoundUse, Binder[F]{w, t.Target}
			break
		}
		callee, ok := w.graph.Entity(t.Target)
		if !ok {
			return UseView[F]{}, ErrScope
		}
		if len(callee.Parameters) != 0 {
			return UseView[F]{}, w.unsupported(o.term, o.declaration, "generic declaration used without arguments")
		}
		view.Kind, view.Callee = DeclarationUse, Declaration[F]{w, t.Target}
	case semantics.Application:
		bindings, ok := w.graph.DirectParameterBindings(o.term)
		if !ok {
			return UseView[F]{}, w.unsupported(o.term, o.declaration, "application without a direct declaration callee and complete arguments")
		}
		callee, _ := w.graph.Term(t.Callee)
		view.Kind, view.Callee = DeclarationUse, Declaration[F]{w, callee.Target}
		for i, b := range bindings {
			view.Actuals = append(view.Actuals, Actual[F]{Binder[F]{w, b.Parameter}, w.child(o, i, b.Argument)})
		}
	case semantics.ProviderExpression:
		view.Kind = ExpressionUse
		for i, operand := range t.Operands {
			view.Operands = append(view.Operands, Operand[F]{Locator: slices.Clone(operand.Locator), Use: w.child(o, i, operand.Term)})
		}
	default:
		return UseView[F]{}, ErrInvalid
	}
	o.view = &view
	return view, nil
}

// unsupported locates a kernel limitation at the term's syntax site, else at
// the declaring site of the entity (for an occurrence, its declaration).
func (w *walk[F]) unsupported(term semantics.TermID, entity semantics.EntityID, message string) error {
	site, ok := w.termSites[term]
	if !ok && entity != 0 {
		if e, found := w.graph.Entity(entity); found && len(e.Declarations) > 0 {
			site, ok = e.Declarations[0], true
		}
	}
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnsupported, message)
	}
	return &SiteError{Err: ErrUnsupported, Site: site, Message: message}
}

// owns reports whether an occurrence handle was issued by this walk.
func (w *walk[F]) owns(u Use[F]) bool { return u.o != nil && u.o.walk == w }
