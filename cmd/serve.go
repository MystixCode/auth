package cmd

import (
	"auth/core"

	"os"

	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the auth",
	Long:  `Start auth and serve it on a http server`,
	Run:   serve,
}

func serve(cmd *cobra.Command, args []string) {
	appCore, err := core.New(cfgFile, isDebug, logFile)
	if err != nil {
		//Log.Debug().Err(err).Msg("initCore error")
		os.Exit(2)
	}
	go appCore.StartServer()

	graceful(appCore, 30*time.Second)
}

func graceful(core *core.Core, timeout time.Duration) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := core.Shutdown(ctx); err != nil {
		core.Log.Error().Err(err).Msg("server shutdown error")
	} else {
		core.Log.Info().Msg("Server stopped gracefully")
	}
}
