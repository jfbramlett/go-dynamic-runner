package validator

import (
	"fmt"
)

type ValidatorFactory interface {
	GetValidator(name string) (Validator, error)
}

func getValidatorFactory() ValidatorFactory {
	validatorFactory := &defaultValidatorFactory{}
	validatorFactory.registerValidators()
	return validatorFactory
}


type defaultValidatorFactory struct {
	validators 		map[string]Validator
}


func (v *defaultValidatorFactory) GetValidator(name string) (Validator, error) {
	validator, found := v.validators[name]
	if !found {
		return nil, fmt.Errorf("failed to find validator named %s", name)
	}

	return validator, nil
}


func (v *defaultValidatorFactory) registerValidators() {
	v.validators = make(map[string]Validator)
	v.validators["default"] = &DefaultValidator{}
}


var GlobalValidatorFactory = getValidatorFactory()