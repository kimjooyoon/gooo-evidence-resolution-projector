package projector

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kimjooyoon/gooo-evidence-resolution-projector/internal/gooo"
	"github.com/kimjooyoon/gooo-evidence-resolution-projector/internal/model"
)

const (
	OutputManifest  = "projection-manifest.json"
	OutputGraph     = "evidence-graph.json"
	OutputEvents    = "projection-events.ndjson"
	OutputUser      = "user-view.md"
	OutputOperator  = "operator-view.md"
	OutputAuditor   = "auditor-view.md"
	OutputReceipt   = "projection-receipt.json"
)

var OutputNames = []string{
	OutputManifest,
	OutputGraph,
	OutputEvents,
	OutputUser,
	OutputOperator,
	OutputAuditor,
	OutputReceipt,
}

type GraphOutput struct {
	GraphID                string              `json:"graph_id"`
	Release                string              `json:"release"`
	ProjectionSchemaVersion string             `json:"projection_schema_version"`
	DecisionOrder           []string            `json:"decision_order"`
	DecisionDigest          string              `json:"decision_digest"`
	Activities              []model.Activity    `json:"activities"`
	Evidence                []model.Evidence    `json:"evidence"`
	Cases                   []CaseOutput        `json:"cases"`
}

type CaseOutput struct {
	ID                        string                    `json:"id"`
	Title                     string                    `json:"title"`
	Summary                   string                    `json:"summary"`
	SummaryRefs               []string                  `json:"summary_refs"`
	ExpectedDecision          string                    `json:"expected_decision"`
	ActivityID                string                    `json:"activity_id"`
	Decision                  string                    `json:"decision"`
	DecisionRank              int                       `json:"decision_rank"`
	DecisionDigest            string                    `json:"decision_digest"`
	Claims                    []model.Claim             `json:"claims"`
	ClaimStates               []ClaimState               `json:"claim_states"`
	ProofChoices              []ProofChoice              `json:"proof_choices"`
	EvidenceIDs               []string                  `json:"evidence_ids"`
	EvidenceDescriptions      []EvidenceDescription     `json:"evidence_descriptions"`
	UnknownFrontier           *model.UnknownFrontier     `json:"unknown_frontier"`
	Counterexamples           []model.Counterexample     `json:"counterexamples"`
	CounterexampleIDs         []string                  `json:"counterexample_ids"`
	CounterexampleDescriptions []CounterexampleDescription `json:"counterexample_descriptions"`
	SourceContext             string                    `json:"source_context"`
	OperatorNotes             string                    `json:"operator_notes"`
	AuditTrace                string                    `json:"audit_trace"`
}

type ClaimState struct {
	ClaimID string `json:"claim_id"`
	State   string `json:"state"`
}

type ProofChoice struct {
	ClaimID string `json:"claim_id"`
	Choice  string `json:"choice"`
}

type EvidenceDescription struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type CounterexampleDescription struct {
	ID          string   `json:"id"`
	ClaimID     string   `json:"claim_id"`
	Description string   `json:"description"`
	EvidenceIDs []string `json:"evidence_ids"`
}

type CaseProjection struct {
	CaseID            string         `json:"case_id"`
	Decision          string         `json:"decision"`
	DecisionDigest    string         `json:"decision_digest"`
	Fields            map[string]any `json:"fields"`
	OmittedFieldIDs   []string       `json:"omitted_field_ids"`
	OmittedFieldCount int            `json:"omitted_field_count"`
	FieldOrder        []string       `json:"-"`
}

type ViewProjection struct {
	ID                string           `json:"id"`
	Label             string           `json:"label"`
	OmittedFieldIDs   []string         `json:"omitted_field_ids"`
	OmittedFieldCount int              `json:"omitted_field_count"`
	Cases             []CaseProjection `json:"cases"`
}

type ProjectionEvent struct {
	Sequence          int            `json:"sequence"`
	EventType         string         `json:"event_type"`
	GraphID           string         `json:"graph_id"`
	View              string         `json:"view"`
	CaseID            string         `json:"case_id"`
	Decision          string         `json:"decision"`
	DecisionDigest    string         `json:"decision_digest"`
	Fields            map[string]any `json:"fields"`
	OmittedFieldIDs   []string       `json:"omitted_field_ids"`
	OmittedFieldCount int            `json:"omitted_field_count"`
}

