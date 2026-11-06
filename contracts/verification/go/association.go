// Package verification associates already captured artifacts and semantic receipts.
// It has no proof language, checker registry, native judgment or mapping policy.
package verification

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

var (
	ErrInvalid = errors.New("invalid verification association")
	ErrLimit   = errors.New("verification association limit")
)

// Participant references are local to an association, not language identities.
type Participant struct {
	Key      semantics.Key
	Snapshot semantic.Snapshot
	Evidence []semantic.Result
}

// Requirement is fixed by receiver composition, separately from candidate inputs.
type Requirement struct {
	Key      semantics.Key
	Language representation.Language
	Evidence []semantic.EvidenceRequirement
}

type Input struct {
	Source       semantics.Key
	Target       semantics.Key
	Participants []Participant
}

type Site struct {
	Participant semantics.Key
	Node        ast.NodeRef
}

// Limits bound the whole association, including all participants together.
type Limits struct {
	Participants  int
	Evidence      int
	MetadataBytes int
	InputNodes    int
	InputEntries  int
	InputBytes    int
	Sites         int
}

func DefaultLimits() Limits {
	return Limits{Participants: 32, Evidence: 128, MetadataBytes: 1 << 20,
		InputNodes: 200000, InputEntries: 1600000, InputBytes: 32 << 20, Sites: 4096}
}

func (l Limits) Valid() bool {
	return l.Participants > 0 && l.Evidence > 0 && l.MetadataBytes > 0 &&
		l.InputNodes > 0 && l.InputEntries > 0 && l.InputBytes > 0 && l.Sites > 0
}

type state struct {
	source, target semantics.Key
	participants   []Participant
	requirements   []Requirement
	catalogue      map[semantics.Key]map[ast.TreeID]ast.NodeID
	limits         Limits
}

// Association is a value handle over private immutable storage. Its validity
// establishes admission/association, never the truth of any requested claim.
type Association struct{ state *state }

func (a Association) Valid() bool                 { return a.state != nil }
func (a Association) Same(other Association) bool { return a.Valid() && a.state == other.state }
func (a Association) Source() semantics.Key {
	if !a.Valid() {
		return semantics.Key{}
	}
	return a.state.source
}
func (a Association) Target() semantics.Key {
	if !a.Valid() {
		return semantics.Key{}
	}
	return a.state.target
}
func (a Association) Limits() Limits {
	if !a.Valid() {
		return Limits{}
	}
	return a.state.limits
}
func cloneParticipant(p Participant) Participant {
	p.Evidence = slices.Clone(p.Evidence)
	return p
}
func (a Association) Participants() []Participant {
	if !a.Valid() {
		return nil
	}
	result := make([]Participant, len(a.state.participants))
	for i, p := range a.state.participants {
		result[i] = cloneParticipant(p)
	}
	return result
}
func (a Association) Participant(key semantics.Key) (Participant, bool) {
	if a.Valid() {
		for _, p := range a.state.participants {
			if p.Key == key {
				return cloneParticipant(p), true
			}
		}
	}
	return Participant{}, false
}
func cloneRequirements(input []Requirement) []Requirement {
	result := slices.Clone(input)
	for i := range result {
		result[i].Evidence = make([]semantic.EvidenceRequirement, len(input[i].Evidence))
		for j, e := range input[i].Evidence {
			result[i].Evidence[j] = e.Clone()
		}
	}
	return result
}
func (a Association) Requirements() []Requirement {
	if !a.Valid() {
		return nil
	}
	return cloneRequirements(a.state.requirements)
}
func (a Association) Contains(site Site) bool {
	return a.Valid() && site.Node.Tree != "" && site.Node.Node > 0 &&
		site.Node.Node <= a.state.catalogue[site.Participant][site.Node.Tree]
}
func (a Association) CheckSites(sites []Site) error {
	if !a.Valid() {
		return ErrInvalid
	}
	if len(sites) > a.state.limits.Sites {
		return ErrLimit
	}
	for _, site := range sites {
		if !a.Contains(site) {
			return ErrInvalid
		}
	}
	return nil
}

func take(left *int, count int) bool {
	if count < 0 || count > *left {
		return false
	}
	*left -= count
	return true
}
func keyBytes(left *int, key semantics.Key) bool {
	return take(left, len(key.Owner)) && take(left, len(key.Name))
}

