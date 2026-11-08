// Package model describes intended native types before AST construction.
package model

import (
	"context"
	"strconv"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
)

var forms = g.NewFormIssuer[s.Language]()

func Form(name string) g.Form[s.Language]      { return forms.Form(name) }
func Name(f g.Form[s.Language]) (string, bool) { return forms.Read(f) }
func Primitive(name string) g.Expression[s.Language] {
	return g.Native(Form("primitive"), []string{name})
}
func Collection(element g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form("sequence"), nil, element)
}
func Field(name string, optional bool, ty g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form("field"), []string{name, strconv.FormatBool(optional)}, ty)
}
func Record(fields ...g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form("record"), nil, fields...)
}
func External(module, name string, args ...g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form("external"), []string{module, name}, args...)
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
	case "primitive":
		if len(d) != 1 || len(a) != 0 || (d[0] != "string" && d[0] != "bool") {
			return g.ErrInvalid
		}
	case "sequence":
		if len(d) != 0 || len(a) != 1 {
			return g.ErrInvalid
		}
	case "field":
		if len(d) != 2 || d[0] == "" || (d[1] != "true" && d[1] != "false") || len(a) != 1 {
			return g.ErrInvalid
		}
	case "record":
		if len(d) != 0 {
			return g.ErrInvalid
		}
		seen := map[string]bool{}
		for _, f := range a {
			n, ok := Name(f.Form())
			v := f.Data()
			if !ok || n != "field" || len(v) != 2 || seen[v[0]] {
				return g.ErrInvalid
			}
			seen[v[0]] = true
		}
	case "external":
		if len(d) != 2 || d[1] == "" {
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
	if d.Name == "" || (n != "alias" && n != "struct") {
		return g.ErrInvalid
	}
	if n == "struct" {
		body, ok := Name(d.Body.Form())
		if !ok || body != "record" {
			return g.ErrInvalid
		}
	}
	return nil
}
