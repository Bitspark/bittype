package generationfixture

// The fixture's own contributions for its standard program: a declaration
// policy, two forms and one exact generic. None is a contribution of
// bittype's structural profile and none reads horos. A Mode makes one of them
// break one rule of the contract, for negative qualification.

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/tools/payload-schema/artifact/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	"github.com/Bitspark/bittype/transform/syntax/horos/payload/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// Publisher is the publisher of the fixture's contributions.
const Publisher = "example.test/generation-fixture"

// GoTarget is the target that the fixture's contributions state.
var GoTarget = string(s.Language{}.ID())

// Reference is the reference of a fixture contribution with the local key key
// and the representation revision revision.
func Reference(key, revision string) g.ContributionReference {
	return g.ContributionReference{
		Identity:       g.ContributionIdentity{Publisher: Publisher, Key: key},
		Meaning:        g.Meaning{Identity: "generation-fixture/" + key, Revision: "1"},
		Representation: revision,
	}
}

// Codec makes the codec of a fixture payload schema from its schema source,
// through the public route: the payload-schema library lowers the source to
// its canonical module and pin, and the shared part of the codec role loads
// it. The payload type of the module is named Payload.
func Codec[P any](reference g.SchemaReference, schemaSource string, encode func(P) (core.Value, error), decode func(core.Value) (P, error)) (g.PayloadCodec, error) {
	module, pin, err := artifact.Module(context.Background(), schemaSource)
	if err != nil {
		return nil, err
	}
	schema, err := payload.Load(reference, module, pin, "Payload")
	if err != nil {
		return nil, err
	}
	return payload.New(schema, encode, decode), nil
}

// Empty is the payload of a fixture contribution that records nothing.
type Empty struct{}

// EmptySchema is the schema reference of Empty.
var EmptySchema = g.SchemaReference{Identity: "generation-fixture/empty", Version: "1"}

// EmptyCodec is the codec of EmptySchema.
func EmptyCodec() (g.PayloadCodec, error) {
	return Codec(EmptySchema, "module fixtures.empty_payload;\n\ntype Payload = { };\n",
		func(Empty) (core.Value, error) { return data.EncodeMap(nil) },
		func(v core.Value) (Empty, error) {
			entries, err := data.ReadMap(v)
			if err == nil && len(entries) != 0 {
				err = errors.New("the empty payload has no field")
			}
			return Empty{}, err
		})
}

// policyPayload is what the fixture's declaration policy records.
type policyPayload struct {
	Name       string
	Nominal    bool
	Erased     bool
	Parameters []string
}

var policySchema = g.SchemaReference{Identity: "generation-fixture/policy", Version: "1"}

