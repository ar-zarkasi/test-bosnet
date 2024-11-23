package controller

import (
	"app/src/constant"
	"app/src/http/request"
	"app/src/http/response"
	"app/src/services"
	"app/utils"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BosController struct {
	Balance services.BalanceService
	Counter services.CounterService
	History services.HistoryService
	Validate *validator.Validate
}

func NewBosController(balance services.BalanceService, counter services.CounterService, history services.HistoryService, validate *validator.Validate) *BosController {
	if validate == nil {
		utils.ErrorFatal(errors.New("validator instance cannot be nil"))
	}
	return &BosController{
		Balance: balance,
		Counter: counter,
		History: history,
		Validate: validate,
	}
}

func (controller *BosController) ListHistory(ctx *gin.Context) {
	id_account := ctx.Param("account")
	from_date := ctx.Query("from")
	to_date := ctx.Query("to")

	var (
		From time.Time
		To *time.Time
	)

	From, err := time.Parse(constant.FORMAT_DATE, from_date)
	if err != nil {
		utils.ErrorResponse(ctx, constant.ValidationError, "From parameter request must be valid format date")
		return
	}

	To = &From

	if to_date != "" {
		ToForm, err := time.Parse(constant.FORMAT_DATE, to_date)
		if err != nil {
			utils.ErrorResponse(ctx, constant.ValidationError, "To Parameter request must be valid format date")
			return
		}
		To = &ToForm
	}

	request := request.RequestHistoryList{
		Account: id_account,
		From: From,
		To: To,
	}

	respond, err := controller.History.HistoryList(request)
	if err != nil {
		utils.ErrorResponse(ctx, constant.ValidationError, err.Error())
		return
	}

	utils.Send(ctx, constant.Success, "History Listed", respond)
}

func (controller *BosController) Transfer(ctx *gin.Context) {
	req := request.RequestTransaction{}
	ctx.ShouldBind(&req)
	valid := controller.Validate.Struct(req)
	if valid != nil {
		utils.ErrorResponse(ctx, constant.ValidationError, valid.Error())
		return
	}
	// Parse to upper currency
	for idx, v := range req.Details {
		req.Details[idx].Currency = strings.ToUpper(v.Currency)
	}

	code, err := controller.History.TransactionTransfer(req)
	if err != nil {
		utils.ErrorResponse(ctx, code, err.Error())
		return
	}

	saldoAccount, code, err := controller.Balance.CekSaldo(req.FromAccount)
	if err != nil {
		utils.ErrorResponse(ctx, code, err.Error())
		return
	}

	saldoTargets := make([]response.ResponseAllSaldo,0)
	for _, w := range req.Details {
		saldo, code, err := controller.Balance.CekSaldo(w.ToAccount)
		if err != nil {
			utils.ErrorResponse(ctx, code, err.Error())
			return;
		}
		saldoTargets = append(saldoTargets, *saldo)
	}

	data := map[string]interface{}{
		"account_source": saldoAccount,
		"account_transferred": saldoTargets,
	}
	
	utils.Send(ctx, constant.Success, "Transfer Completed", data)
}

func (controller *BosController) transactionExecute(req request.RequestSelfTransaction, state string) (int, error) {
	// Parse to upper currency
	for idx, v := range req.Details {
		req.Details[idx].Currency = strings.ToUpper(v.Currency)
	}

	code, err := controller.History.GenericTransaction(req, state)
	if err != nil {
		return code, err
	}

	return constant.Success, nil
}

func (controller *BosController) Tarik(ctx *gin.Context) {
	req := request.RequestSelfTransaction{}
	ctx.ShouldBind(&req)
	valid := controller.Validate.Struct(req)
	if valid != nil {
		utils.ErrorResponse(ctx, constant.ValidationError, valid.Error())
		return
	}

	code, err := controller.transactionExecute(req, constant.TYPE_TRANSACTION_TARIK)
	if err != nil {
		utils.ErrorResponse(ctx, code, err.Error())
		return
	}

	saldoAccount, code, err := controller.Balance.CekSaldo(req.Account)
	if err != nil {
		utils.ErrorResponse(ctx, code, err.Error())
		return
	}

	utils.Send(ctx, code, "Success Tarik Dana", saldoAccount)
}

func (controller *BosController) Setor(ctx *gin.Context) {
	req := request.RequestSelfTransaction{}
	ctx.ShouldBind(&req)
	valid := controller.Validate.Struct(req)
	if valid != nil {
		utils.ErrorResponse(ctx, constant.ValidationError, valid.Error())
		return
	}

	code, err := controller.transactionExecute(req, constant.TYPE_TRANSACTION_SETOR)
	if err != nil {
		utils.ErrorResponse(ctx, code, err.Error())
		return
	}

	saldoAccount, code, err := controller.Balance.CekSaldo(req.Account)
	if err != nil {
		utils.ErrorResponse(ctx, code, err.Error())
		return
	}

	utils.Send(ctx, code, "Success Setor Dana", saldoAccount)
}