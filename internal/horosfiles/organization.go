package horosfiles

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/internal/tree"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/source/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"io"
	"sort"
	"strings"
)

// Placement is an explicit same-language organization decision. Modules are
// canonical module pins, never aliases selected by spelling or filesystem lookup.
type Placement struct {
	ID      files.ID
	Path    files.Path
	Modules []string
}

func OrganizationSetting() semantics.Key {
	return semantics.Key{Owner: "horos/files", Name: "organization-json-1"}
}

type moduleSyntax struct {
	tree ast.AST[s.Language]
	node ast.NodeID
}

func Materialize(snapshot semantic.Snapshot, settings []semantic.Setting, limits source.Limits) ([]Placement, []ast.AST[s.Language], files.ID, error) {
	fail := func(err error) ([]Placement, []ast.AST[s.Language], files.ID, error) { return nil, nil, "", err }
	modules := map[string]moduleSyntax{}
	var inputTrees []ast.AST[s.Language]
	rootPin := ""
	for _, envelope := range snapshot.Inputs() {
		input, err := representation.Unpack[s.Language](snapshot.Language(), envelope)
		if err != nil {
			return fail(err)
		}
		inputTrees = append(inputTrees, input)
		if input.ID == snapshot.Entry().Tree {
			if snapshot.Entry().Node != input.Root {
				return fail(source.ErrInvalid)
			}
			rootPin = RootPin(input)
		}
		for _, child := range input.Nodes[input.Root-1].Children {
			pin := attr(input, child.Node, s.AttributePin)
			if _, found := modules[pin]; found {
				return fail(source.ErrInvalid)
			}
			modules[pin] = moduleSyntax{input, child.Node}
		}
	}
	if rootPin == "" {
		return fail(source.ErrInvalid)
	}
	var placements []Placement
	if len(settings) > 0 {
		// A refusal of the setting is invalid input in the Files category. The
		// value returned before stays joined to it.
		if len(settings) != 1 || settings[0].Key != OrganizationSetting() {
			return fail(fmt.Errorf("%w: %w", source.ErrInvalid, semantic.ErrInvalidRequest))
		}
		decoder := json.NewDecoder(bytes.NewReader(settings[0].Value))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&placements); err != nil {
			return fail(fmt.Errorf("%w: %w", source.ErrInvalid, err))
		}
		var extra any
		if err := decoder.Decode(&extra); err != io.EOF {
			return fail(fmt.Errorf("%w: %w", source.ErrInvalid, semantic.ErrInvalidRequest))
		}
		if repeatedMember(settings[0].Value) {
			return fail(fmt.Errorf("%w: a member of the organization is named twice", source.ErrInvalid))
		}
	} else if context := snapshot.Files(); context != nil {
		byTree := map[ast.TreeID]files.File{}
		for _, link := range context.Links {
			for _, file := range context.Files.Entries {
				if link.File == file.ID {
					byTree[link.Tree] = file
				}
			}
		}
		for _, input := range inputTrees {
			file := byTree[input.ID]
			placements = append(placements, Placement{file.ID, file.Path, ModulePins(input)})
		}
	} else {
		var pins []string
		for pin := range modules {
			pins = append(pins, pin)
		}
		sort.Strings(pins)
		for _, pin := range pins {
			placements = append(placements, Placement{files.ID("module:" + pin), files.Path("modules/" + pin + ".horos"), []string{pin}})
		}
	}
	if len(placements) > limits.Files.MaxFiles {
		return fail(source.ErrLimit)
	}
	inventory := files.Files{}
	for _, placement := range placements {
		inventory.Entries = append(inventory.Entries, files.File{ID: placement.ID, Path: placement.Path})
	}
	if _, err := inventory.Capture(limits.Files); err != nil {
		return fail(filesCategory(err))
	}
	seen := map[string]bool{}
	var output []ast.AST[s.Language]
	var entry files.ID
	work := limits.MaxWork
	for _, placement := range placements {
		if len(placement.Modules) == 0 {
			return fail(source.ErrInvalid)
		}
		out := ast.AST[s.Language]{ID: ast.TreeID("horos/files/" + string(placement.ID))}
		rootModule := ast.NodeID(0)
		var children []ast.Child
		for _, pin := range placement.Modules {
			module, exists := modules[pin]
			if !exists || seen[pin] {
				return fail(source.ErrInvalid)
			}
			seen[pin] = true
			var copyNode func(ast.NodeID, int) (ast.NodeID, error)
			copyNode = func(id ast.NodeID, depth int) (ast.NodeID, error) {
				work--
				if work < 0 || depth > limits.MaxDepth || len(out.Nodes) >= limits.MaxNodes {
					return 0, source.ErrLimit
				}
				node := module.tree.Nodes[id-1]
				out.Nodes = append(out.Nodes, ast.Node[s.Language]{})
				assigned := ast.NodeID(len(out.Nodes))
				// These are fresh current nodes. Resolution re-establishes optional
				// semantic annotations; no old receipt or historical origins are copied.
				node.References = nil
				node.Origin = ast.Origin{}
				node.Attributes = append([]ast.Attribute(nil), node.Attributes...)
				oldChildren := node.Children
				node.Children = nil
				for _, child := range oldChildren {
					next, err := copyNode(child.Node, depth+1)
					if err != nil {
						return 0, err
					}
					node.Children = append(node.Children, ast.Child{Role: child.Role, Node: next})
				}
				out.Nodes[assigned-1] = node
				return assigned, nil
			}
			moduleID, err := copyNode(module.node, 1)
			if err != nil {
				return fail(err)
			}
			children = append(children, tree.Child(s.RoleModule, moduleID))
			if rootModule == 0 {
				rootModule = moduleID
			}
			if pin == rootPin {
				rootModule = moduleID
				entry = placement.ID
			}
		}
		out.Nodes = append(out.Nodes, ast.Node[s.Language]{Kind: s.KindProgram, Children: children, References: []ast.Reference{{Role: s.RoleRoot, Target: ast.NodeRef{Tree: out.ID, Node: rootModule}}}})
		out.Root = ast.NodeID(len(out.Nodes))
		output = append(output, out)
	}
	if len(seen) != len(modules) || entry == "" {
		return fail(source.ErrInvalid)
	}
	return placements, output, entry, nil
}

