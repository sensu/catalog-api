package catalogv1

import (
	"bytes"
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

type Resource map[string]interface{}

type Resources []Resource

func ResourcesFromYAML(data []byte) (Resources, error) {
	resources := Resources{}
	dec := yaml.NewDecoder(bytes.NewReader(data))

	for {
		doc := new(Resource)

		if err := dec.Decode(&doc); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return resources, err
		}
		if doc == nil {
			return resources, errors.New("empty yaml document")
		}
		resources = append(resources, *doc)
	}

	return resources, nil
}
