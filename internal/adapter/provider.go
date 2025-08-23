package adapter

import (
	"github.com/google/wire"
	"go.bankkrud.com/backend/svc/tapmoney/internal/adapter/api"
	"go.bankkrud.com/backend/svc/tapmoney/internal/adapter/http/handler"
	"go.bankkrud.com/backend/svc/tapmoney/internal/adapter/http/server"
	"go.bankkrud.com/backend/svc/tapmoney/internal/adapter/log"
	"go.bankkrud.com/backend/svc/tapmoney/internal/adapter/storage/repo"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/account"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/cbs"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/payment"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/pocket"
	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/transaction"
	pkglog "go.bankkrud.com/backend/svc/tapmoney/internal/pkg/log"
)

var ProviderSet = wire.NewSet(
	log.NewZap, wire.Bind(new(pkglog.Logger), new(*log.ZapLogger)),
	api.NewCBS,
	api.NewAccountAPI, wire.Bind(new(account.Repository), new(*api.AccountAPI)),
	api.NewCBSStatusAPI, wire.Bind(new(cbs.Service), new(*api.CBSStatusAPI)),
	api.NewPaymentGateway, wire.Bind(new(payment.Service), new(*api.PaymentGateway)),
	repo.NewTransactionRepo, wire.Bind(new(transaction.Repository), new(*repo.TransactionRepo)),
	repo.NewPocketRepo, wire.Bind(new(pocket.Repository), new(*repo.PocketRepo)),
	handler.NewTapMoneyHandler,
	server.New,
)
