package target

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
)

// CheckSingleScope qualifies this target's current single-scope construction
// profile. Other target organizations can supply their own program rules.
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
		parameters := map[string]bool{}
		for _, parameter := range d.Description.Parameters {
			if names[parameter.Name] || parameters[parameter.Name] {
				return g.ErrUnsupported
			}
			parameters[parameter.Name] = true
		}
	}
	return nil
}
