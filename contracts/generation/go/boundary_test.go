package generation_test

import (
	"context"
	"errors"
	"testing"

	gen "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func TestSourceRetainsOriginalSnapshotWithoutEvidence(t *testing.T) {
	f := newFixture(t)
	fresh, err := semantic.Capture(context.Background(), f.checker, f.input)
	if err != nil { t.Fatal(err) }
	if fresh.Same(f.source.snapshot) { t.Fatal("fixture recapture reused an identity") }
	if s, err := gen.Begin[sourceLanguage, nativeLanguage](context.Background(), gen.Input{Snapshot: fresh}, f.configuration(nil), gen.DefaultLimits()); !errors.Is(err, errRecord) || s != nil {
		t.Fatalf("source admitted a different capture of identical ASTs: %v", err)
	}
	called := false
	c := f.configuration(func(_ context.Context, r gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error) {
		called = true
		if !r.Want().Snapshot().Same(f.source.snapshot) || !r.Want().Site().Snapshot().Same(f.source.snapshot) {
			t.Fatal("host recaptured or lost the source snapshot")
		}
		return r.Return(f.associations.issue(r.Want(), f.closed, nil))
	})
	s := begin(t, context.Background(), f, c, gen.DefaultLimits())
	if _, err := s.Use(f.scalar); err != nil || !called { t.Fatal(err) }
}

func TestSourceRejectsForeignAndAbsentRecordsBeforeMapping(t *testing.T) {
	for _, name := range []string{"foreign-use", "absent-use", "foreign-subject", "absent-subject", "foreign-binder", "absent-binder", "foreign-declaration", "absent-declaration", "foreign-environment", "absent-environment"} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t)
			query, expected := f.scalar, error(gen.ErrScope)
			switch name {
			case "foreign-use": query = ref(gen.NewSpace[gen.UseRole[sourceLanguage]](), 1)
			case "absent-use": query, expected = ref(f.source.uses, 999), errRecord
			case "foreign-subject": f.source.admitHook = func(u useData) useData { u.Subject = ref(gen.NewSpace[gen.UseRole[sourceLanguage]](), 1); return u }
			case "absent-subject":
				expected = errRecord
				f.source.admitHook = func(u useData) useData { u.Subject = ref(f.source.uses, 999); return u }
			case "foreign-binder": f.source.admitHook = func(u useData) useData { u.Bound = ref(gen.NewSpace[gen.BinderRole[sourceLanguage]](), 1); return u }
			case "absent-binder":
				expected = errRecord
				f.source.admitHook = func(u useData) useData { u.Bound = ref(f.source.binders, 999); return u }
			case "foreign-declaration": f.source.admitHook = func(u useData) useData { u.Callee = ref(gen.NewSpace[gen.DeclarationRole[sourceLanguage]](), 1); return u }
			case "absent-declaration":
				expected = errRecord
				f.source.admitHook = func(u useData) useData { u.Callee = ref(f.source.declarations, 999); return u }
			case "foreign-environment":
				query = f.application
				f.source.openHook = func(o openData) openData { o.Environment = ref(gen.NewSpace[gen.EnvironmentRole[sourceLanguage]](), 1); return o }
			case "absent-environment":
				query, expected = f.application, errRecord
				f.source.openHook = func(o openData) openData { o.Environment = ref(f.source.environments, 999); return o }
			}
			called := false
			c := f.configuration(func(context.Context, gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error) { called = true; return product{}, nil })
			s := begin(t, context.Background(), f, c, gen.DefaultLimits())
			out, err := s.Use(query)
			if !errors.Is(err, expected) || out.Valid() || called { t.Fatalf("unauthenticated source reached mapping: %v", err) }
			failedFinish(t, s, expected)
			if f.native.assemblies != 0 { t.Fatal("failed source admission reached assembly") }
		})
	}
}

