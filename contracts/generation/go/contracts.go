package generation

import (
	"context"
	"errors"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

var (
	ErrInvalid = errors.New("invalid generation operation")
	ErrUnsupported = errors.New("generation capability unavailable")
	ErrScope = errors.New("incompatible generation scope or association")
	ErrLimit = errors.New("generation limit exhausted")
	ErrIncomplete = errors.New("generation obligations incomplete")
)

// Limits bound host work and target-controlled expansion. Providers must also
// bound their own captured records and honor cancellation. They are trusted Go
// collaborators, not a sandbox. Zero/negative limits are invalid.
type Limits struct {
	Operations int
	Depth int
	Plans int
	References int
	Formals int
	Reservations int
	OutputNodes int
	OutputBytes int
}
func DefaultLimits() Limits { return Limits{100000, 128, 4096, 65536, 4096, 4096, 200000, 32 << 20} }
func (l Limits) Valid() bool {
	return l.Operations > 0 && l.Depth > 0 && l.Plans > 0 && l.References > 0 && l.Formals > 0 && l.Reservations > 0 && l.OutputNodes > 0 && l.OutputBytes > 0
}

// Binding describes lexical association. The shorter callback lifetime is on
// Want and Product, not on this record. It contains immutable handles only.
type Binding[F any] struct {
	Run RunID
	Config ConfigID
	Issuer IssuerID
	Subject Use[F]
	Environment Environment[F]
	Frame FrameID
	Signature SignatureID
}
type Captured[F, T any] struct {
	Binding Binding[F]
	Token Token[F, T]
	Payload Fragment[T]
}

// Associations authenticates its immutable record and returns the payload and
// binding from that record. It must not copy Want's expected binding onto an
// unrelated token. There is no public operation accepting Captured directly.
type Associations[F, T any] interface {
	Capture(context.Context, Want[F, T], Token[F, T]) (Captured[F, T], error)
}

// Native owns record authentication, grammar, lexical visibility, symbols,
// complete headers and assembly. Every Check authenticates existing records in
// the configured target stores; nonzero refs or claimed free-symbol lists alone
// are insufficient. Assemble materializes fresh containment at every use and
// applies the saved layout to each body. It must enforce expansion limits before
// allocating nodes. External exports are not fabricated AST NodeRefs.
type Native[F, T any] interface {
	CheckDestination(context.Context, Site[F, T], Destination[T], bool) error
	CheckLayout(context.Context, Site[F, T], Destination[T], Layout[T]) error
	CheckProduct(context.Context, Want[F, T], Fragment[T]) error
	CheckBody(context.Context, Callee[F, T], Fragment[T]) error
	CheckExposure(context.Context, Site[F, T], Declaration[F], Callee[F, T]) error
	Assemble(context.Context, []Publication[F, T], []Definition[F, T], Limits) (representation.ASTSet, error)
}

// The independent capability interfaces are Syntax-generation operations.
// bittype:dimension Syntax
type User[F, T any] interface { Use(context.Context, UseRequest[F, T]) (Product[F, T], error) }
// bittype:dimension Syntax
type Publisher[F, T any] interface { Publish(context.Context, PublishRequest[F, T]) (Publication[F, T], error) }
// bittype:dimension Syntax
type Definer[F, T any] interface { Define(context.Context, DefineRequest[F, T]) (Definition[F, T], error) }
// bittype:dimension Syntax
type Convention[F, T any] interface {
	Introduce(context.Context, IntroduceRequest[F, T]) (Introduction[F, T], error)
	Apply(context.Context, ApplyRequest[F, T]) (Product[F, T], error)
}

// Plan is preselected configuration. Key is local to this captured configuration,
// never semantic declaration identity. A zero Declaration denotes a direct
// mapping or native-only helper. Helpers have no inherited source environment.
// Callbacks and their private records must retain stable policy for the run;
// copying an interface cannot freeze arbitrary implementation state.
type Plan[F, T any] struct {
	Key string
	Declaration Declaration[F]
	User User[F, T]
	Publisher Publisher[F, T]
	Definer Definer[F, T]
	Convention Convention[F, T]
	Associations Associations[F, T]
	Acquire func(context.Context, Site[F, T]) (Destination[T], bool, error)
	// Header supplies the layout without a source generic convention. Source
	// binders must be empty; a native-only helper may have explicit native formals.
	Header func(context.Context, Site[F, T], Destination[T]) (Layout[T], error)
}

type Config[F, T any] struct {
	Source Source[F]
	Native Native[F, T]
	Target representation.Language
	Evidence []semantic.EvidenceRequirement
	Plans []Plan[F, T]
	// Select chooses among the already supplied plans. It is not discovery,
	// fallback, analysis or a core constructor/name switch.
	Select func(UseData[F]) (string, error)
}
type Input struct { Snapshot semantic.Snapshot; Evidence []semantic.Result }

type UseRequest[F, T any] interface {
	Want() Want[F, T]
	Child(Use[F]) (Product[F, T], error)
	Bound() (Product[F, T], error)
	Apply() (Product[F, T], error)
	RequireOwned(Declaration[F]) (Callee[F, T], error)
	RequireHelper(string) (Callee[F, T], error)
	Return(Token[F, T]) (Product[F, T], error)
}
type IntroduceRequest[F, T any] interface {
	Site() Site[F, T]
	Open() OpenData[F]
	Destination() Destination[T]
	FreezeLayout(Layout[T]) (FormalRequest[F, T], error)
}
type FormalBinding[F, T any] struct { Binder Binder[F]; Token Token[F, T] }
type FormalRequest[F, T any] interface {
	WantFor(Binder[F]) (Want[F, T], error)
	Finish([]FormalBinding[F, T]) (Introduction[F, T], error)
}
type ApplyRequest[F, T any] interface {
	Want() Want[F, T]
	Callee() Callee[F, T]
	Actual(Binder[F]) (Use[F], error)
	Represent(Binder[F]) (Product[F, T], error)
	Return(Token[F, T]) (Product[F, T], error)
}
type PublishRequest[F, T any] interface {
	Site() Site[F, T]
	Root() Declaration[F]
	Acquire() (Callee[F, T], error)
	Expose(Callee[F, T]) (Publication[F, T], error)
}
type DefineRequest[F, T any] interface {
	Site() Site[F, T]
	Callee() Callee[F, T]
	Open() OpenData[F]
	Bound(Binder[F]) (Product[F, T], error)
	Use(Use[F]) (Product[F, T], error)
	CompleteBody(Fragment[T]) (Definition[F, T], error)
}
