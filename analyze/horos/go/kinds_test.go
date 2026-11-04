package analyze_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	parse "github.com/Bitspark/bittype/parse/horos/go"
	parsecontract "github.com/Bitspark/bittype/contracts/parse/go"
	provider "github.com/Bitspark/bittype/provider/horos/go"
	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	codec "github.com/bitspark/ontos/codec/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

func properClaim(claim string) value { return d("proper", atom(claim)) }
func constructorClaim(result string, inputs ...value) value {
	return d("constructor", data.EncodeList(inputs), atom(result))
}
func kindParameter(name string, contract value) value { return d("parameter", txt(name), contract) }
func kindDefinition(body value, parameters ...value) value {
	return d("type", data.EncodeList(parameters), atom("type"), body)
}
func applyExpression(callee value, arguments ...value) value {
	return d("apply", callee, data.EncodeList(arguments))
}
func oneFact(t *testing.T, model semantics.Model, subject semantics.Subject, relation semantics.Key) semantics.Term {
	t.Helper()
	values := model.Facts(subject, relation)
	if len(values) != 1 {
		t.Fatalf("expected one %v fact for %+v, got %v", relation, subject, values)
	}
	return term(t, model, values[0])
}

func summaryOperator(t *testing.T, model semantics.Model, summary semantics.Term, name string) semantics.Term {
	t.Helper()
	if summary.Operator.Name == name {
		return summary
	}
	if summary.Operator.Name == "data-summary/and" {
		for _, child := range summary.Operands {
			candidate := term(t, model, child.Term)
			if candidate.Operator.Name == name || candidate.Operator.Name == "data-summary/and" {
				return summaryOperator(t, model, candidate, name)
			}
		}
	}
	t.Fatalf("summary lacks %s: %+v", name, summary)
	return semantics.Term{}
}

