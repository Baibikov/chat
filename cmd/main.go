package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"chat/internal/app/delivery/rest"
	"chat/internal/app/delivery/service/socket"
)

func main() {
	logrus.Info("init ctx")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	c := make(chan os.Signal)
	signal.Notify(c,  syscall.SIGINT)
	go func() {
		oscall := <-c
		logrus.Infof("system call:%+v", oscall)
		cancel()
	}()

	logrus.Info("init socket service")
	s := socket.NewSocket()
	g.Go(func() error {
		return s.ListenAndRoute(ctx)
	})

	handler := rest.New(s)
	http.Handle("/", handler.Setup())


	logrus.Info("init http api")
	g.Go(func() error {
		return http.ListenAndServe(":9091", nil)
	})

	go func() {
		err := g.Wait()
		if err != nil {
			logrus.Debug(err)
			cancel()
		}
	}()

	shutdown(ctx)
}

func shutdown(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			time.Sleep(time.Second*3)
		return

		default:

		}
	}
}

