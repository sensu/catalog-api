package endpoints

import (
	"time"

	catalogapiv1 "github.com/sensu/catalog-api/internal/api/catalogapi/v1"
)

// GET /api/version.json
func GenerateVersionEndpoint(basePath string, sha256 string) error {
	version := catalogapiv1.ReleaseVersion{
		ReleaseSHA256: sha256,
		LastUpdated:   time.Now().Unix(),
	}
	endpoint := catalogapiv1.NewVersionEndpoint(basePath, version)
	return renderJSON(endpoint)
}
