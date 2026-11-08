# bittype/ACX-006 — Common semantic traversal delegated to source integrations

**Status: corrected, 8 November 2026.** The generation kernel owns traversal of
the shared semantic graph. Source integrations contribute only interpretation
and judgments. [GEN-01](../generation-gaps.md#gen-01) records the delivered
division, its completion criteria and its supported scope.
[Research 0013](../../research-docs/0013-generation-source-traversal.md) derived
the replacement interface.

## Actual ownership failure

The public generation Source/Catalog contract asked each source integration for
two different things:

- navigation of declarations, binders, references and application arguments;
- interpretation of source meanings.

horos's `Open` and `Admit` performed both roles. `CheckChild` re-derived the child
relation, and a horos pre-walk built its own term-ownership index. The shared
semantic model already carried all of these structural facts.

The consequence was that every source had to re-create a navigation adapter for
facts the common semantic context already owns. The kernel received a
source-specific reading of common structure instead of traversing the
established graph once. The independent test source demonstrated composability,
but it did so by re-implementing that navigation from hand-built tables.

This finding did not assert a horos dispatch table in the kernel, or missing
source checking. It concerned where common traversal lived. Separating semantic
planning from AST construction in PR #155 did not repair this different boundary.

## Rules involved

- [AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02):
  common structural navigation belongs to the shared mechanism; source
  interpretation belongs to its language owner.
- [AC-06](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-06):
  successful generation through the first reader did not redefine the intended
  kernel/language division or make its correction optional.
- [AC-01](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-01):
  the defect was described explicitly instead of claiming that passing
  composition tests completed the architecture.
- [AC-05](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-05)
  and [AC-09](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-09)
  governed the replacement. Guarantees are preserved through composition, and the
  interface was derived from the domain rather than from the incumbent methods.

## Correction

**Admission.** The kernel admits the request against the source contract. It
requires exactly one evidence requirement that publishes a graph, and unpacks that
graph itself (`generation.AdmitInput`).

**The walk.** The kernel walks the graph directly:

- root site to its declaring entity;
- declaration to its body and real binders;
- references to their targets;
- applications to their callee, with arguments paired by `DirectParameterBindings`;
- provider expressions to their operands, with locators passed through unread.

It reads no entity kind, name, operator, payload, locator meaning or fact.

**Handles.** Occurrence, declaration and binder handles are kernel-issued and
valid for one run. An occurrence carries its declaration, and two paths to a
shared term are two occurrences.

**Scope and captures.** Scope is direct binder ownership: a reference to another
declaration's binder is `ErrScope`. A capture must be a child the kernel minted
for the current occurrence.

**Source contract.** The source supplies `Interpretation`:

- declaration classification and selectors;
- descriptor selectors;
- judgments of declarations and whole occurrences, as Admitted, Unsupported,
  Rejected or Undetermined.

It also supplies a reader that its own mapping libraries use. Mapping callbacks
receive the kernel's views and that reader, never the graph or the
interpretation.

**Whole-use admission.** The kernel judges every actual of a use before any
binding or convention represents it, including actuals the native convention
omits. Undetermined fails as missing evidence and Rejected as `ErrRejected`, so
an uncovered fact is never reported as a rejection.

**Removed with no aliases:**

- the old reader surface: `Open`, `Admit`, `CheckChild`, `Authenticate`,
  `CheckInput`, `Owners`, `RootEnvironment`, `Catalog`, `SourceSession`,
  `ResolveRoots`, `OpenData` and `UseData`;
- the source address spaces, and the effective-subject and formal-subject uses;
- horos's pre-walk;
- the independent source's navigation tables.

[Kernel contract](../../contracts/generation/go/source.go) ·
[Walk and handles](../../contracts/generation/go/occurrence.go) ·
[horos interpretation](../../transform/syntax/horos/source/go/source.go) ·
[Independent source](../../integration/go/testdata/generationfixture/fixture.go)

## Qualification

horos and the independent source drive the same kernel walk. The independent
source now publishes an ordinary graph through its own analysis.

**Positive cases:**

- horos nested `Page<Cell<Person>>`, with a library-owned Cell selected by Cell's
  binder;
- the independent `Page<Page<T>>` and library-owned `Cell<Person>`, where Cell's
  body is never opened;
- binder-keyed reversed arguments.

**Negative cases:**

- omitted actuals that are refuted, uncovered or admitted;
- a binder of another declaration, including through a shared term;
- a capture that is not a child, and a handle kept from another run;
- recursion, parameter-headed applications, and unapplied generic declarations;
- foreign and missing evidence;
- horos applied constructor aliases and curried applications, which now report
  Unsupported at their site rather than malformed input.

Each kernel check was seeded with a mutant that turned its test red. Covered
checks: actual admission, binder ownership, child captures, the
undetermined/rejected distinction, and the structural application check.
[Traversal qualification](../../integration/go/generation_traversal_test.go) ·
[Planning tests](../../integration/go/generation_planning_test.go)

## Limits

**Scope.** Only direct binder ownership is supported; there are no nested generic
scopes. A language that needs nested scopes must contribute its own scope
judgment, and the kernel infers no `Owner`-chain rule.

**Structural reading.** Applications need a direct declaration callee and complete
arguments. Type-level constructor values still need an explicit representation.

**Operands.** Provider-expression operands are visited only when a mapping
captures them. The source's whole-occurrence judgment covers them.

**Identity.** Effective-subject equivalence is not defined.

**horos judgments.** horos reads no membership facts, so its only non-admitted
judgments are Unsupported and Undetermined.

**Out of scope.** Recursive generation (bittype#164) and native enforcement of
source bounds ([GEN-02](../generation-gaps.md#gen-02)) are unchanged. Declared
representation-profile erasure is not evidence about GEN-01, and this correction
does not implement native constraint enforcement.
