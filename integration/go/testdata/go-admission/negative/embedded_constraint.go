package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

// want-error: outside a type constraint
type Fake struct {
	admission.DataEvidence[admission.WireSource, model.Wire]
}