func TestSourceOccurrenceAndActualBinderCoverage(t *testing.T) {
	for _, name := range []string{"wrong-occurrence", "wrong-environment", "missing", "duplicate", "different-binder", "bound-and-callee", "arguments-without-callee"} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t)
			expected := error(gen.ErrInvalid)
			if name == "different-binder" { expected = gen.ErrScope }
			f.source.admitHook = func(u useData) useData {
				switch name {
				case "wrong-occurrence": u.Use = f.alternate
				case "wrong-environment": u.Environment = f.genericEnv
				case "missing": u.Actuals = u.Actuals[:1]
				case "duplicate": u.Actuals[1].Binder = u.Actuals[0].Binder
				case "different-binder": u.Actuals[1].Binder = f.spare
				case "bound-and-callee": u.Bound = f.first
				case "arguments-without-callee": u.Callee = declaration{}
				}
				return u
			}
			s := begin(t, context.Background(), f, f.configuration(nil), gen.DefaultLimits())
			if p, err := s.Use(f.application); !errors.Is(err, expected) || p.Valid() { t.Fatalf("invalid use admitted: %v", err) }
			if f.associations.count != 0 { t.Fatal("invalid application reached a mapping authority") }
			failedFinish(t, s, expected)
		})
	}
}

func TestIgnoredActualStillRequiresWholeUseAdmission(t *testing.T) {
	f := newFixture(t)
	f.source.invalidUses[f.scalar] = true
	called := false
	c := f.configuration(func(_ context.Context, r gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error) {
		called = true
		// This mapping intentionally never represents either actual.
		return r.Return(f.associations.issue(r.Want(), f.closed, nil))
	})
	s := begin(t, context.Background(), f, c, gen.DefaultLimits())
	if p, err := s.Use(f.application); !errors.Is(err, errAdmission) || p.Valid() || called { t.Fatalf("ignored invalid actual bypassed admission: %v", err) }
	failedFinish(t, s, errAdmission)
}

func TestSourceAndRequestConstructionSlicesDoNotAlias(t *testing.T) {
	f := newFixture(t)
	// Deliberately return source-owned construction storage to exercise the
	// host's copy, independently of the ordinary fixture's Clone implementation.
	f.source.admitHook = func(useData) useData { return f.source.useRecords[f.application] }
	c := f.configuration(func(_ context.Context, r gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error) {
		value := r.Want().Use()
		if value.Actuals[0].Use != f.scalar { t.Fatal("selector mutated the admitted source use") }
		value.Actuals[0].Use = f.alternate
		if r.Want().Use().Actuals[0].Use != f.scalar { t.Fatal("Want accessor exposed retained actual storage") }
		return r.Return(f.associations.issue(r.Want(), f.closed, nil))
	})
	c.Select = func(u useData) (string, error) { u.Actuals[0].Use = f.alternate; return "direct", nil }
	s := begin(t, context.Background(), f, c, gen.DefaultLimits())
	if _, err := s.Use(f.application); err != nil { t.Fatal(err) }
	if f.source.useRecords[f.application].Actuals[0].Use != f.scalar { t.Fatal("source storage was mutated") }
}

func TestProductsCannotReplayCallbackTickets(t *testing.T) {
	for _, sealCurrent := range []bool{false, true} {
		t.Run(map[bool]string{false: "without-current-return", true: "after-current-return"}[sealCurrent], func(t *testing.T) {
			f := newFixture(t)
			var previous product
			c := f.configuration(func(_ context.Context, r gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error) {
				if !previous.Valid() { return r.Return(f.associations.issue(r.Want(), f.closed, nil)) }
				if sealCurrent {
					if _, err := r.Return(f.associations.issue(r.Want(), f.closed, nil)); err != nil { return product{}, err }
				}
				return previous, nil
			})
			s := begin(t, context.Background(), f, c, gen.DefaultLimits())
			var err error
			previous, err = s.Use(f.scalar)
			if err != nil { t.Fatal(err) }
			if p, err := s.Use(f.scalar); !errors.Is(err, gen.ErrInvalid) || p.Valid() { t.Fatalf("replayed product accepted: %v", err) }
			failedFinish(t, s, gen.ErrInvalid)
		})
	}
}

