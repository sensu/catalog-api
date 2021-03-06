package catalogv1

import (
	"errors"
	"fmt"
	"strings"

	metav1 "github.com/sensu/catalog-api/internal/api/metadata/v1"
)

type Prompt struct {
	Type     string                 `json:"type" yaml:"type"`
	Name     string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Body     string                 `json:"body,omitempty" yaml:"body,omitempty"`
	Title    string                 `json:"title,omitempty" yaml:"title,omitempty"`
	Input    map[string]interface{} `json:"input,omitempty" yaml:"input,omitempty"`
	Required bool                   `json:"required,omitempty" yaml:"required,omitempty"`
}

type ResourcePatch struct {
	Resource ResourcePatchRef         `json:"resource" yaml:"resource"`
	Patches  []map[string]interface{} `json:"patches" yaml:"patches"`
}

type ResourcePatchRef struct {
	Type       string `json:"type" yaml:"type"`
	ApiVersion string `json:"api_version" yaml:"api_version"`
	Name       string `json:"name" yaml:"name"`
}

type PostInstall struct {
	Type  string `json:"type" yaml:"type"`
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	Body  string `json:"body,omitempty" yaml:"body,omitempty"`
}

func (p PostInstall) Validate() error {
	switch p.Type {
	case "markdown":
		if p.Body == "" {
			return errors.New("body cannot be empty for type markdown")
		}
		if p.Title != "" {
			return errors.New("title must be empty for type markdown")
		}
	case "section":
	default:
		return fmt.Errorf("invalid type")
	}
	return nil
}

type Integration struct {
	Metadata           metav1.Metadata `json:"metadata" yaml:"metadata"`
	DisplayName        string          `json:"display_name" yaml:"display_name"`
	Class              string          `json:"class" yaml:"class"`
	Contributors       []string        `json:"contributors" yaml:"contributors"`
	Provider           string          `json:"provider" yaml:"provider"`
	ShortDescription   string          `json:"short_description" yaml:"short_description"`
	SupportedPlatforms []string        `json:"supported_platforms" yaml:"supported_platforms"`
	Tags               []string        `json:"tags" yaml:"tags"`
	Prompts            []Prompt        `json:"prompts,omitempty" yaml:"prompts,omitempty"`
	ResourcePatches    []ResourcePatch `json:"resource_patches,omitempty" yaml:"resource_patches,omitempty"`
	PostInstall        []PostInstall   `json:"post_install,omitempty" yaml:"post_install,omitempty"`
}

func FixtureIntegration(namespace, name string) Integration {
	return Integration{
		Class:        "community",
		Contributors: []string{"@artem", "@olha"},
		DisplayName:  strings.Title(name),
		Metadata:     metav1.FixtureMetadata(namespace, name),
		Prompts: []Prompt{
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
		},
		Provider: "alerts",
		ResourcePatches: []ResourcePatch{
			{
				Resource: ResourcePatchRef{
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
		},
		ShortDescription:   "lorem ipsum",
		SupportedPlatforms: []string{"linux", "darwin"},
		Tags:               []string{"tag1", "tag2"},
	}
}

func (i Integration) Validate() error {
	if i.Metadata.Namespace == "" {
		return errors.New("namespace cannot be empty")
	}
	if i.Metadata.Name == "" {
		return errors.New("name cannot be empty")
	}
	if i.DisplayName == "" {
		return errors.New("display_name cannot be empty")
	}
	if !isValidClass(i.Class) {
		return fmt.Errorf("class must be one of %s", validClasses())
	}
	if !isValidProvider(i.Provider) {
		return fmt.Errorf("provider must be one of %s, got: %s", validProviders(), i.Provider)
	}
	if i.ShortDescription == "" {
		return errors.New("short_description cannot be empty")
	}
	if len(i.Contributors) == 0 {
		return errors.New("one or more contributors must be defined")
	}

	return nil
}

func validProviders() []string {
	return []string{
		"alerts",
		"deregistration",
		"discovery",
		"events",
		"incidents",
		"metrics",
		"monitoring",
		"remediation",
	}
}

func isValidProvider(provider string) bool {
	for _, p := range validProviders() {
		if p == provider {
			return true
		}
	}
	return false
}

func validClasses() []string {
	return []string{
		"community",
		"partner",
		"supported",
		"enterprise",
	}
}

func isValidClass(class string) bool {
	for _, c := range validClasses() {
		if c == class {
			return true
		}
	}
	return false
}
