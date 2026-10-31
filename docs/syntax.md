# Target-language AST declarations

**Status, 31 October 2026:** TypeScript, Go, Rust and C++ syntax vocabularies
are implemented in Go on the [generic AST](ast.md). Each package supplies
`Language`, `AST`, `Node`, `Kind`, and constants for kinds, roles and attributes.

The shared LanguageProvider interface, native text providers, parsers,
AST-to-CST formatters, checks and generic printer are implemented. Canonical AST,
actual language mappings and public builders remain deferred. The syntax engine
now provides [explicit registration and route execution](engine.md).
See the [native text design](parsing-formatting.md).

The exported carrier fields remain unchecked; the separate language checker
validates supported type syntax. These are not complete compiler ASTs.

## Identity and registration

[contracts/language/go](../contracts/language/go/language.go) declares the agnostic interface:

```go
type ID string

type LanguageProvider interface {
    ID() ID
    Name() string
}
```

| Provider | ID | Display name |
| --- | --- | --- |
| [typescript.Language](../syntax/typescript/go/nodes.go) | `typescript` | TypeScript |
| [golang.Language](../syntax/golang/go/nodes.go) | `go` | Go |
| [rust.Language](../syntax/rust/go/nodes.go) | `rust` | Rust |
| [cpp.Language](../syntax/cpp/go/nodes.go) | `cpp` | C++ |

The zero-value providers implement this interface. The methods return fixed
metadata and perform no IO, discovery or registration. IDs are exact,
case-sensitive language-family keys; names are display text. They are not
toolchain versions, schema versions or structural fingerprints.

**Owner correction, 31 October 2026:** there will be a central, agnostic language
AST registry. The earlier statement that there would be no language registry
is superseded. The syntax engine implements it through typed/profile descriptors
and explicit native-provider adaptation, without concrete syntax imports.
Duplicate keys and directed edges are rejected. Configured descriptors require
construction; the native markers above keep usable zero values. There are no
hidden defaults or automatic init-time registrations.

Each provider's `Language` is also a distinct named Go marker:

```go
type AST = ast.AST[Language]
type Node = ast.Node[Language]
type Kind = ast.Kind[Language]
```

These aliases use the existing carrier, not parallel tree storage.
The generic AST remains unconstrained by the identity interface. An independent
consumer can implement the interface and use its own marker without editing
either foundation. Equal ID strings do not establish Go type identity or
grammar compatibility, and do not permit converting an `AST[L]` into another
language's AST.

## Shared notation and ownership

The tables below specify the intended data shapes, not executable checkers.
Kind names omit the exported `Kind` prefix. Attribute and role names show
their actual string values; the packages export their corresponding
`Attribute...` and `Role...` constants.

- A role without a suffix means exactly one child; `?` means zero or one,
  `*` means zero or more, and `+` means one or more. Repeated children retain
  order. A type expression means one of that language's type-expression kinds.
- Attributes are unique by name. Names, qualified-name segments and package
  qualifiers are individual identifiers unless a row explicitly permits other
  text. Compound types and generic arguments belong in child nodes.
- Boolean attributes are the text `true` or `false`; absence means false.
  Other defaults are stated below. No arbitrary source-code attribute exists.
- An optional reference with role `target` on a named type use, name segment
  or lifetime use can retain its resolved declaration/import/binder. It is a
  qualified `ast.NodeRef`, not a containment edge. Omitting it means no resolved
  link is supplied, not an invented target. Resolution is not implemented.
- Generic parameter declarations and their uses remain separate nodes.
  Same-spelled parameters in different scopes can point to different binders.
  Applications retain their constructor/path and ordered arguments.
- All nodes can carry the generic `Origin`. Neither origin nor a target address
  proves semantic identity, ownership or compatible behavior.

Only the listed forms and attributes belong to this vocabulary. Separate native
checkers enforce structure, roles, lexical atoms and supported contextual rules.
Complete compiler/type checking and cross-tree resolution remain outside them.
Parsing retains original source; formatting applies the documented
[normalizations and layout](parsing-formatting.md).

Function-type and method-signature nodes describe native language syntax.
They contain no executable bodies. They do not give the Function library,
Stream, Cell or another behavioral contract a privileged kind in the generic
AST, identity interface, pipeline or registry.

## TypeScript

