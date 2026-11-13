package horos

// Names of canonical horos data. The horos tree carries canonical values
// literally: a descriptor operand keeps the exact canonical layout, and labels
// and field wrappers keep their literal roles. So every package that writes or
// reads the tree, or builds the value it mirrors, agrees on these names.
//
// Their owner is the horos library, which reads the values and exports none of
// the names (Bitspark/horos#53). This declaration is a stand-in until it does.
// The values are held by canonical module values and by their pins. Changing
// one is a decision of the horos library, not of this package.

// The ontos tag namespaces of canonical horos values: declaration forms and
// descriptor constructors.
const (
	DeclarationTagPrefix = "horos/decl/"
	DescriptorTagPrefix  = "horos/"
)

// Labels that follow a prefix.
//
// ArrowLabel follows DeclarationTagPrefix in the tag of an arrow kind. It is a
// label and not a whole tag because the model reader appends it to the prefix
// it is configured with.
//
// RecordLabel follows DescriptorTagPrefix in the tag of a record payload. It is
// the constructor name of a record: the model reader takes that name from the
// tree's constructor attribute and appends it to its prefix, and never uses
// this constant. Its one user today is the source reader, which builds the
// whole tag to compare with. The parser, which writes the constructor into the
// tree, still spells it by literal.
const (
	ArrowLabel  = "arrow"
	RecordLabel = "record"
)

// FieldTag labels the wrapper of one record field: the tuple
// (FieldTag, name, FieldRequired or FieldOptional, type). The text stages
// write and compare the whole atom, so it is declared whole.
const FieldTag = DescriptorTagPrefix + "field"

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
