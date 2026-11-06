package report

import (
	"fmt"
	"go/token"
	"go/types"
	"slices"
)

// TypeString is a display spelling, not a Go identity: it omits the owners of
// private members in unnamed types. Keep compiler types at every decision point.
type typeSet []types.Type

func (s typeSet) contains(t types.Type) bool {
	return slices.ContainsFunc(s, func(existing types.Type) bool { return types.Identical(existing, t) })
}

func (s *typeSet) add(t types.Type) {
	if !s.contains(t) {
		*s = append(*s, t)
	}
}

func sameNamed(a, b *types.Named) bool {
	if a == nil || b == nil {
		return a == b
	}
	return types.Identical(a, b)
}

func sameLanguage(a, b language) bool {
	return a.id == b.id && sameNamed(a.marker, b.marker)
}

func sameEndpoint(a, b *language) bool {
	if a == nil || b == nil {
		return a == b
	}
	return sameLanguage(*a, *b)
}

func functionOrigin(fn *types.Func) *types.Func {
	if fn == nil {
		return nil
	}
	return fn.Origin()
}

// A delegate can be shared by several operations in one contract/language
// context. That does not make the operations or their receivers interchangeable.
func sameTargetContext(a, b *operation) bool {
	return a.kind == b.kind && sameNamed(a.contract, b.contract) &&
		sameLanguage(a.language, b.language) && sameEndpoint(a.from, b.from) && sameEndpoint(a.to, b.to)
}

func sameOperation(a, b *operation) bool {
	return sameTargetContext(a, b) && functionOrigin(a.public) == functionOrigin(b.public) &&
		types.Identical(a.receiver, b.receiver)
}

func sameCallLink(a, b callLink) bool {
	return functionOrigin(a.caller) == functionOrigin(b.caller) &&
		functionOrigin(a.callee) == functionOrigin(b.callee) && a.pos == b.pos &&
		types.Identical(a.receiver, b.receiver) && slices.EqualFunc(a.arguments, b.arguments, types.Identical)
}

// Anchors are presentation keys. Reject collisions rather than allowing a label
// to collapse distinct compiler identities. Go identity already deduplicated
// operations before this check; no hash or spelling decides equivalence here.
func validateOperationAnchors(m *matrix) error {
	seen := map[string]*operation{}
	for _, op := range operations(m) {
		if previous := seen[op.id]; previous != nil {
			return fmt.Errorf("ambiguous report anchor for %s and %s", operationTitle(previous), operationTitle(op))
		}
		seen[op.id] = op
	}
	return nil
}

// Matrix rows abstract one selected language marker, not arbitrary type
// identity. This transform preserves Go objects (including private member
// ownership, nominal origins, tags and constraints); types.Identical compares
// the resulting types. The fresh parameter is shared only within one discovery.
func rowParameter() *types.TypeParam {
	return types.NewTypeParam(types.NewTypeName(token.NoPos, nil, "L", nil), types.NewInterfaceType(nil, nil).Complete())
}

func rowType(t types.Type, marker *types.Named, parameter *types.TypeParam) (types.Type, error) {
	memo := map[types.Type]types.Type{}
	var visit func(types.Type) (types.Type, error)
	variables := func(tuple *types.Tuple) (*types.Tuple, error) {
		vars := make([]*types.Var, tuple.Len())
		for i := range vars {
			v := tuple.At(i)
			t, err := visit(v.Type())
			if err != nil {
				return nil, err
			}
			vars[i] = types.NewVar(v.Pos(), v.Pkg(), v.Name(), t)
		}
		return types.NewTuple(vars...), nil
	}
	visit = func(t types.Type) (types.Type, error) {
		t = types.Unalias(t)
		if types.Identical(t, marker) {
			return parameter, nil
		}
		if cached := memo[t]; cached != nil {
			return cached, nil
		}
		var result types.Type
		switch t := t.(type) {
		case *types.Basic, *types.TypeParam:
			result = t
		case *types.Named:
			if t.TypeArgs().Len() == 0 {
				return t, nil // Never rewrite a named type's underlying definition.
			}
			args := make([]types.Type, t.TypeArgs().Len())
			for i := range args {
				var err error
				args[i], err = visit(t.TypeArgs().At(i))
				if err != nil {
					return nil, err
				}
			}
			var err error
			// The parameter is only a row abstraction, never an admitted program
			// instantiation. The original contract has already been checked.
			result, err = types.Instantiate(nil, t.Origin(), args, false)
			if err != nil {
				return nil, err
			}
		case *types.Pointer:
			elem, err := visit(t.Elem())
			if err != nil {
				return nil, err
			}
			result = types.NewPointer(elem)
		case *types.Slice:
			elem, err := visit(t.Elem())
			if err != nil {
				return nil, err
			}
			result = types.NewSlice(elem)
		case *types.Array:
			elem, err := visit(t.Elem())
			if err != nil {
				return nil, err
			}
			result = types.NewArray(elem, t.Len())
		case *types.Map:
			key, err := visit(t.Key())
			if err != nil {
				return nil, err
			}
			elem, err := visit(t.Elem())
			if err != nil {
				return nil, err
			}
			result = types.NewMap(key, elem)
		case *types.Chan:
			elem, err := visit(t.Elem())
			if err != nil {
				return nil, err
			}
			result = types.NewChan(t.Dir(), elem)
		case *types.Struct:
			fields, tags := make([]*types.Var, t.NumFields()), make([]string, t.NumFields())
			for i := range fields {
				field := t.Field(i)
				typ, err := visit(field.Type())
				if err != nil {
					return nil, err
				}
				fields[i] = types.NewField(field.Pos(), field.Pkg(), field.Name(), typ, field.Embedded())
				tags[i] = t.Tag(i)
			}
			result = types.NewStruct(fields, tags)
		case *types.Tuple:
			var err error
			result, err = variables(t)
			if err != nil {
				return nil, err
			}
		case *types.Signature:
			if t.TypeParams().Len() != 0 || t.RecvTypeParams().Len() != 0 {
				return nil, fmt.Errorf("open generic signature in matrix row")
			}
			params, err := variables(t.Params())
			if err != nil {
				return nil, err
			}
			results, err := variables(t.Results())
			if err != nil {
				return nil, err
			}
			// Receiver identity is not part of Go function-type identity.
			result = types.NewSignatureType(t.Recv(), nil, nil, params, results, t.Variadic())
		case *types.Interface:
			methods, embeds := make([]*types.Func, t.NumExplicitMethods()), make([]types.Type, t.NumEmbeddeds())
			for i := range methods {
				method := t.ExplicitMethod(i)
				typ, err := visit(method.Type())
				if err != nil {
					return nil, err
				}
				methods[i] = types.NewFunc(method.Pos(), method.Pkg(), method.Name(), typ.(*types.Signature))
			}
			for i := range embeds {
				var err error
				embeds[i], err = visit(t.EmbeddedType(i))
				if err != nil {
					return nil, err
				}
			}
			iface := types.NewInterfaceType(methods, embeds)
			if t.IsImplicit() {
				iface.MarkImplicit()
			}
			result = iface.Complete()
		case *types.Union:
			terms := make([]*types.Term, t.Len())
			for i := range terms {
				term := t.Term(i)
				typ, err := visit(term.Type())
				if err != nil {
					return nil, err
				}
				terms[i] = types.NewTerm(term.Tilde(), typ)
			}
			result = types.NewUnion(terms)
		default:
			return nil, fmt.Errorf("unsupported Go type in matrix row: %T", t)
		}
		memo[t] = result
		return result, nil
	}
	return visit(t)
}
