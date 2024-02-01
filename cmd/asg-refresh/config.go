package main

import (
	"flag"
	"os"

	"github.com/natemarks/asg-refresh/version"
	"github.com/rs/zerolog"
)

// Config is the configuration for the application
type Config struct {
	update bool
	asgID  string
	debug  bool
}

// GetLogger returns a logger for the application
func (c Config) GetLogger() (log zerolog.Logger) {
	log = zerolog.New(os.Stdout).With().Str("version", version.Version).Timestamp().Logger()
	log.Level(zerolog.InfoLevel)
	if c.debug {
		log = log.Level(zerolog.DebugLevel)
	}
	return log
}

// GetConfig returns the configuration for the application
func GetConfig() (config Config, err error) {
	// Define flags
	filePtr := flag.String("file", "", "Path to the file")
	updatePtr := flag.Bool("update", false, "update the secret value if it exists")
	debugPtr := flag.Bool("debug", false, "Enable debug mode")

	// Parse command line arguments
	flag.Parse()
	config.asgID = *filePtr
	config.update = *updatePtr
	config.debug = *debugPtr

	return config, nil
}
