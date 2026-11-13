package generation

// PROTOTYPE (bittype#196): the selected representation plan and its
// serialized form. Use entries and representation nodes (row 18), the
// identity of a node (row 19), applications (row 20), logical symbol
// identities (row 23) and the version axes (row 16). The layout and the
// encoding are provisional (row 30); the meanings are what the gate tests.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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

// PlanApplication links the binders of an applied declaration to the native
// arguments and to the adapter arguments of one use. The two lists are tied to
// binders and never matched by position (row 20).
type PlanApplication struct {
	Declaration string          `json:"declaration"` // the node of the open declaration
	Native      []PlanArgument  `json:"native,omitempty"`
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
	Binder       *PlanBinder       `json:"binder,omitempty"`
	Application  *PlanApplication  `json:"application,omitempty"`
	// Members are the spelled names of the native members, in native order.
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
	// lazy holds, under the mutant "freeze-at-serialization" only, the values
	// that rules returned.
	lazy map[string]lazyPayload
}
type lazyPayload struct {
	codec PayloadCodec
	value any
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

// Serialize writes the plan. Under the mutant it encodes the retained values
// only now, which is what guarantee 5 forbids.
func (p RepresentationPlan) Serialize() ([]byte, error) {
	if !p.Valid() {
		return nil, ErrInvalid
	}
	data := p.data
	if len(p.lazy) > 0 {
		data.Nodes = slices.Clone(data.Nodes)
		for i, n := range data.Nodes {
			if l, ok := p.lazy[n.ID]; ok {
				frozen, err := l.codec.Freeze(l.value)
				if err != nil {
					return nil, err
				}
				data.Nodes[i].Payload = hex.EncodeToString(frozen)
			}
		}
	}
	return json.MarshalIndent(data, "", " ")
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
			for _, arg := range append(slices.Clone(a.Native), a.Adapters...) {
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
// representation. Two uses with equal content share a node. Two different
// contributions never do, because the reference is part of the content.
func nodeIdentity(n PlanNode) string {
	n.ID = ""
	if PrototypeMutants["node-identity-ignores-links"] {
		n.Links, n.Application = nil, nil
	}
	b, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:16])
}

// planBuilder collects the plan while the planner runs.
type planBuilder struct {
	data  planData
	index map[string]int
	lazy  map[string]lazyPayload
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
