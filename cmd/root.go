package cmd

import (
	"fmt"
	"os"

	"github.com/frsfahd/go-proxy/api"
	"github.com/frsfahd/go-proxy/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configs config.Config
	port    int
	origin  string
	cfgFile string
)
var version = "0.0.1"
var rootCmd = &cobra.Command{
	Use:     "go-proxy",
	Short:   "simple cache proxy",
	Version: version,
	Long: `A Cache forward-proxy for any server, for documentation
please take a look at https://github.com/frsfahd/go-proxy`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here

		if port != 0 && origin != "" {
			api.RunAPI(port, origin, configs)
		} else {
			fmt.Println(`⚠️ please insert PORT number and ORIGIN address !!!`)
			os.Exit(1)
		}

	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 0, "local port for the proxy")
	rootCmd.PersistentFlags().StringVarP(&origin, "origin", "o", "", "target server")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "yaml config file for redis connection. if not supplied it will try to look up at config.yaml in current dir")
	rootCmd.MarkPersistentFlagRequired("port")
	rootCmd.MarkPersistentFlagRequired("origin")

}
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find config in current directory.

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("config.yaml not found...")
		} else {
			fmt.Printf("cannot read config file: %s\n", viper.ConfigFileUsed())
		}
		os.Exit(1)
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())

	err = viper.Unmarshal(&configs)
	if err != nil {
		fmt.Printf("unable to decode config to struct: %v\n", err)
		os.Exit(1)
	}

}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
