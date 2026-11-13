# Advice for consult
**Research doc:** doc_0dae7dd0-2f2a-4904-aedd-e8a094145cf5.md
**Document:** doc_0dae7dd0-2f2a-4904-aedd-e8a094145cf5
**Advisor:** ChatGPT, model not read from the answer page
**Session:** nsess_72548a3d-b823-455b-9e27-d243c3ab62e0
**Run:** run_36aec98c-ac83-446d-a057-f81b327ed237
**Chat:** https://chatgpt.com/c/6af7161e-c138-83eb-9fa1-abc80b321376
**Verify:** nightfall consult verify run_36aec98c-ac83-446d-a057-f81b327ed237 --output <this file>

I’ve read the document in full, including the architectural rules, the nine classification cases, the competing design options, the work in flight, and the five questions in section 7.

My overall conclusion is that the proposed split is sound, but the document needs one important refinement: **a model is not the code that validates or interprets it. A model is the declaration against which that code operates.** This distinction provides a much stronger foundation for answering all five questions.

Below is my response, written as an architectural consultation rather than a review of the existing implementation. I distinguish recommendations from requirements already fixed by the owner's direction.

# Architectural consultation: Models and operations

**Subject:** Response to Research 0018 — Models and operations\
**Date:** 13 November 2026

## Executive recommendation

I recommend the following decisions:

| Question                | Recommendation                                                                                                                                                                                                                      |
| ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1. Model vs. operation  | Define a model as a declarative specification of admissible structures, concepts and relationships. Interpretation, validation and processing belong to operations.                                                                 |
| 2. One language's model | Keep syntax, analysed semantics and intended type constructs as distinct representations, with explicitly declared relationships. Adopt U1 initially, with selective future unification where semantic equivalence is demonstrated. |
| 3. Declarative models   | Adopt D1 as the first implementation direction. Declare structural schemas and fact/form schemas as data, with interpreters and generators remaining code.                                                                          |
| 4. Package placement    | Adopt a refined P3: two roots, with operation-independent carriers and declarations under `model/shared`, and operation contracts and mechanisms under `ops`.                                                                       |
| 5. What is missing      | Make model identity, versioning, provenance, evidence and coverage explicit. Do not confuse consistency, validity, capability and correctness.                                                                                      |

The most important architectural principle is:

**Models declare. Operations interpret, establish and transform. Evidence records what an operation has established about a particular model and artifact.**

This principle should govern package placement, API design, the declaration format and the migration sequence.

---

# 1. The line between model and operation

## 1.1 A better classification test

I would refine the test in section 6.1.

A model is a declaration of a domain's concepts, structures, relationships, identities and admissible forms.

An operation is an executable mechanism that consumes an artifact, a model, or evidence about an artifact, and produces a result, a decision or new evidence.

The distinction is not whether a function is pure, whether it is small, or whether it happens to be implemented in the same package as a type.

The test should be:

**Does this responsibility declare what is true or admissible, or does it determine whether something is true, interpret it, or produce something from it?**

This gives us four categories:

1. **Declaration:** What structures, names, relationships and constraints exist? Model.
2. **Interpretation:** What does a declared structure mean, or how should it be read? Operation.
3. **Evaluation:** Does a particular value conform to the declaration? Operation.
4. **Establishment:** Does an artifact correspond to an external specification, compiler judgment or another model? Operation.

A model may declare a constraint. It must not execute the mechanism that evaluates that constraint.

For example:

- A declaration that a Go field has one name and one type child is model.
- A grammar interpreter that checks whether a particular Go field has those properties is operation.
- A declaration that a name must be a valid Go identifier is model.
- The function that evaluates whether a particular string is a valid identifier is operation.

This also answers an important potential objection: executable code is not automatically an operation merely because it is code. A generator that turns a declarative schema into source code is an operation, while the generated source may itself be a model declaration.

Conversely, data is not automatically a model merely because it is represented as a struct or a map. An operation request, a receipt or a planner's working state is still operation-owned data.

## 1.2 The three owner rules

