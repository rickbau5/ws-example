package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/rickbau5/ws-example/cmd"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"
)

var errShuttingDown = errors.New("shutting down")

func main() {
	grp, ctx := errgroup.WithContext(context.Background())

	grp.Go(func() error { return cmd.Server(ctx) })
	grp.Go(func() error { return awaitSignal(ctx) })

	if err := grp.Wait(); err != nil {
		if errors.Is(err, errShuttingDown) {
			return
		}

		_, _ = fmt.Fprintf(os.Stderr, "shutting down with error: %s", err)
		os.Exit(1)
	}
}

func awaitSignal(ctx context.Context) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, unix.SIGTERM)

	defer signal.Reset(os.Interrupt, unix.SIGTERM)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case sig := <-sigs:
		return fmt.Errorf("%w: %s", errShuttingDown, sig)
	}
}
