# Research: Truthful capability reports from owned definitions and tests

**ID:** 0012
**Date:** 6 November 2026
**Status:** delivering
**Graded:** https://github.com/Bitspark/bittype/pull/126#issuecomment-6689094045
**Run-ID:** run_43b89675-3878-47ea-93fd-2452942f7f4d
**Reviewed:** https://github.com/Bitspark/bittype/pull/126#issuecomment-6688471895
**Owner:** codex-bittype-rust-syntax
**Issue:** https://github.com/Bitspark/bittype/pull/126
**Decision recipient:** research owner, coordinating with the Bittype integration session
**Code baseline:** Bittype 02489506c817a93e71878c55a254da5587791fd4

## Question

How should Bittype implement detailed capability reports behind its automatically
generated provider/transformer matrices, extracting bounded factual inventories
from owned Go definitions and associating explicit test cases/results, without
inventing completeness, semantic authority or architectural conformance?

The owner approved the documented proposal and now requests expert consultation
followed by full implementation. We can refine the proposal before implementation.
This is documentation tooling and qualification integration, not new language
semantics, native generation behavior or runtime bindings.
Delivery means the complete reporting/extraction/result-association mechanism,
scope links for existing providers, inventories for the existing native grammar
format, and useful associations to current tests. It does not mean inventing an
exhaustive per-feature test catalogue or closing every native-language gap.
Unassociated areas must remain visibly unassessed rather than a TypeScript-only
demonstration being labelled a complete rollout.

## Context

### System overview and vocabulary

Bittype is a public Go module providing typed syntax infrastructure for Horos,
TypeScript, Go, Rust and C++. Horos is the canonical type/declaration language;
its semantics belong to the independent Horos library. An AST (abstract syntax
tree) stores types/declarations; a CST (concrete syntax tree) adds tokens and
trivia (whitespace/comments).
Native tools currently cover bounded types-only subsets, not full compilers.

The generic engine registers separately supplied operations. A language marker
is a distinct Go type implementing `language.Provider` (`ID`, `Name`). There are
role interfaces for parsing, syntax checking, formatting, printing, reference
resolution and semantic analysis. An ASTTransformer has explicit source/target
language identities and transforms admitted source ASTs using explicitly supplied
semantics and mappings. Configuration/profile identity is separate from language.
Here admission means carrier structure and applicable syntax checks, not semantic
validity; semantic evidence is supplied separately. A profile names a documented
behavioral boundary, while configuration supplies actual collaborators, compilers
and environment choices. Equal profile labels do not establish equal execution context.

Syntax operations act on AST/CST/text. Files operations act on an explicitly
captured source collection and semantic results. Those labels are dimensions,
not inferred directory categories. A provider implements one capability; an
application composes them. A populated matrix cell means an implementation exists
under discovered construction requirements, not complete language support.

For example, TypeScript syntax now admits conditional and mapped types, including
`infer` and key remapping, but its bounded semantic exporter does not cover them.
The generic CST printer can print an admitted CST independently of native
semantics. Four direct Horos-to-native transformer edges exist, but their selected
representation mappings do not support every Horos expression or establish that
native callers obey every source constraint.

The owner wants the matrix to stay small, with each implementation linking to a
report of scope, extracted forms, associated cases and available executed results.
There is no agreed denominator for "percent complete", and adding a finite test
case is not proof that all combinations of a feature work.

### Architectural constraints and reusable prior advice

Shared mechanisms own identity, association, lifecycle, resource/failure guarantees
and composition. Domain meanings enter through independently supplied owners.
Adding an extraction format within the promised boundary changes its extractor
and consumer composition, not a central format/language catalogue or unrelated
extractors. A plugin label or empty interface accepting arbitrary payloads does
not establish that boundary. Passing examples do not establish architecture.

Prior Bittype consultations 0007 (type-owned mappings), 0010 (language-owned
judgments) and 0011 (scoped association products) support those distinctions.
They do not answer this documentation-extraction design. We checked all eleven
existing collection entries, current research files/branches and live runs;
there is no existing consultation on capability reports. We are not importing
generation/proof protocols into this tool just because association matters here.