I interpret the three statements in the owner's direction as follows.

**A model contains no operation.**

A model package must not contain a mechanism that interprets or evaluates its declarations.

It may contain:

- Type declarations.
- Identifiers and vocabulary.
- Declarative constraints.
- Schema data.
- Identity and version information.
- Relationships between declared concepts.

It must not contain:

- Parsers.
- Checkers.
- Interpreters.
- Readers that reconstruct domain meaning from encoded artifacts.
- Planners.
- Builders.
- Validators that execute the declaration.

**A model imports no operation.**

The dependency direction is unambiguous:

`ops → model`

The reverse is forbidden.

An operation can depend on multiple models, but a model cannot depend on an operation to complete its own definition.

If two operations need a common value type, that type must be owned independently of either operation.

**Operations work on the model.**

An operation consumes a model as an explicit input to its mechanism.

This does not mean every operation must receive an entire model object. It can receive a particular schema, vocabulary or model identity.

The important point is that the operation does not own the declaration it interprets.

## 1.3 Consistency and correctness

I agree with the proposed distinction, with one refinement.

Consistency is agreement with the declared structure and rules of the model.

Correctness is agreement with something beyond the artifact's own declared structure.

The four levels should be retained, but level 3 should be described more carefully.

| Level                      | Question                                                                            | Nature               |
| -------------------------- | ----------------------------------------------------------------------------------- | -------------------- |
| 1. Structural consistency  | Is the generic data structure internally sound?                                     | Consistency          |
| 2. Model consistency       | Does the artifact conform to the language's declared model?                         | Consistency          |
| 3. Language validity       | Does the actual language implementation accept and correctly interpret the program? | External correctness |
| 4. Cross-model correctness | Does an artifact in one model correctly represent something in another model?       | External correctness |

This is a useful separation because it prevents the word "valid" from concealing four different guarantees.

A structurally valid tree is not necessarily a grammatically valid Go tree.

A grammatically valid Go tree is not necessarily a semantically valid Go program.

A semantically valid Go program is not necessarily a correct representation of a horos type.

These guarantees should never be inferred from one another.

## 1.4 Should checked state be separate?

Yes.

I strongly recommend that the checked state be represented separately from the unchecked artifact.

The tree implementation already demonstrates the right direction:

- An ordinary tree is data.
- A structure check establishes structural consistency.
- An opaque envelope represents the admitted, immutable result.

I would generalize this pattern.

Conceptually:

- `Artifact`: an ordinary model instance.
- `Evidence`: a record of an operation's assessment.
- `AdmittedArtifact`: an artifact associated with evidence satisfying a particular requirement.

Evidence should identify:

- The artifact or snapshot assessed.
- The model identity and version.
- The operation and its version or profile.
- The guarantees actually established.
- The outcome and any diagnostics.

A receipt should not simply mean "valid". It should mean that a particular operation established a particular set of guarantees about a particular artifact.

This is especially important when an artifact is checked against different versions of a model or different levels of guarantee.

The model does not become responsible for proving itself correct. The operation establishes evidence, and the evidence travels with the resulting admitted state.

## 1.5 Verdicts on cases (a) to (i)

| Case                                        | Verdict                                                                    | Reason                                                                                                                                                             |
| ------------------------------------------- | -------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| (a) Constructors and checks                 | Split. Declarations belong in model; executable validation belongs in ops. | Constructors that enforce invariants are executing validation.                                                                                                     |
| (b) Typed readers                           | Ops.                                                                       | They interpret graph structures and reconstruct domain-level views. The fact schemas they read belong in model.                                                    |
| (c) Structure checks and grammar tables     | Split.                                                                     | The generic tree structure and language grammar declarations are model; the interpreters and checkers are operations.                                              |
| (d) Generic carriers                        | Model.                                                                     | AST, CST, semantic graph and file data types are operation-independent artifact representations. Any packing, checking or publication mechanism must be separated. |
| (e) Contract packages                       | Split by responsibility.                                                   | Operation-independent value types belong in model; interfaces, request protocols, admission and execution belong in ops.                                           |
| (f) Engine, registry and planner            | Ops.                                                                       | They select, coordinate and execute contributions.                                                                                                                 |
| (g) Mapping identities and decision schemas | Model declarations.                                                        | They describe stable identities and the data contract of a contribution. Their selection, codec execution and interpretation belong in ops.                        |
| (h) Lexical rules                           | Split.                                                                     | The language's lexical constraints are model; functions that evaluate strings against them are operations.                                                         |
| (i) Intended forms                          | Split.                                                                     | The declared forms and their structure can be model. Construction, validation, interpretation and coverage decisions belong in ops.                                |