func TestSameAssociationCanBeAcceptedInANewCallback(t *testing.T) {
	f := newFixture(t)
	var saved token
	c := f.configuration(func(_ context.Context, r gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error) {
		if !saved.Valid() { saved = f.associations.issue(r.Want(), f.closed, nil) }
		return r.Return(saved)
	})
	s := begin(t, context.Background(), f, c, gen.DefaultLimits())
	for range 2 {
		if _, err := s.Use(f.scalar); err != nil { t.Fatal(err) }
	}
	if f.associations.count != 2 { t.Fatal("reused token skipped its configured authority") }
}

func TestRetainedRequestCannotIssueAfterCallback(t *testing.T) {
	f := newFixture(t)
	var retained gen.UseRequest[sourceLanguage, nativeLanguage]
	var saved token
	c := f.configuration(func(_ context.Context, r gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error) {
		retained = r
		saved = f.associations.issue(r.Want(), f.closed, nil)
		return r.Return(saved)
	})
	s := begin(t, context.Background(), f, c, gen.DefaultLimits())
	if _, err := s.Use(f.scalar); err != nil { t.Fatal(err) }
	if !retained.Want().Snapshot().Same(f.source.snapshot) { t.Fatal("read-only inspection expired with issuance") }
	if p, err := retained.Return(saved); !errors.Is(err, gen.ErrInvalid) || p.Valid() { t.Fatalf("late callback issuance accepted: %v", err) }
	if f.associations.count != 1 { t.Fatal("expired request reached authority") }
	failedFinish(t, s, gen.ErrInvalid)
}

func TestAssociationBindingMustMatchTheWholeRequest(t *testing.T) {
	for _, name := range []string{"subject", "environment", "run", "configuration", "issuer", "frame-with-closed-payload", "signature", "token"} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t)
			var other callee
			c := f.configuration(func(_ context.Context, r gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error) {
				token := f.associations.issue(r.Want(), f.closed, func(c *captured) {
					switch name {
					case "subject": c.Binding.Subject = f.alternate
					case "environment": c.Binding.Environment = f.genericEnv
					case "run": c.Binding.Run = gen.RunID{}
					case "configuration": c.Binding.Config = gen.ConfigID{}
					case "issuer": c.Binding.Issuer = gen.IssuerID{}
					case "frame-with-closed-payload": c.Binding.Frame = other.Frame()
					case "signature": c.Binding.Signature = other.Signature()
					case "token": c.Token = ref(f.associations.space, 999)
					}
				})
				return r.Return(token)
			})
			c.Plans = append(c.Plans, f.genericPlan(false, nil))
			s := begin(t, context.Background(), f, c, gen.DefaultLimits())
			var err error
			other, err = s.Acquire(f.generic)
			if err != nil { t.Fatal(err) }
			if p, err := s.Use(f.scalar); !errors.Is(err, gen.ErrScope) || p.Valid() { t.Fatalf("mismatched support accepted: %v", err) }
			failedFinish(t, s, gen.ErrScope)
		})
	}
}

func TestAssociationAndNativeRecordsAreAuthenticated(t *testing.T) {
	for _, name := range []string{"foreign-token", "absent-token", "foreign-fragment", "absent-fragment", "native-foreign-frame"} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t)
			var other callee
			expected := error(errRecord)
			if name == "native-foreign-frame" { expected = gen.ErrScope }
			c := f.configuration(func(_ context.Context, r gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error) {
				switch name {
				case "foreign-token": return r.Return(newAssociations().issue(r.Want(), f.closed, nil))
				case "absent-token": return r.Return(ref(f.associations.space, 999))
				}
				payload := f.closed
				switch name {
				case "foreign-fragment": payload = ref(gen.NewSpace[gen.FragmentRole[nativeLanguage]](), 1)
				case "absent-fragment": payload = ref(f.native.fragments, 999)
				case "native-foreign-frame": payload = f.native.fragment(other.Frame())
				}
				return r.Return(f.associations.issue(r.Want(), payload, nil))
			})
			c.Plans = append(c.Plans, f.genericPlan(false, nil))
			s := begin(t, context.Background(), f, c, gen.DefaultLimits())
			var err error
			other, err = s.Acquire(f.generic)
			if err != nil { t.Fatal(err) }
			if p, err := s.Use(f.scalar); !errors.Is(err, expected) || p.Valid() { t.Fatalf("unauthenticated product accepted: %v", err) }
			failedFinish(t, s, expected)
		})
	}
}

