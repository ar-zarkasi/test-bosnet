package services

import (
	"app/src/constant"
	"app/src/http/request"
	interfaces "app/src/interface"
	"app/src/models"
	"app/utils"
	"fmt"
	"time"
)

type SettingService struct {
	Repo interfaces.SettingInterface
	Counter *CounterService
	History *HistoryService
}

func NewSettings(settingRepository interfaces.SettingInterface, counter CounterService, history HistoryService) (*SettingService) {
	return &SettingService{Repo: settingRepository, Counter: &counter, History: &history}
}

func (s *SettingService) InitDeploy() {
	key := "initialize"
	value := "true"

	model := s.Repo.Find(key)
	if model == nil {
		trx := s.Repo.BeginTransaction()
		fmt.Println("Begin Initialize")

		counter, err := s.Counter.InitCounter()
		if err != nil {
			trx.Rollback()
			fmt.Println("Initialize Counter Failed")
			utils.ErrorFatal(err)
			return
		}

		err = s.initHistory(*counter)
		if err != nil {
			trx.Rollback()
			fmt.Println("Initialize History Failed")
			utils.ErrorFatal(err)
			return
		}

		_, err = s.Repo.Create(key, value)
		if err != nil {
			trx.Rollback()
			fmt.Println("Initialize Failed")
			utils.ErrorFatal(err)
			return
		}

		trx.Commit()
	}

	fmt.Println("Has Initialized")
}

func (s *SettingService) initHistory(counter models.Counter) error {
	multipleData := make([]request.DataRequestHistory,0)
	account1 := "000108757484"
	account2 := "000109999999"
	account3 := "000108888888"
	counterId := counter.SzCounterId
	date := time.Date(2020, 12, 31, 16, 34, 0, 0, time.Local)

	data1 := request.DataRequestHistory{
		CounterId: counterId,
		Account: account1,
		Currency: constant.CURR_IDR,
		TypeTransaction: constant.TYPE_TRANSACTION_SETOR,
		Amount: 34500000.00,
		DateTransaction: &date,
	}
	data2 := request.DataRequestHistory{
		CounterId: counterId,
		Account: account1,
		Currency: constant.CURR_SGD,
		TypeTransaction: constant.TYPE_TRANSACTION_SETOR,
		Amount: 125.8750,
		DateTransaction: &date,
	}
	data3 := request.DataRequestHistory{
		CounterId: counterId,
		Account: account2,
		Currency: constant.CURR_IDR,
		TypeTransaction: constant.TYPE_TRANSACTION_SETOR,
		Amount: 1250.00,
		DateTransaction: &date,
	}
	data4 := request.DataRequestHistory{
		CounterId: counterId,
		Account: account2,
		Currency: constant.CURR_SGD,
		TypeTransaction: constant.TYPE_TRANSACTION_SETOR,
		Amount: 128.00,
		DateTransaction: &date,
	}
	data5 := request.DataRequestHistory{
		CounterId: counterId,
		Account: account2,
		Currency: constant.CURR_SGD,
		TypeTransaction: constant.TYPE_TRANSACTION_TRANSFER,
		Amount: -125.75,
		DateTransaction: &date,
	}
	data6 := request.DataRequestHistory{
		CounterId: counterId,
		Account: account3,
		Currency: constant.CURR_SGD,
		TypeTransaction: constant.TYPE_TRANSACTION_TRANSFER,
		Amount: 125.75,
		DateTransaction: &date,
	}
	multipleData = append(multipleData, data1, data2, data3, data4, data5, data6)
	for idx, v := range multipleData {
		_, err := s.History.CreateHistory(v, nil)
		if err != nil {
			fmt.Println("Error Create Initialize History",idx, err)
			return err
		}
	}
	return nil
}