package consumer

import (
	"admission.example/admission"
	"admission.example/foreign"
	"admission.example/model"
)

// want-error: unexported method dataEvidence
type Bad = admission.DataPage[admission.PersonSource, model.Person, foreign.PersonEvidence]
