package metadatav1

type Metadata struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

func FixtureMetadata(namespace, name string) Metadata {
	return Metadata{
		Namespace: namespace,
		Name:      name,
	}
}
