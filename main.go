package main

import (
	websocket "gokripto/controllers"
	"gokripto/database"
	routes "gokripto/routes"
	"log"
	"net/http"
	"sync"

	httpmiddleware "gokripto/prom"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "gokripto/docs"
)

// @title           Go Crypto
// @version         1.0
// @description     Crypto currency app.
// @contact.name   API Support
// @contact.url    https://github.com/SiHecate
// @host      localhost:8080
// @securityDefinitions.basic  BasicAuth
func main() {
	database.Connect()
	app := fiber.New()

	app.Use(cors.New(cors.Config{ // Cross-Origin Resource Sharing
		AllowCredentials: true,
	}))

	routes.Setup(app)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		websocket.StartWebSocket(app)
		wg.Done()
	}()
	wg.Wait()

	app.Listen(":8080")

	// Create non-global registry.
	registry := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle(
		"/metrics",
		httpmiddleware.New(
			registry, nil).
			WrapHandler("/metrics", promhttp.HandlerFor(
				registry,
				promhttp.HandlerOpts{}),
			))

	log.Fatalln(http.ListenAndServe(":9090", nil))

	select {}

}
