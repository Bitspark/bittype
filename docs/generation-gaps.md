# Unfinished generation boundaries

**Recorded 8 November 2026.** The declarative generation change
([PR #155](https://github.com/Bitspark/bittype/pull/155)) separates semantic
planning from native AST construction. The kernel-owned semantic traversal
corrects the source responsibility split (GEN-01). Native enforcement of source
kind constraints (GEN-02) remains open.

Keep these references stable and update their status and completion criteria in
place when implementation changes. A successful `Page<Cell>` example or an unrelated
generation improvement does not close either item.

| Reference | Status | Classification |
| --- | --- | --- |
| [GEN-01](#gen-01) | Corrected 8 November 2026: kernel-owned traversal with source judgments | Architectural responsibility split; recorded as [bittype/ACX-006](architecture-cases/0006-generation-source-traversal.md) |
| [GEN-02](#gen-02) | Open; no native enforcement route implemented | Unsupported guarantee with missing native-plan expressiveness and target encodings |

The implemented [mapping SPI](generation-spi.md) resolves representation choices
before construction. The native constructor receives a closed program with no
source reader, mapping callbacks or source-substitution templates. TypeScript
Files consumes its checked output and verifies bindings after placement. These
delivered boundaries remain distinct from the open item GEN-02.
The remaining Rust/C++ Files ownership defects retain their own
[ACX-004 record](architecture-cases/0004-generation-file-syntax.md).

## GEN-01

**Common semantic navigation is kernel-owned. Corrected 8 November 2026.**

### What was wrong

The planner used to ask each source integration for navigation of the
established semantic graph as well as source interpretation. horos's `Open`
enumerated a declaration's parameters and body. Its `Admit` followed references
and applications and paired actuals with callee binders, alongside horos sort
and scope checks. `CheckChild` re-derived the child relation. The independent
test source reproduced the same navigation from hand-built tables. Each new
source integration would have had to re-implement traversal of facts that the
shared graph already records. [bittype/ACX-006](architecture-cases/0006-generation-source-traversal.md)
records the defect. [Research 0013](../research-docs/0013-generation-source-traversal.md)
derived the replacement.

### Delivered division

The kernel and the source now divide the work as follows.
[SPI](generation-spi.md#semantic-traversal-and-source-interpretation)

**Kernel**, in `contracts/generation/go`:

1. It admits the request and unpacks the one graph that the designated evidence
   publishes.
2. It resolves root sites to their declaring entities.
3. It opens declarations from their body and real binders.
4. It follows references, and pairs application arguments with real binders
   through `DirectParameterBindings`.
5. It mints occurrence, declaration and binder handles for one run, and checks
   that every capture is a structural child of the current occurrence.
6. It enforces direct-binder scope.

**Source**, through `Interpretation` and its own reader:

- classifies the entities it generates and gives each a stable selector;
- names its provider expressions;
- judges declarations and whole occurrences from the captured evidence, with the
  verdict Admitted, Unsupported, Rejected or Undetermined;
- interprets payloads and operand meanings for its own mapping libraries.

For `People = Page<Cell<Person>>`:

| Decision | Owner |
| --- | --- |
| Follow the references to Page, Cell and Person | Kernel, over the captured graph |
| Follow `Page.T` to `Cell<Person>` and `Cell.V` to Person | Kernel, over `DirectParameterBindings` |
| Interpret horos constructors, proper sorts and module-pinned identity | horos interpretation |
| Represent Cell natively, selecting its actual by Cell's binder | Cell's mapping library |
| Admit evidence, preserve scope, and judge every actual | Kernel mechanism with source judgments |

### Completion criteria

1. The kernel consumes declaration, reference and application structure through
   the shared graph. No source supplies a traversal:
   - horos's interpretation reads only facts and payloads;
   - the independent source publishes an ordinary graph through its own analysis.
2. The source-owned contract is narrow. It covers classification and selectors,
   descriptor selectors, and judgments. Undetermined, meaning uncovered evidence,
   fails as `semantic.ErrEvidenceUnavailable` and Rejected as `ErrRejected`.
   horos reports a missing sort fact as Undetermined; it is no longer a malformed
   request. There is no hidden analysis.
3. These guarantees stay enforced:
   - capture identity, through `AdmitInput`;
   - real binders and source environments, through occurrences that carry their
     declaration;
   - foreign-reference rejection, through handles valid only in their run and
     graph;
   - limits and whole-output failure.

   Every actual of a declaration use is judged before any binding or convention
   represents it, including actuals the convention omits.
4. horos and the independent source drive the same walk. The independent source
   covers nested `Page<Page<T>>` and a library-owned `Cell<Person>`. horos covers
   `Page<Cell<Person>>`. The negative cases cover both sources:
   - omitted actuals: refuted, uncovered and admitted;
   - binders of another declaration, including through a shared term;
   - non-child captures, and handles kept from another run;
   - recursion, parameter-headed and unapplied applications;
   - foreign evidence;
   - horos applied constructor aliases and curried forms, now Unsupported.
5. The superseded surface is removed with no aliases:
   - `Source.Open`, `Admit`, `CheckChild`, `Authenticate`, `CheckInput`, `Owners`
     and `RootEnvironment`;
   - `Catalog`, `SourceSession`, `ResolveRoots`, `OpenData` and `UseData`;
   - the source address spaces and the effective-subject and formal-subject uses;
   - the horos pre-walk and its reader navigation.

[Qualification](../integration/go/generation_traversal_test.go) ·
[Independent source](../integration/go/testdata/generationfixture/fixture.go)

### Supported scope and limits

- **Scope** is direct binder ownership by the declaration being planned. Nested
  generic scopes are not supported. A language that needs them must contribute its
  own scope judgment; the kernel adds no `Owner`-chain rule.
- **No structural view.** The following are unsupported at their site:
  - applications without a direct declaration callee and complete arguments;
  - a generic declaration referenced without arguments.

  Type-level constructor values still need an explicit representation.
- **Operands.** Provider-expression operands are visited only when a mapping
  captures them. The source's judgment of the whole occurrence covers its
  operands. For horos, the analysis checks every operand's formation contract.
- **Memoization.** Occurrences are distinguished by path, and no effective-subject
  equivalence is defined. Judgments are memoized per occurrence and declarations
  per entity.
- **Rejection.** horos reads no membership facts, so it has no Rejected path
  today. Unsupported and Undetermined are its outcomes.
- **Out of scope.** Recursive generation remains unsupported
  ([#164](https://github.com/Bitspark/bittype/issues/164)). Native enforcement of
  source bounds is [GEN-02](#gen-02).

## GEN-02

**Preserve source constraints for consumers writing native applications.**

### Current guarantee and visible limitation

Horos checks source applications before generation. The supplied
`pagination-representation/1` profile explicitly permits erasing proper-kind
constraints from native generic parameters. The declared
`pagination-enforcement/1` profile is unsupported: the mapping library rejects
registration with `generation.ErrUnsupported`. It never falls back to the
representation profile. [Profile implementation](../transform/syntax/horos/mapping/go/mapping.go)

For example, with the ordinary Horos prelude imported as h:

```horos
type TextBox<T: h.Text> = { value: T; };
```

The representation output has this equivalent TypeScript shape:

```typescript
interface TextBox<T> {
  value: T;
}
const box: TextBox<boolean> = { value: true };
```

The native application is accepted even though the corresponding boolean
application violates the Horos bound. The source checker has not disappeared;
its restriction is absent from the generated native API.

The current construction data contract is:

```go
type NativeParameter struct {
    Name   string
    Origin ast.NodeRef
}
```

The planner supplies these records and native construction consumes them.
There is no parameter-bound or constraint field.
[Native declaration contract](../contracts/generation/go/native.go)

The implemented generic convention can add, omit or reorder native arguments:

```go
type GenericConvention[T any] struct {
    Parameters []ParameterDescription
    Bindings   []Expression[T]
    Arguments  []Expression[T]
}
```

A mapping supplies this convention and the planner composes it. It does not
provide native bound syntax, evidence issuance or an admission theorem by
itself. [Convention contract](../contracts/generation/go/description.go)

### Remaining design and completion criteria

This is more than supplying four additional mapping functions. The native plan
must be able to describe the chosen target constraints or descriptor/evidence
encoding, and target construction must realize it. A justified general contract
extension may be needed; that is distinct from installing special cases for
Data, Text or another named kind. The enforcement SPI is not finalized.

Completion for an advertised target/profile requires:

1. State its supported source/kind domain, native argument forms, consumer
   operations, checking point, dependency/issuer assumptions and alias/wrapper
   promise. Track target-specific coverage explicitly; one target's success
   does not establish all four.
2. Define the native intent and construction contract from those requirements.
   Horos retains source judgment meaning; mapping libraries own the chosen
   encoding; native integrations own their language's forms. The kernel does
   not interpret named kinds or introduce a universal proof tuple.
3. Preserve the intended source-to-native association, parameter assumptions
   and result obligations through open and nested applications, including
   source arguments omitted from the payload representation.
4. Qualify accepted and rejected native consumers at the stated checking point,
   and establish the guarantee over the stated domain. Compiler acceptance,
   conditional mapping proofs and source membership answer different questions;
   a few successful examples do not prove unrestricted native admission.
5. Reject unsupported enforcement requests without weaker output. Update this
   status only for the implemented and qualified profiles.

The [native constraint proposal](native-kind-constraints.md),
[generic admission discussion](generic-elaboration.md) and
[Go admission proposal](go-native-admission.md) contain candidate encodings and
known pitfalls. A direct bound such as `T extends string` is not a general
translation of Horos kinds: a kind union and a native union classify different
things. Extra evidence arguments, private names and nominal wrappers likewise
need an actual admission argument.

`Page<Cell<Person>>` currently qualifies representation composition and binding
preservation. Unrestricted `Page<T>` introduces no Data constraint. The fixture
does not prove that handwritten native applications preserve a constrained
DataPage or TextBox contract. Runtime value validation and adapters are outside
both open items.
