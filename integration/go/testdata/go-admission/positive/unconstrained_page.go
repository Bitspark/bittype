package consumer

import "admission.example/model"

// Canonical Page remains unconstrained; only DataPage has the Data obligation.
var _ model.Page[model.Wire]
var _ model.Page[func()]
var _ model.Page[chan int]
