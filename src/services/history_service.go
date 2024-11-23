package services

import (
	"app/src/constant"
	"app/src/http/request"
	"app/src/http/response"
	interfaces "app/src/interface"
	"app/src/models"
	"app/utils"
	"errors"
	"fmt"
	"time"
)

type HistoryService struct {
	Repo interfaces.HistoryInterface
	Balance BalanceService
	Counter CounterService
}

func NewHistoryService(HistoryRepository interfaces.HistoryInterface, balanceService *BalanceService, counterService *CounterService) (*HistoryService) {
	return &HistoryService{Repo: HistoryRepository,Balance: *balanceService, Counter: *counterService}
}

func (s *HistoryService) CreateHistory(req request.DataRequestHistory, continuesly *bool) (int, error) {
	tx := s.Repo.BeginTransaction()
	akun := s.Balance.CekAccount(req.Account, req.Currency)
	counterModel, err := s.Counter.FindOne(req.CounterId)
	if err != nil {
		fmt.Println("Counter Not Initiated, Create a new one")
		mdl, err := s.Counter.InitCounter()
		if err != nil {
			tx.Rollback()
			return constant.InternalServerError, err
		}
		counterModel = mdl
	}

	dateNow := time.Now()
	if req.DateTransaction != nil {
		dateNow = *req.DateTransaction
	}
	LastNumber := int(counterModel.ILastNumber)
	trx_last_string := fmt.Sprintf("%0*d", 5, LastNumber)
	total := s.Repo.CountTransaction(trx_last_string)
	mustDouble := utils.IsInArithmeticSequence(LastNumber)
	checkNumDouble := mustDouble && total > 1
	checkNumSingle := !mustDouble && total == 1
	fromRetry := false
	if continuesly != nil {
		fromRetry = *continuesly
	}
	if checkNumDouble || checkNumSingle || fromRetry {
		LastNumber++;
		err := s.Counter.UpdateCounter(LastNumber, counterModel)
		if err != nil {
			tx.Rollback()
			return constant.InternalServerError, err
		}
	}

	trx_id, err := s.generateTransactionNumber(LastNumber, dateNow)
	if err != nil {
		tx.Rollback()
		return constant.InternalServerError, err
	}

	model := models.History{
		SzTransactionId: *trx_id,
		SzAccountId: akun.SzAccountId,
		SzCurrencyId: akun.SzCurrencyId,
		DtmTransaction: dateNow,
		DecAmount: req.Amount,
		SzNote: req.TypeTransaction,
	}

	err = s.Repo.Create(&model)
	if err != nil {
		tx.Rollback()
		retry := true
		return s.CreateHistory(req, &retry)
		// return constant.BadRequest, err
	}

	tx.Commit()
	// go s.CalculateUpdateSaldo(model.SzAccountId, model.SzCurrencyId)
	s.CalculateUpdateSaldo(model.SzAccountId, model.SzCurrencyId)
	return constant.Success, nil
}

func (s *HistoryService) generateTransactionNumber(lastNumber int, Times time.Time) (*string, error) {
	n := 5;

	strPart := fmt.Sprintf("%0*d", n, 0)
	anotherStr := fmt.Sprintf("%0*d", n, lastNumber)
	timeStr := utils.DateToStringFormat(Times, constant.FORMAT_DATE2)
	
	result := timeStr + "-" + strPart + "." + anotherStr

	return &result, nil
}

func (s *HistoryService) CalculateUpdateSaldo(Account string, Currency string) {
	filter := map[string]interface{}{
		"szAccountId": Account,
		"szCurrencyId": Currency,
	}
	saldoAkhir := float64(0)

	transactions, err := s.Repo.FindTransaction(&filter, "dtmTransaction", false)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range transactions {
		if v.SzNote == constant.TYPE_TRANSACTION_SETOR || v.SzNote == constant.TYPE_TRANSACTION_TRANSFER {
			saldoAkhir += v.DecAmount
		} else if v.SzNote == constant.TYPE_TRANSACTION_TARIK {
			saldoAkhir -= v.DecAmount
		}
	}

	code, err := s.Balance.Update(Account, Currency, saldoAkhir)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Balance Update SuccessFully", Account, Currency, code)
	return
}

func (s *HistoryService) HistoryList(req request.RequestHistoryList) ([]response.ResponseHistoryList, error) {

	fromR := time.Date(req.From.Year(),req.From.Month(),req.From.Day(), 0, 0, 0, 0, req.From.Location());
	toR := time.Date(req.To.Year(), req.To.Month(), req.To.Day(), 23,59,59,999999999, req.To.Location());
	filter := map[string]interface{}{
		"szAccountId": req.Account,
		"dateBetween": map[string]string{
			"from": fromR.Format(constant.FORMAT_DATETIME),
			"to": toR.Format(constant.FORMAT_DATETIME),
		},
	}
	by := "dtmTransaction"
	descending := true
	lists, err := s.Repo.FindTransaction(&filter, by, descending)
	if err != nil {
		return nil, err
	}
	resp := make([]response.ResponseHistoryList,0)
	for _,v := range lists {
		temp := response.ResponseHistoryList{
			TransactionId: v.SzTransactionId,
			Account: v.SzAccountId,
			Currency: v.SzCurrencyId,
			TypeTransaction: v.SzNote,
			TransactionDate: v.DtmTransaction.Format(constant.FORMAT_DATETIME),
			Amount: v.DecAmount,
		}
		resp = append(resp, temp)
	}

	return resp, nil
}

