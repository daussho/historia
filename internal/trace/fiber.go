package trace

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type fiberHandler func(ctx *fiber.Ctx) error

func FiberHandler(handler fiberHandler) fiberHandler {
	return func(ctx *fiber.Ctx) error {
		method := ctx.Request().Header.Method()
		uri := ctx.Request().URI().RequestURI()

		rootSegment := fmt.Sprintf("[%s]%s", method, uri)

		span, ctx := StartSpanWithFiberCtx(ctx, rootSegment, map[string]string{
			"method": string(method),
			"uri":    string(uri),
			"host":   string(ctx.Hostname()),
		})
		defer func() {
			span.FinishAndSubmit(ctx.UserContext())
		}()

		return handler(ctx)
	}
}

func StartSpanWithFiberCtx(ctx *fiber.Ctx, segmentName string, tags map[string]string) (*Span, *fiber.Ctx) {
	span, newCtx := StartSpanWithCtx(ctx.UserContext(), segmentName, tags)
	ctx.SetUserContext(newCtx)

	return span, ctx
}
