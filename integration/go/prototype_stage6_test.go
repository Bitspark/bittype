package integration_test

// PROTOTYPE (bittype#196), stage 6: optional fields by delegation (row 11,
// D3). Rows 12, 13, 14 and 35 of gate lines 2 and 4. The evidence for row 11.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	tsdeclarationsfacets "github.com/Bitspark/bittype/completeness/horos-to-typescript/facets/declarations/go"
	tslistfacets "github.com/Bitspark/bittype/completeness/horos-to-typescript/facets/list/go"
	tsoptionalfacets "github.com/Bitspark/bittype/completeness/horos-to-typescript/facets/optional/go"
	tsrecordfacets "github.com/Bitspark/bittype/completeness/horos-to-typescript/facets/record/go"
	tstextfacets "github.com/Bitspark/bittype/completeness/horos-to-typescript/facets/text/go"
	tsvalues "github.com/Bitspark/bittype/completeness/typescript/values/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	ggtarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

var standinConversion = protoRequest{products: []string{tsvalues.Product}, configuration: []g.Setting{{Name: "exactOptionalPropertyTypes", Value: "true"}}}

// standinBundle is the TypeScript bundle with its facet packages, which name
// the stand-in's factories, and the stand-in product.
func standinBundle(t *testing.T, from representation.Language) bundle[ts.Language] {
	t.Helper()
	b := typescriptBundle(t, from)
	for _, facets := range []func() ([]g.Facet, error){tstextfacets.Facets, tslistfacets.Facets, tsrecordfacets.Facets, tsdeclarationsfacets.Facets, tsoptionalfacets.Facets} {
		b.facets = append(b.facets, protoMust[[]g.Facet](t)(facets())...)
	}
	b.products = []g.Product{tsvalues.Conversion()}
	return b
}

// A vector of the stand-in. With the direction "encode" it starts from a
// native value; {"$undefined": true} in it stands for a present undefined.
type standinVector struct {
	Root      string `json:"root"`
	Direction string `json:"direction,omitempty"`
	Value     any    `json:"value,omitempty"`
	Native    any    `json:"native,omitempty"`
}
type standinFailure struct {
	Kind    string `json:"kind"`
	Message string `json:"message"`
}
type standinResult struct {
	Root   string          `json:"root"`
	Native any             `json:"native"`
	Value  any             `json:"value"`
	Error  *standinFailure `json:"error"`
}
type standinOutput struct {
	Results []standinResult   `json:"results"`
	Memo    map[string]string `json:"memo"`
	Failure *standinFailure   `json:"failure"`
}

// runStandin runs the TypeScript stand-in with node, in another process. It
// is given the serialized plan, the snapshot and the vectors as files.
func runStandin(t *testing.T, plan, snapshot []byte, vectors []standinVector, env ...string) standinOutput {
	t.Helper()
	dir, err := os.MkdirTemp(os.Getenv("GOTMPDIR"), "bittype-prototype-standin-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	write := func(name string, content []byte) string {
		file := filepath.Join(dir, name)
		if err := os.WriteFile(file, content, 0o644); err != nil {
			t.Fatal(err)
		}
		return file
	}
	encoded, err := json.Marshal(vectors)
	if err != nil {
		t.Fatal(err)
	}
	if vectors == nil {
		encoded = []byte("[]")
	}
	reader, err := filepath.Abs(filepath.Join("..", "..", "completeness", "typescript", "standin", "reader.mjs"))
	if err != nil {
		t.Fatal(err)
	}
	args := []string{reader, "--plan", write("plan.json", plan), "--vectors", write("vectors.json", encoded)}
	if snapshot != nil {
		args = append(args, "--snapshot", write("snapshot.json", snapshot))
	}
	node := os.Getenv("BITTYPE_TYPESCRIPT_NODE")
	if node == "" {
		node = "node"
	}
	command := exec.Command(node, args...)
	command.Env = append(os.Environ(), env...)
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err != nil {
		t.Fatalf("the stand-in did not run: %v\n%s", err, stderr.String())
	}
	var out standinOutput
	if err := json.Unmarshal(stdout.Bytes(), &out); err != nil {
		t.Fatal(err, stdout.String())
	}
	return out
}

const protoOptionalModule = "module example.optional;\n\ntype Person = { name: text; email?: text; };\n"

