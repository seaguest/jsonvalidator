package jsonvalidator

import (
	"encoding/json"
	"fmt"
	"regexp"
)

func ValidateJson(tpl, input []byte) (valid bool, err error) {
	var tmpl map[string]interface{}
	if err := json.Unmarshal(tpl, &tmpl); err != nil {
		return false, err
	}

	var src map[string]interface{}
	if err := json.Unmarshal(input, &src); err != nil {
		return false, err
	}
	return ValidateMap(tmpl, src)
}

// ValidateMap the map recursively
func ValidateMap(tmpl, src map[string]interface{}) (valid bool, err error) {
	valid = true
	for k, v := range tmpl {
		if innerMap, ok := v.(map[string]interface{}); ok {
			var subValid bool
			subValid, err = ValidateMap(innerMap, src[k].(map[string]interface{}))
			if err != nil {
				return false, err
			}
			valid = valid && subValid
		} else {
			valid = valid && reValidate(v.(string), src[k])
			if !valid {
				err = fmt.Errorf("key [%s] expected [%s], but received [%s]", k, fmt.Sprint(v), fmt.Sprint(src[k]))
				return false, err
			}
		}
		if !valid {
			return
		}
	}
	return
}

func reValidate(reg string, v interface{}) bool {
	exp := regexp.MustCompile(reg)
	return exp.Match([]byte(fmt.Sprint(v)))
}
