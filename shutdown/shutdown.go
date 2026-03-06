package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func WaitForShutdown(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	cancel()
}

//Usage

// ctx, cancel := context.WithCancel(context.Background())

// go shutdown.WaitForShutdown(cancel)

// runServer(ctx)
