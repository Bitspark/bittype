package engine

import (
	"fmt"
	"sort"

	"github.com/Bitspark/bittype/language/go"
)

// Plan is an immutable route, not a claim that an input's semantics can be
// mapped. Execute revalidates it against the executing engine's registry.
// A zero plan is invalid; a one-language plan is a checked identity operation.
type Plan struct{ route []language.ID }

func (p Plan) Languages() []language.ID { return append([]language.ID(nil), p.route...) }
func (p Plan) FromID() language.ID {
	if len(p.route) == 0 {
		return ""
	}
	return p.route[0]
}
func (p Plan) ToID() language.ID {
	if len(p.route) == 0 {
		return ""
	}
	return p.route[len(p.route)-1]
}

// Plan selects the unique shortest directed path. Path counts saturate at two,
// so discovering ambiguity does not enumerate exponentially many routes.
func (e *Engine) Plan(from, to language.ID) (Plan, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, id := range []language.ID{from, to} {
		if _, ok := e.languages[id]; !ok {
			return Plan{}, fmt.Errorf("%w: %s", ErrUnknownLanguage, id)
		}
	}
	distance := map[language.ID]int{from: 0}
	count := map[language.ID]int{from: 1}
	parent := map[language.ID]language.ID{}
	queue := []language.ID{from}
	for head := 0; head < len(queue); head++ {
		current := queue[head]
		next := make([]language.ID, 0, len(e.edges[current]))
		for id := range e.edges[current] {
			next = append(next, id)
		}
		sort.Slice(next, func(i, j int) bool { return next[i] < next[j] })
		for _, id := range next {
			proposed := distance[current] + 1
			known, seen := distance[id]
			if !seen {
				distance[id], count[id], parent[id] = proposed, count[current], current
				queue = append(queue, id)
			} else if known == proposed {
				count[id] += count[current]
				if count[id] > 2 {
					count[id] = 2
				}
			}
		}
	}
	if count[to] == 0 {
		return Plan{}, fmt.Errorf("%w: %s -> %s", ErrNoPath, from, to)
	}
	if count[to] > 1 {
		return Plan{}, fmt.Errorf("%w: %s -> %s", ErrAmbiguousPath, from, to)
	}
	route := make([]language.ID, distance[to]+1)
	current := to
	for i := len(route) - 1; i >= 0; i-- {
		route[i] = current
		current = parent[current]
	}
	return Plan{route: route}, nil
}

// PlanVia selects an explicit route including both endpoints. It may be longer
// than the shortest route, but repeated languages and missing edges are errors.
func (e *Engine) PlanVia(route ...language.ID) (Plan, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if _, _, err := e.resolve(route); err != nil {
		return Plan{}, err
	}
	return Plan{route: append([]language.ID(nil), route...)}, nil
}

// resolve requires a held read lock and returns a snapshot for execution.
func (e *Engine) resolve(route []language.ID) ([]registration, []edge, error) {
	if len(route) == 0 || len(route) > len(e.languages) {
		return nil, nil, ErrInvalidRoute
	}
	seen := make(map[language.ID]bool, len(route))
	registrations := make([]registration, len(route))
	edges := make([]edge, 0, len(route)-1)
	for i, id := range route {
		if seen[id] {
			return nil, nil, fmt.Errorf("%w: repeated %s", ErrInvalidRoute, id)
		}
		seen[id] = true
		r, ok := e.languages[id]
		if !ok {
			return nil, nil, fmt.Errorf("%w: %s", ErrUnknownLanguage, id)
		}
		registrations[i] = r
		if i > 0 {
			step, ok := e.edges[route[i-1]][id]
			if !ok {
				return nil, nil, fmt.Errorf("%w: %s -> %s", ErrNoPath, route[i-1], id)
			}
			edges = append(edges, step)
		}
	}
	return registrations, edges, nil
}
