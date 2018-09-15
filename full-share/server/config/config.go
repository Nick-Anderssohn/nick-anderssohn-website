package config

import (
	"encoding/json"
	"io/ioutil"
)

// StructConfig stores the global config values for the application.
type StructConfig struct {
	RunReaper bool `json:"RunReaper"`
	TestMode  bool `json:"TestMode"`
}

// ApplicationConfig stores the global config values for the application.
var ApplicationConfig StructConfig

func init() {
	readConfig()
}

func readConfig() {
	configBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(configBytes, &ApplicationConfig)
}
