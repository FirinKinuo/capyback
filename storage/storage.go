package storage

import (
	"context"
	"errors"
	"io"
	"strings"
)

type Type string

func (t Type) String() string {
	return string(t)
}

func (t *Type) Set(s string) error {
	switch Type(s) {
	case SwiftStorageType:
		*t = SwiftStorageType
	default:
		return UndefinedStorageTypeErr
	}

	return nil
}

func (t Type) Type() string {
	return "storage"
}

const (
	SwiftStorageType Type = "swift"
)

var (
	AvailableStorageType    = []Type{SwiftStorageType}
	UndefinedStorageTypeErr = errors.New("undefined storage type")
)

func StringAvailableStorages() string {
	storageTypeStrings := make([]string, 0, len(AvailableStorageType))

	for _, storage := range AvailableStorageType {
		storageTypeStrings = append(storageTypeStrings, storage.String())
	}

	return strings.Join(storageTypeStrings, ", ")
}

type Storager interface {
	Authenticate(ctx context.Context) error
	Write(ctx context.Context, content io.Reader, params WriteParams) error
}

type WriteParams interface {
	SetName(name string)
}
