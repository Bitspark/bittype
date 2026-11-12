# Advice for consult
**Research doc:** doc_430dcd8b-84f8-4a7d-8c1e-a9352fe5bcc7.md
**Document:** doc_430dcd8b-84f8-4a7d-8c1e-a9352fe5bcc7
**Advisor:** ChatGPT, model not read from the answer page
**Session:** nsess_25dea434-37fe-4bf4-a588-be22a43d3f05
**Run:** run_cd247478-16f7-4b16-8b65-52be1f89ceb6
**Chat:** https://chatgpt.com/c/6af62fa6-d274-83ed-ba52-351bf0dc7a86
**Verify:** nightfall consult verify run_cd247478-16f7-4b16-8b65-52be1f89ceb6 --output <this file>

**I would choose S2 for placement, U2 for collisions, A generalized to every target for unplaced syntax, P2 for the Go package name, and F1 for Rust fields.** I would also make the evidence requirement depend on the admitted placement and require an explicit environment from the first Rust slice.

The accepted separation between generation and Files is sound. Several assumptions around it need changing before the proposed slices can meet your stability test.

The most consequential issue is that **the literal freeze conflicts with the deferred Rust external-intent and C++ generation corrections**. Those corrections change existing generated output in precisely the ways your Question section counts as failures. I would keep your freeze and normalize those generation outputs now, while leaving unsupported placement checks unavailable.

I use your section numbers below for claims about bittype. The external citations support language and compiler facts; this is a design review of the supplied baseline excerpts, not a verification of the repository.

## 1. Placement: S2, with shared admission of structural invariants

**Choose S2: language-owned public placement types, backed by one shared admission mechanism.** The common invariants justify shared validation. They do not establish that the consumer-facing structure should be identical.

The proposed `Placement[Scope]` assumes that a complete layout is a collection of declaration–file–scope assignments. That is a useful *normalized view*, but Rust already needs more information to construct that view.

### What the shared mechanism should own

I would give the shared mechanism a normalized view with this conceptual shape:

```text
PlacementView
    InputCapture
    Files
        FileID -> RelativePath
    Definitions
        DeclarationSite -> FileID
```

The language owns the associated scope information and compilation structure. The shared mechanism establishes:

- Every assigned declaration belongs to the identified sealed output.
- Every generated placement unit has exactly one definition assignment.
- File identities and paths satisfy the shared file contract.
- Assignments reference existing files.
- Resource limits hold.
- The resulting publication corresponds to that admitted input.

An opaque scope identity may accompany this view for association or diagnostics. The shared validator should not interpret it or perform a generic `(scope, name)` collision check.

Make the coverage operation generic in the assigned item identity if that usefully serves horos as well. A module pin and a native declaration site can both participate in “each expected item exactly once” without pretending they have the same domain meaning.

This satisfies AC-02 and AC-03: one owner implements the common invariant, while independently owned contributions supply their interpretations.

### What each language’s placement should contain

| Language | Consumer-facing organization |
|---|---|
| TypeScript | Files and their declaration assignments, under an explicitly established module interpretation |
| Go | Packages, their package names, their member files, and declaration assignments to those files |
| Rust | Crates and module trees, the physical realization of each module body, and declaration assignments to modules |
| C++ headers | Files and declaration assignments to namespace descriptors, with namespaces independent of files |

These should remain **one placement object per request**. Having a file table and a module table inside that object does not create the forbidden second table of imports.

#### Go: state the package name once per package

**Choose package-level naming over repeating the name in every file assignment.** Package identity and package membership deserve explicit representation.

For example:

```text
GoPlacement
    InputCapture
    Files: FileID -> OutputRelativePath
    Packages:
        PackageID
        Name
        MemberFiles
    Definitions: DeclarationSite -> FileID
```

Here `PackageID` is a placement-local identity. The native import path is established through the environment’s layout mapping.

There is no need to duplicate each file’s directory in its declaration assignment. Validate that the member files of a package occupy the required directory. An empty package can be refused; avoiding that invalid state is not a reason to repeat package names throughout the input.

