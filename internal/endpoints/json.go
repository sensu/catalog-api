package endpoints

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func renderJSON(endpoint APIEndpoint) error {
	contents, err := json.Marshal(endpoint.GetData())
	if err != nil {
		return fmt.Errorf("error generating endpoint: %w", err)
	}

	outputPath := endpoint.GetOutputPath()

	// ensure the parent directory exists
	parent := filepath.Dir(outputPath)
	if err := os.MkdirAll(parent, 0700); err != nil {
		return fmt.Errorf("error creating endpoint parent directory: %w", err)
	}

	// write the endpoint contents to the output path
	if err := os.WriteFile(outputPath, contents, 0600); err != nil {
		return fmt.Errorf("error creating endpoint file: %w", err)
	}

	return nil
}
