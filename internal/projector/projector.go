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
	"github.com/kimjooyoon/gooo-evidence-resolution-projector/internal/ir"
	"github.com/kimjooyoon/gooo-evidence-resolution-projector/internal/model"
)

const (
	OutputManifest = "projection-manifest.json"
	OutputGraph    = "evidence-graph.json"
	OutputEvents   = "projection-events.ndjson"
	OutputUser     = "user-view.md"
	OutputOperator = "operator-view.md"
	OutputAuditor  = "auditor-view.md"
	OutputReceipt  = "projection-receipt.json"
	OutputReviewer = "reviewer-view.md"
	OutputLanguageMaintainer = "language-maintainer-view.md"
	OutputDossier  = "projection-dossier.md"
	OutputArtifact = "projection-artifact.json"
)

var OutputNames = []string{
	OutputManifest,
	OutputGraph,
	OutputEvents,
	OutputUser,
	OutputOperator,
	OutputAuditor,
	OutputReceipt,
	OutputReviewer,
	OutputLanguageMaintainer,
	OutputDossier,
	OutputArtifact,
}

type GraphOutput struct {
	GraphID                 string           `json:"graph_id"`
	Release                 string           `json:"release"`
	ProjectionSchemaVersion string           `json:"projection_schema_version"`
	DecisionOrder           []string         `json:"decision_order"`
	DecisionDigest          string           `json:"decision_digest"`
	Activities              []model.Activity `json:"activities"`
	Evidence                []model.Evidence `json:"evidence"`
	CanonicalNodes          []model.CanonicalNode `json:"canonical_nodes"`
	CanonicalEdges          []model.CanonicalEdge `json:"canonical_edges"`
	CanonicalSpec            model.CanonicalSpec `json:"canonical_spec"`
	AuthorityBoundary       string           `json:"authority_boundary"`
	ProofCells              []model.ProofCell `json:"proof_cells"`
	Cases                   []CaseOutput     `json:"cases"`
}

type CaseOutput struct {
	ID                         string                      `json:"id"`
	Title                      string                      `json:"title"`
	Summary                    string                      `json:"summary"`
	SummaryRefs                []string                    `json:"summary_refs"`
	ExpectedDecision           string                      `json:"expected_decision"`
	ActivityID                 string                      `json:"activity_id"`
	Decision                   string                      `json:"decision"`
	DecisionRank               int                         `json:"decision_rank"`
	DecisionDigest             string                      `json:"decision_digest"`
	Claims                     []model.Claim               `json:"claims"`
	ClaimStates                []ClaimState                `json:"claim_states"`
	ProofChoices               []ProofChoice               `json:"proof_choices"`
	EvidenceIDs                []string                    `json:"evidence_ids"`
	EvidenceDescriptions       []EvidenceDescription       `json:"evidence_descriptions"`
	UnknownFrontier            *model.UnknownFrontier      `json:"unknown_frontier"`
	Counterexamples            []model.Counterexample      `json:"counterexamples"`
	CounterexampleIDs          []string                    `json:"counterexample_ids"`
	CounterexampleDescriptions []CounterexampleDescription `json:"counterexample_descriptions"`
	SourceContext              string                      `json:"source_context"`
	OperatorNotes              string                      `json:"operator_notes"`
	AuditTrace                 string                      `json:"audit_trace"`
	ImmutableIdentity          string                      `json:"immutable_identity"`
	AuthorityBoundary          string                      `json:"authority_boundary"`
	ClaimEdgeIDs               []string                    `json:"claim_edge_ids"`
	EvidenceEdgeIDs            []string                    `json:"evidence_edge_ids"`
	CounterexampleEdgeIDs      []string                    `json:"counterexample_edge_ids"`
	ProvenanceEdgeIDs          []string                    `json:"provenance_edge_ids"`
	ProofCells                 []model.ProofCell          `json:"proof_cells"`
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
	VisibleNodes      int            `json:"visible_nodes"`
	HiddenNodes       int            `json:"hidden_nodes"`
	FoldedEdges       int            `json:"folded_edges"`
	LostFields        int            `json:"lost_fields"`
	Loss              LossManifest   `json:"loss"`
	FieldOrder        []string       `json:"-"`
}

type ViewProjection struct {
	ID                   string           `json:"id"`
	Role                 string           `json:"role"`
	Label                string           `json:"label"`
	RequestedResolution  string           `json:"requested_resolution"`
	CanonicalGraphSHA256 string           `json:"canonical_graph_sha256"`
	VisibleNodes         int              `json:"visible_nodes"`
	HiddenNodes          int              `json:"hidden_nodes"`
	FoldedEdges          int              `json:"folded_edges"`
	LostFields           int              `json:"lost_fields"`
	OmittedFieldIDs      []string         `json:"omitted_field_ids"`
	OmittedFieldCount    int              `json:"omitted_field_count"`
	Loss                 LossManifest     `json:"loss_manifest"`
	Cases                []CaseProjection `json:"cases"`
}

type LossManifest struct {
	Role                 string   `json:"role"`
	RequestedResolution  string   `json:"requested_resolution"`
	CanonicalGraphSHA256 string   `json:"canonical_graph_sha256"`
	VisibleNodes         int      `json:"visible_nodes"`
	HiddenNodes          int      `json:"hidden_nodes"`
	FoldedEdges          int      `json:"folded_edges"`
	LostFields           int      `json:"lost_fields"`
	HiddenNodeIDs        []string `json:"hidden_node_ids"`
	FoldedEdgeIDs        []string `json:"folded_edge_ids"`
	LostFieldIDs         []string `json:"lost_field_ids"`
	Policy               string   `json:"policy"`
}

type ProjectionEvent struct {
	Sequence             int            `json:"sequence"`
	EventType            string         `json:"event_type"`
	GraphID              string         `json:"graph_id"`
	Role                 string         `json:"role"`
	RequestedResolution  string         `json:"requested_resolution"`
	CanonicalGraphSHA256 string         `json:"canonical_graph_sha256"`
	View                 string         `json:"view"`
	CaseID               string         `json:"case_id"`
	Decision             string         `json:"decision"`
	DecisionDigest       string         `json:"decision_digest"`
	Fields               map[string]any `json:"fields"`
	OmittedFieldIDs      []string       `json:"omitted_field_ids"`
	OmittedFieldCount    int            `json:"omitted_field_count"`
	VisibleNodes         int            `json:"visible_nodes"`
	HiddenNodes          int            `json:"hidden_nodes"`
	FoldedEdges          int            `json:"folded_edges"`
	LostFields           int            `json:"lost_fields"`
	Loss                 LossManifest   `json:"loss_manifest"`
}

