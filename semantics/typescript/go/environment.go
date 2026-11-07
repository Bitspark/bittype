package typescript

import (
	"bytes"
	"encoding/json"
	"reflect"
	"slices"
	"sort"
	"unicode/utf8"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	lit "github.com/Bitspark/bittype/internal/typescript"
)

// PackageTypes is the explicitly selected declaration entry of a public package
// export. Its text participates in receipt identity; it is never an output file.
type PackageTypes struct {
	Module      string
	Declaration string
}

func ImportConfiguration(packages []PackageTypes) ([]semantic.Setting, error) {
	if len(packages) > 64 {
		return nil, semantic.ErrLimit
	}
	packages = slices.Clone(packages)
	if packages == nil {
		packages = []PackageTypes{}
	}
	sort.Slice(packages, func(i, j int) bool { return packages[i].Module < packages[j].Module })
	size := 0
	for i, p := range packages {
		size += len(p.Module) + len(p.Declaration)
		if size > 768<<10 {
			return nil, semantic.ErrLimit
		}
		if !lit.PackageSpecifier(p.Module) || p.Declaration == "" || !utf8.ValidString(p.Declaration) || i > 0 && p.Module == packages[i-1].Module {
			return nil, semantic.ErrInvalidRequest
		}
	}
	data, err := json.Marshal(packages)
	if err != nil {
		return nil, err
	}
	return semantic.CaptureSettings(append(Configuration(), semantic.Setting{Key: Key("named-type-imports-1"), Value: data}), semantic.DefaultLimits())
}

// Imports decodes the receipt-bound environment. Enabled with an empty list
// authorizes local imports only. The original captured program profile is distinct.
func Imports(settings []semantic.Setting) (packages []PackageTypes, enabled bool, err error) {
	base := Configuration()
	if reflect.DeepEqual(settings, base) {
		return nil, false, nil
	}
	if len(settings) != len(base)+1 || !reflect.DeepEqual(settings[:len(base)], base) {
		return nil, false, semantic.ErrInvalidRequest
	}
	last := settings[len(base)]
	if last.Key != Key("named-type-imports-1") || len(last.Value) > semantic.DefaultLimits().MaxConfigurationBytes {
		return nil, false, semantic.ErrInvalidRequest
	}
	if err := json.Unmarshal(last.Value, &packages); err != nil {
		return nil, false, semantic.ErrInvalidRequest
	}
	want, err := ImportConfiguration(packages)
	if err != nil {
		return nil, false, err
	}
	if !bytes.Equal(want[len(base)].Value, last.Value) {
		return nil, false, semantic.ErrInvalidRequest
	}
	return packages, true, nil
}
func EvidenceFor(settings []semantic.Setting) (semantic.EvidenceRequirement, error) {
	if !MatchesConfiguration(settings) {
		return semantic.EvidenceRequirement{}, semantic.ErrInvalidRequest
	}
	contract := Evidence()
	contract.Configuration = settings
	return contract.Clone(), nil
}
