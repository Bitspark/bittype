# bittype/ACX-011: Common Files semantic orchestration is implemented per language

**Status.** Identified on 13 November 2026 at main
`a28febc44830d49ec3c2f22cd70c7ded9d0e5028`. Implementation correction pending.

## Ownership failure

The Horos and TypeScript Files implementations each construct a semantic capture
from concrete units, invoke the supplied resolver and analyzer in order, and
admit the resulting evidence. Both already use the common snapshot, unit,
resolution, analysis and receipt contracts. Nevertheless, this orchestration
is implemented separately inside two language integrations.

These are the implemented entry-point signatures; `s.Language` names the
respective language in each package:

```go
// internal/horosfiles
func (p Stages) Establish(ctx context.Context, units []source.Unit[s.Language],
    entry files.ID, limits semantic.Limits) (source.Publication[s.Language], error)

// format/files/typescript/provider/go
func (p Stages) Capture(ctx context.Context, units []source.Unit[s.Language],
    entry files.ID, limits semantic.Limits) (f.FilesInput[s.Language], error)
```

The first is called by Horos Files parsing and formatting. The second is a public
entry point for existing TypeScript files; its private `establish` method is
also used when the formatter checks newly placed files. Consumers explicitly
supply the collaborators in each package's `Stages` value. The relevant
[Horos implementation](https://github.com/Bitspark/bittype/blob/a28febc44830d49ec3c2f22cd70c7ded9d0e5028/internal/horosfiles/stages.go#L129)
and [TypeScript implementation](https://github.com/Bitspark/bittype/blob/a28febc44830d49ec3c2f22cd70c7ded9d0e5028/format/files/typescript/provider/go/format.go#L59)
perform the same invariant sequence:

1. Bound the units' nodes and printed bytes and retain their file/tree links.
2. Capture them using the supplied checker.
3. Resolve that exact snapshot with the supplied resolver and options.
4. Analyze it with that resolution result and the supplied analyzer and options.
5. Admit the requested semantic evidence against the same snapshot.

This finding follows the responsibility and data dependencies, not textual
similarity. The sequence requires no Horos descriptor, module pin, TypeScript
import syntax or native typing rule. Its cancellation, resource accounting,
snapshot association and evidence guarantees are common. Maintaining those
guarantees currently involves both implementations; another Files integration
would otherwise need its own assembly of the same mechanism.

## What remains language-owned

The surrounding work is genuinely different and must remain so:

- Horos selects its entry module, interprets module pins and the module closure,
  derives module/import groups and relations, and checks canonical module
  agreement. Its profile requirements remain explicit.
- TypeScript interprets its compiler import environment, constructs native
  modules and import/export syntax, and checks the resulting bindings.
- Each composition chooses the required evidence, configured implementations
  and options. The common mechanism must not choose a language or weaken a
  requested profile.

The agnostic [Files contracts and Source publisher](../source.md) already
protect publication association. That does not make the separately implemented
sequence agnostic. Nor is `Stages.CheckPublication` wholly generic:
its Horos module-closure and canonical-content checks require Horos knowledge.

## Rules and required correction

[AC-11](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-11)
places the invariant mechanism in the agnostic layer from its first
implementation; [AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02)
keeps language meanings with their owners.
[AC-05](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-05)
requires preservation of the configured contracts, limits, evidence association,
failure categories and publication guarantees during correction.

Provide one agnostic implementation of the explicitly requested semantic
sequence, with typed units, supplied stages, options and evidence requirements.
Have both Files integrations use it and retain their own organization and
meaning checks. This does not require every language to use every optional stage,
or dictate a universal analysis order beyond the operation being composed.
The replacement public contract is not specified by this audit; derive it from
these responsibilities under AC-09, not by copying either current `Stages`.

## Qualification and limits

Source inspection only; no production behavior changed and no new execution
test was run. The audit does not establish an unsound receipt or accepted
incorrect program. Correction must qualify the same common sequence with both
language integrations, including configured evidence, foreign receipts,
cancellation and resource/error behavior, while preserving each language's
organization and binding checks.

[Abstraction levels](../abstraction-levels.md) ·
[Package boundaries](../package-boundaries.md) ·
[Case index](README.md)
