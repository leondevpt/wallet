package server

import (
	"context"
	//"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	//"go.opencensus.io/stats/view"
	//"log"
	"math"
	"net"
	"time"
	pb "github.com/leondevpt/wallet/trxservice/api/v1"
	"github.com/leondevpt/wallet/trxservice/internal/middleware"
	"github.com/leondevpt/wallet/trxservice/version"

	//"trxservice/pkg/metric"
	"github.com/leondevpt/wallet/trxservice/pkg/setting"
	"github.com/leondevpt/wallet/trxservice/pkg/trace"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	//"go.opencensus.io/plugin/ocgrpc"
	//"go.opencensus.io/stats/view"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	//"go.opencensus.io/examples/exporter"
)

// GrpcServer implements a gRPC Server for the Order service
type GrpcServer struct {
	Server        *grpc.Server
	errCh         chan error
	listener      net.Listener
	traceProvider *sdktrace.TracerProvider
}

// NewGrpcServer is a convenience func to create a GrpcServer
func NewGrpcServer(service pb.TrxServiceServer, cfg *setting.Config, zapLogger *zap.Logger) (*GrpcServer, error) {
	/*
		addr := fmt.Sprintf(":%d", cfg.App.GrpcPort)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			return nil, err
		}
		fmt.Println("Listening at ", addr)
	*/

	var tp *sdktrace.TracerProvider

	/*
	view.RegisterExporter(&exporter.PrintExporter{})
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Fatal(err)
	}

	 */

	/*
		metrics, err := metric.CreateMetrics(cfg.Metrics.URL, cfg.Metrics.ServiceName)
		if err != nil {
			zapLogger.Sugar().Errorf("CreateMetrics Error: %s", err)
		}
	*/

	zapLogger.Sugar().Infow(
		"Metrics available URL: %s, ServiceName: %s",
		cfg.Metrics.URL,
		cfg.Metrics.ServiceName,
	)

	// 设置一元拦截器
	interceptors := []grpc.UnaryServerInterceptor{}
	interceptors = append(interceptors,
		middleware.AccessLog, middleware.ErrorLog,
		middleware.Recovery,
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(zapLogger),
		grpc_auth.UnaryServerInterceptor(myAuthFunction),
		grpc_recovery.UnaryServerInterceptor(),
	)

	//grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	if cfg.Trace.Enable {
		tp = trace.Init(context.Background(), cfg.Trace.ServiceName, version.Version, cfg.Trace.Endpoint)
		interceptors = append(interceptors, trace.NewGRPUnaryServerInterceptor()) //设置trace拦截器进行埋点)
	}

	opts := []grpc.ServerOption{
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
		grpc.MaxRecvMsgSize(math.MaxInt64),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				/*
					When the connection reaches its max-age, it will be closed and will trigger a re-resolve from the client.
					If new instances were added in the meantime, the client will see them now and open connections to them as well.
				*/
				MaxConnectionAge: time.Minute * 5,
			},
		),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			interceptors...,
		)),
	}
	// 初始化grpc对象
	server := grpc.NewServer(opts...)
	// 注册服务
	pb.RegisterTrxServiceServer(server, service)
	reflection.Register(server)

	return &GrpcServer{
		Server:        server,
		errCh:         make(chan error, 1),
		traceProvider: tp,
	}, nil
}

// Start starts the server in the background, pushing any error to the error channel
func (g GrpcServer) Start() {
	go func() {
		if err := g.Server.Serve(g.listener); err != nil {
			g.errCh <- err
		}
	}()
}

// Stop stops the gRPC server
func (g GrpcServer) Stop() {
	if g.traceProvider != nil {
		g.traceProvider.Shutdown(context.Background())
	}
	g.Server.GracefulStop()
}

// Error returns the server's error channel
func (g GrpcServer) Error() chan error {
	return g.errCh
}

func myAuthFunction(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
