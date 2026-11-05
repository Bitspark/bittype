package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

// want-error: does not satisfy
type Bad = admission.DataPage[admission.PersonSource, model.Person, *admission.PersonEvidence]
