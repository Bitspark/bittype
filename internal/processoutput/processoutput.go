// Package processoutput collects what a child process writes to its streams,
// up to a bound. It knows no language and no tool. A caller states the bound
// and the function that cancels its process, and decides what a passed bound
// means to its own callers.
package processoutput

import (
	"errors"
	"sync"
)

// ErrBound is the error of a refused write. A caller reads Bound.Exceeded for
// its decision: what os/exec returns for a process whose stream was refused
// depends on how the process ended.
var ErrBound = errors.New("process output passes its bound")

// A Bound is the number of bytes that the writers drawing on it may collect
// together. The first write that would pass it is refused and cancels the
// process, and from then on every write of every writer drawing on it is
// refused. One bound for one stream bounds that stream. One bound for several
// streams, of one process or of several, is a budget they share.
type Bound struct {
	mu        sync.Mutex
	remaining int
	exceeded  bool
	cancel    func()
}

// NewBound returns a bound of the given number of bytes. A bound below zero
// refuses every write. cancel ends the process whose streams the bound's
// writers collect. It is called once, by the first refused write, after the
// bound reports itself passed and outside the bound's lock, so it may call
// Exceeded. It is required: a bound that cannot end its process does not stop
// the process from writing.
func NewBound(bytes int, cancel func()) *Bound {
	if cancel == nil {
		panic("processoutput: a bound needs the function that cancels its process")
	}
	return &Bound{remaining: bytes, cancel: cancel}
}

// Writer returns a new, empty writer that draws on the bound.
func (b *Bound) Writer() *Writer {
	return &Writer{bound: b}
}

// Exceeded reports whether a write was refused.
func (b *Bound) Exceeded() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.exceeded
}

// A Writer collects one stream. Write is its only method that takes bytes.
// os/exec copies a stream with io.Copy, which calls the writer's ReadFrom when
// it has one, so a writer that embeds bytes.Buffer is filled past its bound
// without a call of Write. A Writer that no Bound returned refuses every
// write.
type Writer struct {
	bound *Bound
	data  []byte
}

// Write appends data when the bound has room for all of it. Otherwise it
// appends nothing, cancels the process if no write was refused before, and
// returns ErrBound. It is safe for the concurrent copies of several streams.
func (w *Writer) Write(data []byte) (int, error) {
	b := w.bound
	if b == nil {
		return 0, ErrBound
	}
	b.mu.Lock()
	if !b.exceeded && len(data) <= b.remaining {
		b.remaining -= len(data)
		w.data = append(w.data, data...)
		b.mu.Unlock()
		return len(data), nil
	}
	first := !b.exceeded
	b.exceeded = true
	b.mu.Unlock()
	// Outside the lock: the function belongs to the caller and may ask the
	// bound, which reports itself passed by now. The zero Bound has no
	// function and refuses all the same.
	if first && b.cancel != nil {
		b.cancel()
	}
	return 0, ErrBound
}

// Bytes returns the stream that the writer has collected, or nothing once
// its bound is passed. From the first refused write on, no writer of the
// bound takes another byte, so what each of them holds may be a stream with
// its end cut off, and a cut stream must not read as a whole one. The slice
// is the writer's own. A caller reads it after its process has ended and does
// not change it.
func (w *Writer) Bytes() []byte {
	if w.bound == nil {
		return nil
	}
	w.bound.mu.Lock()
	defer w.bound.mu.Unlock()
	if w.bound.exceeded {
		return nil
	}
	return w.data
}
