package projector

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestProjectionPreservesInvariantFieldsAcrossViews(t *testing.T) {
	root := filepath.Join("..", "..")
	result, err := LoadAndProject(filepath.Join(root, ".gooo", "released.gooo"))
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Graph.Cases) != 9 {
		t.Fatalf("cases = %d, want 9", len(result.Graph.Cases))
	}
	if len(result.Views) != 4 {
		t.Fatalf("views = %d, want 4", len(result.Views))
	}
	if len(result.Events) != 36 {
		t.Fatalf("events = %d, want 36", len(result.Events))
	}
	for caseIndex := range result.Graph.Cases {
		first := result.Views[0].Cases[caseIndex]
		for viewIndex := 1; viewIndex < len(result.Views); viewIndex++ {
			other := result.Views[viewIndex].Cases[caseIndex]
			for fieldID, value := range first.Fields {
				if otherValue, ok := other.Fields[fieldID]; ok && !reflect.DeepEqual(value, otherValue) {
					t.Fatalf("case %s field %s changed between views", first.CaseID, fieldID)
				}
			}
			if first.Decision != other.Decision || first.DecisionDigest != other.DecisionDigest {
				t.Fatalf("case %s decision changed between views", first.CaseID)
			}
		}
	}
}

func TestProjectionConformancePreservesFrontiersAndLossVectors(t *testing.T) {
	root := filepath.Join("..", "..")
	result, err := LoadAndProject(filepath.Join(root, ".gooo", "released.gooo"))
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyConformance(result); err != nil {
		t.Fatal(err)
	}
	for _, view := range result.Views {
		if view.CanonicalGraphSHA256 != result.Manifest.CanonicalGraphSHA256 {
			t.Fatalf("role %s is not digest-bound", view.Role)
		}
	}
}
func TestWriteOutputsHasExactlySevenKinds(t *testing.T) {
	root := filepath.Join("..", "..")
	result, err := LoadAndProject(filepath.Join(root, ".gooo", "released.gooo"))
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	if err := WriteOutputs(result, directory); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != len(OutputNames) {
		t.Fatalf("output files = %d, want %d", len(entries), len(OutputNames))
	}
	for _, name := range OutputNames {
		if _, err := os.Stat(filepath.Join(directory, name)); err != nil {
			t.Fatalf("missing output %s: %v", name, err)
		}
	}
}
