package consumer

import "admission.example/admission"

// want-error: does not satisfy
type Bad = admission.DataPage[admission.PersonSource, func(), admission.PersonEvidence]
