package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

func pageOf[S, N any, E admission.DataEvidence[S, N]](_ E, item N) admission.DataPage[S, N, E] {
	return model.Page[N]{Items: []N{item}}
}

// want-error: missing method dataEvidence
var _ func(func(admission.WireSource) model.Wire, model.Wire) model.Page[model.Wire] = pageOf
