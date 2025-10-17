package main

import (
	"context"
	"crypto/tls"
	"github.com/SigNoz/sample-golang-app/controllers"
	"github.com/SigNoz/sample-golang-app/models"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
	signozKey    = os.Getenv("OTEL_EXPORTER_OTLP_HEADERS")
)

func initTracer() func(context.Context) error {

	/*	var secureOption otlptracegrpc.Option

		if strings.ToLower(insecure) == "false" || insecure == "0" || strings.ToLower(insecure) == "f" {
			secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		} else {
			secureOption = otlptracegrpc.WithInsecure()
		}
		/*
			exporter, err := otlptrace.New(
				context.Background(),
				otlptracegrpc.NewClient(
					secureOption,
					otlptracegrpc.WithEndpoint(collectorURL),
				),
			)


	*/
	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(collectorURL),
		otlptracehttp.WithURLPath("/v1/traces"),
		otlptracehttp.WithTLSClientConfig(&tls.Config{}),
		otlptracehttp.WithHeaders(map[string]string{
			"signoz-ingestion-key": signozKey}),
	)

	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	return exporter.Shutdown
}

func main() {
	cleanup := initTracer()
	defer cleanup(context.Background())

	r := gin.Default()
	models.ConnectDatabase()
	r.Use(otelgin.Middleware(serviceName))
	// Routes
	r.GET("/books", controllers.FindBooks)
	r.GET("/books/:id", controllers.FindBook)
	r.POST("/books", controllers.CreateBook)
	r.PATCH("/books/:id", controllers.UpdateBook)
	r.DELETE("/books/:id", controllers.DeleteBook)

	// Run the server
	r.Run(":8090")
}