type Inventory struct {
	GoooActivities         int            `json:"gooo_activities"`
	CanonicalCases         int            `json:"canonical_cases"`
	CasesByDecision        map[string]int `json:"cases_by_decision"`
	ReaderResolutions      int            `json:"reader_resolutions"`
	CallerOwnedOutputKinds int            `json:"caller_owned_output_kinds"`
	ProjectionEvents       int            `json:"projection_events"`
	RepositoryWrites       int            `json:"repository_writes"`
	InputRepoMutations     int            `json:"input_repo_mutations"`
	SourceMutations        int            `json:"source_mutations"`
	CrossProjectRequiredGates int         `json:"cross_project_required_gates"`
	RuntimeSideEffects     int            `json:"runtime_side_effects"`
}

type ViewManifest struct {
	ID                   string       `json:"id"`
	Role                 string       `json:"role"`
	Label                string       `json:"label"`
	RequestedResolution  string       `json:"requested_resolution"`
	CanonicalGraphSHA256 string       `json:"canonical_graph_sha256"`
	VisibleNodes         int          `json:"visible_nodes"`
	HiddenNodes          int          `json:"hidden_nodes"`
	FoldedEdges          int          `json:"folded_edges"`
	LostFields           int          `json:"lost_fields"`
	OmittedFieldIDs      []string     `json:"omitted_field_ids"`
	OmittedFieldCount    int          `json:"omitted_field_count"`
	Loss                 LossManifest `json:"loss_manifest"`
}

type Manifest struct {
	ManifestVersion        string         `json:"manifest_version"`
	GraphID                string         `json:"graph_id"`
	Release                string         `json:"release"`
	GoooSource             string         `json:"gooo_source"`
	GoooSourceSHA256       string         `json:"gooo_source_sha256"`
	CanonicalGraphSHA256   string         `json:"canonical_graph_sha256"`
	DecisionDigest         string         `json:"decision_digest"`
	DecisionOrder          []string       `json:"decision_order"`
	Inventory              Inventory      `json:"inventory"`
	CallerOwnedOutputKinds []string       `json:"caller_owned_output_kinds"`
	Views                  []ViewManifest `json:"views"`
	SemanticIRSchema        string         `json:"semantic_ir_schema"`
	ExpansionEvaluations    []ExpansionEvaluation `json:"expansion_evaluations"`
	ProofCells              []model.ProofCell `json:"proof_cells"`
	AuthorityBoundary       string         `json:"authority_boundary"`
}

type ReceiptFile struct {
	Name      string `json:"name"`
	SizeBytes int    `json:"size_bytes"`
	SHA256    string `json:"sha256"`
}

type Receipt struct {
	ReceiptVersion         string        `json:"receipt_version"`
	Status                 string        `json:"status"`
	GraphID                string        `json:"graph_id"`
	Release                string        `json:"release"`
	DecisionDigest         string        `json:"decision_digest"`
	Inventory              Inventory     `json:"inventory"`
	CallerOwnedOutputKinds []string      `json:"caller_owned_output_kinds"`
	Files                  []ReceiptFile `json:"content_files"`
	SelfDigestPolicy       string        `json:"self_digest_policy"`
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
	IR             ir.Document
	Expansions     []ExpansionEvaluation
}

type ExpansionRequest struct {
	Role                  string
	RequestedResolution   string
	OriginalGraphAvailable bool
	RoleAvailable         bool
	LossPolicyAvailable   bool
	RestoredFieldIDs      []string
}

type ExpansionEvaluation struct {
	ID                    string                 `json:"id"`
	Role                  string                 `json:"role"`
	Decision              string                 `json:"decision"`
	Reason                string                 `json:"reason"`
	UnknownFrontier       *model.UnknownFrontier `json:"unknown_frontier"`
	RefutingCounterexampleID string              `json:"refuting_counterexample_id"`
}

type ProjectionArtifact struct {
	ArtifactVersion       string                 `json:"artifact_version"`
	Graph                 GraphOutput            `json:"graph"`
	SemanticIRSchema      string                 `json:"semantic_ir_schema"`
	SemanticIR            SemanticIRArtifact     `json:"semantic_ir"`
	CanonicalGraphSHA256  string                 `json:"canonical_graph_sha256"`
	Projections           []ViewProjection       `json:"projections"`
	ExpansionEvaluations  []ExpansionEvaluation  `json:"expansion_evaluations"`
	ProofCells            []model.ProofCell      `json:"proof_cells"`
	Runtime               RuntimeBoundary        `json:"runtime"`
}

type SemanticIRArtifact struct {
	SchemaVersion      string                `json:"schema_version"`
	SourceGraphID      string                `json:"source_graph_id"`
	SourceRelease      string                `json:"source_release"`
	CanonicalNodeCount int                   `json:"canonical_node_count"`
	CanonicalEdgeCount int                   `json:"canonical_edge_count"`
	Nodes              []model.CanonicalNode `json:"nodes"`
	Edges              []model.CanonicalEdge `json:"edges"`
}

