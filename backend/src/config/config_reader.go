package config

import "gopkg.in/ini.v1"

type AppConfig struct {
	Config *Config
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}

func (a *AppConfig) Load(path string) error {
	cfgFile, err := ini.Load(path)
	if err != nil {
		return err
	}

	config := &Config{}

	err = cfgFile.MapTo(config)
	if err != nil {
		return err
	}

	a.Config = config

	return nil
}