func TestHigherOrderFactsRetainConstructorIdentityAndDataPrecision(t *testing.T) {
	typeInput := properClaim("type")
	dataInput := properClaim("data")
	ordinary := constructorClaim("type", typeInput)
	dataPreserving := constructorClaim("data", dataInput)
	wrappedBody := h("record", atom("closed"), h("field", txt("value"), atom("req"), applyExpression(param(0), param(1))))
	root := module("example.kinds", mapping(),
		entry("Page", kindDefinition(page(param(0)), kindParameter("T", typeInput))),
		entry("Ignore", kindDefinition(h("scalar", atom("utf8-text")), kindParameter("T", typeInput))),
		entry("Send", kindDefinition(h("wire", param(0)), kindParameter("T", typeInput))),
		entry("Wrapped", kindDefinition(wrappedBody, kindParameter("F", ordinary), kindParameter("T", typeInput))),
		entry("Compose", kindDefinition(applyExpression(param(0), applyExpression(param(1), param(2))), kindParameter("F", ordinary), kindParameter("G", ordinary), kindParameter("T", typeInput))),
		entry("Invoke", kindDefinition(applyExpression(param(0), param(1), param(2)), kindParameter("H", constructorClaim("type", ordinary, typeInput)), kindParameter("F", ordinary), kindParameter("T", typeInput))),
		entry("DataWrapped", kindDefinition(wrappedBody, kindParameter("F", dataPreserving), kindParameter("D", dataInput))),
		entry("Person", definition(person())),
		entry("WrappedPage", definition(local("Wrapped", local("Page"), local("Person")))),
		entry("WrappedIgnore", definition(local("Wrapped", local("Ignore"), h("wire", local("Person"))))),
		entry("WrappedSend", definition(local("Wrapped", local("Send"), local("Person")))),
		entry("Nested", definition(local("Compose", local("Page"), local("Page"), local("Person")))),
		entry("Invoked", definition(local("Invoke", local("Wrapped"), local("Page"), local("Person")))),
		entry("DataPeople", definition(local("DataWrapped", local("Page"), local("Person")))))
	tree, err := resolve.Resolve(root, nil, "higher-order")
	if err != nil {
		t.Fatal(err)
	}
	model := requireModel(t, run(t, tree, nil))
	for name, expected := range map[string]string{
		"Page": "conditional", "Ignore": "proven", "Send": "failed", "Wrapped": "conditional", "Compose": "conditional",
		"WrappedPage": "proven", "WrappedIgnore": "proven", "WrappedSend": "failed", "Nested": "proven", "Invoked": "proven", "DataPeople": "proven",
	} {
		id, _ := entity(t, model, "example.kinds", name)
		status := oneFact(t, model, semantics.Subject{Entity: id}, resolve.DataStatusRelation())
		if status.Operator.Name != "data-status/"+expected || !bytes.Equal(status.Data, codec.Encode(atom(expected))) {
			t.Fatalf("%s Data status: %+v", name, status)
		}
	}
	wrappedID, wrapped := entity(t, model, "example.kinds", "Wrapped")
	f := wrapped.Parameters[0].Entity
	parameter, _ := model.Entity(f)
	if len(parameter.Parameters) != 0 || wrapped.Parameters[0].Sort != (semantics.Key{Owner: "horos", Name: "constructor"}) {
		t.Fatal("constructor interface fabricated declaration parameters")
	}
	contract := oneFact(t, model, semantics.Subject{Entity: f}, resolve.InterfaceRelation())
	if contract.Operator.Name != "interface/constructor" || len(contract.Operands) != 1 || len(contract.Operands[0].Locator) != 1 || contract.Operands[0].Locator[0] != 0 {
		t.Fatal("constructor interface structure missing")
	}
	input := term(t, model, contract.Operands[0].Term)
	if input.Operator.Name != "interface/proper" || !bytes.Equal(input.Data, codec.Encode(atom("type"))) {
		t.Fatal("constructor input admission was not retained")
	}
	bodyTerm := term(t, model, wrapped.Body)
	applicationID := bodyTerm.Operands[0].Term
	application := term(t, model, applicationID)
	callee := term(t, model, application.Callee)
	if application.Kind != semantics.Application || application.Target != 0 || callee.Kind != semantics.Reference || callee.Target != f {
		t.Fatal("parameter-headed application lost its callee binder")
	}
	if _, ok := model.DirectParameterBindings(applicationID); ok {
		t.Fatal("constructor interface was mistaken for direct declaration formals")
	}
	probe := summaryOperator(t, model, oneFact(t, model, semantics.Subject{Entity: wrappedID}, resolve.DataSummaryRelation()), "data-summary/probe")
	if len(probe.Operands) != 1 {
		t.Fatal("abstract application was collapsed into a failed Data certificate")
	}
	retained := term(t, model, probe.Operands[0].Term)
	if retained.Kind != semantics.Application || term(t, model, retained.Callee).Target != f || len(retained.Arguments) != 1 {
		t.Fatal("summary probe lost typed constructor application structure")
	}
	if len(model.Facts(semantics.Subject{Term: probe.Operands[0].Term}, resolve.DataSummaryRelation())) != 0 {
		t.Fatal("probe recursively exported its own summary")
	}
	identity := oneFact(t, model, semantics.Subject{Entity: f}, resolve.IdentityRelation())
	if identity.Operator.Name != "identity/parameter" {
		t.Fatal("parameter lost its canonical identity fact")
	}
	pageID, _ := entity(t, model, "example.kinds", "Page")
	appliedID, _ := entity(t, model, "example.kinds", "WrappedPage")
	found := false
	for _, dependency := range model.DependenciesFrom(appliedID, resolve.ReferenceDependency()) {
		found = found || dependency.To == pageID
	}
	if !found {
		t.Fatal("unapplied constructor reference was omitted from dependencies")
	}
	composeID, compose := entity(t, model, "example.kinds", "Compose")
	composeProbe := summaryOperator(t, model, oneFact(t, model, semantics.Subject{Entity: composeID}, resolve.DataSummaryRelation()), "data-summary/probe")
	outer := term(t, model, composeProbe.Operands[0].Term)
	inner := term(t, model, outer.Arguments[0])
	if term(t, model, outer.Callee).Target != compose.Parameters[0].Entity || inner.Kind != semantics.Application || term(t, model, inner.Callee).Target != compose.Parameters[1].Entity || term(t, model, inner.Arguments[0]).Target != compose.Parameters[2].Entity {
		t.Fatal("Compose probe lost nested F/G/T identities")
	}
	_, invoke := entity(t, model, "example.kinds", "Invoke")
	higher := oneFact(t, model, semantics.Subject{Entity: invoke.Parameters[0].Entity}, resolve.InterfaceRelation())
	shape := oneFact(t, model, semantics.Subject{Entity: invoke.Parameters[0].Entity}, resolve.KindRelation())
	if len(higher.Operands) != 2 || term(t, model, higher.Operands[0].Term).Operator.Name != "interface/constructor" || len(shape.Operands) != 2 || term(t, model, shape.Operands[0].Term).Operator.Name != "kind/constructor" {
		t.Fatal("higher-order interface or erased kind inputs were flattened")
	}
}

