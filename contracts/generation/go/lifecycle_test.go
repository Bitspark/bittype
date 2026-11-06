package generation_test

import (
	"context"
	"errors"
	"testing"

	gen "github.com/Bitspark/bittype/contracts/generation/go"
)

func TestMissingCapabilitiesFailOnlyWhenDemanded(t *testing.T) {
	for _, name := range []string{"use", "publish", "define", "convention", "acquire", "header"} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t)
			c := f.configuration(nil)
			p := f.plainPlan(true)
			switch name {
			case "use":
				c.Plans[0].User = nil
			case "publish":
				p.Publisher = nil
			case "define":
				p.Definer = nil
			case "convention":
				p = f.genericPlan(false, nil)
				p.Convention = nil
			case "acquire":
				p.Acquire = nil
			case "header":
				p.Header = nil
			}
			c.Plans = append(c.Plans, p)
			// Construction does not demand every capability from every mapping.
			s := begin(t, context.Background(), f, c, gen.DefaultLimits())
			var err error
			switch name {
			case "use":
				var out product
				out, err = s.Use(f.scalar)
				if out.Valid() {
					t.Fatal("missing Use returned a product")
				}
			case "publish":
				var out gen.Publication[sourceLanguage, nativeLanguage]
				out, err = s.Publish(f.plain)
				if out.Valid() {
					t.Fatal("missing Publish returned a publication")
				}
			case "define":
				owned, acquireErr := s.Acquire(f.plain)
				if acquireErr != nil {
					t.Fatal(acquireErr)
				}
				var out gen.Definition[sourceLanguage, nativeLanguage]
				out, err = s.Define(owned)
				if out.Valid() {
					t.Fatal("missing Define returned a definition")
				}
			case "convention":
				_, err = s.Acquire(f.generic)
			default:
				_, err = s.Acquire(f.plain)
			}
			if !errors.Is(err, gen.ErrUnsupported) {
				t.Fatalf("missing capability did not fail explicitly: %v", err)
			}
			failedFinish(t, s, gen.ErrUnsupported)
		})
	}
}

