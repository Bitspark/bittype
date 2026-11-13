package rust

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unicode"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
)

// Environment states what a captured Rust crate is compiled under. A consumer
// states every field, and no field has a default. The compiler is pinned by
// the profile and is not a field. The sysroot is a path given to the resolver
// and has no identity here.
//
// Each field admits the one value below that this profile covers. Another
// value is refused with ErrUnsupported, and none is accepted and ignored.
// Later coverage admits more values of the same fields.
type Environment struct {
	// Edition is the Rust edition of the crate.
	Edition string
	// Target is the target triple the crate is checked for.
	Target string
	// Crate is the identity and the kind of the crate.
	Crate Crate
	// Dependencies are the crates the captured crate may name beside the
	// sysroot. This profile admits none.
	Dependencies []Dependency
}

// Crate is the identity and the kind of the captured crate. The compiler
// analyzes under Name, so every compiler path of the model starts with it.
type Crate struct {
	Name string
	Kind string
}

// Dependency is one crate that the captured crate may name. Name is the name
// the captured crate uses for it. No profile admits a dependency yet, so the
// rest of its description is not fixed.
type Dependency struct {
	Name string
}

// The values of an Environment that this profile covers.
const (
	Edition2024    = "2024"
	TargetLinuxX64 = "x86_64-unknown-linux-gnu"
	KindLibrary    = "lib"
)

const environmentKey = "environment-1"

// ConfigurationFor is the configuration of resolution and analysis under the
// environment. Equal environments give identical settings. A value outside
// this profile is refused with ErrUnsupported, and a value that is no Rust
// value at all with semantic.ErrInvalidRequest.
func ConfigurationFor(environment Environment) ([]semantic.Setting, error) {
	if err := environment.check(); err != nil {
		return nil, err
	}
	environment.Dependencies = []Dependency{}
	data, err := json.Marshal(environment)
	if err != nil {
		return nil, err
	}
	return semantic.CaptureSettings([]semantic.Setting{
		{Key: Key("compiler"), Value: []byte(CompilerVersion)},
		{Key: Key(environmentKey), Value: data},
		{Key: Key("helper"), Value: []byte(HelperProfile)},
	}, semantic.DefaultLimits())
}

// EnvironmentOf reads the environment that settings were built from. It
// accepts exactly what ConfigurationFor returns and refuses everything else
// with the error ConfigurationFor would give, or with
// semantic.ErrInvalidRequest for settings it cannot have returned.
func EnvironmentOf(settings []semantic.Setting) (Environment, error) {
	var environment Environment
	for _, setting := range settings {
		if setting.Key == Key(environmentKey) {
			// A value that does not decode, or decodes loosely, leaves an
			// environment whose settings differ from these. The comparison
			// below is the one check and refuses it.
			_ = json.Unmarshal(setting.Value, &environment)
		}
	}
	want, err := ConfigurationFor(environment)
	if err != nil {
		return Environment{}, err
	}
	if !reflect.DeepEqual(settings, want) {
		return Environment{}, semantic.ErrInvalidRequest
	}
	environment.Dependencies = nil
	return environment, nil
}

// EvidenceFor is the analysis evidence requirement for exactly this
// configuration. There is no requirement without an environment.
func EvidenceFor(settings []semantic.Setting) (semantic.EvidenceRequirement, error) {
	if _, err := EnvironmentOf(settings); err != nil {
		return semantic.EvidenceRequirement{}, err
	}
	requirement := semantic.EvidenceRequirement{Operation: analysis.Kind(), Profile: AnalysisProfile, Configuration: settings, Reports: Requirements(), Semantics: true}
	return requirement.Clone(), nil
}

func (e Environment) check() error {
	for _, field := range []struct{ name, value, supported string }{
		{"edition", e.Edition, Edition2024},
		{"target", e.Target, TargetLinuxX64},
		{"crate kind", e.Crate.Kind, KindLibrary},
	} {
		if field.value == "" {
			return fmt.Errorf("%w: the Rust environment states no %s", semantic.ErrInvalidRequest, field.name)
		}
		if field.value != field.supported {
			return fmt.Errorf("%w: %s %q, this profile covers %q", ErrUnsupported, field.name, field.value, field.supported)
		}
	}
	if err := checkCrateName(e.Crate.Name); err != nil {
		return err
	}
	if len(e.Dependencies) != 0 {
		return fmt.Errorf("%w: %d dependency crates, this profile covers the sysroot only", ErrUnsupported, len(e.Dependencies))
	}
	return nil
}

// checkCrateName admits a name that another crate can write as the first
// segment of a path. The compiler itself takes any letters, digits and
// underscores. Such a name outside the admitted ones is unsupported, and
// anything else is no crate name.
func checkCrateName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: the Rust environment states no crate name", semantic.ErrInvalidRequest)
	}
	identifier := name != "_"
	for i, r := range name {
		if r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return fmt.Errorf("%w: %q is no Rust crate name", semantic.ErrInvalidRequest, name)
		}
		ascii := r == '_' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || i > 0 && r >= '0' && r <= '9'
		identifier = identifier && ascii
	}
	if !identifier || keywords[name] {
		return fmt.Errorf("%w: crate name %q, this profile covers an ASCII identifier that is no keyword", ErrUnsupported, name)
	}
	return nil
}

// The strict and the reserved keywords of edition 2024. None can start a path
// that names a crate.
var keywords = func() map[string]bool {
	out := map[string]bool{}
	for _, word := range []string{
		"as", "async", "await", "break", "const", "continue", "crate", "dyn", "else", "enum", "extern", "false", "fn", "for", "if", "impl", "in", "let", "loop", "match", "mod", "move", "mut", "pub", "ref", "return", "self", "Self", "static", "struct", "super", "trait", "true", "type", "unsafe", "use", "where", "while",
		"abstract", "become", "box", "do", "final", "gen", "macro", "override", "priv", "try", "typeof", "unsized", "virtual", "yield",
	} {
		out[word] = true
	}
	return out
}()
