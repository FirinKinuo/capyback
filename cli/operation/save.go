package operation

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/FirinKinuo/capyback/application"
	"github.com/FirinKinuo/capyback/archive"
	"github.com/FirinKinuo/capyback/cli/flag"
	"github.com/FirinKinuo/capyback/config"
	"github.com/FirinKinuo/capyback/pipe"
	"github.com/FirinKinuo/capyback/storage"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ErrorMultipleFilesWithoutName is an error when handling more than one file or directory without backup name.
var ErrorMultipleFilesWithoutName = errors.New(
	"backup name is required when handling more than one file or directory. " +
		"See the help section (enter --help) for further information",
)

// ErrorNoResourcesToBackup is an error when no resources were specified for backup.
var ErrorNoResourcesToBackup = errors.New("resources for backup were not specified")

// Save is a command for save new backup to storage.
type Save struct {
	command   *cobra.Command
	appConfig *config.Config

	resources  []string
	backupName string

	storageFlagSet *flag.StorageFlagSet
	configFlagSet  *flag.ConfigFlagSet
	archiveFlagSet *flag.ArchiveFlagSet

	storager storage.Storager
	archiver archive.Archiver
}

// NewSave creates a new Save.
func NewSave(defaultConfigPath string) *Save {
	save := &Save{
		storageFlagSet: flag.NewStorageFlagSet(),
		configFlagSet:  flag.NewConfigFlagSet(defaultConfigPath),
		archiveFlagSet: flag.NewArchiveFlagSet(archive.DefaultFormat),
	}

	command := &cobra.Command{
		Use:   "save [FILE/DIR...]",
		Short: "Save new backup",
		Args:  cobra.MinimumNArgs(1),
		Run:   save.run,
	}

	command.PersistentFlags().AddFlagSet(save.FlagSet())

	save.command = command

	return save
}

// FlagSet returns a flag set for save command.
func (s *Save) FlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("save", pflag.PanicOnError)

	flagSet.StringVarP(
		&s.backupName,
		"name",
		"o",
		"",
		"backup name, example: \"my-backup@01.02.2006.tar.zst\". Required when handling more than one file or directory.",
	)

	flagSet.AddFlagSet(s.storageFlagSet.FlagSet())
	flagSet.AddFlagSet(s.configFlagSet.FlagSet())
	flagSet.AddFlagSet(s.archiveFlagSet.FlagSet())

	return flagSet
}

func (s *Save) Command() *cobra.Command {
	return s.command
}

func (s *Save) validateBackupName() error {
	if len(s.resources) > 1 && s.backupName == "" {
		return ErrorMultipleFilesWithoutName
	}

	return nil
}

func (s *Save) configureBackupName() error {
	// If there is only one resource for backup and the name was not set with a flag
	// Then we use the name of the resource itself as the backup name
	if len(s.resources) == 1 && s.backupName == "" {
		s.backupName = s.resources[0]
	}

	err := s.validateBackupName()
	if err != nil {
		return fmt.Errorf("validate backupName: %w", err)
	}

	s.backupName = fmt.Sprintf("%s.%s", s.backupName, s.archiveFlagSet.Format)

	return nil
}

// configure configures the save command from flag sets.
func (s *Save) configure(args []string) error {
	s.resources = args

	if len(s.resources) < 1 {
		return ErrorNoResourcesToBackup
	}

	err := s.configureBackupName()
	if err != nil {
		return fmt.Errorf("configure backupName: %w", err)
	}

	if s.configFlagSet.Path != "" {
		s.appConfig, err = s.configFlagSet.ReadYamlConfig()
		if err != nil {
			return fmt.Errorf("read yaml config: %w", err)
		}
	}

	s.archiver, err = archive.IdentifyArchiver(s.backupName)
	if err != nil {
		return fmt.Errorf("identify archiver: %w", err)
	}

	backupStorage, err := s.appConfig.Storage.ReadStorage()
	if err != nil {
		return fmt.Errorf("read storager: %w", err)
	}

	s.storager = backupStorage

	return nil
}

func (s *Save) performBackup(ctx context.Context) error {
	inMemoryPipe, err := pipe.NewPipe(pipe.InMemoryPipeType)
	if err != nil {
		log.Fatal("create new pipe", "err", err)
	}
	defer func() {
		inMemoryPipe.CloseWrite()
		inMemoryPipe.CloseRead()
	}()

	backup := application.NewBackup(inMemoryPipe, s.storager, s.archiver)

	writeParams, err := s.appConfig.Storage.ReadWriteParams()
	if err != nil {
		return fmt.Errorf("read write params: %w", err)
	}
	writeParams.SetName(s.backupName)

	err = backup.Save(ctx, s.resources, writeParams)
	if err != nil {
		return fmt.Errorf("backup save: %w", err)
	}

	return nil
}

func (s *Save) run(_ *cobra.Command, args []string) {
	err := s.configure(args)
	if err != nil {
		log.Fatal("configure", "err", err)
	}

	log.Infof("Create new backup: %s", s.backupName)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err = s.performBackup(ctx)
	if err != nil {
		select {
		case <-ctx.Done():
			log.Info("Backup cancelled")

		default:
			log.Fatal("perform backup", "err", err)
		}
	}
}
