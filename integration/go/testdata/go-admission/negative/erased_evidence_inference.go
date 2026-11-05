package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

func makePage[E admission.DataEvidence[admission.WireSource, model.Wire]]() admission.DataPage[
	admission.WireSource, model.Wire, E,
] {
	return model.Page[model.Wire]{}
}

// The aliased result erases E. Assignment must not infer an unauthorized witness.
// want-error: cannot infer E|does not satisfy|missing method dataEvidence
var _ func() model.Page[model.Wire] = makePage