The owner has now assigned stable architecture references: AC-01 architecture is
correctness; AC-02 semantic ownership; AC-03 additive extension; AC-04 the same
public mechanism for independent contributions; AC-05 preserved composition
guarantees; AC-06 limited coverage cannot redefine architecture; AC-07 evidence
supports only established claims; AC-08 general contract changes need a general
justification. Identified violations must enter the maintained ACX-numbered
violation/correction corpus with the actual dependency/ownership failure and its
repair state. Those architectural cases are human explanations, distinct from
the `bittype:case` test associations proposed below; neither supplies a score.
This capability-report work has not yet identified a new architectural violation.

Scope constraints:

- Existing matrices remain source-derived; no second manual capability catalogue.
- Three proposed comment directives attach definitions, scope documents and tests.
  They do not change runtime methods, AST shapes or language semantics.
- Keep authored scope/exclusions separate from extracted facts and observations.
- No completeness or architecture score, keyword-based compliance check, backup,
  extra provenance ledger or administrative validation workflow. Existing test
  execution context may be used to avoid crediting stale/inapplicable observations.
- Go source loaded by this tool is inspected, not executed. Configured extractor
  Go code is trusted tooling, not magically sandboxed by a read-only interface.
- Portable builds/tests/generation run on GitHub/shared workers, not the owner's
  workstation. Editing, Git and consultation coordination are local.
- Preserve current matrix discovery, explicit provider identity, generic dependency
  availability, direct transformer directions and authored documentation.

Ownership, truthful association and bounded failure behavior are fixed. Directive
names/syntax, extractor API shape, package placement and static-versus-executed
publication are proposed mechanisms the expert may replace or simplify.

### Relevant code: source loading and typed discovery

The current command is `go run ./tools/provider-matrix/go`. Its path is
`run -> load -> discover -> render -> replaceRegion -> write`.

These are the actual small contracts that discovery matches, in separate packages:

```go
// contracts/language/go/language.go:11,22-25, package language
type ID string
type Provider interface {
    ID() ID
    Name() string
}
```

```go
// contracts/check/go/provider.go:12-15, package check
type Provider[L any] interface {
    language.Provider
    Check(ast.View[L]) error
}
```

```go
// contracts/parse/go/provider.go:9-12, package parse
type Provider[I, O any] interface {
    language.Provider
    Parse(context.Context, I) (O, error)
}
```

Here `ast.View[L]` is Bittype's borrowed AST view. In the source-discovery code
below `ast` instead means the standard library's `go/ast` package.

`tools/provider-matrix/go/load.go:23-49,59-64` (relevant definitions and call):

```go
type packageInfo struct {
    ImportPath, Dir, Export string
    GoFiles, CgoFiles []string
    Module *struct { Path, GoVersion string; Main bool }
    Error *struct{ Err string }
}
type sourcePackage struct {
    types *types.Package
    info *types.Info
    files []*ast.File
}
type declaration struct { pkg *sourcePackage; node *ast.FuncDecl }
type program struct {
    root, module string
    fset *token.FileSet
    packages map[string]*sourcePackage
    functions map[*types.Func]declaration
    dimensions map[*types.TypeName]dimensionMetadata
}

func load(root string) (*program, error) {
    command := exec.Command("go", "list", "-deps", "-export", "-json", "./...")
    command.Dir = root
    var stderr bytes.Buffer
    command.Stderr = &stderr
    data, err := command.Output()
    // Decode metadata; parse main-module GoFiles with parser.ParseComments;
    // types.Config.Check each package; use export data for external dependencies.
}
```

The actual function handles errors and package cycles. It skips tools packages
as discovery candidates; production filtering also excludes contracts, trees,
engine, internal and integration infrastructure. `go list -export` may compile
dependency export data, so generation is a remote workload. No provider constructor
or initializer is executed. Currently there is no context deadline/output byte
limit on this command and no overall source-file/typechecking budget. Bounding only
new extractor recursion would not make the entire tool bounded.

Tests are not loaded: `TestGoFiles`/`XTestGoFiles` are absent. A new test-source path
must not let synthetic test packages become provider implementations or replace
production Go symbol identities. External dependency source/docs also need explicit
selection/loading; an export-only package cannot supply source annotations.

