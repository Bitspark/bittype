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
	// The profile has a new version: the vocabulary lost a member that a
	// reader of the earlier version reads.
	if n.ResolutionProfile != "rustc-1.90.0/captured-types-2/resolution-1" || n.AnalysisProfile != "rustc-1.90.0/captured-types-2/analysis-1" {
		t.Fatal(n.ResolutionProfile, n.AnalysisProfile)
	}
}

// A value that the pinned compiler takes and this profile does not cover is
// unsupported coverage, any other text is invalid input, and no error is
// both. No field is accepted and ignored: every field has a refused value.
// Validity is judged before coverage, over every field.
func TestEnvironmentRefusals(t *testing.T) {
	type change func(*n.Environment)
	edition := func(value string) change { return func(e *n.Environment) { e.Edition = value } }
	target := func(value string) change { return func(e *n.Environment) { e.Target = value } }
	kind := func(value string) change { return func(e *n.Environment) { e.Crate.Kind = value } }
	name := func(value string) change { return func(e *n.Environment) { e.Crate.Name = value } }
	dependencies := func(names ...string) change {
		return func(e *n.Environment) {
			for _, name := range names {
				e.Dependencies = append(e.Dependencies, n.Dependency{Name: name})
			}
		}
	}
	both := func(changes ...change) change {
		return func(e *n.Environment) {
			for _, c := range changes {
				c(e)
			}
		}
	}
	cases := map[string]struct {
		change change
		want   error
	}{
		// An edition: the compiler's four, and no other text.
		"edition 2015":            {edition("2015"), n.ErrUnsupported},
		"edition 2018":            {edition("2018"), n.ErrUnsupported},
		"edition 2021":            {edition("2021"), n.ErrUnsupported},
		"the edition future":      {edition("future"), semantic.ErrInvalidRequest},
		"an edition yet to come":  {edition("2027"), semantic.ErrInvalidRequest},
		"the edition banana":      {edition("banana"), semantic.ErrInvalidRequest},
		"an edition with a space": {edition("2024 "), semantic.ErrInvalidRequest},
		"no edition":              {edition(""), semantic.ErrInvalidRequest},
		// A crate kind: the compiler's seven, and no other text.
		"an rlib":                  {kind("rlib"), n.ErrUnsupported},
		"a static library":         {kind("staticlib"), n.ErrUnsupported},
		"a dynamic library":        {kind("dylib"), n.ErrUnsupported},
		"a C dynamic library":      {kind("cdylib"), n.ErrUnsupported},
		"a binary crate":           {kind("bin"), n.ErrUnsupported},
		"a procedural macro":       {kind("proc-macro"), n.ErrUnsupported},
		"the unstable kind sdylib": {kind("sdylib"), semantic.ErrInvalidRequest},
		"the kind banana":          {kind("banana"), semantic.ErrInvalidRequest},
		"the kind Lib":             {kind("Lib"), semantic.ErrInvalidRequest},
		"two kinds":                {kind("lib,bin"), semantic.ErrInvalidRequest},
		"no crate kind":            {kind(""), semantic.ErrInvalidRequest},
		// A target: a well-formed text is unsupported, whether or not a
		// compiler has such a target.
		"another target":                          {target("aarch64-unknown-linux-gnu"), n.ErrUnsupported},
		"a target with a dot":                     {target("thumbv8m.main-none-eabi"), n.ErrUnsupported},
		"a well-formed text that names no target": {target("banana"), n.ErrUnsupported},
		"a target in upper case":                  {target("X86_64-unknown-linux-gnu"), n.ErrUnsupported},
		"a target with a space":                   {target("not a target"), semantic.ErrInvalidRequest},
		"a target file":                           {target("x86_64-unknown-linux-gnu.json"), semantic.ErrInvalidRequest},
		"a target path":                           {target("targets/custom"), semantic.ErrInvalidRequest},
		"a target outside ASCII":                  {target("ziel-ü"), semantic.ErrInvalidRequest},
		"no target":                               {target(""), semantic.ErrInvalidRequest},
		// A crate name: what the compiler takes, and of that an ASCII
		// identifier that is no keyword and no crate of the sysroot.
		"a crate name of digits":                  {name("1abc"), n.ErrUnsupported},
		"an underscore":                           {name("_"), n.ErrUnsupported},
		"a crate name of letters outside ASCII":   {name("über"), n.ErrUnsupported},
		"a crate name with a superscript digit":   {name("x²"), n.ErrUnsupported},
		"a crate name of a roman numeral":         {name("ⅷ"), n.ErrUnsupported},
		"a crate name with a digit outside ASCII": {name("a٣"), n.ErrUnsupported},
		"no crate name":                           {name(""), semantic.ErrInvalidRequest},
		"a hyphen in the crate name":              {name("my-pages"), semantic.ErrInvalidRequest},
		"an option as crate name":                 {name("--cfg"), semantic.ErrInvalidRequest},
		"a space in the crate name":               {name("a b"), semantic.ErrInvalidRequest},
		"a raw identifier as crate name":          {name("r#type"), semantic.ErrInvalidRequest},
		"a combining mark in the crate name":      {name("á"), semantic.ErrInvalidRequest},
		"a line break in the crate name":          {name("pages\n"), semantic.ErrInvalidRequest},
		// A dependency: its name is judged as a crate name, and then no
		// dependency is covered.
		"one dependency":                               {dependencies("serde"), n.ErrUnsupported},
		"a dependency named like a keyword":            {dependencies("crate"), n.ErrUnsupported},
		"a dependency with a name outside ASCII":       {dependencies("x²"), n.ErrUnsupported},
		"a dependency without a name":                  {dependencies(""), semantic.ErrInvalidRequest},
		"a dependency with a name that is none":        {dependencies("a b!"), semantic.ErrInvalidRequest},
		"a second dependency with a name that is none": {dependencies("serde", "my-crate"), semantic.ErrInvalidRequest},
		// Validity is judged before coverage, over every field.
		"edition 2021 and no target":                              {both(edition("2021"), target("")), semantic.ErrInvalidRequest},
		"edition 2021 and the kind banana":                        {both(edition("2021"), kind("banana")), semantic.ErrInvalidRequest},
		"a binary crate and no crate name":                        {func(e *n.Environment) { e.Crate = n.Crate{Kind: "bin"} }, semantic.ErrInvalidRequest},
		"another target and a hyphenated name":                    {both(target("wasm32-unknown-unknown"), name("my-pages")), semantic.ErrInvalidRequest},
		"one dependency and no edition":                           {both(edition(""), dependencies("serde")), semantic.ErrInvalidRequest},
		"edition 2021 and a dependency without a name":            {both(edition("2021"), dependencies("")), semantic.ErrInvalidRequest},
		"a keyword as crate name and a dependency without a name": {both(name("crate"), dependencies("")), semantic.ErrInvalidRequest},
		"edition 2021 and one dependency":                         {both(edition("2021"), dependencies("serde")), n.ErrUnsupported},
		"the zero environment":                                    {func(e *n.Environment) { *e = n.Environment{} }, semantic.ErrInvalidRequest},
	}
	// Every strict and every reserved keyword of edition 2024 is refused.
	for _, word := range []string{
		"as", "async", "await", "break", "const", "continue", "crate", "dyn", "else", "enum", "extern", "false", "fn", "for", "if", "impl", "in", "let", "loop", "match", "mod", "move", "mut", "pub", "ref", "return", "self", "Self", "static", "struct", "super", "trait", "true", "type", "unsafe", "use", "where", "while",
		"abstract", "become", "box", "do", "final", "gen", "macro", "override", "priv", "try", "typeof", "unsized", "virtual", "yield",
	} {
		cases["the keyword "+word+" as crate name"] = struct {
			change change
			want   error
		}{name(word), n.ErrUnsupported}
	}
	// Every crate that the pinned sysroot loads for a captured crate is refused.
	for _, crate := range []string{
		"std", "core", "alloc", "addr2line", "adler2", "cfg_if", "compiler_builtins", "gimli", "hashbrown", "libc", "memchr", "miniz_oxide", "object", "panic_unwind", "rustc_demangle", "rustc_std_workspace_alloc", "rustc_std_workspace_core", "std_detect", "unwind",
	} {
		cases["the sysroot crate "+crate+" as crate name"] = struct {
			change change
			want   error
		}{name(crate), n.ErrUnsupported}
	}
	for name, c := range cases {
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
	// Admitted: an ASCII identifier, also one that starts like a sysroot
	// crate or is a weak keyword, and a sysroot crate that is not loaded.
	for _, name := range []string{"pages", "Pages", "_pages", "__", "pages_2", "corelib", "stdx", "core_", "alloc2", "test", "proc_macro", "union", "macro_rules", "raw", "safe", "auto", "default"} {
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
	// The requirement asks for the three reports, each complete and holding.
	if len(evidence.Reports) != 3 || !reflect.DeepEqual(evidence.Reports, n.Requirements()) {
		t.Fatal("the evidence requirement does not require the reports", evidence.Reports)
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
		// A value that is no value of the compiler is invalid, also in settings.
		"the edition banana":                 {environment(`{"Edition":"banana","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[]}`), semantic.ErrInvalidRequest},
		"a dependency without a name":        {environment(`{"Edition":"2024","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[{"Name":""}]}`), semantic.ErrInvalidRequest},
		"a crate named like a sysroot crate": {environment(`{"Edition":"2024","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"core","Kind":"lib"},"Dependencies":[]}`), n.ErrUnsupported},
		"edition 2021":                       {environment(`{"Edition":"2021","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[]}`), n.ErrUnsupported},
		"another target":                     {environment(`{"Edition":"2024","Target":"wasm32-unknown-unknown","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[]}`), n.ErrUnsupported},
		"one dependency":                     {environment(`{"Edition":"2024","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[{"Name":"serde"}]}`), n.ErrUnsupported},
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
