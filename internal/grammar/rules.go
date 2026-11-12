// Package grammar shares checking mechanics, never a language's rules.
package grammar

import (
	"fmt"
	"github.com/Bitspark/bittype/check/syntax/shared/ast/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Attribute struct {
	Required bool
	Accept   func(string) bool
}
type Role[L any] struct {
	Name     string
	Min, Max int
	Kinds    []ast.Kind[L]
}
type Rule[L any] struct {
	Attributes map[string]Attribute
	Roles      []Role[L]
	Reference  bool
}

// Rules is one language's admission table. Its keys and the kinds of its roles
// are that language's node kinds, so a table of one language cannot be applied
// to a tree of another. Role and attribute names are strings.
type Rules[L any] map[ast.Kind[L]]Rule[L]

func Attr(required bool, accept func(string) bool) Attribute { return Attribute{required, accept} }
func One[L any](name string, kinds ...ast.Kind[L]) Role[L]   { return Role[L]{name, 1, 1, kinds} }
func Opt[L any](name string, kinds ...ast.Kind[L]) Role[L]   { return Role[L]{name, 0, 1, kinds} }
func Many[L any](name string, min int, kinds ...ast.Kind[L]) Role[L] {
	return Role[L]{name, min, -1, kinds}
}
func Set(values ...string) func(string) bool {
	return func(s string) bool {
		for _, v := range values {
			if s == v {
				return true
			}
		}
		return false
	}
}
func Bool(s string) bool     { return s == "true" || s == "false" }
func Nonempty(s string) bool { return s != "" && utf8.ValidString(s) }
func String(s string) bool   { return utf8.ValidString(s) && !strings.ContainsRune(s, 0) }

var decimal = regexp.MustCompile("^(0|[1-9][0-9]*)$")
var signed = regexp.MustCompile("^-?(0|[1-9][0-9]*)$")

func Decimal(s string) bool  { return decimal.MatchString(s) }
func Positive(s string) bool { return Decimal(s) && s != "0" }
func Signed(s string) bool   { return signed.MatchString(s) }
func Identifier(s string) bool {
	if s == "" || !utf8.ValidString(s) {
		return false
	}
	for i, r := range s {
		if r != '_' && !unicode.IsLetter(r) && !(i > 0 && unicode.IsDigit(r)) {
			return false
		}
	}
	return true
}
func Attributes(pairs ...any) map[string]Attribute {
	m := map[string]Attribute{}
	for i := 0; i < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1].(Attribute)
	}
	return m
}

func Check[L any](tree ast.AST[L], root ast.Kind[L], rules Rules[L]) error {
	if err := check.Structure(tree); err != nil {
		return err
	}
	if tree.Nodes[tree.Root-1].Kind != root {
		return text.AtNode(tree, tree.Root, "wrong root kind")
	}
	for i, n := range tree.Nodes {
		id := ast.NodeID(i + 1)
		bad := func(message string) error { return text.AtNode(tree, id, message) }
		rule, ok := rules[n.Kind]
		if !ok {
			return bad("unsupported node kind: " + string(n.Kind))
		}
		present := map[string]bool{}
		for _, a := range n.Attributes {
			spec, ok := rule.Attributes[a.Name]
			if !ok || (spec.Accept != nil && !spec.Accept(a.Text)) {
				return bad("unknown or invalid attribute: " + a.Name)
			}
			present[a.Name] = true
		}
		for name, a := range rule.Attributes {
			if a.Required && !present[name] {
				return bad("missing attribute: " + name)
			}
		}
		counts := make([]int, len(rule.Roles))
		last := -1
		for _, c := range n.Children {
			index := -1
			for j, r := range rule.Roles {
				if c.Role == r.Name {
					index = j
					break
				}
			}
			if index < 0 || index < last {
				return bad("unknown or out-of-order child role: " + c.Role)
			}
			last = index
			counts[index]++
			childKind := tree.Nodes[c.Node-1].Kind
			allowed := false
			for _, kind := range rule.Roles[index].Kinds {
				if kind == childKind {
					allowed = true
					break
				}
			}
			if !allowed {
				return bad(fmt.Sprintf("%s cannot contain %s", c.Role, childKind))
			}
		}
		for j, r := range rule.Roles {
			if counts[j] < r.Min || (r.Max >= 0 && counts[j] > r.Max) {
				return bad("wrong number of " + r.Name + " children")
			}
		}
		for _, ref := range n.References {
			if !rule.Reference || ref.Role != "target" {
				return bad("unsupported reference role: " + ref.Role)
			}
		}
	}
	return nil
}
