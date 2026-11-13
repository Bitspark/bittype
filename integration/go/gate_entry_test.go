package integration_test

// Row 39 of the contribution contract (bittype#196), the second layer of
// "What a test must show": "Mutants must disagree on presence, layout, slot
// identity or environment." The test here compares the plan entry of each
// record and each application with the generated native declaration. It reads
// the plan as a reader does: the frozen payloads through their codecs, the
// links, the member lists and the ties of an application's arguments to
// binders. It reads the native side from the generated text. The two sides
// meet in this file only.

import (
	"context"
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
)

// gateNative is what a generated native text declares: for each record
// declaration its members in order, for each declaration its type parameters
// in order, and for every other declaration the type it is declared as.
type gateNative struct {
	records    map[string][]gateNativeMember
	parameters map[string][]string
	aliases    map[string]string
}

type gateNativeMember struct {
	name, typ string
	optional  bool
}

func gateNativeGo(t *testing.T, text string) gateNative {
	t.Helper()
	parsed := gateParseGo(t, text)
	out := gateNative{map[string][]gateNativeMember{}, parsed.parameters, parsed.aliases}
	for name, fields := range parsed.fields {
		out.records[name] = []gateNativeMember{}
		for _, field := range fields {
			out.records[name] = append(out.records[name], gateNativeMember{field.name, field.typ, false})
		}
	}
	return out
}

var (
	gateTSInterface = regexp.MustCompile(`^(?:export )?interface (\w+)(?:<([^>]*)>)? \{$`)
	gateTSMember    = regexp.MustCompile(`^  (\w+)(\?)?: (.*);$`)
	gateTSAlias     = regexp.MustCompile(`^(?:export )?type (\w+)(?:<([^>]*)>)? = (.*);$`)
)

// gateNativeTS reads the declarations of generated TypeScript text as the
// formatter of this repository prints them: an interface with one member on a
// line for a record, and a type alias on one line for every other declaration.
// A line it does not read fails the test, so no declaration is passed over.
func gateNativeTS(t *testing.T, text string) gateNative {
	t.Helper()
	out := gateNative{map[string][]gateNativeMember{}, map[string][]string{}, map[string]string{}}
	parameters := func(name, list string) {
		if list != "" {
			out.parameters[name] = strings.Split(list, ", ")
		}
	}
	open := ""
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimRight(line, "\r")
		switch {
		case open != "" && line == "}":
			open = ""
		case open != "":
			m := gateTSMember.FindStringSubmatch(line)
			if m == nil {
				t.Fatalf("unread member line %q in\n%s", line, text)
			}
			out.records[open] = append(out.records[open], gateNativeMember{m[1], m[3], m[2] == "?"})
		case gateTSInterface.MatchString(line):
			m := gateTSInterface.FindStringSubmatch(line)
			open = m[1]
			out.records[open] = []gateNativeMember{}
			parameters(m[1], m[2])
		case gateTSAlias.MatchString(line):
			m := gateTSAlias.FindStringSubmatch(line)
			out.aliases[m[1]] = m[3]
			parameters(m[1], m[2])
		case strings.TrimSpace(line) != "":
			t.Fatalf("unread line %q in\n%s", line, text)
		}
	}
	return out
}

// gateEntryReader is what a reader of one composition knows about its
// contributions: how each selected representation is spelled in the target.
// It is the test's own statement of that, made apart from the projectors.
type gateEntryReader struct {
	codecs []g.PayloadCodec
	// use spells the target type of a selected use from the local key of its
	// contribution, the decoded view of its frozen payload and the spellings
	// of its linked dependencies.
	use func(key string, view any, links []string) (string, error)
	// layout gives the members of a record as its entry states them: from the
	// local key of the record's contribution, the decoded view of its frozen
	// payload and its links in order. Each member names the link that its type
	// belongs to.
	layout func(key string, view any, links []gateLinked) ([]gateLaid, error)
	// position is the place of a source binder among the target parameters of
	// a declaration, read from the parameters that its frozen record states.
	position func(parameters []string, binder int) (int, error)
	// apply spells an application of a declaration to target arguments.
	apply func(name string, args []string) string
}

