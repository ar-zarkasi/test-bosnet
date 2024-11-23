package services

import (
	"app/src/http/request"
	"app/src/http/response"
	interfaces "app/src/interface"
	"app/src/models"
	"errors"
	"fmt"
)

type CounterService struct {
	Repo interfaces.CounterInterface
}

func NewCounterService(CounterRepository interfaces.CounterInterface) (*CounterService) {
	return &CounterService{
		Repo: CounterRepository,
	}
}

func (s *CounterService) Create(req request.DataCounter) (*string, error) {
	find := map[string]interface{} {
		"szCounterId": req.CounterId,
	}
	check := s.Repo.Find(&find)
	if check != nil && len(check) > 0 {
		fmt.Println("Counter Exists", check)
		return &check[0].SzCounterId, nil
	}
	var model models.Counter
	err := s.Repo.Create(req, &model)
	if err != nil {
		return nil, err
	}

	return &model.SzCounterId, nil
}

func (s *CounterService) UpdateCounter(Count int, model *models.Counter) error {
	model.ILastNumber = uint(Count)
	err := s.Repo.Update(model)
	return err
}

func (s *CounterService) Lists() ([]response.ResponseCounter) {
	models := s.Repo.Find(nil)
	list_response := make([]response.ResponseCounter, 0)
	for _, v := range models {
		var resp response.ResponseCounter
		resp.CounterId = v.SzCounterId
		resp.LastNumber = v.ILastNumber
		list_response = append(list_response, resp)
	}
	return list_response
}

func (s *CounterService) FindOne(counterId string) (*models.Counter, error) {
	filter := map[string]interface{}{
		"szCounterId": counterId,
	}
	models := s.Repo.Find(&filter)
	if len(models) == 0 {
		return nil, errors.New("Counter Not Found")
	}

	return &models[0], nil
}

func (s *CounterService) InitCounter() (*models.Counter, error) {
	data := request.DataCounter{
		CounterId: "001-COU",
		LastNumber: 1,
	}
	var model models.Counter
	check, _ := s.FindOne(data.CounterId)
	if check != nil {
		return check, nil
	}
	err := s.Repo.Create(data, &model)
	if err != nil {
		return nil, err
	}
	return &model, nil
}