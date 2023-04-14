package logger

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

const (
	RefererTag        = "referer"
	UserAgentTag      = "user_agent"
	IPTag             = "ip"
	IPsTag            = "ips"
	LatencyTag        = "latency"
	LatencyHumanTag   = "latency_human"
	StatusTag         = "status"
	PathTag           = "path"
	UrlTag            = "url"
	MethodTag         = "method"
	BodyTag           = "body"
	BytesReceivedTag  = "bytes_received"
	BytesSentTag      = "bytes_sent"
	RequestHeadersTag = "request_headers"
	QueryStringTag    = "query_string"
)

type (
	CustomFunc func(*fiber.Ctx, *slog.Record, error)
	Config     struct {
		Filter       func(*fiber.Ctx) bool
		CustomAttr   []CustomFunc
		BuiltinAttrs []string
		Logger       *slog.Logger
	}
)

var defaultConfig = Config{
	BuiltinAttrs: []string{
		RefererTag,
		UserAgentTag,
		IPTag,
		IPsTag,
		LatencyTag,
		LatencyHumanTag,
		StatusTag,
		PathTag,
		UrlTag,
		MethodTag,
		BodyTag,
		BytesReceivedTag,
		BytesSentTag,
		RequestHeadersTag,
		QueryStringTag,
	},
	Logger: slog.Default(),
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return defaultConfig
	}
	cfg := config[0]
	if cfg.BuiltinAttrs == nil {
		cfg.BuiltinAttrs = defaultConfig.BuiltinAttrs
	}
	if cfg.Logger == nil {
		cfg.Logger = defaultConfig.Logger
	}

	return cfg
}
