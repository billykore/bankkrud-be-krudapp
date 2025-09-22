package adapter

import (
	"github.com/google/wire"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/api"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/http/handler"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/http/server"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/service"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/storage/repo"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/account"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/cbs"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/payment"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transaction"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transfer"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
)

var ProviderSet = wire.NewSet(
	api.NewCBS,
	api.NewAccountAPI, wire.Bind(new(account.Service), new(*api.AccountAPI)),
	api.NewCBSStatusAPI, wire.Bind(new(cbs.Service), new(*api.CBSStatusAPI)),
	api.NewTransferAPI, wire.Bind(new(transfer.Service), new(*api.TransferAPI)),
	api.NewPaymentGateway, wire.Bind(new(payment.Service), new(*api.PaymentGateway)),
	repo.NewTransactionRepo, wire.Bind(new(transaction.Repository), new(*repo.TransactionRepo)),
	repo.NewUserRepo, wire.Bind(new(user.Repository), new(*repo.UserRepo)),
	service.NewAuthService, wire.Bind(new(user.AuthService), new(*service.AuthService)),
	handler.NewTransferHandler,
	handler.NewTapMoneyHandler,
	handler.NewAuthenticationHandler,
	handler.NewUserHandler,
	server.NewHTTP,
)
