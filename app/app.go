// Copyright (c) 2017-2025 The qitmeer developers

package app

import (
	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/grpc"
	"github.com/Qitmeer/llama.go/wrapper"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

type App struct {
	ctx   *cli.Context
	cfg   *config.Config
	grSer *grpc.Service
}

func NewApp(ctx *cli.Context, cfg *config.Config) *App {
	app := &App{
		ctx:   ctx,
		cfg:   cfg,
		grSer: grpc.New(ctx, cfg),
	}
	return app
}

func (a *App) Start() error {
	err := initLog(a.cfg)
	if err != nil {
		return err
	}
	log.Info("Start App")
	err = a.cfg.Load()
	if err != nil {
		return err
	}
	if a.cfg.Interactive {
		return wrapper.LlamaInteractive(a.cfg)
	} else if a.cfg.IsLonely() {
		content, err := wrapper.LlamaGenerate(a.cfg)
		if err != nil {
			return err
		}
		log.Info(content)
		return nil
	}
	return a.grSer.Start()
}

func (a *App) Stop() error {
	log.Info("Stop App")
	return nil
}
