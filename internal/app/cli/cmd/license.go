package cmd

import (
	"fmt"
	"goslings/internal/utils"

	"github.com/spf13/cobra"
)

var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Print the license of goslings",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(utils.License)
	},
}
