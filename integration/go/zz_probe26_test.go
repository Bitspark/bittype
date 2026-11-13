package integration_test

import (
	"fmt"
	"testing"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	n "github.com/Bitspark/bittype/semantics/rust/go"
)

// Throwaway probe for row 26 of research 0016. Not committed.
func TestProbeRow26AbsoluteStdPath(t *testing.T) {
	for _, text := range []string{
		"pub type Text = ::std::string::String;",
		"pub type Text = std::string::String;",
		"pub type Text = String;",
		"pub struct Person { pub name: ::std::string::String, pub tags: ::std::vec::Vec<::std::string::String>, pub nick: ::std::option::Option<::std::string::String>, pub ok: bool }",
	} {
		fmt.Println("### source:", text)
		snapshot, trees, _ := rsCapture(t, text)
		_, result := rsResults(t, snapshot)
		for _, key := range []string{"captured-type-coverage", "bindings", "captured-type-validity"} {
			r, _ := result.Report(n.Key(key))
			fmt.Printf("  report %s: status=%v verdict=%v\n", key, r.Status, r.Verdict)
		}
		for _, d := range result.Diagnostics() {
			fmt.Printf("  diagnostic: %s %s %v\n", d.Code.Name, d.Message, d.At)
		}
		if !result.HasSemantics() {
			fmt.Println("  no semantics published")
			continue
		}
		m, err := n.Open(result)
		if err != nil {
			t.Fatal(err)
		}
		data := m.Graph().Data()
		for _, b := range data.Bindings {
			node := trees[0].Nodes[b.At.Site.Node-1]
			name := ""
			for _, a := range node.Attributes {
				if a.Name == "name" {
					name = a.Text
				}
			}
			e := data.Entities[b.Target-1]
			fmt.Printf("  binding: node %d (%s %q) role=%s -> entity %d kind=%s name=%s declarations=%d\n", b.At.Site.Node, node.Kind, name, b.At.Role.Name, b.Target, e.Kind.Name, e.Name, len(e.Declarations))
		}
		_ = semantic.Holds
	}
}
