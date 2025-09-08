package adapter

import (
	"github.com/google/wire"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/api"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/http/handler"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/http/server"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/storage/repo"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/account"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/cbs"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/payment"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transaction"
)

var ProviderSet = wire.NewSet(
	api.NewCBS,
	api.NewAccountAPI, wire.Bind(new(account.Service), new(*api.AccountAPI)),
	api.NewCBSStatusAPI, wire.Bind(new(cbs.Service), new(*api.CBSStatusAPI)),
	api.NewPaymentGateway, wire.Bind(new(payment.Service), new(*api.PaymentGateway)),
	repo.NewTransactionRepo, wire.Bind(new(transaction.Repository), new(*repo.TransactionRepo)),
	handler.NewTapMoneyHandler,
	server.NewHTTP,
)