### The third category: contracts

I would not introduce a third architectural root.

There is, however, a useful distinction between a model and an operation contract.

A contract describes the boundary through which an operation is used. It includes the interface, protocol, request and result semantics.

A contract is not necessarily a model. An operation-specific request type does not become a language model simply because it contains data.

The correct question is whether a value has independent domain meaning outside the operation.

If it does, its declaration belongs in model.

If it exists only to express an operation's invocation, execution or result, it belongs in ops.

This is more precise than placing all value types in model or introducing a third root for everything difficult to classify.

---

# 2. What one language's model should be

## 2.1 One language, multiple representations

I would not define a language's model as one universal representation.

A language has several different representations, each answering a different question.

For Go, at least three are necessary:

1. **Syntax model:** What structures can represent Go type syntax?
2. **Semantic model:** What concepts and relationships does analysis establish about an actual Go program?
3. **Type-intent model:** What type constructs can be described as intended output before a Go program exists?

These are related models, not necessarily identical models.

The distinction is essential.

Syntax describes what can be written.

Semantics describes what a particular program means according to the language implementation.

Intent describes what a proposed program is supposed to construct.

An intended program does not yet have compiler-established identities or bindings. An analysed program does.

Therefore, the same name can sometimes be used across representations, but the representations must not be treated as interchangeable.

## 2.2 U1, U2 or U3?

My recommendation is U1 as the initial architecture, with selective adoption of U2 where equivalence can be established.

I would not adopt U3.

The reasons are different for each option.

**U1: Separate vocabularies with declared relationships.**

This is the safest starting point.

It acknowledges that syntax, analysis and intent have different identities, lifetimes and guarantees.

It also avoids forcing a universal representation onto language concepts that are not actually equivalent.

However, U1 must not mean that every relationship remains an undocumented switch statement. The relationship between a form, its syntax representation and its analysed meaning should be explicit and testable.

**U2: Shared names for intent and analysed concepts.**

This is worth adopting selectively.

When an intended construct and an analysed construct genuinely represent the same language-level concept, they should ideally share the concept's identity.

Their status can differ:

- Intended: a proposed instance of the concept.
- Established: a compiler-analysed instance of the concept.

But the data structures must preserve the differences in identity, structure, provenance and lifetime.

For example, a planned reference to a declaration is not the same thing as a compiler-established binding to an entity.

A common name must not conceal that distinction.

**U3: One description from which all three representations are generated.**

I would reject this as the initial design.

The syntax model contains concepts that have no corresponding type-level semantic or intent concept.

Examples include imports, file layout and parameter-list syntax.

Conversely, semantic analysis introduces identities, bindings and facts that cannot be reduced to syntax-tree structure.

The relationship between a form and syntax is not necessarily one-to-one. Rust and C++ can require multiple syntax nodes to express one intended construct.

A universal description would therefore become a complex meta-language for three different domains.

That is not justified by the present evidence.

## 2.3 Applying this to Go

The Go examples in section 3.3 are particularly revealing.

