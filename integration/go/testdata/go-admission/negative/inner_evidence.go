package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

// want-error: does not satisfy
type Bad = admission.DataPage[
	admission.PageSource[admission.PersonSource],
	model.Page[model.Person], admission.PersonEvidence,
]
