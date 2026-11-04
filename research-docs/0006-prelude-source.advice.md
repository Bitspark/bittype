# Advice: Authoring the Horos prelude in Horos

See [the research brief and evaluated recommendation](0006-prelude-source.md).

## Recommendation

**Adopt option 1: maintain the complete `horos.kinds` module in Horos source, prepare its canonical Ontos encoding with a Bittype-owned tool, and ship that asset with Horos.** Runtime Horos should decode the asset and admit it through the existing module loader and general kind calculus.

The owner’s concern is justified. The present representation is already an ordinary module, but the maintained definitions still live in Go. Replacing `buildPrelude` with an authoritative `.horos` file addresses that concern directly.

I would keep this change focused:

- Make the named kind definitions source authored.
- Expose the existing syntax-to-canonical conversion through a small Bittype command.
- Preserve Horos’s current selection of the standard formation dependency.
- Retain the exact-prelude source-location accommodations while canonical-only use remains supported.
- Treat broader compiled-module support as a separate capability.

The supplied implementation excerpts already establish that module loading and kind admission are shared with ordinary declarations. The remaining special responsibilities are **authoring, dependency selection, and availability**, so those are the boundaries this change should address. [Source: research brief](0006-prelude-source.md)

## 1. What should be intrinsic?

There are three responsibilities to keep distinct: the language’s primitive vocabulary, the definitions expressed using that vocabulary, and the language configuration that selects particular definitions for formation obligations.

| Responsibility | Recommended home |
|---|---|
| Immutable values, canonical encoding, and structural identity | Ontos |
| Revision-4 declaration forms, sorts, general kind operators, sequence patterns, and guarded fixed points | Horos implementation |
| Primitive descriptor forms, arities, type-child positions, selector interpretation, literal treatment, normalization, and recursion rules | Horos implementation |
| Existing ground well-formedness and literal rules | Horos implementation |
| Which primitive operands require `top(P)` and which require the selected standard data contract | Horos formation policy |
| Definitions of `Type`, `Scalar`, `List`, `Record`, `Data`, and the other standard kind families | Ordinary `horos.kinds.horos` source |
| Supplying the standard module to the resolver | Standard-library composition |
| Readable syntax and its conversion to canonical declarations | Bittype’s Horos frontend |

This leaves a meaningful intrinsic base. Horos source can express `Data` because the implementation already understands node observations, designated type children, Boolean kind expressions, and guarded recursion. Those mechanisms do not need `Data` to explain their own operation.

### Primitive formation may legitimately refer to the standard `Data`

The current `descriptorOperandKind` policy can remain in Go. It specifies which contract a primitive operand must satisfy. The recursive body of that contract should come entirely from the source-authored standard module.

A precise interpretation is:

> The language configuration selects an ordinary module, and primitive formation refers to that module’s `Data` declaration.

The supplied `preludeOwner`, `GroundKind`, and operand-kind resolution already provide that scope. A reference such as `preludeReference("Data")` is therefore reasonable wiring: its binding is determined by the selected module’s identity and declaration scope. It does not require the evaluator to recognize `Data` as a reserved kind.

This gives the intended behavior:

- The standard `Data` receives its meaning from its ordinary kind definition.
- A local kind named `Data` cannot redirect primitive formation.
- Importing another module’s `Data` does not replace the configured formation dependency.
- A differently named or separately defined kind may describe an equivalent set without becoming the declaration selected by the language configuration.
- Existing primitive and ground well-formedness rules remain authoritative.

Likewise, omitted contracts should continue to elaborate to intrinsic `top(P)`. The library definition `Type = top(P)` is useful vocabulary, but omission should not acquire a dependency on an implicit source name.

The substantive authoring change is to remove the complete Go-authored kind bodies. Small Go constructor calls that express `top(P)` or reference the selected library declaration can remain. Removing every occurrence of the string `"Data"` from implementation code would obscure this dependency without improving the separation of responsibilities.

### Why this prelude has no semantic bootstrap cycle

The note gives the decisive simplifying fact: the standard module contains only kinds and has no imports. Its definitions can be admitted using the primitive shape vocabulary and general kind calculus before the resolver selects `GroundKind` and admits type declarations.

