package v1

import (
	"path"
)

// GET /api/version.json
type VersionEndpoint struct {
	outputPath string
	data       ReleaseVersion
}

func (e VersionEndpoint) GetOutputPath() string { return e.outputPath }
func (e VersionEndpoint) GetData() interface{}  { return e.data }

type ReleaseVersion struct {
	ReleaseSHA256 string `json:"release_sha256" yaml:"release_sha256"`
	LastUpdated   int64  `json:"last_updated" yaml:"last_updated"`
}

func NewVersionEndpoint(basePath string, version ReleaseVersion) VersionEndpoint {
	outputPath := path.Join(
		basePath,
		"version.json")

	return VersionEndpoint{
		outputPath: outputPath,
		data:       version,
	}
}
