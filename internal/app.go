package internal

import (
	"context"
	"github.com/aleri-godays/frontend/internal/config"
	"github.com/aleri-godays/frontend/internal/http"
	"github.com/aleri-godays/frontend/internal/project"
	"github.com/aleri-godays/frontend/internal/timetracking"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type App struct {
	server      *http.Server
	conf        *config.Config
	traceCloser io.Closer
}

func NewApp(conf *config.Config) *App {
	configureLogger(conf.LogLevel)

	a := &App{
		conf:        conf,
		traceCloser: initTracer(conf.ServiceName),
		server: http.NewServer(
			conf,
			project.NewClient(conf),
			timetracking.NewClient(conf)),
	}

	return a
}

func (a *App) Run() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	//start webserver
	go func() {
		log.WithFields(log.Fields{
			"port":      a.conf.HTTPPort,
			"log_level": a.conf.LogLevel,
			"version":   a.conf.Version,
		}).Info("starting server")
		if err := a.server.Start(); err != nil {
			log.WithFields(log.Fields{
				"reason": err,
			}).Info("shutting down the server")
			quit <- os.Interrupt
		}
	}()

	<-quit
	log.Info("received a stop signal")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	//close http connections
	if err := a.server.Shutdown(ctx); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("errors during http server shutdown")
	} else {
		log.Info("http server closed")
	}

	a.traceCloser.Close()
}

func initTracer(serviceName string) io.Closer {
	defcfg := jaegerconfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	cfg, err := defcfg.FromEnv()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not parse jaeger env vars")
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not initialize jaeger tracer")
	}
	opentracing.SetGlobalTracer(tracer)

	return closer
}

func configureLogger(logLevel string) {
	switch strings.ToLower(logLevel) {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}
}
