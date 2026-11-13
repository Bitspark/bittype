// Package model describes intended native types before AST construction.
package model

import (
	"context"
	"strconv"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
)

var forms = g.NewFormIssuer[s.Language]()

// The names of the forms this model issues, each declared here once. Form
// takes one of them and Name returns one. The first five are the expression
// forms that the constructors below build. FormAlias and FormDefined are the
// declaration forms.
const (
	FormPrimitive = "primitive"
	FormSlice     = "slice"
	FormField     = "field"
	FormRecord    = "record"
	FormExternal  = "external"
	FormAlias     = "alias"
	FormDefined   = "defined"
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
	return g.Apply(Form(FormPrimitive), []string{name})
}
func Collection(element g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Apply(Form(FormSlice), nil, element)
}
func Field(name string, optional bool, ty g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Apply(Form(FormField), []string{name, strconv.FormatBool(optional)}, ty)
}
func Record(fields ...g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Apply(Form(FormRecord), nil, fields...)
}
func External(module, name string, args ...g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Apply(Form(FormExternal), []string{module, name}, args...)
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
	case FormSlice:
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
		return g.ErrUnsupported
	default:
		return g.ErrUnsupported
	}
	return nil
}
func CheckDeclaration(ctx context.Context, d g.TargetDeclaration[s.Language]) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	n, ok := Name(d.Form)
	if !ok {
		return g.ErrScope
	}
	if d.Name == "" || (n != FormAlias && n != FormDefined) {
		return g.ErrInvalid
	}
	if n == FormDefined {
		body, ok := Name(d.Body.Form())
		if !ok || body != FormRecord {
			return g.ErrInvalid
		}
	}
	return nil
}

// EvidenceRecordForm is the evidence kind "the selected native type is this
// model's record form". It belongs to this vocabulary, because the fact
// concerns the native side of this target (row 7). The declaration
// contribution asks for it about a declaration's body (bittype#196, row 3).
var EvidenceRecordForm = g.EvidenceKind{Identity: "bittype/golang/record-form", Version: "1"}

// CheckMember is the member vocabulary of this model (row 13): a member of a
// record is one field form.
func CheckMember(ctx context.Context, e g.Expression[s.Language]) error {
	if name, ok := Name(e.Form()); !ok || name != FormField {
		return g.ErrInvalid
	}
	return CheckExpression(ctx, e)
}
