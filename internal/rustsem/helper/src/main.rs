#![feature(rustc_private)]

extern crate rustc_driver;
extern crate rustc_hir;
extern crate rustc_interface;
extern crate rustc_middle;
extern crate rustc_session;
extern crate rustc_span;

use base64::{Engine, engine::general_purpose::STANDARD};
use rustc_driver::{Callbacks, Compilation};
use rustc_hir::{
    self as hir,
    def::{DefKind, Res},
    def_id::DefId,
    intravisit::{self, Visitor},
};
use rustc_interface::interface;
use rustc_middle::ty::{self, Ty, TyCtxt, TypeVisitableExt};
use rustc_span::{FileName, Span, source_map::FileLoader};
use serde::Deserialize;
use serde_json::{Value, json};
use std::{
    collections::HashMap,
    io::{self, Read},
    path::Path,
    sync::Arc,
};

#[derive(Deserialize)]
#[serde(rename_all = "PascalCase", deny_unknown_fields)]
struct Input {
    version: u32,
    units: Vec<Unit>,
}
#[derive(Deserialize)]
#[serde(rename_all = "PascalCase", deny_unknown_fields)]
struct Unit {
    tree: String,
    text: String,
    nodes: Vec<Node>,
}
#[derive(Deserialize)]
#[serde(rename_all = "PascalCase", deny_unknown_fields)]
struct Node {
    #[serde(rename = "ID")]
    id: u32,
    kind: String,
    start: usize,
    end: usize,
    attributes: Option<Vec<Value>>,
    children: Option<Vec<Value>>,
    references: Option<Vec<Value>>,
}
impl Node {
    fn attribute(&self, key: &str) -> Option<&str> {
        self.attributes
            .as_deref()
            .unwrap_or(&[])
            .iter()
            .find(|a| a["Name"] == key)
            .and_then(|a| a["Text"].as_str())
    }
}
fn key(name: &str) -> Value {
    json!({"Owner":"rust","Name":name})
}
fn diagnostic(unit: Option<&Unit>, node: Option<&Node>, name: &str, message: &str) -> Value {
    let mut d = json!({"Code":key(name),"Subject":key("captured-type-coverage"),"Message":message});
    if let (Some(u), Some(n)) = (unit, node) {
        d["At"] = json!({"Tree":u.tree,"Node":n.id});
    }
    d
}
fn output(status: &str, facts: Option<Value>, diagnostics: Vec<Value>) -> Value {
    json!({"Version":1,"Compiler":env!("BITTYPE_RUSTC_VERSION"),"Status":status,"Facts":facts,"Diagnostics":diagnostics})
}
fn unsupported(unit: &Unit) -> Option<&Node> {
    unit.nodes.iter().find(|n| match n.kind.as_str() {
        "file" | "module" | "type-alias" | "struct" | "tuple-struct" | "enum" | "unit-variant"
        | "tuple-variant" | "record-variant" | "field" | "tuple-field" | "trait"
        | "type-parameter" | "where-predicate" | "relaxed-bound" | "path-type" | "path-segment"
        | "primitive-type" | "raw-pointer-type" | "slice-type" | "array-type" | "tuple-type"
        | "never-type" => false,
        "trait-impl" => n
            .children
            .as_deref()
            .unwrap_or(&[])
            .iter()
            .any(|c| c["Role"] == "member"),
        "reference-type" => !n
            .children
            .as_deref()
            .unwrap_or(&[])
            .iter()
            .any(|c| c["Role"] == "lifetime"),
        "lifetime-reference" => n.attribute("name") != Some("static"),
        _ => true,
    })
}

