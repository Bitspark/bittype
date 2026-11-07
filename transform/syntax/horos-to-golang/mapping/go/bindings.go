package generation

import (
	"context"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	target "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// Provider opens a selected mapping using these concrete source and target APIs.
type Provider struct {
	rule rules.Rule[s.Language, target.Builder[h.Language]]
}

// bittype:scope docs/generation-composition.md#construction-and-selection
func (p Provider) Open(ctx context.Context, r g.MappingRequest[source.Reader, target.Builder[h.Language]]) (g.Contribution[h.Language, s.Language], error) {
	return p.rule.Open(ctx, r)
}
func Text(profile string) Provider {
	return Provider{rules.BuiltinRule(profile, "scalar", "utf8-text", textMapping)}
}
func Bool(profile string) Provider {
	return Provider{rules.BuiltinRule(profile, "scalar", "bool", boolMapping)}
}
func List(profile string) Provider {
	return Provider{rules.BuiltinRule(profile, "list-of", "", listMapping)}
}
func Record(profile string) Provider {
	return Provider{rules.BuiltinRule(profile, "record", "closed", recordMapping)}
}
func Binder(profile string) Provider {
	return Provider{rules.BinderRule[s.Language, target.Builder[h.Language]](profile)}
}
func Declarations(profile string) Provider {
	return Provider{rules.DeclarationRule[s.Language, target.Builder[h.Language]](profile)}
}
func Optional(profile string) Provider { return Provider{rules.HelperRule(profile, optionalHelper)} }

// Providers supplies this library's representations. Consumers can register any subset.
func Providers(profile string) []g.MappingProvider[h.Language, s.Language, source.Reader, target.Builder[h.Language]] {
	return []g.MappingProvider[h.Language, s.Language, source.Reader, target.Builder[h.Language]]{Text(profile), Bool(profile), List(profile), Record(profile), Binder(profile), Declarations(profile), Optional(profile)}
}
