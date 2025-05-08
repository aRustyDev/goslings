package conf

import (
	"strings"

	"goslings/internal/auth/shared"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("brood")              // Config file name without extension
	viper.SetConfigType("yaml")               // Config file type
	viper.AddConfigPath(".")                  // 1st: Look for the config file in the current directory
	viper.AddConfigPath("./configs")          // 2nd: Look for the config file in the local configs directory
	viper.AddConfigPath("~/.config/goslings") // 3rd: Look for the config file in the home config directory
	viper.SetConfigFile("./configs/brood.yaml")
	/*
	   AutomaticEnv will check for an environment variable any time a viper.Get request is made.
	   It will apply the following rules.
	       It will check for an environment variable with a name matching the key uppercased and prefixed with the EnvPrefix if set.
	*/
	viper.AutomaticEnv()
	viper.SetEnvPrefix("env")                              // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // this is useful e.g. want to use . in Get() calls, but environmental variables to use _ delimiters (e.g. app.port -> APP_PORT)

	log.Infof("Reading config file %s", viper.GetViper().ConfigFileUsed())
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
	// Bind the nested keys to the DISTINCT environment variables
	viper.BindEnv("auth.simple.user", "GOSLING_USER")
	viper.BindEnv("auth.simple.pass", "GOSLING_PASS")
	viper.BindEnv("auth.app.id", "GOSLING_APP_ID")
	viper.BindEnv("auth.app.secret", "GOSLING_APP_SECRET")
	viper.BindEnv("msft.tenant", "GOSLING_TENANT")
	viper.BindEnv("msft.subscription", "GOSLING_SUBSCRIPTION")
	viper.BindEnv("msft.usgov.cloud", "GOSLING_USGOV_CLOUD")
	viper.BindEnv("msft.usgov.exo", "GOSLING_USGOV_EXO")
	viper.BindEnv("msft.m365auth", "GOSLING_M365_AUTH")
	viper.BindEnv("msft.msgtrace", "GOSLING_EXO_MSG_TRACE")

	viper.SetDefault("author", "Adam Smith <developer@gh.arusty.dev>")
	viper.SetDefault("license", "agpl3")

}

// getConfigFromViper simulates getting config from a viper-based config package
func GetAuthConfig() *shared.AuthParams {
	log.Info("Extracting configs to AuthParams")
	ap := &shared.AuthParams{
		Username:            viper.GetString("auth.simple.user"),
		Password:            viper.GetString("auth.simple.pass"),
		TenantID:            viper.GetString("msft.tenant"),
		ClientID:            viper.GetString("auth.app.id"),
		ClientSecret:        viper.GetString("auth.app.secret"),
		SubscriptionID:      viper.GetString("msft.subscription"),
		UsGovernment:        viper.GetBool("msft.usgov.cloud"),
		ExoUSGovernment:     viper.GetBool("msft.usgov.exo"),
		M365Enabled:         viper.GetBool("msft.m365auth"),
		MessageTraceEnabled: viper.GetBool("msft.msgtrace"),
	}
	return ap
}
