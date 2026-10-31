package golangcheck

import (
	"go/token"
	"github.com/Bitspark/bittype/ast/go"
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/text/go"
)

func identifier(v string) bool { return g.Identifier(v) && !token.Lookup(v).IsKeyword() }
var types = []string{"named-type","qualified-type","type-application","struct-type","interface-type","function-type","pointer-type","slice-type","array-type","map-type","channel-type"}
var constraints = append(append([]string{}, types...), "type-union","type-approximation")
var rules = g.Rules{
	"file": {Attributes:g.Attributes("package",g.Attr(true,identifier)),Roles:[]g.Role{g.Many("import",0,"import"),g.Many("declaration",0,"type-declaration","type-alias")}},
	"import": {Attributes:g.Attributes("path",g.Attr(true,g.Nonempty),"alias",g.Attr(false,func(v string)bool{return v=="."||identifier(v)}))},
	"type-declaration": {Attributes:g.Attributes("name",g.Attr(true,identifier)),Roles:[]g.Role{g.Many("type-parameter",0,"type-parameter"),g.One("type",types...)}},
	"type-alias": {Attributes:g.Attributes("name",g.Attr(true,identifier)),Roles:[]g.Role{g.Many("type-parameter",0,"type-parameter"),g.One("type",types...)}},
	"type-parameter": {Attributes:g.Attributes("name",g.Attr(true,identifier)),Roles:[]g.Role{g.One("constraint",constraints...)}},
	"named-type": {Attributes:g.Attributes("name",g.Attr(true,identifier)),Reference:true},
	"qualified-type": {Attributes:g.Attributes("name",g.Attr(true,identifier),"qualifier",g.Attr(true,identifier)),Reference:true},
	"type-application": {Roles:[]g.Role{g.One("constructor","named-type","qualified-type"),g.Many("argument",1,types...)}},
	"struct-type": {Roles:[]g.Role{g.Many("field",0,"field","embedded-field")}},
	"field": {Attributes:g.Attributes("name",g.Attr(true,identifier),"tag",g.Attr(false,g.String)),Roles:[]g.Role{g.One("type",types...)}},
	"embedded-field": {Attributes:g.Attributes("tag",g.Attr(false,g.String)),Roles:[]g.Role{g.One("type","named-type","qualified-type","type-application","pointer-type")}},
	"interface-type": {Roles:[]g.Role{g.Many("member",0,"method","embedded-type")}},
	"method": {Attributes:g.Attributes("name",g.Attr(true,identifier)),Roles:[]g.Role{g.One("signature","function-type")}},
	"embedded-type": {Roles:[]g.Role{g.One("type",constraints...)}},
	"function-type": {Roles:[]g.Role{g.Many("parameter",0,"parameter"),g.Many("result",0,"parameter")}},
	"parameter": {Attributes:g.Attributes("name",g.Attr(false,identifier),"variadic",g.Attr(false,g.Bool)),Roles:[]g.Role{g.One("type",types...)}},
	"pointer-type": {Roles:[]g.Role{g.One("element",types...)}},
	"slice-type": {Roles:[]g.Role{g.One("element",types...)}},
	"array-type": {Attributes:g.Attributes("length",g.Attr(true,g.Decimal)),Roles:[]g.Role{g.One("element",types...)}},
	"map-type": {Roles:[]g.Role{g.One("key",types...),g.One("element",types...)}},
	"channel-type": {Attributes:g.Attributes("direction",g.Attr(false,g.Set("both","send","receive"))),Roles:[]g.Role{g.One("element",types...)}},
	"type-union": {Roles:[]g.Role{g.Many("term",2,append(types,"type-approximation")...)}},
	"type-approximation": {Roles:[]g.Role{g.One("type",types...)}},
}

func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t,"file",rules); err != nil { return err }
	v := tree.View[s.Language]{Tree:t}
	for i,n := range t.Nodes {
		id := ast.NodeID(i+1)
		if n.Kind == s.KindFunctionType {
			for _, role := range []string{"parameter","result"} {
				params := v.Children(id,role); named := -1
				for j,p := range params {
					mode := 0; if v.Has(p,"name"){mode=1}
					if named >= 0 && named != mode{return text.AtNode(t,p,"cannot mix named and unnamed parameters")}; named=mode
					if v.Flag(p,"variadic") && (role!="parameter" || j!=len(params)-1) {return text.AtNode(t,p,"variadic parameter must be final input")}
				}
			}
		}
		if n.Kind == s.KindEmbeddedField {
			e:=v.Child(id,"type")
			if v.Kind(e)=="pointer-type" {
				kind:=v.Kind(v.Child(e,"element"))
				if kind!="named-type" && kind!="qualified-type" && kind!="type-application" {return text.AtNode(t,e,"embedded pointer requires a named type")}
			}
		}
	}
	return nil
}

type Provider struct { s.Language }
func (Provider) Check(v ast.View[s.Language]) error {return Check(v)}
var _ text.CheckerProvider[s.Language] = Provider{}

