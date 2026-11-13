// Package mapping holds the identities that the TypeScript mapping side
// declares: the one generation profile of a pair whose source is TypeScript,
// and the publisher of the contributions that state it. It holds identities
// only and imports nothing.
//
// It exists for the two other pairs of bittype#196 (decision O43), which show
// that the contribution contract and its kernel assume no native target.
// TypeScript as a source is not supported beyond what the pair
// typescript-to-horos covers. That pair is done-when 1 of bittype#173 in the
// smallest form; its items 2 to 5 are open there.
package mapping

// CanonicalizationProfile is the profile of the contributions whose source is
// TypeScript and whose target is horos. It promises a horos type whose values
// are the JSON-like values of the TypeScript type under these explicit
// choices, and it covers nothing for which a choice is not stated here:
//
//   - a TypeScript object type whose properties are all required and none
//     readonly is a closed horos record with one required field for each
//     property, named as the property is spelled;
//   - string is the horos scalar utf8-text. Stated limit: a TypeScript string
//     is a sequence of UTF-16 code units and may hold an unpaired surrogate,
//     which is no horos text. The profile promises the well-formed strings;
//   - boolean is the horos scalar bool.
//
// number has no choice here and is refused: "Future canonicalization needs
// explicit semantic choices: native numeric and nullable forms do not uniquely
// identify source domains or absence semantics" (docs/generation.md). The
// same sentence is why an optional property is refused.
//
// Nothing checks these promises. No binding is built for the pair.
const CanonicalizationProfile = "typescript-canonicalization/1"

// Publisher names the publisher of the contributions that state
// CanonicalizationProfile: the library that publishes them, which is bittype
// (row 15). The name authenticates nothing: until a composition lock binds it
// (row 15a), registration in one process is the binding.
//
// COPY (AC-11, finding AG-6 of the two other pairs): the same name is spelled
// in transform/syntax/horos/mapping/go, because the vocabulary of a
// contribution may import the mapping identities of its own source language
// only. The publisher is bittype whatever the source language is: one name
// with a spelling per source language, and no agnostic home.
const Publisher = "github.com/Bitspark/bittype"
