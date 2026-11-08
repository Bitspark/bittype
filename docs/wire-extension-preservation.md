# Optional Wire extensions: concept preservation

**Owner directions, 3–4 November 2026. Revision-4 kinds implemented in Go source.**
Horos owns canonical type and kind semantics; Bittype owns tooling. One declaration
language describes data and capability types. General structural kinds replace
closed Type/Data bounds and their dedicated evaluator. Ontos retains distinct
ground and optional mixed value domains. This page preserves those boundaries
without retaining superseded implementations.

## Authoritative homes

- [Ontos/wire](https://github.com/Bitspark/ontos/blob/main/docs/design/0007-optional-wire-extension.md)
  owns mixed values, ground embedding and partial recovery.
- [Horos declarations](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md)
  own canonical modules, ordinary kinds, type-level values and ground extraction.
- [Bittype's frontend](horos.md), [native mappings](semantic-generation.md),
  [representation design](type-representations.md) and
  [runtime boundary](runtime-values.md) separate static tooling from adapters.

Ground values and ground judgments retain their meanings. Deixis remains a
generic value structure, not a type model or AST. Hydration stays in Wire/runtime.
The accepted [general-kind design](kinds.md) changes none of those runtime laws.

## Preserved concepts

The earlier proposal at
[`9c868646c09a449c6e682328840434126be1301f`](https://github.com/Bitspark/bittype/tree/9c868646c09a449c6e682328840434126be1301f)
remains historical context. Its mixed-type, lifting, recipe and runtime ideas do
not override current canonical semantics or create a compatibility surface.

| Concept | Current home and disposition |
| --- | --- |
| Ground atoms/tuples, equality, codec and data embeddings | Ground Ontos remains unchanged; no Wire is inserted into its containers or codec. |
| Deixis opaque own values and exact child maps | Used by the optional mixed profile without new Deixis kinds or automatic equality requirements. |
| Ground projection and positional recognition | Mixed admission preserves dense canonical keys and rejects foreign structures. |
| None, absent child, empty atom/tuple and null | Remain distinct; None denotes a tuple node, not absence or null. |
| Mixed tuples with opaque sending leaves | Owned by optional ontos/wire, with immutable structure and explicit handle aliasing. |
| Ground recovery | Partial; succeeds for admitted wire-free values, which use the existing ground codec. |
| Ground any/complement/exact/recursion/refinements | Meanings remain fixed. Ordinary ground-kind membership plus materialization and ground WF protects extraction. |
| Wire message type and send-only authority | Unified `H("wire", E)` remains a static form. Syntax traversal never follows a receiver; runtime recognition is separate. |
| Tuple/compound/list/closed-record/union descriptions | Same primitive forms in one language; ordinary structural kinds classify their nested type descriptions. |
| Modules, imports and definitions | Revision 4 retains the module envelope and adds kind/type definition categories; identities are not semantic equivalence or protocol identities. |
| Indexed binders and ordered arguments | Preserve lexical scope, saturated application, caller/definition environments and literal opacity. Contracts are general kind expressions. |
| Constructor parameters/results/aliases | General type-level values retain their semantic sort independently of declaration parameter count. |
| Open Page and applied Page | One authored Page retains generic bindings and ordinary applications. |
| Page containing an ordinary capability-bearing Cell | Uses the same Page identity. A Data-constrained alias requires membership in the ordinary imported ground kind. |
| Template lifting | Remains removed: no API, transformer, derived Page or compatibility path. |
| Complete declaration context | Required imports retain ordinary identity and binding rules; supplied source and bundled standard modules are distinguished explicitly. |
| Type and recipe identities | Exact declarations and target profiles select recipes; shape or erased native spelling never recovers ownership. |
| Function, get/set Cell and next/cancel Stream | Stay library-owned declarations and behaviors; Horos's examples do not establish a production protocol. |
| Per-language native recipes | Remain a separate design with explicit structured types, argument mappings and dependency choices. |
| Generic adapter composition | Remains future runtime work, including lazy bidirectional parameter conversion. |
| Endpoint ownership and lifetimes | Existing runtime laws remain intact; static kind membership grants no endpoint authority. |
| Dynamic descriptions | Remain descriptions beside generic instances; they do not create compiled native types or establish behavior. |
| AST/CST, typed Semantics and provider stages | Existing agnostic architecture remains. Horos projects ordinary contracts and memberships without kernel kind rules. |
| Syntax and Files | Existing multi-Program and typed Source/Semantics composition remains; missing standard-module source receives no invented Files location. |
| Native representations | Explicit representation permits documented kind-constraint erasure; unavailable enforcement returns Unsupported. |
| Types-only syntax versus executable adapters | The current mappings produce types/interfaces, not runtime implementations. |
| No behavioral catalogue | Function, Stream, Cell and future libraries require no general generator/runtime/Wire/hydration switch or fingerprint. |

## General kinds and the ground boundary

`Data` is an ordinary recursive kind in the explicitly imported `horos.kinds`
module. Flat, Rows, Tree and renamed or imported equivalent definitions use the
same membership/inclusion machinery. Primitive operand requirements use ordinary
kind expressions bound to the chosen prelude definitions. A user's local kind
called Data does not change those requirements.

The kind algebra classifies structural type descriptions, not all possible
runtime inhabitants. An optional Wire field prevents Data membership just as a
required one does. A record with no inhabitants does not thereby lose its Wire
structure. `exact` payloads remain opaque data even when they resemble type or
kind tags. Ground recursion is represented by lexical backedges for structural
observation; its separate value-descent admission rule remains unchanged.

Ground extraction requires a closed proper value, membership in the ordinary
ground kind, materialization and unchanged ground well-formedness. A kind proof
does not rewrite a descriptor or turn a live capability into ground data.
Seq/set/map, open records, intersection/complement and recursive descriptors keep
their existing operand restrictions through ordinary contracts. No mixed
complement, universal live equality or recursive capability law is introduced.

## Clean replacement

The canonical Bittype-owned type grammar and separate Horos/Wire frontend remain
superseded. The closed Type/Data interfaces, DataSummary, dedicated Data evaluator
and revision-3 reader are also replaced. There are no automatic compatibility
readers, schema Data wrappers, second language IDs or executable legacy paths.
Historical research remains historical.

Canonical definitions move to revision 4 with ordinary kind contracts and explicit
result kinds. Prelude names use ordinary imports; changed module identities flow
through those imports. Literal payloads and actual binding relationships retain
their meanings. The module envelope has no new schema-digest field.

Bitwire/Bitruntime's value and hydration behavior is independent of this migration.
No transport, runtime recognition, export/proxy, endpoint lifetime or live-value
codec is added or changed here.

## Distinct judgments

Static syntax admission, canonical module resolution, structural kind membership,
value conformance, mixed-profile admission, live capability recognition, authority
and protocol adherence remain different operations. A successful kind query does
not answer the others. Native representation also differs from enforcement for
arbitrary native clients.

The source uses public Horos Go 0.10.2. Existing historical research and
reported implementation states do not establish delivery of this replacement.
