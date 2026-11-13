package processoutput

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"slices"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

// The child process of the tests below is this test binary. A test starts it
// with one argument, a file of this name, and a process that is started so is
// the child and runs no test.
const outputChildFile = "bittype-process-output-child.json"

// outputChild is what the child does. It writes Stdout bytes to its output
// stream and Stderr bytes to its error stream and exits with status 0. With
// Hold it connects to that address first and does not exit by itself: it ends
// when its parent kills it or when the test closes the connection.
type outputChild struct {
	Stdout int
	Stderr int
	Hold   string
}

func TestMain(m *testing.M) {
	if len(os.Args) == 2 && filepath.Base(os.Args[1]) == outputChildFile {
		runOutputChild(os.Args[1])
	}
	os.Exit(m.Run())
}

func runOutputChild(path string) {
	var child outputChild
	data, err := os.ReadFile(path)
	if err != nil || json.Unmarshal(data, &child) != nil {
		os.Exit(2)
	}
	var hold net.Conn
	if child.Hold != "" {
		// A stream that the parent has closed must not end the child either.
		signal.Ignore(syscall.SIGPIPE)
		if hold, err = net.Dial("tcp", child.Hold); err != nil {
			os.Exit(2)
		}
	}
	os.Stdout.Write(bytes.Repeat([]byte{'o'}, child.Stdout))
	os.Stderr.Write(bytes.Repeat([]byte{'e'}, child.Stderr))
	if hold != nil {
		io.Copy(io.Discard, hold)
		os.Exit(3)
	}
	os.Exit(0)
}

// counted returns a cancel function and the number of its calls.
func counted() (func(), *atomic.Int32) {
	var calls atomic.Int32
	return func() { calls.Add(1) }, &calls
}

// The two types have these exported methods and no other. Only Write takes
// bytes. A type that embeds bytes.Buffer also has ReadFrom, which io.Copy
// prefers to Write, and that is how os/exec copied a stream past a bound.
func TestOnlyWriteTakesBytes(t *testing.T) {
	for _, c := range []struct {
		value any
		want  []string
	}{
		{&Writer{}, []string{"Bytes", "Write"}},
		{&Bound{}, []string{"Exceeded", "Writer"}},
	} {
		kind := reflect.TypeOf(c.value)
		var names []string
		for i := range kind.NumMethod() {
			names = append(names, kind.Method(i).Name)
		}
		if !slices.Equal(names, c.want) {
			t.Errorf("%v has the methods %v, not %v", kind, names, c.want)
		}
	}
}

// A write that reaches the bound is taken whole. A write that would pass it
// is refused whole, cancels once and leaves the bound passed: a later write
// is refused too, also one that would have fitted.
func TestBoundRefusesAndCancelsOnce(t *testing.T) {
	cancel, calls := counted()
	bound := NewBound(5, cancel)
	w := bound.Writer()
	if count, err := w.Write([]byte("abc")); count != 3 || err != nil {
		t.Fatal("a write within the bound", count, err)
	}
	if count, err := w.Write([]byte("de")); count != 2 || err != nil || bound.Exceeded() || calls.Load() != 0 {
		t.Fatal("a write of exactly the rest of the bound", count, err, bound.Exceeded(), calls.Load())
	}
	if count, err := w.Write(nil); count != 0 || err != nil || bound.Exceeded() {
		t.Fatal("an empty write at the bound", count, err)
	}
	if count, err := w.Write([]byte("f")); count != 0 || !errors.Is(err, ErrBound) || !bound.Exceeded() || calls.Load() != 1 {
		t.Fatal("a write past the bound", count, err, bound.Exceeded(), calls.Load())
	}
	if count, err := w.Write(nil); count != 0 || !errors.Is(err, ErrBound) || calls.Load() != 1 {
		t.Fatal("a write after the bound was passed", count, err, calls.Load())
	}
	if got := string(w.Bytes()); got != "abcde" {
		t.Fatalf("the writer holds %q", got)
	}

	// Nothing of a refused write is kept, and room that is left stays unused.
	cancel, calls = counted()
	bound = NewBound(3, cancel)
	w = bound.Writer()
	if _, err := w.Write([]byte("ab")); err != nil {
		t.Fatal(err)
	}
	if count, err := w.Write([]byte("cd")); count != 0 || !errors.Is(err, ErrBound) {
		t.Fatal("a write that fits in part", count, err)
	}
	if count, err := w.Write([]byte("c")); count != 0 || !errors.Is(err, ErrBound) || calls.Load() != 1 {
		t.Fatal("a write that would have fitted, after a refused one", count, err, calls.Load())
	}
	if got := string(w.Bytes()); got != "ab" {
		t.Fatalf("the writer holds %q", got)
	}

	// A bound below zero refuses every write.
	cancel, calls = counted()
	bound = NewBound(-1, cancel)
	if _, err := bound.Writer().Write(nil); !errors.Is(err, ErrBound) || !bound.Exceeded() || calls.Load() != 1 {
		t.Fatal("an empty write on a bound below zero", err)
	}
}

