package utils

import "errors"

// Required 判断是否为空
func Required(m map[string]interface{}, r map[string]string) error {
	var s string
	for k, v := range r {
		if _, ok := m[k]; !ok {
			s += v + " "
		}
	}
	if s != "" {
		return errors.New(s)
	}
	return nil
}
