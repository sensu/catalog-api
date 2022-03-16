package catalogmanager

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	catalogv1 "github.com/sensu/catalog-api/internal/api/catalog/v1"
	catalogapiv1 "github.com/sensu/catalog-api/internal/api/catalogapi/v1"
	"github.com/sensu/catalog-api/internal/catalogloader"
	mockcatalogloader "github.com/sensu/catalog-api/internal/catalogloader/mocks"
	"github.com/sensu/catalog-api/internal/integrationloader"
	mockintegrationloader "github.com/sensu/catalog-api/internal/integrationloader/mocks"
	"github.com/sensu/catalog-api/internal/types"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	log.Logger = log.Output(zerolog.Nop())
	os.Exit(m.Run())
}

func newCatalogManager(tb testing.TB) CatalogManager {
	stagingDir := tb.TempDir()
	releaseDir := tb.TempDir()

	return CatalogManager{
		config: Config{
			StagingDir: stagingDir,
			ReleaseDir: releaseDir,
		},
	}
}

func defaultIntegrations() types.Integrations {
	return types.Integrations{
		types.FixtureIntegrationVersion("foo", "bar", 0, 1, 0),
		types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
		types.FixtureIntegrationVersion("example_ns", "example", 1, 3, 0),
		types.FixtureIntegrationVersion("example_ns", "other", 4, 5, 9),
	}
}

func setupEndpointTest(tb testing.TB, integrations types.Integrations) (CatalogManager, error) {
	cl := mockcatalogloader.Loader{}
	cl.On("LoadIntegrations").Return(integrations, nil)

	for _, integration := range integrations {
		config := catalogv1.FixtureIntegration(integration.Namespace, integration.Name)
		images := integrationloader.Images{
			"image_1.png": "images png data 1",
			"image_2.png": "images png data 2",
		}

		il := mockintegrationloader.Loader{}
		il.On("LoadConfig").Return(config, nil)
		il.On("LoadResources").Return(`[{"api_version": "core/v2"}]`, nil)
		il.On("LoadLogo").Return("png data", nil)
		il.On("LoadReadme").Return("readme markdown", nil)
		il.On("LoadChangelog").Return("changelog markdown", nil)
		il.On("LoadImages").Return(images, nil)

		cl.On("NewIntegrationLoader", integration).Return(&il)
	}

	m := newCatalogManager(tb)
	m.loader = &cl

	if err := m.ProcessCatalog(); err != nil {
		return m, err
	}
	return m, nil
}

