package service

import (
	"context"
)

type Auth struct{}

func (a *Auth) GetAppKey() string {
	return "test-appkey"
}

func (a *Auth) GetAppSecret() string {
	return "test-appsecret"
}

func (a *Auth) Check(ctx context.Context) error {
	/*
	md, _ := metadata.FromIncomingContext(ctx)

	var appKey, appSecret string
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}
	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return errcode.TogRPCError(errcode.Unauthorized)
	}
*/
	return nil
}