// gateLinked is one link of a record's entry as the reader resolved it from
// the plan: for a child the type selected for it, and for a field position
// the local key of the contribution selected for it, the member's name that
// the request states and the type selected for its operand.
type gateLinked struct {
	index   int
	kind    string // child | field
	spelled string
	stated  string
	field   string
}

// gateLaid is one member of a record as its entry gives it. link is the link
// that the member's type belongs to, or -1 for a member whose type the
// contribution's revision fixes.
type gateLaid struct {
	member    gateNativeMember
	link      int
	delegated bool
}

type gateDeclared struct {
	name       string
	nominal    bool
	parameters []string
}

type gateEntryCounts struct{ records, own, delegated, aliases, applications, binders int }

func (r gateEntryReader) view(plan g.RepresentationPlan, n g.PlanNode) (any, error) {
	if n.Schema == nil {
		return nil, fmt.Errorf("node %s of kind %s has no payload schema", n.ID, n.Kind)
	}
	i := slices.IndexFunc(r.codecs, func(c g.PayloadCodec) bool { return c.Schema() == *n.Schema })
	if i < 0 {
		return nil, fmt.Errorf("the reader has no codec for schema %s/%s", n.Schema.Identity, n.Schema.Version)
	}
	frozen, err := hex.DecodeString(n.Payload)
	if err != nil {
		return nil, err
	}
	if err := g.CheckReceipt(n.Receipt, *n.Schema, frozen); err != nil {
		return nil, err
	}
	return r.codecs[i].View(frozen)
}

// declared reads the frozen record of a declaration: its name, whether it is
// a defined declaration, and its target parameters in their order.
func (r gateEntryReader) declared(plan g.RepresentationPlan, n g.PlanNode) (gateDeclared, error) {
	if n.Kind != "declaration" {
		return gateDeclared{}, fmt.Errorf("node %s is a %s and no declaration", n.ID, n.Kind)
	}
	view, err := r.view(plan, n)
	if err != nil {
		return gateDeclared{}, err
	}
	v := reflect.ValueOf(view)
	name, nominal, parameters := v.FieldByName("Name"), v.FieldByName("Nominal"), v.FieldByName("Parameters")
	if !name.IsValid() || !nominal.IsValid() || !parameters.IsValid() {
		return gateDeclared{}, fmt.Errorf("the declaration payload %T has not the fields Name, Nominal and Parameters", view)
	}
	return gateDeclared{name.String(), nominal.Bool(), slices.Clone(parameters.Interface().([]string))}, nil
}

