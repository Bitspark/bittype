package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

type People = admission.DataPage[
	admission.PersonSource, model.Person, admission.PersonEvidence,
]
type PeopleEvidence = admission.PageEvidence[
	admission.PersonSource, model.Person, admission.PersonEvidence,
]
type PeoplePages = admission.DataPage[
	admission.PageSource[admission.PersonSource], People, PeopleEvidence,
]
type PeoplePagesEvidence = admission.PageEvidence[
	admission.PageSource[admission.PersonSource], People, PeopleEvidence,
]
type ThreePages = admission.DataPage[
	admission.PageSource[admission.PageSource[admission.PersonSource]],
	PeoplePages, PeoplePagesEvidence,
]

var _ model.Page[model.Page[model.Person]] = PeoplePages{}
var _ model.Page[model.Page[model.Page[model.Person]]] = ThreePages{}
var _ PeoplePages = admission.Nested[
	admission.PersonSource, model.Person, admission.PersonEvidence,
]{}
