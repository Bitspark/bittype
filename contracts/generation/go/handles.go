package generation

import "github.com/Bitspark/bittype/contracts/semantic/go"

type operation[F, T any] struct {
	run *Session[F, T]
	plan *planState[F, T]
	ticket Ticket
	frame FrameID
	environment Environment[F]
	callee *calleeState[F, T]
	live, completed bool
}

// Site carries host identity. Its accessors remain inspectable after a callback;
// the callback's issuance methods expire when it returns.
type Site[F, T any] struct { operation *operation[F, T] }
func (s Site[F, T]) Valid() bool { return s.operation != nil }
func (s Site[F, T]) Run() RunID { if !s.Valid() { return RunID{} }; return s.operation.run.id }
func (s Site[F, T]) Config() ConfigID { if !s.Valid() { return ConfigID{} }; return s.operation.run.configuration }
func (s Site[F, T]) Issuer() IssuerID { if !s.Valid() { return IssuerID{} }; return s.operation.plan.issuer }
func (s Site[F, T]) Frame() FrameID { if !s.Valid() { return FrameID{} }; return s.operation.frame }
func (s Site[F, T]) Signature() SignatureID { if !s.Valid() || s.operation.callee == nil { return SignatureID{} }; return s.operation.callee.signature }
func (s Site[F, T]) Ticket() Ticket { if !s.Valid() { return Ticket{} }; return s.operation.ticket }
func (s Site[F, T]) Snapshot() semantic.Snapshot { if !s.Valid() { return semantic.Snapshot{} }; return s.operation.run.snapshot }
func (s Site[F, T]) Limits() Limits { if !s.Valid() { return Limits{} }; return s.operation.run.limits }
func (s Site[F, T]) Environment() Environment[F] { if !s.Valid() { return Environment[F]{} }; return s.operation.environment }

type Want[F, T any] struct { site Site[F, T]; use UseData[F] }
func (w Want[F, T]) Valid() bool { return w.site.Valid() && w.use.Subject.Valid() }
func (w Want[F, T]) Site() Site[F, T] { return w.site }
func (w Want[F, T]) Snapshot() semantic.Snapshot { return w.site.Snapshot() }
func (w Want[F, T]) Use() UseData[F] { return w.use.Clone() }
func (w Want[F, T]) Callee() Callee[F, T] { if !w.site.Valid() { return Callee[F, T]{} }; return Callee[F, T]{w.site.operation.callee} }
func (w Want[F, T]) Binding() Binding[F] {
	return Binding[F]{Run: w.site.Run(), Config: w.site.Config(), Issuer: w.site.Issuer(), Subject: w.use.Subject, Environment: w.use.Environment, Frame: w.site.Frame(), Signature: w.site.Signature()}
}

type productState[F, T any] struct { captured Captured[F, T]; ticket Ticket }
type Product[F, T any] struct { state *productState[F, T] }
func (p Product[F, T]) Valid() bool { return p.state != nil }
func (p Product[F, T]) Binding() Binding[F] { if !p.Valid() { return Binding[F]{} }; return p.state.captured.Binding }
func (p Product[F, T]) Payload() Fragment[T] { if !p.Valid() { return Fragment[T]{} }; return p.state.captured.Payload }
func (p Product[F, T]) Token() Token[F, T] { if !p.Valid() { return Token[F, T]{} }; return p.state.captured.Token }

type phase uint8
const (reserved phase = iota; introducing; ready; defining; defined)
type calleeState[F, T any] struct {
	run *Session[F, T]
	plan *planState[F, T]
	open OpenData[F]
	frame FrameID
	signature SignatureID
	destination Destination[T]
	layout Layout[T]
	owned bool
	phase phase
	formals map[Binder[F]]Product[F, T]
}
type Callee[F, T any] struct { state *calleeState[F, T] }
func (c Callee[F, T]) Valid() bool { return c.state != nil }
func (c Callee[F, T]) Open() OpenData[F] { if !c.Valid() { return OpenData[F]{} }; return c.state.open.Clone() }
func (c Callee[F, T]) Destination() Destination[T] { if !c.Valid() { return Destination[T]{} }; return c.state.destination }
func (c Callee[F, T]) Layout() Layout[T] { if !c.Valid() { return Layout[T]{} }; return c.state.layout }
func (c Callee[F, T]) Owned() bool { return c.Valid() && c.state.owned }
func (c Callee[F, T]) Frame() FrameID { if !c.Valid() { return FrameID{} }; return c.state.frame }
func (c Callee[F, T]) Signature() SignatureID { if !c.Valid() { return SignatureID{} }; return c.state.signature }
func (c Callee[F, T]) Key() string { if !c.Valid() { return "" }; return c.state.plan.value.Key }

type introductionState[F, T any] struct { callee *calleeState[F, T]; ticket Ticket; formals map[Binder[F]]Product[F, T] }
type Introduction[F, T any] struct { state *introductionState[F, T] }
func (i Introduction[F, T]) Valid() bool { return i.state != nil }

type publicationState[F, T any] struct { callee Callee[F, T]; root Declaration[F]; ticket Ticket }
type Publication[F, T any] struct { state *publicationState[F, T] }
func (p Publication[F, T]) Valid() bool { return p.state != nil }
func (p Publication[F, T]) Callee() Callee[F, T] { if !p.Valid() { return Callee[F, T]{} }; return p.state.callee }
func (p Publication[F, T]) Root() Declaration[F] { if !p.Valid() { return Declaration[F]{} }; return p.state.root }

type definitionState[F, T any] struct { callee Callee[F, T]; body Fragment[T]; ticket Ticket }
type Definition[F, T any] struct { state *definitionState[F, T] }
func (d Definition[F, T]) Valid() bool { return d.state != nil }
func (d Definition[F, T]) Callee() Callee[F, T] { if !d.Valid() { return Callee[F, T]{} }; return d.state.callee }
func (d Definition[F, T]) Body() Fragment[T] { if !d.Valid() { return Fragment[T]{} }; return d.state.body }
