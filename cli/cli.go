package cli

import "github.com/spf13/cobra"

// CommandProvider is a provider for command.
// It is used for register commands in Capyback.
type CommandProvider interface {
	Command() *cobra.Command
}
