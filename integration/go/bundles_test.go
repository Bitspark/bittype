package integration_test

// builtinBundleNames lists the builtin contributions of a pair from horos to
// a native language, in the order a test registers them. Registration order
// is no preference. The record contribution delegates the presence of an
// optional field to the optional contribution through the planner: neither
// names the other, and no consumer creates a key for them.
var builtinBundleNames = []string{"text", "boolean", "list", "record", "declarations", "optional"}
