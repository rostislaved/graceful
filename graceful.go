package graceful

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

type graceful struct {
	processes []process
	logger    *slog.Logger
}

func New(processes ...process) graceful {
	return graceful{
		processes: processes,
		logger:    slog.Default(),
	}
}

func (gr graceful) Start(ctx context.Context) error {
	g, gCtx := errgroup.WithContext(ctx)

	for _, process := range gr.processes {
		if process.disabled {
			continue
		}

		f := func() error {
			err := process.starter.Start(gCtx)
			if err != nil {
				gr.logger.Error(err.Error())
				gr.logger.Info("Start graceful shutdown")

				return err
			}

			return nil
		}

		g.Go(f)
	}

	err := g.Wait()
	if err != nil {
		gr.logger.Info("Application stopped gracefully")

		return err
	}

	gr.logger.Info("Every process stopped by itself with no error")

	return nil
}

func (gr graceful) SetLogger(l *slog.Logger) graceful {
	gr.logger = l

	return gr
}
