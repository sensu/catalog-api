package types

import (
	"fmt"

	"golang.org/x/mod/semver"
)

// IntegrationVersion is a representation of a single integration version
type IntegrationVersion struct {
	Name          string
	Namespace     string
	Major         int
	Minor         int
	Patch         int
	Prerelease    string
	BuildMetadata string
	GitTag        string
	GitRef        string
}

func (i IntegrationVersion) String() string {
	return fmt.Sprintf("%s/%s:%s", i.Namespace, i.Name, i.SemVer())
}

func (i IntegrationVersion) SemVer() string {
	version := fmt.Sprintf("%d.%d.%d", i.Major, i.Minor, i.Patch)

	if i.Prerelease != "" {
		version = fmt.Sprintf("%s-%s", version, i.Prerelease)
	}

	if i.BuildMetadata != "" {
		version = fmt.Sprintf("%s+%s", version, i.BuildMetadata)
	}

	return version
}

type IntegrationVersionsSlice []IntegrationVersion

func (v IntegrationVersionsSlice) LatestVersion() IntegrationVersion {
	latestVersion := IntegrationVersion{}

	for i, version := range v {
		if i != 0 {
			if semver.Compare(latestVersion.SemVer(), version.SemVer()) == -1 {
				latestVersion = version
			}
		} else {
			latestVersion = version
		}
	}

	return latestVersion
}

// VersionedIntegrations is a mapping of integration names to
// IntegrationVersions
type VersionedIntegrations map[string]IntegrationVersionsSlice

// NamespacedIntegrations is a mapping of namespaces to VersionedIntegrations
type NamespacedIntegrations map[string]VersionedIntegrations
