# bittype/ACX-012: Verification-law role classification fixes one language pair

**Status.** Identified on 13 November 2026 at main
`a28febc44830d49ec3c2f22cd70c7ded9d0e5028`. Implementation correction pending.
This case concerns the shared architecture test, not production rule selection.

## Ownership failure

The architecture test describes verification laws as independently supplied
contributions. Its common role classifier nevertheless recognizes only one
literal language-pair path. This is the implemented function:

```go
func isVerificationLaw(relative string) bool {
    s := strings.Split(relative, "/")
    return len(s) == 6 &&
        strings.HasPrefix(relative, "verification/files/horos-to-go/mappings/") &&
        s[5] == "go"
}
```

[`architecture_roles_test.go`](https://github.com/Bitspark/bittype/blob/a28febc44830d49ec3c2f22cd70c7ded9d0e5028/integration/go/architecture_roles_test.go#L149)
implements it. `packageRole` calls it to select the law role and its
no-law-to-law-import rule. The general architecture checker invokes
`packageRole`; unclassified production packages fail. A corresponding law
under `verification/files/typescript-to-rust/mappings/example/go` receives no
role from this branch. Classifying it requires editing shared test machinery
or adding a hand entry.

The current comment justifies the fixed pair because it is the only one and
the directory says `go` where the language identity is `golang`. Neither
reason makes contribution-role recognition a Horos semantic responsibility.
The same file already derives generation-mapping roles from
`<source>-to-<native>` without fixing either language.

## Boundary and required correction

Under [AC-11](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-11)
the invariant role mechanism must be agnostic. Under
[AC-03](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-03)
adding another contribution within this layout boundary must not require
another hard-coded language case in that mechanism.

Derive verification-law roles from the declared pair and contribution shape.
Keep pair-owned semantic dependencies explicit and scoped; this is not
authorization to let every law import every language or model package.
Separate any necessary path-to-language vocabulary from the generic classifier.
Preserve the prohibition on importing another law and the existing restrictions
on selected implementations. Correct the existing test; add no new heuristic
gate or scoring system.

The public Horos/Go law contract and its concrete projections legitimately refer
to those languages. This finding does not require a new production verification
pair or establish that all of its orchestration is correctly placed.
[ACX-001](0001-verification-composition.md) remains the record of the qualified
mapping-extension boundary within Horos/Go.

## Qualification and limits

Source inspection establishes the literal path restriction. No seeded package
or new execution test was run in this documentation audit. A correction should
demonstrate independent pairs receiving the same role mechanism and forbidden
cross-contribution imports still failing. The earlier Files coverage restriction
was a similar test-model failure and remains corrected:
[ACX-003](0003-files-stage-coverage.md).

[Package boundaries](../package-boundaries.md) · [Case index](README.md)