That is enough. There is no need for a seed definition of `Data`, privileged admission, or an initial Go-authored copy. The implementation provides the calculus; the module supplies an expression in that calculus. This follows the loading order already described in the note’s resolver excerpts. [Source: loading and formation configuration](0006-prelude-source.md#actual-canonical-loading-and-formation-configuration)

## 2. The smallest useful compilation boundary

I recommend a **generic Bittype-owned command whose contract is one ordinary Horos module in, one canonical Ontos module encoding out**.

The command should perform four operations:

1. Read the supplied source.
2. Invoke the existing Horos parser.
3. Convert the admitted syntax tree to canonical module values using the existing module reader.
4. Encode the selected module with the ordinary Ontos codec.

It should have no knowledge of the name `horos.kinds`.

### A public Go API is optional

The research note identifies the missing capability accurately, but the minimum public boundary can be a command. A command inside the Bittype module can use `internal/horosmodel`; it does not require making that implementation package public.

Internally, the operation currently performed by `reader.module` can be exposed to the command through a small helper. `SourcePins` can use the same helper so that module emission and identity calculation continue to share the existing lowering logic.

If an actual library caller also needs canonical values from an admitted AST, a public Horos-specific operation is justified. Its natural result would associate module syntax nodes with their canonical `core.Value` values. That retains correspondence for callers working with several module sections. It belongs to Bittype’s Horos provider, with no Horos-specific operation added to the language-agnostic kernel.

For this task, however, **the command is sufficient**. I would not make a second public API a prerequisite.

Since the parser supports multiple module sections, the command needs an unambiguous output contract. Requiring one module is enough for the prelude. Supporting an explicit module selection is also reasonable if it fits the existing command conventions. Neither choice requires a new bundle format.

### Keep this operation separate from semantic resolution

The generated result is a **canonical module asset**. It is not a resolved `Program`.

Syntax lowering can encode the prescribed source representation, including default sorts and lexical parameter slots, while leaving ordinary declaration references for the semantic resolver. There is no need to evaluate `Data`, resolve imported declarations, or admit kind definitions merely to emit their representation.

Consequently, the preparation path should not call `resolve.Model` and then extract `Value()` from its result. That route would unnecessarily bring the installed formation dependency into preparation of its replacement.

The public terminology should make this boundary clear. “Encode module” or “lower source” communicates the operation particularly well. “Compile” is also usable if its documentation explicitly describes the output as an unresolved canonical module.

The existing `SourcePins` path is strong evidence that the necessary conversion already exists: it obtains canonical values before semantic resolution and computes identities from them. The change is to make those values usable as output. [Source: readable syntax and dependency-cycle discussion](0006-prelude-source.md#readable-syntax-and-the-real-dependency-cycle)

## 3. A concrete authoring, packaging, and bootstrap arrangement

### Authoritative source and distribution asset

Maintain one complete `horos.kinds.horos` file in the Horos repository. It must contain the entire current standard module, including the families omitted from the illustrative transcription.

The canonical asset is generated from that source and shipped in the Horos distribution. It is not independently edited.

For the Go implementation, place the asset beneath the package that embeds it and within the released Go module. Keeping the authoritative source beside it is a simple way to ensure that both reach Go-module source consumers.

For example, these would be reasonable **proposed** locations:

| Item | Proposed location | Role |
|---|---|---|
| Authoritative module source | `core/go/declarations/standard/horos.kinds.horos` | Maintained definitions |
| Canonical module asset | `core/go/declarations/standard/horos.kinds.ontos` | Generated runtime representation |
| Runtime wrapper | Existing declarations package | Embed, decode, expose, and supply the module |
| Preparation command | Bittype | Parse and encode ordinary Horos source |

The precise directory and asset extension can follow project conventions. The relevant constraint is that Go embedding reads package-relative files and cannot reach outside the package’s Go module. Embedding the generated bytes permits runtime access without locating a file on disk. [Go Packages](https://pkg.go.dev/embed?utm_source=chatgpt.com)

The existing `Prelude()`, `PreludePin()`, and `PreludeImport()` API can remain. Its backing value becomes the decoded asset. There is no architectural benefit in coupling this authoring correction to a public API rename.

### Bootstrap through an already available tool

Let **H0** be the currently available Horos release and **H1** the release that adopts the source-authored prelude.

The coherent sequence is:

1. **Add the preparation command to Bittype while Bittype still builds against H0.**  
   The note says the necessary syntax and canonical encoding are already supported. This step exposes that capability without requiring H1.

2. **Build or publish the Bittype command as an independently usable executable.**  
   It is available before the candidate Horos package needs its new asset.

3. **Author the complete standard module in Horos source.**  
   Include every existing export, sort, selector, sequence pattern, and recursive definition that belongs to the module.

4. **Run the available command to produce H1’s canonical asset.**  
   This path parses and lowers the source. It does not resolve the new module through H0’s prelude.

5. **Replace `buildPrelude` and the Go-authored definitions with embedded-asset decoding.**  
   Helpers needed only to author those definitions can disappear. Ordinary formation-reference construction can remain.

6. **Publish H1 with the source and generated asset already present.**  
   Building or importing Horos does not invoke Bittype.

7. **Update Bittype’s Horos dependency in its normal subsequent release.**

This sequence preserves the runtime dependency direction:

| Context | Dependency |
|---|---|
| Preparation tool build | Bittype → available Horos release → Ontos |
| Preparation input | Horos source supplied to the Bittype command |
| New Horos runtime | Horos → Ontos |
| Later Bittype runtime | Bittype → new Horos release → Ontos |

The preparation executable may link H0’s declarations package and may therefore contain or initialize H0’s prelude. That does not make H0’s prelude part of the new module’s lowering semantics. The required separation is that **the preparation operation never resolves or admits the new source against the installed prelude**.

An independently built executable also avoids an easy build-order trap: a development workspace could redirect the tool’s Horos dependency to the incomplete H1 checkout whose embedded asset has not yet been produced. The preparation recipe should use the already available tool, with no need to build that unfinished target first.

### Ordinary updates

For later changes, edit the maintained source and run a compatible available encoder to regenerate the distribution asset. The generated asset remains part of the released Horos source package, so ordinary downstream builds have no generation step.

If a future change needs new syntax, frontend support can be made available before using that syntax in the prelude. A future change to the declaration language itself may require a coordinated release sequence; the current source-authoring migration does not require solving that broader problem.

## 4. Runtime loading should stay ordinary

Once decoded, the asset should enter the existing module pool as the selected standard dependency. Keep the current ordinary sequence:

1. Load the selected standard module and the root through the shared loader.
2. Follow ordinary pinned imports.
3. Read and bind kind expressions.
4. Admit kind definitions through the general calculus.
5. Resolve the selected ground contract in the standard module’s scope.
6. Admit type declarations using that contract and the existing formation rules.

This retains ordinary declaration identity and resource limits. Distribution as an embedded asset gives the module no special semantic privilege.

It also answers “loaded through the same machinery as any other module” precisely. Horos’s semantic entry point accepts canonical values. The source-authored prelude can enter that same entry point after preparation, just as other readable source is converted to canonical declarations before resolution.

### Preserve one identity model

`PreludePin()` should continue to reflect the identity of the canonical module value. The import-free prelude source does not need to contain its own identity.

An authoring-only migration can preserve the existing identity **if the source lowers to exactly the same canonical value**. Semantic equivalence alone is insufficient: a different canonical expression can describe an equivalent kind while still changing the module’s identity.

Conversely, source formatting or comments do not introduce a new identity when the canonical value remains unchanged.

This is the existing language rule. No separate source identity, asset identity, or sidecar mechanism is needed.

## 5. Source availability and source locations

Source authorship and source availability are separate promises. Maintaining the prelude in `.horos` does not mean every semantic session has parsed that file.

I would specify the minimum policy as follows:

> Horos ships the maintained prelude source alongside its canonical asset. Canonical bytes are sufficient for runtime use. Source-aware callers may explicitly supply the maintained source through Bittype’s ordinary text and Files pipelines.

This gives three distinct relationships:

| Relationship | What it provides |
|---|---|
| Selected formation dependency | The standard declaration used by primitive formation |
| Explicit source-language import | Lexical access to names such as `h.Data` |
| Supplied source tree | Actual syntax sites and text locations |

None implies the other two. Supplying prelude source does not create an implicit lexical import. An explicit import does not guarantee that source text is available. Formation can use the standard module without either a user import or a source tree.

### Keep the exact-prelude accommodations now

The three existing accommodations remain appropriate for this change:

- The selected prelude may be absent from source correspondence.
- Its semantic entities may lack original syntax sites.
- Files composition may omit a physical file relation when no prelude file was supplied.

These are narrow accommodations for an injected canonical dependency. Moving authoring to Horos source does not eliminate the condition that requires them.

Keep them tied to the **exact selected module identity**. A module-name test, a declaration-name test, or a structural resemblance test would weaken the boundary. Another module named `horos.kinds` cannot supply source sites for the selected module merely because its name matches.

When the maintained source is supplied, the existing ordinary path should provide locations. The note already says this works, so that path does not need replacement. [Source: source locations and Files composition](0006-prelude-source.md#source-locations-and-files-composition)

### Make the source-distribution promise explicit

Shipping source in the repository and Go-module distribution supports callers that supply that file. It does not promise that an installed executable exposes the source text.

If standard-library navigation must work inside a standalone Bittype integration without a caller-supplied file, a small additional convenience is reasonable: embed the **same maintained `.horos` file** and expose its text through an accessor.

That accessor introduces no Bittype dependency into Horos. Explicit host composition can pass the returned text to the normal parser and Files pipeline. Locations then refer to actual supplied text, with an appropriate source label.

I would add this only if that experience is required now. The source-authoring correction itself needs no new text-access API.

### Diagnostics without source

When source is absent, diagnostics can identify the standard module and declaration, and attach a failed consumer obligation to the consumer’s own syntax site. They should not invent prelude line and column locations.

Projection from a canonical module can still produce a useful generated declaration view. Its nodes have locations in that generated view, if any; they do not recover the maintained source’s original comments, formatting, or offsets.

For this small module, supplying the real source through the existing pipeline is sufficient whenever original locations are wanted. A new source-map distribution format would add unnecessary scope.

### When general compiled-module support becomes justified

A broader facility becomes useful when Bittype needs to accept arbitrary canonical-only dependencies.

At that point, the input contract should distinguish:

- Source inputs, for which syntax correspondence remains required.
- Canonical-only inputs, whose semantic entities may legitimately lack original syntax sites.

Semantic dependency relationships would remain available for both. Physical file relationships would exist where actual files exist. Supplied source could attach through the existing module-identity correspondence.

That general facility would subsume the current prelude exceptions. The important step is explicit support for canonical-only inputs; globally relaxing missing-site errors would lose the existing guarantee for source-derived entities.

This belongs in the relevant provider and Files composition. The shared Bittype kernel should acquire no knowledge of `Data` or the standard Horos module.

## 6. Qualifications and overlooked consequences

### The selected `Data` is part of the release’s formation semantics

Moving `Data` into a library source file does not make changes to its body inconsequential library edits.

Because primitive formation selects that definition, changing it can affect admission even for a program that never explicitly imports `horos.kinds`, within the bounds of the intrinsic formation and ground well-formedness rules.

That is an existing consequence of the design. Document the selected standard module as part of the Horos implementation’s formation behavior.

A related distinction matters when versions differ:

| Reference | Binding |
|---|---|
| A consumer’s explicit `h.Data` | The declaration named by that consumer’s pinned import |
| A primitive operand’s standard data contract | The declaration selected by the resolver’s formation configuration |

Supplying an older or newer explicitly imported standard module does not itself change the second binding. This deserves a clear compatibility statement, without introducing a new profile or registry.

### The illustrative source is not yet the complete module

The transcription omits families referenced by `Data` and other current exports. The migration must move the whole module.

In particular, preserve the existing observation semantics: unions retain their union shape, literal payloads remain opaque, and recursion follows the existing normalized graph behavior. The source-authoring change should not become an occasion to alter matching rules or simplify the definition into a different canonical structure.

### “Ordinary prelude” does not require identical availability policy

The desired result can accurately claim:

- Ordinary source definitions.
- Ordinary canonical representation.
- Ordinary declaration binding.
- Ordinary kind admission and evaluation.

The language still selects and supplies one standard formation dependency. That is coherent.

The Rust comparison is useful at this level: its documentation separately describes a library prelude and the policy that imports it by default. Horos should retain its own explicit lexical-import behavior rather than adopt Rust’s availability rules. [Rust](https://doc.rust-lang.org/core/prelude/index.html?utm_source=chatgpt.com)

For clarity, documentation could describe `horos.kinds` as the **standard kinds module**, and separately describe its role as the selected formation dependency. The existing `Prelude` API names can remain.

### Canonical distribution does not supply missing language implementation

The asset is portable as an Ontos value. An implementation still needs revision-4 declarations and the general kind calculus to interpret it.

Given the note’s description of the other implementations, the immediate benefit is a shared source and distribution representation. It does not by itself extend those implementations beyond ground conformance.

## 7. Scope I would approve

The coherent change consists of one complete Horos source module, one generated canonical asset, one Bittype-owned source-to-canonical command, and a Horos wrapper that decodes and supplies the asset through the existing resolver.

I would retain the current formation selection and exact-prelude source accommodations. I would defer moving the frontend below Horos, requiring every application to assemble the standard dependency manually, configurable prelude replacement, and general compiled-module/source-map infrastructure.

**This resolves the owner’s authoring concern completely:** the standard kind definitions are maintained in Horos, the existing frontend produces their ordinary canonical representation, and the existing semantic machinery determines their meaning. Go retains the primitive language implementation and the small amount of composition needed to make its selected standard dependency available.
