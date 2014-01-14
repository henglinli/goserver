package message

import (
	"errors"
	"reflect"
)

//
func Call(in interface{}, name string) ([]reflect.Value, error) {
	value := reflect.ValueOf(in)
	if value.IsValid() != true {
		return nil, errors.New("invalid interface")
	}

	function := value.MethodByName(name)
	if function.IsNil() != true {
		return nil, errors.New("funtion not found")
	}

	result := function.Call(nil)
	switch result[0].Kind() {
	case reflect.String:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	default:
		return nil, errors.New("function not return string")
	}
	return result, nil
}

//
func GetString(in []reflect.Value) (string, error) {
	if in[0].Kind() != reflect.String {
		return "", errors.New("Not String")
	}
	return in[0].String(), nil
}

//
func GetUint(in []reflect.Value) (uint64, error) {
	if in[0].Kind() == reflect.String {
		return 0, errors.New("Not uint")
	}
	return in[0].Uint(), nil
}
