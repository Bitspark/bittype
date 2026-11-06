package report

import (
	"fmt"
	"io"
	"slices"
	"strings"

	observation "github.com/Bitspark/bittype/tools/provider-matrix/observation/go"
)

func observationSnapshot(a *annotations, config Config) observation.Snapshot {
	snapshot := observation.Snapshot{Context: *config.Context}
	for _, test := range a.tests {
		snapshot.Declarations = append(snapshot.Declarations, observation.TestDeclaration{ID: test.id, TestedPackage: test.pkg, Name: test.name, Variant: test.variant, File: test.file})
	}
	return snapshot
}

func renderObservations(b io.Writer, p *program, a *annotations, config Config, sourceURL string) {
	fmt.Fprint(b, "## Invocation observations\n\nEach outcome belongs to one emitted test occurrence in one invocation. Authored profiles and expectations remain separate. Child events do not create independently declared cases.\n\n")
	snapshot := observationSnapshot(a, config)
	assessments := make([]observation.Assessment, len(config.Executions))
	for i, execution := range config.Executions {
		assessment := observation.Associate(snapshot, execution)
		assessments[i] = assessment
		inv := execution.Invocation
		fmt.Fprintf(b, "### Invocation %d: %s\n\n", i+1, escape(inv.CommandID))
		exit := "not observed"
		if inv.ExitCode != nil {
			exit = fmt.Sprint(*inv.ExitCode)
		}
		fmt.Fprintf(b, "Run: %s; attempt: %s; job: %s; matrix: %s.\n\nCommand: %s.\n\nRevision: %s; module: %s (%s); harness: %s.\n\nBuild: %s, %s/%s, CGO=%s, GOFLAGS=%s, build flags=%s.\n\nProcess exit: %s; capture complete: %t; protocol complete: %t; applicability: **%s**.\n\n", escape(inv.RunID), escape(inv.Attempt), escape(inv.Job), escape(inv.Matrix), escape(strings.Join(inv.Args, " ")), escape(inv.Context.Revision), escape(inv.Context.ModulePath), escape(inv.Context.ModuleRevision), escape(inv.Context.Harness), escape(inv.Context.GoVersion), escape(inv.Context.GOOS), escape(inv.Context.GOARCH), escape(inv.Context.CGOEnabled), escape(inv.Context.GOFLAGS), escape(strings.Join(inv.Context.BuildFlags, " ")), exit, inv.CaptureComplete, execution.ProtocolComplete, assessment.Status)
		for _, reason := range assessment.Reasons {
			fmt.Fprintf(b, "- Applicability: %s.\n", escape(reason))
		}
		for _, diagnostic := range execution.Diagnostics {
			fmt.Fprintf(b, "- Capture diagnostic %s: %s.\n", escape(diagnostic.Code), escape(diagnostic.Message))
		}
		for _, event := range execution.Builds {
			if event.Action == "build-fail" {
				fmt.Fprintf(b, "- Build failure: %s (no invented per-test failure).\n", escape(event.ImportPath))
			}
		}
		if len(inv.Context.InputChanges) > 0 {
			fmt.Fprintf(b, "\nInput changes: %s.\n", escape(strings.Join(inv.Context.InputChanges, ", ")))
		}
		fmt.Fprintln(b, "\n| Emitted test name | Package | Occurrence | Outcome | Association |\n| --- | --- | --- | --- | --- |")
		for _, association := range assessment.Observations {
			run := association.Run
			label := string(association.Status)
			if association.Reason != "" {
				label += ": " + association.Reason
			}
			fmt.Fprintf(b, "| %s | %s | %d | %s | %s |\n", escape(run.Name), escape(run.Package), run.Occurrence, run.Outcome, escape(label))
		}
		if len(assessment.Observations) == 0 {
			fmt.Fprintln(b, "| No named test observation supplied | | | | |")
		}
		fmt.Fprint(b, "\nPackage events, stderr and build output remain in the accompanying raw capture; package success never creates a named test observation.\n\n")
	}
	fmt.Fprint(b, "## Associated test declarations\n\n")
	tests := map[string]testDeclaration{}
	for _, c := range a.cases {
		tests[c.test.id] = c.test
	}
	ids := make([]string, 0, len(tests))
	for id := range tests {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	for _, id := range ids {
		test := tests[id]
		pos := p.fset.PositionFor(test.pos, false)
		fmt.Fprintf(b, "<a id=\"test-%s\"></a>\n\n### [%s](%s)\n\nDeclaration: %s / %s / %s.\n\n", id, escape(test.name), sourceTarget(p.root, sourceURL, test.file, pos.Line), escape(test.pkg), escape(test.variant), escape(test.file))
		found := false
		for i, assessment := range assessments {
			for _, association := range assessment.Observations {
				if association.DeclarationID == id && association.Status == observation.Matched {
					found = true
					fmt.Fprintf(b, "- Invocation %d, occurrence %d: **%s**, exact emitted name %s.\n", i+1, association.Run.Occurrence, association.Run.Outcome, escape(association.Run.Name))
				}
			}
		}
		if !found {
			fmt.Fprintln(b, "No matching execution observation supplied for this declaration. See invocation applicability and ambiguous-name diagnostics above.")
		}
		fmt.Fprintln(b)
	}
}