func policyCodec() (g.PayloadCodec, error) {
	names := []string{"name", "nominal", "erased", "parameters"}
	return Codec(policySchema, "module fixtures.policy_payload;\n\ntype Payload = { name: text; nominal: bool; erased: bool; parameters: list<text>; };\n",
		func(p policyPayload) (core.Value, error) {
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
			values := []core.Value{name, data.EncodeBool(p.Nominal), data.EncodeBool(p.Erased), data.EncodeList(parameters)}
			entries := make([][2]core.Value, len(names))
			for i, field := range names {
				key, err := data.EncodeText(field)
				if err != nil {
					return nil, err
				}
				entries[i] = [2]core.Value{key, values[i]}
			}
			return data.EncodeMap(entries)
		},
		func(v core.Value) (policyPayload, error) {
			var p policyPayload
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
			if len(fields) != len(names) {
				return p, errors.New("the policy payload has other fields than its schema")
			}
			if p.Name, err = data.ReadText(fields["name"]); err != nil {
				return p, err
			}
			if p.Nominal, err = data.ReadBool(fields["nominal"]); err != nil {
				return p, err
			}
			if p.Erased, err = data.ReadBool(fields["erased"]); err != nil {
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

// standardCodecs are made once: lowering a schema source is not free, and a
// codec holds no state of a plan.
var standardCodecs struct {
	once          sync.Once
	empty, policy g.PayloadCodec
	err           error
}

type rule = g.Rule[SourceLanguage, Reader]
type nameStep = g.NameStep[SourceLanguage, Reader]
type decideStep = g.DecideStep[SourceLanguage, Reader]
type request = g.Request[SourceLanguage]
type projection = g.Projection[s.Language]
type projected = g.Projected[s.Language]

// cursorHelper is the local name of the helper that the page contribution
// owns, and markerHelper of the one that the person contribution owns.
const (
	cursorHelper = "cursor"
	markerHelper = "marker"
)

// erased reports whether the policy describes a declaration as string and
// represents neither its parameter nor its body. Ignore omits its source
// parameter natively; uses must still admit it. Shadow's body is never
// represented; the kernel still admits it.
func erased(name string) bool { return name == "Ignore" || name == "Shadow" }

// standard is the composition of the standard program under a mode. kept is
// where a mode keeps a use between two plans.
func standard(from representation.Language, mode Mode, kept *g.Use[SourceLanguage]) (Composition, error) {
	selector := func(key string) (g.Selector, error) { return g.NewSelector(from, key) }
	standardCodecs.once.Do(func() {
		if standardCodecs.empty, standardCodecs.err = EmptyCodec(); standardCodecs.err == nil {
			standardCodecs.policy, standardCodecs.err = policyCodec()
		}
	})
	if standardCodecs.err != nil {
		return Composition{}, standardCodecs.err
	}
	person, err := selector("expression/person-body")
	if err != nil {
		return Composition{}, err
	}
	page, _ := selector("expression/page-body")
	cell, _ := selector("declaration/Cell")
	content := HelperAccess
	if mode == LargeHelperValue {
		content = strings.Repeat("a", LargeValueBytes)
	}
	access, err := helperAccess.Issue(content)
	if err != nil {
		return Composition{}, err
	}
	if mode == UnissuedHelperValue {
		access = g.IssuedValue{}
	}
	empty := func(context.Context, decideStep) (g.Decision, error) { return g.Decision{Payload: Empty{}}, nil }
	contributions := []Contribution{
		{
			Reference: Reference("person", "golang/1"), Target: GoTarget, Profile: Profile, Domain: g.ForForm(person), Schema: EmptySchema,
			Rule: rule{Decide: empty},
			// Under TwoHelpersOfOneName the person contribution owns a helper of
			// its own. It is spelled Cursor, like the helper that the page
			// contribution owns, and its body is another primitive. The two are
			// two declarations: the key of an owned declaration is its owner's
			// reference and local name, and the spelling is not the key.
			Owned: []g.OwnedDescription[s.Language]{{Name: markerHelper, Declaration: g.TargetDeclaration[s.Language]{Name: "Cursor", Form: model.Form("alias"), Body: model.Primitive("bool"), Values: []g.IssuedValue{access}}}},
			Project: func(_ context.Context, p projection) (projected, error) {
				if mode != TwoHelpersOfOneName {
					return projected{Type: model.Record(model.Field("Name", false, model.Primitive("string"))), Members: []g.Member{{Name: "Name"}}}, nil
				}
				marker, err := p.Owned(markerHelper)
				if err != nil {
					return projected{}, err
				}
				return projected{Type: model.Record(model.Field("Name", false, model.Primitive("string")), model.Field("Marker", false, marker)), Members: []g.Member{{Name: "Name"}, {Name: "Marker"}}}, nil
			},
		},
		{
			Reference: Reference("page", "golang/1"), Target: GoTarget, Profile: Profile, Domain: g.ForForm(page), Schema: EmptySchema,
			Rule: rule{Decide: empty, Name: func(_ context.Context, step nameStep) ([]request, error) {
				view := step.Subject.Use
				element, ok := view.Operand(0)
				if !ok {
					return nil, g.ErrScope
				}
				switch mode {
				case StrayCapture:
					element = view.Use
				case ForeignHandle:
					if kept.Valid() {
						element = *kept
					} else {
						*kept = element
					}
				case RecipeStray:
					*kept = element
				}
				return []request{g.Reference(element)}, nil
			}},
			Owned: []g.OwnedDescription[s.Language]{{Name: cursorHelper, Declaration: g.TargetDeclaration[s.Language]{Name: "Cursor", Form: model.Form("alias"), Body: model.Primitive("string"), Values: []g.IssuedValue{access}}}},
			Project: func(_ context.Context, p projection) (projected, error) {
				cursor, err := p.Owned(cursorHelper)
				if err != nil {
					return projected{}, err
				}
				return projected{Type: model.Record(model.Field("Items", false, model.Collection(p.Dependency(0).Type)), model.Field("Cursor", false, cursor)), Members: []g.Member{{Name: "Items"}, {Name: "Cursor"}}}, nil
			},
		},
		{
			// Cell's contribution represents every Cell use itself; Cell's body is never opened.
			Reference: Reference("cell", "golang/1"), Target: GoTarget, Profile: Profile, Domain: g.ForExact(cell), Schema: EmptySchema,
			Rule: rule{Decide: empty, Name: func(_ context.Context, step nameStep) ([]request, error) {
				reader, err := step.Reader()
				if err != nil {
					return nil, err
				}
				view := step.Subject.Use
				if mode == PatternStray {
					return []request{g.Reference(view.Use)}, nil
				}
				for _, a := range view.Actuals {
					if name, err := reader.ParameterName(a.Binder); err != nil {
						return nil, err
					} else if name == "V" {
						return []request{g.Reference(a.Use)}, nil
					}
				}
				return nil, g.ErrScope
			}},
			Project: func(_ context.Context, p projection) (projected, error) {
				return projected{Type: model.Collection(p.Dependency(0).Type)}, nil
			},
		},
		{
			Reference: Reference("declarations", "golang/1"), Target: GoTarget, Profile: Profile, Domain: g.ForDeclarationPolicy(), Schema: policySchema,
			Rule: rule{
				Name: func(_ context.Context, step nameStep) ([]request, error) {
					reader, err := step.Reader()
					if err != nil {
						return nil, err
					}
					name, err := reader.Name(step.Subject.Declaration.Declaration)
					if err != nil || erased(name) {
						return nil, err
					}
					body := step.Subject.Declaration.Body
					if mode == RecipeStray && name == "Outer" && kept.Valid() {
						body = *kept
					}
					return []request{g.Reference(body)}, nil
				},
				Decide: func(_ context.Context, step decideStep) (g.Decision, error) {
					reader, err := step.Reader()
					if err != nil {
						return g.Decision{}, err
					}
					name, err := reader.Name(step.Subject.Declaration.Declaration)
					if err != nil {
						return g.Decision{}, err
					}
					p := policyPayload{Name: name, Erased: erased(name), Nominal: name == "Person" || name == "Page"}
					if !p.Erased {
						for i := range step.Subject.Declaration.Formals {
							p.Parameters = append(p.Parameters, fmt.Sprintf("T%d", i))
						}
					}
					return g.Decision{Payload: p, Transparent: !p.Nominal && !p.Erased}, nil
				},
			},
			Declare: func(_ context.Context, p projection) (g.DeclarationDescription[s.Language], error) {
				zero := g.DeclarationDescription[s.Language]{}
				view, ok := p.Payload().(policyPayload)
				if !ok {
					return zero, g.ErrInvalid
				}
				if view.Erased {
					// The declaration is string: no native parameter, and its one binder means string.
					return g.DeclarationDescription[s.Language]{Name: view.Name, Form: model.Form("alias"), Body: model.Primitive("string"), Generic: g.GenericConvention[s.Language]{Bindings: []g.Expression[s.Language]{model.Primitive("string")}}}, nil
				}
				form := model.Form("alias")
				if view.Nominal {
					form = model.Form("defined")
				}
				value, err := intent(mode, view.Name)
				if err != nil {
					return zero, err
				}
				d := g.DeclarationDescription[s.Language]{Name: view.Name, Form: form, Values: []g.IssuedValue{value}, Body: p.Dependency(0).Type}
				if (mode == SameName || mode == ParameterNamedSame) && (view.Name == "Person" || view.Name == "Page") {
					d.Name = "Same"
				}
				if mode == ManyEmptyValues {
					empty, err := namingIntent.Issue("")
					if err != nil {
						return zero, err
					}
					d.Values = slices.Repeat([]g.IssuedValue{empty}, LargeValueBytes)
				}
				for i, name := range view.Parameters {
					d.Generic.Parameters = append(d.Generic.Parameters, g.ParameterDescription{Name: name})
					d.Generic.Bindings = append(d.Generic.Bindings, g.Parameter[s.Language](i))
					d.Generic.Arguments = append(d.Generic.Arguments, g.Argument[s.Language](i))
				}
				if view.Name == "Page" {
					// An additional native parameter that no source parameter stands
					// for. The declaration contribution owns it: it adds the
					// parameter, the argument of every application and the member
					// that uses it.
					storage := "Storage"
					switch mode {
					case ParameterNamedSame:
						storage = "Same"
					case SameParameters:
						storage = "T0"
					}
					d.Generic.Parameters = append(d.Generic.Parameters, g.ParameterDescription{Name: storage})
					d.Generic.Arguments = append(d.Generic.Arguments, model.Collection(g.Argument[s.Language](0)))
					members := d.Body.Arguments()
					if len(members) != 2 {
						return zero, g.ErrScope
					}
					d.Body = model.Record(members[0], model.Field("Storage", false, g.Parameter[s.Language](1)), members[1])
				}
				return d, nil
			},
		},
	}
	return Composition{Contributions: contributions, Codecs: []g.PayloadCodec{standardCodecs.empty, standardCodecs.policy}}, nil
}

// intent is the value the declaration policy states under a mode.
func intent(mode Mode, name string) (g.IssuedValue, error) {
	switch mode {
	case UnissuedValue:
		return g.IssuedValue{}, nil
	case LargeValue:
		return namingIntent.Issue(strings.Repeat("n", LargeValueBytes))
	}
	return namingIntent.Issue(name)
}