Discovery finds concrete exported markers under `syntax/` implementing the
language interface; ID/Name must return constant strings. It discovers provider
interfaces under contracts and matches actual Go method-set signatures/type
arguments. Current generic providers have one language parameter; public
constructor signatures supply alternative sets of required provider collaborators.
A least fixed point admits only available sets: a dependency cycle cannot create
availability. Typed transformer factories use typed call traversal through whole
factory bodies with concrete generic substitution, retaining only direct endpoints.
Unlike provider delegation below, this does not prove runtime reachability or
returned-value flow; it must not be reported as such a proof.

`tools/provider-matrix/go/discover.go:17-36,236-239,365`:

```go
type location struct { label string; pos token.Pos }
type row struct { level, contract string; cells map[string][]location }
type matrix struct {
    languages []language
    rows map[string]*row
    edges map[string]map[string][]location
}
type contribution struct {
    language language
    contract *types.Named
    implementation *types.Func
    alternatives [][]*types.Named
}
func contractName(contract *types.Named, lang language) string {
    name := types.TypeString(contract, func(p *types.Package) string { return p.Name() })
    marker := types.TypeString(lang.marker, func(p *types.Package) string { return p.Name() })
    return strings.ReplaceAll(name, marker, "L")
}
// Once a contribution is admitted:
m.rows[key].cells[c.language.id] = append(m.rows[key].cells[c.language.id],
    location{c.implementation.Name(), c.implementation.Pos()})
```

`language` stores the discovered ID/name and `*types.Named` marker. The display
contract string uses package names rather than full import paths. Detailed reports
must retain real contract, marker, receiver and operation identities instead of
reconstructing association from a display row, name or source position.

`load.go:229-258` follows only transparent single-return delegation:

```go
func (p *program) entry(fn *types.Func) *types.Func {
    seen := map[*types.Func]bool{}
    for fn != nil {
        fn = fn.Origin()
        if seen[fn] { break }
        seen[fn] = true
        decl, ok := p.functions[fn]
        if !ok || decl.node.Body == nil || len(decl.node.Body.List) != 1 { break }
        ret, ok := decl.node.Body.List[0].(*ast.ReturnStmt)
        if !ok || len(ret.Results) != 1 { break }
        call, ok := ret.Results[0].(*ast.CallExpr)
        if !ok { break }
        next := p.called(decl.pkg, call.Fun)
        if next == nil || next.Pkg() == nil || !p.implementation(next.Pkg().Path()) { break }
        if _, ok := p.functions[next.Origin()]; !ok { break }
        fn = next
    }
    return fn
}
```

`called` resolves identifiers/selectors/type-indexed calls through `types.Info`.
Only the final entry survives today. A scope comment on a public wrapper and a
case targeting its delegate need the retained conservative chain, not arbitrary
call-graph reachability. Dynamic dispatch does not establish delegation.

### Relevant code: structured grammar and opaque procedural meaning

`internal/grammar/rules.go:15-44,66-71`:

```go
type Attribute struct { Required bool; Accept func(string) bool }
type Role struct { Name string; Min, Max int; Kinds []string }
type Rule struct { Attributes map[string]Attribute; Roles []Role; Reference bool }
type Rules map[string]Rule
func Attr(required bool, accept func(string) bool) Attribute { return Attribute{required, accept} }
func One(name string, kinds ...string) Role { return Role{name, 1, 1, kinds} }
func Opt(name string, kinds ...string) Role { return Role{name, 0, 1, kinds} }
func Many(name string, min int, kinds ...string) Role { return Role{name, min, -1, kinds} }
func Set(values ...string) func(string) bool {
    return func(s string) bool {
        for _, v := range values { if s == v { return true } }
        return false
    }
}
func Attributes(pairs ...any) map[string]Attribute {
    m := map[string]Attribute{}
    for i := 0; i < len(pairs); i += 2 { m[pairs[i].(string)] = pairs[i+1].(Attribute) }
    return m
}
```

`check/syntax/typescript/go/check.go:14-39` defines a `types` string slice of
allowed type-expression node kinds and the actual grammar map. Selected entries:
`g` means Bittype's `internal/grammar`, `s` in checker signatures means
`syntax/typescript/go`, and `ast` in those signatures means `trees/ast/go`.

