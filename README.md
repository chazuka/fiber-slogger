## Archived

This project is archived in favor of the official Fiber slog middleware:

- [`github.com/gofiber/contrib/fiberslog`](https://github.com/gofiber/contrib/tree/main/fiberslog)

This package remains here for existing users, but new projects should prefer the official package above.

It is an opinionated [Fiber](https://gofiber.io/) request logger middleware using [slog](https://pkg.go.dev/log/slog).

It allows adding attributes from anywhere within the handler.

## Usage

Registering middleware with custom config:

```go
app.Use(logger.New(logger.Config{
    // skip logger with filter
    Filter: func(ctx *fiber.Ctx) bool {
        return ctx.Path() == "/exclude"
    },
    // customize attributes with builtin tag
    BuiltinAttrs []string{PathTag,StatusTag}
    // adding custom attributes with custom function
    CustomAttr: []logger.CustomFunc{
		func(c *fiber.Ctx, r *slog.Record, e error) {
            r.AddAttrs(slog.String("custom_key", "custom_value"))
        },
    },
    // assign slog logger
    Logger: slog.Default(),
}))
```

Adding extra attributes from anywhere within a handler:

```go
router.Post("/", func(ctx *fiber.Ctx) error {
    ...
    logger.AddAttrs(ctx, slog.Any("params", params))
    logger.Add(ctx, "params", params))
    ...
})
```
