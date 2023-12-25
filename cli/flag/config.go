package flag

import (
	"fmt"
	"github.com/FirinKinuo/capyback/config"
	"github.com/spf13/pflag"
)

// ConfigFlagSet is a flag set for configuration application.
type ConfigFlagSet struct {
	Path string
}

// NewConfigFlagSet creates a new ConfigFlagSet.
func NewConfigFlagSet(defaultPath string) *ConfigFlagSet {
	return &ConfigFlagSet{Path: defaultPath}
}

// FlagSet returns a flag set for configuration application.
func (c *ConfigFlagSet) FlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("config", pflag.PanicOnError)

	flagSet.StringVarP(&c.Path, "config", "c", c.Path, "config path")

	return flagSet
}

// ReadYamlConfig reads a yaml config from the path.
func (c *ConfigFlagSet) ReadYamlConfig() (*config.Config, error) {
	yamlConfig := config.NewConfig()
	err := yamlConfig.ReadYaml(c.Path)
	if err != nil {
		return nil, fmt.Errorf("read yaml config: %w", err)
	}

	return yamlConfig, nil
}
