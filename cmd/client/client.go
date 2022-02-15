package main

import (
	"strings"
	"time"

	"github.com/emyrk/grow/game/world"

	"github.com/emyrk/grow/game/events"

	"github.com/emyrk/grow/internal/testdata"

	"github.com/emyrk/grow/client/network"

	"github.com/emyrk/grow/client/render"

	"golang.org/x/xerrors"

	mycmd "github.com/emyrk/grow/cmd"
	"github.com/emyrk/grow/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
)

func init() {
	clientCmd.Flags().StringP("address", "a", "ws://localhost:8060", "Server address")
	mycmd.RootCmd.AddCommand(clientCmd)
	mycmd.RootCmd.AddCommand(localClient)
}

func main() {
	mycmd.RootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return xerrors.Errorf("Run the 'client' subcommand")
	}
	mycmd.Run()
}

var localClient = &cobra.Command{
	Use:   "local",
	Short: "Local copy of the game",
	RunE: func(cmd *cobra.Command, args []string) error {
		//ctx := cmd.Context()
		logger := mycmd.MustLogger(cmd)

		// TODO: Get all these game settings from the server
		gD := testdata.TestGame()
		me := world.RandomPlayer()
		gc := game.NewGameClient(logger, gD.GameCfg)
		err := gc.SendGameMessage(game.NewEventMsgPayload([]events.Event{
			&events.PlayerJoin{
				PlayerID: me.ID,
				Color:    me.Color,
				Team:     me.Team,
			},
		}))
		if err != nil {
			return err
		}
		gr := render.NewGameRenderer(gc, me)

		ebiten.SetWindowSize(testdata.ScreenWidth, testdata.ScreenHeight)
		ebiten.SetWindowTitle("Game")
		ebiten.SetWindowResizable(true)
		if err := ebiten.RunGame(gr); err != nil {
			logger.Fatal().Err(err).Msg("game crashed")
		}
		return nil
	},
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client of the game",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := mycmd.MustLogger(cmd)
		address, _ := cmd.Flags().GetString("address")

		var nc *network.NetworkClient
		var err error
		for {
			nc, err = network.Connect(ctx, logger, strings.TrimRight(address, "/")+"/socket")
			if err != nil {
				logger.Err(err).Msg("connect to server, will try again")
				time.Sleep(time.Second)
				continue
			}
			break
		}

		msgs := nc.ReadMessages(ctx)
		// TODO: Get all these game settings from the server
		gD := testdata.TestGame()
		me := world.RandomPlayer()

		gc := game.NewGameClient(logger, gD.GameCfg).UseServer(
			nc.SendGameMessage(ctx),
		)
		go network.HandleSocketMessages(ctx, gc, msgs)
		gr := render.NewGameRenderer(gc, me)
		err = gc.SendGameMessage(game.NewEventMsgPayload([]events.Event{
			&events.PlayerJoin{
				PlayerID: me.ID,
				Color:    me.Color,
				Team:     me.Team,
			},
		}))
		if err != nil {
			return xerrors.Errorf("player join evt: %w", err)
		}

		ebiten.SetWindowSize(gD.GameCfg.Width, gD.GameCfg.Height)
		ebiten.SetWindowTitle("Game")
		ebiten.SetWindowResizable(true)
		if err := ebiten.RunGame(gr); err != nil {
			logger.Fatal().Err(err).Msg("game crashed")
		}

		return nil
	},
}
