package generation

import (
	"context"
	"slices"
	"github.com/Bitspark/bittype/contracts/semantic/go"
)

type Formal[F any] struct { Binder Binder[F]; Subject Use[F] }
type Actual[F any] struct { Binder Binder[F]; Use Use[F] }

// These are provider construction records. Only the configured source can admit
// them; writing a literal is not a route to a host-owned product or request.
type OpenData[F any] struct {
	Declaration Declaration[F]
	Environment Environment[F]
	Body Use[F]
	Formals []Formal[F]
}
func (o OpenData[F]) Clone() OpenData[F] { o.Formals = slices.Clone(o.Formals); return o }

type UseData[F any] struct {
	Use Use[F]
	Subject Use[F]
	Environment Environment[F]
	Callee Declaration[F]
	Bound Binder[F]
	Actuals []Actual[F]
}
func (u UseData[F]) Clone() UseData[F] { u.Actuals = slices.Clone(u.Actuals); return u }

type SourceOwners[F any] struct {
	Declarations Owner[DeclarationRole[F]]
	Binders Owner[BinderRole[F]]
	Uses Owner[UseRole[F]]
	Environments Owner[EnvironmentRole[F]]
}
func (o SourceOwners[F]) Valid() bool {
	return o.Declarations.Valid() && o.Binders.Valid() && o.Uses.Valid() && o.Environments.Valid()
}

// References requests authentication of existing, correctly typed source
// records. An owner match alone does not establish existence or applicability.
type References[F any] struct {
	Declarations []Declaration[F]
	Binders []Binder[F]
	Uses []Use[F]
	Environments []Environment[F]
}

// Source consumes already supplied receipts/queries. It performs no hidden
// resolution or analysis. Its records and policy remain immutable during a run.
type Source[F any] interface {
	CheckInput(context.Context, semantic.Snapshot, []semantic.Result) error
	Owners() SourceOwners[F]
	RootEnvironment() Environment[F]
	Authenticate(context.Context, References[F]) error
	Open(context.Context, Declaration[F]) (OpenData[F], error)
	// Admit establishes the whole use's applicability under its original
	// assumptions using supplied evidence, including actuals the mapping may
	// never represent. Authenticate only checks record ownership/existence.
	// Use identifies the occurrence; Subject is the source-certified effective
	// subject. Missing evidence must fail, never trigger hidden analysis.
	Admit(context.Context, Environment[F], Use[F]) (UseData[F], error)
}
