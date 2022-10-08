package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OdyseeTeam/gody-cdn/cleanup"
	"github.com/OdyseeTeam/gody-cdn/configs"
	"github.com/OdyseeTeam/gody-cdn/store"
	"github.com/OdyseeTeam/mirage/config"
	"github.com/OdyseeTeam/mirage/metadata"
	"github.com/OdyseeTeam/mirage/optimizer"
	http "github.com/OdyseeTeam/mirage/server"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/stop"
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
		stopper := stop.New()
		config.InitializeConfiguration()

		ds, err := store.NewDiskStore(viper.GetString("disk_cache.path"), 2)
		if err != nil {
			logrus.Fatal(errors.FullTrace(err))
		}
		localDsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", viper.GetString("local_db.user"), viper.GetString("local_db.password"), viper.GetString("local_db.host"), viper.GetString("local_db.database"))
		dbs := store.NewDBBackedStore(ds, localDsn)
		cacheParams := configs.ObjectCacheParams{
			Path: viper.GetString("disk_cache.path"),
			Size: viper.GetString("disk_cache.size"),
		}
		go cleanup.SelfCleanup(dbs, dbs, stopper, cacheParams, 30*time.Second)
		metadataDsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", viper.GetString("metadata_db.user"), viper.GetString("metadata_db.password"), viper.GetString("metadata_db.host"), viper.GetString("metadata_db.database"))
		metadataManager, err := metadata.Init(metadataDsn)
		if err != nil {
			logrus.Fatal(errors.FullTrace(err))
		}
		httpServer := http.NewServer(optimizer.NewOptimizer(), dbs, metadataManager)
		err = httpServer.Start(":6456")
		if err != nil {
			logrus.Fatal(err)
		}
		defer httpServer.Shutdown()

		interruptChan := make(chan os.Signal, 1)
		signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
		<-interruptChan
		stopper.StopAndWait()
	},
}
