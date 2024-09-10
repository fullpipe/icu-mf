package message

import (
	"fmt"
	"math"
	"strconv"
)

// type Context interface {
// 	String(name string) (string, error)
// 	Any(name string) (any, error)
// }

type Context map[string]any

func (c Context) String(name string) (string, error) {
	v, ok := c[name]
	if !ok {
		return "", fmt.Errorf("argument %s not exists", name)
	}

	return fmt.Sprint(v), nil
}

func (c Context) Int64(key string) (int64, error) {
	v, ok := c[key]
	if !ok {
		return 0, fmt.Errorf("argument %s not exists", key)
	}

	switch i := v.(type) {
	case int:
		return int64(i), nil
	case int8:
		return int64(i), nil
	case int16:
		return int64(i), nil
	case int32:
		return int64(i), nil
	case int64:
		return int64(i), nil
	case uint:
		if i > math.MaxInt64 {
			return 0, fmt.Errorf("unable to convert uint %v to int64 from arg %s", v, key)
		}

		return int64(i), nil //nolint:gosec
	case uint8:
		return int64(i), nil
	case uint16:
		return int64(i), nil
	case uint32:
		return int64(i), nil
	case uint64:
		if i > math.MaxInt64 {
			return 0, fmt.Errorf("unable to convert uint64 %v to int64 from arg %s", v, key)
		}

		return int64(i), nil //nolint:gosec
	case float32:
		return int64(i), nil
	case float64:
		return int64(i), nil
	case string:
		fl, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return 0, err
		}

		return int64(fl), nil
	default:
		return 0, fmt.Errorf("unable to convert %v to int64 from arg %s", v, key)
	}
}

func (c Context) Float64(key string) (float64, error) {
	v, ok := c[key]
	if !ok {
		return 0, fmt.Errorf("argument %s not exists", key)
	}

	switch i := v.(type) {
	case int:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case float64:
		return float64(i), nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case string:
		fl, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return 0, err
		}

		return fl, nil
	default:
		return 0, fmt.Errorf("unable to convert %v to float from arg %s", v, key)
	}
}

func (c Context) Any(name string) (any, error) {
	v, ok := c[name]
	if !ok {
		return "", fmt.Errorf("argument %s not exists", name)
	}

	return v, nil
}

func (c Context) Set(name string, value any) {
	c[name] = value
}

// var _ Context = &context{}
