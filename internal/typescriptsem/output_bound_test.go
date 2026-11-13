package typescriptsem

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/internal/processoutput/processoutputtest"
	parser "github.com/Bitspark/bittype/parse/syntax/typescript/provider/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
)

// The child process of the output-bound tests is this test binary, and what
// it does is the shared fixture's. The bridge starts the executable it is
// given as node, with an environment that holds no variable of the test: it
// evaluates the helper source and passes the compiler module as the last
// argument. So a test names a description file of the fixture as the module,
// and a process that is started with node's arguments and such a file is the
// child and runs no test.
func TestMain(m *testing.M) {
	if len(os.Args) == 5 && os.Args[1] == "--input-type=commonjs" && os.Args[2] == "--eval" && filepath.Base(os.Args[4]) == processoutputtest.File {
		var child processoutputtest.Child
		processoutputtest.Described(os.Args[4], &child)
		child.Run()
	}
	os.Exit(m.Run())
}

// The bridge bounds what its helper may write: 16 MiB on the output stream and
// 16 KiB on the error stream. Every case starts a real process that writes to
// its streams, and the bridge receives them the way it receives its helper's.
// A stream of exactly the bound is read whole. One byte more is
// semantic.ErrLimit and no output, whether the process then exits or stays.
func TestHelperOutputBound(t *testing.T) {
	const outputBound, errorBound = 16 << 20, 16 << 10
	const response = `{"Version":1,"Status":"incomplete"}`
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	text := "type Text = string;"
	parsed, err := parser.Parse(context.Background(), parse.Input{Text: text, SourceID: "text", TreeID: "text"})
	if err != nil {
		t.Fatal(err)
	}
	trees := []s.AST{parsed.ASTView()}
	source := &files.Context{Files: files.Files{Entries: []files.File{{ID: "text", Path: "text.ts", Text: text}}}, Links: []files.SyntaxLink{{File: "text", Tree: "text", Source: "text"}}}
	check := func(path string) (Output, error) {
		return Check(context.Background(), self, path, trees, source, nil, n.Configuration())
	}
	refused := func(t *testing.T, out Output, err error) {
		t.Helper()
		if !errors.Is(err, semantic.ErrLimit) || !reflect.DeepEqual(out, Output{}) {
			t.Errorf("the bridge returned status %q, %d diagnostics and the error %v", out.Status, len(out.Diagnostics), err)
		}
	}
	for _, c := range []struct {
		name           string
		stdout, stderr int
		over           bool
	}{
		{"the output stream at its bound", outputBound, 0, false},
		{"the output stream one byte over its bound", outputBound + 1, 0, true},
		{"the error stream at its bound", len(response), errorBound, false},
		{"the error stream one byte over its bound", len(response), errorBound + 1, true},
	} {
		child := processoutputtest.Child{Input: true, Response: response, Stdout: c.stdout, Stderr: c.stderr}
		if !c.over {
			t.Run(c.name, func(t *testing.T) {
				if out, err := check(processoutputtest.Describe(t, child)); err != nil || out.Status != statusIncomplete {
					t.Errorf("the bridge returned status %q and the error %v", out.Status, err)
				}
				processoutputtest.AwaitNoProcessGoroutine(t)
			})
			continue
		}
		t.Run(c.name+", and the process exits", func(t *testing.T) {
			out, err := check(processoutputtest.Describe(t, child))
			refused(t, out, err)
			processoutputtest.AwaitNoProcessGoroutine(t)
		})
		t.Run(c.name+", and the process stays", func(t *testing.T) {
			// The bridge's own deadline is 30 seconds, and the harness waits half
			// of that. A bridge that has not returned by then did not cancel its
			// helper at the bound.
			hold := processoutputtest.Listen(t)
			child.Hold = hold.Address()
			path := processoutputtest.Describe(t, child)
			type result struct {
				out Output
				err error
			}
			r := processoutputtest.AwaitCancelled(t, hold, func() result {
				out, err := check(path)
				return result{out, err}
			})
			refused(t, r.out, r.err)
			processoutputtest.AwaitNoProcessGoroutine(t)
		})
	}
}
