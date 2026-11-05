package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

// want-error: missing method dataEvidence
type Bad = admission.DataPage[
	admission.PersonSource, model.Person, func(admission.PersonSource) model.Person,
]
