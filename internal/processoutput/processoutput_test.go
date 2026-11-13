package processoutput

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Bitspark/bittype/internal/processoutput/processoutputtest"
)

// The child process of the tests below is this test binary. A test starts it
// with one argument, a description file of the shared fixture, and a process
// that is started so is the child and runs no test.
func TestMain(m *testing.M) {
	if len(os.Args) == 2 && filepath.Base(os.Args[1]) == processoutputtest.File {
		var child processoutputtest.Child
		processoutputtest.Described(os.Args[1], &child)
		child.Run()
	}
	os.Exit(m.Run())
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
// is refused too, also one that would have fitted. A writer of a passed
// bound returns no bytes, because its stream may be cut.
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
	if got := string(w.Bytes()); got != "abcde" {
		t.Fatalf("at its bound the writer holds %q", got)
	}
	if count, err := w.Write([]byte("f")); count != 0 || !errors.Is(err, ErrBound) || !bound.Exceeded() || calls.Load() != 1 {
		t.Fatal("a write past the bound", count, err, bound.Exceeded(), calls.Load())
	}
	if count, err := w.Write(nil); count != 0 || !errors.Is(err, ErrBound) || calls.Load() != 1 {
		t.Fatal("a write after the bound was passed", count, err, calls.Load())
	}
	if got := w.Bytes(); got != nil {
		t.Fatalf("past its bound the writer returns %q", got)
	}

	// A write that fits in part is refused whole, and room that is left stays
	// unused.
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
	if got := w.Bytes(); got != nil {
		t.Fatalf("past its bound the writer returns %q", got)
	}

	// A bound below zero refuses every write.
	cancel, calls = counted()
	bound = NewBound(-1, cancel)
	if _, err := bound.Writer().Write(nil); !errors.Is(err, ErrBound) || !bound.Exceeded() || calls.Load() != 1 {
		t.Fatal("an empty write on a bound below zero", err)
	}
}

