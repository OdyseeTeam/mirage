package cmd

import (
	"fmt"
	"os"

	"github.com/OdyseeTeam/mirage/internal/version"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.PersistentFlags().BoolP("debugmode", "d", false, "turns on debug mode for the application command.")
	rootCmd.PersistentFlags().BoolP("tracemode", "t", false, "turns on trace mode for the application command, very verbose logging.")
	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		logrus.Panic(err)
	}
}

var rootCmd = &cobra.Command{
	Use:     "mirage",
	Short:   "Odysee image processing server",
	Version: version.FullName(),
	Long:    `compressor/caching/distribution/proxy server for thumbnails/static content`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute executes the root command and is the entry point of the application from main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
