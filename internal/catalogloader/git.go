package catalogloader

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"
	"github.com/sensu/catalog-api/internal/integrationloader"
	"github.com/sensu/catalog-api/internal/types"
)

const semverRegex = `(?P<Major>0|[1-9]\d*)\.(?P<Minor>0|[1-9]\d*)\.(?P<Patch>0|[1-9]\d*)(?:-(?P<Prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<BuildMetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`

var (
	ErrUnmatchedGitTag = errors.New("unmatched git tag")

	sourceGit = "git"
)

type GitLoader struct {
	repo                *git.Repository
	integrationsDirName string
}

func NewGitLoader(repo *git.Repository, integrationsDirName string) GitLoader {
	return GitLoader{
		repo:                repo,
		integrationsDirName: integrationsDirName,
	}
}

func (l GitLoader) NewIntegrationLoader(integration types.IntegrationVersion) integrationloader.Loader {
	tagName := integration.TagName()
	integrationPath := integration.Path(l.integrationsDirName)
	return integrationloader.NewGitLoader(l.repo, tagName, integrationPath)
}

func (l GitLoader) LoadIntegrations() (types.Integrations, error) {
	integrations := types.Integrations{}

	tags, err := l.repo.Tags()
	if err != nil {
		return integrations, fmt.Errorf("error determining git tags: %w", err)
	}

	tags.ForEach(func(tagRef *plumbing.Reference) error {
		logger := log.With().Str("tag", tagRef.Name().Short()).Logger()

		iv, err := getIntegrationVersionFromGitTag(tagRef)
		if err != nil {
			if errors.Is(err, ErrUnmatchedGitTag) {
				logger.Warn().Str("reason", err.Error()).Msg("Skipping integration version")
			} else {
				return fmt.Errorf("error parsing git tag - tag: %s, err: %w", tagRef.Name().Short(), err)
			}
			return nil
		}

		integrations = append(integrations, iv)

		logger.Info().
			Str("name", iv.Name).
			Str("namespace", iv.Namespace).
			Str("version", iv.SemVer()).
			Str("source", sourceGit).
			Msg("Found integration version")

		return nil
	})

	return integrations, nil
}

func getIntegrationVersionFromGitTag(tagRef *plumbing.Reference) (types.IntegrationVersion, error) {
	var iv types.IntegrationVersion

	gitTag := tagRef.Name().Short()
	gitRef := tagRef.Hash().String()

	expr := fmt.Sprintf("^(?P<IntegrationNamespace>[a-z0-9_-]+)/(?P<IntegrationName>[a-z0-9_-]+)/%s$", semverRegex)
	r := regexp.MustCompile(expr)
	groupNames := r.SubexpNames()
	groupValues := r.FindStringSubmatch(gitTag)
	if len(groupValues) == 0 {
		return iv, ErrUnmatchedGitTag
	}

	groups := map[string]string{}
	for i, groupValue := range groupValues {
		groups[groupNames[i]] = groupValue
	}

	// verify that all of the required regex group keys were set
	namespace, ok := groups["IntegrationNamespace"]
	if !ok {
		return iv, errors.New("key not found in regex match: IntegrationNamespace")
	}
	name, ok := groups["IntegrationName"]
	if !ok {
		return iv, errors.New("key not found in regex match: IntegrationName")
	}
	majorStr, ok := groups["Major"]
	if !ok {
		return iv, errors.New("key not found in regex match: Major")
	}
	minorStr, ok := groups["Minor"]
	if !ok {
		return iv, errors.New("key not found in regex match: Minor")
	}
	patchStr, ok := groups["Patch"]
	if !ok {
		return iv, errors.New("key not found in regex match: Patch")
	}
	prerelease, ok := groups["Prerelease"]
	if !ok {
		return iv, errors.New("key not found in regex match: Prerelease")
	}
	buildMetadata, ok := groups["BuildMetadata"]
	if !ok {
		return iv, errors.New("key not found in regex match: BuildMetadata")
	}

	// convert major, minor & patch versions to integers
	major, err := strconv.Atoi(majorStr)
	if err != nil {
		return iv, fmt.Errorf("error converting major version to integer: %w", err)
	}
	minor, err := strconv.Atoi(minorStr)
	if err != nil {
		return iv, fmt.Errorf("error converting minor version to integer: %w", err)
	}
	patch, err := strconv.Atoi(patchStr)
	if err != nil {
		return iv, fmt.Errorf("error converting patch version to integer: %w", err)
	}

	iv = types.IntegrationVersion{
		Name:          name,
		Namespace:     namespace,
		Major:         major,
		Minor:         minor,
		Patch:         patch,
		Prerelease:    prerelease,
		BuildMetadata: buildMetadata,
		GitTag:        gitTag,
		GitRef:        gitRef,
		Source:        sourceGit,
	}

	return iv, nil
}
