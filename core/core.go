package core

import (
	"auth/api"
	"auth/conf"
	"auth/log"
	"auth/services"
	"auth/services/health"

	"context"
	"net/http"
	"sync"
	"time"

	// "github.com/go-playground/validator/v10"
	// "git.bitcubix.io/go/validation"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Core struct {
	Log        *log.Logger
	Conf       *conf.Config
	Database   *gorm.DB
	httpServer *http.Server
	router     *mux.Router
	handler    http.Handler
	api        *api.Api
	services   *services.Services
	state      health.State
	// Validator  *validation.Validator
	// wg holds registered processes for graceful shutdown
	wg *sync.WaitGroup
	// context holds global context
	context globalContext
	// shutdownFuncs runs at shutdown
	shutdownFuncs []func() error
}

type globalContext struct {
	cancel context.CancelFunc
	ctx    context.Context
}

// New Core
func New(cfgFile string, isDebug bool, logFile string) (*Core, error) {
	c := &Core{}

	c.state = health.StateStarting
	ctx, cancel := context.WithCancel(context.Background())
	c.wg = &sync.WaitGroup{}
	c.context = globalContext{
		cancel: cancel,
		ctx:    ctx,
	}
	c.Conf = c.newConf(cfgFile)
	c.Log = c.newLog(isDebug, logFile)
	// c.translator = c.newTranslator()
	// c.Validator = c.newValidator()
	c.Database = c.NewDatabase()
	c.services = c.newServices()
	// TODO: router with return
	//c.router = c.newRouter()
	c.newRouter()
	c.api = c.newApi()
	c.state = health.StateRunning

	c.Log.Info().Msg("New core done")
	return c, nil
}

func (c *Core) StartServer() {

	var err error
	c.Log.Info().Msg("Starting server")

	c.httpServer = &http.Server{
		Addr:         c.Conf.Server.URL(),
		Handler:      c.handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	listen := func() error {
		if err = c.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			c.Log.Error().Err(err).Msg("failed to start http server")
			return err
		}

		return nil
	}

	c.Log.Info().Str("URL", c.httpServer.Addr).Msg("Server listen")

	for {
		if c.state == health.StateStopping || c.state == health.StateStopped {
			c.Log.Debug().Msgf("Skipping restarts of server because app is not in running state. State is:", c.state)
			return
		}
		if err = listen(); err != nil {
			time.Sleep(2 * time.Second)
			c.Log.Error().Err(err).Msgf("Error on", c.httpServer.Addr)
			c.Log.Debug().Msgf("Restarting server", c.httpServer.Addr)
			continue
		}
		return
	}
}

func (c *Core) Shutdown(ctx context.Context) error {

	if c.state != health.StateRunning {
		c.Log.Warn().Msg("server cannot be shutdown since current state is not 'running'")
		return nil
	}

	c.state = health.StateStopping
	c.Log.Debug().Msg("Server stopping gracefully..")
	defer func() {
		c.state = health.StateStopped
	}()

	c.context.cancel()
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	// loop tru shutdownfuncs
	for _, fn := range c.shutdownFuncs {
		funcErr := fn()
		if funcErr != nil {
			c.Log.Error().Err(funcErr).Msg("Shutdown func failed")
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}

	return nil
}

// registerShutdownFunc is used to register a new shutdown function
func (c *Core) registerShutdownFunc(fn func() error) {
	c.shutdownFuncs = append(c.shutdownFuncs, fn)
}
