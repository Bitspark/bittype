package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

// want-error: invalid recursive type|cannot refer to itself
type Loop = admission.PageEvidence[admission.PersonSource, model.Person, Loop]
