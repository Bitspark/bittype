# bittype/ACX-004 — Generation emits file syntax before placement

**Status: open, 7 November 2026.** The [generation/Files decision](../source.md#from-generated-declarations-to-files)
assigns concrete import and export syntax to the language's Files formatter.

**Observed behavior.** The Cell mapping correctly requests the public type
`@bitspark/cell` exports as `Cell`. The [TypeScript target](../../transform/syntax/typescript/target/go/builder.go)
then calls `imports.Emit` while assembling the logical AST, allocating a local
import binding before the consumer has placed a declaration. Its
[declaration builder](../../transform/syntax/typescript/target/go/symbols.go)
also marks each reserved interface or alias `exported: true`. Later, the
[Files formatter](../../format/files/typescript/provider/go/relocate.go)
discards the old imports, allocates bindings per output unit and consults
export visibility when a use crosses units. The logical unit's syntax is
doing double duty as the record of which dependency was intended.

The same distinction applies outside TypeScript. The [C++ target](../../transform/syntax/cpp/target/go/assemble.go)
materializes registered headers as include nodes during declaration assembly;
knowing the required public header is distinct from placing its directive in
an output unit. The [Rust target](../../transform/syntax/rust/target/go/symbols.go)
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
describes this qualification. No implementation repair is claimed here.
