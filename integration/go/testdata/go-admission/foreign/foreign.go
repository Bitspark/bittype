package foreign

import (
	"admission.example/admission"
	"admission.example/model"
)

// Same spelling and function signature, different private method identity.
type PersonEvidence func(admission.PersonSource) model.Person

func (PersonEvidence) dataEvidence() {}

type ReexportedEvidence = admission.PersonEvidence
type ReexportedSource = admission.PersonSource
type ReexportedPerson = model.Person
