package config

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/caarlos0/env"
)

// Flags ...
type Flags struct {
	Dir      string `env:"AUTOBOT_DIR"`
	Server   string `env:"AUTOBOT_SERVER"`
	Help     bool
	Validate bool
	Verbose  bool
	Version  bool
}

// Config ...
type Config struct {
	// Verbose toggles the verbosity
	Verbose bool
	// LogLevel is the level with with to log for this config
	LogLevel string `mapstructure:"log_level"`
	// LogFormat is the format that is used for logging
	LogFormat string `mapstructure:"log_format"`
	// ReloadSignal ...
	ReloadSignal syscall.Signal
	// TermSignal ...
	TermSignal syscall.Signal
	// KillSignal ...
	KillSignal syscall.Signal
	// File...
	File string
	// FileMode ...
	FileMode os.FileMode
	// Flags ...
	Flags Flags
	// Stdin ...
	Stdin *os.File
	// Stdout ...
	Stdout *os.File
	// Stderr ...
	Stderr *os.File
}

// New ...
func New() *Config {
	return &Config{
		File:         ".autobot.yml",
		KillSignal:   syscall.SIGINT,
		LogFormat:    "text",
		LogLevel:     "warn",
		ReloadSignal: syscall.SIGHUP,
		TermSignal:   syscall.SIGTERM,
		Stdin:        os.Stdin,
		Stdout:       os.Stdout,
		Stderr:       os.Stderr,
	}
}

// InitDefaultConfig() ...
func (c *Config) InitDefaultConfig() error {
	cwd, err := c.Cwd()
	if err != nil {
		return err
	}
	c.File = filepath.Join(cwd, c.File)

	err = env.Parse(&c.Flags)
	if err != nil {
		return err
	}

	return nil
}

// Cwd ...
func (c *Config) Cwd() (string, error) {
	return os.Getwd()
}
