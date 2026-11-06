package generationfixture

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
)

type association struct {
	captured g.Captured[SourceLanguage, s.Language]
	support  g.Fragment[s.Language]
}

// Rules owns one mapping library's relation. The shared SPI sees only the
// authenticated token and native payload, not the meaning of Support.
type Rules struct {
	space                  g.Space[g.AssociationRole[SourceLanguage, s.Language]]
	records                map[g.Token[SourceLanguage, s.Language]]association
	native                 *Native
	applicationDeclaration g.Declaration[SourceLanguage]
	applicationInline      bool
}

func newRules(native *Native) *Rules {
	return &Rules{space: g.NewSpace[g.AssociationRole[SourceLanguage, s.Language]](), records: make(map[g.Token[SourceLanguage, s.Language]]association), native: native}
}

func (r *Rules) issue(ctx context.Context, want g.Want[SourceLanguage, s.Language], payload, support g.Fragment[s.Language]) (g.Token[SourceLanguage, s.Language], error) {
	if len(r.records) >= want.Site().Limits().References {
		return g.Token[SourceLanguage, s.Language]{}, g.ErrLimit
	}
	if err := r.native.CheckProduct(ctx, want, payload); err != nil {
		return g.Token[SourceLanguage, s.Language]{}, err
	}
	if err := r.native.CheckProduct(ctx, want, support); err != nil {
		return g.Token[SourceLanguage, s.Language]{}, err
	}
	token, _ := r.space.Ref(uint64(len(r.records) + 1))
	r.records[token] = association{captured: g.Captured[SourceLanguage, s.Language]{Binding: want.Binding(), Token: token, Payload: payload}, support: support}
	return token, nil
}

func (r *Rules) Capture(ctx context.Context, want g.Want[SourceLanguage, s.Language], token g.Token[SourceLanguage, s.Language]) (g.Captured[SourceLanguage, s.Language], error) {
	if err := ctx.Err(); err != nil {
		return g.Captured[SourceLanguage, s.Language]{}, err
	}
	record, ok := r.records[token]
	if !ok || !r.space.Owner().Owns(token) {
		return g.Captured[SourceLanguage, s.Language]{}, g.ErrScope
	}
	// Return what this authority actually saved; never overwrite its binding
	// with the requesting Want's binding.
	return record.captured, nil
}

type supportReader struct {
	owners []*Rules
	source *Source
}

func (r supportReader) readActual(ctx context.Context, want g.Want[SourceLanguage, s.Language], actual g.Use[SourceLanguage], product g.Product[SourceLanguage, s.Language]) (g.Fragment[s.Language], error) {
	if !want.Snapshot().Same(r.source.snapshot) {
		return g.Fragment[s.Language]{}, g.ErrScope
	}
	expected, err := r.source.Admit(ctx, want.Use().Environment, actual)
	if err != nil {
		return g.Fragment[s.Language]{}, err
	}
	binding := product.Binding()
	if !product.Valid() || binding.Run != want.Site().Run() || binding.Config != want.Site().Config() || binding.Frame != want.Site().Frame() || binding.Environment != expected.Environment || binding.Subject != expected.Subject {
		return g.Fragment[s.Language]{}, g.ErrScope
	}
	for _, owner := range r.owners {
		if !owner.space.Owner().Owns(product.Token()) {
			continue
		}
		record, ok := owner.records[product.Token()]
		if !ok || record.captured.Binding != product.Binding() || record.captured.Payload != product.Payload() {
			return g.Fragment[s.Language]{}, g.ErrScope
		}
		return record.support, nil
	}
	return g.Fragment[s.Language]{}, g.ErrUnsupported
}

// application owns the relation for this fixture mapping. The candidate cannot
// attach an argument's support to a Page result: the rule independently derives
// both native expressions from its authenticated actual premise before issuance.
func (r *Rules) application(ctx context.Context, want g.Want[SourceLanguage, s.Language], callee g.Callee[SourceLanguage, s.Language], actual g.Use[SourceLanguage], argument g.Product[SourceLanguage, s.Language], payload, support g.Fragment[s.Language], read supportReader) (g.Token[SourceLanguage, s.Language], error) {
	open := callee.Open()
	if !r.applicationDeclaration.Valid() || r.applicationDeclaration != open.Declaration || want.Use().Callee != open.Declaration || len(open.Formals) != 1 || len(want.Use().Actuals) != 1 || want.Use().Actuals[0].Binder != open.Formals[0].Binder || want.Use().Actuals[0].Use != actual {
		return g.Token[SourceLanguage, s.Language]{}, g.ErrScope
	}
	argumentSupport, err := read.readActual(ctx, want, actual, argument)
	if err != nil {
		return g.Token[SourceLanguage, s.Language]{}, err
	}
	var expectedPayload g.Fragment[s.Language]
	if r.applicationInline {
		expectedPayload, err = r.native.slice(argument.Payload())
	} else {
		expectedPayload, err = r.native.reference(callee, argument.Payload(), argumentSupport)
	}
	if err != nil {
		return g.Token[SourceLanguage, s.Language]{}, err
	}
	expectedSupport, err := r.native.slice(argument.Payload())
	if err != nil {
		return g.Token[SourceLanguage, s.Language]{}, err
	}
	if !r.native.sameExpression(payload, expectedPayload) || !r.native.sameExpression(support, expectedSupport) {
		return g.Token[SourceLanguage, s.Language]{}, g.ErrScope
	}
	return r.issue(ctx, want, payload, support)
}

