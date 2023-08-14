package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"transfer/internal/app/pkg/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

func Init(ctx context.Context) (func(), error) {

	injector, injectorCleanFunc, err := BuildInjector()
	if err != nil {
		return nil, err
	}

	InitLoger()
	httpServerCleanfunc := InitHTTPServer(ctx, injector.Engine)

	return func() {
		injectorCleanFunc()
		httpServerCleanfunc()
	}, nil
}

func InitLoger() {
	defer logger.Sync()
	if viper.GetViper().GetBool("Global.Debug") {
		logger.SetLevel(zapcore.Level(-1))
		return
	}
	logger.SetLevel(zapcore.Level(viper.GetViper().GetInt("log.Level")))
}

func InitHTTPServer(ctx context.Context, handler http.Handler) func() {

	host := viper.GetViper().GetString("HTTP.Host")

	port := viper.GetViper().GetInt("HTTP.Port")
	addr := fmt.Sprintf("%s:%d", host, port)

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  viper.GetViper().GetDuration("HTTP.ReadTimeout") * time.Second,
		WriteTimeout: viper.GetViper().GetDuration("HTTP.WriteTimeout") * time.Second,
		IdleTimeout:  viper.GetViper().GetDuration("HTTP.IdleTimeout") * time.Second,
	}

	go func() {
		logger.Infof("HTTP server is running at %s.", addr)

		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}

	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(60))
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal(err.Error())
		}
	}
}

// Run 启动gin
func Run(ctx context.Context) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := Init(ctx)
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		logger.Infof("Receive signal[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFunc()
	logger.Infof("Server exit")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}
