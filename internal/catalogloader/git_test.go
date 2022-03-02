package catalogloader

import (
	"reflect"
	"testing"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/sensu/catalog-api/internal/types"
)

func TestGitLoader_LoadIntegrations(t *testing.T) {
	type fields struct {
		repo                *git.Repository
		integrationsDirName string
	}
	tests := []struct {
		name    string
		fields  fields
		want    types.Integrations
		wantErr bool
	}{
		{
			name: "no tags exist",
			fields: fields{
				repo: &git.Repository{
					Storer: memory.NewStorage(),
				},
			},
			want:    types.Integrations{},
			wantErr: false,
		},
		{
			name: "tags exist with no matches",
			fields: fields{
				repo: func() *git.Repository {
					m := memory.NewStorage()

					o := &plumbing.MemoryObject{}
					o.SetType(plumbing.TagObject)
					if _, err := o.Write([]byte{}); err != nil {
						t.Fatal(err)
					}

					ref := plumbing.NewHashReference(plumbing.NewTagReferenceName("test"), o.Hash())
					m.ReferenceStorage.SetReference(ref)

					return &git.Repository{Storer: m}
				}(),
			},
			want: func() types.Integrations {
				integrations := types.Integrations{}
				return integrations
			}(),
			wantErr: false,
		},
		{
			name: "tags exist only matches",
			fields: fields{
				repo: func() *git.Repository {
					m := memory.NewStorage()

					o := &plumbing.MemoryObject{}
					o.SetType(plumbing.TagObject)
					if _, err := o.Write([]byte{}); err != nil {
						t.Fatal(err)
					}

					firstRef := plumbing.NewHashReference(plumbing.NewTagReferenceName("example_ns/example/1.2.3"), o.Hash())
					secondRef := plumbing.NewHashReference(plumbing.NewTagReferenceName("example_ns/example/1.3.0"), o.Hash())
					m.ReferenceStorage.SetReference(firstRef)
					m.ReferenceStorage.SetReference(secondRef)

					return &git.Repository{Storer: m}
				}(),
			},
			want: func() types.Integrations {
				integrations := types.Integrations{
					types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
					types.FixtureIntegrationVersion("example_ns", "example", 1, 3, 0),
				}
				return integrations
			}(),
			wantErr: false,
		},
		{
			name: "tags exist partial matches",
			fields: fields{
				repo: func() *git.Repository {
					m := memory.NewStorage()

					o := &plumbing.MemoryObject{}
					o.SetType(plumbing.TagObject)
					if _, err := o.Write([]byte{}); err != nil {
						t.Fatal(err)
					}

					matchedRef := plumbing.NewHashReference(plumbing.NewTagReferenceName("example_ns/example/1.2.3"), o.Hash())
					unmatchedRef := plumbing.NewHashReference(plumbing.NewTagReferenceName("test"), o.Hash())
					m.ReferenceStorage.SetReference(matchedRef)
					m.ReferenceStorage.SetReference(unmatchedRef)

					return &git.Repository{Storer: m}
				}(),
			},
			want: func() types.Integrations {
				integrations := types.Integrations{
					types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
				}
				return integrations
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := GitLoader{
				repo:                tt.fields.repo,
				integrationsDirName: tt.fields.integrationsDirName,
			}
			got, err := l.LoadIntegrations()
			if (err != nil) != tt.wantErr {
				t.Errorf("GitLoader.LoadIntegrations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitLoader.LoadIntegrations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getIntegrationVersionFromGitTag(t *testing.T) {
	type args struct {
		tagRef *plumbing.Reference
	}
	tests := []struct {
		name       string
		args       args
		want       types.IntegrationVersion
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "unmatched",
			args: args{
				tagRef: func() *plumbing.Reference {
					o := &plumbing.MemoryObject{}
					o.SetType(plumbing.TagObject)
					if _, err := o.Write([]byte{}); err != nil {
						t.Fatal(err)
					}

					ref := plumbing.NewHashReference(plumbing.NewTagReferenceName("test"), o.Hash())
					return ref
				}(),
			},
			want:       types.IntegrationVersion{},
			wantErr:    true,
			wantErrMsg: "unmatched git tag",
		},
		{
			name: "matched",
			args: args{
				tagRef: func() *plumbing.Reference {
					o := &plumbing.MemoryObject{}
					o.SetType(plumbing.TagObject)
					if _, err := o.Write([]byte{}); err != nil {
						t.Fatal(err)
					}

					ref := plumbing.NewHashReference(plumbing.NewTagReferenceName("example_ns/example/1.2.3"), o.Hash())
					return ref
				}(),
			},
			want:    types.FixtureIntegrationVersion("example_ns", "example", 1, 2, 3),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getIntegrationVersionFromGitTag(tt.args.tagRef)
			if (err != nil) != tt.wantErr {
				t.Errorf("getIntegrationVersionFromGitTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.wantErrMsg {
				t.Errorf("getIntegrationVersionFromGitTag() error msg = %v, wantErr %v", err.Error(), tt.wantErrMsg)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getIntegrationVersionFromGitTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
