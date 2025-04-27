// Copyright (c) 2017-2025 The qitmeer developers

package app

import (
	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/version"
	"github.com/urfave/cli/v2"
	"os"
)

func Run() error {
	a := App{}
	app := &cli.App{
		Name:    "",
		Version: version.String(),
		Authors: []*cli.Author{
			&cli.Author{
				Name: "Qitmeer",
			},
		},
		Copyright:            "(c) 2020 Qitmeer",
		Usage:                "Llama",
		Flags:                config.AppFlags,
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			err := a.Start(config.Conf)
			if err != nil {
				return err
			}
			return a.Stop()
		},
	}

	return app.Run(os.Args)
}
