package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/run"

	"github.com/ghostsquad/s3-file-explorer/internal/clioptions/iostreams"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

const binaryName = "fn"

func main() {
	// TODO make the output stream configurable
	fmt.Fprintf(os.Stdout, "%s version: %s %s %s %s", binaryName, version, commit, date, builtBy)

	ctx := context.Background()

	streams := iostreams.NewStdIOStreams()
	r := setupRouter(streams)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	var group run.Group

	group.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	group.Add(func() error {
		return srv.ListenAndServe()
	}, func(err error) {
		// https://github.com/gin-gonic/gin#graceful-shutdown-or-restart

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		fmt.Fprintf(streams.ErrOut(), "rungroup received error: %s", err)

		if err = srv.Shutdown(ctx); err != nil {
			fmt.Fprintf(streams.ErrOut(), "Server forced to shutdown: %s", err)
		}
	})

	err := group.Run()
	if err != nil {
		// TODO cleanup, these error print outs are probably very repetitive
		_, err2 := fmt.Fprintf(streams.ErrOut(), "final error: %s", err)
		if err2 != nil {
			panic(err2)
		}
		os.Exit(1)
	}
}

func setupRouter(streams iostreams.IOStreams) *gin.Engine {
	logger := gin.LoggerWithWriter(streams.Out())
	r := gin.New()
	r.Use(logger, gin.Recovery())

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		c.JSON(http.StatusOK, gin.H{"user": user, "value": "test"})
	})

	return r
}
