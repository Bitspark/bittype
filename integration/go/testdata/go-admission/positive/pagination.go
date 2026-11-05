package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

type People = admission.DataPage[
	admission.PersonSource, model.Person, admission.PersonEvidence,
]

var _ model.Page[model.Person] = People{}
var _ People = model.Page[model.Person]{}
