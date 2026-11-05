package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

type Forward[S, N any, E admission.DataEvidence[S, N]] = admission.DataPage[S, N, E]

type Nested[S, N any, E admission.DataEvidence[S, N]] = admission.DataPage[
	admission.PageSource[S], model.Page[N], admission.PageEvidence[S, N, E],
]

func forward[S, N any, E admission.DataEvidence[S, N]](
	value admission.DataPage[S, N, E],
) model.Page[N] {
	return value
}

var _ func(model.Page[model.Person]) model.Page[model.Person] = forward[
	admission.PersonSource, model.Person, admission.PersonEvidence,
]
var _ model.Page[model.Page[model.Person]] = Nested[
	admission.PersonSource, model.Person, admission.PersonEvidence,
]{}
