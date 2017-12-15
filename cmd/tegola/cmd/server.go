package cmd

import (
	"github.com/spf13/cobra"
	"github.com/airmap/tegola/server"
)

var (
	serverPort      string
	defaultHTTPPort = ":8080"
)

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "Use tegola as a tile server",
	Long:  `Use tegola as a vector tile server. Maps tiles will be served at /maps/:map_name/:z/:x/:y`,
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()

		//	check config for server port setting
		//	if you set the port via the comand line it will override the port setting in the config
		if serverPort == defaultHTTPPort && conf.Webserver.Port != "" {
			serverPort = conf.Webserver.Port
		}

		//	set our server version
		server.Version = Version
		server.HostName = conf.Webserver.HostName

		//	start our webserver
		server.Start(serverPort)
	},
}