// spell spells the target type of a node from the plan alone. With adapters
// it follows the adapter arguments of each application, and otherwise the
// target arguments. Each argument is placed by the binder it is tied to.
func (r gateEntryReader) spell(plan g.RepresentationPlan, id string, adapters bool, counts *gateEntryCounts) (string, error) {
	n, ok := plan.Node(id)
	if !ok {
		return "", fmt.Errorf("the plan has no node %q", id)
	}
	switch n.Kind {
	case "binder":
		scope, ok := plan.Symbol(n.Binder.Scope)
		if !ok {
			return "", fmt.Errorf("the binder's scope %q is no symbol of the plan", n.Binder.Scope)
		}
		declaration, _ := plan.Node(scope.Node)
		declared, err := r.declared(plan, declaration)
		if err != nil {
			return "", err
		}
		position, err := r.position(declared.parameters, n.Binder.Index)
		if err != nil {
			return "", err
		}
		counts.binders++
		return declared.parameters[position], nil
	case "application":
		declaration, _ := plan.Node(n.Application.Declaration)
		declared, err := r.declared(plan, declaration)
		if err != nil {
			return "", err
		}
		arguments, what := n.Application.Target, "target"
		if adapters {
			arguments, what = n.Application.Adapters, "adapter"
		}
		if len(arguments) != declaration.Binders {
			return "", fmt.Errorf("the application of %s has %d %s arguments for %d binders", declared.name, len(arguments), what, declaration.Binders)
		}
		args := make([]string, len(declared.parameters))
		for binder := 0; binder < declaration.Binders; binder++ {
			tied := slices.IndexFunc(arguments, func(a g.PlanArgument) bool { return a.Binder == binder })
			if tied < 0 {
				return "", fmt.Errorf("the application of %s ties no %s argument to binder %d: %+v", declared.name, what, binder, arguments)
			}
			position, err := r.position(declared.parameters, binder)
			if err != nil {
				return "", err
			}
			if args[position], err = r.spell(plan, arguments[tied].Node, adapters, counts); err != nil {
				return "", err
			}
		}
		counts.applications++
		return r.apply(declared.name, args), nil
	case "use":
		var links []string
		for _, link := range n.Links {
			if link.Kind != "child" || link.Node == "" {
				return "", fmt.Errorf("the use %s has a link of kind %q that a type does not spell", n.Source, link.Kind)
			}
			spelled, err := r.spell(plan, link.Node, adapters, counts)
			if err != nil {
				return "", err
			}
			links = append(links, spelled)
		}
		view, err := r.view(plan, n)
		if err != nil {
			return "", err
		}
		return r.use(n.Contribution.Identity.Key, view, links)
	}
	return "", fmt.Errorf("a node of kind %s has no target type", n.Kind)
}

// record compares the entry of one record with the native record: which
// members are present, their order, which are delegated, and which link each
// member's type belongs to.
func (r gateEntryReader) record(t *testing.T, plan g.RepresentationPlan, name string, n g.PlanNode, native []gateNativeMember, adapters bool, counts *gateEntryCounts) {
	t.Helper()
	var nativeNames []string
	for _, member := range native {
		nativeNames = append(nativeNames, member.name)
	}
	// Presence and layout: the member list of the entry is the member list of
	// the native record.
	if !slices.Equal(n.Members, nativeNames) {
		t.Errorf("%s: the entry has the members %v and the native record has %v", name, n.Members, nativeNames)
		return
	}
	view, err := r.view(plan, n)
	if err != nil {
		t.Errorf("%s: %v", name, err)
		return
	}
	// The links of the entry, each resolved from the plan alone.
	var links []gateLinked
	for i, link := range n.Links {
		linked := gateLinked{index: i, kind: link.Kind, stated: link.Stated}
		operand := link.Node
		switch link.Kind {
		case "child":
		case "field":
			// A delegated member: the node of the field position has the
			// contribution selected for the role and one link, its operand.
			field, _ := plan.Node(link.Node)
			if field.Kind != "field" || field.Contribution == nil || len(field.Links) != 1 {
				t.Errorf("%s: link %d is delegated to %+v", name, i, field)
				return
			}
			linked.field, operand = field.Contribution.Identity.Key, field.Links[0].Node
		default:
			t.Errorf("%s: link %d has the kind %q", name, i, link.Kind)
			return
		}
		if linked.spelled, err = r.spell(plan, operand, adapters, counts); err != nil {
			t.Errorf("%s: link %d: %v", name, i, err)
			return
		}
		links = append(links, linked)
	}
	laid, err := r.layout(n.Contribution.Identity.Key, view, links)
	if err != nil {
		t.Errorf("%s: %v", name, err)
		return
	}
	if len(laid) != len(native) {
		t.Errorf("%s: the entry gives %d members and the native record has %d", name, len(laid), len(native))
		return
	}
	for i, member := range laid {
		got := native[i]
		// Presence and layout again, now from the frozen payload and the stated
		// names: the member at this place has this name.
		if member.member.name != got.name {
			t.Errorf("%s: member %d of the entry is %q (delegated: %v) and the native member is %q", name, i, member.member.name, member.delegated, got.name)
			continue
		}
		// Slot identity, and for a delegated member the form of the contribution
		// selected for the field position.
		if member.member.typ != got.typ || member.member.optional != got.optional {
			t.Errorf("%s.%s: link %d of the entry gives %s (optional %v, delegated %v), and the native member is %s (optional %v)", name, got.name, member.link, member.member.typ, member.member.optional, member.delegated, got.typ, got.optional)
		}
		if member.delegated {
			counts.delegated++
		} else {
			counts.own++
		}
	}
	counts.records++
}

