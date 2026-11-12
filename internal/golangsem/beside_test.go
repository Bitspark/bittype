package golangsem

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"
)

// The work is provably still running, held by the test, when the context ends.
// The call returns the context's error without waiting for it.
func TestBesideReturnsWhenContextEnds(t *testing.T) {
	entered, release := make(chan struct{}), make(chan struct{})
	defer close(release) // lets the abandoned work end
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	returned := make(chan error, 1)
	go func() {
		_, err := Beside(ctx, func() (int, error) {
			close(entered)
			<-release
			return 1, nil
		})
		returned <- err
	}()
	// The bounds only turn a hang into a failure; nothing is inferred from time.
	select {
	case <-entered:
	case <-time.After(time.Minute):
		t.Fatal("the work never started")
	}
	cancel()
	select {
	case err := <-returned:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Minute):
		t.Fatal("the call waited for work that cannot end")
	}
}

// Work that ends is returned as it is, its error included.
func TestBesideReturnsTheWorksResult(t *testing.T) {
	if v, err := Beside(context.Background(), func() (int, error) { return 7, nil }); v != 7 || err != nil {
		t.Fatal(v, err)
	}
	failed := errors.New("failed")
	if v, err := Beside(context.Background(), func() (int, error) { return 3, failed }); v != 3 || err != failed {
		t.Fatal(v, err)
	}
}

// A panic inside the work surfaces in the caller's goroutine, where it did
// before the work ran beside the call.
func TestBesidePanicReachesTheCaller(t *testing.T) {
	var recovered any
	func() {
		defer func() { recovered = recover() }()
		_, _ = Beside(context.Background(), func() (int, error) { panic("work panic") })
	}()
	if recovered != "work panic" {
		t.Fatal("recovered", recovered)
	}
}

// Work whose goroutine ends without returning is an error, never a zero value
// that reads as success.
func TestBesideWorkThatEndsWithoutReturningIsAnError(t *testing.T) {
	v, err := Beside(context.Background(), func() (*int, error) {
		runtime.Goexit()
		return nil, nil
	})
	if v != nil || !errors.Is(err, errEnded) {
		t.Fatal(v, err)
	}
}