// Writers of one bound share it, and a passed bound refuses all of them.
// Writers of two bounds do not touch each other.
func TestWritersShareTheirBound(t *testing.T) {
	cancel, calls := counted()
	shared, other := NewBound(5, cancel), NewBound(5, cancel)
	first, second, apart := shared.Writer(), shared.Writer(), other.Writer()
	if _, err := first.Write([]byte("123")); err != nil {
		t.Fatal(err)
	}
	if _, err := second.Write([]byte("45")); err != nil {
		t.Fatal(err)
	}
	if _, err := apart.Write([]byte("12345")); err != nil || other.Exceeded() {
		t.Fatal("a writer of another bound", err)
	}
	if _, err := second.Write([]byte("6")); !errors.Is(err, ErrBound) || !shared.Exceeded() || other.Exceeded() || calls.Load() != 1 {
		t.Fatal("the second writer passes the shared bound", err, calls.Load())
	}
	if _, err := first.Write(nil); !errors.Is(err, ErrBound) || calls.Load() != 1 {
		t.Fatal("the first writer after the second passed the shared bound", err, calls.Load())
	}
	if _, err := apart.Write(nil); err != nil {
		t.Fatal("a writer of another bound after the shared one was passed", err)
	}
	if string(first.Bytes()) != "123" || string(second.Bytes()) != "45" || string(apart.Bytes()) != "12345" {
		t.Fatalf("the writers hold %q, %q and %q", first.Bytes(), second.Bytes(), apart.Bytes())
	}
}

// No value of either type collects without a bound: a bound needs its cancel
// function, a writer that no bound returned refuses, and so does a writer of
// the zero bound.
func TestNoWriterWithoutABound(t *testing.T) {
	func() {
		defer func() {
			if recover() == nil {
				t.Error("a bound without a cancel function was returned")
			}
		}()
		NewBound(1, nil)
	}()
	var w Writer
	if count, err := w.Write([]byte("a")); count != 0 || !errors.Is(err, ErrBound) || w.Bytes() != nil {
		t.Error("a writer that no bound returned", count, err)
	}
	var zero Bound
	w2 := zero.Writer()
	if count, err := w2.Write([]byte("a")); count != 0 || !errors.Is(err, ErrBound) || !zero.Exceeded() || len(w2.Bytes()) != 0 {
		t.Error("a writer of the zero bound", count, err)
	}
}

// os/exec copies the output stream and the error stream of a process in two
// goroutines. Writers that share a bound never hold more than the bound
// together, each holds only its own stream, and the bound cancels once.
func TestConcurrentWritersKeepTheBound(t *testing.T) {
	const writers, each, limit = 8, 4096, 8*4096 - 1000
	cancel, calls := counted()
	bound := NewBound(limit, cancel)
	collected := make([]*Writer, writers)
	start := make(chan struct{})
	var done sync.WaitGroup
	for i := range collected {
		collected[i] = bound.Writer()
		done.Add(1)
		go func() {
			defer done.Done()
			<-start
			for range each {
				collected[i].Write([]byte{byte('a' + i)})
			}
		}()
	}
	close(start)
	done.Wait()
	total := 0
	for i, w := range collected {
		data := w.Bytes()
		total += len(data)
		if len(bytes.Trim(data, string(rune('a'+i)))) != 0 {
			t.Errorf("writer %d holds bytes of another stream", i)
		}
	}
	if total != limit || !bound.Exceeded() || calls.Load() != 1 {
		t.Errorf("the writers hold %d bytes of a bound of %d, and the bound cancelled %d times", total, limit, calls.Load())
	}
}

