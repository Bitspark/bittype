package rust_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	n "github.com/Bitspark/bittype/semantics/rust/go"
)

func supported() n.Environment {
	return n.Environment{Edition: n.Edition2024, Target: n.TargetLinuxX64, Crate: n.Crate{Name: "pages", Kind: n.KindLibrary}}
}

// The configuration of an environment is part of every receipt. These are its
// exact keys and bytes.
func TestEnvironmentConfigurationBytes(t *testing.T) {
	settings, err := n.ConfigurationFor(supported())
	if err != nil {
		t.Fatal(err)
	}
	want := []semantic.Setting{
		{Key: n.Key("compiler"), Value: []byte("1.90.0")},
		{Key: n.Key("environment-1"), Value: []byte(`{"Edition":"2024","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[]}`)},
		{Key: n.Key("helper"), Value: []byte("captured-types-2")},
	}
	if !reflect.DeepEqual(settings, want) {
		t.Fatalf("%s", settings)
	}
	// A nil and an empty dependency collection are one environment.
	empty := supported()
	empty.Dependencies = []n.Dependency{}
	if again, err := n.ConfigurationFor(empty); err != nil || !reflect.DeepEqual(again, want) {
		t.Fatal(again, err)
	}
	environment, err := n.EnvironmentOf(settings)
	if err != nil || !reflect.DeepEqual(environment, supported()) {
		t.Fatal(environment, err)
	}
	other := supported()
	other.Crate.Name = "Other_2"
	different, err := n.ConfigurationFor(other)
	if err != nil || reflect.DeepEqual(different, settings) {
		t.Fatal("the crate name is not part of the configuration", different, err)
	}
	if n.ResolutionProfile != "rustc-1.90.0/captured-types-2/resolution-1" || n.AnalysisProfile != "rustc-1.90.0/captured-types-2/analysis-1" {
		t.Fatal(n.ResolutionProfile, n.AnalysisProfile)
	}
}

// A value outside the profile is unsupported coverage, a value that is no
// Rust value is invalid input, and no error is both. No field is accepted and
// ignored: every field has a refused value.
func TestEnvironmentRefusals(t *testing.T) {
	for name, c := range map[string]struct {
		change func(*n.Environment)
		want   error
	}{
		"edition 2021":                          {func(e *n.Environment) { e.Edition = "2021" }, n.ErrUnsupported},
		"an edition yet to come":                {func(e *n.Environment) { e.Edition = "2027" }, n.ErrUnsupported},
		"another target":                        {func(e *n.Environment) { e.Target = "aarch64-unknown-linux-gnu" }, n.ErrUnsupported},
		"a binary crate":                        {func(e *n.Environment) { e.Crate.Kind = "bin" }, n.ErrUnsupported},
		"one dependency":                        {func(e *n.Environment) { e.Dependencies = []n.Dependency{{Name: "serde"}} }, n.ErrUnsupported},
		"a keyword as crate name":               {func(e *n.Environment) { e.Crate.Name = "crate" }, n.ErrUnsupported},
		"a crate name of digits":                {func(e *n.Environment) { e.Crate.Name = "1abc" }, n.ErrUnsupported},
		"an underscore":                         {func(e *n.Environment) { e.Crate.Name = "_" }, n.ErrUnsupported},
		"a crate name of letters outside ASCII": {func(e *n.Environment) { e.Crate.Name = "über" }, n.ErrUnsupported},
		"no edition":                            {func(e *n.Environment) { e.Edition = "" }, semantic.ErrInvalidRequest},
		"no target":                             {func(e *n.Environment) { e.Target = "" }, semantic.ErrInvalidRequest},
		"no crate kind":                         {func(e *n.Environment) { e.Crate.Kind = "" }, semantic.ErrInvalidRequest},
		"no crate name":                         {func(e *n.Environment) { e.Crate.Name = "" }, semantic.ErrInvalidRequest},
		"a hyphen in the crate name":            {func(e *n.Environment) { e.Crate.Name = "my-pages" }, semantic.ErrInvalidRequest},
		"an option as crate name":               {func(e *n.Environment) { e.Crate.Name = "--cfg" }, semantic.ErrInvalidRequest},
		// Validity is judged before coverage, over every field.
		"edition 2021 and no target":           {func(e *n.Environment) { e.Edition, e.Target = "2021", "" }, semantic.ErrInvalidRequest},
		"a binary crate and no crate name":     {func(e *n.Environment) { e.Crate = n.Crate{Kind: "bin"} }, semantic.ErrInvalidRequest},
		"another target and a hyphenated name": {func(e *n.Environment) { e.Target, e.Crate.Name = "wasm32-unknown-unknown", "my-pages" }, semantic.ErrInvalidRequest},
		"one dependency and no edition":        {func(e *n.Environment) { e.Edition, e.Dependencies = "", []n.Dependency{{Name: "serde"}} }, semantic.ErrInvalidRequest},
		"the zero environment":                 {func(e *n.Environment) { *e = n.Environment{} }, semantic.ErrInvalidRequest},
	} {
		environment := supported()
		c.change(&environment)
		settings, err := n.ConfigurationFor(environment)
		if settings != nil || !errors.Is(err, c.want) {
			t.Errorf("%s: %v, want %v", name, err, c.want)
		}
		if errors.Is(err, n.ErrUnsupported) && errors.Is(err, semantic.ErrInvalidRequest) {
			t.Errorf("%s: both unsupported and invalid: %v", name, err)
		}
	}
	for _, name := range []string{"pages", "Pages", "_pages", "pages_2", "std", "union"} {
		environment := supported()
		environment.Crate.Name = name
		if _, err := n.ConfigurationFor(environment); err != nil {
			t.Errorf("crate name %s: %v", name, err)
		}
	}
}

