package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

// want-error: does not satisfy
type Bad = admission.DataPage[
	admission.PageSource[admission.PersonSource], model.Person,
	admission.PageEvidence[admission.PersonSource, model.Person, admission.PersonEvidence],
]
