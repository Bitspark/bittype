// Package pipeline defines typed transformation contracts.
//
// It imports no AST, language, generator, formatter or application package.
// Execution, composition, scheduling and discovery are not implemented here.
package pipeline

// Stage is the contract of a transformation from I to O.
//
// A stage must treat its input and reachable storage as read-only and own any
// mutable output storage it creates. An error means the output must not be used.
// Stage does not enforce those behavioral requirements or prove that an AST is
// well-formed. The owner of a concrete stage defines its checks and diagnostics.
//
// The nil value denotes no configured stage and must not be invoked.
// No default stage, registry, cancellation policy or hidden IO is installed.
type Stage[I, O any] func(I) (O, error)
