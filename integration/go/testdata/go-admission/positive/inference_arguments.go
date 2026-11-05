package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

func pageOf[S, N any, E admission.DataEvidence[S, N]](_ E, items []N) admission.DataPage[S, N, E] {
	return model.Page[N]{Items: items}
}

var _ model.Page[model.Person] = pageOf(admission.PersonEvidence(nil), []model.Person{})
var _ model.Page[model.Page[model.Person]] = pageOf(
	admission.PageEvidence[admission.PersonSource, model.Person, admission.PersonEvidence](nil),
	[]model.Page[model.Person]{},
)
