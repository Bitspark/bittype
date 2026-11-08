// @ts-check
// Private compiler helper. Input is captured text plus independently checked AST
// correspondence. Inputs may be generated and always belong to an explicit capture.
const fs = require("node:fs");
const path = require("node:path");
/** @type {typeof import("typescript")} */
const ts = require(process.argv[1]);
const key = name => ({ Owner: "typescript", Name: name });
const site = (unit, node) => ({ Tree: unit.Tree, Node: node.ID });
const attr = (node, name) => (node.Attributes || []).find(a => a.Name === name)?.Text || "";
const children = (node, role) => (node.Children || []).filter(c => c.Role === role).map(c => c.Node);
class Failure extends Error {
  constructor(status, message, at = undefined) { super(message); this.status = status; this.at = at; }
}
function diagnostic(status, message, at) {
  return { Code: key(status), Subject: key(status === "invalid" ? "captured-type-validity" : "captured-type-coverage"), Message: message, ...(at ? { At: at } : {}) };
}
function unsupported(message, at = undefined) { throw new Failure("unsupported", message, at); }
function unwrap(node) { while (ts.isParenthesizedTypeNode(node)) node = node.type; return node; }

function main(request) {
  if (request.Version !== 1 || ts.version !== "5.9.3") throw new Failure("environment", "Expected helper protocol 1 and TypeScript 5.9.3");
  const supported = new Set(["source-file", "type-alias", "interface", "type-parameter", "type-reference", "type-application", "keyword-type", "literal-type", "array-type", "tuple-type", "tuple-element", "object-type", "union-type", "intersection-type", "property", "type-operator"]);
  if (request.NamedImports) { supported.add("import"); supported.add("import-binding"); }
  for (const unit of request.Units) for (const node of unit.Nodes) {
    if (!supported.has(node.Kind) || node.Kind === "type-operator" && attr(node, "operator") !== "readonly" || node.Kind === "keyword-type" && attr(node, "name") === "this") unsupported(`Syntax kind ${node.Kind} is outside captured-types-1`, site(unit, node));
    if (children(node, "default").length) unsupported("Generic defaults are outside captured-types-1", site(unit, node));
    if (children(node, "type-parameter").length > 16) unsupported("At most 16 generic parameters are supported", site(unit, node));
  }
  const options = { target: ts.ScriptTarget.ES2020, module: ts.ModuleKind.ESNext, strict: true, exactOptionalPropertyTypes: true, noUncheckedIndexedAccess: true, noEmit: true, types: [], typeRoots: [], lib: ["lib.es2020.d.ts"], skipLibCheck: false };
  const sourceFiles = new Map();
  const unitByFile = new Map();
  for (const unit of request.Units) {
    const name = "/input/" + unit.Path;
    sourceFiles.set(name, ts.createSourceFile(name, unit.Text, options.target, true, ts.ScriptKind.TS));
    unitByFile.set(name, unit);
  }
  // Package declarations are receipt-bound inputs. They have their own virtual
  // filenames and never enter the inventory of owned units that will be printed.
  const packages = new Map(), packageByFile = new Map();
  for (const [i, pkg] of (request.Packages || []).entries()) {
    const name = "/packages/" + i + "/entry.d.ts";
    const sf = ts.createSourceFile(name, pkg.Declaration, options.target, true, ts.ScriptKind.TS);
    // Use the pinned compiler's parsed directives. A comment inside a literal
    // is ordinary type data and must not be treated as a diagnostic control.
    if ("checkJsDirective" in sf && sf.checkJsDirective || "commentDirectives" in sf && Array.isArray(sf.commentDirectives) && sf.commentDirectives.length)
      unsupported("Package declarations must preserve native diagnostics");
    if (!ts.isExternalModule(sf) || sf.hasNoDefaultLib || sf.referencedFiles.length || sf.typeReferenceDirectives.length || sf.libReferenceDirectives.length || sf.statements.some(st => !ts.isInterfaceDeclaration(st) && !ts.isTypeAliasDeclaration(st) && !ts.isImportDeclaration(st)))
      unsupported("Package input must be a declaration module of interfaces, aliases and named type imports");
    if (sf.statements.some(st => ts.isImportDeclaration(st) && (!st.importClause?.isTypeOnly || !st.importClause.namedBindings || !ts.isNamedImports(st.importClause.namedBindings))))
      unsupported("Package imports must be named type imports");
    sourceFiles.set(name, sf); packages.set(pkg.Module, name); packageByFile.set(name, pkg.Module);
  }
  function resolveImport(spec, importer) {
    if (!request.NamedImports) return undefined;
    const pkg = packages.get(spec);
    if (pkg) return { resolvedFileName: pkg, extension: ts.Extension.Dts, isExternalLibraryImport: true };
    if (!importer.startsWith("/input/") || !/^(?:\.\.?\/).*\.js$/.test(spec)) return undefined;
    const file = path.posix.resolve(path.posix.dirname(importer), spec.slice(0, -3) + ".ts");
    if (!file.startsWith("/input/") || !unitByFile.has(file)) return undefined;
    let canonical = path.posix.relative(path.posix.dirname(importer), file).slice(0, -3) + ".js";
    if (!canonical.startsWith("../")) canonical = "./" + canonical;
    return spec === canonical ? { resolvedFileName: file, extension: ts.Extension.Ts } : undefined;
  }
  const libDir = path.dirname(process.argv[1]);
  const libraries = new Map(fs.readdirSync(libDir).filter(name => /^lib\.[a-zA-Z0-9.]+\.d\.ts$/.test(name)).map(name => ["/compiler/" + name, path.join(libDir, name)]));
  const textFor = name => sourceFiles.get(name)?.text ?? (libraries.has(name) ? fs.readFileSync(libraries.get(name), "utf8") : undefined);
  /** @type {import("typescript").CompilerHost} */
  const host = {
    getSourceFile(name, version) {
      if (sourceFiles.has(name)) return sourceFiles.get(name);
      const text = textFor(name); return text === undefined ? undefined : ts.createSourceFile(name, text, version, true);
    },
    getDefaultLibFileName: () => "/compiler/lib.es2020.d.ts",
    getDefaultLibLocation: () => "/compiler",
    writeFile() { throw new Failure("environment", "The semantic helper must not emit files"); },
    getCurrentDirectory: () => "/input",
    getCanonicalFileName: name => name,
    useCaseSensitiveFileNames: () => true,
    getNewLine: () => "\n",
    fileExists: name => sourceFiles.has(name) || libraries.has(name),
    readFile: textFor,
    directoryExists: name => name === "/input" || name === "/compiler",
    getDirectories: () => [],
    resolveModuleNames: (names, importer) => names.map(spec => resolveImport(spec, importer)),
  };
  const program = ts.createProgram([...sourceFiles.keys()], options, host);
  const checker = program.getTypeChecker();
  const nativeByNode = new Map();
  const nodeSites = new Map();
  const keywordKinds = new Set([ts.SyntaxKind.AnyKeyword, ts.SyntaxKind.UnknownKeyword, ts.SyntaxKind.NeverKeyword, ts.SyntaxKind.VoidKeyword, ts.SyntaxKind.UndefinedKeyword, ts.SyntaxKind.BooleanKeyword, ts.SyntaxKind.NumberKeyword, ts.SyntaxKind.BigIntKeyword, ts.SyntaxKind.StringKeyword, ts.SyntaxKind.SymbolKeyword, ts.SyntaxKind.ObjectKeyword]);
  function compatible(kind, raw) {
    const node = unwrap(raw);
    switch (kind) {
      case "import": return ts.isImportDeclaration(node);
      case "import-binding": return ts.isImportSpecifier(node);
      case "type-alias": return ts.isTypeAliasDeclaration(node);
      case "interface": return ts.isInterfaceDeclaration(node);
      case "type-parameter": return ts.isTypeParameterDeclaration(node);
      case "type-reference": return ts.isTypeReferenceNode(node) || ts.isIdentifier(node) || ts.isExpressionWithTypeArguments(node);
      case "type-application": return (ts.isTypeReferenceNode(node) || ts.isExpressionWithTypeArguments(node)) && !!node.typeArguments?.length;
      case "keyword-type": return keywordKinds.has(node.kind);
      case "literal-type": return ts.isLiteralTypeNode(node);
      case "array-type": return ts.isArrayTypeNode(node);
      case "tuple-type": return ts.isTupleTypeNode(node);
      case "tuple-element": return ts.isNamedTupleMember(node) || ts.isTypeNode(node);
      case "object-type": return ts.isTypeLiteralNode(node);
      case "union-type": return ts.isUnionTypeNode(node);
      case "intersection-type": return ts.isIntersectionTypeNode(node);
      case "property": return ts.isPropertySignature(node);
      case "type-operator": return ts.isTypeOperatorNode(node);
      default: return false;
    }
  }
  for (const unit of request.Units) {
    const sf = sourceFiles.get("/input/" + unit.Path);
    const all = [];
    function visit(node) { all.push(node); ts.forEachChild(node, child => { visit(child); }); }
    visit(sf);
    for (const node of unit.Nodes) {
      let native;
      if (node.Kind === "source-file") native = sf;
      else {
        const candidates = all.filter(n => n.getStart(sf) === node.Start && compatible(node.Kind, n));
        candidates.sort((a, b) => Number(b.end === node.End) - Number(a.end === node.End) || Number(ts.isIdentifier(a)) - Number(ts.isIdentifier(b)) || a.end - b.end);
        native = candidates[0];
      }
      if (!native) throw new Failure("mismatch", `No compiler syntax corresponds to ${node.Kind}`, site(unit, node));
      nativeByNode.set(unit.Tree + ":" + node.ID, unwrap(native));
      // Only declaration sites participate in symbol identity. Expressions and
      // tuple wrappers can share compiler nodes without replacing those sites.
      if (["type-alias", "interface", "type-parameter", "property", "import-binding"].includes(node.Kind)) nodeSites.set(unwrap(native), site(unit, node));
    }
  }
  const diagnostics = [...program.getOptionsDiagnostics(), ...program.getGlobalDiagnostics(), ...program.getSyntacticDiagnostics(), ...program.getSemanticDiagnostics()];
  if (diagnostics.length) {
    const d = diagnostics[0]; const unit = d.file && unitByFile.get(d.file.fileName);
    const at = unit && unit.Nodes.filter(n => n.Start <= d.start && n.End >= d.start).sort((a, b) => a.End - a.Start - (b.End - b.Start))[0];
    const status = unit || d.file && packageByFile.has(d.file.fileName) ? "invalid" : "environment";
    return { Version: 1, Status: status, Diagnostics: [diagnostic(status, `TS${d.code}: ${ts.flattenDiagnosticMessageText(d.messageText, "\n")}`, at && site(unit, at))] };
  }

  const graph = { Entities: [], Terms: [], Bindings: [], Expressions: [], Dependencies: [], Facts: [] };
  const entityBySymbol = new Map(); const symbolByEntity = new Map(); const referenceByEntity = new Map();
  const typeIDs = new Map(); const typesByID = new Map(); const declaredOwners = new Map(); const fileOwners = new Map();
  let depth = 0;
  const count = () => {
    if (graph.Entities.length + graph.Terms.length + graph.Facts.length + graph.Bindings.length + graph.Expressions.length > 32768) throw new Failure("incomplete", "TypeScript semantic record limit exceeded");
  };
  function term(value) { count(); graph.Terms.push(value); return graph.Terms.length; }
  function expression(kind, data, operands = []) { return { Kind: 3, Operator: key(kind), Encoding: key("json-1"), Data: Buffer.from(JSON.stringify(data)).toString("base64"), Operands: operands.map((id, i) => ({ Locator: [i], Term: id })) }; }
  function fact(subject, relation, value) { count(); graph.Facts.push({ Subject: subject, Relation: key(relation), Term: value }); }
  function reference(id) { if (!referenceByEntity.has(id)) referenceByEntity.set(id, term({ Kind: 1, Target: id })); return referenceByEntity.get(id); }
  graph.Entities.push({ Kind: key("program"), Name: "captured", Declarations: [], Parameters: [] });
  for (const unit of request.Units) {
    const sf = sourceFiles.get("/input/" + unit.Path);
    graph.Entities.push({ Kind: key(ts.isExternalModule(sf) ? "module" : "script"), Name: unit.Path, Owner: 1, Declarations: [site(unit, unit.Nodes.find(n => n.Kind === "source-file"))], Parameters: [] });
    fileOwners.set(sf, graph.Entities.length);
  }
  for (const [file, specifier] of packageByFile) {
    graph.Entities.push({ Kind: key("package"), Name: specifier, Owner: 1, Declarations: [], Parameters: [] });
    fileOwners.set(sourceFiles.get(file), graph.Entities.length);
  }
  function packageSymbol(symbol) {
    return symbol?.getDeclarations()?.some(d => packageByFile.has(d.getSourceFile().fileName));
  }
  function namedPackageType(symbol) {
    return packageSymbol(symbol) && symbol.getDeclarations().some(d => ts.isInterfaceDeclaration(d) || ts.isTypeAliasDeclaration(d));
  }
  function resolvedSymbol(symbol) { return symbol && symbol.flags & ts.SymbolFlags.Alias ? checker.getAliasedSymbol(symbol) : symbol; }
  function captured(declaration) { return nodeSites.get(declaration); }
  function entity(symbol) {
    if (!symbol) throw new Failure("mismatch", "Missing compiler symbol");
    if (entityBySymbol.has(symbol)) return entityBySymbol.get(symbol);
    count();
    const declarations = symbol.getDeclarations() || [];
    const declared = declarations.filter(d => captured(d));
    if (declared.length > 1 && declared.some(d => ts.isInterfaceDeclaration(d))) unsupported("Interface declaration merging is outside captured-types-1", captured(declared[0]));
    const first = declared[0] || declarations[0];
    let kind = "builtin";
    if (declared.length) kind = ts.isTypeAliasDeclaration(first) ? "alias" : ts.isInterfaceDeclaration(first) ? "interface" : ts.isTypeParameterDeclaration(first) ? "parameter" : "property";
    else if (first && ts.isTypeParameterDeclaration(first)) kind = "parameter";
    else if (packageSymbol(symbol)) kind = "package-type";
    const id = graph.Entities.length + 1; entityBySymbol.set(symbol, id); symbolByEntity.set(id, symbol);
    graph.Entities.push({ Kind: key(kind), Name: symbol.getName(), Owner: 0, Declarations: declared.map(captured), Parameters: [] });
    if (first) {
      let parent = first.parent;
      while (parent && !ts.isSourceFile(parent) && !ts.isTypeAliasDeclaration(parent) && !ts.isInterfaceDeclaration(parent)) parent = parent.parent;
      if (parent && ts.isSourceFile(parent)) graph.Entities[id - 1].Owner = fileOwners.get(parent) || 0;
      else if (parent?.name) graph.Entities[id - 1].Owner = entity(checker.getSymbolAtLocation(parent.name));
    }
    if (first && (ts.isTypeAliasDeclaration(first) || ts.isInterfaceDeclaration(first))) {
      for (const parameter of first.typeParameters || []) {
        const pid = entity(checker.getSymbolAtLocation(parameter.name)); graph.Entities[pid - 1].Owner = id;
        graph.Entities[id - 1].Parameters.push({ Entity: pid, Sort: key("type") });
      }
    }
    return id;
  }
  function symbolFor(raw) {
    const node = unwrap(raw);
    return checker.getSymbolAtLocation(ts.isTypeReferenceNode(node) ? node.typeName : ts.isExpressionWithTypeArguments(node) ? node.expression : node.name || node);
  }
  // Allocate captured declarations before following compiler types, preserving
  // real binder identities even when aliases disappear from native types.
  for (const unit of request.Units) for (const node of unit.Nodes) if (["type-alias", "interface", "type-parameter", "property"].includes(node.Kind)) {
    const native = nativeByNode.get(unit.Tree + ":" + node.ID); const symbol = symbolFor(native); entity(symbol);
    if (node.Kind === "type-alias" || node.Kind === "interface") {
      const type = checker.getDeclaredTypeOfSymbol(symbol);
      if (!declaredOwners.has(type)) declaredOwners.set(type, symbol);
    }
  }
  function ownSymbol(symbol) { return symbol?.getDeclarations()?.some(d => captured(d)); }
  // External mapped/utility types can synthesize modifiers whose meaning is not
  // captured by declaration-owned property metadata. Admit only compiler-known
  // array constructors from the standard library and explicitly supplied package types.
  for (const unit of request.Units) for (const node of unit.Nodes) if (node.Kind === "type-reference") {
    const symbol = resolvedSymbol(symbolFor(nativeByNode.get(unit.Tree + ":" + node.ID)));
    if (!ownSymbol(symbol) && !packageSymbol(symbol) && !checker.isArrayType(checker.getDeclaredTypeOfSymbol(symbol))) unsupported("External type references other than array constructors are outside captured-types-1", site(unit, node));
  }
  function properties(type) {
    if (checker.getSignaturesOfType(type, ts.SignatureKind.Call).length || checker.getSignaturesOfType(type, ts.SignatureKind.Construct).length || checker.getIndexInfosOfType(type).length) unsupported("Callable and indexed object types are outside captured-types-1");
    const metadata = [], operands = [];
    for (const property of checker.getPropertiesOfType(type)) {
      const declaration = property.valueDeclaration || property.declarations?.[0];
      if (!declaration || !captured(declaration)) unsupported("External object members are outside captured-types-1");
      metadata.push({ Name: property.getName(), Optional: !!(property.flags & ts.SymbolFlags.Optional), Readonly: ts.canHaveModifiers(declaration) && !!ts.getModifiers(declaration)?.some(m => m.kind === ts.SyntaxKind.ReadonlyKeyword), Entity: entity(property) });
      operands.push(nativeType(checker.getTypeOfSymbolAtLocation(property, declaration)));
    }
    return term(expression("object", metadata, operands));
  }
  function nativeType(type) {
    if (typeIDs.has(type)) return typeIDs.get(type);
    if (++depth > 128) throw new Failure("incomplete", "TypeScript semantic depth exceeded");
    const id = term({}); typeIDs.set(type, id); typesByID.set(id, type);
    let value;
    const f = type.flags;
    /** @type {[number,string][]} */
    const primitives = [[ts.TypeFlags.Any, "any"], [ts.TypeFlags.Unknown, "unknown"], [ts.TypeFlags.Never, "never"], [ts.TypeFlags.Void, "void"], [ts.TypeFlags.Undefined, "undefined"], [ts.TypeFlags.Null, "null"], [ts.TypeFlags.String, "string"], [ts.TypeFlags.Number, "number"], [ts.TypeFlags.Boolean, "boolean"], [ts.TypeFlags.BigInt, "bigint"], [ts.TypeFlags.ESSymbol, "symbol"], [ts.TypeFlags.NonPrimitive, "object"]];
    const primitive = primitives.find(([flag]) => f & flag);
    // A captured alias may name an instantiation of a package type. Use the
    // package constructor to retain its effective arguments. The syntax graph
    // independently retains the binding to the local alias at each named use.
    const owner = [type.aliasSymbol, type.symbol].find(namedPackageType)
      || [type.aliasSymbol, type.symbol].find(ownSymbol) || declaredOwners.get(type);
    let structure;
    if (primitive) value = expression("primitive", primitive[1]);
    else if (f & (ts.TypeFlags.StringLiteral | ts.TypeFlags.NumberLiteral | ts.TypeFlags.BooleanLiteral | ts.TypeFlags.BigIntLiteral)) value = expression("literal", checker.typeToString(type));
    else if (f & ts.TypeFlags.TypeParameter) value = { Kind: 1, Target: entity(type.symbol) };
    else {
      // References break recursive type graphs; definition structure is a fact,
      // not a containment edge back into the named type.
      if (owner) {
        const original = checker.getDeclaredTypeOfSymbol(owner);
        if (original !== type) {
          const args = type.aliasSymbol === owner ? type.aliasTypeArguments : type.objectFlags & ts.ObjectFlags.Reference ? checker.getTypeArguments(type) : undefined;
          if (!args || args.length !== graph.Entities[entity(owner) - 1].Parameters.length) unsupported("Cannot expose effective arguments for this native type");
          value = { Kind: 2, Callee: reference(entity(owner)), Arguments: args.map(nativeType) };
        } else value = { Kind: 1, Target: entity(owner) };
      }
      // A selected package type retains its compiler identity, arguments and
      // checked uses. Its private structural graph is outside this profile.
      if (packageSymbol(owner) || packageSymbol(type.symbol)) {
        if (!value) unsupported("Anonymous package structure has no captured type representation");
      } else {
        if (checker.isTupleType(type)) structure = term(expression("tuple", { Flags: type.target.elementFlags, Readonly: type.target.readonly }, checker.getTypeArguments(type).map(nativeType)));
        else if (checker.isArrayType(type)) structure = term(expression("array", { Readonly: type.symbol?.getName() === "ReadonlyArray" }, checker.getTypeArguments(type).map(nativeType)));
        else if (f & ts.TypeFlags.Union) structure = term(expression("union", null, type.types.map(nativeType)));
        else if (f & ts.TypeFlags.Intersection) structure = term(expression("intersection", null, type.types.map(nativeType)));
        else if (f & ts.TypeFlags.Object) { structure = properties(type); fact({ Term: id }, "properties", structure); }
        else unsupported(`Native type is outside captured-types-1: ${checker.typeToString(type)}`);
        if (value) fact({ Term: id }, "structure", structure);
        else value = graph.Terms[structure - 1];
      }
    }
    graph.Terms[id - 1] = value;
    fact({ Term: id }, "native-type", term(expression("boolean", true)));
    fact({ Term: id }, "display", term({ Kind: 3, Operator: key("display"), Encoding: key("utf8"), Data: Buffer.from(checker.typeToString(type)).toString("base64") }));
    depth--; return id;
  }
  function expressionType(unit, node) {
    const native = nativeByNode.get(unit.Tree + ":" + node.ID);
    const type = nativeType(checker.getTypeAtLocation(native));
    let symbolic = type;
    if (node.Kind === "type-reference") symbolic = reference(entity(resolvedSymbol(symbolFor(native))));
    else if (node.Kind === "type-application") {
      const target = entity(resolvedSymbol(symbolFor(native)));
      const args = children(node, "argument").map(id => nativeType(checker.getTypeAtLocation(nativeByNode.get(unit.Tree + ":" + id))));
      if (args.length !== graph.Entities[target - 1].Parameters.length) unsupported("Application does not have complete effective arguments", site(unit, node));
      symbolic = term({ Kind: 2, Callee: reference(target), Arguments: args });
    }
    // One shared reference may occur at several sites with one declared type;
    // per-occurrence type associations belong to the occurrence, not that ref.
    graph.Expressions.push({ At: { Site: site(unit, node), Role: key("type-expression"), Ordinal: 0 }, Term: symbolic });
    graph.Expressions.push({ At: { Site: site(unit, node), Role: key("native-type"), Ordinal: 0 }, Term: type });
    return symbolic;
  }
  const expressions = new Set(["type-reference", "type-application", "keyword-type", "literal-type", "array-type", "tuple-type", "object-type", "union-type", "intersection-type", "type-operator"]);
  for (const unit of request.Units) for (const node of unit.Nodes) {
    const native = nativeByNode.get(unit.Tree + ":" + node.ID);
    if (node.Kind === "import-binding") {
      const target = entity(resolvedSymbol(checker.getSymbolAtLocation(native.name)));
      graph.Bindings.push({ At: { Site: site(unit, node), Role: key("import"), Ordinal: 0 }, Target: target });
    }
    if (node.Kind === "type-reference") {
      const symbol = symbolFor(native), actual = resolvedSymbol(symbol), target = entity(actual);
      const sites = symbol?.flags & ts.SymbolFlags.Alias ? symbol.getDeclarations().map(captured).filter(Boolean) : graph.Entities[target - 1].Declarations;
      for (const ref of node.References || []) if (ref.Role !== "target" || !sites.some(d => d.Tree === ref.Target.Tree && d.Node === ref.Target.Node)) throw new Failure("mismatch", "Supplied reference disagrees with the compiler", site(unit, node));
      graph.Bindings.push({ At: { Site: site(unit, node), Role: key("name"), Ordinal: 0 }, Target: target });
    }
    if (expressions.has(node.Kind)) expressionType(unit, node);
    if (["type-alias", "interface", "type-parameter"].includes(node.Kind)) {
      const symbol = symbolFor(native), id = entity(symbol), type = nativeType(checker.getDeclaredTypeOfSymbol(symbol));
      fact({ Entity: id }, "declared-type", type);
      if (node.Kind === "type-parameter" && native.constraint) fact({ Entity: id }, "constraint", nativeType(checker.getTypeFromTypeNode(native.constraint)));
      graph.Entities[id - 1].Body = type;
    }
  }
  // Record the actual public declarations of each captured module.
  for (const unit of request.Units) {
    const sf = sourceFiles.get("/input/" + unit.Path);
    const moduleSymbol = checker.getSymbolAtLocation(sf);
    if (!moduleSymbol) continue;
    for (const exported of checker.getExportsOfModule(moduleSymbol)) {
      const symbol = resolvedSymbol(exported);
      if (symbol.getDeclarations()?.some(d => d.getSourceFile() === sf && captured(d))) {
        fact({ Entity: fileOwners.get(sf) }, "named-export", reference(entity(symbol)));
      }
    }
  }
  // Property entities share the compiler's declared non-missing value type.
  // Presence and readonly flags remain separate object metadata. These are not
  // flow-sensitive read types or writable acceptance judgments.
  for (const [id, symbol] of symbolByEntity) if (graph.Entities[id - 1].Kind.Name === "property") {
    const d = symbol.valueDeclaration || symbol.declarations[0]; graph.Entities[id - 1].Body = nativeType(checker.getTypeOfSymbolAtLocation(symbol, d));
  }
  count();
  const result = { Version: 1, Status: "ok", Facts: graph, Diagnostics: [], Assignable: undefined };
  if (request.Query) {
    const source = typesByID.get(request.Query.Source), target = typesByID.get(request.Query.Target);
    if (!source || !target) throw new Failure("mismatch", "Query refers to foreign native types");
    result.Assignable = checker.isTypeAssignableTo(source, target);
  }
  return result;
}
try {
  const request = JSON.parse(fs.readFileSync(0, "utf8"));
  process.stdout.write(JSON.stringify(main(request)));
} catch (error) {
  const failure = error instanceof Failure ? error : new Failure("environment", String(error));
  process.stdout.write(JSON.stringify({ Version: 1, Status: failure.status, Diagnostics: [diagnostic(failure.status, failure.message, failure.at)] }));
}
