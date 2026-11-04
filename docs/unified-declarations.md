# Unified Horos frontend

**Owner decision, 4 November 2026. Implemented in the candidate; final public dependency qualification is pending.**
The [Horos revision-2 work](https://github.com/Bitspark/horos/issues/27) replaces
both prior declaration profiles with one module language and explicit Type/Data
bounds. This work migrates Bittype and maintained consumers; it implements no
native generator or runtime binding.

Keep syntax/horos, check/horos, resolve/horos, analyze/horos and provider/horos.
Delete the horoswire families and superseded profiles without aliases. The one
Horos marker uses new syntax/resolution/analysis profiles for the new contract.
Syntax admission remains independent of semantic resolution. Parameter bounds
are source facts on Horos AST nodes and provider-owned semantic terms; the
agnostic engine does not interpret them. Analysis consumes explicit resolution
receipts through the existing shared semantic infrastructure.

Page<T> is authored once with an unconstrained Type slot. DataPage<D: Data> is an
ordinary constrained alias. DataPage<Person> succeeds; Page<Cell<Person>> also
succeeds when Cell is an ordinary capability-bearing library declaration;
DataPage<Cell<Person>> fails. Bound forwarding, nested applications, imported
definitions and preserved caller scopes must be tested. A bound cannot be
inferred from a library name or merely from the presence of data arguments.

Migrate old ground slots to explicit Data and old mixed slots to Type. Repin
canonical documents after migration; preserve literal payloads and original
generic source. Replace the two columns in the provider matrix with one Horos
column and tighten architecture dependency allowances to that frontend.

The installed-consumer gate must use a verified public Horos release and Ontos
0.16.0, without local replace directives. Required validation includes contract
and architecture checks, syntax/semantics/receipt tests, bounds and pagination,
race, fuzz, native compiler fixtures and docs. Run portable checks on hosted
runners and accept receipts only for the exact committed source.

The independent template-lifting removal and Horos 0.8.0 dependency delivery
are complete. This change starts from verified Bittype main 972fefe82ccb and
Horos main 690bbe54e1c8. Publish/consume the next breaking Horos version, with
no compatibility readers, old namespace aliases or executable archives.
