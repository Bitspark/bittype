package golangsem

import (
	"context"
	"errors"
)

// errEnded reports native work whose goroutine ended without returning.
var errEnded = errors.New("golangsem: native work ended without a result")

// Beside runs native work beside the call and returns the context's error as
// soon as the context ends. go/types cannot be interrupted, in a package check
// as little as in a type comparison or an instantiation, so every entry that
// reaches it runs its whole native work through here.
//
// Abandoned work keeps its goroutine until it ends by itself, and its result is
// dropped. It may still read what the caller passed and whatever the work
// shares, so a caller that receives the context's error must not modify either.
// A panic of the work surfaces in the caller's goroutine, where it did when the
// work ran there.
func Beside[T any](ctx context.Context, work func() (T, error)) (T, error) {
	type outcome struct {
		value    T
		err      error
		returned bool
		panicked any
	}
	done := make(chan outcome, 1)
	go func() {
		var o outcome
		defer func() {
			if !o.returned {
				o.panicked = recover()
			}
			done <- o
		}()
		o.value, o.err = work()
		o.returned = true
	}()
	var zero T
	select {
	case o := <-done:
		switch {
		case o.returned:
			return o.value, o.err
		case o.panicked != nil:
			panic(o.panicked)
		default:
			// runtime.Goexit, or a nil panic where the runtime allows one.
			return zero, errEnded
		}
	case <-ctx.Done():
		return zero, ctx.Err()
	}
}
