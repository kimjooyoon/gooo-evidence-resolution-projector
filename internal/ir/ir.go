// Package ir lowers the data-only .gooo source into the semantic graph used by
// the projector. The source owns meaning; this package only materializes the
// declared identities and typed relationships for deterministic projection.
package ir

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kimjooyoon/gooo-evidence-resolution-projector/internal/model"
)

const SchemaVersion = "gooo/evidence-resolution-projector/semantic-ir/v1"

type Document struct {
	Source model.SourceDocument
	Nodes  []model.CanonicalNode
	Edges  []model.CanonicalEdge
}

func Lower(document model.SourceDocument) (Document, error) {
	graph := document.Graph
	if graph.Projection.Canonical.AuthorityBoundary == "" {
		return Document{}, fmt.Errorf("canonical authority boundary is required")
	}

	nodes := make([]model.CanonicalNode, 0)
	addNode := func(id, kind, label string) {
		nodes = append(nodes, model.CanonicalNode{
			ID:                id,
			Kind:              kind,
			Label:             label,
			ImmutableIdentity: true,
			AuthorityBoundary: graph.Projection.Canonical.AuthorityBoundary,
		})
	}
	addNode("graph:"+graph.ID, "graph", graph.ID)
	addNode("authority:"+graph.Projection.Canonical.AuthorityBoundary, "authority", graph.Projection.Canonical.AuthorityBoundary)
	for _, activity := range graph.Activities {
		addNode(activity.ID, "activity", activity.Name)
	}
	for _, item := range graph.Cases {
		addNode(item.ID, "case", item.Title)
		for _, claim := range item.Claims {
			addNode(claim.ID, "claim", claim.ID)
		}
		if item.UnknownFrontier != nil {
			addNode(unknownNodeID(item.ID), "unknown_frontier", item.ID+" causal frontier")
		}
		for _, counterexample := range item.Counterexamples {
			addNode(counterexample.ID, "counterexample", counterexample.Description)
		}
	}
	for _, evidence := range graph.Evidence {
		addNode(evidence.ID, "evidence", evidence.Description)
		addNode(provenanceNodeID(evidence.SourceRef), "provenance", evidence.SourceRef)
	}

	seenNodes := map[string]bool{}
	for _, node := range nodes {
		if seenNodes[node.ID] {
			return Document{}, fmt.Errorf("duplicate semantic IR node %q", node.ID)
		}
		seenNodes[node.ID] = true
	}

	edges := make([]model.CanonicalEdge, 0)
	addEdge := func(id, kind, from, to string) error {
		if !seenNodes[from] || !seenNodes[to] {
			return fmt.Errorf("semantic IR edge %q has unknown endpoint", id)
		}
		edges = append(edges, model.CanonicalEdge{
			ID:                id,
			Kind:              kind,
			From:              from,
			To:                to,
			AuthorityBoundary: graph.Projection.Canonical.AuthorityBoundary,
		})
		return nil
	}
	if err := addEdge("edge:graph-authority", "authority_boundary", "graph:"+graph.ID, "authority:"+graph.Projection.Canonical.AuthorityBoundary); err != nil {
		return Document{}, err
	}
	for _, activity := range graph.Activities {
		if err := addEdge("edge:activity-binding:"+activity.ID, "activity_binding", activity.ID, activity.Binding.CaseID); err != nil {
			return Document{}, err
		}
	}
	for _, item := range graph.Cases {
		for _, claim := range item.Claims {
			if err := addEdge("edge:case-claim:"+item.ID+":"+claim.ID, "case_claim", item.ID, claim.ID); err != nil {
				return Document{}, err
			}
			for _, evidenceID := range claim.EvidenceIDs {
				if err := addEdge("edge:claim-evidence:"+claim.ID+":"+evidenceID, "claim_evidence", claim.ID, evidenceID); err != nil {
					return Document{}, err
				}
			}
		}
		if item.UnknownFrontier != nil {
			for _, claim := range item.Claims {
				if claim.State == graph.Projection.SemanticRules.UnknownState || (claim.Utility != nil && claim.Utility.Kind == "external" && len(claim.Utility.EvidenceIDs) == 0) {
					if err := addEdge("edge:claim-unknown-frontier:"+claim.ID, "claim_unknown_frontier", claim.ID, unknownNodeID(item.ID)); err != nil {
						return Document{}, err
					}
				}
			}
		}
		for _, counterexample := range item.Counterexamples {
			if err := addEdge("edge:claim-counterexample:"+counterexample.ClaimID+":"+counterexample.ID, "claim_counterexample", counterexample.ClaimID, counterexample.ID); err != nil {
				return Document{}, err
			}
			for _, evidenceID := range counterexample.EvidenceIDs {
				if err := addEdge("edge:counterexample-evidence:"+counterexample.ID+":"+evidenceID, "counterexample_evidence", counterexample.ID, evidenceID); err != nil {
					return Document{}, err
				}
			}
		}
	}
	for _, evidence := range graph.Evidence {
		if err := addEdge("edge:evidence-provenance:"+evidence.ID, "evidence_provenance", evidence.ID, provenanceNodeID(evidence.SourceRef)); err != nil {
			return Document{}, err
		}
	}

	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })
	sort.Slice(edges, func(i, j int) bool { return edges[i].ID < edges[j].ID })
	return Document{Source: document, Nodes: nodes, Edges: edges}, nil
}

func unknownNodeID(caseID string) string {
	return "unknown-frontier:" + caseID
}

func provenanceNodeID(sourceRef string) string {
	return "provenance:" + strings.TrimSpace(sourceRef)
}
