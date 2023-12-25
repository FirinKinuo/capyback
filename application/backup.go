package application

import (
	"context"
	"fmt"
	"github.com/FirinKinuo/capyback/archive"
	"github.com/FirinKinuo/capyback/pipe"
	"github.com/FirinKinuo/capyback/storage"
	"github.com/charmbracelet/log"
)

// Backup is the application that creates a backup of the files and writes it to the storage.
type Backup struct {
	pipe     pipe.Piper
	storage  storage.Storager
	archiver archive.Archiver
}

// NewBackup constructs a new Backup application.
func NewBackup(p pipe.Piper, s storage.Storager, a archive.Archiver) *Backup {
	return &Backup{
		pipe:     p,
		storage:  s,
		archiver: a,
	}
}

// Save creates a backup of the files and writes it to the storage.
func (t *Backup) Save(ctx context.Context, files []string, writeParams storage.WriteParams) error {
	log.Info("Archiving", "format", t.archiver.Format())

	go t.archive(ctx, files)

	log.Info("Attempting to authenticate to storage")
	err := t.storage.Authenticate(ctx)
	if err != nil {
		return fmt.Errorf("authenticate storage: %w", err)
	}

	log.Info("Authentication to storage succeeded.")

	log.Info("Writing to storage")
	err = t.storage.Write(ctx, t.pipe, writeParams)
	if err != nil {
		return fmt.Errorf("write to storage: %w", err)
	}

	log.Info("Writing to storage completed successfully")
	return nil
}

// archive creates an archive of the files and writes it to the pipe.
func (t *Backup) archive(ctx context.Context, files []string) {
	err := t.archiver.Archive(ctx, t.pipe, files)
	if err != nil {
		t.pipe.CloseWriteWithErr(fmt.Errorf("archive: %w", err))
	}

	t.pipe.CloseWrite()
}
