// Package admission is an independently authored compiler fixture, not an
// implementation of Horos analysis, Logos proofs, or Bittype generation.
package admission

import "admission.example/model"

type PersonSource struct{}
type WireSource struct{}
type PageSource[S any] struct{}

type DataEvidence[S, N any] interface {
	~func(S) N
	dataEvidence()
}

// This fixture takes the Person membership/representation judgment as given.
// No Go compilation result establishes that source judgment.
type PersonEvidence func(PersonSource) model.Person

func (PersonEvidence) dataEvidence() {}

type PageEvidence[S, N any, E DataEvidence[S, N]] func(PageSource[S]) model.Page[N]

func (PageEvidence[S, N, E]) dataEvidence() {}

type DataPage[S, N any, E DataEvidence[S, N]] = model.Page[N]

type Nested[S, N any, E DataEvidence[S, N]] = DataPage[
	PageSource[S], model.Page[N], PageEvidence[S, N, E],
]

type Ignore[S, N any, E DataEvidence[S, N]] = string