type Inventory struct {
	GoooActivities            int            `json:"gooo_activities"`
	CanonicalCases            int            `json:"canonical_cases"`
	CasesByDecision           map[string]int `json:"cases_by_decision"`
	ReaderResolutions         int            `json:"reader_resolutions"`
	CallerOwnedOutputKinds    int            `json:"caller_owned_output_kinds"`
	ProjectionEvents          int            `json:"projection_events"`
	InputRepoMutations        int            `json:"input_repo_mutations"`
	SourceMutations           int            `json:"source_mutations"`
	RuntimeSideEffects        int            `json:"runtime_side_effects"`
}

type ViewManifest struct {
	ID                string   `json:"id"`
	Label             string   `json:"label"`
	OmittedFieldIDs   []string `json:"omitted_field_ids"`
	OmittedFieldCount int      `json:"omitted_field_count"`
}

type Manifest struct {
	ManifestVersion       string         `json:"manifest_version"`
	GraphID               string         `json:"graph_id"`
	Release               string         `json:"release"`
	GoooSource            string         `json:"gooo_source"`
	GoooSourceSHA256      string         `json:"gooo_source_sha256"`
	CanonicalGraphSHA256  string         `json:"canonical_graph_sha256"`
	DecisionDigest        string         `json:"decision_digest"`
	DecisionOrder         []string       `json:"decision_order"`
	Inventory             Inventory      `json:"inventory"`
	CallerOwnedOutputKinds []string      `json:"caller_owned_output_kinds"`
	Views                 []ViewManifest `json:"views"`
}

type ReceiptFile struct {
	Name      string `json:"name"`
	SizeBytes int    `json:"size_bytes"`
	SHA256    string `json:"sha256"`
}

type Receipt struct {
	ReceiptVersion        string       `json:"receipt_version"`
	Status                string       `json:"status"`
	GraphID               string       `json:"graph_id"`
	Release               string       `json:"release"`
	DecisionDigest        string       `json:"decision_digest"`
	Inventory             Inventory    `json:"inventory"`
	CallerOwnedOutputKinds []string    `json:"caller_owned_output_kinds"`
	Files                 []ReceiptFile `json:"content_files"`
	SelfDigestPolicy      string       `json:"self_digest_policy"`
}

type Result struct {
	Source         model.SourceDocument
	SourceBytes    []byte
	Graph          GraphOutput
	Views          []ViewProjection
	Events         []ProjectionEvent
	Inventory      Inventory
	Manifest       Manifest
	Receipt        Receipt
	OutputContents map[string][]byte
}

func LoadAndProject(sourcePath string) (Result, error) {
	document, sourceBytes, err := gooo.Load(sourcePath)
	if err != nil {
		return Result{}, err
	}
	result, err := Project(document, sourceBytes)
	if err != nil {
		return Result{}, err
	}
	return result, nil
}