// Source lookup is deliberately closed. Only the captured Input::Str and the
// explicitly supplied sysroot's compiler metadata participate in this profile.
struct ClosedFiles;
impl FileLoader for ClosedFiles {
    fn file_exists(&self, _: &Path) -> bool {
        false
    }
    fn read_file(&self, _: &Path) -> io::Result<String> {
        Err(io::ErrorKind::PermissionDenied.into())
    }
    fn read_binary_file(&self, _: &Path) -> io::Result<Arc<[u8]>> {
        Err(io::ErrorKind::PermissionDenied.into())
    }
}
struct Driver {
    unit: Unit,
    result: Option<Value>,
}
impl Callbacks for Driver {
    fn config(&mut self, config: &mut interface::Config) {
        config.input = rustc_session::config::Input::Str {
            name: FileName::Custom("captured.rs".into()),
            input: self.unit.text.clone(),
        };
        config.file_loader = Some(Box::new(ClosedFiles));
    }
    fn after_analysis<'tcx>(&mut self, _: &interface::Compiler, tcx: TyCtxt<'tcx>) -> Compilation {
        if tcx.dcx().has_errors().is_some() {
            self.result = Some(output("compiler-errors", None, vec![]));
            return Compilation::Stop;
        }
        let result = Exporter::new(tcx, &self.unit).export();
        self.result = Some(match result {
            Ok(facts) => output("ok", Some(facts), vec![]),
            Err((status, message)) => output(
                status,
                None,
                vec![diagnostic(Some(&self.unit), None, status, &message)],
            ),
        });
        Compilation::Stop
    }
}