**Derive import paths only after making the coordinate system explicit.** Go’s module rule relates a package path to its directory *relative to the module root*. An output-relative directory is sufficient only when the environment explicitly relates the output root to that module root. [The Go Programming Language](https://go.dev/ref/mod?utm_source=chatgpt.com)

For example, output path `pages/types.go` can mean either:

```text
<module root>/pages/types.go
<module root>/generated/pages/types.go
```

Those imply different package paths. Put the authoritative root mapping in the environment; do not silently assume the first interpretation.

#### Rust: represent module bodies and their storage

**Choose a module tree with a body-realization choice.** Arbitrary pairs of file and module path are insufficient.

Rust compiles a crate from a root source file. A source file describes a module and may contain inline descendant modules. An out-of-line module declaration loads another file. These relationships are part of the language’s compilation structure. [The Rust Reference](https://doc.rust-lang.org/reference/crates-and-source-files.html?utm_source=chatgpt.com)

A useful placement can express:

```text
Files
    root   -> src/lib.rs
    people -> src/people.rs

Modules
    []           -> SourceFile(root)
    ["pages"]    -> Inline
    ["people"]   -> SourceFile(people)

Definitions
    Page   -> ["pages"]
    Cells  -> ["pages"]
    Person -> ["people"]
```

`Inline` means inline in the parent module’s body; its physical file follows from that relationship. Moving `people` inline changes one body-realization choice.

This also represents a necessary case that the current nonempty-assignment rule excludes: **a crate-root file containing only module scaffolding, with every generated type in child modules**.

Slice 1 can admit only the root module in one source file. The rest of this structure can already have a defined meaning and return `Unsupported`.

### Assignment order and entry

**Choose unordered placement membership, with deterministic emission order owned by Files.**

Assignment enumeration should not implicitly promise source order. Files must be able to schedule declarations and supporting syntax according to the target language. Representation-owned order, such as field order or type-parameter order, remains part of the generated declaration.

If consumer-controlled declaration ordering belongs to your intended domain, model it explicitly as ordering constraints now. Otherwise, do not accidentally acquire that obligation from JSON array order.

**Remove a universal semantic `Entry` from placement.** The meaningful roots differ:

- Go has package compilation units.
- Rust has a crate root.
- C++ header checking has translation-unit contexts.
- TypeScript has compiler input roots and module configuration.

A snapshot may retain a distinguished node for addressing or presentation. Document that it neither establishes a native root nor limits which published files must be checked. Horos can continue to require its own semantic entry.

### Tests for the division

Before freezing the types, require fixtures for these cases:

| Fixture | What it tests |
|---|---|
| Two Go files in one package | File identity differs from scope identity |
| Two directories with the same Go package name | Package name differs from package identity |
| Output rooted below a Go module subdirectory | Paths have an explicit coordinate system |
| Rust root containing only module declarations | Scaffolding files need no owned type declaration |
| Inline Rust module with an out-of-line child | Module topology determines physical realization |
| One C++ namespace reopened across headers | Scope is independent of file |
| Assignment-list permutations | Enumeration does not change placement meaning |
| One definition plus supporting forward declarations | Exactly-once ownership does not mean exactly one syntax occurrence |

Initially unsupported cases should still be representable and fail with a specific unsupported-case diagnosis. That is a useful test of final structure before implementing coverage.

## 2. Names: U2, with explicit naming intent and exposure

### Collisions belong to the language’s placed scopes

**Choose U2 for top-level collisions, while retaining declaration-local checks during generation.**

Remove the flat-name rule from both the shared builder and the corrected targets. Retain checks whose meaning is wholly inside a declaration: duplicate type parameters, malformed identifiers, and similar local defects.

The existing ban on a type parameter matching *any declaration anywhere in the generated output* crosses that boundary. Removing it requires attention to actual lookup: a parameter can still shadow a declaration used inside its body. Files must either perform an explicitly permitted alpha-renaming or refuse that placement. The binding check establishes which happened.

Also avoid replacing the flat rule with a generic uniqueness map. Languages have different identifier namespaces and declaration-merging rules. Even a native merge that the compiler accepts may violate your requirement to keep two intended declarations distinct.

The TypeScript correction’s retention of flat uniqueness is an implementation fact, not a reason to retain it under AC-09. C++ can temporarily retain a contribution-owned coverage restriction while its correction is incomplete; the shared builder should carry no such policy.

### Add a complete declaration inventory before removing unique-name lookup

**Choose explicit inventory access over syntax walking.**

Every generated placement unit should have an inventory record containing:

- Its output-qualified identity and declaration site.
- Its naming intent.
- Its source provenance, where applicable.
- Its role: requested declaration, dependency, or generated helper.
- Any grouping information needed to select it without interpreting native syntax.

Do not assume one horos declaration produces exactly one native declaration. Specializations and helpers can make that relationship more complicated.

A site remains a good address inside one immutable output. It is not automatically a durable selector across regeneration. Bind a placement to its input capture so that an old site cannot accidentally select a different node in a new output.

This is a prerequisite for U2, not a convenience to add afterwards.

### The hash prefix should become an explicit policy

**Remove the mandatory module-pin prefix from the default naming policy.** Offer stable qualification through an explicitly selected naming policy when a consumer needs it.

The prefix has no established semantic responsibility in the supplied documents. It also has undesirable consequences independently of flat scopes:

- A content change changes the module pin and therefore potentially every prefixed public name.
- A truncated pin is not a uniqueness guarantee.
- Different versions can coexist in different native scopes without different local names.

Ordinary native names are a reasonable default. If two declarations named `Person` are deliberately placed in one conflicting scope, refuse the placement. The consumer can choose separate scopes or an explicit naming policy.

Do not introduce “add a hash only when a collision appears.” That makes adding an unrelated declaration capable of renaming an existing API.

Fix the builder’s comparator at the same time. Equal names need an equality result and a deterministic identity-based tie-breaker. Node IDs assigned by the resulting order cannot themselves supply that tie-breaker.

### Go’s final declaration spelling belongs to Files

**Reject statement 18’s inference that Files can never change a declaration’s spelling.** Generation should choose naming intent; Files should realize it under exposure.

Go exportedness depends on an uppercase initial and declaration context. Preserving every generated uppercase binder while requiring an unexposed declaration to remain unexported is incompatible with assigning module visibility to Files. [The Go Programming Language](https://go.dev/ref/spec?utm_source=chatgpt.com)

I would make the naming contract distinguish:

1. **Declaration identity:** authoritative for references.
2. **Naming intent:** the mapping’s chosen name or naming prescription.
3. **Realized native spelling:** the identifier published by Files.

A Go-owned naming prescription should define the permitted public and internal spellings, including escaping and reserved-word behavior. Files selects the appropriate realization, rewrites linked uses, checks collisions, and reports the realized name.

This does not require names to fluctuate across placements. A fixed output has fixed exposure intent, so an exposed `Person` can remain `Person` in every admitted placement.

It also does not authorize arbitrary collision repair. If realization would produce an invalid identifier, capture a builtin, or conflict with another declaration, Files must use an explicitly authorized naming rule or fail.

### Define exposure’s audience before freezing the output

**Add explicit access intent; keep `Exposed()` as the list of requested roots.**

Your document correctly identifies that “requested root” and “complete visibility policy” have been conflated. Statement 21 is an assumption, and it should not become an accidental universal contract.

For the first profiles, I recommend these meanings:

| Language | Required public access for an exposed declaration |
|---|---|
| TypeScript | Nameable by an importer through its advertised module export |
| Go | Nameable from another package through its advertised package/name |
| Rust | Nameable from another crate through its advertised public path |

For nonexposed declarations, the initial policy should authorize only the language’s local visibility. Additional access, such as crate-internal access, must be recorded as an explicit authorization before placement.

A useful abstract contract is:

> Required access must be established, and realized access must remain within the authorized access.

The audience and access rules belong to the target contribution. Generation records the obligation relative to the eventual owning unit; it need not choose a concrete package or module.

This leaves room for future restricted visibility without reinterpreting today’s private declarations.

## 3. Unplaced syntax: A, generalized to every target

**Choose a declaration-sequence or item-sequence fragment root in every target vocabulary, together with phase-specific admission.**

The root should name a syntax category. It should not claim that all its declarations inhabit one native scope.

You do not need D’s stronger implication that fragments and completed files are different languages or require duplicated grammar implementations. Reuse the declaration grammar under different admissible roots. `generation.Output` and placed `Source` remain distinct phase types; the fragment root preserves that distinction when a reader receives the underlying tree.

The resulting reader contracts should be explicit:

| Reader | Unplaced fragment | Placed file or compilation unit |
|---|---|---|
| Syntax checker | Checks declaration grammar and fragment restrictions | Checks complete native file/unit grammar |
| Syntax formatter and printer | Produces fragment text, without claiming it is compilable source | Produces the complete represented source |
| Generated-Files formatter | Accepts the fragment and its intent | Refuses it as generation input |
| Compiler bridge | Refuses it as a native program | Analyzes it under the declared compilation environment |

Roslyn provides a useful comparison: it has a syntax entry point for parsing a member declaration without requiring a compilation unit. The lesson is that syntax categories can be distinct within one language API. [Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/api/microsoft.codeanalysis.csharp.syntaxfactory.parsememberdeclaration?view=roslyn-dotnet-4.13.0\&utm_source=chatgpt.com)

### Go retains its complete-file guarantee

**Keep the package clause mandatory on Go’s `file` root.** Add a separate fragment root for declarations.

Making the package optional on `file` would weaken every existing reader’s contract. A fragment printer can print declaration text without a package clause; the Files publication path still requires a completed file.

### Rust needs phase information, not an `undecided` visibility attribute

**Use the fragment root to distinguish undecided item visibility from placed privacy.**

In a fragment, absent item visibility means Files has not realized the access intent. In a completed source root, absence has Rust’s ordinary private meaning.

Do not introduce `visibility: undecided` as though it were source syntax. It cannot round-trip through the native language. Continue to express placed private items canonically as the parser does.

**Keep F1 for fields.** Record representation-owned field accessibility in the native model and preserve it through placement. Fragment admission can forbid placement-owned item visibility while permitting representation-owned field visibility.

### Correct the C++ premise

**A bare declaration sequence can be a complete C++ translation unit.** Namespace wrapping is a placement choice, not a grammar requirement. The C++ draft expressly permits a translation unit consisting of a declaration sequence. [eel.is](https://eel.is/c%2B%2Bdraft/basic.link?utm_source=chatgpt.com)

C++ still needs the fragment distinction: an unplaced sequence must not accidentally acquire global-scope meaning merely because its printed text happens to be valid source.

### What the sealed output should guarantee

I would freeze the output around these properties:

- Checked fragment syntax.
- A complete inventory of owned placement units.
- Naming and access intent.
- Intended destinations for every name-sensitive occurrence.
- Snapshot-associated identities and provenance.

It need not contain compiler-established facts about a hypothetical native program. Thus I accept hypothesis 6.9.5 with a qualification: unplaced output still needs strong grammar, identity, and intent-completeness checks.

## 4. Evidence: P2, explicit destinations, and a request-specific evidence policy

### Reconstruction has a legitimate but bounded role

**Reject hypothesis 6.9.4’s distinction between copied and derived values.** That is not the boundary between sound and unsound evidence.

A structural comparison can establish that a placed package clause has the requested spelling. It can likewise establish that an item has a declared `pub` modifier. Whether the value was copied or computed is immaterial to that syntactic proposition.

The important distinctions are:

- Structural checks establish properties of the represented syntax.
- Compiler facts establish resolved destinations and native accessibility.
- Re-running the same layout function establishes consistency with that function; it does not independently establish that the function implements the placement specification.

For example, if `placeGenerated` mistakenly ignores the requested package name, invoking it twice reproduces the mistake. Compare the actual name directly against the admitted consumer input.

Keep reconstruction as a useful integrity check. Supplement it with explicit predicates for preserved declarations, authorized changes, ownership, bindings, and exposure.

### Required facts and comparisons

| Property | Recommendation | Required comparison |
|---|---|---|
| Go package name | **P2** | Consumer’s admitted package name against a compiler package-name fact |
| Go package identity | Compiler ownership | Every relocated declaration’s owner against the package identity derived under the explicit environment mapping |
| Go `string` and `bool` | Recorded external intent | Each resolved use against the intended predeclared entity |
| Rust declared item visibility | Compiler visibility fact | Declared visibility against the permitted realization of access intent |
| Rust public exposure | Audience-specific nameability | Advertised public path is accessible from the promised audience and resolves to the intended declaration |
| Actual generated cross-scope uses | Compiler acceptance and bindings | Each use is accessible and binds to its intended destination |

#### Go package name and identity

**Choose P2 because the compiler already exposes the required observation.** `go/types.Package` provides both `Name()` and `Path()`. Publish the name separately from the path and compare both according to their distinct meanings. [Go Packages](https://pkg.go.dev/go/types?utm_source=chatgpt.com)

There is an important limit to the package-identity claim: `go/types.Config.Check` receives the identifying package path from its caller. It does not discover the filesystem’s module organization. The check establishes ownership consistently within the declared environment; it cannot certify where the consumer later writes the files. [Go Packages](https://pkg.go.dev/go/types?utm_source=chatgpt.com)

Moving a declaration between packages therefore preserves a **correspondence to generated identity**, not one unchanged native nominal identity across placements.

#### Rust visibility

**Publish declared visibility and effective public accessibility as separate facts.**

For declared visibility, use a resolved semantic form such as:

```text
Public
Restricted(scope entity)
```

That is more durable than a `Public` boolean and need not expose rustc’s internal API.

For exposure, check nameability from the promised audience through the advertised public path. A Rust item’s accessibility depends on ancestor modules and reexports, so its own `pub` modifier does not establish the whole property. [The Rust Reference](https://doc.rust-lang.org/reference/visibility-and-privacy.html?utm_source=chatgpt.com)

Also avoid treating rustc’s general “reachable” classification as synonymous with public nameability. Its effective-visibility model distinguishes direct access, access through reexports, and reachability through interfaces; the latter can include types that consumers cannot name. [Rust](https://doc.rust-lang.org/stable/nightly-rustc/rustc_middle/middle/privacy/enum.Level.html?utm_source=chatgpt.com)

For a root-only first slice, these checks are simple. Their meanings should already cover nested modules and public paths.

Require fact coverage for every relevant item. Missing visibility evidence must not silently mean “private.”

#### F1 does not eliminate public-interface checks

**Keep field visibility in the representation, and separately check that the promised representation is usable.**

Consider:

```rust
pub struct Public {
    pub field: Hidden,
}

struct Hidden;
```

Rust has a warn-by-default `private_interfaces` lint for interfaces involving more-private types. Therefore successful compilation alone is insufficient to establish the intended transparent public API. [The rustc book](https://doc.rust-lang.org/rustc/lints/listing/warn-by-default.html?utm_source=chatgpt.com)

For the first transparent-record and alias profile, require an appropriate privacy report or a conservative rejection of relevant private-interface diagnostics. An opaque representation can have different obligations.

Files must never respond by exporting `Hidden` without authorization.

### Record Go predeclared destinations now

**Use explicit Go-issued predeclared symbols.** An API such as `Predeclared("string")` expresses the intended category more clearly than making an empty package field carry two meanings.

An empty package can remain the symbol package’s private encoding. It should not let malformed package-symbol information accidentally become a builtin request.

Compare against the actual predeclared entity, not merely an entity with the name `string`. Both a declaration and a lexical parameter can capture an unqualified identifier.

Change slice 1’s refusal from:

> An output with an external use fails.

to:

> An output with an unsupported external destination fails.

Predeclared types are external to the generated declaration set and require no import.

### Record Rust external destinations now as well

**Move Rust’s external-intent representation into the initial correction.**

These three requirements cannot coexist:

1. Every named use has an intended destination.
2. The output’s recorded intent never changes after slice 1.
3. Rust’s existing external paths acquire intended destinations only in a later slice.

The `ExternalUses` mechanism already provides the relevant shared category. Record Rust-issued symbols for standard and library-owned destinations now. Files may continue to emit absolute standard paths where that is the chosen realization.

Admit only symbols whose destinations the bridge can actually establish under the supplied environment. If the pinned sysroot cannot yet be checked precisely, add that observation or refuse those uses.

Dependency spelling is not enough: rustc permits an external dependency’s locally visible name to differ from the compiled crate name. Bind the symbol to the environment-selected dependency and let Files choose the local spelling. [The rustc book](https://doc.rust-lang.org/rustc/command-line-arguments.html?utm_source=chatgpt.com)

### Constructors and evidence requirements

**Require an immutable language-owned environment at construction, and derive the exact evidence requirement per admitted request.**

Schematically:

```text
NewGenerated(stages, environment)

Format(output, placement)
    admit placement
    construct expected compilation manifest
    derive evidence requirement(environment, manifest)
    construct and analyze placed source
    verify publication against original intent and admitted placement
```

The current constructor-fixed `Evidence` value is too narrow if it cannot incorporate placement-dependent compilation roots and contexts.

I would change the shared contract now to accept a request-policy factory, conceptually:

```text
Prepare(input) -> required evidence + publication guard
```

The factory is supplied by the language contribution. It establishes expected evidence before accepting the provider’s publication.

The language-owned compilation manifest should identify:

| Language | Manifest content |
|---|---|
| Go | Package identities, import paths, and member file IDs |
| Rust | Crate identity, root file, module hierarchy, and physical module realization |
| C++ | Translation-unit contexts and their participating headers and configuration |

The environment identifies compiler/toolchain, language version or edition, target, dependency artifacts, and relevant configuration. The receipt binds those inputs, the compilation manifest, and the exact published snapshot.

The permanent requirement is a rule: **evidence for this exact source under this exact environment and compilation context**. The rule stays fixed even though its values differ between requests.

Represent analysis contexts explicitly. One-context profiles suffice initially; later C++ analysis must not require pretending that one header site has only one compiler occurrence across every compilation.

### Initial admission domains

**Use the Go import environment from slice 1, and introduce the Rust environment from slice 1.**

| | Go slice 1 | Rust slice 1 |
|---|---|---|
| Environment | Explicit module/output-root mapping and compiler configuration | Explicit pinned compiler/sysroot, edition, target, crate identity and kind |
| Placement | One ordinary package, within current file limits | One library crate, root module, one captured file |
| Destinations | Local declarations, parameters, supported predeclared symbols | Local declarations, parameters, and precisely checkable sysroot symbols |
| Visibility | Final spelling and exportedness realize recorded access intent | Declared visibility and public nameability realize recorded access intent |
| Refusals | Multiple packages, package imports, unsupported externals or build modes | Additional modules/files, arbitrary dependency crates, reexports, unsupported visibility forms |
| Failure behavior | No partial `Source` | No partial `Source` |

The Rust environment’s dependency collection can be empty while arbitrary dependencies are unsupported. Its crate and configuration fields can initially accept only the supported values. Later coverage expands accepted values without introducing a new constructor or silently changing defaults.

### Negative tests must survive ordinary compilation

Require tests in which the incorrect output is otherwise valid:

- Change the Go package name.
- Assign the correct Go declaration to the wrong package identity.
- Shadow `string` with a valid local declaration.
- Remove `pub` from an exposed Rust item unused internally.
- Add `pub` to an unexposed Rust item.
- Keep an item public but make its containing module inaccessible.
- Supply evidence for the correct bytes under the wrong compilation root.
- Publish a transparent record with an unusable private field type.

These are tests of the stated obligations. Successful compilation cannot substitute for them.

## 5. Additional issues that can force later interface changes

### A. TypeScript files are not automatically modules

**Establish module status explicitly.**

A TypeScript file with no imports or exports can be treated as a script, putting its declarations in shared global scope. The handbook gives `export {};` as a way to make a file a module without exporting declarations. [Modules](https://www.typescriptlang.org/docs/handbook/2/modules.html?utm_source=chatgpt.com)

The supplied excerpts do not establish whether the provider already forces module treatment, so I would not call this a demonstrated baseline bug. I would require a test now: two files, each containing only an unexposed declaration with the same name.

Require either explicit module syntax or a receipt-bound configuration that establishes the desired interpretation, and check the compiler’s module classification.

### B. Go filenames affect the program an ordinary build sees

**Include build-file selection in the supported environment, or refuse build-sensitive filenames initially.**

Explicitly passing ASTs to `go/types` does not establish that an ordinary Go build selects those files. Names such as `types_linux.go`, `types_amd64.go`, `types_test.go`, and `_types.go` carry selection consequences. [Go Packages](https://pkg.go.dev/cmd/go?utm_source=chatgpt.com)

This differs from your known handwritten-sibling-file limitation: the generated files themselves may be excluded.

The first slice can adopt a deliberately restricted filename policy. Define ordinary packages, command packages, and test packages as distinct supported or refused cases.

### C. C++ namespaces do not provide Rust-style privacy

**Keep exposure policy language-owned and refuse privacy promises the C++ header profile cannot implement.**

A named namespace organizes names; it does not give namespace members class-style private access. Putting a helper in `detail` is not compiler-enforced non-nameability. An unnamed namespace changes identity/linkage behavior and is not a general replacement for private generated types. [eel.is](https://eel.is/c%2B%2Bdraft/namespace.def?utm_source=chatgpt.com)

This is why the universal output should record access requirements and authorizations, rather than equating “not a requested root” with one universal native privacy operation.

Under a strict non-nameability requirement, some header placements must fail. Publicly nameable support declarations can be allowed only when that access was authorized before placement.

### D. C++ evidence needs compilation contexts and occurrence multiplicity

**Distinguish definition placement from compiler occurrences.**

A header participates in translation units through inclusion and preprocessing. The same header can be analyzed in different contexts. The C++ compilation model does not equate source files and translation units. [eel.is](https://eel.is/c%2B%2Bdraft/lex.separate?utm_source=chatgpt.com)

Clang’s compilation database illustrates the appropriate identity separation: it records compilation commands and working directories and permits multiple configurations for the same source file. [Clang](https://clang.llvm.org/docs/JSONCompilationDatabase.html?utm_source=chatgpt.com)

Consequently:

- Exactly-once assignment applies to ownership of a generated definition.
- A declaration may have supporting forward-declaration sites.
- A source occurrence may participate in several analysis contexts.
- Compiler bindings and entities are interpreted within those contexts.

Do not freeze a global bijection from generated site to compiler entity.

### E. Destination identity must survive aliases and lowering

**Specify what counts as the same intended destination.**

A check must not substitute structural type equality for declaration identity. Two aliases or two structurally equal types can remain distinct intended destinations under your accepted rule.

For external symbols, distinguish the selected public export from its underlying canonical type. A semantic provider may need to establish correspondence through aliases or reexports while preserving which public endpoint generation selected.

The same issue affects source provenance: one source declaration may yield several native declarations, and one shared helper may serve several sources. The inventory and relocation contract should accommodate those relationships now.

### F. Failure categories and unsupported fields need stable meanings

**Separate invalid input, unsupported coverage, missing evidence, and exceeded limits.**

AC-05 explicitly includes failure meaning. A future slice should turn an `Unsupported` case into success; it should not require consumers to reinterpret an old `Invalid` error.

Similarly, do not accept an environment or placement field and silently ignore it until a later slice. That changes the meaning of an already accepted input. Unsupported nondefault values must fail immediately.

## 6. The literal freeze requires a broader initial correction

**Normalize all existing generated outputs before the cutoff.** Your documented deferrals otherwise violate your stated test:

| Deferred correction | Later change prohibited by the freeze |
|---|---|
| Rust absolute paths acquire explicit external intent | What existing generated output records about intent |
| C++ generation stops wrapping declarations in a configured namespace | Existing root/containment and scope meaning |
| C++ generation stops emitting includes | Existing file syntax and dependency intent |

This is not a contradiction in the accepted generation/Files separation. It is a contradiction between the proposed sequencing and the stronger freeze criterion.

To retain the literal criterion:

1. Move every existing target to the final fragment contract now.
2. Record every existing name-sensitive destination now, including Rust and C++ external destinations.
3. Remove placement-owned namespace/include syntax from C++ generation now.
4. Keep checked C++ `Source` publication unavailable until compiler-backed evidence exists.
5. Add C++ Files later as a new contribution using the already established generation boundary.

Generating checked fragments does not require implementing C++ semantics. Publishing checked C++ files does.

If legacy C++ generation remains unchanged until its semantic tooling arrives, the guarantee must explicitly exclude that legacy output. That is a coherent narrower commitment, but it is a relaxation of the test in your Question section.

### What I would require before the two slices merge

The first delivered boundary should already contain:

- Language-owned placement types with shared coverage validation.
- A complete, capture-bound declaration inventory.
- Stable naming intent and explicit access intent.
- Distinguishable fragment roots across existing targets.
- Complete named-reference intent.
- Explicit environments and placement-dependent compilation manifests.
- Package-name, ownership, binding, and visibility observations sufficient for admitted cases.
- Refusals that identify unsupported coverage without changing existing input meanings.

That preserves the useful part of AC-10: the first implementation can accept a very small domain, while its identities, responsibilities, and proof obligations already describe the domain you intend to grow.