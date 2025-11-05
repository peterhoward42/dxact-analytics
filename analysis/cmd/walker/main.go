package main

import (
	"fmt"
	"os"

	"github.com/peterhoward42/dxact-analytics/analysis/lib"
	"github.com/peterhoward42/dxact-analytics/analysis/lib/visitors/visitorimplementations"
)

// This program recursively walks the hierarchical file system in which events are stored,
// and performs does some simple aggregate counting operations on what it finds.
//
// For example how many discrete users have launched DrawExact? Etc.
func main() {
	// XXXX todo make this a command line argument.
	root := os.ExpandEnv("$HOME/scratch/drawexact-telemetry/events")

	// The file walker is generic and takes no responsibility for the analysis. Instead we pass
	// in a Vistor implementation to do the analysis.
	visitor := visitorimplementations.NewSimpleCounter()
	walker := lib.NewSimpleWalker(visitor)

	err := walker.Walk(root)
	if err != nil {
		fmt.Printf("XXXX SimpleWalk completed with err: %s\n", err.Error())
		os.Exit(1)
	}

	report := visitor.Report()
	fmt.Printf("XXXX visitor report: %s\n", report)
}