| Concept                                                  | Recommended treatment                                                                                                                                         |
| -------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Slice                                                    | Same conceptual identity across intent and analysis. Different representation structures remain.                                                              |
| Defined type                                             | Same language-level declaration concept. The record-only body restriction is a generation coverage limitation, not a defining property of Go.                 |
| Alias                                                    | Same conceptual identity where the declaration itself is concerned. The syntax, intent and analysis representations still differ.                             |
| Field                                                    | Keep the field declaration, field syntax node and analysed field entity explicitly related. Do not assume their structures are interchangeable.               |
| Primitive / builtin / basic                              | Do not unify without qualification. The intended primitive denotes a proposed type; the analysed builtin denotes an entity and its compiler-established type. |
| Record / struct                                          | Related, but not identical. The intended record omits properties such as Go field tags and embedding.                                                         |
| External                                                 | Keep distinct from a compiler-established external type. Intent names a destination; analysis establishes an identity.                                        |
| Type parameter                                           | Share the language-level concept where appropriate, but keep binder identity, constraints and references explicit.                                            |
| Generic application                                      | Preserve the distinction between a syntactic application, an intended application and an analysed application.                                                |
| Package                                                  | Keep package identity and package syntax distinct. A package is not simply another type form.                                                                 |
| Interface, union, signature, tuple and boolean operators | Retain as semantic concepts even where generation does not currently support them.                                                                            |

One particularly important correction concerns the current defined-type check.

The restriction that a Go defined type must have a record body is not part of Go's model.

It is a restriction of the current generation contribution.

That restriction must therefore be owned by the target operation, not by the language model.

## 2.4 What should the intended-forms package become?

I would move away from calling it simply a transform model.

It is a declaration of the target language's supported type-intent vocabulary.

Its forms are target-owned concepts, but its current constructors, checks and coverage decisions do not all belong together.

I would separate:

- The declaration of available type forms.
- The schema of each form.
- The generic representation of an intent expression.
- The operation that constructs and validates intent.
- The target interpreter that builds the syntax tree.

The first two are model.

The generic expression machinery is a shared carrier.

The last three are operations.

This also gives horos a natural place as a target: it can declare its own target-intent vocabulary without having to be treated as a special case of native generation.

## 2.5 Comparable architectural patterns

Several established systems illustrate why the distinctions matter.

**Roslyn (.NET compiler platform).**

Roslyn distinguishes syntax trees from semantic models. A syntax tree represents written program structure, while the semantic model provides compiler-established information about symbols, types and bindings.

This is a strong precedent for retaining separate syntax and semantic representations.

The lesson is not to collapse them into one AST.

**LLVM.**

LLVM uses a common intermediate representation and separate passes that transform or analyse it.

The IR is not the optimization pass. The pass operates on the IR and has its own responsibilities and guarantees.

This is a good example of the separation between a shared representation and the operations that consume it.

**Protocol Buffers.**

A schema declares message structures, while generators and runtimes interpret those declarations and produce language-specific representations.

This illustrates the value of making schemas independent of the code generated from them.

It also illustrates the importance of explicit schema identity and evolution.

The common lesson is that a representation, its interpretation and the evidence produced by processing it should remain distinct.

---

# 3. Declared models: how far should they be data?

## 3.1 My recommendation: D1 first

I recommend D1 as the first step.

It is the smallest direction that provides a meaningful improvement without introducing another description language prematurely.

The project already has substantial declarative content:

- Node kinds.
- Attributes and child roles.
- Grammar-table entries.
- Semantic vocabulary keys.
- Intended forms.
- Lexical tables.
- Simple mapping patterns.

The problem is not that everything is currently code.

The problem is that important structural contracts are distributed between data, functions, bridges and readers.

The most valuable first improvement is to make these contracts explicit.

## 3.2 What should be declared?

I would introduce three independently versioned schema families.

**A. Syntax schemas**

Declare:

- Node kinds.
- Attribute names and types.
- Required and optional attributes.
- Child roles.
- Cardinalities.
- Allowed child kinds.
- Reference roles.
- Ordering constraints where declaratively expressible.

**B. Semantic schemas**

Declare:

- Entity kinds.
- Term operators.
- Binding roles.
- Fact relations.
- Parameter sorts.
- Fact payload schemas.
- Operand cardinalities and relationships.
- Encoding identities and versions.

