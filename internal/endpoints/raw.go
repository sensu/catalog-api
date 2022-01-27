package endpoints

import (
	"fmt"
	"os"
	"path/filepath"
)

func renderRaw(endpoint APIEndpoint) error {
	contents, ok := endpoint.GetData().(string)
	if !ok {
		return fmt.Errorf("endpoint data is not a string")
	}

	outputPath := endpoint.GetOutputPath()

	// ensure the parent directory exists
	parent := filepath.Dir(outputPath)
	if err := os.MkdirAll(parent, 0700); err != nil {
		return fmt.Errorf("error creating endpoint parent directory: %w", err)
	}

	// write the endpoint contents to the output path
	if err := os.WriteFile(outputPath, []byte(contents), 0600); err != nil {
		return fmt.Errorf("error creating endpoint file: %w", err)
	}

	return nil
}
