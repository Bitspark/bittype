package generation_test

import (
	"context"
	"errors"
	"testing"

	gen "github.com/Bitspark/bittype/contracts/generation/go"
)

type ledgerNative struct {
	*nativeStore
	publications []gen.Publication[sourceLanguage, nativeLanguage]
	definitions  []gen.Definition[sourceLanguage, nativeLanguage]
}

func (n *ledgerNative) Assemble(ctx context.Context, publications []gen.Publication[sourceLanguage, nativeLanguage], definitions []gen.Definition[sourceLanguage, nativeLanguage], limits gen.Limits) (gen.Assembly, error) {
	n.publications = append([]gen.Publication[sourceLanguage, nativeLanguage](nil), publications...)
	n.definitions = append([]gen.Definition[sourceLanguage, nativeLanguage](nil), definitions...)
	return n.nativeStore.Assemble(ctx, publications, definitions, limits)
}

func TestCopiedSessionHandlesShareTheCompleteInvocation(t *testing.T) {
	f := newFixture(t)
	native := &ledgerNative{nativeStore: f.native}
	c := f.configuration(nil)
	c.Native = native
	c.Plans = append(c.Plans, f.plainPlan(true))
	s := begin(t, context.Background(), f, c, gen.DefaultLimits())
	acquiring, publishing := *s, *s

	owned, err := acquiring.Acquire(f.plain)
	if err != nil {
		t.Fatal(err)
	}
	publication, err := publishing.Publish(f.plain)
	if err != nil {
		t.Fatal(err)
	}
	if publication.Callee().Signature() != owned.Signature() {
		t.Fatal("copied handle acquired a different canonical callee")
	}
	// A callee acquired through one handle can be defined through another.
	definition, err := publishing.Define(owned)
	if err != nil {
		t.Fatal(err)
	}
	out, err := s.Finish()
	if err != nil || !out.Snapshot().Valid() {
		t.Fatalf("original handle could not commit the shared invocation: %v", err)
	}
	if len(native.publications) != 1 || native.publications[0].Root() != f.plain || native.publications[0].Callee().Signature() != owned.Signature() {
		t.Fatal("publication through a copied handle was absent from assembly")
	}
	if len(native.definitions) != 1 || native.definitions[0].Callee().Signature() != owned.Signature() || native.definitions[0].Body() != definition.Body() {
		t.Fatal("definition through a copied handle was absent from assembly")
	}
	failedFinish(t, &acquiring, gen.ErrInvalid)
	failedFinish(t, &publishing, gen.ErrInvalid)
	if native.assemblies != 1 {
		t.Fatal("a copied handle committed the finished invocation again")
	}
}

func TestFailureThroughAnySessionCopyPoisonsAllHandles(t *testing.T) {
	for _, failOriginal := range []bool{false, true} {
		name := "copied-handle"
		if failOriginal {
			name = "original-handle"
		}
		t.Run(name, func(t *testing.T) {
			f := newFixture(t)
			s := begin(t, context.Background(), f, f.configuration(nil), gen.DefaultLimits())
			copied := *s
			failing, observing := &copied, s
			if failOriginal {
				failing, observing = s, &copied
			}
			if out, err := failing.Use(use{}); !errors.Is(err, gen.ErrScope) || out.Valid() {
				t.Fatalf("invalid use did not fail the invocation: %v", err)
			}
			if out, err := observing.Use(f.scalar); !errors.Is(err, gen.ErrScope) || out.Valid() {
				t.Fatalf("failure did not reach the other session handle: %v", err)
			}
			failedFinish(t, s, gen.ErrScope)
			failedFinish(t, &copied, gen.ErrScope)
			if f.native.assemblies != 0 {
				t.Fatal("a copied handle assembled a failed invocation")
			}
		})
	}
}

