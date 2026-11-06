package report

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"
	"unicode"
)

// Dimension metadata belongs to source declarations, not a package/name table.
// This index exists only inside the documentation generator.
type dimensionMetadata struct {
	name   string
	shared bool
}

func dimensionComment(groups ...*ast.CommentGroup) (dimensionMetadata, bool, error) {
	var metadata dimensionMetadata
	found := false
	for _, group := range groups {
		if group == nil {
			continue
		}
		// CommentGroup.Text removes directive-like comments. Read the original
		// tokens instead, retaining both line and block documentation comments.
		for _, comment := range group.List {
			text := strings.TrimPrefix(comment.Text, "//")
			text = strings.TrimSuffix(strings.TrimPrefix(text, "/*"), "*/")
			for _, line := range strings.Split(text, "\n") {
				line = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "*"))
				if !strings.HasPrefix(line, "bittype:") {
					continue
				}
				if found {
					return metadata, false, fmt.Errorf("a type must have only one dimension/shared directive")
				}
				fields := strings.Fields(line)
				switch fields[0] {
				case "bittype:dimension":
					if len(fields) != 2 || !dimensionLabel(fields[1]) {
						return metadata, false, fmt.Errorf("bittype:dimension requires one non-reserved identifier label")
					}
					metadata.name = fields[1]
				case "bittype:shared":
					if len(fields) != 1 {
						return metadata, false, fmt.Errorf("bittype:shared takes no argument")
					}
					metadata.shared = true
				default:
					return metadata, false, fmt.Errorf("unknown dimension directive %q", fields[0])
				}
				found = true
			}
		}
	}
	return metadata, found, nil
}

func dimensionLabel(name string) bool {
	if name == "" || name == "Shared" || name == "Unclassified" || name == "Mixed" {
		return false
	}
	for i, ch := range name {
		if !unicode.IsLetter(ch) && (i == 0 || (!unicode.IsDigit(ch) && ch != '_' && ch != '-')) {
			return false
		}
	}
	return true
}

func (p *program) readDimensions(pkg *sourcePackage) error {
	if p.dimensions == nil {
		p.dimensions = map[*types.TypeName]dimensionMetadata{}
	}
	for _, file := range pkg.files {
		for _, decl := range file.Decls {
			group, ok := decl.(*ast.GenDecl)
			if !ok || group.Tok != token.TYPE {
				continue
			}
			if len(group.Specs) != 1 {
				_, found, err := dimensionComment(group.Doc)
				if err != nil || found {
					return fmt.Errorf("%s: put dimension metadata on individual types in a grouped declaration", p.fset.Position(group.Pos()))
				}
			}
			for _, spec := range group.Specs {
				typeSpec := spec.(*ast.TypeSpec)
				var groupDoc *ast.CommentGroup
				if len(group.Specs) == 1 {
					groupDoc = group.Doc
				}
				metadata, found, err := dimensionComment(groupDoc, typeSpec.Doc)
				if err != nil {
					return fmt.Errorf("%s: %w", p.fset.Position(typeSpec.Pos()), err)
				}
				if found {
					obj, ok := pkg.info.Defs[typeSpec.Name].(*types.TypeName)
					if !ok {
						return fmt.Errorf("%s: no Go type object for dimension metadata", p.fset.Position(typeSpec.Pos()))
					}
					if obj.IsAlias() {
						return fmt.Errorf("%s: aliases inherit dimension metadata from their target type", p.fset.Position(typeSpec.Pos()))
					}
					p.dimensions[obj] = metadata
				}
			}
		}
	}
	return nil
}

func (p *program) dimensionOf(t types.Type) (dimensionMetadata, bool) {
	var obj *types.TypeName
	switch t := t.(type) {
	case *types.Named:
		obj = t.Origin().Obj()
	}
	metadata, found := p.dimensions[obj]
	return metadata, found
}

func (p *program) level(contract *types.Named) string {
	names := map[string]bool{}
	metadata, _ := p.dimensionOf(contract)
	if metadata.name != "" {
		names[metadata.name] = true
	}
	seen := map[types.Type]bool{}
	var visit func(types.Type)
	visit = func(t types.Type) {
		if t == nil || seen[t] {
			return
		}
		seen[t] = true
		if boundary, found := p.dimensionOf(t); found {
			if boundary.name != "" {
				names[boundary.name] = true
			}
			// A classified carrier owns its boundary. Its storage, type arguments
			// and helper methods do not change the declared dimension.
			return
		}
		switch t := t.(type) {
		case *types.Alias:
			visit(t.Rhs())
		case *types.Named:
			for i := 0; i < t.TypeArgs().Len(); i++ {
				visit(t.TypeArgs().At(i))
			}
		case *types.Pointer:
			visit(t.Elem())
		case *types.Slice:
			visit(t.Elem())
		case *types.Array:
			visit(t.Elem())
		case *types.Chan:
			visit(t.Elem())
		case *types.Map:
			visit(t.Key())
			visit(t.Elem())
		case *types.Tuple:
			for i := 0; i < t.Len(); i++ {
				visit(t.At(i).Type())
			}
		case *types.Signature:
			visit(t.Params())
			visit(t.Results())
		}
	}
	iface := contract.Underlying().(*types.Interface)
	for i := 0; i < iface.NumMethods(); i++ {
		visit(iface.Method(i).Type())
	}
	var dimensions []string
	for name := range names {
		dimensions = append(dimensions, name)
	}
	slices.Sort(dimensions)
	switch len(dimensions) {
	case 0:
		if metadata.shared {
			return "Shared"
		}
		return "Unclassified"
	case 1:
		return dimensions[0]
	default:
		return "Mixed: " + strings.Join(dimensions, " + ")
	}
}
