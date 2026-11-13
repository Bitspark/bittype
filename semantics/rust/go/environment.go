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
// environment. Equal environments give identical settings. A value that is no
// Rust value at all is refused with semantic.ErrInvalidRequest, and a value
// outside this profile with ErrUnsupported. Validity is judged first, over
// every field, so an environment with both kinds of value is invalid.
func ConfigurationFor(environment Environment) ([]semantic.Setting, error) {
	if err := environment.check(); err != nil {
		return nil, err
	}
	return semantic.CaptureSettings(environment.settings(), semantic.DefaultLimits())
}

// EnvironmentOf reads the environment that settings were built from. It
// accepts exactly what ConfigurationFor returns. Settings in another form,
// which ConfigurationFor returns for no environment, are refused with
// semantic.ErrInvalidRequest before any value is judged. Settings in that form
// are refused as ConfigurationFor refuses their environment.
func EnvironmentOf(settings []semantic.Setting) (Environment, error) {
	var environment Environment
	for _, setting := range settings {
		if setting.Key == Key(environmentKey) {
			// A value that does not decode, or decodes loosely, leaves an
			// environment whose settings differ from these. The comparison
			// below is the one check of the form and refuses it.
			_ = json.Unmarshal(setting.Value, &environment)
		}
	}
	if !reflect.DeepEqual(settings, environment.settings()) {
		return Environment{}, semantic.ErrInvalidRequest
	}
	if err := environment.check(); err != nil {
		return Environment{}, err
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

// settings is the one form of an environment as settings. It judges no value.
func (e Environment) settings() []semantic.Setting {
	if e.Dependencies == nil {
		e.Dependencies = []Dependency{}
	}
	// A value of strings and slices of strings always encodes.
	data, _ := json.Marshal(e)
	return []semantic.Setting{
		{Key: Key("compiler"), Value: []byte(CompilerVersion)},
		{Key: Key(environmentKey), Value: data},
		{Key: Key("helper"), Value: []byte(HelperProfile)},
	}
}

func (e Environment) check() error {
	fields := []struct{ name, value, supported string }{
		{"edition", e.Edition, Edition2024},
		{"target", e.Target, TargetLinuxX64},
		{"crate kind", e.Crate.Kind, KindLibrary},
	}
	for _, field := range fields {
		if field.value == "" {
			return fmt.Errorf("%w: the Rust environment states no %s", semantic.ErrInvalidRequest, field.name)
		}
	}
	if e.Crate.Name == "" {
		return fmt.Errorf("%w: the Rust environment states no crate name", semantic.ErrInvalidRequest)
	}
	// The compiler takes letters, digits and underscores as a crate name.
	for _, r := range e.Crate.Name {
		if r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return fmt.Errorf("%w: %q is no Rust crate name", semantic.ErrInvalidRequest, e.Crate.Name)
		}
	}
	for _, field := range fields {
		if field.value != field.supported {
			return fmt.Errorf("%w: %s %q, this profile covers %q", ErrUnsupported, field.name, field.value, field.supported)
		}
	}
	if !identifier(e.Crate.Name) || keywords[e.Crate.Name] {
		return fmt.Errorf("%w: crate name %q, this profile covers an ASCII identifier that is no keyword", ErrUnsupported, e.Crate.Name)
	}
	if len(e.Dependencies) != 0 {
		return fmt.Errorf("%w: %d dependency crates, this profile covers the sysroot only", ErrUnsupported, len(e.Dependencies))
	}
	return nil
}

// identifier reports an ASCII identifier other than a lone underscore: a name
// that another crate can write as the first segment of a path.
func identifier(name string) bool {
	for i, r := range name {
		if r != '_' && (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (i == 0 || r < '0' || r > '9') {
			return false
		}
	}
	return name != "_"
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
