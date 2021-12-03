package gofibersentry

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/LegalForceLawRAPC/go-template/api/utils"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/gofiber/fiber/v2"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
)

func SentryInit() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              viper.GetString("SENTRY_DSN"),
		Environment:      viper.GetString("ENVIRONMENT"),
		AttachStacktrace: true,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			return event
		},
		// Or provide a custom sampler:
		TracesSampler: sentry.TracesSamplerFunc(func(ctx sentry.SamplingContext) sentry.Sampled {
			return sentry.SampledTrue
		}),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
}

type contextKey int

const ContextKey = contextKey(1)

type Handler struct {
	repanic         bool
	waitForDelivery bool
	timeout         time.Duration
}

type Options struct {
	// Repanic configures whether Sentry should repanic after recovery, in most cases it should be set to false,
	// as fasthttp doesn't include it's own Recovery handler.
	Repanic bool
	// WaitForDelivery configures whether you want to block the request before moving forward with the response.
	// Because fasthttp doesn't include it's own Recovery handler, it will restart the application,
	// and event won't be delivered otherwise.
	WaitForDelivery bool
	// Timeout for the event delivery requests.
	Timeout time.Duration
}

func New(options Options) *Handler {
	timeout := options.Timeout
	if timeout == 0 {
		timeout = 2 * time.Second
	}
	return &Handler{
		repanic:         options.Repanic,
		timeout:         timeout,
		waitForDelivery: options.WaitForDelivery,
	}
}

// Handle wraps fasthttp.RequestHandler and recovers from caught panics.
func (h *Handler) Handle(ctx *fiber.Ctx) error {
	if ctx.Method() == "GET" || ctx.Method() == "POST" || ctx.Method() == "DELETE" || ctx.Method() == "PATCH" {
		goCtx := ctx.Context()
		go func(ct *fiber.Ctx, goCtx *fasthttp.RequestCtx) {
			defer utils.Recover()
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				defer utils.Recover()
				if ct != nil {
					r := http.Request{}
					err := fasthttpadaptor.ConvertRequest(goCtx, &r, false)
					if err != nil {
						scope.SetRequest(&r)
					}
					// prevent capturing body in case of get method
					if ct.Method() != http.MethodGet {
						var body map[string]interface{}
						err := ct.BodyParser(&body)
						if err != nil {
							return
						}
						if ct.Request().Body() != nil {
							scope.SetRequestBody(ct.Request().Body())
						}
					}
				}
			})
		}(ctx, goCtx)
	}
	return ctx.Next()
}

func (h *Handler) RecoverWithSentry(hub *sentry.Hub, ctx *fiber.Ctx) {
	if err := recover(); err != nil {
		eventID := hub.RecoverWithContext(
			context.WithValue(context.Background(), sentry.RequestContextKey, ctx),
			err,
		)
		if eventID != nil && h.waitForDelivery {
			hub.Flush(h.timeout)
		}
		if h.repanic {
			panic(err)
		}
	}
}

func Convert(ctx *fiber.Ctx) *http.Request {
	log.Printf("Goroutine Request: %v", ctx.Request())
	r := new(http.Request)

	r.Method = ctx.Method()

	// Headers
	r.Header = make(http.Header)

	r.Header.Add("Host", ctx.Hostname())
	ctx.Request().Header.VisitAll(func(key, value []byte) {
		r.Header.Add(string(key), string(value))
	})

	r.Host = ctx.Hostname()

	// Cookies
	ctx.Request().Header.VisitAllCookie(func(key, value []byte) {
		r.AddCookie(&http.Cookie{Name: string(key), Value: string(value)})
	})

	// Env
	r.RemoteAddr = ctx.IP()

	// QueryString
	r.URL.RawQuery = string(ctx.Request().RequestURI())

	// Body
	r.Body = ioutil.NopCloser(bytes.NewReader(ctx.Request().Body()))

	return r
}