// Writers of one bound share it, and a passed bound refuses all of them and
// returns the bytes of none: no writer of it takes another byte, so each may
// hold a cut stream. Writers of two bounds do not touch each other.
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
	if string(first.Bytes()) != "123" || string(second.Bytes()) != "45" {
		t.Fatalf("at their bound the writers hold %q and %q", first.Bytes(), second.Bytes())
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
	if first.Bytes() != nil || second.Bytes() != nil || string(apart.Bytes()) != "12345" {
		t.Fatalf("past the shared bound the writers return %q, %q and %q", first.Bytes(), second.Bytes(), apart.Bytes())
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

// The cancel function belongs to the caller, so the bound calls it outside
// its own lock and after it reports itself passed. A cancel function that
// asks the bound and writes to it shows both: called inside the lock, it
// never returns, and this test fails when its wait passes.
func TestCancelMayAskTheBound(t *testing.T) {
	type seen struct {
		exceeded bool
		bytes    []byte
		err      error
	}
	asked := make(chan seen, 1)
	var bound *Bound
	var w *Writer
	bound = NewBound(2, func() {
		_, err := bound.Writer().Write(nil)
		asked <- seen{bound.Exceeded(), w.Bytes(), err}
	})
	w = bound.Writer()
	if _, err := w.Write([]byte("ab")); err != nil {
		t.Fatal(err)
	}
	written := make(chan error, 1)
	go func() {
		_, err := w.Write([]byte("c"))
		written <- err
	}()
	select {
	case err := <-written:
		if !errors.Is(err, ErrBound) {
			t.Fatal("a write past the bound", err)
		}
	case <-time.After(processoutputtest.Wait):
		t.Fatalf("a write past the bound had not returned after %v: its cancel function asked the bound and was called inside the bound's lock", processoutputtest.Wait)
	}
	// Write returned after cancel, so what cancel saw is there.
	select {
	case s := <-asked:
		if !s.exceeded || s.bytes != nil || !errors.Is(s.err, ErrBound) {
			t.Errorf("in its cancel function the bound reported passed: %v, returned the bytes %q, and answered a write with %v", s.exceeded, s.bytes, s.err)
		}
	default:
		t.Fatal("the write past the bound did not call the cancel function")
	}
}

// os/exec copies the output stream and the error stream of a process in two
// goroutines. Writers that share a bound never take more than the bound
// together, each holds only its own stream, and the bound cancels once.
// While they write, another goroutine reads the bound and the writers, as a
// caller may: under the race detector an unguarded read fails here.
func TestConcurrentWritersKeepTheBound(t *testing.T) {
	const writers, each = 8, 4096
	// write releases the writers at once on one bound and returns the writers
	// and how many bytes Write took in all.
	write := func(bound *Bound) ([]*Writer, int) {
		collected := make([]*Writer, writers)
		for i := range collected {
			collected[i] = bound.Writer()
		}
		start := make(chan struct{})
		var taken atomic.Int64
		var done sync.WaitGroup
		for i := range collected {
			done.Add(1)
			go func() {
				defer done.Done()
				<-start
				for range each {
					count, _ := collected[i].Write([]byte{byte('a' + i)})
					taken.Add(int64(count))
				}
			}()
		}
		stop := make(chan struct{})
		var reader sync.WaitGroup
		reader.Add(1)
		go func() {
			defer reader.Done()
			<-start
			for {
				select {
				case <-stop:
					return
				default:
					bound.Exceeded()
					for _, w := range collected {
						w.Bytes()
					}
				}
			}
		}()
		close(start)
		done.Wait()
		close(stop)
		reader.Wait()
		return collected, int(taken.Load())
	}

	// Exactly the bound: every write is taken, and each writer holds its own.
	cancel, calls := counted()
	bound := NewBound(writers*each, cancel)
	collected, taken := write(bound)
	if taken != writers*each || bound.Exceeded() || calls.Load() != 0 {
		t.Errorf("at the bound the writers took %d bytes of %d, and the bound cancelled %d times", taken, writers*each, calls.Load())
	}
	for i, w := range collected {
		if data := w.Bytes(); len(data) != each || len(bytes.Trim(data, string(rune('a'+i)))) != 0 {
			t.Errorf("writer %d holds %d bytes, or bytes of another stream", i, len(data))
		}
	}

	// More than the bound: the writers take exactly the bound and no more.
	const limit = writers*each - 1000
	cancel, calls = counted()
	bound = NewBound(limit, cancel)
	collected, taken = write(bound)
	if taken != limit || !bound.Exceeded() || calls.Load() != 1 {
		t.Errorf("past the bound the writers took %d bytes of a bound of %d, and the bound cancelled %d times", taken, limit, calls.Load())
	}
	for i, w := range collected {
		if w.Bytes() != nil {
			t.Errorf("writer %d returns bytes past the bound", i)
		}
	}
}

// The writers receive the streams of a real process through os/exec, the way
// every caller uses them. Streams of exactly their bound are read whole. One
// byte more is a passed bound, a cancelled process and no bytes, whether the
// process then exits or stays. One bound per stream bounds each stream by
// itself, and one bound for both streams bounds their sum.
func TestProcessStreamsKeepTheirBound(t *testing.T) {
	const outputBound, errorBound = 256 << 10, 64 << 10
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
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
		// This call is the test's own, and so is its wait for the copies of the
		// streams after the process has ended: it only keeps a copy that does
		// not end from holding the test. A caller's shorter wait is the
		// caller's limit and is not what this test is about.
		command.WaitDelay = processoutputtest.Wait
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
		child := processoutputtest.Child{Stdout: c.stdout, Stderr: c.stderr}
		if c.over == "" {
			t.Run(c.name, func(t *testing.T) {
				r := collect(processoutputtest.Describe(t, child), c.shared)
				if r.err != nil || r.output.Exceeded() || r.errors.Exceeded() || r.cancelled || len(r.stdout) != c.stdout || len(r.stderr) != c.stderr {
					t.Errorf("the process ended with %v, and the writers hold %d and %d bytes", r.err, len(r.stdout), len(r.stderr))
				}
				if len(bytes.Trim(r.stdout, string(rune(processoutputtest.OutputFill)))) != 0 || len(bytes.Trim(r.stderr, string(rune(processoutputtest.ErrorFill)))) != 0 {
					t.Error("a writer holds bytes of the other stream")
				}
			})
			continue
		}
		passed := func(t *testing.T, r result) {
			t.Helper()
			// With one bound per stream, only the bound of the longer stream is
			// passed. The other stream keeps its own bound: the cancellation may
			// have cut it, and it never holds more.
			over, other, cut, rest, limit := r.output, r.errors, r.stdout, r.stderr, errorBound
			if c.over == "error" {
				over, other, cut, rest, limit = r.errors, r.output, r.stderr, r.stdout, outputBound
			}
			if !over.Exceeded() || !r.cancelled || cut != nil {
				t.Errorf("bound passed: %v, cancelled: %v, and the writer of the passed bound returns %d bytes", over.Exceeded(), r.cancelled, len(cut))
			}
			if c.shared && rest != nil {
				t.Errorf("the other writer of the passed bound returns %d bytes", len(rest))
			}
			if !c.shared && (other.Exceeded() || len(rest) > limit) {
				t.Errorf("the other bound passed: %v, and its writer holds %d bytes of a bound of %d", other.Exceeded(), len(rest), limit)
			}
		}
		t.Run(c.name+", and the process exits", func(t *testing.T) {
			passed(t, collect(processoutputtest.Describe(t, child), c.shared))
		})
		t.Run(c.name+", and the process stays", func(t *testing.T) {
			// Nothing but the bound's cancel function ends this process.
			hold := processoutputtest.Listen(t)
			child.Hold = hold.Address()
			path := processoutputtest.Describe(t, child)
			passed(t, processoutputtest.AwaitCancelled(t, hold, func() result { return collect(path, c.shared) }))
		})
	}
}
