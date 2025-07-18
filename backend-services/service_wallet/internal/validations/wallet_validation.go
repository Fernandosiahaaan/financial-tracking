package validation

import (
	"fmt"
	"service-wallet/internal/models/request"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func ValidationCreateWallet(params request.CreateWalletRequest) (errRedaksi, errSystem error) {
	err := validation.ValidateStruct(&params,
		validation.Field(&params.ReqID, validation.Required),
		validation.Field(&params.Name, validation.Required),
		validation.Field(&params.Type, validation.Required),
		validation.Field(&params.UserId, validation.Required),
	)
	if err != nil {
		return fmt.Errorf("Data 'ReqId', 'Name', 'Type', 'UserId' must fill."), fmt.Errorf("failed validation struct request kuota cabang. err : %v", err)
	}

	errRedaksi, errSystem = ValidateStrInt("Balance", "0", &params.Balance)
	if errRedaksi != nil || errSystem != nil {
		return
	}

	return ValidateUUID("User ID", params.UserId)
}

func ValidationUpdateWallet(params request.UpdateWalletRequest) (errRedaksi, errSystem error) {
	err := validation.ValidateStruct(&params,
		validation.Field(&params.ReqID, validation.Required),
		validation.Field(&params.WalletID, validation.Required),
		validation.Field(&params.Name, validation.Required),
		validation.Field(&params.Type, validation.Required),
	)
	if err != nil {
		return fmt.Errorf("Data 'ReqId', 'Name', 'Type', 'UserId' must fill."), fmt.Errorf("failed validation struct request kuota cabang. err : %v", err)
	}

	errRedaksi, errSystem = ValidateStrInt("Balance", "0", &params.Balance)
	if errRedaksi != nil || errSystem != nil {
		return
	}

	return ValidateUUID("Wallet ID", params.WalletID)
}

func ValidationGetListWallet(params request.GetListWalletRequest) (errRedaksi, errSystem error) {
	err := validation.ValidateStruct(&params,
		validation.Field(&params.ReqID, validation.Required),
	)
	if err != nil {
		return fmt.Errorf("Data 'ReqId' must fill."), fmt.Errorf("failed validation struct request kuota cabang. err : %v", err)
	}

	errRedaksi, errSystem = ValidateStrInt("PageItem", "", &params.PageItem)
	if errRedaksi != nil || errSystem != nil {
		return
	}

	return ValidateStrInt("Page", "", &params.Page)
}

func ValidateStrInt(key, deafultValue string, value *string) (errRedaksi, errSystem error) {
	if *value == "" {
		*value = deafultValue
	} else if _, err := strconv.Atoi(*value); (*value != "") && (err != nil) {
		return fmt.Errorf("Data '%s' : '%s' is not value.", key, *value), fmt.Errorf("failed convert str to int of params.%s(%s). err : %v", key, *value, err)
	}

	return nil, nil
}

func ValidateUUID(key, id string) (errRedaksi, errSystem error) {

	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%s('%s') is not valid system", key, id), fmt.Errorf("failed parse uuid of WalletId '%s'. err : %v", id, err)
	}

	return nil, nil

}
