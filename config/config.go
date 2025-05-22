// Copyright (c) 2017-2025 The qitmeer developers

package config

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
	"math"
	"path/filepath"
	"runtime"
)

const (
	defaultLogLevel     = "info"
	DefaultGrpcEndpoint = "localhost:50051"
	defaultNPredict     = 512
)

var (
	defaultHomeDir     = "."
	defaultSwaggerFile = filepath.Join(defaultHomeDir, "swagger.json")
)

var (
	Conf = &Config{}

	LogLevel = &cli.StringFlag{
		Name:        "log-level",
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

	Prompt = &cli.StringFlag{
		Name:        "prompt",
		Aliases:     []string{"p"},
		Usage:       "Provide a prompt directly as a command-line option.",
		Destination: &Conf.Prompt,
	}

	NGpuLayers = &cli.IntFlag{
		Name:        "n-gpu-layers",
		Aliases:     []string{"ngl"},
		Usage:       "When compiled with GPU support, this option allows offloading some layers to the GPU for computation. Generally results in increased performance.",
		Value:       -1,
		Destination: &Conf.NGpuLayers,
	}

	NPredict = &cli.IntFlag{
		Name:        "n-predict",
		Aliases:     []string{"n"},
		Usage:       "Set the number of tokens to predict when generating text. Adjusting this value can influence the length of the generated text.",
		Value:       defaultNPredict,
		Destination: &Conf.NPredict,
	}

	Interactive = &cli.BoolFlag{
		Name:        "interactive",
		Aliases:     []string{"i"},
		Usage:       "Run the program in interactive mode, allowing you to provide input directly and receive real-time responses",
		Value:       false,
		Destination: &Conf.Interactive,
	}

	Seed = &cli.UintFlag{
		Name:        "seed",
		Aliases:     []string{"s"},
		Usage:       "Set the random number generator (RNG) seed (default: -1, -1 = random seed).",
		Value:       math.MaxUint32,
		Destination: &Conf.Seed,
	}

	AppFlags = []cli.Flag{
		LogLevel,
		Model,
		CtxSize,
		Prompt,
		NGpuLayers,
		NPredict,
		Interactive,
		Seed,
	}
)

type Config struct {
	LogLevel    string
	Model       string
	CtxSize     int
	Prompt      string
	NGpuLayers  int
	NPredict    int
	Interactive bool
	Seed        uint
}

func (c *Config) Load() error {
	log.Debug("Try to load config")
	if len(c.Model) <= 0 {
		return fmt.Errorf("No config model")
	}
	return nil
}

func (c *Config) IsLonely() bool {
	return len(c.Prompt) > 0 || c.Interactive
}

func defaultNGpuLayers() int {
	switch runtime.GOOS {
	case "darwin":
		return -1
	}
	return 0
}
