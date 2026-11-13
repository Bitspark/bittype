package typescriptsem

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	parser "github.com/Bitspark/bittype/parse/syntax/typescript/provider/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
)

// The child process of the output-bound tests is this test binary. The bridge
// starts the executable it is given as node, with an environment that holds no
// variable of the test: it evaluates the helper source and passes the compiler
// module as the last argument. So a test names a file of this name as the
// module, and a process that is started with node's arguments and such a file
// is the child and runs no test.
const outputChildFile = "bittype-bridge-output-child.json"

// outputChild is what the child does. It reads its input stream to the end,
// writes Response and then spaces to its output stream until it has written
// Stdout bytes, writes Stderr bytes to its error stream and exits with status
// 0. With Hold it connects to that address first and does not exit by itself:
// it ends when the bridge kills it or when the test closes the connection.
type outputChild struct {
	Response string
	Stdout   int
	Stderr   int
	Hold     string
}

func TestMain(m *testing.M) {
	if len(os.Args) == 5 && os.Args[1] == "--input-type=commonjs" && os.Args[2] == "--eval" && filepath.Base(os.Args[4]) == outputChildFile {
		runOutputChild(os.Args[4])
	}
	os.Exit(m.Run())
}

func runOutputChild(path string) {
	var child outputChild
	data, err := os.ReadFile(path)
	if err != nil || json.Unmarshal(data, &child) != nil || child.Stdout < len(child.Response) {
		os.Exit(2)
	}
	var hold net.Conn
	if child.Hold != "" {
		// A stream that the bridge has closed must not end the child either.
		signal.Ignore(syscall.SIGPIPE)
		if hold, err = net.Dial("tcp", child.Hold); err != nil {
			os.Exit(2)
		}
	}
	io.Copy(io.Discard, os.Stdin)
	os.Stdout.Write(append([]byte(child.Response), bytes.Repeat([]byte{' '}, child.Stdout-len(child.Response))...))
	os.Stderr.Write(bytes.Repeat([]byte{' '}, child.Stderr))
	if hold != nil {
		io.Copy(io.Discard, hold)
		os.Exit(3)
	}
	os.Exit(0)
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
	// describe writes what the child does to the file that names the child.
	describe := func(t *testing.T, child outputChild) string {
		t.Helper()
		data, err := json.Marshal(child)
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(t.TempDir(), outputChildFile)
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		return path
	}
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
		child := outputChild{Response: response, Stdout: c.stdout, Stderr: c.stderr}
		if !c.over {
			t.Run(c.name, func(t *testing.T) {
				if out, err := check(describe(t, child)); err != nil || out.Status != statusIncomplete {
					t.Errorf("the bridge returned status %q and the error %v", out.Status, err)
				}
				awaitNoProcessGoroutine(t)
			})
			continue
		}
		t.Run(c.name+", and the process exits", func(t *testing.T) {
			out, err := check(describe(t, child))
			refused(t, out, err)
			awaitNoProcessGoroutine(t)
		})
		t.Run(c.name+", and the process stays", func(t *testing.T) {
			// The bridge's own deadline is 30 seconds. A bridge that returns
			// later than this did not cancel its helper at the bound.
			const wait = 15 * time.Second
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()
			child.Hold = listener.Addr().String()
			path := describe(t, child)
			connected := make(chan net.Conn, 1)
			go func() {
				connection, _ := listener.Accept()
				connected <- connection
			}()
			type result struct {
				out Output
				err error
			}
			returned := make(chan result, 1)
			go func() {
				out, err := check(path)
				returned <- result{out, err}
			}()
			var connection net.Conn
			select {
			case connection = <-connected:
			case r := <-returned:
				t.Fatal("the bridge returned before its helper had started", r.err)
			}
			if connection == nil {
				t.Fatal("the helper did not connect")
			}
			defer connection.Close()
			var r result
			select {
			case r = <-returned:
			case <-time.After(wait):
				t.Errorf("the bridge had not returned %v after its helper started: the helper was not cancelled at the bound", wait)
				connection.Close()
				r = <-returned
			}
			refused(t, r.out, r.err)
			// The child never closes this connection and never exits by itself,
			// so the end of the connection is the end of the process.
			connection.SetReadDeadline(time.Now().Add(wait))
			if _, err := io.Copy(io.Discard, connection); errors.Is(err, os.ErrDeadlineExceeded) {
				t.Error("the helper process outlived the call of the bridge")
			}
			awaitNoProcessGoroutine(t)
		})
	}
}

// The bridge waits for its process and for the copies of the process's
// streams, so no goroutine of os/exec outlives a call.
func awaitNoProcessGoroutine(t *testing.T) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	stacks := make([]byte, 1<<20)
	for {
		all := stacks[:runtime.Stack(stacks, true)]
		if !bytes.Contains(all, []byte("os/exec.")) {
			return
		}
		if time.Now().After(deadline) {
			t.Errorf("a goroutine of os/exec outlived the call of the bridge:\n%s", all)
			return
		}
		time.Sleep(time.Millisecond)
	}
}
