package archive

import (
	"context"
	"io"
)

type Archiver interface {
	Format() string
	Archive(ctx context.Context, out io.Writer, files []string) error
}