// endpoint: /version.json
func TestVersionEndpoint(t *testing.T) {
	integrations := defaultIntegrations()
	m, err := setupEndpointTest(t, integrations)
	if err != nil {
		t.Fatal(err)
	}

	checksum, err := m.config.StagingChecksum()
	if err != nil {
		t.Fatal(err)
	}

	endpoint := path.Join(m.config.ReleaseDir, "version.json")
	b, err := ioutil.ReadFile(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	version := catalogapiv1.ReleaseVersion{}
	if err := json.Unmarshal(b, &version); err != nil {
		t.Fatal(err)
	}
	if version.ReleaseSHA256 != checksum {
		t.Errorf("release_sha_256 got = %v, want %v", version.ReleaseSHA256, checksum)
	}
	if version.LastUpdated < 1 {
		t.Errorf("last_updated field got = %v, want >= 1", version.LastUpdated)
	}
}

// endpoint: /:release_sha256/v1/catalog.json
func TestCatalogEndpoint(t *testing.T) {
	integrations := defaultIntegrations()
	m, err := setupEndpointTest(t, integrations)
	if err != nil {
		t.Fatal(err)
	}

	checksum, err := m.config.StagingChecksum()
	if err != nil {
		t.Fatal(err)
	}

	endpoint := path.Join(m.config.ReleaseDir, checksum, "v1", "catalog.json")
	b, err := ioutil.ReadFile(endpoint)
	if err != nil {
		t.Fatal(err)
	}

	catalog := catalogapiv1.Catalog{}
	if err := json.Unmarshal(b, &catalog); err != nil {
		t.Fatal(err)
	}
	if catalog.NamespacedIntegrations == nil {
		t.Fatal("namespaced_integrations is nil")
	}
	if len(catalog.NamespacedIntegrations) == 0 {
		t.Fatal("namespaced_integrations is empty")
	}

	wantNsCounts := map[string]int{
		"foo":        1,
		"example_ns": 2,
	}
	for namespace, want := range wantNsCounts {
		nsIntegrations, ok := catalog.NamespacedIntegrations[namespace]
		if !ok {
			t.Fatalf("namespace not found in catalog endpoint: %s", namespace)
		}
		got := len(nsIntegrations)
		if got != want {
			t.Fatalf("namespace count mismatch: namespace = %v, got = %v, want %v",
				namespace, got, want)
		}
	}

	tests := []struct {
		name         string
		integration  types.IntegrationVersion
		wantNotFound bool
	}{
		{
			name:        "foo/bar/0.1.0",
			integration: integrations[0],
		},
		{
			name:         "example_ns/example/1.2.3",
			integration:  integrations[1],
			wantNotFound: true,
		},
		{
			name:        "example_ns/example/1.3.0",
			integration: integrations[2],
		},
		{
			name:        "example_ns/other/4.5.9",
			integration: integrations[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namespace := tt.integration.Namespace
			name := tt.integration.Name
			version := tt.integration.SemVer()

			integration, err := catalog.GetIntegration(namespace, name, version)
			gotNotFound := (err != nil)
			if gotNotFound != tt.wantNotFound {
				t.Errorf("wantNotFound: got = %v, want %v", gotNotFound, tt.wantNotFound)
				return
			}
			if err != nil {
				return
			}

			// Class
			wantClass := "community"
			if integration.Class != wantClass {
				t.Errorf("class mismatch: got = %v, want %v",
					integration.Class, wantClass)
			}
			// Contributors
			wantContributors := []string{"@artem", "@olha"}
			if !reflect.DeepEqual(integration.Contributors, wantContributors) {
				t.Errorf("contributors mismatch: got = %v, want %v",
					integration.Contributors, wantContributors)
			}
			// DisplayName
			wantDisplayName := strings.Title(name)
			if integration.DisplayName != wantDisplayName {
				t.Errorf("display name mismatch: got = %v, want %v",
					integration.DisplayName, wantDisplayName)
			}
			// Metadata Namespace
			if integration.Metadata.Namespace != namespace {
				t.Errorf("metadata namespace mismatch: got = %v, want %v",
					integration.Metadata.Namespace, namespace)
			}
			// Metadata Name
			if integration.Metadata.Name != name {
				t.Errorf("metadata name mismatch: got = %v, want %v",
					integration.Metadata.Name, name)
			}
			// Prompts
			wantPrompts := *new([]catalogv1.Prompt)
			if !reflect.DeepEqual(integration.Prompts, wantPrompts) {
				t.Errorf("prompts mismatch: got = %v, want %v",
					integration.Prompts, wantPrompts)
			}
			// Provider
			wantProvider := "alerts"
			if integration.Provider != wantProvider {
				t.Errorf("provider mismatch: got = %v, want %v",
					integration.Provider, wantProvider)
			}
			// ResourcePatches
			wantResourcePatches := *new([]catalogv1.ResourcePatch)
			if !reflect.DeepEqual(integration.ResourcePatches, wantResourcePatches) {
				t.Errorf("resource patches mismatch: got = %v, want %v",
					integration.ResourcePatches, wantResourcePatches)
			}
			// ShortDescription
			wantShortDescription := "lorem ipsum"
			if integration.ShortDescription != wantShortDescription {
				t.Errorf("short description mismatch: got = %v, want %v",
					integration.ShortDescription, wantShortDescription)
			}
			// SupportedPlatforms
			wantSupportedPlatforms := []string{"linux", "darwin"}
			if !reflect.DeepEqual(integration.SupportedPlatforms, wantSupportedPlatforms) {
				t.Errorf("supported platforms mismatch: got = %v, want %v",
					integration.SupportedPlatforms, wantSupportedPlatforms)
			}
			// Tags
			wantTags := []string{"tag1", "tag2"}
			if !reflect.DeepEqual(integration.Tags, wantTags) {
				t.Errorf("tags mismatch: got = %v, want %v",
					integration.Tags, wantTags)
			}
			// Version
			if integration.Version != version {
				t.Errorf("version mismatch: got = %v, want %v",
					integration.Version, version)
			}
		})
	}
}

// endpoint: /:release_sha256/v1/:namespace/:integration.json
func TestIntegrationEndpoint(t *testing.T) {
	integrations := defaultIntegrations()
	m, err := setupEndpointTest(t, integrations)
	if err != nil {
		t.Fatal(err)
	}

	checksum, err := m.config.StagingChecksum()
	if err != nil {
		t.Fatal(err)
	}

	endpoint := path.Join(m.config.ReleaseDir, checksum, "v1", "example_ns", "example.json")
	b, err := ioutil.ReadFile(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	integration := catalogapiv1.IntegrationWithVersions{}
	if err := json.Unmarshal(b, &integration); err != nil {
		t.Fatal(err)
	}

	// Class
	wantClass := "community"
	if integration.Class != wantClass {
		t.Errorf("class mismatch: got = %v, want %v",
			integration.Class, wantClass)
	}
	// Contributors
	wantContributors := []string{"@artem", "@olha"}
	if !reflect.DeepEqual(integration.Contributors, wantContributors) {
		t.Errorf("contributors mismatch: got = %v, want %v",
			integration.Contributors, wantContributors)
	}
	// DisplayName
	wantDisplayName := strings.Title("Example")
	if integration.DisplayName != wantDisplayName {
		t.Errorf("display name mismatch: got = %v, want %v",
			integration.DisplayName, wantDisplayName)
	}
	// Metadata Namespace
	wantNamespace := "example_ns"
	if integration.Metadata.Namespace != wantNamespace {
		t.Errorf("metadata namespace mismatch: got = %v, want %v",
			integration.Metadata.Namespace, wantNamespace)
	}
	// Metadata Name
	wantName := "example"
	if integration.Metadata.Name != wantName {
		t.Errorf("metadata name mismatch: got = %v, want %v",
			integration.Metadata.Name, wantName)
	}
	// Prompts
	wantPrompts := []catalogv1.Prompt{
		{
			Type:  "section",
			Title: "Example Section",
		},
		{
			Type: "question",
			Name: "employer",
			Input: map[string]interface{}{
				"type":        "string",
				"title":       "Who does Number Two work for?",
				"description": "Provide the name of the person",
				"default":     "Dr. Evil",
			},
		},
	}
	if !reflect.DeepEqual(integration.Prompts, wantPrompts) {
		t.Errorf("prompts mismatch: got = %v, want %v",
			integration.Prompts, wantPrompts)
	}
	// Provider
	wantProvider := "alerts"
	if integration.Provider != wantProvider {
		t.Errorf("provider mismatch: got = %v, want %v",
			integration.Provider, wantProvider)
	}
	// ResourcePatches
	wantResourcePatches := []catalogv1.ResourcePatch{
		{
			Resource: catalogv1.ResourcePatchRef{
				Type:       "Handler",
				ApiVersion: "core/v2",
				Name:       "foo_handler",
			},
			Patches: []map[string]interface{}{
				{
					"path":  "/spec/id",
					"op":    "replace",
					"value": "[[employer]]",
				},
			},
		},
	}
	if !reflect.DeepEqual(integration.ResourcePatches, wantResourcePatches) {
		t.Errorf("resource patches mismatch: got = %v, want %v",
			integration.ResourcePatches, wantResourcePatches)
	}
	// ShortDescription
	wantShortDescription := "lorem ipsum"
	if integration.ShortDescription != wantShortDescription {
		t.Errorf("short description mismatch: got = %v, want %v",
			integration.ShortDescription, wantShortDescription)
	}
	// SupportedPlatforms
	wantSupportedPlatforms := []string{"linux", "darwin"}
	if !reflect.DeepEqual(integration.SupportedPlatforms, wantSupportedPlatforms) {
		t.Errorf("supported platforms mismatch: got = %v, want %v",
			integration.SupportedPlatforms, wantSupportedPlatforms)
	}
	// Tags
	wantTags := []string{"tag1", "tag2"}
	if !reflect.DeepEqual(integration.Tags, wantTags) {
		t.Errorf("tags mismatch: got = %v, want %v",
			integration.Tags, wantTags)
	}
	// Version
	wantVersion := "1.3.0"
	if integration.Version != wantVersion {
		t.Errorf("version mismatch: got = %v, want %v",
			integration.Version, wantVersion)
	}
	// Versions
	wantVersions := []string{"1.2.3", "1.3.0"}
	if !reflect.DeepEqual(integration.Versions, wantVersions) {
		t.Errorf("versions mismatch: got = %v, want %v",
			integration.Versions, wantVersions)
	}
}

// endpoint: /:release_sha256/v1/:namespace/:name/:version.json
func TestIntegrationVersionEndpoint(t *testing.T) {
	integrations := defaultIntegrations()
	m, err := setupEndpointTest(t, integrations)
	if err != nil {
		t.Fatal(err)
	}

	checksum, err := m.config.StagingChecksum()
	if err != nil {
		t.Fatal(err)
	}

	endpoint := path.Join(m.config.ReleaseDir, checksum, "v1", "example_ns", "example", "1.2.3.json")
	b, err := ioutil.ReadFile(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	integration := catalogapiv1.IntegrationVersion{}
	if err := json.Unmarshal(b, &integration); err != nil {
		t.Fatal(err)
	}

	// Class
	wantClass := "community"
	if integration.Class != wantClass {
		t.Errorf("class mismatch: got = %v, want %v",
			integration.Class, wantClass)
	}
	// Contributors
	wantContributors := []string{"@artem", "@olha"}
	if !reflect.DeepEqual(integration.Contributors, wantContributors) {
		t.Errorf("contributors mismatch: got = %v, want %v",
			integration.Contributors, wantContributors)
	}
	// DisplayName
	wantDisplayName := strings.Title("Example")
	if integration.DisplayName != wantDisplayName {
		t.Errorf("display name mismatch: got = %v, want %v",
			integration.DisplayName, wantDisplayName)
	}
	// Metadata Namespace
	wantNamespace := "example_ns"
	if integration.Metadata.Namespace != wantNamespace {
		t.Errorf("metadata namespace mismatch: got = %v, want %v",
			integration.Metadata.Namespace, wantNamespace)
	}
	// Metadata Name
	wantName := "example"
	if integration.Metadata.Name != wantName {
		t.Errorf("metadata name mismatch: got = %v, want %v",
			integration.Metadata.Name, wantName)
	}
	// Prompts
	wantPrompts := []catalogv1.Prompt{
		{
			Type:  "section",
			Title: "Example Section",
		},
		{
			Type: "question",
			Name: "employer",
			Input: map[string]interface{}{
				"type":        "string",
				"title":       "Who does Number Two work for?",
				"description": "Provide the name of the person",
				"default":     "Dr. Evil",
			},
		},
	}
	if !reflect.DeepEqual(integration.Prompts, wantPrompts) {
		t.Errorf("prompts mismatch: got = %v, want %v",
			integration.Prompts, wantPrompts)
	}
	// Provider
	wantProvider := "alerts"
	if integration.Provider != wantProvider {
		t.Errorf("provider mismatch: got = %v, want %v",
			integration.Provider, wantProvider)
	}
	// ResourcePatches
	wantResourcePatches := []catalogv1.ResourcePatch{
		{
			Resource: catalogv1.ResourcePatchRef{
				Type:       "Handler",
				ApiVersion: "core/v2",
				Name:       "foo_handler",
			},
			Patches: []map[string]interface{}{
				{
					"path":  "/spec/id",
					"op":    "replace",
					"value": "[[employer]]",
				},
			},
		},
	}
	if !reflect.DeepEqual(integration.ResourcePatches, wantResourcePatches) {
		t.Errorf("resource patches mismatch: got = %v, want %v",
			integration.ResourcePatches, wantResourcePatches)
	}
	// ShortDescription
	wantShortDescription := "lorem ipsum"
	if integration.ShortDescription != wantShortDescription {
		t.Errorf("short description mismatch: got = %v, want %v",
			integration.ShortDescription, wantShortDescription)
	}
	// SupportedPlatforms
	wantSupportedPlatforms := []string{"linux", "darwin"}
	if !reflect.DeepEqual(integration.SupportedPlatforms, wantSupportedPlatforms) {
		t.Errorf("supported platforms mismatch: got = %v, want %v",
			integration.SupportedPlatforms, wantSupportedPlatforms)
	}
	// Tags
	wantTags := []string{"tag1", "tag2"}
	if !reflect.DeepEqual(integration.Tags, wantTags) {
		t.Errorf("tags mismatch: got = %v, want %v",
			integration.Tags, wantTags)
	}
	// Version
	wantVersion := "1.2.3"
	if integration.Version != wantVersion {
		t.Errorf("version mismatch: got = %v, want %v",
			integration.Version, wantVersion)
	}
}

// endpoint: /:release_sha256/v1/:namespace/:name/versions.json
func TestIntegrationVersionsEndpoint(t *testing.T) {
	integrations := defaultIntegrations()
	m, err := setupEndpointTest(t, integrations)
	if err != nil {
		t.Fatal(err)
	}

	checksum, err := m.config.StagingChecksum()
	if err != nil {
		t.Fatal(err)
	}

	endpoint := path.Join(m.config.ReleaseDir, checksum, "v1", "example_ns", "example", "versions.json")
	b, err := ioutil.ReadFile(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	versions := catalogapiv1.IntegrationVersions{}
	if err := json.Unmarshal(b, &versions); err != nil {
		t.Fatal(err)
	}

	wantVersions := catalogapiv1.IntegrationVersions{"1.2.3", "1.3.0"}
	if !reflect.DeepEqual(versions, wantVersions) {
		t.Errorf("versions mismatch: got = %v, want %v",
			versions, wantVersions)
	}
}

// endpoint: /:release_sha256/v1/:namespace/:name/:version/sensu-resources.json
func TestIntegrationVersionResourcesEndpoint(t *testing.T) {
	integrations := defaultIntegrations()
	m, err := setupEndpointTest(t, integrations)
	if err := m.ProcessCatalog(); err != nil {
		t.Fatal(err)
	}

	checksum, err := m.config.StagingChecksum()
	if err != nil {
		t.Fatal(err)
	}

	endpoint := path.Join(m.config.ReleaseDir, checksum, "v1", "example_ns", "example", "1.2.3", "sensu-resources.json")
	b, err := ioutil.ReadFile(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	resources := []map[string]interface{}{}
	if err := json.Unmarshal(b, &resources); err != nil {
		t.Fatal(err)
	}

	want := []map[string]interface{}{
		{
			"api_version": "core/v2",
		},
	}
	if !reflect.DeepEqual(resources, want) {
		t.Errorf("resources mismatch: got = %v, want %v",
			resources, want)
	}
}

// endpoint: /:release_sha256/v1/:namespace/:name/:version/README.md
func TestIntegrationVersionReadmeEndpoint(t *testing.T) {
	integrations := defaultIntegrations()
	m, err := setupEndpointTest(t, integrations)
	if err := m.ProcessCatalog(); err != nil {
		t.Fatal(err)
	}

	checksum, err := m.config.StagingChecksum()
	if err != nil {
		t.Fatal(err)
	}

	endpoint := path.Join(m.config.ReleaseDir, checksum, "v1", "example_ns", "example", "1.2.3", "README.md")
	b, err := ioutil.ReadFile(endpoint)
	if err != nil {
		t.Fatal(err)
	}

	got := string(b)
	want := "readme markdown"
	if !reflect.DeepEqual(got, want) {
		t.Errorf("readme mismatch: got = %v, want %v",
			got, want)
	}
}

// endpoint: /:release_sha256/v1/:namespace/:name/:version/CHANGELOG.md
func TestIntegrationVersionChangelogEndpoint(t *testing.T) {
	integrations := defaultIntegrations()
	m, err := setupEndpointTest(t, integrations)
	if err := m.ProcessCatalog(); err != nil {
		t.Fatal(err)
	}

	checksum, err := m.config.StagingChecksum()
	if err != nil {
		t.Fatal(err)
	}

	endpoint := path.Join(m.config.ReleaseDir, checksum, "v1", "example_ns", "example", "1.2.3", "CHANGELOG.md")
	b, err := ioutil.ReadFile(endpoint)
	if err != nil {
		t.Fatal(err)
	}

	got := string(b)
	want := "changelog markdown"
	if !reflect.DeepEqual(got, want) {
		t.Errorf("changelog mismatch: got = %v, want %v",
			got, want)
	}
}

// endpoint: /:release_sha256/v1/:namespace/:name/:version/logo.png
func TestIntegrationVersionLogoEndpoint(t *testing.T) {
	integrations := defaultIntegrations()
	m, err := setupEndpointTest(t, integrations)
	if err := m.ProcessCatalog(); err != nil {
		t.Fatal(err)
	}

	checksum, err := m.config.StagingChecksum()
	if err != nil {
		t.Fatal(err)
	}

	endpoint := path.Join(m.config.ReleaseDir, checksum, "v1", "example_ns", "example", "1.2.3", "logo.png")
	b, err := ioutil.ReadFile(endpoint)
	if err != nil {
		t.Fatal(err)
	}

	got := string(b)
	want := "png data"
	if !reflect.DeepEqual(got, want) {
		t.Errorf("logo mismatch: got = %v, want %v",
			got, want)
	}
}

// endpoint: /:release_sha256/v1/:namespace/:name/:version/img/:image
func TestIntegrationVersionImageEndpoint(t *testing.T) {
	integrations := defaultIntegrations()
	m, err := setupEndpointTest(t, integrations)
	if err := m.ProcessCatalog(); err != nil {
		t.Fatal(err)
	}

	checksum, err := m.config.StagingChecksum()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		namespace   string
		integration string
		version     string
		image       string
		want        string
		wantErr     bool
	}{
		{
			name:        "foo/bar/0.1.0 image_1.png",
			namespace:   "foo",
			integration: "bar",
			version:     "0.1.0",
			image:       "image_1.png",
			want:        "images png data 1",
		},
		{
			name:        "foo/bar/0.1.0 image_2.png",
			namespace:   "foo",
			integration: "bar",
			version:     "0.1.0",
			image:       "image_2.png",
			want:        "images png data 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := path.Join(
				m.config.ReleaseDir,
				checksum,
				"v1",
				tt.namespace,
				tt.integration,
				tt.version,
				"img",
				tt.image,
			)

			b, err := ioutil.ReadFile(endpoint)
			if err != nil {
				t.Fatal(err)
			}

			if string(b) != tt.want {
				t.Errorf("image data mismatch: got = %v, want %v", string(b), tt.want)
			}
		})
	}
}

func TestCatalogManager_ProcessCatalog(t *testing.T) {
	type fields struct {
		config Config
		loader catalogloader.Loader
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "error when loading integrations fails",
			fields: fields{
				config: Config{
					StagingDir: t.TempDir(),
					ReleaseDir: t.TempDir(),
				},
				loader: func() catalogloader.Loader {
					l := mockcatalogloader.Loader{}
					l.On("LoadIntegrations").
						Return(types.Integrations{}, errors.New("load integrations error"))
					return &l
				}(),
			},
			wantErr: true,
		},
		{
			name: "error when load config fails",
			fields: fields{
				config: Config{
					StagingDir: t.TempDir(),
					ReleaseDir: t.TempDir(),
				},
				loader: func() catalogloader.Loader {
					integrations := types.Integrations{
						types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
					}

					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").Return(catalogv1.Integration{}, errors.New("read error"))

					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").Return(integrations, nil)
					cl.On("NewIntegrationLoader", mock.Anything).Return(&il)

					return &cl
				}(),
			},
			wantErr: true,
		},
		{
			name: "error when config validation fails",
			fields: fields{
				config: Config{
					StagingDir: t.TempDir(),
					ReleaseDir: t.TempDir(),
				},
				loader: func() catalogloader.Loader {
					integrations := types.Integrations{
						types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
					}

					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").Return(catalogv1.Integration{}, nil)

					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").Return(integrations, nil)
					cl.On("NewIntegrationLoader", mock.Anything).Return(&il)

					return &cl
				}(),
			},
			wantErr: true,
		},
		{
			name: "error when load resources fails",
			fields: fields{
				config: Config{
					StagingDir: t.TempDir(),
					ReleaseDir: t.TempDir(),
				},
				loader: func() catalogloader.Loader {
					config := catalogv1.FixtureIntegration("foo", "bar")
					integrations := types.Integrations{
						types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
					}

					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").Return(config, nil)
					il.On("LoadResources").Return("", errors.New("read error"))

					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").Return(integrations, nil)
					cl.On("NewIntegrationLoader", mock.Anything).Return(&il)

					return &cl
				}(),
			},
			wantErr: true,
		},
		{
			name: "error when load logo fails",
			fields: fields{
				config: Config{
					StagingDir: t.TempDir(),
					ReleaseDir: t.TempDir(),
				},
				loader: func() catalogloader.Loader {
					config := catalogv1.FixtureIntegration("foo", "bar")
					integrations := types.Integrations{
						types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
					}

					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").Return(config, nil)
					il.On("LoadResources").Return("", nil)
					il.On("LoadLogo").Return("", errors.New("read error"))

					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").Return(integrations, nil)
					cl.On("NewIntegrationLoader", mock.Anything).Return(&il)

					return &cl
				}(),
			},
			wantErr: true,
		},
		{
			name: "error when load readme fails",
			fields: fields{
				config: Config{
					StagingDir: t.TempDir(),
					ReleaseDir: t.TempDir(),
				},
				loader: func() catalogloader.Loader {
					config := catalogv1.FixtureIntegration("foo", "bar")
					integrations := types.Integrations{
						types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
					}

					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").Return(config, nil)
					il.On("LoadResources").Return("", nil)
					il.On("LoadLogo").Return("", nil)
					il.On("LoadReadme").Return("", errors.New("read error"))

					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").Return(integrations, nil)
					cl.On("NewIntegrationLoader", mock.Anything).Return(&il)

					return &cl
				}(),
			},
			wantErr: true,
		},
		{
			name: "error when load readme fails",
			fields: fields{
				config: Config{
					StagingDir: t.TempDir(),
					ReleaseDir: t.TempDir(),
				},
				loader: func() catalogloader.Loader {
					config := catalogv1.FixtureIntegration("foo", "bar")
					integrations := types.Integrations{
						types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
					}

					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").Return(config, nil)
					il.On("LoadResources").Return("", nil)
					il.On("LoadLogo").Return("", nil)
					il.On("LoadReadme").Return("", nil)
					il.On("LoadChangelog").Return("", errors.New("read error"))

					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").Return(integrations, nil)
					cl.On("NewIntegrationLoader", mock.Anything).Return(&il)

					return &cl
				}(),
			},
			wantErr: true,
		},
		{
			name: "error when load images fails",
			fields: fields{
				config: Config{
					StagingDir: t.TempDir(),
					ReleaseDir: t.TempDir(),
				},
				loader: func() catalogloader.Loader {
					config := catalogv1.FixtureIntegration("foo", "bar")
					integrations := types.Integrations{
						types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
					}

					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").Return(config, nil)
					il.On("LoadResources").Return("", nil)
					il.On("LoadLogo").Return("", nil)
					il.On("LoadReadme").Return("", nil)
					il.On("LoadChangelog").Return("", nil)
					il.On("LoadImages").Return(integrationloader.Images{}, errors.New("read error"))

					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").Return(integrations, nil)
					cl.On("NewIntegrationLoader", mock.Anything).Return(&il)

					return &cl
				}(),
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				config: Config{
					StagingDir: t.TempDir(),
					ReleaseDir: t.TempDir(),
				},
				loader: func() catalogloader.Loader {
					config := catalogv1.FixtureIntegration("foo", "bar")
					integrations := types.Integrations{
						types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
					}

					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").Return(config, nil)
					il.On("LoadResources").Return("", nil)
					il.On("LoadLogo").Return("", nil)
					il.On("LoadReadme").Return("", nil)
					il.On("LoadChangelog").Return("", nil)
					il.On("LoadImages").Return(integrationloader.Images{}, nil)

					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").Return(integrations, nil)
					cl.On("NewIntegrationLoader", mock.Anything).Return(&il)

					return &cl
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := CatalogManager{
				config: tt.fields.config,
				loader: tt.fields.loader,
			}
			if err := m.ProcessCatalog(); (err != nil) != tt.wantErr {
				t.Errorf("CatalogManager.ProcessCatalog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