Source vocabulary: [nodes.go](../syntax/typescript/go/nodes.go).
Native reference: [TypeScript object types](https://www.typescriptlang.org/docs/handbook/2/objects.html).

| Kind | Attributes | Children |
| --- | --- | --- |
| SourceFile | — | `import*` Import; `declaration*` TypeAlias/Interface |
| Import | `module`: decoded module specifier | `binding+` ImportBinding, or one NamespaceImport |
| ImportBinding | `name`: imported identifier; optional `alias`: local identifier | — |
| NamespaceImport | `name`: local namespace identifier | — |
| TypeAlias | `name`; `exported` | `type-parameter*`; `type` |
| Interface | `name`; `exported` | `type-parameter*`; `base*` named/applied types; `member*` |
| TypeParameter | `name` | `constraint?`; `default?` types |
| TypeReference | `name` | — |
| QualifiedType | `name`: final identifier | `qualifier` TypeReference/QualifiedType |
| TypeApplication | — | `constructor` TypeReference/QualifiedType; `argument+` types |
| KeywordType | `name`: native type keyword | — |
| LiteralType | `literal-kind`; `text`: exact literal spelling | — |
| ArrayType | — | `element` type |
| TupleType | — | `element*` TupleElement |
| TupleElement | optional `name`; `optional`; `rest` | `type` |
| ObjectType | — | `member*` Property/Method/CallSignature/IndexSignature |
| UnionType, IntersectionType | — | `type+` (at least two) |
| FunctionType | — | `type-parameter*`; `parameter*`; `result` type |
| Parameter | `name`; `optional`; `rest` | `type` |
| Property | `name`: exact property text; `optional`; `readonly` | `type` |
| Method | `name`: exact property text; `optional` | `signature` FunctionType |
| CallSignature | — | `signature` FunctionType |
| IndexSignature | `readonly` | `parameter` Parameter; `type` |
| TypeOperator | `operator`: `keyof`, `readonly` or `unique` | `type` |
| IndexedAccessType | — | `object` type; `index` type |

Imports are type-only imports. An absent binding alias means its imported name.
Interface members use the same member forms as ObjectType. Property text is
not an expression: a later formatter can quote a non-identifier property name.
Computed property keys are outside this subset.

KeywordType admits `any`, `unknown`, `never`, `void`, `undefined`,
`boolean`, `number`, `bigint`, `string`, `symbol`, `object` and `this`.
Literal kinds are `string`, `number`, `bigint`, `boolean` and `null`.
Their text is one matching native literal (including its string quotes or
numeric sign), never a template interpolation or arbitrary expression.
Checks cover rest/optional ordering, tuple naming consistency and direct operator
applicability. Semantic index-key and imported-type constraints remain compiler
obligations. Numeric property keys are admitted only as canonical safe decimal
integers; other spellings need numeric-to-property-text decoding.

Conditional/mapped/inferred/template-literal types, value declarations,
class bodies, decorators and executable expressions are outside this slice.

## Go

Source vocabulary: [nodes.go](../syntax/golang/go/nodes.go).
Native reference: [Go types](https://go.dev/ref/spec#Types).

| Kind | Attributes | Children |
| --- | --- | --- |
| File | `package`: package identifier | `import*`; `declaration*` TypeDeclaration/TypeAlias |
| Import | `path`: decoded import path; optional `alias` | — |
| TypeDeclaration, TypeAlias | `name` | `type-parameter*`; `type` |
| TypeParameter | `name` | `constraint` type |
| NamedType | `name`: local or predeclared identifier | — |
| QualifiedType | `qualifier`: package identifier; `name` | — |
| TypeApplication | — | `constructor` NamedType/QualifiedType; `argument+` types |
| StructType | — | `field*` Field/EmbeddedField |
| Field | `name`; optional `tag`: decoded tag text | `type` |
| EmbeddedField | optional `tag` | `type` |
| InterfaceType | — | `member*` Method/EmbeddedType |
| Method | `name` | `signature` FunctionType |
| EmbeddedType | — | `type`: type or constraint TypeUnion |
| FunctionType | — | `parameter*` Parameter; `result*` Parameter |
| Parameter | optional `name`; `variadic` | `type` |
| PointerType, SliceType | — | `element` type |
| ArrayType | `length`: nonnegative decimal integer | `element` type |
| MapType | — | `key` type; `element` type |
| ChannelType | `direction`: `both` (default), `send` or `receive` | `element` type |
| TypeUnion | — | `term+`: types or TypeApproximation |
| TypeApproximation | — | `type` |

Names such as `any`, `comparable` and `string` are NamedType uses, not
canonical primitive kinds. An import alias may also be `_` or `.`.
Defined types and aliases remain different kinds. Each field or parameter
node represents one field/parameter, rather than a grouped name list.

A variadic parameter represents `...T` with `T` as its type child. Only the
last input parameter may be variadic; results may not be. Interface method
signatures cannot introduce method-specific type parameters. Constraint unions
and approximations belong to constraint positions, not arbitrary value unions.

Array lengths are literal integers in this initial vocabulary. Constant
expressions, constant/variable/function declarations, method implementations
and statements are outside this slice.

## Rust

Source vocabulary: [nodes.go](../syntax/rust/go/nodes.go).
Native reference: [Rust types](https://doc.rust-lang.org/reference/types.html).

| Kind | Attributes | Children |
| --- | --- | --- |
| File | — | `item*` Module/Use/TypeAlias/Struct/TupleStruct/Enum/Trait |
| Module | `name`; `visibility` | `item*` |
| Use | optional `alias` | `path` PathType without generic arguments |
| TypeAlias | `name`; `visibility` | `generic-parameter*`; `where*`; `type` |
| Struct | `name`; `visibility` | `generic-parameter*`; `where*`; `field*` Field |
| TupleStruct | `name`; `visibility` | `generic-parameter*`; `where*`; `field*` TupleField |
| Enum | `name`; `visibility` | `generic-parameter*`; `where*`; `variant*` |
| UnitVariant | `name` | — |
| TupleVariant | `name` | `field*` TupleField |
| RecordVariant | `name` | `field*` Field |
| Field | `name`; `visibility` | `type` |
| TupleField | `visibility` | `type` |
| Trait | `name`; `visibility` | `generic-parameter*`; `bound*`; `where*`; `member*` MethodSignature/AssociatedType |
| AssociatedType | `name` | `bound*` |
| MethodSignature | `name` | `generic-parameter*`; `where*`; `receiver?`; `parameter*`; `result?` type |
| Receiver | `receiver-form`: `value`, `shared` or `mutable` | `lifetime?` LifetimeReference for a borrowed receiver |
| Parameter | `name` for a method; optional in FunctionType | `type` |
| TypeParameter | `name` | `bound*`; `default?` type |
| LifetimeParameter | `name` without apostrophe | `bound*` LifetimeReference |
| LifetimeReference | `name` without apostrophe | — |
| WherePredicate | — | `subject` type/lifetime; `bound+` |
| PathType | `absolute` | `segment+` PathSegment |
| PathSegment | `name` | `argument*` types/LifetimeReference |
| PrimitiveType | `name`: native primitive | — |
| ReferenceType | `mutable` | `element` type; `lifetime?` |
| RawPointerType | `mutable` | `element` type |
| SliceType | — | `element` type |
| ArrayType | `length`: nonnegative decimal integer | `element` type |
| TupleType | — | `element*` types |
| NeverType | — | — |
| FunctionType | `unsafe`; optional `abi`: decoded ABI name | `parameter*`; `result?` type |
| TraitObjectType | — | `bound+` |

Generic parameter children are TypeParameter or LifetimeParameter. Bounds are
PathType or LifetimeReference; lifetime binders only have lifetime bounds.
`visibility` is `private` (default), `public` or `crate`. Enum variant
fields have inherited visibility, so only the default is admitted there.
Primitive names are `bool`, `char`, `str`, `f32`, `f64`,
`i8/i16/i32/i64/i128/isize` and `u8/u16/u32/u64/u128/usize`.
Unit is an empty TupleType. A missing result denotes unit.
Reference mutability and raw-pointer mutability are explicit; a false raw
pointer flag denotes `*const`. Path segments may use native `crate`,
`self`, `super` and `Self` names where Rust permits them.

FunctionType is a bare function-pointer signature. Trait method receivers are
separate nodes, never ordinary first parameters. TraitObjectType represents
explicit `dyn` bounds. Default generic types belong only where Rust allows
them, not on method-specific parameters.

Unit structs have no distinct AST form and are rejected.
Const generics, higher-ranked binders, associated-type argument bindings,
qualified associated-type projections, generic associated types, `impl Trait`,
macros, attributes, impl blocks and executable bodies are outside this slice.

## C++

Source vocabulary: [nodes.go](../syntax/cpp/go/nodes.go).
Native reference: [C++ declaration grammar](https://eel.is/c++draft/dcl.pre).

| Kind | Attributes | Children |
| --- | --- | --- |
| TranslationUnit | — | `include*`; `declaration*` Namespace/Record/TypeAlias/Enum |
| Include | `path`: header text without delimiters; `system` | — |
| Namespace | `name` | `declaration*` |
| Record | `name`; `record-key`: `struct` or `class`; `forward`; optional `access` when nested | `template-parameter*`; `base*`; `member*` |
| TypeAlias | `name`; optional `access` when nested | `template-parameter*`; `type` |
| Enum | `name`; `scoped`; `forward`; optional `access` when nested | `underlying?` type; `enumerator*` |
| Enumerator | `name`; optional `value`: signed decimal integer | — |
| Base | optional `access`; `virtual` | `type` NamedType |
| Field | `name`; optional `access` | `type` |
| Method | `name`; optional `access`; `const`; `volatile`; `virtual`; `pure`; `noexcept`; `ref-qualifier`; `variadic` | `parameter*`; `result` type |
| Destructor | optional `access`; `virtual`; `pure`; `noexcept` | — |
| Parameter | optional `name` | `type` |
| TypeParameter | `name`; `parameter-key`: `typename` (default) or `class` | `default?` type |
| NamedType | `global`; `typename` | `segment+` NameSegment |
| NameSegment | `name`; `template` | `argument*` types |
| BuiltinType | `name`: native builtin spelling | — |
| PointerType, LValueReferenceType, RValueReferenceType | — | `element` type |
| CVQualifiedType | `const`; `volatile` (at least one true) | `element` type |
| ArrayType | `length`: positive decimal integer | `element` type |
| FunctionType | `noexcept`; `variadic` | `parameter*`; `result` type |
| MemberPointerType | — | `owner` NamedType; `element` type |

Record members are Field, Method, Destructor, TypeAlias, Record or Enum.
Access is `public`, `protected` or `private`; absent member/base access uses
the enclosing record's native default. A true system include uses angle
delimiters; false uses quotes. Namespace nesting uses child nodes. The global
flag means leading `::`; name segments retain their own template arguments.
The `typename` and `template` flags describe dependent-name disambiguators.
A forward Record has no bases or members; a forward Enum has no enumerators.
Unscoped enum forward declarations require an underlying type. Forward forms
preserve the distinction between naming a type and defining an empty type.

Builtin spellings cover `void`, `bool`, the native character types, signed
and unsigned integer type specifiers, and `float`, `double`, `long double`.
Library types such as `std::int64_t` use NamedType; this AST assigns no
fixed-width meaning to native `long`.
`ref-qualifier` is `none` (default), `lvalue` or `rvalue`.
Pure declarations must be virtual. A destructor takes its name from its record;
this vocabulary supplies a declaration only, not a destructor implementation.
The noexcept flag requests an unconditional noexcept specification; omission
leaves the native default rather than asserting a particular exception outcome.

Pointer and qualifier nesting preserves declarator structure: a pointer to
const T differs from a const pointer to T. NamedType represents applications
such as `std::vector<T>` using separate name segments and argument nodes.
There is no source-string fallback for a difficult declarator.

Non-type/template-template parameters, parameter packs, specialization,
constraints/requires expressions, computed array bounds, initializers,
constructors, operator declarations, macros and function bodies are outside
this slice. Qualified member-function-pointer types are rejected because native
FunctionType has no method CV/ref fields. The maintained native text gates
provide separately scoped compiler qualification.

## Evidence and limits

[Consumer fixtures](../integration/go/testdata/syntax/typescript.go) and the
equivalent [Go](../integration/go/testdata/syntax/golang.go),
[Rust](../integration/go/testdata/syntax/rust.go) and
[C++](../integration/go/testdata/syntax/cpp.go) fixtures construct native
`Page<T>`, a `Person`, an applied `PeoplePage`, and a generic reader
signature. They keep same-spelled T binders distinct through target links.
They are independently authored examples, not canonical-to-native mappings.

The [tests](../integration/go/syntax_test.go) exercise provider metadata through
the common interface, type-check those consumers, reject mixed target ASTs,
and check every pair of public marker/AST/node/kind types for non-assignability.
A [separately owned provider](../integration/go/testdata/syntax/extension.go)
implements the same interface without any central language list.
The [architecture checks](../integration/go/architecture_test.go) enforce that
syntax packages depend only on the generic AST and language identity contracts
plus the generic CST carrier and standard-library utilities. The three foundations do not depend
on each other or on concrete languages.

The original contract tests establish Go API usability and isolation. Additional
native text tests cover checking, retained source and independent AST formatting.
Native compiler fixtures run separately on remote runners. These claims do not
establish executable behavior, registry operation or canonical equivalence.
The pipeline contract is unchanged; AST adds a borrowed view interface.
