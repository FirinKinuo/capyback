package storage

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Config struct {
	StorageType   Type           `yaml:"type"`
	StorageParams map[string]any `yaml:"params"`
}

func (c *Config) convertParamsMapTo(to any) error {
	yamlStorageParams, err := yaml.Marshal(c.StorageParams)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	err = yaml.Unmarshal(yamlStorageParams, to)
	if err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	return nil
}

func (c *Config) ReadStorage() (Storager, error) {
	switch c.StorageType {
	case SwiftStorageType:
		swiftStorageConfig := &SwiftStorageConfig{}

		err := c.convertParamsMapTo(swiftStorageConfig)
		if err != nil {
			return nil, fmt.Errorf("convert params map: %w", err)
		}

		return NewSwiftStorage(swiftStorageConfig), nil

	default:
		return nil, UndefinedStorageTypeErr
	}
}

func (c *Config) ReadWriteParams() (WriteParams, error) {
	switch c.StorageType {
	case SwiftStorageType:
		swiftWriteParams := &SwiftWriteParams{}

		err := c.convertParamsMapTo(swiftWriteParams)
		if err != nil {
			return nil, fmt.Errorf("convert params map: %w", err)
		}

		return swiftWriteParams, nil

	default:
		return nil, UndefinedStorageTypeErr
	}
}