```go
var rules = g.Rules{
    "conditional-type": {Roles: []g.Role{
        g.One("check", types...), g.One("constraint", types...),
        g.One("true", types...), g.One("false", types...),
    }},
    "infer-type": {Roles: []g.Role{g.One("type-parameter", "type-parameter")}},
    "mapped-type": {
        Attributes: g.Attributes(
            "readonly-modifier", g.Attr(false, g.Set("add", "remove")),
            "optional-modifier", g.Attr(false, g.Set("add", "remove"))),
        Roles: []g.Role{g.One("type-parameter", "type-parameter"),
            g.Opt("name-type", types...), g.One("type", types...)},
    },
    // Other real entries include identifiers checked by custom functions.
}
```

`Check` first calls `g.Check(tree, "source-file", rules)` (line 52), then
procedurally enforces infer placement and disallows mapped/inferred parameter
defaults (lines 59-85). Extracting the map can describe structural child roles,
attribute presence and finite `Set` values, but cannot claim those procedural
restrictions or that mutable variables retain initial values at execution time.
Custom predicates need opaque source links, not a heuristic interpretation.

The grammar-format extractor can recognize exact owning helper/type identities
and bounded literals/constant/reference expressions. It must not execute helper
functions, dispatch by language name or treat any map named `rules` as this format.

### Existing annotations and proposed format

Existing `bittype:dimension Syntax|Files` and `bittype:shared` directives on type
declarations classify operation boundary carriers. Dimension labels are open,
source-owned identifiers. `dimensions.go:20-56` reads raw comment tokens, because
`ast.CommentGroup.Text()` can remove directive-like comments. It accepts only its
known directives on type declarations and errors on unknown ones. New directives
must not weaken those rules or accidentally interpret code fences/trailing comments.

The proposal adds these declaration doc-comment lines:

```text
// bittype:grammar <symbol-ref> [profile=<value>]
// bittype:scope <document-ref> [profile=<value>]
// bittype:case target=<symbol-ref> feature=<key> case=<key> expects=<outcome> [profile=<value>] [language=<symbol-ref>]
```

Example (hypothetical comments, existing operations):

```go
// bittype:grammar rules
// bittype:scope docs/syntax.md#typescript
func Check(view ast.View[s.Language]) error { /* existing body */ }

// bittype:case target=parse.Parse feature=mapped-type case=key-remapping expects=accept
func TestMappedKeyRemapping(t *testing.T) { /* assert the relevant structure */ }
```

There is one directive per standalone `//` doc-comment line; unquoted tokens or
Go double-quoted string values. Unknown/repeated named keys, unresolved symbols
and invalid attachment sites diagnose. `grammar` attaches to an operation/delegate
and resolves a definition symbol; multiple roots are explicit and conflicts are
not silently overwritten. `scope` links to an existing module-root-relative
Markdown section, preserving the owning prose as authority. No absolute paths,
path escapes or ambient source fetches. Existing import aliases/local symbols or
`import/path#Receiver.Method` / `import/path#Function` references resolve through
actual Go type information. Qualified references avoid imports used only in comments.

`case` attaches to a named Go test function. Feature/case keys are nonempty UTF-8
values scoped to owner; grammar keys preserve existing node-kind strings. Named
cases may also concern semantic or transform operations without AST-node keys.
`expects=accept|reject` labels what the owning test asserts, not its execution
state. A semantic invalidity verdict is distinct from an operation error.
Multiple explicit target associations are allowed if the test exercises them;
the host cannot prove assertion adequacy. Precise dynamic subtest declarations
are deferred; function-level families do not become individual scenario counts.

An absent profile means unspecified, never all profiles. A present profile is the
owner's exact existing profile label, not a new registration. The annotation does
not configure the test. Shared generic targets need an explicit language marker
if typed association remains ambiguous; a transformer factory supplies its actual
direct source/target endpoints. No route-wide or all-profile credit is inferred.

### Proposed extractor boundary and unresolved design choices

The proposal puts a tooling-only contract under
`tools/provider-matrix/capability/go`, the initial implementation under
`tools/provider-matrix/extract/grammar/go`, and selection in the existing command.
No runtime package imports it. An independent owner can provide a different
extractor, selected by exact definition/type identity, through the same contract.