// gateLinkPerMember is the layout of a record whose contribution names one
// request for each field (the horos record policy): a child link is a member
// that the record lays out itself, named by the next name of its frozen
// payload, and a field link is a delegated member, named in the request, in
// the form that field gives for the contribution selected for the position.
func gateLinkPerMember(record string, field func(key, operand string) (typ string, optional bool, err error)) func(string, any, []gateLinked) ([]gateLaid, error) {
	return func(key string, view any, links []gateLinked) ([]gateLaid, error) {
		if key != record {
			return nil, fmt.Errorf("the reader lays out no record of %s", key)
		}
		members := reflect.ValueOf(view).FieldByName("Members")
		if !members.IsValid() {
			return nil, fmt.Errorf("the record payload %T has no field Members", view)
		}
		own := members.Interface().([]string)
		var laid []gateLaid
		next := 0
		for _, link := range links {
			if link.kind == "field" {
				typ, optional, err := field(link.field, link.spelled)
				if err != nil {
					return nil, err
				}
				laid = append(laid, gateLaid{gateNativeMember{link.stated, typ, optional}, link.index, true})
				continue
			}
			if next >= len(own) {
				return nil, fmt.Errorf("link %d is a member the record lays out, and its payload names only %v", link.index, own)
			}
			laid = append(laid, gateLaid{gateNativeMember{own[next], link.spelled, false}, link.index, false})
			next++
		}
		if next != len(own) {
			return nil, fmt.Errorf("the payload names the members %v and the entry links %d of them", own, next)
		}
		return laid, nil
	}
}

// check compares every planned source declaration with the native declaration
// of its name.
func (r gateEntryReader) check(t *testing.T, plan g.RepresentationPlan, native gateNative, adapters bool) gateEntryCounts {
	t.Helper()
	var counts gateEntryCounts
	for _, symbol := range plan.Symbols() {
		if symbol.Role == "owned" {
			continue
		}
		declaration, _ := plan.Node(symbol.Node)
		declared, err := r.declared(plan, declaration)
		if err != nil {
			t.Errorf("%s: %v", symbol.ID, err)
			continue
		}
		name := declared.name
		if symbol.Intent != name {
			t.Errorf("%s: the symbol's intent is %q and the frozen record names %q", symbol.ID, symbol.Intent, name)
		}
		if !slices.Equal(native.parameters[name], declared.parameters) {
			t.Errorf("%s: the entry has the parameters %v and the native declaration has %v", name, declared.parameters, native.parameters[name])
		}
		if len(declaration.Links) != 1 {
			t.Errorf("%s: the declaration has %d links", name, len(declaration.Links))
			continue
		}
		body, _ := plan.Node(declaration.Links[0].Node)
		members, record := native.records[name]
		alias, aliased := native.aliases[name]
		switch {
		case declared.nominal && record:
			r.record(t, plan, name, body, members, adapters, &counts)
			if adapters && !t.Failed() {
				// Once more through the target arguments of each application,
				// where nothing differed so far. The counts are those of the first
				// pass.
				r.record(t, plan, name, body, members, false, &gateEntryCounts{})
			}
		case !declared.nominal && aliased:
			for _, byAdapters := range []bool{false, true} {
				if byAdapters && !adapters {
					continue
				}
				spelled, err := r.spell(plan, body.ID, byAdapters, &counts)
				if err != nil {
					t.Errorf("%s: %v", name, err)
					continue
				}
				if spelled != alias {
					t.Errorf("%s: the entry is %s (adapter arguments: %v) and the native declaration is %s", name, spelled, byAdapters, alias)
				}
			}
			counts.aliases++
		default:
			t.Errorf("%s: the entry is nominal=%v, and the native text has a record: %v, an alias: %v", name, declared.nominal, record, aliased)
		}
	}
	return counts
}

