package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

type Fake func(admission.WireSource) model.Wire

func (Fake) dataEvidence() {}

func pageOf[S, N any, E admission.DataEvidence[S, N]](_ E, item N) admission.DataPage[S, N, E] {
	return model.Page[N]{Items: []N{item}}
}

// want-error: unexported method dataEvidence
var _ = pageOf(Fake(nil), model.Wire{})
