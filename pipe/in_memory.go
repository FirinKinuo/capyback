package pipe

import "io"

// InMemoryPipe describes a pipe that operates in-memory
type InMemoryPipe struct {
	writer *io.PipeWriter // writer represents the writer end of the pipe
	reader *io.PipeReader // reader represents the reader end of the pipe
}

// Write writes to the pipe's writer end
func (mp *InMemoryPipe) Write(p []byte) (n int, err error) {
	return mp.writer.Write(p)
}

// Read reads from the pipe's reader end
func (mp *InMemoryPipe) Read(p []byte) (n int, err error) {
	return mp.reader.Read(p)
}

// CloseWrite closes the pipe's writer end
func (mp *InMemoryPipe) CloseWrite() {
	_ = mp.writer.Close()
}

// CloseRead closes the pipe's reader end
func (mp *InMemoryPipe) CloseRead() {
	_ = mp.reader.Close()
}

// CloseWriteWithErr closes the pipe's writer end and associates an error with it
func (mp *InMemoryPipe) CloseWriteWithErr(err error) {
	_ = mp.writer.CloseWithError(err)
}

// CloseReadWithErr closes the pipe's reader end and associates an error with it
func (mp *InMemoryPipe) CloseReadWithErr(err error) {
	_ = mp.reader.CloseWithError(err)
}

// NewInMemoryPipe initializes and returns a new in-memory pipe
func NewInMemoryPipe() *InMemoryPipe {
	pipeReader, pipeWriter := io.Pipe()

	return &InMemoryPipe{
		writer: pipeWriter,
		reader: pipeReader,
	}
}
