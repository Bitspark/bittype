package consumer

import "admission.example/admission"

type Fake = admission.PersonEvidence

// want-error: cannot define new methods on non-local type
func (Fake) dataEvidence() {}
