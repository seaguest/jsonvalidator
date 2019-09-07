package jsonvalidator

import (
	"encoding/json"
	"fmt"
	"regexp"
)

func ValidateJson(tpl, src []byte) (valid bool, err error) {
	var tplItf interface{}
	if err := json.Unmarshal(tpl, &tplItf); err != nil {
		return false, err
	}

	var srcItf interface{}
	if err := json.Unmarshal(src, &srcItf); err != nil {
		return false, err
	}
	return Validate(tplItf, srcItf)
}

func Validate(tpl, src interface{}) (valid bool, err error) {
	switch v := tpl.(type) {
	case []interface{}:
		return validateSlice(v, src.([]interface{}))
	case map[string]interface{}:
		return validateMap(v, src.(map[string]interface{}))
	default:
		err := fmt.Errorf("non-supported type")
		return false, err
	}
}

// validate the map recursively
func validateMap(tpl, src map[string]interface{}) (valid bool, err error) {
	valid = true
	for k, v := range tpl {
		if subMap, ok := v.(map[string]interface{}); ok {
			var subValid bool
			subValid, err = validateMap(subMap, src[k].(map[string]interface{}))
			if err != nil {
				return false, err
			}
			valid = valid && subValid
		} else if subMap, ok := v.([]interface{}); ok {
			var subValid bool
			subValid, err = validateSlice(subMap, src[k].([]interface{}))
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

// validate the slice recursively
func validateSlice(tpl, src []interface{}) (valid bool, err error) {
	if len(tpl) != len(src) {
		err = fmt.Errorf("size not match expected [%d], but received [%d]", len(tpl), len(src))
		return false, err
	}

	valid = true
	for idx, tv := range tpl {
		if subMap, ok := tv.([]interface{}); ok {
			var subValid bool
			subValid, err = validateSlice(subMap, src[idx].([]interface{}))
			if err != nil {
				return false, err
			}
			valid = valid && subValid
		} else if subMap, ok := tv.(map[string]interface{}); ok {
			var subValid bool
			subValid, err = validateMap(subMap, src[idx].(map[string]interface{}))
			if err != nil {
				return false, err
			}
			valid = valid && subValid
		} else {
			valid = valid && reValidate(tv.(string), src[idx])
			if !valid {
				err = fmt.Errorf("index [%d] expected [%s], but received [%s]", idx, fmt.Sprint(tv), fmt.Sprint(src[idx]))
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