// Rows 13 and 35, O14: "Each state of an optional field has its own vector,
// in Go and in the stand-in, and a present `undefined` is refused."
//
// Shows: the record delegates presence through a field-state protocol on both
// sides, and never reads the optional contribution's payload. Absent and
// present with a value convert there and back in the lowered Go and in the
// stand-in. Present with null is refused on both sides, by the operand's
// adapter, because no delivered family admits null. In the stand-in a present
// undefined is refused by the protocol; red when the protocol reads it as
// absence, which is what the service library's record adapter does today.
// Does not show: a null that a type admits (#211). It does not show a present
// undefined in Go: a Go member is a value of its type, so that state does not
// exist there. It does not show the service library's adapter interface: the
// slot kind is theirs (row 11), and one seat wrote both sides here.
func TestPrototypeStage6FourStatesOfAnOptionalField(t *testing.T) {
	if testing.Short() {
		t.Skip("compiles and runs Go, and runs node")
	}
	from, input := protoSource(t, protoOptionalModule)
	snapshot := protoSnapshot(t, protoOptionalModule)
	absent, present := map[string]any{"name": "Ada"}, map[string]any{"name": "Ada", "email": "ada@example.org"}
	null := map[string]any{"name": "Ada", "email": nil}
	// Go: the lowered conversion.
	dir, err := conversionProgram(t, protoOptionalModule, conversionBundle(t, from), "", snapshot)
	if err != nil {
		t.Fatal(err)
	}
	results, err := runConversion(t, dir, []conversionVector{{Root: "Person", Value: absent}, {Root: "Person", Value: present}, {Root: "Person", Value: null}})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Error != "" || !reflect.DeepEqual(results[0].Value, absent) || results[1].Error != "" || !reflect.DeepEqual(results[1].Value, present) {
		t.Errorf("Go: absent and present came back as %v and %v", results[0], results[1])
	} else {
		t.Log("GREEN Go: an absent field and a present field with a value convert there and back, each as itself")
	}
	if !strings.Contains(results[2].Error, "text is expected") {
		t.Errorf("Go: a present null came back as %v", results[2])
	} else {
		t.Logf("GREEN Go: a present null is refused: %s", results[2].Error)
	}
	// The stand-in.
	plan, err := protoPlanTS(t, input, standinBundle(t, from), standinConversion)
	if err != nil {
		t.Fatal(err)
	}
	serialized, err := plan.Plan.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	vectors := []standinVector{
		{Root: "Person", Value: absent},
		{Root: "Person", Value: present},
		{Root: "Person", Direction: "encode", Native: absent},
		{Root: "Person", Direction: "encode", Native: present},
		{Root: "Person", Direction: "encode", Native: null},
		{Root: "Person", Direction: "encode", Native: map[string]any{"name": "Ada", "email": map[string]any{"$undefined": true}}},
	}
	check := func(env ...string) error {
		out := runStandin(t, serialized, snapshot, vectors, env...)
		if out.Failure != nil {
			return fmt.Errorf("the stand-in failed: %+v", out.Failure)
		}
		r := out.Results
		if r[0].Error != nil || !reflect.DeepEqual(r[0].Value, absent) || !reflect.DeepEqual(r[0].Native, absent) {
			return fmt.Errorf("absent came back as %+v", r[0])
		}
		if r[1].Error != nil || !reflect.DeepEqual(r[1].Value, present) {
			return fmt.Errorf("present came back as %+v", r[1])
		}
		if r[2].Error != nil || !reflect.DeepEqual(r[2].Value, absent) || r[3].Error != nil || !reflect.DeepEqual(r[3].Value, present) {
			return fmt.Errorf("native absent and present encode as %+v and %+v", r[2], r[3])
		}
		if r[4].Error == nil || r[4].Error.Kind != "invalid native value" {
			return fmt.Errorf("a present null encodes as %+v", r[4])
		}
		if r[5].Error == nil || !strings.Contains(r[5].Error.Message, "present and undefined") {
			return fmt.Errorf("a present undefined was not refused: it encodes as %v", r[5].Value)
		}
		return nil
	}
	if err := check(); err != nil {
		t.Error(err)
	} else {
		t.Log("GREEN the stand-in: absent, present with a value, present with null (refused) and present with undefined (refused) each have their own vector and outcome")
	}
	if err := check("BITTYPE_STANDIN_MUTANT=undefined-as-absent"); err == nil {
		t.Error("a protocol that reads a present undefined as absence is not detected")
	} else {
		t.Logf("RED   under a protocol that reads a present undefined as absence: %v", err)
	}
}

// smuggling wraps the Go record contribution as a fixture that keeps the
// operand and the role selector of the first optional field it meets, and
// presents them at a later record: the operand of another record, or its own
// operand with the role issued for another pair.
func smuggling(base golangContribution, mode string) golangContribution {
	base.Reference.Identity = g.ContributionIdentity{Publisher: fixturePublisher, Key: "smuggling-record"}
	var kept *struct {
		use  g.Use[h.Language]
		role g.Role[h.Language]
	}
	base.Rule.Name = func(_ context.Context, step horosName) ([]horosRequest, error) {
		fields, err := step.Reader.Fields(step.Subject.Use)
		var requests []horosRequest
		for _, f := range fields {
			switch {
			case !f.Role.Valid():
				requests = append(requests, g.Reference(f.Use))
			case kept == nil:
				kept = &struct {
					use  g.Use[h.Language]
					role g.Role[h.Language]
				}{f.Use, f.Role}
				requests = append(requests, g.Field(f.Use, f.Role, f.Name))
			case mode == "operand":
				requests = append(requests, g.Field(kept.use, kept.role, f.Name))
			default:
				requests = append(requests, g.Field(f.Use, kept.role, f.Name))
			}
		}
		return requests, err
	}
	return base
}