type RuntimeBoundary struct {
	RepositoryWrites       int    `json:"repository_writes"`
	SourceMutations        int    `json:"source_mutations"`
	CrossProjectRequiredGates int  `json:"cross_project_required_gates"`
	OutputScope            string `json:"output_scope"`
	Authority              string `json:"authority"`
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
	semanticIR, err := ir.Lower(document)
	if err != nil {
		return Result{}, fmt.Errorf("lower Gooo source to semantic IR: %w", err)
	}
	document = semanticIR.Source
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
		caseOutput, err := evaluateCase(document.Graph, sourceCase, evidenceByID, activityByCase, decisionRanks, semanticIR)
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
		DecisionOrder:           append([]string(nil), document.Graph.Projection.DecisionOrder...),
		DecisionDigest:          decisionDigest,
		Activities:              append([]model.Activity(nil), document.Graph.Activities...),
		Evidence:                append([]model.Evidence(nil), document.Graph.Evidence...),
		CanonicalNodes:          append([]model.CanonicalNode(nil), semanticIR.Nodes...),
		CanonicalEdges:          append([]model.CanonicalEdge(nil), semanticIR.Edges...),
		CanonicalSpec:           document.Graph.Canonical,
		AuthorityBoundary:       document.Graph.Canonical.AuthorityBoundary,
		ProofCells:              append([]model.ProofCell(nil), document.Graph.Projection.ProofCells...),
		Cases:                   cases,
	}

	graphBytes, err := jsonBytes(graph)
	if err != nil {
		return Result{}, fmt.Errorf("encode evidence graph: %w", err)
	}
	canonicalGraphSHA256 := digest(graphBytes)
	views, events, err := buildViews(document.Graph, graph, semanticIR, canonicalGraphSHA256)
	if err != nil {
		return Result{}, err
	}
	expansions := buildExpansionEvaluations(document.Graph.Projection)

	inventory := Inventory{
		GoooActivities:         len(document.Graph.Activities),
		CanonicalCases:         len(cases),
		CasesByDecision:        counts,
		ReaderResolutions:      len(document.Graph.Projection.ReaderRoles),
		CallerOwnedOutputKinds: len(OutputNames),
		ProjectionEvents:       len(events),
		RepositoryWrites:      0,
		InputRepoMutations:     0,
		SourceMutations:        0,
		CrossProjectRequiredGates: 0,
		RuntimeSideEffects:     0,
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
		viewBytes[roleOutputName(view.Role)] = data
	}
	if reviewer, ok := viewBytes[OutputReviewer]; ok {
		viewBytes[OutputOperator] = reviewer
	}

	manifest := Manifest{
		ManifestVersion:        "1.1.0",
		GraphID:                document.Graph.ID,
		Release:                document.Graph.Release,
		GoooSource:             ".gooo/released.gooo",
		GoooSourceSHA256:       digest(sourceBytes),
		CanonicalGraphSHA256:   canonicalGraphSHA256,
		DecisionDigest:         decisionDigest,
		DecisionOrder:          append([]string(nil), document.Graph.Projection.DecisionOrder...),
		Inventory:              inventory,
		CallerOwnedOutputKinds: append([]string(nil), OutputNames...),
		Views:                  viewManifests(views),
		SemanticIRSchema:       ir.SchemaVersion,
		ExpansionEvaluations:   expansions,
		ProofCells:             append([]model.ProofCell(nil), document.Graph.Projection.ProofCells...),
		AuthorityBoundary:      document.Graph.Canonical.AuthorityBoundary,
	}
	manifestBytes, err := jsonBytes(manifest)
	if err != nil {
		return Result{}, fmt.Errorf("encode projection manifest: %w", err)
	}

	contents := map[string][]byte{
		OutputManifest: manifestBytes,
		OutputGraph:    graphBytes,
		OutputEvents:   eventsBytes,
		OutputUser:              viewBytes[OutputUser],
		OutputOperator:          viewBytes[OutputOperator],
		OutputAuditor:           viewBytes[OutputAuditor],
		OutputReviewer:          viewBytes[OutputReviewer],
		OutputLanguageMaintainer: viewBytes[OutputLanguageMaintainer],
	}
	dossier, err := renderDossier(document.Graph, graph, views, expansions, canonicalGraphSHA256)
	if err != nil {
		return Result{}, err
	}
	contents[OutputDossier] = dossier
	artifact, err := jsonBytes(ProjectionArtifact{
		ArtifactVersion:      "1.1.0",
		Graph:                graph,
		SemanticIRSchema:     ir.SchemaVersion,
		SemanticIR: SemanticIRArtifact{
			SchemaVersion:      ir.SchemaVersion,
			SourceGraphID:      document.Graph.ID,
			SourceRelease:      document.Graph.Release,
			CanonicalNodeCount: len(semanticIR.Nodes),
			CanonicalEdgeCount: len(semanticIR.Edges),
			Nodes:              append([]model.CanonicalNode(nil), semanticIR.Nodes...),
			Edges:              append([]model.CanonicalEdge(nil), semanticIR.Edges...),
		},
		CanonicalGraphSHA256: canonicalGraphSHA256,
		Projections:          views,
		ExpansionEvaluations: expansions,
		ProofCells:           append([]model.ProofCell(nil), document.Graph.Projection.ProofCells...),
		Runtime: RuntimeBoundary{
			RepositoryWrites:          0,
			SourceMutations:            0,
			CrossProjectRequiredGates:  0,
			OutputScope:                "CALLER_OWNED_OUTPUT_ONLY",
			Authority:                  "GITHUB_ACTIONS",
		},
	})
	if err != nil {
		return Result{}, fmt.Errorf("encode projection artifact: %w", err)
	}
	contents[OutputArtifact] = artifact
	receipt := makeReceipt(document.Graph, inventory, decisionDigest, contents)
	receiptBytes, err := jsonBytes(receipt)
	if err != nil {
		return Result{}, fmt.Errorf("encode projection receipt: %w", err)
	}
	contents[OutputReceipt] = receiptBytes

	result := Result{
		Source:         document,
		SourceBytes:    append([]byte(nil), sourceBytes...),
		Graph:          graph,
		Views:          views,
		Events:         events,
		Inventory:      inventory,
		Manifest:       manifest,
		Receipt:        receipt,
		OutputContents: contents,
		IR:             semanticIR,
		Expansions:     expansions,
	}
	if err := VerifyConformance(result); err != nil {
		return Result{}, err
	}
	return result, nil
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
	roleIDs := map[string]bool{}
	requiredRoles := map[string]bool{"USER": false, "REVIEWER": false, "LANGUAGE_MAINTAINER": false, "AUDITOR": false}
	for _, role := range projection.ReaderRoles {
		if role.ID == "" || role.Label == "" || role.RequestedResolution == "" || roleIDs[role.ID] {
			return fmt.Errorf("reader roles must have unique identity, label, and requested resolution")
		}
		roleIDs[role.ID] = true
		if _, ok := requiredRoles[role.ID]; !ok {
			return fmt.Errorf("unsupported reader role %q", role.ID)
		}
		requiredRoles[role.ID] = true
		for _, fieldID := range role.MandatoryFields {
			if !fieldIDsContains(projection.Fields, fieldID) {
				return fmt.Errorf("role %q requires unknown field %q", role.ID, fieldID)
			}
		}
		for _, fieldID := range role.LostFields {
			field, ok := projectionField(projection.Fields, fieldID)
			if !ok || field.Invariant {
				return fmt.Errorf("role %q loses unknown or invariant field %q", role.ID, fieldID)
			}
		}
		for _, kind := range append(append([]string{}, role.HiddenNodeKinds...), role.FoldedEdgeKinds...) {
			if kind == "" {
				return fmt.Errorf("role %q declares an empty node or edge kind", role.ID)
			}
		}
	}
	if len(projection.ReaderRoles) != len(requiredRoles) {
		return fmt.Errorf("exactly four reader roles are required")
	}
	for roleID, present := range requiredRoles {
		if !present {
			return fmt.Errorf("required reader role %q is missing", roleID)
		}
	}
	if len(projection.Canonical.NodeKinds) == 0 || len(projection.Canonical.EdgeKinds) == 0 || projection.Canonical.ImmutableIdentity == "" || projection.Canonical.AuthorityBoundary == "" {
		return fmt.Errorf("canonical node, edge, identity, and authority declarations are required")
	}
	for _, requiredKind := range projection.Canonical.RequiredNodeKinds {
		if !contains(projection.Canonical.NodeKinds, requiredKind) {
			return fmt.Errorf("required canonical node kind %q is not declared", requiredKind)
		}
	}
	for _, requiredKind := range projection.Canonical.RequiredEdgeKinds {
		if !contains(projection.Canonical.EdgeKinds, requiredKind) {
			return fmt.Errorf("required canonical edge kind %q is not declared", requiredKind)
		}
	}
	if projection.Expansion.MissingInputDecision != projection.SemanticRules.UnknownState || projection.Expansion.RestoredHiddenDecision != projection.SemanticRules.RefutedState || projection.Expansion.RefutingCounterexampleID == "" || !completeUnknownFrontier(&projection.Expansion.MissingInputFrontier) {
		return fmt.Errorf("expansion policy must fail closed to UNKNOWN or REFUTED with a complete frontier")
	}
	if projection.LossPolicy.MissingInputDecision != projection.SemanticRules.UnknownState {
		return fmt.Errorf("loss policy missing-input decision must be UNKNOWN")
	}
	if len(projection.ProofCells) != 9 {
		return fmt.Errorf("proof cell denominator must contain exactly nine cells")
	}
	proofCellIDs := map[string]bool{}
	proofChoices := map[string]bool{"FOUNDATION": true, "COHERENCE": true, "REGRESSION": true}
	indicators := map[string]bool{"DRIVER": true, "OUTCOME": true, "GUARDRAIL": true}
	for _, cell := range projection.ProofCells {
		if cell.ID == "" || proofCellIDs[cell.ID] || cell.CaseID == "" || cell.ActivityID == "" || !proofChoices[cell.ProofChoice] || !indicators[cell.Indicator] {
			return fmt.Errorf("proof cells must have unique IDs and fixed proof/indicator vocabularies")
		}
		proofCellIDs[cell.ID] = true
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
	for _, role := range projection.ReaderRoles {
		if !sameStrings(role.LostFields, projection.LossPolicy.AllowedLostFields[role.ID]) || !sameStrings(role.HiddenNodeKinds, projection.LossPolicy.AllowedHiddenNodeKinds[role.ID]) || !sameStrings(role.FoldedEdgeKinds, projection.LossPolicy.AllowedFoldedEdgeKinds[role.ID]) {
			return fmt.Errorf("role %q does not match the declared loss policy", role.ID)
		}
		for _, fieldID := range projection.LossPolicy.NeverLoseFields {
			if contains(role.LostFields, fieldID) {
				return fmt.Errorf("role %q loses protected field %q", role.ID, fieldID)
			}
		}
		for _, kind := range projection.LossPolicy.NeverHideNodeKinds {
			if contains(role.HiddenNodeKinds, kind) {
				return fmt.Errorf("role %q hides protected node kind %q", role.ID, kind)
			}
		}
		for _, kind := range projection.LossPolicy.NeverFoldEdgeKinds {
			if contains(role.FoldedEdgeKinds, kind) {
				return fmt.Errorf("role %q folds protected edge kind %q", role.ID, kind)
			}
		}
		for _, field := range projection.Fields {
			if field.Invariant && !contains(role.MandatoryFields, field.ID) {
				return fmt.Errorf("role %q does not declare invariant field %q as mandatory", role.ID, field.ID)
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
	proofCellCases := map[string]bool{}
	proofCellActivities := map[string]bool{}
	for _, cell := range projection.ProofCells {
		if !caseIDs[cell.CaseID] || !activityIDs[cell.ActivityID] || activityIDByCase[cell.CaseID] != cell.ActivityID {
			return fmt.Errorf("proof cell %q is not bound to its canonical case activity", cell.ID)
		}
		proofCellCases[cell.CaseID] = true
		proofCellActivities[cell.ActivityID] = true
	}
	if len(proofCellCases) != len(graph.Cases) || len(proofCellActivities) != len(graph.Activities) {
		return fmt.Errorf("proof cells must cover every canonical case and activity")
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

func evaluateCase(graph model.Graph, sourceCase model.Case, evidenceByID map[string]model.Evidence, activityByCase map[string]model.Activity, ranks map[string]int, semanticIR ir.Document) (CaseOutput, error) {
	if activityByCase[sourceCase.ID].ID == "" {
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
	proofCells := make([]model.ProofCell, 0)
	for _, cell := range graph.Projection.ProofCells {
		if cell.CaseID == sourceCase.ID {
			proofCells = append(proofCells, cell)
		}
	}
	if len(proofCells) == 0 {
		return CaseOutput{}, fmt.Errorf("case %q has no proof cell", sourceCase.ID)
	}
	claimEdgeIDs, evidenceEdgeIDs, counterexampleEdgeIDs, provenanceEdgeIDs := caseEdgeIDs(sourceCase, semanticIR.Edges)
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
		ImmutableIdentity:          sourceCase.ID,
		AuthorityBoundary:          graph.Canonical.AuthorityBoundary,
		ClaimEdgeIDs:               claimEdgeIDs,
		EvidenceEdgeIDs:            evidenceEdgeIDs,
		CounterexampleEdgeIDs:      counterexampleEdgeIDs,
		ProvenanceEdgeIDs:          provenanceEdgeIDs,
		ProofCells:                 proofCells,
	}, nil
}

func buildViews(graph model.Graph, output GraphOutput, semanticIR ir.Document, canonicalGraphSHA256 string) ([]ViewProjection, []ProjectionEvent, error) {
	views := make([]ViewProjection, 0, len(graph.Projection.ReaderRoles))
	events := make([]ProjectionEvent, 0, len(graph.Projection.ReaderRoles)*len(output.Cases))
	sequence := 1
	for _, role := range graph.Projection.ReaderRoles {
		loss, err := makeLossManifest(graph, role, semanticIR, canonicalGraphSHA256)
		if err != nil {
			return nil, nil, err
		}
		projection := ViewProjection{
			ID:                   strings.ToLower(role.ID),
			Role:                 role.ID,
			Label:                role.Label,
			RequestedResolution:  role.RequestedResolution,
			CanonicalGraphSHA256: canonicalGraphSHA256,
			VisibleNodes:         loss.VisibleNodes,
			HiddenNodes:          loss.HiddenNodes,
			FoldedEdges:          loss.FoldedEdges,
			LostFields:           loss.LostFields,
			OmittedFieldIDs:      append([]string(nil), loss.LostFieldIDs...),
			OmittedFieldCount:    loss.LostFields,
			Loss:                 loss,
			Cases:             make([]CaseProjection, 0, len(output.Cases)),
		}
		for _, item := range output.Cases {
			fields, err := projectFieldsForRole(graph.Projection.Fields, role, item)
			if err != nil {
				return nil, nil, err
			}
			caseProjection := CaseProjection{
				CaseID:             item.ID,
				Decision:           item.Decision,
				DecisionDigest:     output.DecisionDigest,
				Fields:             fields,
				OmittedFieldIDs:    append([]string(nil), loss.LostFieldIDs...),
				OmittedFieldCount:  loss.LostFields,
				VisibleNodes:       loss.VisibleNodes,
				HiddenNodes:        loss.HiddenNodes,
				FoldedEdges:        loss.FoldedEdges,
				LostFields:         loss.LostFields,
				Loss:               loss,
				FieldOrder:         includedFieldIDsForRole(graph.Projection.Fields, role),
			}
			projection.Cases = append(projection.Cases, caseProjection)
			events = append(events, ProjectionEvent{
				Sequence:          sequence,
				EventType:         "projection",
				GraphID:           output.GraphID,
				Role:              role.ID,
				RequestedResolution: role.RequestedResolution,
				CanonicalGraphSHA256: canonicalGraphSHA256,
				View:              strings.ToLower(role.ID),
				CaseID:            item.ID,
				Decision:          item.Decision,
				DecisionDigest:    output.DecisionDigest,
				Fields:            fields,
				OmittedFieldIDs:   append([]string(nil), loss.LostFieldIDs...),
				OmittedFieldCount: loss.LostFields,
				VisibleNodes:      loss.VisibleNodes,
				HiddenNodes:       loss.HiddenNodes,
				FoldedEdges:       loss.FoldedEdges,
				LostFields:        loss.LostFields,
				Loss:              loss,
			})
			sequence++
		}
		views = append(views, projection)
	}
	return views, events, nil
}

func makeLossManifest(graph model.Graph, role model.ReaderRole, semanticIR ir.Document, canonicalGraphSHA256 string) (LossManifest, error) {
	hiddenNodeIDs := make([]string, 0)
	for _, node := range semanticIR.Nodes {
		if contains(role.HiddenNodeKinds, node.Kind) && !contains(graph.Canonical.RequiredNodeKinds, node.Kind) {
			hiddenNodeIDs = append(hiddenNodeIDs, node.ID)
		}
	}
	foldedEdgeIDs := make([]string, 0)
	for _, edge := range semanticIR.Edges {
		if contains(role.FoldedEdgeKinds, edge.Kind) && !contains(graph.Canonical.RequiredEdgeKinds, edge.Kind) {
			foldedEdgeIDs = append(foldedEdgeIDs, edge.ID)
		}
	}
	for _, fieldID := range role.MandatoryFields {
		if contains(role.LostFields, fieldID) {
			return LossManifest{}, fmt.Errorf("role %q loses mandatory field %q", role.ID, fieldID)
		}
	}
	return LossManifest{
		Role:                 role.ID,
		RequestedResolution:  role.RequestedResolution,
		CanonicalGraphSHA256: canonicalGraphSHA256,
		VisibleNodes:         len(semanticIR.Nodes) - len(hiddenNodeIDs),
		HiddenNodes:          len(hiddenNodeIDs),
		FoldedEdges:          len(foldedEdgeIDs),
		LostFields:           len(role.LostFields),
		HiddenNodeIDs:        hiddenNodeIDs,
		FoldedEdgeIDs:        foldedEdgeIDs,
		LostFieldIDs:         append([]string(nil), role.LostFields...),
		Policy:               "only declared context may be hidden, folded, or lost; invariant identity, authority, decisions, frontiers, and refutations remain visible",
	}, nil
}

func caseEdgeIDs(sourceCase model.Case, edges []model.CanonicalEdge) ([]string, []string, []string, []string) {
	claimIDs := make(map[string]bool, len(sourceCase.Claims))
	evidenceIDs := map[string]bool{}
	counterexampleIDs := map[string]bool{}
	for _, claim := range sourceCase.Claims {
		claimIDs[claim.ID] = true
		for _, evidenceID := range claim.EvidenceIDs {
			evidenceIDs[evidenceID] = true
		}
	}
	for _, counterexample := range sourceCase.Counterexamples {
		counterexampleIDs[counterexample.ID] = true
		for _, evidenceID := range counterexample.EvidenceIDs {
			evidenceIDs[evidenceID] = true
		}
	}
	claimEdges := make([]string, 0)
	evidenceEdges := make([]string, 0)
	counterexampleEdges := make([]string, 0)
	provenanceEdges := make([]string, 0)
	for _, edge := range edges {
		switch edge.Kind {
		case "case_claim":
			if edge.From == sourceCase.ID {
				claimEdges = append(claimEdges, edge.ID)
			}
		case "claim_evidence":
			if claimIDs[edge.From] {
				evidenceEdges = append(evidenceEdges, edge.ID)
			}
		case "claim_counterexample", "counterexample_evidence":
			if claimIDs[edge.From] || counterexampleIDs[edge.From] || counterexampleIDs[edge.To] {
				counterexampleEdges = append(counterexampleEdges, edge.ID)
			}
		case "evidence_provenance":
			if evidenceIDs[edge.From] {
				provenanceEdges = append(provenanceEdges, edge.ID)
			}
		}
	}
	return claimEdges, evidenceEdges, counterexampleEdges, provenanceEdges
}

func buildExpansionEvaluations(projection model.Projection) []ExpansionEvaluation {
	requests := []struct {
		id      string
		request ExpansionRequest
	}{
		{id: "expansion-closed", request: ExpansionRequest{Role: "USER", RequestedResolution: "SUMMARY", OriginalGraphAvailable: true, RoleAvailable: true, LossPolicyAvailable: true}},
		{id: "expansion-unknown-missing-input", request: ExpansionRequest{Role: "USER", RequestedResolution: "SUMMARY", OriginalGraphAvailable: false, RoleAvailable: true, LossPolicyAvailable: true}},
		{id: "expansion-refuted-restored-hidden-fact", request: ExpansionRequest{Role: "USER", RequestedResolution: "SUMMARY", OriginalGraphAvailable: true, RoleAvailable: true, LossPolicyAvailable: true, RestoredFieldIDs: []string{"source_context"}}},
	}
	results := make([]ExpansionEvaluation, 0, len(requests))
	for _, item := range requests {
		evaluation := EvaluateExpansion(projection, item.request)
		evaluation.ID = item.id
		results = append(results, evaluation)
	}
	return results
}

// EvaluateExpansion deliberately treats expansion as a claim about authority,
// not as a convenience deserialization operation. Missing inputs are UNKNOWN;
// restoring a hidden fact as though the projection were authoritative is
// REFUTED.
func EvaluateExpansion(projection model.Projection, request ExpansionRequest) ExpansionEvaluation {
	if request.Role == "" || request.RequestedResolution == "" || !request.OriginalGraphAvailable || !request.RoleAvailable || !request.LossPolicyAvailable {
		frontier := projection.Expansion.MissingInputFrontier
		return ExpansionEvaluation{
			Role:            request.Role,
			Decision:        projection.Expansion.MissingInputDecision,
			Reason:          "the original graph, reader role, or loss policy is unavailable",
			UnknownFrontier: &frontier,
		}
	}
	if len(request.RestoredFieldIDs) > 0 {
		return ExpansionEvaluation{
			Role:                      request.Role,
			Decision:                  projection.Expansion.RestoredHiddenDecision,
			Reason:                    "re-expansion cannot restore hidden facts as canonical authority",
			RefutingCounterexampleID: projection.Expansion.RefutingCounterexampleID,
		}
	}
	return ExpansionEvaluation{
		Role:     request.Role,
		Decision: projection.SemanticRules.ClosedState,
		Reason:   "the original graph, reader role, and loss policy are present and no hidden fact was restored",
	}
}

func projectFieldsForRole(fields []model.ProjectionField, role model.ReaderRole, item CaseOutput) (map[string]any, error) {
	data, err := objectMap(item)
	if err != nil {
		return nil, fmt.Errorf("normalize case %q: %w", item.ID, err)
	}
	projected := make(map[string]any)
	for _, field := range fields {
		if contains(role.LostFields, field.ID) {
			continue
		}
		value, found := pathValue(data, field.Source)
		if !found {
			if field.Invariant && !strings.HasPrefix(field.Source, "unknown_frontier.") {
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
	fmt.Fprintf(&buffer, "- reader role: `%s`\n", view.Role)
	fmt.Fprintf(&buffer, "- requested resolution: `%s`\n", view.RequestedResolution)
	fmt.Fprintf(&buffer, "- canonical graph digest: `%s`\n", view.CanonicalGraphSHA256)
	fmt.Fprintf(&buffer, "- decision digest: `%s`\n", graph.DecisionDigest)
	fmt.Fprintf(&buffer, "- visible nodes: `%d`\n", view.VisibleNodes)
	fmt.Fprintf(&buffer, "- hidden nodes: `%d`\n", view.HiddenNodes)
	fmt.Fprintf(&buffer, "- folded edges: `%d`\n", view.FoldedEdges)
	fmt.Fprintf(&buffer, "- lost fields: `%d`\n", view.LostFields)
	fmt.Fprintf(&buffer, "- lost field IDs: %s\n", formatIDs(view.OmittedFieldIDs))
	fmt.Fprintf(&buffer, "- hidden node IDs: %s\n", formatIDs(view.Loss.HiddenNodeIDs))
	fmt.Fprintf(&buffer, "- folded edge IDs: %s\n", formatIDs(view.Loss.FoldedEdgeIDs))
	buffer.WriteString("- projection contract: only declared context is loss-declared; immutable identity, authority boundary, evidence, decisions, causal frontiers, and refuting counterexamples remain reverse-referenceable.\n\n")
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

func renderDossier(graph model.Graph, output GraphOutput, views []ViewProjection, expansions []ExpansionEvaluation, canonicalGraphSHA256 string) ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString("# Gooo evidence-resolution projection dossier\n\n")
	fmt.Fprintf(&buffer, "Canonical graph: `%s`\n\n", graph.ID)
	fmt.Fprintf(&buffer, "Canonical graph digest: `%s`\n\n", canonicalGraphSHA256)
	buffer.WriteString("The dossier is generated from `.gooo` through semantic IR. A projection may hide or fold declared detail, but it cannot change decision, immutable identity, authority boundary, causal frontier, evidence, or a refuting counterexample. Decision precedence is `REFUTED > UNKNOWN > CLOSED`.\n\n")
	buffer.WriteString("## Projection loss vectors\n\n")
	buffer.WriteString("| Role | Requested resolution | Visible nodes | Hidden nodes | Folded edges | Lost fields | Canonical graph digest |\n| --- | --- | ---: | ---: | ---: | ---: | --- |\n")
	for _, view := range views {
		fmt.Fprintf(&buffer, "| `%s` | `%s` | %d | %d | %d | %d | `%s` |\n", view.Role, view.RequestedResolution, view.VisibleNodes, view.HiddenNodes, view.FoldedEdges, view.LostFields, view.CanonicalGraphSHA256)
	}
	buffer.WriteString("\n## Case decisions and retained causal facts\n\n")
	for _, item := range output.Cases {
		fmt.Fprintf(&buffer, "### `%s` — `%s`\n\n", item.ID, item.Decision)
		fmt.Fprintf(&buffer, "- immutable identity: `%s`\n- authority boundary: `%s`\n- evidence IDs: %s\n- refuting counterexample IDs: %s\n", item.ImmutableIdentity, item.AuthorityBoundary, formatIDs(item.EvidenceIDs), formatIDs(item.CounterexampleIDs))
		if item.UnknownFrontier != nil {
			frontier := item.UnknownFrontier
			fmt.Fprintf(&buffer, "- UNKNOWN frontier: `%s` / `%s` / `%s` / `%s` / `%s` / `%s`\n", frontier.Stage, frontier.Step, frontier.Reason, frontier.UnknownClass, frontier.NextOperation, frontier.BlockedBy)
		}
		for _, cell := range item.ProofCells {
			fmt.Fprintf(&buffer, "- proof cell `%s`: `%s` / `%s`\n", cell.ID, cell.ProofChoice, cell.Indicator)
		}
		buffer.WriteString("\n")
	}
	buffer.WriteString("## Expansion evaluations\n\n")
	for _, evaluation := range expansions {
		fmt.Fprintf(&buffer, "- `%s`: `%s` — %s", evaluation.ID, evaluation.Decision, evaluation.Reason)
		if evaluation.RefutingCounterexampleID != "" {
			fmt.Fprintf(&buffer, "; refuting counterexample `%s`", evaluation.RefutingCounterexampleID)
		}
		if evaluation.UnknownFrontier != nil {
			frontier := evaluation.UnknownFrontier
			fmt.Fprintf(&buffer, "; frontier `%s/%s/%s/%s/%s/%s`", frontier.Stage, frontier.Step, frontier.Reason, frontier.UnknownClass, frontier.NextOperation, frontier.BlockedBy)
		}
		buffer.WriteString("\n")
	}
	buffer.WriteString("\nRuntime authority: repository_writes=`0`; source_mutations=`0`; cross_project_required_gates=`0`; output scope is caller-owned only; verification authority is GitHub Actions.\n")
	return buffer.Bytes(), nil
}

func roleOutputName(role string) string {
	switch role {
	case "USER":
		return OutputUser
	case "REVIEWER":
		return OutputReviewer
	case "LANGUAGE_MAINTAINER":
		return OutputLanguageMaintainer
	case "AUDITOR":
		return OutputAuditor
	default:
		return strings.ToLower(role) + "-view.md"
	}
}

func includedFieldIDsForRole(fields []model.ProjectionField, role model.ReaderRole) []string {
	order := make([]string, 0, len(fields))
	for _, field := range fields {
		if !contains(role.LostFields, field.ID) {
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
		ReceiptVersion:         "1.1.0",
		Status:                 "verified",
		GraphID:                graph.ID,
		Release:                graph.Release,
		DecisionDigest:         decisionDigest,
		Inventory:              inventory,
		CallerOwnedOutputKinds: append([]string(nil), OutputNames...),
		Files:                  files,
		SelfDigestPolicy:       "receipt digest is verified externally because a file cannot contain its own final digest",
	}
}

func viewManifests(views []ViewProjection) []ViewManifest {
	result := make([]ViewManifest, 0, len(views))
	for _, view := range views {
		result = append(result, ViewManifest{
			ID:                   view.ID,
			Role:                 view.Role,
			Label:                view.Label,
			RequestedResolution:  view.RequestedResolution,
			CanonicalGraphSHA256: view.CanonicalGraphSHA256,
			VisibleNodes:         view.VisibleNodes,
			HiddenNodes:          view.HiddenNodes,
			FoldedEdges:          view.FoldedEdges,
			LostFields:           view.LostFields,
			OmittedFieldIDs:      append([]string{}, view.OmittedFieldIDs...),
			OmittedFieldCount:    view.OmittedFieldCount,
			Loss:                 view.Loss,
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
		ID              string                      `json:"id"`
		Decision        string                      `json:"decision"`
		DecisionRank    int                         `json:"decision_rank"`
		ClaimStates     []ClaimState                `json:"claim_states"`
		ProofChoices    []ProofChoice               `json:"proof_choices"`
		ProofCells      []model.ProofCell           `json:"proof_cells"`
		EvidenceIDs     []string                    `json:"evidence_ids"`
		ClaimEdgeIDs    []string                    `json:"claim_edge_ids"`
		EvidenceEdgeIDs []string                    `json:"evidence_edge_ids"`
		CounterexampleEdgeIDs []string              `json:"counterexample_edge_ids"`
		ProvenanceEdgeIDs []string                  `json:"provenance_edge_ids"`
		ImmutableIdentity string                    `json:"immutable_identity"`
		AuthorityBoundary string                    `json:"authority_boundary"`
		Counterexamples []CounterexampleDescription `json:"counterexamples"`
		UnknownFrontier *model.UnknownFrontier      `json:"unknown_frontier"`
	}
	materials := make([]material, 0, len(cases))
	for _, item := range cases {
		materials = append(materials, material{
			ID:              item.ID,
			Decision:        item.Decision,
			DecisionRank:    item.DecisionRank,
			ClaimStates:     item.ClaimStates,
			ProofChoices:    item.ProofChoices,
			ProofCells:      item.ProofCells,
			EvidenceIDs:     item.EvidenceIDs,
			ClaimEdgeIDs:    item.ClaimEdgeIDs,
			EvidenceEdgeIDs: item.EvidenceEdgeIDs,
			CounterexampleEdgeIDs: item.CounterexampleEdgeIDs,
			ProvenanceEdgeIDs: item.ProvenanceEdgeIDs,
			ImmutableIdentity: item.ImmutableIdentity,
			AuthorityBoundary: item.AuthorityBoundary,
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

// VerifyConformance checks the projection invariants without introducing a
// second semantic evaluator. It compares every role to the already-evaluated
// canonical case and verifies that required causal facts remain present.
func VerifyConformance(result Result) error {
	if len(result.Views) != 4 {
		return fmt.Errorf("conformance requires four reader-role projections")
	}
	if result.Manifest.CanonicalGraphSHA256 != digest(result.OutputContents[OutputGraph]) {
		return fmt.Errorf("projection manifest is not bound to the canonical graph digest")
	}
	graphCases := make(map[string]CaseOutput, len(result.Graph.Cases))
	for _, item := range result.Graph.Cases {
		graphCases[item.ID] = item
	}
	seenRoles := map[string]bool{}
	for _, view := range result.Views {
		if seenRoles[view.Role] {
			return fmt.Errorf("duplicate role projection %q", view.Role)
		}
		seenRoles[view.Role] = true
		if view.CanonicalGraphSHA256 != result.Manifest.CanonicalGraphSHA256 || view.Loss.CanonicalGraphSHA256 != result.Manifest.CanonicalGraphSHA256 {
			return fmt.Errorf("role %q is not bound to the canonical graph digest", view.Role)
		}
		if view.VisibleNodes < 0 || view.HiddenNodes < 0 || view.FoldedEdges < 0 || view.LostFields < 0 || view.VisibleNodes+view.HiddenNodes != len(result.IR.Nodes) {
			return fmt.Errorf("role %q has invalid projection vector", view.Role)
		}
		if view.HiddenNodes != len(view.Loss.HiddenNodeIDs) || view.FoldedEdges != len(view.Loss.FoldedEdgeIDs) || view.LostFields != len(view.Loss.LostFieldIDs) {
			return fmt.Errorf("role %q has an invalid loss manifest", view.Role)
		}
		if len(view.Cases) != len(result.Graph.Cases) {
			return fmt.Errorf("role %q does not project every canonical case", view.Role)
		}
		for _, projected := range view.Cases {
			canonical, ok := graphCases[projected.CaseID]
			if !ok || projected.Decision != canonical.Decision || projected.DecisionDigest != result.Graph.DecisionDigest {
				return fmt.Errorf("role %q changes case %q decision or digest", view.Role, projected.CaseID)
			}
			if projected.VisibleNodes != view.VisibleNodes || projected.HiddenNodes != view.HiddenNodes || projected.FoldedEdges != view.FoldedEdges || projected.LostFields != view.LostFields {
				return fmt.Errorf("role %q case %q does not expose the projection vector", view.Role, projected.CaseID)
			}
			for _, field := range result.Source.Graph.Projection.Fields {
				if field.Invariant && !contains(view.Loss.LostFieldIDs, field.ID) {
					if _, ok := projected.Fields[field.ID]; !ok {
						return fmt.Errorf("role %q hides invariant field %q", view.Role, field.ID)
					}
				}
			}
			if projected.Fields["immutable_identity"] != canonical.ImmutableIdentity || projected.Fields["authority_boundary"] != canonical.AuthorityBoundary {
				return fmt.Errorf("role %q changes identity or authority for case %q", view.Role, projected.CaseID)
			}
			if canonical.UnknownFrontier != nil {
				for _, fieldID := range []string{"unknown.stage", "unknown.step", "unknown.reason", "unknown.unknown_class", "unknown.next_operation", "unknown.blocked_by"} {
					value, present := projected.Fields[fieldID]
					if !present || value == nil || fmt.Sprint(value) == "" {
						return fmt.Errorf("role %q hides UNKNOWN frontier field %q", view.Role, fieldID)
					}
				}
			}
			if canonical.Decision == result.Source.Graph.Projection.SemanticRules.RefutedState {
				value, present := projected.Fields["counterexample_ids"]
				if !present || value == nil || len(canonical.CounterexampleIDs) == 0 {
					return fmt.Errorf("role %q hides refuting counterexamples for case %q", view.Role, projected.CaseID)
				}
			}
		}
	}
	for _, role := range []string{"USER", "REVIEWER", "LANGUAGE_MAINTAINER", "AUDITOR"} {
		if !seenRoles[role] {
			return fmt.Errorf("missing role projection %q", role)
		}
	}
	if len(result.Expansions) != 3 || result.Expansions[0].Decision != result.Source.Graph.Projection.SemanticRules.ClosedState || result.Expansions[1].Decision != result.Source.Graph.Projection.SemanticRules.UnknownState || result.Expansions[2].Decision != result.Source.Graph.Projection.SemanticRules.RefutedState {
		return fmt.Errorf("expansion conformance does not preserve CLOSED, UNKNOWN, REFUTED outcomes")
	}
	if result.Expansions[1].UnknownFrontier == nil || !completeUnknownFrontier(result.Expansions[1].UnknownFrontier) || result.Expansions[2].RefutingCounterexampleID == "" {
		return fmt.Errorf("expansion conformance lost UNKNOWN frontier or refuting counterexample")
	}
	return nil
}

func projectionField(fields []model.ProjectionField, fieldID string) (model.ProjectionField, bool) {
	for _, field := range fields {
		if field.ID == fieldID {
			return field, true
		}
	}
	return model.ProjectionField{}, false
}

func fieldIDsContains(fields []model.ProjectionField, fieldID string) bool {
	_, ok := projectionField(fields, fieldID)
	return ok
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

func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
