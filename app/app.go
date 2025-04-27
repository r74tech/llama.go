// Copyright (c) 2017-2025 The qitmeer developers

package app

import "github.com/Qitmeer/llama.go/config"

type App struct {
}

func (a *App) Start(cfg *config.Config) error {
	err := cfg.Load()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) Stop() error {
	return nil
}
