# Public model dependencies

Bittype's [go.mod](../go.mod) pins public Horos Go 0.10.2,
Ontos 0.16.0 and the public Logos contract at
`v0.25.0`. There are no local replacements, private module paths or
compatibility aliases. Consumers upgrade together across the declaration revision.

The Horos frontend consumes revision 4: ordinary kind definitions, general kind
contracts, constructor-valued results, checked symbolic bodies and explicit
membership queries. The closed Type/Data interfaces and dedicated Data summaries
are removed. Ground Horos checking and Ontos value meanings are unchanged.

Horos owns kind meaning and primitive formation. Bittype owns source syntax,
AST projection and adaptation into its agnostic semantic graph. Generic trees,
contracts, shared semantics, the engine and native syntax stages retain their
model-independent import boundaries. The native pagination readers consume
Horos-owned facts under an explicit representation profile.

The [joint checker](../verification/logos/go/joint.go) invokes the typed checker;
Logos-aware verification interpretations and mapping laws also use `github.com/Bitspark/logos-contract/contract/go/typed` directly. Its
[public source revision](https://github.com/Bitspark/logos-contract/commit/eef2c26467bee408f36a8efee38ad99d89ea5b6b)
provides the check-only typed foundation. Agnostic association and native
projection packages have no Logos imports. No hosted service, credential or
private repository is needed to install Bittype.

The ordinary `horos.kinds` prelude is bundled by Horos and exposed as a canonical
module value and import. Its complete definitions are maintained in
[Horos source](https://github.com/Bitspark/horos/blob/main/core/go/declarations/standard/horos.kinds.horos),
lowered by the [Bittype module command](../tools/horos-module/README.md), and
shipped as a canonical asset. Horos decodes and admits that asset through the
ordinary module loader; its Go-authored definition builders are removed.
Downstream Horos builds do not run Bittype. Its exact identity is used for primitive operand
contracts. Source authors import it explicitly to refer to its names; there is no
implicit binding called Data or Type. An absent bundled prelude file does not
cause Bittype to fabricate a source location or physical file dependency.

Ontos 0.16.0 supplies the optional mixed value carrier. Static Bittype consumes
ground module descriptions and adds no production dependency on Ontos/Wire,
Bitwire or Bitruntime. Hydration and live capability behavior stay with their
existing owners.

Existing consumer fixtures use the public declaration API: obtain a Program-owned
value, apply it to admitted arguments, and materialize its descriptor or request
ground-descriptor extraction. Package pins describe source dependencies, not
service deployment state. No extra administrative validation or provenance workflow is an implicit
prerequisite. Relevant software tests and typechecking still apply under the
owner's clarified rule; unrun checks are never reported as passing.
