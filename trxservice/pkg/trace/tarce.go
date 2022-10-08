package trace

import (
	"context"
	"log"
	_ "os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// Init configures an OpenTelemetry exporter and trace provider
func Init(ctx context.Context, serviceName, version, endpoint string) *sdktrace.TracerProvider {
	//New otlp exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))

	if err != nil {
		log.Fatal(err)
	}
	//Resource 设置上报 Token，也可以直接配置环境变量来设置 token: OTEL_RESOURCE_ATTRIBUTES=token=xxxxxxxxx 如 config.yaml 里已配置，此处可忽略
	r, err := resource.New(ctx, []resource.Option{

		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(version),
			//attribute.KeyValue{Key: "token", Value: attribute.StringValue("{上报token}")},
		),
	}...)

	if err != nil {
		log.Fatal(err)
	}
	//创建一个新的TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(r),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
