package report

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
)

// The loader starts the command named go that the search path gives, with the
// environment of its own process. So a test that needs to know what the
// command writes puts a copy of this test binary first in the search path
// under that name, and names in this variable a file that says what the copy
// does. A process that is started with the variable and with one of the two
// subcommands the loader runs is the child and runs no test.
const goChildVariable = "BITTYPE_REPORT_GO_CHILD"

// goChild is what the child does for one subcommand. It writes Stdout bytes
// to its output stream and Stderr bytes to its error stream and exits with
// status 0. For env the output stream starts with a build context that the
// loader accepts, and the rest is spaces. For list it is spaces only: no
// package. With Hold it connects to that address first and does not exit by
// itself: it ends when the loader kills it or when the test closes the
// connection.
type goChild struct {
	Stdout int
	Stderr int
	Hold   string
}

func TestMain(m *testing.M) {
	if path := os.Getenv(goChildVariable); path != "" && len(os.Args) > 1 && (os.Args[1] == "env" || os.Args[1] == "list") {
		runGoChild(path, os.Args[1])
	}
	os.Exit(m.Run())
}

// goChildContext is the build context the child reports for go env.
func goChildContext() string {
	return fmt.Sprintf(`{"GOOS":%q,"GOARCH":%q,"GOFLAGS":"","GOVERSION":%q}`, runtime.GOOS, runtime.GOARCH, runtime.Version())
}

func runGoChild(path, subcommand string) {
	var children map[string]goChild
	data, err := os.ReadFile(path)
	if err != nil || json.Unmarshal(data, &children) != nil {
		os.Exit(2)
	}
	child := children[subcommand]
	response := ""
	if subcommand == "env" {
		response = goChildContext()
	}
	if child.Stdout < len(response) {
		os.Exit(2)
	}
	var hold net.Conn
	if child.Hold != "" {
		// A stream that the loader has closed must not end the child either.
		signal.Ignore(syscall.SIGPIPE)
		if hold, err = net.Dial("tcp", child.Hold); err != nil {
			os.Exit(2)
		}
	}
	os.Stdout.Write(append([]byte(response), bytes.Repeat([]byte{' '}, child.Stdout-len(response))...))
	os.Stderr.Write(bytes.Repeat([]byte{' '}, child.Stderr))
	if hold != nil {
		io.Copy(io.Discard, hold)
		os.Exit(3)
	}
	os.Exit(0)
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
func useGoChild(t *testing.T, directory string, children map[string]goChild) {
	t.Helper()
	data, err := json.Marshal(children)
	if err != nil {
		t.Fatal(err)
	}
	description := filepath.Join(t.TempDir(), "children.json")
	if err := os.WriteFile(description, data, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", directory)
	t.Setenv(goChildVariable, description)
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
	within := map[string]goChild{
		"env":  {Stdout: first + 1000, Stderr: 300 << 10},
		"list": {Stdout: 400 << 10, Stderr: budget - first - 1000 - 300<<10 - 400<<10},
	}
	command := installGoChild(t)
	run := func() (*program, error) {
		return load(context.Background(), t.TempDir(), LoadLimits{MaxMetadataBytes: budget}, nil)
	}
	t.Run("all four streams together at the budget", func(t *testing.T) {
		useGoChild(t, command, within)
		// The child lists no package, so a loader that read every stream whole
		// fails at its next step and says so.
		if p, err := run(); p != nil || err == nil || !strings.HasPrefix(err.Error(), "no main Go module at ") {
			t.Errorf("the loader returned the error %v", err)
		}
	})
	const buildContext, packageMetadata = "load Go build context: ", "load Go package metadata: "
	for _, c := range []struct {
		name string
		// env and list are added to the streams of the case at the budget.
		env, list goChild
		// passes names the command whose write passes the budget, and step the
		// step of the loader that runs it.
		passes, step string
	}{
		// The first command stays within the budget, so the byte counts when
		// the second command writes.
		{"one byte more in the output stream of the first command", goChild{Stdout: 1}, goChild{}, "list", packageMetadata},
		{"one byte more in the error stream of the first command", goChild{Stderr: 1}, goChild{}, "list", packageMetadata},
		{"one byte more in the output stream of the second command", goChild{}, goChild{Stdout: 1}, "list", packageMetadata},
		{"one byte more in the error stream of the second command", goChild{}, goChild{Stderr: 1}, "list", packageMetadata},
		// The first command passes the budget by itself.
		{"the output stream of the first command one byte over the budget", goChild{Stdout: budget + 1 - within["env"].Stdout - within["env"].Stderr}, goChild{}, "env", buildContext},
		{"the error stream of the first command one byte over the budget", goChild{Stderr: budget + 1 - within["env"].Stdout - within["env"].Stderr}, goChild{}, "env", buildContext},
	} {
		over := map[string]goChild{
			"env":  {Stdout: within["env"].Stdout + c.env.Stdout, Stderr: within["env"].Stderr + c.env.Stderr},
			"list": {Stdout: within["list"].Stdout + c.list.Stdout, Stderr: within["list"].Stderr + c.list.Stderr},
		}
		refused := func(t *testing.T, p *program, err error) {
			t.Helper()
			if p != nil || err == nil || err.Error() != c.step+exceeded {
				t.Errorf("the loader returned the error %v", err)
			}
		}
		t.Run(c.name+", and the process exits", func(t *testing.T) {
			useGoChild(t, command, over)
			p, err := run()
			refused(t, p, err)
		})
		t.Run(c.name+", and the process stays", func(t *testing.T) {
			// The loader has no deadline of its own here. A loader that has not
			// returned after this time did not cancel its command at the budget.
			const wait = 15 * time.Second
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()
			staying := map[string]goChild{"env": over["env"], "list": over["list"]}
			child := staying[c.passes]
			child.Hold = listener.Addr().String()
			staying[c.passes] = child
			useGoChild(t, command, staying)
			connected := make(chan net.Conn, 1)
			go func() {
				connection, _ := listener.Accept()
				connected <- connection
			}()
			type result struct {
				p   *program
				err error
			}
			returned := make(chan result, 1)
			go func() {
				p, err := run()
				returned <- result{p, err}
			}()
			var connection net.Conn
			select {
			case connection = <-connected:
			case r := <-returned:
				t.Fatal("the loader returned before its command had started", r.err)
			}
			if connection == nil {
				t.Fatal("the command did not connect")
			}
			defer connection.Close()
			var r result
			select {
			case r = <-returned:
			case <-time.After(wait):
				t.Errorf("the loader had not returned %v after its command started: the command was not cancelled at the budget", wait)
				connection.Close()
				r = <-returned
			}
			refused(t, r.p, r.err)
			// The child never closes this connection and never exits by itself,
			// so the end of the connection is the end of the process.
			connection.SetReadDeadline(time.Now().Add(wait))
			if _, err := io.Copy(io.Discard, connection); errors.Is(err, os.ErrDeadlineExceeded) {
				t.Error("the command outlived the call of the loader")
			}
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
