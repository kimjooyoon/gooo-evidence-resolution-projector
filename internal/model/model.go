package model

// SourceDocument is the data-only Gooo program consumed by the runtime.
type SourceDocument struct {
	Language        string `json:"language"`
	LanguageVersion string `json:"language_version"`
	Graph           Graph  `json:"graph"`
}

type Graph struct {
	ID         string     `json:"id"`
	Release    string     `json:"release"`
	Projection Projection `json:"projection"`
	Activities []Activity `json:"activities"`
	Evidence   []Evidence `json:"evidence"`
	Cases      []Case     `json:"cases"`
}

type Projection struct {
	SchemaVersion string            `json:"schema_version"`
	DecisionOrder []string          `json:"decision_order"`
	SemanticRules SemanticRules     `json:"semantic_rules"`
	Fields        []ProjectionField `json:"fields"`
	Groups        []ProjectionGroup `json:"groups"`
	Views         []ViewSpec        `json:"views"`
	ReaderRoles   []ReaderRole      `json:"reader_roles"`
	LossPolicy    LossPolicy        `json:"loss_policy"`
	Canonical     CanonicalSpec     `json:"canonical"`
	Expansion     ExpansionRules    `json:"expansion"`
	ProofCells    []ProofCell       `json:"proof_cells"`
}

type SemanticRules struct {
	RefutedState                     string `json:"refuted_state"`
	UnknownState                     string `json:"unknown_state"`
	ClosedState                      string `json:"closed_state"`
	ExternalUtilityZeroEvidenceState string `json:"external_utility_zero_evidence_state"`
}

type ProjectionField struct {
	ID              string   `json:"id"`
	Label           string   `json:"label"`
	Source          string   `json:"source"`
	Invariant       bool     `json:"invariant"`
	AllowedOmission []string `json:"allowed_omission"`
}

type ProjectionGroup struct {
	ID         string   `json:"id"`
	FieldIDs   []string `json:"field_ids"`
	ExactCount int      `json:"exact_count"`
}

type ViewSpec struct {
	ID    string   `json:"id"`
	Label string   `json:"label"`
	Omits []string `json:"omits"`
}

// ReaderRole is an append-only role/resolution contract layered on the
// original v0.1 views. It declares what may be hidden, folded, or lost.
type ReaderRole struct {
	ID                  string   `json:"id"`
	Label               string   `json:"label"`
	RequestedResolution string   `json:"requested_resolution"`
	MandatoryFields     []string `json:"mandatory_fields"`
	HiddenNodeKinds     []string `json:"hidden_node_kinds"`
	FoldedEdgeKinds     []string `json:"folded_edge_kinds"`
	LostFields          []string `json:"lost_fields"`
}

type LossPolicy struct {
	AllowedLostFields      map[string][]string `json:"allowed_lost_fields"`
	AllowedHiddenNodeKinds map[string][]string `json:"allowed_hidden_node_kinds"`
	AllowedFoldedEdgeKinds map[string][]string `json:"allowed_folded_edge_kinds"`
	NeverLoseFields        []string            `json:"never_lose_fields"`
	NeverHideNodeKinds     []string            `json:"never_hide_node_kinds"`
	NeverFoldEdgeKinds     []string            `json:"never_fold_edge_kinds"`
	MissingInputDecision   string              `json:"missing_input_decision"`
}

type CanonicalSpec struct {
	NodeKinds         []string `json:"node_kinds"`
	EdgeKinds         []string `json:"edge_kinds"`
	RequiredNodeKinds []string `json:"required_node_kinds"`
	RequiredEdgeKinds []string `json:"required_edge_kinds"`
	ImmutableIdentity string   `json:"immutable_identity"`
	AuthorityBoundary string   `json:"authority_boundary"`
}

type ExpansionRules struct {
	MissingInputDecision     string          `json:"missing_input_decision"`
	RestoredHiddenDecision   string          `json:"restored_hidden_decision"`
	RefutingCounterexampleID string          `json:"refuting_counterexample_id"`
	MissingInputFrontier     UnknownFrontier `json:"missing_input_frontier"`
}

type ProofCell struct {
	ID          string `json:"id"`
	CaseID      string `json:"case_id"`
	ActivityID  string `json:"activity_id"`
	ProofChoice string `json:"proof_choice"`
	Indicator   string `json:"indicator"`
}

type CanonicalNode struct {
	ID                string `json:"id"`
	Kind              string `json:"kind"`
	Label             string `json:"label"`
	ImmutableIdentity bool   `json:"immutable_identity"`
	AuthorityBoundary string `json:"authority_boundary"`
}

type CanonicalEdge struct {
	ID                string `json:"id"`
	Kind              string `json:"kind"`
	From              string `json:"from"`
	To                string `json:"to"`
	AuthorityBoundary string `json:"authority_boundary"`
}

type Activity struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Binding Binding `json:"binding"`
}

type Binding struct {
	CaseID string `json:"case_id"`
}

type Evidence struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	SourceRef   string `json:"source_ref"`
}

type Case struct {
	ID               string           `json:"id"`
	Title            string           `json:"title"`
	Summary          string           `json:"summary"`
	SummaryRefs      []string         `json:"summary_refs"`
	ExpectedDecision string           `json:"expected_decision"`
	ActivityID       string           `json:"activity_id"`
	SourceContext    string           `json:"source_context"`
	OperatorNotes    string           `json:"operator_notes"`
	AuditTrace       string           `json:"audit_trace"`
	Claims           []Claim          `json:"claims"`
	UnknownFrontier  *UnknownFrontier `json:"unknown_frontier"`
	Counterexamples  []Counterexample `json:"counterexamples"`
}

type Claim struct {
	ID          string   `json:"id"`
	State       string   `json:"state"`
	ProofChoice string   `json:"proof_choice"`
	EvidenceIDs []string `json:"evidence_ids"`
	Utility     *Utility `json:"utility"`
}

type Utility struct {
	Kind        string   `json:"kind"`
	Name        string   `json:"name"`
	EvidenceIDs []string `json:"evidence_ids"`
}

type UnknownFrontier struct {
	Stage         string `json:"stage"`
	Step          string `json:"step"`
	Reason        string `json:"reason"`
	UnknownClass  string `json:"unknown_class"`
	NextOperation string `json:"next_operation"`
	BlockedBy     string `json:"blocked_by"`
}

type Counterexample struct {
	ID          string   `json:"id"`
	ClaimID     string   `json:"claim_id"`
	Description string   `json:"description"`
	EvidenceIDs []string `json:"evidence_ids"`
}
