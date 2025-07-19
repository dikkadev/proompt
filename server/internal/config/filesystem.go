package config

import (
	"github.com/spf13/afero"
)

// GetFilesystem returns the configured filesystem
// For now returns OS filesystem, but will support other backends later
func (c *Config) GetFilesystem() afero.Fs {
	// TODO: Support different filesystem backends based on config
	// For now, always return OS filesystem
	return afero.NewOsFs()
}
