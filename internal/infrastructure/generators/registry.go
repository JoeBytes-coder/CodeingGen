package generators

import (
	"configgen/internal/domain"
	"fmt"
)

var generators = map[string]Generator{
	"compose":    &ComposeGenerator{},
	"k8s":        &K8sGenerator{},
	"dockerfile": &DockerfileGenerator{},
	"kustomize":  &KustomizeGenerator{},
}

func GetGenerator(genType string) (Generator, error) {
	gen, exists := generators[genType]
	if !exists {
		return nil, fmt.Errorf("unsupported generator type: %s", genType)
	}
	return gen, nil
}

func GenerateConfig(req domain.ConfigRequest) (domain.ConfigResult, error) {
	gen, err := GetGenerator(req.Type)
	if err != nil {
		return domain.ConfigResult{}, err
	}

	config, err := gen.Generate(req)
	if err != nil {
		return domain.ConfigResult{}, err
	}

	return domain.ConfigResult{Type: req.Type, Config: config}, nil
}