This is particularly important because the semantic graph currently declares vocabulary names but not the complete shapes of the facts written and read by the bridges.

**C. Intent schemas**

Declare:

- Form identities.
- Form data fields.
- Argument structure.
- Declaration forms.
- Parameter and reference structures.
- The relationships between intent forms and language-level concepts.

These three schemas describe different things and should not be forced into one schema type merely because they all contain names and fields.

## 3.3 What remains code?

Code remains necessary.

A declaration cannot replace every mechanism.

| Responsibility                              | Data or code?                                                             |
| ------------------------------------------- | ------------------------------------------------------------------------- |
| Node names and attributes                   | Data                                                                      |
| Child roles and cardinalities               | Data                                                                      |
| Simple grammar constraints                  | Data                                                                      |
| Complex relationships between nodes         | Code, unless a concrete need justifies extending the declarative language |
| Lexical character tables and reserved words | Data                                                                      |
| Evaluating lexical constraints              | Code                                                                      |
| Semantic entity and fact shapes             | Data                                                                      |
| Compiler-specific extraction                | Code                                                                      |
| Interpreting semantic fact payloads         | Code                                                                      |
| Intent form shapes                          | Data                                                                      |
| Building intent expressions                 | Operation code                                                            |
| Building concrete syntax trees              | Code                                                                      |
| Printing and formatting                     | Code                                                                      |
| Coverage and unsupported-feature policies   | Operation code                                                            |

The important principle is that code implementing an operation is not a defect in a declarative model.

A declarative model does not mean a system without code.

It means that the facts and constraints that can be stated as declarations are not unnecessarily embedded in executable implementations.

## 3.4 Why not D2?

The horos kind notation is interesting, but it does not currently express everything required.

It describes ordered tree structures, alternatives, repetitions, recursion and parameters.

It does not currently describe:

- Named child roles.
- General attribute predicates.
- Semantic fact shapes.
- The relationships between entities and terms.
- The full requirements of intended forms.

Using it now would either require a significant extension or leave essential parts of the model in code.

There is also a larger architectural issue: the shared model infrastructure would acquire a dependency on horos.

That conflicts with the existing architectural boundary.

D2 may become useful if a concrete case demonstrates that its expressive power is worth extending, but it should not be the starting point.

## 3.5 Why not D3?

D3 is a reasonable eventual possibility, but it needs a stronger justification.

A new schema language introduces:

- A parser.
- A specification.
- Versioning and compatibility rules.
- Tooling.
- A generator.
- A second implementation language boundary.
- A new mechanism that must itself be maintained.

D1 can already establish whether declarative descriptions provide sufficient value.

A new language should only be introduced when a concrete declaration cannot be expressed adequately in the existing mechanisms.

This follows the project's extension rule.

## 3.6 Identity and versioning

Every declared model needs an explicit identity and revision.

I recommend:

- Stable model identity.
- Explicit schema version.
- Language identity where applicable.
- Representation identity.
- Compatibility rules.
- A content digest for reproducible identification.

The identity must not depend on the Go package path.

A package move should not change the identity of the model.

Likewise, two declarations with identical contents but different semantic identities must not automatically be treated as the same model.

A revision must also distinguish between:

- An incompatible structural change.
- An additive declaration.
- A change in interpretation.
- A change in coverage.

Those are different events and should not be hidden behind one version number.

## 3.7 Keeping declarations honest

Declarations alone do not guarantee that implementations obey them.

I would establish a conformance suite:

1. The checker validates trees against the declared syntax schema.
2. Parsers produce trees that pass that checker.
3. Formatters produce trees that pass the same checker.
4. Bridges publish semantic graphs that pass generic integrity and provider-specific fact validation.
5. Semantic readers are tested against the declared fact schemas.
6. Targets produce trees consistent with the intent schemas.
7. Tests compare generated artifacts against their expected model identities and versions.

For the TypeScript and Rust bridges, which are not written in Go, the schema must be consumable independently of Go.

A Go generator can produce constants or schema bindings for the other implementation languages.

