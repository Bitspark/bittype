// Package processoutputtest is the child process and the harness of the tests
// that bound what a child process writes. It knows no language and no tool:
// the child writes stated numbers of bytes to its streams, and the harness
// tells a call that cancelled its process from one that did not.
//
// A package that uses it keeps what is its own: how its test binary knows
// that it was started as the child, what the output stream starts with,
// whether the child reads its input stream, and the call under test.
package processoutputtest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
	"time"
)

// File is the name of the file that describes a child. A test binary that is
// started with such a file where its package passes one is the child.
const File = "bittype-process-output-child.json"

// Wait bounds each wait of the harness for an event. It only passes when the
// event does not come, and then the test fails.
const Wait = 15 * time.Second

// What the child fills its output stream and its error stream with. A space
// keeps a response that is a JSON value one.
const (
	OutputFill = ' '
	ErrorFill  = '.'
)

// Child is what a child process does. With Input it reads its input stream to
// the end. It writes Response and then OutputFill to its output stream until
// it has written Stdout bytes, writes Stderr bytes of ErrorFill to its error
// stream and exits with status 0. With Hold it connects to that address
// before all of this and does not exit by itself: it ends when its parent
// kills it or when the test closes the connection.
type Child struct {
	Input    bool
	Response string
	Stdout   int
	Stderr   int
	Hold     string
}

// Describe writes a description to a new file named File and returns its
// path. A description is a Child, or what a package's own TestMain reads to
// choose one.
func Describe(t testing.TB, description any) string {
	t.Helper()
	data, err := json.Marshal(description)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), File)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

// Described reads the description at a path in the child. A child that
// cannot read its description exits with status 2.
func Described(path string, description any) {
	data, err := os.ReadFile(path)
	if err != nil || json.Unmarshal(data, description) != nil {
		os.Exit(2)
	}
}

// Run does what the child does and ends the process. It does not return. A
// child whose description is not one, or that cannot connect, exits with
// status 2. A child that stays exits with status 3 when its connection ends.
func (c Child) Run() {
	if c.Stdout < len(c.Response) || c.Stderr < 0 {
		os.Exit(2)
	}
	var hold net.Conn
	if c.Hold != "" {
		// A stream that the parent has closed must not end the child either.
		signal.Ignore(syscall.SIGPIPE)
		var err error
		if hold, err = net.Dial("tcp", c.Hold); err != nil {
			os.Exit(2)
		}
	}
	if c.Input {
		io.Copy(io.Discard, os.Stdin)
	}
	os.Stdout.Write(append([]byte(c.Response), bytes.Repeat([]byte{OutputFill}, c.Stdout-len(c.Response))...))
	os.Stderr.Write(bytes.Repeat([]byte{ErrorFill}, c.Stderr))
	if hold != nil {
		io.Copy(io.Discard, hold)
		os.Exit(3)
	}
	os.Exit(0)
}

// A Hold is where a child that stays connects. The test puts its address
// into the Hold of the child's description.
type Hold struct {
	listener net.Listener
}

// Listen opens a hold for the rest of the test.
func Listen(t testing.TB) *Hold {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { listener.Close() })
	return &Hold{listener: listener}
}

// Address is the address at which the child connects.
func (h *Hold) Address() string {
	return h.listener.Addr().String()
}

// AwaitCancelled runs a call that starts one child with the hold's address
// and returns what the call returned. The child never exits by itself, so the
// call returns only when it ends its process. The harness waits for three
// events, each at most Wait: the child connects, the call returns, and the
// connection ends, which only the end of the process does. A call that has
// not returned by then did not cancel its process, and the test fails.
//
// Nothing orders the return of the call and the receipt of the connection in
// this process: a child that made its parent return had connected before it
// wrote, and the test may still learn of the return first. So a return that
// comes first is kept, and the connection is awaited after it.
func AwaitCancelled[R any](t testing.TB, h *Hold, call func() R) R {
	t.Helper()
	connected := make(chan net.Conn, 1)
	go func() {
		connection, _ := h.listener.Accept()
		connected <- connection
	}()
	returned := make(chan R, 1)
	go func() { returned <- call() }()
	var connection net.Conn
	var result R
	early := false
	select {
	case connection = <-connected:
	case result = <-returned:
		early = true
		select {
		case connection = <-connected:
		case <-time.After(Wait):
			t.Fatalf("the call returned and its process had not connected %v later: %+v", Wait, result)
		}
	}
	if connection == nil {
		t.Fatal("the process did not connect")
	}
	defer connection.Close()
	if !early {
		select {
		case result = <-returned:
		case <-time.After(Wait):
			t.Errorf("the call had not returned %v after its process started: the process was not cancelled", Wait)
			// The end of the connection ends the child, and with it the call.
			connection.Close()
			result = <-returned
		}
	}
	// The child never closes this connection and never exits by itself, so
	// the end of the connection is the end of the process.
	connection.SetReadDeadline(time.Now().Add(Wait))
	if _, err := io.Copy(io.Discard, connection); errors.Is(err, os.ErrDeadlineExceeded) {
		t.Error("the process outlived the call that started it")
	}
	return result
}

// AwaitNoProcessGoroutine waits until no goroutine of os/exec is left. A
// caller that waits for its process and for the copies of the process's
// streams leaves none behind a call.
func AwaitNoProcessGoroutine(t testing.TB) {
	t.Helper()
	deadline := time.Now().Add(Wait)
	stacks := make([]byte, 1<<20)
	for {
		all := stacks[:runtime.Stack(stacks, true)]
		if !bytes.Contains(all, []byte("os/exec.")) {
			return
		}
		if time.Now().After(deadline) {
			t.Errorf("a goroutine of os/exec outlived the call:\n%s", all)
			return
		}
		time.Sleep(time.Millisecond)
	}
}
