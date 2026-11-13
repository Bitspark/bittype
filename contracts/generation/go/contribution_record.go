package generation

// The selected representation plan and its serialized form (bittype#196).
// Use entries and representation nodes (row 18), the identity of a node (row
// 19), applications (row 20), logical symbol identities (row 23) and the
// version axes (row 16).
//
// What is final and what is provisional (row 30). The meanings are final: what
// each identity means and who issues it, the kinds and needs of a link, the
// generic conventions and requirements, and what a node's identity is over.
// The layout and the encoding are provisional and belong to #195: the
// serialized form is JSON under the envelope version below, which is
// explicitly unstable. A later change may alter the layout, the encoding and
// the digest procedure of a node's identity under a new envelope version. It
// may not alter what an identity means.

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"slices"

	"github.com/Bitspark/bittype/trees/ast/go"
)

// EnvelopeVersion is the first version axis: how the shared graph is
// represented. It is explicitly unstable.
const EnvelopeVersion = "bittype.plan/0-unstable"

// ErrLinking: two inputs of a reader do not belong together: a plan and the
// realization manifest of another output, or a plan and a source snapshot of
// another pin. It is not ErrIncomplete.
var ErrLinking = errors.New("plan linking error")

// ErrVersion: a reader met a required value of a version axis that it does
// not know. Readers wrap it with the axis.
var ErrVersion = errors.New("unknown required version")

type PlanSite struct {
	Tree string `json:"tree"`
	Node uint32 `json:"node"`
}

type PlanFact struct {
	Kind    EvidenceKind `json:"kind"`
	Outcome string       `json:"outcome"`
	Content string       `json:"content,omitempty"`
}

// PlanAnswer is an answer the planner issued, as a link of the frozen record.
type PlanAnswer struct {
	Kind    EvidenceKind `json:"kind"`
	Outcome string       `json:"outcome"`
	Content string       `json:"content,omitempty"`
	// Required: the rule declared the fact as a requirement of the binder.
	Required bool `json:"required,omitempty"`
	// Reference is the contribution reference of the dependency the answer is
	// bound to. It is absent for a dependency bound to a binder.
	Reference *ContributionReference `json:"reference,omitempty"`
}

// PlanLink is one resolved request of a rule: a link that the kernel carries.
type PlanLink struct {
	Kind string `json:"kind"` // child | field
	Need string `json:"need"` // reference | evidence
	// Stated is what the rule stated in the request (O3).
	Stated string `json:"stated,omitempty"`
	// Position is the source's own label of the child among the subject's
	// children, passed through unread: an operand locator, or the index of a
	// binder for an actual. A reader with the source semantics interprets it.
	Position []uint32 `json:"position,omitempty"`
	// Role is the key of the role selector of a field position.
	Role    string       `json:"role,omitempty"`
	Node    string       `json:"node"`
	Answers []PlanAnswer `json:"answers,omitempty"`
}

type PlanBinder struct {
	Scope string `json:"scope"` // the logical symbol of the declaration that binds it
	Index int    `json:"index"`
}

// PlanArgument ties a node to one binder of the applied declaration.
type PlanArgument struct {
	Binder int    `json:"binder"`
	Node   string `json:"node"`
}

type PlanRequirement struct {
	Binder int          `json:"binder"`
	Kind   EvidenceKind `json:"kind"`
}

type PlanDischarge struct {
	Binder    int                    `json:"binder"`
	Kind      EvidenceKind           `json:"kind"`
	Outcome   string                 `json:"outcome"` // proved | deferred
	Reference *ContributionReference `json:"reference,omitempty"`
}

// PlanApplication links the binders of an applied declaration to the target
// arguments and to the adapter arguments of one use. The two lists are tied to
// binders and never matched by position (row 20).
type PlanApplication struct {
	Declaration string          `json:"declaration"` // the node of the open declaration
	Target      []PlanArgument  `json:"target,omitempty"`
	Adapters    []PlanArgument  `json:"adapters,omitempty"`
	Discharges  []PlanDischarge `json:"discharges,omitempty"`
}

