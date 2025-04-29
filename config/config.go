// Copyright (c) 2017-2025 The qitmeer developers

package config

import (
	"github.com/urfave/cli/v2"
)

const (
	defaultLogLevel = "info"
)

var (
	Conf = &Config{}

	LogLevel = &cli.StringFlag{
		Name:        "log_level",
		Aliases:     []string{"l"},
		Usage:       "Logging level {trace, debug, info, warn, error}",
		Value:       defaultLogLevel,
		Destination: &Conf.LogLevel,
	}

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
		LogLevel,
		Model,
		CtxSize,
	}
)

type Config struct {
	LogLevel string
	Model    string
	CtxSize  int
}

func (c *Config) Load() error {
	return nil
}
