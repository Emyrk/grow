package main

import (
	"image/color"

	"golang.org/x/xerrors"

	mycmd "github.com/emyrk/grow/cmd"
	"github.com/emyrk/grow/game"
	"github.com/emyrk/grow/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
)

func init() {
	mycmd.RootCmd.AddCommand(clientCmd)
}

func main() {
	mycmd.RootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return xerrors.Errorf("Run the 'client' subcommand")
	}
	mycmd.Run()
}

const (
	screenWidth  = 600
	screenHeight = 600
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client of the game",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := mycmd.MustLogger(cmd)

		me := world.NewPlayer(0, color.RGBA{
			// 844a93
			R: 0x84,
			G: 0x4a,
			B: 0x93,
			A: 0xff,
		})
		players := world.NewPlayerSet()
		me = players.AddPlayer(me)
		g := game.NewGame(logger, screenWidth, screenHeight, players, me)
		ebiten.SetWindowSize(screenWidth, screenHeight)
		ebiten.SetWindowTitle("Game")
		ebiten.SetWindowResizable(true)
		if err := ebiten.RunGame(g); err != nil {
			logger.Fatal().Err(err).Msg("game crashed")
		}

		return nil
	},
}
