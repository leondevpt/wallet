//go:build wireinject
// +build wireinject

package main

import (
	"github.com/leondevpt/wallet/trxservice/internal/biz"
	"github.com/leondevpt/wallet/trxservice/internal/data"
	"github.com/leondevpt/wallet/trxservice/internal/server"
	"github.com/leondevpt/wallet/trxservice/internal/service"
	"github.com/leondevpt/wallet/trxservice/pkg/setting"

	"github.com/google/wire"
	"go.uber.org/zap"
)

/*
// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Node, *conf.Registry, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
*/

// wireApp init kratos application.
func wireApp(cfg *setting.Config, logger *zap.Logger) (app, error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
