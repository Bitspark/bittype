package rustsem

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/internal/processoutput/processoutputtest"
	n "github.com/Bitspark/bittype/semantics/rust/go"
)

// The child process of the output-bound tests is this test binary, and what
// it does is the shared fixture's. The bridge starts its helper with an empty
// environment and one argument, the sysroot. So a test names a description
// file of the fixture as the sysroot, and a process that is started with such
// a file as its only argument is the child and runs no test.
func TestMain(m *testing.M) {
	if len(os.Args) == 2 && filepath.Base(os.Args[1]) == processoutputtest.File {
		var child processoutputtest.Child
		processoutputtest.Described(os.Args[1], &child)
		child.Run()
	}
	os.Exit(m.Run())
}

// The bridge bounds what its helper may write: 16 MiB on the output stream and
// 256 KiB on the error stream. Every case starts a real process that writes to
// its streams, and the bridge receives them the way it receives its helper's.
// A stream of exactly the bound is read whole. One byte more is
// semantic.ErrLimit and no output, whether the process then exits or stays.
func TestHelperOutputBound(t *testing.T) {
	const outputBound, errorBound = 16 << 20, 256 << 10
	response := `{"Version":1,"Profile":"` + n.HelperProfile + `","Compiler":"rustc ` + n.CompilerVersion + ` (test)","Status":"incomplete"}`
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	trees, source := compilerFixture(t)
	check := func(path string) (Output, error) {
		return Check(context.Background(), self, path, environment, trees, source)
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