// Row 12: "A request at a field position with the operand of another record
// fails".
//
// Shows: the kernel keeps its scope check at a field position. The operand
// must be a child of the asking record's use, and the role selector must have
// been issued by the source for exactly that parent and operand. A rule
// cannot present an arbitrary operand or an arbitrary role. Red under a
// kernel without the check.
// Does not show: that the source issues a role only for an optional field.
// That is the source integration's own fact.
func TestPrototypeStage6FieldPositionKeepsTheScopeCheck(t *testing.T) {
	from, input := protoSource(t, "module example.scope;\n\ntype A = { x: { p?: text; }; y: { q?: text; }; };\n")
	b := golangBundle(t, from)
	for mode, want := range map[string]string{"operand": "a capture must be a child", "role": "the role selector was not issued for this parent and operand"} {
		redGreen(t, "a request at a field position with the "+mode+" of another record fails", func() error {
			_, err := protoPlanGo(t, input, b.replace(smuggling(b.find(t, "horos-to-golang/record"), mode)), protoRequest{})
			if !errors.Is(err, g.ErrScope) || !strings.Contains(err.Error(), want) {
				return fmt.Errorf("not rejected: %v", err)
			}
			return nil
		}, "field-ignores-parent")
	}
}

// withoutMembers is the Go target without its member vocabulary.
type withoutMembers struct{ g.Target[gs.Language] }

// Row 14: "A member fragment that the target's member vocabulary does not
// admit is refused, with the class that the prototype fixes".
//
// Shows: the answer of a contribution at a field position is checked against
// the member vocabulary that the target states. A fragment that is no member
// is refused as unsupported, with the member vocabulary as the owner of the
// reason and the contribution named. A target that states no member
// vocabulary admits no fragment. A fragment of two members is admitted by the
// Go model, which checks each member alone.
// Does not show: a layout shared across fields. That needs the record policy
// replaced (row 14) and is out of the extension domain.
func TestPrototypeStage6MemberFragmentOutsideTheVocabulary(t *testing.T) {
	from, input := protoSource(t, protoOptionalModule)
	b := golangBundle(t, from)
	fragment := func(key string, members func(p g.Projection[gs.Language]) []g.Expression[gs.Language]) golangContribution {
		c := b.find(t, "horos-to-golang/optional")
		c.Reference.Identity = g.ContributionIdentity{Publisher: fixturePublisher, Key: key}
		c.Project = func(_ context.Context, p g.Projection[gs.Language]) (g.Projected[gs.Language], error) {
			built := members(p)
			table := make([]g.Member, len(built))
			for i := range built {
				table[i] = g.Member{Name: fmt.Sprintf("%s%d", p.Stated(), i)}
			}
			return g.Projected[gs.Language]{Fragment: built, Members: table}, nil
		}
		return c
	}
	noMember := fragment("no-member", func(p g.Projection[gs.Language]) []g.Expression[gs.Language] {
		return []g.Expression[gs.Language]{model.Collection(p.Dependency(0).Type)}
	})
	_, err := protoPlanGo(t, input, b.replace(noMember), protoRequest{})
	refusal := expectRefusal(t, "a fragment that is no member of the Go model", err, "unsupported")
	if refusal.Owner != "member vocabulary of "+golangTarget || !strings.HasPrefix(refusal.Reason, g.ReasonMemberFragment) || refusal.Contribution != noMember.Reference {
		t.Errorf("the refusal is owned by %q with the reason %q for %s", refusal.Owner, refusal.Reason, refusal.Contribution)
	}
	_, err = protoPlanWith[gs.Language](t, input, nativeChecker(t, golangText()), withoutMembers{ggtarget.Provider{Package: "generated"}}, b, protoRequest{})
	if refusal := expectRefusal(t, "a field position under a target without a member vocabulary", err, "unsupported"); refusal.Reason != g.ReasonNoMemberVocabulary {
		t.Errorf("the reason is %q", refusal.Reason)
	}
	twoMembers := fragment("two-members", func(p g.Projection[gs.Language]) []g.Expression[gs.Language] {
		return []g.Expression[gs.Language]{model.Field(p.Stated()+"Present", false, model.Primitive(model.PrimitiveBool)), model.Field(p.Stated(), false, p.Dependency(0).Type)}
	})
	plan, err := protoPlanGo(t, input, b.replace(twoMembers), protoRequest{})
	if err != nil {
		t.Fatalf("a fragment of two members is not admitted: %v", err)
	}
	fields, _, _ := goFieldTypes(t, protoText[gs.Language](t, golangText(), plan))
	if fields["Person"]["EmailPresent"] != "bool" || fields["Person"]["Email"] != "string" || len(bodyNode(t, plan.Plan, "Person").Members) != 3 {
		t.Errorf("the record with a fragment of two members is %v with the members %v", fields["Person"], bodyNode(t, plan.Plan, "Person").Members)
	} else {
		t.Logf("GREEN a fragment of two members is admitted by the Go model: %v", fields["Person"])
	}
}
