package cmd

import (
	"context"
	"fmt"
	"os"

	"goslings/internal/about"
	"goslings/internal/conf"

	"github.com/spf13/cobra"
)

// define flags and handle configuration
func init() {

	cobra.OnInitialize(conf.InitConfig)

	rootCmd.AddCommand(honkCmd)
	rootCmd.AddCommand(confCmd)
	rootCmd.AddCommand(dumpCmd)
	rootCmd.AddCommand(authCmd)

	rootCmd.PersistentFlags().StringP("author", "a", "Adam Smith", "Author name for copyright attribution")
	// rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
}

var rootCmd = &cobra.Command{
	Use:     "gosling",
	Short:   "Goslings is a cloud-native MSFT Cloud IR tool",
	Version: about.Version,
	Long: `The Goslings Tool is a robust and flexible hunt and incident
			response tool that adds novel authentication and data gathering
			methods in order to run a full investigation against a customer’s
			Azure Active Directory (AzureAD), Azure, and M365 environments.
                Complete documentation is available at https://github.com/arustydev/goslings
                https://github.com/cisagov/untitledgoosetool`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		// TODO: Output subcommand options
	},
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
