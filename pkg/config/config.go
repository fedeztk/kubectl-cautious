package config

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

// holds info about the configuration of the application under HOME/.kube/cautious.yaml/yml
// contexts:
//   - name: prod # This should match a context in your `KUBECONFIG` file
//     action:
//       - name: apply
//       - name: delete
//         dry-run: true

// Config holds the Contexts for the plugin
type Config struct {
	Contexts []Context
}

// Context holds the configuration for a specific context
type Context struct {
	Name    string
	Actions []Action
}

// Action represents a kube action and whether it should be a dry-run
type Action struct {
	Name   string
	DryRun bool `mapstructure:"dry-run" yaml:"dry-run"`
}

// ReadConfig reads the cautious.yaml
func ReadConfig() (*Config, error) {
	viper.SetConfigName("cautious")                                // name of config file (without extension)
	viper.SetConfigType("yaml")                                    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".kube")) // call multiple times to add many search paths
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// generate a default config file and write it to disk
			viper.Set("contexts", []Context{
				{
					Name: "prod",
					Actions: []Action{
						{
							Name:   "apply",
							DryRun: false,
						},
						{
							Name:   "delete",
							DryRun: true,
						},
					},
				},
			})

			err := viper.SafeWriteConfig()
			if err != nil {
				return nil, err
			}

			log.Info("kubectl-cautious config file generated under .kube directory")
		}
	}

	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
