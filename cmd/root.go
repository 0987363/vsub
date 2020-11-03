package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var RootCmd = &cobra.Command{
	Use:   "vsub",
	Short: "v2 subscribe backend server",
}

var configFilePath string

func init() {
	RootCmd.PersistentFlags().StringVarP(
		&configFilePath, "config", "c", "", "Path to the config file",
	)
}

func LoadConfiguration(cmd *cobra.Command, args []string) {
	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
	} else {
		viper.SetConfigName("vsub")
		viper.AddConfigPath("/etc/heifeng")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		if viper.ConfigFileUsed() == "" {
			log.Fatalf("Unable to find configuration file.")
		}
		log.Fatalf("Failed to load %s: %v", viper.ConfigFileUsed(), err)
	} else {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}
}
