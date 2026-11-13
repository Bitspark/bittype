// Package model describes intended native types before AST construction.
package model

import (
	"context"
	"slices"
	"strconv"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
)

var forms = g.NewFormIssuer[s.Language]()

// The names of the forms this model issues, each declared here once. Form
// takes one of them and Name returns one. The first five are the expression
// forms that the constructors below build. FormAlias and FormStruct are the
// declaration forms.
const (
	FormPrimitive = "primitive"
	FormSequence  = "sequence"
	FormField     = "field"
	FormRecord    = "record"
	FormExternal  = "external"
	FormAlias     = "alias"
	FormStruct    = "struct"
)

// The names a primitive form carries, each declared here once. Primitive
// takes one of them.
const (
	PrimitiveString = "string"
	PrimitiveBool   = "bool"
)

func Form(name string) g.Form[s.Language]      { return forms.Form(name) }
func Name(f g.Form[s.Language]) (string, bool) { return forms.Read(f) }
func Primitive(name string) g.Expression[s.Language] {
	return g.Native(Form(FormPrimitive), []string{name})
}
func Collection(element g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormSequence), nil, element)
}
func Field(name string, optional bool, ty g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormField), []string{name, strconv.FormatBool(optional)}, ty)
}
func Record(fields ...g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormRecord), nil, fields...)
}

// Std is the name of the crate std of the standard library, as an external
// form states it. It is declared here once: a mapping passes it, and the Rust
// target maps it to the dependency of that crate.
const Std = "std"

// External is the public item at path below the root of the crate that is
// named crate, applied to args. The crate is separate from the path: the form
// carries the crate's name first and then the segments of the path. This
// model checks the shape of the form only. The Rust target admits a form
// whose crate is a dependency that a destination can name. That is Std today,
// and the target refuses every other crate.
func External(crate string, path []string, args ...g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormExternal), append([]string{crate}, path...), args...)
}
func CheckExpression(ctx context.Context, e g.Expression[s.Language]) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	name, ok := Name(e.Form())
	if !ok {
		return g.ErrScope
	}
	d, a := e.Data(), e.Arguments()
	switch name {
	case FormPrimitive:
		if len(d) != 1 || len(a) != 0 || (d[0] != PrimitiveString && d[0] != PrimitiveBool) {
			return g.ErrInvalid
		}
	case FormSequence:
		if len(d) != 0 || len(a) != 1 {
			return g.ErrInvalid
		}
	case FormField:
		if len(d) != 2 || d[0] == "" || (d[1] != "true" && d[1] != "false") || len(a) != 1 {
			return g.ErrInvalid
		}
	case FormRecord:
		if len(d) != 0 {
			return g.ErrInvalid
		}
		seen := map[string]bool{}
		for _, f := range a {
			n, ok := Name(f.Form())
			v := f.Data()
			if !ok || n != FormField || len(v) != 2 || seen[v[0]] {
				return g.ErrInvalid
			}
			seen[v[0]] = true
		}
	case FormExternal:
		if len(d) < 2 || slices.Contains(d, "") {
			return g.ErrInvalid
		}
	default:
		return g.ErrUnsupported
	}
	return nil
}
func CheckDeclaration(ctx context.Context, d g.NativeDeclaration[s.Language]) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	n, ok := Name(d.Form)
	if !ok {
		return g.ErrScope
	}
	if d.Name == "" || (n != FormAlias && n != FormStruct) {
		return g.ErrInvalid
	}
	if n == FormStruct {
		body, ok := Name(d.Body.Form())
		if !ok || body != FormRecord {
			return g.ErrInvalid
		}
	}
	return nil
}
