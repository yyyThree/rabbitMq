package helper

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"reflect"
	"time"
)

// 基于反射，校验任意值是否为空
func IsEmpty(i interface{}) bool {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func GenerateUuid() string {
	return uuid.NewV4().String()
}

func IsStruct(s interface{}) bool {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// uninitialized zero value of a struct
	if v.Kind() == reflect.Invalid {
		return false
	}

	return v.Kind() == reflect.Struct
}

// 结构体转JSON
func StructToJson(s interface{}) []byte {
	if !IsStruct(s) {
		return nil
	}

	js, err := json.Marshal(s)
	if err != nil {
		return nil
	}
	return js
}

func FormatDateNow() string {
	return time.Now().Format("2006-01-02")
}

func FormatDateNowBySlash() string {
	return time.Now().Format("2006/01/02")
}
