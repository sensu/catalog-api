package metadatav1

type Metadata struct {
	Name        string            `json:"name" yaml:"name"`
	Namespace   string            `json:"namespace" yaml:"namespace"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

func FixtureMetadata(namespace, name string) Metadata {
	return Metadata{
		Namespace: namespace,
		Name:      name,
	}
}
