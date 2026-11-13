package grammar_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/text/go"
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/model/syntax/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// toy is a language of this test only. Its table is written as a language
// writes one: a composite literal of the model type, with predicates named.
type toy struct{}

const (
	kindList ast.Kind[toy] = "list"
	kindItem ast.Kind[toy] = "item"
	kindNote ast.Kind[toy] = "note"
)

func word() syntax.Predicate    { return syntax.Predicate{Owner: "toy", Name: "word"} }
func unbound() syntax.Predicate { return syntax.Predicate{Owner: "toy", Name: "unbound"} }

func toyRules() syntax.Rules[toy] {
	return syntax.Rules[toy]{
		kindList: {
			Attributes: map[string]syntax.Attribute{"title": {Required: true, Predicate: syntax.Nonempty()}, "free": {}},
			Roles: []syntax.Role[toy]{
				{Name: "head", Max: 1, Kinds: []ast.Kind[toy]{kindNote}},
				{Name: "item", Min: 1, Max: -1, Kinds: []ast.Kind[toy]{kindItem}},
			},
		},
		kindItem: {
			Attributes: map[string]syntax.Attribute{
				"name":  {Required: true, Predicate: word()},
				"mode":  {OneOf: []string{"on", "off"}},
				"count": {Predicate: syntax.Decimal(), OneOf: []string{"1", "2", "x"}},
				"flag":  {Predicate: syntax.Bool()},
			},
			Reference: true,
		},
		kindNote: {},
	}
}

func toyEvaluators() g.Evaluators {
	return g.Evaluators{word(): func(value string) bool { return value != "" && strings.ToLower(value) == value }}
}

// toyTree is list(title) with one item(name).
func toyTree() ast.AST[toy] {
	return ast.AST[toy]{ID: "toy", Root: 2, Nodes: []ast.Node[toy]{
		{Kind: kindItem, Attributes: []ast.Attribute{{Name: "name", Text: "first"}}},
		{Kind: kindList, Attributes: []ast.Attribute{{Name: "title", Text: "Things"}}, Children: []ast.Child{{Role: "item", Node: 1}}},
	}}
}

func code(err error) text.Code {
	var diagnostic *text.Diagnostic
	if errors.As(err, &diagnostic) {
		return diagnostic.Code
	}
	return ""
}

