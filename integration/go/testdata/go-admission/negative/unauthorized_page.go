package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

// want-error: does not satisfy
type Forged = admission.PageEvidence[admission.WireSource, model.Wire, admission.PersonEvidence]
