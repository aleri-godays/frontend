package http

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

//AddRequestIDToContext is a middleware that create a logger with a request id
func AddRequestIDToContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			c.Set("request_id", requestID)
			return next(c)
		}
	}
}

//AddLoggerToContext is a middleware that create a logger with a request id
func AddLoggerToContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			logger := log.WithFields(log.Fields{
				"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
			})
			c.Set("logger", logger)
			return next(c)
		}
	}
}

func AddLoginToContext() echo.MiddlewareFunc {
	skipper := func(c echo.Context) bool {
		p := c.Path()

		allowed := map[string]bool{
			"/":             true,
			"/health":       true,
			"/login":        true,
			"/logout":       true,
			"/authcallback": true,
		}
		if ok, exists := allowed[p]; ok && exists {
			return true
		}

		matched, err := regexp.MatchString("^/static/", p)
		if err != nil {
			log.WithFields(log.Fields{
				"path":  p,
				"error": err,
			}).Error("regex failed")
			return false
		}

		return matched
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}
			user := c.Get("user").(*jwt.Token)
			token := c.Get("user").(*jwt.Token)
			c.Set("user_jwt", token.Raw)
			claims := user.Claims.(jwt.MapClaims)
			c.Set("login", claims["user"])
			return next(c)
		}
	}
}

//Logger is a middleware that logs all requests
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			var err error
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			req := c.Request()
			res := c.Response()
			reqSize := req.Header.Get(echo.HeaderContentLength)
			if reqSize == "" {
				reqSize = "0"
			}

			fields := log.Fields{
				"path":          req.RequestURI,
				"method":        req.Method,
				"status":        res.Status,
				"request_size":  reqSize,
				"response_size": res.Size,
				"duration":      stop.Sub(start).String(),
				"error":         err,
			}

			if err == nil {
				fields["error"] = ""
			}
			logger := c.Get("logger").(*log.Entry)
			logger.WithFields(fields).Info("request")

			return err
		}
	}
}

func Tracing() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//no tracing for static files
			if c.Path() == "/*" || c.Path() == "/metrics" {
				return next(c)
			}

			tracer := opentracing.GlobalTracer()
			req := c.Request()
			opName := "HTTP " + req.Method + " URL: " + c.Path()

			var span opentracing.Span
			if ctx, err := tracer.Extract(opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(req.Header)); err != nil {
				span = tracer.StartSpan(opName)
			} else {
				span = tracer.StartSpan(opName, ext.RPCServerOption(ctx))
			}

			ext.HTTPMethod.Set(span, req.Method)
			ext.HTTPUrl.Set(span, req.URL.String())
			ext.Component.Set(span, "rest")

			req = req.WithContext(opentracing.ContextWithSpan(req.Context(), span))
			c.SetRequest(req)

			c.Set("span", span)

			defer func() {
				status := c.Response().Status
				committed := c.Response().Committed
				ext.HTTPStatusCode.Set(span, uint16(status))
				if status >= http.StatusInternalServerError || !committed {
					ext.Error.Set(span, true)
					span.SetBaggageItem("errorval", strconv.Itoa(http.StatusInternalServerError))
				}
				span.Finish()
			}()

			return next(c)
		}
	}
}

func Instrumenting() echo.MiddlewareFunc {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "godays_request_by_http_status_code",
		},
		[]string{"status"},
	)
	prometheus.MustRegister(requestCounter)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "godays_request_duration",
			Help:    "time for reponse",
			Buckets: []float64{0.001, 0.002, 0.003, 0.004, 0.005, 0.01, 0.1},
		},
		[]string{"method"},
	)
	prometheus.MustRegister(requestDuration)

	bla := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "godays_request_duration_2",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method"},
	)
	prometheus.MustRegister(bla)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/*" || c.Path() == "/metrics" {
				return next(c)
			}

			defer func(start time.Time) {
				d := time.Since(start)
				requestDuration.WithLabelValues(c.Request().Method).Observe(d.Seconds())
				bla.WithLabelValues(c.Request().Method).Observe(d.Seconds())
			}(time.Now())

			err := next(c)

			var code int
			if err == nil {
				code = c.Response().Status
			} else if echoError, ok := err.(*echo.HTTPError); ok {
				code = echoError.Code
			} else {
				code = 500
			}

			requestCounter.WithLabelValues(strconv.Itoa(code)).Inc()

			return err
		}
	}
}

//
//func InstrumentingRequestCount() echo.MiddlewareFunc {
//	return func(next echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) error {
//
//		}
//	}
//}
