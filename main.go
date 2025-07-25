package main
func main() {
	app := fiber.New()
	api := app.Group("/api")
	apiV1 := api.Group("/v1")
	log.Fatal(app.Listen(":7000"))
}
