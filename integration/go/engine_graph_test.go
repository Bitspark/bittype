package integration_test

import (
	"errors"
	"github.com/Bitspark/bittype/trees/representation/go"
	"reflect"
	"testing"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/engine/go"
)

// Exhaustively enumerate small directed graphs and compare the planner with
// independently enumerated simple paths, rather than duplicating its BFS.
func TestEngineGraphAgainstSimplePathOracle(t *testing.T) {
	ids := []language.ID{"one", "two", "three", "four"}
	var pairs [][2]int
	for i := range ids {
		for j := range ids {
			if i != j {
				pairs = append(pairs, [2]int{i, j})
			}
		}
	}
	definitions := make([]representation.Language, len(ids))
	for i, id := range ids {
		definitions[i] = definition[echo](t, id)
	}
	for mask := 0; mask < (1 << len(pairs)); mask++ {
		e := engine.New()
		register(t, e, definitions...)
		var adjacent [4][4]bool
		for bit, pair := range pairs {
			if mask&(1<<bit) != 0 {
				adjacent[pair[0]][pair[1]] = true
				connect(t, e, &directTransformer{from: ids[pair[0]], to: ids[pair[1]], run: func(representation.AST) (representation.AST, error) {
					return representation.AST{}, errors.New("oracle only")
				}})
			}
		}
		for from := range ids {
			for to := range ids {
				best, count := 5, 0
				var expected []language.ID
				var visit func(int, []language.ID, uint8)
				visit = func(current int, path []language.ID, seen uint8) {
					if current == to {
						if len(path) < best {
							best, count = len(path), 1
							expected = append([]language.ID(nil), path...)
						} else if len(path) == best {
							count++
						}
						return
					}
					for next := range ids {
						if adjacent[current][next] && seen&(1<<next) == 0 {
							visit(next, append(path, ids[next]), seen|(1<<next))
						}
					}
				}
				visit(from, []language.ID{ids[from]}, 1<<from)
				plan, err := e.Plan(ids[from], ids[to])
				switch {
				case count == 0:
					if !errors.Is(err, engine.ErrNoPath) {
						t.Fatalf("graph %d %d -> %d: %v", mask, from, to, err)
					}
				case count > 1:
					if !errors.Is(err, engine.ErrAmbiguousPath) {
						t.Fatalf("graph %d ambiguity lost: %v", mask, err)
					}
				default:
					if err != nil || !reflect.DeepEqual(plan.Languages(), expected) {
						t.Fatalf("graph %d route %v expected %v: %v", mask, plan.Languages(), expected, err)
					}
				}
			}
		}
	}
}