func TestCheckAppliesTableOfModelType(t *testing.T) {
	if err := g.Check(toyTree(), kindList, toyRules(), toyEvaluators()); err != nil {
		t.Fatal(err)
	}
	item := func(attributes ...ast.Attribute) ast.AST[toy] {
		tree := toyTree()
		tree.Nodes[0].Attributes = append(tree.Nodes[0].Attributes, attributes...)
		return tree
	}
	list := func(change func(*ast.AST[toy])) ast.AST[toy] {
		tree := toyTree()
		change(&tree)
		return tree
	}
	for name, test := range map[string]struct {
		tree ast.AST[toy]
		want string
	}{
		"a language's predicate admits":          {item(ast.Attribute{Name: "mode", Text: "on"}), ""},
		"a language's predicate refuses":         {list(func(tree *ast.AST[toy]) { tree.Nodes[0].Attributes[0].Text = "First" }), "unknown or invalid attribute: name"},
		"a shared predicate refuses":             {list(func(tree *ast.AST[toy]) { tree.Nodes[1].Attributes[0].Text = "" }), "unknown or invalid attribute: title"},
		"a shared predicate admits":              {item(ast.Attribute{Name: "flag", Text: "false"}), ""},
		"a shared predicate refuses a text":      {item(ast.Attribute{Name: "flag", Text: "no"}), "unknown or invalid attribute: flag"},
		"a set admits its member":                {item(ast.Attribute{Name: "mode", Text: "off"}), ""},
		"a set refuses another text":             {item(ast.Attribute{Name: "mode", Text: "auto"}), "unknown or invalid attribute: mode"},
		"a set refuses a member in another case": {item(ast.Attribute{Name: "mode", Text: "On"}), "unknown or invalid attribute: mode"},
		"a set refuses a prefix of a member":     {item(ast.Attribute{Name: "mode", Text: "of"}), "unknown or invalid attribute: mode"},
		"a set refuses an extension of a member": {item(ast.Attribute{Name: "mode", Text: "offs"}), "unknown or invalid attribute: mode"},
		"a set refuses the empty text":           {item(ast.Attribute{Name: "mode", Text: ""}), "unknown or invalid attribute: mode"},
		"predicate and set both admit":           {item(ast.Attribute{Name: "count", Text: "2"}), ""},
		"the set admits and the predicate not":   {item(ast.Attribute{Name: "count", Text: "x"}), "unknown or invalid attribute: count"},
		"the predicate admits and the set not":   {item(ast.Attribute{Name: "count", Text: "3"}), "unknown or invalid attribute: count"},
		"an entry without a constraint admits": {list(func(tree *ast.AST[toy]) {
			tree.Nodes[1].Attributes = append(tree.Nodes[1].Attributes, ast.Attribute{Name: "free", Text: "\x00 any"})
		}), ""},
		"an attribute outside the table":  {item(ast.Attribute{Name: "other", Text: "x"}), "unknown or invalid attribute: other"},
		"a required attribute is missing": {list(func(tree *ast.AST[toy]) { tree.Nodes[1].Attributes = nil }), "missing attribute: title"},
		"a role has too few children":     {ast.AST[toy]{ID: "toy", Root: 1, Nodes: []ast.Node[toy]{{Kind: kindList, Attributes: []ast.Attribute{{Name: "title", Text: "T"}}}}}, "wrong number of item children"},
		"a role outside the table":        {list(func(tree *ast.AST[toy]) { tree.Nodes[1].Children[0].Role = "entry" }), "unknown or out-of-order child role: entry"},
		"a kind the role does not hold":   {list(func(tree *ast.AST[toy]) { tree.Nodes[0] = ast.Node[toy]{Kind: kindNote} }), "item cannot contain note"},
		"a kind outside the table":        {list(func(tree *ast.AST[toy]) { tree.Nodes[0].Kind = "stray" }), "unsupported node kind: stray"},
		"a reference where the kind refers": {list(func(tree *ast.AST[toy]) {
			tree.Nodes[0].References = []ast.Reference{{Role: "target", Target: ast.NodeRef{Tree: "other", Node: 1}}}
		}), ""},
		"a reference of another role": {list(func(tree *ast.AST[toy]) {
			tree.Nodes[0].References = []ast.Reference{{Role: "import", Target: ast.NodeRef{Tree: "other", Node: 1}}}
		}), "unsupported reference role: import"},
		"a reference where the kind does not": {list(func(tree *ast.AST[toy]) {
			tree.Nodes[1].References = []ast.Reference{{Role: "target", Target: ast.NodeRef{Tree: "other", Node: 1}}}
		}), "unsupported reference role: target"},
		"the root is of another kind":           {ast.AST[toy]{ID: "toy", Root: 1, Nodes: []ast.Node[toy]{{Kind: kindNote}}}, "wrong root kind"},
		"an optional role may be absent or one": {ast.AST[toy]{ID: "toy", Root: 3, Nodes: []ast.Node[toy]{{Kind: kindNote}, {Kind: kindItem, Attributes: []ast.Attribute{{Name: "name", Text: "a"}}}, {Kind: kindList, Attributes: []ast.Attribute{{Name: "title", Text: "T"}}, Children: []ast.Child{{Role: "head", Node: 1}, {Role: "item", Node: 2}}}}}, ""},
		"roles keep the table's order":          {ast.AST[toy]{ID: "toy", Root: 3, Nodes: []ast.Node[toy]{{Kind: kindNote}, {Kind: kindItem, Attributes: []ast.Attribute{{Name: "name", Text: "a"}}}, {Kind: kindList, Attributes: []ast.Attribute{{Name: "title", Text: "T"}}, Children: []ast.Child{{Role: "item", Node: 2}, {Role: "head", Node: 1}}}}}, "unknown or out-of-order child role: head"},
	} {
		err := g.Check(test.tree, kindList, toyRules(), toyEvaluators())
		if test.want == "" && err != nil {
			t.Errorf("%s: %v", name, err)
		}
		if test.want != "" && (code(err) != text.InvalidAST || !strings.HasSuffix(err.Error(), test.want)) {
			t.Errorf("%s: got %v, want invalid-ast %q", name, err, test.want)
		}
	}
}

