package micasa

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/derWhity/micasa/internal/log"
	"github.com/derWhity/micasa/internal/models"
	"github.com/pkg/errors"
)

// ConfigService gives the authenticated user access to parts of the application's configuration
type ConfigService interface {
	Load() error
	// LoadFromFile loads the configuration from the given JSON file and returns it
	LoadFromFile(filename string) error
	// Write writes the current application configuration to the default file name
	Write() error
	// WriteToFile writes the current application configuration to a JSON file
	WriteToFile(filename string) error
	// GetConfig retuns the current application configuration
	GetConfig() models.Configuration
}

// -- ConfigService implementation -------------------------------------------------------------------------------------

// Simple index structure to speed up whitelist lookups
type whitelistIdx struct {
	sync.RWMutex
	data map[string]bool
}

type configService struct {
	configFilename string
	config         *models.Configuration
	logger         log.Logger
}

// NewConfigService creates a new configuration service instance with the given default file name
func NewConfigService(configFilename string, logger log.Logger) ConfigService {
	return &configService{
		configFilename: configFilename,
		logger:         logger,
	}
}

// Load loads the application config from its default file location
func (s *configService) Load() error {
	return s.LoadFromFile(s.configFilename)
}

// LoadFromFile loads the configuration from the given JSON file and returns it
func (s *configService) LoadFromFile(filename string) error {
	s.logger.Info("Loading configuration file", log.FldFile, filename)
	conf, err := models.GetDefaultConfig()
	if err != nil {
		return errors.Wrap(err, "LoadFromFile: Failed to create default config")
	}
	f, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "LoadFromFile: cannot load configuration file")
	}
	defer f.Close()
	if err = json.NewDecoder(f).Decode(&conf); err != nil {
		return errors.Wrap(err, "LoadFromFile: Failed to decode configuration file")
	}
	s.config = conf
	return nil
}

// Write writes the current application configuration to the default file name
func (s *configService) Write() error {
	return s.WriteToFile(s.configFilename)
}

// WriteToFile writes the current application configuration to a JSON file
func (s *configService) WriteToFile(filename string) error {
	s.logger.Info("Writing configuration file", log.FldFile, filename)
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "WriteToFile: Cannot open configuration file '%s' to write to", filename)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	conf := s.GetConfig()
	if err := enc.Encode(&conf); err != nil {
		return errors.Wrap(err, "WriteToFile: Failed to serialize configuration data")
	}
	return nil
}

// GetConfig retuns the current application configuration
func (s *configService) GetConfig() models.Configuration {
	var ret models.Configuration
	if s.config != nil {
		ret = *s.config
	} else {
		if tmp, err := models.GetDefaultConfig(); err == nil {
			ret = *tmp
		}
	}
	return ret
}
