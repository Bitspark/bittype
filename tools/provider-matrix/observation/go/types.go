// Package observation imports invocation-scoped Go test observations. It does not
// interpret provider profiles, case expectations, or language capabilities.
package observation

// Schema is the supported capture-file version.
const Schema = 1

type Setting struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ModuleSelection struct {
	Path     string `json:"path"`
	Version  string `json:"version,omitempty"`
	Revision string `json:"revision,omitempty"`
}

type Replacement struct {
	Path        string `json:"path"`
	Version     string `json:"version,omitempty"`
	Replacement string `json:"replacement"`
	WithVersion string `json:"with_version,omitempty"`
	Revision    string `json:"revision,omitempty"`
}

// SourceContext records the source/build selection, independently of test
// outcomes. Complete is a capture statement, not a substitute for the required
// identities below. Modified, unknown, or unequal selections cannot be joined.
type SourceContext struct {
	Revision       string            `json:"revision"`
	ModulePath     string            `json:"module_path"`
	ModuleRevision string            `json:"module_revision"`
	Harness        string            `json:"harness"`
	GoVersion      string            `json:"go_version"`
	GOOS           string            `json:"goos"`
	GOARCH         string            `json:"goarch"`
	CGOEnabled     string            `json:"cgo_enabled"`
	GOFLAGS        string            `json:"go_flags"`
	BuildFlags     []string          `json:"build_flags"`
	Workspace      string            `json:"workspace"`
	Modules        []ModuleSelection `json:"modules"`
	Replacements   []Replacement     `json:"replacements"`
	InputState     string            `json:"input_state"`
	InputChanges   []string          `json:"input_changes"`
	Complete       bool              `json:"complete"`
	Problems       []string          `json:"problems,omitempty"`
}

// Invocation is written beside stdout.jsonl and stderr.log. ExitCode is absent
// when the capture did not observe the child process completing.
type Invocation struct {
	Schema          int           `json:"schema"`
	Context         SourceContext `json:"source"`
	RunID           string        `json:"run_id"`
	Attempt         string        `json:"attempt"`
	Job             string        `json:"job"`
	Matrix          string        `json:"matrix"`
	CommandID       string        `json:"command_id"`
	Args            []string      `json:"args"`
	Environment     []Setting     `json:"environment"`
	ExitCode        *int          `json:"exit_code"`
	CaptureComplete bool          `json:"capture_complete"`
}

type TestDeclaration struct {
	ID            string
	TestedPackage string
	Name          string
	Variant       string
	File          string
}

type Snapshot struct {
	Context      SourceContext
	Declarations []TestDeclaration
}

type Outcome string

const (
	Passed     Outcome = "pass"
	Failed     Outcome = "fail"
	Skipped    Outcome = "skip"
	Incomplete Outcome = "incomplete"
)

// Occurrence distinguishes repeated executions of an identical emitted name.
// Child events retain their exact names; they do not create authored cases.
type TestRun struct {
	Package    string
	Name       string
	Occurrence int
	Started    bool
	Outcome    Outcome
	Elapsed    float64
}

type PackageRun struct {
	Package     string
	Occurrence  int
	Started     bool
	Outcome     Outcome
	Elapsed     float64
	FailedBuild string
}

type BuildEvent struct {
	ImportPath string
	Action     string
	Output     string
}

type Diagnostic struct {
	Code    string
	Message string
	Event   int
}

type Execution struct {
	Invocation       Invocation
	Tests            []TestRun
	Packages         []PackageRun
	Builds           []BuildEvent
	Stderr           string
	Diagnostics      []Diagnostic
	ProtocolComplete bool
}

type AssociationStatus string

const (
	Matched      AssociationStatus = "matched"
	Mismatched   AssociationStatus = "mismatched"
	Insufficient AssociationStatus = "insufficient"
	Ambiguous    AssociationStatus = "ambiguous"
	Unassociated AssociationStatus = "unassociated"
)

type Association struct {
	Run           TestRun
	Status        AssociationStatus
	Reason        string
	DeclarationID string
	Candidates    []string
}

type Assessment struct {
	Status       AssociationStatus
	Reasons      []string
	Observations []Association
}

// Limits bound imported artifacts. Zero selects the documented defaults.
type Limits struct {
	StdoutBytes int64
	StderrBytes int64
	Events      int
}
