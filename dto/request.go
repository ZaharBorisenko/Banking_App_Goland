package dto

import (
	"github.com/ZaharBorisenko/Banking_App_Goland/models"
	"github.com/google/uuid"
	"math/rand"
)

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func NewAccount(firstName string, lastName string) *models.Account {
	return &models.Account{
		ID:        uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(10000000)),
		Balance:   0,
	}
}
