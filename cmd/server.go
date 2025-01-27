package cmd

import (
	"iht/pkg/cfg"
	"iht/pkg/server"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start the server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		port, _ := cmd.Flags().GetInt("port")
		if port == 0 {
			port = cfg.Conf.Server.Port
		}

		dir, _ := cmd.Flags().GetString("dir")
		if dir == "" {
			dir = cfg.Conf.Server.Dir
		}

		server.Serve(port, dir)

	},
}

func init() {
	serverCmd.PersistentFlags().IntP("port", "p", 7777, "server listen port")
	serverCmd.PersistentFlags().StringP("dir", "d", "./web", "server path")

	rootCmd.AddCommand(serverCmd)
}
