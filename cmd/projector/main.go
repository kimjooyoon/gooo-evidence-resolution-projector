package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kimjooyoon/gooo-evidence-resolution-projector/internal/projector"
)

func main() {
	source := flag.String("source", ".gooo/released.gooo", "path to the data-only Gooo source")
	output := flag.String("out", "outputs", "caller-owned output directory")
	check := flag.Bool("check", false, "generate and validate the eleven caller-owned outputs")
	flag.Parse()

	result, err := projector.LoadAndProject(*source)
	if err != nil {
		fail(err)
	}
	if err := projector.WriteOutputs(result, *output); err != nil {
		fail(err)
	}
	if *check {
		fmt.Printf("verified graph=%s decision_digest=%s outputs=%d\n", result.Graph.GraphID, result.Graph.DecisionDigest, len(result.OutputContents))
		return
	}
	fmt.Printf("generated graph=%s decision_digest=%s outputs=%d\n", result.Graph.GraphID, result.Graph.DecisionDigest, len(result.OutputContents))
}
func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
