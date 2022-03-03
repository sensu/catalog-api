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

// endpoint: /version.json
func TestVersionEndpoint(t *testing.T) {
	integrations := types.Integrations{
		types.FixtureIntegrationVersion("foo", "bar", 0, 1, 0),
		types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
		types.FixtureIntegrationVersion("example_ns", "example", 1, 3, 0),
		types.FixtureIntegrationVersion("example_ns", "other", 4, 5, 9),
	}

	cl := mockcatalogloader.Loader{}
	cl.On("LoadIntegrations").Return(integrations, nil)

	for _, integration := range integrations {
		config := catalogv1.FixtureIntegration(integration.Namespace, integration.Name)

		il := mockintegrationloader.Loader{}
		il.On("LoadConfig").Return(config, nil)
		il.On("LoadResources").Return("", nil)
		il.On("LoadLogo").Return("", nil)
		il.On("LoadReadme").Return("", nil)
		il.On("LoadChangelog").Return("", nil)
		il.On("LoadImages").Return(integrationloader.Images{}, nil)

		cl.On(
			"NewIntegrationLoader",
			integration.Namespace,
			integration.Name,
			integration.SemVer(),
		).Return(&il)
	}

	m := newCatalogManager(t)
	m.loader = &cl

	if err := m.ProcessCatalog(); err != nil {
		t.Fatalf("TestVersionEndpoint() error = %v", err)
	}

	checksum, err := m.config.StagingChecksum()
	if err != nil {
		t.Fatalf("TestVersionEndpoint() %v", err)
	}

	endpoint := path.Join(m.config.ReleaseDir, "version.json")
	versionBytes, err := ioutil.ReadFile(endpoint)
	if err != nil {
		t.Fatalf("TestVersionEndpoint() error reading version.json: %v", err)
	}
	version := catalogapiv1.ReleaseVersion{}
	if err := json.Unmarshal(versionBytes, &version); err != nil {
		t.Fatalf("TestVersionEndpoint() error marshalling version.json: %v", err)
	}
	if version.ReleaseSHA256 != checksum {
		t.Errorf("TestVersionEndpoint() release_sha_256 got = %v, want %v", version.ReleaseSHA256, checksum)
	}
	if version.LastUpdated < 1 {
		t.Errorf("TestVersionEndpoint() last_updated field got = %v, want >= 1", version.LastUpdated)
	}
}

