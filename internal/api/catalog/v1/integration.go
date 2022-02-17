package v1

import (
	"errors"
	"fmt"

	metav1 "github.com/sensu/catalog-api/internal/api/metadata/v1"
)

type Integration struct {
	Metadata           metav1.Metadata `json:"metadata" yaml:"metadata"`
	Class              string          `json:"class" yaml:"class"`
	Contributors       []string        `json:"contributors" yaml:"contributors"`
	Provider           string          `json:"provider" yaml:"provider"`
	ShortDescription   string          `json:"short_description" yaml:"short_description"`
	SupportedPlatforms []string        `json:"supported_platforms" yaml:"supported_platforms"`
	Tags               []string        `json:"tags" yaml:"tags"`
}

func (i Integration) Validate() error {
	if i.Metadata.Namespace == "" {
		return errors.New("namespace cannot be empty")
	}
	if i.Metadata.Name == "" {
		return errors.New("name cannot be empty")
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
		"application",
		"agent",
		"agent/check",
		"agent/discovery",
		"agent/monitoring",
		"backend",
		"backend/alert",
		"backend/incidents",
		"backend/metrics",
		"backend/events",
		"backend/deregistration",
		"backend/remediation",
		"backend/other",
		"cli",
		"cli/command",
		"universal",
		"universal/runtime",
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
