package cmd

import (
	"github.com/OdyseeTeam/mirage/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var compressionEnabled bool

func init() {
	serveCmd.PersistentFlags().BoolVar(&compressionEnabled, "enable_compression", true, "should optimizations of images be performed")
	//Bind to Viper
	err := viper.BindPFlags(serveCmd.PersistentFlags())
	if err != nil {
		logrus.Panic(err)
	}
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs the Mirage server",
	Long:  `Runs the Mirage server`,
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config.InitializeConfiguration()
		//server.Start()
	},
}
