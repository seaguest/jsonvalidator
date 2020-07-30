package jsonvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator"
)

const (
	validateTypeInt               = "int"
	validateTypeFloat             = "float"
	validateTypeString            = "string"
	validateTypeRegularExpression = "re"

	typeErrMsg = "tag [%s] expected [%s], but recieved [%+v]"
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
	return validate(tplItf, srcItf, "")
}

func validate(tpl, src interface{}, tag string) (err error) {
	switch val := tpl.(type) {
	case []interface{}:
		// src must be slice
		switch sval := src.(type) {
		case []interface{}:
			if len(val) != len(sval) {
				return fmt.Errorf("node [%s] not match", tag)
			}

			for idx, v := range val {
				err = validate(v, sval[idx], "")
				if err != nil {
					return
				}
			}
		default:
			return fmt.Errorf("node [%s] type inconsistent", tag)
		}
	case map[string]interface{}:
		// src must be slice
		switch sval := src.(type) {
		case map[string]interface{}:
			for k, v := range val {
				err = validate(v, sval[k], k)
				if err != nil {
					return
				}
			}
		default:
			return fmt.Errorf("node [%s] type inconsistent", tag)
		}
	default:
		return validateVar(tag, val.(string), src)
	}
	return
}

func validateVar(tag, exp string, v interface{}) (err error) {
	ss := strings.Split(exp, "|")
	if len(ss) != 2 {
		return fmt.Errorf("invalid exp [%s]", exp)
	}
	typ := ss[0]
	tags := ss[1]

	wrapErr := func(e error) error {
		if e != nil {
			return fmt.Errorf(strings.Replace(e.Error(), "''", fmt.Sprintf("'%s'", tag), -1))
		}
		return nil
	}

	switch typ {
	case validateTypeInt:
		if reflect.TypeOf(v) != reflect.TypeOf(float64(0)) {
			return fmt.Errorf(typeErrMsg, tag, typ, reflect.TypeOf(v))
		}
		return wrapErr(vd.Var(v, tags))
	case validateTypeFloat:
		if reflect.TypeOf(v) != reflect.TypeOf(float64(0)) {
			return fmt.Errorf(typeErrMsg, tag, typ, reflect.TypeOf(v))
		}
		return wrapErr(vd.Var(v, tags))
	case validateTypeString:
		if reflect.TypeOf(v) != reflect.TypeOf("") {
			return fmt.Errorf(typeErrMsg, tag, typ, reflect.TypeOf(v))
		}
		return wrapErr(vd.Var(v, tags))
	case validateTypeRegularExpression:
		exp := regexp.MustCompile(tags)
		macth := exp.Match([]byte(fmt.Sprint(v)))
		if !macth {
			return fmt.Errorf(typeErrMsg, tag, tags, fmt.Sprint(v))
		}
		return
	default:
		return
	}
}
