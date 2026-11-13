package horos

// Names of canonical horos data. The horos library owns them, reads the values
// that carry them and exports none of them. Until it does, this file is where
// bittype spells the ones below. A horos stage that writes or reads one takes
// it from here.
//
// The values are held by canonical module values and by their pins. Changing
// one is a decision of the horos library, not of this package.

// FieldTag labels the wrapper of one record field: the tuple
// (FieldTag, name, FieldRequired or FieldOptional, type). It is in the
// namespace of the descriptor tags, which semantics/horos/go publishes as
// DescriptorTagPrefix.
const FieldTag = "horos/field"

// The words of a record payload. A record is open or closed, and each of its
// fields is required or optional. The source keyword "open" is another name
// with the same spelling: it belongs to the text grammar of the parser and the
// formatter.
const (
	RecordOpen    = "open"
	RecordClosed  = "closed"
	FieldRequired = "req"
	FieldOptional = "opt"
)
