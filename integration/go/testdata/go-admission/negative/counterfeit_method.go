package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

type Fake func(admission.WireSource) model.Wire

func (Fake) dataEvidence() {}

// want-error: unexported method dataEvidence
type Bad = admission.DataPage[admission.WireSource, model.Wire, Fake]