A request should carry the selected operation/definition/profile plus a bounded
source/type view. A result should contain normalized structural observations:
feature key, child role/cardinality/allowed-kind constraints, attribute presence
and finite values, reference allowance and opaque custom-check links. Association
must be stamped/validated by the host; the extractor cannot create providers,
claim tests passed or attach records to another request. Unknown extraction
support is explicit and selection conflicts do not fall back in registration order.

We need advice on the precise boundary. Exposing raw Go AST/type maps is mutable,
even if an interface says read-only. Returning `any`/arbitrary JSON or giving
extractors rendering callbacks would weaken association and allow unstructured
claims. A narrow read API can instead supply normalized expression/symbol records,
but should not become a second universal Go compiler. Trusted in-process Go code
cannot be fully sandboxed or forcibly cancelled by interface convention alone.
Resource/trust guarantees must say what is actually bounded.

An independent fixture using Go struct-field declarations rather than the current
helper/map encoding should exercise the same extension contract while orchestration
stays fixed. Another language using the same grammar table is not sufficient
evidence of a different definition-format contribution.

### Existing test execution and publication

`.github/workflows/go.yml:116,127,130,156,161,166,213` executes these relevant
commands in separate remote steps/jobs (shown verbatim apart from grouping):

```sh
go test -count=1 -v ./...
go test -race -count=1 -run '^(TestEngineConcurrent|TestGoSemanticConcurrent|TestTypeScriptSemanticConcurrent|TestRustSemanticConcurrent)' ./integration/go
go test -race -count=1 -run '^TestJointConcurrentChecks$' ./verification/integration/go
go test -mod=readonly -count=1 -v ./...  # in the independently installed consumer
go test -run='^$' -fuzz=FuzzNativeText -fuzztime=20s ./integration/go
go test -run='^$' -fuzz=FuzzSourceRoundTrip -fuzztime=20s ./parse/syntax/horos/go
go test -count=1 -v -run '^(TestNativeCompilerFixtures|TestGenerationSPIHostedGoCompiler)$' ./integration/go
```

The contracts job checks Go 1.27.3 and prepares TypeScript 5.9.3 plus Rust 1.90.0
semantic helpers. The native job supplies compiler fixture environments separately.
Native compiler fixtures skip without that environment in the contracts job.
Both jobs check `HEAD == EXPECTED_REVISION`; a green package/job alone does not
establish that a specific case ran. Current commands emit text, not test JSON.
Go test events identify package/test/outcome, not the provider profile actually
selected inside a test. CI can supply its toolchain/environment context, but the
profile annotation currently remains an authored assertion about the test.
Matching a test event does not upgrade that assertion into an observed profile fact.

Go's `test2json` contract distinguishes test events (Package, Test) from package
events (no Test); parallel events can interleave. Cached results/subtest filtering
and parent aggregate outcomes must not manufacture individual scenario credit.
Useful primary references: https://pkg.go.dev/cmd/test2json and
https://pkg.go.dev/go/ast#CommentGroup.Text.

The current matrix update workflow is independent of tests. It runs on main pushes,
merged-PR events (including skip-ci merges) and manual dispatch, serialized into
one updater. `.github/workflows/provider-matrices.yml:44-53,63-92`:

```yaml
- uses: actions/checkout@11d5960a326750d5838078e36cf38b85af677262
  with:
    ref: main
    persist-credentials: false
- name: Generate tables from implemented providers
  run: go run ./tools/provider-matrix/go
```

The publication step then uses a separate credentialed environment:

```sh
if git diff --quiet -- docs/providers.md; then exit 0; fi
branch=automation/provider-matrices
base=$(git rev-parse HEAD)
git switch -c "$branch"
git add -- docs/providers.md
git -c user.name='github-actions[bot]' -c user.email='41898282+github-actions[bot]@users.noreply.github.com' commit -m 'docs: regenerate provider matrices [skip ci]'
# Enable ephemeral gh credential helper only now; fetch main.
# If origin/main advanced, dispatch a fresh updater and stop.
# Force-with-lease the exclusively owned automation branch, create/update
# its docs-only PR and squash merge that exact head; delete that branch.
```

Thus a source S's test result is not automatically applicable to the updater's
latest main, and committing a report creates another revision D immediately.
The draft proposal permits an optional explicit results input but does not yet
settle how to present observations without stale 'current' pass labels.
Concretely, the contracts job can pass after its native case skips, the native
job may still be running, and the independent updater can already publish another
documentation commit. None establishes that native case execution has completed.

