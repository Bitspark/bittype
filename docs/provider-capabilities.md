# Provider capability details

Generated from selected production declarations and explicit test associations.
These reports describe bounded implementation scope, not language completeness.

Inspection build: go1.27.3, linux/amd64; Go flags ; explicit build flags .

Definitions describe selected initializers; mutation, operation bodies and procedural semantics are not inferred. Profiles and expectations are authored assertions. A test function's observation is not separate execution of every associated case, nor evidence that a wrapper or every branch ran.

[Report contract](capability-reports.md) · [Provider matrices](providers.md) · [Existing CI invocations](https://github.com/Bitspark/bittype/actions/workflows/go.yml)

No execution observations are included in this committed inventory. Per-invocation CI artifacts contain observations with their actual source, build selection and process result. Missing observations do not mean a test never ran.

<a id="cap-MTE6dHJhbnNmb3JtZXI1NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJhbnNmb3Jtcy9ob3Jvcy10by1nby9nbyNOZXc2NTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3RyYW5zZm9ybS9nby5BU1RUcmFuc2Zvcm1lcjA6NTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nbyNMYW5ndWFnZTU6aG9yb3M1MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nbyNMYW5ndWFnZTI6Z28"></a>

## Horos → Go / github.com/Bitspark/bittype/transforms/horos-to-go/go#New

Contract: github.com/Bitspark/bittype/contracts/transform/go.ASTTransformer.

Public operation: [github.com/Bitspark/bittype/transforms/horos-to-go/go#New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-go/go/transform.go#L15). Receiver: .

Direct endpoints: github.com/Bitspark/bittype/syntax/horos/go.Language → github.com/Bitspark/bittype/syntax/golang/go.Language. Typed calls in the factory body establish this discovered edge; runtime reachability, returned-value flow, configured mappings and composed routes are not inferred.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/transforms/horos-to-go/go#New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-go/go/transform.go#L15) → [github.com/Bitspark/bittype/generation/syntax/go#NewTransformer](https://github.com/Bitspark/bittype/blob/main/generation/syntax/go/transform.go#L16) ([call](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-go/go/transform.go#L19)); receiver ; type arguments github.com/Bitspark/bittype/syntax/horos/go.Language, github.com/Bitspark/bittype/syntax/golang/go.Language.
- [github.com/Bitspark/bittype/generation/syntax/go#NewTransformer](https://github.com/Bitspark/bittype/blob/main/generation/syntax/go/transform.go#L16) → [github.com/Bitspark/bittype/contracts/transform/go#NewTransformer](https://github.com/Bitspark/bittype/blob/main/contracts/transform/go/adapter.go#L40) ([call](https://github.com/Bitspark/bittype/blob/main/generation/syntax/go/transform.go#L20)); receiver ; type arguments github.com/Bitspark/bittype/syntax/horos/go.Language, github.com/Bitspark/bittype/syntax/golang/go.Language.

Traversal boundary: typed transformer boundary.

### Declared target [github.com/Bitspark/bittype/transforms/horos-to-go/go#New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-go/go/transform.go#L15)

Authored scope: [docs/generation-composition.md#construction-and-selection](https://github.com/Bitspark/bittype/blob/main/docs/generation-composition.md#construction-and-selection); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-MTE6dHJhbnNmb3JtZXI1ODpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJhbnNmb3Jtcy9ob3Jvcy10by1jcHAvZ28jTmV3NjU6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy90cmFuc2Zvcm0vZ28uQVNUVHJhbnNmb3JtZXIwOjUyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28jTGFuZ3VhZ2U1Omhvcm9zNTA6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28jTGFuZ3VhZ2UzOmNwcA"></a>

## Horos → C++ / github.com/Bitspark/bittype/transforms/horos-to-cpp/go#New

Contract: github.com/Bitspark/bittype/contracts/transform/go.ASTTransformer.

Public operation: [github.com/Bitspark/bittype/transforms/horos-to-cpp/go#New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-cpp/go/transform.go#L15). Receiver: .

Direct endpoints: github.com/Bitspark/bittype/syntax/horos/go.Language → github.com/Bitspark/bittype/syntax/cpp/go.Language. Typed calls in the factory body establish this discovered edge; runtime reachability, returned-value flow, configured mappings and composed routes are not inferred.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/transforms/horos-to-cpp/go#New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-cpp/go/transform.go#L15) → [github.com/Bitspark/bittype/generation/syntax/go#NewTransformer](https://github.com/Bitspark/bittype/blob/main/generation/syntax/go/transform.go#L16) ([call](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-cpp/go/transform.go#L19)); receiver ; type arguments github.com/Bitspark/bittype/syntax/horos/go.Language, github.com/Bitspark/bittype/syntax/cpp/go.Language.
- [github.com/Bitspark/bittype/generation/syntax/go#NewTransformer](https://github.com/Bitspark/bittype/blob/main/generation/syntax/go/transform.go#L16) → [github.com/Bitspark/bittype/contracts/transform/go#NewTransformer](https://github.com/Bitspark/bittype/blob/main/contracts/transform/go/adapter.go#L40) ([call](https://github.com/Bitspark/bittype/blob/main/generation/syntax/go/transform.go#L20)); receiver ; type arguments github.com/Bitspark/bittype/syntax/horos/go.Language, github.com/Bitspark/bittype/syntax/cpp/go.Language.

Traversal boundary: typed transformer boundary.

### Declared target [github.com/Bitspark/bittype/transforms/horos-to-cpp/go#New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-cpp/go/transform.go#L15)

Authored scope: [docs/generation-composition.md#construction-and-selection](https://github.com/Bitspark/bittype/blob/main/docs/generation-composition.md#construction-and-selection); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-MTE6dHJhbnNmb3JtZXI1OTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJhbnNmb3Jtcy9ob3Jvcy10by1ydXN0L2dvI05ldzY1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvdHJhbnNmb3JtL2dvLkFTVFRyYW5zZm9ybWVyMDo1MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvI0xhbmd1YWdlNTpob3JvczUxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nbyNMYW5ndWFnZTQ6cnVzdA"></a>

## Horos → Rust / github.com/Bitspark/bittype/transforms/horos-to-rust/go#New

Contract: github.com/Bitspark/bittype/contracts/transform/go.ASTTransformer.

Public operation: [github.com/Bitspark/bittype/transforms/horos-to-rust/go#New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-rust/go/transform.go#L15). Receiver: .

Direct endpoints: github.com/Bitspark/bittype/syntax/horos/go.Language → github.com/Bitspark/bittype/syntax/rust/go.Language. Typed calls in the factory body establish this discovered edge; runtime reachability, returned-value flow, configured mappings and composed routes are not inferred.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/transforms/horos-to-rust/go#New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-rust/go/transform.go#L15) → [github.com/Bitspark/bittype/generation/syntax/go#NewTransformer](https://github.com/Bitspark/bittype/blob/main/generation/syntax/go/transform.go#L16) ([call](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-rust/go/transform.go#L19)); receiver ; type arguments github.com/Bitspark/bittype/syntax/horos/go.Language, github.com/Bitspark/bittype/syntax/rust/go.Language.
- [github.com/Bitspark/bittype/generation/syntax/go#NewTransformer](https://github.com/Bitspark/bittype/blob/main/generation/syntax/go/transform.go#L16) → [github.com/Bitspark/bittype/contracts/transform/go#NewTransformer](https://github.com/Bitspark/bittype/blob/main/contracts/transform/go/adapter.go#L40) ([call](https://github.com/Bitspark/bittype/blob/main/generation/syntax/go/transform.go#L20)); receiver ; type arguments github.com/Bitspark/bittype/syntax/horos/go.Language, github.com/Bitspark/bittype/syntax/rust/go.Language.

Traversal boundary: typed transformer boundary.

### Declared target [github.com/Bitspark/bittype/transforms/horos-to-rust/go#New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-rust/go/transform.go#L15)

Authored scope: [docs/generation-composition.md#construction-and-selection](https://github.com/Bitspark/bittype/blob/main/docs/generation-composition.md#construction-and-selection); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-MTE6dHJhbnNmb3JtZXI2NTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJhbnNmb3Jtcy9ob3Jvcy10by10eXBlc2NyaXB0L2dvI05ldzY1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvdHJhbnNmb3JtL2dvLkFTVFRyYW5zZm9ybWVyMDo1MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvI0xhbmd1YWdlNTpob3JvczU3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nbyNMYW5ndWFnZTEwOnR5cGVzY3JpcHQ"></a>

## Horos → TypeScript / github.com/Bitspark/bittype/transforms/horos-to-typescript/go#New

Contract: github.com/Bitspark/bittype/contracts/transform/go.ASTTransformer.

Public operation: [github.com/Bitspark/bittype/transforms/horos-to-typescript/go#New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-typescript/go/transform.go#L15). Receiver: .

Direct endpoints: github.com/Bitspark/bittype/syntax/horos/go.Language → github.com/Bitspark/bittype/syntax/typescript/go.Language. Typed calls in the factory body establish this discovered edge; runtime reachability, returned-value flow, configured mappings and composed routes are not inferred.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/transforms/horos-to-typescript/go#New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-typescript/go/transform.go#L15) → [github.com/Bitspark/bittype/generation/syntax/go#NewTransformer](https://github.com/Bitspark/bittype/blob/main/generation/syntax/go/transform.go#L16) ([call](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-typescript/go/transform.go#L19)); receiver ; type arguments github.com/Bitspark/bittype/syntax/horos/go.Language, github.com/Bitspark/bittype/syntax/typescript/go.Language.
- [github.com/Bitspark/bittype/generation/syntax/go#NewTransformer](https://github.com/Bitspark/bittype/blob/main/generation/syntax/go/transform.go#L16) → [github.com/Bitspark/bittype/contracts/transform/go#NewTransformer](https://github.com/Bitspark/bittype/blob/main/contracts/transform/go/adapter.go#L40) ([call](https://github.com/Bitspark/bittype/blob/main/generation/syntax/go/transform.go#L20)); receiver ; type arguments github.com/Bitspark/bittype/syntax/horos/go.Language, github.com/Bitspark/bittype/syntax/typescript/go.Language.

Traversal boundary: typed transformer boundary.

### Declared target [github.com/Bitspark/bittype/transforms/horos-to-typescript/go#New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-typescript/go/transform.go#L15)

Authored scope: [docs/generation-composition.md#construction-and-selection](https://github.com/Bitspark/bittype/blob/main/docs/generation-composition.md#construction-and-selection); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODppZGVudGl0eTE6IzU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvbGFuZ3VhZ2UvZ28uUHJvdmlkZXI1MDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nby5MYW5ndWFnZTUwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvY3BwL2dvI0xhbmd1YWdlMzpjcHA"></a>

## C++ language identity

Contract: github.com/Bitspark/bittype/contracts/language/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

### Declared target [github.com/Bitspark/bittype/syntax/cpp/go#Language.ID](https://github.com/Bitspark/bittype/blob/main/syntax/cpp/go/nodes.go#L16)

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

### Declared target [github.com/Bitspark/bittype/syntax/cpp/go#Language.Name](https://github.com/Bitspark/bittype/blob/main/syntax/cpp/go/nodes.go#L20)

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODppZGVudGl0eTE6IzU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvbGFuZ3VhZ2UvZ28uUHJvdmlkZXI1MTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28uTGFuZ3VhZ2U1MTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28jTGFuZ3VhZ2U0OnJ1c3Q"></a>

## Rust language identity

Contract: github.com/Bitspark/bittype/contracts/language/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

### Declared target [github.com/Bitspark/bittype/syntax/rust/go#Language.ID](https://github.com/Bitspark/bittype/blob/main/syntax/rust/go/nodes.go#L16)

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

### Declared target [github.com/Bitspark/bittype/syntax/rust/go#Language.Name](https://github.com/Bitspark/bittype/blob/main/syntax/rust/go/nodes.go#L20)

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODppZGVudGl0eTE6IzU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvbGFuZ3VhZ2UvZ28uUHJvdmlkZXI1MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlNTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nbyNMYW5ndWFnZTU6aG9yb3M"></a>

## Horos language identity

Contract: github.com/Bitspark/bittype/contracts/language/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

### Declared target [github.com/Bitspark/bittype/syntax/horos/go#Language.ID](https://github.com/Bitspark/bittype/blob/main/syntax/horos/go/nodes.go#L13)

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

### Declared target [github.com/Bitspark/bittype/syntax/horos/go#Language.Name](https://github.com/Bitspark/bittype/blob/main/syntax/horos/go/nodes.go#L16)

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODppZGVudGl0eTE6IzU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvbGFuZ3VhZ2UvZ28uUHJvdmlkZXI1MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nby5MYW5ndWFnZTUzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvZ29sYW5nL2dvI0xhbmd1YWdlMjpnbw"></a>

## Go language identity

Contract: github.com/Bitspark/bittype/contracts/language/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

### Declared target [github.com/Bitspark/bittype/syntax/golang/go#Language.ID](https://github.com/Bitspark/bittype/blob/main/syntax/golang/go/nodes.go#L16)

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

### Declared target [github.com/Bitspark/bittype/syntax/golang/go#Language.Name](https://github.com/Bitspark/bittype/blob/main/syntax/golang/go/nodes.go#L20)

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODppZGVudGl0eTE6IzU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvbGFuZ3VhZ2UvZ28uUHJvdmlkZXI1NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28uTGFuZ3VhZ2U1NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28jTGFuZ3VhZ2UxMDp0eXBlc2NyaXB0"></a>

## TypeScript language identity

Contract: github.com/Bitspark/bittype/contracts/language/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

### Declared target [github.com/Bitspark/bittype/syntax/typescript/go#Language.ID](https://github.com/Bitspark/bittype/blob/main/syntax/typescript/go/nodes.go#L16)

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

### Declared target [github.com/Bitspark/bittype/syntax/typescript/go#Language.Name](https://github.com/Bitspark/bittype/blob/main/syntax/typescript/go/nodes.go#L20)

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjU3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9nbyNQcm92aWRlci5QcmludDE5NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3ByaW50L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zb3VyY2UvZ28uU291cmNlW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvY3BwL2dvLkxhbmd1YWdlXSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2ZpbGVzL2dvLkZpbGVzXTEwMzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJpbnQvZmlsZXMvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28uTGFuZ3VhZ2VdNTA6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28jTGFuZ3VhZ2UzOmNwcA"></a>

## C++ / github.com/Bitspark/bittype/print/files/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/source/go.Source&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;, github.com/Bitspark/bittype/files/go.Files&#93;.

Public operation: [github.com/Bitspark/bittype/print/files/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L13). Receiver: github.com/Bitspark/bittype/print/files/go.Provider&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider + github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;, string&#93;.

### Declared target [github.com/Bitspark/bittype/print/files/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L13)

Authored scope: [docs/source.md#symmetric-role-contracts](https://github.com/Bitspark/bittype/blob/main/docs/source.md#symmetric-role-contracts); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjU3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9nbyNQcm92aWRlci5QcmludDE5ODpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3ByaW50L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zb3VyY2UvZ28uU291cmNlW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nby5MYW5ndWFnZV0sIGdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9maWxlcy9nby5GaWxlc10xMDQ6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3ByaW50L2ZpbGVzL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nby5MYW5ndWFnZV01MTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28jTGFuZ3VhZ2U0OnJ1c3Q"></a>

## Rust / github.com/Bitspark/bittype/print/files/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/source/go.Source&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;, github.com/Bitspark/bittype/files/go.Files&#93;.

Public operation: [github.com/Bitspark/bittype/print/files/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L13). Receiver: github.com/Bitspark/bittype/print/files/go.Provider&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider + github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;, string&#93;.

### Declared target [github.com/Bitspark/bittype/print/files/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L13)

Authored scope: [docs/source.md#symmetric-role-contracts](https://github.com/Bitspark/bittype/blob/main/docs/source.md#symmetric-role-contracts); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjU3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9nbyNQcm92aWRlci5QcmludDE5OTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3ByaW50L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zb3VyY2UvZ28uU291cmNlW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvZmlsZXMvZ28uRmlsZXNdMTA1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXTUyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28jTGFuZ3VhZ2U1Omhvcm9z"></a>

## Horos / github.com/Bitspark/bittype/print/files/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/source/go.Source&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;, github.com/Bitspark/bittype/files/go.Files&#93;.

Public operation: [github.com/Bitspark/bittype/print/files/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L13). Receiver: github.com/Bitspark/bittype/print/files/go.Provider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider + github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;, string&#93;.

### Declared target [github.com/Bitspark/bittype/print/files/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L13)

Authored scope: [docs/source.md#symmetric-role-contracts](https://github.com/Bitspark/bittype/blob/main/docs/source.md#symmetric-role-contracts); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| multi-file-pagination | receipt-preserving-round-trips | file-content-and-placement-preserved | unspecified (authored profile absent) | [TestFilesPaginationAndReceiptPreservingRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/files_test.go#L56) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjU3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9nbyNQcm92aWRlci5QcmludDIwMDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3ByaW50L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zb3VyY2UvZ28uU291cmNlW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvZ29sYW5nL2dvLkxhbmd1YWdlXSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2ZpbGVzL2dvLkZpbGVzXTEwNjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJpbnQvZmlsZXMvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28uTGFuZ3VhZ2VdNTM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28jTGFuZ3VhZ2UyOmdv"></a>

## Go / github.com/Bitspark/bittype/print/files/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/source/go.Source&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;, github.com/Bitspark/bittype/files/go.Files&#93;.

Public operation: [github.com/Bitspark/bittype/print/files/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L13). Receiver: github.com/Bitspark/bittype/print/files/go.Provider&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider + github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;, string&#93;.

### Declared target [github.com/Bitspark/bittype/print/files/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L13)

Authored scope: [docs/source.md#symmetric-role-contracts](https://github.com/Bitspark/bittype/blob/main/docs/source.md#symmetric-role-contracts); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjU3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9nbyNQcm92aWRlci5QcmludDIwNDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3ByaW50L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zb3VyY2UvZ28uU291cmNlW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nby5MYW5ndWFnZV0sIGdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9maWxlcy9nby5GaWxlc10xMTA6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3ByaW50L2ZpbGVzL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nby5MYW5ndWFnZV01NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28jTGFuZ3VhZ2UxMDp0eXBlc2NyaXB0"></a>

## TypeScript / github.com/Bitspark/bittype/print/files/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/source/go.Source&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;, github.com/Bitspark/bittype/files/go.Files&#93;.

Public operation: [github.com/Bitspark/bittype/print/files/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L13). Receiver: github.com/Bitspark/bittype/print/files/go.Provider&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider + github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;, string&#93;.

### Declared target [github.com/Bitspark/bittype/print/files/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L13)

Authored scope: [docs/source.md#symmetric-role-contracts](https://github.com/Bitspark/bittype/blob/main/docs/source.md#symmetric-role-contracts); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9jcHAvZ28jUHJvdmlkZXIuQ2hlY2sxMDc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9jaGVjay9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nby5MYW5ndWFnZV01MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJvdmlkZXIvY3BwL2dvLlByb3ZpZGVyNTA6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28jTGFuZ3VhZ2UzOmNwcA"></a>

## C++ / github.com/Bitspark/bittype/provider/cpp/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/provider/cpp/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L22). Receiver: github.com/Bitspark/bittype/provider/cpp/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/cpp/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L22) → [github.com/Bitspark/bittype/check/syntax/cpp/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L57) ([call](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L22)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/check/syntax/cpp/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L57)

Related declaration; its cases and observations do not qualify the public wrapper.

Definition: [github.com/Bitspark/bittype/check/syntax/cpp/go#rules](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L29); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### translation-unit ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L30))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key translation-unit.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| include | 0 | unbounded | include |
| declaration | 0 | unbounded | namespace, record, type-alias, enum |

No feature-key case association at this exact authored profile.

#### include ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L31))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key include.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| path | true | opaque: Predicate present; behavior not interpreted |
| system | false | true, false |

No feature-key case association at this exact authored profile.

#### namespace ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L32))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key namespace.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| declaration | 0 | unbounded | namespace, record, type-alias, enum |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### record ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L33))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key record.

Roles: 3 declared. Attributes: unknown: Attribute collection could not be decoded. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| template-parameter | 0 | unbounded | type-parameter |
| base | 0 | unbounded | base |
| member | 0 | unbounded | field, method, destructor, type-alias, record, enum |

No feature-key case association at this exact authored profile.

#### type-alias ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L34))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-alias.

Roles: 2 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| template-parameter | 0 | unbounded | type-parameter |
| type | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| access | false | public, protected, private |

No feature-key case association at this exact authored profile.

#### enum ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L35))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key enum.

Roles: 2 declared. Attributes: unknown: Attribute collection could not be decoded. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| underlying | 0 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |
| enumerator | 0 | unbounded | enumerator |

No feature-key case association at this exact authored profile.

#### enumerator ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L36))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key enumerator.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| value | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### base ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L37))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key base.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| access | false | public, protected, private |
| virtual | false | true, false |

No feature-key case association at this exact authored profile.

#### field ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L38))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key field.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| access | false | public, protected, private |

No feature-key case association at this exact authored profile.

#### method ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L39))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key method.

Roles: 2 declared. Attributes: unknown: Attribute collection could not be decoded. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 0 | unbounded | parameter |
| result | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

#### destructor ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L40))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key destructor.

Roles: 0 declared. Attributes: 4 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| access | false | public, protected, private |
| virtual | false | true, false |
| pure | false | true, false |
| noexcept | false | true, false |

No feature-key case association at this exact authored profile.

#### parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L41))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L42))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-parameter.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| default | 0 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| parameter-key | false | typename, class |

No feature-key case association at this exact authored profile.

#### named-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L43))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key named-type.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| segment | 1 | unbounded | name-segment |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| global | false | true, false |
| typename | false | true, false |

No feature-key case association at this exact authored profile.

#### name-segment ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L44))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key name-segment.

Roles: 1 declared. Attributes: 2 declared. References: target.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| argument | 0 | unbounded | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| template | false | true, false |

No feature-key case association at this exact authored profile.

#### builtin-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L45))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key builtin-type.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### pointer-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L46))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key pointer-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

#### lvalue-reference-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L47))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key lvalue-reference-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

#### rvalue-reference-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L48))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key rvalue-reference-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

#### cv-qualified-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L49))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key cv-qualified-type.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| const | false | true, false |
| volatile | false | true, false |

No feature-key case association at this exact authored profile.

#### array-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L50))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key array-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| length | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### function-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L51))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key function-type.

Roles: 2 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 0 | unbounded | parameter |
| result | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| noexcept | false | true, false |
| variadic | false | true, false |

No feature-key case association at this exact authored profile.

#### member-pointer-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L52))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key member-pointer-type.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| owner | 1 | 1 | named-type |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

Authored scope: [docs/syntax.md](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/provider/cpp/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L22)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| ast-admission | malformed-carriers | reject | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L191) (external) |

<a id="cap-ODpwcm92aWRlcjU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9jcHAvZ28jUHJvdmlkZXIuUGFyc2UyMDc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wYXJzZS9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3BhcnNlL2dvLklucHV0LCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nby5MYW5ndWFnZV1dNTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3Byb3ZpZGVyL2NwcC9nby5Qcm92aWRlcjUwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvY3BwL2dvI0xhbmd1YWdlMzpjcHA"></a>

## C++ / github.com/Bitspark/bittype/provider/cpp/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/provider/cpp/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L19). Receiver: github.com/Bitspark/bittype/provider/cpp/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/cpp/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L19) → [github.com/Bitspark/bittype/parse/syntax/cpp/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/cpp/go/parse.go#L25) ([call](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L20)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/cpp/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/cpp/go/parse.go#L25)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/provider/cpp/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L19)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L118) (external) |
| source-boundary | unsupported-and-malformed-fixtures | located-error-without-output | unspecified (authored profile absent) | [TestUnsupportedAndMalformedSources](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L140) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9jcHAvZ28jUHJvdmlkZXIuUHJpbnQxNjE6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wcmludC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nby5MYW5ndWFnZV0sIHN0cmluZ101MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJvdmlkZXIvY3BwL2dvLlByb3ZpZGVyNTA6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28jTGFuZ3VhZ2UzOmNwcA"></a>

## C++ / github.com/Bitspark/bittype/provider/cpp/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/provider/cpp/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L26). Receiver: github.com/Bitspark/bittype/provider/cpp/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/cpp/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L26) → [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L27)); receiver ; type arguments github.com/Bitspark/bittype/syntax/cpp/go.Language.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/provider/cpp/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L26)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L104) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L118) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvZ28jUHJvdmlkZXIuUHJpbnQxNjE6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wcmludC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nby5MYW5ndWFnZV0sIHN0cmluZ10xMDQ6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3ByaW50L3N5bnRheC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nby5MYW5ndWFnZV01MDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nbyNMYW5ndWFnZTM6Y3Bw"></a>

## C++ / github.com/Bitspark/bittype/print/syntax/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/print/syntax/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L25). Receiver: github.com/Bitspark/bittype/print/syntax/go.Provider&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/print/syntax/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L25) → [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L26)); receiver ; type arguments L.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvZ28jUHJvdmlkZXIuUHJpbnQxNjI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wcmludC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28uTGFuZ3VhZ2VdLCBzdHJpbmddMTA1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvLkxhbmd1YWdlXTUxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nbyNMYW5ndWFnZTQ6cnVzdA"></a>

## Rust / github.com/Bitspark/bittype/print/syntax/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/print/syntax/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L25). Receiver: github.com/Bitspark/bittype/print/syntax/go.Provider&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/print/syntax/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L25) → [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L26)); receiver ; type arguments L.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvZ28jUHJvdmlkZXIuUHJpbnQxNjM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wcmludC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXSwgc3RyaW5nXTEwNjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJpbnQvc3ludGF4L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdNTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nbyNMYW5ndWFnZTU6aG9yb3M"></a>

## Horos / github.com/Bitspark/bittype/print/syntax/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/print/syntax/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L25). Receiver: github.com/Bitspark/bittype/print/syntax/go.Provider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/print/syntax/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L25) → [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L26)); receiver ; type arguments L.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvZ28jUHJvdmlkZXIuUHJpbnQxNjQ6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wcmludC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nby5MYW5ndWFnZV0sIHN0cmluZ10xMDc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3ByaW50L3N5bnRheC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nby5MYW5ndWFnZV01MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nbyNMYW5ndWFnZTI6Z28"></a>

## Go / github.com/Bitspark/bittype/print/syntax/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/print/syntax/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L25). Receiver: github.com/Bitspark/bittype/print/syntax/go.Provider&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/print/syntax/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L25) → [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L26)); receiver ; type arguments L.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| go-concrete-pipeline | pagination-and-interface | source-retained | unspecified (authored profile absent) | [TestGoConcretePipeline](https://github.com/Bitspark/bittype/blob/main/integration/go/text_test.go#L15) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvZ28jUHJvdmlkZXIuUHJpbnQxNjg6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wcmludC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28uTGFuZ3VhZ2VdLCBzdHJpbmddMTExOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvLkxhbmd1YWdlXTU3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nbyNMYW5ndWFnZTEwOnR5cGVzY3JpcHQ"></a>

## TypeScript / github.com/Bitspark/bittype/print/syntax/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/print/syntax/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L25). Receiver: github.com/Bitspark/bittype/print/syntax/go.Provider&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/print/syntax/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L25) → [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/provider.go#L26)); receiver ; type arguments L.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjU5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9jcHAvZ28jUHJvdmlkZXIuRm9ybWF0MjUzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvZm9ybWF0L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9hc3QvZ28uVmlld1tnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nby5MYW5ndWFnZV0sIGdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9jc3QvZ28uQ1NUW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvY3BwL2dvLkxhbmd1YWdlXV01MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJvdmlkZXIvY3BwL2dvLlByb3ZpZGVyNTA6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28jTGFuZ3VhZ2UzOmNwcA"></a>

## C++ / github.com/Bitspark/bittype/provider/cpp/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/provider/cpp/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L23). Receiver: github.com/Bitspark/bittype/provider/cpp/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/cpp/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L23) → [github.com/Bitspark/bittype/format/syntax/cpp/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/cpp/go/format.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L24)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/cpp/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/cpp/go/format.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/provider/cpp/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/cpp/go/syntax.go#L23)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L104) (external) |
| ast-admission | malformed-carriers | rejection-without-output | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L191) (external) |
| native-compiler-qualification | formatted-fixture-contract | compiler-fixture-contract | unspecified (authored profile absent) | [TestNativeCompilerFixtures](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L263) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L118) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjU5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ydXN0L2dvI1Byb3ZpZGVyLkNoZWNrMTA4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvY2hlY2svZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvLkxhbmd1YWdlXTUzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ydXN0L2dvLlByb3ZpZGVyNTE6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvI0xhbmd1YWdlNDpydXN0"></a>

## Rust / github.com/Bitspark/bittype/provider/rust/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/provider/rust/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L22). Receiver: github.com/Bitspark/bittype/provider/rust/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/rust/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L22) → [github.com/Bitspark/bittype/check/syntax/rust/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L72) ([call](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L22)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/check/syntax/rust/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L72)

Related declaration; its cases and observations do not qualify the public wrapper.

Definition: [github.com/Bitspark/bittype/check/syntax/rust/go#rules](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L31); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### file ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L32))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key file.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| item | 0 | unbounded | module, use, type-alias, struct, tuple-struct, enum, trait, trait-impl |

No feature-key case association at this exact authored profile.

#### module ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L33))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key module.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| item | 0 | unbounded | module, use, type-alias, struct, tuple-struct, enum, trait, trait-impl |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### use ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L34))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key use.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| path | 1 | 1 | path-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-alias ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L35))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-alias.

Roles: 3 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| where | 0 | unbounded | where-predicate |
| type | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### struct ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L36))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key struct.

Roles: 3 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| where | 0 | unbounded | where-predicate |
| field | 0 | unbounded | field |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### tuple-struct ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L37))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-struct.

Roles: 3 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| where | 0 | unbounded | where-predicate |
| field | 0 | unbounded | tuple-field |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### enum ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L38))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key enum.

Roles: 3 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| where | 0 | unbounded | where-predicate |
| variant | 0 | unbounded | unit-variant, tuple-variant, record-variant |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### unit-variant ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L39))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key unit-variant.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### tuple-variant ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L40))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-variant.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| field | 0 | unbounded | tuple-field |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### record-variant ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L41))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key record-variant.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| field | 0 | unbounded | field |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### field ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L42))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key field.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### tuple-field ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L43))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-field.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### trait ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L44))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key trait.

Roles: 4 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| bound | 0 | unbounded | path-type, lifetime-reference |
| where | 0 | unbounded | where-predicate |
| member | 0 | unbounded | method-signature, associated-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### trait-impl ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L45))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key trait-impl.

Roles: 5 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| trait | 1 | 1 | path-type |
| subject | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |
| where | 0 | unbounded | where-predicate |
| member | 0 | unbounded | type-alias |

No feature-key case association at this exact authored profile.

#### associated-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L46))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key associated-type.

Roles: 3 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| bound | 0 | unbounded | path-type, lifetime-reference, relaxed-bound |
| where | 0 | unbounded | where-predicate |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### method-signature ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L47))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key method-signature.

Roles: 5 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| where | 0 | unbounded | where-predicate |
| receiver | 0 | 1 | receiver |
| parameter | 0 | unbounded | parameter |
| result | 0 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### receiver ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L48))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key receiver.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| lifetime | 0 | 1 | lifetime-reference |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| receiver-form | true | value, shared, mutable |

No feature-key case association at this exact authored profile.

#### parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L49))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L50))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-parameter.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| bound | 0 | unbounded | path-type, lifetime-reference, relaxed-bound |
| default | 0 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### lifetime-parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L51))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key lifetime-parameter.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| bound | 0 | unbounded | lifetime-reference |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### lifetime-reference ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L52))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key lifetime-reference.

Roles: 0 declared. Attributes: 1 declared. References: target.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### where-predicate ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L53))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key where-predicate.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| subject | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type, lifetime-reference |
| bound | 1 | unbounded | path-type, lifetime-reference, relaxed-bound |

No feature-key case association at this exact authored profile.

#### relaxed-bound ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L54))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key relaxed-bound.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| path | 1 | 1 | path-type |

No feature-key case association at this exact authored profile.

#### path-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L55))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key path-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| segment | 1 | unbounded | path-segment |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| absolute | false | true, false |

No feature-key case association at this exact authored profile.

#### qualified-path-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L56))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key qualified-path-type.

Roles: 3 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| subject | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |
| trait | 0 | 1 | path-type |
| segment | 1 | unbounded | path-segment |

No feature-key case association at this exact authored profile.

#### path-segment ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L57))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key path-segment.

Roles: 1 declared. Attributes: 1 declared. References: target.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| argument | 0 | unbounded | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type, lifetime-reference, associated-type-binding |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### associated-type-binding ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L58))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key associated-type-binding.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| path | 1 | 1 | path-segment |
| type | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

No feature-key case association at this exact authored profile.

#### primitive-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L59))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key primitive-type.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | bool, char, str, f32, f64, i8, i16, i32, i64, i128, isize, u8, u16, u32, u64, u128, usize |

No feature-key case association at this exact authored profile.

#### reference-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L60))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key reference-type.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |
| lifetime | 0 | 1 | lifetime-reference |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| mutable | false | true, false |

No feature-key case association at this exact authored profile.

#### raw-pointer-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L61))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key raw-pointer-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| mutable | false | true, false |

No feature-key case association at this exact authored profile.

#### slice-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L62))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key slice-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

No feature-key case association at this exact authored profile.

#### array-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L63))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key array-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| length | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### tuple-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L64))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 0 | unbounded | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

No feature-key case association at this exact authored profile.

#### never-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L65))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key never-type.

Roles: 0 declared. Attributes: 0 declared. References: none.

No feature-key case association at this exact authored profile.

#### function-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L66))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key function-type.

Roles: 2 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 0 | unbounded | parameter |
| result | 0 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| unsafe | false | true, false |
| abi | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### trait-object-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L67))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key trait-object-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| bound | 1 | unbounded | path-type, lifetime-reference |

No feature-key case association at this exact authored profile.

Authored scope: [docs/syntax.md#rust](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#rust); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| associated-type-projections | independent-ast-and-references | accept | unspecified (authored profile absent) | [TestAssociatedTypeAdmission](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check_test.go#L36) (external) |
| associated-type-projections | invalid-shapes-and-ordering | located-rejection | unspecified (authored profile absent) | [TestInvalidAssociatedTypeShapes](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check_test.go#L53) (external) |
| generic-associated-declarations | independent-ast | accept | unspecified (authored profile absent) | [TestGenericAssociatedDeclarationAdmission](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/declarations_test.go#L45) (external) |
| generic-associated-declarations | relaxed-bound-context | context-sensitive-admission | unspecified (authored profile absent) | [TestRelaxedBoundContext](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/declarations_test.go#L53) (external) |
| generic-associated-declarations | invalid-declaration-contexts | located-rejection | unspecified (authored profile absent) | [TestGenericAssociatedContextRules](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/declarations_test.go#L73) (external) |

### Declared target [github.com/Bitspark/bittype/provider/rust/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L22)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| ast-admission | malformed-carriers | reject | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L191) (external) |

<a id="cap-ODpwcm92aWRlcjU5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ydXN0L2dvI1Byb3ZpZGVyLlBhcnNlMjA4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcGFyc2UvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wYXJzZS9nby5JbnB1dCwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvLkxhbmd1YWdlXV01MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJvdmlkZXIvcnVzdC9nby5Qcm92aWRlcjUxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nbyNMYW5ndWFnZTQ6cnVzdA"></a>

## Rust / github.com/Bitspark/bittype/provider/rust/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/provider/rust/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L19). Receiver: github.com/Bitspark/bittype/provider/rust/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/rust/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L19) → [github.com/Bitspark/bittype/parse/syntax/rust/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse.go#L25) ([call](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L20)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/rust/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse.go#L25)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md#rust](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#rust); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| generic-associated-declarations | declaration-and-impl-structure | structure-preserved | unspecified (authored profile absent) | [TestGenericAssociatedDeclarationStructure](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/declarations_test.go#L16) (external) |
| generic-associated-declarations | unresolved-syntax | accept | unspecified (authored profile absent) | [TestGenericAssociatedSyntaxWithoutResolution](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/declarations_test.go#L46) (external) |
| generic-associated-declarations | malformed-declarations | located-error-without-output | unspecified (authored profile absent) | [TestMalformedGenericAssociatedDeclarations](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/declarations_test.go#L67) (external) |
| type-only-impl-scope | unsupported-impl-forms | unsupported-without-output | unspecified (authored profile absent) | [TestImplScopeRejection](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/declarations_test.go#L94) (external) |
| associated-type-projections | structured-projections-and-bindings | structure-preserved | unspecified (authored profile absent) | [TestAssociatedTypeStructure](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse_test.go#L17) (external) |
| associated-type-projections | unresolved-projection-syntax | accept | unspecified (authored profile absent) | [TestAssociatedTypeSyntaxWithoutResolution](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse_test.go#L57) (external) |
| associated-type-projections | malformed-and-bounded-projections | located-error-without-output | unspecified (authored profile absent) | [TestMalformedAssociatedTypes](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse_test.go#L80) (external) |
| associated-type-projections | associated-bound-syntax | unsupported-without-output | unspecified (authored profile absent) | [TestAssociatedTypeBoundsRemainUnsupported](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse_test.go#L106) (external) |

### Declared target [github.com/Bitspark/bittype/provider/rust/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L19)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L118) (external) |
| source-boundary | unsupported-and-malformed-fixtures | located-error-without-output | unspecified (authored profile absent) | [TestUnsupportedAndMalformedSources](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L140) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjU5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ydXN0L2dvI1Byb3ZpZGVyLlByaW50MTYyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcHJpbnQvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvLkxhbmd1YWdlXSwgc3RyaW5nXTUzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ydXN0L2dvLlByb3ZpZGVyNTE6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvI0xhbmd1YWdlNDpydXN0"></a>

## Rust / github.com/Bitspark/bittype/provider/rust/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/provider/rust/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L26). Receiver: github.com/Bitspark/bittype/provider/rust/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/rust/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L26) → [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L27)); receiver ; type arguments github.com/Bitspark/bittype/syntax/rust/go.Language.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/provider/rust/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L26)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L104) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L118) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjY0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvaG9yb3MvZ28jUHJvdmlkZXIuQ2hlY2sxMDk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9jaGVjay9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXTU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvaG9yb3MvZ28uUHJvdmlkZXI1MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvI0xhbmd1YWdlNTpob3Jvcw"></a>

## Horos / github.com/Bitspark/bittype/check/syntax/horos/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/check/syntax/horos/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/horos/go/check.go#L18). Receiver: github.com/Bitspark/bittype/check/syntax/horos/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/check/syntax/horos/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/horos/go/check.go#L18) → [github.com/Bitspark/bittype/check/syntax/horos/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/horos/go/check.go#L13) ([call](https://github.com/Bitspark/bittype/blob/main/check/syntax/horos/go/check.go#L18)); receiver ; type arguments .

Traversal boundary: callee is dynamic or outside implementation source.

### Declared target [github.com/Bitspark/bittype/check/syntax/horos/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/horos/go/check.go#L13)

Related declaration; its cases and observations do not qualify the public wrapper.

Definition: [github.com/Bitspark/bittype/internal/declsyntax#rules](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L66); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### program ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L67))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key program.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| module | 1 | unbounded | module |

No feature-key case association at this exact authored profile.

#### module ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L68))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key module.

Roles: 2 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| import | 0 | unbounded | import |
| declaration | 0 | unbounded | type-definition, kind-definition |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| pin | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### import ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L69))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key import.

Roles: 0 declared. Attributes: 3 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| alias | true | opaque: Predicate present; behavior not interpreted |
| name | true | opaque: Predicate present; behavior not interpreted |
| pin | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-definition ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L70))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-definition.

Roles: 3 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 0 | unbounded | type-parameter |
| result | 1 | 1 | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |
| body | 1 | 1 | descriptor, parameter, reference, application |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-definition ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L71))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-definition.

Roles: 3 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 0 | unbounded | kind-parameter |
| result | 1 | 1 | proper-sort, arrow-sort |
| body | 1 | 1 | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-parameter ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L72))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-parameter.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| contract | 1 | 1 | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-parameter ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L73))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-parameter.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| sort | 1 | 1 | proper-sort, arrow-sort |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### parameter ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L74))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### reference ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L75))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key reference.

Roles: 0 declared. Attributes: 3 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| scope | true | local, imported |
| name | true | opaque: Predicate present; behavior not interpreted |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### application ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L76))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key application.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| callee | 1 | 1 | descriptor, parameter, reference, application |
| argument | 1 | unbounded | descriptor, parameter, reference, application |

No feature-key case association at this exact authored profile.

#### descriptor ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L77))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key descriptor.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| operand | 0 | unbounded | descriptor, parameter, reference, application, atom, tuple |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| constructor | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### atom ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L78))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key atom.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| hex | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### tuple ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L79))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 0 | unbounded | descriptor, parameter, reference, application, atom, tuple |

No feature-key case association at this exact authored profile.

#### proper-sort ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L80))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key proper-sort.

Roles: 0 declared. Attributes: 0 declared. References: none.

No feature-key case association at this exact authored profile.

#### arrow-sort ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L81))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key arrow-sort.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| input | 1 | unbounded | proper-sort, arrow-sort |
| result | 1 | 1 | proper-sort, arrow-sort |

No feature-key case association at this exact authored profile.

#### kind-top ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L82))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-top.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| sort | 1 | 1 | proper-sort, arrow-sort |

No feature-key case association at this exact authored profile.

#### kind-bottom ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L83))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-bottom.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| sort | 1 | 1 | proper-sort, arrow-sort |

No feature-key case association at this exact authored profile.

#### kind-parameter-reference ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L84))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-parameter-reference.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-reference ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L85))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-reference.

Roles: 0 declared. Attributes: 3 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| scope | true | local, imported |
| name | true | opaque: Predicate present; behavior not interpreted |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-application ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L86))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-application.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| callee | 1 | 1 | kind-reference |
| argument | 1 | unbounded | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### kind-union ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L87))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-union.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| sort | 1 | 1 | proper-sort, arrow-sort |
| operand | 2 | unbounded | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### kind-intersection ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L88))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-intersection.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| sort | 1 | 1 | proper-sort, arrow-sort |
| operand | 2 | unbounded | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### kind-arrow ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L89))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-arrow.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| input | 1 | unbounded | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |
| result | 1 | 1 | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### kind-node ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L90))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-node.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| selector | 0 | unbounded | selector |
| pattern | 1 | 1 | sequence-epsilon, sequence-item, sequence-concat, sequence-alternative, sequence-star |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| constructor | true | any, atom, atom-exact, atom-len, atom-in, text-len, key-grammar, text-non-blank, lowercase-uuid, scalar, exact, wire, seq, list-of, set-of, not, map-of, tuple, compound, union, and, record |

No feature-key case association at this exact authored profile.

#### selector ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L91))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key selector.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | selector, mode |
| value | true | int, utf8-text, bool, open, closed |

No feature-key case association at this exact authored profile.

#### kind-recursive ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L92))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-recursive.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| equation | 1 | unbounded | kind-equation |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-equation ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L93))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-equation.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| body | 1 | 1 | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | false | opaque: Predicate present; behavior not interpreted |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-recursive-variable ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L94))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-recursive-variable.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| depth | true | opaque: Predicate present; behavior not interpreted |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### sequence-epsilon ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L95))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-epsilon.

Roles: 0 declared. Attributes: 0 declared. References: none.

No feature-key case association at this exact authored profile.

#### sequence-item ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L96))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-item.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### sequence-concat ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L97))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-concat.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 2 | unbounded | sequence-epsilon, sequence-item, sequence-concat, sequence-alternative, sequence-star |

No feature-key case association at this exact authored profile.

#### sequence-alternative ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L98))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-alternative.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 2 | unbounded | sequence-epsilon, sequence-item, sequence-concat, sequence-alternative, sequence-star |

No feature-key case association at this exact authored profile.

#### sequence-star ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L99))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-star.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | sequence-epsilon, sequence-item, sequence-concat, sequence-alternative, sequence-star |

No feature-key case association at this exact authored profile.

Authored scope: [docs/horos.md#syntax-vocabulary](https://github.com/Bitspark/bittype/blob/main/docs/horos.md#syntax-vocabulary); profile: unspecified (authored profile absent).

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjY0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L2NwcC9nbyNQcm92aWRlci5Gb3JtYXQyNTM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9mb3JtYXQvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2FzdC9nby5WaWV3W2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvY3BwL2dvLkxhbmd1YWdlXSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28uTGFuZ3VhZ2VdXTU3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L2NwcC9nby5Qcm92aWRlcjUwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvY3BwL2dvI0xhbmd1YWdlMzpjcHA"></a>

## C++ / github.com/Bitspark/bittype/format/syntax/cpp/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/format/syntax/cpp/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/cpp/go/format.go#L316). Receiver: github.com/Bitspark/bittype/format/syntax/cpp/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/format/syntax/cpp/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/cpp/go/format.go#L316) → [github.com/Bitspark/bittype/format/syntax/cpp/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/cpp/go/format.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/format/syntax/cpp/go/format.go#L317)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/cpp/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/cpp/go/format.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjY0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvaG9yb3MvZ28jUHJvdmlkZXIuUGFyc2UyMDk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wYXJzZS9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3BhcnNlL2dvLklucHV0LCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXV01ODpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcGFyc2Uvc3ludGF4L2hvcm9zL2dvLlByb3ZpZGVyNTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nbyNMYW5ndWFnZTU6aG9yb3M"></a>

## Horos / github.com/Bitspark/bittype/parse/syntax/horos/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/parse/syntax/horos/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse.go#L566). Receiver: github.com/Bitspark/bittype/parse/syntax/horos/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/parse/syntax/horos/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse.go#L566) → [github.com/Bitspark/bittype/parse/syntax/horos/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse.go#L61) ([call](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse.go#L567)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/horos/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse.go#L61)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/horos-syntax.md#pipeline-and-limits](https://github.com/Bitspark/bittype/blob/main/docs/horos-syntax.md#pipeline-and-limits); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjY1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvZ29sYW5nL2dvI1Byb3ZpZGVyLkNoZWNrMTEwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvY2hlY2svZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28uTGFuZ3VhZ2VdNTk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NoZWNrL3N5bnRheC9nb2xhbmcvZ28uUHJvdmlkZXI1MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nbyNMYW5ndWFnZTI6Z28"></a>

## Go / github.com/Bitspark/bittype/check/syntax/golang/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/check/syntax/golang/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L87). Receiver: github.com/Bitspark/bittype/check/syntax/golang/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/check/syntax/golang/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L87) → [github.com/Bitspark/bittype/check/syntax/golang/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L45) ([call](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L87)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/check/syntax/golang/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L45)

Related declaration; its cases and observations do not qualify the public wrapper.

Definition: [github.com/Bitspark/bittype/check/syntax/golang/go#rules](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L17); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### file ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L18))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key file.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| import | 0 | unbounded | import |
| declaration | 0 | unbounded | type-declaration, type-alias |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| package | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### import ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L19))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key import.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| path | true | opaque: Predicate present; behavior not interpreted |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-declaration ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L20))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-declaration.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 0 | unbounded | type-parameter |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-alias ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L21))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-alias.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 0 | unbounded | type-parameter |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L22))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-parameter.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| constraint | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type, type-union, type-approximation |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### named-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L23))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key named-type.

Roles: 0 declared. Attributes: 1 declared. References: target.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### qualified-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L24))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key qualified-type.

Roles: 0 declared. Attributes: 2 declared. References: target.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| qualifier | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-application ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L25))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-application.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| constructor | 1 | 1 | named-type, qualified-type |
| argument | 1 | unbounded | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

#### struct-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L26))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key struct-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| field | 0 | unbounded | field, embedded-field |

No feature-key case association at this exact authored profile.

#### field ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L27))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key field.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| tag | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### embedded-field ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L28))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key embedded-field.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| tag | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### interface-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L29))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key interface-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| member | 0 | unbounded | method, embedded-type |

No feature-key case association at this exact authored profile.

#### method ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L30))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key method.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| signature | 1 | 1 | function-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### embedded-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L31))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key embedded-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type, type-union, type-approximation |

No feature-key case association at this exact authored profile.

#### function-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L32))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key function-type.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 0 | unbounded | parameter |
| result | 0 | unbounded | parameter |

No feature-key case association at this exact authored profile.

#### parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L33))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | false | opaque: Predicate present; behavior not interpreted |
| variadic | false | true, false |

No feature-key case association at this exact authored profile.

#### pointer-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L34))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key pointer-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

#### slice-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L35))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key slice-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

#### array-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L36))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key array-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| length | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### map-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L37))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key map-type.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| key | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

#### channel-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L38))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key channel-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| direction | false | both, send, receive |

No feature-key case association at this exact authored profile.

#### type-union ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L39))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-union.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| term | 2 | unbounded | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type, type-approximation |

No feature-key case association at this exact authored profile.

#### type-approximation ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L40))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-approximation.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

Authored scope: [docs/syntax.md#go](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#go); profile: unspecified (authored profile absent).

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjY1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvZmlsZXMvaG9yb3MvZ28jUHJvdmlkZXIuRm9ybWF0MjcyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvZm9ybWF0L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvZm9ybWF0L2dvLkZpbGVzUmVxdWVzdFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3NvdXJjZS9nby5Tb3VyY2VbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nby5MYW5ndWFnZV1dNTg6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2Zvcm1hdC9maWxlcy9ob3Jvcy9nby5Qcm92aWRlcjUyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28jTGFuZ3VhZ2U1Omhvcm9z"></a>

## Horos / github.com/Bitspark/bittype/format/files/horos/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/contracts/format/go.FilesRequest&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;, github.com/Bitspark/bittype/source/go.Source&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/format/files/horos/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/files/horos/go/format.go#L54). Receiver: github.com/Bitspark/bittype/format/files/horos/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

### Declared target [github.com/Bitspark/bittype/format/files/horos/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/files/horos/go/format.go#L54)

Authored scope: [docs/source.md#horos-composition-and-pagination](https://github.com/Bitspark/bittype/blob/main/docs/source.md#horos-composition-and-pagination); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| multi-file-pagination | receipt-preserving-round-trips | fresh-output-and-receipt | unspecified (authored profile absent) | [TestFilesPaginationAndReceiptPreservingRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/files_test.go#L56) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjY1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L3J1c3QvZ28jUHJvdmlkZXIuRm9ybWF0MjU1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvZm9ybWF0L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9hc3QvZ28uVmlld1tnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28uTGFuZ3VhZ2VdLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28uTGFuZ3VhZ2VdXTU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L3J1c3QvZ28uUHJvdmlkZXI1MTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28jTGFuZ3VhZ2U0OnJ1c3Q"></a>

## Rust / github.com/Bitspark/bittype/format/syntax/rust/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/format/syntax/rust/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/rust/go/format.go#L339). Receiver: github.com/Bitspark/bittype/format/syntax/rust/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/format/syntax/rust/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/rust/go/format.go#L339) → [github.com/Bitspark/bittype/format/syntax/rust/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/rust/go/format.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/format/syntax/rust/go/format.go#L340)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/rust/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/rust/go/format.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md#rust](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#rust); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjY1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvZ29sYW5nL2dvI1Byb3ZpZGVyLlBhcnNlMjEwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcGFyc2UvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wYXJzZS9nby5JbnB1dCwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28uTGFuZ3VhZ2VdXTU5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvZ29sYW5nL2dvLlByb3ZpZGVyNTM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28jTGFuZ3VhZ2UyOmdv"></a>

## Go / github.com/Bitspark/bittype/parse/syntax/golang/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/parse/syntax/golang/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/golang/go/parse.go#L280). Receiver: github.com/Bitspark/bittype/parse/syntax/golang/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/parse/syntax/golang/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/golang/go/parse.go#L280) → [github.com/Bitspark/bittype/parse/syntax/golang/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/golang/go/parse.go#L26) ([call](https://github.com/Bitspark/bittype/blob/main/parse/syntax/golang/go/parse.go#L281)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/golang/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/golang/go/parse.go#L26)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md#go](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#go); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| go-concrete-pipeline | pagination-and-interface | round-trip-preserved | unspecified (authored profile absent) | [TestGoConcretePipeline](https://github.com/Bitspark/bittype/blob/main/integration/go/text_test.go#L15) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjY1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci90eXBlc2NyaXB0L2dvI1Byb3ZpZGVyLkNoZWNrMTE0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvY2hlY2svZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvLkxhbmd1YWdlXTU5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci90eXBlc2NyaXB0L2dvLlByb3ZpZGVyNTc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvI0xhbmd1YWdlMTA6dHlwZXNjcmlwdA"></a>

## TypeScript / github.com/Bitspark/bittype/provider/typescript/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/provider/typescript/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L22). Receiver: github.com/Bitspark/bittype/provider/typescript/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/typescript/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L22) → [github.com/Bitspark/bittype/check/syntax/typescript/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L52) ([call](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L22)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/check/syntax/typescript/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L52)

Related declaration; its cases and observations do not qualify the public wrapper.

Definition: [github.com/Bitspark/bittype/check/syntax/typescript/go#rules](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L18); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### source-file ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L19))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key source-file.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| import | 0 | unbounded | import |
| declaration | 0 | unbounded | type-alias, interface |

No feature-key case association at this exact authored profile.

#### import ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L20))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key import.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| binding | 1 | unbounded | import-binding, namespace-import |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| module | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### import-binding ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L21))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key import-binding.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### namespace-import ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L22))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key namespace-import.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-alias ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L23))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-alias.

Roles: 2 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 0 | unbounded | type-parameter |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| exported | false | true, false |

No feature-key case association at this exact authored profile.

#### interface ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L24))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key interface.

Roles: 3 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 0 | unbounded | type-parameter |
| base | 0 | unbounded | type-reference, qualified-type, type-application |
| member | 0 | unbounded | property, method, call-signature, index-signature |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| exported | false | true, false |

No feature-key case association at this exact authored profile.

#### type-parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L25))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-parameter.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| constraint | 0 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| default | 0 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-reference ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L26))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-reference.

Roles: 0 declared. Attributes: 1 declared. References: target.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### qualified-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L27))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key qualified-type.

Roles: 1 declared. Attributes: 1 declared. References: target.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| qualifier | 1 | 1 | type-reference, qualified-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-application ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L28))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-application.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| constructor | 1 | 1 | type-reference, qualified-type |
| argument | 1 | unbounded | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### keyword-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L29))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key keyword-type.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | any, unknown, never, void, undefined, boolean, number, bigint, string, symbol, object, this |

No feature-key case association at this exact authored profile.

#### literal-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L30))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key literal-type.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| literal-kind | true | string, number, bigint, boolean, null |
| text | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### array-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L31))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key array-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### tuple-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L32))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 0 | unbounded | tuple-element |

No feature-key case association at this exact authored profile.

#### tuple-element ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L33))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-element.

Roles: 1 declared. Attributes: 3 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | false | opaque: Predicate present; behavior not interpreted |
| optional | false | true, false |
| rest | false | true, false |

No feature-key case association at this exact authored profile.

#### object-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L34))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key object-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| member | 0 | unbounded | property, method, call-signature, index-signature |

No feature-key case association at this exact authored profile.

#### union-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L35))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key union-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 2 | unbounded | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### intersection-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L36))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key intersection-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 2 | unbounded | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### conditional-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L37))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key conditional-type.

Roles: 4 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| check | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| constraint | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| true | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| false | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### infer-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L38))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key infer-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 1 | 1 | type-parameter |

No feature-key case association at this exact authored profile.

#### mapped-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L39))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key mapped-type.

Roles: 3 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 1 | 1 | type-parameter |
| name-type | 0 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| readonly-modifier | false | add, remove |
| optional-modifier | false | add, remove |

No feature-key case association at this exact authored profile.

#### function-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L40))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key function-type.

Roles: 3 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 0 | unbounded | type-parameter |
| parameter | 0 | unbounded | parameter |
| result | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L41))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

Roles: 1 declared. Attributes: 3 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| optional | false | true, false |
| rest | false | true, false |

No feature-key case association at this exact authored profile.

#### property ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L42))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key property.

Roles: 1 declared. Attributes: 3 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| optional | false | true, false |
| readonly | false | true, false |

No feature-key case association at this exact authored profile.

#### method ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L43))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key method.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| signature | 1 | 1 | function-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| optional | false | true, false |

No feature-key case association at this exact authored profile.

#### call-signature ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L44))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key call-signature.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| signature | 1 | 1 | function-type |

No feature-key case association at this exact authored profile.

#### index-signature ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L45))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key index-signature.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 1 | 1 | parameter |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| readonly | false | true, false |

No feature-key case association at this exact authored profile.

#### type-operator ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L46))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-operator.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| operator | true | keyof, readonly, unique |

No feature-key case association at this exact authored profile.

#### indexed-access-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L47))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key indexed-access-type.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| object | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| index | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

Authored scope: [docs/syntax.md#typescript](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#typescript); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| conditional-mapped-syntax | independent-ast | accept | unspecified (authored profile absent) | [TestConditionalMappedAdmission](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/types_test.go#L33) (external) |
| conditional-mapped-syntax | invalid-shapes-and-contexts | located-rejection | unspecified (authored profile absent) | [TestInvalidConditionalMappedShapes](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/types_test.go#L41) (external) |

### Declared target [github.com/Bitspark/bittype/provider/typescript/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L22)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| ast-admission | malformed-carriers | reject | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L191) (external) |

<a id="cap-ODpwcm92aWRlcjY1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci90eXBlc2NyaXB0L2dvI1Byb3ZpZGVyLlBhcnNlMjE0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcGFyc2UvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wYXJzZS9nby5JbnB1dCwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvLkxhbmd1YWdlXV01OTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJvdmlkZXIvdHlwZXNjcmlwdC9nby5Qcm92aWRlcjU3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nbyNMYW5ndWFnZTEwOnR5cGVzY3JpcHQ"></a>

## TypeScript / github.com/Bitspark/bittype/provider/typescript/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/provider/typescript/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L19). Receiver: github.com/Bitspark/bittype/provider/typescript/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/typescript/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L19) → [github.com/Bitspark/bittype/parse/syntax/typescript/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/go/parse.go#L42) ([call](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L20)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/typescript/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/go/parse.go#L42)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md#typescript](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#typescript); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| conditional-mapped-syntax | structured-binders-and-modifiers | structure-preserved | unspecified (authored profile absent) | [TestConditionalMappedStructure](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/go/types_test.go#L17) (external) |
| conditional-mapped-syntax | unresolved-syntax | accept | unspecified (authored profile absent) | [TestConditionalMappedWithoutResolution](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/go/types_test.go#L48) (external) |
| conditional-mapped-syntax | malformed-and-bounded-input | located-error-without-output | unspecified (authored profile absent) | [TestMalformedConditionalMappedTypes](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/go/types_test.go#L71) (external) |

### Declared target [github.com/Bitspark/bittype/provider/typescript/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L19)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L118) (external) |
| source-boundary | unsupported-and-malformed-fixtures | located-error-without-output | unspecified (authored profile absent) | [TestUnsupportedAndMalformedSources](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L140) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjY1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci90eXBlc2NyaXB0L2dvI1Byb3ZpZGVyLlByaW50MTY4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcHJpbnQvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvLkxhbmd1YWdlXSwgc3RyaW5nXTU5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci90eXBlc2NyaXB0L2dvLlByb3ZpZGVyNTc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvI0xhbmd1YWdlMTA6dHlwZXNjcmlwdA"></a>

## TypeScript / github.com/Bitspark/bittype/provider/typescript/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/provider/typescript/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L26). Receiver: github.com/Bitspark/bittype/provider/typescript/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/typescript/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L26) → [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L27)); receiver ; type arguments github.com/Bitspark/bittype/syntax/typescript/go.Language.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/provider/typescript/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L26)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L104) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L118) (external) |
| concrete-admission | malformed-concrete-tree | rejection-without-output | unspecified (authored profile absent) | [TestMalformedConcreteTree](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L223) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjY2OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9hbmFseXplL2ZpbGVzL3J1c3QvZ28jUHJvdmlkZXIuQW5hbHl6ZTExMTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL2FuYWx5c2lzL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nby5MYW5ndWFnZV01ODpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvYW5hbHl6ZS9maWxlcy9ydXN0L2dvLlByb3ZpZGVyNTE6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvI0xhbmd1YWdlNDpydXN0"></a>

## Rust / github.com/Bitspark/bittype/analyze/files/rust/go#Provider.Analyze

Contract: github.com/Bitspark/bittype/contracts/analysis/go.Provider&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/analyze/files/rust/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/go/analyze.go#L57). Receiver: github.com/Bitspark/bittype/analyze/files/rust/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/analyze/files/rust/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/go/analyze.go#L57) → [github.com/Bitspark/bittype/analyze/files/rust/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/go/analyze.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/go/analyze.go#L58)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/analyze/files/rust/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/go/analyze.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/native-semantics.md#first-rust-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-rust-profile); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/analyze/files/rust/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/go/analyze.go#L57)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | resolution-graph-preserved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| rust-semantic-subset | native-types-and-generics | native-identity-and-substitution-preserved | unspecified (authored profile absent) | [TestRustSemanticNativeTypes](https://github.com/Bitspark/bittype/blob/main/integration/go/rust_semantics_test.go#L197) (external) |
| rust-semantic-subset | bounds-and-profile-outcomes | qualified-reports-without-partial-model | unspecified (authored profile absent) | [TestRustSemanticBoundsAndOutcomes](https://github.com/Bitspark/bittype/blob/main/integration/go/rust_semantics_test.go#L259) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjY2OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L2hvcm9zL2dvI1Byb3ZpZGVyLkZvcm1hdDI1NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL2Zvcm1hdC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvYXN0L2dvLlZpZXdbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nby5MYW5ndWFnZV0sIGdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9jc3QvZ28uQ1NUW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdXTU5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L2hvcm9zL2dvLlByb3ZpZGVyNTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nbyNMYW5ndWFnZTU6aG9yb3M"></a>

## Horos / github.com/Bitspark/bittype/format/syntax/horos/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/format/syntax/horos/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/go/format.go#L467). Receiver: github.com/Bitspark/bittype/format/syntax/horos/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/format/syntax/horos/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/go/format.go#L467) → [github.com/Bitspark/bittype/format/syntax/horos/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/go/format.go#L42) ([call](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/go/format.go#L468)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/horos/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/go/format.go#L42)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/horos-syntax.md#pipeline-and-limits](https://github.com/Bitspark/bittype/blob/main/docs/horos-syntax.md#pipeline-and-limits); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| ast-admission | unspellable-and-nil-ast | rejection-without-output | unspecified (authored profile absent) | [TestFormatterRejectsUnspellableASTWithoutPartialOutput](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/go/format_test.go#L163) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjY2OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci90eXBlc2NyaXB0L2dvI1Byb3ZpZGVyLkZvcm1hdDI2NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL2Zvcm1hdC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvYXN0L2dvLlZpZXdbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvLkxhbmd1YWdlXSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvLkxhbmd1YWdlXV01OTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJvdmlkZXIvdHlwZXNjcmlwdC9nby5Qcm92aWRlcjU3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nbyNMYW5ndWFnZTEwOnR5cGVzY3JpcHQ"></a>

## TypeScript / github.com/Bitspark/bittype/provider/typescript/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/provider/typescript/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L23). Receiver: github.com/Bitspark/bittype/provider/typescript/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/typescript/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L23) → [github.com/Bitspark/bittype/format/syntax/typescript/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/typescript/go/format.go#L24) ([call](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L24)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/typescript/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/typescript/go/format.go#L24)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md#typescript](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#typescript); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/provider/typescript/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/syntax.go#L23)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L104) (external) |
| snapshot-ownership | input-mutation-after-format | output-unchanged | unspecified (authored profile absent) | [TestFormatterSnapshotOwnership](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L121) (external) |
| ast-admission | malformed-carriers | rejection-without-output | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L191) (external) |
| native-compiler-qualification | formatted-fixture-contract | compiler-fixture-contract | unspecified (authored profile absent) | [TestNativeCompilerFixtures](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L263) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L118) (external) |
| conditional-mapped-syntax | independent-ast-formatting | ast-preserved-and-golden-text | unspecified (authored profile absent) | [TestIndependentTypeScriptMappedFormatter](https://github.com/Bitspark/bittype/blob/main/integration/go/typescript_conditional_mapped_test.go#L37) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjY2OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9yZXNvbHZlL2ZpbGVzL3J1c3QvZ28jUHJvdmlkZXIuUmVzb2x2ZTExMzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3Jlc29sdXRpb24vZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvLkxhbmd1YWdlXTU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9yZXNvbHZlL2ZpbGVzL3J1c3QvZ28uUHJvdmlkZXI1MTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28jTGFuZ3VhZ2U0OnJ1c3Q"></a>

## Rust / github.com/Bitspark/bittype/resolve/files/rust/go#Provider.Resolve

Contract: github.com/Bitspark/bittype/contracts/resolution/go.Provider&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/resolve/files/rust/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/rust/go/resolver.go#L43). Receiver: github.com/Bitspark/bittype/resolve/files/rust/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

### Declared target [github.com/Bitspark/bittype/resolve/files/rust/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/rust/go/resolver.go#L43)

Authored scope: [docs/native-semantics.md#first-rust-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-rust-profile); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | captured-graph-resolved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| rust-semantic-subset | native-types-and-generics | native-identity-and-substitution-preserved | unspecified (authored profile absent) | [TestRustSemanticNativeTypes](https://github.com/Bitspark/bittype/blob/main/integration/go/rust_semantics_test.go#L197) (external) |
| rust-semantic-subset | bounds-and-profile-outcomes | qualified-reports-without-partial-model | unspecified (authored profile absent) | [TestRustSemanticBoundsAndOutcomes](https://github.com/Bitspark/bittype/blob/main/integration/go/rust_semantics_test.go#L259) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjY3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9hbmFseXplL2ZpbGVzL2hvcm9zL2dvI1Byb3ZpZGVyLkFuYWx5emUxMTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9hbmFseXNpcy9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXTU5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9hbmFseXplL2ZpbGVzL2hvcm9zL2dvLlByb3ZpZGVyNTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nbyNMYW5ndWFnZTU6aG9yb3M"></a>

## Horos / github.com/Bitspark/bittype/analyze/files/horos/go#Provider.Analyze

Contract: github.com/Bitspark/bittype/contracts/analysis/go.Provider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/analyze/files/horos/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/horos/go/analyze.go#L75). Receiver: github.com/Bitspark/bittype/analyze/files/horos/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/analyze/files/horos/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/horos/go/analyze.go#L75) → [github.com/Bitspark/bittype/analyze/files/horos/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/horos/go/analyze.go#L37) ([call](https://github.com/Bitspark/bittype/blob/main/analyze/files/horos/go/analyze.go#L76)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/analyze/files/horos/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/horos/go/analyze.go#L37)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/analysis.md#horos-profile](https://github.com/Bitspark/bittype/blob/main/docs/analysis.md#horos-profile); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjY3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L2dvbGFuZy9nbyNQcm92aWRlci5Gb3JtYXQyNTk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9mb3JtYXQvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2FzdC9nby5WaWV3W2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvZ29sYW5nL2dvLkxhbmd1YWdlXSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28uTGFuZ3VhZ2VdXTYwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L2dvbGFuZy9nby5Qcm92aWRlcjUzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvZ29sYW5nL2dvI0xhbmd1YWdlMjpnbw"></a>

## Go / github.com/Bitspark/bittype/format/syntax/golang/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/format/syntax/golang/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/golang/go/format.go#L177). Receiver: github.com/Bitspark/bittype/format/syntax/golang/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/format/syntax/golang/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/golang/go/format.go#L177) → [github.com/Bitspark/bittype/format/syntax/golang/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/golang/go/format.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/format/syntax/golang/go/format.go#L178)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/golang/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/golang/go/format.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md#go](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#go); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| go-concrete-pipeline | pagination-and-interface | stable-formatting | unspecified (authored profile absent) | [TestGoConcretePipeline](https://github.com/Bitspark/bittype/blob/main/integration/go/text_test.go#L15) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjY3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci90eXBlc2NyaXB0L2dvI1Byb3ZpZGVyLkFuYWx5emUxMTc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9hbmFseXNpcy9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28uTGFuZ3VhZ2VdNTk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3Byb3ZpZGVyL3R5cGVzY3JpcHQvZ28uUHJvdmlkZXI1NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28jTGFuZ3VhZ2UxMDp0eXBlc2NyaXB0"></a>

## TypeScript / github.com/Bitspark/bittype/provider/typescript/go#Provider.Analyze

Contract: github.com/Bitspark/bittype/contracts/analysis/go.Provider&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/provider/typescript/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/files.go#L33). Receiver: github.com/Bitspark/bittype/provider/typescript/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/typescript/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/files.go#L33) → [github.com/Bitspark/bittype/analyze/files/typescript/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/go/analyze.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/files.go#L34)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/analyze/files/typescript/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/go/analyze.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/native-semantics.md#first-typescript-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-typescript-profile); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjY3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci90eXBlc2NyaXB0L2dvI1Byb3ZpZGVyLlJlc29sdmUxMTk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9yZXNvbHV0aW9uL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nby5MYW5ndWFnZV01OTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJvdmlkZXIvdHlwZXNjcmlwdC9nby5Qcm92aWRlcjU3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nbyNMYW5ndWFnZTEwOnR5cGVzY3JpcHQ"></a>

## TypeScript / github.com/Bitspark/bittype/provider/typescript/go#Provider.Resolve

Contract: github.com/Bitspark/bittype/contracts/resolution/go.Provider&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/provider/typescript/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/files.go#L30). Receiver: github.com/Bitspark/bittype/provider/typescript/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/typescript/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/files.go#L30) → [github.com/Bitspark/bittype/resolve/files/typescript/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/typescript/go/resolver.go#L43) ([call](https://github.com/Bitspark/bittype/blob/main/provider/typescript/go/files.go#L31)); receiver github.com/Bitspark/bittype/resolve/files/typescript/go.Provider; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/resolve/files/typescript/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/typescript/go/resolver.go#L43)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/native-semantics.md#first-typescript-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-typescript-profile); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | captured-graph-resolved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| typescript-semantic-subset | generic-and-native-facts | native-identity-and-substitution-preserved | unspecified (authored profile absent) | [TestTypeScriptSemanticGenericsAndNativeFacts](https://github.com/Bitspark/bittype/blob/main/integration/go/typescript_semantics_test.go#L205) (external) |
| typescript-semantic-subset | rejected-and-unsupported-fixtures | no-partial-model | unspecified (authored profile absent) | [TestTypeScriptSemanticReports](https://github.com/Bitspark/bittype/blob/main/integration/go/typescript_semantics_test.go#L249) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjY3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9yZXNvbHZlL2ZpbGVzL2hvcm9zL2dvI1Byb3ZpZGVyLlJlc29sdmUxMTQ6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9yZXNvbHV0aW9uL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdNTk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3Jlc29sdmUvZmlsZXMvaG9yb3MvZ28uUHJvdmlkZXI1MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvI0xhbmd1YWdlNTpob3Jvcw"></a>

## Horos / github.com/Bitspark/bittype/resolve/files/horos/go#Provider.Resolve

Contract: github.com/Bitspark/bittype/contracts/resolution/go.Provider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/resolve/files/horos/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/horos/go/resolver.go#L57). Receiver: github.com/Bitspark/bittype/resolve/files/horos/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/resolve/files/horos/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/horos/go/resolver.go#L57) → [github.com/Bitspark/bittype/resolve/files/horos/go#ResolveReferences](https://github.com/Bitspark/bittype/blob/main/resolve/files/horos/go/resolver.go#L50) ([call](https://github.com/Bitspark/bittype/blob/main/resolve/files/horos/go/resolver.go#L58)); receiver ; type arguments .

Traversal boundary: callee is dynamic or outside implementation source.

### Declared target [github.com/Bitspark/bittype/resolve/files/horos/go#ResolveReferences](https://github.com/Bitspark/bittype/blob/main/resolve/files/horos/go/resolver.go#L50)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/resolution.md#horos-composition](https://github.com/Bitspark/bittype/blob/main/docs/resolution.md#horos-composition); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjY4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9hbmFseXplL2ZpbGVzL2dvbGFuZy9nbyNQcm92aWRlci5BbmFseXplMTEzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvYW5hbHlzaXMvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28uTGFuZ3VhZ2VdNjA6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2FuYWx5emUvZmlsZXMvZ29sYW5nL2dvLlByb3ZpZGVyNTM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28jTGFuZ3VhZ2UyOmdv"></a>

## Go / github.com/Bitspark/bittype/analyze/files/golang/go#Provider.Analyze

Contract: github.com/Bitspark/bittype/contracts/analysis/go.Provider&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/analyze/files/golang/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/go/analyze.go#L57). Receiver: github.com/Bitspark/bittype/analyze/files/golang/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/analyze/files/golang/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/go/analyze.go#L57) → [github.com/Bitspark/bittype/analyze/files/golang/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/go/analyze.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/go/analyze.go#L58)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/analyze/files/golang/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/go/analyze.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/native-semantics.md#first-go-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-go-profile); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/analyze/files/golang/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/go/analyze.go#L57)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | resolution-graph-preserved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| go-semantic-subset | rejected-and-unsupported-fixtures | no-partial-model | unspecified (authored profile absent) | [TestGoSemanticReportsAndCoverage](https://github.com/Bitspark/bittype/blob/main/integration/go/native_semantics_test.go#L182) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjY4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9yZXNvbHZlL2ZpbGVzL2dvbGFuZy9nbyNQcm92aWRlci5SZXNvbHZlMTE1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcmVzb2x1dGlvbi9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nby5MYW5ndWFnZV02MDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcmVzb2x2ZS9maWxlcy9nb2xhbmcvZ28uUHJvdmlkZXI1MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nbyNMYW5ndWFnZTI6Z28"></a>

## Go / github.com/Bitspark/bittype/resolve/files/golang/go#Provider.Resolve

Contract: github.com/Bitspark/bittype/contracts/resolution/go.Provider&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/resolve/files/golang/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/go/resolver.go#L82). Receiver: github.com/Bitspark/bittype/resolve/files/golang/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/resolve/files/golang/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/go/resolver.go#L82) → [github.com/Bitspark/bittype/resolve/files/golang/go#Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/go/resolver.go#L25) ([call](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/go/resolver.go#L83)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/resolve/files/golang/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/go/resolver.go#L82)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | captured-graph-resolved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| go-semantic-subset | rejected-and-unsupported-fixtures | no-partial-model | unspecified (authored profile absent) | [TestGoSemanticReportsAndCoverage](https://github.com/Bitspark/bittype/blob/main/integration/go/native_semantics_test.go#L182) (external) |

### Declared target [github.com/Bitspark/bittype/resolve/files/golang/go#Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/go/resolver.go#L25)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/native-semantics.md#first-go-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-go-profile); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjY5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvdHlwZXNjcmlwdC9nbyNQcm92aWRlci5DaGVjazExNDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL2NoZWNrL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nby5MYW5ndWFnZV02MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY2hlY2svc3ludGF4L3R5cGVzY3JpcHQvZ28uUHJvdmlkZXI1NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28jTGFuZ3VhZ2UxMDp0eXBlc2NyaXB0"></a>

## TypeScript / github.com/Bitspark/bittype/check/syntax/typescript/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/check/syntax/typescript/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L177). Receiver: github.com/Bitspark/bittype/check/syntax/typescript/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/check/syntax/typescript/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L177) → [github.com/Bitspark/bittype/check/syntax/typescript/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L52) ([call](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L177)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/check/syntax/typescript/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L52)

Related declaration; its cases and observations do not qualify the public wrapper.

Definition: [github.com/Bitspark/bittype/check/syntax/typescript/go#rules](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L18); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### source-file ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L19))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key source-file.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| import | 0 | unbounded | import |
| declaration | 0 | unbounded | type-alias, interface |

No feature-key case association at this exact authored profile.

#### import ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L20))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key import.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| binding | 1 | unbounded | import-binding, namespace-import |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| module | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### import-binding ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L21))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key import-binding.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### namespace-import ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L22))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key namespace-import.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-alias ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L23))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-alias.

Roles: 2 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 0 | unbounded | type-parameter |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| exported | false | true, false |

No feature-key case association at this exact authored profile.

#### interface ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L24))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key interface.

Roles: 3 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 0 | unbounded | type-parameter |
| base | 0 | unbounded | type-reference, qualified-type, type-application |
| member | 0 | unbounded | property, method, call-signature, index-signature |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| exported | false | true, false |

No feature-key case association at this exact authored profile.

#### type-parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L25))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-parameter.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| constraint | 0 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| default | 0 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-reference ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L26))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-reference.

Roles: 0 declared. Attributes: 1 declared. References: target.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### qualified-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L27))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key qualified-type.

Roles: 1 declared. Attributes: 1 declared. References: target.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| qualifier | 1 | 1 | type-reference, qualified-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-application ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L28))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-application.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| constructor | 1 | 1 | type-reference, qualified-type |
| argument | 1 | unbounded | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### keyword-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L29))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key keyword-type.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | any, unknown, never, void, undefined, boolean, number, bigint, string, symbol, object, this |

No feature-key case association at this exact authored profile.

#### literal-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L30))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key literal-type.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| literal-kind | true | string, number, bigint, boolean, null |
| text | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### array-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L31))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key array-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### tuple-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L32))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 0 | unbounded | tuple-element |

No feature-key case association at this exact authored profile.

#### tuple-element ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L33))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-element.

Roles: 1 declared. Attributes: 3 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | false | opaque: Predicate present; behavior not interpreted |
| optional | false | true, false |
| rest | false | true, false |

No feature-key case association at this exact authored profile.

#### object-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L34))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key object-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| member | 0 | unbounded | property, method, call-signature, index-signature |

No feature-key case association at this exact authored profile.

#### union-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L35))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key union-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 2 | unbounded | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### intersection-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L36))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key intersection-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 2 | unbounded | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### conditional-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L37))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key conditional-type.

Roles: 4 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| check | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| constraint | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| true | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| false | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### infer-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L38))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key infer-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 1 | 1 | type-parameter |

No feature-key case association at this exact authored profile.

#### mapped-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L39))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key mapped-type.

Roles: 3 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 1 | 1 | type-parameter |
| name-type | 0 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| readonly-modifier | false | add, remove |
| optional-modifier | false | add, remove |

No feature-key case association at this exact authored profile.

#### function-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L40))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key function-type.

Roles: 3 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 0 | unbounded | type-parameter |
| parameter | 0 | unbounded | parameter |
| result | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L41))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

Roles: 1 declared. Attributes: 3 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| optional | false | true, false |
| rest | false | true, false |

No feature-key case association at this exact authored profile.

#### property ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L42))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key property.

Roles: 1 declared. Attributes: 3 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| optional | false | true, false |
| readonly | false | true, false |

No feature-key case association at this exact authored profile.

#### method ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L43))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key method.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| signature | 1 | 1 | function-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| optional | false | true, false |

No feature-key case association at this exact authored profile.

#### call-signature ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L44))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key call-signature.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| signature | 1 | 1 | function-type |

No feature-key case association at this exact authored profile.

#### index-signature ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L45))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key index-signature.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 1 | 1 | parameter |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| readonly | false | true, false |

No feature-key case association at this exact authored profile.

#### type-operator ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L46))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-operator.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| operator | true | keyof, readonly, unique |

No feature-key case association at this exact authored profile.

#### indexed-access-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L47))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key indexed-access-type.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| object | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| index | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

Authored scope: [docs/syntax.md#typescript](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#typescript); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| conditional-mapped-syntax | independent-ast | accept | unspecified (authored profile absent) | [TestConditionalMappedAdmission](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/types_test.go#L33) (external) |
| conditional-mapped-syntax | invalid-shapes-and-contexts | located-rejection | unspecified (authored profile absent) | [TestInvalidConditionalMappedShapes](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/types_test.go#L41) (external) |

<a id="cap-ODpwcm92aWRlcjY5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvdHlwZXNjcmlwdC9nbyNQcm92aWRlci5QYXJzZTIxNDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3BhcnNlL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcGFyc2UvZ28uSW5wdXQsIGdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9jc3QvZ28uQ1NUW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nby5MYW5ndWFnZV1dNjM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3BhcnNlL3N5bnRheC90eXBlc2NyaXB0L2dvLlByb3ZpZGVyNTc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvI0xhbmd1YWdlMTA6dHlwZXNjcmlwdA"></a>

## TypeScript / github.com/Bitspark/bittype/parse/syntax/typescript/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/parse/syntax/typescript/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/go/parse.go#L568). Receiver: github.com/Bitspark/bittype/parse/syntax/typescript/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/parse/syntax/typescript/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/go/parse.go#L568) → [github.com/Bitspark/bittype/parse/syntax/typescript/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/go/parse.go#L42) ([call](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/go/parse.go#L569)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/typescript/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/go/parse.go#L42)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md#typescript](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#typescript); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| conditional-mapped-syntax | structured-binders-and-modifiers | structure-preserved | unspecified (authored profile absent) | [TestConditionalMappedStructure](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/go/types_test.go#L17) (external) |
| conditional-mapped-syntax | unresolved-syntax | accept | unspecified (authored profile absent) | [TestConditionalMappedWithoutResolution](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/go/types_test.go#L48) (external) |
| conditional-mapped-syntax | malformed-and-bounded-input | located-error-without-output | unspecified (authored profile absent) | [TestMalformedConditionalMappedTypes](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/go/types_test.go#L71) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjYwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ob3Jvcy9nbyNQcm92aWRlci5DaGVjazEwOTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL2NoZWNrL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdNTQ6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3Byb3ZpZGVyL2hvcm9zL2dvLlByb3ZpZGVyNTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nbyNMYW5ndWFnZTU6aG9yb3M"></a>

## Horos / github.com/Bitspark/bittype/provider/horos/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/provider/horos/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L22). Receiver: github.com/Bitspark/bittype/provider/horos/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/horos/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L22) → [github.com/Bitspark/bittype/check/syntax/horos/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/horos/go/check.go#L13) ([call](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L22)); receiver ; type arguments .

Traversal boundary: callee is dynamic or outside implementation source.

### Declared target [github.com/Bitspark/bittype/check/syntax/horos/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/horos/go/check.go#L13)

Related declaration; its cases and observations do not qualify the public wrapper.

Definition: [github.com/Bitspark/bittype/internal/declsyntax#rules](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L66); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### program ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L67))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key program.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| module | 1 | unbounded | module |

No feature-key case association at this exact authored profile.

#### module ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L68))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key module.

Roles: 2 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| import | 0 | unbounded | import |
| declaration | 0 | unbounded | type-definition, kind-definition |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| pin | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### import ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L69))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key import.

Roles: 0 declared. Attributes: 3 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| alias | true | opaque: Predicate present; behavior not interpreted |
| name | true | opaque: Predicate present; behavior not interpreted |
| pin | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-definition ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L70))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-definition.

Roles: 3 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 0 | unbounded | type-parameter |
| result | 1 | 1 | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |
| body | 1 | 1 | descriptor, parameter, reference, application |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-definition ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L71))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-definition.

Roles: 3 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 0 | unbounded | kind-parameter |
| result | 1 | 1 | proper-sort, arrow-sort |
| body | 1 | 1 | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-parameter ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L72))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-parameter.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| contract | 1 | 1 | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-parameter ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L73))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-parameter.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| sort | 1 | 1 | proper-sort, arrow-sort |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### parameter ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L74))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### reference ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L75))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key reference.

Roles: 0 declared. Attributes: 3 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| scope | true | local, imported |
| name | true | opaque: Predicate present; behavior not interpreted |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### application ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L76))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key application.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| callee | 1 | 1 | descriptor, parameter, reference, application |
| argument | 1 | unbounded | descriptor, parameter, reference, application |

No feature-key case association at this exact authored profile.

#### descriptor ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L77))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key descriptor.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| operand | 0 | unbounded | descriptor, parameter, reference, application, atom, tuple |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| constructor | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### atom ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L78))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key atom.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| hex | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### tuple ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L79))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 0 | unbounded | descriptor, parameter, reference, application, atom, tuple |

No feature-key case association at this exact authored profile.

#### proper-sort ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L80))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key proper-sort.

Roles: 0 declared. Attributes: 0 declared. References: none.

No feature-key case association at this exact authored profile.

#### arrow-sort ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L81))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key arrow-sort.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| input | 1 | unbounded | proper-sort, arrow-sort |
| result | 1 | 1 | proper-sort, arrow-sort |

No feature-key case association at this exact authored profile.

#### kind-top ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L82))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-top.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| sort | 1 | 1 | proper-sort, arrow-sort |

No feature-key case association at this exact authored profile.

#### kind-bottom ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L83))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-bottom.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| sort | 1 | 1 | proper-sort, arrow-sort |

No feature-key case association at this exact authored profile.

#### kind-parameter-reference ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L84))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-parameter-reference.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-reference ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L85))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-reference.

Roles: 0 declared. Attributes: 3 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| scope | true | local, imported |
| name | true | opaque: Predicate present; behavior not interpreted |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-application ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L86))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-application.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| callee | 1 | 1 | kind-reference |
| argument | 1 | unbounded | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### kind-union ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L87))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-union.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| sort | 1 | 1 | proper-sort, arrow-sort |
| operand | 2 | unbounded | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### kind-intersection ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L88))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-intersection.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| sort | 1 | 1 | proper-sort, arrow-sort |
| operand | 2 | unbounded | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### kind-arrow ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L89))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-arrow.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| input | 1 | unbounded | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |
| result | 1 | 1 | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### kind-node ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L90))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-node.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| selector | 0 | unbounded | selector |
| pattern | 1 | 1 | sequence-epsilon, sequence-item, sequence-concat, sequence-alternative, sequence-star |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| constructor | true | any, atom, atom-exact, atom-len, atom-in, text-len, key-grammar, text-non-blank, lowercase-uuid, scalar, exact, wire, seq, list-of, set-of, not, map-of, tuple, compound, union, and, record |

No feature-key case association at this exact authored profile.

#### selector ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L91))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key selector.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | selector, mode |
| value | true | int, utf8-text, bool, open, closed |

No feature-key case association at this exact authored profile.

#### kind-recursive ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L92))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-recursive.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| equation | 1 | unbounded | kind-equation |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-equation ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L93))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-equation.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| body | 1 | 1 | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | false | opaque: Predicate present; behavior not interpreted |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-recursive-variable ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L94))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-recursive-variable.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| depth | true | opaque: Predicate present; behavior not interpreted |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### sequence-epsilon ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L95))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-epsilon.

Roles: 0 declared. Attributes: 0 declared. References: none.

No feature-key case association at this exact authored profile.

#### sequence-item ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L96))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-item.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### sequence-concat ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L97))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-concat.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 2 | unbounded | sequence-epsilon, sequence-item, sequence-concat, sequence-alternative, sequence-star |

No feature-key case association at this exact authored profile.

#### sequence-alternative ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L98))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-alternative.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 2 | unbounded | sequence-epsilon, sequence-item, sequence-concat, sequence-alternative, sequence-star |

No feature-key case association at this exact authored profile.

#### sequence-star ([source](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L99))

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-star.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | sequence-epsilon, sequence-item, sequence-concat, sequence-alternative, sequence-star |

No feature-key case association at this exact authored profile.

Authored scope: [docs/horos.md#syntax-vocabulary](https://github.com/Bitspark/bittype/blob/main/docs/horos.md#syntax-vocabulary); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/provider/horos/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L22)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| independent-engine-stages | pagination-source | independently-composable | unspecified (authored profile absent) | [TestHorosSourceThroughIndependentEngineStages](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_text_test.go#L27) (external) |
| horos-text-vocabulary | declared-constructor-fixtures | canonical-round-trip-preserved | unspecified (authored profile absent) | [TestCompleteSourceVocabulary](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse_test.go#L80) (external) |
| syntax-semantics-boundary | semantically-invalid-syntax-fixtures | syntax-accepted | unspecified (authored profile absent) | [TestSyntaxDoesNotResolveOrCheckContracts](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse_test.go#L171) (external) |

<a id="cap-ODpwcm92aWRlcjYwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ob3Jvcy9nbyNQcm92aWRlci5QYXJzZTIwOTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3BhcnNlL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcGFyc2UvZ28uSW5wdXQsIGdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9jc3QvZ28uQ1NUW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdXTU0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ob3Jvcy9nby5Qcm92aWRlcjUyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28jTGFuZ3VhZ2U1Omhvcm9z"></a>

## Horos / github.com/Bitspark/bittype/provider/horos/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/provider/horos/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L19). Receiver: github.com/Bitspark/bittype/provider/horos/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/horos/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L19) → [github.com/Bitspark/bittype/parse/syntax/horos/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse.go#L61) ([call](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L20)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/horos/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse.go#L61)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/horos-syntax.md#pipeline-and-limits](https://github.com/Bitspark/bittype/blob/main/docs/horos-syntax.md#pipeline-and-limits); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/provider/horos/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L19)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| horos-modules-and-generics | pagination-and-imported-library | canonical-identity-preserved | unspecified (authored profile absent) | [TestCanonicalPaginationAndImportedLibraryRoundTrip](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/go/format_test.go#L55) (external) |
| independent-engine-stages | pagination-source | independently-composable | unspecified (authored profile absent) | [TestHorosSourceThroughIndependentEngineStages](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_text_test.go#L27) (external) |
| horos-text-vocabulary | declared-constructor-fixtures | canonical-round-trip-preserved | unspecified (authored profile absent) | [TestCompleteSourceVocabulary](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse_test.go#L80) (external) |
| syntax-semantics-boundary | semantically-invalid-syntax-fixtures | syntax-accepted | unspecified (authored profile absent) | [TestSyntaxDoesNotResolveOrCheckContracts](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse_test.go#L171) (external) |
| source-boundary | whole-input-rejection | located-error-without-output | unspecified (authored profile absent) | [TestRejectWholeInput](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse_test.go#L192) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjYwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ob3Jvcy9nbyNQcm92aWRlci5QcmludDE2MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3ByaW50L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9jc3QvZ28uQ1NUW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdLCBzdHJpbmddNTQ6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3Byb3ZpZGVyL2hvcm9zL2dvLlByb3ZpZGVyNTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nbyNMYW5ndWFnZTU6aG9yb3M"></a>

## Horos / github.com/Bitspark/bittype/provider/horos/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/provider/horos/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L26). Receiver: github.com/Bitspark/bittype/provider/horos/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/horos/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L26) → [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L27)); receiver ; type arguments github.com/Bitspark/bittype/syntax/horos/go.Language.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/provider/horos/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L26)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| horos-modules-and-generics | pagination-and-imported-library | canonical-identity-preserved | unspecified (authored profile absent) | [TestCanonicalPaginationAndImportedLibraryRoundTrip](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/go/format_test.go#L55) (external) |
| independent-engine-stages | pagination-source | independently-composable | unspecified (authored profile absent) | [TestHorosSourceThroughIndependentEngineStages](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_text_test.go#L27) (external) |
| horos-text-vocabulary | declared-constructor-fixtures | canonical-round-trip-preserved | unspecified (authored profile absent) | [TestCompleteSourceVocabulary](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse_test.go#L80) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjYwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ydXN0L2dvI1Byb3ZpZGVyLkZvcm1hdDI1NTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL2Zvcm1hdC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvYXN0L2dvLlZpZXdbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvLkxhbmd1YWdlXSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvLkxhbmd1YWdlXV01MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJvdmlkZXIvcnVzdC9nby5Qcm92aWRlcjUxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nbyNMYW5ndWFnZTQ6cnVzdA"></a>

## Rust / github.com/Bitspark/bittype/provider/rust/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/provider/rust/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L23). Receiver: github.com/Bitspark/bittype/provider/rust/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/rust/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L23) → [github.com/Bitspark/bittype/format/syntax/rust/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/rust/go/format.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L24)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/rust/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/rust/go/format.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md#rust](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#rust); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/provider/rust/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/syntax.go#L23)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L104) (external) |
| ast-admission | malformed-carriers | rejection-without-output | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L191) (external) |
| native-compiler-qualification | formatted-fixture-contract | compiler-fixture-contract | unspecified (authored profile absent) | [TestNativeCompilerFixtures](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L263) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L118) (external) |
| associated-type-projections | independent-ast-formatting | ast-preserved-and-golden-text | unspecified (authored profile absent) | [TestIndependentRustAssociatedFormatter](https://github.com/Bitspark/bittype/blob/main/integration/go/rust_associated_types_test.go#L40) (external) |
| generic-associated-declarations | independent-ast-formatting | ast-preserved-and-golden-text | unspecified (authored profile absent) | [TestIndependentRustDeclarationFormatter](https://github.com/Bitspark/bittype/blob/main/integration/go/rust_declarations_test.go#L41) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjYxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9nb2xhbmcvZ28jUHJvdmlkZXIuQ2hlY2sxMTA6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9jaGVjay9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nby5MYW5ndWFnZV01NTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJvdmlkZXIvZ29sYW5nL2dvLlByb3ZpZGVyNTM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28jTGFuZ3VhZ2UyOmdv"></a>

## Go / github.com/Bitspark/bittype/provider/golang/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/provider/golang/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L22). Receiver: github.com/Bitspark/bittype/provider/golang/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/golang/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L22) → [github.com/Bitspark/bittype/check/syntax/golang/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L45) ([call](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L22)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/check/syntax/golang/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L45)

Related declaration; its cases and observations do not qualify the public wrapper.

Definition: [github.com/Bitspark/bittype/check/syntax/golang/go#rules](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L17); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### file ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L18))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key file.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| import | 0 | unbounded | import |
| declaration | 0 | unbounded | type-declaration, type-alias |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| package | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### import ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L19))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key import.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| path | true | opaque: Predicate present; behavior not interpreted |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-declaration ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L20))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-declaration.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 0 | unbounded | type-parameter |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-alias ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L21))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-alias.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 0 | unbounded | type-parameter |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L22))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-parameter.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| constraint | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type, type-union, type-approximation |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### named-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L23))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key named-type.

Roles: 0 declared. Attributes: 1 declared. References: target.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### qualified-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L24))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key qualified-type.

Roles: 0 declared. Attributes: 2 declared. References: target.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| qualifier | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-application ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L25))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-application.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| constructor | 1 | 1 | named-type, qualified-type |
| argument | 1 | unbounded | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

#### struct-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L26))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key struct-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| field | 0 | unbounded | field, embedded-field |

No feature-key case association at this exact authored profile.

#### field ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L27))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key field.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| tag | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### embedded-field ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L28))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key embedded-field.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| tag | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### interface-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L29))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key interface-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| member | 0 | unbounded | method, embedded-type |

No feature-key case association at this exact authored profile.

#### method ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L30))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key method.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| signature | 1 | 1 | function-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### embedded-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L31))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key embedded-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type, type-union, type-approximation |

No feature-key case association at this exact authored profile.

#### function-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L32))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key function-type.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 0 | unbounded | parameter |
| result | 0 | unbounded | parameter |

No feature-key case association at this exact authored profile.

#### parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L33))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | false | opaque: Predicate present; behavior not interpreted |
| variadic | false | true, false |

No feature-key case association at this exact authored profile.

#### pointer-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L34))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key pointer-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

#### slice-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L35))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key slice-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

#### array-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L36))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key array-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| length | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### map-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L37))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key map-type.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| key | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

#### channel-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L38))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key channel-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| direction | false | both, send, receive |

No feature-key case association at this exact authored profile.

#### type-union ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L39))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-union.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| term | 2 | unbounded | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type, type-approximation |

No feature-key case association at this exact authored profile.

#### type-approximation ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L40))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-approximation.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

Authored scope: [docs/syntax.md#go](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#go); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/provider/golang/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L22)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| ast-admission | malformed-carriers | reject | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L191) (external) |

<a id="cap-ODpwcm92aWRlcjYxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9nb2xhbmcvZ28jUHJvdmlkZXIuUGFyc2UyMTA6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wYXJzZS9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3BhcnNlL2dvLklucHV0LCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nby5MYW5ndWFnZV1dNTU6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3Byb3ZpZGVyL2dvbGFuZy9nby5Qcm92aWRlcjUzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvZ29sYW5nL2dvI0xhbmd1YWdlMjpnbw"></a>

## Go / github.com/Bitspark/bittype/provider/golang/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/provider/golang/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L19). Receiver: github.com/Bitspark/bittype/provider/golang/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/golang/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L19) → [github.com/Bitspark/bittype/parse/syntax/golang/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/golang/go/parse.go#L26) ([call](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L20)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/golang/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/golang/go/parse.go#L26)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md#go](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#go); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| go-concrete-pipeline | pagination-and-interface | round-trip-preserved | unspecified (authored profile absent) | [TestGoConcretePipeline](https://github.com/Bitspark/bittype/blob/main/integration/go/text_test.go#L15) (external) |

### Declared target [github.com/Bitspark/bittype/provider/golang/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L19)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L118) (external) |
| source-boundary | unsupported-and-malformed-fixtures | located-error-without-output | unspecified (authored profile absent) | [TestUnsupportedAndMalformedSources](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L140) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjYxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9nb2xhbmcvZ28jUHJvdmlkZXIuUHJpbnQxNjQ6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wcmludC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nby5MYW5ndWFnZV0sIHN0cmluZ101NTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJvdmlkZXIvZ29sYW5nL2dvLlByb3ZpZGVyNTM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28jTGFuZ3VhZ2UyOmdv"></a>

## Go / github.com/Bitspark/bittype/provider/golang/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/provider/golang/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L26). Receiver: github.com/Bitspark/bittype/provider/golang/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/golang/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L26) → [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L27)); receiver ; type arguments github.com/Bitspark/bittype/syntax/golang/go.Language.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| go-concrete-pipeline | pagination-and-interface | source-retained | unspecified (authored profile absent) | [TestGoConcretePipeline](https://github.com/Bitspark/bittype/blob/main/integration/go/text_test.go#L15) (external) |

### Declared target [github.com/Bitspark/bittype/provider/golang/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L26)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L104) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L118) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjYxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ob3Jvcy9nbyNQcm92aWRlci5Gb3JtYXQyNTc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9mb3JtYXQvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2FzdC9nby5WaWV3W2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXV01NDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJvdmlkZXIvaG9yb3MvZ28uUHJvdmlkZXI1MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvI0xhbmd1YWdlNTpob3Jvcw"></a>

## Horos / github.com/Bitspark/bittype/provider/horos/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/provider/horos/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L23). Receiver: github.com/Bitspark/bittype/provider/horos/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/horos/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L23) → [github.com/Bitspark/bittype/format/syntax/horos/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/go/format.go#L42) ([call](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L24)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/horos/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/go/format.go#L42)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/horos-syntax.md#pipeline-and-limits](https://github.com/Bitspark/bittype/blob/main/docs/horos-syntax.md#pipeline-and-limits); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| ast-admission | unspellable-and-nil-ast | rejection-without-output | unspecified (authored profile absent) | [TestFormatterRejectsUnspellableASTWithoutPartialOutput](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/go/format_test.go#L163) (external) |

### Declared target [github.com/Bitspark/bittype/provider/horos/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/syntax.go#L23)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| horos-modules-and-generics | pagination-and-imported-library | canonical-identity-preserved | unspecified (authored profile absent) | [TestCanonicalPaginationAndImportedLibraryRoundTrip](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/go/format_test.go#L55) (external) |
| independent-engine-stages | pagination-source | independently-composable | unspecified (authored profile absent) | [TestHorosSourceThroughIndependentEngineStages](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_text_test.go#L27) (external) |
| horos-text-vocabulary | declared-constructor-fixtures | canonical-round-trip-preserved | unspecified (authored profile absent) | [TestCompleteSourceVocabulary](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse_test.go#L80) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjYxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ydXN0L2dvI1Byb3ZpZGVyLkFuYWx5emUxMTE6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9hbmFseXNpcy9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28uTGFuZ3VhZ2VdNTM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3Byb3ZpZGVyL3J1c3QvZ28uUHJvdmlkZXI1MTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28jTGFuZ3VhZ2U0OnJ1c3Q"></a>

## Rust / github.com/Bitspark/bittype/provider/rust/go#Provider.Analyze

Contract: github.com/Bitspark/bittype/contracts/analysis/go.Provider&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/provider/rust/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/files.go#L33). Receiver: github.com/Bitspark/bittype/provider/rust/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/rust/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/files.go#L33) → [github.com/Bitspark/bittype/analyze/files/rust/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/go/analyze.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/files.go#L34)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/analyze/files/rust/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/go/analyze.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/native-semantics.md#first-rust-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-rust-profile); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjYxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ydXN0L2dvI1Byb3ZpZGVyLlJlc29sdmUxMTM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9yZXNvbHV0aW9uL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nby5MYW5ndWFnZV01MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJvdmlkZXIvcnVzdC9nby5Qcm92aWRlcjUxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nbyNMYW5ndWFnZTQ6cnVzdA"></a>

## Rust / github.com/Bitspark/bittype/provider/rust/go#Provider.Resolve

Contract: github.com/Bitspark/bittype/contracts/resolution/go.Provider&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/provider/rust/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/files.go#L30). Receiver: github.com/Bitspark/bittype/provider/rust/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/rust/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/files.go#L30) → [github.com/Bitspark/bittype/resolve/files/rust/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/rust/go/resolver.go#L43) ([call](https://github.com/Bitspark/bittype/blob/main/provider/rust/go/files.go#L31)); receiver github.com/Bitspark/bittype/resolve/files/rust/go.Provider; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/resolve/files/rust/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/rust/go/resolver.go#L43)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/native-semantics.md#first-rust-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-rust-profile); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | captured-graph-resolved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| rust-semantic-subset | native-types-and-generics | native-identity-and-substitution-preserved | unspecified (authored profile absent) | [TestRustSemanticNativeTypes](https://github.com/Bitspark/bittype/blob/main/integration/go/rust_semantics_test.go#L197) (external) |
| rust-semantic-subset | bounds-and-profile-outcomes | qualified-reports-without-partial-model | unspecified (authored profile absent) | [TestRustSemanticBoundsAndOutcomes](https://github.com/Bitspark/bittype/blob/main/integration/go/rust_semantics_test.go#L259) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjYyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvY3BwL2dvI1Byb3ZpZGVyLkNoZWNrMTA3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvY2hlY2svZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28uTGFuZ3VhZ2VdNTY6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NoZWNrL3N5bnRheC9jcHAvZ28uUHJvdmlkZXI1MDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nbyNMYW5ndWFnZTM6Y3Bw"></a>

## C++ / github.com/Bitspark/bittype/check/syntax/cpp/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/check/syntax/cpp/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L129). Receiver: github.com/Bitspark/bittype/check/syntax/cpp/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/check/syntax/cpp/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L129) → [github.com/Bitspark/bittype/check/syntax/cpp/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L57) ([call](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L129)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/check/syntax/cpp/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L57)

Related declaration; its cases and observations do not qualify the public wrapper.

Definition: [github.com/Bitspark/bittype/check/syntax/cpp/go#rules](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L29); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### translation-unit ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L30))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key translation-unit.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| include | 0 | unbounded | include |
| declaration | 0 | unbounded | namespace, record, type-alias, enum |

No feature-key case association at this exact authored profile.

#### include ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L31))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key include.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| path | true | opaque: Predicate present; behavior not interpreted |
| system | false | true, false |

No feature-key case association at this exact authored profile.

#### namespace ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L32))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key namespace.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| declaration | 0 | unbounded | namespace, record, type-alias, enum |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### record ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L33))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key record.

Roles: 3 declared. Attributes: unknown: Attribute collection could not be decoded. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| template-parameter | 0 | unbounded | type-parameter |
| base | 0 | unbounded | base |
| member | 0 | unbounded | field, method, destructor, type-alias, record, enum |

No feature-key case association at this exact authored profile.

#### type-alias ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L34))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-alias.

Roles: 2 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| template-parameter | 0 | unbounded | type-parameter |
| type | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| access | false | public, protected, private |

No feature-key case association at this exact authored profile.

#### enum ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L35))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key enum.

Roles: 2 declared. Attributes: unknown: Attribute collection could not be decoded. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| underlying | 0 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |
| enumerator | 0 | unbounded | enumerator |

No feature-key case association at this exact authored profile.

#### enumerator ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L36))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key enumerator.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| value | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### base ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L37))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key base.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| access | false | public, protected, private |
| virtual | false | true, false |

No feature-key case association at this exact authored profile.

#### field ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L38))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key field.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| access | false | public, protected, private |

No feature-key case association at this exact authored profile.

#### method ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L39))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key method.

Roles: 2 declared. Attributes: unknown: Attribute collection could not be decoded. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 0 | unbounded | parameter |
| result | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

#### destructor ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L40))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key destructor.

Roles: 0 declared. Attributes: 4 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| access | false | public, protected, private |
| virtual | false | true, false |
| pure | false | true, false |
| noexcept | false | true, false |

No feature-key case association at this exact authored profile.

#### parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L41))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L42))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-parameter.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| default | 0 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| parameter-key | false | typename, class |

No feature-key case association at this exact authored profile.

#### named-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L43))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key named-type.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| segment | 1 | unbounded | name-segment |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| global | false | true, false |
| typename | false | true, false |

No feature-key case association at this exact authored profile.

#### name-segment ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L44))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key name-segment.

Roles: 1 declared. Attributes: 2 declared. References: target.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| argument | 0 | unbounded | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| template | false | true, false |

No feature-key case association at this exact authored profile.

#### builtin-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L45))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key builtin-type.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### pointer-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L46))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key pointer-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

#### lvalue-reference-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L47))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key lvalue-reference-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

#### rvalue-reference-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L48))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key rvalue-reference-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

#### cv-qualified-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L49))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key cv-qualified-type.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| const | false | true, false |
| volatile | false | true, false |

No feature-key case association at this exact authored profile.

#### array-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L50))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key array-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| length | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### function-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L51))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key function-type.

Roles: 2 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 0 | unbounded | parameter |
| result | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| noexcept | false | true, false |
| variadic | false | true, false |

No feature-key case association at this exact authored profile.

#### member-pointer-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L52))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key member-pointer-type.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| owner | 1 | 1 | named-type |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

Authored scope: [docs/syntax.md](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md); profile: unspecified (authored profile absent).

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjYyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvY3BwL2dvI1Byb3ZpZGVyLlBhcnNlMjA3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcGFyc2UvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wYXJzZS9nby5JbnB1dCwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28uTGFuZ3VhZ2VdXTU2OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvY3BwL2dvLlByb3ZpZGVyNTA6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28jTGFuZ3VhZ2UzOmNwcA"></a>

## C++ / github.com/Bitspark/bittype/parse/syntax/cpp/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/parse/syntax/cpp/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/cpp/go/parse.go#L441). Receiver: github.com/Bitspark/bittype/parse/syntax/cpp/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/parse/syntax/cpp/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/cpp/go/parse.go#L441) → [github.com/Bitspark/bittype/parse/syntax/cpp/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/cpp/go/parse.go#L25) ([call](https://github.com/Bitspark/bittype/blob/main/parse/syntax/cpp/go/parse.go#L442)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/cpp/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/cpp/go/parse.go#L25)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjYyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9nb2xhbmcvZ28jUHJvdmlkZXIuRm9ybWF0MjU5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvZm9ybWF0L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9hc3QvZ28uVmlld1tnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nby5MYW5ndWFnZV0sIGdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9jc3QvZ28uQ1NUW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvZ29sYW5nL2dvLkxhbmd1YWdlXV01NTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJvdmlkZXIvZ29sYW5nL2dvLlByb3ZpZGVyNTM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28jTGFuZ3VhZ2UyOmdv"></a>

## Go / github.com/Bitspark/bittype/provider/golang/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/provider/golang/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L23). Receiver: github.com/Bitspark/bittype/provider/golang/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/golang/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L23) → [github.com/Bitspark/bittype/format/syntax/golang/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/golang/go/format.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L24)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/golang/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/golang/go/format.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md#go](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#go); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| go-concrete-pipeline | pagination-and-interface | stable-formatting | unspecified (authored profile absent) | [TestGoConcretePipeline](https://github.com/Bitspark/bittype/blob/main/integration/go/text_test.go#L15) (external) |

### Declared target [github.com/Bitspark/bittype/provider/golang/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/syntax.go#L23)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L62) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L104) (external) |
| ast-admission | malformed-carriers | rejection-without-output | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L191) (external) |
| native-compiler-qualification | formatted-fixture-contract | compiler-fixture-contract | unspecified (authored profile absent) | [TestNativeCompilerFixtures](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L263) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L118) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjYyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ob3Jvcy9nbyNQcm92aWRlci5BbmFseXplMTEyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvYW5hbHlzaXMvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nby5MYW5ndWFnZV01NDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJvdmlkZXIvaG9yb3MvZ28uUHJvdmlkZXI1MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvI0xhbmd1YWdlNTpob3Jvcw"></a>

## Horos / github.com/Bitspark/bittype/provider/horos/go#Provider.Analyze

Contract: github.com/Bitspark/bittype/contracts/analysis/go.Provider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/provider/horos/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/files.go#L40). Receiver: github.com/Bitspark/bittype/provider/horos/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/horos/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/files.go#L40) → [github.com/Bitspark/bittype/analyze/files/horos/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/horos/go/analyze.go#L37) ([call](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/files.go#L41)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/analyze/files/horos/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/horos/go/analyze.go#L37)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/analysis.md#horos-profile](https://github.com/Bitspark/bittype/blob/main/docs/analysis.md#horos-profile); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/provider/horos/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/files.go#L40)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | analysis-evidence-produced | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| independent-engine-stages | pagination-source | independently-composable | unspecified (authored profile absent) | [TestHorosSourceThroughIndependentEngineStages](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_text_test.go#L27) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjYyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ob3Jvcy9nbyNQcm92aWRlci5SZXNvbHZlMTE0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcmVzb2x1dGlvbi9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXTU0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9ob3Jvcy9nby5Qcm92aWRlcjUyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28jTGFuZ3VhZ2U1Omhvcm9z"></a>

## Horos / github.com/Bitspark/bittype/provider/horos/go#Provider.Resolve

Contract: github.com/Bitspark/bittype/contracts/resolution/go.Provider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/provider/horos/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/files.go#L37). Receiver: github.com/Bitspark/bittype/provider/horos/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/horos/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/files.go#L37) → [github.com/Bitspark/bittype/resolve/files/horos/go#ResolveReferences](https://github.com/Bitspark/bittype/blob/main/resolve/files/horos/go/resolver.go#L50) ([call](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/files.go#L38)); receiver ; type arguments .

Traversal boundary: callee is dynamic or outside implementation source.

### Declared target [github.com/Bitspark/bittype/provider/horos/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/provider/horos/go/files.go#L37)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | captured-graph-resolved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| independent-engine-stages | pagination-source | independently-composable | unspecified (authored profile absent) | [TestHorosSourceThroughIndependentEngineStages](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_text_test.go#L27) (external) |

### Declared target [github.com/Bitspark/bittype/resolve/files/horos/go#ResolveReferences](https://github.com/Bitspark/bittype/blob/main/resolve/files/horos/go/resolver.go#L50)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/resolution.md#horos-composition](https://github.com/Bitspark/bittype/blob/main/docs/resolution.md#horos-composition); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjYzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvcnVzdC9nbyNQcm92aWRlci5DaGVjazEwODpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL2NoZWNrL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nby5MYW5ndWFnZV01NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY2hlY2svc3ludGF4L3J1c3QvZ28uUHJvdmlkZXI1MTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28jTGFuZ3VhZ2U0OnJ1c3Q"></a>

## Rust / github.com/Bitspark/bittype/check/syntax/rust/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/check/syntax/rust/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L192). Receiver: github.com/Bitspark/bittype/check/syntax/rust/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/check/syntax/rust/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L192) → [github.com/Bitspark/bittype/check/syntax/rust/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L72) ([call](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L192)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/check/syntax/rust/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L72)

Related declaration; its cases and observations do not qualify the public wrapper.

Definition: [github.com/Bitspark/bittype/check/syntax/rust/go#rules](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L31); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### file ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L32))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key file.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| item | 0 | unbounded | module, use, type-alias, struct, tuple-struct, enum, trait, trait-impl |

No feature-key case association at this exact authored profile.

#### module ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L33))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key module.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| item | 0 | unbounded | module, use, type-alias, struct, tuple-struct, enum, trait, trait-impl |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### use ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L34))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key use.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| path | 1 | 1 | path-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-alias ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L35))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-alias.

Roles: 3 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| where | 0 | unbounded | where-predicate |
| type | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### struct ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L36))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key struct.

Roles: 3 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| where | 0 | unbounded | where-predicate |
| field | 0 | unbounded | field |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### tuple-struct ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L37))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-struct.

Roles: 3 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| where | 0 | unbounded | where-predicate |
| field | 0 | unbounded | tuple-field |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### enum ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L38))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key enum.

Roles: 3 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| where | 0 | unbounded | where-predicate |
| variant | 0 | unbounded | unit-variant, tuple-variant, record-variant |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### unit-variant ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L39))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key unit-variant.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### tuple-variant ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L40))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-variant.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| field | 0 | unbounded | tuple-field |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### record-variant ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L41))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key record-variant.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| field | 0 | unbounded | field |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### field ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L42))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key field.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### tuple-field ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L43))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-field.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### trait ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L44))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key trait.

Roles: 4 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| bound | 0 | unbounded | path-type, lifetime-reference |
| where | 0 | unbounded | where-predicate |
| member | 0 | unbounded | method-signature, associated-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### trait-impl ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L45))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key trait-impl.

Roles: 5 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| trait | 1 | 1 | path-type |
| subject | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |
| where | 0 | unbounded | where-predicate |
| member | 0 | unbounded | type-alias |

No feature-key case association at this exact authored profile.

#### associated-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L46))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key associated-type.

Roles: 3 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| bound | 0 | unbounded | path-type, lifetime-reference, relaxed-bound |
| where | 0 | unbounded | where-predicate |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### method-signature ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L47))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key method-signature.

Roles: 5 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| where | 0 | unbounded | where-predicate |
| receiver | 0 | 1 | receiver |
| parameter | 0 | unbounded | parameter |
| result | 0 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### receiver ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L48))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key receiver.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| lifetime | 0 | 1 | lifetime-reference |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| receiver-form | true | value, shared, mutable |

No feature-key case association at this exact authored profile.

#### parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L49))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L50))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-parameter.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| bound | 0 | unbounded | path-type, lifetime-reference, relaxed-bound |
| default | 0 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### lifetime-parameter ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L51))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key lifetime-parameter.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| bound | 0 | unbounded | lifetime-reference |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### lifetime-reference ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L52))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key lifetime-reference.

Roles: 0 declared. Attributes: 1 declared. References: target.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### where-predicate ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L53))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key where-predicate.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| subject | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type, lifetime-reference |
| bound | 1 | unbounded | path-type, lifetime-reference, relaxed-bound |

No feature-key case association at this exact authored profile.

#### relaxed-bound ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L54))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key relaxed-bound.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| path | 1 | 1 | path-type |

No feature-key case association at this exact authored profile.

#### path-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L55))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key path-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| segment | 1 | unbounded | path-segment |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| absolute | false | true, false |

No feature-key case association at this exact authored profile.

#### qualified-path-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L56))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key qualified-path-type.

Roles: 3 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| subject | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |
| trait | 0 | 1 | path-type |
| segment | 1 | unbounded | path-segment |

No feature-key case association at this exact authored profile.

#### path-segment ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L57))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key path-segment.

Roles: 1 declared. Attributes: 1 declared. References: target.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| argument | 0 | unbounded | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type, lifetime-reference, associated-type-binding |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### associated-type-binding ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L58))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key associated-type-binding.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| path | 1 | 1 | path-segment |
| type | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

No feature-key case association at this exact authored profile.

#### primitive-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L59))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key primitive-type.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | bool, char, str, f32, f64, i8, i16, i32, i64, i128, isize, u8, u16, u32, u64, u128, usize |

No feature-key case association at this exact authored profile.

#### reference-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L60))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key reference-type.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |
| lifetime | 0 | 1 | lifetime-reference |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| mutable | false | true, false |

No feature-key case association at this exact authored profile.

#### raw-pointer-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L61))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key raw-pointer-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| mutable | false | true, false |

No feature-key case association at this exact authored profile.

#### slice-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L62))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key slice-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

No feature-key case association at this exact authored profile.

#### array-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L63))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key array-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| length | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### tuple-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L64))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 0 | unbounded | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

No feature-key case association at this exact authored profile.

#### never-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L65))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key never-type.

Roles: 0 declared. Attributes: 0 declared. References: none.

No feature-key case association at this exact authored profile.

#### function-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L66))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key function-type.

Roles: 2 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 0 | unbounded | parameter |
| result | 0 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| unsafe | false | true, false |
| abi | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### trait-object-type ([source](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L67))

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key trait-object-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| bound | 1 | unbounded | path-type, lifetime-reference |

No feature-key case association at this exact authored profile.

Authored scope: [docs/syntax.md#rust](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#rust); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| associated-type-projections | independent-ast-and-references | accept | unspecified (authored profile absent) | [TestAssociatedTypeAdmission](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check_test.go#L36) (external) |
| associated-type-projections | invalid-shapes-and-ordering | located-rejection | unspecified (authored profile absent) | [TestInvalidAssociatedTypeShapes](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check_test.go#L53) (external) |
| generic-associated-declarations | independent-ast | accept | unspecified (authored profile absent) | [TestGenericAssociatedDeclarationAdmission](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/declarations_test.go#L45) (external) |
| generic-associated-declarations | relaxed-bound-context | context-sensitive-admission | unspecified (authored profile absent) | [TestRelaxedBoundContext](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/declarations_test.go#L53) (external) |
| generic-associated-declarations | invalid-declaration-contexts | located-rejection | unspecified (authored profile absent) | [TestGenericAssociatedContextRules](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/declarations_test.go#L73) (external) |

<a id="cap-ODpwcm92aWRlcjYzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9maWxlcy9ob3Jvcy9nbyNQcm92aWRlci5QYXJzZTI3MDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3BhcnNlL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcGFyc2UvZ28uRmlsZXNSZXF1ZXN0W2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc291cmNlL2dvLlNvdXJjZVtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXV01NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcGFyc2UvZmlsZXMvaG9yb3MvZ28uUHJvdmlkZXI1MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvI0xhbmd1YWdlNTpob3Jvcw"></a>

## Horos / github.com/Bitspark/bittype/parse/files/horos/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.FilesRequest&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;, github.com/Bitspark/bittype/source/go.Source&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/parse/files/horos/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/files/horos/go/parse.go#L49). Receiver: github.com/Bitspark/bittype/parse/files/horos/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

### Declared target [github.com/Bitspark/bittype/parse/files/horos/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/files/horos/go/parse.go#L49)

Authored scope: [docs/source.md#horos-composition-and-pagination](https://github.com/Bitspark/bittype/blob/main/docs/source.md#horos-composition-and-pagination); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| multi-file-pagination | receipt-preserving-round-trips | module-context-preserved | unspecified (authored profile absent) | [TestFilesPaginationAndReceiptPreservingRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/files_test.go#L56) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjYzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvcnVzdC9nbyNQcm92aWRlci5QYXJzZTIwODpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3BhcnNlL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcGFyc2UvZ28uSW5wdXQsIGdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9jc3QvZ28uQ1NUW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nby5MYW5ndWFnZV1dNTc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3BhcnNlL3N5bnRheC9ydXN0L2dvLlByb3ZpZGVyNTE6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvI0xhbmd1YWdlNDpydXN0"></a>

## Rust / github.com/Bitspark/bittype/parse/syntax/rust/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/parse/syntax/rust/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse.go#L612). Receiver: github.com/Bitspark/bittype/parse/syntax/rust/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/parse/syntax/rust/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse.go#L612) → [github.com/Bitspark/bittype/parse/syntax/rust/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse.go#L25) ([call](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse.go#L613)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/rust/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse.go#L25)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md#rust](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#rust); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| generic-associated-declarations | declaration-and-impl-structure | structure-preserved | unspecified (authored profile absent) | [TestGenericAssociatedDeclarationStructure](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/declarations_test.go#L16) (external) |
| generic-associated-declarations | unresolved-syntax | accept | unspecified (authored profile absent) | [TestGenericAssociatedSyntaxWithoutResolution](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/declarations_test.go#L46) (external) |
| generic-associated-declarations | malformed-declarations | located-error-without-output | unspecified (authored profile absent) | [TestMalformedGenericAssociatedDeclarations](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/declarations_test.go#L67) (external) |
| type-only-impl-scope | unsupported-impl-forms | unsupported-without-output | unspecified (authored profile absent) | [TestImplScopeRejection](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/declarations_test.go#L94) (external) |
| associated-type-projections | structured-projections-and-bindings | structure-preserved | unspecified (authored profile absent) | [TestAssociatedTypeStructure](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse_test.go#L17) (external) |
| associated-type-projections | unresolved-projection-syntax | accept | unspecified (authored profile absent) | [TestAssociatedTypeSyntaxWithoutResolution](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse_test.go#L57) (external) |
| associated-type-projections | malformed-and-bounded-projections | located-error-without-output | unspecified (authored profile absent) | [TestMalformedAssociatedTypes](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse_test.go#L80) (external) |
| associated-type-projections | associated-bound-syntax | unsupported-without-output | unspecified (authored profile absent) | [TestAssociatedTypeBoundsRemainUnsupported](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse_test.go#L106) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjYzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9nb2xhbmcvZ28jUHJvdmlkZXIuQW5hbHl6ZTExMzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL2FuYWx5c2lzL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvZ29sYW5nL2dvLkxhbmd1YWdlXTU1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9nb2xhbmcvZ28uUHJvdmlkZXI1MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nbyNMYW5ndWFnZTI6Z28"></a>

## Go / github.com/Bitspark/bittype/provider/golang/go#Provider.Analyze

Contract: github.com/Bitspark/bittype/contracts/analysis/go.Provider&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/provider/golang/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/files.go#L33). Receiver: github.com/Bitspark/bittype/provider/golang/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/golang/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/files.go#L33) → [github.com/Bitspark/bittype/analyze/files/golang/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/go/analyze.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/files.go#L34)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/analyze/files/golang/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/go/analyze.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/native-semantics.md#first-go-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-go-profile); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjYzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcm92aWRlci9nb2xhbmcvZ28jUHJvdmlkZXIuUmVzb2x2ZTExNTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3Jlc29sdXRpb24vZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28uTGFuZ3VhZ2VdNTU6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3Byb3ZpZGVyL2dvbGFuZy9nby5Qcm92aWRlcjUzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvZ29sYW5nL2dvI0xhbmd1YWdlMjpnbw"></a>

## Go / github.com/Bitspark/bittype/provider/golang/go#Provider.Resolve

Contract: github.com/Bitspark/bittype/contracts/resolution/go.Provider&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/provider/golang/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/files.go#L30). Receiver: github.com/Bitspark/bittype/provider/golang/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/provider/golang/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/files.go#L30) → [github.com/Bitspark/bittype/resolve/files/golang/go#Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/go/resolver.go#L25) ([call](https://github.com/Bitspark/bittype/blob/main/provider/golang/go/files.go#L31)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/resolve/files/golang/go#Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/go/resolver.go#L25)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/native-semantics.md#first-go-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-go-profile); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjcxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L3R5cGVzY3JpcHQvZ28jUHJvdmlkZXIuRm9ybWF0MjY3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvZm9ybWF0L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9hc3QvZ28uVmlld1tnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28uTGFuZ3VhZ2VdLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28uTGFuZ3VhZ2VdXTY0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L3R5cGVzY3JpcHQvZ28uUHJvdmlkZXI1NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28jTGFuZ3VhZ2UxMDp0eXBlc2NyaXB0"></a>

## TypeScript / github.com/Bitspark/bittype/format/syntax/typescript/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/format/syntax/typescript/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/typescript/go/format.go#L281). Receiver: github.com/Bitspark/bittype/format/syntax/typescript/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/format/syntax/typescript/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/typescript/go/format.go#L281) → [github.com/Bitspark/bittype/format/syntax/typescript/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/typescript/go/format.go#L24) ([call](https://github.com/Bitspark/bittype/blob/main/format/syntax/typescript/go/format.go#L282)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/typescript/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/typescript/go/format.go#L24)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/syntax.md#typescript](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#typescript); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjcyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9hbmFseXplL2ZpbGVzL3R5cGVzY3JpcHQvZ28jUHJvdmlkZXIuQW5hbHl6ZTExNzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL2FuYWx5c2lzL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nby5MYW5ndWFnZV02NDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvYW5hbHl6ZS9maWxlcy90eXBlc2NyaXB0L2dvLlByb3ZpZGVyNTc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvI0xhbmd1YWdlMTA6dHlwZXNjcmlwdA"></a>

## TypeScript / github.com/Bitspark/bittype/analyze/files/typescript/go#Provider.Analyze

Contract: github.com/Bitspark/bittype/contracts/analysis/go.Provider&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/analyze/files/typescript/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/go/analyze.go#L57). Receiver: github.com/Bitspark/bittype/analyze/files/typescript/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/analyze/files/typescript/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/go/analyze.go#L57) → [github.com/Bitspark/bittype/analyze/files/typescript/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/go/analyze.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/go/analyze.go#L58)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/analyze/files/typescript/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/go/analyze.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

Authored scope: [docs/native-semantics.md#first-typescript-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-typescript-profile); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/analyze/files/typescript/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/go/analyze.go#L57)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | resolution-graph-preserved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| typescript-semantic-subset | generic-and-native-facts | native-identity-and-substitution-preserved | unspecified (authored profile absent) | [TestTypeScriptSemanticGenericsAndNativeFacts](https://github.com/Bitspark/bittype/blob/main/integration/go/typescript_semantics_test.go#L205) (external) |
| typescript-semantic-subset | rejected-and-unsupported-fixtures | no-partial-model | unspecified (authored profile absent) | [TestTypeScriptSemanticReports](https://github.com/Bitspark/bittype/blob/main/integration/go/typescript_semantics_test.go#L249) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjcyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9yZXNvbHZlL2ZpbGVzL3R5cGVzY3JpcHQvZ28jUHJvdmlkZXIuUmVzb2x2ZTExOTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3Jlc29sdXRpb24vZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvLkxhbmd1YWdlXTY0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9yZXNvbHZlL2ZpbGVzL3R5cGVzY3JpcHQvZ28uUHJvdmlkZXI1NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28jTGFuZ3VhZ2UxMDp0eXBlc2NyaXB0"></a>

## TypeScript / github.com/Bitspark/bittype/resolve/files/typescript/go#Provider.Resolve

Contract: github.com/Bitspark/bittype/contracts/resolution/go.Provider&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/resolve/files/typescript/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/typescript/go/resolver.go#L43). Receiver: github.com/Bitspark/bittype/resolve/files/typescript/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

### Declared target [github.com/Bitspark/bittype/resolve/files/typescript/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/typescript/go/resolver.go#L43)

Authored scope: [docs/native-semantics.md#first-typescript-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-typescript-profile); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | captured-graph-resolved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| typescript-semantic-subset | generic-and-native-facts | native-identity-and-substitution-preserved | unspecified (authored profile absent) | [TestTypeScriptSemanticGenericsAndNativeFacts](https://github.com/Bitspark/bittype/blob/main/integration/go/typescript_semantics_test.go#L205) (external) |
| typescript-semantic-subset | rejected-and-unsupported-fixtures | no-partial-model | unspecified (authored profile absent) | [TestTypeScriptSemanticReports](https://github.com/Bitspark/bittype/blob/main/integration/go/typescript_semantics_test.go#L249) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