The declarations must remain the source of truth, rather than generated language bindings becoming independent authorities.

---

# 4. Carrying the split out

## 4.1 Package organization

I recommend a refined version of P3.

However, I would retain only two architectural roots.

The third conceptual category in P3, shared operation-independent values, should live under the model root rather than become another root.

A possible structure:

```
model/
    shared/
        tree/
        graph/
        files/
        representation/

    language/
        horos/
            syntax/
            semantics/
            intent/
            symbols/
        typescript/
            syntax/
            semantics/
            intent/
            symbols/
        golang/
            syntax/
            semantics/
            intent/
            symbols/
        rust/
            syntax/
            semantics/
            intent/
            symbols/
        cpp/
            syntax/
            semantics/

    relations/
        ...

    contributions/
        ...

ops/
    contracts/
        ...

    engine/
        ...

    parse/
        syntax/
        files/

    check/
        syntax/
        files/

    format/
        syntax/
        files/

    print/
        syntax/
        files/

    resolve/
        ...

    analyze/
        ...

    transform/
        ...

    verification/
        ...

    internal/
        ...

    integration/
        ...

    tools/
        ...
```

This is a conceptual organization, not a proposal to preserve every existing package name.

The important properties are:

- `model` contains no operation.
- `model` imports no operation.
- `ops` imports model declarations and carriers.
- Shared operation-independent data types have one owner.
- Language-specific models remain independently identifiable.
- Cross-language relationships have an explicit home.

## 4.2 What happens to the contracts?

The contract packages should be divided by responsibility.

Operation interfaces, requests, execution protocols and admission mechanisms belong in `ops/contracts`.

Generic carriers belong in `model/shared`.

Operation-independent types that describe persistent or transferable domain artifacts belong with their corresponding model.

A receipt, however, is not necessarily a model artifact merely because it contains data.

If its identity and meaning are inseparable from the guarantees of one operation, the receipt belongs to that operation's contract.

This avoids the mistake of moving all structs into model.

It also avoids allowing an operation's contract to become the owner of a domain model merely because that is where the first consumer happened to need it.

## 4.3 What happens to the language-by-stage matrix?

The existing stage-by-language matrix should remain under ops.

The owner's new root split changes the first architectural dimension, not the need for consistent organization among operation implementations.

Thus:

```
ops/
    transform/
        syntax/
            golang/
            rust/
            typescript/
```

The exact existing path pattern can remain useful below `ops`.

The model root should instead be organized around the kinds of models being declared, with language identity as a further dimension.

This is the significant conceptual difference.

The stage is an organizational property of an operation.

It is not an organizational property of a language model.

## 4.4 What about mappings and cross-language schemas?

The migration's new contribution identities and decision-record schemas should have explicit identities independent of the transform implementation.

The decision-record schema is a declaration.

The codec is an operation.

The mapping implementation is an operation contribution.

The relationship between a source family and a target language should not be mistaken for the model of either language.

I would place the stable declarations under a dedicated cross-language or contribution-model area of `model`, and the implementation under `ops/transform`.

This also accommodates transformations in either direction.

There should be no assumption that horos is always the source or that a programming language is always the target.

The shared contract should continue to describe source and target roles, with available transformations determined by registered contributions.

## 4.5 The first move: the intended-forms packages

I agree with the research owner's identification of the intended-forms packages as an immediate candidate for separation.

However, I would modify the proposed step in section 6.5.

Moving only coverage answers is insufficient under the owner's stricter rule that model contains no operations.

The checks are executable evaluation mechanisms too.

Therefore, I would move:

- Coverage decisions to the target.
- Well-formedness evaluation to an operation-side checker.
- The declarative form descriptions to the model side.
- The constructors and expression machinery according to whether they are model declarations or operation-specific construction mechanisms.

The model should not retain a function merely because the function evaluates an invariant that the model declares.

That would preserve the exact conceptual ambiguity the root split is intended to remove.

---

# 5. What the architecture is not yet seeing

I would add five distinctions to the document.

