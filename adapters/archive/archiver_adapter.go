package archive

import (
	"context"
	"fmt"
	"github.com/mholt/archiver/v4"
	"io"
	"path/filepath"
)

type ArchiverAdapter struct {
	archiver archiver.Archival
}

func NewArchiverAdapter(archiver archiver.Archival) *ArchiverAdapter {
	return &ArchiverAdapter{archiver: archiver}
}

func (a *ArchiverAdapter) Format() string {
	archiverExtension := a.archiver.Name()
	if len(archiverExtension) > 1 {
		return archiverExtension[1:]
	}

	return archiverExtension
}

func (a *ArchiverAdapter) Archive(ctx context.Context, output io.Writer, files []string) error {
	archiverFiles, err := a.convertFilesToArchiveFiles(files)
	if err != nil {
		return fmt.Errorf("prepare file list: %w", err)
	}

	return a.archiver.Archive(ctx, output, archiverFiles)
}

func (a *ArchiverAdapter) convertFilesToArchiveFiles(files []string) ([]archiver.File, error) {
	filesMap := make(map[string]string, len(files))

	for _, file := range files {
		filesMap[file] = filepath.Base(file)
	}

	archiveFiles, err := archiver.FilesFromDisk(nil, filesMap)
	if err != nil {
		return nil, fmt.Errorf("convert source path to archive files: %w", err)
	}

	return archiveFiles, nil
}
