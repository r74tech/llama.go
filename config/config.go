// Copyright (c) 2017-2025 The qitmeer developers

package config

import (
	"github.com/urfave/cli/v2"
)

var (
	Conf = &Config{}

	Model = &cli.StringFlag{
		Name:        "model",
		Aliases:     []string{"m"},
		Usage:       "Specify the path to the LLaMA model file",
		Destination: &Conf.Model,
	}

	CtxSize = &cli.IntFlag{
		Name:        "ctx-size",
		Aliases:     []string{"c"},
		Usage:       "Set the size of the prompt context. The default is 4096, but if a LLaMA model was built with a longer context, increasing this value will provide better results for longer input/inference",
		Value:       4096,
		Destination: &Conf.CtxSize,
	}

	AppFlags = []cli.Flag{
		Model,
		CtxSize,
	}
)

type Config struct {
	Model   string
	CtxSize int
}

func (c *Config) Load() error {
	return nil
}
