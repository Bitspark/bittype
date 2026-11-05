package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

type Source admission.PersonSource

// want-error: does not satisfy
type Bad = admission.DataPage[Source, model.Person, admission.PersonEvidence]
