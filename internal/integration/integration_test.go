//go:build integration

package integration

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"

	"github.com/kimjooyoon/gooo-evidence-resolution-projector/internal/projector"
)

func TestProjectorIntegrationRoundTrip(t *testing.T) {
	root := filepath.Join("..", "..")
	result, err := projector.LoadAndProject(filepath.Join(root, ".gooo", "released.gooo"))
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	if err := projector.WriteOutputs(result, directory); err != nil {
		t.Fatal(err)
	}
	graph, err := os.ReadFile(filepath.Join(directory, projector.OutputGraph))
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(graph)
	if result.Manifest.CanonicalGraphSHA256 != fmtDigest(sum[:]) {
		t.Fatalf("manifest graph digest mismatch")
	}
}
func fmtDigest(data []byte) string {
	const hex = "0123456789abcdef"
	output := make([]byte, len(data)*2)
	for index, value := range data {
		output[index*2] = hex[value>>4]
		output[index*2+1] = hex[value&15]
	}
	return string(output)
}
