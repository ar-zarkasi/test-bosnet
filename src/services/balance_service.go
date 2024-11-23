package services

import (
	"app/src/constant"
	"app/src/http/request"
	"app/src/http/response"
	interfaces "app/src/interface"
	"app/src/models"
	"fmt"
)

type BalanceService struct {
	Repo interfaces.BalanceInterface
}

func NewBalanceService(BalanceRepository interfaces.BalanceInterface) (*BalanceService) {
	return &BalanceService{Repo: BalanceRepository}
}

func (s *BalanceService) Create(Account string, Currency string, Amount *float64) (*models.Balance, int, error) {
	var (
		model models.Balance
		saldo float64
	)
	find, _ := s.Repo.FindOne(Account, Currency)
	if find != nil {
		return find, constant.Success,nil
	}
	saldo = 0 // init default saldo
	if Amount != nil {
		saldo = *Amount
	}
	data := request.DataBalance{
		Account: Account,
		Currency: Currency,
		Balance: saldo,
	}
	err := s.Repo.Create(data, &model)
	if err != nil {
		return nil, constant.InternalServerError, err
	}

	return &model, constant.SuccessCreate, nil
}

func (s *BalanceService) Update(Account string, Currency string, Amount float64) (int, error) {
	model, err := s.Repo.FindOne(Account, Currency)
	if err != nil {
		data := request.DataBalance{
			Account: Account,
			Currency: Currency,
			Balance: Amount,
		}
		err = s.Repo.Create(data, model)
		if err != nil {
			return constant.InternalServerError, err
		}
		return constant.SuccessCreate, nil
	}

	err = s.Repo.UpdateSaldo(Amount, model)
	if err != nil {
		return constant.InternalServerError, err
	}

	return constant.Success, nil
}

func (s *BalanceService) CekSaldo(Account string) (*response.ResponseAllSaldo, int, error){
	allSaldo := make([]response.ResponseSaldoOnly,0)
	var resp response.ResponseAllSaldo
	resp.Account = Account
	models, err := s.Repo.FindAllByAccount(Account)
	if err != nil {
		return nil, constant.NotFound, err
	}
	for _, v := range models {
		saldo := response.ResponseSaldoOnly{
			Currency: v.SzCurrencyId,
			Balance: *v.DecAmount,
		}
		allSaldo = append(allSaldo, saldo)
	}
	resp.Wallet = allSaldo
	return &resp, constant.Success, nil
}

func (s *BalanceService) CekAccount(Account string, Currency string) *models.Balance {
	model, err := s.Repo.FindOne(Account, Currency)
	if err != nil {
		fmt.Println("Error Load Model Balance", err)
		saldoawal := float64(0)
		m, _, _ := s.Create(Account, Currency, &saldoawal)
		model = m
		return model
	}
	return model
}