func Project(document model.SourceDocument, sourceBytes []byte) (Result, error) {
	if err := validateDocument(document); err != nil {
		return Result{}, err
	}

	evidenceByID := make(map[string]model.Evidence, len(document.Graph.Evidence))
	for _, evidence := range document.Graph.Evidence {
		evidenceByID[evidence.ID] = evidence
	}
	activityByCase := make(map[string]model.Activity, len(document.Graph.Activities))
	for _, activity := range document.Graph.Activities {
		activityByCase[activity.Binding.CaseID] = activity
	}

	cases := make([]CaseOutput, 0, len(document.Graph.Cases))
	decisionRanks := rankByDecision(document.Graph.Projection.DecisionOrder)
	for _, sourceCase := range document.Graph.Cases {
		caseOutput, err := evaluateCase(document.Graph, sourceCase, evidenceByID, activityByCase, decisionRanks)
		if err != nil {
			return Result{}, err
		}
		cases = append(cases, caseOutput)
	}
	counts := map[string]int{}
	for _, item := range cases {
		counts[item.Decision]++
	}
	semantic := document.Graph.Projection.SemanticRules
	if counts[semantic.ClosedState] != 3 || counts[semantic.UnknownState] != 3 || counts[semantic.RefutedState] != 3 {
		return Result{}, fmt.Errorf("released graph must contain three cases for each semantic decision: %#v", counts)
	}

	decisionDigest, err := computeDecisionDigest(cases)
	if err != nil {
		return Result{}, err
	}
	for index := range cases {
		cases[index].DecisionDigest = decisionDigest
	}

	graph := GraphOutput{
		GraphID:                 document.Graph.ID,
		Release:                 document.Graph.Release,
		ProjectionSchemaVersion: document.Graph.Projection.SchemaVersion,
		DecisionOrder:            append([]string(nil), document.Graph.Projection.DecisionOrder...),
		DecisionDigest:           decisionDigest,
		Activities:               append([]model.Activity(nil), document.Graph.Activities...),
		Evidence:                 append([]model.Evidence(nil), document.Graph.Evidence...),
		Cases:                    cases,
	}

	views, events, err := buildViews(document.Graph, graph)
	if err != nil {
		return Result{}, err
	}

	inventory := Inventory{
		GoooActivities:         len(document.Graph.Activities),
		CanonicalCases:         len(cases),
		CasesByDecision:        counts,
		ReaderResolutions:      len(views),
		CallerOwnedOutputKinds: len(OutputNames),
		ProjectionEvents:       len(events),
		InputRepoMutations:     0,
		SourceMutations:        0,
		RuntimeSideEffects:     0,
	}

	graphBytes, err := jsonBytes(graph)
	if err != nil {
		return Result{}, fmt.Errorf("encode evidence graph: %w", err)
	}
	eventsBytes, err := encodeEvents(events)
	if err != nil {
		return Result{}, err
	}
	viewBytes := make(map[string][]byte, len(views))
	for _, view := range views {
		data, err := renderView(view, graph)
		if err != nil {
			return Result{}, err
		}
		viewBytes[view.ID+"-view.md"] = data
	}

	manifest := Manifest{
		ManifestVersion:        "1.0.0",
		GraphID:                document.Graph.ID,
		Release:                document.Graph.Release,
		GoooSource:             ".gooo/released.gooo",
		GoooSourceSHA256:       digest(sourceBytes),
		CanonicalGraphSHA256:   digest(graphBytes),
		DecisionDigest:         decisionDigest,
		DecisionOrder:          append([]string(nil), document.Graph.Projection.DecisionOrder...),
		Inventory:              inventory,
		CallerOwnedOutputKinds: append([]string(nil), OutputNames...),
		Views:                  viewManifests(views),
	}
	manifestBytes, err := jsonBytes(manifest)
	if err != nil {
		return Result{}, fmt.Errorf("encode projection manifest: %w", err)
	}

	contents := map[string][]byte{
		OutputManifest: manifestBytes,
		OutputGraph:    graphBytes,
		OutputEvents:   eventsBytes,
		OutputUser:     viewBytes["user-view.md"],
		OutputOperator: viewBytes["operator-view.md"],
		OutputAuditor:  viewBytes["auditor-view.md"],
	}
	receipt := makeReceipt(document.Graph, inventory, decisionDigest, contents)
	receiptBytes, err := jsonBytes(receipt)
	if err != nil {
		return Result{}, fmt.Errorf("encode projection receipt: %w", err)
	}
	contents[OutputReceipt] = receiptBytes

	return Result{
		Source:         document,
		SourceBytes:    append([]byte(nil), sourceBytes...),
		Graph:          graph,
		Views:          views,
		Events:         events,
		Inventory:      inventory,
		Manifest:       manifest,
		Receipt:        receipt,
		OutputContents: contents,
	}, nil
}

func WriteOutputs(result Result, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	expected := make(map[string]bool, len(OutputNames))
	for _, name := range OutputNames {
		expected[name] = true
	}
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return fmt.Errorf("read output directory: %w", err)
	}
	for _, entry := range entries {
		if !expected[entry.Name()] {
			return fmt.Errorf("unexpected caller-owned output: %s", entry.Name())
		}
	}
	for _, name := range OutputNames {
		data, ok := result.OutputContents[name]
		if !ok {
			return fmt.Errorf("missing generated output: %s", name)
		}
		if err := os.WriteFile(filepath.Join(outputDir, name), data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}
	return nil
}

