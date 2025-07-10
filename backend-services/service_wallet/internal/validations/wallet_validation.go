package validation

import (
	"fmt"
	"service-wallet/internal/models/request"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func ValidationCreateWallet(params request.CreateWallet) (errRedaksi, errSystem error) {
	err := validation.ValidateStruct(&params,
		validation.Field(&params.ReqID, validation.Required),
		validation.Field(&params.Name, validation.Required),
		validation.Field(&params.Type, validation.Required),
		validation.Field(&params.UserId, validation.Required),
	)
	if err != nil {
		return fmt.Errorf("Data 'ReqId', 'Name', 'Type', 'UserId' must fill."), fmt.Errorf("failed validation struct request kuota cabang. err : %v", err)
	}

	if params.Balance == "" {
		params.Balance = "0"
	} else if _, err = strconv.Atoi(params.Balance); (params.Balance != "") && (err != nil) {
		return fmt.Errorf("Data 'Balance' : '%s' is not value.", params.Balance), fmt.Errorf("failed convert str to int of params.Balance(%s). err : %v", params.Balance, err)
	}

	if _, err := uuid.Parse(params.UserId); err != nil {
		return fmt.Errorf("user id('%s') is not valid system", params.UserId), fmt.Errorf("failed parse uuid of params.UserId '%s'. err : %v", params.UserId, err)
	}

	return nil, nil
}
