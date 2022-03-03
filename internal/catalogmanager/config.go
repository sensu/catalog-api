package catalogmanager

import (
	"errors"
	"fmt"

	"github.com/sensu/catalog-api/internal/util"
)

type Config struct {
	StagingDir          string
	ReleaseDir          string
	IntegrationsDirName string
}

func (c Config) validate() error {
	if c.StagingDir == "" {
		return errors.New("staging dir must not be empty")
	}
	if c.ReleaseDir == "" {
		return errors.New("release dir must not be empty")
	}
	return nil
}

func (c Config) StagingChecksum() (string, error) {
	// calculate the sha256 checksum of the generated api
	checksum, err := util.CalculateDirChecksum(c.StagingDir, "staging")
	if err != nil {
		return "", fmt.Errorf("error calculating checksum of staging dir: %w", err)
	}
	return checksum, nil
}
