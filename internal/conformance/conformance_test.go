package conformance

import (
	"path/filepath"
	"testing"

	"github.com/kimjooyoon/gooo-evidence-resolution-projector/internal/projector"
)

func TestReleasedGoooConformance(t *testing.T) {
	root := filepath.Join("..", "..")
	result, err := projector.LoadAndProject(filepath.Join(root, ".gooo", "released.gooo"))
	if err != nil {
		t.Fatal(err)
	}
	if result.Inventory.GoooActivities != 9 || result.Inventory.CanonicalCases != 9 {
		t.Fatalf("inventory activities=%d cases=%d, want 9 and 9", result.Inventory.GoooActivities, result.Inventory.CanonicalCases)
	}
	if result.Inventory.CasesByDecision["CLOSED"] != 3 || result.Inventory.CasesByDecision["UNKNOWN"] != 3 || result.Inventory.CasesByDecision["REFUTED"] != 3 {
		t.Fatalf("decision inventory = %#v", result.Inventory.CasesByDecision)
	}
	if result.Inventory.ReaderResolutions != 3 || result.Inventory.CallerOwnedOutputKinds != 7 || result.Inventory.ProjectionEvents != 27 {
		t.Fatalf("resolution/output/event inventory = %d/%d/%d", result.Inventory.ReaderResolutions, result.Inventory.CallerOwnedOutputKinds, result.Inventory.ProjectionEvents)
	}
	for _, view := range result.Views {
		if view.OmittedFieldCount != len(view.OmittedFieldIDs) {
			t.Fatalf("view %s omission count mismatch", view.ID)
		}
		for _, item := range view.Cases {
			if item.OmittedFieldCount != len(item.OmittedFieldIDs) {
				t.Fatalf("case %s omission count mismatch", item.CaseID)
			}
			if _, ok := item.Fields["decision"]; !ok {
				t.Fatalf("case %s hides decision", item.CaseID)
			}
			if _, ok := item.Fields["evidence_ids"]; !ok {
				t.Fatalf("case %s hides evidence IDs", item.CaseID)
			}
			if item.Decision == "REFUTED" {
				if _, ok := item.Fields["counterexample_ids"]; !ok {
					t.Fatalf("case %s hides refuting evidence", item.CaseID)
				}
			}
		}
	}
}
