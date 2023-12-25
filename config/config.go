package config

import (
	"fmt"
	"github.com/FirinKinuo/capyback/storage"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Storage storage.Config `yaml:"storage"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) readYamlFile(path string) (bytesRead []byte, err error) {
	bytesRead, err = os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return bytesRead, nil
}

// ReadYaml reads configuration from yaml file
func (c *Config) ReadYaml(path string) error {
	yamlBytesRead, err := c.readYamlFile(path)
	if err != nil {
		return fmt.Errorf("read yaml: %w", err)
	}

	err = yaml.Unmarshal(yamlBytesRead, c)
	if err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	return nil
}

func (c *Config) makeConfigFolder(path string) error {
	err := os.MkdirAll(path, 0774)
	if err != nil {
		return fmt.Errorf("make dir: %w", err)
	}

	return nil
}

func (c *Config) createYamlFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0660)
	if err != nil {
		return fmt.Errorf("write yaml: %w", err)
	}

	return nil
}

// CreateYaml creates yaml config file at path with current configuration
func (c *Config) CreateYaml(path string) error {
	err := c.makeConfigFolder(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("make config dir: %w", err)
	}

	marshaledYaml, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	err = c.createYamlFile(path, marshaledYaml)
	if err != nil {
		return fmt.Errorf("create yaml: %w", err)
	}

	return nil
}
