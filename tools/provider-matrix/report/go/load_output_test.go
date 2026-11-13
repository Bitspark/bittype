package report

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/internal/processoutput/processoutputtest"
)

// The loader starts the command named go that the search path gives, with the
// environment of its own process. So a test that needs to know what the
// command writes puts a copy of this test binary first in the search path
// under that name, and names in this variable a description file of the
// shared fixture. The file holds one child for each of the two subcommands
// the loader runs. A process that is started with the variable and with one
// of them is that child and runs no test.
const goChildVariable = "BITTYPE_REPORT_GO_CHILD"

func TestMain(m *testing.M) {
	if path := os.Getenv(goChildVariable); path != "" && len(os.Args) > 1 && (os.Args[1] == "env" || os.Args[1] == "list") {
		var children map[string]processoutputtest.Child
		processoutputtest.Described(path, &children)
		children[os.Args[1]].Run()
	}
	os.Exit(m.Run())
}

// goChildContext is a build context that the loader accepts: what the child
// for go env starts its output stream with. The child for go list writes
// spaces only, which is a list of no package.
func goChildContext() string {
	return fmt.Sprintf(`{"GOOS":%q,"GOARCH":%q,"GOFLAGS":"","GOVERSION":%q}`, runtime.GOOS, runtime.GOARCH, runtime.Version())
}

// installGoChild writes a copy of this test binary under the name go into a
// new directory and returns the directory.
func installGoChild(t *testing.T) string {
	t.Helper()
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	binary, err := os.ReadFile(self)
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	name := "go"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	if err := os.WriteFile(filepath.Join(directory, name), binary, 0700); err != nil {
		t.Fatal(err)
	}
	return directory
}

// useGoChild makes the copy in the directory the loader's go command for the
// rest of the test, and makes it do what the description says.
func useGoChild(t *testing.T, directory string, children map[string]processoutputtest.Child) {
	t.Helper()
	t.Setenv("PATH", directory)
	t.Setenv(goChildVariable, processoutputtest.Describe(t, children))
}

// The loader bounds what its Go commands may write: MaxMetadataBytes for the
// output streams and the error streams of both commands together. Every case
// starts real processes that write to their streams, and the loader receives
// them the way it receives the Go command's. Streams of exactly the budget
// are read whole, and the loader goes on to its next step. One byte more on
// any of the four streams makes the load fail with the budget's own message
// and no program, whether the process that passes the budget then exits or
// stays.
func TestLoadBoundsGoCommandOutput(t *testing.T) {
	const budget = 1 << 20
	exceeded := fmt.Sprintf("Go command output exceeds %d bytes", budget)
	first := len(goChildContext())
	// The four streams fill the budget together, and no stream and no command
	// comes near it alone. So a loader that bounds each stream or each command
	// by itself, or that leaves one stream out, accepts a case below that must
	// be refused.
	within := map[string]processoutputtest.Child{
		"env":  {Response: goChildContext(), Stdout: first + 1000, Stderr: 300 << 10},
		"list": {Stdout: 400 << 10, Stderr: budget - first - 1000 - 300<<10 - 400<<10},
	}
	command := installGoChild(t)
	type result struct {
		p   *program
		err error
	}
	run := func(root string) result {
		p, err := load(context.Background(), root, LoadLimits{MaxMetadataBytes: budget}, nil)
		return result{p, err}
	}
	t.Run("all four streams together at the budget", func(t *testing.T) {
		useGoChild(t, command, within)
		// The child lists no package, so a loader that read every stream whole
		// fails at its next step and says so.
		if r := run(t.TempDir()); r.p != nil || r.err == nil || !strings.HasPrefix(r.err.Error(), "no main Go module at ") {
			t.Errorf("the loader returned the error %v", r.err)
		}
	})
	const buildContext, packageMetadata = "load Go build context: ", "load Go package metadata: "
	// more is what a case adds to the streams of one command at the budget.
	type more struct{ stdout, stderr int }
	alone := budget + 1 - within["env"].Stdout - within["env"].Stderr
	for _, c := range []struct {
		name      string
		env, list more
		// passes names the command whose write passes the budget, and step the
		// step of the loader that runs it.
		passes, step string
	}{
		// The first command stays within the budget, so the byte counts when
		// the second command writes.
		{"one byte more in the output stream of the first command", more{stdout: 1}, more{}, "list", packageMetadata},
		{"one byte more in the error stream of the first command", more{stderr: 1}, more{}, "list", packageMetadata},
		{"one byte more in the output stream of the second command", more{}, more{stdout: 1}, "list", packageMetadata},
		{"one byte more in the error stream of the second command", more{}, more{stderr: 1}, "list", packageMetadata},
		// The first command passes the budget by itself.
		{"the output stream of the first command one byte over the budget", more{stdout: alone}, more{}, "env", buildContext},
		{"the error stream of the first command one byte over the budget", more{stderr: alone}, more{}, "env", buildContext},
	} {
		// over returns the children of the case. hold is where the command that
		// passes the budget stays, or empty for one that exits.
		over := func(hold string) map[string]processoutputtest.Child {
			children := map[string]processoutputtest.Child{}
			for subcommand, add := range map[string]more{"env": c.env, "list": c.list} {
				child := within[subcommand]
				child.Stdout += add.stdout
				child.Stderr += add.stderr
				if subcommand == c.passes {
					child.Hold = hold
				}
				children[subcommand] = child
			}
			return children
		}
		refused := func(t *testing.T, r result) {
			t.Helper()
			if r.p != nil || r.err == nil || r.err.Error() != c.step+exceeded {
				t.Errorf("the loader returned the error %v", r.err)
			}
		}
		t.Run(c.name+", and the process exits", func(t *testing.T) {
			useGoChild(t, command, over(""))
			refused(t, run(t.TempDir()))
		})
		t.Run(c.name+", and the process stays", func(t *testing.T) {
			// The loader has no deadline of its own here. A loader that has not
			// returned when the harness stops waiting did not cancel its command
			// at the budget.
			hold := processoutputtest.Listen(t)
			useGoChild(t, command, over(hold.Address()))
			root := t.TempDir()
			refused(t, processoutputtest.AwaitCancelled(t, hold, func() result { return run(root) }))
		})
	}
}

// The same bound with the Go command itself and no stand-in: its answer to
// go env is longer than a budget of sixteen bytes, so the load of a module
// that has nothing wrong fails with the budget's own message.
func TestLoadFailsWhenTheGoCommandPassesTheBudget(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.test/budget\n\ngo 1.27\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "budget.go"), []byte("package budget\n"), 0600); err != nil {
		t.Fatal(err)
	}
	const budget = 16
	p, err := load(context.Background(), root, LoadLimits{MaxMetadataBytes: budget}, nil)
	if want := fmt.Sprintf("load Go build context: Go command output exceeds %d bytes", budget); p != nil || err == nil || err.Error() != want {
		t.Fatalf("the loader returned the error %v", err)
	}
	// With the default budget the same module loads.
	if p, err := load(context.Background(), root, LoadLimits{}, nil); err != nil || p == nil || p.module != "example.test/budget" {
		t.Fatal("the module does not load within the default budget", err)
	}
}
