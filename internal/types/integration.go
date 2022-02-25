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

type Integrations []IntegrationVersion

func (i Integrations) LatestVersion() IntegrationVersion {
	latestVersion := IntegrationVersion{}

	for j, version := range i {
		if j != 0 {
			if semver.Compare(latestVersion.SemVer(), version.SemVer()) == -1 {
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

func (n NamespacedIntegrations) FilterByLatestVersions() NamespacedIntegrations {
	nsIntegrations := NamespacedIntegrations{}
	for namespace, integrations := range n {
		nsIntegrations[namespace] = append(nsIntegrations[namespace], integrations.LatestVersion())
	}
	return nsIntegrations
}

// IntegrationVersions is a mapping of integration names to Integrations
type IntegrationVersions map[string]Integrations
