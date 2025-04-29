// Copyright (c) 2017-2025 The qitmeer developers

package app

import (
	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/wrapper"
)

type App struct {
}

func (a *App) Start(cfg *config.Config) error {
	err := initLog(cfg)
	if err != nil {
		return err
	}
	err = cfg.Load()
	if err != nil {
		return err
	}
	return wrapper.LlamaApp()
}

func (a *App) Stop() error {
	return nil
}
