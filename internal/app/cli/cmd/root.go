package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	// homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// define flags and handle configuration
func init() {

	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(honkCmd)
	rootCmd.AddCommand(confCmd)
	rootCmd.AddCommand(dumpCmd)
	rootCmd.AddCommand(authCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")

	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))

	viper.SetDefault("author", "Adam Smith <developer@gh.arusty.dev>")
	viper.SetDefault("license", "agpl3")
}

func initConfig() {
	viper.SetConfigName("brood") // Config file name without extension
	viper.SetConfigType("yaml")  // Config file type
	viper.AddConfigPath(".")     // Look for the config file in the current directory

	/*
	   AutomaticEnv will check for an environment variable any time a viper.Get request is made.
	   It will apply the following rules.
	       It will check for an environment variable with a name matching the key uppercased and prefixed with the EnvPrefix if set.
	*/
	viper.AutomaticEnv()
	viper.SetEnvPrefix("env")                              // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // this is useful e.g. want to use . in Get() calls, but environmental variables to use _ delimiters (e.g. app.port -> APP_PORT)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Set up environment variable mappings if necessary
	/*
	   BindEnv takes one or more parameters. The first parameter is the key name, the rest are the name of the environment variables to bind to this key.
	   If more than one are provided, they will take precedence in the specified order. The name of the environment variable is case sensitive.
	   If the ENV variable name is not provided, then Viper will automatically assume that the ENV variable matches the following format: prefix + "_" + the key name in ALL CAPS.
	   When you explicitly provide the ENV variable name (the second parameter), it does not automatically add the prefix.
	       For example if the second parameter is "id", Viper will look for the ENV variable "ID".
	*/
	viper.BindEnv("app.name", "APP_NAME") // Bind the app.name key to the APP_NAME environment variable

}

var rootCmd = &cobra.Command{
	Use:   "gosling",
	Short: "Goslings is a cloud-native MSFT Cloud IR tool",
	Long: `The Goslings Tool is a robust and flexible hunt and incident
			response tool that adds novel authentication and data gathering
			methods in order to run a full investigation against a customer’s
			Azure Active Directory (AzureAD), Azure, and M365 environments.
                Complete documentation is available at https://github.com/arustydev/goslings
                https://github.com/cisagov/untitledgoosetool`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
