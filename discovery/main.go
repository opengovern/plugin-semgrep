package main

import (
	"context"
	"fmt"
	"github.com/opengovern/og-describer-template/discovery/pkg"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	if err := WorkerCommand().ExecuteContext(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func WorkerCommand() *cobra.Command {
	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cmd.SilenceUsage = true
			logger, err := zap.NewProduction()
			if err != nil {
				return err
			}

			w, err := pkg.NewWorker(
				logger,
				cmd.Context(),
			)
			if err != nil {
				return err
			}

			return w.Run(ctx)
		},
	}

	return cmd
}
