package misc

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func ParseNumber(src any, target any) error {
	switch target.(type) {
	case *int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64,
		*float32, *float64:
		rv := reflect.ValueOf(target)
		if rv.IsNil() {
			return fmt.Errorf("target must be a non nil pointer")
		}
		elem := rv.Elem()
		switch val := src.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64:
			elem.Set(reflect.ValueOf(src).Convert(elem.Type()))
			return nil
		case *int, *int8, *int16, *int32, *int64,
			*uint, *uint8, *uint16, *uint32, *uint64,
			*float32, *float64:
			elem.Set(reflect.ValueOf(src).Elem().Convert(elem.Type()))
		case string:
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				elem.Set(reflect.ValueOf(f).Convert(elem.Type()))
				return nil
			} else {
				return err
			}
		case []byte:
			if f, err := strconv.ParseFloat(string(val), 64); err == nil {
				elem.Set(reflect.ValueOf(f).Convert(elem.Type()))
				return nil
			} else {
				return err
			}
		default:
			return fmt.Errorf("cannot convert type of %T to a number", src)
		}
		return nil
	default:
		return fmt.Errorf("target must be a pointer of number")
	}
}

func ToNumber[T Number](v any) (T, error) {
	var zero T

	switch val := v.(type) {
	case T:
		return val, nil
	default:
		err := ParseNumber(v, &zero)
		return zero, err
	}
}

func IsStrEmptyAndWhitespace(s string) bool {
	return len(s) == 0 || regexp.MustCompile(`^\s*$`).MatchString(s)
}