// Evidence and the analysis contract exist only for the settings of a
// supported environment.
func TestEvidenceRequiresAnEnvironment(t *testing.T) {
	settings, err := n.ConfigurationFor(supported())
	if err != nil {
		t.Fatal(err)
	}
	evidence, err := n.EvidenceFor(settings)
	if err != nil || evidence.Operation != analysis.Kind() || evidence.Profile != n.AnalysisProfile || !reflect.DeepEqual(evidence.Configuration, settings) || !evidence.Semantics {
		t.Fatal(evidence, err)
	}
	contract, err := n.AnalysisContractFor(settings)
	if err != nil || contract.Semantic.Profile != n.AnalysisProfile || contract.Resolution == nil || contract.Resolution.Profile != n.ResolutionProfile || !reflect.DeepEqual(contract.Resolution.Configuration, settings) {
		t.Fatal(contract, err)
	}
	environment := func(value string) []semantic.Setting {
		changed := append([]semantic.Setting(nil), settings...)
		changed[1].Value = []byte(value)
		return changed
	}
	whole := string(settings[1].Value)
	for name, c := range map[string]struct {
		settings []semantic.Setting
		want     error
	}{
		"no settings": {nil, semantic.ErrInvalidRequest},
		"the configuration of the earlier helper profile": {[]semantic.Setting{
			{Key: n.Key("compiler"), Value: []byte("1.90.0")},
			{Key: n.Key("environment"), Value: []byte("edition=2024;target=x86_64-unknown-linux-gnu;crate=captured;lib;sysroot-only;no-codegen")},
			{Key: n.Key("helper"), Value: []byte("captured-types-1")},
		}, semantic.ErrInvalidRequest},
		"the earlier helper profile with an environment": {[]semantic.Setting{settings[0], settings[1], {Key: n.Key("helper"), Value: []byte("captured-types-1")}}, semantic.ErrInvalidRequest},
		"another compiler":         {[]semantic.Setting{{Key: n.Key("compiler"), Value: []byte("1.91.0")}, settings[1], settings[2]}, semantic.ErrInvalidRequest},
		"no environment":           {[]semantic.Setting{settings[0], settings[2]}, semantic.ErrInvalidRequest},
		"two environments":         {[]semantic.Setting{settings[0], settings[1], settings[1], settings[2]}, semantic.ErrInvalidRequest},
		"one more setting":         {append(append([]semantic.Setting(nil), settings...), semantic.Setting{Key: n.Key("extra"), Value: []byte("1")}), semantic.ErrInvalidRequest},
		"another order":            {[]semantic.Setting{settings[2], settings[1], settings[0]}, semantic.ErrInvalidRequest},
		"no JSON":                  {environment(`edition=2024`), semantic.ErrInvalidRequest},
		"an unknown member":        {environment(`{"Edition":"2024","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[],"Features":[]}`), semantic.ErrInvalidRequest},
		"a member in another case": {environment(`{"edition":"2024","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[]}`), semantic.ErrInvalidRequest},
		"a repeated member":        {environment(`{"Edition":"2021","Edition":"2024","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[]}`), semantic.ErrInvalidRequest},
		"another spelling":         {environment(whole + " "), semantic.ErrInvalidRequest},
		"a second value":           {environment(whole + whole), semantic.ErrInvalidRequest},
		"no dependency member":     {environment(`{"Edition":"2024","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"}}`), semantic.ErrInvalidRequest},
		"null":                     {environment(`null`), semantic.ErrInvalidRequest},
		// The form is judged before any value: these state edition 2021 too.
		"an unknown member and another edition":  {environment(`{"Edition":"2021","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[],"Features":[]}`), semantic.ErrInvalidRequest},
		"another spelling and another edition":   {environment(`{"Edition":"2021", "Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[]}`), semantic.ErrInvalidRequest},
		"one more setting and another edition":   {append(environment(`{"Edition":"2021","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[]}`), semantic.Setting{Key: n.Key("extra"), Value: []byte("1")}), semantic.ErrInvalidRequest},
		"the earlier helper and another edition": {[]semantic.Setting{settings[0], environment(`{"Edition":"2021","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[]}`)[1], {Key: n.Key("helper"), Value: []byte("captured-types-1")}}, semantic.ErrInvalidRequest},
		"edition 2021":                           {environment(`{"Edition":"2021","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[]}`), n.ErrUnsupported},
		"another target":                         {environment(`{"Edition":"2024","Target":"wasm32-unknown-unknown","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[]}`), n.ErrUnsupported},
		"one dependency":                         {environment(`{"Edition":"2024","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[{"Name":"serde"}]}`), n.ErrUnsupported},
	} {
		if _, err := n.EnvironmentOf(c.settings); !errors.Is(err, c.want) || errors.Is(err, n.ErrUnsupported) && errors.Is(err, semantic.ErrInvalidRequest) {
			t.Errorf("%s: environment: %v, want %v", name, err, c.want)
		}
		if evidence, err := n.EvidenceFor(c.settings); !errors.Is(err, c.want) || evidence.Profile != "" {
			t.Errorf("%s: evidence: %v, want %v", name, err, c.want)
		}
		if contract, err := n.AnalysisContractFor(c.settings); !errors.Is(err, c.want) || contract.Resolution != nil {
			t.Errorf("%s: analysis contract: %v, want %v", name, err, c.want)
		}
	}
}