func gateBracketApply(open, close string) func(string, []string) string {
	return func(name string, args []string) string {
		if len(args) == 0 {
			return name
		}
		return name + open + strings.Join(args, ", ") + close
	}
}

func gateSamePosition(parameters []string, binder int) (int, error) {
	if binder < 0 || binder >= len(parameters) {
		return 0, fmt.Errorf("binder %d of a declaration with the parameters %v", binder, parameters)
	}
	return binder, nil
}

// gateGolangReader is the reader of the builtin horos-to-golang set.
func gateGolangReader(codecs []g.PayloadCodec) gateEntryReader {
	return gateEntryReader{
		codecs: codecs,
		use: func(key string, _ any, links []string) (string, error) {
			switch {
			case key == "horos-to-golang/text" && len(links) == 0:
				return "string", nil
			case key == "horos-to-golang/boolean" && len(links) == 0:
				return "bool", nil
			case key == "horos-to-golang/list" && len(links) == 1:
				return "[]" + links[0], nil
			}
			return "", fmt.Errorf("the reader spells no use of %s with %d links", key, len(links))
		},
		layout: gateLinkPerMember("horos-to-golang/record", func(key, operand string) (string, bool, error) {
			if key != "horos-to-golang/optional" {
				return "", false, fmt.Errorf("the reader spells no field of %s", key)
			}
			return "BittypeOptional[" + operand + "]", false, nil
		}),
		position: gateSamePosition,
		apply:    gateBracketApply("[", "]"),
	}
}

// gateTypescriptReader is the reader of the builtin horos-to-typescript set.
func gateTypescriptReader(codecs []g.PayloadCodec) gateEntryReader {
	return gateEntryReader{
		codecs: codecs,
		use: func(key string, _ any, links []string) (string, error) {
			switch {
			case key == "horos-to-typescript/text" && len(links) == 0:
				return "string", nil
			case key == "horos-to-typescript/boolean" && len(links) == 0:
				return "boolean", nil
			case key == "horos-to-typescript/list" && len(links) == 1:
				return links[0] + "[]", nil
			}
			return "", fmt.Errorf("the reader spells no use of %s with %d links", key, len(links))
		},
		layout: gateLinkPerMember("horos-to-typescript/record", func(key, operand string) (string, bool, error) {
			if key != "horos-to-typescript/optional" {
				return "", false, fmt.Errorf("the reader spells no field of %s", key)
			}
			return operand, true, nil
		}),
		position: gateSamePosition,
		apply:    gateBracketApply("<", ">"),
	}
}

