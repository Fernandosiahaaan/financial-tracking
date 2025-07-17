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

	return ValidateUUID("User ID", params.UserId)
}

func ValidationUpdateWallet(params request.UpdateWallet) (errRedaksi, errSystem error) {
	err := validation.ValidateStruct(&params,
		validation.Field(&params.ReqID, validation.Required),
		validation.Field(&params.WalletID, validation.Required),
		validation.Field(&params.Name, validation.Required),
		validation.Field(&params.Type, validation.Required),
	)
	if err != nil {
		return fmt.Errorf("Data 'ReqId', 'Name', 'Type', 'UserId' must fill."), fmt.Errorf("failed validation struct request kuota cabang. err : %v", err)
	}

	if params.Balance == "" {
		params.Balance = "0"
	} else if _, err = strconv.Atoi(params.Balance); (params.Balance != "") && (err != nil) {
		return fmt.Errorf("Data 'Balance' : '%s' is not value.", params.Balance), fmt.Errorf("failed convert str to int of params.Balance(%s). err : %v", params.Balance, err)
	}

	return ValidateUUID("Wallet ID", params.WalletID)
}

func ValidateUUID(key, id string) (errRedaksi, errSystem error) {

	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%s('%s') is not valid system", key, id), fmt.Errorf("failed parse uuid of WalletId '%s'. err : %v", id, err)
	}

	return nil, nil

}