func TestNongenericAndExternalMappingsOmitUnneededCapabilities(t *testing.T) {
	t.Run("nongeneric-without-convention", func(t *testing.T) {
		f := newFixture(t)
		c := f.configuration(nil)
		c.Plans = append(c.Plans, f.plainPlan(true))
		s := begin(t, context.Background(), f, c, gen.DefaultLimits())
		publication, err := s.Publish(f.plain)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := s.Define(publication.Callee()); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Finish(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("external-without-publish-or-define", func(t *testing.T) {
		f := newFixture(t)
		p := f.genericPlan(false, nil)
		p.Definer = nil
		c := f.configuration(nil)
		c.Plans = append(c.Plans, p)
		c.Select = func(u useData) (string, error) {
			if u.Callee == f.generic {
				return "generic", nil
			}
			return "direct", nil
		}
		s := begin(t, context.Background(), f, c, gen.DefaultLimits())
		if _, err := s.Use(f.application); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Finish(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestPublicationFromAnotherInvocationCannotBeReplayed(t *testing.T) {
	f := newFixture(t)
	var previous gen.Publication[sourceLanguage, nativeLanguage]
	p := f.plainPlan(false)
	p.Publisher = publishFunc(func(_ context.Context, r gen.PublishRequest[sourceLanguage, nativeLanguage]) (gen.Publication[sourceLanguage, nativeLanguage], error) {
		callee, err := r.Acquire()
		if err != nil {
			return gen.Publication[sourceLanguage, nativeLanguage]{}, err
		}
		current, err := r.Expose(callee)
		if err != nil {
			return gen.Publication[sourceLanguage, nativeLanguage]{}, err
		}
		if previous.Valid() {
			return previous, nil
		}
		return current, nil
	})
	c := f.configuration(nil)
	c.Plans = append(c.Plans, p)
	first := begin(t, context.Background(), f, c, gen.DefaultLimits())
	var err error
	previous, err = first.Publish(f.plain)
	if err != nil {
		t.Fatal(err)
	}
	second := begin(t, context.Background(), f, c, gen.DefaultLimits())
	if out, err := second.Publish(f.plain); !errors.Is(err, gen.ErrInvalid) || out.Valid() {
		t.Fatalf("saved publication replaced a current exposure: %v", err)
	}
	failedFinish(t, second, gen.ErrInvalid)
}

func TestDefinitionCannotBeReplayedForAnotherCallee(t *testing.T) {
	f := newFixture(t)
	var previous gen.Definition[sourceLanguage, nativeLanguage]
	replay := defineFunc(func(_ context.Context, r gen.DefineRequest[sourceLanguage, nativeLanguage]) (gen.Definition[sourceLanguage, nativeLanguage], error) {
		current, err := r.CompleteBody(f.closed)
		if err != nil {
			return gen.Definition[sourceLanguage, nativeLanguage]{}, err
		}
		if previous.Valid() {
			return previous, nil
		}
		return current, nil
	})
	plain, generic := f.plainPlan(true), f.genericPlan(true, nil)
	plain.Definer, generic.Definer = replay, replay
	c := f.configuration(nil)
	c.Plans = append(c.Plans, plain, generic)
	s := begin(t, context.Background(), f, c, gen.DefaultLimits())
	one, err := s.Acquire(f.plain)
	if err != nil {
		t.Fatal(err)
	}
	two, err := s.Acquire(f.generic)
	if err != nil {
		t.Fatal(err)
	}
	previous, err = s.Define(one)
	if err != nil {
		t.Fatal(err)
	}
	if out, err := s.Define(two); !errors.Is(err, gen.ErrInvalid) || out.Valid() {
		t.Fatalf("saved definition filled another callee: %v", err)
	}
	failedFinish(t, s, gen.ErrInvalid)
}

func TestDefinitionCanOnlyFillItsReservationOnce(t *testing.T) {
	f := newFixture(t)
	c := f.configuration(nil)
	c.Plans = append(c.Plans, f.plainPlan(true))
	s := begin(t, context.Background(), f, c, gen.DefaultLimits())
	owned, err := s.Acquire(f.plain)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Define(owned); err != nil {
		t.Fatal(err)
	}
	if out, err := s.Define(owned); !errors.Is(err, gen.ErrInvalid) || out.Valid() {
		t.Fatalf("reservation filled twice: %v", err)
	}
	failedFinish(t, s, gen.ErrInvalid)
}

func TestReentrantSignatureAcquisitionFails(t *testing.T) {
	for _, phase := range []string{"reservation", "introduction"} {
		t.Run(phase, func(t *testing.T) {
			f := newFixture(t)
			var s *session
			p := f.genericPlan(false, nil)
			if phase == "reservation" {
				p.Acquire = func(context.Context, site) (gen.Destination[nativeLanguage], bool, error) {
					_, err := s.Acquire(f.generic)
					return gen.Destination[nativeLanguage]{}, false, err
				}
			} else {
				original := p.Convention
				p.Convention = convention{
					introduce: func(context.Context, gen.IntroduceRequest[sourceLanguage, nativeLanguage]) (gen.Introduction[sourceLanguage, nativeLanguage], error) {
						_, err := s.Acquire(f.generic)
						return gen.Introduction[sourceLanguage, nativeLanguage]{}, err
					},
					apply: original.Apply,
				}
			}
			c := f.configuration(nil)
			c.Plans = append(c.Plans, p)
			s = begin(t, context.Background(), f, c, gen.DefaultLimits())
			if out, err := s.Acquire(f.generic); !errors.Is(err, gen.ErrInvalid) || out.Valid() {
				t.Fatalf("unfinished signature escaped: %v", err)
			}
			failedFinish(t, s, gen.ErrInvalid)
		})
	}
}

func TestInvalidOrReplacedFormalLayoutsFail(t *testing.T) {
	for _, name := range []string{"foreign-owner", "absent-record", "other-destination", "freeze-twice"} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t)
			p := f.genericPlan(false, nil)
			original := p.Convention
			p.Convention = convention{
				introduce: func(_ context.Context, r gen.IntroduceRequest[sourceLanguage, nativeLanguage]) (gen.Introduction[sourceLanguage, nativeLanguage], error) {
					layout := f.native.layout(r.Destination())
					switch name {
					case "foreign-owner":
						layout = ref(gen.NewSpace[gen.LayoutRole[nativeLanguage]](), 1)
					case "absent-record":
						layout = ref(f.native.layouts, 999)
					case "other-destination":
						layout = f.native.layout(f.native.destination(false))
					}
					_, err := r.FreezeLayout(layout)
					if err == nil && name == "freeze-twice" {
						_, err = r.FreezeLayout(layout)
					}
					return gen.Introduction[sourceLanguage, nativeLanguage]{}, err
				},
				apply: original.Apply,
			}
			c := f.configuration(nil)
			c.Plans = append(c.Plans, p)
			s := begin(t, context.Background(), f, c, gen.DefaultLimits())
			expected := error(errRecord)
			if name == "freeze-twice" {
				expected = gen.ErrInvalid
			}
			if out, err := s.Acquire(f.generic); !errors.Is(err, expected) || out.Valid() {
				t.Fatalf("invalid layout admitted: %v", err)
			}
			failedFinish(t, s, expected)
		})
	}
}

func TestNativeHelperDoesNotInheritSourceOrNativeScope(t *testing.T) {
	for _, name := range []string{"source-use", "source-binder", "native-frame"} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t)
			var other callee
			helper := f.plainPlan(true)
			helper.Key, helper.Declaration, helper.Publisher = "helper", declaration{}, nil
			helper.Definer = defineFunc(func(_ context.Context, r gen.DefineRequest[sourceLanguage, nativeLanguage]) (gen.Definition[sourceLanguage, nativeLanguage], error) {
				if r.Site().Environment().Valid() || len(r.Open().Formals) != 0 {
					t.Fatal("native helper inherited source scope")
				}
				switch name {
				case "source-use":
					_, err := r.Use(f.scalar)
					return gen.Definition[sourceLanguage, nativeLanguage]{}, err
				case "source-binder":
					_, err := r.Bound(f.first)
					return gen.Definition[sourceLanguage, nativeLanguage]{}, err
				default:
					return r.CompleteBody(f.native.fragment(other.Frame()))
				}
			})
			c := f.configuration(nil)
			c.Plans = append(c.Plans, f.genericPlan(false, nil), helper)
			s := begin(t, context.Background(), f, c, gen.DefaultLimits())
			var err error
			other, err = s.Acquire(f.generic)
			if err != nil {
				t.Fatal(err)
			}
			owned, err := s.AcquireHelper("helper")
			if err != nil {
				t.Fatal(err)
			}
			if out, err := s.Define(owned); !errors.Is(err, gen.ErrScope) || out.Valid() {
				t.Fatalf("helper captured caller scope: %v", err)
			}
			failedFinish(t, s, gen.ErrScope)
		})
	}
}
