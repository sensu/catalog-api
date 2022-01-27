package types

import (
	"fmt"

	catalogv1 "github.com/sensu/catalog-api/internal/api/catalog/v1"
	metav1 "github.com/sensu/catalog-api/internal/api/metadata/v1"
	"gopkg.in/yaml.v3"
)

type RawWrapper struct {
	Type       string          `json:"type" yaml:"type"`
	APIVersion string          `json:"api_version" yaml:"api_version"`
	Metadata   metav1.Metadata `json:"metadata" yaml:"metadata"`
	Value      yaml.Node       `json:"spec" yaml:"spec"`
}

func RawWrapperFromYAMLBytes(bytes []byte) (RawWrapper, error) {
	raw := RawWrapper{}
	if err := yaml.Unmarshal([]byte(bytes), &raw); err != nil {
		return raw, fmt.Errorf("error unmarshaling raw wrapper: %w", err)
	}
	return raw, nil
}

type Wrapper struct {
	Type       string          `json:"type" yaml:"type"`
	APIVersion string          `json:"api_version" yaml:"api_version"`
	Metadata   metav1.Metadata `json:"metadata" yaml:"metadata"`
	Value      interface{}     `json:"spec" yaml:"spec"`
}

func (w Wrapper) TypeVersion() string {
	return fmt.Sprintf("%s.%s", w.APIVersion, w.Type)
}

func WrapperFromRawWrapper(raw RawWrapper) (Wrapper, error) {
	wrap := Wrapper{
		Type:       raw.Type,
		APIVersion: raw.APIVersion,
		Metadata:   raw.Metadata,
	}
	switch wrap.TypeVersion() {
	case "catalog/v1.Integration":
		integration := catalogv1.Integration{
			Metadata: raw.Metadata,
		}
		if err := raw.Value.Decode(&integration); err != nil {
			return wrap, fmt.Errorf("failed to decode raw value %s: %w", wrap.TypeVersion(), err)
		}
		if err := integration.Validate(); err != nil {
			return wrap, fmt.Errorf("failed to validate %s: %w", wrap.TypeVersion(), err)
		}
		wrap.Value = integration
		return wrap, nil
	default:
		return wrap, fmt.Errorf("invalid resource type version: %s", wrap.TypeVersion())
	}
}
