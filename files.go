package main

import (
	"encoding/json"
	"github.com/spf13/afero"
)

const (
	defaultAddress        = ":8080"
	defaultHTTPMethod     = "GET"
	defaultLatency        = "0ms"
	defaultHTTPStatusCode = 200
)

type ConfigFile struct {
	Addr      string     `json:"address"`
	Responses []Response `json:"responses"`
}

type Response struct {
	Status   int    `json:"status"`
	Method   string `json:"method"`
	Path     string `json:"path"`
	Latency  string `json:"latency"`
	JsonBody string `json:"json_body"`
}

func (c *ConfigFile) enforceDefaults() {
	if c.Addr == "" {
		c.Addr = defaultAddress
	}
	for i := range c.Responses {
		if c.Responses[i].Method == "" {
			c.Responses[i].Method = defaultHTTPMethod
		}
		if c.Responses[i].Latency == "" {
			c.Responses[i].Latency = defaultLatency
		}
		if c.Responses[i].Status == 0 {
			c.Responses[i].Status = defaultHTTPStatusCode
		}
	}
}

func loadConfig(fs afero.Fs) (ConfigFile, error) {
	file, err := afero.ReadFile(fs, DefaultConfigurationFileName)
	if err != nil {
		return ConfigFile{}, err
	}

	configFile := ConfigFile{}
	err = json.Unmarshal(file, &configFile)
	if err != nil {
		return ConfigFile{}, err
	}
	configFile.enforceDefaults()
	return configFile, nil
}

func isFirstRun(fs afero.Fs) (bool, error) {
	exists, err := afero.Exists(fs, DefaultConfigurationFileName)
	if err != nil {
		return false, err
	}
	return !exists, nil
}

func configureFirstRun(fs afero.Fs) error {
	err := afero.WriteFile(fs, DefaultConfigurationFileName, []byte(DefaultConfigurationFileContent), 0644)
	if err != nil {
		return err
	}
	err = afero.WriteFile(fs, CustomBodyExampleFileName, []byte(CustomBodyExampleFileContent), 0644)
	if err != nil {
		return err
	}
	return nil
}