// The writers receive the streams of a real process through os/exec, the way
// every caller uses them. Streams of exactly their bound are read whole. One
// byte more is a passed bound and a cancelled process, whether the process
// then exits or stays. One bound per stream bounds each stream by itself, and
// one bound for both streams bounds their sum.
func TestProcessStreamsKeepTheirBound(t *testing.T) {
	const outputBound, errorBound = 256 << 10, 64 << 10
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
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
	type result struct {
		err            error
		output, errors *Bound
		stdout, stderr []byte
		cancelled      bool
	}
	// collect starts the child with one bound per stream, or with one bound
	// of the sum of both for the two streams, and waits for it.
	collect := func(path string, shared bool) result {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var r result
		if shared {
			r.output = NewBound(outputBound+errorBound, cancel)
			r.errors = r.output
		} else {
			r.output, r.errors = NewBound(outputBound, cancel), NewBound(errorBound, cancel)
		}
		stdout, stderr := r.output.Writer(), r.errors.Writer()
		command := exec.CommandContext(ctx, self, path)
		command.Stdout, command.Stderr = stdout, stderr
		command.WaitDelay = time.Second
		r.err = command.Run()
		r.stdout, r.stderr, r.cancelled = stdout.Bytes(), stderr.Bytes(), ctx.Err() != nil
		return r
	}
	for _, c := range []struct {
		name           string
		stdout, stderr int
		shared         bool
		// over names the stream whose bytes pass a bound: "", "output" or "error".
		over string
	}{
		{"each stream at its own bound", outputBound, errorBound, false, ""},
		{"the output stream one byte over its own bound", outputBound + 1, errorBound, false, "output"},
		{"the error stream one byte over its own bound", outputBound, errorBound + 1, false, "error"},
		{"both streams at their shared bound, with a longer output stream", outputBound + 1000, errorBound - 1000, true, ""},
		{"both streams at their shared bound, with a longer error stream", outputBound - 1000, errorBound + 1000, true, ""},
		{"both streams one byte over their shared bound, with a longer output stream", outputBound + 1001, errorBound - 1000, true, "output"},
		{"both streams one byte over their shared bound, with a longer error stream", outputBound - 1000, errorBound + 1001, true, "error"},
	} {
		child := outputChild{Stdout: c.stdout, Stderr: c.stderr}
		if c.over == "" {
			t.Run(c.name, func(t *testing.T) {
				r := collect(describe(t, child), c.shared)
				if r.err != nil || r.output.Exceeded() || r.errors.Exceeded() || r.cancelled || len(r.stdout) != c.stdout || len(r.stderr) != c.stderr {
					t.Errorf("the process ended with %v, and the writers hold %d and %d bytes", r.err, len(r.stdout), len(r.stderr))
				}
				if len(bytes.Trim(r.stdout, "o")) != 0 || len(bytes.Trim(r.stderr, "e")) != 0 {
					t.Error("a writer holds bytes of the other stream")
				}
			})
			continue
		}
		passed := func(t *testing.T, r result) {
			t.Helper()
			// With one bound per stream, only the bound of the longer stream is
			// passed, and what counts against it is that stream alone.
			over, other, held, limit := r.output, r.errors, len(r.stdout), outputBound
			if c.over == "error" {
				over, other, held, limit = r.errors, r.output, len(r.stderr), errorBound
			}
			if c.shared {
				held, limit = len(r.stdout)+len(r.stderr), outputBound+errorBound
			}
			if !over.Exceeded() || !c.shared && other.Exceeded() || !r.cancelled || held > limit {
				t.Errorf("bound passed: %v, the other bound passed: %v, cancelled: %v, and %d bytes are held of a bound of %d", over.Exceeded(), other.Exceeded(), r.cancelled, held, limit)
			}
		}
		t.Run(c.name+", and the process exits", func(t *testing.T) {
			passed(t, collect(describe(t, child), c.shared))
		})
		t.Run(c.name+", and the process stays", func(t *testing.T) {
			// Nothing but the bound's cancel function ends this process. A call
			// that has not returned after this time was not cancelled.
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
			returned := make(chan result, 1)
			go func() { returned <- collect(path, c.shared) }()
			var connection net.Conn
			select {
			case connection = <-connected:
			case r := <-returned:
				t.Fatal("the process ended before it had connected", r.err)
			}
			if connection == nil {
				t.Fatal("the process did not connect")
			}
			defer connection.Close()
			var r result
			select {
			case r = <-returned:
			case <-time.After(wait):
				t.Errorf("the process had not ended %v after it started: it was not cancelled at the bound", wait)
				connection.Close()
				r = <-returned
			}
			passed(t, r)
			// The child never closes this connection and never exits by itself,
			// so the end of the connection is the end of the process.
			connection.SetReadDeadline(time.Now().Add(wait))
			if _, err := io.Copy(io.Discard, connection); errors.Is(err, os.ErrDeadlineExceeded) {
				t.Error("the process outlived the call that started it")
			}
		})
	}
}
