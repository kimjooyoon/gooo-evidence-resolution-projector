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
