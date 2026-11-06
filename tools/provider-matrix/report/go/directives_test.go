package report

import (
	"go/ast"
	"strings"
	"testing"
)

func comments(lines ...string) *ast.CommentGroup {
	g := &ast.CommentGroup{}
	for _, line := range lines { g.List = append(g.List, &ast.Comment{Text: line}) }
	return g
}

func TestCapabilityDirectiveGrammar(t *testing.T) {
	ds, err := rawDirectives(comments(`// bittype:case target=p.Parse feature="mapped type" case=remap expects="shape preserved" profile="strict v1"`))
	if err != nil || len(ds) != 1 { t.Fatalf("%v %#v", err, ds) }
	if ds[0].values["feature"] != "mapped type" || ds[0].values["expects"] != "shape preserved" || !ds[0].profile.present { t.Fatal(ds) }
	ds, err = rawDirectives(comments("// ```go", "// bittype:grammar examples", "// ```", "//    bittype:grammar indented", "// bittype:grammar rules"))
	if err != nil || len(ds) != 1 || ds[0].argument != "rules" { t.Fatalf("examples became annotations: %#v %v", ds, err) }
}

func TestCapabilityDirectiveErrors(t *testing.T) {
	for _, input := range []string{
		"grammar", "grammar rules extra", "scope docs/a.md profile=x profile=x",
		"case target=p.Parse feature=f case=c", "case target=p.Parse feature=f case=c expects=ok unknown=x",
		`case target=p.Parse feature="" case=c expects=ok`, `scope "unterminated`, "unknown x",
		"grammar " + strings.Repeat("x", 4100),
	} {
		t.Run(input[:min(60,len(input))], func(t *testing.T) { if _, err := rawDirectives(comments("// bittype:"+input)); err == nil { t.Fatal("accepted malformed directive", input) } })
	}
}
