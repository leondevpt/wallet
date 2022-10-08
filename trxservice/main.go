package main

import (
	"context"
	"flag"
	"fmt"
	//"go.opencensus.io/zpages"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	pb "github.com/leondevpt/wallet/trxservice/api/v1"
	"github.com/leondevpt/wallet/trxservice/internal/logger"
	"github.com/leondevpt/wallet/trxservice/internal/server"
	"github.com/leondevpt/wallet/trxservice/internal/util"
	"github.com/leondevpt/wallet/trxservice/pkg/setting"
	"github.com/leondevpt/wallet/trxservice/version"

	"github.com/leondevpt/wallet/trxservice/global"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	//"go.opencensus.io/zpages"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

var (
	flagconf    string
	versionShow bool
)

func init() {
	flag.StringVar(&flagconf, "conf", "config.yaml", "config path, eg: -conf config.yaml")
	flag.BoolVar(&versionShow, "version", false, "show version info, eg: -version")
}

func main() {
	flag.Parse()

	if versionShow {
		version.ShowVersion()
		return
	}
	setting.Init()

	fmt.Printf("Cfg:%v\n", setting.Conf)
	// 注册自定义的 静态resolver
	builder := util.NewStaticResolverBuilder(map[string][]string{global.TronNode: setting.Conf.App.Node_Addr})
	resolver.Register(builder)

	lg := logger.NewZapLogger()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	app, err := wireApp(setting.Conf, lg)
	if err != nil {
		lg.Sugar().Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := run(ctx, app); err != nil {
			lg.Sugar().Fatal(err)
		}
	}()

	<-done
	cancel()

	// 等待5秒后退出
	c, newCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer newCancel()
	<-c.Done()
	fmt.Println("shutdown")
	close(done)
	fmt.Println("done")

}

// run starts the app, handling any REST or gRPC server error
// and as well as app shutdown
func run(ctx context.Context, app app) error {
	if err := app.start(); err != nil {
		return err
	}
	defer app.shutdown()

	select {
	case grpcErr := <-app.grpcServer.Error():
		return grpcErr
	case <-ctx.Done():
		return nil
	}
}

// app is a convenience wrapper for all things needed to start
// and shutdown the Order microservice
type app struct {
	grpcServer *server.GrpcServer
}

// start starts the REST and gRPC Servers in the background
func (a app) start() error {
	//a.grpcServer.Start() // also non blocking :-)
	// 给客户端连接注册TrxServiceHandler 和Mux转发到grpc的端口
	httpMux := runHttpServer()
	grpcS := a.grpcServer.Server
	gatewayMux := runGrpcGatewayServer()
	httpMux.Handle("/", gatewayMux)

	return http.ListenAndServe(fmt.Sprintf(":%d", setting.Conf.GrpcPort), grpcHandlerFunc(grpcS, httpMux))

}

// stop shuts down the servers
func (a app) shutdown() error {
	a.grpcServer.Stop()
	return nil
}

// newApp creates a new app with REST & gRPC servers
// this func performs all app related initialization
func newApp(gs *server.GrpcServer) (app, error) {
	return app{
		grpcServer: gs,
	}, nil
}

// grpcHandlerFunc 根据请求头判断是grpc请求还是grpc-gateway请求
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func runHttpServer() *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		hostName, _ := os.Hostname()
		_, _ = w.Write([]byte(`pong` + hostName))
	})
	//zpages.Handle(serveMux, "/debug")
	serveMux.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	return serveMux
}

func runGrpcGatewayServer() *gwruntime.ServeMux {
	endpoint := fmt.Sprintf("0.0.0.0:%d", setting.Conf.GrpcPort)

	gwmux := gwruntime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterTrxServiceHandlerFromEndpoint(context.Background(), gwmux, endpoint, opts)
	if err != nil {
		log.Fatal(err)
	}
	return gwmux
}

func httpServer() error {
	ctxr := context.Background()
	ctx, cancel := context.WithCancel(ctxr)
	defer cancel()
	mux := gwruntime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	//與原生gRPC不同點在這邊，需要做http與grpc的對應
	err := pb.RegisterTrxServiceHandlerFromEndpoint(ctx, mux, ":8787", opts)
	if err != nil {
		return err
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", setting.Conf.GrpcPort), mux)
}
