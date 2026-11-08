# bittype/ACX-006 — Common semantic traversal delegated to source integrations

**Status: correction pending independent verification, 8 November 2026.** [#191](https://github.com/Bitspark/bittype/pull/191) delivered the
kernel-owned traversal, but independent verification by `codex-hydrated-type-design` (board claim `bittype-source-traversal-20261108`, record 000003) rejected it: the public planning
boundary still let a source interpret other evidence, and planning ignored a
selector's language. [#192](https://github.com/Bitspark/bittype/pull/192) corrects both, and a diagnostic-site defect. The case stays open until that
correction is independently verified. The generation kernel owns traversal of
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
- [AC-05](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-05)
  and [AC-07](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-07)
  were violated by the first delivery itself. Its exported planner accepted a
  capture and an already opened source as separate fields, and its rule lookup
  dropped the selector's language. Both composed guarantees therefore failed at
  the public boundary, while the case already claimed the correction. Green
  tests over the registry path did not establish that claim.

## Correction

**Admission.** `Plan` admits the input against the source provider's own
contract. It requires exactly one evidence requirement that publishes a graph,
unpacks that graph itself, and opens the provider over exactly the resulting
capture. No public value pairs an opened source with a capture, and handles
resolve only against that capture. Interpreting that capture is the provider's responsibility: the
kernel cannot check where a provider's answers come from.

**The walk.** The kernel walks the graph directly:

- root site to its declaring entity;
- declaration to its body and real binders;
- references to their targets;
- applications to their callee, with arguments paired by `DirectParameterBindings`;
- provider expressions to their operands, with locators passed through unread.

It reads no entity kind, name, operator, payload, locator meaning or fact.

**Handles.** Occurrence, declaration and binder handles are kernel-issued and
valid for one run. An occurrence carries its declaration, and two paths to a
shared term are two occurrences. A source resolves a handle against its capture
(`In(Captured)`), not just the graph. Two captures that share a graph but differ
in evidence are different, so an interpretation opened over other evidence
cannot read this run's handles.

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

**Selectors.** Mapping selectors are frozen in the capture's source language.
Every selector the source returns must have that same language and profile.
Otherwise planning fails with `ErrScope` before any exact or default rule is
chosen, so key comparison is exact.

**Diagnostic sites.** A kernel failure is reported at its term's syntax site when
the graph records exactly one. A term recorded at several sites is shared, and
the graph does not say which of them this occurrence is, so the failure is
reported at the declaration the occurrence is in. A use whose callee the source
does not generate, or whose callee selector is outside the source language, is
a failure of that use, not of the callee.

**Whole-use admission.** The kernel admits an occurrence with its whole subtree
before any binding or convention represents it. Every occurrence in the subtree
needs a structural reading in scope and an Admitted judgment. That covers actuals
the native convention omits, content nested in them, and operands no mapping
captures. A generated declaration's body is admitted even if its recipe does not
represent it. Undetermined fails as missing evidence and Rejected as `ErrRejected`, so
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

- omitted actuals that are refuted, uncovered or admitted; refuted content
  nested in an omitted actual; a foreign binder or unapplied generic as an
  omitted actual; a foreign binder in an uncaptured operand, and in a body the
  recipe never represents;
- a declaration the source rejects, and a zero verdict;
- a binder of another declaration, including through a shared term;
- a capture that is not a child (through both `Mapped` and `Pattern.Capture`),
  a declaration recipe mapping a non-body use, and a handle kept from another run;
- recursion, parameter-headed applications, and unapplied generic declarations;
- foreign and missing evidence;
- horos applied constructor aliases, curried applications and
  constructor-sorted parameters, which now report Unsupported at their site
  rather than malformed input;
- an interpretation opened over other evidence of the same graph, through the
  registry and through the exported planner;
- a source selector that differs only in its language's profile, identity or
  Go marker, for an exact binding, an exact declaration recipe and the default
  recipe;
- a shared term recorded at two syntax sites, and failures of a use reported at
  the use rather than at its callee.

Before #191 merged, two independent reviews found gaps, and both were fixed.

- The first revision let omitted actuals skip the structural scope check and left
  content nested in them unjudged. Admission now covers whole subtrees.
- A later revision did not admit a body the recipe never represents. Generated
  bodies are now always admitted.

## Verification of the first delivery

Independent verification by `codex-hydrated-type-design` (board claim `bittype-source-traversal-20261108`, record 000003) rejected #191 at
`49114c5`. It found criteria 1 and 4 of GEN-01 qualified, and three defects:

1. **Evidence (blocking).** `PlanningInput` carried the capture and an already
   opened source as independent fields. Two captures shared one snapshot and
   graph and differed only in a second report. A source opened over the holding
   report planned the rejecting capture, and the output was valid. Checking
   graph identity cannot tell these evidence sets apart.
2. **Selectors (blocking).** A selector whose language profile differed, with
   the same key, still chose the local descriptor binding, Cell binding and Ignore
   recipe. Registration checked mapping selectors, but planning compared keys
   only.
3. **Diagnostic site.** The index kept the first syntax site of a term. A failure
   while planning Borrow through the term it shares with Lend pointed to Lend.

[#192](https://github.com/Bitspark/bittype/pull/192) corrects each defect:

1. `Plan` takes the input and the source provider, admits the evidence and opens
   the provider itself. Handles resolve only against their exact capture.
2. A source selector outside the capture's source language fails as `ErrScope`.
3. A term recorded at several sites is reported at the occurrence's declaration.

The reviewer's regressions are part of the suite: `TestKernelPlansExactlyTheAdmittedEvidence`,
`TestKernelSelectorsAreQualified` and `TestKernelLocatesSharedTermsAtTheirDeclaration`.
Each has a correctly composed control.

An independent code review of the correction before merge found that two callee
failures were still reported at the callee's site. It also found that the docs
claimed the kernel guarantees the provider's interpretation, that the site index
had become quadratic, and that the registry and the exported planner ordered
description and evidence errors differently. All are fixed in the same PR, which
also tests every part of a selector's language identity. A second review of
the fixes found tests that could not tell a use's own term site from its
declaration, or the declaration fallback from the last recorded site, and an
untested check order and horos handle checks. Those tests are added, and the
remaining overclaims are corrected. The case stays open until the correction
is independently verified.

Each kernel check was seeded with a mutant that turned its test red. Covered
checks:

- admission of actuals, of operands and of unrepresented bodies;
- binder ownership;
- child captures through `Mapped` and `Pattern.Capture`;
- body-only recipes;
- the undetermined/rejected distinction;
- the structural application check;
- the declaration judgment;
- the zero verdict;
- diagnostic class isolation;
- located horos diagnostics;
- handle resolution against the exact capture, and a source opened over a
  second admission of the same input;
- selector qualification, for declaration and expression selectors, and over
  the whole language identity rather than its profile or identity alone;
- reporting at a single recorded term site, at the declaration for a term with
  several (whichever is recorded last), with repeated records at one site as one
  site, and at the use for a callee or selector failure;
- description checks before evidence admission, on both paths;
- the horos reader's handle checks.
[Traversal qualification](../../integration/go/generation_traversal_test.go) ·
[Planning tests](../../integration/go/generation_planning_test.go)

## Limits

**Scope.** Only direct binder ownership is supported; there are no nested generic
scopes. A language that needs nested scopes must contribute its own scope
judgment, and the kernel infers no `Owner`-chain rule.

**Structural reading.** Applications need a direct declaration callee and complete
arguments. Type-level constructor values still need an explicit representation.

**Operands.** Every provider-expression operand is admitted, but represented
only when a mapping captures it.

**Identity.** Effective-subject equivalence is not defined.

**Diagnostic sites.** The graph does not relate syntax sites to declarations. A
shared term recorded at only one site is therefore reported there, even for an
occurrence under another declaration.

**Provider integrity.** The kernel opens the provider over its capture and
resolves handles only against it. It cannot check where a provider's answers
come from. A provider that serves a cached interpretation unchanged cannot
resolve the run's handles: `In` reports false, and a source that checks it, as
horos and the independent source do, fails with `ErrScope`. A provider that rebinds a
cached interpretation to the new capture while keeping verdicts from earlier
evidence is a defect of that provider, and the kernel does not detect it.

**horos judgments.** horos reads no membership facts, so its only non-admitted
judgments are Unsupported and Undetermined.

**Out of scope.** Recursive generation (bittype#164) and native enforcement of
source bounds ([GEN-02](../generation-gaps.md#gen-02)) are unchanged. Declared
representation-profile erasure is not evidence about GEN-01, and this correction
does not implement native constraint enforcement.
