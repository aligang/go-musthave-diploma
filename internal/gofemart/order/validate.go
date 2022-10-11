package order

import (
	"errors"
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
	var err error
	digitSequence, err := strconv.ParseUint(orderId, 10, 64)
	if !CheckLuhn(digitSequence) {
		return errors.New("invalid orderId: invalid checksum")

	}
	return err
}