// Capture admits receipts against their own snapshots and receiver requirements.
// It performs no resolution, analysis, proof checking, filesystem or network IO.
// Input and requirement construction storage must not be mutated during this call.
func Capture(ctx context.Context, input Input, requirements []Requirement, limits Limits) (Association, error) {
	if ctx == nil {
		return Association{}, ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return Association{}, err
	}
	if limits == (Limits{}) {
		limits = DefaultLimits()
	}
	if !limits.Valid() || !input.Source.Valid() || !input.Target.Valid() ||
		len(input.Participants) == 0 || len(input.Participants) != len(requirements) {
		return Association{}, ErrInvalid
	}
	if len(input.Participants) > limits.Participants {
		return Association{}, ErrLimit
	}
	metadata, receipts := limits.MetadataBytes, limits.Evidence
	if !keyBytes(&metadata, input.Source) || !keyBytes(&metadata, input.Target) {
		return Association{}, ErrLimit
	}
	required := make(map[semantics.Key]Requirement, len(requirements))
	for _, r := range requirements {
		if !r.Key.Valid() || r.Language.Validate() != nil {
			return Association{}, ErrInvalid
		}
		if _, exists := required[r.Key]; exists {
			return Association{}, ErrInvalid
		}
		if !keyBytes(&metadata, r.Key) || !take(&metadata, len(r.Language.ID())) ||
			!take(&metadata, len(r.Language.Name())) || !take(&metadata, len(r.Language.Profile())) ||
			!take(&receipts, len(r.Evidence)) {
			return Association{}, ErrLimit
		}
		if err := semantic.ValidateEvidenceRequirements(r.Evidence); err != nil {
			return Association{}, err
		}
		for _, e := range r.Evidence {
			if !keyBytes(&metadata, e.Operation) || !take(&metadata, len(e.Profile)) {
				return Association{}, ErrLimit
			}
			for _, setting := range e.Configuration {
				if !keyBytes(&metadata, setting.Key) || !take(&metadata, len(setting.Value)) {
					return Association{}, ErrLimit
				}
			}
			for _, report := range e.Reports {
				if !keyBytes(&metadata, report.Obligation) {
					return Association{}, ErrLimit
				}
			}
		}
		required[r.Key] = r
	}
	s := &state{source: input.Source, target: input.Target, limits: limits,
		catalogue: make(map[semantics.Key]map[ast.TreeID]ast.NodeID, len(input.Participants))}
	nodes, entries, inputBytes := limits.InputNodes, limits.InputEntries, limits.InputBytes
	for _, p := range input.Participants {
		if err := ctx.Err(); err != nil {
			return Association{}, err
		}
		r, found := required[p.Key]
		if !found || !p.Snapshot.Valid() || !r.Language.Compatible(p.Snapshot.Language()) {
			return Association{}, ErrInvalid
		}
		if _, exists := s.catalogue[p.Key]; exists {
			return Association{}, ErrInvalid
		}
		if len(p.Evidence) != len(r.Evidence) {
			return Association{}, ErrInvalid
		}
		if !keyBytes(&metadata, p.Key) {
			return Association{}, ErrLimit
		}
		catalogue := make(map[ast.TreeID]ast.NodeID)
		for _, tree := range p.Snapshot.Inputs() {
			info, err := tree.Inspect()
			if err != nil {
				return Association{}, err
			}
			if !take(&nodes, info.Nodes) || !take(&entries, info.Entries) || !take(&inputBytes, info.TextBytes) {
				return Association{}, ErrLimit
			}
			catalogue[info.ID] = ast.NodeID(info.Nodes)
		}
		if files := p.Snapshot.Files(); files != nil {
			for _, file := range files.Files.Entries {
				if !take(&entries, 1) || !take(&inputBytes, len(file.ID)) ||
					!take(&inputBytes, len(file.Path)) || !take(&inputBytes, len(file.Text)) {
					return Association{}, ErrLimit
				}
			}
			for _, link := range files.Links {
				if !take(&entries, 1) || !take(&inputBytes, len(link.File)) ||
					!take(&inputBytes, len(link.Tree)) || !take(&inputBytes, len(link.Source)) {
					return Association{}, ErrLimit
				}
			}
		}
		evidence, err := semantic.AdmitEvidence(p.Snapshot, r.Evidence, p.Evidence)
		if err != nil {
			return Association{}, fmt.Errorf("participant %s/%s: %w", p.Key.Owner, p.Key.Name, err)
		}
		s.participants = append(s.participants, Participant{Key: p.Key, Snapshot: p.Snapshot, Evidence: evidence})
		s.catalogue[p.Key] = catalogue
	}
	if _, ok := s.catalogue[input.Source]; !ok {
		return Association{}, ErrInvalid
	}
	if _, ok := s.catalogue[input.Target]; !ok {
		return Association{}, ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return Association{}, err
	}
	s.requirements = cloneRequirements(requirements)
	return Association{state: s}, nil
}
