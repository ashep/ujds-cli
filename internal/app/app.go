package app

import (
	"context"

	ujdscli "github.com/ashep/ujds/sdk/client"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/ashep/ujds-cli/internal/command"
)

type App struct {
	rootCmd *cobra.Command
}

func New(cfg Config, l zerolog.Logger) *App {
	return &App{
		rootCmd: command.New(ujdscli.New(cfg.Host, cfg.Token, nil), l),
	}
}

func (a *App) Run(ctx context.Context, args []string) error {
	return a.rootCmd.ExecuteContext(ctx)
}
