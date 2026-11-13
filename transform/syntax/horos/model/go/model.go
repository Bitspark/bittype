// Package model describes intended horos types before AST construction: the
// intent model of horos as the TARGET of a generation.
//
// It exists for the pair typescript-to-horos, the second of the two other
// pairs of bittype#196 (decision O43), which show that the contribution
// contract and its kernel assume no native target. It declares the forms that
// pair needs and no other: the scalars text and bool, a required field, a
// closed record, and a type definition.
//
// horos has two other vocabularies in this module: the syntax vocabulary
// (syntax/horos/go), whose node kinds a tree is made of, and the analysed
// vocabulary (semantics/horos/go), whose names an analysed graph is read by.
// This is a third, and whether a language should have one model that several
// operations work on is an open question of the owner (bittype#305). Nothing
// is decided here: the package is at the path that the role of a target's
// model has for every language.
//
// COPY (AC-11, finding AG-7 of the two other pairs): the form issuer with
// Form and Name, the refusal of a foreign form with ErrScope, and the shape
// of CheckExpression and CheckDeclaration (the context, the form read back, a
// switch over the forms that checks the number of data and of arguments) are
// the same in the model of every target language. The TypeScript and the C++
// model differ in two form names. Only the form names, what a form's data
// may be and what a declaration's body may be are a language's. The
// invariant part has no agnostic implementation, so this package is its
// fifth copy. That the pair is green does not show that the mechanism has an
// owner.
package model

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
)

var forms = g.NewFormIssuer[s.Language]()

// The names of the forms this model issues, each declared here once. The
// first three are expression forms. FormDefinition is the declaration form.
const (
	FormScalar     = "scalar"
	FormField      = "field"
	FormRecord     = "record"
	FormDefinition = "type-definition"
)

// The selectors a scalar form carries: the two that this model declares of
// the selectors of the horos scalar descriptor.
const (
	ScalarText = "utf8-text"
	ScalarBool = "bool"
)

func Form(name string) g.Form[s.Language]      { return forms.Form(name) }
func Name(f g.Form[s.Language]) (string, bool) { return forms.Read(f) }

// Scalar is a horos scalar descriptor with the given selector.
func Scalar(selector string) g.Expression[s.Language] {
	return g.Apply(Form(FormScalar), []string{selector})
}

// Field is a required field of a record.
func Field(name string, ty g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Apply(Form(FormField), []string{name}, ty)
}

// Record is a closed record of the given fields.
func Record(fields ...g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Apply(Form(FormRecord), nil, fields...)
}

// CheckExpression admits an expression of this model's forms.
//
// Well-formedness of a form: a scalar carries one selector and no argument; a
// field carries a name that is not empty and one type; a record carries no
// data, its members are fields, and no two have one name.
//
// Coverage of this model, and no rule of horos: a scalar's selector is text
// or bool (horos has int too), every field is required and every record is
// closed, because the model has no form that says otherwise.
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
	case FormScalar:
		if len(d) != 1 || len(a) != 0 || (d[0] != ScalarText && d[0] != ScalarBool) {
			return g.ErrInvalid
		}
	case FormField:
		if len(d) != 1 || d[0] == "" || len(a) != 1 {
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
			if !ok || n != FormField || len(v) != 1 || seen[v[0]] {
				return g.ErrInvalid
			}
			seen[v[0]] = true
		}
	default:
		return g.ErrUnsupported
	}
	return nil
}

// CheckDeclaration admits a type definition with a name. A horos type
// definition names any type expression, so this model asks nothing of the
// body's form, and it has one declaration form: there is no choice between
// an alias and a defined declaration to make, and so no evidence kind for a
// declaration contribution to ask about a body. Whether the name is a horos
// identifier is a spelling rule, which the target checks with the lexical
// helper of horos: a model imports no helper.
//
// Coverage of this model, and no rule of horos: a definition has no
// parameter. A horos type definition may have parameters, and this model has
// no form for a use of one.
func CheckDeclaration(ctx context.Context, d g.TargetDeclaration[s.Language]) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	n, ok := Name(d.Form)
	if !ok {
		return g.ErrScope
	}
	if d.Name == "" || n != FormDefinition {
		return g.ErrInvalid
	}
	if len(d.Parameters) != 0 {
		return g.ErrUnsupported
	}
	return nil
}
