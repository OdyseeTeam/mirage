package cmd

import (
	"github.com/OdyseeTeam/mirage/internal/version"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Mirage",
	Long:  `All software has versions. This is Mirage's`,
	Run: func(cmd *cobra.Command, args []string) {
		println(version.FullName())
	},
}
