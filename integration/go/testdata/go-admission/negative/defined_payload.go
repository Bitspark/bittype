package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

type Person model.Person

// want-error: does not satisfy
type Bad = admission.DataPage[admission.PersonSource, Person, admission.PersonEvidence]
