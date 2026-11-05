package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

// want-error: does not satisfy
type Bad = admission.DataPage[admission.WireSource, model.Person, admission.PersonEvidence]