// A table that names a predicate without an evaluator is not applied to any
// tree. The failure is unsupported: no tree was judged.
func TestCheckRefusesTableWithUnboundPredicate(t *testing.T) {
	rules := toyRules()
	rules[kindNote] = syntax.Rule[toy]{Attributes: map[string]syntax.Attribute{"text": {Predicate: unbound()}}}
	want := "unsupported: grammar predicate toy/unbound has no evaluator"
	// The valid tree holds no node of the kind whose entry names the predicate.
	// The other two are no trees at all: the table is refused before either is read.
	for name, tree := range map[string]ast.AST[toy]{
		"a valid tree that does not use the entry": toyTree(),
		"an empty tree":                   {},
		"a tree with a root out of range": {ID: "toy", Root: 7, Nodes: []ast.Node[toy]{{Kind: kindNote}}},
	} {
		err := g.Check(tree, kindList, rules, toyEvaluators())
		if code(err) != text.Unsupported || err.Error() != want {
			t.Errorf("%s: got %v, want %q", name, err, want)
		}
	}
	// A language's own name with no entry at all, and with an entry that holds no function.
	for name, evaluators := range map[string]g.Evaluators{
		"no evaluators":       nil,
		"an empty binding":    {},
		"an entry without fn": {word(): nil},
	} {
		err := g.Check(toyTree(), kindList, toyRules(), evaluators)
		if want := "unsupported: grammar predicate toy/word has no evaluator"; code(err) != text.Unsupported || err.Error() != want {
			t.Errorf("%s: got %v, want %q", name, err, want)
		}
	}
	// A name of the shared owner that the interpreter does not bind.
	rules = toyRules()
	rules[kindNote] = syntax.Rule[toy]{Attributes: map[string]syntax.Attribute{"text": {Predicate: syntax.Predicate{Owner: syntax.Owner, Name: "word"}}}}
	err := g.Check(toyTree(), kindList, rules, toyEvaluators())
	if want := "unsupported: grammar predicate bittype/word has no evaluator"; code(err) != text.Unsupported || err.Error() != want {
		t.Errorf("got %v, want %q", err, want)
	}
}

// A name is the pair of its owner and its name. A name with an empty owner is
// a name like any other: it is no absent predicate.
func TestCheckTakesNameWithEmptyOwnerAsName(t *testing.T) {
	bare := syntax.Predicate{Name: "word"}
	rules := toyRules()
	rule := rules[kindItem]
	rule.Attributes["name"] = syntax.Attribute{Required: true, Predicate: bare}
	err := g.Check(toyTree(), kindList, rules, toyEvaluators())
	if want := "unsupported: grammar predicate /word has no evaluator"; code(err) != text.Unsupported || err.Error() != want {
		t.Errorf("unbound: got %v, want %q", err, want)
	}
	evaluators := toyEvaluators()
	evaluators[bare] = func(value string) bool { return value == "second" }
	if err := g.Check(toyTree(), kindList, rules, evaluators); code(err) != text.InvalidAST || !strings.HasSuffix(err.Error(), "unknown or invalid attribute: name") {
		t.Errorf("bound and refusing: got %v", err)
	}
	tree := toyTree()
	tree.Nodes[0].Attributes[0].Text = "second"
	if err := g.Check(tree, kindList, rules, evaluators); err != nil {
		t.Errorf("bound and admitting: %v", err)
	}
}

// The report of a table's defect is determined. Of several unbound names the
// least by owner and then by name is reported, whatever order the table is
// read in. A rebound shared name comes before an unbound name, and an unbound
// name before an empty set.
func TestCheckReportsTheLeastDefect(t *testing.T) {
	rules := func() syntax.Rules[toy] {
		rules := toyRules()
		rules[kindNote] = syntax.Rule[toy]{Attributes: map[string]syntax.Attribute{
			"a": {Predicate: syntax.Predicate{Owner: "toy", Name: "b"}},
			"b": {Predicate: syntax.Predicate{Owner: "toy", Name: "a"}},
			"c": {Predicate: syntax.Predicate{Owner: "zoo", Name: "a"}},
			"d": {Predicate: syntax.Predicate{Owner: "other", Name: "z"}},
		}}
		rules["extra"] = syntax.Rule[toy]{Attributes: map[string]syntax.Attribute{"e": {Predicate: syntax.Predicate{Owner: "other", Name: "y"}}}}
		return rules
	}
	// Go visits a map in an order that differs from run to run, so the table is
	// read many times.
	for range 200 {
		err := g.Check(toyTree(), kindList, rules(), toyEvaluators())
		if want := "unsupported: grammar predicate other/y has no evaluator"; code(err) != text.Unsupported || err.Error() != want {
			t.Fatalf("got %v, want %q", err, want)
		}
	}
	table := rules()
	table["extra"].Attributes["f"] = syntax.Attribute{OneOf: []string{}}
	evaluators := toyEvaluators()
	if err := g.Check(toyTree(), kindList, table, evaluators); err == nil || err.Error() != "unsupported: grammar predicate other/y has no evaluator" {
		t.Errorf("an unbound name and an empty set: %v", err)
	}
	evaluators[syntax.Bool()] = func(string) bool { return true }
	if err := g.Check(toyTree(), kindList, table, evaluators); err == nil || err.Error() != "unsupported: grammar evaluators bind the shared predicate bittype/bool" {
		t.Errorf("a rebound shared name, an unbound name and an empty set: %v", err)
	}
}

