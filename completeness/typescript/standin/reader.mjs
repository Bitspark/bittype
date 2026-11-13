// PROTOTYPE (bittype#196, row 31). A second reader in another process and
// another language: a stand-in for a TypeScript binding backend. It is run
// with node from a Go test. It reads the serialized plan and the source
// snapshot, and nothing else of the planning process, and converts values.
//
//   node reader.mjs --plan plan.json [--snapshot snapshot.json] --vectors vectors.json
//
// It prints one JSON object: the results of the vectors and a memo key for
// each requested root, or the failure that stopped it, with its class.
//
// What it shows and what it does not is in PROTOTYPE.md and on the design
// page under "What a stand-in proves".

import { createHash } from "node:crypto";
import { readFileSync } from "node:fs";
import path from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";
import { data, decode, fromHex } from "./ontos.mjs";

// What this reader knows of each version axis that it can check.
const KNOWN = {
  envelope: ["bittype.plan/0-unstable"],
  semantics: ["horos/declarations-4"],
  runtime: ["bittype-standin/1"],
  facetContract: "1",
};
const VALUE_FACTORY = "bittype.completeness/typescript/value-factory";
const FIELD_FACTORY = "bittype.completeness/typescript/field-state-factory";

class Failure extends Error {
  constructor(kind, message) {
    super(message);
    this.kind = kind;
  }
}

const key = (r) => `${r.Identity.Publisher}/${r.Identity.Key}@${r.Meaning.Identity}/${r.Meaning.Revision}#${r.Representation}`;
const same = (a, b) => JSON.stringify(a ?? []) === JSON.stringify(b ?? []);

function receipt(schema, frozen) {
  const h = createHash("sha256");
  h.update(`bittype.receipt/0\0${schema.Identity}\0${schema.Version}\0`);
  h.update(frozen);
  return h.digest("hex");
}

