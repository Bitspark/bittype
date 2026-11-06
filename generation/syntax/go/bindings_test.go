package generation_test

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	c "github.com/Bitspark/bittype/generation/syntax/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"testing"
)

type source struct{}

func (source) ID() language.ID { return "independent" }
func (source) Name() string    { return "Independent" }
func TestExactBindingsAndGuardedOverrides(t *testing.T) {
	lang, _ := representation.NewLanguage[source](source{}, "v1")
	key, _ := c.NewSelector(lang, "owner/module/type")
	mapping := func(id string) c.Mapping[source, source] {
		return c.Mapping[source, source]{ID: id, PluginID: "plugin", Plan: func(g.Declaration[source]) (g.Plan[source, source], error) { return g.Plan[source, source]{}, nil }}
	}
	base := c.Binding[source, source]{Selector: key, Mapping: mapping("base-v1")}
	replacement := c.Binding[source, source]{Selector: key, Mapping: mapping("new-v1")}
	if _, err := c.NewBindings(lang, []c.Binding[source, source]{base}, []c.Override[source, source]{{Binding: replacement, ExpectedID: "base-v1"}}, nil); err != nil {
		t.Fatal(err)
	}
	other, _ := c.NewSelector(lang, "different/declaration")
	missing := replacement
	missing.Selector = other
	for name, overrides := range map[string][]c.Override[source, source]{"missing": {{Binding: missing, ExpectedID: "base-v1"}}, "stale": {{Binding: replacement, ExpectedID: "base-v0"}}, "same-id": {{Binding: base, ExpectedID: "base-v1"}}, "duplicate": {{Binding: replacement, ExpectedID: "base-v1"}, {Binding: replacement, ExpectedID: "base-v1"}}} {
		t.Run(name, func(t *testing.T) {
			if _, err := c.NewBindings(lang, []c.Binding[source, source]{base}, overrides, nil); err == nil {
				t.Fatal("invalid override accepted")
			}
		})
	}
	if _, err := c.NewBindings(lang, []c.Binding[source, source]{base, base}, nil, nil); err == nil {
		t.Fatal("duplicate exact binding accepted")
	}
	wrong, _ := representation.NewLanguage[source](source{}, "v2")
	foreign := base
	foreign.Selector, _ = c.NewSelector(wrong, key.Key())
	if _, err := c.NewBindings(lang, []c.Binding[source, source]{foreign}, nil, nil); err == nil {
		t.Fatal("foreign source profile accepted")
	}
}

type renamed struct{}
func (renamed) ID() language.ID { return "independent" }
func (renamed) Name() string { return "New display name" }
func TestSelectorDisplayNameIsNotIdentity(t *testing.T) {
 a, _ := representation.NewLanguage[source](source{}, "v1")
 b, _ := representation.NewLanguage[source](renamed{}, "v1")
 ka, _ := c.NewSelector(a, "one")
 kb, _ := c.NewSelector(b, "one")
 mapping := c.Mapping[source,source]{ID:"base",PluginID:"owner",Plan:func(g.Declaration[source])(g.Plan[source,source],error){return g.Plan[source,source]{},nil}}
 entries := []c.Binding[source,source]{{Selector:ka,Mapping:mapping},{Selector:kb,Mapping:mapping}}
 if _,err:=c.NewBindings(a,entries,nil,nil);err==nil{t.Fatal("display name created a second identity")}
 replacement:=mapping;replacement.ID="new"
 if _,err:=c.NewBindings(a,entries[:1],[]c.Override[source,source]{{Binding:c.Binding[source,source]{Selector:kb,Mapping:replacement},ExpectedID:"base"}},nil);err!=nil{t.Fatal(err)}
}
