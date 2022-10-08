package biz

import (
	"context"
	"math/big"

	"go.uber.org/zap"
)

// TrxRepo is a trx repo.
type TrxRepo interface {
	ListTransactions(ctx context.Context, addr, token string, pageNum, pageSize, queryType int) ([]interface{}, int, error)
}

type TrxUsecase struct {
	repo TrxRepo
	log  *zap.Logger
	cli  *TronCli
}

// NewTrxUsecase new a Trx usecase.
func NewTrxUsecase(repo TrxRepo, logger *zap.Logger, cli *TronCli) *TrxUsecase {
	return &TrxUsecase{repo: repo, log: logger, cli: cli}
}

func (t *TrxUsecase) GetBalance(ctx context.Context, addr string) (*big.Int, error) {
	return t.cli.GetBalance(ctx, addr)
}

func (t *TrxUsecase) GetTRC20TokenBalance(ctx context.Context, addr, contractAddr string) (*big.Int, error) {
	return t.cli.GetTRC20TokenBalance(ctx, addr, contractAddr)
}
