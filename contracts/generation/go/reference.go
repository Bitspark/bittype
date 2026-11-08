// Package generation defines type-owned generation operations and their guarded,
// explicitly driven invocation. It interprets no source or native type rules.
package generation

// Space is an issuer's address space. Keep it with the owner of immutable,
// append-only records. Ref allocates no record and establishes no admission.
type Space[K any] struct{ owner *domain[K] }
type domain[K any] struct{ identity byte }

// Ref is a typed, owner-qualified address, not a proof or an existing-record claim.
type Ref[K any] struct {
	owner *domain[K]
	slot  uint64
}

// Owner exposes identity checks without exposing the issuing space.
type Owner[K any] struct{ owner *domain[K] }

func NewSpace[K any]() Space[K]    { return Space[K]{&domain[K]{identity: 1}} }
func (s Space[K]) Owner() Owner[K] { return Owner[K]{s.owner} }
func (s Space[K]) Ref(slot uint64) (Ref[K], bool) {
	if s.owner == nil || slot == 0 {
		return Ref[K]{}, false
	}
	return Ref[K]{s.owner, slot}, true
}
func (s Space[K]) Index(r Ref[K]) (uint64, bool) {
	if !s.Owner().Owns(r) {
		return 0, false
	}
	return r.slot, true
}
func (o Owner[K]) Valid() bool        { return o.owner != nil }
func (o Owner[K]) Owns(r Ref[K]) bool { return o.owner != nil && r.owner == o.owner && r.slot != 0 }
func (r Ref[K]) Valid() bool          { return r.owner != nil && r.slot != 0 }

type DeclarationRole[F any] struct{}
type BinderRole[F any] struct{}
type UseRole[F any] struct{}
type EnvironmentRole[F any] struct{}
type Declaration[F any] = Ref[DeclarationRole[F]]
type Binder[F any] = Ref[BinderRole[F]]
type Use[F any] = Ref[UseRole[F]]
type Environment[F any] = Ref[EnvironmentRole[F]]
type identity struct{ value byte }

func fresh() *identity { return &identity{value: 1} }
