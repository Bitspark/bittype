# Public model dependencies

Bittype's [go.mod](../go.mod) pins public Horos Go 0.10.1 and
Ontos 0.16.0. There are no local replacements, private module paths or
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

The ordinary `horos.kinds` prelude is bundled by Horos and exposed as a canonical
module value and import. Its exact identity is used for primitive operand
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
service deployment state. No validation or provenance workflow is a prerequisite
under the owner's current delivery rule, and unrun checks are not reported as
passing.