func TestIntroductionCoversExactSourceBinders(t *testing.T) {
	for _, name := range []string{"missing", "duplicate", "different-binder", "wrong-support", "reordered"} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t)
			alter := func(bindings []formalBinding) []formalBinding {
				switch name {
				case "missing": return bindings[:1]
				case "duplicate": bindings[1] = bindings[0]
				case "different-binder": bindings[1].Binder = f.spare
				case "wrong-support": bindings[1].Token = bindings[0].Token
				case "reordered": bindings[0], bindings[1] = bindings[1], bindings[0]
				}
				return bindings
			}
			c := f.configuration(nil)
			c.Plans = append(c.Plans, f.genericPlan(false, alter))
			s := begin(t, context.Background(), f, c, gen.DefaultLimits())
			out, err := s.Acquire(f.generic)
			if name == "reordered" {
				if err != nil || !out.Valid() { t.Fatalf("binder identity was treated as position: %v", err) }
				return
			}
			expected := error(gen.ErrInvalid)
			if name == "different-binder" || name == "wrong-support" { expected = gen.ErrScope }
			if !errors.Is(err, expected) || out.Valid() { t.Fatalf("invalid formal coverage accepted: %v", err) }
			failedFinish(t, s, expected)
		})
	}
}

func TestRetainedFormalRequestCannotIssueAfterIntroduction(t *testing.T) {
	f := newFixture(t)
	var retained gen.FormalRequest[sourceLanguage, nativeLanguage]
	p := f.genericPlan(false, nil)
	original := p.Convention
	p.Convention = convention{
		introduce: func(_ context.Context, r gen.IntroduceRequest[sourceLanguage, nativeLanguage]) (gen.Introduction[sourceLanguage, nativeLanguage], error) {
			var err error
			retained, err = r.FreezeLayout(f.native.layout(r.Destination()))
			if err != nil { return gen.Introduction[sourceLanguage, nativeLanguage]{}, err }
			var bindings []formalBinding
			for _, formal := range r.Open().Formals {
				w, err := retained.WantFor(formal.Binder)
				if err != nil { return gen.Introduction[sourceLanguage, nativeLanguage]{}, err }
				bindings = append(bindings, formalBinding{Binder: formal.Binder, Token: f.associations.issue(w, f.closed, nil)})
			}
			return retained.Finish(bindings)
		},
		apply: original.Apply,
	}
	c := f.configuration(nil)
	c.Plans = append(c.Plans, p)
	s := begin(t, context.Background(), f, c, gen.DefaultLimits())
	if _, err := s.Acquire(f.generic); err != nil { t.Fatal(err) }
	if w, err := retained.WantFor(f.first); !errors.Is(err, gen.ErrInvalid) || w.Valid() { t.Fatalf("formal request remained live: %v", err) }
	failedFinish(t, s, gen.ErrInvalid)
}