func validateDocument(document model.SourceDocument) error {
	graph := document.Graph
	if document.Language == "" || document.LanguageVersion == "" {
		return fmt.Errorf("Gooo language metadata is required")
	}
	if graph.ID == "" || graph.Release == "" {
		return fmt.Errorf("graph id and release are required")
	}
	projection := graph.Projection
	if projection.SchemaVersion == "" || len(projection.DecisionOrder) == 0 {
		return fmt.Errorf("projection schema and decision order are required")
	}
	seenDecisions := map[string]bool{}
	for _, state := range projection.DecisionOrder {
		if state == "" || seenDecisions[state] {
			return fmt.Errorf("decision order must contain unique non-empty states")
		}
		seenDecisions[state] = true
	}
	for _, state := range []string{
		projection.SemanticRules.RefutedState,
		projection.SemanticRules.UnknownState,
		projection.SemanticRules.ClosedState,
		projection.SemanticRules.ExternalUtilityZeroEvidenceState,
	} {
		if !seenDecisions[state] {
			return fmt.Errorf("semantic state %q is absent from decision order", state)
		}
	}

	viewIDs := map[string]bool{}
	for _, view := range projection.Views {
		if view.ID == "" || viewIDs[view.ID] {
			return fmt.Errorf("reader view IDs must be unique and non-empty")
		}
		viewIDs[view.ID] = true
	}
	if len(projection.Views) != 3 {
		return fmt.Errorf("exactly three reader resolutions are required")
	}
	fieldIDs := map[string]bool{}
	fieldByID := map[string]model.ProjectionField{}
	for _, field := range projection.Fields {
		if field.ID == "" || field.Source == "" || fieldIDs[field.ID] {
			return fmt.Errorf("projection field IDs and sources must be unique and non-empty")
		}
		fieldIDs[field.ID] = true
		fieldByID[field.ID] = field
		if field.Invariant && len(field.AllowedOmission) != 0 {
			return fmt.Errorf("invariant field %q has an allowed omission", field.ID)
		}
		for _, viewID := range field.AllowedOmission {
			if !viewIDs[viewID] {
				return fmt.Errorf("field %q omits unknown view %q", field.ID, viewID)
			}
		}
	}
	for _, view := range projection.Views {
		omits := map[string]bool{}
		for _, fieldID := range view.Omits {
			if !fieldIDs[fieldID] || omits[fieldID] {
				return fmt.Errorf("view %q has an invalid omission", view.ID)
			}
			omits[fieldID] = true
		}
		for fieldID, field := range fieldByID {
			declared := false
			for _, allowedView := range field.AllowedOmission {
				if allowedView == view.ID {
					declared = true
				}
			}
			if declared != omits[fieldID] {
				return fmt.Errorf("field %q and view %q disagree about omission", fieldID, view.ID)
			}
		}
	}
	for _, group := range projection.Groups {
		if group.ID == "" || group.ExactCount != len(group.FieldIDs) || group.ExactCount == 0 {
			return fmt.Errorf("projection group %q does not declare an exact field count", group.ID)
		}
		seen := map[string]bool{}
		for _, fieldID := range group.FieldIDs {
			if !fieldIDs[fieldID] || seen[fieldID] {
				return fmt.Errorf("projection group %q has an invalid field", group.ID)
			}
			seen[fieldID] = true
		}
	}
	if len(projection.Groups) == 0 {
		return fmt.Errorf("at least one exact projection group is required")
	}

	if len(graph.Activities) != 9 || len(graph.Cases) != 9 {
		return fmt.Errorf("released graph must contain exactly nine activities and nine cases")
	}
	activityIDs := map[string]bool{}
	activityCaseIDs := map[string]bool{}
	activityIDByCase := map[string]string{}
	for _, activity := range graph.Activities {
		if activity.ID == "" || activity.Name == "" || activityCaseIDs[activity.Binding.CaseID] || activityIDs[activity.ID] {
			return fmt.Errorf("activities must be unique and bind one-to-one to cases")
		}
		activityIDs[activity.ID] = true
		activityCaseIDs[activity.Binding.CaseID] = true
		activityIDByCase[activity.Binding.CaseID] = activity.ID
	}
	evidenceIDs := map[string]bool{}
	for _, evidence := range graph.Evidence {
		if evidence.ID == "" || evidence.Description == "" || evidence.SourceRef == "" || evidenceIDs[evidence.ID] {
			return fmt.Errorf("evidence IDs and content must be unique and non-empty")
		}
		evidenceIDs[evidence.ID] = true
	}
	caseIDs := map[string]bool{}
	claimIDs := map[string]bool{}
	counterexampleIDs := map[string]bool{}
	for _, item := range graph.Cases {
		if item.ID == "" || item.Title == "" || item.Summary == "" || item.ExpectedDecision == "" || item.ActivityID == "" || caseIDs[item.ID] {
			return fmt.Errorf("cases must have unique identity and required content")
		}
		caseIDs[item.ID] = true
		if !activityIDs[item.ActivityID] || !activityCaseIDs[item.ID] || activityIDByCase[item.ID] != item.ActivityID {
			return fmt.Errorf("case %q is not bound to exactly one activity", item.ID)
		}
		if !strings.Contains(item.Summary, item.ID) || len(item.SummaryRefs) == 0 {
			return fmt.Errorf("case %q summary must contain stable references", item.ID)
		}
		if len(item.Claims) == 0 {
			return fmt.Errorf("case %q has no claims", item.ID)
		}
		for _, claim := range item.Claims {
			if claim.ID == "" || claim.State == "" || claim.ProofChoice == "" || claimIDs[claim.ID] {
				return fmt.Errorf("case %q contains an invalid claim", item.ID)
			}
			claimIDs[claim.ID] = true
			for _, evidenceID := range claim.EvidenceIDs {
				if !evidenceIDs[evidenceID] {
					return fmt.Errorf("claim %q references unknown evidence %q", claim.ID, evidenceID)
				}
			}
			if claim.Utility != nil {
				for _, evidenceID := range claim.Utility.EvidenceIDs {
					if !evidenceIDs[evidenceID] {
						return fmt.Errorf("utility on claim %q references unknown evidence %q", claim.ID, evidenceID)
					}
				}
			}
		}
		for _, counterexample := range item.Counterexamples {
			if counterexample.ID == "" || counterexample.ClaimID == "" || counterexample.Description == "" || counterexampleIDs[counterexample.ID] {
				return fmt.Errorf("case %q contains an invalid counterexample", item.ID)
			}
			counterexampleIDs[counterexample.ID] = true
			if !claimIDs[counterexample.ClaimID] {
				return fmt.Errorf("counterexample %q references unknown claim", counterexample.ID)
			}
			for _, evidenceID := range counterexample.EvidenceIDs {
				if !evidenceIDs[evidenceID] {
					return fmt.Errorf("counterexample %q references unknown evidence", counterexample.ID)
				}
			}
		}
	}
	for _, activity := range graph.Activities {
		if !caseIDs[activity.Binding.CaseID] {
			return fmt.Errorf("activity %q references unknown case", activity.ID)
		}
	}
	knownIDs := map[string]bool{}
	for id := range caseIDs {
		knownIDs[id] = true
	}
	for id := range claimIDs {
		knownIDs[id] = true
	}
	for id := range evidenceIDs {
		knownIDs[id] = true
	}
	for id := range counterexampleIDs {
		knownIDs[id] = true
	}
	for _, item := range graph.Cases {
		for _, ref := range item.SummaryRefs {
			if !knownIDs[ref] || !strings.Contains(item.Summary, ref) {
				return fmt.Errorf("case %q has an invalid stable summary reference %q", item.ID, ref)
			}
		}
	}
	return nil
}

