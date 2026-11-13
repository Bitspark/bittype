// Package readonlyprofile is a test fixture: a family of horos-to-typescript
// contributions under a generation profile of its own. It is written the way
// an author of a second profile would write one. It states its own profile
// identity, its own publisher and its own payload schemas, and it builds its
// contributions from the parts that the horos mapping side offers to every
// profile: the source selectors, the common declaration policy, the horos
// source reader and the TypeScript model.
//
// The family differs from the builtin structural family, so that a request
// shows which of the two it received:
//
//   - a horos list becomes ReadonlyArray<E> and not E[];
//   - every declaration becomes an alias, also one whose body is a record;
//   - records are not mapped at all.
//
// The horos mapping package and the generation registry know nothing of this
// package or of its profile.
package readonlyprofile

import (
	"context"
	"errors"
	"sync"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	recipe "github.com/Bitspark/bittype/transform/syntax/horos/mapping/declarations/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// Profile is this family's own promise: lists are read-only views. Its
// author states it and answers for it.
const Profile = "integration-fixture/readonly-collections/1"

// Publisher is the publisher of the family's contributions.
const Publisher = "example.test/readonly-collections"

// reference is the reference of the family's contribution with the local key
// key.
func reference(key string) g.ContributionReference {
	return g.ContributionReference{
		Identity:       g.ContributionIdentity{Publisher: Publisher, Key: key},
		Meaning:        g.Meaning{Identity: "integration-fixture/readonly-collections/" + key, Revision: "1"},
		Representation: "typescript/1",
	}
}

// ListReference is the reference of the list contribution. A test reads it in
// the inspectable plan.
var ListReference = reference("list")

// Contribution is the type of the family's contributions, the same type that
// the builtin horos-to-typescript mapping packages return.
type Contribution = g.Contribution[h.Language, s.Language, source.Reader]

// Composition is the family as a consumer's composition registers it: the
// contributions, and one codec for each payload schema that they state.
type Composition struct {
	Contributions []Contribution
	Codecs        []g.PayloadCodec
}

// empty is the payload of a contribution that records nothing, and
// emptySchema its schema.
type empty struct{}

var emptySchema = g.SchemaReference{Identity: "integration-fixture/readonly-collections/empty", Version: "1"}

// declarationSchema is the schema of what the declaration policy records.
var declarationSchema = g.SchemaReference{Identity: "integration-fixture/readonly-collections/declaration", Version: "1"}

// declarationPayload is what the declaration policy records: the native name
// and the names of the native parameters. Every declaration becomes an alias,
// so the form is no part of the decision.
type declarationPayload struct {
	Name       string
	Parameters []string
}

func emptyCodec() (g.PayloadCodec, error) {
	return f.Codec(emptySchema, "module fixtures.readonly_empty_payload;\n\ntype Payload = { };\n",
		func(empty) (core.Value, error) { return data.EncodeMap(nil) },
		func(v core.Value) (empty, error) {
			entries, err := data.ReadMap(v)
			if err == nil && len(entries) != 0 {
				err = errors.New("the empty payload has no field")
			}
			return empty{}, err
		})
}

func declarationCodec() (g.PayloadCodec, error) {
	return f.Codec(declarationSchema, "module fixtures.readonly_declaration_payload;\n\ntype Payload = { name: text; parameters: list<text>; };\n",
		func(p declarationPayload) (core.Value, error) {
			name, err := data.EncodeText(p.Name)
			if err != nil {
				return nil, err
			}
			parameters := make([]core.Value, len(p.Parameters))
			for i, parameter := range p.Parameters {
				if parameters[i], err = data.EncodeText(parameter); err != nil {
					return nil, err
				}
			}
			var entries [][2]core.Value
			for _, field := range []struct {
				name  string
				value core.Value
			}{{"name", name}, {"parameters", data.EncodeList(parameters)}} {
				key, err := data.EncodeText(field.name)
				if err != nil {
					return nil, err
				}
				entries = append(entries, [2]core.Value{key, field.value})
			}
			return data.EncodeMap(entries)
		},
		func(v core.Value) (declarationPayload, error) {
			var p declarationPayload
			entries, err := data.ReadMap(v)
			if err != nil {
				return p, err
			}
			fields := map[string]core.Value{}
			for _, entry := range entries {
				key, err := data.ReadText(entry[0])
				if err != nil {
					return p, err
				}
				fields[key] = entry[1]
			}
			if len(fields) != 2 || fields["name"] == nil || fields["parameters"] == nil {
				return p, errors.New("the declaration payload has other fields than its schema")
			}
			if p.Name, err = data.ReadText(fields["name"]); err != nil {
				return p, err
			}
			parameters, err := data.ReadList(fields["parameters"])
			if err != nil {
				return p, err
			}
			p.Parameters = make([]string, len(parameters))
			for i, parameter := range parameters {
				if p.Parameters[i], err = data.ReadText(parameter); err != nil {
					return p, err
				}
			}
			return p, nil
		})
}

// codecs are made once: lowering a schema source is not free, and a codec
// holds no state of a plan.
var codecs struct {
	once               sync.Once
	empty, declaration g.PayloadCodec
	err                error
}

type rule = g.Rule[h.Language, source.Reader]
type nameStep = g.NameStep[h.Language, source.Reader]
type decideStep = g.DecideStep[h.Language, source.Reader]
type projection = g.Projection[s.Language]
type projected = g.Projected[s.Language]

// target is the target that the family's contributions state.
var target = string(s.Language{}.ID())

// Family returns the family for the source language from. It takes no
// profile: the profile is part of what each contribution states.
func Family(from representation.Language) (Composition, error) {
	codecs.once.Do(func() {
		if codecs.empty, codecs.err = emptyCodec(); codecs.err == nil {
			codecs.declaration, codecs.err = declarationCodec()
		}
	})
	if codecs.err != nil {
		return Composition{}, codecs.err
	}
	text, err := source.Builtin(from, "scalar", "utf8-text")
	if err != nil {
		return Composition{}, err
	}
	boolean, err := source.Builtin(from, "scalar", "bool")
	if err != nil {
		return Composition{}, err
	}
	list, err := source.Builtin(from, "list-of", "")
	if err != nil {
		return Composition{}, err
	}
	nothing := func(context.Context, decideStep) (g.Decision, error) { return g.Decision{Payload: empty{}}, nil }
	primitive := func(key string, selector g.Selector, name string) Contribution {
		return Contribution{
			Reference: reference(key), Target: target, Profile: Profile, Domain: g.ForForm(selector), Schema: emptySchema,
			Rule: rule{Decide: nothing},
			Project: func(context.Context, projection) (projected, error) {
				return projected{Type: model.Primitive(name)}, nil
			},
		}
	}
	alias := model.Form(model.FormAlias)
	contributions := []Contribution{
		primitive("text", text, model.PrimitiveString),
		primitive("boolean", boolean, model.PrimitiveBool),
		{
			Reference: ListReference, Target: target, Profile: Profile, Domain: g.ForForm(list), Schema: emptySchema,
			Rule: rule{Decide: nothing, Name: func(_ context.Context, step nameStep) ([]g.Request[h.Language], error) {
				element, err := step.Reader.Element(step.Subject.Use)
				if err != nil {
					return nil, err
				}
				return []g.Request[h.Language]{g.Reference(element)}, nil
			}},
			Project: func(_ context.Context, p projection) (projected, error) {
				if p.Dependencies() != 1 {
					return projected{}, g.ErrScope
				}
				return projected{Type: model.External("", "ReadonlyArray", p.Dependency(0).Type)}, nil
			},
		},
		{
			// The common declaration policy states no profile. It is used here
			// with the alias form for every declaration, so every declaration
			// is transparent.
			Reference: reference("declarations"), Target: target, Profile: Profile, Domain: g.ForDeclarationPolicy(), Schema: declarationSchema,
			Rule: rule{
				Name: func(_ context.Context, step nameStep) ([]g.Request[h.Language], error) {
					return recipe.Requests(step, model.EvidenceRecordForm), nil
				},
				Decide: func(_ context.Context, step decideStep) (g.Decision, error) {
					d, err := recipe.Decide(step, model.EvidenceRecordForm, Publisher)
					if err != nil {
						return g.Decision{}, err
					}
					return g.Decision{Payload: declarationPayload{Name: d.Name, Parameters: d.Parameters}, Transparent: true}, nil
				},
			},
			Declare: func(_ context.Context, p projection) (g.DeclarationDescription[s.Language], error) {
				payload, ok := p.Payload().(declarationPayload)
				if !ok {
					return g.DeclarationDescription[s.Language]{}, g.ErrInvalid
				}
				return recipe.Describe(p, recipe.Decision{Name: payload.Name, Parameters: payload.Parameters}, alias, alias)
			},
		},
	}
	return Composition{Contributions: contributions, Codecs: []g.PayloadCodec{codecs.empty, codecs.declaration}}, nil
}