## 5.1 Identity is not representation

A model's identity must survive:

- A change of implementation language.
- A change of package path.
- A change of generated code.
- A change of the consumer.

This is particularly important for the Go and TypeScript semantic vocabularies.

A model identity must not be inferred from a Go marker type or a package name.

## 5.2 Well-formedness is not capability

The document already identifies this problem in the intended-forms packages.

An unsupported form is not necessarily an invalid form.

For example:

- An external type can be structurally well-formed.
- The Go target may not support it today.
- The Rust target may support only certain external crate destinations.

Those are three distinct statements.

I recommend that every operation distinguish:

- Invalid: the value violates the model.
- Unsupported: the value is valid, but the operation cannot process it.
- Incomplete: required information is unavailable.
- Rejected: a policy explicitly refuses the request.

This should be consistent across all operations.

## 5.3 A model is not necessarily complete

A language can have a declared model that covers more than any particular operation supports.

This is already visible in Go.

The syntax model contains pointers, arrays, maps and channels, even though current generation supports only a subset.

The semantic model contains concepts that mappings do not generate.

The model must not be narrowed merely to accommodate the present capabilities of a particular operation.

This follows directly from AC-06.

## 5.4 Evidence must be scoped

Evidence must establish precisely what was checked, against which model and under which profile.

A successful parse does not prove semantic validity.

A successful target build does not prove representation correctness.

A successful semantic analysis does not prove that a generated type correctly represents a source type.

This is why the four-level distinction is useful.

Evidence should not be transferrable to another model version without an explicit compatibility argument.

## 5.5 A model should not contain its own implementation history

A model describes a domain, not the mechanism that happens to process it today.

A Go model should not be considered more authoritative because its schema is expressed as Go structs.

A Rust model should not be considered different because a bridge is written in Rust.

The declaration and its implementation must be independently identifiable.

This becomes particularly important when the same model is consumed by several operations or implementation languages.

---

# 6. Recommended order of work

The migration should not be driven by the current directory layout or the amount of code that would need to move.

AC-09 makes that explicit.

I recommend the following sequence.

### Step 1 — Decide the architectural vocabulary

Before moving anything, formally define:

- Model.
- Operation.
- Operation contract.
- Model declaration.
- Model instance.
- Evidence.
- Admitted artifact.
- Coverage.

Approve the classification test from section 1.

Without this, different contributors will apply the same root split differently.

### Step 2 — Decide the package ownership rules

Approve:

- Two roots only.
- `model/shared` for operation-independent carriers.
- Language models organized by representation.
- Cross-language declarations under model.
- Operation interfaces and execution protocols under ops.
- No model-to-ops dependency.

Also decide whether the existing foundation, kernel, vocabulary, contribution and composition roles remain below the two roots.

My recommendation is that they do.

### Step 3 — Decide the representation relationships

Adopt U1 initially.

Explicitly document the relationships between syntax, semantics and intent.

Do not introduce U2 for a concept merely because its names happen to match.

Make the comparison criteria part of the design review.

### Step 4 — Define the declarative schema boundary

Adopt D1.

Start with syntax schemas and then add semantic fact schemas and intent schemas.

Keep the schema interpreter and all executable checks in ops.

Do not build a new schema language before the limitations of D1 are demonstrated.

### Step 5 — Correct misplaced responsibilities

Apply AC-10.

Move the responsibilities that are demonstrably in the wrong place as independent changes.

I would start with the intended-forms checks and coverage policies, then separate operation protocols from model-side carriers and declarations.

Do not wait for every operation to be implemented before establishing the correct ownership.

### Step 6 — Retarget the work in flight

I would not merge the work in flight into the present package structure merely to move it again.

The migration branch should be retargeted to the agreed architecture.

In particular:

- Complete the mapping contract gate against the new ownership boundaries.
- Build golang-to-typescript using the same source/target-neutral contract.
- Build typescript-to-horos against the new horos target model.
- Place the Go and Rust Files formatters under the final ops structure.