async function read(options) {
  const plan = JSON.parse(readFileSync(options.plan, "utf8"));
  if (!KNOWN.envelope.includes(plan.envelope)) throw new Failure("version", `envelope protocol "${plan.envelope}"`);
  if (!KNOWN.semantics.includes(plan.source.semantics)) throw new Failure("version", `source semantics "${plan.source.semantics}"`);

  // The selected contributions. Each names its factory through the facet that
  // crossed into the plan. The stand-in loads the factory by module and export.
  const here = path.dirname(fileURLToPath(import.meta.url));
  const contributions = new Map();
  for (const entry of plan.contributions) {
    const loaded = {};
    for (const facet of entry.facets ?? []) {
      const slot = facet.kind.identity === VALUE_FACTORY ? "value" : facet.kind.identity === FIELD_FACTORY ? "field" : null;
      if (!slot) continue;
      if (facet.kind.version !== KNOWN.facetContract || facet.mediaType !== "application/json") {
        throw new Failure("version", `facet contract ${facet.kind.identity} version "${facet.kind.version}"`);
      }
      const { module, export: name, runtime } = facet.value;
      if (!KNOWN.runtime.includes(runtime)) throw new Failure("version", `generated-code and runtime compatibility "${runtime}"`);
      const factory = (await import(pathToFileURL(path.join(here, module)).href))[name];
      if (!factory) throw new Failure("linking", `the module ${module} has no factory ${name}`);
      if (!factory.schemas.some((s) => s.identity === entry.schema.Identity && s.version === entry.schema.Version)) {
        throw new Failure("version", `payload schema ${entry.schema.Identity} version "${entry.schema.Version}"`);
      }
      const implemented = `${entry.reference.Meaning.Identity}/${entry.reference.Meaning.Revision}#${entry.reference.Representation}`;
      if (!factory.implements.includes(implemented)) throw new Failure("version", `representation semantics ${implemented}`);
      loaded[slot] = factory;
    }
    contributions.set(key(entry.reference), loaded);
  }

  // The source snapshot: the second input. Without it the reader cannot name
  // a canonical field, and it says so. A snapshot of another pin does not
  // belong to this plan.
  if (!options.snapshot) throw new Failure("incomplete", "a reader needs the source snapshot, and none was given");
  const snapshot = JSON.parse(readFileSync(options.snapshot, "utf8"));
  if (snapshot.semantics !== plan.source.semantics) throw new Failure("linking", "the snapshot is of another source semantics");
  for (const symbol of plan.symbols) {
    if (symbol.id.startsWith("source:") && !snapshot.declarations[symbol.id.slice(7)]) {
      throw new Failure("linking", `the snapshot has no declaration ${symbol.id.slice(7)}: it is of another module pin`);
    }
  }

  const nodes = new Map(plan.nodes.map((n) => [n.id, n]));

  function child(descriptor, link) {
    if (!descriptor) return undefined;
    if (descriptor.kind === "record") {
      const field = descriptor.fields.find((f) => same(f.position, link.position));
      if (field) return field.type;
    }
    if (descriptor.kind === "list" && descriptor.element) return descriptor.element.type;
    return descriptor;
  }

  function make(node, scope, descriptor) {
    const loaded = contributions.get(key(node.contribution)) ?? {};
    const factory = node.kind === "field" ? loaded.field : loaded.value;
    if (!factory) throw new Failure("missing facet", `no factory for ${key(node.contribution)}`);
    const frozen = fromHex(node.payload);
    if (receipt(node.schema, frozen) !== node.receipt) throw new Failure("invalid", `the receipt of node ${node.id} was not issued for its payload`);
    const links = (node.links ?? []).map((link) => {
      const built = build(link.node, scope, node.kind === "declaration" ? descriptor : child(descriptor, link));
      return link.kind === "field" ? { kind: "field", stated: link.stated, state: built } : { kind: "child", stated: link.stated, adapter: built };
    });
    const canonical = (i) => {
      const field = descriptor?.kind === "record" && descriptor.fields.find((f) => same(f.position, node.links[i].position));
      if (!field) throw new Failure("linking", `the snapshot has no child at position ${JSON.stringify(node.links[i].position)}`);
      return field.name;
    };
    return factory.make({ payload: data(decode(frozen)), links, canonical, members: node.members ?? [] });
  }

  // A planned declaration is built with one adapter for each binder: the
  // adapter arguments of the application. Open generics stay open.
  function declaration(node, adapters) {
    return make(node, { adapters }, snapshot.declarations[node.source]);
  }

  function build(id, scope, descriptor) {
    const node = nodes.get(id);
    if (!node) throw new Failure("incomplete", `no node ${id}`);
    switch (node.kind) {
      case "binder":
        return scope.adapters[node.binder.index];
      case "application": {
        const callee = nodes.get(node.application.declaration);
        const adapters = [];
        for (const a of node.application.adapters ?? []) adapters[a.binder] = build(a.node, scope, descriptor?.arguments?.[a.binder]);
        for (let i = 0; i < (callee.binders ?? 0); i++) {
          if (!adapters[i]) throw new Failure("incomplete", `the application ${node.source} has no adapter argument for binder ${i}: the plan was not asked for bindings`);
        }
        return declaration(callee, adapters);
      }
      case "use":
      case "field":
        return make(node, scope, descriptor);
      default:
        throw new Failure("unsupported", `a ${node.kind} node has no conversion`);
    }
  }

  // A vector's native value may hold a present undefined, which JSON cannot
  // say: {"$undefined": true} stands for it.
  const untag = (v) => {
    if (Array.isArray(v)) return v.map(untag);
    if (v !== null && typeof v === "object") {
      if (v.$undefined === true) return undefined;
      return Object.fromEntries(Object.entries(v).map(([k, x]) => [k, untag(x)]));
    }
    return v;
  };

  const roots = new Map();
  const memo = {};
  for (const id of plan.roots) {
    const symbol = plan.symbols.find((s) => s.id === id);
    const node = nodes.get(symbol.node);
    const name = id.slice(id.lastIndexOf("/") + 1);
    // What a memo token is derived from is in the plan: the node of the use,
    // which is identified by its decisions, its application and the nodes
    // selected for its dependencies.
    memo[name] = createHash("sha256").update("memo\0" + node.links[0].node).digest("hex").slice(0, 16);
    if (!node.binders) roots.set(name, node);
  }

  const results = [];
  for (const vector of JSON.parse(readFileSync(options.vectors, "utf8"))) {
    try {
      const node = roots.get(vector.root);
      if (!node) throw new Failure("invalid", `no root ${vector.root}`);
      const adapter = declaration(node, []);
      if (vector.direction === "encode") {
        results.push({ root: vector.root, value: adapter.encode(untag(vector.native)) });
      } else {
        const native = adapter.decode(vector.value);
        results.push({ root: vector.root, native, value: adapter.encode(native) });
      }
    } catch (e) {
      results.push({ root: vector.root, error: { kind: e.kind ?? "error", message: e.message } });
    }
  }
  return { results, memo };
}

const options = {};
for (let i = 2; i + 1 < process.argv.length; i += 2) options[process.argv[i].replace(/^--/, "")] = process.argv[i + 1];
read(options).then(
  (out) => process.stdout.write(JSON.stringify(out)),
  (e) => process.stdout.write(JSON.stringify({ failure: { kind: e.kind ?? "error", message: e.message } })),
);
