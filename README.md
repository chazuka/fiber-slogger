It's an opionated [fiber](https://gofiber.io/) request logger middleware using [slog](https://pkg.go.dev/golang.org/x/exp/slog) library.

It allows adding attribute from anywhere within the handler.

**Usage**

registering middleware with custom config

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

Adding extra attribute from anywhere within handler

```go
router.Post("/", func(ctx *fiber.Ctx) error {
    ...
    logger.AddAttrs(ctx, slog.Any("params", params))
    logger.Add(ctx, "params", params))
    ...
})
```
