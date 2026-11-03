# Optional Wire extensions: concept preservation

**Owner direction, 3 November 2026. Documentation only.** The owner selected
`ontos/wire` and `horos/wire` as optional extensions, replacing the Bittype-owned
contract model, then explicitly required preservation of prior concepts and an
expert consultation. This page makes that preservation reviewable; it is not a
claim that the proposed extensions, generators or runtime adapters exist.

## Decision and authoritative homes

- [Ontos optional profile](https://github.com/Bitspark/ontos/blob/main/docs/design/0007-optional-wire-extension.md)
  owns mixed values over Deixis, ground embedding and partial recovery.
- [Horos optional profile](https://github.com/Bitspark/horos/blob/main/docs/spec/wire-types.md)
  owns mixed type/declaration meaning and partial lifting.
- [Bittype tooling](hydrated-types.md), [native representations](type-representations.md)
  and [runtime boundary](runtime-values.md) explain how those definitions are
  projected, generated and eventually interpreted.

Ground cores retain their existing domains and identities. The optional profile
is a sibling to the frozen ground stack, not another conservative L1/L2 data
reading. Deixis remains a generic value structure, not a type model or AST.
Hydration stays above ground semantics in existing Wire/runtime machinery.

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
| Data embedding and ground any/complements/exact/recursion/refinements | Retained by Data(G) in horos/wire; closed ground semantics are not widened. |
| Wire message type and send-only authority | Retained in horos/wire; syntax traversal never follows a receiver. Runtime recognition is separate. |
| Structural tuple/compound/list/closed record/union types | Retained in the bounded mixed profile with existing ground layouts and ground record keys. |
| Modules, exact import pins and declarations | Owned by each explicit Horos profile; exact pins are not semantic equivalence or protocol identities. |
| Indexed generic slots and ordered arguments | Retained with lexical binding, explicit arity, caller-before-callee substitution and literal opacity. |
| Open Page<T> and Page<Person> | Ground forms remain; native generation retains generics instead of expanding every application. |
| Page<Cell<Person>> without an authored duplicate Page | Explicit partial lifting derives a distinct mixed declaration; its binder can accept mixed types. |
| Ordered partial lifting and unsupported forms | Retained as Horos-owned semantics; Data boundaries and unsupported diagnostics remain explicit. |
| Full source validation versus derived dependency closure | All source imports validate first; output follows retained applications and stops at Data. Root/Other/Alias counterexample remains in the lifting docs. |
| Source, derived and recipe identities | Distinct profile-aware identities and provenance; never recover ownership from shape or native erasure. |
| Function<A,B>, get/set Cell<T>, next/cancel Stream<T> | Exact illustrative modules move to Horos's optional specification; behavior stays library-owned. |
| Per-language native recipes and dependency locks | Retained in Bittype's representation design for TypeScript, Go, Rust and C++; no raw source or implicit discovery. |
| Generic adapter composition | Retained as future runtime work: PageAdapter(CellAdapter(PersonAdapter)), with lazy bidirectional parameter conversion. |
| Endpoint ownership, borrowed sending views, scope/export lifetimes | Existing runtime laws preserved; future adapters need explicit cleanup and ownership evidence. |
| Dynamic type descriptions | Retained as descriptions beside generic instances; they do not create compiled native types or prove behavior. |
| AST/CST, semantic facts, resolver/analyzer SPIs and engine routes | Existing agnostic architecture retained; optional frontend supplies semantics and distinct endpoints, not core cases. |
| Native type-only syntax versus executable adapters | Retained boundary; adapter generation needs separately qualified syntax/checking/formatting/compiler support. |
| No central behavioral catalogue | Function, Stream, Cell and future libraries require no generator/runtime/Wire/hydration switch, fingerprint or built-in registration. |

## Deliberate replacements, not compatibility surfaces

The Bittype-owned `bittype/contract/*` grammar and planned `contract` language
directories were never implemented. They are replaced by the Horos-owned optional
profile and proposed `horoswire` tooling directories. There is no alias accepting
old tags, automatic envelope conversion or dual canonical authority.

The older proposal for a Wire constructor in ground Horos revision 2 is also
superseded. Reusing the conceptual name `horos/wire` does not reinstate that
ground-domain widening. Ground `horos/decl/*` revision 1 remains unchanged; the
optional grammar has its own envelope and semantic profile.

Ground-only open forms such as `not(T)`, `set-of(T)` or guarded recursion with
an open slot are not silently removed or weakened: they remain ground Horos
concepts. The initial partial lift rejects them, while closed specializations
can be explicitly embedded. Mixed top/complement, live set/map keys, recursive
named applications and broader lifting require separately defined laws and
qualification before they can be admitted. No current syntax promises them.

The proposal does not delete current Bitwire/Bitruntime carriers. Their pinned
behavior is recorded in [runtime values](runtime-values.md). Future alignment
may use optimized native representations, provided profile laws and runtime
ownership are demonstrated. Moving ownership alone does not remove duplication.

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
3. Data(any)/ground complement from any hypothetical mixed-domain operations.
4. Ground-membership/layout preservation after lifting and binding from equality
   of ASTs, declaration pins or native type names.
5. Retained applications from closed embedding, including unused-import errors,
   open unsupported forms, exact literals and caller-side scope.
6. Explicit per-language recipes from lookalikes, implicit fallback and accidental
   registration-order selection. Native compiler success is not protocol proof.
7. Typed compiled adapters from dynamic descriptions, and borrowed views from
   owned endpoint/scope resources. Binding failures need explicit cleanup laws.

These are future acceptance criteria, not tests claimed to have run for a
documentation change. The current delivery validates links, examples, consistency,
the preservation inventory and unchanged source boundaries through repository gates.
