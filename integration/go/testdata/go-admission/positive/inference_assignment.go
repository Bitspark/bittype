package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

func pageOf[S, N any, E admission.DataEvidence[S, N]](_ E, item N) admission.DataPage[S, N, E] {
	return model.Page[N]{Items: []N{item}}
}

var _ func(admission.PersonEvidence, model.Person) model.Page[model.Person] = pageOf

type PeopleEvidence = admission.PageEvidence[
	admission.PersonSource, model.Person, admission.PersonEvidence,
]

var _ func(PeopleEvidence, model.Page[model.Person]) model.Page[model.Page[model.Person]] = pageOf
