package validator

import (
	"errors"
	"strings"
)

type Validator = func(interface{}) error

func New(description string) (Validator, error) {
	// TODO allow multi constraint
	if description == "" {
		return nil, nil
	}
	validators := []Validator{}
	constraints := strings.Split(description, ";")
	for _, c := range constraints {
		switch {
		case strings.EqualFold(c, "required"):
			validators = append(validators, func(obj interface{}) error {
				if s, ok := obj.(string); ok {
					if len(s) > 0 {
						return nil
					}
				}
				return errors.New("invalid format")
			})
		}
	}
	return func(val interface{}) error {
		for _, validator := range validators {
			if err := validator(val); err != nil {
				return err
			}
		}
		return nil
    }, nil
}