func evaluateCase(graph model.Graph, sourceCase model.Case, evidenceByID map[string]model.Evidence, activityByCase map[string]model.Activity, ranks map[string]int) (CaseOutput, error) {
	if activityByCase[sourceCase.ID].ID == "" || activityByCase[sourceCase.ID].ID == "" {
		return CaseOutput{}, fmt.Errorf("case %q has no activity binding", sourceCase.ID)
	}
	bestState := ""
	bestRank := len(ranks) + 1
	claimStates := make([]ClaimState, 0, len(sourceCase.Claims))
	proofChoices := make([]ProofChoice, 0, len(sourceCase.Claims))
	claims := make([]model.Claim, 0, len(sourceCase.Claims))
	evidenceIDs := make([]string, 0)
	seenEvidence := map[string]bool{}
	for _, sourceClaim := range sourceCase.Claims {
		state := sourceClaim.State
		if sourceClaim.Utility != nil && sourceClaim.Utility.Kind == "external" && len(sourceClaim.Utility.EvidenceIDs) == 0 {
			state = graph.Projection.SemanticRules.ExternalUtilityZeroEvidenceState
		}
		rank, ok := ranks[state]
		if !ok {
			return CaseOutput{}, fmt.Errorf("claim %q uses a state outside Gooo decision order", sourceClaim.ID)
		}
		if rank < bestRank {
			bestRank = rank
			bestState = state
		}
		claimStates = append(claimStates, ClaimState{ClaimID: sourceClaim.ID, State: state})
		proofChoices = append(proofChoices, ProofChoice{ClaimID: sourceClaim.ID, Choice: sourceClaim.ProofChoice})
		claim := sourceClaim
		claim.State = state
		claims = append(claims, claim)
		for _, evidenceID := range sourceClaim.EvidenceIDs {
			if !seenEvidence[evidenceID] {
				evidenceIDs = append(evidenceIDs, evidenceID)
				seenEvidence[evidenceID] = true
			}
		}
	}
	if bestState == "" {
		return CaseOutput{}, fmt.Errorf("case %q did not produce a decision", sourceCase.ID)
	}
	hasUnknown := false
	hasRefuted := false
	for _, claimState := range claimStates {
		if claimState.State == graph.Projection.SemanticRules.UnknownState {
			hasUnknown = true
		}
		if claimState.State == graph.Projection.SemanticRules.RefutedState {
			hasRefuted = true
		}
	}
	if hasUnknown && sourceCase.UnknownFrontier == nil {
		return CaseOutput{}, fmt.Errorf("UNKNOWN case %q is missing its causal frontier", sourceCase.ID)
	}
	if hasUnknown && !completeUnknownFrontier(sourceCase.UnknownFrontier) {
		return CaseOutput{}, fmt.Errorf("UNKNOWN case %q does not contain six complete frontier fields", sourceCase.ID)
	}
	if hasRefuted && len(sourceCase.Counterexamples) == 0 {
		return CaseOutput{}, fmt.Errorf("REFUTED case %q is missing a counterexample", sourceCase.ID)
	}
	if sourceCase.ExpectedDecision != bestState {
		return CaseOutput{}, fmt.Errorf("case %q expected %q but evaluated %q", sourceCase.ID, sourceCase.ExpectedDecision, bestState)
	}
	for _, ref := range sourceCase.SummaryRefs {
		if !strings.Contains(sourceCase.Summary, ref) {
			return CaseOutput{}, fmt.Errorf("case %q summary does not reference %q", sourceCase.ID, ref)
		}
	}
	evidenceDescriptions := make([]EvidenceDescription, 0, len(evidenceIDs))
	for _, evidenceID := range evidenceIDs {
		evidenceDescriptions = append(evidenceDescriptions, EvidenceDescription{ID: evidenceID, Description: evidenceByID[evidenceID].Description})
	}
	counterexampleIDs := make([]string, 0, len(sourceCase.Counterexamples))
	counterexampleDescriptions := make([]CounterexampleDescription, 0, len(sourceCase.Counterexamples))
	for _, counterexample := range sourceCase.Counterexamples {
		counterexampleIDs = append(counterexampleIDs, counterexample.ID)
		counterexampleDescriptions = append(counterexampleDescriptions, CounterexampleDescription{
			ID:          counterexample.ID,
			ClaimID:     counterexample.ClaimID,
			Description: counterexample.Description,
			EvidenceIDs: append([]string(nil), counterexample.EvidenceIDs...),
		})
	}
	return CaseOutput{
		ID:                         sourceCase.ID,
		Title:                      sourceCase.Title,
		Summary:                    sourceCase.Summary,
		SummaryRefs:                append([]string(nil), sourceCase.SummaryRefs...),
		ExpectedDecision:           sourceCase.ExpectedDecision,
		ActivityID:                 sourceCase.ActivityID,
		Decision:                   bestState,
		DecisionRank:               bestRank,
		Claims:                     claims,
		ClaimStates:                claimStates,
		ProofChoices:               proofChoices,
		EvidenceIDs:                evidenceIDs,
		EvidenceDescriptions:       evidenceDescriptions,
		UnknownFrontier:            sourceCase.UnknownFrontier,
		Counterexamples:            append([]model.Counterexample(nil), sourceCase.Counterexamples...),
		CounterexampleIDs:          counterexampleIDs,
		CounterexampleDescriptions: counterexampleDescriptions,
		SourceContext:              sourceCase.SourceContext,
		OperatorNotes:              sourceCase.OperatorNotes,
		AuditTrace:                 sourceCase.AuditTrace,
	}, nil
}