type useFunc func(context.Context, g.UseRequest[SourceLanguage, s.Language]) (g.Product[SourceLanguage, s.Language], error)

func (f useFunc) Use(ctx context.Context, request g.UseRequest[SourceLanguage, s.Language]) (g.Product[SourceLanguage, s.Language], error) {
	return f(ctx, request)
}

type defineFunc func(context.Context, g.DefineRequest[SourceLanguage, s.Language]) (g.Definition[SourceLanguage, s.Language], error)

func (f defineFunc) Define(ctx context.Context, request g.DefineRequest[SourceLanguage, s.Language]) (g.Definition[SourceLanguage, s.Language], error) {
	return f(ctx, request)
}

type publisher struct{}

func (publisher) Publish(_ context.Context, request g.PublishRequest[SourceLanguage, s.Language]) (g.Publication[SourceLanguage, s.Language], error) {
	callee, err := request.Acquire()
	if err != nil {
		return g.Publication[SourceLanguage, s.Language]{}, err
	}
	return request.Expose(callee)
}

type genericMapping struct {
	native              *Native
	rules               *Rules
	read                supportReader
	name                string
	inline              bool
	argumentOnlySupport bool
}

func (m *genericMapping) Introduce(ctx context.Context, request g.IntroduceRequest[SourceLanguage, s.Language]) (g.Introduction[SourceLanguage, s.Language], error) {
	open := request.Open()
	if len(open.Formals) != 1 {
		return g.Introduction[SourceLanguage, s.Language]{}, g.ErrUnsupported
	}
	layout, err := m.native.freeze(request.Site(), request.Destination(), m.name, "S"+m.name)
	if err != nil {
		return g.Introduction[SourceLanguage, s.Language]{}, err
	}
	formals, err := request.FreezeLayout(layout)
	if err != nil {
		return g.Introduction[SourceLanguage, s.Language]{}, err
	}
	want, err := formals.WantFor(open.Formals[0].Binder)
	if err != nil {
		return g.Introduction[SourceLanguage, s.Language]{}, err
	}
	value, err := m.native.formal(layout, 0)
	if err != nil {
		return g.Introduction[SourceLanguage, s.Language]{}, err
	}
	support, err := m.native.formal(layout, 1)
	if err != nil {
		return g.Introduction[SourceLanguage, s.Language]{}, err
	}
	token, err := m.rules.issue(ctx, want, value, support)
	if err != nil {
		return g.Introduction[SourceLanguage, s.Language]{}, err
	}
	return formals.Finish([]g.FormalBinding[SourceLanguage, s.Language]{{Binder: open.Formals[0].Binder, Token: token}})
}

func (m *genericMapping) Apply(ctx context.Context, request g.ApplyRequest[SourceLanguage, s.Language]) (g.Product[SourceLanguage, s.Language], error) {
	formal := request.Callee().Open().Formals[0].Binder
	actual, err := request.Actual(formal)
	if err != nil {
		return g.Product[SourceLanguage, s.Language]{}, err
	}
	argument, err := request.Represent(formal)
	if err != nil {
		return g.Product[SourceLanguage, s.Language]{}, err
	}
	argumentSupport, err := m.read.readActual(ctx, request.Want(), actual, argument)
	if err != nil {
		return g.Product[SourceLanguage, s.Language]{}, err
	}
	var payload g.Fragment[s.Language]
	if m.inline {
		payload, err = m.native.slice(argument.Payload())
	} else {
		payload, err = m.native.reference(request.Callee(), argument.Payload(), argumentSupport)
	}
	if err != nil {
		return g.Product[SourceLanguage, s.Language]{}, err
	}
	// The support for a whole Page<A> is []Rep(A), not A's support. The next
	// enclosing application must consume this newly derived result support.
	support, err := m.native.slice(argument.Payload())
	if err != nil {
		return g.Product[SourceLanguage, s.Language]{}, err
	}
	if m.argumentOnlySupport {
		support = argumentSupport
	}
	token, err := m.rules.application(ctx, request.Want(), request.Callee(), actual, argument, payload, support, m.read)
	if err != nil {
		return g.Product[SourceLanguage, s.Language]{}, err
	}
	return request.Return(token)
}

func (m *genericMapping) Use(_ context.Context, request g.UseRequest[SourceLanguage, s.Language]) (g.Product[SourceLanguage, s.Language], error) {
	return request.Apply()
}

func defineSource(_ context.Context, request g.DefineRequest[SourceLanguage, s.Language]) (g.Definition[SourceLanguage, s.Language], error) {
	body, err := request.Use(request.Open().Body)
	if err != nil {
		return g.Definition[SourceLanguage, s.Language]{}, err
	}
	return request.CompleteBody(body.Payload())
}

func nonGeneric(native *Native, rules *Rules) useFunc {
	return func(ctx context.Context, request g.UseRequest[SourceLanguage, s.Language]) (g.Product[SourceLanguage, s.Language], error) {
		callee, err := request.RequireOwned(request.Want().Use().Callee)
		if err != nil {
			return g.Product[SourceLanguage, s.Language]{}, err
		}
		payload, err := native.reference(callee)
		if err != nil {
			return g.Product[SourceLanguage, s.Language]{}, err
		}
		support, err := native.text()
		if err != nil {
			return g.Product[SourceLanguage, s.Language]{}, err
		}
		token, err := rules.issue(ctx, request.Want(), payload, support)
		if err != nil {
			return g.Product[SourceLanguage, s.Language]{}, err
		}
		return request.Return(token)
	}
}