One considered implementation is to commit only deterministic inventories,
scope links and named case associations, while generating reports with observed
results as artifacts/job summaries of the existing test jobs for their exact
source/configuration. The static detail page links to the ordinary CI route and
does not chase mutable 'latest success'. This avoids a new ledger, implicit remote
fetching and re-running tests merely because generated documentation changed.
It still fully implements explicit result association for consumers providing
matching events/context. We want advice on whether this split is the right model.

### Detailed reports and failure behavior

The proposed generated detail file is `docs/provider-capabilities.md`. Each
implemented provider/factory cell retains its source link and adds a details link.
Unavailable cells stay `/`. Multiple implementations remain separately identified.
Report anchors must derive from canonical operation/contract/marker/profile
identity, not line numbers, case counts, hashes or human-readable names alone.
Both files must be rendered/validated before mutation/publication. Two OS file
writes are not a filesystem transaction; the existing Git commit is the external
publication boundary. Errors must prevent publishing half-updated links.

Each provider section includes identity, authored scope links, extracted definitions
with explicit extraction limits, named case associations with expected outcome,
and available execution results. A report may have no grammar extractor (e.g.
printer or semantic provider) yet valid scope/case associations. Transformer
reports additionally name exact endpoints and selected mapping/profile context;
target AST vocabulary does not prove source representability or constraint
preservation. No decoder silently weakens an unsupported/missing interpretation.

We must distinguish at least: declared structure; opaque procedural check;
missing extractor; explicitly partial extraction; defined case without execution;
pass; fail; skip; mismatch in supplied evidence; known authored exclusion; and
architectural violations documented independently. Making those all one scalar
'support' status would be misleading.

### Options considered and rejected

1. Manual feature/progress catalogue: rejected because it duplicates implementation
   availability and drifts; intended scope remains authored where already owned.
2. AST-kind/switch/test-count heuristics as support: rejected. A parser can lack a
   declared node; a formatter may print a form not covered semantically; tests
   can miss branches and combinations.
3. Executing definitions/predicates to discover support: rejected. Source inspection
   must not execute inspected packages or infer semantic claims from runtime shapes.
4. One central grammar decoder teaching every language's special cases: rejected.
   The grammar-format owner may interpret its schema; host orchestration may not
   learn each independent format or library meaning.
5. Raw unrestricted callbacks/payloads with a plugin label: rejected; missing
   identity, scope, resource and result-shape guarantees make composition dishonest.
6. Green CI as per-case result or persisted 'latest pass' after source advances:
   rejected. Explicit named events and applicable context are required, with no
   claim of whole-feature or architectural conformance.
7. New reporting service, ledger, approvals or artificial checksum workflows:
   rejected. Extend existing source tooling and test/artifact publication only.

## Questions for the expert

1. What precise vocabulary and report model best separates definition inventories,
   author assertions, test associations and observations while staying useful for
   finding unassessed areas within an explicitly declared scope? How should partial extraction and custom procedural
   checks be expressed without implying semantic completeness?
2. What minimal extractor request/result contract preserves independent format
   ownership, exact association and bounded interpretation? How should it expose
   Go source/type information without mutable host state, an arbitrary payload
   escape hatch or an unnecessary second compiler? Describe a materially different
   contribution that would test this boundary.
3. How should annotation target/profile/language identities be resolved for generic
   providers, wrappers and direct transformer factories, especially test-source
   loading and multiple implementations sharing one displayed cell? What pitfalls
   or simplifications would you recommend in the proposed grammar/scope/case syntax?
4. How should named Go test observations and source/configuration context be
   associated and published using these existing workflows? Assess the static
   committed report plus per-run result artifact approach, including skip/cache,
   subtests, separate jobs, failed tests and automatic docs-only commits.
   What minimum execution context can be established reliably, and which
   profile/configuration claims must remain authored rather than observed facts?
5. What have we missed? Recommend a conceptually sound, appropriately bounded
   implementation sequence and meaningful counterexamples/tests, including any
   changes you think the documented proposal itself needs before implementation.

## Applied