func TestImportedSummarySubstitutionRetainsCallerBindersAndDefinitionScope(t *testing.T) {
	proper := properClaim("type")
	constructor := constructorClaim("type", proper)
	library := module("example.library", mapping(),
		entry("Wrap", kindDefinition(applyExpression(param(0), param(1)), kindParameter("F", constructor), kindParameter("T", proper))),
		entry("Page", kindDefinition(h("list-of", param(0)), kindParameter("T", proper))))
	libraryPin := pin(t, library)
	root := module("example.caller", mapping(entry("lib", d("import", txt("example.library"), core.NewAtom(libraryPin[:])))),
		entry("Page", kindDefinition(h("wire", param(0)), kindParameter("T", proper))),
		entry("Forward", kindDefinition(imported("lib", "Wrap", param(0), param(1)), kindParameter("F", constructor), kindParameter("T", proper))),
		entry("Known", kindDefinition(imported("lib", "Wrap", imported("lib", "Page"), param(0)), kindParameter("T", proper))),
		entry("Applied", definition(local("Known", h("scalar", atom("utf8-text"))))))
	tree, err := resolve.Resolve(root, []value{library}, "imported-kinds")
	if err != nil {
		t.Fatal(err)
	}
	model := requireModel(t, run(t, tree, nil))
	forwardID, forward := entity(t, model, "example.caller", "Forward")
	wrapID, wrap := entity(t, model, "example.library", "Wrap")
	probe := summaryOperator(t, model, oneFact(t, model, semantics.Subject{Entity: forwardID}, resolve.DataSummaryRelation()), "data-summary/probe")
	application := term(t, model, probe.Operands[0].Term)
	if term(t, model, application.Callee).Target != wrapID || len(application.Arguments) != 2 || term(t, model, application.Arguments[0]).Target != forward.Parameters[0].Entity || term(t, model, application.Arguments[1]).Target != forward.Parameters[1].Entity || forward.Parameters[0].Entity == wrap.Parameters[0].Entity {
		t.Fatal("forwarded summary captured the imported declaration's same-spelled binders")
	}
	knownID, known := entity(t, model, "example.caller", "Known")
	libraryPageID, libraryPage := entity(t, model, "example.library", "Page")
	knownProbe := summaryOperator(t, model, oneFact(t, model, semantics.Subject{Entity: knownID}, resolve.DataSummaryRelation()), "data-summary/probe")
	knownApplication := term(t, model, knownProbe.Operands[0].Term)
	if term(t, model, knownApplication.Callee).Target != wrapID || len(knownApplication.Arguments) != 2 || term(t, model, knownApplication.Arguments[0]).Target != libraryPageID || term(t, model, knownApplication.Arguments[1]).Target != known.Parameters[0].Entity {
		t.Fatal("known constructor probe lost its definition target or caller arguments")
	}
	need := summaryOperator(t, model, oneFact(t, model, semantics.Subject{Entity: libraryPageID}, resolve.DataSummaryRelation()), "data-summary/need")
	if len(need.Operands) != 1 || term(t, model, need.Operands[0].Term).Target != libraryPage.Parameters[0].Entity || libraryPage.Parameters[0].Entity == known.Parameters[0].Entity {
		t.Fatal("evaluating a caller rewrote the imported Page summary's own T binder")
	}
	rootPin := pin(t, root)
	identity := oneFact(t, model, semantics.Subject{Entity: known.Parameters[0].Entity}, resolve.IdentityRelation())
	expected := core.NewTuple(core.NewAtom(rootPin[:]), txt("Known"), integer(0))
	if !bytes.Equal(identity.Data, codec.Encode(expected)) {
		t.Fatal("same-spelled parameter lost exact module/declaration/slot identity")
	}
	appliedID, _ := entity(t, model, "example.caller", "Applied")
	if oneFact(t, model, semantics.Subject{Entity: appliedID}, resolve.DataStatusRelation()).Operator.Name != "data-status/proven" {
		t.Fatal("imported Page was rebound to the caller's Wire-bearing Page")
	}
}