// filesCategory joins the Files category to a value of the files package, so
// that errors.Is holds for both. It returns any other error as it is.
func filesCategory(err error) error {
	switch {
	case errors.Is(err, files.ErrLimit):
		return fmt.Errorf("%w: %w", source.ErrLimit, err)
	case errors.Is(err, files.ErrInvalid):
		return fmt.Errorf("%w: %w", source.ErrInvalid, err)
	}
	return err
}

// repeatedMember reports a JSON object that names one member twice. The decoder
// keeps the last value, so the first would be accepted and ignored. Names are
// compared as the decoder matches them, without regard to case. The caller has
// decoded the value into its type, which bounds the nesting.
func repeatedMember(data []byte) bool {
	decoder := json.NewDecoder(bytes.NewReader(data))
	repeated := false
	var value func()
	value = func() {
		token, err := decoder.Token()
		if err != nil {
			return
		}
		switch token {
		case json.Delim('{'):
			var names []string
			for decoder.More() {
				key, _ := decoder.Token()
				name, _ := key.(string)
				for _, seen := range names {
					repeated = repeated || strings.EqualFold(seen, name)
				}
				names = append(names, name)
				value()
			}
			decoder.Token()
		case json.Delim('['):
			for decoder.More() {
				value()
			}
			decoder.Token()
		}
	}
	value()
	return repeated
}

func ProgramIdentity(snapshot semantic.Snapshot) (string, []string, error) {
	var root string
	var pins []string
	for _, envelope := range snapshot.Inputs() {
		input, err := representation.Unpack[s.Language](snapshot.Language(), envelope)
		if err != nil {
			return "", nil, err
		}
		pins = append(pins, ModulePins(input)...)
		if input.ID == snapshot.Entry().Tree {
			root = RootPin(input)
		}
	}
	sort.Strings(pins)
	return root, pins, nil
}
