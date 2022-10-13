package order

import (
	"errors"
	"github.com/theplant/luhn"
	"strconv"
)

func ValidateId(orderId string) error {
	var err error
	_, err = strconv.ParseUint(orderId, 10, 64)
	if err != nil {
		errors.New("invalid orderId: invalid character set")
	}
	return err
}

func ValidateIdFormat(orderId string) error {
	order, _ := strconv.ParseUint(orderId, 10, 32)
	if luhn.Valid(int(order)) {
		return nil
	}
	return errors.New("id format is invalid")
}
