package main

import (
	"context"
	"encoding/base64"

	"google.golang.org/grpc/credentials"
)

/*
 gRPC 默认提供用于自定义认证 Token 的接口，它的作用是将所需的安全认证信息添加到每个 RPC 方法的上下文中
*/
var _ credentials.PerRPCCredentials = &Auth{}

type Auth struct {
	AppKey    string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return false
}

type BasicAuth struct {
	UserName string
	Password string
}

// 实现GetRequestMetadata方法 将用户凭证转换成请求meta元数据
func (b BasicAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	auth := b.UserName + ":" + b.Password
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + encoded,
	}, nil
}

// Basic Authentication 不建议使用明文传输，要使用加密通信
func (b BasicAuth) RequireTransportSecurity() bool {
	return true
}
