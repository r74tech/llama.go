// Copyright (c) 2017-2025 The qitmeer developers

package app

import (
	"fmt"
	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/wrapper"
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
	err = a.cfg.Load()
	if err != nil {
		return err
	}
	if a.cfg.IsLonely() {
		return wrapper.LlamaApp(a.cfg)
	}
	return fmt.Errorf("The server mode is still under development")
}

func (a *App) Stop() error {
	return nil
}
