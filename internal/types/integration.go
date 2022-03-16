package types

import (
	"fmt"
	"path"

	semver "github.com/Masterminds/semver/v3"
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
	Source        string
}

func (i IntegrationVersion) Path(base string) string {
	return path.Join(base, i.Namespace, i.Name)
}

func (i IntegrationVersion) TagName() string {
	return fmt.Sprintf("%s/%s/%s", i.Namespace, i.Name, i.SemVer())
}

func FixtureIntegrationVersion(namespace, name string, major, minor, patch int) IntegrationVersion {
	return IntegrationVersion{
		Name:          name,
		Namespace:     namespace,
		Major:         major,
		Minor:         minor,
		Patch:         patch,
		Prerelease:    "",
		BuildMetadata: "",
		GitTag:        fmt.Sprintf("%s/%s/%d.%d.%d", namespace, name, major, minor, patch),
		GitRef:        "d994c6bb648123a17e8f70a966857c546b2a6f94",
		Source:        "git",
	}
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

type Integrations []IntegrationVersion

func (i Integrations) LatestVersion() IntegrationVersion {
	latestVersion := IntegrationVersion{}

	for j, version := range i {
		if j != 0 {
			latest := semver.MustParse(latestVersion.SemVer())
			next := semver.MustParse(version.SemVer())
			if next.GreaterThan(latest) {
				latestVersion = version
			}
		} else {
			latestVersion = version
		}
	}

	return latestVersion
}

func (i Integrations) ByNamespace() NamespacedIntegrations {
	integrations := NamespacedIntegrations{}
	for _, integrationVersion := range i {
		namespace := integrationVersion.Namespace
		integrations[namespace] = append(integrations[namespace], integrationVersion)
	}
	return integrations
}

func (i Integrations) ByName() IntegrationVersions {
	integrations := IntegrationVersions{}
	for _, integrationVersion := range i {
		name := integrationVersion.Name
		integrations[name] = append(integrations[name], integrationVersion)
	}
	return integrations
}

func (i Integrations) FilterByNamespace(namespace string) Integrations {
	integrations := Integrations{}
	for _, integrationVersion := range i {
		if integrationVersion.Namespace == namespace {
			integrations = append(integrations, integrationVersion)
		}
	}
	return integrations
}

func (i Integrations) Versions() []string {
	versions := []string{}
	for _, integrationVersion := range i {
		versions = append(versions, integrationVersion.SemVer())
	}
	return versions
}

// NamespacedIntegrations is a mapping of namespaces to Integrations
type NamespacedIntegrations map[string]Integrations

// IntegrationVersions is a mapping of integration names to Integrations
type IntegrationVersions map[string]Integrations
