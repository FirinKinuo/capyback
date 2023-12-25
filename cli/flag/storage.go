package flag

import (
	"fmt"

	"github.com/FirinKinuo/capyback/storage"

	"github.com/spf13/pflag"
)

// StorageFlagSet is a flag set for command with storage using.
type StorageFlagSet struct {
	StorageType storage.Type

	swift *SwiftStorageFlagSet
}

// NewStorageFlagSet creates a new StorageFlagSet.
func NewStorageFlagSet() *StorageFlagSet {
	return &StorageFlagSet{
		swift: NewSwiftStorageFlagSet(),
	}
}

// FlagSet returns a flag set for command with storage using.
func (s *StorageFlagSet) FlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("storage", pflag.PanicOnError)

	flagSet.Var(
		&s.StorageType,
		"storage",
		fmt.Sprintf("Type of storage (%s)", storage.StringAvailableStorages()),
	)

	flagSet.AddFlagSet(s.swift.FlagSet())

	return flagSet
}

// WriteParams returns a write params for storage.
func (s *StorageFlagSet) WriteParams() storage.WriteParams {
	return s.swift.WriteParams()
}

// SwiftStorageFlagSet is a flag set for Swift Storage configuration.
type SwiftStorageFlagSet struct {
	Container   string
	Hash        string
	ContentType string
}

// NewSwiftStorageFlagSet creates a new SwiftStorageFlagSet.
func NewSwiftStorageFlagSet() *SwiftStorageFlagSet {
	return &SwiftStorageFlagSet{}
}

// WriteParams returns a write params for Swift Storage.
func (s *SwiftStorageFlagSet) WriteParams() storage.WriteParams {
	return &storage.SwiftWriteParams{
		Container:   s.Container,
		Hash:        s.Hash,
		ContentType: s.ContentType,
	}
}

// FlagSet returns a flag set for Swift Storage configuration.
func (s *SwiftStorageFlagSet) FlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("swift-storage", pflag.PanicOnError)

	flagSet.StringVar(
		&s.Container,
		"swift-container",
		"",
		"Specify the container in Swift Storage where the objects are stored. This flag is required.",
	)
	flagSet.StringVar(
		&s.Hash,
		"swift-hash",
		"",
		"If you know the MD5 hash of the object ahead of time then set the Hash parameter and it will be sent to the server and the server will check the MD5 itself after the upload.",
	)
	flagSet.StringVar(
		&s.ContentType,
		"swift-content-type",
		"",
		"Set the content type of the object in Swift Storage. This value is used by the system to understand how to handle the object.",
	)
	return flagSet
}
