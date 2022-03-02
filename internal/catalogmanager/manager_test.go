package catalogmanager

import (
	"errors"
	"testing"

	catalogv1 "github.com/sensu/catalog-api/internal/api/catalog/v1"
	"github.com/sensu/catalog-api/internal/catalogloader"
	mockcatalogloader "github.com/sensu/catalog-api/internal/catalogloader/mocks"
	"github.com/sensu/catalog-api/internal/integrationloader"
	mockintegrationloader "github.com/sensu/catalog-api/internal/integrationloader/mocks"
	"github.com/sensu/catalog-api/internal/types"
	"github.com/stretchr/testify/mock"
)

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
					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").
						Return(catalogv1.Integration{}, errors.New("read error"))
					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").
						Return(types.Integrations{
							types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
						}, nil)
					cl.On("NewIntegrationLoader", mock.Anything, mock.Anything, mock.Anything).
						Return(&il)
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
					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").
						Return(catalogv1.Integration{}, nil)
					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").
						Return(types.Integrations{
							types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
						}, nil)
					cl.On("NewIntegrationLoader", mock.Anything, mock.Anything, mock.Anything).
						Return(&il)
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
					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").
						Return(catalogv1.FixtureIntegration(), nil)
					il.On("LoadResources").
						Return("", errors.New("read error"))
					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").
						Return(types.Integrations{
							types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
						}, nil)
					cl.On("NewIntegrationLoader", mock.Anything, mock.Anything, mock.Anything).
						Return(&il)
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
					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").
						Return(catalogv1.FixtureIntegration(), nil)
					il.On("LoadResources").
						Return("", nil)
					il.On("LoadLogo").
						Return("", errors.New("read error"))
					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").
						Return(types.Integrations{
							types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
						}, nil)
					cl.On("NewIntegrationLoader", mock.Anything, mock.Anything, mock.Anything).
						Return(&il)
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
					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").
						Return(catalogv1.FixtureIntegration(), nil)
					il.On("LoadResources").
						Return("", nil)
					il.On("LoadLogo").
						Return("", nil)
					il.On("LoadReadme").
						Return("", errors.New("read error"))
					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").
						Return(types.Integrations{
							types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
						}, nil)
					cl.On("NewIntegrationLoader", mock.Anything, mock.Anything, mock.Anything).
						Return(&il)
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
					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").
						Return(catalogv1.FixtureIntegration(), nil)
					il.On("LoadResources").
						Return("", nil)
					il.On("LoadLogo").
						Return("", nil)
					il.On("LoadReadme").
						Return("", nil)
					il.On("LoadChangelog").
						Return("", errors.New("read error"))
					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").
						Return(types.Integrations{
							types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
						}, nil)
					cl.On("NewIntegrationLoader", mock.Anything, mock.Anything, mock.Anything).
						Return(&il)
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
					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").
						Return(catalogv1.FixtureIntegration(), nil)
					il.On("LoadResources").
						Return("", nil)
					il.On("LoadLogo").
						Return("", nil)
					il.On("LoadReadme").
						Return("", nil)
					il.On("LoadChangelog").
						Return("", nil)
					il.On("LoadImages").
						Return(integrationloader.Images{}, errors.New("read error"))
					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").
						Return(types.Integrations{
							types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
						}, nil)
					cl.On("NewIntegrationLoader", mock.Anything, mock.Anything, mock.Anything).
						Return(&il)
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
					il := mockintegrationloader.Loader{}
					il.On("LoadConfig").
						Return(catalogv1.FixtureIntegration(), nil)
					il.On("LoadResources").
						Return("", nil)
					il.On("LoadLogo").
						Return("", nil)
					il.On("LoadReadme").
						Return("", nil)
					il.On("LoadChangelog").
						Return("", nil)
					il.On("LoadImages").
						Return(integrationloader.Images{}, nil)
					cl := mockcatalogloader.Loader{}
					cl.On("LoadIntegrations").
						Return(types.Integrations{
							types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
						}, nil)
					cl.On("NewIntegrationLoader", mock.Anything, mock.Anything, mock.Anything).
						Return(&il)
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

func TestCatalogManager_ProcessNamespace(t *testing.T) {
	type fields struct {
		config Config
		loader catalogloader.Loader
	}
	type args struct {
		namespace    string
		integrations types.Integrations
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := CatalogManager{
				config: tt.fields.config,
				loader: tt.fields.loader,
			}
			if err := m.ProcessNamespace(tt.args.namespace, tt.args.integrations); (err != nil) != tt.wantErr {
				t.Errorf("CatalogManager.ProcessNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCatalogManager_ProcessIntegrationVersions(t *testing.T) {
	type fields struct {
		config Config
		loader catalogloader.Loader
	}
	type args struct {
		namespace       string
		integrationName string
		integrations    types.Integrations
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := CatalogManager{
				config: tt.fields.config,
				loader: tt.fields.loader,
			}
			if err := m.ProcessIntegrationVersions(tt.args.namespace, tt.args.integrationName, tt.args.integrations); (err != nil) != tt.wantErr {
				t.Errorf("CatalogManager.ProcessIntegrationVersions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCatalogManager_ProcessIntegrationVersion(t *testing.T) {
	type fields struct {
		config Config
		loader catalogloader.Loader
	}
	type args struct {
		version types.IntegrationVersion
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := CatalogManager{
				config: tt.fields.config,
				loader: tt.fields.loader,
			}
			if err := m.ProcessIntegrationVersion(tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("CatalogManager.ProcessIntegrationVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
