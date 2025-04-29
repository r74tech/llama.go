package app

import (
	"github.com/urfave/cli/v2"
)

func commands() []*cli.Command {
	cmds := []*cli.Command{}
	cmds = append(cmds, downloadCmd())
	return cmds
}

func downloadCmd() *cli.Command {
	return &cli.Command{
		Name:        "download",
		Aliases:     []string{"r"},
		Category:    "llama",
		Usage:       "Download model",
		Description: "Download model",
		Action: func(ctx *cli.Context) error {
			return nil
		},
	}
}