// gateIndependentReader is the reader of the fixture composition for the
// independent source. The place of a binder among the target parameters is
// what the frozen record of the declaration states: the parameter named
// "T<binder>". The fixture's records record nothing: the layout of each is
// what its contribution's revision fixes, and only the types come from links.
func gateIndependentReader(codecs []g.PayloadCodec) gateEntryReader {
	return gateEntryReader{
		codecs: codecs,
		use: func(key string, view any, links []string) (string, error) {
			switch {
			case (key == "independent-word" || key == "independent-carry") && len(links) == 0:
				return "string", nil
			case key == "independent-cell" && len(links) == 1:
				return "[]" + links[0], nil
			case key == "independent-person" && len(links) == 0:
				// The record form where it is no body of a declaration.
				return "struct{Name string}", nil
			case key == "independent-probe" && len(links) == 1:
				if direct, ok := view.(gateDirect); ok && direct.Direct {
					return "[]" + links[0], nil
				}
				return "[][]" + links[0], nil
			}
			return "", fmt.Errorf("the reader spells no use of %s with %d links", key, len(links))
		},
		layout: func(key string, _ any, links []gateLinked) ([]gateLaid, error) {
			own := func(name, typ string, link int) gateLaid {
				return gateLaid{gateNativeMember{name, typ, false}, link, false}
			}
			for _, link := range links {
				if link.kind != "child" && key != "independent-outcome" {
					return nil, fmt.Errorf("the record of %s has a link of kind %s", key, link.kind)
				}
			}
			switch {
			case key == "independent-person" && len(links) == 0:
				return []gateLaid{own("Name", "string", -1)}, nil
			case key == "independent-page" && len(links) == 1:
				return []gateLaid{own("Items", "[]"+links[0].spelled, 0)}, nil
			case key == "independent-pair" && len(links) == 2:
				return []gateLaid{own("First", links[0].spelled, 0), own("Second", links[1].spelled, 1)}, nil
			case key == "independent-outcome":
				// The mock record delegates every field to the mock family, which
				// represents a field as a slice of its operand.
				var laid []gateLaid
				for _, link := range links {
					if link.kind != "field" || link.field != "mock-optional-nullable" {
						return nil, fmt.Errorf("link %d of the mock record is a %s of %q", link.index, link.kind, link.field)
					}
					laid = append(laid, gateLaid{gateNativeMember{link.stated, "[]" + link.spelled, false}, link.index, true})
				}
				return laid, nil
			}
			return nil, fmt.Errorf("the reader lays out no record of %s with %d links", key, len(links))
		},
		position: func(parameters []string, binder int) (int, error) {
			position := slices.Index(parameters, fmt.Sprintf("T%d", binder))
			if position < 0 {
				return 0, fmt.Errorf("no parameter of %v stands for binder %d", parameters, binder)
			}
			return position, nil
		},
		apply: gateBracketApply("[", "]"),
	}
}

