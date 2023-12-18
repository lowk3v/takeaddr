package config

import (
	_ "embed"
	"github.com/lowk3v/micro-tool-template/pkg/log"
	"gopkg.in/yaml.v3"
	"os"
)

var Config Yaml

//go:embed config.yml
var configYml string

var Log log.Logger

type Yaml struct {
	// Your config load from config.yml here
}

func init() {
	Log = *log.New("info")

	// Load Config yml
	err := yaml.Unmarshal([]byte(configYml), &Config)
	if err != nil {
		Log.Errorf("Error loading config: %s", err)
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
