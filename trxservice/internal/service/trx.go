package service

import (
	"context"
	"fmt"
	"strings"
	pb "github.com/leondevpt/wallet/trxservice/api/v1"
	"github.com/leondevpt/wallet/trxservice/internal/biz"
	"github.com/leondevpt/wallet/trxservice/pkg/setting"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var _ pb.TrxServiceServer = &TrxService{}

type TrxService struct {
	auth *Auth
	uc   *biz.TrxUsecase
	pb.UnimplementedTrxServiceServer
	log *zap.Logger
}

func NewTrxService(uc *biz.TrxUsecase, log *zap.Logger) pb.TrxServiceServer {
	return &TrxService{uc: uc, log: log, auth: &Auth{}}
}

func (s *TrxService) GetTrxBalance(c context.Context, req *pb.GetTrxBalanceRequest) (*pb.GetTrxBalanceReply, error) {
	if err := s.auth.Check(c); err != nil {
		return nil, err
	}
	bigBalance, err := s.uc.GetBalance(c, req.Address)
	if err != nil {
		s.log.Sugar().Errorw("GetTrxBalance", "addr", req.Address, "err", err)
		return nil, err
	}

	d := decimal.NewFromBigInt(bigBalance, 0)
	result := d.Div(decimal.New(1, 6))
	s.log.Sugar().Infow("GetTrxBalance", "addr", req.Address, "balance", result.String())
	return &pb.GetTrxBalanceReply{Balance: result.String()}, nil
}

func (s *TrxService) GetTRC20TokenBalance(c context.Context, req *pb.GetTRC20TokenBalanceRequest) (*pb.GetTRC20TokenBalanceReply, error) {
	if err := s.auth.Check(c); err != nil {
		return nil, err
	}
	ok, tokenInfo := checkTokenSupport(req.Token)
	if !ok {
		return nil, fmt.Errorf("token %s not support", req.Token)
	}

	bigBalance, err := s.uc.GetTRC20TokenBalance(c, req.Address, tokenInfo.ContractAddr)
	if err != nil {
		s.log.Sugar().Errorw("GetTRC20TokenBalance", "addr", req.Address, "token", req.Token, "err", err)
		return nil, err
	}

	d := decimal.NewFromBigInt(bigBalance, 0)

	result := d.Div(decimal.New(1, int32(tokenInfo.Decimal)))

	s.log.Sugar().Infow("GetTRC20TokenBalance", "addr", req.Address, "token", req.Token, "balance", result.String())

	return &pb.GetTRC20TokenBalanceReply{Token: req.Token, Balance: result.String()}, nil
}

func checkTokenSupport(tokenSymbol string) (bool, setting.Token) {
	for k, info := range setting.Conf.TokenList {
		if strings.ToUpper(k) == strings.ToUpper(tokenSymbol) {
			return true, info
		}
	}
	return false, setting.Token{}
}