func buildViews(graph model.Graph, output GraphOutput) ([]ViewProjection, []ProjectionEvent, error) {
	views := make([]ViewProjection, 0, len(graph.Projection.Views))
	events := make([]ProjectionEvent, 0, len(graph.Projection.Views)*len(output.Cases))
	sequence := 1
	for _, view := range graph.Projection.Views {
		omitted := omittedFields(graph.Projection.Fields, view.ID)
		projection := ViewProjection{
			ID:                view.ID,
			Label:             view.Label,
			OmittedFieldIDs:   omitted,
			OmittedFieldCount: len(omitted),
			Cases:             make([]CaseProjection, 0, len(output.Cases)),
		}
		for _, item := range output.Cases {
			fields, err := projectFields(graph.Projection.Fields, view.ID, item)
			if err != nil {
				return nil, nil, err
			}
			caseProjection := CaseProjection{
				CaseID:            item.ID,
				Decision:          item.Decision,
				DecisionDigest:    output.DecisionDigest,
				Fields:            fields,
				OmittedFieldIDs:   omitted,
				OmittedFieldCount: len(omitted),
				FieldOrder:        includedFieldIDs(graph.Projection.Fields, view.ID),
			}
			projection.Cases = append(projection.Cases, caseProjection)
			events = append(events, ProjectionEvent{
				Sequence:          sequence,
				EventType:         "projection",
				GraphID:           output.GraphID,
				View:              view.ID,
				CaseID:            item.ID,
				Decision:          item.Decision,
				DecisionDigest:    output.DecisionDigest,
				Fields:            fields,
				OmittedFieldIDs:   omitted,
				OmittedFieldCount: len(omitted),
			})
			sequence++
		}
		views = append(views, projection)
	}
	return views, events, nil
}

