package main

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	HTTP           HTTPConfig           `mapstructure:"http"`
	Authentication AuthenticationConfig `mapstructure:"authentication"`
}

type HTTPConfig struct {
	ListenAddr string `mapstructure:"listen_addr"`
}

type AuthenticationConfig struct {
	ClientId     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	Realm        string `mapstructure:"realm"`
	Hostname     string `mapstructure:"hostname"`
}

// Parse parses the configuration from the environment.

// Parse parses config from the environment, config file and flags.
func Parse() (*Config, error) {
	const configFileFlag = "config"

	{
		pflag.String(configFileFlag, "", "Config file")
		pflag.String("http.listen_addr", "", "HTTP server listen addr")
		pflag.String("authentication.client_id", "https", "Keycloak client id")
		pflag.String("authentication.client_secret", "", "Keycloak client secret")
		pflag.String("authentication.realm", "", "Keycloak realm")
		pflag.String("authentication.hostname", "", "Keycloak hostname")
		pflag.Parse()
		_ = viper.BindPFlags(pflag.CommandLine)
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	configFile := viper.GetString(configFileFlag)
	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	cfg := &Config{}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("viper unmarshall: %w", err)
	}

	return cfg, nil
}
