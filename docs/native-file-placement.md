# Go, Rust and C++ declaration placement

**Proposal, 8 November 2026.** These examples specify the intended Files boundary.
TypeScript's [generated Files handoff](typescript-type-imports.md) is implemented;
Go, Rust and C++ Files formatters and their concrete placement APIs are not.
Rust currently chooses public visibility and C++ emits includes during generation;
those ownership defects remain open in [ACX-004](architecture-cases/0004-generation-file-syntax.md).
Go's package clause and C++'s namespace are fixed by target configuration during
generation; the same record lists them as pending instances.

## Inputs and ownership

Generate the logical declarations once:

```text
Page<T> = { items: list<T>; next?: text; }
Person = { name: text; }
Cells = Page<cells.Cell<Person>>
```

Page, Person and Cells are requested output declarations. The selected Cell
mapping names an external native public type. Generation owns the emitted
declarations and their intended uses, not the Cell package or its implementation.
The native model must retain that external destination without constructing a
file directive. External-symbol support for a future target must be represented
by that target's model and the shared output contract, not guessed from a name.

The consumer supplies one placement plus the selected dependency environment:

| Input | Proposed meaning |
| --- | --- |
| Owned declaration placement | Output file and native scope for each emitted declaration/helper; entries identify declarations, not just spellings |
| Public roots | Generation's exposed declaration sites; placement cannot make an inaccessible dependency public merely to obtain a valid file |
| External destination | Library-selected public package/crate/header symbol and mapped arguments |
| Dependency environment | Consumer-selected module/crate/include roots, actual dependency artifacts and relevant native configuration |

There is no second table telling Files which imports to emit. Files follows each
intended use under the placement and the language's scope rules. Changing only
Person's placement does not rerun semantic mapping or choose another Cell type.

The implemented shared handoff is in `contracts/format/go`:

```go
type GeneratedFilesInput struct {
    Output        generation.Output
    Configuration []semantic.Setting
    Limits        FilesLimits
}
type Provider[I, O any] interface {
    language.Provider
    Format(context.Context, I) (O, error)
}
```

A native Files implementation supplies
`Provider[GeneratedFilesRequest[L], source.Source[L]]`; the shared formatter
invokes it and admits publication through the configured contract. The full
[generated formatting contract](../contracts/format/go/generated.go) requires
evidence, printing and a check of the publication against its input. The concrete
Go/Rust/C++ settings and checks still need specification and implementation.
The following tables are proposed inputs/output, not shipped Go structs or APIs.

## Go

Assume the selected dependency exports `cell.Cell[T]` at import path
`example.org/cell`. A package import path and an output filename are different
identities; the consumer's Go module layout relates them.

| Placement | Go in the file containing Cells |
| --- | --- |
| Page/Cells in `pages/page.go`, package `example.org/app/pages`; Person in `people/person.go`, package `example.org/app/people` | Import `example.org/app/people` and `example.org/cell`; use `Page[cell.Cell[people.Person]]` |
| Person in `pages/person.go`, same package as Page | No import of Person's file; use `Page[cell.Cell[Person]]`; Cell import remains |
| Move Person from `pages/person.go` into `pages/page.go` | Same type use and package scope; still no Person import; Cell import remains |

The proposed split-package result includes:

```go
package pages

import (
    "example.org/app/people"
    "example.org/cell"
)

type Cells = Page[cell.Cell[people.Person]]
```

This excerpt omits the independently generated Page and presence helper.
Imports are file-scoped in Go, so another file in the same package needs its own
imports for its own uses. Moving a declaration across packages also changes its
native package identity and can introduce a dependency cycle. Files must check
visibility, identity correspondence and package cycles; it cannot silently export
an unexported type or treat same-package movement as a package import.

Package aliases must avoid names in the use's lexical scope. The Cell mapping
identifies its public type; Files chooses the local package alias. No Cell case
belongs in the Go target or shared orchestration.

## Rust

Assume the selected dependency crate is available as `bitspark_cell` and exports
`Cell<T>`. Cargo package selection/renaming belongs to the dependency environment;
it is not inferred from a source type's name.

| Placement | Rust in module `pages` |
| --- | --- |
| Page/Cells in `src/pages.rs`; Person in module `people` at `src/people.rs` | `use crate::people::Person;` and `use bitspark_cell::Cell;`; alias `Page<Cell<Person>>` |
| Person moves into `src/pages.rs` and module `pages` | Remove the Person import; retain the external Cell import |
| Person's bytes move into an inline `mod people` in the same physical file | Person still belongs to another module; use a suitable module path/import and satisfy its visibility |

The proposed split-module result includes:

```rust
use crate::people::Person;
use bitspark_cell::Cell;

pub type Cells = Page<Cell<Person>>;
```

This excerpt assumes Page and the crate's module declarations are supplied.
Files must construct `mod` declarations or inline modules consistently with the
placement, and realize public-root visibility under the crate/module rules.
Two files may contribute through one module arrangement; one file may contain
several modules. A common filename alone never makes Person directly visible.
Lifetimes, bounds and generic argument conventions remain part of the selected
native representation, not file placement.

## C++

Assume a C++20 header profile and the external type `bitspark::cell::Cell<T>`,
whose owning library identifies the public header `<bitspark/cell.hpp>`. Include
search roots and dependency versions come from the consuming build.

| Placement | Directives and names in `pages/page.hpp` |
| --- | --- |
| Page/Cells in namespace `app::pages`; Person in `model/person.hpp`, namespace `app::model` | Include the Person header and Cell header; alias `Page<bitspark::cell::Cell<app::model::Person>>` |
| Person moves into `pages/page.hpp`, retaining namespace `app::model` | Remove its header include; retain its qualified name and the Cell include |
| Person moves into the same header and namespace as Page | Remove its header include; Person can be referenced in that scope; retain the Cell include |

The proposed split-header result includes:

```cpp
#include "../model/person.hpp"
#include <bitspark/cell.hpp>

namespace app::pages {
using Cells = Page<bitspark::cell::Cell<app::model::Person>>;
}
```

This excerpt omits Page, standard-library dependencies and header guards.
The Files profile owns include directives, guards, declaration order and the
chosen relative/public-header policy. The type mapping still owns representation
choices such as vector versus another collection. A forward declaration is not
automatically enough: required completeness depends on the selected native use.
Files needs the target's rules or must reject unsupported dependency arrangements.
C++ modules would require a separately specified Files profile; they are not an
implicit alternative to the header profile.

## Required checks before claiming support

For each target, the same generated output must support split and co-located
Person placements without changing its selected Cell destination. Include the
same-package Go case, same-file/different-module Rust case and
same-header/different-namespace C++ case. Check collisions, visibility, missing
external destinations, dependency cycles/completeness and public roots.

The shared contract associates sites and evidence; language integrations establish
native scope and destination meaning. Go/Rust have bounded semantic tooling, but
that does not itself implement placement. C++ semantic evidence remains absent.
Successful printing or compilation alone does not establish that every intended
use bound to the right declaration. No runtime bindings or adapters are implied.
