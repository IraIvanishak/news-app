package middlewares

const (
	DateTimeLayout = "15:04:05 02-01-2006"
)

func RoutLoggerMiddlewareInitializer() {
	// .App().Use(logger.New(logger.Config{
	// 	Format:     "[${ip}]:${port} ${method} ${path} ${status} pid=${pid}\n",
	// 	TimeFormat: DateTimeLayout,
	// }))
}
