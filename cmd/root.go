package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tonyalaribe/ninja/core"
	"github.com/tonyalaribe/ninja/datalayer"
	_ "github.com/tonyalaribe/ninja/datalayer/mongodb"
	"github.com/tonyalaribe/ninja/uilayer"
)

var config Config
var rootCmd = &cobra.Command{
	Use:   "ninja",
	Short: "Ninja is a dynamic api engine",
	Long:  `Ninja lets you build powerful api's(REST, graphql, grpc, etc) for your apps and web applications using a very simple interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		datastore, err := datalayer.Connect(config.DBConfig.DriverType, config.DBConfig)
		if err != nil {
			log.Fatalf("unable to initialize datalayer with error: `%v`", err)
			return
		}

		manager, err := core.New(core.UseDataStore(datastore))
		if err != nil {
			log.Fatalf("unable to initialize core with error: `%v`", err)
			return
		}
		uilayer.Register(manager)
	},
}

func Execute() {
	var cfgFile string
	cobra.OnInitialize(initConfig(cfgFile))
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ninja.yaml)")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Config struct {
	IsProduction bool               `mapstructure:"is_production"`
	ShortName    string             `mapstructure:"short_name"`
	LongName     string             `mapstructure:"long_name"`
	DBConfig     datalayer.DBConfig `mapstructure:"db_config"`
}

func initConfig(cfgFile string) func() {
	return func() {
		// Don't forget to read config either from cfgFile or from home directory!
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Search config in home and current directory with name ".ninja" (without extension).
			viper.AddConfigPath(home)
			viper.AddConfigPath(".")
			viper.SetConfigName(".ninja")
		}

		var err error
		if err = viper.ReadInConfig(); err != nil {
			fmt.Println("Unable to read config:", err)
			os.Exit(1)
		}

		err = viper.Unmarshal(&config)
		if err != nil {
			log.Panicf("unable to decode into struct, %v", err)
		}

		if !config.IsProduction {
			log.Println("In Development Mode. Logging configuration data:")
			indentedConfig, _ := json.MarshalIndent(config, "", "\t")
			log.Printf("\n%s\n\n", indentedConfig)
		}
	}
}
