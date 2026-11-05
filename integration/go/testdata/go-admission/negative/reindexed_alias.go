package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

type Fake[S, N any] = admission.PersonEvidence

// want-error: does not satisfy
type Bad = admission.DataPage[
	admission.WireSource, model.Wire, Fake[admission.WireSource, model.Wire],
]
