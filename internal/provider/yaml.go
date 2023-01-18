package provider

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var tagResolvers = make(map[string]func(*yaml.Node) (*yaml.Node, error))
var tagResolversMutex = &sync.Mutex{}

type CustomTagProcessor struct {
	target interface{}
}

func (i *CustomTagProcessor) UnmarshalYAML(value *yaml.Node) error {
	resolved, err := resolveTags(value)
	if err != nil {
		return err
	}
	return resolved.Decode(i.target)
}

func resolveTags(node *yaml.Node) (*yaml.Node, error) {
	tagResolversMutex.Lock()
	for tag, fn := range tagResolvers {
		if node.Tag == tag {
			return fn(node)
		}
	}
	tagResolversMutex.Unlock()
	if node.Kind == yaml.SequenceNode || node.Kind == yaml.MappingNode {
		var err error
		for i := range node.Content {
			node.Content[i], err = resolveTags(node.Content[i])
			if err != nil {
				return nil, err
			}
		}
	}
	return node, nil
}

func resolveEnv(node *yaml.Node) (*yaml.Node, error) {
	if node.Kind != yaml.ScalarNode {
		return nil, errors.New("!env on a non-scalar node")
	}
	value := os.Getenv(node.Value)
	if value == "" {
		return nil, fmt.Errorf("environment variable %v not set", node.Value)
	}
	node.Value = value
	return node, nil
}

func AddResolvers(tag string, fn func(*yaml.Node) (*yaml.Node, error)) {
	tagResolversMutex.Lock()
	tagResolvers[tag] = fn
	tagResolversMutex.Unlock()
}

func YamlUnmarshal(in []byte, out interface{}) error {
	AddResolvers("!env", resolveEnv)
	err := yaml.Unmarshal(in, &CustomTagProcessor{out})
	return err
}
