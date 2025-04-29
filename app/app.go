// Copyright (c) 2017-2025 The qitmeer developers

package app

import (
	"fmt"
	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/wrapper"
	"github.com/ethereum/go-ethereum/log"
)

type App struct {
	cfg *config.Config
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
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
		return wrapper.LlamaGenerate(a.cfg)
	}
	return fmt.Errorf("The server mode is still under development")
}

func (a *App) Stop() error {
	log.Info("Stop App")
	return nil
}
