package order

import (
	"errors"
	"github.com/theplant/luhn"
	"strconv"
)

func ValidateID(orderID string) error {
	var err error
	_, err = strconv.ParseUint(orderID, 10, 64)
	if err != nil {
		return errors.New("invalid orderID: invalid character set")
	}
	return err
}

func ValidateIDFormat(orderID string) error {
	order, _ := strconv.ParseUint(orderID, 10, 64)
	if luhn.Valid(int(order)) {
		return nil
	}
	return errors.New("id format is invalid")
}
