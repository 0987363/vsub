package cmd

import (
	"github.com/0987363/vsub/handlers"
	"github.com/0987363/vsub/middleware"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultAddress = ":10090"

var serveCmd = &cobra.Command{
	Use:    "serve",
	Short:  "Start black server",
	PreRun: LoadConfiguration,
	Run:    serve,
}

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP(
		"address",
		"a",
		defaultAddress,
		"Address the server binds to",
	)
	viper.BindPFlag("address", serveCmd.Flags().Lookup("address"))
}

func serve(cmd *cobra.Command, args []string) {
	middleware.ConnectSession(viper.GetString("authentication.cookie_key"))

	if err := middleware.ConnectDB(viper.GetString("database.mongodb")); err != nil {
		log.Fatalln("connect to db failed: ", err)
	}

	handlers.RootMux.Run(viper.GetString("address"))
}
