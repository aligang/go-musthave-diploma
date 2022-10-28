package account

import (
	"errors"
	"fmt"
)

func ValidateCredentials(account *CustomerAccount) error {
	if len(account.Login) == 0 {
		return errors.New("invalid input data; login may not be empty")
	}
	if len(account.Password) == 0 {
		return fmt.Errorf("invalid input data; password may not be empty")
	}
	return nil
}