func TestAssemblyCannotReenterTheInvocation(t *testing.T) {
	for _, operation := range []string{"use", "acquire", "helper", "publish", "finish"} {
		t.Run(operation, func(t *testing.T) {
			f := newFixture(t)
			c := f.configuration(nil)
			helper := f.plainPlan(true)
			helper.Key, helper.Declaration, helper.Publisher = "helper", declaration{}, nil
			c.Plans = append(c.Plans, f.plainPlan(false), helper)
			s := begin(t, context.Background(), f, c, gen.DefaultLimits())
			copied := *s
			attempted, published := false, false
			var reentryErr error
			f.native.assemblyHook = func() error {
				// Bound this deliberately faulty collaborator even if recursive
				// Finish is not rejected before a second call to Assemble.
				if attempted {
					return nil
				}
				attempted = true
				switch operation {
				case "use":
					var out product
					out, reentryErr = copied.Use(f.scalar)
					published = out.Valid()
				case "acquire":
					var out callee
					out, reentryErr = copied.Acquire(f.plain)
					published = out.Valid()
				case "helper":
					var out callee
					out, reentryErr = copied.AcquireHelper("helper")
					published = out.Valid()
				case "publish":
					var out gen.Publication[sourceLanguage, nativeLanguage]
					out, reentryErr = copied.Publish(f.plain)
					published = out.Valid()
				case "finish":
					var out gen.Output
					out, reentryErr = copied.Finish()
					published = out.Snapshot().Valid()
				}
				// Swallow the nested failure and return the ordinary valid AST.
				// The host must still prevent the outer commit from succeeding.
				return nil
			}
			failedFinish(t, s, gen.ErrInvalid)
			if !attempted || !errors.Is(reentryErr, gen.ErrInvalid) || published {
				t.Fatalf("assembly reentry escaped the commit boundary: attempted=%v published=%v error=%v", attempted, published, reentryErr)
			}
			failedFinish(t, &copied, gen.ErrInvalid)
			if f.native.assemblies != 1 {
				t.Fatalf("reentry invoked assembly again: %d calls", f.native.assemblies)
			}
		})
	}
}

func TestZeroSessionHandlesReturnErrorsWithoutOutput(t *testing.T) {
	operations := []struct {
		name string
		run  func(*session) (bool, error)
	}{
		{"use", func(s *session) (bool, error) {
			out, err := s.Use(use{})
			return out.Valid(), err
		}},
		{"acquire", func(s *session) (bool, error) {
			out, err := s.Acquire(declaration{})
			return out.Valid(), err
		}},
		{"helper", func(s *session) (bool, error) {
			out, err := s.AcquireHelper("helper")
			return out.Valid(), err
		}},
		{"publish", func(s *session) (bool, error) {
			out, err := s.Publish(declaration{})
			return out.Valid(), err
		}},
		{"define", func(s *session) (bool, error) {
			out, err := s.Define(callee{})
			return out.Valid(), err
		}},
		{"finish", func(s *session) (bool, error) {
			out, err := s.Finish()
			return out.Snapshot().Valid(), err
		}},
	}
	for _, operation := range operations {
		for _, nilHandle := range []bool{false, true} {
			name := operation.name + "/zero"
			if nilHandle {
				name = operation.name + "/nil"
			}
			t.Run(name, func(t *testing.T) {
				var zero session
				handle := &zero
				if nilHandle {
					handle = nil
				}
				if published, err := operation.run(handle); !errors.Is(err, gen.ErrInvalid) || published {
					t.Fatalf("invalid session returned a product or wrong error: published=%v error=%v", published, err)
				}
			})
		}
	}
}

func TestRecursivePublicationOfTheSameRootFailsTheInvocation(t *testing.T) {
	f := newFixture(t)
	var s *session
	p := f.plainPlan(false)
	original := p.Publisher
	calls, nestedPublished := 0, false
	var nestedErr error
	p.Publisher = publishFunc(func(ctx context.Context, r gen.PublishRequest[sourceLanguage, nativeLanguage]) (gen.Publication[sourceLanguage, nativeLanguage], error) {
		calls++
		current, err := original.Publish(ctx, r)
		if err != nil {
			return gen.Publication[sourceLanguage, nativeLanguage]{}, err
		}
		if calls == 1 {
			var nested gen.Publication[sourceLanguage, nativeLanguage]
			nested, nestedErr = s.Publish(r.Root())
			nestedPublished = nested.Valid()
		}
		// Returning a previously completed current result cannot hide the
		// recursive operation's failure from the invocation ledger.
		return current, nil
	})
	c := f.configuration(nil)
	c.Plans = append(c.Plans, p)
	s = begin(t, context.Background(), f, c, gen.DefaultLimits())
	if out, err := s.Publish(f.plain); !errors.Is(err, gen.ErrInvalid) || out.Valid() {
		t.Fatalf("recursive publication succeeded: %v", err)
	}
	if calls != 1 || !errors.Is(nestedErr, gen.ErrInvalid) || nestedPublished {
		t.Fatalf("same-root reentry reached its publisher: calls=%d published=%v error=%v", calls, nestedPublished, nestedErr)
	}
	failedFinish(t, s, gen.ErrInvalid)
	if f.native.assemblies != 0 {
		t.Fatal("recursive publication reached output assembly")
	}
}
