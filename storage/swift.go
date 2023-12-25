package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/ncw/swift/v2"
)

const (
	envSwiftStorageUsername    = "SWIFT_STORAGE_USERNAME"
	envSwiftStorageApiKey      = "SWIFT_STORAGE_API_KEY"
	envSwiftStorageAuthUrl     = "SWIFT_STORAGE_AUTH_URL"
	envSwiftStorageRegion      = "SWIFT_STORAGE_REGION"
	envSwiftStorageUserAgent   = "SWIFT_STORAGE_USER_AGENT"
	envSwiftStorageAuthVersion = "SWIFT_STORAGE_AUTH_VERSION"
	envSwiftStorageDomain      = "SWIFT_STORAGE_DOMAIN"
	envSwiftStorageTenant      = "SWIFT_STORAGE_TENANT"
)

const (
	swiftMinSupportedAuthVersion = 1
	swiftMaxSupportedAuthVersion = 3
)

type SwiftStorageConfig struct {
	UserName    string `yaml:"user-name"`
	ApiKey      string `yaml:"api-key"`
	AuthUrl     string `yaml:"auth-url"`
	Region      string `yaml:"region"`
	UserAgent   string `yaml:"user-agent"`
	AuthVersion int    `yaml:"auth-version"`
	Domain      string `yaml:"domain"`
	Tenant      string `yaml:"tenant"`
}

func (s *SwiftStorageConfig) ReadFromEnviron() error {
	s.UserName = os.Getenv(envSwiftStorageUsername)
	s.ApiKey = os.Getenv(envSwiftStorageApiKey)
	s.AuthUrl = os.Getenv(envSwiftStorageAuthUrl)
	s.Region = os.Getenv(envSwiftStorageRegion)
	s.UserAgent = os.Getenv(envSwiftStorageUserAgent)
	s.Domain = os.Getenv(envSwiftStorageDomain)
	s.Tenant = os.Getenv(envSwiftStorageTenant)
	authVersion, err := s.readAuthVersion()
	if err != nil {
		return fmt.Errorf("read auth version: %w", err)
	}

	s.AuthVersion = authVersion

	return nil
}

func (s *SwiftStorageConfig) readAuthVersion() (int, error) {
	authVersionStr := os.Getenv(envSwiftStorageAuthVersion)
	authVersion, err := strconv.Atoi(authVersionStr)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", envSwiftStorageAuthVersion)
	}
	if authVersion < swiftMinSupportedAuthVersion || authVersion > swiftMaxSupportedAuthVersion {
		return 0, fmt.Errorf("%s is not within the supported range", envSwiftStorageAuthVersion)
	}
	return authVersion, nil
}

type SwiftWriteParams struct {
	Container   string `yaml:"container"`
	ObjectName  string `yaml:"-"`
	Hash        string `yaml:"-"`
	ContentType string `yaml:"-"`
}

func (s *SwiftWriteParams) SetName(name string) {
	s.ObjectName = name
}

type SwiftStorage struct {
	conn *swift.Connection
}

func NewSwiftStorage(config *SwiftStorageConfig) *SwiftStorage {
	conn := &swift.Connection{
		UserName:    config.UserName,
		AuthUrl:     config.AuthUrl,
		ApiKey:      config.ApiKey,
		Region:      config.Region,
		AuthVersion: config.AuthVersion,
		UserAgent:   config.UserAgent,
		Tenant:      config.Tenant,
		Domain:      config.Domain,
	}

	return &SwiftStorage{conn: conn}
}

func (s *SwiftStorage) Authenticate(ctx context.Context) error {
	return s.conn.Authenticate(ctx)
}

func (s *SwiftStorage) Write(ctx context.Context, content io.Reader, params WriteParams) error {
	swiftParams, ok := params.(*SwiftWriteParams)
	if !ok {
		return errors.New("params is not of type *SwiftWriteParams")
	}

	err := s.putObject(ctx, content, swiftParams)
	if err != nil {
		return fmt.Errorf("write to swift storage: %v", err)
	}

	return nil
}

func (s *SwiftStorage) putObject(ctx context.Context, content io.Reader, swiftParams *SwiftWriteParams) error {
	checkHash := swiftParams.Hash != ""

	_, err := s.conn.ObjectPut(
		ctx,
		swiftParams.Container,
		swiftParams.ObjectName,
		content,
		checkHash,
		swiftParams.Hash,
		swiftParams.ContentType,
		nil,
	)

	return err
}
