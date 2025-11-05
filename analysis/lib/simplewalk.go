package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	visitor "github.com/peterhoward42/dxact-analytics/analysis/lib/visitors"
	"github.com/peterhoward42/dxact-analytics/lib"
)

// XXXX todo
type SimpleWalker struct {
	visitor visitor.Visitor
}

func NewSimpleWalker(visitor visitor.Visitor) *SimpleWalker {
	return &SimpleWalker{visitor: visitor}
}

// XXXX todo all of of this file
func (sw *SimpleWalker) Walk(directoryPath string) (err error) {
	err = filepath.WalkDir(directoryPath, sw.processNode)
	return
}

func (sw *SimpleWalker) processNode(path string, d os.DirEntry, err error) error {
	if err != nil {
		fmt.Printf("processNode() has been passed an error for this path %s: %v\n", err, path)
		return nil
	}
	if d.IsDir() {
		return nil
	}
	contents, ok := sw.readFile(path)
	if !ok {
		return nil
	}
	var event lib.EventPayload
	if err := json.Unmarshal(contents, &event); err != nil {
		fmt.Printf("processNode() json.Unmarshal error on this node %s: %s\n", path, err.Error())
		return nil
	}

	// Validate the payload
	if err := event.Validate(); err != nil {
		fmt.Printf("raw event validation error on this node %s: %s\n", path, err.Error())
		return nil
	}

	// Call the visitor provided at construction time to process this event.
	return sw.visitor.Visit(event, path)
}

func (sw *SimpleWalker) readFile(path string) (contents []byte, ok bool) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Cannot open  (continuing) %s: %v\n", path, err)
		return contents, false
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		fmt.Printf("ReadAll() problem (continuing) %s: %v\n", path, err)
		return contents, false
	}
	return bytes, true
}