// PlanNode is a representation node: the decisions and the dependencies of a
// selected representation. Its ID is its identity (row 19): a digest over
// everything below except diagnostic provenance, which a node does not hold.
type PlanNode struct {
	ID   string `json:"id"`
	Kind string `json:"kind"` // use | field | declaration | application | binder | owned
	// Source is the source identity of the form or declaration, or a role key.
	Source       string                 `json:"source,omitempty"`
	Contribution *ContributionReference `json:"contribution,omitempty"`
	Schema       *SchemaReference       `json:"schema,omitempty"`
	// Payload is the frozen value in hexadecimal and Receipt the kernel's record
	// that this value passed validation against Schema.
	Payload string     `json:"payload,omitempty"`
	Receipt string     `json:"receipt,omitempty"`
	Links   []PlanLink `json:"links,omitempty"`
	Facts   []PlanFact `json:"facts,omitempty"`
	// Symbol is the logical symbol identity of a planned declaration.
	Symbol       string            `json:"symbol,omitempty"`
	Binders      int               `json:"binders,omitempty"`
	Requirements []PlanRequirement `json:"requirements,omitempty"`
	// Transparent: the declaration contribution decided an alias in the target, so a
	// use of the declaration is represented as its body is (O35).
	Transparent bool             `json:"transparent,omitempty"`
	Binder      *PlanBinder      `json:"binder,omitempty"`
	Application *PlanApplication `json:"application,omitempty"`
	// Members are the spelled names of the members of the target type, in its order.
	Members []string `json:"members,omitempty"`
}

// PlanUse is a use entry: one selection at one source use, with its
// provenance. Several use entries refer to one node.
type PlanUse struct {
	ID    int      `json:"id"`
	Node  string   `json:"node"`
	Site  PlanSite `json:"site"`
	Scope string   `json:"scope,omitempty"`
}

// PlanSymbol is a planned declaration by its logical symbol identity, with the
// spelled intent. The name and the place come from the realization.
type PlanSymbol struct {
	ID         string     `json:"id"`
	Role       string     `json:"role"` // requested | dependency | owned
	Intent     string     `json:"intent"`
	Node       string     `json:"node"`
	Provenance []PlanSite `json:"provenance,omitempty"`
}

// PlanFacet is a facet value that crossed into the serialized plan through
// the transport descriptor of its kind (row 26a).
type PlanFacet struct {
	Kind      FacetKindReference `json:"kind"`
	MediaType string             `json:"mediaType"`
	Value     json.RawMessage    `json:"value"`
}

type PlanContribution struct {
	Reference ContributionReference `json:"reference"`
	Target    string                `json:"target"`
	Profile   string                `json:"profile"`
	Domain    string                `json:"domain"`
	Schema    SchemaReference       `json:"schema"`
	// Facts are the facts the contribution declared at registration.
	Facts  []PlanFact  `json:"facts,omitempty"`
	Facets []PlanFacet `json:"facets,omitempty"`
}

type Setting struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PlanTarget struct {
	Language      string    `json:"language"`
	Configuration []Setting `json:"configuration,omitempty"`
}

type PlanSource struct {
	Language string `json:"language"`
	// Semantics is the fifth version axis: which source interpretation the
	// checked references use. The source integration states it.
	Semantics string `json:"semantics"`
}

type planData struct {
	Envelope      string             `json:"envelope"`
	Profile       string             `json:"profile"`
	Products      []string           `json:"products"`
	Target        PlanTarget         `json:"target"`
	Source        PlanSource         `json:"source"`
	Contributions []PlanContribution `json:"contributions"`
	Nodes         []PlanNode         `json:"nodes"`
	Uses          []PlanUse          `json:"uses"`
	Symbols       []PlanSymbol       `json:"symbols"`
	Roots         []string           `json:"roots"`
}

// RepresentationPlan is the selected representation plan. It holds no pointer
// of the planning process: every identity in it is a value.
type RepresentationPlan struct {
	owner *identity
	data  planData
	index map[string]int
}

func (p RepresentationPlan) Valid() bool           { return p.owner != nil }
func (p RepresentationPlan) Profile() string       { return p.data.Profile }
func (p RepresentationPlan) Products() []string    { return slices.Clone(p.data.Products) }
func (p RepresentationPlan) Target() PlanTarget    { return p.data.Target }
func (p RepresentationPlan) Source() PlanSource    { return p.data.Source }
func (p RepresentationPlan) Nodes() []PlanNode     { return slices.Clone(p.data.Nodes) }
func (p RepresentationPlan) Uses() []PlanUse       { return slices.Clone(p.data.Uses) }
func (p RepresentationPlan) Symbols() []PlanSymbol { return slices.Clone(p.data.Symbols) }
func (p RepresentationPlan) Roots() []string       { return slices.Clone(p.data.Roots) }
func (p RepresentationPlan) Contributions() []PlanContribution {
	return slices.Clone(p.data.Contributions)
}
func (p RepresentationPlan) Node(id string) (PlanNode, bool) {
	i, ok := p.index[id]
	if !ok {
		return PlanNode{}, false
	}
	return p.data.Nodes[i], true
}
func (p RepresentationPlan) Symbol(id string) (PlanSymbol, bool) {
	for _, s := range p.data.Symbols {
		if s.ID == id {
			return s, true
		}
	}
	return PlanSymbol{}, false
}
func (p RepresentationPlan) Contribution(r ContributionReference) (PlanContribution, bool) {
	for _, c := range p.data.Contributions {
		if c.Reference == r {
			return c, true
		}
	}
	return PlanContribution{}, false
}