// Row 39, the second layer ("Mutants that disagree"), gate line 1.
//
// Shows, for the Go and the TypeScript output of the pagination module, of a
// module with a generic of two binders, and for a fixture composition over
// the independent source whose declaration policy reverses the target
// parameters of one generic and whose mock record delegates every field: the
// plan entry and the generated native declaration agree in
//   - presence: the member list of a record's entry is the member list of the
//     native record, the members the record lays out are the names of its
//     frozen payload, and a member is delegated in the entry exactly where
//     the native member has the form of the contribution selected for the
//     field position;
//   - layout: the members and the target parameters are in the same order;
//   - slot identity: the type of each native member is the type of the link at
//     that member's place, spelled from the plan alone;
//   - environment: each target argument and each adapter argument of an
//     application, placed by the binder it is tied to, spells the native type
//     argument at that place, and a binder inside a declaration spells the
//     native parameter of that binder.
//
// Does not show: that the entry is right about the source. A rule that names
// its fields in another order gives an entry and a native type that agree with
// each other. It does not show a fragment of more than one member, and it
// reads the native side from the printed text and not from the native tree.
// The spellings are this test's own statement of what each contribution
// represents, so a contribution that changes its representation changes this
// reader too. The fixture's own records record nothing, so for them the
// presence and the layout rest on the member list of the entry alone. No
// module here has a record with a member that applies a generic of two
// binders: the ties of such an application are compared in the aliases.
//
// Red under, by kind. Presence: R39-presence-native-only,
// R39-presence-entry-only, R39-presence-golang-delegation,
// R39-presence-typescript-optional-form, R39-presence-fixture-entry. Layout:
// R39-layout-entry, R39-layout-native, R39-layout-native-parameters,
// R39-layout-golang-native, R39-layout-typescript-entry,
// R39-layout-fixture-entry. Slot identity: R39-slot-native-crossed,
// R39-slot-entry-crossed,
// R39-slot-golang-member-type, R39-slot-typescript-member-type. Environment:
// R39-environment-target-binder, R39-environment-adapter-binder,
// R39-environment-binder-node, R39-environment-native-arguments,
// R39-environment-fixture-arguments.
func TestGateEntryAgreesWithNativeDeclaration(t *testing.T) {
	type expectation struct{ records, own, delegated, aliases int }
	type planned struct {
		golang, typescript g.GenerationPlan
		goCodecs, tsCodecs []g.PayloadCodec
	}
	request := gateRequest{products: []string{gateBindings.Identity}}
	plans := func(t *testing.T, input transform.Input) planned {
		t.Helper()
		language := input.Source.Language()
		goBundle, tsBundle := gateGolang(t, language), gateTypescript(t, language)
		goBundle.products, tsBundle.products = []g.Product{gateBindings}, []g.Product{gateBindings}
		return planned{
			golang: gateMust[g.GenerationPlan](t)(gatePlanGo(t, input, goBundle, request)), typescript: gateMust[g.GenerationPlan](t)(gatePlanTS(t, input, tsBundle, request)),
			goCodecs: goBundle.codecs, tsCodecs: tsBundle.codecs,
		}
	}
	_, _, pagination := generationInput(t, paginationModule())
	_, twoBinders := gateSource(t, gatePagination)
	for _, module := range []struct {
		name  string
		input transform.Input
		want  expectation
	}{
		{"pagination", pagination, expectation{records: 2, own: 3, delegated: 1, aliases: 4}},
		{"two-binders", twoBinders, expectation{records: 3, own: 5, delegated: 2, aliases: 4}},
	} {
		p := plans(t, module.input)
		for _, target := range []struct {
			name   string
			plan   g.RepresentationPlan
			native func(t *testing.T) gateNative
			reader gateEntryReader
		}{
			{"golang", p.golang.Representation(), func(t *testing.T) gateNative {
				return gateNativeGo(t, gateText[gs.Language](t, golangText(), p.golang))
			}, gateGolangReader(p.goCodecs)},
			{"typescript", p.typescript.Representation(), func(t *testing.T) gateNative {
				return gateNativeTS(t, gateText[ts.Language](t, typescriptText(), p.typescript))
			}, gateTypescriptReader(p.tsCodecs)},
		} {
			t.Run(module.name+"/"+target.name, func(t *testing.T) {
				counts := target.reader.check(t, target.plan, target.native(t), true)
				if got := (expectation{counts.records, counts.own, counts.delegated, counts.aliases}); got != module.want {
					t.Errorf("compared %+v, expected %+v", got, module.want)
				}
				if counts.applications == 0 || counts.binders == 0 {
					t.Errorf("compared %d applications and %d binders", counts.applications, counts.binders)
				}
				t.Logf("compared %d records with %d members they lay out and %d delegated members, %d aliases, %d applications and %d binders", counts.records, counts.own, counts.delegated, counts.aliases, counts.applications, counts.binders)
			})
		}
	}
	t.Run("fixture", func(t *testing.T) {
		fixture := gateFixture(t)
		composition := indComposition(t, fixture, indOptions{reversed: []string{"Pair"}})
		composition.Products, composition.Requested = []g.Product{gateBindings}, []string{gateBindings.Identity}
		plan, err := fixture.PlanContributions(context.Background(), composition, "Words", "People", "Cells", "Outcomes")
		if err != nil {
			t.Fatal(err)
		}
		text := gateText[gs.Language](t, golangText(), plan)
		native := gateNativeGo(t, text)
		// The decision that the source does not determine: Pair's target
		// parameters are in the reverse order of its binders.
		if !slices.Equal(native.parameters["Pair"], []string{"T1", "T0"}) || native.aliases["Words"] != "Pair[Person, Word]" {
			t.Fatalf("the fixture's reversed declaration is not generated as configured:\n%s", text)
		}
		counts := gateIndependentReader(composition.Codecs).check(t, plan.Representation(), native, true)
		if want := (gateEntryCounts{records: 4, own: 4, delegated: 2, aliases: 5}); counts.records != want.records || counts.own != want.own || counts.delegated != want.delegated || counts.aliases != want.aliases || counts.applications == 0 || counts.binders == 0 {
			t.Errorf("compared %+v, expected %+v and some applications and binders", counts, want)
		}
		t.Logf("compared %d records with %d members they lay out and %d delegated members, %d aliases, %d applications and %d binders", counts.records, counts.own, counts.delegated, counts.aliases, counts.applications, counts.binders)
	})
}
