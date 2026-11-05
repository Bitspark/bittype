package consumer

import (
	"admission.example/admission"
	"admission.example/foreign"
	"admission.example/model"
)

type SourceAlias = admission.PersonSource
type PayloadAlias = model.Person
type EvidenceAlias = admission.PersonEvidence
type ReindexedAlias[S, N any] = EvidenceAlias

type People = admission.DataPage[SourceAlias, PayloadAlias, EvidenceAlias]
type AlsoPeople = admission.DataPage[
	foreign.ReexportedSource, foreign.ReexportedPerson, foreign.ReexportedEvidence,
]

var _ People = AlsoPeople{}
var _ People = admission.DataPage[
	SourceAlias, PayloadAlias, ReindexedAlias[SourceAlias, PayloadAlias],
]{}
