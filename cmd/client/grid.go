package main

import (
	"github.com/emyrk/grow/client/render"
	mycmd "github.com/emyrk/grow/cmd"
	"github.com/emyrk/grow/game/world/grid"
	"github.com/emyrk/grow/internal/testdata"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
)

func init() {
	mycmd.RootCmd.AddCommand(gridCmd)
}

var gridCmd = &cobra.Command{
	Use: "grid",
	RunE: func(cmd *cobra.Command, args []string) error {
		//ctx := cmd.Context()
		logger := mycmd.MustLogger(cmd)
		g := render.NewGridRenderer(logger, grid.NewGrid(testdata.ScreenWidth, testdata.ScreenHeight))

		ebiten.SetWindowSize(g.Width, g.Height)
		ebiten.SetWindowTitle("GridTesting")
		ebiten.SetWindowResizable(true)
		if err := ebiten.RunGame(g); err != nil {
			logger.Fatal().Err(err).Msg("game crashed")
		}
		return nil
	},
}
