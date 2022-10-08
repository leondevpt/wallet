package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	pb "github.com/leondevpt/wallet/trxservice/api/v1"
	"github.com/leondevpt/wallet/trxservice/pkg/trace"

	"github.com/joho/godotenv"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Printf("load env file err:%v\n", err)
	}

	auth := Auth{
		AppKey:    "test-appkey",
		AppSecret: "test-appsecret",
	}

	view.RegisterExporter(&exporter.PrintExporter{})
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}

	tp := trace.Init(context.Background(), "trxservice-client", "1.0.0", "http://127.0.0.1:14268/api/traces")

	defer tp.Shutdown(context.Background())

	//ctx := context.Background()
	// grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name))
	opts := []grpc.DialOption{grpc.WithPerRPCCredentials(&auth), grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		grpc.WithUnaryInterceptor(trace.NewGRPUnaryClientInterceptor())}

	// grpc.UseCompressor(gzip.Name) 为了更好地节省带宽，在rpc调用的客户端和服务端都开启压缩
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	serverAddr := GetEnv("SERVER_ADDR", "localhost:50051")

	//conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("did not connect:%v", err)
	}
	defer conn.Close()

	cli := pb.NewTrxServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	addr := "TMxN1hCZMwcsyQn7CBHizJ7TNDT4nrmpdh"

	trxBalance, err := cli.GetTrxBalance(ctx, &pb.GetTrxBalanceRequest{Address: addr})
	if err != nil {
		log.Fatalf("Could not get trx balance: %v", err)
	}
	log.Printf("TRX balance:%s\n", trxBalance.Balance)

	newCtx, cancel2 := context.WithTimeout(context.Background(), time.Second)
	defer cancel2()

	usdtBalance, err := cli.GetTRC20TokenBalance(newCtx, &pb.GetTRC20TokenBalanceRequest{Address: addr, Token: "USDT"})

	if err != nil {
		log.Fatalf("Could not get trc20 usdt balance: %v", err)
	}
	log.Printf("USDT balance: %s\n", usdtBalance.Balance)

}

func GetEnv(name string, defaultValue string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return defaultValue
}