// A set that is present and empty is a defect of the table: it would admit no
// text. The table is refused before a tree is read, as unsupported. An absent
// set, which is a nil slice, restricts nothing.
func TestCheckRefusesTableWithEmptySet(t *testing.T) {
	rules := toyRules()
	rules[kindItem].Attributes["mode"] = syntax.Attribute{OneOf: []string{}}
	want := "unsupported: grammar attribute item.mode holds an empty set"
	for name, tree := range map[string]ast.AST[toy]{
		"a valid tree that does not use the entry": toyTree(),
		"an empty tree": {},
	} {
		if err := g.Check(tree, kindList, rules, toyEvaluators()); code(err) != text.Unsupported || err.Error() != want {
			t.Errorf("%s: got %v, want %q", name, err, want)
		}
	}
	// Of several empty sets one is reported, the same in every run.
	rules[kindList].Attributes["free"] = syntax.Attribute{OneOf: make([]string, 0, 4)}
	rules[kindItem].Attributes["flag"] = syntax.Attribute{Predicate: syntax.Bool(), OneOf: []string{}}
	for range 200 {
		if err := g.Check(toyTree(), kindList, rules, toyEvaluators()); err == nil || err.Error() != "unsupported: grammar attribute item.flag holds an empty set" {
			t.Fatalf("got %v", err)
		}
	}
	// An absent set restricts nothing, with and without a predicate.
	rules = toyRules()
	rules[kindItem].Attributes["mode"] = syntax.Attribute{OneOf: nil}
	rules[kindItem].Attributes["flag"] = syntax.Attribute{Predicate: syntax.Bool(), OneOf: nil}
	tree := toyTree()
	tree.Nodes[0].Attributes = append(tree.Nodes[0].Attributes, ast.Attribute{Name: "mode", Text: "any text"}, ast.Attribute{Name: "flag", Text: "true"})
	if err := g.Check(tree, kindList, rules, toyEvaluators()); err != nil {
		t.Errorf("an absent set: %v", err)
	}
}

// The shared names have one binding. A caller cannot give one of them another
// evaluator, and cannot add a name under their owner.
func TestCheckRefusesSecondBindingOfSharedPredicate(t *testing.T) {
	for name, predicate := range map[string]syntax.Predicate{
		"a shared name":              syntax.Nonempty(),
		"a new name under the owner": {Owner: syntax.Owner, Name: "word"},
	} {
		evaluators := toyEvaluators()
		evaluators[predicate] = func(string) bool { return true }
		err := g.Check(toyTree(), kindList, toyRules(), evaluators)
		want := "unsupported: grammar evaluators bind the shared predicate " + predicate.Owner + "/" + predicate.Name
		if code(err) != text.Unsupported || err.Error() != want {
			t.Errorf("%s: got %v, want %q", name, err, want)
		}
	}
}

// Each shared name is bound to the evaluator that states it.
func TestSharedPredicatesAreBound(t *testing.T) {
	for _, test := range []struct {
		predicate      syntax.Predicate
		admits, refuse []string
	}{
		{syntax.Bool(), []string{"true", "false"}, []string{"", "True", "1"}},
		{syntax.Decimal(), []string{"0", "7", "120"}, []string{"", "07", "-1", "1.0"}},
		{syntax.Nonempty(), []string{"a", " "}, []string{"", "\xff"}},
		{syntax.Positive(), []string{"1", "120"}, []string{"0", "", "-1", "01"}},
		{syntax.Signed(), []string{"0", "-1", "12"}, []string{"", "-", "+1", "-01"}},
		{syntax.String(), []string{"", "a b"}, []string{"a\x00b", "\xff"}},
	} {
		rules := syntax.Rules[toy]{kindNote: {Attributes: map[string]syntax.Attribute{"text": {Predicate: test.predicate}}}}
		check := func(value string) error {
			tree := ast.AST[toy]{ID: "toy", Root: 1, Nodes: []ast.Node[toy]{{Kind: kindNote, Attributes: []ast.Attribute{{Name: "text", Text: value}}}}}
			return g.Check(tree, kindNote, rules, nil)
		}
		for _, value := range test.admits {
			if err := check(value); err != nil {
				t.Errorf("%s refuses %q: %v", test.predicate.Name, value, err)
			}
		}
		for _, value := range test.refuse {
			if err := check(value); code(err) != text.InvalidAST {
				t.Errorf("%s admits %q: %v", test.predicate.Name, value, err)
			}
		}
	}
}
