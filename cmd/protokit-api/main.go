// Package main implements the runtime for the protokit-api binary.
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	host = kingpin.Flag("host", "<help>").Envar("HOST").String()
	port = kingpin.Flag("port", "<help>").Envar("PORT").Default("8080").Int()
)

func main() {
	ctx := context.Background()

	kingpin.Parse()

	log, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("[protokit-api] unable to initialize logger: %v", err))
	}

	term := make(chan os.Signal, 4)
	signal.Notify(term, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(time.Second * 60))
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message": "Hello, World!"}`)
	})

	httpServer := &http.Server{
		Addr:         net.JoinHostPort(*host, strconv.Itoa(*port)),
		Handler:      r,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
	}

	go func() {
		log.Info("http server starting", zap.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal("http server closed unexpectedly", zap.Error(err))
				return
			}

			log.Info("http server closed", zap.Error(err))
		}
	}()

	<-term

	timeout := time.Second * 15
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	log.Info("shutting down", zap.Duration("timeout", timeout))

	g.Go(func() error {
		return httpServer.Shutdown(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Error("shutdown error occured", zap.Error(err))
	}
}
