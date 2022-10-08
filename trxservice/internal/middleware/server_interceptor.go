package middleware

import (
	"context"
	"runtime/debug"
	"time"
	"github.com/leondevpt/wallet/trxservice/pkg/errcode"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//requestLog := "access request log: method: %s, begin_time: %d, request: %v"
	//log.Printf(requestLog, info.FullMethod, beginTime, req)
	beginTime := time.Now()
	beginTimeUnix := beginTime.Local().Unix()
	zap.S().Infof("access request log: method: %s, begin_time: %d, request: %v",
		info.FullMethod, beginTimeUnix, req)

	resp, err := handler(ctx, req)

	endTimeUnix := time.Now().Local().Unix()
	//responseLog := "access response log: method: %s, begin_time: %d, end_time: %d, cost:%s,response: %v"
	//log.Printf(responseLog, info.FullMethod, beginTimeUnix, endTimeUnix, time.Since(beginTime), resp)
	zap.S().Infof("access response log: method: %s, begin_time: %d, end_time: %d, cost:%s,response: %v",
		info.FullMethod, beginTimeUnix, endTimeUnix, time.Since(beginTime), resp)
	return resp, err
}

// 普通错误记录的日志拦截器
func ErrorLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		s := errcode.FromError(err)
		//errLog := "error log: method: %s, code: %v, message: %v, details: %v"
		//log.Printf(errLog, info.FullMethod, s.Code(), s.Err().Error(), s.Details())
		zap.S().Infof("error log: method: %s, code: %v, message: %v, details: %v", info.FullMethod, s.Code(), s.Err().Error(), s.Details())
	}
	return resp, err
}

// 异常捕抓拦截器
func Recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		if e := recover(); e != nil {
			//recoveryLog := "recovery log: method: %s, message: %v, stack: %s"
			//log.Printf(recoveryLog, info.FullMethod, e, string(debug.Stack()[:]))
			zap.S().Info("recovery log: method: %s, message: %v, stack: %s", info.FullMethod, e, string(debug.Stack()[:]))
		}
	}()

	return handler(ctx, req)
}

// ServerInterceptor Server 端的拦截器
func ServerInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		spanContext, err := tracer.Extract(
			opentracing.TextMap,
			MDCarrier{md},
		)

		if err != nil && err != opentracing.ErrSpanContextNotFound {
			grpclog.Errorf("extract from metadata err: %v", err)
		} else {
			span := tracer.StartSpan(
				info.FullMethod,
				ext.RPCServerOption(spanContext),
				opentracing.Tag{Key: string(ext.Component), Value: "gRPC Server"},
				ext.SpanKindRPCServer,
			)
			defer span.Finish()

			ctx = opentracing.ContextWithSpan(ctx, span)
		}

		return handler(ctx, req)

	}

}

/*
func RequestIDClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string, req, resp interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
	) (err error) {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}

		value := ctx.Value(trace.RequestID)
		if requestID, ok := value.(string); ok && requestID != "" {
			md[string(trace.RequestID)] = []string{requestID}
		}
		return invoker(metadata.NewOutgoingContext(ctx, md), method, req, resp, cc, opts...)
	}
}

func RequestIDServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}
		// Set request ID for context.
		requestIDs := md[string(trace.RequestID)]
		if len(requestIDs) >= 1 {
			ctx = context.WithValue(ctx, trace.RequestID, requestIDs[0])
			return handler(ctx, req)
		}

		// Generate request ID and set context if not exists.
		requestID := id.NewHex32()
		ctx = context.WithValue(ctx, trace.RequestID, requestID)
		return handler(ctx, req)
	}
}
*/
