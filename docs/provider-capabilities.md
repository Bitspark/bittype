# Provider capability details

Generated from selected production declarations and explicit test associations.
These reports describe bounded implementation scope, not language completeness.

Inspection build: go1.27.3, linux/amd64; Go flags ; explicit build flags .

Definitions describe selected initializers; mutation, operation bodies and procedural semantics are not inferred. Profiles and expectations are authored assertions. A test function's observation is not separate execution of every associated case, nor evidence that a wrapper or every branch ran.

[Report contract](capability-reports.md) · [Provider matrices](providers.md) · [Existing CI invocations](https://github.com/Bitspark/bittype/actions/workflows/go.yml)

For a compact view, see the [language support overview](language-support.md).

No execution observations are included in this committed inventory. Per-invocation CI artifacts contain observations with their actual source, build selection and process result. Missing observations do not mean a test never ran.

<a id="cap-MTY6bWFwcGluZyBwcm92aWRlcjg1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmFuc2Zvcm0vc3ludGF4L2hvcm9zLXRvLWdvbGFuZy9tYXBwaW5nL2dvI1Byb3ZpZGVyLk9wZW4zNzA6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9nZW5lcmF0aW9uL2dvLk1hcHBpbmdQcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nby5MYW5ndWFnZSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyYW5zZm9ybS9zeW50YXgvaG9yb3Mvc291cmNlL2dvLlJlYWRlciwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyYW5zZm9ybS9zeW50YXgvZ29sYW5nL3RhcmdldC9nby5CdWlsZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdXTgwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmFuc2Zvcm0vc3ludGF4L2hvcm9zLXRvLWdvbGFuZy9tYXBwaW5nL2dvLlByb3ZpZGVyNTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nbyNMYW5ndWFnZTU6aG9yb3M1MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nbyNMYW5ndWFnZTI6Z28"></a>

## Horos → Go / github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/go#Provider.Open

Contract: github.com/Bitspark/bittype/contracts/generation/go.MappingProvider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language, github.com/Bitspark/bittype/syntax/golang/go.Language, github.com/Bitspark/bittype/transform/syntax/horos/source/go.Reader, github.com/Bitspark/bittype/transform/syntax/golang/target/go.Builder&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/go#Provider.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-golang/mapping/go/bindings.go#L19). Receiver: github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/go.Provider.

Selectable mapping from github.com/Bitspark/bittype/syntax/horos/go.Language to github.com/Bitspark/bittype/syntax/golang/go.Language. Its public Go contract identifies the required source and native APIs. The consumer supplies a registration and request.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/go#Provider.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-golang/mapping/go/bindings.go#L19) → [github.com/Bitspark/bittype/transform/syntax/horos/mapping/go#Rule.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos/mapping/go/provider.go#L16) ([call](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-golang/mapping/go/bindings.go#L20)); receiver github.com/Bitspark/bittype/transform/syntax/horos/mapping/go.Rule&#91;github.com/Bitspark/bittype/syntax/golang/go.Language, github.com/Bitspark/bittype/transform/syntax/golang/target/go.Builder&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/go#Provider.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-golang/mapping/go/bindings.go#L19)

<a id="binding-OTE1Ok1UWTZiV0Z3Y0dsdVp5QndjbTkyYVdSbGNqZzFPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzkwY21GdWMyWnZjbTB2YzNsdWRHRjRMMmh2Y205ekxYUnZMV2R2YkdGdVp5OXRZWEJ3YVc1bkwyZHZJMUJ5YjNacFpHVnlMazl3Wlc0ek56QTZaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMMk52Ym5SeVlXTjBjeTluWlc1bGNtRjBhVzl1TDJkdkxrMWhjSEJwYm1kUWNtOTJhV1JsY2x0bmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDJodmNtOXpMMmR2TGt4aGJtZDFZV2RsTENCbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDJkdmJHRnVaeTluYnk1TVlXNW5kV0ZuWlN3Z1oybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzUnlZVzV6Wm05eWJTOXplVzUwWVhndmFHOXliM012YzI5MWNtTmxMMmR2TGxKbFlXUmxjaXdnWjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNSeVlXNXpabTl5YlM5emVXNTBZWGd2WjI5c1lXNW5MM1JoY21kbGRDOW5ieTVDZFdsc1pHVnlXMmRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl6ZVc1MFlYZ3ZhRzl5YjNNdloyOHVUR0Z1WjNWaFoyVmRYVGd3T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5MGNtRnVjMlp2Y20wdmMzbHVkR0Y0TDJodmNtOXpMWFJ2TFdkdmJHRnVaeTl0WVhCd2FXNW5MMmR2TGxCeWIzWnBaR1Z5TlRJNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzTjViblJoZUM5b2IzSnZjeTluYnlOTVlXNW5kV0ZuWlRVNmFHOXliM00xTXpwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDJkdmJHRnVaeTluYnlOTVlXNW5kV0ZuWlRJNloyODg1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmFuc2Zvcm0vc3ludGF4L2hvcm9zLXRvLWdvbGFuZy9tYXBwaW5nL2dvI1Byb3ZpZGVyLk9wZW41OnNjb3BlMTojNTc6ZG9jcy9nZW5lcmF0aW9uLWNvbXBvc2l0aW9uLm1kI2NvbnN0cnVjdGlvbi1hbmQtc2VsZWN0aW9uNTpmYWxzZTA6"></a>

Authored scope: [docs/generation-composition.md#construction-and-selection](https://github.com/Bitspark/bittype/blob/main/docs/generation-composition.md#construction-and-selection); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-MTY6bWFwcGluZyBwcm92aWRlcjg5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmFuc2Zvcm0vc3ludGF4L2hvcm9zLXRvLXR5cGVzY3JpcHQvbWFwcGluZy9nbyNQcm92aWRlci5PcGVuMzc4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvZ2VuZXJhdGlvbi9nby5NYXBwaW5nUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nby5MYW5ndWFnZSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvLkxhbmd1YWdlLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJhbnNmb3JtL3N5bnRheC9ob3Jvcy9zb3VyY2UvZ28uUmVhZGVyLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJhbnNmb3JtL3N5bnRheC90eXBlc2NyaXB0L3RhcmdldC9nby5CdWlsZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdXTg0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmFuc2Zvcm0vc3ludGF4L2hvcm9zLXRvLXR5cGVzY3JpcHQvbWFwcGluZy9nby5Qcm92aWRlcjUyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28jTGFuZ3VhZ2U1Omhvcm9zNTc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvI0xhbmd1YWdlMTA6dHlwZXNjcmlwdA"></a>

## Horos → TypeScript / github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/go#Provider.Open

Contract: github.com/Bitspark/bittype/contracts/generation/go.MappingProvider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language, github.com/Bitspark/bittype/syntax/typescript/go.Language, github.com/Bitspark/bittype/transform/syntax/horos/source/go.Reader, github.com/Bitspark/bittype/transform/syntax/typescript/target/go.Builder&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/go#Provider.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-typescript/mapping/go/bindings.go#L19). Receiver: github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/go.Provider.

Selectable mapping from github.com/Bitspark/bittype/syntax/horos/go.Language to github.com/Bitspark/bittype/syntax/typescript/go.Language. Its public Go contract identifies the required source and native APIs. The consumer supplies a registration and request.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/go#Provider.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-typescript/mapping/go/bindings.go#L19) → [github.com/Bitspark/bittype/transform/syntax/horos/mapping/go#Rule.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos/mapping/go/provider.go#L16) ([call](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-typescript/mapping/go/bindings.go#L20)); receiver github.com/Bitspark/bittype/transform/syntax/horos/mapping/go.Rule&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language, github.com/Bitspark/bittype/transform/syntax/typescript/target/go.Builder&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/go#Provider.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-typescript/mapping/go/bindings.go#L19)

<a id="binding-OTU0Ok1UWTZiV0Z3Y0dsdVp5QndjbTkyYVdSbGNqZzVPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzkwY21GdWMyWnZjbTB2YzNsdWRHRjRMMmh2Y205ekxYUnZMWFI1Y0dWelkzSnBjSFF2YldGd2NHbHVaeTluYnlOUWNtOTJhV1JsY2k1UGNHVnVNemM0T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdloyVnVaWEpoZEdsdmJpOW5ieTVOWVhCd2FXNW5VSEp2ZG1sa1pYSmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOW9iM0p2Y3k5bmJ5NU1ZVzVuZFdGblpTd2daMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOTBlWEJsYzJOeWFYQjBMMmR2TGt4aGJtZDFZV2RsTENCbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmRISmhibk5tYjNKdEwzTjViblJoZUM5b2IzSnZjeTl6YjNWeVkyVXZaMjh1VW1WaFpHVnlMQ0JuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2ZEhKaGJuTm1iM0p0TDNONWJuUmhlQzkwZVhCbGMyTnlhWEIwTDNSaGNtZGxkQzluYnk1Q2RXbHNaR1Z5VzJkcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5emVXNTBZWGd2YUc5eWIzTXZaMjh1VEdGdVozVmhaMlZkWFRnME9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOTBjbUZ1YzJadmNtMHZjM2x1ZEdGNEwyaHZjbTl6TFhSdkxYUjVjR1Z6WTNKcGNIUXZiV0Z3Y0dsdVp5OW5ieTVRY205MmFXUmxjalV5T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5emVXNTBZWGd2YUc5eWIzTXZaMjhqVEdGdVozVmhaMlUxT21odmNtOXpOVGM2WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNONWJuUmhlQzkwZVhCbGMyTnlhWEIwTDJkdkkweGhibWQxWVdkbE1UQTZkSGx3WlhOamNtbHdkQTg5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmFuc2Zvcm0vc3ludGF4L2hvcm9zLXRvLXR5cGVzY3JpcHQvbWFwcGluZy9nbyNQcm92aWRlci5PcGVuNTpzY29wZTE6IzU3OmRvY3MvZ2VuZXJhdGlvbi1jb21wb3NpdGlvbi5tZCNjb25zdHJ1Y3Rpb24tYW5kLXNlbGVjdGlvbjU6ZmFsc2UwOg"></a>

Authored scope: [docs/generation-composition.md#construction-and-selection](https://github.com/Bitspark/bittype/blob/main/docs/generation-composition.md#construction-and-selection); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-MTY6bWFwcGluZyBwcm92aWRlcjgyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmFuc2Zvcm0vc3ludGF4L2hvcm9zLXRvLWNwcC9tYXBwaW5nL2dvI1Byb3ZpZGVyLk9wZW4zNjQ6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9nZW5lcmF0aW9uL2dvLk1hcHBpbmdQcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nby5MYW5ndWFnZSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyYW5zZm9ybS9zeW50YXgvaG9yb3Mvc291cmNlL2dvLlJlYWRlciwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyYW5zZm9ybS9zeW50YXgvY3BwL3RhcmdldC9nby5CdWlsZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdXTc3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmFuc2Zvcm0vc3ludGF4L2hvcm9zLXRvLWNwcC9tYXBwaW5nL2dvLlByb3ZpZGVyNTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nbyNMYW5ndWFnZTU6aG9yb3M1MDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nbyNMYW5ndWFnZTM6Y3Bw"></a>

## Horos → C++ / github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/go#Provider.Open

Contract: github.com/Bitspark/bittype/contracts/generation/go.MappingProvider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language, github.com/Bitspark/bittype/syntax/cpp/go.Language, github.com/Bitspark/bittype/transform/syntax/horos/source/go.Reader, github.com/Bitspark/bittype/transform/syntax/cpp/target/go.Builder&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/go#Provider.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-cpp/mapping/go/bindings.go#L19). Receiver: github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/go.Provider.

Selectable mapping from github.com/Bitspark/bittype/syntax/horos/go.Language to github.com/Bitspark/bittype/syntax/cpp/go.Language. Its public Go contract identifies the required source and native APIs. The consumer supplies a registration and request.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/go#Provider.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-cpp/mapping/go/bindings.go#L19) → [github.com/Bitspark/bittype/transform/syntax/horos/mapping/go#Rule.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos/mapping/go/provider.go#L16) ([call](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-cpp/mapping/go/bindings.go#L20)); receiver github.com/Bitspark/bittype/transform/syntax/horos/mapping/go.Rule&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language, github.com/Bitspark/bittype/transform/syntax/cpp/target/go.Builder&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/go#Provider.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-cpp/mapping/go/bindings.go#L19)

<a id="binding-ODk2Ok1UWTZiV0Z3Y0dsdVp5QndjbTkyYVdSbGNqZ3lPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzkwY21GdWMyWnZjbTB2YzNsdWRHRjRMMmh2Y205ekxYUnZMV053Y0M5dFlYQndhVzVuTDJkdkkxQnliM1pwWkdWeUxrOXdaVzR6TmpRNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwyTnZiblJ5WVdOMGN5OW5aVzVsY21GMGFXOXVMMmR2TGsxaGNIQnBibWRRY205MmFXUmxjbHRuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMmh2Y205ekwyZHZMa3hoYm1kMVlXZGxMQ0JuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMk53Y0M5bmJ5NU1ZVzVuZFdGblpTd2daMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM1J5WVc1elptOXliUzl6ZVc1MFlYZ3ZhRzl5YjNNdmMyOTFjbU5sTDJkdkxsSmxZV1JsY2l3Z1oybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzUnlZVzV6Wm05eWJTOXplVzUwWVhndlkzQndMM1JoY21kbGRDOW5ieTVDZFdsc1pHVnlXMmRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl6ZVc1MFlYZ3ZhRzl5YjNNdloyOHVUR0Z1WjNWaFoyVmRYVGMzT21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5MGNtRnVjMlp2Y20wdmMzbHVkR0Y0TDJodmNtOXpMWFJ2TFdOd2NDOXRZWEJ3YVc1bkwyZHZMbEJ5YjNacFpHVnlOVEk2WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNONWJuUmhlQzlvYjNKdmN5OW5ieU5NWVc1bmRXRm5aVFU2YUc5eWIzTTFNRHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMk53Y0M5bmJ5Tk1ZVzVuZFdGblpUTTZZM0J3ODI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyYW5zZm9ybS9zeW50YXgvaG9yb3MtdG8tY3BwL21hcHBpbmcvZ28jUHJvdmlkZXIuT3BlbjU6c2NvcGUxOiM1Nzpkb2NzL2dlbmVyYXRpb24tY29tcG9zaXRpb24ubWQjY29uc3RydWN0aW9uLWFuZC1zZWxlY3Rpb241OmZhbHNlMDo"></a>

Authored scope: [docs/generation-composition.md#construction-and-selection](https://github.com/Bitspark/bittype/blob/main/docs/generation-composition.md#construction-and-selection); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-MTY6bWFwcGluZyBwcm92aWRlcjgzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmFuc2Zvcm0vc3ludGF4L2hvcm9zLXRvLXJ1c3QvbWFwcGluZy9nbyNQcm92aWRlci5PcGVuMzY2OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvZ2VuZXJhdGlvbi9nby5NYXBwaW5nUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nby5MYW5ndWFnZSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvLkxhbmd1YWdlLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJhbnNmb3JtL3N5bnRheC9ob3Jvcy9zb3VyY2UvZ28uUmVhZGVyLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJhbnNmb3JtL3N5bnRheC9ydXN0L3RhcmdldC9nby5CdWlsZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdXTc4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmFuc2Zvcm0vc3ludGF4L2hvcm9zLXRvLXJ1c3QvbWFwcGluZy9nby5Qcm92aWRlcjUyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28jTGFuZ3VhZ2U1Omhvcm9zNTE6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvI0xhbmd1YWdlNDpydXN0"></a>

## Horos → Rust / github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/go#Provider.Open

Contract: github.com/Bitspark/bittype/contracts/generation/go.MappingProvider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language, github.com/Bitspark/bittype/syntax/rust/go.Language, github.com/Bitspark/bittype/transform/syntax/horos/source/go.Reader, github.com/Bitspark/bittype/transform/syntax/rust/target/go.Builder&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/go#Provider.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-rust/mapping/go/bindings.go#L19). Receiver: github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/go.Provider.

Selectable mapping from github.com/Bitspark/bittype/syntax/horos/go.Language to github.com/Bitspark/bittype/syntax/rust/go.Language. Its public Go contract identifies the required source and native APIs. The consumer supplies a registration and request.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/go#Provider.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-rust/mapping/go/bindings.go#L19) → [github.com/Bitspark/bittype/transform/syntax/horos/mapping/go#Rule.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos/mapping/go/provider.go#L16) ([call](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-rust/mapping/go/bindings.go#L20)); receiver github.com/Bitspark/bittype/transform/syntax/horos/mapping/go.Rule&#91;github.com/Bitspark/bittype/syntax/rust/go.Language, github.com/Bitspark/bittype/transform/syntax/rust/target/go.Builder&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/go#Provider.Open](https://github.com/Bitspark/bittype/blob/main/transform/syntax/horos-to-rust/mapping/go/bindings.go#L19)

<a id="binding-OTA0Ok1UWTZiV0Z3Y0dsdVp5QndjbTkyYVdSbGNqZ3pPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzkwY21GdWMyWnZjbTB2YzNsdWRHRjRMMmh2Y205ekxYUnZMWEoxYzNRdmJXRndjR2x1Wnk5bmJ5TlFjbTkyYVdSbGNpNVBjR1Z1TXpZMk9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOWpiMjUwY21GamRITXZaMlZ1WlhKaGRHbHZiaTluYnk1TllYQndhVzVuVUhKdmRtbGtaWEpiWjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNONWJuUmhlQzlvYjNKdmN5OW5ieTVNWVc1bmRXRm5aU3dnWjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNONWJuUmhlQzl5ZFhOMEwyZHZMa3hoYm1kMVlXZGxMQ0JuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2ZEhKaGJuTm1iM0p0TDNONWJuUmhlQzlvYjNKdmN5OXpiM1Z5WTJVdloyOHVVbVZoWkdWeUxDQm5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZkSEpoYm5ObWIzSnRMM041Ym5SaGVDOXlkWE4wTDNSaGNtZGxkQzluYnk1Q2RXbHNaR1Z5VzJkcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5emVXNTBZWGd2YUc5eWIzTXZaMjh1VEdGdVozVmhaMlZkWFRjNE9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOTBjbUZ1YzJadmNtMHZjM2x1ZEdGNEwyaHZjbTl6TFhSdkxYSjFjM1F2YldGd2NHbHVaeTluYnk1UWNtOTJhV1JsY2pVeU9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOXplVzUwWVhndmFHOXliM012WjI4alRHRnVaM1ZoWjJVMU9taHZjbTl6TlRFNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzTjViblJoZUM5eWRYTjBMMmR2STB4aGJtZDFZV2RsTkRweWRYTjA4MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJhbnNmb3JtL3N5bnRheC9ob3Jvcy10by1ydXN0L21hcHBpbmcvZ28jUHJvdmlkZXIuT3BlbjU6c2NvcGUxOiM1Nzpkb2NzL2dlbmVyYXRpb24tY29tcG9zaXRpb24ubWQjY29uc3RydWN0aW9uLWFuZC1zZWxlY3Rpb241OmZhbHNlMDo"></a>

Authored scope: [docs/generation-composition.md#construction-and-selection](https://github.com/Bitspark/bittype/blob/main/docs/generation-composition.md#construction-and-selection); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODppZGVudGl0eTE6IzU4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvbGFuZ3VhZ2UvZ28uUHJvdmlkZXI1MDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nby5MYW5ndWFnZTUwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvY3BwL2dvI0xhbmd1YWdlMzpjcHA"></a>

## C++ language identity

Contract: github.com/Bitspark/bittype/contracts/language/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

### Declared target [github.com/Bitspark/bittype/syntax/cpp/go#Language.ID](https://github.com/Bitspark/bittype/blob/main/syntax/cpp/go/nodes.go#L16)

<a id="binding-MjQ3Ok9EcHBaR1Z1ZEdsMGVURTZJelU0T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdmJHRnVaM1ZoWjJVdloyOHVVSEp2ZG1sa1pYSTFNRHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMk53Y0M5bmJ5NU1ZVzVuZFdGblpUVXdPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl6ZVc1MFlYZ3ZZM0J3TDJkdkkweGhibWQxWVdkbE16cGpjSEE1MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nbyNMYW5ndWFnZS5JRDU6c2NvcGUxOiM0MDpkb2NzL3N5bnRheC5tZCNpZGVudGl0eS1hbmQtcmVnaXN0cmF0aW9uNTpmYWxzZTA6"></a>

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

### Declared target [github.com/Bitspark/bittype/syntax/cpp/go#Language.Name](https://github.com/Bitspark/bittype/blob/main/syntax/cpp/go/nodes.go#L20)

<a id="binding-MjQ3Ok9EcHBaR1Z1ZEdsMGVURTZJelU0T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdmJHRnVaM1ZoWjJVdloyOHVVSEp2ZG1sa1pYSTFNRHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMk53Y0M5bmJ5NU1ZVzVuZFdGblpUVXdPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl6ZVc1MFlYZ3ZZM0J3TDJkdkkweGhibWQxWVdkbE16cGpjSEE1NTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nbyNMYW5ndWFnZS5OYW1lNTpzY29wZTE6IzQwOmRvY3Mvc3ludGF4Lm1kI2lkZW50aXR5LWFuZC1yZWdpc3RyYXRpb241OmZhbHNlMDo"></a>

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

<a id="binding-MjUxOk9EcHBaR1Z1ZEdsMGVURTZJelU0T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdmJHRnVaM1ZoWjJVdloyOHVVSEp2ZG1sa1pYSTFNVHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMM0oxYzNRdloyOHVUR0Z1WjNWaFoyVTFNVHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMM0oxYzNRdloyOGpUR0Z1WjNWaFoyVTBPbkoxYzNRNTQ6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvI0xhbmd1YWdlLklENTpzY29wZTE6IzQwOmRvY3Mvc3ludGF4Lm1kI2lkZW50aXR5LWFuZC1yZWdpc3RyYXRpb241OmZhbHNlMDo"></a>

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

### Declared target [github.com/Bitspark/bittype/syntax/rust/go#Language.Name](https://github.com/Bitspark/bittype/blob/main/syntax/rust/go/nodes.go#L20)

<a id="binding-MjUxOk9EcHBaR1Z1ZEdsMGVURTZJelU0T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdmJHRnVaM1ZoWjJVdloyOHVVSEp2ZG1sa1pYSTFNVHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMM0oxYzNRdloyOHVUR0Z1WjNWaFoyVTFNVHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMM0oxYzNRdloyOGpUR0Z1WjNWaFoyVTBPbkoxYzNRNTY6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvI0xhbmd1YWdlLk5hbWU1OnNjb3BlMTojNDA6ZG9jcy9zeW50YXgubWQjaWRlbnRpdHktYW5kLXJlZ2lzdHJhdGlvbjU6ZmFsc2UwOg"></a>

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

<a id="binding-MjU1Ok9EcHBaR1Z1ZEdsMGVURTZJelU0T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdmJHRnVaM1ZoWjJVdloyOHVVSEp2ZG1sa1pYSTFNanBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMmh2Y205ekwyZHZMa3hoYm1kMVlXZGxOVEk2WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNONWJuUmhlQzlvYjNKdmN5OW5ieU5NWVc1bmRXRm5aVFU2YUc5eWIzTTU1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28jTGFuZ3VhZ2UuSUQ1OnNjb3BlMTojNDA6ZG9jcy9zeW50YXgubWQjaWRlbnRpdHktYW5kLXJlZ2lzdHJhdGlvbjU6ZmFsc2UwOg"></a>

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

### Declared target [github.com/Bitspark/bittype/syntax/horos/go#Language.Name](https://github.com/Bitspark/bittype/blob/main/syntax/horos/go/nodes.go#L16)

<a id="binding-MjU1Ok9EcHBaR1Z1ZEdsMGVURTZJelU0T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdmJHRnVaM1ZoWjJVdloyOHVVSEp2ZG1sa1pYSTFNanBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMmh2Y205ekwyZHZMa3hoYm1kMVlXZGxOVEk2WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNONWJuUmhlQzlvYjNKdmN5OW5ieU5NWVc1bmRXRm5aVFU2YUc5eWIzTTU3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28jTGFuZ3VhZ2UuTmFtZTU6c2NvcGUxOiM0MDpkb2NzL3N5bnRheC5tZCNpZGVudGl0eS1hbmQtcmVnaXN0cmF0aW9uNTpmYWxzZTA6"></a>

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

<a id="binding-MjU0Ok9EcHBaR1Z1ZEdsMGVURTZJelU0T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdmJHRnVaM1ZoWjJVdloyOHVVSEp2ZG1sa1pYSTFNenBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMmR2YkdGdVp5OW5ieTVNWVc1bmRXRm5aVFV6T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5emVXNTBZWGd2WjI5c1lXNW5MMmR2STB4aGJtZDFZV2RsTWpwbmJ3NTY6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28jTGFuZ3VhZ2UuSUQ1OnNjb3BlMTojNDA6ZG9jcy9zeW50YXgubWQjaWRlbnRpdHktYW5kLXJlZ2lzdHJhdGlvbjU6ZmFsc2UwOg"></a>

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

### Declared target [github.com/Bitspark/bittype/syntax/golang/go#Language.Name](https://github.com/Bitspark/bittype/blob/main/syntax/golang/go/nodes.go#L20)

<a id="binding-MjU0Ok9EcHBaR1Z1ZEdsMGVURTZJelU0T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdmJHRnVaM1ZoWjJVdloyOHVVSEp2ZG1sa1pYSTFNenBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMmR2YkdGdVp5OW5ieTVNWVc1bmRXRm5aVFV6T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5emVXNTBZWGd2WjI5c1lXNW5MMmR2STB4aGJtZDFZV2RsTWpwbmJ3NTg6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28jTGFuZ3VhZ2UuTmFtZTU6c2NvcGUxOiM0MDpkb2NzL3N5bnRheC5tZCNpZGVudGl0eS1hbmQtcmVnaXN0cmF0aW9uNTpmYWxzZTA6"></a>

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

<a id="binding-Mjc2Ok9EcHBaR1Z1ZEdsMGVURTZJelU0T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdmJHRnVaM1ZoWjJVdloyOHVVSEp2ZG1sa1pYSTFOenBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMM1I1Y0dWelkzSnBjSFF2WjI4dVRHRnVaM1ZoWjJVMU56cG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwzUjVjR1Z6WTNKcGNIUXZaMjhqVEdGdVozVmhaMlV4TURwMGVYQmxjMk55YVhCMDYwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nbyNMYW5ndWFnZS5JRDU6c2NvcGUxOiM0MDpkb2NzL3N5bnRheC5tZCNpZGVudGl0eS1hbmQtcmVnaXN0cmF0aW9uNTpmYWxzZTA6"></a>

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

### Declared target [github.com/Bitspark/bittype/syntax/typescript/go#Language.Name](https://github.com/Bitspark/bittype/blob/main/syntax/typescript/go/nodes.go#L20)

<a id="binding-Mjc2Ok9EcHBaR1Z1ZEdsMGVURTZJelU0T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdmJHRnVaM1ZoWjJVdloyOHVVSEp2ZG1sa1pYSTFOenBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMM1I1Y0dWelkzSnBjSFF2WjI4dVRHRnVaM1ZoWjJVMU56cG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwzUjVjR1Z6WTNKcGNIUXZaMjhqVEdGdVozVmhaMlV4TURwMGVYQmxjMk55YVhCMDYyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nbyNMYW5ndWFnZS5OYW1lNTpzY29wZTE6IzQwOmRvY3Mvc3ludGF4Lm1kI2lkZW50aXR5LWFuZC1yZWdpc3RyYXRpb241OmZhbHNlMDo"></a>

Authored scope: [docs/syntax.md#identity-and-registration](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#identity-and-registration); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| language-identity | stable-registration-and-name | expected-identity | unspecified (authored profile absent) | [TestLanguageIdentityInterface](https://github.com/Bitspark/bittype/blob/main/integration/go/syntax_test.go#L26) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvZ29sYW5nL3Byb3ZpZGVyL2dvI1Byb3ZpZGVyLkNoZWNrMTEwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvY2hlY2svZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28uTGFuZ3VhZ2VdNjg6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NoZWNrL3N5bnRheC9nb2xhbmcvcHJvdmlkZXIvZ28uUHJvdmlkZXI1MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nbyNMYW5ndWFnZTI6Z28"></a>

## Go / github.com/Bitspark/bittype/check/syntax/golang/provider/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/check/syntax/golang/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/provider/go/check.go#L87). Receiver: github.com/Bitspark/bittype/check/syntax/golang/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/check/syntax/golang/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/provider/go/check.go#L87) → [github.com/Bitspark/bittype/check/syntax/golang/provider/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/provider/go/check.go#L45) ([call](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/provider/go/check.go#L87)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/check/syntax/golang/provider/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/provider/go/check.go#L45)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NDQzOk9EcHdjbTkyYVdSbGNqYzBPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYUdWamF5OXplVzUwWVhndloyOXNZVzVuTDNCeWIzWnBaR1Z5TDJkdkkxQnliM1pwWkdWeUxrTm9aV05yTVRFd09tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOWpiMjUwY21GamRITXZZMmhsWTJzdloyOHVVSEp2ZG1sa1pYSmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOW5iMnhoYm1jdloyOHVUR0Z1WjNWaFoyVmROamc2WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDJOb1pXTnJMM041Ym5SaGVDOW5iMnhoYm1jdmNISnZkbWxrWlhJdloyOHVVSEp2ZG1sa1pYSTFNenBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMmR2YkdGdVp5OW5ieU5NWVc1bmRXRm5aVEk2WjI4NjU6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NoZWNrL3N5bnRheC9nb2xhbmcvcHJvdmlkZXIvZ28jQ2hlY2s3OmdyYW1tYXI2NTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY2hlY2svc3ludGF4L2dvbGFuZy9wcm92aWRlci9nbyNydWxlczA6NTpmYWxzZTA6"></a>

Definition: [github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/provider/go/check.go#L17); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### file (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key file.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| import | 0 | unbounded | import |
| declaration | 0 | unbounded | type-declaration, type-alias |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| package | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### import (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key import.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| path | true | opaque: Predicate present; behavior not interpreted |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-declaration (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-declaration.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 0 | unbounded | type-parameter |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-alias (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-alias.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 0 | unbounded | type-parameter |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-parameter (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-parameter.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| constraint | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type, type-union, type-approximation |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### named-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key named-type.

Roles: 0 declared. Attributes: 1 declared. References: target.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### qualified-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key qualified-type.

Roles: 0 declared. Attributes: 2 declared. References: target.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| qualifier | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-application (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-application.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| constructor | 1 | 1 | named-type, qualified-type |
| argument | 1 | unbounded | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

#### struct-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key struct-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| field | 0 | unbounded | field, embedded-field |

No feature-key case association at this exact authored profile.

#### field (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key field.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| tag | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### embedded-field (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key embedded-field.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| tag | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### interface-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key interface-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| member | 0 | unbounded | method, embedded-type |

No feature-key case association at this exact authored profile.

#### method (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key method.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| signature | 1 | 1 | function-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### embedded-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key embedded-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type, type-union, type-approximation |

No feature-key case association at this exact authored profile.

#### function-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key function-type.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 0 | unbounded | parameter |
| result | 0 | unbounded | parameter |

No feature-key case association at this exact authored profile.

#### parameter (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | false | opaque: Predicate present; behavior not interpreted |
| variadic | false | true, false |

No feature-key case association at this exact authored profile.

#### pointer-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key pointer-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

#### slice-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key slice-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

#### array-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key array-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| length | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### map-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key map-type.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| key | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

#### channel-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key channel-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| direction | false | both, send, receive |

No feature-key case association at this exact authored profile.

#### type-union (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-union.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| term | 2 | unbounded | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type, type-approximation |

No feature-key case association at this exact authored profile.

#### type-approximation (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/golang/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-approximation.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, qualified-type, type-application, struct-type, interface-type, function-type, pointer-type, slice-type, array-type, map-type, channel-type |

No feature-key case association at this exact authored profile.

<a id="binding-NDQzOk9EcHdjbTkyYVdSbGNqYzBPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYUdWamF5OXplVzUwWVhndloyOXNZVzVuTDNCeWIzWnBaR1Z5TDJkdkkxQnliM1pwWkdWeUxrTm9aV05yTVRFd09tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOWpiMjUwY21GamRITXZZMmhsWTJzdloyOHVVSEp2ZG1sa1pYSmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOW5iMnhoYm1jdloyOHVUR0Z1WjNWaFoyVmROamc2WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDJOb1pXTnJMM041Ym5SaGVDOW5iMnhoYm1jdmNISnZkbWxrWlhJdloyOHVVSEp2ZG1sa1pYSTFNenBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMmR2YkdGdVp5OW5ieU5NWVc1bmRXRm5aVEk2WjI4NjU6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NoZWNrL3N5bnRheC9nb2xhbmcvcHJvdmlkZXIvZ28jQ2hlY2s1OnNjb3BlMTojMTc6ZG9jcy9zeW50YXgubWQjZ281OmZhbHNlMDo"></a>

Authored scope: [docs/syntax.md#go](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#go); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/check/syntax/golang/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/provider/go/check.go#L87)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| ast-admission | malformed-carriers | reject | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L193) (external) |

<a id="cap-ODpwcm92aWRlcjc0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvZmlsZXMvaG9yb3MvcHJvdmlkZXIvZ28jUHJvdmlkZXIuRm9ybWF0MjcyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvZm9ybWF0L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvZm9ybWF0L2dvLkZpbGVzUmVxdWVzdFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3NvdXJjZS9nby5Tb3VyY2VbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nby5MYW5ndWFnZV1dNjc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2Zvcm1hdC9maWxlcy9ob3Jvcy9wcm92aWRlci9nby5Qcm92aWRlcjUyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28jTGFuZ3VhZ2U1Omhvcm9z"></a>

## Horos / github.com/Bitspark/bittype/format/files/horos/provider/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/contracts/format/go.FilesRequest&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;, github.com/Bitspark/bittype/source/go.Source&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/format/files/horos/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/files/horos/provider/go/format.go#L54). Receiver: github.com/Bitspark/bittype/format/files/horos/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

### Declared target [github.com/Bitspark/bittype/format/files/horos/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/files/horos/provider/go/format.go#L54)

<a id="binding-NjYwOk9EcHdjbTkyYVdSbGNqYzBPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzltYjNKdFlYUXZabWxzWlhNdmFHOXliM012Y0hKdmRtbGtaWEl2WjI4alVISnZkbWxrWlhJdVJtOXliV0YwTWpjeU9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOWpiMjUwY21GamRITXZabTl5YldGMEwyZHZMbEJ5YjNacFpHVnlXMmRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYjI1MGNtRmpkSE12Wm05eWJXRjBMMmR2TGtacGJHVnpVbVZ4ZFdWemRGdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwyaHZjbTl6TDJkdkxreGhibWQxWVdkbFhTd2daMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM052ZFhKalpTOW5ieTVUYjNWeVkyVmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOW9iM0p2Y3k5bmJ5NU1ZVzVuZFdGblpWMWROamM2WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDJadmNtMWhkQzltYVd4bGN5OW9iM0p2Y3k5d2NtOTJhV1JsY2k5bmJ5NVFjbTkyYVdSbGNqVXlPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl6ZVc1MFlYZ3ZhRzl5YjNNdloyOGpUR0Z1WjNWaFoyVTFPbWh2Y205ejc0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvZmlsZXMvaG9yb3MvcHJvdmlkZXIvZ28jUHJvdmlkZXIuRm9ybWF0NTpzY29wZTE6IzQ3OmRvY3Mvc291cmNlLm1kI2hvcm9zLWNvbXBvc2l0aW9uLWFuZC1wYWdpbmF0aW9uNTpmYWxzZTA6"></a>

Authored scope: [docs/source.md#horos-composition-and-pagination](https://github.com/Bitspark/bittype/blob/main/docs/source.md#horos-composition-and-pagination); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| multi-file-pagination | receipt-preserving-round-trips | fresh-output-and-receipt | unspecified (authored profile absent) | [TestFilesPaginationAndReceiptPreservingRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/files_test.go#L85) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L3J1c3QvcHJvdmlkZXIvZ28jUHJvdmlkZXIuRm9ybWF0MjU1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvZm9ybWF0L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9hc3QvZ28uVmlld1tnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28uTGFuZ3VhZ2VdLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28uTGFuZ3VhZ2VdXTY3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L3J1c3QvcHJvdmlkZXIvZ28uUHJvdmlkZXI1MTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28jTGFuZ3VhZ2U0OnJ1c3Q"></a>

## Rust / github.com/Bitspark/bittype/format/syntax/rust/provider/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/format/syntax/rust/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/rust/provider/go/format.go#L339). Receiver: github.com/Bitspark/bittype/format/syntax/rust/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/format/syntax/rust/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/rust/provider/go/format.go#L339) → [github.com/Bitspark/bittype/format/syntax/rust/provider/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/rust/provider/go/format.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/format/syntax/rust/provider/go/format.go#L340)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/rust/provider/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/rust/provider/go/format.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NjM1Ok9EcHdjbTkyYVdSbGNqYzBPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzltYjNKdFlYUXZjM2x1ZEdGNEwzSjFjM1F2Y0hKdmRtbGtaWEl2WjI4alVISnZkbWxrWlhJdVJtOXliV0YwTWpVMU9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOWpiMjUwY21GamRITXZabTl5YldGMEwyZHZMbEJ5YjNacFpHVnlXMmRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzkwY21WbGN5OWhjM1F2WjI4dVZtbGxkMXRuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMM0oxYzNRdloyOHVUR0Z1WjNWaFoyVmRMQ0JuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2ZEhKbFpYTXZZM04wTDJkdkxrTlRWRnRuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMM0oxYzNRdloyOHVUR0Z1WjNWaFoyVmRYVFkzT21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5bWIzSnRZWFF2YzNsdWRHRjRMM0oxYzNRdmNISnZkbWxrWlhJdloyOHVVSEp2ZG1sa1pYSTFNVHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMM0oxYzNRdloyOGpUR0Z1WjNWaFoyVTBPbkoxYzNRNjU6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2Zvcm1hdC9zeW50YXgvcnVzdC9wcm92aWRlci9nbyNGb3JtYXQ1OnNjb3BlMTojMTk6ZG9jcy9zeW50YXgubWQjcnVzdDU6ZmFsc2UwOg"></a>

Authored scope: [docs/syntax.md#rust](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#rust); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/format/syntax/rust/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/rust/provider/go/format.go#L339)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L106) (external) |
| ast-admission | malformed-carriers | rejection-without-output | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L193) (external) |
| native-compiler-qualification | formatted-fixture-contract | compiler-fixture-contract | unspecified (authored profile absent) | [TestNativeCompilerFixtures](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L265) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L115) (external) |
| associated-type-projections | independent-ast-formatting | ast-preserved-and-golden-text | unspecified (authored profile absent) | [TestIndependentRustAssociatedFormatter](https://github.com/Bitspark/bittype/blob/main/integration/go/rust_associated_types_test.go#L40) (external) |
| generic-associated-declarations | independent-ast-formatting | ast-preserved-and-golden-text | unspecified (authored profile absent) | [TestIndependentRustDeclarationFormatter](https://github.com/Bitspark/bittype/blob/main/integration/go/rust_declarations_test.go#L41) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvZ29sYW5nL3Byb3ZpZGVyL2dvI1Byb3ZpZGVyLlBhcnNlMjEwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcGFyc2UvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wYXJzZS9nby5JbnB1dCwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28uTGFuZ3VhZ2VdXTY4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvZ29sYW5nL3Byb3ZpZGVyL2dvLlByb3ZpZGVyNTM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28jTGFuZ3VhZ2UyOmdv"></a>

## Go / github.com/Bitspark/bittype/parse/syntax/golang/provider/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/parse/syntax/golang/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/golang/provider/go/parse.go#L280). Receiver: github.com/Bitspark/bittype/parse/syntax/golang/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/parse/syntax/golang/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/golang/provider/go/parse.go#L280) → [github.com/Bitspark/bittype/parse/syntax/golang/provider/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/golang/provider/go/parse.go#L26) ([call](https://github.com/Bitspark/bittype/blob/main/parse/syntax/golang/provider/go/parse.go#L281)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/golang/provider/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/golang/provider/go/parse.go#L26)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NTc2Ok9EcHdjbTkyYVdSbGNqYzBPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3WVhKelpTOXplVzUwWVhndloyOXNZVzVuTDNCeWIzWnBaR1Z5TDJkdkkxQnliM1pwWkdWeUxsQmhjbk5sTWpFd09tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOWpiMjUwY21GamRITXZjR0Z5YzJVdloyOHVVSEp2ZG1sa1pYSmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMMk52Ym5SeVlXTjBjeTl3WVhKelpTOW5ieTVKYm5CMWRDd2daMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM1J5WldWekwyTnpkQzluYnk1RFUxUmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOW5iMnhoYm1jdloyOHVUR0Z1WjNWaFoyVmRYVFk0T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5d1lYSnpaUzl6ZVc1MFlYZ3ZaMjlzWVc1bkwzQnliM1pwWkdWeUwyZHZMbEJ5YjNacFpHVnlOVE02WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNONWJuUmhlQzluYjJ4aGJtY3ZaMjhqVEdGdVozVmhaMlV5T21kdjY1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvZ29sYW5nL3Byb3ZpZGVyL2dvI1BhcnNlNTpzY29wZTE6IzE3OmRvY3Mvc3ludGF4Lm1kI2dvNTpmYWxzZTA6"></a>

Authored scope: [docs/syntax.md#go](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#go); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| go-concrete-pipeline | pagination-and-interface | round-trip-preserved | unspecified (authored profile absent) | [TestGoConcretePipeline](https://github.com/Bitspark/bittype/blob/main/integration/go/text_test.go#L15) (external) |

### Declared target [github.com/Bitspark/bittype/parse/syntax/golang/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/golang/provider/go/parse.go#L280)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L115) (external) |
| source-boundary | unsupported-and-malformed-fixtures | located-error-without-output | unspecified (authored profile absent) | [TestUnsupportedAndMalformedSources](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L137) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvc2hhcmVkL3Byb3ZpZGVyL2dvI1Byb3ZpZGVyLlByaW50MTY0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcHJpbnQvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28uTGFuZ3VhZ2VdLCBzdHJpbmddMTIzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvc2hhcmVkL3Byb3ZpZGVyL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvZ29sYW5nL2dvLkxhbmd1YWdlXTUzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvZ29sYW5nL2dvI0xhbmd1YWdlMjpnbw"></a>

## Go / github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25). Receiver: github.com/Bitspark/bittype/print/syntax/shared/provider/go.Provider&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25) → [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L26)); receiver ; type arguments L.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NTkwOk9EcHdjbTkyYVdSbGNqYzBPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3Y21sdWRDOXplVzUwWVhndmMyaGhjbVZrTDNCeWIzWnBaR1Z5TDJkdkkxQnliM1pwWkdWeUxsQnlhVzUwTVRZME9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOWpiMjUwY21GamRITXZjSEpwYm5RdloyOHVVSEp2ZG1sa1pYSmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM1J5WldWekwyTnpkQzluYnk1RFUxUmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOW5iMnhoYm1jdloyOHVUR0Z1WjNWaFoyVmRMQ0J6ZEhKcGJtZGRNVEl6T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5d2NtbHVkQzl6ZVc1MFlYZ3ZjMmhoY21Wa0wzQnliM1pwWkdWeUwyZHZMbEJ5YjNacFpHVnlXMmRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl6ZVc1MFlYZ3ZaMjlzWVc1bkwyZHZMa3hoYm1kMVlXZGxYVFV6T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5emVXNTBZWGd2WjI5c1lXNW5MMmR2STB4aGJtZDFZV2RsTWpwbmJ3NjU6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3ByaW50L3N5bnRheC9zaGFyZWQvcHJvdmlkZXIvZ28jUHJpbnQ1OnNjb3BlMTojNDk6ZG9jcy9wYXJzaW5nLWZvcm1hdHRpbmcubWQjY2FycmllcnMtYW5kLW93bmVyc2hpcDU6ZmFsc2UwOg"></a>

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| go-concrete-pipeline | pagination-and-interface | source-retained | unspecified (authored profile absent) | [TestGoConcretePipeline](https://github.com/Bitspark/bittype/blob/main/integration/go/text_test.go#L15) (external) |

### Declared target [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L106) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L115) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvc2hhcmVkL3Byb3ZpZGVyL2dvI1Byb3ZpZGVyLlByaW50MTY4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcHJpbnQvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvLkxhbmd1YWdlXSwgc3RyaW5nXTEyNzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJpbnQvc3ludGF4L3NoYXJlZC9wcm92aWRlci9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28uTGFuZ3VhZ2VdNTc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvI0xhbmd1YWdlMTA6dHlwZXNjcmlwdA"></a>

## TypeScript / github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25). Receiver: github.com/Bitspark/bittype/print/syntax/shared/provider/go.Provider&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25) → [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L26)); receiver ; type arguments L.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NjE4Ok9EcHdjbTkyYVdSbGNqYzBPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3Y21sdWRDOXplVzUwWVhndmMyaGhjbVZrTDNCeWIzWnBaR1Z5TDJkdkkxQnliM1pwWkdWeUxsQnlhVzUwTVRZNE9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOWpiMjUwY21GamRITXZjSEpwYm5RdloyOHVVSEp2ZG1sa1pYSmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM1J5WldWekwyTnpkQzluYnk1RFUxUmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOTBlWEJsYzJOeWFYQjBMMmR2TGt4aGJtZDFZV2RsWFN3Z2MzUnlhVzVuWFRFeU56cG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjSEpwYm5RdmMzbHVkR0Y0TDNOb1lYSmxaQzl3Y205MmFXUmxjaTluYnk1UWNtOTJhV1JsY2x0bmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDNSNWNHVnpZM0pwY0hRdloyOHVUR0Z1WjNWaFoyVmROVGM2WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNONWJuUmhlQzkwZVhCbGMyTnlhWEIwTDJkdkkweGhibWQxWVdkbE1UQTZkSGx3WlhOamNtbHdkQTY1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvc2hhcmVkL3Byb3ZpZGVyL2dvI1ByaW50NTpzY29wZTE6IzQ5OmRvY3MvcGFyc2luZy1mb3JtYXR0aW5nLm1kI2NhcnJpZXJzLWFuZC1vd25lcnNoaXA1OmZhbHNlMDo"></a>

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L106) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L115) (external) |
| concrete-admission | malformed-concrete-tree | rejection-without-output | unspecified (authored profile absent) | [TestMalformedConcreteTree](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L220) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvc2hhcmVkL3Byb3ZpZGVyL2dvI1Byb3ZpZGVyLlByaW50MTYxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcHJpbnQvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28uTGFuZ3VhZ2VdLCBzdHJpbmddMTIwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvc2hhcmVkL3Byb3ZpZGVyL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvY3BwL2dvLkxhbmd1YWdlXTUwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvY3BwL2dvI0xhbmd1YWdlMzpjcHA"></a>

## C++ / github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25). Receiver: github.com/Bitspark/bittype/print/syntax/shared/provider/go.Provider&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25) → [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L26)); receiver ; type arguments L.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NTc5Ok9EcHdjbTkyYVdSbGNqYzBPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3Y21sdWRDOXplVzUwWVhndmMyaGhjbVZrTDNCeWIzWnBaR1Z5TDJkdkkxQnliM1pwWkdWeUxsQnlhVzUwTVRZeE9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOWpiMjUwY21GamRITXZjSEpwYm5RdloyOHVVSEp2ZG1sa1pYSmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM1J5WldWekwyTnpkQzluYnk1RFUxUmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOWpjSEF2WjI4dVRHRnVaM1ZoWjJWZExDQnpkSEpwYm1kZE1USXdPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3Y21sdWRDOXplVzUwWVhndmMyaGhjbVZrTDNCeWIzWnBaR1Z5TDJkdkxsQnliM1pwWkdWeVcyZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOXplVzUwWVhndlkzQndMMmR2TGt4aGJtZDFZV2RsWFRVd09tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOXplVzUwWVhndlkzQndMMmR2STB4aGJtZDFZV2RsTXpwamNIQTY1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvc2hhcmVkL3Byb3ZpZGVyL2dvI1ByaW50NTpzY29wZTE6IzQ5OmRvY3MvcGFyc2luZy1mb3JtYXR0aW5nLm1kI2NhcnJpZXJzLWFuZC1vd25lcnNoaXA1OmZhbHNlMDo"></a>

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L106) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L115) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvc2hhcmVkL3Byb3ZpZGVyL2dvI1Byb3ZpZGVyLlByaW50MTYyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcHJpbnQvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvLkxhbmd1YWdlXSwgc3RyaW5nXTEyMTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJpbnQvc3ludGF4L3NoYXJlZC9wcm92aWRlci9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28uTGFuZ3VhZ2VdNTE6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvI0xhbmd1YWdlNDpydXN0"></a>

## Rust / github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25). Receiver: github.com/Bitspark/bittype/print/syntax/shared/provider/go.Provider&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25) → [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L26)); receiver ; type arguments L.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NTg0Ok9EcHdjbTkyYVdSbGNqYzBPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3Y21sdWRDOXplVzUwWVhndmMyaGhjbVZrTDNCeWIzWnBaR1Z5TDJkdkkxQnliM1pwWkdWeUxsQnlhVzUwTVRZeU9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOWpiMjUwY21GamRITXZjSEpwYm5RdloyOHVVSEp2ZG1sa1pYSmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM1J5WldWekwyTnpkQzluYnk1RFUxUmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOXlkWE4wTDJkdkxreGhibWQxWVdkbFhTd2djM1J5YVc1blhURXlNVHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2Y0hKcGJuUXZjM2x1ZEdGNEwzTm9ZWEpsWkM5d2NtOTJhV1JsY2k5bmJ5NVFjbTkyYVdSbGNsdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwzSjFjM1F2WjI4dVRHRnVaM1ZoWjJWZE5URTZaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOXlkWE4wTDJkdkkweGhibWQxWVdkbE5EcHlkWE4wNjU6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3ByaW50L3N5bnRheC9zaGFyZWQvcHJvdmlkZXIvZ28jUHJpbnQ1OnNjb3BlMTojNDk6ZG9jcy9wYXJzaW5nLWZvcm1hdHRpbmcubWQjY2FycmllcnMtYW5kLW93bmVyc2hpcDU6ZmFsc2UwOg"></a>

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L106) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L115) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9zeW50YXgvc2hhcmVkL3Byb3ZpZGVyL2dvI1Byb3ZpZGVyLlByaW50MTYzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcHJpbnQvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nby5MYW5ndWFnZV0sIHN0cmluZ10xMjI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3ByaW50L3N5bnRheC9zaGFyZWQvcHJvdmlkZXIvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nby5MYW5ndWFnZV01MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvI0xhbmd1YWdlNTpob3Jvcw"></a>

## Horos / github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;, string&#93;.

Public operation: [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25). Receiver: github.com/Bitspark/bittype/print/syntax/shared/provider/go.Provider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25) → [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/print.go#L14) ([call](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L26)); receiver ; type arguments L.

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/print.go#L14)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NTkwOk9EcHdjbTkyYVdSbGNqYzBPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3Y21sdWRDOXplVzUwWVhndmMyaGhjbVZrTDNCeWIzWnBaR1Z5TDJkdkkxQnliM1pwWkdWeUxsQnlhVzUwTVRZek9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOWpiMjUwY21GamRITXZjSEpwYm5RdloyOHVVSEp2ZG1sa1pYSmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM1J5WldWekwyTnpkQzluYnk1RFUxUmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOW9iM0p2Y3k5bmJ5NU1ZVzVuZFdGblpWMHNJSE4wY21sdVoxMHhNakk2WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNCeWFXNTBMM041Ym5SaGVDOXphR0Z5WldRdmNISnZkbWxrWlhJdloyOHVVSEp2ZG1sa1pYSmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOW9iM0p2Y3k5bmJ5NU1ZVzVuZFdGblpWMDFNanBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMmh2Y205ekwyZHZJMHhoYm1kMVlXZGxOVHBvYjNKdmN3NjU6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3ByaW50L3N5bnRheC9zaGFyZWQvcHJvdmlkZXIvZ28jUHJpbnQ1OnNjb3BlMTojNDk6ZG9jcy9wYXJzaW5nLWZvcm1hdHRpbmcubWQjY2FycmllcnMtYW5kLW93bmVyc2hpcDU6ZmFsc2UwOg"></a>

Authored scope: [docs/parsing-formatting.md#carriers-and-ownership](https://github.com/Bitspark/bittype/blob/main/docs/parsing-formatting.md#carriers-and-ownership); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/shared/provider/go/provider.go#L25)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| horos-modules-and-generics | pagination-and-imported-library | canonical-identity-preserved | unspecified (authored profile absent) | [TestCanonicalPaginationAndImportedLibraryRoundTrip](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/provider/go/format_test.go#L59) (external) |
| independent-engine-stages | pagination-source | independently-composable | unspecified (authored profile absent) | [TestHorosSourceThroughIndependentEngineStages](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_text_test.go#L27) (external) |
| horos-text-vocabulary | declared-constructor-fixtures | canonical-round-trip-preserved | unspecified (authored profile absent) | [TestCompleteSourceVocabulary](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/provider/go/parse_test.go#L85) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9hbmFseXplL2ZpbGVzL3J1c3QvcHJvdmlkZXIvZ28jUHJvdmlkZXIuQW5hbHl6ZTExMTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL2FuYWx5c2lzL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nby5MYW5ndWFnZV02NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvYW5hbHl6ZS9maWxlcy9ydXN0L3Byb3ZpZGVyL2dvLlByb3ZpZGVyNTE6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvI0xhbmd1YWdlNDpydXN0"></a>

## Rust / github.com/Bitspark/bittype/analyze/files/rust/provider/go#Provider.Analyze

Contract: github.com/Bitspark/bittype/contracts/analysis/go.Provider&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/analyze/files/rust/provider/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/provider/go/analyze.go#L57). Receiver: github.com/Bitspark/bittype/analyze/files/rust/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/analyze/files/rust/provider/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/provider/go/analyze.go#L57) → [github.com/Bitspark/bittype/analyze/files/rust/provider/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/provider/go/analyze.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/provider/go/analyze.go#L58)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/analyze/files/rust/provider/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/provider/go/analyze.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NDQ0Ok9EcHdjbTkyYVdSbGNqYzFPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzloYm1Gc2VYcGxMMlpwYkdWekwzSjFjM1F2Y0hKdmRtbGtaWEl2WjI4alVISnZkbWxrWlhJdVFXNWhiSGw2WlRFeE1UcG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZZMjl1ZEhKaFkzUnpMMkZ1WVd4NWMybHpMMmR2TGxCeWIzWnBaR1Z5VzJkcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5emVXNTBZWGd2Y25WemRDOW5ieTVNWVc1bmRXRm5aVjAyTnpwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdllXNWhiSGw2WlM5bWFXeGxjeTl5ZFhOMEwzQnliM1pwWkdWeUwyZHZMbEJ5YjNacFpHVnlOVEU2WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNONWJuUmhlQzl5ZFhOMEwyZHZJMHhoYm1kMVlXZGxORHB5ZFhOMDY2OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9hbmFseXplL2ZpbGVzL3J1c3QvcHJvdmlkZXIvZ28jQW5hbHl6ZTU6c2NvcGUxOiM0Mzpkb2NzL25hdGl2ZS1zZW1hbnRpY3MubWQjZmlyc3QtcnVzdC1wcm9maWxlNTpmYWxzZTA6"></a>

Authored scope: [docs/native-semantics.md#first-rust-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-rust-profile); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/analyze/files/rust/provider/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/provider/go/analyze.go#L57)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | resolution-graph-preserved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| rust-semantic-subset | native-types-and-generics | native-identity-and-substitution-preserved | unspecified (authored profile absent) | [TestRustSemanticNativeTypes](https://github.com/Bitspark/bittype/blob/main/integration/go/rust_semantics_test.go#L202) (external) |
| rust-semantic-subset | bounds-and-profile-outcomes | qualified-reports-without-partial-model | unspecified (authored profile absent) | [TestRustSemanticBoundsAndOutcomes](https://github.com/Bitspark/bittype/blob/main/integration/go/rust_semantics_test.go#L264) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L2hvcm9zL3Byb3ZpZGVyL2dvI1Byb3ZpZGVyLkZvcm1hdDI1NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL2Zvcm1hdC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvYXN0L2dvLlZpZXdbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nby5MYW5ndWFnZV0sIGdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9jc3QvZ28uQ1NUW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdXTY4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L2hvcm9zL3Byb3ZpZGVyL2dvLlByb3ZpZGVyNTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nbyNMYW5ndWFnZTU6aG9yb3M"></a>

## Horos / github.com/Bitspark/bittype/format/syntax/horos/provider/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/format/syntax/horos/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/provider/go/format.go#L467). Receiver: github.com/Bitspark/bittype/format/syntax/horos/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/format/syntax/horos/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/provider/go/format.go#L467) → [github.com/Bitspark/bittype/format/syntax/horos/provider/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/provider/go/format.go#L42) ([call](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/provider/go/format.go#L468)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/horos/provider/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/provider/go/format.go#L42)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NjQzOk9EcHdjbTkyYVdSbGNqYzFPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzltYjNKdFlYUXZjM2x1ZEdGNEwyaHZjbTl6TDNCeWIzWnBaR1Z5TDJkdkkxQnliM1pwWkdWeUxrWnZjbTFoZERJMU56cG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZZMjl1ZEhKaFkzUnpMMlp2Y20xaGRDOW5ieTVRY205MmFXUmxjbHRuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2ZEhKbFpYTXZZWE4wTDJkdkxsWnBaWGRiWjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNONWJuUmhlQzlvYjNKdmN5OW5ieTVNWVc1bmRXRm5aVjBzSUdkcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5MGNtVmxjeTlqYzNRdloyOHVRMU5VVzJkcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5emVXNTBZWGd2YUc5eWIzTXZaMjh1VEdGdVozVmhaMlZkWFRZNE9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOW1iM0p0WVhRdmMzbHVkR0Y0TDJodmNtOXpMM0J5YjNacFpHVnlMMmR2TGxCeWIzWnBaR1Z5TlRJNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzTjViblJoZUM5b2IzSnZjeTluYnlOTVlXNW5kV0ZuWlRVNmFHOXliM002NjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvZm9ybWF0L3N5bnRheC9ob3Jvcy9wcm92aWRlci9nbyNGb3JtYXQ1OnNjb3BlMTojNDA6ZG9jcy9ob3Jvcy1zeW50YXgubWQjcGlwZWxpbmUtYW5kLWxpbWl0czU6ZmFsc2UwOg"></a>

Authored scope: [docs/horos-syntax.md#pipeline-and-limits](https://github.com/Bitspark/bittype/blob/main/docs/horos-syntax.md#pipeline-and-limits); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| ast-admission | unspellable-and-nil-ast | rejection-without-output | unspecified (authored profile absent) | [TestFormatterRejectsUnspellableASTWithoutPartialOutput](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/provider/go/format_test.go#L167) (external) |

### Declared target [github.com/Bitspark/bittype/format/syntax/horos/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/provider/go/format.go#L467)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| horos-modules-and-generics | pagination-and-imported-library | canonical-identity-preserved | unspecified (authored profile absent) | [TestCanonicalPaginationAndImportedLibraryRoundTrip](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/provider/go/format_test.go#L59) (external) |
| independent-engine-stages | pagination-source | independently-composable | unspecified (authored profile absent) | [TestHorosSourceThroughIndependentEngineStages](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_text_test.go#L27) (external) |
| horos-text-vocabulary | declared-constructor-fixtures | canonical-round-trip-preserved | unspecified (authored profile absent) | [TestCompleteSourceVocabulary](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/provider/go/parse_test.go#L85) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9yZXNvbHZlL2ZpbGVzL3J1c3QvcHJvdmlkZXIvZ28jUHJvdmlkZXIuUmVzb2x2ZTExMzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3Jlc29sdXRpb24vZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvLkxhbmd1YWdlXTY3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9yZXNvbHZlL2ZpbGVzL3J1c3QvcHJvdmlkZXIvZ28uUHJvdmlkZXI1MTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28jTGFuZ3VhZ2U0OnJ1c3Q"></a>

## Rust / github.com/Bitspark/bittype/resolve/files/rust/provider/go#Provider.Resolve

Contract: github.com/Bitspark/bittype/contracts/resolution/go.Provider&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/resolve/files/rust/provider/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/rust/provider/go/resolver.go#L44). Receiver: github.com/Bitspark/bittype/resolve/files/rust/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

### Declared target [github.com/Bitspark/bittype/resolve/files/rust/provider/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/rust/provider/go/resolver.go#L44)

<a id="binding-NDQ3Ok9EcHdjbTkyYVdSbGNqYzFPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl5WlhOdmJIWmxMMlpwYkdWekwzSjFjM1F2Y0hKdmRtbGtaWEl2WjI4alVISnZkbWxrWlhJdVVtVnpiMngyWlRFeE16cG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZZMjl1ZEhKaFkzUnpMM0psYzI5c2RYUnBiMjR2WjI4dVVISnZkbWxrWlhKYloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzTjViblJoZUM5eWRYTjBMMmR2TGt4aGJtZDFZV2RsWFRZM09tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOXlaWE52YkhabEwyWnBiR1Z6TDNKMWMzUXZjSEp2ZG1sa1pYSXZaMjh1VUhKdmRtbGtaWEkxTVRwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDNKMWMzUXZaMjhqVEdGdVozVmhaMlUwT25KMWMzUTc1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9yZXNvbHZlL2ZpbGVzL3J1c3QvcHJvdmlkZXIvZ28jUHJvdmlkZXIuUmVzb2x2ZTU6c2NvcGUxOiM0Mzpkb2NzL25hdGl2ZS1zZW1hbnRpY3MubWQjZmlyc3QtcnVzdC1wcm9maWxlNTpmYWxzZTA6"></a>

Authored scope: [docs/native-semantics.md#first-rust-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-rust-profile); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | captured-graph-resolved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| rust-semantic-subset | native-types-and-generics | native-identity-and-substitution-preserved | unspecified (authored profile absent) | [TestRustSemanticNativeTypes](https://github.com/Bitspark/bittype/blob/main/integration/go/rust_semantics_test.go#L202) (external) |
| rust-semantic-subset | bounds-and-profile-outcomes | qualified-reports-without-partial-model | unspecified (authored profile absent) | [TestRustSemanticBoundsAndOutcomes](https://github.com/Bitspark/bittype/blob/main/integration/go/rust_semantics_test.go#L264) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc2OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9hbmFseXplL2ZpbGVzL2hvcm9zL3Byb3ZpZGVyL2dvI1Byb3ZpZGVyLkFuYWx5emUxMTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9hbmFseXNpcy9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXTY4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9hbmFseXplL2ZpbGVzL2hvcm9zL3Byb3ZpZGVyL2dvLlByb3ZpZGVyNTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nbyNMYW5ndWFnZTU6aG9yb3M"></a>

## Horos / github.com/Bitspark/bittype/analyze/files/horos/provider/go#Provider.Analyze

Contract: github.com/Bitspark/bittype/contracts/analysis/go.Provider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/analyze/files/horos/provider/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/horos/provider/go/analyze.go#L75). Receiver: github.com/Bitspark/bittype/analyze/files/horos/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/analyze/files/horos/provider/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/horos/provider/go/analyze.go#L75) → [github.com/Bitspark/bittype/analyze/files/horos/provider/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/horos/provider/go/analyze.go#L37) ([call](https://github.com/Bitspark/bittype/blob/main/analyze/files/horos/provider/go/analyze.go#L76)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/analyze/files/horos/provider/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/horos/provider/go/analyze.go#L37)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NDUxOk9EcHdjbTkyYVdSbGNqYzJPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzloYm1Gc2VYcGxMMlpwYkdWekwyaHZjbTl6TDNCeWIzWnBaR1Z5TDJkdkkxQnliM1pwWkdWeUxrRnVZV3g1ZW1VeE1USTZaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMMk52Ym5SeVlXTjBjeTloYm1Gc2VYTnBjeTluYnk1UWNtOTJhV1JsY2x0bmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDJodmNtOXpMMmR2TGt4aGJtZDFZV2RsWFRZNE9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOWhibUZzZVhwbEwyWnBiR1Z6TDJodmNtOXpMM0J5YjNacFpHVnlMMmR2TGxCeWIzWnBaR1Z5TlRJNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzTjViblJoZUM5b2IzSnZjeTluYnlOTVlXNW5kV0ZuWlRVNmFHOXliM002NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvYW5hbHl6ZS9maWxlcy9ob3Jvcy9wcm92aWRlci9nbyNBbmFseXplNTpzY29wZTE6IzMwOmRvY3MvYW5hbHlzaXMubWQjaG9yb3MtcHJvZmlsZTU6ZmFsc2UwOg"></a>

Authored scope: [docs/analysis.md#horos-profile](https://github.com/Bitspark/bittype/blob/main/docs/analysis.md#horos-profile); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/analyze/files/horos/provider/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/horos/provider/go/analyze.go#L75)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | analysis-evidence-produced | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| independent-engine-stages | pagination-source | independently-composable | unspecified (authored profile absent) | [TestHorosSourceThroughIndependentEngineStages](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_text_test.go#L27) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc2OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L2dvbGFuZy9wcm92aWRlci9nbyNQcm92aWRlci5Gb3JtYXQyNTk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9mb3JtYXQvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2FzdC9nby5WaWV3W2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvZ29sYW5nL2dvLkxhbmd1YWdlXSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28uTGFuZ3VhZ2VdXTY5OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L2dvbGFuZy9wcm92aWRlci9nby5Qcm92aWRlcjUzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvZ29sYW5nL2dvI0xhbmd1YWdlMjpnbw"></a>

## Go / github.com/Bitspark/bittype/format/syntax/golang/provider/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/format/syntax/golang/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/golang/provider/go/format.go#L177). Receiver: github.com/Bitspark/bittype/format/syntax/golang/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/format/syntax/golang/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/golang/provider/go/format.go#L177) → [github.com/Bitspark/bittype/format/syntax/golang/provider/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/golang/provider/go/format.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/format/syntax/golang/provider/go/format.go#L178)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/golang/provider/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/golang/provider/go/format.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NjQ2Ok9EcHdjbTkyYVdSbGNqYzJPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzltYjNKdFlYUXZjM2x1ZEdGNEwyZHZiR0Z1Wnk5d2NtOTJhV1JsY2k5bmJ5TlFjbTkyYVdSbGNpNUdiM0p0WVhReU5UazZaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMMk52Ym5SeVlXTjBjeTltYjNKdFlYUXZaMjh1VUhKdmRtbGtaWEpiWjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNSeVpXVnpMMkZ6ZEM5bmJ5NVdhV1YzVzJkcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5emVXNTBZWGd2WjI5c1lXNW5MMmR2TGt4aGJtZDFZV2RsWFN3Z1oybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzUnlaV1Z6TDJOemRDOW5ieTVEVTFSYloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzTjViblJoZUM5bmIyeGhibWN2WjI4dVRHRnVaM1ZoWjJWZFhUWTVPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzltYjNKdFlYUXZjM2x1ZEdGNEwyZHZiR0Z1Wnk5d2NtOTJhV1JsY2k5bmJ5NVFjbTkyYVdSbGNqVXpPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl6ZVc1MFlYZ3ZaMjlzWVc1bkwyZHZJMHhoYm1kMVlXZGxNanBuYnc2NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvZm9ybWF0L3N5bnRheC9nb2xhbmcvcHJvdmlkZXIvZ28jRm9ybWF0NTpzY29wZTE6IzE3OmRvY3Mvc3ludGF4Lm1kI2dvNTpmYWxzZTA6"></a>

Authored scope: [docs/syntax.md#go](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#go); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| go-concrete-pipeline | pagination-and-interface | stable-formatting | unspecified (authored profile absent) | [TestGoConcretePipeline](https://github.com/Bitspark/bittype/blob/main/integration/go/text_test.go#L15) (external) |

### Declared target [github.com/Bitspark/bittype/format/syntax/golang/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/golang/provider/go/format.go#L177)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L106) (external) |
| ast-admission | malformed-carriers | rejection-without-output | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L193) (external) |
| native-compiler-qualification | formatted-fixture-contract | compiler-fixture-contract | unspecified (authored profile absent) | [TestNativeCompilerFixtures](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L265) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L115) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc2OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9yZXNvbHZlL2ZpbGVzL2hvcm9zL3Byb3ZpZGVyL2dvI1Byb3ZpZGVyLlJlc29sdmUxMTQ6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9yZXNvbHV0aW9uL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdNjg6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3Jlc29sdmUvZmlsZXMvaG9yb3MvcHJvdmlkZXIvZ28uUHJvdmlkZXI1MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvI0xhbmd1YWdlNTpob3Jvcw"></a>

## Horos / github.com/Bitspark/bittype/resolve/files/horos/provider/go#Provider.Resolve

Contract: github.com/Bitspark/bittype/contracts/resolution/go.Provider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/resolve/files/horos/provider/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/horos/provider/go/resolver.go#L57). Receiver: github.com/Bitspark/bittype/resolve/files/horos/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/resolve/files/horos/provider/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/horos/provider/go/resolver.go#L57) → [github.com/Bitspark/bittype/resolve/files/horos/provider/go#ResolveReferences](https://github.com/Bitspark/bittype/blob/main/resolve/files/horos/provider/go/resolver.go#L50) ([call](https://github.com/Bitspark/bittype/blob/main/resolve/files/horos/provider/go/resolver.go#L58)); receiver ; type arguments .

Traversal boundary: callee is dynamic or outside implementation source.

### Declared target [github.com/Bitspark/bittype/resolve/files/horos/provider/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/horos/provider/go/resolver.go#L57)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | captured-graph-resolved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| independent-engine-stages | pagination-source | independently-composable | unspecified (authored profile absent) | [TestHorosSourceThroughIndependentEngineStages](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_text_test.go#L27) (external) |

### Declared target [github.com/Bitspark/bittype/resolve/files/horos/provider/go#ResolveReferences](https://github.com/Bitspark/bittype/blob/main/resolve/files/horos/provider/go/resolver.go#L50)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NDU0Ok9EcHdjbTkyYVdSbGNqYzJPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl5WlhOdmJIWmxMMlpwYkdWekwyaHZjbTl6TDNCeWIzWnBaR1Z5TDJkdkkxQnliM1pwWkdWeUxsSmxjMjlzZG1VeE1UUTZaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMMk52Ym5SeVlXTjBjeTl5WlhOdmJIVjBhVzl1TDJkdkxsQnliM1pwWkdWeVcyZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOXplVzUwWVhndmFHOXliM012WjI4dVRHRnVaM1ZoWjJWZE5qZzZaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM0psYzI5c2RtVXZabWxzWlhNdmFHOXliM012Y0hKdmRtbGtaWEl2WjI4dVVISnZkbWxrWlhJMU1qcG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwyaHZjbTl6TDJkdkkweGhibWQxWVdkbE5UcG9iM0p2Y3c3NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcmVzb2x2ZS9maWxlcy9ob3Jvcy9wcm92aWRlci9nbyNSZXNvbHZlUmVmZXJlbmNlczU6c2NvcGUxOiMzNjpkb2NzL3Jlc29sdXRpb24ubWQjaG9yb3MtY29tcG9zaXRpb241OmZhbHNlMDo"></a>

Authored scope: [docs/resolution.md#horos-composition](https://github.com/Bitspark/bittype/blob/main/docs/resolution.md#horos-composition); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9hbmFseXplL2ZpbGVzL2dvbGFuZy9wcm92aWRlci9nbyNQcm92aWRlci5BbmFseXplMTEzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvYW5hbHlzaXMvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28uTGFuZ3VhZ2VdNjk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2FuYWx5emUvZmlsZXMvZ29sYW5nL3Byb3ZpZGVyL2dvLlByb3ZpZGVyNTM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9nb2xhbmcvZ28jTGFuZ3VhZ2UyOmdv"></a>

## Go / github.com/Bitspark/bittype/analyze/files/golang/provider/go#Provider.Analyze

Contract: github.com/Bitspark/bittype/contracts/analysis/go.Provider&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/analyze/files/golang/provider/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/provider/go/analyze.go#L57). Receiver: github.com/Bitspark/bittype/analyze/files/golang/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/analyze/files/golang/provider/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/provider/go/analyze.go#L57) → [github.com/Bitspark/bittype/analyze/files/golang/provider/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/provider/go/analyze.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/provider/go/analyze.go#L58)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/analyze/files/golang/provider/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/provider/go/analyze.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NDUyOk9EcHdjbTkyYVdSbGNqYzNPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzloYm1Gc2VYcGxMMlpwYkdWekwyZHZiR0Z1Wnk5d2NtOTJhV1JsY2k5bmJ5TlFjbTkyYVdSbGNpNUJibUZzZVhwbE1URXpPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYjI1MGNtRmpkSE12WVc1aGJIbHphWE12WjI4dVVISnZkbWxrWlhKYloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzTjViblJoZUM5bmIyeGhibWN2WjI4dVRHRnVaM1ZoWjJWZE5qazZaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMMkZ1WVd4NWVtVXZabWxzWlhNdloyOXNZVzVuTDNCeWIzWnBaR1Z5TDJkdkxsQnliM1pwWkdWeU5UTTZaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOW5iMnhoYm1jdloyOGpUR0Z1WjNWaFoyVXlPbWR2Njg6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2FuYWx5emUvZmlsZXMvZ29sYW5nL3Byb3ZpZGVyL2dvI0FuYWx5emU1OnNjb3BlMTojNDE6ZG9jcy9uYXRpdmUtc2VtYW50aWNzLm1kI2ZpcnN0LWdvLXByb2ZpbGU1OmZhbHNlMDo"></a>

Authored scope: [docs/native-semantics.md#first-go-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-go-profile); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/analyze/files/golang/provider/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/provider/go/analyze.go#L57)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | resolution-graph-preserved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| go-semantic-subset | rejected-and-unsupported-fixtures | no-partial-model | unspecified (authored profile absent) | [TestGoSemanticReportsAndCoverage](https://github.com/Bitspark/bittype/blob/main/integration/go/native_semantics_test.go#L188) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9yZXNvbHZlL2ZpbGVzL2dvbGFuZy9wcm92aWRlci9nbyNQcm92aWRlci5SZXNvbHZlMTE1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcmVzb2x1dGlvbi9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nby5MYW5ndWFnZV02OTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcmVzb2x2ZS9maWxlcy9nb2xhbmcvcHJvdmlkZXIvZ28uUHJvdmlkZXI1MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nbyNMYW5ndWFnZTI6Z28"></a>

## Go / github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve

Contract: github.com/Bitspark/bittype/contracts/resolution/go.Provider&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/provider/go/resolver.go#L82). Receiver: github.com/Bitspark/bittype/resolve/files/golang/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/provider/go/resolver.go#L82) → [github.com/Bitspark/bittype/resolve/files/golang/provider/go#Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/provider/go/resolver.go#L25) ([call](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/provider/go/resolver.go#L83)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/provider/go/resolver.go#L82)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | captured-graph-resolved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| go-semantic-subset | rejected-and-unsupported-fixtures | no-partial-model | unspecified (authored profile absent) | [TestGoSemanticReportsAndCoverage](https://github.com/Bitspark/bittype/blob/main/integration/go/native_semantics_test.go#L188) (external) |

### Declared target [github.com/Bitspark/bittype/resolve/files/golang/provider/go#Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/provider/go/resolver.go#L25)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NDU1Ok9EcHdjbTkyYVdSbGNqYzNPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl5WlhOdmJIWmxMMlpwYkdWekwyZHZiR0Z1Wnk5d2NtOTJhV1JsY2k5bmJ5TlFjbTkyYVdSbGNpNVNaWE52YkhabE1URTFPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYjI1MGNtRmpkSE12Y21WemIyeDFkR2x2Ymk5bmJ5NVFjbTkyYVdSbGNsdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwyZHZiR0Z1Wnk5bmJ5NU1ZVzVuZFdGblpWMDJPVHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2Y21WemIyeDJaUzltYVd4bGN5OW5iMnhoYm1jdmNISnZkbWxrWlhJdloyOHVVSEp2ZG1sa1pYSTFNenBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMmR2YkdGdVp5OW5ieU5NWVc1bmRXRm5aVEk2WjI4Njg6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3Jlc29sdmUvZmlsZXMvZ29sYW5nL3Byb3ZpZGVyL2dvI1Jlc29sdmU1OnNjb3BlMTojNDE6ZG9jcy9uYXRpdmUtc2VtYW50aWNzLm1kI2ZpcnN0LWdvLXByb2ZpbGU1OmZhbHNlMDo"></a>

Authored scope: [docs/native-semantics.md#first-go-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-go-profile); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjc4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvdHlwZXNjcmlwdC9wcm92aWRlci9nbyNQcm92aWRlci5DaGVjazExNDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL2NoZWNrL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nby5MYW5ndWFnZV03MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY2hlY2svc3ludGF4L3R5cGVzY3JpcHQvcHJvdmlkZXIvZ28uUHJvdmlkZXI1NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28jTGFuZ3VhZ2UxMDp0eXBlc2NyaXB0"></a>

## TypeScript / github.com/Bitspark/bittype/check/syntax/typescript/provider/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/check/syntax/typescript/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/provider/go/check.go#L177). Receiver: github.com/Bitspark/bittype/check/syntax/typescript/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/check/syntax/typescript/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/provider/go/check.go#L177) → [github.com/Bitspark/bittype/check/syntax/typescript/provider/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/provider/go/check.go#L52) ([call](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/provider/go/check.go#L177)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/check/syntax/typescript/provider/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/provider/go/check.go#L52)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NDc2Ok9EcHdjbTkyYVdSbGNqYzRPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYUdWamF5OXplVzUwWVhndmRIbHdaWE5qY21sd2RDOXdjbTkyYVdSbGNpOW5ieU5RY205MmFXUmxjaTVEYUdWamF6RXhORHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2WTI5dWRISmhZM1J6TDJOb1pXTnJMMmR2TGxCeWIzWnBaR1Z5VzJkcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5emVXNTBZWGd2ZEhsd1pYTmpjbWx3ZEM5bmJ5NU1ZVzVuZFdGblpWMDNNanBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2WTJobFkyc3ZjM2x1ZEdGNEwzUjVjR1Z6WTNKcGNIUXZjSEp2ZG1sa1pYSXZaMjh1VUhKdmRtbGtaWEkxTnpwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDNSNWNHVnpZM0pwY0hRdloyOGpUR0Z1WjNWaFoyVXhNRHAwZVhCbGMyTnlhWEIwNjk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NoZWNrL3N5bnRheC90eXBlc2NyaXB0L3Byb3ZpZGVyL2dvI0NoZWNrNzpncmFtbWFyNjk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NoZWNrL3N5bnRheC90eXBlc2NyaXB0L3Byb3ZpZGVyL2dvI3J1bGVzMDo1OmZhbHNlMDo"></a>

Definition: [github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/provider/go/check.go#L18); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### source-file (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key source-file.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| import | 0 | unbounded | import |
| declaration | 0 | unbounded | type-alias, interface |

No feature-key case association at this exact authored profile.

#### import (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key import.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| binding | 1 | unbounded | import-binding, namespace-import |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| module | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### import-binding (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key import-binding.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### namespace-import (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key namespace-import.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-alias (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-alias.

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

#### interface (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key interface.

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

#### type-parameter (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-parameter.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| constraint | 0 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| default | 0 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-reference (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-reference.

Roles: 0 declared. Attributes: 1 declared. References: target.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### qualified-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key qualified-type.

Roles: 1 declared. Attributes: 1 declared. References: target.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| qualifier | 1 | 1 | type-reference, qualified-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-application (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-application.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| constructor | 1 | 1 | type-reference, qualified-type |
| argument | 1 | unbounded | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### keyword-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key keyword-type.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | any, unknown, never, void, undefined, boolean, number, bigint, string, symbol, object, this |

No feature-key case association at this exact authored profile.

#### literal-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key literal-type.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| literal-kind | true | string, number, bigint, boolean, null |
| text | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### array-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key array-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### tuple-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 0 | unbounded | tuple-element |

No feature-key case association at this exact authored profile.

#### tuple-element (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-element.

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

#### object-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key object-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| member | 0 | unbounded | property, method, call-signature, index-signature |

No feature-key case association at this exact authored profile.

#### union-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key union-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 2 | unbounded | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### intersection-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key intersection-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 2 | unbounded | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### conditional-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key conditional-type.

Roles: 4 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| check | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| constraint | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| true | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| false | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### infer-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key infer-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 1 | 1 | type-parameter |

No feature-key case association at this exact authored profile.

#### mapped-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key mapped-type.

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

#### function-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key function-type.

Roles: 3 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type-parameter | 0 | unbounded | type-parameter |
| parameter | 0 | unbounded | parameter |
| result | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

#### parameter (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

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

#### property (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key property.

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

#### method (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key method.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| signature | 1 | 1 | function-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| optional | false | true, false |

No feature-key case association at this exact authored profile.

#### call-signature (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key call-signature.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| signature | 1 | 1 | function-type |

No feature-key case association at this exact authored profile.

#### index-signature (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key index-signature.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 1 | 1 | parameter |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| readonly | false | true, false |

No feature-key case association at this exact authored profile.

#### type-operator (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-operator.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| operator | true | keyof, readonly, unique |

No feature-key case association at this exact authored profile.

#### indexed-access-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/typescript/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key indexed-access-type.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| object | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |
| index | 1 | 1 | type-reference, qualified-type, type-application, keyword-type, literal-type, array-type, tuple-type, object-type, union-type, intersection-type, conditional-type, infer-type, mapped-type, function-type, type-operator, indexed-access-type |

No feature-key case association at this exact authored profile.

<a id="binding-NDc2Ok9EcHdjbTkyYVdSbGNqYzRPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYUdWamF5OXplVzUwWVhndmRIbHdaWE5qY21sd2RDOXdjbTkyYVdSbGNpOW5ieU5RY205MmFXUmxjaTVEYUdWamF6RXhORHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2WTI5dWRISmhZM1J6TDJOb1pXTnJMMmR2TGxCeWIzWnBaR1Z5VzJkcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5emVXNTBZWGd2ZEhsd1pYTmpjbWx3ZEM5bmJ5NU1ZVzVuZFdGblpWMDNNanBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2WTJobFkyc3ZjM2x1ZEdGNEwzUjVjR1Z6WTNKcGNIUXZjSEp2ZG1sa1pYSXZaMjh1VUhKdmRtbGtaWEkxTnpwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDNSNWNHVnpZM0pwY0hRdloyOGpUR0Z1WjNWaFoyVXhNRHAwZVhCbGMyTnlhWEIwNjk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NoZWNrL3N5bnRheC90eXBlc2NyaXB0L3Byb3ZpZGVyL2dvI0NoZWNrNTpzY29wZTE6IzI1OmRvY3Mvc3ludGF4Lm1kI3R5cGVzY3JpcHQ1OmZhbHNlMDo"></a>

Authored scope: [docs/syntax.md#typescript](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#typescript); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| conditional-mapped-syntax | independent-ast | accept | unspecified (authored profile absent) | [TestConditionalMappedAdmission](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/provider/go/types_test.go#L33) (external) |
| conditional-mapped-syntax | invalid-shapes-and-contexts | located-rejection | unspecified (authored profile absent) | [TestInvalidConditionalMappedShapes](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/provider/go/types_test.go#L41) (external) |

### Declared target [github.com/Bitspark/bittype/check/syntax/typescript/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/provider/go/check.go#L177)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| ast-admission | malformed-carriers | reject | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L193) (external) |

<a id="cap-ODpwcm92aWRlcjc4OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvdHlwZXNjcmlwdC9wcm92aWRlci9nbyNQcm92aWRlci5QYXJzZTIxNDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3BhcnNlL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcGFyc2UvZ28uSW5wdXQsIGdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9jc3QvZ28uQ1NUW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nby5MYW5ndWFnZV1dNzI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3BhcnNlL3N5bnRheC90eXBlc2NyaXB0L3Byb3ZpZGVyL2dvLlByb3ZpZGVyNTc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvI0xhbmd1YWdlMTA6dHlwZXNjcmlwdA"></a>

## TypeScript / github.com/Bitspark/bittype/parse/syntax/typescript/provider/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/parse/syntax/typescript/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/provider/go/parse.go#L568). Receiver: github.com/Bitspark/bittype/parse/syntax/typescript/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/parse/syntax/typescript/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/provider/go/parse.go#L568) → [github.com/Bitspark/bittype/parse/syntax/typescript/provider/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/provider/go/parse.go#L42) ([call](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/provider/go/parse.go#L569)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/typescript/provider/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/provider/go/parse.go#L42)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NjEwOk9EcHdjbTkyYVdSbGNqYzRPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3WVhKelpTOXplVzUwWVhndmRIbHdaWE5qY21sd2RDOXdjbTkyYVdSbGNpOW5ieU5RY205MmFXUmxjaTVRWVhKelpUSXhORHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2WTI5dWRISmhZM1J6TDNCaGNuTmxMMmR2TGxCeWIzWnBaR1Z5VzJkcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdmNHRnljMlV2WjI4dVNXNXdkWFFzSUdkcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5MGNtVmxjeTlqYzNRdloyOHVRMU5VVzJkcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5emVXNTBZWGd2ZEhsd1pYTmpjbWx3ZEM5bmJ5NU1ZVzVuZFdGblpWMWROekk2WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNCaGNuTmxMM041Ym5SaGVDOTBlWEJsYzJOeWFYQjBMM0J5YjNacFpHVnlMMmR2TGxCeWIzWnBaR1Z5TlRjNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzTjViblJoZUM5MGVYQmxjMk55YVhCMEwyZHZJMHhoYm1kMVlXZGxNVEE2ZEhsd1pYTmpjbWx3ZEE2OTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcGFyc2Uvc3ludGF4L3R5cGVzY3JpcHQvcHJvdmlkZXIvZ28jUGFyc2U1OnNjb3BlMTojMjU6ZG9jcy9zeW50YXgubWQjdHlwZXNjcmlwdDU6ZmFsc2UwOg"></a>

Authored scope: [docs/syntax.md#typescript](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#typescript); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| conditional-mapped-syntax | structured-binders-and-modifiers | structure-preserved | unspecified (authored profile absent) | [TestConditionalMappedStructure](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/provider/go/types_test.go#L17) (external) |
| conditional-mapped-syntax | unresolved-syntax | accept | unspecified (authored profile absent) | [TestConditionalMappedWithoutResolution](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/provider/go/types_test.go#L48) (external) |
| conditional-mapped-syntax | malformed-and-bounded-input | located-error-without-output | unspecified (authored profile absent) | [TestMalformedConditionalMappedTypes](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/provider/go/types_test.go#L71) (external) |

### Declared target [github.com/Bitspark/bittype/parse/syntax/typescript/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/provider/go/parse.go#L568)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L115) (external) |
| source-boundary | unsupported-and-malformed-fixtures | located-error-without-output | unspecified (authored profile absent) | [TestUnsupportedAndMalformedSources](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L137) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjcxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvY3BwL3Byb3ZpZGVyL2dvI1Byb3ZpZGVyLkNoZWNrMTA3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvY2hlY2svZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28uTGFuZ3VhZ2VdNjU6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NoZWNrL3N5bnRheC9jcHAvcHJvdmlkZXIvZ28uUHJvdmlkZXI1MDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nbyNMYW5ndWFnZTM6Y3Bw"></a>

## C++ / github.com/Bitspark/bittype/check/syntax/cpp/provider/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/check/syntax/cpp/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/provider/go/check.go#L129). Receiver: github.com/Bitspark/bittype/check/syntax/cpp/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/check/syntax/cpp/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/provider/go/check.go#L129) → [github.com/Bitspark/bittype/check/syntax/cpp/provider/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/provider/go/check.go#L57) ([call](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/provider/go/check.go#L129)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/check/syntax/cpp/provider/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/provider/go/check.go#L57)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NDI4Ok9EcHdjbTkyYVdSbGNqY3hPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYUdWamF5OXplVzUwWVhndlkzQndMM0J5YjNacFpHVnlMMmR2STFCeWIzWnBaR1Z5TGtOb1pXTnJNVEEzT21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdlkyaGxZMnN2WjI4dVVISnZkbWxrWlhKYloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzTjViblJoZUM5amNIQXZaMjh1VEdGdVozVmhaMlZkTmpVNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwyTm9aV05yTDNONWJuUmhlQzlqY0hBdmNISnZkbWxrWlhJdloyOHVVSEp2ZG1sa1pYSTFNRHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMk53Y0M5bmJ5Tk1ZVzVuZFdGblpUTTZZM0J3NjI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NoZWNrL3N5bnRheC9jcHAvcHJvdmlkZXIvZ28jQ2hlY2s3OmdyYW1tYXI2MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY2hlY2svc3ludGF4L2NwcC9wcm92aWRlci9nbyNydWxlczA6NTpmYWxzZTA6"></a>

Definition: [github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/provider/go/check.go#L29); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### translation-unit (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key translation-unit.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| include | 0 | unbounded | include |
| declaration | 0 | unbounded | namespace, record, type-alias, enum |

No feature-key case association at this exact authored profile.

#### include (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key include.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| path | true | opaque: Predicate present; behavior not interpreted |
| system | false | true, false |

No feature-key case association at this exact authored profile.

#### namespace (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key namespace.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| declaration | 0 | unbounded | namespace, record, type-alias, enum |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### record (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key record.

Roles: 3 declared. Attributes: unknown: Attribute collection could not be decoded. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| template-parameter | 0 | unbounded | type-parameter |
| base | 0 | unbounded | base |
| member | 0 | unbounded | field, method, destructor, type-alias, record, enum |

No feature-key case association at this exact authored profile.

#### type-alias (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-alias.

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

#### enum (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key enum.

Roles: 2 declared. Attributes: unknown: Attribute collection could not be decoded. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| underlying | 0 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |
| enumerator | 0 | unbounded | enumerator |

No feature-key case association at this exact authored profile.

#### enumerator (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key enumerator.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| value | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### base (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key base.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| access | false | public, protected, private |
| virtual | false | true, false |

No feature-key case association at this exact authored profile.

#### field (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key field.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| access | false | public, protected, private |

No feature-key case association at this exact authored profile.

#### method (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key method.

Roles: 2 declared. Attributes: unknown: Attribute collection could not be decoded. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| parameter | 0 | unbounded | parameter |
| result | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

#### destructor (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key destructor.

Roles: 0 declared. Attributes: 4 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| access | false | public, protected, private |
| virtual | false | true, false |
| pure | false | true, false |
| noexcept | false | true, false |

No feature-key case association at this exact authored profile.

#### parameter (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-parameter (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-parameter.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| default | 0 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| parameter-key | false | typename, class |

No feature-key case association at this exact authored profile.

#### named-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key named-type.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| segment | 1 | unbounded | name-segment |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| global | false | true, false |
| typename | false | true, false |

No feature-key case association at this exact authored profile.

#### name-segment (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key name-segment.

Roles: 1 declared. Attributes: 2 declared. References: target.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| argument | 0 | unbounded | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| template | false | true, false |

No feature-key case association at this exact authored profile.

#### builtin-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key builtin-type.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### pointer-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key pointer-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

#### lvalue-reference-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key lvalue-reference-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

#### rvalue-reference-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key rvalue-reference-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

#### cv-qualified-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key cv-qualified-type.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| const | false | true, false |
| volatile | false | true, false |

No feature-key case association at this exact authored profile.

#### array-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key array-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| length | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### function-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key function-type.

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

#### member-pointer-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/cpp/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key member-pointer-type.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| owner | 1 | 1 | named-type |
| element | 1 | 1 | named-type, builtin-type, pointer-type, lvalue-reference-type, rvalue-reference-type, cv-qualified-type, array-type, function-type, member-pointer-type |

No feature-key case association at this exact authored profile.

<a id="binding-NDI4Ok9EcHdjbTkyYVdSbGNqY3hPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYUdWamF5OXplVzUwWVhndlkzQndMM0J5YjNacFpHVnlMMmR2STFCeWIzWnBaR1Z5TGtOb1pXTnJNVEEzT21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdlkyaGxZMnN2WjI4dVVISnZkbWxrWlhKYloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzTjViblJoZUM5amNIQXZaMjh1VEdGdVozVmhaMlZkTmpVNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwyTm9aV05yTDNONWJuUmhlQzlqY0hBdmNISnZkbWxrWlhJdloyOHVVSEp2ZG1sa1pYSTFNRHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMk53Y0M5bmJ5Tk1ZVzVuZFdGblpUTTZZM0J3NjI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NoZWNrL3N5bnRheC9jcHAvcHJvdmlkZXIvZ28jQ2hlY2s1OnNjb3BlMTojMTY6ZG9jcy9zeW50YXgubWQjYzU6ZmFsc2UwOg"></a>

Authored scope: [docs/syntax.md#c](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#c); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/check/syntax/cpp/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/provider/go/check.go#L129)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| ast-admission | malformed-carriers | reject | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L193) (external) |

<a id="cap-ODpwcm92aWRlcjcxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvY3BwL3Byb3ZpZGVyL2dvI1Byb3ZpZGVyLlBhcnNlMjA3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcGFyc2UvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wYXJzZS9nby5JbnB1dCwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28uTGFuZ3VhZ2VdXTY1OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvY3BwL3Byb3ZpZGVyL2dvLlByb3ZpZGVyNTA6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28jTGFuZ3VhZ2UzOmNwcA"></a>

## C++ / github.com/Bitspark/bittype/parse/syntax/cpp/provider/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/parse/syntax/cpp/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/cpp/provider/go/parse.go#L441). Receiver: github.com/Bitspark/bittype/parse/syntax/cpp/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/parse/syntax/cpp/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/cpp/provider/go/parse.go#L441) → [github.com/Bitspark/bittype/parse/syntax/cpp/provider/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/cpp/provider/go/parse.go#L25) ([call](https://github.com/Bitspark/bittype/blob/main/parse/syntax/cpp/provider/go/parse.go#L442)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/cpp/provider/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/cpp/provider/go/parse.go#L25)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NTYyOk9EcHdjbTkyYVdSbGNqY3hPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3WVhKelpTOXplVzUwWVhndlkzQndMM0J5YjNacFpHVnlMMmR2STFCeWIzWnBaR1Z5TGxCaGNuTmxNakEzT21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdmNHRnljMlV2WjI4dVVISnZkbWxrWlhKYloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwyTnZiblJ5WVdOMGN5OXdZWEp6WlM5bmJ5NUpibkIxZEN3Z1oybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzUnlaV1Z6TDJOemRDOW5ieTVEVTFSYloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzTjViblJoZUM5amNIQXZaMjh1VEdGdVozVmhaMlZkWFRZMU9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOXdZWEp6WlM5emVXNTBZWGd2WTNCd0wzQnliM1pwWkdWeUwyZHZMbEJ5YjNacFpHVnlOVEE2WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNONWJuUmhlQzlqY0hBdloyOGpUR0Z1WjNWaFoyVXpPbU53Y0E2MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcGFyc2Uvc3ludGF4L2NwcC9wcm92aWRlci9nbyNQYXJzZTU6c2NvcGUxOiMxNjpkb2NzL3N5bnRheC5tZCNjNTpmYWxzZTA6"></a>

Authored scope: [docs/syntax.md#c](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#c); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/parse/syntax/cpp/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/cpp/provider/go/parse.go#L441)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L115) (external) |
| source-boundary | unsupported-and-malformed-fixtures | located-error-without-output | unspecified (authored profile absent) | [TestUnsupportedAndMalformedSources](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L137) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjcyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvcnVzdC9wcm92aWRlci9nbyNQcm92aWRlci5DaGVjazEwODpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL2NoZWNrL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nby5MYW5ndWFnZV02NjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY2hlY2svc3ludGF4L3J1c3QvcHJvdmlkZXIvZ28uUHJvdmlkZXI1MTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28jTGFuZ3VhZ2U0OnJ1c3Q"></a>

## Rust / github.com/Bitspark/bittype/check/syntax/rust/provider/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/check/syntax/rust/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/provider/go/check.go#L192). Receiver: github.com/Bitspark/bittype/check/syntax/rust/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/check/syntax/rust/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/provider/go/check.go#L192) → [github.com/Bitspark/bittype/check/syntax/rust/provider/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/provider/go/check.go#L72) ([call](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/provider/go/check.go#L192)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/check/syntax/rust/provider/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/provider/go/check.go#L72)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NDM1Ok9EcHdjbTkyYVdSbGNqY3lPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYUdWamF5OXplVzUwWVhndmNuVnpkQzl3Y205MmFXUmxjaTluYnlOUWNtOTJhV1JsY2k1RGFHVmphekV3T0RwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdlkyOXVkSEpoWTNSekwyTm9aV05yTDJkdkxsQnliM1pwWkdWeVcyZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOXplVzUwWVhndmNuVnpkQzluYnk1TVlXNW5kV0ZuWlYwMk5qcG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZZMmhsWTJzdmMzbHVkR0Y0TDNKMWMzUXZjSEp2ZG1sa1pYSXZaMjh1VUhKdmRtbGtaWEkxTVRwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDNKMWMzUXZaMjhqVEdGdVozVmhaMlUwT25KMWMzUTYzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvcnVzdC9wcm92aWRlci9nbyNDaGVjazc6Z3JhbW1hcjYzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvcnVzdC9wcm92aWRlci9nbyNydWxlczA6NTpmYWxzZTA6"></a>

Definition: [github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/provider/go/check.go#L31); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### file (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key file.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| item | 0 | unbounded | module, use, type-alias, struct, tuple-struct, enum, trait, trait-impl |

No feature-key case association at this exact authored profile.

#### module (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key module.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| item | 0 | unbounded | module, use, type-alias, struct, tuple-struct, enum, trait, trait-impl |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### use (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key use.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| path | 1 | 1 | path-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-alias (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-alias.

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

#### struct (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key struct.

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

#### tuple-struct (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-struct.

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

#### enum (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key enum.

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

#### unit-variant (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key unit-variant.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### tuple-variant (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-variant.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| field | 0 | unbounded | tuple-field |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### record-variant (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key record-variant.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| field | 0 | unbounded | field |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### field (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key field.

Roles: 1 declared. Attributes: 2 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### tuple-field (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-field.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| visibility | false | private, public, crate |

No feature-key case association at this exact authored profile.

#### trait (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key trait.

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

#### trait-impl (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key trait-impl.

Roles: 5 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| generic-parameter | 0 | unbounded | type-parameter, lifetime-parameter |
| trait | 1 | 1 | path-type |
| subject | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |
| where | 0 | unbounded | where-predicate |
| member | 0 | unbounded | type-alias |

No feature-key case association at this exact authored profile.

#### associated-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key associated-type.

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

#### method-signature (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key method-signature.

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

#### receiver (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key receiver.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| lifetime | 0 | 1 | lifetime-reference |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| receiver-form | true | value, shared, mutable |

No feature-key case association at this exact authored profile.

#### parameter (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| type | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-parameter (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key type-parameter.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| bound | 0 | unbounded | path-type, lifetime-reference, relaxed-bound |
| default | 0 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### lifetime-parameter (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key lifetime-parameter.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| bound | 0 | unbounded | lifetime-reference |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### lifetime-reference (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key lifetime-reference.

Roles: 0 declared. Attributes: 1 declared. References: target.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### where-predicate (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key where-predicate.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| subject | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type, lifetime-reference |
| bound | 1 | unbounded | path-type, lifetime-reference, relaxed-bound |

No feature-key case association at this exact authored profile.

#### relaxed-bound (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key relaxed-bound.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| path | 1 | 1 | path-type |

No feature-key case association at this exact authored profile.

#### path-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key path-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| segment | 1 | unbounded | path-segment |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| absolute | false | true, false |

No feature-key case association at this exact authored profile.

#### qualified-path-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key qualified-path-type.

Roles: 3 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| subject | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |
| trait | 0 | 1 | path-type |
| segment | 1 | unbounded | path-segment |

No feature-key case association at this exact authored profile.

#### path-segment (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key path-segment.

Roles: 1 declared. Attributes: 1 declared. References: target.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| argument | 0 | unbounded | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type, lifetime-reference, associated-type-binding |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### associated-type-binding (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key associated-type-binding.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| path | 1 | 1 | path-segment |
| type | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

No feature-key case association at this exact authored profile.

#### primitive-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key primitive-type.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | bool, char, str, f32, f64, i8, i16, i32, i64, i128, isize, u8, u16, u32, u64, u128, usize |

No feature-key case association at this exact authored profile.

#### reference-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key reference-type.

Roles: 2 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |
| lifetime | 0 | 1 | lifetime-reference |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| mutable | false | true, false |

No feature-key case association at this exact authored profile.

#### raw-pointer-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key raw-pointer-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| mutable | false | true, false |

No feature-key case association at this exact authored profile.

#### slice-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key slice-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

No feature-key case association at this exact authored profile.

#### array-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key array-type.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| length | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### tuple-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 0 | unbounded | path-type, qualified-path-type, primitive-type, reference-type, raw-pointer-type, slice-type, array-type, tuple-type, never-type, function-type, trait-object-type |

No feature-key case association at this exact authored profile.

#### never-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key never-type.

Roles: 0 declared. Attributes: 0 declared. References: none.

No feature-key case association at this exact authored profile.

#### function-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key function-type.

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

#### trait-object-type (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/check/syntax/rust/provider/go#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key trait-object-type.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| bound | 1 | unbounded | path-type, lifetime-reference |

No feature-key case association at this exact authored profile.

<a id="binding-NDM1Ok9EcHdjbTkyYVdSbGNqY3lPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYUdWamF5OXplVzUwWVhndmNuVnpkQzl3Y205MmFXUmxjaTluYnlOUWNtOTJhV1JsY2k1RGFHVmphekV3T0RwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdlkyOXVkSEpoWTNSekwyTm9aV05yTDJkdkxsQnliM1pwWkdWeVcyZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOXplVzUwWVhndmNuVnpkQzluYnk1TVlXNW5kV0ZuWlYwMk5qcG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZZMmhsWTJzdmMzbHVkR0Y0TDNKMWMzUXZjSEp2ZG1sa1pYSXZaMjh1VUhKdmRtbGtaWEkxTVRwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDNKMWMzUXZaMjhqVEdGdVozVmhaMlUwT25KMWMzUTYzOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvcnVzdC9wcm92aWRlci9nbyNDaGVjazU6c2NvcGUxOiMxOTpkb2NzL3N5bnRheC5tZCNydXN0NTpmYWxzZTA6"></a>

Authored scope: [docs/syntax.md#rust](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#rust); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| associated-type-projections | independent-ast-and-references | accept | unspecified (authored profile absent) | [TestAssociatedTypeAdmission](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/provider/go/check_test.go#L36) (external) |
| associated-type-projections | invalid-shapes-and-ordering | located-rejection | unspecified (authored profile absent) | [TestInvalidAssociatedTypeShapes](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/provider/go/check_test.go#L53) (external) |
| generic-associated-declarations | independent-ast | accept | unspecified (authored profile absent) | [TestGenericAssociatedDeclarationAdmission](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/provider/go/declarations_test.go#L45) (external) |
| generic-associated-declarations | relaxed-bound-context | context-sensitive-admission | unspecified (authored profile absent) | [TestRelaxedBoundContext](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/provider/go/declarations_test.go#L53) (external) |
| generic-associated-declarations | invalid-declaration-contexts | located-rejection | unspecified (authored profile absent) | [TestGenericAssociatedContextRules](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/provider/go/declarations_test.go#L73) (external) |

### Declared target [github.com/Bitspark/bittype/check/syntax/rust/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/provider/go/check.go#L192)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| ast-admission | malformed-carriers | reject | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L193) (external) |

<a id="cap-ODpwcm92aWRlcjcyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9maWxlcy9ob3Jvcy9wcm92aWRlci9nbyNQcm92aWRlci5QYXJzZTI3MDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3BhcnNlL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcGFyc2UvZ28uRmlsZXNSZXF1ZXN0W2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc291cmNlL2dvLlNvdXJjZVtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXV02NjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcGFyc2UvZmlsZXMvaG9yb3MvcHJvdmlkZXIvZ28uUHJvdmlkZXI1MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvI0xhbmd1YWdlNTpob3Jvcw"></a>

## Horos / github.com/Bitspark/bittype/parse/files/horos/provider/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.FilesRequest&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;, github.com/Bitspark/bittype/source/go.Source&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/parse/files/horos/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/files/horos/provider/go/parse.go#L49). Receiver: github.com/Bitspark/bittype/parse/files/horos/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

### Declared target [github.com/Bitspark/bittype/parse/files/horos/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/files/horos/provider/go/parse.go#L49)

<a id="binding-NjU0Ok9EcHdjbTkyYVdSbGNqY3lPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3WVhKelpTOW1hV3hsY3k5b2IzSnZjeTl3Y205MmFXUmxjaTluYnlOUWNtOTJhV1JsY2k1UVlYSnpaVEkzTURwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdlkyOXVkSEpoWTNSekwzQmhjbk5sTDJkdkxsQnliM1pwWkdWeVcyZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOWpiMjUwY21GamRITXZjR0Z5YzJVdloyOHVSbWxzWlhOU1pYRjFaWE4wVzJkcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5emVXNTBZWGd2YUc5eWIzTXZaMjh1VEdGdVozVmhaMlZkTENCbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMyOTFjbU5sTDJkdkxsTnZkWEpqWlZ0bmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDJodmNtOXpMMmR2TGt4aGJtZDFZV2RsWFYwMk5qcG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjR0Z5YzJVdlptbHNaWE12YUc5eWIzTXZjSEp2ZG1sa1pYSXZaMjh1VUhKdmRtbGtaWEkxTWpwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDJodmNtOXpMMmR2STB4aGJtZDFZV2RsTlRwb2IzSnZjdzcyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9maWxlcy9ob3Jvcy9wcm92aWRlci9nbyNQcm92aWRlci5QYXJzZTU6c2NvcGUxOiM0Nzpkb2NzL3NvdXJjZS5tZCNob3Jvcy1jb21wb3NpdGlvbi1hbmQtcGFnaW5hdGlvbjU6ZmFsc2UwOg"></a>

Authored scope: [docs/source.md#horos-composition-and-pagination](https://github.com/Bitspark/bittype/blob/main/docs/source.md#horos-composition-and-pagination); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| multi-file-pagination | receipt-preserving-round-trips | module-context-preserved | unspecified (authored profile absent) | [TestFilesPaginationAndReceiptPreservingRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/files_test.go#L85) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjcyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvcnVzdC9wcm92aWRlci9nbyNQcm92aWRlci5QYXJzZTIwODpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3BhcnNlL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvcGFyc2UvZ28uSW5wdXQsIGdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9jc3QvZ28uQ1NUW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nby5MYW5ndWFnZV1dNjY6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3BhcnNlL3N5bnRheC9ydXN0L3Byb3ZpZGVyL2dvLlByb3ZpZGVyNTE6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvI0xhbmd1YWdlNDpydXN0"></a>

## Rust / github.com/Bitspark/bittype/parse/syntax/rust/provider/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/parse/syntax/rust/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/provider/go/parse.go#L612). Receiver: github.com/Bitspark/bittype/parse/syntax/rust/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/parse/syntax/rust/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/provider/go/parse.go#L612) → [github.com/Bitspark/bittype/parse/syntax/rust/provider/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/provider/go/parse.go#L25) ([call](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/provider/go/parse.go#L613)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/rust/provider/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/provider/go/parse.go#L25)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NTY4Ok9EcHdjbTkyYVdSbGNqY3lPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3WVhKelpTOXplVzUwWVhndmNuVnpkQzl3Y205MmFXUmxjaTluYnlOUWNtOTJhV1JsY2k1UVlYSnpaVEl3T0RwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdlkyOXVkSEpoWTNSekwzQmhjbk5sTDJkdkxsQnliM1pwWkdWeVcyZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOWpiMjUwY21GamRITXZjR0Z5YzJVdloyOHVTVzV3ZFhRc0lHZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOTBjbVZsY3k5amMzUXZaMjh1UTFOVVcyZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOXplVzUwWVhndmNuVnpkQzluYnk1TVlXNW5kV0ZuWlYxZE5qWTZaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM0JoY25ObEwzTjViblJoZUM5eWRYTjBMM0J5YjNacFpHVnlMMmR2TGxCeWIzWnBaR1Z5TlRFNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzTjViblJoZUM5eWRYTjBMMmR2STB4aGJtZDFZV2RsTkRweWRYTjA2MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcGFyc2Uvc3ludGF4L3J1c3QvcHJvdmlkZXIvZ28jUGFyc2U1OnNjb3BlMTojMTk6ZG9jcy9zeW50YXgubWQjcnVzdDU6ZmFsc2UwOg"></a>

Authored scope: [docs/syntax.md#rust](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#rust); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| generic-associated-declarations | declaration-and-impl-structure | structure-preserved | unspecified (authored profile absent) | [TestGenericAssociatedDeclarationStructure](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/provider/go/declarations_test.go#L16) (external) |
| generic-associated-declarations | unresolved-syntax | accept | unspecified (authored profile absent) | [TestGenericAssociatedSyntaxWithoutResolution](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/provider/go/declarations_test.go#L46) (external) |
| generic-associated-declarations | malformed-declarations | located-error-without-output | unspecified (authored profile absent) | [TestMalformedGenericAssociatedDeclarations](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/provider/go/declarations_test.go#L67) (external) |
| type-only-impl-scope | unsupported-impl-forms | unsupported-without-output | unspecified (authored profile absent) | [TestImplScopeRejection](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/provider/go/declarations_test.go#L94) (external) |
| associated-type-projections | structured-projections-and-bindings | structure-preserved | unspecified (authored profile absent) | [TestAssociatedTypeStructure](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/provider/go/parse_test.go#L17) (external) |
| associated-type-projections | unresolved-projection-syntax | accept | unspecified (authored profile absent) | [TestAssociatedTypeSyntaxWithoutResolution](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/provider/go/parse_test.go#L57) (external) |
| associated-type-projections | malformed-and-bounded-projections | located-error-without-output | unspecified (authored profile absent) | [TestMalformedAssociatedTypes](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/provider/go/parse_test.go#L80) (external) |
| associated-type-projections | associated-bound-syntax | unsupported-without-output | unspecified (authored profile absent) | [TestAssociatedTypeBoundsRemainUnsupported](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/provider/go/parse_test.go#L106) (external) |

### Declared target [github.com/Bitspark/bittype/parse/syntax/rust/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/provider/go/parse.go#L612)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L115) (external) |
| source-boundary | unsupported-and-malformed-fixtures | located-error-without-output | unspecified (authored profile absent) | [TestUnsupportedAndMalformedSources](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L137) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjczOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvaG9yb3MvcHJvdmlkZXIvZ28jUHJvdmlkZXIuQ2hlY2sxMDk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9jaGVjay9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXTY3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jaGVjay9zeW50YXgvaG9yb3MvcHJvdmlkZXIvZ28uUHJvdmlkZXI1MjpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvI0xhbmd1YWdlNTpob3Jvcw"></a>

## Horos / github.com/Bitspark/bittype/check/syntax/horos/provider/go#Provider.Check

Contract: github.com/Bitspark/bittype/contracts/check/go.Provider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/check/syntax/horos/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/horos/provider/go/check.go#L18). Receiver: github.com/Bitspark/bittype/check/syntax/horos/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/check/syntax/horos/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/horos/provider/go/check.go#L18) → [github.com/Bitspark/bittype/check/syntax/horos/provider/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/horos/provider/go/check.go#L13) ([call](https://github.com/Bitspark/bittype/blob/main/check/syntax/horos/provider/go/check.go#L18)); receiver ; type arguments .

Traversal boundary: callee is dynamic or outside implementation source.

### Declared target [github.com/Bitspark/bittype/check/syntax/horos/provider/go#Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/horos/provider/go/check.go#L13)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NDQyOk9EcHdjbTkyYVdSbGNqY3pPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYUdWamF5OXplVzUwWVhndmFHOXliM012Y0hKdmRtbGtaWEl2WjI4alVISnZkbWxrWlhJdVEyaGxZMnN4TURrNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwyTnZiblJ5WVdOMGN5OWphR1ZqYXk5bmJ5NVFjbTkyYVdSbGNsdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwyaHZjbTl6TDJkdkxreGhibWQxWVdkbFhUWTNPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYUdWamF5OXplVzUwWVhndmFHOXliM012Y0hKdmRtbGtaWEl2WjI4dVVISnZkbWxrWlhJMU1qcG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwyaHZjbTl6TDJkdkkweGhibWQxWVdkbE5UcG9iM0p2Y3c2NDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY2hlY2svc3ludGF4L2hvcm9zL3Byb3ZpZGVyL2dvI0NoZWNrNzpncmFtbWFyNTM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2ludGVybmFsL2RlY2xzeW50YXgjcnVsZXMwOjU6ZmFsc2UwOg"></a>

Definition: [github.com/Bitspark/bittype/internal/declsyntax#rules](https://github.com/Bitspark/bittype/blob/main/internal/declsyntax/check.go#L66); format: github.com/Bitspark/bittype/internal/grammar#Rules; extraction: **partial**; profile: unspecified (authored profile absent).

#### program (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key program.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| module | 1 | unbounded | module |

No feature-key case association at this exact authored profile.

#### module (source unavailable)

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

#### import (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key import.

Roles: 0 declared. Attributes: 3 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| alias | true | opaque: Predicate present; behavior not interpreted |
| name | true | opaque: Predicate present; behavior not interpreted |
| pin | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### type-definition (source unavailable)

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

#### kind-definition (source unavailable)

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

#### type-parameter (source unavailable)

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

#### kind-parameter (source unavailable)

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

#### parameter (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key parameter.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### reference (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key reference.

Roles: 0 declared. Attributes: 3 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| scope | true | local, imported |
| name | true | opaque: Predicate present; behavior not interpreted |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### application (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key application.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| callee | 1 | 1 | descriptor, parameter, reference, application |
| argument | 1 | unbounded | descriptor, parameter, reference, application |

No feature-key case association at this exact authored profile.

#### descriptor (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key descriptor.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| operand | 0 | unbounded | descriptor, parameter, reference, application, atom, tuple |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| constructor | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### atom (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key atom.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| hex | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### tuple (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key tuple.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 0 | unbounded | descriptor, parameter, reference, application, atom, tuple |

No feature-key case association at this exact authored profile.

#### proper-sort (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key proper-sort.

Roles: 0 declared. Attributes: 0 declared. References: none.

No feature-key case association at this exact authored profile.

#### arrow-sort (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key arrow-sort.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| input | 1 | unbounded | proper-sort, arrow-sort |
| result | 1 | 1 | proper-sort, arrow-sort |

No feature-key case association at this exact authored profile.

#### kind-top (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-top.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| sort | 1 | 1 | proper-sort, arrow-sort |

No feature-key case association at this exact authored profile.

#### kind-bottom (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-bottom.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| sort | 1 | 1 | proper-sort, arrow-sort |

No feature-key case association at this exact authored profile.

#### kind-parameter-reference (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-parameter-reference.

Roles: 0 declared. Attributes: 1 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-reference (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-reference.

Roles: 0 declared. Attributes: 3 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| scope | true | local, imported |
| name | true | opaque: Predicate present; behavior not interpreted |
| alias | false | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-application (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-application.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| callee | 1 | 1 | kind-reference |
| argument | 1 | unbounded | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### kind-union (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-union.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| sort | 1 | 1 | proper-sort, arrow-sort |
| operand | 2 | unbounded | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### kind-intersection (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-intersection.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| sort | 1 | 1 | proper-sort, arrow-sort |
| operand | 2 | unbounded | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### kind-arrow (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-arrow.

Roles: 2 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| input | 1 | unbounded | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |
| result | 1 | 1 | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### kind-node (source unavailable)

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

#### selector (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key selector.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| name | true | selector, mode |
| value | true | int, utf8-text, bool, open, closed |

No feature-key case association at this exact authored profile.

#### kind-recursive (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-recursive.

Roles: 1 declared. Attributes: 1 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| equation | 1 | unbounded | kind-equation |

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### kind-equation (source unavailable)

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

#### kind-recursive-variable (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key kind-recursive-variable.

Roles: 0 declared. Attributes: 2 declared. References: none.

| Attribute | Required | Values or predicate |
| --- | --- | --- |
| depth | true | opaque: Predicate present; behavior not interpreted |
| index | true | opaque: Predicate present; behavior not interpreted |

No feature-key case association at this exact authored profile.

#### sequence-epsilon (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-epsilon.

Roles: 0 declared. Attributes: 0 declared. References: none.

No feature-key case association at this exact authored profile.

#### sequence-item (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-item.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | kind-top, kind-bottom, kind-parameter-reference, kind-reference, kind-application, kind-union, kind-intersection, kind-arrow, kind-node, kind-recursive, kind-recursive-variable |

No feature-key case association at this exact authored profile.

#### sequence-concat (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-concat.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 2 | unbounded | sequence-epsilon, sequence-item, sequence-concat, sequence-alternative, sequence-star |

No feature-key case association at this exact authored profile.

#### sequence-alternative (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-alternative.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 2 | unbounded | sequence-epsilon, sequence-item, sequence-concat, sequence-alternative, sequence-star |

No feature-key case association at this exact authored profile.

#### sequence-star (source unavailable)

Feature identity: definition github.com/Bitspark/bittype/internal/declsyntax#rules / format github.com/Bitspark/bittype/internal/grammar#Rules / key sequence-star.

Roles: 1 declared. Attributes: 0 declared. References: none.

| Role | Minimum | Maximum | Allowed kinds |
| --- | --- | --- | --- |
| element | 1 | 1 | sequence-epsilon, sequence-item, sequence-concat, sequence-alternative, sequence-star |

No feature-key case association at this exact authored profile.

<a id="binding-NDQyOk9EcHdjbTkyYVdSbGNqY3pPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYUdWamF5OXplVzUwWVhndmFHOXliM012Y0hKdmRtbGtaWEl2WjI4alVISnZkbWxrWlhJdVEyaGxZMnN4TURrNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwyTnZiblJ5WVdOMGN5OWphR1ZqYXk5bmJ5NVFjbTkyYVdSbGNsdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwyaHZjbTl6TDJkdkxreGhibWQxWVdkbFhUWTNPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzlqYUdWamF5OXplVzUwWVhndmFHOXliM012Y0hKdmRtbGtaWEl2WjI4dVVISnZkbWxrWlhJMU1qcG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwyaHZjbTl6TDJkdkkweGhibWQxWVdkbE5UcG9iM0p2Y3c2NDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY2hlY2svc3ludGF4L2hvcm9zL3Byb3ZpZGVyL2dvI0NoZWNrNTpzY29wZTE6IzMxOmRvY3MvaG9yb3MubWQjc3ludGF4LXZvY2FidWxhcnk1OmZhbHNlMDo"></a>

Authored scope: [docs/horos.md#syntax-vocabulary](https://github.com/Bitspark/bittype/blob/main/docs/horos.md#syntax-vocabulary); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/check/syntax/horos/provider/go#Provider.Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/horos/provider/go/check.go#L18)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| independent-engine-stages | pagination-source | independently-composable | unspecified (authored profile absent) | [TestHorosSourceThroughIndependentEngineStages](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_text_test.go#L27) (external) |
| horos-text-vocabulary | declared-constructor-fixtures | canonical-round-trip-preserved | unspecified (authored profile absent) | [TestCompleteSourceVocabulary](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/provider/go/parse_test.go#L85) (external) |
| syntax-semantics-boundary | semantically-invalid-syntax-fixtures | syntax-accepted | unspecified (authored profile absent) | [TestSyntaxDoesNotResolveOrCheckContracts](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/provider/go/parse_test.go#L176) (external) |

<a id="cap-ODpwcm92aWRlcjczOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L2NwcC9wcm92aWRlci9nbyNQcm92aWRlci5Gb3JtYXQyNTM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9mb3JtYXQvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2FzdC9nby5WaWV3W2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvY3BwL2dvLkxhbmd1YWdlXSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3RyZWVzL2NzdC9nby5DU1RbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9jcHAvZ28uTGFuZ3VhZ2VdXTY2OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L2NwcC9wcm92aWRlci9nby5Qcm92aWRlcjUwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvY3BwL2dvI0xhbmd1YWdlMzpjcHA"></a>

## C++ / github.com/Bitspark/bittype/format/syntax/cpp/provider/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/format/syntax/cpp/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/cpp/provider/go/format.go#L316). Receiver: github.com/Bitspark/bittype/format/syntax/cpp/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/format/syntax/cpp/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/cpp/provider/go/format.go#L316) → [github.com/Bitspark/bittype/format/syntax/cpp/provider/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/cpp/provider/go/format.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/format/syntax/cpp/provider/go/format.go#L317)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/cpp/provider/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/cpp/provider/go/format.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NjI3Ok9EcHdjbTkyYVdSbGNqY3pPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzltYjNKdFlYUXZjM2x1ZEdGNEwyTndjQzl3Y205MmFXUmxjaTluYnlOUWNtOTJhV1JsY2k1R2IzSnRZWFF5TlRNNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwyTnZiblJ5WVdOMGN5OW1iM0p0WVhRdloyOHVVSEp2ZG1sa1pYSmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM1J5WldWekwyRnpkQzluYnk1V2FXVjNXMmRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl6ZVc1MFlYZ3ZZM0J3TDJkdkxreGhibWQxWVdkbFhTd2daMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM1J5WldWekwyTnpkQzluYnk1RFUxUmJaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOWpjSEF2WjI4dVRHRnVaM1ZoWjJWZFhUWTJPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzltYjNKdFlYUXZjM2x1ZEdGNEwyTndjQzl3Y205MmFXUmxjaTluYnk1UWNtOTJhV1JsY2pVd09tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOXplVzUwWVhndlkzQndMMmR2STB4aGJtZDFZV2RsTXpwamNIQTY0OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L2NwcC9wcm92aWRlci9nbyNGb3JtYXQ1OnNjb3BlMTojMTY6ZG9jcy9zeW50YXgubWQjYzU6ZmFsc2UwOg"></a>

Authored scope: [docs/syntax.md#c](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#c); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/format/syntax/cpp/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/cpp/provider/go/format.go#L316)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L106) (external) |
| ast-admission | malformed-carriers | rejection-without-output | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L193) (external) |
| native-compiler-qualification | formatted-fixture-contract | compiler-fixture-contract | unspecified (authored profile absent) | [TestNativeCompilerFixtures](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L265) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L115) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjczOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wYXJzZS9zeW50YXgvaG9yb3MvcHJvdmlkZXIvZ28jUHJvdmlkZXIuUGFyc2UyMDk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wYXJzZS9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3BhcnNlL2dvLklucHV0LCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXV02NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcGFyc2Uvc3ludGF4L2hvcm9zL3Byb3ZpZGVyL2dvLlByb3ZpZGVyNTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nbyNMYW5ndWFnZTU6aG9yb3M"></a>

## Horos / github.com/Bitspark/bittype/parse/syntax/horos/provider/go#Provider.Parse

Contract: github.com/Bitspark/bittype/contracts/parse/go.Provider&#91;github.com/Bitspark/bittype/contracts/parse/go.Input, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/parse/syntax/horos/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/provider/go/parse.go#L566). Receiver: github.com/Bitspark/bittype/parse/syntax/horos/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/parse/syntax/horos/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/provider/go/parse.go#L566) → [github.com/Bitspark/bittype/parse/syntax/horos/provider/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/provider/go/parse.go#L61) ([call](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/provider/go/parse.go#L567)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/parse/syntax/horos/provider/go#Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/provider/go/parse.go#L61)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NTc1Ok9EcHdjbTkyYVdSbGNqY3pPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3WVhKelpTOXplVzUwWVhndmFHOXliM012Y0hKdmRtbGtaWEl2WjI4alVISnZkbWxrWlhJdVVHRnljMlV5TURrNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwyTnZiblJ5WVdOMGN5OXdZWEp6WlM5bmJ5NVFjbTkyYVdSbGNsdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZZMjl1ZEhKaFkzUnpMM0JoY25ObEwyZHZMa2x1Y0hWMExDQm5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZkSEpsWlhNdlkzTjBMMmR2TGtOVFZGdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwyaHZjbTl6TDJkdkxreGhibWQxWVdkbFhWMDJOenBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2Y0dGeWMyVXZjM2x1ZEdGNEwyaHZjbTl6TDNCeWIzWnBaR1Z5TDJkdkxsQnliM1pwWkdWeU5USTZaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOW9iM0p2Y3k5bmJ5Tk1ZVzVuZFdGblpUVTZhRzl5YjNNNjQ6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3BhcnNlL3N5bnRheC9ob3Jvcy9wcm92aWRlci9nbyNQYXJzZTU6c2NvcGUxOiM0MDpkb2NzL2hvcm9zLXN5bnRheC5tZCNwaXBlbGluZS1hbmQtbGltaXRzNTpmYWxzZTA6"></a>

Authored scope: [docs/horos-syntax.md#pipeline-and-limits](https://github.com/Bitspark/bittype/blob/main/docs/horos-syntax.md#pipeline-and-limits); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/parse/syntax/horos/provider/go#Provider.Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/provider/go/parse.go#L566)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| horos-modules-and-generics | pagination-and-imported-library | canonical-identity-preserved | unspecified (authored profile absent) | [TestCanonicalPaginationAndImportedLibraryRoundTrip](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/provider/go/format_test.go#L59) (external) |
| independent-engine-stages | pagination-source | independently-composable | unspecified (authored profile absent) | [TestHorosSourceThroughIndependentEngineStages](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_text_test.go#L27) (external) |
| horos-text-vocabulary | declared-constructor-fixtures | canonical-round-trip-preserved | unspecified (authored profile absent) | [TestCompleteSourceVocabulary](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/provider/go/parse_test.go#L85) (external) |
| syntax-semantics-boundary | semantically-invalid-syntax-fixtures | syntax-accepted | unspecified (authored profile absent) | [TestSyntaxDoesNotResolveOrCheckContracts](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/provider/go/parse_test.go#L176) (external) |
| source-boundary | whole-input-rejection | located-error-without-output | unspecified (authored profile absent) | [TestRejectWholeInput](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/provider/go/parse_test.go#L197) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjczOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9zaGFyZWQvcHJvdmlkZXIvZ28jUHJvdmlkZXIuUHJpbnQxOTc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wcmludC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc291cmNlL2dvLlNvdXJjZVtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nby5MYW5ndWFnZV0sIGdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9maWxlcy9nby5GaWxlc10xMTk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3ByaW50L2ZpbGVzL3NoYXJlZC9wcm92aWRlci9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nby5MYW5ndWFnZV01MDpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2NwcC9nbyNMYW5ndWFnZTM6Y3Bw"></a>

## C++ / github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/source/go.Source&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;, github.com/Bitspark/bittype/files/go.Files&#93;.

Public operation: [github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/shared/provider/go/print.go#L13). Receiver: github.com/Bitspark/bittype/print/files/shared/provider/go.Provider&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/cpp/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider + github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/cpp/go.Language&#93;, string&#93;.

### Declared target [github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/shared/provider/go/print.go#L13)

<a id="binding-NjI0Ok9EcHdjbTkyYVdSbGNqY3pPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3Y21sdWRDOW1hV3hsY3k5emFHRnlaV1F2Y0hKdmRtbGtaWEl2WjI4alVISnZkbWxrWlhJdVVISnBiblF4T1RjNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwyTnZiblJ5WVdOMGN5OXdjbWx1ZEM5bmJ5NVFjbTkyYVdSbGNsdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjMjkxY21ObEwyZHZMbE52ZFhKalpWdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwyTndjQzluYnk1TVlXNW5kV0ZuWlYwc0lHZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOW1hV3hsY3k5bmJ5NUdhV3hsYzEweE1UazZaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM0J5YVc1MEwyWnBiR1Z6TDNOb1lYSmxaQzl3Y205MmFXUmxjaTluYnk1UWNtOTJhV1JsY2x0bmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDJOd2NDOW5ieTVNWVc1bmRXRm5aVjAxTURwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDJOd2NDOW5ieU5NWVc1bmRXRm5aVE02WTNCdzczOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9zaGFyZWQvcHJvdmlkZXIvZ28jUHJvdmlkZXIuUHJpbnQ1OnNjb3BlMTojMzk6ZG9jcy9zb3VyY2UubWQjc3ltbWV0cmljLXJvbGUtY29udHJhY3RzNTpmYWxzZTA6"></a>

Authored scope: [docs/source.md#symmetric-role-contracts](https://github.com/Bitspark/bittype/blob/main/docs/source.md#symmetric-role-contracts); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjczOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9zaGFyZWQvcHJvdmlkZXIvZ28jUHJvdmlkZXIuUHJpbnQxOTg6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wcmludC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc291cmNlL2dvLlNvdXJjZVtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3J1c3QvZ28uTGFuZ3VhZ2VdLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvZmlsZXMvZ28uRmlsZXNdMTIwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9zaGFyZWQvcHJvdmlkZXIvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ydXN0L2dvLkxhbmd1YWdlXTUxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvcnVzdC9nbyNMYW5ndWFnZTQ6cnVzdA"></a>

## Rust / github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/source/go.Source&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;, github.com/Bitspark/bittype/files/go.Files&#93;.

Public operation: [github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/shared/provider/go/print.go#L13). Receiver: github.com/Bitspark/bittype/print/files/shared/provider/go.Provider&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/rust/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider + github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/rust/go.Language&#93;, string&#93;.

### Declared target [github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/shared/provider/go/print.go#L13)

<a id="binding-NjMwOk9EcHdjbTkyYVdSbGNqY3pPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3Y21sdWRDOW1hV3hsY3k5emFHRnlaV1F2Y0hKdmRtbGtaWEl2WjI4alVISnZkbWxrWlhJdVVISnBiblF4T1RnNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwyTnZiblJ5WVdOMGN5OXdjbWx1ZEM5bmJ5NVFjbTkyYVdSbGNsdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjMjkxY21ObEwyZHZMbE52ZFhKalpWdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwzSjFjM1F2WjI4dVRHRnVaM1ZoWjJWZExDQm5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZabWxzWlhNdloyOHVSbWxzWlhOZE1USXdPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3Y21sdWRDOW1hV3hsY3k5emFHRnlaV1F2Y0hKdmRtbGtaWEl2WjI4dVVISnZkbWxrWlhKYloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwzTjViblJoZUM5eWRYTjBMMmR2TGt4aGJtZDFZV2RsWFRVeE9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOXplVzUwWVhndmNuVnpkQzluYnlOTVlXNW5kV0ZuWlRRNmNuVnpkQTczOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9zaGFyZWQvcHJvdmlkZXIvZ28jUHJvdmlkZXIuUHJpbnQ1OnNjb3BlMTojMzk6ZG9jcy9zb3VyY2UubWQjc3ltbWV0cmljLXJvbGUtY29udHJhY3RzNTpmYWxzZTA6"></a>

Authored scope: [docs/source.md#symmetric-role-contracts](https://github.com/Bitspark/bittype/blob/main/docs/source.md#symmetric-role-contracts); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjczOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9zaGFyZWQvcHJvdmlkZXIvZ28jUHJvdmlkZXIuUHJpbnQxOTk6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wcmludC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc291cmNlL2dvLlNvdXJjZVtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2hvcm9zL2dvLkxhbmd1YWdlXSwgZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2ZpbGVzL2dvLkZpbGVzXTEyMTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvcHJpbnQvZmlsZXMvc2hhcmVkL3Byb3ZpZGVyL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvaG9yb3MvZ28uTGFuZ3VhZ2VdNTI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC9ob3Jvcy9nbyNMYW5ndWFnZTU6aG9yb3M"></a>

## Horos / github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/source/go.Source&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;, github.com/Bitspark/bittype/files/go.Files&#93;.

Public operation: [github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/shared/provider/go/print.go#L13). Receiver: github.com/Bitspark/bittype/print/files/shared/provider/go.Provider&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/horos/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider + github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/horos/go.Language&#93;, string&#93;.

### Declared target [github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/shared/provider/go/print.go#L13)

<a id="binding-NjM1Ok9EcHdjbTkyYVdSbGNqY3pPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3Y21sdWRDOW1hV3hsY3k5emFHRnlaV1F2Y0hKdmRtbGtaWEl2WjI4alVISnZkbWxrWlhJdVVISnBiblF4T1RrNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwyTnZiblJ5WVdOMGN5OXdjbWx1ZEM5bmJ5NVFjbTkyYVdSbGNsdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjMjkxY21ObEwyZHZMbE52ZFhKalpWdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwyaHZjbTl6TDJkdkxreGhibWQxWVdkbFhTd2daMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMMlpwYkdWekwyZHZMa1pwYkdWelhURXlNVHBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2Y0hKcGJuUXZabWxzWlhNdmMyaGhjbVZrTDNCeWIzWnBaR1Z5TDJkdkxsQnliM1pwWkdWeVcyZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOXplVzUwWVhndmFHOXliM012WjI4dVRHRnVaM1ZoWjJWZE5USTZaMmwwYUhWaUxtTnZiUzlDYVhSemNHRnlheTlpYVhSMGVYQmxMM041Ym5SaGVDOW9iM0p2Y3k5bmJ5Tk1ZVzVuZFdGblpUVTZhRzl5YjNNNzM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3ByaW50L2ZpbGVzL3NoYXJlZC9wcm92aWRlci9nbyNQcm92aWRlci5QcmludDU6c2NvcGUxOiMzOTpkb2NzL3NvdXJjZS5tZCNzeW1tZXRyaWMtcm9sZS1jb250cmFjdHM1OmZhbHNlMDo"></a>

Authored scope: [docs/source.md#symmetric-role-contracts](https://github.com/Bitspark/bittype/blob/main/docs/source.md#symmetric-role-contracts); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| multi-file-pagination | receipt-preserving-round-trips | file-content-and-placement-preserved | unspecified (authored profile absent) | [TestFilesPaginationAndReceiptPreservingRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/files_test.go#L85) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjczOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9zaGFyZWQvcHJvdmlkZXIvZ28jUHJvdmlkZXIuUHJpbnQyMDA6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wcmludC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc291cmNlL2dvLlNvdXJjZVtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nby5MYW5ndWFnZV0sIGdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9maWxlcy9nby5GaWxlc10xMjI6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3ByaW50L2ZpbGVzL3NoYXJlZC9wcm92aWRlci9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nby5MYW5ndWFnZV01MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L2dvbGFuZy9nbyNMYW5ndWFnZTI6Z28"></a>

## Go / github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/source/go.Source&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;, github.com/Bitspark/bittype/files/go.Files&#93;.

Public operation: [github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/shared/provider/go/print.go#L13). Receiver: github.com/Bitspark/bittype/print/files/shared/provider/go.Provider&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/golang/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider + github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/golang/go.Language&#93;, string&#93;.

### Declared target [github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/shared/provider/go/print.go#L13)

<a id="binding-NjM1Ok9EcHdjbTkyYVdSbGNqY3pPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3Y21sdWRDOW1hV3hsY3k5emFHRnlaV1F2Y0hKdmRtbGtaWEl2WjI4alVISnZkbWxrWlhJdVVISnBiblF5TURBNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwyTnZiblJ5WVdOMGN5OXdjbWx1ZEM5bmJ5NVFjbTkyYVdSbGNsdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjMjkxY21ObEwyZHZMbE52ZFhKalpWdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwyZHZiR0Z1Wnk5bmJ5NU1ZVzVuZFdGblpWMHNJR2RwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzltYVd4bGN5OW5ieTVHYVd4bGMxMHhNakk2WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNCeWFXNTBMMlpwYkdWekwzTm9ZWEpsWkM5d2NtOTJhV1JsY2k5bmJ5NVFjbTkyYVdSbGNsdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwyZHZiR0Z1Wnk5bmJ5NU1ZVzVuZFdGblpWMDFNenBuYVhSb2RXSXVZMjl0TDBKcGRITndZWEpyTDJKcGRIUjVjR1V2YzNsdWRHRjRMMmR2YkdGdVp5OW5ieU5NWVc1bmRXRm5aVEk2WjI4NzM6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3ByaW50L2ZpbGVzL3NoYXJlZC9wcm92aWRlci9nbyNQcm92aWRlci5QcmludDU6c2NvcGUxOiMzOTpkb2NzL3NvdXJjZS5tZCNzeW1tZXRyaWMtcm9sZS1jb250cmFjdHM1OmZhbHNlMDo"></a>

Authored scope: [docs/source.md#symmetric-role-contracts](https://github.com/Bitspark/bittype/blob/main/docs/source.md#symmetric-role-contracts); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjczOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9zaGFyZWQvcHJvdmlkZXIvZ28jUHJvdmlkZXIuUHJpbnQyMDQ6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL2NvbnRyYWN0cy9wcmludC9nby5Qcm92aWRlcltnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc291cmNlL2dvLlNvdXJjZVtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28uTGFuZ3VhZ2VdLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvZmlsZXMvZ28uRmlsZXNdMTI2OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9zaGFyZWQvcHJvdmlkZXIvZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvLkxhbmd1YWdlXTU3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nbyNMYW5ndWFnZTEwOnR5cGVzY3JpcHQ"></a>

## TypeScript / github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print

Contract: github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/source/go.Source&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;, github.com/Bitspark/bittype/files/go.Files&#93;.

Public operation: [github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/shared/provider/go/print.go#L13). Receiver: github.com/Bitspark/bittype/print/files/shared/provider/go.Provider&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: github.com/Bitspark/bittype/contracts/language/go.Provider + github.com/Bitspark/bittype/contracts/print/go.Provider&#91;github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;, string&#93;.

### Declared target [github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print](https://github.com/Bitspark/bittype/blob/main/print/files/shared/provider/go/print.go#L13)

<a id="binding-NjYzOk9EcHdjbTkyYVdSbGNqY3pPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl3Y21sdWRDOW1hV3hsY3k5emFHRnlaV1F2Y0hKdmRtbGtaWEl2WjI4alVISnZkbWxrWlhJdVVISnBiblF5TURRNloybDBhSFZpTG1OdmJTOUNhWFJ6Y0dGeWF5OWlhWFIwZVhCbEwyTnZiblJ5WVdOMGN5OXdjbWx1ZEM5bmJ5NVFjbTkyYVdSbGNsdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjMjkxY21ObEwyZHZMbE52ZFhKalpWdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwzUjVjR1Z6WTNKcGNIUXZaMjh1VEdGdVozVmhaMlZkTENCbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdlptbHNaWE12WjI4dVJtbHNaWE5kTVRJMk9tZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOXdjbWx1ZEM5bWFXeGxjeTl6YUdGeVpXUXZjSEp2ZG1sa1pYSXZaMjh1VUhKdmRtbGtaWEpiWjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNONWJuUmhlQzkwZVhCbGMyTnlhWEIwTDJkdkxreGhibWQxWVdkbFhUVTNPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl6ZVc1MFlYZ3ZkSGx3WlhOamNtbHdkQzluYnlOTVlXNW5kV0ZuWlRFd09uUjVjR1Z6WTNKcGNIUTczOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9wcmludC9maWxlcy9zaGFyZWQvcHJvdmlkZXIvZ28jUHJvdmlkZXIuUHJpbnQ1OnNjb3BlMTojMzk6ZG9jcy9zb3VyY2UubWQjc3ltbWV0cmljLXJvbGUtY29udHJhY3RzNTpmYWxzZTA6"></a>

Authored scope: [docs/source.md#symmetric-role-contracts](https://github.com/Bitspark/bittype/blob/main/docs/source.md#symmetric-role-contracts); profile: unspecified (authored profile absent).

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

No associated cases found in the selected test sources. This does not establish absence of tests.

<a id="cap-ODpwcm92aWRlcjgwOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L3R5cGVzY3JpcHQvcHJvdmlkZXIvZ28jUHJvdmlkZXIuRm9ybWF0MjY3OmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9jb250cmFjdHMvZm9ybWF0L2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS90cmVlcy9hc3QvZ28uVmlld1tnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28uTGFuZ3VhZ2VdLCBnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvdHJlZXMvY3N0L2dvLkNTVFtnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28uTGFuZ3VhZ2VdXTczOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L3R5cGVzY3JpcHQvcHJvdmlkZXIvZ28uUHJvdmlkZXI1NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28jTGFuZ3VhZ2UxMDp0eXBlc2NyaXB0"></a>

## TypeScript / github.com/Bitspark/bittype/format/syntax/typescript/provider/go#Provider.Format

Contract: github.com/Bitspark/bittype/contracts/format/go.Provider&#91;github.com/Bitspark/bittype/trees/ast/go.View&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;, github.com/Bitspark/bittype/trees/cst/go.CST&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;&#93;.

Public operation: [github.com/Bitspark/bittype/format/syntax/typescript/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/typescript/provider/go/format.go#L281). Receiver: github.com/Bitspark/bittype/format/syntax/typescript/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/format/syntax/typescript/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/typescript/provider/go/format.go#L281) → [github.com/Bitspark/bittype/format/syntax/typescript/provider/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/typescript/provider/go/format.go#L24) ([call](https://github.com/Bitspark/bittype/blob/main/format/syntax/typescript/provider/go/format.go#L282)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/format/syntax/typescript/provider/go#Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/typescript/provider/go/format.go#L24)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-Njg0Ok9EcHdjbTkyYVdSbGNqZ3dPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzltYjNKdFlYUXZjM2x1ZEdGNEwzUjVjR1Z6WTNKcGNIUXZjSEp2ZG1sa1pYSXZaMjhqVUhKdmRtbGtaWEl1Um05eWJXRjBNalkzT21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5amIyNTBjbUZqZEhNdlptOXliV0YwTDJkdkxsQnliM1pwWkdWeVcyZHBkR2gxWWk1amIyMHZRbWwwYzNCaGNtc3ZZbWwwZEhsd1pTOTBjbVZsY3k5aGMzUXZaMjh1Vm1sbGQxdG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwzUjVjR1Z6WTNKcGNIUXZaMjh1VEdGdVozVmhaMlZkTENCbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmRISmxaWE12WTNOMEwyZHZMa05UVkZ0bmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDNSNWNHVnpZM0pwY0hRdloyOHVUR0Z1WjNWaFoyVmRYVGN6T21kcGRHaDFZaTVqYjIwdlFtbDBjM0JoY21zdlltbDBkSGx3WlM5bWIzSnRZWFF2YzNsdWRHRjRMM1I1Y0dWelkzSnBjSFF2Y0hKdmRtbGtaWEl2WjI4dVVISnZkbWxrWlhJMU56cG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZjM2x1ZEdGNEwzUjVjR1Z6WTNKcGNIUXZaMjhqVEdGdVozVmhaMlV4TURwMGVYQmxjMk55YVhCMDcxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9mb3JtYXQvc3ludGF4L3R5cGVzY3JpcHQvcHJvdmlkZXIvZ28jRm9ybWF0NTpzY29wZTE6IzI1OmRvY3Mvc3ludGF4Lm1kI3R5cGVzY3JpcHQ1OmZhbHNlMDo"></a>

Authored scope: [docs/syntax.md#typescript](https://github.com/Bitspark/bittype/blob/main/docs/syntax.md#typescript); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/format/syntax/typescript/provider/go#Provider.Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/typescript/provider/go/format.go#L281)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| provider-contract | native-fixture-contract | contract-preserved | unspecified (authored profile absent) | [TestNativeProviderConformance](https://github.com/Bitspark/bittype/blob/main/integration/go/conformance_test.go#L59) (external) |
| independent-ast-formatting | independent-ast-goldens | golden-text | unspecified (authored profile absent) | [TestIndependentFormatterGoldens](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L106) (external) |
| snapshot-ownership | input-mutation-after-format | output-unchanged | unspecified (authored profile absent) | [TestFormatterSnapshotOwnership](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L123) (external) |
| ast-admission | malformed-carriers | rejection-without-output | unspecified (authored profile absent) | [TestMalformedASTs](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L193) (external) |
| native-compiler-qualification | formatted-fixture-contract | compiler-fixture-contract | unspecified (authored profile absent) | [TestNativeCompilerFixtures](https://github.com/Bitspark/bittype/blob/main/integration/go/formatter_fixtures_test.go#L265) (external) |
| native-text | valid-fixture-round-trips | round-trip-preserved | unspecified (authored profile absent) | [TestNativeTextRoundTrips](https://github.com/Bitspark/bittype/blob/main/integration/go/native_text_test.go#L115) (external) |
| conditional-mapped-syntax | independent-ast-formatting | ast-preserved-and-golden-text | unspecified (authored profile absent) | [TestIndependentTypeScriptMappedFormatter](https://github.com/Bitspark/bittype/blob/main/integration/go/typescript_conditional_mapped_test.go#L37) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjgxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9hbmFseXplL2ZpbGVzL3R5cGVzY3JpcHQvcHJvdmlkZXIvZ28jUHJvdmlkZXIuQW5hbHl6ZTExNzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL2FuYWx5c2lzL2dvLlByb3ZpZGVyW2dpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9zeW50YXgvdHlwZXNjcmlwdC9nby5MYW5ndWFnZV03MzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvYW5hbHl6ZS9maWxlcy90eXBlc2NyaXB0L3Byb3ZpZGVyL2dvLlByb3ZpZGVyNTc6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvI0xhbmd1YWdlMTA6dHlwZXNjcmlwdA"></a>

## TypeScript / github.com/Bitspark/bittype/analyze/files/typescript/provider/go#Provider.Analyze

Contract: github.com/Bitspark/bittype/contracts/analysis/go.Provider&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/analyze/files/typescript/provider/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/provider/go/analyze.go#L57). Receiver: github.com/Bitspark/bittype/analyze/files/typescript/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

Static call relationships (no evidence transfer):

- [github.com/Bitspark/bittype/analyze/files/typescript/provider/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/provider/go/analyze.go#L57) → [github.com/Bitspark/bittype/analyze/files/typescript/provider/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/provider/go/analyze.go#L23) ([call](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/provider/go/analyze.go#L58)); receiver ; type arguments .

Traversal boundary: not a single-return static call.

### Declared target [github.com/Bitspark/bittype/analyze/files/typescript/provider/go#Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/provider/go/analyze.go#L23)

Related declaration; its cases and observations do not qualify the public wrapper.

<a id="binding-NDg2Ok9EcHdjbTkyYVdSbGNqZ3hPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzloYm1Gc2VYcGxMMlpwYkdWekwzUjVjR1Z6WTNKcGNIUXZjSEp2ZG1sa1pYSXZaMjhqVUhKdmRtbGtaWEl1UVc1aGJIbDZaVEV4TnpwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdlkyOXVkSEpoWTNSekwyRnVZV3g1YzJsekwyZHZMbEJ5YjNacFpHVnlXMmRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl6ZVc1MFlYZ3ZkSGx3WlhOamNtbHdkQzluYnk1TVlXNW5kV0ZuWlYwM016cG5hWFJvZFdJdVkyOXRMMEpwZEhOd1lYSnJMMkpwZEhSNWNHVXZZVzVoYkhsNlpTOW1hV3hsY3k5MGVYQmxjMk55YVhCMEwzQnliM1pwWkdWeUwyZHZMbEJ5YjNacFpHVnlOVGM2WjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNONWJuUmhlQzkwZVhCbGMyTnlhWEIwTDJkdkkweGhibWQxWVdkbE1UQTZkSGx3WlhOamNtbHdkQTcyOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9hbmFseXplL2ZpbGVzL3R5cGVzY3JpcHQvcHJvdmlkZXIvZ28jQW5hbHl6ZTU6c2NvcGUxOiM0OTpkb2NzL25hdGl2ZS1zZW1hbnRpY3MubWQjZmlyc3QtdHlwZXNjcmlwdC1wcm9maWxlNTpmYWxzZTA6"></a>

Authored scope: [docs/native-semantics.md#first-typescript-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-typescript-profile); profile: unspecified (authored profile absent).

### Declared target [github.com/Bitspark/bittype/analyze/files/typescript/provider/go#Provider.Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/provider/go/analyze.go#L57)

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | resolution-graph-preserved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| typescript-semantic-subset | generic-and-native-facts | native-identity-and-substitution-preserved | unspecified (authored profile absent) | [TestTypeScriptSemanticGenericsAndNativeFacts](https://github.com/Bitspark/bittype/blob/main/integration/go/typescript_semantics_test.go#L209) (external) |
| typescript-semantic-subset | rejected-and-unsupported-fixtures | no-partial-model | unspecified (authored profile absent) | [TestTypeScriptSemanticReports](https://github.com/Bitspark/bittype/blob/main/integration/go/typescript_semantics_test.go#L253) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.

<a id="cap-ODpwcm92aWRlcjgxOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9yZXNvbHZlL2ZpbGVzL3R5cGVzY3JpcHQvcHJvdmlkZXIvZ28jUHJvdmlkZXIuUmVzb2x2ZTExOTpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvY29udHJhY3RzL3Jlc29sdXRpb24vZ28uUHJvdmlkZXJbZ2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3N5bnRheC90eXBlc2NyaXB0L2dvLkxhbmd1YWdlXTczOmdpdGh1Yi5jb20vQml0c3BhcmsvYml0dHlwZS9yZXNvbHZlL2ZpbGVzL3R5cGVzY3JpcHQvcHJvdmlkZXIvZ28uUHJvdmlkZXI1NzpnaXRodWIuY29tL0JpdHNwYXJrL2JpdHR5cGUvc3ludGF4L3R5cGVzY3JpcHQvZ28jTGFuZ3VhZ2UxMDp0eXBlc2NyaXB0"></a>

## TypeScript / github.com/Bitspark/bittype/resolve/files/typescript/provider/go#Provider.Resolve

Contract: github.com/Bitspark/bittype/contracts/resolution/go.Provider&#91;github.com/Bitspark/bittype/syntax/typescript/go.Language&#93;.

Public operation: [github.com/Bitspark/bittype/resolve/files/typescript/provider/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/typescript/provider/go/resolver.go#L44). Receiver: github.com/Bitspark/bittype/resolve/files/typescript/provider/go.Provider.

Language marker: github.com/Bitspark/bittype/syntax/typescript/go.Language.

Construction alternatives: no provider collaborators.

### Declared target [github.com/Bitspark/bittype/resolve/files/typescript/provider/go#Provider.Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/typescript/provider/go/resolver.go#L44)

<a id="binding-NDg4Ok9EcHdjbTkyYVdSbGNqZ3hPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl5WlhOdmJIWmxMMlpwYkdWekwzUjVjR1Z6WTNKcGNIUXZjSEp2ZG1sa1pYSXZaMjhqVUhKdmRtbGtaWEl1VW1WemIyeDJaVEV4T1RwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdlkyOXVkSEpoWTNSekwzSmxjMjlzZFhScGIyNHZaMjh1VUhKdmRtbGtaWEpiWjJsMGFIVmlMbU52YlM5Q2FYUnpjR0Z5YXk5aWFYUjBlWEJsTDNONWJuUmhlQzkwZVhCbGMyTnlhWEIwTDJkdkxreGhibWQxWVdkbFhUY3pPbWRwZEdoMVlpNWpiMjB2UW1sMGMzQmhjbXN2WW1sMGRIbHdaUzl5WlhOdmJIWmxMMlpwYkdWekwzUjVjR1Z6WTNKcGNIUXZjSEp2ZG1sa1pYSXZaMjh1VUhKdmRtbGtaWEkxTnpwbmFYUm9kV0l1WTI5dEwwSnBkSE53WVhKckwySnBkSFI1Y0dVdmMzbHVkR0Y0TDNSNWNHVnpZM0pwY0hRdloyOGpUR0Z1WjNWaFoyVXhNRHAwZVhCbGMyTnlhWEIwODE6Z2l0aHViLmNvbS9CaXRzcGFyay9iaXR0eXBlL3Jlc29sdmUvZmlsZXMvdHlwZXNjcmlwdC9wcm92aWRlci9nbyNQcm92aWRlci5SZXNvbHZlNTpzY29wZTE6IzQ5OmRvY3MvbmF0aXZlLXNlbWFudGljcy5tZCNmaXJzdC10eXBlc2NyaXB0LXByb2ZpbGU1OmZhbHNlMDo"></a>

Authored scope: [docs/native-semantics.md#first-typescript-profile](https://github.com/Bitspark/bittype/blob/main/docs/native-semantics.md#first-typescript-profile); profile: unspecified (authored profile absent).

| Feature | Case | Authored expectation | Authored profile | Test declaration |
| --- | --- | --- | --- | --- |
| semantic-conformance | independent-provider-registration | captured-graph-resolved | unspecified (authored profile absent) | [TestSemanticProviders](https://github.com/Bitspark/bittype/blob/main/integration/go/horos_test.go#L50) (external) |
| typescript-semantic-subset | generic-and-native-facts | native-identity-and-substitution-preserved | unspecified (authored profile absent) | [TestTypeScriptSemanticGenericsAndNativeFacts](https://github.com/Bitspark/bittype/blob/main/integration/go/typescript_semantics_test.go#L209) (external) |
| typescript-semantic-subset | rejected-and-unsupported-fixtures | no-partial-model | unspecified (authored profile absent) | [TestTypeScriptSemanticReports](https://github.com/Bitspark/bittype/blob/main/integration/go/typescript_semantics_test.go#L253) (external) |

No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.