func projectFields(fields []model.ProjectionField, viewID string, item CaseOutput) (map[string]any, error) {
	data, err := objectMap(item)
	if err != nil {
		return nil, fmt.Errorf("normalize case %q: %w", item.ID, err)
	}
	projected := make(map[string]any)
	for _, field := range fields {
		if contains(field.AllowedOmission, viewID) {
			continue
		}
		value, found := pathValue(data, field.Source)
		if !found {
			if field.Invariant {
				return nil, fmt.Errorf("invariant field %q is absent for case %q", field.ID, item.ID)
			}
			value = nil
		}
		projected[field.ID] = value
	}
	return projected, nil
}

func objectMap(item CaseOutput) (map[string]any, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	var object map[string]any
	if err := json.Unmarshal(data, &object); err != nil {
		return nil, err
	}
	return object, nil
}

func pathValue(object map[string]any, path string) (any, bool) {
	var current any = object
	for _, part := range strings.Split(path, ".") {
		mapping, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		value, ok := mapping[part]
		if !ok {
			return nil, false
		}
		if value == nil {
			return nil, true
		}
		current = value
	}
	return current, true
}

func renderView(view ViewProjection, graph GraphOutput) ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString("# ")
	buffer.WriteString(view.Label)
	buffer.WriteString(" projection\n\n")
	fmt.Fprintf(&buffer, "- graph ID: `%s`\n", graph.GraphID)
	fmt.Fprintf(&buffer, "- decision digest: `%s`\n", graph.DecisionDigest)
	fmt.Fprintf(&buffer, "- omitted field IDs: %s\n", formatIDs(view.OmittedFieldIDs))
	fmt.Fprintf(&buffer, "- omitted field count: `%d`\n", view.OmittedFieldCount)
	buffer.WriteString("- projection contract: omitted context is loss-declared; invariant evidence and decisions remain reverse-referenceable.\n\n")
		for _, item := range view.Cases {
		fmt.Fprintf(&buffer, "## %s\n\n", item.CaseID)
		fmt.Fprintf(&buffer, "- decision: `%s`\n", item.Decision)
		fmt.Fprintf(&buffer, "- decision digest: `%s`\n", item.DecisionDigest)
		for _, field := range item.FieldOrder {
			fmt.Fprintf(&buffer, "- %s: %s\n", field, formatValue(item.Fields[field]))
		}
		buffer.WriteString("\n")
	}
	return buffer.Bytes(), nil
}