func TestParameterCalleeFailurePreservesArgumentAndBinderSites(t *testing.T) {
	for _, test := range []struct {
		name, source, code string
		child int
	}{
		{"admission", "module example.bad; type Bad<F: Data -> Data> = F<wire<text>>;", "bound-violation", 1},
		{"shape", "module example.bad; type Bad<F> = F<text>;", "kind-mismatch", 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			parsed, err := parse.Parse(parsecontract.Input{SourceID: "bad.horos", TreeID: "bad", Text: test.source})
			if err != nil {
				t.Fatal(err)
			}
			tree := parsed.ASTView()
			result := run(t, tree, nil)
			requireFailure(t, result, semantic.Complete, semantic.Rejected, test.code, "source")
			bad := body(tree, find(tree, s.KindDeclaration, "Bad"))
			location := tree.Nodes[bad-1].Children[test.child].Node
			if result.Diagnostics()[0].At == nil || result.Diagnostics()[0].At.Node != location {
				t.Fatal("primary rejection lost its callee/argument source path")
			}
			binder := find(tree, s.KindTypeParameter, "F")
			found := false
			for _, diagnostic := range result.Diagnostics()[1:] {
				found = found || diagnostic.Code.Name == "callee-declaration" && diagnostic.At != nil && diagnostic.At.Node == binder
			}
			if !found {
				t.Fatal("constructor-parameter rejection lost its declared binder site")
			}
		})
	}
}

type missingTypeFacts struct{ s.Language }

func (missingTypeFacts) Resolve(_ context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	var reports []semantic.Report
	for _, obligation := range resolve.Contract().Semantic.Obligations {
		if obligation.Key == resolve.TypeFacts() {
			continue
		}
		report := semantic.Report{Obligation: obligation.Key, Status: semantic.Complete}
		if obligation.Key == resolve.Validity() || obligation.Key == resolve.ReferenceResolution() {
			report.Verdict = semantic.Holds
		}
		reports = append(reports, report)
	}
	return request.Publish(semantic.Publication{Reports: reports, Facts: &semantics.Data{}})
}

func TestTypeFactCoverageIsRequiredBeforeModelPublication(t *testing.T) {
	tree := fixture(t)
	language, input := horosInput(t, tree)
	checker, err := check.NewChecker[s.Language](language, provider.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := semantic.Capture(context.Background(), checker, semantic.Input{Entry: ast.NodeRef{Tree: tree.ID, Node: tree.Root}, Trees: input.Trees})
	if err != nil {
		t.Fatal(err)
	}
	resolver, err := resolution.NewResolver[s.Language](language, resolve.Contract(), missingTypeFacts{})
	if err != nil {
		t.Fatal(err)
	}
	result, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{})
	if err == nil || result.Valid() {
		t.Fatal("missing type-facts enumeration established a model")
	}
}

func TestRevisionTwoResolutionProfileIsRejected(t *testing.T) {
	tree := fixture(t)
	language, input := horosInput(t, tree)
	checker, err := check.NewChecker[s.Language](language, provider.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := semantic.Capture(context.Background(), checker, semantic.Input{Entry: input.Entry, Trees: input.Trees})
	if err != nil {
		t.Fatal(err)
	}
	contract := resolve.Contract()
	contract.Semantic.Profile = "horos-declarations-2/resolution-1"
	resolver, err := resolution.NewResolver[s.Language](language, contract, provider.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	if result, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{}); err == nil || result.Valid() {
		t.Fatal("superseded resolution profile was accepted")
	}
}
