package pipe

import (
	"errors"
	"io"
	"strings"
)

// Type defines the type of the pipe
type Type string

var (
	// UndefinedPipeErr is the error that is returned when the pipe type is not defined
	UndefinedPipeErr = errors.New("undefined pipe type")
)

const (
	// InMemoryPipeType is a pipe type that uses in-memory storage
	InMemoryPipeType Type = "in-memory"
)

// String method returns the string representation of the Type
func (t Type) String() string {
	return string(t)
}

// MarshalText method converts the Type to a []byte
func (t Type) MarshalText() (text []byte, err error) {
	return []byte(t), nil
}

// UnmarshalText method converts a []byte to a Type
func (t *Type) UnmarshalText(text []byte) error {
	s := strings.ToLower(string(text))

	switch s {
	case InMemoryPipeType.String():
		*t = InMemoryPipeType
	default:
		// return an error if the provided text does not match a known Type
		return UndefinedPipeErr
	}

	return nil
}

// Piper is an interface for pipes that includes methods for ReadWrite Closer with error handling
type Piper interface {
	io.ReadWriter
	CloseWrite()
	CloseWriteWithErr(err error)
	CloseRead()
	CloseReadWithErr(err error)
}

// NewPipe function returns a new Piper of the provided pipeType, or an error if the pipeType is not defined
func NewPipe(pipeType Type) (Piper, error) {
	switch pipeType {
	case InMemoryPipeType:
		return NewInMemoryPipe(), nil
	default:
		return nil, UndefinedPipeErr
	}
}