// endpoint: /:release_sha256/v1/catalog.json
func TestCatalogEndpoint(t *testing.T) {
	integrations := types.Integrations{
		types.FixtureIntegrationVersion("foo", "bar", 0, 1, 0),
		types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
		types.FixtureIntegrationVersion("example_ns", "example", 1, 3, 0),
		types.FixtureIntegrationVersion("example_ns", "other", 4, 5, 9),
	}

	cl := mockcatalogloader.Loader{}
	cl.On("LoadIntegrations").Return(integrations, nil)

	for _, integration := range integrations {
		config := catalogv1.FixtureIntegration(integration.Namespace, integration.Name)

		il := mockintegrationloader.Loader{}
		il.On("LoadConfig").Return(config, nil)
		il.On("LoadResources").Return("", nil)
		il.On("LoadLogo").Return("", nil)
		il.On("LoadReadme").Return("", nil)
		il.On("LoadChangelog").Return("", nil)
		il.On("LoadImages").Return(integrationloader.Images{}, nil)

		cl.On(
			"NewIntegrationLoader",
			integration.Namespace,
			integration.Name,
			integration.SemVer(),
		).Return(&il)
	}

	m := newCatalogManager(t)
	m.loader = &cl

	if err := m.ProcessCatalog(); err != nil {
		t.Fatalf("TestCatalogEndpoint() error = %v", err)
	}

	checksum, err := m.config.StagingChecksum()
	if err != nil {
		t.Fatalf("TestCatalogEndpoint() error = %v", err)
	}

	endpoint := path.Join(m.config.ReleaseDir, checksum, "v1", "catalog.json")
	catalogBytes, err := ioutil.ReadFile(endpoint)
	if err != nil {
		t.Fatalf("TestCatalogEndpoint() error reading catalog.json: %v", err)
	}
	catalog := catalogapiv1.Catalog{}
	if err := json.Unmarshal(catalogBytes, &catalog); err != nil {
		t.Fatalf("TestCatalogEndpoint() error marshalling catalog.json: %v", err)
	}
	if catalog.NamespacedIntegrations == nil {
		t.Fatal("TestCatalogEndpoint() namespaced_integrations is nil")
	}
	if len(catalog.NamespacedIntegrations) == 0 {
		t.Fatal("TestCatalogEndpoint() namespaced_integrations is empty")
	}

	tests := []struct {
		name        string
		index       int
		cindex      int
		namespace   string
		integration string
		version     string
		wantErr     bool
	}{
		{
			name:        "foo/bar/0.1.0",
			index:       0,
			cindex:      0,
			namespace:   "foo",
			integration: "bar",
			version:     "0.1.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cNamespace, ok := catalog.NamespacedIntegrations[tt.namespace]
			if !ok {
				t.Fatalf("TestCatalogEndpoint() namespace not found in catalog endpoint: %s", tt.namespace)
			}
			if len(cNamespace) <= tt.cindex {
				t.Fatalf("TestCatalogEndpoint() integration not found in catalog: namespace = %s, index = %d", tt.namespace, tt.cindex)
			}
			cIntegration := cNamespace[tt.cindex]
			if len(integrations) <= tt.index {
				t.Fatalf("TestCatalogEndpoint() integration not found: namespace = %s, index = %d", tt.namespace, tt.index)
			}
			wIntegration := integrations[tt.index]

			// Class
			wantClass := "community"
			if cIntegration.Class != wantClass {
				t.Errorf("TestCatalogEndpoint() class mismatch: got = %v, want %v",
					cIntegration.Class, wantClass)
			}
			// Contributors
			wantContributors := []string{"@artem", "@olha"}
			if !reflect.DeepEqual(cIntegration.Contributors, wantContributors) {
				t.Errorf("TestCatalogEndpoint() contributors mismatch: got = %v, want %v",
					cIntegration.Contributors, wantContributors)
			}
			// DisplayName
			wantDisplayName := strings.Title(tt.integration)
			if cIntegration.DisplayName != wantDisplayName {
				t.Errorf("TestCatalogEndpoint() display name mismatch: got = %v, want %v",
					cIntegration.DisplayName, wantDisplayName)
			}
			// Metadata Namespace
			if cIntegration.Metadata.Namespace != wIntegration.Namespace {
				t.Errorf("TestCatalogEndpoint() metadata namespace mismatch: got = %v, want %v",
					cIntegration.Metadata.Namespace, wIntegration.Namespace)
			}
			// Metadata Name
			if cIntegration.Metadata.Name != wIntegration.Name {
				t.Errorf("TestCatalogEndpoint() metadata name mismatch: got = %v, want %v",
					cIntegration.Metadata.Name, wIntegration.Name)
			}
			// Prompts
			wantPrompts := *new([]catalogv1.Prompt)
			if !reflect.DeepEqual(cIntegration.Prompts, wantPrompts) {
				t.Errorf("TestCatalogEndpoint() prompts mismatch: got = %v, want %v",
					cIntegration.Prompts, wantPrompts)
			}
			// Provider
			wantProvider := "alerts"
			if cIntegration.Provider != wantProvider {
				t.Errorf("TestCatalogEndpoint() provider mismatch: got = %v, want %v",
					cIntegration.Provider, wantProvider)
			}
			// ResourcePatches
			wantResourcePatches := *new([]catalogv1.ResourcePatch)
			if !reflect.DeepEqual(cIntegration.ResourcePatches, wantResourcePatches) {
				t.Errorf("TestCatalogEndpoint() resource patches mismatch: got = %v, want %v",
					cIntegration.ResourcePatches, wantResourcePatches)
			}
			// ShortDescription
			wantShortDescription := "lorem ipsum"
			if cIntegration.ShortDescription != wantShortDescription {
				t.Errorf("TestCatalogEndpoint() short description mismatch: got = %v, want %v",
					cIntegration.ShortDescription, wantShortDescription)
			}
			// SupportedPlatforms
			wantSupportedPlatforms := []string{"linux", "darwin"}
			if !reflect.DeepEqual(cIntegration.SupportedPlatforms, wantSupportedPlatforms) {
				t.Errorf("TestCatalogEndpoint() supported platforms mismatch: got = %v, want %v",
					cIntegration.SupportedPlatforms, wantSupportedPlatforms)
			}
			// Tags
			wantTags := []string{"tag1", "tag2"}
			if !reflect.DeepEqual(cIntegration.Tags, wantTags) {
				t.Errorf("TestCatalogEndpoint() tags mismatch: got = %v, want %v",
					cIntegration.Tags, wantTags)
			}
			// Version
			if cIntegration.Version != tt.version {
				t.Errorf("TestCatalogEndpoint() version mismatch: got = %v, want %v",
					cIntegration.Version, tt.version)
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
					cl.On(
						"NewIntegrationLoader",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).Return(&il)

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
					cl.On(
						"NewIntegrationLoader",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).Return(&il)

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
					cl.On(
						"NewIntegrationLoader",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).Return(&il)

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
					cl.On(
						"NewIntegrationLoader",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).Return(&il)
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
					cl.On(
						"NewIntegrationLoader",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).Return(&il)

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
					cl.On(
						"NewIntegrationLoader",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).Return(&il)

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
					cl.On(
						"NewIntegrationLoader",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).Return(&il)

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
					cl.On(
						"NewIntegrationLoader",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).Return(&il)

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
