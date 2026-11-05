package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

type Fake admission.PersonEvidence

// want-error: missing method dataEvidence
type Bad = admission.DataPage[admission.PersonSource, model.Person, Fake]