Evaluation is against Bittype `902121e4581a10914923a359c329fb4af1176025`,
which merges the owner's architecture rules into the submitted code baseline.
Attempt 1 submitted 29,532 bytes, SHA-256
`96ebf4613f685c0975d49def09f9d45dc2d1f51d5d6bf3e82b0c50aae954494b`,
to run `run_43b89675-3878-47ea-93fd-2452942f7f4d`. The local frozen input
still matches. Nightfall completed and verified the local advice against its
stored result, SHA-256
`80556cec8f7ef1351df5b2d3583b0ff7c3c2e93687a4742fdabbd955d07f11ec`.
The living document differed from the submitted input only in lifecycle metadata.
The raw advice's sandbox link is not needed: its substantive answer is present.

The research owner evaluated the complete answer, with independent source-bound
reviews of discovery/source access, grammar interpretation, and CI observations.
The advice confirms the existing direction; these refinements are part of the
authorized implementation, not a new runtime or semantic-model task.

| # | Recommendation | Verdict | Action | Owner | Link | Validation |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | Keep identity, authored scope, structural facts, cases, execution and applicability separate | HOLDS | Use separate records and report sections; no percentage, support score or architecture verdict | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | Report tests distinguish absent associations, unknown extraction and absent observations |
| 2 | Definition-centered extractor; host retains provider/profile/test authority | HOLDS | Public tooling contract takes only definition/source/limits; validate and associate output in the host | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | Invalid selectors/results and foreign spans fail; no destination IDs in extractor output |
| 3 | Field-level knowledge and namespaced features | HOLDS | Preserve known/unknown/unconstrained/opaque fields, root/format/key identity and unresolved entries | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | Custom predicates, undecodable kinds/map keys and duplicate-root cases |
| 4 | Immutable source strings and narrow existing Go fact queries | HOLDS | Extractors parse private copies; source views return copied records, physical spans and canonical symbols; no mutable compiler objects | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | Span/selection/alias/constant tests and independent extractor fixture |
| 5 | Exact owning helper identities, static interpretation only | HOLDS | Grammar extractor recognizes its owned types/helpers; new defined lookalikes do not match; unsupported expressions stay explicit | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | Unrelated Rules/One, aliases, opaque functions, cycles and partial roles |
| 6 | Public composable host plus materially different extractor | HOLDS | Move the host out of package main; command composes grammar extractor; independent struct-format fixture uses the same public runner | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | External-package fixture supplies struct-field schema without host/renderer changes |
| 7 | End-to-end resource guarantees, not extractor-only budgets | NEEDS ADAPTATION | Bound loading/output/source/records/queries/traversal and use cancellation plus the existing remote job deadline. Trusted in-process code is not a hard CPU/memory sandbox | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | Oversized input/output and traversal failures; document non-preemptible Go checker/extractor limit |
| 8 | Retain canonical contribution identities and static wrapper relationships | HOLDS | Keep contract arguments, operation/receiver and marker/endpoints before rendering; wrapper navigation never transfers observations | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | Existing availability/factory tests plus rewriting wrapper, generic ambiguity and stable-anchor tests |
| 9 | Isolated test-declaration index and exact event-name join | HOLDS | Parse active TestGoFiles/XTestGoFiles without interpreting bodies or admitting test-only implementations; duplicate runner names remain ambiguous | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | External test aliases, invalid signatures, test-only providers and duplicate TestFoo |
| 10 | Strict raw-comment directives and owner-defined expectations | NEEDS ADAPTATION | Keep existing dimension attachment rules; add standalone grammar/scope/case parser, nonempty authored expects labels and safe module-root Markdown fragments | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | Unknown/repeated keys, quoted labels, fences, wrong attachments, symlink/path escapes and fragment tests |
| 11 | Stable injective anchors and explicit inspection build context | HOLDS | Encode canonical tuples without hashes; omit timestamps/HEAD from committed reports; distinguish absent profile from any string | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | Reordered discovery and documentation-only revision produce stable output; profile collision tests |
| 12 | Observe executions first; preserve Go JSON protocol and failures | HOLDS | Import named events, build events, subtests, occurrence/selection, stderr, completion and process exit; never infer from Output text or package pass | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | Skip, repeat, incomplete, malformed stream, early pass then failure, build failure and filter tests |
| 13 | Explicit applicable invocation context and authored profiles | NEEDS ADAPTATION | Capture existing checked HEAD and allowlisted configuration, disclose preparation changes, keep independent consumer harness separate. Match exact declared context; do not observe profile selection from test JSON | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | Revision/module/replacement/context mismatch tests; distinct native and installed-consumer reports |
| 14 | Commit static reports; publish observations per CI invocation | HOLDS | Extend existing Go jobs with JSON capture and artifacts, preserving test selections/failures. Static page links to the existing workflow, never mutable latest-pass assertions | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | Exact-revision remote runs show contracts/native differences and retain failures |
| 15 | Validate/render both documents before publication; useful current rollout | HOLDS | Every discovered contribution gets details; annotate all current provider scopes/native grammars and useful current tests. Extend existing updater's generated-file set together | research owner | [PR126](https://github.com/Bitspark/bittype/pull/126) | Both-matrix detail regression, write-failure/no-publication tests and remote generated artifacts |
| 16 | Stronger isolation, cross-module source loading and dynamic-case precision when required | DEFERRED | Add worker-process isolation if untrusted extractors become a supported input; add explicit external source selections when a consumer needs dependency annotations; add definition/subtest qualifiers only with an ambiguous real use case. Current missing evidence stays explicit | research owner | — | Trigger is a supported consumer requiring those capabilities, not an undocumented promise |

