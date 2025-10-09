package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// DumpBodyWithTracer creates a middleware that dumps the request and response bodies.
func DumpBodyWithTracer(tracer trace.Tracer) echo.MiddlewareFunc {
	return echomiddleware.BodyDump(func(ctx echo.Context, reqBody, respBody []byte) {
		if !strings.Contains(ctx.Request().URL.Path, "/auth") {
			_, span := tracer.Start(ctx.Request().Context(), "bodydump")
			defer span.End()
			span.SetAttributes(
				attribute.String("request.body", string(reqBody)),
				attribute.String("response.body", string(respBody)),
			)
			if resp := ctx.Response(); resp.Status != http.StatusOK {
				span.SetStatus(codes.Error, http.StatusText(resp.Status))
			}
		}
	})
}
