package report

import (
	"context"
	"fmt"
	"strings"
	"testing"

	capability "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
)

func TestPublicProviderLinksDistinguishSharedDelegation(t *testing.T) {
	p := fixture(t, map[string]string{
		"print/syntax/shared/provider/go": `package printing
import (id "example.test/fixture/contracts/language/go"; c "example.test/fixture/trees/cst/go")
type Provider[L any] struct{identity id.Provider}
func New[L any](identity id.Provider)(Provider[L],error){return Provider[L]{identity},nil}
func(p Provider[L])ID()id.ID{return p.identity.ID()}
func(p Provider[L])Name()string{return p.identity.Name()}
//line misleading.go:900
func(p Provider[L])Print(v c.CST[L])(string,error){return Print(v)}
func Print[L any](c.CST[L])(string,error){return "",nil}
`,
		"print/syntax/ember/provider/go": `package emberprovider
import (s "example.test/fixture/syntax/ember/go"; c "example.test/fixture/trees/cst/go"; printing "example.test/fixture/print/syntax/shared/provider/go")
type Provider struct{s.Language}
func(Provider)Print(v c.CST[s.Language])(string,error){return printing.Print(v)}
`,
	})
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	locations := m.rows["print.Provider[cst.CST[L], string]"].cells["ember"]
	if len(locations) != 2 {
		t.Fatal("fixture needs the two public providers", locations)
	}
	links, err := p.links(locations, "https://example.test/source", "details.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"[print/syntax/shared/provider/go.Provider&#91;syntax/ember/go.Language&#93;.Print](https://example.test/source/print/syntax/shared/provider/go/implementation.go#L8)",
		"[print/syntax/ember/provider/go.Provider.Print](https://example.test/source/print/syntax/ember/provider/go/implementation.go#L4)",
	} {
		if !strings.Contains(links, want) {
			t.Errorf("missing distinct public-method link %s in %s", want, links)
		}
	}
	if strings.Contains(links, "print/syntax/shared/provider/go/implementation.go#L9") || strings.Contains(links, "misleading.go") {
		t.Fatal("shared delegate used as the public source link", links)
	}
	a := &annotations{operations: map[capability.Symbol][]authored{}}
	details, err := p.renderDetails(context.Background(), m, a, Config{}, "https://example.test/source", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, loc := range locations {
		if len(loc.op.chain) != 1 || loc.op.chain[0].callee.Name() != "Print" || loc.op.chain[0].callee.Pkg().Path() != p.module+"/print/syntax/shared/provider/go" {
			t.Fatal("shared delegation was lost", loc.op)
		}
		if !strings.Contains(links, "details.md#cap-"+loc.op.id) || !strings.Contains(details, fmt.Sprintf("id=\"cap-%s\"", loc.op.id)) {
			t.Fatal("operation-specific evidence link changed")
		}
	}
	if !strings.Contains(details, "Static call relationships (no evidence transfer)") || !strings.Contains(details, "print/syntax/shared/provider/go/implementation.go#L9") {
		t.Fatal("delegate source missing from details")
	}
	providers, _, err := p.render(m, "https://example.test/source", "details.md")
	if err != nil {
		t.Fatal(err)
	}
	grid, err := collectProviderGrid(m)
	if err != nil {
		t.Fatal(err)
	}
	_, overview, err := p.renderProviderGrid(grid, m.languages[0], "https://example.test/source", "details.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(providers, links) || !strings.Contains(overview, links) {
		t.Fatal("matrix and overview use different public links")
	}
}

func TestPublicLinksCannotFallBackToTheDelegate(t *testing.T) {
	p := fixture(t)
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	op := *m.rows["print.Provider[cst.CST[L], string]"].cells["ember"][0].op
	if op.entry == nil {
		t.Fatal("fixture needs a retained delegate")
	}
	op.public = nil
	for _, loc := range []location{{op: &op}, {}} {
		if _, err := p.links([]location{loc}, "https://example.test/source"); err == nil {
			t.Fatal("missing public identity silently substituted")
		}
	}
}
