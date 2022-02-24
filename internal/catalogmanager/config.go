package catalogmanager

import "errors"

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
