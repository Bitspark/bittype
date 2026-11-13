package target

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
)

// CheckDistinctParameters is the name rule that lies inside one declaration:
// no two parameters of one declaration share a name. It reads nothing but the
// names of each declaration's own parameters, so it knows no language and no
// scope across declarations. It refuses with ErrUnsupported. Whether a language
// asks for the rule is its target's choice: a target selects it in its program
// rule, alone or as part of CheckSingleScope.
func CheckDistinctParameters[T any](ctx context.Context, p g.Program[T]) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	for _, d := range p.Declarations() {
		parameters := map[string]bool{}
		for _, parameter := range d.Description.Parameters {
			if parameters[parameter.Name] {
				return g.ErrUnsupported
			}
			parameters[parameter.Name] = true
		}
	}
	return nil
}

// CheckSingleScope qualifies this target's current single-scope construction
// profile. Other target organizations can supply their own program rules. It
// states two rules of its own, that no two declarations share a name and that
// no parameter is named like a declaration, and it selects
// CheckDistinctParameters for the rule inside a declaration.
func CheckSingleScope[T any](ctx context.Context, p g.Program[T]) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	names := map[string]bool{}
	for _, d := range p.Declarations() {
		if names[d.Description.Name] {
			return g.ErrUnsupported
		}
		names[d.Description.Name] = true
	}
	for _, d := range p.Declarations() {
		for _, parameter := range d.Description.Parameters {
			if names[parameter.Name] {
				return g.ErrUnsupported
			}
		}
	}
	return CheckDistinctParameters(ctx, p)
}
