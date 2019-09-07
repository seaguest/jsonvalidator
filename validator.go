package jsonvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator"
)

var vd *validator.Validate

func init() {
	vd = validator.New()
}

func Validate(tpl, src []byte) (err error) {
	var tplItf interface{}
	if err := json.Unmarshal(tpl, &tplItf); err != nil {
		return err
	}

	var srcItf interface{}
	if err := json.Unmarshal(src, &srcItf); err != nil {
		return err
	}
	return validate(tplItf, srcItf)
}

func validate(tpl, src interface{}) (err error) {
	switch v := tpl.(type) {
	case []interface{}:
		return validateSlice(v, src.([]interface{}))
	case map[string]interface{}:
		return validateMap(v, src.(map[string]interface{}))
	default:
		return fmt.Errorf("non-supported type")
	}
}

func validateVar(id, exp string, v interface{}) (err error) {
	ss := strings.Split(exp, "|")
	if len(ss) != 2 {
		err = fmt.Errorf("invalid exp [%s]", exp)
		return
	}
	typ := ss[0]
	tags := ss[1]

	switch typ {
	case "int":
		if reflect.TypeOf(v) != reflect.TypeOf(float64(0)) {
			err = fmt.Errorf("id [%s] expected [%s], but recieved [%s]", id, typ, reflect.TypeOf(v))
			return
		}

		vv := v.(float64)
		if vv != float64(int64(vv)) {
			err = fmt.Errorf("id [%s] expected [%s], but recieved [%s]", id, typ, reflect.TypeOf(v))
			return
		}
		return vd.Var(vv, tags)
	case "float":
		if reflect.TypeOf(v) != reflect.TypeOf(float64(0)) {
			err = fmt.Errorf("id [%s] expected [%s], but recieved [%s]", id, typ, reflect.TypeOf(v))
			return
		}

		return vd.Var(v, tags)
	case "string":
		if reflect.TypeOf(v) != reflect.TypeOf("") {
			err = fmt.Errorf("id [%s] expected [%s], but recieved [%s]", id, typ, reflect.TypeOf(v))
			return
		}
		return vd.Var(v, tags)
	case "re":
		exp := regexp.MustCompile(tags)
		macth := exp.Match([]byte(fmt.Sprint(v)))
		if !macth {
			return fmt.Errorf("id [%s] expected [%s], but received [%s]", id, fmt.Sprint(exp), fmt.Sprint(v))
		}
		return nil
	default:
		return nil
	}
}

// vd the map recursively
func validateMap(tpl, src map[string]interface{}) (err error) {
	for k, v := range tpl {
		if subMap, ok := v.(map[string]interface{}); ok {
			err = validateMap(subMap, src[k].(map[string]interface{}))
			if err != nil {
				return
			}
		} else if subMap, ok := v.([]interface{}); ok {
			err = validateSlice(subMap, src[k].([]interface{}))
			if err != nil {
				return
			}
		} else {
			err = validateVar("", v.(string), src[k])
			if err != nil {
				return
			}
		}
	}
	return
}

// vd the map slice recursively
func validateSlice(tpl, src []interface{}) (err error) {
	if len(tpl) != len(src) {
		err = fmt.Errorf("size not match expected [%d], but received [%d]", len(tpl), len(src))
		return
	}

	for idx, tv := range tpl {
		if subMap, ok := tv.([]interface{}); ok {
			err = validateSlice(subMap, src[idx].([]interface{}))
			if err != nil {
				return
			}
		} else if subMap, ok := tv.(map[string]interface{}); ok {
			err = validateMap(subMap, src[idx].(map[string]interface{}))
			if err != nil {
				return
			}
		} else {
			err = validateVar("", tv.(string), src[idx])
			if err != nil {
				return
			}
		}
	}
	return
}
