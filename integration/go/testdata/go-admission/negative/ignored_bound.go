package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

// Even a result that mentions none of its parameters must check the bound.
// want-error: does not satisfy
type Bad = admission.Ignore[admission.WireSource, model.Wire, admission.PersonEvidence]
