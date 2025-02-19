package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Parse(ptr interface{}) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}
	return set(ptr, "default")
}

func set(ptr interface{}, tag string) error {

	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		environmentVar := os.Getenv(t.Field(i).Tag.Get("env"))
		if environmentVar != "" {
			if err := setField(v.Field(i), environmentVar); err != nil {
				return err
			}
		} else if defaultVal := t.Field(i).Tag.Get(tag); defaultVal != "-" {
			if err := setField(v.Field(i), defaultVal); err != nil {
				return err
			}
		}
	}
	return nil
}

func setField(field reflect.Value, defaultVal string) error {

	if !field.CanSet() {
		return fmt.Errorf("can't set value")
	}

	switch field.Kind() {
	case reflect.Bool:
		if val, err := strconv.ParseBool(defaultVal); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val, err := strconv.ParseInt(defaultVal, 10, 64); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val, err := strconv.ParseUint(defaultVal, 10, 64); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
		}
	case reflect.Float32, reflect.Float64:
		if val, err := strconv.ParseFloat(defaultVal, 64); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
		}
	case reflect.String:
		field.Set(reflect.ValueOf(defaultVal).Convert(field.Type()))
	case reflect.Struct:
		if field.Type().String() == "time.Time" {
			if val, err := time.Parse(time.RFC3339, defaultVal); err == nil {
				field.Set(reflect.ValueOf(val).Convert(field.Type()))
			}
		}
	case reflect.Slice:
		fmt.Println(field.Type().String())
		if field.Type().String() == "[]uint8" {
			field.Set(reflect.ValueOf([]byte(defaultVal)).Convert(field.Type()))
		} else if field.Type().String() == "[]string" {
			field.Set(reflect.ValueOf(strings.Split(defaultVal, ",")).Convert(field.Type()))
		} else if field.Type().String() == "[]int" {
			arr := strings.Split(defaultVal, ",")
			var arrInt []int
			for _, v := range arr {
				if val, err := strconv.Atoi(v); err == nil {
					arrInt = append(arrInt, val)
				}
			}
			field.Set(reflect.ValueOf(arrInt).Convert(field.Type()))
		}
	default:
		return fmt.Errorf("unsupported type: %s", field.Kind())
	}

	return nil
}
