package telemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"io"
	"log"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type Telemetry struct {
	Log   *log.Logger
	tFile *os.File
	tp    *trace.TracerProvider
}

func NewTelemetry(file *os.File, l *log.Logger) *Telemetry {
	return &Telemetry{
		Log:   l,
		tFile: file,
	}
}

// newExporter returns a console exporter.
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("mercari"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}

func (t *Telemetry) InitProvider() *trace.TracerProvider {
	exp, err := newExporter(t.tFile)
	if err != nil {
		t.Log.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	t.tp = tp
	return t.tp
}

func (t *Telemetry) Shutdown() {
	if t.tp == nil {
		panic("Should Init Provider First")
	}
	if err := t.tp.Shutdown(context.Background()); err != nil {
		t.Log.Fatal(err)
	}
}

func (t *Telemetry) StartTracerProvider() {
	if t.tp == nil {
		panic("Should Init Provider First")
	}
	otel.SetTracerProvider(t.tp)
}
