# bittype/ACX-004 — Generation emits file syntax before placement

**TypeScript corrected on 8 November 2026. The C++, Rust and Go cases remain open.** The [generation/Files decision](../source.md#from-generated-declarations-to-files)
assigns concrete import and export syntax to the language's Files formatter.

**Original behavior, 7 November.** The Cell mapping correctly requests the public type
`@bitspark/cell` exports as `Cell`. The [TypeScript target](https://github.com/Bitspark/bittype/blob/dbb7ff2/transform/syntax/typescript/target/go/builder.go)
then calls `imports.Emit` while assembling the logical AST, allocating a local
import binding before the consumer has placed a declaration. Its
[declaration builder](https://github.com/Bitspark/bittype/blob/dbb7ff2/transform/syntax/typescript/target/go/symbols.go)
also marks each reserved interface or alias `exported: true`. Later, the
[Files formatter](../../format/files/typescript/provider/go/relocate.go)
discards the old imports, allocates bindings per output unit and consults
export visibility when a use crosses units. The logical unit's syntax is
doing double duty as the record of which dependency was intended.

The same distinction applies outside TypeScript. The [C++ target](https://github.com/Bitspark/bittype/blob/dbb7ff2/transform/syntax/cpp/target/go/assemble.go)
materializes registered headers as include nodes during declaration assembly;
knowing the required public header is distinct from placing its directive in
an output unit. The [Rust target](https://github.com/Bitspark/bittype/blob/dbb7ff2/transform/syntax/rust/target/go/symbols.go)
sets module visibility to public on its reserved structs and aliases. These
are further instances to address, not evidence that native Files placement
has been implemented for those targets.

**Rules.** [AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02)
assigns representation choice to the mapping and concrete module spelling to
Files. One local alias or export modifier is a decision about a module, not a
property of the type being generated. [AC-06](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-06)
applies because generating a single convenient TypeScript module made its
organization a prerequisite for later placements. The successful
[split and co-located examples](../../integration/go/typescript_placement_test.go)
qualify the final imports; they do not establish that the two stages own
different decisions.

**Required correction and qualification.** Preserve each use's target and
publication intent through target semantic admission without using a synthetic
module import or export flag to convey them. Apply one consumer placement in
Files. Generate `Cells = Page<Cell<Person>>` once, then format it twice:
moving `Person` into the file containing `Cells` must remove only the relative
import. Both outputs must preserve bindings and intended visibility, including
in the presence of a conflicting local name; an inaccessible target must fail
instead of being rebound or silently exported. The [planner](../generation-planner.md#where-the-plan-stops)
describes this qualification.

**Implemented TypeScript correction.** Generation retains exposed sites and emitted external use occurrences. Files allocates imports and aliases after placement and checks the compiler's named exports and each relocated destination. The qualification uses one output with a local Cell and an external use of Cell. With Person in a second unit, the relative import is present. Moving Person alongside Page removes its relative import and retains the package import with the alias `Cell_1`. Missing or inaccessible destinations fail publication. [Output](../../contracts/generation/go/output.go) · [Generated Files input](../../format/files/typescript/provider/go/generated.go) · [Placement tests](../../integration/go/typescript_placement_test.go)

**Further instances identified on 10 November 2026. Pending; no correction is implemented.** The [Go target](https://github.com/Bitspark/bittype/blob/33a3287/transform/syntax/golang/target/go/provider.go#L79) writes one package clause, taken from target configuration, on the single file node it constructs. The [C++ target](https://github.com/Bitspark/bittype/blob/33a3287/transform/syntax/cpp/target/go/provider.go#L159) wraps every declaration in one configured namespace. The [placement proposal](../native-file-placement.md#required-checks-before-claiming-support) takes the native scope of each emitted declaration from the consumer's placement. It requires one generated output to be placed both across packages or namespaces and co-located. That is impossible while generation fixes a single package or namespace, and a Files stage that overwrote it would repeat the double duty described above. These are the same AC-02 and AC-06 failures as the import and visibility cases. [Kernel composition](../kernel-composition.md) documents Package and Namespace as target configuration; that describes the present behaviour, not the corrected ownership. **Required correction:** generation carries no package or namespace choice, and Go Files and C++ Files read it from placement. Tracked in [#161](https://github.com/Bitspark/bittype/issues/161) and [#160](https://github.com/Bitspark/bittype/issues/160), which record the decision and its reasons.

**Rust external paths, identified 10 November 2026. Pending.** The [Rust target](https://github.com/Bitspark/bittype/blob/ca5dc80/transform/syntax/rust/target/go/provider.go#L118-L119) spells an external library type as an absolute path inside the declaration and reports no external use. An absolute path stays valid wherever the declaration is placed, so no placement is broken today. The [placement proposal](../native-file-placement.md#rust) nevertheless assigns the import of an external type to Rust Files, and the generated output records no intent that Files could check. Tracked in [#159](https://github.com/Bitspark/bittype/issues/159).

**Order of correction, 11 November 2026. Pending; nothing in this section is implemented.** The owner directed: "if something is placed incorrectly, it is incorrect and must change. Something that is missing can be healed additively, i.e. there is nothing wrong requiring removal." And: "no compromise on architecture, but practical deferral is allowed, if the huge effort justifies it." ([AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)). Each open instance is therefore corrected in two steps. The first step moves the misplaced decision to the language's generated-Files formatter, for the placements that existing semantic evidence can check, and refuses every other placement. The second step adds coverage and changes nothing the first delivered. The first step cannot skip its check: [`NewGeneratedFilesFormatter`](../../contracts/format/go/generated.go) refuses a formatter that has no check or no semantic evidence.

| Instance | First step | State |
| --- | --- | --- |
| Go package clause | Generation decides no package. Go Files writes the clause for one package and no import. [#161](https://github.com/Bitspark/bittype/issues/161), slice 1 | Can start. It does not wait for imported-package semantics |
| Rust visibility | Generation chooses no visibility. Rust Files decides it for one module. [#159](https://github.com/Bitspark/bittype/issues/159), slice 1 | Can start. It does not wait for multi-file crates |
| Rust external paths | The target reports an external use and Rust Files decides the import | Deferred until the first item of [#170](https://github.com/Bitspark/bittype/issues/170) supplies another crate's semantics |
| C++ includes and namespace | Generation emits neither. C++ Files decides both for one header and one namespace. [#160](https://github.com/Bitspark/bittype/issues/160) | Deferred until [#167](https://github.com/Bitspark/bittype/issues/167) |

**Why C++ is deferred.** C++ has no resolution or analysis provider, so a C++ generated-Files formatter cannot be constructed under the shared contract. Removing the includes and the namespace without one would leave the generated header unable to compile. Accepting a formatter without semantic evidence would weaken a guarantee the shared mechanism owns. This is a deferral for effort with a stated trigger. It is not an accepted design, and both C++ instances stay open violations.

The C++ header directives and Rust module visibility described above require separate target and Files work. The declarative planner and mapping redesign are implemented; they do not resolve these separate Files ownership cases.