func (s *HistoryService) TransactionTransfer(req request.RequestTransaction) (int, error) {
	typeTransaction := constant.TYPE_TRANSACTION_TRANSFER
	counter, _ := s.Counter.InitCounter()
	TransferError := make([]request.DataRequestHistory,0)
	TransferPrepareDetail := make([]request.PrepareRequestTransactionDetail,0)

	for _, tprepare := range req.Details {
		_ = s.Balance.CekAccount(req.FromAccount, tprepare.Currency)
		TransferPrepareDetail = append(TransferPrepareDetail, request.PrepareRequestTransactionDetail{
			Transaction: typeTransaction,
			Amount: -tprepare.Amount,
			Currency: tprepare.Currency,
		})
	}
	TransferPrepare := request.PrepareRequestTransaction{
		Account: req.FromAccount,
		Detail: TransferPrepareDetail,
	}

	// cek minimum saldo
	_, err := s.CekMinimumAfterTransaction(TransferPrepare)
	if err != nil {
		return constant.BadRequest, err
	}

	transaction := s.Repo.BeginTransaction()
	errorCount := 0
	for _, trx := range req.Details {
		form := request.DataRequestHistory{
			CounterId: counter.SzCounterId,
			Account: req.FromAccount,
			Currency: trx.Currency,
			Amount: -trx.Amount,
			TypeTransaction: typeTransaction,
		}
		_, err := s.CreateHistory(form, nil)
		if err != nil {
			errorCount++
			TransferError = append(TransferError, form)
			continue
		}
		
		form.Account = trx.ToAccount
		form.Amount = trx.Amount
		_, err = s.CreateHistory(form, nil)
		if err != nil {
			errorCount++
			TransferError = append(TransferError, form)
		}
	}
	if errorCount > 0 {
		transaction.Rollback()
		countError := len(TransferError)
		return constant.ServiceBroken, errors.New("Failed "+string(countError)+" Transaction")
	}

	transaction.Commit()
	return constant.Success, nil
}

func (s *HistoryService) CekMinimumAfterTransaction(data request.PrepareRequestTransaction) (string, error) {
	wallet, _, _ := s.Balance.CekSaldo(data.Account)
	for idx, w := range wallet.Wallet {
		for _, v := range data.Detail {
			if w.Currency == v.Currency {
				if v.Transaction == constant.TYPE_TRANSACTION_SETOR || v.Transaction == constant.TYPE_TRANSACTION_TRANSFER {
					wallet.Wallet[idx].Balance = w.Balance + v.Amount
				} else if v.Transaction == constant.TYPE_TRANSACTION_TARIK {
					wallet.Wallet[idx].Balance = w.Balance - v.Amount
				}
			}
		}
	}
	// pengecekan akhir
	for _, w := range wallet.Wallet {
		if w.Balance < 0 {
			return w.Currency, errors.New("Account "+w.Currency+", Insufficient Balance")
		}
	}

	fmt.Println("After Check", wallet.Wallet)

	return "", nil
}

func (s *HistoryService) GenericTransaction(req request.RequestSelfTransaction, typeTransaction string) (int, error) {
	if typeTransaction == constant.TYPE_TRANSACTION_TRANSFER {
		return constant.BadRequest, errors.ErrUnsupported
	}
	counter, _ := s.Counter.InitCounter()

	TarikError := make([]request.DataRequestHistory,0)
	TarikPrepareDetail := make([]request.PrepareRequestTransactionDetail,0)

	for _, tprepare := range req.Details {
		_ = s.Balance.CekAccount(req.Account, tprepare.Currency)
		TarikPrepareDetail = append(TarikPrepareDetail, request.PrepareRequestTransactionDetail{
			Transaction: typeTransaction,
			Amount: tprepare.Amount,
			Currency: tprepare.Currency,
		})
	}
	TarikPrepare := request.PrepareRequestTransaction{
		Account: req.Account,
		Detail: TarikPrepareDetail,
	}

	// cek minimum saldo
	_, err := s.CekMinimumAfterTransaction(TarikPrepare)
	if err != nil {
		return constant.BadRequest, err
	}

	transaction := s.Repo.BeginTransaction()
	errorCount := 0
	for _, trx := range req.Details {
		form := request.DataRequestHistory{
			CounterId: counter.SzCounterId,
			Account: req.Account,
			Currency: trx.Currency,
			Amount: trx.Amount,
			TypeTransaction: typeTransaction,
		}
		_, err := s.CreateHistory(form, nil)
		if err != nil {
			errorCount++
			TarikError = append(TarikError, form)
		}
	}
	if errorCount > 0 {
		transaction.Rollback()
		countError := len(TarikError)
		return constant.ServiceBroken, errors.New("Failed "+string(countError)+" Transaction")
	}

	transaction.Commit()
	return constant.Success, nil
}