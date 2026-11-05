package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

type Fake struct{ admission.PersonEvidence }

// want-error: does not satisfy
type Bad = admission.DataPage[admission.PersonSource, model.Person, Fake]