// Serialize writes the plan. Every payload in it was frozen when its rule
// returned: nothing is encoded now.
func (p RepresentationPlan) Serialize() ([]byte, error) {
	if !p.Valid() {
		return nil, ErrInvalid
	}
	return json.MarshalIndent(p.data, "", " ")
}

// Digest is the content digest of the serialized plan.
func (p RepresentationPlan) Digest() (string, error) {
	b, err := p.Serialize()
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

// ParsePlan reads a serialized plan in a process that need not have planned
// it. It refuses an envelope version it does not know, a node whose identity
// is not the digest of its content, a payload whose receipt was issued for
// another value, and a link to a node the plan does not hold.
func ParsePlan(serialized []byte) (RepresentationPlan, error) {
	var data planData
	if err := json.Unmarshal(serialized, &data); err != nil {
		return RepresentationPlan{}, fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	if data.Envelope != EnvelopeVersion {
		return RepresentationPlan{}, fmt.Errorf("%w: envelope protocol %q", ErrVersion, data.Envelope)
	}
	p := RepresentationPlan{owner: fresh(), data: data, index: map[string]int{}}
	for i, n := range data.Nodes {
		if n.ID != nodeIdentity(n) {
			return RepresentationPlan{}, fmt.Errorf("%w: node %s is not identified by its content", ErrInvalid, n.ID)
		}
		if n.Schema != nil {
			frozen, err := hex.DecodeString(n.Payload)
			if err != nil {
				return RepresentationPlan{}, fmt.Errorf("%w: payload of node %s", ErrInvalid, n.ID)
			}
			if err := CheckReceipt(n.Receipt, *n.Schema, frozen); err != nil {
				return RepresentationPlan{}, err
			}
		}
		if _, dup := p.index[n.ID]; dup {
			return RepresentationPlan{}, fmt.Errorf("%w: node %s twice", ErrInvalid, n.ID)
		}
		p.index[n.ID] = i
	}
	missing := func(id string) bool { _, ok := p.index[id]; return !ok }
	for _, n := range data.Nodes {
		for _, l := range n.Links {
			if l.Node != "" && missing(l.Node) {
				return RepresentationPlan{}, fmt.Errorf("%w: node %s links to a node the plan does not hold", ErrIncomplete, n.ID)
			}
		}
		if a := n.Application; a != nil {
			if missing(a.Declaration) {
				return RepresentationPlan{}, fmt.Errorf("%w: application %s without its declaration", ErrIncomplete, n.ID)
			}
			for _, arg := range append(slices.Clone(a.Target), a.Adapters...) {
				if missing(arg.Node) {
					return RepresentationPlan{}, fmt.Errorf("%w: application %s without an argument", ErrIncomplete, n.ID)
				}
			}
		}
	}
	for _, u := range data.Uses {
		if missing(u.Node) {
			return RepresentationPlan{}, fmt.Errorf("%w: use entry %d without its node", ErrIncomplete, u.ID)
		}
	}
	for _, s := range data.Symbols {
		if missing(s.Node) {
			return RepresentationPlan{}, fmt.Errorf("%w: symbol %s without its node", ErrIncomplete, s.ID)
		}
	}
	return p, nil
}

// receipt is the kernel's record that one frozen value passed validation
// against one schema: a digest over the schema reference and the value. Only
// the planner issues one, at the moment a rule returns (guarantee 6).
func receipt(schema SchemaReference, frozen []byte) string {
	h := sha256.New()
	fmt.Fprintf(h, "bittype.receipt/0\x00%s\x00%s\x00", schema.Identity, schema.Version)
	h.Write(frozen)
	return hex.EncodeToString(h.Sum(nil))
}

// CheckReceipt refuses a receipt that was issued for another frozen value or
// another schema. A backend calls it before it reads a payload.
func CheckReceipt(issued string, schema SchemaReference, frozen []byte) error {
	if issued == "" || issued != receipt(schema, frozen) {
		return fmt.Errorf("%w: the receipt was not issued for this frozen value and schema %s/%s", ErrPayload, schema.Identity, schema.Version)
	}
	return nil
}

// nodeIdentity is the identity of a node: a digest over the selected semantic
// representation (row 19). It covers the recorded decisions, the selected
// contribution with both parts of its reference, the payload schema, the
// links with the answers that were issued, the stated facts, the generic
// conventions and requirements and the identity of a nominal declaration. It
// covers no diagnostic provenance: a node holds none. Two uses with equal
// content share a node. Two different contributions never do, because the
// reference is part of the content.
//
// The digest is taken over a framing of the content that this function
// defines, and not over the serialized form, so that a change of the encoding
// leaves an identity as it is. The receipt is left out: it is derived from
// the schema and the payload, which are covered.
func nodeIdentity(n PlanNode) string {
	w := identityWriter{sha256.New()}
	w.text("bittype.node/0")
	w.text(n.Kind)
	w.text(n.Source)
	w.reference(n.Contribution)
	w.flag(n.Schema != nil)
	if n.Schema != nil {
		w.texts(n.Schema.Identity, n.Schema.Version)
	}
	w.text(n.Payload)
	w.count(len(n.Links))
	for _, l := range n.Links {
		w.texts(l.Kind, l.Need, l.Stated, l.Role, l.Node)
		w.count(len(l.Position))
		for _, position := range l.Position {
			w.count(int(position))
		}
		w.count(len(l.Answers))
		for _, a := range l.Answers {
			w.texts(a.Kind.Identity, a.Kind.Version, a.Outcome, a.Content)
			w.flag(a.Required)
			w.reference(a.Reference)
		}
	}
	w.count(len(n.Facts))
	for _, f := range n.Facts {
		w.texts(f.Kind.Identity, f.Kind.Version, f.Outcome, f.Content)
	}
	w.text(n.Symbol)
	w.count(n.Binders)
	w.count(len(n.Requirements))
	for _, r := range n.Requirements {
		w.count(r.Binder)
		w.texts(r.Kind.Identity, r.Kind.Version)
	}
	w.flag(n.Transparent)
	w.flag(n.Binder != nil)
	if n.Binder != nil {
		w.text(n.Binder.Scope)
		w.count(n.Binder.Index)
	}
	w.flag(n.Application != nil)
	if a := n.Application; a != nil {
		w.text(a.Declaration)
		for _, arguments := range [][]PlanArgument{a.Target, a.Adapters} {
			w.count(len(arguments))
			for _, argument := range arguments {
				w.count(argument.Binder)
				w.text(argument.Node)
			}
		}
		w.count(len(a.Discharges))
		for _, d := range a.Discharges {
			w.count(d.Binder)
			w.texts(d.Kind.Identity, d.Kind.Version, d.Outcome)
			w.reference(d.Reference)
		}
	}
	w.count(len(n.Members))
	w.texts(n.Members...)
	return hex.EncodeToString(w.h.Sum(nil)[:16])
}

// identityWriter frames content for a digest: every text is written with its
// length, so that no two contents have one framing.
type identityWriter struct{ h hash.Hash }

func (w identityWriter) count(n int) {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(n))
	w.h.Write(b[:])
}
func (w identityWriter) text(s string) {
	w.count(len(s))
	w.h.Write([]byte(s))
}
func (w identityWriter) texts(values ...string) {
	for _, s := range values {
		w.text(s)
	}
}
func (w identityWriter) flag(b bool) {
	if b {
		w.count(1)
	} else {
		w.count(0)
	}
}
func (w identityWriter) reference(r *ContributionReference) {
	w.flag(r != nil)
	if r != nil {
		w.texts(r.Identity.Publisher, r.Identity.Key, r.Meaning.Identity, r.Meaning.Revision, r.Representation)
	}
}

// planBuilder collects the plan while the planner runs.
type planBuilder struct {
	data  planData
	index map[string]int
}

func (b *planBuilder) node(n PlanNode) string {
	n.ID = nodeIdentity(n)
	if _, ok := b.index[n.ID]; !ok {
		b.index[n.ID] = len(b.data.Nodes)
		b.data.Nodes = append(b.data.Nodes, n)
	}
	return n.ID
}
func (b *planBuilder) use(node string, site ast.NodeRef, scope string) {
	b.data.Uses = append(b.data.Uses, PlanUse{ID: len(b.data.Uses) + 1, Node: node, Site: planSite(site), Scope: scope})
}
func planSite(site ast.NodeRef) PlanSite { return PlanSite{string(site.Tree), uint32(site.Node)} }
func planFacts(facts []Fact) []PlanFact {
	var out []PlanFact
	for _, f := range facts {
		out = append(out, PlanFact{f.Kind, f.Outcome.String(), f.Content})
	}
	return out
}

// ---------------------------------------------------------------------------
// What a reader knows (row 16, row 27): "the planner emits them and the
// plan's reader refuses".

// ReaderKnowledge is what one reader of a plan can interpret: the payload
// schemas it has codecs or decoders for, the contribution references it has
// an implementation for, and the source interpretations it reads. The
// envelope protocol is checked by ParsePlan. Generated-code and runtime
// compatibility is stated in a facet and checked by the reader of that facet.
type ReaderKnowledge struct {
	Schemas         []SchemaReference
	Contributions   []ContributionReference
	SourceSemantics []string
}

// CheckKnown refuses a plan that requires a value of a version axis that the
// reader does not know. Each axis is refused under its own name.
func (p RepresentationPlan) CheckKnown(k ReaderKnowledge) error {
	if !slices.Contains(k.SourceSemantics, p.data.Source.Semantics) {
		return fmt.Errorf("%w: source semantics %q", ErrVersion, p.data.Source.Semantics)
	}
	for _, c := range p.data.Contributions {
		if !slices.Contains(k.Schemas, c.Schema) {
			return fmt.Errorf("%w: payload schema %s version %q", ErrVersion, c.Schema.Identity, c.Schema.Version)
		}
		if slices.Contains(k.Contributions, c.Reference) {
			continue
		}
		part := "contribution " + c.Reference.String()
		for _, known := range k.Contributions {
			if known.Identity != c.Reference.Identity || known.Meaning.Identity != c.Reference.Meaning.Identity {
				continue
			}
			if known.Meaning.Revision != c.Reference.Meaning.Revision {
				part = fmt.Sprintf("meaning revision %q of %s", c.Reference.Meaning.Revision, c.Reference.Meaning.Identity)
			} else {
				part = fmt.Sprintf("representation revision %q of %s/%s", c.Reference.Representation, c.Reference.Identity.Publisher, c.Reference.Identity.Key)
			}
		}
		return fmt.Errorf("%w: representation semantics: %s", ErrVersion, part)
	}
	return nil
}

// ---------------------------------------------------------------------------
// The realization manifest (row 23, O16).

// RealizedSymbol is the name and the place that a language's Files stage gave
// one planned declaration.
type RealizedSymbol struct {
	Symbol string `json:"symbol"` // the logical symbol identity of the plan
	Name   string `json:"name"`
	Module string `json:"module"`
}

// RealizationManifest relates a plan, an output and its realization. It is
// serializable. It identifies the plan by its content digest, the generation
// run, the target configuration and the placement, and gives the realized
// symbol of each logical symbol. Who writes it is not specified (#195); on
// this branch the test composition does.
type RealizationManifest struct {
	Plan      string           `json:"plan"`
	Run       string           `json:"run"`
	Target    PlanTarget       `json:"target"`
	Placement string           `json:"placement"`
	Symbols   []RealizedSymbol `json:"symbols"`
}

// LinkRealization gives a reader the realized symbol of each logical symbol
// of a plan. A manifest of another plan or another target is a linking error.
// It is not an incomplete plan: nothing is missing, two inputs do not belong
// together.
func LinkRealization(p RepresentationPlan, m RealizationManifest) (map[string]RealizedSymbol, error) {
	digest, err := p.Digest()
	if err != nil {
		return nil, err
	}
	if m.Plan != digest {
		return nil, fmt.Errorf("%w: the realization manifest is of plan %s and the plan is %s", ErrLinking, m.Plan, digest)
	}
	if m.Target.Language != p.data.Target.Language || !slices.Equal(m.Target.Configuration, p.data.Target.Configuration) {
		return nil, fmt.Errorf("%w: the realization manifest is of another target or configuration", ErrLinking)
	}
	realized := map[string]RealizedSymbol{}
	for _, s := range m.Symbols {
		if _, dup := realized[s.Symbol]; dup || s.Name == "" || s.Module == "" {
			return nil, fmt.Errorf("%w: the realization of symbol %s", ErrInvalid, s.Symbol)
		}
		realized[s.Symbol] = s
	}
	for _, s := range p.data.Symbols {
		if _, ok := realized[s.ID]; !ok {
			return nil, fmt.Errorf("%w: symbol %s has no realization", ErrIncomplete, s.ID)
		}
	}
	return realized, nil
}
