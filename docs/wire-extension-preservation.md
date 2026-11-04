# Optional Wire extensions: concept preservation

**Owner directions, 3-4 November 2026.** Horos owns canonical types; Bittype
owns tooling. The later unification decision removes the two declaration profiles
in favor of Type/Data bounds. Ontos retains distinct ground and optional mixed
value domains. This ledger preserves earlier concepts without retaining obsolete
code or overriding the current grammar.

## Decision and authoritative homes

- [Ontos/wire](https://github.com/Bitspark/ontos/blob/main/docs/design/0007-optional-wire-extension.md)
  owns mixed values, ground embedding and partial recovery.
- [Unified Horos](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md)
  owns one declaration grammar, Type/Data bounds and specialization.
- [Bittype tooling](horos.md), [native representations](type-representations.md)
  and [runtime boundary](runtime-values.md) separate implemented static APIs
  from future generators and adapters.

Ground values/checks retain their meanings. Deixis remains a generic value
structure, not a type model or AST. Hydration remains in Wire/runtime.

## Preservation ledger

The previous Bittype proposal is pinned at
[`9c868646c09a449c6e682328840434126be1301f`](https://github.com/Bitspark/bittype/tree/9c868646c09a449c6e682328840434126be1301f).
Its [mixed declarations](https://github.com/Bitspark/bittype/blob/9c868646c09a449c6e682328840434126be1301f/docs/hydrated-types.md),
[lifting/recipes](https://github.com/Bitspark/bittype/blob/9c868646c09a449c6e682328840434126be1301f/docs/type-representations.md)
and [runtime explanation](https://github.com/Bitspark/bittype/blob/9c868646c09a449c6e682328840434126be1301f/docs/runtime-values.md)
are retained historical evidence. This ledger names both preserved ideas and
deliberate replacements rather than silently deleting them from current guidance.

| Concept | Current home and disposition |
| --- | --- |
| Atom/tuple ground values, equality, codec and data embeddings | Ground Ontos unchanged; no Wire inserted into existing containers or codec. |
| Deixis opaque own values and exact child maps | Reused for the optional mixed profile; no new Deixis kinds or automatic Eq requirement. |
| Ground projection and exact positional recognition | Existing projection retained; mixed admission preserves its dense keys and rejects foreign structures. |
| None, absent child, empty atom, empty tuple and null distinctions | Preserved; None denotes a tuple node, not absence or null. |
| Mixed tuples with opaque sending leaves | Moved to optional ontos/wire; immutable data/structure and handle aliasing are explicit obligations. |
| Ground recovery and canonical bytes | Recovery is partial and succeeds exactly for admitted wire-free values; encode recovered ground with the existing Ontos codec. |
| Data embedding and ground any/complements/exact/recursion/refinements | Ground descriptors occur directly; Data bounds and DataDescriptor preserve the data-only boundary. |
| Wire message type and send-only authority | Retained as unified H("wire", E); syntax traversal never follows a receiver. Runtime recognition is separate. |
| Structural tuple/compound/list/closed record/union types | Same constructors in one language; data classification follows operands and ground layouts remain exact. |
| Modules, exact import pins and declarations | Owned by one Horos revision-2 grammar; exact pins are not semantic equivalence or protocol identities. |
| Indexed generic slots and ordered arguments | Retained with explicit Type/Data bounds, lexical binding, arity, caller-before-callee substitution and literal opacity. |
| Open Page<T> and Page<Person> | Ground forms remain; native generation retains generics instead of expanding every application. |
| Page<Cell<Person>> without an authored duplicate Page | Page<Type> accepts both domains; DataPage<Data> is an ordinary constrained alias of the same Page. |
| Template lifting | Removed completely on 4 November: no API, transformer, derivation metadata or compatibility path. |
| Complete declaration context | Reachable imports, including unused declared imports, validate. Unrelated pool values are bounded/hashed. Exact dependencies and explicit Data bounds remain. |
| Source and recipe identities | Profile-aware declaration identities and provenance; never recover ownership from shape or native erasure. |
| Function<A,B>, get/set Cell<T>, next/cancel Stream<T> | Exact illustrative modules live in Horos's unified specification; behavior stays library-owned. |
| Per-language native recipes and dependency locks | Retained in Bittype's representation design for TypeScript, Go, Rust and C++; no raw source or implicit discovery. |
| Generic adapter composition | Retained as future runtime work: PageAdapter(CellAdapter(PersonAdapter)), with lazy bidirectional parameter conversion. |
| Endpoint ownership, borrowed sending views, scope/export lifetimes | Existing runtime laws preserved; future adapters need explicit cleanup and ownership evidence. |
| Dynamic type descriptions | Retained as descriptions beside generic instances; they do not create compiled native types or prove behavior. |
| AST/CST, semantic facts, resolver/analyzer SPIs and engine routes | Existing agnostic architecture retained; one frontend exports bounds through existing terms, with no kernel change. |
| Native type-only syntax versus executable adapters | Retained boundary; adapter generation needs separately qualified syntax/checking/formatting/compiler support. |
| No central behavioral catalogue | Function, Stream, Cell and future libraries require no generator/runtime/Wire/hydration switch, fingerprint or built-in registration. |

## Deliberate replacements, not compatibility surfaces

The original Bittype-owned contract grammar was superseded by Horos ownership.
The separate Horos/Wire profile was subsequently superseded by one revision-2
Horos declaration language. Remove old stage directories, profile IDs, envelope
readers, schema Data wrappers and closed-embedding APIs. There are no aliases,
automatic conversions or duplicate canonical authorities. Template lifting stays
removed. Historical sources and review evidence remain immutable.

Old data-only slots migrate explicitly to Data; old mixed slots migrate to Type.
Recompute canonical module pins from rewritten dependency leaves. Preserve literal
payloads. Open not/set/map/recursion operands remain supported with Data bounds;
capability operands are rejected until separately specified laws exist. The one
grammar does not introduce mixed top/complement, live equality or recursive named
applications. Ground any and conformance never admit live capabilities.

Bitwire/Bitruntime's Ontos 0.16.0 native value migration is independent. Unifying
declarations changes no transport, hydration, endpoint ownership or runtime code.

## Review and implementation obligations

[Research 0003](../research-docs/0003-optional-wire-extensions.md) records the
reviewed brief, frozen input, recoverable service submission, raw expert advice
and recommendation-by-recommendation evaluation. Research 0001's historical
type/generic work and [research 0002](../research-docs/0002-semantic-analysis-kernel.md)
remain evidence; they do not override the current profiles or encoding rules.

Before claiming implementation, independently authored cases must distinguish:

1. Mixed structure admission from runtime recognition, operational liveness and
   protocol adherence. A Send-shaped object or reference-looking data is no authority.
2. Captured bytes/keys from mutable caller storage, and repeated occurrences of
   the same handle from cloned resources. Structure traversal performs no IO.
3. Ground any/complement from any hypothetical mixed-domain operations.
4. Closed data recovery and generic specialization from
   equality of ASTs, declaration pins or native type names.
5. Retained open applications from closed specialization, including unused-import
   errors, generic bounds, exact literals and caller-side scope.
6. Explicit per-language recipes from lookalikes, implicit fallback and accidental
   registration-order selection. Native compiler success is not protocol proof.
7. Typed compiled adapters from dynamic descriptions, and borrowed views from
   owned endpoint/scope resources. Binding failures need explicit cleanup laws.

Value-profile checks qualify structure and immutable capture in their owning
repositories. Horos tests qualify unified static bounds, direct pagination,
specialization, diagnostics and exhaustion. Bittype adds AST round trips, bound
facts, immutable receipts, origins and dependency boundaries. Native recipes and
runtime adapters remain future acceptance criteria, not claims of this change.
Research 0003 records the historical design review; the owner's later unification
supersedes its profile and lifting recommendations. Exact-source hosted gates and
independent review qualify each actual delivered revision.