func TestCancellationBeforeAndAfterCollaborators(t *testing.T) {
	for _, phase := range []string{"before-begin", "source-input", "mapping", "assembly"} {
		t.Run(phase, func(t *testing.T) {
			f := newFixture(t)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := f.configuration(nil)
			switch phase {
			case "before-begin": cancel()
			case "source-input": f.source.checkHook = cancel
			case "mapping": c = f.configuration(func(_ context.Context, r gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error) {
				p, err := r.Return(f.associations.issue(r.Want(), f.closed, nil))
				cancel()
				return p, err
			})
			case "assembly": f.native.assemblyHook = func() error { cancel(); return nil }
			}
			s, err := gen.Begin[sourceLanguage, nativeLanguage](ctx, gen.Input{Snapshot: f.source.snapshot}, c, gen.DefaultLimits())
			if phase == "before-begin" || phase == "source-input" {
				if !errors.Is(err, context.Canceled) || s != nil { t.Fatalf("canceled input committed: %v", err) }
				if phase == "before-begin" && f.source.checks != 0 { t.Fatal("canceled input reached source") }
				return
			}
			if err != nil { t.Fatal(err) }
			p, err := s.Use(f.scalar)
			if phase == "mapping" {
				if !errors.Is(err, context.Canceled) || p.Valid() { t.Fatalf("canceled mapping committed: %v", err) }
			} else if err != nil { t.Fatal(err) }
			failedFinish(t, s, context.Canceled)
		})
	}
}

func TestHostLimitsAbortTheInvocation(t *testing.T) {
	for _, name := range []string{"operations", "references", "depth", "output-nodes", "output-bytes"} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t)
			limits := gen.DefaultLimits()
			c := f.configuration(nil)
			switch name {
			case "operations": limits.Operations = 1
			case "references": limits.References = 1
			case "depth":
				limits.Depth = 1
				c = f.configuration(func(_ context.Context, r gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error) { return r.Child(f.scalar) })
			case "output-nodes": limits.OutputNodes = 1
			case "output-bytes": limits.OutputBytes = 1
			}
			s := begin(t, context.Background(), f, c, limits)
			p, err := s.Use(f.scalar)
			if name == "output-nodes" || name == "output-bytes" {
				if err != nil || !p.Valid() { t.Fatal(err) }
			} else if !errors.Is(err, gen.ErrLimit) || p.Valid() { t.Fatalf("work limit did not abort: %v", err) }
			failedFinish(t, s, gen.ErrLimit)
			assemblies := f.native.assemblies
			failedFinish(t, s, gen.ErrLimit)
			if f.native.assemblies != assemblies { t.Fatal("failed Finish retried assembly") }
		})
	}
}

func TestFailedFinishIsTerminalAndCannotRepairMissingDefinitions(t *testing.T) {
	f := newFixture(t)
	c := f.configuration(nil)
	c.Plans = append(c.Plans, f.genericPlan(true, nil))
	s := begin(t, context.Background(), f, c, gen.DefaultLimits())
	owned, err := s.Acquire(f.generic)
	if err != nil { t.Fatal(err) }
	failedFinish(t, s, gen.ErrIncomplete)
	if d, err := s.Define(owned); !errors.Is(err, gen.ErrIncomplete) || d.Valid() { t.Fatalf("incomplete invocation was reopened: %v", err) }
	failedFinish(t, s, gen.ErrIncomplete)
	if f.native.assemblies != 0 { t.Fatal("unfilled reservation reached assembly") }
}

func TestAssemblyFailureAndSuccessfulFinishAreBothTerminal(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "failure"}[fail], func(t *testing.T) {
			f := newFixture(t)
			if fail { f.native.assemblyHook = func() error { return errRecord } }
			s := begin(t, context.Background(), f, f.configuration(nil), gen.DefaultLimits())
			if fail {
				failedFinish(t, s, errRecord)
				failedFinish(t, s, errRecord)
			} else {
				out, err := s.Finish()
				if err != nil || len(out.Trees) == 0 { t.Fatal(err) }
				again, err := s.Finish()
				if !errors.Is(err, gen.ErrInvalid) || len(again.Trees) != 0 || again.Entry != (ast.NodeRef{}) { t.Fatalf("completed invocation published twice: %v", err) }
			}
			if f.native.assemblies != 1 { t.Fatal("terminal invocation retried assembly") }
		})
	}
}
