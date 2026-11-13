// Package grammar shares checking mechanics, never a language's rules.
package grammar

import (
	"fmt"
	"github.com/Bitspark/bittype/check/syntax/shared/ast/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/model/syntax/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Evaluators binds predicate names to the functions that evaluate them. The
// package that calls Check supplies one entry for each predicate of its own
// language that its table names.
type Evaluators map[syntax.Predicate]func(string) bool

// shared binds the predicates that every language's table can use. No caller
// binds them a second time.
var shared = Evaluators{
	syntax.Bool():     Bool,
	syntax.Decimal():  Decimal,
	syntax.Nonempty(): Nonempty,
	syntax.Positive(): Positive,
	syntax.Signed():   Signed,
	syntax.String():   String,
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

// Words returns the evaluator of a word list: it reports whether a text is
// one of the words. A lexical helper builds one for each list that its
// language's lexical part declares.
func Words(words []string) func(string) bool {
	set := make(map[string]bool, len(words))
	for _, word := range words {
		set[word] = true
	}
	return func(s string) bool { return set[s] }
}

// evaluator returns the function bound to a name, or nil.
func evaluator(name syntax.Predicate, evaluators Evaluators) func(string) bool {
	if name.Owner == syntax.Owner {
		return shared[name]
	}
	return evaluators[name]
}

// bound reports whether every predicate that the table names has an
// evaluator, and whether the caller left the shared names alone. A table that
// fails here was not applied to any tree, so the failure is unsupported and
// never an invalid tree. Of several names the least is reported.
func bound[L any](rules syntax.Rules[L], evaluators Evaluators) error {
	var rebound, missing *syntax.Predicate
	least := func(found **syntax.Predicate, name syntax.Predicate) {
		if *found == nil || name.Owner < (*found).Owner || name.Owner == (*found).Owner && name.Name < (*found).Name {
			*found = &name
		}
	}
	for name := range evaluators {
		if name.Owner == syntax.Owner {
			least(&rebound, name)
		}
	}
	if rebound != nil {
		return &text.Diagnostic{Code: text.Unsupported, Message: "grammar evaluators bind the shared predicate " + rebound.Owner + "/" + rebound.Name}
	}
	for _, rule := range rules {
		for _, attribute := range rule.Attributes {
			if attribute.Predicate != (syntax.Predicate{}) && evaluator(attribute.Predicate, evaluators) == nil {
				least(&missing, attribute.Predicate)
			}
		}
	}
	if missing != nil {
		return &text.Diagnostic{Code: text.Unsupported, Message: "grammar predicate " + missing.Owner + "/" + missing.Name + " has no evaluator"}
	}
	return nil
}

// accepts reports whether a text satisfies an attribute's declaration.
func accepts(attribute syntax.Attribute, evaluators Evaluators, value string) bool {
	if attribute.Predicate != (syntax.Predicate{}) && !evaluator(attribute.Predicate, evaluators)(value) {
		return false
	}
	return len(attribute.OneOf) == 0 || slices.Contains(attribute.OneOf, value)
}

// Check applies one language's grammar table to a tree. The table names its
// predicates; evaluators binds the names of the caller's language, and the
// shared names are bound here. A table that names a predicate without an
// evaluator fails every check before the tree is read.
func Check[L any](tree ast.AST[L], root ast.Kind[L], rules syntax.Rules[L], evaluators Evaluators) error {
	if err := bound(rules, evaluators); err != nil {
		return err
	}
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
			if !ok || !accepts(spec, evaluators, a.Text) {
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
