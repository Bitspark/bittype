package consumer

import (
	"admission.example/admission"
	"admission.example/model"
)

// Arbitrary values of an already authorized evidence type do not change its
// indexed source/native pair. The admission claim is about types, not values.
var _ admission.PersonEvidence = (func(admission.PersonSource) model.Person)(nil)
var _ admission.PersonEvidence = admission.PersonEvidence(nil)

type Ignored = admission.Ignore[
	admission.PersonSource, model.Person, admission.PersonEvidence,
]

var _ Ignored = ""
