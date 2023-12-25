package flag

import (
	"github.com/spf13/pflag"
)

// ArchiveFlagSet is a flag set for command with archiving.
type ArchiveFlagSet struct {
	Format string
}

// NewArchiveFlagSet creates a new ArchiveFlagSet.
func NewArchiveFlagSet(defaultFormat string) *ArchiveFlagSet {
	return &ArchiveFlagSet{Format: defaultFormat}
}

// FlagSet returns a flag set for command with archiving.
func (a *ArchiveFlagSet) FlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("config", pflag.PanicOnError)

	flagSet.StringVarP(&a.Format, "format", "f", a.Format, "archive format")

	return flagSet
}
