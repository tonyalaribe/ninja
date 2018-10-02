package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tonyalaribe/ninja/core"
	"github.com/tonyalaribe/ninja/uilayer"
)

var rootCmd = &cobra.Command{
	Use:   "ninja",
	Short: "Ninja is a dynamic api engine",
	Long:  `Ninja lets you build powerful api's(REST, graphql, grpc, etc) for your apps and web applications using a very simple interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager := core.New()
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

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Unable to read config:", err)
			os.Exit(1)
		}
	}
}
