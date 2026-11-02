package horosbridge

import (
	"strconv"
	"strings"
	"testing"

	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func TestParameterIndexDecodeHasMachineIntegerWorkBound(t *testing.T) {
	decode := func(spelling string) (bool, error) {
		input := s.AST{ID: "index", Root: 1, Nodes: []s.Node{{
			Kind: s.KindParameter, Attributes: []ast.Attribute{{Name: s.AttributeIndex, Text: spelling}},
		}}}
		value, err := (reader{tree: input}).value(1)
		return value != nil, err
	}
	for _, spelling := range []string{"0", strconv.Itoa(int(^uint(0) >> 1))} {
		if present, err := decode(spelling); !present || err != nil {
			t.Fatalf("representable slot rejected: %q: %v", spelling, err)
		}
	}
	// Horos's declaration slot API is int-valued. Huge decimal spellings cannot
	// denote a valid larger slot and must fail before constructing an Ontos int.
	if present, err := decode(strings.Repeat("9", 64<<10)); present || err == nil {
		t.Fatal("unrepresentable slot was converted instead of rejected at the decoding boundary")
	}
}
