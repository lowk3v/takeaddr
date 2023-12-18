package config

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var Config Yaml

//go:embed config.yml
var configYml string

type Yaml struct {
	// Your config load from config.yml here
}

func init() {
	// Load Config yml
	err := yaml.Unmarshal([]byte(configYml), &Config)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error loading config: %s", err)
		os.Exit(1)
	}
}

func CustomConfig(cfgPath string) error {
	// Open config file
	file, err := os.Open(cfgPath)
	if err != nil {
		return err
	}

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&Config); err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}