#[derive(Default)]
struct Catalogue {
    definitions: Vec<(DefId, Span, &'static str)>,
    paths: Vec<(Span, Res)>,
}
impl<'tcx> Visitor<'tcx> for Catalogue {
    fn visit_item(&mut self, item: &'tcx hir::Item<'tcx>) {
        // rustc injects a standard-library prelude import; it is compiler context,
        // not an authored declaration in the capture.
        if item.span.is_dummy() || item.span.from_expansion() {
            return;
        }
        let kind = match item.kind {
            hir::ItemKind::Mod(..) => "module",
            hir::ItemKind::TyAlias(..) => "type-alias",
            hir::ItemKind::Struct(..) => "struct",
            hir::ItemKind::Enum(..) => "enum",
            hir::ItemKind::Trait(..) => "trait",
            hir::ItemKind::Impl(..) => "trait-impl",
            _ => "unsupported",
        };
        self.definitions
            .push((item.owner_id.def_id.to_def_id(), item.span, kind));
        intravisit::walk_item(self, item);
    }
    fn visit_generic_param(&mut self, p: &'tcx hir::GenericParam<'tcx>) {
        self.definitions
            .push((p.def_id.to_def_id(), p.span, "type-parameter"));
        intravisit::walk_generic_param(self, p);
    }
    fn visit_field_def(&mut self, f: &'tcx hir::FieldDef<'tcx>) {
        self.definitions
            .push((f.def_id.to_def_id(), f.span, "field"));
        intravisit::walk_field_def(self, f);
    }
    fn visit_variant(&mut self, v: &'tcx hir::Variant<'tcx>) {
        self.definitions
            .push((v.def_id.to_def_id(), v.span, "variant"));
        intravisit::walk_variant(self, v);
    }
    fn visit_path(&mut self, path: &hir::Path<'tcx>, _: hir::HirId) {
        for (i, s) in path.segments.iter().enumerate() {
            let res = if matches!(s.res, Res::Err) && i + 1 == path.segments.len() {
                path.res
            } else {
                s.res
            };
            self.paths.push((s.ident.span, res));
        }
        intravisit::walk_path(self, path);
    }
}
type ExportResult<T> = Result<T, (&'static str, String)>;
fn unsupported_error(message: impl Into<String>) -> (&'static str, String) {
    ("unsupported", message.into())
}
struct Exporter<'tcx, 'a> {
    tcx: TyCtxt<'tcx>,
    unit: &'a Unit,
    entities: Vec<Value>,
    terms: Vec<Value>,
    bindings: Vec<Value>,
    expressions: Vec<Value>,
    facts: Vec<Value>,
    defs: HashMap<DefId, u32>,
    params: HashMap<(DefId, u32), u32>,
    types: HashMap<(Ty<'tcx>, Option<DefId>), u32>,
    sites: HashMap<DefId, u32>,
    depth: usize,
}
impl<'tcx, 'a> Exporter<'tcx, 'a> {
    fn new(tcx: TyCtxt<'tcx>, unit: &'a Unit) -> Self {
        Self {
            tcx,
            unit,
            entities: vec![],
            terms: vec![],
            bindings: vec![],
            expressions: vec![],
            facts: vec![],
            defs: HashMap::new(),
            params: HashMap::new(),
            types: HashMap::new(),
            sites: HashMap::new(),
            depth: 0,
        }
    }
    fn site(&self, id: u32) -> Value {
        json!({"Tree":self.unit.tree,"Node":id})
    }
    fn path(&self, did: DefId) -> String {
        format!(
            "{}{}",
            self.tcx.crate_name(did.krate),
            self.tcx.def_path(did).to_string_no_crate_verbose()
        )
    }
    fn find(&self, span: Span, kind: &str) -> Option<u32> {
        let lo = self.tcx.sess.source_map().lookup_byte_offset(span.lo());
        let start = lo.pos.0 as usize;
        let end = start + (span.hi() - span.lo()).0 as usize;
        self.unit
            .nodes
            .iter()
            .filter(|n| {
                let matches = match kind {
                    "struct" => n.kind == "struct" || n.kind == "tuple-struct",
                    "field" => n.kind == "field" || n.kind == "tuple-field",
                    "variant" => n.kind.ends_with("-variant"),
                    _ => n.kind == kind,
                };
                matches
                    && n.start <= start
                    && n.end >= end
                    && (kind != "path-segment" || n.start == start)
            })
            .min_by_key(|n| n.end - n.start)
            .map(|n| n.id)
    }
    fn term(&mut self, t: Value) -> u32 {
        self.terms.push(t);
        self.terms.len() as u32
    }
    fn reference(&mut self, target: u32) -> u32 {
        self.term(json!({"Kind":1,"Target":target}))
    }
    fn expression(&mut self, operator: &str, data: Value, args: Vec<u32>) -> u32 {
        let bytes = serde_json::to_vec(&data).unwrap();
        let operands: Vec<_> = args
            .into_iter()
            .enumerate()
            .map(|(i, t)| json!({"Locator":[i],"Term":t}))
            .collect();
        self.term(json!({"Kind":3,"Operator":key(operator),"Encoding":key("json-1"),"Data":STANDARD.encode(bytes),"Operands":operands}))
    }
    fn fact(&mut self, subject: Value, name: &str, term: u32) {
        self.facts
            .push(json!({"Subject":subject,"Relation":key(name),"Term":term}));
    }
    fn entity_fact(&mut self, id: u32, name: &str, term: u32) {
        self.fact(json!({"Entity":id}), name, term);
    }
    fn type_fact(&mut self, id: u32, name: &str, term: u32) {
        self.fact(json!({"Term":id}), name, term);
    }
    fn entity(&mut self, did: DefId) -> ExportResult<u32> {
        if let Some(&id) = self.defs.get(&did) {
            return Ok(id);
        }
        if matches!(self.tcx.def_kind(did), DefKind::TyParam) {
            let owner = self.tcx.parent(did);
            let p = self
                .tcx
                .generics_of(owner)
                .own_params
                .iter()
                .find(|p| p.def_id == did)
                .ok_or_else(|| unsupported_error("missing native binder"))?;
            return self.parameter(owner, p.index);
        }
        let kind = match self.tcx.def_kind(did) {
            DefKind::Mod => "module",
            DefKind::Struct => "struct",
            DefKind::Enum => "enum",
            DefKind::Variant => "variant",
            DefKind::Field => "field",
            DefKind::TyAlias => "alias",
            DefKind::Trait => "trait",
            DefKind::Impl { .. } => "impl",
            _ => {
                return Err(unsupported_error(format!(
                    "native declaration {:?}",
                    self.tcx.def_kind(did)
                )));
            }
        };
        let id = self.entities.len() as u32 + 1;
        self.defs.insert(did, id);
        let sites: Vec<_> = self
            .sites
            .get(&did)
            .map(|&n| vec![self.site(n)])
            .unwrap_or_default();
        self.entities
            .push(json!({"Kind":key(kind),"Name":self.path(did),"Declarations":sites}));
        if let Some(parent) = self.tcx.opt_parent(did) {
            let owner = self.entity(parent)?;
            self.entities[id as usize - 1]["Owner"] = json!(owner);
        }
        if matches!(
            self.tcx.def_kind(did),
            DefKind::Struct
                | DefKind::Enum
                | DefKind::TyAlias
                | DefKind::Trait
                | DefKind::Impl { .. }
        ) {
            let mut parameters = vec![];
            for p in &self.tcx.generics_of(did).own_params {
                let entity = self.parameter(did, p.index)?;
                parameters.push(json!({"Entity":entity,"Sort":key("type")}));
            }
            self.entities[id as usize - 1]["Parameters"] = json!(parameters);
        }
        Ok(id)
    }
    fn parameter(&mut self, owner: DefId, index: u32) -> ExportResult<u32> {
        let g = self.tcx.generics_of(owner);
        if (index as usize) < g.parent_count {
            return self.parameter(g.parent.unwrap(), index);
        }
        if let Some(&id) = self.params.get(&(owner, index)) {
            return Ok(id);
        }
        let p = g.param_at(index as usize, self.tcx);
        if !matches!(p.kind, ty::GenericParamDefKind::Type { .. }) {
            return Err(unsupported_error("lifetime/const generic parameters"));
        }
        // Trait Self reuses its trait's DefId in rustc. A binder is therefore
        // identified by (generic owner, slot), not by DefId alone.
        let entity_owner = self.entity(owner)?;
        if let Some(&id) = self.params.get(&(owner, index)) {
            return Ok(id);
        }
        let id = self.entities.len() as u32 + 1;
        self.params.insert((owner, index), id);
        let sites: Vec<_> = if p.def_id == owner {
            vec![]
        } else {
            self.sites
                .get(&p.def_id)
                .map(|&n| vec![self.site(n)])
                .unwrap_or_default()
        };
        self.entities.push(json!({"Kind":key("parameter"),"Name":format!("{}::{}",self.path(owner),p.name),"Owner":entity_owner,"Declarations":sites}));
        Ok(id)
    }
    fn arguments(
        &mut self,
        args: ty::GenericArgsRef<'tcx>,
        owner: Option<DefId>,
    ) -> ExportResult<Vec<u32>> {
        args.iter()
            .map(|a| match a.kind() {
                ty::GenericArgKind::Type(t) => self.native(t, owner),
                _ => Err(unsupported_error("lifetime/const native arguments")),
            })
            .collect()
    }
    fn native(&mut self, t: Ty<'tcx>, owner: Option<DefId>) -> ExportResult<u32> {
        if self.entities.len() + self.terms.len() + self.facts.len() >= 32768 {
            return Err(("incomplete", "native semantic record limit".into()));
        }
        if self.depth >= 128 {
            return Err(("incomplete", "native type expansion depth".into()));
        }
        self.depth += 1;
        let result = self.native_inner(t, owner);
        self.depth -= 1;
        result
    }
    fn native_inner(&mut self, t: Ty<'tcx>, owner: Option<DefId>) -> ExportResult<u32> {
        let context = if t.has_param() { owner } else { None };
        if let Some(&id) = self.types.get(&(t, context)) {
            return Ok(id);
        }
        if self.types.len() > 16000 {
            return Err(("incomplete", "native type graph limit".into()));
        }
        let term = match *t.kind() {
            ty::Param(p) => {
                let owner =
                    context.ok_or_else(|| unsupported_error("parameter without native owner"))?;
                let entity = self.parameter(owner, p.index)?;
                json!({"Kind":1,"Target":entity})
            }
            ty::Adt(adt, args) => {
                let entity = self.entity(adt.did())?;
                let callee = self.reference(entity);
                let args = self.arguments(args, context)?;
                if args.is_empty() {
                    json!({"Kind":1,"Target":entity})
                } else {
                    json!({"Kind":2,"Callee":callee,"Arguments":args})
                }
            }
            ty::Bool | ty::Char | ty::Int(_) | ty::Uint(_) | ty::Float(_) | ty::Str | ty::Never => {
                let term = self.expression("primitive", json!(t.to_string()), vec![]);
                self.types.insert((t, context), term);
                self.decorate(term, t);
                return Ok(term);
            }
            ty::Tuple(types) => {
                let args = types
                    .iter()
                    .map(|t| self.native(t, context))
                    .collect::<ExportResult<Vec<_>>>()?;
                let term = self.expression("tuple", Value::Null, args);
                self.types.insert((t, context), term);
                self.decorate(term, t);
                return Ok(term);
            }
            ty::Slice(element) => {
                let element = self.native(element, context)?;
                let term = self.expression("slice", Value::Null, vec![element]);
                self.types.insert((t, context), term);
                self.decorate(term, t);
                return Ok(term);
            }
            ty::Array(element, length) => {
                let length = length
                    .try_to_target_usize(self.tcx)
                    .ok_or_else(|| unsupported_error("unevaluated array length"))?;
                let element = self.native(element, context)?;
                let term = self.expression("array", json!(length), vec![element]);
                self.types.insert((t, context), term);
                self.decorate(term, t);
                return Ok(term);
            }
            ty::RawPtr(element, mutability) => {
                let element = self.native(element, context)?;
                let term = self.expression(
                    "pointer",
                    json!({"Mutable":mutability==hir::Mutability::Mut}),
                    vec![element],
                );
                self.types.insert((t, context), term);
                self.decorate(term, t);
                return Ok(term);
            }
            ty::Ref(region, element, mutability) => {
                if !matches!(region.kind(), ty::ReStatic) {
                    return Err(unsupported_error("non-static region"));
                }
                let element = self.native(element, context)?;
                let term = self.expression(
                    "reference",
                    json!({"Lifetime":"static","Mutable":mutability==hir::Mutability::Mut}),
                    vec![element],
                );
                self.types.insert((t, context), term);
                self.decorate(term, t);
                return Ok(term);
            }
            _ => return Err(unsupported_error(format!("native type {t}"))),
        };
        let id = self.term(term);
        self.types.insert((t, context), id);
        self.decorate(id, t);
        if let ty::Adt(adt, args) = *t.kind() {
            if adt.did().is_local() {
                let mut variants = vec![];
                let mut variant_terms = vec![];
                for variant in adt.variants() {
                    let mut fields = vec![];
                    let mut field_terms = vec![];
                    for f in &variant.fields {
                        let entity = self.entity(f.did)?;
                        let ft = self
                            .native(self.tcx.type_of(f.did).instantiate(self.tcx, args), context)?;
                        fields.push(json!({"Name":f.name.to_string(),"Public":self.tcx.visibility(f.did).is_public(),"Entity":entity}));
                        field_terms.push(ft);
                    }
                    let fields = self.expression("fields", json!(fields), field_terms);
                    if adt.is_struct() {
                        self.type_fact(id, "fields", fields);
                    } else {
                        variants.push(json!({"Name":variant.name.to_string(),"Entity":self.entity(variant.def_id)?}));
                        variant_terms.push(fields);
                    }
                }
                if adt.is_enum() {
                    let variants = self.expression("variants", json!(variants), variant_terms);
                    self.type_fact(id, "variants", variants);
                }
            }
        }
        Ok(id)
    }
    fn decorate(&mut self, id: u32, t: Ty<'tcx>) {
        let marker = self.expression("native-type", json!(true), vec![]);
        self.type_fact(id, "native-type", marker);
        let display=self.term(json!({"Kind":3,"Operator":key("display"),"Encoding":key("utf8"),"Data":STANDARD.encode(t.to_string().as_bytes())}));
        self.type_fact(id, "display", display);
    }
    fn predicates(&mut self, did: DefId) -> ExportResult<()> {
        let entity = self.entity(did)?;
        for (clause, _) in self.tcx.predicates_of(did).predicates {
            let bound = clause.kind();
            if !bound.bound_vars().is_empty() {
                return Err(unsupported_error("higher-ranked predicate"));
            }
            match bound.skip_binder() {
                ty::ClauseKind::Trait(predicate) => {
                    let target = self.entity(predicate.trait_ref.def_id)?;
                    let callee = self.reference(target);
                    let args = self.arguments(predicate.trait_ref.args, Some(did))?;
                    let app = self.term(json!({"Kind":2,"Callee":callee,"Arguments":args}));
                    let term = self.expression(
                        "trait-predicate",
                        json!({"Polarity":format!("{:?}",predicate.polarity)}),
                        vec![app],
                    );
                    self.entity_fact(entity, "predicate", term);
                }
                _ => return Err(unsupported_error(format!("predicate {clause:?}"))),
            }
        }
        Ok(())
    }
    fn export(mut self) -> ExportResult<Value> {
        let mut cat = Catalogue::default();
        self.tcx.hir_visit_all_item_likes_in_crate(&mut cat);
        for &(did, span, kind) in &cat.definitions {
            let site = self.find(span, kind).ok_or_else(|| {
                unsupported_error(format!(
                    "no captured site for {} ({kind})",
                    self.tcx.def_path_str(did)
                ))
            })?;
            self.sites.insert(did, site);
        }
        for &(did, _, _) in &cat.definitions {
            let entity = self.entity(did)?;
            let kind = self.tcx.def_kind(did);
            if matches!(
                kind,
                DefKind::Struct | DefKind::Enum | DefKind::TyAlias | DefKind::Field
            ) {
                let owner = if kind == DefKind::Field {
                    self.tcx.parent(did)
                } else {
                    did
                };
                let owner = if self.tcx.def_kind(owner) == DefKind::Variant {
                    self.tcx.parent(owner)
                } else {
                    owner
                };
                let t = self.native(self.tcx.type_of(did).instantiate_identity(), Some(owner))?;
                self.entity_fact(entity, "declared-type", t);
                self.entities[entity as usize - 1]["Body"] = json!(t);
                let at = json!({"Site":self.site(self.sites[&did]),"Role":key("declared-type")});
                self.expressions.push(json!({"At":at,"Term":t}));
            }
            if matches!(
                kind,
                DefKind::Struct
                    | DefKind::Enum
                    | DefKind::TyAlias
                    | DefKind::Trait
                    | DefKind::Impl { .. }
            ) {
                self.predicates(did)?;
                for p in &self.tcx.generics_of(did).own_params {
                    if matches!(
                        p.kind,
                        ty::GenericParamDefKind::Type {
                            has_default: true,
                            ..
                        }
                    ) {
                        let entity = self.parameter(did, p.index)?;
                        let term = self
                            .native(self.tcx.type_of(p.def_id).instantiate_identity(), Some(did))?;
                        self.entity_fact(entity, "default", term);
                    }
                }
                if let DefKind::Impl { .. } = kind {
                    if let Some(trait_ref) = self.tcx.impl_trait_ref(did) {
                        let tr = trait_ref.instantiate_identity();
                        let entity = self.entity(did)?;
                        let target = self.entity(tr.def_id)?;
                        let callee = self.reference(target);
                        let args = self.arguments(tr.args, Some(did))?;
                        let app = self.term(json!({"Kind":2,"Callee":callee,"Arguments":args}));
                        self.entity_fact(entity, "implements", app);
                    }
                }
            }
        }
        let mut targets = HashMap::new();
        for (span, res) in cat.paths {
            let Some(site) = self.find(span, "path-segment") else {
                continue;
            };
            let target = match res {
                Res::Def(_, did) => self.entity(did)?,
                Res::SelfTyParam { trait_ } => self.parameter(trait_, 0)?,
                Res::SelfTyAlias { alias_to, .. } => self.entity(alias_to)?,
                _ => return Err(unsupported_error(format!("native path resolution {res:?}"))),
            };
            if let Some(previous) = targets.insert(site, target) {
                if previous != target {
                    return Err(unsupported_error("conflicting native path sites"));
                } else {
                    continue;
                }
            }
            let at = json!({"Site":self.site(site),"Role":key("target")});
            self.bindings.push(json!({"At":at,"Target":target}));
        }
        for node in &self.unit.nodes {
            if node.kind == "path-segment" && !targets.contains_key(&node.id) {
                return Err(unsupported_error(format!(
                    "unmapped path segment {}",
                    node.id
                )));
            }
            for r in node.references.as_deref().unwrap_or(&[]) {
                let Some(&target) = targets.get(&node.id) else {
                    return Err(("mismatch", "reference has no native binding".into()));
                };
                if r["Role"] != "target"
                    || !self.entities[target as usize - 1]["Declarations"]
                        .as_array()
                        .unwrap()
                        .contains(&r["Target"])
                {
                    return Err((
                        "mismatch",
                        "supplied reference disagrees with native resolution".into(),
                    ));
                }
            }
        }
        if self.entities.len()
            + self.terms.len()
            + self.facts.len()
            + self.bindings.len()
            + self.expressions.len()
            > 32768
        {
            return Err(("incomplete", "native semantic record limit".into()));
        }
        Ok(
            json!({"Entities":self.entities,"Terms":self.terms,"Bindings":self.bindings,"Expressions":self.expressions,"Facts":self.facts}),
        )
    }
}
fn main() {
    let sysroot = std::env::args().nth(1).expect("explicit sysroot");
    assert!(Path::new(&sysroot).is_absolute());
    let mut bytes = vec![];
    io::stdin().take(2 << 20).read_to_end(&mut bytes).unwrap();
    let input: Input = serde_json::from_slice(&bytes).expect("private protocol input");
    assert_eq!(input.version, 1);
    if input.units.len() != 1 {
        println!(
            "{}",
            output(
                "unsupported",
                None,
                vec![diagnostic(
                    None,
                    None,
                    "units",
                    "one captured crate file is required"
                )]
            )
        );
        return;
    }
    let unit = input.units.into_iter().next().unwrap();
    if let Some(node) = unsupported(&unit) {
        println!(
            "{}",
            output(
                "unsupported",
                None,
                vec![diagnostic(
                    Some(&unit),
                    Some(node),
                    "syntax",
                    "syntax is outside the captured Rust profile"
                )]
            )
        );
        return;
    }
    let mut driver = Driver { unit, result: None };
    let args = [
        "bittype-rust-semantics",
        "/__bittype__/captured.rs",
        "--sysroot",
        &sysroot,
        "--crate-name",
        "captured",
        "--crate-type",
        "lib",
        "--edition",
        "2024",
        "--target",
        "x86_64-unknown-linux-gnu",
        "--emit",
        "metadata",
        "--error-format",
        "json",
        "-A",
        "warnings",
    ]
    .map(String::from)
    .to_vec();
    let result =
        rustc_driver::catch_fatal_errors(|| rustc_driver::run_compiler(&args, &mut driver));
    println!(
        "{}",
        driver.result.unwrap_or_else(|| {
            if result.is_err() {
                output("compiler-errors", None, vec![])
            } else {
                output(
                    "environment",
                    None,
                    vec![diagnostic(
                        None,
                        None,
                        "driver",
                        "compiler did not finish analysis",
                    )],
                )
            }
        })
    );
}
