package integration_test

const (
	admitsOwnTree     = "admits the tree it built with its own language's grammar before publishing it"
	admitsOwnInput    = "admits the tree it is given with its own language's grammar before formatting it"
	verifiesOwnSource = "reparses captured source with its own language's parser to verify that the captured tree corresponds to it"
	readsResolution   = "reads the resolution contract identity that its analysis is bound to"
)

// providerExceptions lists every import of a stage provider package
// (<stage>/<dimension>/<language>/provider/go) by a production package that is
// not a composition. Only a composition chooses a provider, so each such
// import is admitted by name: the importer, the provider, the exact exported
// symbols used and the reason. The check fails on an import without an entry,
// on a use of any other symbol, and on an entry or symbol that is not used.
//
// An entry without an issue is accepted. An entry with an issue is an open
// finding of bittype/ACX-009; that issue removes the import and this entry.
var providerExceptions = []providerException{
	// Accepted: the caller needs its own grammar's answer, not the provider a
	// composition selected.
	{Importer: "parse/syntax/cpp/provider/go", Provider: "check/syntax/cpp/provider/go", Symbols: []string{"Check"}, Reason: admitsOwnTree},
	{Importer: "parse/syntax/golang/provider/go", Provider: "check/syntax/golang/provider/go", Symbols: []string{"Check"}, Reason: admitsOwnTree},
	{Importer: "parse/syntax/horos/provider/go", Provider: "check/syntax/horos/provider/go", Symbols: []string{"Check"}, Reason: admitsOwnTree},
	{Importer: "parse/syntax/rust/provider/go", Provider: "check/syntax/rust/provider/go", Symbols: []string{"Check"}, Reason: admitsOwnTree},
	{Importer: "parse/syntax/typescript/provider/go", Provider: "check/syntax/typescript/provider/go", Symbols: []string{"Check"}, Reason: admitsOwnTree},
	{Importer: "format/syntax/cpp/provider/go", Provider: "check/syntax/cpp/provider/go", Symbols: []string{"Check"}, Reason: admitsOwnInput},
	{Importer: "format/syntax/golang/provider/go", Provider: "check/syntax/golang/provider/go", Symbols: []string{"Check"}, Reason: admitsOwnInput},
	{Importer: "format/syntax/horos/provider/go", Provider: "check/syntax/horos/provider/go", Symbols: []string{"Check"}, Reason: admitsOwnInput},
	{Importer: "format/syntax/rust/provider/go", Provider: "check/syntax/rust/provider/go", Symbols: []string{"Check"}, Reason: admitsOwnInput},
	{Importer: "format/syntax/typescript/provider/go", Provider: "check/syntax/typescript/provider/go", Symbols: []string{"Check"}, Reason: admitsOwnInput},
	{Importer: "resolve/files/horos/provider/go", Provider: "check/syntax/horos/provider/go", Symbols: []string{"Check"}, Reason: "admits the tree it is given with its own language's grammar before reading its model"},
	{Importer: "internal/rustsem", Provider: "parse/syntax/rust/provider/go", Symbols: []string{"Parse"}, Reason: verifiesOwnSource},
	{Importer: "internal/typescriptsem", Provider: "parse/syntax/typescript/provider/go", Symbols: []string{"Parse"}, Reason: verifiesOwnSource},

	// Open until #220: an identity that a provider package publishes and a
	// vocabulary package should.
	{Importer: "analyze/files/golang/provider/go", Provider: "resolve/files/golang/provider/go", Symbols: []string{"Contract"}, Reason: readsResolution, Issue: 220},

	// Open until #221: a query service that shares a package with a provider.
	{Importer: "verification/files/horos-to-go/go", Provider: "analyze/files/golang/provider/go", Symbols: []string{"Admission"}, Reason: "takes the Go query service's admission type, which lives in the analysis provider package", Issue: 221},
}
