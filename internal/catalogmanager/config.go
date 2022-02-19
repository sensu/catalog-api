package catalogmanager

import "errors"

type Config struct {
	RepoDir         string
	StagingDir      string
	ReleaseDir      string
	IntegrationsDir string
}

func (c Config) validate() error {
	if c.RepoDir == "" {
		return errors.New("repo dir must not be empty")
	}
	if c.StagingDir == "" {
		return errors.New("staging dir must not be empty")
	}
	if c.ReleaseDir == "" {
		return errors.New("release dir must not be empty")
	}
	return nil
}
