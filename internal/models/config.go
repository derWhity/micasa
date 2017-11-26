package models

import (
	"path"

	"github.com/kardianos/osext"
)

// Configuration is the application's main configuration structure
type Configuration struct {
	// The directory where MiCasa stores all of its data - defaults to the ./data subdirectory of the folder, the
	// executable resides in
	DataDir string `json:"dataDir"`
	// The IP address to listen at - including the port number
	ListenAddress string `json:"listenAddress"`
}

// GetDefaultConfig returns the default configuration values for the application
func GetDefaultConfig() (*Configuration, error) {
	execDir, err := osext.ExecutableFolder()
	if err != nil {
		return nil, err
	}
	return &Configuration{
		DataDir:       path.Join(execDir, "data"),
		ListenAddress: ":3000",
	}, nil
}