Source discovery remains explicitly rooted in one selected module using the
existing Bittype repository layout. Making the tooling runner importable does
not imply arbitrary multi-module provider discovery. An external consumer can
compose extractors over that runner's selected source snapshot. Export-only
dependencies do not become ambient source/document inputs.

The initial advice evaluation found no additional existing architectural
violation. Subsequent implementation inspection exposed a concrete dependency
identity failure: display names could discharge an unrelated provider contract.
[ACX-007](https://github.com/Bitspark/platform/blob/main/docs/concepts/architecture-cases/0007-bittype-matrix-contract-identity.md)
records that counterexample and the canonical-key correction landed in PR126.
The separately tracked ACX-001 verification-composition defect remains outside
this task. The reporting extension must meet AC-02 through AC-05 from its first
slice; report facts and tests remain qualified under AC-07.

The renderer regression `TestRendererLinksCapabilityDetails` failed on the
pre-implementation renderer at `bb631f69280275f3c336c1bd0d0f15b498f7862f` in
[run 42246957504](https://github.com/Bitspark/bittype/actions/runs/42246957504).
It is retained in the importable report package for the corrected behavior.

### Delivered result

The initial implementation was delivered by [PR126](https://github.com/Bitspark/bittype/pull/126),
merged as `a582c5bc9327b38b75b16dba59fc558067ed1171`. Its final head `f209457`
passed [Go contracts and native qualification](https://github.com/Bitspark/bittype/actions/runs/42254515087)
and [documentation checks](https://github.com/Bitspark/bittype/actions/runs/42254514914).
The merged source tree is the tested tree. The tooling packages additionally passed
`go test ./tools/provider-matrix/...` locally on Windows.

The successful suite includes annotation/source/identity regressions, exact grammar
helper interpretation, the external-package struct extractor, bounded-input and
invalid-result failures, observation decoding and context association, and the
real failing-child capture test. Existing full Go tests, vet, TypeScript/Rust
typechecking, concurrency checks, native/Horos fuzzing, independently installed
public consumption and native compiler fixtures also passed.

The real invocation artifacts were inspected: repository, native and race/fuzz
captures match their independently selected reporting contexts. The successful
installed-consumer invocation is explicitly mismatched to repository declarations,
so it cannot credit their cases. Static inventories contain no current-pass claim.
Both matrix/detail outputs are committed, and the existing updater publishes the
pair. Row 16's explicit deferrals remain the supported boundary, not accepted
implementation left open. No runtime or semantic-provider behavior was added.

Independent delivery review reopened row 8 and ACX-007: package-qualified Go type
strings still omit ownership of private fields/methods in unnamed type arguments.
The original named-contract regressions passed but did not qualify this case.
The follow-up replaces string-based semantic comparisons with actual Go type
identity and adds distinct-owner unnamed-argument regressions. This qualification
is pending; the consultation is delivering again until the correction lands.
