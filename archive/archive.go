package archive

import (
	"fmt"
	archiveAdapter "github.com/FirinKinuo/capyback/adapters/archive"
	"github.com/mholt/archiver/v4"
)

const DefaultFormat = "tar.zst"

// IdentifyArchiver is a function to identify the archiving method of a file.
func IdentifyArchiver(file string) (Archiver, error) {
	format, _, err := archiver.Identify(file, nil)
	if err != nil {
		return nil, fmt.Errorf("identify: %w", err)
	}

	// If we successfully identified the format, adapt it using NewArchiverAdapter
	// and return the related Archiver
	return archiveAdapter.NewArchiverAdapter(format.(archiver.Archival)), nil
}
