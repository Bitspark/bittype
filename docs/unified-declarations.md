# Unified Horos frontend delivery

**Historical delivery, 4 November 2026.** The unification described here was
implemented with public Horos 0.9.0; exact-source gates and installed-consumer
checks qualified that revision. The subsequent [revision-3 kind system](kinds.md)
replaces its bound-only parameter representation with proper/constructor
interfaces and result claims. This page preserves the earlier delivery evidence;
it is not the current canonical grammar or native-generation status.

The [Horos revision-2 work](https://github.com/Bitspark/horos/issues/27) replaced
both prior declaration profiles with one module language and explicit Type/Data
bounds. That work migrated Bittype and maintained consumers without adding
native generators or runtime bindings. Subsequent [pagination mappings](semantic-generation.md)
implement ordinary native generation; runtime bindings remain deferred.

The retained families are syntax/horos, check/horos, resolve/horos, analyze/horos
and provider/horos. The horoswire families and superseded profiles were removed
without aliases. The one Horos marker used new syntax/resolution/analysis profiles.
Syntax admission remained independent of semantic resolution. That revision
projected parameter bounds into Horos AST nodes and provider-owned semantic
terms; the agnostic engine did not interpret them. Analysis consumed explicit
resolution receipts through the shared semantic infrastructure. Revision 3
preserves that boundary while replacing the representation and profile IDs.

Page<T> is authored once with an unconstrained Type slot. DataPage<D: Data> is an
ordinary constrained alias. DataPage<Person> succeeds; Page<Cell<Person>> also
succeeds when Cell is an ordinary capability-bearing library declaration;
DataPage<Cell<Person>> fails. Bound forwarding, nested applications, imported
definitions and preserved caller scopes must be tested. A bound cannot be
inferred from a library name or merely from the presence of data arguments.

The migration changed old ground slots to explicit Data and old mixed slots to
Type, repinned canonical documents, and preserved literal payloads and original
generic source. It replaced the provider matrix's two columns with one Horos
column and tightened architecture dependency allowances to that frontend.

The maintained platform checkout inventory found no additional direct declaration
API users outside Horos and Bittype. This migration covers the concrete frontend,
canonical fixtures, provider/architecture tests and isolated installed consumer.
Ground-only descriptor consumers have no declaration-format rewrite to perform.

The installed-consumer gate for this delivery used a verified public Horos release
and Ontos 0.16.0, without local replace directives. Required validation included
contract and architecture checks, syntax/semantics/receipt tests, bounds and pagination,
race, fuzz, native compiler fixtures and docs. Run portable checks on hosted
runners; accept receipts only for the exact committed source.

The independent template-lifting removal and Horos 0.8.0 dependency delivery
are complete. This change starts from verified Bittype main 972fefe82ccb and
Horos main 690bbe54e1c8. It now consumes Horos 0.9.0 from source e2896317a1f4,
with no compatibility readers, old namespace aliases or executable archives.
