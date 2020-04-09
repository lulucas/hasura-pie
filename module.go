package pie

import (
	"gopkg.in/asaskevich/govalidator.v9"
	"reflect"
)

type Module interface {
	BeforeCreated(bc BeforeCreatedContext)
	Created(cc CreatedContext)
}

func getModuleName(module Module) string {
	name := reflect.TypeOf(module).Elem().Name()
	return govalidator.CamelCaseToUnderscore(name)
}
