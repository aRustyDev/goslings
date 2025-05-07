package conf

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("brood")              // Config file name without extension
	viper.SetConfigType("yaml")               // Config file type
	viper.AddConfigPath(".")                  // 1st: Look for the config file in the current directory
	viper.AddConfigPath("./configs")          // 2nd: Look for the config file in the local configs directory
	viper.AddConfigPath("~/.config/goslings") // 3rd: Look for the config file in the home config directory

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

	viper.SetDefault("author", "Adam Smith <developer@gh.arusty.dev>")
	viper.SetDefault("license", "agpl3")

}
