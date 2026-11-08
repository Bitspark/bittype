# Unfinished generation boundaries

**Open work, recorded 8 November 2026.** The declarative generation change
([PR #155](https://github.com/Bitspark/bittype/pull/155)) separates semantic
planning from native AST construction. It does not complete the source-reader
responsibility split or native enforcement of source kind constraints.

Keep these references stable and update their status and completion criteria in
place when implementation changes. A successful `Page<Cell>` example or an unrelated
generation improvement does not close either item.

| Reference | Status | Classification |
| --- | --- | --- |
| [GEN-01](#gen-01) | Open; replacement SPI not finalized | Architectural responsibility split; recorded as [bittype/ACX-006](architecture-cases/0006-generation-source-traversal.md) |
| [GEN-02](#gen-02) | Open; no native enforcement route implemented | Unsupported guarantee with missing native-plan expressiveness and target encodings |

The implemented [mapping SPI](generation-spi.md) resolves representation choices
before construction. The native constructor receives a closed program with no
source reader, mapping callbacks or source-substitution templates. TypeScript
Files consumes its checked output and verifies bindings after placement. These
delivered boundaries remain distinct from the open items below.
The remaining Rust/C++ Files ownership defects retain their own
[ACX-004 record](architecture-cases/0004-generation-file-syntax.md).

## GEN-01

**Move common semantic navigation from source integrations into the kernel.**

### Current contract and concrete problem

The following are implemented method excerpts, not the proposed replacement.
Horos's generation integration implements them; the planner invokes them.
[Source](../contracts/generation/go/source.go) and
[Catalog](../contracts/generation/go/mapping.go) define the full interfaces.

```go
// Source[F] excerpts
Open(context.Context, Declaration[F]) (OpenData[F], error)
Admit(context.Context, Environment[F], Use[F]) (UseData[F], error)
CheckChild(context.Context, Environment[F], Use[F], Use[F]) error

type Catalog[F any] interface {
    Declarations() []Declaration[F]
    DeclarationSelector(Declaration[F]) (Selector, bool, error)
    UseSelector(UseData[F]) (Selector, error)
}
```

The reader currently supplies both navigation of established semantic structure
and interpretation of source meanings. In particular, Horos's `Open` enumerates
the declaration's parameters and body; Admit follows references/applications
and pairs actuals with callee formals, alongside Horos-specific sort and scope
checks. [Current implementation](../transform/syntax/horos/source/go/source.go)

The shared semantic model already represents entities, references, applications
and argument bindings. One relevant implemented query is:

```go
func (m Semantics[L]) DirectParameterBindings(id TermID) ([]ArgumentBinding, bool)
```

The shared model implements this query; a common traversal can consume it.
It pairs direct applications structurally. It does not prove applicability,
interpret kinds, handle every possible application form, or authorize a use
under arbitrary assumptions. [Query and limits](../semantics/go/model.go)

For `Page<Cell<Person>>`, the intended division is:

| Decision | Intended owner |
| --- | --- |
| Follow established references to Page, Cell and Person | Kernel over the captured semantic graph |
| Follow `Page.T` to `Cell<Person>` and `Cell.V` to Person | Kernel over established application bindings |
| Interpret Horos constructors, kind judgments and pinned identity | Horos integration |
| Select a native representation of Cell | Cell's mapping library |
| Admit source evidence and preserve capture/scope association | Shared contract guarantees with source-owned judgments |

This is an unfinished architectural correction, not an optional simplification
or a missing output feature. The kernel has no Horos-specific dispatch, but
making each source integration reproduce common graph navigation is still the
wrong responsibility boundary. Source checking exists today; it must survive
the correction.

### Remaining design and completion criteria

The narrower SPI is not finalized. Define which operations are graph mechanics
and which questions require source-owned interpretation. Do not preserve the
existing method names, handles or division of work merely because they exist.
Do not move `Admit` wholesale into the kernel: its current responsibilities need
to be separated, including the evidence and assumption checks.

Completion requires:

1. The kernel consumes established declaration/reference/application structure
   through the shared semantic context. A source integration no longer has to
   supply a parallel traversal of those same facts.
2. Narrow source-owned contracts retain constructor/operand interpretation,
   stable selector identity and applicability judgments. Unknown or uncovered
   facts remain distinct from failed judgments; there is no hidden analysis.
3. Capture identity, real binders, source environments, foreign-reference
   rejection, limits and whole-output failure remain enforced. Omitted native
   arguments still require source admission.
4. Horos and an independent source exercise the same traversal, including nested
   `Page<Cell<Person>>` composition and invalid scope/evidence cases. Inspection
   of the change establishes ownership; tests alone do not establish it.
5. Superseded reader interfaces and duplicated traversal are removed. Update
   this item and ACX-006 with the correction and its actual supported scope.

The existing independent-source fixture qualifies today's public composition;
it still implements the broad reader and therefore does not close GEN-01.
[Fixture](../integration/go/testdata/generationfixture/fixture.go) ·
[Context matrix](abstraction-levels.md#agnostic-kernel-and-language-specific-rules)

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
