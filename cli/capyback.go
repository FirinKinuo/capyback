package cli

import (
	"github.com/FirinKinuo/capyback/cli/operation"

	"github.com/spf13/cobra"
)

const (
	capybackUse  = "capyback"
	capybackDesc = "Capyback – tool for efficient backups, compressing, and storing them into your storage."
)

// Capyback is a root for start application from cli.
type Capyback struct {
	command *cobra.Command
}

// NewCapyback creates a new Capyback.
func NewCapyback(version string, defaultConfigPath string) *Capyback {
	capyback := &Capyback{}

	capyback.command = &cobra.Command{
		Use:     capybackUse,
		Short:   capybackDesc,
		Long:    capybackDesc,
		Version: version,
	}

	defaultCommands := []CommandProvider{
		operation.NewSave(defaultConfigPath),
	}

	capyback.RegisterCommands(defaultCommands...)

	return capyback
}

// RegisterCommands registers commands in Capyback.
func (c *Capyback) RegisterCommands(providers ...CommandProvider) {
	for _, provider := range providers {
		c.command.AddCommand(provider.Command())
	}
}

// Execute executes the command.
func (c *Capyback) Execute() error {
	return c.command.Execute()
}
