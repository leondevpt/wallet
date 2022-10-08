package data

import (
	"context"
	"github.com/leondevpt/wallet/trxservice/internal/biz"

	"go.uber.org/zap"
)

type trxrRepo struct {
	data *Data
	log  *zap.Logger
}

// NewTrxRepo .
func NewTrxRepo(data *Data, logger *zap.Logger) biz.TrxRepo {
	return &trxrRepo{
		data: data,
		log:  logger,
	}
}

func (trxrRepo) ListTransactions(ctx context.Context, addr, token string, pageNum, pageSize, queryType int) ([]interface{}, int, error) {
	return nil, 0, nil
}
