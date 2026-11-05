# Unified Horos frontend

**Owner decisions, 4 November 2026. Revision 4 implemented in Go source.**
Horos has one declaration language for data and capability types. The replacement
uses [general structural kinds](kinds.md), with ordinary named/parameterized
definitions, unions/intersections, structural patterns, guarded recursive groups
and constructor contracts. Bittype consumes public Horos Go 0.10.1.

## One language and one authority

Horos owns the [canonical grammar](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md),
membership, inclusion, application and ground extraction. Bittype owns
`syntax/horos`, `parse/syntax/horos`, `check/syntax/horos`, `resolve/files/horos`, `analyze/files/horos`,
`format/syntax/horos` and `provider/horos`, with their private bridge. The language ID
stays `horos`; its syntax and semantic profiles advance together.

Revision 4 retains the existing module envelope and ordinary imports. One
declarations map contains kind and type definitions. There is no schema-digest
field, separate formation artifact or parallel Horos/Wire language. Ground-only
consumers retain their existing ground API meanings.

`Type` and `Data` are ordinary definitions in the explicitly imported
`horos.kinds` prelude. Primitive sorts/top/bottom do not inspect those names.
Flat, Rows and Tree use the same machinery without checker or generator changes.
Primitive operand contracts bind to ordinary prelude declaration identities;
a same-spelled local definition cannot change them.

## Generics and results

One `Page<T>` accepts admitted proper types. `DataPage<D: h.Data>` is an ordinary
constrained definition using that Page, where `h` is a prelude import. The known
Page body can establish requested guarantees through general kind queries. It is
not copied or lifted for each use.

Type definitions may return constructor values. Their result kind must have the
corresponding sort; omission means `top(P)`, not inferred higher-order shape.
Applications retain an expression callee, argument order, definition scope and
caller bindings. Open bodies and ignored arguments remain subject to original
contracts. Kind definitions have sort-annotated parameters and ordinary imports.

Checked symbolic values, declared contracts and individual membership results
replace Claim/Interface and DataSummary. No principal kind or enumeration of all
true memberships is promised. Ground extraction requires a closed proper value,
membership in the ordinary ground kind, materialization and unchanged ground
well-formedness.

## Bittype composition

Syntax admission stays separate from semantic resolution. The resolver asks Horos
to establish meaning and exports provider-owned terms through typed
`Semantics[horos.Language]`. Analysis consumes its explicit receipt. The shared
graph and engine learn no Data, Scalar, Record, List or Wire kind rules.

The implemented [Syntax and Files composition](source.md) remains intact:
multiple Program trees share an explicit captured entry and supplied module pool.
File paths do not confer visibility. Files formatting establishes a new capture
and semantic association for its output rather than reusing receipts from changed
syntax.

The four [native pagination mappings](semantic-generation.md) use explicit
representation or enforcement profiles. Representation permits documented kind
constraint erasure for Horos-admitted uses. Enforcement returns Unsupported until
a target encoding exists. Kind definitions do not become native type exports;
unsupported exported type definitions are never silently omitted.

## Clean replacement

Migrate canonical documents to revision 4, replace bound/interface fields with
ordinary kind expressions, add result kinds and preserve explicit bindings.
Authors using prelude names add ordinary imports. Rewritten module identities
flow through normal imports; literal payloads retain their exact values.

Remove superseded readers, the closed Type/Data proof route, schema Data wrappers
and compatibility aliases. The earlier separate Horos/Wire packages and template
lifting remain removed. Historical research does not define a second executable
language.

Ground Ontos/Horos, optional Ontos/wire values and Bitwire/Bitruntime hydration
retain their boundaries. Function, Stream, Cell and future contracts remain
library-owned. This migration changes no endpoint ownership, runtime behavior or
transport contract. The [preservation map](wire-extension-preservation.md) records
those distinctions.