The two additional language pairs are especially valuable because they test whether the abstraction is genuinely direction-independent.

### Step 7 — Establish architectural enforcement

The import test should enforce:

1. No imports from ops into model.
2. No operation implementation under model.
3. No cross-language dependency except explicitly declared cross-language contributions and relationships.
4. No operation-specific private implementation imported by model.
5. No circular dependencies between model declarations.

The test should derive rules from the new paths and a small number of explicit exceptions.

However, the test cannot decide whether a package has been conceptually classified correctly.

That remains a review responsibility.

A green architecture test establishes that forbidden dependencies are absent. It does not prove that the package belongs in the correct root.

### Step 8 — Establish conformance and evolution tests

Once the new ownership is established:

- Test model declarations against their interpreters.
- Test all parsers and formatters against syntax schemas.
- Test semantic writers against fact schemas.
- Test semantic readers against the same schemas.
- Test targets against intent schemas.
- Test schema identity and version compatibility.
- Test that adding a contribution does not require changing shared orchestration.

This final point is particularly important for the migration gate.

---

# 7. Final answers to the five questions

**Question 1: What is the line between model and operation?**

A model declares the admissible concepts and structures. An operation interprets, evaluates, processes or establishes something about them.

Consistency is internal to a model. Correctness requires evidence against something external.

The four levels should remain, with level 3 explicitly treated as external language validity.

Checked states should be represented separately, with scoped evidence.

**Question 2: What is one language's model?**

A language has multiple related models: syntax, semantics and intent.

They should not be collapsed into one universal representation.

Adopt U1 initially, selectively unify concepts under U2 where identity, lifetime, authority, guarantees and failure behavior are genuinely compatible, and reject U3 as the starting point.

**Question 3: How far can a model be data?**

As far as its declarations can express the structures, constraints, identities and relationships without executing an interpreter.

Adopt D1 first.

Use separate declarative schemas for syntax, semantic facts and intent.

Keep parsing, evaluation, compiler integration, formatting and interpretation as operations.

**Question 4: How should the split be carried out?**

Use two roots, with shared operation-independent models and carriers under model, and contracts and executable mechanisms under ops.

Preserve the stage-by-language organization inside ops.

Make the architectural decisions before moving packages, then move misplaced responsibilities first and retarget the in-flight work rather than landing it against an obsolete layout.

**Question 5: What are we not seeing?**

The most important missing distinction is between:

- What a model declares.
- What an artifact contains.
- What an operation establishes.
- What an operation supports.
- What evidence actually proves.

These should become explicit architectural concepts.

---

# Closing recommendation

I would approve the owner's root direction without qualification.

I would also approve the proposed separation of language-independent and language-specific responsibilities, but I would not approve a migration that merely moves existing packages under two new directory names.

The architecture needs to establish a deeper ownership boundary.

The decisive principle is:

**A model must be complete as a declaration without depending on any operation that consumes it. An operation must be replaceable without changing the identity or meaning of the model it consumes.**

This does not mean that models must be universally expressive or that operations must be generic. It means each side has a clear and independent authority.

The model owns what exists and what its declared structures mean.

The operation owns how those declarations are interpreted, checked, transformed or used to establish evidence.

That is the separation I would build the new architecture around.

## **My strongest additional observation**

There is one issue I would prioritize above the individual package moves: **the current design sometimes treats a model's vocabulary, its schema, its interpreter and its coverage policy as if they were one architectural responsibility.**

The Go examples expose this particularly clearly:

- `FormRecord` is a declaration.
- The shape of a record expression is a schema.
- `CheckExpression` is an evaluator.
- The Go target's decision to refuse external types is a capability policy.
- The target's conversion of `record` into `struct-type` is an interpretation.

These are five different responsibilities, despite currently being closely related in the code.

Separating them gives the project a much more stable foundation for adding languages, transformations and compiler integrations.

**One qualification:** My proposed P3 refinement and D1-first approach are recommendations, not consequences already established by the owner's rules. They should be recorded as explicit design decisions before migration begins.