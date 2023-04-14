package logger

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

type contextKey struct {
	name string
}

func (p contextKey) String() string {
	return fmt.Sprintf("context key: %s", p.name)
}

var logRecordContextKey = contextKey{"request_log"}

func AddAttrs(c *fiber.Ctx, attrs ...slog.Attr) {
	if r, ok := c.Locals(logRecordContextKey).(*slog.Record); ok {
		r.AddAttrs(attrs...)
	}
}

func Add(c *fiber.Ctx, args ...any) {
	if r, ok := c.Locals(logRecordContextKey).(*slog.Record); ok {
		r.Add(args...)
	}
}

func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)
	return func(c *fiber.Ctx) error {
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}

		startedAt := time.Now()
		r := slog.NewRecord(startedAt, slog.LevelInfo, "", 0)

		// available via local context in chain handler, see: @Record
		c.Locals(logRecordContextKey, &r)

		err := c.Next()

		latency := time.Since(startedAt)
		defer func(e error) {
			if e != nil {
				r.Level = slog.LevelError
			}
			var (
				attrs = make([]slog.Attr, len(cfg.BuiltinAttrs))
				attr  slog.Attr
			)
			for i, key := range cfg.BuiltinAttrs {
				switch key {
				case RefererTag:
					attr = slog.String(RefererTag, c.Get(fiber.HeaderReferer))
				case UserAgentTag:
					attr = slog.Any(UserAgentTag, c.Get(fiber.HeaderUserAgent))
				case IPTag:
					attr = slog.String(IPTag, c.IP())
				case IPsTag:
					attr = slog.String(IPsTag, c.Get(fiber.HeaderXForwardedFor))
				case LatencyTag:
					attr = slog.Duration(LatencyTag, latency.Round(time.Microsecond))
				case LatencyHumanTag:
					attr = slog.String(LatencyHumanTag, fmt.Sprintf("%7v", latency))
				case StatusTag:
					attr = slog.Int(StatusTag, c.Response().StatusCode())
				case PathTag:
					attr = slog.String(PathTag, c.Path())
				case UrlTag:
					attr = slog.String(UrlTag, c.OriginalURL())
				case MethodTag:
					attr = slog.String(MethodTag, c.Method())
				case BodyTag:
					attr = slog.Any(BodyTag, c.Body())
				case BytesReceivedTag:
					attr = slog.Any(BytesReceivedTag, len(c.Request().Body()))
				case BytesSentTag:
					attr = func() slog.Attr {
						sent := 0
						if c.Response().Header.ContentLength() > 0 {
							sent = len(c.Response().Body())
						}
						return slog.Int(BytesSentTag, sent)
					}()
				case RequestHeadersTag:
					rh := fiber.Map{}
					for k, v := range c.GetReqHeaders() {
						rh[k] = v
					}
					attr = slog.Any(RequestHeadersTag, rh)
				case QueryStringTag:
					rh := fiber.Map{}
					qs := c.Request().URI().QueryArgs()
					qs.VisitAll(func(key, value []byte) {
						rh[string(key[:])] = string(value[:])
					})
					attr = slog.Any(QueryStringTag, rh)
				}
				attrs[i] = attr
			}
			r.AddAttrs(attrs...)
			if cfg.CustomAttr != nil {
				for _, fn := range cfg.CustomAttr {
					fn(c, &r, e)
				}
			}
			_ = cfg.Logger.Handler().Handle(c.Context(), r)
		}(err)
		return err
	}
}