func includedFieldIDs(fields []model.ProjectionField, viewID string) []string {
	order := make([]string, 0, len(fields))
	for _, field := range fields {
		if !contains(field.AllowedOmission, viewID) {
			order = append(order, field.ID)
		}
	}
	return order
}

func formatValue(value any) string {
	if text, ok := value.(string); ok {
		return text
	}
	data, err := json.Marshal(value)
	if err != nil {
		return "null"
	}
	return string(data)
}

func formatIDs(ids []string) string {
	if len(ids) == 0 {
		return "(none)"
	}
	values := make([]string, len(ids))
	for index, id := range ids {
		values[index] = "`" + id + "`"
	}
	return strings.Join(values, ", ")
}

func encodeEvents(events []ProjectionEvent) ([]byte, error) {
	var buffer bytes.Buffer
	for _, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			return nil, fmt.Errorf("encode projection event %d: %w", event.Sequence, err)
		}
		buffer.Write(data)
		buffer.WriteByte('\n')
	}
	return buffer.Bytes(), nil
}

func makeReceipt(graph model.Graph, inventory Inventory, decisionDigest string, contents map[string][]byte) Receipt {
	files := make([]ReceiptFile, 0, len(OutputNames)-1)
	for _, name := range OutputNames {
		if name == OutputReceipt {
			continue
		}
		data := contents[name]
		files = append(files, ReceiptFile{Name: name, SizeBytes: len(data), SHA256: digest(data)})
	}
	return Receipt{
		ReceiptVersion:        "1.0.0",
		Status:                "verified",
		GraphID:               graph.ID,
		Release:               graph.Release,
		DecisionDigest:        decisionDigest,
		Inventory:             inventory,
		CallerOwnedOutputKinds: append([]string(nil), OutputNames...),
		Files:                 files,
		SelfDigestPolicy:      "receipt digest is verified externally because a file cannot contain its own final digest",
	}
}

func viewManifests(views []ViewProjection) []ViewManifest {
	result := make([]ViewManifest, 0, len(views))
	for _, view := range views {
		result = append(result, ViewManifest{
			ID:                view.ID,
			Label:             view.Label,
			OmittedFieldIDs:   append([]string(nil), view.OmittedFieldIDs...),
			OmittedFieldCount: view.OmittedFieldCount,
		})
	}
	return result
}

func omittedFields(fields []model.ProjectionField, viewID string) []string {
	result := make([]string, 0)
	for _, field := range fields {
		if contains(field.AllowedOmission, viewID) {
			result = append(result, field.ID)
		}
	}
	return result
}

func rankByDecision(order []string) map[string]int {
	ranks := make(map[string]int, len(order))
	for index, state := range order {
		ranks[state] = index
	}
	return ranks
}

func computeDecisionDigest(cases []CaseOutput) (string, error) {
	type material struct {
		ID                 string                    `json:"id"`
		Decision           string                    `json:"decision"`
		DecisionRank       int                       `json:"decision_rank"`
		ClaimStates        []ClaimState              `json:"claim_states"`
		ProofChoices       []ProofChoice              `json:"proof_choices"`
		EvidenceIDs        []string                  `json:"evidence_ids"`
		Counterexamples    []CounterexampleDescription `json:"counterexamples"`
		UnknownFrontier    *model.UnknownFrontier     `json:"unknown_frontier"`
	}
	materials := make([]material, 0, len(cases))
	for _, item := range cases {
		materials = append(materials, material{
			ID:              item.ID,
			Decision:        item.Decision,
			DecisionRank:    item.DecisionRank,
			ClaimStates:     item.ClaimStates,
			ProofChoices:    item.ProofChoices,
			EvidenceIDs:     item.EvidenceIDs,
			Counterexamples: item.CounterexampleDescriptions,
			UnknownFrontier: item.UnknownFrontier,
		})
	}
	data, err := json.Marshal(materials)
	if err != nil {
		return "", err
	}
	return digest(data), nil
}

func completeUnknownFrontier(frontier *model.UnknownFrontier) bool {
	return frontier != nil && frontier.Stage != "" && frontier.Step != "" && frontier.Reason != "" && frontier.UnknownClass != "" && frontier.NextOperation != "" && frontier.BlockedBy != ""
}

func jsonBytes(value any) ([]byte, error) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func digest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
