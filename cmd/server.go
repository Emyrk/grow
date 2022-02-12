package cmd

import (
	"context"

	"github.com/emyrk/grow/server"
	"github.com/spf13/cobra"
)

func init() {
	srvCommand.Flags().IntP("port", "p", 8080, "port for the game server")
	RootCmd.AddCommand(srvCommand)
}

var srvCommand = &cobra.Command{
	Use:   "server",
	Short: "Run a game server",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")
		log := MustLogger(cmd)
		ctx, cancel := context.WithCancel(cmd.Context())

		cfg := server.ServerConfig{
			Port: port,
			Log:  log,
		}
		gs := server.NewGameServer(cfg, nil)
		go func() {
			err := gs.Start()
			if err != nil {
				log.Err(err).Msg("start server")
			}
			cancel()
		}()

		<-ctx.Done()
		err := gs.Close()
		return err
	},
}
